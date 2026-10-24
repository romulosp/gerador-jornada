package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
)

// toolsPanel separa páginas, catálogo e estrutura, como nos fluxos do legado.
func (e *Editor) toolsPanel() fyne.CanvasObject {
	// pages mantém a lista de telas com sua própria rolagem.
	pages := container.NewBorder(nil, widget.NewButtonWithIcon("Adicionar tela", theme.ContentAddIcon(), e.addScreen), nil, nil,
		container.NewVScroll(container.NewVBox(section("TELAS DA JORNADA"), e.Screens)))
	// components torna todas as famílias alcançáveis sem rolar a lista de páginas.
	components := container.NewBorder(widget.NewLabel("Arraste para a imagem\nou clique para inserir."), nil, nil, nil,
		container.NewVScroll(container.NewVBox(e.palette(), widget.NewSeparator(), section("ELEMENTOS DO LEGADO"), e.presetsPalette(), widget.NewSeparator(), e.legacyPalette())))
	// structure seleciona controles que estejam cobertos por outras camadas.
	structure := container.NewBorder(widget.NewLabel("Selecione um componente\npelo seu nome na tela."), nil, nil, nil, container.NewVScroll(e.Structure))
	// tabs preserva o espaço da paleta e não mistura suas alturas com a lista de telas.
	tabs := container.NewAppTabs(container.NewTabItem("Estrutura projeto", e.workspacePanel()), container.NewTabItem("Telas", pages), container.NewTabItem("Paleta", components), container.NewTabItem("Estrutura", structure))
	tabs.SelectIndex(0)
	return tabs
}

// refreshStructure mantém nomes, seleção e ordem de camadas sincronizados com o canvas.
func (e *Editor) refreshStructure() {
	if e.page() == nil {
		e.Structure.RemoveAll()
		return
	}
	if e.StructurePage == e.PageID && e.StructureRevision == e.Session.Document.Project.Revision && len(e.Structure.Objects) == len(e.page().Components) {
		// object conserva os botões existentes; somente mudanças de realce provocam refresh.
		for index, object := range e.Structure.Objects {
			if item, ok := object.(*widget.Button); ok {
				// importance reflete a seleção sem reconstruir a lista nem seus callbacks.
				importance := widget.MediumImportance
				if e.isSelected(e.page().Components[len(e.page().Components)-1-index].ID) {
					importance = widget.HighImportance
				}
				if item.Importance != importance {
					item.Importance = importance
					item.Refresh()
				}
			}
		}
		return
	}
	e.StructurePage, e.StructureRevision = e.PageID, e.Session.Document.Project.Revision
	e.Structure.RemoveAll()
	// index percorre as camadas do primeiro plano para o fundo.
	for index := len(e.page().Components) - 1; index >= 0; index-- {
		// component mantém a identidade correta no callback de seleção.
		component := e.page().Components[index]
		// item apresenta o nome lógico sem depender da legenda visível.
		item := widget.NewButton(component.Name, func() {
			e.Window.Canvas().Unfocus()
			if e.setSelection([]string{component.ID}) {
				e.selectionChanged()
			}
		})
		item.Alignment = widget.ButtonAlignLeading
		if e.isSelected(component.ID) {
			item.Importance = widget.HighImportance
		}
		e.Structure.Add(item)
	}
	if len(e.page().Components) == 0 {
		e.Structure.Add(widget.NewLabel("Nenhum componente nesta tela."))
	}
}

// familyName relaciona o nome amigável ao tipo correspondente no legado.
func familyName(kind catalog.Kind) string {
	switch kind {
	case catalog.Button:
		return "Botão · EFLButton"
	case catalog.Label:
		return "Rótulo · EFLLabel"
	case catalog.Title:
		return "Título · EFLLabelTitulo1"
	case catalog.TextField:
		return "Campo de texto · EFLTextField"
	case catalog.Radio:
		return "Seleção · EFLRadio"
	case catalog.Report:
		return "Relatório · EFLReportEvent"
	case catalog.NumericKeyboard:
		return "Teclado · EFLVirtualKeyboardNumeric"
	case catalog.BarcodeInput:
		return "Código de barras · entrada"
	case catalog.BarcodeProcessed:
		return "Código de barras processado · ficha"
	case catalog.ApplicationMenu:
		return "Menu da aplicação · árvore"
	}
	return string(kind)
}
