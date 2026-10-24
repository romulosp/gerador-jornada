//go:build splashsmoke

// Command splash-smoke verifica a apresentação e a passagem launcher/editor no driver gráfico real.
package main

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/desktop"
	"gerador-jornada.local/editor/internal/application"
)

// visit percorre somente os controles da janela pertencente ao teste.
func visit(object fyne.CanvasObject, action func(fyne.CanvasObject), seen map[fyne.CanvasObject]bool) {
	if object == nil || seen[object] {
		return
	}
	seen[object] = true
	action(object)
	// children inclui renderers públicos usados pelos popups Fyne.
	var children []fyne.CanvasObject
	if c, ok := object.(*fyne.Container); ok {
		children = c.Objects
	} else if w, ok := object.(fyne.Widget); ok {
		children = test.WidgetRenderer(w).Objects()
	}
	for _, child := range children {
		visit(child, action, seen)
	}
}

// tap encontra o botão antes de acioná-lo; o callback pode substituir/destruir o conteúdo atual.
func tap(object fyne.CanvasObject, text string) {
	// selected conserva somente a instância realmente renderizada.
	var selected *widget.Button
	visit(object, func(object fyne.CanvasObject) {
		if button, ok := object.(*widget.Button); ok && button.Text == text && selected == nil {
			selected = button
		}
	}, map[fyne.CanvasObject]bool{})
	if selected == nil {
		panic("botão ausente: " + text)
	}
	test.Tap(selected)
}

// main observa a execução com vídeo completo, atraso opcional e criação/reabertura posterior.
func main() {
	// directory isola os projetos/configurações da regressão, sem tocar projetos do usuário.
	directory := os.Getenv("JORNADA_REPRO_DIR")
	if directory == "" {
		panic("JORNADA_REPRO_DIR obrigatório")
	}
	if err := os.MkdirAll(directory, 0700); err != nil {
		panic(err)
	}
	// delay é somente da ferramenta de teste; a aplicação distribuída nunca usa espera artificial.
	delayMS, _ := strconv.Atoi(os.Getenv("JORNADA_SPLASH_DELAY_MS"))
	mode := os.Getenv("JORNADA_SPLASH_MODE")
	desktop.InitApplicationLog()
	defer desktop.CloseApplicationLog()
	app := app.NewWithID("br.local.geradorjornada.splash-smoke")
	var launcher *desktop.Launcher
	var editor *desktop.Editor
	// began e loadedAt permitem demonstrar que nenhuma condição isolada liberou o fluxo.
	began := time.Now()
	loadedAt := time.Time{}
	var splash *desktop.Splash
	splash = desktop.NewSplash(app, func(window fyne.Window) fyne.Window {
		if mode == "cancel" || mode == "fail" {
			panic("cancelamento/falha abriu a tela seguinte")
		}
		if time.Since(began) < 10*time.Second || loadedAt.IsZero() {
			panic("transição prematura")
		}
		launcher = desktop.NewLauncherInWindow(app, window, func(document application.Document, path string) {
			editor = desktop.New(app, document)
			editor.SetProjectPath(path)
			editor.ShowWindow()
			launcher.Window.Close()
		})
		return launcher.Window
	}, func() {
		if mode == "cancel" || mode == "fail" {
			panic("cancelamento/falha foi tratado como sucesso")
		}
		fmt.Printf("READY mode=%s elapsed_ms=%d loaded_ms=%d\n", mode, time.Since(began).Milliseconds(), loadedAt.Sub(began).Milliseconds())
		if launcher.Window.FullScreen() {
			panic("launcher não restaurou janela normal")
		}
		if mode == "create" {
			tap(launcher.Window.Content(), "Novo projeto em branco")
			popup := launcher.Window.Canvas().Overlays().Top()
			visit(popup, func(object fyne.CanvasObject) {
				if entry, ok := object.(*widget.Entry); ok {
					if entry.Text == "Meu projeto" {
						entry.SetText("Projeto splash")
					} else {
						entry.SetText(directory)
					}
				}
			}, map[fyne.CanvasObject]bool{})
			tap(popup, "Criar")
		} else if mode == "recent" {
			tap(launcher.Window.Content(), "Abrir projeto")
		}
		go func() {
			time.Sleep(500 * time.Millisecond)
			fyne.Do(func() {
				// target e filename identificam somente a evidência da janela já revelada.
				target := launcher.Window
				if mode == "create" || mode == "recent" {
					if editor == nil || editor.Session.Document.Project.Name != "Projeto splash" {
						panic("editor não abriu depois da splash")
					}
					editor.EnableWindowMode()
					editor.FitToWindow()
					target = editor.Window
				}
				filename := filepath.Join(directory, "splash-"+mode+".png")
				file, err := os.Create(filename)
				if err != nil {
					panic(err)
				}
				if err := png.Encode(file, target.Canvas().Capture()); err != nil {
					panic(err)
				}
				file.Close()
				fmt.Println("PASS", mode)
				app.Quit()
			})
		}()
	})
	splash.Show(func(ctx context.Context, report func(float64, string)) error {
		if mode == "fail" {
			return fmt.Errorf("Falha de recurso injetada pela ferramenta de teste")
		}
		if err := desktop.LoadStartupResources(ctx, report); err != nil {
			return err
		}
		// remaining injeta carga lenta depois de preparar os recursos, sem atrasar a thread da interface.
		remaining := time.Duration(delayMS)*time.Millisecond - time.Since(began)
		if remaining > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(remaining):
			}
		}
		loadedAt = time.Now()
		return nil
	})
	if mode == "cancel" || mode == "fail" {
		// O fechamento real deve parar áudio e workers sem executar next/finished.
		go func() { time.Sleep(2 * time.Second); fyne.Do(func() { splash.Window.Close() }) }()
	}
	app.Run()
	if mode == "cancel" || mode == "fail" {
		fmt.Println("PASS", mode, "sem liberar launcher")
	}
}
