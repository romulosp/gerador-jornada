# Tasks: Domínio Catálogo

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md) e contrato compartilhado.
**Prerequisites**: preparação/toolchain da feature coordenadora; regras/documentação revisadas.
**Tests**: checks executados em 2026-10-18, conforme evidências abaixo. Caixas abertas distinguem cobertura/implementação parcial de trabalho não iniciado; não representam ausência de código.
**Status**: reconciliação documental concluída conforme constituição 1.2.0; domínio com implementação e lacunas identificadas, sem aceite integral.

**Instrução vigente — 2026-10-21:** [T017 aberta](#t017-builtin-filtrado) verifica [FR-012](spec.md#builtin-filtrado-fr-012) sob [001/T132](../001-gerador-jornada/tasks.md#t132-builtin-filtrado), sem duplicar implementação. Inventários/checks de 214 definições/47 diagnósticos são históricos **antes da redução**; estado atual e roteiro na [seção 14 da coordenadora](../001-gerador-jornada/quickstart.md#14-builtin-filtrado). Caixas anteriores e checklist 13/2 preservados.

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

## Expansão pós-v1 — todo o acervo Java/EFL (2026-10-19)

Todas as novas caixas abertas; FR-007–FR-011; referência da cobertura C01–C26 em [research.md](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19). Checklist 13/2 preservado, prosseguimento já autorizado.

- [ ] T009 Classificar saída integral de 001/T108 por família/variante/preset/helper/estado, completar defaults geométricos e contratos de cada classe/skin; registrar proveniência e bloqueios inclusive C04 ValorSurpresaSuperSete, Password/Dialog e delegates ausentes. Depende do inventário T108 da coordenadora, não da biblioteca gerada T110; não duplicar parser. FR-007/FR-009–FR-011.
- [ ] T010 Evoluir Kind/Definition/Properties/Copy/Validate para catálogo 2 com união e metadados tipados de variantes/estados/capacidades; contrato de referência skin/face/overrides sem Fyne/PPT/XML/I/O no domínio. Preservar discriminadores atuais; integrar por lote com 002/T010 e 004/T009, sem habilitar tipos incompletos. FR-008–FR-010; depende T009.
- [ ] T011 Implementar regras/defaults/cópia/fixtures de Button e twoLabel, Label/Title/todas skins aplicáveis, ToggleButton/grupos e Radio simples; classificar preset label+toggle de ValorSurpresaSuperSete. Preservar IDs e semântica de seleção; 007–009/FR-008, 011/FR-008. Depende T010; consumidores em 001/T111–T116, preset em 004/T011.
- [ ] T012 Implementar regras/defaults/cópia/fixtures de TextField e todos documentos/máscaras inventariados, TextArea, PasswordField, ComboBox/UF e Table/headers/colunas/células. Provar limites de máscara e seleção local; nenhuma validação fiscal inventada ou credencial real. 010/FR-008; depende T010, C06–C10/C26. Recursos/evidência insuficientes mantêm homologação bloqueada.
- [ ] T013 Implementar regras/defaults/cópia/fixtures de NumericKeyboard, AlphaKeyboard e calculadora como variante/helper comprovado; contratos de foco/cursor/modificadores/alvos compatíveis da mesma tela. Novo numérico 200 × 350, 14 teclas, sem resize de existentes. Distinguir legenda LIMPAR de VK_BACK_SPACE Java. Depende T012/contrato de campo; 013/T001–T003 detalha fixture/aceitação sem repetir modelo.
- [ ] T014 Implementar regras/defaults/cópia/fixtures de todas variantes Report (Report/2/3/Event/Event2/Notificacao/SelecaoMultipla/Internal), linhas tipadas/seleção/altura e ProgressBar; partir do piloto Event sem assumir que homologa subclasses. Depende T010 e partes Panel de T015 necessárias; 012/T001–T003 detalha piloto. FR-008/FR-011.
- [ ] T015 Implementar regras/defaults/cópia/fixtures de Panel/Panel2, capacidades scroll/clip, Dialog e contratos de presets Frame/014–018; não criar Kind por seta/skin nem recriar janela SO. Fornecer primeiro o subconjunto de dados Panel necessário a T014, depois contêineres completos. Depende T010; ParentID/expansão de presets ficam em 004/T010–T011; referência ausente bloqueia homologação.
- [ ] T016 Executar matriz automatizada de defaults/enums/propriedades/cópia/capacidades de cada família/variante/estado e casos inválidos, incluindo zero/falso/ausência e IDs/seleções; atualizar rastreabilidade granular com revisão/resultados reais. Depende T011–T015 por lote. Integração de codecs/UI/PPT em 001/T118, homologação/aceite em T120–T122 da coordenadora, não nesta caixa.

T001–T008 e estados antigos preservados. T010 é fundação; T011–T015 podem evoluir por sublotes sem ciclos de dependência, com contrato de parte Panel pronto antes do piloto. Fontes/importador/skin (001/T108–T110), persistência (002/T009–T012), composição (004/T009–T012) e consumidores (001/T111–T116) têm donos distintos; não duplicar implementação aqui.

**Parcial implementado na continuação:** catálogo2/16 famílias e galerias estáticas
Java/XML; Profile/Library por documento, revisão arquivada exata e captura integral
dos metadados da biblioteca original, inclusive variantes/revisões não usadas e todos
os estados/dependências, sem novos binários. Não são screenshots
Java/homologação DPI. Style/Fonts/Aliases/Resources ainda não atendem todo o alvo de
appearance normalizado/presença/RGBA/descritores. Calculadora sem cálculos e referência
incompleta impedem encerrar lotes integralmente. Após correção integral de Capture,
conforme main: gofmt/vet PASS, **15 suítes Go compiladas, não executadas**, bloqueio
prévio por política Windows preservado. Python **28 PASS** (15 importador + 4 packaging
+ 9 Java 6, incluindo erro esperado de `String.strip`, não homologação EFL); importer
`--check` PASS, mesmo inventário 214/8/270/47. Distribuição Windows amd64 **gerada e
verificada**; checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos
preserva o EXE. Identidade do EXE
e manual pendente no [quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
T016 continua aberto sem execução integral. Nenhuma nova caixa concluída por compilação
ou documentação.

<a id="t017-builtin-filtrado"></a>
## Verificação da disponibilidade do builtin filtrado — 2026-10-21

- [ ] T017 [US1] Verificar localmente [FR-012](spec.md#builtin-filtrado-fr-012), clarificação de FR-007/FR-010, consumindo a geração e integração de [001/FR-058/T132](../001-gerador-jornada/tasks.md#t132-builtin-filtrado), sem reimplementar filtro, parser, codecs ou consumidores. Conferir disponibilidade por família/identidade tipada no builtin para novos projetos: **206 definições ativas / 8 fontes / 270 recursos / 31 diagnósticos**, **8 definições adiadas / 16 referências missing**, **214 originais** preservadas; não exigir 206 opções em cada seletor das 16 famílias. As oito variantes de [001/FR-058](../001-gerador-jornada/spec.md#builtin-filtrado-fr-058) não aparecem para nova inserção/troca nessa revisão, mas constam integralmente de `deferred_definitions` com `type`/`name`/`definition`/`reasons`, sem apagar `source_tree`, `xml_sha256`, origens/pins/bytes/fontes. Verificar retirada da definição inteira por referência `imagePath`/`slider`/`font:type` ausente/inválida, caminho recusado ou tipo incompatível; filho afetado adia composição pai/todos os filhos, sem órfãos nem exclusão de homônimos independentes.

	**Verificação mista — ABERTA:** **19 Python do importador PASS**, incluindo descoberta/execução explícitas dos quatro novos métodos após corrigir o aninhamento; **4 Python de empacotamento e importador `--check` PASS** para a skin **`legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a`**. Regressões Go de biblioteca/seletores/Capture/perfis próprios e arquivados implementadas; gofmt/vet e **compile-only de 15 suítes PASS, ZERO execuções Go**, respeitando política Windows. Build/pacote e integridade somente leitura PASS; dois EXEs idênticos, identidades reais e limites na seção 14. Compilação não executa as regressões nem comprova novo round-trip de disco; execução Go e roteiros manuais continuam pendentes.

	**Aceitação local pendente:** novos projetos usam builtin filtrado completo, não subset por uso; preservar inserções válidas após salvar/reabrir e variantes/estados/referências de projetos anteriores e perfis próprios/arquivados, sem migração adicional, poda, upgrade, enriquecimento ou fallback builtin. Reimportação somente em futura change com recursos exatos, origem/hash/autorização, regeneração e validação. Manifest/Profile/project/catalog/produto **1/1/3/2/1.0.0.0** e gates de fontes/licenças, `DELEGATE_MISSING`, `FIDELITY_UNVERIFIED`, Java/PPT inalterados. Roteiro na [seção 14](../001-gerador-jornada/quickstart.md#14-builtin-filtrado); manuais/DPI/Office/aceite continuam em 001/T120–T122, sem promessa de 100% de fidelidade.

T001–T016, estados e evidências preservados; T017 não conclui a matriz ampla T016 nem duplica T132: trata disponibilidade/invariantes locais, enquanto 001 coordena integração e evidência dos adaptadores. Checklist 13/2 e aprovações permanecem inalterados.

<a id="tarefas-codigos-barras-2026-10-23"></a>
## Incremento de domínio — códigos de barras (2026-10-23)

**Status atualizado — 2026-10-23**: o marco original era planejamento; tipos, defaults, cópia e regras dos três compostos estão implementados no catálogo 3. As fixtures Go existem, mas não foram executadas. 031/032 detalham contratos; o modelo continua neste domínio.

### Fundação — tipos e capacidades

- [X] T018 Definir tipos `barcodeInput`/`barcodeProcessed`, propriedades/defaults/capacidades e cópia profunda em internal/catalog/model.go e internal/catalog/copy.go; fixtures em internal/catalog/model_test.go. FR-013/FR-015; contrato compartilhado de 2026-10-23; entrada vazia e processado com um bloco/CNPJ/Nome vazios. **Código concluído:** catálogo 3, defaults 610 × 60/530 × 120 e consumidores implementados, sem mudar instâncias/perfis antigos; fixtures presentes, não executadas. T019/T020 usam arquivos próprios.

### US1 — entrada validada (031/US1–US3)

- [ ] T019 [US1] Implementar regras puras e fixtures em internal/catalog/barcode.go e internal/catalog/barcode_test.go (novos arquivos planejados): DV local Módulo 10, geral Módulo 11, conversão exata 44/47, zeros, valor em centavos/banco/fator/CAIXA e SIGCB experimental. FR-014; depende T018. Cobrir F01/F02, caracteres inválidos, 43/45/46/48 posições, DVs adulterados, candidato 10/11 e ausência de resultado parcial; comparar com classes reais quando recebidas antes de afirmar equivalência SISPL. Sem UI/IO, data por relógio ou cálculo de pagamento.

### US2 — parametrização processada (032/US1–US2)

- [X] T020 [P] [US2] Implementar estrutura/defaults/validação/cópia e fixtures em internal/catalog/barcode_processed.go e internal/catalog/barcode_processed_test.go, após T018. FR-013/FR-015; “ao menos um bloco e uma linha por bloco”; UUIDs únicos, tipos `cpf/cnpj/name/tradeName/corporateName`, títulos/valores Unicode/vazios, opções repetidas permitidas. Trocar rótulo conserva valor literal; cópia não compartilha slices. Nenhum CPF/CNPJ/nome é extraído de boleto ou validado fiscalmente. **Código concluído; fixtures presentes, não executadas.** Pode evoluir em paralelo a T019, sem alterar o mesmo arquivo de fundação.

**Ordem/integração**: T018 → T019/T020 → 004/T015–T016 e 002/T013–T014 → consumidores 001/T135–T138. Regras de domínio não são reimplementadas na coordenadora. Execução dos testes e homologação permanecem pendentes, sem herdar PASS histórico.

<a id="tarefas-menu-aplicacao-2026-10-23"></a>
## Incremento — 033/US1, árvore de menu (P1)

<a id="t022-correcao-menu-xml"></a>
### Correção após teste parcial — 2026-10-23

- [X] T022 [US1] Recuperar e preservar `material-legado/menu/menu.xml`; substituir a transcrição parcial e relatórios sintéticos por `menu_legacy_generated.go`, gerado por `scripts/import_legacy_menu.py`. Factory preserva árvore ativa, ordem, espaços e Movement*/Shortcut, com novos UUIDs; dez jogos incluem MILIONÁRIA e excluem LotoGol comentado. BOLÃO/digital e relatórios reais importados integralmente. `menu_xml_test.go` compara cada nó com o XML independente; novas instâncias apenas. Skins/consumidores em 001/T152 e evidências na seção 24 do quickstart; homologação não é concluída por esta tarefa.

Caixa encerra somente código/fixtures, não checks/build/aceite. Aprovação parcial do usuário com correções; estados históricos/checklists preservados. Hashes/resultados anteriores não validam esta revisão.

**Goal**: FR-016; modelo puro reutilizado por todos os consumidores. **Independent Test**: criar/copiar/validar seed e árvore editada sem UI/IO. Fundação 003/T018 primeiro; nenhum novo setup de projeto.

- [X] T021 [US1] Acrescentar ApplicationMenu="applicationMenu" e Properties.Menu ao registro central em internal/catalog/model.go/copy.go, com modelos em internal/catalog/menu.go, seed em menu_defaults.go/menu_legacy_generated.go e testes vizinhos. Implementar factory supplyID com Enabled=true e seed serviços/jogos/aninhamento/Bolão da fonte fornecida. Validar layouts root/standard/financial/bolao, inicial existente, UUIDs únicos, “profundidade <=16 (raiz=1)” e “total de MenuItem <=1000”, menus vazios válidos, ciclos e exclusão submenu+destino. Cópia profunda sem alias; metadados inertes. Fixtures 16/17, 1000/1001, duplicados, Unicode/vazios e cópia independente presentes, não executadas. **Código concluído:** default 800 × 600, inicial por ID explícito, `Legend` literal independente de MovementName e `Category` livre/Unicode, inclusive vazia. Depende T018; não duplica matemática ou ficha. Seed transcrito da fonte não comprova fidelidade visual sem imagens/homologação.

**Dependências/paralelismo**: após registro serial de T018/T021, arquivos próprios de menu/barcode e consumidores não duplicam o modelo. T018/T020/T021 concluem somente implementação. **T019 permanece aberta pelo escopo de comprovação:** regras puras Módulo 10/11, conversão 44/47, centavos/banco/fator/CAIXA e SIGCB experimental e fixtures F01/F02/erros/limites estão em barcode.go/barcode_test.go, mas não foram executadas nem comparadas às classes reais SISPL/SIGCB ausentes. Não há cálculo de pagamento ou homologação bancária. Integração em 001/T138/T143; manual em T139/T144. Estados históricos e checklists preservados.

**Evidência limitada:** 27 Python e importador `--check` PASS informados pelo principal. Vet e 14 suítes Go compile-only PASS e o build Windows gerado são anteriores às últimas alterações de Jornada/Demonstração; o principal repetirá checks/build. Nenhum teste Go executado neste incremento por política Windows, sem contorno; sem novos hashes nem validação da última revisão. Faltas de imagens, DPI, PowerPoint e aceite permanecem pendentes, sem reaproveitar PASS Go histórico.

### Evidência parcial do recorte visual — 2026-10-24

A revisão coordenada 001/T148 implementa o default numérico 200 × 350 para novas inserções, sem mudar bounds persistidos; aplica defaults visuais por família na apresentação. `internal/skin` oferece as definições disponíveis e nomes reais de fonte, com oito definições ausentes fora da oferta. Isso não fecha T013 nem as demais tarefas da expansão: seus tipos/contratos mais amplos permanecem pendentes. Evidência: seção 22 do quickstart integrado.

- [X] T023 Preservar ausência literal de tituloPrincipal e variantes financeiras/digitais/Bolão na fábrica conforme PnlMenu.java; fortalecer comparação de Title com o XML independente. Nenhum parser runtime, campo persistido ou migração automática. Metadados continuam inertes e cadastro movimento.xml permanece pendente em 001/T154; consumidores e evidência em 001/T153/quickstart seção 25.
