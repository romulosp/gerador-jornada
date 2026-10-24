# Quickstart: Splash em janela com áudio

**Revisão vigente: [1.0.0.11 em janela](#revisão-vigente--10011-em-janela).** As seções 1.0.0.9/1.0.0.10 são histórico; o fullscreen foi substituído por solicitação posterior.

## Histórico técnico — 1.0.0.9 (substituído pela revisão atual abaixo)

Revisão técnica **1.0.0.9**, 2026-10-24, Windows 10 x64 / Go 1.27.2 / Fyne 2.8.2. Base Git 2319beb05 + alterações locais (sem commit de entrega). Entrega: `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`, **37851648 bytes**, SHA-256 **`af4400c379ad87096a693f30ee97275fa565f6ba0fa1d7196ae4ef388f4aaa03`**. A cópia `dist/JornadaStudio.exe` tem os mesmos bytes. Não foram gerados pacotes ou outros alvos.

## Resultado técnico histórico observado

| Cenário | Evidência | Resultado |
| --- | --- | --- |
| Barreira de inicialização | `startup.TestBothCompletionsRequired`: duas ordens de conclusão, fim sem início, erro, cancelamento e liberação única | PASS |
| Integridade de mídia | `builtinassets.TestSplashV2Integrity`: bytes do v2, hash, trilhas vide/soun, cache corrompido reparado e extração cancelada | PASS |
| Visual mínimo e autoria | `desktop.TestSplashOnlyVideoAndLauncherCredit`: splash contém somente a superfície; autoria literal no launcher | PASS |
| Recursos e paridade | `TestStartupResourcePreparation`, caches imutáveis e regressões existentes de fontes/alpha/skins/menu/PPT | PASS |
| Geometria nativa | `nativevideo.TestEntireFrameFits`: a própria função C++ em seis tamanhos; razão 16/9 convertida para 1280×720; pixels não quadrados 4/3; contain com limite 100% | PASS |
| Intermitência | Quatro processos consecutivos do EXE; capturas aos ~2 s e ~8 s de cada reprodução; conteúdo colorido significativo nas oito imagens; áudio selecionado e ordenação de eventos em cada sessão | PASS técnico observado, sem garantia para hardware/DPI não homologados |
| Fluxo completo | Ferramenta nativa: create e recent passaram; o editor permaneceu aberto após splash/launcher | PASS |
| Carregamento lento | Carga artificial de 14 s na ferramenta; tela liberada em **14.971 ms**, depois do vídeo e da carga. Carga rápida create: **330 ms**, liberação em **11.766 ms** | PASS |
| Cancelamento/falha | Fechar em 2 s e falha injetada: áudio/worker encerrados, sem liberar launcher | PASS |
| Utilitário -open | EXE com uma jornada temporária; apresentação completa e abertura posterior do editor | PASS técnico; captura específica abaixo |

O EXE final usa uma classe de janela própria para vídeo, com WM_PAINT/WM_SIZE. A configuração e o primeiro Play são serializados na thread do HWND após mídia pronta; callbacks livres não tocam a janela. A superfície STATIC anterior e a interpretação de 16/9 como pixels causaram regressões intermediárias identificadas pelo usuário. Eventos PLAY/END e processo com código zero foram insuficientes para provar imagem; a validação final inclui capturas reais repetidas.

## Logs e capturas

- [Checks completos](../../apps/gerador-jornada/build/036-check.log): gofmt, dois importadores --check, go vet, três testes Python de binários e go test ./....
- [Build Windows](../../apps/gerador-jornada/build/036-build.log), somente EXE.
- [Quatro inicializações verificadas](../../apps/gerador-jornada/build/036-repeated-startups.log), com [eventos da primeira](../../apps/gerador-jornada/build/036-startup-1.jsonl).
- [Imagem no início](../../apps/gerador-jornada/build/036-video-1-early.png) e [perto do fim](../../apps/gerador-jornada/build/036-video-1-late.png), em janela fullscreen. Captura da área da própria janela; marca de ativação do Windows pertence ao sistema, não ao app/vídeo.
- [Launcher final com autoria](../../apps/gerador-jornada/build/036-launcher-4.png), depois da transição.
- [Criar/reabrir/carga lenta/cancelar/falhar](../../apps/gerador-jornada/build/036-final-flow.log), com [editor após criar](../../apps/gerador-jornada/build/036-final-flow/splash-create.png).
- [Abertura explícita](../../apps/gerador-jornada/build/036-open-project.log) e [editor -open](../../apps/gerador-jornada/build/036-editor-native.png).

## Roteiro humano histórico — 1.0.0.9

1. Fechar instâncias antigas e abrir o EXE **1.0.0.9** do caminho acima.
2. Conferir vídeo v2 com áudio, tela cheia e quadro integral; sem autoria, progresso ou texto sobreposto. Margens pretas são intencionais; não ampliar além dos pixels originais.
3. Aguardar vídeo inteiro mesmo com recursos prontos, observar fade e launcher em janela normal.
4. Conferir o rodapé do launcher: **Desenvolvido por Rômulo Penha - Canais Parceiro - 03**, sem vírgula.
5. Criar projeto em branco/com exemplo e reabrir recente; o editor deve permanecer aberto.
6. Repetir algumas inicializações para procurar a falha intermitente; fechar durante uma apresentação e conferir que áudio para e launcher não aparece depois.
7. Homologar som percebido, proporção, fullscreen, transição e DPI 100/150/200% no computador de uso.

**Aceite humano pendente em T009**. As capturas foram revisadas e o player informou trilha de áudio selecionada, sem mute imposto pelo app. Volume/mudo/saída do Windows e percepção auditiva precisam ser conferidos pelo usuário.

## Repetir testes no desenvolvimento

```powershell
cd apps/gerador-jornada
./scripts/check.ps1
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
# Windows gráfico, Pillow disponível só no ambiente de desenvolvimento; toca áudio:
python -X utf8 scripts/test_splash_native.py --runs 4
# Modelo puro, sem janela:
go test ./internal/startup
```

Para a ferramenta de integração (fonte, sem distribuição adicional), carregar scripts/env.ps1, definir JORNADA_REPRO_DIR e JORNADA_CONFIG_DIR numa pasta sob build, e executar `go run -tags splashsmoke ./cmd/splash-smoke`. JORNADA_SPLASH_MODE aceita create/recent/slow/cancel/fail; JORNADA_SPLASH_DELAY_MS=14000 injeta carga lenta. A aplicação entregue não usa esse atraso artificial.

## Separação de modelos e responsabilidades

- [035/data-model](../035-inicializacao/data-model.md): regras puras em internal/startup.
- [036/data-model](data-model.md): estado de janela/player/cache nos adaptadores, sem contaminar Document.
- [036/plan](plan.md) e [DESIGN](DESIGN.md): integração e visual vigente. Formato de projeto continua **3/4**.

## Revisão atual — 1.0.0.10

**Histórico: substituída pela 1.0.0.11 em janela, abaixo.**

2026-10-24, Windows 10 x64, Go 1.27.2/Fyne 2.8.2, executor técnico: Codex. Revisão do checkout com alterações locais; aceite humano ainda pendente. **Esta seção substitui as instruções de 1.0.0.9 acima.**

Entrega somente **JornadaStudio.exe** Windows amd64, **37860864 bytes**, SHA-256 **`c71eb5d0a0c342e2576e0d748214ded693a7831bd51bfc819c7f979a3bedb7ad`**, em `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`; cópia `dist/JornadaStudio.exe` idêntica. Formato do projeto continua 3/4, sem alteração. Não gerar ZIP, instalador ou outro alvo sem solicitação explícita.

### Comportamento e resultado técnico

| Cenário | Resultado observado |
| --- | --- |
| Primeira configuração por usuário | Vídeo inteiro com áudio selecionado; Presented=true e ShowOnStartup=false publicados; launcher com checkbox desmarcado. Processo completo em 13,17 s; vídeo reportou 10.159 ms. |
| Reabertura padrão | Sem player/fullscreen/áudio; recursos prontos antes do launcher. Processo completo em 1,28 s. |
| Preferência marcada | Nova reprodução completa e checkbox restaurado marcado; processo em 12,72 s. |
| Preferência desmarcada | Reabertura direta em 1,27 s; ShowOnStartup=false preservado. |
| Interação do checkbox | Teste com Check real: marcar, reconstruir launcher, restaurar marcação, desmarcar; preservar Presented e a escolha ao publicar outro campo. PASS. |
| Falha de configuração | Preservar arquivo inválido, diagnóstico e reversão da marcação que não foi salva. PASS. |
| Sem splash e recurso inválido | Bloquear tela seguinte; não emitir eventos falsos de vídeo. PASS. |
| Fullscreen | Janela ocupa o monitor inteiro; superfície nativa ocupa a maior área proporcional. Regressão C++/Go em seis dimensões, incluindo 1080p/4K; original 1280×720 é ampliado para preencher telas 16:9. |
| Imagem repetida | Quatro reproduções do EXE nesta revisão: duas no ciclo de preferência e mais duas repetições habilitadas. Capturas no início e perto do fim, sem tela preta no ambiente testado; geometria e faixa de áudio verificadas. |
| -open com preferência desabilitada | EXE abriu editor após preparar recursos, sem splash; hash do projeto permaneceu igual. |
| Checks e build | gofmt, importadores de skin/menu, go vet, três testes de identificação dos binários e go test ./...: PASS. Build EXE com recursos Windows: PASS. |

Tempos incluem captura final e encerramento automático da ferramenta; não são promessa de desempenho. Áudio habilitado/selecionado é evidência do player; ouvir a saída e homologar DPI/hardware do usuário continua obrigatório no aceite humano. As capturas foram inspecionadas: vídeo presente em fullscreen e launcher normal com autoria e opção.

### Evidências atuais

- [Checks](../../apps/gerador-jornada/build/036-v10-check.log) e [build](../../apps/gerador-jornada/build/036-v10-build.log).
- [Ciclo primeiro uso/opção](../../apps/gerador-jornada/build/036-v10-policy.log), [eventos da primeira](../../apps/gerador-jornada/build/036-v10-policy/first-events.jsonl) e [sem apresentação](../../apps/gerador-jornada/build/036-v10-policy/default-reopen-events.jsonl).
- [Vídeo no início](../../apps/gerador-jornada/build/036-v10-policy/first-video-early.png), [perto do fim](../../apps/gerador-jornada/build/036-v10-policy/first-video-late.png), [launcher desmarcado](../../apps/gerador-jornada/build/036-v10-policy/disabled-launcher.png) e [marcado](../../apps/gerador-jornada/build/036-v10-policy/enabled-launcher.png).
- [Repetições extras](../../apps/gerador-jornada/build/036-v10-repeated.log) com [captura inicial](../../apps/gerador-jornada/build/036-v10-repeated/video-2-early.png) e [final](../../apps/gerador-jornada/build/036-v10-repeated/video-2-late.png).
- [Editor com -open sem splash](../../apps/gerador-jornada/build/036-v10-open/editor.png) e [eventos](../../apps/gerador-jornada/build/036-v10-open/events.jsonl).

### Roteiro humano vigente

1. Fechar instâncias antigas e abrir o EXE 1.0.0.10 acima. Na primeira execução desta configuração de usuário, a splash deve aparecer uma vez e mostrar somente vídeo, em tela cheia com áudio; aguardar seu término e o fade.
2. Conferir imagem completa na maior área proporcional da tela; em monitor 16:9 deve preencher a tela. Em outra proporção, faixas pretas preservam as bordas. Não existe limite ao tamanho original nem corte de bordas imposto pelo app.
3. Na seleção de projetos, conferir a autoria no rodapé e **Mostrar splash screen ao iniciar** abaixo de abrir projeto. A opção começa desmarcada.
4. Fechar e reabrir: deve ir direto ao launcher, depois de preparar recursos, sem vídeo/som/fade/fullscreen intermediário.
5. Marcar, fechar e reabrir: deve exibir o vídeo inteiro novamente. Desmarcar, fechar e reabrir: deve voltar a abrir direto. Repetir com a opção marcada para procurar ausência intermitente de imagem.
6. Criar projeto em branco ou com exemplo e reabrir recente; não deve haver encerramento inesperado. Conferir nome/pasta e autoria no launcher.
7. Com apresentação habilitada, fechar durante o vídeo: áudio deve parar e launcher não deve aparecer depois. Para repetir primeiro uso no desenvolvimento, usar **outra pasta temporária** em JORNADA_CONFIG_DIR, sem apagar preferências existentes.
8. Conferir som percebido, proporção, fullscreen, fade e DPI 100/150/200% no computador de uso. Comunicar resultado para encerrar 036/T009 e 034/T025; testes técnicos não significam aprovação humana.

### Reproduzir no desenvolvimento

```powershell
cd apps/gerador-jornada
./scripts/check.ps1
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
# Exigem sessão gráfica Windows; reproduzem áudio, usam configurações isoladas em build:
python -X utf8 scripts/test_startup_policy_native.py
python -X utf8 scripts/test_splash_native.py --runs 2
```

O primeiro smoke usa fixtures JSON para as alterações entre processos; o Check real e falhas de gravação são exercitados pelos testes Go. A configuração real do usuário não foi alterada pelos testes. A política e os modelos separados estão em 035/036 e a opção visual em 034. O caminho sem vídeo carrega os recursos antes de mostrar a janela normal, sem player ou eventos de reprodução simulados.

## Revisão vigente — 1.0.0.11 em janela

Solicitação posterior do usuário: reproduzir o vídeo **sem fullscreen**, pois a revisão anterior ainda não foi satisfatória no seu computador. Esta seção substitui a apresentação das revisões 1.0.0.9/1.0.0.10 acima.

**Entrega**: somente JornadaStudio.exe Windows amd64 em `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`, **37860864 bytes**, SHA-256 **`74aa094afa2e98a6e8d5f213ac735307b4ada6c8591b0c4c78bbeb2bff506c85`**. Cópia `dist/JornadaStudio.exe` idêntica. Executor técnico: Codex; ambiente Windows 10 x64 / Go 1.27.2 / Fyne 2.8.2, 2026-10-24. Aceite humano pendente.

### Mudança e evidência

Splash em janela normal centralizada, inicialmente 960×540 unidades Fyne, com barra de título e bordas, sem fullscreen/maximização automática e sem troca desse modo no fade. O vídeo se ajusta proporcionalmente ao cliente nativo da janela. Áudio, espera de vídeo+carga, preferência por usuário e autoria apenas no launcher continuam vigentes; nenhuma alteração no domínio 035 ou formato .jornada.

- [Checks completos](../../apps/gerador-jornada/build/036-v11-check.log): gofmt, importadores de skin/menu, go vet, três testes de identificação de binários e go test ./... aprovados.
- [Build EXE](../../apps/gerador-jornada/build/036-v11-build.log) aprovado.
- [Ciclo nativo de preferência](../../apps/gerador-jornada/build/036-v11-policy.log): primeiro uso com vídeo e opção falsa (12,85 s); reabertura direta (1,27 s); opção habilitada com vídeo (12,51 s); desabilitada com abertura direta (1,28 s). Os tempos incluem captura/encerramento do smoke.
- [Captura inicial em janela](../../apps/gerador-jornada/build/036-v11-policy/first-video-early.png) e [perto do fim](../../apps/gerador-jornada/build/036-v11-policy/first-video-late.png), inspecionadas: imagem presente, moldura/barra de título visíveis. Smoke verificou janela diferente do monitor, não maximizada, cliente separado da moldura e quadro dentro do cliente.
- [Eventos](../../apps/gerador-jornada/build/036-v11-policy/first-events.jsonl): fullscreen=false, áudio selecionado sem mute e término real em 10.161 ms antes da transição. [Launcher com preferência](../../apps/gerador-jornada/build/036-v11-policy/enabled-launcher.png).

Duas reproduções do ciclo completo acima passaram com capturas inicial/final. Duas tentativas adicionais do smoke repetido foram interrompidas por fechamento da janela (evento startup_cancelled) antes da captura final; **não contam como repetições aprovadas**. [Primeira interrupção](../../apps/gerador-jornada/build/036-v11-repeated-interrupted.log) e [segunda tentativa](../../apps/gerador-jornada/build/036-v11-repeated.log). O modo em janela foi implementado e verificado localmente; não afirmar que a falha de imagem foi definitivamente resolvida no computador do usuário. Percepção auditiva e estabilidade nesse ambiente ainda precisam de teste.

### Teste manual desta revisão

1. Fechar instâncias antigas e abrir o EXE 1.0.0.11 indicado acima. Se a primeira apresentação já foi consumida, marcar **Mostrar splash screen ao iniciar** no launcher, fechar e abrir novamente.
2. A splash deve abrir em **janela normal**, centralizada, com bordas/barra de título; não deve ocupar o monitor inteiro ou maximizar automaticamente.
3. Conferir se o vídeo aparece e se o áudio é ouvido. Não deve haver barra de progresso, texto ou autoria sobrepostos ao conteúdo. Aguardar o vídeo inteiro; depois deve ocorrer o fade para o launcher normal.
4. Conferir autoria no rodapé do launcher. Desmarcar a opção, fechar e reabrir: deve ir direto ao launcher sem vídeo. Não apagar configurações existentes para forçar a primeira apresentação.
5. Repetir com a opção marcada para avaliar a ausência intermitente de imagem. Fechar durante o vídeo deve cancelar e parar áudio, sem abrir o launcher depois.
6. Informar o resultado no computador de uso; 036/T009 permanece aberto para imagem, áudio, proporção, fade e DPI. Checks locais não substituem esse aceite.

Ferramentas de desenvolvimento: `scripts/test_startup_policy_native.py` e `scripts/test_splash_native.py --runs 2` em sessão gráfica Windows, usando configurações isoladas em build. Ambas capturam somente a janela do PID criado pelo teste e validam o cliente nativo; não alterar preferências reais do usuário.
