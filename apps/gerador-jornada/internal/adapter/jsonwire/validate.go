// Package jsonwire valida JSON de fronteira sem normalização silenciosa ou dependência de transporte.
package jsonwire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

// Validate recusa chaves duplicadas/Unicode inválido, com limites explícitos de bytes e profundidade.
func Validate(data []byte, maxBytes, maxDepth int) error {
	if maxBytes < 1 || maxDepth < 1 || len(data) > maxBytes {
		return fmt.Errorf("resourceLimit: JSON excede limite de bytes")
	}
	if !utf8.Valid(data) || !json.Valid(data) {
		return fmt.Errorf("invalidArgument: JSON ou UTF-8 inválido")
	}
	if err := validateSurrogates(data); err != nil {
		return err
	}
	// decoder percorre tokens sem perder duplicatas por conversão prematura para map.
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := validateJSONValue(decoder, 0, maxDepth); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("invalidArgument: conteúdo JSON adicional")
	}
	return nil
}

// validateJSONValue mantém um conjunto de chaves por objeto e limita profundidade.
func validateJSONValue(decoder *json.Decoder, depth, maxDepth int) error {
	if depth > maxDepth {
		return fmt.Errorf("resourceLimit: profundidade JSON excedida")
	}
	// token preserva números como texto e delimita objetos/arrays.
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("invalidArgument: token JSON inválido")
	}
	// delimiter distingue valor escalar de coleção sem criar modelo de negócio paralelo.
	delimiter, container := token.(json.Delim)
	if !container {
		return nil
	}
	if delimiter == '{' {
		// keys pertence somente ao objeto atual, inclusive quando existem objetos aninhados.
		keys := map[string]bool{}
		for decoder.More() {
			// keyToken/key são o nome decodificado; escapes equivalentes também colidem.
			keyToken, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("invalidArgument: chave JSON inválida")
			}
			key, ok := keyToken.(string)
			if !ok || keys[key] {
				return fmt.Errorf("invalidArgument: chave JSON duplicada ou inválida")
			}
			keys[key] = true
			if err := validateJSONValue(decoder, depth+1, maxDepth); err != nil {
				return err
			}
		}
	} else if delimiter == '[' {
		for decoder.More() {
			if err := validateJSONValue(decoder, depth+1, maxDepth); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("invalidArgument: delimitador JSON inválido")
	}
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("invalidArgument: coleção JSON incompleta")
	}
	return nil
}

// validateSurrogates recusa escapes UTF-16 isolados que encoding/json substituiria.
func validateSurrogates(data []byte) error {
	// inside distingue escapes em strings de pontuação e whitespace do envelope.
	inside := false
	// index percorre bytes originais, antes de qualquer transformação Unicode.
	for index := 0; index < len(data); index++ {
		if data[index] == '"' {
			inside = !inside
			continue
		}
		if !inside || data[index] != '\\' {
			continue
		}
		if index+1 >= len(data) {
			return fmt.Errorf("invalidArgument: escape incompleto")
		}
		if data[index+1] != 'u' {
			index++
			continue
		}
		if index+6 > len(data) {
			return fmt.Errorf("invalidArgument: escape Unicode incompleto")
		}
		// first é uma unidade UTF-16, sem normalizar os bytes retidos.
		first, err := strconv.ParseUint(string(data[index+2:index+6]), 16, 16)
		if err != nil {
			return fmt.Errorf("invalidArgument: escape Unicode inválido")
		}
		if first >= 0xDC00 && first <= 0xDFFF {
			return fmt.Errorf("invalidArgument: surrogate isolado")
		}
		if first >= 0xD800 && first <= 0xDBFF {
			if index+12 > len(data) || data[index+6] != '\\' || data[index+7] != 'u' {
				return fmt.Errorf("invalidArgument: surrogate sem par")
			}
			// second precisa ser a unidade baixa imediatamente adjacente.
			second, err := strconv.ParseUint(string(data[index+8:index+12]), 16, 16)
			if err != nil || second < 0xDC00 || second > 0xDFFF {
				return fmt.Errorf("invalidArgument: par surrogate inválido")
			}
			index += 11
		} else {
			index += 5
		}
	}
	return nil
}
