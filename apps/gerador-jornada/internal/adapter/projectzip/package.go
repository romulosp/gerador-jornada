package projectzip

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// Save codifica o projeto completo e substitui seu destino atomicamente.
// Documentos antigos recebem uma cópia de segurança integral antes da primeira escrita 3/4.
func Save(ctx context.Context, destination string, document application.Document) error {
	// content contém os bytes ou a composição intermediária verificada para a operação atual.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	content, err := Encode(document)
	if err != nil {
		return err
	}
	if _, err = Decode(content); err != nil {
		return err
	}
	if err = backupLegacy(ctx, destination, document.Project.CatalogVersion); err != nil {
		return err
	}
	return atomicfile.Write(ctx, destination, content)
}

// backupLegacy preserva o arquivo original antes do primeiro save de catálogo incompatível.
// A falha de cópia impede a sobrescrita, conforme o contrato de formato.
func backupLegacy(ctx context.Context, destination string, catalogVersion int) error {
	if !project.RequiresCatalogMigration(catalogVersion) {
		return nil
	}
	// original contém os bytes atuais do projeto legado, quando ele já existe.
	original, err := os.ReadFile(destination)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// backup conserva também backups de migrações anteriores, sem sobrescrever o original guardado.
	backup := destination + ".bak"
	if saved, err := os.ReadFile(backup); err == nil {
		if bytes.Equal(saved, original) {
			return nil
		}
		backup = fmt.Sprintf("%s.catalog-%d.bak", destination, catalogVersion)
		if saved, err := os.ReadFile(backup); err == nil {
			if bytes.Equal(saved, original) {
				return nil
			}
			return fmt.Errorf("Backup de migração já existe e difere do original: %s", backup)
		} else if !os.IsNotExist(err) {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return atomicfile.Write(ctx, backup, original)
}

// Load lê um pacote com tamanho limitado sem substituir a sessão quando ocorre erro.
func Load(source string) (application.Document, error) {
	// file mantém o arquivo aberto para limitar seu tamanho ou publicar uma gravação temporária.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	file, err := os.Open(source)
	if err != nil {
		return application.Document{}, err
	}
	defer file.Close()
	// info impede a leitura de uma entrada comprimida sem limite.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	info, err := file.Stat()
	if err != nil {
		return application.Document{}, err
	}
	if info.Size() > 250<<20 {
		return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: pacote excede 250 MiB")
	}
	// data limita os dados mesmo se o arquivo mudar após a consulta de tamanho.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	data, err := io.ReadAll(io.LimitReader(file, (250<<20)+1))
	if err != nil {
		return application.Document{}, err
	}
	return Decode(data)
}

// Encode produz entradas ZIP determinísticas com recursos verificados.
func Encode(document application.Document) ([]byte, error) {
	// output conserva a API existente; a codificação única também pode escrever diretamente em staging.
	var output bytes.Buffer
	if err := Write(context.Background(), &output, document); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Decode valida caminhos ZIP, limites de tamanho, esquema e integridade das imagens.
func Decode(data []byte) (application.Document, error) {
	// reader lê o ZIP sem extrair caminhos não confiáveis para o disco.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: %w", err)
	}
	if len(reader.File) > 1000 {
		return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: mais de 1000 entradas")
	}
	// entries reúne o manifesto e os recursos incorporados com conteúdo limitado e validado.
	entries := map[string][]byte{}
	// seen detecta identidades ou nomes repetidos, inclusive diferenças de caixa em caminhos Windows.
	seen := map[string]bool{}
	// total soma os bytes efetivamente descomprimidos.
	total := int64(0)
	// entry representa uma entrada do pacote; caminhos são validados antes do uso.
	for _, entry := range reader.File {
		if !safePath(entry.Name) || seen[strings.ToLower(entry.Name)] || entry.Mode()&os.ModeSymlink != 0 {
			return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: caminho inválido ou repetido")
		}
		seen[strings.ToLower(entry.Name)] = true
		// limit se aplica aos bytes efetivos, sem confiar no cabeçalho ZIP.
		limit := int64(25 << 20)
		if entry.Name == "manifest.json" {
			limit = 10 << 20
		}
		// stream fornece bytes limitados e valida o checksum da entrada.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		stream, err := entry.Open()
		if err != nil {
			return application.Document{}, err
		}
		// content contém os bytes ou a composição intermediária verificada para a operação atual.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		content, err := io.ReadAll(io.LimitReader(stream, limit+1))
		stream.Close()
		if err != nil {
			return application.Document{}, err
		}
		total += int64(len(content))
		if int64(len(content)) > limit || total > 250<<20 {
			return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: limite descomprimido excedido")
		}
		entries[entry.Name] = content
	}
	// manifest descreve a revisão do projeto e é obrigatório mesmo em rascunhos.
	// ok informa se a consulta encontrou um valor válido.
	manifest, ok := entries["manifest.json"]
	if !ok {
		return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: manifesto ausente")
	}
	// document contém a composição tipada restaurada somente após validação estrutural do JSON.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	document, err := DecodeManifest(manifest)
	if err != nil {
		return document, err
	}
	delete(entries, "manifest.json")
	// asset descreve uma imagem incorporada imutável e permite validar seu conteúdo.
	for _, asset := range document.Project.Assets {
		// content contém os bytes ou a composição intermediária verificada para a operação atual.
		// ok informa se a consulta encontrou um valor válido.
		content, ok := entries[asset.Path]
		if !ok {
			return document, fmt.Errorf("PROJECT_RESOURCE_MISSING: %s", asset.Path)
		}
		if err = validateAsset(asset.Path, content, asset.SHA256, asset.SizeBytes); err != nil {
			return document, err
		}
		document.Resources[asset.Path] = content
		delete(entries, asset.Path)
	}
	if len(entries) > 0 {
		return document, fmt.Errorf("PROJECT_SCHEMA_INVALID: recurso não declarado")
	}
	return document, nil
}

// safePath aceita somente o manifesto e o espaço de recursos do contrato v1.
func safePath(name string) bool {
	if strings.ContainsAny(name, "\\:") || strings.HasPrefix(name, "/") || path.Clean(name) != name || strings.Contains(name, "..") {
		return false
	}
	return name == "manifest.json" || (strings.HasPrefix(name, "assets/") && (strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg")) && strings.Count(name, "/") == 1)
}

// validateAsset verifica bytes, nome do recurso e dimensões suportadas pelo decodificador.
func validateAsset(name string, data []byte, hash string, size int64) error {
	if !safePath(name) || name == "manifest.json" || int64(len(data)) != size || size > 25<<20 {
		return fmt.Errorf("PROJECT_RESOURCE_INVALID: %s", name)
	}
	// digest detecta conteúdo alterado independentemente do CRC do ZIP.
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != hash {
		return fmt.Errorf("PROJECT_RESOURCE_INVALID: hash %s", name)
	}
	// config limita as dimensões da imagem antes de alocá-la na interface.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < 1 || config.Height < 1 || config.Width > 8192 || config.Height > 8192 {
		return fmt.Errorf("PROJECT_RESOURCE_INVALID: dimensões %s", name)
	}
	return nil
}
