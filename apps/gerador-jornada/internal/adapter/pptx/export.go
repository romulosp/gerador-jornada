package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/scenario"
	"gerador-jornada.local/editor/internal/screen"
)

// Options seleciona percursos sem omitir roteiros inválidos silenciosamente.
type Options struct {
	// Exploration inclui o percurso manual.
	Exploration bool
	// ScenarioIDs seleciona roteiros; nil seleciona todos os não vazios.
	ScenarioIDs []string
}

// menuLink vincula uma entrada de menu ao número de um slide.
type menuLink struct {
	// Caption identifica o percurso.
	Caption string
	// Target informa o índice de destino, começando em um.
	Target int
}

// slidePlan descreve uma página antes da criação de XML.
type slidePlan struct {
	// Page referencia uma composição de tela, quando aplicável.
	Page *screen.Screen
	// Visit referencia uma ocorrência da demonstração.
	Visit *scenario.Occurrence
	// Links contém entradas de menu ou de encerramento.
	Links []menuLink
	// Caption identifica o slide no arquivo.
	Caption string
	// Next determina o avanço automático para a página seguinte.
	Next bool
}

// deck mantém somente estado local de uma exportação.
type deck struct {
	// Document é o snapshot de origem.
	Document application.Document
	// Width é a maior largura lógica das telas.
	Width int
	// Height inclui a faixa inferior de navegação.
	Height int
	// Slides contém o plano completo e determinístico.
	Slides []slidePlan
	// Exploration associa IDs de tela aos seus slides manuais.
	Exploration map[string]int
}

// Export valida e publica o arquivo completo; cancelamento preserva o destino anterior.
func Export(ctx context.Context, path string, document application.Document, options Options, progress func(float64)) error {
	// data contém o pacote validado em memória.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	data, err := Build(ctx, document, options, progress)
	if err != nil {
		return err
	}
	return atomicfile.Write(ctx, path, data)
}

// Build gera o pacote determinístico sem exigir PowerPoint instalado.
func Build(ctx context.Context, document application.Document, options Options, progress func(float64)) ([]byte, error) {
	// problems reúne os diagnósticos que impedem execução ou exportação.
	if problems := document.Validate(); len(problems) > 0 {
		return nil, fmt.Errorf("%s", application.ProblemsText(problems))
	}
	// plan resolve todos os índices antes de escrever relacionamentos.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	plan, err := planDeck(document, options)
	if err != nil {
		return nil, err
	}
	// parts contém partes OPC ordenadas na publicação.
	parts := plan.baseParts()
	// index determina o número estável de cada slide.
	for index, slide := range plan.Slides {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		// content e links pertencem exclusivamente ao slide atual.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		content, links, err := plan.renderSlide(slide)
		if err != nil {
			return nil, err
		}
		parts[fmt.Sprintf("ppt/slides/slide%d.xml", index+1)] = content
		parts[fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", index+1)] = relations(links)
		if progress != nil {
			progress(float64(index+1) / float64(len(plan.Slides)) * .9)
		}
	}
	// asset é um recurso decorativo já validado pela persistência.
	for _, asset := range document.Project.Assets {
		parts["ppt/media/"+asset.ID+extension(asset.Path)] = string(document.Resources[asset.Path])
	}
	// names define a ordem estável das entradas ZIP.
	names := []string{}
	// name identifica uma entrada relativa do pacote ou arquivo gerado.
	for name := range parts {
		names = append(names, name)
	}
	sort.Strings(names)
	// output armazena o pacote final.
	output := &bytes.Buffer{}
	// writer normaliza os cabeçalhos e o horário das entradas.
	writer := zip.NewWriter(output)
	// name identifica uma parte independente do pacote.
	for _, name := range names {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		// entry recebe XML ou mídia sem caminhos externos.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		entry, err := writer.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
		if err != nil {
			return nil, err
		}
		if _, err = entry.Write([]byte(parts[name])); err != nil {
			return nil, err
		}
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	if err = ValidatePackage(output.Bytes()); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// planDeck separa exploração livre de demonstrações temporizadas.
func planDeck(document application.Document, options Options) (*deck, error) {
	// plan começa com o menu, sem destinos ainda resolvidos.
	plan := &deck{Document: document, Slides: []slidePlan{{Caption: document.Project.Name}}, Exploration: map[string]int{}}
	// page determina o tamanho único do arquivo.
	for _, page := range document.Screens {
		plan.Width = max(plan.Width, page.Size.Width)
		plan.Height = max(plan.Height, page.Size.Height)
	}
	plan.Height += 40
	if options.Exploration {
		// index localiza a tela no snapshot, preservando sua ordem.
		for index := range plan.Document.Screens {
			// page resolve a composição corrente sem confundir sua identidade com o nome exibido.
			page := &plan.Document.Screens[index]
			plan.Exploration[page.ID] = len(plan.Slides) + 1
			plan.Slides = append(plan.Slides, slidePlan{Page: page, Caption: page.Name})
		}
		plan.Slides[0].Links = append(plan.Slides[0].Links, menuLink{"Explorar a jornada", plan.Exploration[document.Journey.StartScreenID]})
	}
	// selected mantém a seleção explícita ou o padrão de todos os roteiros.
	selected := map[string]bool{}
	if options.ScenarioIDs != nil {
		// id identifica de forma estável o item ou recurso atual.
		for _, id := range options.ScenarioIDs {
			selected[id] = true
		}
	}
	// script é validado por inteiro antes de adicionar seus slides.
	for _, script := range document.Scenarios {
		if (options.ScenarioIDs == nil && len(script.Steps) == 0) || (options.ScenarioIDs != nil && !selected[script.ID]) {
			continue
		}
		delete(selected, script.ID)
		// visits representam inclusive as revisitas à mesma tela.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		visits, err := script.Interpret(document.Screens, document.Journey)
		if err != nil {
			return nil, err
		}
		// start identifica a primeira ocorrência do roteiro.
		start := len(plan.Slides) + 1
		plan.Slides[0].Links = append(plan.Slides[0].Links, menuLink{"Demonstração: " + script.Name, start})
		// index referencia cada visita por ponteiro estável.
		for index := range visits {
			plan.Slides = append(plan.Slides, slidePlan{Page: plan.Document.Screen(visits[index].ScreenID), Visit: &visits[index], Caption: script.Name, Next: true})
		}
		plan.Slides = append(plan.Slides, slidePlan{Caption: "Demonstração concluída", Links: []menuLink{{"Repetir demonstração", start}, {"Voltar ao menu", 1}}})
	}
	if len(selected) > 0 {
		return nil, fmt.Errorf("SCENARIO_EMPTY: seleção de roteiro inexistente")
	}
	if len(plan.Slides[0].Links) == 0 {
		return nil, fmt.Errorf("PPT_TARGET_INVALID: selecione um percurso")
	}
	return plan, nil
}

// extension preserva apenas formatos decorativos aceitos.
func extension(path string) string {
	if strings.HasSuffix(path, ".jpg") {
		return ".jpg"
	}
	return ".png"
}

// renderSlide converte componentes e roteiro em XML editável.
func (d *deck) renderSlide(plan slidePlan) (string, string, error) {
	// render controla IDs de formas e vínculos por slide.
	render := newSlideRenderer(d, plan)
	render.content.WriteString(treeStart)
	if plan.Page == nil {
		render.menu()
	} else {
		render.page()
	}
	render.content.WriteString("</p:spTree>")
	// transition fica separado da árvore de animações.
	transition := ""
	if plan.Visit != nil {
		// duration expressa em milissegundos o tempo da ocorrência exportada.
		duration := plan.Visit.DurationMS
		if len(plan.Visit.Events) == 0 || plan.Visit.Events[len(plan.Visit.Events)-1].Step.Kind != "navigate" {
			duration += 1000
		}
		transition = fmt.Sprintf(`<p:transition advClick="0" advTm="%d"/>`, max(1, duration))
	}
	// timing contém somente animações nativas do roteiro selecionado.
	timing := render.timingXML()
	return `<?xml version="1.0" encoding="UTF-8"?><p:sld ` + slideNamespaces + `><p:cSld name="` + escape(plan.Caption) + `"><p:bg><p:bgPr>` + solid("#FFFFFF") + `<a:effectLst/></p:bgPr></p:bg>` + render.content.String() + `</p:cSld><p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>` + transition + timing + `</p:sld>`, render.links.String(), nil
}

// slideRenderer resolve formas, vínculos e efeitos de uma página.
type slideRenderer struct {
	// deck fornece recursos e dimensões do arquivo.
	deck *deck
	// plan descreve a visita ou tela exploratória.
	plan slidePlan
	// content acumula formas do slide.
	content strings.Builder
	// links acumula relacionamentos locais.
	links strings.Builder
	// effects acumula nós temporais do roteiro.
	effects strings.Builder
	// nextID garante unicidade de formas.
	nextID int
	// timeID garante unicidade de nós temporais.
	timeID int
	// values associa campos às suas formas de valor.
	values map[string]int
	// valueShapes conserva geometria e estilo dos campos.
	valueShapes map[string]presentation.Primitive
	// offsetX centraliza telas menores.
	offsetX int
	// offsetY centraliza telas menores acima da faixa de navegação.
	offsetY int
}

// newSlideRenderer inicia contadores e o vínculo de layout obrigatório.
func newSlideRenderer(d *deck, plan slidePlan) *slideRenderer {
	// renderer armazena somente estado efêmero.
	renderer := &slideRenderer{deck: d, plan: plan, nextID: 2, timeID: 3, values: map[string]int{}, valueShapes: map[string]presentation.Primitive{}}
	renderer.links.WriteString(relation("rIdLayout", "slideLayout", "../slideLayouts/slideLayout1.xml"))
	if plan.Page != nil {
		renderer.offsetX = (d.Width - plan.Page.Size.Width) / 2
		renderer.offsetY = (d.Height - 40 - plan.Page.Size.Height) / 2
	}
	return renderer
}

// id reserva um identificador de forma.
func (r *slideRenderer) id() int { r.nextID++; return r.nextID }

// link registra um destino interno específico do slide.
func (r *slideRenderer) link(target int) string {
	// id também é único quando dois botões levam ao mesmo destino.
	id := fmt.Sprintf("rIdLink%d", r.id())
	r.links.WriteString(relation(id, "slide", fmt.Sprintf("slide%d.xml", target)))
	return id
}

// menu desenha as entradas de exploração e roteiro.
func (r *slideRenderer) menu() {
	r.content.WriteString(shape(r.id(), "Título", presentation.Primitive{Bounds: screen.Bounds{X: 48, Y: 48, Width: r.deck.Width - 96, Height: 60}, Text: r.plan.Caption, Color: "#0F172A", FontSize: 28, Bold: true}, "", ""))
	r.content.WriteString(shape(r.id(), "Descrição", presentation.Primitive{Bounds: screen.Bounds{X: 48, Y: 112, Width: r.deck.Width - 96, Height: 32}, Text: "Escolha um percurso para conhecer a jornada.", Color: "#64748B", FontSize: 15}, "", ""))
	// index define o espaçamento entre ações.
	// item representa um roteiro examinado na filtragem.
	for index, item := range r.plan.Links {
		r.content.WriteString(shape(r.id(), item.Caption, presentation.Primitive{Bounds: screen.Bounds{X: 48, Y: 180 + index*64, Width: r.deck.Width - 96, Height: 48}, Text: item.Caption, Fill: "#EFF6FF", Color: "#1D4ED8", FontSize: 18, Bold: true}, r.link(item.Target), ""))
	}
}

// page expande a composição e mantém controles como grupos editáveis.
func (r *slideRenderer) page() {
	// page aponta para a composição imutável.
	page := r.plan.Page
	r.content.WriteString(shape(r.id(), "Fundo", presentation.Primitive{Bounds: screen.Bounds{X: r.offsetX, Y: r.offsetY, Width: page.Size.Width, Height: page.Size.Height}, Fill: page.Background.Color}, "", ""))
	if page.Background.AssetID != "" {
		r.image(presentation.Primitive{Bounds: screen.Bounds{X: r.offsetX, Y: r.offsetY, Width: page.Size.Width, Height: page.Size.Height}, AssetID: page.Background.AssetID})
	}
	// group é um componente expandido sem dependência da janela.
	for _, group := range presentation.Build(*page) {
		group.Bounds.X += r.offsetX
		group.Bounds.Y += r.offsetY
		// link é aplicado somente à exploração.
		link := ""
		if r.plan.Visit == nil {
			// target resolve o destino da operação na coleção ou no mapa corrente.
			if target := r.deck.Exploration[r.deck.Document.Journey.Destination(group.ComponentID)]; target > 0 {
				link = r.link(target)
			}
		}
		r.content.WriteString(groupStart(r.id(), group, link))
		// primitive é uma forma ou texto interno do controle.
		for _, primitive := range group.Shapes {
			primitive.Bounds.X += r.offsetX
			primitive.Bounds.Y += r.offsetY
			if primitive.AssetID != "" {
				r.image(primitive)
				continue
			}
			// id identifica inclusive a caixa vazia de valor do campo.
			id := r.id()
			r.content.WriteString(shape(id, group.ComponentID, primitive, "", ""))
			if primitive.Role == "value" {
				r.values[group.ComponentID] = id
				r.valueShapes[group.ComponentID] = primitive
			}
		}
		r.content.WriteString("</p:grpSp>")
	}
	if r.plan.Visit == nil {
		r.content.WriteString(shape(r.id(), "Menu", presentation.Primitive{Bounds: screen.Bounds{X: 12, Y: r.deck.Height - 30, Width: 140, Height: 24}, Text: "← Menu da jornada", Color: "#2563EB", FontSize: 12}, r.link(1), ""))
	} else {
		r.animate()
	}
}

// image incorpora apenas recursos decorativos já declarados.
func (r *slideRenderer) image(primitive presentation.Primitive) {
	// asset procura a imagem pelo identificador estável.
	for _, asset := range r.deck.Document.Project.Assets {
		if asset.ID != primitive.AssetID {
			continue
		}
		// id mantém forma e relacionamento únicos.
		id := r.id()
		// rid identifica a relação entre a forma e o recurso de imagem no pacote PPTX.
		rid := fmt.Sprintf("rIdImage%d", id)
		r.links.WriteString(relation(rid, "image", "../media/"+asset.ID+extension(asset.Path)))
		r.content.WriteString(fmt.Sprintf(`<p:pic><p:nvPicPr><p:cNvPr id="%d" name="%s"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr><p:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></p:blipFill><p:spPr>%s<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:pic>`, id, escape(asset.ID), rid, transform(primitive.Bounds)))
	}
}
