# Contract: GitHub desktop

**Cobertura**: spec 037 / FR-001–004 e FR-015–018.

## Autorização

`Connect(ctx)` cria tentativa única, listener `127.0.0.1:0` e rota `/oauth/github/callback`, antes de `fyne.App.OpenURL`. State/verifier criptográficos ficam em memória; usar PKCE S256 e limite de 10 minutos. Nova tentativa invalida anterior. Validar GET, Host/porta, caminho, state em tempo constante e uso único antes da troca de código. Token exchange ocorre no backend Go do desktop, não em JavaScript.

Callback responde apenas mensagem para voltar ao Jornada Studio, com no-store/no-referrer e sem dados pessoais. Após identidade e credencial persistidas, `fyne.Do` exibe a janela existente (`Show`/`RequestFocus`) e abre seletor. Se Windows impedir foco automático, o estado já estará pronto ao retornar. Não abrir processo novo ou registrar protocolo do sistema. Cancelar encerra listener e invalida sessão.

Refresh, quando fornecido, é serializado por conta e atualiza o cofre. 401 definitivo/falha de refresh exige reconexão. Renovar token não autoriza nova publicação. Ausência de refresh é suportada com nova autorização quando necessária.

## REST

Base `https://api.github.com`; Bearer token só em memória; Accept application/vnd.github+json, User-Agent JornadaStudio e X-GitHub-Api-Version 2022-11-28. TLS validado; redirecionamentos não podem repassar Authorization a outro host.

| Requisição | Uso |
|---|---|
| `GET /user` | id/login/name, conferidos após cada autorização. |
| `GET /user/repos?visibility=all&affiliation=owner,collaborator,organization_member&per_page=100&sort=full_name` | id, owner.login, name, default_branch, permissions.push, archived, disabled. |
| `GET /repos/{owner}/{repo}` | Revalidar ID, destino, branch e escrita antes de vincular/publicar. |

Seguir Link rel=next somente no host/esquema/endpoint esperado. Deduplicar por ID; busca parcial nunca se apresenta como completa. Busca global carrega páginas restantes; erro repete só a página afetada, preservando itens anteriores.

401 pede reconexão; 403 diferencia política/SSO/limite quando possível; 404 informa inacessível; 429/limites respeitam Retry-After/reset, sem repetição imediata. Transporte indisponível não apaga vínculo. Permissão de escrita não garante aceitação pelas políticas de branch.

## Interfaces internas e erros

No consumidor application, separar interfaces de autorização (`Connect`, `Disconnect`), consulta (`CurrentUser`, `ListRepositories`, `GetRepository`), credencial (`Load`, `Save`, `Delete`) e preferências (`Load`, `Save`). Rede recebe context.Context. Adaptadores retornam structs concretas. UI recebe Account/Repository/resultados sem tokens.

Erros tipados: canceled, authorizationExpired, authorizationDenied, reconnectRequired, permissionDenied, rateLimited, networkUnavailable, credentialStoreUnavailable, invalidCallback, configurationMissing. Texto português acionável e causa sanitizada; não registrar corpos brutos OAuth.

## Distribuição

O build Windows também aceita `JORNADA_GITHUB_PROPERTIES` apontando a um arquivo externo com `CLI_ID_GITHUB_OAUTH` e `SECRET_GITHUB_OAUTH`. Valores não são impressos, não são copiados para preferências e o Go gerado ignorado é removido ao finalizar. Variáveis diretas completas prevalecem; pares incompletos são recusados.

## Diagnóstico seguro — FR-019

Reutilizar o Log de aplicação JSONL e correlacionar toda ação por `operation_id`, independente dos segredos OAuth. Registrar listener pronto (host literal, porta e caminho sem query), navegador aberto, callback recebido/validado/rejeitado/repetido, troca OAuth, identidade, cofre e consultas REST. HTTP registra status, duração, formato de resposta, request ID validado, limites numéricos e flags SSO; DNS/TCP/TLS registram fases, quantidades, versão/cipher e classes de erro, sem nomes privados ou dados de certificados.

Não serializar `Error()` de causas livres, request/response bodies, headers de autorização ou URLs completas. Código de erro OAuth é aceito somente por allowlist; error_description é sempre omitida. Callback informa apenas flags de validação/presença, nunca state/code/verifier/challenge. Nenhum valor de Client ID/Secret é logado. Se o navegador não retornar, registrar cancelamento/timeout e ausência do callback; o aplicativo não inspeciona páginas externas.

Implementar build que recebe `JORNADA_GITHUB_CLIENT_ID`/`JORNADA_GITHUB_CLIENT_SECRET` do processo e gera configuração Go temporária ignorada pelo Git, sem registrar valores no inventário ou log. Exigir par completo e limpar arquivo temporário ao finalizar. Credencial pública distribuída não é token pessoal nem garantia de confidencialidade.

Build sem cadastro desabilita conectar com explicação, preservando abertura e save. Proprietário fornece cadastro real antes da homologação; não inventar client_id/homepage, usar app de terceiros ou embutir token de desenvolvimento.
