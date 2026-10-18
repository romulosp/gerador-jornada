package desktop

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/scenario"
)

// refreshScenarios apresenta roteiros e etapas como ações explícitas editáveis.
func (e *Editor) refreshScenarios() {
	e.ScenarioPanel.RemoveAll()
	e.ScenarioPanel.Add(section("DEMONSTRAÇÃO GUIADA"))
	e.ScenarioPanel.Add(widget.NewLabel("Monte etapas de cursor, clique, texto, espera e navegação. Clique é visual; navegar é uma etapa separada."))
	e.ScenarioPanel.Add(widget.NewButtonWithIcon("Novo roteiro", theme.ContentAddIcon(), func() {
		e.check(e.Session.Change(func(document *application.Document) error {
			document.Scenarios = append(document.Scenarios, scenario.Scenario{ID: project.NewID(), Name: fmt.Sprintf("Demonstração %d", len(document.Scenarios)+1), StartScreenID: document.Journey.StartScreenID, Steps: []scenario.Step{}})
			return nil
		}))
		e.changed()
	}))
	// script corresponde a um roteiro independente.
	for _, script := range e.Session.Document.Scenarios {
		e.ScenarioPanel.Add(e.scenarioCard(script))
	}
}

// scenarioCard mantém configuração e passos de um roteiro em um único cartão.
func (e *Editor) scenarioCard(script scenario.Scenario) fyne.CanvasObject {
	// items contém os passos em ordem de execução.
	items := container.NewVBox()
	// name permite renomear o roteiro sem perder seus IDs.
	name := entry(script.Name)
	items.Add(container.NewBorder(nil, nil, nil, widget.NewButton("Renomear", func() { e.editScenario(script.ID, func(target *scenario.Scenario) { target.Name = name.Text }) }), name))
	// initial oferece somente telas conhecidas como entrada.
	// initialLabels reúne as legendas das telas que podem iniciar o roteiro.
	initialLabels := []string{}
	// initialIDs associa as legendas de telas iniciais aos identificadores persistidos.
	initialIDs := map[string]string{}
	// current mantém a seleção corrente ou a tela resolvida para editar o roteiro.
	current := ""
	// page é uma possível primeira visita.
	// index localiza o item atual ou sua posição de destino na coleção ordenada.
	for index, page := range e.Session.Document.Screens {
		// label diferencia nomes iguais com um índice visível.
		label := fmt.Sprintf("%d · %s", index+1, page.Name)
		initialLabels = append(initialLabels, label)
		initialIDs[label] = page.ID
		if page.ID == script.StartScreenID {
			current = label
		}
	}
	// initial é preparado antes do callback para evitar um comando ao renderizar.
	initial := widget.NewSelect(initialLabels, nil)
	initial.SetSelected(current)
	initial.OnChanged = func(value string) {
		e.editScenario(script.ID, func(target *scenario.Scenario) { target.StartScreenID = initialIDs[value] })
	}
	items.Add(widget.NewForm(widget.NewFormItem("Tela inicial", initial)))
	// index identifica o passo que pode ser editado/reordenado.
	// step representa a etapa cuja referência será preservada ou verificada.
	for index, step := range script.Steps {
		// position permanece válido até a próxima atualização do painel.
		position := index
		// label informa ação, tempo e tela sem expor detalhes internos.
		label := fmt.Sprintf("%02d  ·  %s  ·  %d ms", index+1, step.Kind, step.DurationMS)
		items.Add(container.NewBorder(nil, nil, widget.NewLabel(label), container.NewHBox(widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() { e.stepDialog(script.ID, position) }), widget.NewButton("↑", func() { e.reorderStep(script.ID, position, -1) }), widget.NewButton("↓", func() { e.reorderStep(script.ID, position, 1) }), widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			e.editScenario(script.ID, func(target *scenario.Scenario) {
				target.Steps = append(target.Steps[:position], target.Steps[position+1:]...)
			})
		}))))
	}
	items.Add(container.NewHBox(widget.NewButtonWithIcon("Adicionar etapa", theme.ContentAddIcon(), func() { e.stepDialog(script.ID, -1) }), widget.NewButtonWithIcon("Reproduzir", theme.MediaPlayIcon(), func() { e.playScenario(script) }), widget.NewButtonWithIcon("Duplicar", theme.ContentCopyIcon(), func() {
		e.check(e.Session.Change(func(document *application.Document) error {
			// duplicate vem de uma cópia já isolada da transação.
			duplicate := *findScenario(document, script.ID)
			duplicate.ID = project.NewID()
			duplicate.Name += " — cópia"
			duplicate.Steps = append([]scenario.Step{}, duplicate.Steps...)
			// index localiza o item atual ou sua posição de destino na coleção ordenada.
			for index := range duplicate.Steps {
				duplicate.Steps[index].ID = project.NewID()
			}
			document.Scenarios = append(document.Scenarios, duplicate)
			return nil
		}))
		e.changed()
	}), widget.NewButtonWithIcon("Remover", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Remover roteiro", "Remover este roteiro e suas etapas?", func(yes bool) {
			if yes {
				e.check(e.Session.Change(func(document *application.Document) error {
					// scripts reúne os roteiros preservados após a remoção do roteiro selecionado.
					scripts := []scenario.Scenario{}
					// item representa um roteiro examinado na filtragem.
					for _, item := range document.Scenarios {
						if item.ID != script.ID {
							scripts = append(scripts, item)
						}
					}
					document.Scenarios = scripts
					return nil
				}))
				e.changed()
			}
		}, e.Window)
	})))
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := script.Interpret(e.Session.Document.Screens, e.Session.Document.Journey); err != nil {
		items.Add(widget.NewLabel("Ajuste necessário: " + err.Error()))
	}
	return widget.NewCard(script.Name, "", items)
}

// findScenario resolve um roteiro na transação, sem ponteiros da sessão anterior.
func findScenario(document *application.Document, id string) *scenario.Scenario {
	// index localiza a entidade a modificar.
	for index := range document.Scenarios {
		if document.Scenarios[index].ID == id {
			return &document.Scenarios[index]
		}
	}
	return nil
}

// editScenario publica uma alteração explícita em um único comando.
func (e *Editor) editScenario(id string, operation func(*scenario.Scenario)) {
	if e.check(e.Session.Change(func(document *application.Document) error {
		// target resolve o destino da operação na coleção ou no mapa corrente.
		if target := findScenario(document, id); target != nil {
			operation(target)
		}
		return nil
	})) {
		e.changed()
	}
}

// reorderStep muda a posição preservando a identidade e permitindo diagnosticar tela incompatível.
func (e *Editor) reorderStep(id string, index, delta int) {
	e.editScenario(id, func(script *scenario.Scenario) {
		// target resolve o destino da operação na coleção ou no mapa corrente.
		if target := index + delta; target >= 0 && target < len(script.Steps) {
			script.Steps[index], script.Steps[target] = script.Steps[target], script.Steps[index]
		}
	})
}

// stepDialog oferece propriedades compatíveis com a ação selecionada.
func (e *Editor) stepDialog(scriptID string, index int) {
	// script é consultado apenas para preparar o formulário.
	script := findScenario(&e.Session.Document, scriptID)
	if script == nil {
		return
	}
	// step é um novo passo até que exista uma posição editada.
	step := scenario.Step{ID: project.NewID(), Kind: "wait", ScreenID: script.StartScreenID, DurationMS: 500}
	// limit determina quantas navegações precedem o novo passo.
	limit := index
	if limit < 0 {
		limit = len(script.Steps)
	}
	// previous fornece uma etapa anterior para resolver o contexto da etapa editada.
	for _, previous := range script.Steps[:limit] {
		if previous.Kind == "navigate" {
			step.ScreenID = e.Session.Document.Journey.Destination(previous.ComponentID)
		}
	}
	if index >= 0 {
		step = script.Steps[index]
	}
	// kind seleciona a semântica, não uma função serializada.
	kind := widget.NewSelect([]string{"moveCursor", "click", "typeText", "wait", "navigate"}, nil)
	kind.SetSelected(step.Kind)
	// pages e pageIDs representam as composições possíveis.
	pages := []string{}
	// pageIDs associa legendas de telas a suas identidades persistidas.
	pageIDs := map[string]string{}
	// selectedPage mantém a tela selecionada no formulário da etapa.
	selectedPage := ""
	// number contém o valor inteiro convertido de uma propriedade geométrica.
	// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
	for number, page := range e.Session.Document.Screens {
		// label diferencia nomes iguais com um índice visível.
		label := fmt.Sprintf("%d · %s", number+1, page.Name)
		pages = append(pages, label)
		pageIDs[label] = page.ID
		if page.ID == step.ScreenID {
			selectedPage = label
		}
	}
	// page permite corrigir uma referência após reordenar etapas.
	page := widget.NewSelect(pages, nil)
	page.SetSelected(selectedPage)
	// targets vincula legendas de controles aos IDs corretos.
	targets := map[string]string{}
	// target indica um componente ou uma posição absoluta.
	target := widget.NewSelect([]string{}, nil)
	// refreshTargets restringe alvos conforme a ação e a tela.
	refreshTargets := func() {
		targets = map[string]string{"Posição absoluta": ""}
		// labels reúne as legendas exibidas sem expor identificadores internos.
		labels := []string{}
		if kind.Selected == "moveCursor" {
			labels = append(labels, "Posição absoluta")
		}
		// current mantém a seleção corrente ou a tela resolvida para editar o roteiro.
		if current := e.Session.Document.Screen(pageIDs[page.Selected]); current != nil {
			// component identifica o controle da tela que será consultado ou alterado.
			for _, component := range current.Components {
				if (kind.Selected == "typeText" && component.Kind != catalog.TextField) || (kind.Selected == "navigate" && component.Kind != catalog.Button) || (kind.Selected == "click" && component.Kind != catalog.Button && component.Kind != catalog.TextField) {
					continue
				}
				labels = append(labels, component.Name)
				targets[component.Name] = component.ID
			}
		}
		target.Options = labels
		target.ClearSelected()
		// label diferencia nomes iguais com um índice visível.
		// id identifica de forma estável o item ou recurso atual.
		for label, id := range targets {
			if id == step.ComponentID {
				target.SetSelected(label)
				break
			}
		}
		target.Refresh()
	}
	kind.OnChanged = func(string) { refreshTargets() }
	page.OnChanged = func(string) { refreshTargets() }
	refreshTargets()
	// duration e text são rascunhos até confirmação.
	duration := entry(strconv.Itoa(step.DurationMS))
	// text mantém o texto editado ou seus caracteres para a operação atual.
	text := entry("")
	if step.Text != nil {
		text.SetText(*step.Text)
	}
	// x e y só são usados no movimento por coordenadas.
	x := entry("0")
	// y mantém a coordenada vertical editada no formulário.
	y := entry("0")
	if step.Point != nil {
		x.SetText(fmt.Sprint(step.Point.X))
		y.SetText(fmt.Sprint(step.Point.Y))
	}
	// form mantém todos os campos visíveis para deixar sua semântica explícita.
	form := widget.NewForm(widget.NewFormItem("Ação", kind), widget.NewFormItem("Tela corrente", page), widget.NewFormItem("Componente", target), widget.NewFormItem("Duração (ms)", duration), widget.NewFormItem("Texto (digitação)", text), widget.NewFormItem("X (posição absoluta)", x), widget.NewFormItem("Y (posição absoluta)", y))
	dialog.ShowCustomConfirm("Etapa da demonstração", "Salvar etapa", "Cancelar", container.NewVBox(form, widget.NewLabel("wait não usa alvo. click é apenas visual.\nnavigate segue a ação configurada no botão.")), func(yes bool) {
		if !yes {
			return
		}
		// milliseconds é validado antes da alteração.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		milliseconds, err := strconv.Atoi(duration.Text)
		if !e.check(err) {
			return
		}
		if milliseconds < 1 || milliseconds > 600000 {
			e.check(fmt.Errorf("duração deve estar entre 1 e 600000 ms"))
			return
		}
		// next contém apenas campos aplicáveis à ação.
		next := scenario.Step{ID: step.ID, Kind: kind.Selected, ScreenID: pageIDs[page.Selected], DurationMS: milliseconds}
		if next.Kind != "wait" {
			next.ComponentID = targets[target.Selected]
		}
		if next.Kind == "typeText" {
			// value contém o valor corrente em validação, serialização ou edição.
			value := text.Text
			next.Text = &value
		}
		if next.Kind == "moveCursor" && next.ComponentID == "" {
			// px contém a coordenada horizontal convertida para unidades lógicas.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			px, err := strconv.ParseFloat(x.Text, 64)
			if !e.check(err) {
				return
			}
			// py contém a coordenada vertical convertida para unidades lógicas.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			py, err := strconv.ParseFloat(y.Text, 64)
			if !e.check(err) {
				return
			}
			next.Point = &scenario.Point{X: px, Y: py}
		}
		e.editScenario(scriptID, func(document *scenario.Scenario) {
			if index < 0 {
				document.Steps = append(document.Steps, next)
			} else {
				document.Steps[index] = next
			}
		})
	}, e.Window)
}
