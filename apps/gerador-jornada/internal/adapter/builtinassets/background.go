// Package builtinassets incorpora recursos visuais do legado sem caminhos externos em execução.
package builtinassets

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// background contém o fundo original de 800 × 600 solicitado para a área de desenho.
//
//go:embed assets/fundo-loterias.jpg
var background []byte

// titleStrip contém a faixa decorativa original; seu título continua texto editável.
//
//go:embed assets/faixa-titulo.png
var titleStrip []byte

// DefaultBackgrounds aplica o fundo às novas telas sem substituir imagens já escolhidas.
func DefaultBackgrounds(document *application.Document) {
	// page identifica cada composição nova que ainda não tem uma imagem definida.
	for index := range document.Screens {
		if document.Screens[index].Background.AssetID == "" {
			document.Screens[index].Background.AssetID = incorporate(document, background, ".jpg", "image/jpeg")
		}
	}
}

// UseBackground aplica o fundo incorporado à tela indicada; a aplicação controla o histórico.
func UseBackground(document *application.Document, pageID string) error {
	// page precisa existir antes de incorporar ou alterar recursos.
	page := document.Screen(pageID)
	if page == nil {
		return fmt.Errorf("tela inexistente para aplicar o fundo")
	}
	page.Background.AssetID = incorporate(document, background, ".jpg", "image/jpeg")
	return nil
}

// TitleStrip incorpora uma faixa reutilizável que pode decorar componentes de título.
func TitleStrip(document *application.Document) string {
	return incorporate(document, titleStrip, ".png", "image/png")
}

// incorporate deduplica por conteúdo e mantém todos os bytes dentro do projeto portátil.
func incorporate(document *application.Document, data []byte, extension, mediaType string) string {
	// digest identifica o recurso independentemente do nome e da tela.
	digest := sha256.Sum256(data)
	// checksum permite compartilhar um único recurso entre várias composições.
	checksum := hex.EncodeToString(digest[:])
	// asset reutiliza um recurso existente com os mesmos bytes declarados.
	for _, asset := range document.Project.Assets {
		if asset.SHA256 == checksum && asset.MediaType == mediaType {
			if document.Resources == nil {
				document.Resources = map[string][]byte{}
			}
			document.Resources[asset.Path] = append([]byte{}, data...)
			return asset.ID
		}
	}
	// id mantém o mesmo contrato UUID dos recursos importados pelo autor.
	id := project.NewID()
	// asset descreve um arquivo portátil, sem referência à instalação do editor.
	asset := project.Asset{ID: id, Path: "assets/" + id + extension, MediaType: mediaType,
		SHA256: checksum, SizeBytes: int64(len(data))}
	if document.Resources == nil {
		document.Resources = map[string][]byte{}
	}
	document.Project.Assets = append(document.Project.Assets, asset)
	document.Resources[asset.Path] = append([]byte{}, data...)
	return id
}
