# Feature Specification: Paleta — ComponentTab

**Feature Branch**: `023-paleta-componentes` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/ComponentTab/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Paleta oferece sete famílias base, cinco presets e três compostos em ELEMENTOS DO LEGADO. Clique/drop usam as definições tipadas; imagens de skin são aplicadas na cena, não como novos itens inseríveis.

**Pendências e limites:** Galeria de todas as skins/16 famílias e metadados de disponibilidade completos continuam planejados. Panel/Label/ScrollBar declarados no acervo não se tornam automaticamente novos Kinds de paleta.

**Fontes verificadas:** [palette.go](../../apps/gerador-jornada/internal/adapter/desktop/palette.go), [tools.go](../../apps/gerador-jornada/internal/adapter/desktop/tools.go), [model.go](../../apps/gerador-jornada/internal/catalog/model.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero encontrar todas as famílias e conjuntos do legado para inserir na tela.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve exibir as sete famílias e os conjuntos auxiliares em uma aba própria.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve inserir por clique ou por drop em uma coordenada da imagem.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve separar famílias primitivas de conjuntos prontos e relacionar nomes amigáveis aos nomes EFL.

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

Como autor, quero trocar de seleção ou modo e voltar ao painel sem perda ou alteração implícita dos dados confirmados.

**Why this priority**: O mesmo trabalho precisa continuar coerente ao mudar de contexto.
**Independent Test**: Usar uma fixture com texto acentuado e geometria conhecida, executar o percurso abaixo e comparar os dados antes/depois.

**Acceptance Scenarios**:

1. **Given** uma configuração confirmada, **When** o autor executa o percurso de preservação, **Then** o sistema deve manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
2. **Given** uma operação cancelada ou recusada, **When** o autor retorna ao trabalho, **Then** nenhum dado confirmado foi perdido nem parcialmente substituído.

### Edge Cases

- Paleta maior que a janela, rolagem antes do arraste, drop em painel lateral e nenhum projeto aberto.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Exibir as sete famílias e os conjuntos auxiliares em uma aba própria.
- **FR-002**: Inserir por clique ou por drop em uma coordenada da imagem.
- **FR-003**: Separar famílias primitivas de conjuntos prontos e relacionar nomes amigáveis aos nomes EFL.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Definições do catálogo e conjuntos disponíveis; não persistir itens de paleta dentro do projeto.
- **Estado de interação**: Disponível; arrastando; drop aceito; drop recusado; sem tela.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `023-paleta-componentes/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 003-catalogo.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Oferta planejada — códigos de barras (2026-10-23)

- **FR-008**: Expor “Código de barras” e “Código de barras processado” como dois compostos do Catálogo, não como oito campos soltos ou um Report renomeado. Consumir [003/FR-013](../003-catalogo/spec.md#codigos-barras-2026-10-23) e inserir pelos comandos de Tela, sem oferecer tipo apenas por existir enum; a cadeia de codec/cena/prévia/PPT precisa estar pronta. Inserção usa os defaults identificados em 031/032, sem números ou dados pessoais preenchidos.

**Aceitação**: cada opção cria uma instância com ID próprio e configuração independente; processado oferece edição por bloco no inspector. Oferta ainda não implementada. Integração única em [001/T136](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); modelos/defaults pertencem a 003 e criação a 004.

## Itens legados — lote autorizado 031/032/033 (2026-10-23)

- **FR-009**: Agrupar **Código de barras**, **Código de barras processado** e **Menu aplicação** em **Itens legados**, sem misturar os compostos com partes internas ou gerar uma tela por submenu. Cada inserção aplica defaults do Catálogo e cria uma instância independente pelos comandos de Tela; label amigável não é Kind novo nem implica execução EFL.

**Aceitação**: as três opções ficam no mesmo grupo quando suas cadeias técnicas estiverem prontas; menu nasce com árvore de serviços/jogos do usuário, habilitação e inicial próprios. Preparação de formulário/oferta em [001/T136](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23) e [001/T141](../001-gerador-jornada/tasks.md#tarefas-menu-aplicacao-2026-10-23); liberação técnica barcode em T138 e agrupamento conjunto em T143, sem ciclo ou liberação só por enum. Usuário autorizou implementar, não homologou oferta/visual por esta revisão. Nenhum backlog duplicado nesta feature.
