package desktop

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// legacyEntry preserva seleção, cursor e clipboard do Entry sobre uma cena EFL.
// O renderer nativo pinta somente texto/cursor, com fundo e borda transparentes.
type legacyEntry struct {
	widget.Entry
	// document fornece os recursos da prévia isolada.
	document application.Document
	// component contém a configuração confirmada do campo.
	component screen.Component
	// focused seleciona default/focused sem persistir foco no projeto.
	focused bool
}

// newLegacyEntry cria a captura de texto com as mesmas fontes e bounds da cena.
func newLegacyEntry(document application.Document, component screen.Component, value string) (*legacyEntry, fyne.CanvasObject) {
	// field conserva a edição transitória enquanto a sessão valida o valor.
	field := &legacyEntry{document: document, component: component}
	field.component.Bounds.X, field.component.Bounds.Y = 0, 0
	field.ExtendBaseWidget(field)
	field.SetText(value)
	// surface aplica tipografia apenas ao campo, sem alterar o tema do editor.
	surface := container.NewThemeOverride(field, legacyEntryTheme{field: field})
	place(surface, component.Bounds)
	return field, surface
}

// FocusGained altera o fundo EFL e conserva a implementação nativa do cursor.
func (e *legacyEntry) FocusGained() { e.focused = true; e.Entry.FocusGained(); e.Refresh() }

// FocusLost restaura o estado visual normal sem alterar o valor da visita.
func (e *legacyEntry) FocusLost() { e.focused = false; e.Entry.FocusLost(); e.Refresh() }

// state informa somente o estado transitório de apresentação.
func (e *legacyEntry) state() string {
	if !e.component.Enabled {
		return "disabled"
	}
	if e.focused {
		return "focused"
	}
	return "default"
}

// textPrimitive resolve fonte e cor no mesmo caminho usado pelo canvas e pelo PPT.
func (e *legacyEntry) textPrimitive() presentation.Primitive {
	// primitive identifica o valor animável sem depender da quantidade de fundos/bordas.
	for _, primitive := range presentation.ComponentState(e.component, e.state()).Shapes {
		if primitive.Role == "value" {
			return primitive
		}
	}
	return presentation.Primitive{FontSize: 12, Color: "#000000", FontAsset: "/fonts/arial.ttf"}
}

// CreateRenderer combina a decoração EFL com a captura nativa de texto transparente.
func (e *legacyEntry) CreateRenderer() fyne.WidgetRenderer {
	return &legacyEntryRenderer{field: e, text: e.Entry.CreateRenderer(), background: container.NewWithoutLayout()}
}

// legacyEntryRenderer mantém a decoração sob o texto, a seleção e o cursor.
type legacyEntryRenderer struct {
	// field fornece estado local e geometria explícita.
	field *legacyEntry
	// text mantém as regras de edição e acessibilidade do Fyne.
	text fyne.WidgetRenderer
	// background desenha somente as partes sem o papel value.
	background *fyne.Container
}

// Layout mantém o campo no tamanho persistido, sem mínimo de altura temático.
func (r *legacyEntryRenderer) Layout(size fyne.Size) {
	r.field.component.Bounds.Width, r.field.component.Bounds.Height = int(size.Width), int(size.Height)
	r.background.Resize(size)
	r.text.Layout(size)
	r.Refresh()
}

// MinSize deixa o contêiner externo respeitar os limites lógicos do documento.
func (r *legacyEntryRenderer) MinSize() fyne.Size { return fyne.NewSize(1, 1) }

// Objects coloca a decoração atrás da edição nativa transparente.
func (r *legacyEntryRenderer) Objects() []fyne.CanvasObject {
	return append([]fyne.CanvasObject{r.background}, r.text.Objects()...)
}

// Refresh atualiza a cena conforme o foco sem duplicar o texto do valor.
func (r *legacyEntryRenderer) Refresh() {
	// group elimina somente o valor; bordas e fundos conservam a ordem da cena.
	group := presentation.ComponentState(r.field.component, r.field.state())
	// parts contém a decoração, nunca uma segunda cópia do conteúdo digitado.
	parts := []presentation.Primitive{}
	// primitive percorre as partes da configuração confirmada.
	for _, primitive := range group.Shapes {
		if primitive.Role != "value" {
			parts = append(parts, primitive)
		}
	}
	group.Shapes = parts
	r.background.Objects = drawGroup(r.field.document, group)
	r.background.Refresh()
	r.text.Refresh()
}

// Destroy delega a liberação dos recursos de captura de texto ao renderer nativo.
func (r *legacyEntryRenderer) Destroy() { r.text.Destroy() }

// legacyEntryTheme permite desenhar texto nativo sobre a decoração da cena.
type legacyEntryTheme struct {
	// field informa a face/cor atuais, incluindo o estado focused.
	field *legacyEntry
}

// Color deixa o fundo transparente e usa a cor da mesma primitiva compartilhada.
func (t legacyEntryTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameInputBackground, theme.ColorNameInputBorder, theme.ColorNameBackground:
		return color.Transparent
	case theme.ColorNameForeground, theme.ColorNameDisabled, theme.ColorNamePrimary, theme.ColorNamePlaceHolder:
		return rgb(t.field.textPrimitive().Color)
	case theme.ColorNameSelection:
		return color.NRGBA{R: 128, G: 176, B: 224, A: 160}
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

// Font devolve os bytes reais da fonte declarada no campo.
func (t legacyEntryTheme) Font(style fyne.TextStyle) fyne.Resource {
	if resource := sceneFont(t.field.textPrimitive().FontAsset); resource != nil {
		return resource
	}
	return theme.DefaultTheme().Font(style)
}

// Icon conserva somente os ícones funcionais do capturador de texto.
func (t legacyEntryTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size elimina bordas/arredondamentos nativos e usa o tamanho lógico do texto da cena.
func (t legacyEntryTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return float32(t.field.textPrimitive().FontSize)
	case theme.SizeNameInputRadius, theme.SizeNameInputBorder:
		return 0
	case theme.SizeNameInnerPadding:
		return 3
	}
	return theme.DefaultTheme().Size(name)
}
