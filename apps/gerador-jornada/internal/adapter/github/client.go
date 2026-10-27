package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/application"
)

// Client concentra somente transporte e conversão dos contratos do provedor.
type Client struct {
	// Diagnostics é opcional e recebe somente metadados técnicos sanitizados.
	Diagnostics *syncdiag.Logger
	// rateMutex/retryAt impedem repetir automaticamente uma requisição limitada pelo servidor.
	rateMutex sync.Mutex
	retryAt   time.Time
	// http possui timeout e política de redirecionamento sem vazamento de Authorization.
	http *http.Client
	// api e tokenEndpoint permitem fixtures HTTP locais no próprio pacote.
	api string
	// tokenEndpoint recebe as trocas OAuth em produção ou no fixture local.
	tokenEndpoint string
	// authorizeEndpoint e credenciais identificam o cadastro OAuth público do produto.
	authorizeEndpoint string
	// id é o identificador público do cadastro OAuth.
	id string
	// secret é o valor distribuído do cliente nativo, não o token do usuário.
	secret string
}

// New cria um transporte sem depender de Git, navegador embutido ou configuração global.
// id e secret identificam o cadastro distribuído; o retorno ainda não conecta a conta.
func New(id, secret string) *Client {
	return &Client{id: id, secret: secret, api: "https://api.github.com", tokenEndpoint: "https://github.com/login/oauth/access_token", authorizeEndpoint: "https://github.com/login/oauth/authorize", http: &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// req é o próximo destino e via contém a cadeia anterior a validar.
		if len(via) > 5 || len(via) > 0 && (req.URL.Host != via[0].URL.Host || req.URL.Scheme != via[0].URL.Scheme) {
			return fmt.Errorf("Redirecionamento GitHub recusado.")
		}
		return nil
	}}}
}

// request envia dados e devolve somente erros sanitizados, nunca o corpo OAuth.
// ctx limita I/O; method/address definem a chamada; access é usado só no header.
// body é serializado quando presente e result recebe o JSON esperado no sucesso.
func (c *Client) request(ctx context.Context, method, address, access string, body any, result any) (requestErr error) {
	ctx = application.WithSyncOperation(ctx)
	// started mede a requisição sem incorporar a URL de autorização ao log.
	started := time.Now()
	c.Diagnostics.Event(ctx, "http.request.start", map[string]any{"method": method, "http_route": c.route(address)})
	defer func() { c.Diagnostics.Complete(ctx, "http.request.complete", started, requestErr) }()
	if err := c.checkRate(); err != nil {
		return err
	}
	// payload limita a serialização a parâmetros explicitamente fornecidos.
	var payload io.Reader
	if body != nil {
		// data contém o corpo serializado; err impede envio de parâmetros inválidos.
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(data)
	}
	// req mantém credencial no header, fora de URLs e diagnósticos.
	req, err := http.NewRequestWithContext(ctx, method, address, payload)
	if err != nil {
		return fmt.Errorf("Requisição GitHub inválida.")
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// A troca OAuth usa negociação JSON própria, diferente da versão da API REST.
	if address == c.tokenEndpoint {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("User-Agent", "JornadaStudio")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if access != "" {
		req.Header.Set("Authorization", "Bearer "+access)
	}
	if c.Diagnostics != nil {
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), c.Diagnostics.HTTPTrace(ctx, started)))
	}
	// response jamais é impresso integralmente.
	response, err := c.http.Do(req)
	if err != nil {
		c.Diagnostics.Complete(ctx, "http.transport.failed", started, err)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("GitHub indisponível; verifique sua conexão e tente novamente.")
	}
	defer response.Body.Close()
	c.Diagnostics.Event(ctx, "http.response", syncdiag.HTTPFields(response))
	c.observeRate(response)
	if response.StatusCode < 200 || response.StatusCode > 299 {
		if address == c.tokenEndpoint {
			// failure ignora descrições livres e quaisquer tokens eventualmente presentes no corpo.
			var failure struct {
				Error string `json:"error"`
			}
			if json.NewDecoder(io.LimitReader(response.Body, 64<<10)).Decode(&failure) == nil {
				c.Diagnostics.Event(ctx, "oauth.provider.error", map[string]any{"provider_error": syncdiag.ProviderError(failure.Error)})
			}
		}
		switch response.StatusCode {
		case 401:
			return application.SyncProblem("reconnectRequired", "Autorização GitHub expirada ou revogada; reconecte a conta.")
		case 403:
			return application.SyncProblem("permissionDenied", "Acesso GitHub bloqueado por permissão, política, SSO ou limite; confira sua autorização.")
		case 404:
			return fmt.Errorf("Repositório GitHub inexistente ou inacessível.")
		case 429:
			return application.SyncProblem("rateLimited", "Limite GitHub atingido; aguarde antes de tentar novamente.")
		}
		return fmt.Errorf("GitHub recusou a operação (HTTP %d).", response.StatusCode)
	}
	if result == nil {
		return nil
	}
	if err = json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(result); err != nil {
		c.Diagnostics.Complete(ctx, "http.response.decode.failed", started, err)
		return fmt.Errorf("Resposta GitHub inválida.")
	}
	return nil
}

// Identity confirma a conta autenticada em cada autorização.
// ctx controla a consulta e token autoriza o acesso; identidade incompleta retorna erro.
func (c *Client) Identity(ctx context.Context, token application.GitHubCredential) (application.GitHubAccount, error) {
	// account não inclui nenhuma credencial.
	var account application.GitHubAccount
	// err conserva a falha HTTP antes de validar os campos públicos retornados.
	err := c.request(ctx, http.MethodGet, c.api+"/user", token.AccessToken, nil, &account)
	if err == nil && (account.ID <= 0 || account.Login == "") {
		err = fmt.Errorf("Identidade GitHub inválida.")
	}
	return account, err
}

// repositoryWire mantém o formato REST fora dos modelos de negócio.
type repositoryWire struct {
	// ID, Name e Owner identificam o destino.
	ID int64 `json:"id"`
	// Name é o nome físico atual do repositório no provedor.
	Name string `json:"name"`
	// Owner contém apenas a identificação pública do proprietário.
	Owner struct {
		// Login identifica o usuário ou organização proprietária.
		Login string `json:"login"`
	} `json:"owner"`
	// Branch e Visibility vêm do provedor, sem inferência por nome.
	Branch string `json:"default_branch"`
	// Visibility distingue público, privado e interno conforme o servidor.
	Visibility string `json:"visibility"`
	// Private suporta respostas antigas sem visibility.
	Private bool `json:"private"`
	// Archived, Disabled e Permissions indicam elegibilidade básica.
	Archived bool `json:"archived"`
	// Disabled registra indisponibilidade administrativa do destino.
	Disabled bool `json:"disabled"`
	// Permissions inclui somente a permissão usada para seleção básica.
	Permissions struct {
		// Push informa escrita básica, ainda sujeita às políticas da branch.
		Push bool `json:"push"`
	} `json:"permissions"`
}

// model converte o wire sem copiar URLs arbitrárias recebidas pela rede.
func (w repositoryWire) model() application.GitHubRepository {
	// visibility conserva compatibilidade com respostas que só incluem private.
	visibility := w.Visibility
	if visibility == "" {
		visibility = "public"
		if w.Private {
			visibility = "private"
		}
	}
	// branch define main apenas para repositório ainda sem referência informada.
	branch := w.Branch
	if branch == "" {
		branch = "main"
	}
	return application.GitHubRepository{ID: w.ID, Owner: w.Owner.Login, Name: w.Name, Branch: branch, Visibility: visibility, Writable: w.Permissions.Push, Archived: w.Archived, Disabled: w.Disabled}
}

// Repositories carrega todas as páginas antes de declarar a lista completa.
// ctx cancela a paginação e token autoriza cada página; falha não retorna lista parcial como completa.
func (c *Client) Repositories(ctx context.Context, token application.GitHubCredential) (items []application.GitHubRepository, listErr error) {
	ctx = application.WithSyncOperation(ctx)
	// started e pageNumber permitem localizar paginação interrompida sem expor nomes privados.
	started := time.Now()
	pageNumber := 0
	c.Diagnostics.Event(ctx, "repositories.start", nil)
	defer func() { c.Diagnostics.Complete(ctx, "repositories.complete", started, listErr) }()
	// next nasce no endpoint conhecido e somente aceita continuadores do mesmo endpoint.
	next := c.api + "/user/repos?visibility=all&affiliation=owner,collaborator,organization_member&per_page=100&sort=full_name"
	// seen evita loops de paginação e duplicação por ID.
	seen := map[string]bool{}
	// ids deduplica repositórios por sua identidade estável entre páginas.
	ids := map[int64]bool{}
	// repositories acumula somente metadados convertidos sem credenciais.
	repositories := []application.GitHubRepository{}
	for next != "" {
		pageNumber++
		c.Diagnostics.Event(ctx, "repositories.page.start", map[string]any{"page": pageNumber})
		if err := c.checkRate(); err != nil {
			return nil, err
		}
		if seen[next] {
			return nil, fmt.Errorf("Paginação GitHub inválida.")
		}
		seen[next] = true
		// req precisa do header Link além do corpo JSON.
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, next, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "JornadaStudio")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		if c.Diagnostics != nil {
			req = req.WithContext(httptrace.WithClientTrace(req.Context(), c.Diagnostics.HTTPTrace(ctx, time.Now())))
		}
		// response é encerrada em cada página, sem acumular conexões.
		response, err := c.http.Do(req)
		if err != nil {
			c.Diagnostics.Complete(ctx, "repositories.transport.failed", started, err)
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, fmt.Errorf("Não foi possível carregar os repositórios; tente novamente.")
		}
		c.observeRate(response)
		c.Diagnostics.Event(ctx, "repositories.page.response", syncdiag.HTTPFields(response))
		if response.StatusCode != 200 {
			response.Body.Close()
			return nil, fmt.Errorf("Não foi possível listar repositórios (HTTP %d); confira autorização, SSO e limites.", response.StatusCode)
		}
		// page é convertida antes de avançar a próxima URL.
		var page []repositoryWire
		err = json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(&page)
		response.Body.Close()
		if err != nil {
			c.Diagnostics.Complete(ctx, "repositories.decode.failed", started, err)
			return nil, fmt.Errorf("Listagem GitHub inválida.")
		}
		c.Diagnostics.Event(ctx, "repositories.page.decoded", map[string]any{"page": pageNumber, "items": len(page)})
		// wire é cada repositório da página, ainda no formato REST.
		for _, wire := range page {
			if wire.ID > 0 && !ids[wire.ID] {
				ids[wire.ID] = true
				repositories = append(repositories, wire.model())
			}
		}
		next = ""
		// part é cada relação Link a examinar para localizar o continuador next.
		for _, part := range strings.Split(response.Header.Get("Link"), ",") {
			if !strings.Contains(part, `rel="next"`) {
				continue
			}
			// start e end delimitam a URL entre os sinais de menor e maior no header Link.
			start, end := strings.Index(part, "<"), strings.Index(part, ">")
			if start < 0 || end <= start {
				return nil, fmt.Errorf("Paginação GitHub inválida.")
			}
			// address é o continuador recebido; err impede aceitar uma URL malformada.
			address, err := url.Parse(part[start+1 : end])
			if err != nil {
				return nil, fmt.Errorf("Paginação GitHub inválida.")
			}
			// base fixa o host/esquema autorizados para a paginação.
			base, _ := url.Parse(c.api)
			if address.Scheme != base.Scheme || address.Host != base.Host || address.Path != "/user/repos" || address.User != nil {
				return nil, fmt.Errorf("Destino de paginação GitHub recusado.")
			}
			next = address.String()
		}
	}
	c.Diagnostics.Event(ctx, "repositories.loaded", map[string]any{"total_items": len(repositories)})
	return repositories, nil
}

// route identifica apenas o contrato, nunca owner/repo, query, Client ID ou URL livre.
func (c *Client) route(address string) string {
	if address == c.tokenEndpoint {
		return "oauth.token"
	}
	if address == c.api+"/user" {
		return "user.get"
	}
	if strings.HasPrefix(address, c.api+"/user/repos?") {
		return "repositories.list"
	}
	if strings.HasPrefix(address, c.api+"/repos/") {
		return "repository.get"
	}
	return "other"
}

// checkRate exige aguardar o prazo informado sem bloquear a thread da interface.
func (c *Client) checkRate() error {
	c.rateMutex.Lock()
	defer c.rateMutex.Unlock()
	if time.Now().Before(c.retryAt) {
		return fmt.Errorf("Limite GitHub atingido; tente novamente após %s.", c.retryAt.Local().Format("15:04:05"))
	}
	return nil
}

// observeRate interpreta somente headers de limite, sem registrar a resposta do usuário.
func (c *Client) observeRate(response *http.Response) {
	if response.StatusCode != 429 && (response.StatusCode != 403 || response.Header.Get("X-RateLimit-Remaining") != "0") {
		return
	}
	// until mantém um intervalo conservador quando o servidor não informar prazo válido.
	until := time.Now().Add(time.Minute)
	if seconds, err := strconv.Atoi(response.Header.Get("Retry-After")); err == nil && seconds >= 0 && seconds <= 86400 {
		until = time.Now().Add(time.Duration(seconds) * time.Second)
	} else if date, err := http.ParseTime(response.Header.Get("Retry-After")); err == nil {
		until = date
	}
	if reset, err := strconv.ParseInt(response.Header.Get("X-RateLimit-Reset"), 10, 64); err == nil && time.Unix(reset, 0).After(until) {
		until = time.Unix(reset, 0)
	}
	c.rateMutex.Lock()
	defer c.rateMutex.Unlock()
	if until.After(c.retryAt) {
		c.retryAt = until
	}
}

// Repository revalida identidade, branch e permissão antes de cada operação remota.
// repo fornece o destino esperado e token autoriza a consulta sob ctx; ID diferente
// retorna erro, sem redirecionar o vínculo para outro repositório.
func (c *Client) Repository(ctx context.Context, token application.GitHubCredential, repo application.GitHubRepository) (application.GitHubRepository, error) {
	// wire é restrito a campos conhecidos e revalidado contra o ID selecionado.
	var wire repositoryWire
	// err conserva a falha HTTP antes de validar os campos públicos retornados.
	err := c.request(ctx, http.MethodGet, c.api+"/repos/"+url.PathEscape(repo.Owner)+"/"+url.PathEscape(repo.Name), token.AccessToken, nil, &wire)
	if err != nil {
		return repo, err
	}
	if wire.ID != repo.ID {
		return repo, fmt.Errorf("O destino GitHub mudou; selecione novamente o repositório.")
	}
	return wire.model(), nil
}
