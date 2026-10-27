package gitworkspace

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/adapter/userconfig"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/synctest"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
)

// init usa servidor Go incorporado no transporte local de fixtures, evitando git-upload-pack externo.
func init() {
	client.InstallProtocol("file", fixtureTransport{Transport: server.NewClient(server.DefaultLoader)})
}

// fixtureTransport ignora haves desconhecidos como um servidor Git completo, sem executar git-upload-pack.
type fixtureTransport struct{ transport.Transport }

// fixtureSession conserva os objetos do remoto apenas para filtrar negociação do servidor de teste.
type fixtureSession struct {
	transport.UploadPackSession
	// objects é o armazenamento bare descartável, nunca um repositório GitHub real.
	objects storer.Storer
}

// NewUploadPackSession adapta a negociação limitada do servidor Go usado pelos fixtures locais.
func (f fixtureTransport) NewUploadPackSession(endpoint *transport.Endpoint, auth transport.AuthMethod) (transport.UploadPackSession, error) {
	// session e objects pertencem ao mesmo endpoint descartável.
	session, err := f.Transport.NewUploadPackSession(endpoint, auth)
	if err != nil {
		return nil, err
	}
	objects, err := server.DefaultLoader.Load(endpoint)
	if err != nil {
		session.Close()
		return nil, err
	}
	return &fixtureSession{UploadPackSession: session, objects: objects}, nil
}

// UploadPack não anuncia como disponíveis commits locais que o servidor nunca recebeu.
func (s *fixtureSession) UploadPack(ctx context.Context, request *packp.UploadPackRequest) (*packp.UploadPackResponse, error) {
	// filtered conserva wants/capabilities, mas remove haves que não existem no remoto.
	filtered := *request
	filtered.Haves = nil
	for _, hash := range request.Haves {
		if err := s.objects.HasEncodedObject(hash); err == nil {
			filtered.Haves = append(filtered.Haves, hash)
		} else if !errors.Is(err, plumbing.ErrObjectNotFound) {
			return nil, err
		}
	}
	return s.UploadPackSession.UploadPack(ctx, &filtered)
}

// TestPublishedJourneysNavigation verifica o conjunto real ZIP e a navegação após publicação e reabertura.
func TestPublishedJourneysNavigation(t *testing.T) {
	// engine usa remoto bare local; os documentos não contêm credenciais ou URLs remotas.
	engine, root, remote, destination := fixture(t)
	synctest.Journeys(t, root)
	before, err := os.ReadFile(filepath.Join(root, "Principal.jornada"))
	if err != nil {
		t.Fatal(err)
	}
	if err = engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	if err = publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	binding, err := engine.Binding(root)
	if err != nil {
		t.Fatal(err)
	}
	repo, err := git.PlainOpen(remote)
	if err != nil {
		t.Fatal(err)
	}
	files, err := loadTree(repo, binding.Published)
	if err != nil || string(files["Principal.jornada"].Data) != string(before) {
		t.Fatal("bytes ZIP alterados")
	}
	// restored recebe somente os pacotes remotos em outra pasta, sem pastas do software original.
	restored := filepath.Join(t.TempDir(), "restaurado")
	if err = os.Mkdir(restored, 0700); err != nil {
		t.Fatal(err)
	}
	if err = engine.Bind(context.Background(), restored, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	root = restored
	// restoredBytes comprova que o projeto reaberto vem exatamente do pacote publicado.
	restoredBytes, err := os.ReadFile(filepath.Join(root, "Principal.jornada"))
	if err != nil || string(restoredBytes) != string(before) {
		t.Fatal("pacote remoto não restaurou os bytes salvos")
	}
	// document é reaberto pelo codec real, sem reconstruir vínculos no teste.
	document, err := projectzip.Load(filepath.Join(root, "Principal.jornada"))
	if err != nil {
		t.Fatal(err)
	}
	preview, err := application.NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	preview.CurrentJourney = "Principal.jornada"
	preview.LoadJourney = func(name string) (application.Document, error) { return projectzip.Load(filepath.Join(root, name)) }
	menu := document.Screens[0].Components[0]
	if _, err = preview.MenuActivate(menu.ID, menu.Properties.Menu.Root.Items[0].ID); err != nil {
		t.Fatal(err)
	}
	if _, err = preview.MenuActivate(menu.ID, menu.Properties.Menu.Root.Items[0].Submenu.Items[1].ID); err != nil || preview.CurrentJourney != "Destino.jornada" {
		t.Fatalf("destino perdido: %v", err)
	}
}

// TestAbruptProcessRecovery encerra um processo real antes dos defers para verificar locks e intenção persistida.
func TestAbruptProcessRecovery(t *testing.T) {
	// engine e diretórios permanecem vivos após a saída do subprocesso isolado.
	engine, root, remote, destination := fixture(t)
	put(t, root, "A.jornada", "bytes preservados")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// executable é o binário de testes atual, sem ferramentas externas ou credenciais.
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(executable, "-test.run=^TestCrashProcessHelper$")
	command.Env = append(os.Environ(), "JORNADA_CRASH_ROOT="+root, "JORNADA_CRASH_REMOTE="+remote, "JORNADA_CRASH_CONFIG="+engine.Store.Root)
	if err = command.Run(); err == nil {
		t.Fatal("subprocesso não foi interrompido")
	}
	if err = engine.Recover(context.Background(), root); err != nil {
		t.Fatal(err)
	}
	binding, err := engine.Binding(root)
	if err != nil || binding.Pending == "" {
		t.Fatal("commit interrompido perdido")
	}
	if err = publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
}

// TestCrashProcessHelper é executado somente pelo harness de encerramento abrupto acima.
func TestCrashProcessHelper(t *testing.T) {
	// root é sempre um diretório descartável explicitamente passado pelo teste pai.
	root := os.Getenv("JORNADA_CRASH_ROOT")
	if root == "" {
		return
	}
	engine := New(&userconfig.Store{Root: os.Getenv("JORNADA_CRASH_CONFIG")})
	engine.remoteURL = func(application.GitHubRepository) string { return os.Getenv("JORNADA_CRASH_REMOTE") }
	engine.fail = func(phase string) error {
		if phase == "reference" {
			os.Exit(23)
		}
		return nil
	}
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
}

// TestBindingCollisionAndRollback preserva arquivos anteriores e desfaz apenas materialização própria incompleta.
func TestBindingCollisionAndRollback(t *testing.T) {
	// source publica conteúdo remoto que será integrado por outro workspace.
	engine, source, _, destination := fixture(t)
	put(t, source, "A.jornada", "remoto")
	if err := engine.Bind(context.Background(), source, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	if err := publish(t, engine, source); err != nil {
		t.Fatal(err)
	}
	// second fica ao lado da raiz original, sem herdar seus metadados Git.
	second := filepath.Join(filepath.Dir(source), "second")
	if err := os.Mkdir(second, 0700); err != nil {
		t.Fatal(err)
	}
	put(t, second, "A.jornada", "local diferente")
	if err := engine.Bind(context.Background(), second, destination, application.GitHubCredential{}); err == nil {
		t.Fatal("colisão aceita")
	}
	data, _ := os.ReadFile(filepath.Join(second, "A.jornada"))
	if string(data) != "local diferente" {
		t.Fatal("conteúdo local perdido")
	}
	if _, err := os.Stat(filepath.Join(second, ".git")); !os.IsNotExist(err) {
		t.Fatal("vínculo parcial publicado")
	}
	if err := os.Remove(filepath.Join(second, "A.jornada")); err != nil {
		t.Fatal(err)
	}
	engine.fail = func(phase string) error {
		if phase == "file:A.jornada" {
			return errors.New("falha parcial")
		}
		return nil
	}
	if err := engine.Bind(context.Background(), second, destination, application.GitHubCredential{}); err == nil {
		t.Fatal("falha não induzida")
	}
	engine.fail = nil
	if err := engine.Recover(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(second, "A.jornada")); !os.IsNotExist(err) {
		t.Fatal("arquivo parcial não foi recuperado")
	}
	if err := engine.Bind(context.Background(), second, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
}

// fixture prepara remoto bare e workspace distintos, sem Git instalado ou rede externa.
func fixture(t *testing.T) (*Engine, string, string, application.GitHubRepository) {
	t.Helper()
	// base isola todos os efeitos deste teste.
	base := t.TempDir()
	// remote identifica o repositório bare descartável usado como destino do teste.
	remote := filepath.Join(base, "remote")
	// root identifica exclusivamente a pasta descartável usada pelo cenário.
	root := filepath.Join(base, "workspace")
	// err garante que o remoto bare de teste existe antes de exercitar o vínculo.
	if _, err := git.PlainInit(remote, true); err != nil {
		t.Fatal(err)
	}
	// err exige a pasta temporária do workspace antes de executar o cenário.
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	// engine e os caminhos retornados pela fixture isolam estado, remoto e workspace deste cenário.
	engine := New(&userconfig.Store{Root: filepath.Join(base, "preferences")})
	engine.remoteURL = func(application.GitHubRepository) string { return remote }
	return engine, root, remote, application.GitHubRepository{ID: 1, Owner: "fixture", Name: "project", Branch: "main", Writable: true}
}

// put grava uma fixture local e falha imediatamente se não puder preservar os bytes fornecidos.
func put(t *testing.T, root, name, content string) {
	t.Helper()
	// err impede usar fixture com conteúdo ausente ou parcialmente preparado.
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

// remoteCommit simula outro cliente, mantendo histórico real no remoto descartável.
func remoteCommit(t *testing.T, remote, name, content string) string {
	t.Helper()
	// directory recebe um checkout próprio para evitar interferência com o workspace sob teste.
	directory := filepath.Join(t.TempDir(), "clone")
	// repo acessa somente o repositório exato da raiz sob teste ou recuperação.
	repo, err := git.PlainClone(directory, false, &git.CloneOptions{URL: remote, ReferenceName: plumbing.NewBranchReferenceName("main")})
	if err != nil {
		t.Fatal(err)
	}
	put(t, directory, name, content)
	// worktree prepara somente as alterações sintéticas do outro cliente de teste.
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = worktree.Add(name); err != nil {
		t.Fatal(err)
	}
	// oid fixa o commit criado ou lido para conferir integridade e ancestralidade.
	oid, err := worktree.Commit("Outro cliente", &git.CommitOptions{Author: &object.Signature{Name: "Outro", Email: "fixture@example.invalid", When: time.Now()}})
	if err != nil {
		t.Fatal(err)
	}
	if err = repo.Push(&git.PushOptions{RefSpecs: []config.RefSpec{"refs/heads/main:refs/heads/main"}}); err != nil {
		t.Fatal(err)
	}
	return oid.String()
}

// publish revisa e confirma usando somente identidade sintética.
func publish(t *testing.T, e *Engine, root string) error {
	t.Helper()
	// review fixa o estado que será revalidado na confirmação.
	review, err := e.Review(context.Background(), root, application.GitHubCredential{})
	if err != nil {
		return err
	}
	review.AccountID = 7
	return e.Publish(context.Background(), review, "Atualiza jornadas", application.GitHubAccount{ID: 7, Login: "fixture"}, application.GitHubCredential{})
}

// TestFirstPublicationEmptyRemote comprova o caso sem HEAD, commit raiz sem pais e envio do conjunto.
func TestFirstPublicationEmptyRemote(t *testing.T) {
	// engine/root/remote são integralmente descartáveis.
	engine, root, remote, destination := fixture(t)
	put(t, root, "Principal.jornada", "bytes A → Destino.jornada")
	put(t, root, "Destino.jornada", "bytes B")
	put(t, root, "backup.tmp", "não publicar")
	// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// repo acessa somente o repositório exato da raiz sob teste ou recuperação.
	repo, err := git.PlainOpen(root)
	if err != nil {
		t.Fatal(err)
	}
	// before e bOK identificam conteúdo e presença da base comum deste caminho.
	before, _ := head(repo)
	if before != "" {
		t.Fatal("vínculo criou commit")
	}
	if err = publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// local identifica o commit local que deve corresponder exatamente à referência remota.
	local, err := head(repo)
	if err != nil {
		t.Fatal(err)
	}
	// commit fixa a árvore do OID solicitado antes da leitura dos blobs.
	commit, err := repo.CommitObject(plumbing.NewHash(local))
	if err != nil {
		t.Fatal(err)
	}
	if len(commit.ParentHashes) != 0 {
		t.Fatal("primeiro commit não é raiz")
	}
	// remoteRepo permite conferir o destino bare sem usar o Git instalado.
	remoteRepo, err := git.PlainOpen(remote)
	if err != nil {
		t.Fatal(err)
	}
	// ref representa a branch remota usada para confirmar a entrega do commit.
	ref, err := remoteRepo.Reference(plumbing.NewBranchReferenceName("main"), true)
	if err != nil || ref.Hash().String() != local {
		t.Fatalf("push não confirmado: %v", err)
	}
	// files mantém o conjunto de arquivos a comparar sem materializar conteúdo não validado.
	files, err := loadTree(remoteRepo, local)
	if err != nil || len(files) != 2 || string(files["Principal.jornada"].Data) != "bytes A → Destino.jornada" {
		t.Fatalf("conteúdo incorreto: %v", err)
	}
	if err = publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// same permite provar que repetir a publicação não criou commit vazio.
	same, _ := head(repo)
	if same != local {
		t.Fatal("commit vazio foi criado")
	}
}

// TestPendingRetryAndUncertainPush confere identidade de OID antes e depois de aceitação remota sem resposta.
func TestPendingRetryAndUncertainPush(t *testing.T) {
	// phase escolhe o ponto de falha para distinguir envio ainda não tentado de resposta perdida.
	for _, phase := range []string{"before-push", "after-push"} {
		t.Run(phase, func(t *testing.T) {
			// t isola as asserções da fase de falha selecionada neste subteste.
			// fixture usa nova raiz/cofre para cada fase de falha.
			engine, root, _, destination := fixture(t)
			put(t, root, "A.jornada", "conteúdo")
			// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
			if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
				t.Fatal(err)
			}
			// step identifica a fronteira transacional onde a fixture deve interromper o fluxo.
			engine.fail = func(step string) error {
				if step == phase {
					return errors.New("interrupção simulada")
				}
				return nil
			}
			// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
			if err := publish(t, engine, root); err == nil {
				t.Fatal("falha não induzida")
			}
			// pending conserva o OID ainda não confirmado antes da repetição de envio.
			pending, err := engine.Binding(root)
			if err != nil || pending.Pending == "" {
				t.Fatalf("commit perdido: %v", err)
			}
			engine.fail = nil
			if err = publish(t, engine, root); err != nil {
				t.Fatal(err)
			}
			// completed permite conferir que a confirmação esvaziou a pendência sem duplicar o commit.
			completed, err := engine.Binding(root)
			if err != nil || completed.Pending != "" || completed.Published != pending.Pending {
				t.Fatalf("retentativa duplicou ou perdeu commit: %v", err)
			}
		})
	}
}

// TestRemoteDisjointAndConflict preserva novidade remota e bloqueia ZIP alterado por ambos os lados.
func TestRemoteDisjointAndConflict(t *testing.T) {
	// engine cria um histórico comum pelo primeiro envio.
	engine, root, remote, destination := fixture(t)
	put(t, root, "A.jornada", "inicial")
	// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	remoteCommit(t, remote, "Remota.jornada", "remoto")
	put(t, root, "A.jornada", "local")
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
	data, err := os.ReadFile(filepath.Join(root, "Remota.jornada"))
	if err != nil || string(data) != "remoto" {
		t.Fatal("novidade remota perdida")
	}
	remoteCommit(t, remote, "A.jornada", "concorrente")
	put(t, root, "A.jornada", "ainda local")
	if _, err = engine.Review(context.Background(), root, application.GitHubCredential{}); err == nil {
		t.Fatal("conflito aceito")
	}
	data, _ = os.ReadFile(filepath.Join(root, "A.jornada"))
	if string(data) != "ainda local" {
		t.Fatal("local sobrescrito")
	}
}

// TestPendingDivergentMerge conserva o commit pendente como segundo pai sem rebase.
func TestPendingDivergentMerge(t *testing.T) {
	// engine começa com ancestral comum publicado.
	engine, root, remote, destination := fixture(t)
	put(t, root, "A.jornada", "inicial")
	// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	put(t, root, "A.jornada", "edição pendente")
	// step identifica a fronteira transacional onde a fixture deve interromper o fluxo.
	engine.fail = func(step string) error {
		if step == "before-push" {
			return errors.New("rede")
		}
		return nil
	}
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err == nil {
		t.Fatal("falha esperada")
	}
	// pending conserva o OID ainda não confirmado antes da repetição de envio.
	pending, _ := engine.Binding(root)
	engine.fail = nil
	// other identifica o avanço remoto feito por um cliente independente da aplicação.
	other := remoteCommit(t, remote, "README.md", "novo remoto")
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err != nil {
		t.Fatalf("integração local: %v; causa: %v", err, errors.Unwrap(err))
	}
	// repo acessa somente o repositório exato da raiz sob teste ou recuperação.
	repo, _ := git.PlainOpen(root)
	// oid fixa o commit criado ou lido para conferir integridade e ancestralidade.
	oid, _ := head(repo)
	// commit fixa a árvore do OID solicitado antes da leitura dos blobs.
	commit, err := repo.CommitObject(plumbing.NewHash(oid))
	if err != nil {
		t.Fatal(err)
	}
	if len(commit.ParentHashes) != 2 || commit.ParentHashes[0].String() != other || commit.ParentHashes[1].String() != pending.Pending {
		t.Fatal("pais não preservados")
	}
}

// TestStagingAndStaleReview rejeita conteúdo alheio preparado e snapshot que mudou após a revisão.
func TestStagingAndStaleReview(t *testing.T) {
	// fixture mantém staging e publicação dentro de diretórios temporários.
	engine, root, _, destination := fixture(t)
	put(t, root, "A.jornada", "inicial")
	// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// review captura o snapshot que deve continuar válido até a confirmação.
	review, err := engine.Review(context.Background(), root, application.GitHubCredential{})
	if err != nil {
		t.Fatal(err)
	}
	put(t, root, "A.jornada", "mudou")
	review.AccountID = 7
	if err = engine.Publish(context.Background(), review, "mensagem", application.GitHubAccount{ID: 7, Login: "fixture"}, application.GitHubCredential{}); err == nil {
		t.Fatal("snapshot obsoleto aceito")
	}
	put(t, root, "outro.txt", "externo")
	// repo acessa somente o repositório exato da raiz sob teste ou recuperação.
	repo, _ := git.PlainOpen(root)
	// worktree prepara somente as alterações sintéticas do outro cliente de teste.
	worktree, _ := repo.Worktree()
	if _, err = worktree.Add("outro.txt"); err != nil {
		t.Fatal(err)
	}
	if _, err = engine.Review(context.Background(), root, application.GitHubCredential{}); err == nil {
		t.Fatal("staging alheio aceito")
	}
}

// TestCommitCrashRecovery reconhece o OID planejado após mover a ref sem salvar o estado.
func TestCommitCrashRecovery(t *testing.T) {
	// engine interrompe o primeiro commit antes da persistência do vínculo pendente.
	engine, root, _, destination := fixture(t)
	put(t, root, "A.jornada", "inicial")
	// err impede avaliar publicação ou recuperação sem concluir o vínculo de teste.
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	// step identifica a fronteira transacional onde a fixture deve interromper o fluxo.
	engine.fail = func(step string) error {
		if step == "reference" {
			return errors.New("processo encerrado")
		}
		return nil
	}
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err == nil {
		t.Fatal("interrupção não executada")
	}
	// repo acessa somente o repositório exato da raiz sob teste ou recuperação.
	repo, _ := git.PlainOpen(root)
	// planned identifica o commit que a recuperação deve reconhecer após a queda simulada.
	planned, _ := head(repo)
	engine.fail = nil
	// err confirma se a recuperação reconheceu a intenção persistida sem perder o commit.
	if err := engine.Recover(context.Background(), root); err != nil {
		t.Fatal(err)
	}
	// binding permite conferir que a recuperação persistiu o commit planejado como pendente.
	binding, _ := engine.Binding(root)
	if binding.Pending != planned {
		t.Fatal("intenção do commit perdida")
	}
	// err distingue o resultado de envio esperado pelo cenário de falha ou sucesso.
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
}
