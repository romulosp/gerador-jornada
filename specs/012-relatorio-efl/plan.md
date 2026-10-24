# Implementation Plan: piloto Relatório Java/EFL

**Date**: 2026-10-19 | **Spec**: [spec.md](spec.md) | **Status**: plano pós-v1 autorizado, implementação/verificação e homologação pendentes.

## Escopo e responsabilidades

FR-008–FR-009: Event é o slice piloto; todas as oito classes Report continuam no escopo total. Esta feature detalha contrato/fixture/aceitação do relatório, sem criar novo domínio ou duplicar implementação. Mesma toolchain Go/Fyne da coordenadora; regras puras, recursos offline e texto PPT editável, Java pausado.

| Trabalho | Dono único |
| --- | --- |
| Inventário/fonte/skin | 001/T108–T110 e 003/T009, matriz C11–C14 em [research.md](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19) |
| Modelo/variantes/linhas/cópia/validação | 003/T014; partes Panel necessárias em 003/T015 |
| Instância/comandos/hierarquia | 004/T009–T012 |
| Migração/recursos/round-trip | 002/T009–T012 |
| Cena/métricas/desktop/inspector/prévia/PPT | 001/T111–T116 |
| Galeria/regressão/entrega | 001/T117–T118/T123/T125 |
| Manual/aprovação | 001/T120–T122, com recorte do piloto em tasks local |

## Contrato visual e sequência

1. Reconciliar [DESIGN.md](DESIGN.md) por parte com Java/XML/React identificados: topo 26/rodapé 32 comprovados por constantes; título 18 é React, entradas Java 20/26/custom; não assumir equivalência. Registrar método/skin/estado, métricas, recorte e bloqueios na matriz existente.
2. Definir fixture Event mínima, com título/linhas Unicode, IDs estáveis, bounds/overrides, estado vazio, seleção inicial e excedente. Comparação precisa de referência independente do Java/legado; screenshot novo não é oracle.
3. Consumir implementações dos donos acima para editor→prévia→save/reopen→PPT, verificando partes separadas e editabilidade. Os paths atuais são catalog/model.go, presentation/scene.go, application/preview.go e adapters desktop/projectzip/pptx; não criar segunda cópia dessas responsabilidades aqui.
4. Repetir por variante/estado. Resultado piloto não aprova Report2/3/Event2/Notificacao/SelecaoMultipla/Internal.

## Teste e gates

Automatizado específico do piloto verifica a fixture integrada (tarefa local T003), sem substituir testes de domínio/importador/codec. Recursos/delegates/referências ausentes bloqueiam fidelidade das partes dependentes. Manual Windows DPI 100/150/200% e PowerPoint 2016/365 registra revisão/hash/ambiente/executor/observado; aprovação exige usuário sobre essa revisão. Este plano não preenche quickstart nem conclui tarefas técnicas/manuais.