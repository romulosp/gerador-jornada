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
	return result
}
