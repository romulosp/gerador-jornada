package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
)

// presetsPalette oferece os elementos auxiliares como composições arrastáveis e editáveis.
func (e *Editor) presetsPalette() fyne.CanvasObject {
	// items separa conjuntos prontos das sete famílias primitivas.
	items := container.NewVBox()
	// definition identifica a composição que será inserida ao clicar ou soltar.
	for _, definition := range application.Presets() {
		// item reutiliza o mesmo comportamento de arraste da paleta principal.
		item := &paletteItem{Editor: e, PresetID: definition.ID}
		item.Text = definition.Name
		item.Icon = theme.ContentAddIcon()
		item.OnTapped = func() {
			if item.PresetID == "header" {
				e.insertPreset(item.PresetID, max(0, e.page().Size.Width-426), 6)
			} else if item.PresetID == "pageTitle" {
				e.insertPreset(item.PresetID, 0, 100)
			} else {
				e.insertPreset(item.PresetID, 48, 153)
			}
		}
		item.ExtendBaseWidget(item)
		items.Add(item)
	}
	return items
}

// insertPreset confirma um conjunto inteiro e sua decoração em uma única ação reversível.
func (e *Editor) insertPreset(presetID string, x, y int) {
	// ids permite selecionar o primeiro item sem persistir seleção no documento.
	var ids []string
	if e.check(e.Session.Change(func(document *application.Document) error {
		// err preserva a falha antes de publicar qualquer item do conjunto.
		var err error
		ids, err = application.InsertPreset(document, e.PageID, presetID, x, y, builtinassets.TitleStrip)
		if err != nil {
			return err
		}
		return nil
	})) {
		e.Selected = ids[0]
		e.changed()
	}
}
