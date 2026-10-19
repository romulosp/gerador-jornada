package desktop

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/screen"
)

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
	destination := newSelect(options, nil)
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
