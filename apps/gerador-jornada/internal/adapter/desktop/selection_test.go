package desktop

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// clickBoard exercita seleção no mesmo ponto lógico, independentemente do zoom da interface.
func clickBoard(editor *Editor, x, y int, modifier fyne.KeyModifier) {
	editor.Board.MouseDown(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary, Modifier: modifier,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(24+float32(x)*editor.Zoom, 24+float32(y)*editor.Zoom)}})
}

// TestMultipleSelectionAndHistory cobre alternância, arraste coletivo, caixa, cursor e foco.
func TestMultipleSelectionAndHistory(t *testing.T) {
	// app usa um canvas independente, sem afirmar homologação do driver Windows.
	app := test.NewApp()
	defer app.Quit()
	// editor fornece título e rótulo separados verticalmente na composição de referência.
	editor := New(app, application.Example())
	defer editor.close()
	clickBoard(editor, 50, 50, 0)
	clickBoard(editor, 50, 110, fyne.KeyModifierControl)
	if len(editor.selectionIDs()) != 2 {
		t.Fatal("Ctrl+clique não adicionou o segundo componente")
	}
	clickBoard(editor, 50, 50, 0)
	if len(editor.selectionIDs()) != 2 {
		t.Fatal("clicar em um membro não pode desmontar o conjunto antes do arraste")
	}
	editor.Board.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 10, DY: 10}})
	editor.Board.DragEnd()
	if editor.page().Components[0].Bounds.X != 58 || editor.page().Components[1].Bounds.X != 58 {
		t.Fatal("o arraste precisa mover os dois controles")
	}
	if len(editor.Session.UndoStack) != 1 {
		t.Fatal("movimento coletivo gerou múltiplos comandos")
	}
	if err := editor.Session.Undo(); err != nil || editor.page().Components[0].Bounds.X != 48 {
		t.Fatalf("movimento não reversível: %v", err)
	}
	// field representa um campo com foco: setas de desenho devem ser ignoradas.
	field := widget.NewEntry()
	editor.Inspector.Add(field)
	editor.Window.Canvas().Focus(field)
	editor.moveSelection(10, 0)
	if editor.page().Components[0].Bounds.X != 48 {
		t.Fatal("atalho de desenho capturou foco textual")
	}
	editor.Window.Canvas().Unfocus()
	editor.handleCanvasKey(&fyne.KeyEvent{Name: fyne.KeyRight}, fyne.KeyModifierShift)
	if editor.page().Components[0].Bounds.X != 58 || editor.page().Components[1].Bounds.X != 58 {
		t.Fatal("Shift+seta precisa mover o conjunto em dez unidades")
	}
	if err := editor.Session.Undo(); err != nil {
		t.Fatal(err)
	}
	clickBoard(editor, 25, 25, 0)
	editor.Board.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 730, DY: 115}})
	editor.Board.DragEnd()
	if len(editor.selectionIDs()) != 2 {
		t.Fatal("caixa de seleção não encontrou os dois textos")
	}
	editor.selectAll()
	if len(editor.selectionIDs()) != len(editor.page().Components) {
		t.Fatal("selecionar todos incompleto")
	}
	editor.setSelection([]string{editor.page().Components[0].ID})
	editor.Board.MouseMoved(&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(74, 74)}})
	if editor.Board.Cursor() != desktop.CrosshairCursor {
		t.Fatal("componente selecionado precisa indicar movimentação")
	}
	// bounds fornece a posição da alça inferior direita da seleção única.
	bounds := editor.page().Components[0].Bounds
	editor.Board.MouseMoved(&desktop.MouseEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(24+float32(bounds.X+bounds.Width-2), 24+float32(bounds.Y+bounds.Height-2))}})
	if editor.Board.Cursor() != desktop.NWSEResizeCursor {
		t.Fatal("alça precisa indicar redimensionamento diagonal")
	}
}

// TestStickyPropertiesAndTextAlignment confirma o rascunho sem mover a faixa ou perder o texto.
func TestStickyPropertiesAndTextAlignment(t *testing.T) {
	// app mantém widgets e callbacks no mesmo ambiente de teste.
	app := test.NewApp()
	defer app.Quit()
	// editor recebe uma tela vazia para inserir uma faixa explícita.
	editor := New(app, application.NewDocument("Faixa"))
	defer editor.close()
	editor.insertPreset("pageTitle", 0, 100)
	// before permite verificar que alinhamento do texto não afeta o elemento.
	before := editor.page().ComponentByID(editor.Selected).Bounds
	// form contém o rascunho; a ação de salvar fica no rodapé separado.
	form := editor.Inspector.Objects[1].(*widget.Form)
	for _, item := range form.Items {
		if item.Text == "Alinhamento do texto" {
			item.Widget.(*widget.Select).SetSelected("Centro")
		}
		if item.Text == "Texto" {
			item.Widget.(*widget.Entry).SetText("Novo título com acentuação")
		}
	}
	if editor.page().ComponentByID(editor.Selected).Properties.Title.Text == "Novo título com acentuação" {
		t.Fatal("rascunho foi publicado antes da confirmação")
	}
	if len(editor.InspectorActions.Objects) != 2 {
		t.Fatal("ações fixas não disponíveis fora da rolagem")
	}
	test.Tap(editor.InspectorActions.Objects[0].(*widget.Button))
	// component contém somente os dados confirmados pelo botão explícito.
	component := editor.page().ComponentByID(editor.Selected)
	if component.Style.TextAlign != "center" || component.Bounds != before || component.Properties.Title.Text != "Novo título com acentuação" {
		t.Fatal("salvar propriedades não preservou conteúdo, alinhamento e geometria")
	}
	if err := editor.Session.Undo(); err != nil || editor.page().ComponentByID(editor.Selected).Style.TextAlign != "right" {
		t.Fatalf("alinhamento não reversível: %v", err)
	}
	editor.refreshInspector()
	// revision comprova que descartar o formulário não é uma edição de documento.
	revision := editor.Session.Document.Project.Revision
	form = editor.Inspector.Objects[1].(*widget.Form)
	for _, item := range form.Items {
		if item.Text == "Texto" {
			item.Widget.(*widget.Entry).SetText("Rascunho descartado")
		}
	}
	test.Tap(editor.InspectorActions.Objects[1].(*widget.Button))
	if editor.Session.Document.Project.Revision != revision || editor.page().ComponentByID(editor.Selected).Properties.Title.Text == "Rascunho descartado" {
		t.Fatal("descartar publicou conteúdo ou revisão")
	}
}
