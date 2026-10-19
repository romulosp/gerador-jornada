package desktop

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// TestFullScreenPreservesEditing comprova alternância reversível sem publicar alterações no projeto.
func TestFullScreenPreservesEditing(t *testing.T) {
	// app usa o driver de teste, sem afirmar que uma janela Win32 foi maximizada de fato.
	app := test.NewApp()
	defer app.Quit()
	// editor conserva seleção, zoom e formulário em andamento enquanto muda apenas a janela.
	editor := New(app, application.Example())
	defer editor.close()
	if editor.Window.FullScreen() {
		t.Fatal("tela cheia foi aplicada antes da exibição nativa")
	}
	editor.ShowWindow()
	if !editor.Window.FullScreen() || !editor.Mode.fullScreen.Load() || editor.Mode.button.Text != "Restaurar" {
		t.Fatal("editor não iniciou em tela cheia com ação de retorno")
	}
	// Restaurar permite validar também a reentrada sem reiniciar o editor.
	test.Tap(editor.Mode.button)
	if editor.Window.FullScreen() {
		t.Fatal("botão não restaurou a janela após a abertura em tela cheia")
	}
	clickBoard(editor, 50, 50, 0)
	editor.setZoom(.5)
	// before preserva o documento para comparar após entrada e saída do modo.
	before := editor.Session.Document.Copy()
	// selected identifica o controle que não deve ser desmarcado ao sair por Escape.
	selected := editor.Selected
	// field fornece um rascunho que não pode ser salvo ou descartado pela mudança de janela.
	field := editor.Inspector.Objects[1].(*widget.Form).Items[0].Widget.(*widget.Entry)
	field.SetText("Rascunho ainda não salvo")
	editor.Window.Canvas().Focus(field)
	test.Tap(editor.Mode.button)
	if !editor.Window.FullScreen() || editor.Mode.button.Text != "Restaurar" {
		t.Fatal("botão não entrou em tela cheia com ação de retorno")
	}
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyEscape}, 0)
	if editor.Window.FullScreen() || editor.Mode.button.Text != "Tela cheia" {
		t.Fatal("Escape não restaurou a janela")
	}
	if field.Text != "Rascunho ainda não salvo" || editor.Selected != selected || editor.Zoom != .5 || !reflect.DeepEqual(before, editor.Session.Document) || len(editor.Session.UndoStack) != 0 {
		t.Fatal("tela cheia alterou formulário, seleção, zoom, documento ou histórico")
	}
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyF12}, 0)
	if !editor.Window.FullScreen() {
		t.Fatal("F12 encaminhado ao canvas não entrou em tela cheia")
	}
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyF12}, 0)
	if editor.Window.FullScreen() {
		t.Fatal("segundo F12 não restaurou a janela")
	}
	editor.Mode.close()
	editor.Mode.applyWindowAction(windowEnter)
	if editor.Window.FullScreen() {
		t.Fatal("ação enfileirada após fechar o controle alterou a janela")
	}
}

// TestWindowsWindowMessages protege maximização, restauração e teclas sem consumir fechar/minimizar.
func TestWindowsWindowMessages(t *testing.T) {
	// cases representa mensagens Win32 documentadas e eventos que devem continuar no Fyne.
	cases := []struct {
		// message e command representam o tipo de mensagem e seu argumento principal.
		message, command uintptr
		// flags preserva o bit de repetição de tecla do Windows.
		flags uintptr
		// full representa o estado transitório de tela cheia antes do evento.
		full bool
		// want é a intenção que pode ser enfileirada no adaptador UI.
		want windowAction
	}{
		{winSystemCommand, winMaximize | 2, 0, false, windowEnter},
		{winSystemCommand, winRestore, 0, true, windowExit},
		{winSystemCommand, winRestore, 0, false, windowPass},
		{winSystemCommand, 0xF060, 0, true, windowPass},
		{winSystemCommand, 0xF020, 0, true, windowPass},
		{winKeyDown, winF12, 0, false, windowToggle},
		{winKeyDown, winF12, 1 << 30, true, windowConsume},
		{winKeyDown, winEscape, 0, true, windowExit},
		{winKeyDown, winEscape, 0, false, windowPass},
		{winKeyUp, winF12, 0, false, windowPass},
	}
	for _, item := range cases {
		if got := windowMessageAction(item.message, item.command, item.flags, item.full); got != item.want {
			t.Fatalf("mensagem %x/%x: ação %d, esperava %d", item.message, item.command, got, item.want)
		}
	}
}
