package desktop

import (
	"slices"
	"strconv"

	"fyne.io/fyne/v2/widget"
)

// inspectorActions fixa a confirmação fora da rolagem e separa rascunho de edição confirmada.
func (e *Editor) inspectorActions(form *widget.Form) {
	// baseline conserva valores iniciais sem substituir callbacks de validação dos campos.
	baseline := inspectorValues(form)
	e.InspectorHasDraft = func() bool { return !slices.Equal(baseline, inspectorValues(form)) }
	// submit conserva a validação do formulário antes de ocultar seus botões internos.
	submit := form.OnSubmit
	form.OnSubmit = nil
	form.Refresh()
	// save sempre permanece acessível na base do painel de propriedades.
	save := widget.NewButton("Salvar alterações", submit)
	save.Importance = widget.HighImportance
	e.InspectorActions.Add(save)
	e.InspectorActions.Add(widget.NewButton("Descartar", e.refreshInspector))
}

// inspectorValues consulta valores editáveis para proteger rascunhos ao encerrar o programa.
func inspectorValues(form *widget.Form) []string {
	// values preserva a ordem dos campos do formulário sem serializar widgets no projeto.
	values := make([]string, 0, len(form.Items))
	for _, item := range form.Items {
		// field identifica somente controles com valores que o usuário pode editar.
		switch field := item.Widget.(type) {
		case *widget.Entry:
			values = append(values, field.Text)
		case *widget.Select:
			values = append(values, field.Selected)
		case *widget.Check:
			values = append(values, strconv.FormatBool(field.Checked))
		}
	}
	return values
}
