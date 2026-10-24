package projectzip

import (
	"encoding/json"
	"fmt"
)

// migrateMenuTargets exige o campo no catálogo 4 e preenche apenas sua ausência nos históricos.
func migrateMenuTargets(raw json.RawMessage, version int) (json.RawMessage, error) {
	// properties conserva campos desconhecidos para a rejeição posterior pelo decoder tipado.
	properties := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &properties); err != nil {
		return nil, err
	}
	// migrate visita a árvore já limitada pelo tamanho do manifesto e sem criar rotas implícitas.
	var migrate func(json.RawMessage, int) (json.RawMessage, error)
	migrate = func(raw json.RawMessage, depth int) (json.RawMessage, error) {
		if depth > 16 {
			return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: profundidade de menu excedida")
		}
		// node e items conservam a ordem e todos os atributos originais.
		node := map[string]json.RawMessage{}
		if err := json.Unmarshal(raw, &node); err != nil {
			return nil, err
		}
		items := []map[string]json.RawMessage{}
		if err := json.Unmarshal(node["items"], &items); err != nil {
			return nil, err
		}
		// item recebe destino vazio somente se o formato anterior não o possuía.
		for _, item := range items {
			if value, exists := item["targetJourney"]; exists {
				if version < 4 || string(value) == "null" {
					return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: targetJourney incompatível com o catálogo")
				}
			} else if version >= 4 {
				return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: campo obrigatório targetJourney")
			} else {
				item["targetJourney"] = json.RawMessage(`""`)
			}
			if child := item["submenu"]; len(child) > 0 && string(child) != "null" {
				// migrated conserva a identidade do ramo e migra seus filhos de forma independente.
				migrated, err := migrate(child, depth+1)
				if err != nil {
					return nil, err
				}
				item["submenu"] = migrated
			}
		}
		// content publica a coleção apenas depois da migração de todos os descendentes.
		content, err := json.Marshal(items)
		if err != nil {
			return nil, err
		}
		node["items"] = content
		return json.Marshal(node)
	}
	// root contém a árvore principal, sem mover jornadas ou arquivos durante a leitura.
	root, err := migrate(properties["root"], 1)
	if err != nil {
		return nil, err
	}
	properties["root"] = root
	return json.Marshal(properties)
}
