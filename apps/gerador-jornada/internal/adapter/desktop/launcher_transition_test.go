package desktop

import (
	"path/filepath"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"gerador-jornada.local/editor/internal/application"
)

// themeTrackingSettings observa em que momento o tema global é publicado no driver.
type themeTrackingSettings struct {
	fyne.Settings
	// BeforeSet permite detectar uma troca de tema depois de registrar janela ainda não exibida.
	BeforeSet func()
}

// SetTheme registra a publicação e delega ao adaptador real de preferências de teste.
func (s *themeTrackingSettings) SetTheme(theme fyne.Theme) {
	s.BeforeSet()
	s.Settings.SetTheme(theme)
}

// themeTrackingApp conserva o driver Fyne e substitui apenas o observador de Settings.
type themeTrackingApp struct {
	fyne.App
	// SettingsObserver preserva a mesma instância para todas as janelas do teste.
	SettingsObserver *themeTrackingSettings
}

// Settings devolve a preferência observada sem trocar o driver nem o aplicativo corrente.
func (a *themeTrackingApp) Settings() fyne.Settings { return a.SettingsObserver }

// TestLauncherEditorThemeOrder impede reaplicar o tema com o editor nativo ainda sem viewport.
func TestLauncherEditorThemeOrder(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// base e app conservam o driver enquanto o observador registra a quantidade de janelas existentes.
	base := test.NewApp()
	defer base.Quit()
	base.Preferences().SetString(themePreferenceKey, "Escuro")
	// O driver de teste já pode ter uma janela auxiliar; somente janelas novas importam para a ordem.
	initialWindows := len(base.Driver().AllWindows())
	windowsAtTheme := []int{}
	app := &themeTrackingApp{App: base}
	app.SettingsObserver = &themeTrackingSettings{Settings: base.Settings(), BeforeSet: func() {
		windowsAtTheme = append(windowsAtTheme, len(base.Driver().AllWindows()))
	}}
	var editor *Editor
	var launcher *Launcher
	launcher = NewLauncher(app, func(document application.Document, path string) {
		editor = New(app, document)
		editor.SetProjectPath(path)
		editor.ShowWindow()
		launcher.Window.Close()
	})
	launcher.Window.Show()
	if err := launcher.createProject("teste", t.TempDir(), false); err != nil {
		t.Fatal(err)
	}
	defer func() { editor.close() }()
	if !reflect.DeepEqual(windowsAtTheme, []int{initialWindows}) {
		t.Fatalf("tema reaplicado depois de criar janela: %v", windowsAtTheme)
	}
	if editor.ThemeSelect.Selected != "Escuro" || editor.Session.Document.Project.Name != "teste" {
		t.Fatal("transição perdeu tema ou projeto")
	}
	if !reflect.DeepEqual(editor.Title.Color, app.Settings().Theme().Color(theme.ColorNameForeground, theme.VariantLight)) {
		t.Fatal("título não recebeu a cor do tema antes de exibir a janela")
	}
	// recent repete a abertura do arquivo já criado com outra janela de boas-vindas.
	path := editor.FilePath
	editor.close()
	launcher = NewLauncher(app, func(document application.Document, path string) {
		editor = New(app, document)
		editor.SetProjectPath(path)
		editor.ShowWindow()
		launcher.Window.Close()
	})
	launcher.Window.Show()
	launcher.openRecent(RecentProject{Path: path})
	if editor.FilePath != filepath.Clean(path) || !reflect.DeepEqual(windowsAtTheme, []int{initialWindows}) {
		t.Fatal("reabertura reaplicou tema durante construção da janela")
	}
	// A seleção feita pelo usuário continua funcionando depois de exibir a janela.
	editor.ThemeSelect.SetSelected("Azul")
	if base.Preferences().String(themePreferenceKey) != "Azul" || len(windowsAtTheme) != 2 {
		t.Fatal("mudança explícita de tema deixou de funcionar")
	}
}
