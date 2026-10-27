package mcp

import (
	"context"
	"strconv"

	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// edit executa um comando compartilhado com CAS/replay/cancelamento e resposta pequena de IDs.
func edit(ctx context.Context, runtime *runtime, operation string, request mutationRequest, arguments any, apply func(*application.Session) (map[string]string, error)) (*sdk.CallToolResult, error) {
	// result é armazenado pelo coordenador antes de responder, inclusive nas repetições de requestId.
	result, err := runEdit(ctx, runtime, operation, request, arguments, apply)
	if err != nil {
		return nil, err
	}
	return successResult(result.Values, strconv.FormatUint(result.Revision, 10)), nil
}

// runEdit permite respostas tipadas específicas sem repetir a política de CAS/replay do coordenador.
func runEdit(ctx context.Context, runtime *runtime, operation string, request mutationRequest, arguments any, apply func(*application.Session) (map[string]string, error)) (application.AuthoringEditResult, error) {
	if err := validIDs(request.SessionID, request.RequestID); err != nil {
		return application.AuthoringEditResult{}, err
	}
	return runtime.authoring.EditContext(ctx, request.SessionID, request.RequestID, requestDigest(operation, arguments), request.ExpectedRevision.number(), apply)
}

// registerScreens adapta comandos sem transferir regras de tamanho/nome ao protocolo.
func registerScreens(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "tela_criar", Description: "Cria tela em memória com os mesmos defaults do editor."},
		func(ctx context.Context, request screenCreateRequest) (*sdk.CallToolResult, error) {
			return edit(ctx, runtime, "tela_criar", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				// id é gerado pelo caso de uso comum, com fundo da tela inicial ou recurso incorporado.
				id, err := session.CreateScreen(request.Name, session.Document.Journey.StartScreenID, application.ScreenChanges{Width: request.Width, Height: request.Height, BackgroundColor: request.BackgroundColor}, builtinassets.UseBackground)
				return map[string]string{"screenId": id}, err
			})
		})
	registerTool(server, runtime, &sdk.Tool{Name: "tela_configurar", Description: "Configura campos de tela sem reposicionar controles implicitamente."},
		func(ctx context.Context, request screenUpdateRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			return edit(ctx, runtime, "tela_configurar", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				// changes permanece tipado na aplicação, sem DTO de transporte no domínio.
				changes := application.ScreenChanges{Name: request.Changes.Name, FileID: request.Changes.FileID, Width: request.Changes.Width, Height: request.Changes.Height, BackgroundColor: request.Changes.BackgroundColor}
				err := session.Change(func(document *application.Document) error {
					return application.ConfigureScreen(document, request.ScreenID, changes)
				})
				return map[string]string{"screenId": request.ScreenID}, err
			})
		})
	registerTool(server, runtime, &sdk.Tool{Name: "tela_duplicar", Description: "Duplica tela preservando regras de remapeamento e ordem do editor."},
		func(ctx context.Context, request screenRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			// result conserva o mapa integral do primeiro sucesso, sem recriar UUIDs no replay.
			result, err := runEdit(ctx, runtime, "tela_duplicar", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				if session.Document.Screen(request.ScreenID) == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				// id é devolvido pelo mesmo comando que remapeia componentes/referências internas.
				id, mapping, err := session.DuplicateScreenWithMapping(request.ScreenID)
				mapping["screenId"] = id
				return mapping, err
			})
			if err != nil {
				return nil, err
			}
			// mapping é cópia do cache de resposta, separada da identidade principal.
			id := result.Values["screenId"]
			delete(result.Values, "screenId")
			return successResult(struct {
				ScreenID string            `json:"screenId"`
				IDMap    map[string]string `json:"idMap"`
			}{id, result.Values}, strconv.FormatUint(result.Revision, 10)), nil
		})
}
