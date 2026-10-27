package consent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/adapter/workspaceguard"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// targetSnapshot conserva identidade/precondição do destino autorizado, sem publicar nada.
type targetSnapshot struct {
	// identity é lógica/física, sem caminho de arquivo fornecido no commit posterior.
	identity string
	baseline workspaceguard.Baseline
}

// stagedArtifact mantém o único handle gravável internamente e depois apenas acesso de leitura em fluxo.
type stagedArtifact struct {
	// mu impede fechar o arquivo enquanto uma publicação autorizada futura lê seu conteúdo.
	mu   sync.Mutex
	file *os.File
	hash string
	size int64
	// closeErr conserva falha incerta em vez de declarar liberação que não foi comprovada.
	closeErr error
	// closed impede usar uma referência obtida imediatamente antes de Release retirar o registro.
	closed bool
}

// Freezer conserva somente artefatos privados/IDs; o documento de domínio não fica retido aqui.
type Freezer struct {
	// guard/profile/raízes são dependências fixas da composição, não configuração alterável por ferramenta.
	guard                                         *workspaceguard.Guard
	profile                                       ReviewProfile
	workspaceRoot, workspaceIdentity              string
	privateRoot, privatePhysical, privateIdentity string
	// mu serializa registro/seleção de diretório; IO de recurso usa sua própria trava.
	mu        sync.Mutex
	artifacts map[string]*stagedArtifact
}

// NewFreezer valida raízes e referências de host, sem criar pasta, staging ou autorização.
// privateRoot é configuração do operador e deve ficar fora da raiz física do workspace.
func NewFreezer(privateRoot, workspaceRoot string, profile ReviewProfile) (*Freezer, error) {
	if !validProfile(profile) {
		return nil, &application.AuthoringError{Code: "authorizationUnavailable"}
	}
	// guard/lease capturam identidade física antes de aceitar a composição.
	guard, err := workspaceguard.NewDefault()
	if err != nil {
		return nil, err
	}
	lease, err := guard.Open(workspaceRoot)
	if err != nil {
		return nil, err
	}
	defer lease.Close()
	workspacePhysical, err := lease.PhysicalPath()
	if err != nil {
		return nil, err
	}
	canonical, physical, err := resolvePrivate(guard, privateRoot)
	if err != nil {
		return nil, err
	}
	if withinPhysical(workspacePhysical, physical) {
		return nil, &application.AuthoringError{Code: "permissionDenied"}
	}
	return &Freezer{guard: guard, profile: profile, workspaceRoot: lease.Root, workspaceIdentity: lease.Identity, privateRoot: canonical, privatePhysical: physical, artifacts: map[string]*stagedArtifact{}}, nil
}

// withinPhysical compara nomes NT normalizados, incluindo aliases de volume resolvidos por handle.
func withinPhysical(root, candidate string) bool {
	root = strings.TrimSuffix(root, `\`)
	return strings.EqualFold(root, candidate) || strings.HasPrefix(strings.ToLower(candidate), strings.ToLower(root)+`\`)
}

// resolvePrivate resolve a parte existente antes de qualquer mkdir e conserva os segmentos ainda ausentes.
func resolvePrivate(guard *workspaceguard.Guard, root string) (string, string, error) {
	if root == "" {
		return "", "", os.ErrInvalid
	}
	// absolute/remaining identificam localização operacional, sem interpretar conteúdo de uma jornada como caminho.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	current, remaining := filepath.Clean(absolute), []string{}
	for {
		// info distingue diretório existente de componente inválido, sem criar descendentes.
		info, err := os.Stat(current)
		if err == nil {
			if !info.IsDir() {
				return "", "", os.ErrInvalid
			}
			break
		}
		if !os.IsNotExist(err) {
			return "", "", err
		}
		name := filepath.Base(current)
		if !workspaceguard.ValidChildName(name) {
			return "", "", os.ErrInvalid
		}
		remaining = append([]string{name}, remaining...)
		parent := filepath.Dir(current)
		if parent == current {
			return "", "", os.ErrInvalid
		}
		current = parent
	}
	// lease resolve caminhos físicos/aliases da parte existente, nunca apenas concatenação textual de raiz autorizada.
	lease, err := guard.Open(current)
	if err != nil {
		return "", "", err
	}
	defer lease.Close()
	physical, err := lease.PhysicalPath()
	if err != nil {
		return "", "", err
	}
	canonical := lease.Root
	for _, name := range remaining {
		canonical = filepath.Join(canonical, name)
		physical = strings.TrimSuffix(physical, `\`) + `\` + name
	}
	return canonical, physical, nil
}

// target revalida base/ausência sem alocar pacote integral e sem segurar trava enquanto aguarda humano.
func (freezer *Freezer) target(ctx context.Context, request application.AuthoringFreezeRequest) (targetSnapshot, string, error) {
	if !catalog.ValidJourneyFilename(request.Options.Filename) || !workspaceguard.ValidChildName(request.Options.Filename) {
		return targetSnapshot{}, "", &application.AuthoringError{Code: "invalidArgument"}
	}
	// lease conserva a raiz física durante identidade e hash do destino.
	lease, err := freezer.guard.Open(freezer.workspaceRoot)
	if err != nil {
		return targetSnapshot{}, "", err
	}
	defer lease.Close()
	if lease.Identity != freezer.workspaceIdentity {
		return targetSnapshot{}, "", &application.AuthoringError{Code: "externalModification"}
	}
	physical, err := lease.PhysicalPath()
	if err != nil {
		return targetSnapshot{}, "", err
	}
	baseline, err := lease.BaselineContext(ctx, request.Options.Filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return targetSnapshot{}, "", err
	}
	if request.Options.Intent == "create" {
		if baseline.Exists {
			return targetSnapshot{}, "", &application.AuthoringError{Code: "externalModification"}
		}
	} else if request.Options.Intent == "replace" {
		// source deve ser exatamente a base aberta, não outro arquivo existente escolhido pelo agente.
		source := request.Baseline
		if !source.Exists || source.RootIdentity != freezer.workspaceIdentity || !strings.EqualFold(source.Filename, request.Options.Filename) || !baseline.Exists || source.Identity != baseline.Identity || source.SHA256 != baseline.SHA256 {
			return targetSnapshot{}, "", &application.AuthoringError{Code: "externalModification"}
		}
	} else {
		return targetSnapshot{}, "", &application.AuthoringError{Code: "invalidArgument"}
	}
	return targetSnapshot{identity: freezer.workspaceIdentity + "/" + strings.ToLower(request.Options.Filename), baseline: baseline}, physical, nil
}

// transient cria somente após revalidar finalidade/identidade da pasta privada, sem abrir caminho arbitrário.
func (freezer *Freezer) transient(workspacePhysical, name string) (*os.File, error) {
	freezer.mu.Lock()
	defer freezer.mu.Unlock()
	// canonical/physical detectam troca de um ancestral antes da criação operacional de diretórios.
	canonical, physical, err := resolvePrivate(freezer.guard, freezer.privateRoot)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(physical, freezer.privatePhysical) || withinPhysical(workspacePhysical, physical) {
		return nil, &application.AuthoringError{Code: "externalModification"}
	}
	if err := os.MkdirAll(canonical, 0700); err != nil {
		return nil, err
	}
	lease, err := freezer.guard.Acquire(canonical)
	if err != nil {
		return nil, err
	}
	defer lease.Close()
	actual, err := lease.PhysicalPath()
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(actual, freezer.privatePhysical) || (freezer.privateIdentity != "" && freezer.privateIdentity != lease.Identity) {
		return nil, &application.AuthoringError{Code: "externalModification"}
	}
	freezer.privateIdentity = lease.Identity
	return lease.CreateTransient(name)
}

// Freeze escreve ZIP limitado diretamente em staging delete-on-close e produz revisão canônica sobre seu hash.
// Nunca escreve filename no workspace nem verifica/aprova uma decisão humana.
func (freezer *Freezer) Freeze(ctx context.Context, request application.AuthoringFreezeRequest) (frozen application.AuthoringFrozenReview, resultErr error) {
	// target/physical são verificados antes de criar qualquer artefato privado.
	target, physical, err := freezer.target(ctx, request)
	if err != nil {
		return frozen, err
	}
	if err := ctx.Err(); err != nil {
		return frozen, err
	}
	artifactID := project.NewID()
	file, err := freezer.transient(physical, artifactID+".stage")
	if err != nil {
		return frozen, err
	}
	retained := false
	defer func() {
		if !retained {
			if closeErr := file.Close(); closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
				// Uma falha incerta permanece rastreada para não perder a reserva de staging.
				freezer.mu.Lock()
				freezer.artifacts[artifactID] = &stagedArtifact{file: file, closeErr: closeErr}
				freezer.mu.Unlock()
				frozen.ArtifactID, resultErr = artifactID, closeErr
			}
		}
	}()
	// digest é atualizado nos mesmos blocos efetivamente escritos, sem segunda cópia do ZIP.
	digest := sha256.New()
	if err := projectzip.Write(ctx, io.MultiWriter(file, digest), request.Document); err != nil {
		if errors.Is(err, projectzip.ErrLimit) {
			return frozen, &application.AuthoringError{Code: "resourceLimit"}
		}
		return frozen, err
	}
	if err := file.Sync(); err != nil {
		return frozen, err
	}
	info, err := file.Stat()
	if err != nil {
		return frozen, err
	}
	// Verify usa o mesmo handle e confirma legibilidade/integridade sem carregar ZIP/imagens completos.
	if err := projectzip.Verify(ctx, file, info.Size()); err != nil {
		if errors.Is(err, projectzip.ErrLimit) {
			return frozen, &application.AuthoringError{Code: "resourceLimit"}
		}
		return frozen, err
	}
	contentHash := hex.EncodeToString(digest.Sum(nil))
	review, err := encodeReview(freezer.profile, request, target, contentHash)
	if err != nil {
		return frozen, err
	}
	if err := ctx.Err(); err != nil {
		return frozen, err
	}
	freezer.mu.Lock()
	freezer.artifacts[artifactID] = &stagedArtifact{file: file, hash: contentHash, size: info.Size()}
	freezer.mu.Unlock()
	retained = true
	return application.AuthoringFrozenReview{ArtifactID: artifactID, ContentHash: contentHash, ContentSizeBytes: info.Size(), Bytes: review}, nil
}

// Release fecha o handle privado e nunca remove um caminho interpretado a partir do ID.
func (freezer *Freezer) Release(artifactID string) error {
	freezer.mu.Lock()
	// artifact é procurado somente no registro deste processo; identidade desconhecida não apaga nada.
	artifact := freezer.artifacts[artifactID]
	if artifact == nil {
		freezer.mu.Unlock()
		return nil
	}
	freezer.mu.Unlock()
	artifact.mu.Lock()
	defer artifact.mu.Unlock()
	if artifact.closed {
		return nil
	}
	if artifact.closeErr != nil {
		return artifact.closeErr
	}
	if err := artifact.file.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		artifact.closeErr = err
		return err
	}
	artifact.closed = true
	freezer.mu.Lock()
	delete(freezer.artifacts, artifactID)
	freezer.mu.Unlock()
	return nil
}

// WriteArtifact fornece bytes congelados por referência interna, sem abrir path informado pelo agente.
// Destino pertence a uma publicação autorizada futura; esta função não publica nem autoriza por si só.
func (freezer *Freezer) WriteArtifact(ctx context.Context, artifactID string, destination io.Writer) error {
	freezer.mu.Lock()
	// artifact mantém handle/imutabilidade mesmo se o antigo nome da pasta mudar.
	artifact := freezer.artifacts[artifactID]
	if artifact == nil {
		freezer.mu.Unlock()
		return &application.AuthoringError{Code: "staleReview"}
	}
	freezer.mu.Unlock()
	artifact.mu.Lock()
	defer artifact.mu.Unlock()
	if artifact.closed {
		return &application.AuthoringError{Code: "staleReview"}
	}
	if artifact.closeErr != nil {
		return artifact.closeErr
	}
	if _, err := artifact.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	digest := sha256.New()
	buffer, total := make([]byte, 64<<10), int64(0)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		read, err := artifact.file.Read(buffer)
		if read > 0 {
			written, writeErr := destination.Write(buffer[:read])
			if writeErr != nil {
				return writeErr
			}
			if written != read {
				return io.ErrShortWrite
			}
			_, _ = digest.Write(buffer[:read])
			total += int64(read)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	if total != artifact.size || hex.EncodeToString(digest.Sum(nil)) != artifact.hash {
		return &application.AuthoringError{Code: "externalModification"}
	}
	return nil
}

// A asserção mantém o provider limitado ao contrato consumido pela aplicação.
var _ application.AuthoringReviewFreezer = (*Freezer)(nil)
