package screen

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// FileIDLimit limita o nome-base do arquivo; a extensão não faz parte dos 40 caracteres.
const FileIDLimit = 40

// CleanFileID remove acentos e caracteres não alfanuméricos sem alterar letras ASCII ou números.
func CleanFileID(value string) string {
	// result mantém somente caracteres válidos para os nomes-base adotados pelo gerador.
	var result strings.Builder
	// character percorre a forma decomposta para descartar marcas de acentuação.
	for _, character := range norm.NFD.String(value) {
		if character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z' || character >= '0' && character <= '9' {
			result.WriteRune(character)
		}
	}
	return result.String()
}

// SuggestFileID sugere um nome válido; telas usam Pnl e ignoram números de ordenação no início.
// A sugestão inicial pode ser encurtada; alterações explícitas são validadas sem truncamento silencioso.
func SuggestFileID(name string, panel bool) string {
	// words permite separar palavras antes de remover pontuação e acentuação.
	words := strings.FieldsFunc(norm.NFD.String(name), func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsDigit(character) && !unicode.Is(unicode.Mn, character)
	})
	// value conserva cada palavra com inicial maiúscula, sem caracteres especiais.
	value := ""
	for _, word := range words {
		// clean remove as marcas Unicode antes de aplicar a capitalização ASCII.
		clean := CleanFileID(word)
		if clean != "" {
			value += strings.ToUpper(clean[:1]) + clean[1:]
		}
	}
	if panel {
		value = strings.TrimLeft(value, "0123456789")
		if strings.HasPrefix(strings.ToLower(value), "pnl") {
			value = value[3:]
		}
		if value == "" {
			value = "Tela"
		}
		value = "Pnl" + value
	} else {
		if value == "" {
			value = "Componente"
		}
		if value[0] >= '0' && value[0] <= '9' {
			value = "Componente" + value
		}
		if reservedFileID(value) {
			value += "Item"
		}
	}
	return value[:min(len(value), FileIDLimit)]
}

// ValidateFileID rejeita nomes vazios, especiais, reservados, longos ou sem o prefixo de tela.
func ValidateFileID(value string, panel bool) error {
	if len(value) == 0 || len(value) > FileIDLimit {
		return fmt.Errorf("ID deve conter de 1 a 40 caracteres, somente letras sem acento e números")
	}
	// character verifica ASCII diretamente, sem alocar normalizações para IDs já válidos.
	for _, character := range value {
		if !((character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9')) {
			return fmt.Errorf("ID permite somente letras sem acento e números")
		}
	}
	if !((value[0] >= 'A' && value[0] <= 'Z') || (value[0] >= 'a' && value[0] <= 'z')) {
		return fmt.Errorf("ID deve começar com uma letra")
	}
	if panel && (!strings.HasPrefix(value, "Pnl") || len(value) == 3) {
		return fmt.Errorf("ID da tela deve começar com Pnl e conter um nome, como PnlIdentificacao")
	}
	if reservedFileID(value) {
		return fmt.Errorf("ID reservado pelo Windows; escolha outro nome")
	}
	return nil
}

// reservedFileID identifica dispositivos que não podem ser usados como nome de arquivo Windows.
func reservedFileID(value string) bool {
	// upper permite tratar o sistema de arquivos sem diferenciação de caixa.
	upper := strings.ToUpper(value)
	return upper == "CON" || upper == "PRN" || upper == "AUX" || upper == "NUL" ||
		len(upper) == 4 && (strings.HasPrefix(upper, "COM") || strings.HasPrefix(upper, "LPT")) && upper[3] >= '1' && upper[3] <= '9'
}
