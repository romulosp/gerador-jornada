package javaefl

import (
	"context"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
)

// TestProfileGate impede apresentar código não homologado como uma saída validada.
func TestProfileGate(t *testing.T) {
	// err precisa indicar a dependência externa sem criar fontes.
	err := Export(context.Background(), t.TempDir()+"/fontes", application.Example(), CurrentProfile())
	if err == nil || !strings.Contains(err.Error(), "JAVA_PROFILE_UNVERIFIED") {
		t.Fatal("perfil não homologado foi publicado")
	}
}

// TestNamesAndEscaping verifica compatibilidade lexical sem depender de EFL.
func TestNamesAndEscaping(t *testing.T) {
	// value contém o valor corrente em validação, serialização ou edição.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if value, err := Identifier("Confirmação final", true); err != nil || value != "ConfirmacaoFinal" {
		t.Fatalf("normalização incorreta: %s %v", value, err)
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := Identifier("class", false); err == nil {
		t.Fatal("palavra reservada aceita")
	}
	if ValidatePackageName("br.gov.class") == nil {
		t.Fatal("pacote reservado aceito")
	}
	if !strings.Contains(Literal("a\n\"b"), "\\n") {
		t.Fatal("quebra não escapada")
	}
	if strings.Contains(Literal("😀"), "\\U") {
		t.Fatal("escape Go inválido em Java")
	}
	if strings.Contains(Comment("*/ \\u000a"), "*/") || strings.Contains(Comment("*/ \\u000a"), "\\u") {
		t.Fatal("comentário inseguro")
	}
}
