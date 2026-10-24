package desktop

import "golang.org/x/sys/windows"

// publishWorkspaceFile move o temporário completo sem substituir um destino concorrente.
// Sem REPLACE_EXISTING, também funciona em volumes Windows que não suportam hard links.
func publishWorkspaceFile(source, destination string) error {
	// sourcePath e destinationPath conservam Unicode para a API nativa do Windows.
	sourcePath, err := windows.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	destinationPath, err := windows.UTF16PtrFromString(destination)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(sourcePath, destinationPath, windows.MOVEFILE_WRITE_THROUGH)
}
