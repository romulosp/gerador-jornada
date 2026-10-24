package builtinassets

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
)

// SplashVideoSHA256 identifica os bytes originais de video-splash-programa-v2.mp4.
const SplashVideoSHA256 = "3aa7a762e16da881509be1e7802c84f24e5d5848e14f8db4986332846c7bf3b3"

// splashVideo incorpora vídeo e áudio originais, sem recompressão nem dependência do material externo.
//
//go:embed assets/splash-v2.mp4
var splashVideo []byte

// SplashVideoPath disponibiliza uma cópia íntegra no cache local para o decoder nativo.
// O cache usa hash e escrita atômica; nenhum arquivo de projeto é consultado ou modificado.
func SplashVideoPath(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	// directory pertence ao cache do usuário e nunca ao diretório de execução do EXE.
	directory, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if fmt.Sprintf("%x", sha256.Sum256(splashVideo)) != SplashVideoSHA256 {
		return "", fmt.Errorf("Vídeo incorporado não corresponde à versão v2 esperada.")
	}
	directory = filepath.Join(directory, "GeradorJornada", "splash", SplashVideoSHA256)
	if err := os.MkdirAll(directory, 0700); err != nil {
		return "", err
	}
	filename := filepath.Join(directory, "splash-v2.mp4")
	if data, err := os.ReadFile(filename); err == nil && fmt.Sprintf("%x", sha256.Sum256(data)) == SplashVideoSHA256 {
		return filename, nil
	}
	if err := atomicfile.Write(ctx, filename, splashVideo); err != nil {
		return "", err
	}
	return filename, nil
}
