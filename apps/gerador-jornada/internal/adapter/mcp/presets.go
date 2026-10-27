package mcp

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPresets usa o mesmo caso de uso integral e os mesmos recursos da paleta manual.
func registerPresets(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "preset_inserir", Description: "Insere composição instalada com recursos/decoração legados numa única edição."},
		func(ctx context.Context, request presetRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			if request.X < 0 || request.Y < 0 {
				return nil, &application.AuthoringError{Code: "invalidArgument"}
			}
			// result mantém IDs ordenados no cache pequeno de edição, incluindo replay.
			result, err := runEdit(ctx, runtime, "preset_inserir", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				if session.Document.Screen(request.ScreenID) == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				values := map[string]string{}
				err := session.Change(func(document *application.Document) error {
					// ids e decoração são produzidos integralmente pela aplicação, sem pós-processamento exclusivo da GUI.
					ids, err := application.InsertPreset(document, request.ScreenID, request.PresetID, request.X, request.Y, builtinassets.TitleStrip)
					if err != nil {
						return err
					}
					for index, id := range ids {
						values[fmt.Sprintf("%04d", index)] = id
					}
					return nil
				})
				return values, err
			})
			if err != nil {
				return nil, err
			}
			// keys restaura a ordem original, sem depender de iteração de map de Go.
			keys, ids := []string{}, []string{}
			for key := range result.Values {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				ids = append(ids, result.Values[key])
			}
			return successResult(struct {
				ComponentIDs []string `json:"componentIds"`
			}{ids}, strconv.FormatUint(result.Revision, 10)), nil
		})
}
