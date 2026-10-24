# Contract: Eventos de inicialização

O coordenador serializa notificações de carregamento e reprodução. `PlaybackCompleted` só tem efeito depois de `PlaybackStarted`. `ResourcesCompleted(nil)` marca Loaded; erro conserva falha. `BeginTransition` retorna true exatamente uma vez após ambos, sem falha ou Close. Fechar impede transição posterior. Erro nunca é convertido em término normal. Consumidor: `adapter/desktop.Splash`; produtor de reprodução: `adapter/nativevideo`. O domínio não conhece esses tipos concretos.

A consulta `Preference.ShouldShow()` ocorre antes de criar player. Se falsa, não emitir PlaybackStarted/Completed fictícios: carregar recursos, tratar erro e revelar a janela normal diretamente. Se verdadeira, manter o contrato completo da Session. Após liberação e antes de montar a tela seguinte, o adaptador pode publicar Presented=true; nunca habilitar ShowOnStartup automaticamente.
