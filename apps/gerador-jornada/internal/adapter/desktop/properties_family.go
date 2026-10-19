package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// propertyFields delega cada família ao formulário coeso correspondente.
func (e *Editor) propertyFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	switch component.Kind {
	case catalog.Button:
		return e.buttonFields(form, component)
	case catalog.TextField:
		return e.fieldFields(form, component)
	case catalog.Radio, catalog.Report:
		return e.optionFields(form, component)
	case catalog.NumericKeyboard:
		return e.keyboardFields(form, component)
	case catalog.Label, catalog.Title:
		// content permite compartilhar o formulário textual sem unir variantes do modelo.
		content := entry("")
		if component.Properties.Label != nil {
			content.SetText(component.Properties.Label.Text)
		} else {
			content.SetText(component.Properties.Title.Text)
		}
		form.Append("Texto", content)
		return func(target *screen.Component) error {
			if target.Properties.Label != nil {
				target.Properties.Label.Text = content.Text
			} else {
				target.Properties.Title.Text = content.Text
			}
			return nil
		}
	}
	return func(*screen.Component) error { return nil }
}
