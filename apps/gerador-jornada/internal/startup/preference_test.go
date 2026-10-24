package startup

import "testing"

// TestPresentationPolicy protege primeira execução, repetição optativa e desmarcação.
func TestPresentationPolicy(t *testing.T) {
	// preference percorre o ciclo sem IO nem eventos simulados de reprodução.
	preference := Preference{}
	if !preference.ShouldShow() || preference.ShowOnStartup {
		t.Fatal("primeira execução deve apresentar, mantendo a repetição desmarcada")
	}
	preference.Presented = true
	if preference.ShouldShow() {
		t.Fatal("segunda execução não deve apresentar por padrão")
	}
	preference.ShowOnStartup = true
	if !preference.ShouldShow() {
		t.Fatal("opção marcada deve repetir a apresentação")
	}
	preference.ShowOnStartup = false
	if preference.ShouldShow() {
		t.Fatal("desmarcar deve impedir apresentações posteriores")
	}
}
