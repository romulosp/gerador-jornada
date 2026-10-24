package presentation

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestSkinAppliesFontAndColor confere a resolução de fonte/cor do acervo no texto.
func TestSkinAppliesFontAndColor(t *testing.T) {
	// component usa a skin label_titulo_1.
	component, _ := screen.NewComponent(catalog.Label, "rotulo", 20, 20)
	component.Style.Skin = "label_titulo_1"
	component.Properties.Label.Text = "Título"
	// group contém a caixa de texto com os valores do skin.
	group := Component(component)
	if len(group.Shapes) != 1 {
		t.Fatalf("formas inesperadas: %d", len(group.Shapes))
	}
	text := group.Shapes[0]
	if text.FontSize != 14 || text.Color != "#3E578D" {
		t.Fatalf("skin não aplicado: %+v", text)
	}
}

// TestSkinAppliesBackgroundAndImage confere fundo e imagem decorativa do acervo.
func TestSkinAppliesBackgroundAndImage(t *testing.T) {
	// component usa uma skin de rótulo com imagem embutida.
	component, _ := screen.NewComponent(catalog.Label, "logon", 20, 20)
	component.Style.Skin = "label_logon"
	// group contém a imagem antes do texto.
	group := Component(component)
	if len(group.Shapes) < 2 || group.Shapes[0].SkinAsset != "/azul/label_logon.png" {
		t.Fatalf("imagem de skin ausente: %+v", group.Shapes)
	}
}

// TestDefaultUsesLegacyColors garante que ausência de override resolve o padrão do legado.
func TestDefaultUsesLegacyColors(t *testing.T) {
	// component não escolhe skin.
	component, _ := screen.NewComponent(catalog.Label, "rotulo", 20, 20)
	group := Component(component)
	if len(group.Shapes) != 1 || group.Shapes[0].Color != "#5064AA" || group.Shapes[0].FontSize != 12 {
		t.Fatalf("padrão alterado: %+v", group.Shapes)
	}
}
