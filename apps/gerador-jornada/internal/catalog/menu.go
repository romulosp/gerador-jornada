package catalog

import (
	"fmt"

	"gerador-jornada.local/editor/internal/project"
)

// Layouts aceitos pela árvore de menus.
const (
	// MenuLayoutRoot identifica a raiz com serviços à esquerda e jogos à direita.
	MenuLayoutRoot = "root"
	// MenuLayoutStandard identifica uma grade comum de itens.
	MenuLayoutStandard = "standard"
	// MenuLayoutFinancial identifica a grade financeira de botões maiores.
	MenuLayoutFinancial = "financial"
	// MenuLayoutBolao identifica a composição digital de três botões com legenda.
	MenuLayoutBolao = "bolao"
)

// Limites estruturais da árvore de menus.
const (
	// MenuMaxDepth limita a profundidade incluindo a raiz como nível um.
	MenuMaxDepth = 16
	// MenuMaxItems limita o total de itens somando todos os descendentes.
	MenuMaxItems = 1000
)

// MenuItem representa uma folha ou um ramo da árvore de menus.
type MenuItem struct {
	// ID identifica o item de forma estável dentro do documento.
	ID string
	// Label contém o rótulo do botão, inclusive Unicode e vazio.
	Label string
	// Legend contém a legenda ilustrativa literal, independente dos metadados.
	Legend string
	// Category contém a categoria editável, inclusive vazia.
	Category string
	// Variant referencia um estilo suportado pelo catálogo, sem código executável.
	Variant string
	// Enabled controla a ativação explícita do item.
	Enabled bool
	// Submenu contém o ramo descendente; nil representa uma folha.
	Submenu *SubmenuProperties
	// TargetScreenID referencia uma tela de destino autorada, quando folha.
	TargetScreenID string
	// TargetJourney referencia somente um arquivo .jornada relativo à pasta do workspace.
	TargetJourney string
	// MovementCode contém metadado legado inerte.
	MovementCode string
	// MovementOption contém metadado legado inerte.
	MovementOption string
	// MovementName contém metadado legado inerte.
	MovementName string
	// Shortcut contém atalho descritivo inerte, nunca uma rota implícita.
	Shortcut string
}

// SubmenuProperties representa um nó ordenado da árvore.
type SubmenuProperties struct {
	// ID identifica o nó de forma estável dentro do documento.
	ID string
	// Name contém o nome editável do nó.
	Name string
	// Title contém o título exibido na faixa; vazio na raiz.
	Title string
	// Layout seleciona root, standard, financial ou bolao.
	Layout string
	// Items mantém a ordem visual; lista vazia é válida.
	Items []MenuItem
}

// ApplicationMenuProperties representa a árvore persistida e o submenu inicial.
type ApplicationMenuProperties struct {
	// Root contém o nó raiz por valor.
	Root SubmenuProperties
	// InitialMenuID referencia o UUID do submenu exibido por padrão.
	InitialMenuID string
}

// NewApplicationMenuProperties cria o menu padrão com a árvore de serviços e jogos.
func NewApplicationMenuProperties(supplyID func() string) *ApplicationMenuProperties {
	return seedApplicationMenu(supplyID)
}

// Validate verifica layouts, identidades, limites, ciclos e a escolha inicial.
func (p ApplicationMenuProperties) Validate() error {
	// seen detecta UUIDs repetidos entre nós e itens.
	seen := map[string]bool{}
	// nodeIDs reúne todos os submenus para validar a escolha inicial.
	nodeIDs := map[string]bool{}
	// active detecta um ciclo de ponteiros em memória antes da recursão.
	active := map[*SubmenuProperties]bool{}
	// count acumula itens inclusive desabilitados.
	count := 0
	// walk percorre a árvore com profundidade limitada.
	var walk func(node *SubmenuProperties, depth int, root bool) error
	walk = func(node *SubmenuProperties, depth int, root bool) error {
		if node == nil {
			return fmt.Errorf("submenu ausente")
		}
		if active[node] {
			return fmt.Errorf("ciclo de submenu")
		}
		if depth > MenuMaxDepth {
			return fmt.Errorf("profundidade de menu excede %d", MenuMaxDepth)
		}
		if !project.ValidID(node.ID) || seen[node.ID] {
			return fmt.Errorf("submenu sem identidade única")
		}
		seen[node.ID] = true
		nodeIDs[node.ID] = true
		allowed := node.Layout == MenuLayoutStandard || node.Layout == MenuLayoutFinancial || node.Layout == MenuLayoutBolao
		if root {
			allowed = node.Layout == MenuLayoutRoot
		}
		if !allowed {
			return fmt.Errorf("layout de menu inválido: %s", node.Layout)
		}
		active[node] = true
		// item percorre a ordem preservada do nó.
		for _, item := range node.Items {
			count++
			if count > MenuMaxItems {
				return fmt.Errorf("menu excede %d itens", MenuMaxItems)
			}
			if !project.ValidID(item.ID) || seen[item.ID] {
				return fmt.Errorf("item de menu sem identidade única")
			}
			seen[item.ID] = true
			// destinations conta alternativas mutuamente exclusivas, sem resolver arquivos no domínio.
			destinations := 0
			if item.Submenu != nil {
				destinations++
			}
			if item.TargetScreenID != "" {
				destinations++
			}
			if item.TargetJourney != "" {
				destinations++
			}
			if destinations > 1 {
				return fmt.Errorf("item de menu combina submenu e destino")
			}
			if item.TargetJourney != "" && !ValidJourneyFilename(item.TargetJourney) {
				return fmt.Errorf("destino de jornada inválido: %s", item.TargetJourney)
			}
			if item.Submenu != nil {
				if err := walk(item.Submenu, depth+1, false); err != nil {
					return err
				}
			}
		}
		delete(active, node)
		return nil
	}
	if err := walk(&p.Root, 1, true); err != nil {
		return err
	}
	if !nodeIDs[p.InitialMenuID] {
		return fmt.Errorf("menu inicial inexistente")
	}
	return nil
}

// FindSubmenu localiza um nó pelo identificador estável, inclusive a raiz.
func (p ApplicationMenuProperties) FindSubmenu(id string) *SubmenuProperties {
	// visited impede recursão infinita quando a árvore recebida ainda não foi validada.
	visited := map[*SubmenuProperties]bool{}
	// find percorre a árvore sem depender de títulos repetidos.
	var find func(node *SubmenuProperties) *SubmenuProperties
	find = func(node *SubmenuProperties) *SubmenuProperties {
		if node == nil || visited[node] {
			return nil
		}
		visited[node] = true
		if node.ID == id {
			return node
		}
		// item percorre os ramos descendentes.
		for _, item := range node.Items {
			if found := find(item.Submenu); found != nil {
				return found
			}
		}
		return nil
	}
	return find(&p.Root)
}

// Initial devolve o submenu escolhido como inicial.
func (p ApplicationMenuProperties) Initial() *SubmenuProperties {
	return p.FindSubmenu(p.InitialMenuID)
}

// PathTo devolve a cadeia de nós da raiz até o identificador solicitado.
func (p ApplicationMenuProperties) PathTo(id string) []*SubmenuProperties {
	// path acumula os nós visitados durante a busca em profundidade.
	path := []*SubmenuProperties{}
	// visited impede recursão infinita em árvores ainda não validadas.
	visited := map[*SubmenuProperties]bool{}
	// find devolve a cadeia encontrada ou nil.
	var find func(node *SubmenuProperties) []*SubmenuProperties
	find = func(node *SubmenuProperties) []*SubmenuProperties {
		if node == nil || visited[node] {
			return nil
		}
		visited[node] = true
		path = append(path, node)
		if node.ID == id {
			return append([]*SubmenuProperties{}, path...)
		}
		// item percorre os ramos descendentes.
		for _, item := range node.Items {
			if found := find(item.Submenu); found != nil {
				return found
			}
		}
		path = path[:len(path)-1]
		return nil
	}
	return find(&p.Root)
}

// WalkItems percorre todos os itens na ordem da árvore, inclusive ramos desabilitados.
// A visita recebe um ponteiro somente para leitura; a função interrompe em erro ou ciclo.
func (p ApplicationMenuProperties) WalkItems(visit func(*MenuItem) error) error {
	// visited impede recursão infinita em uma árvore ainda não validada.
	visited := map[*SubmenuProperties]bool{}
	// count limita a travessia ao mesmo teto da validação.
	count := 0
	// walk percorre um nó e seus descendentes.
	var walk func(node *SubmenuProperties, depth int) error
	walk = func(node *SubmenuProperties, depth int) error {
		if node == nil {
			return nil
		}
		if visited[node] {
			return fmt.Errorf("ciclo de submenu")
		}
		if depth > MenuMaxDepth {
			return fmt.Errorf("profundidade de menu excede %d", MenuMaxDepth)
		}
		visited[node] = true
		// index percorre os itens do nó.
		for index := range node.Items {
			count++
			if count > MenuMaxItems {
				return fmt.Errorf("menu excede %d itens", MenuMaxItems)
			}
			if err := visit(&node.Items[index]); err != nil {
				return err
			}
			if err := walk(node.Items[index].Submenu, depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(&p.Root, 1)
}

// FindMenuItem localiza um item pelo identificador em toda a árvore.
func (p ApplicationMenuProperties) FindMenuItem(id string) *MenuItem {
	// found preserva o item localizado durante a travessia.
	var found *MenuItem
	_ = p.WalkItems(func(item *MenuItem) error {
		if item.ID == id {
			found = item
		}
		return nil
	})
	return found
}

// CollectMenuIDs devolve todas as identidades internas para validação global.
func CollectMenuIDs(p ApplicationMenuProperties) []string {
	// ids preserva a ordem de nós e itens.
	ids := []string{}
	// visited impede recursão infinita em árvores ainda não validadas.
	visited := map[*SubmenuProperties]bool{}
	// walk percorre a árvore confirmada.
	var walk func(node *SubmenuProperties)
	walk = func(node *SubmenuProperties) {
		if node == nil || visited[node] {
			return
		}
		visited[node] = true
		ids = append(ids, node.ID)
		// item percorre os itens do nó.
		for _, item := range node.Items {
			ids = append(ids, item.ID)
			walk(item.Submenu)
		}
	}
	walk(&p.Root)
	return ids
}

// DuplicateApplicationMenu copia a árvore com identidades novas e remapeia a escolha inicial.
// supplyID entrega UUIDs novos sem compartilhar slices ou ponteiros com a origem.
func DuplicateApplicationMenu(source ApplicationMenuProperties, supplyID func() string) ApplicationMenuProperties {
	// mapping associa cada identidade antiga à nova para remapear a escolha inicial.
	mapping := map[string]string{}
	// clone copia um nó e seus descendentes com identidades novas.
	var clone func(node *SubmenuProperties) *SubmenuProperties
	clone = func(node *SubmenuProperties) *SubmenuProperties {
		if node == nil {
			return nil
		}
		// copy preserva textos e layout, separando os itens mutáveis.
		copy := *node
		copy.ID = supplyID()
		mapping[node.ID] = copy.ID
		copy.Items = []MenuItem{}
		// item percorre a ordem original.
		for _, item := range node.Items {
			// itemCopy separa identidade, ponteiro e destino autorreferente quando necessário.
			itemCopy := item
			itemCopy.ID = supplyID()
			mapping[item.ID] = itemCopy.ID
			itemCopy.Submenu = clone(item.Submenu)
			copy.Items = append(copy.Items, itemCopy)
		}
		return &copy
	}
	// result recebe a raiz clonada.
	result := ApplicationMenuProperties{}
	root := clone(&source.Root)
	result.Root = *root
	// initial resolve a nova identidade; raiz permanece o padrão quando não houver escolha.
	result.InitialMenuID = mapping[source.InitialMenuID]
	if result.InitialMenuID == "" {
		result.InitialMenuID = result.Root.ID
	}
	return result
}

// RemapMenuTargets troca identidades de tela de destino conforme o mapa de uma cópia de tela.
// Destinos externos são preservados; referências à própria tela copiada usam o novo identificador.
func RemapMenuTargets(menu *ApplicationMenuProperties, screenMapping map[string]string) {
	if menu == nil {
		return
	}
	// walk percorre os itens para remapear somente destinos presentes no mapa.
	var walk func(node *SubmenuProperties)
	walk = func(node *SubmenuProperties) {
		if node == nil {
			return
		}
		// index percorre os itens mutáveis do nó.
		for index := range node.Items {
			if target, ok := screenMapping[node.Items[index].TargetScreenID]; ok {
				node.Items[index].TargetScreenID = target
			}
			walk(node.Items[index].Submenu)
		}
	}
	walk(&menu.Root)
}
