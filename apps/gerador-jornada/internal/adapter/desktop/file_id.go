package desktop

import (
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/screen"
)

// fileIDEntry mantém edição ASCII com validação do limite e do prefixo obrigatório das telas.
// A limpeza de acentos é imediata; excesso de tamanho e nomes inválidos impedem salvar.
func fileIDEntry(value string, panel bool) *widget.Entry {
	// field mantém o rascunho do ID, independente do nome amigável e do UUID.
	field := entry(value)
	field.SetPlaceHolder("Até 40 caracteres, sem acentos")
	field.Validator = func(text string) error { return screen.ValidateFileID(text, panel) }
	field.OnChanged = func(text string) {
		// clean impede que acentos, espaços ou símbolos permaneçam no atributo de arquivo.
		clean := screen.CleanFileID(text)
		if clean != text {
			field.SetText(clean)
		}
	}
	return field
}
