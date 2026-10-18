package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
)

// TestEditableNavigationAndTiming valida o pacote gerado sem afirmar homologação Office.
func TestEditableNavigationAndTiming(t *testing.T) {
	// document exercita sete famílias e revisita da tela inicial.
	document := application.Example()
	// data contém exploração e demonstração.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = ValidatePackage(data); err != nil {
		t.Fatal(err)
	}
	// archive permite inspecionar partes reais do arquivo.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	// timing e links comprovam presença dos comportamentos requeridos no XML.
	timing := false
	// links registra se o pacote contém navegação entre slides.
	links := false
	// text registra se o slide contém o texto editável esperado.
	text := false
	// file representa cada parte produzida.
	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".xml") {
			continue
		}
		// stream abre o conteúdo de uma entrada do pacote para inspeção.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		stream, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		// content contém os bytes lidos da entrada para verificar o conteúdo gerado.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		content, err := io.ReadAll(stream)
		stream.Close()
		if err != nil {
			t.Fatal(err)
		}
		timing = timing || bytes.Contains(content, []byte("<p:animMotion"))
		links = links || bytes.Contains(content, []byte("ppaction://hlinksldjump"))
		text = text || bytes.Contains(content, []byte("<a:t"))
	}
	if !timing || !links || !text {
		t.Fatal("faltam movimento, links ou texto editável")
	}
	// again comprova reprodução determinística das mesmas opções.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	again, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, again) {
		t.Fatal("PPT não determinístico")
	}
}

// TestCancellationAndInvalidRoute impede saída parcial e omissão silenciosa de roteiro.
func TestCancellationAndInvalidRoute(t *testing.T) {
	// ctx já está cancelado antes da geração.
	// cancel permite interromper a operação e liberar os recursos do contexto.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := Build(ctx, application.Example(), Options{Exploration: true}, nil); err != context.Canceled {
		t.Fatalf("cancelamento incorreto: %v", err)
	}
	// document contém uma etapa inválida que deve impedir demonstração.
	document := application.Example()
	document.Scenarios[0].Steps[0].ScreenID = "ausente"
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := Build(context.Background(), document, Options{Exploration: true}, nil); err == nil {
		t.Fatal("roteiro inválido omitido")
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := Build(context.Background(), document, Options{Exploration: true, ScenarioIDs: []string{}}, nil); err != nil {
		t.Fatal("exploração independente foi bloqueada", err)
	}
}
