//go:build !windows

package desktop

import (
	"fmt"
	"fyne.io/fyne/v2"
)

// openVideo informa o limite da entrega Windows em vez de simular uma reprodução concluída.
func (s *Splash) openVideo(path string) error {
	return fmt.Errorf("Reprodução da splash disponível na entrega Windows x64.")
}

// setSplashOpacity mantém compilação futura sem afirmar suporte visual de outra plataforma.
func setSplashOpacity(window fyne.Window, alpha uint8) {}
