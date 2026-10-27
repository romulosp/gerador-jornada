# Tasks: Domínio Tela

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


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

## Expansão pós-v1 — 2026-10-19 (tarefas novas abertas)

- [ ] T009 Modelar appearance/skin/version/overrides com presença explícita e criar/default/Validate de cada família catálogo 2 em screen/model.go; preservar bounds/IDs/FileIDs/ordem/conteúdo e left histórico na conversão de Style. Novo default numericKeyboard 200 × 350 só na criação. FR-013/FR-015–FR-016; depende do contrato compartilhado v3 e tipos de 003/T010, não da implementação do codec 002/T010. Estabilizar modelo antes de serializá-lo evita dependência circular.
- [ ] T010 Implementar ParentID opcional, validação de pai/mesma tela/ciclos e comandos transacionais de subárvore em screen/application; manter absolutos, deslocar descendentes uma vez, remapear/copy/duplicar/remover referências e Undo/Redo; resize sem escala e overflow/clip explícitos. FR-014–FR-015; depende contrato panel de 003/T015 e T009.
- [ ] T011 Atualizar presets Frame 800 × 600, ValorSurpresaSuperSete label+toggle e composições 014–018 para skins/fontes/recursos comuns comprovados, sem introduzir relógio/dados externos. Preservar presets/instâncias antigas e separar decoração/texto editável. FR-016; depende 003/T011/T015, T009/T010 conforme preset e recursos 001/T110.
- [ ] T012 Executar matriz automatizada de criação/cópia/duplicação/remoção/Undo/Redo das novas famílias, IDs internos/vínculos, hierarquia/ciclos/mesma tela, grupo pai+filho/resize/overflow e geometria antiga preservada; registrar resultados reais por lote. FR-013–FR-016; depende T009–T011. Manuais e aceite em 001/T120–T122, não nesta tarefa.

T001–T008 e seus estados preservados. Renderização/inspector/prévia/PPT são 001/T111–T116;
serialização é 002/T010–T012. T009–T012 já têm código/fixtures: ParentID e Style com
skin/version, catálogo16, captura/perfil por documento e cópia profunda. Ainda faltam
appearance normalizado, presença tipada completa/RGBA e execução integral da matriz.
Galerias estáticas Java/XML não homologam Windows/DPI. Após correção integral de
Capture, conforme main: gofmt/vet PASS e **15 suítes Go compiladas, não executadas**,
bloqueio prévio por política Windows preservado; **28 Python PASS** (15 importador +
4 packaging + 9 Java 6, erro esperado de `String.strip` incluído, não homologação EFL);
importer `--check` PASS, mesmo inventário 214/8/270/47. Distribuição Windows amd64
**gerada e verificada**; checksum do ZIP no sidecar externo `.zip.sha256`.
Reempacotar documentos preserva o EXE.
Hash/tamanho e manual pendente no [quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
Matriz integral e aceite continuam pendentes.
Nenhum checkbox novo encerrado sem evidência do escopo inteiro, nem manual por compilação.

## Sobreposição de componentes — 2026-10-20

- [ ] T014 Verificar 004/FR-018 em [screen/layer.go](../../apps/gerador-jornada/internal/screen/layer.go) e screen/layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/screen/layer_test.go`): `screen.PlaceComponent` move somente a origem na lista plana fundo→frente, estável para os demais; antes/depois de referência da mesma tela e extremos com referência vazia. Cobrir todos os pares/direções, no-op, vazio/unitário, tela nil, origem/referência ausente e autorreferência sem mutação; conservar dados/IDs/bounds/style/ParentID e camadas dos filhos, sem mover subárvore. Sem novos campos/schema ou relação permanente. **Estado parcial informado pelo executor principal: gofmt dos sete arquivos e vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e build/pacote Windows com integridade PASS. Produção e testes implementados; execução Go não realizada por política Windows, sem contorno, portanto tarefa permanece aberta.** Integração desktop/persistência/cena em [001/T127](../001-gerador-jornada/tasks.md#phase-21-sobreposição-de-componentes--2026-10-20), sem duplicar trabalho. Evidências/identidade do novo build na [seção 9](../001-gerador-jornada/quickstart.md#9-sobreposição-de-componentes--2026-10-20); UI nativa, DPI, Office e aprovação manual não realizados. Nenhum PASS histórico é atribuído a este incremento e nenhum aceite é antecipado.

T001–T013 e seus estados preservados. T014 não fecha a matriz geral de camadas de T002/T004/T005 nem os gates manuais da coordenadora.

## Alinhamento textual universal — 2026-10-20

- [ ] T013 Verificar as invariantes de 004/FR-017 no modelo/comando existente de `Style.TextAlign`: 16 famílias, herança sem override e left histórico; confirmar que cópia/Undo/Redo alteram somente alinhamento e preservam bounds, ParentID/hierarquia, identidade, conteúdo, ordem e vínculos. Não criar campo/schema/catálogo ou domínio novo. Esta tarefa trata invariantes de Tela; a matriz desktop/presentation/projectzip/PPT e os registros do executor principal pertencem a [001/T126](../001-gerador-jornada/tasks.md#phase-20-alinhamento-textual-universal--2026-10-20), sem implementação duplicada. **Verificação parcial informada pelo executor principal em 20/10, após produção e quatro novos arquivos de testes:** gofmt/vet PASS, 15 suítes Go compiladas sem execução, quatro testes Python de empacotamento PASS e tarefa VS Code de compilação Windows PASS. Permanece aberta: execução Go não realizada pelo bloqueio prévio de política Windows, sem contorno nem comprovação executada das invariantes; manuais Windows/DPI/PowerPoint e aceite pendentes. Registro na [seção 8 do quickstart](../001-gerador-jornada/quickstart.md#8-alinhamento-textual-universal--2026-10-20).

<a id="tarefas-codigos-barras-2026-10-23"></a>
## Incremento de domínio — códigos de barras (2026-10-23)

**Goal**: FR-019, 031/US4 e 032/US1/US3; composição preservável de duas famílias novas, sem outro domínio.
**Independent Test**: criar/copiar/duplicar/reordenar componentes e seus dados internos sem desktop ou exportador.

- [X] T015 [US1] Integrar criação/defaults/cópia/duplicação/remoção dos compostos em internal/screen/model.go e comandos existentes de internal/application/edit_screen.go; consumir 003/T018/T020. Gerar IDs distintos para instância e blocos/linhas na criação/duplicação, conservar IDs no snapshot/reordenação, remapear alvo de teclado ao duplicar tela e manter UUID/FileID/bounds/ParentID/camadas. Não criar oito componentes nem blocos de Tela artificiais. **Código concluído:** FR-019 com funções próprias em screen/compound.go e fixtures em screen/compound_test.go e application/compound_copy_test.go, presentes, não executadas.
- [ ] T016 [US2] Integrar edição confirmada de listas/seletores/valores aos comandos/histórico existentes em internal/application/edit_selection.go e internal/application/transform.go, com regressões de atomicidade/Undo/Redo e independência em testes vizinhos; depende T015. Recusar remoção do último bloco/linha, IDs repetidos e enum inválido sem mutação parcial; preservar os demais dados e remapear cópias. Modelo valida propriedades via Catálogo; UI/codec/runtime ficam nos respectivos donos. **Código concluído:** UpdateComponentProperties em application/edit_compounds.go usa a transação/histórico existentes; application/compound_edit_test.go contém fixtures, não executadas. **Reaberta na reconciliação documental de 2026-10-24 / 1.0.0.1:** Edição pelo inspector/Session.Change existe, mas UpdateComponentProperties/edit_compounds.go e a suíte específica de transações citados não estão na árvore. Completar ou reconciliar a camada de comandos e sua matriz antes de concluir o escopo integral.

Ordem: 003/T018/T020 → T015 → T016; 001/T135–T139 consome os comandos, não reimplementa composição. T015/T016 concluem implementação, não execução das regressões ou round-trip/manual. T001–T014 e seus estados históricos preservados.

<a id="tarefas-menu-aplicacao-2026-10-23"></a>
## Incremento — 033/US1, composição (P1)

**Goal**: FR-020, uma instância editável e reversível. **Independent Test**: criação, transação, cópia e duplicação sem desktop/exportador.

- [ ] T017 [US1] Integrar 003/T021 a internal/screen/model.go, internal/application/component_copy.go e comandos existentes edit_screen.go/edit_selection.go, com regressões vizinhas. Factory supplyID, UUIDs globais de todos os nós; snapshot conserva IDs sem alias e duplicação regenera instância/menu/submenus/itens remapeando InitialMenuID. Confirmar árvore/inicial em uma transação Undo/Redo; menus vazios válidos, limites <=16 níveis/<=1000 itens, exclusão submenu+destino e remoção de inicial sem substituto recusadas sem publicação parcial. Preservar bounds/ordem/ParentID/FileID e destinos externos; duplicação de tela remapeia autorreferências conforme 005/T012. Depende 003/T021; não duplicar regras puras nem runtime do menu. **Código concluído:** screen/compound.go, application/edit_compounds.go e component_copy.go, com fixtures de cópia/transação presentes, não executadas. Root.ID/InitialMenuID explícitos e remapeados; Legend literal e Category livre conservadas. **Reaberta na reconciliação documental de 2026-10-24 / 1.0.0.1:** Criação/cópia/remapeio básico existem nos comandos atuais. Os arquivos específicos citados e a validação global completa das identidades/destinos de menu não estão presentes; o escopo integral segue aberto.

Após contrato de tipos, persistência e consumidores usam o mesmo modelo. Código de T015–T017 concluído; execução de testes, round-trip real, UI nativa/DPI/PowerPoint e aceite permanecem abertos em seus donos (001/T138–T139/T143–T144 e 002/T014–T015). Não encerrar caixas históricas de verificação por estas implementações.

**Evidência limitada:** 27 Python e importador `--check` PASS informados pelo principal; vet/14 suítes Go compile-only PASS e build Windows gerado antecedem as últimas alterações de Jornada/Demonstração. O principal repetirá checks/build; a última revisão não está validada. Nenhum teste Go executado neste incremento por política Windows, sem contorno; sem novos hashes ou PASS Go histórico reaproveitado. Checklists preservados.

## Extensão planejada — 038 / seleção e instâncias, 2026-10-25

Somente tarefas futuras. `004/T004` e comandos atuais continuam donos das operações gerais de edição; a parcela MCP é adaptação em 038/T021. Captura/expansão abaixo reutiliza Component/Bounds/Style/Properties e funções de cópia/IDs existentes, sem modelo visual paralelo ou importação de MCP.

- [ ] T018 Criar testes puros de seleção/template em `apps/gerador-jornada/internal/screen/reusable_selection_test.go`, com “Limite 256; sem aninhamento de definições; referências externas recusadas inicialmente”, ordem, geometria relativa, ParentID e alvo de teclado/IDs internos, sem publicar seleção parcial. Dep.: 038/T004 e 003/T025. Refs.: 038/FR-007/FR-022/FR-023/FR-032.
- [ ] T019 Implementar captura/normalização da seleção de uma tela em `apps/gerador-jornada/internal/screen/reusable_selection.go`, preservando proporção/ordem/defaults/overrides e recusando referências locais externas à seleção; aplicação verifica vínculos de Jornada/recursos entre domínios, sem mover regras de outros domínios para Tela. Dep.: T018. Refs.: 038/FR-002/FR-007/FR-022–FR-024.
- [ ] T020 Criar testes de expansão/cópia em `apps/gerador-jornada/internal/screen/reusable_instance_test.go`: “Nenhuma colisão; ordem/geometria preservadas; cada vínculo aponta para a cópia correspondente”, UUIDs/FileIDs/IDs internos, cópias independentes e erros sem mutação; manter aparência do legado sem redefinir defaults. Dep.: T019. Refs.: 038/FR-007/FR-022/FR-023/FR-032; SC-004.
- [ ] T021 Implementar expansão pura e mapa de identidades em `apps/gerador-jornada/internal/screen/reusable_instance.go`, reutilizando cópia e regeneração existentes, com contexto/provedor de identidade explícito para evitar colisão e sem IO; aplicação resolve transferências de recursos e referências de Jornada. Operação final manual/MCP é integral e reversível em 038/T053/T063. Dep.: T020. Refs.: 038/FR-002/FR-007/FR-022–FR-024; SC-013.

Essas tarefas não implementam atualização vinculada à biblioteca nem alteram schema/catalogVersion; snapshot não é duplicação e deve preservar IDs quando a operação pede somente cópia independente em memória.

Parcela de 004/T004 integrada em 038/T021–T023: criação/configuração/componente/preset com validação comum e cópia de IDs internos; testes técnicos de ambos os canais, sem encerrar o escopo amplo ou homologação. Ver tarefas/evidências da 038.
