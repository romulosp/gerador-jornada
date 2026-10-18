// Package catalog define as sete famílias tipadas de componentes.
package catalog

import (
	"fmt"
	"regexp"
)

// Kind identifica uma família estável de componentes do catálogo.
type Kind string

const (
	// Button representa um controle acionável; nas propriedades, é preenchido somente para button.
	Button Kind = "button"
	// Label representa texto simples; nas propriedades, é preenchido somente para label.
	Label Kind = "label"
	// Title representa texto destacado; nas propriedades, é preenchido somente para title.
	Title Kind = "title"
	// TextField representa uma entrada editável.
	TextField Kind = "textField"
	// Radio representa uma seleção única; nas propriedades, é preenchido somente para radio.
	Radio Kind = "radio"
	// Report representa uma lista visual; nas propriedades, é preenchido somente para report.
	Report Kind = "report"
	// NumericKeyboard representa um teclado de entrada numérica.
	NumericKeyboard Kind = "numericKeyboard"
)

// Definition descreve uma família sem depender da implementação da interface.
type Definition struct {
	// Kind identifica uma família estável de componentes do catálogo.
	Kind Kind
	// Name contém o nome exibido ao usuário; nomes de componentes são únicos na tela.
	Name string
	// Width expressa a largura em unidades lógicas.
	Width int
	// Height expressa a altura em unidades lógicas.
	Height int
}

// Option representa uma opção ordenada ou linha de relatório.
type Option struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Text contém o texto Unicode exibido; typeText permite inclusive uma string vazia.
	Text string
}

// ButtonProperties define o conteúdo visual de um controle de ação.
type ButtonProperties struct {
	// Text contém o texto Unicode exibido; typeText permite inclusive uma string vazia.
	Text string
	// Variant aceita confirm, cancel, yes ou no.
	Variant string
}

// TextProperties define o conteúdo de rótulos e títulos.
type TextProperties struct {
	// Text contém o texto Unicode exibido; typeText permite inclusive uma string vazia.
	Text string
}

// FieldProperties define o valor inicial e a intenção de formatação da entrada.
type FieldProperties struct {
	// InitialValue é restaurado a cada entrada na tela.
	InitialValue string
	// Format declara a intenção de entrada suportada.
	Format string
}

// RadioProperties define um conjunto de seleção única.
type RadioProperties struct {
	// Options contém a coleção ordenada de opções.
	Options []Option
	// Orientation aceita vertical ou horizontal.
	Orientation string
	// SelectedOptionID permanece vazio quando não há seleção inicial.
	SelectedOptionID string
}

// ReportProperties define uma lista visual de relatório.
type ReportProperties struct {
	// Title contém o cabeçalho do relatório.
	Title string
	// Rows contém a coleção ordenada de linhas.
	Rows []Option
}

// KeyboardProperties define o alvo do teclado na mesma tela.
type KeyboardProperties struct {
	// TargetComponentID identifica um campo de texto ou permanece vazio.
	TargetComponentID string
}

// Properties contém uma união discriminada com exatamente uma variante tipada.
type Properties struct {
	// Button representa um controle acionável; nas propriedades, é preenchido somente para button.
	Button *ButtonProperties
	// Label representa texto simples; nas propriedades, é preenchido somente para label.
	Label *TextProperties
	// Title representa texto destacado; nas propriedades, é preenchido somente para title.
	Title *TextProperties
	// Field é preenchido somente para textField.
	Field *FieldProperties
	// Radio representa uma seleção única; nas propriedades, é preenchido somente para radio.
	Radio *RadioProperties
	// Report representa uma lista visual; nas propriedades, é preenchido somente para report.
	Report *ReportProperties
	// Keyboard é preenchido somente para numericKeyboard.
	Keyboard *KeyboardProperties
}

// Definitions retorna uma nova paleta ordenada.
func Definitions() []Definition {
	return []Definition{{Button, "Botão", 110, 60}, {Label, "Rótulo", 150, 20}, {Title, "Título", 300, 28}, {TextField, "Campo", 150, 25}, {Radio, "Seleção", 180, 60}, {Report, "Relatório", 300, 345}, {NumericKeyboard, "Teclado numérico", 200, 280}}
}

// Lookup resolve uma família pelo discriminador estável.
func Lookup(kind Kind) (Definition, bool) {
	// definition contém as dimensões padrão da família solicitada no catálogo.
	for _, definition := range Definitions() {
		if definition.Kind == kind {
			return definition, true
		}
	}
	return Definition{}, false
}

// Validate rejeita discriminadores desconhecidos, variantes misturadas e valores fora das enumerações permitidas.
func (p Properties) Validate(kind Kind) error {
	// count impede a retenção silenciosa de propriedades de outra família.
	count := 0
	if p.Button != nil {
		count++
	}
	if p.Label != nil {
		count++
	}
	if p.Title != nil {
		count++
	}
	if p.Field != nil {
		count++
	}
	if p.Radio != nil {
		count++
	}
	if p.Report != nil {
		count++
	}
	if p.Keyboard != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("uma única variante de propriedades é necessária")
	}
	switch kind {
	case Button:
		if p.Button != nil && contains(p.Button.Variant, "confirm", "cancel", "yes", "no") {
			return nil
		}
	case Label:
		if p.Label != nil {
			return nil
		}
	case Title:
		if p.Title != nil {
			return nil
		}
	case TextField:
		if p.Field != nil && contains(p.Field.Format, "text", "number", "cellular", "cep", "cpf", "date", "ddd", "nis") {
			return nil
		}
	case Radio:
		if p.Radio != nil && len(p.Radio.Options) > 0 && contains(p.Radio.Orientation, "vertical", "horizontal") {
			// seen detecta identidades ou nomes repetidos, inclusive diferenças de caixa em caminhos Windows.
			seen := map[string]bool{}
			// selected registra a validade da escolha inicial.
			selected := p.Radio.SelectedOptionID == ""
			// option representa uma escolha ordenada.
			for _, option := range p.Radio.Options {
				if option.ID == "" || seen[option.ID] {
					return fmt.Errorf("opção inválida ou repetida")
				}
				seen[option.ID] = true
				selected = selected || option.ID == p.Radio.SelectedOptionID
			}
			if selected {
				return nil
			}
		}
	case Report:
		if p.Report != nil {
			return nil
		}
	case NumericKeyboard:
		if p.Keyboard != nil {
			return nil
		}
	}
	return fmt.Errorf("propriedades inválidas para %s", kind)
}

// contains verifica se o valor pertence a uma enumeração fechada pequena.
func contains(value string, allowed ...string) bool {
	// candidate contém o valor candidato, preservado enquanto não houver colisão.
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

// colorPattern valida a notação RGB utilizada pelo projeto.
var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// ValidColor aceita cores RGB no formato do catálogo.
func ValidColor(value string) bool { return colorPattern.MatchString(value) }
