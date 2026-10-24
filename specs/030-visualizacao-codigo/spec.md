# Feature Specification: Código e diagnóstico — Code

**Feature Branch**: `030-visualizacao-codigo` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Interface de código planejada; Java/EFL suspenso; homologação pendente
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Code/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** A ação Java/EFL está visível e desabilitada; o gerador interno e os testes estruturais de fontes Java foram preservados. A revisão visual não habilita exportação nem executa Java.

**Pendências e limites:** A interface completa de arquivos/código e diagnóstico com perfil EFL homologado continua planejada. Sem biblioteca real não há comprovação de integração Java/EFL.

**Fontes verificadas:** [editor.go](../../apps/gerador-jornada/internal/adapter/desktop/editor.go), [sources.go](../../apps/gerador-jornada/internal/adapter/javaefl/sources.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero inspecionar os fontes por tela e entender impedimentos antes da geração Java.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve apresentar arquivos organizados por tela, fábrica, navegação e integração quando houver perfil homologado.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve exibir diagnóstico claro de perfil não verificado enquanto faltarem bibliotecas reais.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve preservar o alvo Java 1.6 e a documentação em português, sem permitir executar conteúdo editado como template.

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

- Sem JARs EFL, nomes em colisão, pacote inválido, arquivo extenso e código Unicode. O legado gerava templates inconsistentes; essa falha não deve ser reproduzida.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Apresentar arquivos organizados por tela, fábrica, navegação e integração quando houver perfil homologado.
- **FR-002**: Exibir diagnóstico claro de perfil não verificado enquanto faltarem bibliotecas reais.
- **FR-003**: Preservar o alvo Java 1.6 e a documentação em português, sem permitir executar conteúdo editado como template.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Arquivos gerados e diagnósticos derivados do snapshot; seleção de arquivo e rolagem transitórias.
- **Estado de interação**: Perfil indisponível; validação recusada; fontes disponíveis; arquivo selecionado; exportação cancelada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `030-visualizacao-codigo/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Suspensão temporária de acesso

- **FR-008**: Enquanto a capacidade Java/EFL estiver suspensa por decisão do usuário, sua ação permanece visível e desabilitada. Não remover o gerador nem desbloquear perfil sem homologação.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

### Escopo vigente — integração indisponível retirada

A decisão posterior [001/FR-059](../001-gerador-jornada/spec.md#dependencias-ausentes-fr-059) retira Java/EFL do produto vigente: sem ação, callback, adaptador gerador ou ferramentas Java 6. A visualização Java descrita anteriormente não é funcionalidade entregue nem ação apenas suspensa. Eventual reintrodução/visualização exige change com a biblioteca real, implementação e validação próprias. Os documentos e IDs antigos continuam legíveis, sem integração executável. A atualização é de escopo, não declaração de implementação ou aprovação desta tela; acompanhar [001/T133](../001-gerador-jornada/tasks.md#t133-dependencias-ausentes).
