# Design: Splash em janela e autoria no launcher

**Versão visual vigente: 1.0.0.11.** Substitui o fullscreen da 1.0.0.10 por solicitação do usuário, que ainda relata falha de imagem.

Splash: janela normal centralizada, com barra de título/bordas padrão, área inicial 960×540 unidades Fyne. Não abrir fullscreen ou maximizar automaticamente. Área interna preta com o vídeo v2 original e áudio; manter proporção e mostrar o quadro inteiro. Não sobrepor textos, autoria, barra, spinner ou percentual ao conteúdo. Fechar cancela, sem abrir o launcher.

O player usa GetClientRect em pixels físicos; o encaixe proporcional considera somente o cliente, sem bordas/barra de título e sem repetir escala de DPI. A imagem ocupa a maior área proporcional disponível, com faixas pretas se necessário. A superfície mantém tratamento de pintura/tamanho.

Após vídeo e recursos completos, fade out de aproximadamente 220 ms, encerrar player, montar o launcher na janela normal existente e fade in de aproximadamente 220 ms. Não há troca de modo fullscreen. A autoria fica somente no rodapé do launcher:

Desenvolvido por Rômulo Penha - Canais Parceiro - 03

A primeira apresentação e a opção “Mostrar splash screen ao iniciar”, desmarcada por padrão, continuam conforme 035/034. Alterar a opção afeta as próximas inicializações. Não modificar bytes do vídeo, skins ou fonte dos componentes. O teste nativo deve verificar moldura/cliente, ausência de fullscreen/maximização, imagem no início e perto do fim e eventos de áudio/término. A reprodução no computador do usuário ainda exige teste e aceite explícito.
