package application

import (
	"context"
	"errors"
	"testing"
)

// memoryDocumentSource simula IO para provar admissão idempotente sem Fyne ou filesystem.
type memoryDocumentSource struct {
	// document/baseline representam a leitura congelada; failure simula origem removida.
	document Document
	baseline AuthoringBaseline
	failure  error
	// reads permite detectar reabertura indevida depois da primeira admissão.
	reads int
}

// ListJourneys fornece inventário sintético sem percorrer a fonte de documentos.
func (source *memoryDocumentSource) ListJourneys(context.Context) (AuthoringInventory, error) {
	return AuthoringInventory{Revision: "inventory-1", Entries: []AuthoringJourney{{Filename: "base.jornada", Name: "Base"}}}, nil
}

// OpenDocument conserva a relação entre bytes decodificados e sua precondição.
func (source *memoryDocumentSource) OpenDocument(_ context.Context, name string) (Document, AuthoringBaseline, error) {
	source.reads++
	return source.document, source.baseline, source.failure
}

// NewDocument representa o construtor comum que o adaptador injeta nos dois canais.
func (source *memoryDocumentSource) NewDocument(name string, example bool) Document {
	return NewDocument(name)
}

// TestAuthoringWorkspaceAdmissionReplay não volta a ler uma origem que já foi admitida.
func TestAuthoringWorkspaceAdmissionReplay(t *testing.T) {
	// source simula uma abertura cuja origem desaparece depois do sucesso inicial.
	source := &memoryDocumentSource{document: NewDocument("Base"), baseline: AuthoringBaseline{Exists: true, Identity: "file-1", SHA256: "hash-1", RootIdentity: "root-1", Filename: "base.jornada"}}
	core := NewAuthoring(DefaultAuthoringLimits(), nil)
	workspace := NewAuthoringWorkspace(core, source)
	requestID := "00000000-0000-4000-8000-000000000001"
	first, err := workspace.Open(context.Background(), requestID, "digest-1", "base.jornada")
	if err != nil {
		t.Fatal(err)
	}
	source.failure = errors.New("origem removida")
	second, err := workspace.Open(context.Background(), requestID, "digest-1", "base.jornada")
	if err != nil || second != first || source.reads != 1 {
		t.Fatalf("replay reabriu a origem: %+v (%v)", second, err)
	}
	if _, err := workspace.Open(context.Background(), requestID, "digest-2", "base.jornada"); err == nil || source.reads != 1 {
		t.Fatal("ID reutilizado com outros argumentos")
	}
	// snapshot é isolado do documento retornado pela fonte.
	source.document.Project.Name = "Alterada fora"
	snapshot, err := core.Inspect(first.ID)
	if err != nil || snapshot.Document.Project.Name != "Base" {
		t.Fatal("estado compartilhado com IO")
	}
}

// TestAuthoringWorkspaceNewDoesNotPersist cria rascunho usando o construtor injetado.
func TestAuthoringWorkspaceNewDoesNotPersist(t *testing.T) {
	// source não oferece método de gravação; criação só admite o documento na memória.
	source := &memoryDocumentSource{}
	core := NewAuthoring(DefaultAuthoringLimits(), nil)
	workspace := NewAuthoringWorkspace(core, source)
	info, err := workspace.Create(context.Background(), "00000000-0000-4000-8000-000000000002", "new-digest", "Nova", false)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := core.Inspect(info.ID)
	if err != nil || !snapshot.DocumentDirty || snapshot.Document.Project.Name != "Nova" || source.reads != 0 {
		t.Fatal("criação não permaneceu transitória")
	}
	if _, err := workspace.Open(context.Background(), "00000000-0000-4000-8000-000000000003", "bad-name", "../fora.jornada"); err == nil || source.reads != 0 {
		t.Fatal("nome interpretado como caminho")
	}
}
