//go:build !windows

package desktop

// installNativeWindowMode mantém a maximização do sistema nas plataformas sem integração Win32.
// Botão de tela cheia e teclas encaminhadas pelo canvas continuam usando a API pública Fyne.
func installNativeWindowMode(_ *windowMode) (func(), error) { return nil, nil }
