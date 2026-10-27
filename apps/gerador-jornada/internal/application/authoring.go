package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/project"
)

// AuthoringLimits fixa quotas locais; não é configuração recebida do agente.
type AuthoringLimits struct {
	// Sessions limita documentos simultâneos no processo.
	Sessions int
	// ConcurrentCalls limita acesso concorrente aos snapshots do registro.
	ConcurrentCalls int
	// EditRequests limita identidades de edição lembradas, sem eviction.
	EditRequests int
	// Proposals limita anexos transitórios por sessão.
	Proposals int
	// Controls limita o histórico recente de controles, além do high-water.
	Controls int
	// Commits limita posições consumidas, incertas ou reservadas por sessão.
	Commits int
	// RetainedBytes inclui documentos, histórico, cache e reserva de metadados.
	RetainedBytes int64
	// FinalizationReserveBytes inclui revisão e metadados por sessão antes da admissão; buffers pesados usam a reserva global.
	FinalizationReserveBytes int64
	// FinalizationWorkBytes reserva o buffer global de uma finalização, sem fila ilimitada ou clone de ZIP.
	FinalizationWorkBytes int64
	// StagingBytes limita reservas lógicas de finalização no processo.
	StagingBytes int64
	// TargetStagingBytes reserva o limite de pacote por alvo persistível.
	TargetStagingBytes int64
	// Idle e ReviewTTL separam inatividade de validade de revisão.
	Idle, ReviewTTL time.Duration
}

// DefaultAuthoringLimits corresponde aos limites iniciais do contrato 038.
func DefaultAuthoringLimits() AuthoringLimits {
	return AuthoringLimits{Sessions: 8, ConcurrentCalls: 4, EditRequests: 1024, Proposals: 32, Controls: 16, Commits: 64,
		RetainedBytes: 256 << 20, FinalizationReserveBytes: 1 << 20, FinalizationWorkBytes: 64 << 20, StagingBytes: 1 << 30, TargetStagingBytes: 250 << 20, Idle: 30 * time.Minute, ReviewTTL: 5 * time.Minute}
}

// AuthoringError fornece código estável sem transcrever conteúdo do documento.
type AuthoringError struct {
	// Code distingue contrato, quota, expiração, conflito e concorrência.
	Code string
}

// Error conserva uma mensagem técnica curta, sem nomes ou textos do projeto.
func (e *AuthoringError) Error() string { return "autoria: " + e.Code }

// authoringCode extrai somente o código seguro de um erro de aplicação.
func authoringCode(err error) string {
	// typed identifica um erro de autoria mesmo quando envolvido por outro erro.
	var typed *AuthoringError
	if errors.As(err, &typed) {
		return typed.Code
	}
	return ""
}

// AuthoringBaseline descreve a base fornecida pelo adaptador, sem fazer IO.
type AuthoringBaseline struct {
	// Exists diferencia abertura de arquivo e rascunho novo ainda não persistido.
	Exists bool
	// Identity, SHA256 e SizeBytes serão revalidados pelo adaptador sob gate.
	Identity, SHA256 string
	SizeBytes        int64
	// RootIdentity e Filename vinculam a origem lógica, sem aceitar caminho arbitrário.
	RootIdentity, Filename string
}

// AuthoringSessionInfo é o resultado pequeno e idempotente da admissão.
type AuthoringSessionInfo struct {
	// ID não é caminho e permanece tombstone depois de fechar/expirar.
	ID string
	// Revision é a revisão inicial; consulta posterior usa Inspect.
	Revision uint64
	// Name conserva o resumo original para replay sem reler origem ou estado alterado.
	Name string
	// ExpiresAt informa o prazo inicial; replay da admissão conserva esse resultado sem renovar a sessão.
	ExpiresAt time.Time
}

// AuthoringSnapshot contém uma cópia independente e quotas observáveis.
type AuthoringSnapshot struct {
	// ID e Revision identificam sessão e revisão atual do documento.
	ID       string
	Revision uint64
	// Document é independente do estado mutável mantido pela aplicação.
	Document Document
	// DocumentDirty e ProposalDirty conservam pendências de origens distintas.
	DocumentDirty, ProposalDirty, Dirty bool
	// ControlGeneration identifica a próxima preparação deliberada de controle.
	ControlGeneration uint64
	// EditRequests e CommitPositions informam consumo sem liberar a quota por replay.
	EditRequests, CommitPositions int
	// ReservedTargets/Proposals/ConsumedCommits/UncertainIntents separam cotas e reservas pendentes.
	ReservedTargets, Proposals, ConsumedCommits, UncertainIntents int
	// StagingReservedBytes não inclui aposentadorias globais já destacadas da sessão.
	StagingReservedBytes int64
	// RecoveryRequired identifica resultado incerto que deve ser reconciliado antes de nova edição/encerramento.
	RecoveryRequired bool
	// ExpiresAt informa o prazo de inatividade da sessão.
	ExpiresAt time.Time
}

// AuthoringEditResult permite devolver IDs remapeados sem DTO MCP no domínio.
type AuthoringEditResult struct {
	// Revision é a revisão confirmada pela aplicação.
	Revision uint64
	// Values contém somente resultados pequenos, como identidades criadas.
	Values map[string]string
}

// AuthoringControlTicket identifica preparação, nunca aprovação de persistência.
type AuthoringControlTicket struct {
	// Revision/Closed/Discarded conservam o resultado do controle sem consultar estado posterior no replay.
	Revision          uint64
	Closed, Discarded bool
	// Generation avança em edição ou controle novos, sem reset de quota.
	Generation uint64
	// ReviewID e CommitRequestID serão vinculados ao conteúdo congelado pelo caso de uso.
	ReviewID, CommitRequestID string
	// ExpiresAt é fixo e não se renova quando o cliente repete a solicitação.
	ExpiresAt int64
}

// authoringRequest mantém digest e resultado da primeira edição identificada.
type authoringRequest struct {
	// kind distingue edição de documento e anexo mesmo se o chamador reutilizar digest.
	kind string
	// digest não é recalculado a partir de estado que mudou depois da edição.
	digest string
	// result/code preservam sucesso ou falha sem executar o callback novamente.
	result AuthoringEditResult
	code   string
}

// authoringControl guarda a identidade recebida, anterior ao avanço da geração.
type authoringControl struct {
	// requestID/digest/generation distinguem replay de novo comando.
	requestID, digest string
	generation        uint64
	// ticket é o resultado imutável da preparação ainda corrente.
	ticket AuthoringControlTicket
}

// authoringAdmission permite repetir criação/abertura sem criar outro documento.
type authoringAdmission struct {
	// digest vincula a identidade aos argumentos originais.
	digest string
	// info conserva a resposta inicial mesmo que a sessão depois seja encerrada.
	info AuthoringSessionInfo
}

// authoringIntent ocupa a mesma posição desde a intenção até sua reconciliação.
type authoringIntent struct {
	// target é identidade lógica, nunca nome de arquivo interpretado pelo domínio.
	target string
	// resolved impede produzir novo efeito com uma identidade já finalizada.
	resolved bool
}

// authoringState protege um histórico de domínio e metadados externos ao manifesto.
type authoringState struct {
	// mu permite uma operação por sessão; concorrência incompatível retorna busy.
	mu sync.Mutex
	// id/model/base conservam identidade, edição atual e estado para descarte.
	id    string
	model *Session
	base  Document
	// baseline não concede autoridade para gravar; somente descreve a base lida.
	baseline AuthoringBaseline
	// unsaved/proposalDirty compõem dirty sem mudar Project.Revision por um anexo.
	unsaved, proposalDirty bool
	// closed/expired impedem recriação silenciosa por replay.
	closed, expired bool
	// closedControl guarda somente o resultado terminal do fechamento, sem documento ou artefato.
	closedControl *authoringControl
	// touched controla inatividade separadamente da validade da revisão.
	touched time.Time
	// edits conserva todas as identidades de edição aceitas até encerrar.
	edits map[string]authoringRequest
	// generation/recent conservam high-water e histórico de controle limitado.
	generation uint64
	recent     []authoringControl
	// reserved/intents/consumed preservam o invariante de finalização.
	reserved map[string]bool
	intents  map[string]authoringIntent
	consumed int
	// retained/staging mantêm a parte desta sessão nas quotas do processo.
	retained, staging int64
	// cacheBytes conta respostas de sucesso além da reserva de metadados fixos.
	cacheBytes int64
	// drafts/draftBytes mantêm anexos da própria sessão fora do manifesto e do undo.
	drafts     map[string]AuthoringDraft
	draftBytes int64
	// frozen conserva artefato e revisão anexados à sessão, sem registro paralelo MCP.
	frozen *authoringReviewState
}

// Authoring coordena sessões dos adaptadores sem importar Fyne, MCP ou filesystem.
type Authoring struct {
	// mu protege registro/quotas; nunca espera sessão bloqueada enquanto está retido.
	mu sync.Mutex
	// limits/now são fixados pelo operador de composição, não por ferramentas do agente.
	limits AuthoringLimits
	now    func() time.Time
	// active limita chamadas concorrentes sem criar filas de trabalho ilimitadas.
	active chan struct{}
	// finalizing permite um encode/verify por vez, mantendo consultas e demais comandos disponíveis.
	finalizing chan struct{}
	// sessions/admissions mantêm identidades e idempotência antes/depois da admissão.
	sessions   map[string]*authoringState
	admissions map[string]authoringAdmission
	// retained/staging contabilizam quotas lógicas, sem prometer heap Go exato.
	retained, staging int64
	// retirements mantém staging cobrado até liberação física fora da trava global.
	retirements []authoringRetirement
	// retiring conta também o lote que já saiu da fila e ainda está fechando handles fora das travas.
	retiring int
}

// NewAuthoring cria um registro com limites válidos e relógio substituível em testes.
// Limites inválidos indicam erro de composição, nunca entrada de cliente MCP.
func NewAuthoring(limits AuthoringLimits, now func() time.Time) *Authoring {
	if limits.Sessions < 1 || limits.ConcurrentCalls < 1 || limits.EditRequests < 1 || limits.Proposals < 1 || limits.Controls < 1 || limits.Commits < 1 || limits.RetainedBytes <= 0 || limits.FinalizationReserveBytes <= 0 || limits.FinalizationWorkBytes <= 0 || limits.FinalizationWorkBytes >= limits.RetainedBytes || limits.TargetStagingBytes <= 0 || limits.StagingBytes <= 0 || limits.Idle <= 0 || limits.ReviewTTL <= 0 {
		panic("limites de autoria inválidos")
	}
	if now == nil {
		now = time.Now
	}
	return &Authoring{limits: limits, now: now, active: make(chan struct{}, limits.ConcurrentCalls), finalizing: make(chan struct{}, 1), sessions: map[string]*authoringState{}, admissions: map[string]authoringAdmission{}, retained: limits.FinalizationWorkBytes}
}

// Admit cria snapshot isolado; replay com argumentos iguais devolve a resposta inicial.
// fingerprint é digest calculado pela aplicação/adaptador sobre todos os argumentos.
func (a *Authoring) Admit(requestID, fingerprint string, document Document, baseline AuthoringBaseline) (AuthoringSessionInfo, error) {
	defer a.drainReleases()
	if !project.ValidID(requestID) || fingerprint == "" || len(fingerprint) > 256 {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "invalidArgument"}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if prior, ok := a.admissions[requestID]; ok {
		if prior.digest != fingerprint {
			return AuthoringSessionInfo{}, &AuthoringError{Code: "invalidArgument"}
		}
		return prior.info, nil
	}
	// count desconta sessões inativas sem bloquear uma mutação concorrente.
	count := 0
	// state pertence ao registro e só é inspecionado depois de TryLock.
	for _, state := range a.sessions {
		if !state.mu.TryLock() {
			count++
			continue
		}
		if !state.closed && !state.expired && !state.pendingIntent() && a.now().Sub(state.touched) >= a.limits.Idle {
			a.retireLocked(state, true)
		}
		if !state.closed && !state.expired {
			count++
		}
		state.mu.Unlock()
	}
	if count >= a.limits.Sessions {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "resourceLimit"}
	}
	if err := document.ValidateFileIDs(); err != nil {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "invalidArgument"}
	}
	// model/base são cópias independentes; a reserva de metadados protege finalização.
	model := NewSession(document)
	base := document.Copy()
	// retained soma conteúdo e metadados, incluindo cache inicial da admissão.
	retained := authoringSessionWeight(model, base) + a.limits.FinalizationReserveBytes
	// admissionBytes conserva também o nome original no tombstone de replay.
	admissionBytes := int64(512 + len(document.Project.Name))
	if a.retained+retained+admissionBytes > a.limits.RetainedBytes {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "resourceLimit"}
	}
	// stage reserva o pacote de um documento novo antes de admitir seu trabalho.
	stage := int64(0)
	if !baseline.Exists {
		stage = a.limits.TargetStagingBytes
	}
	if a.staging+stage > a.limits.StagingBytes {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "resourceLimit"}
	}
	// info é a única resposta conservada no cache de criação por processo.
	info := AuthoringSessionInfo{ID: project.NewID(), Revision: model.Document.Project.Revision, Name: model.Document.Project.Name, ExpiresAt: a.now().Add(a.limits.Idle)}
	// state não possui referência a documento/widget do chamador.
	state := &authoringState{id: info.ID, model: model, base: base, baseline: baseline, unsaved: !baseline.Exists, touched: a.now(), edits: map[string]authoringRequest{}, reserved: map[string]bool{}, intents: map[string]authoringIntent{}, drafts: map[string]AuthoringDraft{}, retained: retained, staging: stage}
	if !baseline.Exists {
		state.reserved["document"] = true
	}
	a.sessions[info.ID] = state
	a.admissions[requestID] = authoringAdmission{digest: fingerprint, info: info}
	a.retained += retained + admissionBytes
	a.staging += stage
	return info, nil
}

// acquire obtém registro e sessão sem inverter ordem de locks ou esperar indefinidamente.
func (a *Authoring) acquire(id string) (*authoringState, func(), error) {
	select {
	case a.active <- struct{}{}:
	default:
		return nil, nil, &AuthoringError{Code: "busy"}
	}
	a.mu.Lock()
	state := a.sessions[id]
	a.mu.Unlock()
	if state == nil {
		<-a.active
		return nil, nil, &AuthoringError{Code: "sessionExpired"}
	}
	if !state.mu.TryLock() {
		<-a.active
		return nil, nil, &AuthoringError{Code: "busy"}
	}
	// release libera a sessão antes de devolver a vaga global.
	release := func() { state.mu.Unlock(); <-a.active }
	if state.closed || state.expired {
		release()
		return nil, nil, &AuthoringError{Code: "sessionExpired"}
	}
	if a.now().Sub(state.touched) >= a.limits.Idle && !state.pendingIntent() {
		a.mu.Lock()
		a.retireLocked(state, true)
		a.mu.Unlock()
		release()
		a.drainReleases()
		return nil, nil, &AuthoringError{Code: "sessionExpired"}
	}
	// Uma intenção incerta após inatividade continua somente para consulta/reconciliação, sem renovar a sessão.
	if a.now().Sub(state.touched) < a.limits.Idle {
		state.touched = a.now()
	}
	return state, release, nil
}

// Inspect devolve snapshot independente e não mantém lock depois da resposta.
func (a *Authoring) Inspect(id string) (AuthoringSnapshot, error) {
	return a.inspect(id, true, true)
}

// InspectMetadata copia o documento sem bytes de recursos para consultas que não os transportam.
// Metadados, coleções e revisões permanecem independentes, sem duplicar imagens para responder R0.
func (a *Authoring) InspectMetadata(id string) (AuthoringSnapshot, error) {
	return a.inspect(id, true, false)
}

// SessionStatus consulta ciclo de vida e quotas sem copiar documento, histórico ou recursos.
func (a *Authoring) SessionStatus(id string) (AuthoringSnapshot, error) {
	return a.inspect(id, false, false)
}

// inspect constrói a cópia solicitada mantendo a mesma proteção/expiração de sessão.
func (a *Authoring) inspect(id string, includeDocument, includeResources bool) (AuthoringSnapshot, error) {
	// state/release conservam uma única operação sobre a sessão em consulta.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringSnapshot{}, err
	}
	defer release()
	// dirty combina documento novo/alterado e anexos pendentes sem misturá-los.
	dirty := state.unsaved || state.model.Dirty()
	// document remove bytes antes de copiar; não altera o mapa de recursos mantido pela sessão.
	document := Document{}
	if includeDocument {
		document = state.model.Document
		if !includeResources {
			document.Resources = nil
		}
		document = document.Copy()
	}
	// uncertain conta somente intents sem resultado comprovado, sem duplicar posições resolvidas.
	uncertain := 0
	for _, intent := range state.intents {
		if !intent.resolved {
			uncertain++
		}
	}
	return AuthoringSnapshot{ID: state.id, Revision: state.model.Document.Project.Revision, Document: document, DocumentDirty: dirty, ProposalDirty: state.proposalDirty, Dirty: dirty || state.proposalDirty, ControlGeneration: state.generation, EditRequests: len(state.edits), CommitPositions: state.positions(), ReservedTargets: len(state.reserved), Proposals: len(state.drafts), ConsumedCommits: state.consumed, UncertainIntents: uncertain, StagingReservedBytes: state.staging, RecoveryRequired: state.pendingIntent(), ExpiresAt: state.touched.Add(a.limits.Idle)}, nil
}

// ExpireIdle aposenta sessões ociosas sem esperar uma mutação que esteja em andamento.
// Intents incertos permanecem contabilizados para reconciliação e não ganham prazo novo por consulta.
func (a *Authoring) ExpireIdle() {
	a.mu.Lock()
	// states permite observar cada sessão sem bloquear sua trava sob o mutex global.
	states := make([]*authoringState, 0, len(a.sessions))
	for _, state := range a.sessions {
		states = append(states, state)
	}
	a.mu.Unlock()
	for _, state := range states {
		if !state.mu.TryLock() {
			continue
		}
		if !state.closed && !state.expired && !state.pendingIntent() && a.now().Sub(state.touched) >= a.limits.Idle {
			a.mu.Lock()
			a.retireLocked(state, true)
			a.mu.Unlock()
		}
		state.mu.Unlock()
	}
	a.drainReleases()
}

// AuthoringUsage descreve reservas operacionais, sem documentos, segredos ou caminhos.
type AuthoringUsage struct {
	// RetainedBytes inclui buffers globais reservados, conteúdo, histórico, caches e metadados de sessão.
	RetainedBytes int64
	// StagingReservedBytes inclui aposentadorias cuja liberação física ainda não terminou.
	StagingReservedBytes int64
	// PendingCleanups/ActiveCalls/FinalizationInProgress não concedem permissões adicionais.
	PendingCleanups, ActiveCalls int
	FinalizationInProgress       bool
}

// Usage oferece visão segura do processo para monitoramento de cotas e limpeza, sem esperar IO.
func (a *Authoring) Usage() AuthoringUsage {
	a.mu.Lock()
	defer a.mu.Unlock()
	return AuthoringUsage{RetainedBytes: a.retained, StagingReservedBytes: a.staging, PendingCleanups: a.retiring, ActiveCalls: len(a.active), FinalizationInProgress: len(a.finalizing) != 0}
}

// Edit aplica exatamente um comando Session, preservando conteúdo/undo em falha.
// apply recebe somente estado candidato e deve usar os comandos atômicos Session.
func (a *Authoring) Edit(id, requestID, fingerprint string, expected uint64, apply func(*Session) (map[string]string, error)) (AuthoringEditResult, error) {
	return a.EditContext(context.Background(), id, requestID, fingerprint, expected, apply)
}

// EditContext impede publicar um candidato cancelado; replay mantém a mesma identidade.
func (a *Authoring) EditContext(ctx context.Context, id, requestID, fingerprint string, expected uint64, apply func(*Session) (map[string]string, error)) (AuthoringEditResult, error) {
	if !project.ValidID(requestID) || fingerprint == "" || len(fingerprint) > 256 {
		return AuthoringEditResult{}, &AuthoringError{Code: "invalidArgument"}
	}
	// state/release mantêm o mutex apenas durante trabalho local, nunca interação humana.
	state, release, err := a.acquire(id)
	if err != nil {
		return AuthoringEditResult{}, err
	}
	defer release()
	if cached, ok := state.edits[requestID]; ok {
		if cached.digest != fingerprint || cached.kind != "edit" {
			return AuthoringEditResult{}, &AuthoringError{Code: "invalidArgument"}
		}
		if cached.code != "" {
			return AuthoringEditResult{}, &AuthoringError{Code: cached.code}
		}
		return copyAuthoringResult(cached.result), nil
	}
	if len(state.edits) >= a.limits.EditRequests || state.consumed >= a.limits.Commits {
		return AuthoringEditResult{}, &AuthoringError{Code: "resourceLimit"}
	}
	if state.pendingIntent() {
		return AuthoringEditResult{}, &AuthoringError{Code: "recoveryRequired"}
	}
	// fail recorda a identidade sem conservar mensagem do domínio ou publicar candidato.
	fail := func(code string) (AuthoringEditResult, error) {
		state.edits[requestID] = authoringRequest{kind: "edit", digest: fingerprint, code: code}
		return AuthoringEditResult{}, &AuthoringError{Code: code}
	}
	if expected != state.model.Document.Project.Revision {
		return fail("staleReview")
	}
	if ctx.Err() != nil {
		return fail("canceled")
	}
	if apply == nil {
		return fail("invalidArgument")
	}
	if state.model.Sequence == math.MaxUint64 || state.generation == math.MaxUint64 {
		return fail("resourceLimit")
	}
	// candidate separa o documento; snapshots históricos existentes permanecem imutáveis.
	candidate := *state.model
	candidate.Document = state.model.Document.Copy()
	candidate.UndoStack = append([]Document(nil), state.model.UndoStack...)
	candidate.RedoStack = append([]Document(nil), state.model.RedoStack...)
	// values reúne IDs pequenos; erro nunca preserva mutação parcial do callback.
	values, err := apply(&candidate)
	if err != nil {
		// code conserva somente classificação tipada; texto do erro de domínio não entra no cache.
		code := authoringCode(err)
		if code == "" {
			code = "invalidArgument"
		}
		return fail(code)
	}
	if ctx.Err() != nil {
		return fail("canceled")
	}
	if candidate.Sequence != state.model.Sequence+1 || candidate.Document.Project.Revision != candidate.Sequence || candidate.SavedRevision != state.model.SavedRevision {
		return fail("invalidArgument")
	}
	if len(values) > 65536 {
		return fail("resourceLimit")
	}
	// encodedValues limita o resultado antes de publicar a mutação, inclusive mapas grandes de remapeamento.
	encodedValues, encodeErr := json.Marshal(values)
	if encodeErr != nil || len(encodedValues) > 512<<10 {
		return fail("resourceLimit")
	}
	// key/value são resultados de aplicação, limitados antes de entrar no cache.
	for key, value := range values {
		if len(key) > 128 || len(value) > 256 {
			return fail("resourceLimit")
		}
	}
	// weight inclui todo o histórico retido, recursos e reserva de finalização/cache.
	// cacheCost inclui chaves/valores retidos para replay, sem eviction silenciosa.
	cacheCost := int64(len(requestID) + len(fingerprint) + 128)
	for key, value := range values {
		cacheCost += int64(len(key) + len(value) + 32)
	}
	weight := authoringSessionWeight(&candidate, state.base) + a.limits.FinalizationReserveBytes + state.cacheBytes + cacheCost + state.draftBytes
	a.mu.Lock()
	if a.retained-state.retained+weight > a.limits.RetainedBytes {
		a.mu.Unlock()
		return fail("resourceLimit")
	}
	if err := a.reserveLocked(state, "document"); err != nil {
		a.mu.Unlock()
		return fail("resourceLimit")
	}
	a.retained += weight - state.retained
	state.retained = weight
	a.mu.Unlock()
	state.model = &candidate
	state.cacheBytes += cacheCost
	state.generation++
	state.recent = nil
	// result é copiado tanto para cache quanto para resposta, sem alias de mapas.
	result := AuthoringEditResult{Revision: candidate.Sequence, Values: values}
	state.edits[requestID] = authoringRequest{kind: "edit", digest: fingerprint, result: copyAuthoringResult(result)}
	return copyAuthoringResult(result), nil
}

// Control prepara identidade de revisão ou executa descarte/fechamento explícitos.
// A identidade retornada não autoriza IO e a revisão final será congelada pelo caso de uso.
func (a *Authoring) Control(id, requestID, fingerprint string, expected, generation uint64, operation string) (AuthoringControlTicket, error) {
	if !project.ValidID(requestID) || fingerprint == "" || len(fingerprint) > 256 {
		return AuthoringControlTicket{}, &AuthoringError{Code: "invalidArgument"}
	}
	// state/release protegem a geração e o cache corrente de controle.
	state, release, err := a.acquire(id)
	if err != nil {
		if authoringCode(err) == "sessionExpired" {
			return a.closedControlReplay(id, requestID, fingerprint, generation, operation)
		}
		return AuthoringControlTicket{}, err
	}
	defer a.drainReleases()
	defer release()
	return a.controlLocked(state, requestID, fingerprint, expected, generation, operation)
}

// controlLocked aplica controle com a sessão já protegida, compartilhando replay com PrepareReview.
func (a *Authoring) controlLocked(state *authoringState, requestID, fingerprint string, expected, generation uint64, operation string) (AuthoringControlTicket, error) {
	// prior devolve replay antes de comparar high-water ou validar conteúdo posterior.
	prior, found, err := a.controlLookupLocked(state, requestID, fingerprint, expected, generation)
	if err != nil || found {
		return prior, err
	}
	if operation != "review" && operation != "discard" && operation != "close" && operation != "closeDiscard" {
		return AuthoringControlTicket{}, &AuthoringError{Code: "invalidArgument"}
	}
	if operation == "close" && (state.unsaved || state.model.Dirty() || state.proposalDirty) {
		return AuthoringControlTicket{}, &AuthoringError{Code: "invalidArgument"}
	}
	if state.pendingIntent() {
		return AuthoringControlTicket{}, &AuthoringError{Code: "recoveryRequired"}
	}
	if operation == "discard" && state.model.Sequence == math.MaxUint64 {
		return AuthoringControlTicket{}, &AuthoringError{Code: "resourceLimit"}
	}
	if state.generation == math.MaxUint64 {
		return AuthoringControlTicket{}, &AuthoringError{Code: "resourceLimit"}
	}
	state.generation++
	// ticket possui IDs exclusivos para revisão/commit, diferentes da preparação.
	ticket := AuthoringControlTicket{Revision: state.model.Document.Project.Revision, Generation: state.generation, ReviewID: project.NewID(), CommitRequestID: project.NewID(), ExpiresAt: a.now().Add(a.limits.ReviewTTL).Unix()}
	state.recent = append(state.recent, authoringControl{requestID: requestID, digest: fingerprint, generation: generation, ticket: ticket})
	if len(state.recent) > a.limits.Controls {
		state.recent = state.recent[len(state.recent)-a.limits.Controls:]
	}
	if operation == "discard" {
		// restored conserva Sequence monotônica, impedindo reutilizar revisão antiga após descarte.
		restored := NewSession(state.base)
		restored.Sequence = state.model.Sequence + 1
		restored.Document.Project.Revision = restored.Sequence
		restored.SavedRevision = restored.Sequence
		state.model = restored
		ticket.Revision, ticket.Discarded = restored.Sequence, true
		state.recent[len(state.recent)-1].ticket = ticket
		state.unsaved = false
		state.proposalDirty = false
		state.drafts = map[string]AuthoringDraft{}
		state.draftBytes = 0
		a.mu.Lock()
		// held permanece no processo enquanto o staging está em liberação.
		held := a.retireFrozenLocked(state)
		a.staging -= state.staging - held
		state.staging = 0
		state.reserved = map[string]bool{}
		// weight devolve também o orçamento do histórico descartado, mantendo cache/controle.
		weight := authoringSessionWeight(restored, state.base) + a.limits.FinalizationReserveBytes + state.cacheBytes
		a.retained += weight - state.retained
		state.retained = weight
		a.mu.Unlock()
	}
	if operation == "close" || operation == "closeDiscard" {
		ticket.Closed, ticket.Discarded = true, operation == "closeDiscard"
		// closedRecord conserva a resposta terminal depois que retireLocked libera o conteúdo da sessão.
		closedRecord := authoringControl{requestID: requestID, digest: fingerprint, generation: generation, ticket: ticket}
		a.mu.Lock()
		a.retireLocked(state, false)
		a.mu.Unlock()
		state.closedControl = &closedRecord
	}
	return ticket, nil
}

// closedControlReplay devolve somente resultado terminal conhecido, sem reabrir ou ressuscitar uma sessão.
func (a *Authoring) closedControlReplay(id, requestID, fingerprint string, generation uint64, operation string) (AuthoringControlTicket, error) {
	a.mu.Lock()
	// state é tombstone da mesma identidade de sessão; não concede uma origem nova.
	state := a.sessions[id]
	a.mu.Unlock()
	if state == nil {
		return AuthoringControlTicket{}, &AuthoringError{Code: "sessionExpired"}
	}
	if !state.mu.TryLock() {
		return AuthoringControlTicket{}, &AuthoringError{Code: "busy"}
	}
	defer state.mu.Unlock()
	if state.closed && !state.expired && state.closedControl != nil {
		// prior conserva o único controle terminal, com parâmetros já vinculados pelo digest.
		prior := state.closedControl
		if prior.requestID == requestID && prior.generation == generation {
			wanted := "close"
			if prior.ticket.Discarded {
				wanted = "closeDiscard"
			}
			if prior.digest != fingerprint || operation != wanted {
				return AuthoringControlTicket{}, &AuthoringError{Code: "invalidArgument"}
			}
			return prior.ticket, nil
		}
	}
	return AuthoringControlTicket{}, &AuthoringError{Code: "sessionExpired"}
}

// controlLookupLocked consulta replay/CAS sem alterar geração ou aposentar a revisão existente.
func (a *Authoring) controlLookupLocked(state *authoringState, requestID, fingerprint string, expected, generation uint64) (AuthoringControlTicket, bool, error) {
	if len(state.recent) > 0 {
		// last é conferido antes de recusar a geração já avançada pelo primeiro sucesso.
		last := state.recent[len(state.recent)-1]
		if last.generation == generation && last.requestID == requestID {
			if last.digest != fingerprint {
				return AuthoringControlTicket{}, false, &AuthoringError{Code: "invalidArgument"}
			}
			return last.ticket, true, nil
		}
	}
	if generation < state.generation {
		return AuthoringControlTicket{}, false, &AuthoringError{Code: "controlRetired"}
	}
	if generation != state.generation || expected != state.model.Document.Project.Revision {
		return AuthoringControlTicket{}, false, &AuthoringError{Code: "staleReview"}
	}
	return AuthoringControlTicket{}, false, nil
}

// ReserveTarget admite um alvo persistível sem consumir outra posição nas edições seguintes.
func (a *Authoring) ReserveTarget(id, target string) error {
	// state/release mantêm a mesma ordem de coordenação das mutações locais.
	state, release, err := a.acquire(id)
	if err != nil {
		return err
	}
	defer release()
	if target == "" || len(target) > 256 {
		return &AuthoringError{Code: "invalidArgument"}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.reserveLocked(state, target)
}

// reserveLocked pressupõe locks do registro e da sessão, sem IO ou espera humana.
func (a *Authoring) reserveLocked(state *authoringState, target string) error {
	if state.reserved[target] {
		return nil
	}
	for _, intent := range state.intents {
		if !intent.resolved && intent.target == target {
			return &AuthoringError{Code: "recoveryRequired"}
		}
	}
	if state.positions() >= a.limits.Commits || a.staging+a.limits.TargetStagingBytes > a.limits.StagingBytes {
		return &AuthoringError{Code: "resourceLimit"}
	}
	state.reserved[target] = true
	state.staging += a.limits.TargetStagingBytes
	a.staging += a.limits.TargetStagingBytes
	return nil
}

// BeginIntent converte reserva em posição incerta; não grava nem declara sucesso.
func (a *Authoring) BeginIntent(id, target, commitID string) error {
	// state/release serializam consumo da posição sem depender de vaga de edição.
	state, release, err := a.acquire(id)
	if err != nil {
		return err
	}
	defer release()
	if !project.ValidID(commitID) {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if prior, ok := state.intents[commitID]; ok {
		if prior.target != target {
			return &AuthoringError{Code: "invalidArgument"}
		}
		if prior.resolved {
			return &AuthoringError{Code: "authorizationConsumed"}
		}
		return nil
	}
	if !state.reserved[target] {
		return &AuthoringError{Code: "resourceLimit"}
	}
	if state.pendingIntent() {
		return &AuthoringError{Code: "busy"}
	}
	delete(state.reserved, target)
	state.intents[commitID] = authoringIntent{target: target}
	return nil
}

// ResolveIntent reconcilia a mesma posição; sucesso só pode vir do adaptador que provou o efeito.
// Falha comprovadamente sem efeito devolve a reserva ao alvo, sem gastar novo commit.
func (a *Authoring) ResolveIntent(id, commitID string, committed bool) error {
	// state/release protegem resultado sem criar 65ª posição em recuperação.
	state, release, err := a.acquire(id)
	if err != nil {
		return err
	}
	defer release()
	// intent deve existir no registro; entrada desconhecida não autoriza publicação.
	intent, ok := state.intents[commitID]
	if !ok {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if intent.resolved {
		return nil
	}
	if !committed {
		delete(state.intents, commitID)
		state.reserved[intent.target] = true
		return nil
	}
	intent.resolved = true
	state.intents[commitID] = intent
	state.consumed++
	a.mu.Lock()
	state.staging -= a.limits.TargetStagingBytes
	a.staging -= a.limits.TargetStagingBytes
	a.mu.Unlock()
	if intent.target == "document" {
		state.model.SavedRevision = state.model.Document.Project.Revision
		state.unsaved = false
		state.base = state.model.Document.Copy()
	}
	// key/draft marcam somente o alvo persistível; análise sem target continua pendente.
	for key, draft := range state.drafts {
		if draft.Target == intent.target {
			draft.Dirty = false
			state.drafts[key] = draft
		}
	}
	state.proposalDirty = false
	for _, draft := range state.drafts {
		state.proposalDirty = state.proposalDirty || draft.Dirty
	}
	return nil
}

// positions conta posições uma única vez, incluindo intents ainda sem resultado.
func (state *authoringState) positions() int {
	// total começa pelos commits concluídos e alvos ainda reservados.
	total := state.consumed + len(state.reserved)
	// intent incerta já ocupava sua posição antes da gravação.
	for _, intent := range state.intents {
		if !intent.resolved {
			total++
		}
	}
	return total
}

// pendingIntent impede que descarte ou edição encubram uma publicação ainda incerta.
func (state *authoringState) pendingIntent() bool {
	// intent ocupa a mesma posição até o adaptador comprovar seu resultado.
	for _, intent := range state.intents {
		if !intent.resolved {
			return true
		}
	}
	return false
}

// retireLocked libera conteúdo transitório e conserva tombstone sem renovar sessão.
func (a *Authoring) retireLocked(state *authoringState, expired bool) {
	a.retained -= state.retained - 1024
	// held preserva staging físico durante limpeza fora da trava global.
	held := a.retireFrozenLocked(state)
	a.staging -= state.staging - held
	state.closed = !expired
	state.expired = expired
	state.retained = 1024
	state.staging = 0
	state.model = nil
	state.base = Document{}
	state.edits = nil
	state.recent = nil
	state.closedControl = nil
	state.reserved = nil
	state.intents = nil
	state.drafts = nil
	state.draftBytes = 0
}

// authoringSessionWeight mede dados retidos sem incluir cópia duplicada de recurso no JSON.
func authoringSessionWeight(model *Session, base Document) int64 {
	// total inclui documento/base e cada snapshot de undo/redo retido.
	// current/baseWeight reservam também a futura base salva antes de admitir edição.
	current := authoringDocumentWeight(model.Document)
	baseWeight := authoringDocumentWeight(base)
	if current > baseWeight {
		baseWeight = current
	}
	total := current + baseWeight
	// document é imutável durante esta estimativa sob o mutex da sessão.
	for _, document := range model.UndoStack {
		total += authoringDocumentWeight(document)
	}
	for _, document := range model.RedoStack {
		total += authoringDocumentWeight(document)
	}
	return total
}

// authoringDocumentWeight conta metadados/textos e bytes reais dos recursos.
func authoringDocumentWeight(document Document) int64 {
	// resources são lidos sem mutar o mapa; a cópia por valor remove somente o header local.
	resources := document.Resources
	document.Resources = nil
	// metadata/err tratam tamanho inválido como quota excedida, nunca como zero.
	metadata, err := json.Marshal(document)
	if err != nil {
		return math.MaxInt64 / 4
	}
	total := int64(len(metadata))
	// data inclui recursos incorporados também em cada snapshot histórico.
	for _, data := range resources {
		total += int64(len(data))
	}
	return total
}

// copyAuthoringResult impede que o chamador altere a resposta conservada para replay.
func copyAuthoringResult(result AuthoringEditResult) AuthoringEditResult {
	// copied possui seu próprio mapa de identidades retornadas.
	copied := AuthoringEditResult{Revision: result.Revision, Values: map[string]string{}}
	// key/value são dados pequenos já limitados na fronteira de aplicação.
	for key, value := range result.Values {
		copied.Values[key] = value
	}
	return copied
}

// String descreve apenas quantidade de sessões, sem conteúdo ou caminhos pessoais.
func (a *Authoring) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("autoria: %d identidades", len(a.sessions))
}

// CloseAll descarta somente conteúdo transitório ao terminar o transporte.
// Publicação real deve reconciliar seu journal antes de chamar este encerramento.
func (a *Authoring) CloseAll() {
	defer a.drainReleases()
	a.mu.Lock()
	// states permite liberar o registro antes de esperar operações locais terminarem.
	states := make([]*authoringState, 0, len(a.sessions))
	for _, state := range a.sessions {
		states = append(states, state)
	}
	a.mu.Unlock()
	// state não é aguardado enquanto a trava global estiver retida.
	for _, state := range states {
		state.mu.Lock()
		if !state.closed && !state.expired && !state.pendingIntent() {
			a.mu.Lock()
			a.retireLocked(state, false)
			a.mu.Unlock()
		}
		state.mu.Unlock()
	}
}
