package startup

import (
	"errors"
	"testing"
)

// TestBothCompletionsRequired verifica carga rápida/lenta, falha, cancelamento e liberação única.
func TestBothCompletionsRequired(t *testing.T) {
	for _, videoFirst := range []bool{false, true} {
		// session recebe eventos nas duas ordens possíveis; nenhuma delas libera prematuramente.
		session := &Session{}
		session.PlaybackCompleted()
		if session.VideoEnded {
			t.Fatal("fim sem reprodução foi aceito")
		}
		session.PlaybackStarted()
		if videoFirst {
			session.PlaybackCompleted()
		} else {
			session.ResourcesCompleted(nil)
		}
		if session.BeginTransition() {
			t.Fatal("apenas uma condição liberou a tela")
		}
		if videoFirst {
			session.ResourcesCompleted(nil)
		} else {
			session.PlaybackCompleted()
		}
		if !session.BeginTransition() || session.BeginTransition() {
			t.Fatal("transição deve ocorrer exatamente uma vez")
		}
	}
	for _, cancel := range []bool{false, true} {
		// session conserva bloqueio mesmo que os demais eventos cheguem após erro ou cancelamento.
		session := &Session{}
		if cancel {
			session.Close()
		} else {
			session.Fail(errors.New("falha de áudio/vídeo"))
		}
		session.PlaybackStarted()
		session.PlaybackCompleted()
		session.ResourcesCompleted(nil)
		if session.BeginTransition() {
			t.Fatal("falha/cancelamento liberou o launcher")
		}
	}
}
