package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/screen"
)

// configureBackground torna a configuração acessível mesmo com um componente selecionado.
func (e *Editor) configureBackground() {
	if e.page() == nil {
		return
	}
	// chooser fecha antes de executar outro diálogo ou comando de edição.
	var chooser dialog.Dialog
	// original aplica exatamente a imagem fornecida no material legado.
	original := widget.NewButton("Usar fundo Loterias CAIXA", func() {
		chooser.Hide()
		e.useDefaultBackground()
	})
	// imported permite substituir o fundo por PNG ou JPEG incorporado no projeto.
	imported := widget.NewButton("Importar outra imagem…", func() {
		chooser.Hide()
		e.importImage("", true)
	})
	// remove conserva a cor da tela e permite desfazer a remoção da imagem.
	remove := widget.NewButton("Remover imagem de fundo", func() {
		chooser.Hide()
		if e.check(e.Session.Change(func(document *application.Document) error {
			document.Screen(e.PageID).Background.AssetID = ""
			return nil
		})) {
			e.changed()
		}
	})
	chooser = dialog.NewCustom("Fundo da tela", "Fechar", container.NewVBox(
		widget.NewLabel("O fundo fica atrás dos componentes e acompanha o projeto."), original, imported, remove), e.Window)
	chooser.Show()
}

// useDefaultBackground publica a troca de fundo como uma única ação reversível.
func (e *Editor) useDefaultBackground() {
	if e.check(e.Session.Change(func(document *application.Document) error {
		return builtinassets.UseBackground(document, e.PageID)
	})) {
		e.changed()
	}
}

// resourceImage resolve uma imagem incorporada e aplica geometria lógica sem interceptar eventos.
func resourceImage(document application.Document, assetID string, bounds screen.Bounds) *canvas.Image {
	// asset vincula o identificador estável ao conteúdo salvo no pacote.
	for _, asset := range document.Project.Assets {
		if asset.ID == assetID {
			// data permanece local ao documento, sem filesystem nem conexão de rede.
			data := document.Resources[asset.Path]
			if len(data) == 0 {
				return nil
			}
			// image conserva o recurso proporcional à superfície lógica configurada.
			image := canvas.NewImageFromResource(fyne.NewStaticResource(asset.Path, data))
			image.FillMode = canvas.ImageFillStretch
			place(image, bounds)
			return image
		}
	}
	return nil
}

// pageBackground desenha cor e imagem antes dos controles na prévia e na demonstração.
func pageBackground(document application.Document, page screen.Screen) []fyne.CanvasObject {
	// bounds define a área lógica completa, preservada em todas as visitas.
	bounds := screen.Bounds{Width: page.Size.Width, Height: page.Size.Height}
	// surface fornece a cor alternativa atrás da imagem opcional.
	surface := canvas.NewRectangle(rgb(page.Background.Color))
	place(surface, bounds)
	// objects respeita a ordem do fundo para a frente.
	objects := []fyne.CanvasObject{surface}
	if image := resourceImage(document, page.Background.AssetID, bounds); image != nil {
		objects = append(objects, image)
	}
	return objects
}
