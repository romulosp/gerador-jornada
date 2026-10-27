# Tasks: Conexão do workspace com GitHub

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contratos](contracts/), [DESIGN.md](DESIGN.md) e [quickstart.md](quickstart.md).
**Feature**: `037-integracao-github` — identificador documental, sem branch criada.
**Estado**: incremento técnico 1.0.0.12 implementado e empacotado em 2026-10-24. Checkboxes registram somente escopo técnico comprovado; testes de integração UI adicionais, homologação GitHub real, matriz DPI e aceite permanecem explicitamente pendentes. Evidências no quickstart.
**Tests**: incluídos porque spec/aceitação, contratos e quickstart exigem verificação de persistência, navegação, falhas e publicação. Escrever testes comportamentais antes da implementação da capacidade; demonstrar falha antes da correção quando viável. Erro de compilação por símbolo ausente não é evidência de regressão comportamental.
**Organization**: fases por história, seguidas por validação transversal. Cada ID é local a esta feature; referências externas devem ser qualificadas pelo caminho.

## Formato e convenções

- Cada tarefa começa com checkbox pendente e ID sequencial. `[P]` identifica arquivos distintos que podem ser trabalhados simultaneamente após suas dependências; não autoriza execução antes delas.
- Os caminhos de código são relativos à raiz do repositório. Arquivos novos são destinos planejados, não prova de implementação existente. Um módulo Go, composição manual e domínio existente; não criar banco, backend hospedado ou framework de DI.
- Toda tarefa de código inclui GoDoc/comentários em português conforme a constituição, tratamento de erros sanitizado e seus checks pertinentes. Testes não podem publicar em repositório real por padrão.
- Tarefa técnica pode terminar com testes automatizados e evidência pertinente; tarefas manuais e aceite possuem IDs separados. Não marcar itens para transpor gate nem confundir cadastro OAuth ausente com conclusão funcional.

## Fase 1: Setup compartilhado

**Objetivo**: preparar dependência incorporada, fixtures isoladas e configuração de distribuição, sem mudar o formato `.jornada`.

- [X] T001 Fixar `github.com/go-git/go-git/v5 v5.19.3` e transitivas em `apps/gerador-jornada/go.mod` e `apps/gerador-jornada/go.sum`; registrar licença Apache-2.0 e inventário em `apps/gerador-jornada/README.md`, sem Git do PATH ou atualização `@latest`.
- [X] T002 [P] Criar fixtures reutilizáveis de servidores HTTP, repositórios go-git temporários e injeção de falhas em `apps/gerador-jornada/internal/synctest/fixtures.go`, incluindo duas jornadas interligadas, arquivos alheios e snapshots comparáveis; usar apenas recursos descartáveis. Dep.: T001.
- [X] T003 [P] Implementar geração/limpeza de configuração OAuth privada em `apps/gerador-jornada/scripts/build.py` e `apps/gerador-jornada/scripts/build.ps1`, com testes em `apps/gerador-jornada/scripts/test_configure_github.py`, leitura do par de variáveis definido no contrato e stub em `apps/gerador-jornada/internal/adapter/github/config.go`; ignorar arquivo gerado em `.gitignore`, não emitir valores em logs/inventário e manter build sem cadastro utilizável. Dep.: T001.

## Fase 2: Fundação bloqueante

**Objetivo**: estabelecer contratos internos, persistência segura e coordenação local antes das histórias.

- [X] T004 Definir estados comuns, erros tipados sanitizados e interfaces pequenas de rede/cofre/configuração/Git no consumidor em `apps/gerador-jornada/internal/application/sync_model.go` e `apps/gerador-jornada/internal/application/workspace_sync.go`; usar context.Context em I/O, structs concretas nos adaptadores e nenhum token nos resultados da UI. Cobertura: FR-015/017/018.
- [ ] T005 [P] Escrever testes de caminho de configuração isolado, escrita atômica, corrupção, versão desconhecida e ausência de segredos em `apps/gerador-jornada/internal/adapter/userconfig/store_test.go`, verificando que estado inválido não é sobrescrito. Dep.: T004.
- [X] T006 Implementar armazenamento comum em `apps/gerador-jornada/internal/adapter/userconfig/store.go`, reutilizando `apps/gerador-jornada/internal/adapter/atomicfile/write.go`: `JORNADA_CONFIG_DIR` ou os.UserConfigDir()/JornadaStudio, com a restrição literal "JSON local: `version: 1`, lowerCamelCase, instantes UTC RFC3339"; fornecer primitives para conta, vínculo e journal sem gravar no workspace. Dep.: T005.
- [X] T007 [P] Implementar mutex por raiz, trava entre processos e contexto cancelável em `apps/gerador-jornada/internal/adapter/gitworkspace/gate_windows.go`, com testes em `apps/gerador-jornada/internal/adapter/gitworkspace/inspect_test.go`; garantir "Uma mutação por raiz entre processos da aplicação", detectar locks Git e liberar recursos em erro/cancelamento. Dep.: T004; raiz canônica será validada por T024.
- [X] T008 Refatorar o salvamento de `apps/gerador-jornada/internal/adapter/desktop/files.go` para callback explícito de conclusão/cancelamento/erro, com regressão em `apps/gerador-jornada/internal/adapter/desktop/github_panel_test.go`; preservar Ctrl+S/Salvar como, revisão salva e confirmação de migração, sem iniciar rede ou publicação. Dep.: T006–T007. Cobertura: FR-008/009/018.

**Checkpoint**: T001–T008 concluídas permitem iniciar US1. O editor continua abrindo/salvando sem OAuth configurado ou internet.

## Fase 3: US1 — Conectar conta e escolher repositório (P1, primeiro incremento)

**Objetivo**: autorizar pelo navegador, retornar à mesma janela e selecionar destino elegível, sem inicializar Git ou publicar ainda.
**Teste independente**: conta de teste/fixtures com mais de 100 repositórios, privados/organização e sem escrita; retorno válido, negado, cancelado e repetido. A seleção deve entregar Repository ao fluxo de vínculo futuro sem executar envio.

- [X] T009 [US1] Definir Account/Repository em `apps/gerador-jornada/internal/application/sync_model.go` e registro secreto em `apps/gerador-jornada/internal/adapter/credential/store_windows.go`, preservando as restrições "Uma ativa; identidade pelo ID estável", "Escrita e destino ativo exigidos para selecionar" e "Exclusivo do cofre, nunca em JSON comum"; refreshToken/expiresAt são opcionais. Cobertura: FR-004/017.
- [ ] T010 [P] [US1] Escrever testes OAuth/loopback em `apps/gerador-jornada/internal/adapter/github/oauth_test.go`: PKCE, state, método/Host/rota, porta efêmera, expiração de 10 minutos, negação, cancelamento, callback repetido e encerramento do listener, sem credenciais reais. Dep.: T009. Cobertura: FR-001/002/018.
- [ ] T011 [P] [US1] Escrever testes REST em `apps/gerador-jornada/internal/adapter/github/client_test.go`: headers, paginação Link além de 100, deduplicação por ID, busca completa, listas vazias, privados/organização, escrita, arquivados, 401/403/404/429 e bloqueio de redirecionamento de token para outro host. Dep.: T009. Cobertura: FR-003/017.
- [ ] T012 [P] [US1] Escrever testes de blob versionado, Load/Save/Delete, ausência/erro de cofre e namespace isolado em `apps/gerador-jornada/internal/adapter/credential/store_windows_test.go`; limpar somente credenciais de teste e não usar fallback em texto puro. Dep.: T009. Cobertura: FR-017.
- [X] T013 [US1] Implementar OAuth App público com S256/state, navegador injetável e loopback de uso único em `apps/gerador-jornada/internal/adapter/github/oauth.go`, incluindo troca de código pelo Go, scopes repo/offline_access, resposta no-store/no-referrer sem dados pessoais, timeout e cancelamento; não presumir confidencialidade do client_secret distribuído. Dep.: T010 e T003. Cobertura: FR-001/002.
- [X] T014 [US1] Implementar identidade/listagem/revalidação em `apps/gerador-jornada/internal/adapter/github/client.go` conforme endpoints do contrato, timeout REST de 30 s, TLS, paginação validada e erros sanitizados; não misturar type com visibility/affiliation nem afirmar busca completa antes do fim. Dep.: T011. Cobertura: FR-003/004/017.
- [X] T015 [US1] Implementar WinCred por x/sys/windows em `apps/gerador-jornada/internal/adapter/credential/store_windows.go` e stub indisponível em `apps/gerador-jornada/internal/adapter/credential/store_other.go`; usar chave por configuração/accountID, validar limites e manter token fora dos retornos públicos. Dep.: T012. Cobertura: FR-017.
- [X] T016 [US1] Implementar Connect e seleção de repositórios em `apps/gerador-jornada/internal/application/github_connection.go`, persistência AccountPreference em `apps/gerador-jornada/internal/adapter/userconfig/store.go` e testes em `apps/gerador-jornada/internal/application/github_connection_test.go`; conferir identidade/scopes, manter conta anterior em cancelamento e só declarar conexão concluída após cofre e preferências consistentes. Dep.: T013–T015. Cobertura: FR-001–004/015/017/018.
- [X] T017 [US1] Criar controles de conexão em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go` e incorporá-los em `apps/gerador-jornada/internal/adapter/desktop/workspace.go`; seguir DESIGN.md, incluir ícone vetorial estático licenciado em `apps/gerador-jornada/internal/adapter/desktop/assets/github.svg`, desabilitar sem workspace e manter cancelamento/erro acionáveis. Dep.: T016. Cobertura: FR-001/015.
- [X] T018 [US1] Implementar seletor 720×520 redimensionável em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go`, com lista rolável, busca completa em background, motivo de inelegibilidade e confirmação somente com escrita; entregar seleção ao coordenador, sem vínculo implícito. Dep.: T017. Cobertura: FR-003/004.
- [X] T019 [US1] Compor dependências em `apps/gerador-jornada/cmd/gerador-jornada/main.go` e `apps/gerador-jornada/internal/adapter/desktop/editor.go`; executar rede em goroutines com contextos e resultados por fyne.Do, Show/RequestFocus na janela existente e abertura do seletor após autenticação válida. Dep.: T018. Cobertura: FR-002/018.
- [ ] T020 [US1] Verificar entrada, retorno à mesma janela, lista vazia/paginada, cancelamento e build sem cadastro nos testes `apps/gerador-jornada/internal/adapter/desktop/github_panel_test.go`, registrando checks técnicos do incremento em `specs/037-integracao-github/quickstart.md`, sem marcar M01/M03 reais como executados por fixtures. Dep.: T019. Cobertura: FR-001–004/015/018; SC-001.

**Checkpoint**: US1 verificável por autorização/listagem/seleção, sem depender de US2–US4. Não anunciar publicação pronta neste marco.

## Fase 4: US2 — Vincular pasta e integrar conteúdo existente (P1)

**Objetivo**: usar a raiz exata, aceitar remoto vazio ou com conteúdo compatível e interromper colisões sem perder dados.
**Teste independente**: repositórios temporários e cópias de workspace com inventário de bytes e histórico antes/depois; nenhum commit/push no vínculo. Testes do adaptador usam Repository e credencial de fixture; integração UI depende de US1.

- [X] T021 [US2] Definir WorkspaceBinding/OperationJournal/BindingPlan em `apps/gerador-jornada/internal/application/sync_model.go`, preservando "Um destino por raiz; branch padrão confirmada" e "Intenção e OID planejado gravados antes de mover HEAD; nenhum segredo"; ID de plano opaco e snapshot validado. Cobertura: FR-005/006/007/018.
- [ ] T022 [P] [US2] Escrever testes de raiz exata, Git em ancestral, origin/branch incompatíveis, staging, locks e recursos não suportados em `apps/gerador-jornada/internal/adapter/gitworkspace/inspect_test.go`, incluindo nomes reservados, caixa, acentos, symlink/junction e colisões arquivo/diretório. Dep.: T021. Cobertura: FR-005–007/018.
- [X] T023 [P] [US2] Escrever testes de vínculo vazio/com conteúdo e falhas parciais em `apps/gerador-jornada/internal/adapter/gitworkspace/binding_test.go`, verificando igualdade de bytes, ausência de envio e recuperação sem excluir arquivos anteriores. Dep.: T021. Cobertura: FR-007/008/018.
- [X] T024 [US2] Implementar inspeção/canonicalização em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, rejeitando bare/linked worktree/sparse/submódulo/operação incompleta, origem incompatível e estrutura que escape da raiz; resolver comparação Windows sem alterar nome físico. Dep.: T022. Cobertura: FR-005/006/018.
- [X] T025 [US2] Implementar fetch privado e consulta de refs em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, com token somente no transporte, refspec explícito, contexto cancelável de até 15 minutos e branch REST/default main apenas quando omitida, exibida antes de confirmar. Dep.: T024 e T014. Cobertura: FR-006/007/017/018.
- [X] T026 [US2] Implementar inventário e plano de materialização em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`: comparar tipo/bytes de homônimos, preservar exclusivos locais e remotos e bloquear caminhos inválidos, .git, LFS e colisões antes da primeira escrita. Dep.: T025. Cobertura: FR-007.
- [X] T027 [US2] Implementar journal/backup privado em `apps/gerador-jornada/internal/adapter/gitworkspace/journal.go` e persistência de vínculo em `apps/gerador-jornada/internal/adapter/userconfig/store.go`, com pré-imagens/hashes, fases atômicas e limpeza somente após conclusão; divergência externa preserva recuperação e localização do backup. Dep.: T026 e T006. Cobertura: FR-007/018.
- [X] T028 [US2] Implementar aplicação inicial em `apps/gerador-jornada/internal/adapter/gitworkspace/journal.go`: preparar .git privadamente, índice/branch na base remota, publicar apenas ausentes compatíveis e gravar vínculo por último; sem commit/push ou substituição de arquivo anterior. Dep.: T027. Cobertura: FR-005–008.
- [X] T029 [US2] Implementar vínculo de Git existente em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, preservando config/histórico compatíveis e bloqueando ausência de ancestral comum ou commits locais desconhecidos; preparar avanço não conflitante necessário ao vínculo sem enviar histórico. Dep.: T028. Cobertura: FR-005–007.
- [X] T030 [US2] Implementar InspectWorkspace/PlanBinding/ApplyBinding em `apps/gerador-jornada/internal/application/workspace_sync.go`, com revalidação de ID/permissão/branch e estado anterior em erro, plano invalidado por mudança externa e trava por raiz. Dep.: T029 e T016. Cobertura: FR-004–007/018.
- [X] T031 [US2] Ligar seleção ao vínculo em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go`, mostrando pasta/destino/conteúdo remoto e colisões; usar callback T008 para salvar e pedir confirmação antes de recarregar jornada ativa, recusando substituição de sessão por pacote inválido. Dep.: T030 e T018. Cobertura: FR-007/015/018.
- [X] T032 [US2] Implementar recuperação local idempotente em `apps/gerador-jornada/internal/adapter/gitworkspace/journal.go` com testes em `apps/gerador-jornada/internal/adapter/gitworkspace/fault_test.go`; restaurar só bytes ainda pertencentes à operação e bloquear alteração externa, sem reset hard ou apagar .git anterior. Dep.: T031. Cobertura: FR-007/018.
- [ ] T033 [US2] Verificar fluxo completo de vínculo e árvore em `apps/gerador-jornada/internal/application/workspace_sync_test.go` e `apps/gerador-jornada/internal/adapter/desktop/github_panel_test.go`, incluindo recarga cancelada, falha de escrita e pacote remoto inválido; registrar resultados técnicos em `specs/037-integracao-github/quickstart.md`. Dep.: T032. Cobertura: FR-005–008/018; SC-003.

**Checkpoint**: workspace vinculado e recuperável, sem envio automático. US2 não depende da implementação de Publicar para demonstrar valor.

## Fase 5: US3 — Publicar jornadas e retomar envio (P1)

**Objetivo**: revisão explícita, commit restrito às jornadas e push confirmado; integrar novidades disjuntas e preservar OIDs pendentes.
**Teste independente**: workspace vinculado com duas jornadas interligadas e remoto temporário, comparando blobs/árvores/ancestralidade, staging alheio, conflito, falha de envio e nova tentativa. Nenhum fixture publica no GitHub real.

- [X] T034 [US3] Definir Change/PublicationReview/PublicationResult em `apps/gerador-jornada/internal/application/sync_model.go`, preservando "Arquivo regular na raiz `.jornada`; excluir só se rastreado" e "Imutável, efêmera, mensagem não vazia e snapshot revalidado"; SHA-256 de trabalho e OID Git são distintos. Cobertura: FR-009–013.
- [X] T035 [P] [US3] Escrever testes de inclusão/alteração/exclusão, conjunto de duas jornadas, bytes ZIP, backups/temporários/subpastas, staging alheio, mensagem vazia e revisão obsoleta em `apps/gerador-jornada/internal/adapter/gitworkspace/engine_test.go`. Dep.: T034. Cobertura: FR-009–012.
- [X] T036 [P] [US3] Escrever testes de merge conservador em `apps/gerador-jornada/internal/adapter/gitworkspace/engine_test.go`: fast-forward com edição local, caminhos disjuntos, bytes idênticos, delete/modify, arquivo/diretório, zero/múltiplas bases e pais remoto/local preservados. Dep.: T034. Cobertura: FR-014.
- [X] T037 [P] [US3] Escrever testes de push rejeitado, permissão/política, resultado incerto após aceitação remota, queda entre commit e estado persistido e retentativa sem duplicar OID em `apps/gerador-jornada/internal/adapter/gitworkspace/engine_test.go`. Dep.: T034. Cobertura: FR-013/014/018.
- [X] T038 [US3] Implementar revisão de arquivos em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, enumerando explicitamente todos os `.jornada` regulares da raiz e exclusões rastreadas, calculando snapshot e bloqueando índice preparado externamente; sem seleção parcial ou inferência por pasta recursiva. Dep.: T035. Cobertura: FR-009–012.
- [X] T039 [US3] Coordenar salvar/revisar e diagnóstico de referências afetadas por exclusão em `apps/gerador-jornada/internal/application/workspace_sync.go` e `apps/gerador-jornada/internal/adapter/desktop/files.go`; respeitar callback T008, cancelar/erro de migração e usar destino explícito do contrato 034 sem reparar vínculos automaticamente. Dep.: T038. Cobertura: FR-008/009/011.
- [X] T040 [US3] Implementar commit por blobs/árvore em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, herdando base e auditando somente deltas elegíveis, sem All/add recursivo/filtros; mostrar autoria GitHub/noreply, persistir intenção/OID/pais antes de atualizar ref e manter índice/worktree recuperáveis. Dep.: T039 e T027. Cobertura: FR-010–013/017.
- [X] T041 [US3] Implementar comparação B/L/R por presença/tipo/modo/OID em `apps/gerador-jornada/internal/adapter/gitworkspace/tree.go`, exigindo uma base comum e journal de commits locais; conflito não escolhe versão, merge disjunto tem primeiro pai remoto/segundo local e diff auditado contra remoto, sem rebase. Dep.: T036 e T040. Cobertura: FR-014.
- [X] T042 [US3] Implementar aplicação de novidades remotas em `apps/gerador-jornada/internal/adapter/gitworkspace/journal.go`, com preflight/journal, preservação de alterações locais, invalidação de revisão após novo fetch e confirmação de recarga da jornada ativa via ponte desktop. Dep.: T041 e T031. Cobertura: FR-007/014/018.
- [X] T043 [US3] Implementar push/verificação de ancestralidade em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, refspec somente da branch padrão, sem tags/force; revalidar destino e só confirmar publicado quando OID estiver alcançável no remoto. Dep.: T037 e T042. Cobertura: FR-013/014/017.
- [X] T044 [US3] Implementar reconciliação/reenvio em `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go`, reutilizando OID comprovado, identificando entrega já aceita e integrando/bloqueando divergência; novas edições aguardam outra publicação e HEAD inesperado exige recuperação, sem recriar commits. Dep.: T043 e T032. Cobertura: FR-013/018.
- [X] T045 [US3] Implementar ReviewPublication/Publish/RetryPublication em `apps/gerador-jornada/internal/application/workspace_sync.go` com ID opaco, conta/destino/snapshot revalidados, mensagem editável não vazia, ausência de novidades sem commit vazio e falha remota sem impedir save local. Dep.: T044. Cobertura: FR-008–014/018.
- [X] T046 [US3] Implementar diálogo/botão Publicar em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go`, seguindo DESIGN.md: lista completa, autoria, mensagem, referências afetadas, alcance da permissão e confirmação; apresentar pendência, Verificar envio e Tentar publicar novamente sem exibir sucesso incerto. Dep.: T045. Cobertura: FR-009/013/015.
- [X] T047 [US3] Integrar gate e contextos em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go`, `apps/gerador-jornada/internal/adapter/desktop/files.go` e `apps/gerador-jornada/internal/adapter/desktop/workspace.go`; serializar salvar/criar/importar só durante mutação local, manter rede fora da thread Fyne e permitir edição após erro/cancelamento remoto. Dep.: T046. Cobertura: FR-008/018.
- [ ] T048 [US3] Verificar coordenador/UI em `apps/gerador-jornada/internal/application/workspace_sync_test.go` e `apps/gerador-jornada/internal/adapter/desktop/github_panel_test.go`, cobrindo save sem rede, confirmação, migração cancelada, reenvio separado de novas edições, exclusão referenciada e navegação entre duas jornadas publicadas. Dep.: T047. Cobertura: FR-008–015/018; SC-002/004/005/007.
- [X] T049 [US3] Completar harness de crash/rede e interferência externa em `apps/gerador-jornada/internal/adapter/gitworkspace/fault_test.go`, cobrindo cada fronteira de journal/ref/worktree, política de branch e perda da resposta de push; registrar OIDs e hashes sanitizados em `specs/037-integracao-github/quickstart.md`. Dep.: T048. Cobertura: FR-013/014/018; SC-003/005.

**Checkpoint**: US1–US3 permitem conexão, vínculo e publicação explícita, com falhas exercitadas. Persistência em disco já existe, mas experiência completa de reinício/logout depende de US4.

## Fase 6: US4 — Retomar trabalho e desconectar (P2)

**Objetivo**: restaurar conta/vínculo/pendências, recuperar operações interrompidas, renovar autorização e desconectar sem excluir jornadas.
**Teste independente**: reiniciar após conexão, falha/resultado incerto do push e logout; testar outra conta, revogação e cofre ausente em namespace isolado.

- [ ] T050 [P] [US4] Escrever testes de reinício com conta/vínculo/pendências e JSON inválido em `apps/gerador-jornada/internal/application/github_connection_test.go`, incluindo leitura offline e publicação histórica sem garantia de sincronismo atual. Dep.: T049. Cobertura: FR-015/016.
- [ ] T051 [P] [US4] Escrever testes de refresh opcional/rotacionado, revogação, troca de conta e falha de exclusão de credencial em `apps/gerador-jornada/internal/application/github_connection_test.go` e `apps/gerador-jornada/internal/adapter/github/security_test.go`, sem alterar destino silenciosamente. Dep.: T049. Cobertura: FR-004/016/017.
- [X] T052 [US4] Implementar carga de conta/vínculo/pendências em `apps/gerador-jornada/internal/application/github_connection.go` usando configuração/cofre existentes; preservar as restrições "Credencial ausente pede reconexão sem impedir edição local" e "Última publicação é fato histórico, não garantia de sincronismo atual". Dep.: T050. Cobertura: FR-015/016/018.
- [X] T053 [US4] Implementar refresh em `apps/gerador-jornada/internal/adapter/github/oauth.go` e seu uso em `apps/gerador-jornada/internal/application/github_connection.go`, serializado por conta, persistindo token rotacionado antes de trocar estado e exigindo reconexão quando não houver recuperação; renovar não autoriza push. Dep.: T051 e T052. Cobertura: FR-016/017.
- [X] T054 [US4] Coordenar recuperação na abertura do workspace em `apps/gerador-jornada/internal/application/workspace_sync.go`, reutilizando recovery/journal/pending de US2/US3; reconhecer HEAD planejado, anterior ou inesperado, sem envio automático por mensagem/autor e sem sobrescrever mudança externa. Dep.: T053. Cobertura: FR-013/016/018.
- [X] T055 [US4] Implementar desconexão e troca de conta em `apps/gerador-jornada/internal/application/github_connection.go`: cancelar rede, remover credencial anterior com confirmação de sucesso, preservar vínculo/histórico/pendências, revalidar escrita na nova conta e bloquear novos envios se remoção falhar; não afirmar revogação do OAuth App no site. Dep.: T054. Cobertura: FR-004/016/017.
- [X] T056 [US4] Implementar Status local e Atualizar estado explícito em `apps/gerador-jornada/internal/application/workspace_sync.go`, distinguindo alterações locais, pendingPush/outcomeUnknown, última publicação confirmada e remoto não verificado, sem rede obrigatória na abertura. Dep.: T055. Cobertura: FR-015/016/018.
- [X] T057 [US4] Conectar retomada/status/logout em `apps/gerador-jornada/internal/adapter/desktop/github_panel.go` e `apps/gerador-jornada/internal/adapter/desktop/editor.go`, exibindo conta/destino/branch e ações Atualizar estado/Ver repositório/Desconectar, com foco/erros e confirmação em português. Dep.: T056. Cobertura: FR-015–018.
- [ ] T058 [US4] Verificar ciclo nativo com namespace WinCred isolado em `apps/gerador-jornada/internal/adapter/desktop/github_panel_test.go`, incluindo reinício offline, perda de token, troca de conta, limpeza de credencial e preservação de bytes/OIDs; registrar resultado técnico em `specs/037-integracao-github/quickstart.md`. Dep.: T057. Cobertura: FR-015–018; SC-006.

**Checkpoint**: todas as histórias têm implementação/testes técnicos próprios; homologação real e aprovação continuam separadas abaixo.

## Fase 7: Validação transversal e entrega

- [X] T059 Revisar rastreabilidade, comentários em português e consistência dos caminhos/assinaturas em `specs/037-integracao-github/spec.md`, `specs/037-integracao-github/plan.md`, `specs/037-integracao-github/data-model.md`, `specs/037-integracao-github/contracts/` e `specs/037-integracao-github/tasks.md`; sincronizar somente artefatos de domínios efetivamente afetados, sem mudar escopo para esconder lacuna. Dep.: T058.
- [ ] T060 Executar testes focados, `apps/gerador-jornada/scripts/check.ps1` e checks de empacotamento em `apps/gerador-jornada/scripts/test_build.py`; verificar timeout/cancelamento e feedback até 250 ms no harness, e registrar comandos/resultados observados em `specs/037-integracao-github/quickstart.md`. Dep.: T059.
- [X] T061 Gerar pacote Windows x64 com `apps/gerador-jornada/scripts/build.py`, incluindo licenças/ícone e build sem cadastro OAuth; verificar PE/arquitetura, inventário sem segredos e funcionamento local em computador sem Git, registrando caminho real, versão e SHA-256 em `specs/037-integracao-github/quickstart.md`. Dep.: T060.
- [ ] T062 Preparar cadastro OAuth real e repositórios descartáveis com o proprietário do produto conforme `specs/037-integracao-github/contracts/github-desktop.md`; fornecer valores ao ambiente do build, gerar pacote configurado por `apps/gerador-jornada/scripts/build.py` e registrar em `specs/037-integracao-github/quickstart.md` apenas disponibilidade, revisão/hash e ambiente, nunca credenciais. Dep.: T061; requer recursos externos reais, não substituíveis por valores fictícios.
- [ ] T063 Executar homologação GitHub real M01–M08/M11/M12 e paginação real quando houver volume autorizado, seguindo `specs/037-integracao-github/quickstart.md`; inspecionar bytes, refs, Git/config/cofre/logs e registrar executor, resultado observado e evidência sanitizada da mesma revisão. Dep.: T062. Reprovação ou cenário indisponível permanece explícito, sem copiar resultado esperado como observado.
- [ ] T064 Executar M09/M10/M13 com harness determinístico de `apps/gerador-jornada/internal/adapter/gitworkspace/fault_test.go` e integração da aplicação, registrando reinício real, conflitos e recuperação em `specs/037-integracao-github/quickstart.md`; corrigir falhas e repetir cenários afetados antes de concluir. Dep.: T061 e T049.
- [ ] T065 Homologar M14 segundo `specs/037-integracao-github/DESIGN.md` em Windows, DPI 100/125/150%, claro/escuro e nomes longos; guardar capturas sanitizadas da entrada/seletor/revisão/pendência/conflito e evidência de editor/prévia sem regressão em `specs/037-integracao-github/quickstart.md`. Dep.: T061.
- [ ] T066 Consolidar evidências e disponibilizar executável/roteiro ao usuário em `specs/037-integracao-github/quickstart.md`; registrar pronto_para_aprovacao somente após T063–T065 aprovadas na mesma revisão e registrar aprovado apenas após manifestação humana explícita sobre ela. Dep.: T063–T065; aprovação não autoriza merge/push/publicação do código do produto.

## Dependências e ordem de execução

## Correções do teste manual — 2026-10-25

- [X] T077 Remover a confirmação por digitação rejeitada no teste 1.0.0.13: seletor em `desktop/github_panel.go` usa a seleção diretamente e abre revisão após vínculo bem-sucedido; teste de UI deve comprovar ID selecionado, filtro, cancelamento e ausência de segundo formulário. Cobertura: FR-022/SC-010.
- [X] T078 Classificar colisão de jornadas e registrar fase do vínculo nos logs em `gitworkspace/journal.go` e `engine.go`, com regressão que preserve arquivos e bloqueie envio; não atribuir a causa da tentativa real a erro genérico sem prova. Cobertura: FR-007/018/019.
- [X] T079 Sincronizar artefatos, executar checks e gerar EXE 1.0.0.14 configurado em caminho separado para não interromper a instância antiga em uso; registrar hash e roteiro curto. Reteste GitHub real permanece humano.

- [X] T071 Registrar reprovação de 1.0.0.12 e sincronizar spec/plan/research/contratos/DESIGN/quickstart com FR-007/020–022, sem alterar o formato `.jornada` nem apagar evidências anteriores.
- [X] T072 Reproduzir em `internal/adapter/desktop/launcher_test.go` a pasta duplicada e em `internal/adapter/gitworkspace/journey_only_test.go` a materialização de código-fonte remoto, colisão não elegível, preservação remota e push confirmado para o segundo destino; provar falha antes da correção. Caminhos sob `apps/gerador-jornada/`. Cobertura: SC-009/010.
- [X] T073 Corrigir `internal/adapter/desktop/launcher.go` para criar o pacote diretamente na pasta escolhida, preservando homônimos e ajustando instruções do diálogo. Cobertura: FR-020/021.
- [X] T074 Restringir preflight/journal/materialização a `.jornada` da raiz em `internal/adapter/gitworkspace/`, preservando a árvore Git completa e índice com skip-worktree para os demais arquivos; manter staging externo bloqueado e recuperação coberta. Cobertura: FR-007/010/018; SC-009.
- [X] T075 Corrigir confirmação e resultado em `internal/adapter/desktop/github_panel.go`, testar seleção após filtro e conferir identidade completa no coordenador. Corrigir confirmação com remoto vazio para não reutilizar ref em cache; registrar destino por identidade opaca nos logs. Cobertura: FR-019/022; SC-010.
- [X] T076 Executar regressões, checks e build 1.0.0.13 configurado; entregar EXE/hash e roteiro de reteste em quickstart.md. Preservar os arquivos e o histórico da tentativa anterior; não limpar nem publicar em repositórios reais durante o conserto. Reteste humano permanece pendente.

## Incremento de diagnóstico solicitado pelo usuário — 2026-10-24

- [X] T067 [US4] Implementar correlação independente de OAuth e logger seguro por allowlist em `apps/gerador-jornada/internal/application/sync_trace.go` e `apps/gerador-jornada/internal/adapter/syncdiag/diagnostics.go`, com classes técnicas de erro, DNS/TCP/TLS e sem mensagens/corpos livres. Cobertura: FR-019.
- [X] T068 [US4] Instrumentar navegador, listener, callback, troca OAuth, REST, cofre, transações Git e falhas inesperadas em `apps/gerador-jornada/internal/adapter/github/`, `apps/gerador-jornada/internal/adapter/gitworkspace/` e `apps/gerador-jornada/internal/adapter/desktop/github_panel.go`; apresentar identificador de diagnóstico no erro da aplicação. Cobertura: FR-019.
- [X] T069 [US4] Verificar exclusão de Client ID/Secret, tokens, code/state/PKCE e descrição livre maliciosa em `apps/gerador-jornada/internal/adapter/github/diagnostics_test.go` e `apps/gerador-jornada/internal/adapter/syncdiag/diagnostics_test.go`, preservando request ID validado, código conhecido e sequência de etapas. Cobertura: FR-019/SC-008.
- [X] T070 Documentar uso do Log de aplicação, limites de observação do navegador externo, build e pendências reais em `apps/gerador-jornada/README.md` e `specs/037-integracao-github/quickstart.md`. Dep.: T067–T069.

**Evidência técnica atual**: check.ps1 aprovado; teste Git compilado executado com PATH sem Git aprovado (29 casos/subcasos); build x64 configurado gerado; captura nativa obtida. T005/T010/T011/T012/T020/T022/T033/T048/T050/T051/T058 mantêm caixas abertas quando a descrição pede matriz adicional além dos testes já existentes. T060 tem checks aprovados, mas a medição dedicada de 250 ms ainda não foi registrada. T062 tem cadastro e configuração externa consumidos, mas repositórios/consentimento reais ainda precisam da etapa humana. Nenhuma dessas pendências foi apagada para declarar conclusão integral.

## Ordem após o incremento de diagnóstico

T067–T070 foram acrescentadas por solicitação explícita posterior, sem renumerar as 66 tarefas originais. A demonstração manual deve usar a mesma revisão configurada com os logs. O restante desta seção preserva a ordem original das capacidades; as caixas abertas acima continuam sendo trabalho restante.

Ordem padrão: sequencial por ID, exceto grupos explicitamente independentes. Dependências por faixa incluem todos os IDs da faixa. Completude funcional não é inferida da existência de arquivos.

```text
T001 → {T002, T003} → T004
T004 → {T005 → T006, T007} → T008
T008 → US1 (T009–T020) → US2 (T021–T033)
US2 → US3 (T034–T049) → US4 (T050–T058)
US4 → T059 → T060 → T061 → T062
T062 → {T063, T064, T065} → T066
```

- US1 depende da fundação. US2 usa a conta/seleção de US1 na UI; seus algoritmos são testáveis com fixtures. US3 depende do vínculo/journal de US2. US4 usa credenciais/vínculo/commit persistidos pelas histórias anteriores.
- Testes precedem implementação correspondente: T010→T013, T011→T014, T012→T015, T022→T024, T023→T026–T029, T035→T038–T040, T036→T041–T042, T037→T043–T044, T050→T052 e T051→T053/T055.
- Cadastro real só bloqueia T062/T063 e demonstrações remotas reais; não interrompe código, fixtures ou validação offline. Não declarar o recurso homologado sem esses resultados.
- Histórias não são paralelas integralmente: compartilham composição/estado/UI. Cada incremento tem teste independente de comportamento, sem afirmar independência de infraestrutura inexistente.

## Exemplos de execução paralela

| Grupo | Pré-requisito | Tarefas em arquivos distintos |
|---|---|---|
| Setup | T001 | T002 (fixtures) e T003 (build/configuração). |
| Fundação | T004 | T005 (testes de config) e T007 (gate); T006 espera T005. |
| US1 | T009 | T010 (OAuth), T011 (REST) e T012 (WinCred), todos testes de contrato distintos. |
| US2 | T021 | T022 (inspeção) e T023 (vínculo/falhas), antes dos adaptadores correspondentes. |
| US3 | T034 | T035 (commit), T036 (integração) e T037 (envio), com repositórios temporários próprios. |
| US4 | T049 | T050 (restore) e T051 (refresh/logout), sem compartilhar credenciais reais. |

T063–T065 só podem ser executadas simultaneamente em ambientes/cópias independentes; a atualização de `quickstart.md` deve ser serializada. `[P]` descreve possibilidade técnica, não exige subagentes ou execução concorrente.

## Estratégia de implementação e entrega

1. **Primeiro incremento verificável**: T001–T020, autenticação e seleção. É MVP da US1; não representa a publicação solicitada inteira. Demo real exige cadastro do proprietário; fixtures permitem avanço técnico.
2. **Workspace conectado**: concluir US2; validar remoto vazio/com conteúdo e bloqueios sem perda, antes de permitir publicar.
3. **Fluxo principal completo**: concluir US3; provar commit do conjunto, preservação de nomes/bytes, integração disjunta e retentativa por OID.
4. **Uso contínuo**: concluir US4 e checks transversais; disponibilizar build utilizável com estado de teste manual explícito.
5. **Homologação e aprovação**: executar tarefas externas/manuais, corrigir e retestar. Nenhuma caixa é marcada nesta geração. A análise cruzada `speckit-analyze` deve ocorrer antes de implementar.

## Rastreabilidade por requisito

Os requisitos abaixo são sempre `specs/037-integracao-github/spec.md/FR-xxx`; tarefas e cenários têm namespace desta feature.

| Requisito | Implementação principal | Verificação planejada |
|---|---|---|
| FR-001 | T013/T016–T019 | T010/T020; M01 |
| FR-002 | T013/T016/T019 | T010/T020; M01/M02 |
| FR-003 | T014/T018 | T011/T020; M03 |
| FR-004 | T009/T016/T030/T055 | T020/T051/T058; M01/M12 |
| FR-005 | T024/T028–T030 | T022/T023/T033; M04/M05 |
| FR-006 | T025/T029/T030 | T022/T033; M04/M05 |
| FR-007 | T026–T032/T042 | T023/T033/T049; M05/M13 |
| FR-008 | T008/T028/T039/T045/T047 | T033/T048; M04/M06/M07 |
| FR-009 | T008/T034/T038/T039/T046 | T035/T048; M06/M07 |
| FR-010 | T038/T040/T041 | T035/T048; M06/M08 |
| FR-011 | T039/T040 | T035/T048; M06/M08 |
| FR-012 | T038/T045 | T035/T048; M08 |
| FR-013 | T040/T043–T045/T054 | T037/T049/T058; M09/M13 |
| FR-014 | T041–T045 | T036/T037/T049; M10/M11 |
| FR-015 | T017/T046/T052/T056/T057 | T020/T048/T050/T058; M12/M14 |
| FR-016 | T052–T056 | T050/T051/T058; M09/M12 |
| FR-017 | T003/T015/T040/T043/T053/T055 | T012/T051/T058/T061/T063; M12 |
| FR-018 | T007/T008/T013/T027/T032/T042/T047/T054 | T010/T023/T049/T058; M02/M11/M13 |

SC-001 → T020/T063; SC-002 → T048/T063; SC-003 → T033/T049/T063/T064; SC-004 → T048/T058/T063; SC-005 → T037/T049/T064; SC-006 → T058/T063; SC-007 → T035/T048/T061/T063. Os casos M01–M14 são definidos no quickstart, não resultados já obtidos.

## Coordenação de gravações com a 038 — 2026-10-25

Tarefas futuras decorrentes da fronteira compartilhada, sem nova autorização de rede/publicação. O gate comum é implementado por 038/T011; estas tarefas tratam somente a integração da engine Git já existente, não refazem OAuth, merge ou journal geral.

- [ ] T080 Criar regressão de writer Git versus GUI/MCP em `apps/gerador-jornada/internal/adapter/gitworkspace/authoring_gate_test.go`, usando fixtures locais sem rede/conta: trava por raiz física, recuperação pendente, alterações concorrentes, preservação de bytes e separação salvar/publicar. Dep.: 038/T007 e contrato de gate 038/T011. Refs.: 038/FR-012/FR-013/FR-017/FR-038.
- [ ] T081 Integrar `apps/gerador-jornada/internal/adapter/gitworkspace/engine.go` e `journal.go` ao gate comum de `internal/adapter/workspaceguard`, compartilhando namespace/ordem de coordenação na aplicação/materialização e sem segurar a raiz durante consentimento ou rede; revalidar base e tratar recuperação sem restaurar arquivo mudado por terceiro. Dep.: T080 e 038/T011. Refs.: 038/FR-012/FR-013/FR-017; integração final 038/T040.

Somente a confirmação específica de publicação continua autorizando commit/push Git. Host MCP autorizado a salvar não recebe autoridade para publicar; testes/homologação remota antigos mantêm seus estados.
