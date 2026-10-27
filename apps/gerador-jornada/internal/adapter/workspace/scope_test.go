package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestMaterialScopesSeparatePurposes não transforma um ID ou caminho em autorização nova.
func TestMaterialScopesSeparatePurposes(t *testing.T) {
	// root é concedida apenas como evidência, embora seja pasta acessível ao processo.
	root := t.TempDir()
	scopes, err := NewMaterialScopes([]string{root}, nil)
	if err != nil {
		t.Fatal(err)
	}
	roots := scopes.Roots(LegacyEvidence)
	if len(roots) != 1 || len(scopes.Roots(LibraryImport)) != 0 {
		t.Fatal("finalidades confundidas")
	}
	if _, err := scopes.Open(LibraryImport, roots[0].ID); err == nil {
		t.Fatal("ID atravessou finalidade")
	}
	if _, err := scopes.Open(LegacyEvidence, root); err == nil {
		t.Fatal("caminho virou autorização")
	}
	// lease de evidência é incapaz de publicar mesmo com nome válido.
	lease, err := scopes.Open(LegacyEvidence, roots[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := lease.CreateNew(context.Background(), "indevido.jornada", []byte("bytes")); err == nil {
		t.Fatal("leitura concedeu escrita")
	}
	lease.Close()
	roots[0].ID = "alterado"
	if scopes.Roots(LegacyEvidence)[0].ID == "alterado" {
		t.Fatal("configuração compartilhada com cliente")
	}
}

// TestMaterialScopesRejectReboundRoot conserva a autorização sobre a pasta física original.
func TestMaterialScopesRejectReboundRoot(t *testing.T) {
	// parent mantém todos os movimentos da fixture no mesmo diretório temporário.
	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	scopes, err := NewMaterialScopes(nil, []string{root})
	if err != nil {
		t.Fatal(err)
	}
	rootID := scopes.Roots(LibraryImport)[0].ID
	if err := os.Rename(root, filepath.Join(parent, "original")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := scopes.Open(LibraryImport, rootID); err == nil {
		t.Fatal("autorização seguiu substituição do caminho")
	}
}
