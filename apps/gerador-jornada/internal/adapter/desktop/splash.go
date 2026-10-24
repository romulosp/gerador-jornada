package desktop

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/startup"
)

// splashVideoState traduz eventos técnicos do adaptador multimídia sem expor COM ao coordenador.
type splashVideoState struct {
	Started, Ended, Audio bool
	Err                   error
}

// splashPlayback é a fronteira mínima consumida pela janela durante a reprodução.
type splashPlayback interface {
	Snapshot() splashVideoState
	Resize() bool
	Close()
}

// Splash coordena apresentação e trabalhos assíncronos; as condições de liberação pertencem a startup.
type Splash struct {
	// Window é reutilizada pelo launcher para a transição sem fechar a aplicação.
	Window fyne.Window
	// VideoArea ocupa o cliente da janela; textos e autoria pertencem somente ao launcher.
	VideoArea *canvas.Rectangle
	// gate contém somente estado de domínio e é acessado na thread de UI.
	gate startup.Session
	// ctx e cancel encerram tarefas e animações quando o usuário fecha a janela.
	ctx    context.Context
	cancel context.CancelFunc
	// player pertence à thread da interface, assim como sua criação e destruição nativas.
	player splashPlayback
	// next monta a tela seguinte; finished é chamado somente depois da transição completa.
	next     func(fyne.Window) fyne.Window
	finished func()
	// startedAt permite observar o tempo real de reprodução sem usar tempo para liberá-la.
	startedAt time.Time
}

// NewSplash monta a tela, sem iniciar trabalhos antes de Show.
// next recebe a janela existente e devolve a janela que será revelada após o fade.
func NewSplash(app fyne.App, next func(fyne.Window) fyne.Window, finished func()) *Splash {
	prepareApplicationTheme(app, app.Preferences().String(themePreferenceKey))
	// ctx permite interromper carregamento e áudio junto com a janela, sem avançar para o launcher.
	ctx, cancel := context.WithCancel(context.Background())
	splash := &Splash{Window: app.NewWindow("Jornada Studio · Inicializando"), ctx: ctx, cancel: cancel, next: next, finished: finished}
	splash.VideoArea = canvas.NewRectangle(color.Black)
	splash.VideoArea.SetMinSize(fyne.NewSize(800, 450))
	splash.Window.SetPadded(false)
	splash.Window.SetContent(splash.VideoArea)
	splash.Window.Resize(fyne.NewSize(960, 540))
	splash.Window.CenterOnScreen()
	splash.Window.SetOnClosed(splash.close)
	return splash
}

// Show inicia reprodução e carregamento em paralelo. load pode acrescentar leitura de documento -open.
func (s *Splash) Show(load func(context.Context, func(float64, string)) error) {
	// A apresentação usa janela normal para não depender da troca de modo fullscreen do driver.
	s.Window.SetFullScreen(false)
	s.Window.Show()
	LogApplication("splash_open", map[string]any{"video": "v2", "audio_enabled": true, "fullscreen": false})
	// Recursos e extração do MP4 não bloqueiam o loop gráfico nem dependem um do outro.
	go func() {
		path, err := builtinassets.SplashVideoPath(s.ctx)
		fyne.Do(func() {
			if s.gate.Closed || s.gate.Err != nil {
				return
			}
			if err != nil {
				s.fail(err)
				return
			}
			if err := s.openVideo(path); err != nil {
				s.fail(err)
			}
		})
	}()
	go func() {
		if load == nil {
			load = LoadStartupResources
		}
		// O progresso permanece parte do contrato do carregador, sem controles ou legendas na apresentação.
		err := load(s.ctx, func(progress float64, message string) {})
		fyne.Do(func() {
			if s.gate.Closed {
				return
			}
			s.gate.ResourcesCompleted(err)
			if err != nil {
				s.fail(err)
				return
			}
			LogApplication("startup_resources_ready", nil)
			s.maybeTransition()
		})
	}()
	go func() {
		// ticker observa eventos do player; um relógio jamais substitui PLAYBACK_ENDED.
		ticker := time.NewTicker(40 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				fyne.Do(s.poll)
			}
		}
	}()
}

// poll mantém o vídeo na região reservada e encaminha eventos para o domínio.
func (s *Splash) poll() {
	if s.gate.Closed || s.gate.Transitioning || s.player == nil || s.gate.Err != nil {
		return
	}
	// O adaptador consulta pixels reais do cliente nativo, sem reaplicar a escala lógica do Fyne.
	if s.player.Resize() {
		s.VideoArea.Refresh()
	}
	state := s.player.Snapshot()
	if state.Err != nil {
		s.fail(state.Err)
		return
	}
	if state.Started && !s.gate.Playing {
		s.gate.PlaybackStarted()
		s.startedAt = time.Now()
		LogApplication("splash_video_started", map[string]any{"audio_track": state.Audio, "muted": false})
	}
	if state.Ended && !s.gate.VideoEnded {
		s.gate.PlaybackCompleted()
		LogApplication("splash_video_ended", map[string]any{"elapsed_ms": time.Since(s.startedAt).Milliseconds(), "audio_track": state.Audio})
	}
	s.maybeTransition()
}

// fail conserva a tela com diagnóstico; falha de mídia nunca permite ignorar o vídeo obrigatório.
func (s *Splash) fail(err error) {
	s.gate.Fail(err)
	if s.player != nil {
		s.player.Close()
		s.player = nil
	}
	s.cancel()
	dialog.ShowError(fmt.Errorf("Não foi possível inicializar: %w", err), s.Window)
	LogApplication("startup_failed", map[string]any{"phase": "splash"})
}

// maybeTransition aplica fade out/in somente depois de ambas as condições obrigatórias do domínio.
func (s *Splash) maybeTransition() {
	if !s.gate.BeginTransition() {
		return
	}
	if s.next == nil {
		s.fail(fmt.Errorf("Tela inicial não configurada."))
		return
	}
	LogApplication("splash_transition_start", nil)
	go func() {
		if !s.fade(s.Window, false) {
			return
		}
		// next só pode tocar widgets na thread principal; o player é fechado antes da troca de conteúdo.
		var destination fyne.Window
		fyne.DoAndWait(func() {
			if s.gate.Closed {
				return
			}
			if s.player != nil {
				s.player.Close()
				s.player = nil
			}
			// A janela normal recebe o layout de boas-vindas durante a parte invisível do fade.
			s.Window.SetOnClosed(nil)
			destination = s.next(s.Window)
			if destination != nil {
				setSplashOpacity(destination, 0)
			}
		})
		if destination == nil {
			return
		}
		if !s.fade(destination, true) {
			return
		}
		fyne.Do(func() {
			s.cancel()
			LogApplication("startup_ready", nil)
			if s.finished != nil {
				s.finished()
			}
		})
	}()
}

// fade modifica somente opacidade nativa; trabalho do domínio já terminou e nenhuma mídia é cortada.
func (s *Splash) fade(window fyne.Window, reveal bool) bool {
	// steps distribui 220 ms em quadros curtos, sem bloquear a thread gráfica.
	const steps = 14
	for step := 0; step <= steps; step++ {
		select {
		case <-s.ctx.Done():
			return false
		case <-time.After(220 * time.Millisecond / steps):
		}
		alpha := uint8(255 * step / steps)
		if !reveal {
			alpha = 255 - alpha
		}
		fyne.DoAndWait(func() {
			if !s.gate.Closed {
				setSplashOpacity(window, alpha)
			}
		})
	}
	return true
}

// close libera a reprodução na thread da janela e impede qualquer liberação posterior.
func (s *Splash) close() {
	if !s.gate.Closed && s.gate.Err == nil {
		LogApplication("startup_cancelled", nil)
	}
	s.gate.Close()
	s.cancel()
	if s.player != nil {
		s.player.Close()
		s.player = nil
	}
}
