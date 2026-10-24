package projectzip

import (
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestSkinPersistence preserva a escolha de skin no manifesto 3/3.
func TestSkinPersistence(t *testing.T) {
	// document contém um rótulo personalizado pelo acervo.
	document := application.NewDocument("skin")
	component, _ := screen.NewComponent(catalog.Label, "rotulo", 20, 20)
	component.Style.Skin = "label_titulo_1"
	document.Screens[0].Components = append(document.Screens[0].Components, component)
	data, err := EncodeManifest(document)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := DecodeManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Screens[0].Components[0].Style.Skin != "label_titulo_1" {
		t.Fatalf("skin perdido: %+v", restored.Screens[0].Components[0].Style)
	}
}
