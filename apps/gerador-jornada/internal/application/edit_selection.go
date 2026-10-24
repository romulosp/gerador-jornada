package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// DeleteComponents remove o conjunto e seus bindings de origem em uma única ação reversível.
// Referências externas, como roteiro e teclado, permanecem diagnosticáveis em vez de serem apagadas.
func (s *Session) DeleteComponents(pageID string, ids []string) error {
	return s.Change(func(document *Document) error {
		// page e seleção precisam ser válidas antes de publicar a remoção.
		page := document.Screen(pageID)
		if page == nil {
			return fmt.Errorf("tela inexistente")
		}
		if _, err := screen.SelectionBounds(*page, ids); err != nil {
			return err
		}
		// removed identifica o conjunto a eliminar sem alterar a ordem dos irmãos.
		removed := map[string]bool{}
		for _, id := range ids {
			removed[id] = true
		}
		// remaining conserva os controles que não pertencem ao comando.
		remaining := []screen.Component{}
		for _, component := range page.Components {
			if !removed[component.ID] {
				remaining = append(remaining, component)
			}
		}
		page.Components = remaining
		//Filhos de um componente removido permanecem na tela como raízes, sem referência inválida.
		for _, id := range ids {
			page.OrphanChildren(id)
		}
		// actions conserva bindings de outros componentes e suas identidades.
		actions := []journey.Action{}
		for _, action := range document.Journey.Actions {
			if !removed[action.SourceComponentID] {
				actions = append(actions, action)
			}
		}
		document.Journey.Actions = actions
		return nil
	})
}

// DuplicateComponents clona o conjunto, remapeia vínculos internos e conserva espaçamento nas bordas.
func (s *Session) DuplicateComponents(pageID string, ids []string) ([]string, error) {
	// created permite selecionar as novas identidades após a confirmação do comando.
	created := []string{}
	// err preserva o resultado da transação sem publicar cópias parciais.
	err := s.Change(func(document *Document) error {
		// page fornece a composição de origem e destino das cópias.
		page := document.Screen(pageID)
		if page == nil {
			return fmt.Errorf("tela inexistente")
		}
		if _, err := screen.SelectionBounds(*page, ids); err != nil {
			return err
		}
		// mapping garante um único novo identificador por componente solicitado.
		mapping := map[string]string{}
		for _, id := range ids {
			if mapping[id] != "" {
				continue
			}
			// clone separa propriedades mutáveis e mantém o estilo do original.
			clone := *page.ComponentByID(id)
			clone.Properties = clone.Properties.Copy()
			// compostos recebem identidades internas novas para não colidirem com o original.
			clone.Properties.RegenerateCompoundIDs(project.NewID)
			clone.ID = project.NewID()
			clone.Name = UniqueName(*page, clone.Name+"_copia")
			clone.FileID = UniqueComponentFileID(*page, clone.FileID)
			mapping[id] = clone.ID
			created = append(created, clone.ID)
			page.Components = append(page.Components, clone)
		}
		for _, id := range created {
			// clone pode apontar a um campo também duplicado dentro do conjunto.
			clone := page.ComponentByID(id)
			if clone.Properties.Keyboard != nil && mapping[clone.Properties.Keyboard.TargetComponentID] != "" {
				clone.Properties.Keyboard.TargetComponentID = mapping[clone.Properties.Keyboard.TargetComponentID]
			}
			// Pai copiado no mesmo conjunto recebe a nova identidade; pai externo é preservado.
			if mapped := mapping[clone.ParentID]; mapped != "" {
				clone.ParentID = mapped
			}
		}
		for _, action := range append([]journey.Action{}, document.Journey.Actions...) {
			if mapping[action.SourceComponentID] != "" {
				action.SourceComponentID = mapping[action.SourceComponentID]
				document.Journey.Actions = append(document.Journey.Actions, action)
			}
		}
		return screen.MoveSelection(page, created, 10, 10)
	})
	return created, err
}
