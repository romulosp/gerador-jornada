package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"gerador-jornada.local/editor/internal/adapter/workspaceguard"
	"gerador-jornada.local/editor/internal/application"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxStructuredBytes reserva espaço para Content, ID escapada e envelope JSON-RPC.
const maxStructuredBytes = MaxFrameBytes - (8 << 10)

// successEnvelope mantém dados de negócio em um único objeto estruturado.
type successEnvelope struct {
	// OK/Revision indicam resultado e revisão decimal sem perda de precisão no cliente.
	OK       bool   `json:"ok"`
	Revision string `json:"revision,omitempty"`
	// Result é preenchido apenas com DTOs do adaptador, nunca widgets ou recursos binários.
	Result any `json:"result"`
}

// errorEnvelope separa falha segura de detalhes internos de filesystem/validação.
type errorEnvelope struct {
	// OK permanece falso em toda falha de ferramenta.
	OK bool `json:"ok"`
	// Error conserva somente metadados contratuais permitidos.
	Error errorDetail `json:"error"`
}

// errorDetail informa recuperação sem transcrever valores ou caminhos internos.
type errorDetail struct {
	// Code/Message são estáticos e independentes do conteúdo recebido.
	Code    string `json:"code"`
	Message string `json:"message"`
	// Retryable/StatePreserved orientam repetição sem conceder consentimento.
	Retryable      bool `json:"retryable"`
	StatePreserved bool `json:"statePreserved"`
}

// successResult limita a saída antes da escrita; nunca duplica StructuredContent em Content.
func successResult(result any, revision string) *sdk.CallToolResult {
	// encoded já mede a forma escapada que entrará no frame final.
	encoded, err := json.Marshal(successEnvelope{OK: true, Revision: revision, Result: result})
	if err != nil || len(encoded) > maxStructuredBytes {
		return failureCode("resourceLimit")
	}
	return &sdk.CallToolResult{Content: []sdk.Content{&sdk.TextContent{Text: "Operação concluída; resultado no conteúdo estruturado."}}, StructuredContent: json.RawMessage(encoded)}
}

// failureResult traduz somente classes conhecidas; o texto do erro original não sai do processo.
func failureResult(err error) *sdk.CallToolResult { return failureCode(errorCode(err)) }

// errorCode aplica allowlist para que uma mensagem arbitrária não vire código de protocolo.
func errorCode(err error) string {
	// typed identifica erros estáveis do coordenador sem importar DTOs no domínio.
	var typed *application.AuthoringError
	if errors.As(err, &typed) {
		switch typed.Code {
		case "invalidArgument", "unknownProperty", "invalidReference", "capabilityUnavailable", "permissionDenied", "authorizationUnavailable", "invalidAuthorization", "authorizationExpired", "authorizationConsumed", "staleReview", "controlRetired", "externalModification", "sessionExpired", "busy", "resourceLimit", "canceled", "persistenceFailed", "recoveryRequired":
			return typed.Code
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return "canceled"
	}
	if errors.Is(err, workspaceguard.ErrBusy) {
		return "busy"
	}
	if errors.Is(err, os.ErrPermission) {
		return "permissionDenied"
	}
	if errors.Is(err, application.ErrJourneyMissing) || errors.Is(err, os.ErrNotExist) {
		return "invalidReference"
	}
	return "invalidArgument"
}

// failureCode produz resposta curta; desconhecidos são reduzidos a invalidArgument.
func failureCode(code string) *sdk.CallToolResult {
	code = errorCode(&application.AuthoringError{Code: code})
	// message usa apenas a identidade segura conhecida, sem interpolar dados de origem.
	message := "Operação recusada: " + code + "."
	return &sdk.CallToolResult{IsError: true, Content: []sdk.Content{&sdk.TextContent{Text: message}}, StructuredContent: errorEnvelope{Error: errorDetail{Code: code, Message: message, Retryable: code == "busy", StatePreserved: code != "recoveryRequired" && code != "persistenceFailed"}}}
}
