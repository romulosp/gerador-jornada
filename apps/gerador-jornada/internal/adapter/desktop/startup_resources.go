package desktop

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"strings"

	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// LoadStartupResources valida e aquece os recursos utilizados pelo editor, sem construir sua janela.
// report recebe progresso 0–1; o chamador encaminha as mensagens à thread de interface.
func LoadStartupResources(ctx context.Context, report func(float64, string)) error {
	// library verifica hashes e tamanhos das skins antes de disponibilizar qualquer componente.
	report(0, "Conferindo componentes e recursos…")
	library, err := skin.Builtin()
	if err != nil {
		return err
	}
	paths := library.ResourcePaths()
	for index, path := range paths {
		// O mesmo cache é consumido posteriormente pelo canvas e pela prévia.
		if err := ctx.Err(); err != nil {
			return err
		}
		if strings.HasSuffix(strings.ToLower(path), ".ttf") {
			if sceneFont(path) == nil {
				return fmt.Errorf("Carregar fonte: %s", path)
			}
		} else if _, err := decodedSkinImage(path); err != nil {
			return fmt.Errorf("Carregar imagem %s: %w", path, err)
		}
		report(.8*float64(index+1)/float64(len(paths)), "Carregando imagens e fontes…")
	}
	// definitions valida defaults dos dez tipos, incluindo a árvore incorporada do menu.
	definitions := append(catalog.Definitions(), catalog.LegacyDefinitions()...)
	for _, definition := range definitions {
		if _, err := screen.NewComponent(definition.Kind, "inicializacao", 0, 0); err != nil {
			return err
		}
	}
	// document valida também o fundo e o ícone incorporados, sem publicar um projeto fictício.
	document := application.NewDocument("Inicialização")
	builtinassets.DefaultBackgrounds(&document)
	for _, data := range document.Resources {
		if _, _, err := image.DecodeConfig(bytes.NewReader(data)); err != nil {
			return err
		}
	}
	if _, _, err := image.DecodeConfig(bytes.NewReader(builtinassets.IconPNG())); err != nil {
		return err
	}
	report(.95, "Conferindo preferências e arquivos locais…")
	// Recentes corrompidos são opcionais: o launcher apresenta o diagnóstico e permite criar outro projeto.
	_, _ = loadRecents()
	report(1, "Carregamento concluído · aguardando o vídeo terminar…")
	return ctx.Err()
}
