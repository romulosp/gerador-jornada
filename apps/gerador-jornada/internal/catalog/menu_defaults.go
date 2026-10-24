package catalog

// legacyMenuNode conserva os seis atributos XML sem interpretar códigos ou opções.
type legacyMenuNode struct {
	// Name conserva nome, acentos e espaços literais.
	Name string
	// Title conserva o título principal, usado quando o nó abre um submenu.
	Title string
	// MovementCode, MovementOption e MovementName são metadados inertes.
	MovementCode, MovementOption, MovementName string
	// Shortcut conserva o atalho descritivo, sem gerar comandos.
	Shortcut string
	// Children contém somente filhos ativos, na ordem original do XML.
	Children []legacyMenuNode
}

// seedApplicationMenu instancia a árvore importada com novos UUIDs, sem migrar documentos existentes.
func seedApplicationMenu(supplyID func() string) *ApplicationMenuProperties {
	// buildNode converte os dados gerados em propriedades tipadas e independentes.
	var buildNode func(legacyMenuNode, string, bool) SubmenuProperties
	buildNode = func(source legacyMenuNode, category string, root bool) SubmenuProperties {
		// node conserva título e identidade separados do item que abre o submenu.
		node := SubmenuProperties{ID: supplyID(), Name: source.Name, Title: source.Title, Layout: MenuLayoutStandard, Items: []MenuItem{}}
		if root {
			node.Layout, node.Title = MenuLayoutRoot, ""
		} else if source.Name == "OUTROS SERVIÇOS FINANCEIROS" {
			node.Layout = MenuLayoutFinancial
		} else if source.Name == "BOLÕES DIGITAIS" {
			node.Layout = MenuLayoutBolao
		}
		// child conserva ordem e literais, inclusive jogos especiais e metadados nos ramos.
		for index, child := range source.Children {
			// itemCategory separa apenas as colunas raiz; descendentes herdam a categoria local.
			itemCategory := category
			if root {
				itemCategory = "SERVIÇOS"
				if MenuGameSkin(child.Name) != "" {
					itemCategory = "JOGOS"
				}
			}
			// item preserva campos de negócio; Variant referencia somente apresentação.
			item := MenuItem{ID: supplyID(), Label: child.Name, Category: itemCategory, Variant: MenuGameSkin(child.Name), Enabled: true, MovementCode: child.MovementCode, MovementOption: child.MovementOption, MovementName: child.MovementName, Shortcut: child.Shortcut}
			if node.Layout == MenuLayoutFinancial {
				item.Variant = "botao_menu2"
				if len(child.Children) > 0 {
					item.Variant = "botao_menu2_seta"
				}
			}
			if child.Name == "BOLÕES DIGITAIS" {
				item.Variant = "botao_menu_peq"
			}
			if len(child.Children) > 0 {
				// submenu ganha UUID próprio, sem compartilhar ponteiros entre novas instâncias.
				submenu := buildNode(child, itemCategory, false)
				item.Submenu = &submenu
			}
			if node.Layout == MenuLayoutBolao && index < 3 {
				// legends são complementos visuais da referência Java, não substitutos dos rótulos XML.
				legends := []string{"Visualizar cotas digitais para venda na lotérica.", "Consulta bolões de concursos abertos para comercialização.", "Consulta bolões de concursos encerrados e não prescritos."}
				item.Legend = legends[index]
			}
			node.Items = append(node.Items, item)
		}
		return node
	}
	// root usa a ordem XML; a apresentação separa serviços e jogos sem reescrever a árvore.
	root := buildNode(legacyMenuSource, "SERVIÇOS", true)
	return &ApplicationMenuProperties{Root: root, InitialMenuID: root.ID}
}
