package projectzip

import (
	"encoding/json"
	"fmt"
)

// requireFields exige a presença das chaves obrigatórias em um objeto JSON.
func requireFields(object map[string]json.RawMessage, fields ...string) error {
	// name identifica o campo obrigatório corrente.
	for _, name := range fields {
		value, ok := object[name]
		if !ok || string(value) == "null" {
			return fmt.Errorf("PROJECT_SCHEMA_INVALID: campo obrigatório %s", name)
		}
	}
	return nil
}

// requireObject decodifica um objeto JSON não nulo e não vazio.
func requireObject(raw json.RawMessage) (map[string]json.RawMessage, error) {
	// object recebe o conteúdo estrito do campo.
	object := map[string]json.RawMessage{}
	if len(raw) == 0 || string(raw) == "null" || json.Unmarshal(raw, &object) != nil || object == nil {
		return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: objeto obrigatório")
	}
	return object, nil
}

// requireArray garante que o campo é um array JSON, inclusive vazio.
func requireArray(raw json.RawMessage) ([]json.RawMessage, error) {
	// items recebe a coleção ordenada do contrato.
	items := []json.RawMessage{}
	if len(raw) == 0 || string(raw) == "null" || json.Unmarshal(raw, &items) != nil {
		return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: array obrigatório")
	}
	return items, nil
}

// validateCompoundPresence valida os campos obrigatórios que o decoder tipado não distingue de zero.
func validateCompoundPresence(kind string, raw json.RawMessage) error {
	switch kind {
	case "barcodeInput":
		// properties contém somente o rótulo persistido.
		object, err := requireObject(raw)
		if err != nil {
			return err
		}
		return requireFields(object, "label")
	case "barcodeProcessed":
		// object contém a lista de blocos persistida.
		object, err := requireObject(raw)
		if err != nil {
			return err
		}
		if err = requireFields(object, "blocks"); err != nil {
			return err
		}
		// blocks exige ao menos um bloco no contrato de domínio.
		blocks, err := requireArray(object["blocks"])
		if err != nil || len(blocks) == 0 {
			return fmt.Errorf("PROJECT_SCHEMA_INVALID: blocks obrigatório")
		}
		// blockRaw representa um bloco serializado.
		for _, blockRaw := range blocks {
			block, err := requireObject(blockRaw)
			if err != nil {
				return err
			}
			if err = requireFields(block, "id", "title", "rows"); err != nil {
				return err
			}
			rows, err := requireArray(block["rows"])
			if err != nil || len(rows) == 0 {
				return fmt.Errorf("PROJECT_SCHEMA_INVALID: rows obrigatório")
			}
			// rowRaw representa uma linha serializada.
			for _, rowRaw := range rows {
				row, err := requireObject(rowRaw)
				if err != nil {
					return err
				}
				if err = requireFields(row, "id", "labelKind", "value"); err != nil {
					return err
				}
			}
		}
		return nil
	case "applicationMenu":
		// object contém a raiz e a escolha inicial.
		object, err := requireObject(raw)
		if err != nil {
			return err
		}
		if err = requireFields(object, "root", "initialMenuId"); err != nil {
			return err
		}
		root, err := requireObject(object["root"])
		if err != nil {
			return err
		}
		return validateSubmenuPresence(root)
	}
	return nil
}

// validateSubmenuPresence valida recursivamente os campos obrigatórios de um submenu.
func validateSubmenuPresence(node map[string]json.RawMessage) error {
	if err := requireFields(node, "id", "name", "title", "layout", "items"); err != nil {
		return err
	}
	items, err := requireArray(node["items"])
	if err != nil {
		return err
	}
	// itemRaw representa um item serializado do submenu.
	for _, itemRaw := range items {
		item, err := requireObject(itemRaw)
		if err != nil {
			return err
		}
		if err = requireFields(item, "id", "label", "legend", "category", "variant", "enabled", "targetScreenId", "movementCode", "movementOption", "movementName", "shortcut"); err != nil {
			return err
		}
		// enabled precisa ser booleano explícito no contrato.
		if string(item["enabled"]) != "true" && string(item["enabled"]) != "false" {
			return fmt.Errorf("PROJECT_SCHEMA_INVALID: enabled booleano obrigatório")
		}
		// submenu ausente/null representa folha; objeto exige validação recursiva.
		if submenuRaw, ok := item["submenu"]; ok && string(submenuRaw) != "null" {
			submenu, err := requireObject(submenuRaw)
			if err != nil {
				return err
			}
			if err = validateSubmenuPresence(submenu); err != nil {
				return err
			}
		}
	}
	return nil
}
