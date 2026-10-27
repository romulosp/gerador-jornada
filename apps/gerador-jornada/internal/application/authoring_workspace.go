package application

import (
	"context"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// AuthoringJourney descreve uma entrada autorizada, inclusive diagnóstico de corrupção.
type AuthoringJourney struct {
	// Filename é nome regular na raiz autorizada, nunca caminho fornecido pelo agente.
	Filename string
	// Name é conteúdo de exibição e não deve ser reinterpretado como instrução.
	Name string
	// Diagnostic é código seguro por entrada, sem mensagens internas de IO.
	Diagnostic string
}

// AuthoringInventory vincula a paginação ao conteúdo de um inventário da raiz.
type AuthoringInventory struct {
	// Revision muda se identidade, nome, diagnóstico ou conteúdo das entradas mudar.
	Revision string
	// Entries é cópia independente, em ordem determinística.
	Entries []AuthoringJourney
}

// AuthoringDocumentSource fornece somente leitura e construção, sem autoridade para publicar.
// A implementação compartilha o codec e o construtor usados pelo editor manual.
type AuthoringDocumentSource interface {
	ListJourneys(context.Context) (AuthoringInventory, error)
	OpenDocument(context.Context, string) (Document, AuthoringBaseline, error)
	NewDocument(name string, example bool) Document
}

// AuthoringWorkspace coordena a fonte fixada na composição e o único registro de sessões.
type AuthoringWorkspace struct {
	// core é compartilhado com edição/revisão; source não é substituível por ferramenta.
	core   *Authoring
	source AuthoringDocumentSource
}

// NewAuthoringWorkspace exige dependências concretas na composição, antes de receber ferramentas.
func NewAuthoringWorkspace(core *Authoring, source AuthoringDocumentSource) *AuthoringWorkspace {
	if core == nil || source == nil {
		panic("fonte de autoria ausente")
	}
	return &AuthoringWorkspace{core: core, source: source}
}

// List retorna somente entradas oferecidas pela fonte autorizada, sem novos caminhos.
func (workspace *AuthoringWorkspace) List(ctx context.Context) (AuthoringInventory, error) {
	return workspace.source.ListJourneys(ctx)
}

// admissionReplay consulta a identidade antes de qualquer nova leitura ou construção.
func (a *Authoring) admissionReplay(requestID, fingerprint string) (AuthoringSessionInfo, bool, error) {
	if !project.ValidID(requestID) || fingerprint == "" || len(fingerprint) > 256 {
		return AuthoringSessionInfo{}, false, &AuthoringError{Code: "invalidArgument"}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	// prior conserva inclusive tombstones; repetir não ressuscita uma sessão encerrada.
	prior, found := a.admissions[requestID]
	if found && prior.digest != fingerprint {
		return AuthoringSessionInfo{}, false, &AuthoringError{Code: "invalidArgument"}
	}
	return prior.info, found, nil
}

// Open admite os bytes e a base da mesma leitura; replay não volta ao filesystem.
func (workspace *AuthoringWorkspace) Open(ctx context.Context, requestID, fingerprint, filename string) (AuthoringSessionInfo, error) {
	if err := ctx.Err(); err != nil {
		return AuthoringSessionInfo{}, err
	}
	// prior é verificado antes de abrir uma origem que pode ter mudado ou desaparecido.
	prior, found, err := workspace.core.admissionReplay(requestID, fingerprint)
	if err != nil || found {
		return prior, err
	}
	if !catalog.ValidJourneyFilename(filename) {
		return AuthoringSessionInfo{}, &AuthoringError{Code: "invalidArgument"}
	}
	// document/baseline são fornecidos juntos pelo adaptador de IO.
	document, baseline, err := workspace.source.OpenDocument(ctx, filename)
	if err != nil {
		return AuthoringSessionInfo{}, err
	}
	if err := ctx.Err(); err != nil {
		return AuthoringSessionInfo{}, err
	}
	return workspace.core.Admit(requestID, fingerprint, document, baseline)
}

// Create admite um rascunho com defaults comuns, sem gravar arquivo ou escolher destino.
func (workspace *AuthoringWorkspace) Create(ctx context.Context, requestID, fingerprint, name string, example bool) (AuthoringSessionInfo, error) {
	if err := ctx.Err(); err != nil {
		return AuthoringSessionInfo{}, err
	}
	// prior evita gerar novas identidades ao repetir a mesma solicitação.
	prior, found, err := workspace.core.admissionReplay(requestID, fingerprint)
	if err != nil || found {
		return prior, err
	}
	if err := ValidateAuthoringName(name); err != nil {
		return AuthoringSessionInfo{}, err
	}
	return workspace.core.Admit(requestID, fingerprint, workspace.source.NewDocument(name, example), AuthoringBaseline{})
}
