package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// AddComponent adiciona uma instância do catálogo em um comando reversível.
func (s *Session) AddComponent(pageID string, kind catalog.Kind, x, y int) (string, error) {
	// created registra a identidade permanente atribuída pelo comando.
	created := ""
	// err preserva o resultado da operação sem publicar estado parcial.
	err := s.Change(func(document *Document) error {
		// page representa a composição atual, resolvida sem alterar o documento original.
		page := document.Screen(pageID)
		if page == nil {
			return fmt.Errorf("tela inexistente")
		}
		// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		component, err := screen.NewComponent(kind, fmt.Sprintf("%s%d", kind, len(page.Components)+1), x, y)
		if err != nil {
			return err
		}
		component.Name = UniqueName(*page, component.Name)
		component.Bounds = component.Bounds.Clamp(page.Size)
		page.Components = append(page.Components, component)
		created = component.ID
		return nil
	})
	return created, err
}

// UniqueName gera um nome sem colisões durante a edição, sem renomear durante a exportação.
func UniqueName(page screen.Screen, base string) string {
	// suffix aumenta até o nome solicitado ficar disponível.
	for suffix := 0; ; suffix++ {
		// candidate contém o valor candidato, preservado enquanto não houver colisão.
		candidate := base
		if suffix > 0 {
			candidate = fmt.Sprintf("%s_%d", base, suffix)
		}
		// exists indica colisão com outra instância.
		exists := false
		// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
		for _, component := range page.Components {
			exists = exists || component.Name == candidate
		}
		if !exists {
			return candidate
		}
	}
}

// DuplicateScreen copia composição, referências internas e navegação para a própria tela.
func (s *Session) DuplicateScreen(id string) (string, error) {
	// created registra a identidade permanente atribuída pelo comando.
	created := ""
	// err preserva o resultado da operação sem publicar estado parcial.
	err := s.Change(func(document *Document) error {
		// original é lido da transação independente.
		original := document.Screen(id)
		if original == nil {
			return fmt.Errorf("tela inexistente")
		}
		// duplicate recebe coleções mutáveis independentes.
		duplicate := *original
		duplicate.Name += " — cópia"
		duplicate.ID = project.NewID()
		duplicate.Components = append([]screen.Component{}, original.Components...)
		// mapping remapeia referências da mesma tela e ações copiadas.
		mapping := map[string]string{}
		// index localiza o item atual na coleção sem alterar sua ordem.
		for index := range duplicate.Components {
			mapping[duplicate.Components[index].ID] = project.NewID()
			duplicate.Components[index].ID = mapping[duplicate.Components[index].ID]
			duplicate.Components[index].Properties = duplicate.Components[index].Properties.Copy()
		}
		// index localiza o item atual na coleção sem alterar sua ordem.
		for index := range duplicate.Components {
			if duplicate.Components[index].Properties.Keyboard != nil {
				duplicate.Components[index].Properties.Keyboard.TargetComponentID = mapping[duplicate.Components[index].Properties.Keyboard.TargetComponentID]
			}
		}
		// action contém uma associação declarativa de botão, preservada salvo remoção de sua origem.
		for _, action := range append([]journey.Action{}, document.Journey.Actions...) {
			// target resolve o destino da operação na coleção ou no mapa corrente.
			// ok informa se a consulta encontrou um valor válido.
			if target, ok := mapping[action.SourceComponentID]; ok {
				action.SourceComponentID = target
				action.SourceScreenID = duplicate.ID
				if action.TargetScreenID == id {
					action.TargetScreenID = duplicate.ID
				}
				document.Journey.Actions = append(document.Journey.Actions, action)
			}
		}
		created = duplicate.ID
		document.Screens = append(document.Screens, duplicate)
		return nil
	})
	return created, err
}

// RemoveScreen remove as ações de saída da tela e preserva entradas como referências diagnosticáveis.
func (s *Session) RemoveScreen(id string) error {
	return s.Change(func(document *Document) error {
		// pages preserva as composições ordenadas e seus campos para validação estrita posterior.
		pages := []screen.Screen{}
		// page representa a composição atual, resolvida sem alterar o documento original.
		for _, page := range document.Screens {
			if page.ID != id {
				pages = append(pages, page)
			}
		}
		// actions exclui as rotas pertencentes é tela removida.
		actions := []journey.Action{}
		// action contém uma associação declarativa de botão, preservada salvo remoção de sua origem.
		for _, action := range document.Journey.Actions {
			if action.SourceScreenID != id {
				actions = append(actions, action)
			}
		}
		document.Screens = pages
		document.Journey.Actions = actions
		return nil
	})
}
