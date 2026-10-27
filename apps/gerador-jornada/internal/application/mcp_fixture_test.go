package application

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"reflect"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/journey"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// mcpFixture cria duas telas, dez Kinds e presets com identidades determinísticas.
// A imagem pequena é dado de teste; não substitui recursos ou homologação do legado.
func mcpFixture(t *testing.T) Document {
	t.Helper()
	// sequence fornece UUIDs v4 únicos sem depender da entropia na saída da fixture.
	sequence := 0
	// nextID mantém as identidades estáveis entre criações independentes.
	nextID := func() string {
		sequence++
		return fmt.Sprintf("00000000-0000-4000-8000-%012x", sequence)
	}
	// document conserva os contratos reais, substituindo somente IDs aleatórios.
	document := NewDocument("Fixture MCP 038")
	document.Project.ID = nextID()
	document.Screens = []screen.Screen{screen.New("Origem"), screen.New("Destino")}
	// index percorre as duas telas para fixar identidades e área de teste suficiente.
	for index := range document.Screens {
		document.Screens[index].ID = nextID()
		document.Screens[index].Size = screen.Size{Width: 2000, Height: 1600}
	}
	document.Journey = journey.Journey{StartScreenID: document.Screens[0].ID, Actions: []journey.Action{}}
	// kinds reúne as famílias base e os três compostos realmente instalados.
	kinds := append(catalog.Definitions(), catalog.LegacyDefinitions()...)
	// definition e index posicionam cada família na mesma tela para comparação.
	for index, definition := range kinds {
		// component aplica a factory real e conserva skin vazia para testar defaults.
		component, err := screen.NewComponent(definition.Kind, fmt.Sprintf("item%d", index+1), (index%2)*900, (index/2)*250)
		if err != nil {
			t.Fatal(err)
		}
		component.ID = nextID()
		component.Properties.RegenerateCompoundIDs(nextID)
		if component.Properties.Radio != nil {
			// optionIndex normaliza identidades de opções que não são compostos.
			for optionIndex := range component.Properties.Radio.Options {
				component.Properties.Radio.Options[optionIndex].ID = nextID()
			}
		}
		document.Screens[0].Components = append(document.Screens[0].Components, component)
	}
	// assetID reserva a imagem sintética que também representa a decoração nesta fixture sem UI.
	assetID := nextID()
	// preset insere cada conjunto real no documento sem chamar widgets.
	for _, preset := range Presets() {
		if _, err := InsertPreset(&document, document.Screens[1].ID, preset.ID, 20, 20, func(_ *Document) string { return assetID }); err != nil {
			t.Fatal(err)
		}
	}
	// index normaliza IDs produzidos pelos presets sem alterar o conteúdo.
	for index := range document.Screens[1].Components {
		document.Screens[1].Components[index].ID = nextID()
	}
	// bitmap/encoded fornecem imagem incorporada com alpha verificável.
	bitmap := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	bitmap.SetNRGBA(0, 0, color.NRGBA{R: 10, G: 20, B: 30, A: 127})
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, bitmap); err != nil {
		t.Fatal(err)
	}
	// data e digest representam bytes e identidade de recurso portável.
	data := encoded.Bytes()
	digest := sha256.Sum256(data)
	document.Project.Assets = []project.Asset{{ID: assetID, Path: "assets/fixture.png", MediaType: "image/png", SHA256: fmt.Sprintf("%x", digest), SizeBytes: int64(len(data))}}
	document.Resources["assets/fixture.png"] = append([]byte(nil), data...)
	document.Screens[1].Background.AssetID = document.Project.Assets[0].ID
	// button/field/keyboard dão origem a referências internas e navegação explícita.
	button := document.Screens[0].Components[0].ID
	field := document.Screens[0].Components[3].ID
	document.Screens[0].Components[6].Properties.Keyboard.TargetComponentID = field
	document.Journey.Actions = []journey.Action{{SourceScreenID: document.Screens[0].ID, SourceComponentID: button, Kind: "navigate", TargetScreenID: document.Screens[1].ID}}
	// override mantém exemplo de fonte/cor explícitos sem aplicar skin artificial.
	override := &document.Screens[0].Components[1]
	override.Style.FontFamily = "Arial"
	override.Style.FontSize = 18
	override.Style.TextColor = "#000000"
	EnsureFileIDs(&document)
	return document
}

// TestMCPFixture verifica determinismo, integridade e independência das cópias.
func TestMCPFixture(t *testing.T) {
	// first e second devem conter os mesmos dados, sem compartilhar coleções.
	first, second := mcpFixture(t), mcpFixture(t)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("fixture não determinística")
	}
	if len(first.Screens) != 2 || len(first.Screens[0].Components) != 10 {
		t.Fatal("fixture sem duas telas e dez Kinds")
	}
	if err := first.ValidateFileIDs(); err != nil {
		t.Fatal(err)
	}
	// diagnostic precisa estar vazio: fixture não representa rascunho corrompido.
	for _, diagnostic := range first.Validate() {
		if diagnostic.Severity == "error" {
			t.Fatalf("fixture inválida: %s", diagnostic.Code)
		}
	}
	second.Resources["assets/fixture.png"][0] = 0
	if first.Resources["assets/fixture.png"][0] == 0 {
		t.Fatal("recursos compartilhados")
	}
}
