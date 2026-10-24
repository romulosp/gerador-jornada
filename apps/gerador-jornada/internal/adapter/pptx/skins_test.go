package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// TestSkinImageIsEmbedded confere que a imagem decorativa do skin entra no pacote.
func TestSkinImageIsEmbedded(t *testing.T) {
	// document contém um rótulo com skin de imagem.
	document := application.NewDocument("skin")
	component, _ := screen.NewComponent(catalog.Label, "logon", 20, 20)
	component.Style.Skin = "label_logon"
	document.Screens[0].Components = append(document.Screens[0].Components, component)
	// data contém o PPTX gerado.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// reader percorre as partes do pacote.
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	// media conta as imagens incorporadas e slide contém o XML do slide da tela.
	media := 0
	slide := ""
	for _, entry := range reader.File {
		if strings.HasPrefix(entry.Name, "ppt/media/skin-") {
			media++
		}
		if entry.Name == "ppt/slides/slide2.xml" {
			file, _ := entry.Open()
			buffer := new(bytes.Buffer)
			_, _ = buffer.ReadFrom(file)
			file.Close()
			slide = buffer.String()
		}
	}
	if media == 0 || !strings.Contains(slide, "<p:pic>") {
		t.Fatalf("imagem de skin não incorporada: media=%d", media)
	}
	if !strings.Contains(slide, "Título") && !strings.Contains(slide, "Rótulo") {
		t.Fatal("texto editável do componente ausente")
	}
}

// TestUnavailableSkinBlocksPublication impede exportar um controle com substituição silenciosa.
func TestUnavailableSkinBlocksPublication(t *testing.T) {
	// document conserva uma skin antiga cujo recurso ainda não foi fornecido.
	document := application.NewDocument("recurso indisponível")
	// component mantém a escolha do autor sem migrar ou apagar seu estilo.
	component, _ := screen.NewComponent(catalog.Button, "botao", 0, 0)
	component.Style.Skin = "botao_loteca_especial"
	document.Screens[0].Components = append(document.Screens[0].Components, component)
	// data nunca deve conter um pacote parcial quando a aparência está indisponível.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err == nil || len(data) != 0 || !strings.Contains(err.Error(), "skin indisponível") {
		t.Fatalf("publicação sem aparência resolvida: bytes=%d erro=%v", len(data), err)
	}
}

// TestLegacyShapeKeepsAlphaAndFont verifica contratos DrawingML sem rasterizar o texto.
func TestLegacyShapeKeepsAlphaAndFont(t *testing.T) {
	// primitive representa uma legenda com face física, alfa parcial e espaçamento explícito.
	primitive := presentation.Primitive{Bounds: screen.Bounds{Width: 164, Height: 68}, Text: "MEIOS DE PAGAMENTO", Fill: "#123456", FillAlpha: 128, FontFamily: "Futura BdCn BT", FontSize: 18, Italic: true, VerticalAlign: "center", TextAlign: "center", Wrap: true, LineSpacing: 0.8}
	// xml conserva as propriedades em formas e texto DrawingML editáveis.
	xml := shape(2, "menu", primitive, "", "")
	// token representa cada propriedade pública que o PowerPoint precisa reproduzir.
	for _, token := range []string{`typeface="Futura BdCn BT"`, `i="1"`, `anchor="ctr"`, `<a:alpha val="50196"/>`, `<a:spcPct val="80000"/>`, `MEIOS DE PAGAMENTO`, `txBox="1"`} {
		if !strings.Contains(xml, token) {
			t.Fatalf("propriedade ausente: %s", token)
		}
	}
}
