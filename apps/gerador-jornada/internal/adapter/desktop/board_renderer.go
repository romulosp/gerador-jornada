package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// boardGroup conserva formas de um controle e suas posições antes de uma translação provisória.
type boardGroup struct {
	// Bounds identifica a geometria utilizada para gerar as formas.
	Bounds screen.Bounds
	// Objects contém primitivas reutilizáveis do componente.
	Objects []fyne.CanvasObject
	// Positions conserva posições base, sem acumular erros entre quadros de arraste.
	Positions []fyne.Position
}

// boardRenderer mantém cache visual por sessão, tela, revisão e zoom.
type boardRenderer struct {
	// board fornece estado de interação sem misturá-lo ao modelo persistido.
	board *Board
	// objects contém a lista final de desenho do quadro corrente.
	objects []fyne.CanvasObject
	// session detecta abertura de outro snapshot mesmo com IDs/revisões iguais.
	session *application.Session
	// pageID, revision e zoom formam a chave de invalidação das primitivas.
	pageID   string
	revision uint64
	zoom     float32
	// background mantém cor e imagem já resolvidas no documento.
	background []fyne.CanvasObject
	// groups permite reutilizar formas sem reconstruir todos os controles ao selecionar.
	groups map[string]*boardGroup
	// outlines reutiliza contornos de seleção conforme a quantidade ativa.
	outlines []*canvas.Rectangle
	// overlays mantém contornos acima de todos os componentes, inclusive dos que cobrem a seleção.
	overlays []fyne.CanvasObject
	// marquee representa a caixa provisória de seleção em área vazia.
	marquee *canvas.Rectangle
}

// Layout preserva geometria absoluta definida pelo domínio Tela.
func (r *boardRenderer) Layout(_ fyne.Size) {}

// MinSize delega o tamanho lógico e a escala ao artboard.
func (r *boardRenderer) MinSize() fyne.Size { return r.board.MinSize() }

// Objects expõe apenas as formas do quadro atual.
func (r *boardRenderer) Objects() []fyne.CanvasObject { return r.objects }

// Destroy não possui recursos externos para liberar.
func (r *boardRenderer) Destroy() {}

// Refresh reutiliza a página estável e altera somente geometria provisória e realces de seleção.
func (r *boardRenderer) Refresh() {
	// editor identifica o documento e a escala do quadro.
	editor := r.board.Editor
	// page fornece os dados confirmados da superfície ativa.
	page := editor.page()
	if page == nil {
		r.objects = []fyne.CanvasObject{widget.NewLabel("Adicione uma tela para começar.")}
		return
	}
	editor.selectionIDs()
	if r.session != editor.Session || r.pageID != page.ID || r.revision != editor.Session.Document.Project.Revision || r.zoom != editor.Zoom {
		r.rebuild(*page)
	}
	r.objects = append(r.objects[:0], r.background...)
	r.overlays = r.overlays[:0]
	// selectedIndex reutiliza contornos em vez de alocar novos a cada clique.
	selectedIndex := 0
	// component determina a ordem de camadas do documento, sem reconstruir suas formas.
	for _, component := range page.Components {
		// group é uma composição em cache da revisão atual.
		group := r.groups[component.ID]
		// bounds será substituído apenas durante um resize provisório.
		bounds := component.Bounds
		if r.board.Dragging && r.board.Resizing && component.ID == editor.Selected {
			bounds = r.board.Pending
		}
		if group == nil || group.Bounds != bounds {
			// copy permite expandir um resize sem publicar dados no documento.
			copy := component
			copy.Bounds = bounds
			group = r.component(copy)
			r.groups[component.ID] = group
		}
		// dx e dy descrevem uma translação visual do grupo ativo.
		dx, dy := float32(0), float32(0)
		if r.board.Dragging && !r.board.Marquee && !r.board.Resizing && editor.isSelected(component.ID) {
			dx, dy = float32(r.board.OffsetX)*editor.Zoom, float32(r.board.OffsetY)*editor.Zoom
			bounds.X += r.board.OffsetX
			bounds.Y += r.board.OffsetY
		}
		// index e object reposicionam as mesmas formas sem alterar o tamanho nem recriar imagens.
		for index, object := range group.Objects {
			// position é derivada da origem estável para evitar acumulação de arredondamento.
			position := group.Positions[index].Add(fyne.NewPos(dx, dy))
			if object.Position() != position {
				object.Move(position)
			}
		}
		r.objects = append(r.objects, group.Objects...)
		if editor.isSelected(component.ID) {
			// outline contém apenas a decoração da interface, nunca as formas exportadas.
			outline := r.outline(selectedIndex)
			r.place(outline, bounds)
			r.overlays = append(r.overlays, outline)
			selectedIndex++
			if len(editor.Selection) == 1 {
				// handle só aparece quando redimensionar um único controle é permitido.
				handle := r.outline(selectedIndex)
				handle.FillColor = rgb("#2563EB")
				r.place(handle, screen.Bounds{X: bounds.X + bounds.Width - 8, Y: bounds.Y + bounds.Height - 8, Width: 8, Height: 8})
				r.overlays = append(r.overlays, handle)
				selectedIndex++
			}
		}
	}
	r.objects = append(r.objects, r.overlays...)
	if r.board.Marquee {
		if r.marquee == nil {
			r.marquee = canvas.NewRectangle(transparent())
			r.marquee.StrokeColor, r.marquee.StrokeWidth = rgb("#2563EB"), 1
		}
		r.place(r.marquee, r.board.Area)
		r.objects = append(r.objects, r.marquee)
	}
	canvas.Refresh(r.board)
}

// rebuild invalida recursos somente quando dados confirmados ou escala mudam.
func (r *boardRenderer) rebuild(page screen.Screen) {
	// editor fornece o snapshot confirmado e a escala corrente.
	editor := r.board.Editor
	r.session, r.pageID, r.revision, r.zoom = editor.Session, page.ID, editor.Session.Document.Project.Revision, editor.Zoom
	r.groups = map[string]*boardGroup{}
	// surface mantém o fundo de cor sob a imagem opcional.
	surface := canvas.NewRectangle(rgb(page.Background.Color))
	surface.StrokeColor, surface.StrokeWidth = rgb("#CBD5E1"), 1
	r.place(surface, screen.Bounds{Width: page.Size.Width, Height: page.Size.Height})
	r.background = []fyne.CanvasObject{surface}
	if image := resourceImage(editor.Session.Document, page.Background.AssetID, screen.Bounds{Width: page.Size.Width, Height: page.Size.Height}); image != nil {
		r.place(image, screen.Bounds{Width: page.Size.Width, Height: page.Size.Height})
		r.background = append(r.background, image)
	}
	for _, component := range page.Components {
		r.groups[component.ID] = r.component(component)
	}
}

// component expande somente o controle que ainda não possui formas compatíveis em cache.
func (r *boardRenderer) component(component screen.Component) *boardGroup {
	// group mantém posições fixas e formas independentes das alças do editor.
	group := &boardGroup{Bounds: component.Bounds}
	// primitive é uma parte compartilhada com a prévia e o PPT.
	for _, primitive := range presentation.Component(component).Shapes {
		// object já contém a escala lógica; somente a margem do canvas é acrescentada aqui.
		for _, object := range scenePrimitiveObjects(r.board.Editor.Session.Document, primitive, r.board.Editor.Zoom) {
			object.Move(object.Position().Add(fyne.NewPos(24, 24)))
			group.Objects = append(group.Objects, object)
		}
	}

	for _, object := range group.Objects {
		group.Positions = append(group.Positions, object.Position())
	}
	return group
}

// outline reutiliza realces e restaura seus estilos ao alternar entre alça e contorno.
func (r *boardRenderer) outline(index int) *canvas.Rectangle {
	for len(r.outlines) <= index {
		r.outlines = append(r.outlines, canvas.NewRectangle(transparent()))
	}
	// outline pode ter sido uma alça no quadro anterior; o preenchimento deve voltar a transparente.
	outline := r.outlines[index]
	outline.FillColor, outline.StrokeColor, outline.StrokeWidth = transparent(), rgb("#2563EB"), 2
	return outline
}

// place converte coordenadas uma única vez, mantendo margem externa e escala consistentes.
func (r *boardRenderer) place(object fyne.CanvasObject, bounds screen.Bounds) {
	// zoom altera a visualização, nunca a geometria persistida.
	zoom := r.board.Editor.Zoom
	object.Move(fyne.NewPos(float32(bounds.X)*zoom+24, float32(bounds.Y)*zoom+24))
	object.Resize(fyne.NewSize(float32(bounds.Width)*zoom, float32(bounds.Height)*zoom))
}

// transparent fornece preenchimento sem opacidade para alças e formas sem fundo.
func transparent() colorTransparent { return colorTransparent{} }

// colorTransparent implementa canais totalmente transparentes.
type colorTransparent struct{}

// RGBA retorna ausência de opacidade ao renderer.
func (colorTransparent) RGBA() (uint32, uint32, uint32, uint32) { return 0, 0, 0, 0 }
