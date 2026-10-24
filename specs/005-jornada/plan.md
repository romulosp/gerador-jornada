# Implementation Plan: Domínio Jornada

## Situação corrente — revisão 1.0.0.1

Journey mantém StartScreenID e Actions none/navigate para botões. Folhas de menu guardam TargetScreenID em suas propriedades e o runtime trata esse destino. Não há NavigationGraph/Graph ou índice completo de integridade/alcançabilidade de menus na árvore atual.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Branch**: `005-jornada` (sem branch Git criada) | **Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Input**: spec.md deste domínio e [plano coordenador](../001-gerador-jornada/plan.md).

## Summary

Modelar tela inicial e navegação de botões como dados explícitos, compartilhados pela prévia e pelas saídas. Implementar em tipos pequenos, operações explícitas e testes puros. Não adicionar um serviço que concentre outros domínios.

## Technical Context

**Language/Version**: mesma toolchain Go fixada pela feature coordenadora; nenhum go.mod adicional.
**Primary Dependencies**: tipos/biblioteca padrão e dependências de domínio descritas abaixo.
**Storage**: domínio não faz IO; serialização é responsabilidade de projectzip.
**Testing**: testes Go junto ao pacote, usando fixtures e invariantes.
**Target Platform**: núcleo independente de SO; integrado ao editor Windows.
**Project Type**: pacote de domínio interno do monólito.
**Performance Goals**: contribuir para NFR-007 da feature pai; nenhuma goroutine adicional nas regras puras.
**Constraints**: documentação integral conforme constituição; sem Fyne, XML, JSON DTO ou templates Java no domínio.
**Scale/Scope**: conjunto nominal da feature pai; responsabilidades limitadas a Jornada.

## Constitution Check

Conforme princípios I/II: modelo coeso e sem adaptadores. IV/V exigem documentação integral e testes de regras/erros; os checks atuais aprovados não substituem a revisão completa nem a matriz ainda aberta em tasks.md. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo. VII: rastreabilidade reconciliada com estado parcial explícito, sem declarar o domínio aprovado.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/journey/model.go` e `model_test.go`; validação entre telas/componentes em `internal/application/document.go`; execução isolada em `internal/application/preview.go` e `session_test.go`. Os caminhos são relativos a apps/gerador-jornada. Não criar operations.go/validation.go apenas para repetir uma divisão planejada; dividir por responsabilidade quando a evolução justificar.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

Recebe índice imutável de referências permitido pela aplicação, sem importar Screen ou adaptadores. Scenario usa resolução de rota para interpretar navigate. Desktop e Java executam a ação; PPT traduz em links.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar none implícito, ciclo, auto-navegação, fonte não botão, destino apagado, inicial ausente e binding duplicado. Verificar caso válido, falha e preservação de estado; `go test ./internal/journey` já é executável; cobertura atual e casos restantes estão identificados em tasks.md, sem equiparar teste parcial à matriz completa.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.


## Estado reconciliado — 2026-10-18

Binding none/navigate, resolução por ID e validação de origem/destino estão implementados. O warning de tela inalcançável previsto em FR-004 ainda não é emitido; permanece em T006. Diagnóstico clicável/localização completa também depende de 001-gerador-jornada/T030.

Checks desta revisão: `scripts/check.ps1 -Jdk17 C:\Desenvolvimento\jdk-17.0.11` aprovados; escopo e evidências em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18). Não implica teste manual nem aprovação do domínio.

## Integração catálogo 2 — 2026-10-19

**Continuação implementada:** a composição da aplicação já transporta
`Document.AppearanceProfile *skin.Profile`, copiado profundamente; desktop/PPT
consomem Library imutável por documento e revisão arquivada exata, sem fallback
para desconhecida. Capture mantém todos os metadados da biblioteca original em cópias
privadas: definições/fontes/aliases/recursos e todas as revisões arquivadas, inclusive
não usadas e todos os estados. Não enriquece perfis privados com builtin nem faz
upgrade; uso só limita diagnósticos e acréscimo de entrada missing. Desvio arquitetural
pelo bug crítico save/reopen → AddComponent `defaultconfirm`, preserva inserir famílias/
trocar variantes após reopen. Não é incorporação binária integral: alvo SUBSET de bytes
usados e permitidos permanece; payload selado mais amplo altera o hash. Sem copiar
novos binários sem autorização, mesmo com JSON allowed. Jornada não passa a persistir
fontes/skins ou estados de UI; IDs/rotas permanecem independentes da aparência.
Style/DTOs atuais não concluem normalização appearance/overrides/RGBA/descritores.
Galerias estáticas16 Java/XML não homologam execução Java ou Windows/DPI.
Após correção integral de Capture, conforme main: gofmt/vet PASS; **15 suítes Go
compiladas, não executadas**, bloqueio prévio por política Windows preservado;
**28 Python PASS** (15 importador + 4 packaging + 9 Java 6, incluindo erro esperado
de `String.strip`, não homologação EFL); importer `--check` PASS, mesmo inventário
214/8/270/47. Distribuição Windows amd64 **gerada e verificada**;
checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos preserva o EXE.
Hash/tamanho e manual pendente no [quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
T009/T010 continuam abertos para escopo integral e execução.

Jornada continua com none/navigate por referência, não um motor de eventos de negócio. Evoluir o índice imutável recebido da aplicação para capacidades explícitas de origem/alvo; validar twoLabel como button e recusar novos Kinds sem capacidade de navegação. Não importar Fyne/skins/PPT nem guardar valores runtime em Journey.

T009–T010 locais cobrem resolução/diagnósticos/regressão com novos tipos, sem repetir os metadados de 003/T010 ou estado runtime de 001/T114. application/document.go coordena referências, preview.go executa reset isolado e desktop/pptx traduzem ações já validadas. Contrato compartilhado [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md#expansão-pós-v1--capacidades-e-estado-local-2026-10-19); aparência/migração nunca remapeia rotas por nome/posição. EFL permanece gate externo pausado.

<a id="plano-codigos-barras-2026-10-23"></a>
## Verificação de fronteira — códigos de barras (2026-10-23)

FR-008/T011 amplia fixtures de `internal/journey` e do índice em `internal/application/document.go`, sem alterar o modelo de ações. Consumir capacidades do Catálogo após 003/T018; recusar compostos como origem de navigate e preservar ações dos botões existentes. 001/T135 controla a conclusão local e reset da entrada, não o domínio Journey. Nenhuma UI/validação matemática é acrescentada a Journey. Regressão A→B→A pode usar aplicação após o runtime, sem duplicar testes de widget. Esta revisão é documental; contrato de campos e regras de roteamento permanece [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md#limites-codigos-barras-2026-10-23).

<a id="plano-menu-aplicacao-2026-10-23"></a>
## Plano de domínio — destinos autorados do menu

FR-009/T012 estende validação pura de destinos em `internal/journey/model.go` e índice de referências em `internal/application/document.go` após 003/T021/004/T017. Aplicação fornece projeção de folhas com sourceScreenID/componentID/itemID/targetScreenID e enabled efetivo; Journey não importa catálogo/desktop/ZIP e não armazena cópia da árvore. Reusar regras de tela existente, impacto e grafo; não criar Action por item ou duplicar target em dois lugares.

TargetScreenID em Properties.Menu é a configuração única; autoridade de validação/remoção/destino pertence a Jornada. Auto/ciclos são explícitos; árvore local nunca entra como tela no grafo. Duplicar componente mantém destinos; duplicar tela remapeia destinos à própria tela antiga para a nova e preserva externos, coordenado com 004/T017. Runtime é 001/T140 e não executa metadados. Contrato em [navigation-scenario.md](../001-gerador-jornada/contracts/navigation-scenario.md#navegacao-menu-aplicacao-2026-10-23). Sem novo Kind de Action/Step; decisão de clique demonstrativo fica com coordenação.
