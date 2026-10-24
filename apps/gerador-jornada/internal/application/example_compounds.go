package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/screen"
)

// CompoundsExample cria um projeto de demonstração com os compostos 031, 032 e 033.
// Os valores digitados permanecem vazios; o menu nasce com a árvore inicial documentada.
func CompoundsExample() Document {
	// document reúne duas telas para separar a entrada/ficha do menu de tela inteira.
	document := NewDocument("Teste dos compostos")
	document.Screens = []screen.Screen{screen.New("01 · Códigos de barras"), screen.New("02 · Menu da aplicação")}
	// add insere uma instância nova no índice de tela indicado.
	add := func(page int, kind catalog.Kind, name string, bounds screen.Bounds) string {
		// component contém a identidade e as propriedades padrão do composto.
		component, _ := screen.NewComponent(kind, name, bounds.X, bounds.Y)
		component.Bounds = bounds
		document.Screens[page].Components = append(document.Screens[page].Components, component)
		return component.ID
	}
	add(0, catalog.BarcodeInput, "codigoBarras", screen.Bounds{X: 40, Y: 80, Width: 610, Height: 60})
	add(0, catalog.BarcodeProcessed, "codigoProcessado", screen.Bounds{X: 40, Y: 190, Width: 530, Height: 120})
	// title usa uma skin textual do acervo para demonstrar a personalização.
	titleID := add(0, catalog.Title, "tituloCompostos", screen.Bounds{X: 360, Y: 20, Width: 400, Height: 28})
	title := document.Screens[0].ComponentByID(titleID)
	title.Properties.Title.Text = "Compostos 031 · 032"
	title.Style.Skin = "label_titulo_1"
	// button demonstra a imagem de skin de botão do acervo.
	buttonID := add(0, catalog.Button, "botaoSkin", screen.Bounds{X: 40, Y: 340, Width: 200, Height: 60})
	button := document.Screens[0].ComponentByID(buttonID)
	button.Properties.Button.Text = "Confirmar"
	button.Style.Skin = "botao_confirmar"
	add(1, catalog.ApplicationMenu, "menuAplicacao", screen.Bounds{X: 0, Y: 0, Width: 800, Height: 600})
	// journey mantém a primeira tela como entrada, sem ações implícitas nos compostos.
	document.Journey = journey.Journey{StartScreenID: document.Screens[0].ID, Actions: []journey.Action{}}
	return document
}
