# Feature Specification: Tela inicial (launcher) e Estrutura de Projeto

**Feature Branch**: `034-launcher-estrutura-projeto` (identificador documental; nenhuma branch Git criada)

**Created**: 2026-10-24

**Status**: Implementada tecnicamente na revisão 1.0.0.5; aceite humano pendente

**Input**: Pedido do usuário: ao abrir o programa, mostrar uma tela inicial no estilo Android Studio com os projetos recentes e as opções de criar novo projeto (em branco, somente com o fundo atual, ou com exemplo — o cenário de carregamento atual). Adicionar, no canto esquerdo, uma aba "Estrutura projeto" que lista, em árvore, todos os arquivos `.jornada` presentes no workspace em uso, com botão de três pontos para criar ou importar uma jornada. Quando o componente Menu aplicação estiver na jornada e uma folha sem destino for clicada na prévia (hoje "SELEÇÃO SIMULADA: NENHUM DESTINO CONFIGURADO"), sugerir imediatamente a criação de uma jornada com o nome do botão; a nova jornada aparece na árvore e, ao clicar novamente, a prévia navega para ela.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Iniciar a partir da tela de boas-vindas (Priority: P1)

Como autor, ao abrir o programa quero ver os projetos recentes e escolher entre criar um projeto em branco (somente com o fundo atual) ou com o exemplo (o cenário de carregamento atual), em vez de cair direto num documento.

**Why this priority**: É a porta de entrada; determina como o autor começa qualquer trabalho e evita criar ou descartar documentos sem querer.

**Independent Test**: Abrir o programa sem argumentos e observar a tela inicial; criar um projeto em branco e um com exemplo; fechar e reabrir para conferir que os recentes aparecem.

**Acceptance Scenarios**:

1. **Given** o programa inicia sem um documento aberto, **When** a janela é exibida, **Then** aparece a tela inicial com a lista de projetos recentes e as ações de criar novo projeto.
2. **Given** a tela inicial, **When** o autor escolhe "Novo projeto em branco", **Then** é criado um documento novo somente com o fundo atual, sem conteúdo de exemplo.
3. **Given** a tela inicial, **When** o autor escolhe "Novo projeto com exemplo", **Then** é criado o cenário de carregamento atual (exemplo editável), preservando o nome escolhido.
4. **Given** um projeto aberto e salvo, **When** o programa é reaberto, **Then** esse projeto aparece na lista de recentes e pode ser reaberto.

---

### User Story 2 - Ver e gerenciar jornadas na Estrutura projeto (Priority: P1)

Como autor, quero uma aba "Estrutura projeto" no canto esquerdo que mostre, em árvore, todas as jornadas (arquivos `.jornada`) presentes no workspace em uso, e permita criar ou importar uma jornada.

**Why this priority**: Materializa o workspace com múltiplas jornadas e é a base para a criação rápida e a navegação entre jornadas.

**Independent Test**: Com um workspace aberto, abrir a aba Estrutura projeto e ver a árvore; usar o menu de três pontos para criar uma jornada e importar uma `.jornada` existente; conferir que a árvore reflete as mudanças.

**Acceptance Scenarios**:

1. **Given** um workspace com jornadas, **When** o autor abre a aba "Estrutura projeto", **Then** a árvore lista as jornadas do workspace (ex.: "projeto principal" com "Arrecadação de boleto de cobrança").
2. **Given** o botão de três pontos da árvore, **When** o autor escolhe "Criar jornada", **Then** nasce uma nova jornada em branco e ela aparece na árvore.
3. **Given** o botão de três pontos da árvore, **When** o autor escolhe "Importar jornada existente", **Then** uma `.jornada` existente é incorporada ao workspace e passa a aparecer na árvore, sem alterar o arquivo de origem.

---

### User Story 3 - Criar e navegar para a jornada de um item do menu (Priority: P1)

Como autor, quando um item do Menu aplicação sem destino é clicado na prévia, quero receber uma sugestão imediata para criar uma jornada com o nome do item; depois de criada, quero que a prévia navegue para essa jornada.

**Why this priority**: Converte um beco sem saída ("SELEÇÃO SIMULADA: NENHUM DESTINO CONFIGURADO") em uma ação produtiva e fecha o ciclo com a árvore e a navegação.

**Independent Test**: Em uma jornada com o Menu aplicação, clicar em uma folha sem destino na prévia; aceitar a sugestão de criação; conferir a nova jornada na árvore e, ao clicar novamente, ser direcionado a ela.

**Acceptance Scenarios**:

1. **Given** o Menu aplicação na jornada e uma folha sem destino configurado, **When** o autor clica na folha durante a prévia, **Then** a aplicação informa a seleção simulada e sugere imediatamente criar uma jornada com o nome da folha.
2. **Given** a sugestão aceita, **When** a jornada é criada, **Then** ela aparece na árvore da Estrutura projeto sob o projeto principal e o item passa a ter aquele destino.
3. **Given** a jornada criada a partir do item, **When** o autor executa a jornada inicial (que mostra o Menu aplicação) e clica no mesmo item, **Then** a prévia navega para a jornada criada.

### Edge Cases

- Nome vazio ou repetido ao criar jornada; nome derivado de item de menu com caracteres inválidos para nome de arquivo.
- Projeto recente que foi movido, renomeado ou excluído: a entrada deve indicar indisponibilidade, sem quebrar a tela inicial.
- Workspace sem nenhuma `.jornada`; importação de arquivo corrompido ou de versão incompatível.
- Item de menu que já possui destino; item desabilitado; item com submenu (não é folha) clicado.
- Criar a jornada sugerida e cancelar: a sugestão não deve criar nada nem alterar o projeto sem confirmação.
- Jornada destino apagada depois de configurada: o item deve voltar a indicar destino ausente, sem navegação implícita.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Ao abrir sem documento, o sistema DEVE exibir uma tela inicial com a lista de projetos recentes e as ações de novo projeto, sem abrir um documento automaticamente.
- **FR-002**: A tela inicial DEVE permitir criar projeto em branco contendo somente o fundo atual, e projeto com exemplo contendo o cenário de carregamento atual (exemplo editável), com nome escolhido pelo autor.
- **FR-003**: O sistema DEVE manter a lista de projetos recentes localmente e permitir reabrir um recente; uma entrada indisponível DEVE ser indicada sem interromper as demais.
- **FR-004**: O painel esquerdo DEVE ganhar a aba "Estrutura projeto" que lista, em árvore, todas as jornadas (arquivos `.jornada`) do workspace em uso, com o projeto principal como raiz e as jornadas como descendentes.
- **FR-005**: A aba DEVE oferecer um menu de três pontos com "Criar jornada" e "Importar jornada existente", usando o componente de árvore; criar e importar DEVE refletir imediatamente na árvore.
- **FR-006**: Importar uma `.jornada` DEVE incorporá-la ao workspace sem modificar o arquivo de origem; um arquivo inválido ou incompatível DEVE ser recusado com diagnóstico localizado.
- **FR-007**: O sistema DEVE tratar o workspace como a pasta do sistema em uso: o nó raiz ("projeto principal") representa essa pasta e as jornadas são os arquivos `.jornada` presentes nela; os projetos recentes DEVE referenciar essa pasta (ou sua `.jornada` principal).
- **FR-008**: Quando uma folha do Menu aplicação sem destino for ativada na prévia, o sistema DEVE informar a seleção simulada e sugerir imediatamente a criação de uma jornada com o nome da folha.
- **FR-009**: Ao aceitar a sugestão, o sistema DEVE criar a jornada, vinculá-la como destino do item e exibi-la na árvore da Estrutura projeto.
- **FR-010**: Ao ativar um item cujo destino é uma jornada, a prévia DEVE exibir a tela inicial da jornada destino e oferecer um comando de voltar explícito que retorna à jornada de origem, mantendo um empilhamento simples de jornadas visitadas.
- **FR-011**: O sistema NÃO DEVE inferir rota pelo nome, atalho ou metadado; apenas destino explicitamente confirmado pelo autor habilita a navegação para a jornada.
- **FR-012**: Cancelar a sugestão ou a criação NÃO DEVE publicar estado parcial nem alterar o projeto; o histórico (desfazer/refazer) DEVE preservar o último estado confirmado.
- **FR-013**: Manter controles e mensagens em português; preservar nomes com acentos/Unicode e normalizar apenas o necessário para nome de arquivo, sem quebrar referências.

### Key Entities *(include if feature involves data)*

- **Workspace/Projeto**: pasta do sistema em uso que agrupa as jornadas; o nó raiz da árvore representa essa pasta, sem criar um arquivo de workspace novo.
- Correção integrada com 037/FR-020 em 2026-10-25: criar projeto salva o pacote diretamente na pasta escolhida. O nome define o `.jornada`, sem acrescentar subpasta; homônimos anteriores permanecem preservados.
- **Jornada**: documento `.jornada` (telas, roteiro e roteiros de demonstração), descendente do workspace na árvore.
- **Projeto recente**: entrada da tela inicial que referencia um workspace/jornada aberto anteriormente.
- **Item de menu (folha)**: origem do atalho de criação; pode passar a apontar para uma jornada como destino.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Um autor abre o programa, identifica os recentes e cria um projeto em branco ou com exemplo em até 3 interações, sem editar ou descartar um documento indevido.
- **SC-002**: 100% das `.jornada` válidas do workspace aparecem na árvore da Estrutura projeto; criar e importar refletem imediatamente, sem alterar arquivos de origem.
- **SC-003**: Ao clicar numa folha sem destino, a sugestão de criação aparece imediatamente; após aceitar, a jornada aparece na árvore e a próxima ativação navega para ela, com zero mudanças no projeto quando o autor cancela.
- **SC-004**: Nenhum caso inválido (importação corrompida, recente indisponível, nome inválido, item com submenu) publica estado parcial ou perde o último estado confirmado.

## Assumptions

- Esta feature estende os domínios Projeto (workspace/múltiplas jornadas) e Jornada (navegação entre jornadas); a tela inicial e a aba Estrutura projeto são componentes de interface, não novos domínios, conforme a constituição.
- O workspace é a pasta do sistema em uso (não um novo arquivo/formato de workspace); importar uma jornada significa trazê-la para essa pasta sem alterar a origem.
- "Fundo atual" é o fundo incorporado padrão (Loterias CAIXA) já aplicado a novas telas; "cenário de carregamento atual" é o exemplo editável carregado hoje ao abrir o programa.
- A lista de projetos recentes é local ao usuário, sem telemetria nem sincronização em rede (fora do escopo).
- O componente de árvore reutiliza o mesmo padrão visual e hierárquico dos componentes de estrutura existentes, sem criar um tipo de árvore novo.
- Navegação entre jornadas é um destino explícito confirmado pelo autor; nenhuma inferência por nome, atalho ou metadado.
- A saída Java/EFL e o PPT não fazem parte desta entrega; a navegação entre jornadas é, inicialmente, comportamento de prévia do editor.

## Entrada pela splash e autoria — integração 035/036 / 1.0.0.9

Quando há apresentação, o launcher é construído na janela devolvida pela splash, depois de vídeo e recursos concluídos e do fade. Sem apresentação, abre diretamente após recursos prontos. Sua funcionalidade de criar/abrir/recentes permanece nesta spec. O rodapé exibe somente a autoria solicitada, centralizada: Desenvolvido por Rômulo Penha - Canais Parceiro - 03. A splash não contém autoria. Modelo do ciclo em 035 e apresentação em 036; manter seus backlogs separados.

## Preferência da apresentação — 1.0.0.10

- **FR-014**: Na seleção de projetos, abaixo de “Abrir projeto existente”, oferecer “Mostrar splash screen ao iniciar”, desmarcado por padrão. Restaurar a escolha salva e persistir alterações imediatamente; marcar repete nas próximas aberturas, desmarcar suprime. Não iniciar vídeo ao clicar no controle. Erro de gravação reverte a marcação e informa diagnóstico. Política de primeira execução em 035/FR-008 e persistência/coordenação em 036/FR-013–FR-014.

A seleção pode ser revelada após a splash ou diretamente quando a preferência dispensa apresentação; autoria literal continua somente no rodapé do launcher. Não criar/abrir projeto ao alterar a opção.
