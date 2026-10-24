package catalog

import (
	"fmt"

	"gerador-jornada.local/editor/internal/project"
)

// Tipos de rótulo aceitos pela ficha de código de barras processado.
const (
	// BarcodeLabelCPF identifica uma linha de CPF.
	BarcodeLabelCPF = "cpf"
	// BarcodeLabelCNPJ identifica uma linha de CNPJ.
	BarcodeLabelCNPJ = "cnpj"
	// BarcodeLabelName identifica uma linha de nome.
	BarcodeLabelName = "name"
	// BarcodeLabelTradeName identifica uma linha de nome fantasia.
	BarcodeLabelTradeName = "tradeName"
	// BarcodeLabelCorporateName identifica uma linha de razão social.
	BarcodeLabelCorporateName = "corporateName"
)

// BarcodeRow representa uma linha ordenada de rótulo tipado e valor ilustrativo.
type BarcodeRow struct {
	// ID identifica a linha de forma estável dentro do documento.
	ID string
	// LabelKind seleciona o texto do rótulo derivado, nunca uma segunda string concorrente.
	LabelKind string
	// Value contém o texto Unicode ilustrativo, inclusive vazio.
	Value string
}

// BarcodeBlock representa um bloco ordenado com título opcional e linhas não vazias.
type BarcodeBlock struct {
	// ID identifica o bloco de forma estável dentro do documento.
	ID string
	// Title contém o título Unicode opcional, inclusive vazio.
	Title string
	// Rows mantém a ordem visual das linhas.
	Rows []BarcodeRow
}

// BarcodeProcessedProperties define a ficha de leitura parametrizável.
type BarcodeProcessedProperties struct {
	// Blocks mantém a ordem dos blocos; deve conter ao menos um bloco.
	Blocks []BarcodeBlock
}

// NewBarcodeProcessedProperties cria a ficha padrão com um bloco e duas linhas vazias.
// supplyID entrega UUIDs novos, mantendo o domínio independente de relógio e entropia própria.
func NewBarcodeProcessedProperties(supplyID func() string) *BarcodeProcessedProperties {
	return &BarcodeProcessedProperties{Blocks: []BarcodeBlock{{
		ID:    supplyID(),
		Title: "",
		Rows: []BarcodeRow{
			{ID: supplyID(), LabelKind: BarcodeLabelCNPJ, Value: ""},
			{ID: supplyID(), LabelKind: BarcodeLabelName, Value: ""},
		},
	}}}
}

// NewBarcodeRow cria uma linha padrão de nome com valor vazio e identidade nova.
func NewBarcodeRow(supplyID func() string) BarcodeRow {
	return BarcodeRow{ID: supplyID(), LabelKind: BarcodeLabelName, Value: ""}
}

// NewBarcodeBlock cria um bloco vazio com uma linha padrão e identidade nova.
func NewBarcodeBlock(supplyID func() string) BarcodeBlock {
	return BarcodeBlock{ID: supplyID(), Title: "", Rows: []BarcodeRow{NewBarcodeRow(supplyID)}}
}

// Validate rejeita estrutura vazia, identidades repetidas ou inválidas e rótulo desconhecido.
func (p BarcodeProcessedProperties) Validate() error {
	if len(p.Blocks) == 0 {
		return fmt.Errorf("ao menos um bloco é necessário")
	}
	// seen impede que blocos ou linhas compartilhem identidade.
	seen := map[string]bool{}
	// blockIndex percorre os blocos na ordem confirmada.
	for blockIndex, block := range p.Blocks {
		if !project.ValidID(block.ID) || seen[block.ID] {
			return fmt.Errorf("bloco %d sem identidade única", blockIndex+1)
		}
		seen[block.ID] = true
		if len(block.Rows) == 0 {
			return fmt.Errorf("bloco %d sem linhas", blockIndex+1)
		}
		// rowIndex percorre as linhas preservando a ordem visual.
		for rowIndex, row := range block.Rows {
			if !project.ValidID(row.ID) || seen[row.ID] {
				return fmt.Errorf("linha %d do bloco %d sem identidade única", rowIndex+1, blockIndex+1)
			}
			seen[row.ID] = true
			if _, ok := BarcodeLabelText(row.LabelKind); !ok {
				return fmt.Errorf("linha %d do bloco %d com rótulo desconhecido", rowIndex+1, blockIndex+1)
			}
		}
	}
	return nil
}

// CollectBarcodeProcessedIDs devolve todas as identidades internas para validação global.
func CollectBarcodeProcessedIDs(p BarcodeProcessedProperties) []string {
	// ids preserva a ordem de blocos e linhas.
	ids := []string{}
	// block percorre a configuração confirmada.
	for _, block := range p.Blocks {
		ids = append(ids, block.ID)
		// row percorre as linhas do bloco.
		for _, row := range block.Rows {
			ids = append(ids, row.ID)
		}
	}
	return ids
}
