package desktop

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"
	"gerador-jornada.local/editor/internal/application"
)

// TestRecentsAndLauncher percorre persistência, entrada removida e os dois modos sem documento implícito.
func TestRecentsAndLauncher(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app e opened permitem detectar qualquer carregamento automático indevido.
	app := test.NewApp()
	defer app.Quit()
	opened := []application.Document{}
	launcher := NewLauncher(app, func(document application.Document, path string) { opened = append(opened, document) })
	defer launcher.Window.Close()
	if len(opened) != 0 {
		t.Fatal("launcher abriu documento automaticamente")
	}
	base := t.TempDir()
	if err := launcher.createProject("Projeto em branco — ação", base, false); err != nil {
		t.Fatal(err)
	}
	if err := launcher.createProject("Projeto com exemplo", base, true); err != nil {
		t.Fatal(err)
	}
	if len(opened) != 2 || len(opened[0].Screens) != 1 || len(opened[0].Screens[0].Components) != 0 || len(opened[0].Project.Assets) == 0 || len(opened[1].Screens[0].Components) == 0 || opened[1].Project.Name != "Projeto com exemplo" {
		t.Fatal("modos em branco/exemplo divergentes")
	}
	entries, err := loadRecents()
	if err != nil || len(entries) != 2 || entries[0].Name != "Projeto com exemplo" || entries[0].OpenedAt.IsZero() {
		t.Fatalf("recentes não persistidos: %+v / %v", entries, err)
	}
	if err := rememberProject(entries[0].Path, entries[0].Name); err != nil {
		t.Fatal(err)
	}
	if repeated, _ := loadRecents(); len(repeated) != 2 {
		t.Fatal("recente duplicado")
	}
	if err := os.Remove(entries[1].Path); err != nil {
		t.Fatal(err)
	}
	if entries[1].Available() {
		t.Fatal("recente removido continua disponível")
	}
	launcher.refreshRecents()
	if len(launcher.Recents.Objects) != 2 {
		t.Fatal("um recente indisponível impediu listar o outro")
	}
	launcher.openRecent(entries[0])
	if len(opened) != 3 {
		t.Fatal("recente válido não reabriu")
	}
	launcher.Window.Show()
	// capture é evidência opcional do launcher real, sem distribuir outro formato de compilação.
	if capture := os.Getenv("JORNADA_LAUNCHER_CAPTURE"); capture != "" {
		file, err := os.Create(capture)
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(file, launcher.Window.Canvas().Capture()); err != nil {
			file.Close()
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	// Uma configuração corrompida não pode impedir a tela de boas-vindas nem criar documento.
	config, _ := recentFile()
	if err := os.WriteFile(filepath.Clean(config), []byte("{quebrado"), 0600); err != nil {
		t.Fatal(err)
	}
	launcher.refreshRecents()
	if launcher.Message.Text == "" {
		t.Fatal("configuração corrompida sem diagnóstico")
	}
}
