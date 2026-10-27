package projectzip

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"image"
	"io"
	"os"
	"strings"
)

// verifyingReader valida cancelamento/tamanho/hash dos mesmos bytes descomprimidos consumidos.
type verifyingReader struct {
	// ctx/source/digest limitam trabalho sem alocar uma imagem ou ZIP completos.
	ctx    context.Context
	source io.Reader
	digest hash.Hash
	// limit/count são bytes efetivos e não valores confiados ao cabeçalho ZIP.
	limit, count int64
	// terminal conserva CRC/erro que um decodificador de cabeçalho possa ler junto de bytes válidos.
	terminal error
}

// Read acumula bytes/hash em bloco limitado e conserva EOF/CRC do leitor ZIP original.
func (reader *verifyingReader) Read(data []byte) (int, error) {
	if reader.terminal != nil {
		return 0, reader.terminal
	}
	if err := reader.ctx.Err(); err != nil {
		return 0, err
	}
	// size mantém o limite mesmo quando o consumidor pede um buffer grande.
	size := min(len(data), 64<<10)
	read, err := reader.source.Read(data[:size])
	reader.count += int64(read)
	if reader.count > reader.limit {
		reader.terminal = ErrLimit
		return 0, ErrLimit
	}
	if read > 0 {
		_, _ = reader.digest.Write(data[:read])
	}
	if err != nil {
		reader.terminal = err
	}
	return read, err
}

// Verify confere pacote congelado pelo ReaderAt sem reconstruir/copiar os recursos em memória.
// Reutiliza DecodeManifest e as mesmas regras de ZIP/imagens; não modifica nem fecha source.
func Verify(ctx context.Context, source io.ReaderAt, size int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if size < 0 || size > 250<<20 {
		return ErrLimit
	}
	// archive fornece offsets do staging mantido aberto, sem reabrir caminho do filesystem.
	archive, err := zip.NewReader(source, size)
	if err != nil {
		return fmt.Errorf("PROJECT_CORRUPT: %w", err)
	}
	if len(archive.File) > 1000 {
		return ErrLimit
	}
	entries, seen := map[string]*zip.File{}, map[string]bool{}
	for _, entry := range archive.File {
		if !safePath(entry.Name) || seen[strings.ToLower(entry.Name)] || entry.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("PROJECT_CORRUPT: caminho inválido ou repetido")
		}
		seen[strings.ToLower(entry.Name)] = true
		entries[entry.Name] = entry
	}
	// manifest é limitado a 10 MiB; os demais recursos são validados apenas em fluxo.
	manifest, found := entries["manifest.json"]
	if !found {
		return fmt.Errorf("PROJECT_CORRUPT: manifesto ausente")
	}
	stream, err := manifest.Open()
	if err != nil {
		return err
	}
	reader := &verifyingReader{ctx: ctx, source: stream, digest: sha256.New(), limit: 10 << 20}
	data, err := io.ReadAll(reader)
	stream.Close()
	if err != nil {
		return err
	}
	document, err := DecodeManifest(data)
	if err != nil {
		return err
	}
	delete(entries, "manifest.json")
	total := int64(len(data))
	buffer := make([]byte, 64<<10)
	for _, asset := range document.Project.Assets {
		if err := ctx.Err(); err != nil {
			return err
		}
		// entry deve pertencer exatamente ao manifesto, sem recurso implícito ou caminho externo.
		entry, found := entries[asset.Path]
		if !found {
			return fmt.Errorf("PROJECT_RESOURCE_MISSING: recurso não encontrado")
		}
		if !safePath(asset.Path) || asset.Path == "manifest.json" || asset.SizeBytes > 25<<20 {
			return fmt.Errorf("PROJECT_RESOURCE_INVALID: metadados")
		}
		stream, err := entry.Open()
		if err != nil {
			return err
		}
		// reader é compartilhado por DecodeConfig e consumo restante, mantendo CRC/hash de todo o recurso.
		reader := &verifyingReader{ctx: ctx, source: stream, digest: sha256.New(), limit: 25 << 20}
		config, _, configErr := image.DecodeConfig(reader)
		if configErr != nil || config.Width < 1 || config.Height < 1 || config.Width > 8192 || config.Height > 8192 {
			stream.Close()
			return fmt.Errorf("PROJECT_RESOURCE_INVALID: dimensões")
		}
		_, copyErr := io.CopyBuffer(io.Discard, reader, buffer)
		stream.Close()
		if copyErr != nil {
			return copyErr
		}
		if reader.count != asset.SizeBytes || hex.EncodeToString(reader.digest.Sum(nil)) != asset.SHA256 {
			return fmt.Errorf("PROJECT_RESOURCE_INVALID: hash/tamanho")
		}
		total += reader.count
		if total > 250<<20 {
			return ErrLimit
		}
		delete(entries, asset.Path)
	}
	if len(entries) != 0 {
		return fmt.Errorf("PROJECT_SCHEMA_INVALID: recurso não declarado")
	}
	return nil
}
