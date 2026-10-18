// Package scenario interpreta etapas finitas de demonstração sem depender de interface ou XML.
package scenario

import (
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/screen"
	"github.com/rivo/uniseg"
)

// Point localiza o cursor da demonstração em unidades lógicas; nas etapas, é opcional.
type Point struct {
	// X expressa a coordenada horizontal em unidades lógicas.
	X float64
	// Y expressa a coordenada vertical em unidades lógicas.
	Y float64
}

// Step contém a ação declarativa original com temporização explícita.
type Step struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Kind aceita moveCursor, click, typeText, wait ou navigate.
	Kind string
	// ScreenID identifica a tela afetada ou visitada e deve corresponder é tela atual na execução.
	ScreenID string
	// ComponentID identifica o controle afetado ou o alvo opcional da ação.
	ComponentID string
	// Point localiza o cursor da demonstração em unidades lógicas; nas etapas, é opcional.
	Point *Point
	// Text contém o texto Unicode exibido; typeText permite inclusive uma string vazia.
	Text *string
	// DurationMS expressa a duração em milissegundos; etapas aceitam de 1 a 600000.
	DurationMS int
}

// Scenario contém uma demonstração finita e ordenada.
type Scenario struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Name contém o nome exibido ao usuário; nomes de componentes são únicos na tela.
	Name string
	// StartScreenID identifica a visita inicial; pode ficar vazio em jornadas ainda sem tela de entrada.
	StartScreenID string
	// Steps contém as etapas executadas em ordem.
	Steps []Step
}

// Event acrescenta o contexto de execução a uma etapa interpretada.
type Event struct {
	// Step contém a ação declarativa original com temporização explícita.
	Step Step
	// StartMS expressa o instante inicial em milissegundos relativos é ocorrência.
	StartMS int
	// From registra a posição do cursor antes da ação.
	From Point
	// To registra a posição resultante do cursor.
	To Point
	// PreviousText registra o conteúdo do campo antes da substituição.
	PreviousText string
}

// Occurrence representa uma visita distinta, incluindo novas visitas é mesma tela.
type Occurrence struct {
	// ScreenID identifica a tela afetada ou visitada e deve corresponder é tela atual na execução.
	ScreenID string
	// Events contém os eventos relativos é visita atual.
	Events []Event
	// DurationMS soma a duração de todas as etapas desta ocorrência.
	DurationMS int
}

// Graphemes divide o texto visível sem separar caracteres combinantes.
func Graphemes(value string) []string {
	// result preserva os valores imutáveis e separa as coleções mutáveis da origem.
	result := []string{}
	// iterator utiliza os limites Unicode dos agrupamentos de grafemas.
	iterator := uniseg.NewGraphemes(value)
	for iterator.Next() {
		result = append(result, iterator.Str())
	}
	return result
}

// Interpret valida e resolve todo o roteiro antes de produzir eventos.
func (s Scenario) Interpret(screens []screen.Screen, routes journey.Journey) ([]Occurrence, error) {
	if len(s.Steps) == 0 {
		return nil, fmt.Errorf("SCENARIO_EMPTY: %s", s.Name)
	}
	// index localiza o item atual na coleção sem alterar sua ordem.
	index := map[string]screen.Screen{}
	// page representa a composição atual, resolvida sem alterar o documento original.
	for _, page := range screens {
		index[page.ID] = page
	}
	// ok informa se a consulta encontrou um valor válido.
	if _, ok := index[s.StartScreenID]; !ok {
		return nil, fmt.Errorf("JOURNEY_START_MISSING: %s", s.Name)
	}
	// visits começa com uma ocorrência inicial.
	visits := []Occurrence{{ScreenID: s.StartScreenID, Events: []Event{}}}
	// cursor é reiniciado a cada navegação.
	cursor := Point{}
	// values armazena os valores de campos somente da ocorrência atual.
	values := map[string]string{}
	// seen detecta identidades ou nomes repetidos, inclusive diferenças de caixa em caminhos Windows.
	seen := map[string]bool{}
	// step é interpretado somente após a conclusão da etapa anterior.
	for _, step := range s.Steps {
		// visit aponta para a ocorrência atual na saída.
		visit := &visits[len(visits)-1]
		if step.ID == "" || seen[step.ID] {
			return nil, fmt.Errorf("STEP_KIND_INVALID: identificador repetido")
		}
		seen[step.ID] = true
		if step.ScreenID != visit.ScreenID {
			return nil, fmt.Errorf("STEP_SCREEN_MISMATCH: %s", step.ID)
		}
		if step.DurationMS < 1 || step.DurationMS > 600000 {
			return nil, fmt.Errorf("STEP_DURATION_INVALID: %s", step.ID)
		}
		// page representa a composição atual, resolvida sem alterar o documento original.
		page := index[visit.ScreenID]
		// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
		component := page.ComponentByID(step.ComponentID)
		// event registra o estado do cursor e do texto antes da ação.
		event := Event{Step: step, StartMS: visit.DurationMS, From: cursor, To: cursor}
		switch step.Kind {
		case "moveCursor":
			if (step.Point == nil) == (component == nil) || step.Text != nil {
				return nil, fmt.Errorf("STEP_TARGET_MISSING: %s", step.ID)
			}
			if step.Point != nil {
				event.To = *step.Point
			} else {
				event.To = center(component.Bounds)
			}
			if event.To.X < 0 || event.To.Y < 0 || event.To.X > float64(page.Size.Width) || event.To.Y > float64(page.Size.Height) {
				return nil, fmt.Errorf("STEP_VALUE_INVALID: %s", step.ID)
			}
		case "click":
			if component == nil || (component.Kind != catalog.Button && component.Kind != catalog.TextField) || step.Point != nil || step.Text != nil {
				return nil, fmt.Errorf("STEP_TARGET_MISSING: %s", step.ID)
			}
			event.To = center(component.Bounds)
		case "typeText":
			if component == nil || component.Properties.Field == nil || step.Text == nil || step.Point != nil {
				return nil, fmt.Errorf("STEP_VALUE_INVALID: %s", step.ID)
			}
			if component.Properties.Field.Format == "number" && strings.Trim(*step.Text, "0123456789") != "" {
				return nil, fmt.Errorf("STEP_VALUE_INVALID: somente dígitos")
			}
			// previous contém o valor inicial ou o último valor digitado durante a execução.
			// ok informa se a consulta encontrou um valor válido.
			previous, ok := values[component.ID]
			if !ok {
				previous = component.Properties.Field.InitialValue
			}
			event.PreviousText = previous
			values[component.ID] = *step.Text
		case "wait":
			if step.ComponentID != "" || step.Point != nil || step.Text != nil {
				return nil, fmt.Errorf("STEP_VALUE_INVALID: %s", step.ID)
			}
		case "navigate":
			if component == nil || component.Kind != catalog.Button || step.Point != nil || step.Text != nil {
				return nil, fmt.Errorf("STEP_TARGET_MISSING: %s", step.ID)
			}
			// ok informa se a consulta encontrou um valor válido.
			if _, ok := index[routes.Destination(component.ID)]; !ok {
				return nil, fmt.Errorf("NAVIGATION_TARGET_MISSING: %s", step.ID)
			}
		default:
			return nil, fmt.Errorf("STEP_KIND_INVALID: %s", step.ID)
		}
		visit.Events = append(visit.Events, event)
		visit.DurationMS += step.DurationMS
		cursor = event.To
		if step.Kind == "navigate" {
			visits = append(visits, Occurrence{ScreenID: routes.Destination(component.ID), Events: []Event{}})
			cursor = Point{}
			values = map[string]string{}
		}
	}
	return visits, nil
}

// center retorna o centro lógico do componente.
func center(bounds screen.Bounds) Point {
	return Point{float64(bounds.X) + float64(bounds.Width)/2, float64(bounds.Y) + float64(bounds.Height)/2}
}
