# Data Model: Domínio Demonstração

**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/scenario`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Scenario contém id, name, startScreenId e steps. Step contém id, kind, screenId, componentId, point, text e durationMs, com nulos e campos obrigatórios por kind conforme contrato. Ocorrência deriva de cada visita e acumula tempo; cursor inicia em (0,0).

## Invariantes

IDs de steps únicos no projeto; duração inteira 1..600000 ms; etapas são sequenciais sem sobreposição. moveCursor usa ponto ou componente, nunca ambos. typeText exige textField e substitui valor runtime. navigate exige button com binding válido; todas as etapas usam screenId corrente.

## Limites e relacionamento com outros domínios

Recebe snapshot resolvido de telas/componentes e resolução Journey, preparado pela aplicação. Interpretador devolve eventos puros. Player da aplicação controla relógio e pausa; adaptadores desktop/PPT apenas consomem eventos.

Não grava uso do mouse, não cria macros, não contém lógica condicional e não emite XML PPT ou eventos Fyne.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: intervalo temporal, campos por kind, alvo removido, tela errada, A-B-A, auto-navegação, typeText vazio, pausa e reinício. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-017 a FR-021 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Nenhuma persistência ou integração foi executada nesta entrega.
