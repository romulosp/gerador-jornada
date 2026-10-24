package desktop

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"gerador-jornada.local/editor/internal/adapter/nativevideo"
)

// windowsSplashVideo traduz o estado do player sem transportar dependências de plataforma ao domínio.
type windowsSplashVideo struct{ *nativevideo.Player }

// Snapshot converte os eventos atômicos sem callbacks Go em threads COM.
func (p windowsSplashVideo) Snapshot() splashVideoState {
	state := p.Player.Snapshot()
	return splashVideoState{Started: state.Started, Ended: state.Ended, Audio: state.Audio, Err: state.Err}
}

// splashNativeHandle acessa somente a extensão pública do Fyne após Show.
func splashNativeHandle(window fyne.Window) uintptr {
	native, ok := window.(driver.NativeWindow)
	if !ok {
		return 0
	}
	var handle uintptr
	native.RunNative(func(context any) {
		if target, ok := context.(driver.WindowsWindowContext); ok {
			handle = target.HWND
		}
	})
	return handle
}

// openVideo ocupa toda a área da apresentação, sem sobrepor qualquer texto.
func (s *Splash) openVideo(path string) error {
	// handle só é consultado depois de Show, evitando repetir o defeito de janela ainda não criada.
	handle := splashNativeHandle(s.Window)
	if handle == 0 {
		return fmt.Errorf("Janela Windows ainda não disponível para o vídeo.")
	}
	player, err := nativevideo.Open(handle, path)
	if err != nil {
		return err
	}
	s.player = windowsSplashVideo{player}
	return nil
}

// setSplashOpacity aplica o fade à janela real, sem modificar tema ou conteúdo dos componentes.
func setSplashOpacity(window fyne.Window, alpha uint8) {
	if handle := splashNativeHandle(window); handle != 0 {
		nativevideo.SetOpacity(handle, alpha)
	}
}
