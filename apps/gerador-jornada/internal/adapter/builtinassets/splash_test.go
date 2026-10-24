package builtinassets

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

// TestSplashV2Integrity verifica origem, áudio/vídeo e cache sem perder bytes do MP4.
func TestSplashV2Integrity(t *testing.T) {
	t.Setenv("LOCALAPPDATA", t.TempDir())
	if fmt.Sprintf("%x", sha256.Sum256(splashVideo)) != SplashVideoSHA256 || !bytes.Contains(splashVideo, []byte("soun")) || !bytes.Contains(splashVideo, []byte("vide")) {
		t.Fatal("MP4 v2 perdeu integridade ou uma das trilhas")
	}
	// path identifica somente o cache isolado do teste; uma cópia corrompida precisa ser substituída.
	path, err := SplashVideoPath(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("cache corrompido"), 0600); err != nil {
		t.Fatal(err)
	}
	again, err := SplashVideoPath(context.Background())
	if err != nil || again != path {
		t.Fatalf("cache não reparado: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, splashVideo) {
		t.Fatalf("bytes originais alterados: %v", err)
	}
	// ctx cancelado impede iniciar extração; não representa conclusão bem-sucedida.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := SplashVideoPath(ctx); err == nil {
		t.Fatal("cancelamento ignorado")
	}
}
