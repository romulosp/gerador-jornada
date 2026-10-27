package desktop

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"path/filepath"
	"testing"
)

// gitDialogWidgets encontra controles reais do diálogo sem executar operações de rede.
// object é percorrido uma vez por visited; entries/lists recebem os widgets editáveis do fluxo.
func gitDialogWidgets(object fyne.CanvasObject, entries *[]*widget.Entry, lists *[]*widget.List, visited map[fyne.CanvasObject]bool) {
	if object == nil || visited[object] {
		return
	}
	visited[object] = true
	if control, ok := object.(*widget.Entry); ok {
		*entries = append(*entries, control)
		return
	}
	if control, ok := object.(*widget.List); ok {
		*lists = append(*lists, control)
		return
	}
	// children inclui o popup e seus renderers; entradas/listas não precisam ser expandidas.
	var children []fyne.CanvasObject
	if group, ok := object.(*fyne.Container); ok {
		children = group.Objects
	} else if control, ok := object.(fyne.Widget); ok {
		children = test.WidgetRenderer(control).Objects()
	}
	for _, child := range children {
		gitDialogWidgets(child, entries, lists, visited)
	}
}

// TestGitRepositoryFilterAndExactConfirmation usa o destino selecionado sem pedir seu nome novamente.
func TestGitRepositoryFilterAndExactConfirmation(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app/editor estão desconectados e o teste para antes de confirmar qualquer rede.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Destino"))
	defer editor.Window.Close()
	editor.SetProjectPath(filepath.Join(t.TempDir(), "A.jornada"))
	editor.saveWithCompletion(false, func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})
	editor.Window.Show()
	// selected captura a identidade que seria entregue à conexão, sem autenticação ou rede real.
	var selected application.GitHubRepository
	editor.chooseGitRepository([]application.GitHubRepository{
		{ID: 1, Owner: "fixture", Name: "gerador-jornada", Branch: "main", Writable: true},
		{ID: 2, Owner: "fixture", Name: "repositorio_jornadas", Branch: "main", Writable: true},
	}, func(repo application.GitHubRepository) { selected = repo })
	// entries/lists apontam para busca e lista virtualizada do seletor em uso.
	var entries []*widget.Entry
	var lists []*widget.List
	gitDialogWidgets(editor.Window.Canvas().Overlays().Top(), &entries, &lists, map[fyne.CanvasObject]bool{})
	if len(entries) != 1 || len(lists) != 1 {
		t.Fatal("seletor sem busca/lista")
	}
	lists[0].Select(0)
	confirm := findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), "Conectar repositório", map[fyne.CanvasObject]bool{})
	if confirm == nil || confirm.Disabled() {
		t.Fatal("seleção inicial não habilitou ação")
	}
	entries[0].SetText("repositorio_jornadas")
	if !confirm.Disabled() {
		t.Fatal("filtro conservou destino anterior")
	}
	lists[0].Select(0)
	test.Tap(confirm)
	if selected.ID != 2 || selected.FullName() != "fixture/repositorio_jornadas" {
		t.Fatal("conexão trocou o destino selecionado")
	}
	// O seletor deve terminar sem solicitar o nome novamente.
	entries = nil
	lists = nil
	gitDialogWidgets(editor.Window.Canvas().Overlays().Top(), &entries, &lists, map[fyne.CanvasObject]bool{})
	if len(entries) != 0 {
		t.Fatal("selecao abriu um segundo formulario para redigitar o repositorio")
	}

}

// TestGitRepositoryCancelDoesNotConnect mantém cancelamento sem efeito remoto ou nova confirmação.
func TestGitRepositoryCancelDoesNotConnect(t *testing.T) {
	// app/editor são isolados; o callback detecta qualquer conexão indevida.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Cancelar seleção"))
	defer editor.Window.Close()
	editor.Window.Show()
	called := false
	editor.chooseGitRepository([]application.GitHubRepository{{ID: 2, Owner: "fixture", Name: "destino", Writable: true}}, func(application.GitHubRepository) { called = true })
	cancel := findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), "Cancelar", map[fyne.CanvasObject]bool{})
	if cancel == nil {
		t.Fatal("seletor não oferece cancelar")
	}
	test.Tap(cancel)
	if called {
		t.Fatal("cancelamento iniciou conexão")
	}
}

// TestGitSaveCompletionAndCancel protege publicação contra salvamento incompleto e migração recusada.
func TestGitSaveCompletionAndCancel(t *testing.T) {
	// app/editor são isolados da instância nativa e não têm cliente GitHub configurado.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Teste de salvamento"))
	editor.SetProjectPath(filepath.Join(t.TempDir(), "A.jornada"))
	editor.Session.Document.Project.Revision++
	called := 0
	editor.saveWithCompletion(false, func(err error) {
		if err != nil {
			t.Error(err)
		}
		called++
	})
	if called != 1 || editor.Session.Dirty() {
		t.Fatal("publicação não recebeu conclusão íntegra")
	}
	// A versão legada deve pedir confirmação antes da próxima escrita.
	editor.Session.Document.Project.CatalogVersion = 3
	editor.Session.Document.Project.Revision++
	editor.Window.Show()
	var outcome error
	editor.saveWithCompletion(false, func(err error) { outcome = err; called++ })
	cancel := findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), "No", map[fyne.CanvasObject]bool{})
	if cancel == nil {
		cancel = findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), "Não", map[fyne.CanvasObject]bool{})
	}
	if cancel == nil {
		cancel = findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), "Cancel", map[fyne.CanvasObject]bool{})
	}
	if cancel == nil {
		t.Fatal("confirmação de migração não exibida")
	}
	test.Tap(cancel)
	if called != 2 || !errors.Is(outcome, errSaveCanceled) || !editor.Session.Dirty() {
		t.Fatal("cancelamento não preservou edição")
	}
}

// TestGitPanelDoesNotHideTree conserva as jornadas no viewport e bloqueia publicação sem conta/vínculo.
func TestGitPanelDoesNotHideTree(t *testing.T) {
	// app e configurações ficam separados dos dados reais do usuário.
	app := test.NewApp()
	defer app.Quit()
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	editor := New(app, application.NewDocument("Workspace local"))
	editor.EnableGitHub()
	editor.SetProjectPath(filepath.Join(t.TempDir(), "A.jornada"))
	editor.Window.Show()
	if !editor.GitPublishButton.Disabled() || editor.GitLabel.Text == "" {
		t.Fatal("estado inicial inválido")
	}
	if editor.GitLabel.Size().Height > 45 {
		t.Fatal("estado administrativo consumiu a árvore")
	}
}
