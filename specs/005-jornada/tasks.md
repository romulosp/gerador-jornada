# Tasks: Domínio Jornada

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


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

## Expansão pós-v1 — 2026-10-19 (tarefas novas abertas)

- [ ] T009 Atualizar resolução/validação de referências e capacidades de origem/alvo em journey/application para catálogo 2; button twoLabel conserva none/navigate, novas seleções/teclados não ganham ação externa implícita; diagnosticar alvos incompatíveis/fora da tela mantendo bindings/UUIDs. FR-006–FR-007; depende contratos de 003/T010–T015, sem implementar estado de prévia de 001/T114.
- [ ] T010 Executar regressões A→B→A, twoLabel, tipos sem capacidade navigate, referências removidas, snapshot/reset de estados iniciais e migração de aparência sem alteração de rotas; confirmar que Prévia Tela Atual não navega. FR-006–FR-007; depende T009 e 001/T114 para integração runtime, sem duplicar seus testes de widget. Registrar revisão/resultados; homologação/aprovação seguem 001/T120–T122.

Nenhuma nova conclusão técnica/manual atribuída nesta revisão; T001–T008 preservados. Java segue pausado; não exigir serviços/periféricos ou `.page`.

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
## Incremento de verificação — códigos de barras (2026-10-23)

- [ ] T011 [US1] Cobrir FR-008 em testes de internal/journey e do índice de referências de internal/application/document.go: `barcodeInput`/`barcodeProcessed` não originam navigate, DV válido/erro/reinício e troca de rótulo não alteram bindings ou inicial; botão explícito conserva rota. Após 001/T135, exercitar A→B→A com entrada vazia e configuração processada intacta e Prévia Tela Atual sem navegação. Depende 003/T018 para capacidades, sem mover matemática ou sessão para Journey. Contrato em navigation-scenario.md, 031/FR-007/FR-012 e 032/FR-009; não criar motor condicional ou segundo estado de runtime.

**T011 permanece aberta como verificação:** há código de capacidades/projeção em application/journey_compounds.go e fixtures em journey/journey_compounds_test.go, application/journey_compounds_test.go e testes de prévia, sem execução Go. Incluem ausência de navigate nos barcodes, preservação de bindings, retorno/reset e escopo de Tela Atual; não comprovam execução A→B→A ou round-trip real. Runtime continua 001/T135/T138, sem duplicação neste domínio. T001–T010 e suas caixas históricas preservadas.

<a id="tarefas-menu-aplicacao-2026-10-23"></a>
## Incremento — 033/US2, destinos explícitos (P1)

**Goal**: FR-009, referência de folha validada, sem rota inferida. **Independent Test**: grafo/índice com folhas habilitadas/desabilitadas sem desktop.

- [ ] T012 [US1] Estender internal/journey e projeção/índice em internal/application/document.go, com testes vizinhos, para validar destino explícito de folha por screen/component/item IDs. Depende contrato 003/T021/004/T017. Incluir destinos de ramos desabilitados na integridade/impacto de remoção; alcançabilidade considera habilitação efetiva. Recusar destino inexistente/submenu+destino, conservar ciclos/autorrotas explícitos, manter 031/032 sem navigate e folha sem destino sem rota. Duplicação de componente conserva destinos; duplicação de tela remapeia só autorreferências com 004/T017. Não duplicar Action/TargetScreenID nem importar árvore/UI no domínio. Consumidor 001/T140 verifica Prévia Tela Atual sem transição; metadados de movimento/atalho nunca são lookup ou dispatch. **Código concluído:** Document.NavigationGraph() em application/journey_compounds.go e Graph em journey/graph.go, com Validate/References/Reachable, snapshot sem alias e habilitação por item/ancestrais/ParentID. Fixtures de integridade, impacto, ciclos, limites e duplicação presentes, não executadas. **Reaberta na reconciliação documental de 2026-10-24 / 1.0.0.1:** MenuActivate navega para folhas autoradas, mas NavigationGraph, journey/graph.go e o índice global de referências/alcançabilidade citados não existem. O runtime parcial não conclui a integridade requerida.

T012 conclui implementação, não prova executada do grafo/runtime. **T006 histórica continua aberta:** Reachable já existe como consulta pura, mas isso não implementa o warning de inalcançável descrito naquela caixa nem a localização clicável de 001/T030. T011 permanece verificação exclusiva de barcode. Integração mista em 001/T138/T143, manual/aceite em T139/T144; estados históricos e checklists preservados.

**Evidência limitada:** 27 Python e importador `--check` PASS informados pelo principal. Vet e 14 suítes Go compile-only PASS e o build Windows gerado são **anteriores às últimas alterações**, inclusive NavigationGraph/References/Reachable e suas fixtures; o principal repetirá checks/build. Não afirmar compilação ou validação final dessas adições. Nenhum teste Go executado neste incremento por política Windows, sem contorno; sem novos hashes, homologação nativa/DPI/PPT/SISPL ou reaproveitamento de PASS Go histórico.
