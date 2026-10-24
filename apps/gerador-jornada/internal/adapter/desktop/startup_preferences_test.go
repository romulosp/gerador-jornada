package desktop

import (
	"context"
	"errors"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/startup"
)

// TestLauncherSplashPreferencePersistence percorre a opção real e a reabertura por usuário.
func TestLauncherSplashPreferencePersistence(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app mantém as janelas de teste sem iniciar reprodução nativa.
	app := test.NewApp()
	defer app.Quit()
	launcher := NewLauncher(app, nil)
	defer launcher.Window.Close()
	path, err := startupPreferenceFile()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) || launcher.ShowSplash.Checked {
		t.Fatal("construir launcher gravou preferência ou marcou a opção padrão")
	}
	if err := updateStartupPreference(func(saved *startup.Preference) { saved.Presented = true }); err != nil {
		t.Fatal(err)
	}
	// A alteração do launcher não pode apagar a conclusão publicada depois de sua construção.
	test.Tap(launcher.ShowSplash)
	saved, err := loadStartupPreference()
	if err != nil || !saved.Presented || !saved.ShowOnStartup || !saved.ShouldShow() {
		t.Fatalf("opção não persistiu ou apagou primeiro uso: %+v / %v", saved, err)
	}
	reopened := NewLauncher(app, nil)
	defer reopened.Window.Close()
	if !reopened.ShowSplash.Checked {
		t.Fatal("reabertura não restaurou opção marcada")
	}
	test.Tap(reopened.ShowSplash)
	saved, err = loadStartupPreference()
	if err != nil || !saved.Presented || saved.ShouldShow() {
		t.Fatalf("desmarcar não suprimiu próximas apresentações: %+v / %v", saved, err)
	}
	// A publicação da conclusão também deve preservar a opção caso ela já esteja marcada.
	test.Tap(reopened.ShowSplash)
	if err := updateStartupPreference(func(current *startup.Preference) { current.Presented = true }); err != nil {
		t.Fatal(err)
	}
	saved, _ = loadStartupPreference()
	if !saved.ShowOnStartup {
		t.Fatal("conclusão sobrescreveu opção marcada")
	}
}

// TestSplashPreferenceWriteFailure conserva a configuração e reverte a marcação não salva.
func TestSplashPreferenceWriteFailure(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// launcher nasce com configuração válida; a corrupção posterior simula uma falha real de persistência.
	app := test.NewApp()
	defer app.Quit()
	launcher := NewLauncher(app, nil)
	defer launcher.Window.Close()
	path, _ := startupPreferenceFile()
	const corrupted = "{configuração inválida"
	if err := os.WriteFile(path, []byte(corrupted), 0600); err != nil {
		t.Fatal(err)
	}
	test.Tap(launcher.ShowSplash)
	data, err := os.ReadFile(path)
	if err != nil || string(data) != corrupted || launcher.ShowSplash.Checked || launcher.Message.Text == "" {
		t.Fatal("falha de gravação perdeu configuração, confirmou seleção ou ocultou diagnóstico")
	}
	reopened := NewLauncher(app, nil)
	defer reopened.Window.Close()
	if !reopened.ShowSplash.Disabled() || reopened.Message.Text == "" {
		t.Fatal("configuração inválida foi oferecida para sobrescrita silenciosa")
	}
}

// TestSkippedSplashStillLoads verifica que pular vídeo não pula recursos ou erros obrigatórios.
func TestSkippedSplashStillLoads(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	if err := updateStartupPreference(func(saved *startup.Preference) { saved.Presented = true }); err != nil {
		t.Fatal(err)
	}
	// app e etapas observam somente a ordem externa; nenhum evento fictício de vídeo é necessário.
	app := test.NewApp()
	defer app.Quit()
	loaded, opened, ready := false, false, false
	window := StartDesktop(app, func(current fyne.Window) fyne.Window {
		if !loaded {
			t.Fatal("janela abriu antes dos recursos")
		}
		opened = true
		launcher := NewLauncherInWindow(app, current, nil)
		return launcher.Window
	}, func(context.Context, func(float64, string)) error {
		loaded = true
		return nil
	}, func() { ready = true })
	defer window.Close()
	if !loaded || !opened || !ready || window.FullScreen() {
		t.Fatal("sem splash deve carregar e revelar janela normal")
	}
	opened, ready = false, false
	failedWindow := StartDesktop(app, func(current fyne.Window) fyne.Window {
		opened = true
		return current
	}, func(context.Context, func(float64, string)) error {
		return errors.New("recurso obrigatório indisponível")
	}, func() { ready = true })
	defer failedWindow.Close()
	if opened || ready {
		t.Fatal("falha de recursos liberou tela seguinte ao pular splash")
	}
}
