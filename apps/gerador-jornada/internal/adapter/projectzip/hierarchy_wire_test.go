package projectzip

import (
	"encoding/json"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestHierarchyPersistence preserva ParentID e Enabled e aplica defaults em projetos antigos.
func TestHierarchyPersistence(t *testing.T) {
	// document contém um pai e um filho desabilitado.
	document := application.NewDocument("hierarquia")
	parent, err := screen.NewComponent(catalog.Button, "pai", 20, 20)
	if err != nil {
		t.Fatal(err)
	}
	child, err := screen.NewComponent(catalog.Label, "filho", 40, 40)
	if err != nil {
		t.Fatal(err)
	}
	child.ParentID = parent.ID
	child.Enabled = false
	document.Screens[0].Components = append(document.Screens[0].Components, parent, child)
	// data contém o manifesto 3/3 com os novos campos.
	data, err := EncodeManifest(document)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := DecodeManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Screens[0].Components[1].ParentID != parent.ID || restored.Screens[0].Components[1].Enabled {
		t.Fatalf("hierarquia não preservada: %+v", restored.Screens[0].Components)
	}
	// old remove os campos para simular um projeto anterior.
	root := map[string]any{}
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatal(err)
	}
	for _, page := range root["screens"].([]any) {
		for _, element := range page.(map[string]any)["components"].([]any) {
			delete(element.(map[string]any), "enabled")
			delete(element.(map[string]any), "parentId")
		}
	}
	old, _ := json.Marshal(root)
	migrated, err := DecodeManifest(old)
	if err != nil {
		t.Fatal(err)
	}
	for _, component := range migrated.Screens[0].Components {
		if !component.Enabled || component.ParentID != "" {
			t.Fatalf("defaults antigos incorretos: %+v", component)
		}
	}
}
