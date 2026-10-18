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
		container.NewVScroll(container.NewVBox(e.palette(), widget.NewSeparator(), section("ELEMENTOS DO LEGADO"), e.presetsPalette())))
	// structure seleciona controles que estejam cobertos por outras camadas.
	structure := container.NewBorder(widget.NewLabel("Selecione um componente\npelo seu nome na tela."), nil, nil, nil, container.NewVScroll(e.Structure))
	// tabs preserva o espaço da paleta e não mistura suas alturas com a lista de telas.
	tabs := container.NewAppTabs(container.NewTabItem("Telas", pages), container.NewTabItem("Paleta", components), container.NewTabItem("Estrutura", structure))
	tabs.SelectIndex(1)
	return tabs
}

// refreshStructure mantém nomes, seleção e ordem de camadas sincronizados com o canvas.
func (e *Editor) refreshStructure() {
	e.Structure.RemoveAll()
	if e.page() == nil {
		return
	}
	// index percorre as camadas do primeiro plano para o fundo.
	for index := len(e.page().Components) - 1; index >= 0; index-- {
		// component mantém a identidade correta no callback de seleção.
		component := e.page().Components[index]
		// item apresenta o nome lógico sem depender da legenda visível.
		item := widget.NewButton(component.Name, func() {
			e.Selected = component.ID
			e.Window.Canvas().Unfocus()
			e.refreshInspector()
			e.refreshStructure()
			e.Board.Refresh()
		})
		item.Alignment = widget.ButtonAlignLeading
		if component.ID == e.Selected {
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
	}
	return string(kind)
}
