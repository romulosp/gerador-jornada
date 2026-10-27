package mcp

import (
	"context"
	"strconv"

	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// setValue aplica somente presença explícita; zero/vazio são valores deliberados, não omissão.
func setValue[T any](target *T, value *T) {
	if value != nil {
		*target = *value
	}
}

// adaptComponentChanges traduz DTOs para tipos reais, sem map livre no caso de uso.
func adaptComponentChanges(before screen.Component, changes componentChanges) (application.ComponentChanges, error) {
	// result conserva omissões dos campos comuns.
	result := application.ComponentChanges{Name: changes.Name, FileID: changes.FileID, Enabled: changes.Enabled}
	if changes.Bounds != nil {
		result.Bounds = &screen.Bounds{X: changes.Bounds.X, Y: changes.Bounds.Y, Width: changes.Bounds.Width, Height: changes.Bounds.Height}
	}
	if changes.Style != nil {
		// style copia overrides atuais antes de aplicar campos explicitamente informados.
		style, patch := before.Style, changes.Style
		setValue(&style.FontFamily, patch.FontFamily)
		setValue(&style.FontSize, patch.FontSize)
		setValue(&style.TextColor, patch.TextColor)
		setValue(&style.FillColor, patch.FillColor)
		setValue(&style.BorderColor, patch.BorderColor)
		setValue(&style.BackgroundAssetID, patch.BackgroundAssetID)
		setValue(&style.TextAlign, patch.TextAlign)
		setValue(&style.Skin, patch.Skin)
		result.Style = &style
	}
	if changes.Properties != nil {
		if changes.Properties.value == nil || changes.Properties.value.family() != string(before.Kind) {
			return result, &application.AuthoringError{Code: "invalidArgument"}
		}
		// properties conserva a variante atual e converte cada patch para campos concretos.
		properties := before.Properties.Copy()
		switch patch := changes.Properties.value.(type) {
		case buttonPatch:
			setValue(&properties.Button.Text, patch.Text)
			setValue(&properties.Button.Variant, patch.Variant)
		case textPatch:
			if properties.Label != nil {
				setValue(&properties.Label.Text, patch.Text)
			} else if properties.Title != nil {
				setValue(&properties.Title.Text, patch.Text)
			} else {
				return result, &application.AuthoringError{Code: "invalidArgument"}
			}
		case fieldPatch:
			setValue(&properties.Field.InitialValue, patch.InitialValue)
			setValue(&properties.Field.Format, patch.Format)
		case radioPatch:
			if patch.Options != nil {
				properties.Radio.Options = adaptOptions(*patch.Options)
			}
			setValue(&properties.Radio.Orientation, patch.Orientation)
			setValue(&properties.Radio.SelectedOptionID, patch.SelectedOptionID)
		case reportPatch:
			if patch.Rows != nil {
				properties.Report.Rows = adaptOptions(*patch.Rows)
			}
			setValue(&properties.Report.Title, patch.Title)
		case keyboardPatch:
			setValue(&properties.Keyboard.TargetComponentID, patch.TargetComponentID)
		default:
			return result, &application.AuthoringError{Code: "capabilityUnavailable"}
		}
		result.Properties = &properties
	}
	return result, nil
}

// adaptOptions conserva texto/UUID sem gerar IDs silenciosamente durante uma atualização.
func adaptOptions(options []optionDTO) []catalog.Option {
	// result é coleção nova, sem referência mutável ao DTO da requisição.
	result := make([]catalog.Option, 0, len(options))
	for _, option := range options {
		result = append(result, catalog.Option{ID: option.ID, Text: option.Text})
	}
	return result
}

// registerComponents registra somente comandos I1 sobre Kinds já instalados.
func registerComponents(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "componente_inserir", Description: "Insere Kind instalado com defaults do catálogo e do legado."},
		func(ctx context.Context, request insertRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			if request.X < 0 || request.Y < 0 {
				return nil, &application.AuthoringError{Code: "invalidArgument"}
			}
			return edit(ctx, runtime, "componente_inserir", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				if session.Document.Screen(request.ScreenID) == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				// id/defaults são produzidos pelo mesmo comando utilizado pelo editor.
				id, err := session.AddComponent(request.ScreenID, catalog.Kind(request.Kind), request.X, request.Y)
				return map[string]string{"componentId": id}, err
			})
		})
	registerTool(server, runtime, &sdk.Tool{Name: "componente_atualizar", Description: "Atualiza campos comuns e patch fechado da família; não converte Kind."},
		func(ctx context.Context, request componentUpdateRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID, request.ComponentID); err != nil {
				return nil, err
			}
			return edit(ctx, runtime, "componente_atualizar", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				// page/component são consultados no candidato protegido da edição.
				page := session.Document.Screen(request.ScreenID)
				if page == nil || page.ComponentByID(request.ComponentID) == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				changes, err := adaptComponentChanges(*page.ComponentByID(request.ComponentID), request.Changes)
				if err != nil {
					return nil, err
				}
				return map[string]string{"componentId": request.ComponentID}, session.ConfigureComponent(request.ScreenID, request.ComponentID, changes, builtinassets.ValidateAppearanceEdit)
			})
		})
	registerTool(server, runtime, &sdk.Tool{Name: "componente_mover", Description: "Move componente com as mesmas restrições geométricas e de hierarquia do editor."},
		func(ctx context.Context, request componentMoveRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID, request.ComponentID); err != nil {
				return nil, err
			}
			if request.X < 0 || request.Y < 0 {
				return nil, &application.AuthoringError{Code: "invalidArgument"}
			}
			// result conserva também a geometria confirmada, sem reler estado posterior num replay.
			result, err := runEdit(ctx, runtime, "componente_mover", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				// component fornece a origem para o comando coletivo já usado pela GUI.
				page := session.Document.Screen(request.ScreenID)
				if page == nil || page.ComponentByID(request.ComponentID) == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				component := page.ComponentByID(request.ComponentID)
				if err := session.MoveComponents(request.ScreenID, []string{request.ComponentID}, request.X-component.Bounds.X, request.Y-component.Bounds.Y); err != nil {
					return nil, err
				}
				bounds := session.Document.Screen(request.ScreenID).ComponentByID(request.ComponentID).Bounds
				return map[string]string{"componentId": request.ComponentID, "x": strconv.Itoa(bounds.X), "y": strconv.Itoa(bounds.Y), "width": strconv.Itoa(bounds.Width), "height": strconv.Itoa(bounds.Height)}, nil
			})
			if err != nil {
				return nil, err
			}
			// x/y/width/height são escalares produzidos pelo servidor, armazenados junto da resposta original.
			x, _ := strconv.Atoi(result.Values["x"])
			y, _ := strconv.Atoi(result.Values["y"])
			width, _ := strconv.Atoi(result.Values["width"])
			height, _ := strconv.Atoi(result.Values["height"])
			return successResult(struct {
				ComponentID string    `json:"componentId"`
				Bounds      boundsDTO `json:"bounds"`
			}{result.Values["componentId"], boundsDTO{X: x, Y: y, Width: width, Height: height}}, strconv.FormatUint(result.Revision, 10)), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "componente_remover", Description: "Remove componente sob política explícita; referências externas precisam ser resolvidas antes."},
		func(ctx context.Context, request componentRemoveRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.ScreenID, request.ComponentID); err != nil {
				return nil, err
			}
			return edit(ctx, runtime, "componente_remover", request.mutationRequest, request, func(session *application.Session) (map[string]string, error) {
				return map[string]string{"componentId": request.ComponentID}, session.DeleteComponentWithPolicy(request.ScreenID, request.ComponentID, request.ReferencePolicy)
			})
		})
}
