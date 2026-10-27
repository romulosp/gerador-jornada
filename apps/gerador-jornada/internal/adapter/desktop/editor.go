package desktop

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/buildinfo"
	"gerador-jornada.local/editor/internal/adapter/gitworkspace"
	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/screen"
)

// Editor conserva apenas sessão e seleção da interface.
type Editor struct {
	// GitDiagnostics reutiliza o log local para rastrear etapas sem material confidencial.
	GitDiagnostics *syncdiag.Logger
	// GitStatusText mantém detalhes em diálogo sem consumir a altura da árvore.
	GitStatusText string
	// GitVerifiedAt registra somente a última consulta remota bem-sucedida desta sessão.
	GitVerifiedAt time.Time
	// GitSync coordena casos de uso sem expor credenciais à interface.
	GitSync *application.WorkspaceSync
	// GitHub e GitWorkspace mantêm integração opcional fora dos documentos.
	GitHub       *application.GitHubConnection
	GitWorkspace *gitworkspace.Engine
	// GitLabel e ações representam somente estado administrativo da conexão.
	GitLabel              *widget.Label
	GitConnectButton      *widget.Button
	GitPublishButton      *widget.Button
	GitRepositoriesButton *widget.Button
	GitDisconnectButton   *widget.Button
	GitRefreshButton      *widget.Button
	// GitOpenButton abre apenas o destino já vinculado no navegador padrão.
	GitOpenButton *widget.Button
	// GitBusy e GitCancel protegem a operação assíncrona da janela atual.
	GitBusy   bool
	GitCancel func()
	// GitRecoveredRoot evita recuperar/recarregar o mesmo documento em cada refresh.
	GitRecoveredRoot string
	GitRecoveryError error
	// GitClosed impede atualizar widgets após fechar a janela.
	GitClosed bool
	// App controla eventos e janelas Fyne.
	App fyne.App
	// Window é a janela principal.
	Window fyne.Window
	// Mode mantém tela cheia e integração nativa fora dos domínios e do documento persistido.
	Mode *windowMode
	// Session coordena o documento e seu histórico.
	Session *application.Session
	// PageID identifica a tela em edição.
	PageID string
	// Selected identifica um componente selecionado.
	Selected string
	// Selection contém IDs selecionados somente na interface, em ordem estável.
	Selection []string
	// DeleteButton permite excluir somente quando há componentes válidos selecionados.
	DeleteButton *widget.Button
	// DeleteDialog impede confirmações simultâneas para a mesma operação de exclusão.
	DeleteDialog *dialog.ConfirmDialog
	// CurrentScreenPreview abre a tela ativa e fica desabilitado quando não há tela existente.
	CurrentScreenPreview *widget.Button
	// FilePath é o destino explícito atual.
	FilePath string
	// WorkspacePath identifica a pasta escolhida enquanto a jornada ainda não tem caminho salvo.
	WorkspacePath string
	// WorkspaceTree e WorkspaceEntries conservam a árvore de arquivos, separada das camadas da tela.
	WorkspaceTree    *widget.Tree
	WorkspaceEntries []WorkspaceEntry
	// Zoom afeta somente o desenho.
	Zoom float32
	// Board desenha e manipula os componentes.
	Board *Board
	// Screens exibe a navegação entre telas.
	Screens *fyne.Container
	// Inspector exibe propriedades da seleção.
	Inspector *fyne.Container
	// InspectorActions mantém Salvar alterações e Descartar fora da área rolável.
	InspectorActions *fyne.Container
	// InspectorHasDraft detecta valores ainda não aplicados no formulário atual, sem editar o domínio.
	InspectorHasDraft func() bool
	// Structure relaciona os controles da tela e permite selecioná-los por nome.
	Structure *fyne.Container
	// StructurePage e StructureRevision identificam o conteúdo já montado na árvore.
	StructurePage     string
	StructureRevision uint64
	// Drawing mantém a área disponível para ajustar o zoom à janela.
	Drawing *container.Scroll
	// ZoomLabel informa a escala efetiva sem um valor fixo enganoso.
	ZoomLabel *widget.Label
	// Status informa revisões, falhas e resultados.
	Status *widget.Label
	// Title exibe o nome do projeto no cabeçalho.
	Title *canvas.Text
	// Brand apresenta a identificação do produto com a cor do tema.
	Brand *canvas.Text
	// Subtitle acompanha o contraste do tema sem alterar o conteúdo do projeto.
	Subtitle *canvas.Text
	// ThemeSelect oferece preferências locais de aparência do editor.
	ThemeSelect *widget.Select
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
	// O tema precisa estar instalado antes de registrar uma janela sem viewport no driver em execução.
	prepareApplicationTheme(app, app.Preferences().String(themePreferenceKey))
	// editor recebe uma sessão completamente isolada.
	editor := &Editor{App: app, Window: app.NewWindow("Gerador de Jornada " + buildinfo.Version()), Session: application.NewSession(document), Zoom: 1}
	// A janela principal encerra o aplicativo também quando existem prévias abertas.
	editor.Window.SetMaster()
	editor.Mode = newWindowMode(editor.Window)
	editor.Mode.exitButton.OnTapped = editor.requestClose
	if len(document.Screens) > 0 {
		// A tela inicial da jornada abre primeiro; projetos antigos sem referência válida usam a primeira tela.
		editor.PageID = document.Journey.StartScreenID
		if editor.PageID == "" || document.Screen(editor.PageID) == nil {
			editor.PageID = document.Screens[0].ID
		}
	}
	editor.Status = widget.NewLabel("Pronto · funcionamento offline")
	editor.Title = canvas.NewText(document.Project.Name, app.Settings().Theme().Color(theme.ColorNameForeground, theme.VariantLight))
	editor.Title.TextSize = 24
	editor.Title.TextStyle.Bold = true
	editor.initTheme()
	editor.Screens = container.NewVBox()
	editor.Inspector = container.NewVBox()
	editor.InspectorActions = container.NewVBox()
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
	editor.Window.SetCloseIntercept(editor.requestClose)
	return editor
}

// close encerra timers e libera a janela sem alterar dados persistidos.
func (e *Editor) close() {
	e.GitClosed = true
	if e.GitCancel != nil {
		e.GitCancel()
	}
	e.Mode.close()
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
	brand := canvas.NewText("JORNADA STUDIO · "+buildinfo.Version(), e.App.Settings().Theme().Color(theme.ColorNamePrimary, theme.VariantLight))
	e.Brand = brand
	brand.TextSize = 12
	brand.TextStyle.Bold = true
	// subtitle orienta a tarefa principal.
	subtitle := canvas.NewText("Desenhe experiências. Conecte telas. Apresente sua jornada.", e.App.Settings().Theme().Color(theme.ColorNamePlaceHolder, theme.VariantLight))
	e.Subtitle = subtitle
	subtitle.TextSize = 13
	// identity mantém hierarquia tipográfica clara.
	identity := container.NewVBox(e.Title, subtitle)
	// windowControls mantém aparência e encerramento acessíveis sem decoração nativa.
	windowControls := container.NewHBox(brand, layout.NewSpacer(), widget.NewButtonWithIcon("Log de aplicação", theme.InfoIcon(), e.showApplicationLog), widget.NewLabel("Tema"), e.ThemeSelect, e.Mode.button, e.Mode.exitButton)
	// actions disponibiliza fluxos completos de arquivo e exportação.
	// java mantém a capacidade identificável enquanto seu acesso está suspenso pelo usuário.
	java := widget.NewButtonWithIcon("Java / EFL", theme.ComputerIcon(), e.exportJava)
	java.Disable()
	e.CurrentScreenPreview = widget.NewButtonWithIcon("Prévia Tela Atual", theme.VisibilityIcon(), e.previewCurrentScreen)
	actions := container.NewHBox(widget.NewButtonWithIcon("Novo", theme.DocumentCreateIcon(), e.newProject), widget.NewButtonWithIcon("Abrir", theme.FolderOpenIcon(), e.openProject), widget.NewButtonWithIcon("Salvar", theme.DocumentSaveIcon(), func() { e.saveProject(false) }), widget.NewButtonWithIcon("Prévia", theme.VisibilityIcon(), e.preview), e.CurrentScreenPreview, widget.NewButtonWithIcon("PowerPoint", theme.MediaPlayIcon(), e.exportPPT), java)
	return container.NewPadded(container.NewVBox(windowControls, container.NewBorder(nil, nil, identity, nil, container.NewVBox(layout.NewSpacer(), actions, layout.NewSpacer())), widget.NewSeparator()))
}

// design organiza paleta, área de desenho e inspector em painéis redimensionáveis.
func (e *Editor) design() fyne.CanvasObject {
	// left organiza telas e componentes sem mistura de propriedades.
	left := e.toolsPanel()
	// DeleteButton compartilha a confirmação usada pelo menu contextual e pela tecla Delete.
	e.DeleteButton = widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), e.removeComponent)
	e.DeleteButton.Disable()
	// controls contém ações de histórico, zoom e organização.
	controls := container.NewHBox(widget.NewButtonWithIcon("", theme.ContentUndoIcon(), func() { e.check(e.Session.Undo()); e.changed() }), widget.NewButtonWithIcon("", theme.ContentRedoIcon(), func() { e.check(e.Session.Redo()); e.changed() }), widget.NewSeparator(), widget.NewButton("−", func() { e.setZoom(e.Zoom - .25) }), e.ZoomLabel, widget.NewButton("+", func() { e.setZoom(e.Zoom + .25) }), widget.NewButton("Ajustar", e.fitZoom), e.DeleteButton, layout.NewSpacer(), widget.NewButtonWithIcon("Fundo", theme.FileImageIcon(), e.configureBackground))
	e.Drawing = container.NewScroll(e.Board)
	// drawing preserva a área de desenho em uma superfície rolável.
	drawing := container.NewBorder(container.NewPadded(controls), nil, nil, nil, e.Drawing)
	// right contém o inspector com largura mínima agradável.
	right := container.NewBorder(section("PROPRIEDADES"), e.InspectorActions, nil, nil, container.NewScroll(e.Inspector))
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
	// label segue o tema automaticamente, incluindo contraste em superfícies escuras.
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
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
		// O log registra somente a classe do erro; a mensagem exibida ao usuário não é transcrita.
		LogApplication("ui_error", map[string]any{"code": "UI_ERROR"})
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
	if e.CurrentScreenPreview != nil {
		if e.page() == nil {
			e.CurrentScreenPreview.Disable()
		} else {
			e.CurrentScreenPreview.Enable()
		}
	}
	e.refreshSelectionActions()
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
