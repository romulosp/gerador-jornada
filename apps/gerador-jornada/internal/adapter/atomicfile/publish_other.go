//go:build !windows

package atomicfile

import "os"

// publish publica a gravação temporária por substituição no mesmo volume, com sincronização no Windows.
func publish(source, destination string) error { return os.Rename(source, destination) }
