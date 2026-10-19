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
	e.InspectorActions.RemoveAll()
	if e.page() == nil {
		e.Inspector.Add(widget.NewLabel("Crie uma tela para começar."))
		return
	}
	if count := len(e.selectionIDs()); count > 1 {
		e.groupInspector(count)
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
	// fileID é o nome-base explícito da tela; o UUID interno não é editável pelo formulário.
	fileID := fileIDEntry(page.FileID, true)
	// width mantém a largura informada no formulário até a confirmação.
	width := entry(strconv.Itoa(page.Size.Width))
	// height mantém a altura informada no formulário até a confirmação.
	height := entry(strconv.Itoa(page.Size.Height))
	// background mantém a cor de fundo informada antes de aplicar a alteração.
	background := entry(page.Background.Color)
	// projectName permite renomear o documento independentemente das telas.
	projectName := entry(e.Session.Document.Project.Name)
	// form só publica alterações quando seu botão é acionado.
	form := widget.NewForm(widget.NewFormItem("Projeto", projectName), widget.NewFormItem("Tela", name), widget.NewFormItem("ID", fileID), widget.NewFormItem("Largura", width), widget.NewFormItem("Altura", height), widget.NewFormItem("Fundo", background))
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
			document.Screen(e.PageID).FileID = fileID.Text
			document.Screen(e.PageID).Size = screen.Size{Width: w, Height: h}
			document.Screen(e.PageID).Background.Color = background.Text
			return nil
		})) {
			e.changed()
		}
	}
	e.inspectorActions(form)
	e.Inspector.Add(widget.NewLabelWithStyle("Configuração da tela", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(form)
	e.Inspector.Add(widget.NewLabel("ID do arquivo: Pnl + nome, até 40 caracteres.\nExemplo: PnlIdentificacao → PnlIdentificacao.java"))
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
	// fileID define o nome-base do arquivo do componente, sem exigir Pnl.
	fileID := fileIDEntry(component.FileID, false)
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
	form := widget.NewForm(widget.NewFormItem("Nome", name), widget.NewFormItem("ID", fileID), widget.NewFormItem("X", x), widget.NewFormItem("Y", y), widget.NewFormItem("Largura", width), widget.NewFormItem("Altura", height), widget.NewFormItem("Cor do texto", color), widget.NewFormItem("Fonte (0 = padrão)", font))
	// applyProperties altera apenas a variante de catálogo selecionada.
	applyProperties := e.propertyFields(form, component)
	// alignment edita o parágrafo dentro dos bounds, sem mover o componente.
	alignment := newSelect([]string{"Esquerda", "Centro", "Direita"}, nil)
	alignment.SetSelected("Esquerda")
	if component.Style.TextAlign == "center" {
		alignment.SetSelected("Centro")
	} else if component.Style.TextAlign == "right" {
		alignment.SetSelected("Direita")
	}
	form.Append("Alinhamento do texto", alignment)
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
			target.FileID = fileID.Text
			target.Bounds = screen.Bounds{X: numbers[0], Y: numbers[1], Width: numbers[2], Height: numbers[3]}
			target.Style.TextColor = color.Text
			target.Style.FontSize = numbers[4]
			target.Style.TextAlign = map[string]string{"Esquerda": "left", "Centro": "center", "Direita": "right"}[alignment.Selected]
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			if err := applyProperties(target); err != nil {
				return err
			}
			return page.Validate()
		})) {
			e.changed()
		}
	}
	e.inspectorActions(form)
	e.Inspector.Add(widget.NewLabelWithStyle(familyName(component.Kind), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(form)
	e.Inspector.Add(widget.NewLabel("ID do arquivo: até 40 letras/números, sem acentos.\nPrefixo Pnl não é obrigatório nos componentes."))
	e.Inspector.Add(container.NewGridWithColumns(2, widget.NewButton("Centralizar H", func() { e.align("h") }), widget.NewButton("Centralizar V", func() { e.align("v") }), widget.NewButton("À esquerda", func() { e.align("left") }), widget.NewButton("À direita", func() { e.align("right") }), widget.NewButton("Frente ↑", func() { e.layer(1) }), widget.NewButton("Fundo ↓", func() { e.layer(-1) })))
	e.Inspector.Add(widget.NewButton("Imagem decorativa…", func() { e.importImage(component.ID, false) }))
	e.Inspector.Add(container.NewGridWithColumns(2, widget.NewButton("Duplicar", e.duplicateComponent), widget.NewButton("Remover", e.removeComponent)))
	if component.Kind == catalog.Button {
		e.Inspector.Add(e.navigationControl(component))
	}
}
