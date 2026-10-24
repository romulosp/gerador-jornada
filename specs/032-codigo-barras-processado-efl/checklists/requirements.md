# Specification Quality Checklist: Código de barras processado — blocos parametrizáveis

**Purpose**: revisar a clareza/completude dos requisitos antes da implementação por domínio, conforme GitHub Spec Kit.
**Created**: 2026-10-23
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] Requisitos orientados ao autor e à composição da ficha, sem tratar parametrização como processamento bancário real.
- [x] Texto em português e seções obrigatórias do template preenchidas, com histórias priorizadas/testes independentes.
- [x] Apresentação/inspector distintos do modelo canônico, sem domínio artificial ou implementação centralizada no formulário.
- [x] Geometria proposta e anexo ausente identificados no DESIGN, sem alegação de medidas observadas ou paridade comprovada.

## Requirement Completeness

- [x] Blocos/linhas ordenados, IDs, títulos/valores vazios, quantidades mínimas e defaults estão definidos. [FR-001/FR-002/FR-005/FR-006]
- [x] As cinco opções do combo e seus rótulos visíveis estão especificados sem ambiguidade. [US2; FR-003]
- [x] Cada valor/tipo pertence a uma linha, com independência entre blocos e preservação literal ao trocar rótulo. [FR-002–FR-004]
- [x] Localização dos combos no inspector, cancelamento e transação atômica estão explícitos; não são controles runtime da jornada. [US2; FR-006/FR-008; Assumptions]
- [x] Cópia profunda/duplicação/remapeio de IDs, reordenação, Undo/Redo e persistência têm critérios de aceitação. [US3; FR-007]
- [x] Dados ilustrativos/Unicode/pontuação e ausência de máscara/validação fiscal/consulta estão delimitados. [Edge Cases; FR-004/FR-009]
- [x] Ausência de vínculo automático com entrada e de navegação/processamento bancário está especificada. [US3; FR-009]
- [x] Overflow, acessibilidade e preservação do último estado válido estão definidos. [Edge Cases; FR-006/FR-010]

## Feature Readiness

- [x] FR-001/FR-002/FR-005/FR-006 cobertos por US1/bordas; FR-003/FR-004 por US2; FR-007–FR-010 por US3/bordas/DESIGN.
- [x] SC-001–SC-004 são mensuráveis e independentes de framework; nenhuma marcação de qualidade significa execução funcional.
- [x] Domínios 002–006 e coordenação001 têm responsabilidades e tarefas qualificadas, sem dados ou backlogs paralelos na pasta do componente.
- [x] Dependência visual pendente e premissas de UX estão declaradas; não há marcador indefinido ou obrigação de consulta/Java escondida.

## Notes

Checklist embutido de qualidade mantido por `speckit-specify`; não checklist customizado de homologação ou teste executado. Marcação avalia a escrita dos requisitos. O componente continua Draft e a imagem não foi recebida; proposta geométrica, localização dos combos e limites de escopo são explícitos. Modelo/plano/tarefas seguem os domínios responsáveis, com [roteiro manual/evidências pendentes](../../001-gerador-jornada/quickstart.md#17-codigos-barras-2026-10-23). Não iniciar oferta ou declarar fidelidade por existir documentação.

Revisão documental pelo agente em 2026-10-23: 16/16 itens satisfeitos. Revisão independente não apontou contradição na parametrização/ownership; dependências de liberação foram explicitadas na coordenadora. Links/IDs e pré-requisitos dos domínios verificados; não foi executado teste funcional ou aprovado visual do segundo anexo.