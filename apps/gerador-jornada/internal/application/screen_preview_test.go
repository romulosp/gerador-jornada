package application

import (
	"reflect"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
)

// TestScreenPreviewScope comprova que a tela selecionada independe da jornada e não navega para outra tela.
func TestScreenPreviewScope(t *testing.T) {
	// document fornece três telas conectadas, com a primeira definida como inicial.
	document := Example()
	// before permite detectar mudanças de conteúdo, recursos ou identidade no projeto original.
	before := document.Copy()
	// pageID escolhe explicitamente a segunda tela, sem modificar a tela inicial.
	pageID := document.Screens[1].ID
	// state contém somente a tela solicitada e seu estado de interação isolado.
	state, err := NewScreenPreview(document, pageID)
	if err != nil {
		t.Fatal(err)
	}
	if state.ScreenID != pageID || len(state.Document.Screens) != 1 || len(state.Document.Journey.Actions) != 0 || len(state.Document.Scenarios) != 0 {
		t.Fatal("prévia deve conter apenas a tela selecionada, sem jornada ou roteiro")
	}
	// component identifica botões que originalmente levariam a outras telas.
	for _, component := range document.Screens[1].Components {
		if component.Kind == catalog.Button {
			if err := state.Click(component.ID); err != nil || state.ScreenID != pageID {
				t.Fatalf("botão não pode sair da tela atual: %v", err)
			}
		}
	}
	if err := state.Enter(document.Screens[0].ID); err == nil {
		t.Fatal("outra tela não deve pertencer à sessão isolada")
	}
	state.Document.Screens[0].Components[0].Properties.Title.Text = "Alteração somente da prévia"
	if !reflect.DeepEqual(document, before) {
		t.Fatal("a prévia alterou o documento original")
	}
}

// TestScreenPreviewValidationAndValues permite revisar tela válida apesar de erros em outras telas/jornada.
func TestScreenPreviewValidationAndValues(t *testing.T) {
	// document conserva uma primeira tela válida com campo numérico e teclado.
	document := Example()
	// pageID é a tela que precisa continuar disponível durante a edição incompleta do projeto.
	pageID := document.Screens[0].ID
	document.Screens[1].Size.Width = 0
	document.Journey.StartScreenID = "inexistente"
	document.Journey.Actions[0].TargetScreenID = "inexistente"
	if _, err := NewPreview(document); err == nil {
		t.Fatal("prévia do projeto deveria continuar validando toda a jornada")
	}
	// state valida somente o escopo visual solicitado, sem reutilizar a navegação quebrada.
	state, err := NewScreenPreview(document, pageID)
	if err != nil {
		t.Fatal(err)
	}
	// component fornece o campo cuja edição deve permanecer local à sessão.
	for _, component := range document.Screens[0].Components {
		if component.Kind == catalog.TextField {
			if err := state.SetValue(component.ID, "123"); err != nil {
				t.Fatal(err)
			}
			if state.Values[component.ID] != "123" || document.Screens[0].ComponentByID(component.ID).Properties.Field.InitialValue != "" {
				t.Fatal("edição de campo precisa permanecer na prévia")
			}
		}
	}
	if _, err := NewScreenPreview(document, "inexistente"); err == nil {
		t.Fatal("tela ausente precisa ser rejeitada")
	}
	document.Screens[0].Size.Width = 0
	if _, err := NewScreenPreview(document, pageID); err == nil {
		t.Fatal("dados inválidos da própria tela precisam ser rejeitados")
	}
}
