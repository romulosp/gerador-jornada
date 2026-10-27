package projectzip

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/application"
)

// ErrLimit identifica limite contratual de pacote/entrada antes de qualquer publicação.
var ErrLimit = errors.New("PROJECT_RESOURCE_LIMIT")

// streamWriter limita pacote comprimido e observa cancelamento em cada escrita efetiva.
type streamWriter struct {
	// ctx/destination pertencem ao trabalho local; used conta bytes efetivamente aceitos.
	ctx         context.Context
	destination io.Writer
	used        int64
}

// Write não encaminha bloco que excederia o limite de pacote, nem omite erro de short write.
func (writer *streamWriter) Write(data []byte) (int, error) {
	if err := writer.ctx.Err(); err != nil {
		return 0, err
	}
	if int64(len(data)) > (250<<20)-writer.used {
		return 0, fmt.Errorf("%w: pacote excede 250 MiB", ErrLimit)
	}
	// written/err são preservados para que nenhum chamador publique um ZIP incompleto.
	written, err := writer.destination.Write(data)
	writer.used += int64(written)
	if err == nil && written != len(data) {
		err = io.ErrShortWrite
	}
	return written, err
}

// Write codifica diretamente no destino, sem montar ZIP integral nem copiar recursos do documento.
// O chamador controla staging/rollback; erro pode deixar bytes parciais no destino transitório.
// Ordem, datas e compressão são idênticas a Encode; não fecha o destino recebido.
func Write(ctx context.Context, destination io.Writer, document application.Document) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// manifest é o único conteúdo JSON intermediário; imagens permanecem nos slices imutáveis originais.
	manifest, err := EncodeManifest(document)
	if err != nil {
		return err
	}
	if len(manifest) > 10<<20 {
		return fmt.Errorf("%w: manifesto excede 10 MiB", ErrLimit)
	}
	if len(document.Project.Assets)+1 > 1000 {
		return fmt.Errorf("%w: mais de 1000 entradas", ErrLimit)
	}
	// entries/seen validam o pacote inteiro antes de produzir qualquer cabeçalho ZIP.
	entries := map[string][]byte{"manifest.json": manifest}
	seen := map[string]bool{"manifest.json": true}
	total := int64(len(manifest))
	for _, asset := range document.Project.Assets {
		if err := ctx.Err(); err != nil {
			return err
		}
		// content referencia bytes já retidos, sem uma nova cópia integral em memória.
		content, found := document.Resources[asset.Path]
		if !found {
			return fmt.Errorf("PROJECT_RESOURCE_MISSING: %s", asset.Path)
		}
		if seen[strings.ToLower(asset.Path)] {
			return fmt.Errorf("PROJECT_RESOURCE_INVALID: caminho repetido")
		}
		if err := validateAsset(asset.Path, content, asset.SHA256, asset.SizeBytes); err != nil {
			return err
		}
		seen[strings.ToLower(asset.Path)] = true
		total += int64(len(content))
		if total > 250<<20 {
			return fmt.Errorf("%w: total descomprimido excede 250 MiB", ErrLimit)
		}
		entries[asset.Path] = content
	}
	// names mantém a ordem lexical da codificação original, independente da ordem dos mapas.
	names := []string{}
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	writer := zip.NewWriter(&streamWriter{ctx: ctx, destination: destination})
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return err
		}
		// entry recebe blocos limitados; nenhuma extração ou interpretação de caminho ocorre.
		entry, err := writer.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
		if err != nil {
			return err
		}
		content := entries[name]
		for offset := 0; offset < len(content); {
			if err := ctx.Err(); err != nil {
				return err
			}
			// end limita trabalho/cancelamento também quando o recurso comprime quase inteiramente.
			end := min(offset+(64<<10), len(content))
			written, err := entry.Write(content[offset:end])
			if err != nil {
				return err
			}
			if written != end-offset {
				return io.ErrShortWrite
			}
			offset = end
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return writer.Close()
}
