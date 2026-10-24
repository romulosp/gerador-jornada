# Feature Specification: Domínio Demonstração

**Feature Branch**: `006-demonstracao` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Roteiros finitos, interpretação temporal, cursor/clique, texto Unicode por grafemas e navegação explícita; player e exportador usam as cenas dos componentes existentes. Compostos aparecem no estado confirmado/inicial.

**Pendências e limites:** typeText permanece restrito aos campos suportados; não opera entrada composta ou árvore de menu. Relógio controlado, reordenação de todos os roteiros, expansão de capacidades e homologação das animações no Office continuam pendentes.

**Fontes verificadas:** [model.go](../../apps/gerador-jornada/internal/scenario/model.go), [playback.go](../../apps/gerador-jornada/internal/application/playback.go), [export_test.go](../../apps/gerador-jornada/internal/adapter/pptx/export_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Definir um percurso demonstrativo (Priority: P1)

Interpretar roteiro explícito como sequência finita de estados e eventos temporizados, reutilizável por prévia e PPT.

**Why this priority**: Uma demonstração precisa repetir cliques, digitação e mudanças de tela com ordem e tempos previsíveis, sem gravar vídeo.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** um roteiro de A para B, **When** suas etapas são interpretadas, **Then** clique produz destaque visual e navigate segue o binding explicitamente, sem disparo oculto de navegação.

---

### User Story 2 - Reproduzir tempo e revisitas (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um roteiro A → B → A, **When** ele é reproduzido, **Then** são produzidas três ocorrências e a revisita a A começa com valores iniciais; pausa não consome tempo.

### Edge Cases

Cobrir intervalo temporal, campos por kind, alvo removido, tela errada, A-B-A, auto-navegação, typeText vazio, pausa e reinício. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Representar Scenario com UUID, nome, inicial e steps ordenados; lista vazia salvável, mas não demonstrável.
- **FR-002**: Suportar moveCursor, click, typeText, wait e navigate com duração inteira de 1 a 600000 ms.
- **FR-003**: Validar tela corrente e alvo em cada etapa; click é visual, navigate muda tela ao fim da duração usando a ação da jornada.
- **FR-004**: Produzir ocorrências distintas para toda navegação, incluindo revisita e auto-navegação, com valores e cursor reinicializados.
- **FR-005**: Produzir revelação de texto por cluster Unicode; valores são runtime, sem alterar o projeto.
- **FR-006**: Permitir reprodução/pausa/continuação/reinício via relógio controlado pela aplicação, sem depender de Fyne ou Open XML.

### Key Entities

Scenario, Step, ocorrência, evento temporizado e estado de execução; os três últimos são derivados e não persistidos. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A sequência de referência produz estados e tempos determinísticos; pausa e reinício não vazam estado; toda etapa inválida identifica seu ID e a causa.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-017 a FR-021. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md).
- Não grava uso do mouse, não cria macros, não contém lógica condicional e não emite XML PPT ou eventos Fyne.
- Evidências da base implementada estão no backlog/quickstart da coordenadora; esta atualização é documental e não executa novos testes.

## Atualização solicitada em 2026-10-18 — Transições da apresentação

- **FR-007**: A exportação permite fade, deslizar, revelar ou nenhuma transição, com fade de 300 ms por padrão; não substituir animações de cursor, clique e digitação por transições entre slides.
- **FR-008**: Transições não criam avanço automático na exploração nem no menu/final; ocorrências da demonstração conservam o tempo de avanço do roteiro. Homologar o custo visual de entrada separadamente dos tempos das etapas.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-007, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

## Expansão pós-v1 autorizada — 2026-10-19

- **FR-009**: Validar alvos por capacidades tipadas do catálogo 2; typeText alcança TextArea/PasswordField somente quando contratados e com conteúdo ilustrativo/mascaramento suportado. Click continua visual e navigate explícito apenas de button, inclusive twoLabel; não mudar toggle/seleção/modal por click demonstrativo implícito.
- **FR-010**: Produzir quadros/cena determinísticos com fonte/skin e estados locais/recorte aplicáveis, reiniciados por visita/restart, isolados do documento. Novas famílias sem operação de roteiro suportada recebem diagnóstico de capacidade; não prometer gravação/interação ausente nem introduzir relógio real ou negócio externo.

**Aceitação**: roteiro com novas entradas e A→B→A preserva tempos/ordem/Unicode, restaura estados e produz o mesmo quadro na prévia/PPT; click em toggle não o altera implicitamente. Nova aparência não muda IDs de etapas/rotas nem redimensiona componentes. FR-009–FR-010 ampliam os alvos originalmente limitados a textField sem alterar a semântica dos cinco kinds de Step. T009–T010 locais; 001/T114–T116 integra consumidores; homologação/aceite permanecem separados, Java pausado.

<a id="codigos-barras-2026-10-23"></a>
## Incremento planejado — representação dos códigos de barras (2026-10-23)

- **FR-011**: Preservar no quadro demonstrativo a configuração de [031](../031-codigo-barras-efl/spec.md) e [032](../032-codigo-barras-processado-efl/spec.md), com ordem/textos/estilo e sem importar estado da prévia livre. Entrada começa vazia e ficha processada mostra seus valores configurados. Nesta entrega, os compostos não ganham `typeText`, processamento por click ou nova etapa de roteiro: uma etapa que exigir capacidade não suportada gera diagnóstico localizado. PPT mostra formas/textos editáveis, não executa validação, leitora, combo de parametrização ou consulta bancária.

**Aceitação**: mesma configuração e instante produzem o mesmo quadro em prévia/PPT; revisita e restart não herdam código digitado na prévia livre. Roteiro tentando digitar na ficha de leitura ou na entrada sem contrato de etapa é recusado com referência de etapa, não simulado como sucesso. Exploração estática continua exportável se válida; roteiro selecionado inválido não é omitido silenciosamente. [T011](tasks.md#tarefas-codigos-barras-2026-10-23) verifica capacidades/quadros; cena/exportador são 001/T137. Digitação animada de boleto ou interação bancária ficam fora deste incremento.

<a id="menu-aplicacao-2026-10-23"></a>
## Incremento autorizado — quadro estático de menu

- **FR-012**: Produzir quadro de [033](../033-menu-aplicacao-efl/spec.md) a partir da árvore/inicial confirmados, nunca do caminho da prévia livre. Seleção inicial pode ser raiz ou submenu existente; repetir snapshot/tempo produz o mesmo quadro, com títulos/ordem/estilo/enabled. Não oferecer typeText no menu ou nos compostos 031/032; metadados de movimento não criam etapas ou ações. Clique de roteiro continua visual; extensão por item depende de decisão explícita da coordenação e contrato/capacidades sincronizados antes de habilitar.

**Aceitação**: revisita/restart e exploração livre não contaminam quadro exportado; inicial removido/ausente gera erro localizado. PPT estático mantém textos/formas editáveis, sem afirmar interação de árvore. Roteiro selecionado com capacidade não suportada é recusado, não omitido. [T012](tasks.md#tarefas-menu-aplicacao-2026-10-23) cobre dados/capacidades; cena/exportador em 001/T142.
