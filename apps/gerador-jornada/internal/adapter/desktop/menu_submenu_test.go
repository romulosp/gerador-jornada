package desktop

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestGameSubmenuRedraw percorre a prévia realmente reconstruída, não apenas o caminho no runtime.
func TestGameSubmenuRedraw(t *testing.T) {
	// app e window exercitam os widgets e a substituição da página após cada clique.
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("Submenus dos jogos")
	// closed evita fechar duas vezes a janela de teste após acionar o botão Sair real.
	closed := false
	window.SetOnClosed(func() { closed = true })
	defer func() {
		if !closed {
			window.Close()
		}
	}()
	// document e component preservam uma referência independente para verificar isolamento.
	document := application.NewDocument("Submenus do legado")
	component, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	component.Bounds = screen.Bounds{Width: 800, Height: 600}
	document.Screens[0].Components = append(document.Screens[0].Components, component)
	builtinassets.DefaultBackgrounds(&document)
	// session começa na raiz da tela atual, sem navegação externa.
	session, err := application.NewScreenPreview(document, document.Screens[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	// status expõe somente a seleção simulada da visita.
	status := widget.NewLabel("")
	// redraw precisa substituir o conteúdo para capturar os defeitos que um teste só de estado omite.
	var redraw func()
	redraw = func() {
		window.SetContent(previewPage(session, window, redraw, status, container.NewHBox()))
	}
	redraw()
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
	// tap usa um controle da vista atual e permite verificar o conteúdo após o callback.
	tap := func(id string) {
		// object procura somente a identidade, sem inferir ações por rótulo ou movimento.
		for _, object := range window.Content().(*fyne.Container).Objects {
			if control, ok := object.(*legacyControl); ok && control.menuItem.ID == id {
				test.Tap(control)
				return
			}
		}
		t.Fatalf("controle visível ausente: %s", id)
	}
	// capture registra opcionalmente a janela já reconstruída, sem gerar artefatos de distribuição.
	capture := func(name string) {
		if dir := os.Getenv("JORNADA_SUBMENU_CAPTURE"); dir != "" {
			// file mantém a evidência exclusivamente no diretório informado para o teste.
			file, err := os.Create(filepath.Join(dir, name+".png"))
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
	capture("submenu-raiz")
	// index percorre os quatro jogos cujo ramo contém normal/especial no menu.xml ativo.
	for _, index := range []int{0, 1, 2, 5} {
		// rootItem conserva o ramo original e seus metadados antes de qualquer clique.
		rootItem := component.Properties.Menu.Root.Items[index]
		tap(rootItem.ID)
		// node precisa corresponder à vista que foi efetivamente reconstruída.
		node := session.MenuNode(component.ID)
		if node == nil || node.ID != rootItem.Submenu.ID || len(node.Items) != 2 {
			t.Fatalf("ramo não aberto: %s", rootItem.Label)
		}
		if index == 0 {
			capture("submenu-megasena")
		}
		// item representa duas folhas com o mesmo código, mas nomes de movimento distintos.
		for _, item := range node.Items {
			tap(item.ID)
			// selected comprova que normal/especial não foram fundidos pelo código compartilhado.
			selected, ok := session.MenuSelectedItem(component.ID)
			if !ok || selected.ID != item.ID || selected.MovementName != item.MovementName || !strings.Contains(status.Text, item.MovementName) {
				t.Fatalf("movimento selecionado incorretamente: %+v / %s", selected, status.Text)
			}
		}
		tap("command/back")
		if session.MenuNode(component.ID).ID != component.Properties.Menu.Root.ID {
			t.Fatal("Voltar não reconstruiu a raiz")
		}
		if _, ok := session.MenuSelectedItem(component.ID); ok {
			t.Fatal("Voltar manteve seleção de outro nó")
		}
	}
	tap(component.Properties.Menu.Root.Items[0].ID)
	tap("command/cancel")
	if session.MenuNode(component.ID).ID != component.Properties.Menu.Root.ID {
		t.Fatal("Cancelar não restaurou a raiz")
	}
	if component.Properties.Menu.Root.Items[0].Submenu.Items[0].MovementName != "Mega-Sena" || component.Properties.Menu.InitialMenuID != component.Properties.Menu.Root.ID {
		t.Fatal("interação modificou o documento")
	}
	tap("command/exit")
	// runtime registra somente o encerramento desta visita, sem alterar a raiz confirmada.
	runtime, _ := session.MenuState(component.ID)
	if !runtime.Exited {
		t.Fatal("Sair não encerrou a sessão local")
	}
}
