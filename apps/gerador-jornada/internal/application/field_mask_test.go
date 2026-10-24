package application

import (
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestFieldFormatIntent valida a intenção de formato sem alegar máscara fiscal completa.
func TestFieldFormatIntent(t *testing.T) {
	// cases associa formato, valor e aceitação esperada.
	cases := []struct {
		format string
		value  string
		ok     bool
	}{
		{"text", "Ação 100%", true},
		{"number", "12345", true},
		{"number", "12a", false},
		{"cellular", "11999998888", true},
		{"cellular", "119999", false},
		{"cep", "01001000", true},
		{"cep", "0100100", false},
		{"cpf", "12345678901", true},
		{"cpf", "123", false},
		{"date", "31/12/2026", true},
		{"date", "2026-12-31", false},
		{"ddd", "11", true},
		{"ddd", "1", false},
		{"nis", "12345678901", true},
		{"nis", "123", false},
		{"", "", true},
	}
	// testCase percorre a matriz de formatos.
	for _, testCase := range cases {
		if err := validateFieldValue(testCase.format, testCase.value); (err == nil) != testCase.ok {
			t.Fatalf("formato %q valor %q: erro=%v, esperado ok=%v", testCase.format, testCase.value, err, testCase.ok)
		}
	}
}

// TestPreviewSetValueUsesFormat confere que o campo rejeitado não altera o valor da visita.
func TestPreviewSetValueUsesFormat(t *testing.T) {
	// document contém um campo de CPF.
	document := NewDocument("formato")
	field, _ := screen.NewComponent(catalog.TextField, "cpf", 20, 20)
	field.Properties.Field.Format = "cpf"
	document.Screens[0].Components = append(document.Screens[0].Components, field)
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	if err := preview.SetValue(field.ID, "123"); err == nil {
		t.Fatal("CPF curto aceito")
	}
	if preview.Values[field.ID] != "" {
		t.Fatal("valor inválido foi publicado")
	}
	if err := preview.SetValue(field.ID, "12345678901"); err != nil {
		t.Fatal(err)
	}
	if preview.Values[field.ID] != "12345678901" {
		t.Fatal("valor válido não foi gravado")
	}
}
