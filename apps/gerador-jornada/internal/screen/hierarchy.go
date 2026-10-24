package screen

import "fmt"

// Children devolve os componentes cujo ParentID é o identificador indicado, na ordem de camadas.
func (s Screen) Children(id string) []string {
	// children preserva a ordem da lista plana.
	children := []string{}
	// component percorre os componentes da tela.
	for _, component := range s.Components {
		if component.ParentID == id {
			children = append(children, component.ID)
		}
	}
	return children
}

// Descendants devolve a subárvore do componente, sem incluir o próprio identificador.
// A travessia é limitada pela quantidade de componentes e não entra em ciclo.
func (s Screen) Descendants(id string) []string {
	// result preserva a ordem de descoberta em largura.
	result := []string{}
	// seen protege contra ciclo em uma composição ainda não validada.
	seen := map[string]bool{id: true}
	// queue contém os pais cujos filhos ainda serão visitados.
	queue := []string{id}
	// index percorre a fila até esgotar a subárvore.
	for index := 0; index < len(queue) && len(result) < len(s.Components); index++ {
		// child é um filho direto do nó atual.
		for _, child := range s.Children(queue[index]) {
			if seen[child] {
				continue
			}
			seen[child] = true
			result = append(result, child)
			queue = append(queue, child)
		}
	}
	return result
}

// IsAncestor informa se ancestor pertence à cadeia de pais do componente indicado.
func (s Screen) IsAncestor(ancestor, id string) bool {
	// component é o nó cuja cadeia será percorrida.
	component := s.ComponentByID(id)
	// seen protege contra ciclo.
	seen := map[string]bool{}
	for component != nil && component.ParentID != "" {
		if seen[component.ParentID] {
			return false
		}
		seen[component.ParentID] = true
		if component.ParentID == ancestor {
			return true
		}
		component = s.ComponentByID(component.ParentID)
	}
	return false
}

// ValidateHierarchy verifica pai existente na mesma tela, autorreferência e ciclos.
func (s Screen) ValidateHierarchy() error {
	// component percorre os nós para validar cada vínculo.
	for _, component := range s.Components {
		if component.ParentID == "" {
			continue
		}
		if component.ParentID == component.ID {
			return fmt.Errorf("%s: componente não pode ser pai de si mesmo", component.Name)
		}
		if s.ComponentByID(component.ParentID) == nil {
			return fmt.Errorf("%s: pai inexistente na mesma tela", component.Name)
		}
		// current percorre a cadeia de ancestrais em busca de ciclo.
		current := s.ComponentByID(component.ParentID)
		// seen limita a travessia ao tamanho da tela.
		seen := map[string]bool{component.ID: true}
		for current != nil {
			if seen[current.ID] {
				return fmt.Errorf("%s: ciclo de hierarquia", component.Name)
			}
			seen[current.ID] = true
			if current.ParentID == "" {
				break
			}
			current = s.ComponentByID(current.ParentID)
		}
	}
	return nil
}

// OrphanChildren limpa o ParentID dos filhos diretos de um componente removido.
// A operação mantém os filhos na mesma tela e não altera sua geometria.
func (s *Screen) OrphanChildren(id string) {
	// index percorre os componentes mutáveis.
	for index := range s.Components {
		if s.Components[index].ParentID == id {
			s.Components[index].ParentID = ""
		}
	}
}
