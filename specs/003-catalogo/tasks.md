# Tasks: Domínio Catálogo

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: exigidos pelo plano aprovado. Todos os itens são futuros e permanecem desmarcados.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [ ] T001 Conferir campos, invariantes e fronteiras com o contrato em specs/003-catalogo/data-model.md; resolver divergência antes de criar código.

## Phase 2: Foundational

- [ ] T002 Preparar fixtures e casos de contrato de Catálogo em internal/catalog/model_test.go: sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível.
- [ ] T003 Implementar tipos/valores coesos em internal/catalog/model.go conforme data-model.md, sem importar UI/IO nem duplicar entidades de outro pacote.

## Phase 3: User Story 1 - Escolher um componente conhecido (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Implementar resolução de famílias, defaults e metadados do inspector em internal/catalog/operations.go, com retornos tipados e sem alteração parcial em erro.
- [ ] T005 [US1] Verificar cenário válido e invariantes de US1 em internal/catalog/operations_test.go, incluindo preservação de IDs/ordem aplicável.

## Phase 4: User Story 2 - Editar apenas propriedades válidas (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [ ] T006 [US2] Implementar propriedades tipadas e capacidades declaradas em internal/catalog/validation.go, com erros que a aplicação consiga localizar.
- [ ] T007 [US2] Testar os limites e erros de US2 em internal/catalog/validation_test.go, demonstrando ausência de publicação parcial.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Revisar GoDoc de todos os elementos e comentários de variáveis em internal/catalog/; executar testes do pacote e registrar resultado em specs/003-catalogo/tasks.md.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Pode iniciar o modelo após a preparação da toolchain; integração depende dos checkpoints da feature principal.

## Parallel Opportunities

Após os tipos base, revisão de testes e validação podem ocorrer em arquivos distintos se as interfaces estiverem estabilizadas. Não paralelizar alterações conflitantes no mesmo modelo. Não adicionar interfaces somente para simular paralelismo.

## Implementation Strategy

Concluir a US1 com fixture mínima, depois US2 e integração. Definição de pronto: requisitos locais cobertos, todos os elementos documentados, testes pertinentes aprovados e nenhuma importação proibida.

## Traceability

FR-001 a FR-005 → T002–T007; SC-001/SC-002 → T005/T007; constituição/documentação → T008. Testes independentes ficam nos pacotes; validação ponta a ponta em quickstart.md da feature coordenadora.
