# Data Model: Domínio Tela

## Situação corrente — revisão 1.0.0.1

Screen mantém ID/Name/FileID, Size, Background e Components em ordem de camadas. Component mantém ID/Name/FileID, Kind, Bounds, Style, Enabled, ParentID e Properties. Style tem FontFamily/FontSize/TextColor/FillColor/BorderColor/BackgroundAssetID/TextAlign/Skin; não possui FillAlpha nem appearance tipada.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/screen`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Screen contém id, fileId, name, size, background e components. Component contém id, fileId, name, kind, bounds, style e properties. Bounds contém x, y, width e height; Style usa os campos comuns do contrato do catálogo. Não há árvore de componentes ou agrupamento persistido na v1.

## Invariantes

x/y inteiros não negativos; width/height inteiros positivos; x+width e y+height dentro da tela para validação executável. Nomes lógicos únicos na tela. Ordem do array é camada. Dimensões de drop são ajustadas na aplicação; regras recusam estado inválido de operação confirmada.

## Limites e relacionamento com outros domínios

Importa catalog para tipos e validação de família. Usa valores/IDs, não widgets. Aplicação agrupa alterações de tela e referências de Journey/Scenario num comando reversível único. Presentation consome geometria sem modificar entidades.

Não implementa Fyne, coordenadas de ponteiro, pilha undo/redo, serialização, roteiro ou exportação.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: bordas, dimensões inválidas, nome duplicado, z-order, cópia de teclado, auto-referência e resize de tela sem reposicionamento oculto. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-002, FR-004 a FR-008, FR-013, FR-016 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.

## Alinhamento textual e transformação coletiva

Style.TextAlign aceita vazio/left/center/right; vazio mantém o comportamento anterior. SelectionBounds e MoveSelection recebem IDs como argumento, sem adicionar seleção à entidade Screen. ConstrainTranslation limita o grupo inteiro mantendo distâncias. Rascunho de formulário não modifica o componente até a ação Salvar alterações.

## ID de arquivo

FileID é independente do ID UUID e do Name Unicode. Aplicar limites/prefixo de 004-tela/FR-011–FR-012, incluindo unicidade local dos componentes e, na aplicação, unicidade de telas. Campos existem na entidade de composição e não em Style ou Properties. A persistência usa fileId no manifesto v2 e a exportação usa o valor exato em nomes de arquivo/classe.


## Correspondência com o código atual

Modelo, bounds, seleção coletiva, IDs de arquivo, duplicação e remoção reversível implementados. Parte dos comandos de composição ainda é coordenada no desktop por Session.Change. A matriz completa de bordas/duplicação e todos os alinhamentos previstos ainda requerem conferência; DPI nativo permanece manual.

## Modelo alvo pós-v1 — 2026-10-19

**Alvo ainda não integralmente implementado:** `component.appearance`
(skinId/skinVersion/overrides tipados), presença distinguindo zero/falso, FontFace ID e
RGBA canônico seguem o [contrato](../001-gerador-jornada/contracts/project-format.md#expansão-pós-v1--schema-3-e-catálogo-2-2026-10-19).
**Implementado no modelo corrente:** ParentID, Enabled e Style com os campos listados no quadro inicial. ParentID é opcional e aponta à mesma tela, sem ciclos; bounds são absolutos e mover o pai desloca descendentes. A ordem canônica permanece em Components. Alpha de skin existe na primitiva, não em Style.

O perfil por documento descrito nos registros de 2026-10-19 é histórico/plano de expansão; não existe na revisão 1.0.0.1. A implementação atual usa Style e biblioteca incorporada global, sem captura/selamento de revisões privadas. Evidências antigas permanecem no quickstart e não comprovam o estado deste build.

**Expansão pendente:** appearance normalizado com referência de skin/versão, presença tipada e RGBA; grupos contíguos/recorte de contêiner e cobertura das famílias/variantes adicionais. Essas regras não introduzem callbacks/widgets ou fontes persistidas concorrentes.



## Ordem canônica e colocação de camada — 2026-10-20

FR-018: `Screen.Components` permanece coleção plana ordenada **do fundo para a frente**, única fonte persistida da ordem de camadas. `screen.PlaceComponent(componentID, referenceID string, inFront bool) error` move somente a entrada escolhida: imediatamente depois/antes da referência ou fim/início com referência vazia. A ordem relativa de todas as demais entradas é estável; posição já satisfeita não muda a coleção. Validar tela/origem/referência e recusar autorreferência antes de mutar; erro preserva estado integral.

Nenhum campo `inFront`, `referenceID`, z-index extra, modo ou relação permanente é acrescentado ao documento. Schema 3/catalogVersion 2 não mudam. Os parâmetros representam um comando pontual; a referência não passa a restringir edições posteriores. Salvar/reabrir usa a ordem canônica existente. Modo neutro e referência do formulário pertencem ao rascunho do adaptador e não ao modelo.

**Precisão sobre a hierarquia descrita acima:** ParentID não torna a reordenação uma operação de subárvore. Somente o componente indicado muda de índice; filhos conservam suas próprias posições na lista e não acompanham a camada do pai. Não alterar bounds, estilo, conteúdo, IDs/FileIDs, ParentID ou vínculos. O histórico da aplicação captura a ordem e as outras propriedades da mesma transação atômica; falha de qualquer campo não publica edição parcial.

Regressões de domínio e round-trip/cena estão implementadas. O executor principal informou gofmt dos sete arquivos e vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e build/pacote Windows com integridade PASS. T014 local e [001/T127](../001-gerador-jornada/tasks.md#phase-21-sobreposição-de-componentes--2026-10-20) permanecem abertas: compilação não comprova invariantes em execução; testes Go não executados por política Windows, sem contorno. Evidências e identidade na [seção 9 do quickstart](../001-gerador-jornada/quickstart.md#9-sobreposição-de-componentes--2026-10-20); UI nativa, DPI, Office e aprovação manual não realizados. Resultados e hashes anteriores são históricos, não validação desta operação.

## Uso universal de TextAlign — 2026-10-20

FR-017 reutiliza `Component.Style.TextAlign`, sem campo/enum novo, mudança de schema 3/catalogVersion 2 ou domínio adicional. Valores explícitos `left/center/right` afetam somente texto nas 16 famílias; o vazio selecionado por **Padrão da skin** deixa a resolução à skin/default existente. Isso é distinto do vazio de arquivos antigos, normalizado para `left` pela migração já existente.

O rascunho só publica após Salvar alterações e participa do histórico da aplicação; seleção/reabertura carrega o valor confirmado. Bounds, ParentID, hierarquia, IDs, ordem, conteúdo, vínculos e recursos permanecem iguais. Cena e adaptadores conservam fundos/hit areas e alinham cada texto em sua região interna (célula, área de relatório, legenda de tecla), inclusive opções e Sem seleção no popup do combo. Popup não é novo estado persistido. Após produção e quatro novos arquivos de testes, o executor principal informou gofmt/vet PASS e 15 suítes Go compiladas em 20/10, sem executar os testes pelo bloqueio prévio de política Windows, sem contorno. A compilação não comprova as invariantes em execução; T013 e manuais permanecem abertos. Evidências técnicas, quatro Python de empacotamento PASS e build Windows PASS na [seção 8 do quickstart](../001-gerador-jornada/quickstart.md#8-alinhamento-textual-universal--2026-10-20).

<a id="modelo-codigos-barras-2026-10-23"></a>
## Instâncias planejadas — códigos de barras (2026-10-23)

FR-019 reutiliza `Component.ID/FileID/Kind/Bounds/Style/Properties` e ParentID existentes. Não acrescenta entidade de boleto a Screen; tipos das propriedades pertencem ao [Catálogo](../003-catalogo/data-model.md#modelo-codigos-barras-2026-10-23). Subcampos da entrada têm índices derivados, não UUIDs de Component. Blocos/linhas processados têm UUIDs internos únicos, sem FileIDs e sem entradas próprias em Components.

Snapshot/Copy mantém identidade e separa dados mutáveis; duplicação gera novas identidades de composto/blocos/linhas e remapeia referências locais, não usa títulos como chave. Reordenar conserva IDs; remover não reutiliza ID para outro item. Unicidade global inclui os itens internos novos. `targetComponentId` de teclado pode apontar à entrada composta habilitada da mesma tela; índice do campo ativo é transitório na aplicação, não um vínculo persistido de oito componentes.

Trocar `labelKind` ou `value` não altera geometria, camada ou outros valores. A aplicação valida e publica o conjunto em uma única transação com rollback/histórico existentes. Canvas/inspector e runtime não adicionam foco/popup/modal ao modelo de Tela.

<a id="modelo-menu-aplicacao-2026-10-23"></a>
## Instância planejada — Menu aplicação

FR-020 reutiliza Component com Kind/Properties canônicos em [003/data-model.md](../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23). Uma instância contém árvore, não novos componentes por item. Menu/submenu e MenuItem possuem UUIDs internos globalmente únicos, sem FileIDs. Estado de expansão/seleção do inspector e caminho runtime não pertencem a Screen.

Snapshot mantém IDs, copia profundamente todos os slices/ponteiros. Duplicação cria novo Component.ID/FileID e novos IDs em Root, cada submenu e cada item; mapa de remapeio atualiza InitialMenuID, nunca label/name. TargetScreenID referencia tela real e não é remapeado pelo mapa de nós: duplicar componente conserva destinos; duplicar tela remapeia apenas destinos autorreferentes à tela duplicada para a nova tela, preservando destinos externos, conforme 005. Reordenação e edição de texto conservam identidade. Remover inicial sem substituição válida impede confirmação; remover ramo elimina descendentes sem deixar referência inicial pendurada.
