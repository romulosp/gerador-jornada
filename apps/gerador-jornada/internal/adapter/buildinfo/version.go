// Package buildinfo fornece a identidade da distribuição sem alterar versões dos formatos de domínio.
package buildinfo

import (
	_ "embed"
	"strings"
)

// version contém a mesma fonte utilizada para os metadados do executável Windows.
//
//go:embed VERSION
var version string

// Version retorna a versão do software em quatro partes, sem espaços ou quebras de linha.
func Version() string { return strings.TrimSpace(version) }
