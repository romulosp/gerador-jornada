package mcp

import (
	"testing"
)

// TestRevisionValueKeepsUint64 garante representação decimal precisa antes de chamar o domínio.
func TestRevisionValueKeepsUint64(t *testing.T) {
	// input usa o mesmo tipo embutido pelos comandos de mutação.
	type input struct {
		Revision revisionValue `json:"revision"`
	}
	value, err := decodeArguments[input]([]byte(`{"revision":"18446744073709551615"}`))
	if err != nil || value.Revision.number() != ^uint64(0) {
		t.Fatal("máximo uint64 perdido")
	}
	// invalid distingue overflow, coerção numérica e representação não canônica.
	for _, invalid := range []string{`"18446744073709551616"`, `18446744073709551615`, `"01"`, `"-1"`, `"1e2"`} {
		if _, err := decodeArguments[input]([]byte(`{"revision":` + invalid + `}`)); err == nil {
			t.Fatalf("revisão inválida aceita: %s", invalid)
		}
	}
}

// TestMutationDTOKeepsClosedEnvelope impede ampliar o contrato por campos desconhecidos.
func TestMutationDTOKeepsClosedEnvelope(t *testing.T) {
	// input combina metadados e entidade de uma operação instalada, com JSON plano.
	input := `{"sessionId":"00000000-0000-4000-8000-000000000001","requestId":"00000000-0000-4000-8000-000000000002","expectedRevision":"1","screenId":"00000000-0000-4000-8000-000000000003","kind":"button","x":0,"y":0}`
	request, err := decodeArguments[insertRequest]([]byte(input))
	if err != nil || request.ExpectedRevision.number() != 1 || request.Kind != "button" {
		t.Fatalf("DTO plano inválido: %+v (%v)", request, err)
	}
	if err := validIDs(request.SessionID, request.RequestID, request.ScreenID); err != nil {
		t.Fatal(err)
	}
	if err := validIDs("00000000-0000-4000-8000-00000000000A"); err == nil {
		t.Fatal("UUID não canônico aceito")
	}
	// changes não aceita null como remoção de nome; omissão é o patch sem alteração.
	nullPatch := `{"sessionId":"00000000-0000-4000-8000-000000000001","requestId":"00000000-0000-4000-8000-000000000002","expectedRevision":"1","screenId":"00000000-0000-4000-8000-000000000003","changes":{"name":null}}`
	if _, err := decodeArguments[screenUpdateRequest]([]byte(nullPatch)); err == nil {
		t.Fatal("null aceito como patch genérico")
	}
}

// TestComponentFamilyDTORejectsOtherProperties conserva a separação entre Kinds e patches de família.
func TestComponentFamilyDTORejectsOtherProperties(t *testing.T) {
	// input oferece apenas a união, exercitando o mesmo schema usado por componente_atualizar.
	type input struct {
		Properties propertiesPatch `json:"properties"`
	}
	// valid cobre as sete famílias instaladas com propriedades avançadas I1.
	for _, valid := range []string{
		`{"kind":"button","text":"Continuar","variant":"confirm"}`, `{"kind":"label","text":"Nome"}`,
		`{"kind":"title","text":"Título"}`, `{"kind":"textField","initialValue":"","format":"text"}`,
		`{"kind":"radio","options":[{"id":"00000000-0000-4000-8000-000000000001","text":"Uma"}]}`,
		`{"kind":"report","title":"Relatório","rows":[]}`, `{"kind":"numericKeyboard","targetComponentId":""}`,
	} {
		if _, err := decodeArguments[input]([]byte(`{"properties":` + valid + `}`)); err != nil {
			t.Fatalf("família válida recusada: %s (%v)", valid, err)
		}
	}
	// invalid não pode ampliar a união nem aceitar campo de outra família.
	for _, invalid := range []string{`{"kind":"button","initialValue":"x"}`, `{"kind":"label","variant":"confirm"}`, `{"kind":"unknown","text":"x"}`, `{"kind":"barcodeInput","text":"x"}`, `{"text":"sem kind"}`} {
		if _, err := decodeArguments[input]([]byte(`{"properties":` + invalid + `}`)); err == nil {
			t.Fatalf("patch fora do contrato aceito: %s", invalid)
		}
	}
}
