package catalog

// seedApplicationMenu monta a árvore inicial de serviços e jogos.
// supplyID entrega UUIDs novos; a árvore é somente dados inertes, sem execução de movimentos.
func seedApplicationMenu(supplyID func() string) *ApplicationMenuProperties {
	// root agrega os doze serviços e os dez jogos das duas colunas.
	root := SubmenuProperties{ID: supplyID(), Name: "Menu da aplicação", Title: "", Layout: MenuLayoutRoot}
	// node cria um submenu vazio com layout explícito.
	node := func(name string, layout string) *SubmenuProperties {
		return &SubmenuProperties{ID: supplyID(), Name: name, Title: name, Layout: layout, Items: []MenuItem{}}
	}
	// branch cria um item de ramo habilitado.
	branch := func(label string, submenu *SubmenuProperties) MenuItem {
		return MenuItem{ID: supplyID(), Label: label, Category: "SERVIÇOS", Enabled: true, Submenu: submenu}
	}
	// leaf cria uma folha habilitada sem destino, apenas ilustrativa.
	leaf := func(label string) MenuItem {
		return MenuItem{ID: supplyID(), Label: label, Category: "SERVIÇOS", Enabled: true}
	}
	// services preserva a ordem documentada dos serviços da coluna esquerda.
	services := []MenuItem{
		leaf("BOLETO"),
		branch("CONVÊNIOS", node("CONVÊNIOS", MenuLayoutStandard)),
		branch("BENEFÍCIOS", node("BENEFÍCIOS", MenuLayoutStandard)),
		branch("FINANCEIROS", node("FINANCEIROS", MenuLayoutFinancial)),
		branch("ATENDIMENTO", node("ATENDIMENTO", MenuLayoutStandard)),
		branch("RESULTADOS", node("RESULTADOS", MenuLayoutStandard)),
		branch("APOIO", node("APOIO", MenuLayoutStandard)),
	}
	// bolao contém o Bolão ordinário e o Bolões Digitais como último filho.
	bolao := node("BOLÃO", MenuLayoutStandard)
	bolao.Items = append(bolao.Items, leaf("BOLÃO"))
	digital := node("BOLÕES DIGITAIS", MenuLayoutBolao)
	digital.Items = []MenuItem{
		{ID: supplyID(), Label: "BOLÕES DIGITAIS", Legend: "Visualizar cotas digitais para venda na lotérica.", Category: "SERVIÇOS", Enabled: true},
		{ID: supplyID(), Label: "CONCURSOS ABERTOS", Legend: "Consulta bolões de concursos abertos para comercialização.", Category: "SERVIÇOS", Enabled: true},
		{ID: supplyID(), Label: "CONCURSOS ENCERRADOS", Legend: "Consulta bolões de concursos encerrados e não prescritos.", Category: "SERVIÇOS", Enabled: true},
	}
	bolao.Items = append(bolao.Items, branch("BOLÕES DIGITAIS", digital))
	services = append(services, branch("BOLÃO", bolao))
	// reports reúne a árvore estrutural documentada de 27 títulos e 148 folhas.
	services = append(services, branch("RELATÓRIOS", seedReportsMenu(supplyID)))
	// paga contém os três códigos literais documentados.
	paga := node("PAGA PREMIADOS", MenuLayoutStandard)
	paga.Items = []MenuItem{
		{ID: supplyID(), Label: "PROGNÓSTICOS", Category: "SERVIÇOS", Enabled: true, MovementCode: "1500"},
		{ID: supplyID(), Label: "LOTERIA FEDERAL", Category: "SERVIÇOS", Enabled: true, MovementCode: "1510"},
		{ID: supplyID(), Label: "INSTANTÂNEA", Category: "SERVIÇOS", Enabled: true, MovementCode: "1513"},
	}
	services = append(services, branch("PAGA PREMIADOS", paga))
	// meios contém PIX, cartões e entrada em dinheiro com opções literais.
	meios := node("MEIOS DE PAGAMENTO", MenuLayoutStandard)
	pix := node("PIX", MenuLayoutStandard)
	pix.Items = []MenuItem{
		{ID: supplyID(), Label: "RECEBIMENTO", Category: "SERVIÇOS", Enabled: true, MovementOption: "99"},
		{ID: supplyID(), Label: "CONSULTAR RECEBIMENTO", Category: "SERVIÇOS", Enabled: true, MovementOption: "99"},
	}
	meios.Items = []MenuItem{
		branch("PIX", pix),
		{ID: supplyID(), Label: "CARTÃO CRÉDITO/DÉBITO", Category: "SERVIÇOS", Enabled: true, MovementOption: "99"},
		{ID: supplyID(), Label: "ENTRADA EM DINHEIRO", Category: "SERVIÇOS", Enabled: true, MovementOption: "2"},
	}
	services = append(services, branch("MEIOS DE PAGAMENTO", meios))
	services = append(services, leaf("OFERTAR"))
	root.Items = append(root.Items, services...)
	// games preserva dez jogos evidenciados pelas skins legadas.
	games := []string{"MEGA-SENA", "LOTOFÁCIL", "QUINA", "LOTOMANIA", "LOTECA", "LOTOGOL", "DUPLA SENA", "DIA DE SORTE", "SUPER SETE", "TIMEMANIA"}
	// name identifica o jogo atual da coluna direita.
	for _, name := range games {
		root.Items = append(root.Items, MenuItem{ID: supplyID(), Label: name, Category: "JOGOS", Enabled: true})
	}
	// menu fixa a raiz como estado inicial, sem escolher submenu por título.
	return &ApplicationMenuProperties{Root: root, InitialMenuID: root.ID}
}
