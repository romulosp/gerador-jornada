# Implementation Plan: Domínio Projeto

## Situação corrente — revisão 1.0.0.1

Project contém ID, Name, Revision, CatalogVersion e Assets; Document agrega telas, jornada, roteiros, JavaOptions e Resources. Writer 3/3; pares de leitura reconhecidos: 1/1, 2/1, 3/2 e 3/3. Campos desconhecidos são recusados. Perfis visuais privados por documento não foram implementados.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Branch**: `002-projeto` (sem branch Git criada) | **Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Input**: spec.md deste domínio e [plano coordenador](../001-gerador-jornada/plan.md).

## Summary

Definir o documento de uma jornada, sua identidade, revisão e recursos, sem conhecer ZIP, janela ou exportadores. Implementar em tipos pequenos, operações explícitas e testes puros. Não adicionar um serviço que concentre outros domínios.

## Technical Context

**Language/Version**: mesma toolchain Go fixada pela feature coordenadora; nenhum go.mod adicional.
**Primary Dependencies**: tipos/biblioteca padrão e dependências de domínio descritas abaixo.
**Storage**: domínio não faz IO; serialização é responsabilidade de projectzip.
**Testing**: testes Go junto ao pacote, usando fixtures e invariantes.
**Target Platform**: núcleo independente de SO; integrado ao editor Windows.
**Project Type**: pacote de domínio interno do monólito.
**Performance Goals**: contribuir para NFR-007 da feature pai; nenhuma goroutine adicional nas regras puras.
**Constraints**: documentação integral conforme constituição; sem Fyne, XML, JSON DTO ou templates Java no domínio.
**Scale/Scope**: conjunto nominal da feature pai; responsabilidades limitadas a Projeto.

## Constitution Check

Conforme princípios I/II: modelo coeso e sem adaptadores. IV/V exigem documentação integral e testes de regras/erros; os checks atuais aprovados não substituem a revisão completa nem a matriz ainda aberta em tasks.md. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo. VII: rastreabilidade reconciliada com estado parcial explícito, sem declarar o domínio aprovado.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/project/model.go` e `model_test.go`; composição/snapshot em `internal/application/document.go` e `session_test.go`; manifesto/ZIP em `internal/adapter/projectzip/wire.go`, `package.go` e respectivos testes. Os caminhos são relativos a apps/gerador-jornada. Não criar operations.go/validation.go apenas para repetir uma divisão planejada; dividir por responsabilidade quando a evolução justificar.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

A aplicação resolve referências entre módulos e constrói snapshot. O adaptador projectzip converte DTOs e bytes; Project não importa screen/journey/scenario para executar suas regras nem seus adaptadores. A composição do documento é montada na aplicação, evitando um agregado com todos os comportamentos.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/project-format.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar IDs duplicados, nome vazio, revisão negativa, recurso compartilhado e isolamento de snapshot. Verificar caso válido, falha e preservação de estado; `go test ./internal/project` já é executável; cobertura atual e casos restantes estão identificados em tasks.md, sem equiparar teste parcial à matriz completa.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.


## Estado reconciliado — 2026-10-18

Identidade e metadados implementados; snapshot, manifesto v2/migração v1 e persistência integrados. Remoção segura de recursos e matriz completa de falhas/recuperação ainda precisam de cobertura; não concluir T006/T007 apenas pelo round-trip.

Checks desta revisão: `scripts/check.ps1 -Jdk17 C:\Desenvolvimento\jdk-17.0.11` aprovados; escopo e evidências em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18). Não implica teste manual nem aprovação do domínio.

## Plano de evolução para schema 3/catalog 2 — 2026-10-19

Implementado parcialmente na continuação: `application.Document.AppearanceProfile
*skin.Profile` (não campo de Project), cópia profunda, biblioteca imutável e captura
de metadados por `internal/adapter/visualprofile`. `projectzip` lê v1/v2/v3 e escreve
v3/catalog2, captura antes de serializar e valida perfil/bytes/limites. `appearanceProfile`
é opcional na leitura v3 em transição, proibido em v1/v2. `atomicfile` e fluxo de salvar
preservam original/cópia antes do primeiro commit migrado. Parsing continua adaptador.

Captura conserva todos os metadados da biblioteca original em cópias privadas:
definições, fontes, aliases, recursos e todas as revisões arquivadas (até 16 diretas),
inclusive não usadas e todos os estados. Não enriquece perfis privados com builtin
nem faz upgrade; uso corrente só limita diagnósticos e acréscimo de entrada missing.
Desvio arquitetural justificado pelo bug crítico save/reopen → AddComponent
`defaultconfirm`: preserva inserir famílias/trocar variantes após reopen, não é
incorporação binária integral. Alvo SUBSET de bytes usados e permitidos permanece;
payload selado mais amplo altera o hash. Não copia novos binários sem autorização,
mesmo com JSON `allowed`.
`ResourceData` exige AssetID/hash/tamanho/MIME de embedded ou builtin de bytes exatos
para referenced. `Library` resolve revisão exata, sem fallback para desconhecida.
Hash determinístico inclui arquivos/hashes e exclui somente o hash próprio.

Não duplicar o importador/fontes/biblioteca de 001/T108–T110: consumir seus dados/referências aprovados; validação do pacote do usuário é responsabilidade distinta da validação do acervo de desenvolvimento. O modelo pode ser preparado com fixtures antes de todos os recursos estarem disponíveis, mas fidelidade dependente permanece bloqueada. Preservar as restrições offline e limites existentes, sem paths absolutos em perfil/assets.

T009–T012 mantêm trabalho restante: normalização `component.appearance`, overrides
com presença completa/RGBA e descritores-alvo `fontFaces/resourceIds` versus os atuais
`Style`/`Fonts/Aliases/Resources`; portabilidade licenciada e execução da matriz.
Família desconhecida continua literal em Style com diagnóstico, não descriptor
unresolved implementado. Não ocultar o alvo para tratar a transição como concluída.
Após correção integral de Capture, conforme main: gofmt/vet PASS e **15 suítes Go
compiladas, não executadas**, bloqueio prévio por política Windows preservado.
Python **28 PASS** (15 importador + 4 packaging + 9 Java 6, incluindo erro esperado
de `String.strip`, não homologação EFL); importer `--check` PASS, mesmo inventário
214/8/270/47. Distribuição Windows amd64 **gerada e verificada**; checksum do ZIP no
sidecar externo `.zip.sha256`. Reempacotar documentos preserva o EXE. Identidade do EXE e manual pendente no
[quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
001/T118 integra, T120/T121 homologa e T122 obtém aceite; continuam abertos pelo escopo
restante. Histórico v2 acima não instrui o writer a permanecer em v2.

<a id="plano-codigos-barras-2026-10-23"></a>
## Plano de domínio — códigos de barras (2026-10-23)

FR-011/T013–T014, após estabilização de 003/T018/T020 e 004/T015. Atualizar tradução estrita de Properties em `internal/adapter/projectzip/wire.go` e fixtures no mesmo adaptador; Project continua sem JSON/ZIP/UI ou nova entidade financeira. O [contrato](../001-gerador-jornada/contracts/project-format.md#formato-codigos-barras-2026-10-23) define o alvo 3/3 e recusa de leitor antigo; documentação não modifica writer corrente 3/2.

Manter reader 1/1, 2/1, 3/2 e acrescentar 3/3; upgrade ao salvar informa incompatibilidade e preserva original antes de sobrescrever, mesmo quando schema continua 3 e somente catalogVersion muda. Reutilizar salvamento atômico/cópia existente, sem duplicar IO. Arquivos sem novos componentes não ganham instâncias, valores ou recursos por migração. Perfis próprios/arquivados permanecem exatos, sem enriquecer definições para supostamente suportar os compostos.

Testes locais cobrem ida/volta, versão incompatível, enum/UUID/estrutura inválidos, valores Unicode/vazios, ausência de runtime e falha de publicação preservando original. Integração desktop/PPT e evidência ficam em 001/T138–T139; nenhum teste de adaptador equivale a validação visual ou financeira.

<a id="plano-menu-aplicacao-2026-10-23"></a>
## Plano de domínio — codec de menu

FR-012/T015 estende `internal/adapter/projectzip/wire.go` e testes vizinhos após 003/T021/004/T017. Reutilizar matriz/backup de T013–T014, sem segundo mecanismo ou catálogo 4. DTO recursivo estrito segue [project-format.md](../001-gerador-jornada/contracts/project-format.md#formato-menu-aplicacao-2026-10-23); Enabled ausente não se confunde com false explícito: defaults só na factory, não no reader.

Limites de manifesto/ZIP continuam; validar profundidade/quantidade na fronteira antes de tradução recursiva ilimitada. Traduzir Properties.Menu sem copiar regra de negócio para Project. Campos extras/enums/IDs inválidos falham localmente; caminho runtime nunca é aceito como configuração. Abrir preserva bytes, salvar faz publicação atômica com cópia original quando incompatível. Perfis/revisões/assets não são enriquecidos pelo menu. Integração em 001/T143; manual em T144; nada executado nesta entrega documental.
