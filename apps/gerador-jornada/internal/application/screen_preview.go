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
	// A Prévia Tela Atual valida apenas a tela isolada: destinos externos de menu
	// continuam autorados e não podem impedir a exploração local.
	if err := page.Validate(); err != nil {
		return nil, err
	}
	// isolated limita o escopo e conserva os componentes e recursos da tela.
	isolated := document.Copy()
	// copiedPage pertence à cópia profunda, impedindo que a prévia altere o documento original.
	copiedPage := isolated.Screen(screenID)
	isolated.Screens = []screen.Screen{*copiedPage}
	isolated.Journey = journey.Journey{StartScreenID: screenID}
	isolated.Scenarios = nil
	// preview conserva o snapshot reduzido sem autorizar navegação entre telas.
	preview := &Preview{Document: isolated, AllowScreenNavigation: false}
	if err := preview.Enter(screenID); err != nil {
		return nil, err
	}
	return preview, nil
}
