package projectzip

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// TestRoundTripAllFamilies protege o contrato portátil e o discriminador de propriedades.
func TestRoundTripAllFamilies(t *testing.T) {
	// session recebe todas as famílias suportadas.
	session := application.NewSession(application.NewDocument("Acentos: ação"))
	// definition fornece cada família do catálogo.
	for _, definition := range catalog.Definitions() {
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if _, err := session.AddComponent(session.Document.Screens[0].ID, definition.Kind, 0, 0); err != nil {
			t.Fatal(err)
		}
	}
	// data representa o arquivo portable completo.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	data, err := Encode(session.Document)
	if err != nil {
		t.Fatal(err)
	}
	// restored precisa manter IDs, conteúdo e tipos.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	restored, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Project.Name != session.Document.Project.Name || len(restored.Screens[0].Components) != 7 {
		t.Fatal("perda de conteúdo")
	}
	// again comprova determinismo do ZIP e do manifesto.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	again, err := Encode(restored)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, again) {
		t.Fatal("serialização não determinística")
	}
}

// TestRejectUnknownFieldAndTraversal recusa extensões silenciosas e caminhos externos.
func TestRejectUnknownFieldAndTraversal(t *testing.T) {
	// data contém um manifesto válido antes da alteração indevida.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	data, err := EncodeManifest(application.NewDocument("Teste"))
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), `"revision": 0`, `"revision": 0, "surpresa": true`, 1))
	if _, err = DecodeManifest(data); err == nil {
		t.Fatal("campo desconhecido aceito")
	}
	// buffer simula um ZIP malicioso sem gravar no disco.
	buffer := &bytes.Buffer{}
	// writer insere uma entrada fora da raiz permitida.
	writer := zip.NewWriter(buffer)
	// entry contém o caminho inseguro.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	entry, err := writer.Create("../arquivo")
	if err != nil {
		t.Fatal(err)
	}
	entry.Write([]byte("x"))
	writer.Close()
	if _, err = Decode(buffer.Bytes()); err == nil {
		t.Fatal("path traversal aceito")
	}
}
