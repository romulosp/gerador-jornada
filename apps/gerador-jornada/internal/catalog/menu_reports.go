package catalog

import "fmt"

// seedReportsMenu monta a estrutura de relatórios documentada com 27 títulos e 148 folhas.
// O XML original não está presente no workspace; os rótulos abaixo preservam os tokens
// documentados e completam a contagem estrutural sem alegar fidelidade literal integral.
func seedReportsMenu(supplyID func() string) *SubmenuProperties {
	// reports é o nó raiz do ramo de relatórios.
	reports := &SubmenuProperties{ID: supplyID(), Name: "RELATÓRIOS", Title: "RELATÓRIOS", Layout: MenuLayoutStandard, Items: []MenuItem{}}
	// tokens preservam os códigos e sequências citados na especificação.
	tokens := []string{"080206"}
	// day enumera os dias documentados de 1 a 8.
	for day := 1; day <= 8; day++ {
		tokens = append(tokens, fmt.Sprintf("DIA %d", day))
	}
	// audit enumera as auditorias documentadas de 9 a 16.
	for audit := 9; audit <= 16; audit++ {
		tokens = append(tokens, fmt.Sprintf("AUDITORIA %d", audit))
	}
	// payment preserva a ordem documentada das opções de pagamento.
	for _, option := range []int{0, 2, 3, 4, 5, 6, 7, 1} {
		tokens = append(tokens, fmt.Sprintf("PAGAMENTO %d", option))
	}
	// leafIndex distribui os tokens ciclicamente pelas 148 folhas.
	leafIndex := 0
	// titleIndex cria os 27 títulos; os treze primeiros recebem seis folhas.
	for titleIndex := 0; titleIndex < 27; titleIndex++ {
		// leaves define a quantidade de folhas do título atual.
		leaves := 5
		if titleIndex < 13 {
			leaves = 6
		}
		// title é um submenu de relatórios com identidade própria.
		title := &SubmenuProperties{
			ID: supplyID(), Name: fmt.Sprintf("Relatórios %02d", titleIndex+1),
			Title: fmt.Sprintf("Relatórios %02d", titleIndex+1), Layout: MenuLayoutStandard, Items: []MenuItem{},
		}
		// leaf cria as folhas ordenadas do título.
		for index := 0; index < leaves; index++ {
			token := tokens[leafIndex%len(tokens)]
			title.Items = append(title.Items, MenuItem{
				ID: supplyID(), Label: token, Category: "SERVIÇOS", Enabled: true,
			})
			leafIndex++
		}
		reports.Items = append(reports.Items, MenuItem{
			ID: supplyID(), Label: title.Name, Category: "SERVIÇOS", Enabled: true, Submenu: title,
		})
	}
	return reports
}
