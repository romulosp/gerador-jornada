package projectzip

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// visitSerializedMenu percorre a fixture para simular precisamente o catálogo anterior.
func visitSerializedMenu(node map[string]any, visit func(map[string]any)) {
	// items mantém a ordem e a identidade de cada ramo durante a transformação de teste.
	items, _ := node["items"].([]any)
	for _, raw := range items {
		item := raw.(map[string]any)
		visit(item)
		if child, ok := item["submenu"].(map[string]any); ok {
			visitSerializedMenu(child, visit)
		}
	}
}

// TestJourneyCatalogMigration verifica contrato estrito, Unicode e backup integral antes do save 3/4.
func TestJourneyCatalogMigration(t *testing.T) {
	// document possui uma árvore completa, incluindo folhas e submenus legados.
	document := compoundDocument(t)
	menu := document.Screens[0].Components[2].Properties.Menu
	leaf := &menu.Initial().Items[0]
	leaf.TargetScreenID = ""
	leaf.TargetJourney = "Arrecadação de cobrança.jornada"
	data, err := Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := Decode(data)
	if err != nil || restored.Screens[0].Components[2].Properties.Menu.Initial().Items[0].TargetJourney != leaf.TargetJourney {
		t.Fatalf("destino Unicode não preservado: %v", err)
	}
	manifest, err := EncodeManifest(document)
	if err != nil {
		t.Fatal(err)
	}
	root := manifestObject(t, manifest)
	root["catalogVersion"] = 3
	visitSerializedMenu(propertiesByKind(t, root, "applicationMenu")["root"].(map[string]any), func(item map[string]any) { delete(item, "targetJourney") })
	legacyManifest, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	legacy := archiveWithManifest(t, legacyManifest)
	migrated, err := Decode(legacy)
	if err != nil || migrated.Project.CatalogVersion != 3 || migrated.Screens[0].Components[2].Properties.Menu.Initial().Items[0].TargetJourney != "" {
		t.Fatalf("migração inventou destino ou recusou 3/3: %v", err)
	}
	path := filepath.Join(t.TempDir(), "legado.jornada")
	if err := os.WriteFile(path, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	if err := Save(context.Background(), path, migrated); err != nil {
		t.Fatal(err)
	}
	backup, err := os.ReadFile(path + ".bak")
	if err != nil || !bytes.Equal(backup, legacy) {
		t.Fatalf("original 3/3 não preservado integralmente: %v", err)
	}
	current, err := Load(path)
	if err != nil || current.Project.CatalogVersion != 4 {
		t.Fatalf("save não publicou 3/4: %v", err)
	}
	// Um backup indisponível impede a migração e mantém os bytes da origem.
	blocked := filepath.Join(t.TempDir(), "bloqueada.jornada")
	if err := os.WriteFile(blocked, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(blocked+".bak", 0755); err != nil {
		t.Fatal(err)
	}
	if err := Save(context.Background(), blocked, migrated); err == nil {
		t.Fatal("migração prosseguiu sem conservar original")
	}
	unchanged, _ := os.ReadFile(blocked)
	if !bytes.Equal(unchanged, legacy) {
		t.Fatal("falha de backup alterou original")
	}
	// fixtures inválidas exercitam ausência no catálogo atual e presença no histórico.
	for _, version := range []int{3, 4} {
		invalid := manifestObject(t, manifest)
		invalid["catalogVersion"] = version
		if version == 4 {
			visitSerializedMenu(propertiesByKind(t, invalid, "applicationMenu")["root"].(map[string]any), func(item map[string]any) { delete(item, "targetJourney") })
		}
		patched, _ := json.Marshal(invalid)
		if _, err := Decode(archiveWithManifest(t, patched)); err == nil || !strings.Contains(err.Error(), "PROJECT_SCHEMA_INVALID") {
			t.Fatalf("contrato de catálogo %d ignorado: %v", version, err)
		}
	}
}
