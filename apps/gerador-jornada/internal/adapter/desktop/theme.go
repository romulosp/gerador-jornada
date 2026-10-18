// Package desktop implementa a interface Fyne sem colocar regras de negócio nos widgets.
package desktop

import (
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// jornadaTheme adapta a identidade visual sem alterar o tema global do Windows.
type jornadaTheme struct{}

// Color define superfícies claras e contraste da navegação.
func (jornadaTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameFocus:
		return rgb("#2563EB")
	case theme.ColorNameBackground:
		return rgb("#F8FAFC")
	case theme.ColorNameForeground:
		return rgb("#0F172A")
	case theme.ColorNameInputBackground, theme.ColorNameMenuBackground:
		return rgb("#FFFFFF")
	case theme.ColorNameButton:
		return rgb("#E8EEF6")
	case theme.ColorNameDisabled:
		return rgb("#94A3B8")
	case theme.ColorNameSelection:
		return rgb("#DBEAFE")
	case theme.ColorNameSeparator:
		return rgb("#E2E8F0")
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

// Font preserva fontes empacotadas pelo Fyne para distribuição offline.
func (jornadaTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }

// Icon reutiliza ícones vetoriais consistentes com o desktop.
func (jornadaTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size ajusta espaçamento e raio sem afetar a geometria persistida das telas.
func (jornadaTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNamePadding {
		return 10
	}
	if name == theme.SizeNameText {
		return 14
	}
	return theme.DefaultTheme().Size(name)
}

// rgb converte uma cor validada de projeto para o renderer.
func rgb(value string) color.Color {
	// number representa os três canais RGB.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	number, err := strconv.ParseUint(strings.TrimPrefix(value, "#"), 16, 32)
	if err != nil {
		return color.NRGBA{R: 100, G: 116, B: 139, A: 255}
	}
	return color.NRGBA{R: uint8(number >> 16), G: uint8(number >> 8), B: uint8(number), A: 255}
}
