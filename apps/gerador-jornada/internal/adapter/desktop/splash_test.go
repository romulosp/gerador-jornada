package desktop

import (
	"context"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// TestSplashOnlyVideoAndLauncherCredit protege a última decisão visual, sem textos sobre o vídeo.
func TestSplashOnlyVideoAndLauncherCredit(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app e splash exercitam apenas construção/layout; reprodução real pertence ao smoke nativo.
	app := test.NewApp()
	defer app.Quit()
	splash := NewSplash(app, nil, nil)
	defer splash.Window.Close()
	if splash.Window.Content() != splash.VideoArea {
		t.Fatal("splash deve conter somente a superfície de vídeo")
	}
	launcher := NewLauncherInWindow(app, splash.Window, nil)
	if !hasWorkspaceLabel(launcher.Window.Content(), StudioCredit, map[fyne.CanvasObject]bool{}) {
		t.Fatal("autoria literal não foi apresentada no rodapé do launcher")
	}
	if StudioCredit != "Desenvolvido por Rômulo Penha - Canais Parceiro - 03" {
		t.Fatal("autoria não corresponde ao texto solicitado")
	}
}

// TestStartupResourcePreparation verifica carga operacional e cancelamento antes da preparação.
func TestStartupResourcePreparation(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// progress registra a conclusão real das etapas, mesmo sem controles visíveis na splash.
	progress := 0.0
	if err := LoadStartupResources(context.Background(), func(value float64, message string) { progress = value }); err != nil {
		t.Fatal(err)
	}
	if progress != 1 {
		t.Fatal("preparação não concluiu seus recursos")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := LoadStartupResources(ctx, func(value float64, message string) {}); err == nil {
		t.Fatal("carregamento cancelado foi aceito")
	}
}
