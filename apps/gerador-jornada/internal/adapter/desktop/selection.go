package desktop

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// selectionIDs concilia a seleção múltipla com ações existentes que definem um alvo primário.
func (e *Editor) selectionIDs() []string {
	if e.Selected == "" || e.page() == nil {
		e.Selection = nil
		return nil
	}
	// hasPrimary detecta se outra ação, como inserir ou trocar tela, substituiu a seleção.
	hasPrimary := false
	for _, id := range e.Selection {
		hasPrimary = hasPrimary || id == e.Selected
	}
	if !hasPrimary {
		e.Selection = []string{e.Selected}
	}
	// valid remove referências eliminadas por uma edição ou histórico.
	valid := []string{}
	for _, id := range e.Selection {
		if e.page().ComponentByID(id) != nil {
			valid = append(valid, id)
		}
	}
	e.Selection = valid
	return valid
}

// isSelected informa se uma identidade participa do conjunto ativo da interface.
func (e *Editor) isSelected(id string) bool {
	if len(e.Selection) == 0 {
		return id == e.Selected && id != ""
	}
	for _, selected := range e.Selection {
		if selected == id {
			return true
		}
	}
	return false
}

// setSelection publica apenas estado visual e informa se a seleção realmente mudou.
func (e *Editor) setSelection(ids []string) bool {
	// previous conserva a ordem dos alvos antes de aplicar uma alteração de seleção.
	previous := append([]string{}, e.selectionIDs()...)
	// seen elimina referências repetidas sem reordenar os controles.
	seen := map[string]bool{}
	e.Selection = nil
	for _, id := range ids {
		if !seen[id] && e.page() != nil && e.page().ComponentByID(id) != nil {
			seen[id] = true
			e.Selection = append(e.Selection, id)
		}
	}
	e.Selected = ""
	if len(e.Selection) > 0 {
		e.Selected = e.Selection[len(e.Selection)-1]
	}
	if len(previous) != len(e.Selection) {
		return true
	}
	for index := range previous {
		if previous[index] != e.Selection[index] {
			return true
		}
	}
	return false
}

// selectionChanged atualiza somente os painéis afetados, preservando o restante da janela.
func (e *Editor) selectionChanged() {
	e.refreshSelectionActions()
	e.refreshInspector()
	e.refreshStructure()
	e.Board.Refresh()
}

// refreshSelectionActions reflete a seleção válida sem modificar o documento ou seu histórico.
func (e *Editor) refreshSelectionActions() {
	if e.DeleteButton == nil {
		return
	}
	if len(e.selectionIDs()) == 0 {
		e.DeleteButton.Disable()
	} else {
		e.DeleteButton.Enable()
	}
}

// selectAll seleciona os controles da tela sem capturar Ctrl+A de um campo textual.
func (e *Editor) selectAll() {
	if e.page() == nil || e.Window.Canvas().Focused() != nil {
		return
	}
	// ids preserva a ordem visual da composição.
	ids := []string{}
	for _, component := range e.page().Components {
		ids = append(ids, component.ID)
	}
	if e.setSelection(ids) {
		e.selectionChanged()
	}
}

// groupInspector oferece comandos coletivos sem fingir editar propriedades individuais em lote.
func (e *Editor) groupInspector(count int) {
	e.Inspector.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d componentes selecionados", count), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	e.Inspector.Add(widget.NewLabel("Arraste o conjunto ou use as setas.\nShift + setas move 10 unidades.\nCtrl/Shift + clique alterna a seleção."))
	// clear retorna ao formulário de tela sem alterar o documento.
	e.Inspector.Add(widget.NewButton("Limpar seleção", func() { e.setSelection(nil); e.selectionChanged() }))
	e.Inspector.Add(widget.NewButton("Duplicar seleção", e.duplicateComponent))
	e.Inspector.Add(widget.NewButton("Remover seleção", e.removeComponent))
}
