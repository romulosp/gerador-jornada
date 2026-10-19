package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/screen"
)

// fieldFields monta somente propriedades da família correspondente, sem persistência ou exportação.
func (e *Editor) fieldFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// content conserva o rascunho textual até a confirmação explícita.
	content := entry("")

	content.SetText(component.Properties.Field.InitialValue)
	form.Append("Valor inicial", content)
	// format preserva os nomes de intenção sem validar CPF como regra fiscal.
	format := newSelect([]string{"text", "number", "cellular", "cep", "cpf", "date", "ddd", "nis"}, nil)
	format.SetSelected(component.Properties.Field.Format)
	form.Append("Formato", format)
	return func(target *screen.Component) error {
		target.Properties.Field.InitialValue = content.Text
		target.Properties.Field.Format = format.Selected
		return nil
	}
}
