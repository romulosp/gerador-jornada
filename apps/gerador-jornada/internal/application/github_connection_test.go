package application

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"testing"
	"time"
)

// syncFixture implementa fronteiras isoladas sem rede, disco ou cofre reais.
type syncFixture struct {
	// account e token representam respostas sintéticas do provedor.
	account GitHubAccount
	token   GitHubCredential
	// records e credentials mantêm estado em memória por namespace.
	records     map[string][]byte
	credentials map[string]GitHubCredential
	// denyDelete e refreshCount permitem induzir falha e conferir renovação.
	denyDelete   bool
	refreshCount int
}

// wrongRepositoryDirectory simula provedor que respondeu outra identidade ao destino selecionado.
type wrongRepositoryDirectory struct{ *syncFixture }

// Repository altera só o ID para provar que nome e permissão não bastam para liberar operação.
func (f wrongRepositoryDirectory) Repository(_ context.Context, _ GitHubCredential, repo GitHubRepository) (GitHubRepository, error) {
	repo.ID++
	return repo, nil
}

// TestConnectionRejectsAnotherRepositoryID impede publicação em outro destino com mesmo nome aparente.
func TestConnectionRejectsAnotherRepositoryID(t *testing.T) {
	// connection conserva cofre/preferências isolados, substituindo só a resposta de diretório.
	connection, fixture := connectionFixture()
	if _, err := connection.Connect(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	connection.directory = wrongRepositoryDirectory{fixture}
	called := false
	err := connection.WithRepository(context.Background(), GitHubRepository{ID: 42, Owner: "fixture", Name: "destino", Writable: true}, func(GitHubCredential, GitHubAccount, GitHubRepository) error { called = true; return nil })
	if err == nil || called {
		t.Fatal("operação liberada para identidade diferente")
	}
}

// Authorize devolve token de teste sem abrir navegador.
func (f *syncFixture) Authorize(context.Context, func(*url.URL) error) (GitHubCredential, error) {
	return f.token, nil
}

// Refresh comprova que uma renovação não publica nenhuma alteração.
func (f *syncFixture) Refresh(context.Context, GitHubCredential) (GitHubCredential, error) {
	f.refreshCount++
	return GitHubCredential{AccessToken: "renewed", Scopes: "repo", ExpiresAt: time.Now().Add(time.Hour)}, nil
}

// Identity representa a conta atualmente autorizada pelo fixture.
func (f *syncFixture) Identity(context.Context, GitHubCredential) (GitHubAccount, error) {
	return f.account, nil
}

// Repositories produz uma lista vazia válida sem chamadas externas.
func (f *syncFixture) Repositories(context.Context, GitHubCredential) ([]GitHubRepository, error) {
	return []GitHubRepository{}, nil
}

// Repository preserva o destino selecionado durante a revalidação.
func (f *syncFixture) Repository(_ context.Context, _ GitHubCredential, repo GitHubRepository) (GitHubRepository, error) {
	return repo, nil
}

// Read preserva semântica de ausência de uma preferência.
func (f *syncFixture) Read(key string, destination any) error {
	data, ok := f.records[key]
	if !ok {
		return os.ErrNotExist
	}
	return json.Unmarshal(data, destination)
}

// Write guarda somente o JSON recebido, permitindo inspecionar vazamento de token.
func (f *syncFixture) Write(_ context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err == nil {
		f.records[key] = data
	}
	return err
}

// Load lê uma credencial de namespace isolado.
func (f *syncFixture) Load(key string) (GitHubCredential, error) {
	token, ok := f.credentials[key]
	if !ok {
		return token, os.ErrNotExist
	}
	return token, nil
}

// Save substitui somente a credencial explicitamente indicada.
func (f *syncFixture) Save(key string, token GitHubCredential) error {
	f.credentials[key] = token
	return nil
}

// Delete simula recusa do cofre sem apagar dados.
func (f *syncFixture) Delete(key string) error {
	if f.denyDelete {
		return errors.New("cofre recusou remoção")
	}
	delete(f.credentials, key)
	return nil
}

// connectionFixture monta todas as fronteiras sem efeitos pessoais no ambiente.
func connectionFixture() (*GitHubConnection, *syncFixture) {
	// fixture usa IDs e tokens deliberadamente sintéticos.
	fixture := &syncFixture{account: GitHubAccount{ID: 1, Login: "fixture"}, token: GitHubCredential{AccessToken: "synthetic", Scopes: "repo"}, records: map[string][]byte{}, credentials: map[string]GitHubCredential{}}
	return NewGitHubConnection(fixture, fixture, fixture, fixture, "test"), fixture
}

// TestConnectionRefreshDisconnect confirma renovação, remoção e metadados sem token.
func TestConnectionRefreshDisconnect(t *testing.T) {
	// connection e fixture compartilham somente memória deste teste.
	connection, fixture := connectionFixture()
	if _, err := connection.Connect(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	// preference permite localizar a credencial sem copiá-la para os dados comuns.
	var preference GitHubAccountPreference
	if err := fixture.Read("github.json", &preference); err != nil {
		t.Fatal(err)
	}
	token := fixture.credentials[preference.CredentialKey]
	token.ExpiresAt = time.Now().Add(-time.Minute)
	token.RefreshToken = "refresh"
	fixture.credentials[preference.CredentialKey] = token
	if _, err := connection.Repositories(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fixture.refreshCount != 1 {
		t.Fatal("refresh ausente")
	}
	if err := connection.Disconnect(context.Background()); err != nil {
		t.Fatal(err)
	}
	account, err := connection.Account()
	if err != nil || account.ID != 0 || len(fixture.credentials) != 0 {
		t.Fatal("desconexão incompleta")
	}
}

// TestConnectionRequiresRepoScope recusa autorização incompleta antes de persistir credenciais.
func TestConnectionRequiresRepoScope(t *testing.T) {
	connection, fixture := connectionFixture()
	fixture.token.Scopes = "read:user"
	if _, err := connection.Connect(context.Background(), nil); err == nil {
		t.Fatal("scope insuficiente aceito")
	}
	if len(fixture.credentials) != 0 || len(fixture.records) != 0 {
		t.Fatal("autorização incompleta persistida")
	}
}

// TestConnectionRetainsAccountWhenRemovalFails impede trocar a conta ativa deixando segredo anterior abandonado.
func TestConnectionRetainsAccountWhenRemovalFails(t *testing.T) {
	connection, fixture := connectionFixture()
	if _, err := connection.Connect(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	fixture.account = GitHubAccount{ID: 2, Login: "second"}
	fixture.denyDelete = true
	if _, err := connection.Connect(context.Background(), nil); err == nil {
		t.Fatal("falha de exclusão ignorada")
	}
	account, err := connection.Account()
	if err != nil || account.ID != 1 {
		t.Fatal("conta substituída apesar da falha de remoção")
	}
}
