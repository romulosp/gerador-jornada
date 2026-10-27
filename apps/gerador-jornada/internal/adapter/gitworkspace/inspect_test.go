package gitworkspace

import (
	"context"
	"gerador-jornada.local/editor/internal/application"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"path/filepath"
	"testing"
)

// TestUnsafeRepositoryConfigurations impede inicializar por cima de metadados ou operações alheias.
func TestUnsafeRepositoryConfigurations(t *testing.T) {
	for _, kind := range []string{"bare", "linked", "index-lock", "branch", "origin", "ancestor"} {
		t.Run(kind, func(t *testing.T) {
			// engine/root são descartáveis mesmo quando a configuração Git é deliberadamente inválida.
			engine, root, _, destination := fixture(t)
			switch kind {
			case "bare":
				if _, err := git.PlainInit(root, true); err != nil {
					t.Fatal(err)
				}
			case "linked":
				put(t, root, ".git", "gitdir: missing")
			default:
				if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
					t.Fatal(err)
				}
				repo, err := git.PlainOpen(root)
				if err != nil {
					t.Fatal(err)
				}
				switch kind {
				case "index-lock":
					put(t, filepath.Join(root, ".git"), "index.lock", "external")
				case "branch":
					if err = repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("other"))); err != nil {
						t.Fatal(err)
					}
				case "origin":
					configuration, err := repo.Config()
					if err != nil {
						t.Fatal(err)
					}
					configuration.Remotes["origin"].URLs = []string{"https://github.com/other/repo.git"}
					if err = repo.SetConfig(configuration); err != nil {
						t.Fatal(err)
					}
				case "ancestor":
					root = filepath.Join(root, "nested")
					if err = os.Mkdir(root, 0700); err != nil {
						t.Fatal(err)
					}
				}
			}
			if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err == nil {
				t.Fatal("configuração insegura aceita")
			}
		})
	}
}

// TestRecoveryIsOptionalForLocalNestedProject não bloqueia salvar projeto local sem journal por causa do Git ancestral.
func TestRecoveryIsOptionalForLocalNestedProject(t *testing.T) {
	// engine recebe uma raiz sem vínculo dentro de outro repositório, sem executar integração.
	engine, parent, _, destination := fixture(t)
	if err := engine.Bind(context.Background(), parent, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(parent, "local-only")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := engine.Recover(context.Background(), root); err != nil {
		t.Fatal("edição local bloqueada por integração opcional")
	}
}

// TestLocalStatusDoesNotUseNetwork distingue jornadas alteradas sem autenticação ou transporte.
func TestLocalStatusDoesNotUseNetwork(t *testing.T) {
	// engine começa com vínculo válido e um arquivo ainda não publicado.
	engine, root, _, destination := fixture(t)
	put(t, root, "A.jornada", "local")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	engine.remoteURL = func(application.GitHubRepository) string { t.Fatal("status local tentou consultar remoto"); return "" }
	changes, err := engine.LocalChanges(root)
	if err != nil || len(changes) != 1 || changes[0].Kind != "inclusão" {
		t.Fatal("alterações locais não identificadas")
	}
}
