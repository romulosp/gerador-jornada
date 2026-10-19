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
type jornadaTheme struct {
	// name identifica a paleta da interface, sem alterar cores do documento.
	name string
}

// Color resolve superfícies e estados da paleta escolhida, sem seguir o tema do sistema.
func (t jornadaTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.name == "Escuro" {
		switch name {
		case theme.ColorNamePrimary, theme.ColorNameFocus:
			return rgb("#60A5FA")
		case theme.ColorNameBackground:
			return rgb("#111827")
		case theme.ColorNameForeground:
			return rgb("#F1F5F9")
		case theme.ColorNameForegroundOnPrimary:
			return rgb("#111827")
		case theme.ColorNameInputBackground, theme.ColorNameMenuBackground:
			return rgb("#1F2937")
		case theme.ColorNameButton:
			return rgb("#334155")
		case theme.ColorNameDisabled:
			return rgb("#94A3B8")
		case theme.ColorNamePlaceHolder:
			return rgb("#CBD5E1")
		case theme.ColorNameSelection:
			return rgb("#1E40AF")
		case theme.ColorNameSeparator:
			return rgb("#475569")
		}
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
	if t.name == "Azul" {
		switch name {
		case theme.ColorNamePrimary, theme.ColorNameFocus:
			return rgb("#174EA6")
		case theme.ColorNameBackground:
			return rgb("#EDF4FC")
		case theme.ColorNameForeground:
			return rgb("#142D4E")
		case theme.ColorNameButton:
			return rgb("#DCE9F8")
		case theme.ColorNamePlaceHolder:
			return rgb("#486581")
		case theme.ColorNameSelection:
			return rgb("#C9DDF5")
		}
	}
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
	case theme.ColorNamePlaceHolder:
		return rgb("#64748B")
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
