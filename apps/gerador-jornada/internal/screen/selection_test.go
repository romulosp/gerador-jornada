package screen

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// TestGroupMovementPreservesSpacing limita o grupo pela borda sem aproximar seus membros.
func TestGroupMovementPreservesSpacing(t *testing.T) {
	// page contém dois elementos com espaçamento conhecido e um terceiro que não deve mover.
	page := New("Grupo")
	for index := 0; index < 3; index++ {
		// component ocupa cinquenta unidades e mantém sua posição independente.
		component, _ := NewComponent(catalog.Label, string(rune('A'+index)), 10+index*100, 20)
		component.Bounds.Width, component.Bounds.Height = 50, 20
		page.Components = append(page.Components, component)
	}
	// ids seleciona somente os dois primeiros controles.
	ids := []string{page.Components[0].ID, page.Components[1].ID}
	if err := MoveSelection(&page, ids, 1000, -1000); err != nil {
		t.Fatal(err)
	}
	if page.Components[1].Bounds.X+50 != page.Size.Width || page.Components[1].Bounds.X-page.Components[0].Bounds.X != 100 || page.Components[0].Bounds.Y != 0 || page.Components[2].Bounds.X != 210 {
		t.Fatal("a limitação precisa preservar o conjunto e não mover itens não selecionados")
	}
	// before permite provar que uma referência inválida não provoca alteração parcial.
	before := page.Components[0].Bounds
	if err := MoveSelection(&page, []string{ids[0], "ausente"}, 10, 10); err == nil || page.Components[0].Bounds != before {
		t.Fatal("referência inválida alterou o grupo")
	}
}
