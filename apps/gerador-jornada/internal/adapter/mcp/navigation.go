package mcp

import (
	"context"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/journey"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNavigation reutiliza regras de Jornada sem criar comandos de negócio por Kind.
func registerNavigation(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "jornada_definir_inicio", Description: "Define tela inicial existente no rascunho."},
		func(ctx context.Context, request screenRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			return edit(ctx, runtime, "jornada_definir_inicio", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				return map[string]string{"screenId": request.ScreenID}, session.Change(func(document *application.Document) error {
					return application.SetJourneyStart(document, request.ScreenID)
				})
			})
		})
	registerTool(server, runtime, &sdk.Tool{Name: "navegacao_vincular", Description: "Liga botão a tela explícita ou remove sua ação; não abre arquivo externo."},
		func(ctx context.Context, request navigationRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.SourceScreenID, request.SourceComponentID); err != nil {
				return nil, err
			}
			if request.TargetScreenID != "" {
				if err := validIDs(request.TargetScreenID); err != nil {
					return nil, err
				}
			}
			return edit(ctx, runtime, "navegacao_vincular", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				return map[string]string{"componentId": request.SourceComponentID}, session.Change(func(document *application.Document) error {
					return application.BindNavigation(document, journey.Action{SourceScreenID: request.SourceScreenID, SourceComponentID: request.SourceComponentID, TargetScreenID: request.TargetScreenID, Kind: request.Action})
				})
			})
		})
}
