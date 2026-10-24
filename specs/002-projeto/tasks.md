# Tasks: Domínio Projeto

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


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

## Expansão pós-v1 — 2026-10-19 (tarefas novas abertas)

- [ ] T009 Modelar perfil/descritores de fontes/recursos tipados em application.Document/internal/skin, mantendo metadados/assets em project/model.go; testar referências compartilhadas/imutabilidade. FR-008/FR-010. Alvo normativo SUBSET de bytes usados e permitidos preservado; desvio arquitetural atual conserva integralmente metadados da biblioteca original para editabilidade futura, não todos os binários. Depende contrato de 001/T108–T110, não reimplementa importador. Parcial: AppearanceProfile *skin.Profile, Copy, Library e Capture implementados; descritores normalizados e portabilidade licenciada integral ainda pendentes.
- [ ] T010 Evoluir projectzip/wire.go para schema 3/catalog 2 com reader v1/v2/v3 estrito, aparência/overrides/ParentID e todas Properties do catálogo novo; migrar sem perda de IDs/FileIDs/conteúdo/bounds/ordem/vínculos/overrides, fixando left histórico. Novos defaults não redimensionam instâncias antigas. FR-009; depende T009 e contratos 003/T010, 004/T009–T010; por lote.
- [ ] T011 Evoluir projectzip/package.go e fluxo de salvamento para TTF/perfil/recursos usados com limites/MIME/hash/safePath, aviso de incompatibilidade e cópia íntegra original antes do primeiro save migrado/publicação atômica; erro preserva sessão/arquivo. FR-008–FR-010; depende T009/T010. I/O permanece em adaptadores/aplicação.
- [ ] T012 Executar fixtures v1/v2→v3 e round-trip de cada família/variante/perfil, igualdade de identidades/bounds/referências/overrides, incluindo família legada explícita conhecida/desconhecida, descriptor unresolved e fontes usadas só por override; testar corrupção/hash/licença/fontes ausentes, ZIP/backup/falha de publicação/snapshot. Registrar revisão e resultado observado; FR-008–FR-010. Depende T010/T011; teste automatizado não fecha 001/T120–T122.

IDs T001–T008 e evidências históricas preservados. Caixas novas continuam abertas
pelo escopo integral, não por ausência total de código. T010/T011 já têm schema3,
reader1/2/3, perfil opcional transitório v3 (proibido v1/v2), captura/validação e
backup; ainda não writer `component.appearance` normalizado/overrides/RGBA ou todos
os descritores-alvo. `Fonts/Aliases/Resources` e `Style` são transitórios.

T009/T012 têm fixtures de perfil/arquivos/cópia/hash/bytes, compiladas, **não executadas**.
Perfil schema1 admite até 16 revisões diretas; hash inclui arquivos/hashes, exclui só
o próprio. Capture preserva todos os metadados originais em cópias privadas:
definições/fontes/aliases/recursos e todas as revisões arquivadas, inclusive não usadas
e todos os estados. Não enriquece perfis privados com builtin nem faz upgrade;
uso só limita diagnósticos e acréscimo de entrada missing. Conserva metadados
referenced/unknown, sem acrescentar binários sem autorização nem com JSON allowed.
Desvio justificado pelo bug crítico save/reopen → AddComponent `defaultconfirm`;
inserir famílias/trocar variantes deve continuar possível após reopen. Payload
selado mais amplo altera o hash, não autoriza incorporação binária integral.

T012 inclui regressões novas save/reopen/addButton/variant/transient states/custom
namespace: conservar definições de todos os estados e revisões não usadas, inserir
famílias/trocar variantes após reopen, sem builtin/upgrade de perfil privado e sem
novos bytes. Incluídas na compilação pelo main **após a correção integral de Capture**;
não executadas nem declaradas PASS de testes Go. Não fechar T012 pelo escopo restante.
ResourceData exige AssetID/hash/tamanho/MIME ou builtin exato; revisão desconhecida
não recebe fallback. Evidências do main: gofmt/vet PASS e **15 suítes Go compiladas,
não executadas**, bloqueio prévio por política Windows preservado. Python **28 PASS**
(15 importador + 4 packaging + 9 Java 6, erro esperado de `String.strip` incluído,
não homologação EFL); importer `--check` PASS, mesmo inventário 214/8/270/47.
Distribuição Windows amd64 **gerada e verificada**; checksum do ZIP no sidecar
externo `.zip.sha256`. Reempacotar documentos preserva o EXE. Identidade do EXE e manual pendente no
[quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
Homologação/aprovação continuam 001/T120–T122;
integração 001/T118 não duplica persistência. Nenhum checkbox encerrado por compilação.

<a id="tarefas-codigos-barras-2026-10-23"></a>
## Incremento de domínio — códigos de barras (2026-10-23)

**Goal**: preservar os dois compostos sem gravar runtime; FR-011, 031/US4 e 032/US3. O marco inicial foi documental; agora há codec dos três compostos implementado, sem comprovação executada de round-trip.
**Independent Test**: codec/ZIP com duas configurações e formatos anteriores, sem UI ou PPT.

- [ ] T013 [US1] Evoluir internal/adapter/projectzip/wire.go e fluxo existente de salvamento para schema3/catalog3 e Properties `barcodeInput`/`barcodeProcessed`, preservando reader 1/1, 2/1 e 3/2, identidade/ordem/Unicode/perfis exatos. Depende 003/T018/T020 e 004/T015. Não persistir valores capturados, foco/cache/validação/modal; sinalizar primeiro save incompatível e preservar cópia íntegra original antes da publicação mesmo se só catalogVersion mudar. Não alterar allowlist de recursos ou duplicar regras de Catálogo no DTO. **Código concluído:** writer/reader 3/3 e backup também em 3/2→3/3; validação tipada estrita em compound_wire.go. Fixtures presentes, não executadas; comprovação de disco/falhas continua T014. **Reaberta na reconciliação documental de 2026-10-24 / 1.0.0.1:** Writer 3/3, propriedades compostas e backup existem; preservação de perfis visuais exatos citada nesta tarefa não existe no modelo atual. A compatibilidade completa permanece pendente.
- [ ] T014 [US1] Cobrir FR-011 em internal/adapter/projectzip/package_test.go e testes vizinhos de wire/migração: round-trip das duas famílias, bloco/linha/seleção/valor vazio e acentuado, UUIDs/ordem, enum/campo desconhecido, recusa do catálogo novo pelo reader antigo e preservação da sessão/original em falhas de backup/publicação. Depende T013. Comprovar ausência de runtime no manifesto e preservação de perfis próprios/arquivados sem builtin/upgrade. Registrar resultado observado; compilação não é execução de testes.

Integração única em 001/T138, manual/aceite em 001/T139. Não reimplementar modelos 003, instâncias 004 ou cena/PPT. T001–T012 e seus estados permanecem intactos.

<a id="tarefas-menu-aplicacao-2026-10-23"></a>
## Incremento — 033/US1 e US3, persistência (P1)

**Goal**: FR-012, árvore/inicial confirmados persistíveis. **Independent Test**: round-trip ZIP de duas árvores, sem UI.

- [ ] T015 [US1] Estender internal/adapter/projectzip/wire.go e testes package_test.go/vizinhos para applicationMenu/root/initialMenuId e submenus/itens do contrato. Depende 003/T021, 004/T017 e matriz/backup 002/T013–T014; reutilizar schema3/catalog3 para 031/032/033. Preservar UUIDs, ordem, Unicode/vazios, enabled=false, metadados/destinos e menus vazios items=[]; recusar ausências obrigatórias, campos/layouts desconhecidos, UUIDs duplicados, >16 níveis, >1000 itens e submenu+destino antes de publicar. Comprovar ausência de runtime, reader histórico 1/1 2/1 3/2 e novo 3/3, falha de backup sem overwrite e preservação exata de perfis. Não criar segundo backup/migrador nem gravar código executável.

**T015 permanece aberta pelo critério misto:** codec de menu e fixtures em internal/adapter/projectzip/compound_test.go e compound_migration_test.go presentes, incluindo `Menu.Root.ID` e `InitialMenuID` explícitos (wire `root.id`/`initialMenuId`), `Legend` literal e `Category` livre/Unicode. Cópia de IDs e backup 3/2→3/3 implementados; os testes de round-trip, compatibilidade, ausência de runtime e falhas de backup/publicação **não foram executados**, nem há novo round-trip real comprovado. T014 também permanece aberta: código/fixtures não equivalem a resultado observado. 001/T138/T143 integram e T139/T144 registram manual/aceite. T001–T012 e checklists históricos preservados.

**Limite das evidências atuais:** 27 Python e importador `--check` PASS informados pelo principal; vet e 14 suítes Go compile-only PASS, assim como o build Windows gerado, antecedem as últimas alterações em Jornada/Demonstração e serão repetidos/recompilados pelo principal. Nenhum teste Go executado neste incremento por política Windows; nenhum PASS histórico valida a última revisão. Sem novos hashes, homologação nativa/DPI/PPT/SISPL ou aceite nesta atualização documental.
