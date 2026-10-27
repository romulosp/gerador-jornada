package gitworkspace

import (
	"context"
	"errors"
	"gerador-jornada.local/editor/internal/application"
	"os"
	"path/filepath"
	"testing"
)

// TestJournalBoundaries testa recuperação de cada fronteira mutável com novidade remota real.
func TestJournalBoundaries(t *testing.T) {
	// phase identifica intenção, escrita de arquivo, índice, referência e vínculo persistido.
	for _, phase := range []string{"intent", "file:Remota.jornada", "index", "reference", "binding"} {
		t.Run(phase, func(t *testing.T) {
			// engine mantém configuração e remoto separados da raiz de trabalho.
			engine, root, remote, destination := fixture(t)
			put(t, root, "A.jornada", "inicial")
			if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
				t.Fatal(err)
			}
			if err := publish(t, engine, root); err != nil {
				t.Fatal(err)
			}
			remoteCommit(t, remote, "Remota.jornada", "novidade remota")
			put(t, root, "A.jornada", "edição local preservada")
			engine.fail = func(step string) error {
				if step == phase {
					return errors.New("interrupção controlada")
				}
				return nil
			}
			if err := publish(t, engine, root); err == nil {
				t.Fatal("falha não induzida")
			}
			engine.fail = nil
			if err := engine.Recover(context.Background(), root); err != nil {
				t.Fatal(err)
			}
			if err := publish(t, engine, root); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(root, "A.jornada"))
			if err != nil || string(data) != "edição local preservada" {
				t.Fatal("edição local perdeu bytes")
			}
			data, err = os.ReadFile(filepath.Join(root, "Remota.jornada"))
			if err != nil || string(data) != "novidade remota" {
				t.Fatal("novidade remota perdida")
			}
		})
	}
}

// TestRecoveryDoesNotOverwriteExternalEdit preserva a edição feita depois de uma materialização interrompida.
func TestRecoveryDoesNotOverwriteExternalEdit(t *testing.T) {
	// engine cria histórico remoto conhecido antes de induzir uma colisão durante a recuperação.
	engine, root, remote, destination := fixture(t)
	put(t, root, "A.jornada", "inicial")
	if err := engine.Bind(context.Background(), root, destination, application.GitHubCredential{}); err != nil {
		t.Fatal(err)
	}
	if err := publish(t, engine, root); err != nil {
		t.Fatal(err)
	}
	remoteCommit(t, remote, "Remota.jornada", "remoto")
	put(t, root, "A.jornada", "edição")
	engine.fail = func(step string) error {
		if step == "file:Remota.jornada" {
			return errors.New("interrupção")
		}
		return nil
	}
	if err := publish(t, engine, root); err == nil {
		t.Fatal("falha não induzida")
	}
	engine.fail = nil
	put(t, root, "Remota.jornada", "alterado externamente")
	if err := engine.Recover(context.Background(), root); err == nil {
		t.Fatal("edição externa não bloqueou rollback")
	}
	data, _ := os.ReadFile(filepath.Join(root, "Remota.jornada"))
	if string(data) != "alterado externamente" {
		t.Fatal("edição externa perdida")
	}
}

// TestWindowsPathValidation rejeita nomes ambíguos antes de materializar conteúdo remoto.
func TestWindowsPathValidation(t *testing.T) {
	for _, name := range []string{"../outside", "/absolute", ".git/config", "CON.txt", "folder/aux", "bad:name", "trailing. ", "folder\\file"} {
		if err := validPath(name); err == nil {
			t.Fatalf("caminho inseguro aceito: %s", name)
		}
	}
	if err := validateTree(map[string]entry{"A.jornada": {}, "a.jornada": {}}); err == nil {
		t.Fatal("colisão de caixa aceita")
	}
	if err := validateTree(map[string]entry{"folder": {}, "folder/file": {}}); err == nil {
		t.Fatal("colisão arquivo/pasta aceita")
	}
}
