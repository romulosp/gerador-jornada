package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/screen"
)

// NewScreenPreview valida e copia somente a tela indicada, sem executar a jornada ou seus roteiros.
// Campos, escolhas e recursos continuam disponíveis; botões não navegam para outras telas.
func NewScreenPreview(document Document, screenID string) (*Preview, error) {
	// page resolve a tela solicitada, independentemente da tela inicial da jornada.
	page := document.Screen(screenID)
	if page == nil {
		return nil, fmt.Errorf("selecione uma tela existente para visualizar")
	}
	// isolated limita o escopo antes da validação e da cópia profunda feita por NewPreview.
	isolated := document
	isolated.Screens = []screen.Screen{*page}
	isolated.Journey = journey.Journey{StartScreenID: screenID}
	isolated.Scenarios = nil
	return NewPreview(isolated)
}
