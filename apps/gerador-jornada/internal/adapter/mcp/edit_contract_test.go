package mcp

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// editArgs constrói argumentos a partir da revisão observada, nunca de contador do cliente.
func editArgs(t *testing.T, runtime *runtime, sessionID string) map[string]any {
	t.Helper()
	// snapshot fornece CAS atual e preserva conteúdo original para comparação posterior.
	snapshot, err := runtime.authoring.Inspect(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"sessionId": sessionID, "requestId": project.NewID(), "expectedRevision": strconv.FormatUint(snapshot.Revision, 10)}
}

// TestMCPTextLimitAndMetadataQuery preserva Unicode integral e recusa excesso antes de mudar a revisão.
func TestMCPTextLimitAndMetadataQuery(t *testing.T) {
	// runtime/client usam transporte SDK, incluindo schema fechado do wrapper opcional properties.
	runtime, client := testClient(t)
	document := application.NewDocument("Unicode")
	info, err := runtime.authoring.Admit(project.NewID(), "unicode-fixture", document, application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	pageID := document.Screens[0].ID
	args := editArgs(t, runtime, info.ID)
	args["screenId"], args["kind"], args["x"], args["y"] = pageID, "button", 0, 0
	inserted := callResult[struct {
		ComponentID string `json:"componentId"`
	}](t, client, "componente_inserir", args)
	text := strings.Repeat("😀", 65536)
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["changes"] = pageID, inserted.ComponentID, map[string]any{"properties": map[string]any{"kind": "button", "text": text}}
	callResult[struct{}](t, client, "componente_atualizar", args)
	// query deve devolver o valor inteiro, sem truncamento, duplicação textual ou pacote base64.
	query := callResult[pageResult[struct {
		Properties struct {
			Text string `json:"text"`
		} `json:"properties"`
	}]](t, client, "jornada_obter", map[string]any{"sessionId": info.ID, "screenId": pageID})
	if len(query.Items) != 1 || query.Items[0].Properties.Text != text {
		t.Fatal("consulta alterou Unicode")
	}
	before, _ := runtime.authoring.Inspect(info.ID)
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["changes"] = pageID, inserted.ComponentID, map[string]any{"properties": map[string]any{"kind": "button", "text": text + "x"}}
	callFailure(t, client, "componente_atualizar", args, "invalidArgument")
	after, _ := runtime.authoring.Inspect(info.ID)
	if before.Revision != after.Revision || !reflect.DeepEqual(before.Document, after.Document) {
		t.Fatal("excesso alterou estado")
	}
}

// callFailure confere erro estruturado sem depender de texto livre ou erro interno da SDK.
func callFailure(t *testing.T, session *sdk.ClientSession, name string, args any, code string) {
	t.Helper()
	// result precisa existir no canal de ferramenta, com estado preservado.
	result, err := session.CallTool(context.Background(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil || !result.IsError {
		t.Fatalf("%s deveria recusar: %+v (%v)", name, result, err)
	}
	encoded, _ := json.Marshal(result.StructuredContent)
	var failure errorEnvelope
	if json.Unmarshal(encoded, &failure) != nil || failure.Error.Code != code || !failure.Error.StatePreserved {
		t.Fatalf("falha divergente: %s", encoded)
	}
}

// TestEditContractsForInstalledKinds cobre dez inserções, sete patches e alterações comuns.
func TestEditContractsForInstalledKinds(t *testing.T) {
	// document vazio permite observar defaults de todas as inserções pelo protocolo real.
	runtime, client := testClient(t)
	document := application.NewDocument("Contratos I1")
	info, err := runtime.authoring.Admit(project.NewID(), "fixture", document, application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	pageID := document.Screens[0].ID
	components := map[catalog.Kind]string{}
	for _, definition := range append(catalog.Definitions(), catalog.LegacyDefinitions()...) {
		// args contém somente campos da inserção; nunca conteúdo genérico de outro Kind.
		args := editArgs(t, runtime, info.ID)
		args["screenId"], args["kind"], args["x"], args["y"] = pageID, string(definition.Kind), 0, 0
		inserted := callResult[struct {
			ComponentID string `json:"componentId"`
		}](t, client, "componente_inserir", args)
		if !project.ValidID(inserted.ComponentID) {
			t.Fatal("inserção sem UUID novo")
		}
		components[definition.Kind] = inserted.ComponentID
	}
	// snapshot referencia opções reais para alterar seleção sem fabricar identidade de outra entidade.
	snapshot, _ := runtime.authoring.Inspect(info.ID)
	if err := snapshot.Document.Screen(pageID).Validate(); err != nil {
		t.Fatalf("inserção produziu composição inválida: %v", err)
	}
	radio := snapshot.Document.Screen(pageID).ComponentByID(components[catalog.Radio]).Properties.Radio
	patches := []struct {
		kind  catalog.Kind
		patch map[string]any
	}{
		{catalog.Button, map[string]any{"kind": "button", "text": "Olá 😀", "variant": "cancel"}},
		{catalog.Label, map[string]any{"kind": "label", "text": "Rótulo novo"}},
		{catalog.Title, map[string]any{"kind": "title", "text": "Título novo"}},
		{catalog.TextField, map[string]any{"kind": "textField", "initialValue": "123", "format": "number"}},
		{catalog.Radio, map[string]any{"kind": "radio", "orientation": "horizontal", "selectedOptionId": radio.Options[1].ID}},
		{catalog.Report, map[string]any{"kind": "report", "title": "Resumo", "rows": []map[string]any{{"id": project.NewID(), "text": "Linha"}}}},
		{catalog.NumericKeyboard, map[string]any{"kind": "numericKeyboard", "targetComponentId": components[catalog.TextField]}},
	}
	for _, patch := range patches {
		t.Logf("patch de família: %s", patch.kind)
		// args exerce o wrapper fechado e a conversão para mudança tipada de aplicação.
		args := editArgs(t, runtime, info.ID)
		args["screenId"], args["componentId"], args["changes"] = pageID, components[patch.kind], map[string]any{"properties": patch.patch}

		callResult[struct{}](t, client, "componente_atualizar", args)
	}
	// advanced é composto inserível, mas mantém suas propriedades avançadas fora de I1.
	advanced := editArgs(t, runtime, info.ID)
	advanced["screenId"], advanced["componentId"], advanced["changes"] = pageID, components[catalog.BarcodeInput], map[string]any{"enabled": false, "name": "Boleto"}
	callResult[struct{}](t, client, "componente_atualizar", advanced)
	// common verifica alterações geométricas/visuais tipadas também numa família avançada não editável em propriedades.
	common := editArgs(t, runtime, info.ID)
	common["screenId"], common["componentId"], common["changes"] = pageID, components[catalog.BarcodeInput], map[string]any{"bounds": map[string]any{"x": 10, "y": 10, "width": 600, "height": 60}, "style": map[string]any{"fontFamily": "Arial", "textColor": "#123456"}}
	callResult[struct{}](t, client, "componente_atualizar", common)
	// invalid não pode alterar documento, revisão nem histórico.
	before, _ := runtime.authoring.Inspect(info.ID)
	invalid := editArgs(t, runtime, info.ID)
	invalid["screenId"], invalid["componentId"], invalid["changes"] = pageID, components[catalog.Button], map[string]any{"properties": map[string]any{"kind": "button", "initialValue": "indevido"}}
	callFailure(t, client, "componente_atualizar", invalid, "invalidArgument")
	after, _ := runtime.authoring.Inspect(info.ID)
	if !reflect.DeepEqual(before.Document, after.Document) || before.Revision != after.Revision {
		t.Fatal("falha de contrato alterou o estado")
	}
	invalid = editArgs(t, runtime, info.ID)
	invalid["screenId"], invalid["componentId"], invalid["changes"] = pageID, components[catalog.BarcodeInput], map[string]any{"properties": map[string]any{"kind": "barcodeInput"}}
	callFailure(t, client, "componente_atualizar", invalid, "capabilityUnavailable")
}

// TestMoveAndRemovalPolicies não remove referência externa nem repete movimento depois de outro comando.
func TestMoveAndRemovalPolicies(t *testing.T) {
	// model cria referências reais de teclado e navegação antes da admissão.
	model := application.NewSession(application.NewDocument("Referências"))
	pageID := model.Document.Screens[0].ID
	field, _ := model.AddComponent(pageID, catalog.TextField, 0, 0)
	keyboard, _ := model.AddComponent(pageID, catalog.NumericKeyboard, 0, 0)
	button, _ := model.AddComponent(pageID, catalog.Button, 0, 0)
	model.Document.Screen(pageID).ComponentByID(keyboard).Properties.Keyboard.TargetComponentID = field
	if err := application.BindNavigation(&model.Document, journey.Action{SourceScreenID: pageID, SourceComponentID: button, Kind: "navigate", TargetScreenID: pageID}); err != nil {
		t.Fatal(err)
	}
	runtime, client := testClient(t)
	info, err := runtime.authoring.Admit(project.NewID(), "reference-fixture", model.Document, application.AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	args := editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["referencePolicy"] = pageID, field, "removeOwnedReferences"
	callFailure(t, client, "componente_remover", args, "invalidReference")
	// move é repetido depois de nova edição; replay deve devolver geometria original e não refazer movimento.
	move := editArgs(t, runtime, info.ID)
	move["screenId"], move["componentId"], move["x"], move["y"] = pageID, keyboard, 100, 100
	first := callResult[struct {
		Bounds boundsDTO `json:"bounds"`
	}](t, client, "componente_mover", move)
	if first.Bounds.X != 100 || first.Bounds.Y != 100 {
		t.Fatal("posição confirmada ausente")
	}
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["changes"] = pageID, keyboard, map[string]any{"properties": map[string]any{"kind": "numericKeyboard", "targetComponentId": ""}}
	callResult[struct{}](t, client, "componente_atualizar", args)
	replayed := callResult[struct {
		Bounds boundsDTO `json:"bounds"`
	}](t, client, "componente_mover", move)
	if replayed != first {
		t.Fatal("replay mudou resposta original")
	}
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["referencePolicy"] = pageID, field, "reject"
	callResult[struct{}](t, client, "componente_remover", args)
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["referencePolicy"] = pageID, button, "reject"
	callFailure(t, client, "componente_remover", args, "invalidReference")
	args = editArgs(t, runtime, info.ID)
	args["screenId"], args["componentId"], args["referencePolicy"] = pageID, button, "removeOwnedReferences"
	callResult[struct{}](t, client, "componente_remover", args)
	after, _ := runtime.authoring.Inspect(info.ID)
	if len(after.Document.Journey.Actions) != 0 || len(after.Document.Screen(pageID).Components) != 1 || after.Document.Screen(pageID).ComponentByID(keyboard) == nil {
		t.Fatal("remoção afetou entidades externas ou conservou binding órfão")
	}
}
