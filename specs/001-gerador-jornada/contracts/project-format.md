# Contract: formato de projeto .jornada

## Implementação conferida — 1.0.0.1

Writer vigente: **schemaVersion 3/catalogVersion 3**. Pares reconhecidos: **1/1, 2/1, 3/2 e 3/3**, com rejeição de campos desconhecidos. `wire.go` não aceita `appearanceProfile` nem `appearance`; esse perfil, seu selamento, arquivos de fontes e revisões arquivadas descritos na expansão abaixo são alvo/histórico e não capacidade atual. Projetos antigos com campos não suportados podem ser recusados, mesmo com par de versões reconhecido; não substituir o original nesses casos.

Raiz efetivamente escrita: `id`, `name`, `revision`, `catalogVersion`, `assets`, `schemaVersion`, `screens`, `journey`, `scenarios`, `javaOptions`. Instâncias gravam `id`, `name`, `fileId`, `kind`, `bounds`, `style`, `enabled`, `parentId`, `properties`. Style contém somente `fontFamily`, `fontSize`, `textColor`, `fillColor`, `borderColor`, `backgroundAssetId`, `textAlign` e `skin`; valores zero são omitidos. Não há override de alpha em Style. Biblioteca/fontes de skin acompanham o executável, não são bytes privados incorporados no documento.

Salvar catálogo antigo sinaliza migração e preserva cópia `.bak` do destino existente antes da publicação. Testes atuais em `projectzip/{package,skin_wire,compound}_test.go` cobrem fatias desses percursos; a matriz completa de disco/recuperação e perfis permanece aberta. Exemplos v1/v2 abaixo são históricos e regras de compatibilidade; não representam o writer vigente.

Código atual: [wire do projeto](../../../apps/gerador-jornada/internal/adapter/projectzip/wire.go), [cena](../../../apps/gerador-jornada/internal/presentation/scene.go) e [exportador PPT](../../../apps/gerador-jornada/internal/adapter/pptx/export.go). Evidência executada na [seção 22 do quickstart](../quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). O restante do contrato conserva regras normativas e evolução planejada; divergências apontadas aqui não foram apagadas dos requisitos.


**Version**: 3 implementada, reader1/2/3, catalogVersion2; perfil schema1 e transição/normalização parcial em 2026-10-19 | **Consumers**: editor e adaptador projectzip | **Requirements**: FR-001, FR-009, FR-010, FR-051, NFR-005.

As seções anteriores à expansão abaixo documentam o wire format v1/v2 e sua migração **histórica**, não instruem o writer atual a gravar v2. A expansão distingue o v3 implementado do contrato-alvo de normalização ainda pendente e prevalece para a escrita atual; não torna leitores antigos compatíveis.

## Pacote e wire format

Arquivo ZIP com `manifest.json` UTF-8 na raiz e imagens em `assets/<UUID>.png` ou `.jpg`. Sem links externos, código, templates arbitrários, JARs ou scripts no projeto.

Manifesto usa camelCase, sem funções. Os campos de raiz são:

| Campo | Tipo | Regra |
| --- | --- | --- |
| schemaVersion | integer | Writer: 2. Reader: 1 com migração, ou 2 com validação de fileId. |
| catalogVersion | integer | Exatamente 1. |
| id, name | string | UUID e nome não vazio. |
| revision | integer | Inteiro não negativo; identifica a revisão salva. |
| screens | array | Screen na ordem de apresentação. Pode ser vazio em rascunho. |
| journey | object | startScreenId (UUID ou string vazia quando ausente; reader também aceita null) e actions (array). |
| scenarios | array | Scenario; vazio quando não há demonstração. |
| assets | array | Asset; vazio quando não há imagens. |
| javaOptions | object | packageName e profileId, strings; podem estar vazios no rascunho. |

Screen contém `id`, `name`, `fileId`, `size:{width,height}`, `background:{color,assetId}` e `components`. Component contém `id`, `name`, `fileId`, `kind`, `bounds:{x,y,width,height}`, `style` e `properties`. assetId vazio significa sem imagem; o reader também aceita null. Style permite os campos descritos no modelo, omitindo os não personalizados.

Todos os arrays são obrigatórios, inclusive vazios. Propriedades e campos de outras entidades seguem [catálogo](component-catalog.md) e [jornada/roteiro](navigation-scenario.md). Campos desconhecidos são recusados pelo leitor, com localização; não são ignorados silenciosamente.

### Exemplo mínimo válido como rascunho

```json
{
  "schemaVersion": 2,
  "catalogVersion": 1,
  "id": "8be7fa20-bf51-4d26-9bc9-cae9fd075d01",
  "name": "Atendimento",
  "revision": 0,
  "screens": [],
  "journey": { "startScreenId": "", "actions": [] },
  "scenarios": [],
  "assets": [],
  "javaOptions": { "packageName": "", "profileId": "" }
}
```

Esse arquivo pode ser salvo/reaberto; exportação retorna `JOURNEY_START_MISSING`.

## Ausência de referências opcionais no wire format

O writer atual representa ausência em campos Go string por `""`, incluindo
startScreenId, targetScreenId de ação none, selectedOptionId, assetId e alvos opcionais.
O reader também aceita `null` para esses campos, normalizando-o para a string vazia.
Nos modelos conceituais, “nulo” descreve ausência de referência, não exige JSON null.
Opcionais tipados por ponteiro, como Step.point e Step.text, conservam JSON null;
texto vazio de typeText continua sendo um valor explícito diferente de null.
Essa descrição registra o comportamento do writer/reader já existente, sem migração
adicional nem mudança de schemaVersion 2.

## Recursos e limites de leitura

Asset contém `id`, `path`, `mediaType`, `sha256` hexadecimal minúsculo e `sizeBytes`. Aceitar somente PNG/JPEG decodificáveis e dimensões até 8192 × 8192. Sem SVG executável ou URLs remotas nas versões 1 e 2.

Limites iniciais para leitura: manifesto descomprimido até 10 MiB, recurso até 25 MiB, pacote até 250 MiB descomprimidos no total e até 1000 entradas. Verificar os bytes efetivamente lidos, não apenas o tamanho informado no cabeçalho. Mensagem deve indicar o limite excedido.

Recusar caminhos absolutos, unidades Windows, `..`, barras invertidas, nomes duplicados ou equivalentes sem distinção de maiúsculas, symlinks e entradas fora de manifest/assets. Nenhuma extração pode escapar da área temporária. Hash e tamanho precisam corresponder aos bytes.

Versão desconhecida gera `PROJECT_VERSION_UNSUPPORTED`; corrupção/ZIP inválido, `PROJECT_CORRUPT`; tipo/campo desconhecido, `PROJECT_SCHEMA_INVALID`. Erros de referência em conteúdo legível podem abrir rascunho com diagnósticos de domínio.

## Gravação e publicação

1. Capturar snapshot da revisão a salvar e validar estrutura/recursos. Permitir diagnósticos semânticos para salvar rascunho.
2. Escrever ZIP temporário no mesmo volume do destino, fechar streams, sincronizar e reler sua integridade.
3. Publicar com operação de substituição atômica compatível com Windows. Se não houver garantia de substituir com segurança, falhar preservando o arquivo anterior.
4. Atualizar revisão salva somente após publicação. Edições posteriores à revisão capturada permanecem pendentes.

Solicitar confirmação apenas para substituir um destino escolhido que já exista; autosave não dispara esse diálogo. Falta de permissão, disco cheio e arquivo bloqueado são erros recuperáveis.

Ordenar entradas e manter metadados ZIP normalizados para comparação determinística. Não gravar seleção, cursor, zoom, janela, undo/redo ou valores runtime.

## Autosave e recuperação

Após comando confirmado, debounce de 2 s grava cópia no diretório de recuperação do usuário: Windows `%LOCALAPPDATA%/GeradorJornada/recovery/<projectId>.jornada`; Linux `$XDG_CONFIG_HOME/GeradorJornada/recovery` (ou `~/.config`); macOS `~/Library/Application Support/GeradorJornada/recovery`. A resolução pertence ao adaptador; erro não deve cair em diretório relativo. Máximo uma cópia válida por projeto, substituída com o mesmo procedimento seguro.

Autosave não sobrescreve arquivo explícito nem marca sessão limpa. Em interrupção, oferecer recuperar a revisão mais recente válida; recuperação abre como alterado e exige salvar explicitamente. Descartar recuperação requer ação do usuário. Autosave falho é indicado sem interromper edição.

Salvar/fechar com sucesso pode remover a recuperação apenas quando ela não contém revisão posterior à salva. Em diretório sem permissão, explicar que recuperação automática está indisponível.

## Testes de contrato

<a id="formato-codigos-barras-2026-10-23"></a>
### Extensão planejada de catálogo para códigos de barras (2026-10-23)

**Planejado, não implementado nesta entrega documental.** O writer vigente continua schema 3/catalog 2; o incremento de [002/FR-011](../../002-projeto/spec.md#codigos-barras-2026-10-23) tem alvo **schema 3/catalog 3**, preservando envelope e limites ZIP/perfil. Novas famílias e campos tipados estão no [contrato de Catálogo](component-catalog.md#contrato-codigos-barras-2026-10-23). A versão de catálogo muda porque o leitor antigo recusa Kinds desconhecidos; não apresentar o novo conteúdo como catálogo 2 compatível. Não mudar schema do perfil, versão de software ou bytes de recursos por associação.

| Par de entrada schema/catalog | Leitor planejado | Próxima escrita planejada |
| --- | --- | --- |
| 1/1, 2/1 | Migrações históricas preservadas | 3/3, com aviso/cópia íntegra antes de sobrescrever |
| 3/2 | Preservar conteúdo, geometria e perfis exatos | 3/3, com aviso/cópia por incompatibilidade de catálogo |
| 3/3 | Validação estrita dos tipos/enum/IDs | 3/3 |
| Outros | Recusar versão/catálogo incompatível sem alterar sessão/original | Nenhuma |

`barcodeInput.properties` contém apenas o rótulo configurado. Seus campos capturados, foco, leitura parcial, DV, estado concluído, erros, caches e derivados não são serializados. `barcodeProcessed.properties.blocks` é lista ordenada de blocos `id/title/rows`; linhas contêm `id/labelKind/value`. Strings vazias de título/valor são válidas; arrays mínimos, UUIDs e enums seguem Catálogo; não inferir valores fiscais ou IDs pelo conteúdo. Campos de instância/estilo não mudam de dono nem são duplicados dentro de Properties.

Writer não inclui rótulo derivado de enum, popup do inspector ou referência implícita entre os compostos. Reader recusa campo/enum desconhecido sem descarte ou conversão para Report. Abrir não escreve; salvamento migrado conserva original, preserva IDs/ordem/Unicode, não redimensiona e não enriquece perfil próprio/arquivado com builtin. Não há exportação reversa para catálogo 2. A migração só altera o identificador de compatibilidade e traduções contratadas; não insere compostos em arquivos anteriores.

Testes locais em [002/T013–T014](../../002-projeto/tasks.md#tarefas-codigos-barras-2026-10-23); integração e resultado observado em [001/T138–T139](../tasks.md#tarefas-codigos-barras-2026-10-23). Esta seção supera os alvos históricos de escrita 3/2 **somente para o incremento futuro**, não descreve o executável corrente.

<a id="formato-menu-aplicacao-2026-10-23"></a>
### Extensão conjunta planejada 031/032/033 — árvore de menu

Aplicar a matriz de leitura **1/1, 2/1, 3/2 e 3/3** e escrita alvo **schemaVersion=3/catalogVersion=3** da seção anterior aos três novos tipos. Não mudar schema para4 nem catalogVersion para4; base anterior validada3/2 permanece identificada como histórica até validação principal. Primeiro save incompatível conserva original íntegro antes da publicação; falha de cópia impede sobrescrita, abrir nunca regrava. Mesmo envelope/allowlist/limites ZIP/perfis, sem nova entrada ou recurso executável.

Para `kind: applicationMenu`, o corpo `properties` corresponde a `Properties.Menu` do Go (não um segundo invólucro `properties.menu`):

| Objeto wire | Campos | Semântica |
| --- | --- | --- |
| properties | `root`, `initialMenuId` | root objeto obrigatório; inicial string UUID existente obrigatório. |
| submenu | `id`, `name`, `title`, `layout`, `items` | Strings obrigatórias; items array obrigatório, inclusive []; root usa root, descendentes standard/financial/bolao. |
| item | `id`, `label`, `legend`, `category`, `variant`, `enabled`, `submenu`, `targetScreenId`, `movementCode`, `movementOption`, `movementName`, `shortcut` | Todos os campos textuais e enabled obrigatórios; submenu opcional ausente/null ou objeto válido. legend é conteúdo ilustrativo literal, inclusive Unicode/vazio, independente de movementName e sem fallback para metadados. |

Writer emite strings vazias quando valor vazio é permitido e booleano enabled sempre, inclusive false; não usa omitempty para perder configuração. TargetScreenId vazio é folha sem destino; com submenu não pode ser não vazio. Reader não inventa habilitação/defaults a partir de ausência nem infere inicial por título. Items=null não equivale ao array vazio autorado. Strings name/title/label/category/metadados admitem Unicode/vazio; IDs/layouts seguem invariantes. Variant desconhecida é preservada como string e diagnosticada pela capacidade visual, não aproximada ou executada.

Validar campos estritos, UUIDs globais, inicial resolvido, profundidade<=16/itens<=1000 por árvore e exclusão submenu+destino antes de publicação/tradução recursiva sem limite. Integridade de TargetScreenId é Jornada; conteúdo estruturalmente legível pode seguir política geral de rascunho, mas não ser exportado como válido. Não adicionar Action duplicada por item; referência persistida única fica na árvore. Caminho atual, modo da vista, foco/pressionado, seleção simulada, exit flag e expansão do inspector são proibidos no manifesto.

002/T015 estende o codec e fixtures existentes, não replica migração/backup002/T013–T014. Round-trip inclui falso explícito, menus vazios, todos os campos de metadados, inicial em descendente, IDs/ordem/Unicode, limites16/17 e1000/1001 e rejeição sem troca de sessão. Nenhuma alteração de versão de software, schema do perfil, bytes de recursos ou enriquecimento de revisões por associação com catálogo3.

### Matriz histórica da base

Round-trip de todas as famílias; preservação de IDs/ordem; campos desconhecidos; versão futura; JSON truncado; UUID repetido; path traversal; recurso alterado; limites ZIP; falha antes da substituição; arquivo bloqueado; recuperação e edição concorrente ao salvamento.

## Nome-base de arquivos e migração v1 → v2

Screen.fileId e Component.fileId são obrigatórios na versão 2. O campo `id` continua sendo UUID interno, e `name` continua sendo nome amigável Unicode. FileID não substitui UUID em ações, roteiros, recursos ou vínculos de teclado.

Regra comum: `^[A-Za-z][A-Za-z0-9]{0,39}$`, sem extensão. Telas exigem adicionalmente prefixo exato Pnl e ao menos um caractere após ele; prefixo conta no limite de 40. Componentes não exigem Pnl. Recusar dispositivos Windows CON, PRN, AUX, NUL, COM1–COM9 e LPT1–LPT9 sem distinguir caixa. Unicidade: telas por projeto e componentes por tela, sem distinguir caixa.

O leitor aceita v1 e preenche somente fileIds ausentes, sem alterar UUIDs ou conteúdo. Sugestão de tela remove a ordenação numérica inicial e acentos: `01 · Identificação` → `PnlIdentificacao`. Sugestões vazias usam PnlTela/Componente; colisões de sugestões recebem sufixos 2, 3, … com espaço reservado dentro de 40 caracteres. Valores explícitos não são renomeados. O próximo salvamento grava v2; o arquivo original não é escrito durante a leitura. V2 com valor ausente, inválido ou duplicado gera PROJECT_SCHEMA_INVALID; versões fora de 1/2 geram PROJECT_VERSION_UNSUPPORTED. Software 1.0.0.0 e catalogVersion 1 permanecem independentes da versão do manifesto.

Exercitar migração, estabilidade das sugestões, UUID/navegação preservados, round-trip v2, colisões de caixa, 40/41 caracteres, acentos, caminhos e nomes Windows reservados.

## Expansão pós-v1 — schema 3 e catálogo 2 (2026-10-19)

### Compatibilidade e preservação

Tabela reconciliada com o reader/writer 1.0.0.1. O perfil privado descrito no plano histórico abaixo não é acrescentado durante a migração atual.

| Entrada | Leitura pelo novo editor | Escrita seguinte |
| --- | --- | --- |
| schema 1/catalog 1 | Sugestão de FileIDs e campos compatíveis; nenhum overwrite ao abrir | schema 3/catalog 3 |
| schema 2/catalog 1 | Preservar FileIDs e campos compatíveis; nenhum overwrite ao abrir | schema 3/catalog 3 |
| schema 3/catalog 2 | Campos suportados; recusar appearanceProfile ou outro campo desconhecido, sem perda do original | schema 3/catalog 3 |
| schema 3/catalog 3 | Validação estrita dos campos e tipos atuais | schema 3/catalog 3 |
| Demais versões/combinações | PROJECT_VERSION_UNSUPPORTED ou diagnóstico de catálogo incompatível | Nenhuma escrita |

Migração preserva UUIDs/FileIDs explícitos, nomes/textos/opções/linhas, ordem, vínculos de teclado/navegação/roteiro, bounds e overrides. Defaults novos só se aplicam à criação, não à geometria de arquivos antigos: numericKeyboard 200 × 280 continua assim; novas inserções usam 200 × 350. Não inferir override comparando valor ao default atual: campo explicitamente gravado continua explícito mesmo se igual ao default antigo. Traduzir Style legado para overrides com presença; textAlign ausente/vazio antigo vira left explícito para preservar semântica histórica. Nenhum ParentID é inferido de sobreposição.

Primeiro salvamento migrado informa que editores antigos não abrem v3 e preserva cópia íntegra do pacote original antes da publicação atômica, sem substituir backup existente silenciosamente. Falha ao preservar a cópia impede sobrescrever o original; usuário pode escolher outro destino. Leitura nunca regrava. Não há requisito de exportação reversa v3→v2 nem importação `.page`.

### Campos e recursos implementados — continuação de 2026-10-19

O writer atual grava schema 3/catalog 2 e captura `application.Document.AppearanceProfile`
antes de serializar. O campo **pertence a Document, não a project.Project**. Na leitura,
`appearanceProfile` é opcional no schema 3 durante a transição; quando presente deve
ser objeto não nulo válido. Schemas 1/2 recusam esse campo. Ausência permite usar
builtin; presença inválida nunca é substituída por builtin.

O DTO implementado é `skin.Profile`: `schema=1`, `id`, `version`, `sha256`,
`definitions`, `fonts`, `aliases`, `resources`, `diagnostics`, `archived`.
`archived` admite até 16 revisões diretas, com versões únicas, sem novo nível.
Cada recurso contém `info` (caminho lógico/hash/tamanho/MIME e demais metadados),
`assetId`, `license` (`unknown/allowed/restricted`) e `status`
(`embedded/referenced/blocked/missing`). Caminho lógico não é caminho físico.

Hash: SHA-256 hexadecimal minúsculo do JSON compacto de `profilePayload`, na ordem
de campos declarada no código, mapas ordenados por `encoding/json`, ordem dos slices
significativa e `nil` distinto de coleção vazia. Exclui **somente `sha256` da própria
revisão**; inclui todos os dados de `archived` e seus hashes individuais. `Seal`
calcula primeiro os hashes arquivados, depois a raiz; `ValidateProfile` confere
estrutura/referências/limites/hashes. Integridade não comprova licença ou procedência.

`skin.Library` mantém snapshot privado imutável, com cópia profunda nas saídas.
`ForVersion` aceita vazio/versão corrente ou versão arquivada direta exata; versão
desconhecida retorna indisponível, sem fallback. `Lookup` exige estado exato.

`visualprofile.Capture` conserva **todos os metadados da biblioteca original** em
cópias privadas: definições, fontes, aliases, recursos e todas as revisões arquivadas,
inclusive versões não usadas e todos os estados. Não enriquece perfis privados com
builtin nem faz upgrade. O uso corrente limita **somente diagnósticos e acréscimo
de entrada missing**, não a retenção de metadados. Recursos builtin conservam
metadados `referenced/unknown`; faltas e overrides desconhecidos recebem diagnóstico.
Conserva assets/bytes já incorporados, **não acrescenta binários sem autorização**,
mesmo se JSON declara `allowed`.

**Desvio arquitetural justificado — bug crítico save/reopen → AddComponent:** podar
o catálogo quebrava a inserção de Button com `defaultconfirm` após reabrir. Snapshot
integral de metadados preserva inserir famílias e trocar variantes no futuro, todos
os estados e namespace customizado, sem completar perfil privado com builtin ou
atualizar sua revisão. Não é incorporação integral binária: o alvo normativo SUBSET
de bytes usados e permitidos permanece. O payload selado mais amplo altera o hash;
nenhum hash novo é declarado aqui. Autorização jurídica e transporte automático
continuam pendentes; normalização appearance/overrides tipados/RGBA também.

`ResourceData` resolve `embedded` pelo `AssetID` único do documento, conferindo
hash/tamanho/MIME de Asset, perfil e bytes reais. `referenced` só pode usar bytes
builtin com hash/tamanho/MIME exatamente coincidentes, também verificados por
conteúdo. Revisão desconhecida, conteúdo diferente, recurso ausente/bloqueado ou
vínculo inválido falham; não há busca por homônimo, fonte instalada ou filesystem.
TTF passa pelo parser SFNT; imagens aceitas por essa fronteira são PNG/JPEG e GIF
estático de um frame validado. Isso não amplia indiscriminadamente a allowlist ZIP.

**Mapeamento transitório e lacunas mantidas:** componente ainda persiste `style`,
incluindo `skinName`, `skinVersion`, `fontFamily`, tamanho, RGB, `fontItalic`,
`fontUnderline`, `spacement`, `disabled` e `fillAlpha` opcional. `FillAlpha=nil`
herda e zero é transparente; os demais campos não fornecem ainda presença tipada
completa. Família desconhecida é preservada literalmente em `Style.FontFamily`,
com diagnóstico, não em um descritor `requestedFamily/unresolved` implementado.
`Fonts/Aliases/Resources` atendem ao perfil atual, mas não são os nomes/descritores
normalizados `fontFaces/resourceIds` abaixo. O writer **não emite ainda**
`component.appearance` normalizado com `style` vazio, nem normaliza todas as cores
para RGBA. Essas diferenças permanecem trabalho de 002/T009–T012 e 004/T009;
não se altera o alvo para ocultá-las.

### Campos e recursos novos (contrato alvo ainda não integralmente implementado)

O desenho histórico propunha `schemaVersion=3`, `catalogVersion=2` com `appearanceProfile` obrigatório; essa proposta não foi entregue nesta árvore. O formato atual 3/3 já grava somente os campos conhecidos. Introduzir o perfil no futuro exige rever o contrato de versão e compatibilidade, sem atribuí-lo aos pares atuais por documentação. O perfil planejado contém com `id`, `version`, `sha256`, `definitions`, `fontFaces` e `resourceIds`. Essa obrigatoriedade e esses descritores não descrevem o reader transitório atual acima; qualquer mudança futura do wire exige compatibilidade explicitada. definitions contém snapshot tipado das variantes usadas com todos os seus estados/dependências, não XML/código. fontFaces contém descritores tipados das fontes usadas pela skin e por overrides, inclusive sem associação a uma skin. sha256 cobre a representação canônica determinística do perfil, excluindo somente o hash próprio, incluindo revisões arquivadas. resourceIds referencia assets realmente incluídos. Metadados de proveniência/licença acompanham definições/recursos, sem caminhos absolutos.

A restrição de metadados às variantes usadas no desenho original acima **não
descreve Capture corrente**: o desvio arquitetural documentado mantém todos os
metadados originais para editabilidade futura. O SUBSET normativo continua aplicado
aos bytes usados e permitidos. Não remover os alvos de normalização dos descritores,
appearance/overrides tipados/RGBA por causa desse desvio.

Descritor de face tem `id`, `resolutionStatus` (resolved/unresolved/blocked), `resourceId` quando resolvido, nome interno/família/subfamília/estilo/peso/métricas comprovados e origem/licença/permissões. `requestedFamily` opcional conserva literalmente a família solicitada por override legado, mesmo que ainda não resolvida. Não inventar nome interno/métricas ou incluir asset fictício quando irresolvido; os campos dependentes ficam ausentes e o diagnóstico localiza o descriptor/componente. Resolved exige bytes válidos e referência existente; sem permissão de incorporação o recurso não é incluído e fidelidade portátil fica bloqueada. Flags OS/2 e nome não comprovam isoladamente licença legal.

Migração de `style.fontFamily` explícito cria/reusa descritor preservando a string original e produz override `fontFaceId` presente. Só associar a face real por regra/metadados/evidência documentados, não por arquivo de nome parecido ou fonte disponível na máquina. Família desconhecida mantém descritor unresolved e permite preservar rascunho/round-trip, mas bloqueia fidelidade dependente; nunca vira Arial/default silenciosamente. Style antigo sem família não cria override artificial de fonte. Essa regra também vale para famílias que coincidem com o default antigo.

Component mantém seus campos, acrescenta `parentId` opcional (ausente/null/"" = raiz) e `appearance` opcional `{skinId, skinVersion, overrides}`. Referência explícita inválida não recebe fallback silencioso. `style` legado é compatível para leitura; escritor v3 emite aparência normalizada e style vazio, sem dois valores concorrentes. Conflitos em entrada v3 com propriedade presente em ambos são PROJECT_SCHEMA_INVALID. `overrides` usa campos tipados: fontFaceId, fontSize, fontWeight, italic, textColor, fillColor, borderColor, backgroundAssetId, textAlign, lineHeight, spacement. Ausência herda; zero/falso explícitos não são ausência e continuam sujeitos à validação de cada propriedade. Cores novas usam RGBA canônico #RRGGBBAA; #RRGGBB legado normaliza com FF, sem perder transparência.

IDs/versões/nome interno de fonte e permissões pertencem aos descritores; bytes não entram em definições. Asset conserva id/path/mediaType/sha256/sizeBytes e acrescenta categoria tipada (`image`, `font`) e metadados de origem/licença para os recursos novos. Imagens seguem PNG/JPEG e formatos adicionais apenas se inventariados, decodificáveis e contratados antes de habilitar; não habilitar SVG/código ou GIF animado indiscriminadamente. TTF válido usa extensão .ttf e mediaType font/ttf. Perfil é dado tipado no manifesto, não script nem caminho externo. Sem fontes instaladas obrigatórias para o editor.

Os limites de manifesto, entrada, descompressão e total do pacote anteriores também se aplicam a fontes/perfil; nenhuma exceção por ser recurso embutido. Conferir bytes, tamanho, hash, MIME, referências, parser TTF, chaves de skin e permissões de redistribuição. Arquivos/fontes não permitidos não são copiados automaticamente. Não carregar/executar conteúdo como código. SafePath continua recusando traversal, symlink, duplicidade de caixa e entradas fora da allowlist manifest/assets. **Alvo normativo binário SUBSET:** incorporar somente bytes de dependências usadas e permitidas, deduplicados, não todos os binários do EXE. Isso não limita o snapshot integral de metadados adotado pelo desvio arquitetural acima; nem JSON `allowed` autoriza cópia automática.

### Testes adicionais obrigatórios

Regressões novas: save/reopen/addButton (`AddComponent` com `defaultconfirm`), inserção
de outra família e troca de variant após reopen, preservação de transient states
nas definições (sem persistir estado runtime), custom namespace sem enriquecimento
builtin/upgrade e revisões não usadas conservadas em cópias privadas. Conferir que
uso só limita diagnósticos/missing, sem novos binários mesmo com JSON `allowed`.
Incluídas na compilação pelo main **após a correção integral de Capture**:
`gofmt`/`go vet ./...` PASS e `go test -c -o ./build/compile-tests/ ./...` PASS de
compilação de **15 suítes, não executadas**; bloqueio prévio por política Windows
preservado, sem PASS de testes Go. Python **28 PASS** (15 importador + 4 packaging +
9 Java 6, incluindo erro esperado de `String.strip`, não homologação EFL); importador
`--check` PASS, mesmo inventário 214/8/270/47. Novo EXE Windows amd64 construído;
distribuição Windows gerada e verificada; checksum do ZIP no sidecar externo
`.zip.sha256`. Reempacotar documentos preserva o EXE. Evidências,
hash/tamanho e manual pendente no [quickstart](../quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
A matriz de fixtures e o aceite continuam pendentes.

Fixtures v1/v2→v3 e round-trip de cada família/variante com recursos; backup/aviso/falha de publicação; igualdade de IDs/bounds/overrides/ordem/referências antes e depois; left histórico; campos zero/falso versus ausentes; catálogo/Kind/estado/campo desconhecido; perfil/hash/fontes corrompidos ou ausentes; licença; limites ZIP inclusive TTF; ParentID/ciclos/mesma tela; snapshot/reabertura estável sem o acervo original. Responsável: [002-projeto/tasks.md](../../002-projeto/tasks.md), T009–T012; integrações sem duplicar persistência em 001/T118.
