package pptx

import "fmt"

// transitionXML combina efeito nativo e avanço, sem criar temporizador no percurso exploratório.
func transitionXML(effect string, advanceMS int) string {
	// child determina somente o efeito de entrada; timing do roteiro continua separado.
	child := `<p:fade/>`
	switch effect {
	case "push":
		child = `<p:push dir="l"/>`
	case "wipe":
		child = `<p:wipe dir="r"/>`
	case "none":
		child = ""
	}
	// advance só existe para ocorrências temporizadas da demonstração.
	advance := ""
	if advanceMS > 0 {
		advance = fmt.Sprintf(` advTm="%d"`, advanceMS)
	}
	if effect == "none" {
		return `<p:transition advClick="0"` + advance + `/>`
	}
	return `<p:transition spd="fast" p14:dur="300" advClick="0"` + advance + `>` + child + `</p:transition>`
}
