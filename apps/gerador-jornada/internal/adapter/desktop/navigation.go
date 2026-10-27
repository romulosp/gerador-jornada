package desktop

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/screen"
)

// refreshJourney apresenta fontes, destinos e diagnósticos sem esconder vínculos quebrados.
func (e *Editor) refreshJourney() {
	e.JourneyPanel.RemoveAll()
	e.JourneyPanel.Add(section("FLUXO DA JORNADA"))
	e.JourneyPanel.Add(widget.NewLabel("Conecte os botões às telas. Os mesmos destinos serão usados na prévia e no PowerPoint."))
	// page organiza os vínculos por composição de origem.
	for _, page := range e.Session.Document.Screens {
		// items reúne os botões da tela.
		items := container.NewVBox()
		if page.ID == e.Session.Document.Journey.StartScreenID {
			items.Add(widget.NewLabelWithStyle("Tela inicial", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		}
		// component fornece um botão configurável.
		for _, component := range page.Components {
			if component.Kind == catalog.Button {
				// copy conserva a identidade utilizada no callback.
				copy := component
				items.Add(container.NewBorder(nil, nil, widget.NewLabel(component.Name), nil, e.navigationFor(page.ID, copy)))
			}
		}
		e.JourneyPanel.Add(widget.NewCard(page.Name, "", items))
	}
	// problems informa referências quebradas e erros de composição.
	problems := e.Session.Document.Validate()
	if len(problems) > 0 {
		e.JourneyPanel.Add(widget.NewCard("Ajustes necessários", "", widget.NewLabel(application.ProblemsText(problems))))
	} else {
		e.JourneyPanel.Add(widget.NewButtonWithIcon("Experimentar a jornada", theme.VisibilityIcon(), e.preview))
	}
}

// navigationFor reutiliza o inspector sem trocar a seleção da tela em edição.
func (e *Editor) navigationFor(pageID string, component screen.Component) fyne.CanvasObject {
	// options relaciona rótulos legíveis a IDs persistidos.
	options := []string{"Sem navegação"}
	// targets associa as legendas visíveis aos identificadores dos destinos.
	targets := map[string]string{"Sem navegação": ""}
	// selected mantém a legenda selecionada no formulário.
	selected := "Sem navegação"
	// page é um destino disponível.
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
	// control é configurado antes de registrar seu callback.
	control := newSelect(options, nil)
	control.SetSelected(selected)
	control.OnChanged = func(value string) {
		e.check(e.Session.Change(func(document *application.Document) error {
			// kind seleciona o discriminador da ação conforme o destino informado.
			kind := "navigate"
			if targets[value] == "" {
				kind = "none"
			}
			return application.BindNavigation(document, journey.Action{SourceScreenID: pageID, SourceComponentID: component.ID, Kind: kind, TargetScreenID: targets[value]})
		}))
		e.changed()
	}
	return control
}
