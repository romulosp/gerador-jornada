package desktop

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
)

// TestDragOverBackground verifica que o fundo não interfere com seleção, arraste, zoom ou histórico.
func TestDragOverBackground(t *testing.T) {
	// app permite executar o mesmo renderer e comandos com um canvas de teste.
	app := test.NewApp()
	defer app.Quit()
	// document usa o recurso real incorporado e as sete famílias de referência.
	document := application.Example()
	builtinassets.DefaultBackgrounds(&document)
	// editor mantém o fundo independente da geometria do componente escolhido.
	editor := New(app, document)
	defer editor.close()
	editor.setZoom(.5)
	// original fornece um título que será selecionado por coordenadas sobre a imagem.
	original := editor.page().Components[0]
	editor.Board.MouseDown(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(24+float32(original.Bounds.X+4)*.5, 24+float32(original.Bounds.Y+4)*.5)}})
	editor.Board.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 10, DY: 15}})
	editor.Board.DragEnd()
	if editor.page().Components[0].Bounds.X != original.Bounds.X+20 || editor.page().Components[0].Bounds.Y != original.Bounds.Y+30 {
		t.Fatal("o arraste deve respeitar coordenadas lógicas sobre o fundo")
	}
	if err := editor.Session.Undo(); err != nil || editor.page().Components[0].Bounds != original.Bounds {
		t.Fatalf("arraste não reversível: %v", err)
	}
	// state reutiliza o projeto para verificar o fundo da execução interativa.
	state, err := application.NewPreview(editor.Session.Document)
	if err != nil {
		t.Fatal(err)
	}
	// preview mantém a imagem antes dos controles interativos na pilha visual.
	preview := previewPage(state, func() {}).(*fyne.Container)
	if len(preview.Objects) < 3 {
		t.Fatal("prévia incompleta")
	}
	if _, ok := preview.Objects[1].(*canvas.Image); !ok {
		t.Fatal("a imagem precisa ficar atrás dos componentes na prévia")
	}
	// previous conserva o tamanho antes de inserir o conjunto auxiliar.
	previous := len(editor.page().Components)
	editor.insertPreset("header", 374, 6)
	if len(editor.page().Components) != previous+8 {
		t.Fatal("cabeçalho incompleto")
	}
	if err := editor.Session.Undo(); err != nil || len(editor.page().Components) != previous {
		t.Fatalf("o conjunto deve ser desfeito em uma única operação: %v", err)
	}
}
