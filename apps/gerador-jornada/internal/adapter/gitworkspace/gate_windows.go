//go:build windows

package gitworkspace

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
)

// lock abre um handle sem compartilhamento; encerrar o processo libera a trava automaticamente.
// path é o arquivo privado de trava; retorna a função de liberação ou erro de acesso/ocupação.
func lock(path string) (func(), error) {
	// err impede prosseguir sem o diretório privado necessário à trava ou recuperação.
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	// name permanece válido durante a chamada CreateFile.
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	// handle impede outra instância de alterar a mesma raiz enquanto estiver aberto.
	handle, err := windows.CreateFile(name, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, fmt.Errorf("Já existe uma operação Git neste workspace; aguarde e tente novamente.")
	}
	return func() { _ = windows.CloseHandle(handle) }, nil
}
