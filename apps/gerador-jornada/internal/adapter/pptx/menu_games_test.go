package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/skin"
)

// TestMenuGamePPTImages comprova que o PPT incorpora cada jogo original, inclusive disabled.
func TestMenuGamePPTImages(t *testing.T) {
	// document usa uma tela composta real; o menu será ajustado antes de exportar.
	document, _, _, _ := compoundDocument(t)
	// menu permite configurar um jogo desabilitado e outro com legenda autorada.
	menu := componentByKind(t, &document.Screens[0], catalog.ApplicationMenu)
	menu.Properties.Menu.Root.Items[0].Label = "Apostar agora"
	// Milionária é o décimo jogo ativo; LotoGol e Federal são acréscimos explícitos só da fixture.
	menu.Properties.Menu.Root.Items[9].Enabled = false
	menu.Properties.Menu.Root.Items = append(menu.Properties.Menu.Root.Items, catalog.MenuItem{ID: project.NewID(), Label: "LOTOGOL", Category: "JOGOS", Variant: "botao_lotogol", Enabled: true}, catalog.MenuItem{ID: project.NewID(), Label: "LOTERIA FEDERAL", Category: "JOGOS", Variant: "botao_federal", Enabled: true})
	// data contém o pacote estático produzido pela mesma cena do canvas.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// archive permite comparar cada mídia com os arquivos originais sem rasterizar a tela.
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	// media reúne os bytes realmente publicados no pacote.
	media := [][]byte{}
	// file percorre as partes OPC sem depender de seus nomes por hash.
	for _, file := range archive.File {
		if strings.HasPrefix(file.Name, "ppt/media/") {
			// reader e content permitem validar o conteúdo, em vez de só contar imagens.
			reader, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			content, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
			media = append(media, content)
		}
	}
	// path é a referência literal das doze decorações, incluindo a exceção Dupla Sena.
	for _, path := range []string{"/botao_megasena.png", "/botao_lotofacil.png", "/botao_quina.png", "/botao_lotomania.png", "/botao_loteca.png", "/botao_lotogol.png", "/botao_dupla_sena.png", "/botao_dia_de_sorte.png", "/botao_super_sete.png", "/botao_timemania.png", "/botao_milionaria_disabled.png", "/botao_federal.png"} {
		// original representa o PNG exato, com seu ícone e wordmark incorporados.
		original, err := skin.MustBuiltin().Asset(path)
		if err != nil {
			t.Fatal(err)
		}
		// found precisa comprovar a incorporação do recurso, não apenas sua oferta no catálogo.
		found := false
		// content permite comparar sem depender da lógica de deduplicação do exportador.
		for _, content := range media {
			found = found || bytes.Equal(original, content)
		}
		if !found {
			t.Fatalf("PPT sem imagem original: %s", path)
		}
	}
	// content verifica que o autor conserva texto editável, sem duplicar marcas do bitmap.
	content := slideContent(t, archive, "ppt/slides/slide2.xml")
	if !strings.Contains(content, "Apostar agora") || strings.Contains(content, ">+MILIONÁRIA</a:t>") || strings.Contains(content, ">MEGA-SENA</a:t>") {
		t.Fatal("legenda autorada perdida ou marca duplicada sobre PNG")
	}
}
