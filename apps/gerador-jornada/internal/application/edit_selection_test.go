package application

import "testing"

// TestSelectionDuplicateAndDelete comprova que duplicar e remover conjuntos são ações reversíveis.
func TestSelectionDuplicateAndDelete(t *testing.T) {
	// session usa a fixture com teclado vinculado a campo da mesma tela.
	session := NewSession(Example())
	// page e ids selecionam campo e teclado para verificar remapeamento do vínculo.
	page := &session.Document.Screens[0]
	ids := []string{page.Components[2].ID, page.Components[3].ID}
	// initial preserva a quantidade de controles antes da duplicação.
	initial := len(page.Components)
	// created identifica os novos controles sem reutilizar IDs da origem.
	created, err := session.DuplicateComponents(page.ID, ids)
	if err != nil {
		t.Fatal(err)
	}
	page = &session.Document.Screens[0]
	if len(page.Components) != initial+2 || page.ComponentByID(created[1]).Properties.Keyboard.TargetComponentID != created[0] {
		t.Fatal("duplicação perdeu o vínculo interno entre teclado e campo")
	}
	if err := session.DeleteComponents(page.ID, created); err != nil || len(session.Document.Screens[0].Components) != initial {
		t.Fatalf("remoção coletiva incorreta: %v", err)
	}
	if err := session.Undo(); err != nil || len(session.Document.Screens[0].Components) != initial+2 {
		t.Fatalf("remoção não reversível em um comando: %v", err)
	}
}
