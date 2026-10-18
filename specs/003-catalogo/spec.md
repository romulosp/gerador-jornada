# Feature Specification: Domínio Catálogo

**Feature Branch**: `003-catalogo` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — domínio em implementação; homologação e aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Escolher um componente conhecido (Priority: P1)

Definir as sete famílias, seus defaults, propriedades tipadas e capacidades sem executar código do projeto.

**Why this priority**: Uma mesma definição deve orientar paleta, inspector e validação, evitando interpretações diferentes entre editor e saídas.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** o catálogo versão 1, **When** a paleta solicita componentes disponíveis, **Then** recebe exatamente botão, rótulo, título, campo, seleção, relatório e teclado numérico, com identificadores estáveis.

---

### User Story 2 - Editar apenas propriedades válidas (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um campo de texto, **When** se tenta atribuir uma propriedade de outra família ou enum desconhecido, **Then** a operação é recusada sem alterar o componente; propriedades suportadas ficam preservadas.

### Edge Cases

Cobrir sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Disponibilizar exatamente os kinds button, label, title, textField, radio, report e numericKeyboard no catálogo inicial.
- **FR-002**: Definir defaults, tipos, enums e estilos conforme o contrato de catálogo; versão inicial igual a 1.
- **FR-003**: Validar propriedades por structs de família, sem map[string]any atravessando os casos de uso.
- **FR-004**: Expor definições para paleta/inspector sem depender de Fyne, JSON, XML ou templates Java.
- **FR-005**: Tratar capacidade Java como mapeamento de perfil do adaptador; nome EFL observado não comprova API nem permite ignorar propriedade alterada.

### Key Entities

ComponentDefinition, Kind, PropertyDefinition e structs ButtonProperties, LabelProperties, TitleProperties, TextFieldProperties, RadioProperties, ReportProperties e NumericKeyboardProperties. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: As sete famílias têm defaults válidos e round-trip sem perda; 100% dos enums/propriedades fora do contrato são diagnosticados nos testes de referência.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-003, FR-005, FR-015, FR-022; NFR-002. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [component-catalog.md](../001-gerador-jornada/contracts/component-catalog.md).
- Não implementa widget, drag/drop, perfis EFL externos, plugins executáveis ou templates por projeto.
- Esta especificação define o alvo; evidências de implementação e testes estão no quickstart da feature coordenadora, sem homologação implícita.

## Especificações individuais dos componentes

O detalhamento abaixo mantém este domínio como responsável pelo modelo e pelas regras compartilhadas. Cada feature tem contrato visual próprio; sua existência documental não conclui implementação ou homologação.

- [Botão — EFLButton](../007-botao-efl/spec.md) — [DESIGN.md](../007-botao-efl/DESIGN.md).
- [Rótulo — EFLLabel](../008-rotulo-efl/spec.md) — [DESIGN.md](../008-rotulo-efl/DESIGN.md).
- [Título — EFLLabelTitulo1](../009-titulo-efl/spec.md) — [DESIGN.md](../009-titulo-efl/DESIGN.md).
- [Campo de texto — EFLTextField](../010-campo-texto-efl/spec.md) — [DESIGN.md](../010-campo-texto-efl/DESIGN.md).
- [Seleção — EFLRadio](../011-selecao-efl/spec.md) — [DESIGN.md](../011-selecao-efl/DESIGN.md).
- [Relatório/lista — EFLReportEvent](../012-relatorio-efl/spec.md) — [DESIGN.md](../012-relatorio-efl/DESIGN.md).
- [Teclado numérico — EFLVirtualKeyboardNumeric](../013-teclado-numerico-efl/spec.md) — [DESIGN.md](../013-teclado-numerico-efl/DESIGN.md).
- [Paleta — ComponentTab](../023-paleta-componentes/spec.md) — [DESIGN.md](../023-paleta-componentes/DESIGN.md).
