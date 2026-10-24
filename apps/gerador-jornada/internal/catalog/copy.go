package catalog

// Copy cria uma cópia independente, separando todos os dados mutáveis da origem.
func (p Properties) Copy() Properties {
	// result preserva os valores imutáveis e separa as coleções mutáveis da origem.
	result := p
	if p.Button != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Button
		result.Button = &v
	}
	if p.Label != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Label
		result.Label = &v
	}
	if p.Title != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Title
		result.Title = &v
	}
	if p.Field != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Field
		result.Field = &v
	}
	if p.Keyboard != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Keyboard
		result.Keyboard = &v
	}
	if p.Radio != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Radio
		v.Options = append([]Option{}, p.Radio.Options...)
		result.Radio = &v
	}
	if p.Report != nil {
		// v copia o valor apontado antes de separar suas coleções mutáveis.
		v := *p.Report
		v.Rows = append([]Option{}, p.Report.Rows...)
		result.Report = &v
	}
	if p.BarcodeInput != nil {
		// v copia o rótulo imutável da entrada.
		v := *p.BarcodeInput
		result.BarcodeInput = &v
	}
	if p.BarcodeProcessed != nil {
		// v separa cada bloco e suas linhas da configuração original.
		v := *p.BarcodeProcessed
		v.Blocks = []BarcodeBlock{}
		// block representa um bloco copiado com linhas independentes.
		for _, block := range p.BarcodeProcessed.Blocks {
			blockCopy := block
			blockCopy.Rows = append([]BarcodeRow{}, block.Rows...)
			v.Blocks = append(v.Blocks, blockCopy)
		}
		result.BarcodeProcessed = &v
	}
	if p.Menu != nil {
		result.Menu = copyApplicationMenu(p.Menu)
	}
	return result
}

// copyApplicationMenu copia a árvore preservando identidades, ordem e textos.
func copyApplicationMenu(source *ApplicationMenuProperties) *ApplicationMenuProperties {
	// clone separa ponteiros e slices sem trocar UUIDs.
	var clone func(node *SubmenuProperties) *SubmenuProperties
	clone = func(node *SubmenuProperties) *SubmenuProperties {
		if node == nil {
			return nil
		}
		// copy preserva o nó e suas propriedades imutáveis.
		copy := *node
		copy.Items = []MenuItem{}
		// item percorre a ordem original para separar os ramos.
		for _, item := range node.Items {
			itemCopy := item
			itemCopy.Submenu = clone(item.Submenu)
			copy.Items = append(copy.Items, itemCopy)
		}
		return &copy
	}
	// result mantém a escolha inicial e a raiz copiada.
	result := *source
	result.Root = *clone(&source.Root)
	return &result
}
