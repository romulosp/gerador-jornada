package desktop

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
	"image"
	"sync"
)

// skinImages compartilha imagens imutáveis decodificadas durante a inicialização.
var skinImages sync.Map

// cachedSkinImage preserva pixels decodificados e proveniência dos bytes consumidos pela cena.
type cachedSkinImage struct {
	// Pixels é imutável e Resource conserva o mesmo PNG/JPEG original, incluindo alpha.
	Pixels   image.Image
	Resource fyne.Resource
}

// decodedSkinImage aquece o mesmo cache utilizado pelo canvas e pela prévia.
func decodedSkinImage(path string) (*cachedSkinImage, error) {
	if cached, ok := skinImages.Load(path); ok {
		return cached.(*cachedSkinImage), nil
	}
	// data e decoded preservam os bytes originais e seus pixels RGBA.
	data, err := skin.MustBuiltin().Asset(path)
	if err != nil {
		return nil, err
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resource := &cachedSkinImage{Pixels: decoded, Resource: fyne.NewStaticResource(path, data)}
	skinImages.Store(path, resource)
	return resource, nil
}

// skinImage resolve uma imagem decorativa do acervo de skins incorporado.
// O caminho lógico nunca é lido do disco em execução.
func skinImage(path string, bounds screen.Bounds) fyne.CanvasObject {
	// data contém os bytes embutidos do recurso.
	data, err := decodedSkinImage(path)
	if err != nil {
		return nil
	}
	// image usa a mesma escala lógica das demais primitivas.
	image := canvas.NewImageFromImage(data.Pixels)
	image.Resource = data.Resource
	image.FillMode = canvas.ImageFillStretch
	place(image, bounds)
	return image
}
