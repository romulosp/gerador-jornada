// Package atomicfile publica arquivos completos preservando a saída anterior em caso de falha.
package atomicfile

import (
	"context"
	"os"
	"path/filepath"
)

// Write grava ao lado do destino e publica somente após sincronização bem-sucedida.
func Write(ctx context.Context, path string, data []byte) error {
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := ctx.Err(); err != nil {
		return err
	}
	// file mantém o arquivo aberto para limitar seu tamanho ou publicar uma gravação temporária.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	file, err := os.CreateTemp(filepath.Dir(path), ".jornada-*.tmp")
	if err != nil {
		return err
	}
	// staging identifica o arquivo temporário até a publicação bem-sucedida.
	staging := file.Name()
	defer os.Remove(staging)
	if _, err = file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return publish(staging, path)
}
