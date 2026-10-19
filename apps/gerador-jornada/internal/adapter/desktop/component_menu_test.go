package desktop

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// TestComponentContextActions cobre alvo, propriedades, foco e exclusão cancelável pelo menu real.
func TestComponentContextActions(t *testing.T) {
	// app usa o driver de testes sem executar o EXE distribuído.
	app := test.NewApp()
	defer app.Quit()
	// editor fornece componentes com posições conhecidas na primeira tela.
	editor := New(app, application.Example())
	defer editor.close()
	if !editor.DeleteButton.Disabled() {
		t.Fatal("Excluir precisa iniciar desabilitado")
	}
	// before conserva a quantidade e a revisão anteriores a qualquer ação confirmada.
	before := len(editor.page().Components)
	// revision permite detectar alterações indevidas ao abrir/cancelar ações.
	revision := editor.Session.Document.Project.Revision
	editor.setZoom(.5)
	test.TapSecondaryAt(editor.Board, fyne.NewPos(49, 49))
	editor.setZoom(1)
	if editor.Board.ContextMenu == nil || len(editor.Board.ContextMenu.Items) != 2 || editor.DeleteButton.Disabled() {
		t.Fatal("clique direito precisa selecionar o alvo e oferecer duas ações")
	}
	// target permite verificar que a edição e a exclusão se referem ao componente clicado.
	target := editor.page().Components[0].ID
	if editor.Selected != target {
		t.Fatal("menu selecionou um componente diferente do alvo")
	}
	test.Tap(editor.Board.ContextMenu.Items[0].(fyne.Tappable))
	if _, ok := editor.Window.Canvas().Focused().(*widget.Entry); !ok {
		t.Fatal("Alterar propriedades precisa direcionar o foco ao formulário")
	}
	// Delete dentro de um campo deve continuar sendo edição textual.
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyDelete}, 0)
	if editor.DeleteDialog != nil {
		t.Fatal("Delete com foco textual não pode excluir componente")
	}
	test.TapSecondaryAt(editor.Board, fyne.NewPos(74, 74))
	test.Tap(editor.Board.ContextMenu.Items[1].(fyne.Tappable))
	if editor.DeleteDialog == nil || len(editor.page().Components) != before {
		t.Fatal("menu deve pedir confirmação antes de excluir")
	}
	editor.DeleteDialog.Dismiss()
	if editor.DeleteDialog != nil || len(editor.page().Components) != before || editor.Session.Document.Project.Revision != revision {
		t.Fatal("a resposta Não deve preservar componentes, revisão e seleção")
	}
	test.Tap(editor.DeleteButton)
	if editor.DeleteDialog == nil {
		t.Fatal("botão Excluir não abriu confirmação")
	}
	editor.DeleteDialog.Confirm()
	if editor.page().ComponentByID(target) != nil || !editor.DeleteButton.Disabled() || len(editor.Session.UndoStack) != 1 {
		t.Fatal("Sim precisa excluir uma vez, limpar seleção e desabilitar o botão")
	}
	if err := editor.Session.Undo(); err != nil {
		t.Fatal(err)
	}
	editor.changed()
	if editor.page().ComponentByID(target) == nil || !editor.DeleteButton.Disabled() {
		t.Fatal("desfazer deve restaurar o componente sem inventar seleção")
	}
	test.TapSecondaryAt(editor.Board, fyne.NewPos(25, 25))
	if editor.Board.ContextMenu != nil {
		t.Fatal("área vazia não pode abrir ações de componente")
	}
}

// TestDeleteSelectionKeyboardAndScreen cobre confirmação única, grupos e troca de tela.
func TestDeleteSelectionKeyboardAndScreen(t *testing.T) {
	// app isola a árvore de widgets e as confirmações deste cenário.
	app := test.NewApp()
	defer app.Quit()
	// editor começa sem seleção e sem histórico de edição.
	editor := New(app, application.Example())
	defer editor.close()
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyDelete}, 0)
	if editor.DeleteDialog != nil {
		t.Fatal("Delete sem seleção não pode abrir confirmação")
	}
	clickBoard(editor, 50, 50, 0)
	clickBoard(editor, 50, 110, fyne.KeyModifierControl)
	// before distingue a remoção do grupo de uma exclusão individual indevida.
	before := len(editor.page().Components)
	test.TapSecondaryAt(editor.Board, fyne.NewPos(74, 74))
	if len(editor.selectionIDs()) != 2 {
		t.Fatal("menu de um membro precisa preservar a seleção coletiva")
	}
	editor.Board.ContextMenu.Hide()
	editor.Window.Canvas().Unfocus()
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyDelete}, 0)
	// confirmation precisa continuar sendo a mesma caso Delete seja pressionado novamente.
	confirmation := editor.DeleteDialog
	editor.removeComponent()
	if confirmation == nil || editor.DeleteDialog != confirmation {
		t.Fatal("a mesma operação não pode acumular confirmações")
	}
	confirmation.Confirm()
	if len(editor.page().Components) != before-2 || !editor.DeleteButton.Disabled() || len(editor.Session.UndoStack) != 1 {
		t.Fatal("Delete precisa excluir o grupo em uma ação reversível")
	}
	if err := editor.Session.Undo(); err != nil {
		t.Fatal(err)
	}
	editor.changed()
	clickBoard(editor, 50, 50, 0)
	if editor.DeleteButton.Disabled() {
		t.Fatal("seleção restaurada precisa habilitar Excluir")
	}
	clickBoard(editor, 1, 1, 0)
	if !editor.DeleteButton.Disabled() {
		t.Fatal("clique vazio precisa limpar a seleção e desabilitar Excluir")
	}
	clickBoard(editor, 50, 50, 0)
	test.Tap(editor.Screens.Objects[1].(*widget.Button))
	if !editor.DeleteButton.Disabled() || len(editor.selectionIDs()) != 0 {
		t.Fatal("trocar de tela precisa desabilitar a exclusão anterior")
	}
}
