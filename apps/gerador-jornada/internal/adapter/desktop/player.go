package desktop

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/scenario"
)

// playScenario usa os mesmos eventos interpretados que alimentam o exportador PPT.
func (e *Editor) playScenario(script scenario.Scenario) {
	// document é estável durante toda a reprodução.
	document := e.Session.Document.Copy()
	// visits validam o roteiro antes de abrir uma janela.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	visits, err := script.Interpret(document.Screens, document.Journey)
	if !e.check(err) {
		return
	}
	// window mantém sua reprodução independente do editor.
	window := e.App.NewWindow("Demonstração · " + script.Name)
	// content recebe cada quadro renderizado.
	content := container.NewStack()
	// elapsed é alterado somente na thread de UI.
	elapsed := 0
	// playing e closed controlam comandos sem acesso concorrente fora da UI.
	playing := true
	// closed impede atualizações da demonstração após o fechamento da janela.
	closed := false
	// status exibe o instante e a etapa de execução.
	status := widget.NewLabel("Reproduzindo")
	// update reconstrói os objetos visuais do quadro corrente.
	update := func() {
		// frame é uma derivação pura, sem mutação do documento.
		frame := application.FrameAt(document, visits, elapsed)
		// objects renderiza a visita sem adicionar controles interativos.
		objects := pageBackground(document, frame.Page)
		// component contém seu valor runtime.
		for _, component := range frame.Page.Components {
			objects = append(objects, previewPrimitives(document, component)...)
		}
		// cursor permanece vetorial na prévia.
		cursor := canvas.NewText("➤", rgb("#2563EB"))
		cursor.TextSize = 26
		cursor.Move(fyne.NewPos(float32(frame.Cursor.X), float32(frame.Cursor.Y)))
		objects = append(objects, cursor)
		if frame.Click {
			// circle desenha o realce visual do clique sobre a tela apresentada.
			circle := canvas.NewCircle(transparent())
			circle.StrokeColor = rgb("#F59E0B")
			circle.StrokeWidth = 3
			circle.Move(fyne.NewPos(float32(frame.Cursor.X)-10, float32(frame.Cursor.Y)-10))
			circle.Resize(fyne.NewSize(28, 28))
			objects = append(objects, circle)
		}
		content.Objects = []fyne.CanvasObject{container.NewWithoutLayout(objects...)}
		content.Refresh()
		status.SetText(fmt.Sprintf("%0.1f s · %s", float64(elapsed)/1000, frame.Page.Name))
		if frame.Finished {
			playing = false
			status.SetText("Demonstração concluída")
		}
	}
	// pause evita consumir tempo enquanto a apresentação está pausada.
	pause := widget.NewButton("Pausar / continuar", func() { playing = !playing })
	// restart reinicia estado e tempo desde o começo.
	restart := widget.NewButton("Reiniciar", func() { elapsed = 0; playing = true; update() })
	window.SetContent(container.NewBorder(container.NewHBox(pause, restart, status), nil, nil, nil, content))
	window.Resize(fyne.NewSize(860, 720))
	// ctx encerra o ticker ao fechar a janela.
	// cancel permite interromper a operação e liberar os recursos do contexto.
	ctx, cancel := context.WithCancel(context.Background())
	window.SetCloseIntercept(func() { closed = true; cancel(); window.SetCloseIntercept(nil); window.Close() })
	update()
	window.Show()
	go func() {
		// ticker apenas agenda trabalho na UI; não lê seu estado mutável.
		ticker := time.NewTicker(40 * time.Millisecond)
		defer ticker.Stop()
		// previous mede tempo real sem incluir intervalos pausados.
		previous := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			// now fornece o instante do temporizador para calcular o avanço da demonstração.
			case now := <-ticker.C:
				// delta conserva o tempo efetivamente decorrido entre quadros.
				delta := int(now.Sub(previous) / time.Millisecond)
				previous = now
				fyne.Do(func() {
					if closed || !playing {
						return
					}
					elapsed += delta
					update()
				})
			}
		}
	}()
}
