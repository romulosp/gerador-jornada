package desktop

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"gerador-jornada.local/editor/internal/startup"
)

// StartDesktop aplica a política por usuário antes de criar player ou janela em tela cheia.
// Sem apresentação, prepara recursos antes de exibir a janela normal; a carga continua obrigatória.
// next monta o destino na thread de UI; finished indica que ele pode receber interação/captura.
func StartDesktop(app fyne.App, next func(fyne.Window) fyne.Window, load func(context.Context, func(float64, string)) error, finished func()) fyne.Window {
	// preference decide apenas a apresentação, sem dispensar as etapas de preparação.
	preference, err := loadStartupPreference()
	if err != nil {
		LogApplication("startup_preference_failed", nil)
	}
	if load == nil {
		load = LoadStartupResources
	}
	if preference.ShouldShow() {
		// splash confirma a primeira apresentação somente após vídeo e recursos completos.
		splash := NewSplash(app, func(current fyne.Window) fyne.Window {
			// markError não bloqueia autoria; a próxima abertura repetirá se a gravação falhar.
			markError := updateStartupPreference(func(saved *startup.Preference) { saved.Presented = true })
			destination := next(current)
			if markError != nil && destination != nil {
				LogApplication("startup_preference_failed", nil)
				dialog.ShowError(fmt.Errorf("Não foi possível registrar a apresentação da splash: %w", markError), destination)
			}
			return destination
		}, finished)
		splash.Show(load)
		return splash.Window
	}
	LogApplication("splash_skipped", map[string]any{"show_on_startup": false})
	// O tema deve preceder NewWindow para preservar o viewport nativo do launcher/editor.
	prepareApplicationTheme(app, app.Preferences().String(themePreferenceKey))
	window := app.NewWindow("Jornada Studio · Inicializando")
	if err := load(context.Background(), func(float64, string) {}); err != nil {
		LogApplication("startup_failed", map[string]any{"phase": "resources"})
		window.Show()
		dialog.ShowError(fmt.Errorf("Não foi possível inicializar: %w", err), window)
		return window
	}
	LogApplication("startup_resources_ready", nil)
	window = next(window)
	if window != nil {
		window.Show()
		LogApplication("startup_ready", nil)
		if finished != nil {
			finished()
		}
	}
	return window
}
