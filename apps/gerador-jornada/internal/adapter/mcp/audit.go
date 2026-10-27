package mcp

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/application"
)

// audit emite apenas metadados permitidos; não recebe argumentos, documentos ou assinaturas.
type audit struct {
	// mu impede intercalar linhas de chamadas concorrentes no stderr escolhido na composição.
	mu sync.Mutex
	// destination não é stdout do protocolo; nil desabilita registro sem criar arquivo.
	destination io.Writer
}

// auditEvent é um formato fechado sem campo de texto livre, caminho ou identidade pessoal.
type auditEvent struct {
	// Operation/Code pertencem à allowlist estática de operações e resultados.
	Operation string `json:"operation"`
	Code      string `json:"code"`
	// Revision é decimal uint64 validado; DurationMS é duração operacional em milissegundos.
	Revision   string `json:"revision,omitempty"`
	DurationMS int64  `json:"durationMs"`
}

// newAudit recebe stderr explicitamente, sem logger desktop ou estado persistente.
func newAudit(destination io.Writer) *audit { return &audit{destination: destination} }

// record reduz qualquer campo não permitido antes de serializar uma linha independente.
func (audit *audit) record(operation, code, revision string, duration time.Duration) {
	if audit == nil || audit.destination == nil {
		return
	}
	switch operation {
	case "aplicacao_capacidades", "catalogo_listar", "catalogo_descrever", "catalogo_listar_presets", "catalogo_listar_skins", "workspace_listar_jornadas", "jornada_criar", "jornada_abrir", "jornada_obter", "tela_criar", "tela_configurar", "tela_duplicar", "componente_inserir", "componente_atualizar", "componente_mover", "componente_remover", "preset_inserir", "jornada_definir_inicio", "navegacao_vincular", "jornada_validar", "jornada_revisar_alteracoes", "jornada_salvar", "sessao_obter", "sessao_descartar", "sessao_encerrar":
	default:
		operation = "unknown"
	}
	if code != "ok" {
		code = errorCode(&application.AuthoringError{Code: code})
	}
	if parsed, err := strconv.ParseUint(revision, 10, 64); err != nil || strconv.FormatUint(parsed, 10) != revision {
		revision = ""
	}
	if duration < 0 {
		duration = 0
	}
	audit.mu.Lock()
	defer audit.mu.Unlock()
	// Encode não pode publicar conteúdo pessoal mesmo quando o consumidor de stderr falha.
	_ = json.NewEncoder(audit.destination).Encode(auditEvent{Operation: operation, Code: code, Revision: revision, DurationMS: duration.Milliseconds()})
}
