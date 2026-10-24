package presentation

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// ValidateStyle recusa definições e fontes indisponíveis sem alterar o documento original.
// A interface pode apresentar o diagnóstico; o exportador deve interromper a publicação.
func ValidateStyle(component screen.Component) error {
	if component.Style.Skin != "" {
		if _, ok := skin.MustBuiltin().Lookup(component.Kind, component.Style.Skin); !ok {
			return fmt.Errorf("skin indisponível: %s", component.Style.Skin)
		}
	}
	if component.Style.FontFamily != "" && fontAsset(component.Style.FontFamily) == "" {
		return fmt.Errorf("fonte não incorporada: %s", component.Style.FontFamily)
	}
	if component.Kind == catalog.ApplicationMenu && component.Properties.Menu != nil {
		// checkNode verifica também submenus ainda não visitados pela prévia.
		var checkNode func(catalog.SubmenuProperties) error
		checkNode = func(node catalog.SubmenuProperties) error {
			// item identifica a definição exata de cada parte da árvore.
			for _, item := range node.Items {
				if item.Variant != "" {
					if _, ok := skin.MustBuiltin().Lookup(catalog.Button, item.Variant); !ok {
						return fmt.Errorf("skin do item %s indisponível: %s", item.Label, item.Variant)
					}
				}
				if item.Submenu != nil {
					if err := checkNode(*item.Submenu); err != nil {
						return err
					}
				}
			}
			return nil
		}
		return checkNode(component.Properties.Menu.Root)
	}
	return nil
}
