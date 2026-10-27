package application

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// GitHubAuthorizer isola autorização e renovação sem acoplar os casos de uso ao HTTP.
type GitHubAuthorizer interface {
	// Authorize usa a tentativa e o abridor de URL fornecidos, retornando token ou falha.
	Authorize(context.Context, func(*url.URL) error) (GitHubCredential, error)
	// Refresh renova a credencial fornecida sob contexto cancelável ou exige reconexão.
	Refresh(context.Context, GitHubCredential) (GitHubCredential, error)
}

// GitHubDirectory consulta apenas identidade e destinos autorizados.
type GitHubDirectory interface {
	// Identity confere a identidade associada à credencial fornecida.
	Identity(context.Context, GitHubCredential) (GitHubAccount, error)
	// Repositories consulta os destinos autorizados pela credencial fornecida.
	Repositories(context.Context, GitHubCredential) ([]GitHubRepository, error)
	// Repository revalida o destino informado com a credencial atual.
	Repository(context.Context, GitHubCredential, GitHubRepository) (GitHubRepository, error)
}

// GitHubConnection coordena cofre e preferências; tokens não saem para a UI.
type GitHubConnection struct {
	// Observe recebe etapas e erros somente no backend, sem incluir credenciais nos campos públicos.
	Observe func(context.Context, string, error)
	// mutex serializa refresh, troca de conta e logout dentro do processo.
	mutex sync.Mutex
	// authorizer e directory são fronteiras de rede testáveis.
	authorizer GitHubAuthorizer
	// directory consulta identidade e destinos sem depender da UI.
	directory GitHubDirectory
	// vault e preferences separam credenciais de metadados locais.
	vault CredentialVault
	// preferences guarda somente identidade e referência opaca ao cofre.
	preferences SyncPreferences
	// namespace distingue configurações portáteis e testes.
	namespace string
}

// NewGitHubConnection recebe dependências explícitas, sem localizar widgets ou arquivos por conta própria.
// auth e directory fornecem rede; vault protege tokens; prefs mantém metadados;
// namespace isola a configuração. O retorno não inicia uma conexão implicitamente.
func NewGitHubConnection(auth GitHubAuthorizer, directory GitHubDirectory, vault CredentialVault, prefs SyncPreferences, namespace string) *GitHubConnection {
	return &GitHubConnection{authorizer: auth, directory: directory, vault: vault, preferences: prefs, namespace: namespace}
}

// Account lê a identidade local sem exigir internet.
// Retorna a conta persistida ou a falha de leitura/versão das preferências.
func (c *GitHubConnection) Account() (GitHubAccount, error) {
	// preference não contém tokens.
	var preference GitHubAccountPreference
	// err distingue ausência de configuração de uma leitura concluída.
	err := c.preferences.Read("github.json", &preference)
	return preference.Account, err
}

// Connect somente substitui a conta ativa após identidade e persistência válidas.
// ctx cancela a tentativa e open abre a URL no navegador do adaptador.
// Retorna a identidade verificada ou falha de autorização, cofre ou persistência.
func (c *GitHubConnection) Connect(ctx context.Context, open func(*url.URL) error) (connected GitHubAccount, connectErr error) {
	ctx = WithSyncOperation(ctx)
	c.report(ctx, "connection.start", nil)
	defer func() { c.report(ctx, "connection.complete", connectErr) }()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// token é mantido no backend durante a autorização.
	token, err := c.authorizer.Authorize(ctx, open)
	if err != nil {
		return GitHubAccount{}, err
	}
	// granted exige o escopo necessário para listar e publicar também repositórios privados.
	granted := false
	for _, scope := range strings.Fields(strings.ReplaceAll(token.Scopes, ",", " ")) {
		if scope == "repo" {
			granted = true
		}
	}
	if !granted {
		return GitHubAccount{}, fmt.Errorf("Autorize o acesso a repositórios (repo) para conectar o GitHub.")
	}
	// account confirma a identidade atual junto ao provedor.
	account, err := c.directory.Identity(ctx, token)
	c.report(ctx, "connection.identity.complete", err)
	if err != nil {
		return GitHubAccount{}, err
	}
	// key faz parte de um namespace privado; nenhuma credencial é salva no JSON.
	key := fmt.Sprintf("JornadaStudio/github/%s/%d", c.namespace, account.ID)
	// old registra somente a referência de credencial anterior para remoção posterior.
	var old GitHubAccountPreference
	if err = c.preferences.Read("github.json", &old); err != nil && !os.IsNotExist(err) {
		c.report(ctx, "preferences.account.read.failed", err)
		return GitHubAccount{}, err
	}
	// previous permite compensar falha de persistência sem abandonar a conta anterior.
	previous, previousErr := c.vault.Load(old.CredentialKey)
	// switching exige remover a credencial anterior antes de publicar a nova identidade.
	switching := old.CredentialKey != "" && old.CredentialKey != key
	if switching {
		if err = c.vault.Delete(old.CredentialKey); err != nil {
			return GitHubAccount{}, err
		}
	}
	if err = c.vault.Save(key, token); err != nil {
		c.report(ctx, "credentials.save.failed", err)
		if switching && previousErr == nil {
			_ = c.vault.Save(old.CredentialKey, previous)
		}
		return GitHubAccount{}, err
	}
	c.report(ctx, "credentials.saved", nil)
	if err = c.preferences.Write(ctx, "github.json", GitHubAccountPreference{Version: 1, Account: account, CredentialKey: key}); err != nil {
		c.report(ctx, "preferences.account.save.failed", err)
		_ = c.vault.Delete(key)
		if old.CredentialKey != "" && previousErr == nil {
			_ = c.vault.Save(old.CredentialKey, previous)
		}
		return GitHubAccount{}, err
	}
	c.report(ctx, "preferences.account.saved", nil)
	return account, nil
}

// credential renova o acesso sob mutex e não inicia publicação.
// ctx controla a eventual renovação; falha de leitura/renovação exige reconexão.
func (c *GitHubConnection) credential(ctx context.Context) (GitHubCredential, error) {
	// preference referencia somente o cofre da conta ativa.
	var preference GitHubAccountPreference
	// err impede usar preferência ausente, inválida ou ilegível como conexão válida.
	if err := c.preferences.Read("github.json", &preference); err != nil {
		c.report(ctx, "preferences.account.read.failed", err)
		return GitHubCredential{}, fmt.Errorf("Conecte sua conta GitHub antes de continuar.")
	}
	if preference.CredentialKey == "" {
		return GitHubCredential{}, fmt.Errorf("Reconecte sua conta GitHub.")
	}
	// token é renovado antes do vencimento, quando o provedor fornecer refresh.
	token, err := c.vault.Load(preference.CredentialKey)
	c.report(ctx, "credentials.load.complete", err)
	if err != nil {
		return token, err
	}
	if !token.ExpiresAt.IsZero() && time.Until(token.ExpiresAt) < time.Minute {
		token, err = c.authorizer.Refresh(ctx, token)
		c.report(ctx, "credentials.refresh.complete", err)
		if err != nil {
			return GitHubCredential{}, err
		}
		if err = c.vault.Save(preference.CredentialKey, token); err != nil {
			return GitHubCredential{}, err
		}
	}
	return token, nil
}

// Repositories obtém a lista completa antes de habilitar a seleção global.
// ctx controla o I/O; o retorno contém metadados, sem devolver a credencial à UI.
func (c *GitHubConnection) Repositories(ctx context.Context) ([]GitHubRepository, error) {
	ctx = WithSyncOperation(ctx)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// token nunca é devolvido ao chamador de interface.
	token, err := c.credential(ctx)
	if err != nil {
		return nil, err
	}
	return c.directory.Repositories(ctx, token)
}

// WithRepository revalida destino/permissão e disponibiliza credencial somente à operação backend.
// ctx controla a rede, repo fixa o destino esperado e operation recebe a conta e
// o destino revalidados sob mutex; seu erro é devolvido sem iniciar nova tentativa.
func (c *GitHubConnection) WithRepository(ctx context.Context, repo GitHubRepository, operation func(GitHubCredential, GitHubAccount, GitHubRepository) error) error {
	ctx = WithSyncOperation(ctx)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// token, account e fresh representam uma única identidade autorizada.
	token, err := c.credential(ctx)
	if err != nil {
		return err
	}
	// account confirma a identidade atual junto ao provedor.
	account, err := c.directory.Identity(ctx, token)
	c.report(ctx, "repository.identity.complete", err)
	if err != nil {
		return err
	}
	// expected é a conta local usada para impedir troca silenciosa de identidade.
	expected, err := c.Account()
	if err != nil {
		return err
	}
	if account.ID != expected.ID {
		return fmt.Errorf("Conta GitHub mudou; reconecte antes de publicar.")
	}
	// fresh revalida permissões, nome e branch antes de executar a operação.
	fresh, err := c.directory.Repository(ctx, token, repo)
	c.report(ctx, "repository.validation.complete", err)
	if err != nil {
		return err
	}
	if fresh.ID != repo.ID {
		return SyncProblem("repositoryMismatch", "O repositório retornado difere do destino confirmado; selecione novamente.")
	}
	if !fresh.Writable || fresh.Archived || fresh.Disabled {
		return fmt.Errorf("Repositório sem escrita ou indisponível para publicação.")
	}
	if repo.Branch != "" && fresh.Branch != repo.Branch || repo.FullName() != fresh.FullName() {
		return fmt.Errorf("Nome ou branch padrão do repositório mudou; revise o vínculo antes de publicar.")
	}
	return operation(token, account, fresh)
}

// Disconnect remove a credencial antes de declarar a conta desconectada.
// ctx limita a escrita do estado sem conta; falha de cofre preserva sua referência.
func (c *GitHubConnection) Disconnect(ctx context.Context) (disconnectErr error) {
	ctx = WithSyncOperation(ctx)
	c.report(ctx, "disconnect.start", nil)
	defer func() { c.report(ctx, "disconnect.complete", disconnectErr) }()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// preference é preservada se o cofre recusar remoção.
	var preference GitHubAccountPreference
	// err impede usar preferência ausente, inválida ou ilegível como conexão válida.
	if err := c.preferences.Read("github.json", &preference); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if preference.CredentialKey != "" {
		// err conserva falha de remoção para não declarar logout fictício.
		if err := c.vault.Delete(preference.CredentialKey); err != nil {
			c.report(ctx, "credentials.delete.failed", err)
			return err
		}
		c.report(ctx, "credentials.deleted", nil)
	}
	return c.preferences.Write(ctx, "github.json", GitHubAccountPreference{Version: 1})
}

// report encaminha o erro original somente à fronteira de diagnóstico, que nunca serializa Error().
func (c *GitHubConnection) report(ctx context.Context, event string, err error) {
	if c.Observe != nil {
		c.Observe(ctx, event, err)
	}
}
