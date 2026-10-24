# Feature Specification: Relatório/lista — EFLReportEvent

**Feature Branch**: `012-relatorio-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/ListComponents/EFLReportEvent/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Relatório 300 × 345 com título e linhas simples. Cena recompõe painel, cabeçalho e rodapé originais; título/linhas usam label_report_titulo/label_report e permanecem editáveis no PPT.

**Pendências e limites:** As oito classes Java de Report, linhas estruturadas, seleção/rolagem e diagnóstico integral de overflow não estão implementados. Hoje a cena limita as linhas ao corpo visível; isso não conclui FR-008–FR-009.

**Fontes verificadas:** [legacy_parts.go](../../apps/gerador-jornada/internal/presentation/legacy_parts.go), [legacy_style_test.go](../../apps/gerador-jornada/internal/presentation/legacy_style_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero organizar um cabeçalho e linhas de informação visual.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve editar título e linhas ordenadas, preservando IDs independentes dos textos.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve apresentar cabeçalho separado das linhas sem apagá-lo ao atualizar o conteúdo.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve representar a lista sem conectar banco de dados nem executar regras de negócio ao selecionar uma linha.

---

### User Story 2 - Personalizar a aparência com previsibilidade (Priority: P2)

Como autor ou responsável pelo visual, quero ajustar a aparência deste componente sem alterar sua identidade nem suas regras.

**Why this priority**: Permite reproduzir o visual legado e evoluí-lo de forma controlada.
**Independent Test**: Usar o estado de referência do DESIGN.md, alterar uma propriedade visual por vez e comparar resultado e comportamento antes/depois.

**Acceptance Scenarios**:

1. **Given** o estado de referência, **When** uma propriedade visual permitida muda, **Then** a aparência correspondente muda e os três comportamentos da US1 permanecem válidos.
2. **Given** uma personalização fora dos limites, **When** o autor tenta confirmá-la, **Then** o autor recebe diagnóstico no contexto e o último estado válido permanece íntegro.

---

### User Story 3 - Preservar conteúdo e contexto (Priority: P2)

Como autor, quero salvar, reabrir e exportar a composição sem perda ou alteração implícita dos dados confirmados.

**Why this priority**: O mesmo trabalho precisa continuar coerente ao mudar de contexto.
**Independent Test**: Usar uma fixture com texto acentuado e geometria conhecida, executar o percurso abaixo e comparar os dados antes/depois.

**Acceptance Scenarios**:

1. **Given** uma configuração confirmada, **When** o autor executa o percurso de preservação, **Then** o sistema deve salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
2. **Given** uma operação cancelada ou recusada, **When** o autor retorna ao trabalho, **Then** nenhum dado confirmado foi perdido nem parcialmente substituído.

### Edge Cases

- Sem linhas, linhas longas, ordem alterada e conteúdo excedente. Definir rolagem na prévia e representação do recorte visível no PPT; não descartar dados do projeto.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Editar título e linhas ordenadas, preservando IDs independentes dos textos.
- **FR-002**: Apresentar cabeçalho separado das linhas sem apagá-lo ao atualizar o conteúdo.
- **FR-003**: Representar a lista sem conectar banco de dados nem executar regras de negócio ao selecionar uma linha.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Título, linhas, ordem, dimensões, texto, cores e decoração de cabeçalho/corpo/rodapé.
- **Estado de interação**: Lista vazia; preenchida; linha destacada; conteúdo maior que a área; seleção de edição.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `012-relatorio-efl/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 003-catalogo.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## ID das instâncias

Toda instância visual desta responsabilidade possui ID de arquivo conforme [Tela](../004-tela/spec.md#identidade-de-arquivo--2026-10-18) e [Catálogo](../003-catalogo/spec.md#atributo-compartilhado-das-instâncias). Campo comum no [painel de propriedades](../028-propriedades/spec.md#campo-id-comum), sem prefixo Pnl obrigatório para componentes e sem substituir UUIDs ou conteúdo visual. Em conjuntos prontos, a regra vale para cada controle efetivamente criado.

## Piloto e expansão Java/EFL autorizados — 2026-10-19

- **FR-008**: Usar EFLReportEvent como piloto completo do fluxo domínio→editor→prévia→salvar/reabrir v3→PPT editável, com topo/título/corpo/linhas/rodapé, seleção, vazio/overflow/scroll e fontes/skins efetivas. Não tratar o corpo simplificado atual como paridade; reconciliar medidas React/Java antes de afirmar fidelidade.
- **FR-009**: Cobrir como variantes de Report todas as classes Report, Report2, Report3, ReportEvent, ReportEvent2, ReportNotificacao, ReportSelecaoMultipla e EFLInternalReport, com tipos de linha/seleção/altura e partes comprovados por classe. Piloto/reuso não homologa as demais variantes. Sem banco/serviços/regras externas; faltas de recurso/delegate/referência mantêm entradas bloqueadas.

**Aceitação**: fixture Event com título/linhas acentuadas, seleção inicial, vazio e excedente conserva IDs/bounds/overrides no round-trip, preserva cabeçalho e recorte, muda seleção/scroll só em runtime e exporta formas/textos editáveis com face real. Comparar referência independente por estado e DPI; cada outra classe exige fixture/evidência próprias. Modelo em 003/T014, partes Panel em 003/T015, persistência 002/T009–T012 e consumidores 001/T111–T116. [plan.md](plan.md)/[tasks.md](tasks.md) locais verificam o piloto sem duplicar esses donos; manuais/aceite 001/T120–T122. Java permanece pausado; nenhuma homologação nesta documentação.
