package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
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
	// Originals conserva os bounds de cada membro antes de mover a seleção.
	Originals map[string]screen.Bounds
	// GroupBounds limita a translação do conjunto sem comprimir seu espaçamento.
	GroupBounds screen.Bounds
	// OffsetX e OffsetY expressam a translação provisória em unidades lógicas.
	OffsetX, OffsetY int
	// EmptyAnchor inicia uma caixa de seleção ao arrastar uma área vazia.
	EmptyAnchor bool
	// Marquee indica desenho de uma caixa de seleção, sem edição do documento.
	Marquee bool
	// Anchor e Area conservam a origem e o retângulo provisório de seleção.
	Anchor fyne.Position
	Area   screen.Bounds
	// Additive conserva a seleção anterior quando Ctrl ou Shift estava pressionado.
	Additive []string
	// Pointer identifica o cursor contextual indicado ao driver desktop.
	Pointer desktop.Cursor
	// ContextMenu conserva o menu aberto sobre o componente, sem persistir esse estado.
	ContextMenu *widget.PopUpMenu
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
	// hit identifica o controle superior sem reconstruir propriedades durante o teste geométrico.
	hit := b.componentAt(event.Position)
	// ids conserva a seleção anterior para alternância e caixas aditivas.
	ids := append([]string{}, b.Editor.selectionIDs()...)
	b.Additive = nil
	// additive distingue alternância explícita do clique para iniciar movimentação coletiva.
	additive := event.Modifier&(fyne.KeyModifierControl|fyne.KeyModifierShift) != 0
	if additive {
		b.Additive = append([]string{}, ids...)
	}
	if hit == "" {
		if !additive {
			ids = nil
		}
	} else if additive {
		// next alterna somente o alvo clicado, preservando os demais controles.
		next := []string{}
		for _, id := range ids {
			if id != hit {
				next = append(next, id)
			}
		}
		if !b.Editor.isSelected(hit) {
			next = append(next, hit)
		}
		ids = next
	} else if !b.Editor.isSelected(hit) {
		ids = []string{hit}
	}
	// changed evita reconstruir formulários ao clicar novamente no conjunto já ativo.
	changed := b.Editor.setSelection(ids)
	b.EmptyAnchor = hit == ""
	b.Anchor = fyne.NewPos(float32(x), float32(y))
	b.Marquee, b.Resizing = false, false
	b.OffsetX, b.OffsetY = 0, 0
	b.Originals = map[string]screen.Bounds{}
	for _, id := range b.Editor.selectionIDs() {
		b.Originals[id] = page.ComponentByID(id).Bounds
	}
	if len(b.Originals) > 0 {
		b.GroupBounds, _ = screen.SelectionBounds(*page, b.Editor.selectionIDs())
		b.Original = page.ComponentByID(b.Editor.Selected).Bounds
		b.Pending = b.Original
		b.Resizing = len(b.Originals) == 1 && hit != "" && x > b.Original.X+b.Original.Width-12 && y > b.Original.Y+b.Original.Height-12
	}
	b.Delta = fyne.Position{}
	b.Dragging = false
	if changed {
		b.Editor.selectionChanged()
	}
}

// MouseUp não confirma duas vezes o mesmo arraste.
func (b *Board) MouseUp(_ *desktop.MouseEvent) {}

// Dragged atualiza apenas a geometria provisória do canvas.
func (b *Board) Dragged(event *fyne.DragEvent) {
	if b.Editor.page() == nil {
		return
	}
	b.Dragging = true
	b.Delta = b.Delta.Add(fyne.NewPos(event.Dragged.DX, event.Dragged.DY))
	if b.EmptyAnchor {
		b.Marquee = true
		// end descreve o canto móvel em coordenadas lógicas, inclusive para arraste invertido.
		end := b.Anchor.Add(fyne.NewPos(b.Delta.X/b.Editor.Zoom, b.Delta.Y/b.Editor.Zoom))
		b.Area = screen.Bounds{X: int(min(b.Anchor.X, end.X)), Y: int(min(b.Anchor.Y, end.Y)), Width: max(1, int(abs(end.X-b.Anchor.X))), Height: max(1, int(abs(end.Y-b.Anchor.Y)))}
		b.Pointer = desktop.CrosshairCursor
		b.Refresh()
		return
	}
	if len(b.Originals) == 0 {
		return
	}
	b.Pending = b.Original
	if b.Resizing {
		b.Pending.Width += int(b.Delta.X / b.Editor.Zoom)
		b.Pending.Height += int(b.Delta.Y / b.Editor.Zoom)
	} else {
		b.OffsetX, b.OffsetY = screen.ConstrainTranslation(b.GroupBounds, b.Editor.page().Size, int(b.Delta.X/b.Editor.Zoom), int(b.Delta.Y/b.Editor.Zoom))
		b.Pending.X += b.OffsetX
		b.Pending.Y += b.OffsetY
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
	if b.Marquee {
		b.Marquee = false
		// ids junta a seleção anterior e os alvos que intersectam a caixa.
		ids := append(append([]string{}, b.Additive...), screen.IntersectingIDs(*b.Editor.page(), b.Area)...)
		if b.Editor.setSelection(ids) {
			b.Editor.selectionChanged()
		}
		b.Refresh()
		return
	}
	if !b.Resizing {
		if b.OffsetX != 0 || b.OffsetY != 0 {
			b.Editor.check(b.Editor.Session.MoveComponents(b.Editor.PageID, b.Editor.selectionIDs(), b.OffsetX, b.OffsetY))
			b.Editor.changed()
		} else {
			b.Refresh()
		}
		return
	}
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

// abs calcula a magnitude do deslocamento do ponteiro na seleção retangular.
func abs(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}

// Cursor fornece o cursor contextual, consultado pelo driver a cada movimento.
func (b *Board) Cursor() desktop.Cursor {
	if b.Pointer != nil {
		return b.Pointer
	}
	return desktop.DefaultCursor
}

// MouseIn atualiza a intenção visual ao entrar no artboard.
func (b *Board) MouseIn(event *desktop.MouseEvent) { b.MouseMoved(event) }

// MouseOut restaura o cursor padrão ao sair da área de desenho.
func (b *Board) MouseOut() { b.Pointer = desktop.DefaultCursor }

// MouseMoved distingue movimentação em cruz e redimensionamento diagonal sem reconstruir formas.
func (b *Board) MouseMoved(event *desktop.MouseEvent) {
	b.Pointer = desktop.DefaultCursor
	if b.Editor.page() == nil {
		return
	}
	// x e y são coordenadas lógicas usadas somente no hit-test.
	x, y := b.logical(event.Position)
	for _, id := range b.Editor.selectionIDs() {
		// bounds pertence ao controle atual, sem dependência de widgets de propriedades.
		bounds := b.Editor.page().ComponentByID(id).Bounds
		if x >= bounds.X && x <= bounds.X+bounds.Width && y >= bounds.Y && y <= bounds.Y+bounds.Height {
			b.Pointer = desktop.CrosshairCursor
			if len(b.Editor.Selection) == 1 && x > bounds.X+bounds.Width-12 && y > bounds.Y+bounds.Height-12 {
				b.Pointer = desktop.NWSEResizeCursor
			}
			return
		}
	}
}

// logical remove margem e escala sem contaminar os dados.
func (b *Board) logical(position fyne.Position) (int, int) {
	return int((position.X - 24) / b.Editor.Zoom), int((position.Y - 24) / b.Editor.Zoom)
}
