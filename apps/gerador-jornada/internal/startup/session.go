// Package startup define o ciclo de inicialização, sem UI, relógio, áudio ou acesso a arquivos.
package startup

// Session impede liberar a aplicação antes do carregamento e do término real da reprodução.
// Seu coordenador deve serializar eventos; o domínio não cria goroutines nem depende de plataforma.
type Session struct {
	// Loaded indica que todas as etapas de preparação obrigatórias terminaram com sucesso.
	Loaded bool
	// Playing e VideoEnded representam os eventos reais do reprodutor, não duração estimada.
	Playing, VideoEnded bool
	// Err bloqueia a liberação quando carregamento ou reprodução falham.
	Err error
	// Closed cancela o ciclo e Transitioning garante uma única passagem de tela.
	Closed, Transitioning bool
}

// ResourcesCompleted confirma a preparação ou conserva a falha para apresentação pelo adaptador.
func (s *Session) ResourcesCompleted(err error) {
	if err != nil {
		s.Fail(err)
	} else {
		s.Loaded = true
	}
}

// PlaybackStarted confirma que o vídeo realmente começou.
func (s *Session) PlaybackStarted() { s.Playing = true }

// PlaybackCompleted confirma o fim apenas de uma reprodução iniciada.
func (s *Session) PlaybackCompleted() {
	if s.Playing {
		s.VideoEnded = true
	}
}

// Fail conserva a primeira falha e nunca a converte em término bem-sucedido.
func (s *Session) Fail(err error) {
	if s.Err == nil {
		s.Err = err
	}
}

// BeginTransition libera exatamente uma vez depois de ambas as condições obrigatórias.
func (s *Session) BeginTransition() bool {
	if s.Closed || s.Transitioning || s.Err != nil || !s.Loaded || !s.Playing || !s.VideoEnded {
		return false
	}
	s.Transitioning = true
	return true
}

// Close impede callbacks tardios de avançar a inicialização cancelada.
func (s *Session) Close() { s.Closed = true }
