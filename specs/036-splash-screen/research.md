# Research: Reprodução e apresentação

O usuário substituiu o vídeo original por v2, pediu áudio, depois fullscreen sem progresso/legenda e finalmente removeu o rodapé da splash, movendo a autoria para o launcher. Esta é a regra vigente; propostas intermediárias não fazem parte do design final.

O protótipo MCI local recusou o MP4 (erro 277); usar GIF/quadros descartaria áudio e não atenderia à solicitação. MFPlay oferece reprodução nativa com eventos de término e áudio no mesmo player, sem WebView ou executável externo. Referência primária: [Getting Started with MFPlay](https://learn.microsoft.com/en-us/windows/win32/medfound/getting-started-with-mfplay). É uma API nativa existente; a documentação recomenda MediaPlayer/IMFMediaEngine para novas plataformas. A escolha evita framework de UI adicional; rever quando houver mudança de plataforma/requisitos de mídia. Se o Windows não oferecer o subsistema multimídia necessário, apresentar erro, sem alegar reprodução concluída.

O MP4 original fica incorporado e é materializado por hash no cache porque o decoder recebe URL local. Não reconverter ou remover a trilha sonora. Carregar os recursos enquanto toca o vídeo, mas manter a decisão de avanço no domínio 035. Tempo de 10,006 s é metadado/teste, jamais mecanismo de liberação. O cache de imagens conserva pixels e bytes de proveniência usados pelas regressões de skin.

A captura Canvas.Capture do Fyne não inclui a superfície de vídeo nativa. Evidência durante o vídeo usa captura da própria janela visível; captura GL é suficiente após a transição. Som do sistema, mute do Windows e percepção auditiva continuam sob controle/validação do usuário.

## Correções de enquadramento e exibição intermitente — 1.0.0.9

O usuário observou corte/zoom e, depois, ausência intermitente de imagem. O encaixe agora usa GetClientRect como dimensão física, sem multiplicação da escala lógica Fyne; origem normalizada é sempre 0/0/1/1. O cálculo preserva proporção, centraliza, reduz quando necessário e limita a ampliação a 100% da resolução original.

Houve uma regressão intermediária: a saída de razão de aspecto 16/9 foi tomada como resolução, resultando num vídeo de 16×9 pixels. DisplaySize converte a razão utilizando a altura nativa; testes da própria função C++ cobrem 1280×720 + 16/9 e pixels não quadrados. FitRect é verificado em seis dimensões de tela, inclusive retrato e tela pequena.

O controle STATIC inicial não implementava o protocolo de pintura de vídeo e a reprodução começava antes da configuração estar serializada. A janela dedicada JornadaSplashVideoSurface trata WM_PAINT/WM_SIZE conforme o player; mídia pronta posta uma mensagem à thread do HWND, que configura quadro/proporção/áudio e chama Play. Callbacks COM não manipulam o HWND. A superfície protege a pintura da janela pai, evita MoveWindow redundante e repinta margens descobertas. Falha intermitente exige repetição com capturas reais; PLAY e processo encerrado com código zero não são evidência visual suficiente.

A documentação primária de pintura e letterbox está em [Getting Started with MFPlay](https://learn.microsoft.com/en-us/windows/win32/medfound/getting-started-with-mfplay). A ferramenta de desenvolvimento scripts/test_splash_native.py registra imagem no início e perto do fim em várias execuções; amostra cromática recusa tela preta ou vídeo minúsculo, além da revisão visual dos PNGs. Não é dependência do usuário final.

## Pedido vigente de fullscreen e primeira execução — 1.0.0.10

O limite de 100% mantinha o MP4 1280×720 pequeno em monitores maiores; foi removido. Usar a maior área proporcional da tela mantém o quadro integral e o rodapé do próprio vídeo. Os tratamentos de pintura e serialização nativos permanecem; não mudar codec, remover áudio nem introduzir player externo.

Preferência por usuário usa o mesmo diretório dos recentes, em arquivo próprio, sem afetar .jornada. O valor inicial mantém o checkbox desmarcado mas exige a primeira apresentação. A conclusão é registrada depois das condições de domínio e antes de construir o launcher. O caminho sem vídeo carrega os mesmos recursos antes de mostrar a janela normal; preserva erros obrigatórios, sem reproduzir áudio ou apresentar uma splash vazia.

## Relato posterior e teste em janela — 1.0.0.11

O usuário informou que a revisão fullscreen ainda não resolveu a exibição e pediu reproduzir sem esse modo. Aplicar a mudança solicitada, sem assumir causa raiz confirmada: janela normal centralizada e mesmo player/vídeo/áudio. Remover a troca fullscreen no início e na transição reduz as alterações de geometria e permite conferir a imagem com moldura visível. Captura local e eventos não comprovam reprodução no hardware do usuário; conservar o aceite pendente. Atualizar o smoke para validar área cliente, que agora difere do retângulo externo pela barra de título/bordas.
