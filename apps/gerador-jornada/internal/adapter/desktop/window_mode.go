package desktop

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// windowMode controla somente a apresentação da janela; não integra o documento ou seu histórico.
type windowMode struct {
	// window recebe o modo de tela cheia oferecido pelo Fyne.
	window fyne.Window
	// button permite entrar/sair do modo mesmo sem a barra de título do sistema.
	button *widget.Button
	// exitButton permite encerrar sem restaurar a janela; fica oculto fora de tela cheia.
	exitButton *widget.Button
	// exitDialog protege alterações sem empilhar solicitações de encerramento.
	exitDialog *dialog.ConfirmDialog
	// fullScreen compartilha somente o estado necessário ao callback nativo Windows.
	fullScreen atomic.Bool
	// closed impede ações enfileiradas depois do fechamento da janela.
	closed atomic.Bool
	// uninstall restaura o procedimento nativo anterior quando a captura está instalada.
	uninstall func()
}

// newWindowMode prepara o controle; a entrada inicial em tela cheia ocorre após mostrar a janela.
func newWindowMode(window fyne.Window) *windowMode {
	// mode conserva o estado transitório da janela principal.
	mode := &windowMode{window: window}
	mode.button = widget.NewButtonWithIcon("Tela cheia", theme.ViewFullScreenIcon(), mode.toggle)
	mode.exitButton = widget.NewButtonWithIcon("Sair", theme.LogoutIcon(), nil)
	mode.exitButton.Hide()
	return mode
}

// setFullScreen modifica somente a janela, mantendo zoom, seleção, conteúdo e formulário em edição.
func (m *windowMode) setFullScreen(full bool) {
	if m.closed.Load() || m.window.FullScreen() == full {
		return
	}
	m.fullScreen.Store(full)
	m.window.SetFullScreen(full)
	if full {
		m.button.SetText("Restaurar")
		m.exitButton.Show()
	} else {
		m.button.SetText("Tela cheia")
		m.exitButton.Hide()
	}
}

// toggle alterna entre a janela normal e a ocupação completa do monitor corrente.
func (m *windowMode) toggle() { m.setFullScreen(!m.window.FullScreen()) }

// handleKey prioriza F12 e a saída por Escape sem consumir atalhos normais de edição.
func (m *windowMode) handleKey(event *fyne.KeyEvent) bool {
	if event.Name == fyne.KeyF12 {
		m.toggle()
		return true
	}
	if event.Name == fyne.KeyEscape && m.window.FullScreen() {
		m.setFullScreen(false)
		return true
	}
	return false
}

// close desativa a ponte nativa sem encerrar outros programas nem alterar seus estilos de janela.
func (m *windowMode) close() {
	if !m.closed.CompareAndSwap(false, true) {
		return
	}
	if m.uninstall != nil {
		m.uninstall()
		m.uninstall = nil
	}
}

// EnableWindowMode instala a integração nativa após a janela real existir; chamadas repetidas são seguras.
func (e *Editor) EnableWindowMode() {
	if e.Mode == nil || e.Mode.closed.Load() || e.Mode.uninstall != nil {
		return
	}
	// uninstall captura somente a janela deste editor, sem polling de seleção ou renderização.
	uninstall, err := installNativeWindowMode(e.Mode)
	if err != nil {
		e.Status.SetText("Maximização automática indisponível; use Tela cheia: " + err.Error())
		return
	}
	e.Mode.uninstall = uninstall
}

// ShowWindow exibe a janela nativa antes de entrar em tela cheia pela mesma operação usada por F12.
// A ordem preserva dimensões de restauração e evita depender de uma configuração anterior à criação.
func (e *Editor) ShowWindow() {
	e.Window.Show()
	e.EnableWindowMode()
	e.Mode.setFullScreen(true)
}
