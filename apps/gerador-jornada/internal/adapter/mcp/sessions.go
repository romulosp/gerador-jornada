package mcp

import (
	"context"
	"strconv"
	"time"

	"gerador-jornada.local/editor/internal/application"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// sessionQuotas distingue conteúdo admitido, posições de commit e reservas ainda pendentes.
type sessionQuotas struct {
	// EditRequests/Proposals não são zerados para abrir espaço artificialmente no replay.
	EditRequests int `json:"editRequests"`
	Proposals    int `json:"proposals"`
	// CommitPositions soma commits consumidos, intenções incertas e alvos reservados.
	CommitPositions  int `json:"commitPositions"`
	ConsumedCommits  int `json:"consumedCommits"`
	UncertainIntents int `json:"uncertainIntents"`
	ReservedTargets  int `json:"reservedTargets"`
	// StagingReservedBytes representa a parcela da sessão, em bytes.
	StagingReservedBytes int64 `json:"stagingReservedBytes"`
}

// processUsage informa pressão operacional, inclusive limpezas físicas que continuam cobradas.
type processUsage struct {
	// RetainedBytes/StagingReservedBytes são orçamentos lógicos do processo, em bytes.
	RetainedBytes        int64 `json:"retainedBytes"`
	StagingReservedBytes int64 `json:"stagingReservedBytes"`
	// PendingCleanups conta artefatos já destacados e ainda não liberados fisicamente.
	PendingCleanups        int  `json:"pendingCleanups"`
	ActiveCalls            int  `json:"activeCalls"`
	FinalizationInProgress bool `json:"finalizationInProgress"`
}

// sessionStatus transmite ciclo de vida sem enviar o documento ou suas imagens.
type sessionStatus struct {
	// SessionID/Revision/ControlGeneration são identidade e contadores decimais canônicos.
	SessionID         string `json:"sessionId"`
	Revision          string `json:"revision"`
	ControlGeneration string `json:"controlGeneration"`
	// Dirty agrega trabalho de documento e anexos, protegendo ambos no fechamento.
	DocumentDirty bool `json:"documentDirty"`
	ProposalDirty bool `json:"proposalDirty"`
	Dirty         bool `json:"dirty"`
	// RecoveryRequired impede tratar uma publicação incerta como descarte comum.
	RecoveryRequired bool `json:"recoveryRequired"`
	// ExpiresAt usa UTC RFC3339; consultar sessão ativa renova somente inatividade, nunca revisão.
	ExpiresAt string `json:"expiresAt"`
	// Quotas/ProcessUsage separam reservas da sessão da pressão global.
	Quotas       sessionQuotas `json:"quotas"`
	ProcessUsage processUsage  `json:"processUsage"`
	// Capabilities corresponde às ferramentas instaladas, sem promessa de gravação autorizada.
	Capabilities installedCapabilities `json:"capabilities"`
}

// sessionControlResult conserva somente o efeito imutável do controle para replay exato.
type sessionControlResult struct {
	// SessionID/Revision/ControlGeneration identificam o resultado original, não estado consultado depois.
	SessionID         string `json:"sessionId"`
	Revision          string `json:"revision"`
	ControlGeneration string `json:"controlGeneration"`
	// Closed/Discarded descrevem o efeito transitório; não provam autorização persistente.
	Closed    bool `json:"closed"`
	Discarded bool `json:"discarded"`
}

// sessionControl aplica CAS e high-water no núcleo e nunca publica arquivos.
func sessionControl(runtime *runtime, request controlRequest, operation string, args any) (*sdk.CallToolResult, error) {
	if err := validIDs(request.SessionID, request.RequestID); err != nil {
		return nil, err
	}
	// ticket permite repetir a resposta original mesmo depois de o documento ter sido liberado.
	ticket, err := runtime.authoring.Control(request.SessionID, request.RequestID, requestDigest(operation, args), request.ExpectedRevision.number(), request.ControlGeneration.number(), operation)
	if err != nil {
		return nil, err
	}
	// revision permanece string, inclusive acima do limite inteiro seguro do JavaScript.
	revision := strconv.FormatUint(ticket.Revision, 10)
	return successResult(sessionControlResult{SessionID: request.SessionID, Revision: revision, ControlGeneration: strconv.FormatUint(ticket.Generation, 10), Closed: ticket.Closed, Discarded: ticket.Discarded}, revision), nil
}

// registerSessions oferece controles transitórios separados das quotas de edição de conteúdo.
func registerSessions(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "sessao_obter", Description: "Consulta dirty agregado, geração, prazo, cotas e capacidades sem transportar documento.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request sessionRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.SessionID); err != nil {
				return nil, err
			}
			// snapshot consulta somente ciclo de vida, sem clonar recursos ou metadados grandes.
			snapshot, err := runtime.authoring.SessionStatus(request.SessionID)
			if err != nil {
				return nil, err
			}
			// usage é uma observação global independente do CAS da sessão e não concede permissões.
			usage := runtime.authoring.Usage()
			status := sessionStatus{SessionID: snapshot.ID, Revision: strconv.FormatUint(snapshot.Revision, 10), ControlGeneration: strconv.FormatUint(snapshot.ControlGeneration, 10), DocumentDirty: snapshot.DocumentDirty, ProposalDirty: snapshot.ProposalDirty, Dirty: snapshot.Dirty, RecoveryRequired: snapshot.RecoveryRequired, ExpiresAt: snapshot.ExpiresAt.UTC().Format(time.RFC3339),
				Quotas:       sessionQuotas{EditRequests: snapshot.EditRequests, Proposals: snapshot.Proposals, CommitPositions: snapshot.CommitPositions, ConsumedCommits: snapshot.ConsumedCommits, UncertainIntents: snapshot.UncertainIntents, ReservedTargets: snapshot.ReservedTargets, StagingReservedBytes: snapshot.StagingReservedBytes},
				ProcessUsage: processUsage{RetainedBytes: usage.RetainedBytes, StagingReservedBytes: usage.StagingReservedBytes, PendingCleanups: usage.PendingCleanups, ActiveCalls: usage.ActiveCalls, FinalizationInProgress: usage.FinalizationInProgress}, Capabilities: runtime.capabilities()}
			return successResult(status, status.Revision), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "sessao_descartar", Description: "Descarta alterações/anexos e revisão, restaura base conhecida e libera reservas; não salva."},
		func(_ context.Context, request controlRequest) (*sdk.CallToolResult, error) {
			return sessionControl(runtime, request, "discard", request)
		})
	registerTool(server, runtime, &sdk.Tool{Name: "sessao_encerrar", Description: "Encerra sessão; trabalho pendente exige discardUnsaved=true. Não salva nem autoriza gravação."},
		func(_ context.Context, request closeRequest) (*sdk.CallToolResult, error) {
			// operation distingue fechamento limpo de descarte explicitamente informado.
			operation := "close"
			if request.DiscardUnsaved {
				operation = "closeDiscard"
			}
			return sessionControl(runtime, request.controlRequest, operation, request)
		})
}

// sweepIdle observa inatividade até cancelamento/EOF e nunca espera uma sessão ocupada.
func sweepIdle(ctx context.Context, authoring *application.Authoring, interval time.Duration, done chan<- struct{}) {
	defer close(done)
	// ticker limita a varredura sem criar uma goroutine por sessão ou por chamada.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			authoring.ExpireIdle()
		}
	}
}
