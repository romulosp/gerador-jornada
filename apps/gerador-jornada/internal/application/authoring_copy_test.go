package application

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// TestDuplicateRegeneratesAllInstanceIDs impede compartilhar UUID de opção/linha entre cópias.
func TestDuplicateRegeneratesAllInstanceIDs(t *testing.T) {
	// document inclui dez Kinds e partes compostas; rádio/relatório ganham seleção e linha explícitas.
	document := mcpFixture(t)
	pageID := document.Screens[0].ID
	document.Screens[0].Components[4].Properties.Radio.SelectedOptionID = document.Screens[0].Components[4].Properties.Radio.Options[1].ID
	document.Screens[0].Components[5].Properties.Report.Rows = []catalog.Option{{ID: project.NewID(), Text: "Linha original"}}
	prior := documentEntityIDs(&document, "")
	session := NewSession(document)
	id, mapping, err := session.DuplicateScreenWithMapping(pageID)
	if err != nil {
		t.Fatal(err)
	}
	clone := session.Document.Screen(id)
	if clone == nil || prior[id] || mapping[pageID] != id {
		t.Fatal("identidade da tela não remapeada")
	}
	for index, component := range clone.Components {
		// original/oldIDs/newIDs verificam correspondência da cópia em ordem estável.
		original := document.Screens[0].Components[index]
		if prior[component.ID] || mapping[original.ID] != component.ID {
			t.Fatal("identidade de componente reaproveitada")
		}
		oldIDs, newIDs := original.Properties.InstanceIDs(), component.Properties.InstanceIDs()
		if len(oldIDs) != len(newIDs) {
			t.Fatal("cópia perdeu partes internas")
		}
		for position, newID := range newIDs {
			if prior[newID] || mapping[oldIDs[position]] != newID {
				t.Fatal("identidade interna reaproveitada ou ausente no mapa")
			}
		}
	}
	if clone.Components[4].Properties.Radio.SelectedOptionID != mapping[document.Screens[0].Components[4].Properties.Radio.SelectedOptionID] {
		t.Fatal("seleção inicial não remapeada")
	}
	if clone.Components[6].Properties.Keyboard.TargetComponentID != mapping[document.Screens[0].Components[3].ID] {
		t.Fatal("teclado não remapeado")
	}
	// O comando de seleção usa a mesma regra para opções e linhas, sem criar tela adicional.
	selection := NewSession(document)
	ids, err := selection.DuplicateComponents(pageID, []string{document.Screens[0].Components[4].ID, document.Screens[0].Components[5].ID})
	if err != nil {
		t.Fatal(err)
	}
	for _, copiedID := range ids {
		for _, internalID := range selection.Document.Screen(pageID).ComponentByID(copiedID).Properties.InstanceIDs() {
			if prior[internalID] {
				t.Fatal("cópia da seleção reaproveitou opção/linha")
			}
		}
	}
}
