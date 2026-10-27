package application

import (
	"context"
	"crypto/sha256"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

// reviewFreezer simula somente a fronteira de artefato; não concede consentimento ou escreve jornada.
type reviewFreezer struct {
	// calls/releases verificam replay e descarte de artefatos aposentados.
	calls, releases int
	// payload é revisão sintética opaca, usada apenas para testar coordenação e cópias.
	payload []byte
	// inspect recebe a visão protegida para verificar ausência de clone integral de recursos.
	inspect func(AuthoringFreezeRequest)
}

// blockedReviewRelease segura somente IO de limpeza para provar ausência de mutex global retido.
type blockedReviewRelease struct {
	reviewFreezer
	// entered/finish sincronizam a limpeza sem sleeps; once sinaliza somente a primeira tentativa.
	entered, finish chan struct{}
	once            sync.Once
}

// Release espera liberação explícita apenas na fixture; produção fecha handle nativo local.
func (freezer *blockedReviewRelease) Release(id string) error {
	freezer.once.Do(func() { close(freezer.entered) })
	<-freezer.finish
	return freezer.reviewFreezer.Release(id)
}

// TestReviewRetirementKeepsBudgetUntilPhysicalRelease não libera cota nem bloqueia consultas durante limpeza.
func TestReviewRetirementKeepsBudgetUntilPhysicalRelease(t *testing.T) {
	// limits deixa uma única reserva física, tornando visível qualquer liberação antecipada.
	limits := DefaultAuthoringLimits()
	limits.StagingBytes = limits.TargetStagingBytes
	core := NewAuthoring(limits, nil)
	first, err := core.Admit("00000000-0000-4000-8000-000000000001", "first", NewDocument("Primeira"), AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	other, err := core.Admit("00000000-0000-4000-8000-000000000002", "other", NewDocument("Consulta"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	freezer := &blockedReviewRelease{reviewFreezer: reviewFreezer{payload: []byte(`{"review":true}`)}, entered: make(chan struct{}), finish: make(chan struct{})}
	defer func() {
		select {
		case <-freezer.finish:
		default:
			close(freezer.finish)
		}
	}()
	prepared, err := core.PrepareReview(context.Background(), first.ID, "10000000-0000-4000-8000-000000000001", "prepare", first.Revision, 0, AuthoringReviewOptions{Filename: "new.jornada", Intent: "create", AllowDraft: true}, freezer)
	if err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() {
		_, err := core.Control(first.ID, "20000000-0000-4000-8000-000000000001", "discard", first.Revision, prepared.Ticket.Generation, "discard")
		finished <- err
	}()
	<-freezer.entered
	// queried deve terminar embora outra chamada esteja aguardando liberação física.
	queried := make(chan error, 1)
	go func() { _, err := core.InspectMetadata(other.ID); queried <- err }()
	select {
	case err := <-queried:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("IO de limpeza reteve a trava global")
	}
	if usage := core.Usage(); usage.PendingCleanups != 1 || usage.StagingReservedBytes != limits.TargetStagingBytes {
		t.Fatal("reserva física liberada antes do handle")
	}
	if _, err := core.Admit("00000000-0000-4000-8000-000000000003", "blocked", NewDocument("Bloqueada"), AuthoringBaseline{}); authoringCode(err) != "resourceLimit" {
		t.Fatal("nova admissão reutilizou staging ainda aberto")
	}
	close(freezer.finish)
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
	if usage := core.Usage(); usage.PendingCleanups != 0 || usage.StagingReservedBytes != 0 {
		t.Fatal("limpeza concluída não liberou reserva")
	}
}

// TestIdleExpirationPreservesUncertainIntent conserva a posição original até a reconciliação, sem renovar prazo.
func TestIdleExpirationPreservesUncertainIntent(t *testing.T) {
	// now permite atravessar o limite sem aguardar meia hora no teste.
	now := time.Unix(1800000000, 0)
	core := NewAuthoring(DefaultAuthoringLimits(), func() time.Time { return now })
	info, err := core.Admit("00000000-0000-4000-8000-000000000001", "new", NewDocument("Incerta"), AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	commitID := "10000000-0000-4000-8000-000000000001"
	if err := core.BeginIntent(info.ID, "document", commitID); err != nil {
		t.Fatal(err)
	}
	now = now.Add(31 * time.Minute)
	core.ExpireIdle()
	snapshot, err := core.InspectMetadata(info.ID)
	if err != nil || !snapshot.RecoveryRequired || snapshot.CommitPositions != 1 || snapshot.UncertainIntents != 1 || !snapshot.ExpiresAt.Before(now) {
		t.Fatal("inatividade apagou ou renovou intenção incerta")
	}
	if err := core.ResolveIntent(info.ID, commitID, false); err != nil {
		t.Fatal(err)
	}
	core.ExpireIdle()
	if _, err := core.InspectMetadata(info.ID); authoringCode(err) != "sessionExpired" {
		t.Fatal("sessão foi ressuscitada pela recuperação")
	}
}

// Freeze devolve bytes imutáveis já produzidos pelo adaptador, sem fonte de autorização.
func (freezer *reviewFreezer) Freeze(ctx context.Context, request AuthoringFreezeRequest) (AuthoringFrozenReview, error) {
	freezer.calls++
	if freezer.inspect != nil {
		freezer.inspect(request)
	}
	return AuthoringFrozenReview{ArtifactID: fmt.Sprintf("artifact-%d", freezer.calls), ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte("frozen-content"))), ContentSizeBytes: 14, Bytes: append([]byte{}, freezer.payload...)}, nil
}

// Release aposenta artefato, sem depender de confirmação humana ou nova vaga de edição.
func (freezer *reviewFreezer) Release(artifactID string) error { freezer.releases++; return nil }

// TestReviewReplayRetainsExactBytes conserva validade, identidade e dados protegidos sem nova codificação.
func TestReviewReplayRetainsExactBytes(t *testing.T) {
	// now/core usam tempo controlado para distinguir replay de renovação de validade.
	now := time.Unix(1800000000, 0)
	core := NewAuthoring(DefaultAuthoringLimits(), func() time.Time { return now })
	document := NewDocument("Revisão")
	document.Resources["test"] = []byte{1, 2, 3}
	info, err := core.Admit("00000000-0000-4000-8000-000000000001", "new", document, AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	freezer := &reviewFreezer{payload: []byte(`{"syntheticReview":true}`)}
	freezer.inspect = func(request AuthoringFreezeRequest) {
		// address compara os bytes já retidos; Freeze não recebe outro documento integral clonado.
		if &request.Document.Resources["test"][0] != &core.sessions[info.ID].model.Document.Resources["test"][0] {
			t.Fatal("preparação clonou recursos")
		}
	}
	requestID := "10000000-0000-4000-8000-000000000001"
	options := AuthoringReviewOptions{Filename: "nova.jornada", Intent: "create", AllowDraft: true}
	first, err := core.PrepareReview(context.Background(), info.ID, requestID, "review-1", info.Revision, 0, options, freezer)
	if err != nil {
		t.Fatal(err)
	}
	if first.Ticket.ReviewID == requestID || first.Ticket.CommitRequestID == requestID || first.Ticket.CommitRequestID == first.Ticket.ReviewID {
		t.Fatal("identidades de preparação/revisão/commit misturadas")
	}
	first.Bytes[0] = 'x'
	now = now.Add(time.Minute)
	replayed, err := core.PrepareReview(context.Background(), info.ID, requestID, "review-1", info.Revision, 0, options, freezer)
	if err != nil || freezer.calls != 1 || replayed.Bytes[0] != '{' || replayed.Ticket.ExpiresAt != now.Add(4*time.Minute).Unix() {
		t.Fatal("replay mudou bytes/validade ou recodificou")
	}
	// discard aposenta o único artefato e libera reservas antes de retornar.
	snapshot, _ := core.InspectMetadata(info.ID)
	if _, err := core.Control(info.ID, "20000000-0000-4000-8000-000000000001", "discard", snapshot.Revision, snapshot.ControlGeneration, "discard"); err != nil {
		t.Fatal(err)
	}
	if freezer.releases != 1 {
		t.Fatal("descarte não liberou staging")
	}
}

// TestReviewPreparationSurvivesEditQuota não depende de nova posição de edição nem altera o documento.
func TestReviewPreparationSurvivesEditQuota(t *testing.T) {
	// limits reduz apenas IDs de edição; preparação deve usar a reserva já admitida.
	limits := DefaultAuthoringLimits()
	limits.EditRequests = 1
	core := NewAuthoring(limits, nil)
	info, err := core.Admit("00000000-0000-4000-8000-000000000001", "new", NewDocument("Inicial"), AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	edit, err := core.Edit(info.ID, "10000000-0000-4000-8000-000000000001", "edit", info.Revision, func(session *Session) (map[string]string, error) {
		return nil, session.Change(func(document *Document) error { document.Project.Name = "Alterada"; return nil })
	})
	if err != nil {
		t.Fatal(err)
	}
	before, _ := core.Inspect(info.ID)
	freezer := &reviewFreezer{payload: []byte(`{"syntheticReview":true}`)}
	first, err := core.PrepareReview(context.Background(), info.ID, "20000000-0000-4000-8000-000000000001", "prepare", edit.Revision, before.ControlGeneration, AuthoringReviewOptions{Filename: "new.jornada", Intent: "create", AllowDraft: true}, freezer)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := core.Inspect(info.ID)
	if before.EditRequests != after.EditRequests || before.Revision != after.Revision || !reflect.DeepEqual(before.Document, after.Document) || after.CommitPositions != 1 {
		t.Fatal("preparação consumiu edição/mudou documento ou duplicou reserva")
	}
	// second substitui a revisão anterior e fecha seu artefato antes de abrir outro.
	second, err := core.PrepareReview(context.Background(), info.ID, "20000000-0000-4000-8000-000000000002", "prepare-2", edit.Revision, after.ControlGeneration, AuthoringReviewOptions{Filename: "new.jornada", Intent: "create", AllowDraft: true}, freezer)
	if err != nil || first.Ticket.ReviewID == second.Ticket.ReviewID || freezer.releases != 1 {
		t.Fatal("repreparo acumulou ou reutilizou revisão")
	}
	core.CloseAll()
	if freezer.releases != 2 {
		t.Fatal("EOF não liberou revisão corrente")
	}
}
