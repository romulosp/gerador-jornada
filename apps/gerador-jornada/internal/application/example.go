package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/scenario"
	"gerador-jornada.local/editor/internal/screen"
)

// Example fornece uma jornada inicial editável que demonstra as sete famílias.
func Example() Document {
	// document é um projeto independente, sem dados de usuário.
	document := NewDocument("Jornada de atendimento")
	document.Screens = []screen.Screen{screen.New("01 · Identificação"), screen.New("02 · Confirmação"), screen.New("03 · Conclusão")}
	// add aplica defaults e geometria explícita sem criar dependência de UI.
	add := func(page int, kind catalog.Kind, name string, bounds screen.Bounds, text string) string {
		// component recebe identidade e propriedades tipadas.
		component, _ := screen.NewComponent(kind, name, bounds.X, bounds.Y)
		component.Bounds = bounds
		if component.Properties.Title != nil {
			component.Properties.Title.Text = text
		}
		if component.Properties.Label != nil {
			component.Properties.Label.Text = text
		}
		if component.Properties.Button != nil {
			component.Properties.Button.Text = text
		}
		document.Screens[page].Components = append(document.Screens[page].Components, component)
		return component.ID
	}
	add(0, catalog.Title, "tituloIdentificacao", screen.Bounds{X: 48, Y: 48, Width: 680, Height: 40}, "Vamos começar seu atendimento")
	add(0, catalog.Label, "instrucao", screen.Bounds{X: 48, Y: 108, Width: 700, Height: 26}, "Informe seu número para seguir com a jornada.")
	// field identifica o alvo do teclado numérico.
	field := add(0, catalog.TextField, "numeroAtendimento", screen.Bounds{X: 48, Y: 180, Width: 350, Height: 44}, "")
	document.Screens[0].ComponentByID(field).Properties.Field.Format = "number"
	// keyboard demonstra o vínculo entre controles.
	keyboard := add(0, catalog.NumericKeyboard, "teclado", screen.Bounds{X: 500, Y: 170, Width: 230, Height: 270}, "")
	document.Screens[0].ComponentByID(keyboard).Properties.Keyboard.TargetComponentID = field
	// next é a ação de entrada na tela de confirmação.
	next := add(0, catalog.Button, "continuar", screen.Bounds{X: 48, Y: 258, Width: 190, Height: 48}, "Continuar")
	add(1, catalog.Title, "tituloConfirmacao", screen.Bounds{X: 48, Y: 48, Width: 700, Height: 40}, "Confira os detalhes")
	add(1, catalog.Label, "instrucaoConfirmacao", screen.Bounds{X: 48, Y: 108, Width: 700, Height: 26}, "Escolha uma opção e confirme o atendimento.")
	// report demonstra dados visuais sem regras de negócio externas.
	report := add(1, catalog.Report, "resumo", screen.Bounds{X: 48, Y: 164, Width: 380, Height: 260}, "")
	document.Screens[1].ComponentByID(report).Properties.Report = &catalog.ReportProperties{Title: "Resumo do atendimento", Rows: []catalog.Option{{ID: project.NewID(), Text: "Serviço: atendimento presencial"}, {ID: project.NewID(), Text: "Unidade: agência de demonstração"}}}
	add(1, catalog.Radio, "opcoes", screen.Bounds{X: 470, Y: 164, Width: 260, Height: 88}, "")
	// confirm avança para o encerramento.
	confirm := add(1, catalog.Button, "confirmar", screen.Bounds{X: 48, Y: 460, Width: 190, Height: 48}, "Confirmar")
	add(2, catalog.Title, "tituloConclusao", screen.Bounds{X: 48, Y: 100, Width: 700, Height: 42}, "Tudo certo. Atendimento confirmado!")
	add(2, catalog.Label, "mensagemFinal", screen.Bounds{X: 48, Y: 170, Width: 700, Height: 28}, "Você concluiu uma jornada com três telas conectadas.")
	// restart permite demonstrar uma revisita à primeira tela.
	restart := add(2, catalog.Button, "reiniciar", screen.Bounds{X: 48, Y: 250, Width: 210, Height: 48}, "Novo atendimento")
	document.Journey.StartScreenID = document.Screens[0].ID
	document.Journey.Actions = []journey.Action{{SourceScreenID: document.Screens[0].ID, SourceComponentID: next, Kind: "navigate", TargetScreenID: document.Screens[1].ID}, {SourceScreenID: document.Screens[1].ID, SourceComponentID: confirm, Kind: "navigate", TargetScreenID: document.Screens[2].ID}, {SourceScreenID: document.Screens[2].ID, SourceComponentID: restart, Kind: "navigate", TargetScreenID: document.Screens[0].ID}}
	// text é o preenchimento demonstrativo, sem alterar o valor inicial do campo.
	text := "123456"
	// script contém etapas explícitas, incluindo a revisita final.
	script := scenario.Scenario{ID: project.NewID(), Name: "Atendimento completo", StartScreenID: document.Screens[0].ID, Steps: []scenario.Step{}}
	// appendStep mantém IDs únicos e nulabilidade apropriada.
	appendStep := func(kind, page, component string, duration int, value *string) {
		script.Steps = append(script.Steps, scenario.Step{ID: project.NewID(), Kind: kind, ScreenID: page, ComponentID: component, DurationMS: duration, Text: value})
	}
	appendStep("moveCursor", document.Screens[0].ID, field, 600, nil)
	appendStep("click", document.Screens[0].ID, field, 300, nil)
	appendStep("typeText", document.Screens[0].ID, field, 1500, &text)
	// route descreve a sequência demonstrativa de cada botão.
	for _, route := range document.Journey.Actions {
		appendStep("wait", route.SourceScreenID, "", 700, nil)
		appendStep("moveCursor", route.SourceScreenID, route.SourceComponentID, 500, nil)
		appendStep("click", route.SourceScreenID, route.SourceComponentID, 300, nil)
		appendStep("navigate", route.SourceScreenID, route.SourceComponentID, 200, nil)
	}
	document.Scenarios = []scenario.Scenario{script}
	return document
}
