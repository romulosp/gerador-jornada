# Feature Specification: Formulários de propriedades — Properties

**Feature Branch**: `028-propriedades` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Draft — especificação individual; aceite funcional e visual pendente
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Properties/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero personalizar conteúdo, geometria e aparência com validação antes de aplicar.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve personalizar cor, fonte, preenchimento, borda e imagens quando suportados pelo componente.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve aplicar todos os valores válidos em uma ação reversível e conservar o estado anterior em caso de erro.

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

- Inteiro inválido, cor fora do formato, nome vazio, referências inválidas e texto com caracteres especiais. O Properties legado era um placeholder; o comportamento funcional estava em ToolsRight.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
- **FR-002**: Personalizar cor, fonte, preenchimento, borda e imagens quando suportados pelo componente.
- **FR-003**: Aplicar todos os valores válidos em uma ação reversível e conservar o estado anterior em caso de erro.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Valores editáveis do modelo e rascunho de formulário separado; tipos de entrada apropriados a número, texto, opção e cor.
- **Estado de interação**: Não alterado; alterado; inválido; aplicado; seleção trocada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `028-propriedades/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Confirmação de alterações

- **FR-008**: Manter botão Salvar alterações visível fora da rolagem do formulário, inclusive quando o último campo estiver fora da área visível.
- **FR-009**: Aplicar o rascunho válido em uma ação reversível e atualizar a tela imediatamente após salvar; Descartá-lo pelo botão Descartar restaura os valores confirmados.
- **FR-010**: Exibir Alinhamento do texto, com esquerda, centro e direita, em controles textuais. A alteração modifica texto dentro da caixa, sem mover o elemento.
- **FR-011**: Com seleção múltipla, mostrar quantidade e comandos coletivos; não exibir formulário de um único componente como se editasse todos.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
3. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-010, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
4. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-011, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

### Cenário concreto do rodapé

Selecionar um componente com formulário maior que o painel, rolar até o início e alterar o texto. Salvar alterações continua visível na base do painel; clicar nele publica uma revisão e atualiza o canvas. Alterar novamente e clicar Descartar restaura o último texto confirmado, sem nova revisão. Uma largura inválida não publica nenhum campo do rascunho.

## Campo ID comum

- **FR-012**: Exibir campo ID nas propriedades de tela e de todas as famílias, independente de Nome. Explicar nome-base sem extensão e limite de 40 caracteres; exemplificar PnlIdentificacao na tela. Remover acentos/espaços/símbolos imediatamente na edição; não truncar silenciosamente valor digitado acima do limite. Validar prefixo Pnl somente para telas, início em letra e nomes Windows reservados.
- **FR-013**: Salvar alterações deve confirmar FileID, validar colisão em seu contexto e produzir uma ação reversível, sem mudar UUID ou texto exibido do componente. Se inválido, manter o último estado do projeto e informar o motivo; Descartar retorna ao ID confirmado.
