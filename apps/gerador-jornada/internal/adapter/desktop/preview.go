package desktop

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// preview abre uma execução independente sem contaminar o histórico de edição.
func (e *Editor) preview() {
	// state é validado e isolado pela aplicação.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	state, err := application.NewPreview(e.Session.Document)
	if !e.check(err) {
		return
	}
	e.showPreview(state, "Prévia · "+e.Session.Document.Project.Name, "Prévia interativa · campos e escolhas são reiniciados a cada visita")
}

// previewCurrentScreen abre somente a tela ativa, sem mudar a tela inicial ou percorrer a jornada.
func (e *Editor) previewCurrentScreen() {
	// state contém uma cópia da tela atual sem ações de navegação entre telas.
	state, err := application.NewScreenPreview(e.Session.Document, e.PageID)
	if !e.check(err) {
		return
	}
	e.showPreview(state, "Prévia Tela Atual · "+state.Document.Screen(state.ScreenID).Name, "Tela atual · campos e escolhas interativos; navegação entre telas desativada")
}

// showPreview compartilha a apresentação visual sem confundir os escopos de tela e projeto.
func (e *Editor) showPreview(state *application.Preview, title, description string) {
	// window é exclusiva da sessão de prévia e não altera a seleção do editor.
	window := e.App.NewWindow(title)
	// redraw reconstrói apenas após uma navegação, não a cada tecla.
	var redraw func()
	redraw = func() {
		window.SetContent(container.NewBorder(widget.NewLabel(description), nil, nil, nil, container.NewScroll(previewPage(state, redraw))))
		// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
		page := state.Document.Screen(state.ScreenID)
		window.Resize(fyne.NewSize(float32(page.Size.Width+48), float32(page.Size.Height+96)))
	}
	redraw()
	window.Show()
}

// previewPage combina primitivas visuais e controles realmente interativos.
func previewPage(state *application.Preview, navigate func()) fyne.CanvasObject {
	// page é a composição de execução atual.
	page := state.Document.Screen(state.ScreenID)
	// objects conserva geometria absoluta e ordem de desenho.
	objects := pageBackground(state.Document, *page)
	// entries permite ao teclado atualizar seu campo de destino.
	entries := map[string]*widget.Entry{}
	// component corresponde a um controle do snapshot de execução.
	for _, component := range page.Components {
		// copy mantém o alvo correto nos callbacks.
		copy := component
		switch component.Kind {
		case catalog.Button:
			// button dispara a ação definida pela jornada.
			button := widget.NewButton(component.Properties.Button.Text, func() {
				if state.Click(copy.ID) == nil {
					navigate()
				}
			})
			button.Importance = widget.HighImportance
			place(button, component.Bounds)
			objects = append(objects, button)
		case catalog.TextField:
			// field aceita valores runtime sem modificar o documento original.
			field := entry(state.Values[component.ID])
			field.OnChanged = func(value string) {
				if state.SetValue(copy.ID, value) != nil {
					field.SetText(state.Values[copy.ID])
				}
			}
			entries[component.ID] = field
			place(field, component.Bounds)
			objects = append(objects, field)
		case catalog.Radio:
			// choices conserva IDs distintos mesmo com legendas iguais.
			choices := container.NewVBox()
			// option é uma seleção identificada do grupo.
			for _, option := range component.Properties.Radio.Options {
				// optionCopy conserva a opção correta no fechamento.
				optionCopy := option
				// button apresenta a marca de seleção atual.
				button := widget.NewButton(option.Text, func() { state.Choices[copy.ID] = optionCopy.ID; navigate() })
				if state.Choices[copy.ID] == option.ID {
					button.Importance = widget.HighImportance
				}
				choices.Add(button)
			}
			if component.Properties.Radio.Orientation == "horizontal" {
				choices.Layout = layout.NewHBoxLayout()
			}
			place(choices, component.Bounds)
			objects = append(objects, choices)
		case catalog.NumericKeyboard:
			// keyboard é montado após os campos para resolver o vínculo.
		default:
			objects = append(objects, previewPrimitives(state.Document, component)...)
		}
	}
	// component fornece teclados que podem referenciar campos posteriores na lista.
	for _, component := range page.Components {
		if component.Kind == catalog.NumericKeyboard {
			// target resolve a entrada da visita, sem armazenar widgets no domínio.
			target := entries[component.Properties.Keyboard.TargetComponentID]
			// keys cria o arranjo das doze teclas.
			keys := container.NewGridWithColumns(3)
			// key corresponde a um comando explícito do teclado.
			for _, key := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "Limpar", "0", "⌫"} {
				// value conserva o comando no callback.
				value := key
				keys.Add(widget.NewButton(key, func() {
					if target == nil {
						return
					}
					if value == "Limpar" {
						target.SetText("")
					} else if value == "⌫" {
						// text mantém o texto editado ou seus caracteres para a operação atual.
						text := []rune(target.Text)
						if len(text) > 0 {
							target.SetText(string(text[:len(text)-1]))
						}
					} else {
						target.SetText(target.Text + value)
					}
				}))
			}
			place(keys, component.Bounds)
			objects = append(objects, keys)
		}
	}
	return container.NewWithoutLayout(objects...)
}

// previewPrimitives expande o conteúdo não interativo sem captura de tela.
func previewPrimitives(document application.Document, component screen.Component) []fyne.CanvasObject {
	// objects contém caixas e textos posicionados.
	objects := []fyne.CanvasObject{}
	// primitive é uma forma independente do renderer de edição.
	for _, primitive := range presentation.Component(component).Shapes {
		if primitive.AssetID != "" {
			if image := resourceImage(document, primitive.AssetID, primitive.Bounds); image != nil {
				objects = append(objects, image)
			}
			continue
		}
		if primitive.Fill != "" || primitive.Stroke != "" {
			// box é um fundo ou contorno editável no modelo.
			box := canvas.NewRectangle(transparent())
			if primitive.Fill != "" {
				box.FillColor = rgb(primitive.Fill)
			}
			if primitive.Stroke != "" {
				box.StrokeColor = rgb(primitive.Stroke)
				box.StrokeWidth = 1
			}
			place(box, primitive.Bounds)
			objects = append(objects, box)
		}
		if strings.TrimSpace(primitive.Text) != "" {
			// text preserva a tipografia lógica da composição.
			text := canvas.NewText(primitive.Text, rgb(primitive.Color))
			text.TextSize = float32(primitive.FontSize)
			text.TextStyle.Bold = primitive.Bold
			setTextAlignment(text, primitive.TextAlign)
			place(text, primitive.Bounds)
			objects = append(objects, text)
		}
	}
	return objects
}

// place aplica bounds lógicos ao widget da sessão de execução.
func place(object fyne.CanvasObject, bounds screen.Bounds) {
	object.Move(fyne.NewPos(float32(bounds.X), float32(bounds.Y)))
	object.Resize(fyne.NewSize(float32(bounds.Width), float32(bounds.Height)))
}
