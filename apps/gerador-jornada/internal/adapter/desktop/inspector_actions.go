package desktop

import "fyne.io/fyne/v2/widget"

// inspectorActions fixa a confirmação fora da rolagem e separa rascunho de edição confirmada.
func (e *Editor) inspectorActions(form *widget.Form) {
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
