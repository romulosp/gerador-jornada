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
	"gerador-jornada.local/editor/internal/adapter/pptx"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
)

// main compõe adaptadores; regras de negócio permanecem nos pacotes internos.
func main() {
	// examplePPT permite validar o exportador sem interação com a janela.
	examplePPT := flag.String("export-example", "", "Gera um PPTX de exemplo no caminho indicado.")
	// exampleProject permite abrir o conjunto de referência no editor.
	exampleProject := flag.String("save-example", "", "Salva um projeto .jornada de exemplo.")
	// smokeScreenshot captura uma janela real e encerra a prova de inicialização.
	smokeScreenshot := flag.String("smoke-screenshot", "", "Captura a janela inicial em PNG e encerra.")
	flag.Parse()
	// document contém um exemplo totalmente editável e sem dependências externas.
	document := application.Example()
	builtinassets.DefaultBackgrounds(&document)
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
	// editor recebe as dependências concretas sem um contêiner de DI.
	editor := desktop.New(desktopApp, document)
	if *smokeScreenshot != "" {
		go func() {
			time.Sleep(2 * time.Second)
			fyne.DoAndWait(func() { editor.EnableWindowMode(); editor.FitToWindow() })
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
				if err = png.Encode(file, editor.Window.Canvas().Capture()); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				file.Close()
				desktopApp.Quit()
			})
		}()
	} else {
		go func() {
			time.Sleep(400 * time.Millisecond)
			fyne.Do(func() { editor.EnableWindowMode(); editor.FitToWindow(); editor.OfferRecovery() })
		}()
	}
	// A abertura aplica tela cheia somente depois de criar e exibir a janela nativa.
	editor.ShowWindow()
	desktopApp.Run()
}

// fail preserva o código de saída dos utilitários de validação.
func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
