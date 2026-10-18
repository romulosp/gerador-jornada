# Data Model: Domínio Tela

**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/screen`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Screen contém id, name, size, background e components. Component contém id, name, kind, bounds, style e properties. Bounds contém x, y, width e height; Style usa os campos comuns do contrato do catálogo. Não há árvore de componentes ou agrupamento persistido na v1.

## Invariantes

x/y inteiros não negativos; width/height inteiros positivos; x+width e y+height dentro da tela para validação executável. Nomes lógicos únicos na tela. Ordem do array é camada. Dimensões de drop são ajustadas na aplicação; regras recusam estado inválido de operação confirmada.

## Limites e relacionamento com outros domínios

Importa catalog para tipos e validação de família. Usa valores/IDs, não widgets. Aplicação agrupa alterações de tela e referências de Journey/Scenario num comando reversível único. Presentation consome geometria sem modificar entidades.

Não implementa Fyne, coordenadas de ponteiro, pilha undo/redo, serialização, roteiro ou exportação.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-002, FR-004 a FR-008, FR-013, FR-016 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Nenhuma persistência ou integração foi executada nesta entrega.
