# Specification Quality Checklist: Código de barras — entrada validada

**Purpose**: revisar a qualidade documental antes da implementação por domínio, conforme GitHub Spec Kit.
**Created**: 2026-10-23
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] Comportamento/valor de uso separados de caminhos de código e desenho técnico, mantidos no DESIGN e nos planos/contratos dos domínios.
- [x] Texto em português, com atores, fluxo principal e responsabilidade da entrada definidos.
- [x] Todas as seções obrigatórias do template preenchidas, com histórias priorizadas e testes independentes.
- [x] Decomposição por domínio explícita, sem criar domínio de boleto ou backlog concorrente na pasta do componente.

## Requirement Completeness

- [x] Oito campos visuais distinguidos dos cinco blocos lógicos, com comprimentos e momento correto de validação definidos. [FR-001–FR-004]
- [x] DVs locais/geral, conversão44/47, zeros e rejeição48 têm regras e fixtures localizáveis. [FR-002–FR-005; contrato compartilhado]
- [x] Captura válida/inválida, modo de leitura, normalização, cancelamento e texto/ação exatos do diálogo estão especificados. [US2; FR-005–FR-007]
- [x] Recuperação manual de par e erro geral são distinguíveis da limpeza por captura integral, sem avanço inválido. [US1; FR-002–FR-006]
- [x] Banco/valor/fator/CAIXA/SIGCB e ausência de dados derivados inválidos estão delimitados; SIGCB não é anunciado homologado. [US3; FR-008; Assumptions]
- [x] Escopo de hardware/PIX/ciclo de vencimento/consulta/pagamento e navegação implícita está delimitado. [Edge Cases; FR-012]
- [x] Reset, habilitação, teclado numérico e isolamento entre instâncias/documento estão especificados. [FR-007/FR-009]
- [x] Ausência de imagem e das classes de módulo aparece como limitação/gate, sem marcador de esclarecimento indefinido ou evidência inventada. [Assumptions]

## Feature Readiness

- [x] FR-001–FR-004/FR-007 cobertos por US1 e bordas; FR-005–FR-007 por US2; FR-008/FR-012 por US3; FR-009–FR-011 por US4/bordas.
- [x] SC-001–SC-005 são mensuráveis, independentes de framework e definidos como resultados esperados, não como testes aprovados.
- [x] Modelos/contratos/planos/tarefas nos domínios 002–006 e coordenação001 são rastreáveis por IDs qualificados, com responsabilidades sem duplicação.
- [x] Limites da saída PPT/roteiro e ausência de Java estão explícitos; checklists não autorizam homologação ou implementação fora dos domínios.

## Notes

Checklist embutido de qualidade mantido pelo fluxo `speckit-specify`, não checklist customizado de testes de implementação. Marcação significa revisão da escrita dos requisitos, não comportamento executado. Esta spec permanece Draft; implementação, teste manual, anexos, comparação com módulos reais e aceite são gates separados no [quickstart](../../001-gerador-jornada/quickstart.md#17-codigos-barras-2026-10-23). Formalizar análise `spec/plan/tasks` nos domínios responsáveis antes de iniciar o código; 031 contém somente o pacote documental solicitado.

Revisão documental pelo agente em 2026-10-23: 16/16 itens satisfeitos. Revisão independente identificou cobertura de colagem44 e dependências de oferta; ambos os pontos foram explicitados antes do fechamento deste checklist. Pré-requisitos Spec Kit dos seis artefatos coordenador/domínios e links/IDs verificados, sem testes Go ou aprovação funcional.