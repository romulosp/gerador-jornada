# Feature Specification: Faixa de título da página

**Feature Branch**: `018-faixa-titulo` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [004-tela](../004-tela/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Design/page.css). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Faixa de título é preset de largura da tela × 28. Cria título comum com decoração incorporada, texto branco de tamanho 18, família Arial e alinhamento à direita; não muda os bounds ao alinhar o texto.

**Pendências e limites:** Esse preset não se confunde com a faixa interna do menu nem com o título simples de tamanho 14. Equivalência completa com os tokens Java e validação de texto longo/DPI continuam pendentes.

**Fontes verificadas:** [presets.go](../../apps/gerador-jornada/internal/application/presets.go), [presets.go](../../apps/gerador-jornada/internal/adapter/desktop/presets.go), [background.go](../../apps/gerador-jornada/internal/adapter/builtinassets/background.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero adicionar uma faixa decorativa com título de página editável.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve inserir a faixa original com um título editável por cima.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve mostrar, ocultar ou remover a faixa como um componente explícito, sem gerar texto invisível no modelo.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve alterar título e largura preservando separação entre recurso decorativo e texto.

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

- Sem imagem, título muito longo e página estreita. Remoção deve deixar de exportar tanto a faixa quanto o texto removido.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Inserir a faixa original com um título editável por cima.
- **FR-002**: Mostrar, ocultar ou remover a faixa como um componente explícito, sem gerar texto invisível no modelo.
- **FR-003**: Alterar título e largura preservando separação entre recurso decorativo e texto.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Texto do título, faixa de fundo, posição, dimensões, fonte e cor.
- **Estado de interação**: Visível; removida; texto vazio; selecionada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `018-faixa-titulo/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 004-tela.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Alinhamento do texto da faixa

- **FR-008**: Oferecer alinhamento horizontal do texto da faixa: esquerda, centro e direita; preservar X, Y, largura e altura da faixa ao mudar somente essa propriedade.
- **FR-009**: Conservar alinhamento ao salvar/reabrir e reproduzi-lo no editor, prévia, demonstração e PPT, mantendo texto editável. A faixa pronta inicia com texto à direita.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

### Cenário concreto do texto

Dada uma faixa em X=0, Y=100, largura=800 e altura=28, escolher Centro no campo Alinhamento do texto e Salvar alterações. O texto aparece centralizado e os quatro valores geométricos não mudam. Repetir para Esquerda/Direita, salvar/reabrir e conferir o mesmo alinhamento no PPT. Undo restaura o alinhamento anterior.

## ID das instâncias

Toda instância visual desta responsabilidade possui ID de arquivo conforme [Tela](../004-tela/spec.md#identidade-de-arquivo--2026-10-18) e [Catálogo](../003-catalogo/spec.md#atributo-compartilhado-das-instâncias). Campo comum no [painel de propriedades](../028-propriedades/spec.md#campo-id-comum), sem prefixo Pnl obrigatório para componentes e sem substituir UUIDs ou conteúdo visual. Em conjuntos prontos, a regra vale para cada controle efetivamente criado.

## Aparência pós-v1 — 2026-10-19

- **FR-010**: Resolver faixas/títulos aplicáveis pela biblioteca skin Java/EFL e recursos/hash comprovados, com fonte/face/RGBA/alinhamento/espaçamento e texto editável separado da decoração. Preservar IDs/bounds/título/overrides antigos; textAlign vazio antigo continua left, novo preset continua right. Trocar skin não redimensiona nem recolore pelo tema da chrome.

**Aceitação**: variantes/recursos incompletos continuam bloqueados; round-trip v3 e editor/prévia/PPT mantêm alinhamento independente de geometria, título editável e imagem correta por hash. Dono [004/T011](../004-tela/tasks.md), matriz C23, consumidores 001/T111–T116; manuais/aceite T120–T122. Sem homologação fictícia ou reativação Java.
