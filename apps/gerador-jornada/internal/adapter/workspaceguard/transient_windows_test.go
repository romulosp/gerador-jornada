//go:build windows

package workspaceguard

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

// TestTransientFileIsImmutableAndDeletedOnClose conserva staging privado sem efeito definitivo.
func TestTransientFileIsImmutableAndDeletedOnClose(t *testing.T) {
	// root/guard ficam fora de qualquer workspace real e usam exclusão operacional de teste.
	root := t.TempDir()
	guard, _ := New(t.TempDir())
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	file, err := lease.CreateTransient("staging.tmp")
	lease.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := file.Write([]byte("frozen")); err != nil {
		t.Fatal(err)
	}
	if writer, err := os.OpenFile(filepath.Join(root, "staging.tmp"), os.O_WRONLY, 0); err == nil {
		writer.Close()
		t.Fatal("outro escritor alterou staging aberto")
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "staging.tmp")); !os.IsNotExist(err) {
		t.Fatal("temporário sobreviveu ao último handle")
	}
}

// TestTransientKernelCleanup encerra subprocesso sem defer/Close para comprovar descarte pelo SO.
func TestTransientKernelCleanup(t *testing.T) {
	// root/config são as únicas pastas fornecidas ao auxiliar descartável.
	root, config := t.TempDir(), t.TempDir()
	command := exec.Command(os.Args[0], "-test.run=^TestTransientChild$")
	command.Env = append(os.Environ(), "JORNADA_TRANSIENT_CHILD_ROOT="+root, "JORNADA_TRANSIENT_CHILD_CONFIG="+config)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("subprocesso: %s (%v)", output, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatal("queda deixou staging persistente")
	}
}

// TestTransientChild só encerra seu próprio processo quando executado como auxiliar explícito.
func TestTransientChild(t *testing.T) {
	// root vazio distingue a suíte normal do processo nativo descartável.
	root := os.Getenv("JORNADA_TRANSIENT_CHILD_ROOT")
	if root == "" {
		t.Skip("auxiliar de teste")
	}
	guard, err := New(os.Getenv("JORNADA_TRANSIENT_CHILD_CONFIG"))
	if err != nil {
		t.Fatal(err)
	}
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	file, err := lease.CreateTransient("crash.tmp")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte("private")); err != nil {
		t.Fatal(err)
	}
	if err := file.Sync(); err != nil {
		t.Fatal(err)
	}
	// ExitProcess evita os defers Go somente deste auxiliar; kernel fecha seus handles.
	windows.ExitProcess(0)
}

// TestBaselineContextMatchesSnapshot evita ler pacote inteiro ao revalidar precondição de revisão.
func TestBaselineContextMatchesSnapshot(t *testing.T) {
	// root/lease oferecem bytes estáveis a duas APIs com o mesmo significado de identidade/hash.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "base.jornada"), []byte("baseline"), 0600); err != nil {
		t.Fatal(err)
	}
	guard, _ := New(t.TempDir())
	lease, err := guard.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	_, expected, err := lease.Snapshot("base.jornada")
	if err != nil {
		t.Fatal(err)
	}
	actual, err := lease.BaselineContext(context.Background(), "base.jornada")
	if err != nil || expected != actual {
		t.Fatalf("precondições divergentes: %+v %+v (%v)", expected, actual, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := lease.BaselineContext(ctx, "base.jornada"); err != context.Canceled {
		t.Fatal("hash ignorou cancelamento")
	}
	if _, err := lease.CreateTransient("forbidden.tmp"); err == nil {
		t.Fatal("lease de leitura criou staging")
	}
}
