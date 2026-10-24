# Quickstart: Tela inicial (launcher) e Estrutura de Projeto

Guia de validação executável da feature 034. Complementa o quickstart da coordenadora (001) sem duplicar o roteiro geral.

## Pré-requisitos

- Windows 10/11 x64; ferramentas de build do Go disponíveis no ambiente de desenvolvimento (o usuário final não precisa de Go).
- Repositório do editor em `apps/gerador-jornada`.

## Build e execução

```powershell
cd apps/gerador-jornada
./scripts/check.ps1
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
./dist/windows-amd64/JornadaStudio.exe
```

Para gerar jornadas de teste em uma pasta (sem abrir a janela):

```powershell
go run ./cmd/gerador-jornada -save-example C:\tmp\ws\principal.jornada
go run ./cmd/gerador-jornada -save-compounds C:\tmp\ws\menu.jornada
```

## Cenários de validação

### C1 — Tela inicial (projetos recentes e novo projeto)

1. Abrir o programa sem `-open` e sem documento.
2. **Esperado**: aparece a tela inicial, sem abrir o exemplo automaticamente.
3. Criar "Novo projeto em branco" (escolher nome e confirmar pasta): nasce um documento só com o fundo padrão, sem conteúdo de exemplo.
4. Fechar e reabrir o EXE, criar "Novo projeto com exemplo": nasce o cenário de carregamento atual (exemplo editável), com o nome escolhido.
5. Fechar e reabrir: os projetos recentes aparecem e podem ser reabertos; uma entrada movida/excluída é indicada como indisponível.

### C2 — Estrutura projeto (árvore de jornadas)

1. Abrir (ou criar) uma jornada dentro de `C:\tmp\ws\`.
2. No painel esquerdo, abrir a aba "Estrutura projeto".
3. **Esperado**: a árvore mostra a raiz (nome da pasta) e as `.jornada` presentes, ex.: `ws` → `principal` e `menu`.
4. Menu de três pontos → "Criar jornada": nasce uma `.jornada` em branco e ela aparece na árvore.
5. Menu de três pontos → "Importar jornada existente": uma `.jornada` de fora é copiada para a pasta (origem intacta) e aparece na árvore; um arquivo inválido é recusado com diagnóstico.

### C3 — Criação rápida e navegação entre jornadas

1. Em uma jornada com o Menu aplicação, executar a prévia e clicar numa folha sem destino.
2. **Esperado**: mensagem de seleção simulada + sugestão imediata de criar uma jornada com o nome da folha.
3. Aceitar: a jornada aparece na árvore da Estrutura projeto e a folha passa a ter esse destino.
4. Executar a jornada inicial e clicar na mesma folha.
5. **Esperado**: a prévia mostra a tela inicial da jornada criada; "Voltar à jornada anterior" retorna à jornada de origem.
6. Repetir a sugestão e cancelar: bytes da origem, arquivos e histórico não mudam. Remover um destino numa cópia de teste e conferir o diagnóstico; Prévia Tela Atual nunca segue destinos.

## Checks

```powershell
cd apps/gerador-jornada
go test ./...
```

## Critérios de aceite (resumo)

- `SC-001..SC-004` da [spec.md](./spec.md): criar projeto em ≤3 interações; 100% das `.jornada` válidas na árvore; sugestão imediata + navegação após criar; cancelamento sem estado parcial.

## Resultado técnico — 2026-10-24 / 1.0.0.4

Base Git: `2319beb05 + árvore de trabalho local (sem commit de entrega)`. Windows 10 x64; Go 1.27.2/Fyne 2.8.2. Executável: `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`, **33334784 bytes**, SHA-256 **`366f6a7276d781bb7193d15fde70dbd12f0a0f8ec7837744495e2655f47a9378`**. A cópia de compatibilidade `dist/JornadaStudio.exe` tem os mesmos bytes. Somente o EXE foi compilado para entrega; nenhum pacote ou outro alvo foi gerado.

| Cenário | Execução observada | Resultado |
| --- | --- | --- |
| C1 | `TestRecentsAndLauncher`: startup sem abrir sessão, criação real em branco/com exemplo, nome/fundo, recentes persistidos/reabertura e entrada removida; JSON corrompido localizado. Startup do EXE sem `-open`. | PASS técnico; [captura nativa](../../apps/gerador-jornada/build/034-launcher-native.png). Caminhos longos quebram linhas e a indisponibilidade continua legível. |
| C2 | `TestWorkspaceCreateImport` e `TestWorkspacePreviewConfirmation`: raiz não recursiva, colisões sem sobrescrita, criação/importação com bytes da origem intactos, arquivo corrompido recusado, atualização e seleção real do Tree. Renderer verificado para raiz e nome do arquivo. | PASS técnico; [árvore no EXE](../../apps/gerador-jornada/build/034-workspace-native.png). `TestWorkspaceHundredJourneys`: listagem/validação de 100 pacotes com fundo incorporado em **54,96 ms** nesta execução; não é garantia para qualquer volume/tamanho de recurso. |
| C3 | Clique nos controles legados reconstruídos: folha abre confirmação, Cancelar mantém bytes/árvore/histórico, Criar grava vínculo e arquivo, próximo clique navega, voltar restaura submenu; seleção da árvore reabre o destino. Testes do núcleo conferem StartScreenID, arquivo ausente, escopo Tela Atual e round-trip. Cancelamento de IO testa rollback e Undo/Redo. | PASS técnico; [prévia de retorno à origem](../../apps/gerador-jornada/build/034-preview-test.png), captura produzida pelo backend Fyne de teste. |
| Compatibilidade | `TestJourneyCatalogMigration`: TargetJourney Unicode sobrevive ao round-trip; 3/3 migra sem inventar rota; save publica 3/4 e guarda backup byte a byte; backup indisponível impede sobrescrita; campo ausente em 4 ou presente em versão histórica é recusado. | PASS técnico. |

Logs: [baseline](../../apps/gerador-jornada/build/034-baseline.log), [cenários detalhados](../../apps/gerador-jornada/build/034-scenarios.log), [check completo](../../apps/gerador-jornada/build/034-check.log), [build](../../apps/gerador-jornada/build/034-build.log) e [startup nativo](../../apps/gerador-jornada/build/034-native.log). `check.ps1` executou gofmt, os dois importadores `--check`, `go vet ./...`, três testes Python do empacotamento e `go test ./...`; todos passaram. Check opcional de toolchain Java não foi solicitado nesta feature. Hooks de extensão ausentes em `.specify/extensions.yml`.

**Correção identificada na evidência visual:** a primeira captura deixou a árvore vazia porque o Tree tentava percorrer a raiz implícita não configurada. Definir `Root = "workspace"` tornou a pasta e os filhos visíveis; a regressão passou a inspecionar os labels renderizados, além de ChildUIDs. Folhas exibem nome físico e nome do projeto em duas linhas para painéis estreitos.

**Limites e aceite:** C1/C2/C3 foram executados tecnicamente com fixtures e cliques programáticos, mais startup/capturas do EXE. Não houve aceite humano de usabilidade, medição das três interações por um autor, DPI 100/150/200%, leitor de tela ou homologação Office. Esses resultados continuam pendentes em 001/T156. Um arquivo de origem ainda não salvo recebe orientação para salvar antes de vincular; a confirmação informa o salvamento das alterações pendentes. Desfazer o vínculo conserva o arquivo criado. Mudanças externas na pasta são relidas por **⋮ → Atualizar**. Nenhum XML de movimentos ausente foi reconstruído por esta feature.

## Correção de encerramento — 1.0.0.5

**Relato:** criar gravava o projeto, mas fechava ao entrar no editor; escolher o recente também fechava. O JSONL fornecido tinha somente startup da 1.0.0.4. A reprodução em janela nativa confirmou panic de ponteiro nulo em `GetWin32Window`, chamado por `setDarkMode` quando `initTheme` reaplicava o tema com uma janela ainda não exibida. O arquivo do projeto não causava o encerramento.

**Correção:** `prepareApplicationTheme` antes de NewWindow, idempotência do tema, seleção inicial sem OnChanged e cor inicial do título obtida do tema. O callback manual continua ativo depois da inicialização. `editor_open_start`/`editor_open_ready` acrescentam diagnóstico da transição sem conteúdo ou caminho do documento. Formato continua **3/4**; não houve migração adicional.

**Verificação observada:** oito fluxos nativos passaram — criar em branco/reabrir em Claro, Escuro e Azul; criar exemplo/reabrir em Escuro. O teste usa driver gráfico real e controles renderizados, com uma configuração de recentes isolada e fixtures em build; não toca os projetos fornecidos pelo usuário. A regressão de unidade `TestLauncherEditorThemeOrder` observa a ordem de SetTheme e preserva a troca manual para Azul. Reabertura nativa final repetida depois da correção de cor do título. Checks completos e build passaram; smoke do EXE final sem argumentos e com `-open` aprovado. Homologação humana continua independente em 001/T156.

Logs: [panic reproduzido antes da correção](../../apps/gerador-jornada/build/035-repro.log), [reprodução corrigida](../../apps/gerador-jornada/build/035-repro-fixed.log), [oito fluxos nativos](../../apps/gerador-jornada/build/034-launcher-native-regression.log), [reabertura nativa final](../../apps/gerador-jornada/build/034-launcher-native-final.log), [checks](../../apps/gerador-jornada/build/034-launcher-fix-check.log), [build](../../apps/gerador-jornada/build/034-launcher-fix-build.log) e [startup do EXE](../../apps/gerador-jornada/build/034-launcher-exe-final.log). Capturas: [exemplo reaberto pelo launcher nativo](../../apps/gerador-jornada/build/034-launcher-fix/Exemplo/editor-recent-Escuro.png) e [editor no EXE final](../../apps/gerador-jornada/build/034-launcher-fix-editor-exe.png).

EXE **1.0.0.5**, **33336320 bytes**, SHA-256 **`2576917419e883744ff411ee362506935f808f6f9a01eb83242b70b563845568`**. Gerado somente `JornadaStudio.exe` Windows x64; cópia `dist/JornadaStudio.exe` idêntica. A revisão anterior e seu hash na seção acima são históricos.

### Repetir a regressão nativa

Em sessão gráfica Windows, a partir de `apps/gerador-jornada`:

```powershell
. ./scripts/env.ps1
$env:JORNADA_REPRO_DIR = Join-Path (Get-Location) 'build/launcher-smoke'
New-Item -ItemType Directory -Force -Path $env:JORNADA_REPRO_DIR | Out-Null
$env:JORNADA_CONFIG_DIR = Join-Path $env:JORNADA_REPRO_DIR 'config'
$env:JORNADA_REPRO_THEME = 'Escuro' # repetir com Claro e Azul
$env:JORNADA_REPRO_MODE = 'blank' # example para o outro modo
 go run -tags launchersmoke ./scripts/launcher_smoke
$env:JORNADA_REPRO_MODE = 'recent'
 go run -tags launchersmoke ./scripts/launcher_smoke
```

A ferramenta é fonte de teste, não distribuição adicional: abre o launcher, aciona os botões, preenche nome/pasta, cria ou reabre, verifica que o editor permanece ativo e captura a janela. `go test ./...` normal não executa esse pacote com janela; use a tag explicitamente na sessão gráfica.

## Opção de apresentação — 1.0.0.10

Na tela inicial, “Mostrar splash screen ao iniciar” aparece abaixo de “Abrir projeto existente”, desmarcado por padrão. `TestLauncherSplashPreferencePersistence` usa o Check real para marcar, reconstruir o launcher, restaurar a marcação e desmarcar; a escolha sobrevive ao fechamento e não cria documento. Testes de erro preservam configuração e revertem marcação não salva. O EXE passou por primeira abertura com splash, reabertura direta, preferência habilitada com splash e desabilitada sem splash. [Evidências e roteiro de aprovação](../036-splash-screen/quickstart.md#revisão-atual--10010). Aceite humano permanece em T025, sem encerrar gates anteriores.
