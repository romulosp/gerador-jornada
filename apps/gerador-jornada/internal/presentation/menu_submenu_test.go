package presentation

import (
	"testing"

	"gerador-jornada.local/editor/internal/screen"
)

// TestMegaSubmenuLegacyScene confere a faixa e a geometria literais do PnlMenu.java.
func TestMegaSubmenuLegacyScene(t *testing.T) {
	// component e menu representam uma nova instância com os dois filhos do XML original.
	component, menu := newMenuComponent(t)
	// node identifica o ramo por seu item raiz, não por um movimento compartilhado.
	node := menu.Root.Items[0].Submenu
	// group precisa conter a faixa específica e os botões normal/especial centralizados.
	group := ApplicationMenuState(component, node)
	if len(primitivesByName(group, "faixa")) != 0 {
		t.Fatal("Mega-Sena sem tituloPrincipal ganhou faixa principal inventada")
	}
	// header é o PNG horizontal já desenhado pelo legado, sem duplicação de texto.
	header := primitivesByName(group, "faixa-jogo")
	if len(header) != 1 || header[0].SkinAsset != "/azul/icon_megasena_grande.png" || header[0].Text != "" || header[0].Bounds != (screen.Bounds{X: 0, Y: 150, Width: 800, Height: 27}) {
		t.Fatalf("faixa da Mega-Sena divergente: %+v", header)
	}
	// index percorre os dois botões e compara caixa pintada e região clicável.
	for index, item := range node.Items {
		// shapes e expected comprovam as medidas Java 164×52, em x220/400 e y310.
		shapes := primitivesByName(group, item.ID)
		expected := screen.Bounds{X: 220 + index*180, Y: 310, Width: 164, Height: 52}
		if len(shapes) != 1 || shapes[0].Bounds != expected || shapes[0].Text != "" {
			t.Fatalf("botão %d fora do layout legado: %+v", index, shapes)
		}
		if areas := MenuHitAreas(component, node); areas[index].Bounds != expected {
			t.Fatalf("clique %d fora da pintura: %+v", index, areas[index])
		}
	}
	if len(primitivesByName(group, "command/back")) != 1 || len(primitivesByName(group, "command/cancel")) != 1 {
		t.Fatal("submenu sem Voltar/Cancelar originais")
	}
}
