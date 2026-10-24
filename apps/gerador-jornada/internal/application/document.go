// Package application coordena operações de domínio, sessões e trabalhos de exportação.
package application

import (
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/scenario"
	"gerador-jornada.local/editor/internal/screen"
)

// JavaOptions define a configuração de destino Java e o perfil de integração EFL.
type JavaOptions struct {
	// PackageName define o espaço de nomes Java gerado.
	PackageName string
	// ProfileID seleciona um perfil de integração incorporado.
	ProfileID string
}

// Document compõe domínios independentes sem incluir o estado da interface.
type Document struct {
	// Project contém a identidade e os metadados dos recursos sem importar outros domínios.
	Project project.Project
	// Screens contém as composições visuais ordenadas.
	Screens []screen.Screen
	// Journey declara a tela inicial e as associações de navegação dos botões.
	Journey journey.Journey
	// Scenarios contém os roteiros finitos de demonstração.
	Scenarios []scenario.Scenario
	// JavaOptions define a configuração de destino Java e o perfil de integração EFL.
	JavaOptions JavaOptions
	// Resources associa caminhos de recursos incorporados aos seus bytes imutáveis.
	Resources map[string][]byte
}

// Diagnostic localiza um problema de domínio ou exportação que exige ação.
type Diagnostic struct {
	// Code identifica uma categoria estável de problema.
	Code string
	// Severity aceita error ou warning.
	Severity string
	// Message explica o problema em português.
	Message string
	// ScreenID identifica a tela afetada ou visitada e deve corresponder é tela atual na execução.
	ScreenID string
	// ComponentID identifica o controle afetado ou o alvo opcional da ação.
	ComponentID string
	// ScenarioID identifica o roteiro afetado.
	ScenarioID string
}

// NewDocument cria um rascunho com uma tela inicial utilizável.
func NewDocument(name string) Document {
	// page representa a composição atual, resolvida sem alterar o documento original.
	page := screen.New("Tela 1")
	return Document{Project: project.New(name), Screens: []screen.Screen{page}, Journey: journey.Journey{StartScreenID: page.ID, Actions: []journey.Action{}}, Scenarios: []scenario.Scenario{}, Resources: map[string][]byte{}}
}

// Screen retorna a composição mutável pelo identificador permanente, ou nil quando ausente.
func (d *Document) Screen(id string) *screen.Screen {
	// index localiza o item atual na coleção sem alterar sua ordem.
	for index := range d.Screens {
		if d.Screens[index].ID == id {
			return &d.Screens[index]
		}
	}
	return nil
}

// Copy cria uma cópia independente, separando todos os dados mutáveis da origem.
func (d Document) Copy() Document {
	// result preserva os valores imutáveis e separa as coleções mutáveis da origem.
	result := d
	result.Project.Assets = append([]project.Asset{}, d.Project.Assets...)
	result.Screens = append([]screen.Screen{}, d.Screens...)
	// index localiza o item atual na coleção sem alterar sua ordem.
	for index := range result.Screens {
		result.Screens[index].Components = append([]screen.Component{}, d.Screens[index].Components...)
		// componentIndex localiza uma variante discriminada de propriedades.
		for componentIndex := range result.Screens[index].Components {
			result.Screens[index].Components[componentIndex].Properties = d.Screens[index].Components[componentIndex].Properties.Copy()
		}
	}
	result.Journey.Actions = append([]journey.Action{}, d.Journey.Actions...)
	result.Scenarios = append([]scenario.Scenario{}, d.Scenarios...)
	// index localiza o item atual na coleção sem alterar sua ordem.
	for index := range result.Scenarios {
		result.Scenarios[index].Steps = append([]scenario.Step{}, d.Scenarios[index].Steps...)
		// stepIndex localiza ponteiros opcionais que exigem cópias independentes.
		for stepIndex := range result.Scenarios[index].Steps {
			if d.Scenarios[index].Steps[stepIndex].Point != nil {
				// v separa o valor opcional para impedir alterações na origem.
				v := *d.Scenarios[index].Steps[stepIndex].Point
				result.Scenarios[index].Steps[stepIndex].Point = &v
			}
			if d.Scenarios[index].Steps[stepIndex].Text != nil {
				// v separa o valor opcional para impedir alterações na origem.
				v := *d.Scenarios[index].Steps[stepIndex].Text
				result.Scenarios[index].Steps[stepIndex].Text = &v
			}
		}
	}
	result.Resources = map[string][]byte{}
	// name identifica o caminho relativo da entrada ou do recurso incorporado imutável.
	// data contém os bytes imutáveis do recurso incorporado.
	for name, data := range d.Resources {
		result.Resources[name] = append([]byte{}, data...)
	}
	return result
}

// Validate resolve referências entre domínios e reúne diagnósticos para execução e exportação.
func (d Document) Validate() []Diagnostic {
	// problems acumula erros em português sem descartar referências inválidas.
	problems := []Diagnostic{}
	// add mantém consistente a construção dos diagnósticos.
	add := func(code, message, page, component string) {
		problems = append(problems, Diagnostic{Code: code, Severity: "error", Message: message, ScreenID: page, ComponentID: component})
	}
	if err := d.ValidateFileIDs(); err != nil {
		add("FILE_ID_INVALID", err.Error(), "", "")
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := d.Project.Validate(); err != nil {
		add("PROJECT_SCHEMA_INVALID", err.Error(), "", "")
	}
	// seen detecta identidades ou nomes repetidos, inclusive diferenças de caixa em caminhos Windows.
	seen := map[string]bool{d.Project.ID: true}
	// register rejeita identidades duplicadas ou malformadas.
	register := func(id, page, component string) {
		if !project.ValidID(id) || seen[id] {
			add("PROJECT_SCHEMA_INVALID", "Identificador inválido ou repetido", page, component)
		}
		seen[id] = true
	}
	// assetIDs identifica recursos declarados com bytes portáteis no documento.
	assetIDs := map[string]bool{}
	// asset percorre os metadados do projeto.
	for _, asset := range d.Project.Assets {
		assetIDs[asset.ID] = assetIDs[asset.ID] || len(d.Resources[asset.Path]) > 0
	}
	// page representa a composição atual, resolvida sem alterar o documento original.
	for _, page := range d.Screens {
		register(page.ID, page.ID, "")
		if page.Background.AssetID != "" && !assetIDs[page.Background.AssetID] {
			add("RESOURCE_MISSING", "Imagem de fundo ausente do pacote", page.ID, "")
		}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if err := page.Validate(); err != nil {
			add("SCREEN_INVALID", err.Error(), page.ID, "")
		}
		// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
		for _, component := range page.Components {
			register(component.ID, page.ID, component.ID)
			if component.Style.BackgroundAssetID != "" && !assetIDs[component.Style.BackgroundAssetID] {
				add("RESOURCE_MISSING", "Imagem decorativa ausente do pacote", page.ID, component.ID)
			}
			// identidades internas de compostos também são únicas em todo o documento.
			for _, internalID := range component.Properties.CompoundIDs() {
				register(internalID, page.ID, component.ID)
			}
			if component.Properties.Keyboard != nil && component.Properties.Keyboard.TargetComponentID != "" {
				// target deve identificar um campo editável ou uma entrada de código de barras na mesma tela.
				target := page.ComponentByID(component.Properties.Keyboard.TargetComponentID)
				if target == nil || (target.Kind != catalog.TextField && target.Kind != catalog.BarcodeInput) {
					add("KEYBOARD_TARGET_INVALID", "Teclado deve apontar para um campo ou entrada da mesma tela", page.ID, component.ID)
				}
			}
			if component.Properties.Menu != nil {
				// walk verifica destinos explícitos inclusive em ramos desabilitados.
				_ = component.Properties.Menu.WalkItems(func(item *catalog.MenuItem) error {
					if item.TargetScreenID != "" && d.Screen(item.TargetScreenID) == nil {
						add("MENU_TARGET_MISSING", "Destino de menu inexistente", page.ID, component.ID)
					}
					return nil
				})
			}
		}
	}
	if d.Screen(d.Journey.StartScreenID) == nil {
		add("JOURNEY_START_MISSING", "Defina uma tela inicial existente", "", "")
	}
	// bindings impede associações de clique duplicadas.
	bindings := map[string]bool{}
	// action contém uma associação declarativa de botão, preservada salvo remoção de sua origem.
	for _, action := range d.Journey.Actions {
		// page representa a composição atual, resolvida sem alterar o documento original.
		page := d.Screen(action.SourceScreenID)
		if page == nil {
			add("NAVIGATION_SOURCE_INVALID", "Tela de origem inexistente", action.SourceScreenID, action.SourceComponentID)
			continue
		}
		// button identifica a origem declarada da navegação.
		button := page.ComponentByID(action.SourceComponentID)
		if button == nil || button.Kind != catalog.Button || bindings[action.SourceComponentID] {
			add("NAVIGATION_SOURCE_INVALID", "Origem deve ser um botão com ação única", page.ID, action.SourceComponentID)
		}
		bindings[action.SourceComponentID] = true
		if action.Kind == "navigate" {
			if d.Screen(action.TargetScreenID) == nil {
				add("NAVIGATION_TARGET_MISSING", "Destino de navegação inexistente", page.ID, action.SourceComponentID)
			}
		} else if action.Kind != "none" || action.TargetScreenID != "" {
			add("NAVIGATION_SOURCE_INVALID", "Ação inválida", page.ID, action.SourceComponentID)
		}
	}
	return problems
}

// ProblemsText formata diagnósticos em português para a linha de comando e os diálogos desktop.
func ProblemsText(problems []Diagnostic) string {
	// lines contém o resumo exibido ao usuário sem os conteúdos dos campos.
	lines := []string{}
	// problem contém um diagnóstico que exige ação.
	for _, problem := range problems {
		lines = append(lines, fmt.Sprintf("%s: %s [%s %s]", problem.Code, problem.Message, problem.ScreenID, problem.ComponentID))
	}
	return strings.Join(lines, "\n")
}
