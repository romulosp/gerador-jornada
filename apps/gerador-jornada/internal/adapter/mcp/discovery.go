package mcp

import (
	"context"
	"sort"

	"gerador-jornada.local/editor/internal/adapter/buildinfo"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/skin"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// kindRequest seleciona um Kind instalado, sem tratar preset/skin como tipo novo.
type kindRequest struct {
	// Kind é validado contra o catálogo real antes de consultar qualquer metadado.
	Kind string `json:"kind"`
}

// skinRequest acrescenta paginação à família visual autorizada na consulta.
type skinRequest struct {
	// kindRequest/pageRequest são filtros, nunca autorização de leitura externa.
	kindRequest
	pageRequest
}

// targetCapabilities separa pintura/navegação suportadas de geração Java não ofertada.
type targetCapabilities struct {
	// Draft/Preview/PPTX refletem suporte instalado; PPTX é representação estática.
	Draft   bool `json:"draft"`
	Preview bool `json:"preview"`
	PPTX    bool `json:"pptx"`
	// Java permanece falso na interface MCP I1.
	Java       bool   `json:"java"`
	PPTXMode   string `json:"pptxMode"`
	JavaReason string `json:"javaReason"`
}

// catalogItem descreve somente uma família ou composto nativo inserível.
type catalogItem struct {
	// Kind/Name/Category separam tipo, nome e agrupamento da paleta.
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Category string `json:"category"`
	// Width/Height são defaults reais de criação, sujeitos ao mesmo clamp do editor.
	Width  int `json:"width"`
	Height int `json:"height"`
	// AdvancedProperties indica disponibilidade do patch específico I1, não recursos da GUI.
	AdvancedProperties bool   `json:"advancedProperties"`
	AdvancedReasonCode string `json:"advancedReasonCode,omitempty"`
	// Operations lista ferramentas instaladas para esta família, sem permissões novas implícitas.
	Operations []string `json:"operations"`
	// Targets/ImplicitActions não transformam inserção em ação de negócio automática.
	Targets         targetCapabilities `json:"targets"`
	ImplicitActions bool               `json:"implicitActions"`
}

// propertyDescriptor traduz metadados puros sem serializar regras do inspector como widgets.
type propertyDescriptor struct {
	// Name/Type/Description identificam conteúdo e finalidade.
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	// EnumValues/MaxCodePoints conservam enums de domínio e unidades textuais explícitas.
	EnumValues    []string `json:"enumValues,omitempty"`
	MaxCodePoints int      `json:"maxCodePoints,omitempty"`
}

// defaultContent oferece defaults sem inventar UUID reutilizável para opções novas.
type defaultContent struct {
	// Text/Variant/InitialValue/Format são valores da variante instalada.
	Text         *string `json:"text,omitempty"`
	Variant      *string `json:"variant,omitempty"`
	InitialValue *string `json:"initialValue,omitempty"`
	Format       *string `json:"format,omitempty"`
	// Orientation/OptionTexts/IDs descrevem opções cujo UUID nasce em cada inserção.
	Orientation      *string  `json:"orientation,omitempty"`
	OptionTexts      []string `json:"optionTexts,omitempty"`
	IDs              string   `json:"ids,omitempty"`
	SelectedOptionID *string  `json:"selectedOptionId,omitempty"`
	// Title/RowCount descrevem relatório sem transformar quantidade em patch de coleção.
	Title             *string `json:"title,omitempty"`
	RowCount          *int    `json:"rowCount,omitempty"`
	TargetComponentID *string `json:"targetComponentId,omitempty"`
	// Constructor informa defaults de composto usados na inserção, sem prometer patch avançado.
	Constructor string `json:"constructor,omitempty"`
}

// catalogDescriptor reúne defaults/capacidades da mesma versão instalada.
type catalogDescriptor struct {
	// catalogItem conserva os mesmos campos da listagem de Kinds.
	catalogItem
	// Properties/Defaults não concedem execução ou publicação.
	Properties []propertyDescriptor `json:"properties"`
	Defaults   defaultContent       `json:"defaults"`
	// DefaultAppearance aplica o legado por família mesmo sem skin explícita.
	DefaultAppearance string   `json:"defaultAppearance"`
	CommonChanges     []string `json:"commonChanges"`
}

// presetItem descreve composição instalada sem atribuir-lhe Kind serializado novo.
type presetItem struct {
	// ID/Name/Width/Height são os mesmos valores usados pela paleta manual.
	ID     string `json:"id"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// skinResource informa apenas metadados validados de bytes incorporados.
type skinResource struct {
	// Path é referência lógica interna, não caminho livre no computador do operador.
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	SizeBytes int    `json:"sizeBytes"`
	// FontFamily é extraída da face física, não inventada a partir do nome XML.
	FontFamily string `json:"fontFamily,omitempty"`
}

// skinItem separa oferta disponível e diagnóstico de definição quebrada.
type skinItem struct {
	// Name/Type não são novos Kinds de componente.
	Name       string `json:"name"`
	Type       string `json:"type"`
	Available  bool   `json:"available"`
	ReasonCode string `json:"reasonCode,omitempty"`
	// States/Resources/Missing não incluem imagem/fontes em base64 nem conteúdos pessoais.
	States    []string       `json:"states"`
	Resources []skinResource `json:"resources"`
	Missing   []string       `json:"missing,omitempty"`
}

// installedCatalog mantém a mesma ordenação da paleta: sete famílias e três compostos.
func installedCatalog() []catalogItem {
	// result adapta definições instaladas, sem antecipar novas famílias do backlog.
	result := []catalogItem{}
	for _, definition := range append(catalog.Definitions(), catalog.LegacyDefinitions()...) {
		// advanced deriva dos descritores realmente suportados pelo patch I1.
		advanced := len(catalog.PropertyDefinitions(definition.Kind)) != 0
		category, reason := "family", ""
		if !advanced {
			category, reason = "legacyCompound", "capabilityUnavailable"
		}
		result = append(result, catalogItem{Kind: string(definition.Kind), Name: definition.Name, Category: category, Width: definition.Width, Height: definition.Height, AdvancedProperties: advanced, AdvancedReasonCode: reason, Operations: []string{"componente_inserir", "componente_atualizar", "componente_mover", "componente_remover", "jornada_obter"}, Targets: targetCapabilities{Draft: true, Preview: true, PPTX: true, PPTXMode: "static", JavaReason: "Java/EFL não ofertado por MCP I1"}})
	}
	return result
}

// installedDefaults lê valores do construtor comum; IDs novos não são expostos como defaults fixos.
func installedDefaults(kind catalog.Kind) defaultContent {
	// properties usa o mesmo construtor de screen.NewComponent, sem ler material-legado em runtime.
	properties, _ := catalog.DefaultProperties(kind, project.NewID)
	switch kind {
	case catalog.Button:
		return defaultContent{Text: &properties.Button.Text, Variant: &properties.Button.Variant}
	case catalog.Label:
		return defaultContent{Text: &properties.Label.Text}
	case catalog.Title:
		return defaultContent{Text: &properties.Title.Text}
	case catalog.TextField:
		return defaultContent{InitialValue: &properties.Field.InitialValue, Format: &properties.Field.Format}
	case catalog.Radio:
		// texts descreve conteúdo inicial; identidades são geradas apenas na inserção efetiva.
		texts := []string{}
		for _, option := range properties.Radio.Options {
			texts = append(texts, option.Text)
		}
		return defaultContent{Orientation: &properties.Radio.Orientation, OptionTexts: texts, IDs: "newUUIDPerInsertion", SelectedOptionID: &properties.Radio.SelectedOptionID}
	case catalog.Report:
		// count conserva zero explícito sem emitir um ID fictício de linha.
		count := len(properties.Report.Rows)
		return defaultContent{Title: &properties.Report.Title, RowCount: &count}
	case catalog.NumericKeyboard:
		return defaultContent{TargetComponentID: &properties.Keyboard.TargetComponentID}
	default:
		return defaultContent{Constructor: "catalog.DefaultProperties", IDs: "newUUIDPerInsertion"}
	}
}

// skinItems enumera somente definições íntegras e diagnósticos pertinentes ao tipo da família.
func skinItems(kind catalog.Kind) ([]skinItem, error) {
	// library já verifica tamanho/hash dos recursos incorporados antes de oferecer nomes.
	library, err := skin.Builtin()
	if err != nil {
		return nil, &application.AuthoringError{Code: "capabilityUnavailable"}
	}
	items := []skinItem{}
	for _, name := range library.Names(kind) {
		// definition não é modificada; a resposta usa coleções independentes.
		definition, _ := library.Lookup(kind, name)
		item := skinItem{Name: name, Type: definition.Type, Available: true, States: []string{}, Resources: []skinResource{}}
		paths := map[string]bool{}
		for stateName, state := range definition.States {
			item.States = append(item.States, stateName)
			if state.ImagePath != "" {
				paths[state.ImagePath] = true
			}
			if state.Font != nil && state.Font.Type != "" {
				paths[state.Font.Type] = true
			}
		}
		sort.Strings(item.States)
		for path := range paths {
			// resource contém integridade/fonte física, sem transportar os bytes da imagem ou fonte.
			resource, found := library.ResourceInfo(path)
			if !found {
				return nil, &application.AuthoringError{Code: "capabilityUnavailable"}
			}
			item.Resources = append(item.Resources, skinResource{Path: path, SHA256: resource.SHA256, SizeBytes: resource.Size, FontFamily: resource.Family})
		}
		sort.Slice(item.Resources, func(i, j int) bool { return item.Resources[i].Path < item.Resources[j].Path })
		items = append(items, item)
	}
	for _, unavailable := range library.Unavailable() {
		if unavailable.Type == skin.TypeForKind(kind) {
			items = append(items, skinItem{Name: unavailable.Name, Type: unavailable.Type, ReasonCode: "missingResource", States: []string{}, Resources: []skinResource{}, Missing: unavailable.Missing})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

// registerDiscovery oferece consultas tipadas sem acesso a filesystem externo ou ações novas.
func registerDiscovery(server *sdk.Server, runtime *runtime) {
	// revision é fixa para recursos incorporados deste processo; cursor também é autenticado por processo.
	revision := buildinfo.Version() + "/catalog-4"
	registerTool(server, runtime, &sdk.Tool{Name: "catalogo_listar", Description: "Lista Kinds instalados; skins e presets permanecem distintos.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request pageRequest) (*sdk.CallToolResult, error) {
			page, err := paginate(runtime.pages, "catalog", revision, installedCatalog(), request)
			if err != nil {
				return nil, err
			}
			return successResult(page, ""), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "catalogo_descrever", Description: "Descreve propriedades, defaults e capacidades instaladas de um Kind.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request kindRequest) (*sdk.CallToolResult, error) {
			for _, item := range installedCatalog() {
				if item.Kind != request.Kind {
					continue
				}
				// descriptor traduz metadados puros; nenhuma regra de Fyne entra no catálogo.
				descriptor := catalogDescriptor{catalogItem: item, Properties: []propertyDescriptor{}, Defaults: installedDefaults(catalog.Kind(item.Kind)), DefaultAppearance: "legacyByFamily", CommonChanges: []string{"name", "fileId", "enabled", "bounds", "style"}}
				for _, property := range catalog.PropertyDefinitions(catalog.Kind(item.Kind)) {
					descriptor.Properties = append(descriptor.Properties, propertyDescriptor{Name: property.Name, Type: property.Type, Description: property.Description, EnumValues: property.EnumValues, MaxCodePoints: property.MaxCodePoints})
				}
				return successResult(descriptor, ""), nil
			}
			return nil, &application.AuthoringError{Code: "invalidArgument"}
		})
	registerTool(server, runtime, &sdk.Tool{Name: "catalogo_listar_presets", Description: "Lista composições instaladas e dimensões reais.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request pageRequest) (*sdk.CallToolResult, error) {
			// items espelha o caso de uso comum, sem duplicar nomes na borda.
			items := []presetItem{}
			for _, preset := range application.Presets() {
				items = append(items, presetItem{ID: preset.ID, Name: preset.Name, Width: preset.Width, Height: preset.Height})
			}
			page, err := paginate(runtime.pages, "presets", revision, items, request)
			if err != nil {
				return nil, err
			}
			return successResult(page, ""), nil
		})
	registerTool(server, runtime, &sdk.Tool{Name: "catalogo_listar_skins", Description: "Lista skins incorporadas íntegras e seus metadados de recurso.", Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}},
		func(_ context.Context, request skinRequest) (*sdk.CallToolResult, error) {
			if _, found := catalog.Lookup(catalog.Kind(request.Kind)); !found {
				return nil, &application.AuthoringError{Code: "invalidArgument"}
			}
			// items não inclui bytes de fontes/imagens nem qualquer arquivo do computador do cliente.
			items, err := skinItems(catalog.Kind(request.Kind))
			if err != nil {
				return nil, err
			}
			page, err := paginate(runtime.pages, "skins/"+request.Kind, revision, items, request.pageRequest)
			if err != nil {
				return nil, err
			}
			return successResult(page, ""), nil
		})
}
