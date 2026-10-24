# Feature Specification: Teclado numérico — EFLVirtualKeyboardNumeric

**Feature Branch**: `013-teclado-numerico-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/ListComponents/EFLVirtualKeyboardNumeric/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Novo teclado 200 × 350 com 14 teclas Java: 7–9, 4–6, 1–3, < / 0 / >, TAB e LIMPAR. Bounds antigos são preservados e as partes escaladas. LIMPAR emite Backspace; setas respeitam cursor e TAB percorre campos habilitados na prévia.

**Pendências e limites:** Teclado alfanumérico e variantes/overrides tipados de toda a expansão não estão disponíveis. O alvo barcodeInput existe, mas a matriz completa de bloqueio/avanço por par com TAB/setas ainda requer validação própria. Não tratar LIMPAR como clear-all.

**Fontes verificadas:** [legacy_parts.go](../../apps/gerador-jornada/internal/presentation/legacy_parts.go), [preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go), [legacy_preview_test.go](../../apps/gerador-jornada/internal/adapter/desktop/legacy_preview_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero compor uma entrada numérica com um teclado vinculado a um campo da tela.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve arrastar o teclado, redimensioná-lo e escolher um campo de destino da mesma tela.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve digitar números, apagar e limpar o campo na prévia sem alterar seus valores iniciais persistidos.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve reproduzir o arranjo legado 7–9, 4–6, 1–3, controles de cursor e 0, TAB e LIMPAR; explicitar diferenças enquanto a versão atual ainda usa doze teclas.

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

Como autor, quero salvar, reabrir e exportar a composição sem perda ou alteração implícita dos dados confirmados.

**Why this priority**: O mesmo trabalho precisa continuar coerente ao mudar de contexto.
**Independent Test**: Usar uma fixture com texto acentuado e geometria conhecida, executar o percurso abaixo e comparar os dados antes/depois.

**Acceptance Scenarios**:

1. **Given** uma configuração confirmada, **When** o autor executa o percurso de preservação, **Then** o sistema deve salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
2. **Given** uma operação cancelada ou recusada, **When** o autor retorna ao trabalho, **Then** nenhum dado confirmado foi perdido nem parcialmente substituído.

### Edge Cases

- Alvo fora da tela, alvo removido, campo vazio ao apagar, Unicode, cursor no início/fim e TAB sem próximo campo. No PPT as teclas são formas; a digitação vem do roteiro.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Arrastar o teclado, redimensioná-lo e escolher um campo de destino da mesma tela.
- **FR-002**: Digitar números, apagar e limpar o campo na prévia sem alterar seus valores iniciais persistidos.
- **FR-003**: Reproduzir o arranjo legado 7–9, 4–6, 1–3, controles de cursor e 0, TAB e LIMPAR; explicitar diferenças enquanto a versão atual ainda usa doze teclas.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Campo de destino, posição, tamanho, fonte, cor e skins das teclas normais e auxiliares.
- **Estado de interação**: Sem alvo; alvo válido; alvo removido; tecla pressionada; foco; selecionado no editor.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `013-teclado-numerico-efl/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 003-catalogo.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## ID das instâncias

Toda instância visual desta responsabilidade possui ID de arquivo conforme [Tela](../004-tela/spec.md#identidade-de-arquivo--2026-10-18) e [Catálogo](../003-catalogo/spec.md#atributo-compartilhado-das-instâncias). Campo comum no [painel de propriedades](../028-propriedades/spec.md#campo-id-comum), sem prefixo Pnl obrigatório para componentes e sem substituir UUIDs ou conteúdo visual. Em conjuntos prontos, a regra vale para cada controle efetivamente criado.

## Expansão Java/EFL autorizada — 2026-10-19

- **FR-008**: Novo default numérico 200 × 350, com 14 teclas Java e posições internas comprovadas; aplicar somente a novas inserções, mantendo bounds/IDs/vínculos/overrides dos projetos antigos (inclusive 200 × 280). Reconciliar comando efetivo: legenda LIMPAR Java emite VK_BACK_SPACE, não prova clear-all; conservar requisito de limpar do editor como ação local explícita distinta até evidência do tratamento, sem chamar a divergência de paridade homologada.
- **FR-009**: Integrar foco/cursor/TAB e alvos habilitados/compatíveis da mesma tela com cena/skin comum; estados são runtime isolado. AlphaKeyboard 800 × 200 integra o Catálogo e não duplica esta família. A pintura é local Go; referências a delegates externos ausentes foram removidas conforme [001/FR-059](../001-gerador-jornada/spec.md#dependencias-ausentes-fr-059). Recursos/métricas e limitações reais continuam verificáveis.
- **FR-010**: Oferecer somente Padrão/numeric/alpha. Calculadora sem layout não integra oferta, helper de composição, desenho, hit areas ou comandos. Conservar leitura/valor em arquivos antigos no salvamento neutro, sem converter para numeric; somente escolha explícita substitui a variante. Reintrodução exige change com implementação real. Verificações em [001/T133](../001-gerador-jornada/tasks.md#t133-dependencias-ausentes).

**Aceitação**: fixture com dois campos comprova layout 7–9/4–6/1–3, setas/0/TAB/LIMPAR, foco/cursor/limites e reset sem alterar valores iniciais; comparar comando e efeito observado, não inferir comportamento da legenda. Salvar/reabrir arquivo antigo não causa resize; novo teclado usa 200 × 350 e PPT mantém teclas/textos editáveis, sem interação real no PowerPoint. [plan.md](plan.md)/[tasks.md](tasks.md) locais detalham verificação; modelo 003/T013, runtime 001/T114–T115 e PPT T116. Sem reativar Java ou afirmar homologação.

## Alvo composto de código de barras — planejamento 2026-10-23

- **FR-011**: Aceitar [031 — entrada de código de barras](../031-codigo-barras-efl/spec.md) como alvo compatível na mesma tela, escrevendo somente no subcampo ativo da sessão e compartilhando seus limites/validações. TAB não ultrapassa par inválido, apagar atua no campo ativo e desabilitado não recebe entrada. Não criar oito vínculos persistidos nem transformar o processado 032, somente leitura, em alvo de digitação.

**Aceitação**: teclado físico e numérico levam ao mesmo bloqueio/avanço para a fixture F01 e seus DVs adulterados; duplicar tela remapeia alvo ao composto copiado. Modelo/capacidade em 003/T018, composição em 004/T015 e runtime em [001/T135](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23), sem duplicar algoritmo ou tarefas locais de teclado. Incremento parcialmente implementado; cobertura restante e homologação seguem pendentes no estado verificado acima.
