package catalog

import "testing"

// TestPropertiesVariants protege o discriminador e os enums do catálogo.
func TestPropertiesVariants(t *testing.T) {
	// properties contém deliberadamente duas variantes.
	properties := Properties{Button: &ButtonProperties{Variant: "confirm"}, Label: &TextProperties{}}
	if properties.Validate(Button) == nil {
		t.Fatal("duas variantes foram aceitas")
	}
	properties.Label = nil
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := properties.Validate(Button); err != nil {
		t.Fatal(err)
	}
	properties.Button.Variant = "execute"
	if properties.Validate(Button) == nil {
		t.Fatal("enum desconhecido foi aceito")
	}
	if len(Definitions()) != 7 {
		t.Fatal("catálogo deve conter sete famílias")
	}
}
