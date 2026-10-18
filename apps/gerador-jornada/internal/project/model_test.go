package project

import "testing"

// TestIdentity verifica formato e independência das identidades.
func TestIdentity(t *testing.T) {
	// first e second precisam ser identidades diferentes e canônicas.
	first := New("A")
	// second fornece uma segunda identidade para verificar a independência dos projetos.
	second := New("A")
	if !ValidID(first.ID) || !ValidID(second.ID) || first.ID == second.ID {
		t.Fatal("UUID inválido ou repetido")
	}
	first.Name = ""
	if first.Validate() == nil {
		t.Fatal("nome vazio aceito")
	}
}
