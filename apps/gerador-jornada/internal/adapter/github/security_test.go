package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
)

// TestInvalidCallbackAndReplay exige state correto e uma única troca mesmo com retorno repetido.
func TestInvalidCallbackAndReplay(t *testing.T) {
	// calls só conta trocas de código sintéticas, nunca tokens de usuário.
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		fmt.Fprint(w, `{"access_token":"fixture","scope":"repo"}`)
	}))
	defer server.Close()
	client := New("id", "secret")
	client.tokenEndpoint = server.URL
	_, err := client.Authorize(context.Background(), func(address *url.URL) error {
		// callback é recebido no listener real da aplicação e não no servidor de tokens.
		callback, _ := url.Parse(address.Query().Get("redirect_uri"))
		query := callback.Query()
		query.Set("state", "invalid")
		query.Set("code", "fixture")
		callback.RawQuery = query.Encode()
		response, err := http.Get(callback.String())
		if err != nil {
			return err
		}
		response.Body.Close()
		if response.StatusCode != 400 {
			t.Error("state inválido aceito")
		}
		query.Set("state", address.Query().Get("state"))
		callback.RawQuery = query.Encode()
		for attempt := 0; attempt < 2; attempt++ {
			response, err = http.Get(callback.String())
			if err != nil {
				return err
			}
			response.Body.Close()
			if response.StatusCode != 200 {
				t.Error("retorno válido recusado")
			}
		}
		return nil
	})
	if err != nil || calls.Load() != 1 {
		t.Fatalf("retorno/replay incorreto: %v", err)
	}
}

// TestAuthorizationCancellation invalida a sessão sem depender do prazo de dez minutos.
func TestAuthorizationCancellation(t *testing.T) {
	// ctx permite cancelar a tentativa enquanto o navegador seria aberto.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := New("id", "secret")
	if _, err := client.Authorize(ctx, func(*url.URL) error { cancel(); return nil }); err == nil {
		t.Fatal("cancelamento ignorado")
	}
}

// TestPaginationCannotSendTokenToAnotherHost bloqueia Link externo antes de qualquer chamada.
func TestPaginationCannotSendTokenToAnotherHost(t *testing.T) {
	// leaked representa qualquer acesso indevido ao segundo servidor, mesmo sem analisar token.
	var leaked atomic.Int32
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1); fmt.Fprint(w, "[]") }))
	defer foreign.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", "<"+foreign.URL+`/user/repos>; rel="next"`)
		fmt.Fprint(w, "[]")
	}))
	defer server.Close()
	client := New("id", "secret")
	client.api = server.URL
	if _, err := client.Repositories(context.Background(), fixtureToken()); err == nil || leaked.Load() != 0 {
		t.Fatal("paginação externa não foi bloqueada")
	}
}

// TestRateLimitPreventsImmediateRetry preserva Retry-After sem sleep ou loop de rede.
func TestRateLimitPreventsImmediateRetry(t *testing.T) {
	// calls deve continuar em um mesmo quando o chamador clicar novamente.
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(429)
	}))
	defer server.Close()
	client := New("id", "secret")
	client.api = server.URL
	if _, err := client.Identity(context.Background(), fixtureToken()); err == nil {
		t.Fatal("429 ignorado")
	}
	if _, err := client.Identity(context.Background(), fixtureToken()); err == nil || calls.Load() != 1 {
		t.Fatal("limite não respeitado")
	}
}
