# Implementation Plan: Splash e mídia nativa

Um módulo Go, composição manual, Windows x64. `cmd/gerador-jornada/main.go` conecta carregamento e navegação; o núcleo permanece em `internal/startup` (035). Adaptadores: `desktop/splash.go`, `splash_windows.go`, `splash_other.go`, `startup_resources.go`; recurso em `builtinassets/splash.go` e `assets/splash-v2.mp4`; reprodução em `nativevideo/player_windows.go/.cc/.h`.

MFPlay/Media Foundation nativos reproduzem MP4 e áudio numa superfície filha do HWND Fyne, sem navegador, processo de player ou instalação de codec distribuída pelo projeto. CGO usa a toolchain existente; C++ com bibliotecas de suporte estáticas no EXE e DLLs multimídia do próprio Windows. Callbacks COM somente em C++; Go observa campos atômicos, na thread gráfica, e traduz eventos para o domínio. Quando há apresentação, recursos carregam em goroutines com contexto; somente fyne.Do toca UI. Não reaplicar tema depois de registrar janela ainda não exibida.

A splash abre em janela normal centralizada (960×540 unidades Fyne), sem fullscreen ou maximização automática, e reaproveita sua janela para o launcher. Fade out/in por opacidade nativa de aproximadamente 220 ms por fase; término do vídeo vem de PLAYBACK_ENDED, nunca de tempo fixo. Player é encerrado antes de trocar conteúdo. A janela normal recebe o layout/tamanho de boas-vindas durante a parte invisível da transição, sem trocar modo de exibição.

Verificar domínio, integridade da mídia, caches/cena, controles mínimos, regressão do launcher, startup nativo, carga rápida/lenta e cancelamento. Build somente por scripts/build.ps1. Evidências e limites em quickstart.

## Ajustes comprovados na implementação

O player cria sua própria classe de janela e serializa configuração/Play por mensagem após MEDIAITEM_SET. WM_PAINT/WM_SIZE atualizam o vídeo; a área pai conserva pintura das margens. Resize observa GetClientRect diretamente. DisplaySize interpreta razão de aspecto e FitRect faz contain com ampliação/redução até a maior área disponível, sem limite de tamanho original; ambas as funções nativas têm regressões Go/CGO. O teste repetido do EXE é scripts/test_splash_native.py; testes de carga lenta/criar/reabrir/cancelar/falhar ficam em cmd/splash-smoke com build tag splashsmoke. Não existe fallback que salte ou silencie o vídeo.

## Inicialização optativa — 1.0.0.10

`desktop/startup.go` concentra a escolha anterior à criação do player. Caminho com vídeo usa Splash e registra conclusão antes de construir launcher; caminho sem vídeo prepara recursos de forma síncrona antes de mostrar a primeira janela normal, sem criar player, áudio ou fullscreen. A carga local é curta; não se adiciona uma tela de espera que pareça uma splash preta.

`desktop/startup_preferences.go` resolve o diretório comum dos recentes, lê `startup-preferences.json` e publica alterações com `atomicfile.Write`. Ausência não é erro; JSON inválido é preservado e diagnosticado. Cada atualização relê o estado para preservar o outro campo. Não escrever no construtor do checkbox. Campo alterado só é confirmado na UI depois da gravação; falha reverte a seleção. Testes de integração em `startup_preferences_test.go`; smoke do EXE em `scripts/test_startup_policy_native.py`.

## Execução em janela — 1.0.0.11

Em `desktop/splash.go`, remover a entrada em fullscreen e a saída no fade; mostrar uma janela padrão centralizada, 960×540 unidades Fyne. Log splash_open inclui fullscreen=false. Player continua com o mesmo quadro original, áudio e encaixe proporcional ao cliente nativo. No smoke, comparar a janela com o monitor para recusar fullscreen/maximização, exigir barra de título e medir o cliente separado da moldura antes de conferir a superfície. Domínio, política e armazenamento não mudam. Checks/build e execução real em quickstart; não declarar resolvida a falha reportada no hardware do usuário sem validação dele.
