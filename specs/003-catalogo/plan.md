# Implementation Plan: Domínio Catálogo

## Situação corrente — revisão 1.0.0.1

catalog.Definition contém Kind/Name/Width/Height. Properties é uma união tipada de dez Kinds; catalogVersion=4 não significa implementação das 16 famílias planejadas. As skins estáticas vêm de internal/skin e não são entidades persistidas de Catálogo.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Branch**: `003-catalogo` (sem branch Git criada) | **Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Input**: spec.md deste domínio e [plano coordenador](../001-gerador-jornada/plan.md).

## Summary

**Referência vigente:** consultar a situação corrente 1.0.0.1 acima. A regra de filtragem de dependências ausentes permanece; os inventários de perfis 214/206 e suas APIs são históricos/alvo de expansão.

Definir as sete famílias, seus defaults, propriedades tipadas e capacidades sem executar código do projeto. Implementar em tipos pequenos, operações explícitas e testes puros. Não adicionar um serviço que concentre outros domínios.

## Technical Context

**Language/Version**: mesma toolchain Go fixada pela feature coordenadora; nenhum go.mod adicional.
**Primary Dependencies**: tipos/biblioteca padrão e dependências de domínio descritas abaixo.
**Storage**: domínio não faz IO; serialização é responsabilidade de projectzip.
**Testing**: testes Go junto ao pacote, usando fixtures e invariantes.
**Target Platform**: núcleo independente de SO; integrado ao editor Windows.
**Project Type**: pacote de domínio interno do monólito.
**Performance Goals**: contribuir para NFR-007 da feature pai; nenhuma goroutine adicional nas regras puras.
**Constraints**: documentação integral conforme constituição; sem Fyne, XML, JSON DTO ou templates Java no domínio.
**Scale/Scope**: conjunto nominal da feature pai; responsabilidades limitadas a Catálogo.

## Constitution Check

Conforme princípios I/II: modelo coeso e sem adaptadores. IV/V exigem documentação integral e testes de regras/erros; os checks atuais aprovados não substituem a revisão completa nem a matriz ainda aberta em tasks.md. Conforme VI: sem interfaces para cada struct, DI framework, banco ou serviço externo. VII: rastreabilidade reconciliada com estado parcial explícito, sem declarar o domínio aprovado.

## Project Structure

### Documentation (this feature)

spec.md, plan.md, data-model.md e tasks.md. Pesquisa e contratos compartilhados permanecem na feature coordenadora, sem cópia divergente de formatos.

### Source Code (apps/gerador-jornada)

`internal/catalog/model.go`, `copy.go` e `model_test.go`; criação de defaults em `internal/screen/model.go`; round-trip das sete famílias em `internal/adapter/projectzip/package_test.go`. Os caminhos são relativos a apps/gerador-jornada. Não criar operations.go/validation.go apenas para repetir uma divisão planejada; dividir por responsabilidade quando a evolução justificar.

**Structure Decision**: um pacote por domínio, sem classes centrais ou novo módulo Go.

## Interfaces and Data Flow

Catálogo é definição embutida e sem IO. Screen consome suas definições; apresentação consulta metadados para expandir formas. Perfil EFL no adaptador resolve suporte efetivo, sem importar o adaptador no domínio.

Operações recebem valores tipados, devolvem estado/resultados concretos e erros de regra. A aplicação traduz os erros para Diagnostic localizado. Não executar callbacks de UI ou gravar arquivos dentro das operações.

Normas de formato/propriedades: [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md). Fonte de verdade dos campos/invariantes: [data-model.md](data-model.md), alinhado ao contrato compartilhado. Mudança incompatível deve atualizar os dois na mesma entrega.

## Validation and Integration

Testar sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível. Verificar caso válido, falha e preservação de estado; `go test ./internal/catalog` já é executável; cobertura atual e casos restantes estão identificados em tasks.md, sem equiparar teste parcial à matriz completa.

Integrar conforme [tasks.md](tasks.md) deste domínio e checkpoints da [feature coordenadora](../001-gerador-jornada/tasks.md). Não marcar checkpoint de integração antes dos testes de domínio.

## Complexity Tracking

Sem exceção à constituição. Tipos concretos bastam para regras internas; fronteiras de IO permanecem fora deste pacote.


## Estado reconciliado — 2026-10-18

Sete famílias, união tipada, Lookup e validação de enums implementados. Definition contém Kind/Name/Width/Height; metadados completos de propriedades/capacidades ainda não são um modelo declarativo. Inspector usa formulários explícitos por família. Perfil EFL segue no adaptador e não está homologado.

Checks desta revisão: `scripts/check.ps1 -Jdk17 C:\Desenvolvimento\jdk-17.0.11` aprovados; escopo e evidências em [quickstart da coordenadora](../001-gerador-jornada/quickstart.md#reconciliação-documental-por-domínio--2026-10-18). Não implica teste manual nem aprovação do domínio.

## Plano catálogo 2 — expansão aprovada em 2026-10-19

A base histórica acima tem sete famílias; catálogo 2 já tem 16, sem encerrar paridade
de todo o acervo. `catalog/model.go`/`copy.go` evoluem com definições estáticas e
Properties/variantes tipadas, sem map livre ou factories Fyne/PPT no núcleo.
`internal/skin` contém Profile/Library imutável por documento; bytes builtin continuam
em legacyskin e `internal/adapter/visualprofile` captura/resolve as dependências
exatas, incluindo arquivos diretos. Referência desconhecida não recebe fallback;
JSON allowed não autoriza cópia automática de binários. Propriedades comuns são de
Tela; metadados/arquivos ficam em Document, não em catalog.Definition.

Galerias estáticas das 16 famílias com origem Java/XML estão implementadas; não são
screenshots Java nem homologação DPI. `Style` e DTOs Fonts/Aliases/Resources são
transitórios: normalização appearance/overrides/RGBA/descritores continua alvo aberto.
Após correção integral de Capture, conforme main: gofmt/vet PASS; **15 suítes Go
compiladas, não executadas**, bloqueio prévio por política Windows preservado.
Python **28 PASS** (15 importador + 4 packaging + 9 Java 6, erro esperado de `String.strip`
incluído, não homologação EFL); importer `--check` PASS, mesmo inventário 214/8/270/47.
Distribuição Windows amd64 **gerada e verificada**; checksum do ZIP no sidecar externo
`.zip.sha256`. Reempacotar documentos preserva o EXE. Identidade do EXE e manual
pendente no [quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
Paridade, execução de testes Go e aceite continuam pendentes.

T009 fecha classificação granular e defaults com referência; T010 prepara união/metadados/capacidades. Lotes T011–T015 implementam somente regras/modelos/defaults/deep-copy do catálogo. Criação e comandos de composição usam 004/T009–T012; wire/ZIP usa 002/T009–T012; cena/desktop/prévia/PPT usam 001/T111–T116, sem duplicação. Cada lote tem fixtures de modelo e suporte dos consumidores antes de liberação; T016 executa a matriz de regressão do catálogo.

Começar por Button/Label/Title e partes Panel necessárias ao ReportEvent piloto (012), depois C1 seleção/duas legendas, C2 campos/combos/tabela/teclados (013 numérico), C3 demais Reports/Progress e C4 contêineres/dialog/presets. Tipo incompleto é visível como bloqueado, não anunciado pronto. Recursos/fontes em outras_fontes estão presentes por nome; licença/métricas, assets e delegates ausentes continuam gates independentes de fidelidade. Java permanece pausado; catálogo 2 não homologa perfis EFL.

<a id="plano-builtin-filtrado"></a>
## Disponibilidade operacional após redução do builtin — 2026-10-21

[FR-012](spec.md#builtin-filtrado-fr-012) clarifica FR-007/FR-010: conservar a classificação/inventário completo, mas não oferecer as oito definições adiadas para nova inserção/troca no builtin filtrado. A implementação do filtro pertence ao importador de desenvolvimento em [001/plano](../001-gerador-jornada/plan.md#plano-builtin-filtrado), não ao domínio Catálogo. Não duplicar parser, persistência, modelos ou geração em T009–T016.

- Consumir somente definições ativas da revisão builtin em novos projetos, por família e identidade `(Type, Name)`. Falhas de `imagePath`/`slider`/`font:type` retiram a definição inteira, inclusive estados futuros; falha de filho retira pai/todos os filhos pela composição real, sem órfãos nem exclusão de homônimos/prefixos independentes.
- Inventário existente: **206 definições ativas / 8 fontes / 270 recursos / 31 diagnósticos**, **8 definições adiadas / 16 referências missing**, **214 originais** preservadas. `deferred_definitions` aditivo mantém `type`/`name`/`definition`/`reasons`; `source_tree`, `xml_sha256`, origens/pins/bytes/fontes continuam. Nomes das oito entradas ficam em [001/FR-058](../001-gerador-jornada/spec.md#builtin-filtrado-fr-058), sem lista divergente ou Kind novo.
- Skin **`legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a`**, manifest/Profile/project/catalog/produto **1/1/3/2/1.0.0.0** inalterados. As 16 famílias continuam; total de definições não é quantidade de opções de cada seletor.
- Manter biblioteca ativa completa, não subset por uso nem remoção de fontes/bytes sem uso. Perfil próprio, revisão arquivada e documento existente conservam seu conjunto original; sem migração extra, poda, upgrade, enriquecimento builtin ou fallback. A redução de novo builtin não é política de conversão de documentos.
- [T017](tasks.md#t017-builtin-filtrado) verifica a disponibilidade por identidade e a preservação local, consumindo a integração de [001/T132](../001-gerador-jornada/tasks.md#t132-builtin-filtrado). Domínios e adaptadores mantêm seus donos; esta é verificação mista, não segunda implementação ou encerramento da paridade.

**Evidência delimitada:** **19 Python do importador PASS**, incluindo descoberta/execução dos quatro novos métodos após corrigir seu aninhamento; **4 Python de empacotamento e importador `--check` PASS**. Regressões Go de biblioteca/seletores/Capture/perfis próprios e arquivados implementadas; gofmt/vet e **compile-only de 15 suítes PASS, ZERO execuções Go**, sem contorno da política Windows. Build/pacote e integridade somente leitura PASS, dois EXEs idênticos; identidades reais na [seção 14](../001-gerador-jornada/quickstart.md#14-builtin-filtrado). Não herdar hashes antigos nem equiparar compilação a execução. Roteiros de inserção/seletor e projeto anterior/próprio/arquivado, UI nativa/DPI/Office e aceite continuam pendentes.

Reimportar em change futura com recursos exatos, origem/hash/licença, regeneração e validação de composição/estados/consumidores. Preservar fontes/licenças e gates `DELEGATE_MISSING`/`FIDELITY_UNVERIFIED`/Java/PPT; 31 diagnósticos não significam acervo sem bloqueios ou 100% fiel. Caixas antigas, checklist 13/2, manuais/DPI/Office e aceite permanecem inalterados.

<a id="plano-codigos-barras-2026-10-23"></a>
## Plano de domínio — códigos de barras (2026-10-23)

FR-013–FR-015 são incremento novo; [031](../031-codigo-barras-efl/spec.md) e [032](../032-codigo-barras-processado-efl/spec.md) detalham comportamento, não novos pacotes/domínios. Fonte dos dados/enums e premissas matemáticas no [contrato](../001-gerador-jornada/contracts/component-catalog.md#contrato-codigos-barras-2026-10-23), decisões na [pesquisa coordenadora](../001-gerador-jornada/research.md#pesquisa-codigos-barras-2026-10-23).

1. T018 prepara união tipada/defaults/capacidades/cópia em `internal/catalog/model.go` e `copy.go`; extrair propriedades para novos arquivos coesos em `internal/catalog` se necessário, sem criar pacote financeiro ou importar UI/IO. `barcodeInput` e `barcodeProcessed` não são variantes de Report nem oito textFields soltos.
2. T019 implementa cálculo/conversão/extração puros em arquivo planejado `internal/catalog/barcode.go` com testes vizinhos. Validar comprimento e caracteres antes de recortar; usar centavos inteiros e invalidar derivados após qualquer alteração de entrada. Não copiar os bugs 48/47, campo5 validando campo1 ou cache inconsistentes da classe de exemplo. SIGCB segue premissa expressa, sem afirmar homologação.
3. T020 implementa validação/cópia/defaults de blocos processados e fixtures de independência em `internal/catalog/barcode_processed.go` (planejado) e testes vizinhos; texto literal e seletores são dados, não consultas ou transformação fiscal. Modelo de T018 estabilizado permite T019/T020 em paralelo em arquivos distintos.

002/T013–T014 traduz wire; 004/T015–T016 cuida de instâncias/IDs/edição. Aplicação/desktop, cena e PPT ficam em 001/T135–T139. Nenhum algoritmo de DV deve ser refeito em Fyne ou no exportador. Tipos novos só são oferecidos quando a cadeia de suporte e compatibilidade estiver pronta. Testes de domínio, revisão visual e aceite são verificações distintas; nenhuma caixa concluída por esta atualização.

<a id="plano-menu-aplicacao-2026-10-23"></a>
## Plano de domínio — Menu aplicação

FR-016/T021 estende a fundação 003/T018, sem reabrir/duplicar T019/T020. Arquivos próprios planejados `internal/catalog/application_menu.go` e `internal/catalog/submenu.go` separam menu, submenu e regras; registro de Kind/Properties/cópia fica nos arquivos centrais existentes. API acordada no [modelo](data-model.md#modelo-menu-aplicacao-2026-10-23). Factory pura recebe callback de UUIDs e cria seed com enabled=true; nenhuma geração de identidade via UI, importação Java/XML em runtime ou pacote financeiro.

Validação/cópia percorrem árvore com limites e detecção de ciclo antes de recursão ilimitada, preservam strings literais, IDs e ordem; verificam inicial existente, layouts fechados e exclusão submenu/destino. Unicidade global no documento é aplicação, existência de tela é Jornada. Mudanças centrais serializadas com T018; depois arquivos próprios podem avançar sem concorrer com barcode.go/barcode_processed.go. Catálogo 3 cobre os três compostos. Defaults geométricos do menu devem vir do Java fornecido; valores não disponíveis não viram medidas inventadas. Constituição I/II/VI/VII atendida no desenho; implementação e evidência ainda pendentes.

## Fonte integral recuperada — 1.0.0.2

O gerador `scripts/import_legacy_menu.py` converte `material-legado/menu/menu.xml` em `internal/catalog/menu_legacy_generated.go`; `menu_defaults.go` cria propriedades e UUIDs. O Catálogo continua sem dependência de filesystem, Fyne ou XML em runtime. `menu_reports.go` sintético foi removido. `menu_xml_test.go` lê a referência independente somente em teste. Skins são resolvidas pela apresentação, conforme 033/FR-012.

## Defaults conferidos com o Java — 1.0.0.3

A fábrica pura preserva Title ausente e referencia skins financeiras/digitais e marcas reutilizadas no Bolão. PnlMenu.java/UIMenuControle.java ficam preservados em material-legado/menu como referências de desenvolvimento. Catálogo não carrega XML ou Java em runtime. Dados de movimento continuam referências inertes; conferir cadastro externo depende da fonte real ainda indisponível.

## Carregamento antecipado — integração 035/036

A inicialização enumera os recursos operacionais por skin.Library.ResourcePaths, confere integridade e aquece caches de imagens/fontes consumidos por canvas/prévia. Os defaults dos dez Kinds são validados sem publicar um documento fictício. Nenhum atributo de domínio/catálogo ou número de versão muda por esse cache; núcleo do ciclo fica em internal/startup, adaptador de carga em desktop/startup_resources.go e tarefas em 036.
