package mcp

import (
	"encoding/json"
	"reflect"
	"sync"

	"gerador-jornada.local/editor/internal/application"
	"github.com/google/jsonschema-go/jsonschema"
)

// boundsDTO exige geometria integral em coordenadas lógicas.
type boundsDTO struct {
	// X/Y são posição; Width/Height são dimensões positivas na validação comum.
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// stylePatch adapta somente campos reais de screen.Style, sem inventar alpha separado.
type stylePatch struct {
	// FontFamily/FontSize preservam defaults legados quando vazios/zero deliberados.
	FontFamily *string `json:"fontFamily,omitempty"`
	FontSize   *int    `json:"fontSize,omitempty"`
	// TextColor/FillColor/BorderColor são overrides explícitos sujeitos à validação comum.
	TextColor   *string `json:"textColor,omitempty"`
	FillColor   *string `json:"fillColor,omitempty"`
	BorderColor *string `json:"borderColor,omitempty"`
	// BackgroundAssetID referencia recurso já incorporado, nunca um arquivo externo.
	BackgroundAssetID *string `json:"backgroundAssetId,omitempty"`
	// TextAlign/Skin mantêm enums e compatibilidade da apresentação instalada.
	TextAlign *string `json:"textAlign,omitempty"`
	Skin      *string `json:"skin,omitempty"`
}

// propertyKind fornece o discriminador comum aos sete patches suportados em I1.
type propertyKind struct {
	// Kind precisa coincidir com a entidade existente, sem conversão implícita.
	Kind string `json:"kind"`
}

// family identifica um patch tipado sem expor map de propriedades ao domínio.
func (kind propertyKind) family() string { return kind.Kind }

// familyPatch permite somente variantes concretas construídas pelo decoder do adaptador.
type familyPatch interface{ family() string }

// buttonPatch altera texto/variante do botão instalado.
type buttonPatch struct {
	// propertyKind distingue conteúdo de botão de campos de outras famílias.
	propertyKind
	// Text/Variant omitidos conservam valores atuais.
	Text    *string `json:"text,omitempty"`
	Variant *string `json:"variant,omitempty"`
}

// textPatch serve rótulo e título, conservando seus Kinds distintos.
type textPatch struct {
	// propertyKind aceita label/title exclusivamente no schema discriminado.
	propertyKind
	// Text é conteúdo Unicode sem normalização.
	Text *string `json:"text,omitempty"`
}

// fieldPatch altera somente valor inicial e formato instalado do campo.
type fieldPatch struct {
	// propertyKind não permite aplicar este patch a outro Kind.
	propertyKind
	// InitialValue/Format têm limites e enums do mesmo caso de uso manual.
	InitialValue *string `json:"initialValue,omitempty"`
	Format       *string `json:"format,omitempty"`
}

// optionDTO conserva identidades explícitas na troca integral de opções/linhas.
type optionDTO struct {
	// ID/Text são identificador estável e conteúdo literal, com validações independentes.
	ID   string `json:"id"`
	Text string `json:"text"`
}

// radioPatch representa seleção única sem campos de relatório ou botão.
type radioPatch struct {
	// propertyKind identifica a variante radio.
	propertyKind
	// Options troca a lista inteira; omitir difere de enviar lista vazia.
	Options *[]optionDTO `json:"options,omitempty"`
	// Orientation/SelectedOptionID conservam enums/referências do domínio.
	Orientation      *string `json:"orientation,omitempty"`
	SelectedOptionID *string `json:"selectedOptionId,omitempty"`
}

// reportPatch representa conteúdo do relatório, independente de captura visual.
type reportPatch struct {
	// propertyKind identifica a variante report.
	propertyKind
	// Title/Rows são texto e lista integral de linhas identificadas.
	Title *string      `json:"title,omitempty"`
	Rows  *[]optionDTO `json:"rows,omitempty"`
}

// keyboardPatch altera somente referência ao campo na mesma tela.
type keyboardPatch struct {
	// propertyKind identifica numericKeyboard, sem novos comportamentos implícitos.
	propertyKind
	// TargetComponentID vazio desvincula; omissão conserva o vínculo.
	TargetComponentID *string `json:"targetComponentId,omitempty"`
}

// propertiesPatch conserva uma união discriminada fechada, nunca propriedades livres.
type propertiesPatch struct {
	// value existe apenas no adaptador e é convertido para mudanças tipadas da aplicação.
	value familyPatch
}

// UnmarshalJSON seleciona schema/struct da família antes de aceitar qualquer propriedade.
func (patch *propertiesPatch) UnmarshalJSON(raw []byte) error {
	// head só seleciona a variante; o decoder fechado seguinte valida todos os campos.
	var head propertyKind
	if err := json.Unmarshal(raw, &head); err != nil {
		return &application.AuthoringError{Code: "invalidArgument"}
	}
	// err é preservado sem expor os valores do conteúdo ao cliente.
	var err error
	switch head.Kind {
	case "button":
		patch.value, err = decodeArguments[buttonPatch](raw)
	case "label", "title":
		patch.value, err = decodeArguments[textPatch](raw)
	case "textField":
		patch.value, err = decodeArguments[fieldPatch](raw)
	case "radio":
		patch.value, err = decodeArguments[radioPatch](raw)
	case "report":
		patch.value, err = decodeArguments[reportPatch](raw)
	case "numericKeyboard":
		patch.value, err = decodeArguments[keyboardPatch](raw)
	case "barcodeInput", "barcodeProcessed", "applicationMenu":
		return &application.AuthoringError{Code: "capabilityUnavailable"}
	default:
		return &application.AuthoringError{Code: "invalidArgument"}
	}
	return err
}

// MarshalJSON mantém a forma plana do contrato sem serializar o campo interno value.
func (patch propertiesPatch) MarshalJSON() ([]byte, error) { return json.Marshal(patch.value) }

// componentChanges fecha campos comuns e a união tipada de propriedades.
type componentChanges struct {
	// Name/FileID têm regras próprias, sem confundir nome visual e arquivo gerado.
	Name   *string `json:"name,omitempty"`
	FileID *string `json:"fileId,omitempty"`
	// Enabled/Bounds/Style são patches tipados, sem acesso a ParentID ou coleção interna.
	Enabled *bool       `json:"enabled,omitempty"`
	Bounds  *boundsDTO  `json:"bounds,omitempty"`
	Style   *stylePatch `json:"style,omitempty"`
	// Properties é fechada por Kind e não permite converter tipos.
	Properties *propertiesPatch `json:"properties,omitempty"`
}

// componentUpdateRequest conserva entidade alvo e patch tipado numa única edição.
type componentUpdateRequest struct {
	// mutationRequest mantém identidade e precondição da edição.
	mutationRequest
	// ScreenID/ComponentID devem pertencer ao mesmo documento da sessão.
	ScreenID    string `json:"screenId"`
	ComponentID string `json:"componentId"`
	// Changes não transporta JSON livre ao domínio.
	Changes componentChanges `json:"changes"`
}

// propertySchemasOnce/propertySchema inicializam a união imutável antes de resolver DTOs.
var propertySchemasOnce sync.Once
var propertySchema *jsonschema.Schema

// componentSchemaOptions mapeia a união interna para o schema discriminado oferecido ao cliente.
func componentSchemaOptions() *jsonschema.ForOptions {
	propertySchemasOnce.Do(func() {
		// union é construída com tipos concretos; nenhum schema é recebido do agente.
		union := &jsonschema.Schema{Type: "object", OneOf: []*jsonschema.Schema{}}
		// variant associa o tipo Go ao conjunto fechado de Kinds que pode representar.
		for _, variant := range []struct {
			kind  reflect.Type
			names []any
		}{
			{reflect.TypeFor[buttonPatch](), []any{"button"}}, {reflect.TypeFor[textPatch](), []any{"label", "title"}},
			{reflect.TypeFor[fieldPatch](), []any{"textField"}}, {reflect.TypeFor[radioPatch](), []any{"radio"}},
			{reflect.TypeFor[reportPatch](), []any{"report"}}, {reflect.TypeFor[keyboardPatch](), []any{"numericKeyboard"}},
		} {
			// schema fecha campos por família e usa enum no discriminador.
			schema, err := jsonschema.ForType(variant.kind, nil)
			if err != nil {
				panic("schema de família inválido")
			}
			schema.Properties["kind"].Enum = variant.names
			union.OneOf = append(union.OneOf, schema)
		}
		propertySchema = union
	})
	return &jsonschema.ForOptions{TypeSchemas: map[reflect.Type]*jsonschema.Schema{reflect.TypeFor[propertiesPatch](): propertySchema}}
}
