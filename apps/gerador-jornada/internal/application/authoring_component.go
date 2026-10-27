package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// ComponentChanges recebe valores tipados completos por campo alterado, sem DTO MCP.
type ComponentChanges struct {
	// Name/FileID/Enabled preservam os campos omitidos na solicitação do autor.
	Name, FileID *string
	Enabled      *bool
	// Bounds/Style/Properties já foram adaptados da família correta, sem maps livres.
	Bounds     *screen.Bounds
	Style      *screen.Style
	Properties *catalog.Properties
}

// ConfigureComponent aplica candidato validado em uma única transação reversível.
// appearance valida fatos dos recursos incorporados no adaptador, sem importar UI/IO na aplicação.
func (session *Session) ConfigureComponent(pageID, componentID string, changes ComponentChanges, appearance func(screen.Component, screen.Component) error) error {
	return session.Change(func(document *Document) error {
		// page/target pertencem à cópia provisória e nunca são consultados por caminho externo.
		page := document.Screen(pageID)
		if page == nil {
			return &AuthoringError{Code: "invalidReference"}
		}
		target := page.ComponentByID(componentID)
		if target == nil {
			return &AuthoringError{Code: "invalidReference"}
		}
		before := *target
		if changes.Name != nil {
			target.Name = *changes.Name
		}
		if changes.FileID != nil {
			target.FileID = *changes.FileID
		}
		if changes.Enabled != nil {
			target.Enabled = *changes.Enabled
		}
		if changes.Bounds != nil {
			target.Bounds = *changes.Bounds
		}
		if changes.Style != nil {
			target.Style = *changes.Style
		}
		if changes.Properties != nil {
			target.Properties = changes.Properties.Copy()
		}
		if err := ValidateComponentEdit(document, page, before, *target); err != nil {
			return err
		}
		if appearance != nil {
			if err := appearance(before, *target); err != nil {
				return err
			}
		}
		return page.Validate()
	})
}

// DeleteComponentWithPolicy preserva referências externas e permite remover somente bindings próprios.
func (session *Session) DeleteComponentWithPolicy(pageID, componentID, policy string) error {
	if policy != "reject" && policy != "removeOwnedReferences" {
		return &AuthoringError{Code: "invalidArgument"}
	}
	// page/target delimitam referências antes de delegar à remoção compartilhada existente.
	page := session.Document.Screen(pageID)
	if page == nil || page.ComponentByID(componentID) == nil {
		return &AuthoringError{Code: "invalidReference"}
	}
	for _, other := range page.Components {
		if other.ID == componentID {
			continue
		}
		if other.ParentID == componentID || (other.Properties.Keyboard != nil && other.Properties.Keyboard.TargetComponentID == componentID) {
			return &AuthoringError{Code: "invalidReference"}
		}
	}
	for _, scenario := range session.Document.Scenarios {
		for _, step := range scenario.Steps {
			if step.ComponentID == componentID {
				return &AuthoringError{Code: "invalidReference"}
			}
		}
	}
	if policy == "reject" {
		for _, action := range session.Document.Journey.Actions {
			if action.SourceComponentID == componentID {
				return &AuthoringError{Code: "invalidReference"}
			}
		}
	}
	return session.DeleteComponents(pageID, []string{componentID})
}
