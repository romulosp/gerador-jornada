//go:build launchersmoke

package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/desktop"
	"gerador-jornada.local/editor/internal/application"
)

// walk percorre os controles nativos de teste pelo renderer público do Fyne.
func walk(object fyne.CanvasObject, visit func(fyne.CanvasObject), seen map[fyne.CanvasObject]bool) {
	if object == nil || seen[object] {
		return
	}
	seen[object] = true
	visit(object)
	var children []fyne.CanvasObject
	if c, ok := object.(*fyne.Container); ok {
		children = c.Objects
	} else if w, ok := object.(fyne.Widget); ok {
		children = test.WidgetRenderer(w).Objects()
	}
	for _, child := range children {
		walk(child, visit, seen)
	}
}

// tap aciona o botão depois da busca, sem percorrer renderers destruídos pelo callback.
func tap(object fyne.CanvasObject, label string) {
	var button *widget.Button
	walk(object, func(o fyne.CanvasObject) {
		if b, ok := o.(*widget.Button); ok && b.Text == label && button == nil {
			button = b
		}
	}, map[fyne.CanvasObject]bool{})
	if button == nil {
		panic("botão ausente: " + label)
	}
	test.Tap(button)
}

// main executa o fluxo real no driver gráfico; go test headless não cobre falhas de HWND.
func main() {
	a := app.NewWithID("br.local.geradorjornada.launcher-smoke")
	if theme := os.Getenv("JORNADA_REPRO_THEME"); theme != "" {
		a.Preferences().SetString("editor.theme", theme)
	}
	var editor *desktop.Editor
	var launcher *desktop.Launcher
	launcher = desktop.NewLauncher(a, func(document application.Document, path string) {
		fmt.Println("CALLBACK editor.New")
		editor = desktop.New(a, document)
		fmt.Println("SetProjectPath")
		editor.SetProjectPath(path)
		fmt.Println("ShowWindow")
		editor.ShowWindow()
		fmt.Println("Close launcher")
		launcher.Window.Close()
		fmt.Println("callback terminou")
		go func() {
			time.Sleep(400 * time.Millisecond)
			fyne.Do(func() { editor.EnableWindowMode(); editor.FitToWindow(); editor.OfferRecovery() })
		}()
	})
	launcher.Window.Show()
	go func() {
		time.Sleep(time.Second)
		if os.Getenv("JORNADA_REPRO_MODE") == "recent" {
			fyne.DoAndWait(func() { tap(launcher.Window.Content(), "Abrir projeto") })
		} else {
			label := "Novo projeto em branco"
			if os.Getenv("JORNADA_REPRO_MODE") == "example" {
				label = "Novo projeto com exemplo"
			}
			fyne.DoAndWait(func() { tap(launcher.Window.Content(), label) })
			time.Sleep(200 * time.Millisecond)
			fyne.DoAndWait(func() {
				popup := launcher.Window.Canvas().Overlays().Top()
				walk(popup, func(o fyne.CanvasObject) {
					if e, ok := o.(*widget.Entry); ok {
						if e.Text == "Meu projeto" {
							e.SetText("teste")
						} else {
							e.SetText(os.Getenv("JORNADA_REPRO_DIR"))
						}
					}
				}, map[fyne.CanvasObject]bool{})
				tap(popup, "Criar")
			})
		}
		time.Sleep(2 * time.Second)
		fyne.DoAndWait(func() {
			if editor == nil {
				panic("editor não abriu")
			}
			f, err := os.Create(filepath.Join(os.Getenv("JORNADA_REPRO_DIR"), "editor-"+os.Getenv("JORNADA_REPRO_MODE")+"-"+os.Getenv("JORNADA_REPRO_THEME")+".png"))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err := png.Encode(f, editor.Window.Canvas().Capture()); err != nil {
				panic(err)
			}
			fmt.Println("PASS editor permaneceu aberto")
			a.Quit()
		})
	}()
	a.Run()
	fmt.Println("Run terminou")
}
