package desktop

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/atomicfile"
	"gerador-jornada.local/editor/internal/startup"
)

// startupPreferenceFile usa a mesma configuração por usuário dos recentes, fora do projeto e EXE.
func startupPreferenceFile() (string, error) {
	// recents fornece o diretório comum, incluindo isolamento por JORNADA_CONFIG_DIR.
	recents, err := recentFile()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(recents), "startup-preferences.json"), nil
}

// loadStartupPreference conserva arquivo inválido para diagnóstico, sem sobrescrevê-lo.
func loadStartupPreference() (startup.Preference, error) {
	// preference assume a primeira execução quando não existe uma configuração anterior.
	var preference startup.Preference
	path, err := startupPreferenceFile()
	if err != nil {
		return preference, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return preference, nil
	}
	if err != nil {
		return preference, err
	}
	err = json.Unmarshal(data, &preference)
	return preference, err
}

// updateStartupPreference relê antes de alterar um campo, preservando a outra escolha persistida.
func updateStartupPreference(change func(*startup.Preference)) error {
	// preference contém o último estado publicado, inclusive alterações feitas pelo launcher.
	preference, err := loadStartupPreference()
	if err != nil {
		return err
	}
	change(&preference)
	data, err := json.MarshalIndent(preference, "", "  ")
	if err != nil {
		return err
	}
	path, err := startupPreferenceFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return atomicfile.Write(context.Background(), path, data)
}

// splashPreferenceCheck monta a opção sem gravar durante a construção do launcher.
func (l *Launcher) splashPreferenceCheck() *widget.Check {
	// preference e check representam a escolha salva, independente do primeiro uso obrigatório.
	preference, err := loadStartupPreference()
	check := widget.NewCheck("Mostrar splash screen ao iniciar", nil)
	check.SetChecked(preference.ShowOnStartup)
	if err != nil {
		check.Disable()
		l.Message.SetText("Não foi possível ler a preferência da splash: " + err.Error())
		return check
	}
	check.OnChanged = func(show bool) {
		if err := updateStartupPreference(func(current *startup.Preference) { current.ShowOnStartup = show }); err != nil {
			// Reverter sem disparar outra gravação mantém a interface coerente com o arquivo anterior.
			handler := check.OnChanged
			check.OnChanged = nil
			check.SetChecked(preference.ShowOnStartup)
			check.OnChanged = handler
			l.Message.SetText("Não foi possível salvar a preferência da splash: " + err.Error())
			return
		}
		preference.ShowOnStartup = show
		LogApplication("splash_preference_changed", map[string]any{"show_on_startup": show})
	}
	return check
}
