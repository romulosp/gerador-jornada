# Implementation Plan: Domínio Tela

## Situação corrente — revisão 1.0.0.1

Screen mantém ID/Name/FileID, Size, Background e Components em ordem de camadas. Component mantém ID/Name/FileID, Kind, Bounds, Style, Enabled, ParentID e Properties. Style tem FontFamily/FontSize/TextColor/FillColor/BorderColor/BackgroundAssetID/TextAlign/Skin; não possui FillAlpha nem appearance tipada.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


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

Conforme princípios I/II: modelo coeso e sem adaptadores. IV/V exigem documentação integral e testes de regras/erros; os checks atuais aprovados não substituem a revisão completa nem a matriz ainda aberta em tasks.md. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo. VII: rastreabilidade reconciliada com estado parcial explícito, sem declarar o domínio aprovado.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/screen/model.go`, `selection.go`, `file_id.go` e testes; comandos em `internal/application/edit_screen.go`, `edit_selection.go` e `transform.go`; alinhamento/camadas em `internal/adapter/desktop/arrangement.go`. Os caminhos são relativos a apps/gerador-jornada. Não criar operations.go/validation.go apenas para repetir uma divisão planejada; dividir por responsabilidade quando a evolução justificar.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

Importa catalog para tipos e validação de família. Usa valores/IDs, não widgets. Aplicação agrupa alterações de tela e referências de Journey/Scenario num comando reversível único. Presentation consome geometria sem modificar entidades.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Verificar caso válido, falha e preservação de estado; `go test ./internal/screen` já é executável; cobertura atual e casos restantes estão identificados em tasks.md, sem equiparar teste parcial à matriz completa.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.


## Estado reconciliado — 2026-10-18

Modelo, bounds, seleção coletiva, IDs de arquivo, duplicação e remoção reversível implementados. Parte dos comandos de composição ainda é coordenada no desktop por Session.Change. A matriz completa de bordas/duplicação e todos os alinhamentos previstos ainda requerem conferência; DPI nativo permanece manual.

Checks desta revisão: `scripts/check.ps1 -Jdk17 C:\Desenvolvimento\jdk-17.0.11` aprovados; escopo e evidências em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18). Não implica teste manual nem aprovação do domínio.

## Evolução de composição/aparência — 2026-10-19

`screen/model.go` já tem ParentID opcional e Style com SkinName/SkinVersion e campos
suportados; não tipos Fyne/parser Java. Presença tipada completa, RGBA e
`component.appearance` normalizado continuam alvo, não implementação entregue.
Defaults/criação consomem catálogo 2; validação garante tipo compatível, pai contêiner
da mesma tela e ausência de ciclos. ParentID não converte coordenadas persistidas:
traduzir absolutos para locais apenas na pintura (001/T111), mantendo z-order entre
irmãos e subárvore. Não reordenar/redimensionar automaticamente ao migrar.

Perfil de documento já é `Document.AppearanceProfile *skin.Profile`, com Library
imutável e revisão arquivada exata. Capture guarda todos os metadados da biblioteca
original em cópias privadas: definições, fontes, aliases, recursos e todas as revisões
arquivadas, inclusive não usadas e todos os estados. Não enriquece perfis privados
com builtin nem faz upgrade; uso só limita diagnósticos e acréscimo de entrada missing.
Desvio arquitetural pelo bug crítico save/reopen → AddComponent `defaultconfirm`,
preserva inserir famílias/trocar variantes após reopen. Não é incorporação binária
integral: alvo SUBSET de bytes usados e permitidos permanece, payload selado mais
amplo altera o hash. Não acrescenta binários sem autorização nem por JSON allowed.
Desktop/PPT usam perfil e cache SHA/face exata; galerias estáticas16 Java/XML não são
screenshots Java/homologação DPI. Após correção integral de Capture, conforme main:
gofmt/vet PASS, **15 suítes Go compiladas, não executadas**, bloqueio prévio por
política Windows preservado; **28 Python PASS** (15 importador + 4 packaging + 9 Java 6,
erro esperado de `String.strip` incluído, não homologação EFL); importer `--check`
PASS, mesmo inventário 214/8/270/47. Distribuição Windows amd64 **gerada e verificada**;
checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos preserva
o EXE. Identidade do EXE e manual pendente no
[quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
T009–T012 continuam abertos pelo escopo restante.

Casos de uso em application/edit_screen.go/edit_selection.go/transform.go e cópias do documento incluem novas propriedades, hierarquia, referências de grupo/teclado/opções/linhas. Calcular fechamento dos descendentes antes de mover/duplicar/remover; transação única sem deslocamento duplo. Resize do pai não escala filhos; validar clip/overflow e diagnosticar o excedente conforme contrato da família. Presets em application/presets.go reutilizam controles, sem novo domínio e sem integrações reais.

T009–T012 locais tratam modelo/comandos/presets/testes. Projeto 002/T010 serializa; Catálogo 003/T010–T015 define tipos; coordenadora 001/T111–T116 pinta/edita/executa/exporta. 014–018 têm somente suas regras visuais afetadas. Sem novo módulo Go; fidelidade e aceite dependem dos manuais/recursos ainda abertos.

## Sobreposição de componentes — 2026-10-20

FR-018/T014: implementado [screen.PlaceComponent](../../apps/gerador-jornada/internal/screen/layer.go), operação coesa de reordenação estável da lista plana `Components` (fundo → frente). Referência vazia aplica extremo; referência válida aplica adjacência antes/depois. Validar participantes antes da escrita, retornar erro sem mutação e tratar posição já satisfeita como no-op. Não introduzir campos, schema/catálogo novo, constraints relacionais ou serviço central.

**Delimitação do contrato de hierarquia acima:** a colocação de camada implementada não aplica ordenação de subárvore/grupo. Move somente a entrada selecionada na lista canônica; filhos não acompanham o pai, e ParentID/bounds/estilo/conteúdo permanecem intactos. Qualquer composição derivada na cena não muda essa semântica da operação nem cria relação persistida frente/atrás.

O adaptador inspector_layer.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_layer.go`) integra `componentInspector`: modos neutro/relativos/extremos e referência da mesma tela são rascunho; a chamada ocorre no mesmo `Session.Change` das outras propriedades, com rollback e histórico existentes. Botões Frente ↑/Fundo ↓ em arrangement continuam um passo, sem substituição. Persistência usa apenas a ordem já existente de Components.

Testes de domínio em layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/screen/layer_test.go`); integração desktop/projectzip/presentation em [001/T127](../001-gerador-jornada/tasks.md#phase-21-sobreposição-de-componentes--2026-10-20), sem duplicar implementação. Produção e regressões implementadas. O executor principal informou gofmt dos sete arquivos e vet PASS, 15 suítes Go **somente compiladas**, quatro testes Python de empacotamento PASS e build/pacote Windows com integridade PASS; evidências e identidade na [seção 9](../001-gerador-jornada/quickstart.md#9-sobreposição-de-componentes--2026-10-20). T014/T127 permanecem abertas: testes Go não executados por política Windows, sem contorno; UI nativa, DPI, Office e aprovação manual não realizados. Checks/hashes anteriores pertencem somente às revisões históricas.

## Incremento de alinhamento textual universal — 2026-10-20

Reutilizar `Style.TextAlign` e sua validação/cópia/histórico existentes para FR-017, sem mudar schema/catálogo nem criar domínio. T013 local verifica invariantes do modelo: a edição textual não muda bounds, ParentID, hierarquia, identidade, conteúdo, ordem ou vínculos; manter left histórico e herança ao remover override. Não implementar um segundo mecanismo de alinhamento ou duplicar a matriz dos adaptadores.

Integração em [001/T126](../001-gerador-jornada/tasks.md#phase-20-alinhamento-textual-universal--2026-10-20): campo sempre disponível nas 16 famílias, resolução textual compartilhada e popup de combo; cena mantém regiões textuais, fundos e hit areas. Quatro novos arquivos de regressões desktop/presentation/projectzip/PPT contemplam 48 combinações (16 × três overrides), sem evidência de execução. Em 20/10, após produção e testes novos, o executor principal informou gofmt/vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e tarefa VS Code de compilação Windows PASS; registro na [seção 8 do quickstart](../001-gerador-jornada/quickstart.md#8-alinhamento-textual-universal--2026-10-20). T013 e T126 permanecem abertas com verificação parcial; execução Go não realizada pelo bloqueio prévio Windows, sem contorno. Manuais Windows/DPI/PowerPoint e aceite continuam abertos, sem mudança nos gates Java/fontes/PPT.

<a id="plano-codigos-barras-2026-10-23"></a>
## Plano de domínio — códigos de barras (2026-10-23)

FR-019/T015–T016: evoluir `internal/screen/model.go` (NewComponent/cópia/duplicação) e comandos existentes de composição da aplicação, sem introduzir novo agregado de negócio ou widgets no domínio. Catálogo 003/T018/T020 fornece propriedades/defaults; Tela atribui IDs de blocos/linhas na criação, verifica invariantes de composição e remapeia cópias/vínculos.

Manter entrada e processado como uma instância cada, inclusive mover/redimensionar/ordenar/remover. O histórico da aplicação conserva comandos confirmados de edição de listas/seletores/valores numa transação, sem enfileirar operação por campo digitado do runtime. Teclado pode apontar ao composto compatível; subcampo ativo é responsabilidade da sessão 001/T135, não oito alvos de Tela.

Testes locais de identidade/cópia/ordem e comandos ficam em screen/application. Serialização é 002/T013–T014; inspector/cena/prévia/PPT são 001/T135–T138. Não repetir algoritmo de DV no comando de composição. A proposta visual de 032 aguarda imagem e não comprova default legado.

<a id="plano-menu-aplicacao-2026-10-23"></a>
## Plano de domínio — composição do menu

FR-020/T017 consome 003/T021 em `internal/screen/model.go`, `internal/application/component_copy.go` e comandos existentes de edição/histórico. Factory recebe callback de identidade da composição; copiar snapshot conserva UUIDs, duplicação regenera nós e remapeia inicial por mapa old→new. Não converter itens em Screen.Components ou ParentID internos. Remapeio de destinos ao duplicar tela consulta a política explícita de Jornada em 005/T012; duplicar componente não cria tela de destino.

Editar/remover/reordenar ramo e escolher menu inicial passam por uma transação existente, validação de catálogo e índice global de UUIDs; erro preserva documento/rascunho confirmado. Bounds/z-order/FileID/ParentID e referências não alteradas permanecem. Nenhum callback Fyne, runtime de navegação ou DTO ZIP entra em Tela. Registro central é serial com T015–T016; formulários/cena usam os comandos após contrato estável. Testes e aprovação pendentes.
