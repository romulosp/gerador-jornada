package desktop

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// barcodeLabelChoices devolve as cinco opções visíveis e sua relação com o enum do domínio.
func barcodeLabelChoices() ([]string, map[string]string, map[string]string) {
	// choices preserva a ordem fixa exigida pela ficha processada.
	choices := []string{"CPF", "CNPJ", "Nome", "Nome fantasia", "Razão social"}
	// kinds associa a legenda visível ao discriminador persistido.
	kinds := map[string]string{
		"CPF":           catalog.BarcodeLabelCPF,
		"CNPJ":          catalog.BarcodeLabelCNPJ,
		"Nome":          catalog.BarcodeLabelName,
		"Nome fantasia": catalog.BarcodeLabelTradeName,
		"Razão social":  catalog.BarcodeLabelCorporateName,
	}
	// labels associa o discriminador à legenda visível, sem gravar texto derivado no modelo.
	labels := map[string]string{}
	// choice representa uma legenda visível conhecida do enum.
	for choice := range kinds {
		labels[kinds[choice]] = choice
	}
	return choices, kinds, labels
}

// barcodeBlockLabel apresenta a posição e o título sem depender de título único.
func barcodeBlockLabel(index int, title string) string {
	// shown contém o título atual ou um marcador explícito de ausência.
	shown := title
	if shown == "" {
		shown = "(sem título)"
	}
	return fmt.Sprintf("%d · %s", index+1, shown)
}

// barcodeRowLabel apresenta posição, rótulo derivado e valor de forma curta.
func barcodeRowLabel(index int, row catalog.BarcodeRow) string {
	// text contém o rótulo derivado do enum.
	text, ok := catalog.BarcodeLabelText(row.LabelKind)
	if !ok {
		text = row.LabelKind + ":"
	}
	// value contém o valor atual ou um marcador explícito de vazio.
	value := row.Value
	if value == "" {
		value = "(vazio)"
	}
	return fmt.Sprintf("%d · %s %s", index+1, text, value)
}

// barcodeBlockIndexByID localiza o bloco pela identidade estável.
func barcodeBlockIndexByID(blocks []catalog.BarcodeBlock, id string) int {
	// index percorre a ordem confirmada do rascunho.
	for index := range blocks {
		if blocks[index].ID == id {
			return index
		}
	}
	return -1
}

// barcodeBlockByID devolve o bloco mutável associado à identidade, ou nil.
func barcodeBlockByID(blocks []catalog.BarcodeBlock, id string) *catalog.BarcodeBlock {
	// index contém a posição encontrada ou -1.
	index := barcodeBlockIndexByID(blocks, id)
	if index < 0 {
		return nil
	}
	return &blocks[index]
}

// barcodeRowIndexByID localiza a linha pela identidade estável.
func barcodeRowIndexByID(rows []catalog.BarcodeRow, id string) int {
	// index percorre a ordem visual das linhas.
	for index := range rows {
		if rows[index].ID == id {
			return index
		}
	}
	return -1
}

// barcodeRowByID devolve a linha mutável associada à identidade, ou nil.
func barcodeRowByID(rows []catalog.BarcodeRow, id string) *catalog.BarcodeRow {
	// index contém a posição encontrada ou -1.
	index := barcodeRowIndexByID(rows, id)
	if index < 0 {
		return nil
	}
	return &rows[index]
}

// barcodeLabelByID devolve o texto visível associado à identidade em um mapa de opções.
func barcodeLabelByID(labels map[string]string, id string) string {
	// label representa uma opção visível candidata.
	for label, candidate := range labels {
		if candidate == id {
			return label
		}
	}
	return ""
}

// barcodeProcessedFields monta o editor de rascunho dos blocos e linhas sem tocar no confirmado.
func (e *Editor) barcodeProcessedFields(form *widget.Form, component screen.Component) func(*screen.Component) error {
	// draft preserva uma cópia isolada dos blocos confirmados até a confirmação explícita.
	draft := catalog.BarcodeProcessedProperties{}
	if component.Properties.BarcodeProcessed != nil {
		// copied contém a variante separada da origem, sem slices compartilhados.
		copied := component.Properties.Copy().BarcodeProcessed
		draft = *copied
	}
	// supplyID entrega UUIDs novos apenas aos itens criados neste rascunho.
	supplyID := project.NewID
	// blockID e rowID identificam a seleção local por identidade, nunca por título.
	blockID, rowID := "", ""
	if len(draft.Blocks) > 0 {
		blockID = draft.Blocks[0].ID
		if len(draft.Blocks[0].Rows) > 0 {
			rowID = draft.Blocks[0].Rows[0].ID
		}
	}
	// choices, kinds e labels traduzem o enum de rótulo para o seletor local.
	choices, kinds, labels := barcodeLabelChoices()
	// blockSelect e rowSelect escolhem itens por posição visível e identidade interna.
	blockSelect := newSelect(nil, nil)
	rowSelect := newSelect(nil, nil)
	// blockTitle, labelSelect e value mantêm o rascunho textual do bloco e da linha.
	blockTitle := entry("")
	labelSelect := newSelect(choices, nil)
	value := entry("")
	// message informa recusas locais sem abrir diálogo nem publicar estado.
	message := widget.NewLabel("")
	message.Wrapping = fyne.TextWrapWord
	// syncing impede que atualizações programáticas disparem callbacks de seleção.
	syncing := false
	// blockLabels e rowLabels associam o texto visível ao identificador estável.
	blockLabels := map[string]string{}
	rowLabels := map[string]string{}
	// flush copia os campos visíveis para o rascunho antes de trocar de seleção ou confirmar.
	flush := func() {
		// block é o bloco selecionado, localizado pela identidade.
		block := barcodeBlockByID(draft.Blocks, blockID)
		if block == nil {
			return
		}
		block.Title = blockTitle.Text
		// row é a linha selecionada dentro do bloco atual.
		row := barcodeRowByID(block.Rows, rowID)
		if row == nil {
			return
		}
		if kind, ok := kinds[labelSelect.Selected]; ok {
			row.LabelKind = kind
		}
		row.Value = value.Text
	}
	// loadFields recarrega os controles a partir da seleção e do rascunho.
	loadFields := func() {
		syncing = true
		// block é o bloco selecionado, localizado pela identidade.
		block := barcodeBlockByID(draft.Blocks, blockID)
		if block == nil {
			blockTitle.SetText("")
			labelSelect.ClearSelected()
			value.SetText("")
			syncing = false
			return
		}
		blockTitle.SetText(block.Title)
		// row é a linha selecionada dentro do bloco atual.
		row := barcodeRowByID(block.Rows, rowID)
		labelSelect.ClearSelected()
		if row == nil {
			value.SetText("")
			syncing = false
			return
		}
		if label, ok := labels[row.LabelKind]; ok {
			labelSelect.SetSelected(label)
		}
		value.SetText(row.Value)
		syncing = false
	}
	// refreshSelectors reconstrói as listas visíveis a partir das identidades do rascunho.
	refreshSelectors := func() {
		syncing = true
		// blockOptions preserva a ordem dos blocos no seletor.
		blockOptions := []string{}
		blockLabels = map[string]string{}
		// block indexado recebe posição e título apenas para exibição.
		for index, block := range draft.Blocks {
			// label contém posição e título, nunca uma chave persistida.
			label := barcodeBlockLabel(index, block.Title)
			blockOptions = append(blockOptions, label)
			blockLabels[label] = block.ID
		}
		blockSelect.Options = blockOptions
		blockSelect.ClearSelected()
		if label := barcodeLabelByID(blockLabels, blockID); label != "" {
			blockSelect.SetSelected(label)
		}
		// block é o bloco cujas linhas serão listadas.
		block := barcodeBlockByID(draft.Blocks, blockID)
		// rowOptions preserva a ordem visual das linhas do bloco atual.
		rowOptions := []string{}
		rowLabels = map[string]string{}
		if block != nil {
			// row indexado recebe posição, rótulo derivado e valor apenas para exibição.
			for index, row := range block.Rows {
				// label contém posição, rótulo e valor, nunca uma chave persistida.
				label := barcodeRowLabel(index, row)
				rowOptions = append(rowOptions, label)
				rowLabels[label] = row.ID
			}
		}
		rowSelect.Options = rowOptions
		rowSelect.ClearSelected()
		if label := barcodeLabelByID(rowLabels, rowID); label != "" {
			rowSelect.SetSelected(label)
		}
		blockSelect.Refresh()
		rowSelect.Refresh()
		syncing = false
		loadFields()
	}
	// addBlock acrescenta um bloco vazio com uma linha padrão e identidade nova.
	addBlock := func() {
		flush()
		draft.Blocks = append(draft.Blocks, catalog.NewBarcodeBlock(supplyID))
		// created contém o bloco recém-criado na posição final.
		created := draft.Blocks[len(draft.Blocks)-1]
		blockID, rowID = created.ID, ""
		if len(created.Rows) > 0 {
			rowID = created.Rows[0].ID
		}
		message.SetText("")
		refreshSelectors()
	}
	// removeBlock recusa a última unidade e preserva as identidades restantes.
	removeBlock := func() {
		if len(draft.Blocks) <= 1 {
			message.SetText("Mantenha ao menos um bloco.")
			return
		}
		flush()
		// index contém a posição do bloco selecionado.
		index := barcodeBlockIndexByID(draft.Blocks, blockID)
		if index < 0 {
			return
		}
		draft.Blocks = append(draft.Blocks[:index], draft.Blocks[index+1:]...)
		if index >= len(draft.Blocks) {
			index = len(draft.Blocks) - 1
		}
		blockID = draft.Blocks[index].ID
		rowID = ""
		if len(draft.Blocks[index].Rows) > 0 {
			rowID = draft.Blocks[index].Rows[0].ID
		}
		message.SetText("")
		refreshSelectors()
	}
	// addRow acrescenta uma linha padrão ao bloco selecionado.
	addRow := func() {
		flush()
		// block é o bloco que receberá a nova linha.
		block := barcodeBlockByID(draft.Blocks, blockID)
		if block == nil {
			return
		}
		block.Rows = append(block.Rows, catalog.NewBarcodeRow(supplyID))
		rowID = block.Rows[len(block.Rows)-1].ID
		message.SetText("")
		refreshSelectors()
	}
	// removeRow recusa a última linha do bloco e preserva as demais identidades.
	removeRow := func() {
		// block é o bloco que perderá a linha selecionada.
		block := barcodeBlockByID(draft.Blocks, blockID)
		if block == nil {
			return
		}
		if len(block.Rows) <= 1 {
			message.SetText("Mantenha ao menos uma linha no bloco.")
			return
		}
		flush()
		// index contém a posição da linha selecionada.
		index := barcodeRowIndexByID(block.Rows, rowID)
		if index < 0 {
			return
		}
		block.Rows = append(block.Rows[:index], block.Rows[index+1:]...)
		if index >= len(block.Rows) {
			index = len(block.Rows) - 1
		}
		rowID = block.Rows[index].ID
		message.SetText("")
		refreshSelectors()
	}
	// moveBlock troca o bloco selecionado com seu vizinho na ordem visual.
	moveBlock := func(delta int) {
		flush()
		// index contém a posição atual do bloco selecionado.
		index := barcodeBlockIndexByID(draft.Blocks, blockID)
		// target contém a posição vizinha solicitada.
		target := index + delta
		if index < 0 || target < 0 || target >= len(draft.Blocks) {
			message.SetText("O bloco já está no limite.")
			return
		}
		draft.Blocks[index], draft.Blocks[target] = draft.Blocks[target], draft.Blocks[index]
		message.SetText("")
		refreshSelectors()
	}
	// moveRow troca a linha selecionada com sua vizinha no mesmo bloco.
	moveRow := func(delta int) {
		flush()
		// block é o bloco que contém a linha selecionada.
		block := barcodeBlockByID(draft.Blocks, blockID)
		if block == nil {
			return
		}
		// index contém a posição atual da linha selecionada.
		index := barcodeRowIndexByID(block.Rows, rowID)
		// target contém a posição vizinha solicitada.
		target := index + delta
		if index < 0 || target < 0 || target >= len(block.Rows) {
			message.SetText("A linha já está no limite.")
			return
		}
		block.Rows[index], block.Rows[target] = block.Rows[target], block.Rows[index]
		message.SetText("")
		refreshSelectors()
	}
	blockSelect.OnChanged = func(label string) {
		if syncing {
			return
		}
		flush()
		// id contém a identidade do bloco escolhido.
		id := blockLabels[label]
		if id == "" {
			return
		}
		blockID = id
		rowID = ""
		if block := barcodeBlockByID(draft.Blocks, blockID); block != nil && len(block.Rows) > 0 {
			rowID = block.Rows[0].ID
		}
		message.SetText("")
		refreshSelectors()
	}
	rowSelect.OnChanged = func(label string) {
		if syncing {
			return
		}
		flush()
		// id contém a identidade da linha escolhida.
		id := rowLabels[label]
		if id == "" {
			return
		}
		rowID = id
		message.SetText("")
		refreshSelectors()
	}
	labelSelect.OnChanged = func(string) {
		if syncing {
			return
		}
		flush()
		message.SetText("")
		refreshSelectors()
	}
	form.Append("Bloco", blockSelect)
	form.Append("Título do bloco", blockTitle)
	form.Append("Linha", rowSelect)
	form.Append("Rótulo", labelSelect)
	form.Append("Valor", value)
	form.Append("Blocos", container.NewGridWithColumns(2,
		widget.NewButton("Adicionar bloco", addBlock),
		widget.NewButton("Remover bloco", removeBlock),
		widget.NewButton("Bloco ↑", func() { moveBlock(-1) }),
		widget.NewButton("Bloco ↓", func() { moveBlock(1) }),
	))
	form.Append("Linhas", container.NewGridWithColumns(2,
		widget.NewButton("Adicionar linha", addRow),
		widget.NewButton("Remover linha", removeRow),
		widget.NewButton("Subir", func() { moveRow(-1) }),
		widget.NewButton("Descer", func() { moveRow(1) }),
	))
	form.Append("", message)
	refreshSelectors()
	return func(target *screen.Component) error {
		flush()
		if target.Properties.BarcodeProcessed == nil {
			return fmt.Errorf("componente sem propriedades de código de barras processado")
		}
		if err := draft.Validate(); err != nil {
			return err
		}
		// blocks separa a configuração confirmada das listas ainda usadas pelo formulário.
		blocks := make([]catalog.BarcodeBlock, len(draft.Blocks))
		// index percorre os blocos do rascunho na ordem visual.
		for index, block := range draft.Blocks {
			// blockCopy separa as linhas mutáveis mantendo identidade e ordem.
			blockCopy := block
			blockCopy.Rows = append([]catalog.BarcodeRow{}, block.Rows...)
			blocks[index] = blockCopy
		}
		target.Properties.BarcodeProcessed.Blocks = blocks
		return target.Properties.Validate(target.Kind)
	}
}
