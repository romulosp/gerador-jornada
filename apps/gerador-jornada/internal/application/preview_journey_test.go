package application

import (
	"errors"
	"reflect"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestPreviewJourneyRoundTrip verifica entrada inicial, restauração de runtime, escopo e falha sem substituição.
func TestPreviewJourneyRoundTrip(t *testing.T) {
	// source e destination são documentos distintos; os nomes não resolvem destinos implicitamente.
	source := NewDocument("Origem")
	destination := NewDocument("Destino")
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	menu.Properties.Menu.Root.Items[0].Submenu.Items[1].TargetJourney = "destino.jornada"
	source.Screens[0].Components = append(source.Screens[0].Components, menu)
	preview, err := NewPreview(source)
	if err != nil {
		t.Fatal(err)
	}
	preview.CurrentJourney = "principal.jornada"
	preview.LoadJourney = func(filename string) (Document, error) {
		if filename != "destino.jornada" {
			t.Fatalf("rota inferida: %s", filename)
		}
		return destination, nil
	}
	if _, err := preview.MenuActivate(menu.ID, menu.Properties.Menu.Root.Items[0].ID); err != nil {
		t.Fatal(err)
	}
	// item é a folha especial no submenu; o código compartilhado não determina a jornada.
	item := menu.Properties.Menu.Root.Items[0].Submenu.Items[1]
	if outcome, err := preview.MenuActivate(menu.ID, item.ID); err != nil || outcome != MenuNavigated {
		t.Fatalf("navegação: %v %v", outcome, err)
	}
	if preview.Document.Project.ID != destination.Project.ID || preview.ScreenID != destination.Journey.StartScreenID || len(preview.JourneyHistory) != 1 {
		t.Fatal("destino não entrou pela tela inicial")
	}
	if err := preview.BackJourney(); err != nil {
		t.Fatal(err)
	}
	if preview.CurrentJourney != "principal.jornada" || preview.MenuNode(menu.ID).ID != menu.Properties.Menu.Root.Items[0].Submenu.ID {
		t.Fatal("Voltar perdeu jornada ou caminho de origem")
	}
	// before conserva o estado que uma falha de IO não pode substituir.
	before := preview.Document.Copy()
	preview.LoadJourney = func(string) (Document, error) { return Document{}, ErrJourneyMissing }
	if _, err := preview.MenuActivate(menu.ID, item.ID); !errors.Is(err, ErrJourneyMissing) || !reflect.DeepEqual(before, preview.Document) || len(preview.JourneyHistory) != 0 {
		t.Fatalf("falha substituiu origem: %v", err)
	}
	// current representa a prévia restrita; seu loader nunca pode ser invocado.
	current, err := NewScreenPreview(source, source.Screens[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	current.LoadJourney = func(string) (Document, error) { t.Fatal("Prévia Tela Atual seguiu destino"); return destination, nil }
	_, _ = current.MenuActivate(menu.ID, menu.Properties.Menu.Root.Items[0].ID)
	if outcome, err := current.MenuActivate(menu.ID, item.ID); err != nil || outcome != MenuSimulated || current.Document.Project.ID != source.Project.ID {
		t.Fatalf("escopo restrito: %v %v", outcome, err)
	}
}
