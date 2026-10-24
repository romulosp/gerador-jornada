package projectzip

import (
	"encoding/json"
	"testing"

	"gerador-jornada.local/editor/internal/application"
)

// TestFileIDPersistenceAndMigration cobre ida/volta v2, leitura v1 e recusa de v2 incompleta.
func TestFileIDPersistenceAndMigration(t *testing.T) {
	// document usa IDs explícitos distintos do nome amigável.
	document := application.Example()
	document.Screens[0].FileID = "PnlMinhaIdentificacao"
	document.Screens[0].Components[0].FileID = "TituloPrincipal"
	// data contém o novo manifesto com nomes de arquivo separados do UUID.
	data, err := EncodeManifest(document)
	if err != nil {
		t.Fatal(err)
	}
	// root permite simular a estrutura real de um arquivo anterior, sem fileId.
	root := map[string]any{}
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatal(err)
	}
	if root["schemaVersion"] != float64(3) {
		t.Fatal("nova estrutura deve declarar versão 3")
	}
	// restored comprova preservação de UUID e ID de arquivo com valores explícitos.
	restored, err := DecodeManifest(data)
	if err != nil || restored.Screens[0].ID != document.Screens[0].ID || restored.Screens[0].FileID != "PnlMinhaIdentificacao" || restored.Screens[0].Components[0].FileID != "TituloPrincipal" {
		t.Fatalf("perda de identidade: %v", err)
	}
	for _, item := range root["screens"].([]any) {
		// page representa uma tela serializada sem os atributos recém-introduzidos.
		page := item.(map[string]any)
		delete(page, "fileId")
		for _, element := range page["components"].([]any) {
			delete(element.(map[string]any), "fileId")
		}
	}
	data, _ = json.Marshal(root)
	if _, err := DecodeManifest(data); err == nil {
		t.Fatal("v2 sem ID obrigatório não pode ser aceita")
	}
	root["schemaVersion"] = 1
	root["catalogVersion"] = 1
	data, _ = json.Marshal(root)
	restored, err = DecodeManifest(data)
	if err != nil || restored.Screens[0].FileID != "PnlIdentificacao" || restored.Screens[0].ID != document.Screens[0].ID {
		t.Fatalf("migração v1 incorreta: %v", err)
	}
	if _, err := EncodeManifest(restored); err != nil {
		t.Fatalf("projeto migrado não pode ser salvo: %v", err)
	}
	root["schemaVersion"] = 999
	data, _ = json.Marshal(root)
	if _, err := DecodeManifest(data); err == nil {
		t.Fatal("versão desconhecida foi aceita")
	}
}
