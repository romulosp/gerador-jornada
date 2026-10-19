package builtinassets

import _ "embed"

// appIcon guarda a imagem escolhida pelo autor para identificar o programa.
//
//go:embed assets/app-icon.png
var appIcon []byte

// IconPNG fornece o recurso incorporado sem permitir alteração dos bytes originais.
func IconPNG() []byte {
	return append([]byte(nil), appIcon...)
}
