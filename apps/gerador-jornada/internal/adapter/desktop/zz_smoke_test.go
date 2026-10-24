package desktop

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// smokeForm devolve o formulário de propriedades montado no inspector.
func smokeForm(editor *Editor) *widget.Form {
	for _, object := range editor.Inspector.Objects {
		if form, ok := object.(*widget.Form); ok {
			return form
		}
	}
	return nil
}

// smokeSave aciona o botão fixo de salvar alterações.
func smokeSave(editor *Editor) {
	editor.InspectorActions.Objects[0].(*widget.Button).OnTapped()
}

// TestSmokeCompoundInspectors confere que os três compostos inserem e publicam uma transação.
func TestSmokeCompoundInspectors(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.Example())
	defer editor.close()
	// base preserva a composição de exemplo antes das três inserções.
	base := len(editor.page().Components)
	for _, kind := range []catalog.Kind{catalog.BarcodeInput, catalog.BarcodeProcessed, catalog.ApplicationMenu} {
		editor.insert(kind, 40, 40)
		if editor.Selected == "" {
			t.Fatalf("inserção de %s não selecionou", kind)
		}
		editor.refresh()
		if smokeForm(editor) == nil {
			t.Fatalf("formulário ausente para %s", kind)
		}
		before := editor.Session.Document.Project.Revision
		smokeSave(editor)
		if editor.Session.Document.Project.Revision != before+1 {
			t.Fatalf("salvar não publicou %s", kind)
		}
	}
	// A lista de componentes permanece com uma instância por composto, não oito campos soltos.
	if count := len(editor.page().Components); count != base+3 {
		t.Fatalf("componentes = %d, esperado %d", count, base+3)
	}
}
