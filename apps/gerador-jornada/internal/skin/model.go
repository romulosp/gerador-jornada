// Package skin interpreta as definições visuais do acervo Java/EFL incorporado.
package skin

import (
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"

	"gerador-jornada.local/editor/internal/catalog"
)

// assets contém o skin.xml e os recursos resolvidos do material legado.
//
//go:embed assets
var assets embed.FS

// Color representa uma cor RGBA do skin.
type Color struct {
	// R contém o canal vermelho de 0 a 255.
	R int `json:"R"`
	// G contém o canal verde de 0 a 255.
	G int `json:"G"`
	// B contém o canal azul de 0 a 255.
	B int `json:"B"`
	// A contém o alfa de 0 a 255.
	A int `json:"A"`
}

// Hex devolve a cor no formato #RRGGBB, ignorando o alfa.
func (c Color) Hex() string {
	if c.A == 0 {
		return ""
	}
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// Font descreve a face, o tamanho e a cor declarados em um estado.
type Font struct {
	// Type contém o caminho lógico da fonte no acervo.
	Type string `json:"Type"`
	// Size expressa o tamanho em unidades lógicas.
	Size int `json:"Size"`
	// Bold indica o peso forte declarado.
	Bold bool `json:"Bold"`
	// Italic indica o estilo itálico declarado.
	Italic bool `json:"Italic"`
	// Color contém a cor do texto.
	Color Color `json:"Color"`
}

// State representa um estado visual nomeado do skin.
type State struct {
	// Spacement informa o fator entre linhas declarado no estado original.
	Spacement float64 `json:"Spacement"`
	// Name identifica default, focused, disabled, touched, pressed ou released.
	Name string `json:"Name"`
	// ImagePath contém o caminho lógico da imagem, quando houver.
	ImagePath string `json:"ImagePath"`
	// Background contém a cor de fundo declarada.
	Background Color `json:"Background"`
	// Font contém a tipografia declarada.
	Font *Font `json:"Font"`
}

// Definition reúne os estados de um elemento nomeado do skin.
type Definition struct {
	// Type identifica panel, label, button, textField, toggleButton, table etc.
	Type string `json:"Type"`
	// Name identifica a definição no acervo.
	Name string `json:"Name"`
	// States associa o nome do estado à sua definição.
	States map[string]State `json:"States"`
}

// Library mantém as definições incorporadas e o acesso aos recursos.
type Library struct {
	// definitions associa tipo e nome à definição interpretada.
	definitions map[string]map[string]Definition
	// assets contém os recursos embutidos por caminho lógico.
	assets fs.FS
	// resources conserva hashes e famílias reais dos recursos incorporados.
	resources map[string]Resource
}

// builtin guarda a biblioteca carregada uma única vez.
var builtin struct {
	// once garante a leitura única do XML.
	once sync.Once
	// library contém a biblioteca pronta.
	library *Library
	// err preserva uma falha de leitura para consulta.
	err error
}

// Builtin devolve a biblioteca incorporada, carregando o XML na primeira chamada.
func Builtin() (*Library, error) {
	builtin.once.Do(func() {
		builtin.library, builtin.err = load()
	})
	return builtin.library, builtin.err
}

// MustBuiltin devolve a biblioteca incorporada ou interrompe uma falha de integridade interna.
func MustBuiltin() *Library {
	library, err := Builtin()
	if err != nil || library == nil {
		panic(fmt.Sprintf("biblioteca visual incorporada inválida: %v", err))
	}
	return library
}

// Resource descreve a integridade e a família tipográfica de um recurso incorporado.
type Resource struct {
	// Bold identifica o peso físico no cabeçalho SFNT, independente das flags XML.
	Bold bool `json:"Bold"`
	// Italic identifica o itálico físico da face SFNT.
	Italic bool `json:"Italic"`
	// SHA256 identifica os bytes importados sem depender do nome do arquivo.
	SHA256 string `json:"SHA256"`
	// Size informa o número de bytes esperados.
	Size int `json:"Size"`
	// Family informa a família SFNT real, vazia para imagens.
	Family string `json:"Family"`
}

// load lê somente dados gerados; nenhum XML Java é interpretado em execução.
func load() (*Library, error) {
	// sub isola os caminhos lógicos dos recursos incorporados.
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	// data contém a biblioteca estática gerada pelo importador de desenvolvimento.
	data, err := fs.ReadFile(sub, "library.json")
	if err != nil {
		return nil, err
	}
	// manifest conserva definições e integridade de todos os recursos disponíveis.
	var manifest struct {
		// Definitions agrupa as skins por tipo original.
		Definitions map[string]map[string]Definition
		// Resources descreve os bytes usados pelas definições operacionais.
		Resources map[string]Resource
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	// path e resource verificam a integridade antes de disponibilizar qualquer skin.
	for path, resource := range manifest.Resources {
		// content permite detectar recursos ausentes ou alterados após a importação.
		content, err := fs.ReadFile(sub, strings.TrimPrefix(path, "/"))
		if err != nil {
			return nil, fmt.Errorf("recurso de skin %s: %w", path, err)
		}
		if len(content) != resource.Size || fmt.Sprintf("%x", sha256.Sum256(content)) != resource.SHA256 {
			return nil, fmt.Errorf("recurso de skin corrompido: %s", path)
		}
	}
	return &Library{definitions: manifest.Definitions, resources: manifest.Resources, assets: sub}, nil
}

// FontFamily devolve a família real do recurso SFNT, sem inferência pelo nome.
func (l *Library) FontFamily(resourcePath string) string {
	return l.resources[resourcePath].Family
}

// FontStyle devolve as flags físicas da face, sem sintetizar estilos do XML Java.
func (l *Library) FontStyle(resourcePath string) (bool, bool) {
	// resource conserva os metadados dos mesmos bytes usados para desenhar o texto.
	resource := l.resources[resourcePath]
	return resource.Bold, resource.Italic
}

// StateByType consulta partes compostas, como os painéis de um relatório.
func (l *Library) StateByType(componentType, name, stateName string) (State, bool) {
	// definition conserva a definição exata, sem substituir skins desconhecidas.
	definition, ok := l.definitions[componentType][name]
	if !ok {
		return State{}, false
	}
	// state representa o estado solicitado, com herança apenas dentro da mesma definição.
	state, exists := definition.States[stateName]
	if !exists {
		state, exists = definition.States["default"]
	}
	return state, exists
}

// TypeForKind relaciona a família do editor ao tipo do skin legado.
func TypeForKind(kind catalog.Kind) string {
	switch kind {
	case catalog.Button:
		return "button"
	case catalog.TextField, catalog.BarcodeInput:
		return "textField"
	case catalog.Radio:
		return "toggleButton"
	case catalog.Report:
		return "table"
	case catalog.NumericKeyboard:
		return "button"
	case catalog.Label, catalog.Title, catalog.BarcodeProcessed:
		return "label"
	case catalog.ApplicationMenu:
		return "button"
	}
	return "label"
}

// Names devolve as definições disponíveis para a família, em ordem alfabética.
func (l *Library) Names(kind catalog.Kind) []string {
	// definitions contém o mapa do tipo correspondente.
	definitions := l.definitions[TypeForKind(kind)]
	// names preserva a ordem estável exigida pelo seletor.
	names := make([]string, 0, len(definitions))
	// name percorre as definições disponíveis.
	for name := range definitions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Lookup devolve uma definição pelo tipo e nome.
func (l *Library) Lookup(kind catalog.Kind, name string) (Definition, bool) {
	// definitions contém o mapa do tipo correspondente.
	definitions := l.definitions[TypeForKind(kind)]
	// definition é a entrada consultada.
	definition, ok := definitions[name]
	return definition, ok
}

// State devolve o estado solicitado, com fallback para default.
func (l *Library) State(kind catalog.Kind, name, stateName string) (State, bool) {
	// definition é a entrada do skin escolhido.
	definition, ok := l.Lookup(kind, name)
	if !ok {
		return State{}, false
	}
	if state, exists := definition.States[stateName]; exists {
		return state, true
	}
	// state é o estado padrão quando o nome não existe.
	state, exists := definition.States["default"]
	return state, exists
}

// Asset devolve os bytes de um recurso lógico, como /azul/fundo.jpg.
func (l *Library) Asset(resourcePath string) ([]byte, error) {
	// normalized remove a barra inicial para consultar o sistema de arquivos.
	normalized := strings.TrimPrefix(resourcePath, "/")
	if normalized == "" || !fs.ValidPath(normalized) || path.Ext(normalized) == ".xml" {
		return nil, fmt.Errorf("recurso de skin vazio")
	}
	return fs.ReadFile(l.assets, normalized)
}

// HasAsset informa se o recurso lógico está disponível.
func (l *Library) HasAsset(path string) bool {
	_, err := l.Asset(path)
	return err == nil
}
