package github

import (
	"context"
	"encoding/json"
	"fmt"
	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// TestOAuthDiagnosticsExcludeCredentials acompanha falha real loopback sem registrar nenhum segredo da tentativa.
func TestOAuthDiagnosticsExcludeCredentials(t *testing.T) {
	// mutex protege os registros emitidos pelas goroutines do listener e do cliente HTTP.
	var mutex sync.Mutex
	records := []string{}
	secrets := []string{"fixture-confidential-client-id", "fixture-confidential-client-secret", "fixture-confidential-token", "fixture-confidential-code"}
	observer := syncdiag.New(func(event string, fields map[string]any) {
		mutex.Lock()
		defer mutex.Unlock()
		data, _ := json.Marshal(fields)
		records = append(records, event+" "+string(data))
	})
	// server devolve descrição maliciosa para comprovar que corpos livres não chegam ao log.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameters := map[string]string{}
		_ = json.NewDecoder(r.Body).Decode(&parameters)
		mutex.Lock()
		secrets = append(secrets, parameters["code_verifier"])
		mutex.Unlock()
		w.Header().Set("X-GitHub-Request-Id", "ABCD:1234:5678")
		fmt.Fprint(w, `{"error":"bad_verification_code","error_description":"fixture-confidential-client-secret","access_token":"fixture-confidential-token"}`)
	}))
	defer server.Close()
	client := New(secrets[0], secrets[1])
	client.tokenEndpoint = server.URL
	client.Diagnostics = observer
	_, err := client.Authorize(context.Background(), func(address *url.URL) error {
		mutex.Lock()
		secrets = append(secrets, address.Query().Get("state"), address.Query().Get("code_challenge"))
		mutex.Unlock()
		callback, _ := url.Parse(address.Query().Get("redirect_uri"))
		query := callback.Query()
		query.Set("state", address.Query().Get("state"))
		query.Set("code", "fixture-confidential-code")
		callback.RawQuery = query.Encode()
		response, err := http.Get(callback.String())
		if err == nil {
			response.Body.Close()
		}
		return err
	})
	if err == nil {
		t.Fatal("falha OAuth esperada")
	}
	mutex.Lock()
	defer mutex.Unlock()
	joined := strings.Join(records, "\n")
	for _, secret := range secrets {
		if secret != "" && strings.Contains(joined, secret) {
			t.Fatal("dado confidencial encontrado no diagnóstico")
		}
	}
	for _, event := range []string{"oauth.listener.ready", "oauth.browser.open.complete", "oauth.callback.received", "oauth.callback.validated", "oauth.exchange.start", "oauth.provider.error", "oauth.authorize.complete"} {
		if !strings.Contains(joined, event) {
			t.Fatalf("etapa ausente: %s", event)
		}
	}
	if !strings.Contains(joined, "bad_verification_code") || !strings.Contains(joined, "ABCD:1234:5678") {
		t.Fatal("diagnóstico seguro do provedor ausente")
	}
}
