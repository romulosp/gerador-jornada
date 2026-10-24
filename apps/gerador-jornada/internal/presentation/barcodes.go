package presentation

import (
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// BarcodeInputWidth e BarcodeInputHeight definem a geometria-base da entrada de código de barras.
const (
	// BarcodeInputWidth expressa a largura de referência em unidades lógicas.
	BarcodeInputWidth = 610
	// BarcodeInputHeight expressa a altura de referência em unidades lógicas.
	BarcodeInputHeight = 60
)

// barcodeBaseFields contém as posições e larguras dos oito campos na escala-base.
var barcodeBaseFields = [8]struct{ x, width int }{
	{30, 50}, {85, 50}, {140, 50}, {195, 60}, {260, 50}, {315, 60}, {380, 15}, {400, 110},
}

// barcodeSeparators contém os três pontos separadores na escala-base.
var barcodeSeparators = [3]struct{ x, y int }{{81, 27}, {191, 27}, {311, 27}}

// BarcodeState contém somente o estado transitório de uma instância na prévia.
// Nenhum destes valores é persistido no projeto ou exportado para o PowerPoint.
type BarcodeState struct {
	// Values contém os oito campos visuais já digitados.
	Values [8]string
	// Active identifica o campo que recebe digitação, de zero a sete.
	Active int
	// Locked indica entrada aceita e somente leitura.
	Locked bool
	// Reading indica modo explícito de captura integral.
	Reading bool
	// Message contém o erro local da instância, sempre textual.
	Message string
}

// BarcodeFieldBounds devolve as oito áreas clicáveis na mesma escala da pintura.
func BarcodeFieldBounds(bounds screen.Bounds) [8]screen.Bounds {
	// result preserva a ordem visual dos campos.
	result := [8]screen.Bounds{}
	// index percorre a geometria de referência.
	for index, field := range barcodeBaseFields {
		result[index] = screen.Bounds{
			X:      bounds.X + field.x*bounds.Width/BarcodeInputWidth,
			Y:      bounds.Y + 20*bounds.Height/BarcodeInputHeight,
			Width:  field.width * bounds.Width / BarcodeInputWidth,
			Height: 20 * bounds.Height / BarcodeInputHeight,
		}
	}
	return result
}

// BarcodeFieldLimit devolve o limite de dígitos do campo visual indicado.
func BarcodeFieldLimit(index int) int {
	if index < 0 || index >= len(catalog.BarcodeFieldLimits) {
		return 0
	}
	return catalog.BarcodeFieldLimits[index]
}

// BarcodeSeparatorBounds devolve as três marcas separadoras na mesma escala.
func BarcodeSeparatorBounds(bounds screen.Bounds) [3]screen.Bounds {
	// result preserva a ordem dos separadores.
	result := [3]screen.Bounds{}
	// index percorre a geometria de referência.
	for index, separator := range barcodeSeparators {
		result[index] = screen.Bounds{
			X:     bounds.X + separator.x*bounds.Width/BarcodeInputWidth,
			Y:     bounds.Y + separator.y*bounds.Height/BarcodeInputHeight,
			Width: 5 * bounds.Width / BarcodeInputWidth, Height: 10 * bounds.Height / BarcodeInputHeight,
		}
	}
	return result
}

// BarcodeInput expande a entrada sem estado de runtime, para canvas e PPT estáticos.
func BarcodeInput(component screen.Component) Group {
	return BarcodeInputState(component, BarcodeState{})
}

// BarcodeInputState expande a entrada com os valores e estados da prévia.
// A pintura e as áreas clicáveis usam o mesmo cálculo de escala.
func BarcodeInputState(component screen.Component, state BarcodeState) Group {
	// group preserva a identidade do composto.
	group := Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds, Shapes: []Primitive{}}
	// labelBounds posiciona o rótulo conforme o DESIGN.
	labelBounds := screen.Bounds{
		X:      component.Bounds.X + 30*component.Bounds.Width/BarcodeInputWidth,
		Y:      component.Bounds.Y + 2*component.Bounds.Height/BarcodeInputHeight,
		Width:  200 * component.Bounds.Width / BarcodeInputWidth,
		Height: 20 * component.Bounds.Height / BarcodeInputHeight,
	}
	group.Shapes = append(group.Shapes, Primitive{Bounds: labelBounds, Text: component.Properties.BarcodeInput.Label, Color: "#5064AA", FontSize: 12, Name: "rótulo"})
	// fields contém as áreas pintadas e clicáveis.
	fields := BarcodeFieldBounds(component.Bounds)
	// index percorre os oito campos para pintar estado e conteúdo.
	for index, field := range fields {
		// fill, textColor e border representam normal, foco e desabilitado.
		fill, textColor, border := "#FFFFFF", "#000000", "#CBD5E1"
		if state.Locked {
			fill, textColor = "#E6EEFF", "#808080"
		} else if state.Active == index {
			fill, textColor, border = "#10508E", "#FFFFFF", "#10508E"
		}
		group.Shapes = append(group.Shapes, Primitive{Bounds: field, Fill: fill, Stroke: border, Name: fmt.Sprintf("campo-%d", index+1)})
		group.Shapes = append(group.Shapes, Primitive{Bounds: inset(field, 4, 2), Text: state.Values[index], Color: textColor, FontSize: 12, TextAlign: "center", Name: fmt.Sprintf("valor-%d", index+1)})
	}
	// separator pinta os três pontos sem criar identidade de componente.
	for _, separator := range BarcodeSeparatorBounds(component.Bounds) {
		group.Shapes = append(group.Shapes, Primitive{Bounds: separator, Fill: "#5064AA", Name: "separador"})
	}
	if state.Message != "" {
		group.Shapes = append(group.Shapes, Primitive{Bounds: screen.Bounds{X: component.Bounds.X, Y: component.Bounds.Y + component.Bounds.Height - 14*component.Bounds.Height/BarcodeInputHeight, Width: component.Bounds.Width, Height: 12 * component.Bounds.Height / BarcodeInputHeight}, Text: state.Message, Color: "#B91C1C", FontSize: 10, Name: "mensagem"})
	}
	compoundTypography(component, &group)
	return group
}

// BarcodeProcessed expande a ficha de leitura com blocos e linhas na mesma ordem do modelo.
func BarcodeProcessed(component screen.Component) Group {
	// group preserva a identidade do composto processado.
	group := Group{ComponentID: component.ID, Kind: component.Kind, Bounds: component.Bounds, Shapes: []Primitive{}}
	// scale converte a proposta geométrica de 530 × 120 para o tamanho atual.
	scaleX := func(value int) int { return value * component.Bounds.Width / 530 }
	scaleY := func(value int) int { return value * component.Bounds.Height / 120 }
	// y acumula a posição vertical do próximo bloco.
	y := 12
	// blockIndex identifica o bloco na ordem confirmada.
	for blockIndex, block := range component.Properties.BarcodeProcessed.Blocks {
		if blockIndex > 0 {
			y += 12
		}
		if block.Title != "" {
			group.Shapes = append(group.Shapes, Primitive{Bounds: screen.Bounds{X: component.Bounds.X + scaleX(12), Y: component.Bounds.Y + scaleY(y), Width: scaleX(506), Height: scaleY(20)}, Text: block.Title, Color: "#1E3A8A", FontSize: 12, Bold: true, Name: fmt.Sprintf("bloco-%d-título", blockIndex+1)})
			y += 24
		}
		// row percorre as linhas ordenadas do bloco.
		for rowIndex, row := range block.Rows {
			// label contém o texto derivado do enum, nunca uma string concorrente.
			label, ok := catalog.BarcodeLabelText(row.LabelKind)
			if !ok {
				label = row.LabelKind + ":"
			}
			group.Shapes = append(group.Shapes, Primitive{Bounds: screen.Bounds{X: component.Bounds.X + scaleX(12), Y: component.Bounds.Y + scaleY(y), Width: scaleX(202), Height: scaleY(24)}, Text: label, Color: "#334155", FontSize: 11, Name: fmt.Sprintf("bloco-%d-linha-%d-rótulo", blockIndex+1, rowIndex+1)})
			group.Shapes = append(group.Shapes, Primitive{Bounds: screen.Bounds{X: component.Bounds.X + scaleX(222), Y: component.Bounds.Y + scaleY(y), Width: scaleX(296), Height: scaleY(24)}, Text: row.Value, Color: "#0F172A", FontSize: 11, Wrap: true, Name: fmt.Sprintf("bloco-%d-linha-%d-valor", blockIndex+1, rowIndex+1)})
			y += 24
		}
	}
	compoundTypography(component, &group)
	return group
}

// compoundTypography aplica a face real e overrides sem substituir a geometria das fichas.
func compoundTypography(component screen.Component, group *Group) {
	// index e primitive mantêm cada texto editável na sua posição original.
	for index, primitive := range group.Shapes {
		if primitive.Text == "" && !strings.HasPrefix(primitive.Name, "valor-") {
			continue
		}
		// asset distingue as entradas Arial dos rótulos Futura do acervo.
		asset := "/fonts/futura_mh.ttf"
		if strings.HasPrefix(primitive.Name, "valor-") {
			asset = "/fonts/arial.ttf"
		}
		if component.Style.FontFamily != "" {
			asset = fontAsset(component.Style.FontFamily)
		}
		primitive.FontAsset, primitive.FontFamily = asset, skin.MustBuiltin().FontFamily(asset)
		// bold e italic são propriedades físicas dos mesmos bytes incorporados.
		bold, italic := skin.MustBuiltin().FontStyle(asset)
		primitive.Bold, primitive.Italic = primitive.Bold || bold, italic
		if component.Style.FontSize > 0 {
			primitive.FontSize = component.Style.FontSize
		}
		if component.Style.TextColor != "" {
			primitive.Color = component.Style.TextColor
		}
		if component.Style.TextAlign != "" {
			primitive.TextAlign = component.Style.TextAlign
		}
		group.Shapes[index] = primitive
	}
}
