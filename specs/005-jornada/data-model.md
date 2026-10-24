# Data Model: Domínio Jornada

## Situação corrente — revisão 1.0.0.1

Journey mantém StartScreenID e Actions none/navigate para botões. Folhas de menu guardam TargetScreenID em suas propriedades e o runtime trata esse destino. Não há NavigationGraph/Graph ou índice completo de integridade/alcançabilidade de menus na árvore atual.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/journey`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Journey contém startScreenId e actions. ButtonAction contém sourceScreenId, sourceComponentId, kind e targetScreenId. IDs são referências, nunca índices ou nomes. O grafo deriva da lista de ações; não persistir uma segunda cópia divergente.

## Invariantes

Uma ação por botão; fonte deve ser button da tela informada. none implica destino nulo; navigate implica destino existente. Não exigir grafo acíclico. Sem inicial, salvar rascunho é permitido; executar não.

## Limites e relacionamento com outros domínios

Recebe índice imutável de referências permitido pela aplicação, sem importar Screen ou adaptadores. Scenario usa resolução de rota para interpretar navigate. Desktop e Java executam a ação; PPT traduz em links.

Não armazena valores de campo, histórico de navegação, widgets, slides ou listeners Java.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: none implícito, ciclo, auto-navegação, fonte não botão, destino apagado, inicial ausente e binding duplicado. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/navigation-scenario.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-011 a FR-013, FR-014, FR-023 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.


## Correspondência com o código atual

Binding none/navigate, resolução por ID e validação de origem/destino estão implementados. O warning de tela inalcançável previsto em FR-004 ainda não é emitido; permanece em T006. Diagnóstico clicável/localização completa também depende de 001-gerador-jornada/T030.

## Capacidades do catálogo 2 — alvo 2026-10-19

Journey/Action não acrescentam callbacks, skins, foco/seleção/scroll/modal ou dados de serviço. O índice de aplicação resolve sourceScreenId/sourceComponentId por ID e capacidade tipada, preservando a restrição navigate de button (inclusive twoLabel). Eventos de seleção/teclado de novas famílias são locais, não bindings de navegação implícitos. Não persistir grafo ou estado de visita duplicado.

Revisita restaura valores/seleções/estados iniciais tipados no snapshot de prévia, sem mudar o documento. Capacidades de texto/foco/grupo/contêiner e referências locais são validadas na fronteira de aplicação; aparência nova/migração não altera bindings/startScreenId. T009–T010 locais testam resolução/erros/revisita com esses índices; 001/T114 implementa sessão runtime, sem duplicar trabalho.

O perfil por documento descrito nos registros de 2026-10-19 é histórico/plano de expansão; não existe na revisão 1.0.0.1. A implementação atual usa Style e biblioteca incorporada global, sem captura/selamento de revisões privadas. Evidências antigas permanecem no quickstart e não comprovam o estado deste build.

A expansão não deve transferir perfis para Journey nem remapear rotas ao trocar aparência. O índice global de destinos de menu permanece pendente.

<a id="modelo-codigos-barras-2026-10-23"></a>
## Jornada inalterada — códigos de barras (2026-10-23)

FR-008 não introduz campos em Journey/Action, nova ação ou entidades de boleto. Os dois compostos não têm capacidade de origem `navigate`; índice da aplicação os reconhece pelos tipos sem adaptá-los a botões. Dados configuráveis processados pertencem às propriedades do Catálogo/Tela; a sessão de entrada pertence à aplicação. `startScreenId`, UUIDs de origem/destino e bindings existentes não mudam por validação ou edição de rótulo. Revisita restaura defaults sem serializar conteúdo capturado.

<a id="modelo-menu-aplicacao-2026-10-23"></a>
## Referências planejadas — folhas de menu

FR-009 mantém Journey/Action de botões; não trata ApplicationMenu inteiro como button nem cria binding paralelo por folha. TargetScreenID de MenuItem é o único destino persistido autorado. Uma projeção transitória fornecida pela aplicação identifica `(sourceScreenID, componentID, itemID, targetScreenID, enabled efetivo)` para validação e diagnósticos; índice conserva também referências em ramos desabilitados para detectar destinos removidos. O modelo da árvore continua em Catálogo/Tela.

Grafo considera arestas explícitas de folhas habilitadas para alcançabilidade e valida todos os destinos configurados para integridade. IDs de item localizam origem, não substituem Screen.ID do destino. Remoção de tela lista impactos; retirar destino exige comando confirmado e reversível, sem redirect silencioso. Caminho/root/submenu/Voltar/Cancelar/Sair e seleção simulada são sessão da aplicação, não entidades de Journey. Não há ações financeiras, de autenticação ou condições por código de movimento.
