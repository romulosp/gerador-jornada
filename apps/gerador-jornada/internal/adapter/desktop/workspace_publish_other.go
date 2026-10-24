//go:build !windows

package desktop

import "os"

// publishWorkspaceFile publica um link exclusivo; o chamador remove o temporário após a publicação.
func publishWorkspaceFile(source, destination string) error {
	return os.Link(source, destination)
}
