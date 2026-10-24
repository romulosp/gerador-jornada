package catalog

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/project"
)

// testIDs devolve uma fonte determinística de UUIDs v4 canônicos.
func testIDs() func() string {
	// sequence cresce a cada chamada e usa somente dígitos hexadecimais válidos.
	sequence := 0
	return func() string {
		sequence++
		return fmt.Sprintf("00000000-0000-4000-8000-%012d", sequence)
	}
}

// emptyMenu monta um nó válido com layout explícito e lista de itens vazia.
func emptyMenu(id, layout string) *SubmenuProperties {
	return &SubmenuProperties{ID: id, Name: "Menu", Title: "Menu", Layout: layout, Items: []MenuItem{}}
}

// leafMenu monta um item folha habilitado na categoria informada.
func leafMenu(id, label, category string) MenuItem {
	return MenuItem{ID: id, Label: label, Category: category, Enabled: true}
}

// findLayout devolve o primeiro nó com o layout procurado.
func findLayout(node *SubmenuProperties, layout string) *SubmenuProperties {
	if node == nil {
		return nil
	}
	if node.Layout == layout {
		return node
	}
	for _, item := range node.Items {
		if found := findLayout(item.Submenu, layout); found != nil {
			return found
		}
	}
	return nil
}

// findLabel devolve o primeiro item com o rótulo procurado.
func findLabel(node *SubmenuProperties, label string) *MenuItem {
	if node == nil {
		return nil
	}
	for index := range node.Items {
		if node.Items[index].Label == label {
			return &node.Items[index]
		}
		if found := findLabel(node.Items[index].Submenu, label); found != nil {
			return found
		}
	}
	return nil
}

// validEmptyMenu produz a menor árvore válida, com inicial na raiz.
func validEmptyMenu(ids func() string) ApplicationMenuProperties {
	root := emptyMenu(ids(), MenuLayoutRoot)
	return ApplicationMenuProperties{Root: *root, InitialMenuID: root.ID}
}

// nestedMenu encadeia níveis de submenu e devolve o identificador do último nó.
func nestedMenu(ids func() string, levels int) (ApplicationMenuProperties, string) {
	root := emptyMenu(ids(), MenuLayoutRoot)
	node := root
	for depth := 2; depth <= levels; depth++ {
		child := emptyMenu(ids(), MenuLayoutStandard)
		node.Items = []MenuItem{{ID: ids(), Label: fmt.Sprintf("Nível %d", depth), Enabled: true, Submenu: child}}
		node = child
	}
	return ApplicationMenuProperties{Root: *root, InitialMenuID: root.ID}, node.ID
}

// wideMenu cria uma raiz com a quantidade informada de folhas.
func wideMenu(ids func() string, items int) ApplicationMenuProperties {
	root := emptyMenu(ids(), MenuLayoutRoot)
	for index := 0; index < items; index++ {
		root.Items = append(root.Items, leafMenu(ids(), fmt.Sprintf("Item %d", index+1), "SERVIÇOS"))
	}
	return ApplicationMenuProperties{Root: *root, InitialMenuID: root.ID}
}

// TestBarcodeInputDefaults confere o único campo persistido e o rótulo inicial contratado.
func TestBarcodeInputDefaults(t *testing.T) {
	// properties contém apenas rótulo, sem valores de runtime.
	properties := NewBarcodeInputProperties()
	if properties.Label != "Número do Código de Barras:" {
		t.Fatalf("rótulo padrão divergente: %q", properties.Label)
	}
	if properties.Label != BarcodeDefaultLabel {
		t.Fatalf("rótulo padrão deve usar a constante publicada: %q", BarcodeDefaultLabel)
	}
	// o tipo persistido não pode ganhar valores digitados, foco ou derivados.
	if fields := reflect.TypeOf(BarcodeInputProperties{}).NumField(); fields != 1 {
		t.Fatalf("entrada deve persistir somente o rótulo: %d campos", fields)
	}
	if err := properties.Validate(); err != nil {
		t.Fatalf("rótulo padrão recusado: %v", err)
	}
	// rótulo Unicode e vazio continuam válidos, pois são conteúdo autoral.
	for _, label := range []string{"", "Número do Código de Barras: ação"} {
		properties.Label = label
		if err := properties.Validate(); err != nil {
			t.Fatalf("rótulo %q recusado: %v", label, err)
		}
	}
}

// TestLegacyDefinitions fixa as dimensões e a separação das três famílias legadas.
func TestLegacyDefinitions(t *testing.T) {
	// want associa cada composto à definição contratada.
	want := map[Kind]Definition{
		BarcodeInput:     {Kind: BarcodeInput, Name: "Código de barras", Width: 610, Height: 60},
		BarcodeProcessed: {Kind: BarcodeProcessed, Name: "Código de barras processado", Width: 530, Height: 120},
		ApplicationMenu:  {Kind: ApplicationMenu, Name: "Menu da aplicação", Width: 800, Height: 600},
	}
	for kind, definition := range want {
		got, ok := Lookup(kind)
		if !ok || got != definition {
			t.Fatalf("definição de %s divergente: %+v", kind, got)
		}
	}
	// os compostos ficam em grupo próprio, sem inflar as sete famílias base.
	if len(Definitions()) != 7 || len(LegacyDefinitions()) != 3 {
		t.Fatalf("grupos do catálogo divergentes: %d/%d", len(Definitions()), len(LegacyDefinitions()))
	}
}

// TestBarcodeProcessedDefaults confere um bloco de título vazio com CNPJ e Nome vazios.
func TestBarcodeProcessedDefaults(t *testing.T) {
	ids := testIDs()
	// properties nasce com exatamente um bloco e duas linhas ordenadas.
	properties := NewBarcodeProcessedProperties(ids)
	if len(properties.Blocks) != 1 {
		t.Fatalf("uma nova ficha deve ter um bloco: %d", len(properties.Blocks))
	}
	block := properties.Blocks[0]
	if block.Title != "" {
		t.Fatalf("título inicial deve ser vazio: %q", block.Title)
	}
	if len(block.Rows) != 2 {
		t.Fatalf("um novo bloco deve ter duas linhas: %d", len(block.Rows))
	}
	if block.Rows[0].LabelKind != BarcodeLabelCNPJ || block.Rows[0].Value != "" {
		t.Fatalf("primeira linha deve ser CNPJ vazio: %+v", block.Rows[0])
	}
	if block.Rows[1].LabelKind != BarcodeLabelName || block.Rows[1].Value != "" {
		t.Fatalf("segunda linha deve ser Nome vazio: %+v", block.Rows[1])
	}
	// idsList confirma a ordem bloco, linha, linha e a unicidade global.
	idsList := CollectBarcodeProcessedIDs(*properties)
	want := []string{block.ID, block.Rows[0].ID, block.Rows[1].ID}
	if !slices.Equal(idsList, want) {
		t.Fatalf("ordem de IDs divergente: %v", idsList)
	}
	if err := properties.Validate(); err != nil {
		t.Fatalf("defaults recusados: %v", err)
	}
	// row e newBlock usam identidades novas e os defaults documentados.
	row := NewBarcodeRow(ids)
	if row.LabelKind != BarcodeLabelName || row.Value != "" || !project.ValidID(row.ID) {
		t.Fatalf("nova linha inválida: %+v", row)
	}
	fresh := NewBarcodeBlock(ids)
	if len(fresh.Rows) != 1 || fresh.Rows[0].LabelKind != BarcodeLabelName || !project.ValidID(fresh.ID) {
		t.Fatalf("novo bloco inválido: %+v", fresh)
	}
}

// TestApplicationMenuDefaults confere raiz, duas colunas, inicial e legendas do Bolão.
func TestApplicationMenuDefaults(t *testing.T) {
	ids := testIDs()
	// properties contém a árvore seed documentada.
	properties := NewApplicationMenuProperties(ids)
	if properties.InitialMenuID != properties.Root.ID {
		t.Fatalf("menu inicial deve ser a raiz: %q != %q", properties.InitialMenuID, properties.Root.ID)
	}
	if properties.Root.Layout != MenuLayoutRoot || properties.Root.Title != "" {
		t.Fatalf("raiz com layout ou título incorretos: %+v", properties.Root)
	}
	// services e games preservam as duas colunas da raiz.
	services, games := []MenuItem{}, []MenuItem{}
	for _, item := range properties.Root.Items {
		if item.Category == "JOGOS" {
			games = append(games, item)
		} else {
			services = append(services, item)
		}
	}
	if len(services) != 12 || len(games) != 10 {
		t.Fatalf("raiz deve ter 12 serviços e 10 jogos: %d/%d", len(services), len(games))
	}
	// wantServices e wantGames fixam os rótulos e a ordem do contrato.
	wantServices := []string{"  ARRECADA BOLETO  DE COBRANÇA", "ARRECADAÇÃO DE CONVÊNIOS", "BENEFÍCIOS SOCIAIS", "OUTROS SERVIÇOS FINANCEIROS", "ATENDIMENTO", "ÚLTIMOS RESULTADOS", "APOIO", "BOLÃO", "RELATÓRIOS", "PAGA PREMIADOS", "MEIOS DE PAGAMENTO", "OFERTAR"}
	wantGames := []string{"MEGA-SENA", "LOTOFÁCIL", "QUINA", "LOTOMANIA", "LOTECA", "DUPLA-SENA", "TIMEMANIA", "DIA DE SORTE", "SUPER SETE", "MILIONÁRIA"}
	// index percorre as duas colunas para comparar cada rótulo.
	for index, label := range wantServices {
		if services[index].Label != label {
			t.Fatalf("serviço %d divergente: %q", index+1, services[index].Label)
		}
	}
	for index, label := range wantGames {
		if games[index].Label != label {
			t.Fatalf("jogo %d divergente: %q", index+1, games[index].Label)
		}
	}
	if err := properties.Validate(); err != nil {
		t.Fatalf("seed recusado: %v", err)
	}
	// initial resolve a raiz por identidade, nunca por título.
	if initial := properties.Initial(); initial == nil || initial.ID != properties.Root.ID {
		t.Fatalf("menu inicial não resolve a raiz: %+v", initial)
	}
	// digital contém os três botões com legenda do layout bolao.
	digital := findLayout(&properties.Root, MenuLayoutBolao)
	if digital == nil {
		t.Fatal("BOLÕES DIGITAIS com layout bolao não encontrado")
	}
	if len(digital.Items) != 3 {
		t.Fatalf("BOLÕES DIGITAIS deve ter três itens: %d", len(digital.Items))
	}
	wantLegends := []string{
		"Visualizar cotas digitais para venda na lotérica.",
		"Consulta bolões de concursos abertos para comercialização.",
		"Consulta bolões de concursos encerrados e não prescritos.",
	}
	wantLabels := []string{"VENDER COTAS", "BOLÕES MARKETPLACE", "BOLÕES CONCLUÍDOS"}
	for index := range digital.Items {
		if digital.Items[index].Legend != wantLegends[index] {
			t.Fatalf("legenda %d divergente: %q", index+1, digital.Items[index].Legend)
		}
		if digital.Items[index].Label != wantLabels[index] {
			t.Fatalf("rótulo %d divergente: %q", index+1, digital.Items[index].Label)
		}
		if !digital.Items[index].Enabled || digital.Items[index].Category != "SERVIÇOS" {
			t.Fatalf("item %d do bolão deve estar habilitado em SERVIÇOS: %+v", index+1, digital.Items[index])
		}
	}
	// bolao é o oitavo item root e contém BOLÕES DIGITAIS como último filho.
	bolao := findLabel(&properties.Root, "BOLÃO")
	if bolao == nil || bolao.Submenu == nil || bolao.Submenu.Layout != MenuLayoutStandard {
		t.Fatalf("BOLÃO ordinário deve usar layout standard: %+v", bolao)
	}
	if len(bolao.Submenu.Items) != 10 || bolao.Submenu.Items[9].Submenu != digital {
		t.Fatalf("BOLÕES DIGITAIS deve ser o último filho de BOLÃO: %+v", bolao.Submenu.Items)
	}
}

// TestBarcodeProcessedValidation recusa estruturas vazias, enum desconhecido e IDs repetidos.
func TestBarcodeProcessedValidation(t *testing.T) {
	// newProperties cria uma ficha padrão com identidades determinísticas e exclusivas.
	newProperties := func() *BarcodeProcessedProperties { return NewBarcodeProcessedProperties(testIDs()) }
	t.Run("defaults válidos", func(t *testing.T) {
		if err := newProperties().Validate(); err != nil {
			t.Fatalf("defaults recusados: %v", err)
		}
	})
	t.Run("sem blocos", func(t *testing.T) {
		properties := BarcodeProcessedProperties{}
		if err := properties.Validate(); err == nil {
			t.Fatal("ficha sem blocos foi aceita")
		}
	})
	t.Run("sem linhas", func(t *testing.T) {
		ids := testIDs()
		properties := BarcodeProcessedProperties{Blocks: []BarcodeBlock{{ID: ids(), Rows: []BarcodeRow{}}}}
		if err := properties.Validate(); err == nil {
			t.Fatal("bloco sem linhas foi aceito")
		}
	})
	t.Run("rótulo desconhecido", func(t *testing.T) {
		properties := newProperties()
		properties.Blocks[0].Rows[0].LabelKind = "rg"
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "rótulo desconhecido") {
			t.Fatalf("enum desconhecido aceito ou erro genérico: %v", err)
		}
	})
	t.Run("UUID repetido", func(t *testing.T) {
		properties := newProperties()
		properties.Blocks[0].Rows[1].ID = properties.Blocks[0].Rows[0].ID
		if err := properties.Validate(); err == nil {
			t.Fatal("linhas com UUID repetido foram aceitas")
		}
	})
	t.Run("UUID inválido", func(t *testing.T) {
		properties := newProperties()
		properties.Blocks[0].ID = "bloco-1"
		if err := properties.Validate(); err == nil {
			t.Fatal("UUID inválido foi aceito")
		}
	})
	t.Run("bloco repetido com linha", func(t *testing.T) {
		properties := newProperties()
		properties.Blocks = append(properties.Blocks, BarcodeBlock{ID: properties.Blocks[0].Rows[0].ID, Rows: []BarcodeRow{NewBarcodeRow(testIDs())}})
		if err := properties.Validate(); err == nil {
			t.Fatal("UUID compartilhado entre bloco e linha foi aceito")
		}
	})
}

// TestApplicationMenuValidation cobre vazio, limites, identidades, ciclo e referências.
func TestApplicationMenuValidation(t *testing.T) {
	t.Run("menus vazios válidos", func(t *testing.T) {
		ids := testIDs()
		properties := validEmptyMenu(ids)
		if err := properties.Validate(); err != nil {
			t.Fatalf("raiz vazia recusada: %v", err)
		}
		// submenu vazio também é uma composição válida durante a autoria.
		child := emptyMenu(ids(), MenuLayoutStandard)
		properties.Root.Items = []MenuItem{{ID: ids(), Label: "Vazio", Enabled: true, Submenu: child}}
		if err := properties.Validate(); err != nil {
			t.Fatalf("submenu vazio recusado: %v", err)
		}
		// folha habilitada e folha desabilitada convivem sem inventar navegação.
		properties.Root.Items = append(properties.Root.Items, leafMenu(ids(), "Desabilitado", "SERVIÇOS"))
		properties.Root.Items[len(properties.Root.Items)-1].Enabled = false
		if err := properties.Validate(); err != nil {
			t.Fatalf("folha desabilitada recusada: %v", err)
		}
	})
	t.Run("profundidade 16 aceita", func(t *testing.T) {
		properties, lastID := nestedMenu(testIDs(), MenuMaxDepth)
		if err := properties.Validate(); err != nil {
			t.Fatalf("profundidade %d recusada: %v", MenuMaxDepth, err)
		}
		if path := properties.PathTo(lastID); len(path) != MenuMaxDepth {
			t.Fatalf("caminho deve conter %d níveis: %d", MenuMaxDepth, len(path))
		}
	})
	t.Run("profundidade 17 recusada", func(t *testing.T) {
		properties, _ := nestedMenu(testIDs(), MenuMaxDepth+1)
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "profundidade") {
			t.Fatalf("profundidade %d foi aceita: %v", MenuMaxDepth+1, err)
		}
	})
	t.Run("1000 itens aceitos", func(t *testing.T) {
		properties := wideMenu(testIDs(), MenuMaxItems)
		if err := properties.Validate(); err != nil {
			t.Fatalf("%d itens recusados: %v", MenuMaxItems, err)
		}
	})
	t.Run("1001 itens recusados", func(t *testing.T) {
		properties := wideMenu(testIDs(), MenuMaxItems+1)
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "itens") {
			t.Fatalf("%d itens foram aceitos: %v", MenuMaxItems+1, err)
		}
	})
	t.Run("UUID repetido", func(t *testing.T) {
		ids := testIDs()
		properties := validEmptyMenu(ids)
		properties.Root.Items = []MenuItem{leafMenu(ids(), "A", "SERVIÇOS"), leafMenu(ids(), "B", "SERVIÇOS")}
		properties.Root.Items[1].ID = properties.Root.Items[0].ID
		if err := properties.Validate(); err == nil {
			t.Fatal("itens com UUID repetido foram aceitos")
		}
	})
	t.Run("UUID repetido entre nó e item", func(t *testing.T) {
		ids := testIDs()
		properties := validEmptyMenu(ids)
		properties.Root.Items = []MenuItem{leafMenu(ids(), "A", "SERVIÇOS")}
		properties.Root.Items[0].ID = properties.Root.ID
		if err := properties.Validate(); err == nil {
			t.Fatal("item reutilizou o UUID da raiz")
		}
	})
	t.Run("submenu com destino", func(t *testing.T) {
		ids := testIDs()
		properties := validEmptyMenu(ids)
		properties.Root.Items = []MenuItem{{
			ID: ids(), Label: "Ramo", Enabled: true,
			Submenu: emptyMenu(ids(), MenuLayoutStandard), TargetScreenID: project.NewID(),
		}}
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "submenu e destino") {
			t.Fatalf("item com submenu e destino foi aceito: %v", err)
		}
	})
	t.Run("inicial inexistente", func(t *testing.T) {
		properties := validEmptyMenu(testIDs())
		properties.InitialMenuID = project.NewID()
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "menu inicial") {
			t.Fatalf("inicial inexistente foi aceito: %v", err)
		}
	})
	t.Run("ciclo de ponteiro", func(t *testing.T) {
		ids := testIDs()
		root := emptyMenu(ids(), MenuLayoutRoot)
		// cycle aponta para si mesmo por um item descendente.
		cycle := emptyMenu(ids(), MenuLayoutStandard)
		cycle.Items = []MenuItem{{ID: ids(), Label: "Ciclo", Enabled: true, Submenu: cycle}}
		root.Items = []MenuItem{{ID: ids(), Label: "Entrada", Enabled: true, Submenu: cycle}}
		properties := ApplicationMenuProperties{Root: *root, InitialMenuID: root.ID}
		if err := properties.Validate(); err == nil || !strings.Contains(err.Error(), "ciclo") {
			t.Fatalf("ciclo de submenu foi aceito: %v", err)
		}
	})
	t.Run("layouts inválidos", func(t *testing.T) {
		ids := testIDs()
		// raiz precisa declarar layout root.
		root := emptyMenu(ids(), MenuLayoutStandard)
		if err := (ApplicationMenuProperties{Root: *root, InitialMenuID: root.ID}).Validate(); err == nil {
			t.Fatal("raiz com layout standard foi aceita")
		}
		// descendentes não podem reutilizar o layout root.
		properties := validEmptyMenu(ids)
		properties.Root.Items = []MenuItem{{ID: ids(), Label: "Ramo", Enabled: true, Submenu: emptyMenu(ids(), MenuLayoutRoot)}}
		if err := properties.Validate(); err == nil {
			t.Fatal("descendente com layout root foi aceito")
		}
	})
	t.Run("initial em descendente", func(t *testing.T) {
		ids := testIDs()
		properties := validEmptyMenu(ids)
		child := emptyMenu(ids(), MenuLayoutStandard)
		properties.Root.Items = []MenuItem{{ID: ids(), Label: "Ramo", Enabled: true, Submenu: child}}
		properties.InitialMenuID = child.ID
		if err := properties.Validate(); err != nil {
			t.Fatalf("inicial em descendente recusado: %v", err)
		}
		if path := properties.PathTo(child.ID); len(path) != 2 {
			t.Fatalf("caminho até o descendente deve ter dois níveis: %d", len(path))
		}
	})
}

// TestCompoundPropertiesValidate confere a união discriminada dos três tipos novos.
func TestCompoundPropertiesValidate(t *testing.T) {
	ids := testIDs()
	// cases associa cada união ao Kind que deve aceitá-la.
	cases := []struct {
		name       string
		properties Properties
		kind       Kind
		valid      bool
	}{
		{"barcodeInput", Properties{BarcodeInput: NewBarcodeInputProperties()}, BarcodeInput, true},
		{"barcodeProcessed", Properties{BarcodeProcessed: NewBarcodeProcessedProperties(ids)}, BarcodeProcessed, true},
		{"applicationMenu", Properties{Menu: NewApplicationMenuProperties(ids)}, ApplicationMenu, true},
		{"barcodeInput com menu", Properties{BarcodeInput: NewBarcodeInputProperties(), Menu: NewApplicationMenuProperties(ids)}, BarcodeInput, false},
		{"barcodeInput como menu", Properties{BarcodeInput: NewBarcodeInputProperties()}, ApplicationMenu, false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.properties.Validate(testCase.kind)
			if testCase.valid && err != nil {
				t.Fatalf("união válida recusada: %v", err)
			}
			if !testCase.valid && err == nil {
				t.Fatal("união inválida foi aceita")
			}
		})
	}
}

// TestPropertiesCopyIsolatesCompounds garante que a cópia não compartilha slices nem ponteiros.
func TestPropertiesCopyIsolatesCompounds(t *testing.T) {
	t.Run("barcodeInput", func(t *testing.T) {
		source := Properties{BarcodeInput: NewBarcodeInputProperties()}
		copied := source.Copy()
		if copied.BarcodeInput == source.BarcodeInput {
			t.Fatal("ponteiro da entrada foi compartilhado")
		}
		copied.BarcodeInput.Label = "Cópia"
		if source.BarcodeInput.Label != BarcodeDefaultLabel {
			t.Fatalf("original alterado: %q", source.BarcodeInput.Label)
		}
	})
	t.Run("barcodeProcessed", func(t *testing.T) {
		ids := testIDs()
		source := Properties{BarcodeProcessed: NewBarcodeProcessedProperties(ids)}
		// originalID preserva a identidade do primeiro bloco usada na comparação.
		originalID := source.BarcodeProcessed.Blocks[0].ID
		copied := source.Copy()
		if copied.BarcodeProcessed == source.BarcodeProcessed {
			t.Fatal("ponteiro da ficha foi compartilhado")
		}
		copied.BarcodeProcessed.Blocks[0].Title = "Cópia"
		copied.BarcodeProcessed.Blocks[0].Rows[0].Value = "novo valor"
		copied.BarcodeProcessed.Blocks[0].Rows[0].ID = ""
		copied.BarcodeProcessed.Blocks = append(copied.BarcodeProcessed.Blocks, BarcodeBlock{ID: ids(), Rows: []BarcodeRow{NewBarcodeRow(ids)}})
		if len(source.BarcodeProcessed.Blocks) != 1 {
			t.Fatal("slice de blocos foi compartilhado")
		}
		if source.BarcodeProcessed.Blocks[0].Title != "" || source.BarcodeProcessed.Blocks[0].Rows[0].Value != "" {
			t.Fatalf("conteúdo original alterado: %+v", source.BarcodeProcessed.Blocks[0])
		}
		if source.BarcodeProcessed.Blocks[0].ID != originalID || source.BarcodeProcessed.Blocks[0].Rows[0].ID == "" {
			t.Fatal("identidades originais alteradas pela cópia")
		}
		// a linha copiada é um valor independente, não um alias do slice original.
		if &copied.BarcodeProcessed.Blocks[0].Rows[0] == &source.BarcodeProcessed.Blocks[0].Rows[0] {
			t.Fatal("linhas compartilham o mesmo endereço")
		}
	})
	t.Run("menu", func(t *testing.T) {
		ids := testIDs()
		source := Properties{Menu: NewApplicationMenuProperties(ids)}
		// originalRootID e originalInitial fixam a configuração da origem.
		originalRootID := source.Menu.Root.ID
		originalInitial := source.Menu.InitialMenuID
		copied := source.Copy()
		if copied.Menu == source.Menu || &copied.Menu.Root == &source.Menu.Root {
			t.Fatal("ponteiro da árvore foi compartilhado")
		}
		if copied.Menu.Root.Items[1].Submenu == source.Menu.Root.Items[1].Submenu {
			t.Fatal("ponteiro de submenu foi compartilhado")
		}
		copied.Menu.Root.Name = "Outro nome"
		copied.Menu.Root.Items[0].Label = "ALTERADO"
		copied.Menu.Root.Items[1].Submenu.Name = "OUTRO"
		copied.Menu.Root.Items[1].Submenu.Items = append(copied.Menu.Root.Items[1].Submenu.Items, leafMenu(ids(), "Novo", "SERVIÇOS"))
		copied.Menu.Root.Items = append(copied.Menu.Root.Items, leafMenu(ids(), "Extra", "SERVIÇOS"))
		copied.Menu.InitialMenuID = copied.Menu.Root.Items[0].ID
		if source.Menu.Root.Name != "Principal" || source.Menu.Root.Items[0].Label != "MEGA-SENA" {
			t.Fatalf("nó ou item original alterado: %+v", source.Menu.Root.Items[0])
		}
		if source.Menu.Root.Items[1].Submenu.Name != "LOTOFÁCIL" || len(source.Menu.Root.Items[1].Submenu.Items) != 2 {
			t.Fatalf("submenu original alterado: %+v", source.Menu.Root.Items[1].Submenu)
		}
		if len(source.Menu.Root.Items) != 22 || source.Menu.Root.ID != originalRootID || source.Menu.InitialMenuID != originalInitial {
			t.Fatal("estrutura ou escolha inicial original alterada")
		}
	})
}

// TestRegenerateCompoundIDs troca identidades internas e remapeia o menu inicial.
func TestRegenerateCompoundIDs(t *testing.T) {
	t.Run("barcode processado", func(t *testing.T) {
		ids := testIDs()
		properties := Properties{BarcodeProcessed: NewBarcodeProcessedProperties(ids)}
		oldIDs := properties.CompoundIDs()
		properties.RegenerateCompoundIDs(ids)
		newIDs := properties.CompoundIDs()
		if len(newIDs) != len(oldIDs) {
			t.Fatalf("quantidade de identidades mudou: %d != %d", len(newIDs), len(oldIDs))
		}
		// oldSet comprova que nenhuma identidade antiga foi reutilizada.
		oldSet := map[string]bool{}
		for _, id := range oldIDs {
			oldSet[id] = true
		}
		for _, id := range newIDs {
			if oldSet[id] {
				t.Fatalf("identidade antiga reutilizada: %s", id)
			}
			if !project.ValidID(id) {
				t.Fatalf("identidade nova inválida: %s", id)
			}
		}
		if err := properties.Validate(BarcodeProcessed); err != nil {
			t.Fatalf("ficha regenerada inválida: %v", err)
		}
		// textos e ordem continuam os mesmos após a regeneração.
		block := properties.BarcodeProcessed.Blocks[0]
		if block.Rows[0].LabelKind != BarcodeLabelCNPJ || block.Rows[1].LabelKind != BarcodeLabelName {
			t.Fatalf("rótulos alterados na regeneração: %+v", block.Rows)
		}
	})
	t.Run("menu remapeia inicial", func(t *testing.T) {
		ids := testIDs()
		menu := NewApplicationMenuProperties(ids)
		digital := findLayout(&menu.Root, MenuLayoutBolao)
		if digital == nil {
			t.Fatal("BOLÕES DIGITAIS não encontrado")
		}
		// initialTarget aponta para um descendente de layout bolao.
		initialTarget := digital.ID
		menu.InitialMenuID = initialTarget
		properties := Properties{Menu: menu}
		oldRootID := properties.Menu.Root.ID
		oldIDs := properties.CompoundIDs()
		properties.RegenerateCompoundIDs(ids)
		newIDs := properties.CompoundIDs()
		if len(newIDs) != len(oldIDs) {
			t.Fatalf("quantidade de identidades mudou: %d != %d", len(newIDs), len(oldIDs))
		}
		// oldSet reúne raiz, submenus e itens para provar os três níveis regenerados.
		oldSet := map[string]bool{}
		for _, id := range oldIDs {
			oldSet[id] = true
		}
		for _, id := range newIDs {
			if oldSet[id] {
				t.Fatalf("identidade antiga reutilizada: %s", id)
			}
		}
		if properties.Menu.Root.ID == oldRootID || properties.Menu.InitialMenuID == initialTarget {
			t.Fatal("raiz ou inicial não foram regeneradas")
		}
		// initial resolve o mesmo nó lógico, agora com identidade nova.
		initial := properties.Menu.Initial()
		if initial == nil || initial.ID != properties.Menu.InitialMenuID || initial.Name != "BOLÕES DIGITAIS" {
			t.Fatalf("inicial não foi remapeado para o nó equivalente: %+v", initial)
		}
		if err := properties.Validate(ApplicationMenu); err != nil {
			t.Fatalf("menu regenerado inválido: %v", err)
		}
	})
}

// barcodeF01Digitized contém a linha digitável de referência do contrato.
const barcodeF01Digitized = "00190500954014481606906809350314337370000000100"

// barcodeF01Read contém o código de barras lido correspondente a F01.
const barcodeF01Read = "00193373700000001000500940144816060680935031"

// barcodeF01Fields contém os oito campos visuais derivados de F01.
var barcodeF01Fields = []string{"00190", "50095", "40144", "816069", "06809", "350314", "3", "37370000000100"}

// TestBarcodeFieldsFixtures confere o recorte dos oito campos e seus limites.
func TestBarcodeFieldsFixtures(t *testing.T) {
	// fields preserva a ordem visual dos oito campos de F01.
	fields, err := BarcodeFieldsFromDigitized(barcodeF01Digitized)
	if err != nil {
		t.Fatalf("recorte de F01 falhou: %v", err)
	}
	if !slices.Equal(fields, barcodeF01Fields) {
		t.Fatalf("campos de F01 divergentes: %v", fields)
	}
	// limits repete os limites contratados na ordem dos campos.
	limits := []int{5, 5, 5, 6, 5, 6, 1, 14}
	for index, field := range fields {
		if len(field) != limits[index] || len(field) != BarcodeFieldLimits[index] {
			t.Fatalf("campo %d com comprimento divergente: %q", index+1, field)
		}
	}
	// fromRead converte o código lido nos mesmos oito valores.
	fromRead, err := BarcodeFieldsFromRead(barcodeF01Read)
	if err != nil {
		t.Fatalf("recorte do código lido falhou: %v", err)
	}
	if !slices.Equal(fromRead, fields) {
		t.Fatalf("código lido e linha digitável divergem: %v", fromRead)
	}
	// f02Read e f02Digitized comprovam o mesmo recorte no banco 104.
	f02Read := "10498100000000001001234560000100010000000003"
	f02Digitized := "10491234566000010001000000000034810000000000100"
	// readFields e digitizedFields precisam coincidir posição a posição.
	readFields, err := BarcodeFieldsFromRead(f02Read)
	if err != nil {
		t.Fatalf("recorte de F02 lido falhou: %v", err)
	}
	digitizedFields, err := BarcodeFieldsFromDigitized(f02Digitized)
	if err != nil {
		t.Fatalf("recorte de F02 digitado falhou: %v", err)
	}
	if !slices.Equal(readFields, digitizedFields) || len(readFields) != 8 {
		t.Fatalf("F02 com representações divergentes: %v != %v", readFields, digitizedFields)
	}
	// uma linha de 44 posições não pode ser usada como linha digitável.
	if _, err := BarcodeFieldsFromDigitized(barcodeF01Read); err == nil {
		t.Fatal("código lido de 44 posições foi aceito como linha digitável")
	}
}

// TestParseTypedFieldsFixtures valida a digitação dos oito campos de F01.
func TestParseTypedFieldsFixtures(t *testing.T) {
	// typed contém os oito campos na ordem visual de F01.
	typed := barcodeF01Fields
	barcode, err := ParseTypedFields(typed)
	if err != nil {
		t.Fatalf("campos digitados de F01 falharam: %v", err)
	}
	// reference é o mesmo código aceito pela captura integral.
	reference, err := ParseBarcode(barcodeF01Digitized)
	if err != nil {
		t.Fatalf("captura integral de F01 falhou: %v", err)
	}
	if barcode != reference {
		t.Fatalf("digitação e captura divergem: %+v != %+v", barcode, reference)
	}
	// invalid reúne entradas estruturais recusadas antes da aritmética.
	invalid := map[string][]string{
		"nenhum campo":           nil,
		"sete campos":            typed[:7],
		"nove campos":            append(slices.Clone(typed), "0"),
		"campo acima do limite":  slices.Clone(typed),
		"caractere não numérico": slices.Clone(typed),
		"entrada incompleta":     slices.Clone(typed),
	}
	invalid["campo acima do limite"][0] = "001900"
	invalid["caractere não numérico"][0] = "0019a"
	invalid["entrada incompleta"][7] = "373700000001"
	for name, fields := range invalid {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseTypedFields(fields); err == nil {
				t.Fatalf("entrada %q foi aceita", name)
			}
		})
	}
}

// TestValidateTypedBlockFixtures cobre F03 e F04 nos três blocos lógicos.
func TestValidateTypedBlockFixtures(t *testing.T) {
	// fields contém os oito campos válidos de F01.
	fields := barcodeF01Fields
	for block := 1; block <= 3; block++ {
		if err := ValidateTypedBlock(fields, block); err != nil {
			t.Fatalf("bloco %d válido recusado: %v", block, err)
		}
	}
	// F03 troca o DV local de cada par e exige erro apenas no bloco adulterado.
	for block := 1; block <= 3; block++ {
		mutated := slices.Clone(fields)
		// pairEnd aponta para o último campo visual do par, onde mora o DV local.
		pairEnd := 2*(block-1) + 1
		mutated[pairEnd] = nextDigit(mutated[pairEnd])
		if err := ValidateTypedBlock(mutated, block); err == nil || !strings.Contains(err.Error(), "inválido") {
			t.Fatalf("F03 do bloco %d foi aceito: %v", block, err)
		}
		for other := 1; other <= 3; other++ {
			if other == block {
				continue
			}
			if err := ValidateTypedBlock(mutated, other); err != nil {
				t.Fatalf("bloco %d válido foi contaminado pelo bloco %d: %v", other, block, err)
			}
		}
		if _, err := ParseTypedFields(mutated); err == nil {
			t.Fatalf("F03 do bloco %d foi aceito na digitação completa", block)
		}
	}
	// pares incompletos, acima do limite ou fora do intervalo são recusados.
	short := slices.Clone(fields)
	short[2] = ""
	if err := ValidateTypedBlock(short, 2); err == nil {
		t.Fatal("bloco incompleto foi aceito")
	}
	long := slices.Clone(fields)
	long[4] = "068090"
	if err := ValidateTypedBlock(long, 3); err == nil {
		t.Fatal("campo acima do limite foi aceito")
	}
	for _, block := range []int{0, 4, -1} {
		if err := ValidateTypedBlock(fields, block); err == nil {
			t.Fatalf("bloco %d fora do intervalo foi aceito", block)
		}
	}
	if err := ValidateTypedBlock(fields[:7], 1); err == nil {
		t.Fatal("lista com sete campos foi aceita")
	}
}

// TestValidateTypedBlockRejectsGeneralDV reproduz F04: DVs locais válidos e geral adulterado.
func TestValidateTypedBlockRejectsGeneralDV(t *testing.T) {
	// mutatedRead troca somente a posição do DV geral no código lido.
	mutatedRead := barcodeF01Read[:4] + "4" + barcodeF01Read[5:]
	if mutatedRead == barcodeF01Read {
		t.Fatal("mutação do DV geral sem efeito")
	}
	// digitized converte estruturalmente sem revalidar a aritmética.
	digitized := readToDigitized(mutatedRead)
	fields, err := BarcodeFieldsFromDigitized(digitized)
	if err != nil {
		t.Fatalf("recorte do F04 falhou: %v", err)
	}
	// os três DVs locais permanecem válidos, pois o DV geral não compõe os pares.
	for block := 1; block <= 3; block++ {
		if err := ValidateTypedBlock(fields, block); err != nil {
			t.Fatalf("F04 contaminou o bloco %d: %v", block, err)
		}
	}
	if _, err := ParseTypedFields(fields); err == nil {
		t.Fatal("F04 com DV geral inválido foi aceito")
	}
	if _, err := ParseBarcode(mutatedRead); err == nil {
		t.Fatal("F04 foi aceito na captura integral")
	}
}

// TestCompleteLastField confere o preenchimento de zeros à direita do último campo.
func TestCompleteLastField(t *testing.T) {
	fields := slices.Clone(barcodeF01Fields)
	// trimmed remove os dois zeros finais do fator/valor de F01.
	trimmed := slices.Clone(fields)
	trimmed[7] = "373700000001"
	completed := CompleteLastField(trimmed)
	if completed[7] != fields[7] {
		t.Fatalf("preenchimento divergente: %q", completed[7])
	}
	if trimmed[7] != "373700000001" {
		t.Fatal("entrada original foi alterada pelo preenchimento")
	}
	// barcode comprova que a conclusão explícita produz o mesmo código válido.
	barcode, err := ParseTypedFields(completed)
	if err != nil {
		t.Fatalf("código completado foi recusado: %v", err)
	}
	reference, err := ParseBarcode(barcodeF01Digitized)
	if err != nil {
		t.Fatalf("captura integral de F01 falhou: %v", err)
	}
	if barcode != reference {
		t.Fatalf("código completado diverge do original: %+v", barcode)
	}
	// campo vazio recebe exatamente 14 zeros e continua exigindo o DV geral.
	empty := slices.Clone(fields)
	empty[7] = ""
	completedEmpty := CompleteLastField(empty)
	if completedEmpty[7] != strings.Repeat("0", 14) {
		t.Fatalf("campo vazio preenchido incorretamente: %q", completedEmpty[7])
	}
	if _, err := ParseTypedFields(completedEmpty); err == nil {
		t.Fatal("conclusão com campo vazio foi aceita sem validar o DV geral")
	}
	// listas fora do formato de oito campos são devolvidas sem alteração.
	for _, short := range [][]string{nil, {"123"}} {
		unchanged := CompleteLastField(short)
		if !slices.Equal(unchanged, short) {
			t.Fatalf("lista de %d campos foi alterada: %v", len(short), unchanged)
		}
	}
}

// nextDigit troca o último algarismo de um campo por outro sempre diferente.
func nextDigit(field string) string {
	if field == "" {
		return "0"
	}
	// last contém o algarismo atual, que será deslocado ciclicamente.
	last := field[len(field)-1]
	next := byte('0' + (last-'0'+1)%10)
	return field[:len(field)-1] + string(next)
}
