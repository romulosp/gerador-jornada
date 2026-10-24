package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// skinImage resolve uma imagem decorativa do acervo de skins incorporado.
// O caminho lógico nunca é lido do disco em execução.
func skinImage(path string, bounds screen.Bounds) fyne.CanvasObject {
	// data contém os bytes embutidos do recurso.
	data, err := skin.MustBuiltin().Asset(path)
	if err != nil || len(data) == 0 {
		return nil
	}
	// image usa a mesma escala lógica das demais primitivas.
	image := canvas.NewImageFromResource(fyne.NewStaticResource(path, data))
	image.FillMode = canvas.ImageFillStretch
	place(image, bounds)
	return image
}
