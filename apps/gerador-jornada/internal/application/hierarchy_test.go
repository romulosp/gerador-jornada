package application

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestDuplicateScreenRemapsHierarchy confere que pai e filho copiados mantêm o vínculo interno.
func TestDuplicateScreenRemapsHierarchy(t *testing.T) {
	// document contém um pai habilitado e um filho desabilitado.
	document := NewDocument("hierarquia")
	parent, _ := screen.NewComponent(catalog.Button, "pai", 20, 20)
	child, _ := screen.NewComponent(catalog.Label, "filho", 40, 40)
	child.ParentID = parent.ID
	child.Enabled = false
	document.Screens[0].Components = append(document.Screens[0].Components, parent, child)
	session := NewSession(document)
	// duplicated contém a nova tela com identidades próprias.
	duplicated, err := session.DuplicateScreen(document.Screens[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	page := session.Document.Screen(duplicated)
	if page == nil || len(page.Components) != 2 {
		t.Fatal("tela duplicada incompleta")
	}
	if page.Components[1].ParentID != page.Components[0].ID {
		t.Fatalf("ParentID não remapeado: %+v", page.Components)
	}
	if page.Components[1].Enabled {
		t.Fatal("estado desabilitado não preservado na cópia")
	}
	if page.Components[0].ID == parent.ID || page.Components[1].ID == child.ID {
		t.Fatal("identidade reutilizada na cópia")
	}
}

// TestDeleteComponentsOrphansChildren mantém filhos na tela sem referência inválida.
func TestDeleteComponentsOrphansChildren(t *testing.T) {
	// document contém pai e filho na mesma tela.
	document := NewDocument("hierarquia")
	parent, _ := screen.NewComponent(catalog.Button, "pai", 20, 20)
	child, _ := screen.NewComponent(catalog.Label, "filho", 40, 40)
	child.ParentID = parent.ID
	document.Screens[0].Components = append(document.Screens[0].Components, parent, child)
	session := NewSession(document)
	if err := session.DeleteComponents(document.Screens[0].ID, []string{parent.ID}); err != nil {
		t.Fatal(err)
	}
	page := session.Document.Screens[0]
	if len(page.Components) != 1 || page.Components[0].ParentID != "" {
		t.Fatalf("filho não foi preservado como raiz: %+v", page.Components)
	}
	if err := page.Validate(); err != nil {
		t.Fatalf("tela inválida após remoção: %v", err)
	}
}
