package desktop

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/application"
)

// TestCurrentScreenPreviewButton valida a janela da tela ativa sem substituir a prévia do projeto.
func TestCurrentScreenPreviewButton(t *testing.T) {
	// app fornece janelas controláveis sem executar o binário de distribuição.
	app := test.NewApp()
	defer app.Quit()
	// editor mantém a segunda tela selecionada, diferente da tela inicial do projeto.
	editor := New(app, application.Example())
	defer editor.close()
	editor.PageID = editor.Session.Document.Screens[1].ID
	editor.refresh()
	// before conserva o estado que deve permanecer inalterado pela visualização.
	before := editor.Session.Document.Copy()
	test.Tap(editor.CurrentScreenPreview)
	// windows permite inspecionar a janela produzida pelo botão real.
	windows := app.Driver().AllWindows()
	// current é a janela exclusiva da tela selecionada.
	current := windows[len(windows)-1]
	if current == editor.Window || current.Title() != "Prévia Tela Atual · "+editor.page().Name {
		t.Fatal("botão não abriu a prévia da tela selecionada")
	}
	current.Close()
	editor.preview()
	windows = app.Driver().AllWindows()
	// project mantém o fluxo anterior de execução do projeto completo.
	project := windows[len(windows)-1]
	if project.Title() != "Prévia · "+before.Project.Name {
		t.Fatal("prévia do projeto foi substituída pelo modo de tela")
	}
	project.Close()
	if !reflect.DeepEqual(before, editor.Session.Document) || len(editor.Session.UndoStack) != 0 || editor.PageID != before.Screens[1].ID {
		t.Fatal("abrir prévias não pode alterar documento, histórico ou tela ativa")
	}
	editor.PageID = ""
	editor.refresh()
	if !editor.CurrentScreenPreview.Disabled() {
		t.Fatal("sem tela atual o botão precisa ficar desabilitado")
	}
}
