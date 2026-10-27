package mcp

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

// TestFrameReader rejeita entradas ambíguas antes de fornecer bytes à SDK.
func TestFrameReader(t *testing.T) {
	// prefix/suffix permitem montar exatamente um frame completo no limite.
	prefix, suffix := `{"jsonrpc":"2.0","id":1,"method":"ping","params":{"padding":"`, `"}}`+"\n"
	valid := prefix + strings.Repeat("a", MaxFrameBytes-len(prefix)-len(suffix)) + suffix
	// cases distingue excesso, ID inválido, chave duplicada e Unicode malformado.
	cases := []struct {
		name, input string
		valid       bool
	}{
		{"maximum", valid, true}, {"over", strings.TrimSuffix(valid, "\n") + " \n", false},
		{"duplicate", `{"id":1,"params":{"value":1,"value":2}}` + "\n", false},
		{"surrogate", `{"id":"\ud800"}` + "\n", false},
		{"lowSurrogate", `{"id":"\udc00"}` + "\n", false},
		{"pair", `{"id":"\ud83d\ude00"}` + "\n", true},
		{"literalEscape", `{"id":"\\ud800"}` + "\n", true},
		{"unsafeInteger", `{"id":9007199254740992}` + "\n", false},
		{"longID", `{"id":"` + strings.Repeat("x", 129) + `"}` + "\n", false},
		{"invalidUTF8", "{\"id\":\"" + string([]byte{255}) + "\"}\n", false},
		{"missingLF", `{"id":1}`, false},
	}
	// test é executado por valor; nenhuma fixture depende de rede ou filesystem.
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			reader := newFrameReader(io.NopCloser(strings.NewReader(test.input)))
			data, err := io.ReadAll(reader)
			if test.valid {
				if err != nil || string(data) != test.input {
					t.Fatalf("frame válido recusado: %v", err)
				}
			} else if err == nil || len(data) != 0 {
				t.Fatalf("entrada inválida foi entregue parcialmente: %d (%v)", len(data), err)
			}
		})
	}
}

// testOutput conserva bytes observados no pipe sem fechar stdout do processo de teste.
type testOutput struct{ bytes.Buffer }

// Close não altera o buffer e permite conferir que fragmento final foi descartado.
func (o *testOutput) Close() error { return nil }

// TestFrameWriter não publica fragmento antes do LF ou antes de validar limite.
func TestFrameWriter(t *testing.T) {
	// output/writer simulam escrita fragmentada da SDK, sem espelhar uma implementação específica.
	output := &testOutput{}
	writer := &frameWriter{destination: output}
	if _, err := writer.Write([]byte(`{"ok":true}`)); err != nil {
		t.Fatal(err)
	}
	if output.Len() != 0 {
		t.Fatal("publicou frame sem LF")
	}
	if _, err := writer.Write([]byte("\n")); err != nil {
		t.Fatal(err)
	}
	if output.String() != "{\"ok\":true}\n" {
		t.Fatal("frame alterado")
	}
	output.Reset()
	if _, err := writer.Write(bytes.Repeat([]byte{'a'}, MaxFrameBytes+1)); err == nil || output.Len() != 0 {
		t.Fatalf("frame excessivo publicado: %d (%v)", output.Len(), err)
	}
	if _, err := writer.Write([]byte("{}\n")); err == nil || output.Len() != 0 {
		t.Fatal("writer retomou depois de falha terminal")
	}
}

// blockedOutput representa pipe cujo cliente deixou de consumir respostas.
type blockedOutput struct {
	// entered/closed sincronizam bloqueio e cancelamento sem depender de sleep.
	entered, closed chan struct{}
}

// Write só termina quando Close cancela a saída, como ocorre com pipe desconectado.
func (o *blockedOutput) Write(data []byte) (int, error) {
	close(o.entered)
	<-o.closed
	return 0, io.ErrClosedPipe
}

// Close libera a escrita bloqueada no teste de encerramento por EOF/cancelamento.
func (o *blockedOutput) Close() error { close(o.closed); return nil }

// TestFrameWriterClose não espera mutex retido por escrita bloqueada no consumidor.
func TestFrameWriterClose(t *testing.T) {
	// output/writer simulam consumidor desconectado enquanto existe resposta em andamento.
	output := &blockedOutput{entered: make(chan struct{}), closed: make(chan struct{})}
	writer := &frameWriter{destination: output}
	// finished comprova liberação da goroutine sem deixar processo residual.
	finished := make(chan struct{})
	go func() { defer close(finished); _, _ = writer.Write([]byte("{}\n")) }()
	<-output.entered
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("Close bloqueou atrás de Write")
	}
}

// TestDecodeArguments mantém objetos fechados sem aceitar dado extra silenciosamente.
func TestDecodeArguments(t *testing.T) {
	// input exemplifica DTO tipado; o decoder não entrega maps livres ao domínio.
	type input struct {
		Name string `json:"name"`
	}
	if _, err := decodeArguments[input]([]byte(`{"name":"ok","extra":true}`)); err == nil {
		t.Fatal("campo extra aceito")
	}
	if _, err := decodeArguments[input]([]byte(`{"Name":"ok"}`)); err == nil {
		t.Fatal("nome fora do schema aceito por comparação case-insensitive")
	}
	if _, err := decodeArguments[input]([]byte(`{}`)); err == nil {
		t.Fatal("campo obrigatório ausente aceito")
	}
	if _, err := decodeArguments[input]([]byte(`{"name":"a","name":"b"}`)); err == nil {
		t.Fatal("chave duplicada aceita")
	}
	if _, err := decodeArguments[input]([]byte(`{"name":"\ud800"}`)); err == nil {
		t.Fatal("Unicode inválido aceito")
	}
	if value, err := decodeArguments[input]([]byte(`{"name":"😀"}`)); err != nil || value.Name != "😀" {
		t.Fatalf("Unicode válido alterado: %+v (%v)", value, err)
	}
}
