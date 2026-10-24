package desktop

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// previewSurface recebe foco e eventos sem desenhar um Entry opaco.
type previewSurface struct {
	widget.BaseWidget
	// onTapped resolve o clique nas áreas internas do composto.
	onTapped func(*fyne.PointEvent)
	// onRune encaminha a digitação para o comando da instância.
	onRune func(rune)
	// onKey trata Enter, Backspace e Escape.
	onKey func(*fyne.KeyEvent)
	// onShortcut trata colagem explícita.
	onShortcut func(fyne.Shortcut)
}

// newPreviewSurface cria uma superfície focável e transparente.
func newPreviewSurface(onTapped func(*fyne.PointEvent), onRune func(rune), onKey func(*fyne.KeyEvent), onShortcut func(fyne.Shortcut)) *previewSurface {
	// surface reutiliza o comportamento base de um widget Fyne.
	surface := &previewSurface{onTapped: onTapped, onRune: onRune, onKey: onKey, onShortcut: onShortcut}
	surface.ExtendBaseWidget(surface)
	return surface
}

// CreateRenderer desenha somente um retângulo transparente.
func (s *previewSurface) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

// Tapped encaminha o ponto relativo ao widget.
func (s *previewSurface) Tapped(event *fyne.PointEvent) {
	if s.onTapped != nil {
		s.onTapped(event)
	}
}

// FocusGained mantém o foco sem alterar o canvas durante a transição.
func (s *previewSurface) FocusGained() {}

// FocusLost mantém o foco sem alterar o canvas durante a transição.
func (s *previewSurface) FocusLost() {}

// TypedRune encaminha um caractere digitado.
func (s *previewSurface) TypedRune(value rune) {
	if s.onRune != nil {
		s.onRune(value)
	}
}

// TypedKey encaminha uma tecla de controle.
func (s *previewSurface) TypedKey(event *fyne.KeyEvent) {
	if s.onKey != nil {
		s.onKey(event)
	}
}

// TypedShortcut encaminha atalhos, inclusive colagem.
func (s *previewSurface) TypedShortcut(shortcut fyne.Shortcut) {
	if s.onShortcut != nil {
		s.onShortcut(shortcut)
	}
}

// insetBounds reserva margens internas sem introduzir dimensões negativas.
func insetBounds(bounds screen.Bounds, x, y int) screen.Bounds {
	return screen.Bounds{X: bounds.X + x, Y: bounds.Y + y, Width: max(1, bounds.Width-2*x), Height: max(1, bounds.Height-2*y)}
}

// drawGroup converte uma cena em objetos Fyne posicionados.
func drawGroup(document application.Document, group presentation.Group) []fyne.CanvasObject {
	// objects mantém a ordem de todas as partes, usando o mesmo renderer do canvas.
	objects := []fyne.CanvasObject{}
	// primitive é a parte atual da cena compartilhada.
	for _, primitive := range group.Shapes {
		objects = append(objects, scenePrimitiveObjects(document, primitive, 1)...)
	}
	return objects
}

// barcodeView desenha a entrada e mantém os objetos sincronizados com o estado.
type barcodeView struct {
	// component contém a instância exibida.
	component screen.Component
	// state contém o estado transitório da aplicação.
	state *application.BarcodeRuntime
	// session é a prévia dona do estado.
	session *application.Preview
	// window recebe o foco e o diálogo de erro.
	window fyne.Window
	// status apresenta o resultado local da instância.
	status *widget.Label
	// tools recebe os comandos explícitos da prévia.
	tools *fyne.Container
	// surface captura digitação, clique e colagem.
	surface *previewSurface
	// rectangles e texts são atualizados sem recriar a página.
	rectangles [8]*canvas.Rectangle
	texts      [8]*canvas.Text
	// message exibe o erro local abaixo dos campos.
	message *canvas.Text
	// objects contém a lista pronta para o contêiner da página.
	objects []fyne.CanvasObject
	// modalOpen impede diálogos duplicados.
	modalOpen bool
}

// newBarcodeView monta a entrada com a geometria compartilhada.
func newBarcodeView(component screen.Component, state *application.BarcodeRuntime, session *application.Preview, window fyne.Window, status *widget.Label, tools *fyne.Container) *barcodeView {
	// view agrupa os objetos e comandos de uma instância.
	view := &barcodeView{component: component, state: state, session: session, window: window, status: status, tools: tools}
	// bounds contém as oito áreas pintadas e clicáveis.
	bounds := presentation.BarcodeFieldBounds(component.Bounds)
	// index percorre os campos para criar fundo e texto.
	for index := range bounds {
		view.rectangles[index] = canvas.NewRectangle(rgb("#FFFFFF"))
		view.rectangles[index].StrokeColor, view.rectangles[index].StrokeWidth = rgb("#CBD5E1"), 1
		place(view.rectangles[index], bounds[index])
		view.texts[index] = canvas.NewText("", rgb("#000000"))
		view.texts[index].TextSize = 12
		setTextAlignment(view.texts[index], "center")
		place(view.texts[index], insetBounds(bounds[index], 4, 2))
	}
	// separators preserva as três marcas da geometria de referência.
	separators := presentation.BarcodeSeparatorBounds(component.Bounds)
	// separatorObjects contém as marcas separadoras.
	separatorObjects := []fyne.CanvasObject{}
	for _, separator := range separators {
		mark := canvas.NewRectangle(rgb("#5064AA"))
		place(mark, separator)
		separatorObjects = append(separatorObjects, mark)
	}
	// label apresenta o rótulo configurável.
	label := canvas.NewText(component.Properties.BarcodeInput.Label, rgb("#5064AA"))
	label.TextSize = 12
	labelBounds := screen.Bounds{X: component.Bounds.X + 30*component.Bounds.Width/presentation.BarcodeInputWidth, Y: component.Bounds.Y + 2*component.Bounds.Height/presentation.BarcodeInputHeight, Width: 200 * component.Bounds.Width / presentation.BarcodeInputWidth, Height: 20 * component.Bounds.Height / presentation.BarcodeInputHeight}
	place(label, labelBounds)
	view.message = canvas.NewText("", rgb("#B91C1C"))
	view.message.TextSize = 10
	place(view.message, screen.Bounds{X: component.Bounds.X, Y: component.Bounds.Y + component.Bounds.Height - 14*component.Bounds.Height/presentation.BarcodeInputHeight, Width: component.Bounds.Width, Height: 12 * component.Bounds.Height / presentation.BarcodeInputHeight})
	// surface fica por último para receber os eventos.
	view.surface = newPreviewSurface(view.tapped, view.typedRune, view.typedKey, view.typedShortcut)
	place(view.surface, component.Bounds)
	// objects preserva a ordem de pintura.
	view.objects = []fyne.CanvasObject{label}
	for index := range bounds {
		view.objects = append(view.objects, view.rectangles[index], view.texts[index])
	}
	view.objects = append(view.objects, separatorObjects...)
	view.objects = append(view.objects, view.message, view.surface)
	view.buildTools()
	view.sync()
	return view
}

// objectsList devolve a lista de formas e a superfície de eventos.
func (v *barcodeView) objectsList() []fyne.CanvasObject { return v.objects }

// buildTools acrescenta os comandos explícitos da entrada.
func (v *barcodeView) buildTools() {
	v.tools.Add(widget.NewButton("Ler · 44 + Enter", func() {
		if err := v.session.BarcodeBeginRead(v.component.ID); err == nil {
			v.status.SetText("Leitura ativa: digite 44 dígitos e pressione Enter.")
			v.focus()
			v.sync()
		}
	}))
	v.tools.Add(widget.NewButton("Cancelar leitura", func() {
		_ = v.session.BarcodeCancelRead(v.component.ID)
		v.status.SetText("Leitura cancelada; o rascunho foi preservado.")
		v.sync()
	}))
	v.tools.Add(widget.NewButton("Colar 44/47", func() {
		_ = v.session.BarcodePaste(v.component.ID, v.window.Clipboard().Content())
		v.afterCommand()
	}))
	v.tools.Add(widget.NewButton("Concluir · Enter", func() {
		_ = v.session.BarcodeFinish(v.component.ID)
		v.afterCommand()
	}))
	v.tools.Add(widget.NewButton("Reiniciar", func() {
		_ = v.session.BarcodeRestart(v.component.ID)
		v.status.SetText("Entrada reiniciada.")
		v.focus()
		v.sync()
	}))
}

// sync atualiza cores, textos e mensagem sem reconstruir a página.
func (v *barcodeView) sync() {
	for index := range v.rectangles {
		fill, textColor, border := "#FFFFFF", "#000000", "#CBD5E1"
		if v.state.Locked {
			fill, textColor = "#E6EEFF", "#808080"
		} else if !v.state.Reading && v.state.Active == index {
			fill, textColor, border = "#10508E", "#FFFFFF", "#10508E"
		}
		v.rectangles[index].FillColor = rgb(fill)
		v.rectangles[index].StrokeColor = rgb(border)
		v.rectangles[index].Refresh()
		v.texts[index].Text = v.state.Values[index]
		v.texts[index].Color = rgb(textColor)
		v.texts[index].Refresh()
	}
	v.message.Text = v.state.Message
	if v.state.Reading {
		v.message.Text = fmt.Sprintf("Leitura: %d/44", len(v.state.Capture))
	}
	v.message.Refresh()
	if v.modalOpen {
		return
	}
	if v.state.Locked {
		v.status.SetText(formatBarcodeStatus(v.state.Result))
	}
}

// tapped seleciona o campo do composto pela posição do clique.
func (v *barcodeView) tapped(event *fyne.PointEvent) {
	v.window.Canvas().Focus(v.surface)
	if v.state.Locked || v.state.Reading {
		return
	}
	// point converte a posição relativa ao widget para coordenadas lógicas absolutas.
	point := fyne.NewPos(float32(v.component.Bounds.X)+event.Position.X, float32(v.component.Bounds.Y)+event.Position.Y)
	// index percorre as áreas dos oito campos.
	for index, bounds := range presentation.BarcodeFieldBounds(v.component.Bounds) {
		if point.X >= float32(bounds.X) && point.X <= float32(bounds.X+bounds.Width) && point.Y >= float32(bounds.Y) && point.Y <= float32(bounds.Y+bounds.Height) {
			// Somente campo liberado recebe foco; blocos posteriores ficam bloqueados.
			if !v.session.BarcodeFieldAllowed(v.component.ID, index) {
				return
			}
			v.state.Active = index
			v.sync()
			return
		}
	}
}

// typedRune encaminha dígitos ou a captura integral.
func (v *barcodeView) typedRune(value rune) {
	_ = v.session.TypeBarcodeDigit(v.component.ID, value)
	v.sync()
}

// typedKey trata Enter, Backspace e Escape.
func (v *barcodeView) typedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyBackspace:
		_ = v.session.BarcodeBackspace(v.component.ID)
		v.sync()
	case fyne.KeyReturn, fyne.KeyEnter:
		if v.state.Reading {
			_ = v.session.BarcodeSubmitRead(v.component.ID, v.state.Capture)
		} else {
			_ = v.session.BarcodeFinish(v.component.ID)
		}
		v.afterCommand()
	case fyne.KeyEscape:
		_ = v.session.BarcodeCancelRead(v.component.ID)
		v.sync()
	}
}

// typedShortcut trata Ctrl+V quando a superfície tem foco.
func (v *barcodeView) typedShortcut(shortcut fyne.Shortcut) {
	if _, ok := shortcut.(*fyne.ShortcutPaste); !ok {
		return
	}
	_ = v.session.BarcodePaste(v.component.ID, v.window.Clipboard().Content())
	v.afterCommand()
}

// afterCommand sincroniza o estado e abre o diálogo de captura inválida.
func (v *barcodeView) afterCommand() {
	v.sync()
	if v.state.Modal && !v.modalOpen {
		v.modalOpen = true
		// errorDialog usa o texto exato contratado e uma única ação.
		errorDialog := dialog.NewCustom("Código de barras inválido", "OK, entendi", widget.NewLabel("O código lido ou colado não foi aceito."), v.window)
		errorDialog.SetOnClosed(func() {
			v.modalOpen = false
			_ = v.session.BarcodeDismissModal(v.component.ID)
			v.focus()
			v.sync()
		})
		errorDialog.Show()
	}
	if !v.modalOpen && v.state.Locked {
		v.status.SetText(formatBarcodeStatus(v.state.Result))
	}
}

// focus coloca o foco na superfície da entrada.
func (v *barcodeView) focus() { v.window.Canvas().Focus(v.surface) }

// menuView desenha um menu e encaminha cliques ao estado local.
type menuView struct {
	// component contém a instância exibida.
	component screen.Component
	// state contém o caminho local da aplicação.
	state *application.MenuRuntime
	// session é a prévia dona do estado.
	session *application.Preview
	// window recebe foco e comandos de saída.
	window fyne.Window
	// status apresenta seleção simulada ou bloqueio.
	status *widget.Label
	// redraw reconstrói a página após abrir um submenu.
	redraw func()
	// objects contém a lista pronta para o contêiner.
	objects []fyne.CanvasObject
}

// newMenuView monta o nó visível do menu da aplicação.
func newMenuView(component screen.Component, session *application.Preview, window fyne.Window, status *widget.Label, redraw func()) *menuView {
	// view agrupa os objetos e comandos de uma instância.
	view := &menuView{component: component, session: session, window: window, status: status, redraw: redraw}
	view.state, _ = session.MenuState(component.ID)
	// node contém o submenu inicial ou o caminho da sessão.
	node := session.MenuNode(component.ID)
	// group mantém faixas e legendas; cada botão recebe seu próprio estado de pressão.
	group := presentation.ApplicationMenuState(component, node)
	if node != nil {
		// decorations exclui partes dos botões que serão pintadas pelos controles transparentes.
		decorations := group
		decorations.Shapes = nil
		// primitive conserva somente a decoração compartilhada externa aos itens.
		for _, primitive := range group.Shapes {
			if primitive.Role == "menu-control" {
				continue
			}
			// isItem impede desenhar duas vezes os pixels semitransparentes da imagem.
			isItem := false
			// item relaciona a forma à identidade estável, inclusive sua legenda autorada.
			for _, item := range node.Items {
				isItem = isItem || primitive.Name == item.ID || strings.HasPrefix(primitive.Name, item.ID+"/")
			}
			if !isItem {
				decorations.Shapes = append(decorations.Shapes, primitive)
			}
		}
		view.objects = drawGroup(session.Document, decorations)
		// area conserva as mesmas regiões de clique usadas no canvas e na cena estática.
		for _, area := range presentation.MenuHitAreas(component, node) {
			// item encontra o conteúdo tipado correspondente à região exibida.
			for _, item := range node.Items {
				if item.ID != area.ItemID {
					continue
				}
				// part recebe geometria local e habilitação combinada de item e composto.
				part := component
				part.Bounds, part.Enabled = area.Bounds, component.Enabled && item.Enabled
				// control conserva ícone, fonte e estados originais sem widget temático opaco.
				control := newLegacyControl(session.Document, part, func() { view.activate(item.ID) })
				control.menuItem = &item
				place(control, area.Bounds)
				view.objects = append(view.objects, control)
				break
			}
		}
	}
	view.buildTools()
	return view
}

// buildTools pinta os comandos locais no próprio rodapé usando a cena e os estados originais.
func (v *menuView) buildTools() {
	// command conserva o mesmo ID e a mesma caixa da cena estática.
	for _, command := range presentation.MenuCommands(v.component, v.session.MenuNode(v.component.ID)) {
		// part e item criam somente um capturador transitório, sem acrescentar itens ao documento.
		part := v.component
		part.Bounds = command.Bounds
		item := catalog.MenuItem{ID: "command/" + command.ID, Label: command.Label, Variant: command.Variant, Enabled: true}
		// control é transparente e compartilha default/touched/disabled do botão original.
		control := newLegacyControl(v.session.Document, part, func() { v.command(command.ID) })
		control.menuItem = &item
		place(control, command.Bounds)
		v.objects = append(v.objects, control)
	}
}

// command executa somente navegação local ou fechamento da prévia, nunca movimentos externos.
func (v *menuView) command(commandID string) {
	switch commandID {
	case "back":
		_ = v.session.MenuBack(v.component.ID)
	case "cancel":
		_ = v.session.MenuCancel(v.component.ID)
	case "exit":
		_ = v.session.MenuExit(v.component.ID)
		v.window.Close()
		return
	}
	v.status.SetText("")
	v.redraw()
}

// activate encaminha a intenção pelo ID; pressione/solte altera somente a imagem transitória.
func (v *menuView) activate(itemID string) {
	// outcome informa navegação, simulação ou bloqueio sem executar movimentos externos.
	outcome, err := v.session.MenuActivate(v.component.ID, itemID)
	if err != nil {
		v.status.SetText(err.Error())
		if errors.Is(err, application.ErrJourneyMissing) && v.session.OnMissingJourney != nil {
			v.session.OnMissingJourney(v.component.ID, itemID)
		}
		return
	}
	switch outcome {
	case application.MenuOpened, application.MenuNavigated:
		v.redraw()
	case application.MenuSimulated:
		// item mantém a distinção nomeMovimento entre normal/especial que compartilham código.
		item, ok := v.session.MenuSelectedItem(v.component.ID)
		if ok {
			// details usa somente metadados fornecidos; código/atalho não são rotas executáveis.
			details := "Seleção simulada: " + strings.TrimSpace(item.Label)
			if item.MovementName != "" {
				details += " · movimento: " + item.MovementName
			}
			if item.MovementCode != "" {
				details += " · código: " + item.MovementCode
			}
			if item.MovementOption != "" {
				details += " · opção: " + item.MovementOption
			}
			if item.TargetScreenID != "" || item.TargetJourney != "" {
				v.status.SetText(details + " · navegação entre telas desativada nesta prévia.")
			} else {
				v.status.SetText(details + " · nenhum destino configurado.")
			}
			if item.TargetScreenID == "" && item.TargetJourney == "" && v.session.OnMissingJourney != nil {
				v.session.OnMissingJourney(v.component.ID, itemID)
			}
		}
	case application.MenuBlocked:
		v.status.SetText("Item desabilitado; nenhuma ação executada.")
	}
}

// formatBarcodeStatus apresenta somente dados derivados, sem repetir o código completo.
func formatBarcodeStatus(barcode catalog.Barcode) string {
	// caixa descreve a classificação local, sem alegar homologação bancária.
	caixa := "não"
	if barcode.IsCaixa {
		caixa = "sim (premissa local não homologada)"
	}
	// factor descreve a ausência do fator sem inventar data de vencimento.
	factor := fmt.Sprintf("%d", barcode.Factor)
	if !barcode.HasFactor {
		factor = "ausente"
	}
	// sigcb descreve a classificação sob a premissa local, sem alegar homologação.
	sigcb := "não"
	if barcode.Sigcb {
		sigcb = "sim (premissa local não homologada)"
	}
	// title só existe quando a classificação SIGCB é verdadeira.
	title := barcode.TitleNumber
	if title == "" {
		title = "indisponível"
	}
	return fmt.Sprintf("Aceito · banco %s · valor R$ %.2f · fator %s · CAIXA: %s · SIGCB: %s · título %s", barcode.Bank, float64(barcode.ValueCents)/100, factor, caixa, sigcb, title)
}
