package desktop

// windowAction identifica apenas intenções transitórias de apresentação da janela.
type windowAction uint8

const (
	// winSystemCommand representa comandos nativos da barra e do menu da janela.
	winSystemCommand uintptr = 0x0112
	// winKeyDown e winKeyUp identificam pressão/liberação de teclas sem modificador Alt.
	winKeyDown uintptr = 0x0100
	winKeyUp   uintptr = 0x0101
	// winDestroy identifica o término da vida do HWND, inclusive quando o aplicativo encerra diretamente.
	winDestroy uintptr = 0x0082
	// winMaximize e winRestore identificam as intenções nativas de ampliar/restaurar a janela.
	winMaximize uintptr = 0xF030
	winRestore  uintptr = 0xF120
	// winF12 e winEscape são os códigos de teclas tratados pela apresentação da janela.
	winF12    uintptr = 0x7B
	winEscape uintptr = 0x1B
)

const (
	// windowPass encaminha a mensagem original para o procedimento do Fyne/GLFW.
	windowPass windowAction = iota
	// windowEnter converte Maximizar em tela cheia.
	windowEnter
	// windowExit restaura a janela sem fechar projeto ou diálogos.
	windowExit
	// windowToggle alterna o modo a partir de F12.
	windowToggle
	// windowConsume ignora repetição automática da tecla funcional já pressionada.
	windowConsume
)

// windowMessageAction traduz mensagens Win32 sem acessar HWND ou estado de negócio.
// Os quatro bits inferiores de WM_SYSCOMMAND são reservados pelo Windows.
func windowMessageAction(message, command, flags uintptr, full bool) windowAction {
	if message == winSystemCommand {
		switch command & 0xFFF0 {
		case winMaximize:
			return windowEnter
		case winRestore:
			if full {
				return windowExit
			}
		}
	}
	if message == winKeyDown {
		if command == winF12 {
			if flags&(1<<30) != 0 {
				return windowConsume
			}
			return windowToggle
		}
		if command == winEscape && full {
			return windowExit
		}
	}
	return windowPass
}

// applyWindowAction executa a intenção na thread de UI, sem publicar comandos de domínio.
func (m *windowMode) applyWindowAction(action windowAction) {
	if m.closed.Load() {
		return
	}
	switch action {
	case windowEnter:
		m.setFullScreen(true)
	case windowExit:
		m.setFullScreen(false)
	case windowToggle:
		m.toggle()
	}
}
