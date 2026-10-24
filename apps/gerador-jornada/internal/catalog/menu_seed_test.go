package catalog

import (
	"testing"

	"gerador-jornada.local/editor/internal/project"
)

// TestMenuSeedStructure confere a fábrica documentada de serviços, jogos e relatórios.
func TestMenuSeedStructure(t *testing.T) {
	// menu contém a árvore inicial de uma nova instância.
	menu := NewApplicationMenuProperties(project.NewID)
	if err := menu.Validate(); err != nil {
		t.Fatalf("seed inválido: %v", err)
	}
	if menu.InitialMenuID != menu.Root.ID || menu.Root.Layout != MenuLayoutRoot {
		t.Fatalf("raiz/inicial incorretos: %+v", menu.Root)
	}
	// services e games contam as duas colunas da raiz.
	services, games := 0, 0
	for _, item := range menu.Root.Items {
		if item.Category == "JOGOS" {
			games++
		} else {
			services++
		}
	}
	if services != 12 || games != 10 {
		t.Fatalf("raiz deve ter 12 serviços e 10 jogos: %d/%d", services, games)
	}
	// bolao contém BOLÕES DIGITAIS como último filho, em layout bolao.
	bolao := menu.FindSubmenu(menu.Root.Items[7].Submenu.ID)
	if bolao == nil || len(bolao.Items) != 2 {
		t.Fatalf("submenu BOLÃO inesperado: %+v", bolao)
	}
	digital := bolao.Items[len(bolao.Items)-1].Submenu
	if digital == nil || digital.Layout != MenuLayoutBolao || len(digital.Items) != 3 {
		t.Fatalf("BOLÕES DIGITAIS inválido: %+v", digital)
	}
	// legends preserva os textos exatos do contrato.
	legends := []string{
		"Visualizar cotas digitais para venda na lotérica.",
		"Consulta bolões de concursos abertos para comercialização.",
		"Consulta bolões de concursos encerrados e não prescritos.",
	}
	for index, legend := range legends {
		if digital.Items[index].Legend != legend {
			t.Fatalf("legenda %d divergente: %q", index+1, digital.Items[index].Legend)
		}
	}
	// reports contém 27 títulos e 148 folhas no total.
	reports := menu.Root.Items[8].Submenu
	if reports == nil || len(reports.Items) != 27 {
		t.Fatalf("RELATÓRIOS deve ter 27 títulos: %+v", reports)
	}
	leaves := 0
	for _, title := range reports.Items {
		leaves += len(title.Submenu.Items)
	}
	if leaves != 148 {
		t.Fatalf("RELATÓRIOS deve ter 148 folhas: %d", leaves)
	}
}

// TestMenuDuplicateRegeneratesIDs confere que a duplicação troca todas as identidades internas.
func TestMenuDuplicateRegeneratesIDs(t *testing.T) {
	// source é a árvore padrão usada como origem.
	source := NewApplicationMenuProperties(project.NewID).Root
	// duplicated recebe a cópia com IDs novos.
	duplicated := DuplicateApplicationMenu(ApplicationMenuProperties{Root: source, InitialMenuID: source.ID}, project.NewID)
	// original e cópia não compartilham nenhum UUID.
	seen := map[string]bool{}
	for _, id := range CollectMenuIDs(ApplicationMenuProperties{Root: source, InitialMenuID: source.ID}) {
		seen[id] = true
	}
	for _, id := range CollectMenuIDs(duplicated) {
		if seen[id] {
			t.Fatalf("ID interno repetido na cópia: %s", id)
		}
	}
	if duplicated.InitialMenuID != duplicated.Root.ID {
		t.Fatalf("menu inicial não foi remapeado: %s", duplicated.InitialMenuID)
	}
}
