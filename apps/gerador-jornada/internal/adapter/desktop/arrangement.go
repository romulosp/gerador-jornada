package desktop

import "gerador-jornada.local/editor/internal/application"

// align publica um ajuste em relação à superfície da tela.
func (e *Editor) align(kind string) {
	e.check(e.Session.Change(func(document *application.Document) error {
		// page e component são os alvos da operação reversível.
		page := document.Screen(e.PageID)
		// component identifica o controle da tela que será consultado ou alterado.
		component := page.ComponentByID(e.Selected)
		if component == nil {
			return nil
		}
		switch kind {
		case "h":
			component.Bounds.X = (page.Size.Width - component.Bounds.Width) / 2
		case "v":
			component.Bounds.Y = (page.Size.Height - component.Bounds.Height) / 2
		case "left":
			component.Bounds.X = 0
		case "right":
			component.Bounds.X = page.Size.Width - component.Bounds.Width
		}
		return nil
	}))
	e.changed()
}

// layer move somente uma posição na ordem de camadas.
func (e *Editor) layer(delta int) {
	e.check(e.Session.Change(func(document *application.Document) error {
		// page contém o z-order atual.
		page := document.Screen(e.PageID)
		// index procura o componente selecionado.
		for index := range page.Components {
			if page.Components[index].ID == e.Selected {
				// next calcula a posição pretendida e permite verificar os limites antes de mover.
				if next := index + delta; next >= 0 && next < len(page.Components) {
					page.Components[index], page.Components[next] = page.Components[next], page.Components[index]
				}
				break
			}
		}
		return nil
	}))
	e.changed()
}
