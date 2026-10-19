package desktop

import (
	"fmt"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"golang.org/x/sys/windows"
)

// installNativeWindowMode intercepta somente maximização e teclas de tela cheia desta janela.
// Mensagens restantes continuam no procedimento original; não altera a barra de tarefas do sistema.
func installNativeWindowMode(mode *windowMode) (func(), error) {
	// native restringe o acesso à extensão pública de janela nativa do Fyne.
	native, ok := mode.window.(driver.NativeWindow)
	if !ok {
		return nil, nil
	}
	// handle pertence exclusivamente ao editor e é obtido depois de Show.
	handle := uintptr(0)
	native.RunNative(func(context any) {
		if window, ok := context.(driver.WindowsWindowContext); ok {
			handle = window.HWND
		}
	})
	if handle == 0 {
		return nil, fmt.Errorf("janela nativa ainda não está disponível")
	}
	// user32 utiliza a DLL do sistema, independente do diretório de execução do aplicativo.
	user32 := windows.NewLazySystemDLL("user32.dll")
	// getProcedure e setProcedure mantêm a cadeia de procedimentos da própria janela.
	getProcedure := user32.NewProc("GetWindowLongPtrW")
	setProcedure := user32.NewProc("SetWindowLongPtrW")
	// callProcedure delega mensagens não tratadas ao Fyne/GLFW.
	callProcedure := user32.NewProc("CallWindowProcW")
	// index corresponde a GWLP_WNDPROC (-4), representado sem truncar o ponteiro de 64 bits.
	index := ^uintptr(3)
	// original permanece atômico porque a instalação e o callback podem usar threads distintas.
	var original atomic.Uintptr
	// previous permite delegar desde a primeira mensagem recebida durante a instalação.
	previous, _, err := getProcedure.Call(handle, index)
	if previous == 0 {
		return nil, fmt.Errorf("não foi possível ler o procedimento da janela: %v", err)
	}
	original.Store(previous)
	// escapeHeld impede que repetições da tecla vazem para a seleção depois de restaurar a janela.
	escapeHeld := false
	// callback é retido pelo runtime Windows durante a vida do processo; não recebe ponteiros Go nativos.
	var callback uintptr
	callback = windows.NewCallback(func(hwnd, message, command, flags uintptr) uintptr {
		if message == winDestroy {
			mode.closed.Store(true)
			// current preserva eventuais procedimentos inseridos depois desta captura.
			current, _, _ := getProcedure.Call(hwnd, index)
			if current == callback {
				_, _, _ = setProcedure.Call(hwnd, index, original.Load())
			}
		}
		if !mode.closed.Load() {
			if command == winEscape && escapeHeld && (message == winKeyDown || message == winKeyUp) {
				escapeHeld = message != winKeyUp
				return 0
			}
			// action consulta somente estado atômico; operações de UI são enfileiradas.
			action := windowMessageAction(message, command, flags, mode.fullScreen.Load())
			if action != windowPass {
				if command == winEscape && message == winKeyDown {
					escapeHeld = true
				}
				if action != windowConsume {
					fyne.Do(func() { mode.applyWindowAction(action) })
				}
				return 0
			}
		}
		// result conserva o processamento original de fechar, minimizar, foco, DPI e desenho.
		result, _, _ := callProcedure.Call(original.Load(), hwnd, message, command, flags)
		return result
	})
	// replaced confirma que a substituição foi efetivada e preserva eventual cadeia anterior.
	replaced, _, err := setProcedure.Call(handle, index, callback)
	if replaced == 0 {
		return nil, fmt.Errorf("não foi possível capturar a maximização: %v", err)
	}
	original.Store(replaced)
	return func() {
		// current evita sobrescrever um procedimento que outro adaptador tenha instalado depois.
		current, _, _ := getProcedure.Call(handle, index)
		if current == callback {
			_, _, _ = setProcedure.Call(handle, index, original.Load())
		}
	}, nil
}
