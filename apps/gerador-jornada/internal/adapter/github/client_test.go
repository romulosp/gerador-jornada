package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestRepositoriesPaginationVisibility verifica completude e campos exigidos pelo seletor.
func TestRepositoriesPaginationVisibility(t *testing.T) {
	// server fornece duas páginas sem depender de conta ou rede externa.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w publica a página sintética e r permite conferir a autorização recebida.
		if r.Header.Get("Authorization") != "Bearer fixture" {
			t.Error("token ausente")
		}
		if r.URL.Query().Get("page") == "2" {
			fmt.Fprint(w, `[{"id":100,"name":"repetido","owner":{"login":"autor"}},{"id":101,"name":"privado","owner":{"login":"org"},"visibility":"private","permissions":{"push":true}}]`)
			return
		}
		w.Header().Set("Link", "<http://"+r.Host+`/user/repos?page=2>; rel="next"`)
		fmt.Fprint(w, "[")
		// id ocupa a primeira página inteira para detectar término prematuro após cem itens.
		for id := 1; id <= 100; id++ {
			if id > 1 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `{"id":%d,"name":"publico","owner":{"login":"autor"},"visibility":"public","permissions":{"push":false}}`, id)
		}
		fmt.Fprint(w, "]")
	}))
	defer server.Close()
	// client aponta exclusivamente ao servidor local.
	client := New("id", "secret")
	client.api = server.URL
	// repositories e err permitem conferir completude e visibilidade após percorrer ambas as páginas.
	repositories, err := client.Repositories(context.Background(), fixtureToken())
	if err != nil || len(repositories) != 101 || repositories[100].Visibility != "private" || !repositories[100].Writable {
		t.Fatalf("resultado inesperado: %+v %v", repositories, err)
	}
}

// TestAuthorizePKCEReturn exercita callback real loopback sem abrir navegador externo.
func TestAuthorizePKCEReturn(t *testing.T) {
	// server simula a resposta da troca de código sem acessar o provedor real.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w devolve o token sintético; r representa a chamada local que substitui o provedor.
		if r.Header.Get("Accept") != "application/json" {
			t.Error("troca OAuth não solicitou resposta JSON")
		}
		fmt.Fprint(w, `{"access_token":"fixture","scope":"repo","expires_in":3600}`)
	}))
	defer server.Close()
	// client conserva callback loopback real e troca token só no fixture.
	client := New("id", "secret")
	client.tokenEndpoint = server.URL
	// token e err permitem conferir o retorno após o navegador simulado concluir o callback.
	token, err := client.Authorize(context.Background(), func(address *url.URL) error {
		// address é a URL de autorização entregue ao navegador simulado.
		if address.Query().Get("code_challenge_method") != "S256" {
			t.Error("PKCE ausente")
		}
		// callback inclui exclusivamente o state correto da tentativa.
		callback, _ := url.Parse(address.Query().Get("redirect_uri"))
		// query reproduz o state recebido e acrescenta um código sintético de uso único.
		query := callback.Query()
		query.Set("state", address.Query().Get("state"))
		query.Set("code", "one-use")
		callback.RawQuery = query.Encode()
		// response e err confirmam o envio ao listener real e permitem liberar o corpo HTTP.
		response, err := http.Get(callback.String())
		if err == nil {
			response.Body.Close()
		}
		return err
	})
	if err != nil || token.AccessToken != "fixture" {
		t.Fatalf("autorização falhou: %v", err)
	}
}
