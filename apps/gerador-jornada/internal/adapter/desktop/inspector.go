package desktop

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// entry cria um campo sem disparar comandos durante a montagem do inspector.
func entry(value string) *widget.Entry {
	// field mantém o campo de entrada usado na edição da propriedade.
	field := widget.NewEntry()
	field.SetText(value)
	return field
}

// refreshInspector seleciona propriedades de tela ou do controle ativo.
func (e *Editor) refreshInspector() {
	e.Inspector.RemoveAll()
	if e.page() == nil {
		e.Inspector.Add(widget.NewLabel("Crie uma tela para começar."))
		return
	}
	// component identifica o controle da tela que será consultado ou alterado.
	if component := e.page().ComponentByID(e.Selected); component != nil {
		e.componentInspector(*component)
	} else {
		e.screenInspector()
	}
}

// screenInspector permite editar a tela sem gerar título ou cabeçalho ocultos.
func (e *Editor) screenInspector() {
	// page fornece os valores que estão sendo editados.
	page := e.page()
	// name, width, height e background mantêm o rascunho do formulário.
	name := entry(page.Name)
	// width mantém a largura informada no formulário até a confirmação.
	width := entry(strconv.Itoa(page.Size.Width))
	// height mantém a altura informada no formulário até a confirmação.
	height := entry(strconv.Itoa(page.Size.Height))
	// background mantém a cor de fundo informada antes de aplicar a alteração.
	background := entry(page.Background.Color)
	// projectName permite renomear o documento independentemente das telas.
	projectName := entry(e.Session.Document.Project.Name)
	// form só publica alterações quando seu botão é acionado.
	form := widget.NewForm(widget.NewFormItem("Projeto", projectName), widget.NewFormItem("Tela", name), widget.NewFormItem("Largura", width), widget.NewFormItem("Altura", height), widget.NewFormItem("Fundo", background))
	form.SubmitText = "Aplicar à tela"
	form.CancelText = ""
	form.OnSubmit = func() {
		// w e h são validados antes de publicar qualquer campo.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		w, err := strconv.Atoi(width.Text)
		if !e.check(err) {
			return
		}
		// h contém a altura inteira validada antes da transação.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		h, err := strconv.Atoi(height.Text)
		if !e.check(err) {
			return
		}
		if w < 1 || h < 1 || w > 8192 || h > 8192 || strings.TrimSpace(name.Text) == "" || strings.TrimSpace(projectName.Text) == "" || !catalog.ValidColor(background.Text) {
			e.check(fmt.Errorf("nome, dimensões ou cor inválidos"))
			return
		}
		if e.check(e.Session.Change(func(document *application.Document) error {
			document.Project.Name = projectName.Text
			document.Screen(e.PageID).Name = name.Text
			document.Screen(e.PageID).Size = screen.Size{Width: w, Height: h}
			document.Screen(e.PageID).Background.Color = background.Text
			return nil
		})) {
			e.changed()
		}
	}
	e.Inspector.Add(widget.NewLabelWithStyle("Configuração da tela", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(form)
	e.Inspector.Add(widget.NewButton("Definir como tela inicial", func() {
		e.check(e.Session.Change(func(document *application.Document) error { document.Journey.StartScreenID = e.PageID; return nil }))
		e.changed()
	}))
	e.Inspector.Add(widget.NewButton("Fundo da tela…", e.configureBackground))
	e.Inspector.Add(widget.NewLabel("Selecione um componente\npara editar suas propriedades."))
}

// componentInspector cria campos tipados sem concentrar geração ou persistência na UI.
func (e *Editor) componentInspector(component screen.Component) {
	// name, x, y, width e height são o rascunho geométrico.
	name := entry(component.Name)
	// x mantém a coordenada horizontal editada no formulário.
	x := entry(strconv.Itoa(component.Bounds.X))
	// y mantém a coordenada vertical editada no formulário.
	y := entry(strconv.Itoa(component.Bounds.Y))
	// width mantém a largura informada no formulário até a confirmação.
	width := entry(strconv.Itoa(component.Bounds.Width))
	// height mantém a altura informada no formulário até a confirmação.
	height := entry(strconv.Itoa(component.Bounds.Height))
	// color e font são overrides opcionais.
	color := entry(component.Style.TextColor)
	// font mantém o tamanho de fonte informado antes da validação.
	font := entry(strconv.Itoa(component.Style.FontSize))
	// form mantém os campos comuns.
	form := widget.NewForm(widget.NewFormItem("Nome", name), widget.NewFormItem("X", x), widget.NewFormItem("Y", y), widget.NewFormItem("Largura", width), widget.NewFormItem("Altura", height), widget.NewFormItem("Cor do texto", color), widget.NewFormItem("Fonte (0 = padrão)", font))
	// applyProperties altera apenas a variante de catálogo selecionada.
	applyProperties := e.propertyFields(form, component)
	form.SubmitText = "Aplicar alterações"
	form.CancelText = ""
	form.OnSubmit = func() {
		// numbers contém os valores geométricos convertidos antes da transação.
		numbers := []int{}
		// field é um valor inteiro informado no formulário.
		for _, field := range []*widget.Entry{x, y, width, height, font} {
			// number contém o valor inteiro convertido de uma propriedade geométrica.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			number, err := strconv.Atoi(field.Text)
			if !e.check(err) {
				return
			}
			numbers = append(numbers, number)
		}
		if color.Text != "" && !catalog.ValidColor(color.Text) {
			e.check(fmt.Errorf("use cor no formato #RRGGBB"))
			return
		}
		if e.check(e.Session.Change(func(document *application.Document) error {
			// page pertence ao snapshot provisório.
			page := document.Screen(e.PageID)
			// target é o controle editado.
			target := page.ComponentByID(component.ID)
			if target == nil {
				return fmt.Errorf("componente inexistente")
			}
			target.Name = name.Text
			target.Bounds = screen.Bounds{X: numbers[0], Y: numbers[1], Width: numbers[2], Height: numbers[3]}
			target.Style.TextColor = color.Text
			target.Style.FontSize = numbers[4]
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			if err := applyProperties(target); err != nil {
				return err
			}
			return page.Validate()
		})) {
			e.changed()
		}
	}
	e.Inspector.Add(widget.NewLabelWithStyle(familyName(component.Kind), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(form)
	e.Inspector.Add(container.NewGridWithColumns(2, widget.NewButton("Centralizar H", func() { e.align("h") }), widget.NewButton("Centralizar V", func() { e.align("v") }), widget.NewButton("À esquerda", func() { e.align("left") }), widget.NewButton("À direita", func() { e.align("right") }), widget.NewButton("Frente ↑", func() { e.layer(1) }), widget.NewButton("Fundo ↓", func() { e.layer(-1) })))
	e.Inspector.Add(widget.NewButton("Imagem decorativa…", func() { e.importImage(component.ID, false) }))
	e.Inspector.Add(container.NewGridWithColumns(2, widget.NewButton("Duplicar", e.duplicateComponent), widget.NewButton("Remover", e.removeComponent)))
	if component.Kind == catalog.Button {
		e.Inspector.Add(e.navigationControl(component))
	}
}

// propertyFields delega a edição de conteúdo à variante correspondente.
func (e *Editor) propertyFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// content conserva texto simples ou valor inicial.
	content := entry("")
	switch component.Kind {
	case catalog.Button:
		content.SetText(component.Properties.Button.Text)
		form.Append("Texto", content)
		// variant restringe a escolha às skins previstas.
		variant := widget.NewSelect([]string{"confirm", "cancel", "yes", "no"}, nil)
		variant.SetSelected(component.Properties.Button.Variant)
		form.Append("Tipo", variant)
		return func(target *screen.Component) error {
			target.Properties.Button.Text = content.Text
			target.Properties.Button.Variant = variant.Selected
			return nil
		}
	case catalog.Label:
		content.SetText(component.Properties.Label.Text)
		form.Append("Texto", content)
		return func(target *screen.Component) error { target.Properties.Label.Text = content.Text; return nil }
	case catalog.Title:
		content.SetText(component.Properties.Title.Text)
		form.Append("Texto", content)
		return func(target *screen.Component) error { target.Properties.Title.Text = content.Text; return nil }
	case catalog.TextField:
		content.SetText(component.Properties.Field.InitialValue)
		form.Append("Valor inicial", content)
		// format preserva os nomes de intenção sem validar CPF como regra fiscal.
		format := widget.NewSelect([]string{"text", "number", "cellular", "cep", "cpf", "date", "ddd", "nis"}, nil)
		format.SetSelected(component.Properties.Field.Format)
		form.Append("Formato", format)
		return func(target *screen.Component) error {
			target.Properties.Field.InitialValue = content.Text
			target.Properties.Field.Format = format.Selected
			return nil
		}
	case catalog.Radio, catalog.Report:
		// lines permite vírgulas dentro de uma opção.
		lines := widget.NewMultiLineEntry()
		lines.SetMinRowsVisible(3)
		// title é aplicado somente a relatórios.
		title := entry("")
		// options preserva IDs de linhas existentes enquanto possível.
		options := []catalog.Option{}
		if component.Properties.Radio != nil {
			options = component.Properties.Radio.Options
		} else {
			options = component.Properties.Report.Rows
			title.SetText(component.Properties.Report.Title)
			form.Append("Cabeçalho", title)
		}
		// texts fornece uma linha por item, sem usar vírgulas como separador.
		texts := []string{}
		// option fornece o texto de uma opção preservando sua ordem.
		for _, option := range options {
			texts = append(texts, option.Text)
		}
		lines.SetText(strings.Join(texts, "\n"))
		form.Append("Uma opção por linha", lines)
		// orientation só afeta seleção.
		orientation := widget.NewSelect([]string{"vertical", "horizontal"}, nil)
		if component.Properties.Radio != nil {
			orientation.SetSelected(component.Properties.Radio.Orientation)
			form.Append("Orientação", orientation)
		}
		return func(target *screen.Component) error {
			// result mantém ordem e identidade das opções não removidas.
			result := []catalog.Option{}
			if lines.Text != "" {
				// index localiza o item atual ou sua posição de destino na coleção ordenada.
				// text mantém o texto editado ou seus caracteres para a operação atual.
				for index, text := range strings.Split(lines.Text, "\n") {
					// id identifica de forma estável o item ou recurso atual.
					id := project.NewID()
					if index < len(options) {
						id = options[index].ID
					}
					result = append(result, catalog.Option{ID: id, Text: text})
				}
			}
			if target.Properties.Radio != nil {
				target.Properties.Radio.Options = result
				target.Properties.Radio.Orientation = orientation.Selected
				target.Properties.Radio.SelectedOptionID = ""
			} else {
				target.Properties.Report.Rows = result
				target.Properties.Report.Title = title.Text
			}
			return nil
		}
	case catalog.NumericKeyboard:
		// targets é o vocabulário visível dos campos da mesma tela.
		targets := map[string]string{"Sem vínculo": ""}
		// labels reúne as legendas exibidas sem expor identificadores internos.
		labels := []string{"Sem vínculo"}
		// selected mantém a legenda selecionada no formulário.
		selected := "Sem vínculo"
		// field é um possível destino de entrada.
		for _, field := range e.page().Components {
			if field.Kind == catalog.TextField {
				labels = append(labels, field.Name)
				targets[field.Name] = field.ID
				if field.ID == component.Properties.Keyboard.TargetComponentID {
					selected = field.Name
				}
			}
		}
		// target permite vincular o teclado sem expor IDs.
		target := widget.NewSelect(labels, nil)
		target.SetSelected(selected)
		form.Append("Campo de destino", target)
		return func(control *screen.Component) error {
			control.Properties.Keyboard.TargetComponentID = targets[target.Selected]
			return nil
		}
	}
	return func(*screen.Component) error { return nil }
}

// navigationControl edita o binding do botão sem serializar callbacks.
func (e *Editor) navigationControl(component screen.Component) fyne.CanvasObject {
	// options e targets separam legendas de identidades persistidas.
	options := []string{"Sem navegação"}
	// targets associa as legendas visíveis aos identificadores dos destinos.
	targets := map[string]string{"Sem navegação": ""}
	// selected mantém a legenda selecionada no formulário.
	selected := "Sem navegação"
	// page é um destino existente, com índice para diferenciar nomes iguais.
	// index localiza o item atual ou sua posição de destino na coleção ordenada.
	for index, page := range e.Session.Document.Screens {
		// label diferencia nomes iguais com um índice visível.
		label := fmt.Sprintf("%d · %s", index+1, page.Name)
		options = append(options, label)
		targets[label] = page.ID
		if page.ID == e.Session.Document.Journey.Destination(component.ID) {
			selected = label
		}
	}
	// destination só confirma a alteração depois da seleção do usuário.
	destination := widget.NewSelect(options, nil)
	destination.SetSelected(selected)
	destination.OnChanged = func(value string) {
		e.check(e.Session.Change(func(document *application.Document) error {
			// kind seleciona o discriminador da ação conforme o destino informado.
			kind := "navigate"
			if targets[value] == "" {
				kind = "none"
			}
			return document.Journey.Bind(journey.Action{SourceScreenID: e.PageID, SourceComponentID: component.ID, Kind: kind, TargetScreenID: targets[value]})
		}))
		e.changed()
	}
	return container.NewVBox(section("AO CLICAR"), destination)
}

// align publica um ajuste em relação à superfície da tela.
func (e *Editor) align(kind string) {
	e.check(e.Session.Change(func(document *application.Document) error {
		// page e component são os alvos da operação reversível.
		page := document.Screen(e.PageID)
		// component identifica o controle da tela que será consultado ou alterado.
		component := page.ComponentByID(e.Selected)
		if component == nil {
			return nil
		}
		switch kind {
		case "h":
			component.Bounds.X = (page.Size.Width - component.Bounds.Width) / 2
		case "v":
			component.Bounds.Y = (page.Size.Height - component.Bounds.Height) / 2
		case "left":
			component.Bounds.X = 0
		case "right":
			component.Bounds.X = page.Size.Width - component.Bounds.Width
		}
		return nil
	}))
	e.changed()
}

// layer move somente uma posição na ordem de camadas.
func (e *Editor) layer(delta int) {
	e.check(e.Session.Change(func(document *application.Document) error {
		// page contém o z-order atual.
		page := document.Screen(e.PageID)
		// index procura o componente selecionado.
		for index := range page.Components {
			if page.Components[index].ID == e.Selected {
				// next calcula a posição pretendida e permite verificar os limites antes de mover.
				if next := index + delta; next >= 0 && next < len(page.Components) {
					page.Components[index], page.Components[next] = page.Components[next], page.Components[index]
				}
				break
			}
		}
		return nil
	}))
	e.changed()
}
