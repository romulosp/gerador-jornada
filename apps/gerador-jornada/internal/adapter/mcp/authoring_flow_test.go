package mcp

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// TestAuthoringFlowTwoScreens preserva identidades, ordem e edição atômica pelos comandos MCP.
func TestAuthoringFlowTwoScreens(t *testing.T) {
	// runtime/client conservam o mesmo núcleo da composição distribuída, sem persistência.
	runtime, client := testClient(t)
	created := callResult[struct {
		SessionID string `json:"sessionId"`
	}](t, client, "jornada_criar", map[string]any{"requestId": project.NewID(), "name": "Duas telas", "mode": "blank"})
	if !project.ValidID(created.SessionID) {
		t.Fatal("criação sem sessão")
	}
	initial, _ := runtime.authoring.Inspect(created.SessionID)
	firstPage := initial.Document.Screens[0].ID
	// second é criada com o mesmo tamanho/defaults disponíveis no editor.
	args := editArgs(t, runtime, created.SessionID)
	args["name"] = "Conclusão"
	second := callResult[struct {
		ScreenID string `json:"screenId"`
	}](t, client, "tela_criar", args)
	args = editArgs(t, runtime, created.SessionID)
	args["screenId"], args["kind"], args["x"], args["y"] = firstPage, "button", 20, 20
	button := callResult[struct {
		ComponentID string `json:"componentId"`
	}](t, client, "componente_inserir", args)
	args = editArgs(t, runtime, created.SessionID)
	args["sourceScreenId"], args["sourceComponentId"], args["targetScreenId"], args["action"] = firstPage, button.ComponentID, second.ScreenID, "navigate"
	callResult[struct{}](t, client, "navegacao_vincular", args)
	before, _ := runtime.authoring.Inspect(created.SessionID)
	if len(before.Document.Screens) != 2 || len(before.Document.Journey.Actions) != 1 {
		t.Fatal("fluxo não publicado no rascunho")
	}
	// invalid aponta para UUID não pertencente ao documento; nenhum destino externo é procurado.
	invalid := editArgs(t, runtime, created.SessionID)
	invalid["sourceScreenId"], invalid["sourceComponentId"], invalid["targetScreenId"], invalid["action"] = firstPage, button.ComponentID, project.NewID(), "navigate"
	callFailure(t, client, "navegacao_vincular", invalid, "invalidReference")
	after, _ := runtime.authoring.Inspect(created.SessionID)
	if before.Revision != after.Revision || !reflect.DeepEqual(before.Document, after.Document) {
		t.Fatal("referência inválida alterou fluxo")
	}
	args = editArgs(t, runtime, created.SessionID)
	args["screenId"] = firstPage
	duplicate := callResult[struct {
		ScreenID string            `json:"screenId"`
		IDMap    map[string]string `json:"idMap"`
	}](t, client, "tela_duplicar", args)
	final, _ := runtime.authoring.Inspect(created.SessionID)
	if len(final.Document.Screens) != 3 || final.Document.Screens[0].ID != firstPage || final.Document.Screens[1].ID != second.ScreenID || duplicate.ScreenID == firstPage {
		t.Fatal("duplicação perdeu ordem/identidades")
	}
	if duplicate.IDMap[firstPage] != duplicate.ScreenID || duplicate.IDMap[button.ComponentID] != final.Document.Screen(duplicate.ScreenID).Components[0].ID {
		t.Fatal("mapa de IDs da cópia ausente ou divergente")
	}
	if err := final.Document.ValidateFileIDs(); err != nil {
		t.Fatal(err)
	}
	callResult[struct{}](t, client, "jornada_validar", map[string]any{"sessionId": created.SessionID, "target": "draft"})
}

// TestHistoricalContentRemainsReadableAndUnsaved não transforma limite de nova entrada em migração destrutiva.
func TestHistoricalContentRemainsReadableAndUnsaved(t *testing.T) {
	// document/source simulam pacote histórico válido com texto acima do limite atual de autoria.
	document := application.NewDocument("Histórico")
	session := application.NewSession(document)
	id, err := session.AddComponent(document.Screens[0].ID, catalog.Title, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Repeat("😀", 65537)
	session.Document.Screens[0].ComponentByID(id).Properties.Title.Text = text
	root := t.TempDir()
	source := filepath.Join(root, "historica.jornada")
	if err := projectzip.Save(context.Background(), source, session.Document); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	runtime, client := testClientAt(t, root)
	opened := callResult[admissionResult](t, client, "jornada_abrir", map[string]any{"requestId": project.NewID(), "filename": "historica.jornada"})
	query := callResult[pageResult[struct {
		Properties struct {
			Text string `json:"text"`
		} `json:"properties"`
	}]](t, client, "jornada_obter", map[string]any{"sessionId": opened.SessionID, "screenId": document.Screens[0].ID})
	if len(query.Items) != 1 || query.Items[0].Properties.Text != text {
		t.Fatal("leitura histórica truncou conteúdo")
	}
	args := editArgs(t, runtime, opened.SessionID)
	args["screenId"], args["componentId"], args["changes"] = document.Screens[0].ID, id, map[string]any{"enabled": false}
	callResult[struct{}](t, client, "componente_atualizar", args)
	after, err := os.ReadFile(source)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("rascunho persistiu sem revisão/consentimento")
	}
	snapshot, _ := runtime.authoring.Inspect(opened.SessionID)
	if snapshot.Document.Screens[0].ComponentByID(id).Properties.Title.Text != text {
		t.Fatal("edição independente apagou texto histórico")
	}
}

// TestScreenConfigurationPreservesGeometry recusa redução incompatível sem publicar campos anteriores ao erro.
func TestScreenConfigurationPreservesGeometry(t *testing.T) {
	// runtime/client aplicam comandos pela borda real de schemas e replay.
	runtime, client := testClient(t)
	document := application.NewDocument("Geometria")
	info, err := runtime.authoring.Admit(project.NewID(), "screen-fixture", document, application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	args := editArgs(t, runtime, info.ID)
	args["screenId"], args["kind"], args["x"], args["y"] = document.Screens[0].ID, "button", 700, 500
	callResult[struct{}](t, client, "componente_inserir", args)
	before, _ := runtime.authoring.Inspect(info.ID)
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["changes"] = document.Screens[0].ID, map[string]any{"name": "Não publicar", "width": 300}
	callFailure(t, client, "tela_configurar", args, "invalidArgument")
	after, _ := runtime.authoring.Inspect(info.ID)
	if !reflect.DeepEqual(before.Document, after.Document) || before.Revision != after.Revision {
		t.Fatal("configuração inválida publicou campos parciais")
	}
}
