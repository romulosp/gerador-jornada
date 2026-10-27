# Research: Integração GitHub desktop

**Data**: 2026-10-24. Pesquisa e decisões reconciliadas com a implementação de código 1.0.0.12. Dependências e testes locais já foram executados; autorização no navegador e publicação no GitHub real continuam sem homologação registrada nesta atualização. Teste local não substitui integração externa real.

## Autorização desktop

**Decisão**: OAuth App público com authorization code, PKCE S256, state e callback loopback. O GitHub exige client_secret na troca e documenta que, em cliente nativo, ele é distribuído e recuperável. Não tratá-lo como fronteira de confiança. Token de usuário é credencial protegida. [Boas práticas oficiais](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/best-practices-for-creating-an-oauth-app).

**Justificativa**: retorno automático sem backend próprio. Registrar callback base `http://127.0.0.1/oauth/github/callback`; porta efêmera pela exceção loopback. O proprietário fornece homepage real do produto, sem URL inventada. Fluxo e parâmetros seguem [autorização OAuth](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps).

**Alternativas**: device flow exige código manual e não atende ao retorno aprovado; GitHub App acrescenta instalação por repositório; broker hospedado acrescenta operação não solicitada. Revisar se surgir exigência de token limitado ao repositório ou proibição de credencial pública distribuída.

## Escopos e seleção

**Decisão**: solicitar `repo` e `offline_access`, suportar token com ou sem refresh/expiração e conferir escopos efetivos. O seletor limita o destino operacional; não limita o alcance de `repo`. Não solicitar workflow, e-mail privado ou administração. [Escopos](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps).

**Justificativa**: repositórios privados e de organização precisam do acesso autorizado. Listagem REST paginada por ID, visibilidade convertida para o modelo e escrita revalidada antes de publicar; arquivados/desabilitados não são destinos elegíveis. Todas as páginas são obtidas em background antes de exibir a busca global, sem Carregar mais incremental visível. Uma falha não é apresentada como lista completa; FR-003 continua exigindo todos os destinos autorizados. [REST de repositórios](https://docs.github.com/en/rest/repos/repos#list-repositories-for-the-authenticated-user).

**Alternativas**: listar só públicos ou próprios omite cenários pedidos; PAT manual altera o fluxo. Políticas de organização/SSO ainda podem limitar acesso.

## Credenciais

**Decisão**: Windows Credential Manager, CRED_TYPE_GENERIC, persistência local, chave por configuração e conta. Usar x/sys/windows existente para WinCred; JSON só guarda referência opaca. [CredWriteW](https://learn.microsoft.com/en-us/windows/win32/api/wincred/nf-wincred-credwritew), [CredReadW](https://learn.microsoft.com/en-us/windows/win32/api/wincred/nf-wincred-credreadw).

**Justificativa**: proteção nativa sem criptografia própria; falha do cofre não produz conexão persistente fictícia. O blob usa protectedRecord com version=1 e token, recusando versão desconhecida ou token vazio. Round-trip/remoção nativos usam token sintético em namespace descartável; não comprovam login GitHub real.

**Alternativas**: texto puro é rejeitado; DPAPI com arquivo acrescenta gestão desnecessária. Não há fallback inseguro.

## Mecanismo Git

**Decisão**: go-git v5.19.3, Apache-2.0, exige Go 1.26 e atende ao módulo Go 1.27.2. Versão/transitivas já estão fixadas em go.mod/go.sum; avisos de licença e custo do executável devem acompanhar a verificação da build real. [Release](https://github.com/go-git/go-git/releases/tag/v5.19.3), [go.mod](https://raw.githubusercontent.com/go-git/go-git/v5.19.3/go.mod), [licença](https://github.com/go-git/go-git/blob/v5.19.3/LICENSE).

**Justificativa**: init/fetch/objetos/push HTTPS incorporados; token apenas em memória. Merge geral não é pressuposto: integração conservadora própria usa árvores e commits com dois pais. [API versionada](https://pkg.go.dev/github.com/go-git/go-git/v5@v5.19.3), [objetos e MergeBase](https://pkg.go.dev/github.com/go-git/go-git/v5@v5.19.3/plumbing/object).

**Alternativas**: Git do PATH cria pré-requisito; Git empacotado aumenta distribuição/manutenção; API Contents não cria Git local. Rever se LFS, assinatura ou recursos avançados se tornarem requisitos.

## Integração, autoria e retomada

**Decisão**: comparar árvores por caminho, presença, modo e bytes exatos; alterações disjuntas integram, mesmo arquivo divergente bloqueia. OIDs identificam commits/árvores e SHA-256 identifica o conteúdo revisado sem confundir os dois formatos. Usar journal anterior à mutação, snapshot imutável e verificação de refs. Bloquear índice com staging alheio; construir commit somente com `.jornada` elegíveis. Pendência é OID comprovado pela intenção transacional, não mensagem ou autor.

**Justificativa**: bloqueio de toda divergência impediria US3.9; merge textual de ZIP não tem contrato confiável. Primeiro pai de merge conservador é remoto, segundo local: delta contra o remoto contém só jornadas locais, preservando ambos os históricos. Retentativa preserva OIDs sem rebase.

**Autoria escolhida**: a revisão mostra o login público; o commit usa nome público GitHub ou login e e-mail noreply `<id>+<login>@users.noreply.github.com`. Sem alterar configuração Git global. O GitHub aceita e-mail noreply para autoria, conforme [configuração de e-mail de commit](https://docs.github.com/en/account-and-profile/how-tos/email-preferences/setting-your-commit-email-address). Política de assinatura/PR pode bloquear o envio; não contorná-la.

**Alternativas**: `add .`, All, force push, reset destrutivo e recriar commits em retentativas são rejeitados. Contrato completo em [workspace-git.md](contracts/workspace-git.md).

## Primeiro envio e verificação de recuperação

**Decisão**: Head/Remote vazios representam ausência de commit, não objetos Git nulos a carregar. A base inicial é árvore vazia; primeiro commit tem zero pais e cria a referência da branch padrão. A implementação local cobre esse fluxo com `TestFirstPublicationEmptyRemote`.

**Justificativa**: vincular remoto vazio não pode criar commit/push, mas a publicação posterior precisa funcionar sem ancestral existente. Retentativas preservam o mesmo OID. Os testes também exercitam dois pais em integração disjunta, colisão/rollback e encerramento abrupto de um subprocesso antes dos defers, preservando intenção e lock.

**Alternativas**: criar um commit artificial durante o vínculo violaria a publicação explícita; tratar OID vazio como hash válido falharia no primeiro uso. Os testes usam repositórios bare descartáveis e não comprovam permissões/restrições do GitHub real.

## Coordenação e apresentação

**Decisão**: `application.WorkspaceSync` expõe Bind, ReviewPublication e Publish e mantém tokens no backend. Os arquivos planejados foram consolidados em `gitworkspace/engine.go`, `tree.go`, `journal.go`, `status.go` e `desktop/github_panel.go`; não é uma mudança de domínio ou formato.

**Justificativa**: revisão por valor, OIDs e digests revalidados substitui o ID opaco proposto, mantendo rejeição de snapshot obsoleto. Painel com status truncado, detalhes em ⋮ e uma ação principal preserva a altura da árvore. Ícone Octicons MIT é incorporado e acompanha o tema. A lista completa em background simplifica busca sem omitir páginas.

**Alternativas**: manter todos os botões e linhas empilhados ocultava jornadas no painel estreito; lista incremental exigiria distinguir pesquisa parcial de completa. A decisão mantém as capacidades e seus critérios de aceitação; não dispensa teste visual/DPI ou autorização real.

## Revisão após teste manual de 2026-10-25

**Evidência**: o workspace testado estava vinculado a `romulosp/gerador-jornada`, diferente do destino vazio `romulosp/repositorio_jornadas` mostrado pelo usuário. O log de push e a ref local origin/main correspondem ao primeiro vínculo; não comprovam envio ao destino esperado. A criação também acrescentava o nome do projeto à pasta escolhida.

**Decisão**: usar a pasta escolhida diretamente e integrar no disco somente `.jornada` da raiz. Conservar a árvore completa nos objetos e commits para não excluir arquivos remotos anteriores; índice Git formato 3 marca os demais caminhos com `SkipWorktree`. A inspeção compara entradas/estágios/modos/hashes do índice com HEAD sem varrer o worktree omitido. [Contrato oficial do índice go-git](https://pkg.go.dev/github.com/go-git/go-git/v5/plumbing/format/index#Entry). Arquivos locais não elegíveis ficam intocados, inclusive em colisões com o remoto. Journals novos incluem somente mutações de jornadas; recovery continua reconhecendo journals anteriores.

**Destino e confirmação**: filtro invalida a seleção; antes de vincular, exigir digitação de owner/name completo. Registrar ID numérico e hash do destino, sem caminhos privados ou credenciais. Conclusão de publicação mostra owner/name, branch e link do commit confirmado. Remoto vazio deve retornar ausência, nunca aproveitar ref de cache. Preservar o vínculo incompatível com diagnóstico; não reescrever origin nem limpar o histórico existente.

**Alternativas rejeitadas**: checkout completo polui a pasta; remover arquivos do commit remoto perderia conteúdo preexistente; limpar automaticamente a pasta antiga não distingue arquivos originais de cópias integradas e pode causar perda. Usuário deve poder continuar pelos pacotes em outra pasta, preservando a tentativa anterior. Salvar continua local e conectar não autoriza publicação automática.

## Configuração externa e dependências operacionais

**Revisão de UX / 1.0.0.14**: a digitação adicional introduzida na 1.0.0.13 foi rejeitada no teste manual de 2026-10-25. O seletor passa o destino escolhido por valor, sem segundo formulário, e abre a revisão após o vínculo. Isso não autoriza push automático: a publicação ainda exige confirmar os arquivos e a mensagem. A tentativa real 1.0.0.13 falhou antes da transação de arquivos/push, com erro sem código; a causa exata não pode ser reconstruída apenas desse log. Jornadas locais e a cópia anterior do destino têm homônimos com bytes distintos; o novo diagnóstico testa colisão com `error_code=collision`, fase `workspace_integration` e correlação, sem atribuir esse resultado simulado à tentativa real.

O proprietário já forneceu o cadastro externo com os dois atributos necessários. O build lê `JORNADA_GITHUB_PROPERTIES`, seleciona apenas `CLI_ID_GITHUB_OAUTH` e `SECRET_GITHUB_OAUTH` e gera configuração Go temporária ignorada pelo Git; `build.ps1` remove o gerado ao finalizar. O par direto de ambiente continua aceito para automação. Não registrar valores nem a localização privada do arquivo em artefatos públicos.

A presença dos atributos é pré-requisito de build, não teste de autorização: callback cadastrado, permissões de organização/branch, fluxo no navegador e publicação ainda exigem homologação real. Checks globais, revisão/hash da build e resultados observados ficam em quickstart.md. Sem cadastro, o editor local continua disponível e a conexão informa indisponibilidade.
