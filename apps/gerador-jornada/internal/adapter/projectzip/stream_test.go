package projectzip

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"testing"

	"gerador-jornada.local/editor/internal/application"
)

// TestStreamMatchesExistingPackage preserva bytes determinísticos sem exigir ZIP integral em memória.
func TestStreamMatchesExistingPackage(t *testing.T) {
	// document tem todas as famílias base; referência cobre o codec já utilizado pelo editor.
	document := application.Example()
	reference, err := Encode(document)
	if err != nil {
		t.Fatal(err)
	}
	var destination bytes.Buffer
	if err := Write(context.Background(), &destination, document); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reference, destination.Bytes()) {
		t.Fatal("stream alterou o pacote determinístico")
	}
	if _, err := Decode(destination.Bytes()); err != nil {
		t.Fatal(err)
	}
}

// TestStreamKeepsPreChangeGolden compara hash calculado no codec anterior à extração Write.
func TestStreamKeepsPreChangeGolden(t *testing.T) {
	// document fixa apenas UUIDs que o construtor gera aleatoriamente; demais defaults permanecem reais.
	document := application.NewDocument("Baseline stream 038")
	document.Project.ID = "00000000-0000-4000-8000-000000000001"
	document.Screens[0].ID = "00000000-0000-4000-8000-000000000002"
	document.Journey.StartScreenID = document.Screens[0].ID
	var destination bytes.Buffer
	if err := Write(context.Background(), &destination, document); err != nil {
		t.Fatal(err)
	}
	if destination.Len() != 439 || fmt.Sprintf("%x", sha256.Sum256(destination.Bytes())) != "071d6970f7dd6b0ca04ffc5ee6db6d14af35c620be5b01dc21c2ba596cc6d55b" {
		t.Fatal("bytes do formato mudaram durante a extração")
	}
}

// shortDestination simula falha real de armazenamento sem depender de falta de espaço no SO.
type shortDestination struct{}

// Write informa short write para que a fronteira interrompa a finalização.
func (shortDestination) Write(data []byte) (int, error) { return 0, io.ErrShortWrite }

// TestStreamCancellationAndFailure não transforma ZIP parcial em artefato publicável.
func TestStreamCancellationAndFailure(t *testing.T) {
	// ctx já cancelado deve impedir até a primeira escrita do manifesto.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var destination bytes.Buffer
	if err := Write(ctx, &destination, application.NewDocument("Cancelada")); !errors.Is(err, context.Canceled) || destination.Len() != 0 {
		t.Fatal("cancelamento não foi respeitado antes de escrever")
	}
	if err := Write(context.Background(), shortDestination{}, application.NewDocument("Falha")); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("falha de destino perdida: %v", err)
	}
}
