package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	desktopdriver "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// legacyControl mantém a interação Fyne sobre a mesma cena usada no canvas e no PPT.
type legacyControl struct {
	widget.DisableableWidget
	// document fornece somente os recursos do snapshot de prévia.
	document application.Document
	// component conserva o conteúdo tipado local, sem publicar alterações de interação.
	component screen.Component
	// onTapped executa a intenção do controle sem modificar sua decoração.
	onTapped func()
	// focused e pressed selecionam os estados visuais transitórios.
	focused, pressed bool
	// takeFocus preserva o cursor do campo ao clicar no teclado virtual.
	takeFocus bool
}

// newLegacyControl cria um controle cuja geometria interna começa em zero.
func newLegacyControl(document application.Document, component screen.Component, onTapped func()) *legacyControl {
	// control recebe uma cópia da instância para manter o runtime isolado.
	control := &legacyControl{document: document, component: component, onTapped: onTapped, takeFocus: true}
	control.component.Bounds.X, control.component.Bounds.Y = 0, 0
	control.ExtendBaseWidget(control)
	if !component.Enabled {
		control.Disable()
	}
	return control
}

// CreateRenderer reutiliza as primitivas de apresentação sem botão opaco sobreposto.
func (c *legacyControl) CreateRenderer() fyne.WidgetRenderer {
	// renderer mantém a cena abaixo do realce local de foco.
	renderer := &legacyControlRenderer{control: c, content: container.NewWithoutLayout()}
	renderer.Refresh()
	return renderer
}

// Tapped dispara somente controles habilitados e conserva foco por teclado.
func (c *legacyControl) Tapped(_ *fyne.PointEvent) {
	if c.Disabled() {
		return
	}
	if c.takeFocus {
		if surface := fyne.CurrentApp().Driver().CanvasForObject(c); surface != nil {
			surface.Focus(c)
		}
	}
	if c.onTapped != nil {
		c.onTapped()
	}
}

// FocusGained mostra o foco sem deslocar a composição.
func (c *legacyControl) FocusGained() { c.focused = true; c.Refresh() }

// FocusLost remove o estado de foco e qualquer pressão interrompida.
func (c *legacyControl) FocusLost() { c.focused, c.pressed = false, false; c.Refresh() }

// TypedRune não insere texto em controles acionáveis.
func (c *legacyControl) TypedRune(_ rune) {}

// TypedKey permite acionar o controle com Enter ou espaço.
func (c *legacyControl) TypedKey(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyReturn || event.Name == fyne.KeySpace {
		c.Tapped(nil)
	}
}

// MouseDown seleciona a imagem touched somente para o botão principal.
func (c *legacyControl) MouseDown(event *desktopdriver.MouseEvent) {
	if !c.Disabled() && event.Button == desktopdriver.MouseButtonPrimary {
		c.pressed = true
		c.Refresh()
	}
}

// MouseUp restaura a imagem normal; a ativação pertence ao evento Tapped.
func (c *legacyControl) MouseUp(_ *desktopdriver.MouseEvent) { c.pressed = false; c.Refresh() }

// legacyControlRenderer conserva a cena e seu indicador de foco no mesmo widget.
type legacyControlRenderer struct {
	// control fornece o conteúdo e o estado transitório.
	control *legacyControl
	// content agrupa as partes derivadas da cena.
	content *fyne.Container
}

// Layout acompanha somente o tamanho explícito atribuído ao controle.
func (r *legacyControlRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
	if r.control.component.Bounds.Width != int(size.Width) || r.control.component.Bounds.Height != int(size.Height) {
		r.control.component.Bounds.Width, r.control.component.Bounds.Height = int(size.Width), int(size.Height)
		r.Refresh()
	}
}

// MinSize permite os bounds do projeto sem impor o mínimo de um botão temático.
func (r *legacyControlRenderer) MinSize() fyne.Size { return fyne.NewSize(1, 1) }

// Objects expõe o contêiner que mantém a ordem interna da cena.
func (r *legacyControlRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.content} }

// Destroy não possui recursos externos ou goroutines para encerrar.
func (r *legacyControlRenderer) Destroy() {}

// Refresh resolve o estado original e acrescenta apenas o realce de foco local.
func (r *legacyControlRenderer) Refresh() {
	// state seleciona a imagem original pressionada sem gravá-la no documento.
	state := "default"
	if r.control.pressed {
		state = "touched"
	}
	r.content.Objects = drawGroup(r.control.document, presentation.ComponentState(r.control.component, state))
	if r.control.focused && !r.control.Disabled() {
		// outline é decoração de interação, ausente do canvas estático e do PPT.
		outline := canvas.NewRectangle(transparent())
		outline.StrokeColor, outline.StrokeWidth = rgb("#F59E0B"), 1
		place(outline, r.control.component.Bounds)
		r.content.Objects = append(r.content.Objects, outline)
	}
	r.content.Refresh()
}
