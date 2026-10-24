package desktop

import (
	"image/color"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// sceneFonts reutiliza os mesmos bytes SFNT por face sem carregar fontes do sistema.
var sceneFonts sync.Map

// sceneFont resolve somente recursos incorporados, sem consultar caminhos externos.
func sceneFont(path string) fyne.Resource {
	if path == "" {
		return nil
	}
	// cached evita reconstruir o recurso a cada quadro ou tecla.
	if cached, ok := sceneFonts.Load(path); ok {
		// resource é o único tipo armazenado neste cache privado.
		resource, _ := cached.(fyne.Resource)
		return resource
	}
	// data conserva a face física original do legado.
	data, err := skin.MustBuiltin().Asset(path)
	if err != nil {
		return nil
	}
	// resource compartilha os bytes imutáveis com os renderers.
	resource := fyne.NewStaticResource(path, data)
	sceneFonts.Store(path, resource)
	return resource
}

// scenePrimitiveObjects traduz uma parte da cena uma única vez para os consumidores Fyne.
// zoom altera somente a apresentação, nunca a geometria persistida.
func scenePrimitiveObjects(document application.Document, primitive presentation.Primitive, zoom float32) []fyne.CanvasObject {
	// objects conserva imagem, forma e texto na ordem original.
	objects := []fyne.CanvasObject{}
	// position aplica escala uniforme, sem acrescentar margens do editor à prévia.
	position := func(object fyne.CanvasObject, bounds screen.Bounds) {
		object.Move(fyne.NewPos(float32(bounds.X)*zoom, float32(bounds.Y)*zoom))
		object.Resize(fyne.NewSize(float32(bounds.Width)*zoom, float32(bounds.Height)*zoom))
		objects = append(objects, object)
	}
	if primitive.SkinAsset != "" {
		if image := skinImage(primitive.SkinAsset, primitive.Bounds); image != nil {
			position(image, primitive.Bounds)
		}
	}
	if primitive.AssetID != "" {
		if image := resourceImage(document, primitive.AssetID, primitive.Bounds); image != nil {
			position(image, primitive.Bounds)
		}
	}
	if primitive.Fill != "" || primitive.Stroke != "" {
		// fill respeita alfa parcial e ausência de preenchimento.
		fill := color.Color(transparent())
		if primitive.Fill != "" {
			// parsed transforma RGB sem perder o alfa importado.
			parsed := color.NRGBAModel.Convert(rgb(primitive.Fill)).(color.NRGBA)
			if primitive.FillAlpha > 0 {
				parsed.A = uint8(primitive.FillAlpha)
			}
			fill = parsed
		}
		if primitive.Geometry == "ellipse" {
			// circle dispensa a presença dos glifos de marcador na fonte escolhida.
			circle := canvas.NewCircle(fill)
			if primitive.Stroke != "" {
				circle.StrokeColor, circle.StrokeWidth = rgb(primitive.Stroke), zoom
			}
			position(circle, primitive.Bounds)
		} else {
			// rectangle conserva fundo e contorno independentes da tipografia.
			rectangle := canvas.NewRectangle(fill)
			if primitive.Stroke != "" {
				rectangle.StrokeColor, rectangle.StrokeWidth = rgb(primitive.Stroke), zoom
			}
			position(rectangle, primitive.Bounds)
		}
	}
	if primitive.Text != "" {
		objects = append(objects, sceneTextObjects(primitive, zoom)...)
	}
	return objects
}

// sceneTextObjects mede a face real e quebra legendas sem substituí-las por bitmaps.
// O fator entre linhas vem do spacement original, inclusive 0,8 nos menus.
func sceneTextObjects(primitive presentation.Primitive, zoom float32) []fyne.CanvasObject {
	// newText mantém a mesma face, estilo físico e cor para medição e pintura.
	newText := func(value string) *canvas.Text {
		// text conserva os bytes SFNT em vez de procurar uma fonte parecida no sistema.
		text := canvas.NewText(value, rgb(primitive.Color))
		text.FontSource = sceneFont(primitive.FontAsset)
		text.TextSize = float32(primitive.FontSize) * zoom
		text.TextStyle.Bold, text.TextStyle.Italic = primitive.Bold, primitive.Italic
		setTextAlignment(text, primitive.TextAlign)
		return text
	}
	// width corresponde à caixa tipográfica lógica na escala exibida.
	width := float32(primitive.Bounds.Width) * zoom
	// lines mantém quebras explícitas e aplica a reserva total TEXTBOX_BORDER=10 do Java.
	lines := []string{}
	// paragraph conserva as quebras fornecidas pelo autor.
	for _, paragraph := range strings.Split(primitive.Text, "\n") {
		if !primitive.Wrap {
			lines = append(lines, paragraph)
			continue
		}
		// line acumula palavras completas, como o SkinUtil original.
		line := ""
		// word mantém uma palavra inteira mesmo quando ultrapassa a largura disponível.
		for _, word := range strings.Fields(paragraph) {
			// candidate testa a próxima palavra com a mesma fonte usada no desenho.
			candidate := word
			if line != "" {
				candidate = line + " " + word
			}
			if line != "" && newText(candidate).MinSize().Width > width-10*zoom {
				lines = append(lines, line)
				line = word
			} else {
				line = candidate
			}
		}
		lines = append(lines, line)
	}
	// height usa as métricas do renderer da face física, sem estimativa por quantidade de letras.
	height := newText("Ag").MinSize().Height
	// spacing define a distância entre baselines; ausência usa uma linha completa.
	spacing := float32(primitive.LineSpacing)
	if spacing <= 0 {
		spacing = 1
	}
	// top centraliza o bloco tipográfico somente quando a cena solicita.
	top := float32(primitive.Bounds.Y) * zoom
	if primitive.VerticalAlign == "center" {
		top += (float32(primitive.Bounds.Height)*zoom - height - float32(len(lines)-1)*height*spacing) / 2
	}
	// objects preserva as caixas individuais sem acrescentar seletores de tema ou botões Fyne.
	objects := []fyne.CanvasObject{}
	// index e line mantêm a ordem das linhas editáveis da legenda.
	for index, line := range lines {
		// text recebe o mesmo alinhamento horizontal dentro da largura completa.
		text := newText(line)
		text.Move(fyne.NewPos(float32(primitive.Bounds.X)*zoom, top+float32(index)*height*spacing))
		text.Resize(fyne.NewSize(width, height))
		objects = append(objects, text)
	}
	return objects
}
