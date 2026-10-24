package presentation

import (
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// MenuCommand representa um comando local do rodapé, sem acrescentar itens à árvore persistida.
type MenuCommand struct {
	// ID identifica a intenção local: back, cancel ou exit.
	ID string
	// Label e Variant conservam texto e skin originais do comando.
	Label, Variant string
	// Bounds usa as mesmas unidades lógicas da cena e das regiões clicáveis.
	Bounds screen.Bounds
}

// menuScaled transforma uma caixa da referência 800×600 sem alterar bounds confirmados.
func menuScaled(parent, base screen.Bounds) screen.Bounds {
	return screen.Bounds{X: parent.X + base.X*parent.Width/MenuBaseWidth, Y: parent.Y + base.Y*parent.Height/MenuBaseHeight, Width: base.Width * parent.Width / MenuBaseWidth, Height: base.Height * parent.Height / MenuBaseHeight}
}

// MenuCommands devolve os comandos visíveis do PnlMenu: Sair na raiz, Cancelar/Voltar no submenu.
func MenuCommands(component screen.Component, node *catalog.SubmenuProperties) []MenuCommand {
	if node == nil {
		return nil
	}
	if node.Layout == catalog.MenuLayoutRoot {
		return []MenuCommand{{"exit", "SAIR", "botao_menu_peq_vermelho", menuScaled(component.Bounds, screen.Bounds{X: 675, Y: 533, Width: 80, Height: 46})}}
	}
	// commands são locais à prévia; no canvas/PPT conservam somente a aparência.
	commands := []MenuCommand{}
	if node.Layout != catalog.MenuLayoutFinancial {
		commands = append(commands, MenuCommand{"cancel", "CANCELAR", "botao_menu_peq_vermelho", menuScaled(component.Bounds, screen.Bounds{X: 20, Y: 530, Width: 80, Height: 46})})
	}
	return append(commands, MenuCommand{"back", "VOLTAR", "botao_menu_peq", menuScaled(component.Bounds, screen.Bounds{X: 700, Y: 530, Width: 80, Height: 46})})
}

// menuBannerName resolve somente decoração, conforme faixasJogos do Java; nunca cria rotas.
func menuBannerName(node *catalog.SubmenuProperties) string {
	if node == nil || node.Layout == catalog.MenuLayoutRoot {
		return ""
	}
	if strings.HasPrefix(node.Name, "BOLÃO ") && node.Name != "BOLÃO DUPLA-SENA" {
		return ""
	}
	// game conserva a chave exata; a definição XML determina o caminho da imagem horizontal.
	game := catalog.MenuGameSkin(node.Name)
	if game == "" {
		return ""
	}
	return "label_faixa_titulo_secundario_" + strings.TrimPrefix(strings.TrimSuffix(game, "_especial"), "botao_")
}

// menuSubmenuBounds reproduz a grade calculada pelo PnlMenu para cada nó visível.
func menuSubmenuBounds(parent screen.Bounds, node *catalog.SubmenuProperties, item catalog.MenuItem, index int) screen.Bounds {
	if node.Layout == catalog.MenuLayoutBolao {
		return menuScaled(parent, screen.Bounds{X: 50 + index*250, Y: 300, Width: 190, Height: 70})
	}
	if node.Layout == catalog.MenuLayoutFinancial {
		return menuScaled(parent, screen.Bounds{X: 35 + (index/4)*150, Y: 160 + (index%4)*90, Width: 140, Height: 80})
	}
	// columns e rows seguem a fórmula Java: mínimo de duas colunas e até quatro linhas por coluna.
	columns := max(2, (len(node.Items)+3)/4)
	rows := (len(node.Items) + columns - 1) / columns
	// Os quatro casos especiais são visuais do legado, sem efeito sobre identidade ou navegação.
	switch node.Name {
	case "BENEFÍCIOS SOCIAIS":
		columns, rows = 3, 3
	case "APOIO":
		columns, rows = 4, 3
	case "CAIXA TEM":
		columns, rows = 2, 1
	case "PRODUTOS NEGOCIAIS":
		columns, rows = 4, 1
	}
	if item.Label == "BOLÕES DIGITAIS" && item.Submenu != nil && item.Submenu.Layout == catalog.MenuLayoutBolao {
		return menuScaled(parent, screen.Bounds{X: 500, Y: 528, Width: 180, Height: 50})
	}
	// height distingue imagens dos jogos de botões de serviço, sem esticar o ícone.
	height := 68
	if _, branded := catalog.MenuGameLabel(menuItemSkin(item, screen.Component{})); branded {
		height = 52
	}
	return menuScaled(parent, screen.Bounds{X: (800-columns*180)/2 + (index%columns)*180, Y: (700-rows*80)/2 + (index/columns)*80, Width: 164, Height: height})
}
