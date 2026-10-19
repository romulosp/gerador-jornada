# Data Model: Domínio Jornada

**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/journey`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Journey contém startScreenId e actions. ButtonAction contém sourceScreenId, sourceComponentId, kind e targetScreenId. IDs são referências, nunca índices ou nomes. O grafo deriva da lista de ações; não persistir uma segunda cópia divergente.

## Invariantes

Uma ação por botão; fonte deve ser button da tela informada. none implica destino nulo; navigate implica destino existente. Não exigir grafo acíclico. Sem inicial, salvar rascunho é permitido; executar não.

## Limites e relacionamento com outros domínios

Recebe índice imutável de referências permitido pela aplicação, sem importar Screen ou adaptadores. Scenario usa resolução de rota para interpretar navigate. Desktop e Java executam a ação; PPT traduz em links.

Não armazena valores de campo, histórico de navegação, widgets, slides ou listeners Java.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: none implícito, ciclo, auto-navegação, fonte não botão, destino apagado, inicial ausente e binding duplicado. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-011 a FR-013, FR-014, FR-023 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.


## Correspondência com o código atual

Binding none/navigate, resolução por ID e validação de origem/destino estão implementados. O warning de tela inalcançável previsto em FR-004 ainda não é emitido; permanece em T006. Diagnóstico clicável/localização completa também depende de 001-gerador-jornada/T030.
