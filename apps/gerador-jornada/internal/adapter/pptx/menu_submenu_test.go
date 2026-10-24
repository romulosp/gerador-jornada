package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/skin"
)

// TestMegaSubmenuPPT conserva as duas opções, a faixa própria e os comandos sem título inventado.
func TestMegaSubmenuPPT(t *testing.T) {
	// document contém a mesma instância usada pelas fixtures dos três consumidores.
	document, _, _, _ := compoundDocument(t)
	// menu fixa o submenu confirmado; o caminho de uma prévia não participa da exportação.
	menu := componentByKind(t, &document.Screens[0], catalog.ApplicationMenu)
	menu.Properties.Menu.InitialMenuID = menu.Properties.Menu.Root.Items[0].Submenu.ID
	// data mantém o teste em memória, sem gerar uma distribuição PPT adicional.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// archive e media permitem conferir bytes originais, sem confiar apenas nos nomes por hash.
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	media := [][]byte{}
	// file percorre somente imagens incorporadas ao pacote temporário.
	for _, file := range archive.File {
		if !strings.HasPrefix(file.Name, "ppt/media/") {
			continue
		}
		// reader e content representam o payload original publicado.
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
	// path fixa referências independentes para faixa, duas opções e comandos do rodapé.
	for _, path := range []string{"/azul/icon_megasena_grande.png", "/botao_megasena.png", "/botao_megasena_especial.png", "/azul/botao_menu_peq_released.png"} {
		// original contém a imagem exata definida pelo XML do legado.
		original, err := skin.MustBuiltin().Asset(path)
		if err != nil {
			t.Fatal(err)
		}
		// found exige a imagem correspondente no pacote real.
		found := false
		// content compara o recurso sem usar a implementação de deduplicação como oráculo.
		for _, content := range media {
			found = found || bytes.Equal(original, content)
		}
		if !found {
			t.Fatalf("submenu exportado sem %s", path)
		}
	}
	// content verifica texto autorado/editável dos comandos e ausência do título inventado anterior.
	content := slideContent(t, archive, "ppt/slides/slide2.xml")
	if !strings.Contains(content, ">VOLTAR</a:t>") || !strings.Contains(content, ">CANCELAR</a:t>") || strings.Contains(content, ">MEGA-SENA</a:t>") {
		t.Fatal("rodapé ausente ou título duplicado no submenu")
	}
}
