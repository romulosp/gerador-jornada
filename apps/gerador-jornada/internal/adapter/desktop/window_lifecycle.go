package desktop

import "fyne.io/fyne/v2/dialog"

// requestClose compartilha o encerramento do botão Sair com o fechamento nativo da janela.
// Alterações no documento ou rascunhos exigem confirmação; cancelar preserva a sessão e a exibição.
func (e *Editor) requestClose() {
	if e.Mode.closed.Load() || e.Mode.exitDialog != nil {
		return
	}
	if !e.Session.Dirty() && (e.InspectorHasDraft == nil || !e.InspectorHasDraft()) {
		e.close()
		return
	}
	e.Mode.exitDialog = dialog.NewConfirm("Alterações não salvas", "Sair sem salvar as alterações?", func(yes bool) {
		e.Mode.exitDialog = nil
		if yes {
			e.close()
		}
	}, e.Window)
	e.Mode.exitDialog.SetConfirmText("Sim")
	e.Mode.exitDialog.SetDismissText("Não")
	e.Mode.exitDialog.Show()
}
