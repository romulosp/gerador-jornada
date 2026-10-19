# Tasks: Domínio Catálogo

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

## Format: `[ID] [P?] [Story] Description`

IDs Txxx/USx são locais a esta feature. Caminhos de código relativos a apps/gerador-jornada; caminhos specs/ relativos à raiz do workspace. A coordenação de sprints e exportadores permanece em [tasks.md da feature principal](../001-gerador-jornada/tasks.md).

## Phase 1: Setup

- [X] T001 Conferir campos/fronteiras de Catálogo e corrigir o inventário documental. Revisão 2026-10-18 distingue Definition implementado do metamodelo completo ainda planejado, sem importar adaptadores no domínio.

## Phase 2: Foundational

- [ ] T002 Completar fixtures em internal/catalog/model_test.go e adapter/projectzip/package_test.go: já passam união discriminada, enum de botão e round-trip das sete famílias; faltam casos específicos para todos os enums, cores, opções com vírgulas, catálogo incompatível e matriz completa de defaults.
- [X] T003 Implementar Kind, Definition, Option e Properties tipadas em internal/catalog/model.go, com cópia em copy.go. TestPropertiesVariants e TestRoundTripAllFamilies passaram no check de 2026-10-18; pacote não importa UI/IO.

## Phase 3: User Story 1 - Escolher um componente conhecido (Priority: P1)

**Goal**: atender o cenário principal da US1.
**Independent Test**: caso de aceitação US1 executado sem UI e sem geração de saídas.

- [ ] T004 [US1] Completar metadados declarativos para inspector/defaults. Definitions/Lookup estão em internal/catalog/model.go, defaults em screen.NewComponent e formulários no desktop. propertyDefinitions/capabilities do metamodelo não existem integralmente; manter essa evolução explícita, sem mover Fyne para o domínio.
- [ ] T005 [US1] Completar cenário de referência e invariantes de todas as famílias em catalog/model_test.go. O teste atual verifica quantidade/discriminador e alguns enums; round-trip das sete famílias está em projectzip/package_test.go, sem substituir toda a matriz.

## Phase 4: User Story 2 - Editar apenas propriedades válidas (Priority: P2)

**Goal**: comportamento previsível e diagnóstico de erro.
**Independent Test**: cenário US2 com fixture isolada e verificação do estado anterior.

- [ ] T006 [US2] Completar declaração/validação de capacidades por propriedade. Properties.Validate já valida variantes/enums em catalog/model.go; capacidades por perfil/propriedade permanecem parciais e o perfil EFL real está pendente no adaptador (001-gerador-jornada/T055).
- [ ] T007 [US2] Completar a matriz de limites/erros e preservação de estado nos testes do catálogo e da aplicação. TestPropertiesVariants passa, mas não cobre todas as propriedades/combinações inválidas previstas.

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 Concluir revisão integral de GoDoc e variáveis de internal/catalog. Testes do pacote passaram em 2026-10-18; revisão documental de toda declaração e complementação da matriz continuam pendentes.

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

FR-006 (ID comum): 001-gerador-jornada/T093–T095 e 004-tela/FR-011–FR-012; screen/file_id.go, application/file_ids.go e desktop/file_id.go. O catálogo continua sem persistir FileID em Properties; validação fica na composição.
