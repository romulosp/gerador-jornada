package mcp

import (
	"context"
	"sync"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestCallBudgetKeepsCancellationAvailable não consome quinta vaga nem bloqueia cancelamento.
func TestCallBudgetKeepsCancellationAvailable(t *testing.T) {
	// entered/release seguram exatamente quatro chamadas sem temporizações artificiais.
	entered := make(chan struct{}, 4)
	release := make(chan struct{})
	var pending sync.WaitGroup
	budget := callBudget(func(ctx context.Context, method string, request sdk.Request) (sdk.Result, error) {
		if method == "notifications/cancelled" {
			return nil, nil
		}
		entered <- struct{}{}
		<-release
		return nil, nil
	})
	for index := 0; index < 4; index++ {
		pending.Add(1)
		go func() {
			defer pending.Done()
			_, _ = budget(context.Background(), "tools/call", &sdk.CallToolRequest{Params: &sdk.CallToolParamsRaw{Name: "aplicacao_capacidades"}})
		}()
	}
	for index := 0; index < 4; index++ {
		<-entered
	}
	// fifth deve retornar busy imediatamente, enquanto cancelamento continua passando.
	fifth, err := budget(context.Background(), "tools/call", &sdk.CallToolRequest{})
	if err != nil || fifth == nil || !fifth.(*sdk.CallToolResult).IsError {
		t.Fatal("quinta chamada não retornou busy")
	}
	if _, err := budget(context.Background(), "notifications/cancelled", nil); err != nil {
		t.Fatal("cancelamento bloqueado por quota")
	}
	close(release)
	pending.Wait()
}

// TestCallBudgetSetsDeadlines mantém prazos de trabalho distintos da espera humana.
func TestCallBudgetSetsDeadlines(t *testing.T) {
	// tool/expected cobrem operação comum e preparação de pacote.
	for _, test := range []struct {
		tool     string
		expected time.Duration
	}{{"aplicacao_capacidades", 30 * time.Second}, {"jornada_revisar_alteracoes", 120 * time.Second}} {
		budget := callBudget(func(ctx context.Context, _ string, _ sdk.Request) (sdk.Result, error) {
			// deadline precisa existir no contexto efetivo recebido pelo handler.
			deadline, found := ctx.Deadline()
			remaining := time.Until(deadline)
			if !found || remaining > test.expected || remaining < test.expected-time.Second {
				t.Fatalf("prazo incorreto: %s", remaining)
			}
			return nil, nil
		})
		if _, err := budget(context.Background(), "tools/call", &sdk.CallToolRequest{Params: &sdk.CallToolParamsRaw{Name: test.tool}}); err != nil {
			t.Fatal(err)
		}
	}
}
