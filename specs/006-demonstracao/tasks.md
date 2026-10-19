# Tasks: Domínio Demonstração

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [X] T001 Conferir Scenario/Step/Event/Occurrence e separação entre interpretação, relógio e exportação. Revisão 2026-10-18 registra a implementação e a separação temporal ainda pendente, sem declarar o player inteiramente na aplicação.

## Phase 2: Foundational

- [ ] T002 Completar fixtures em scenario/model_test.go e testes da aplicação/player. TestRevisitAndUnicode cobre revisita, grafemas e duração zero; faltam UUID canônico/global de roteiros/etapas entre entidades (001-gerador-jornada/T008), todos os kinds/limites, A-B-A, alvos removidos, typeText vazio, pausa/reinício com relógio controlado e matriz completa de falhas.
- [X] T003 Implementar Scenario, Step, Event e Occurrence em internal/scenario/model.go, sem Fyne/Open XML. TestRevisitAndUnicode passou no check de 2026-10-18.

## Phase 3: User Story 1 - Definir um percurso demonstrativo (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Completar comandos tipados/reordenação e separação temporal. Interpret e validação existem em scenario/model.go; desktop/scenarios.go edita e reordena etapas, mas não reordena a lista de roteiros. Consolidar controle temporal fora de player.go conforme FR-006 e 001-gerador-jornada/T044–T045.
- [ ] T005 [US1] Completar cenário de referência e matriz válida/inválida em scenario/model_test.go e testes da aplicação. A cobertura de revisita/Unicode/duração já passa; não cobre todo o contrato de etapas e preservação de estado.

## Phase 4: User Story 2 - Reproduzir tempo e revisitas (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [X] T006 [US2] Implementar interpretação determinística, grafemas, reset e ocorrências em scenario/model.go (Interpret/Graphemes), consumida por application/playback.go e PPT. TestRevisitAndUnicode e TestEditableNavigationAndTiming aprovados em 2026-10-18; relógio controlado e matriz adicional permanecem em T004/T007.
- [ ] T007 [US2] Completar testes de estados/tempos e todas as falhas de referência, pausa e reinício. Há evidência estrutural em scenario/model_test.go e pptx/export_test.go; execução PowerPoint e tolerância real permanecem em 001-gerador-jornada/T054.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Concluir revisão integral de GoDoc/variáveis em internal/scenario e regras de reprodução. Checks passaram em 2026-10-18, sem comprovar a matriz temporal completa ou homologação visual.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Depende de resolvedor de navegação Journey e de índices imutáveis preparados pela aplicação.

## Parallel Opportunities

Após os tipos base, revisão de testes e validação podem ocorrer em arquivos distintos se as interfaces estiverem estabilizadas. Não paralelizar alterações conflitantes no mesmo modelo. Não adicionar interfaces somente para simular paralelismo.

## Implementation Strategy

Concluir a US1 com fixture mínima, depois US2 e integração. Definição de pronto: requisitos locais cobertos, todos os elementos documentados, testes pertinentes aprovados e nenhuma importação proibida.

## Traceability

FR-001 a FR-006 → T002–T007; SC-001/SC-002 → T005/T007; constituição/documentação → T008. Testes independentes ficam nos pacotes; validação ponta a ponta em quickstart.md da feature coordenadora.


## Evidências da reconciliação — 2026-10-18

Executor: agente, ambiente Windows deste workspace. `scripts/check.ps1 -Jdk17
C:\Desenvolvimento\jdk-17.0.11` passou: formatação, vet, testes dos pacotes Go,
três testes do empacotador e nove testes da toolchain Java estática. Resultados por
pacote e limites em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18).

A existência de um teste aprovado não completa os demais cenários descritos nas
caixas abertas. Não foram executados testes manuais, PowerPoint ou EFL real nesta
revisão. IDs T001–T008 foram preservados; caminhos planejados foram substituídos
pelos arquivos existentes, sem criar novos módulos ou duplicar implementação.

## Requisitos acrescentados e integração

FR-007–FR-008 (transições): 001-gerador-jornada/T084 e adapter/pptx/transition_test.go, aprovados estruturalmente. Homologação visual/tempo em PowerPoint permanece em 001-gerador-jornada/T054; o teste de XML não conclui esse gate.
