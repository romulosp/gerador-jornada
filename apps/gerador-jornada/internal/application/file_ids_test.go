package application

import (
	"reflect"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/screen"
)

// TestFileIDEditingAndCollisions garante edição reversível sem alterar UUIDs ou vínculos.
func TestFileIDEditingAndCollisions(t *testing.T) {
	// session fornece uma jornada real com referências entre três telas.
	session := NewSession(Example())
	// before preserva todas as referências para detectar renomeação indevida de identidades.
	before := session.Document.Copy()
	if err := session.Change(func(document *Document) error {
		document.Screens[0].FileID = "PnlAtendimento"
		document.Screens[0].Components[0].FileID = "TituloPrincipal"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if session.Document.Screens[0].ID != before.Screens[0].ID || !reflect.DeepEqual(session.Document.Journey, before.Journey) {
		t.Fatal("editar o ID de arquivo alterou UUID ou navegação")
	}
	if err := session.Undo(); err != nil || session.Document.Screens[0].FileID != before.Screens[0].FileID {
		t.Fatalf("renomeação não reversível: %v", err)
	}
	// invalid enumera colisões sem distinção de caixa e excesso de tamanho.
	for _, invalid := range []func(*Document){
		func(document *Document) { document.Screens[1].FileID = document.Screens[0].FileID },
		func(document *Document) {
			document.Screens[0].Components[1].FileID = strings.ToLower(document.Screens[0].Components[0].FileID)
		},
		func(document *Document) { document.Screens[0].FileID = "Pnl" + strings.Repeat("A", 38) },
		func(document *Document) { document.Screens[0].Components[0].FileID = "inválido" },
	} {
		// current permite comprovar atomicidade mesmo após sucessivas recusas.
		current := session.Document.Copy()
		if err := session.Change(func(document *Document) error { invalid(document); return nil }); err == nil {
			t.Fatal("ID inválido foi publicado")
		}
		if !reflect.DeepEqual(current, session.Document) {
			t.Fatal("edição recusada alterou o documento")
		}
	}
}

// TestFileIDMigrationAndDuplication verifica sugestões determinísticas, limite e cópias independentes.
func TestFileIDMigrationAndDuplication(t *testing.T) {
	// document simula telas antigas com nomes equivalentes depois de remover acentos.
	document := Example()
	document.Screens[1].Name = "01 · Identificação"
	for index := range document.Screens {
		document.Screens[index].FileID = ""
		for componentIndex := range document.Screens[index].Components {
			document.Screens[index].Components[componentIndex].FileID = ""
		}
	}
	EnsureFileIDs(&document)
	if document.Screens[0].FileID != "PnlIdentificacao" || document.Screens[1].FileID != "PnlIdentificacao2" {
		t.Fatal("migração não resolveu sugestões conflitantes deterministicamente")
	}
	// session permite duplicar sem colisão mesmo no limite de 40 caracteres.
	session := NewSession(document)
	session.Document.Screens[0].FileID = "Pnl" + strings.Repeat("A", 37)
	// original identifica o componente a duplicar no mesmo contexto.
	original := session.Document.Screens[0].Components[0].ID
	session.Document.Screens[0].Components[0].FileID = strings.Repeat("B", 40)
	if _, err := session.DuplicateComponents(session.Document.Screens[0].ID, []string{original}); err != nil {
		t.Fatal(err)
	}
	if _, err := session.DuplicateScreen(session.Document.Screens[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := session.Document.ValidateFileIDs(); err != nil {
		t.Fatal(err)
	}
	for _, page := range session.Document.Screens {
		if len(page.FileID) > screen.FileIDLimit {
			t.Fatal("duplicação ultrapassou o limite")
		}
	}
}
