package mcp

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestSessionIdleSweepReleasesAndStops comprova limpeza sem nova chamada do cliente e término da goroutine.
func TestSessionIdleSweepReleasesAndStops(t *testing.T) {
	// seconds permite avançar o relógio da sessão com leitura concorrente segura pela varredura.
	var seconds atomic.Int64
	seconds.Store(1800000000)
	core := application.NewAuthoring(application.DefaultAuthoringLimits(), func() time.Time { return time.Unix(seconds.Load(), 0) })
	info, err := core.Admit(project.NewID(), "sweep-fixture", application.NewDocument("Ociosa"), application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	seconds.Add(31 * 60)
	// ctx/done asseguram encerramento mesmo se alguma asserção falhar.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go sweepIdle(ctx, core, time.Millisecond, done)
	defer func() { cancel(); <-done }()
	// deadline/poll observam somente quotas; consultar a sessão alteraria seu prazo.
	deadline, poll := time.NewTimer(time.Second), time.NewTicker(time.Millisecond)
	defer deadline.Stop()
	defer poll.Stop()
	for core.Usage().StagingReservedBytes != 0 {
		select {
		case <-deadline.C:
			t.Fatal("varredura não liberou sessão ociosa")
		case <-poll.C:
		}
	}
	if _, err := core.SessionStatus(info.ID); err == nil {
		t.Fatal("sessão ociosa foi reaberta")
	}
}

// sessionObservation verifica o contrato público sem consultar estado privado do adaptador.
type sessionObservation struct {
	// SessionID/Revision/ControlGeneration preservam identidades e contadores decimais.
	SessionID, Revision, ControlGeneration string
	// Dirty distingue trabalho de documento e de proposta.
	Dirty, DocumentDirty, ProposalDirty bool
	// Quotas permite conferir a liberação de reservas sem observar arquivos privados.
	Quotas struct {
		EditRequests, ReservedTargets, Proposals int
		StagingReservedBytes                     int64
	}
}

// sessionControlArgs obtém o CAS e o high-water pela mesma ferramenta disponível ao cliente.
func sessionControlArgs(t *testing.T, client *sdk.ClientSession, id string) map[string]any {
	t.Helper()
	// current é o estado público que vincula a próxima intenção de controle.
	current := callResult[sessionObservation](t, client, "sessao_obter", map[string]any{"sessionId": id})
	return map[string]any{"sessionId": id, "requestId": project.NewID(), "expectedRevision": current.Revision, "controlGeneration": current.ControlGeneration}
}

// TestSessionControlsAggregateDraftsAndReplay protege anexos e mantém resposta terminal sem ressuscitar a sessão.
func TestSessionControlsAggregateDraftsAndReplay(t *testing.T) {
	// root/runtime/client usam pasta vazia para comprovar ausência de persistência dos controles.
	root := t.TempDir()
	runtime, client := testClientAt(t, root)
	// info admite base conhecida limpa, isolando o dirty produzido apenas pelo anexo.
	info, err := runtime.authoring.Admit(project.NewID(), "session-fixture", application.NewDocument("Sessão"), application.AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.authoring.StoreDraft(info.ID, project.NewID(), "proposal-fixture", info.Revision, application.AuthoringDraft{ID: project.NewID(), Revision: 1, Origin: "selection", Data: []byte(`{"name":"Composição"}`), Target: "library/example"}); err != nil {
		t.Fatal(err)
	}
	// before mostra proposta pendente sem alterar revisão ou dirty do documento.
	before := callResult[sessionObservation](t, client, "sessao_obter", map[string]any{"sessionId": info.ID})
	if before.DocumentDirty || !before.ProposalDirty || !before.Dirty || before.Quotas.ReservedTargets != 1 || before.Quotas.Proposals != 1 {
		t.Fatalf("dirty/reservas incorretos: %+v", before)
	}
	// refused mantém todo o trabalho quando o cliente não informa descarte explícito.
	refused := sessionControlArgs(t, client, info.ID)
	refused["discardUnsaved"] = false
	callFailure(t, client, "sessao_encerrar", refused, "invalidArgument")
	// discard conserva resposta idempotente, remove propostas e libera staging sem gravar.
	discard := sessionControlArgs(t, client, info.ID)
	first := callResult[map[string]any](t, client, "sessao_descartar", discard)
	replay := callResult[map[string]any](t, client, "sessao_descartar", discard)
	if !reflect.DeepEqual(first, replay) || first["discarded"] != true {
		t.Fatalf("replay de descarte divergente: %v / %v", first, replay)
	}
	// after verifica que o high-water avançou e a revisão antiga não voltou a ser válida.
	after := callResult[sessionObservation](t, client, "sessao_obter", map[string]any{"sessionId": info.ID})
	if after.Dirty || after.Quotas.ReservedTargets != 0 || after.Quotas.Proposals != 0 || after.Quotas.StagingReservedBytes != 0 || after.Revision == before.Revision {
		t.Fatalf("descarte incompleto: %+v", after)
	}
	// closeArgs encerra a sessão limpa; somente o replay terminal conhecido é permitido.
	closeArgs := sessionControlArgs(t, client, info.ID)
	closeArgs["discardUnsaved"] = false
	closed := callResult[map[string]any](t, client, "sessao_encerrar", closeArgs)
	closedReplay := callResult[map[string]any](t, client, "sessao_encerrar", closeArgs)
	if !reflect.DeepEqual(closed, closedReplay) || closed["closed"] != true {
		t.Fatal("fechamento não conservou resposta terminal")
	}
	callFailure(t, client, "sessao_obter", map[string]any{"sessionId": info.ID}, "sessionExpired")
	closeArgs["requestId"] = project.NewID()
	callFailure(t, client, "sessao_encerrar", closeArgs, "sessionExpired")
	// entries observa a raiz autorizada, que continua sem nenhum arquivo de jornada.
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 || runtime.authoring.Usage().StagingReservedBytes != 0 {
		t.Fatalf("controle produziu arquivo ou reteve reserva: %v (%v)", entries, err)
	}
}

// TestSessionControlsRemainAvailableAtEditLimit não permite que quota de conteúdo impeça descartar/encerrar.
func TestSessionControlsRemainAvailableAtEditLimit(t *testing.T) {
	// limits mantém as políticas reais exceto quota pequena para atingir a fronteira via protocolo.
	limits := application.DefaultAuthoringLimits()
	limits.EditRequests = 1
	runtime, client := testClient(t)
	runtime.authoring.CloseAll()
	runtime.authoring = application.NewAuthoring(limits, nil)
	// info cria rascunho que já precisa de descarte informado no fechamento.
	info, err := runtime.authoring.Admit(project.NewID(), "small-quota", application.NewDocument("Quota"), application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	// edit usa um comando real e esgota a única identidade admitida.
	edit := map[string]any{"sessionId": info.ID, "requestId": project.NewID(), "expectedRevision": strconv.FormatUint(info.Revision, 10), "name": "Outra tela"}
	callResult[struct{}](t, client, "tela_criar", edit)
	edit["requestId"] = project.NewID()
	callFailure(t, client, "tela_criar", edit, "resourceLimit")
	// args continua disponível mesmo quando novas edições são recusadas.
	args := sessionControlArgs(t, client, info.ID)
	callResult[struct{}](t, client, "sessao_descartar", args)
	args = sessionControlArgs(t, client, info.ID)
	args["discardUnsaved"] = true
	callResult[struct{}](t, client, "sessao_encerrar", args)
}

// TestSessionExpiredControlsCannotRevive mantém o prazo absoluto após inatividade, inclusive em controle atrasado.
func TestSessionExpiredControlsCannotRevive(t *testing.T) {
	// now é avançado sem sleep para reproduzir expiração real do registro.
	now := time.Unix(1800000000, 0)
	runtime, client := testClient(t)
	runtime.authoring.CloseAll()
	runtime.authoring = application.NewAuthoring(application.DefaultAuthoringLimits(), func() time.Time { return now })
	// info/args conservam as precondições válidas antes da expiração.
	info, err := runtime.authoring.Admit(project.NewID(), "expiry-fixture", application.NewDocument("Ociosa"), application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	args := sessionControlArgs(t, client, info.ID)
	now = now.Add(31 * time.Minute)
	runtime.authoring.ExpireIdle()
	callFailure(t, client, "sessao_descartar", args, "sessionExpired")
	args["discardUnsaved"] = true
	callFailure(t, client, "sessao_encerrar", args, "sessionExpired")
	// capabilities continua consultável e não constitui criação implícita de sessão.
	callResult[installedCapabilities](t, client, "aplicacao_capacidades", map[string]any{})
	if runtime.authoring.Usage().StagingReservedBytes != 0 {
		t.Fatal("sessão expirada reteve staging")
	}
}
