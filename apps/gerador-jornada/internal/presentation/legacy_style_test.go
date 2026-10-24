package presentation

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestLegacyButtonDefaultAndOverrides impede o retorno do retângulo genérico sem skin.
func TestLegacyButtonDefaultAndOverrides(t *testing.T) {
	// component é um botão sem personalização, como os inseridos pela paleta.
	component, _ := screen.NewComponent(catalog.Button, "confirmar", 30, 40)
	// group contém imagem decorativa e legenda editável independentes.
	group := Component(component)
	if len(group.Shapes) != 2 || group.Shapes[0].SkinAsset != "/azul/botao_aux2_released.png" {
		t.Fatalf("botão sem decoração original: %+v", group.Shapes)
	}
	if group.Shapes[1].Color != "#00FF00" || group.Shapes[1].FontFamily != "Futura Medium" || group.Shapes[1].FontSize != 12 {
		t.Fatalf("tipografia original não resolvida: %+v", group.Shapes[1])
	}
	component.Style.TextColor, component.Style.FontSize = "#123456", 19
	// overridden preserva os valores explícitos sobre os defaults do acervo.
	overridden := Component(component)
	if overridden.Shapes[1].Color != "#123456" || overridden.Shapes[1].FontSize != 19 {
		t.Fatal("override de aparência perdido")
	}
}

// TestLegacyKeyboardHasFourteenKeys verifica o arranjo Java e a escala dos projetos antigos.
func TestLegacyKeyboardHasFourteenKeys(t *testing.T) {
	// component usa a dimensão aprovada para novas inserções.
	component, _ := screen.NewComponent(catalog.NumericKeyboard, "teclado", 10, 20)
	if component.Bounds.Height != 350 {
		t.Fatalf("altura padrão: %d", component.Bounds.Height)
	}
	// group preserva uma imagem e um texto por tecla, sem rasterizar a entrada.
	group := Component(component)
	if len(group.Shapes) != 28 {
		t.Fatalf("esperadas 14 teclas: %d formas", len(group.Shapes))
	}
	if group.Shapes[0].Bounds != (screen.Bounds{X: 15, Y: 25, Width: 52, Height: 52}) || group.Shapes[1].Text != "7" {
		t.Fatalf("primeira tecla fora da referência Java: %+v", group.Shapes[:2])
	}
	if group.Shapes[25].Text != "TAB" || group.Shapes[27].Text != "LIMPAR" {
		t.Fatal("teclas auxiliares ausentes")
	}
	component.Bounds.Height = 280
	// legacy mantém o tamanho persistido e escala apenas as partes derivadas.
	legacy := Component(component)
	if legacy.Bounds.Height != 280 || legacy.Shapes[0].Bounds.Height != 41 {
		t.Fatal("teclado antigo foi redimensionado em vez de preservar seus bounds")
	}
}

// TestLegacyReportKeepsOriginalDecorations verifica as partes do relatório sem cobrir sua imagem.
func TestLegacyReportKeepsOriginalDecorations(t *testing.T) {
	// component contém uma linha de conteúdo editável.
	component, _ := screen.NewComponent(catalog.Report, "relatorio", 20, 20)
	component.Properties.Report.Rows = []catalog.Option{{ID: "linha", Text: "Pagamento recebido"}}
	// group permite inspecionar a decoração sem depender do renderer Fyne.
	group := Component(component)
	// assets e content identificam recursos originais e o texto editável.
	assets, content := map[string]bool{}, false
	for _, primitive := range group.Shapes {
		assets[primitive.SkinAsset] = true
		content = content || primitive.Text == "Pagamento recebido"
	}
	if !assets["/azul/label_tabela_cabecalho.png"] || !assets["/azul/icon_tabela_rodape_2.png"] || !content {
		t.Fatalf("relatório incompleto: %+v", group.Shapes)
	}
}

// TestUnavailableStyleIsVisible impede substituir uma skin inexistente sem diagnóstico.
func TestUnavailableStyleIsVisible(t *testing.T) {
	// component conserva uma escolha antiga indisponível na biblioteca atual.
	component, _ := screen.NewComponent(catalog.Button, "antigo", 0, 0)
	component.Style.Skin = "botao_loteca_especial"
	if ValidateStyle(component) == nil {
		t.Fatal("recurso ausente aceito")
	}
	// group deve informar a causa mantendo os dados originais do componente.
	group := Component(component)
	if len(group.Shapes) != 1 || group.Shapes[0].Color != "#B91C1C" || group.Shapes[0].Text == "" || group.Shapes[0].SkinAsset != "" {
		t.Fatalf("substituição silenciosa de skin: %+v", group.Shapes)
	}
	component.Style.Skin, component.Style.FontFamily = "", "Fonte inexistente"
	if ValidateStyle(component) == nil {
		t.Fatal("fonte não incorporada aceita")
	}
}
