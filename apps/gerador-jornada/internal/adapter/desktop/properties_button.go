package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/screen"
)

// buttonFields monta somente propriedades da família correspondente, sem persistência ou exportação.
func (e *Editor) buttonFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// content conserva o rascunho textual até a confirmação explícita.
	content := entry("")

	content.SetText(component.Properties.Button.Text)
	form.Append("Texto", content)
	// variant restringe a escolha às skins previstas.
	variant := newSelect([]string{"confirm", "cancel", "yes", "no"}, nil)
	variant.SetSelected(component.Properties.Button.Variant)
	form.Append("Tipo", variant)
	return func(target *screen.Component) error {
		target.Properties.Button.Text = content.Text
		target.Properties.Button.Variant = variant.Selected
		return nil
	}
}
