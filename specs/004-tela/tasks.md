# Tasks: Domínio Tela

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: exigidos pelo plano aprovado. Todos os itens são futuros e permanecem desmarcados.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [ ] T001 Conferir campos, invariantes e fronteiras com o contrato em specs/004-tela/data-model.md; resolver divergência antes de criar código.

## Phase 2: Foundational

- [ ] T002 Preparar fixtures e casos de contrato de Tela em internal/screen/model_test.go: bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto.
- [ ] T003 Implementar tipos/valores coesos em internal/screen/model.go conforme data-model.md, sem importar UI/IO nem duplicar entidades de outro pacote.

## Phase 3: User Story 1 - Compor uma tela (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Implementar criação/renomeação/dimensionamento de tela, geometria, alinhamento, camada e manipulação de componentes em internal/screen/operations.go, com retornos tipados e sem alteração parcial em erro; separar tipos/arquivos quando necessário.
- [ ] T005 [US1] Verificar cenário válido e invariantes de US1 em internal/screen/operations_test.go, incluindo preservação de IDs/ordem aplicável.

## Phase 4: User Story 2 - Alterar a composição com previsibilidade (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [ ] T006 [US2] Implementar duplicação com mapa de IDs e validação de limites em internal/screen/validation.go, com erros que a aplicação consiga localizar.
- [ ] T007 [US2] Testar os limites e erros de US2 em internal/screen/validation_test.go, demonstrando ausência de publicação parcial.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Revisar GoDoc de todos os elementos e comentários de variáveis em internal/screen/; executar testes do pacote e registrar resultado em specs/004-tela/tasks.md.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Depende dos tipos/definições de Catálogo.

## Parallel Opportunities

Após os tipos base, revisão de testes e validação podem ocorrer em arquivos distintos se as interfaces estiverem estabilizadas. Não paralelizar alterações conflitantes no mesmo modelo. Não adicionar interfaces somente para simular paralelismo.

## Implementation Strategy

Concluir a US1 com fixture mínima, depois US2 e integração. Definição de pronto: requisitos locais cobertos, todos os elementos documentados, testes pertinentes aprovados e nenhuma importação proibida.

## Traceability

FR-001 a FR-005 → T002–T007; SC-001/SC-002 → T005/T007; constituição/documentação → T008. Testes independentes ficam nos pacotes; validação ponta a ponta em quickstart.md da feature coordenadora.
