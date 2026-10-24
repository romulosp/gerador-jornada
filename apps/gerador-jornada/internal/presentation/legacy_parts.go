package presentation

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// KeyboardKey descreve uma tecla derivada do construtor Java, em coordenadas da tela.
type KeyboardKey struct {
	// Label contém a legenda original, independente do comando de entrada.
	Label string
	// Bounds conserva a mesma região para desenho, foco e clique.
	Bounds screen.Bounds
	// Skin identifica a definição de botão usada pela tecla.
	Skin string
	// Auxiliary distingue TAB e LIMPAR dos botões quadrados.
	Auxiliary bool
}

// KeyboardKeys escala o arranjo Java 200 × 350, inclusive suas margens e área livre.
// Bounds de projetos anteriores são respeitados; esta função não migra o documento.
func KeyboardKeys(bounds screen.Bounds) []KeyboardKey {
	// keys conserva a ordem visual: 7–9, 4–6, 1–3, setas/zero e auxiliares.
	keys := []KeyboardKey{}
	// index e label identificam a posição de cada uma das doze teclas quadradas.
	for index, label := range []string{"7", "8", "9", "4", "5", "6", "1", "2", "3", "<", "0", ">"} {
		keys = append(keys, KeyboardKey{Label: label, Bounds: scalePart(bounds, screen.Bounds{X: 5 + (index%3)*59, Y: 5 + (index/3)*55, Width: 52, Height: 52}, 200, 350), Skin: "botao_teclado"})
	}
	keys = append(keys,
		KeyboardKey{Label: "TAB", Bounds: scalePart(bounds, screen.Bounds{X: 5, Y: 225, Width: 80, Height: 60}, 200, 350), Skin: "botao_teclado_auxiliar_1", Auxiliary: true},
		KeyboardKey{Label: "LIMPAR", Bounds: scalePart(bounds, screen.Bounds{X: 98, Y: 225, Width: 80, Height: 60}, 200, 350), Skin: "botao_teclado_auxiliar_1", Auxiliary: true})
	return keys
}

// scalePart converte medidas de referência em limites absolutos sem reflow.
func scalePart(bounds, part screen.Bounds, width, height int) screen.Bounds {
	return screen.Bounds{X: bounds.X + part.X*bounds.Width/width, Y: bounds.Y + part.Y*bounds.Height/height,
		Width: max(1, part.Width*bounds.Width/width), Height: max(1, part.Height*bounds.Height/height)}
}

// RadioOptionBounds divide a área configurada, mantendo pintura e hit-test consistentes.
func RadioOptionBounds(component screen.Component, index int) screen.Bounds {
	// bounds conserva a origem absoluta e o tamanho explícito do grupo.
	bounds := component.Bounds
	// count impede divisão por zero em rascunhos sem opções.
	count := len(component.Properties.Radio.Options)
	if count == 0 {
		return bounds
	}
	if component.Properties.Radio.Orientation == "horizontal" {
		bounds.X += index * component.Bounds.Width / count
		bounds.Width = max(1, (index+1)*component.Bounds.Width/count-index*component.Bounds.Width/count)
	} else {
		bounds.Y += index * component.Bounds.Height / count
		bounds.Height = max(1, (index+1)*component.Bounds.Height/count-index*component.Bounds.Height/count)
	}
	return bounds
}

// reportShapes recompõe as decorações do relatório com linhas e título editáveis.
// O modelo atual oferece título e linhas simples; variantes Java customizadas permanecem pendentes.
func reportShapes(component screen.Component, stateName string) []Primitive {
	// box preserva as dimensões externas, incluindo cabeçalho 26 e rodapé 32 da base 300 × 345.
	box := component.Bounds
	// shapes contém partes estáticas, sem rolagem ou seleção persistidas.
	shapes := []Primitive{}
	// appendPart aplica a definição exata do painel ou rótulo original.
	appendPart := func(componentType, name, value string, bounds screen.Bounds) {
		// state conserva os tokens da parte sem herdar a skin de outra família.
		state, _ := skin.MustBuiltin().StateByType(componentType, name, stateName)
		shapes = append(shapes, decoration(component, state, bounds)...)
		if value != "" {
			shapes = append(shapes, styledText(component, state, value, inset(bounds, 6, 0)))
		}
	}
	appendPart("panel", "painel_tabela", "", box)
	appendPart("label", "label_tabela_cabecalho", "", scalePart(box, screen.Bounds{Width: 300, Height: 26}, 300, 345))
	appendPart("label", "label_tabela_rodape_2", "", scalePart(box, screen.Bounds{Y: 313, Width: 300, Height: 32}, 300, 345))
	// titleBounds reserva a linha de título do relatório dentro do corpo.
	titleBounds := scalePart(box, screen.Bounds{Y: 26, Width: 292, Height: 26}, 300, 345)
	appendPart("panel", "painel_report_titulo", "", titleBounds)
	appendPart("label", "label_report_titulo", component.Properties.Report.Title, titleBounds)
	// index e row preservam a ordem das linhas e limitam o recorte ao corpo visível.
	for index, row := range component.Properties.Report.Rows {
		// bounds não invade o rodapé; o conteúdo completo permanece salvo no modelo.
		bounds := scalePart(box, screen.Bounds{Y: 52 + index*20, Width: 292, Height: 20}, 300, 345)
		if bounds.Y+bounds.Height > box.Y+313*box.Height/345 {
			break
		}
		appendPart("panel", "painel_report_simples", "", bounds)
		appendPart("label", "label_report", row.Text, bounds)
	}
	// explicit permite personalizar a família table sem apagar a composição por partes.
	if component.Style.Skin != "" {
		// state conserva somente uma skin table escolhida explicitamente.
		state, ok := skin.MustBuiltin().State(catalog.Report, component.Style.Skin, stateName)
		if ok {
			// primitive atualiza as linhas de texto, preservando a decoração original.
			for index, primitive := range shapes {
				if primitive.Text != "" {
					shapes[index] = styledText(component, state, primitive.Text, primitive.Bounds)
				}
			}
		}
	}
	return shapes
}
