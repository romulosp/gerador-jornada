package mcp

import (
	"context"
	"encoding/json"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addTypedTool mantém DTO fechado e erros seguros antes do handler, sem fallback textual da SDK.
// O handler recebe apenas argumentos validados e nunca uma solicitação de aprovação presumida.
func addTypedTool[In any](server *sdk.Server, recorder *audit, tool *sdk.Tool, handler func(context.Context, In) (*sdk.CallToolResult, error)) {
	// definition copia metadados para não modificar objetos já oferecidos por outro servidor.
	definition := *tool
	definition.InputSchema = schemaFor[In]().Schema()
	server.AddTool(&definition, func(ctx context.Context, request *sdk.CallToolRequest) (*sdk.CallToolResult, error) {
		// started contabiliza execução sem registrar argumentos ou documento.
		started := time.Now()
		// arguments assume objeto vazio apenas quando o protocolo omite o objeto inteiro.
		arguments := request.Params.Arguments
		if len(arguments) == 0 {
			arguments = json.RawMessage(`{}`)
		}
		// advanced só classifica recusa de propriedades conhecidas ainda não ofertadas em I1.
		if definition.Name == "componente_atualizar" {
			var advanced struct {
				Changes struct {
					Properties struct {
						Kind string `json:"kind"`
					} `json:"properties"`
				} `json:"changes"`
			}
			if json.Unmarshal(arguments, &advanced) == nil {
				switch advanced.Changes.Properties.Kind {
				case "barcodeInput", "barcodeProcessed", "applicationMenu":
					recorder.record(definition.Name, "capabilityUnavailable", "", time.Since(started))
					return failureCode("capabilityUnavailable"), nil
				}
			}
		}
		input, err := decodeArguments[In](arguments)
		var result *sdk.CallToolResult
		if err == nil {
			err = ctx.Err()
		}
		if err == nil {
			result, err = handler(ctx, input)
		}
		if err != nil {
			result = failureResult(err)
		}
		if result == nil {
			result = failureCode("capabilityUnavailable")
		}
		// code/revision vêm somente do envelope produzido pela aplicação; auditoria revalida ambos.
		code, revision := "ok", ""
		if encoded, marshalErr := json.Marshal(result.StructuredContent); marshalErr == nil {
			var envelope struct {
				Revision string      `json:"revision"`
				Error    errorDetail `json:"error"`
			}
			if json.Unmarshal(encoded, &envelope) == nil {
				revision = envelope.Revision
				if result.IsError {
					code = envelope.Error.Code
				}
			}
		}
		recorder.record(definition.Name, code, revision, time.Since(started))
		return result, nil
	})
}
