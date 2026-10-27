package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestServerFoundationUsesStrictBoundary verifica composição real SDK com a borda compartilhada.
func TestServerFoundationUsesStrictBoundary(t *testing.T) {
	// root/diagnostic isolam workspace e stderr, sem logger desktop nem preferências.
	root := t.TempDir()
	var diagnostic bytes.Buffer
	runtime, err := newRuntime(Config{Workspace: root}, &diagnostic)
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.authoring.CloseAll()
	server := newServer(runtime)
	clientTransport, serverTransport := sdk.NewInMemoryTransports()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatal(err)
	}
	client := sdk.NewClient(&sdk.Implementation{Name: "foundation-harness", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	// result conserva falha curta sem transcrever o campo extra ou seu valor.
	result, err := session.CallTool(ctx, &sdk.CallToolParams{Name: "aplicacao_capacidades", Arguments: map[string]any{"extra": "texto-pessoal"}})
	if err != nil || !result.IsError {
		t.Fatal("borda fechada não aplicada")
	}
	encoded, _ := json.Marshal(result)
	if bytes.Contains(encoded, []byte("texto-pessoal")) {
		t.Fatal("erro transcreveu argumento")
	}
	result, err = session.CallTool(ctx, &sdk.CallToolParams{Name: "aplicacao_capacidades", Arguments: map[string]any{}})
	if err != nil || result.IsError {
		t.Fatal("descoberta falhou")
	}
	encoded, _ = json.Marshal(result.StructuredContent)
	if !bytes.Contains(encoded, []byte(`"persistentWrites":false`)) || !bytes.Contains(encoded, []byte(`"startupProbe":false`)) {
		t.Fatal("estado anunciado não corresponde à fundação")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("startup/consulta criaram arquivo no workspace")
	}
	if !bytes.Contains(diagnostic.Bytes(), []byte(`"operation":"aplicacao_capacidades"`)) {
		t.Fatal("auditoria não foi composta")
	}
}
