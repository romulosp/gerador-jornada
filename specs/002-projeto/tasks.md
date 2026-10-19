# Tasks: Domínio Projeto

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [X] T001 Conferir o modelo e as fronteiras com data-model.md; registrar estado real e lacunas sem alterar requisitos. Evidência: revisão documental de 2026-10-18; metadados em project/model.go, composição em application/document.go e manifesto em projectzip/wire.go.

## Phase 2: Foundational

- [ ] T002 Completar a matriz de fixtures de Projeto em internal/project/model_test.go, application/session_test.go e adapter/projectzip/package_test.go. Já existem UUID/nome, snapshot/rollback, round-trip e migração. Faltam validação/casos para UUID canônico e unicidade global de roteiros/etapas entre entidades (001-gerador-jornada/T008), além dos demais IDs duplicados, limites e recursos compartilhados previstos.
- [X] T003 Implementar tipos de identidade e metadados em internal/project/model.go, sem UI/IO. Project/Asset, New/NewID/ValidID/Validate existentes; TestIdentity aprovado no check de 2026-10-18. Composição fica em application.Document, sem duplicar entidades.

## Phase 3: User Story 1 - Preservar a identidade do documento (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Completar operações de identidade/revisão e índice de recursos. New/Validate e Session.Change já existem em project/model.go e application/session.go; o índice/referências e todas as operações atômicas previstas ainda não têm API/testes completos. Não criar operations.go artificialmente.
- [ ] T005 [US1] Completar cenários válidos e invariantes em project/model_test.go e application/session_test.go. TestIdentity/TestHistoryAndSnapshot cobrem o núcleo; faltam as combinações de identidade/ordem/recursos da matriz original.

## Phase 4: User Story 2 - Gerenciar recursos incorporados (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [ ] T006 [US2] Completar integridade e remoção segura de recursos. adapter/projectzip/package.go já valida bytes/hash/caminhos; application.Document.Copy preserva snapshot. Remoção segura e índice de uso ainda precisam de regra e cobertura; não concluir apenas com leitura ZIP.
- [ ] T007 [US2] Ampliar testes de falha de recurso e ausência de publicação parcial em adapter/projectzip/package_test.go e no adaptador atomicfile. Round-trip, campo desconhecido e traversal já passam; corrupção/hash/limites e falhas de publicação precisam de cenários específicos.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Concluir revisão integral de GoDoc/variáveis em internal/project e fronteiras de composição. O pacote compila e TestIdentity passou no check de 2026-10-18; isso não comprova revisão de toda declaração nem toda a cobertura planejada.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Pode iniciar o modelo após a preparação da toolchain; integração depende dos checkpoints da feature principal.

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

FR-006 (manifesto v2/migração): 001-gerador-jornada/T095, projectzip/file_ids_test.go. FR-007 (recuperação nativa): 001-gerador-jornada/T098 e T013/T069/T070 para testes ainda pendentes; desktop/recovery_path.go. Esses requisitos de adaptador não criam I/O em Project.
