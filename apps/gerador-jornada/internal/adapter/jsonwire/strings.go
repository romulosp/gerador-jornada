package jsonwire

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

// ValidateStrings recusa strings Go inválidas antes que encoding/json as substitua por U+FFFD.
// Bytes binários de imagens/ZIP não são texto e não são percorridos nem copiados.
func ValidateStrings(value any) error {
	return validateStrings(reflect.ValueOf(value), map[uintptr]bool{}, 0)
}

// validateStrings percorre somente dados exportados e limita profundidade/ciclos na fronteira.
func validateStrings(value reflect.Value, ancestors map[uintptr]bool, depth int) error {
	if !value.IsValid() {
		return nil
	}
	if depth > 256 {
		return fmt.Errorf("profundidade de dados excedida")
	}
	switch value.Kind() {
	case reflect.String:
		if !utf8.ValidString(value.String()) {
			return fmt.Errorf("texto UTF-8 inválido")
		}
	case reflect.Interface:
		if !value.IsNil() {
			return validateStrings(value.Elem(), ancestors, depth+1)
		}
	case reflect.Pointer:
		if value.IsNil() {
			return nil
		}
		// address participa somente do caminho atual, permitindo referências compartilhadas sem ciclo.
		address := value.Pointer()
		if ancestors[address] {
			return fmt.Errorf("ciclo nos dados JSON")
		}
		ancestors[address] = true
		defer delete(ancestors, address)
		return validateStrings(value.Elem(), ancestors, depth+1)
	case reflect.Struct:
		// index identifica campos exportados que podem compor a representação JSON.
		for index := 0; index < value.NumField(); index++ {
			if value.Type().Field(index).PkgPath != "" {
				continue
			}
			if err := validateStrings(value.Field(index), ancestors, depth+1); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return nil
		}
		for index := 0; index < value.Len(); index++ {
			if err := validateStrings(value.Index(index), ancestors, depth+1); err != nil {
				return err
			}
		}
	case reflect.Map:
		// entries verifica nomes e valores sem interpretar referência lógica como caminho de arquivo.
		entries := value.MapRange()
		for entries.Next() {
			if err := validateStrings(entries.Key(), ancestors, depth+1); err != nil {
				return err
			}
			if err := validateStrings(entries.Value(), ancestors, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}
