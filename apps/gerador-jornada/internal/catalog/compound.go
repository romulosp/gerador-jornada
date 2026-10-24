package catalog

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
