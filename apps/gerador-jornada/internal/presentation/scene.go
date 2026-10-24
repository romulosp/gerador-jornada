// Package presentation transforma componentes em primitivas independentes de Fyne e XML.
package presentation

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
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
	// Italic identifica a face física em itálico, sem inferência pelo atributo XML.
	Italic bool
	// LineSpacing mantém o fator original de distância entre linhas.
	LineSpacing float64
	// AssetID referencia uma imagem decorativa.
	AssetID string
	// Role identifica o valor animável de um campo.
	Role string
	// TextAlign determina o alinhamento do parágrafo sem mover os limites da forma.
	TextAlign string
	// Name identifica a forma interna para edição e diagnóstico no PowerPoint.
	Name string
	// Wrap permite quebra de linha dentro do limite da forma.
	Wrap bool
	// SkinAsset identifica uma imagem decorativa do acervo de skins incorporado.
	SkinAsset string
	// FontFamily contém a face declarada pelo skin, quando houver.
	FontFamily string
	// FontAsset identifica os bytes da face incorporada usada pelo desktop.
	FontAsset string
	// VerticalAlign centraliza legendas dentro dos seus limites lógicos.
	VerticalAlign string
	// Geometry aceita ellipse para marcadores independentes de glifos.
	Geometry string
	// FillAlpha preserva alfa parcial; zero mantém o padrão opaco se Fill existir.
	FillAlpha int
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

// Component expande a aparência confirmada; estados de interação não alteram o documento.
func Component(component screen.Component) Group {
	return ComponentState(component, "default")
}

// ComponentState resolve uma cena compartilhada por edição, prévia e PowerPoint.
// stateName é transitório (default, focused ou touched), nunca serializado.
func ComponentState(component screen.Component, stateName string) Group {
	// group mantém identidade e ordem das partes derivadas.
	group := Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds}
	// box preserva a geometria confirmada, inclusive a de arquivos antigos.
	box := component.Bounds
	if err := ValidateStyle(component); err != nil {
		// O aviso visível preserva os dados sem desenhar uma substituição silenciosa.
		group.Shapes = []Primitive{{Bounds: box, Text: err.Error(), Color: "#B91C1C", FontSize: 12, FontFamily: "Arial", FontAsset: "/fonts/arial.ttf"}}
		return group
	}
	if !component.Enabled {
		stateName = "disabled"
	}
	// name resolve o default legado somente na apresentação, sem regravar o projeto.
	name := component.Style.Skin
	if name == "" {
		name = defaultSkin(component)
	}
	// state contém os tokens originais da definição exata.
	state, _ := skin.MustBuiltin().State(component.Kind, name, stateName)
	switch component.Kind {
	case catalog.Button:
		group.Shapes = append(group.Shapes, decoration(component, state, box)...)
		// caption conserva a legenda como texto separado da imagem original.
		caption := styledText(component, state, component.Properties.Button.Text, inset(box, 4, 0))
		if caption.TextAlign == "" {
			caption.TextAlign = "center"
		}
		caption.VerticalAlign = "center"
		caption.Wrap = true
		group.Shapes = append(group.Shapes, caption)
	case catalog.Label, catalog.Title:
		group.Shapes = append(group.Shapes, decoration(component, state, box)...)
		group.Shapes = append(group.Shapes, styledText(component, state, Caption(component), box))
	case catalog.TextField:
		group.Shapes = append(group.Shapes, decoration(component, state, box)...)
		// value conserva o alvo estável da digitação animada no PPT.
		value := styledText(component, state, component.Properties.Field.InitialValue, inset(box, 6, 0))
		value.Role, value.VerticalAlign = "value", "center"
		group.Shapes = append(group.Shapes, value)
	case catalog.Radio:
		// index e option distribuem os itens dentro dos bounds, sem altura fixa excedente.
		for index, option := range component.Properties.Radio.Options {
			// bounds é a região compartilhada pela pintura e pelo clique.
			bounds := RadioOptionBounds(component, index)
			if component.Style.Skin != "" {
				// selected usa touched somente na opção escolhida e preserva disabled.
				selected := stateName
				if component.Enabled && option.ID == component.Properties.Radio.SelectedOptionID {
					selected = "touched"
				}
				// optionState resolve o estado da definição escolhida pelo usuário.
				optionState, _ := skin.MustBuiltin().State(component.Kind, name, selected)
				group.Shapes = append(group.Shapes, decoration(component, optionState, bounds)...)
				// caption segue a mesma tipografia das demais partes da cena.
				caption := styledText(component, optionState, option.Text, inset(bounds, 4, 0))
				if caption.TextAlign == "" {
					caption.TextAlign = "center"
				}
				caption.VerticalAlign = "center"
				group.Shapes = append(group.Shapes, caption)
			} else {
				// marker desenha o círculo como forma editável, sem depender de glifo na fonte.
				marker := screen.Bounds{X: bounds.X + 2, Y: bounds.Y + (bounds.Height-12)/2, Width: 12, Height: 12}
				group.Shapes = append(group.Shapes, Primitive{Bounds: marker, Stroke: "#4864A0", Geometry: "ellipse"})
				if option.ID == component.Properties.Radio.SelectedOptionID {
					group.Shapes = append(group.Shapes, Primitive{Bounds: inset(marker, 3, 3), Fill: "#4864A0", Geometry: "ellipse"})
				}
				// caption reserva a área do marcador e mantém a opção em texto independente.
				caption := styledText(component, state, option.Text, screen.Bounds{X: bounds.X + 20, Y: bounds.Y, Width: max(1, bounds.Width-20), Height: bounds.Height})
				if component.Style.FontSize == 0 {
					caption.FontSize = 16
				}
				if component.Style.TextColor == "" {
					caption.Color = "#4864A0"
				}
				caption.VerticalAlign = "center"
				group.Shapes = append(group.Shapes, caption)
			}
		}
	case catalog.Report:
		group.Shapes = reportShapes(component, stateName)
	case catalog.NumericKeyboard:
		// key define a geometria única usada também pelas regiões interativas.
		for _, key := range KeyboardKeys(component.Bounds) {
			// part representa somente uma tecla derivada, nunca uma nova instância persistida.
			part := component
			part.Kind, part.Bounds = catalog.Button, key.Bounds
			part.Style.Skin = key.Skin
			if component.Style.Skin != "" && !key.Auxiliary {
				part.Style.Skin = component.Style.Skin
			}
			part.Properties = catalog.Properties{Button: &catalog.ButtonProperties{Text: key.Label, Variant: "confirm"}}
			group.Shapes = append(group.Shapes, ComponentState(part, stateName).Shapes...)
		}
	case catalog.BarcodeInput:
		group.Shapes = BarcodeInputState(component, BarcodeState{}).Shapes
	case catalog.BarcodeProcessed:
		group.Shapes = BarcodeProcessed(component).Shapes
	case catalog.ApplicationMenu:
		group.Shapes = ApplicationMenu(component).Shapes
	}
	if component.Style.BackgroundAssetID != "" {
		group.Shapes = append([]Primitive{{Bounds: box, AssetID: component.Style.BackgroundAssetID}}, group.Shapes...)
	}
	return group
}

// defaultSkin centraliza os defaults por família sem usar o tema do editor.
func defaultSkin(component screen.Component) string {
	switch component.Kind {
	case catalog.Button:
		switch component.Properties.Button.Variant {
		case "cancel":
			return "botao_cancelar"
		case "yes":
			return "botao_sim"
		case "no":
			return "botao_nao"
		}
		return "botao_confirmar"
	case catalog.Label:
		return "label_texto_azul"
	case catalog.Title:
		return "label_titulo_1"
	case catalog.TextField:
		return "type1"
	case catalog.Radio:
		return "botao_uf"
	}
	return ""
}

// styledText combina tokens do estado com overrides explícitos da instância.
func styledText(component screen.Component, state skin.State, value string, bounds screen.Bounds) Primitive {
	// result conserva texto editável e os limites lógicos do parágrafo.
	result := Primitive{Bounds: bounds, Text: value, FontSize: 12, Color: "#4864A0", TextAlign: component.Style.TextAlign}
	if state.Font != nil {
		result.FontSize, result.Color = state.Font.Size, state.Font.Color.Hex()
		result.FontAsset = state.Font.Type
		result.FontFamily = skin.MustBuiltin().FontFamily(state.Font.Type)
		// O builder Java usa a face física: flags XML não comprovam síntese de peso/itálico.
	}
	if component.Style.FontSize != 0 {
		result.FontSize = component.Style.FontSize
	}
	if component.Style.TextColor != "" {
		result.Color = component.Style.TextColor
	}
	if component.Style.FontFamily != "" {
		result.FontFamily = component.Style.FontFamily
		result.FontAsset = fontAsset(component.Style.FontFamily)
	}
	if result.FontSize == 0 {
		result.FontSize = 12
	}
	result.Bold, result.Italic = skin.MustBuiltin().FontStyle(result.FontAsset)
	result.LineSpacing = state.Spacement
	return result
}

// fontAsset encontra apenas famílias incorporadas, sem inventar equivalência tipográfica.
func fontAsset(family string) string {
	// resource percorre as faces disponíveis do acervo operacional.
	for _, resource := range []string{"/fonts/arial.ttf", "/fonts/futura_mh.ttf", "/fonts/futurah.ttf", "/fonts/futuram.ttf", "/fonts/futurbci.ttf", "/fonts/veramobd.ttf"} {
		if skin.MustBuiltin().FontFamily(resource) == family {
			return resource
		}
	}
	return ""
}

// decoration conserva a imagem original sem sobrepor um retângulo genérico opaco.
func decoration(component screen.Component, state skin.State, bounds screen.Bounds) []Primitive {
	// shapes reserva a ordem fundo, imagem e contorno opcional.
	shapes := []Primitive{}
	// fill aplica o override explícito ou o fundo declarado (inclusive transparente).
	fill := component.Style.FillColor
	if fill == "" {
		fill = state.Background.Hex()
	}
	if fill != "" {
		shapes = append(shapes, Primitive{Bounds: bounds, Fill: fill, FillAlpha: state.Background.A})
	}
	if state.ImagePath != "" && component.Style.FillColor == "" {
		shapes = append(shapes, Primitive{Bounds: bounds, SkinAsset: state.ImagePath, Name: "skin"})
	}
	// border só sintetiza a borda padrão de entrada; a decoração dos botões já a contém.
	border := component.Style.BorderColor
	if border == "" && component.Kind == catalog.TextField {
		border = "#000000"
	}
	if border != "" {
		shapes = append(shapes, Primitive{Bounds: bounds, Stroke: border})
	}
	return shapes
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
	case catalog.BarcodeInput:
		return component.Properties.BarcodeInput.Label
	case catalog.BarcodeProcessed:
		if len(component.Properties.BarcodeProcessed.Blocks) > 0 && component.Properties.BarcodeProcessed.Blocks[0].Title != "" {
			return component.Properties.BarcodeProcessed.Blocks[0].Title
		}
		return "Código de barras processado"
	case catalog.ApplicationMenu:
		return component.Properties.Menu.Root.Name
	}
	return fmt.Sprintf("%s", component.Kind)
}
