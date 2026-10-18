package builtinassets

import (
	"archive/zip"
	"bytes"
	"context"
	"image/jpeg"
	"io"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/adapter/pptx"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
)

// TestBackgroundPortability comprova que o recurso original acompanha projeto e PPT sem caminhos externos.
func TestBackgroundPortability(t *testing.T) {
	// document oferece várias telas para conferir o compartilhamento de um único recurso.
	document := application.Example()
	DefaultBackgrounds(&document)
	DefaultBackgrounds(&document)
	if len(document.Project.Assets) != 1 {
		t.Fatal("o fundo deve ser reutilizado entre as telas")
	}
	// config confirma que o recurso incorporado é a imagem original de 800 × 600.
	config, err := jpeg.DecodeConfig(bytes.NewReader(background))
	if err != nil || config.Width != 800 || config.Height != 600 {
		t.Fatalf("recurso inesperado: %+v %v", config, err)
	}
	// encoded contém um arquivo real no formato portátil do projeto.
	encoded, err := projectzip.Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	// restored precisa conservar os bytes e o mesmo vínculo em todas as telas.
	restored, err := projectzip.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range restored.Screens {
		if page.Background.AssetID != restored.Project.Assets[0].ID {
			t.Fatal("fundo perdido ao reabrir")
		}
	}
	if !bytes.Equal(restored.Resources[restored.Project.Assets[0].Path], background) {
		t.Fatal("imagem alterada ao persistir")
	}
	// deck verifica incorporação real na apresentação, sem homologar o aplicativo PowerPoint.
	deck, err := pptx.Build(context.Background(), restored, pptx.Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// archive permite conferir mídia original e a relação de imagem dos slides.
	archive, err := zip.NewReader(bytes.NewReader(deck), int64(len(deck)))
	if err != nil {
		t.Fatal(err)
	}
	// found diferencia um vínculo XML de um recurso efetivamente incorporado.
	found := false
	for _, file := range archive.File {
		if strings.HasPrefix(file.Name, "ppt/media/") {
			// stream contém a mídia publicada no pacote da apresentação.
			stream, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			// data permite comparar sem depender da renderização do Office.
			data, err := io.ReadAll(stream)
			stream.Close()
			if err != nil {
				t.Fatal(err)
			}
			found = found || bytes.Equal(data, background)
		}
	}
	if !found {
		t.Fatal("fundo ausente do PPT")
	}
}
