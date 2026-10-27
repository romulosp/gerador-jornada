//go:build windows

package credential

import (
	"gerador-jornada.local/editor/internal/adapter/userconfig"
	"gerador-jornada.local/editor/internal/application"
	"testing"
)

// TestWindowsCredentialRoundTrip testa o cofre real em namespace descartável sem tocar credenciais pessoais.
func TestWindowsCredentialRoundTrip(t *testing.T) {
	// key é diferente para cada execução, mesmo com testes simultâneos.
	key := "JornadaStudio/test/" + userconfig.Key(t.TempDir())
	store := Store{}
	if err := store.Save(key, application.GitHubCredential{AccessToken: "synthetic-only", Scopes: "repo"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := store.Delete(key); err != nil {
			t.Error(err)
		}
	})
	token, err := store.Load(key)
	if err != nil || token.AccessToken != "synthetic-only" {
		t.Fatal("cofre não preservou credencial sintética")
	}
	if err = store.Delete(key); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Load(key); err == nil {
		t.Fatal("credencial permaneceu após desconectar")
	}
}
