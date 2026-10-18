// Package projectzip converte estados dos domínios para o formato ZIP portátil versionado.
package projectzip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// key converte siglas Go para os nomes camelCase do contrato.
func key(name string) string {
	name = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name, "ID", "Id"), "MS", "Ms"), "SHA256", "Sha256")
	return strings.ToLower(name[:1]) + name[1:]
}

// wire converte valores tipados somente na fronteira de serialização.
func wire(value reflect.Value) any {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		return wire(value.Elem())
	}
	if value.Type() == reflect.TypeOf(catalog.Properties{}) {
		// index localiza o item atual na coleção sem alterar sua ordem.
		for index := 0; index < value.NumField(); index++ {
			if !value.Field(index).IsNil() {
				return wire(value.Field(index))
			}
		}
		return nil
	}
	switch value.Kind() {
	case reflect.Struct:
		// object contém o objeto de transferência na fronteira de serialização.
		object := map[string]any{}
		// index localiza o item atual na coleção sem alterar sua ordem.
		for index := 0; index < value.NumField(); index++ {
			if value.Type() == reflect.TypeOf(screen.Style{}) && value.Field(index).IsZero() {
				continue
			}
			object[key(value.Type().Field(index).Name)] = wire(value.Field(index))
		}
		return object
	case reflect.Slice:
		// items preserva coleções vazias como arrays, nunca null.
		items := []any{}
		// index localiza o item atual na coleção sem alterar sua ordem.
		for index := 0; index < value.Len(); index++ {
			items = append(items, wire(value.Index(index)))
		}
		return items
	default:
		return value.Interface()
	}
}

// EncodeManifest serializa os campos explícitos da raiz e as propriedades planas dos componentes.
func EncodeManifest(document application.Document) ([]byte, error) {
	// root reúne os metadados e coleções independentes somente na fronteira JSON.
	root := wire(reflect.ValueOf(document.Project)).(map[string]any)
	root["schemaVersion"] = 1
	root["screens"] = wire(reflect.ValueOf(document.Screens))
	root["journey"] = wire(reflect.ValueOf(document.Journey))
	root["scenarios"] = wire(reflect.ValueOf(document.Scenarios))
	root["javaOptions"] = wire(reflect.ValueOf(document.JavaOptions))
	return json.MarshalIndent(root, "", "  ")
}

// DecodeManifest rejeita campos desconhecidos e restaura as variantes tipadas das propriedades.
func DecodeManifest(data []byte) (application.Document, error) {
	// root reúne os metadados e coleções independentes somente na fronteira JSON.
	root := map[string]json.RawMessage{}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := json.Unmarshal(data, &root); err != nil {
		return application.Document{}, fmt.Errorf("PROJECT_CORRUPT: %w", err)
	}
	if string(root["schemaVersion"]) != "1" || string(root["catalogVersion"]) != "1" {
		return application.Document{}, fmt.Errorf("PROJECT_VERSION_UNSUPPORTED")
	}
	// metadata contém os campos pertencentes ao domínio Projeto.
	metadata := map[string]json.RawMessage{}
	// envelope corresponde é composição tipada da aplicação.
	envelope := map[string]json.RawMessage{}
	// field enumera os nomes permitidos na raiz do contrato v1.
	// name identifica uma entrada relativa do pacote ou arquivo gerado.
	// value contém o valor corrente em validação, serialização ou edição.
	for name, value := range root {
		switch name {
		case "schemaVersion":
		case "id", "name", "revision", "catalogVersion", "assets":
			metadata[name] = value
		case "screens", "journey", "scenarios", "javaOptions":
			envelope[name] = value
		default:
			return application.Document{}, fmt.Errorf("PROJECT_SCHEMA_INVALID: campo %s", name)
		}
	}
	// required enumera entradas obrigatórias da raiz, inclusive arrays vazios.
	for _, required := range []string{"id", "name", "revision", "catalogVersion", "assets", "screens", "journey", "scenarios", "javaOptions"} {
		// ok informa se a consulta encontrou um valor válido.
		if _, ok := root[required]; !ok || string(root[required]) == "null" {
			return application.Document{}, fmt.Errorf("PROJECT_SCHEMA_INVALID: campo obrigatório %s", required)
		}
	}
	// screens é normalizado antes da decodificação tipada estrita.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	screens, err := normalizeScreens(envelope["screens"])
	if err != nil {
		return application.Document{}, err
	}
	envelope["screens"] = screens
	envelope["project"], err = json.Marshal(metadata)
	if err != nil {
		return application.Document{}, err
	}
	// content contém os bytes ou a composição intermediária verificada para a operação atual.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	content, err := json.Marshal(envelope)
	if err != nil {
		return application.Document{}, err
	}
	// document contém a composição tipada restaurada somente após validação estrutural do JSON.
	document := application.Document{}
	// decoder rejeita campos desconhecidos nas estruturas e propriedades aninhadas.
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&document); err != nil {
		return document, fmt.Errorf("PROJECT_SCHEMA_INVALID: %w", err)
	}
	if err = decoder.Decode(new(any)); err != io.EOF {
		return document, fmt.Errorf("PROJECT_SCHEMA_INVALID: conteúdo adicional")
	}
	document.Resources = map[string][]byte{}
	if err = document.Project.Validate(); err != nil {
		return document, err
	}
	return document, nil
}

// normalizeScreens expande o discriminador serializado para sua variante concreta em Go.
func normalizeScreens(data []byte) ([]byte, error) {
	// pages preserva as composições ordenadas e seus campos para validação estrita posterior.
	pages := []map[string]json.RawMessage{}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := json.Unmarshal(data, &pages); err != nil {
		return nil, err
	}
	// page representa a composição atual, resolvida sem alterar o documento original.
	for _, page := range pages {
		if string(page["components"]) == "null" || page["components"] == nil {
			return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: components obrigatório")
		}
		// components preserva os campos desconhecidos dos controles para a validação tipada posterior.
		components := []map[string]json.RawMessage{}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if err := json.Unmarshal(page["components"], &components); err != nil {
			return nil, err
		}
		// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
		for _, component := range components {
			// kind identifica a estrutura concreta das propriedades.
			kind := ""
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			if err := json.Unmarshal(component["kind"], &kind); err != nil {
				return nil, err
			}
			// variant mapeia a representação de transferência internamente sem aceitar reflexão escolhida pelo usuário.
			variant := map[string]string{"button": "Button", "label": "Label", "title": "Title", "textField": "Field", "radio": "Radio", "report": "Report", "numericKeyboard": "Keyboard"}[kind]
			if variant == "" || component["properties"] == nil || string(component["properties"]) == "null" {
				return nil, fmt.Errorf("PROJECT_SCHEMA_INVALID: família/propriedades inválidas")
			}
			// wrapped restaura a união esperada pelo decodificador tipado.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			wrapped, err := json.Marshal(map[string]json.RawMessage{variant: component["properties"]})
			if err != nil {
				return nil, err
			}
			component["properties"] = wrapped
		}
		// normalized preserva a ordem dos componentes.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		normalized, err := json.Marshal(components)
		if err != nil {
			return nil, err
		}
		page["components"] = normalized
	}
	return json.Marshal(pages)
}
