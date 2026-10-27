# Tasks: Tela inicial (launcher) e Estrutura de Projeto

**Input**: Design documents from `/specs/034-launcher-estrutura-projeto/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Incluídas onde a constituição exige (invariantes, persistência, navegação e contratos). Homologação manual registrada em quickstart.md.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Baseline verification before changes.

- [X] T001 Verificar baseline: `go build ./...` e `go test ./...` passam em `apps/gerador-jornada` antes de qualquer alteração; registrar a revisão inicial no log.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Infraestrutura compartilhada que TODAS as histórias usam (recentes e workspace).

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Implementar armazenamento de projetos recentes (carregar/salvar JSON por usuário com caminho do workspace, nome de exibição e timestamp da última abertura; entrada indisponível é indicada sem quebrar as demais) em `apps/gerador-jornada/internal/adapter/desktop/recents.go`.
- [X] T003 [P] Implementar resolução do workspace (pasta do sistema em uso a partir do `FilePath` atual ou de pasta escolhida) e listagem de arquivos `.jornada` na raiz da pasta (não recursiva, nome relativo) em `apps/gerador-jornada/internal/adapter/desktop/workspace.go`.

**Checkpoint**: Fundação pronta — as histórias podem começar.

---

## Phase 3: User Story 1 - Tela inicial (launcher) (Priority: P1) 🎯 MVP

**Goal**: Ao abrir sem documento, mostrar projetos recentes e criar projeto em branco (fundo atual) ou com exemplo (cenário de carregamento atual).

**Independent Test**: Abrir sem `-open` → ver launcher; criar em branco e com exemplo; reabrir e ver recentes (entrada movida/excluída indicada como indisponível).

### Implementation for User Story 1

- [X] T004 [US1] Construir a tela inicial (lista de recentes + ações "Novo projeto em branco" e "Novo projeto com exemplo") em `apps/gerador-jornada/internal/adapter/desktop/launcher.go`, usando o armazenamento de T002.
- [X] T005 [US1] Ligar o launcher em `apps/gerador-jornada/cmd/gerador-jornada/main.go`: sem `-open`/documento, exibir a tela inicial em vez de abrir `application.Example()` automaticamente; rotear em branco (`application.NewDocument(nome)` + fundo padrão) e exemplo (`application.Example()` + fundo padrão).
- [X] T006 [US1] Registrar projeto recente ao abrir/salvar em `apps/gerador-jornada/internal/adapter/desktop/files.go` (em `openProject`, `saveProject` e ao criar pelo launcher), preservando caminho/nome/timestamp.
- [X] T007 [US1] Adicionar teste para recentes (persistência/reabertura, entrada indisponível) e para a distinção em branco/exemplo em `apps/gerador-jornada/internal/adapter/desktop/launcher_test.go` (usar fixtures, sem rede).

**Checkpoint**: US1 funcional e testável sozinha.

---

## Phase 4: User Story 2 - Estrutura projeto (Priority: P1)

**Goal**: Aba "Estrutura projeto" no painel esquerdo listando em árvore as `.jornada` do workspace, com criar/importar.

**Independent Test**: Com uma pasta de workspace, abrir a aba e ver a árvore; criar e importar uma jornada; arquivo inválido é recusado sem alterar origem.

### Implementation for User Story 2

- [X] T008 [US2] Adicionar a aba "Estrutura projeto" ao `AppTabs` do painel esquerdo em `apps/gerador-jornada/internal/adapter/desktop/tools.go`, renderizando árvore com raiz = nome da pasta e descendentes = arquivos `.jornada` (usar a listagem de T003).
- [X] T009 [US2] Adicionar o menu de três pontos com "Criar jornada" e "Importar jornada existente" na aba, em `apps/gerador-jornada/internal/adapter/desktop/workspace.go` (controles em PT).
- [X] T010 [US2] Implementar "Criar jornada": gravar um `.jornada` em branco (fundo padrão) na pasta do workspace com nome escolhido e normalizado para arquivo (nome de exibição = `Project.Name`), e refletir na árvore.
- [X] T011 [US2] Implementar "Importar jornada existente": copiar um `.jornada` válido para dentro do workspace sem alterar a origem; recusar arquivo corrompido/versão incompatível com diagnóstico localizado (reusar `projectzip.Decode` para validar antes de copiar).
- [X] T012 [US2] Adicionar teste para listagem/árvore, criação e importação (origem intacta, inválido recusado) em `apps/gerador-jornada/internal/adapter/desktop/workspace_test.go`.

**Checkpoint**: US2 funcional; US1 + US2 funcionam juntas.

---

## Phase 5: User Story 3 - Criação rápida e navegação entre jornadas (Priority: P1)

**Goal**: Folha do Menu aplicação sem destino, ao ser clicada na prévia, sugere criar uma jornada com o nome da folha; a jornada entra na árvore e a prévia navega para ela (substituir prévia + "Voltar").

**Independent Test**: Clicar folha sem destino → sugestão; aceitar → jornada na árvore; clicar de novo → navega; cancelar não publica estado parcial.

### Implementation for User Story 3

- [X] T013 [US3] Adicionar `MenuItem.TargetJourney string` (nome de arquivo `.jornada` relativo ao workspace) em `apps/gerador-jornada/internal/catalog/menu.go` e validar exclusividade mútua com `Submenu` e `TargetScreenID` (máximo um por item; folha = `Submenu == nil`). Restrição do data-model: "Submenu, TargetScreenID e TargetJourney são mutuamente exclusivos".
- [X] T014 [US3] Bump de `catalogVersion` (3 → 4) e estender o codec em `apps/gerador-jornada/internal/adapter/projectzip/wire.go` para ler/gravar `TargetJourney`; leitor aceita pares históricos 1/1, 2/1, 3/2 e 3/3 por migração explícita; primeiro save incompatível avisa e conserva cópia original antes da publicação atômica.
- [X] T015 [US3] Estender `apps/gerador-jornada/internal/application/preview.go` para resolver `TargetJourney`: carregar a `.jornada` destino e entrar na tela inicial dela (`Journey.StartScreenID`), mantendo pilha simples de jornadas e "Voltar" para a origem; caminho percorrido é transitório (não persistido). Restrição do data-model: "Prévia Tela Atual não segue destinos".
- [X] T016 [US3] Em `apps/gerador-jornada/internal/adapter/desktop/preview_compounds_view.go` (`menuView.activate`), no caso `MenuSimulated` de folha sem destino, exibir sugestão imediata para criar uma jornada com o nome da folha (diálogo de confirmação), além da mensagem de seleção simulada.
- [X] T017 [US3] Ao aceitar a sugestão, criar o `.jornada` em branco no workspace com `Project.Name` = rótulo da folha, vincular `TargetJourney` ao item e atualizar a árvore da Estrutura projeto; cancelar não publica estado parcial e o histórico preserva o último estado confirmado.
- [X] T018 [US3] Implementar a navegação: ativar folha com `TargetJourney` carrega a jornada destino (substitui a prévia) e o comando "Voltar" retorna à jornada de origem; nenhuma rota inferida por Label/Shortcut/MovementCode.
- [X] T019 [US3] Adicionar testes para destino entre jornadas (exclusividade, arquivo ausente diagnosticado sem substituição, navegação round-trip, cancelamento sem estado parcial) em `apps/gerador-jornada/internal/catalog/menu_journey_test.go` e `apps/gerador-jornada/internal/application/preview_journey_test.go`.

**Checkpoint**: US3 funcional; as três histórias funcionam juntas.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Rastreabilidade, validação e documentação.

- [X] T020 Atualizar ponteiros de formato/backlog em `specs/002-projeto/spec.md`, `specs/003-catalogo/spec.md` e `specs/005-jornada/spec.md` (novo campo `TargetJourney` e `catalogVersion` 4), sem duplicar o backlog.
- [X] T021 Executar os cenários C1/C2/C3 de `specs/034-launcher-estrutura-projeto/quickstart.md` e registrar revisão/resultado no quickstart da coordenadora (`specs/001-gerador-jornada/quickstart.md`).
- [X] T022 Documentar em português (GoDoc) os novos tipos/funções (recentes, workspace, `TargetJourney`, pilha de navegação) conforme constituição IV.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Sem dependências.
- **Foundational (Phase 2)**: Depende do Setup — BLOQUEIA todas as histórias.
- **User Stories (Phase 3+)**: Todas dependem da Fase 2; depois podem seguir em paralelo (se houver equipe) ou em ordem de prioridade (US1 → US2 → US3).
- **Polish (Phase 6)**: Depende das histórias desejadas concluídas.

### User Story Dependencies

- **US1 (P1)**: começa após a Fase 2; sem dependência das outras histórias.
- **US2 (P1)**: começa após a Fase 2; usa a listagem do workspace (T003); independe de US1.
- **US3 (P1)**: começa após a Fase 2; depende de T003 (workspace) e do campo/codec (T013/T014), mas não de US1/US2 completas.

### Within Each User Story

- Campo/contrato antes de serialização; validação antes de runtime; núcleo antes de integração.

### Parallel Opportunities

- T002 e T003 são [P] (arquivos distintos: `recents.go` e `workspace.go`).
- Após a Fase 2, US1, US2 e US3 podem avançar em paralelo (arquivos majoritariamente distintos).
- Dentro de US3, T013 (catálogo) e T016 (view da prévia) tocam arquivos distintos e podem ser preparados em paralelo antes de T015/T017/T018.

---

## Parallel Example: Fase 2 (Foundational)

```text
Task: "T002 [P] Implementar armazenamento de projetos recentes em apps/gerador-jornada/internal/adapter/desktop/recents.go"
Task: "T003 [P] Implementar resolução do workspace e listagem .jornada em apps/gerador-jornada/internal/adapter/desktop/workspace.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Completar Fase 1 (baseline) e Fase 2 (recentes + workspace).
2. Completar Fase 3 (US1 launcher).
3. **PARAR e VALIDAR**: testar US1 sozinha.
4. Demonstrar se pronto.

### Incremental Delivery

1. Setup + Foundational → fundação pronta.
2. US1 → validar → demonstrar (MVP).
3. US2 → validar → demonstrar (árvore + criar/importar).
4. US3 → validar → demonstrar (criação rápida + navegação).
5. Polish → rastreabilidade + quickstart.

---

## Notes

- [P] tasks = arquivos distintos, sem dependência.
- [Story] liga a tarefa à história para rastreabilidade.
- Cada história é completável e testável de forma independente.
- Validar em checkpoint antes de avançar.
- Evitar: tarefas vagas, conflito no mesmo arquivo, dependências cruzadas que quebrem a independência.

## Evidência de conclusão técnica

T001: baseline em `apps/gerador-jornada/build/034-baseline.log`. T002–T019: código e testes de recentes/workspace/diálogos/navegação/codec, com C1/C2/C3 descritos no [quickstart](quickstart.md). T020: ponteiros em 002/003/005 e contrato 001 atualizados para 3/4. T021: execução técnica e capturas registradas em 034/quickstart e 001/quickstart §26; não significa homologação humana. T022: tipos, funções e invariantes novos documentados em português. Checks/build e smoke nativo aprovados na 1.0.0.4; aceite humano permanece separado em 001/T156. Checklists de revisão preservadas.

## Correção após uso do executável — 1.0.0.5

- [X] T023 Reproduzir o encerramento ao criar e reabrir pelo launcher no driver Windows real. Corrigir a ordem de instalação do tema e inicialização do Select sem callback; preservar troca manual de tema e contraste do título. Acrescentar regressão de Settings/janelas, smoke nativo com criação/reabertura nos três temas e com exemplo, fases de log sem conteúdo do projeto, checks e EXE 1.0.0.5. Evidência em quickstart; aceite humano continua em 001/T156.

- [X] T024 Implementar FR-014: checkbox “Mostrar splash screen ao iniciar” abaixo de abrir projeto, padrão desmarcado, restaurar/salvar escolha imediatamente e reverter em falha. Regressões em desktop/startup_preferences_test.go; política e armazenamento separados em 035/036.
- [ ] T025 Homologar opção da splash com o usuário (marcar, reabrir com vídeo, desmarcar, reabrir direto), junto de 036/T009. Checks e smoke não constituem aceite.

Correção de pasta em 2026-10-25: implementação e regressão pertencem a `037-integracao-github/T072–T073`; este domínio mantém spec/plan/contrato sincronizados sem duplicar tarefa. A criação usa diretamente a pasta escolhida e preserva homônimos.

## Extração planejada para a 038 — 2026-10-25

Estado atual: extração implementada durante a execução autorizada da 038; evidências abaixo. A geração inicial de tarefas foi somente documental. A integração 038/T012 utiliza este resultado; a coordenação de writers está em 038/T011 e a adaptação Git em 037/T080–T081, sem duplicação de implementação.

- [X] T026 Criar regressões de extração em `apps/gerador-jornada/internal/adapter/workspace/store_test.go`, cobrindo as operações reais de `desktop/workspace_store.go`: pasta existente, raiz sem descendentes, nome Windows, links, homônimo, criação/importação exclusiva, exemplo/default de fundo e arquivo ilegível sem omissão silenciosa. Dep.: 038/T005. Refs.: 038/FR-002/FR-005/FR-015/FR-036.
- [X] T027 Extrair somente armazenamento/resolução/criação/importação reutilizáveis para `apps/gerador-jornada/internal/adapter/workspace/store.go`, mantendo regras de domínio fora do IO, codec existente e configuração por raiz autorizada; não adotar restrição Git de repositório ancestral para workspaces locais em geral. Dep.: T026. Refs.: 038/FR-005/FR-013/FR-036.
- [X] T028 Adaptar `apps/gerador-jornada/internal/adapter/desktop/workspace_store.go` e launcher/workspace existentes para o adaptador extraído, com aliases/fachada mínima quando necessário à compatibilidade e testes de regressão em `desktop/workspace_test.go`; preservar modo manual, recentes, tema, defaults e operações sem Git/LLM. Dep.: T027. Refs.: 038/FR-001/FR-002/FR-005; integração 038/T012.

Homologação do modo MCP/EXE é gate da 038. Não marcar tarefas antigas do launcher como aprovadas pela existência desta extração.

Evidência da extração 038: store_test/scope_test e guard_windows_test PASS -race; regressões desktop TestWorkspace/TestLauncher PASS; go vet dos pacotes extraídos/desktop PASS. Fachada mantém chamadas internas e defaults do launcher. A criação/importação novas usam gate/handle; a coordenação de todas as substituições GUI/Git continua em 038/T040 e 037/T080–T081. Nenhum aceite humano ou alteração de tema é inferido desses testes.
