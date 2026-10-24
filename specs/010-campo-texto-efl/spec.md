# Feature Specification: Campo de texto — EFLTextField

**Feature Branch**: `010-campo-texto-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/ListComponents/EFLTextField/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Campo 150 × 25, textField/type1: Arial 12, fundo branco/texto preto; foco azul #10508E/texto branco. legacyEntry usa captura Fyne transparente sobre a cena, preservando cursor, seleção e clipboard.

**Pendências e limites:** Os formatos existentes validam intenção/comprimento/data; não implementam um motor incremental de máscaras nem consulta/validação fiscal. TextArea/PasswordField continuam tipos planejados, não variantes disponíveis deste campo.

**Fontes verificadas:** [legacy_entry.go](../../apps/gerador-jornada/internal/adapter/desktop/legacy_entry.go), [preview.go](../../apps/gerador-jornada/internal/application/preview.go), [legacy_preview_test.go](../../apps/gerador-jornada/internal/adapter/desktop/legacy_preview_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero desenhar e experimentar entradas de texto com intenção de formato definida.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve editar valor inicial e escolher texto, número, celular, CEP, CPF, data, DDD ou NIS.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve digitar na prévia sem gravar o valor runtime no projeto; reiniciar o valor a cada visita.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve vincular um teclado numérico da mesma tela e demonstrar digitação por etapas explícitas.

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

- Valor vazio, colagem com caracteres inválidos em campo numérico, teclado sem alvo e formatos sem validação fiscal. Formato não significa consulta de CPF ou CEP.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Editar valor inicial e escolher texto, número, celular, CEP, CPF, data, DDD ou NIS.
- **FR-002**: Digitar na prévia sem gravar o valor runtime no projeto; reiniciar o valor a cada visita.
- **FR-003**: Vincular um teclado numérico da mesma tela e demonstrar digitação por etapas explícitas.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Salvar e reabrir o conteúdo e a aparência confirmados; exportar a representação visual sem substituir textos por screenshots.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Valor inicial, formato, posição, dimensões, fonte, texto, fundo e borda.
- **Estado de interação**: Vazio; preenchido; foco; entrada recusada; selecionado no editor.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `010-campo-texto-efl/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 003-catalogo.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## ID das instâncias

Toda instância visual desta responsabilidade possui ID de arquivo conforme [Tela](../004-tela/spec.md#identidade-de-arquivo--2026-10-18) e [Catálogo](../003-catalogo/spec.md#atributo-compartilhado-das-instâncias). Campo comum no [painel de propriedades](../028-propriedades/spec.md#campo-id-comum), sem prefixo Pnl obrigatório para componentes e sem substituir UUIDs ou conteúdo visual. Em conjuntos prontos, a regra vale para cada controle efetivamente criado.

## Expansão Java/EFL autorizada — 2026-10-19

- **FR-008**: Ampliar documentos/máscaras de TextField para todas regras locais inventariadas, preservando formatos atuais e tipagem. Integrar foco/cursor/teclados com TextArea e PasswordField novos no domínio, sem confundi-los com variantes de textField ou anunciar validação fiscal/serviço pelo nome. Fontes/skins/estados seguem evidência Java, não widget Fyne; dependência/delegate ausente bloqueia fidelidade.

**Aceitação**: fixture por documento/estado demonstra edição/colagem/limites/cursor/TAB e reset isolados, máscara comprovada e diagnóstico de capacidade ausente; preservar valor inicial, IDs, bounds e overrides no round-trip. Modelo [003/T012](../003-catalogo/tasks.md), teclados 003/T013/013, runtime/cena/PPT 001/T111–T116, manuais/aceite T120–T122. Matriz C06–C08/C26; não usar credenciais reais e não reativar Java.
