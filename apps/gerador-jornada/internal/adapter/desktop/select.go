package desktop

import "fyne.io/fyne/v2/widget"

// newSelect cria seletores com placeholder em português sem depender da tradução global do toolkit.
func newSelect(options []string, changed func(string)) *widget.Select {
	// selectBox mantém o mesmo comportamento nativo com texto de interface explicitamente definido.
	selectBox := &widget.Select{Options: options, OnChanged: changed, PlaceHolder: "Selecione"}
	selectBox.ExtendBaseWidget(selectBox)
	return selectBox
}
