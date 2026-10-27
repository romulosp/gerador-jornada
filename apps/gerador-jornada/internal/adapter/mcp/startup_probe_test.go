//go:build windows

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/sys/windows"
)

// TestDistributedStartupProbe usa o EXE Windows distribuído com pipes reais.
// JORNADA_MCP_PROBE_EXE seleciona o artefato; ausência não comprova o gate G0.
func TestDistributedStartupProbe(t *testing.T) {
	// executable vem do operador de teste, nunca de uma chamada MCP.
	executable := os.Getenv("JORNADA_MCP_PROBE_EXE")
	if executable == "" {
		t.Skip("G0 pendente: configure JORNADA_MCP_PROBE_EXE para o EXE Windows padrão")
	}
	testDistributedStartup(t, executable, true)
}

// TestDistributedFoundation usa a composição atual, sem confundir seu build com G0.
func TestDistributedFoundation(t *testing.T) {
	// executable seleciona explicitamente o build de fontes da fundação atual.
	executable := os.Getenv("JORNADA_MCP_EXE")
	if executable == "" {
		t.Skip("configure JORNADA_MCP_EXE para o EXE Windows atual")
	}
	testDistributedStartup(t, executable, false)
}

// testDistributedStartup comprova negociação/descoberta/EOF sem janela ou efeito persistente.
func testDistributedStartup(t *testing.T, executable string, probe bool) {
	t.Helper()
	// ctx limita todo o processo, incluindo negociação e EOF.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// root é o único workspace que o processo pode receber nesta prova.
	root := t.TempDir()
	// private isola preferências para detectar escrita incidental do desktop.
	private := filepath.Join(root, "config-nao-criada")
	// command conserva o artefato windowsgui e captura somente stderr separado.
	command := exec.CommandContext(ctx, executable, "--mcp", "--workspace", root)
	command.Env = append(os.Environ(), "JORNADA_CONFIG_DIR="+private)
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	// client negocia o perfil da SDK sem substituir o processo por chamada local.
	client := sdk.NewClient(&sdk.Implementation{Name: "jornada-g0-harness", Version: "1"}, nil)
	// session deve negociar, descobrir e consultar sem qualquer janela.
	session, err := client.Connect(ctx, &sdk.CommandTransport{Command: command}, nil)
	if err != nil {
		t.Fatalf("startup/negociação: %v; stderr=%s", err, diagnostic.String())
	}
	defer session.Close()
	// tools comprova oferta mínima real, sem anunciar operações ainda inexistentes.
	tools, err := session.ListTools(ctx, nil)
	if err != nil || len(tools.Tools) == 0 || tools.Tools[0].Name != "aplicacao_capacidades" || (probe && len(tools.Tools) != 1) {
		t.Fatalf("descoberta inesperada: %+v, %v", tools, err)
	}
	assertNoProcessWindows(t, uint32(command.Process.Pid))
	// result é lido do canal estruturado, não do texto livre da resposta.
	result, err := session.CallTool(ctx, &sdk.CallToolParams{Name: "aplicacao_capacidades", Arguments: map[string]any{}})
	if err != nil || result.IsError {
		t.Fatalf("capacidade: %+v, %v", result, err)
	}
	// raw/decoded normalizam StructuredContent sem depender do tipo usado pela SDK.
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		// OK identifica sucesso da ferramenta, separado de disponibilidade R2.
		OK bool `json:"ok"`
		// Result contém o subconjunto de capacidades exigido pelo gate.
		Result struct {
			// PersistentWrites deve permanecer falso na prova sem host.
			PersistentWrites bool `json:"persistentWrites"`
			// StartupProbe diferencia G0 de autoria I1 homologada.
			StartupProbe bool `json:"startupProbe"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil || !decoded.OK || decoded.Result.StartupProbe != probe || decoded.Result.PersistentWrites {
		t.Fatalf("capacidades incompatíveis: %s (%v)", raw, err)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if command.ProcessState == nil || !command.ProcessState.Success() {
		t.Fatalf("processo não encerrou normalmente após EOF: %v", command.ProcessState)
	}
	if _, err := os.Stat(private); !os.IsNotExist(err) {
		t.Fatalf("preferências criadas no headless: %v", err)
	}
	// entries detecta arquivo/log/rascunho persistido sem solicitação.
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("efeito persistente incidental: %+v (%v)", entries, err)
	}
}

// TestDistributedMCPRejectsMixedFlags recusa mistura de modalidades antes de qualquer inicialização desktop.
func TestDistributedMCPRejectsMixedFlags(t *testing.T) {
	// executable precisa ser a revisão atual, não o artefato G0 sem composição completa.
	executable := os.Getenv("JORNADA_MCP_EXE")
	if executable == "" {
		t.Skip("configure JORNADA_MCP_EXE para o EXE Windows atual")
	}
	// name é só identificação do cenário; args nunca vêm de conteúdo de projeto.
	for _, name := range []string{"open", "save-example", "root-without-mcp"} {
		t.Run(name, func(t *testing.T) {
			// root/private/output detectam qualquer efeito incidental do caminho incompatível.
			root := t.TempDir()
			private := filepath.Join(root, "private")
			output := filepath.Join(root, "indevido.jornada")
			args := []string{"--mcp", "--workspace", root, "--" + name, output}
			if name == "root-without-mcp" {
				args = []string{"--legacy-root", root}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, executable, args...)
			command.Env = append(os.Environ(), "JORNADA_CONFIG_DIR="+private)
			if err := command.Run(); err == nil || ctx.Err() != nil {
				t.Fatal("flags incompatíveis não recusadas imediatamente")
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatal("flag incompatível causou efeito persistente")
			}
		})
	}
}

// assertNoProcessWindows consulta janelas Win32 do processo sem controlar a GUI.
// Mesmo uma janela oculta é rejeitada: a prova MCP não deve inicializar desktop.
func assertNoProcessWindows(t *testing.T, processID uint32) {
	t.Helper()
	// user32 fornece somente consultas locais de identidade das janelas.
	user32 := windows.NewLazySystemDLL("user32.dll")
	// enumWindows/getProcess relacionam cada janela de topo ao processo servidor.
	enumWindows := user32.NewProc("EnumWindows")
	getProcess := user32.NewProc("GetWindowThreadProcessId")
	// found é preenchido sincronamente pelo callback Win32 e conserva a janela detectada.
	found := uintptr(0)
	// window/argument são o handle de topo e o contexto não utilizado da enumeração.
	callback := windows.NewCallback(func(window, argument uintptr) uintptr {
		// owner recebe o PID sem criar ou alterar janela.
		var owner uint32
		getProcess.Call(window, uintptr(unsafe.Pointer(&owner)))
		if owner == processID {
			found = window
			return 0
		}
		return 1
	})
	// status/callError impedem interpretar falha da consulta como ausência de janela.
	status, _, callError := enumWindows.Call(callback, 0)
	if status == 0 && found == 0 {
		t.Fatalf("não foi possível consultar janelas Win32: %v", callError)
	}
	if found != 0 {
		t.Fatalf("modo MCP criou janela Win32: %x", found)
	}
}
