package pptx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"
)

// ValidatePackage comprova integridade estrutural e resolução de vínculos OPC.
// Esta verificação não substitui homologação visual no PowerPoint.
func ValidatePackage(data []byte) error {
	// archive lê o pacote sem extrair caminhos.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	// entries permite resolver relacionamentos internos.
	entries := map[string]bool{}
	// file é uma parte ZIP.
	for _, file := range archive.File {
		if entries[file.Name] {
			return fmt.Errorf("PPT_RELATIONSHIP_INVALID: parte repetida")
		}
		entries[file.Name] = true
	}
	// file é validado conforme seu conteúdo XML.
	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".xml") && !strings.HasSuffix(file.Name, ".rels") {
			continue
		}
		// stream fornece a parte sem materializar o pacote inteiro novamente.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		stream, err := file.Open()
		if err != nil {
			return err
		}
		// decoder verifica formação de XML e destinos internos.
		decoder := xml.NewDecoder(stream)
		for {
			// token representa o próximo nó do documento.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				stream.Close()
				return fmt.Errorf("%s: %w", file.Name, err)
			}
			// start identifica relacionamentos declarados em arquivos .rels.
			// ok informa se a consulta encontrou um valor válido.
			start, ok := token.(xml.StartElement)
			if ok && strings.HasSuffix(file.Name, ".rels") && start.Name.Local == "Relationship" {
				// target é o caminho relativo fornecido pela parte.
				target := ""
				// external marca relações externas que não devem ser resolvidas como partes internas.
				external := false
				// attribute seleciona destino e modo de vínculo.
				for _, attribute := range start.Attr {
					if attribute.Name.Local == "Target" {
						target = attribute.Value
					}
					if attribute.Name.Local == "TargetMode" && attribute.Value == "External" {
						external = true
					}
				}
				// base retorna da pasta _rels à pasta da parte de origem.
				base := strings.TrimSuffix(path.Dir(file.Name), "/_rels")
				if file.Name == "_rels/.rels" {
					base = ""
				}
				if external || !entries[path.Clean(path.Join(base, target))] {
					stream.Close()
					return fmt.Errorf("PPT_RELATIONSHIP_INVALID: %s -> %s", file.Name, target)
				}
			}
		}
		stream.Close()
	}
	return nil
}
