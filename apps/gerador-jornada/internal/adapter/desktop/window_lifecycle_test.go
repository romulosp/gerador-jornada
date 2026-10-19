package desktop

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// TestFullScreenExitRequiresConfirmation cobre visibilidade, cancelamento e confirmação sem restaurar.
func TestFullScreenExitRequiresConfirmation(t *testing.T) {
	// app evita afirmar que os testes headless exercitam os estilos nativos do Windows.
	app := test.NewApp()
	defer app.Quit()
	// editor fornece um documento alterado para proteger o encerramento.
	editor := New(app, application.Example())
	defer func() {
		if !editor.Mode.closed.Load() {
			editor.close()
		}
	}()
	if editor.Mode.exitButton.Visible() {
		t.Fatal("Sair não deve aparecer em janela normal")
	}
	editor.ShowWindow()
	if !editor.Mode.exitButton.Visible() {
		t.Fatal("Sair deve aparecer em tela cheia")
	}
	test.Tap(editor.Mode.button)
	if editor.Mode.exitButton.Visible() {
		t.Fatal("Restaurar deve ocultar Sair")
	}
	test.Tap(editor.Mode.button)
	if err := editor.Session.Change(func(document *application.Document) error {
		document.Project.Name = "Projeto alterado"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	// before registra dados que precisam sobreviver à resposta Não.
	before := editor.Session.Document.Copy()
	test.Tap(editor.Mode.exitButton)
	if editor.Mode.exitDialog == nil || editor.Mode.closed.Load() {
		t.Fatal("Sair precisa pedir confirmação antes de descartar alterações")
	}
	// confirmation conserva a instância para verificar que pedidos repetidos não empilham diálogos.
	confirmation := editor.Mode.exitDialog
	editor.requestClose()
	if editor.Mode.exitDialog != confirmation {
		t.Fatal("encerramento abriu confirmações duplicadas")
	}
	confirmation.Dismiss()
	if editor.Mode.exitDialog != nil || editor.Mode.closed.Load() || !editor.Window.FullScreen() || !reflect.DeepEqual(before, editor.Session.Document) {
		t.Fatal("Não deve preservar projeto e tela cheia sem encerrar")
	}
	test.Tap(editor.Mode.exitButton)
	editor.Mode.exitDialog.Confirm()
	if !editor.Mode.closed.Load() {
		t.Fatal("Sim deve encerrar sem exigir restaurar a janela")
	}
}

// TestFullScreenExitCleanDocument fecha diretamente quando não existe alteração pendente no documento.
func TestFullScreenExitCleanDocument(t *testing.T) {
	// app fornece uma janela de teste sem dependência da sessão gráfica.
	app := test.NewApp()
	defer app.Quit()
	// editor inicia com a mesma revisão que consta como salva na sessão.
	editor := New(app, application.Example())
	defer func() {
		if !editor.Mode.closed.Load() {
			editor.close()
		}
	}()
	editor.ShowWindow()
	// field altera apenas o rascunho, sem gerar revisão salva no domínio.
	field := editor.Inspector.Objects[1].(*widget.Form).Items[0].Widget.(*widget.Entry)
	field.SetText("Rascunho sem aplicar")
	test.Tap(editor.Mode.exitButton)
	if editor.Mode.exitDialog == nil || editor.Mode.closed.Load() {
		t.Fatal("rascunho também deve exigir confirmação de saída")
	}
	editor.Mode.exitDialog.Dismiss()
	if field.Text != "Rascunho sem aplicar" || editor.Mode.closed.Load() {
		t.Fatal("Não deve conservar o rascunho sem encerrar")
	}
	editor.refreshInspector()
	test.Tap(editor.Mode.exitButton)
	if !editor.Mode.closed.Load() || editor.Mode.exitDialog != nil {
		t.Fatal("projeto sem alterações deve encerrar pelo botão Sair")
	}
}
