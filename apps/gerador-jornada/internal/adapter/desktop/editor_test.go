package desktop

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// TestEditorComposition verifica montagem dos painéis com o driver de teste Fyne.
func TestEditorComposition(t *testing.T) {
	// app não cria uma janela nativa nem simula homologação visual.
	app := test.NewApp()
	defer app.Quit()
	// editor recebe o conjunto de referência validado.
	editor := New(app, application.Example())
	if editor.Board.MinSize().Width < 800 || len(editor.Tabs.Items) != 3 {
		t.Fatal("layout incompleto")
	}
	// initial conserva o número de componentes antes do comando.
	initial := len(editor.page().Components)
	editor.insert(catalog.Label, 40, 40)
	if len(editor.page().Components) != initial+1 || editor.Selected == "" {
		t.Fatal("inclusão não conectada ao domínio")
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := editor.Session.Undo(); err != nil {
		t.Fatal(err)
	}
	editor.close()
}
