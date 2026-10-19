package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// installShortcuts integra edição e arquivo sem interceptar texto digitado.
func (e *Editor) installShortcuts() {
	e.Window.Canvas().AddShortcut(&fyne.ShortcutSelectAll{}, func(fyne.Shortcut) { e.selectAll() })
	// key e delta oferecem movimento acelerado sem interceptar setas normais de um campo.
	for key, delta := range map[fyne.KeyName][2]int{fyne.KeyLeft: {-10, 0}, fyne.KeyRight: {10, 0}, fyne.KeyUp: {0, -10}, fyne.KeyDown: {0, 10}} {
		e.Window.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: key, Modifier: fyne.KeyModifierShift}, func(fyne.Shortcut) { e.moveSelection(delta[0], delta[1]) })
	}
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
		// modifiers vem do driver Windows; Shift+setas não é emitido como CustomShortcut pelo Fyne.
		modifiers := fyne.KeyModifier(0)
		if driver, ok := e.App.Driver().(desktop.Driver); ok {
			modifiers = driver.CurrentKeyModifiers()
		}
		e.handleCanvasKey(event, modifiers)
	})
}

// handleCanvasKey aplica atalhos geométricos sem consumir teclas de campos com foco.
func (e *Editor) handleCanvasKey(event *fyne.KeyEvent, modifiers fyne.KeyModifier) {
	if e.Mode.handleKey(event) {
		return
	}
	if e.Window.Canvas().Focused() != nil {
		return
	}
	// step define movimento fino ou acelerado em coordenadas lógicas.
	step := 1
	if modifiers&fyne.KeyModifierShift != 0 {
		step = 10
	}
	switch event.Name {
	case fyne.KeyLeft:
		e.moveSelection(-step, 0)
	case fyne.KeyRight:
		e.moveSelection(step, 0)
	case fyne.KeyUp:
		e.moveSelection(0, -step)
	case fyne.KeyDown:
		e.moveSelection(0, step)
	case fyne.KeyDelete:
		e.removeComponent()
	case fyne.KeyEscape:
		e.Board.Dragging, e.Board.Marquee = false, false
		e.setSelection(nil)
		e.selectionChanged()
	}
}
