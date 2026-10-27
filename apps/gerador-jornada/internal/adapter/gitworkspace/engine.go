package gitworkspace

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/adapter/userconfig"
	"gerador-jornada.local/editor/internal/application"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Engine opera somente na raiz explicitamente escolhida e não depende de git.exe.
type Engine struct {
	// Diagnostics registra somente etapas, classificações e identificadores opacos.
	Diagnostics *syncdiag.Logger
	// Store guarda vínculo e intenção fora dos arquivos de jornada.
	Store *userconfig.Store
	// remoteURL é substituído somente por fixtures do próprio pacote.
	remoteURL func(application.GitHubRepository) string
	// fail permite interromper fases de transação em testes determinísticos.
	fail func(string) error
}

// New compõe o adaptador Git com preferências privadas já resolvidas.
// store deve apontar para o namespace privado da instalação; o retorno não faz rede ou checkout.
func New(store *userconfig.Store) *Engine {
	return &Engine{Store: store, remoteURL: func(repo application.GitHubRepository) string {
		// repo fornece somente a identidade pública do destino, sem URL arbitrária ou credenciais.
		return "https://github.com/" + repo.FullName() + ".git"
	}}
}

// canonical exige pasta existente e impede usar acidentalmente um repositório ancestral.
// root é a seleção explícita do usuário; retorna a raiz física ou erro sem criar arquivos.
func canonical(root string) (string, error) {
	// absolute e resolved distinguem a seleção explícita de caminhos relativos/links.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	// resolved conserva a raiz física conferida, evitando usar alias ou metadados externos.
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	// info confere existência e tipo da pasta antes de aceitar seu conteúdo como workspace.
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("Workspace precisa ser uma pasta existente.")
	}
	// Bare mantém metadados diretamente na raiz e não pode receber outro .git aninhado.
	if _, metadataErr := os.Lstat(filepath.Join(resolved, ".git")); os.IsNotExist(metadataErr) {
		if repository, openErr := git.PlainOpen(resolved); openErr == nil {
			if configuration, configErr := repository.Config(); configErr == nil && configuration.Core.IsBare {
				return "", fmt.Errorf("Repositório bare não pode ser usado como workspace de jornadas.")
			}
		}
	}
	// parent percorre os ancestrais para não inicializar Git dentro de outro repositório.
	for parent := filepath.Dir(resolved); parent != filepath.Dir(parent); parent = filepath.Dir(parent) {
		if _, err = os.Lstat(filepath.Join(parent, ".git")); err == nil {
			return "", fmt.Errorf("Workspace está dentro de outro repositório Git; selecione a raiz correta.")
		}
	}
	return resolved, nil
}

// gate obtém trava liberada automaticamente pelo sistema no encerramento do processo.
// root é canonicalizada antes da trava; a função de liberação só é utilizável no sucesso.
func (e *Engine) gate(root string) (string, func(), error) {
	// resolved identifica uma única pasta mesmo quando escolhida por outro alias.
	resolved, err := canonical(root)
	if err != nil {
		return "", nil, err
	}
	// release encerra a trava da raiz ao concluir ou falhar a operação.
	release, err := lock(filepath.Join(e.Store.Root, "git-locks", userconfig.Key(resolved)+".lock"))
	return resolved, release, err
}

// bindingName mantém o vínculo fora do repositório e sem caminho legível no nome.
// root deve ser a raiz física já validada; o retorno é relativo à configuração privada.
func bindingName(root string) string { return "git-workspaces/" + userconfig.Key(root) + ".json" }

// Binding consulta exclusivamente estado local, sem afirmar sincronismo atual.
// root seleciona o workspace; ausência, versão inválida ou raiz divergente permanecem erros distinguíveis.
func (e *Engine) Binding(root string) (application.GitWorkspaceBinding, error) {
	// resolved impede colisão com caminho relativo usado por outro processo.
	resolved, err := canonical(root)
	if err != nil {
		return application.GitWorkspaceBinding{}, err
	}
	// binding recebe somente os metadados não secretos associados à raiz canônica.
	var binding application.GitWorkspaceBinding
	err = e.Store.Read(bindingName(resolved), &binding)
	if err == nil && !strings.EqualFold(binding.Root, resolved) {
		err = fmt.Errorf("Vínculo local não pertence ao workspace selecionado.")
	}
	return binding, err
}

// head devolve vazio para a branch ainda sem commit, inclusive o primeiro envio remoto.
// repo deve estar aberto na raiz exata; outras falhas de leitura não são tratadas como primeiro commit.
func head(repo *git.Repository) (string, error) {
	// ref pode não existir em repositório recém-inicializado.
	ref, err := repo.Head()
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

// inspect recusa operações/configurações que não podem ser preservadas com segurança.
// root fixa o repositório exato; destination e url delimitam branch e origin esperados.
// O retorno é somente um handle validado, sem alterar staging ou configuração.
func inspect(root string, destination application.GitHubRepository, url string) (*git.Repository, error) {
	// metadata precisa ser uma pasta real, nunca linked worktree ou junction.
	metadata := filepath.Join(root, ".git")
	// info confere existência e tipo da pasta antes de aceitar seu conteúdo como workspace.
	info, err := os.Lstat(metadata)
	if err != nil {
		return nil, err
	}
	// resolved conserva a raiz física conferida, evitando usar alias ou metadados externos.
	resolved, err := filepath.EvalSymlinks(metadata)
	if err != nil || !info.IsDir() || !strings.EqualFold(resolved, metadata) {
		return nil, fmt.Errorf("Linked worktree ou metadados Git externos não são suportados.")
	}
	// name identifica o caminho ou rótulo específico examinado nesta etapa.
	for _, name := range []string{"index.lock", "HEAD.lock", "MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "rebase-merge", "rebase-apply", "shallow", "info/sparse-checkout"} {
		if _, err = os.Lstat(filepath.Join(metadata, filepath.FromSlash(name))); err == nil {
			return nil, fmt.Errorf("Operação Git incompleta ou recurso não suportado: %s", name)
		}
	}
	// repo é aberto somente na raiz exata, sem procurar ancestrais.
	repo, err := git.PlainOpen(root)
	if err != nil {
		return nil, err
	}
	// cfg permite conferir ou preparar somente a configuração local deste repositório.
	cfg, err := repo.Config()
	if err != nil {
		return nil, err
	}
	if cfg.Core.IsBare || cfg.Core.Worktree != "" {
		return nil, fmt.Errorf("Configuração de worktree não suportada.")
	}
	// origin deve coincidir com o destino confirmado, sem substituir histórico ou credenciais.
	origin := cfg.Remotes["origin"]
	if origin == nil || len(origin.URLs) != 1 || !strings.EqualFold(origin.URLs[0], url) {
		return nil, fmt.Errorf("Origin difere do repositório selecionado; ajuste externamente sem perder seu histórico.")
	}
	// ref conserva a referência consultada e distingue ausência de falha de leitura.
	ref, err := repo.Reference(plumbing.HEAD, false)
	if err != nil {
		return nil, err
	}
	if ref.Type() != plumbing.SymbolicReference || ref.Target() != plumbing.NewBranchReferenceName(destination.Branch) {
		return nil, fmt.Errorf("A branch ativa difere da branch padrão selecionada.")
	}
	// oid e files fixam o conteúdo esperado do índice sem checkout ou varredura de arquivos omitidos.
	oid, err := head(repo)
	if err != nil {
		return nil, err
	}
	files, err := loadTree(repo, oid)
	if err != nil {
		return nil, err
	}
	// staged conserva entradas completas, incluindo skip-worktree; flags não ocultam staging externo.
	staged, err := repo.Storer.Index()
	if err != nil {
		return nil, err
	}
	if len(staged.Entries) != len(files) {
		return nil, fmt.Errorf("Há arquivos preparados por outro cliente Git; finalize ou retire o staging antes de publicar.")
	}
	// item deve coincidir com o blob e modo de HEAD, independente de sua presença no disco.
	for _, item := range staged.Entries {
		previous, ok := files[item.Name]
		if !ok || item.Stage != 0 || item.IntentToAdd || item.Hash != plumbing.ComputeHash(plumbing.BlobObject, previous.Data) || item.Mode != previous.Mode {
			return nil, fmt.Errorf("Há arquivos preparados por outro cliente Git; finalize ou retire o staging antes de publicar.")
		}
	}
	return repo, nil
}

// auth mantém token exclusivamente em memória no transporte HTTPS.
// address não HTTPS recebe nil para fixtures locais; produção deriva somente URLs GitHub HTTPS.
func auth(address string, token application.GitHubCredential) transport.AuthMethod {
	if !strings.HasPrefix(address, "https://") {
		return nil
	}
	return &githttp.BasicAuth{Username: "x-access-token", Password: token.AccessToken}
}

// fetch atualiza somente a referência remota da branch selecionada.
// ctx controla I/O; binding fixa a branch; address/token definem o transporte sem persistir segredo.
// Retorna OID vazio somente quando a referência não existe, preservando falhas como erros sanitizados.
func fetch(ctx context.Context, repo *git.Repository, binding application.GitWorkspaceBinding, address string, token application.GitHubCredential) (string, error) {
	// branch é validada antes de gerar qualquer refspec.
	branch := plumbing.NewBranchReferenceName(binding.Repository.Branch)
	// err impede construir refspec a partir de uma branch inválida.
	if err := branch.Validate(); err != nil {
		return "", fmt.Errorf("Branch GitHub inválida.")
	}
	// remoteRef fixa o nome local da branch remota selecionada antes do fetch.
	remoteRef := plumbing.NewRemoteReferenceName("origin", binding.Repository.Branch)
	// err distingue remoto vazio/sem novidades de falha que impede confirmar a referência consultada.
	err := repo.FetchContext(ctx, &git.FetchOptions{RemoteName: "origin", Auth: auth(address, token), RefSpecs: []config.RefSpec{config.RefSpec("+" + branch.String() + ":" + remoteRef.String())}})
	if errors.Is(err, transport.ErrEmptyRemoteRepository) {
		// Ausência no servidor não pode ser substituída pela referência antiga em cache.
		return "", nil
	}
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) && !errors.Is(err, transport.ErrEmptyRemoteRepository) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", application.WrapSyncProblem("networkUnavailable", "Não foi possível obter a branch GitHub; verifique rede, permissão e políticas.", err)
	}
	// ref conserva a referência consultada e distingue ausência de falha de leitura.
	ref, err := repo.Reference(remoteRef, true)
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

// ancestor trata ausência de commit como base inicial, sem chamar MergeBase com hash nulo.
// left é o candidato ancestral e right o descendente no armazenamento repo.
func ancestor(repo *git.Repository, left, right string) (bool, error) {
	if left == "" {
		return true, nil
	}
	if right == "" {
		return false, nil
	}
	if left == right {
		return true, nil
	}
	// a identifica o commit potencialmente ancestral, cuja ausência já foi tratada.
	a, err := repo.CommitObject(plumbing.NewHash(left))
	if err != nil {
		return false, err
	}
	// b identifica o commit descendente contra o qual a ancestralidade será verificada.
	b, err := repo.CommitObject(plumbing.NewHash(right))
	if err != nil {
		return false, err
	}
	return a.IsAncestor(b)
}

// Bind inicializa ou reaproveita Git sem criar commit nem publicar durante a conexão.
// ctx cancela I/O; root e destination são a pasta e destino confirmados; token é usado só no transporte.
// Colisão, configuração incompatível ou erro parcial preservam os dados e a intenção recuperável.
func (e *Engine) Bind(ctx context.Context, root string, destination application.GitHubRepository, token application.GitHubCredential) (bindErr error) {
	// started delimita toda a conexão sem registrar caminho, owner ou nome de repositório.
	ctx, started := e.startDiagnostics(ctx, "git.bind.start", root)
	// phase identifica a última etapa iniciada, sem registrar caminhos, documentos ou credenciais.
	phase := "inspection"
	e.Diagnostics.Event(ctx, "git.destination.confirmed", map[string]any{"repository_id": destination.ID, "repository_key": userconfig.Key(destination.FullName()), "branch_present": destination.Branch != ""})
	defer func() {
		if bindErr != nil {
			// fields combina classificação segura com a fase fixa, nunca com Error() livre.
			fields := syncdiag.ErrorFields(bindErr)
			fields["checkpoint"] = phase
			e.Diagnostics.Event(ctx, "git.bind.failed", fields)
		}
		e.Diagnostics.Complete(ctx, "git.bind.complete", started, bindErr)
	}()
	// resolved e release delimitam toda mutação deste workspace.
	resolved, release, err := e.gate(root)
	if err != nil {
		return err
	}
	defer release()
	if err = e.recover(ctx, resolved); err != nil {
		return err
	}
	if destination.Branch == "" {
		destination.Branch = "main"
	}
	if !destination.Writable || destination.Archived || destination.Disabled {
		return fmt.Errorf("Repositório sem escrita.")
	}
	// binding só será publicado depois que arquivos e metadados estiverem consistentes.
	binding := application.GitWorkspaceBinding{Version: 1, Root: resolved, Repository: destination}
	// address é derivado do destino confirmado e não inclui credenciais na URL.
	address := e.remoteURL(destination)
	// existing e loadErr distinguem vínculo anterior compatível de estado inválido ou ausente.
	if existing, loadErr := e.Binding(resolved); loadErr == nil && existing.Repository.ID != destination.ID {
		return fmt.Errorf("Esta pasta já está vinculada a %s, diferente de %s. Use outra pasta com suas jornadas para conectar o destino desejado; o histórico existente será preservado.", existing.Repository.FullName(), destination.FullName())
	} else if loadErr != nil && !os.IsNotExist(loadErr) {
		return loadErr
	} else if loadErr == nil {
		if existing.Pending != "" {
			return fmt.Errorf("Este workspace possui envio pendente; use Publicar para verificar ou reenviar o commit.")
		}
		binding.Published, binding.PublishedAt = existing.Published, existing.PublishedAt
	}
	if _, err = os.Lstat(filepath.Join(resolved, ".git")); err == nil {
		// repo acessa somente o repositório exato desta raiz ou do staging privado.
		repo, err := inspect(resolved, destination, address)
		if err != nil {
			return err
		}
		// local é o OID local observado, vazio enquanto a branch ainda não tiver commit.
		local, err := head(repo)
		if err != nil {
			return err
		}
		// remote é o OID confirmado pelo fetch da branch selecionada, ou vazio no primeiro envio.
		phase = "remote_fetch"
		remote, err := fetch(ctx, repo, binding, address, token)
		if err != nil {
			return err
		}
		// compatible registra a relação de histórico exigida para não enviar commits locais desconhecidos.
		compatible, err := ancestor(repo, local, remote)
		if err != nil {
			return err
		}
		if !compatible {
			return fmt.Errorf("Histórico local adiantado ou incompatível; resolva externamente antes de vincular.")
		}
		// oldTree conserva bytes e modos da referência local anterior à integração.
		phase = "remote_inventory"
		oldTree, err := loadTree(repo, local)
		if err != nil {
			return err
		}
		// newTree é a árvore remota validada que deve ser integrada sem descartar conteúdo local.
		newTree, err := loadTree(repo, remote)
		if err != nil {
			return err
		}
		phase = "workspace_integration"
		return e.apply(ctx, repo, binding, local, remote, oldTree, newTree, false, "")
	} else if !os.IsNotExist(err) {
		return err
	}
	// staging fica no volume do workspace para instalar .git por rename atômico.
	staging, err := os.MkdirTemp(filepath.Dir(resolved), ".jornada-git-")
	if err != nil {
		return err
	}
	defer e.removeStage(resolved, staging)
	// repo acessa somente o repositório exato desta raiz ou do staging privado.
	repo, err := git.PlainInit(staging, true)
	if err != nil {
		return err
	}
	if _, err = repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{address}}); err != nil {
		return err
	}
	if err = repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName(destination.Branch))); err != nil {
		return err
	}
	// remote é o OID confirmado pelo fetch da branch selecionada, ou vazio no primeiro envio.
	phase = "remote_fetch"
	remote, err := fetch(ctx, repo, binding, address, token)
	if err != nil {
		return err
	}
	// files conserva a árvore ou listagem efetiva usada para definir os caminhos elegíveis.
	phase = "remote_inventory"
	files, err := loadTree(repo, remote)
	if err != nil {
		return err
	}
	// cfg permite conferir ou preparar somente a configuração local deste repositório.
	cfg, err := repo.Config()
	if err != nil {
		return err
	}
	cfg.Core.IsBare = false
	if err = repo.SetConfig(cfg); err != nil {
		return err
	}
	if remote != "" {
		if err = repo.Storer.SetReference(plumbing.NewHashReference(plumbing.NewBranchReferenceName(destination.Branch), plumbing.NewHash(remote))); err != nil {
			return err
		}
	}
	phase = "workspace_integration"
	return e.apply(ctx, repo, binding, "", remote, map[string]entry{}, files, true, staging)
}

// working sobrepõe somente jornadas elegíveis aos arquivos já commitados.
// root fornece os bytes salvos e base preserva caminhos não elegíveis; retorna árvore e revisão ordenada.
func working(root string, base map[string]entry) (map[string]entry, []application.GitChange, error) {
	// result preserva todos os caminhos não elegíveis da base.
	result := map[string]entry{}
	// name e item copiam a base sem permitir que arquivos alheios entrem no delta do commit.
	for name, item := range base {
		result[name] = item
	}
	// files conserva a árvore ou listagem efetiva usada para definir os caminhos elegíveis.
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, err
	}
	// paths inclui jornadas atuais e anteriormente rastreadas para representar também exclusões.
	paths := map[string]bool{}
	// name identifica o caminho ou rótulo específico examinado nesta etapa.
	for name := range base {
		if !strings.Contains(name, "/") && strings.EqualFold(filepath.Ext(name), ".jornada") {
			paths[name] = true
		}
	}
	// item é a próxima entrada examinada, sem inferir nomes ou destinos alternativos.
	for _, item := range files {
		if strings.EqualFold(filepath.Ext(item.Name()), ".jornada") {
			paths[item.Name()] = true
		}
	}
	// changes é ordenado para que o snapshot e a revisão sejam determinísticos.
	changes := []application.GitChange{}
	// name identifica o caminho ou rótulo específico examinado nesta etapa.
	for name := range paths {
		// data conserva o resultado específico desta etapa para a validação antes de prosseguir.
		data, exists, err := fileData(root, name)
		if err != nil {
			return nil, nil, err
		}
		// before e beforeOK distinguem o arquivo anteriormente rastreado de uma nova jornada.
		before, beforeOK := base[name]
		// after representa os bytes de trabalho que serão comparados com o estado rastreado.
		after := entry{Data: data, Mode: filemode.Regular}
		if beforeOK {
			after.Mode = before.Mode
		}
		if equal(before, beforeOK, after, exists) {
			continue
		}
		// change descreve somente uma inclusão, alteração ou exclusão elegível da revisão.
		change := application.GitChange{Path: name, Kind: "alteração"}
		if !beforeOK {
			change.Kind = "inclusão"
		}
		if !exists {
			change.Kind = "exclusão"
			delete(result, name)
		} else {
			change.Digest = digest(data)
			result[name] = after
		}
		changes = append(changes, change)
	}
	// i/j estabilizam a ordem dos caminhos para comparar revisões sem depender da ordem do mapa.
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return result, changes, validateTree(result)
}

// prepared reúne dados privados do plano sem expor bytes ou credencial à interface.
type prepared struct {
	// review é a projeção pública revisável.
	review application.GitPublicationReview
	// old e next representam árvores imutáveis antes/depois da integração.
	old map[string]entry
	// next é o resultado que poderá ser materializado somente após a confirmação.
	next map[string]entry
	// parents define zero, um ou dois pais conforme o primeiro envio ou integração.
	parents []plumbing.Hash
	// delivered indica que a pendência já está alcançável no remoto.
	delivered bool
}

// prepare calcula integração e colisões sem modificar arquivos de trabalho.
// ctx controla fetch; repo/binding fixam a raiz e o destino; token não integra a revisão retornada.
func (e *Engine) prepare(ctx context.Context, repo *git.Repository, binding application.GitWorkspaceBinding, token application.GitHubCredential) (prepared, error) {
	// plan fixa os OIDs observados para a confirmação posterior.
	plan := prepared{}
	// local é o OID local observado, vazio enquanto a branch ainda não tiver commit.
	local, err := head(repo)
	if err != nil {
		return plan, err
	}
	// remote é o OID confirmado pelo fetch da branch selecionada, ou vazio no primeiro envio.
	remote, err := fetch(ctx, repo, binding, e.remoteURL(binding.Repository), token)
	if err != nil {
		return plan, err
	}
	// oldTree conserva bytes e modos da referência local anterior à integração.
	oldTree, err := loadTree(repo, local)
	if err != nil {
		return plan, err
	}
	// remoteTree preserva o conteúdo do destino ao combinar novidades com o conjunto local.
	remoteTree, err := loadTree(repo, remote)
	if err != nil {
		return plan, err
	}
	// localTree e changes representam o conteúdo salvo e o conjunto revisável das jornadas.
	localTree, changes, err := working(binding.Root, oldTree)
	if err != nil {
		return plan, err
	}
	plan.review = application.GitPublicationReview{Root: binding.Root, Repository: binding.Repository, Head: local, Remote: remote, Changes: changes, Pending: binding.Pending != ""}
	plan.old = oldTree
	if binding.Pending != "" {
		if local != binding.Pending {
			return plan, fmt.Errorf("HEAD foi alterado por outro cliente; não será enviado histórico desconhecido.")
		}
		plan.delivered, err = ancestor(repo, local, remote)
		if err != nil {
			return plan, err
		}
		if plan.delivered {
			plan.next = oldTree
			return plan, nil
		}
		localTree = oldTree
	} else {
		// compatible registra a relação de histórico exigida para não enviar commits locais desconhecidos.
		compatible, err := ancestor(repo, local, remote)
		if err != nil {
			return plan, err
		}
		if !compatible {
			return plan, fmt.Errorf("Commits locais desconhecidos; resolva externamente antes de publicar.")
		}
	}
	// base é HEAD para avanço normal ou o ancestral comum único para divergência.
	base := oldTree
	// remoteBehind identifica retentativa que pode preservar o mesmo commit sem merge.
	remoteBehind, err := ancestor(repo, remote, local)
	if err != nil {
		return plan, err
	}
	if binding.Pending != "" && remoteBehind {
		plan.next = oldTree
		return plan, nil
	}
	if binding.Pending != "" && local != remote {
		if remote == "" {
			return plan, fmt.Errorf("Branch remota foi removida; confirme o destino externamente.")
		}
		// left é o commit local usado na busca de ancestral comum antes de integrar divergência.
		left, err := repo.CommitObject(plumbing.NewHash(local))
		if err != nil {
			return plan, err
		}
		// right é o commit remoto usado na busca de ancestral comum antes de integrar divergência.
		right, err := repo.CommitObject(plumbing.NewHash(remote))
		if err != nil {
			return plan, err
		}
		// bases deve conter uma única ancestral comum para a integração conservadora suportada.
		bases, err := left.MergeBase(right)
		if err != nil || len(bases) != 1 {
			return plan, fmt.Errorf("Históricos incompatíveis ou múltiplas bases; resolva externamente.")
		}
		base, err = loadTree(repo, bases[0].Hash.String())
		if err != nil {
			return plan, err
		}
		// diff local somente pode conter jornadas; histórico remoto herdado é preservado.
		// name/item identificam o arquivo local a comparar com a base antes de aceitar seu delta.
		for name, item := range localTree {
			// previous e ok identificam o conteúdo da base para auditar deltas locais.
			previous, ok := base[name]
			if !equal(item, true, previous, ok) && (strings.Contains(name, "/") || !strings.EqualFold(filepath.Ext(name), ".jornada")) {
				return plan, fmt.Errorf("Histórico local contém alterações alheias às jornadas.")
			}
		}
	}
	plan.next, err = merge(base, localTree, remoteTree)
	if err != nil {
		return plan, err
	}
	if remote != "" {
		plan.parents = append(plan.parents, plumbing.NewHash(remote))
	}
	if binding.Pending != "" {
		plan.parents = append(plan.parents, plumbing.NewHash(local))
	}
	// name e item identificam conteúdo novo para listar atualizações locais previstas na revisão.
	for name, item := range plan.next {
		// old e ok identificam o estado anterior para apresentar arquivos que serão atualizados.
		old, ok := oldTree[name]
		if journeyPath(name) && !equal(old, ok, item, true) {
			plan.review.RemoteChanges = append(plan.review.RemoteChanges, name)
		}
	}
	// name identifica o caminho ou rótulo específico examinado nesta etapa.
	for name := range oldTree {
		// ok conserva o resultado específico desta etapa para a validação antes de prosseguir.
		if _, ok := plan.next[name]; journeyPath(name) && !ok {
			plan.review.RemoteChanges = append(plan.review.RemoteChanges, name)
		}
	}
	sort.Strings(plan.review.RemoteChanges)
	// preflight protege somente jornadas; demais arquivos locais nunca serão tocados.
	if err = preflight(binding.Root, oldTree, plan.next); err != nil {
		return plan, err
	}
	return plan, nil
}

// Review retorna snapshot sem fazer checkout das novidades remotas ou criar commit novo.
// ctx controla I/O e token autoriza fetch para root; uma intenção local interrompida
// é recuperada antes da revisão, podendo concluir ou desfazer suas escritas conhecidas.
func (e *Engine) Review(ctx context.Context, root string, token application.GitHubCredential) (review application.GitPublicationReview, reviewErr error) {
	ctx, started := e.startDiagnostics(ctx, "git.review.start", root)
	defer func() { e.Diagnostics.Complete(ctx, "git.review.complete", started, reviewErr) }()
	// resolved/release impedem duas consultas de mutação concorrentes no mesmo namespace.
	resolved, release, err := e.gate(root)
	if err != nil {
		return application.GitPublicationReview{}, err
	}
	defer release()
	if err = e.recover(ctx, resolved); err != nil {
		return application.GitPublicationReview{}, err
	}
	// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
	binding, err := e.Binding(resolved)
	if err != nil {
		return application.GitPublicationReview{}, err
	}
	// repo acessa somente o repositório exato desta raiz ou do staging privado.
	repo, err := inspect(resolved, binding.Repository, e.remoteURL(binding.Repository))
	if err != nil {
		return application.GitPublicationReview{}, err
	}
	// plan mantém árvores e OIDs privados usados para revalidar a revisão na confirmação.
	plan, err := e.prepare(ctx, repo, binding, token)
	return plan.review, err
}

// Publish confirma a revisão, cria no máximo um commit e mantém pendência em qualquer falha de push.
// ctx controla I/O; review fixa o snapshot; message e account compõem mensagem/autoria confirmadas.
// token fica somente no transporte. Revisão obsoleta, conflito ou erro de envio não descartam commits.
func (e *Engine) Publish(ctx context.Context, review application.GitPublicationReview, message string, account application.GitHubAccount, token application.GitHubCredential) (publishErr error) {
	ctx, started := e.startDiagnostics(ctx, "git.publish.start", review.Root)
	e.Diagnostics.Event(ctx, "git.destination.confirmed", map[string]any{"repository_id": review.Repository.ID, "repository_key": userconfig.Key(review.Repository.FullName()), "branch_present": review.Repository.Branch != ""})
	e.Diagnostics.Event(ctx, "git.publish.review", map[string]any{"pending": review.Pending, "change_count": len(review.Changes), "remote_change_count": len(review.RemoteChanges)})
	defer func() { e.Diagnostics.Complete(ctx, "git.publish.complete", started, publishErr) }()
	if review.AccountID != account.ID || account.ID <= 0 {
		return application.SyncProblem("staleReview", "A conta mudou após a revisão; revise novamente antes de publicar.")
	}
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("Informe uma mensagem de commit.")
	}
	// resolved/release fixam a raiz até concluir a aplicação transacional.
	resolved, release, err := e.gate(review.Root)
	if err != nil {
		return err
	}
	defer release()
	if err = e.recover(ctx, resolved); err != nil {
		return err
	}
	// binding conserva o destino e as pendências locais antes de revalidar ou publicar.
	binding, err := e.Binding(resolved)
	if err != nil {
		return err
	}
	// repo acessa somente o repositório exato desta raiz ou do staging privado.
	repo, err := inspect(resolved, binding.Repository, e.remoteURL(binding.Repository))
	if err != nil {
		return err
	}
	// plan mantém árvores e OIDs privados usados para revalidar a revisão na confirmação.
	plan, err := e.prepare(ctx, repo, binding, token)
	if err != nil {
		return err
	}
	plan.review.AccountID = account.ID
	if !reflect.DeepEqual(review, plan.review) {
		return application.SyncProblem("staleReview", "A revisão mudou (arquivos, destino ou remoto); revise novamente antes de publicar.")
	}
	if plan.delivered {
		binding.Published, binding.Pending, binding.PublishedAt = binding.Pending, "", time.Now().UTC()
		return e.Store.Write(ctx, bindingName(resolved), binding)
	}
	// target pode permanecer o mesmo OID em retentativa sem novidades remotas.
	target := review.Head
	if !review.Pending && len(review.Changes) == 0 {
		target = review.Remote
	} else if !review.Pending || !reflect.DeepEqual(plan.old, plan.next) || len(plan.parents) == 2 {
		// treeHash identifica a árvore imutável do commit ainda sem referência publicada.
		treeHash, err := writeTree(repo, plan.next)
		if err != nil {
			return err
		}
		// name escolhe a autoria pública, usando login quando o perfil não informar nome.
		name := account.Name
		if strings.TrimSpace(name) == "" {
			name = account.Login
		}
		// signature registra a autoria pública confirmada sem alterar a configuração Git global.
		signature := object.Signature{Name: name, Email: fmt.Sprintf("%d+%s@users.noreply.github.com", account.ID, account.Login), When: time.Now()}
		// commit preserva os pais planejados e os bytes da árvore confirmada pelo usuário.
		commit := &object.Commit{Author: signature, Committer: signature, Message: message, TreeHash: treeHash, ParentHashes: plan.parents}
		// encoded recebe o commit completo antes de gravá-lo como objeto imutável.
		encoded := repo.Storer.NewEncodedObject()
		if err = commit.Encode(encoded); err != nil {
			return err
		}
		// oid é a identidade do commit criado, persistida como pendência antes do push.
		oid, err := repo.Storer.SetEncodedObject(encoded)
		if err != nil {
			return err
		}
		target = oid.String()
		binding.Pending = target
		e.Diagnostics.Event(ctx, "git.commit.prepared", map[string]any{"commit_created": true, "commit_key": userconfig.Key(target), "pending": true})
	}
	if target != review.Head {
		if err = e.apply(ctx, repo, binding, review.Head, target, plan.old, plan.next, false, ""); err != nil {
			return err
		}
	}
	if binding.Pending == "" {
		return nil
	}
	if err = e.hit(ctx, "before-push"); err != nil {
		return fmt.Errorf("Commit salvo no computador; envio pendente.")
	}
	// push nunca inclui tags ou force; confirmação é feita pela ancestralidade remota.
	err = repo.PushContext(ctx, &git.PushOptions{RemoteName: "origin", Auth: auth(e.remoteURL(binding.Repository), token), RefSpecs: []config.RefSpec{config.RefSpec("refs/heads/" + binding.Repository.Branch + ":refs/heads/" + binding.Repository.Branch)}})
	e.Diagnostics.Event(ctx, "git.push.response", syncdiag.ErrorFields(err))
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return application.WrapSyncProblem("pushRejected", "Commit salvo no computador; envio pendente. Verifique rede, permissão e proteção de branch e tente novamente.", err)
	}
	if err = e.hit(ctx, "after-push"); err != nil {
		return application.SyncProblem("outcomeUnknown", "Envio sem confirmação; verifique novamente antes de considerar publicado.")
	}
	// remote é o OID confirmado pelo fetch da branch selecionada, ou vazio no primeiro envio.
	remote, err := fetch(ctx, repo, binding, e.remoteURL(binding.Repository), token)
	if err != nil {
		return fmt.Errorf("Envio sem confirmação; verifique novamente antes de considerar publicado.")
	}
	// confirmed só permite marcar publicado quando o commit é alcançável no destino.
	confirmed, err := ancestor(repo, binding.Pending, remote)
	if err != nil || !confirmed {
		return fmt.Errorf("GitHub ainda não confirmou o commit; envio permanece pendente.")
	}
	binding.Published, binding.Pending, binding.PublishedAt = target, "", time.Now().UTC()
	e.Diagnostics.Event(ctx, "git.push.confirmed", map[string]any{"confirmed": true, "commit_key": userconfig.Key(target)})
	return e.Store.Write(context.WithoutCancel(ctx), bindingName(resolved), binding)
}

// startDiagnostics herda correlação do caso de uso e observa transporte sem URLs ou conteúdo Git.
func (e *Engine) startDiagnostics(ctx context.Context, event, root string) (context.Context, time.Time) {
	ctx = application.WithSyncOperation(ctx)
	started := time.Now()
	if e.Diagnostics != nil {
		ctx = httptrace.WithClientTrace(ctx, e.Diagnostics.HTTPTrace(ctx, started))
	}
	e.Diagnostics.Event(ctx, event, map[string]any{"workspace_key": userconfig.Key(root)})
	return ctx, started
}
