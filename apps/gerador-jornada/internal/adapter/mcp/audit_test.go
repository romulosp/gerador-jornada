package mcp

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestAuditAllowlist não transcreve nome arbitrário de ferramenta, erro ou revisão do cliente.
func TestAuditAllowlist(t *testing.T) {
	// output substitui somente stderr de teste, sem abrir log de arquivo.
	var output bytes.Buffer
	audit := newAudit(&output)
	audit.record("workspace_listar_jornadas", "ok", "12", 3*time.Millisecond)
	audit.record("texto-pessoal token=segredo\n", "C:\\privado", "texto-pessoal", time.Second)
	if strings.Contains(output.String(), "pessoal") || strings.Contains(output.String(), "segredo") || strings.Contains(output.String(), "privado") {
		t.Fatal("dado livre entrou na auditoria")
	}
	if !strings.Contains(output.String(), "workspace_listar_jornadas") || !strings.Contains(output.String(), `"revision":"12"`) {
		t.Fatal("metadados operacionais perdidos")
	}
	if strings.Count(output.String(), "\n") != 2 {
		t.Fatal("registro sofreu injeção de linha")
	}
}
