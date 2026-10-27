package github

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/application"
)

// oauthResponse é interpretada sem expor descrições livres do servidor.
type oauthResponse struct {
	// AccessToken e RefreshToken são persistidos somente no cofre.
	AccessToken string `json:"access_token"`
	// RefreshToken permite renovar quando o cadastro oferecer esse recurso.
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn está em segundos e pode ser omitido pelo cadastro.
	ExpiresIn int64 `json:"expires_in"`
	// Scope registra o alcance efetivamente informado pelo provedor.
	Scope string `json:"scope"`
	// Error informa uma falha de autorização sem transportar corpo bruto à UI.
	Error string `json:"error"`
}

// secret gera entropia independente para state e verifier.
func secret() string {
	// value tem 256 bits de entropia do sistema operacional.
	value := make([]byte, 32)
	// err indica indisponibilidade do gerador criptográfico do sistema.
	if _, err := rand.Read(value); err != nil {
		panic("Gerador criptográfico indisponível")
	}
	return base64.RawURLEncoding.EncodeToString(value)
}

// exchange executa a troca no desktop e converte o prazo sem imprimir tokens.
// ctx controla I/O e parameters recebe as credenciais do cadastro para essa troca.
// Retorna o token ou uma falha sanitizada, sem expor corpos livres do provedor.
func (c *Client) exchange(ctx context.Context, parameters map[string]string) (credential application.GitHubCredential, exchangeErr error) {
	ctx = application.WithSyncOperation(ctx)
	// started marca o início da troca sem registrar código, verifier ou credenciais do cadastro.
	started := time.Now()
	c.Diagnostics.Event(ctx, "oauth.exchange.start", map[string]any{"refresh_available": parameters["grant_type"] == "refresh_token"})
	defer func() { c.Diagnostics.Complete(ctx, "oauth.exchange.complete", started, exchangeErr) }()
	// response contém apenas o contrato OAuth esperado.
	var response oauthResponse
	parameters["client_id"], parameters["client_secret"] = c.id, c.secret
	// err conserva a falha de transporte sem registrar o conteúdo da troca.
	if err := c.request(ctx, http.MethodPost, c.tokenEndpoint, "", parameters, &response); err != nil {
		return application.GitHubCredential{}, err
	}
	if response.Error != "" || response.AccessToken == "" {
		c.Diagnostics.Event(ctx, "oauth.provider.error", map[string]any{"provider_error": syncdiag.ProviderError(response.Error), "success": false})
		return application.GitHubCredential{}, fmt.Errorf("Autorização GitHub não concluída; conecte novamente.")
	}
	// token permanece no backend até ser protegido no cofre.
	token := application.GitHubCredential{AccessToken: response.AccessToken, RefreshToken: response.RefreshToken, Scopes: response.Scope}
	if response.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}
	// repoScope é um indicador de permissão, não o token ou a lista livre de atributos.
	repoScope := false
	for _, scope := range strings.Fields(strings.ReplaceAll(response.Scope, ",", " ")) {
		if scope == "repo" {
			repoScope = true
		}
	}
	c.Diagnostics.Event(ctx, "oauth.token.received", map[string]any{"success": true, "repo_scope": repoScope, "refresh_available": response.RefreshToken != "", "expires_in_seconds": response.ExpiresIn})
	return token, nil
}

// Authorize abre o navegador e espera um único callback válido da tentativa atual.
// parent permite cancelar e open recebe a URL pública com desafio PKCE.
// Retorna credencial após o callback ou falha de cadastro, navegador, rede ou sessão.
func (c *Client) Authorize(parent context.Context, open func(*url.URL) error) (credential application.GitHubCredential, authorizeErr error) {
	parent = application.WithSyncOperation(parent)
	// started permite reconhecer timeout quando nenhum callback tiver chegado.
	started := time.Now()
	c.Diagnostics.Event(parent, "oauth.authorize.start", map[string]any{"oauth_configured": c.id != "" && c.secret != ""})
	defer func() { c.Diagnostics.Complete(parent, "oauth.authorize.complete", started, authorizeErr) }()
	if c.id == "" || c.secret == "" {
		return application.GitHubCredential{}, fmt.Errorf("Este build ainda não tem cadastro OAuth GitHub configurado.")
	}
	// ctx limita a sessão e permite cancelamento ao fechar o aplicativo.
	ctx, cancel := context.WithTimeout(parent, 10*time.Minute)
	defer cancel()
	// listener é aberto antes do navegador para não perder o retorno.
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return application.GitHubCredential{}, fmt.Errorf("Não foi possível receber o retorno local do GitHub.")
	}
	defer listener.Close()
	// port é informação técnica do loopback; nenhum endereço externo é copiado ao log.
	_, portText, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portText)
	c.Diagnostics.Event(ctx, "oauth.listener.ready", map[string]any{"callback_host": "127.0.0.1", "callback_port": port, "callback_path": "/oauth/github/callback"})
	// state e verifier nunca são persistidos ou registrados em log.
	state, verifier := secret(), secret()
	// challenge é o resumo S256 do verifier enviado ao provedor, sem revelar sua origem.
	challenge := sha256.Sum256([]byte(verifier))
	// redirect fixa o endereço efêmero onde esta tentativa receberá seu retorno.
	redirect := "http://" + listener.Addr().String() + "/oauth/github/callback"
	// outcome transporta somente o código temporário ou sinal de negação.
	outcome := make(chan string, 1)
	// once impede que callbacks repetidos iniciem mais de uma troca de código.
	once := sync.Once{}
	// server escuta apenas no loopback e não serve outros recursos.
	server := &http.Server{ReadHeaderTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w responde sem cache e r contém somente o retorno a validar contra a tentativa.
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		// query é inspecionada somente para flags, contagens e validação em tempo constante.
		query := r.URL.Query()
		methodValid, hostValid, pathValid := r.Method == "GET", r.Host == listener.Addr().String(), r.URL.Path == "/oauth/github/callback"
		stateValid := subtle.ConstantTimeCompare([]byte(query.Get("state")), []byte(state)) == 1
		parametersValid := len(query["state"]) == 1 && len(query["code"]) <= 1 && len(query["error"]) <= 1
		c.Diagnostics.Event(ctx, "oauth.callback.received", map[string]any{"method_valid": methodValid, "host_valid": hostValid, "path_valid": pathValid, "state_present": query.Get("state") != "", "state_valid": stateValid, "code_present": query.Get("code") != "", "provider_error": syncdiag.ProviderError(query.Get("error")), "query_parameter_count": len(query)})
		if !methodValid || !hostValid || !pathValid || !stateValid || !parametersValid {
			// reason é uma enumeração local e não inclui parâmetros recebidos do navegador.
			reason := "parameters"
			if !methodValid {
				reason = "method"
			} else if !hostValid {
				reason = "host"
			} else if !pathValid {
				reason = "path"
			} else if !stateValid {
				reason = "state"
			}
			c.Diagnostics.Event(ctx, "oauth.callback.rejected", map[string]any{"reason": reason, "success": false})
			http.Error(w, "Retorno de autorização inválido.", 400)
			return
		}
		// code vazio representa autorização negada, somente com state correto.
		code := r.URL.Query().Get("code")
		// accepted distingue repetição de callback sem revelar ou hashear o código.
		accepted := false
		once.Do(func() { accepted = true; outcome <- code })
		c.Diagnostics.Event(ctx, "oauth.callback.validated", map[string]any{"success": accepted, "replay": !accepted})
		fmt.Fprint(w, "Autorização recebida. Volte ao Jornada Studio.")
	})}
	defer server.Close()
	go func() {
		if serveErr := server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			c.Diagnostics.Complete(ctx, "oauth.listener.failed", started, serveErr)
		}
	}()
	// address contém somente parâmetros públicos e desafio PKCE.
	address, err := url.Parse(c.authorizeEndpoint)
	if err != nil {
		return application.GitHubCredential{}, fmt.Errorf("Cadastro OAuth inválido.")
	}
	// query reúne apenas parâmetros públicos e o desafio da tentativa atual.
	query := address.Query()
	query.Set("client_id", c.id)
	query.Set("redirect_uri", redirect)
	query.Set("scope", "repo offline_access")
	query.Set("state", state)
	query.Set("code_challenge", base64.RawURLEncoding.EncodeToString(challenge[:]))
	query.Set("code_challenge_method", "S256")
	address.RawQuery = query.Encode()
	c.Diagnostics.Event(ctx, "oauth.browser.open.start", nil)
	if err = open(address); err != nil {
		c.Diagnostics.Complete(ctx, "oauth.browser.open.complete", started, err)
		return application.GitHubCredential{}, fmt.Errorf("Não foi possível abrir o navegador padrão.")
	}
	c.Diagnostics.Event(ctx, "oauth.browser.open.complete", map[string]any{"success": true})
	select {
	case <-ctx.Done():
		return application.GitHubCredential{}, ctx.Err()
	// code é o resultado único do callback aceito para esta tentativa.
	case code := <-outcome:
		if ctx.Err() != nil {
			return application.GitHubCredential{}, ctx.Err()
		}
		if code == "" {
			return application.GitHubCredential{}, fmt.Errorf("Autorização GitHub negada ou cancelada.")
		}
		return c.exchange(ctx, map[string]string{"code": code, "redirect_uri": redirect, "code_verifier": verifier})
	}
}

// Refresh renova apenas a autorização e nunca inicia publicação.
// ctx controla a troca e token fornece refreshToken; sua ausência exige reconexão.
func (c *Client) Refresh(ctx context.Context, token application.GitHubCredential) (application.GitHubCredential, error) {
	if token.RefreshToken == "" {
		return token, fmt.Errorf("Autorização expirada; reconecte o GitHub.")
	}
	// refreshed pode omitir escopos que permanecem iguais aos da autorização original.
	refreshed, err := c.exchange(ctx, map[string]string{"grant_type": "refresh_token", "refresh_token": token.RefreshToken})
	if err == nil && refreshed.Scopes == "" {
		refreshed.Scopes = token.Scopes
	}
	return refreshed, err
}
