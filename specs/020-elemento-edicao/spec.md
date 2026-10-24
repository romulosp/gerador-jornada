# Feature Specification: Elemento editável — PageElement

**Feature Branch**: `020-elemento-edicao` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [004-tela](../004-tela/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/PageElement/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Elemento em edição usa a mesma cena de apresentação da jornada, com imagens/textos em sua ordem e contornos/alças acima da composição. Seleção múltipla, movimento e comandos contextuais integram o histórico existente.

**Pendências e limites:** Cobertura de 16 famílias e todos os estados/recortes ainda é alvo de expansão. Fonte visual aplicada não comprova todas as condições de texto longo, zoom ou DPI.

**Fontes verificadas:** [board_renderer.go](../../apps/gerador-jornada/internal/adapter/desktop/board_renderer.go), [scene_render.go](../../apps/gerador-jornada/internal/adapter/desktop/scene_render.go), [component_menu_test.go](../../apps/gerador-jornada/internal/adapter/desktop/component_menu_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero manipular um componente sem confundir sua aparência com as alças do editor.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve selecionar o componente superior sob o ponteiro.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve confirmar um arraste ou resize inteiro como uma única operação de histórico.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve duplicar, remover e mudar camada sem perder vínculos internos válidos.

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

- Componente coberto, arraste além da borda, Escape, campo com foco e clique em área vazia.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Selecionar o componente superior sob o ponteiro.
- **FR-002**: Confirmar um arraste ou resize inteiro como uma única operação de histórico.
- **FR-003**: Duplicar, remover e mudar camada sem perder vínculos internos válidos.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Geometria e ordem persistidas; contorno de seleção, alça e arraste provisório transitórios.
- **Estado de interação**: Não selecionado; selecionado; movendo; redimensionando; operação cancelada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `020-elemento-edicao/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 004-tela.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Seleção e movimentação profissional

- **FR-008**: Ctrl/Shift+clique alterna inclusão na seleção; clique simples em membro selecionado conserva o conjunto para arrastar; clique simples fora do conjunto troca a seleção.
- **FR-009**: Arrastar área vazia desenha caixa de seleção; ao soltar, seleciona os componentes que intersectam a caixa. Ctrl+A seleciona todos; Escape cancela o arraste provisório ou limpa a seleção.
- **FR-010**: Setas movem 1 unidade lógica e Shift+setas movem 10; nenhum atalho de desenho deve capturar digitação em campos com foco.
- **FR-011**: Arraste e setas movimentam todo o conjunto preservando espaçamento; bordas limitam o grupo inteiro. Resize pela alça só é oferecido para seleção única.
- **FR-012**: Exibir cursor de movimentação em cruz sobre componente selecionado e cursor diagonal na alça; não reconstruir todas as propriedades e imagens a cada clique ou movimento.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
3. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-010, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
4. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-011, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
5. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-012, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

- **FR-013**: Ctrl+D duplica a seleção inteira e remapeia vínculos internos; Delete oferece confirmação do conjunto e remove-o em uma ação reversível. Referências externas permanecem diagnosticáveis.

### Cenários concretos de regressão

1. Selecionar A, Ctrl+clicar B e arrastar 20 × 10: os dois mudam pelo mesmo deslocamento, preservam distâncias e um Undo restaura ambos.
2. Arrastar caixa vazia em qualquer direção seleciona A e B que a intersectam; Ctrl+A seleciona todos os componentes da tela, sem editar o documento.
3. Setas movem 1 unidade; Shift+setas, 10. Com campo textual focado, esses atalhos não movem controles.
4. Ctrl+D em campo e teclado cria novos IDs e vincula o teclado copiado ao campo copiado. Delete do conjunto e um Undo removem/restauram todos.
5. Na alça de uma seleção única aparece cursor diagonal; no corpo selecionado, cruz de movimentação. A seleção múltipla não oferece resize coletivo implícito.

## Menu contextual e exclusão — 2026-10-18

- **FR-014**: Botão direito sobre o componente superior no ponto clicado abre Alterar propriedades e Excluir elemento. Converter coordenadas pelo zoom/margem e respeitar camadas. Selecionar o alvo se estiver fora da seleção; manter o conjunto se o alvo já for membro. Em área vazia, não abrir menu de componente.
- **FR-015**: Alterar propriedades deve selecionar somente o componente clicado e direcionar o foco ao formulário do painel direito. A abertura do menu/formulário não modifica dados nem histórico.
- **FR-016**: Delete fora de campos de entrada e Excluir elemento/seleção devem solicitar confirmação explícita Sim/Não antes de chamar a remoção transacional. Não ou fechamento conservam documento, revisão e seleção; Sim remove os alvos apresentados, permite Undo e limpa a seleção. Uma confirmação aberta impede outra para a mesma operação.

### Acceptance Scenarios

1. **Given** um componente sem seleção, **When** clicar com botão direito, **Then** o componente fica selecionado e as duas ações são exibidas; dados permanecem iguais.
2. **Given** o menu do componente, **When** escolher Alterar propriedades, **Then** seu formulário recebe foco; Delete dentro do campo edita texto sem abrir confirmação de exclusão.
3. **Given** um componente selecionado, **When** pressionar Delete e escolher Não, **Then** componente, revisão e histórico permanecem iguais; escolhendo Sim, o componente é excluído e restaurável por Undo.
4. **Given** dois controles selecionados, **When** abrir o menu num membro e confirmar Excluir seleção, **Then** ambos são removidos em uma única transação.
5. **Given** área vazia ou nenhuma seleção, **When** clicar com botão direito ou pressionar Delete, respectivamente, **Then** não há menu/diálogo indevido nem alteração de dados.
