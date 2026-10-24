package projectzip

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// compoundDocument monta uma tela com os três compostos, Unicode, vazio e item desabilitado.
func compoundDocument(t *testing.T) application.Document {
	t.Helper()
	document := application.NewDocument("Projeto 031/032/033 — ação")
	page := &document.Screens[0]
	// input guarda somente o rótulo configurado, sem estado de runtime.
	input, err := screen.NewComponent(catalog.BarcodeInput, "codigoBarras", 10, 10)
	if err != nil {
		t.Fatalf("entrada inválida: %v", err)
	}
	input.Bounds = screen.Bounds{X: 10, Y: 10, Width: 610, Height: 60}
	input.Properties.BarcodeInput.Label = "Número do Código de Barras: ação"
	// processed cobre os cinco rótulos, um valor vazio e um bloco com título vazio.
	processed, err := screen.NewComponent(catalog.BarcodeProcessed, "codigoProcessado", 10, 90)
	if err != nil {
		t.Fatalf("ficha inválida: %v", err)
	}
	processed.Bounds = screen.Bounds{X: 10, Y: 90, Width: 530, Height: 120}
	processed.Properties.BarcodeProcessed = &catalog.BarcodeProcessedProperties{Blocks: []catalog.BarcodeBlock{
		{
			ID: project.NewID(), Title: "Dados — ação",
			Rows: []catalog.BarcodeRow{
				{ID: project.NewID(), LabelKind: catalog.BarcodeLabelCPF, Value: ""},
				{ID: project.NewID(), LabelKind: catalog.BarcodeLabelCNPJ, Value: "12.345.678/0001-90"},
				{ID: project.NewID(), LabelKind: catalog.BarcodeLabelName, Value: "José, ção"},
				{ID: project.NewID(), LabelKind: catalog.BarcodeLabelTradeName, Value: "Nome fantasia — ação"},
				{ID: project.NewID(), LabelKind: catalog.BarcodeLabelCorporateName, Value: "Razão social"},
			},
		},
		{
			ID: project.NewID(), Title: "",
			Rows: []catalog.BarcodeRow{{ID: project.NewID(), LabelKind: catalog.BarcodeLabelName, Value: ""}},
		},
	}}
	// menu usa o seed e recebe menu inicial em descendente, item desabilitado e items vazio.
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menuAplicacao", 0, 0)
	if err != nil {
		t.Fatalf("menu inválido: %v", err)
	}
	menu.Bounds = screen.Bounds{X: 0, Y: 0, Width: 800, Height: 600}
	// paga é o submenu escolhido como inicial para provar o remapeamento persistido.
	paga := menuNodeByLabel(&menu.Properties.Menu.Root, "PAGA PREMIADOS")
	if paga == nil {
		t.Fatal("PAGA PREMIADOS ausente do seed")
	}
	paga.Items[0].Enabled = false
	paga.Items[0].Legend = "Legenda desabilitada — ação"
	paga.Items[0].Category = ""
	paga.Items = append(paga.Items, catalog.MenuItem{
		ID: project.NewID(), Label: "SUBMENU VAZIO", Enabled: false,
		Submenu: &catalog.SubmenuProperties{ID: project.NewID(), Name: "SUBMENU VAZIO", Layout: catalog.MenuLayoutStandard, Items: []catalog.MenuItem{}},
	})
	menu.Properties.Menu.InitialMenuID = paga.ID
	page.Components = append(page.Components, input, processed, menu)
	return document
}

// menuNodeByLabel procura um nó da árvore pelo rótulo do item que o contém.
func menuNodeByLabel(node *catalog.SubmenuProperties, label string) *catalog.SubmenuProperties {
	if node == nil {
		return nil
	}
	for index := range node.Items {
		if node.Items[index].Label == label && node.Items[index].Submenu != nil {
			return node.Items[index].Submenu
		}
		if found := menuNodeByLabel(node.Items[index].Submenu, label); found != nil {
			return found
		}
	}
	return nil
}

// menuNodeByName procura um nó da árvore pelo nome editável.
func menuNodeByName(node *catalog.SubmenuProperties, name string) *catalog.SubmenuProperties {
	if node == nil {
		return nil
	}
	if node.Name == name {
		return node
	}
	for index := range node.Items {
		if found := menuNodeByName(node.Items[index].Submenu, name); found != nil {
			return found
		}
	}
	return nil
}

// assertProcessedEqual compara blocos e linhas preservando ordem e textos.
func assertProcessedEqual(t *testing.T, want, got catalog.BarcodeProcessedProperties) {
	t.Helper()
	if len(want.Blocks) != len(got.Blocks) {
		t.Fatalf("quantidade de blocos divergente: %d != %d", len(want.Blocks), len(got.Blocks))
	}
	for blockIndex := range want.Blocks {
		wantBlock, gotBlock := want.Blocks[blockIndex], got.Blocks[blockIndex]
		if wantBlock.ID != gotBlock.ID || wantBlock.Title != gotBlock.Title {
			t.Fatalf("bloco %d divergente: %+v != %+v", blockIndex+1, gotBlock, wantBlock)
		}
		if len(wantBlock.Rows) != len(gotBlock.Rows) {
			t.Fatalf("linhas do bloco %d divergentes: %d != %d", blockIndex+1, len(gotBlock.Rows), len(wantBlock.Rows))
		}
		for rowIndex := range wantBlock.Rows {
			if wantBlock.Rows[rowIndex] != gotBlock.Rows[rowIndex] {
				t.Fatalf("linha %d do bloco %d divergente: %+v != %+v", rowIndex+1, blockIndex+1, gotBlock.Rows[rowIndex], wantBlock.Rows[rowIndex])
			}
		}
	}
}

// assertMenuEqual compara a árvore inteira, inclusive nils, ordem e campos textuais.
func assertMenuEqual(t *testing.T, want, got *catalog.SubmenuProperties) {
	t.Helper()
	if want == nil || got == nil {
		if want != got {
			t.Fatalf("presença de submenu divergente: %v != %v", want != nil, got != nil)
		}
		return
	}
	if want.ID != got.ID || want.Name != got.Name || want.Title != got.Title || want.Layout != got.Layout {
		t.Fatalf("nó divergente: %+v != %+v", got, want)
	}
	if len(want.Items) != len(got.Items) {
		t.Fatalf("itens do nó %q divergentes: %d != %d", want.Name, len(got.Items), len(want.Items))
	}
	for index := range want.Items {
		wantItem, gotItem := want.Items[index], got.Items[index]
		if wantItem.ID != gotItem.ID || wantItem.Label != gotItem.Label || wantItem.Legend != gotItem.Legend ||
			wantItem.Category != gotItem.Category || wantItem.Variant != gotItem.Variant ||
			wantItem.Enabled != gotItem.Enabled || wantItem.TargetScreenID != gotItem.TargetScreenID || wantItem.TargetJourney != gotItem.TargetJourney ||
			wantItem.MovementCode != gotItem.MovementCode || wantItem.MovementOption != gotItem.MovementOption ||
			wantItem.MovementName != gotItem.MovementName || wantItem.Shortcut != gotItem.Shortcut {
			t.Fatalf("item %d do nó %q divergente: %+v != %+v", index+1, want.Name, gotItem, wantItem)
		}
		assertMenuEqual(t, wantItem.Submenu, gotItem.Submenu)
	}
}

// manifestObject decodifica o manifesto em uma estrutura mutável independente.
func manifestObject(t *testing.T, manifest []byte) map[string]any {
	t.Helper()
	root := map[string]any{}
	if err := json.Unmarshal(manifest, &root); err != nil {
		t.Fatalf("manifesto inválido: %v", err)
	}
	return root
}

// propertiesByKind devolve o objeto de propriedades do componente indicado.
func propertiesByKind(t *testing.T, root map[string]any, kind string) map[string]any {
	t.Helper()
	screens, ok := root["screens"].([]any)
	if !ok || len(screens) == 0 {
		t.Fatalf("manifesto sem telas: %v", root["screens"])
	}
	components, ok := screens[0].(map[string]any)["components"].([]any)
	if !ok {
		t.Fatalf("tela sem componentes")
	}
	for _, raw := range components {
		component, ok := raw.(map[string]any)
		if !ok || component["kind"] != kind {
			continue
		}
		properties, ok := component["properties"].(map[string]any)
		if !ok {
			t.Fatalf("propriedades de %s não são objeto", kind)
		}
		return properties
	}
	t.Fatalf("componente %s ausente do manifesto", kind)
	return nil
}

// archiveWithManifest monta um ZIP mínimo apenas com o manifesto informado.
func archiveWithManifest(t *testing.T, manifest []byte) []byte {
	t.Helper()
	// buffer recebe o pacote em memória, sem tocar o disco.
	buffer := &bytes.Buffer{}
	writer := zip.NewWriter(buffer)
	entry, err := writer.CreateHeader(&zip.FileHeader{Name: "manifest.json", Method: zip.Deflate})
	if err != nil {
		t.Fatalf("entrada do manifesto falhou: %v", err)
	}
	if _, err := entry.Write(manifest); err != nil {
		t.Fatalf("escrita do manifesto falhou: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("ZIP mínimo inválido: %v", err)
	}
	return buffer.Bytes()
}

// zipEntry lê uma parte do pacote sem extrair arquivos.
func zipEntry(t *testing.T, data []byte, name string) []byte {
	t.Helper()
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("ZIP não abre: %v", err)
	}
	file, err := archive.Open(name)
	if err != nil {
		t.Fatalf("parte %s ausente: %v", name, err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("leitura de %s falhou: %v", name, err)
	}
	return content
}

// TestCompoundManifestShape confere versões 3/4 e o corpo exato de cada composto.
func TestCompoundManifestShape(t *testing.T) {
	document := compoundDocument(t)
	data, err := Encode(document)
	if err != nil {
		t.Fatalf("encode falhou: %v", err)
	}
	// manifest contém a revisão e as propriedades planas dos compostos.
	manifest := zipEntry(t, data, "manifest.json")
	root := manifestObject(t, manifest)
	if root["schemaVersion"] != float64(3) || root["catalogVersion"] != float64(4) {
		t.Fatalf("writer deve publicar 3/4: %v/%v", root["schemaVersion"], root["catalogVersion"])
	}
	// inputSerialized preserva somente o rótulo, sem invólucro adicional.
	inputSerialized := propertiesByKind(t, root, "barcodeInput")
	if len(inputSerialized) != 1 || inputSerialized["label"] != document.Screens[0].Components[0].Properties.BarcodeInput.Label {
		t.Fatalf("propriedades da entrada divergentes: %v", inputSerialized)
	}
	// processedBlocks preserva blocos e linhas na ordem confirmada.
	processedSerialized := propertiesByKind(t, root, "barcodeProcessed")
	blocks, ok := processedSerialized["blocks"].([]any)
	if !ok || len(blocks) != 2 {
		t.Fatalf("blocos serializados divergentes: %v", processedSerialized["blocks"])
	}
	firstBlock, ok := blocks[0].(map[string]any)
	if !ok {
		t.Fatalf("bloco serializado inválido: %v", blocks[0])
	}
	if firstBlock["title"] != "Dados — ação" {
		t.Fatalf("título do bloco não preservado: %v", firstBlock["title"])
	}
	rows, ok := firstBlock["rows"].([]any)
	if !ok || len(rows) != 5 {
		t.Fatalf("linhas serializadas divergentes: %v", firstBlock["rows"])
	}
	firstRow, ok := rows[0].(map[string]any)
	if !ok || firstRow["labelKind"] != catalog.BarcodeLabelCPF || firstRow["value"] != "" {
		t.Fatalf("linha vazia não preservada: %v", rows[0])
	}
	// menuSerialized contém root e initialMenuId diretamente, sem segundo invólucro.
	menuSerialized := propertiesByKind(t, root, "applicationMenu")
	if _, ok := menuSerialized["menu"]; ok {
		t.Fatal("properties não pode repetir o invólucro menu")
	}
	if menuSerialized["initialMenuId"] != document.Screens[0].Components[2].Properties.Menu.InitialMenuID {
		t.Fatalf("menu inicial divergente: %v", menuSerialized["initialMenuId"])
	}
	rootSerialized, ok := menuSerialized["root"].(map[string]any)
	if !ok {
		t.Fatalf("root serializado inválido: %v", menuSerialized["root"])
	}
	items, ok := rootSerialized["items"].([]any)
	if !ok || len(items) != len(document.Screens[0].Components[2].Properties.Menu.Root.Items) {
		t.Fatalf("itens da raiz divergentes: %v", rootSerialized["items"])
	}
	firstItem, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("primeiro item inválido: %v", items[0])
	}
	// todos os campos obrigatórios aparecem, inclusive enabled falso e strings vazias.
	for _, field := range []string{"id", "label", "legend", "category", "variant", "enabled", "submenu", "targetScreenId", "targetJourney", "movementCode", "movementOption", "movementName", "shortcut"} {
		if _, ok := firstItem[field]; !ok {
			t.Fatalf("campo %s ausente no item serializado", field)
		}
	}
	// o booleano false precisa ser emitido explicitamente, não omitido.
	if !strings.Contains(string(manifest), `"enabled": false`) {
		t.Fatal("writer omitiu enabled falso")
	}
	// items vazio continua array, nunca null.
	emptySerialized := findSerializedSubmenu(rootSerialized, "SUBMENU VAZIO")
	if emptySerialized == nil {
		t.Fatal("submenu vazio não foi serializado")
	}
	if items, ok := emptySerialized["items"].([]any); !ok || len(items) != 0 {
		t.Fatalf("items vazio não foi emitido como array: %v", emptySerialized["items"])
	}
}

// findSerializedSubmenu procura um nó serializado pelo nome editável.
func findSerializedSubmenu(node map[string]any, name string) map[string]any {
	if node == nil {
		return nil
	}
	if node["name"] == name {
		return node
	}
	items, _ := node["items"].([]any)
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		submenu, ok := item["submenu"].(map[string]any)
		if !ok {
			continue
		}
		if found := findSerializedSubmenu(submenu, name); found != nil {
			return found
		}
	}
	return nil
}

// TestCompoundRoundTrip confere IDs, ordem, Unicode, vazios e enabled falso após Decode.
func TestCompoundRoundTrip(t *testing.T) {
	document := compoundDocument(t)
	data, err := Encode(document)
	if err != nil {
		t.Fatalf("encode falhou: %v", err)
	}
	restored, err := Decode(data)
	if err != nil {
		t.Fatalf("decode falhou: %v", err)
	}
	if len(restored.Screens) != 1 || len(restored.Screens[0].Components) != 3 {
		t.Fatalf("estrutura restaurada divergente: %d telas", len(restored.Screens))
	}
	// kinds fixa a ordem dos três compostos na tela.
	kinds := []catalog.Kind{catalog.BarcodeInput, catalog.BarcodeProcessed, catalog.ApplicationMenu}
	for index, kind := range kinds {
		component := restored.Screens[0].Components[index]
		if component.Kind != kind || component.ID != document.Screens[0].Components[index].ID {
			t.Fatalf("componente %d divergente: %+v", index+1, component)
		}
	}
	// entrada preserva rótulo Unicode e identidade.
	sourceInput := document.Screens[0].Components[0]
	restoredInput := restored.Screens[0].Components[0]
	if restoredInput.Properties.BarcodeInput.Label != sourceInput.Properties.BarcodeInput.Label {
		t.Fatalf("rótulo restaurado divergente: %q", restoredInput.Properties.BarcodeInput.Label)
	}
	// ficha preserva blocos, linhas, ordem e valores vazios.
	sourceProcessed := document.Screens[0].Components[1]
	restoredProcessed := restored.Screens[0].Components[1]
	assertProcessedEqual(t, *sourceProcessed.Properties.BarcodeProcessed, *restoredProcessed.Properties.BarcodeProcessed)
	// menu preserva a árvore completa e a escolha inicial em descendente.
	sourceMenu := document.Screens[0].Components[2].Properties.Menu
	restoredMenu := restored.Screens[0].Components[2].Properties.Menu
	if restoredMenu.InitialMenuID != sourceMenu.InitialMenuID {
		t.Fatalf("menu inicial divergente: %q", restoredMenu.InitialMenuID)
	}
	assertMenuEqual(t, &sourceMenu.Root, &restoredMenu.Root)
	initial := restoredMenu.Initial()
	if initial == nil || initial.Name != "PAGA PREMIADOS" {
		t.Fatalf("inicial não resolve o descendente: %+v", initial)
	}
	if initial.Items[0].Enabled {
		t.Fatal("enabled falso não sobreviveu ao round-trip")
	}
	if initial.Items[0].Category != "" || initial.Items[0].Legend != "Legenda desabilitada — ação" {
		t.Fatalf("strings vazias ou Unicode não preservados: %+v", initial.Items[0])
	}
	// o submenu desabilitado mantém Items como array vazio não nulo.
	empty := menuNodeByName(&restoredMenu.Root, "SUBMENU VAZIO")
	if empty == nil || empty.Items == nil || len(empty.Items) != 0 {
		t.Fatalf("items vazio virou null ou desapareceu: %+v", empty)
	}
	// a nova escrita é idêntica, comprovando determinismo do pacote.
	again, err := Encode(restored)
	if err != nil {
		t.Fatalf("reencode falhou: %v", err)
	}
	if !bytes.Equal(data, again) {
		t.Fatal("round-trip não determinístico")
	}
}

// TestCompoundSaveLoad confere gravação atômica sem cópia de segurança em catálogo 3.
func TestCompoundSaveLoad(t *testing.T) {
	document := compoundDocument(t)
	// target fica no diretório temporário exclusivo do teste.
	target := filepath.Join(t.TempDir(), "projeto.jornada")
	if err := Save(context.Background(), target, document); err != nil {
		t.Fatalf("save falhou: %v", err)
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatal("catálogo 3 não deve gerar cópia de segurança")
	}
	loaded, err := Load(target)
	if err != nil {
		t.Fatalf("load falhou: %v", err)
	}
	if len(loaded.Screens[0].Components) != 3 {
		t.Fatalf("componentes perdidos no arquivo: %d", len(loaded.Screens[0].Components))
	}
	// a ficha restaurada do disco mantém o título e o valor Unicode.
	processed := loaded.Screens[0].Components[1].Properties.BarcodeProcessed
	if processed.Blocks[0].Title != "Dados — ação" || processed.Blocks[0].Rows[2].Value != "José, ção" {
		t.Fatalf("conteúdo Unicode perdido no disco: %+v", processed.Blocks[0])
	}
	if loaded.Screens[0].Components[2].Properties.Menu.Initial().Name != "PAGA PREMIADOS" {
		t.Fatal("menu inicial perdido no disco")
	}
}

// TestCompoundDecodeRejections recusa versões, campos desconhecidos e presenças obrigatórias.
func TestCompoundDecodeRejections(t *testing.T) {
	document := compoundDocument(t)
	// cases alteram uma cópia independente do manifesto antes de empacotar.
	cases := []struct {
		name   string
		code   string
		mutate func(t *testing.T, root map[string]any)
	}{
		{"schema 4 com catálogo 3", "PROJECT_VERSION_UNSUPPORTED", func(t *testing.T, root map[string]any) {
			root["schemaVersion"] = 4
		}},
		{"schema 3 com catálogo 5", "PROJECT_VERSION_UNSUPPORTED", func(t *testing.T, root map[string]any) {
			root["catalogVersion"] = 5
		}},
		{"campo desconhecido na raiz", "PROJECT_SCHEMA_INVALID", func(t *testing.T, root map[string]any) {
			root["surpresa"] = true
		}},
		{"campo desconhecido nas propriedades", "PROJECT_SCHEMA_INVALID", func(t *testing.T, root map[string]any) {
			propertiesByKind(t, root, "barcodeInput")["extra"] = "x"
		}},
		{"enabled ausente", "PROJECT_SCHEMA_INVALID", func(t *testing.T, root map[string]any) {
			items := propertiesByKind(t, root, "applicationMenu")["root"].(map[string]any)["items"].([]any)
			delete(items[0].(map[string]any), "enabled")
		}},
		{"items null", "PROJECT_SCHEMA_INVALID", func(t *testing.T, root map[string]any) {
			propertiesByKind(t, root, "applicationMenu")["root"].(map[string]any)["items"] = nil
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			manifest, err := EncodeManifest(document)
			if err != nil {
				t.Fatalf("manifesto base inválido: %v", err)
			}
			root := manifestObject(t, manifest)
			testCase.mutate(t, root)
			patched, err := json.Marshal(root)
			if err != nil {
				t.Fatalf("patch do manifesto falhou: %v", err)
			}
			if _, err := Decode(archiveWithManifest(t, patched)); err == nil {
				t.Fatal("manifesto inválido foi aceito")
			} else if !strings.Contains(err.Error(), testCase.code) {
				t.Fatalf("código divergente: %v", err)
			}
		})
	}
}

// TestDecodeLegacyVersions aceita 2/1 e 1/1 e completa os fileIds ausentes.
func TestDecodeLegacyVersions(t *testing.T) {
	document := application.NewDocument("Projeto legado")
	// label representa um componente comum de arquivos anteriores aos compostos.
	label, err := screen.NewComponent(catalog.Label, "rotulo", 10, 10)
	if err != nil {
		t.Fatalf("rótulo inválido: %v", err)
	}
	document.Screens[0].Components = append(document.Screens[0].Components, label)
	// base contém o manifesto 3/3 usado como origem das fixtures.
	base, err := EncodeManifest(document)
	if err != nil {
		t.Fatalf("manifesto base inválido: %v", err)
	}
	t.Run("2/1 preserva fileIds", func(t *testing.T) {
		root := manifestObject(t, base)
		root["schemaVersion"] = 2
		root["catalogVersion"] = 1
		patched, err := json.Marshal(root)
		if err != nil {
			t.Fatalf("patch falhou: %v", err)
		}
		restored, err := Decode(archiveWithManifest(t, patched))
		if err != nil {
			t.Fatalf("2/1 recusado: %v", err)
		}
		if restored.Project.CatalogVersion != 1 {
			t.Fatalf("catálogo restaurado divergente: %d", restored.Project.CatalogVersion)
		}
		if restored.Screens[0].FileID != document.Screens[0].FileID || restored.Screens[0].Components[0].FileID != label.FileID {
			t.Fatal("2/1 não preservou os fileIds explícitos")
		}
	})
	t.Run("1/1 preenche fileIds ausentes", func(t *testing.T) {
		root := manifestObject(t, base)
		root["schemaVersion"] = 1
		root["catalogVersion"] = 1
		// um arquivo v1 real não possui fileId em telas nem em componentes.
		for _, rawPage := range root["screens"].([]any) {
			page := rawPage.(map[string]any)
			delete(page, "fileId")
			for _, rawComponent := range page["components"].([]any) {
				delete(rawComponent.(map[string]any), "fileId")
			}
		}
		patched, err := json.Marshal(root)
		if err != nil {
			t.Fatalf("patch falhou: %v", err)
		}
		restored, err := Decode(archiveWithManifest(t, patched))
		if err != nil {
			t.Fatalf("1/1 recusado: %v", err)
		}
		if restored.Screens[0].FileID != "PnlTela1" || restored.Screens[0].Components[0].FileID != "Rotulo" {
			t.Fatalf("sugestões de fileId divergentes: %q/%q", restored.Screens[0].FileID, restored.Screens[0].Components[0].FileID)
		}
		if restored.Screens[0].ID != document.Screens[0].ID {
			t.Fatal("migração 1/1 alterou o UUID da tela")
		}
	})
}
