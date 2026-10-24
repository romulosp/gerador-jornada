package application

import (
	"errors"
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
)

// ErrJourneyMissing distingue destino removido de arquivo corrompido ou versão incompatível.
var ErrJourneyMissing = errors.New("jornada destino ausente")

// EnterJourney carrega e valida antes de substituir a prévia, preservando a origem em caso de falha.
func (p *Preview) EnterJourney(filename string) error {
	if !p.AllowScreenNavigation {
		return fmt.Errorf("Navegação entre jornadas desativada na Prévia Tela Atual.")
	}
	if !catalog.ValidJourneyFilename(filename) || p.LoadJourney == nil {
		return fmt.Errorf("Destino de jornada sem workspace válido: %s", filename)
	}
	if len(p.JourneyHistory) >= 32 {
		return fmt.Errorf("Limite de 32 jornadas visitadas; use Voltar.")
	}
	// document e next permanecem independentes até a validação integral do destino.
	document, err := p.LoadJourney(filename)
	if err != nil {
		return fmt.Errorf("Abrir jornada %s: %w", filename, err)
	}
	next, err := NewPreview(document)
	if err != nil {
		return err
	}
	// previous retém também tela, cursor e caminho local do menu, sem persistir nada no arquivo.
	previous := *p
	previous.JourneyHistory = nil
	next.JourneyHistory = append(append([]*Preview{}, p.JourneyHistory...), &previous)
	next.LoadJourney, next.OnMissingJourney = p.LoadJourney, p.OnMissingJourney
	next.CurrentJourney = filename
	*p = *next
	return nil
}

// BackJourney restaura a jornada anterior e a tela em que o autor deixou a prévia.
func (p *Preview) BackJourney() error {
	if len(p.JourneyHistory) == 0 {
		return fmt.Errorf("Não há jornada anterior.")
	}
	// history e previous removem somente o topo da pilha, mantendo o runtime da origem.
	history := p.JourneyHistory
	previous := history[len(history)-1]
	*p = *previous
	p.JourneyHistory = history[:len(history)-1]
	return nil
}
