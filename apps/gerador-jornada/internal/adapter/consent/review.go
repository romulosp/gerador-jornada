// Package consent adapta revisão/consentimento e staging privados, sem presumir autorização do agente.
package consent

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// hostIdentityPattern restringe IDs operacionais a uma representação segura e independente de nomes pessoais.
var hostIdentityPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// ReviewProfile identifica o host/chave configurados na composição; não prova qualificação nem autoriza R2.
type ReviewProfile struct {
	// HostID/KeyID referenciam configuração do operador, nunca declaração produzida pela LLM.
	HostID, KeyID string
	// ServerInstanceID muda a cada processo e invalida assinaturas antigas para novo efeito.
	ServerInstanceID string
}

// reviewDiagnostic é fechado conforme authorization.schema.json, sem texto livre de erro interno.
type reviewDiagnostic struct {
	// Code/Severity/Message são metadados de revisão apresentados como dados ao humano.
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	// Location identifica entidades, sem revelar caminho arbitrário do computador.
	Location string `json:"location,omitempty"`
}

// reviewDTO conserva exatamente a ordem required do schema HostApproval-v1.
// Campos anuláveis de proposta/materialização pertencem às extensões da biblioteca, não ao agente de jornada.
type reviewDTO struct {
	// ProfileVersion distingue este contrato das versões de software/manifesto.
	ProfileVersion int `json:"profileVersion"`
	// HostID/KeyID/ServerInstanceID são fixados na composição antes de preparar revisão.
	HostID           string `json:"hostId"`
	KeyID            string `json:"keyId"`
	ServerInstanceID string `json:"serverInstanceId"`
	// SessionID/RequestID/ReviewID vinculam sessão, commit separado e revisão corrente.
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
	ReviewID  string `json:"reviewId"`
	// Operation/TargetIdentity não podem ser substituídos no pedido de commit.
	Operation      string `json:"operation"`
	TargetIdentity string `json:"targetIdentity"`
	// BaseExists/BaseHash descrevem precondição; ausência não é permissão para sobrescrever homônimo.
	BaseExists bool    `json:"baseExists"`
	BaseHash   *string `json:"baseHash"`
	// ContentHash identifica bytes finais do artefato e ExpectedRevision mantém uint64 como decimal.
	ContentHash       string `json:"contentHash"`
	ExpectedRevision  string `json:"expectedRevision"`
	ControlGeneration string `json:"controlGeneration"`
	// Proposal/Materialization são null na preparação de jornada; a biblioteca usará variantes próprias.
	Proposal        *struct{} `json:"proposal"`
	Materialization *struct{} `json:"materialization"`
	// ExpiresAt/Challenge têm validade fixa e 32 bytes de desafio, sem renovação por replay.
	ExpiresAt int64  `json:"expiresAt"`
	Challenge string `json:"challenge"`
	// Summary/Diagnostics vêm do aplicativo e não de um resumo livre produzido pelo modelo.
	Summary     []string           `json:"summary"`
	Diagnostics []reviewDiagnostic `json:"diagnostics"`
}

// validProfile exige identidade técnica sem equiparar configuração à homologação do host.
func validProfile(profile ReviewProfile) bool {
	return hostIdentityPattern.MatchString(profile.HostID) && hostIdentityPattern.MatchString(profile.KeyID) && project.ValidID(profile.ServerInstanceID)
}

// encodeReview congela DTO UTF-8 de ordem fixa; o coordenador conserva hash/bytes originais.
func encodeReview(profile ReviewProfile, request application.AuthoringFreezeRequest, target targetSnapshot, contentHash string) ([]byte, error) {
	if !validProfile(profile) {
		return nil, &application.AuthoringError{Code: "authorizationUnavailable"}
	}
	if !project.ValidID(request.SessionID) || !project.ValidID(request.Ticket.ReviewID) || !project.ValidID(request.Ticket.CommitRequestID) || request.Ticket.ReviewID == request.Ticket.CommitRequestID || request.Ticket.Generation == 0 || request.Ticket.ExpiresAt < 1 || request.Ticket.ExpiresAt > 253402300799 || request.Revision != request.Document.Project.Revision {
		return nil, &application.AuthoringError{Code: "invalidArgument"}
	}
	// challenge usa entropia local; nenhuma chave privada ou decisão humana é criada pelo servidor.
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, err
	}
	operation := "createJourney"
	if request.Options.Intent == "replace" {
		operation = "replaceJourney"
	}
	var baseHash *string
	if target.baseline.Exists {
		value := target.baseline.SHA256
		baseHash = &value
	}
	// diagnostics usa códigos/localização estáveis; mensagens internas podem conter nomes/textos e não são transcritas.
	diagnostics := []reviewDiagnostic{{Code: "REVIEW_IS_NOT_AUTHORIZATION", Severity: "info", Message: "Esta revisão não autoriza gravar. Consentimento verificável separado é obrigatório."}}
	for index, problem := range request.Diagnostics {
		if index >= 998 {
			diagnostics = append(diagnostics, reviewDiagnostic{Code: "MORE_DIAGNOSTICS", Severity: "warning", Message: "Há diagnósticos adicionais; consultar a validação paginada da mesma revisão antes de decidir."})
			break
		}
		// location conserva somente UUIDs das entidades diagnosticadas.
		location := strings.Trim(strings.Join([]string{problem.ScreenID, problem.ComponentID, problem.ScenarioID}, "/"), "/")
		diagnostics = append(diagnostics, reviewDiagnostic{Code: problem.Code, Severity: problem.Severity, Message: "Há dados que exigem correção neste alvo.", Location: location})
	}
	for _, summary := range request.Summary {
		if !utf8.ValidString(summary) || utf8.RuneCountInString(summary) > 4096 {
			return nil, &application.AuthoringError{Code: "resourceLimit"}
		}
	}
	if len(request.Summary) > 1000 || utf8.RuneCountInString(target.identity) > 512 {
		return nil, &application.AuthoringError{Code: "resourceLimit"}
	}
	// encoded preserva a ordem de campos acima; não reconstrói um objeto enviado pelo cliente.
	encoded, err := json.Marshal(reviewDTO{ProfileVersion: 1, HostID: profile.HostID, KeyID: profile.KeyID, ServerInstanceID: profile.ServerInstanceID, SessionID: request.SessionID, RequestID: request.Ticket.CommitRequestID, ReviewID: request.Ticket.ReviewID, Operation: operation, TargetIdentity: target.identity, BaseExists: target.baseline.Exists, BaseHash: baseHash, ContentHash: contentHash, ExpectedRevision: strconv.FormatUint(request.Revision, 10), ControlGeneration: strconv.FormatUint(request.Ticket.Generation, 10), ExpiresAt: request.Ticket.ExpiresAt, Challenge: base64.RawURLEncoding.EncodeToString(challenge), Summary: request.Summary, Diagnostics: diagnostics})
	if err != nil {
		return nil, err
	}
	if len(encoded) > 512<<10 {
		return nil, &application.AuthoringError{Code: "resourceLimit"}
	}
	return encoded, nil
}
