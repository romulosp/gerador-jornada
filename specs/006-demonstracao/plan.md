# Implementation Plan: Domínio Demonstração

## Situação corrente — revisão 1.0.0.1

Scenario/Step são interpretados em eventos/ocorrências; application/playback.go e desktop/player.go consomem essa interpretação. Inputs runtime de prévia não são persistidos nos roteiros. typeText de compostos/menu e os tipos novos continuam sem capacidade disponível.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


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

## Integração de aparência/estados do catálogo 2 — 2026-10-19

**Continuação implementada:** `Document.AppearanceProfile *skin.Profile`, cópia
profunda, Library imutável e revisão arquivada exata já são consumidos por desktop/PPT,
com cache SHA/face exata. Capture conserva todos os metadados da biblioteca original
em cópias privadas: definições/fontes/aliases/recursos e todas as revisões arquivadas,
inclusive não usadas e todos os estados. Não enriquece perfis privados com builtin
nem faz upgrade; uso só limita diagnósticos e acréscimo de entrada missing. Desvio
arquitetural pelo bug crítico save/reopen → AddComponent `defaultconfirm`, preserva
inserir famílias/trocar variantes após reopen. Não é incorporação binária integral:
alvo SUBSET de bytes usados e permitidos permanece; payload selado mais amplo altera
o hash. Sem copiar novos binários sem autorização, mesmo com JSON allowed. Isso não insere
skin/fontes em Scenario nem altera o contrato temporal. Style/DTOs transitórios
ainda não concluem appearance/overrides/RGBA/descritores-alvo. Galerias estáticas16
Java/XML não são screenshots Java/homologação Windows/DPI. PPT mantém
FontEmbedding=false/CompleteFidelity=false e erro PPT_FONT_EMBEDDING_UNSUPPORTED para
EmbedFonts/StrictFidelity. Após correção integral de Capture, conforme main:
gofmt/vet PASS; **15 suítes Go compiladas, não executadas**, bloqueio prévio por
política Windows preservado; **28 Python PASS** (15 importador + 4 packaging + 9 Java 6,
incluindo erro esperado de `String.strip`, não homologação EFL); importer `--check`
PASS, mesmo inventário 214/8/270/47. Distribuição Windows amd64 **gerada e verificada**;
checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos preserva
o EXE. Identidade do EXE e manual pendente no
[quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
T009/T010 e PowerPoint/aceite continuam abertos.

Interpretador recebe capacidades/valores iniciais tipados no snapshot, não classes/widgets; FrameAt produz estado/cena coerente com a resolução visual comum de 001/T111. Expandir validação de typeText por capacidade para entradas novas, mantendo click só visual e navigate explícito. Não acrescentar Kind de Step por conveniência de UI; eventos locais sem etapa contratada são prévia livre, não roteiro anunciado pronto.

Ocorrências copiam estados iniciais de novas famílias e recorte/seleções tipados necessários ao quadro, sem salvar uso no documento. Relógio controlado e separação de player continuam pendências anteriores, não resolvidas por documentação. T009–T010 locais testam interpretação/quadros; 001/T114–T116 trata sessão/widgets/PPT. Cada família/variante tem fixture; nenhum automatizado substitui PowerPoint ou aprovação da revisão.

<a id="plano-codigos-barras-2026-10-23"></a>
## Plano de domínio — códigos de barras (2026-10-23)

FR-011/T011: consumir os tipos estabilizados de 003/T018/T020 no índice de capacidades e nas fixtures de `internal/scenario`/`internal/application/playback.go`, mantendo os cinco kinds existentes. Não permitir `typeText` em `barcodeInput`/`barcodeProcessed` nesta entrega; código de captura interativa é a prévia livre 001/T135, não o interpretador ou PPT. Relógio, navegação e regras de etapas não mudam por existir um componente semelhante a campo.

Estado inicial/quadros preservam a ordem e Unicode configurados na ficha processada e entrada vazia. Cena e exportador compartilham apresentação em 001/T137, com textos/formas editáveis e diagnósticos de capacidade antes do job. Não rasterizar a ficha nem inferir transição entrada→processado. Validação temporal/quadros é local; integração/PPT real e aceite ficam 001/T138–T139.

<a id="plano-menu-aplicacao-2026-10-23"></a>
## Plano de domínio — quadro de menu

FR-012/T012 consome 003/T021/004/T017 em `internal/scenario/model.go` e fixtures de `internal/application/playback.go`, sem duplicar cena ou runtime. Inicial do quadro vem de Properties.Menu.InitialMenuID, derivando caminho para o submenu de modo determinístico; não captura stack da prévia. A configuração não muda tempos/kinds de Step. Validação recusa typeText e eventos de item não contratados antes do job; código de movimento/atalho não se torna evento.

Coordenador decide separadamente se estenderá click com item metadata; nesta entrega o contrato seguro é click visual sem mutação e sem novo campo/Kind de Step. Se decidir ampliar, sincronizar Scenario/Step, wire, capacidades, resolução de hit regions, quadros e testes antes de expor; sem isso, somente representação estática do menu inicial. 001/T142 cuida da cena/PPT e T143 da integração; nenhuma promessa de árvore operante no PowerPoint.
