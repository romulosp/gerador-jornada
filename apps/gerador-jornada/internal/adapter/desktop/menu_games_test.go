package desktop

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	desktopdriver "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// TestGameMenuPreview verifica bytes originais na prévia real, pressão, liberação e bloqueio.
func TestGameMenuPreview(t *testing.T) {
	// app e window exercitam o mesmo renderer Fyne do executável sem exigir Office.
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("Jogos com skins originais")
	defer window.Close()
	// document contém apenas o menu, para conferir também a distribuição completa dos jogos.
	document := application.NewDocument("Jogos do legado")
	// component usa a fábrica atual, com variantes persistidas para os dez jogos do menu.
	component, err := screen.NewComponent(catalog.ApplicationMenu, "menuAplicacao", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	component.Bounds = screen.Bounds{Width: 800, Height: 600}
	document.Screens[0].Components = append(document.Screens[0].Components, component)
	// session isola a interação das propriedades confirmadas do documento.
	session, err := application.NewScreenPreview(document, document.Screens[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	// status permite identificar ativação indevida de controles desabilitados.
	status := widget.NewLabel("")
	// view fornece os controles que o usuário realmente pressiona na prévia.
	view := newMenuView(component, session, window, status, func() {})
	window.SetContent(container.NewWithoutLayout(view.objects...))
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
	// capture publica evidência opcional, mantendo a execução habitual livre de artefatos.
	capture := func(suffix string) {
		if path := os.Getenv("JORNADA_GAME_CAPTURE"); path != "" {
			if suffix != "" {
				path = path[:len(path)-len(filepath.Ext(path))] + suffix + ".png"
			}
			// file recebe a imagem do canvas da janela já renderizada.
			file, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := png.Encode(file, window.Canvas().Capture()); err != nil {
				file.Close()
				t.Fatal(err)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
	// imageBytes inspeciona os bytes do recurso efetivamente consumido pelo renderer.
	imageBytes := func(control *legacyControl) []byte {
		// renderer contém a cena do item, e não o desenho padrão de um widget.Button.
		renderer := test.WidgetRenderer(control).(*legacyControlRenderer)
		// object percorre as partes pintadas em busca da decoração original.
		for _, object := range renderer.content.Objects {
			if image, ok := object.(*canvas.Image); ok && image.Resource != nil {
				return image.Resource.Content()
			}
		}
		t.Fatal("controle sem recurso visual")
		return nil
	}
	// million é o controle +Milionária cujo ciclo de pressão será exercitado.
	var million *legacyControl
	// mega permite verificar a abertura do submenu normal/especial pelo controle renderizado.
	var mega *legacyControl
	// count confirma que nenhum dos jogos foi coberto por superfície genérica.
	count := 0
	// object percorre os controles da composição interativa.
	for _, object := range view.objects {
		if control, ok := object.(*legacyControl); ok && control.menuItem.Category == "JOGOS" {
			count++
			if control.menuItem.Variant == "botao_milionaria" {
				million = control
			}
			if control.menuItem.Variant == "botao_megasena" {
				mega = control
			}
			if len(imageBytes(control)) == 0 {
				t.Fatal("jogo sem imagem")
			}
		}
	}
	if count != 10 || million == nil || mega == nil {
		t.Fatalf("jogos incompletos: %d", count)
	}
	capture("")
	// assertImage compara bytes conhecidos sem depender do nome atribuído pelo renderer.
	assertImage := func(path string) {
		// original conserva o arquivo exato mapeado pelo XML.
		original, err := skin.MustBuiltin().Asset(path)
		if err != nil || !bytes.Equal(original, imageBytes(million)) {
			t.Fatalf("estado não usa %s: %v", path, err)
		}
	}
	assertImage("/botao_milionaria.png")
	million.MouseDown(&desktopdriver.MouseEvent{Button: desktopdriver.MouseButtonPrimary})
	assertImage("/botao_milionaria_pressed.png")
	capture("-pressed")
	million.MouseOut()
	assertImage("/botao_milionaria.png")
	million.MouseDown(&desktopdriver.MouseEvent{Button: desktopdriver.MouseButtonPrimary})
	million.MouseUp(&desktopdriver.MouseEvent{Button: desktopdriver.MouseButtonPrimary})
	assertImage("/botao_milionaria.png")
	// Enabled pertence à cópia local; o estado disabled deve prevalecer até sobre uma pressão.
	million.menuItem.Enabled, million.component.Enabled = false, false
	million.Disable()
	million.MouseDown(&desktopdriver.MouseEvent{Button: desktopdriver.MouseButtonPrimary})
	million.Tapped(nil)
	assertImage("/botao_milionaria_disabled.png")
	capture("-disabled")
	if status.Text != "" || !component.Properties.Menu.Root.Items[9].Enabled {
		t.Fatal("item desabilitado executou ação ou alterou o documento")
	}
	mega.Tapped(nil)
	// node comprova que a ativação usa a árvore importada, sem inferir ações pela marca.
	node := session.MenuNode(component.ID)
	if node == nil || node.ID != component.Properties.Menu.Root.Items[0].Submenu.ID || len(node.Items) != 2 || node.Items[1].Variant != "botao_megasena_especial" {
		t.Fatal("controle não abriu o submenu normal/especial original")
	}
}
