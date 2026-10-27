package mcp

import (
	"context"
	"encoding/json"
	"gerador-jornada.local/editor/internal/skin"
	"io"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// testClient compõe servidor/cliente reais da SDK e expõe o registro só para asserções locais.
func testClient(t *testing.T) (*runtime, *sdk.ClientSession) {
	return testClientAt(t, t.TempDir())
}

// testClientAt permite fornecer pacotes de origem em uma raiz temporária explicitamente escolhida.
func testClientAt(t *testing.T, root string) (*runtime, *sdk.ClientSession) {
	t.Helper()
	// runtime usa somente diretório temporário e descarta auditoria em memória.
	runtime, err := newRuntime(Config{Workspace: root}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	server := newServer(runtime)
	clientTransport, serverTransport := sdk.NewInMemoryTransports()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		cancel()
		t.Fatal(err)
	}
	client := sdk.NewClient(&sdk.Implementation{Name: "authoring-harness", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	t.Cleanup(func() { session.Close(); cancel(); runtime.authoring.CloseAll() })
	return runtime, session
}

// callResult mantém o canal estruturado como origem única de dados do teste.
func callResult[T any](t *testing.T, session *sdk.ClientSession, name string, args any) T {
	t.Helper()
	// result é normalizado sem depender da representação interna escolhida pela SDK cliente.
	result, err := session.CallTool(context.Background(), &sdk.CallToolParams{Name: name, Arguments: args})
	if err != nil || result.IsError {
		t.Fatalf("%s falhou: %+v (%v)", name, result, err)
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		OK     bool `json:"ok"`
		Result T    `json:"result"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil || !envelope.OK {
		t.Fatalf("resposta inválida: %s (%v)", encoded, err)
	}
	return envelope.Result
}

// TestDiscoveryDistinguishesKindsSkinsAndPresets oferece somente capacidades instaladas.
func TestDiscoveryDistinguishesKindsSkinsAndPresets(t *testing.T) {
	// session consulta o servidor, sem chamar diretamente funções de descoberta.
	_, session := testClient(t)
	type entry struct {
		Kind               string `json:"kind"`
		Category           string `json:"category"`
		AdvancedProperties bool   `json:"advancedProperties"`
	}
	catalog := callResult[pageResult[entry]](t, session, "catalogo_listar", map[string]any{})
	if len(catalog.Items) != 10 {
		t.Fatalf("Kinds instalados: %d", len(catalog.Items))
	}
	base, compound := 0, 0
	for _, item := range catalog.Items {
		if item.Category == "family" {
			base++
			if !item.AdvancedProperties {
				t.Fatal("família I1 perdeu patch")
			}
		} else if item.Category == "legacyCompound" {
			compound++
			if item.AdvancedProperties {
				t.Fatal("composto anunciou patch avançado inexistente")
			}
		} else {
			t.Fatal("skin/preset anunciado como Kind")
		}
	}
	if base != 7 || compound != 3 {
		t.Fatal("paleta instalada divergente")
	}
	// Skins permanecem recursos disponíveis do tipo, sem imagem/pacote em base64 ou definição quebrada inserível.
	for _, entry := range catalog.Items {
		page := callResult[pageResult[skinItem]](t, session, "catalogo_listar_skins", map[string]any{"kind": entry.Kind})
		for _, item := range page.Items {
			if !item.Available {
				if item.ReasonCode == "" {
					t.Fatal("indisponibilidade sem motivo")
				}
				continue
			}
			for _, resource := range item.Resources {
				actual, found := skin.MustBuiltin().ResourceInfo(resource.Path)
				if !found || actual.SHA256 != resource.SHA256 || actual.Family != resource.FontFamily {
					t.Fatal("recurso/fonte física divergiu do manifesto")
				}
			}
		}
	}
	presets := callResult[pageResult[struct {
		ID string `json:"id"`
	}]](t, session, "catalogo_listar_presets", map[string]any{})
	if len(presets.Items) != 5 {
		t.Fatal("presets instalados divergentes")
	}
	// descriptor deve explicar limites/ação explícita e suspensão do alvo Java.
	descriptor := callResult[struct {
		Kind            string `json:"kind"`
		ImplicitActions bool   `json:"implicitActions"`
		Targets         struct {
			Java bool `json:"java"`
		} `json:"targets"`
	}](t, session, "catalogo_descrever", map[string]any{"kind": "button"})
	if descriptor.Kind != "button" || descriptor.ImplicitActions || descriptor.Targets.Java {
		t.Fatal("capacidade futura ou ação implícita oferecida")
	}
}
