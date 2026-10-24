package application

import (
	"fmt"
	"strings"
	"time"

	"gerador-jornada.local/editor/internal/catalog"
)

// BarcodeRuntime contém somente o estado transitório de uma entrada de código de barras na prévia.
// Nenhum destes valores é persistido no projeto ou enviado aos exportadores.
type BarcodeRuntime struct {
	// Values contém os oito campos visuais na ordem exibida.
	Values [8]string
	// Active identifica o campo que recebe digitação, de zero a sete.
	Active int
	// Locked indica que a entrada foi aceita e está somente para leitura.
	Locked bool
	// Reading indica o modo explícito de captura integral.
	Reading bool
	// Capture acumula a leitura integral até o Enter final.
	Capture string
	// Message contém o erro local da instância, sempre textual.
	Message string
	// Modal indica que o aviso de captura inválida está aberto para esta instância.
	Modal bool
	// Result contém os dados derivados somente depois de uma aceitação.
	Result catalog.Barcode
}

// MenuOutcome classifica o efeito produzido pela ativação de um item de menu.
type MenuOutcome int

const (
	// MenuNoop indica que nenhuma operação foi executada.
	MenuNoop MenuOutcome = iota
	// MenuOpened indica que um submenu habilitado passou a ser o nó visível.
	MenuOpened
	// MenuSimulated indica uma folha sem destino acionável, exibida como seleção simulada.
	MenuSimulated
	// MenuNavigated indica que uma tela de destino autorada foi aberta.
	MenuNavigated
	// MenuBlocked indica que o item está desabilitado e não pode ser ativado.
	MenuBlocked
)

// MenuRuntime contém somente o caminho local e o resultado da última interação com um menu.
// O caminho nunca é gravado no documento ou no histórico do editor.
type MenuRuntime struct {
	// Path contém as identidades dos nós visíveis, da raiz até o submenu atual.
	Path []string
	// SimulatedItemID identifica a última folha acionada sem destino de tela.
	SimulatedItemID string
	// Exited indica que o comando Sair encerrou somente a sessão de prévia.
	Exited bool
}

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
	// Barcodes contém o estado local de cada entrada de código de barras da visita.
	Barcodes map[string]*BarcodeRuntime
	// Menus contém o caminho local de cada menu da aplicação da visita.
	Menus map[string]*MenuRuntime
	// AllowScreenNavigation autoriza folhas de menu a abrirem telas existentes do documento.
	AllowScreenNavigation bool
}

// NewPreview valida a jornada antes de abrir uma sessão isolada.
func NewPreview(document Document) (*Preview, error) {
	// problems reúne os diagnósticos que impedem execução ou exportação.
	if problems := document.Validate(); len(problems) > 0 {
		return nil, fmt.Errorf("%s", ProblemsText(problems))
	}
	// preview conserva seu próprio snapshot e autoriza a navegação entre telas.
	preview := &Preview{Document: document.Copy(), AllowScreenNavigation: true}
	// err registra uma falha ao resolver a tela inicial validada.
	if err := preview.Enter(document.Journey.StartScreenID); err != nil {
		return nil, err
	}
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
	p.Barcodes = map[string]*BarcodeRuntime{}
	p.Menus = map[string]*MenuRuntime{}
	// component fornece os valores iniciais da visita.
	for _, component := range page.Components {
		if component.Properties.Field != nil {
			p.Values[component.ID] = component.Properties.Field.InitialValue
		}
		if component.Properties.Radio != nil {
			p.Choices[component.ID] = component.Properties.Radio.SelectedOptionID
		}
		if component.Properties.BarcodeInput != nil {
			p.Barcodes[component.ID] = &BarcodeRuntime{}
		}
		if component.Properties.Menu != nil {
			p.Menus[component.ID] = &MenuRuntime{Path: initialMenuPath(component.Properties.Menu)}
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

// SetValue valida a intenção de formato declarada sem alegar validação fiscal.
func (p *Preview) SetValue(componentID, value string) error {
	// page identifica a composição corrente.
	page := p.Document.Screen(p.ScreenID)
	if page == nil {
		return fmt.Errorf("tela corrente inexistente")
	}
	// component deve ser um campo da visita atual.
	component := page.ComponentByID(componentID)
	if component == nil || component.Kind != catalog.TextField {
		return fmt.Errorf("campo inexistente")
	}
	if err := validateFieldValue(component.Properties.Field.Format, value); err != nil {
		return err
	}
	p.Values[componentID] = value
	return nil
}

// validateFieldValue aplica a intenção de formato do campo, sem máscara fiscal completa.
func validateFieldValue(format, value string) error {
	if value == "" {
		return nil
	}
	// digits informa se o conteúdo contém somente algarismos ASCII.
	digits := strings.Trim(value, "0123456789") == ""
	switch format {
	case "number":
		if !digits {
			return fmt.Errorf("digite apenas números")
		}
	case "cellular":
		if !digits || (len(value) != 10 && len(value) != 11) {
			return fmt.Errorf("celular deve ter 10 ou 11 dígitos")
		}
	case "cep":
		if !digits || len(value) != 8 {
			return fmt.Errorf("CEP deve ter 8 dígitos")
		}
	case "cpf":
		if !digits || len(value) != 11 {
			return fmt.Errorf("CPF deve ter 11 dígitos")
		}
	case "date":
		if _, err := time.Parse("02/01/2006", value); err != nil {
			return fmt.Errorf("data deve usar dd/mm/aaaa")
		}
	case "ddd":
		if !digits || (len(value) != 2 && len(value) != 3) {
			return fmt.Errorf("DDD deve ter 2 ou 3 dígitos")
		}
	case "nis":
		if !digits || len(value) != 11 {
			return fmt.Errorf("NIS deve ter 11 dígitos")
		}
	}
	return nil
}
