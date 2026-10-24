package application

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// TestCompoundsExampleShape confere o conjunto de teste usado na revisão manual.
func TestCompoundsExampleShape(t *testing.T) {
	// document contém duas telas, com entrada/ficha na primeira e menu na segunda.
	document := CompoundsExample()
	if len(document.Screens) != 2 {
		t.Fatalf("esperadas duas telas, obtidas %d", len(document.Screens))
	}
	// kinds registra os três compostos encontrados.
	kinds := map[catalog.Kind]bool{}
	for _, page := range document.Screens {
		for _, component := range page.Components {
			kinds[component.Kind] = true
			if err := component.Properties.Validate(component.Kind); err != nil {
				t.Fatalf("%s inválido: %v", component.Kind, err)
			}
		}
	}
	for _, kind := range []catalog.Kind{catalog.BarcodeInput, catalog.BarcodeProcessed, catalog.ApplicationMenu} {
		if !kinds[kind] {
			t.Fatalf("composto ausente: %s", kind)
		}
	}
	if problems := document.Validate(); len(problems) > 0 {
		t.Fatalf("documento de teste inválido: %s", ProblemsText(problems))
	}
}
