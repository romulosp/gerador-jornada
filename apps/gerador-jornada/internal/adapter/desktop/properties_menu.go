package desktop

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// menuTreeEntry descreve um nó ou item selecionável do editor de árvore.
type menuTreeEntry struct {
	// Key identifica a entrada por identidade estável, nunca por título.
	Key string
	// Label apresenta caminho e ordem para desambiguar títulos repetidos.
	Label string
	// NodeID contém o UUID do submenu, quando a entrada é um nó.
	NodeID string
	// ItemID contém o UUID do item, quando a entrada é um item.
	ItemID string
	// ParentNodeID identifica o submenu que contém o item.
	ParentNodeID string
}

// menuVisibleText devolve o texto literal ou um marcador explícito de ausência.
func menuVisibleText(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// menuJoinPath acrescenta um segmento sem duplicar a barra da raiz.
func menuJoinPath(parent, segment string) string {
	if parent == "" || parent == "/" {
		return "/" + segment
	}
	return parent + "/" + segment
}

// menuTreeEntries achata a árvore em caminhos, sem criar um widget por folha.
func menuTreeEntries(menu *catalog.ApplicationMenuProperties) []menuTreeEntry {
	// entries preserva a ordem de profundidade usada pelo seletor e pelo menu inicial.
	entries := []menuTreeEntry{}
	if menu == nil {
		return entries
	}
	// visited impede recursão infinita em uma árvore ainda não validada.
	visited := map[*catalog.SubmenuProperties]bool{}
	// ordinal diferencia nós e itens que compartilham o mesmo caminho visível.
	ordinal := 0
	// walk percorre nós e itens em profundidade sem renderizar a árvore inteira.
	var walk func(node *catalog.SubmenuProperties, path string)
	walk = func(node *catalog.SubmenuProperties, path string) {
		if node == nil || visited[node] {
			return
		}
		visited[node] = true
		ordinal++
		entries = append(entries, menuTreeEntry{Key: "node:" + node.ID, Label: fmt.Sprintf("Nó %d · %s", ordinal, path), NodeID: node.ID})
		// item representa uma folha ou um ramo do nó atual.
		for _, item := range node.Items {
			ordinal++
			// itemPath contém o caminho navegável até a legenda do item.
			itemPath := menuJoinPath(path, menuVisibleText(item.Label, "(sem legenda)"))
			entries = append(entries, menuTreeEntry{Key: "item:" + item.ID, Label: fmt.Sprintf("Item %d · %s", ordinal, itemPath), ItemID: item.ID, ParentNodeID: node.ID})
			walk(item.Submenu, itemPath)
		}
	}
	walk(&menu.Root, "/")
	return entries
}

// menuFindItem localiza um item pelo UUID e devolve o nó pai e sua posição.
func menuFindItem(root *catalog.SubmenuProperties, id string) (*catalog.SubmenuProperties, int) {
	// visited impede recursão infinita em uma árvore ainda não validada.
	visited := map[*catalog.SubmenuProperties]bool{}
	// foundNode e foundIndex preservam o resultado da busca.
	var foundNode *catalog.SubmenuProperties
	foundIndex := -1
	// walk percorre os nós em profundidade, primeiro os irmãos e depois os ramos.
	var walk func(node *catalog.SubmenuProperties)
	walk = func(node *catalog.SubmenuProperties) {
		if node == nil || visited[node] || foundNode != nil {
			return
		}
		visited[node] = true
		// index percorre os itens do nó atual.
		for index := range node.Items {
			if node.Items[index].ID == id {
				foundNode, foundIndex = node, index
				return
			}
		}
		// index percorre novamente para descer somente após esgotar os irmãos.
		for index := range node.Items {
			walk(node.Items[index].Submenu)
			if foundNode != nil {
				return
			}
		}
	}
	walk(root)
	return foundNode, foundIndex
}

// menuFindSubmenuOwner localiza o item cujo Submenu possui a identidade informada.
func menuFindSubmenuOwner(root *catalog.SubmenuProperties, id string) (*catalog.SubmenuProperties, int) {
	// visited impede recursão infinita em uma árvore ainda não validada.
	visited := map[*catalog.SubmenuProperties]bool{}
	// foundNode e foundIndex preservam o resultado da busca.
	var foundNode *catalog.SubmenuProperties
	foundIndex := -1
	// walk percorre os nós em profundidade.
	var walk func(node *catalog.SubmenuProperties)
	walk = func(node *catalog.SubmenuProperties) {
		if node == nil || visited[node] || foundNode != nil {
			return
		}
		visited[node] = true
		// index percorre os itens do nó atual.
		for index := range node.Items {
			if child := node.Items[index].Submenu; child != nil && child.ID == id {
				foundNode, foundIndex = node, index
				return
			}
		}
		// index percorre novamente para buscar em ramos descendentes.
		for index := range node.Items {
			walk(node.Items[index].Submenu)
			if foundNode != nil {
				return
			}
		}
	}
	walk(root)
	return foundNode, foundIndex
}

// menuFields monta o editor de rascunho da árvore e do menu inicial.
func (e *Editor) menuFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// draft preserva uma cópia isolada da árvore confirmada até a confirmação explícita.
	draft := catalog.ApplicationMenuProperties{}
	if component.Properties.Menu != nil {
		// copied contém a árvore separada da origem, sem slices ou ponteiros compartilhados.
		copied := component.Properties.Copy().Menu
		draft = *copied
	} else {
		draft = *catalog.NewApplicationMenuProperties(project.NewID)
	}
	// initialID conserva a escolha do menu inicial apenas no rascunho.
	initialID := draft.InitialMenuID
	if draft.FindSubmenu(initialID) == nil {
		initialID = draft.Root.ID
	}
	// selectedKey identifica o nó ou item atual, sem persistir caminho ou foco.
	selectedKey := "node:" + draft.Root.ID
	// syncing impede que atualizações programáticas disparem callbacks de seleção.
	syncing := false
	// message informa recusas locais e mantém o rascunho sem publicar.
	message := widget.NewLabel("")
	message.Wrapping = fyne.TextWrapWord
	// nodeSelect lista nós e itens por caminho, sem criar um controle por folha.
	nodeSelect := newSelect(nil, nil)
	// nameEntry, titleEntry e layoutSelect editam o nó selecionado.
	nameEntry := entry("")
	titleEntry := entry("")
	layoutSelect := newSelect([]string{catalog.MenuLayoutRoot}, nil)
	// labelEntry, legendEntry, categoryEntry, variantEntry e enabledCheck editam o item selecionado.
	labelEntry := entry("")
	legendEntry := widget.NewMultiLineEntry()
	legendEntry.SetMinRowsVisible(2)
	categoryEntry := entry("")
	variantEntry := entry("")
	enabledCheck := widget.NewCheck("", nil)
	// targetSelect vincula uma folha a uma tela existente ou à ausência de destino.
	targetSelect := newSelect(nil, nil)
	// submenuLabel informa o ramo existente ou a ausência de submenu.
	submenuLabel := widget.NewLabel("")
	// initialSelect escolhe o menu inicial por caminho, sem depender de título.
	initialSelect := newSelect([]string{}, nil)
	// targetIDs associa o texto visível ao UUID da tela, incluindo a ausência de destino.
	targetIDs := map[string]string{"Sem destino": ""}
	// targetLabelByID associa o UUID da tela ao texto visível.
	targetLabelByID := map[string]string{"": "Sem destino"}
	// targetOptions preserva a ordem das telas existentes.
	targetOptions := []string{"Sem destino"}
	// index percorre as telas existentes para montar o seletor de destino.
	for index, page := range e.Session.Document.Screens {
		// label diferencia telas com nomes repetidos por posição.
		label := fmt.Sprintf("%d · %s", index+1, page.Name)
		targetOptions = append(targetOptions, label)
		targetIDs[label] = page.ID
		targetLabelByID[page.ID] = label
	}
	targetSelect.Options = targetOptions
	// optionKeys e initialKeys associam o texto visível à identidade estável.
	optionKeys := map[string]string{}
	initialKeys := map[string]string{}
	// selectedNode devolve o nó de contexto: o próprio nó ou o pai do item selecionado.
	selectedNode := func() *catalog.SubmenuProperties {
		if nodeID, ok := strings.CutPrefix(selectedKey, "node:"); ok {
			return draft.FindSubmenu(nodeID)
		}
		if itemID, ok := strings.CutPrefix(selectedKey, "item:"); ok {
			// parent contém o nó que hospeda o item selecionado.
			parent, _ := menuFindItem(&draft.Root, itemID)
			return parent
		}
		return nil
	}
	// selectedItem devolve o nó pai e a posição do item selecionado, quando houver.
	selectedItem := func() (*catalog.SubmenuProperties, int) {
		if itemID, ok := strings.CutPrefix(selectedKey, "item:"); ok {
			return menuFindItem(&draft.Root, itemID)
		}
		return nil, -1
	}
	// flush copia os campos visíveis para o rascunho antes de trocar de seleção ou confirmar.
	flush := func() {
		if nodeID, ok := strings.CutPrefix(selectedKey, "node:"); ok {
			// target contém o nó selecionado, localizado pela identidade.
			target := draft.FindSubmenu(nodeID)
			if target == nil {
				return
			}
			target.Name = nameEntry.Text
			target.Title = titleEntry.Text
			if target.ID == draft.Root.ID {
				target.Layout = catalog.MenuLayoutRoot
			} else if slices.Contains(layoutSelect.Options, layoutSelect.Selected) {
				target.Layout = layoutSelect.Selected
			}
			return
		}
		// parent e index localizam o item selecionado no rascunho.
		parent, index := selectedItem()
		if parent == nil || index < 0 || index >= len(parent.Items) {
			return
		}
		// item contém o item mutável do rascunho.
		item := &parent.Items[index]
		item.Label = labelEntry.Text
		item.Legend = legendEntry.Text
		item.Category = categoryEntry.Text
		item.Variant = variantEntry.Text
		item.Enabled = enabledCheck.Checked
		if item.Submenu != nil {
			item.TargetScreenID = ""
		} else if targetID, ok := targetIDs[targetSelect.Selected]; ok {
			item.TargetScreenID = targetID
		}
	}
	// syncEditor reconstrói as opções e recarrega o nó ou item selecionado.
	syncEditor := func() {
		syncing = true
		// entries contém a árvore achatada em caminhos.
		entries := menuTreeEntries(&draft)
		// labels e labelByKey preservam a ordem e associam texto visível à identidade.
		labels := []string{}
		labelByKey := map[string]string{}
		// initialLabels e initialLabelByID limitam o menu inicial a submenus existentes.
		initialLabels := []string{}
		initialLabelByID := map[string]string{}
		// optionKeys e initialKeys renovam os mapas de seleção a cada reconstrução.
		optionKeys = map[string]string{}
		initialKeys = map[string]string{}
		// entry representa um nó ou item selecionável da árvore achatada.
		for _, entry := range entries {
			labels = append(labels, entry.Label)
			optionKeys[entry.Label] = entry.Key
			labelByKey[entry.Key] = entry.Label
			if entry.NodeID != "" {
				initialLabels = append(initialLabels, entry.Label)
				initialKeys[entry.Label] = entry.NodeID
				initialLabelByID[entry.NodeID] = entry.Label
			}
		}
		if labelByKey[selectedKey] == "" {
			selectedKey = "node:" + draft.Root.ID
		}
		nodeSelect.Options = labels
		nodeSelect.ClearSelected()
		if label := labelByKey[selectedKey]; label != "" {
			nodeSelect.SetSelected(label)
		}
		nodeSelect.Refresh()
		initialSelect.Options = initialLabels
		initialSelect.ClearSelected()
		if label := initialLabelByID[initialID]; label != "" {
			initialSelect.SetSelected(label)
		} else {
			initialID = ""
			message.SetText("O menu inicial foi removido; escolha outro antes de salvar.")
		}
		initialSelect.Refresh()
		// node contém o nó em edição: seleção direta ou pai do item selecionado.
		node := selectedNode()
		if node == nil {
			syncing = false
			return
		}
		// root informa se o nó atual é a raiz, que aceita somente o layout root.
		root := node.ID == draft.Root.ID
		nameEntry.SetText(node.Name)
		titleEntry.SetText(node.Title)
		// layoutOptions restringe a raiz a root e os demais nós aos três layouts de grade.
		layoutOptions := []string{catalog.MenuLayoutStandard, catalog.MenuLayoutFinancial, catalog.MenuLayoutBolao}
		if root {
			layoutOptions = []string{catalog.MenuLayoutRoot}
		}
		layoutSelect.Options = layoutOptions
		layoutSelect.ClearSelected()
		if slices.Contains(layoutOptions, node.Layout) {
			layoutSelect.SetSelected(node.Layout)
		}
		layoutSelect.Refresh()
		if nodeKey, ok := strings.CutPrefix(selectedKey, "node:"); ok && nodeKey == node.ID {
			nameEntry.Enable()
			titleEntry.Enable()
			layoutSelect.Enable()
			labelEntry.Disable()
			legendEntry.Disable()
			categoryEntry.Disable()
			variantEntry.Disable()
			enabledCheck.Disable()
			targetSelect.Disable()
			submenuLabel.SetText("Selecione um item para editar legenda, destino e submenu.")
			syncing = false
			return
		}
		// parent e index localizam o item selecionado para carregar seus campos.
		parent, index := selectedItem()
		if parent == nil || index < 0 || index >= len(parent.Items) {
			syncing = false
			return
		}
		// item contém a folha ou ramo selecionado.
		item := parent.Items[index]
		nameEntry.Disable()
		titleEntry.Disable()
		layoutSelect.Disable()
		labelEntry.Enable()
		labelEntry.SetText(item.Label)
		legendEntry.Enable()
		legendEntry.SetText(item.Legend)
		categoryEntry.Enable()
		categoryEntry.SetText(item.Category)
		variantEntry.Enable()
		variantEntry.SetText(item.Variant)
		enabledCheck.Enable()
		enabledCheck.SetChecked(item.Enabled)
		targetSelect.ClearSelected()
		if item.Submenu != nil {
			submenuLabel.SetText("Submenu: " + menuVisibleText(item.Submenu.Name, "(sem nome)"))
			targetSelect.Disable()
		} else {
			submenuLabel.SetText("Folha sem submenu.")
			targetSelect.Enable()
			if label := targetLabelByID[item.TargetScreenID]; label != "" {
				targetSelect.SetSelected(label)
			}
		}
		syncing = false
	}
	// addItem acrescenta um item habilitado ao nó selecionado ou ao pai do item selecionado.
	addItem := func() {
		flush()
		// parent contém o nó que receberá o novo item.
		parent := selectedNode()
		if parent == nil {
			message.SetText("Selecione um nó ou item.")
			return
		}
		// item nasce habilitado e com identidade nova, coerente com o seed do domínio.
		item := catalog.MenuItem{ID: project.NewID(), Enabled: true}
		parent.Items = append(parent.Items, item)
		selectedKey = "item:" + item.ID
		message.SetText("")
		syncEditor()
	}
	// removeItem remove o item selecionado e seleciona seu nó pai.
	removeItem := func() {
		// parent e index localizam o item selecionado antes do descarte.
		parent, index := selectedItem()
		if parent == nil || index < 0 {
			message.SetText("Selecione um item para remover.")
			return
		}
		flush()
		parent.Items = append(parent.Items[:index], parent.Items[index+1:]...)
		selectedKey = "node:" + parent.ID
		message.SetText("")
		syncEditor()
	}
	// addSubmenu anexa um ramo novo ao item selecionado e remove seu destino.
	addSubmenu := func() {
		// parent e index localizam o item que receberá o ramo.
		parent, index := selectedItem()
		if parent == nil || index < 0 {
			message.SetText("Selecione um item para criar o submenu.")
			return
		}
		flush()
		// item contém o item selecionado após a atualização dos campos.
		item := &parent.Items[index]
		if item.Submenu != nil {
			message.SetText("O item já possui submenu.")
			return
		}
		// created contém o ramo novo com identidade própria e layout padrão.
		created := &catalog.SubmenuProperties{ID: project.NewID(), Name: "Novo submenu", Title: "", Layout: catalog.MenuLayoutStandard, Items: []catalog.MenuItem{}}
		item.Submenu = created
		item.TargetScreenID = ""
		selectedKey = "node:" + created.ID
		message.SetText("O destino foi removido para criar o submenu.")
		syncEditor()
	}
	// removeSubmenu desanexa o ramo do item selecionado ou do nó selecionado não raiz.
	removeSubmenu := func() {
		if parent, index := selectedItem(); parent != nil && index >= 0 {
			if parent.Items[index].Submenu == nil {
				message.SetText("O item selecionado não possui submenu.")
				return
			}
			flush()
			parent.Items[index].Submenu = nil
			selectedKey = "item:" + parent.Items[index].ID
			message.SetText("")
			syncEditor()
			return
		}
		if nodeID, ok := strings.CutPrefix(selectedKey, "node:"); ok {
			if nodeID == draft.Root.ID {
				message.SetText("A raiz não pode ser removida.")
				return
			}
			// owner e ownerIndex localizam o item que contém o ramo selecionado.
			owner, ownerIndex := menuFindSubmenuOwner(&draft.Root, nodeID)
			if owner == nil || ownerIndex < 0 {
				message.SetText("Selecione um submenu para remover.")
				return
			}
			flush()
			owner.Items[ownerIndex].Submenu = nil
			selectedKey = "item:" + owner.Items[ownerIndex].ID
			message.SetText("")
			syncEditor()
			return
		}
		message.SetText("Selecione um submenu para remover.")
	}
	// moveItem troca o item selecionado com seu vizinho no mesmo nó.
	moveItem := func(delta int) {
		// parent e index localizam o item selecionado.
		parent, index := selectedItem()
		if parent == nil || index < 0 {
			message.SetText("Selecione um item para mudar a ordem.")
			return
		}
		// target contém a posição vizinha solicitada.
		target := index + delta
		if target < 0 || target >= len(parent.Items) {
			message.SetText("O item já está no limite.")
			return
		}
		flush()
		parent.Items[index], parent.Items[target] = parent.Items[target], parent.Items[index]
		message.SetText("")
		syncEditor()
	}
	nodeSelect.OnChanged = func(label string) {
		if syncing {
			return
		}
		flush()
		// key contém a identidade escolhida pelo caminho exibido.
		key := optionKeys[label]
		if key == "" {
			return
		}
		selectedKey = key
		message.SetText("")
		syncEditor()
	}
	initialSelect.OnChanged = func(label string) {
		if syncing {
			return
		}
		// id contém o UUID do submenu escolhido como inicial.
		id := initialKeys[label]
		if id == "" {
			return
		}
		initialID = id
		message.SetText("")
		syncEditor()
	}
	form.Append("Nó ou item", nodeSelect)
	form.Append("Nome", nameEntry)
	form.Append("Título", titleEntry)
	form.Append("Layout", layoutSelect)
	form.Append("Legenda", labelEntry)
	form.Append("Legenda explicativa", legendEntry)
	form.Append("Categoria", categoryEntry)
	form.Append("Variante", variantEntry)
	form.Append("Habilitado", enabledCheck)
	form.Append("Destino", targetSelect)
	form.Append("Submenu", submenuLabel)
	form.Append("Menu inicial", initialSelect)
	form.Append("Ações", container.NewGridWithColumns(2,
		widget.NewButton("Adicionar item", addItem),
		widget.NewButton("Remover item", removeItem),
		widget.NewButton("Adicionar submenu", addSubmenu),
		widget.NewButton("Remover submenu", removeSubmenu),
		widget.NewButton("Subir", func() { moveItem(-1) }),
		widget.NewButton("Descer", func() { moveItem(1) }),
	))
	form.Append("", message)
	syncEditor()
	return func(target *screen.Component) error {
		flush()
		if target.Properties.Menu == nil {
			return fmt.Errorf("componente sem propriedades de menu da aplicação")
		}
		// snapshot separa a árvore confirmada das listas ainda usadas pelo formulário.
		snapshot := catalog.Properties{Menu: &draft}.Copy().Menu
		snapshot.InitialMenuID = initialID
		if err := snapshot.Validate(); err != nil {
			return err
		}
		target.Properties.Menu.Root = snapshot.Root
		target.Properties.Menu.InitialMenuID = snapshot.InitialMenuID
		return target.Properties.Validate(target.Kind)
	}
}
