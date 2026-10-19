package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// componentAt encontra o componente superior no ponto local do canvas, respeitando zoom e camadas.
func (b *Board) componentAt(position fyne.Position) string {
	// page contém os componentes da composição atual, em ordem de camadas.
	page := b.Editor.page()
	if page == nil {
		return ""
	}
	// x e y removem a margem e o zoom antes da comparação com a geometria de domínio.
	x, y := b.logical(position)
	// index percorre os controles da frente para o fundo.
	for index := len(page.Components) - 1; index >= 0; index-- {
		// component fornece os limites do candidato ao clique.
		component := page.Components[index]
		if x >= component.Bounds.X && x <= component.Bounds.X+component.Bounds.Width && y >= component.Bounds.Y && y <= component.Bounds.Y+component.Bounds.Height {
			return component.ID
		}
	}
	return ""
}

// TappedSecondary abre ações do componente clicado sem iniciar arraste ou editar o projeto.
func (b *Board) TappedSecondary(event *fyne.PointEvent) {
	if b.ContextMenu != nil {
		b.ContextMenu.Hide()
		b.ContextMenu = nil
	}
	// id identifica o alvo superior; área vazia não oferece exclusão de componente.
	id := b.componentAt(event.Position)
	if id == "" {
		return
	}
	b.Editor.Window.Canvas().Unfocus()
	if !b.Editor.isSelected(id) && b.Editor.setSelection([]string{id}) {
		b.Editor.selectionChanged()
	}
	// pageID protege a ação de propriedades caso a composição mude antes da execução.
	pageID := b.Editor.PageID
	// menu oferece edição individual e exclusão da seleção que contém o alvo clicado.
	menu := fyne.NewMenu("Componente", fyne.NewMenuItem("Alterar propriedades", func() {
		if b.Editor.PageID == pageID && b.Editor.page() != nil && b.Editor.page().ComponentByID(id) != nil {
			b.Editor.setSelection([]string{id})
			b.Editor.selectionChanged()
			b.Editor.focusComponentProperties()
		}
	}), fyne.NewMenuItem("Excluir elemento", b.Editor.removeComponent))
	if len(b.Editor.selectionIDs()) > 1 {
		menu.Items[1].Label = "Excluir seleção"
	}
	b.ContextMenu = widget.NewPopUpMenu(menu, b.Editor.Window.Canvas())
	b.ContextMenu.ShowAtPosition(event.AbsolutePosition)
}

// focusComponentProperties direciona o teclado ao primeiro campo do formulário da seleção.
func (e *Editor) focusComponentProperties() {
	// object localiza o formulário sem acoplar o menu a uma família de componentes.
	for _, object := range e.Inspector.Objects {
		// form concentra os campos comuns e específicos já montados no painel direito.
		form, ok := object.(*widget.Form)
		if !ok {
			continue
		}
		// item identifica o primeiro campo que aceita edição por teclado.
		for _, item := range form.Items {
			if field, ok := item.Widget.(fyne.Focusable); ok {
				e.Window.Canvas().Focus(field)
				return
			}
		}
	}
}
