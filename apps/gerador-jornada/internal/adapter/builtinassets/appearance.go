package builtinassets

import (
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// ValidateAppearanceEdit recusa novas escolhas indisponíveis sem substituir aparências históricas.
// Manual e MCP usam os mesmos recursos incorporados; nenhum arquivo externo é consultado.
func ValidateAppearanceEdit(before, after screen.Component) error {
	// library já validou bytes/hash e famílias físicas no carregamento.
	library, err := skin.Builtin()
	if err != nil {
		return &application.AuthoringError{Code: "capabilityUnavailable"}
	}
	if after.Style.Skin != before.Style.Skin && after.Style.Skin != "" {
		if _, found := library.Lookup(after.Kind, after.Style.Skin); !found {
			return &application.AuthoringError{Code: "capabilityUnavailable"}
		}
	}
	if after.Style.FontFamily != before.Style.FontFamily && after.Style.FontFamily != "" {
		// found exige a face física incorporada, sem fonte parecida ou fallback silencioso.
		found := false
		for _, path := range library.ResourcePaths() {
			found = found || library.FontFamily(path) == after.Style.FontFamily
		}
		if !found {
			return &application.AuthoringError{Code: "capabilityUnavailable"}
		}
	}
	return nil
}
