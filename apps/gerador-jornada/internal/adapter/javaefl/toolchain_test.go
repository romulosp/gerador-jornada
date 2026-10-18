package javaefl

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestGeneratedSourcesJava6 compila fontes reais do gerador com um perfil de teste sem homologar EFL.
func TestGeneratedSourcesJava6(t *testing.T) {
	// jdk habilita a integração externa apenas quando seu ambiente foi informado explicitamente.
	jdk := os.Getenv("JORNADA_TEST_JDK")
	if jdk == "" {
		t.Skip("Defina JORNADA_TEST_JDK para validar fontes gerados no JDK 17 com alvo 1.6.")
	}
	// document fornece as sete famílias, três telas e navegação circular.
	document := application.Example()
	document.JavaOptions.PackageName = "br.exemplo.jornada"
	document.Screens[0].Components[0].Properties.Title.Text = "Olá 😀 \\u000a */"
	// profile usa somente Swing público para testar a estrutura, sem afirmar compatibilidade EFL.
	profile := Profile{ID: "somente-teste-estrutural", PanelType: "javax.swing.JPanel",
		HostSource: "/** Contexto vazio usado exclusivamente pelo teste estrutural. */ public final class HostBindings { }",
		Factories:  map[catalog.Kind]Factory{}}
	// types escolhe tipos de teste; não substitui o perfil EFL do produto.
	types := map[catalog.Kind]string{catalog.Button: "javax.swing.JButton", catalog.Label: "javax.swing.JLabel",
		catalog.Title: "javax.swing.JLabel", catalog.TextField: "javax.swing.JTextField",
		catalog.Radio: "javax.swing.JPanel", catalog.Report: "javax.swing.JPanel", catalog.NumericKeyboard: "javax.swing.JPanel"}
	// kind e className cobrem cada família com uma assinatura Swing real e conhecida.
	for kind, className := range types {
		profile.Factories[kind] = Factory{Type: className, Create: func(component screen.Component) (string, error) {
			if component.Properties.Title != nil {
				return "new " + className + "(" + Literal(component.Properties.Title.Text) + ")", nil
			}
			return "new " + className + "()", nil
		}, Click: func(variable, destination string) string {
			return "        " + variable + ".addActionListener(new java.awt.event.ActionListener() {\n" +
				"            /** Delega o clique ao navegador sem construir outra tela no listener. */\n" +
				"            public void actionPerformed(java.awt.event.ActionEvent event) { navigator.show(" + destination + "); }\n        });\n"
		}}
	}
	// files contém a saída exata de Sources; Export continua bloqueado pelo perfil não homologado.
	files, err := Sources(document, profile)
	if err != nil {
		t.Fatal(err)
	}
	// directory isola os fontes gerados sem publicar um projeto EFL ao usuário.
	directory := t.TempDir()
	// name e content preservam UTF-8 e os caminhos de pacote gerados.
	for name, content := range files {
		// target fica sob o diretório temporário fornecido pelo teste.
		target := filepath.Join(directory, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// script resolve a ferramenta externa sem depender do diretório de execução do usuário.
	script, err := filepath.Abs(filepath.Join("..", "..", "..", "scripts", "java6.py"))
	if err != nil {
		t.Fatal(err)
	}
	// reportPath permite conferir a evidência produzida após todos os passos de validação.
	reportPath := filepath.Join(directory, "compatibilidade.json")
	// command executa o verificador com argumentos separados e cache previamente preparado.
	command := exec.Command("python", "-X", "utf8", script, "check", directory, "--jdk-home", jdk,
		"--offline", "--report", reportPath)
	// output mantém o diagnóstico do compilador disponível em caso de regressão.
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("fontes gerados incompatíveis: %v\n%s", err, output)
	}
	// data contém o relatório da execução atual, nunca uma evidência anterior.
	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatal(err)
	}
	// report exige que o bytecode gerado permaneça no alvo Java 6.
	var report struct {
		// Major identifica a versão binária exigida pela JVM 1.6.
		Major int `json:"classfileMajor"`
		// Status distingue validação estática de homologação EFL.
		Status string `json:"status"`
	}
	if err := json.Unmarshal(data, &report); err != nil || report.Major != 50 || !strings.Contains(report.Status, "aprovada") {
		t.Fatalf("relatório inválido: %s; %v", data, err)
	}
	if CurrentProfile().Verified {
		t.Fatal("o teste estrutural não pode homologar o perfil EFL do produto")
	}
}
