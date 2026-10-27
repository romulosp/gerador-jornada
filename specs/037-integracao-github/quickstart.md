# Quickstart: Validar integração GitHub

## Estado

**pronto_para_teste_manual — correção 1.0.0.14**. O reteste de 1.0.0.13 mostrou confusão com a redigitação do nome e nova falha antes do push. A 1.0.0.14 remove o segundo formulário: selecionar o destino, conectar e revisar a publicação. Jornadas homônimas diferentes continuam protegidas contra sobrescrita e agora recebem classificação de erro e fase nos logs. A captura do usuário identifica `romulosp/gerador-jornada`; para o destino esperado no teste anterior, selecionar `romulosp/repositorio_jornadas`. Não atribuir a divergência ao usuário nem inferir a causa exata do erro genérico antigo. A aprovação humana continua pendente.

## Executável e teste imediato

- Executável: `apps/gerador-jornada/dist/windows-amd64-1.0.0.14/JornadaStudio.exe`.
- Versão: **1.0.0.14**; pasta separada preserva o EXE 1.0.0.13 aberto pelo usuário. O ZIP 1.0.0.12 é histórico e não contém estas correções.
- SHA-256 do EXE: `e909ca46e629e175906b1f4eb7a08a630adcb32aa2fdb064c97af4e3200617d6`.
- Abrir/criar projeto local, clicar **Conectar GitHub**, autorizar no navegador e selecionar repositório com escrita. Conferir o nome na linha Destino selecionado e clicar **Conectar repositório**; não haverá segunda digitação. Após sucesso, revisar os arquivos e clicar **Publicar no GitHub** para confirmar commit e push. A conclusão apresenta destino, branch e link do commit confirmado.
- Build consumiu o arquivo externo autorizado pelo usuário via `JORNADA_GITHUB_PROPERTIES`; os atributos foram encontrados e a configuração Go temporária foi removida. Nenhum valor foi copiado para documentação ou código versionável.

## Evidências técnicas executadas — 2026-10-24

### Simplificação do fluxo — 2026-10-25 / 1.0.0.14

O teste de UI reproduziu a segunda entrada de texto antes da correção. Depois, comprovou seleção do segundo repositório após filtro, ID entregue por valor, ausência de formulário adicional e cancelamento sem conectar. O teste de colisão reproduziu o erro sem código e depois comprovou `error_code=collision`, `checkpoint=workspace_integration`, correlação e preservação dos bytes/ausência de `.git` após o bloqueio. `scripts/check.ps1` passou com gofmt, vet, Go e verificações Python. O EXE configurado foi compilado em pasta versionada sem encerrar o aplicativo antigo. Reteste de login/vínculo/push no GitHub real permanece pendente.

O smoke nativo 1.0.0.14 passou com configurações/cache isolados, PATH sem Git e OAuth configurado. O log não contém os valores reais de Client ID/Secret; o arquivo Go temporário foi removido. Captura: `apps/gerador-jornada/build/037-editor-native-1.0.0.14.png`. PE AMD64, licenças e inventário de hashes verificados. Nenhuma escrita ou push foi realizado nos workspaces/repositórios reais durante esta correção.

### Correção e reteste técnico — 2026-10-25 / 1.0.0.13

Antes da correção, os testes reproduziram materialização de `.agents`, colisão indevida de README, confirmação de remoto vazio com ref em cache e criação na pasta homônima adicional. O teste de UI também falhou por ausência da conferência obrigatória do nome completo. Após corrigir, os testes focados e `scripts/check.ps1` passaram, incluindo `go test ./...`, vet, gofmt, skins/menu e verificações Python de build/configuração.

O harness com dois repositórios Git reais temporários comprovou que só o segundo destino vazio recebe os dois pacotes e que a ref do primeiro permanece inalterada. Outros testes preservam arquivos não elegíveis no remoto, não materializam código-fonte e não substituem README local. O teste de navegação agora publica, integra os pacotes em outra pasta e reabre o ZIP remoto com seus vínculos. O teste compilado rodou sem git.exe no PATH: 20 testes principais aprovados; evidência em `apps/gerador-jornada/build/037-git-no-git-path-1.0.0.13.log`. Esses resultados não substituem reteste no GitHub real.

O smoke do EXE configurado passou em processo nativo com APPDATA/LOCALAPPDATA/configuração isolados e PATH sem Git. Salvou um pacote, abriu-o sem arquivos auxiliares e capturou `apps/gerador-jornada/build/037-editor-native-1.0.0.13.png`; a janela mostra 1.0.0.13. O log confirmou configuração OAuth disponível e foi comparado aos valores reais de Client ID/Secret: zero ocorrências. O arquivo de configuração Go gerado foi removido após o build. Cabeçalho PE AMD64 e inventário SHA256SUMS conferidos; nenhum novo ZIP foi gerado nesta correção.

### Histórico da revisão 1.0.0.12

`scripts/check.ps1` passou: gofmt, skins/menu estáticos, go vet, empacotamento, configuração OAuth e go test ./.... O teste nativo de WinCred gravou/leu/removeu somente credencial sintética em namespace isolado. A suíte Git compilada rodou com PATH limitado ao System32, sem git.exe disponível, aprovando 29 casos/subcasos; o transporte de fixtures usa servidor Go incorporado. Isso comprova independência do executável Git nesse harness, sem afirmar homologação de máquina limpa completa.

Foram exercitados primeiro commit em remoto vazio (zero pais), preservação de bytes ZIP e navegação entre duas jornadas publicadas, paginação REST de 101 repositórios com deduplicação, retorno OAuth loopback com PKCE, callback inválido/repetido/cancelado, mudança de conta após revisão, staging alheio, conflitos, integração disjunta, push incerto e retentativa por OID. O harness encerrou um subprocesso real antes dos defers e recuperou o commit e seu lock próprio; testes também induziram falhas em intenção/arquivo/índice/referência/vínculo e preservaram edição externa durante rollback.

Os testes de diagnóstico preservaram as etapas e códigos seguros do provedor sem registrar valores sintéticos de Client ID/Secret, tokens, code/state/verifier/challenge ou descrição maliciosa. Evidências locais: `apps/gerador-jornada/build/037-git-no-git-path.log` e `apps/gerador-jornada/build/037-editor-native.png`. A captura nativa final usou configurações e cache isolados, confirmou a janela com o botão Conectar GitHub e a configuração OAuth incorporada, sem Git no PATH. Os logs dessa execução foram verificados contra os valores reais de Client ID/Secret, sem ocorrências; os eventos seguros estão em `apps/gerador-jornada/build/037-native-log-safe.json`. Essa prova não substitui o login/push reais nem as capturas de todos os estados/DPI.

## Investigar falha no login ou callback

Usar **Log de aplicação** no topo da janela. Cada erro GitHub apresenta um identificador: filtrar `operation_id` nos eventos `github.*`. O JSONL local fica em `%LOCALAPPDATA%\GeradorJornada\logs` e já possui rotação/retenção. Não enviar o arquivo externo de credenciais para investigar falhas.

Sequência esperada: `oauth.authorize.start` → `oauth.listener.ready` → `oauth.browser.open.complete` → `oauth.callback.received` → `oauth.callback.validated` → `oauth.exchange.start` → `http.response` → `oauth.token.received` → `connection.identity.complete` → `credentials.saved` → `connection.complete`. Eventos possuem o prefixo `github.`. Conferir flags de validação, porta/caminho esperados, duração, HTTP/request ID, erro conhecido OAuth e classes DNS/TLS/Win32.

Sem `callback.received`, o aplicativo ainda não recebeu retorno: verificar a página externa e o redirect cadastrado `http://127.0.0.1/oauth/github/callback`; aguardar conclusão ou cancelar. O programa não inspeciona o conteúdo da página GitHub antes do callback. Client ID/Secret, tokens, código, state, PKCE, URLs com query, corpos livres e conteúdo/mensagem de projeto não são registrados.

## Preparação

Usar Windows x64, toolchains locais e instruções do README da aplicação. Cadastro OAuth real pelo proprietário, callback conforme [contrato desktop](contracts/github-desktop.md); valores de build fornecidos por ambiente do processo, sem escrever credenciais no repositório ou argumentos de comando.

Preparar cópias descartáveis de workspace com duas jornadas interligadas, backups/exportações/temporários e documento legado; conta de teste com repositórios vazio, com conteúdo, sem escrita e com proteção de branch. Paginação além de 100 itens deve ser coberta por fixture HTTP; confirmar também paginação real quando houver esse volume autorizado. Homologar distribuição em computador sem Git instalado.

## Reexecutar os checks

Com o código atual e as ferramentas locais, PowerShell na raiz:

```powershell
Set-Location apps/gerador-jornada
. ./scripts/env.ps1
go test ./internal/application/... ./internal/adapter/github/... ./internal/adapter/gitworkspace/... ./internal/adapter/userconfig/... ./internal/adapter/credential/...
./scripts/check.ps1
python -X utf8 scripts/build.py --target windows --arch amd64
```

Os checks acima já passaram nesta revisão. O empacotador entrega `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`; toda alteração posterior exige registrar a nova revisão/hash e repetir os cenários afetados.

Testes devem incluir HTTP OAuth/REST local, callback inválido/repetido/cancelado, refresh, cofre em namespace isolado e repositórios reais temporários via go-git. Cobrir staging alheio, árvores divergentes, byte-preserving ZIP, colisões Windows e interrupção em cada fronteira do journal. Não usar rede externa por padrão nem credenciais reais nos fixtures.

## Roteiro manual

### Reteste prioritário de 1.0.0.14

Abrir o EXE da pasta versionada e confirmar 1.0.0.14. Usar a pasta local `teste`, selecionar o destino completo pretendido e clicar Conectar repositório. Não é necessário redigitar o nome. Após conexão, o aplicativo abre a revisão; confirmar Publicar no GitHub e conferir o link do commit. Se houver colisão, conferir o arquivo indicado e preservar os pacotes; usar o destino vazio correto em vez de substituir arquivos no repositório de código-fonte. Cancelar deve encerrar o seletor sem conexão ou publicação. Salvar continua disponível localmente.

### Histórico do reteste 1.0.0.13 (substituído pelo fluxo acima)

1. Preservar a pasta antiga e seu histórico. Copiar **todos** os `.jornada` da pasta aninhada para uma pasta limpa, mantendo os nomes exatos. Não copiar `.git`, código-fonte ou demais pastas. Abrir a jornada principal dessa pasta no executável novo. Evitar criar uma jornada homônima antes da cópia.
2. Selecionar **romulosp/repositorio_jornadas** e digitar esse nome completo na confirmação. Se uma pasta ainda vinculada a outro destino for usada, esperar bloqueio identificando ambos os destinos; não alterar origin nem apagar `.git` para contornar o bloqueio.
3. Conferir que a pasta contém os pacotes e somente os metadados Git locais necessários; nenhuma pasta do software é baixada. Publicar, conferir os nomes e confirmar a mensagem de commit.
4. Na conclusão, abrir **Ver commit no GitHub**. O link deve apontar a `romulosp/repositorio_jornadas`, com os `.jornada` na branch informada. Reabrir a jornada e verificar telas, imagens, propriedades e navegação entre pacotes. Salvar continua local; o push depende de Publicar.
5. Separadamente, criar um novo projeto em pasta já chamada `Meu projeto`; esperar o arquivo diretamente nela, sem subpasta homônima. Repetir integração com remoto contendo README/fonte e verificar que nada disso é copiado nem excluído do remoto.

Não houve limpeza da pasta real nem novo push pelo agente durante a correção. A pasta/histórico da tentativa anterior e o vínculo antigo foram preservados. Esta revisão não remove automaticamente arquivos previamente integrados: não é possível distinguir com segurança cópias e arquivos locais preexistentes.

| Caso | Ação | Resultado esperado | Requisitos |
|---|---|---|---|
| M01 | Conectar e autorizar no navegador. | Retorno à mesma janela, conta e seletor corretos. | FR-001/002/004 |
| M02 | Negar/cancelar e simular retorno inválido/duplicado no harness. | Sem vínculo indevido; projeto local preservado. | FR-002/018 |
| M03 | Percorrer/buscar além da primeira página; privado/organização e sem escrita. | Lista completa autorizada; destino inelegível bloqueado com motivo. | FR-003 |
| M04 | Vincular pasta sem Git a remoto vazio. | `.git` na raiz, sem commit/push durante vínculo. | FR-005/006/008 |
| M05 | Integrar remoto com README/fonte; repetir com colisão de jornada/histórico incompatível. | Só jornadas da raiz materializadas; demais arquivos preservados no remoto, sem cópia local; colisão bloqueia. | FR-007/018/020 |
| M06 | Salvar/publicar duas jornadas interligadas e reabrir. | Conjunto no mesmo commit; bytes/vínculos intactos; alheios excluídos. | FR-008–011 |
| M07 | Publicar edição suja e legado com confirmação de migração; cancelar/repetir. | Callback respeita save, cancelamento e erro; sem envio prematuro. | FR-009 |
| M08 | Sem mudanças; staging alheio; exclusão referenciada. | Sem commit vazio; staging bloqueado; referência afetada informada. | FR-010–012 |
| M09 | Falhar após commit e após servidor receber push; reiniciar e repetir. | Mesmo OID entregue uma vez e pendência recuperada. | FR-013/016 |
| M10 | Avançar remoto em arquivo distinto e depois no mesmo `.jornada`. | Integração disjunta preservada; conflito bloqueado sem force. | FR-014 |
| M11 | Proteger branch/revogar token/desligar rede; editar e salvar. | Diagnóstico específico, commits preservados e edição local disponível. | FR-014/018 |
| M12 | Reiniciar, trocar conta e desconectar; inspecionar cofre/disco/logs. | Destino revalidado, credencial removida, jornadas/histórico preservados. | FR-015–017 |
| M13 | Interromper integração/commit no harness; alterar externamente durante recuperação. | Journal recupera ou bloqueia sem sobrescrever mudança externa. | FR-007/013/018 |
| M14 | Conferir controles em DPI/temas previstos. | Layout, TAB/foco e editor/prévia sem regressão. | DESIGN.md; FR-015 |
| M15 | Criar projeto em pasta com mesmo nome; salvar e reabrir só o pacote. | Nenhuma subpasta adicional; projeto e recursos restaurados pelo `.jornada`. | FR-020/021; SC-009 |
| M16 | Buscar e selecionar o segundo repositório, conectar e publicar sem redigitar o nome. | Só o destino selecionado recebe o commit; conexão abre revisão; conclusão abre o commit correto; remoto vazio não comprova entrega por cache. | FR-022; SC-010 |

M01/M03/M04/M06/M11/M12 precisam de GitHub real para homologação integral; simulação não substitui evidência real. M09/M10/M13 também exigem harness determinístico para controlar falhas. Este planejamento não executa nenhuma escrita nos repositórios de teste.

## Evidências

| Revisão/build/SHA-256 | Executor/data/ambiente | Casos | Observado | Evidência | Situação |
|---|---|---|---|---|---|
| 1.0.0.12 / 61f9758d… | Agente / 2026-10-24 / Windows x64 | Provas locais anteriores | Checks e build aprovados, insuficientes para o fluxo real | Logs/captura em build/ | Revisão posteriormente reprovada |
| 1.0.0.12 | Usuário / 2026-10-25 / Windows 10 | Pasta, destino e publicação | Pasta duplicada, fonte no workspace e destino esperado vazio; autenticação retornou, mas vínculo local aponta outro repositório | Relato, listagem, captura e log local | teste_manual_reprovado |
| 1.0.0.13 / a89cfa4c... | Agente / 2026-10-25 / Windows x64 | Regressões de pasta/conteúdo/destino/cache, checks, build e smoke nativo | Testes focados, check.ps1 e abertura do EXE aprovados | Harness local, captura e build | pronto_para_teste_manual; reteste humano pendente |
| 1.0.0.13 / a89cfa4c… | Usuário / 2026-10-25 / Windows 10 | Seleção e vínculo | Redigitação confundiu o fluxo; vínculo falhou antes do push com erro sem código | Captura e log local | teste_manual_reprovado |
| 1.0.0.14 / e909ca46… | Agente / 2026-10-25 / Windows x64 | Seleção sem segunda entrada, cancelamento, colisão, checks e smoke | Todos os checks técnicos executados passaram | EXE/captura em pasta versionada | pronto_para_teste_manual; GitHub real pendente |

Preencher somente após execução. Build e checks aprovados permitem pronto_para_teste_manual; cenários aplicáveis aprovados permitem pronto_para_aprovacao. Aprovado exige manifestação humana sobre a revisão testada. Falhas retornam à correção e reteste.
