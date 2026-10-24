package screen

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// newTestComponents cria três componentes válidos na mesma tela.
func newTestComponents(t *testing.T) Screen {
	t.Helper()
	// page recebe controles independentes para os testes de hierarquia e camada.
	page := New("teste")
	// kinds preserva a ordem das famílias usadas.
	kinds := []catalog.Kind{catalog.Button, catalog.Label, catalog.Title}
	for index, kind := range kinds {
		// component é criado com identidade e propriedades padrão.
		component, err := NewComponent(kind, "componente"+string(rune('A'+index)), 20+index*120, 20)
		if err != nil {
			t.Fatal(err)
		}
		page.Components = append(page.Components, component)
	}
	return page
}

// TestHierarchyValidation recusa autorreferência, pai ausente e ciclo.
func TestHierarchyValidation(t *testing.T) {
	// page contém três componentes válidos.
	page := newTestComponents(t)
	a, b, c := page.Components[0].ID, page.Components[1].ID, page.Components[2].ID
	// pai existente na mesma tela é aceito.
	page.Components[1].ParentID = a
	page.Components[2].ParentID = b
	if err := page.Validate(); err != nil {
		t.Fatalf("hierarquia válida recusada: %v", err)
	}
	// autorreferência é inválida.
	page.Components[0].ParentID = a
	if err := page.Validate(); err == nil {
		t.Fatal("autorreferência aceita")
	}
	// ciclo A→B→C→A é inválido.
	page.Components[0].ParentID = c
	if err := page.Validate(); err == nil {
		t.Fatal("ciclo aceito")
	}
	// pai inexistente é inválido.
	page.Components[0].ParentID = ""
	page.Components[0].ParentID = "00000000-0000-4000-8000-000000000000"
	if err := page.Validate(); err == nil {
		t.Fatal("pai inexistente aceito")
	}
}

// TestMoveSelectionMovesSubtreeOnce confere que pai e filho não recebem o deslocamento duas vezes.
func TestMoveSelectionMovesSubtreeOnce(t *testing.T) {
	// page contém pai e filho na mesma tela.
	page := newTestComponents(t)
	parent := page.Components[0]
	child := page.Components[1]
	page.Components[1].ParentID = parent.ID
	beforeParent, beforeChild := parent.Bounds, child.Bounds
	if err := MoveSelection(&page, []string{parent.ID, child.ID}, 10, 15); err != nil {
		t.Fatal(err)
	}
	if page.Components[0].Bounds.X != beforeParent.X+10 || page.Components[0].Bounds.Y != beforeParent.Y+15 {
		t.Fatal("pai não moveu uma vez")
	}
	if page.Components[1].Bounds.X != beforeChild.X+10 || page.Components[1].Bounds.Y != beforeChild.Y+15 {
		t.Fatal("filho não moveu uma vez")
	}
	// mover apenas o pai também desloca o filho.
	beforeChild = page.Components[1].Bounds
	if err := MoveSelection(&page, []string{parent.ID}, -5, 0); err != nil {
		t.Fatal(err)
	}
	if page.Components[1].Bounds.X != beforeChild.X-5 {
		t.Fatal("subárvore não acompanhou o pai")
	}
}

// TestPlaceComponentOrder cobre referência relativa, extremos e erros sem mutação.
func TestPlaceComponentOrder(t *testing.T) {
	// page contém três componentes na ordem A, B, C.
	page := newTestComponents(t)
	a, b, c := page.Components[0].ID, page.Components[1].ID, page.Components[2].ID
	// C atrás de B resulta em A, C, B.
	if err := PlaceComponent(&page, c, b, false); err != nil {
		t.Fatal(err)
	}
	if page.Components[1].ID != c || page.Components[2].ID != b {
		t.Fatalf("ordem relativa inesperada: %v", page.Components)
	}
	// A para a frente absoluta resulta em C, B, A.
	if err := PlaceComponent(&page, a, "", true); err != nil {
		t.Fatal(err)
	}
	if page.Components[2].ID != a {
		t.Fatal("origem não foi para a frente")
	}
	// C para o fundo absoluto resulta em C, B, A.
	if err := PlaceComponent(&page, c, "", false); err != nil {
		t.Fatal(err)
	}
	if page.Components[0].ID != c {
		t.Fatal("origem não foi para o fundo")
	}
	// erros preservam a ordem e todos os dados.
	order := append([]Component{}, page.Components...)
	if err := PlaceComponent(&page, a, a, true); err == nil {
		t.Fatal("autorreferência aceita")
	}
	if err := PlaceComponent(&page, a, "inexistente", true); err == nil {
		t.Fatal("referência inexistente aceita")
	}
	for index := range order {
		if order[index].ID != page.Components[index].ID {
			t.Fatal("erro alterou a ordem")
		}
	}
}
