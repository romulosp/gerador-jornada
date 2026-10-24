# Feature Specification: Painel direito — ToolsRight

**Feature Branch**: `027-painel-direito` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/ToolsRight/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Painel direito usa formulário em área com rolagem nos dois eixos e rodapé fixo Salvar alterações/Descartar. O divisor faz parte do layout Fyne; conteúdo longo não é um novo dado do projeto.

**Pendências e limites:** A suíte histórica inspector_split_test.go não está na árvore atual. Há implementação do layout, mas falta repetir sua matriz específica de largura/rolagem/foco em janela nativa/DPI.

**Fontes verificadas:** [editor.go](../../apps/gerador-jornada/internal/adapter/desktop/editor.go), [inspector.go](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero apresentar somente as propriedades pertinentes à seleção atual.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve mostrar configuração de tela quando nenhum controle estiver selecionado.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve mostrar propriedades comuns e da família correta ao selecionar um controle.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve separar formulário, validação e aplicação transacional, com rolagem própria.

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

- Trocar seleção com rascunho, componente removido, propriedade não suportada e campos fora da área visível.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Mostrar configuração de tela quando nenhum controle estiver selecionado.
- **FR-002**: Mostrar propriedades comuns e da família correta ao selecionar um controle.
- **FR-003**: Separar formulário, validação e aplicação transacional, com rolagem própria.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

**Clarificação do requisito existente de largura ajustável e rolagem (FR-003/FR-008) — 2026-10-21, sem novo FR:** o divisor entre canvas e painel direito deve alterar a largura efetiva disponível, inclusive com seletores, caminhos de fontes e versões de perfil extensos. O corpo deve rolar nos dois eixos, desacoplando suas dimensões mínimas do painel sem truncar valores ou perder acesso horizontal aos campos. Cabeçalho e rodapé Salvar alterações/Descartar continuam fixos e fora da rolagem; seus controles definem naturalmente a largura mínima. Não exigir largura fixa, estreitar botões ou prometer colapso ilimitado. Redimensionar/rolar preserva documento, histórico, rascunho e semântica de Salvar/Descartar.

**Aceitação da clarificação (pendente):** na janela 1280 × 720, conferir tela, cada uma das 16 famílias e seleção múltipla; mover o divisor e observar largura/posição reais, não somente `Offset`. Com valores longos, alcançar os campos nos dois eixos e manter cabeçalho/rodapé imóveis. Confirmar Salvar/Descartar pelos botões renderizados, sem criar ações individuais fictícias na seleção múltipla. Quatro testes preparados em inspector_split_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_split_test.go`), não executados por política Windows, sem contorno. **Evidência final do executor principal (2026-10-21):** gofmt em dois arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` concluiu a compilação de 15 suítes, **sem executá-las**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade PASS: CRC, todos os SHA256SUMS, sidecar SHA-256, PE amd64, ambas as cópias do EXE idênticas, README definitivo, manifesto com aprovações pendentes e ausência de backups terminados em `~`. Identidades finais e evidências no [quickstart, seção 11](../001-gerador-jornada/quickstart.md#11-redimensionamento-do-painel-de-propriedades--2026-10-21). [001/T129](../001-gerador-jornada/tasks.md#phase-23-redimensionamento-do-painel-de-propriedades--2026-10-21) permanece aberta como verificação mista: sem execução dos testes Go, validação de UI/arraste nativos, DPI, revisão manual ou aceite; `manualApproval` e `redistributionApproval` continuam `pending`.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Rascunho de formulário e foco transitórios; valores confirmados pertencem à tela ou ao componente.
- **Estado de interação**: Propriedades da tela; de componente; sem tela; rascunho inválido; confirmado.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `027-painel-direito/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Ação de confirmação sempre acessível

- **FR-008**: Reservar rodapé fixo para Salvar alterações e Descartar, ligado ao formulário corrente e fora da sua rolagem.
- **FR-009**: Distinguir propriedade de texto, geometria do elemento e comandos para seleção múltipla; trocar seleção não deve publicar rascunho incompleto.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
