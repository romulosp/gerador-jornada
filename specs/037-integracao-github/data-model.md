# Data Model: Conta, vínculo e publicação

Modelo efetivamente implementado em `internal/application/sync_model.go` e no journal do adaptador Git. Não há campo novo em Project/Document/ZIP. Nomes da coluna de persistência correspondem ao JSON atual; tipos de revisão são valores transitórios revalidados, não novos arquivos de negócio.

| Tipo real | Campos principais | Persistência e invariantes |
|---|---|---|
| GitHubAccount | ID int64, Login, Name | `id`, `login`, `name`; identidade por ID estável, uma conta ativa. |
| GitHubCredential | AccessToken, RefreshToken, ExpiresAt, Scopes | Somente no cofre, dentro de protectedRecord `{version: 1, token: ...}`. Vencimento zero indica prazo não informado. |
| GitHubAccountPreference | Version, Account, CredentialKey | `github.json`: `version`, `account`, `credentialKey`; sem token ou connectedAt. |
| GitHubRepository | ID, Owner, Name, Branch, Visibility, Writable, Archived, Disabled | `id`, `owner`, `name`, `branch`, `visibility`, `writable`, `archived`, `disabled`; seleção exige escrita e destino ativo. |
| GitWorkspaceBinding | Version, Root, Repository, Pending, Published, PublishedAt | `version`, `canonicalRoot`, `repository`, `pendingOID`, `lastPublishedOID`, `lastPublishedAt`; branch está dentro de Repository. |
| GitChange | Path, Kind, Digest | Projeção em memória; Kind é inclusão/alteração/exclusão. Digest é SHA-256 dos bytes de trabalho, vazio para exclusão. |
| GitPublicationReview | AccountID, Root, Repository, Head, Remote, Changes, RemoteChanges, Pending | Valor transitório comparado novamente antes de publicar; AccountID fixa a conta revisada; Head/Remote vazios representam referências sem commit. |

**Diagnóstico**: `operation_id` é gerado independentemente de OAuth state/PKCE e propagado por contexto. Eventos possuem sequência, etapa, duração, classificação de falha e metadados técnicos permitidos. O logger não persiste credenciais, parâmetros sensíveis, mensagens livres de causa, corpos HTTP ou conteúdo de projeto; falhas preservam somente tipos/classes/códigos seguros. O mesmo coletor JSONL local já existente recebe os eventos `github.*`.
| journal (adaptador) | LockID, Version, Binding, From, To, Initial, Stage, Files, IndexBefore, IndexAfter | `lockId`, `version`, `binding`, `from`, `to`, `initial`, `stage`, `files`, `indexBefore`, `indexAfter`; intenção anterior às mutações e marcador de propriedade de lock. |
| mutation (adaptador) | Path, Before, After | `path`, `before`, `after`; digests vazios representam ausência. Pré/pós-imagens ficam em arquivos privados separados. |

## Persistência

- Raiz: `JORNADA_CONFIG_DIR` quando definido; senão `os.UserConfigDir()/JornadaStudio`. Conta em `github.json`, vínculo em `git-workspaces/<SHA256(root-normalizada)>.json`, journal/backup em `git-operations/<SHA256(root-normalizada)>/`; uma operação por raiz. Trava privada em `git-locks/<hash>.lock`.
- JSON de preferências usa `version: 1`, lowerCamelCase e datas UTC RFC3339. Corrupção ou versão incompatível é recusada sem substituir o original. Cofre também usa wrapper versão 1 e rejeita token vazio ou blob incompatível.
- Canonicalizar raiz e links existentes; no Windows comparar a identidade da pasta sem distinção de caixa. Arquivos com symlink/junction e caminhos que escapem da raiz não são materializados.
- Cofre genérico local: `JornadaStudio/github/<hash-da-configuração>/<accountID>`. Credencial de usuário nunca vai ao JSON comum, `.jornada`, origin ou log. Namespace dos testes é descartável e separado do real.
- Stage inicial é um diretório temporário irmão da raiz, permitindo instalar `.git` no mesmo volume; é validado antes de remover. Pré/pós-imagens ficam no namespace de recuperação com nomes hash e sufixos `.before`/`.after`.
- O índice anterior/futuro é guardado no journal para distinguir escrita própria de alteração externa. OID planejado é To; o estado de avanço é inferido das referências, do vínculo, dos marcadores e dos bytes conhecidos, em vez de um enum phase persistido.
- A lista de repositórios permanece em memória. Backups da transação são removidos após conclusão; conteúdo divergente de escrita externa bloqueia recuperação automática e preserva localização da cópia.

## Revisão e primeiro commit

A revisão não persiste operationID, author, message ou snapshotDigest: carrega o destino, OIDs, lista de mudanças com seus digests e os caminhos da árvore final que podem exigir atualização local. Mensagem/autoria são fornecidas somente ao confirmar. O adaptador recalcula o plano e compara o valor completo antes de alterar HEAD; revisão desatualizada é recusada.

Head/Remote vazios não são um hash Git nulo a consultar. Com ambos vazios, a base é árvore vazia; o primeiro commit usa zero pais, inclui todo o conjunto elegível e cria a ref da branch padrão. Se apenas o remoto possui commit, ele fornece a base/pai sem substituir jornadas locais incompatíveis. Retentativa usa Pending conhecido; não recria o commit por falha de rede.

## Estados observáveis

**Conta**: desconectada → autorizando → conectada; cancelamento preserva a conta anterior. Token ausente/revogado exige reconexão, refresh é serializado e a remoção do cofre precisa concluir antes de declarar logout.

**Workspace**: sem vínculo → selecionando destino → integrando → vinculado. Vínculo é gravado somente depois de metadados/arquivos consistentes; erro parcial deixa intenção recuperável. Estado local inválido não é tratado como simples ausência.

**Publicação**: alterações locais → revisão → commit local pendente → envio/verificação → última publicação confirmada. Timeout após envio é resultado incerto até fetch comprovar ancestralidade. Edições posteriores aguardam a próxima publicação quando há commit pendente.

Esses nomes descrevem comportamento da aplicação; não são enums de estado salvos no modelo atual. Não existem lastRemoteOID/lastVerifiedAt persistidos: a data da última publicação é fato histórico, e a UI separa consulta remota explícita de informação local.

## Invariantes

Revisão 2026-10-25: nenhum campo de domínio ou formato `.jornada` muda. A árvore Git preserva os arquivos remotos; somente `.jornada` da raiz gera `mutation` no journal e entra em `RemoteChanges`. Entradas restantes usam skip-worktree no índice formato 3. O destino confirmado é a tupla ID/owner/name/branch, revalidada antes do vínculo e da publicação. Remoto vazio não é confirmado a partir de ref de cache.

Uma mutação por raiz no namespace da aplicação; índices/refs e arquivos são revalidados antes de aplicar plano. SHA-256 de arquivo não é OID Git. lastPublishedOID só avança após confirmação remota. Commits locais desconhecidos não são enviados automaticamente. Nenhum valor de erro/UI contém tokens, código OAuth ou verifier.
