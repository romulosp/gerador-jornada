package catalog

import (
	"strings"
	"unicode"
)

// ValidJourneyFilename aceita um único nome portátil, sem diretórios, volumes ou dispositivos Windows.
// A existência e a legibilidade são responsabilidades do adaptador de workspace.
func ValidJourneyFilename(name string) bool {
	if strings.TrimSpace(name) != name || !strings.HasSuffix(strings.ToLower(name), ".jornada") || len(name) <= len(".jornada") || strings.HasSuffix(name, ".") || strings.ContainsAny(name, `/\:<>"|?*`) {
		return false
	}
	// char recusa controles e nomes que poderiam escapar da pasta escolhida.
	for _, char := range name {
		if unicode.IsControl(char) {
			return false
		}
	}
	// stem identifica nomes de dispositivos, inclusive quando seguidos por extensão.
	stem := strings.ToUpper(strings.SplitN(name, ".", 2)[0])
	if stem == "" || strings.HasSuffix(stem, " ") || stem == "CON" || stem == "PRN" || stem == "AUX" || stem == "NUL" {
		return false
	}
	if len(stem) == 4 && (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && stem[3] >= '0' && stem[3] <= '9' {
		return false
	}
	return true
}
