# Implementation Plan: Domínio Tela

**Branch**: `004-tela` (sem branch Git criada) | **Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Input**: spec.md deste domínio e [plano coordenador](../001-gerador-jornada/plan.md).

## Summary

Modelar telas, componentes, geometria e ordem visual com operações coesas e independentes da UI. Implementar em tipos pequenos, operações explícitas e testes puros. Não adicionar um serviço que concentre outros domínios.

## Technical Context

**Language/Version**: mesma toolchain Go fixada pela feature coordenadora; nenhum go.mod adicional.
**Primary Dependencies**: tipos/biblioteca padrão e dependências de domínio descritas abaixo.
**Storage**: domínio não faz IO; serialização é responsabilidade de projectzip.
**Testing**: testes Go junto ao pacote, usando fixtures e invariantes.
**Target Platform**: núcleo independente de SO; integrado ao editor Windows.
**Project Type**: pacote de domínio interno do monólito.
**Performance Goals**: contribuir para NFR-007 da feature pai; nenhuma goroutine adicional nas regras puras.
**Constraints**: documentação integral conforme constituição; sem Fyne, XML, JSON DTO ou templates Java no domínio.
**Scale/Scope**: conjunto nominal da feature pai; responsabilidades limitadas a Tela.

## Constitution Check

Conforme princípios I/II: modelo coeso e sem adaptadores. Conforme IV/V: documentar tipos, funções, métodos, campos e variáveis e testar regras, incluindo erros. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/screen/model.go`, `operations.go`, `validation.go` e testes correspondentes; separar outros arquivos por responsabilidade se necessário. Os caminhos de código desta feature são relativos a apps/gerador-jornada; a implementação segue o backlog local.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

Importa catalog para tipos e validação de família. Usa valores/IDs, não widgets. Aplicação agrupa alterações de tela e referências de Journey/Scenario num comando reversível único. Presentation consome geometria sem modificar entidades.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Verificar caso válido, falha e preservação de estado; executar `go test ./internal/screen` somente após a implementação existir.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.
