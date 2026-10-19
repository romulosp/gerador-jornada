package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// setTextAlignment configura o parágrafo dentro da caixa completa sem alterar sua posição.
func setTextAlignment(text *canvas.Text, alignment string) {
	text.Alignment = fyne.TextAlignLeading
	if alignment == "center" {
		text.Alignment = fyne.TextAlignCenter
	} else if alignment == "right" {
		text.Alignment = fyne.TextAlignTrailing
	}
}
