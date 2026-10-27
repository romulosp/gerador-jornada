# Implementation Plan: Conexão do workspace com GitHub

**Branch**: `037-integracao-github` (identificador documental; nenhuma branch criada) | **Date**: 2026-10-24 | **Spec**: [spec.md](spec.md)

**Input**: `/specs/037-integracao-github/spec.md`, FR-001 a FR-022. Revisão 1.0.0.14 remove a digitação extra conforme teste manual de 2026-10-25: `chooseGitRepository` entrega a seleção por callback, `bindGitRepository` conecta e abre a revisão sem push automático. Colisão gera código seguro e fase de vínculo no logger. Entrega em pasta versionada para preservar o EXE 1.0.0.13 em execução.

## Summary

A revisão de código 1.0.0.14 conecta um workspace local por OAuth App desktop com PKCE e retorno loopback à janela Fyne existente. Lista repositórios por REST e usa go-git incorporado para inicializar, integrar, commitar e publicar por HTTPS. Salvar permanece local; Publicar exige revisão e confirmação. A coordenação em `application.WorkspaceSync` revalida a conta/destino sem entregar tokens à UI.

A implementação possui testes técnicos locais. O cadastro OAuth foi fornecido externamente, mas presença de atributos não comprova autorização no navegador ou publicação no GitHub real. Build, checks globais e resultados manuais são registrados em [quickstart.md](quickstart.md); esta atualização documental não os declara concluídos. Decisões: [research.md](research.md); interface efetiva: [DESIGN.md](DESIGN.md).

## Technical Context

**Correção de 2026-10-25**: o launcher usa a pasta escolhida diretamente. A árvore/histórico Git continua completa para preservar conteúdo remoto, mas preflight e journal só materializam `.jornada` da raiz. O índice marca os demais caminhos como skip-worktree (formato 3), sem ativar sparse checkout externo. A inspeção compara índice com HEAD diretamente, detectando staging sem varrer o conteúdo omitido. Conteúdo não elegível já presente nunca é removido automaticamente. Seleção exige confirmar owner/name; sucesso exibe destino/branch/OID após consulta do servidor. Remoto vazio invalida refs em cache como prova de entrega. Testes e build 1.0.0.13 devem substituir a revisão reprovada para reteste.

**Language/Version**: Go 1.27.2, módulo único `gerador-jornada.local/editor`; software 1.0.0.14.

**Primary Dependencies**: Fyne v2.8.2, golang.org/x/sys e `github.com/go-git/go-git/v5 v5.19.3` (Apache-2.0), fixados em go.mod/go.sum. HTTP/OAuth/PKCE/JSON usam a biblioteca padrão. Não há dependência de git.exe.

**Storage**: `.jornada` inalterado; `.git` na raiz selecionada; conta/vínculo/journal na configuração privada; credencial pessoal no Windows Credential Manager em blob versionado.

**Testing**: testes Go com HTTP e repositórios bare temporários, primeiro commit sem HEAD, retentativa por OID, integração disjunta/conflito, interrupção em subprocesso real, navegação de pacotes ZIP, cofre Windows real com token sintético e callback de salvamento Fyne. Eles não substituem homologação browser/GitHub e DPI reais.

**Target Platform**: Windows 10/11 x64, sem Git instalado, navegador embutido ou backend próprio. Entrega corrente somente `JornadaStudio.exe`; não gerar arquivos de distribuição adicionais sem pedido explícito.

**Project Type**: desktop-app; monólito modular e injeção manual.

**Performance Goals**: rede fora da thread Fyne, resposta administrativa em até 250 ms no harness e cancelamento cooperativo. REST é paginado em blocos de até 100; a UI recebe a lista completa após a carga em background, sem botão incremental Carregar mais.

**Constraints**: uma conta, um repositório por workspace, branch padrão e GitHub.com HTTPS; sem force push ou merge textual de ZIP. Autorização limita-se a 10 minutos, REST a 30 s e operação Git a 15 minutos. Timeout não comprova rejeição do push. O diálogo de operação bloqueia interação com a sessão enquanto o trabalho está ativo; cancelamento/falha libera a sessão após recuperação ou diagnóstico do journal.

**Scale/Scope**: autor único, centenas de jornadas/repositórios paginados; LFS, submódulos, symlinks, linked worktrees, sparse checkout e históricos incompatíveis são recusados. O editor local continua independente de autorização/rede; a integração exige cadastro OAuth no build.

## Constitution Check

Avaliação arquitetural com constituição 1.3.0; conformidade funcional continua dependente das evidências.

| Princípio | Implementação e limite da verificação |
|---|---|
| I. Domínios | WorkspaceSync coordena portas de Git/autorização; filesystem, rede e Win32 permanecem em adaptadores. |
| II. Contratos | Documento e formato `.jornada` intactos; preferências e credenciais não integram o pacote. |
| III. Distribuição/offline | go-git incorporado, loopback temporário e rede opcional; distribuição real sem Git depende da build entregue e teste manual. |
| IV. Documentação | Novos pacotes e variáveis documentados em português; manter revisão nas mudanças posteriores. |
| V. Evidências | Testes locais cobrem persistência, bytes, navegação e recuperação; homologação real e aceite não são inferidos deles. |
| VI. Simplicidade | Um módulo e dependência Git justificada; sem banco, DI framework ou serviço próprio. |
| VII. Rastreabilidade | Artefatos descrevem caminhos consolidados; tarefas/evidências permanecem vinculadas aos requisitos 037. |

Não há alteração de princípio para acomodar a implementação. O listener temporário atende ao retorno desktop explicitamente solicitado.

## Project Structure

### Documentation (this feature)

```text
specs/037-integracao-github/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── DESIGN.md
├── quickstart.md
├── tasks.md
├── checklists/requirements.md
└── contracts/
    ├── github-desktop.md
    └── workspace-git.md
```

[tasks.md](tasks.md) contém o backlog e [quickstart.md](quickstart.md) o registro de revisão/execução. Criar documentos ou atualizar este plano não conclui tarefas técnicas/manuais.

### Source Code (repository root)

```text
apps/gerador-jornada/
├── cmd/gerador-jornada/main.go          # ativa a composição GitHub após construir a janela
├── internal/application/
│   ├── github_connection.go            # autorização, conta, refresh e logout
│   ├── workspace_sync.go               # WorkspaceSync: Bind, ReviewPublication e Publish
│   └── sync_model.go                   # projeções e portas de cofre/configuração
├── internal/adapter/
│   ├── github/{config,client,oauth}.go  # configuração pública, REST e loopback
│   ├── gitworkspace/                   # engine, tree, journal, gate Windows e status local
│   ├── userconfig/store.go             # JSON atômico e chaves por raiz
│   ├── credential/                     # WinCred versionado; stub fora do Windows
│   └── desktop/
│       ├── github_panel.go             # fluxo administrativo compacto e diálogos
│       ├── github_icon.go              # Octicons incorporado, licença MIT
│       └── files.go                    # callback de save/cancelamento/erro
├── internal/synctest/fixtures.go        # pacotes .jornada reais para integração local
└── scripts/configure_github.py         # configuração de build externa e limpeza do gerado
```

**Structure Decision**: consolidar os arquivos planejados em poucos adaptadores coesos, mantendo interfaces no consumidor. A UI usa WorkspaceSync para as ações remotas e não recebe token; o adaptador Git guarda a transação em journal. Não foram criados domínios artificiais nem outro módulo Go.

## Implementação e validação

1. Configuração usa `JORNADA_GITHUB_PROPERTIES` apontando para arquivo externo com `CLI_ID_GITHUB_OAUTH` e `SECRET_GITHUB_OAUTH`. O build também aceita o par direto de ambiente descrito no script, sem misturá-lo com par incompleto. `configure_github.py` gera `config_generated.go` ignorado pelo Git; `build.ps1` remove o temporário em finally. Valores e localização privada não são documentação pública.
2. Vínculo inspeciona raiz/origin/branch, prepara conteúdo fora do workspace e aplica arquivos/índice/referência com journal. Sem HEAD e remoto sem referência, ambos os OIDs são vazios: a primeira publicação usa árvore vazia e commit sem pais, antes de criar a ref da branch padrão.
3. Publicação compara revisão e arquivos novamente, inclui somente `.jornada` da raiz e guarda OID pendente. Integração conservadora usa árvores e commits com dois pais, preservando novidade remota disjunta sem mesclar ZIP. Interferência externa e conflito bloqueiam sem overwrite forçado.
4. Painel GitHub compacto mantém status truncado com detalhes no menu ⋮ e uma ação principal conforme o estado. A lista inteira é obtida em background antes da busca; toda página autorizada continua incluída, preservando FR-003.
5. Callback de save preserva confirmação de migração/cancelamento. Operação falha invalida recuperação cacheada e recupera ou bloqueia mutações locais antes de seguir. Documento ativo só é recarregado se a sessão estiver limpa e o pacote for íntegro.
6. Testes focados locais já exercitam esses mecanismos. Checks globais, build 1.0.0.12, execução no navegador/GitHub real e homologação visual devem ter seus resultados próprios no quickstart; não há aprovação humana implícita.

## Complexity Tracking

Journal e integração por árvores são necessários a FR-007/010/013/014/018; não criam um cliente Git geral. Branches e resolução de conflitos continuam fora do escopo. A lista completa em background e o painel compacto são ajustes de apresentação documentados, não redução do conjunto de repositórios ou dos dados preservados.
