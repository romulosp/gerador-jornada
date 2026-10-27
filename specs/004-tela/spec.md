# Feature Specification: Domínio Tela

## Integração planejada — instanciação da biblioteca 038, 2026-10-25

A [038 v0.4](../038-integracao-mcp-agentes/spec.md) prevê criação manual/MCP sobre os mesmos casos de uso e estilo legado. Tela é dona das instâncias expandidas, geometria, ordem e remapeamento UUID/FileID/IDs internos; não recebe vínculo vivo com biblioteca ou dados do agente. Referências: 038/FR-002/FR-007/FR-022–FR-024/FR-032, US3 e SC-004/SC-009/SC-013. [Modelo](../038-integracao-mcp-agentes/data-model.md) e [plano](../038-integracao-mcp-agentes/plan.md) detalham a divisão. Escopo apenas planejado; nenhuma alteração de código, teste ou homologação foi realizada por este registro.

**Feature Branch**: `004-tela` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Telas 800 × 600 por padrão, bounds explícitos, UUID/FileID separados, Style, Enabled/ParentID e ordem plana de camadas. Criação/cópia/seleção e transações existentes integram os compostos; skins são resolvidas na apresentação.

**Pendências e limites:** Aparência tipada/RGBA por documento, todos os novos tipos e validação global completa das identidades/referências da expansão continuam pendentes. A existência de ParentID não significa contêiner visual/clipping completo nem árvore expandível no painel Estrutura.

**Fontes verificadas:** [model.go](../../apps/gerador-jornada/internal/screen/model.go), [hierarchy.go](../../apps/gerador-jornada/internal/screen/hierarchy.go), [edit_selection.go](../../apps/gerador-jornada/internal/application/edit_selection.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compor uma tela (Priority: P1)

Modelar telas, componentes, geometria e ordem visual com operações coesas e independentes da UI.

**Why this priority**: O desenho precisa manter as mesmas coordenadas e relações ao ser reaberto, apresentado ou convertido para Java.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** uma tela 800 × 600, **When** um componente válido é inserido e movido, **Then** seus bounds são armazenados em unidades lógicas, sem depender do zoom ou DPI da janela.

---

### User Story 2 - Alterar a composição com previsibilidade (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** uma tela com teclado vinculado e botão para si mesma, **When** ela é duplicada, **Then** tela e componentes recebem novos IDs e vínculos internos são remapeados para a cópia.

### Edge Cases

Cobrir bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Criar, renomear, ordenar e dimensionar telas; default 800 × 600, dimensões positivas.
- **FR-002**: Manter componentes com UUID, nome lógico único na tela, kind do catálogo, bounds inteiros, style e properties tipadas.
- **FR-003**: Oferecer inserir, mover, redimensionar, alinhar, duplicar, remover e ordenar camadas como operações de responsabilidade explícita.
- **FR-004**: Exigir bounds dentro da tela para executar/exportar; reduzir tela não pode reposicionar componentes silenciosamente.
- **FR-005**: Duplicar IDs e remapear vínculos internos, preservando destinos externos; histórico e transações entre módulos ficam na aplicação.

### Key Entities

Screen, Component, Bounds e Style; ordem no array é z-order do fundo para frente. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todas as operações de referência preservam inteiros e IDs corretos; zoom/DPI não alteram o modelo; duplicação não mantém IDs originais.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-002, FR-004 a FR-008, FR-013, FR-016. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [component-catalog.md](../001-gerador-jornada/contracts/component-catalog.md).
- Não implementa Fyne, coordenadas de ponteiro, pilha undo/redo, serialização, roteiro ou exportação.
- Esta especificação define o alvo; evidências de implementação e testes estão no quickstart da feature coordenadora, sem homologação implícita.

## Especificações individuais dos componentes

O detalhamento abaixo mantém este domínio como responsável pelo modelo e pelas regras compartilhadas. Cada feature tem contrato visual próprio; sua existência documental não conclui implementação ou homologação.

- [Data e hora — DateInfo](../014-data-hora/spec.md) — [DESIGN.md](../014-data-hora/DESIGN.md).
- [Dados da lotérica — LotericaInfo](../015-dados-loterica/spec.md) — [DESIGN.md](../015-dados-loterica/DESIGN.md).
- [Indicadores — InfraInfo](../016-infraestrutura/spec.md) — [DESIGN.md](../016-infraestrutura/DESIGN.md).
- [Cabeçalho informativo da tela](../017-cabecalho-tela/spec.md) — [DESIGN.md](../017-cabecalho-tela/DESIGN.md).
- [Faixa de título da página](../018-faixa-titulo/spec.md) — [DESIGN.md](../018-faixa-titulo/DESIGN.md).
- [Área de desenho — Design](../019-canvas-edicao/spec.md) — [DESIGN.md](../019-canvas-edicao/DESIGN.md).
- [Elemento editável — PageElement](../020-elemento-edicao/spec.md) — [DESIGN.md](../020-elemento-edicao/DESIGN.md).
- [Lista de telas — PageTab](../022-lista-telas/spec.md) — [DESIGN.md](../022-lista-telas/DESIGN.md).
- [Estrutura — StructureTab](../024-arvore-estrutura/spec.md) — [DESIGN.md](../024-arvore-estrutura/DESIGN.md).

## Atualização solicitada em 2026-10-18 — Seleção e composição de texto

- **FR-006**: Selecionar múltiplos componentes por Ctrl/Shift+clique, caixa de seleção e Ctrl+A; seleção é estado de interface e não deve ser persistida.
- **FR-007**: Mover uma seleção inteira preservando distâncias relativas e limitando o deslocamento pelo conjunto; desfazer o movimento inteiro em uma ação.
- **FR-008**: Persistir alinhamento horizontal do texto como esquerda, centro ou direita, independente do alinhamento geométrico do elemento; projetos anteriores sem a propriedade usam esquerda.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-006, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-007, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
3. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

## Exclusão confirmada — 2026-10-18

- **FR-009**: Remover o conjunto confirmado em uma única transação reversível, preservando identidades/ordem dos demais componentes e referências externas diagnosticáveis. Seleção, menu contextual, disponibilidade do botão e confirmação Sim/Não pertencem ao adaptador; não inserir widgets ou estado de confirmação no modelo de domínio.

**Aceitação**: dado um conjunto de dois componentes selecionados, confirmar sua exclusão remove ambos em um comando e desfazer restaura ambos; cancelar a confirmação não gera revisão ou histórico de domínio.

## Visualização individual — 2026-10-18

- **FR-010**: A composição de uma tela deve poder ser visualizada isoladamente, preservando seus componentes, dimensões, recursos e conteúdo confirmado. O escopo de visualização e os valores temporários são responsabilidade da aplicação/adaptador; a tela de domínio não recebe estado de janela ou navegação de prévia.

## Identidade de arquivo — 2026-10-18

- **FR-011**: Screen e Component têm FileID editável, independente de ID UUID e Name. Exigir nome-base ASCII alfanumérico de até 40 caracteres, início em letra e sem nomes reservados Windows; Screen exige Pnl seguido de nome. Component não exige Pnl. Nome amigável e textos continuam Unicode. Alterar FileID não troca identidades, navegação, teclado ou referências de roteiro.
- **FR-012**: Garantir FileID de telas único no projeto e de componentes único na tela, sem diferenciar caixa. Duplicação/criação pode sugerir nome livre com sufixo dentro de 40 caracteres; edição explícita inválida deve ser rejeitada atomicamente. Troca de Name não altera FileID já confirmado.

**Aceitação**: Identificação sugere PnlIdentificacao; campos de 40 caracteres são válidos e de 41 são recusados; Pnl sozinho não é tela válida. Copiar controle/tela no limite não cria colisão nem ultrapassa 40. Editar somente FileID preserva UUID e vínculos; Undo restaura o ID anterior.

## Expansão pós-v1 autorizada — 2026-10-19

- **FR-013**: Cada instância deve preservar referência de skin/version e overrides tipados presentes, inclusive fonte/face real, RGBA/alinhamento/espaçamento suportados, com precedência override → skin/estado → default comprovado. Migrar Style anterior sem perder customizações; alinhamento vazio antigo permanece esquerda. Aplicar aparência nova a antigos sem alterar bounds/UUID/FileID/conteúdo/ordem/vínculos.
- **FR-014**: Permitir ParentID opcional de contêiner na mesma tela, sem ciclos/autorreferência; manter lista plana e bounds absolutos. Mover pai desloca descendentes uma vez em transação; redimensionar pai não escala filhos implicitamente. Clip/overflow segue contrato do contêiner, sem reposicionamento silencioso.
- **FR-015**: Criar/copiar/duplicar/remover novas famílias/variantes e hierarquias com cópia profunda/remapeio de IDs/referências, inclusive opções/linhas/grupos/alvos; desfazer/refazer preserva conjunto e identidade dos demais. Report/teclado mantêm partes internas derivadas, não subcomponentes persistidos artificiais.
- **FR-016**: Expandir presets Frame 800 × 600 e composições 014–018/ValorSurpresaSuperSete usando famílias/skins comuns; valores informativos continuam estáticos, sem relógio vivo, negócio, autenticação, serviço ou periférico. Novos defaults só na criação (teclado 200 × 350); migrar nunca redimensiona instância antiga 200 × 280.

**Aceitação**: mover panel com filhos e desfazer restaura todos em uma ação, inclusive filho também selecionado sem movimento duplicado; ciclo/pai de outra tela é recusado sem publicação parcial. Abrir fixture antiga preserva bounds/IDs/overrides e recebe nova aparência onde permitido. Salvar/reabrir e duplicar preserva conteúdo com IDs remapeados; resize do pai evidencia overflow sem escalar filho. T009–T012 locais; rendering/eventos ficam em 001/T111–T116. Esta é expansão intencional da composição plana v1, não conclusão técnica ou homologação.

## Sobreposição de componentes — 2026-10-20

- **FR-018**: Oferecer `screen.PlaceComponent(componentID, referenceID string, inFront bool) error` para reposicionar somente o componente indicado na lista plana `Components`, ordenada do fundo para a frente. Com referência da mesma tela, colocar imediatamente depois (à frente) ou antes (atrás); referência vazia significa extremo frontal/fundo. Preservar a ordem relativa de todos os demais componentes. Tela nil, origem ausente, referência ausente ou autorreferência retornam erro sem mutação; posição já satisfeita é no-op. Preservar UUID/FileID, bounds, estilo, conteúdo, vínculos e ParentID. Filhos não acompanham a camada do pai: não mover subárvore nem reinterpretar a operação como ordenação de grupo. Persistir somente a ordem canônica existente, sem novos campos, versão de schema/catálogo ou restrição relacional permanente.

**Aceitação adicional (pendente):** exercitar todos os pares/direções, extremos, adjacência já satisfeita, origem em ambos os lados da referência, tela vazia/unitária e erros sem alteração parcial. Conferir estabilidade dos demais e conteúdo integral das 16 famílias, inclusive pai/filhos com camadas independentes. A aplicação confirma camada e outras propriedades numa transação reversível; formulário/modos/referência são responsabilidade de [028/FR-018](../028-propriedades/spec.md#sobreposição-de-componentes--2026-10-20), não novos dados de Screen. Os comandos antigos Frente ↑/Fundo ↓ permanecem de um passo.

Implementação em [layer.go](../../apps/gerador-jornada/internal/screen/layer.go) e regressões em layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/screen/layer_test.go`); verificação local T014 e integração [001/T127](../001-gerador-jornada/tasks.md#phase-21-sobreposição-de-componentes--2026-10-20). Testes Go não executados por política Windows, sem contorno; vet/compilação integral/build aguardam o executor principal. Manual e aceite pendentes, sem transferir PASS ou hashes históricos ao novo incremento.

## Alinhamento textual universal — 2026-10-20

- **FR-017**: Aplicar o `Style.TextAlign` existente a todas as 16 famílias capazes de apresentar texto, independentemente de conteúdo vazio, sem novas entidades ou versões de schema/catálogo. Esquerda/centro/direita são overrides textuais; **Padrão da skin** remove o override e resolve a herança existente. Preservar a conversão histórica de vazio antigo para esquerda. Alterar alinhamento nunca modifica bounds, ParentID, pais/filhos, ordem, vínculos, fundos ou áreas de interação; cada texto alinha dentro de sua própria região, não pela caixa externa de um composto.

**Aceitação adicional (pendente):** comparar antes/depois nas 16 famílias × três overrides, incluindo componentes com filhos e compostos com células/linhas/teclas; salvar/reabrir e desfazer/refazer preservam alinhamento e os demais dados. O combo conserva a mesma escolha no popup, inclusive Sem seleção. Domínio não recebe estado do popup nem widgets. Invariantes em T013 local; integração e regressões em [001/T126](../001-gerador-jornada/tasks.md#phase-20-alinhamento-textual-universal--2026-10-20). Evidência parcial informada pelo executor principal em 20/10: gofmt/vet PASS, 15 suítes Go compiladas sem execução, quatro testes Python de empacotamento PASS e tarefa VS Code de compilação Windows PASS. T013 permanece aberta: testes Go não executados pelo bloqueio prévio Windows, sem contorno; manuais Windows/DPI/PowerPoint e homologação continuam pendentes.

<a id="codigos-barras-2026-10-23"></a>
## Incremento planejado — instâncias de códigos de barras (2026-10-23)

- **FR-019**: Compor [031 — entrada](../031-codigo-barras-efl/spec.md) e [032 — processado](../032-codigo-barras-processado-efl/spec.md) como instâncias únicas de Tela, consumindo propriedades tipadas do Catálogo. Partes internas não são componentes soltos. Criação aplica defaults somente ao novo objeto; snapshot copia profundamente mantendo IDs, duplicação gera novos IDs de composto/blocos/linhas e remapeia vínculos compatíveis, inclusive teclado. Edição de blocos/linhas/seletores/valores deve ser atômica e reversível pelo histórico da aplicação, preservando IDs dos itens não removidos, bounds, z-order, ParentID e demais propriedades não alteradas.

**Aceitação**: mover/duplicar/desfazer uma entrada não cria oito componentes; alterar um bloco processado não afeta outros blocos ou uma cópia. Duplicar tela gera IDs internos novos e remapeia o teclado para a cópia da entrada na mesma tela. Remover última linha/bloco é recusado sem perda parcial. Seleção/foco/popup permanecem fora de Screen. Trabalho em [T015–T016](tasks.md#tarefas-codigos-barras-2026-10-23); wire em 002 e validações de propriedades/matemática em 003, sem duplicação.

<a id="menu-aplicacao-2026-10-23"></a>
## Incremento autorizado — instância de Menu aplicação

- **FR-020**: Compor [033](../033-menu-aplicacao-efl/spec.md) como uma única instância com árvore tipada do Catálogo. Confirmar alterações da árvore e escolha inicial atomicamente, com Undo/Redo, mantendo IDs dos nós conservados; aceitar menus vazios. Snapshot preserva identidades por cópia profunda; duplicar componente/tela gera novos UUIDs para instância, todos os submenus e itens e remapeia InitialMenuID. Destinos explícitos entre telas seguem política de Jornada, não são inferidos da categoria.

**Aceitação**: mover/redimensionar/reordenar não separa botões internos nem muda árvore. Descartar/erro não publica fragmentos; duplicação não compartilha dados nem reaproveita UUIDs. Remoção do submenu inicial requer substituto autorado válido na mesma transação. Bounds/estilo/ParentID/FileID continuam comuns; caminho local/foco ficam fora de Screen. [T017](tasks.md#tarefas-menu-aplicacao-2026-10-23) é backlog único, sem repetir 004/T015–T016 de barcode.

## Parcela instalada integrada à 038

038/T021–T023 compartilham CreateScreen/ConfigureScreen/ComponentChanges/ValidateComponentEdit com as operações manuais, mantendo contexto de fundo e validação de nome/FileID/geometria. Configuração incompatível recusa sem mover controles. Texto novo usa pontos de código; histórico inalterado é preservado. Cópias de tela/seleção remapeiam opções/linhas/partes e selecionado do rádio, além dos IDs principais. O caso de uso de presets aplica decoração recebida do adaptador; não há pós-processamento exclusivo da GUI. Isso não fecha todo o backlog de Tela ou aprovação visual. Evidência em 038/quickstart.
