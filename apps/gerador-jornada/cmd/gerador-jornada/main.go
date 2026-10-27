// Command gerador-jornada inicia o editor desktop e utilitários locais de validação.
package main

import (
	"context"
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/adapter/desktop"
	mcplocal "gerador-jornada.local/editor/internal/adapter/mcp"
	"gerador-jornada.local/editor/internal/adapter/pptx"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
)

// main compõe adaptadores; regras de negócio permanecem nos pacotes internos.
func main() {
	// mcpMode escolhe o transporte local antes de logger, splash, preferências ou GUI.
	mcpMode := flag.Bool("mcp", false, "Inicia a interface MCP local por stdin/stdout, sem janela.")
	// mcpWorkspace é a única raiz inicial autorizada explicitamente pelo operador.
	mcpWorkspace := flag.String("workspace", "", "Diretório autorizado para a sessão MCP local.")
	// legacyRoots/importRoots são concessões explícitas por finalidade, fixadas antes do transporte.
	var legacyRoots, importRoots []string
	flag.Func("legacy-root", "Pasta de evidências do legado autorizada para leitura MCP; pode repetir.", func(root string) error { legacyRoots = append(legacyRoots, root); return nil })
	flag.Func("import-root", "Pasta de entrada da biblioteca autorizada para leitura MCP; pode repetir.", func(root string) error { importRoots = append(importRoots, root); return nil })
	// examplePPT permite validar o exportador sem interação com a janela.
	examplePPT := flag.String("export-example", "", "Gera um PPTX de exemplo no caminho indicado.")
	// exampleProject permite abrir o conjunto de referência no editor.
	exampleProject := flag.String("save-example", "", "Salva um projeto .jornada de exemplo.")
	// compoundsProject gera a revisão de teste dos compostos 031, 032 e 033.
	compoundsProject := flag.String("save-compounds", "", "Salva um projeto .jornada com os compostos de teste.")
	// compoundsPPT gera a apresentação correspondente aos compostos de teste.
	compoundsPPT := flag.String("export-compounds", "", "Gera um PPTX de exemplo com os compostos.")
	// openProject carrega um projeto .jornada existente para a revisão manual.
	openProject := flag.String("open", "", "Abre um projeto .jornada ao iniciar.")
	// openScreen escolhe a tela inicial da revisão; -1 preserva a tela do projeto.
	openScreen := flag.Int("screen", -1, "Índice da tela inicial ao abrir um projeto.")
	// smokeScreenshot captura uma janela real e encerra a prova de inicialização.
	smokeScreenshot := flag.String("smoke-screenshot", "", "Captura a janela inicial em PNG e encerra.")
	flag.Parse()
	if *mcpMode {
		// incompatible recusa inclusive flags GUI explicitamente enviadas com valor vazio.
		incompatible := flag.NArg() != 0
		// option percorre as flags fornecidas para impedir mistura de modalidades.
		flag.Visit(func(option *flag.Flag) {
			if option.Name != "mcp" && option.Name != "workspace" && option.Name != "legacy-root" && option.Name != "import-root" {
				incompatible = true
			}
		})
		if incompatible {
			fail(fmt.Errorf("--mcp não aceita flags GUI ou utilitários de exportação"))
		}
		fail(mcplocal.Run(context.Background(), os.Stdin, os.Stdout, os.Stderr, mcplocal.Config{Workspace: *mcpWorkspace, LegacyRoots: legacyRoots, ImportRoots: importRoots}))
		return
	}
	if *mcpWorkspace != "" || len(legacyRoots) != 0 || len(importRoots) != 0 {
		fail(fmt.Errorf("--workspace, --legacy-root e --import-root exigem --mcp"))
	}
	// utility executa somente operações explicitamente solicitadas sem janela nem reprodução.
	utility := *exampleProject != "" || *examplePPT != "" || *compoundsProject != "" || *compoundsPPT != ""
	// O log de aplicação nasce antes do Fyne para registrar o startup real do processo.
	desktop.InitApplicationLog()
	defer desktop.CloseApplicationLog()
	// document permanece vazio no launcher; o exemplo só nasce por solicitação explícita.
	var document application.Document
	if *exampleProject != "" || *examplePPT != "" {
		document = application.Example()
		builtinassets.DefaultBackgrounds(&document)
	}
	if *openProject != "" && utility {
		// loaded valida e restaura o projeto antes de abrir a janela.
		loaded, err := projectzip.Load(*openProject)
		fail(err)
		document = loaded
		if *openScreen >= 0 && *openScreen < len(document.Screens) {
			// A tela escolhida serve somente à revisão; o documento carregado não é regravado.
			document.Journey.StartScreenID = document.Screens[*openScreen].ID
		}
	}
	if *compoundsProject != "" {
		// compounds preserva os defaults e a árvore inicial sem valores de runtime.
		compounds := application.CompoundsExample()
		builtinassets.DefaultBackgrounds(&compounds)
		fail(projectzip.Save(context.Background(), *compoundsProject, compounds))
		return
	}
	if *compoundsPPT != "" {
		// compounds contém as duas telas de teste e nenhum roteiro obrigatório.
		compounds := application.CompoundsExample()
		builtinassets.DefaultBackgrounds(&compounds)
		fail(pptx.Export(context.Background(), *compoundsPPT, compounds, pptx.Options{Exploration: true}, nil))
		return
	}
	if *exampleProject != "" {
		fail(projectzip.Save(context.Background(), *exampleProject, document))
		return
	}
	if *examplePPT != "" {
		fail(pptx.Export(context.Background(), *examplePPT, document, pptx.Options{Exploration: true}, nil))
		return
	}
	// desktopApp usa recursos de tema embutidos pelo Fyne.
	desktopApp := app.NewWithID("br.local.geradorjornada")
	desktopApp.SetIcon(fyne.NewStaticResource("jornada-studio.png", builtinassets.IconPNG()))
	// editor só nasce depois de abrir ou criar um documento confirmado pelo autor.
	var editor *desktop.Editor
	var window fyne.Window
	// prepare ajusta a janela após o primeiro layout, sem criar exemplo no launcher.
	prepare := func(current *desktop.Editor) {
		go func() {
			time.Sleep(400 * time.Millisecond)
			fyne.Do(func() { current.EnableWindowMode(); current.FitToWindow(); current.OfferRecovery() })
		}()
	}
	// ready é fechado após recursos prontos e, quando exibida, apresentação completa.
	ready := make(chan struct{})
	window = desktop.StartDesktop(desktopApp, func(current fyne.Window) fyne.Window {
		if *openProject != "" {
			editor = desktop.New(desktopApp, document)
			editor.EnableGitHub()
			editor.SetProjectPath(*openProject)
			editor.ShowWindow()
			window = editor.Window
			current.Close()
			prepare(editor)
			return window
		}
		// launcher reutiliza a janela da splash e mantém a transição até revelar a seleção de projetos.
		var launcher *desktop.Launcher
		launcher = desktop.NewLauncherInWindow(desktopApp, current, func(opened application.Document, path string) {
			desktop.LogApplication("editor_open_start", map[string]any{"source": "launcher"})
			editor = desktop.New(desktopApp, opened)
			editor.EnableGitHub()
			editor.SetProjectPath(path)
			editor.ShowWindow()
			launcher.Window.Close()
			desktop.LogApplication("editor_open_ready", map[string]any{"source": "launcher"})
			prepare(editor)
		})
		window = launcher.Window
		return window
	}, func(ctx context.Context, report func(float64, string)) error {
		if err := desktop.LoadStartupResources(ctx, report); err != nil {
			return err
		}
		if *openProject != "" {
			report(.98, "Abrindo o projeto solicitado…")
			loaded, err := projectzip.Load(*openProject)
			if err != nil {
				return err
			}
			document = loaded
			if *openScreen >= 0 && *openScreen < len(document.Screens) {
				document.Journey.StartScreenID = document.Screens[*openScreen].ID
			}
		}
		report(1, "Carregamento concluído.")
		return nil
	}, func() { close(ready) })
	if *smokeScreenshot != "" {
		go func() {
			<-ready
			time.Sleep(400 * time.Millisecond)
			fyne.DoAndWait(func() {
				if editor != nil {
					editor.EnableWindowMode()
					editor.FitToWindow()
				}
			})
			// O intervalo permite pintar o quadro ajustado antes da captura do framebuffer.
			time.Sleep(200 * time.Millisecond)
			fyne.Do(func() {
				// destination identifica a evidência da janela realmente inicializada.
				destination := *smokeScreenshot
				// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
				if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
					fmt.Fprintln(os.Stderr, err)
					desktopApp.Quit()
					return
				}
				// file recebe a imagem capturada pelo canvas da janela real.
				// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
				file, err := os.Create(destination)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					desktopApp.Quit()
					return
				}
				if err = png.Encode(file, window.Canvas().Capture()); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				file.Close()
				desktopApp.Quit()
			})
		}()
	}
	desktopApp.Run()
}

// fail preserva o código de saída dos utilitários de validação.
func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
