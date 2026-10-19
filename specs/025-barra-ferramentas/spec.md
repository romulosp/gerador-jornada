# Feature Specification: Barra de ferramentas — Header

**Feature Branch**: `025-barra-ferramentas` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Draft — especificação individual; aceite funcional e visual pendente
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Header/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero acessar arquivos, alinhamento, prévia e exportação sem abandonar o desenho.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve disponibilizar novo, abrir, salvar, prévia e as duas saídas com nomes claros.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve oferecer desfazer/refazer, zoom, ajuste da tela e configuração de fundo.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve alinhar à esquerda, direita, topo, base e centros conforme dimensões reais da tela, sem offsets ocultos do legado.

---

### User Story 2 - Personalizar a aparência com previsibilidade (Priority: P2)

Como autor ou responsável pelo visual, quero ajustar a aparência deste componente sem alterar sua identidade nem suas regras.

**Why this priority**: Permite reproduzir o visual legado e evoluí-lo de forma controlada.
**Independent Test**: Usar o estado de referência do DESIGN.md, alterar uma propriedade visual por vez e comparar resultado e comportamento antes/depois.

**Acceptance Scenarios**:

1. **Given** o estado de referência, **When** uma propriedade visual permitida muda, **Then** a aparência correspondente muda e os três comportamentos da US1 permanecem válidos.
2. **Given** uma personalização fora dos limites, **When** o autor tenta confirmá-la, **Then** o autor recebe diagnóstico no contexto e o último estado válido permanece íntegro.

---

### User Story 3 - Preservar conteúdo e contexto (Priority: P2)

Como autor, quero trocar de seleção ou modo e voltar ao painel sem perda ou alteração implícita dos dados confirmados.

**Why this priority**: O mesmo trabalho precisa continuar coerente ao mudar de contexto.
**Independent Test**: Usar uma fixture com texto acentuado e geometria conhecida, executar o percurso abaixo e comparar os dados antes/depois.

**Acceptance Scenarios**:

1. **Given** uma configuração confirmada, **When** o autor executa o percurso de preservação, **Then** o sistema deve manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
2. **Given** uma operação cancelada ou recusada, **When** o autor retorna ao trabalho, **Then** nenhum dado confirmado foi perdido nem parcialmente substituído.

### Edge Cases

- Janela estreita, ação sem projeto, salvamento cancelado e duas exportações concorrentes. Importação .page continua fora do escopo.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Disponibilizar novo, abrir, salvar, prévia e as duas saídas com nomes claros.
- **FR-002**: Oferecer desfazer/refazer, zoom, ajuste da tela e configuração de fundo.
- **FR-003**: Alinhar à esquerda, direita, topo, base e centros conforme dimensões reais da tela, sem offsets ocultos do legado.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Comandos de edição e arquivo; habilitação e estado ocupado transitórios.
- **Estado de interação**: Pronto; sem seleção; alterações não salvas; exportação ativa; erro recuperável.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `025-barra-ferramentas/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Disponibilidade e atalhos

- **FR-008**: Manter Java/EFL visível e desabilitado temporariamente, sem abrir diálogo de exportação; preservar código e contratos de geração para retomada posterior.
- **FR-009**: Oferecer atalhos de seleção e movimento previstos em 020-elemento-edicao, respeitando foco de entrada de texto e histórico.
- **FR-010**: Identificar a aplicação com a imagem escolhida pelo usuário como ícone incorporado da janela, barra de tarefas e EXE Windows; não depender de caminhos externos. O ícone pertence à identidade do editor, não aos componentes exportados da jornada.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

## Exclusão junto ao ajuste de zoom — 2026-10-18

- **FR-011**: Exibir botão Excluir com ícone de lixeira imediatamente à direita de Ajustar. Iniciar desabilitado e habilitar somente enquanto houver ao menos um componente válido selecionado na tela ativa. Atualizar após seleção, clique vazio, Escape, troca de tela, exclusão e histórico.
- **FR-012**: O botão usa a mesma confirmação Sim/Não e a mesma operação de exclusão que o menu contextual e Delete, sem caminhos de remoção divergentes.

**Aceitação**: ao iniciar ou trocar de tela sem seleção, Excluir fica desabilitado. Selecionar componente habilita; Não mantém a seleção e o botão habilitado; Sim remove e desabilita. Limpar seleção também desabilita. Excluir deve permanecer visível logo após Ajustar.

- **FR-013**: Mostrar a versão 1.0.0.0 no título da janela e na identificação do produto, usando a mesma fonte dos metadados FileVersion/ProductVersion do executável; não confundir versão do software com versões de contratos de domínio.

## Botão Prévia Tela Atual — 2026-10-18

- **FR-014**: Acrescentar Prévia Tela Atual ao lado de Prévia no cabeçalho, com ícone de visualização e texto visível. Prévia mantém o percurso completo do projeto; a nova ação usa a tela ativa. Desabilitar a nova ação quando não houver tela existente selecionada e atualizar seu estado na troca/remoção de tela ou projeto.

## Controle de tela cheia

- **FR-015**: Iniciar o editor em tela cheia por padrão com a ação Restaurar disponível. Exibir ação Tela cheia próxima à identificação do produto e mudar sua legenda para Restaurar quando ativa. No Windows, Maximizar e F12 entram no mesmo modo Fyne sem barra de título/tarefas; F12/Escape/botão retornam. Preservar acesso às ferramentas do editor, rascunho, seleção, zoom e histórico. As teclas nativas Windows são interceptadas somente na janela principal e liberadas no fechamento; fechar e minimizar mantêm o comportamento existente.

**Aceitação**: com componente selecionado e nome em rascunho, entrar/sair do modo sem salvar/descartar involuntariamente, nem desmarcar o alvo. Segurar F12 não alterna repetidamente. Em modo normal, Escape continua limpando a seleção quando não houver campo focado.


**Aceitação de FR-015 na abertura**: ao iniciar o executável atualizado, solicitar tela cheia após a exibição nativa; Restaurar deve sair para uma janela normal. A revisão visual deve conferir que não se trata apenas de maximização com barras visíveis.


## Encerramento em tela cheia e temas do editor

### User Story 3 - Encerrar e personalizar o editor (Priority: P1)

Como autor, quero sair sem restaurar a janela e escolher um tema confortável, preservando o visual das telas desenhadas.

**Independent Test**: abrir em tela cheia, localizar Sair, cancelar/confirmar o encerramento com alterações e alternar os três temas sem alterar o projeto.

**Acceptance Scenarios**:

1. **Given** tela cheia, **When** o cabeçalho é exibido, **Then** Sair está visível; após Restaurar fica oculto, reaparecendo ao entrar novamente em tela cheia.
2. **Given** um projeto alterado, **When** Sair ou fechar pelo sistema é acionado, **Then** a mesma confirmação Sim/Não protege o documento; Não preserva a sessão, Sim encerra; não empilhar confirmações.
3. **Given** primeira execução sem preferência válida, **When** o editor abre, **Then** usa Claro; escolher Escuro ou Azul atualiza os controles e persiste a escolha por usuário para próxima abertura.
4. **Given** componente selecionado e propriedades em edição, **When** muda o tema, **Then** seleção, rascunho, zoom, documento, histórico, fundo e cores do projeto permanecem iguais; prévias/PPTX não adotam a paleta do editor.

- **FR-016**: Exibir Sair com ícone no cabeçalho somente em tela cheia, ao lado de Restaurar. Usar o mesmo fluxo de encerramento da janela, protegendo alterações não salvas com Sim/Não e impedindo confirmações duplicadas.
- **FR-017**: Oferecer seletor Tema com Claro, Escuro e Azul. Claro é padrão inicial e fallback para preferência desconhecida. Aplicar imediatamente à interface e salvar em preferências locais do usuário, fora do arquivo .jornada; não seguir automaticamente o tema do sistema.
- **FR-018**: Aplicar contraste coerente a painéis, textos, cabeçalho, menus, diálogos, campos e estados de foco/seleção/desabilitado. Tema não altera dados/canvas do projeto, prévias ou exportadores. Manter os adaptadores de apresentação separados dos domínios existentes.

**Success Criteria**: usuário encerra com uma ação acessível sem restaurar; três temas selecionáveis, preferência reaplicada na reabertura e nenhuma revisão do projeto causada pela troca. Homologação visual em Windows e DPI reais permanece manual.
