package scenario

import (
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// TestRevisitAndUnicode protege a separação entre visitas e caracteres visíveis.
func TestRevisitAndUnicode(t *testing.T) {
	// page representa uma tela com auto-navegação.
	page := screen.New("A")
	// field e button são os alvos tipados do roteiro.
	field, _ := screen.NewComponent(catalog.TextField, "campo", 0, 0)
	// button contém o botão de teste que origina a navegação.
	button, _ := screen.NewComponent(catalog.Button, "seguir", 0, 100)
	page.Components = []screen.Component{field, button}
	// routes aponta o botão para a própria tela.
	routes := journey.Journey{StartScreenID: page.ID, Actions: []journey.Action{{SourceScreenID: page.ID, SourceComponentID: button.ID, Kind: "navigate", TargetScreenID: page.ID}}}
	// text contém um caractere com marca combinante.
	text := "a\u0301😀"
	// script declara duas visitas sem laço implícito.
	script := Scenario{ID: project.NewID(), Name: "Revisita", StartScreenID: page.ID, Steps: []Step{{ID: project.NewID(), Kind: "typeText", ScreenID: page.ID, ComponentID: field.ID, Text: &text, DurationMS: 200}, {ID: project.NewID(), Kind: "navigate", ScreenID: page.ID, ComponentID: button.ID, DurationMS: 100}, {ID: project.NewID(), Kind: "typeText", ScreenID: page.ID, ComponentID: field.ID, Text: &text, DurationMS: 200}}}
	// visits devem conter estados iniciais independentes.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	visits, err := script.Interpret([]screen.Screen{page}, routes)
	if err != nil {
		t.Fatal(err)
	}
	if len(visits) != 2 || visits[1].Events[0].PreviousText != "" {
		t.Fatal("estado da visita anterior vazou")
	}
	if len(Graphemes(text)) != 2 {
		t.Fatal("graphemas foram separados")
	}
	script.Steps[0].DurationMS = 0
	if _, err = script.Interpret([]screen.Screen{page}, routes); err == nil || !strings.Contains(err.Error(), "STEP_DURATION_INVALID") {
		t.Fatal("duração inválida aceita")
	}
}
