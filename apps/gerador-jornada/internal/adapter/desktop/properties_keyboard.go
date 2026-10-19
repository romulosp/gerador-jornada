package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// keyboardFields monta somente propriedades da família correspondente, sem persistência ou exportação.
func (e *Editor) keyboardFields(form *widget.Form, component screen.Component) func(*screen.Component) error {

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
	target := newSelect(labels, nil)
	target.SetSelected(selected)
	form.Append("Campo de destino", target)
	return func(control *screen.Component) error {
		control.Properties.Keyboard.TargetComponentID = targets[target.Selected]
		return nil
	}
}
