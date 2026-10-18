# Contract: navegação e roteiro

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

A prévia usa relógio injetável para testes. O PPT compila os mesmos eventos em efeitos e tempos; não grava vídeo da prévia. Exportação pode selecionar exploração apenas, um roteiro ou todos; opção padrão é exploração mais todos os roteiros não vazios.

Se houver roteiro selecionado inválido, bloquear esse job e indicar como exportar apenas a exploração ou corrigir o roteiro. Não omitir silenciosamente um roteiro selecionado.

## Diagnósticos

`JOURNEY_START_MISSING`, `NAVIGATION_TARGET_MISSING`, `NAVIGATION_SOURCE_INVALID`, `SCENARIO_EMPTY`, `STEP_TARGET_MISSING`, `STEP_SCREEN_MISMATCH`, `STEP_DURATION_INVALID`, `STEP_KIND_INVALID` e `STEP_VALUE_INVALID`. Localizar por IDs e nomes visíveis quando disponíveis; não substituir tela ou alvo automaticamente.

## Exemplos de validação

No cenário A → B → A: roteiro inicia em A, executa moveCursor, click, navigate, wait em B, click e navigate de volta a A. Gerar três ocorrências e confirmar reset de valores na revisita. Testar auto-navegação, alteração posterior do binding, remoção de alvo, pausa, Unicode e duração inválida.
