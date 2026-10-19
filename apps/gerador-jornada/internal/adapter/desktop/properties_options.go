package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
	"strings"
)

// optionFields monta somente propriedades da família correspondente, sem persistência ou exportação.
func (e *Editor) optionFields(form *widget.Form, component screen.Component) func(*screen.Component) error {

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
	orientation := newSelect([]string{"vertical", "horizontal"}, nil)
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
}
