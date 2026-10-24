# Contract: navegação e roteiro

## Implementação conferida — 1.0.0.1

Bindings de botões são `none`/`navigate`; runtime de menu lê o destino autorado de uma folha. `Document.Validate` ainda não possui grafo completo de menus ou warning de inalcançáveis. Esses diagnósticos são requisitos pendentes, não comportamento comprovado pelo fato de navegar na prévia.

Roteiros usam os tipos já existentes de cenário; clique continua visual e navegação é explícita. `typeText` não opera os compostos 031/032/033. Campos, escolhas e caminho de menu na prévia livre não alteram os dados confirmados ou a revisão. Ausência em strings opcionais é gravada como string vazia; texto/point opcionais de etapa usam os ponteiros tipados do modelo.

Código atual: [wire do projeto](../../../apps/gerador-jornada/internal/adapter/projectzip/wire.go), [cena](../../../apps/gerador-jornada/internal/presentation/scene.go) e [exportador PPT](../../../apps/gerador-jornada/internal/adapter/pptx/export.go). Evidência executada na [seção 22 do quickstart](../quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). O restante do contrato conserva regras normativas e evolução planejada; divergências apontadas aqui não foram apagadas dos requisitos.


**Version**: 1 | **Requirements**: FR-011 a FR-013, FR-017 a FR-021.

## Navegação da jornada

Journey contém `startScreenId` (UUID ou null no rascunho) e `actions`. Cada ação tem `sourceScreenId`, `sourceComponentId`, `kind` e `targetScreenId`.

- `kind: none`: targetScreenId nulo; botão não navega.
- `kind: navigate`: targetScreenId deve identificar tela existente.
- No máximo um binding por botão; ausência de binding equivale a none.
- A fonte precisa ser button pertencente à sourceScreenId. Não usar índice da tela como destino.
- Ciclos e auto-navegação são permitidos; telas inalcançáveis produzem warning, não são omitidas da exportação.
- Antes de prévia ou exportação, exigir inicial válida e ações válidas. Após entrar numa tela, inicializar valores de campos e seleção a partir do projeto, sem persistir estado de uso.

Prévia e Java executam navigate no clique normal. Exploração PPT usa link para slide da tela de destino. Não há condicionais, código de usuário, navegação baseada em valor de campo ou chamadas externas na v1.

## Roteiro persistido

Scenario contém `id`, `name`, `startScreenId` e `steps`. Step contém `id`, `kind`, `screenId`, `componentId`, `point`, `text` e `durationMs`. Campos não aplicáveis devem ser nulos; campos obrigatórios por kind não podem faltar.

Duração inteira de 1 a 600000 ms por etapa. As etapas são sequenciais, sem sobreposição e sem saltos para índices anteriores. Roteiro vazio pode ser salvo; não é exportável como demonstração. Todos os IDs de etapas são únicos no projeto.

| Kind | Campos não nulos específicos | Semântica |
| --- | --- | --- |
| moveCursor | point `{x,y}` ou componentId, exatamente um | Mover linearmente até ponto da tela ou centro do componente durante durationMs. |
| click | componentId de button ou textField | Posicionar cursor no centro do alvo e mostrar destaque de clique durante durationMs. Não dispara navegação nem altera dados. |
| typeText | componentId de textField, text | Substituir o valor runtime do campo por text, revelando progressivamente durante durationMs. |
| wait | nenhum | Manter estado atual durante durationMs. |
| navigate | componentId de button com ação navigate | Manter estado por durationMs e, ao fim, seguir o destino do binding da jornada. |

Toda etapa tem screenId igual à tela corrente. A tela corrente começa em Scenario.startScreenId e só muda no fim de navigate. O destino não é duplicado no step; mudar o binding muda o destino e pode invalidar etapas seguintes, com diagnóstico.

**Clique e navegação são ações demonstrativas distintas.** A timeline deve explicitar essa regra e oferecer os dois comandos. Separá-los evita efeitos ocultos e permite pausas entre clique e troca de tela. Na prévia livre o botão continua navegando normalmente.

typeText aceita texto vazio para limpar. Campo number admite apenas dígitos; outras máscaras seguem os limites de simulação definidos no catálogo. A revelação usa caracteres visíveis (clusters Unicode), não bytes; o valor final deve preservar acentos e pares Unicode.

## Interpretação e estados

O interpretador recebe snapshot e Scenario e produz eventos com tempo acumulado, estado inicial/final e ocorrências de telas. Não depende de Fyne nem de PPT XML.

- Ao iniciar: valores iniciais da tela e cursor em (0,0).
- Ao navegar: encerrar ocorrência, entrar na tela de destino com valores iniciais e cursor em (0,0). Revisitas são novas ocorrências.
- typeText define valor apenas no estado runtime; wait e moveCursor não mudam campos.
- Pause congela relógio; continue retoma o tempo restante; restart descarta estado runtime e inicia do começo.
- A última ocorrência permanece visível ao terminar. Reprodução não cria histórico de edição.

Cada ocorrência contém lista de etapas contíguas da mesma visita, mesmo quando sourceScreenId e targetScreenId são iguais. Um roteiro que navega para a própria tela precisa de nova ocorrência, não da continuação do slide anterior.

## Prévia e exportação

A arquitetura alvo exige relógio injetável para testes. Na implementação atual, FrameAt é puro, mas o ticker/pausa/reinício ainda estão em desktop/player.go; a extração e os testes controlados permanecem em 001-gerador-jornada/T045 e 006-demonstracao/T004/T007. O PPT compila os mesmos eventos em efeitos e tempos; não grava vídeo da prévia. Exportação pode selecionar exploração apenas, um roteiro ou todos; opção padrão é exploração mais todos os roteiros não vazios.

Se houver roteiro selecionado inválido, bloquear esse job e indicar como exportar apenas a exploração ou corrigir o roteiro. Não omitir silenciosamente um roteiro selecionado.

## Diagnósticos

`JOURNEY_START_MISSING`, `NAVIGATION_TARGET_MISSING`, `NAVIGATION_SOURCE_INVALID`, `SCENARIO_EMPTY`, `STEP_TARGET_MISSING`, `STEP_SCREEN_MISMATCH`, `STEP_DURATION_INVALID`, `STEP_KIND_INVALID` e `STEP_VALUE_INVALID`. Localizar por IDs e nomes visíveis quando disponíveis; não substituir tela ou alvo automaticamente.

## Exemplos de validação

No cenário A → B → A: roteiro inicia em A, executa moveCursor, click, navigate, wait em B, click e navigate de volta a A. Gerar três ocorrências e confirmar reset de valores na revisita. Testar auto-navegação, alteração posterior do binding, remoção de alvo, pausa, Unicode e duração inválida.

## Expansão pós-v1 — capacidades e estado local (2026-10-19)

O catálogo 2 amplia alvos sem transformar ações de domínio em serviços. Button de duas legendas continua button com none/navigate e UUID original. Toggle, seleção de Combo/UF/Table/Report e eventos de teclado não recebem navigate por similaridade com botão: permanecem comandos/eventos locais tipados. Integrações externas são somente contratos, nunca chamadas reais. Java segue pausado.

O índice preparado pela aplicação deve informar capacidades tipadas (recebe texto, foco, seleção, scroll, contêiner/modal, habilitado) em vez de adaptar todo controle para textField. Teclados validam alvo compatível na mesma tela; alvo ausente é diagnosticado e não resolvido em outra tela. TAB percorre apenas alvos habilitados/compatíveis na ordem contratada da composição, sem alterar vínculo persistido; cursor/modificadores são runtime. Máscaras respeitam apenas regras locais comprovadas e declaradas no catálogo.

`click` do roteiro continua **somente visual**; não muda toggle/seleção/modal nem navega implicitamente. `navigate` continua exclusivo de button com binding. `typeText` pode alcançar textArea/passwordField quando capacidade tipada permitir; conteúdo de password é exclusivamente ilustrativo e sua representação usa mascaramento aprovado, sem revelar texto em logs. Cada nova combinação precisa validação/teste; não aceitar Kind desconhecido nem converter alvo silenciosamente. Nenhum novo Kind de Step é acrescentado por esta revisão; interações locais sem etapa tipada podem ser exercitadas na prévia livre, não anunciadas como roteiro já suportado.

Ao entrar/reentrar numa tela, restaurar valores/seleções/estados iniciais e zerar foco/cursor/scroll/modal transitórios conforme defaults contratados. Ocorrência/quadro contém snapshot de estado local necessário à cena; nunca grava de volta no documento. Mesma configuração e tempo produzem o mesmo quadro em prévia/PPT, sem relógio da máquina, periférico ou consulta externa. Persistência desses novos tipos pertence ao schema 3; T009–T010 dos domínios 005/006 validam capacidades/revisitas e 001/T114–T116 integram os consumidores.

<a id="limites-codigos-barras-2026-10-23"></a>
## Limites planejados — códigos de barras (2026-10-23)

Os compostos [031](../../031-codigo-barras-efl/spec.md) e [032](../../032-codigo-barras-processado-efl/spec.md) não são fontes de `navigate` e não alteram Journey/Action/Step ou seus kinds. Apenas a prévia livre opera entrada por bloco, captura e modal. O estado aceito da entrada não é rota, precondição persistida de botão ou processamento de dados de outra tela. Prosseguir significa concluir a entrada localmente; botões continuam com ações explícitas existentes, sem condição financeira nova no grafo.

Nesta entrega, `typeText` de roteiro não é suportado nos dois compostos; diagnosticar capacidade e ID da etapa antes de demonstrar/exportar o roteiro, sem adaptar o composto a textField ou omitir a etapa. A representação estática é permitida: entrada vazia, ficha com valores confirmados, textos/formas editáveis no PPT. `click` permanece visual e não abre modal, escolhe rótulo ou processa código implicitamente.

Revisita/restart restaura entrada vazia e configuração processada; captura da prévia livre nunca altera Scenario. Digitação animada de boleto, atualização automática de ficha ou navegação condicional requerem especificação de extensão própria antes de implementar. Regras locais em 003, instâncias em 004, verificações de rotas em [005/T011](../../005-jornada/tasks.md#tarefas-codigos-barras-2026-10-23), quadros/capacidades em [006/T011](../../006-demonstracao/tasks.md#tarefas-codigos-barras-2026-10-23) e adaptadores em 001/T135–T139; não duplicar implementações.

<a id="navegacao-menu-aplicacao-2026-10-23"></a>
## Extensão planejada — Menu aplicação

Esta seção amplia a regra histórica “somente botão” **apenas para a folha de menu com destino explicitamente autorado na prévia livre**; não transforma applicationMenu inteiro ou 031/032 em button/navigate. O tipo/árvore é Catálogo, composição Tela, persistência Projeto, autoridade de destino Jornada e execução de sessão aplicação.

1. Abrir submenu muda apenas vista/caminho local da mesma instância. Caminho inicial é derivado de Root até InitialMenuID. Voltar remove um nível (raiz é limite), Cancelar volta à raiz e Sair encerra somente a prévia, sem logoff, rede ou encerramento do editor. Revisita reinicia no inicial autorado.
2. Ativação exige instância e caminho habilitados. Item com Submenu não pode ter TargetScreenID. Folha com destino não vazio pode pedir navegação para a tela existente; destino vazio informa seleção simulada. Label/Category/Variant/Shortcut/MovementCode/MovementOption/MovementName não são comandos ou lookup de rota. Não há movimento bancário, conta, autenticação, QR, timer ou periférico.
3. TargetScreenID na árvore é a referência persistida única; não criar binding concorrente em Journey.actions. A aplicação projeta screen/component/item/target e enabled efetivo para validação pura de Jornada. Validar destinos de todos os ramos, inclusive desabilitados; listar impactos ao remover tela. Alcançabilidade considera somente folhas efetivamente habilitadas; ciclos/autorrotas explícitos continuam permitidos.
4. Duplicar componente conserva destinos existentes; duplicar tela remapeia só destinos autorreferentes à tela duplicada e conserva externos. Identidades de nós/inicial são remapeadas por Tela. Remoção de destino exige comando explícito e diagnóstico, não redirect por nome ou substituição silenciosa.
5. **Prévia Tela Atual** permite submenus/Voltar/Cancelar/seleção simulada/Sair, mas nunca segue destino entre telas. Referências externas continuam persistidas; não bloquear exploração local por tela não incluída no snapshot individual. Prévia completa/exportação mantêm validação global.

### Demonstração e saída estática

Quadro/canvas/PPT estático resolvem InitialMenuID do snapshot confirmado, não o caminho da prévia livre. Estado fixo não significa menu interativo dentro do PPT. Textos/formas são editáveis pela cena compartilhada; nenhum código de movimento é executado. Não há typeText para applicationMenu, barcodeInput ou barcodeProcessed.

Step e seus cinco kinds permanecem; click continua visual sem abrir submenu, selecionar folha ou navegar. Navigate de roteiro continua button com binding, até contrato de capacidade explicitamente aprovado. A coordenação principal pode decidir extensão de click por item metadata, **mas não a habilitar por interpretação implícita**: antes disso precisa sincronizar Step/wire, capacidades, resolução componentID/itemID/hit regions, quadros determinísticos, diagnóstico e compatibilidade. Sem tal decisão, não persistir itemID novo em Step, não criar novo Kind e recusar interações não contratadas por ID da etapa. Não omitir roteiro inválido selecionado.

Donos: [005/T012](../../005-jornada/tasks.md#tarefas-menu-aplicacao-2026-10-23) destinos; [006/T012](../../006-demonstracao/tasks.md#tarefas-menu-aplicacao-2026-10-23) quadros; [001/T140–T144](../tasks.md#tarefas-menu-aplicacao-2026-10-23) sessão/adaptadores/validação/manual. Esta documentação não executou testes nem concede aceite dos gates históricos.
