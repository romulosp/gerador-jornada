package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/adapter/jsonwire"
	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxJSONDepth preserva menus de 16 níveis sem permitir recursão JSON ilimitada.
const maxJSONDepth = 64

// argumentSchemas limita inferência às structs de operação registradas pelo aplicativo.
// Chaves são tipos Go fixos, nunca schemas enviados pelo agente.
var argumentSchemas sync.Map

// schemaFor infere contrato fechado uma vez por DTO e o compartilha somente para leitura.
func schemaFor[T any]() *jsonschema.Resolved {
	// kind é estável durante o processo e não contém conteúdo do cliente.
	kind := reflect.TypeFor[T]()
	if prior, found := argumentSchemas.Load(kind); found {
		return prior.(*jsonschema.Resolved)
	}
	// schema/resolved rejeitam definições inválidas na composição, sem aceitar fallback livre.
	schema, err := jsonschema.For[T](componentSchemaOptions())
	if err != nil {
		panic("DTO MCP sem schema válido")
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		panic("schema MCP não resolvido")
	}
	prior, _ := argumentSchemas.LoadOrStore(kind, resolved)
	return prior.(*jsonschema.Resolved)
}

// frameReader valida uma linha inteira antes de entregar qualquer byte à SDK.
type frameReader struct {
	// source permite interromper leitura bloqueada quando a conexão é encerrada.
	source io.ReadCloser
	// reader usa buffer pequeno; fragmentos de um único frame têm orçamento explícito.
	reader *bufio.Reader
	// pending conserva apenas o frame já validado ainda não consumido pela SDK.
	pending []byte
	// failed torna erro de framing terminal, sem aceitar fragmento como mensagem nova.
	failed error
}

// newFrameReader compõe leitura limitada sem goroutine ou configuração externa.
func newFrameReader(source io.ReadCloser) *frameReader {
	return &frameReader{source: source, reader: bufio.NewReaderSize(source, 32<<10)}
}

// Read só fornece bytes depois de confirmar limite, LF, JSON e Unicode.
func (r *frameReader) Read(destination []byte) (int, error) {
	if len(destination) == 0 {
		return 0, nil
	}
	if r.failed != nil {
		return 0, r.failed
	}
	if len(r.pending) == 0 {
		// frame cresce no máximo até MaxFrameBytes, antes de qualquer decodificação.
		frame := make([]byte, 0, 32<<10)
		for {
			// fragment/err distinguem buffer cheio de EOF e falha real do pipe.
			fragment, err := r.reader.ReadSlice('\n')
			if len(frame)+len(fragment) > MaxFrameBytes {
				r.failed = fmt.Errorf("resourceLimit: frame MCP excede 1 MiB")
				return 0, r.failed
			}
			frame = append(frame, fragment...)
			if err == bufio.ErrBufferFull {
				continue
			}
			if err != nil {
				if err == io.EOF && len(frame) > 0 {
					err = fmt.Errorf("invalidArgument: frame MCP sem LF final")
				}
				r.failed = err
				return 0, err
			}
			if err := validateJSON(frame); err != nil {
				r.failed = err
				return 0, err
			}
			if err := validateRPCID(frame); err != nil {
				r.failed = err
				return 0, err
			}
			r.pending = frame
			break
		}
	}
	// copied informa somente bytes já validados do frame corrente.
	copied := copy(destination, r.pending)
	r.pending = r.pending[copied:]
	return copied, nil
}

// Close interrompe o pipe e não publica conteúdo parcial.
func (r *frameReader) Close() error { return r.source.Close() }

// validateJSON reutiliza a mesma política de Unicode/chaves do codec de pacote, com orçamento do frame MCP.
func validateJSON(data []byte) error { return jsonwire.Validate(data, MaxFrameBytes, maxJSONDepth) }

// validateRPCID aceita notificações sem ID e IDs canônicos pequenos/seguros.
func validateRPCID(frame []byte) error {
	// envelope mantém a ID em bytes até conferir seu tipo/intervalo.
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(frame, &envelope); err != nil || envelope == nil {
		return fmt.Errorf("invalidArgument: envelope JSON-RPC deve ser objeto")
	}
	// id/existing separam ausência legítima em notificação de ID explicitamente inválida.
	id, existing := envelope["id"]
	if !existing {
		return nil
	}
	if len(id) > 0 && id[0] == '"' {
		// text já possui Unicode válido, confirmado antes deste parse.
		var text string
		if json.Unmarshal(id, &text) != nil || len(text) > 128 {
			return fmt.Errorf("invalidArgument: ID JSON-RPC excede 128 bytes")
		}
		return nil
	}
	// value exige inteiro decimal, sem fração/expoente e sem perda em clientes JavaScript.
	value, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil || value > 9007199254740991 || value < -9007199254740991 {
		return fmt.Errorf("invalidArgument: ID JSON-RPC não é inteiro seguro")
	}
	return nil
}

// decodeArguments converte objeto fechado para DTO específico, nunca map de domínio.
func decodeArguments[T any](raw []byte) (T, error) {
	// result só é entregue quando todos os campos e bytes foram conferidos.
	var result T
	if err := validateJSON(raw); err != nil {
		return result, err
	}
	if _, strict := any(result).(interface{ noNullInput() }); strict {
		// tokens distingue null de uma string literal, inclusive em patches aninhados.
		tokens := json.NewDecoder(bytes.NewReader(raw))
		for {
			value, err := tokens.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return result, fmt.Errorf("invalidArgument: objeto inválido")
			}
			if value == nil {
				return result, fmt.Errorf("invalidArgument: null não apaga propriedades")
			}
		}
	}
	// instance confere nomes exatos e campos obrigatórios antes do decode case-insensitive de Go.
	var instance any
	if err := json.Unmarshal(raw, &instance); err != nil || schemaFor[T]().Validate(instance) != nil {
		return result, fmt.Errorf("invalidArgument: objeto fora do contrato")
	}
	// decoder rejeita propriedade desconhecida inclusive em structs aninhadas.
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("invalidArgument: propriedades ou tipos inválidos")
	}
	return result, nil
}

// callBudget aplica concorrência/deadline sem ocupar slot aguardando aprovação humana.
func callBudget(next sdk.MethodHandler) sdk.MethodHandler {
	// slots pertence à instância do servidor, não a uma configuração enviada pela LLM.
	slots := make(chan struct{}, 4)
	// ctx/method/request são a chamada já identificada pela SDK e seu cancelamento.
	return func(ctx context.Context, method string, request sdk.Request) (sdk.Result, error) {
		if method != "tools/call" && method != "resources/read" {
			return next(ctx, method, request)
		}
		select {
		case slots <- struct{}{}:
		default:
			if method == "tools/call" {
				return failureCode("busy"), nil
			}
			return nil, fmt.Errorf("busy: quatro chamadas em andamento")
		}
		defer func() { <-slots }()
		// timeout é fixo por classe de trabalho, sem admitir que o cliente amplie o prazo.
		timeout := 30 * time.Second
		if tool, ok := request.(*sdk.CallToolRequest); ok && tool.Params != nil {
			switch tool.Params.Name {
			case "componente_preparar_importacao", "componente_revisar", "jornada_revisar_alteracoes":
				timeout = 120 * time.Second
			}
		}
		// bounded conserva também o prazo menor/cancelamento já recebido.
		bounded, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return next(bounded, method, request)
	}
}
