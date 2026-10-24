package screen

import "fmt"

// SelectionBounds calcula a caixa conjunta sem armazenar seleção dentro da tela.
func SelectionBounds(page Screen, ids []string) (Bounds, error) {
	// combined acumula as extremidades dos componentes selecionados.
	combined := Bounds{}
	if len(ids) == 0 {
		return combined, fmt.Errorf("nenhum componente selecionado")
	}
	// index e id percorrem somente os alvos explicitamente solicitados.
	for index, id := range ids {
		// component precisa pertencer à composição consultada.
		component := page.ComponentByID(id)
		if component == nil {
			return Bounds{}, fmt.Errorf("componente selecionado inexistente: %s", id)
		}
		if index == 0 {
			combined = component.Bounds
			continue
		}
		// right e bottom preservam as extremidades antes de mover a origem da caixa conjunta.
		right, bottom := max(combined.X+combined.Width, component.Bounds.X+component.Bounds.Width), max(combined.Y+combined.Height, component.Bounds.Y+component.Bounds.Height)
		combined.X, combined.Y = min(combined.X, component.Bounds.X), min(combined.Y, component.Bounds.Y)
		combined.Width, combined.Height = right-combined.X, bottom-combined.Y
	}
	return combined, nil
}

// ConstrainTranslation limita o deslocamento do conjunto inteiro, preservando distâncias relativas.
func ConstrainTranslation(bounds Bounds, size Size, dx, dy int) (int, int) {
	return min(max(dx, -bounds.X), size.Width-bounds.X-bounds.Width), min(max(dy, -bounds.Y), size.Height-bounds.Y-bounds.Height)
}

// MoveSelection traduz alvos válidos com um único deslocamento comum, sem alterar outros controles.
// Selecionar um pai desloca sua subárvore uma única vez, mesmo que um filho também esteja selecionado.
func MoveSelection(page *Screen, ids []string, dx, dy int) error {
	// expanded reúne os alvos e seus descendentes sem repetição.
	expanded := []string{}
	// seen evita mover um alvo duas vezes quando o consumidor repetir sua identidade.
	seen := map[string]bool{}
	for _, id := range ids {
		if page.ComponentByID(id) == nil {
			return fmt.Errorf("componente selecionado inexistente: %s", id)
		}
		// index percorre o alvo e a subárvore.
		for _, moveID := range append([]string{id}, page.Descendants(id)...) {
			if seen[moveID] {
				continue
			}
			seen[moveID] = true
			expanded = append(expanded, moveID)
		}
	}
	// bounds valida todas as referências antes de publicar qualquer geometria.
	bounds, err := SelectionBounds(*page, expanded)
	if err != nil {
		return err
	}
	dx, dy = ConstrainTranslation(bounds, page.Size, dx, dy)
	for _, id := range expanded {
		// component já foi validado pela caixa conjunta.
		component := page.ComponentByID(id)
		component.Bounds.X += dx
		component.Bounds.Y += dy
	}
	return nil
}

// IntersectingIDs retorna componentes atingidos por uma caixa de seleção em ordem de camadas.
func IntersectingIDs(page Screen, area Bounds) []string {
	// ids não persiste no documento; é apenas resultado da consulta geométrica.
	ids := []string{}
	for _, component := range page.Components {
		// bounds fornece o retângulo do candidato.
		bounds := component.Bounds
		if bounds.X < area.X+area.Width && bounds.X+bounds.Width > area.X && bounds.Y < area.Y+area.Height && bounds.Y+bounds.Height > area.Y {
			ids = append(ids, component.ID)
		}
	}
	return ids
}
