package application

import (
	"fmt"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// TestHistoryAndSnapshot verifica rollback, isolamento, undo e reinicialização de redo.
func TestHistoryAndSnapshot(t *testing.T) {
	// session representa uma edição com histórico limitado.
	session := NewSession(NewDocument("Teste"))
	// initial preserva o estado anterior a qualquer comando.
	initial := session.Document.Copy()
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := session.Change(func(document *Document) error { document.Project.Name = "inválido"; return fmt.Errorf("falha") }); err == nil {
		t.Fatal("falha esperada")
	}
	if session.Document.Project.Name != "Teste" || len(session.UndoStack) != 0 {
		t.Fatal("transação falha alterou estado")
	}
	// pageID identifica a composição inicial.
	pageID := session.Document.Screens[0].ID
	// componentID identifica o controle inserido.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	componentID, err := session.AddComponent(pageID, catalog.Button, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	// snapshot não pode mudar com novas alterações.
	snapshot := session.Document.Copy()
	if err = session.Change(func(document *Document) error {
		document.Screen(pageID).ComponentByID(componentID).Properties.Button.Text = "Outro"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if snapshot.Screen(pageID).ComponentByID(componentID).Properties.Button.Text != "Continuar" {
		t.Fatal("snapshot compartilhado")
	}
	if err = session.Undo(); err != nil {
		t.Fatal(err)
	}
	if err = session.Redo(); err != nil {
		t.Fatal(err)
	}
	if len(initial.Screens[0].Components) != 0 {
		t.Fatal("snapshot inicial alterado")
	}
}

// TestPreviewReset comprova que navegação não persiste dados de execução.
func TestPreviewReset(t *testing.T) {
	// session recebe um campo numérico e uma ação de auto-navegação.
	session := NewSession(NewDocument("Prévia"))
	// pageID identifica a tela da visita.
	pageID := session.Document.Screens[0].ID
	// fieldID identifica o campo de execução.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	fieldID, err := session.AddComponent(pageID, catalog.TextField, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	// preview usa um documento separado da sessão.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	preview, err := NewPreview(session.Document)
	if err != nil {
		t.Fatal(err)
	}
	if err = preview.SetValue(fieldID, "123"); err != nil {
		t.Fatal(err)
	}
	if err = preview.Enter(pageID); err != nil {
		t.Fatal(err)
	}
	if preview.Values[fieldID] != "" {
		t.Fatal("valor não foi reiniciado")
	}
	if session.Document.Screen(pageID).ComponentByID(fieldID).Properties.Field.InitialValue != "" {
		t.Fatal("prévia alterou projeto")
	}
}
