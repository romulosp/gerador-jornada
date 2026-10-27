package application

import (
	"context"
	"testing"
)

// syncBackendFixture registra apenas se uma confirmação chegou ao mecanismo de publicação.
type syncBackendFixture struct {
	// repository representa o vínculo selecionado pelo usuário.
	repository GitHubRepository
	// published detecta envio indevido após troca de conta.
	published bool
}

// Binding devolve uma preferência sintética sem acessar arquivos.
func (f *syncBackendFixture) Binding(root string) (GitWorkspaceBinding, error) {
	return GitWorkspaceBinding{Version: 1, Root: root, Repository: f.repository}, nil
}

// Recover não produz efeitos em um fixture sem journal.
func (f *syncBackendFixture) Recover(context.Context, string) error { return nil }

// Review fixa o destino sem conhecer a credencial do usuário.
func (f *syncBackendFixture) Review(_ context.Context, root string, _ GitHubCredential) (GitPublicationReview, error) {
	return GitPublicationReview{Root: root, Repository: f.repository}, nil
}

// Bind somente atende à porta mutante exigida pelo coordenador.
func (f *syncBackendFixture) Bind(context.Context, string, GitHubRepository, GitHubCredential) error {
	return nil
}

// Publish marca o efeito que precisa ser bloqueado pelo teste.
func (f *syncBackendFixture) Publish(context.Context, GitPublicationReview, string, GitHubAccount, GitHubCredential) error {
	f.published = true
	return nil
}

// TestPublicationRejectsAccountChange exige confirmar novamente quando outra conta se torna ativa.
func TestPublicationRejectsAccountChange(t *testing.T) {
	// connection usa somente preferências e cofre em memória.
	connection, provider := connectionFixture()
	if _, err := connection.Connect(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	backend := &syncBackendFixture{repository: GitHubRepository{ID: 3, Owner: "fixture", Name: "project", Branch: "main", Writable: true}}
	service := NewWorkspaceSync(connection, backend, backend)
	review, err := service.ReviewPublication(context.Background(), "fixture")
	if err != nil || review.AccountID != 1 {
		t.Fatal("conta da revisão ausente")
	}
	provider.account = GitHubAccount{ID: 2, Login: "other"}
	if _, err = connection.Connect(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if err = service.Publish(context.Background(), review, "mensagem"); err == nil || backend.published {
		t.Fatal("publicação usou outra conta sem nova revisão")
	}
}
