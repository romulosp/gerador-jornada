package desktop

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
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
	e.InspectorHasDraft = nil
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
			// Nome histórico inalterado não é truncado; nova entrada segue o limite comum.
			if document.Project.Name != projectName.Text {
				if err := application.ValidateAuthoringName(projectName.Text); err != nil {
					return err
				}
			}
			document.Project.Name = projectName.Text
			return application.ConfigureScreen(document, e.PageID, application.ScreenChanges{Name: &name.Text, FileID: &fileID.Text, Width: &w, Height: &h, BackgroundColor: &background.Text})
		})) {
			e.changed()
		}
	}
	e.inspectorActions(form)
	e.Inspector.Add(widget.NewLabelWithStyle("Configuração da tela", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(form)
	e.Inspector.Add(widget.NewLabel("ID do arquivo: Pnl + nome, até 40 caracteres.\nExemplo: PnlIdentificacao → PnlIdentificacao.java"))
	e.Inspector.Add(widget.NewButton("Definir como tela inicial", func() {
		e.check(e.Session.Change(func(document *application.Document) error { return application.SetJourneyStart(document, e.PageID) }))
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
	// fill e border são overrides opcionais; vazio preserva a decoração original.
	fill, border := entry(component.Style.FillColor), entry(component.Style.BorderColor)
	form.Append("Cor de fundo (opcional)", fill)
	form.Append("Cor da borda (opcional)", border)
	// face oferece nomes reais dos recursos incorporados e conserva valores de arquivos anteriores.
	faceOptions := []string{"Padrão da skin", "Arial", "Futura Medium", "Futura Hv BT", "Futura Md BT", "Futura BdCn BT", "Bitstream Vera Sans Mono"}
	if component.Style.FontFamily != "" {
		// found identifica uma face preservada que não pertence à biblioteca atual.
		found := false
		// option compara o nome real sem modificar a preferência persistida.
		for _, option := range faceOptions {
			found = found || option == component.Style.FontFamily
		}
		if !found {
			faceOptions = append(faceOptions, component.Style.FontFamily)
		}
	}
	// face mantém o rascunho da família até Salvar alterações.
	face := widget.NewSelect(faceOptions, nil)
	face.SetSelected("Padrão da skin")
	if component.Style.FontFamily != "" {
		face.SetSelected(component.Style.FontFamily)
	}
	form.Append("Família da fonte", face)
	// applyProperties altera apenas a variante de catálogo selecionada.
	applyProperties := e.propertyFields(form, component)
	// alignment edita o parágrafo dentro dos bounds, sem mover o componente.
	alignment := newSelect([]string{"Padrão da skin", "Esquerda", "Centro", "Direita"}, nil)
	alignment.SetSelected("Padrão da skin")
	if component.Style.TextAlign == "left" {
		alignment.SetSelected("Esquerda")
	} else if component.Style.TextAlign == "center" {
		alignment.SetSelected("Centro")
	} else if component.Style.TextAlign == "right" {
		alignment.SetSelected("Direita")
	}
	form.Append("Alinhamento do texto", alignment)
	// skins oferece as definições do acervo legado compatíveis com a família.
	skinOptions := append([]string{"Padrão do legado"}, skin.MustBuiltin().Names(component.Kind)...)
	skinSelect := widget.NewSelect(skinOptions, nil)
	skinSelect.SetSelected("Padrão do legado")
	if component.Style.Skin != "" {
		skinSelect.SetSelected(component.Style.Skin)
	}
	form.Append("Skin", skinSelect)
	// enabled controla a interação e o estado visual do componente.
	enabled := widget.NewCheck("Habilitado", nil)
	enabled.SetChecked(component.Enabled)
	form.Append("Habilitado", enabled)
	// parentChoices e parentIDs excluem o próprio nó e seus descendentes.
	parentChoices, parentIDs := parentOptions(*e.page(), component)
	parent := widget.NewSelect(parentChoices, nil)
	parent.SetSelected("Nenhum")
	// label e id percorrem as opções para refletir o pai confirmado.
	for label, id := range parentIDs {
		if id == component.ParentID {
			parent.SetSelected(label)
		}
	}
	form.Append("Pai", parent)
	// layer e total apresentam a posição canônica fundo→frente.
	layer, total := currentLayer(*e.page(), component.ID)
	form.Append("Camada", widget.NewLabel(fmt.Sprintf("%d de %d (1 = fundo)", layer, total)))
	// mode oferece os cinco modos contratados; referência só vale para os relativos.
	mode := widget.NewSelect([]string{"Manter posição atual", "À frente do objeto", "Atrás do objeto", "Trazer para a frente", "Enviar para o fundo"}, nil)
	mode.SetSelected("Manter posição atual")
	referenceChoices, referenceIDs := layerOptions(*e.page(), component)
	reference := widget.NewSelect(referenceChoices, nil)
	reference.SetSelected("Selecione um objeto")
	reference.Disable()
	mode.OnChanged = func(value string) {
		if value == "À frente do objeto" || value == "Atrás do objeto" {
			reference.Enable()
			return
		}
		reference.Disable()
	}
	form.Append("Sobreposição", mode)
	form.Append("Objeto de referência", reference)
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
		if (color.Text != "" && !catalog.ValidColor(color.Text)) || (fill.Text != "" && !catalog.ValidColor(fill.Text)) || (border.Text != "" && !catalog.ValidColor(border.Text)) {
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
			// before conserva propriedades anteriores para validar apenas conteúdo realmente alterado.
			before := *target
			before.Properties = target.Properties.Copy()
			target.Name = name.Text
			target.FileID = fileID.Text
			target.Bounds = screen.Bounds{X: numbers[0], Y: numbers[1], Width: numbers[2], Height: numbers[3]}
			target.Style.TextColor = color.Text
			target.Style.FontSize = numbers[4]
			target.Style.FillColor, target.Style.BorderColor = fill.Text, border.Text
			target.Style.FontFamily = ""
			if face.Selected != "Padrão da skin" {
				target.Style.FontFamily = face.Selected
			}
			target.Style.TextAlign = map[string]string{"Padrão da skin": "", "Esquerda": "left", "Centro": "center", "Direita": "right"}[alignment.Selected]
			target.Style.Skin = ""
			if skinSelect.Selected != "Padrão do legado" {
				target.Style.Skin = skinSelect.Selected
			}
			target.Enabled = enabled.Checked
			target.ParentID = parentIDs[parent.Selected]
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			if err := applyProperties(target); err != nil {
				return err
			}
			if err := application.ValidateComponentEdit(document, page, before, *target); err != nil {
				return err
			}
			if err := builtinassets.ValidateAppearanceEdit(before, *target); err != nil {
				return err
			}
			switch mode.Selected {
			case "À frente do objeto", "Atrás do objeto":
				// referência relativa precisa existir na mesma tela.
				referenceID := referenceIDs[reference.Selected]
				if referenceID == "" {
					return fmt.Errorf("selecione um objeto de referência")
				}
				if err := screen.PlaceComponent(page, target.ID, referenceID, mode.Selected == "À frente do objeto"); err != nil {
					return err
				}
			case "Trazer para a frente":
				if err := screen.PlaceComponent(page, target.ID, "", true); err != nil {
					return err
				}
			case "Enviar para o fundo":
				if err := screen.PlaceComponent(page, target.ID, "", false); err != nil {
					return err
				}
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
