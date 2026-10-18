# Feature Specification: Domínio Demonstração

**Feature Branch**: `006-demonstracao` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — domínio especificado; implementação não iniciada
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

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
- Esta entrega contém especificação, não código implementado ou testes executados.
