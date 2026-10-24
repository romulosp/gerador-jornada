package presentation

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// MenuBaseWidth e MenuBaseHeight definem a geometria de referência do menu da aplicação.
const (
	// MenuBaseWidth expressa a largura de referência em unidades lógicas.
	MenuBaseWidth = 800
	// MenuBaseHeight expressa a altura de referência em unidades lógicas.
	MenuBaseHeight = 600
)

// MenuHitArea identifica a região clicável de um item sem depender de legenda ou posição.
type MenuHitArea struct {
	// ComponentID identifica a instância composta na tela.
	ComponentID string
	// ItemID identifica o item ou nó acionado.
	ItemID string
	// Bounds informa a área absoluta em unidades lógicas.
	Bounds screen.Bounds
}

// ApplicationMenu expande o estado inicial confirmado para canvas e PPT estático.
func ApplicationMenu(component screen.Component) Group {
	if component.Properties.Menu == nil {
		return Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds}
	}
	return ApplicationMenuState(component, component.Properties.Menu.Initial())
}

// ApplicationMenuState expande o nó visível com a geometria do layout contratado.
// A prévia usa o caminho local; canvas e PPT usam o menu inicial confirmado.
func ApplicationMenuState(component screen.Component, node *catalog.SubmenuProperties) Group {
	// group preserva a identidade do composto e sua ordem de camadas.
	group := Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds, Shapes: []Primitive{}}
	if node == nil {
		return group
	}
	if node.Layout == catalog.MenuLayoutRoot {
		// services reúne os itens da coluna esquerda; games, os marcados como jogo.
		services, games := []catalog.MenuItem{}, []catalog.MenuItem{}
		for _, item := range node.Items {
			if item.Category == "JOGOS" {
				games = append(games, item)
			} else {
				services = append(services, item)
			}
		}
		// item percorre os serviços em duas colunas alternadas.
		for index, item := range services {
			// bounds posiciona o botão sem reduzir o tamanho de referência.
			bounds := MenuButtonBounds(component.Bounds, catalog.MenuLayoutRoot, index)
			bounds.X += (index % 2) * 180 * component.Bounds.Width / MenuBaseWidth
			if bounds.X+bounds.Width > component.Bounds.X+400*component.Bounds.Width/MenuBaseWidth {
				continue
			}
			group.Shapes = append(group.Shapes, menuButton(bounds, item, component))
		}
		// item percorre os jogos na metade direita.
		for index, item := range games {
			// bounds usa duas colunas à direita e cinco linhas.
			bounds := screen.Bounds{
				X:      component.Bounds.X + (420+(index%2)*180)*component.Bounds.Width/MenuBaseWidth,
				Y:      component.Bounds.Y + (150+(index/2)*73)*component.Bounds.Height/MenuBaseHeight,
				Width:  164 * component.Bounds.Width / MenuBaseWidth,
				Height: 68 * component.Bounds.Height / MenuBaseHeight,
			}
			group.Shapes = append(group.Shapes, menuButton(bounds, item, component))
		}
		return group
	}
	// title reserve a faixa somente em submenus, nunca na raiz.
	title := node.Title
	if title == "" {
		title = node.Name
	}
	// titleState conserva a tipografia branca da faixa principal do legado.
	titleState, _ := skin.MustBuiltin().StateByType("label", "label_faixa_titulo_principal", "default")
	// band mantém texto e decoração sem deslocar os botões do submenu.
	band := styledText(component, titleState, title, screen.Bounds{X: component.Bounds.X + 20*component.Bounds.Width/MenuBaseWidth, Y: component.Bounds.Y + 20*component.Bounds.Height/MenuBaseHeight, Width: 760 * component.Bounds.Width / MenuBaseWidth, Height: 40 * component.Bounds.Height / MenuBaseHeight})
	band.SkinAsset, band.Name, band.TextAlign, band.VerticalAlign = "/azul/faixatituloprincipal.png", "faixa", "center", "center"
	group.Shapes = append(group.Shapes, band)
	if node.Layout == catalog.MenuLayoutBolao {
		// item percorre os três botões com legenda da composição digital.
		for index, item := range node.Items {
			if index >= 3 {
				break
			}
			bounds := screen.Bounds{X: component.Bounds.X + (20+index*180)*component.Bounds.Width/MenuBaseWidth, Y: component.Bounds.Y + 150*component.Bounds.Height/MenuBaseHeight, Width: 164 * component.Bounds.Width / MenuBaseWidth, Height: 68 * component.Bounds.Height / MenuBaseHeight}
			group.Shapes = append(group.Shapes, menuButton(bounds, item, component))
			group.Shapes = append(group.Shapes, Primitive{Bounds: screen.Bounds{X: bounds.X, Y: bounds.Y + bounds.Height + 4, Width: bounds.Width, Height: 34 * component.Bounds.Height / MenuBaseHeight}, Text: item.Legend, Color: "#334155", FontSize: 10, Wrap: true, Name: "legenda"})
		}
		return group
	}
	// columns, buttonWidth, buttonHeight e pitch distinguem a grade financeira da grade comum.
	columns, buttonWidth, buttonHeight, pitch := 4, 164, 60, 70
	if node.Layout == catalog.MenuLayoutFinancial {
		columns, buttonWidth, buttonHeight, pitch = 3, 140, 80, 92
	}
	// item percorre a grade do submenu selecionado.
	for index, item := range node.Items {
		// bounds distribui os itens em colunas e linhas estáveis.
		bounds := screen.Bounds{
			X:      component.Bounds.X + (20+(index%columns)*(buttonWidth+21))*component.Bounds.Width/MenuBaseWidth,
			Y:      component.Bounds.Y + (90+(index/columns)*pitch)*component.Bounds.Height/MenuBaseHeight,
			Width:  buttonWidth * component.Bounds.Width / MenuBaseWidth,
			Height: buttonHeight * component.Bounds.Height / MenuBaseHeight,
		}
		group.Shapes = append(group.Shapes, menuButton(bounds, item, component))
	}
	return group
}

// MenuButtonBounds devolve a área do item na raiz ou na grade comum.
func MenuButtonBounds(bounds screen.Bounds, layout string, index int) screen.Bounds {
	if layout == catalog.MenuLayoutRoot {
		return screen.Bounds{
			X:      bounds.X + 20*bounds.Width/MenuBaseWidth,
			Y:      bounds.Y + (150+(index/2)*73)*bounds.Height/MenuBaseHeight,
			Width:  164 * bounds.Width / MenuBaseWidth,
			Height: 68 * bounds.Height / MenuBaseHeight,
		}
	}
	// columns, buttonWidth, buttonHeight e pitch repetem a grade da pintura.
	columns, buttonWidth, buttonHeight, pitch := 4, 164, 60, 70
	if layout == catalog.MenuLayoutFinancial {
		columns, buttonWidth, buttonHeight, pitch = 3, 140, 80, 92
	}
	return screen.Bounds{
		X:      bounds.X + (20+(index%columns)*(buttonWidth+21))*bounds.Width/MenuBaseWidth,
		Y:      bounds.Y + (90+(index/columns)*pitch)*bounds.Height/MenuBaseHeight,
		Width:  buttonWidth * bounds.Width / MenuBaseWidth,
		Height: buttonHeight * bounds.Height / MenuBaseHeight,
	}
}

// MenuHitAreas devolve as regiões clicáveis do estado visível, na mesma escala da pintura.
func MenuHitAreas(component screen.Component, node *catalog.SubmenuProperties) []MenuHitArea {
	// areas preserva a ordem dos itens exibidos.
	areas := []MenuHitArea{}
	if node == nil {
		return areas
	}
	if node.Layout == catalog.MenuLayoutRoot {
		// services e games repetem a separação usada pela pintura.
		services, games := []catalog.MenuItem{}, []catalog.MenuItem{}
		for _, item := range node.Items {
			if item.Category == "JOGOS" {
				games = append(games, item)
			} else {
				services = append(services, item)
			}
		}
		for index, item := range services {
			area := MenuButtonBounds(component.Bounds, catalog.MenuLayoutRoot, index)
			area.X += (index % 2) * 180 * component.Bounds.Width / MenuBaseWidth
			if area.X+area.Width > component.Bounds.X+400*component.Bounds.Width/MenuBaseWidth {
				continue
			}
			areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: area})
		}
		for index, item := range games {
			areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: screen.Bounds{
				X:      component.Bounds.X + (420+(index%2)*180)*component.Bounds.Width/MenuBaseWidth,
				Y:      component.Bounds.Y + (150+(index/2)*73)*component.Bounds.Height/MenuBaseHeight,
				Width:  164 * component.Bounds.Width / MenuBaseWidth,
				Height: 68 * component.Bounds.Height / MenuBaseHeight,
			}})
		}
		return areas
	}
	// item percorre os itens do submenu visível.
	for index, item := range node.Items {
		if node.Layout == catalog.MenuLayoutBolao && index >= 3 {
			break
		}
		// area reutiliza a geometria da grade correspondente.
		area := screen.Bounds{}
		if node.Layout == catalog.MenuLayoutBolao {
			area = screen.Bounds{X: component.Bounds.X + (20+index*180)*component.Bounds.Width/MenuBaseWidth, Y: component.Bounds.Y + 150*component.Bounds.Height/MenuBaseHeight, Width: 164 * component.Bounds.Width / MenuBaseWidth, Height: 68 * component.Bounds.Height / MenuBaseHeight}
		} else {
			area = MenuButtonBounds(component.Bounds, node.Layout, index)
		}
		areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: area})
	}
	return areas
}

// menuButton pinta um item com estado habilitado ou desabilitado e legenda editável.
func menuButton(bounds screen.Bounds, item catalog.MenuItem, parent screen.Component) Primitive {
	// name resolve a skin pelo papel do item, sem inferir imagem a partir de caminho externo.
	name := "botao_menu"
	if item.Submenu != nil {
		name = "botao_menu_seta"
	}
	if parent.Style.Skin != "" {
		name = parent.Style.Skin
	}
	if item.Variant != "" {
		name = item.Variant
	}
	// stateName mantém a cor e a imagem disabled originais.
	stateName := "default"
	if !item.Enabled {
		stateName = "disabled"
	}
	// state usa somente definições que possuem todos os recursos incorporados.
	state, _ := skin.MustBuiltin().State(catalog.Button, name, stateName)
	// component permite compartilhar o resolvedor de tipografia dos demais botões.
	component := screen.Component{Kind: catalog.Button, Bounds: bounds, Enabled: item.Enabled, Style: parent.Style}
	// result conserva imagem decorativa e texto editável na mesma região lógica.
	result := styledText(component, state, item.Label, bounds)
	result.Wrap = true
	result.SkinAsset, result.Name, result.TextAlign, result.VerticalAlign = state.ImagePath, item.ID, "center", "center"
	if parent.Style.TextAlign != "" {
		result.TextAlign = parent.Style.TextAlign
	}
	result.Stroke = parent.Style.BorderColor
	if parent.Style.FillColor != "" {
		result.SkinAsset, result.Fill = "", parent.Style.FillColor
	}
	return result
}
