package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
)

// TestTransitionsAndAlignedText verifica artefatos reais, editabilidade e avanço manual independente.
func TestTransitionsAndAlignedText(t *testing.T) {
	// document mantém um texto à direita sem mudar sua geometria.
	document := application.Example()
	document.Screens[0].Components[0].Style.TextAlign = "right"
	// saved e restored conferem que o estilo opcional atravessa o formato portátil.
	saved, err := projectzip.Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := projectzip.Decode(saved)
	if err != nil || restored.Screens[0].Components[0].Style.TextAlign != "right" {
		t.Fatalf("alinhamento perdido: %v", err)
	}
	// effect e child cobrem cada escolha pública do diálogo de exportação.
	for effect, child := range map[string]string{"": "<p:fade", "push": "<p:push", "wipe": "<p:wipe", "none": ""} {
		// output contém somente exploração, portanto nenhum slide pode receber avanço temporal.
		output, err := Build(context.Background(), restored, Options{Exploration: true, ScenarioIDs: []string{}, Transition: effect}, nil)
		if err != nil {
			t.Fatal(err)
		}
		// archive permite verificar todas as partes de slide do pacote publicado.
		archive, err := zip.NewReader(bytes.NewReader(output), int64(len(output)))
		if err != nil {
			t.Fatal(err)
		}
		// aligned exige que pelo menos uma forma textual use o parágrafo à direita.
		aligned := false
		for _, file := range archive.File {
			if !strings.HasPrefix(file.Name, "ppt/slides/slide") || !strings.HasSuffix(file.Name, ".xml") {
				continue
			}
			// stream e content representam um slide real, já verificado pelo validador OPC.
			stream, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			content, err := io.ReadAll(stream)
			stream.Close()
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(content, []byte("advTm=")) || !bytes.Contains(content, []byte(`advClick="0"`)) {
				t.Fatal("exploração recebeu avanço automático ou clique livre")
			}
			if child != "" && (!bytes.Contains(content, []byte(child)) || !bytes.Contains(content, []byte(`p14:dur="300"`))) {
				t.Fatalf("transição %s ausente em %s", effect, file.Name)
			}
			if effect == "none" && (bytes.Contains(content, []byte("p14:dur=")) || bytes.Contains(content, []byte("<p:fade"))) {
				t.Fatal("opção sem transição manteve efeito visual")
			}
			aligned = aligned || bytes.Contains(content, []byte(`<a:pPr algn="r"`))
		}
		if !aligned {
			t.Fatal("texto perdeu alinhamento no PPT")
		}
	}
}
