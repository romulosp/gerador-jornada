package presentation

import (
	"fmt"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
)

// newBarcodeComponent cria uma entrada com rótulo e geometria explícitos.
func newBarcodeComponent(t *testing.T, label string, bounds screen.Bounds) screen.Component {
	t.Helper()
	component, err := screen.NewComponent(catalog.BarcodeInput, "codigoBarras", bounds.X, bounds.Y)
	if err != nil {
		t.Fatalf("componente de código de barras inválido: %v", err)
	}
	component.Bounds = bounds
	component.Properties.BarcodeInput.Label = label
	return component
}

// primitivesByName devolve as primitivas com o nome exato, preservando a ordem de pintura.
func primitivesByName(group Group, name string) []Primitive {
	found := []Primitive{}
	for _, primitive := range group.Shapes {
		if primitive.Name == name {
			found = append(found, primitive)
		}
	}
	return found
}

// TestBarcodeFieldBoundsBase fixa a geometria de referência de 610 por 60.
func TestBarcodeFieldBoundsBase(t *testing.T) {
	// bounds reproduz a dimensão da definição do catálogo, sem escala.
	bounds := screen.Bounds{Width: 610, Height: 60}
	// baseFields repete posições e larguras contratadas na ordem visual.
	baseFields := []struct{ x, width int }{{30, 50}, {85, 50}, {140, 50}, {195, 60}, {260, 50}, {315, 60}, {380, 15}, {400, 110}}
	fields := BarcodeFieldBounds(bounds)
	if len(fields) != 8 {
		t.Fatalf("oito áreas clicáveis são obrigatórias: %d", len(fields))
	}
	for index, field := range fields {
		want := screen.Bounds{X: baseFields[index].x, Y: 20, Width: baseFields[index].width, Height: 20}
		if field != want {
			t.Fatalf("campo %d divergente: %+v != %+v", index+1, field, want)
		}
	}
	// baseSeparators repete os três pontos separadores contratados.
	baseSeparators := []struct{ x, y int }{{81, 27}, {191, 27}, {311, 27}}
	separators := BarcodeSeparatorBounds(bounds)
	if len(separators) != 3 {
		t.Fatalf("três separadores são obrigatórios: %d", len(separators))
	}
	for index, separator := range separators {
		want := screen.Bounds{X: baseSeparators[index].x, Y: baseSeparators[index].y, Width: 5, Height: 10}
		if separator != want {
			t.Fatalf("separador %d divergente: %+v != %+v", index+1, separator, want)
		}
	}
	// limits preserva os comprimentos dos oito campos na mesma ordem.
	limits := []int{5, 5, 5, 6, 5, 6, 1, 14}
	for index, limit := range limits {
		if BarcodeFieldLimit(index) != limit {
			t.Fatalf("limite do campo %d divergente: %d", index+1, BarcodeFieldLimit(index))
		}
	}
	if BarcodeFieldLimit(-1) != 0 || BarcodeFieldLimit(8) != 0 {
		t.Fatal("índice fora da faixa deve devolver limite zero")
	}
}

// TestBarcodeFieldBoundsScaled confere que pintura e áreas clicáveis usam a mesma escala.
func TestBarcodeFieldBoundsScaled(t *testing.T) {
	// bounds dobra a referência e desloca a origem para provar a transformação.
	bounds := screen.Bounds{X: 100, Y: 50, Width: 1220, Height: 120}
	fields := BarcodeFieldBounds(bounds)
	wantFields := []screen.Bounds{
		{X: 160, Y: 90, Width: 100, Height: 40},
		{X: 270, Y: 90, Width: 100, Height: 40},
		{X: 380, Y: 90, Width: 100, Height: 40},
		{X: 490, Y: 90, Width: 120, Height: 40},
		{X: 620, Y: 90, Width: 100, Height: 40},
		{X: 730, Y: 90, Width: 120, Height: 40},
		{X: 860, Y: 90, Width: 30, Height: 40},
		{X: 900, Y: 90, Width: 220, Height: 40},
	}
	for index, field := range fields {
		if field != wantFields[index] {
			t.Fatalf("campo %d escalado divergente: %+v != %+v", index+1, field, wantFields[index])
		}
	}
	separators := BarcodeSeparatorBounds(bounds)
	wantSeparators := []screen.Bounds{
		{X: 262, Y: 104, Width: 10, Height: 20},
		{X: 482, Y: 104, Width: 10, Height: 20},
		{X: 722, Y: 104, Width: 10, Height: 20},
	}
	for index, separator := range separators {
		if separator != wantSeparators[index] {
			t.Fatalf("separador %d escalado divergente: %+v != %+v", index+1, separator, wantSeparators[index])
		}
	}
}

// TestBarcodeInputStateShapes confere rótulo, oito fundos, oito textos, separadores e foco.
func TestBarcodeInputStateShapes(t *testing.T) {
	// component usa a geometria de referência e um rótulo Unicode.
	component := newBarcodeComponent(t, "Número do Código de Barras: ação", screen.Bounds{Width: 610, Height: 60})
	// values repete o preenchimento de F01 usado apenas na prévia.
	values := [8]string{"00190", "50095", "40144", "816069", "06809", "350314", "3", "37370000000100"}
	// state mantém o terceiro campo em foco e uma mensagem local.
	state := BarcodeState{Values: values, Active: 2, Message: "Código de barras inválido"}
	group := BarcodeInputState(component, state)
	if group.ComponentID != component.ID || group.Kind != catalog.BarcodeInput || group.Bounds != component.Bounds {
		t.Fatalf("identidade do grupo divergente: %+v", group)
	}
	// 1 rótulo + 8 fundos + 8 textos + 3 separadores + 1 mensagem resultam em 21 formas.
	if len(group.Shapes) != 21 {
		t.Fatalf("quantidade de formas divergente: %d", len(group.Shapes))
	}
	// label é a primeira forma e preserva texto, cor e posição do contrato.
	label := group.Shapes[0]
	if label.Name != "rótulo" || label.Text != component.Properties.BarcodeInput.Label {
		t.Fatalf("rótulo divergente: %+v", label)
	}
	if label.Color != "#5064AA" || label.FontSize != 12 {
		t.Fatalf("estilo do rótulo divergente: %+v", label)
	}
	if label.Bounds.X != 30 || label.Bounds.Y != 2 || label.Bounds.Width != 200 || label.Bounds.Height != 20 {
		t.Fatalf("geometria do rótulo divergente: %+v", label.Bounds)
	}
	// fields reutiliza a mesma função das áreas clicáveis para comparar pintura e hit area.
	fields := BarcodeFieldBounds(component.Bounds)
	for index := 0; index < 8; index++ {
		// fill e value são as formas consecutivas do campo visual.
		fill := group.Shapes[1+2*index]
		value := group.Shapes[2+2*index]
		if fill.Name != fmt.Sprintf("campo-%d", index+1) || value.Name != fmt.Sprintf("valor-%d", index+1) {
			t.Fatalf("campo %d fora da ordem esperada: %q/%q", index+1, fill.Name, value.Name)
		}
		if fill.Bounds != fields[index] {
			t.Fatalf("fundo %d não coincide com a área clicável: %+v", index+1, fill.Bounds)
		}
		if value.Text != values[index] {
			t.Fatalf("texto do campo %d divergente: %q", index+1, value.Text)
		}
	}
	// o terceiro campo recebe foco azul e texto branco.
	focusedFill, focusedValue := group.Shapes[5], group.Shapes[6]
	if focusedFill.Fill != "#10508E" || focusedFill.Stroke != "#10508E" {
		t.Fatalf("fundo em foco divergente: %+v", focusedFill)
	}
	if focusedValue.Color != "#FFFFFF" {
		t.Fatalf("texto em foco divergente: %+v", focusedValue)
	}
	// campos sem foco permanecem brancos com contorno neutro e texto preto.
	if group.Shapes[1].Fill != "#FFFFFF" || group.Shapes[1].Stroke != "#CBD5E1" || group.Shapes[2].Color != "#000000" {
		t.Fatalf("campo sem foco divergente: %+v/%+v", group.Shapes[1], group.Shapes[2])
	}
	// separators ocupam as três posições seguintes e usam a cor da marca.
	separators := BarcodeSeparatorBounds(component.Bounds)
	for index := 0; index < 3; index++ {
		separator := group.Shapes[17+index]
		if separator.Name != "separador" || separator.Fill != "#5064AA" || separator.Bounds != separators[index] {
			t.Fatalf("separador %d divergente: %+v", index+1, separator)
		}
	}
	// message é a última forma, textual e com quebra fora de cor apenas.
	message := group.Shapes[20]
	if message.Name != "mensagem" || message.Text != state.Message || message.Color != "#B91C1C" {
		t.Fatalf("mensagem divergente: %+v", message)
	}
	if message.Bounds.Y != 46 || message.Bounds.Height != 12 {
		t.Fatalf("geometria da mensagem divergente: %+v", message.Bounds)
	}
}

// TestBarcodeInputStateLocked confere o estado desabilitado sobre todos os campos.
func TestBarcodeInputStateLocked(t *testing.T) {
	// component usa a geometria de referência e um rótulo vazio permitido.
	component := newBarcodeComponent(t, "", screen.Bounds{Width: 610, Height: 60})
	// state combina leitura concluída e campo ativo: a trava deve prevalecer.
	group := BarcodeInputState(component, BarcodeState{Active: 1, Locked: true})
	if len(group.Shapes) != 20 {
		t.Fatalf("estado travado não deve exibir mensagem: %d", len(group.Shapes))
	}
	for index := 0; index < 8; index++ {
		fill := group.Shapes[1+2*index]
		value := group.Shapes[2+2*index]
		if fill.Fill != "#E6EEFF" || fill.Stroke != "#CBD5E1" {
			t.Fatalf("fundo travado %d divergente: %+v", index+1, fill)
		}
		if value.Color != "#808080" {
			t.Fatalf("texto travado %d divergente: %+v", index+1, value)
		}
	}
	// a forma sem estado de runtime usa os mesmos limites, sem valores digitados.
	static := BarcodeInput(component)
	if len(static.Shapes) != 20 || static.Shapes[0].Text != "" {
		t.Fatalf("entrada estática divergente: %d formas", len(static.Shapes))
	}
}

// TestBarcodeProcessedDerivedLabels confere os cinco rótulos derivados, a ordem e o Unicode.
func TestBarcodeProcessedDerivedLabels(t *testing.T) {
	// component usa a geometria da definição e um bloco com as cinco opções.
	component, err := screen.NewComponent(catalog.BarcodeProcessed, "codigoProcessado", 0, 0)
	if err != nil {
		t.Fatalf("componente processado inválido: %v", err)
	}
	component.Bounds = screen.Bounds{Width: 530, Height: 120}
	// values preserva acentos, vírgulas e pontuação literal.
	values := []string{"123.456.789-09", "12.345.678/0001-90", "José da Silva", "Ação & Cia.", "Razão Social ção"}
	// kinds segue a ordem visual das cinco opções do combo.
	kinds := []string{catalog.BarcodeLabelCPF, catalog.BarcodeLabelCNPJ, catalog.BarcodeLabelName, catalog.BarcodeLabelTradeName, catalog.BarcodeLabelCorporateName}
	rows := make([]catalog.BarcodeRow, len(kinds))
	for index := range kinds {
		rows[index] = catalog.BarcodeRow{ID: project.NewID(), LabelKind: kinds[index], Value: values[index]}
	}
	component.Properties.BarcodeProcessed = &catalog.BarcodeProcessedProperties{Blocks: []catalog.BarcodeBlock{{
		ID: project.NewID(), Title: "Dados do boleto", Rows: rows,
	}}}
	group := BarcodeProcessed(component)
	if group.ComponentID != component.ID || group.Kind != catalog.BarcodeProcessed {
		t.Fatalf("identidade do grupo divergente: %+v", group)
	}
	// 1 título + 5 pares de rótulo/valor resultam em 11 formas.
	if len(group.Shapes) != 11 {
		t.Fatalf("quantidade de formas divergente: %d", len(group.Shapes))
	}
	title := group.Shapes[0]
	if title.Name != "bloco-1-título" || title.Text != "Dados do boleto" || !title.Bold || title.Color != "#1E3A8A" {
		t.Fatalf("título do bloco divergente: %+v", title)
	}
	// wantLabels repete o texto derivado do enum, nunca uma segunda string persistida.
	wantLabels := []string{"CPF:", "CNPJ:", "Nome:", "Nome fantasia:", "Razão social:"}
	for index := range wantLabels {
		label := group.Shapes[1+2*index]
		value := group.Shapes[2+2*index]
		if label.Name != fmt.Sprintf("bloco-1-linha-%d-rótulo", index+1) || value.Name != fmt.Sprintf("bloco-1-linha-%d-valor", index+1) {
			t.Fatalf("linha %d fora da ordem esperada: %q/%q", index+1, label.Name, value.Name)
		}
		if label.Text != wantLabels[index] {
			t.Fatalf("rótulo %d divergente: %q", index+1, label.Text)
		}
		if value.Text != values[index] {
			t.Fatalf("valor %d não preservou Unicode: %q", index+1, value.Text)
		}
	}
	// a primeira linha mantém a geometria de leitura: rótulo à esquerda, valor à direita.
	if group.Shapes[1].Bounds.X != 12 || group.Shapes[1].Bounds.Width != 202 {
		t.Fatalf("geometria do rótulo divergente: %+v", group.Shapes[1].Bounds)
	}
	if group.Shapes[2].Bounds.X != 222 || group.Shapes[2].Bounds.Width != 296 {
		t.Fatalf("geometria do valor divergente: %+v", group.Shapes[2].Bounds)
	}
}

// TestBarcodeProcessedBlocksAndEmptyValues confere título vazio, valor vazio e ordem dos blocos.
func TestBarcodeProcessedBlocksAndEmptyValues(t *testing.T) {
	// component contém dois blocos: o primeiro sem título e o segundo nomeado.
	component, err := screen.NewComponent(catalog.BarcodeProcessed, "codigoProcessado", 10, 20)
	if err != nil {
		t.Fatalf("componente processado inválido: %v", err)
	}
	component.Bounds = screen.Bounds{X: 10, Y: 20, Width: 530, Height: 120}
	component.Properties.BarcodeProcessed = &catalog.BarcodeProcessedProperties{Blocks: []catalog.BarcodeBlock{
		{
			ID: project.NewID(), Title: "",
			Rows: []catalog.BarcodeRow{{ID: project.NewID(), LabelKind: catalog.BarcodeLabelName, Value: ""}},
		},
		{
			ID: project.NewID(), Title: "Segundo bloco",
			Rows: []catalog.BarcodeRow{{ID: project.NewID(), LabelKind: catalog.BarcodeLabelCPF, Value: "000.000.000-00"}},
		},
	}}
	group := BarcodeProcessed(component)
	// o primeiro bloco não reserva faixa de título, mas mantém o par rótulo/valor vazio.
	if shapes := primitivesByName(group, "bloco-1-título"); len(shapes) != 0 {
		t.Fatalf("título vazio não deve gerar forma: %+v", shapes)
	}
	emptyValue := primitivesByName(group, "bloco-1-linha-1-valor")
	if len(emptyValue) != 1 || emptyValue[0].Text != "" {
		t.Fatalf("valor vazio deve permanecer como caixa editável: %+v", emptyValue)
	}
	if label := primitivesByName(group, "bloco-1-linha-1-rótulo"); len(label) != 1 || label[0].Text != "Nome:" {
		t.Fatalf("rótulo da linha vazia divergente: %+v", label)
	}
	// o segundo bloco vem depois e preserva o título próprio.
	names := []string{}
	for _, shape := range group.Shapes {
		names = append(names, shape.Name)
	}
	wantOrder := []string{
		"bloco-1-linha-1-rótulo", "bloco-1-linha-1-valor",
		"bloco-2-título",
		"bloco-2-linha-1-rótulo", "bloco-2-linha-1-valor",
	}
	if len(names) != len(wantOrder) {
		t.Fatalf("ordem de formas inesperada: %v", names)
	}
	for index, name := range wantOrder {
		if names[index] != name {
			t.Fatalf("forma %d divergente: %q != %q", index+1, names[index], name)
		}
	}
	if title := primitivesByName(group, "bloco-2-título"); len(title) != 1 || title[0].Text != "Segundo bloco" {
		t.Fatalf("título do segundo bloco divergente: %+v", title)
	}
}
