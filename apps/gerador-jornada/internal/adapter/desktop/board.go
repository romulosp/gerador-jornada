package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// Board implementa seleção e manipulação direta sem gravar cada movimento.
type Board struct {
	widget.BaseWidget
	// Editor fornece a sessão e os comandos confirmados.
	Editor *Editor
	// Dragging indica uma manipulação provisória.
	Dragging bool
	// Resizing indica uso da alça inferior direita.
	Resizing bool
	// Original guarda a geometria antes do arraste.
	Original screen.Bounds
	// Pending guarda a geometria provisória.
	Pending screen.Bounds
	// Delta acumula deslocamento em coordenadas de tela.
	Delta fyne.Position
}

// NewBoard cria o canvas customizado com estado efêmero de manipulação.
func NewBoard(editor *Editor) *Board {
	// board controla somente o desenho e eventos de ponteiro.
	board := &Board{Editor: editor}
	board.ExtendBaseWidget(board)
	return board
}

// MinSize acompanha as dimensões lógicas e o zoom corrente.
func (b *Board) MinSize() fyne.Size {
	// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
	if page := b.Editor.page(); page != nil {
		return fyne.NewSize(float32(page.Size.Width)*b.Editor.Zoom+48, float32(page.Size.Height)*b.Editor.Zoom+48)
	}
	return fyne.NewSize(848, 648)
}

// CreateRenderer produz uma superfície livre sem impor layout aos componentes.
func (b *Board) CreateRenderer() fyne.WidgetRenderer {
	// renderer recria primitivas visuais, nunca regras de negócio.
	renderer := &boardRenderer{board: b}
	renderer.Refresh()
	return renderer
}

// MouseDown seleciona o componente superior e identifica a alça de tamanho.
func (b *Board) MouseDown(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonPrimary {
		return
	}
	b.Editor.Window.Canvas().Unfocus()
	// page identifica a composição clicada.
	page := b.Editor.page()
	if page == nil {
		return
	}
	// x e y convertem o ponteiro para unidades lógicas.
	x, y := b.logical(event.Position)
	b.Editor.Selected = ""
	// index percorre as camadas de frente para trás.
	for index := len(page.Components) - 1; index >= 0; index-- {
		// component é um candidato ao clique.
		component := page.Components[index]
		if x >= component.Bounds.X && x <= component.Bounds.X+component.Bounds.Width && y >= component.Bounds.Y && y <= component.Bounds.Y+component.Bounds.Height {
			b.Editor.Selected = component.ID
			b.Original = component.Bounds
			b.Pending = component.Bounds
			b.Resizing = x > component.Bounds.X+component.Bounds.Width-12 && y > component.Bounds.Y+component.Bounds.Height-12
			break
		}
	}
	b.Delta = fyne.Position{}
	b.Dragging = false
	b.Editor.refreshInspector()
	b.Editor.refreshStructure()
	b.Refresh()
}

// MouseUp não confirma duas vezes o mesmo arraste.
func (b *Board) MouseUp(_ *desktop.MouseEvent) {}

// Dragged atualiza apenas a geometria provisória do canvas.
func (b *Board) Dragged(event *fyne.DragEvent) {
	if b.Editor.Selected == "" || b.Editor.page() == nil {
		return
	}
	b.Dragging = true
	b.Delta = b.Delta.Add(fyne.NewPos(event.Dragged.DX, event.Dragged.DY))
	b.Pending = b.Original
	if b.Resizing {
		b.Pending.Width += int(b.Delta.X / b.Editor.Zoom)
		b.Pending.Height += int(b.Delta.Y / b.Editor.Zoom)
	} else {
		b.Pending.X += int(b.Delta.X / b.Editor.Zoom)
		b.Pending.Y += int(b.Delta.Y / b.Editor.Zoom)
	}
	b.Pending = b.Pending.Clamp(b.Editor.page().Size)
	b.Refresh()
}

// DragEnd publica um arraste inteiro como um comando reversível.
func (b *Board) DragEnd() {
	if !b.Dragging {
		return
	}
	b.Dragging = false
	// bounds é a geometria a confirmar.
	bounds := b.Pending
	// id identifica o alvo da operação.
	id := b.Editor.Selected
	b.Editor.check(b.Editor.Session.Change(func(document *application.Document) error {
		// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
		if page := document.Screen(b.Editor.PageID); page != nil {
			// component identifica o controle da tela que será consultado ou alterado.
			if component := page.ComponentByID(id); component != nil {
				component.Bounds = bounds
			}
		}
		return nil
	}))
	b.Editor.changed()
}

// logical remove margem e escala sem contaminar os dados.
func (b *Board) logical(position fyne.Position) (int, int) {
	return int((position.X - 24) / b.Editor.Zoom), int((position.Y - 24) / b.Editor.Zoom)
}

// boardRenderer conserva a lista de objetos Fyne do canvas.
type boardRenderer struct {
	// board identifica o widget que fornece geometria e seleção.
	board *Board
	// objects contém a representação visual atual.
	objects []fyne.CanvasObject
}

// Layout preserva posicionamento absoluto das primitivas.
func (r *boardRenderer) Layout(_ fyne.Size) {}

// MinSize delega dimensões ao artboard.
func (r *boardRenderer) MinSize() fyne.Size { return r.board.MinSize() }

// Objects fornece os objetos sem expor componentes de domínio.
func (r *boardRenderer) Objects() []fyne.CanvasObject { return r.objects }

// Destroy não possui recursos externos a liberar.
func (r *boardRenderer) Destroy() {}

// Refresh reconstitui somente o estado visual da revisão atual.
func (r *boardRenderer) Refresh() {
	// editor concentra referências da sessão visual.
	editor := r.board.Editor
	// objects começa com o fundo externo do canvas.
	objects := []fyne.CanvasObject{}
	// surface cria retângulos posicionados sem layout automático.
	surface := func(bounds screen.Bounds, fill, stroke string) {
		// rectangle é a primitiva Fyne já desacoplada do domínio.
		rectangle := canvas.NewRectangle(rgb(fill))
		if fill == "" {
			rectangle.FillColor = rgb("#FFFFFF")
		}
		if stroke != "" {
			rectangle.StrokeColor = rgb(stroke)
			rectangle.StrokeWidth = 1
		}
		rectangle.Move(fyne.NewPos(float32(bounds.X)*editor.Zoom+24, float32(bounds.Y)*editor.Zoom+24))
		rectangle.Resize(fyne.NewSize(float32(bounds.Width)*editor.Zoom, float32(bounds.Height)*editor.Zoom))
		objects = append(objects, rectangle)
	}
	// page é a composição ativa.
	page := editor.page()
	if page == nil {
		r.objects = []fyne.CanvasObject{widget.NewLabel("Adicione uma tela para começar.")}
		return
	}
	surface(screen.Bounds{Width: page.Size.Width, Height: page.Size.Height}, page.Background.Color, "#CBD5E1")
	if page.Background.AssetID != "" {
		objects = append(objects, r.image(page.Background.AssetID, screen.Bounds{Width: page.Size.Width, Height: page.Size.Height}))
	}
	// original identifica cada componente sem alterar o snapshot.
	for _, original := range page.Components {
		// component conserva uma cópia escalar para arraste provisório.
		component := original
		if component.ID == editor.Selected && r.board.Dragging {
			component.Bounds = r.board.Pending
		}
		// primitive é uma caixa, texto ou imagem editável.
		for _, primitive := range presentation.Component(component).Shapes {
			if primitive.AssetID != "" {
				objects = append(objects, r.image(primitive.AssetID, primitive.Bounds))
				continue
			}
			if primitive.Fill != "" || primitive.Stroke != "" {
				// rectangle usa transparência quando a forma não tem preenchimento.
				rectangle := canvas.NewRectangle(rgb(primitive.Fill))
				if primitive.Fill == "" {
					rectangle.FillColor = transparent()
				}
				if primitive.Stroke != "" {
					rectangle.StrokeColor = rgb(primitive.Stroke)
					rectangle.StrokeWidth = 1
				}
				rectangle.Move(fyne.NewPos(float32(primitive.Bounds.X)*editor.Zoom+24, float32(primitive.Bounds.Y)*editor.Zoom+24))
				rectangle.Resize(fyne.NewSize(float32(primitive.Bounds.Width)*editor.Zoom, float32(primitive.Bounds.Height)*editor.Zoom))
				objects = append(objects, rectangle)
			}
			if primitive.Text != "" {
				// text mantém a tipografia proporcional ao zoom.
				text := canvas.NewText(primitive.Text, rgb(primitive.Color))
				text.TextSize = float32(primitive.FontSize) * editor.Zoom
				text.TextStyle.Bold = primitive.Bold
				text.Move(fyne.NewPos(float32(primitive.Bounds.X)*editor.Zoom+24, float32(primitive.Bounds.Y)*editor.Zoom+24))
				objects = append(objects, text)
			}
		}
		if component.ID == editor.Selected {
			// selection delimita o componente sem cobrir seu conteúdo.
			selection := canvas.NewRectangle(transparent())
			selection.StrokeColor = rgb("#2563EB")
			selection.StrokeWidth = 2
			selection.Move(fyne.NewPos(float32(component.Bounds.X)*editor.Zoom+24, float32(component.Bounds.Y)*editor.Zoom+24))
			selection.Resize(fyne.NewSize(float32(component.Bounds.Width)*editor.Zoom, float32(component.Bounds.Height)*editor.Zoom))
			objects = append(objects, selection)
			surface(screen.Bounds{X: component.Bounds.X + component.Bounds.Width - 8, Y: component.Bounds.Y + component.Bounds.Height - 8, Width: 8, Height: 8}, "#2563EB", "")
		}
	}
	r.objects = objects
	// object recebe a atualização de seu renderer.
	for _, object := range r.objects {
		object.Refresh()
	}
}

// image resolve bytes incorporados sem acessar uma URL.
func (r *boardRenderer) image(id string, bounds screen.Bounds) fyne.CanvasObject {
	// image mostra um recurso local do documento.
	image := canvas.NewImageFromResource(nil)
	// asset é a descrição de imagem incorporada.
	for _, asset := range r.board.Editor.Session.Document.Project.Assets {
		if asset.ID == id {
			image = canvas.NewImageFromResource(fyne.NewStaticResource(asset.Path, r.board.Editor.Session.Document.Resources[asset.Path]))
		}
	}
	image.FillMode = canvas.ImageFillStretch
	image.Move(fyne.NewPos(float32(bounds.X)*r.board.Editor.Zoom+24, float32(bounds.Y)*r.board.Editor.Zoom+24))
	image.Resize(fyne.NewSize(float32(bounds.Width)*r.board.Editor.Zoom, float32(bounds.Height)*r.board.Editor.Zoom))
	return image
}

// transparent fornece preenchimento sem opacidade.
func transparent() colorTransparent { return colorTransparent{} }

// colorTransparent implementa a cor totalmente transparente.
type colorTransparent struct{}

// RGBA retorna canais transparentes para o renderer.
func (colorTransparent) RGBA() (uint32, uint32, uint32, uint32) { return 0, 0, 0, 0 }
