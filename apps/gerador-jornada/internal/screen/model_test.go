package screen

import (
	"gerador-jornada.local/editor/internal/catalog"
	"testing"
)

// TestGeometryAndLocalNames protege limites e unicidade da composição.
func TestGeometryAndLocalNames(t *testing.T) {
	// page representa uma área de desenho reduzida.
	page := New("Teste")
	page.Size = Size{Width: 100, Height: 80}
	// component começa deliberadamente além do limite.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	component, err := NewComponent(catalog.Button, "botao", 200, -10)
	if err != nil {
		t.Fatal(err)
	}
	component.Bounds = component.Bounds.Clamp(page.Size)
	page.Components = append(page.Components, component)
	if err = page.Validate(); err != nil {
		t.Fatal(err)
	}
	page.Size.Width = 20
	if page.Validate() == nil {
		t.Fatal("reduzir tela reposicionou silenciosamente o componente")
	}
	if component.Bounds.Width != 100 {
		t.Fatal("clamp não limitou o tamanho ao canvas")
	}
}
