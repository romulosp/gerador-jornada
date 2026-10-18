// Package journey define a navegação explícita entre telas sem depender de eventos da interface.
package journey

import "fmt"

// Action associa um botão a um destino opcional.
type Action struct {
	// SourceScreenID identifica a tela proprietária da origem.
	SourceScreenID string
	// SourceComponentID identifica o botão de origem.
	SourceComponentID string
	// Kind aceita none ou navigate.
	Kind string
	// TargetScreenID permanece vazio somente para ações none.
	TargetScreenID string
}

// Journey declara a tela inicial e as associações de navegação dos botões.
type Journey struct {
	// StartScreenID identifica a visita inicial; pode ficar vazio em jornadas ainda sem tela de entrada.
	StartScreenID string
	// Actions contém no máximo uma associação por botão.
	Actions []Action
}

// Destination resolve o destino pela identidade do componente; retorna vazio quando não há navegação.
func (j Journey) Destination(componentID string) string {
	// action contém uma associação declarativa de botão, preservada salvo remoção de sua origem.
	for _, action := range j.Actions {
		if action.SourceComponentID == componentID && action.Kind == "navigate" {
			return action.TargetScreenID
		}
	}
	return ""
}

// Bind substitui atomicamente a associação de um botão; a aplicação valida referências entre domínios.
func (j *Journey) Bind(action Action) error {
	if action.SourceComponentID == "" || action.SourceScreenID == "" || (action.Kind != "navigate" && action.Kind != "none") || (action.Kind == "navigate" && action.TargetScreenID == "") || (action.Kind == "none" && action.TargetScreenID != "") {
		return fmt.Errorf("NAVIGATION_SOURCE_INVALID: ação inválida")
	}
	// index localiza o item atual na coleção sem alterar sua ordem.
	for index := range j.Actions {
		if j.Actions[index].SourceComponentID == action.SourceComponentID {
			j.Actions[index] = action
			return nil
		}
	}
	j.Actions = append(j.Actions, action)
	return nil
}
