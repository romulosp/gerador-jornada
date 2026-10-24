package catalog

import (
	"gerador-jornada.local/editor/internal/project"
	"testing"
)

// TestMenuJourneyDestination exige destino portátil, explícito e exclusivo dos demais destinos.
func TestMenuJourneyDestination(t *testing.T) {
	// menu usa o seed completo para provar que a extensão não altera os demais itens.
	menu := NewApplicationMenuProperties(project.NewID)
	item := &menu.Root.Items[0]
	item.TargetJourney = "Arrecadação de boleto.jornada"
	if menu.Validate() == nil {
		t.Fatal("ramo aceitou destino entre jornadas")
	}
	item.Submenu = nil
	item.TargetScreenID = project.NewID()
	if menu.Validate() == nil {
		t.Fatal("dois destinos foram aceitos")
	}
	item.TargetScreenID = ""
	if err := menu.Validate(); err != nil {
		t.Fatal(err)
	}
	// filename percorre escapes, volumes e nomes de dispositivo que não são referências de workspace.
	for _, filename := range []string{"../fora.jornada", `..\fora.jornada`, `C:\fora.jornada`, "/fora.jornada", "CON.jornada", "lpt1.jornada", "arquivo.txt", "inválido?.jornada", "quebra\n.jornada"} {
		item.TargetJourney = filename
		if menu.Validate() == nil {
			t.Fatalf("nome perigoso aceito: %q", filename)
		}
	}
}
