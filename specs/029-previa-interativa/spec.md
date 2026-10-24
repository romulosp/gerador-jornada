# Feature Specification: Prévia interativa — Preview

**Feature Branch**: `029-previa-interativa` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Preview/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Prévia completa parte da inicial e permite rotas; Prévia Tela Atual isola uma tela. Cena, decorações/fontes e regiões de teclado são compartilhadas; campos/opções/menu/código de barras mantêm runtime separado.

**Pendências e limites:** Máscaras incrementais, todas as novas famílias, scroll/modal de todos os tipos e bloqueios integrais de perfis ainda não foram entregues. Os testes executados não homologam os critérios geométricos/visuais no Office ou em todos os DPIs.

**Fontes verificadas:** [preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go), [screen_preview.go](../../apps/gerador-jornada/internal/application/screen_preview.go), [legacy_preview_test.go](../../apps/gerador-jornada/internal/adapter/desktop/legacy_preview_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero percorrer a jornada e experimentar entradas antes de exportar.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve abrir a tela inicial em uma janela independente com o mesmo fundo e controles.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve navegar por botões e operar campos, seleção e teclado com estado isolado.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve restaurar valores iniciais a cada visita, inclusive ao voltar à mesma tela.

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

- Ciclo de navegação, fechar durante callback, projeto inválido e ausência de imagem local. Prévia não executa Java nem regras do sistema real.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Abrir a tela inicial em uma janela independente com o mesmo fundo e controles.
- **FR-002**: Navegar por botões e operar campos, seleção e teclado com estado isolado.
- **FR-003**: Restaurar valores iniciais a cada visita, inclusive ao voltar à mesma tela.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Snapshot do projeto; tela corrente, foco, valores e escolhas runtime não são gravados no documento.
- **Estado de interação**: Aberta; navegando; entrada alterada; tela revisitada; erro de referência; fechada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `029-previa-interativa/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Prévia Tela Atual — 2026-10-18

- **FR-008**: Disponibilizar sessão exclusiva da tela selecionada pelo autor, mesmo que ela não seja a tela inicial. O snapshot inclui somente essa tela e seus recursos disponíveis; não inclui ações de navegação ou roteiro. Botões não saem da tela; campos, escolhas e teclado permanecem interativos.
- **FR-009**: Validar os dados da tela selecionada e referências locais necessárias, sem bloquear por geometria ou navegação incompletas de outras telas. Tela inexistente ou composição local inválida produz erro em português. Não modificar o projeto, sua tela inicial, seleção do editor, revisão ou histórico; alterações de valores existem somente durante a sessão.
- **FR-010**: Manter a prévia de projeto completo iniciada na tela inicial e com navegação; a prévia individual tem título Prévia Tela Atual seguido do nome da tela e informa que a navegação entre telas está desativada. Exibir o conteúdo confirmado do modelo; alterações no formulário ainda não salvas não entram no snapshot.

### Acceptance Scenarios

1. **Given** três telas com a primeira como inicial e a segunda ativa no editor, **When** escolher Prévia Tela Atual, **Then** somente a segunda abre; clicar em botão de navegação não troca de tela.
2. **Given** uma tela atual válida e outra com dimensões inválidas ou destino inexistente na jornada, **When** abrir a prévia individual, **Then** a tela atual é apresentada; a prévia do projeto mantém sua validação completa.
3. **Given** a prévia individual aberta, **When** editar um campo e fechar, **Then** documento, tela inicial, seleção e histórico do editor permanecem iguais.
4. **Given** nenhuma tela existente ativa, **When** consultar o cabeçalho, **Then** Prévia Tela Atual fica desabilitado.


## Isolamento do tema do editor

- **FR-011**: A paleta escolhida no editor não deve recolorir o conteúdo da prévia da jornada. Manter primitivas com cores do modelo; a base inicial isola widgets runtime no tema claro. Na expansão de 2026-10-19, pintar controles pela cena/skin comum, sem widget Fyne opaco substituindo aparência. A decoração descritiva pode acompanhar o editor. Alternar Claro/Escuro/Azul com a prévia aberta não altera controles/snapshot. Responsabilidade desktop, integrada a 025-barra-ferramentas/FR-018.

## Paridade interativa pós-v1 — 2026-10-19

- **FR-012**: Consumir a mesma cena/resolução de skins/fontes/estados do canvas e PPT, com pintura compartilhada e hit regions/eventos explícitos; widgets custom quando necessários. Não encobrir controle fiel com Button/Entry/choices nativos opacos. Biblioteca é dado estático gerado/recursos do adaptador; runtime não lê XML legado nem executa Java.
- **FR-013**: Operar entrada/máscaras, toggles/grupos, Combo/UF/Table/Report, teclados numérico/alfanumérico/calculadora, scroll/contêiner/dialog e demais estados locais suportados por família, sem negócio/periféricos/serviços. Foco/cursor/modificadores/seleção corrente/scroll/modal são isolados e reiniciam por visita; enabled/valores iniciais vêm do snapshot. Preservar escopo sem navegação da Prévia Tela Atual.
- **FR-014**: Localizar fonte/recurso/delegate/estado/capacidade ausente e manter fidelidade bloqueada; prévia exploratória parcial não homologa entrada ou acervo. Preservar bounds/IDs/overrides antigos e tema/chrome, compartilhando face real sem fallback silencioso.

**Aceitação**: fixture de cada classe/skin/estado compara pintura com canvas/PPT/referência independente e eventos com contrato; fechar/revisitar não muda documento/histórico. 001/T114 é sessão, T115 desktop, T111 cena e T120–T122 manuais/aceite. Automatizados/autorização do plano não são homologação; Java permanece pausado.

**Gate objetivo (2026-10-24):** um teste de prévia deve falhar se `Style.Skin == ""` produzir aparência genérica, se a árvore de objetos contiver `widget.Button`/`widget.Entry`/`widget.Select` visível sobre a decoração, ou se as primitivas da prévia divergirem das do canvas/PPT. Fonte física (`canvas.Text.FontSource`) e alpha explícito também são asserções, não inspeção visual. Skill obrigatória: [jornada-visual-parity](../../.agents/skills/jornada-visual-parity/SKILL.md).

## Interação de código de barras — planejamento 2026-10-23

- **FR-015**: Operar [031](../031-codigo-barras-efl/spec.md) com subcampo ativo, validação sequencial por par/DV geral, captura integral atômica e diálogo “Código de barras inválido” / “OK, entendi”; respeitar teclado numérico/habilitação e reset por visita. [032](../032-codigo-barras-processado-efl/spec.md) apresenta somente dados confirmados de seus blocos; combos são autoria. Usar regras puras de Catálogo e estado de sessão isolado, sem gerar processamento financeiro, navegação ou transferência automática entre os dois compostos.

**Aceitação**: erro local bloqueia avanço da entrada, captura inválida limpa só a instância destinatária e confirmação do modal devolve foco ao primeiro campo; outro componente/documento permanece intacto. Revisita começa vazia, ficha preserva configuração. Contrato visual em 031/032; integração em [001/T135/T137](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23), ainda planejada. Periférico real, câmera e fluxo condicional estão fora do escopo; colagem permite verificar sem hardware.

## Navegação local do Menu aplicação — 2026-10-23

- **FR-016**: Operar [033](../033-menu-aplicacao-efl/spec.md) como uma interface com vista raiz/submenu selecionado e caminho local, inicialmente no menu autorado. Abrir submenus habilitados, Voltar um nível e Cancelar para raiz; Sair fecha somente prévia, não editor ou autenticação. Estado isolado por instância/visita, sem gravar no projeto. Folha sem destino informa seleção simulada; somente destino de tela explicitamente autorado e validado por Jornada permite transição na prévia completa. Não inferir rota ou dispatch de texto/categoria/atalho/Movement*.

**Aceitação**: raiz→submenu→aninhado→Voltar conserva contexto; Cancelar retorna raiz mesmo com inicial diferente; fechar/reabrir/revisitar restaura inicial. Enabled efetivo bloqueia ativação. Prévia Tela Atual explora menus mas não segue rotas nem exige destinos externos no snapshot individual. Canvas/PPT usam estado inicial confirmado, não percurso da prévia. Layout financeiro140x80 e Bolão três botões+legendas seguem [033/DESIGN](../033-menu-aplicacao-efl/DESIGN.md), sem fidelidade visual alegada de anexos vazios. Sessão/cena em [001/T140/T142](../001-gerador-jornada/tasks.md#tarefas-menu-aplicacao-2026-10-23); nenhuma operação bancária, timer, conta, periférico, rede, QR ou Java/EFL. Implementação autorizada não encerra aceite manual pendente.

## Estados dos jogos do menu — 1.0.0.2

Controles de menu usam a cena compartilhada por item, com PNG touched ao pressionar, default ao soltar/sair e disabled para bloqueio. Mouse não pode acionar item bloqueado; foco/pressão não alteram a árvore. Normal/especial seguem os ramos reais do XML; nomes/códigos não inferem destinos. Regra canônica em 033/FR-011–FR-012.

## Submenu reconstruído e comandos — 1.0.0.3

Após abrir jogo, a janela deve reconstruir as opções normal/especial com faixa e geometria originais. Voltar/Cancelar no rodapé usam a cena compartilhada e limpam a seleção anterior; Sair encerra somente a prévia. Seleção simulada mostra a folha e seus metadados, diferenciando nomes com mesmo código. Componente desabilitado/sessão encerrada bloqueiam ativação. Verificação deve acionar controles do conteúdo atual da janela, não apenas conferir Path. Cadastro movimento.xml ausente e homologação permanecem pendentes.
