package catalog

import "testing"

// TestParseBarcodeFixtures valida as conversões e os derivados dos casos F01/F02 do contrato.
func TestParseBarcodeFixtures(t *testing.T) {
	// f01Digitized contém a linha digitável de referência com 47 dígitos.
	f01Digitized := "00190500954014481606906809350314337370000000100"
	// f01Read contém o código de barras de 44 dígitos correspondente.
	f01Read := "00193373700000001000500940144816060680935031"
	// cases enumera as duas representações e os derivados esperados.
	cases := []struct {
		name      string
		raw       string
		read      string
		digitized string
		bank      string
		cents     int64
		factor    int
		sigcb     bool
		title     string
	}{
		{"F01 linha digitável", f01Digitized, f01Read, f01Digitized, "001", 100, 3737, false, ""},
		{"F01 lido", f01Read, f01Read, f01Digitized, "001", 100, 3737, false, ""},
		{"F02 lido sintético", "10498100000000001001234560000100010000000003", "10498100000000001001234560000100010000000003", "10491234566000010001000000000034810000000000100", "104", 100, 1000, true, "000100010000000003"},
		{"F02 digitado sintético", "10491234566000010001000000000034810000000000100", "10498100000000001001234560000100010000000003", "10491234566000010001000000000034810000000000100", "104", 100, 1000, true, "000100010000000003"},
	}
	// testCase representa um caso de referência do contrato.
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// barcode contém o resultado validado.
			barcode, err := ParseBarcode(testCase.raw)
			if err != nil {
				t.Fatalf("ParseBarcode(%s) falhou: %v", testCase.name, err)
			}
			if barcode.Read != testCase.read || barcode.Digitized != testCase.digitized {
				t.Fatalf("representações divergentes: %+v", barcode)
			}
			if barcode.Bank != testCase.bank || barcode.ValueCents != testCase.cents || barcode.Factor != testCase.factor {
				t.Fatalf("derivados divergentes: %+v", barcode)
			}
			if barcode.Sigcb != testCase.sigcb || barcode.TitleNumber != testCase.title {
				t.Fatalf("SIGCB/título divergentes: %+v", barcode)
			}
		})
	}
}

// TestParseBarcodeRejectsMutation confere que DVs adulterados não produzem resultado parcial.
func TestParseBarcodeRejectsMutation(t *testing.T) {
	// f01Digitized contém a linha válida usada nas mutações.
	f01Digitized := "00190500954014481606906809350314337370000000100"
	// mutate troca um dígito por outro sem alterar o comprimento.
	mutate := func(value string, index int) string {
		// digit contém o novo algarismo, sempre diferente do original.
		digit := byte('0')
		if value[index] == '0' {
			digit = '1'
		}
		return value[:index] + string(digit) + value[index+1:]
	}
	// candidates reúne as adulterações dos três DVs locais e do DV geral.
	candidates := []string{
		mutate(f01Digitized, 9),
		mutate(f01Digitized, 20),
		mutate(f01Digitized, 31),
		mutate(f01Digitized, 32),
	}
	// candidate é uma linha adulterada.
	for _, candidate := range candidates {
		if candidate == f01Digitized {
			t.Fatalf("mutação sem efeito: %s", candidate)
		}
		if _, err := ParseBarcode(candidate); err == nil {
			t.Fatalf("mutação aceita indevidamente: %s", candidate)
		}
	}
	// invalid contém comprimentos e caracteres recusados pela captura.
	for _, invalid := range []string{"", "123", "0019050095401448160690680935031433737000000010", "001905009540144816069068093503143373700000001000", "0019050095401448160690680935031433737000000010A", "0019050095401448160690680935031433737000000010\n0"} {
		if _, err := NormalizeBarcode(invalid); err == nil {
			t.Fatalf("entrada inválida aceita: %q", invalid)
		}
	}
}

// TestBarcodeLabelText confere os cinco rótulos derivados do enum.
func TestBarcodeLabelText(t *testing.T) {
	// expected associa cada enum ao texto visível exigido.
	expected := map[string]string{"cpf": "CPF:", "cnpj": "CNPJ:", "name": "Nome:", "tradeName": "Nome fantasia:", "corporateName": "Razão social:"}
	// kind e text percorrem o contrato fechado.
	for kind, text := range expected {
		// actual contém o texto derivado para o enum.
		actual, ok := BarcodeLabelText(kind)
		if !ok || actual != text {
			t.Fatalf("rótulo %s: obtido %q, esperado %q", kind, actual, text)
		}
	}
	if _, ok := BarcodeLabelText("unknown"); ok {
		t.Fatal("enum desconhecido não pode receber rótulo aproximado")
	}
}
