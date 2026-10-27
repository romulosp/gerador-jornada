package gitworkspace

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
	"gerador-jornada.local/editor/internal/adapter/userconfig"
	"gerador-jornada.local/editor/internal/application"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
)

// mutation descreve uma escrita reversível; pré-imagens ficam em arquivos privados separados.
type mutation struct {
	// Path é relativo à raiz validada e nunca contém .git.
	Path string `json:"path"`
	// Before e After são digests de conteúdo; string vazia representa ausência.
	Before string `json:"before"`
	// After identifica os bytes planejados; vazio registra exclusão do caminho.
	After string `json:"after"`
}

// journal permite reconhecer a atualização de HEAD mesmo se o processo morrer antes de salvar o vínculo.
type journal struct {
	// LockID comprova propriedade de um index.lock deixado por encerramento abrupto.
	LockID string `json:"lockId"`
	// Version e Binding validam formato e destino da operação interrompida.
	Version int `json:"version"`
	// Binding conserva o destino e as pendências que devem ser publicados ao concluir.
	Binding application.GitWorkspaceBinding `json:"binding"`
	// From e To são OIDs; vazio representa repositório ainda sem commit.
	From string `json:"from"`
	// To identifica a referência planejada, ou vazio enquanto não houver commit.
	To string `json:"to"`
	// Initial e Stage identificam instalação exclusiva de novos metadados.
	Initial bool `json:"initial"`
	// Stage é o diretório privado de metadados ainda não instalado na raiz.
	Stage string `json:"stage,omitempty"`
	// Files registra apenas caminhos efetivamente modificados pela transação.
	Files []mutation `json:"files"`
	// IndexBefore e IndexAfter preservam o staging anterior sem reset destrutivo.
	IndexBefore []byte `json:"indexBefore,omitempty"`
	// IndexAfter permite reconhecer a escrita da operação antes de restaurar o índice anterior.
	IndexAfter []byte `json:"indexAfter,omitempty"`
}

// operationDir mantém intenção e backups fora do workspace.
// root é a raiz canônica; a chave evita escrever seu caminho no nome do diretório privado.
func (e *Engine) operationDir(root string) string {
	return filepath.Join(e.Store.Root, "git-operations", userconfig.Key(root))
}

// removeStage só remove staging criado como irmão da raiz e com prefixo controlado.
// root limita a pasta-pai permitida e stage é conferido antes da remoção recursiva.
func (e *Engine) removeStage(root, stage string) {
	if stage == "" || !strings.EqualFold(filepath.Dir(stage), filepath.Dir(root)) || !strings.HasPrefix(filepath.Base(stage), ".jornada-git-") {
		return
	}
	_ = os.RemoveAll(stage)
}

// preflight bloqueia a substituição de um caminho editado em paralelo ou fora da base conhecida.
// before/after são as árvores de uma operação para root; erros precedem qualquer escrita.
func preflight(root string, before, after map[string]entry) error {
	// paths abrange também exclusões remotas.
	paths := map[string]bool{}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range before {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range after {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range paths {
		if !journeyPath(name) {
			continue
		}
		// old e oldOK representam o estado conhecido antes da integração deste caminho.
		old, oldOK := before[name]
		// next e nextOK representam o conteúdo planejado ou sua ausência no estado futuro.
		next, nextOK := after[name]
		if equal(old, oldOK, next, nextOK) {
			continue
		}
		// current conserva o conteúdo ou caminho efetivamente observado antes de qualquer mutação.
		current, exists, err := fileData(root, name)
		if err != nil {
			return err
		}
		if exists == nextOK && (!exists || bytes.Equal(current, next.Data)) {
			continue
		}
		if exists != oldOK || exists && !bytes.Equal(current, old.Data) {
			return application.SyncProblem("collision", fmt.Sprintf("Já existe uma jornada diferente chamada %s na pasta local e no repositório selecionado. Nenhum arquivo foi substituído. Escolha o repositório correto ou resolva a colisão entre os pacotes antes de conectar/publicar.", name))
		}
	}
	return nil
}

// hit interrompe uma fase exclusivamente quando um fixture instalar a injeção de falha.
// phase identifica a fronteira observável; produção sem fixture sempre retorna nil.
func (e *Engine) hit(ctx context.Context, phase string) error {
	// checkpoint remove o nome do arquivo; somente a fase técnica é registrada.
	checkpoint := strings.SplitN(phase, ":", 2)[0]
	e.Diagnostics.Event(ctx, "git.transaction.checkpoint", map[string]any{"checkpoint": checkpoint})
	if e.fail != nil {
		return e.fail(phase)
	}
	return nil
}

// apply publica arquivos, índice, referência e vínculo com intenção anterior à primeira mutação.
// ctx cancela etapas preparatórias; repo/binding fixam repositório e raiz. from/to são
// OIDs anteriores/futuros, before/after são as árvores e initial/stage delimitam nova
// instalação. Falha preserva journal para recuperação, sem reset ou overwrite indiscriminado.
func (e *Engine) apply(ctx context.Context, repo *git.Repository, binding application.GitWorkspaceBinding, from, to string, before, after map[string]entry, initial bool, stage string) error {
	// err impede iniciar uma alteração após cancelamento do contexto.
	if err := ctx.Err(); err != nil {
		return err
	}
	// err impede iniciar a transação quando o disco diverge das árvores conhecidas.
	if err := preflight(binding.Root, before, after); err != nil {
		return err
	}
	// directory é específico da raiz e já está protegido pelo gate entre processos.
	directory := e.operationDir(binding.Root)
	// err distingue journal ainda pendente antes de iniciar outra transação.
	if _, err := os.Stat(filepath.Join(directory, "journal.json")); err == nil {
		return fmt.Errorf("Há recuperação Git pendente neste workspace.")
	}
	// err impede prosseguir sem o diretório privado necessário à trava ou recuperação.
	if err := os.MkdirAll(directory, 0700); err != nil {
		return err
	}
	// intent guarda OID planejado antes de mover referência ou índice.
	intent := journal{Version: 1, Binding: binding, From: from, To: to, Initial: initial, Stage: stage, Files: []mutation{}}
	// nonce distingue a trava desta intenção de qualquer index.lock de cliente externo.
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	intent.LockID = digest(nonce)
	// paths reúne a união dos caminhos para considerar também exclusões.
	paths := map[string]bool{}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range before {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range after {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range paths {
		if !journeyPath(name) {
			continue
		}
		// previous guarda o valor anterior para recusar colisão ou recuperar somente bytes conhecidos.
		previous, previousOK := before[name]
		// next e nextOK representam o conteúdo planejado ou sua ausência no estado futuro.
		next, nextOK := after[name]
		if equal(previous, previousOK, next, nextOK) {
			continue
		}
		// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
		data, exists, err := fileData(binding.Root, name)
		if err != nil {
			return err
		}
		if exists == nextOK && (!exists || bytes.Equal(data, next.Data)) {
			continue
		}
		// change identifica o caminho e os digests da mutação reversível desta etapa.
		change := mutation{Path: name}
		if exists {
			change.Before = digest(data)
			if err = atomicfile.Write(ctx, filepath.Join(directory, userconfig.Key(name)+".before"), data); err != nil {
				return err
			}
		}
		if nextOK {
			change.After = digest(next.Data)
			if err = atomicfile.Write(ctx, filepath.Join(directory, userconfig.Key(name)+".after"), next.Data); err != nil {
				return err
			}
		}
		intent.Files = append(intent.Files, change)
	}
	// i/j ordenam as mutações por caminho para uma aplicação determinística.
	sort.Slice(intent.Files, func(i, j int) bool { return intent.Files[i].Path < intent.Files[j].Path })
	// stagingIndex contém toda a árvore futura, sem preparar arquivos alheios do working tree.
	stagingIndex := &index.Index{Version: 3}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name, item := range after {
		stagingIndex.Entries = append(stagingIndex.Entries, &index.Entry{Name: name, Hash: plumbing.ComputeHash(plumbing.BlobObject, item.Data), Mode: item.Mode, SkipWorktree: !journeyPath(name)})
	}
	// i/j estabilizam a ordem das entradas antes de serializar o índice futuro.
	sort.Slice(stagingIndex.Entries, func(i, j int) bool { return stagingIndex.Entries[i].Name < stagingIndex.Entries[j].Name })
	// buffer conserva o índice serializado integralmente antes de gravar metadados.
	buffer := bytes.Buffer{}
	// err impede publicar um índice cuja codificação não foi concluída.
	if err := index.NewEncoder(&buffer).Encode(stagingIndex); err != nil {
		return err
	}
	intent.IndexAfter = buffer.Bytes()
	if !initial {
		// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
		data, err := os.ReadFile(filepath.Join(binding.Root, ".git", "index"))
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		intent.IndexBefore = data
	}
	// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
	data, err := json.Marshal(intent)
	if err != nil {
		return err
	}
	if err = atomicfile.Write(ctx, filepath.Join(directory, "journal.json"), data); err != nil {
		return err
	}
	if err = e.hit(ctx, "intent"); err != nil {
		return err
	}
	// metadataLock é reconhecida também por clientes Git externos.
	metadata := filepath.Join(binding.Root, ".git")
	if initial {
		metadata = stage
	}
	// metadataLock impede escrita concorrente do índice por clientes Git que respeitam index.lock.
	metadataLock, err := os.OpenFile(filepath.Join(metadata, "index.lock"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Índice Git ocupado por outra operação.")
	}
	if _, err = metadataLock.WriteString("JornadaStudio:" + intent.LockID); err != nil {
		metadataLock.Close()
		return err
	}
	if err = metadataLock.Sync(); err != nil {
		metadataLock.Close()
		return err
	}
	defer func() { metadataLock.Close(); _ = os.Remove(filepath.Join(metadata, "index.lock")) }()
	// change identifica o caminho e os digests da mutação reversível desta etapa.
	for _, change := range intent.Files {
		// current conserva o conteúdo ou caminho efetivamente observado antes de qualquer mutação.
		current, exists, err := fileData(binding.Root, change.Path)
		if err != nil {
			return err
		}
		if exists != (change.Before != "") || exists && digest(current) != change.Before {
			return fmt.Errorf("Arquivo alterado externamente: %s; recuperação preservada.", change.Path)
		}
		// path fixa o arquivo ou prefixo exato desta etapa dentro da raiz ou do namespace privado.
		path := filepath.Join(binding.Root, filepath.FromSlash(change.Path))
		if change.After == "" {
			if err = os.Remove(path); err != nil {
				return err
			}
		} else {
			// next contém os bytes planejados preservados no namespace privado do journal.
			next, err := os.ReadFile(filepath.Join(directory, userconfig.Key(change.Path)+".after"))
			if err != nil {
				return err
			}
			if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return err
			}
			if change.Before == "" {
				// created mantém o arquivo novo aberto de forma exclusiva sem substituir dados anteriores.
				created, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
				if err != nil {
					return err
				}
				// writeErr conserva falha de gravação sem ocultá-la por sincronização ou fechamento posteriores.
				_, writeErr := created.Write(next)
				// syncErr distingue falha ao sincronizar o arquivo recém-criado.
				syncErr := created.Sync()
				// closeErr impede declarar a criação concluída se o fechamento do arquivo falhar.
				closeErr := created.Close()
				if writeErr != nil {
					return writeErr
				}
				if syncErr != nil {
					return syncErr
				}
				if closeErr != nil {
					return closeErr
				}
			} else if err = atomicfile.Write(ctx, path, next); err != nil {
				return err
			}
		}
		if err = e.hit(ctx, "file:"+change.Path); err != nil {
			return err
		}
	}
	if err = atomicfile.Write(ctx, filepath.Join(metadata, "index"), intent.IndexAfter); err != nil {
		return err
	}
	if err = e.hit(ctx, "index"); err != nil {
		return err
	}
	if initial {
		if err = atomicfile.Write(ctx, filepath.Join(metadata, "jornada-initial"), []byte(userconfig.Key(binding.Root))); err != nil {
			return err
		}
		// O handle do lock precisa fechar antes de mover metadados no Windows.
		metadataLock.Close()
		_ = os.Remove(filepath.Join(metadata, "index.lock"))
		if _, err = os.Lstat(filepath.Join(binding.Root, ".git")); !os.IsNotExist(err) {
			return fmt.Errorf("Metadados Git foram criados externamente; vínculo interrompido.")
		}
		if err = os.Rename(stage, filepath.Join(binding.Root, ".git")); err != nil {
			return err
		}
	} else if from != to && to != "" {
		// A referência é publicada por compare-and-set, nunca por force/reset.
		name := plumbing.NewBranchReferenceName(binding.Repository.Branch)
		// expected recusa publicação quando a referência anterior não for mais a revisada.
		var expected *plumbing.Reference
		if from != "" {
			expected = plumbing.NewHashReference(name, plumbing.NewHash(from))
		} else {
			if _, err = repo.Reference(name, true); !errorsIsMissing(err) {
				return fmt.Errorf("Branch inicial mudou durante a publicação.")
			}
		}
		if err = repo.Storer.CheckAndSetReference(plumbing.NewHashReference(name, plumbing.NewHash(to)), expected); err != nil {
			return fmt.Errorf("HEAD mudou durante a publicação; recuperação preservada.")
		}
	}
	if err = e.hit(ctx, "reference"); err != nil {
		return err
	}
	if err = e.Store.Write(context.WithoutCancel(ctx), bindingName(binding.Root), binding); err != nil {
		return err
	}
	if err = e.hit(ctx, "binding"); err != nil {
		return err
	}
	if err = os.RemoveAll(directory); err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(binding.Root, ".git", "jornada-initial"))
	return nil
}

// errorsIsMissing reconhece ausência de referência sem esconder falha de leitura.
// err deve vir da consulta de referência cuja ausência é permitida no primeiro commit.
func errorsIsMissing(err error) bool { return err == plumbing.ErrReferenceNotFound }

// recover conclui intenção cujo HEAD já foi publicado ou desfaz somente bytes ainda pertencentes à operação.
// ctx controla recuperação local para root; alteração externa preserva backups e exige ação manual.
func (e *Engine) recover(ctx context.Context, root string) error {
	// directory só pode apontar para o namespace privado desta raiz.
	directory := e.operationDir(root)
	// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
	data, err := os.ReadFile(filepath.Join(directory, "journal.json"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	// intent reconstrói o plano original antes de decidir entre concluir e restaurar.
	var intent journal
	if json.Unmarshal(data, &intent) != nil || intent.Version != 1 || !strings.EqualFold(intent.Binding.Root, root) {
		return fmt.Errorf("Journal Git inválido; recuperação manual necessária em %s", directory)
	}
	// lockRoot aponta apenas aos metadados desta raiz ou ao staging privado validado.
	lockRoot := filepath.Join(root, ".git")
	if intent.Initial {
		if _, statErr := os.Stat(lockRoot); os.IsNotExist(statErr) && strings.EqualFold(filepath.Dir(intent.Stage), filepath.Dir(root)) && strings.HasPrefix(filepath.Base(intent.Stage), ".jornada-git-") {
			lockRoot = intent.Stage
		}
	}
	// lockData permite remover somente o lock comprovadamente criado por esta intenção.
	lockData, lockErr := os.ReadFile(filepath.Join(lockRoot, "index.lock"))
	if lockErr == nil {
		if intent.LockID == "" || string(lockData) != "JornadaStudio:"+intent.LockID {
			return fmt.Errorf("Índice ocupado externamente; recuperação preservada em %s", directory)
		}
		if err = os.Remove(filepath.Join(lockRoot, "index.lock")); err != nil {
			return fmt.Errorf("Operação Git ainda ativa; tente a recuperação novamente.")
		}
	} else if !os.IsNotExist(lockErr) {
		return lockErr
	}
	// repo não procura pastas ancestrais.
	repo, openErr := git.PlainOpen(root)
	// current é o OID efetivo de HEAD, comparado aos pontos anterior e planejado da intenção.
	current := ""
	if openErr == nil {
		current, err = head(repo)
		if err != nil {
			return err
		}
	}
	if intent.Initial && openErr == nil {
		// marker comprova a autoria da instalação inicial antes de concluir a recuperação.
		marker, err := os.ReadFile(filepath.Join(root, ".git", "jornada-initial"))
		if err != nil || string(marker) != userconfig.Key(root) || current != intent.To {
			return fmt.Errorf("Metadados alterados externamente; recuperação manual necessária em %s", directory)
		}
		if err = e.Store.Write(ctx, bindingName(root), intent.Binding); err != nil {
			return err
		}
		e.removeStage(root, intent.Stage)
		if err = os.RemoveAll(directory); err != nil {
			return err
		}
		_ = os.Remove(filepath.Join(root, ".git", "jornada-initial"))
		return nil
	}
	if !intent.Initial && openErr == nil && current == intent.To {
		if err = e.Store.Write(ctx, bindingName(root), intent.Binding); err != nil {
			return err
		}
		return os.RemoveAll(directory)
	}
	if !intent.Initial && (openErr != nil || current != intent.From) {
		return fmt.Errorf("HEAD externo inesperado; recuperação manual necessária em %s", directory)
	}
	// position percorre mutações em ordem inversa para desfazer somente as escritas conhecidas.
	for position := len(intent.Files) - 1; position >= 0; position-- {
		// change identifica o caminho e os digests da mutação reversível desta etapa.
		change := intent.Files[position]
		// actual conserva os bytes atuais para não sobrescrever alterações externas durante recuperação.
		actual, exists, err := fileData(root, change.Path)
		if err != nil {
			return err
		}
		if exists == (change.Before != "") && (!exists || digest(actual) == change.Before) {
			continue
		}
		if exists != (change.After != "") || exists && digest(actual) != change.After {
			return fmt.Errorf("Arquivo alterado externamente: %s. Backup preservado em %s", change.Path, directory)
		}
		// path fixa o arquivo ou prefixo exato desta etapa dentro da raiz ou do namespace privado.
		path := filepath.Join(root, filepath.FromSlash(change.Path))
		if change.Before == "" {
			if exists {
				if err = os.Remove(path); err != nil {
					return err
				}
			}
		} else {
			// previous guarda o valor anterior para recusar colisão ou recuperar somente bytes conhecidos.
			previous, err := os.ReadFile(filepath.Join(directory, userconfig.Key(change.Path)+".before"))
			if err != nil || digest(previous) != change.Before {
				return fmt.Errorf("Backup inválido; recuperação manual necessária.")
			}
			if err = atomicfile.Write(ctx, path, previous); err != nil {
				return err
			}
		}
	}
	if !intent.Initial {
		// path fixa o arquivo ou prefixo exato desta etapa dentro da raiz ou do namespace privado.
		path := filepath.Join(root, ".git", "index")
		// actual conserva os bytes atuais para não sobrescrever alterações externas durante recuperação.
		actual, readErr := os.ReadFile(path)
		if readErr != nil && !os.IsNotExist(readErr) {
			return readErr
		}
		if !bytes.Equal(actual, intent.IndexBefore) && !bytes.Equal(actual, intent.IndexAfter) {
			return fmt.Errorf("Índice alterado externamente; recuperação manual necessária em %s", directory)
		}
		if len(intent.IndexBefore) == 0 {
			_ = os.Remove(path)
		} else if err = atomicfile.Write(ctx, path, intent.IndexBefore); err != nil {
			return err
		}
	}
	e.removeStage(root, intent.Stage)
	return os.RemoveAll(directory)
}

// Recover permite retomar operações interrompidas sem fazer rede ou publicar commits.
// root é validada e protegida pelo gate; o retorno informa conflito externo ou falha de recuperação.
func (e *Engine) Recover(ctx context.Context, root string) (recoveryErr error) {
	ctx, started := e.startDiagnostics(ctx, "git.recovery.start", root)
	defer func() { e.Diagnostics.Complete(ctx, "git.recovery.complete", started, recoveryErr) }()
	// absolute/resolved permitem ignorar workspaces sem journal, sem interferir no uso local opcional.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	resolvedPath, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return err
	}
	if _, err = os.Stat(filepath.Join(e.operationDir(resolvedPath), "journal.json")); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	// resolved/release protegem também a recuperação contra outra instância do aplicativo.
	resolved, release, err := e.gate(root)
	if err != nil {
		return err
	}
	defer release()
	return e.recover(ctx, resolved)
}
