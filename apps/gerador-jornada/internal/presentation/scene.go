// Package presentation transforma componentes em primitivas independentes de Fyne e XML.
package presentation

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// Primitive descreve uma forma, texto ou recurso visual editável.
type Primitive struct {
	// Bounds informa a geometria absoluta em unidades lógicas.
	Bounds screen.Bounds
	// Text contém o texto da forma; vazio não elimina a caixa.
	Text string
	// Fill determina o preenchimento RGB, vazio para transparente.
	Fill string
	// Stroke determina a borda RGB opcional.
	Stroke string
	// Color determina a cor do texto.
	Color string
	// FontSize determina o tamanho lógico do texto.
	FontSize int
	// Bold destaca títulos e cabeçalhos.
	Bold bool
	// AssetID referencia uma imagem decorativa.
	AssetID string
	// Role identifica o valor animável de um campo.
	Role string
	// TextAlign determina o alinhamento do parágrafo sem mover os limites da forma.
	TextAlign string
}

// Group reúne primitivas de um componente.
type Group struct {
	// ComponentID vincula o grupo ao componente original.
	ComponentID string
	// Kind preserva a família do catálogo.
	Kind catalog.Kind
	// Bounds informa o limite lógico do grupo.
	Bounds screen.Bounds
	// Shapes mantém a ordem visual interna.
	Shapes []Primitive
}

// Build expande uma tela sem consultar widgets ou alterar o projeto.
func Build(page screen.Screen) []Group {
	// groups preserva a ordem de camadas.
	groups := []Group{}
	// component é uma instância da composição.
	for _, component := range page.Components {
		groups = append(groups, Component(component))
	}
	return groups
}

// Component expande as sete famílias em texto e formas editáveis.
func Component(component screen.Component) Group {
	// group preserva identidade e geometria.
	group := Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds, Shapes: []Primitive{}}
	// box representa a área absoluta do componente.
	box := component.Bounds
	// text cria uma caixa tipográfica sem fundo.
	text := func(value string, bounds screen.Bounds, bold bool) Primitive {
		// size resolve o padrão e o override de fonte.
		size := component.Style.FontSize
		if size == 0 {
			size = 12
		}
		// color resolve a cor padrão sem modificar propriedades.
		color := component.Style.TextColor
		if color == "" {
			color = "#4864A0"
		}
		return Primitive{Bounds: bounds, Text: value, Color: color, FontSize: size, Bold: bold, TextAlign: component.Style.TextAlign}
	}
	// rect cria a forma base do controle.
	rect := func(fill string) Primitive {
		if component.Style.FillColor != "" {
			fill = component.Style.FillColor
		}
		// border resolve o contorno padrão.
		border := component.Style.BorderColor
		if border == "" {
			border = "#CBD5E1"
		}
		return Primitive{Bounds: box, Fill: fill, Stroke: border}
	}
	switch component.Kind {
	case catalog.Button:
		group.Shapes = append(group.Shapes, rect("#2563EB"))
		// caption mantém a legenda editável sobre a forma.
		caption := text(component.Properties.Button.Text, inset(box, 10, 10), true)
		if component.Style.TextColor == "" {
			caption.Color = "#FFFFFF"
		}
		group.Shapes = append(group.Shapes, caption)
	case catalog.Label:
		group.Shapes = append(group.Shapes, text(component.Properties.Label.Text, box, false))
	case catalog.Title:
		// title usa destaque sem exigir uma skin externa.
		title := text(component.Properties.Title.Text, box, true)
		if component.Style.FontSize == 0 {
			title.FontSize = 22
		}
		group.Shapes = append(group.Shapes, title)
	case catalog.TextField:
		group.Shapes = append(group.Shapes, rect("#FFFFFF"))
		// value é o alvo estável para animação de digitação.
		value := text(component.Properties.Field.InitialValue, inset(box, 6, 4), false)
		value.Role = "value"
		group.Shapes = append(group.Shapes, value)
	case catalog.Radio:
		// option distribui opções sem rasterização do grupo.
		// index localiza o item atual ou sua posição de destino na coleção ordenada.
		for index, option := range component.Properties.Radio.Options {
			// bounds posiciona a opção conforme orientação.
			bounds := screen.Bounds{X: box.X, Y: box.Y + index*24, Width: box.Width, Height: 22}
			if component.Properties.Radio.Orientation == "horizontal" {
				bounds.X = box.X + index*box.Width/len(component.Properties.Radio.Options)
				bounds.Y = box.Y
				bounds.Width = box.Width / len(component.Properties.Radio.Options)
			}
			// marker mostra o estado inicial da seleção.
			marker := "○ "
			if option.ID == component.Properties.Radio.SelectedOptionID {
				marker = "● "
			}
			group.Shapes = append(group.Shapes, text(marker+option.Text, bounds, false))
		}
	case catalog.Report:
		group.Shapes = append(group.Shapes, rect("#FFFFFF"))
		// header separa título e linhas do relatório.
		header := Primitive{Bounds: screen.Bounds{X: box.X, Y: box.Y, Width: box.Width, Height: 28}, Fill: "#EFF6FF"}
		group.Shapes = append(group.Shapes, header, text(component.Properties.Report.Title, inset(header.Bounds, 8, 5), true))
		// row representa uma linha ordenada e editável.
		// index localiza o item atual ou sua posição de destino na coleção ordenada.
		for index, row := range component.Properties.Report.Rows {
			group.Shapes = append(group.Shapes, text(row.Text, screen.Bounds{X: box.X + 8, Y: box.Y + 34 + index*24, Width: box.Width - 16, Height: 22}, false))
		}
	case catalog.NumericKeyboard:
		// keys mantém o arranjo previsível do teclado.
		keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "Limpar", "0", "⌫"}
		// index determina linha e coluna de cada tecla.
		// key identifica a tecla correspondente à posição atual do teclado.
		for index, key := range keys {
			// bounds reserva espaçamento entre teclas.
			bounds := screen.Bounds{X: box.X + (index%3)*box.Width/3 + 2, Y: box.Y + (index/3)*box.Height/4 + 2, Width: box.Width/3 - 4, Height: box.Height/4 - 4}
			group.Shapes = append(group.Shapes, Primitive{Bounds: bounds, Fill: "#F1F5F9", Stroke: "#CBD5E1"}, text(key, inset(bounds, 8, 8), true))
		}
	}
	if component.Style.BackgroundAssetID != "" {
		group.Shapes = append([]Primitive{{Bounds: box, AssetID: component.Style.BackgroundAssetID}}, group.Shapes...)
	}
	return group
}

// inset reserva margens internas sem introduzir dimensões negativas.
func inset(bounds screen.Bounds, x, y int) screen.Bounds {
	return screen.Bounds{X: bounds.X + x, Y: bounds.Y + y, Width: max(1, bounds.Width-2*x), Height: max(1, bounds.Height-2*y)}
}

// Caption apresenta um resumo textual para listas de seleção.
func Caption(component screen.Component) string {
	switch component.Kind {
	case catalog.Button:
		return component.Properties.Button.Text
	case catalog.Label:
		return component.Properties.Label.Text
	case catalog.Title:
		return component.Properties.Title.Text
	case catalog.TextField:
		return component.Properties.Field.InitialValue
	case catalog.Report:
		return component.Properties.Report.Title
	}
	return fmt.Sprintf("%s", component.Kind)
}
