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
func ApplicationMenuState(component screen.Component, node *catalog.SubmenuProperties) (group Group) {
	// group preserva a identidade do composto e sua ordem de camadas.
	group = Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds, Shapes: []Primitive{}}
	if node == nil {
		return group
	}
	// O rodapé compartilha cena e hit-test, mas nunca entra na árvore confirmada.
	defer func() {
		// command conserva a intenção local e a skin original do botão de navegação.
		for _, command := range MenuCommands(component, node) {
			// shapes permitem também overrides explícitos sem controles temáticos opacos.
			shapes := menuItemAt(command.Bounds, catalog.MenuItem{ID: "command/" + command.ID, Label: command.Label, Variant: command.Variant, Enabled: true}, component)
			// index identifica as partes do comando, separadas dos itens persistidos.
			for index := range shapes {
				shapes[index].Role = "menu-control"
			}
			group.Shapes = append(group.Shapes, shapes...)
		}
	}()
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
			bounds := menuRootServiceBounds(component.Bounds, index)
			bounds.X += (index % 2) * 180 * component.Bounds.Width / MenuBaseWidth
			if bounds.X+bounds.Width > component.Bounds.X+400*component.Bounds.Width/MenuBaseWidth {
				continue
			}
			group.Shapes = append(group.Shapes, menuItemAt(bounds, item, component)...)
		}
		// item percorre os jogos na metade direita.
		for index, item := range games {
			// bounds usa duas colunas à direita e preserva a ordem dos itens autorados.
			bounds := screen.Bounds{
				X:      component.Bounds.X + (420+(index%2)*180)*component.Bounds.Width/MenuBaseWidth,
				Y:      component.Bounds.Y + (150+(index/2)*77)*component.Bounds.Height/MenuBaseHeight,
				Width:  164 * component.Bounds.Width / MenuBaseWidth,
				Height: 52 * component.Bounds.Height / MenuBaseHeight,
			}
			group.Shapes = append(group.Shapes, menuItemAt(bounds, item, component)...)
		}
		return group
	}
	// A faixa principal só existe quando tituloPrincipal foi fornecido; não inventar título por nome.
	if node.Title != "" {
		// titleState e panelState conservam fonte e decoração próprias do legado.
		titleState, _ := skin.MustBuiltin().StateByType("label", "label_faixa_titulo_principal", "default")
		panelState, _ := skin.MustBuiltin().StateByType("panel", "painel_titulo_principal", "default")
		// band usa a faixa inteira; o texto reserva os 50 pixels da margem direita do Java.
		band := Primitive{Bounds: menuScaled(component.Bounds, screen.Bounds{Y: 100, Width: 800, Height: 28}), SkinAsset: panelState.ImagePath, Name: "faixa-fundo"}
		group.Shapes = append(group.Shapes, band)
		// title conserva texto editável sobre o PNG de fundo compartilhado.
		title := styledText(component, titleState, node.Title, menuScaled(component.Bounds, screen.Bounds{Y: 100, Width: 750, Height: 28}))
		title.Name, title.TextAlign, title.VerticalAlign = "faixa", "right", "center"
		group.Shapes = append(group.Shapes, title)
	}
	// name seleciona a faixa horizontal específica do jogo no mapa original de PnlMenu.
	if name := menuBannerName(node); name != "" {
		// state fornece o PNG completo, com marca e cor já desenhadas.
		state, ok := skin.MustBuiltin().StateByType("label", name, "default")
		if ok {
			group.Shapes = append(group.Shapes, Primitive{Name: "faixa-jogo", SkinAsset: state.ImagePath, Bounds: menuScaled(component.Bounds, screen.Bounds{Y: 150, Width: 800, Height: 27})})
		}
	}
	if node.Layout == catalog.MenuLayoutBolao {
		// item percorre os três botões com legenda da composição digital.
		for index, item := range node.Items {
			if index >= 3 {
				break
			}
			bounds := menuSubmenuBounds(component.Bounds, node, item, index)
			group.Shapes = append(group.Shapes, menuItemAt(bounds, item, component)...)
			// state conserva fonte física e cor do label de legenda do Java.
			state, _ := skin.MustBuiltin().StateByType("label", "label_texto_azul", "default")
			// legend usa a caixa própria de 200×65 logo abaixo do botão.
			legend := styledText(component, state, item.Legend, menuScaled(component.Bounds, screen.Bounds{X: 50 + index*250, Y: 370, Width: 200, Height: 65}))
			legend.Name, legend.Wrap, legend.TextAlign = "legenda", true, "center"
			group.Shapes = append(group.Shapes, legend)
		}
		return group
	}
	// item percorre a grade do nó visível com as mesmas caixas usadas pelo clique.
	for index, item := range node.Items {
		// bounds deriva somente geometria; a identidade e os metadados continuam intactos.
		bounds := menuSubmenuBounds(component.Bounds, node, item, index)
		group.Shapes = append(group.Shapes, menuItemAt(bounds, item, component)...)
	}
	return group
}

// menuRootServiceBounds conserva a caixa dos serviços na raiz; jogos e submenus têm regras próprias.
func menuRootServiceBounds(bounds screen.Bounds, index int) screen.Bounds {
	return menuScaled(bounds, screen.Bounds{X: 20, Y: 150 + (index/2)*73, Width: 164, Height: 68})
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
			area := menuRootServiceBounds(component.Bounds, index)
			area.X += (index % 2) * 180 * component.Bounds.Width / MenuBaseWidth
			if area.X+area.Width > component.Bounds.X+400*component.Bounds.Width/MenuBaseWidth {
				continue
			}
			areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: area})
		}
		for index, item := range games {
			areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: screen.Bounds{
				X:      component.Bounds.X + (420+(index%2)*180)*component.Bounds.Width/MenuBaseWidth,
				Y:      component.Bounds.Y + (150+(index/2)*77)*component.Bounds.Height/MenuBaseHeight,
				Width:  164 * component.Bounds.Width / MenuBaseWidth,
				Height: 52 * component.Bounds.Height / MenuBaseHeight,
			}})
		}
		return areas
	}
	// item percorre os itens do submenu visível.
	for index, item := range node.Items {
		if node.Layout == catalog.MenuLayoutBolao && index >= 3 {
			break
		}
		// area é exatamente a caixa usada na pintura, inclusive nos jogos normal/especial.
		area := menuSubmenuBounds(component.Bounds, node, item, index)
		areas = append(areas, MenuHitArea{ComponentID: component.ID, ItemID: item.ID, Bounds: area})
	}
	return areas
}

// menuItemSkin resolve a chave XML; a variante do item prevalece sobre o estilo do composto.
func menuItemSkin(item catalog.MenuItem, parent screen.Component) string {
	if item.Variant != "" {
		return item.Variant
	}
	if parent.Style.Skin != "" {
		return parent.Style.Skin
	}
	if item.Category == "JOGOS" {
		if name := catalog.MenuGameSkin(item.Label); name != "" {
			return name
		}
	}
	// name mantém a seta de ramo para itens sem marca conhecida.
	name := "botao_menu"
	if item.Submenu != nil {
		name = "botao_menu_seta"
	}
	return name
}

// menuItemAt posiciona a cena do item no menu sem alterar os dados confirmados.
func menuItemAt(bounds screen.Bounds, item catalog.MenuItem, parent screen.Component) []Primitive {
	parent.Bounds = bounds
	return ApplicationMenuItem(parent, item, "default").Shapes
}

// ApplicationMenuItem compartilha a imagem normal, touched ou disabled entre os três consumidores.
// Bounds é a região completa clicável; legenda autorada diferente da marca ganha caixa própria.
func ApplicationMenuItem(parent screen.Component, item catalog.MenuItem, stateName string) Group {
	// bounds conserva a geometria do menu, independentemente da quantidade de primitivas.
	bounds := parent.Bounds
	// name seleciona a definição completa, nunca concatena nomes de arquivos por estado.
	name := menuItemSkin(item, parent)
	if !item.Enabled || !parent.Enabled {
		stateName = "disabled"
	}
	// state usa somente definições que possuem todos os recursos incorporados.
	state, _ := skin.MustBuiltin().State(catalog.Button, name, stateName)
	// component permite compartilhar o resolvedor de tipografia dos demais botões.
	component := screen.Component{Kind: catalog.Button, Bounds: bounds, Enabled: item.Enabled && parent.Enabled, Style: parent.Style}
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
	// group permite preservar texto autorado sem imprimi-lo sobre o logotipo original.
	group := Group{ComponentID: parent.ID, Kind: catalog.ApplicationMenu, Bounds: bounds, Shapes: []Primitive{result}}
	if _, branded := catalog.MenuGameLabel(name); branded && result.SkinAsset != "" {
		group.Shapes[0].Text = ""
		if item.Label != "" && !catalog.MenuGameCaptionMatches(item.Label, name) {
			// caption conserva tipografia editável em faixa separada, dentro da área clicável.
			caption := result
			caption.SkinAsset, caption.Fill, caption.Stroke, caption.Name = "", "", "", item.ID+"/legenda"
			caption.Bounds.Height = min(bounds.Height/2, max(1, int(result.FontSize)+6))
			caption.Bounds.Y += bounds.Height - caption.Bounds.Height
			group.Shapes[0].Bounds.Height -= caption.Bounds.Height
			group.Shapes = append(group.Shapes, caption)
		}
	}
	return group
}
