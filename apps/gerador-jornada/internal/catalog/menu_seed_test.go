package catalog

import (
	"testing"

	"gerador-jornada.local/editor/internal/project"
)

// TestMenuSeedStructure verifica a árvore real, sem completar contagens com títulos sintéticos.
func TestMenuSeedStructure(t *testing.T) {
	// menu recebe UUIDs novos para todos os nós importados.
	menu := NewApplicationMenuProperties(project.NewID)
	if err := menu.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(menu.Root.Items) != 22 || menu.Root.Items[9].Label != "MILIONÁRIA" || menu.Root.Items[9].Variant != "botao_milionaria" {
		t.Fatal("raiz XML divergente")
	}
	// bolao mantém os dez ramos ativos; digital é seu último filho.
	bolao := menu.Root.Items[17].Submenu
	if bolao == nil || len(bolao.Items) != 10 {
		t.Fatal("BOLÃO divergente")
	}
	// digital conserva nomes e códigos XML, além do layout visual de três posições.
	digital := bolao.Items[9]
	if digital.Label != "BOLÕES DIGITAIS" || digital.MovementCode != "803" || digital.Submenu == nil || digital.Submenu.Layout != MenuLayoutBolao {
		t.Fatal("BOLÕES DIGITAIS divergente")
	}
	// reports e counters verificam os 27 títulos descendentes e as 148 folhas reais.
	reports := menu.Root.Items[18].Submenu
	titles, leaves := 0, 0
	var count func(*SubmenuProperties)
	count = func(node *SubmenuProperties) {
		// item representa um título real ou uma folha da árvore original.
		for _, item := range node.Items {
			if item.Submenu != nil {
				titles++
				count(item.Submenu)
			} else {
				leaves++
			}
		}
	}
	count(reports)
	if titles != 27 || leaves != 148 {
		t.Fatalf("relatórios: %d títulos / %d folhas", titles, leaves)
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
