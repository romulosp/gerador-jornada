package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"reflect"
	"strings"
	"testing"
)

// TestAuthoringTextCountsCodePoints preserva Unicode literal e distingue caracteres de bytes/grafemas.
func TestAuthoringTextCountsCodePoints(t *testing.T) {
	// emoji tem quatro bytes, mas um ponto de código depois de decodificar JSON.
	emoji := strings.Repeat("😀", 65536)
	if err := ValidateAuthoringText(emoji, 65536); err != nil {
		t.Fatal("contagem usou bytes ou UTF-16")
	}
	if err := ValidateAuthoringText(emoji+"x", 65536); err == nil {
		t.Fatal("limite de conteúdo ignorado")
	}
	if err := ValidateAuthoringText("e\u0301", 1); err == nil {
		t.Fatal("grafema foi contado como um ponto de código")
	}
	if err := ValidateAuthoringText(string([]byte{0xff}), 65536); err == nil {
		t.Fatal("Unicode inválido aceito")
	}
	if err := ValidateAuthoringName(strings.Repeat("ç", 256)); err != nil {
		t.Fatal(err)
	}
	if err := ValidateAuthoringName(strings.Repeat("ç", 257)); err == nil {
		t.Fatal("nome excedente aceito")
	}
	if err := ValidateAuthoringName(" \t"); err == nil {
		t.Fatal("nome vazio aceito")
	}
}

// TestComponentTextDeltaPreservesHistory permite alterar campos independentes sem apagar texto antigo excedente.
func TestComponentTextDeltaPreservesHistory(t *testing.T) {
	// session representa documento histórico legível cujo título excede o limite atual de nova entrada.
	session := NewSession(NewDocument("Histórico"))
	pageID := session.Document.Screens[0].ID
	id, err := session.AddComponent(pageID, catalog.Title, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	session.Document.Screen(pageID).ComponentByID(id).Properties.Title.Text = strings.Repeat("x", 65537)
	before := session.Document.Copy()
	enabled := false
	if err := session.ConfigureComponent(pageID, id, ComponentChanges{Enabled: &enabled}, nil); err != nil {
		t.Fatal("histórico bloqueou campo independente:", err)
	}
	if session.Document.Screen(pageID).ComponentByID(id).Properties.Title.Text != before.Screen(pageID).ComponentByID(id).Properties.Title.Text {
		t.Fatal("texto histórico foi truncado")
	}
	// properties é cópia tipada, contendo novo valor que deve ser recusado sem mutação.
	properties := session.Document.Screen(pageID).ComponentByID(id).Properties.Copy()
	properties.Title.Text = strings.Repeat("😀", 65537)
	before = session.Document.Copy()
	if err := session.ConfigureComponent(pageID, id, ComponentChanges{Properties: &properties}, nil); err == nil || !reflect.DeepEqual(before, session.Document) {
		t.Fatal("conteúdo excedente alterou documento")
	}
}

// TestOptionIdentityCollision mantém UUIDs de opções independentes das demais entidades do documento.
func TestOptionIdentityCollision(t *testing.T) {
	// session/report permitem provar colisão com a própria identidade e inclusão válida nova.
	session := NewSession(NewDocument("IDs"))
	pageID := session.Document.Screens[0].ID
	id, err := session.AddComponent(pageID, catalog.Report, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	properties := session.Document.Screen(pageID).ComponentByID(id).Properties.Copy()
	properties.Report.Rows = []catalog.Option{{ID: id, Text: "colisão"}}
	if err := session.ConfigureComponent(pageID, id, ComponentChanges{Properties: &properties}, nil); err == nil {
		t.Fatal("opção reutilizou ID de componente")
	}
	properties.Report.Rows = []catalog.Option{{ID: project.NewID(), Text: "nova"}}
	if err := session.ConfigureComponent(pageID, id, ComponentChanges{Properties: &properties}, nil); err != nil {
		t.Fatal("nova opção confundida com colisão própria:", err)
	}
}
