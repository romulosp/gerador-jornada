# Feature Specification: Domínio Projeto

**Feature Branch**: `002-projeto` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — domínio especificado; implementação não iniciada
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preservar a identidade do documento (Priority: P1)

Definir o documento de uma jornada, sua identidade, revisão e recursos, sem conhecer ZIP, janela ou exportadores.

**Why this priority**: O usuário precisa continuar e transportar o trabalho sem perder referências ou depender do navegador.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** um projeto com telas e roteiros, **When** sua identidade, nome ou revisão são consultados após edição e persistência, **Then** os IDs permanecem estáveis e a revisão salva representa exatamente o snapshot gravado.

---

### User Story 2 - Gerenciar recursos incorporados (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um recurso utilizado por duas telas, **When** uma tela é removida, **Then** o recurso ainda referenciado permanece disponível; dados de recurso ausente ou alterado geram diagnóstico localizado.

### Edge Cases

Cobrir IDs duplicados, nome vazio, revisão negativa, recurso compartilhado e isolamento de snapshot. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Representar projeto com UUID estável, nome não vazio, revisão não negativa e coleções ordenadas de telas/roteiros/recursos.
- **FR-002**: Manter UUIDs internos únicos globalmente no projeto; mudanças de nome não alteram referências.
- **FR-003**: Representar imagens PNG/JPEG incorporadas por ID, caminho interno, mediaType, hash SHA-256 e tamanho, sem URLs externas.
- **FR-004**: Permitir documento sem telas como rascunho; sua ausência de inicial deve ser resolvida antes de prévia/exportação.
- **FR-005**: Fornecer visão consistente para snapshot sem compartilhar mapas ou slices mutáveis com edição; autosave e IO pertencem à aplicação/adaptador.

### Key Entities

Project e Asset; metadados JavaOptions são dados do documento, validados pelo alvo Java. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: IDs e ordem preservados em 100% das operações do conjunto de referência; snapshot anterior não muda após nova edição.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-001, FR-009, FR-010; NFR-004, NFR-005. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [project-format.md](../001-gerador-jornada/contracts/project-format.md).
- Não implementa filesystem, ZIP, autosave, undo, UI, decodificação de imagens ou geração de artefatos.
- Esta entrega contém especificação, não código implementado ou testes executados.

## Persistência de IDs de arquivo

- **FR-006**: Gravar manifesto schemaVersion 2 com fileId de telas/componentes. Aceitar schemaVersion 1 por migração explícita: preencher somente IDs ausentes com sugestões sem acentos e sufixos determinísticos para colisões; preservar UUIDs, nomes e referências. Em v2, exigir valores válidos e recusar duplicidade no escopo. Não renomear IDs explícitos silenciosamente. Leitores antigos devem recusar a versão 2 como incompatível.

## Recuperação nos alvos de distribuição

- **FR-007**: O adaptador desktop deve resolver a pasta de recuperação por usuário com APIs nativas, mantendo o caminho anterior no Windows. Linux usa configuração XDG e macOS Application Support. Falha de resolução deve informar indisponibilidade, sem gravar caminho relativo no diretório corrente; Project/Document não carregam caminhos do sistema operacional.
