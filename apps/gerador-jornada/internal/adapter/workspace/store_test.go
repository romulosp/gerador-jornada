package workspace

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
)

// TestWorkspaceStoragePreservesManualFlow conserva o comportamento usado pelo launcher.
func TestWorkspaceStoragePreservesManualFlow(t *testing.T) {
	// store seleciona pasta existente; o adaptador não exige repositório Git.
	store, err := OpenWorkspace(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	first, document, err := store.Create("Nome / Unicode ç", false)
	if err != nil {
		t.Fatal(err)
	}
	if document.Project.Name != "Nome / Unicode ç" || len(document.Screens[0].Components) != 0 || len(document.Project.Assets) == 0 {
		t.Fatal("default de nova jornada alterado")
	}
	second, _, err := store.Create("Nome / Unicode ç", true)
	if err != nil || first == second {
		t.Fatalf("homônimo: %s (%v)", second, err)
	}
	// external é origem escolhida manualmente e deve permanecer byte a byte intacta.
	external := filepath.Join(t.TempDir(), "importada.jornada")
	if err := projectzip.Save(context.Background(), external, document); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(external)
	imported, err := store.Import(external)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(store.Root, imported))
	after, _ := os.ReadFile(external)
	if !bytes.Equal(before, got) || !bytes.Equal(before, after) {
		t.Fatal("importação reescreveu bytes")
	}
	if err := os.Mkdir(filepath.Join(store.Root, "sub"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.Root, "broken.jornada"), []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := projectzip.Save(context.Background(), filepath.Join(store.Root, "sub", "nested.jornada"), document); err != nil {
		t.Fatal(err)
	}
	// entries não deve percorrer descendentes nem esconder arquivo corrompido.
	entries, err := store.List()
	if err != nil || len(entries) != 4 {
		t.Fatalf("inventário incorreto: %+v (%v)", entries, err)
	}
	// invalid confere diagnóstico por entrada em vez de descartar silenciosamente o arquivo.
	invalid := false
	for _, entry := range entries {
		if entry.Filename == "broken.jornada" && entry.Error != "" {
			invalid = true
		}
	}
	if !invalid {
		t.Fatal("arquivo inválido omitido")
	}
	if _, err := store.Resolve("../outside.jornada"); err == nil {
		t.Fatal("caminho externo aceito")
	}
}

// TestWorkspacePinsPhysicalRoot recusa a reutilização do caminho por outra pasta.
func TestWorkspacePinsPhysicalRoot(t *testing.T) {
	// parent permite deslocar a raiz sem cruzar volumes nem tocar dados do usuário.
	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	name, _, err := store.Create("Original", false)
	if err != nil {
		t.Fatal(err)
	}
	// snapshot vincula hash e identidade aos bytes efetivamente decodificados.
	document, baseline, err := store.LoadSnapshot(name)
	if err != nil || !baseline.Exists || baseline.Identity == "" || baseline.SHA256 == "" || document.Project.Name != "Original" {
		t.Fatalf("base incompleta: %+v (%v)", baseline, err)
	}
	if err := os.Rename(root, filepath.Join(parent, "original")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := store.List(); err == nil {
		t.Fatal("inventário aceitou raiz substituída")
	}
	if _, _, err := store.LoadSnapshot(name); err == nil {
		t.Fatal("leitura aceitou raiz substituída")
	}
	if _, _, err := store.Create("Indevida", false); err == nil {
		t.Fatal("escrita aceitou raiz substituída")
	}
	// replacement deve continuar vazia após todas as tentativas.
	replacement, err := os.ReadDir(root)
	if err != nil || len(replacement) != 0 {
		t.Fatalf("efeito na pasta substituta: %v", err)
	}
}

// TestWorkspaceInvalidImportAndExclusivePublish preserva destinos e não publica pacote inválido.
func TestWorkspaceInvalidImportAndExclusivePublish(t *testing.T) {
	// store usa pasta local sem requisito Git nem dependência de LLM.
	store, err := OpenWorkspace(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	name, _, err := store.Create("Original", false)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(store.Root, name))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PublishNew(context.Background(), name, []byte("substituição")); err == nil {
		t.Fatal("homônimo sobrescrito")
	}
	// invalid conserva uma origem corrompida fora da raiz de destino.
	invalid := filepath.Join(t.TempDir(), "invalid.jornada")
	if err := os.WriteFile(invalid, []byte("inválido"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Import(invalid); err == nil {
		t.Fatal("pacote inválido importado")
	}
	after, err := os.ReadFile(filepath.Join(store.Root, name))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("destino original alterado")
	}
	entries, err := os.ReadDir(store.Root)
	if err != nil || len(entries) != 1 {
		t.Fatalf("resíduo de publicação: %v", err)
	}
}
