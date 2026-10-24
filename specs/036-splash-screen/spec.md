# Feature Specification: Splash screen com vídeo e áudio

**Created**: 2026-10-24
**Status**: Implementação técnica 1.0.0.11; validação final registrada no quickstart; aceite humano pendente
**Input vigente**: usar `material-legado/splash_screen/video-splash-programa-v2.mp4` com áudio, apresentação em janela normal centralizada, inicialmente 960×540 unidades Fyne, sem fullscreen ou maximização automática, somente com o vídeo, sem barra, legenda de carga ou autoria. Autoria movida para a seleção de projetos. Quando exibida, transição obrigatória após vídeo e carregamento. Exibir na primeira execução; nas seguintes, somente por preferência marcada no launcher, desmarcada por padrão.

## Escopo e domínio

Apresentação e adaptadores do [domínio Inicialização 035](../035-inicializacao/spec.md). A regra de liberação pertence ao núcleo 035; este incremento implementa mídia, carregamento, janela e animação. Projeto, Catálogo, Tela, Jornada e Demonstração preservam seus modelos e formato.

## Requisitos

- **FR-001**: Quando a política 035/FR-008 exigir, iniciar a interface com vídeo v2 incorporado, em janela normal com barra de título e bordas, sem dependência da pasta de desenvolvimento. Caso contrário carregar os recursos e abrir a janela normal diretamente.
- **FR-002**: Reproduzir a faixa de áudio original habilitada, respeitando volume/mudo do Windows, sem controles de pular, seek ou autoplay de player externo.
- **FR-003**: Carregar em paralelo skins/recursos com integridade, imagens decodificadas, fontes, defaults do catálogo, fundo/ícone e preferências necessárias. Para `-open`, validar o documento durante essa preparação.
- **FR-004**: Encaminhar eventos reais ao modelo 035; recursos prontos não cortam o vídeo e vídeo concluído não dispensa recursos restantes.
- **FR-005**: Quando liberado pelo domínio, aplicar fade out/in e mostrar o launcher na janela normal de boas-vindas. Para `-open`, mostrar o editor após as mesmas condições.
- **FR-006**: A splash deve exibir somente vídeo, preservando proporção; sem barra de progresso, textos de carregamento ou rodapé. O texto **Desenvolvido por Rômulo Penha - Canais Parceiro - 03** fica centralizado no rodapé do launcher, sem vírgula final.
- **FR-007**: Fechar durante a inicialização cancela tarefas, encerra áudio e libera o player, sem abrir a tela seguinte.
- **FR-008**: Falha de vídeo/áudio ou recursos obrigatórios mostra diagnóstico em português e impede avanço; não simular sucesso nem trocar silenciosamente por reprodução muda.
- **FR-009**: Preservar a correção de tema/viewport da 1.0.0.5 durante splash → launcher → editor.
- **FR-010**: Utilitários explícitos de salvar/exportar sem interface não executam splash. A entrega permanece somente JornadaStudio.exe Windows x64.

## Cenários

S1: carga rápida, aguardar vídeo inteiro e só então revelar launcher com autoria. S2: carga lenta, vídeo acaba e a apresentação conserva o último quadro até o carregamento terminar. S3: cancelar/erro não abre launcher. S4: vídeo com áudio ativo e imagem ajustada ao cliente da janela normal, sem textos sobrepostos ao conteúdo. S5: criar/reabrir projeto após transição continua funcionando. S6: `-open` respeita splash e mostra documento solicitado sem alterações indevidas.

Aceite humano separado: avaliar áudio, modo em janela, proporção, fade e legibilidade de autoria no launcher em DPI 100/150/200%. Screenshots e eventos não substituem percepção auditiva no computador do usuário.

## Enquadramento e regressão após teste visual

- **FR-011**: Mostrar integralmente o quadro, sem zoom/corte, respeitando cliente físico e proporção. Ajustar proporcionalmente à maior área possível do cliente da janela normal, sem incluir bordas e barra de título no tamanho de reprodução; centralizar com faixas pretas se necessário. O pedido vigente substitui a apresentação fullscreen das revisões anteriores.
- **FR-012**: Preparar o renderer e atender sua pintura nativa antes/durante Play; demonstrar imagem real no início e perto do fim em inicializações consecutivas, recusando evidência apenas de áudio/eventos.

- **FR-013**: Ler a preferência antes de abrir a janela/player da apresentação. Arquivo ausente exibe uma vez; registrar conclusão sem habilitar repetição. Gravar por usuário em JSON atômico, separado do projeto e versão; preservar configurações inválidas com diagnóstico.
- **FR-014**: Sem apresentação, preparar os mesmos recursos antes de exibir launcher/editor, sem som, fullscreen ou fade de vídeo. Falha de recursos continua bloqueando avanço.

S7: primeira exibição completa → launcher com opção desmarcada → reabrir sem splash. S8: marcar → reabrir com vídeo completo → desmarcar → reabrir sem vídeo. S9: cancelamento/falha da primeira apresentação não registra conclusão. A opção e seu comportamento visual ficam em 034/FR-014; política pura em 035.

## Revisão visual solicitada — 1.0.0.11

O usuário ainda não confirmou imagem no seu computador. A solicitação atual é testar o vídeo **sem fullscreen**: janela normal, centralizada, com moldura padrão e área inicial 960×540 unidades Fyne. Nenhum salto para fullscreen/maximização no início ou na transição. Preservar áudio, espera de vídeo+carga, primeiro uso/opção marcada, autoria somente no launcher e modelo de domínio 035 separado. Evidência técnica local não significa que o problema foi resolvido no ambiente do usuário.
