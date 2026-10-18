# Feature Specification: Domínio Jornada

**Feature Branch**: `005-jornada` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — domínio especificado; implementação não iniciada
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Definir o fluxo entre telas (Priority: P1)

Modelar tela inicial e navegação de botões como dados explícitos, compartilhados pela prévia e pelas saídas.

**Why this priority**: A navegação funcional da prévia legada não é transportada aos listeners Java; a nova versão precisa de um contrato comum.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** telas A e B e um botão de A, **When** a ação navigate aponta a B, **Then** o resolvedor retorna B pelo ID; ausência de ação equivale a none.

---

### User Story 2 - Identificar destinos inválidos (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um destino removido, **When** a jornada é validada, **Then** o diagnóstico identifica origem e destino ausente; o grafo inválido não é silenciosamente corrigido.

### Edge Cases

Cobrir none implícito, ciclo, auto-navegação, fonte não botão, destino apagado, inicial ausente e binding duplicado. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Manter startScreenId opcional no rascunho e obrigatório/existente para execução.
- **FR-002**: Permitir somente ações none e navigate associadas ao click de button, no máximo uma por fonte.
- **FR-003**: Exigir destino existente em navigate e destino nulo em none; fonte pertence à tela informada.
- **FR-004**: Permitir ciclos e auto-navegação; emitir warning de tela inalcançável sem excluí-la.
- **FR-005**: Expor resolução e validação pura, sem cliques de UI, chamadas externas, condições baseadas em campos ou regras de negócio.

### Key Entities

Journey e ButtonAction com origem, kind e destino por IDs. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os vínculos do conjunto A → B → C → A resolvem corretamente; destino apagado sempre produz diagnóstico de origem; ciclos válidos não são recusados.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-011 a FR-013, FR-014, FR-023. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md).
- Não armazena valores de campo, histórico de navegação, widgets, slides ou listeners Java.
- Esta entrega contém especificação, não código implementado ou testes executados.
