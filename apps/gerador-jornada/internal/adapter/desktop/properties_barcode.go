package desktop

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// barcodeInputFields monta somente o rótulo configurável da entrada, sem valores de runtime ou DV.
func (e *Editor) barcodeInputFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// label conserva o rascunho do rótulo até a confirmação explícita.
	label := entry("")
	if component.Properties.BarcodeInput != nil {
		label.SetText(component.Properties.BarcodeInput.Label)
	}
	form.Append("Rótulo", label)
	return func(target *screen.Component) error {
		if target.Properties.BarcodeInput == nil {
			return fmt.Errorf("componente sem propriedades de código de barras")
		}
		// draft reúne a configuração completa antes de substituir o estado confirmado.
		draft := catalog.BarcodeInputProperties{Label: label.Text}
		if err := draft.Validate(); err != nil {
			return err
		}
		*target.Properties.BarcodeInput = draft
		return target.Properties.Validate(target.Kind)
	}
}
