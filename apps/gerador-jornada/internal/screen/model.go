// Package screen implementa a composição visual em coordenadas lógicas.
package screen

import (
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
)

// Size define as dimensões lógicas da superfície de desenho.
type Size struct {
	// Width expressa a largura em unidades lógicas.
	Width int
	// Height expressa a altura em unidades lógicas.
	Height int
}

// Bounds define a geometria lógica armazenada do componente na tela.
type Bounds struct {
	// X expressa a coordenada horizontal em unidades lógicas.
	X int
	// Y expressa a coordenada vertical em unidades lógicas.
	Y int
	// Width expressa a largura em unidades lógicas.
	Width int
	// Height expressa a altura em unidades lógicas.
	Height int
}

// Style define substituições opcionais da aparência padrão do catálogo.
type Style struct {
	// FontFamily define o nome da fonte; vazio utiliza o padrão do catálogo.
	FontFamily string
	// FontSize define o tamanho lógico da fonte; zero utiliza o padrão.
	FontSize int
	// TextColor define a cor RGB opcional do texto.
	TextColor string
	// FillColor define um fundo RGB opcional.
	FillColor string
	// BorderColor define um contorno RGB opcional.
	BorderColor string
	// BackgroundAssetID identifica uma imagem decorativa opcional.
	BackgroundAssetID string
}

// Background define a decoração da tela sem sintetizar controles.
type Background struct {
	// Color define a cor RGB da superfície.
	Color string
	// AssetID identifica uma imagem de fundo opcional.
	AssetID string
}

// Component reúne identidade estável e dados tipados da família.
type Component struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Name contém o nome exibido ao usuário; nomes de componentes são únicos na tela.
	Name string
	// Kind seleciona a variante tipada ou a ação pelos discriminadores estáveis do contrato.
	Kind catalog.Kind
	// Bounds define a geometria lógica armazenada do componente na tela.
	Bounds Bounds
	// Style define substituições opcionais da aparência padrão do catálogo.
	Style Style
	// Properties contém uma união discriminada com exatamente uma variante tipada.
	Properties catalog.Properties
}

// Screen contém componentes visuais planos ordenados do fundo para a frente.
type Screen struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Name contém o nome exibido ao usuário; nomes de componentes são únicos na tela.
	Name string
	// Size define as dimensões lógicas da superfície de desenho.
	Size Size
	// Background define a decoração da tela sem sintetizar controles.
	Background Background
	// Components ordena os componentes do fundo para a frente.
	Components []Component
}

// New cria uma tela vazia de 800 por 600 unidades lógicas com fundo branco.
func New(name string) Screen {
	return Screen{ID: project.NewID(), Name: name, Size: Size{800, 600}, Background: Background{Color: "#FFFFFF"}, Components: []Component{}}
}

// NewComponent cria uma instância do catálogo com novos identificadores de opções.
func NewComponent(kind catalog.Kind, name string, x, y int) (Component, error) {
	// definition contém as dimensões padrão da família solicitada no catálogo.
	// ok informa se a consulta encontrou um valor válido.
	definition, ok := catalog.Lookup(kind)
	if !ok {
		return Component{}, fmt.Errorf("tipo de componente desconhecido")
	}
	// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
	component := Component{ID: project.NewID(), Name: name, Kind: kind, Bounds: Bounds{x, y, definition.Width, definition.Height}}
	switch kind {
	case catalog.Button:
		component.Properties.Button = &catalog.ButtonProperties{Text: "Continuar", Variant: "confirm"}
	case catalog.Label:
		component.Properties.Label = &catalog.TextProperties{Text: "Rótulo"}
	case catalog.Title:
		component.Properties.Title = &catalog.TextProperties{Text: "Título da tela"}
	case catalog.TextField:
		component.Properties.Field = &catalog.FieldProperties{Format: "text"}
	case catalog.Radio:
		component.Properties.Radio = &catalog.RadioProperties{Orientation: "vertical", Options: []catalog.Option{{ID: project.NewID(), Text: "Opção 1"}, {ID: project.NewID(), Text: "Opção 2"}}}
	case catalog.Report:
		component.Properties.Report = &catalog.ReportProperties{Title: "Relatório", Rows: []catalog.Option{}}
	case catalog.NumericKeyboard:
		component.Properties.Keyboard = &catalog.KeyboardProperties{}
	}
	return component, nil
}

// ComponentByID retorna o componente mutável pertencente é tela, ou nil quando ausente.
func (s *Screen) ComponentByID(id string) *Component {
	// index localiza o item atual na coleção sem alterar sua ordem.
	for index := range s.Components {
		if s.Components[index].ID == id {
			return &s.Components[index]
		}
	}
	return nil
}

// Fits informa se a geometria positiva está inteiramente dentro da tela.
func (b Bounds) Fits(size Size) bool {
	return b.X >= 0 && b.Y >= 0 && b.Width > 0 && b.Height > 0 && b.X <= size.Width-b.Width && b.Y <= size.Height-b.Height
}

// Clamp ajusta o posicionamento é área de desenho disponível.
func (b Bounds) Clamp(size Size) Bounds {
	b.Width = min(max(1, b.Width), size.Width)
	b.Height = min(max(1, b.Height), size.Height)
	b.X = min(max(0, b.X), size.Width-b.Width)
	b.Y = min(max(0, b.Y), size.Height-b.Height)
	return b
}

// Validate verifica identidade, dimensões, nomes únicos, geometria, propriedades e cores locais da composição.
func (s Screen) Validate() error {
	if !project.ValidID(s.ID) || strings.TrimSpace(s.Name) == "" || s.Size.Width < 1 || s.Size.Height < 1 {
		return fmt.Errorf("tela inválida")
	}
	// names controla nomes únicos ou sua ordem lexical para empacotamento determinístico.
	names := map[string]bool{}
	// component contém o controle tipado em processamento, validado nas coordenadas e referências da própria tela.
	for _, component := range s.Components {
		if !project.ValidID(component.ID) || strings.TrimSpace(component.Name) == "" || names[component.Name] {
			return fmt.Errorf("componente sem identidade ou com nome repetido")
		}
		names[component.Name] = true
		if !component.Bounds.Fits(s.Size) {
			return fmt.Errorf("%s: posição/tamanho fora da tela", component.Name)
		}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if err := component.Properties.Validate(component.Kind); err != nil {
			return fmt.Errorf("%s: %w", component.Name, err)
		}
		if component.Style.FontSize < 0 {
			return fmt.Errorf("tamanho de fonte inválido")
		}
		// color contém uma substituição opcional de cor.
		for _, color := range []string{component.Style.TextColor, component.Style.FillColor, component.Style.BorderColor} {
			if color != "" && !catalog.ValidColor(color) {
				return fmt.Errorf("cor inválida")
			}
		}
	}
	return nil
}
