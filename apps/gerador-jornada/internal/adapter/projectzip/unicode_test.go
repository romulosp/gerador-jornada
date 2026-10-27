package projectzip

import (
	"bytes"
	"testing"

	"gerador-jornada.local/editor/internal/application"
)

// TestManifestRejectsAmbiguousUnicode impede que uma origem inválida vire U+FFFD ou sobrescreva chave silenciosamente.
func TestManifestRejectsAmbiguousUnicode(t *testing.T) {
	// manifest é válido antes de introduzir corrupção específica, sem alterar os demais campos do formato.
	manifest, err := EncodeManifest(application.NewDocument("Base"))
	if err != nil {
		t.Fatal(err)
	}
	invalid := [][]byte{
		bytes.Replace(manifest, []byte(`"name": "Base"`), []byte(`"name": "\ud800"`), 1),
		bytes.Replace(manifest, []byte(`"name": "Base"`), append(append([]byte(`"name": "`), 0xff), []byte(`"`)...), 1),
		bytes.Replace(manifest, []byte(`"name": "Base"`), []byte(`"name":"Primeiro","name":"Segundo"`), 1),
		bytes.Replace(manifest, []byte(`"name": "Tela 1"`), []byte(`"name":"Uma","\u006eame":"Outra"`), 1),
	}
	for index, data := range invalid {
		if bytes.Equal(data, manifest) {
			t.Fatal("fixture não modificou entrada")
		}
		if _, err := DecodeManifest(data); err == nil {
			t.Fatalf("corrupção %d aceita/normalizada", index)
		}
	}
	// Invalidar o texto Go antes de serializar também deve falhar, em vez de gravar substituição.
	document := application.NewDocument("Base")
	document.Project.Name = string([]byte{0xff})
	if _, err := EncodeManifest(document); err == nil {
		t.Fatal("encoder substituiu texto inválido silenciosamente")
	}
}
