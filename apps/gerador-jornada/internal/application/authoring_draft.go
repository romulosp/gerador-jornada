package application

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"gerador-jornada.local/editor/internal/project"
)

// AuthoringDraft conserva bytes de proposta validados pelo caso de uso específico.
// O registro coordena ownership/CAS; contratos de legado/captura continuam em seus casos de uso.
type AuthoringDraft struct {
	// ID/Revision identificam o anexo e sua revisão int32 independente do documento.
	ID       string
	Revision int
	// BaseRevision liga validação/revisão à base do documento sem alterar Project.Revision.
	BaseRevision uint64
	// Origin distingue selection, import e legacy sem criar outro registro de sessões.
	Origin string
	// Data contém metadados/descrição estruturados; pacotes grandes usam staging reservado.
	Data []byte
	// SHA256 é calculado pela aplicação sobre os bytes retidos, nunca confiado ao agente.
	SHA256 string
	// Target é identidade lógica persistível; análise legacy permanece transitória e usa vazio.
	Target string
	// Dirty diferencia anexo pendente de captura já publicada, sem apagar a análise vinculada.
	Dirty bool
}

// StoreDraft admite/atualiza anexo da sessão por CAS, consumindo edição e orçamento.
// Data deve ser JSON já validado pelo caso de uso; este método não interpreta seu conteúdo.
func (a *Authoring) StoreDraft(id, requestID, fingerprint string, expected uint64, draft AuthoringDraft) (AuthoringDraft, error) {
	if !project.ValidID(requestID) || !project.ValidID(draft.ID) || fingerprint == "" || len(fingerprint) > 256 {
		return AuthoringDraft{}, &AuthoringError{Code: "invalidArgument"}
	}
	// state/release protegem ownership, revisão, bytes e reserva em uma única admissão.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringDraft{}, err
	}
	defer release()
	if prior, ok := state.edits[requestID]; ok {
		if prior.kind != "draft" || prior.digest != fingerprint {
			return AuthoringDraft{}, &AuthoringError{Code: "invalidArgument"}
		}
		if prior.code != "" {
			return AuthoringDraft{}, &AuthoringError{Code: prior.code}
		}
		// revision/base são a resposta original, mesmo que o anexo tenha nova versão depois.
		revision, _ := strconv.Atoi(prior.result.Values["revision"])
		base, _ := strconv.ParseUint(prior.result.Values["baseRevision"], 10, 64)
		return AuthoringDraft{ID: prior.result.Values["id"], Revision: revision, BaseRevision: base, Origin: prior.result.Values["origin"], SHA256: prior.result.Values["sha256"], Target: prior.result.Values["target"], Dirty: true}, nil
	}
	if len(state.edits) >= a.limits.EditRequests || state.consumed >= a.limits.Commits {
		return AuthoringDraft{}, &AuthoringError{Code: "resourceLimit"}
	}
	if state.pendingIntent() {
		return AuthoringDraft{}, &AuthoringError{Code: "recoveryRequired"}
	}
	// reject conserva falha identificada dentro da reserva de metadados da sessão.
	reject := func(code string) (AuthoringDraft, error) {
		state.edits[requestID] = authoringRequest{kind: "draft", digest: fingerprint, code: code}
		return AuthoringDraft{}, &AuthoringError{Code: code}
	}
	if expected != state.model.Document.Project.Revision {
		return reject("staleReview")
	}
	if draft.Origin != "selection" && draft.Origin != "import" && draft.Origin != "legacy" {
		return reject("invalidArgument")
	}
	if (draft.Origin == "legacy" && draft.Target != "") || len(draft.Target) > 256 || len(draft.Data) > 1<<20 || !json.Valid(draft.Data) {
		return reject("invalidArgument")
	}
	// previous/existing determinam CAS e delta de memória antes da mutação do registro.
	previous, existing := state.drafts[draft.ID]
	if (!existing && draft.Revision != 1) || (existing && (previous.Revision == math.MaxInt32 || draft.Revision != previous.Revision+1 || draft.Origin != previous.Origin || draft.Target != previous.Target)) {
		return reject("staleReview")
	}
	if !existing && len(state.drafts) >= a.limits.Proposals {
		return reject("resourceLimit")
	}
	if state.generation == math.MaxUint64 {
		return reject("resourceLimit")
	}
	draft.BaseRevision = expected
	draft.Dirty = true
	// digest é sobre bytes completos do anexo, não resumo de evidência ou prompt reescrito.
	digest := sha256.Sum256(draft.Data)
	draft.SHA256 = fmt.Sprintf("%x", digest)
	// delta/cost contam bytes do anexo e resposta idempotente sem liberar histórico de IDs.
	delta := int64(len(draft.Data) - len(previous.Data))
	if !existing {
		delta += 512
	}
	cost := int64(1024)
	a.mu.Lock()
	if a.retained+delta+cost > a.limits.RetainedBytes {
		a.mu.Unlock()
		return reject("resourceLimit")
	}
	if draft.Target != "" {
		if err := a.reserveLocked(state, draft.Target); err != nil {
			a.mu.Unlock()
			return reject("resourceLimit")
		}
	}
	a.retained += delta + cost
	state.retained += delta + cost
	state.cacheBytes += cost
	state.draftBytes += delta
	a.mu.Unlock()
	draft.Data = append([]byte(nil), draft.Data...)
	state.drafts[draft.ID] = draft
	state.proposalDirty = true
	state.generation++
	state.recent = nil
	state.edits[requestID] = authoringRequest{kind: "draft", digest: fingerprint, result: AuthoringEditResult{Revision: expected, Values: map[string]string{"id": draft.ID, "revision": strconv.Itoa(draft.Revision), "baseRevision": strconv.FormatUint(expected, 10), "origin": draft.Origin, "sha256": draft.SHA256, "target": draft.Target}}}
	// response conserva somente identidade/revisão/digest; conteúdo é consultado por Draft.
	draft.Data = nil
	return draft, nil
}

// Draft devolve cópia independente de um anexo pertencente à sessão informada.
func (a *Authoring) Draft(id, proposalID string) (AuthoringDraft, error) {
	// state/release não permitem consultar anexo de outra sessão pelo UUID isolado.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringDraft{}, err
	}
	defer release()
	// draft/ok conservam o erro localizado sem buscar em qualquer outra sessão.
	draft, ok := state.drafts[proposalID]
	if !ok {
		return AuthoringDraft{}, &AuthoringError{Code: "invalidReference"}
	}
	draft.Data = append([]byte(nil), draft.Data...)
	return draft, nil
}
