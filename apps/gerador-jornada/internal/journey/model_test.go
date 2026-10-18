package journey

import "testing"

// TestBindAndCycle protege a resolução por identidade sem impor grafo acíclico.
func TestBindAndCycle(t *testing.T) {
	// routes permite referenciar a própria tela.
	routes := Journey{StartScreenID: "A"}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := routes.Bind(Action{SourceScreenID: "A", SourceComponentID: "botao", Kind: "navigate", TargetScreenID: "A"}); err != nil {
		t.Fatal(err)
	}
	if routes.Destination("botao") != "A" {
		t.Fatal("auto-navegação rejeitada")
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := routes.Bind(Action{SourceScreenID: "A", SourceComponentID: "botao", Kind: "none"}); err != nil {
		t.Fatal(err)
	}
	if routes.Destination("botao") != "" || len(routes.Actions) != 1 {
		t.Fatal("binding não foi substituído")
	}
}
