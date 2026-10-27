package catalog

import "fmt"

// PropertyDefinition descreve uma propriedade instalada sem depender de widgets ou transporte.
type PropertyDefinition struct {
	// Name/Type/Description são identidade, tipo lógico e finalidade para inspector/descoberta.
	Name, Type, Description string
	// EnumValues limita escolhas aos mesmos valores usados por Validate.
	EnumValues []string
	// MaxCodePoints limita novas entradas textuais; zero indica outro limite de domínio.
	MaxCodePoints int
}

// ButtonVariants devolve cópia das variantes aceitas pela regra de domínio.
func ButtonVariants() []string { return []string{"confirm", "cancel", "yes", "no"} }

// FieldFormats devolve os formatos sem inferir conversões ou ações externas.
func FieldFormats() []string {
	return []string{"text", "number", "cellular", "cep", "cpf", "date", "ddd", "nis"}
}

// RadioOrientations devolve as orientações aceitas pela seleção única.
func RadioOrientations() []string { return []string{"vertical", "horizontal"} }

// PropertyDefinitions descreve somente patches das sete famílias I1, preservando os compostos existentes.
// Lista vazia para composto não remove suas propriedades da GUI nem promete edição MCP avançada.
func PropertyDefinitions(kind Kind) []PropertyDefinition {
	switch kind {
	case Button:
		return []PropertyDefinition{{Name: "text", Type: "string", Description: "Texto do botão", MaxCodePoints: 65536}, {Name: "variant", Type: "enum", Description: "Variante instalada", EnumValues: ButtonVariants()}}
	case Label, Title:
		return []PropertyDefinition{{Name: "text", Type: "string", Description: "Texto exibido", MaxCodePoints: 65536}}
	case TextField:
		return []PropertyDefinition{{Name: "initialValue", Type: "string", Description: "Valor inicial do campo", MaxCodePoints: 65536}, {Name: "format", Type: "enum", Description: "Formato instalado", EnumValues: FieldFormats()}}
	case Radio:
		return []PropertyDefinition{{Name: "options", Type: "collection", Description: "Lista integral de opções identificadas"}, {Name: "orientation", Type: "enum", Description: "Orientação", EnumValues: RadioOrientations()}, {Name: "selectedOptionId", Type: "reference", Description: "Opção inicial ou vazio"}}
	case Report:
		return []PropertyDefinition{{Name: "title", Type: "string", Description: "Título do relatório", MaxCodePoints: 65536}, {Name: "rows", Type: "collection", Description: "Lista integral de linhas identificadas"}}
	case NumericKeyboard:
		return []PropertyDefinition{{Name: "targetComponentId", Type: "reference", Description: "Campo de destino na mesma tela ou vazio"}}
	}
	return []PropertyDefinition{}
}

// DefaultProperties centraliza defaults de criação das dez variantes tipadas instaladas.
// newID é chamado para opções/partes que exigem identidade própria em cada inserção.
func DefaultProperties(kind Kind, newID func() string) (Properties, error) {
	if newID == nil {
		return Properties{}, fmt.Errorf("gerador de identidade ausente")
	}
	// properties contém exatamente uma variante, sem recursos de UI ou IO.
	properties := Properties{}
	switch kind {
	case Button:
		properties.Button = &ButtonProperties{Text: "Continuar", Variant: "confirm"}
	case Label:
		properties.Label = &TextProperties{Text: "Rótulo"}
	case Title:
		properties.Title = &TextProperties{Text: "Título da tela"}
	case TextField:
		properties.Field = &FieldProperties{Format: "text"}
	case Radio:
		properties.Radio = &RadioProperties{Orientation: "vertical", Options: []Option{{ID: newID(), Text: "Opção 1"}, {ID: newID(), Text: "Opção 2"}}}
	case Report:
		properties.Report = &ReportProperties{Title: "Relatório", Rows: []Option{}}
	case NumericKeyboard:
		properties.Keyboard = &KeyboardProperties{}
	case BarcodeInput:
		properties.BarcodeInput = NewBarcodeInputProperties()
	case BarcodeProcessed:
		properties.BarcodeProcessed = NewBarcodeProcessedProperties(newID)
	case ApplicationMenu:
		properties.Menu = NewApplicationMenuProperties(newID)
	default:
		return Properties{}, fmt.Errorf("tipo de componente desconhecido")
	}
	return properties, nil
}
