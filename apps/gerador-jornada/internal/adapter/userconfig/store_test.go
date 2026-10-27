package userconfig

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestStorePreservesCorruption impede que uma atualização apague o único estado recuperável.
func TestStorePreservesCorruption(t *testing.T) {
	// store e record são isolados de configurações pessoais.
	store := &Store{Root: t.TempDir()}
	// record é um contrato local mínimo para exercitar publicação e proteção contra corrupção.
	record := struct {
		// Version identifica o formato de teste suportado pelo armazenamento.
		Version int `json:"version"`
		// Value permite distinguir o conteúdo inicial da corrupção deliberada.
		Value string `json:"value"`
	}{1, "local"}
	// err comprova se a operação aceita somente o estado e caminho permitidos.
	if err := store.Write(context.Background(), "state.json", record); err != nil {
		t.Fatal(err)
	}
	// path recebe corrupção deliberada para comprovar recusa de substituição.
	path := filepath.Join(store.Root, "state.json")
	// err interrompe o teste se a preparação da corrupção não puder ser gravada.
	if err := os.WriteFile(path, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	// err comprova se a operação aceita somente o estado e caminho permitidos.
	if err := store.Write(context.Background(), "state.json", record); err == nil {
		t.Fatal("corrupção sobrescrita")
	}
	// data conserva os bytes observados depois da tentativa recusada de substituição.
	data, _ := os.ReadFile(path)
	if string(data) != "broken" {
		t.Fatal("original modificado")
	}
	// err comprova se a operação aceita somente o estado e caminho permitidos.
	if err := store.Write(context.Background(), "../escape.json", record); err == nil {
		t.Fatal("escape aceito")
	}
}
