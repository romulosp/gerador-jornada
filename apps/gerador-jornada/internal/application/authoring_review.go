package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// AuthoringReviewOptions recebe intenção de revisão, sem autoridade de persistência.
type AuthoringReviewOptions struct {
	// Filename é regular na raiz de composição; Intent aceita create/replace para jornada.
	Filename, Intent string
	// AllowDraft permite revisão de rascunho legível com diagnósticos, nunca corrupção ou aprovação.
	AllowDraft bool
}

// AuthoringFreezeRequest fornece visão protegida ao adaptador confiável de codificação/staging.
// Document/BaselineDocument são somente leitura e não podem ser retidos pelo adaptador após Freeze.
type AuthoringFreezeRequest struct {
	// SessionID/Ticket vinculam o artefato à preparação corrente e à requestId separada do commit.
	SessionID string
	Ticket    AuthoringControlTicket
	// Revision/Options/Baseline preservam conteúdo, intenção e precondição já capturada ao abrir.
	Revision uint64
	Options  AuthoringReviewOptions
	Baseline AuthoringBaseline
	// Document/BaselineDocument usam recursos já retidos, sem clone integral adicional.
	Document, BaselineDocument Document
	// Summary/Diagnostics são evidências do aplicativo, sem resumo produzido pela LLM.
	Summary     []string
	Diagnostics []Diagnostic
}

// AuthoringFrozenReview contém referência de staging e bytes exatos da revisão, sem autorização.
type AuthoringFrozenReview struct {
	// Ticket é atribuído pelo coordenador; nunca é identidade apresentada pelo agente.
	Ticket AuthoringControlTicket
	// Revision identifica o snapshot de domínio que originou a revisão.
	Revision uint64
	// ArtifactID é referência opaca privada; ContentHash/ContentSizeBytes descrevem o pacote final.
	ArtifactID, ContentHash string
	ContentSizeBytes        int64
	// Bytes/Hash são os bytes JSON canônicos congelados e seu SHA-256, não o ZIP inteiro.
	Bytes []byte
	Hash  string
}

// AuthoringReviewFreezer codifica/congela e libera staging sem publicar um destino final.
// Em erro, Freeze deve liberar o parcial ou retornar ArtifactID rastreável para aposentadoria.
// Release é idempotente e não depende de interação humana ou acesso ao registro de autoria.
type AuthoringReviewFreezer interface {
	Freeze(context.Context, AuthoringFreezeRequest) (AuthoringFrozenReview, error)
	Release(artifactID string) error
}

// authoringReviewState anexa a única revisão/artefato corrente à sessão já existente.
type authoringReviewState struct {
	// fingerprint/ready distinguem uma resposta completa de artefato parcial aguardando limpeza.
	fingerprint string
	ready       bool
	// frozen não concede consentimento; release conserva somente provider/ID do artefato.
	frozen  AuthoringFrozenReview
	release func() error
}

// authoringRetirement mantém staging contabilizado enquanto a liberação física não terminou.
type authoringRetirement struct {
	// release executa IO fora dos mutexes de registro e sessão.
	release func() error
	// stage conserva a reserva até Release concluir; erro mantém o registro para nova tentativa.
	stage int64
}

// PrepareReview congela um candidato sob a proteção da sessão, sem consumir edição ou publicar jornada.
// A implementação do freezer recebe perfis/raízes pela composição, nunca do agente nesta chamada.
func (a *Authoring) PrepareReview(ctx context.Context, id, requestID, fingerprint string, expected, generation uint64, options AuthoringReviewOptions, freezer AuthoringReviewFreezer) (AuthoringFrozenReview, error) {
	defer a.drainReleases()
	if !project.ValidID(requestID) || fingerprint == "" || len(fingerprint) > 256 || !catalog.ValidJourneyFilename(options.Filename) || (options.Intent != "create" && options.Intent != "replace") {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "invalidArgument"}
	}
	if freezer == nil {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "authorizationUnavailable"}
	}
	// state/release conservam revisão e geração durante codificação limitada, sem espera humana.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringFrozenReview{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "canceled"}
	}
	// prior consulta replay/CAS antes de validar conteúdo, sem adiantar geração em erro semântico.
	prior, replay, err := a.controlLookupLocked(state, requestID, fingerprint, expected, generation)
	if err != nil {
		return AuthoringFrozenReview{}, err
	}
	if replay && state.frozen != nil && state.frozen.ready && state.frozen.frozen.Ticket == prior {
		if state.frozen.fingerprint != fingerprint {
			return AuthoringFrozenReview{}, &AuthoringError{Code: "invalidArgument"}
		}
		return copyFrozenReview(state.frozen.frozen), nil
	}
	// A preparação pesada usa a reserva global de buffers, sem impedir consultas ou consumir identidade por busy.
	select {
	case a.finalizing <- struct{}{}:
		defer func() { <-a.finalizing }()
	default:
		return AuthoringFrozenReview{}, &AuthoringError{Code: "busy"}
	}
	// problems são calculados sem copiar imagens; rascunho inválido exige intenção explícita.
	problems := state.model.Document.Validate()
	if !options.AllowDraft {
		for _, problem := range problems {
			if problem.Severity == "error" {
				return AuthoringFrozenReview{}, &AuthoringError{Code: "invalidArgument"}
			}
		}
	}
	// ticket conserva prazo da primeira tentativa mesmo quando o adaptador precisa repetir IO após falha.
	ticket := prior
	if !replay {
		ticket, err = a.controlLocked(state, requestID, fingerprint, expected, generation, "review")
		if err != nil {
			return AuthoringFrozenReview{}, err
		}
	}
	if ticket.ExpiresAt <= a.now().Unix() {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "authorizationExpired"}
	}

	// A reserva existe antes de qualquer artefato novo, inclusive revisão de documento aberto limpo.
	a.mu.Lock()
	err = a.reserveLocked(state, "document")
	a.mu.Unlock()
	if err != nil {
		return AuthoringFrozenReview{}, err
	}
	if state.frozen != nil {
		if err := state.frozen.release(); err != nil {
			return AuthoringFrozenReview{}, &AuthoringError{Code: "persistenceFailed"}
		}
		state.frozen = nil
	}
	// request contém apenas aliases imutáveis sob a trava de sessão; o adaptador não os retém.
	request := AuthoringFreezeRequest{SessionID: id, Ticket: ticket, Revision: expected, Options: options, Baseline: state.baseline, Document: state.model.Document, BaselineDocument: state.base, Summary: reviewSummary(state.base, state.model.Document, options), Diagnostics: problems}
	frozen, freezeErr := freezer.Freeze(ctx, request)
	frozen.Ticket, frozen.Revision = ticket, expected
	if frozen.ArtifactID != "" {
		// artifactID evita capturar o documento ou os bytes JSON no callback de liberação.
		artifactID := frozen.ArtifactID
		state.frozen = &authoringReviewState{fingerprint: fingerprint, frozen: frozen, release: func() error { return freezer.Release(artifactID) }}
	}
	if freezeErr != nil {
		if ctx.Err() != nil {
			return AuthoringFrozenReview{}, &AuthoringError{Code: "canceled"}
		}
		// code preserva classificação segura do adaptador, sem copiar sua mensagem de IO.
		code := authoringCode(freezeErr)
		if code == "" {
			code = "persistenceFailed"
		}
		return AuthoringFrozenReview{}, &AuthoringError{Code: code}
	}
	if ctx.Err() != nil {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "canceled"}
	}
	// hash/size/JSON são conferidos antes de tornar qualquer revisão observável.
	contentHash, hashErr := hex.DecodeString(frozen.ContentHash)
	if state.frozen == nil || len(frozen.ArtifactID) > 128 || frozen.ContentSizeBytes < 1 || frozen.ContentSizeBytes > a.limits.TargetStagingBytes || hashErr != nil || len(contentHash) != 32 || hex.EncodeToString(contentHash) != frozen.ContentHash || len(frozen.Bytes) == 0 || len(frozen.Bytes) > 512<<10 || !utf8.Valid(frozen.Bytes) || !json.Valid(frozen.Bytes) {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "resourceLimit"}
	}
	// Bytes são copiados uma vez para retenção; o resultado devolvido tem sua própria cópia pequena.
	frozen.Bytes = append([]byte{}, frozen.Bytes...)
	digest := sha256.Sum256(frozen.Bytes)
	frozen.Hash = hex.EncodeToString(digest[:])
	state.frozen.frozen, state.frozen.ready = frozen, true
	return copyFrozenReview(frozen), nil
}

// CurrentReview retorna apenas a revisão corrente válida; bytes antigos não viram nova autorização.
func (a *Authoring) CurrentReview(id, reviewID string) (AuthoringFrozenReview, error) {
	// state/release protegem revisão e geração contra edição/descarte concorrentes.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringFrozenReview{}, err
	}
	defer release()
	if state.frozen == nil || !state.frozen.ready || state.frozen.frozen.Ticket.ReviewID != reviewID || state.frozen.frozen.Ticket.Generation != state.generation || state.frozen.frozen.Revision != state.model.Document.Project.Revision {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "staleReview"}
	}
	if state.frozen.frozen.Ticket.ExpiresAt <= a.now().Unix() {
		return AuthoringFrozenReview{}, &AuthoringError{Code: "authorizationExpired"}
	}
	return copyFrozenReview(state.frozen.frozen), nil
}

// copyFrozenReview impede alterar bytes retidos ou reapresentar um hash calculado sobre outro conteúdo.
func copyFrozenReview(frozen AuthoringFrozenReview) AuthoringFrozenReview {
	frozen.Bytes = append([]byte{}, frozen.Bytes...)
	return frozen
}

// retireFrozenLocked destaca o artefato sob os mutexes, conservando sua reserva até liberar fisicamente.
func (a *Authoring) retireFrozenLocked(state *authoringState) int64 {
	if state.frozen == nil {
		return 0
	}
	// held não libera staging lógico enquanto o handle/arquivo privado ainda existir.
	held := min(state.staging, a.limits.TargetStagingBytes)
	a.retirements = append(a.retirements, authoringRetirement{release: state.frozen.release, stage: held})
	a.retiring++
	state.frozen = nil
	return held
}

// drainReleases executa IO fora das travas globais; falha conserva reserva e permite tentar limpeza depois.
func (a *Authoring) drainReleases() {
	a.mu.Lock()
	// pending é um lote independente; outro chamador nunca libera o mesmo artefato simultaneamente.
	pending := a.retirements
	a.retirements = nil
	a.mu.Unlock()
	for _, retirement := range pending {
		// err distingue recurso efetivamente fechado de tentativa ainda não concluída.
		err := retirement.release()
		a.mu.Lock()
		if err == nil {
			a.staging -= retirement.stage
			a.retiring--
		} else {
			a.retirements = append(a.retirements, retirement)
		}
		a.mu.Unlock()
	}
}

// reviewSummary descreve mudanças a partir dos dados, sem interpretar intenções do legado ou da LLM.
func reviewSummary(base, current Document, options AuthoringReviewOptions) []string {
	// result sempre identifica operação e impactos agregados, mesmo para documentos extensos.
	result := []string{fmt.Sprintf("Operação %s para %s; nome do projeto (dado): %s.", options.Intent, options.Filename, reviewData(current.Project.Name)), fmt.Sprintf("Telas: %d → %d; recursos: %d → %d.", len(base.Screens), len(current.Screens), len(base.Project.Assets), len(current.Project.Assets))}
	changed, omitted := 0, 0
	for _, page := range current.Screens {
		// previous permite identificar inclusões e alterações sem copiar recursos binários.
		previous := base.Screen(page.ID)
		if previous == nil {
			result = append(result, fmt.Sprintf("Tela nova %s: %s (%d × %d).", page.ID, reviewData(page.Name), page.Size.Width, page.Size.Height))
		}
		for _, component := range page.Components {
			if previous != nil {
				if old := previous.ComponentByID(component.ID); old != nil && reflect.DeepEqual(*old, component) {
					continue
				}
			}
			changed++
			if len(result) >= 120 {
				omitted++
				continue
			}
			// encoded é dado tipado da instância; conteúdo grande recebe referência explícita e hash, sem truncar valor silenciosamente.
			encoded, _ := json.Marshal(component)
			digest := sha256.Sum256(encoded)
			if len(encoded) <= 2048 {
				result = append(result, fmt.Sprintf("Componente proposto %s na tela %s: %s", component.ID, page.ID, string(encoded)))
			} else {
				result = append(result, fmt.Sprintf("Componente proposto %s na tela %s: %d bytes, SHA-256 %x; conteúdo extenso deve ser consultado por ID na revisão %d antes de decidir.", component.ID, page.ID, len(encoded), digest, current.Project.Revision))
			}
		}
	}
	// Remoções são impactos próprios, inclusive quando inclusão e remoção conservam a contagem total.
	removed := 0
	for _, previous := range base.Screens {
		page := current.Screen(previous.ID)
		if page == nil {
			if len(result) < 120 {
				result = append(result, fmt.Sprintf("Tela removida: %s (%s).", previous.ID, reviewData(previous.Name)))
			}
			removed += len(previous.Components)
			continue
		}
		for _, component := range previous.Components {
			if page.ComponentByID(component.ID) == nil {
				removed++
				if len(result) < 120 {
					result = append(result, fmt.Sprintf("Componente removido: %s na tela %s.", component.ID, previous.ID))
				}
			}
		}
	}
	if !reflect.DeepEqual(base.Journey, current.Journey) {
		// encoded/digest identificam navegação explícita; valores grandes são resumidos com indicação clara.
		encoded, _ := json.Marshal(current.Journey)
		digest := sha256.Sum256(encoded)
		if len(encoded) <= 2048 {
			result = append(result, "Navegação proposta (dados): "+string(encoded))
		} else {
			result = append(result, fmt.Sprintf("Navegação proposta: %d vínculos; SHA-256 %x. Consultar os destinos explícitos dos componentes na mesma revisão antes de decidir.", len(current.Journey.Actions), digest))
		}
	}
	result = append(result, fmt.Sprintf("Componentes removidos: %d.", removed))
	result = append(result, fmt.Sprintf("Componentes incluídos/alterados: %d; detalhes adicionais fora do resumo: %d. A revisão vincula o hash do pacote completo; consulta de conteúdo deve conferir a mesma revisão.", changed, omitted))
	return result
}

// reviewData identifica explicitamente prévias de valores históricos longos, sem alterá-los no documento.
func reviewData(value string) string {
	if utf8.RuneCountInString(value) <= 512 {
		return strconv.Quote(value)
	}
	// prefix/digest mantêm resumo limitado e vínculo ao valor completo, sem alegar exibição integral.
	prefix := []rune(value)[:128]
	digest := sha256.Sum256([]byte(value))
	return fmt.Sprintf("prévia explícita %s; total %d pontos de código; SHA-256 %x", strconv.Quote(string(prefix)), utf8.RuneCountInString(value), digest)
}
