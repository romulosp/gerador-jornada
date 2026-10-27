package desktop

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
)

// addScreen acrescenta uma composição vazia em uma única operação.
func (e *Editor) addScreen() {
	// id é criado pelo mesmo caso de uso consumido pelo MCP; a UI fornece apenas contexto e recurso.
	id, err := e.Session.CreateScreen(fmt.Sprintf("Tela %d", len(e.Session.Document.Screens)+1), e.PageID, application.ScreenChanges{}, builtinassets.UseBackground)
	if e.check(err) {
		e.PageID = id
		e.Selected = ""
		e.changed()
	}
}

// duplicateScreen remapeia os vínculos internos usando o caso de uso.
func (e *Editor) duplicateScreen() {
	// id é atribuído pela transação de duplicação.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	id, err := e.Session.DuplicateScreen(e.PageID)
	if e.check(err) {
		e.PageID = id
		e.Selected = ""
		e.changed()
	}
}

// removeScreen mostra o impacto em navegação e roteiros antes de remover.
func (e *Editor) removeScreen() {
	if e.page() == nil {
		return
	}
	// affected conta referências que precisarão ser corrigidas.
	affected := 0
	// action é uma referência de navegação existente.
	for _, action := range e.Session.Document.Journey.Actions {
		if action.TargetScreenID == e.PageID {
			affected++
		}
	}
	// script contém as referências demonstrativas.
	for _, script := range e.Session.Document.Scenarios {
		// step representa a etapa cuja referência será preservada ou verificada.
		for _, step := range script.Steps {
			if step.ScreenID == e.PageID {
				affected++
			}
		}
	}
	dialog.ShowConfirm("Remover tela", fmt.Sprintf("Remover %q? %d referência(s) poderão precisar de correção. A operação pode ser desfeita.", e.page().Name, affected), func(yes bool) {
		if !yes {
			return
		}
		if e.check(e.Session.RemoveScreen(e.PageID)) {
			e.PageID = ""
			if len(e.Session.Document.Screens) > 0 {
				e.PageID = e.Session.Document.Screens[0].ID
			}
			e.Selected = ""
			e.changed()
		}
	}, e.Window)
}

// reorderScreen move a tela sem alterar IDs ou destinos de navegação.
func (e *Editor) reorderScreen(delta int) {
	e.check(e.Session.Change(func(document *application.Document) error {
		// index localiza a tela ativa na ordem visual.
		for index := range document.Screens {
			if document.Screens[index].ID == e.PageID {
				// next calcula a posição pretendida e permite verificar os limites antes de mover.
				if next := index + delta; next >= 0 && next < len(document.Screens) {
					document.Screens[index], document.Screens[next] = document.Screens[next], document.Screens[index]
				}
				break
			}
		}
		return nil
	}))
	e.changed()
}

// duplicateComponent duplica a seleção com vínculos internos em uma única transação.
func (e *Editor) duplicateComponent() {
	if e.page() == nil || e.Selected == "" {
		return
	}
	// created contém as identidades independentes retornadas pela aplicação.
	created, err := e.Session.DuplicateComponents(e.PageID, e.selectionIDs())
	if e.check(err) {
		e.setSelection(created)
		e.changed()
	}
}

// removeComponent confirma a remoção do conjunto mantendo referências externas diagnosticáveis.
func (e *Editor) removeComponent() {
	if e.DeleteDialog != nil || e.page() == nil {
		return
	}
	// ids e pageID capturam exatamente o conjunto mostrado na confirmação.
	ids := append([]string{}, e.selectionIDs()...)
	if len(ids) == 0 {
		return
	}
	pageID := e.PageID
	e.DeleteDialog = dialog.NewConfirm("Excluir componente", fmt.Sprintf("Deseja realmente excluir %d componente(s) selecionado(s)? Referências de roteiro ou teclado poderão precisar de correção. Você pode desfazer.", len(ids)), func(yes bool) {
		e.DeleteDialog = nil
		if yes && e.check(e.Session.DeleteComponents(pageID, ids)) {
			e.setSelection(nil)
			e.changed()
		}
	}, e.Window)
	e.DeleteDialog.SetConfirmText("Sim")
	e.DeleteDialog.SetDismissText("Não")
	e.DeleteDialog.Show()
}

// moveSelection aplica deslocamento de teclado somente fora de um campo em edição.
func (e *Editor) moveSelection(dx, dy int) {
	if e.Window.Canvas().Focused() != nil || e.Selected == "" || e.page() == nil {
		return
	}
	e.check(e.Session.MoveComponents(e.PageID, e.selectionIDs(), dx, dy))
	e.changed()
}
