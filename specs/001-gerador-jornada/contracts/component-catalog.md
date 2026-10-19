# Contract: catálogo de componentes

**Version**: 1 | **Requirements**: FR-002 a FR-008, FR-015, FR-016, FR-022.

## Modelo comum

Kind é identificador estável em inglês; nome exibido e mensagens são em português. Catálogo é embutido no editor. Um projeto referencia `catalogVersion: 1`; não instala componentes externos.

Todo componente tem UUID interno, nome lógico único na tela, fileId de arquivo, bounds inteiros e propriedades específicas. FileID é atributo comum da instância em Tela, não uma variante de Properties; segue o [contrato de projeto v2](project-format.md#nome-base-de-arquivos-e-migração-v1--v2), com até 40 caracteres ASCII e sem Pnl obrigatório. Defaults: x/y do ponto de drop, dimensões da tabela, texto em Arial 12 unidades lógicas, cor `#4864A0`, preenchimento branco e borda cinza quando a família tiver caixa.

Arraste conserva o deslocamento entre ponteiro e origem. Novo drop ou drag que exceda a tela é ajustado para caber; se a tela for menor que o default, reduzir bounds até caber sem valores zero. Alteração de tamanho da tela pode deixar elementos inválidos destacados, sem reposicioná-los.

A ordem do array de componentes define camadas. Zoom de 25% a 400% afeta somente renderização. Coordenadas Java recebem bounds diretamente; PPT usa a conversão do seu contrato.

## Famílias e propriedades

| Kind | Família / referência EFL | Default width × height | Properties persistidas |
| --- | --- | --- | --- |
| button | Botão / EFLButton | 110 × 60 | text string, variant enum confirm/cancel/yes/no |
| label | Rótulo / EFLLabel | 150 × 20 | text string |
| title | Título / EFLLabel com skin de título | 300 × 28 | text string |
| textField | Campo / EFLTextField | 150 × 25 | initialValue string, format enum text/number/cellular/cep/cpf/date/ddd/nis |
| radio | Seleção / EFLButtonGroup a confirmar | 180 × 60 | options array de Option, orientation enum vertical/horizontal, selectedOptionId UUID ou null |
| report | Relatório/lista / EFLReportEvent | 300 × 345 | title string, rows array de Row |
| numericKeyboard | Teclado numérico / EFLVirtualKeyboardNumeric | 200 × 280 | targetComponentId UUID ou null |

Option e Row contêm `id` UUID e `text` string. Arrays preservam ordem; seleção exige pelo menos uma opção; relatório pode estar vazio. IDs de opção/linha devem ser únicos dentro do componente. Não usar strings separadas por vírgula: o próprio texto pode conter vírgulas.

Todos os campos da linha são obrigatórios; texto pode ser vazio, mas nome lógico do componente não. Variantes default: confirm, text, vertical; options começa com duas opções com IDs próprios, report começa sem linhas; seleção e alvo do teclado começam nulos.

O teclado vincula apenas textField da mesma tela, aceita dígitos, apagar último caractere e limpar. Sem alvo, é válido como desenho, mas prévia emite warning e não altera campo algum. Fonte ausente, opção selecionada inexistente e alvo de tipo errado são diagnósticos.

## Aparência e capacidade por alvo

Style opcional: fontFamily string não vazia, fontSize positivo em unidades lógicas, textColor/fillColor/borderColor em `#RRGGBB`, backgroundAssetId UUID ou null e textAlign em `left`, `center` ou `right`. Ausência de textAlign mantém esquerda nos arquivos anteriores. Alinhamento é do texto dentro dos bounds, independente de mover/alinha-los na tela. A faixa pronta usa direita. Ausência das demais propriedades significa default do catálogo; nunca inferir propriedades do widget Fyne.

Na v1 usar texto horizontal de uma linha para label/title/button e linhas explícitas em relatório. Não implementar rich text, HTML, rotação ou layout responsivo. Conteúdo que não cabe recebe warning de overflow; o usuário deve ajustar bounds antes da homologação visual.

Desktop e PPT suportam estilo definido no catálogo e recursos decorativos. Java só aceita personalização que o perfil EFL consiga representar. Antes de exportar, validar cada propriedade alterada; `JAVA_PROPERTY_UNSUPPORTED` impede essa saída, sem ignorar a mudança nem trocar para Swing puro.

A apresentação agrupa primitivas por componente, mas mantém texto e formas internos editáveis. Radio, report e numericKeyboard expandem em várias primitivas; não viram screenshots.

## Interação e limites

Prévia/Java permitem edição e seleção reais. Formatos de documento EFL são nomes de intenção observados no legado; máscara e validação exatas dependem do perfil homologado. Na prévia v1, text aceita Unicode e number aceita dígitos; demais formatos apresentam texto/sugestão de formato, sem afirmar validação de domínio fiscal ou de data.

PPT permite clicar em botões de navegação e demonstrar preenchimento; os campos não recebem teclado nem calculam validações no modo apresentação. O diálogo de exportação deve explicar essa diferença.

Report é lista visual nesta versão; seleção de linha não dispara regra de negócio. NumericKeyboard é representação editável no PPT; preenchimento da demonstração é uma etapa typeText, não um teclado executável no PowerPoint.

Ações de botões ficam em Journey, não como funções em properties. Kind é dado; não resolve classe por reflexão nem executa template do projeto.

## Comandos de edição

Selecionar, mover, redimensionar, duplicar, remover, centralizar horizontal/verticalmente, alinhar bordas à tela e mover uma camada para frente/trás. Duplicação desloca 10 unidades se couber e sempre cria novo ID/nome lógico sem colisão; se não couber, conserva posição.

Atalhos: Delete remove; setas movem 1 unidade; Shift+seta move 10; Ctrl+D duplica; Ctrl+Z desfaz; Ctrl+Y refaz; Ctrl+S salva. Não interceptar esses atalhos de edição de desenho enquanto um campo de texto estiver recebendo entrada. Escape cancela manipulação em andamento.

Um arraste gera uma operação; alterações confirmadas no inspector também. Agrupamento de seleção múltipla e componentes aninhados não fazem parte da v1.

## Validação

Testar cada família desde defaults até serialização, layout e saídas. EFLRadio do legado não prova construtor de grupo: a inconsistência foi registrada em research.md e deve ser corrigida após inspeção da biblioteca real.
