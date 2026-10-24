# Tasks: Adaptadores da splash

- [X] T001 Preservar vídeo v2 original com áudio e hash, incorporar e extrair cache atomicamente.
- [X] T002 Implementar player nativo Windows/COM com eventos reais e liberação de recursos.
- [X] T003 Implementar carregamento antecipado de skins, imagens, fontes e catálogo reutilizando caches operacionais.
- [X] T004 Conectar Splash ao modelo 035 e ao startup GUI, incluindo -open; utilitários sem interface não exibem vídeo.
- [X] T005 Aplicar visual final: fullscreen somente vídeo, sem textos/progresso/rodapé; autoria literal somente no launcher.
- [X] T006 Implementar fade entre apresentação e seleção de projetos, preservando correção de tema da 1.0.0.5.
- [X] T007 Executar testes de mídia/incorporação/carga, cancelamento/falha e regressões; smoke nativo rápido/lento e abertura posterior.
- [X] T008 Executar checks/build EXE, registrar hashes/capturas/eventos e atualizar README/coordenadora/domínios afetados.
- [ ] T009 Homologação humana da revisão vigente em janela: áudio, proporção, imagem, fade, DPI e aceite da revisão final.

- [X] T010 Corrigir enquadramento com cliente físico, origem integral, razão convertida em resolução e contain limitado a 100%; substituir STATIC por janela com pintura/tamanho e Play serializado. Regressões nativas e quatro aberturas com capturas aprovadas na 1.0.0.9; não encerrar T009 por isso.

- [X] T011 Aplicar fullscreen com maior área proporcional, removendo limite de pixels originais; testar dimensões 1080p/4K e preservar quadro/áudio.
- [X] T012 Integrar preferência atômica por usuário e política 035: primeira apresentação única, repetição opcional e caminho sem vídeo com recursos obrigatórios. Testar checkbox/persistência/falhas e sequência de processos reais.
- [X] T013 Registrar checks, build 1.0.0.10, capturas/logs e README; disponibilizar somente EXE para teste do usuário. Homologação continua em T009.

- [X] T014 Implementar apresentação em janela normal centralizada, 960×540 unidades Fyne, sem fullscreen/maximização ou troca de modo na transição; manter áudio, quadro completo e política. Atualizar smoke para validar moldura/cliente e recusar fullscreen.
- [X] T015 Executar checks/build EXE 1.0.0.11 e smoke nativo em janela; registrar revisão, capturas, README e metadados para teste do usuário. Aceite e reprodução no computador do usuário continuam pendentes em T009.
