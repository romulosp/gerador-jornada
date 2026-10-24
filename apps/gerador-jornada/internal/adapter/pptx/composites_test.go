package pptx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"gerador-jornada.local/editor/internal/skin"
)

// compoundDocument monta uma tela com os três compostos e textos Unicode explícitos.
func compoundDocument(t *testing.T) (application.Document, screen.Component, screen.Component, screen.Component) {
	t.Helper()
	document := application.NewDocument("Compostos 031/032/033 — ação")
	page := &document.Screens[0]
	// input usa o rótulo persistido com acento e a geometria de referência.
	input, err := screen.NewComponent(catalog.BarcodeInput, "codigoBarras", 20, 20)
	if err != nil {
		t.Fatalf("entrada inválida: %v", err)
	}
	input.Bounds = screen.Bounds{X: 20, Y: 20, Width: 610, Height: 60}
	input.Properties.BarcodeInput.Label = "Número do Código de Barras: ação"
	// processed contém um bloco com título e as duas linhas padrão preenchidas.
	processed, err := screen.NewComponent(catalog.BarcodeProcessed, "codigoProcessado", 20, 100)
	if err != nil {
		t.Fatalf("ficha inválida: %v", err)
	}
	processed.Bounds = screen.Bounds{X: 20, Y: 100, Width: 530, Height: 120}
	block := &processed.Properties.BarcodeProcessed.Blocks[0]
	block.Title = "Dados do boleto"
	block.Rows[0].Value = "12.345.678/0001-90"
	block.Rows[1].LabelKind = catalog.BarcodeLabelTradeName
	block.Rows[1].Value = "Ação Comércio Ltda"
	// menu usa a árvore seed inteira na geometria de referência.
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menuAplicacao", 0, 0)
	if err != nil {
		t.Fatalf("menu inválido: %v", err)
	}
	menu.Bounds = screen.Bounds{X: 0, Y: 0, Width: 800, Height: 600}
	page.Components = append(page.Components, input, processed, menu)
	return document, input, processed, menu
}

// componentByKind localiza a instância do composto na tela.
func componentByKind(t *testing.T, page *screen.Screen, kind catalog.Kind) *screen.Component {
	t.Helper()
	for index := range page.Components {
		if page.Components[index].Kind == kind {
			return &page.Components[index]
		}
	}
	t.Fatalf("componente %s ausente", kind)
	return nil
}

// findBolao devolve o primeiro nó com layout bolao na árvore do menu.
func findBolao(node *catalog.SubmenuProperties) *catalog.SubmenuProperties {
	if node == nil {
		return nil
	}
	if node.Layout == catalog.MenuLayoutBolao {
		return node
	}
	for _, item := range node.Items {
		if found := findBolao(item.Submenu); found != nil {
			return found
		}
	}
	return nil
}

// slideContent lê uma parte XML do pacote sem extrair o arquivo para o disco.
func slideContent(t *testing.T, archive *zip.Reader, name string) string {
	t.Helper()
	file, err := archive.Open(name)
	if err != nil {
		t.Fatalf("parte %s ausente: %v", name, err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("leitura de %s falhou: %v", name, err)
	}
	return string(content)
}

// TestCompositesExportEditableShapes confere ZIP, textos exatos, formas editáveis e ausência de imagem.
func TestCompositesExportEditableShapes(t *testing.T) {
	document, input, processed, menu := compoundDocument(t)
	// target grava o pacote no diretório temporário exclusivo do teste.
	target := filepath.Join(t.TempDir(), "compostos.pptx")
	if err := Export(context.Background(), target, document, Options{Exploration: true}, nil); err != nil {
		t.Fatalf("exportação falhou: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("pacote não foi gravado: %v", err)
	}
	if err := ValidatePackage(data); err != nil {
		t.Fatalf("pacote PPTX inválido: %v", err)
	}
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("ZIP não abre: %v", err)
	}
	// o primeiro slide é o menu da jornada; o segundo contém a tela composta.
	content := slideContent(t, archive, "ppt/slides/slide2.xml")
	// texts reúne rótulos, títulos e valores que precisam sobreviver como texto editável.
	texts := []string{
		input.Properties.BarcodeInput.Label,
		"Dados do boleto",
		"CNPJ:",
		"Nome fantasia:",
		"12.345.678/0001-90",
		"Ação Comércio Ltda",
		"BOLETO",
		"MEGA-SENA",
	}
	for _, text := range texts {
		if !strings.Contains(content, text) {
			t.Fatalf("texto %q ausente do slide", text)
		}
	}
	// caixas tipográficas confirmam formas editáveis, nunca bitmap.
	if !strings.Contains(content, `<p:cNvSpPr txBox="1"/>`) {
		t.Fatal("formas não são caixas de texto editáveis")
	}
	// names exige a nomenclatura por item e por forma interna do composto.
	names := []string{
		input.ID + "/campo-1",
		input.ID + "/valor-8",
		processed.ID + "/bloco-1-título",
		processed.ID + "/bloco-1-linha-1-rótulo",
		processed.ID + "/bloco-1-linha-2-valor",
		menu.ID + "/" + menu.Properties.Menu.Root.Items[0].ID,
	}
	for _, name := range names {
		if !strings.Contains(content, `name="`+name+`"`) {
			t.Fatalf("forma %q não encontrada no slide", name)
		}
	}
	// Nenhum payload base64 ou screenshot substitui os textos e as formas editáveis.
	lower := strings.ToLower(content)
	for _, forbidden := range []string{"base64"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("saída contém %q, incompatível com forma editável", forbidden)
		}
	}
	// o único vínculo do slide é a volta ao menu da jornada; os compostos são estáticos.
	if count := strings.Count(content, "ppaction://hlinksldjump"); count != 1 {
		t.Fatalf("compostos não podem criar vínculos próprios: %d", count)
	}
	if !strings.Contains(content, `name="Menu"><a:hlinkClick`) {
		t.Fatal("vínculo de volta ao menu da jornada ausente ou deslocado")
	}
	// Cada mídia deve ser uma decoração original conhecida, nunca uma captura da tela.
	// allowed conserva os bytes independentes usados como referência para a exportação.
	allowed := [][]byte{}
	// path identifica as duas decorações dos serviços raiz.
	for _, path := range []string{"/azul/botao_menu_released.png", "/azul/botao_menu_seta_released.png"} {
		// data representa o recurso original e não uma renderização do teste.
		data, err := skin.MustBuiltin().Asset(path)
		if err != nil {
			t.Fatal(err)
		}
		allowed = append(allowed, data)
	}
	// file percorre as mídias incorporadas para conferir sua proveniência por bytes.
	for _, file := range archive.File {
		if strings.HasPrefix(file.Name, "ppt/media/") {
			// reader e data permitem comparar os bytes exatos da parte OPC.
			reader, err := file.Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
			// matched só aceita as decorações originais declaradas acima.
			matched := false
			// original compara sem depender de nome ou digest gerado pelo exportador.
			for _, original := range allowed {
				matched = matched || bytes.Equal(data, original)
			}
			if !matched {
				t.Fatalf("mídia sem origem no skin: %s", file.Name)
			}
		}
	}

}

// TestCompositesExportUsesInitialMenu confere que o PPT estático usa o submenu inicial escolhido.
func TestCompositesExportUsesInitialMenu(t *testing.T) {
	document, _, _, _ := compoundDocument(t)
	menuComponent := componentByKind(t, &document.Screens[0], catalog.ApplicationMenu)
	digital := findBolao(&menuComponent.Properties.Menu.Root)
	if digital == nil {
		t.Fatal("BOLÕES DIGITAIS não encontrado no seed")
	}
	// a escolha inicial aponta para o submenu de três legendas, não para a raiz.
	menuComponent.Properties.Menu.InitialMenuID = digital.ID
	// data contém somente a exploração, sem roteiros temporizados.
	data, err := Build(context.Background(), document, Options{Exploration: true}, nil)
	if err != nil {
		t.Fatalf("exportação do menu inicial falhou: %v", err)
	}
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("ZIP não abre: %v", err)
	}
	content := slideContent(t, archive, "ppt/slides/slide2.xml")
	// legends fixa o conteúdo estático das três legendas do layout bolao.
	legends := []string{
		"Visualizar cotas digitais para venda na lotérica.",
		"Consulta bolões de concursos abertos para comercialização.",
		"Consulta bolões de concursos encerrados e não prescritos.",
	}
	if !strings.Contains(content, "BOLÕES DIGITAIS") {
		t.Fatal("faixa do submenu inicial ausente")
	}
	for _, legend := range legends {
		if !strings.Contains(content, legend) {
			t.Fatalf("legenda %q ausente do quadro estático", legend)
		}
	}
	// a raiz não pode vazar para o quadro quando outro submenu é o inicial.
	if strings.Contains(content, "MEGA-SENA") {
		t.Fatal("quadro estático exibiu itens da raiz em vez do menu inicial")
	}
	// os três botões continuam formas editáveis e sem vínculo ou animação próprios.
	if !strings.Contains(content, `<p:cNvSpPr txBox="1"/>`) {
		t.Fatal("menu estático perdeu caixas de texto editáveis")
	}
	if count := strings.Count(content, "ppaction://hlinksldjump"); count != 1 {
		t.Fatalf("menu estático não pode criar vínculos próprios: %d", count)
	}
	if strings.Count(content, `name="`+menuComponent.ID+`/faixa"`) != 1 {
		t.Fatal("faixa do submenu inicial não foi nomeada para edição")
	}
	for index, item := range digital.Items {
		name := menuComponent.ID + "/" + item.ID
		if !strings.Contains(content, `name="`+name+`"`) {
			t.Fatalf("botão %d do menu inicial ausente: %s", index+1, name)
		}
	}
}
