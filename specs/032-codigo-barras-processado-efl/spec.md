# Feature Specification: Código de barras processado — blocos parametrizáveis

**Feature Branch**: `032-codigo-barras-processado-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-23
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Segundo componente solicitado pelo usuário: apresentar dados de código de barras processado em blocos parametrizáveis, com valores editáveis e escolha por combo de CPF/CNPJ e Nome/Nome fantasia/Razão social. Manter o padrão GitHub Spec Kit e decomposição por domínio.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md#codigos-barras-2026-10-23); instâncias em [004-tela](../004-tela/spec.md#codigos-barras-2026-10-23); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md#codigos-barras-2026-10-23).
**Referência**: “segundo anexo” citado na conversa, mas sem imagem efetivamente recebida. O [DESIGN.md](DESIGN.md) não atribui à imagem posições, cores, títulos ou quantidade de blocos não observados.

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Ficha composta 530 × 120 com blocos/linhas identificados, cinco rótulos permitidos e valores literais; inspector administra rascunho, tipos e ordem. Salvar/reabrir, cópia e PPT usam as propriedades confirmadas.

**Pendências e limites:** A imagem de referência continua ausente; posições/cores do renderer são proposta documentada, não reprodução homologada do anexo. Diagnóstico integral de overflow e matriz completa de identidade/edição permanecem pendentes.

**Fontes verificadas:** [barcode_processed.go](../../apps/gerador-jornada/internal/catalog/barcode_processed.go), [properties_barcode_processed.go](../../apps/gerador-jornada/internal/adapter/desktop/properties_barcode_processed.go), [compound_test.go](../../apps/gerador-jornada/internal/adapter/projectzip/compound_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configurar os dados de cada bloco (Priority: P1)

Como autor da jornada, quero definir blocos de informação com rótulos e valores próprios para exemplificar um resultado processado sem depender de um sistema bancário.

**Why this priority**: O valor do componente é permitir compor dados diferentes em cada bloco, não fixar uma única ficha para todas as jornadas.
**Independent Test**: Inserir o componente, configurar dois blocos com títulos/valores distintos e comparar editor/prévia, sem componente de entrada, código lido ou serviço externo.

**Acceptance Scenarios**:

1. **Given** uma nova instância, **When** consultar as propriedades, **Then** existe um bloco com título vazio e duas linhas iniciais, CNPJ e Nome, ambas com valor vazio; nenhum documento ou nome fictício é inserido automaticamente.
2. **Given** dois blocos, **When** alterar o título ou valor de uma linha do segundo, **Then** somente essa linha/bloco muda; o primeiro permanece integralmente igual.
3. **Given** blocos e linhas existentes, **When** adicionar, reordenar ou remover um item e confirmar, **Then** a ordem visual acompanha a configuração e a operação pode ser desfeita/refeita sem perder conteúdo ou identidades dos itens restantes.
4. **Given** um bloco vazio de valores, **When** visualizar a tela, **Then** os rótulos permanecem com valores em branco, sem inventar dados ou exibir um falso estado de processamento concluído.

---

### User Story 2 - Escolher o significado de cada rótulo pelo combo (Priority: P1)

Como autor, quero escolher se uma linha exibe CPF ou CNPJ e se uma linha de nome exibe Nome, Nome fantasia ou Razão social, mantendo seu valor parametrizável.

**Why this priority**: Atende diretamente à variação de pessoas e nomes solicitada, sem alterar código para cada tela.
**Independent Test**: Em dois blocos, alternar as cinco opções suportadas nos seletores e alterar valores com acentos e pontuação, conferindo independência e preservação ao cancelar ou confirmar.

**Acceptance Scenarios**:

1. **Given** uma linha com CNPJ e valor configurado, **When** selecionar CPF no combo de propriedades, **Then** seu rótulo passa a “CPF:” e o valor permanece disponível para edição; nenhum valor é apagado, consultado ou convertido silenciosamente.
2. **Given** uma linha Nome, **When** escolher Nome fantasia ou Razão social, **Then** o texto correspondente é exibido e o valor permanece o informado pelo autor, inclusive acentos, espaços e vírgulas.
3. **Given** escolhas diferentes em dois blocos, **When** confirmar a configuração, **Then** cada linha conserva sua própria escolha e valor, sem um seletor global substituir os demais.
4. **Given** um rascunho com seleções alteradas, **When** cancelar a edição, **Then** os rótulos, valores, ordem e geometria confirmados não mudam; o popup do combo nunca entra no documento nem na tela exportada.

---

### User Story 3 - Preservar e apresentar a composição (Priority: P2)

Como autor, quero reutilizar e transportar a ficha parametrizada sem transformá-la em resultado real de uma consulta financeira.

**Why this priority**: O componente precisa continuar coerente entre persistência, prévia e apresentação.
**Independent Test**: Duplicar instância e tela, salvar/reabrir e exportar duas configurações distintas; conferir IDs, textos, ordem e editabilidade das formas no PPT.

**Acceptance Scenarios**:

1. **Given** configuração confirmada com mais de um bloco, **When** salvar/reabrir ou desfazer/refazer uma alteração, **Then** títulos, tipos de rótulo, valores, ordem e estilo permanecem exatamente configurados.
2. **Given** uma instância duplicada, **When** editar uma linha da cópia, **Then** o original não muda e a cópia possui novos IDs de componente, blocos e linhas, sem perder valores ou seletores.
3. **Given** a tela na prévia e no PPT, **When** observar o componente, **Then** aparecem os mesmos dados configurados, como conteúdo de leitura; no PPT os rótulos e valores continuam textos editáveis, não screenshot.
4. **Given** também uma entrada de código de barras válida na tela, **When** ela é concluída ou reiniciada, **Then** os nomes/documentos configurados nesta ficha não são preenchidos, apagados ou alterados automaticamente e nenhuma navegação é disparada.

### Edge Cases

- Ao menos um bloco e uma linha por bloco; recusar configuração estruturalmente vazia ou IDs repetidos sem publicação parcial. Títulos e valores vazios são permitidos.
- Nova linha começa como Nome, valor vazio e novo ID. As cinco opções são CPF, CNPJ, Nome, Nome fantasia e Razão social; enum desconhecido é erro localizado, não rótulo escolhido por aproximação.
- Valores são textos ilustrativos, não documentos fiscais validados. Pontuação/máscara digitada permanece literal; troca CPF ↔ CNPJ não calcula DV, não altera comprimento nem completa zeros.
- O seletor altera uma linha por vez; permitir linhas repetidas e combinações diferentes entre blocos. Não inferir que uma pessoa é física/jurídica pelo nome ou valor.
- Título opcional vazio não gera título fictício nem reserva linha de título. Nomes longos exigem quebra/área suficiente ou diagnóstico de overflow, nunca truncamento persistido.
- Remover ou reordenar não troca identidade dos itens restantes. Duplicação não compartilha listas mutáveis ou UUIDs internos.
- Nenhuma linha recebe nome/CPF/CNPJ por decodificação do código de barras: esses dados não são obtidos pelas regras matemáticas da entrada.
- A escolha é no painel de propriedades do autor, não um combo interativo obrigatório desenhado na jornada. Essa premissa de UX pode ser revista com a imagem; a versão atual da spec a torna explícita.
- DPI 100%, 150% e 200%, navegação por teclado nos seletores, texto Unicode e exportação com todos os valores completos.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Oferecer componente composto de leitura com lista ordenada de blocos, cada um com identidade própria, título opcional e linhas ordenadas de rótulo/valor, sem codificar quantidade fixa de blocos a partir de imagem ausente.
- **FR-002**: Permitir editar individualmente todos os valores e títulos, adicionar/remover/reordenar blocos e linhas, exigindo ao menos um bloco com uma linha e mantendo identidade dos itens preservados.
- **FR-003**: Oferecer um combo por linha no painel de propriedades para escolher CPF, CNPJ, Nome, Nome fantasia ou Razão social; o rótulo visível corresponde exatamente à opção selecionada, seguido de dois-pontos.
- **FR-004**: Preservar literalmente o valor da linha ao trocar seu rótulo, sem conversão, consulta ou descarte automático; cada linha/bloco é independente dos demais.
- **FR-005**: Iniciar com um bloco de título vazio e duas linhas CNPJ/Nome com valores vazios; nova linha usa Nome e valor vazio; não oferecer dados pessoais reais ou exemplos automáticos como conteúdo inicial.
- **FR-006**: Validar estrutura, IDs, tipos de rótulo e geometria e recusar alterações inválidas atomicamente, preservando o último estado confirmado; não exigir validade fiscal de um valor meramente ilustrativo.
- **FR-007**: Persistir configuração confirmada, valores, tipos, IDs e ordem; duplicar profundamente com novas identidades internas; integrar histórico e ID de arquivo comum sem converter blocos em componentes de Tela independentes.
- **FR-008**: Apresentar a mesma composição no editor, prévia e PPT com rótulos/valores editáveis e estilo conforme [DESIGN.md](DESIGN.md), sem exportar controles do inspector ou seleção do editor.
- **FR-009**: Ser independente da entrada validada: não executar consulta, pagamento, vínculo implícito, navegação ou inferência de nome/documento; o nome “processado” descreve a representação ilustrativa, não um processamento bancário realizado.
- **FR-010**: Manter interface em português, texto Unicode, foco perceptível e diagnóstico de overflow; cancelamento de rascunho e interação na prévia não alteram o documento.

### Key Entities

- **Configuração processada**: coleção ordenada de blocos, pertencente às propriedades tipadas do Catálogo, sem entidade bancária ou cadastro externo.
- **Bloco parametrizável**: ID, título opcional e linhas ordenadas; não é um domínio nem subcomponente persistido de Tela.
- **Linha parametrizável**: ID, tipo de rótulo selecionado e valor textual ilustrativo; texto do rótulo é derivado do tipo, não uma segunda string concorrente.
- **Instância de Tela**: UUID/FileID, bounds/estilo e configuração processada; popup/seleção de linha do inspector são transitórios.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: As cinco opções de rótulo produzem os cinco textos esperados em 100% dos cenários de dois blocos, sem alterar valores de outra linha ou bloco.
- **SC-002**: Todos os títulos/valores e sua pontuação/Unicode permanecem idênticos após salvar/reabrir, cancelar edição, desfazer/refazer e duplicar; a cópia possui IDs novos em todos os níveis internos.
- **SC-003**: Editor, prévia e PPT apresentam a mesma ordem/conteúdo, com diferença geométrica máxima de uma unidade lógica; todos os rótulos/valores no PPT são textos editáveis.
- **SC-004**: Todos os casos estruturais inválidos da matriz de referência são recusados com localização de componente/bloco/linha, sem perda de estado confirmado ou exposição de valores em logs.

## Assumptions

- A imagem citada como segundo anexo não foi recebida; títulos, papéis de beneficiário/pagador, cores, medidas exatas e quantidade de blocos da referência continuam sem evidência. Não se anuncia paridade visual.
- Por premissa explícita, os combos são controles de autoria no inspector. Na tela da jornada o componente é de leitura; combo runtime, condicionais ou integração automática com a entrada exigem novo requisito.
- Os valores são configurados pelo autor para a demonstração. Não há serviço real ou dados extraídos de boleto para preencher CPF/CNPJ, nome, nome fantasia ou razão social.
- Usar tipos próprios de propriedades do Catálogo, não texto CSV, HTML, script, JSON livre ou um Report renomeado; reutilizar pintura de rótulos/valores sem duplicar domínio.
- Requisitos são locais: usar `032-codigo-barras-processado-efl/FR-xxx`. A entrega é spec + DESIGN + checklist; os modelos/planos/tarefas continuam nos domínios canônicos. Java/EFL não é reativado.

## Divisão por domínio e rastreabilidade

| Responsabilidade | Artefato e trabalho responsável |
| --- | --- |
| Tipos de rótulo, blocos/linhas e validação | [003/FR-013/FR-015 e T018/T020](../003-catalogo/tasks.md#tarefas-codigos-barras-2026-10-23); FR-001–FR-006 |
| Instâncias, configuração atômica e duplicação | [004/FR-019 e T015–T016](../004-tela/tasks.md#tarefas-codigos-barras-2026-10-23); FR-002/FR-006–FR-007 |
| Codec, versões e preservação de valores | [002/FR-011 e T013–T014](../002-projeto/tasks.md#tarefas-codigos-barras-2026-10-23); FR-007 |
| Não criar rotas implícitas | [005/FR-008 e T011](../005-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); FR-009 |
| Conteúdo estático e limites de roteiro | [006/FR-011 e T011](../006-demonstracao/tasks.md#tarefas-codigos-barras-2026-10-23); FR-008–FR-009 |
| Inspector, cena, prévia e PPT | [001/FR-060 e T136–T139](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); FR-003/FR-008/FR-010 |

Este detalhamento não acumula tarefas de implementação. Fonte de verdade dos dados e enums no [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md#contrato-codigos-barras-2026-10-23); revisão documental no [checklist](checklists/requirements.md).

## ID das instâncias

Aplicar [004/FR-011–FR-012](../004-tela/spec.md#identidade-de-arquivo--2026-10-18) somente à instância de Tela. Blocos e linhas têm UUIDs internos para identidade/ordem/cópia, não FileIDs ou arquivos próprios. Trocar o rótulo não renomeia identidades.