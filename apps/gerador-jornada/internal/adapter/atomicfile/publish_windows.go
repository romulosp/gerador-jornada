package atomicfile

import "golang.org/x/sys/windows"

// publish publica a gravação temporária por substituição no mesmo volume, com sincronização no Windows.
func publish(source, destination string) error {
	// sourceUTF16 contém o nome normalizado do arquivo temporário em UTF-16.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	sourceUTF16, err := windows.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	// destinationUTF16 identifica o caminho final em UTF-16.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	destinationUTF16, err := windows.UTF16PtrFromString(destination)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(sourceUTF16, destinationUTF16, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH)
}
