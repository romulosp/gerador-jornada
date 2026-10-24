package application

import (
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// addCompound insere um composto novo na tela e devolve sua identidade.
func addCompound(t *testing.T, document *Document, pageIndex int, kind catalog.Kind, name string, x, y int) string {
	t.Helper()
	// component recebe os defaults do catálogo e a posição solicitada.
	component, err := screen.NewComponent(kind, name, x, y)
	if err != nil {
		t.Fatal(err)
	}
	document.Screens[pageIndex].Components = append(document.Screens[pageIndex].Components, component)
	return component.ID
}

// typeFields digita os valores nos oito campos na ordem informada.
func typeFields(t *testing.T, preview *Preview, id string, fields []string) {
	t.Helper()
	for _, field := range fields {
		for _, value := range field {
			if err := preview.TypeBarcodeDigit(id, value); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// TestBarcodeFieldAllowedBloqueiaBlocos confere que cliques não atravessam bloco inválido.
func TestBarcodeFieldAllowedBloqueiaBlocos(t *testing.T) {
	// document contém uma entrada vazia.
	document := NewDocument("runtime")
	id := addCompound(t, &document, 0, catalog.BarcodeInput, "codigo", 40, 40)
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	if !preview.BarcodeFieldAllowed(id, 0) || preview.BarcodeFieldAllowed(id, 2) {
		t.Fatal("bloqueio inicial incorreto")
	}
	typeFields(t, preview, id, []string{"00190"})
	if !preview.BarcodeFieldAllowed(id, 1) || preview.BarcodeFieldAllowed(id, 2) {
		t.Fatal("campo 2 deveria estar liberado e o bloco seguinte não")
	}
	typeFields(t, preview, id, []string{"50095"})
	if !preview.BarcodeFieldAllowed(id, 2) || preview.BarcodeFieldAllowed(id, 4) {
		t.Fatal("bloco 2 deveria estar liberado após o bloco 1 válido")
	}
}

// TestBarcodeRuntimeF01 valida digitação por blocos, conclusão e dados derivados.
func TestBarcodeRuntimeF01(t *testing.T) {
	// document contém uma única tela com a entrada de referência.
	document := NewDocument("runtime")
	id := addCompound(t, &document, 0, catalog.BarcodeInput, "codigo", 40, 40)
	// preview isola a sessão e inicializa a entrada vazia.
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	typeFields(t, preview, id, []string{"00190", "50095", "40144", "816069", "06809", "350314", "3", "3737000000010"})
	// runtime contém o rascunho antes da conclusão.
	runtime, _ := preview.Barcode(id)
	if runtime.Locked || runtime.Values[7] != "3737000000010" {
		t.Fatalf("rascunho inesperado: %+v", runtime)
	}
	if err := preview.BarcodeFinish(id); err != nil {
		t.Fatal(err)
	}
	runtime, _ = preview.Barcode(id)
	if !runtime.Locked || runtime.Values[7] != "37370000000100" {
		t.Fatalf("conclusão com zeros falhou: %+v", runtime)
	}
	if runtime.Result.Bank != "001" || runtime.Result.ValueCents != 100 || runtime.Result.Factor != 3737 {
		t.Fatalf("derivados inesperados: %+v", runtime.Result)
	}
	// reiniciar limpa somente esta instância.
	if err := preview.BarcodeRestart(id); err != nil {
		t.Fatal(err)
	}
	runtime, _ = preview.Barcode(id)
	if runtime.Locked || runtime.Result.Bank != "" || runtime.Values != [8]string{} {
		t.Fatalf("reinício incompleto: %+v", runtime)
	}
}

// TestBarcodeBlockErrorAndModal cobre erro local por par e captura integral inválida.
func TestBarcodeBlockErrorAndModal(t *testing.T) {
	// document contém a entrada usada nas duas falhas.
	document := NewDocument("runtime")
	id := addCompound(t, &document, 0, catalog.BarcodeInput, "codigo", 40, 40)
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	typeFields(t, preview, id, []string{"00190", "50095", "40144", "816068"})
	runtime, _ := preview.Barcode(id)
	if !strings.Contains(runtime.Message, "Bloco 2 inválido") || runtime.Active != 2 {
		t.Fatalf("erro de bloco inesperado: %+v", runtime)
	}
	if runtime.Values[2] != "40144" || runtime.Values[3] != "816068" {
		t.Fatal("erro de bloco apagou o rascunho")
	}
	// captura inválida limpa a instância e abre o modal.
	if err := preview.BarcodePaste(id, "123"); err != nil {
		t.Fatal(err)
	}
	runtime, _ = preview.Barcode(id)
	if !runtime.Modal || runtime.Message != "Código de barras inválido" || runtime.Values != [8]string{} {
		t.Fatalf("modal inválido: %+v", runtime)
	}
	if err := preview.BarcodeDismissModal(id); err != nil || runtime.Modal {
		t.Fatalf("modal não fechou: %v", err)
	}
}

// TestBarcodeTwoInstancesIsolated confere que digitar em uma entrada não altera a outra.
func TestBarcodeTwoInstancesIsolated(t *testing.T) {
	// document contém duas entradas na mesma tela.
	document := NewDocument("runtime")
	first := addCompound(t, &document, 0, catalog.BarcodeInput, "primeira", 20, 20)
	second := addCompound(t, &document, 0, catalog.BarcodeInput, "segunda", 20, 120)
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	typeFields(t, preview, first, []string{"12345"})
	if runtime, _ := preview.Barcode(second); runtime.Values[0] != "" {
		t.Fatalf("segunda instância recebeu digitação: %+v", runtime)
	}
	// revisitar a tela reinicia ambas as entradas.
	if err := preview.Enter(document.Screens[0].ID); err != nil {
		t.Fatal(err)
	}
	if runtime, _ := preview.Barcode(first); runtime.Values != [8]string{} {
		t.Fatalf("revisita preservou valores: %+v", runtime)
	}
}

// TestMenuRuntimeNavigation cobre abertura, Voltar, Cancelar, folha e bloqueio.
func TestMenuRuntimeNavigation(t *testing.T) {
	// document contém o menu e uma segunda tela para destino explícito.
	document := NewDocument("runtime")
	document.Screens = append(document.Screens, screen.New("02 · Destino"))
	id := addCompound(t, &document, 0, catalog.ApplicationMenu, "menu", 0, 0)
	// menu permite configurar uma folha com destino e um item desabilitado.
	menu := document.Screens[0].ComponentByID(id).Properties.Menu
	// bolao é o item de submenu usado no percurso.
	var bolao *catalog.MenuItem
	for index := range menu.Root.Items {
		if menu.Root.Items[index].Label == "BOLÃO" {
			bolao = &menu.Root.Items[index]
		}
	}
	if bolao == nil {
		t.Fatal("BOLÃO ausente no seed")
	}
	// alvo é uma folha raiz convertida em destino explícito.
	target := &menu.Root.Items[10]
	target.TargetScreenID = document.Screens[1].ID
	// desabilitado permanece habilitado no modelo; a cópia da prévia usa o valor confirmado.
	disabled := &menu.Root.Items[1]
	disabled.Enabled = false
	preview, err := NewPreview(document)
	if err != nil {
		t.Fatal(err)
	}
	// abrir BOLÃO empilha o submenu e Voltar retorna à raiz.
	if outcome, err := preview.MenuActivate(id, bolao.ID); err != nil || outcome != MenuOpened {
		t.Fatalf("abertura inesperada: %v %v", outcome, err)
	}
	node := preview.MenuNode(id)
	if node == nil || node.ID != bolao.Submenu.ID {
		t.Fatalf("nó visível inesperado: %+v", node)
	}
	if err := preview.MenuBack(id); err != nil {
		t.Fatal(err)
	}
	if node = preview.MenuNode(id); node == nil || node.ID != menu.Root.ID {
		t.Fatal("Voltar não retornou à raiz")
	}
	// Cancelar volta à raiz mesmo quando o inicial é submenu.
	if outcome, err := preview.MenuActivate(id, bolao.ID); err != nil || outcome != MenuOpened {
		t.Fatalf("segunda abertura falhou: %v %v", outcome, err)
	}
	if err := preview.MenuCancel(id); err != nil {
		t.Fatal(err)
	}
	// item desabilitado não produz efeito.
	if outcome, err := preview.MenuActivate(id, disabled.ID); err != nil || outcome != MenuBlocked {
		t.Fatalf("item desabilitado: %v %v", outcome, err)
	}
	// folha com destino explícito abre a tela autorada.
	if outcome, err := preview.MenuActivate(id, target.ID); err != nil || outcome != MenuNavigated {
		t.Fatalf("navegação explícita: %v %v", outcome, err)
	}
	if preview.ScreenID != document.Screens[1].ID {
		t.Fatal("tela de destino não foi aberta")
	}
}

// TestMenuCurrentScreenDoesNotNavigate garante o escopo da Prévia Tela Atual.
func TestMenuCurrentScreenDoesNotNavigate(t *testing.T) {
	// document contém o menu e uma segunda tela, mas a prévia isola a primeira.
	document := NewDocument("runtime")
	document.Screens = append(document.Screens, screen.New("02 · Destino"))
	id := addCompound(t, &document, 0, catalog.ApplicationMenu, "menu", 0, 0)
	menu := document.Screens[0].ComponentByID(id).Properties.Menu
	menu.Root.Items[10].TargetScreenID = document.Screens[1].ID
	preview, err := NewScreenPreview(document, document.Screens[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if outcome, err := preview.MenuActivate(id, menu.Root.Items[10].ID); err != nil || outcome != MenuSimulated {
		t.Fatalf("Tela Atual deveria simular: %v %v", outcome, err)
	}
	if preview.ScreenID != document.Screens[0].ID {
		t.Fatal("Tela Atual navegou entre telas")
	}
}
