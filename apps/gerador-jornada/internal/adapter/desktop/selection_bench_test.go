package desktop

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// BenchmarkSelect500 mede seleção alternada em uma tela com quinhentos controles já renderizados.
func BenchmarkSelect500(b *testing.B) {
	// app fornece o mesmo driver de teste para comparação antes/depois.
	app := test.NewApp()
	defer app.Quit()
	// document contém controles pequenos, com posições conhecidas e sem dependências externas.
	document := application.NewDocument("Seleção de 500 controles")
	for index := 0; index < 500; index++ {
		// component distribui rótulos em vinte colunas para alternar seleção sem sobreposição.
		component, _ := screen.NewComponent(catalog.Label, fmt.Sprintf("rotulo%d", index), index%20*40, index/20*22)
		component.Bounds.Width, component.Bounds.Height = 38, 20
		document.Screens[0].Components = append(document.Screens[0].Components, component)
	}
	// editor termina a montagem antes da medição.
	editor := New(app, document)
	defer editor.close()
	editor.Board.Refresh()
	b.ReportAllocs()
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		// event alterna entre dois componentes para incluir uma mudança real de seleção.
		event := &desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(26+float32(index%2*40), 26)}}
		editor.Board.MouseDown(event)
		editor.Board.MouseUp(event)
	}
}
