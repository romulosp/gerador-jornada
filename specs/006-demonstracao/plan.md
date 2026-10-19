# Implementation Plan: Domínio Demonstração

**Branch**: `006-demonstracao` (sem branch Git criada) | **Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Input**: spec.md deste domínio e [plano coordenador](../001-gerador-jornada/plan.md).

## Summary

Interpretar roteiro explícito como sequência finita de estados e eventos temporizados, reutilizável por prévia e PPT. Implementar em tipos pequenos, operações explícitas e testes puros. Não adicionar um serviço que concentre outros domínios.

## Technical Context

**Language/Version**: mesma toolchain Go fixada pela feature coordenadora; nenhum go.mod adicional.
**Primary Dependencies**: tipos/biblioteca padrão e dependências de domínio descritas abaixo.
**Storage**: domínio não faz IO; serialização é responsabilidade de projectzip.
**Testing**: testes Go junto ao pacote, usando fixtures e invariantes.
**Target Platform**: núcleo independente de SO; integrado ao editor Windows.
**Project Type**: pacote de domínio interno do monólito.
**Performance Goals**: contribuir para NFR-007 da feature pai; nenhuma goroutine adicional nas regras puras.
**Constraints**: documentação integral conforme constituição; sem Fyne, XML, JSON DTO ou templates Java no domínio.
**Scale/Scope**: conjunto nominal da feature pai; responsabilidades limitadas a Demonstração.

## Constitution Check

Conforme princípios I/II: modelo coeso e sem adaptadores. IV/V exigem documentação integral e testes de regras/erros; os checks atuais aprovados não substituem a revisão completa nem a matriz ainda aberta em tasks.md. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo. VII: rastreabilidade reconciliada com estado parcial explícito, sem declarar o domínio aprovado.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/scenario/model.go` e `model_test.go`; quadros puros em `internal/application/playback.go`; edição/reprodução em `internal/adapter/desktop/scenarios.go` e `player.go`; animações em `internal/adapter/pptx/timing.go`. Os caminhos são relativos a apps/gerador-jornada. Não criar operations.go/validation.go apenas para repetir uma divisão planejada; dividir por responsabilidade quando a evolução justificar.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

Recebe snapshot resolvido de telas/componentes e resolução Journey, preparado pela aplicação. Interpretador devolve eventos puros. FrameAt deriva quadros na aplicação; relógio/pausa ainda estão em desktop/player.go. A arquitetura alvo exige separar o controle temporal e testá-lo com relógio controlado (001-gerador-jornada/T045). PPT consome eventos puros, sem temporizador de UI.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar intervalo temporal, campos por kind, alvo removido, tela errada, A-B-A, auto-navegação, typeText vazio, pausa e reinício. Verificar caso válido, falha e preservação de estado; `go test ./internal/scenario` já é executável; cobertura atual e casos restantes estão identificados em tasks.md, sem equiparar teste parcial à matriz completa.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.


## Estado reconciliado — 2026-10-18

Interpretador, ocorrências, grafemas, cursor/clique/digitação e transições XML implementados. Pausa/reinício existem na UI, mas o relógio ainda não é abstraído na aplicação; reordenação disponível para etapas, não para a lista de roteiros. T004/T007 e 001-gerador-jornada/T044–T047 conservam essas pendências; execução PowerPoint ainda não homologada.

Checks desta revisão: `scripts/check.ps1 -Jdk17 C:\Desenvolvimento\jdk-17.0.11` aprovados; escopo e evidências em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18). Não implica teste manual nem aprovação do domínio.
