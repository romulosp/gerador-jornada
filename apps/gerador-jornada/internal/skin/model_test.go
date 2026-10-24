package skin

import (
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// TestBuiltinDefinitions confere a leitura do skin.xml incorporado.
func TestBuiltinDefinitions(t *testing.T) {
	// library contém as definições interpretadas.
	library, err := Builtin()
	if err != nil {
		t.Fatal(err)
	}
	// label_titulo_1 deve existir com fonte Futura e cor azul.
	state, ok := library.State(catalog.Label, "label_titulo_1", "default")
	if !ok || state.Font == nil {
		t.Fatal("label_titulo_1 ausente")
	}
	if state.Font.Size != 14 || state.Font.Color.Hex() != "#3E578D" {
		t.Fatalf("fonte inesperada: %+v", state.Font)
	}
	// textField type1 usa Arial 12.
	textState, ok := library.State(catalog.TextField, "type1", "default")
	if !ok || textState.Font == nil || textState.Font.Size != 12 || !strings.Contains(textState.Font.Type, "arial") {
		t.Fatalf("textField inesperado: %+v", textState)
	}
	// button botao_confirmar declara imagem e cor verde.
	buttonState, ok := library.State(catalog.Button, "botao_confirmar", "default")
	if !ok || buttonState.ImagePath == "" || buttonState.Font == nil {
		t.Fatalf("botão inesperado: %+v", buttonState)
	}
	// estado desabilitado cai para default quando ausente.
	fallback, ok := library.State(catalog.Label, "label_titulo_1", "disabled")
	if !ok || fallback.Font == nil {
		t.Fatal("fallback de estado ausente")
	}
}

// TestBuiltinAssets confere o acesso aos recursos embutidos e a ausência segura.
func TestBuiltinAssets(t *testing.T) {
	// library contém o sistema de arquivos embutido.
	library := MustBuiltin()
	if !library.HasAsset("/azul/fundo.jpg") {
		t.Fatal("fundo do skin ausente")
	}
	if library.HasAsset("/azul/inexistente.png") {
		t.Fatal("recurso inexistente aceito")
	}
}

// TestNamesForKind confere o filtro por família e a ordem estável.
func TestNamesForKind(t *testing.T) {
	// names contém as definições de botão.
	names := MustBuiltin().Names(catalog.Button)
	if len(names) < 20 {
		t.Fatalf("poucas skins de botão: %d", len(names))
	}
	// sorted confirma a ordenação exigida pelo seletor.
	for index := 1; index < len(names); index++ {
		if names[index-1] > names[index] {
			t.Fatal("nomes fora de ordem")
		}
	}
}

// TestImportedColorAndFace impede perder preto opaco, transparência ou o nome real SFNT.
func TestImportedColorAndFace(t *testing.T) {
	// library resolve os dados estáticos gerados, com integridade verificada.
	library, err := Builtin()
	if err != nil {
		t.Fatal(err)
	}
	// field usa preto sem atributo alpha no XML: a ausência significa opaco.
	field, _ := library.State(catalog.TextField, "type1", "default")
	if field.Font == nil || field.Font.Color.Hex() != "#000000" || field.Font.Color.A != 255 {
		t.Fatalf("preto opaco perdido: %+v", field.Font)
	}
	// label declara branco transparente: alpha=0 deve continuar transparente.
	label, _ := library.State(catalog.Label, "label_titulo_1", "default")
	if label.Background.A != 0 || label.Background.Hex() != "" {
		t.Fatal("fundo transparente tornou-se branco")
	}
	if family := library.FontFamily("/fonts/futura_mh.ttf"); family != "Futura Medium" {
		t.Fatalf("família SFNT incorreta: %q", family)
	}
	if _, ok := library.Lookup(catalog.Button, "botao_loteca_especial"); ok {
		t.Fatal("definição com recurso ausente foi oferecida")
	}
}
