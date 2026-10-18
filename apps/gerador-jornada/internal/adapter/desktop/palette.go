package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/catalog"
)

// paletteItem permite clicar para inserir ou arrastar para o artboard.
type paletteItem struct {
	widget.Button
	// Editor resolve o destino sem depender de coordenadas globais do Windows.
	Editor *Editor
	// Kind seleciona a família do catálogo.
	Kind catalog.Kind
	// PresetID seleciona um conjunto auxiliar do legado em vez de uma família isolada.
	PresetID string
	// LastPosition conserva a última posição absoluta do ponteiro.
	LastPosition fyne.Position
	// Dragging distingue um clique simples de um arraste real.
	Dragging bool
}

// palette monta as sete famílias com ícones e nomes legíveis.
func (e *Editor) palette() fyne.CanvasObject {
	// items mantém a ordem do catálogo.
	items := container.NewVBox()
	// definition representa uma família arrastável.
	for _, definition := range catalog.Definitions() {
		// item encapsula o comportamento visual da paleta.
		item := &paletteItem{Editor: e, Kind: definition.Kind}
		item.Text = definition.Name
		item.Icon = theme.ContentAddIcon()
		item.Alignment = widget.ButtonAlignLeading
		item.OnTapped = func() { e.insert(item.Kind, 48, 48) }
		item.ExtendBaseWidget(item)
		items.Add(item)
	}
	return items
}

// Dragged memoriza o destino do ponteiro no canvas da janela.
func (p *paletteItem) Dragged(event *fyne.DragEvent) {
	p.LastPosition = event.AbsolutePosition
	p.Dragging = true
}

// DragEnd insere somente quando o destino pertence à área lógica da tela.
func (p *paletteItem) DragEnd() {
	if !p.Dragging {
		return
	}
	p.Dragging = false
	// origin representa a posição do artboard na janela.
	origin := fyne.CurrentApp().Driver().AbsolutePositionForObject(p.Editor.Board)
	// x e y são coordenadas lógicas independentes do zoom.
	x, y := p.Editor.Board.logical(p.LastPosition.Subtract(origin))
	// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
	if page := p.Editor.page(); page != nil && x >= 0 && y >= 0 && x < page.Size.Width && y < page.Size.Height {
		if p.PresetID != "" {
			p.Editor.insertPreset(p.PresetID, x, y)
		} else {
			p.Editor.insert(p.Kind, x, y)
		}
	}
}

// insert confirma a inclusão e seleciona o novo controle.
func (e *Editor) insert(kind catalog.Kind, x, y int) {
	// id é atribuído pelo domínio, não pelo widget.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	id, err := e.Session.AddComponent(e.PageID, kind, x, y)
	if !e.check(err) {
		return
	}
	e.Selected = id
	e.changed()
}
