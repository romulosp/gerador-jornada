package application

import (
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
)

// Preview contém estado de execução isolado da edição.
type Preview struct {
	// Document é uma cópia imutável do projeto executado.
	Document Document
	// ScreenID identifica a tela corrente.
	ScreenID string
	// Values armazena entradas somente da visita atual.
	Values map[string]string
	// Choices armazena a seleção de opções durante a visita.
	Choices map[string]string
}

// NewPreview valida a jornada antes de abrir uma sessão isolada.
func NewPreview(document Document) (*Preview, error) {
	// problems reúne os diagnósticos que impedem execução ou exportação.
	if problems := document.Validate(); len(problems) > 0 {
		return nil, fmt.Errorf("%s", ProblemsText(problems))
	}
	// preview conserva seu próprio snapshot.
	preview := &Preview{Document: document.Copy()}
	preview.Enter(document.Journey.StartScreenID)
	return preview, nil
}

// Enter reinicializa valores a cada visita, incluindo retorno à própria tela.
func (p *Preview) Enter(id string) error {
	// page é a composição de destino.
	page := p.Document.Screen(id)
	if page == nil {
		return fmt.Errorf("tela de destino inexistente")
	}
	p.ScreenID = id
	p.Values = map[string]string{}
	p.Choices = map[string]string{}
	// component fornece os valores iniciais.
	for _, component := range page.Components {
		if component.Properties.Field != nil {
			p.Values[component.ID] = component.Properties.Field.InitialValue
		}
		if component.Properties.Radio != nil {
			p.Choices[component.ID] = component.Properties.Radio.SelectedOptionID
		}
	}
	return nil
}

// Click executa exclusivamente a navegação de um botão.
func (p *Preview) Click(componentID string) error {
	// target resolve a ação sem depender de nome ou posição.
	target := p.Document.Journey.Destination(componentID)
	if target == "" {
		return nil
	}
	return p.Enter(target)
}

// SetValue valida a entrada numérica sem alegar validação fiscal das demais máscaras.
func (p *Preview) SetValue(componentID, value string) error {
	// page identifica a composição corrente.
	page := p.Document.Screen(p.ScreenID)
	// component deve ser um campo da visita atual.
	component := page.ComponentByID(componentID)
	if component == nil || component.Kind != catalog.TextField {
		return fmt.Errorf("campo inexistente")
	}
	if component.Properties.Field.Format == "number" && strings.Trim(value, "0123456789") != "" {
		return fmt.Errorf("digite apenas números")
	}
	p.Values[componentID] = value
	return nil
}
