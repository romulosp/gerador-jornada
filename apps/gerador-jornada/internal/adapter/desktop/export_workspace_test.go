package desktop

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/adapter/pptx"
	"gerador-jornada.local/editor/internal/application"
)

// TestWorkspacePowerPointOutput gera o pacote real junto à jornada, sem solicitar pasta ou arquivo.
func TestWorkspacePowerPointOutput(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app/editor usam um nome exibido diferente do nome físico para detectar destinos incorretos.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Nome exibido"))
	defer editor.Window.Close()
	root := t.TempDir()
	// canonical compara a mesma pasta física, inclusive quando TEMP usa um alias 8.3 do Windows.
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	editor.SetProjectPath(filepath.Join(root, "Meu projeto (2).jornada"))
	// received só será preenchido se o destino automático dispensar o diálogo antigo.
	received := ""
	editor.workspacePPTOutput(func(path string) {
		received = path
		if err := pptx.Export(context.Background(), path, editor.Session.Document, pptx.Options{Exploration: true}, nil); err != nil {
			t.Fatal(err)
		}
	})
	if received != filepath.Join(canonical, "Meu projeto (2).pptx") {
		t.Fatal("destino não usa a pasta e o nome da jornada atual")
	}
	// archive comprova que o arquivo final é um PPTX real, não apenas um caminho calculado.
	archive, err := zip.OpenReader(received)
	if err != nil {
		t.Fatal(err)
	}
	archive.Close()
	// A apresentação anterior deve permanecer intacta quando a substituição for recusada.
	before, err := os.ReadFile(received)
	if err != nil {
		t.Fatal(err)
	}
	editor.Window.Show()
	called := false
	editor.workspacePPTOutput(func(string) { called = true })
	// cancel localiza a resposta negativa na tradução disponível do Fyne.
	var cancel fyne.CanvasObject
	for _, label := range []string{"No", "Não", "Cancel", "Cancelar"} {
		if button := findWorkspaceButton(editor.Window.Canvas().Overlays().Top(), label, map[fyne.CanvasObject]bool{}); button != nil {
			cancel = button
			break
		}
	}
	if cancel == nil {
		t.Fatal("substituição não oferece cancelamento")
	}
	test.Tap(cancel.(fyne.Tappable))
	after, err := os.ReadFile(received)
	if err != nil || called || string(after) != string(before) {
		t.Fatal("cancelamento alterou a apresentação anterior")
	}
}

// TestFirstSaveUsesConfiguredWorkspace salva o primeiro pacote diretamente na pasta já escolhida.
func TestFirstSaveUsesConfiguredWorkspace(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app/editor simulam documento novo associado a um workspace local existente.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Meu projeto"))
	defer editor.Window.Close()
	editor.WorkspacePath = t.TempDir()
	called := false
	editor.saveWithCompletion(false, func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		called = true
	})
	if !called || editor.FilePath != filepath.Join(editor.WorkspacePath, "Meu projeto.jornada") {
		t.Fatal("salvamento pediu outro destino apesar do workspace configurado")
	}
	if _, err := os.Stat(editor.FilePath); err != nil {
		t.Fatal(err)
	}
}
