package desktop

import (
	"os"
	"path/filepath"
	"runtime"
)

// recoveryDirectory usa uma pasta por usuário, preservando o caminho antigo no Windows.
// Linux respeita XDG_CONFIG_HOME e macOS utiliza Application Support; erros impedem caminho relativo.
func recoveryDirectory() (string, error) {
	// base usa o diretório nativo de configuração, sem depender do diretório do executável.
	base, err := os.UserConfigDir()
	if runtime.GOOS == "windows" {
		base, err = os.UserCacheDir()
	}
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "GeradorJornada", "recovery"), nil
}
