package desktop

import "fyne.io/fyne/v2/theme"

// themePreferenceKey separa a preferência local do formato persistido da jornada.
const themePreferenceKey = "editor.theme"

// normalizedTheme recusa valores desconhecidos usando o tema claro como padrão seguro.
func normalizedTheme(name string) string {
	switch name {
	case "Escuro", "Azul":
		return name
	default:
		return "Claro"
	}
}

// initTheme recupera a preferência por usuário e prepara a seleção sem reconstruir o formulário.
func (e *Editor) initTheme() {
	e.ThemeSelect = newSelect([]string{"Claro", "Escuro", "Azul"}, e.applyTheme)
	e.ThemeSelect.SetSelected(normalizedTheme(e.App.Preferences().String(themePreferenceKey)))
}

// applyTheme atualiza apenas a interface e a preferência local, preservando todo estado de edição.
func (e *Editor) applyTheme(name string) {
	name = normalizedTheme(name)
	// selected mantém a paleta independente da configuração visual do sistema operacional.
	selected := jornadaTheme{name: name}
	e.App.Settings().SetTheme(selected)
	e.App.Preferences().SetString(themePreferenceKey, name)
	if e.Title != nil {
		e.Title.Color = selected.Color(theme.ColorNameForeground, theme.VariantLight)
		e.Title.Refresh()
	}
	if e.Brand != nil {
		e.Brand.Color = selected.Color(theme.ColorNamePrimary, theme.VariantLight)
		e.Brand.Refresh()
	}
	if e.Subtitle != nil {
		e.Subtitle.Color = selected.Color(theme.ColorNamePlaceHolder, theme.VariantLight)
		e.Subtitle.Refresh()
	}
}
