package desktop

import (
	"fmt"

	"gerador-jornada.local/editor/internal/screen"
)

// parentOptions monta as opções de contêiner da mesma tela, excluindo o próprio nó e seus descendentes.
func parentOptions(page screen.Screen, component screen.Component) ([]string, map[string]string) {
	// options começa pela ausência de pai.
	options := []string{"Nenhum"}
	// labels associa o texto exibido à identidade estável.
	labels := map[string]string{"Nenhum": ""}
	// index percorre a ordem canônica para numerar os candidatos.
	for index, candidate := range page.Components {
		if candidate.ID == component.ID || page.IsAncestor(component.ID, candidate.ID) {
			continue
		}
		// label usa o ordinal e o nome, sem depender de nome único.
		label := fmt.Sprintf("%d · %s", index+1, candidate.Name)
		options = append(options, label)
		labels[label] = candidate.ID
	}
	return options, labels
}

// layerOptions monta as referências possíveis para posicionamento relativo de camada.
func layerOptions(page screen.Screen, component screen.Component) ([]string, map[string]string) {
	// options começa pela ausência de referência.
	options := []string{"Selecione um objeto"}
	// labels associa o texto exibido à identidade estável.
	labels := map[string]string{"Selecione um objeto": ""}
	// index percorre a ordem canônica para numerar os demais objetos.
	for index, candidate := range page.Components {
		if candidate.ID == component.ID {
			continue
		}
		// label usa o ordinal e o nome, sem depender de nome único.
		label := fmt.Sprintf("%d · %s", index+1, candidate.Name)
		options = append(options, label)
		labels[label] = candidate.ID
	}
	return options, labels
}

// currentLayer devolve a posição canônica do componente, começando em um.
func currentLayer(page screen.Screen, componentID string) (int, int) {
	// index percorre a ordem fundo→frente.
	for index, candidate := range page.Components {
		if candidate.ID == componentID {
			return index + 1, len(page.Components)
		}
	}
	return 0, len(page.Components)
}
