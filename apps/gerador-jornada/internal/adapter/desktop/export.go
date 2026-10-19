package desktop

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/javaefl"
	"gerador-jornada.local/editor/internal/adapter/pptx"
	"gerador-jornada.local/editor/internal/application"
)

// exportPPT permite escolher exploração e demonstrações sem descartar erros silenciosamente.
func (e *Editor) exportPPT() {
	if e.ExportCancel != nil {
		e.check(fmt.Errorf("há uma exportação em andamento"))
		return
	}
	// exploration mantém o percurso manual como padrão.
	exploration := widget.NewCheck("Incluir exploração navegável", nil)
	exploration.SetChecked(true)
	// demonstration mantém roteiros explícitos no arquivo.
	demonstration := widget.NewCheck("Incluir todos os roteiros não vazios", nil)
	demonstration.SetChecked(true)
	// transition permite um acabamento de apresentação sem alterar o roteiro de negócio.
	transition := newSelect([]string{"Fade suave", "Deslizar", "Revelar", "Sem transição"}, nil)
	transition.SetSelected("Fade suave")
	dialog.ShowCustomConfirm("Exportar PowerPoint", "Escolher destino", "Cancelar", container.NewVBox(exploration, demonstration,
		widget.NewForm(widget.NewFormItem("Transição (300 ms)", transition)),
		widget.NewLabel("Textos e controles permanecem editáveis.\nA digitação é simulada, sem entrada real no PPT.")), func(yes bool) {
		if !yes {
			return
		}
		// options diferencia seleção vazia de todos os roteiros.
		options := pptx.Options{Exploration: exploration.Checked, Transition: map[string]string{"Fade suave": "fade", "Deslizar": "push", "Revelar": "wipe", "Sem transição": "none"}[transition.Selected]}
		if !demonstration.Checked {
			options.ScenarioIDs = []string{}
		}
		e.chooseOutput("Salvar apresentação", ".pptx", func(path string) {
			e.startExport(path, func(ctx context.Context, document application.Document, progress func(float64)) error {
				return pptx.Export(ctx, path, document, options, progress)
			})
		})
	}, e.Window)
}

// startExport captura uma revisão e mantém um único job cancelável.
func (e *Editor) startExport(path string, export func(context.Context, application.Document, func(float64)) error) {
	if e.ExportCancel != nil {
		e.check(fmt.Errorf("aguarde ou cancele a exportação atual"))
		return
	}
	// snapshot não compartilha estado com a UI.
	snapshot := e.Session.Document.Copy()
	// ctx permite interromper a geração antes da publicação.
	// cancel permite interromper a operação e liberar os recursos do contexto.
	ctx, cancel := context.WithCancel(context.Background())
	e.ExportCancel = cancel
	// progress informa geração sem prometer aprovação do artefato.
	progress := widget.NewProgressBar()
	// label identifica a revisão que será exportada.
	label := widget.NewLabel(fmt.Sprintf("Exportando revisão %d…", snapshot.Project.Revision))
	// popup oferece cancelamento sem bloquear o processamento da janela.
	popup := dialog.NewCustom("Gerando arquivo", "Cancelar", container.NewVBox(label, progress), e.Window)
	popup.SetOnClosed(cancel)
	popup.Show()
	go func() {
		// err conserva o resultado real da geração e publicação.
		err := export(ctx, snapshot, func(value float64) { fyne.Do(func() { progress.SetValue(value) }) })
		fyne.Do(func() {
			e.ExportCancel = nil
			popup.Hide()
			cancel()
			if errors.Is(err, context.Canceled) {
				e.Status.SetText("Exportação cancelada · destino anterior preservado")
				return
			}
			if !e.check(err) {
				return
			}
			e.Status.SetText("Arquivo gerado · " + filepath.Base(path))
			dialog.ShowInformation("Arquivo gerado", "Salvo em:\n"+path+"\n\nA homologação no PowerPoint e sua aprovação manual continuam pendentes.", e.Window)
		})
	}()
}

// exportJava conserva o gate de homologação sem inventar assinaturas EFL.
func (e *Editor) exportJava() {
	if javaExportPaused {
		return
	}
	// profile só se torna verificado depois de evidência real da biblioteca.
	profile := javaefl.CurrentProfile()
	if !profile.Verified {
		e.check(fmt.Errorf("JAVA_PROFILE_UNVERIFIED\nA exportação Java/EFL precisa dos JARs EFL e de um exemplo funcional de integração Java 1.6.\nO editor e a geração de PPT não dependem dessa biblioteca."))
		return
	}
	// packageName é validado antes de publicar fontes.
	packageName := entry(e.Session.Document.JavaOptions.PackageName)
	dialog.ShowCustomConfirm("Gerar fontes Java/EFL", "Escolher pasta", "Cancelar", widget.NewForm(widget.NewFormItem("Pacote Java", packageName)), func(yes bool) {
		if !yes || !e.check(javaefl.ValidatePackageName(packageName.Text)) {
			return
		}
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if !e.check(err) || uri == nil {
				return
			}
			// destination é sempre um novo subdiretório, sem mesclar fontes editados.
			destination := filepath.Join(uri.Path(), "jornada-java")
			e.startExport(destination, func(ctx context.Context, document application.Document, progress func(float64)) error {
				document.JavaOptions.PackageName = packageName.Text
				document.JavaOptions.ProfileID = profile.ID
				return javaefl.Export(ctx, destination, document, profile)
			})
		}, e.Window)
	}, e.Window)
}

// javaExportPaused preserva a suspensão temporária solicitada, sem remover o gerador.
const javaExportPaused = true
