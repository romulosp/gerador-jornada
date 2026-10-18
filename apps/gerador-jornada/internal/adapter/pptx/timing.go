package pptx

import (
	"fmt"

	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/scenario"
	"gerador-jornada.local/editor/internal/screen"
)

// behavior cria o alvo e a condição temporal de um efeito nativo.
func (r *slideRenderer) behavior(shapeID, start, duration int, attribute string) string {
	r.timeID++
	return fmt.Sprintf(`<p:cBhvr><p:cTn id="%d" dur="%d" fill="hold"><p:stCondLst><p:cond delay="%d"/></p:stCondLst></p:cTn><p:tgtEl><p:spTgt spid="%d"/></p:tgtEl><p:attrNameLst><p:attrName>%s</p:attrName></p:attrNameLst></p:cBhvr>`, r.timeID, max(1, duration), max(0, start), shapeID, attribute)
}

// visibility altera somente a visibilidade da forma mantendo texto editável.
func (r *slideRenderer) visibility(id, start int, visible bool) {
	// value é o literal reconhecido pelo PresentationML.
	value := "hidden"
	if visible {
		value = "visible"
	}
	r.effects.WriteString(`<p:set>` + r.behavior(id, start, 1, "style.visibility") + `<p:to><p:strVal val="` + value + `"/></p:to></p:set>`)
}

// motion usa uma trajetória nativa com coordenadas normalizadas ao slide.
func (r *slideRenderer) motion(id, start, duration int, from, to scenario.Point) {
	r.effects.WriteString(fmt.Sprintf(`<p:animMotion origin="layout" path="M %.7f %.7f L %.7f %.7f E" pathEditMode="relative">`, from.X/float64(r.deck.Width), from.Y/float64(r.deck.Height), to.X/float64(r.deck.Width), to.Y/float64(r.deck.Height)))
	r.effects.WriteString(r.behavior(id, start, duration, "ppt_x"))
	r.effects.WriteString("</p:animMotion>")
}

// animate converte eventos interpretados em cursor, realces e caixas de texto.
func (r *slideRenderer) animate() {
	// cursorID identifica a seta vetorial editável.
	cursorID := r.id()
	r.content.WriteString(shape(cursorID, "Cursor", presentation.Primitive{Bounds: screen.Bounds{X: r.offsetX, Y: r.offsetY, Width: 12, Height: 18}, Fill: "#2563EB", Stroke: "#FFFFFF"}, "", "rtTriangle"))
	// event é uma etapa com estado resolvido pelo domínio.
	for _, event := range r.plan.Visit.Events {
		switch event.Step.Kind {
		case "moveCursor":
			r.motion(cursorID, event.StartMS, event.Step.DurationMS, event.From, event.To)
		case "click":
			r.motion(cursorID, event.StartMS, 1, event.From, event.To)
			// markerID identifica o realce de clique desta etapa.
			markerID := r.id()
			r.content.WriteString(shape(markerID, "Clique", presentation.Primitive{Bounds: screen.Bounds{X: r.offsetX + int(event.To.X) - 10, Y: r.offsetY + int(event.To.Y) - 10, Width: 28, Height: 28}, Stroke: "#F59E0B"}, "", "ellipse"))
			r.visibility(markerID, 0, false)
			r.visibility(markerID, event.StartMS+1, true)
			r.visibility(markerID, event.StartMS+event.Step.DurationMS, false)
		case "typeText":
			r.typing(event)
		}
	}
}

// typing revela prefixos Unicode sem rasterizar o valor do campo.
func (r *slideRenderer) typing(event scenario.Event) {
	// primitive conserva a geometria e tipografia do valor.
	// ok informa se a consulta encontrou um valor válido.
	primitive, ok := r.valueShapes[event.Step.ComponentID]
	if !ok {
		return
	}
	// previousID é a última caixa visível desse campo.
	previousID := r.values[event.Step.ComponentID]
	// clusters mantém os caracteres visíveis inteiros.
	clusters := scenario.Graphemes(*event.Step.Text)
	if len(clusters) == 0 {
		clusters = []string{""}
	}
	// prefix acumula somente caracteres completos.
	prefix := ""
	// index determina o instante de troca de prefixo.
	// cluster preserva um grafema completo durante a reprodução da digitação.
	for index, cluster := range clusters {
		prefix += cluster
		primitive.Text = prefix
		// id identifica o texto editável progressivo.
		id := r.id()
		r.content.WriteString(shape(id, "Digitação", primitive, "", ""))
		// start evita visibilidade inicial acidental.
		start := event.StartMS + (index+1)*event.Step.DurationMS/len(clusters)
		r.visibility(id, 0, false)
		r.visibility(previousID, start, false)
		r.visibility(id, start, true)
		previousID = id
	}
	r.values[event.Step.ComponentID] = previousID
}

// timingXML envolve os efeitos no nó temporal raiz do slide.
func (r *slideRenderer) timingXML() string {
	if r.effects.Len() == 0 {
		return ""
	}
	return `<p:timing><p:tnLst><p:par><p:cTn id="1" dur="indefinite" restart="never" nodeType="tmRoot"><p:childTnLst><p:par><p:cTn id="2" dur="indefinite" fill="hold"><p:stCondLst><p:cond delay="0"/></p:stCondLst><p:childTnLst>` + r.effects.String() + `</p:childTnLst></p:cTn></p:par></p:childTnLst></p:cTn></p:par></p:tnLst></p:timing>`
}
