package gitworkspace

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/adapter/syncdiag"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/synctest"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// seedJourneyRemote grava árvore remota real sem checkout nem dependência de git.exe.
// files é uma fixture sintética e o retorno identifica seu commit na branch main.
func seedJourneyRemote(t *testing.T, remote string, files map[string]entry) string {
	t.Helper()
	// repo e tree pertencem exclusivamente ao repositório temporário do teste.
	repo, err := git.PlainOpen(remote)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := writeTree(repo, files)
	if err != nil {
		t.Fatal(err)
	}
	// signature e commit simulam um cliente remoto com conteúdo arbitrário já publicado.
	signature := object.Signature{Name: "Fixture", Email: "fixture@example.invalid", When: time.Now()}
	commit := &object.Commit{TreeHash: tree, Author: signature, Committer: signature, Message: "Base remota"}
	// previous preserva o histórico quando o teste acrescenta uma atualização remota.
	previous, err := head(repo)
	if err != nil {
		t.Fatal(err)
	}
	if previous != "" {
		commit.ParentHashes = []plumbing.Hash{plumbing.NewHash(previous)}
	}
	// encoded recebe o objeto antes da atualização da referência da fixture.
	encoded := repo.Storer.NewEncodedObject()
	if err = commit.Encode(encoded); err != nil {
		t.Fatal(err)
	}
	oid, err := repo.Storer.SetEncodedObject(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if err = repo.Storer.SetReference(plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), oid)); err != nil {
		t.Fatal(err)
	}
	if err = repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main"))); err != nil {
		t.Fatal(err)
	}
	return oid.String()
}

// TestBindingMaterializesOnlyRootJourneys reproduz o workspace poluído pelo código-fonte remoto.
func TestBindingMaterializesOnlyRootJourneys(t *testing.T) {
	// engine, root e remote isolam os efeitos da integração e da publicação real de objetos Git.
	engine, root, remote, destination := fixture(t)
	// files mistura fonte, preferências e jornada aninhada que não podem surgir no workspace.
	files := map[string]entry{}
	for _, name := range []string{"README.md", ".agents/SKILL.md", ".github/workflows/build.yml", ".specify/spec.md", ".vscode/settings.json", "apps/main.go", "specs/spec.md", "folder/Extra.jornada"} {
		files[name] = entry{Data: []byte("conteúdo remoto: " + name), Mode: filemode.Regular}
	}
	files["Remota.jornada"] = entry{Data: []byte("jornada remota"), Mode: filemode.Regular}
	seedJourneyRemote(t, remote, files)
	synctest.Journeys(t, root)
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// entries comprova ausência de cópias extras, em vez de só inspecionar a lista revisável.
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range entries {
		if item.Name() != ".git" && item.Name() != "Principal.jornada" && item.Name() != "Destino.jornada" && item.Name() != "Remota.jornada" {
			t.Fatalf("arquivo remoto alheio materializado: %s", item.Name())
		}
	}
	if err = publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// remoteRepo e delivered conferem o conteúdo efetivamente publicado, incluindo o histórico preservado.
	remoteRepo, err := git.PlainOpen(remote)
	if err != nil {
		t.Fatal(err)
	}
	oid, err := head(remoteRepo)
	if err != nil {
		t.Fatal(err)
	}
	delivered, err := loadTree(remoteRepo, oid)
	if err != nil {
		t.Fatal(err)
	}
	for name, expected := range files {
		if actual, ok := delivered[name]; !equal(expected, true, actual, ok) {
			t.Fatalf("conteúdo remoto anterior alterado: %s", name)
		}
	}
	if len(delivered) != len(files)+2 {
		t.Fatal("publicação não entregou somente as duas jornadas novas")
	}
}

// TestNonJourneyCollisionDoesNotBlockOrOverwrite conserva arquivos locais alheios ao pacote.
func TestNonJourneyCollisionDoesNotBlockOrOverwrite(t *testing.T) {
	// engine e diretórios são privados; o README local nunca pertence à publicação do Studio.
	engine, root, remote, destination := fixture(t)
	files := map[string]entry{"README.md": {Data: []byte("remoto inicial"), Mode: filemode.Regular}}
	seedJourneyRemote(t, remote, files)
	put(t, root, "README.md", "arquivo local preservado")
	put(t, root, "A.jornada", "jornada local")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	files["README.md"] = entry{Data: []byte("remoto atualizado"), Mode: filemode.Regular}
	seedJourneyRemote(t, remote, files)
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// data verifica a preservação do arquivo mesmo com atualização remota durante publicação.
	data, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil || string(data) != "arquivo local preservado" {
		t.Fatal("arquivo local alheio substituído")
	}
}

// TestEmptyRemoteDoesNotConfirmCachedReference impede sucesso baseado somente em origin/main antigo.
func TestEmptyRemoteDoesNotConfirmCachedReference(t *testing.T) {
	// engine começa com um envio confirmado e portanto uma referência remota em cache.
	engine, root, remote, destination := fixture(t)
	put(t, root, "A.jornada", "conteúdo")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// remoteRepo simula branch removida pelo servidor depois da publicação anterior.
	remoteRepo, err := git.PlainOpen(remote)
	if err != nil {
		t.Fatal(err)
	}
	if err = remoteRepo.Storer.RemoveReference(plumbing.NewBranchReferenceName("main")); err != nil {
		t.Fatal(err)
	}
	localRepo, err := git.PlainOpen(root)
	if err != nil {
		t.Fatal(err)
	}
	binding, err := engine.Binding(root)
	if err != nil {
		t.Fatal(err)
	}
	oid, err := fetch(context.Background(), localRepo, binding, remote, application.GitHubCredential{})
	if err != nil {
		t.Fatal(err)
	}
	if oid != "" {
		t.Fatal("remoto vazio tratado como commit em cache confirmado")
	}
}

// TestSelectedEmptyRepositoryReceivesOnlyJourneys verifica dois destinos distintos com refs reais.
func TestSelectedEmptyRepositoryReceivesOnlyJourneys(t *testing.T) {
	// first é um repositório com fonte que não pode receber a publicação destinada ao segundo.
	engine, root, first, destination := fixture(t)
	firstHead := seedJourneyRemote(t, first, map[string]entry{"apps/main.go": {Data: []byte("fonte remota"), Mode: filemode.Regular}})
	second := filepath.Join(t.TempDir(), "destino-vazio")
	if _, err := git.PlainInit(second, true); err != nil {
		t.Fatal(err)
	}
	destination.ID, destination.Name = 2, "repositorio_jornadas"
	// repo fixa o transporte derivado da identidade selecionada, sem credenciais reais.
	engine.remoteURL = func(repo application.GitHubRepository) string {
		if repo.ID == destination.ID {
			return second
		}
		return first
	}
	synctest.Journeys(t, root)
	// Code e preferências locais preexistentes não podem entrar no commit do pacote.
	if err := os.Mkdir(filepath.Join(root, "apps"), 0700); err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "apps"), "main.go", "fonte local")
	put(t, root, "README.md", "anotação local")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// binding é lido novamente, como na aplicação após reiniciar o workspace.
	binding, err := engine.Binding(root)
	if err != nil || binding.Repository.ID != 2 || binding.Pending != "" || binding.Published == "" {
		t.Fatal("publicação não confirmou o destino escolhido")
	}
	secondRepo, err := git.PlainOpen(second)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := secondRepo.Reference(plumbing.NewBranchReferenceName("main"), true)
	if err != nil || ref.Hash().String() != binding.Published {
		t.Fatal("OID não está na branch real do destino vazio")
	}
	files, err := loadTree(secondRepo, binding.Published)
	if err != nil || len(files) != 2 {
		t.Fatal("destino recebeu conteúdo além dos dois pacotes")
	}
	for name := range files {
		if !journeyPath(name) {
			t.Fatal("arquivo alheio publicado", name)
		}
	}
	firstRepo, err := git.PlainOpen(first)
	if err != nil {
		t.Fatal(err)
	}
	unchanged, err := head(firstRepo)
	if err != nil || unchanged != firstHead {
		t.Fatal("primeiro destino foi alterado indevidamente")
	}
}

// TestBindingCollisionHasDiagnosticCode preserva as jornadas e identifica o bloqueio sem erro genérico.
func TestBindingCollisionHasDiagnosticCode(t *testing.T) {
	// engine e remote reproduzem duas jornadas de mesmo nome com conteúdos diferentes.
	engine, root, remote, destination := fixture(t)
	// failure captura somente o registro sanitizado da tentativa, sem ler mensagens livres.
	var failure map[string]any
	engine.Diagnostics = syncdiag.New(func(event string, fields map[string]any) {
		if event == "github.git.bind.failed" {
			failure = fields
		}
	})
	seedJourneyRemote(t, remote, map[string]entry{"A.jornada": {Data: []byte("remota"), Mode: filemode.Regular}})
	put(t, root, "A.jornada", "local preservada")
	err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{})
	// problem deve levar a classificação segura ao logger, sem registrar os bytes do documento.
	var problem *application.SyncError
	if !errors.As(err, &problem) || problem.Code != "collision" {
		t.Fatalf("colisao sem codigo seguro: %T", err)
	}
	if failure["error_code"] != "collision" || failure["checkpoint"] != "workspace_integration" || failure["operation_id"] == "" {
		t.Fatal("log não identifica fase, código e tentativa da colisão")
	}
	data, readErr := os.ReadFile(filepath.Join(root, "A.jornada"))
	if readErr != nil || string(data) != "local preservada" {
		t.Fatal("colisao substituiu a jornada local")
	}
	if _, statErr := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(statErr) {
		t.Fatal("colisao concluiu vinculo invalido")
	}
}
