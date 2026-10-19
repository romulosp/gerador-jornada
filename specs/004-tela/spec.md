# Feature Specification: Domínio Tela

**Feature Branch**: `004-tela` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — domínio em implementação; homologação e aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compor uma tela (Priority: P1)

Modelar telas, componentes, geometria e ordem visual com operações coesas e independentes da UI.

**Why this priority**: O desenho precisa manter as mesmas coordenadas e relações ao ser reaberto, apresentado ou convertido para Java.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** uma tela 800 × 600, **When** um componente válido é inserido e movido, **Then** seus bounds são armazenados em unidades lógicas, sem depender do zoom ou DPI da janela.

---

### User Story 2 - Alterar a composição com previsibilidade (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** uma tela com teclado vinculado e botão para si mesma, **When** ela é duplicada, **Then** tela e componentes recebem novos IDs e vínculos internos são remapeados para a cópia.

### Edge Cases

Cobrir bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Criar, renomear, ordenar e dimensionar telas; default 800 × 600, dimensões positivas.
- **FR-002**: Manter componentes com UUID, nome lógico único na tela, kind do catálogo, bounds inteiros, style e properties tipadas.
- **FR-003**: Oferecer inserir, mover, redimensionar, alinhar, duplicar, remover e ordenar camadas como operações de responsabilidade explícita.
- **FR-004**: Exigir bounds dentro da tela para executar/exportar; reduzir tela não pode reposicionar componentes silenciosamente.
- **FR-005**: Duplicar IDs e remapear vínculos internos, preservando destinos externos; histórico e transações entre módulos ficam na aplicação.

### Key Entities

Screen, Component, Bounds e Style; ordem no array é z-order do fundo para frente. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todas as operações de referência preservam inteiros e IDs corretos; zoom/DPI não alteram o modelo; duplicação não mantém IDs originais.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-002, FR-004 a FR-008, FR-013, FR-016. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [component-catalog.md](../001-gerador-jornada/contracts/component-catalog.md).
- Não implementa Fyne, coordenadas de ponteiro, pilha undo/redo, serialização, roteiro ou exportação.
- Esta especificação define o alvo; evidências de implementação e testes estão no quickstart da feature coordenadora, sem homologação implícita.

## Especificações individuais dos componentes

O detalhamento abaixo mantém este domínio como responsável pelo modelo e pelas regras compartilhadas. Cada feature tem contrato visual próprio; sua existência documental não conclui implementação ou homologação.

- [Data e hora — DateInfo](../014-data-hora/spec.md) — [DESIGN.md](../014-data-hora/DESIGN.md).
- [Dados da lotérica — LotericaInfo](../015-dados-loterica/spec.md) — [DESIGN.md](../015-dados-loterica/DESIGN.md).
- [Indicadores — InfraInfo](../016-infraestrutura/spec.md) — [DESIGN.md](../016-infraestrutura/DESIGN.md).
- [Cabeçalho informativo da tela](../017-cabecalho-tela/spec.md) — [DESIGN.md](../017-cabecalho-tela/DESIGN.md).
- [Faixa de título da página](../018-faixa-titulo/spec.md) — [DESIGN.md](../018-faixa-titulo/DESIGN.md).
- [Área de desenho — Design](../019-canvas-edicao/spec.md) — [DESIGN.md](../019-canvas-edicao/DESIGN.md).
- [Elemento editável — PageElement](../020-elemento-edicao/spec.md) — [DESIGN.md](../020-elemento-edicao/DESIGN.md).
- [Lista de telas — PageTab](../022-lista-telas/spec.md) — [DESIGN.md](../022-lista-telas/DESIGN.md).
- [Estrutura — StructureTab](../024-arvore-estrutura/spec.md) — [DESIGN.md](../024-arvore-estrutura/DESIGN.md).

## Atualização solicitada em 2026-10-18 — Seleção e composição de texto

- **FR-006**: Selecionar múltiplos componentes por Ctrl/Shift+clique, caixa de seleção e Ctrl+A; seleção é estado de interface e não deve ser persistida.
- **FR-007**: Mover uma seleção inteira preservando distâncias relativas e limitando o deslocamento pelo conjunto; desfazer o movimento inteiro em uma ação.
- **FR-008**: Persistir alinhamento horizontal do texto como esquerda, centro ou direita, independente do alinhamento geométrico do elemento; projetos anteriores sem a propriedade usam esquerda.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-006, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-007, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
3. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

## Exclusão confirmada — 2026-10-18

- **FR-009**: Remover o conjunto confirmado em uma única transação reversível, preservando identidades/ordem dos demais componentes e referências externas diagnosticáveis. Seleção, menu contextual, disponibilidade do botão e confirmação Sim/Não pertencem ao adaptador; não inserir widgets ou estado de confirmação no modelo de domínio.

**Aceitação**: dado um conjunto de dois componentes selecionados, confirmar sua exclusão remove ambos em um comando e desfazer restaura ambos; cancelar a confirmação não gera revisão ou histórico de domínio.

## Visualização individual — 2026-10-18

- **FR-010**: A composição de uma tela deve poder ser visualizada isoladamente, preservando seus componentes, dimensões, recursos e conteúdo confirmado. O escopo de visualização e os valores temporários são responsabilidade da aplicação/adaptador; a tela de domínio não recebe estado de janela ou navegação de prévia.

## Identidade de arquivo — 2026-10-18

- **FR-011**: Screen e Component têm FileID editável, independente de ID UUID e Name. Exigir nome-base ASCII alfanumérico de até 40 caracteres, início em letra e sem nomes reservados Windows; Screen exige Pnl seguido de nome. Component não exige Pnl. Nome amigável e textos continuam Unicode. Alterar FileID não troca identidades, navegação, teclado ou referências de roteiro.
- **FR-012**: Garantir FileID de telas único no projeto e de componentes único na tela, sem diferenciar caixa. Duplicação/criação pode sugerir nome livre com sufixo dentro de 40 caracteres; edição explícita inválida deve ser rejeitada atomicamente. Troca de Name não altera FileID já confirmado.

**Aceitação**: Identificação sugere PnlIdentificacao; campos de 40 caracteres são válidos e de 41 são recusados; Pnl sozinho não é tela válida. Copiar controle/tela no limite não cria colisão nem ultrapassa 40. Editar somente FileID preserva UUID e vínculos; Undo restaura o ID anterior.
