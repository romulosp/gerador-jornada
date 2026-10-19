# Tasks: Domínio Tela

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [X] T001 Conferir entidades/fronteiras de Tela e registrar o mapeamento real em plan.md/data-model.md. Revisão de 2026-10-18 preserva a separação entre composição, comandos da aplicação e interação desktop.

## Phase 2: Foundational

- [ ] T002 Completar fixtures em screen/model_test.go, selection_test.go, file_id_test.go e application/edit_selection_test.go. Já há limites/clamp, seleção/IDs e duplicação; faltam combinações completas de dimensões, nomes, camadas, remapeamento de teclado e remoção referenciada previstas no plano.
- [X] T003 Implementar Screen/Component/Bounds/Style em screen/model.go, regras coletivas em selection.go e IDs em file_id.go. TestGeometryAndLocalNames, TestGroupMovementPreservesSpacing e TestFileIDRules passaram no check de 2026-10-18, sem UI/IO no pacote.

## Phase 3: User Story 1 - Compor uma tela (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Completar operações de criação/geometria/alinhamento/camadas com responsabilidade explícita. Núcleo existe em application/edit_screen.go, transform.go e desktop/arrangement.go/operations.go; conferir todos os alinhamentos previstos e consolidar os comandos que ainda vivem no desktop. Não criar uma classe central ou arquivos apenas para repetir o plano antigo.
- [ ] T005 [US1] Completar matriz de operações válidas e invariantes nos testes existentes de screen/application. Geometria, transação, seleção e IDs já têm evidência; falta cobertura integral de todos os comandos/ordem/camadas e efeitos de remoção.

## Phase 4: User Story 2 - Alterar a composição com previsibilidade (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [X] T006 [US2] Implementar duplicação com mapa de IDs e limites entre os módulos em application/edit_screen.go, edit_selection.go e screen/selection.go. Evidências: TestSelectionDuplicateAndDelete, TestFileIDMigrationAndDuplication e TestGroupMovementPreservesSpacing aprovados em 2026-10-18; casos adicionais seguem em T002/T007.
- [ ] T007 [US2] Completar cenários de falha de duplicação, vínculos de teclado/destinos externos e ausência de alteração parcial nos testes existentes de screen/application. A cobertura atual não fecha toda a matriz original.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Concluir revisão integral de GoDoc/variáveis do domínio Tela e comandos relacionados. Checks do pacote passaram em 2026-10-18; DPI/homologação desktop pertencem aos gates manuais de 001, não a esta caixa.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Depende dos tipos/definições de Catálogo.

## Parallel Opportunities

Após os tipos base, revisão de testes e validação podem ocorrer em arquivos distintos se as interfaces estiverem estabilizadas. Não paralelizar alterações conflitantes no mesmo modelo. Não adicionar interfaces somente para simular paralelismo.

## Implementation Strategy

Concluir a US1 com fixture mínima, depois US2 e integração. Definição de pronto: requisitos locais cobertos, todos os elementos documentados, testes pertinentes aprovados e nenhuma importação proibida.

## Traceability

FR-001 a FR-005 → T002–T007; SC-001/SC-002 → T005/T007; constituição/documentação → T008. Testes independentes ficam nos pacotes; validação ponta a ponta em quickstart.md da feature coordenadora.


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

FR-006–FR-009: 001-gerador-jornada/T080–T083 e T088–T090; screen/selection.go, application/edit_selection.go e testes. FR-010: 001-gerador-jornada/T092 e application/screen_preview_test.go. FR-011–FR-012: 001-gerador-jornada/T093–T096 e testes file_id nos pacotes screen/application/projectzip. São verificações de integração já registradas, não novos comandos a implementar duas vezes.
