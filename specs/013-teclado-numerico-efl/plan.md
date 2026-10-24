# Implementation Plan: teclado numérico Java/EFL

**Date**: 2026-10-19 | **Spec**: [spec.md](spec.md) | **Status**: plano pós-v1 autorizado; implementação/verificação e homologação pendentes.

## Escopo e contratos

FR-008–FR-009 complementam requisitos iniciais: 14 teclas com layout Java 200 × 350 somente para novas inserções; 200 × 280 e outros bounds antigos são preservados na migração. [DESIGN.md](DESIGN.md) registra posições, dimensões e código VK_BACK_SPACE emitido pela legenda LIMPAR. Não confundir intenção da legenda, evento emitido e efeito no campo; clear-all atual exige contrato distinto até comprovação do tratamento Java.

Tipos/formatos/alvos/cópia/defaults pertencem a 003/T013 (campos em T012); instâncias/comandos a 004/T009–T012; wire/migração a 002/T009–T012. Fontes/skins/importador são 001/T108–T110. Cena/board/inspector/sessão/preview/PPT são 001/T111–T116. Esta feature não reimplementa esses mecanismos: mantém fixture, verificação do slice e aceitação. Mesmo Go/Fyne da coordenadora, sem novo domínio/módulo/JVM/runtime.

## Sequência e validação

1. Inventariar contrato de tecla/estado/fonte/skin e efeito de TAB/setas/VK_BACK_SPACE no alvo; registrar lacunas na matriz C15, sem inventar recurso/delegate ausente. Alfanumérico 800 × 200/calculadora continuam C16–C17 do catálogo, não resultado deste piloto.
2. Fixture de dois campos compatíveis, teclado vinculado, texto Unicode/cursor, habilitado/desabilitado, alvo ausente/removido e arquivo v1/v2 de teclado 200 × 280 com overrides/IDs. Novo default e round-trip antigo são cenários distintos.
3. Verificar integrações dos donos únicos e teste local da fixture com cena/entrada/recorte/PPT editável; comportamento de apagar/limpar deve corresponder ao contrato identificado, sem manipular documento pela prévia.
4. Homologar manualmente Windows/DPI 100/150/200% e PowerPoint 2016/365 com referência independente, revisão/hash e executor/observado; exigir aceite da revisão testada. Sem evidência, estado permanece pendente/bloqueado.

Manter Java pausado; sem serviços/periféricos, importação `.page`, segredo real, screenshot como representação funcional ou homologação de outros teclados por similaridade. Tarefas locais em [tasks.md](tasks.md); entrega/evidências futuras a cargo do executor na coordenadora, não preenchidas por este plano.