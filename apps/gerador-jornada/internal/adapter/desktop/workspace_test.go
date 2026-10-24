package desktop

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestWorkspaceCreateImport confere colisões, raiz não recursiva, pacote inválido e origem intacta.
func TestWorkspaceCreateImport(t *testing.T) {
	// workspace fica num diretório temporário, sem tocar projetos reais ou configurações do usuário.
	workspace, err := OpenWorkspace(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	name := "Arrecadação: boleto / cobrança?"
	first, document, err := workspace.Create(name, false)
	if err != nil {
		t.Fatal(err)
	}
	if document.Project.Name != name || len(document.Screens[0].Components) != 0 || len(document.Project.Assets) == 0 {
		t.Fatal("criação em branco perdeu nome ou fundo")
	}
	second, _, err := workspace.Create(name, false)
	if err != nil || second == first {
		t.Fatal("homônimo substituiu a primeira jornada")
	}
	outside := filepath.Join(t.TempDir(), "importada.jornada")
	if err := projectzip.Save(context.Background(), outside, document); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(outside)
	imported, err := workspace.Import(outside)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(outside)
	copyBytes, _ := os.ReadFile(filepath.Join(workspace.Root, imported))
	if !bytes.Equal(before, after) || !bytes.Equal(before, copyBytes) {
		t.Fatal("importação alterou bytes da origem ou da cópia")
	}
	// nested e broken não podem ocultar jornadas válidas nem inserir descendentes recursivamente.
	nested := filepath.Join(workspace.Root, "subpasta")
	if err := os.Mkdir(nested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := projectzip.Save(context.Background(), filepath.Join(nested, "oculta.jornada"), document); err != nil {
		t.Fatal(err)
	}
	broken := filepath.Join(t.TempDir(), "quebrada.jornada")
	if err := os.WriteFile(broken, []byte("corrompido"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Import(broken); err == nil {
		t.Fatal("arquivo inválido foi publicado")
	}
	entries, err := workspace.List()
	if err != nil || len(entries) != 3 {
		t.Fatalf("listagem da raiz: %d / %v", len(entries), err)
	}
	if _, err := workspace.Resolve("../fora.jornada"); err == nil {
		t.Fatal("escape aceito")
	}
	if err := workspace.publishNew(first, []byte("substituição indevida")); err == nil {
		t.Fatal("publicação concorrente sobrescreveu arquivo")
	}
	if _, err := workspace.Load(first); err != nil {
		t.Fatalf("publicação concorrente danificou original: %v", err)
	}
}

// TestWorkspaceHundredJourneys mede a listagem real de cem pacotes, sem travar limite temporal em CI.
func TestWorkspaceHundredJourneys(t *testing.T) {
	// workspace e data isolam a medição do custo de gerar cada documento da fixture.
	workspace, err := OpenWorkspace(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	data, err := projectzip.Encode(NewJourneyDocument("Jornada de escala", false))
	if err != nil {
		t.Fatal(err)
	}
	for index := range 100 {
		if err := workspace.publishNew(fmt.Sprintf("jornada-%03d.jornada", index), data); err != nil {
			t.Fatal(err)
		}
	}
	start := time.Now()
	entries, err := workspace.List()
	elapsed := time.Since(start)
	if err != nil || len(entries) != 100 {
		t.Fatalf("listagem incompleta: %d / %v", len(entries), err)
	}
	t.Logf("Listagem e validação de 100 jornadas: %s", elapsed)
}

// findWorkspaceButton procura um botão realmente renderizado no diálogo ou na janela de teste.
func findWorkspaceButton(object fyne.CanvasObject, text string, visited map[fyne.CanvasObject]bool) *widget.Button {
	if object == nil || visited[object] {
		return nil
	}
	visited[object] = true
	if button, ok := object.(*widget.Button); ok && button.Text == text {
		return button
	}
	// children percorre containers e renderers, incluindo o PopUp que hospeda a confirmação.
	var children []fyne.CanvasObject
	if group, ok := object.(*fyne.Container); ok {
		children = group.Objects
	} else if control, ok := object.(fyne.Widget); ok {
		children = test.WidgetRenderer(control).Objects()
	}
	for _, child := range children {
		if button := findWorkspaceButton(child, text, visited); button != nil {
			return button
		}
	}
	return nil
}

// hasWorkspaceLabel verifica o texto efetivamente criado pelo renderer virtualizado da árvore.
func hasWorkspaceLabel(object fyne.CanvasObject, text string, visited map[fyne.CanvasObject]bool) bool {
	if object == nil || visited[object] {
		return false
	}
	visited[object] = true
	if label, ok := object.(*widget.Label); ok && object.Visible() && strings.Contains(label.Text, text) {
		return true
	}
	// children inclui o scroller e os nós visíveis, sem depender só dos callbacks de dados.
	var children []fyne.CanvasObject
	if group, ok := object.(*fyne.Container); ok {
		children = group.Objects
	} else if control, ok := object.(fyne.Widget); ok {
		children = test.WidgetRenderer(control).Objects()
	}
	for _, child := range children {
		if hasWorkspaceLabel(child, text, visited) {
			return true
		}
	}
	return false
}

// TestWorkspacePreviewConfirmation percorre cliques reais, cancelamento, árvore e navegação C2/C3.
func TestWorkspacePreviewConfirmation(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app, editor e preview compartilham o mesmo workspace de teste, nunca projetos do usuário.
	app := test.NewApp()
	defer app.Quit()
	document := NewJourneyDocument("Principal", false)
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	menu.Bounds = screen.Bounds{Width: 800, Height: 600}
	document.Screens[0].Components = append(document.Screens[0].Components, menu)
	workspace, _ := OpenWorkspace(t.TempDir())
	source := filepath.Join(workspace.Root, "principal.jornada")
	if err := projectzip.Save(context.Background(), source, document); err != nil {
		t.Fatal(err)
	}
	editor := New(app, document)
	defer editor.close()
	editor.SetProjectPath(source)
	editor.Window.Show()
	for _, label := range []string{filepath.Base(workspace.Root), "principal.jornada"} {
		if !hasWorkspaceLabel(editor.WorkspaceTree, label, map[fyne.CanvasObject]bool{}) {
			t.Fatalf("árvore não renderizou %s", label)
		}
	}
	state, err := application.NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	window := app.NewWindow("Criação e navegação entre jornadas")
	defer window.Close()
	editor.configureWorkspacePreview(state, window)
	status := widget.NewLabel("")
	var redraw func()
	redraw = func() { window.SetContent(previewPage(state, window, redraw, status, container.NewHBox())) }
	redraw()
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
	// tap aciona somente um item atualmente visível na cena legada.
	tap := func(id string) {
		for _, object := range window.Content().(*fyne.Container).Objects {
			if control, ok := object.(*legacyControl); ok && control.menuItem.ID == id {
				test.Tap(control)
				return
			}
		}
		t.Fatalf("item não renderizado: %s", id)
	}
	// confirm usa os botões reais do diálogo, verificando inclusive o rótulo em português.
	confirm := func(label string) {
		button := findWorkspaceButton(window.Canvas().Overlays().Top(), label, map[fyne.CanvasObject]bool{})
		if button == nil {
			t.Fatalf("diálogo imediato não ofereceu %s", label)
		}
		test.Tap(button)
	}
	rootItem := menu.Properties.Menu.Root.Items[0]
	leaf := rootItem.Submenu.Items[0]
	tap(rootItem.ID)
	before, _ := os.ReadFile(source)
	tap(leaf.ID)
	confirm("Cancelar")
	after, _ := os.ReadFile(source)
	if !bytes.Equal(before, after) || len(editor.WorkspaceTree.ChildUIDs("workspace")) != 1 || len(editor.Session.UndoStack) != 0 {
		t.Fatal("cancelar diálogo modificou origem, árvore ou histórico")
	}
	tap(leaf.ID)
	confirm("Criar jornada")
	if len(editor.WorkspaceTree.ChildUIDs("workspace")) != 2 {
		t.Fatal("criação confirmada não atualizou árvore")
	}
	tap(leaf.ID)
	if state.Document.Project.Name != leaf.Label || len(state.JourneyHistory) != 1 {
		t.Fatal("clique após criação não abriu destino")
	}
	if err := state.BackJourney(); err != nil || state.Document.Project.ID != document.Project.ID || state.MenuNode(menu.ID).ID != rootItem.Submenu.ID {
		t.Fatalf("voltar não restaurou origem e submenu: %v", err)
	}
	// capture registra a prévia real depois de voltar; o executável tem evidência nativa separada.
	if capture := os.Getenv("JORNADA_WORKSPACE_CAPTURE"); capture != "" {
		redraw()
		file, err := os.Create(capture)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		if err := png.Encode(file, window.Canvas().Capture()); err != nil {
			t.Fatal(err)
		}
	}
	// Selecionar o arquivo na árvore abre o mesmo destino persistido no disco.
	saved, err := projectzip.Load(source)
	if err != nil {
		t.Fatal(err)
	}
	target := saved.Screens[0].Components[0].Properties.Menu.Root.Items[0].Submenu.Items[0].TargetJourney
	editor.WorkspaceTree.Select(target)
	if editor.Session.Document.Project.Name != leaf.Label || filepath.Base(editor.FilePath) != target {
		t.Fatal("seleção da árvore não abriu jornada persistida")
	}
}

// TestLinkedJourneyAtomicHistory confirma vínculo persistido, rollback e histórico sem alterar cancelamento.
func TestLinkedJourneyAtomicHistory(t *testing.T) {
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	// app e editor fornecem o mesmo caso de uso chamado após a confirmação do diálogo.
	app := test.NewApp()
	defer app.Quit()
	document := NewJourneyDocument("Principal", false)
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	document.Screens[0].Components = append(document.Screens[0].Components, menu)
	workspace, _ := OpenWorkspace(t.TempDir())
	path := filepath.Join(workspace.Root, "principal.jornada")
	if err := projectzip.Save(context.Background(), path, document); err != nil {
		t.Fatal(err)
	}
	editor := New(app, document)
	defer editor.close()
	editor.SetProjectPath(path)
	preview, err := application.NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	preview.CurrentJourney, preview.LoadJourney = "principal.jornada", workspace.Load
	// item é uma folha normal de Mega-Sena; criar destino não pode trocar sua identidade ou movimento.
	item := menu.Properties.Menu.Root.Items[0].Submenu.Items[0]
	before, _ := os.ReadFile(path)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := editor.createLinkedJourney(canceled, preview, menu.ID, item.ID); err == nil {
		t.Fatal("operação cancelada confirmou destino")
	}
	after, _ := os.ReadFile(path)
	entries, _ := workspace.List()
	if !bytes.Equal(before, after) || len(entries) != 1 || len(editor.Session.UndoStack) != 0 {
		t.Fatal("cancelamento publicou estado parcial")
	}
	filename, err := editor.createLinkedJourney(context.Background(), preview, menu.ID, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	saved, err := projectzip.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Screens[0].Components[0].Properties.Menu.Root.Items[0].Submenu.Items[0].TargetJourney != filename {
		t.Fatal("vínculo não foi persistido")
	}
	_, _ = preview.MenuActivate(menu.ID, menu.Properties.Menu.Root.Items[0].ID)
	if outcome, err := preview.MenuActivate(menu.ID, item.ID); err != nil || outcome != application.MenuNavigated {
		t.Fatalf("próxima ativação não navegou: %v %v", outcome, err)
	}
	if err := editor.Session.Undo(); err != nil {
		t.Fatal(err)
	}
	if editor.Session.Document.Screens[0].Components[0].Properties.Menu.Root.Items[0].Submenu.Items[0].TargetJourney != "" {
		t.Fatal("Undo não restaurou vínculo anterior")
	}
	if err := editor.Session.Redo(); err != nil {
		t.Fatal(err)
	}
	if editor.Session.Document.Screens[0].Components[0].Properties.Menu.Root.Items[0].Submenu.Items[0].TargetJourney != filename {
		t.Fatal("Redo não restaurou destino confirmado")
	}
}
