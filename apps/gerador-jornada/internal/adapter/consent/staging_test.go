//go:build windows

package consent

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/adapter/workspace"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
	"github.com/google/jsonschema-go/jsonschema"
)

// TestFreezerFollowsSessionReviewLifecycle integra o coordenador ao ZIP privado real, sem host ou gravação final.
func TestFreezerFollowsSessionReviewLifecycle(t *testing.T) {
	// root/private separam jornada e artefato; referências do harness não qualificam autorização.
	root, private := t.TempDir(), filepath.Join(t.TempDir(), "stage")
	freezer, err := NewFreezer(private, root, ReviewProfile{HostID: "harness", KeyID: "key", ServerInstanceID: project.NewID()})
	if err != nil {
		t.Fatal(err)
	}
	core := application.NewAuthoring(application.DefaultAuthoringLimits(), nil)
	defer core.CloseAll()
	info, err := core.Admit(project.NewID(), "lifecycle-fixture", workspace.NewJourneyDocument("Integrada 😀", false), application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	// request/options conservam o mesmo pedido durante replay, incluindo o CAS inicial.
	request := project.NewID()
	options := application.AuthoringReviewOptions{Filename: "integrada.jornada", Intent: "create", AllowDraft: true}
	first, err := core.PrepareReview(context.Background(), info.ID, request, "prepare", info.Revision, 0, options, freezer)
	if err != nil {
		t.Fatal(err)
	}
	testReviewSchema(t, first.Bytes)
	replay, err := core.PrepareReview(context.Background(), info.ID, request, "prepare", info.Revision, 0, options, freezer)
	if err != nil || !reflect.DeepEqual(first, replay) {
		t.Fatalf("replay regenerou revisão/artefato: %v", err)
	}
	// entries deve conter somente o arquivo original; replay não cria outro staging.
	entries, err := os.ReadDir(private)
	if err != nil || len(entries) != 1 {
		t.Fatal("replay criou artefato adicional")
	}
	// second aposenta fisicamente o primeiro artefato antes de disponibilizar a nova revisão.
	second, err := core.PrepareReview(context.Background(), info.ID, project.NewID(), "reprepare", info.Revision, first.Ticket.Generation, options, freezer)
	if err != nil || second.Hash == first.Hash || second.ArtifactID == first.ArtifactID {
		t.Fatalf("nova preparação inválida: %v", err)
	}
	if err := freezer.WriteArtifact(context.Background(), first.ArtifactID, &bytes.Buffer{}); err == nil {
		t.Fatal("artefato aposentado continuou disponível")
	}
	if _, err := core.CurrentReview(info.ID, first.Ticket.ReviewID); err == nil {
		t.Fatal("revisão antiga continuou corrente")
	}
	// discard libera a única referência publicável e sua reserva; nenhum destino final foi criado.
	if _, err := core.Control(info.ID, project.NewID(), "discard", info.Revision, second.Ticket.Generation, "discard"); err != nil {
		t.Fatal(err)
	}
	if err := freezer.WriteArtifact(context.Background(), second.ArtifactID, &bytes.Buffer{}); err == nil {
		t.Fatal("descarte conservou artefato")
	}
	entries, err = os.ReadDir(private)
	if err != nil || len(entries) != 0 {
		t.Fatal("descarte não removeu staging")
	}
	entries, err = os.ReadDir(root)
	if err != nil || len(entries) != 0 || core.Usage().StagingReservedBytes != 0 {
		t.Fatal("revisão publicou destino ou conservou cota")
	}
}

// freezeRequest cria preparação sintética técnica, sem host humano qualificado ou assinatura.
func freezeRequest() application.AuthoringFreezeRequest {
	// document/ticket fornecem as mesmas identidades usadas pelo coordenador real.
	document := workspace.NewJourneyDocument("Staging", false)
	ticket := application.AuthoringControlTicket{Generation: 1, ReviewID: project.NewID(), CommitRequestID: project.NewID(), ExpiresAt: time.Now().Add(5 * time.Minute).Unix()}
	return application.AuthoringFreezeRequest{SessionID: project.NewID(), Ticket: ticket, Document: document, BaselineDocument: document.Copy(), Revision: document.Project.Revision, Options: application.AuthoringReviewOptions{Filename: "new.jornada", Intent: "create", AllowDraft: true}, Summary: []string{"Criar jornada em destino ausente; revisão técnica sem autorização."}}
}

// testReviewSchema valida o objeto interno Review e suas referências locais no schema normativo.
func testReviewSchema(t *testing.T, data []byte) {
	t.Helper()
	// path referencia a especificação do repositório; não é recurso exigido no computador final.
	path := filepath.Join("..", "..", "..", "..", "..", "specs", "038-integracao-mcp-agentes", "contracts", "authorization.schema.json")
	encoded, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var schema jsonschema.Schema
	if err := json.Unmarshal(encoded, &schema); err != nil {
		t.Fatal(err)
	}
	schema.OneOf, schema.Ref = nil, "#/$defs/Review"
	resolved, err := schema.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	var instance any
	if err := json.Unmarshal(data, &instance); err != nil {
		t.Fatal(err)
	}
	if err := resolved.Validate(instance); err != nil {
		t.Fatalf("revisão fora do schema: %v", err)
	}
}

// TestFreezerCreatesCanonicalImmutableReview testa ZIP real, hash/revisão e liberação sem gravar jornada.
func TestFreezerCreatesCanonicalImmutableReview(t *testing.T) {
	// root/private separam dados do autor e staging operacional; profile é apenas referência de harness.
	root, private := t.TempDir(), filepath.Join(t.TempDir(), "private")
	profile := ReviewProfile{HostID: "harness", KeyID: "harness-key", ServerInstanceID: project.NewID()}
	freezer, err := NewFreezer(private, root, profile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(private); !os.IsNotExist(err) {
		t.Fatal("construtor criou diretório antes da preparação")
	}
	request := freezeRequest()
	frozen, err := freezer.Freeze(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	defer freezer.Release(frozen.ArtifactID)
	testReviewSchema(t, frozen.Bytes)
	// data é lido do mesmo handle privado, sem reabrir caminho ou aceitar pacote do cliente.
	var data bytes.Buffer
	if err := freezer.WriteArtifact(context.Background(), frozen.ArtifactID, &data); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data.Bytes())
	if hex.EncodeToString(digest[:]) != frozen.ContentHash || int64(data.Len()) != frozen.ContentSizeBytes {
		t.Fatal("hash/tamanho não correspondem aos bytes congelados")
	}
	decoded, err := projectzip.Decode(data.Bytes())
	if err != nil || decoded.Project.ID != request.Document.Project.ID {
		t.Fatal("artefato não preservou o documento")
	}
	var review reviewDTO
	if err := json.Unmarshal(frozen.Bytes, &review); err != nil {
		t.Fatal(err)
	}
	if review.RequestID != request.Ticket.CommitRequestID || review.ReviewID != request.Ticket.ReviewID || review.HostID != profile.HostID || review.BaseExists || review.BaseHash != nil || review.ContentHash != frozen.ContentHash {
		t.Fatal("revisão não vinculou identidades/base/conteúdo")
	}
	// order conserva a ordem required, evitando reconstrução ambígua no contrato assinado.
	previous := -1
	for _, field := range []string{"profileVersion", "hostId", "keyId", "serverInstanceId", "sessionId", "requestId", "reviewId", "operation", "targetIdentity", "baseExists", "baseHash", "contentHash", "expectedRevision", "controlGeneration", "proposal", "materialization", "expiresAt", "challenge", "summary", "diagnostics"} {
		position := bytes.Index(frozen.Bytes, []byte(`"`+field+`":`))
		if position <= previous {
			t.Fatal("ordem canônica incorreta:", field)
		}
		previous = position
	}
	if err := freezer.Release(frozen.ArtifactID); err != nil {
		t.Fatal(err)
	}
	if err := freezer.Release(frozen.ArtifactID); err != nil {
		t.Fatal("liberação não idempotente")
	}
	entries, err := os.ReadDir(private)
	if err != nil || len(entries) != 0 {
		t.Fatal("artefato sobreviveu à liberação")
	}
	entries, err = os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("preparação publicou arquivo no workspace")
	}
}

// TestFreezerRejectsWorkspacePrivateDirectory recusa escopo inválido antes de criar diretórios.
func TestFreezerRejectsWorkspacePrivateDirectory(t *testing.T) {
	// private ainda não existe; validação deve usar seu ancestral físico, sem escrever no workspace.
	root := t.TempDir()
	private := filepath.Join(root, "not-private", "stage")
	_, err := NewFreezer(private, root, ReviewProfile{HostID: "harness", KeyID: "key", ServerInstanceID: project.NewID()})
	var typed *application.AuthoringError
	if !errors.As(err, &typed) || typed.Code != "permissionDenied" {
		t.Fatalf("raiz privada dentro do workspace aceita: %v", err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("validação de escopo causou escrita")
	}
}

// TestFreezerRejectsChangedBaseAndCleansFailedReview não congela base obsoleta nem conserva ZIP de revisão inválida.
func TestFreezerRejectsChangedBaseAndCleansFailedReview(t *testing.T) {
	// original é a origem escolhida explicitamente; outro escritor a muda depois da leitura.
	root, private := t.TempDir(), filepath.Join(t.TempDir(), "private")
	original := filepath.Join(root, "base.jornada")
	document := workspace.NewJourneyDocument("Original", false)
	if err := projectzip.Save(context.Background(), original, document); err != nil {
		t.Fatal(err)
	}
	store, err := workspace.OpenWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	loaded, baseline, err := store.OpenDocument(context.Background(), "base.jornada")
	if err != nil {
		t.Fatal(err)
	}
	freezer, err := NewFreezer(private, root, ReviewProfile{HostID: "harness", KeyID: "key", ServerInstanceID: project.NewID()})
	if err != nil {
		t.Fatal(err)
	}
	document.Project.Name = "Alterado por terceiro"
	if err := projectzip.Save(context.Background(), original, document); err != nil {
		t.Fatal(err)
	}
	request := freezeRequest()
	request.Document, request.Baseline = loaded, baseline
	request.Options.Filename, request.Options.Intent = "base.jornada", "replace"
	_, err = freezer.Freeze(context.Background(), request)
	var typed *application.AuthoringError
	if !errors.As(err, &typed) || typed.Code != "externalModification" {
		t.Fatalf("base obsoleta aceita: %v", err)
	}
	if _, err := os.Stat(private); !os.IsNotExist(err) {
		t.Fatal("base obsoleta gerou staging")
	}
	// summary inválido é recusado depois de codificar/sincronizar o ZIP, exigindo rollback real do temporário.
	request = freezeRequest()
	request.Summary = []string{strings.Repeat("x", 4097)}
	if _, err := freezer.Freeze(context.Background(), request); err == nil {
		t.Fatal("revisão fora do schema aceita")
	}
	entries, err := os.ReadDir(private)
	if err != nil || len(entries) != 0 {
		t.Fatal("falha depois da codificação deixou pacote privado")
	}
	if len(freezer.artifacts) != 0 {
		t.Fatal("falha conservou referência publicável")
	}
	current, err := projectzip.Load(original)
	if err != nil || current.Project.Name != "Alterado por terceiro" {
		t.Fatal("preparação restaurou/sobrescreveu base de terceiro")
	}
}
