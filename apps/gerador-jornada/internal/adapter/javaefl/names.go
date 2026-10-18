package javaefl

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// reserved registra palavras que não podem nomear pacotes ou identificadores Java SE 6.
var reserved = map[string]bool{"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true, "case": true, "catch": true, "char": true, "class": true, "const": true, "continue": true, "default": true, "do": true, "double": true, "else": true, "enum": true, "extends": true, "final": true, "finally": true, "float": true, "for": true, "goto": true, "if": true, "implements": true, "import": true, "instanceof": true, "int": true, "interface": true, "long": true, "native": true, "new": true, "package": true, "private": true, "protected": true, "public": true, "return": true, "short": true, "static": true, "strictfp": true, "super": true, "switch": true, "synchronized": true, "this": true, "throw": true, "throws": true, "transient": true, "try": true, "void": true, "volatile": true, "while": true, "true": true, "false": true, "null": true}

// packageSegment reconhece um segmento ASCII de namespace.
var packageSegment = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// ValidatePackageName rejeita namespaces inválidos e palavras reservadas.
func ValidatePackageName(value string) error {
	if value == "" {
		return fmt.Errorf("JAVA_PACKAGE_INVALID: informe o pacote de destino")
	}
	// segment representa uma parte do namespace.
	for _, segment := range strings.Split(value, ".") {
		if !packageSegment.MatchString(segment) || reserved[segment] {
			return fmt.Errorf("JAVA_PACKAGE_INVALID: %s", segment)
		}
	}
	return nil
}

// Identifier normaliza nomes sem aplicar sufixo silencioso em colisões.
func Identifier(value string, class bool) (string, error) {
	// source remove marcas de acentuação antes de dividir palavras.
	source := norm.NFD.String(value)
	// words conserva blocos ASCII alfanuméricos.
	words := []string{}
	// word acumula a palavra normalizada antes da composição do identificador Java.
	word := ""
	// character representa uma unidade Unicode.
	for _, character := range source {
		if unicode.Is(unicode.Mn, character) {
			continue
		}
		if character < 128 && (unicode.IsLetter(character) || unicode.IsDigit(character)) {
			word += string(character)
		} else if word != "" {
			words = append(words, word)
			word = ""
		}
	}
	if word != "" {
		words = append(words, word)
	}
	if len(words) == 0 {
		return "", fmt.Errorf("JAVA_IDENTIFIER_INVALID: nome sem caracteres utilizáveis")
	}
	// result é construído em PascalCase ou camelCase.
	result := ""
	// index determina a capitalização da primeira palavra.
	// part contém um segmento do nome que será normalizado para Java.
	for index, part := range words {
		if index > 0 || class {
			result += strings.ToUpper(part[:1]) + part[1:]
		} else {
			result += strings.ToLower(part[:1]) + part[1:]
		}
	}
	if result[0] >= '0' && result[0] <= '9' {
		if class {
			result = "Screen" + result
		} else {
			result = "component" + result
		}
	}
	if reserved[result] || reserved[strings.ToLower(result)] {
		return "", fmt.Errorf("JAVA_IDENTIFIER_INVALID: palavra reservada %s", result)
	}
	return result, nil
}

// Literal cria uma string Java preservando acentos e escapes de controle.
func Literal(value string) string {
	// output preserva UTF-8 e usa somente escapes reconhecidos pelo Java SE 6.
	var output strings.Builder
	output.WriteByte('"')
	// character é um ponto Unicode, não um byte isolado de texto.
	for _, character := range value {
		switch character {
		case '\\':
			output.WriteString("\\\\")
		case '"':
			output.WriteString("\\\"")
		case '\n':
			output.WriteString("\\n")
		case '\r':
			output.WriteString("\\r")
		case '\t':
			output.WriteString("\\t")
		case '\b':
			output.WriteString("\\b")
		case '\f':
			output.WriteString("\\f")
		default:
			if character < 32 || character == 127 {
				fmt.Fprintf(&output, "\\%03o", character)
			} else {
				output.WriteRune(character)
			}
		}
	}
	output.WriteByte('"')
	return output.String()
}

// Comment impede término de Javadoc e interpretação de escapes Unicode maliciosos.
func Comment(value string) string {
	return strings.NewReplacer("*/", "* /", "\\u", "\\ u", "&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(value)
}
