package desktop

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// addScreen acrescenta uma composição vazia em uma única operação.
func (e *Editor) addScreen() {
	// created identifica a composição adicionada.
	created := screen.New(fmt.Sprintf("Tela %d", len(e.Session.Document.Screens)+1))
	if e.check(e.Session.Change(func(document *application.Document) error {
		document.Screens = append(document.Screens, created)
		if current := document.Screen(e.PageID); current != nil && current.Background.AssetID != "" {
			document.Screen(created.ID).Background = current.Background
		} else if err := builtinassets.UseBackground(document, created.ID); err != nil {
			return err
		}
		if document.Journey.StartScreenID == "" {
			document.Journey.StartScreenID = created.ID
		}
		return nil
	})) {
		e.PageID = created.ID
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

// duplicateComponent cria um controle independente e copia somente sua ação.
func (e *Editor) duplicateComponent() {
	if e.page() == nil || e.page().ComponentByID(e.Selected) == nil {
		return
	}
	// selected é o ID antes de publicar a cópia.
	selected := e.Selected
	// created será selecionado após confirmação.
	created := ""
	if e.check(e.Session.Change(func(document *application.Document) error {
		// page representa a tela do comando.
		page := document.Screen(e.PageID)
		// duplicate recebe propriedades profundamente copiadas.
		duplicate := *page.ComponentByID(selected)
		duplicate.Properties = duplicate.Properties.Copy()
		duplicate.ID = project.NewID()
		duplicate.Name = application.UniqueName(*page, duplicate.Name+"_copia")
		duplicate.Bounds.X += 10
		duplicate.Bounds.Y += 10
		duplicate.Bounds = duplicate.Bounds.Clamp(page.Size)
		page.Components = append(page.Components, duplicate)
		created = duplicate.ID
		// action copia somente o binding pertencente ao botão.
		for _, action := range append([]journey.Action{}, document.Journey.Actions...) {
			if action.SourceComponentID == selected {
				action.SourceComponentID = created
				document.Journey.Actions = append(document.Journey.Actions, action)
			}
		}
		return nil
	})) {
		e.Selected = created
		e.changed()
	}
}

// removeComponent preserva referências externas como diagnósticos reparáveis.
func (e *Editor) removeComponent() {
	if e.Selected == "" || e.page() == nil {
		return
	}
	// id identifica o alvo que será removido.
	id := e.Selected
	dialog.ShowConfirm("Remover componente", "Remover o componente? Roteiros ou teclados vinculados poderão precisar de correção. Você pode desfazer.", func(yes bool) {
		if !yes {
			return
		}
		if e.check(e.Session.Change(func(document *application.Document) error {
			// page é a composição afetada.
			page := document.Screen(e.PageID)
			// components conserva os irmãos e sua ordem.
			components := []screen.Component{}
			// component identifica o controle da tela que será consultado ou alterado.
			for _, component := range page.Components {
				if component.ID != id {
					components = append(components, component)
				}
			}
			page.Components = components
			// actions remove somente bindings de que o controle era dono.
			actions := []journey.Action{}
			// action contém a associação de navegação examinada na operação.
			for _, action := range document.Journey.Actions {
				if action.SourceComponentID != id {
					actions = append(actions, action)
				}
			}
			document.Journey.Actions = actions
			return nil
		})) {
			e.Selected = ""
			e.changed()
		}
	}, e.Window)
}

// moveSelection aplica deslocamento de teclado somente fora de um campo em edição.
func (e *Editor) moveSelection(dx, dy int) {
	if e.Window.Canvas().Focused() != nil || e.Selected == "" || e.page() == nil {
		return
	}
	e.check(e.Session.Change(func(document *application.Document) error {
		// page identifica a tela ativa.
		page := document.Screen(e.PageID)
		// component é o controle selecionado.
		component := page.ComponentByID(e.Selected)
		if component == nil {
			return nil
		}
		component.Bounds.X += dx
		component.Bounds.Y += dy
		component.Bounds = component.Bounds.Clamp(page.Size)
		return nil
	}))
	e.changed()
}

// installShortcuts integra edição e arquivo sem interceptar texto digitado.
func (e *Editor) installShortcuts() {
	e.Window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) { e.saveProject(false) })
	e.Window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyZ, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) {
		if e.Window.Canvas().Focused() == nil {
			e.check(e.Session.Undo())
			e.changed()
		}
	})
	e.Window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyY, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) {
		if e.Window.Canvas().Focused() == nil {
			e.check(e.Session.Redo())
			e.changed()
		}
	})
	e.Window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) {
		if e.Window.Canvas().Focused() == nil {
			e.duplicateComponent()
		}
	})
	e.Window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if e.Window.Canvas().Focused() != nil {
			return
		}
		switch event.Name {
		case fyne.KeyLeft:
			e.moveSelection(-1, 0)
		case fyne.KeyRight:
			e.moveSelection(1, 0)
		case fyne.KeyUp:
			e.moveSelection(0, -1)
		case fyne.KeyDown:
			e.moveSelection(0, 1)
		case fyne.KeyDelete:
			e.removeComponent()
		case fyne.KeyEscape:
			e.Board.Dragging = false
			e.Board.Refresh()
		}
	})
}
