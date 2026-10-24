package desktop

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// configureWorkspacePreview injeta leitura e confirmação de autoria sem levar filesystem ao domínio.
func (e *Editor) configureWorkspacePreview(state *application.Preview, window fyne.Window) {
	if workspace, err := e.workspace(); err == nil {
		state.LoadJourney = workspace.Load
	}
	if e.FilePath != "" {
		state.CurrentJourney = filepath.Base(e.FilePath)
	}
	// pending evita confirmações duplicadas somente nesta janela de prévia.
	pending := false
	state.OnMissingJourney = func(componentID, itemID string) {
		if pending {
			return
		}
		// page e component localizam a folha realmente ativada, nunca pelo nome do movimento.
		page := state.Document.Screen(state.ScreenID)
		if page == nil {
			return
		}
		component := page.ComponentByID(componentID)
		if component == nil || component.Properties.Menu == nil || !component.Enabled {
			return
		}
		parent, index := menuFindItem(&component.Properties.Menu.Root, itemID)
		if parent == nil || index < 0 {
			return
		}
		item := parent.Items[index]
		if !item.Enabled || item.Submenu != nil || item.TargetScreenID != "" {
			return
		}
		if state.CurrentJourney == "" {
			dialog.ShowInformation("Salve a jornada de origem", "Salve a jornada e execute a prévia novamente para criar e vincular outra jornada na mesma pasta.", window)
			return
		}
		// origin conserva o contexto mesmo se outra janela modificar a sessão enquanto o diálogo está aberto.
		origin := state.CurrentJourney
		pending = true
		message := fmt.Sprintf("Criar uma jornada em branco chamada %q e vinculá-la a este item?\n\nA ação salva o vínculo e as alterações da jornada de origem. Arquivos antigos recebem cópia antes da atualização para catálogo 4.", item.Label)
		confirmation := dialog.NewConfirm("Criar jornada para este item?", message, func(yes bool) {
			pending = false
			if !yes {
				return
			}
			if origin != state.CurrentJourney {
				dialog.ShowError(fmt.Errorf("A jornada de origem mudou; selecione o item novamente."), window)
				return
			}
			if _, err := e.createLinkedJourney(context.Background(), state, componentID, itemID); err != nil {
				dialog.ShowError(err, window)
			}
		}, window)
		confirmation.SetConfirmText("Criar jornada")
		confirmation.SetDismissText("Cancelar")
		confirmation.Show()
	}
}

// createLinkedJourney confirma duas publicações sem perder a origem em caso de falha da gravação.
// Cancelar o diálogo nunca chama esta operação; ctx cancelado testa e preserva o estado anterior.
func (e *Editor) createLinkedJourney(ctx context.Context, state *application.Preview, componentID, itemID string) (string, error) {
	// workspace e sourcePath ficam restritos à mesma pasta e a arquivo regular já salvo.
	workspace, err := e.workspace()
	if err != nil {
		return "", err
	}
	sourcePath, err := workspace.Resolve(state.CurrentJourney)
	if err != nil {
		return "", fmt.Errorf("Jornada de origem indisponível: %w", err)
	}
	// current evita sobrescrever edições confirmadas desde a abertura da prévia.
	current := e.FilePath != "" && strings.EqualFold(filepath.Clean(e.FilePath), sourcePath)
	var staged *application.Session
	if current {
		// copy conserva o histórico, mas Change sempre escreve outro documento independente.
		copy := *e.Session
		copy.UndoStack = append([]application.Document{}, e.Session.UndoStack...)
		copy.RedoStack = append([]application.Document{}, e.Session.RedoStack...)
		staged = &copy
	} else {
		// source é relida antes da autoria numa jornada visitada que não está aberta no editor.
		source, err := projectzip.Load(sourcePath)
		if err != nil {
			return "", err
		}
		staged = application.NewSession(source)
	}
	// locate sempre usa o UUID da instância e do item no documento mais recente.
	locate := func(document *application.Document) (*catalog.MenuItem, error) {
		for pageIndex := range document.Screens {
			if component := document.Screens[pageIndex].ComponentByID(componentID); component != nil && component.Properties.Menu != nil {
				parent, index := menuFindItem(&component.Properties.Menu.Root, itemID)
				if parent != nil && index >= 0 {
					return &parent.Items[index], nil
				}
			}
		}
		return nil, fmt.Errorf("O item de origem foi removido ou alterado.")
	}
	item, err := locate(&staged.Document)
	if err != nil {
		return "", err
	}
	if !item.Enabled || item.Submenu != nil || item.TargetScreenID != "" {
		return "", fmt.Errorf("O item não é uma folha habilitada sem destino.")
	}
	if item.TargetJourney != "" {
		if _, err := workspace.Load(item.TargetJourney); !errors.Is(err, application.ErrJourneyMissing) {
			return "", fmt.Errorf("O item já possui destino; confira a jornada configurada.")
		}
	}
	// filename deriva do nome somente depois da confirmação; o vínculo gravado é explícito.
	filename, err := JourneyFilename(item.Label)
	if err != nil {
		return "", err
	}
	filename, err = workspace.uniqueName(filename)
	if err != nil {
		return "", err
	}
	name := item.Label
	if err := staged.Change(func(document *application.Document) error {
		// target é localizado novamente na cópia de Change para preservar atomicidade do histórico.
		target, err := locate(document)
		if err != nil {
			return err
		}
		target.TargetJourney = filename
		return targetValidation(document)
	}); err != nil {
		return "", err
	}
	// targetData é verificado antes de qualquer escrita da nova jornada.
	targetData, err := projectzip.Encode(NewJourneyDocument(name, false))
	if err != nil {
		return "", err
	}
	if _, err := projectzip.Decode(targetData); err != nil {
		return "", err
	}
	if err := workspace.publishNew(filename, targetData); err != nil {
		return "", err
	}
	if err := projectzip.Save(ctx, sourcePath, staged.Document); err != nil {
		// rollback remove apenas o arquivo que esta operação publicou de forma exclusiva.
		if cleanup := os.Remove(filepath.Join(workspace.Root, filename)); cleanup != nil {
			return "", fmt.Errorf("%v; não foi possível remover a jornada criada: %w", err, cleanup)
		}
		return "", err
	}
	staged.Document.Project.CatalogVersion = 4
	staged.SavedRevision = staged.Document.Project.Revision
	if current {
		e.Session = staged
	}
	// previewItem publica somente o vínculo na visita, preservando campos e percurso transitórios.
	previewItem, err := locate(&state.Document)
	if err == nil {
		previewItem.TargetJourney = filename
	}
	e.refreshWorkspace()
	if current {
		e.refresh()
	}
	e.Status.SetText("Jornada criada e vinculada: " + filename)
	return filename, nil
}

// targetValidation verifica as uniões tipadas sem exigir que o novo arquivo já exista durante o rascunho.
func targetValidation(document *application.Document) error {
	// page valida toda a composição antes de permitir confirmar o novo vínculo.
	for _, page := range document.Screens {
		if err := page.Validate(); err != nil {
			return err
		}
	}
	return nil
}
