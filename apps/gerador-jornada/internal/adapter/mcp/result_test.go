package mcp

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// TestResultKeepsReviewInStructuredContent evita duplicar uma revisão grande no texto de compatibilidade.
func TestResultKeepsReviewInStructuredContent(t *testing.T) {
	// review simula uma revisão serializada de 512 KiB, além do envelope de negócio.
	review := json.RawMessage(`{"summary":"` + strings.Repeat("x", (512<<10)-14) + `"}`)
	result := successResult(review, "18446744073709551615")
	if result.IsError || len(result.Content) != 1 {
		t.Fatal("revisão no limite foi recusada")
	}
	// wire contém também ID e envelope JSON-RPC reais, já com StructuredContent serializado.
	wire, err := json.Marshal(struct {
		JSONRPC, ID string
		Result      any
	}{"2.0", strings.Repeat("a", 128), result})
	if err != nil || len(wire)+1 > MaxFrameBytes {
		t.Fatalf("frame excedeu limite: %d (%v)", len(wire), err)
	}
	text, err := json.Marshal(result.Content)
	if err != nil || len(text) > 1024 || strings.Contains(string(text), strings.Repeat("x", 100)) {
		t.Fatal("Content duplicou StructuredContent")
	}
}

// TestResultRefusesOversizedObject recusa antes da escrita, sem retornar conteúdo truncado.
func TestResultRefusesOversizedObject(t *testing.T) {
	// payload ultrapassa o orçamento útil, independentemente de caber em um item lógico.
	payload := struct {
		Text string `json:"text"`
	}{strings.Repeat("x", MaxFrameBytes)}
	result := successResult(payload, "")
	if !result.IsError {
		t.Fatal("resposta sem orçamento aceita")
	}
	encoded, _ := json.Marshal(result)
	if len(encoded) > 1024 || !strings.Contains(string(encoded), "resourceLimit") {
		t.Fatal("erro não foi curto e identificável")
	}
}

// TestSafeFailureDoesNotEchoSource mantém conteúdo pessoal fora de erros de protocolo/aplicação.
func TestSafeFailureDoesNotEchoSource(t *testing.T) {
	// source representa erro interno contendo caminho, conteúdo e credencial.
	source := errors.New("C:\\privado\\texto-do-projeto token=segredo")
	encoded, _ := json.Marshal(failureResult(source))
	if strings.Contains(string(encoded), "privado") || strings.Contains(string(encoded), "segredo") {
		t.Fatal("erro interno foi exposto")
	}
}
