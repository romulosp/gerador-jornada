package desktop

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestMenuMovementFields confirma edição de metadados sem misturar as duas folhas de código 702.
func TestMenuMovementFields(t *testing.T) {
	// app e editor fornecem o formulário real e seus destinos disponíveis.
	app := test.NewApp()
	defer app.Quit()
	editor := New(app, application.NewDocument("Metadados do menu"))
	defer editor.close()
	// component mantém o documento original separado do rascunho e do alvo publicado.
	component, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	// special é a segunda folha do ramo Mega-Sena, com o mesmo código da primeira.
	special := component.Properties.Menu.Root.Items[0].Submenu.Items[1]
	// form e apply são exatamente o editor de árvore e sua confirmação atômica.
	form := widget.NewForm()
	apply := editor.menuFields(form, component)
	// fields e selector localizam controles públicos pelo rótulo apresentado ao autor.
	fields := map[string]*widget.Entry{}
	var selector *widget.Select
	// field percorre os controles já montados, sem invocar closures internas de edição.
	for _, field := range form.Items {
		if field.Text == "Nó ou item" {
			selector = field.Widget.(*widget.Select)
		}
		if entry, ok := field.Widget.(*widget.Entry); ok {
			fields[field.Text] = entry
		}
	}
	if selector == nil {
		t.Fatal("seletor de item ausente")
	}
	// entry encontra o caminho de identidade da folha que será editada.
	for _, entry := range menuTreeEntries(component.Properties.Menu) {
		if entry.ItemID == special.ID {
			selector.SetSelected(entry.Label)
		}
	}
	if fields["Nome do movimento"] == nil || fields["Código do movimento"].Text != "702" || fields["Nome do movimento"].Text != "Mega-Sena Especial" {
		t.Fatal("metadados originais não carregados no formulário")
	}
	fields["Código do movimento"].SetText("00702")
	fields["Nome do movimento"].SetText("Mega-Sena Especial — teste")
	fields["Opção do movimento"].SetText("099")
	fields["Atalho legado"].SetText("F12")
	// target recebe cópia independente; confirmar não deve reescrever o componente de origem.
	target := component
	target.Properties = component.Properties.Copy()
	if err := apply(&target); err != nil {
		t.Fatal(err)
	}
	// changed conserva zeros e Unicode em texto; a folha normal não herda os metadados especiais.
	changed := target.Properties.Menu.Root.Items[0].Submenu.Items[1]
	if changed.MovementCode != "00702" || changed.MovementName != "Mega-Sena Especial — teste" || changed.MovementOption != "099" || changed.Shortcut != "F12" || target.Properties.Menu.Root.Items[0].Submenu.Items[0].MovementName != "Mega-Sena" {
		t.Fatalf("confirmação dos metadados divergente: %+v", changed)
	}
	if component.Properties.Menu.Root.Items[0].Submenu.Items[1].MovementName != special.MovementName {
		t.Fatal("rascunho alterou a origem antes de confirmar")
	}
}
