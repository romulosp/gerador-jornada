package application

import (
	"gerador-jornada.local/editor/internal/scenario"
	"gerador-jornada.local/editor/internal/screen"
)

// Frame descreve o estado visual derivado em um instante do roteiro.
type Frame struct {
	// Page é uma cópia da tela com valores runtime aplicados.
	Page screen.Screen
	// Cursor localiza a seta na composição.
	Cursor scenario.Point
	// Click indica se o realce de clique está ativo.
	Click bool
	// Finished indica o fim do roteiro e da permanência final.
	Finished bool
}

// FrameAt calcula um estado reproduzível sem relógio, widget ou efeitos colaterais.
func FrameAt(document Document, visits []scenario.Occurrence, elapsedMS int) Frame {
	// frame é vazio somente para uma lista não interpretada.
	frame := Frame{}
	// selected determina a visita atual, inclusive revisitas.
	selected := 0
	// local é o tempo relativo à ocorrência selecionada.
	local := elapsedMS
	// index percorre as durações acumuladas.
	// visit localiza a ocorrência atual na sequência da demonstração.
	for index, visit := range visits {
		// duration inclui permanência apenas na última ocorrência.
		duration := visit.DurationMS
		if index == len(visits)-1 {
			duration += 1000
		}
		selected = index
		if local < duration {
			break
		}
		if index == len(visits)-1 {
			frame.Finished = true
			break
		}
		local -= duration
	}
	if len(visits) == 0 {
		return frame
	}
	// visit mantém eventos e valores desta visita somente.
	visit := visits[selected]
	// copy impede qualquer mutação dos campos persistidos.
	copy := document.Copy()
	// page resolve a composição da visita.
	page := copy.Screen(visit.ScreenID)
	if page == nil {
		return frame
	}
	frame.Page = *page
	// event aplica alterações concluídas e a interpolação corrente.
	for _, event := range visit.Events {
		if local < event.StartMS {
			break
		}
		// duration limita a interpolação ao intervalo da etapa.
		duration := min(local-event.StartMS, event.Step.DurationMS)
		// progress é uma fração entre zero e um.
		progress := float64(duration) / float64(event.Step.DurationMS)
		switch event.Step.Kind {
		case "moveCursor":
			frame.Cursor = scenario.Point{X: event.From.X + (event.To.X-event.From.X)*progress, Y: event.From.Y + (event.To.Y-event.From.Y)*progress}
		case "click":
			frame.Cursor = event.To
			frame.Click = duration < event.Step.DurationMS
		case "typeText":
			// component representa o campo do estado visual copiado.
			component := frame.Page.ComponentByID(event.Step.ComponentID)
			if component == nil {
				continue
			}
			// clusters mantém unidades visíveis Unicode.
			clusters := scenario.Graphemes(*event.Step.Text)
			// text é o prefixo visível deste instante.
			text := ""
			// cluster preserva um grafema completo durante a reprodução da digitação.
			for _, cluster := range clusters[:min(len(clusters), int(float64(len(clusters))*progress))] {
				text += cluster
			}
			if duration == 0 {
				text = event.PreviousText
			}
			component.Properties.Field.InitialValue = text
		}
	}
	return frame
}
