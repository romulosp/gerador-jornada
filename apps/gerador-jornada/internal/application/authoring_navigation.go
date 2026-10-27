package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
)

// SetJourneyStart altera somente o início para tela existente no documento.
// O chamador usa Session.Change para publicar a decisão manual/MCP no mesmo histórico.
func SetJourneyStart(document *Document, screenID string) error {
	if document.Screen(screenID) == nil {
		return &AuthoringError{Code: "invalidReference"}
	}
	document.Journey.StartScreenID = screenID
	return nil
}

// BindNavigation valida referências explícitas antes de ligar/remover a ação de um botão.
// Não resolve arquivos externos, metadados de menu ou código Java como destino implícito.
func BindNavigation(document *Document, action journey.Action) error {
	// page/source precisam pertencer ao documento e a origem deve ser botão instalado.
	page := document.Screen(action.SourceScreenID)
	if page == nil {
		return &AuthoringError{Code: "invalidReference"}
	}
	source := page.ComponentByID(action.SourceComponentID)
	if source == nil || source.Kind != catalog.Button {
		return &AuthoringError{Code: "invalidReference"}
	}
	if action.Kind == "none" {
		if action.TargetScreenID != "" {
			return &AuthoringError{Code: "invalidArgument"}
		}
		// actions conserva bindings de outros botões, na ordem original.
		actions := []journey.Action{}
		for _, existing := range document.Journey.Actions {
			if existing.SourceComponentID != action.SourceComponentID {
				actions = append(actions, existing)
			}
		}
		document.Journey.Actions = actions
		return nil
	}
	if action.Kind != "navigate" {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if document.Screen(action.TargetScreenID) == nil {
		return &AuthoringError{Code: "invalidReference"}
	}
	return document.Journey.Bind(action)
}
