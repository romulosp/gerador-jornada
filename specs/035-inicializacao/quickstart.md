# Quickstart: Inicialização

Na pasta `apps/gerador-jornada`, executar `go test ./internal/startup -count=1`. O teste `TestBothCompletionsRequired` verifica carga rápida/lenta, início obrigatório, falha, cancelamento e liberação única. A evidência nativa e o roteiro humano são coordenados pela 036/quickstart. Não confundir teste puro com reprodução real de MP4/áudio.

## Evidência técnica final — 1.0.0.9

`TestBothCompletionsRequired` passou no check completo da 036. Integração nativa validou carga rápida, vídeo concluído antes da carga artificial de 14 s, cancelamento e erro sem liberação. As quatro execuções do EXE e os fluxos completos estão em [036/quickstart](../036-splash-screen/quickstart.md). Domínio não depende de Fyne, Windows ou filesystem. Aceite humano da apresentação permanece em 036/T009.

## Política por usuário — revisão 1.0.0.10

`TestPresentationPolicy` cobre valor inicial (mostrar uma vez, repetição falsa), apresentação concluída, habilitar e desabilitar repetição. Testes do adaptador `TestLauncherSplashPreferencePersistence`, `TestSplashPreferenceWriteFailure` e `TestSkippedSplashStillLoads` cobrem persistência por Check real, restauração, preservação de campos, configuração corrompida e bloqueio por falha de recurso sem vídeo. Check completo aprovado em build/036-v10-check.log. Evidência de processos reais e roteiro atualizado em [036/quickstart](../036-splash-screen/quickstart.md#revisão-atual--10010).
