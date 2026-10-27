package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// admissionResult mantém a resposta original pequena, mesmo após editar ou remover a origem.
type admissionResult struct {
	// SessionID/Revision identificam admissão e base inicial; consulta posterior usa jornada_obter.
	SessionID string `json:"sessionId"`
	Revision  string `json:"revision"`
	// Name/Draft são resumo original, não estado posterior obtido de uma sessão alterada.
	Name  string `json:"name"`
	Draft bool   `json:"draft"`
	// ExpiresAt/IdleSeconds/UnsavedPolicy informam prazo inicial e descarte obrigatório antes de editar.
	ExpiresAt     string `json:"expiresAt"`
	IdleSeconds   int64  `json:"idleSeconds"`
	UnsavedPolicy string `json:"unsavedPolicy"`
}

// admitted conserva os metadados originais da admissão, inclusive em replay depois do encerramento.
func admitted(info application.AuthoringSessionInfo, draft bool) admissionResult {
	return admissionResult{SessionID: info.ID, Revision: strconv.FormatUint(info.Revision, 10), Name: info.Name, Draft: draft, ExpiresAt: info.ExpiresAt.UTC().Format(time.RFC3339), IdleSeconds: int64(application.DefaultAuthoringLimits().Idle.Seconds()), UnsavedPolicy: "explicitDiscardRequired;idleOrEOFDiscardsWithoutSaving"}
}

// journeyEntry conserva nome e diagnóstico individual sem expor erro interno de IO.
type journeyEntry struct {
	// Filename é regular na raiz; Name permanece conteúdo literal de exibição.
	Filename   string `json:"filename"`
	Name       string `json:"name"`
	Diagnostic string `json:"diagnostic,omitempty"`
}

// screenSummary retorna geometria/defaults sem widgets ou recursos binários.
type screenSummary struct {
	// ID/Name/FileID separam identidade, nome visual e nome-base do projeto.
	ID     string `json:"id"`
	Name   string `json:"name"`
	FileID string `json:"fileId"`
	// Width/Height/BackgroundColor/BackgroundAssetID são dados reais da tela persistível.
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	BackgroundColor   string `json:"backgroundColor"`
	BackgroundAssetID string `json:"backgroundAssetId"`
	ComponentCount    int    `json:"componentCount"`
}

// componentSummary conserva o modelo tipado sem emitir imagens ou pacote em base64.
type componentSummary struct {
	// ID/Name/FileID/ParentID preservam identidades e hierarquia observável.
	ID       string `json:"id"`
	Name     string `json:"name"`
	FileID   string `json:"fileId"`
	ParentID string `json:"parentId"`
	// Kind/Enabled/Bounds/Style são dados do documento e não estado de widget.
	Kind    string     `json:"kind"`
	Enabled bool       `json:"enabled"`
	Bounds  boundsDTO  `json:"bounds"`
	Style   stylePatch `json:"style"`
	// Properties recebe um DTO concreto de família; compostos avançados usam resumo explicitamente limitado.
	Properties         any  `json:"properties"`
	AdvancedProperties bool `json:"advancedProperties"`
	// NavigationAction/TargetScreenID mostram vínculo explícito do botão, sem inferir ação por Kind.
	NavigationAction string `json:"navigationAction,omitempty"`
	TargetScreenID   string `json:"targetScreenId,omitempty"`
}

// summarizeDocumentComponent acrescenta ao resumo somente a navegação explicitamente registrada.
func summarizeDocumentComponent(document application.Document, component screen.Component) componentSummary {
	// summary preserva dados reais; consulta não cria binding ausente.
	summary := summarizeComponent(component)
	if component.Kind == catalog.Button {
		summary.NavigationAction = "none"
		summary.TargetScreenID = document.Journey.Destination(component.ID)
		if summary.TargetScreenID != "" {
			summary.NavigationAction = "navigate"
		}
	}
	return summary
}

// projectSummary informa contexto sem duplicar listas grandes ou bytes de recursos.
type projectSummary struct {
	// ID/Name/StartScreenID conservam o contexto da página devolvida.
	ID            string `json:"id"`
	Name          string `json:"name"`
	StartScreenID string `json:"startScreenId"`
	// ScreenCount/AssetCount permitem consultas específicas sem transporte integral do documento.
	ScreenCount int `json:"screenCount"`
	AssetCount  int `json:"assetCount"`
}

// documentPage identifica o recorte retornado por jornada_obter.
type documentPage[T any] struct {
	// pageResult conserva cursor/revisão e valores integrais da lista filtrada.
	pageResult[T]
	// View distingue screens/components/component, sem inferir o filtro no cliente.
	View    string         `json:"view"`
	Project projectSummary `json:"project"`
}

// requestDigest vincula replay à operação e aos argumentos tipados, inclusive CAS inicial.
func requestDigest(operation string, args any) string {
	// encoded tem campos fixos e não contém maps de negócio nem espaços semânticos ambíguos.
	encoded, err := json.Marshal(struct {
		Operation string `json:"operation"`
		Arguments any    `json:"arguments"`
	}{operation, args})
	if err != nil {
		panic("DTO de autoria não serializável")
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

// summarizeComponent adapta todas as propriedades I1, sem misturar famílias ou metadados da UI.
func summarizeComponent(component screen.Component) componentSummary {
	// result fornece overrides explícitos; vazio continua default do legado, não cor genérica.
	result := componentSummary{ID: component.ID, Name: component.Name, FileID: component.FileID, ParentID: component.ParentID, Kind: string(component.Kind), Enabled: component.Enabled,
		Bounds: boundsDTO{X: component.Bounds.X, Y: component.Bounds.Y, Width: component.Bounds.Width, Height: component.Bounds.Height},
		Style:  stylePatch{FontFamily: &component.Style.FontFamily, FontSize: &component.Style.FontSize, TextColor: &component.Style.TextColor, FillColor: &component.Style.FillColor, BorderColor: &component.Style.BorderColor, BackgroundAssetID: &component.Style.BackgroundAssetID, TextAlign: &component.Style.TextAlign, Skin: &component.Style.Skin}, AdvancedProperties: true}
	// kind/properties pertencem à cópia do snapshot, sem referência mutável à sessão.
	kind, properties := propertyKind{Kind: result.Kind}, component.Properties
	switch component.Kind {
	case catalog.Button:
		result.Properties = buttonPatch{propertyKind: kind, Text: &properties.Button.Text, Variant: &properties.Button.Variant}
	case catalog.Label:
		result.Properties = textPatch{propertyKind: kind, Text: &properties.Label.Text}
	case catalog.Title:
		result.Properties = textPatch{propertyKind: kind, Text: &properties.Title.Text}
	case catalog.TextField:
		result.Properties = fieldPatch{propertyKind: kind, InitialValue: &properties.Field.InitialValue, Format: &properties.Field.Format}
	case catalog.Radio:
		// options mantém IDs existentes; a consulta não gera identidade nem altera a seleção.
		options := []optionDTO{}
		for _, option := range properties.Radio.Options {
			options = append(options, optionDTO{ID: option.ID, Text: option.Text})
		}
		result.Properties = radioPatch{propertyKind: kind, Options: &options, Orientation: &properties.Radio.Orientation, SelectedOptionID: &properties.Radio.SelectedOptionID}
	case catalog.Report:
		// rows é lista independente para não expor coleção mutável do domínio.
		rows := []optionDTO{}
		for _, row := range properties.Report.Rows {
			rows = append(rows, optionDTO{ID: row.ID, Text: row.Text})
		}
		result.Properties = reportPatch{propertyKind: kind, Title: &properties.Report.Title, Rows: &rows}
	case catalog.NumericKeyboard:
		result.Properties = keyboardPatch{propertyKind: kind, TargetComponentID: &properties.Keyboard.TargetComponentID}
	default:
		result.AdvancedProperties = false
		result.Properties = struct {
			Kind        string `json:"kind"`
			SummaryOnly bool   `json:"summaryOnly"`
			ReasonCode  string `json:"reasonCode"`
		}{result.Kind, true, "capabilityUnavailable"}
	}
	return result
}

// registerProjects oferece leitura/admissão em memória, sem ferramenta implícita de gravação.
func registerProjects(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "workspace_listar_jornadas", Description: "Lista jornadas da raiz autorizada sem percorrer subpastas.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(ctx context.Context, request pageRequest) (*sdk.CallToolResult, error) {
			// inventory vincula paginação ao conteúdo real das entradas, inclusive corrupção.
			inventory, err := runtime.source.List(ctx)
			if err != nil {
				return nil, err
			}
			entries := []journeyEntry{}
			for _, entry := range inventory.Entries {
				entries = append(entries, journeyEntry{Filename: entry.Filename, Name: entry.Name, Diagnostic: entry.Diagnostic})
			}
			page, err := paginate(runtime.pages, "workspace", inventory.Revision, entries, request)
			if err != nil {
				return nil, err
			}
			return successResult(page, ""), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "jornada_criar", Description: "Cria rascunho em memória usando os mesmos defaults do launcher; não salva."},
		func(ctx context.Context, request createRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.RequestID); err != nil {
				return nil, err
			}
			if request.Mode != "blank" && request.Mode != "example" {
				return nil, &application.AuthoringError{Code: "invalidArgument"}
			}
			// info conserva a resposta inicial para replay sem gerar outro documento.
			info, err := runtime.source.Create(ctx, request.RequestID, requestDigest("jornada_criar", request), request.Name, request.Mode == "example")
			if err != nil {
				return nil, err
			}
			revision := strconv.FormatUint(info.Revision, 10)
			return successResult(admitted(info, true), revision), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "jornada_abrir", Description: "Abre snapshot de uma jornada regular da raiz autorizada; não altera a GUI nem a origem."},
		func(ctx context.Context, request openRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.RequestID); err != nil {
				return nil, err
			}
			// info/base derivam da mesma leitura, hash e identidade física no adaptador de workspace.
			info, err := runtime.source.Open(ctx, request.RequestID, requestDigest("jornada_abrir", request), request.Filename)
			if err != nil {
				return nil, err
			}
			revision := strconv.FormatUint(info.Revision, 10)
			return successResult(admitted(info, false), revision), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "jornada_obter", Description: "Consulta telas/componentes paginados sem transportar bytes de imagem ou pacote.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request documentRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.SessionID); err != nil {
				return nil, err
			}
			// snapshot não duplica recursos binários para responder apenas metadados e conteúdo tipado.
			snapshot, err := runtime.authoring.InspectMetadata(request.SessionID)
			if err != nil {
				return nil, err
			}
			revision := strconv.FormatUint(snapshot.Revision, 10)
			summary := projectSummary{ID: snapshot.Document.Project.ID, Name: snapshot.Document.Project.Name, StartScreenID: snapshot.Document.Journey.StartScreenID, ScreenCount: len(snapshot.Document.Screens), AssetCount: len(snapshot.Document.Project.Assets)}
			scope := "document/" + request.SessionID + "/" + request.ScreenID + "/" + request.ComponentID
			if request.ScreenID == "" {
				if request.ComponentID != "" {
					return nil, &application.AuthoringError{Code: "invalidArgument"}
				}
				// screens preserva a ordem original da jornada, sem ordenar por nome.
				screens := []screenSummary{}
				for _, page := range snapshot.Document.Screens {
					screens = append(screens, screenSummary{ID: page.ID, Name: page.Name, FileID: page.FileID, Width: page.Size.Width, Height: page.Size.Height, BackgroundColor: page.Background.Color, BackgroundAssetID: page.Background.AssetID, ComponentCount: len(page.Components)})
				}
				page, err := paginate(runtime.pages, scope, revision, screens, request.pageRequest)
				if err != nil {
					return nil, err
				}
				return successResult(documentPage[screenSummary]{pageResult: page, View: "screens", Project: summary}, revision), nil
			}
			if err := validIDs(request.ScreenID); err != nil {
				return nil, err
			}
			selected := snapshot.Document.Screen(request.ScreenID)
			if selected == nil {
				return nil, &application.AuthoringError{Code: "invalidReference"}
			}
			// components é filtrado antes de paginar; ID desconhecido não vira procura externa.
			components := []componentSummary{}
			view := "components"
			if request.ComponentID != "" {
				if err := validIDs(request.ComponentID); err != nil {
					return nil, err
				}
				component := selected.ComponentByID(request.ComponentID)
				if component == nil {
					return nil, &application.AuthoringError{Code: "invalidReference"}
				}
				components = append(components, summarizeDocumentComponent(snapshot.Document, *component))
				view = "component"
			} else {
				for _, component := range selected.Components {
					components = append(components, summarizeDocumentComponent(snapshot.Document, component))
				}
			}
			page, err := paginate(runtime.pages, scope, revision, components, request.pageRequest)
			if err != nil {
				return nil, err
			}
			return successResult(documentPage[componentSummary]{pageResult: page, View: view, Project: summary}, revision), nil
		})
}
