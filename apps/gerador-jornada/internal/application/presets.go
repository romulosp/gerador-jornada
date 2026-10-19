package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// Preset descreve um conjunto editável formado pelas famílias existentes do catálogo.
type Preset struct {
	// ID identifica a composição pronta sem criar uma nova família serializada.
	ID string
	// Name informa o nome amigável na paleta.
	Name string
	// Width expressa a largura mínima do conjunto em unidades lógicas.
	Width int
	// Height expressa a altura mínima do conjunto em unidades lógicas.
	Height int
}

// Presets relaciona os elementos auxiliares observados no Design do legado.
func Presets() []Preset {
	return []Preset{{"dateInfo", "Data e hora", 126, 32}, {"lotericaInfo", "Dados da lotérica", 222, 36},
		{"infraInfo", "Infraestrutura", 78, 36}, {"header", "Cabeçalho completo", 426, 40},
		{"pageTitle", "Faixa de título da página", 800, 28}}
}

// presetLabel define uma linha de texto relativa à origem de um conjunto pronto.
type presetLabel struct {
	// Name fornece o nome lógico inicial da linha.
	Name string
	// Text contém um valor ilustrativo fixo e editável, sem consulta à máquina ou à rede.
	Text string
	// Bounds expressa a posição relativa em unidades lógicas.
	Bounds screen.Bounds
}

// presetLabels expande os três grupos informativos e o cabeçalho composto.
func presetLabels(id string) []presetLabel {
	switch id {
	case "dateInfo":
		return []presetLabel{{"dataBrasilia", "BSB - 01/01/2026", screen.Bounds{X: 0, Y: 3, Width: 126, Height: 12}},
			{"horaBrasilia", "12:00:00", screen.Bounds{X: 0, Y: 19, Width: 126, Height: 12}}}
	case "lotericaInfo":
		return []presetLabel{{"proprietario", "PROPRIETÁRIO: L009666 - USUÁRIO", screen.Bounds{Width: 222, Height: 12}},
			{"terminal", "Terminal: 000599", screen.Bounds{Y: 12, Width: 222, Height: 12}},
			{"loterica", "LOTÉRICA: 009999", screen.Bounds{Y: 24, Width: 222, Height: 12}}}
	case "infraInfo":
		return []presetLabel{{"gateway", "GATEWAY", screen.Bounds{Width: 78, Height: 12}},
			{"dns", "DNS", screen.Bounds{Y: 12, Width: 78, Height: 12}},
			{"http", "HTTP", screen.Bounds{Y: 24, Width: 78, Height: 12}}}
	case "header":
		// labels preserva o alinhamento das três áreas do cabeçalho original.
		labels := presetLabels("dateInfo")
		// group posiciona os grupos de usuário e infraestrutura após a data.
		for _, group := range []struct {
			// ID seleciona um conjunto informativo já definido.
			ID string
			// X desloca as linhas horizontalmente em unidades lógicas.
			X int
		}{{"lotericaInfo", 126}, {"infraInfo", 348}} {
			// line recebe somente uma translação, sem alterar o conjunto de origem.
			for _, line := range presetLabels(group.ID) {
				line.Bounds.X += group.X
				labels = append(labels, line)
			}
		}
		return labels
	}
	return nil
}

// InsertPreset insere componentes comuns; o chamador controla a transação e o histórico.
// Data, hora e indicadores são conteúdo ilustrativo, sem relógio vivo ou monitoramento de rede.
func InsertPreset(document *Document, pageID, presetID string, x, y int) ([]string, error) {
	// page deve existir antes de construir qualquer componente.
	page := document.Screen(pageID)
	if page == nil {
		return nil, fmt.Errorf("tela inexistente")
	}
	// definition confirma identidade e dimensões antes da alteração.
	var definition *Preset
	// candidate percorre os conjuntos conhecidos, sem aceitar templates externos.
	for _, candidate := range Presets() {
		if candidate.ID == presetID {
			definition = &candidate
			break
		}
	}
	if definition == nil {
		return nil, fmt.Errorf("conjunto de componentes desconhecido")
	}
	if presetID == "pageTitle" {
		definition.Width = page.Size.Width
	}
	if page.Size.Width < definition.Width || page.Size.Height < definition.Height {
		return nil, fmt.Errorf("a tela precisa comportar o conjunto %s (%d × %d)", definition.Name, definition.Width, definition.Height)
	}
	// origin mantém todo o conjunto dentro da tela sem deformar o espaçamento entre linhas.
	origin := (screen.Bounds{X: x, Y: y, Width: definition.Width, Height: definition.Height}).Clamp(page.Size)
	// lines contém apenas componentes de texto comuns e editáveis.
	lines := presetLabels(presetID)
	if presetID == "pageTitle" {
		lines = []presetLabel{{"tituloPagina", "Título da página", screen.Bounds{Width: page.Size.Width, Height: 28}}}
	}
	// ids permite selecionar o primeiro controle e aplicar decoração na camada adaptadora.
	ids := []string{}
	// line fornece conteúdo e geometria de cada componente independente.
	for _, line := range lines {
		// kind mantém a faixa de título na família tipada correspondente.
		kind := catalog.Label
		if presetID == "pageTitle" {
			kind = catalog.Title
		}
		// component recebe UUID e defaults do mesmo catálogo usado pela paleta.
		component, err := screen.NewComponent(kind, UniqueName(*page, line.Name), 0, 0)
		if err != nil {
			return nil, err
		}
		component.FileID = UniqueComponentFileID(*page, component.FileID)
		component.Bounds = line.Bounds
		component.Bounds.X += origin.X
		component.Bounds.Y += origin.Y
		component.Style.FontSize = 11
		component.Style.TextColor = "#424344"
		if component.Properties.Label != nil {
			component.Properties.Label.Text = line.Text
		} else {
			component.Properties.Title.Text = line.Text
			component.Style.FontSize = 18
			component.Style.TextColor = "#FFFFFF"
			component.Style.TextAlign = "right"
		}
		page.Components = append(page.Components, component)
		ids = append(ids, component.ID)
	}
	return ids, nil
}
