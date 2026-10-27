//go:build windows

package mcp

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/project"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestDistributedSessionControls comprova controles e EOF no mesmo executável Windows distribuído.
func TestDistributedSessionControls(t *testing.T) {
	// executable vem do operador; ausência não conta como aprovação deste gate.
	executable := os.Getenv("JORNADA_MCP_EXE")
	if executable == "" {
		t.Skip("configure JORNADA_MCP_EXE para a revisão atual")
	}
	// ctx/root isolam processo e destino sem compartilhar dados de trabalho do usuário.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	root := t.TempDir()
	command := exec.CommandContext(ctx, executable, "--mcp", "--workspace", root)
	command.Env = append(os.Environ(), "JORNADA_CONFIG_DIR="+filepath.Join(root, "private-nao-criada"))
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	client := sdk.NewClient(&sdk.Implementation{Name: "distributed-session-harness", Version: "1"}, nil)
	session, err := client.Connect(ctx, &sdk.CommandTransport{Command: command}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	// created deve anunciar rascunho pendente e exigir descarte para encerrar.
	created := callResult[admissionResult](t, session, "jornada_criar", map[string]any{"requestId": project.NewID(), "name": "Controle distribuído 😀", "mode": "blank"})
	if _, err := time.Parse(time.RFC3339, created.ExpiresAt); err != nil || created.IdleSeconds != 1800 || created.UnsavedPolicy != "explicitDiscardRequired;idleOrEOFDiscardsWithoutSaving" {
		t.Fatal("admissão não informou prazo/política de descarte")
	}
	before := callResult[sessionObservation](t, session, "sessao_obter", map[string]any{"sessionId": created.SessionID})
	if !before.Dirty || !before.DocumentDirty || before.ProposalDirty || before.Quotas.ReservedTargets != 1 {
		t.Fatal("rascunho sem estado/reserva declarados")
	}
	// args é preservado integralmente para demonstrar replay terminal, sem recompor CAS.
	args := sessionControlArgs(t, session, created.SessionID)
	args["discardUnsaved"] = false
	callFailure(t, session, "sessao_encerrar", args, "invalidArgument")
	args["discardUnsaved"] = true
	closed := callResult[sessionControlResult](t, session, "sessao_encerrar", args)
	replay := callResult[sessionControlResult](t, session, "sessao_encerrar", args)
	if closed != replay || !closed.Closed || !closed.Discarded {
		t.Fatal("fechamento/replay não conservou resultado")
	}
	callFailure(t, session, "sessao_obter", map[string]any{"sessionId": created.SessionID}, "sessionExpired")
	// capabilities continua recusando qualificação R2; controles não concedem consentimento.
	capabilities := callResult[installedCapabilities](t, session, "aplicacao_capacidades", map[string]any{})
	if capabilities.PersistentWrites || capabilities.ApprovalProfile.Qualified {
		t.Fatal("controle qualificou host ou habilitou gravação")
	}
	assertNoProcessWindows(t, uint32(command.Process.Pid))
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if command.ProcessState == nil || !command.ProcessState.Success() {
		t.Fatal("EOF não terminou normalmente")
	}
	// entries garante que encerramento não persistiu rascunho, configuração ou journal.
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("controle transitório produziu arquivos")
	}
	if strings.Contains(diagnostic.String(), "Controle distribuído") {
		t.Fatal("auditoria expôs nome do projeto")
	}
}

// TestDistributedAuthoringDraft confere autoria transitória no EXE windowsgui com pipes reais.
func TestDistributedAuthoringDraft(t *testing.T) {
	// executable seleciona o artefato pelo operador; ausência pula teste, sem comprovar o gate.
	executable := os.Getenv("JORNADA_MCP_EXE")
	if executable == "" {
		t.Skip("configure JORNADA_MCP_EXE para a revisão atual")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	root := t.TempDir()
	private := filepath.Join(root, "private-nao-criada")
	command := exec.CommandContext(ctx, executable, "--mcp", "--workspace", root)
	command.Env = append(os.Environ(), "JORNADA_CONFIG_DIR="+private)
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	client := sdk.NewClient(&sdk.Implementation{Name: "distributed-authoring-harness", Version: "1"}, nil)
	session, err := client.Connect(ctx, &sdk.CommandTransport{Command: command}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	// created/overview vêm somente do protocolo, sem acesso ao registro em memória do processo.
	created := callResult[admissionResult](t, session, "jornada_criar", map[string]any{"requestId": project.NewID(), "name": "Distribuída 😀", "mode": "blank"})
	overview := callResult[pageResult[screenSummary]](t, session, "jornada_obter", map[string]any{"sessionId": created.SessionID})
	if len(overview.Items) != 1 {
		t.Fatal("rascunho inicial sem tela")
	}
	pageID := overview.Items[0].ID
	button := callResult[struct {
		ComponentID string `json:"componentId"`
	}](t, session, "componente_inserir", map[string]any{"sessionId": created.SessionID, "requestId": project.NewID(), "expectedRevision": overview.Revision, "screenId": pageID, "kind": "button", "x": 20, "y": 20})
	overview = callResult[pageResult[screenSummary]](t, session, "jornada_obter", map[string]any{"sessionId": created.SessionID})
	second := callResult[struct {
		ScreenID string `json:"screenId"`
	}](t, session, "tela_criar", map[string]any{"sessionId": created.SessionID, "requestId": project.NewID(), "expectedRevision": overview.Revision, "name": "Conclusão"})
	overview = callResult[pageResult[screenSummary]](t, session, "jornada_obter", map[string]any{"sessionId": created.SessionID})
	callResult[struct{}](t, session, "navegacao_vincular", map[string]any{"sessionId": created.SessionID, "requestId": project.NewID(), "expectedRevision": overview.Revision, "sourceScreenId": pageID, "sourceComponentId": button.ComponentID, "targetScreenId": second.ScreenID, "action": "navigate"})
	validation := callResult[validationResult](t, session, "jornada_validar", map[string]any{"sessionId": created.SessionID, "target": "draft"})
	if !validation.Readable || !validation.Valid {
		t.Fatal("rascunho produzido pelo processo não é válido")
	}
	capabilities := callResult[installedCapabilities](t, session, "aplicacao_capacidades", map[string]any{})
	if capabilities.PersistentWrites || capabilities.ApprovalProfile.Qualified {
		t.Fatal("autoria transitória qualificou gravação indevidamente")
	}
	assertNoProcessWindows(t, uint32(command.Process.Pid))
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if command.ProcessState == nil || !command.ProcessState.Success() {
		t.Fatal("EOF não encerrou normalmente")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("autoria transitória persistiu jornada/preferências sem consentimento")
	}
	if strings.Contains(diagnostic.String(), "Distribuída") || strings.Contains(diagnostic.String(), "Conclusão") {
		t.Fatal("auditoria expôs conteúdo de projeto")
	}
}
