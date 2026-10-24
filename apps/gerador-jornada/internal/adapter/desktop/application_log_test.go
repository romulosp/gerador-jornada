package desktop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gerador-jornada.local/editor/internal/application"
)

// TestApplicationLogWritesStartupAndRedacts confere JSONL e ausência de conteúdo do projeto.
func TestApplicationLogWritesStartupAndRedacts(t *testing.T) {
	// dir isola os arquivos do teste.
	dir := t.TempDir()
	// logger grava o startup no arquivo da sessão.
	logger, err := newApplicationLog(dir)
	if err != nil {
		t.Fatal(err)
	}
	logger.write("startup", map[string]any{"version": "teste", "pid": 1})
	logger.write("component_snapshot", map[string]any{"components": 2, "screens": 1})
	if logger.file != nil {
		_ = logger.file.Close()
	}
	// data contém as duas linhas JSONL.
	data, err := os.ReadFile(logger.path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], `"event":"startup"`) || !strings.Contains(lines[1], `"components":2`) {
		t.Fatalf("log inesperado: %s", data)
	}
	if strings.Contains(string(data), "Jornada de atendimento") || strings.Contains(string(data), "componente") {
		t.Fatal("log transcreveu conteúdo do projeto")
	}
}

// TestApplicationLogDedupAndForce confere a janela de dois segundos e a captura explícita.
func TestApplicationLogDedupAndForce(t *testing.T) {
	// dir isola o arquivo do teste.
	logger, err := newApplicationLog(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	// writeDedup grava a primeira captura e ignora a repetição imediata.
	logger.writeDedup("component_snapshot", map[string]any{"components": 1}, "tela/1", false, 2*time.Second)
	logger.writeDedup("component_snapshot", map[string]any{"components": 1}, "tela/1", false, 2*time.Second)
	// force autoriza a repetição mesmo dentro da janela.
	logger.writeDedup("component_snapshot", map[string]any{"components": 1}, "tela/1", true, 2*time.Second)
	if logger.file != nil {
		_ = logger.file.Close()
	}
	data, err := os.ReadFile(logger.path)
	if err != nil {
		t.Fatal(err)
	}
	if count := strings.Count(string(data), `"component_snapshot"`); count != 2 {
		t.Fatalf("dedup inesperado: %d", count)
	}
}

// TestApplicationLogRotation preserva backups e reabre o segmento corrente.
func TestApplicationLogRotation(t *testing.T) {
	// logger usa um limite pequeno para forçar rotação.
	logger, err := newApplicationLog(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	logger.limit = 200
	// index gera linhas suficientes para mais de um segmento.
	for index := 0; index < 12; index++ {
		logger.write("evento", map[string]any{"index": index})
	}
	if logger.file != nil {
		_ = logger.file.Close()
	}
	if _, err := os.Stat(logger.path + ".1"); err != nil {
		t.Fatalf("backup rotacionado ausente: %v", err)
	}
	if logger.size >= logger.limit {
		t.Fatal("segmento corrente não foi reaberto")
	}
}

// TestApplicationLogSnapshotOnlyCounts confere que a captura usa somente contagens.
func TestApplicationLogSnapshotOnlyCounts(t *testing.T) {
	// dir isola o arquivo e o logger global do teste.
	dir := t.TempDir()
	logger, err := newApplicationLog(dir)
	if err != nil {
		t.Fatal(err)
	}
	appLog.mu.Lock()
	previous := appLog.current
	appLog.current = logger
	appLog.mu.Unlock()
	defer func() {
		appLog.mu.Lock()
		appLog.current = previous
		appLog.mu.Unlock()
	}()
	document := application.Example()
	LogApplicationSnapshot(document, document.Screens[0].ID, true)
	if logger.file != nil {
		_ = logger.file.Close()
	}
	data, err := os.ReadFile(filepath.Join(dir, filepath.Base(logger.path)))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `"components":`) || !strings.Contains(text, `"kinds":`) {
		t.Fatalf("snapshot sem contagens: %s", text)
	}
	if strings.Contains(text, "Jornada de atendimento") || strings.Contains(text, "numeroAtendimento") {
		t.Fatal("snapshot transcreveu nome do projeto ou componente")
	}
}
