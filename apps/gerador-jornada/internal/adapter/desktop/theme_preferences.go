package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

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
	// SetSelected dispara OnChanged; a janela ainda não tem viewport nativo neste ponto.
	// O tema global já foi preparado antes de NewWindow e não pode ser reaplicado aqui.
	e.ThemeSelect = newSelect([]string{"Claro", "Escuro", "Azul"}, nil)
	e.ThemeSelect.SetSelected(normalizedTheme(e.App.Preferences().String(themePreferenceKey)))
	e.ThemeSelect.OnChanged = e.applyTheme
}

// prepareApplicationTheme configura o tema antes de criar uma janela e evita reaplicações redundantes.
// No driver Windows, trocar tema com uma janela recém-criada e ainda não exibida causa acesso nativo nulo.
func prepareApplicationTheme(app fyne.App, name string) {
	// selected normaliza a preferência antes de comparar a paleta já instalada.
	selected := jornadaTheme{name: normalizedTheme(name)}
	if current, ok := app.Settings().Theme().(jornadaTheme); ok && current == selected {
		return
	}
	app.Settings().SetTheme(selected)
}

// applyTheme atualiza apenas a interface e a preferência local, preservando todo estado de edição.
func (e *Editor) applyTheme(name string) {
	name = normalizedTheme(name)
	// selected mantém a paleta independente da configuração visual do sistema operacional.
	selected := jornadaTheme{name: name}
	prepareApplicationTheme(e.App, name)
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
