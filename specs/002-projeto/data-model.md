# Data Model: Domínio Projeto

## Extensão planejada — recursos portáteis 038

O [modelo 038](../038-integracao-mcp-agentes/data-model.md) mantém Project/Asset e o codec vigente como fonte de verdade. A biblioteca contém template `.jornada` de uma tela com somente recursos selecionados; instanciação transfere imagens com novos IDs quando necessário e deduplicação por hash, preservando integridade. Sessão/revisão/consentimento não são novos campos persistidos. Nenhum schema/catalogVersion é alterado por MCP; evolução futura exige contrato próprio. Desenho de 2026-10-25, ainda não implementado.

## Situação corrente — revisão 1.0.0.1

Project contém ID, Name, Revision, CatalogVersion e Assets; Document agrega telas, jornada, roteiros, JavaOptions e Resources. Writer 3/3; pares de leitura reconhecidos: 1/1, 2/1, 3/2 e 3/3. Campos desconhecidos são recusados. Perfis visuais privados por documento não foram implementados.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/project`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Project contém id, name, revision, catalogVersion e assets. Asset contém id, path, mediaType, sha256 e sizeBytes. A aplicação compõe Project com screens, journey, scenarios, javaOptions, Resources e `AppearanceProfile *skin.Profile` em Document; perfil não pertence a project.Project. O DTO de persistência é achatado conforme o contrato compartilhado. Revisão é incrementada por comandos confirmados de aplicação; operações que falham não publicam nova revisão.

## Invariantes

UUID v4 canônico; nome não vazio; revision inteiro não negativo; unicidade global de IDs; arrays mantêm ordem. Recursos referenciados não são removidos implicitamente. Projeto vazio é salvável, mas não executável. Erro estrutural impede nova sessão; referência quebrada legível permanece reparável.

## Limites e relacionamento com outros domínios

A aplicação resolve referências entre módulos e constrói snapshot. O adaptador projectzip converte DTOs e bytes; Project não importa screen/journey/scenario para executar suas regras nem seus adaptadores. A composição do documento é montada na aplicação, evitando um agregado com todos os comportamentos.

Não implementa filesystem, ZIP, autosave, undo, UI, decodificação de imagens ou geração de artefatos.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: IDs duplicados, nome vazio, revisão negativa, recurso compartilhado e isolamento de snapshot. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/project-format.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-001, FR-009, FR-010; NFR-004, NFR-005 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.

## Versão de arquivo e identidade de saída

O manifesto atual é escrito em schemaVersion 3/catalogVersion 3. O reader reconhece os pares 1/1, 2/1, 3/2 e 3/3 com seus campos suportados, sem sobrescrever o original na leitura. FileID pertence a Screen/Component do domínio Tela, não ao UUID de Project. Software continua 1.0.0.0. Schema 2/catalog 1 é histórico; preservar referências UUID e recursos durante a migração.


## Correspondência histórica — 2026-10-18

Identidade e metadados implementados; snapshot, manifesto v2/migração v1 e persistência integrados. Remoção segura de recursos e matriz completa de falhas/recuperação ainda precisam de cobertura; não concluir T006/T007 apenas pelo round-trip.

## Perfil visual planejado e lacunas — referência de 2026-10-19

O perfil por documento descrito nos registros de 2026-10-19 é histórico/plano de expansão; não existe na revisão 1.0.0.1. A implementação atual usa Style e biblioteca incorporada global, sem captura/selamento de revisões privadas. O histórico de build e validação continua no quickstart; não o repetir como evidência vigente deste modelo.

A expansão continua exigindo perfil privado imutável, definições/estados, faces/aliases/recursos, hash determinístico e revisões arquivadas com limites, cópia profunda e resolução exata. O alvo de bytes é o subset usado e autorizado; hash ou JSON declarando licença não concedem direito de redistribuição. `appearance`/overrides com presença/RGBA e descritores de recursos exigem DTO, compatibilidade e testes antes de entrar no wire.

Assets do documento atual são PNG/JPEG. TTF/JSON de perfil dentro do pacote, Capture/ResourceData/ValidateProfile e seus contratos completos ainda não são suportados. A seção de situação corrente descreve o writer real; nenhuma migração pode trocar IDs, bounds ou ordem.

<a id="modelo-codigos-barras-2026-10-23"></a>
## Persistência planejada — códigos de barras (2026-10-23)

FR-011 não acrescenta entidade de boleto a Project nem cópia de blocos em Journey. As propriedades da instância de Tela são traduzidas pelo adaptador a partir do [modelo de Catálogo](../003-catalogo/data-model.md#modelo-codigos-barras-2026-10-23). Rótulo de entrada é persistível; seus oito valores capturados e derivados não são. No processado, `blocks[].id/title/rows[].id/labelKind/value` e ordem são configuração persistível, inclusive valores vazios e Unicode.

Alvo de escrita 3/3 e matriz de leitura/migração no [contrato de formato](../001-gerador-jornada/contracts/project-format.md#formato-codigos-barras-2026-10-23). Não alterar schema do perfil visual ou a versão do produto por associação com catalogVersion; estas versões têm ciclos independentes. Snapshot/cópia profunda preserva IDs e não compartilha listas; duplicação de identidades pertence a Tela. Sem novas entradas ZIP ou recursos de hardware, sem valores pessoais em logs de falha.

<a id="modelo-menu-aplicacao-2026-10-23"></a>
## Persistência planejada — Menu aplicação

FR-012 não adiciona agregado ao Project. `applicationMenu` persiste propriedades tipadas `root` e `initialMenuId`; submenus/itens são tradução camelCase do [modelo canônico 003](../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23), detalhada no [contrato wire](../001-gerador-jornada/contracts/project-format.md#formato-menu-aplicacao-2026-10-23). Ordem, IDs e enabled=false são conteúdo; ausência de campo obrigatório não recebe default silencioso. Menus vazios representam `items: []`.

TargetScreenID persiste referência autorada, cuja validade é de Jornada; não duplicar sua rota em Action nem gerar bindings por label. Nenhuma entrada ZIP ou versão de perfil nova. Schema3/catalog3 é alvo conjunto 031/032/033, não migração adicional; matriz/backup pertencem a FR-011/T013–T014. Stack/caminho, seleção simulada e exit flag nunca entram no manifesto.
