package screen

import "fmt"

// PlaceComponent reposiciona somente o componente na lista plana fundo→frente.
// referenceID vazio representa o extremo absoluto; inFront coloca imediatamente à frente.
// A operação não move subárvores nem altera geometria, estilo, vínculos ou ParentID.
func PlaceComponent(page *Screen, componentID, referenceID string, inFront bool) error {
	if page == nil {
		return fmt.Errorf("tela inexistente")
	}
	// sourceIndex localiza a origem na ordem canônica.
	sourceIndex := -1
	// index percorre a lista plana.
	for index := range page.Components {
		if page.Components[index].ID == componentID {
			sourceIndex = index
			break
		}
	}
	if sourceIndex < 0 {
		return fmt.Errorf("componente de origem inexistente")
	}
	if referenceID != "" && referenceID == componentID {
		return fmt.Errorf("componente não pode ser referência de si mesmo")
	}
	// source preserva o valor integral do componente.
	source := page.Components[sourceIndex]
	// remaining remove a origem sem alterar a ordem relativa dos demais.
	remaining := append([]Component{}, page.Components[:sourceIndex]...)
	remaining = append(remaining, page.Components[sourceIndex+1:]...)
	// target define a posição de inserção no restante.
	target := 0
	if inFront {
		target = len(remaining)
	}
	// foundReference informa se a referência foi encontrada.
	foundReference := referenceID == ""
	if referenceID != "" {
		target = -1
		// index percorre o restante para localizar a referência.
		for index := range remaining {
			if remaining[index].ID == referenceID {
				target = index
				foundReference = true
				break
			}
		}
		if !foundReference {
			return fmt.Errorf("componente de referência inexistente")
		}
		// inFront insere imediatamente após a referência; caso contrário, imediatamente antes.
		if inFront {
			target++
		}
	}
	// result insere a origem na posição calculada.
	result := []Component{}
	result = append(result, remaining[:target]...)
	result = append(result, source)
	result = append(result, remaining[target:]...)
	page.Components = result
	return nil
}

// IsPlacedAfter informa se o componente já está imediatamente à frente da referência.
// referenceID vazio verifica o extremo frontal.
func IsPlacedAfter(page Screen, componentID, referenceID string) bool {
	if page.Components == nil {
		return false
	}
	// index percorre a ordem canônica.
	for index := range page.Components {
		if page.Components[index].ID != componentID {
			continue
		}
		if referenceID == "" {
			return index == len(page.Components)-1
		}
		return index > 0 && page.Components[index-1].ID == referenceID
	}
	return false
}
