# Tasks: Domínio Jornada

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [X] T001 Conferir Journey/Action, referências e divisão entre modelo e aplicação. Revisão 2026-10-18 registra Bind/Destination em journey/model.go, validação entre telas em application/document.go e warning de inalcançável ainda pendente em T006.

## Phase 2: Foundational

- [ ] T002 Completar fixtures em journey/model_test.go e application/session_test.go: auto-navegação, troca de binding e reset de visita já passam; faltam casos dedicados para todas as origens inválidas, destino removido, ausência de inicial e warning de inalcançável.
- [X] T003 Implementar Journey e Action em internal/journey/model.go, sem dependência de adaptadores ou Screen. TestBindAndCycle passou no check de 2026-10-18.

## Phase 3: User Story 1 - Definir o fluxo entre telas (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [X] T004 [US1] Implementar campo de inicial, binding atômico e resolução por ID em journey/model.go, com validação entre domínios em application/document.go e edição em desktop/navigation.go. TestBindAndCycle/TestPreviewReset aprovados em 2026-10-18; não inclui warning de alcançabilidade de T006.
- [ ] T005 [US1] Completar o cenário A-B-C-A e invariantes de identidade/ordem. O núcleo já é exercitado por TestBindAndCycle e TestPreviewReset; falta a matriz completa de vínculos/removal/erros especificada.

## Phase 4: User Story 2 - Identificar destinos inválidos (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [ ] T006 [US2] Completar diagnóstico de referências e alcançabilidade. Document.Validate já identifica inicial/origem/destino inválidos e binding duplicado; o warning de tela inalcançável de FR-004 ainda não é emitido. Localização clicável permanece em 001-gerador-jornada/T030.
- [ ] T007 [US2] Completar testes de todos os limites/erros da validação cruzada e ausência de publicação parcial em application/session_test.go, além de journey/model_test.go. Não confundir Bind básico com homologação da jornada completa.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Concluir revisão integral de GoDoc/variáveis em internal/journey e fronteiras de validação. Checks passaram em 2026-10-18; os testes e a revisão ainda não cobrem integralmente a definição de pronto do domínio.

## Dependencies & Execution Order

T001 → T002/T003 → T004/T006 → T005/T007 → T008. Preparar teste significativo antes da implementação correspondente e verificar falha pertinente quando viável.

Recebe referências resolvidas pela aplicação, sem dependência de adaptadores.

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
