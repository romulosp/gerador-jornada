package desktop

import (
	"image/png"
	"os"
	"path/filepath"
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

// TestLegacyPreviewRuntime exercita cursor, TAB, exclusão e isolamento sobre a cena renderizada.
func TestLegacyPreviewRuntime(t *testing.T) {
	// app e window usam o renderer Fyne em memória, sem simular homologação nativa.
	app := test.NewApp()
	defer app.Quit()
	window := app.NewWindow("Referência visual EFL")
	defer window.Close()
	// document reúne as famílias em uma fixture legível e sem dependências externas.
	document := application.NewDocument("Estilo legado")
	// page recebe apenas componentes confirmados; runtime ficará em outro snapshot.
	page := &document.Screens[0]
	// add mantém identidade, posição e conteúdo das famílias inseridas pela paleta.
	add := func(kind catalog.Kind, name string, x, y int) *screen.Component {
		// component é uma nova instância tipada, com defaults reais da paleta.
		component, err := screen.NewComponent(kind, name, x, y)
		if err != nil {
			t.Fatal(err)
		}
		page.Components = append(page.Components, component)
		return &page.Components[len(page.Components)-1]
	}
	add(catalog.Title, "titulo", 40, 80).Properties.Title.Text = "Componentes do legado"
	// firstID e secondID permitem conferir a ordem de TAB sem usar legendas.
	firstID := add(catalog.TextField, "campo", 40, 145).ID
	secondID := add(catalog.TextField, "campoDois", 40, 185).ID
	add(catalog.Button, "confirmar", 40, 245)
	add(catalog.Radio, "escolha", 40, 330)
	// report conserva título e linhas em formas separadas sobre a decoração original.
	report := add(catalog.Report, "relatorio", 250, 140)
	report.Properties.Report.Title = "Atendimento"
	report.Properties.Report.Rows = []catalog.Option{{ID: "linha", Text: "Pagamento recebido"}, {ID: "linha2", Text: "Operação confirmada"}}
	// keyboard vincula as teclas ao primeiro campo, preservando cursores locais.
	keyboard := add(catalog.NumericKeyboard, "teclado", 580, 140)
	keyboard.Properties.Keyboard.TargetComponentID = firstID
	builtinassets.DefaultBackgrounds(&document)
	// state recebe uma cópia independente do documento, incluindo valores e escolhas.
	state, err := application.NewScreenPreview(document, page.ID)
	if err != nil {
		t.Fatal(err)
	}
	// content é a composição de prévia com suas regiões interativas reais.
	content := previewPage(state, window, func() {}, widget.NewLabel(""), container.NewHBox())
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
	// fields e keys permitem acionar os widgets renderizados em vez de comandos artificiais.
	fields := map[string]*legacyEntry{}
	keys := map[string]*legacyControl{}
	// surface confirma o nível de contêiner criado pela prévia.
	surface, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("prévia sem contêiner")
	}
	// object percorre as superfícies exatamente na ordem de pintura.
	for _, object := range surface.Objects {
		switch control := object.(type) {
		case *container.ThemeOverride:
			// field identifica entradas com captura transparente e fonte real.
			field, isField := control.Content.(*legacyEntry)
			if isField {
				fields[field.component.ID] = field
			}
		case *legacyControl:
			if control.component.Kind == catalog.Button {
				keys[control.component.Properties.Button.Text] = control
			}
		}
	}
	if len(fields) != 2 || keys["TAB"] == nil || keys["LIMPAR"] == nil {
		t.Fatal("prévia sem campos ou teclas originais")
	}
	// first é o campo vinculado; suas teclas preservam a posição do cursor.
	first := fields[firstID]
	window.Canvas().Focus(first)
	first.SetText("12")
	first.CursorColumn = 1
	keys["7"].Tapped(nil)
	if first.Text != "172" {
		t.Fatalf("inserção deve respeitar o cursor: %q", first.Text)
	}
	keys["LIMPAR"].Tapped(nil)
	if first.Text != "12" {
		t.Fatalf("LIMPAR deve emitir Backspace: %q", first.Text)
	}
	keys["TAB"].Tapped(nil)
	if window.Canvas().Focused() != fields[secondID] {
		t.Fatal("TAB não avançou ao próximo campo")
	}
	keys["8"].Tapped(nil)
	if fields[secondID].Text != "8" {
		t.Fatal("teclado não acompanhou o foco após TAB")
	}
	if document.Screens[0].ComponentByID(firstID).Properties.Field.InitialValue != "" || document.Screens[0].ComponentByID(secondID).Properties.Field.InitialValue != "" {
		t.Fatal("edição runtime alterou o documento")
	}
	if first.textPrimitive().FontAsset != "/fonts/arial.ttf" {
		t.Fatal("captura de texto perdeu a face original")
	}
	// capture é uma evidência opcional do renderer em memória, separada do aceite humano.
	if capture := os.Getenv("JORNADA_VISUAL_CAPTURE"); capture != "" {
		if err := os.MkdirAll(filepath.Dir(capture), 0755); err != nil {
			t.Fatal(err)
		}
		// file recebe a captura da composição realmente renderizada por este teste.
		file, err := os.Create(capture)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		if err := png.Encode(file, window.Canvas().Capture()); err != nil {
			t.Fatal(err)
		}
	}
}
