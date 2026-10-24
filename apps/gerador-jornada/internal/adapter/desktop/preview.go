package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	// Log registra somente o escopo aberto, sem nomes ou conteúdo do projeto.
	LogApplication("preview_open", map[string]any{"scope": "project"})
	e.showPreview(state, "Prévia · "+e.Session.Document.Project.Name, "Prévia interativa · campos e escolhas são reiniciados a cada visita")
}

// previewCurrentScreen abre somente a tela ativa, sem mudar a tela inicial ou percorrer a jornada.
func (e *Editor) previewCurrentScreen() {
	// state contém uma cópia da tela atual sem ações de navegação entre telas.
	state, err := application.NewScreenPreview(e.Session.Document, e.PageID)
	if !e.check(err) {
		return
	}
	// Log registra somente o escopo aberto, sem nomes ou conteúdo do projeto.
	LogApplication("preview_open", map[string]any{"scope": "screen"})
	e.showPreview(state, "Prévia Tela Atual · "+state.Document.Screen(state.ScreenID).Name, "Tela atual · campos e escolhas interativos; navegação entre telas desativada")
}

// showPreview compartilha a apresentação visual sem confundir os escopos de tela e projeto.
func (e *Editor) showPreview(state *application.Preview, title, description string) {
	// window é exclusiva da sessão de prévia e não altera a seleção do editor.
	window := e.App.NewWindow(title)
	// status apresenta mensagens locais sem reabrir o diálogo de erro do editor.
	status := widget.NewLabel(description)
	// tools reúne os comandos explícitos dos compostos da tela.
	tools := container.NewHBox()
	// closed impede callbacks de tocar um canvas destruído.
	closed := false
	// redraw reconstrói apenas após uma navegação ou troca de submenu.
	var redraw func()
	redraw = func() {
		if closed {
			return
		}
		// Uma nova página precisa recriar os comandos do nó atual.
		tools.RemoveAll()
		// surface isola controles runtime do tema do editor, preservando o visual claro da prévia.
		surface := container.NewThemeOverride(previewPage(state, window, redraw, status, tools), jornadaTheme{})
		window.SetContent(container.NewBorder(container.NewVBox(status, container.NewHScroll(tools)), nil, nil, nil, container.NewScroll(surface)))
		// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
		page := state.Document.Screen(state.ScreenID)
		window.Resize(fyne.NewSize(float32(page.Size.Width+48), float32(page.Size.Height+96)))
	}
	redraw()
	window.SetOnClosed(func() { closed = true })
	window.Show()
}

// previewPage combina primitivas visuais e controles realmente interativos.
func previewPage(state *application.Preview, window fyne.Window, redraw func(), status *widget.Label, tools *fyne.Container) fyne.CanvasObject {
	// page é a composição de execução isolada da sessão do editor.
	page := state.Document.Screen(state.ScreenID)
	// objects preserva a ordem de camadas inclusive para teclados e componentes sobrepostos.
	objects := pageBackground(state.Document, *page)
	// entries e surfaces resolvem vínculos sem depender da ordem dos componentes.
	entries := map[string]*legacyEntry{}
	surfaces := map[string]fyne.CanvasObject{}
	// barcodes permite ao teclado virtual operar o mesmo estado da entrada física.
	barcodes := map[string]*barcodeView{}
	// fields mantém a ordem dos campos habilitados para o comando TAB.
	fields := []*legacyEntry{}
	// component cria capturadores sem acrescentá-los ainda à ordem de pintura.
	for _, component := range page.Components {
		if !component.Enabled || presentation.ValidateStyle(component) != nil {
			continue
		}
		// copy conserva a identidade usada pelos callbacks.
		copy := component
		switch component.Kind {
		case catalog.TextField:
			// field e surface mantêm a captura transparente sobre a decoração da cena.
			field, surface := newLegacyEntry(state.Document, component, state.Values[component.ID])
			field.OnChanged = func(value string) {
				if err := state.SetValue(copy.ID, value); err != nil {
					status.SetText(err.Error())
					if field.Text != state.Values[copy.ID] {
						field.SetText(state.Values[copy.ID])
					}
				}
			}
			entries[component.ID], surfaces[component.ID] = field, surface
			fields = append(fields, field)
		case catalog.BarcodeInput:
			// runtime já foi inicializado pela aplicação, sem publicar valores no documento.
			runtime, ok := state.Barcode(component.ID)
			if !ok {
				runtime = &application.BarcodeRuntime{}
			}
			barcodes[component.ID] = newBarcodeView(copy, runtime, state, window, status, tools)
		}
	}
	// component acrescenta cada controle exatamente na posição canônica de sua camada.
	for _, component := range page.Components {
		// copy impede que callbacks compartilhem o componente da iteração seguinte.
		copy := component
		if !component.Enabled || presentation.ValidateStyle(component) != nil {
			objects = append(objects, previewPrimitives(state.Document, component)...)
			continue
		}
		switch component.Kind {
		case catalog.Button:
			// button mantém a cena original enquanto a aplicação executa a navegação.
			button := newLegacyControl(state.Document, component, func() {
				if state.Click(copy.ID) == nil {
					redraw()
				}
			})
			place(button, component.Bounds)
			objects = append(objects, button)
		case catalog.TextField:
			objects = append(objects, surfaces[component.ID])
		case catalog.Radio:
			objects = append(objects, previewRadio(state, component)...)
		case catalog.NumericKeyboard:
			objects = append(objects, previewKeyboard(state, component, window, entries, fields, barcodes, status)...)
		case catalog.BarcodeInput:
			objects = append(objects, barcodes[component.ID].objectsList()...)
		case catalog.ApplicationMenu:
			objects = append(objects, newMenuView(copy, state, window, status, tools, redraw).objects...)
		default:
			objects = append(objects, previewPrimitives(state.Document, component)...)
		}
	}
	return container.NewWithoutLayout(objects...)
}

// previewRadio compartilha a pintura estática e atualiza somente a seleção da visita.
func previewRadio(state *application.Preview, component screen.Component) []fyne.CanvasObject {
	// objects mantém as regiões derivadas no mesmo nível de camadas do grupo original.
	objects := []fyne.CanvasObject{}
	// controls permite atualizar as marcas sem reconstruir a janela e perder foco/cursor.
	controls := []*legacyControl{}
	// index e option percorrem as escolhas identificadas, inclusive legendas iguais.
	for index, option := range component.Properties.Radio.Options {
		// part é uma única opção derivada do grupo, sem nova identidade persistida.
		part := component
		part.Bounds = presentation.RadioOptionBounds(component, index)
		part.Properties = catalog.Properties{Radio: &catalog.RadioProperties{Options: []catalog.Option{option}, SelectedOptionID: state.Choices[component.ID]}}
		// optionID conserva a escolha correta no callback.
		optionID := option.ID
		// control desenha a mesma marca/cena de uma opção no canvas e no PPT.
		control := newLegacyControl(state.Document, part, func() {
			state.Choices[component.ID] = optionID
			// choice atualiza somente os controles desta visita.
			for _, choice := range controls {
				choice.component.Properties.Radio.SelectedOptionID = optionID
				choice.Refresh()
			}
		})
		place(control, part.Bounds)
		controls = append(controls, control)
		objects = append(objects, control)
	}
	return objects
}

// previewKeyboard usa as mesmas quatorze regiões da cena para clique, foco e cursor.
func previewKeyboard(state *application.Preview, component screen.Component, window fyne.Window, entries map[string]*legacyEntry, fields []*legacyEntry, barcodes map[string]*barcodeView, status *widget.Label) []fyne.CanvasObject {
	// objects contém apenas as partes interativas desta instância.
	objects := []fyne.CanvasObject{}
	// key fornece posição e skin sem duplicar o layout Java no adaptador.
	for _, key := range presentation.KeyboardKeys(component.Bounds) {
		// part conserva overrides explícitos e a definição de cada tecla auxiliar.
		part := component
		part.Kind, part.Bounds, part.Style.Skin = catalog.Button, key.Bounds, key.Skin
		if component.Style.Skin != "" && !key.Auxiliary {
			part.Style.Skin = component.Style.Skin
		}
		part.Properties = catalog.Properties{Button: &catalog.ButtonProperties{Text: key.Label, Variant: "confirm"}}
		// label identifica o comando sem confundi-lo com um dígito ASCII.
		label := key.Label
		// control permanece fora do foco para não deslocar o cursor do campo vinculado.
		control := newLegacyControl(state.Document, part, func() {
			// barcodeTarget tem prioridade quando o vínculo aponta para o composto.
			barcodeTarget := barcodes[component.Properties.Keyboard.TargetComponentID]
			if barcodeTarget != nil {
				switch label {
				case "LIMPAR":
					barcodeTarget.typedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
				case "<", ">", "TAB":
					// direction seleciona outro subcampo habilitado sem editar os valores.
					direction := 1
					if label == "<" {
						direction = -1
					}
					// index procura o próximo subcampo permitido respeitando os limites.
					for index := barcodeTarget.state.Active + direction; index >= 0 && index < 8; index += direction {
						if state.BarcodeFieldAllowed(barcodeTarget.component.ID, index) {
							barcodeTarget.state.Active = index
							break
						}
					}
					barcodeTarget.sync()
				default:
					barcodeTarget.typedRune(rune(label[0]))
				}
				barcodeTarget.focus()
				return
			}
			// target usa o vínculo explícito e acompanha o foco após TAB na mesma visita.
			target := entries[component.Properties.Keyboard.TargetComponentID]
			if target == nil {
				status.SetText("Vincule o teclado a um campo habilitado desta tela.")
				return
			}
			if focused, ok := window.Canvas().Focused().(*legacyEntry); ok {
				if entries[focused.component.ID] == focused {
					target = focused
				}
			}
			window.Canvas().Focus(target)
			switch label {
			case "<":
				target.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
			case ">":
				target.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
			case "LIMPAR":
				target.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
			case "TAB":
				// index encontra o próximo campo habilitado sem focar as teclas virtuais.
				for index, field := range fields {
					if field == target && len(fields) > 0 {
						window.Canvas().Focus(fields[(index+1)%len(fields)])
						break
					}
				}
			default:
				target.TypedRune(rune(label[0]))
			}
		})
		control.takeFocus = false
		place(control, key.Bounds)
		objects = append(objects, control)
	}
	return objects
}

// previewPrimitives expande o conteúdo não interativo sem captura de tela.
func previewPrimitives(document application.Document, component screen.Component) []fyne.CanvasObject {
	return drawGroup(document, presentation.Component(component))
}

// place aplica bounds lógicos ao widget da sessão de execução.
func place(object fyne.CanvasObject, bounds screen.Bounds) {
	object.Move(fyne.NewPos(float32(bounds.X), float32(bounds.Y)))
	object.Resize(fyne.NewSize(float32(bounds.Width), float32(bounds.Height)))
}
