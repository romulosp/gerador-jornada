package desktop

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// workspace resolve a pasta pelo documento aberto ou pela escolha explícita do autor.
func (e *Editor) workspace() (*Workspace, error) {
	// directory nunca é gravado dentro de Project ou de MenuItem.
	directory := e.WorkspacePath
	if e.FilePath != "" {
		directory = filepath.Dir(e.FilePath)
	}
	if directory == "" {
		return nil, fmt.Errorf("Escolha uma pasta ou salve a jornada para definir o workspace.")
	}
	return OpenWorkspace(directory)
}

// SetProjectPath associa uma abertura explícita à pasta e ao histórico local.
func (e *Editor) SetProjectPath(path string) {
	if path == "" {
		return
	}
	// absolute torna o workspace independente do diretório do processo.
	absolute, err := filepath.Abs(path)
	if !e.check(err) {
		return
	}
	e.FilePath, e.WorkspacePath = absolute, filepath.Dir(absolute)
	if err := rememberProject(absolute, e.Session.Document.Project.Name); err != nil {
		e.Status.SetText("Projeto aberto; recentes indisponíveis: " + err.Error())
	}
	e.refreshWorkspace()
}

// refreshWorkspace atualiza a árvore somente em operações de arquivo, sem reler pacotes a cada pintura.
func (e *Editor) refreshWorkspace() {
	if e.WorkspaceTree == nil {
		return
	}
	e.WorkspaceEntries = nil
	if workspace, err := e.workspace(); err == nil {
		// entries mantém os diagnósticos por arquivo, sem ocultar jornadas válidas.
		entries, err := workspace.List()
		if err == nil {
			e.WorkspaceEntries = entries
		} else {
			e.Status.SetText(err.Error())
		}
	}
	e.WorkspaceTree.Refresh()
	e.WorkspaceTree.OpenBranch("workspace")
}

// workspacePanel usa uma árvore Fyne independente da estrutura de componentes da tela.
func (e *Editor) workspacePanel() fyne.CanvasObject {
	// tree lista filhos da pasta de forma virtualizada; arquivos continuam identificados por nome relativo.
	tree := widget.NewTree(
		func(id string) []string {
			if id == "" {
				return []string{"workspace"}
			}
			if id != "workspace" {
				return nil
			}
			// ids preserva a ordem estável da listagem sem depender de nomes de exibição.
			ids := []string{}
			for _, entry := range e.WorkspaceEntries {
				ids = append(ids, entry.Filename)
			}
			return ids
		},
		func(id string) bool { return id == "workspace" },
		func(branch bool) fyne.CanvasObject {
			// label reserva duas linhas nas folhas para manter o nome físico visível em painéis estreitos.
			label := widget.NewLabel("Pasta")
			if !branch {
				label.SetText("arquivo.jornada\nNome da jornada")
			}
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(id string, branch bool, object fyne.CanvasObject) {
			// label é reaproveitado pela árvore, incluindo nomes Unicode.
			label := object.(*widget.Label)
			if branch {
				if workspace, err := e.workspace(); err == nil {
					label.SetText(filepath.Base(workspace.Root))
				} else {
					label.SetText("Escolha uma pasta de projeto")
				}
				return
			}
			for _, entry := range e.WorkspaceEntries {
				if entry.Filename == id {
					text := entry.Filename + "\n" + entry.Name
					if entry.Error != "" {
						text += " · indisponível"
					}
					label.SetText(text)
					return
				}
			}
			label.SetText(id)
		})
	// Root explícita faz a pasta aparecer e permite ao Tree percorrer seus filhos.
	tree.Root = "workspace"
	e.WorkspaceTree = tree
	tree.OnSelected = func(id string) {
		if id == "workspace" {
			return
		}
		e.confirmDiscard(func() {
			// workspace e document só substituem a edição após leitura íntegra.
			workspace, err := e.workspace()
			if !e.check(err) {
				return
			}
			document, err := workspace.Load(id)
			if !e.check(err) {
				return
			}
			e.loadDocument(document, filepath.Join(workspace.Root, id))
		})
	}
	// more mantém as ações de workspace no menu de três pontos pedido na spec.
	more := widget.NewButton("⋮", func() {
		// menu não cria arquivos até a escolha e confirmação de uma ação.
		menu := fyne.NewMenu("Workspace", fyne.NewMenuItem("Criar jornada", e.createWorkspaceJourney), fyne.NewMenuItem("Importar jornada existente", e.importWorkspaceJourney), fyne.NewMenuItem("Escolher pasta", e.chooseWorkspace), fyne.NewMenuItem("Atualizar", e.refreshWorkspace))
		widget.NewPopUpMenu(menu, e.Window.Canvas()).Show()
	})
	tree.OpenBranch("workspace")
	return container.NewBorder(container.NewBorder(nil, nil, widget.NewLabel("Jornadas do projeto"), more), nil, nil, nil, tree)
}

// chooseWorkspace altera somente a pasta de trabalho da árvore, sem publicar um manifesto novo.
func (e *Editor) chooseWorkspace() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if !e.check(err) || uri == nil {
			return
		}
		if e.FilePath != "" && filepath.Clean(filepath.Dir(e.FilePath)) != filepath.Clean(uri.Path()) {
			dialog.ShowInformation("Pasta da jornada", "Abra uma jornada dessa pasta para trocar de workspace. A jornada atual permanece na pasta em que foi salva.", e.Window)
			return
		}
		e.WorkspacePath = uri.Path()
		e.refreshWorkspace()
	}, e.Window)
}

// createWorkspaceJourney confirma nome e modo em branco antes de gravar na pasta atual.
func (e *Editor) createWorkspaceJourney() {
	// workspace é resolvido antes do diálogo para localizar erros de pasta.
	workspace, err := e.workspace()
	if !e.check(err) {
		return
	}
	name := entry("Nova jornada")
	dialog.ShowCustomConfirm("Criar jornada", "Criar", "Cancelar", widget.NewForm(widget.NewFormItem("Nome", name)), func(yes bool) {
		if !yes {
			return
		}
		// filename é publicado sem substituir homônimos; a sessão corrente continua preservada.
		filename, _, err := workspace.Create(name.Text, false)
		if !e.check(err) {
			return
		}
		e.refreshWorkspace()
		e.Status.SetText("Jornada criada: " + filename)
	}, e.Window)
}

// importWorkspaceJourney valida e copia uma origem, sem pedir ao seletor que abra um arquivo para escrita.
func (e *Editor) importWorkspaceJourney() {
	// workspace preserva o destino escolhido antes da seleção de origem.
	workspace, err := e.workspace()
	if !e.check(err) {
		return
	}
	chooser := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if !e.check(err) || reader == nil {
			return
		}
		// source é fechado antes da publicação da cópia; a origem nunca recebe escrita.
		source := reader.URI().Path()
		reader.Close()
		filename, err := workspace.Import(source)
		if !e.check(err) {
			return
		}
		e.refreshWorkspace()
		e.Status.SetText("Jornada importada: " + filename)
	}, e.Window)
	chooser.SetFilter(storage.NewExtensionFileFilter([]string{".jornada"}))
	chooser.Show()
}

// defaultProjectFolder escolhe um destino local inicial editável, sem criar pasta ao abrir o diálogo.
func defaultProjectFolder() string {
	// home é fornecido pelo sistema; a pasta só nasce após confirmação do autor.
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, "JornadaStudio", "Projetos")
}
