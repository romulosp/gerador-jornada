package desktop

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/credential"
	githubapi "gerador-jornada.local/editor/internal/adapter/github"
	"gerador-jornada.local/editor/internal/adapter/gitworkspace"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/adapter/userconfig"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// EnableGitHub compõe adaptadores opcionais após construir a janela, sem iniciar rede.
// Falha de configuração é apresentada no painel, preservando os fluxos locais do editor.
func (e *Editor) EnableGitHub() {
	// store usa o mesmo namespace da instalação e da lista de recentes.
	store, err := userconfig.New()
	if err != nil {
		e.setGitStatus("GitHub indisponível: configuração local inacessível.")
		return
	}
	// client e vault não expõem tokens aos controles.
	client := githubapi.New(githubapi.ClientID, githubapi.ClientSecret)
	e.GitDiagnostics = syncdiag.New(LogApplication)
	client.Diagnostics = e.GitDiagnostics
	e.GitHub = application.NewGitHubConnection(client, client, credential.Store{}, store, userconfig.Key(store.Root))
	e.GitHub.Observe = func(ctx context.Context, event string, err error) {
		e.GitDiagnostics.Event(ctx, event, syncdiag.ErrorFields(err))
	}
	e.GitWorkspace = gitworkspace.New(store)
	e.GitWorkspace.Diagnostics = e.GitDiagnostics
	e.GitSync = application.NewWorkspaceSync(e.GitHub, e.GitWorkspace, e.GitWorkspace)
	e.GitDiagnostics.Event(application.WithSyncOperation(context.Background()), "configuration", map[string]any{"oauth_configured": githubapi.Configured(), "callback_host": "127.0.0.1", "callback_path": "/oauth/github/callback"})
	e.refreshGitLocal()
}

// githubPanel fica separado da árvore e não modifica componentes/skins da jornada.
func (e *Editor) githubPanel() fyne.CanvasObject {
	e.GitLabel = widget.NewLabel("Projeto salvo apenas neste computador.")
	e.GitLabel.Truncation = fyne.TextTruncateEllipsis
	e.GitConnectButton = widget.NewButtonWithIcon("Conectar GitHub", githubIcon(), e.connectGitHub)
	e.GitRepositoriesButton = widget.NewButton("Selecionar repositório", e.selectGitHubRepository)
	e.GitPublishButton = widget.NewButtonWithIcon("Publicar no GitHub", githubIcon(), e.publishGitHub)
	e.GitRefreshButton = widget.NewButton("Atualizar estado", e.refreshGitHubRemote)
	e.GitDisconnectButton = widget.NewButton("Desconectar GitHub", e.disconnectGitHub)
	e.GitOpenButton = widget.NewButton("Ver repositório no navegador", e.openGitRepository)
	e.GitOpenButton.Disable()
	e.GitConnectButton.Show()
	e.GitConnectButton.SetText("Conectar GitHub")
	e.GitConnectButton.OnTapped = e.connectGitHub
	e.GitPublishButton.Hide()
	e.GitConnectButton.Disable()
	e.GitRepositoriesButton.Disable()
	e.GitPublishButton.Disable()
	e.GitRefreshButton.Disable()
	e.GitDisconnectButton.Disable()
	// more conserva ações administrativas em menu compacto para não ocultar a árvore.
	more := widget.NewButton("⋮", e.showGitMenu)
	e.GitPublishButton.Hide()
	return container.NewVBox(widget.NewSeparator(), container.NewBorder(nil, nil, nil, more, e.GitLabel), e.GitConnectButton, e.GitPublishButton)
}

// gitCanWrite bloqueia mutações explícitas durante operação Git ou recuperação diagnosticada.
func (e *Editor) gitCanWrite() bool {
	if e.GitBusy {
		return false
	}
	if e.GitRecoveryError != nil {
		dialog.ShowInformation("Recuperação Git necessária", e.GitRecoveryError.Error(), e.Window)
		return false
	}
	return true
}

// refreshGitLocal consulta metadados locais e tenta recuperar o journal ao entrar em outra raiz.
func (e *Editor) refreshGitLocal() {
	if e.GitLabel == nil || e.GitHub == nil || e.GitBusy {
		return
	}
	e.GitConnectButton.Show()
	e.GitConnectButton.SetText("Conectar GitHub")
	e.GitConnectButton.OnTapped = e.connectGitHub
	e.GitPublishButton.Hide()
	e.GitConnectButton.Disable()
	e.GitRepositoriesButton.Disable()
	e.GitPublishButton.Disable()
	e.GitRefreshButton.Disable()
	e.GitDisconnectButton.Disable()
	e.GitOpenButton.Disable()
	// workspace fixa a raiz administrativa sem colocá-la no documento.
	workspace, err := e.workspace()
	if err != nil {
		e.setGitStatus("Abra ou salve uma jornada para definir a pasta do projeto.")
		return
	}
	if e.GitRecoveredRoot != workspace.Root {
		e.GitVerifiedAt = time.Time{}
		// previous representa os bytes carregados antes de recuperar uma operação interrompida.
		previous, _ := os.ReadFile(e.FilePath)
		e.GitRecoveredRoot = workspace.Root
		e.GitRecoveryError = e.GitWorkspace.Recover(context.Background(), workspace.Root)
		if e.GitRecoveryError != nil {
			e.setGitStatus("Recuperação necessária: " + e.GitRecoveryError.Error())
			e.GitRefreshButton.Enable()
			return
		}
		// current evita manter uma sessão limpa baseada em bytes anteriores ao rollback.
		current, readErr := os.ReadFile(e.FilePath)
		if readErr == nil && string(previous) != string(current) && !e.Session.Dirty() {
			e.reloadGitActive(workspace.Root, e.FilePath)
		}
	}
	if e.GitRecoveryError != nil {
		e.setGitStatus("Recuperação necessária: " + e.GitRecoveryError.Error())
		e.GitRefreshButton.Enable()
		return
	}
	if githubapi.Configured() {
		e.GitConnectButton.Enable()
	}
	// account é uma preferência local, válida mesmo offline.
	account, accountErr := e.GitHub.Account()
	// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
	binding, bindingErr := e.GitWorkspace.Binding(workspace.Root)
	if accountErr != nil && !os.IsNotExist(accountErr) {
		e.setGitStatus("Configuração GitHub inválida; arquivo preservado.")
		return
	}
	if account.ID == 0 {
		// text reúne somente estado administrativo e nomes, sem transportar tokens ou conteúdo dos pacotes.
		text := "Projeto salvo apenas neste computador."
		if bindingErr == nil {
			text = "GitHub desconectado · " + binding.Repository.FullName() + "\nHistórico e vínculo preservados."
		}
		if !githubapi.Configured() {
			text += "\nEste build ainda precisa do cadastro OAuth do Jornada Studio."
		}
		e.setGitStatus(text)
		return
	}
	e.GitDisconnectButton.Enable()
	e.GitRepositoriesButton.Enable()
	e.GitConnectButton.SetText("Selecionar repositório")
	e.GitConnectButton.OnTapped = e.selectGitHubRepository
	if bindingErr != nil {
		if !os.IsNotExist(bindingErr) {
			e.setGitStatus("Vínculo local inválido; arquivo preservado.")
			return
		}
		e.setGitStatus("GitHub: " + account.Login + "\nSelecione o repositório deste workspace.")
		return
	}
	e.GitConnectButton.Hide()
	e.GitPublishButton.Show()
	e.GitPublishButton.Enable()
	e.GitRefreshButton.Enable()
	e.GitOpenButton.Enable()
	// text reúne somente estado administrativo e nomes, sem transportar tokens ou conteúdo dos pacotes.
	text := "GitHub: " + account.Login + "\n" + binding.Repository.FullName() + " · " + binding.Repository.Branch
	if binding.Pending != "" {
		text += "\nCommit salvo; envio pendente de verificação."
	} else if !binding.PublishedAt.IsZero() {
		text += "\nÚltima publicação: " + binding.PublishedAt.Local().Format("02/01 15:04")
	} else {
		text += "\nNenhuma publicação confirmada."
	}
	if e.GitVerifiedAt.IsZero() {
		text += "\nEstado remoto ainda não consultado nesta sessão."
	} else {
		text += "\nÚltima consulta remota: " + e.GitVerifiedAt.Local().Format("15:04:05")
	}
	// changes informa pendências locais sem precisar consultar o GitHub.
	changes, changesErr := e.GitWorkspace.LocalChanges(workspace.Root)
	if changesErr == nil {
		text += fmt.Sprintf("\nAlterações locais: %d arquivo(s).", len(changes))
	} else {
		text += "\nAlterações locais precisam de revisão: " + changesErr.Error()
	}
	e.setGitStatus(text)
}

// gitOperation mantém I/O fora da thread Fyne e permite cancelar sem congelar a janela.
// title descreve a ação; operation recebe contexto cancelável e completed volta à thread
// da janela com o resultado sanitizado. O diálogo impede edição durante a operação.
func (e *Editor) gitOperation(title string, operation func(context.Context) error, completed func(error)) {
	if e.GitBusy || e.GitHub == nil {
		return
	}
	// ctx limita também operações Git; aplicação local interrompida será recuperada pelo journal.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	ctx = application.WithSyncOperation(ctx)
	// started mede toda a ação, sem registrar título livre ou parâmetros de autenticação.
	started := time.Now()
	e.GitDiagnostics.Event(ctx, "operation.start", nil)
	e.GitBusy = true
	e.GitCancel = cancel
	// progress é modal para impedir edições concorrentes na sessão durante a integração de arquivos.
	progress := dialog.NewCustom(title, "Cancelar", container.NewVBox(widget.NewLabel(title), widget.NewProgressBarInfinite()), e.Window)
	progress.SetOnClosed(cancel)
	progress.Show()
	go func() {
		// result não transporta tokens nem resposta bruta de OAuth.
		result := func() (failure error) {
			defer func() {
				if value := recover(); value != nil {
					e.GitDiagnostics.Event(ctx, "operation.panic", syncdiag.PanicFields(value))
					failure = application.SyncProblem("unexpectedFailure", "Erro inesperado na integração GitHub; consulte o Log de aplicação.")
				}
			}()
			return operation(ctx)
		}()
		e.GitDiagnostics.Complete(ctx, "operation.complete", started, result)
		fyne.Do(func() {
			cancel()
			if e.GitClosed {
				return
			}
			progress.Hide()
			e.GitBusy = false
			e.GitCancel = nil
			// Força recuperação ou bloqueio antes de liberar salvar após falha de apply.
			if result != nil {
				e.GitRecoveredRoot = ""
			}
			e.refreshWorkspace()
			if result != nil && !errors.Is(result, context.Canceled) {
				dialog.ShowError(fmt.Errorf("%s\n\nDiagnóstico: %s\nConsulte o botão Log de aplicação.", result.Error(), application.SyncOperationID(ctx)), e.Window)
			}
			if completed != nil {
				completed(result)
			}
		})
	}()
}

// connectGitHub retoma a mesma janela depois da autorização no navegador padrão.
func (e *Editor) connectGitHub() {
	// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
	e.gitOperation("Aguardando autorização no navegador", func(ctx context.Context) error {
		// account não contém credenciais e só confirma a conexão persistida.
		_, err := e.GitHub.Connect(ctx, e.App.OpenURL)
		return err
		// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
	}, func(err error) {
		if err == nil {
			e.GitVerifiedAt = time.Time{}
			e.Window.Show()
			e.Window.RequestFocus()
			e.selectGitHubRepository()
		}
	})
}

// ensureGitSaved impede publicar rascunhos de formulário ou edição ainda não salva.
// next só é chamado após sessão limpa ou salvamento concluído; cancelamento não o executa.
func (e *Editor) ensureGitSaved(next func()) {
	if e.InspectorHasDraft != nil && e.InspectorHasDraft() {
		dialog.ShowInformation("Alterações de propriedades", "Aplique ou descarte as alterações do painel de propriedades antes de continuar.", e.Window)
		return
	}
	if !e.Session.Dirty() && e.FilePath != "" {
		next()
		return
	}
	// yes registra a confirmação explícita antes de salvar, conectar ou publicar.
	dialog.ShowConfirm("Salvar antes de continuar?", "As alterações precisam ser salvas localmente antes da operação GitHub.", func(yes bool) {
		if yes {
			// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
			e.saveWithCompletion(false, func(err error) {
				if err == nil {
					next()
				}
			})
		}
	}, e.Window)
}

// selectGitHubRepository carrega a lista inteira em background antes da busca global.
func (e *Editor) selectGitHubRepository() {
	if e.GitHub == nil {
		return
	}
	// repositories só contém metadados públicos permitidos pela conta.
	var repositories []application.GitHubRepository
	// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
	e.gitOperation("Carregando repositórios GitHub", func(ctx context.Context) error {
		// err preserva a falha da consulta sem apresentar uma lista parcial como completa.
		var err error
		repositories, err = e.GitHub.Repositories(ctx)
		return err
		// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
	}, func(err error) {
		if err == nil {
			e.showGitRepositoryDialog(repositories)
		}
	})
}

// showGitRepositoryDialog exibe visibilidade/escrita e não permite confirmar destino inelegível.
// repositories é a lista autorizada completa; filtrar e selecionar não criam vínculo ou fazem push.
func (e *Editor) showGitRepositoryDialog(repositories []application.GitHubRepository) {
	// destination é o valor selecionado; nenhuma segunda digitação altera essa identidade.
	e.chooseGitRepository(repositories, func(destination application.GitHubRepository) {
		e.ensureGitSaved(func() { e.bindGitRepository(destination) })
	})
}

// chooseGitRepository separa seleção visual da conexão assíncrona para testar o destino por valor.
// confirmed recebe apenas a identidade selecionada; cancelar não chama esse callback.
func (e *Editor) chooseGitRepository(repositories []application.GitHubRepository, confirmed func(application.GitHubRepository)) {
	// filtered mantém IDs independentes dos rótulos exibidos na lista.
	filtered := append([]application.GitHubRepository{}, repositories...)
	// selected identifica a linha atual e é invalidado sempre que a busca muda.
	selected := -1
	// search contém somente o texto de filtro da lista de destinos autorizados.
	search := widget.NewEntry()
	search.SetPlaceHolder("Buscar proprietário ou repositório")
	// detail explica o destino selecionado ou o impedimento de escrita.
	detail := widget.NewLabel("Selecione um repositório com permissão de escrita.")
	detail.Wrapping = fyne.TextWrapWord
	// chooser mantém o seletor aberto até confirmação, atualização ou cancelamento explícitos.
	var chooser dialog.Dialog
	// confirm só permite vincular uma seleção elegível, ainda sujeita à revalidação remota.
	confirm := widget.NewButton("Conectar repositório", func() {
		if selected >= 0 && selected < len(filtered) {
			// destination fixa o repositório escolhido antes de fechar o seletor.
			destination := filtered[selected]
			if !destination.Writable || destination.Archived || destination.Disabled {
				return
			}
			chooser.Hide()
			confirmed(destination)
		}
	})
	confirm.Disable()
	// list virtualiza a apresentação sem confundir os IDs da lista com IDs GitHub.
	list := widget.NewList(func() int { return len(filtered) }, func() fyne.CanvasObject {
		return widget.NewLabel("proprietário/repositório\nvisibilidade · escrita")
		// id seleciona a linha e object é o widget reutilizado somente para sua apresentação.
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		// repo é o destino desta linha virtualizada, independente do rótulo do widget.
		repo := filtered[id]
		// permission apresenta a elegibilidade básica sem prometer aceitação de todas as políticas.
		permission := "sem escrita"
		if repo.Writable && !repo.Archived && !repo.Disabled {
			permission = "com escrita"
		}
		object.(*widget.Label).SetText(repo.FullName() + "\n" + repo.Visibility + " · " + permission)
	})
	// id identifica a seleção atual na lista filtrada, sem alterar o ID GitHub do destino.
	list.OnSelected = func(id widget.ListItemID) {
		selected = id
		// repo é o destino selecionado cuja elegibilidade determina a ação de confirmação.
		repo := filtered[id]
		detail.SetText("Destino selecionado: " + repo.FullName() + " · branch " + repo.Branch)
		confirm.Disable()
		if repo.Writable && !repo.Archived && !repo.Disabled {
			confirm.Enable()
		} else {
			detail.SetText("Sem escrita, arquivado ou indisponível: " + repo.FullName())
		}
	}
	// query é o filtro digitado; alterar a busca invalida qualquer seleção anterior.
	search.OnChanged = func(query string) {
		filtered = nil
		// repo é cada destino autorizado examinado pelo filtro textual, sem alterar sua identidade.
		for _, repo := range repositories {
			if strings.Contains(strings.ToLower(repo.FullName()), strings.ToLower(query)) {
				filtered = append(filtered, repo)
			}
		}
		selected = -1
		detail.SetText("Selecione um repositório com permissão de escrita.")
		confirm.Disable()
		list.UnselectAll()
		list.Refresh()
	}
	if len(repositories) == 0 {
		detail.SetText("Nenhum repositório acessível. Confira sua autorização e tente Atualizar.")
	}
	// content organiza a busca, lista e ações do seletor sem alterar dados do projeto.
	// explanation distingue conectar de enviar, e informa o próximo passo sem outro formulário.
	explanation := widget.NewLabel("Depois de conectar, você revisará os arquivos .jornada e confirmará a publicação.")
	explanation.Wrapping = fyne.TextWrapWord
	content := container.NewBorder(container.NewVBox(search, detail), container.NewVBox(explanation, widget.NewButton("Atualizar lista", func() { chooser.Hide(); e.selectGitHubRepository() }), confirm), nil, nil, list)
	chooser = dialog.NewCustom("Selecionar repositório GitHub", "Cancelar", content, e.Window)
	chooser.Resize(fyne.NewSize(720, 520))
	chooser.Show()
}

// bindGitRepository usa diretamente a seleção e abre a revisão da publicação após conectar.
// destination é a seleção revisada, revalidada pelo backend antes de modificar a pasta.
func (e *Editor) bindGitRepository(destination application.GitHubRepository) {
	// workspace fornece a raiz administrativa atual sem adicionar caminhos ao documento.
	workspace, err := e.workspace()
	if !e.check(err) {
		return
	}
	// root e active são capturados na thread Fyne, antes de começar I/O.
	root, active := workspace.Root, e.FilePath
	// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
	e.gitOperation("Integrando workspace", func(ctx context.Context) error {
		e.GitDiagnostics.Event(ctx, "repository.selection.confirmed", map[string]any{"repository_id": destination.ID, "repository_key": userconfig.Key(destination.FullName())})
		return e.GitSync.Bind(ctx, root, destination)
		// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
	}, func(err error) {
		e.reloadGitActive(root, active)
		if err == nil {
			e.GitVerifiedAt = time.Now()
			e.refreshGitLocal()
			e.Status.SetText("Workspace conectado; use Publicar no GitHub para enviar as jornadas.")
			// A revisão não faz push; o usuário ainda confirma arquivos e mensagem antes do envio.
			e.publishGitHub()
		}
	})
}

// publishGitHub oferece revisão explícita; salvar não chama esta operação.
func (e *Editor) publishGitHub() {
	e.ensureGitSaved(func() {
		// workspace fornece a raiz administrativa atual sem adicionar caminhos ao documento.
		workspace, err := e.workspace()
		if !e.check(err) {
			return
		}
		// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
		_, err = e.GitWorkspace.Binding(workspace.Root)
		if !e.check(err) {
			return
		}
		// review é obtida sem materializar novidades remotas.
		var review application.GitPublicationReview
		// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
		e.gitOperation("Verificando alterações e destino", func(ctx context.Context) error {
			// readErr permanece no backend; a UI recebe somente a revisão pública.
			var readErr error
			review, readErr = e.GitSync.ReviewPublication(ctx, workspace.Root)
			return readErr
			// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
		}, func(err error) {
			if err == nil {
				e.showGitPublication(review)
			}
		})
	})
}

// showGitPublication mostra todo o conjunto e confirma publicação e eventual atualização local.
// review fixa nomes, destino e OIDs revisados; o backend rejeita alterações posteriores do snapshot.
func (e *Editor) showGitPublication(review application.GitPublicationReview) {
	if !review.Pending && len(review.Changes) == 0 && len(review.RemoteChanges) == 0 {
		dialog.ShowInformation("GitHub", "Não há novidades para publicar.", e.Window)
		return
	}
	// lines apresenta inclusões, alterações e exclusões sem revelar conteúdo dos pacotes.
	lines := []string{"Pasta: " + review.Root, "Destino: " + review.Repository.FullName() + " · " + review.Repository.Branch}
	// account identifica a autoria pública exibida ao usuário sem expor credenciais.
	account, _ := e.GitHub.Account()
	lines = append(lines, "Autoria: "+account.Login)
	if review.Pending {
		lines = append(lines, "Reenviar/verificar commit pendente primeiro. Novas edições aguardam a próxima publicação.")
	}
	// change descreve somente uma inclusão, alteração ou exclusão elegível da revisão.
	for _, change := range review.Changes {
		lines = append(lines, change.Kind+": "+change.Path)
		if change.Kind == "exclusão" {
			lines = append(lines, e.gitDeletionReferences(review.Root, change.Path)...)
		}
	}
	if len(review.RemoteChanges) > 0 {
		lines = append(lines, "Arquivos locais serão atualizados após confirmação; a jornada ativa será recarregada se estiver válida.")
	}
	lines = append(lines, "A autorização OAuth pode abranger outros repositórios. Esta ação usa somente o destino acima.")
	// message permite revisar a mensagem antes da criação do commit.
	message := widget.NewEntry()
	message.SetText("Atualiza jornadas do workspace")
	// label apresenta o conjunto completo sem expor os bytes internos dos pacotes.
	label := widget.NewLabel(strings.Join(lines, "\n"))
	label.Wrapping = fyne.TextWrapWord
	// active fixa a jornada que pode ser recarregada depois da transação.
	active := e.FilePath
	// chooser mantém o diálogo da confirmação e seu ciclo de fechamento.
	// yes registra a confirmação explícita antes de salvar, conectar ou publicar.
	chooser := dialog.NewCustomConfirm("Publicar no GitHub", "Publicar no GitHub", "Cancelar", container.NewBorder(nil, widget.NewForm(widget.NewFormItem("Mensagem de commit", message)), nil, nil, container.NewVScroll(label)), func(yes bool) {
		if !yes {
			return
		}
		if strings.TrimSpace(message.Text) == "" {
			dialog.ShowInformation("Mensagem necessária", "Informe uma mensagem de commit.", e.Window)
			return
		}
		// text captura a mensagem confirmada antes de iniciar o trabalho em background.
		text := message.Text
		// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
		e.gitOperation("Publicando jornadas no GitHub", func(ctx context.Context) error {
			return e.GitSync.Publish(ctx, review, text)
			// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
		}, func(err error) {
			e.reloadGitActive(review.Root, active)
			if err == nil {
				e.GitVerifiedAt = time.Now()
				e.refreshGitLocal()
				e.showGitPublicationResult(review)
			}
		})
	}, e.Window)
	chooser.Resize(fyne.NewSize(720, 520))
	chooser.Show()
}

// showGitPublicationResult identifica o destino e o commit cuja presença remota foi verificada.
// review distingue publicação de mera atualização local; metadados ausentes não geram sucesso genérico.
func (e *Editor) showGitPublicationResult(review application.GitPublicationReview) {
	// binding contém somente a confirmação persistida pelo backend depois do push e da consulta remota.
	binding, err := e.GitWorkspace.Binding(review.Root)
	if !e.check(err) {
		return
	}
	if !review.Pending && len(review.Changes) == 0 {
		e.Status.SetText("Workspace atualizado em " + binding.Repository.FullName() + "; nenhuma nova publicação.")
		return
	}
	if binding.Repository.ID != review.Repository.ID || binding.Pending != "" || binding.Published == "" {
		dialog.ShowInformation("Publicação não confirmada", "Confira o destino e tente verificar o envio novamente. O commit local será preservado.", e.Window)
		return
	}
	e.Status.SetText("Publicado em " + binding.Repository.FullName() + " · " + binding.Repository.Branch)
	// address aponta diretamente para o commit no destino efetivamente confirmado, sem tokens.
	address, err := url.Parse("https://github.com/" + url.PathEscape(binding.Repository.Owner) + "/" + url.PathEscape(binding.Repository.Name) + "/commit/" + url.PathEscape(binding.Published))
	if !e.check(err) {
		return
	}
	// details permite conferir o repositório antes de abrir a evidência no navegador.
	details := widget.NewLabel("Repositório: " + binding.Repository.FullName() + "\nBranch: " + binding.Repository.Branch + "\nCommit: " + binding.Published + "\n\nAs alterações dos arquivos .jornada foram enviadas e confirmadas nesse destino.")
	details.Wrapping = fyne.TextWrapWord
	open := widget.NewButton("Ver commit no GitHub", func() { e.check(e.App.OpenURL(address)) })
	chooser := dialog.NewCustom("Publicação confirmada no GitHub", "Fechar", container.NewVBox(details, open), e.Window)
	chooser.Resize(fyne.NewSize(720, 300))
	chooser.Show()
}

// reloadGitActive só substitui uma sessão limpa por pacote íntegro da mesma pasta.
// root e active são capturados antes do I/O; mudança de sessão ou pacote ilegível impede a recarga.
func (e *Editor) reloadGitActive(root, active string) {
	if active == "" || e.FilePath != active || e.Session.Dirty() || !strings.EqualFold(filepath.Dir(active), root) {
		return
	}
	// document só substitui a sessão após leitura bem-sucedida; corrupção permanece diagnosticável na árvore.
	document, err := projectzip.Load(active)
	if err != nil {
		e.Status.SetText("Arquivos atualizados; jornada ativa não foi recarregada: " + err.Error())
		return
	}
	e.loadDocument(document, active)
}

// gitDeletionReferences informa referências locais afetadas sem repará-las por inferência.
// root limita a busca à pasta atual e filename é o destino explícito que será excluído.
func (e *Editor) gitDeletionReferences(root, filename string) []string {
	// notes reúne somente nomes de arquivos e não textos digitados no documento.
	notes := []string{}
	// entries enumera somente arquivos locais para diagnosticar referências afetadas por exclusão.
	entries, err := os.ReadDir(root)
	if err != nil {
		return notes
	}
	// item é a próxima entrada examinada, sem inferir nomes ou destinos alternativos.
	for _, item := range entries {
		if item.IsDir() || !strings.EqualFold(filepath.Ext(item.Name()), ".jornada") {
			continue
		}
		// document é lido pelo mesmo codec da jornada antes de examinar destinos explícitos.
		document, err := projectzip.Load(filepath.Join(root, item.Name()))
		if err != nil {
			continue
		}
		// A validação estruturada das folhas usa o mesmo documento do workspace.
		for _, screen := range document.Screens {
			// component fornece somente as propriedades do menu, sem inferir destinos por texto.
			for _, component := range screen.Components {
				if component.Properties.Menu != nil && menuReferences(component.Properties.Menu.Root, filename) {
					notes = append(notes, "Referência afetada em "+item.Name())
				}
			}
		}
	}
	return notes
}

// menuReferences visita destinos explícitos sem inferir referências por rótulo ou texto.
// menu contém a estrutura validada pelo codec e filename é comparado como nome físico Windows.
func menuReferences(menu catalog.SubmenuProperties, filename string) bool {
	// item é uma folha ou um submenu já validado pelo codec do documento.
	for _, item := range menu.Items {
		if strings.EqualFold(item.TargetJourney, filename) {
			return true
		}
		if item.Submenu != nil && menuReferences(*item.Submenu, filename) {
			return true
		}
	}
	return false
}

// refreshGitHubRemote revalida autorização e consulta pendências sem publicar.
func (e *Editor) refreshGitHubRemote() {
	if e.GitRecoveryError != nil {
		e.GitRecoveredRoot = ""
		e.GitRecoveryError = nil
		e.refreshGitLocal()
		return
	}
	e.publishGitHubReviewOnly()
}

// publishGitHubReviewOnly usa a revisão como consulta remota, sem confirmar commit ou push.
// Eventual journal interrompido é tratado pelo adaptador antes de calcular o snapshot.
func (e *Editor) publishGitHubReviewOnly() {
	// workspace fornece a raiz administrativa atual sem adicionar caminhos ao documento.
	workspace, err := e.workspace()
	if !e.check(err) {
		return
	}
	// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
	_, err = e.GitWorkspace.Binding(workspace.Root)
	if !e.check(err) {
		return
	}
	// ctx limita esta operação de rede à tentativa cancelável apresentada na janela.
	e.gitOperation("Atualizando estado GitHub", func(ctx context.Context) error {
		// readErr conserva a verificação remota sem expor token à interface.
		_, readErr := e.GitSync.ReviewPublication(ctx, workspace.Root)
		return readErr
		// err informa à thread Fyne se a etapa terminou antes de liberar a próxima ação.
	}, func(err error) {
		if err == nil {
			e.GitVerifiedAt = time.Now()
			e.refreshGitLocal()
			e.Status.SetText("Destino verificado; use Publicar para revisar as alterações ou verificar o envio pendente.")
		}
	})
}

// disconnectGitHub remove a credencial da aplicação sem remover arquivos ou histórico.
func (e *Editor) disconnectGitHub() {
	// yes registra a confirmação explícita antes de salvar, conectar ou publicar.
	dialog.ShowConfirm("Desconectar GitHub?", "A credencial desta aplicação será removida. Jornadas, histórico Git e vínculo local serão preservados.", func(yes bool) {
		if yes {
			e.gitOperation("Desconectando GitHub", e.GitHub.Disconnect, nil)
		}
	}, e.Window)
}

// openGitRepository usa somente o destino já vinculado, sem credenciais na URL.
func (e *Editor) openGitRepository() {
	// workspace fornece a raiz administrativa atual sem adicionar caminhos ao documento.
	workspace, err := e.workspace()
	if !e.check(err) {
		return
	}
	// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
	binding, err := e.GitWorkspace.Binding(workspace.Root)
	if !e.check(err) {
		return
	}
	// address é derivado do destino confirmado e não inclui credenciais na URL.
	address, _ := url.Parse("https://github.com/" + binding.Repository.FullName())
	_ = e.App.OpenURL(address)
}

// setGitStatus conserva o texto completo nos detalhes sem consumir a altura da árvore.
func (e *Editor) setGitStatus(value string) {
	e.GitStatusText = value
	e.GitLabel.SetText(strings.ReplaceAll(value, "\n", " · "))
}

// showGitMenu oferece ações secundárias sem esconder as jornadas do workspace.
func (e *Editor) showGitMenu() {
	// action conserva a elegibilidade já calculada pelos controles administrativos.
	action := func(text string, button *widget.Button) *fyne.MenuItem {
		item := fyne.NewMenuItem(text, button.OnTapped)
		item.Disabled = button.Disabled() || e.GitBusy
		return item
	}
	// menu reúne somente operações explícitas; abrir detalhes não consulta a rede.
	menu := fyne.NewMenu("GitHub", fyne.NewMenuItem("Detalhes da conexão", func() { dialog.ShowInformation("GitHub", e.GitStatusText, e.Window) }),
		action("Conectar ou trocar conta", widget.NewButton("", e.connectGitHub)),
		action("Selecionar repositório", e.GitRepositoriesButton), action("Publicar no GitHub", e.GitPublishButton),
		action("Atualizar estado", e.GitRefreshButton), action("Ver repositório no navegador", e.GitOpenButton), action("Desconectar GitHub", e.GitDisconnectButton))
	menu.Items[1].Disabled = !githubapi.Configured() || e.GitBusy || e.GitHub == nil
	widget.NewPopUpMenu(menu, e.Window.Canvas()).Show()
}
