package desktop

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/screen"
)

// Editor conserva apenas sessão e seleção da interface.
type Editor struct {
	// App controla eventos e janelas Fyne.
	App fyne.App
	// Window é a janela principal.
	Window fyne.Window
	// Session coordena o documento e seu histórico.
	Session *application.Session
	// PageID identifica a tela em edição.
	PageID string
	// Selected identifica um componente selecionado.
	Selected string
	// FilePath é o destino explícito atual.
	FilePath string
	// Zoom afeta somente o desenho.
	Zoom float32
	// Board desenha e manipula os componentes.
	Board *Board
	// Screens exibe a navegação entre telas.
	Screens *fyne.Container
	// Inspector exibe propriedades da seleção.
	Inspector *fyne.Container
	// Structure relaciona os controles da tela e permite selecioná-los por nome.
	Structure *fyne.Container
	// Drawing mantém a área disponível para ajustar o zoom à janela.
	Drawing *container.Scroll
	// ZoomLabel informa a escala efetiva sem um valor fixo enganoso.
	ZoomLabel *widget.Label
	// Status informa revisões, falhas e resultados.
	Status *widget.Label
	// Title exibe o nome do projeto no cabeçalho.
	Title *canvas.Text
	// Tabs alterna editor, navegação e roteiro.
	Tabs *container.AppTabs
	// JourneyPanel atualiza destinos e tela inicial.
	JourneyPanel *fyne.Container
	// ScenarioPanel atualiza as etapas explícitas.
	ScenarioPanel *fyne.Container
	// Autosave agenda cópia local de recuperação.
	Autosave *time.Timer
	// ExportCancel cancela o job ativo.
	ExportCancel func()
}

// New monta a interface e mantém composição de dependências fora dos domínios.
func New(app fyne.App, document application.Document) *Editor {
	// editor recebe uma sessão completamente isolada.
	editor := &Editor{App: app, Window: app.NewWindow("Gerador de Jornada"), Session: application.NewSession(document), Zoom: 1}
	if len(document.Screens) > 0 {
		editor.PageID = document.Screens[0].ID
	}
	app.Settings().SetTheme(jornadaTheme{})
	editor.Status = widget.NewLabel("Pronto · funcionamento offline")
	editor.Title = canvas.NewText(document.Project.Name, rgb("#0F172A"))
	editor.Title.TextSize = 24
	editor.Title.TextStyle.Bold = true
	editor.Screens = container.NewVBox()
	editor.Inspector = container.NewVBox()
	editor.Structure = container.NewVBox()
	editor.ZoomLabel = widget.NewLabel("100%")
	editor.JourneyPanel = container.NewVBox()
	editor.ScenarioPanel = container.NewVBox()
	editor.Board = NewBoard(editor)
	editor.Tabs = container.NewAppTabs(container.NewTabItemWithIcon("Desenhar", theme.DocumentCreateIcon(), editor.design()), container.NewTabItemWithIcon("Navegação", theme.NavigateNextIcon(), container.NewVScroll(editor.JourneyPanel)), container.NewTabItemWithIcon("Demonstração", theme.MediaPlayIcon(), container.NewVScroll(editor.ScenarioPanel)))
	editor.Window.SetContent(container.NewBorder(editor.header(), container.NewPadded(editor.Status), nil, nil, editor.Tabs))
	editor.Window.Resize(fyne.NewSize(1280, 720))
	editor.Window.CenterOnScreen()
	editor.installShortcuts()
	editor.refresh()
	editor.Window.SetCloseIntercept(func() {
		if editor.Session.Dirty() {
			dialog.ShowConfirm("Alterações não salvas", "Sair sem salvar as alterações?", func(yes bool) {
				if yes {
					editor.close()
				}
			}, editor.Window)
		} else {
			editor.close()
		}
	})
	return editor
}

// close encerra timers e libera a janela sem alterar dados persistidos.
func (e *Editor) close() {
	if e.Autosave != nil {
		e.Autosave.Stop()
	}
	if e.ExportCancel != nil {
		e.ExportCancel()
	}
	e.Window.SetCloseIntercept(nil)
	e.Window.Close()
}

// page encontra a composição em edição.
func (e *Editor) page() *screen.Screen { return e.Session.Document.Screen(e.PageID) }

// header reúne identidade, ações e informações essenciais sem poluir o canvas.
func (e *Editor) header() fyne.CanvasObject {
	// brand identifica o produto, separado do título da jornada.
	brand := canvas.NewText("JORNADA STUDIO", rgb("#2563EB"))
	brand.TextSize = 12
	brand.TextStyle.Bold = true
	// subtitle orienta a tarefa principal.
	subtitle := canvas.NewText("Desenhe experiências. Conecte telas. Apresente sua jornada.", rgb("#64748B"))
	subtitle.TextSize = 13
	// identity mantém hierarquia tipográfica clara.
	identity := container.NewVBox(brand, e.Title, subtitle)
	// actions disponibiliza fluxos completos de arquivo e exportação.
	actions := container.NewHBox(widget.NewButtonWithIcon("Novo", theme.DocumentCreateIcon(), e.newProject), widget.NewButtonWithIcon("Abrir", theme.FolderOpenIcon(), e.openProject), widget.NewButtonWithIcon("Salvar", theme.DocumentSaveIcon(), func() { e.saveProject(false) }), widget.NewButtonWithIcon("Prévia", theme.VisibilityIcon(), e.preview), widget.NewButtonWithIcon("PowerPoint", theme.MediaPlayIcon(), e.exportPPT), widget.NewButtonWithIcon("Java / EFL", theme.ComputerIcon(), e.exportJava))
	return container.NewPadded(container.NewVBox(container.NewBorder(nil, nil, identity, nil, container.NewVBox(layout.NewSpacer(), actions, layout.NewSpacer())), widget.NewSeparator()))
}

// design organiza paleta, área de desenho e inspector em painéis redimensionáveis.
func (e *Editor) design() fyne.CanvasObject {
	// left organiza telas e componentes sem mistura de propriedades.
	left := e.toolsPanel()
	// controls contém ações de histórico, zoom e organização.
	controls := container.NewHBox(widget.NewButtonWithIcon("", theme.ContentUndoIcon(), func() { e.check(e.Session.Undo()); e.changed() }), widget.NewButtonWithIcon("", theme.ContentRedoIcon(), func() { e.check(e.Session.Redo()); e.changed() }), widget.NewSeparator(), widget.NewButton("−", func() { e.setZoom(e.Zoom - .25) }), e.ZoomLabel, widget.NewButton("+", func() { e.setZoom(e.Zoom + .25) }), widget.NewButton("Ajustar", e.fitZoom), layout.NewSpacer(), widget.NewButtonWithIcon("Fundo", theme.FileImageIcon(), e.configureBackground))
	e.Drawing = container.NewScroll(e.Board)
	// drawing preserva a área de desenho em uma superfície rolável.
	drawing := container.NewBorder(container.NewPadded(controls), nil, nil, nil, e.Drawing)
	// right contém o inspector com largura mínima agradável.
	right := container.NewBorder(section("PROPRIEDADES"), nil, nil, nil, container.NewVScroll(e.Inspector))
	// main separa canvas e propriedades.
	main := container.NewHSplit(drawing, right)
	main.Offset = .72
	// split reserva navegação à esquerda.
	split := container.NewHSplit(container.NewPadded(left), main)
	split.Offset = .20
	return split
}

// section cria um rótulo discreto de painel.
func section(text string) fyne.CanvasObject {
	// label mantém contraste sem competir com o conteúdo.
	label := canvas.NewText(text, color.NRGBA{R: 100, G: 116, B: 139, A: 255})
	label.TextSize = 11
	label.TextStyle.Bold = true
	return container.NewPadded(label)
}

// setZoom altera somente a escala da visualização.
func (e *Editor) setZoom(value float32) {
	e.Zoom = min(max(.25, value), 4)
	e.ZoomLabel.SetText(fmt.Sprintf("%.0f%%", e.Zoom*100))
	e.Board.Refresh()
	e.Drawing.Refresh()
	e.Status.SetText(fmt.Sprintf("Zoom %.0f%% · coordenadas lógicas preservadas", e.Zoom*100))
}

// fitZoom mostra a superfície inteira usando somente a área disponível no editor.
func (e *Editor) fitZoom() {
	if e.page() == nil || e.Drawing == nil {
		return
	}
	// available desconta as margens do artboard e a barra de rolagem.
	available := e.Drawing.Size().Subtract(fyne.NewSize(64, 64))
	if available.Width > 0 && available.Height > 0 {
		e.setZoom(min(available.Width/float32(e.page().Size.Width), available.Height/float32(e.page().Size.Height)))
		e.Drawing.ScrollToTop()
	}
}

// FitToWindow ajusta a tela após o layout inicial da janela, sem mudar dados do projeto.
func (e *Editor) FitToWindow() { e.fitZoom() }

// check apresenta uma falha real sem suprimir seu contexto.
func (e *Editor) check(err error) bool {
	if err != nil {
		dialog.ShowError(err, e.Window)
		return false
	}
	return true
}

// changed atualiza a UI após um comando confirmado e agenda recuperação.
func (e *Editor) changed() { e.Selected = validSelection(e); e.refresh(); e.scheduleRecovery() }

// validSelection elimina apenas a seleção de um objeto já removido.
func validSelection(e *Editor) string {
	// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
	if page := e.page(); page != nil && page.ComponentByID(e.Selected) != nil {
		return e.Selected
	}
	return ""
}

// refresh sincroniza painéis com o estado de aplicação atual.
func (e *Editor) refresh() {
	e.Title.Text = e.Session.Document.Project.Name
	e.Title.Refresh()
	e.refreshScreens()
	e.refreshStructure()
	e.refreshInspector()
	e.refreshJourney()
	e.refreshScenarios()
	e.Board.Refresh()
	e.Drawing.Refresh()
	e.Status.SetText(fmt.Sprintf("%d telas · revisão %d · zoom %.0f%%", len(e.Session.Document.Screens), e.Session.Document.Project.Revision, e.Zoom*100))
}

// refreshScreens reconstrói a lista sem persistir seleção ou foco.
func (e *Editor) refreshScreens() {
	e.Screens.RemoveAll()
	// page representa uma tela selecionável.
	for _, page := range e.Session.Document.Screens {
		// id permanece estável mesmo quando a lista é reordenada.
		id := page.ID
		// button alterna a composição em edição.
		button := widget.NewButton(page.Name, func() { e.PageID = id; e.Selected = ""; e.refresh() })
		if id == e.PageID {
			button.Importance = widget.HighImportance
		}
		e.Screens.Add(button)
	}
	e.Screens.Add(container.NewHBox(widget.NewButtonWithIcon("", theme.ContentCopyIcon(), e.duplicateScreen), widget.NewButtonWithIcon("", theme.DeleteIcon(), e.removeScreen), widget.NewButton("↑", func() { e.reorderScreen(-1) }), widget.NewButton("↓", func() { e.reorderScreen(1) })))
}
