package mcp

import (
	"context"
	"strconv"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// diagnosticDTO preserva código/localização sem ecoar mensagens internas com dados pessoais.
type diagnosticDTO struct {
	// Code/Severity/Message identificam a regra sem copiar textos do projeto.
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	// ScreenID/ComponentID/ScenarioID são localizadores autorizados do documento.
	ScreenID    string `json:"screenId,omitempty"`
	ComponentID string `json:"componentId,omitempty"`
	ScenarioID  string `json:"scenarioId,omitempty"`
}

// validationResult distingue rascunho legível de validade para execução/saída.
type validationResult struct {
	// Target/Readable/Valid não representam autorização ou homologação visual.
	Target      string                    `json:"target"`
	Readable    bool                      `json:"readable"`
	Valid       bool                      `json:"valid"`
	Diagnostics pageResult[diagnosticDTO] `json:"diagnostics"`
	// Limitations informa limites estáticos reais, sem prometer runtime Java/EFL ou Office homologado.
	Limitations []string `json:"limitations"`
}

// registerValidation consulta regras/estilos sem executar a prévia nem chamar exportador.
func registerValidation(server *sdk.Server, runtime *runtime) {
	registerTool(server, runtime, &sdk.Tool{Name: "jornada_validar", Description: "Valida rascunho/prévia/PPT sem salvar, exportar ou executar ações.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(ctx context.Context, request validationRequest) (*sdk.CallToolResult, error) {
			if err := validIDs(request.SessionID); err != nil {
				return nil, err
			}
			if request.Target != "draft" && request.Target != "preview" && request.Target != "pptx" {
				return nil, &application.AuthoringError{Code: "capabilityUnavailable"}
			}
			// snapshot dispensa bytes de imagens; Diagnostics usa o estado protegido para presença real de recursos.
			snapshot, err := runtime.authoring.InspectMetadata(request.SessionID)
			if err != nil {
				return nil, err
			}
			revision, problems, err := runtime.authoring.Diagnostics(ctx, request.SessionID)
			if err != nil {
				return nil, err
			}
			if revision != snapshot.Revision {
				return nil, &application.AuthoringError{Code: "staleReview"}
			}
			items := []diagnosticDTO{}
			for _, problem := range problems {
				items = append(items, diagnosticDTO{Code: problem.Code, Severity: problem.Severity, Message: "Há dados que exigem correção neste alvo.", ScreenID: problem.ScreenID, ComponentID: problem.ComponentID, ScenarioID: problem.ScenarioID})
			}
			if request.Target != "draft" {
				for _, screen := range snapshot.Document.Screens {
					for _, component := range screen.Components {
						if err := ctx.Err(); err != nil {
							return nil, err
						}
						if err := presentation.ValidateStyle(component); err != nil {
							items = append(items, diagnosticDTO{Code: "STYLE_UNAVAILABLE", Severity: "error", Message: "A aparência exige recurso instalado compatível.", ScreenID: screen.ID, ComponentID: component.ID})
						}
					}
				}
			}
			// valid considera a lista inteira antes da paginação, sem omitir diagnóstico após a primeira página.
			valid := true
			for _, item := range items {
				valid = valid && item.Severity != "error"
			}
			decimal := strconv.FormatUint(revision, 10)
			page, err := paginate(runtime.pages, "diagnostics/"+request.SessionID+"/"+request.Target, decimal, items, request.pageRequest)
			if err != nil {
				return nil, err
			}
			return successResult(validationResult{Target: request.Target, Readable: true, Valid: valid, Diagnostics: page, Limitations: []string{"Java/EFL não ofertado por MCP I1", "PPT é saída estática; fonte não é incorporada e homologação no Office permanece independente", "Nenhuma ação nova é inferida pelo Kind"}}, decimal), nil
		})
}
