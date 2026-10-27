package mcp

import (
	"encoding/json"
	"strconv"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// revisionValue mantém uint64 como string decimal no protocolo, sem coerção numérica.
type revisionValue string

// UnmarshalJSON recusa sinal, expoente, zeros iniciais, overflow e tipo numérico JSON.
func (revision *revisionValue) UnmarshalJSON(raw []byte) error {
	// text é conferido antes da conversão para a revisão de domínio.
	var text string
	if json.Unmarshal(raw, &text) != nil {
		return &application.AuthoringError{Code: "invalidArgument"}
	}
	value, err := strconv.ParseUint(text, 10, 64)
	if err != nil || strconv.FormatUint(value, 10) != text {
		return &application.AuthoringError{Code: "invalidArgument"}
	}
	*revision = revisionValue(text)
	return nil
}

// number converte apenas valores já validados pelo decoder de protocolo.
func (revision revisionValue) number() uint64 {
	// value tem erro impossível após UnmarshalJSON, mas zero não concede mutação sem CAS.
	value, _ := strconv.ParseUint(string(revision), 10, 64)
	return value
}

// validIDs confere UUIDs v4 canônicos antes de consultar entidades de domínio.
func validIDs(ids ...string) error {
	for _, id := range ids {
		if !project.ValidID(id) {
			return &application.AuthoringError{Code: "invalidArgument"}
		}
	}
	return nil
}

// mutationRequest reúne identidade de edição e precondição sem conter conteúdo de domínio.
type mutationRequest struct {
	// SessionID/RequestID pertencem a namespaces de sessão/replay, nunca a caminhos.
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
	// ExpectedRevision permanece string no JSON, incluindo o máximo uint64.
	ExpectedRevision revisionValue `json:"expectedRevision"`
}

// noNullInput impede usar null como apagamento genérico nos comandos I1 de mutação.
// Contratos com campos explicitamente anuláveis usam DTOs próprios, sem este marcador.
func (mutationRequest) noNullInput() {}

// sessionRequest identifica consulta de uma única sessão autorizada.
type sessionRequest struct {
	// SessionID não cria sessão caso esteja encerrada ou expirada.
	SessionID string `json:"sessionId"`
}

// createRequest não recebe destino ou aprovação; a criação é transitória.
type createRequest struct {
	// RequestID permite repetir a admissão sem gerar novas identidades.
	RequestID string `json:"requestId"`
	// Name/Mode selecionam o mesmo construtor usado pelo launcher.
	Name string `json:"name"`
	Mode string `json:"mode"`
}

// openRequest aceita somente o nome regular na raiz previamente fixada.
type openRequest struct {
	// RequestID/Filename vinculam replay e origem da primeira leitura.
	RequestID string `json:"requestId"`
	Filename  string `json:"filename"`
}

// documentRequest recorta o snapshot sem transportar pacote, imagem ou caminho externo.
type documentRequest struct {
	// sessionRequest/pageRequest compartilham apenas metadados de consulta.
	sessionRequest
	pageRequest
	// ScreenID/ComponentID são filtros opcionais pertencentes ao mesmo documento.
	ScreenID    string `json:"screenId,omitempty"`
	ComponentID string `json:"componentId,omitempty"`
}

// screenCreateRequest oferece somente propriedades instaladas de criação de tela.
type screenCreateRequest struct {
	// mutationRequest mantém o CAS e a identidade de edição.
	mutationRequest
	// Name não substitui o FileID gerado pelas regras comuns.
	Name string `json:"name"`
	// Width/Height/BackgroundColor omitidos conservam o default compartilhado.
	Width           *int    `json:"width,omitempty"`
	Height          *int    `json:"height,omitempty"`
	BackgroundColor *string `json:"backgroundColor,omitempty"`
}

// screenChanges fecha os campos reais de configuração de tela.
type screenChanges struct {
	// Name/FileID têm regras independentes de conteúdo e nome-base gerado.
	Name   *string `json:"name,omitempty"`
	FileID *string `json:"fileId,omitempty"`
	// Width/Height são dimensões lógicas; BackgroundColor não recebe caminho de recurso.
	Width           *int    `json:"width,omitempty"`
	Height          *int    `json:"height,omitempty"`
	BackgroundColor *string `json:"backgroundColor,omitempty"`
}

// screenUpdateRequest separa a entidade alvo do patch fechado.
type screenUpdateRequest struct {
	// mutationRequest mantém uma edição atômica e idempotente.
	mutationRequest
	// ScreenID/Changes são resolvidos antes de publicar o candidato.
	ScreenID string        `json:"screenId"`
	Changes  screenChanges `json:"changes"`
}

// screenRequest serve comandos de entidade sem parâmetros de alteração livres.
type screenRequest struct {
	// mutationRequest conserva a mesma política de replay dos demais comandos.
	mutationRequest
	// ScreenID precisa pertencer ao documento da sessão.
	ScreenID string `json:"screenId"`
}

// insertRequest seleciona Kind instalado e coordenadas lógicas sem código executável.
type insertRequest struct {
	// mutationRequest identifica uma inserção independente.
	mutationRequest
	// ScreenID/Kind não são nome de skin, preset ou definição de biblioteca.
	ScreenID string `json:"screenId"`
	Kind     string `json:"kind"`
	// X/Y são coordenadas lógicas não negativas.
	X int `json:"x"`
	Y int `json:"y"`
}

// presetRequest seleciona composição instalada sem inventar Kind novo.
type presetRequest struct {
	// mutationRequest conserva uma entrada de histórico para a composição inteira.
	mutationRequest
	// ScreenID/PresetID identificam tela e preset ofertado.
	ScreenID string `json:"screenId"`
	PresetID string `json:"presetId"`
	// X/Y são a origem lógica da composição.
	X int `json:"x"`
	Y int `json:"y"`
}

// componentMoveRequest usa posição absoluta, convertida pelo caso de uso para o comando comum.
type componentMoveRequest struct {
	// mutationRequest conserva precondição e replay.
	mutationRequest
	// ScreenID/ComponentID devem resolver a mesma tela.
	ScreenID    string `json:"screenId"`
	ComponentID string `json:"componentId"`
	// X/Y não recebem unidades físicas ou coordenadas de janela.
	X int `json:"x"`
	Y int `json:"y"`
}

// componentRemoveRequest exige política explícita para referências afetadas.
type componentRemoveRequest struct {
	// mutationRequest aplica remoção como uma edição atômica.
	mutationRequest
	// ScreenID/ComponentID delimitam a entidade sem caminhos de arquivos.
	ScreenID    string `json:"screenId"`
	ComponentID string `json:"componentId"`
	// ReferencePolicy aceita somente reject/removeOwnedReferences no handler.
	ReferencePolicy string `json:"referencePolicy"`
}

// navigationRequest altera somente vínculos de navegação já suportados.
type navigationRequest struct {
	// mutationRequest conserva os mesmos metadados de autoria.
	mutationRequest
	// SourceScreenID/SourceComponentID identificam a origem; destino pode faltar ao remover.
	SourceScreenID    string `json:"sourceScreenId"`
	SourceComponentID string `json:"sourceComponentId"`
	TargetScreenID    string `json:"targetScreenId,omitempty"`
	// Action é navigate/none e não cria ação implícita por Kind.
	Action string `json:"action"`
}

// validationRequest consulta diagnóstico por alvo sem exportar nem gravar.
type validationRequest struct {
	// sessionRequest delimita o documento transitório.
	sessionRequest
	// pageRequest aplica o mesmo orçamento das demais consultas, inclusive a diagnósticos.
	pageRequest
	// Target aceita draft/preview/pptx conforme capacidade instalada.
	Target string `json:"target"`
}

// controlRequest separa controle da quota de edições de conteúdo.
type controlRequest struct {
	// mutationRequest mantém revisão de base e requestId própria da preparação.
	mutationRequest
	// ControlGeneration é high-water decimal e não aprovação humana.
	ControlGeneration revisionValue `json:"controlGeneration"`
}

// closeRequest exige descarte explícito quando existir trabalho pendente.
type closeRequest struct {
	// controlRequest identifica a tentativa sem renovar controles aposentados.
	controlRequest
	// DiscardUnsaved não é autorização de gravação persistente.
	DiscardUnsaved bool `json:"discardUnsaved"`
}

// reviewRequest recebe intenção/destino, mas não assinatura ou consentimento presumido.
type reviewRequest struct {
	// controlRequest usa geração de controle separada da revisão de documento.
	controlRequest
	// Filename/Intent/AllowDraft entram nos bytes congelados da revisão futura.
	Filename   string `json:"filename"`
	Intent     string `json:"intent"`
	AllowDraft bool   `json:"allowDraft"`
}
