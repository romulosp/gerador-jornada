package desktop

import (
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	mcplocal "gerador-jornada.local/editor/internal/adapter/mcp"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// desktopMCPClient conecta o editor de teste ao mesmo servidor stdio usado pelo EXE.
func desktopMCPClient(t *testing.T, root string) *sdk.ClientSession {
	t.Helper()
	// pipes não alteram stdin/stdout do processo que executa os testes Fyne.
	serverInput, clientOutput := io.Pipe()
	clientInput, serverOutput := io.Pipe()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	finished := make(chan error, 1)
	go func() {
		finished <- mcplocal.Run(ctx, serverInput, serverOutput, io.Discard, mcplocal.Config{Workspace: root})
	}()
	client := sdk.NewClient(&sdk.Implementation{Name: "manual-parity-harness", Version: "1"}, nil)
	session, err := client.Connect(ctx, &sdk.IOTransport{Reader: clientInput, Writer: clientOutput}, nil)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		session.Close()
		cancel()
		serverInput.Close()
		serverOutput.Close()
		clientInput.Close()
		clientOutput.Close()
		select {
		case <-finished:
		case <-time.After(time.Second):
			t.Error("servidor não liberou transporte")
		}
	})
	return session
}

// TestManualTextUsesSameCodePointLimit aplica o formulário real e compara o limite ao contrato MCP.
func TestManualTextUsesSameCodePointLimit(t *testing.T) {
	// app/document fornecem botão padrão e evitam dados pessoais no formulário.
	app := test.NewApp()
	defer app.Quit()
	session := application.NewSession(NewJourneyDocument("Unicode manual", false))
	id, err := session.AddComponent(session.Document.Screens[0].ID, catalog.Button, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	editor := New(app, session.Document)
	defer editor.close()
	editor.Selected = id
	editor.refreshInspector()
	// submit consulta o formulário vigente, pois uma aplicação bem-sucedida recria seus widgets.
	submit := func(text string) {
		form := editor.Inspector.Objects[1].(*widget.Form)
		found := false
		for _, item := range form.Items {
			if item.Text == "Texto" {
				item.Widget.(*widget.Entry).SetText(text)
				found = true
			}
		}
		if !found {
			t.Fatal("campo textual ausente")
		}
		test.Tap(editor.InspectorActions.Objects[0].(*widget.Button))
	}
	text := strings.Repeat("😀", 65536)
	submit(text)
	if editor.page().ComponentByID(id).Properties.Button.Text != text {
		t.Fatal("formulário contou bytes ou UTF-16")
	}
	revision := editor.Session.Document.Project.Revision
	submit(text + "x")
	if editor.Session.Document.Project.Revision != revision || editor.page().ComponentByID(id).Properties.Button.Text != text {
		t.Fatal("formulário publicou excesso ou truncou conteúdo")
	}
}

// desktopMCPCall lê somente o envelope estruturado para evitar interpretar texto de compatibilidade.
func desktopMCPCall[T any](t *testing.T, session *sdk.ClientSession, name string, args any) T {
	t.Helper()
	// result conserva a separação entre erro de protocolo e erro de ferramenta.
	result, err := session.CallTool(context.Background(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil || result.IsError {
		t.Fatalf("%s: %+v (%v)", name, result, err)
	}
	encoded, _ := json.Marshal(result.StructuredContent)
	var envelope struct {
		Result T `json:"result"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope.Result
}

// TestManualMCPPresetKeepsLegacyDecoration compara o comando manual com a inserção real por MCP.
func TestManualMCPPresetKeepsLegacyDecoration(t *testing.T) {
	// app/editor exercitam o caminho manual, inclusive a decoração existente da faixa de título.
	app := test.NewApp()
	defer app.Quit()
	base := NewJourneyDocument("Paridade", false)
	// A decoração é pré-incorporada para comparar identidade estável, sem depender de UUID aleatório.
	builtinassets.TitleStrip(&base)
	root := t.TempDir()
	if err := projectzip.Save(context.Background(), filepath.Join(root, "base.jornada"), base); err != nil {
		t.Fatal(err)
	}
	editor := New(app, base)
	defer editor.close()
	editor.insertPreset("pageTitle", 0, 100)
	manual := editor.page().ComponentByID(editor.Selected)
	client := desktopMCPClient(t, root)
	opened := desktopMCPCall[struct {
		SessionID string `json:"sessionId"`
		Revision  string `json:"revision"`
	}](t, client, "jornada_abrir", map[string]any{"requestId": project.NewID(), "filename": "base.jornada"})
	desktopMCPCall[struct{}](t, client, "preset_inserir", map[string]any{"sessionId": opened.SessionID, "requestId": project.NewID(), "expectedRevision": opened.Revision, "screenId": base.Screens[0].ID, "presetId": "pageTitle", "x": 0, "y": 100})
	// query compara dados de cena sem transportar os bytes da imagem ou o pacote completo.
	query := desktopMCPCall[struct {
		Items []struct {
			Bounds     screen.Bounds `json:"bounds"`
			Style      screen.Style  `json:"style"`
			Properties struct {
				Text string `json:"text"`
			} `json:"properties"`
		} `json:"items"`
	}](t, client, "jornada_obter", map[string]any{"sessionId": opened.SessionID, "screenId": base.Screens[0].ID})
	if len(query.Items) != 1 || query.Items[0].Bounds != manual.Bounds || query.Items[0].Style != manual.Style || query.Items[0].Style.BackgroundAssetID == "" || query.Items[0].Properties.Text != manual.Properties.Title.Text {
		t.Fatal("MCP e manual divergiram em decoração/defaults da faixa")
	}
	// O documento de origem continua intacto; nenhum comando de autoria MCP salva implicitamente.
	original, err := projectzip.Load(filepath.Join(root, "base.jornada"))
	if err != nil || len(original.Screens[0].Components) != 0 || original.Screens[0].Background != base.Screens[0].Background {
		t.Fatal("autoria MCP alterou a origem ou seu fundo")
	}
}
