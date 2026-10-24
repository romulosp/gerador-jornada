# Tasks: Domínio Demonstração

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


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

## Expansão pós-v1 — 2026-10-19 (tarefas novas abertas)

- [ ] T009 Evoluir índice/validação de alvos e quadros Interpret/FrameAt para capacidades/estados iniciais do catálogo 2, typeText em novas entradas suportadas e recorte/cena comum; preservar click visual/navigate explícito, reset e isolamento do documento. FR-009–FR-010; depende contratos 003/T010–T015 e 001/T111; widgets/PPT ficam em 001/T114–T116.
- [ ] T010 Executar fixtures de Unicode/typeText/máscaras/mascaramento ilustrativo, twoLabel/navigate, click sem mutação de toggle/seleção/modal, A→B→A/reset/restart e determinismo de quadros/tempos/IDs sem resize; testar alvos sem capacidade com diagnóstico de etapa. FR-009–FR-010; depende T009. Registrar revisão/resultados reais; execução PowerPoint/manual/aceite são 001/T121–T122, não aprovados por teste estrutural.

T001–T008 preservados; nenhuma caixa nova concluída nesta atualização documental. Gate Java e requisitos de relógio/reordenação anteriores continuam abertos onde já registrados.

**Evidência pós-Capture fornecida pelo main:** gofmt/vet PASS; **15 suítes Go
compiladas, não executadas**, bloqueio prévio por política Windows preservado;
**28 Python PASS** (15 importador + 4 packaging + 9 Java 6, incluindo erro esperado de
`String.strip`, não homologação EFL); importer `--check` PASS, mesmo inventário
214/8/270/47. Distribuição Windows amd64 **gerada e verificada**; checksum do ZIP no
sidecar externo `.zip.sha256`. Reempacotar documentos preserva o EXE. Identidade do EXE e manual pendente no
[quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
T009/T010 continuam abertos: compilação/build não executam sua matriz nem homologam
Windows/DPI/PowerPoint ou aceite.

<a id="tarefas-codigos-barras-2026-10-23"></a>
## Incremento de domínio — códigos de barras (2026-10-23)

- [ ] T011 [US1] Verificar/ajustar capacidades e fixtures de quadros em internal/scenario e internal/application/playback.go para FR-011, após 003/T018/T020 e 004/T015: entrada vazia e ficha processada com dados confirmados/Unicode/ordem; nenhuma herança de captura de prévia livre em revisita/restart. Recusar typeText nos dois compostos por capacidade com ID da etapa, sem converter para textField, omitir roteiro inválido ou acrescentar Kind de Step. Click continua visual e navigate explícito; representação editável PPT é 001/T137, integração em T138 e manual em T139 da coordenadora. Contrato navigation-scenario.md, 031/FR-011 e 032/FR-008–FR-009.

**T011 permanece aberta pelo critério de verificação:** capacidades em scenario/model.go ajustadas e fixtures em scenario/compound_test.go e application/playback_compounds_test.go presentes, **não executadas**. Entrada permanece vazia, ficha conserva dados/Unicode/ordem; revisita/restart não herdam captura livre. Click é somente visual e typeText é recusado com STEP_CAPABILITY_UNSUPPORTED e ID da etapa, sem saída parcial ou animação de digitação de boleto. T001–T010 e seus estados históricos preservados; Catálogo, persistência, cena e widgets mantêm seus donos.

<a id="tarefas-menu-aplicacao-2026-10-23"></a>
## Incremento — 033/US3, quadros (P2)

**Goal**: FR-012, estado inicial fixo e capacidades explícitas. **Independent Test**: mesmo snapshot/tempo gera mesmo quadro, sem desktop/PPT.

- [ ] T012 [US1] Estender capacidades/fixtures em internal/scenario/model.go e internal/application/playback.go com testes vizinhos após 003/T021/004/T017: quadro usa InitialMenuID existente (raiz ou submenu), ordem/Unicode/enabled e menus vazios; revisita/restart não herdam caminho da prévia livre. Recusar typeText em applicationMenu e conservar recusa em 031/032 por diagnóstico de etapa. Click permanece visual sem dispatch/mudança de menu; nenhuma extensão de item metadata/Kind antes da decisão do coordenador e sincronização de navigation-scenario/wire/capacidades. Roteiro inválido selecionado bloqueia job sem omissão; cena/PPT pertence a 001/T142, não a esta tarefa.

**T012 permanece aberta pela verificação integrada de quadros e job:** código de capacidade/clique visual e fixtures de menu inicial raiz/submenu/vazio, enabled herdado, Unicode/ordem, FrameAt determinístico, A→B→A, revisita/restart e isolamento da prévia livre presentes em scenario/compound_test.go e application/playback_compounds_test.go, **não executadas**. Click destaca o componente sem dispatch de item, mudança de menu, captura ou navegação; typeText é recusado nos três compostos, inclusive vazio, com diagnóstico de etapa. O código não autoriza omitir roteiro inválido; fixtures de recusa na exportação estão em adapter/pptx/composites_test.go, sem execução ou comprovação de job completo. Nenhum novo Kind/metadado de etapa. T011 permanece própria de barcode; integração mista em 001/T138/T143, manual/aceite em T139/T144. Estados históricos e checklists preservados.

**Evidência limitada:** 27 Python e importador `--check` PASS informados pelo principal. Vet/14 suítes Go compile-only PASS e build Windows gerado são **anteriores às últimas alterações** em scenario/click visual e fixtures de quadros, assim como ao grafo de Jornada. O principal repetirá checks/build; não há validação da última revisão. Nenhum teste Go executado neste incremento por política Windows, sem contorno; nenhum PASS Go histórico reaproveitado. Sem novos hashes, homologação nativa/DPI/PowerPoint/SISPL, recuperação de imagens ausentes ou aceite.
