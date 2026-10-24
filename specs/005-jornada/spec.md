# Feature Specification: Domínio Jornada

**Feature Branch**: `005-jornada` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Bindings none/navigate de botões, tela inicial, destinos explícitos de folhas de menu e execução de prévia independente. Menu sem destino apenas informa seleção simulada; Prévia Tela Atual não navega.

**Pendências e limites:** A aplicação atual não possui NavigationGraph/Graph nem warning completo de inalcançáveis. Integridade global e índice de impacto de destinos internos de menu ainda precisam ser completados; o runtime verifica o destino ao tentar navegar.

**Fontes verificadas:** [model.go](../../apps/gerador-jornada/internal/journey/model.go), [document.go](../../apps/gerador-jornada/internal/application/document.go), [preview_compounds.go](../../apps/gerador-jornada/internal/application/preview_compounds.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Definir o fluxo entre telas (Priority: P1)

Modelar tela inicial e navegação de botões como dados explícitos, compartilhados pela prévia e pelas saídas.

**Why this priority**: A navegação funcional da prévia legada não é transportada aos listeners Java; a nova versão precisa de um contrato comum.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** telas A e B e um botão de A, **When** a ação navigate aponta a B, **Then** o resolvedor retorna B pelo ID; ausência de ação equivale a none.

---

### User Story 2 - Identificar destinos inválidos (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um destino removido, **When** a jornada é validada, **Then** o diagnóstico identifica origem e destino ausente; o grafo inválido não é silenciosamente corrigido.

### Edge Cases

Cobrir none implícito, ciclo, auto-navegação, fonte não botão, destino apagado, inicial ausente e binding duplicado. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Manter startScreenId opcional no rascunho e obrigatório/existente para execução.
- **FR-002**: Permitir somente ações none e navigate associadas ao click de button, no máximo uma por fonte.
- **FR-003**: Exigir destino existente em navigate e destino nulo em none; fonte pertence à tela informada.
- **FR-004**: Permitir ciclos e auto-navegação; emitir warning de tela inalcançável sem excluí-la.
- **FR-005**: Expor resolução e validação pura, sem cliques de UI, chamadas externas, condições baseadas em campos ou regras de negócio.

### Key Entities

Journey e ButtonAction com origem, kind e destino por IDs. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os vínculos do conjunto A → B → C → A resolvem corretamente; destino apagado sempre produz diagnóstico de origem; ciclos válidos não são recusados.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-011 a FR-013, FR-014, FR-023. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md).
- Não armazena valores de campo, histórico de navegação, widgets, slides ou listeners Java.
- Evidências da base implementada estão no backlog/quickstart da coordenadora; esta atualização é documental e não executa novos testes.

## Expansão pós-v1 autorizada — 2026-10-19

- **FR-006**: Resolver capacidades de novas famílias/variantes por índice tipado preparado pela aplicação. Button de duas legendas conserva none/navigate por UUID; Toggle/Combo/Table/Report/teclados não ganham navegação por semelhança visual. Estados e eventos locais não executam serviço/periférico/negócio externo.
- **FR-007**: Preservar/resetar estados iniciais de cada família em cada visita, mantendo foco/cursor/seleção corrente/scroll/modal na sessão runtime isolada; alvo de teclado/grupo deve ser compatível e da mesma tela. Migração de aparência não muda inicial, destinos, bindings ou identidade dos componentes.

**Aceitação**: fixture A→B→A com novas famílias restaura valores iniciais e mantém IDs/rotas; twoLabel navega como button, toggle não vira navigate; alvo inválido/externo é diagnosticado sem substituição automática. Prévia Tela Atual continua sem navegação. T009–T010 locais cuidam da validação/resolução; 001/T114–T115 cuida da sessão/pintura. Java permanece pausado e a expansão não homologa seu perfil.

<a id="codigos-barras-2026-10-23"></a>
## Limite de domínio — códigos de barras (2026-10-23)

- **FR-008**: Conservar `none/navigate` explícitos exclusivamente em botões ao introduzir [031](../031-codigo-barras-efl/spec.md) e [032](../032-codigo-barras-processado-efl/spec.md). Validação/conclusão/reinício de entrada e escolhas de rótulo não são ações de Jornada, não preenchem outra tela e não disparam navegação. Não armazenar códigos, derivados ou fichas processadas no grafo nem acrescentar condicionais baseadas em DV nesta entrega.

**Aceitação**: entrada aceita/recusada ou troca CPF/CNPJ conserva bindings/tela corrente; só botão com ação explícita navega. A→B→A reinicia entrada vazia e mantém configuração processada; Prévia Tela Atual não navega. A conclusão “permite prosseguir” na 031 significa liberação local da entrada, não motor de navegação condicional. Invariantes em [T011](tasks.md#tarefas-codigos-barras-2026-10-23); runtime é 001/T135, regras são 003/T019.

<a id="menu-aplicacao-2026-10-23"></a>
## Incremento autorizado — destinos de folhas do menu

- **FR-009**: Reconhecer destino de tela explicitamente autorado em folha de [033](../033-menu-aplicacao-efl/spec.md), validando existência, impacto de remoção e alcançabilidade pelo mesmo índice de referências de Jornada. Navegação local de submenu/Voltar/Cancelar não é rota; folha sem destino apenas informa seleção simulada, nunca infere rota por Label/Shortcut/MovementCode. Item com submenu não admite destino. Habilitação da instância e do caminho deve impedir ativação desabilitada.

**Aceitação**: folha habilitada com destino válido pode navegar somente na prévia completa; destino inexistente é localizado por tela/componente/item e não recebe substituto. Ciclos/autorrotas explícitos seguem regras existentes. Remover tela informa referências das folhas inclusive em ramos desabilitados; Prévia Tela Atual conserva interação local sem seguir destinos. [T012](tasks.md#tarefas-menu-aplicacao-2026-10-23) valida referências, não implementa runtime. FR-008 continua exclusivo aos códigos de barras; esta extensão não os torna origens de navigate.
