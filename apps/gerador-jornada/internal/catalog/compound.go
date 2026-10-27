package catalog

// InstanceIDs reúne identidades internas de opções, linhas e partes de uma instância.
// Referências a outras entidades não pertencem à lista e não recebem UUID novo por cópia.
func (p Properties) InstanceIDs() []string {
	// result começa pelas partes dos compostos, mantendo sua ordem estrutural estável.
	result := append([]string{}, p.CompoundIDs()...)
	if p.Radio != nil {
		for _, option := range p.Radio.Options {
			result = append(result, option.ID)
		}
	}
	if p.Report != nil {
		for _, row := range p.Report.Rows {
			result = append(result, row.ID)
		}
	}
	return result
}

// RegenerateInstanceIDs recria todas as identidades proprietárias, inclusive seleção inicial do rádio.
func (p *Properties) RegenerateInstanceIDs(supplyID func() string) {
	if p == nil {
		return
	}
	p.RegenerateCompoundIDs(supplyID)
	if p.Radio != nil {
		// selected conserva a correspondência da opção escolhida, sem usar índice como identidade.
		selected := p.Radio.SelectedOptionID
		for index := range p.Radio.Options {
			previous := p.Radio.Options[index].ID
			p.Radio.Options[index].ID = supplyID()
			if previous == selected {
				p.Radio.SelectedOptionID = p.Radio.Options[index].ID
			}
		}
	}
	if p.Report != nil {
		for index := range p.Report.Rows {
			p.Report.Rows[index].ID = supplyID()
		}
	}
}

// CompoundIDs devolve as identidades internas dos compostos que pertencem ao documento.
// Componentes simples não possuem identidades internas e retornam uma lista vazia.
func (p Properties) CompoundIDs() []string {
	switch {
	case p.BarcodeProcessed != nil:
		return CollectBarcodeProcessedIDs(*p.BarcodeProcessed)
	case p.Menu != nil:
		return CollectMenuIDs(*p.Menu)
	}
	return nil
}

// RegenerateCompoundIDs substitui as identidades internas por UUIDs novos.
// A cópia profunda preserva textos, ordem e escolha inicial remapeada.
func (p *Properties) RegenerateCompoundIDs(supplyID func() string) {
	if p == nil {
		return
	}
	if p.BarcodeProcessed != nil {
		// blocks recria blocos e linhas com identidades novas, preservando o conteúdo.
		blocks := make([]BarcodeBlock, len(p.BarcodeProcessed.Blocks))
		for blockIndex, block := range p.BarcodeProcessed.Blocks {
			blockCopy := block
			blockCopy.ID = supplyID()
			blockCopy.Rows = make([]BarcodeRow, len(block.Rows))
			for rowIndex, row := range block.Rows {
				rowCopy := row
				rowCopy.ID = supplyID()
				blockCopy.Rows[rowIndex] = rowCopy
			}
			blocks[blockIndex] = blockCopy
		}
		p.BarcodeProcessed.Blocks = blocks
	}
	if p.Menu != nil {
		// duplicated regenera a árvore inteira e remapeia o menu inicial.
		duplicated := DuplicateApplicationMenu(*p.Menu, supplyID)
		p.Menu = &duplicated
	}
}
