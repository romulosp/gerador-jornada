package desktop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/buildinfo"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
)

// StudioCredit conserva a autoria literal solicitada para o rodapé da seleção de projetos.
const StudioCredit = "Desenvolvido por Rômulo Penha - Canais Parceiro - 03"

// Launcher é a porta de entrada sem documento implícito, com recentes locais e criação confirmada.
type Launcher struct {
	// App e Window recebem somente a interface de boas-vindas, sem abrir sessão de edição.
	App    fyne.App
	Window fyne.Window
	// Recents e Message apresentam o histórico e diagnósticos sem impedir criar outro projeto.
	Recents *fyne.Container
	Message *widget.Label
	// ShowSplash oferece a preferência persistente de repetir a apresentação.
	ShowSplash *widget.Check
	// Open recebe um documento já validado e seu caminho explicitamente escolhido.
	Open func(application.Document, string)
}

// NewLauncher monta a tela inicial com o tema do editor e não cria nenhum projeto ao iniciar.
func NewLauncher(app fyne.App, open func(application.Document, string)) *Launcher {
	// A preferência global é aplicada antes de NewWindow, também quando o driver já está rodando.
	prepareApplicationTheme(app, app.Preferences().String(themePreferenceKey))
	return NewLauncherInWindow(app, app.NewWindow("Jornada Studio · Boas-vindas"), open)
}

// NewLauncherInWindow reutiliza a janela da splash, preservando o fade sem destruir o driver.
func NewLauncherInWindow(app fyne.App, window fyne.Window, open func(application.Document, string)) *Launcher {
	prepareApplicationTheme(app, app.Preferences().String(themePreferenceKey))
	// launcher não é janela master: ao abrir o editor, ela pode fechar sem encerrar o aplicativo.
	launcher := &Launcher{App: app, Window: window, Open: open, Recents: container.NewVBox(), Message: widget.NewLabel("")}
	window.SetTitle("Jornada Studio · Boas-vindas")
	window.SetPadded(true)
	launcher.Message.Wrapping = fyne.TextWrapWord
	// logo reutiliza a identidade já incorporada, sem depender de assets externos.
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("jornada-studio.png", builtinassets.IconPNG()))
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(96, 96))
	launcher.ShowSplash = launcher.splashPreferenceCheck()
	// welcome organiza o começo do trabalho, sem abrir o exemplo automaticamente.
	welcome := container.NewVBox(logo, widget.NewLabelWithStyle("Jornada Studio", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), widget.NewLabelWithStyle("Desenhe experiências. Conecte jornadas.", fyne.TextAlignCenter, fyne.TextStyle{}), widget.NewSeparator(), widget.NewButtonWithIcon("Novo projeto em branco", theme.DocumentCreateIcon(), func() { launcher.newProject(false) }), widget.NewButtonWithIcon("Novo projeto com exemplo", theme.ContentAddIcon(), func() { launcher.newProject(true) }), widget.NewButtonWithIcon("Abrir projeto existente", theme.FolderOpenIcon(), launcher.openProject), launcher.ShowSplash, launcher.Message)
	// recentPanel separa a lista persistida das ações de criação, como uma tela de boas-vindas desktop.
	recentPanel := container.NewBorder(section("PROJETOS RECENTES"), nil, nil, nil, container.NewVScroll(launcher.Recents))
	split := container.NewHSplit(recentPanel, container.NewPadded(container.NewCenter(welcome)))
	split.Offset = .43
	launcher.Window.SetContent(container.NewBorder(container.NewPadded(widget.NewLabel("JORNADA STUDIO · "+buildinfo.Version())), container.NewPadded(widget.NewLabelWithStyle(StudioCredit, fyne.TextAlignCenter, fyne.TextStyle{})), nil, nil, split))
	launcher.Window.Resize(fyne.NewSize(980, 640))
	launcher.Window.CenterOnScreen()
	launcher.refreshRecents()
	return launcher
}

// refreshRecents conserva entradas indisponíveis e oferece os demais projetos normalmente.
func (l *Launcher) refreshRecents() {
	l.Recents.RemoveAll()
	// entries vêm apenas da preferência local, sem descoberta de rede ou telemetria.
	entries, err := loadRecents()
	if err != nil {
		l.Message.SetText("Não foi possível carregar os recentes: " + err.Error())
	}
	if len(entries) == 0 {
		l.Recents.Add(widget.NewLabel("Nenhum projeto recente.\nCrie ou abra um projeto para começar."))
	}
	for _, recent := range entries {
		// name e path quebram linhas separadamente para caminhos longos não ocultarem a identidade.
		name := widget.NewLabelWithStyle(recent.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		name.Wrapping = fyne.TextWrapWord
		path := widget.NewLabel(recent.Path)
		path.Wrapping = fyne.TextWrapBreak
		description := "Última abertura: " + recent.OpenedAt.Local().Format("02/01/2006 15:04")
		available := recent.Available()
		if !available {
			description = "Indisponível · arquivo movido ou removido"
		}
		status := widget.NewLabel(description)
		status.Wrapping = fyne.TextWrapWord
		// button mantém cada entrada isolada; um arquivo removido não bloqueia os demais.
		button := widget.NewButtonWithIcon("Abrir projeto", theme.FolderOpenIcon(), func() { l.openRecent(recent) })
		button.Alignment = widget.ButtonAlignLeading
		if !available {
			button.Disable()
		}
		l.Recents.Add(container.NewPadded(container.NewVBox(name, path, status, button, widget.NewSeparator())))
	}
}

// openRecent valida novamente uma entrada, inclusive se ela mudou após a pintura da lista.
func (l *Launcher) openRecent(recent RecentProject) {
	// path pode apontar à pasta de workspace de uma sessão anterior.
	path := recent.Path
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		// workspace resolve uma jornada legível da pasta sem criar um manifesto novo.
		workspace, err := OpenWorkspace(path)
		if err != nil {
			l.Message.SetText(err.Error())
			return
		}
		entries, err := workspace.List()
		if err != nil {
			l.Message.SetText(err.Error())
			return
		}
		path = ""
		for _, entry := range entries {
			if entry.Error == "" {
				path = filepath.Join(workspace.Root, entry.Filename)
				break
			}
		}
		if path == "" {
			l.Message.SetText("A pasta não contém uma jornada válida. Crie um projeto nessa pasta.")
			return
		}
	}
	// document só é entregue ao editor após o codec verificar o pacote.
	document, err := projectzip.Load(path)
	if err != nil {
		l.Message.SetText("Projeto indisponível: " + err.Error())
		return
	}
	l.finishOpen(document, path)
}

// finishOpen registra a escolha e transfere a sessão somente quando há documento válido.
func (l *Launcher) finishOpen(document application.Document, path string) {
	if err := rememberProject(path, document.Project.Name); err != nil {
		l.Message.SetText("Recentes indisponíveis: " + err.Error())
	}
	if l.Open != nil {
		l.Open(document, path)
	}
}

// createProject grava a jornada inicial em sua pasta, sem sobrescrever arquivos existentes.
func (l *Launcher) createProject(name, base string, example bool) error {
	// filename normaliza apenas a pasta e o arquivo, preservando o nome exibido.
	filename, err := JourneyFilename(name)
	if err != nil {
		return err
	}
	if strings.TrimSpace(base) == "" {
		return fmt.Errorf("Escolha a pasta do projeto.")
	}
	root := filepath.Join(base, strings.TrimSuffix(filename, ".jornada"))
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	workspace, err := OpenWorkspace(root)
	if err != nil {
		return err
	}
	file, document, err := workspace.Create(name, example)
	if err != nil {
		return err
	}
	l.finishOpen(document, filepath.Join(workspace.Root, file))
	return nil
}

// newProject confirma nome e pasta antes de criar o documento em branco ou com exemplo.
func (l *Launcher) newProject(example bool) {
	// name e folder são rascunhos; cancelar o diálogo não cria pasta, arquivo ou recente.
	name := entry("Meu projeto")
	folder := entry(defaultProjectFolder())
	browse := widget.NewButton("Escolher pasta…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				l.Message.SetText(err.Error())
			} else if uri != nil {
				folder.SetText(uri.Path())
			}
		}, l.Window)
	})
	title := "Novo projeto em branco"
	if example {
		title = "Novo projeto com exemplo"
	}
	dialog.ShowCustomConfirm(title, "Criar", "Cancelar", container.NewVBox(widget.NewForm(widget.NewFormItem("Nome", name), widget.NewFormItem("Pasta base", folder)), browse, widget.NewLabel("Será criada uma pasta com o nome do projeto.")), func(yes bool) {
		if !yes {
			return
		}
		if err := l.createProject(name.Text, folder.Text, example); err != nil {
			l.Message.SetText(err.Error())
			dialog.ShowError(err, l.Window)
		}
	}, l.Window)
}

// openProject recebe somente leitura do arquivo escolhido e conserva a tela inicial em caso de erro.
func (l *Launcher) openProject() {
	// chooser limita a apresentação aos arquivos de jornada, sem escrever durante a seleção.
	chooser := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			l.Message.SetText(err.Error())
			return
		}
		if reader == nil {
			return
		}
		path := reader.URI().Path()
		reader.Close()
		l.openRecent(RecentProject{Path: path})
	}, l.Window)
	chooser.SetFilter(storage.NewExtensionFileFilter([]string{".jornada"}))
	chooser.Show()
}
