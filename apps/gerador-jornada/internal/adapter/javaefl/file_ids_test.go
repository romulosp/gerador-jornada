package javaefl

import (
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestSourceFileIDs comprova nomes exatos dos arquivos de tela/componente sem reativar a exportação EFL.
func TestSourceFileIDs(t *testing.T) {
	// document apresenta nomes amigáveis diferentes dos nomes-base confirmados.
	document := application.Example()
	document.JavaOptions.PackageName = "br.exemplo"
	document.Screens[0].FileID = "PnlIdentificacaoFinal"
	document.Screens[0].Components[0].FileID = "TituloAcao"
	// profile serve exclusivamente à estrutura dos fontes; não homologa integração proprietária.
	profile := Profile{PanelType: "javax.swing.JPanel", HostSource: "public final class HostBindings {}", Factories: map[catalog.Kind]Factory{}}
	for _, definition := range catalog.Definitions() {
		profile.Factories[definition.Kind] = Factory{Type: "javax.swing.JPanel", Create: func(screen.Component) (string, error) { return "new javax.swing.JPanel()", nil }, Click: func(string, string) string { return "" }}
	}
	// files deve usar o ID literalmente, inclusive a caixa escolhida pelo autor.
	files, err := Sources(document, profile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(files["src/br/exemplo/screen/PnlIdentificacaoFinal.java"], "class PnlIdentificacaoFinal") {
		t.Fatal("arquivo da tela não usa o ID exato")
	}
	if !strings.Contains(files["src/br/exemplo/component/pnlidentificacaofinal/TituloAcao.java"], "class TituloAcao") {
		t.Fatal("arquivo do componente não usa o ID exato")
	}
	document.Screens[0].Components[0].FileID = "class"
	if _, err := Sources(document, profile); err == nil {
		t.Fatal("palavra reservada Java não pode produzir código inválido")
	}
}
