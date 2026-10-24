// Package nativevideo reproduz áudio e vídeo pelo Media Foundation do Windows, sem player externo.
package nativevideo

/*
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -lmfplay -lmfuuid -luuid -lole32 -luser32 -lgdi32 -lstdc++
#include "player_windows.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// State descreve os eventos reais observados do player, incluindo presença da faixa de áudio.
type State struct {
	// Started e Ended vêm dos eventos PLAY e PLAYBACK_ENDED, nunca de duração estimada.
	Started, Ended bool
	// Audio informa faixa de áudio selecionada; volume do sistema continua sob controle do usuário.
	Audio bool
	// Err informa HRESULT de falha; erro nunca equivale ao término do vídeo.
	Err error
}

// Player mantém somente a instância nativa; Open/Resize/Close devem usar a thread de UI.
type Player struct{ handle unsafe.Pointer }

// Open reproduz o caminho local usando o tamanho físico do cliente HWND, sem conversão de DPI do Fyne.
func Open(parent uintptr, path string) (*Player, error) {
	// filename conserva nomes Unicode ao atravessar a fronteira nativa.
	filename, err := windows.UTF16FromString(path)
	if err != nil {
		return nil, err
	}
	// result recebe o HRESULT da inicialização sem exigir callbacks Go em threads COM.
	var result C.int32_t
	handle := C.jornada_video_open(C.uintptr_t(parent), (*C.uint16_t)(unsafe.Pointer(&filename[0])), &result)
	if handle == nil {
		return nil, fmt.Errorf("Reproduzir vídeo: HRESULT 0x%08X", uint32(result))
	}
	return &Player{handle: handle}, nil
}

// Snapshot lê eventos atômicos sem bloquear a interface.
func (p *Player) Snapshot() State {
	// native contém uma cópia independente dos campos de estado.
	var native C.JornadaVideoSnapshot
	C.jornada_video_snapshot(p.handle, &native)
	state := State{Started: native.started != 0, Ended: native.ended != 0, Audio: native.audio != 0}
	if native.error != 0 {
		state.Err = fmt.Errorf("Reproduzir vídeo: HRESULT 0x%08X", uint32(native.error))
	}
	return state
}

// Resize consulta novamente o cliente físico e encaixa o quadro completo, ocupando a maior área proporcional possível.
// O retorno informa mudança de geometria para repintar somente as margens descobertas.
func (p *Player) Resize() bool {
	return C.jornada_video_resize(p.handle) != 0
}

// Close libera a instância uma única vez e encerra o áudio.
func (p *Player) Close() {
	if p.handle != nil {
		C.jornada_video_close(p.handle)
		p.handle = nil
	}
}

// SetOpacity aplica alfa 0–255 para a transição visual na mesma janela.
func SetOpacity(window uintptr, alpha uint8) {
	C.jornada_window_alpha(C.uintptr_t(window), C.int(alpha))
}

// FitRect devolve x/y/largura/altura em pixels do mesmo encaixe usado pela superfície nativa.
// O vídeo ocupa a maior área proporcional possível sem exceder a tela, inclusive sob DPI elevado.
func FitRect(width, height, sourceWidth, sourceHeight int) [4]int {
	// rect atravessa somente quatro inteiros, sem ponteiro Go retido pela implementação nativa.
	var rect [4]C.int
	C.jornada_video_fit(C.int(width), C.int(height), C.int(sourceWidth), C.int(sourceHeight), &rect[0])
	return [4]int{int(rect[0]), int(rect[1]), int(rect[2]), int(rect[3])}
}

// DisplaySize interpreta a razão MFPlay como proporção e conserva a resolução nativa do quadro.
func DisplaySize(width, height, aspectWidth, aspectHeight int) [2]int {
	// size é o mesmo cálculo utilizado pelo player, sem duplicação de algoritmo no teste.
	var size [2]C.int
	C.jornada_video_display_size(C.int(width), C.int(height), C.int(aspectWidth), C.int(aspectHeight), &size[0])
	return [2]int{int(size[0]), int(size[1])}
}
