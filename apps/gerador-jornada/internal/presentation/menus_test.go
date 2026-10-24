package presentation

import (
	"fmt"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// newMenuComponent cria o menu seed na geometria de referência de 800 por 600.
func newMenuComponent(t *testing.T) (screen.Component, *catalog.ApplicationMenuProperties) {
	t.Helper()
	component, err := screen.NewComponent(catalog.ApplicationMenu, "menuAplicacao", 0, 0)
	if err != nil {
		t.Fatalf("componente de menu inválido: %v", err)
	}
	component.Bounds = screen.Bounds{Width: 800, Height: 600}
	return component, component.Properties.Menu
}

// menuNodeByName procura um nó da árvore pelo nome editável.
func menuNodeByName(node *catalog.SubmenuProperties, name string) *catalog.SubmenuProperties {
	if node == nil {
		return nil
	}
	if node.Name == name {
		return node
	}
	for _, item := range node.Items {
		if found := menuNodeByName(item.Submenu, name); found != nil {
			return found
		}
	}
	return nil
}

// menuNodeByLayout devolve o primeiro nó com o layout procurado.
func menuNodeByLayout(node *catalog.SubmenuProperties, layout string) *catalog.SubmenuProperties {
	if node == nil {
		return nil
	}
	if node.Layout == layout {
		return node
	}
	for _, item := range node.Items {
		if found := menuNodeByLayout(item.Submenu, layout); found != nil {
			return found
		}
	}
	return nil
}

// splitMenuColumns separa os itens da raiz entre serviços e jogos.
func splitMenuColumns(node catalog.SubmenuProperties) ([]catalog.MenuItem, []catalog.MenuItem) {
	services, games := []catalog.MenuItem{}, []catalog.MenuItem{}
	for _, item := range node.Items {
		if item.Category == "JOGOS" {
			games = append(games, item)
		} else {
			services = append(services, item)
		}
	}
	return services, games
}

// TestApplicationMenuRootPaint confere a raiz sem faixa, as duas colunas e o total de botões.
func TestApplicationMenuRootPaint(t *testing.T) {
	component, menu := newMenuComponent(t)
	group := ApplicationMenuState(component, &menu.Root)
	if faixa := primitivesByName(group, "faixa"); len(faixa) != 0 {
		t.Fatalf("raiz não pode pintar faixa de título: %+v", faixa)
	}
	services, games := splitMenuColumns(menu.Root)
	if len(services) != 12 || len(games) != 10 {
		t.Fatalf("raiz deve ter 12 serviços e 10 jogos: %d/%d", len(services), len(games))
	}
	// 12 botões de serviço e 10 de jogo preenchem a raiz sem outras formas.
	if len(group.Shapes) != 22 {
		t.Fatalf("quantidade de botões da raiz divergente: %d", len(group.Shapes))
	}
	// expected preserva a ordem de pintura: serviços e depois jogos.
	expected := append(append([]catalog.MenuItem{}, services...), games...)
	for index, item := range expected {
		shape := group.Shapes[index]
		if shape.Name != item.ID || shape.Text != item.Label {
			t.Fatalf("botão %d divergente: %+v", index+1, shape)
		}
		if shape.SkinAsset == "" || shape.FontAsset != "/fonts/futurbci.ttf" || shape.FontFamily != "Futura BdCn BT" {
			t.Fatalf("estilo do botão %d divergente: %+v", index+1, shape)
		}
	}
	// a primeira coluna começa em x=20, topo 150, com botões 164 por 68.
	if group.Shapes[0].Bounds != (screen.Bounds{X: 20, Y: 150, Width: 164, Height: 68}) {
		t.Fatalf("primeiro botão fora da geometria contratada: %+v", group.Shapes[0].Bounds)
	}
	// a segunda coluna alterna para x=200.
	if group.Shapes[1].Bounds != (screen.Bounds{X: 200, Y: 150, Width: 164, Height: 68}) {
		t.Fatalf("segundo botão fora da geometria contratada: %+v", group.Shapes[1].Bounds)
	}
	// o wrapper estático usa a escolha inicial, que por padrão é a raiz.
	if wrapper := ApplicationMenu(component); len(wrapper.Shapes) != 22 {
		t.Fatalf("menu estático deve usar a raiz inicial: %d", len(wrapper.Shapes))
	}
}

// TestMenuHitAreasRoot confere que o clique coincide exatamente com a pintura da raiz.
func TestMenuHitAreasRoot(t *testing.T) {
	component, menu := newMenuComponent(t)
	group := ApplicationMenuState(component, &menu.Root)
	areas := MenuHitAreas(component, &menu.Root)
	if len(areas) != len(group.Shapes) {
		t.Fatalf("áreas e botões divergem: %d != %d", len(areas), len(group.Shapes))
	}
	for index, area := range areas {
		shape := group.Shapes[index]
		if area.ComponentID != component.ID || area.ItemID != shape.Name {
			t.Fatalf("identidade da área %d divergente: %+v", index+1, area)
		}
		if area.Bounds != shape.Bounds {
			t.Fatalf("geometria da área %d divergente: %+v != %+v", index+1, area.Bounds, shape.Bounds)
		}
	}
}

// TestApplicationMenuSubmenuPaint confere a faixa e a grade de um submenu comum.
func TestApplicationMenuSubmenuPaint(t *testing.T) {
	component, menu := newMenuComponent(t)
	// convenios não possui itens no seed, mas ainda exibe a faixa do título.
	convenios := menuNodeByName(&menu.Root, "CONVÊNIOS")
	if convenios == nil || convenios.Layout != catalog.MenuLayoutStandard {
		t.Fatalf("submenu CONVÊNIOS inválido: %+v", convenios)
	}
	empty := ApplicationMenuState(component, convenios)
	if faixa := primitivesByName(empty, "faixa"); len(faixa) != 1 || faixa[0].Text != "CONVÊNIOS" {
		t.Fatalf("faixa do submenu vazio divergente: %+v", faixa)
	}
	if len(empty.Shapes) != 1 || len(MenuHitAreas(component, convenios)) != 0 {
		t.Fatalf("submenu vazio não deve pintar itens: %d", len(empty.Shapes))
	}
	// meios exibe três itens na grade comum.
	meios := menuNodeByName(&menu.Root, "MEIOS DE PAGAMENTO")
	if meios == nil || meios.Layout != catalog.MenuLayoutStandard || len(meios.Items) != 3 {
		t.Fatalf("submenu MEIOS DE PAGAMENTO inválido: %+v", meios)
	}
	group := ApplicationMenuState(component, meios)
	if len(group.Shapes) != 1+len(meios.Items) {
		t.Fatalf("formas do submenu divergentes: %d", len(group.Shapes))
	}
	if group.Shapes[0].Name != "faixa" || group.Shapes[0].Text != "MEIOS DE PAGAMENTO" {
		t.Fatalf("faixa do submenu divergente: %+v", group.Shapes[0])
	}
	areas := MenuHitAreas(component, meios)
	if len(areas) != len(meios.Items) {
		t.Fatalf("áreas do submenu divergentes: %d", len(areas))
	}
	for index, area := range areas {
		shape := group.Shapes[1+index]
		if area.ItemID != shape.Name || area.ItemID != meios.Items[index].ID {
			t.Fatalf("item %d do submenu divergente: %+v", index+1, area)
		}
		if area.Bounds != shape.Bounds {
			t.Fatalf("área %d do submenu não coincide com a pintura: %+v", index+1, area.Bounds)
		}
	}
}

// TestApplicationMenuBolaoLegends confere os três botões com legenda do layout bolao.
func TestApplicationMenuBolaoLegends(t *testing.T) {
	component, menu := newMenuComponent(t)
	digital := menuNodeByLayout(&menu.Root, catalog.MenuLayoutBolao)
	if digital == nil || len(digital.Items) != 3 {
		t.Fatalf("BOLÕES DIGITAIS inválido: %+v", digital)
	}
	group := ApplicationMenuState(component, digital)
	// 1 faixa + 3 botões + 3 legendas.
	if len(group.Shapes) != 7 {
		t.Fatalf("layout bolao deve produzir sete formas: %d", len(group.Shapes))
	}
	if group.Shapes[0].Name != "faixa" || group.Shapes[0].Text != "BOLÕES DIGITAIS" {
		t.Fatalf("faixa do bolão divergente: %+v", group.Shapes[0])
	}
	// wantLegends fixa os três textos literais do contrato, na ordem.
	wantLegends := []string{
		"Visualizar cotas digitais para venda na lotérica.",
		"Consulta bolões de concursos abertos para comercialização.",
		"Consulta bolões de concursos encerrados e não prescritos.",
	}
	legends := primitivesByName(group, "legenda")
	if len(legends) != 3 {
		t.Fatalf("legendas do bolão divergentes: %d", len(legends))
	}
	for index, legend := range legends {
		if legend.Text != wantLegends[index] {
			t.Fatalf("legenda %d divergente: %q", index+1, legend.Text)
		}
		if legend.Color != "#334155" || legend.FontSize != 10 || !legend.Wrap {
			t.Fatalf("estilo da legenda %d divergente: %+v", index+1, legend)
		}
		button := group.Shapes[1+2*index]
		if button.Name != digital.Items[index].ID || button.Text != digital.Items[index].Label {
			t.Fatalf("botão %d do bolão divergente: %+v", index+1, button)
		}
		if legend.Bounds.Y != button.Bounds.Y+button.Bounds.Height+4 {
			t.Fatalf("legenda %d fora da posição contratada: %+v", index+1, legend.Bounds)
		}
	}
	// apenas os três botões visíveis são áreas clicáveis, sem as legendas.
	areas := MenuHitAreas(component, digital)
	if len(areas) != 3 {
		t.Fatalf("áreas do bolão divergentes: %d", len(areas))
	}
	for index, area := range areas {
		button := group.Shapes[1+2*index]
		if area.ItemID != button.Name || area.Bounds != button.Bounds {
			t.Fatalf("área %d do bolão não coincide com o botão: %+v", index+1, area)
		}
	}
}

// TestApplicationMenuFinancialGrid confere a grade financeira de três colunas e 140 por 80.
func TestApplicationMenuFinancialGrid(t *testing.T) {
	// node é construído à mão porque o seed não oferece filhos financeiros.
	items := []catalog.MenuItem{}
	for index := 0; index < 4; index++ {
		items = append(items, catalog.MenuItem{ID: project.NewID(), Label: fmt.Sprintf("Financeiro %d", index+1), Category: "SERVIÇOS", Enabled: true})
	}
	node := &catalog.SubmenuProperties{ID: project.NewID(), Name: "FINANCEIROS", Title: "FINANCEIROS", Layout: catalog.MenuLayoutFinancial, Items: items}
	component, _ := newMenuComponent(t)
	group := ApplicationMenuState(component, node)
	if len(group.Shapes) != 1+len(items) {
		t.Fatalf("formas da grade financeira divergentes: %d", len(group.Shapes))
	}
	// item percorre a grade financeira conferindo nome, texto e dimensão contratada.
	for index, item := range items {
		shape := group.Shapes[1+index]
		if shape.Name != item.ID || shape.Text != item.Label {
			t.Fatalf("item financeiro %d divergente: %+v", index+1, shape)
		}
		if shape.Bounds.Width != 140 || shape.Bounds.Height != 80 {
			t.Fatalf("botão financeiro %d deve medir 140 por 80: %+v", index+1, shape.Bounds)
		}
	}
	// os três primeiros botões formam a primeira linha e o quarto inicia a segunda.
	first, second, third, fourth := group.Shapes[1].Bounds, group.Shapes[2].Bounds, group.Shapes[3].Bounds, group.Shapes[4].Bounds
	if first.Y != second.Y || second.Y != third.Y || !(first.X < second.X && second.X < third.X) {
		t.Fatalf("primeira linha financeira fora de três colunas: %+v %+v %+v", first, second, third)
	}
	if fourth.Y <= first.Y || fourth.X != first.X {
		t.Fatalf("quarto botão deve iniciar a segunda linha na primeira coluna: %+v", fourth)
	}
	areas := MenuHitAreas(component, node)
	if len(areas) != len(items) {
		t.Fatalf("áreas financeiras divergentes: %d", len(areas))
	}
	for index, area := range areas {
		if area.Bounds != group.Shapes[1+index].Bounds || area.ItemID != items[index].ID {
			t.Fatalf("área financeira %d não coincide com a pintura: %+v", index+1, area)
		}
	}
}
