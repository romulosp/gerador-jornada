package screen

import (
	"strings"
	"testing"
)

// TestFileIDRules protege normalização, prefixo, dispositivos Windows e limite inclusivo de 40.
func TestFileIDRules(t *testing.T) {
	if got := SuggestFileID("01 · Identificação", true); got != "PnlIdentificacao" {
		t.Fatalf("sugestão de tela incorreta: %s", got)
	}
	if got := CleanFileID("Btn Confirmação_01!/.java"); got != "BtnConfirmacao01java" {
		t.Fatalf("caracteres especiais permaneceram no ID: %s", got)
	}
	// valid enumera os dois contextos e o limite exato sem contar extensão.
	for _, valid := range []struct {
		// value é o nome-base permitido.
		value string
		// panel exige Pnl somente para telas.
		panel bool
	}{{"PnlIdentificacao", true}, {"BtnContinuar", false}, {"continuar1", false}, {"Pnl" + strings.Repeat("A", 37), true}, {strings.Repeat("B", 40), false}} {
		if err := ValidateFileID(valid.value, valid.panel); err != nil {
			t.Fatalf("ID válido recusado: %s: %v", valid.value, err)
		}
	}
	// invalid cobre entradas diretas que não devem ser aceitas por UI, importação ou geração.
	for _, invalid := range []string{"", "Identificação", "Btn Continuar", "../arquivo", "a_b", "a-b", "a.java", "1Nome", "CON", "nul", "COM1", "LPT9", strings.Repeat("A", 41)} {
		if err := ValidateFileID(invalid, false); err == nil {
			t.Fatalf("ID inválido aceito: %q", invalid)
		}
	}
	for _, invalid := range []string{"Pnl", "Identificacao", "pnlIdentificacao"} {
		if err := ValidateFileID(invalid, true); err == nil {
			t.Fatalf("prefixo inválido aceito: %s", invalid)
		}
	}
}
