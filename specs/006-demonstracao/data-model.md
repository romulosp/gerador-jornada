# Data Model: Domínio Demonstração

## Situação corrente — revisão 1.0.0.1

Scenario/Step são interpretados em eventos/ocorrências; application/playback.go e desktop/player.go consomem essa interpretação. Inputs runtime de prévia não são persistidos nos roteiros. typeText de compostos/menu e os tipos novos continuam sem capacidade disponível.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/scenario`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Scenario contém id, name, startScreenId e steps. Step contém id, kind, screenId, componentId, point, text e durationMs, com nulos e campos obrigatórios por kind conforme contrato. Ocorrência deriva de cada visita e acumula tempo; cursor inicia em (0,0).

## Invariantes

IDs de steps únicos no projeto; duração inteira 1..600000 ms; etapas são sequenciais sem sobreposição. moveCursor usa ponto ou componente, nunca ambos. typeText exige textField e substitui valor runtime. navigate exige button com binding válido; todas as etapas usam screenId corrente.

## Limites e relacionamento com outros domínios

Recebe snapshot resolvido de telas/componentes e resolução Journey, preparado pela aplicação. Interpretador devolve eventos puros. A arquitetura alvo mantém relógio/pausa na aplicação. Atualmente FrameAt deriva quadros puros, enquanto ticker/pausa/reinício ficam em desktop/player.go; a separação e seus testes controlados permanecem em 001-gerador-jornada/T045 e no backlog local. PPT consome eventos interpretados sem relógio de UI.

Não grava uso do mouse, não cria macros, não contém lógica condicional e não emite XML PPT ou eventos Fyne.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: intervalo temporal, campos por kind, alvo removido, tela errada, A-B-A, auto-navegação, typeText vazio, pausa e reinício. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-017 a FR-021 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.


## Correspondência com o código atual

Interpretador, ocorrências, grafemas, cursor/clique/digitação e transições XML implementados. Pausa/reinício existem na UI, mas o relógio ainda não é abstraído na aplicação; reordenação disponível para etapas, não para a lista de roteiros. T004/T007 e 001-gerador-jornada/T044–T047 conservam essas pendências; execução PowerPoint ainda não homologada.

## Alvos e estados novos — contrato alvo 2026-10-19

Scenario/Step mantêm identidades e cinco kinds. O índice de capacidades amplia textField para textArea/passwordField onde suportado, sem converter Kind desconhecido; password é exemplo de conteúdo, nunca segredo de serviço. Click continua exclusivamente visual, sem mutação de seleção/toggle/modal. Navigate continua button por binding, inclusive twoLabel. Os campos persistidos não recebem callbacks ou código.

Ocorrência/quadro incorpora valores/seleções/estados iniciais, perfil visual imutável e recorte necessário da cena; cursor/foco/scroll/modal transitórios são derivados, não dados de Project. Revisita/restart não herdam interação anterior. Unicode/tempo e opções de transição permanecem iguais. Fonte/skin é resolvida por apresentação/consumidores, não por parser XML em scenario. T009–T010 verificam capacidades, quadro determinístico e diagnóstico por etapa; 001/T116 traduz para PPT e T121 homologa.

O perfil por documento descrito nos registros de 2026-10-19 é histórico/plano de expansão; não existe na revisão 1.0.0.1. A implementação atual usa Style e biblioteca incorporada global, sem captura/selamento de revisões privadas. Evidências antigas permanecem no quickstart e não comprovam o estado deste build.

A expansão de quadros resolve estados/capacidades sem serializar interação livre no Scenario. Relógio controlado, novos tipos e homologação do Office permanecem pendentes.

<a id="modelo-codigos-barras-2026-10-23"></a>
## Quadros planejados — códigos de barras (2026-10-23)

FR-011 não altera Scenario/Step, seus cinco kinds ou tempos. Snapshot de quadro lê configuração dos compostos de Tela/Catálogo: entrada sem valores iniciais e processado com blocos/linhas/valores confirmados. Campo ativo, captura, validação e modal de uma prévia livre não integram o quadro nem são gravados em Scenario. `typeText` não recebe capacidade para esses tipos nesta entrega; tentativa é diagnóstico por ID da etapa. Conteúdo processado continua string literal e rótulo derivado do enum, sem consulta de dados ou evento oculto.

<a id="modelo-menu-aplicacao-2026-10-23"></a>
## Quadro planejado — Menu aplicação

FR-012 mantém Scenario/Step e cinco kinds existentes, sem itemID persistido por suposição. Quadro resolve Root/InitialMenuID do snapshot de Tela e apenas a vista inicial correspondente com título/itens ordenados/enabled; estado runtime de outra sessão é irrelevante. Caminho de ancestrais é derivado, não propriedade nova no projeto. Menu vazio pode formar quadro estático se restante do projeto for válido.

Não há capacidade typeText para applicationMenu/barcodeInput/barcodeProcessed. Click é destaque visual, não seleção local ou dispatch. Extensão futura de click por item permanece decisão coordenadora, com contrato e compatibilidade antes do código. Fonte da configuração em [003/data-model.md](../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23); apresentação editável em 001, não uma cópia de MenuItem em Step.
