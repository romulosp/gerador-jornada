package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/screen"
)

// MoveComponents confirma uma translação coletiva como uma única ação reversível.
func (s *Session) MoveComponents(pageID string, ids []string, dx, dy int) error {
	return s.Change(func(document *Document) error {
		// page precisa existir antes de aplicar as regras geométricas do domínio Tela.
		page := document.Screen(pageID)
		if page == nil {
			return fmt.Errorf("tela inexistente")
		}
		return screen.MoveSelection(page, ids, dx, dy)
	})
}
