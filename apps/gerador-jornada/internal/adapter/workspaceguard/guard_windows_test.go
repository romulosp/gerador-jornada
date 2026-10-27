//go:build windows

package workspaceguard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestGuardExclusiveCreation preserva homônimo e bytes completos em publicação exclusiva.
func TestGuardExclusiveCreation(t *testing.T) {
	// root/config isolam material do usuário e arquivos operacionais privados.
	root, config := t.TempDir(), t.TempDir()
	guard, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	if err := lease.CreateNew(context.Background(), "teste.jornada", []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := lease.CreateNew(context.Background(), "TESTE.jornada", []byte("overwrite")); !errors.Is(err, os.ErrExist) {
		t.Fatalf("homônimo sobrescrito: %v", err)
	}
	// data/base devem descrever exatamente o mesmo handle de arquivo lido.
	data, base, err := lease.Snapshot("teste.jornada")
	if err != nil || !bytes.Equal(data, []byte("first")) || base.SizeBytes != 5 || base.Identity == "" || len(base.SHA256) != 64 {
		t.Fatalf("snapshot: %q %+v (%v)", data, base, err)
	}
	// name cobre nomes Windows que não podem chegar ao filesystem como caminhos.
	for _, name := range []string{"../escape.jornada", `sub\escape.jornada`, "NUL.jornada", "COM1.jornada", "conout$.jornada", "with:stream.jornada", "trailing. ", ".."} {
		if err := lease.CreateNew(context.Background(), name, []byte("no")); err == nil {
			t.Fatalf("nome inválido aceito: %q", name)
		}
	}
	// canceled não deixa temporário ou arquivo final após cancelamento anterior à publicação.
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := lease.CreateNew(canceled, "cancelado.jornada", []byte("no")); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelamento: %v", err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatalf("resíduo de criação: %+v (%v)", entries, err)
	}
}

// TestGuardProcessSerialization comprova exclusão entre processos e aliases de caixa.
func TestGuardProcessSerialization(t *testing.T) {
	// root/config são enviados somente ao subprocesso de teste explícito.
	root, config := t.TempDir(), t.TempDir()
	guard, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	if second, err := guard.Acquire(strings.ToUpper(root)); err == nil {
		second.Close()
		t.Fatal("alias adquiriu a mesma raiz")
	}
	// ctx limita falha de um subprocesso sem introduzir esperas indefinidas.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGuardChild$")
	command.Env = append(os.Environ(), "JORNADA_GUARD_CHILD_ROOT="+root, "JORNADA_GUARD_CHILD_CONFIG="+config)
	output, err := command.CombinedOutput()
	if err != nil || !bytes.Contains(output, []byte("BUSY")) {
		t.Fatalf("trava entre processos: %s (%v)", output, err)
	}
	lease.Close()
	second, err := guard.Acquire(root)
	if err != nil {
		t.Fatalf("trava não liberada: %v", err)
	}
	second.Close()
}

// TestGuardChild só é executado como processo auxiliar pelo próprio teste acima.
func TestGuardChild(t *testing.T) {
	// root seleciona a fixture; execução normal da suíte não cria processo auxiliar.
	root := os.Getenv("JORNADA_GUARD_CHILD_ROOT")
	if root == "" {
		t.Skip("auxiliar de teste")
	}
	guard, err := New(os.Getenv("JORNADA_GUARD_CHILD_CONFIG"))
	if err != nil {
		t.Fatal(err)
	}
	lease, err := guard.Acquire(root)
	if err == nil {
		lease.Close()
		t.Fatal("subprocesso ignorou exclusão")
	}
	fmt.Println("BUSY")
}

// TestGuardNoFollow rejeita links de arquivo sem revelar bytes externos.
func TestGuardNoFollow(t *testing.T) {
	// root/outside conservam dados e link em diretórios distintos.
	root, outside := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "outside.jornada"), []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside.jornada"), filepath.Join(root, "link.jornada")); err != nil {
		// Junction não exige o privilégio de symlink e exercita reparse point real.
		t.Logf("symlink de arquivo indisponível; conferindo junction: %v", err)
		makeTestJunction(t, filepath.Join(root, "link.jornada"), outside)
	}
	guard, _ := New(t.TempDir())
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	if data, _, err := lease.Snapshot("link.jornada"); err == nil || len(data) != 0 {
		t.Fatalf("leitura seguiu link: %q (%v)", data, err)
	}
}

// makeTestJunction cria somente uma fixture de reparse point em diretórios temporários.
func makeTestJunction(t *testing.T, link, target string) {
	t.Helper()
	// command usa script fixo e argumentos de caminho como dados de ambiente.
	command := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; New-Item -ItemType Junction -Path $env:JORNADA_TEST_LINK -Target $env:JORNADA_TEST_TARGET | Out-Null")
	command.Env = append(os.Environ(), "JORNADA_TEST_LINK="+link, "JORNADA_TEST_TARGET="+target)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("fixture de junction: %s (%v)", output, err)
	}
}

// TestGuardRelativeDirectory impede redirecionar publicação ao trocar o alias da raiz.
func TestGuardRelativeDirectory(t *testing.T) {
	// parent mantém inclusive a pasta renomeada dentro do diretório temporário do teste.
	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	moved := filepath.Join(parent, "moved")
	outside := filepath.Join(parent, "outside")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(outside, 0700); err != nil {
		t.Fatal(err)
	}
	guard, _ := New(t.TempDir())
	lease, err := guard.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	if err := os.Rename(root, moved); err != nil {
		t.Fatalf("fixture de troca de raiz: %v", err)
	}
	makeTestJunction(t, root, outside)
	if err := lease.CreateNew(context.Background(), "original.jornada", []byte("complete")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(outside, "original.jornada")); !os.IsNotExist(err) {
		t.Fatalf("escrita saiu da raiz física: %v", err)
	}
	// data comprova que o mesmo handle físico foi utilizado, sem resolver novamente o alias.
	data, err := os.ReadFile(filepath.Join(moved, "original.jornada"))
	if err != nil || string(data) != "complete" {
		t.Fatalf("destino físico: %q (%v)", data, err)
	}
	// entries precisa enumerar o mesmo diretório aberto, inclusive após troca do caminho.
	entries, err := lease.Entries()
	if err != nil || len(entries) != 1 || entries[0].Name() != "original.jornada" {
		t.Fatalf("inventário saiu da raiz física: %v (%v)", entries, err)
	}
}

// TestGuardDefaultNamespace impede separar o gate alterando configurações por cliente.
func TestGuardDefaultNamespace(t *testing.T) {
	// first/second precisam escolher a mesma pasta conhecida do SO, sem criar preferência.
	first, err := NewDefault()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("JORNADA_CONFIG_DIR", t.TempDir())
	second, err := NewDefault()
	if err != nil {
		t.Fatal(err)
	}
	if first.PrivateRoot != second.PrivateRoot {
		t.Fatal("configuração de cliente mudou o namespace do gate")
	}
	// root é temporária; somente a trava opaca é criada na configuração operacional.
	root := t.TempDir()
	lease, err := first.Acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()
	if other, err := second.Acquire(root); err == nil {
		other.Close()
		t.Fatal("outro perfil ignorou exclusão")
	}
	// reader não pode aproveitar um handle R0 para publicar sem passar por Acquire.
	reader, err := first.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := reader.CreateNew(context.Background(), "blocked.jornada", []byte("no")); !errors.Is(err, os.ErrPermission) {
		t.Fatalf("consulta obteve autoridade de escrita: %v", err)
	}
}
