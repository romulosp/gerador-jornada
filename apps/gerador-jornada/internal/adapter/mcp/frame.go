// Package mcp adapta o núcleo de autoria ao protocolo local sem inicializar a GUI.
package mcp

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// MaxFrameBytes limita cada linha completa do protocolo, incluindo o LF final.
const MaxFrameBytes = 1 << 20

// frameWriter só publica linhas completas; erro de limite nunca vaza frame parcial.
type frameWriter struct {
	// mu serializa gravações de respostas concorrentes sem intercalar mensagens.
	mu sync.Mutex
	// destination recebe exclusivamente frames completos do protocolo.
	destination io.WriteCloser
	// pending guarda no máximo um frame dentro do orçamento contratual.
	pending []byte
	// failed é terminal para impedir reutilizar fragmento após falha de protocolo/pipe.
	failed error
	// closed/closeOnce permitem cancelar pipe bloqueado sem esperar o mutex de Write.
	closed    atomic.Bool
	closeOnce sync.Once
	// closeError é publicado pelo sincronismo de closeOnce, sem fechar o pipe duas vezes.
	closeError error
}

// Write aceita fragmentos da SDK e encaminha bytes apenas ao encontrar LF válido.
func (w *frameWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed.Load() {
		return 0, io.ErrClosedPipe
	}
	if w.failed != nil {
		return 0, w.failed
	}
	// consumed contabiliza os bytes aceitos desta chamada, inclusive frames anteriores.
	consumed := 0
	for len(data) > 0 {
		// size limita o fragmento ao LF quando a chamada contém vários frames.
		size := bytes.IndexByte(data, '\n') + 1
		if size == 0 {
			size = len(data)
		}
		if len(w.pending)+size > MaxFrameBytes {
			w.pending = nil
			w.failed = fmt.Errorf("frame MCP excede 1 MiB")
			return consumed, w.failed
		}
		w.pending = append(w.pending, data[:size]...)
		consumed += size
		data = data[size:]
		if w.pending[len(w.pending)-1] == '\n' {
			// written/err detectam falha do pipe sem repetir bytes já publicados.
			written, err := w.destination.Write(w.pending)
			if err == nil && written != len(w.pending) {
				err = io.ErrShortWrite
			}
			w.pending = nil
			if err != nil {
				w.failed = err
				return consumed, err
			}
		}
	}
	return consumed, nil
}

// Close abandona fragmento incompleto e encerra o pipe sem publicar lixo.
func (w *frameWriter) Close() error {
	w.closed.Store(true)
	w.closeOnce.Do(func() { w.closeError = w.destination.Close() })
	if w.mu.TryLock() {
		w.pending = nil
		w.mu.Unlock()
	}
	return w.closeError
}
