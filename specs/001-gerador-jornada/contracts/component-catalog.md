# Contract: catálogo de componentes

## Implementação conferida — 1.0.0.1

O catálogo gravado é **3**, com dez Kinds efetivos: `button`, `label`, `title`, `textField`, `radio`, `report`, `numericKeyboard`, `barcodeInput`, `barcodeProcessed`, `applicationMenu`. Definition tem Kind/Name/Width/Height; não há um metamodelo declarativo completo de capacidades ou todas as 16 famílias da expansão. Ver a [spec do Catálogo](../../003-catalogo/spec.md#estado-verificado-1-0-0-1).

Default numérico vigente: **200 × 350, 14 teclas**; bounds antigos não são migrados pelo renderer. Skins vazias resolvem os defaults do legado e overrides explícitos prevalecem. As 205 chaves de skin, 224 recursos e seis fontes físicas pertencem à biblioteca incorporada; oito definições indisponíveis não entram na oferta. Requisitos de aparência tipada, novos tipos, variantes completas e perfil portátil continuam pendentes. As seções de catálogo 1/2 abaixo preservam histórico e contrato de expansão.

Código atual: [wire do projeto](../../../apps/gerador-jornada/internal/adapter/projectzip/wire.go), [cena](../../../apps/gerador-jornada/internal/presentation/scene.go) e [exportador PPT](../../../apps/gerador-jornada/internal/adapter/pptx/export.go). Evidência executada na [seção 22 do quickstart](../quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). O restante do contrato conserva regras normativas e evolução planejada; divergências apontadas aqui não foram apagadas dos requisitos.


**Version**: catálogo gravado 3; base 1 e plano histórico de expansão 2 preservados abaixo | **Requirements**: FR-002 a FR-008, FR-015, FR-016, FR-022, FR-049–FR-054.

As seções v1 abaixo são histórico da primeira versão. A expansão ao final define o alvo vigente de implementação, sem afirmar que já foi entregue. Java continua pausado.

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
| numericKeyboard | Teclado numérico / EFLVirtualKeyboardNumeric | 200 × 350 (280 é histórico) | targetComponentId UUID ou null |

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

## Catálogo 2 — expansão autorizada 2026-10-19

Preservar discriminadores das sete famílias. Skins/nomes XML são variantes/presets, não novos Kinds por entrada. Propriedades são structs e união tipada, sem JSON livre no domínio. Geometria nova aplica-se somente a novas inserções: button continua 110 × 60 (160 × 44 era divergência no DESIGN 007); numericKeyboard passa de 200 × 280 a 200 × 350 por referência Java. Projetos antigos mantêm bounds; não redimensionar na migração.

| Kind / preset planejado | Contrato semântico mínimo, validado por família | Referências / dono |
| --- | --- | --- |
| button (existente) | Preservar variantes históricas; acrescentar variante twoLabel e segunda legenda tipada, ambas editáveis; ação de Journey continua separada | EFLButton/TwoLabelButton; 003/T011 |
| label/title (existentes) | Texto Unicode, todas skins aplicáveis, alinhamento/face/alpha; ValorSurpresaSuperSete é preset label+toggle opcional, sem cálculo de negócio | EFLLabel; 003/T011 e 004/T011 |
| toggleButton (novo) | text, initialSelected, enabled e groupId opcional; grupo tipado com seleção simples/múltipla e validação de membros na mesma tela | EFLToggleButton/Group; 003/T011 |
| radio (existente) | Opções/IDs/orientação/selectedOptionId permanecem; seleção simples não muda silenciosamente para múltipla | 011; 003/T011 |
| textField (existente) | Preservar formatos; ampliar documentos/máscaras comprovados por enums/regras locais e diagnósticos, sem anunciar validação fiscal só pelo nome | EFLTextField/Document; 003/T012 |
| textArea (novo) | initialValue, formato suportado, multiline, quebra/recorte/scroll tipados | EFLTextArea; 003/T012 |
| passwordField (novo) | initialValue ilustrativo e modo de mascaramento tipado; não credencial real, não logar texto; referência visual insuficiente mantém fidelidade bloqueada | PasswordFieldSkin/Document; 003/T012 |
| comboBox (novo); UF (preset) | options com IDs, selectedOptionId inicial, enabled; UF é lista/preset ordenado, não Kind distinto | EFLComboBox/Uf; 003/T012 |
| table (novo) | Colunas/células/linhas com IDs/ordem e tipos explícitos, seleção inicial simples/múltipla contratada; não texto CSV ou Report disfarçado | EFLTable/Header; 003/T012 |
| report (existente) | title/rows anteriores preservados; variante tipada para Report, Report2, Report3, Event, Event2, Notificacao, SelecaoMultipla e InternalReport; tipos de linha, altura e seleção iniciais por contrato | 012 piloto e 003/T014 |
| numericKeyboard (existente) | targetComponentId preservado, layout completo 14 teclas 7–9/4–6/1–3, setas/0/TAB/LIMPAR; comando efetivo distinguido da legenda | 013; 003/T013 |
| alphaKeyboard (novo) | targetComponentId, layout/modificadores iniciais tipados, base 800 × 200; calculadora como variante/helper local comprovado, sem negócio externo | EFLVirtualKeyboard/Calculator; 003/T013 |
| progressBar (novo) | mínimo/máximo/valor ilustrativo e legenda; mínimo < máximo, valor no intervalo; sem consulta externa | EFLProgressBar; 003/T014 |
| panel (novo) | variante/skin e capacidades de contêiner/scroll/clip; filhos via ParentID, não widgets persistidos | EFLPanel/Panel2/ScrollPane; 003/T015 e 004/T009–T010 |
| dialog (novo) | Composição de conteúdo/estado inicial tipado e overlay local; sem recriar janela SO ou serviços; faltas de evidência bloqueiam homologação | DialogSkin; 003/T015 |
| Frame e 014–018 (presets) | Frame de tela 800 × 600; conjuntos comuns com valores ilustrativos, fontes/skins da mesma biblioteca | 004/T011 |

Defaults geométricos das famílias novas, além dos comprovados acima, **a determinar por evidência em 003/T009 antes de liberar cada família**; não inventar dimensões por analogia. Referências/campos de tabela, linhas Report e grupos exigem IDs estáveis, cópia profunda, remapeio e validação; não liberar enum sem cadeia completa de suporte.

### Contrato de aparência e biblioteca

Biblioteca skin é estática gerada e pura; recursos são resolvidos pelo adaptador e snapshot. Chave de skin inclui família/nome/versão; estado tem chave original e semântica mapeada explicitamente. Cada campo resolve override presente → skin/estado selecionado → default comprovado da família. Estados ausentes só herdam onde o Java comprova regra; não sintetizar pressed/focused/disabled como equivalentes. Estado runtime não é override persistido.

Manter fonte/face/tamanho efetivos, RGBA, imagens/hashes, spacement, baseline, alinhamento, quebra e clipping. Flags XML bold/italic ficam na proveniência: SkinBuilderImpl calcula estilo sem aplicá-lo ao Font. Não transformar declaração em prova visual, nem sintetizar face sem comparação. Recursos/delegates ausentes geram diagnóstico localizado de fidelidade bloqueada; fallback explícito pode permitir rascunho/teste exploratório, mas **nunca** encerrar homologação.

O motor partilhado não usa widget Fyne como default. Texto no PPT continua shape editável; scroll/export representa recorte do estado inicial ou do quadro de roteiro, preservando todo conteúdo no projeto. Rich text/HTML/código arbitrário não são introduzidos; multiline de TextArea e linhas Report são expansão intencional da restrição v1 de uma linha. Agrupamento de seleção atual continua transitório; ParentID de composição é a extensão tipada descrita em Tela, não serialização da seleção.

### Capacidades e gate de liberação

Declarar suporte por família/propriedade/estado/consumidor e distinguir implementado, teste não executado, homologação pendente e bloqueado. Tipo novo só é utilizável no alvo após NewComponent/default, Validate, copy, duplicação/remoção/Undo, codecs, inspector/paleta, cena/prévia e diagnóstico/export correspondentes. A matriz de [research.md](../research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19) é a referência de cobertura. Perfis incompletos aparecem bloqueados; suporte no desktop não habilita Java ou implica PowerPoint homologado.

<a id="contrato-codigos-barras-2026-10-23"></a>
## Contrato planejado — códigos de barras (2026-10-23)

**Status**: contrato de especificação, não código implementado. Dono do modelo e das regras: [003/FR-013–FR-015](../../003-catalogo/spec.md#codigos-barras-2026-10-23). Detalhes de interação/visual em [031](../../031-codigo-barras-efl/spec.md) e [032](../../032-codigo-barras-processado-efl/spec.md); tarefas nos domínios, não duplicadas aqui. A versão de catálogo alvo é **3**, para que leitores do catálogo 2 recusem tipos que não reconhecem; formato/compatibilidade no [contrato de Projeto](project-format.md#formato-codigos-barras-2026-10-23).

### Famílias e dados tipados

| Kind planejado | Propriedades persistidas | Defaults / capacidades |
| --- | --- | --- |
| `barcodeInput` | `label` string | Rótulo “Número do Código de Barras:”; nova instância 610 × 60; oito valores runtime vazios. Entrada/foco/captura na prévia livre; sem typeText de roteiro nesta entrega. |
| `barcodeProcessed` | `blocks` array ordenado não vazio; bloco `{id,title,rows}`; linha `{id,labelKind,value}` | Um bloco com título vazio, duas linhas `cnpj`/`name` e valores vazios; nova linha `name` vazia. Proposta geométrica 530 × 120, não medida observada do anexo. Somente leitura na jornada; combos no inspector. |

`title` e `value` são strings Unicode e podem ser vazias. `rows` deve ser não vazio em cada bloco; `id` é UUID e não se repete no documento. Arrays preservam ordem, não têm interpretação por vírgulas. Tipos de rótulo: `cpf` → “CPF:”, `cnpj` → “CNPJ:”, `name` → “Nome:”, `tradeName` → “Nome fantasia:” e `corporateName` → “Razão social:”. Não persistir outra string de label que possa divergir do enum. Rótulos repetidos são permitidos; enum desconhecido é erro, sem aproximação. Troca de enum conserva literalmente `value`; nenhuma máscara, consulta ou validação fiscal implícita.

Identidade/FileID/bounds/estilo/habilitação são comuns à instância de Tela. Copiar snapshot conserva IDs; duplicação gera IDs novos no composto, blocos e linhas. Valores capturados, campo ativo, sucesso/erro, modo de leitura e caches são runtime por instância, nunca dados de `properties`. Processado não contém referência implícita a entrada e não é alimentado automaticamente por ela.

### Representações de boleto e normalização

Somente boleto de cobrança: leitura óptica de **44 dígitos** e linha digitável de **47 dígitos**. Não aceitar a condição de 48 posições do setter do exemplo: ela conflita com os recortes de 47. Arrecadação de 48 posições, QR/PIX e câmera não são cobertos.

Representação canônica contém apenas `0`–`9` ASCII. Na colagem integral, remover somente espaços ASCII e pontos; recusar outros caracteres ou quebra interna. Em modo de leitora, receber 44 dígitos numa captura e submeter com Enter final separado do conteúdo; não inferir pelo intervalo entre teclas nem dividir antes de validar. Colagem integral de 44 ou 47 usa a mesma validação; colagem parcial em campo ativo respeita seu limite e não trunca excedentes silenciosamente.

Os cinco blocos lógicos são `C1` (10), `C2` (11), `C3` (11), `C4` (1) e `C5` (14). Os oito campos visuais dividem os três primeiros em `5+5`, `5+6`, `5+6`, seguidos de C4 e C5. Os DVs locais são os últimos dígitos de C1/C2/C3; C4 é o DV geral; C5 não tem DV próprio. Prefixo visual de cinco dígitos ainda não permite validar o bloco.

Recortes a seguir usam **posições iniciando em 1, limites inclusivos**:

- **Digitado → lido**: `C1[1..4] + C4 + C5 + C1[5..9] + C2[1..10] + C3[1..10]`, total 44. Validar três DVs locais e o geral antes de aceitar.
- **Lido → digitado**: C1 base = posições lidas `1..4 + 20..24`; C2 base = `25..34`; C3 base = `35..44`. Acrescentar a cada base seu DV Módulo 10; C4 = posição lida 5; C5 = `6..19`. Validar o DV geral lido antes de converter.
- Validar comprimento/caracteres antes de recortar; não aceitar sinal numérico ou converter toda a string em inteiro, pois isso perde zeros. Alteração de qualquer parte invalida resultados derivados/código anteriormente aceito.

### Premissas matemáticas locais

As classes dos módulos não vieram no material fornecido. Estas regras são **premissas explícitas para a implementação local**, não prova de equivalência SISPL. A comparação com classes/vetores independentes reais permanece gate para essa afirmação, sobretudo SIGCB.

- **Módulo 10**: da direita para a esquerda, pesos alternados 2 e 1; somar os algarismos de cada produto. DV = `(10 - (soma % 10)) % 10`. Comparar com último dígito do bloco; não calcular sobre o DV incluído.
- **Módulo 11 do boleto**: remover a posição lida 5; aplicar pesos 2 a 9 ciclicamente, da direita para a esquerda, sobre as 43 posições restantes. Candidato = `11 - (soma % 11)`; candidatos 0, 10 ou 11 tornam-se 1; os demais permanecem. Comparar com posição lida 5. Um DV de C4 isolado só pode ser julgado após a composição completa.
- **Módulo 11 SISPL experimental**: mesmos pesos 2 a 9; candidato = `11 - (soma % 11)`; candidatos 10 ou 11 tornam-se 0. Não dizer “resto 11”: resto de divisão por 11 está entre 0 e 10. Essa premissa não é intercambiável com a regra de boleto e não condiciona a aceitação matemática do boleto.

Se as classes reais divergirem, rever o contrato e as fixtures antes de anunciar compatibilidade, em vez de modificar algoritmo para fazer um teste sintético passar. Nenhum algoritmo é replicado na UI, no codec ou no exportador.

### Dados derivados e classificação

Somente após aceitar a entrada, expor banco `lido[1..3]`, valor em centavos exatos `lido[10..19]`, fator numérico `lido[6..9]` e CAIXA quando banco for `104`. `0000` significa fator ausente; não usar somente o primeiro algarismo como teste de existência. Não converter fator em data civil sem ciclo contratado, inclusive reinício do fator em 2025. Valor zero continua zero, sem inferir cobrança ou valor a pagar.

Classificação SIGCB experimental segue os recortes do exemplo: banco 104, posição lida 30 em `{1,2}`, posição lida 34 entre 1 e 5, módulo SISPL de `lido[20..25]` igual ao dígito em 26 e de `lido[20..43]` igual ao dígito em 44. Em código aceito e classificado, número do título = `lido[27..44]`; nos demais, título indisponível. Para entrada parcial/inválida, todos os derivados são indisponíveis, sem erro de substring ou reaproveitamento do anterior. Não produzir nome/CPF/CNPJ/razão social: esses campos não são derivados por essas regras.

### Casos de referência documentais

| Fixture | Conteúdo / alteração | Resultado esperado |
| --- | --- | --- |
| F01 linha digitável | `00190500954014481606906809350314337370000000100` | 47 dígitos; três DVs locais e geral válidos sob as premissas; banco 001, valor 100 centavos, fator 3737, não CAIXA/SIGCB |
| F01 lido | `00193373700000001000500940144816060680935031` | 44 dígitos; conversão exata ida/volta para F01 |
| F02 lido sintético | `10498100000000001001234560000100010000000003` | Banco 104, 100 centavos, fator 1000; SIGCB verdadeiro **somente sob a premissa SISPL experimental**; título `000100010000000003` |
| F02 digitado sintético | `10491234566000010001000000000034810000000000100` | Conversão correspondente de 47; não é vetor independente de homologação |
| F03 DV local | Trocar o dígito da posição 10 da linha F01 de 5 para 4 | Falhar no bloco 1 e não avançar; repetir a mutação dos DVs nos blocos 2/3 |
| F04 DV geral | Trocar somente posição 5 do lido F01 de 3 para 4 e convertê-lo estruturalmente para linha | DVs locais permanecem válidos, mas conclusão/leitura falham pelo DV geral |
| F05 inválidos estruturais | Vazio como captura submetida; 43/45/46/48 dígitos; sinal, letras, Unicode numérico ou quebra interna | Recusar captura; modal especificado, sem preenchimento parcial |
| F06 módulo em limites | Bases `0`, `6`, `14`, `1` para Módulo 11 | DV boleto 1/1/1/9; DV SISPL experimental 0/0/0/9 |

Fixtures ilustram aritmética; não representam autorização de pagamento ou documento pagável. Uma verificação isolada por script durante a redação não é teste do produto nem validação contra biblioteca bancária real.

### Interação, saída e gates

Erro manual local conserva blocos já válidos e retorna ao início do par; erro geral digitado bloqueia conclusão e retorna ao primeiro campo com rascunho preservado. Captura integral inválida limpa os oito valores/resultados da instância e abre “Código de barras inválido” / “OK, entendi”; confirmação devolve foco ao primeiro campo. Captura cancelada antes de submeter conserva rascunho anterior. Entrada aceita fica somente para leitura até reiniciar; repetir visita começa vazia. Teclado numérico usa o subcampo ativo sem persistir esse índice.

Os combos de 032 são autoria, não widgets exportados. PPT preserva formas/textos editáveis, mas não calcula DV, lê dispositivos, aceita digitação ou faz processamento; roteiro `typeText` nos dois tipos não é suportado neste incremento e recebe diagnóstico. Entrada aceita não cria ação de Jornada nem alimenta ficha processada. Não reintroduzir Java/EFL, dependências externas ou dados pessoais em logs. Imagens ausentes impedem afirmar fidelidade, não autorizam inventar referência visual.

<a id="contrato-menu-aplicacao-2026-10-23"></a>
## Contrato planejado — Menu aplicação (2026-10-23)

**Autoridade**: [003/FR-016 e modelo canônico](../../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23). [033/spec](../../033-menu-aplicacao-efl/spec.md) e [DESIGN](../../033-menu-aplicacao-efl/DESIGN.md) detalham interação/visual, não outro domínio. Catálogo **3** planejado inclui `barcodeInput`, `barcodeProcessed` e `applicationMenu`; não criar catálogo4 por acrescentar 033 ao mesmo lote ainda planejado.

| Conceito | API acordada | Regras de contrato |
| --- | --- | --- |
| Kind/união | `ApplicationMenu="applicationMenu"`; `Properties.Menu *ApplicationMenuProperties` | Só o braço Menu presente nesse Kind; não adaptar a Report/panel ou cadastrar Kind por submenu. |
| Menu | `ApplicationMenuProperties{Root SubmenuProperties; InitialMenuID string}` | Root por valor; inicial é UUID existente, default Root.ID. |
| Submenu | `SubmenuProperties{ID,Name,Title,Layout string; Items []MenuItem}` | Menus vazios válidos; raiz layout root, descendentes standard/financial/bolao; ordem preservada. |
| Item | `MenuItem{ID,Label,Legend,Category,Variant string; Enabled bool; Submenu *SubmenuProperties; TargetScreenID,MovementCode,MovementOption,MovementName,Shortcut string}` | Label é o rótulo do botão; Legend é a legenda ilustrativa Unicode literal, inclusive vazia, independente de MovementName, sem fallback para metadados. Submenu nil é folha; destino vazio é ausência; submenu+destino não vazio é inválido; Enabled explícito. |

Modelos em arquivos separados application_menu.go/submenu.go de Catálogo, composição em Screen, wire no adaptador projectzip; destinos são responsabilidade de Jornada. Factory recebe callback supplyID para todos os UUIDs novos e inicializa Enabled=true; reader não reaplica defaults a valores confirmados. Nome/título/legenda/categoria/metadados podem ser vazios ou Unicode. Category é texto editável; Variant é referência ao estilo suportado, não enum Java executável. Estilo desconhecido gera diagnóstico de fidelidade/capacidade, não fallback silencioso ou código.

UUIDs de menu/submenu/item não se repetem entre si nem no documento. Limites por árvore: <=16 níveis contando raiz=1, <=1000 itens no total, inclusive desabilitados. Rejeitar ciclos em memória antes de recursão ilimitada; snapshot copia profundamente e mantém identidade, duplicação em Tela regenera nós e remapeia inicial. TargetScreenID é referência autorada, não fórmula; nenhuma categoria/variante/legenda/Shortcut/Movement* identifica destino implicitamente.

Seed único conserva grandes ramos de serviços e jogos e aninhamento do Java/XML fornecido, com raiz esquerda/direita, grade standard, financial140x80 e bolaoDigital de três botões+legendas. Não repetir XML completo em contratos nem substituir dados não acessíveis por conteúdo inventado. Geometrias de outras partes exigem fonte Java; sem imagem recebida não se alega fidelidade visual.

O layout `bolao`, não o Name editável, seleciona sempre as três posições de botões com Legend; mais de três itens exige diagnóstico e recorte sem poda. BOLÃO ordinário usa `standard` no seed, sem mudar árvore/rótulos. BOLÕES DIGITAIS usa `bolao` e as Legend exatas, na ordem: “Visualizar cotas digitais para venda na lotérica.”, “Consulta bolões de concursos abertos para comercialização.” e “Consulta bolões de concursos encerrados e não prescritos.”. Navegação permanece na árvore autorada; não buscar nomes na raiz para criar atalhos entre irmãos.

Prévia livre possui apenas vista raiz ou submenu selecionado com caminho local. Voltar remove nível, Cancelar restaura raiz, Sair encerra apenas prévia. Habilitação efetiva impede ativação bloqueada; folha com destino passa por Jornada, sem destino só informa seleção simulada. Nenhum timer/periférico/conta/autenticação/rede/QR/dispatch é executado. Canvas/quadro/PPT estático usam InitialMenuID confirmado; runtime nunca é persistido ou exportado como estado autorado. TypeText não é suportado; click demonstrativo continua visual até decisão/contrato de item metadata pela coordenação.

**Gate**: 003/T021 → 004/T017; codec002/T015, destinos005/T012 e quadros006/T012 → integração001/T140–T143. Oferta em Itens legados requer cadeia completa, não enum. Manual/aceite001/T144 é posterior à preparação técnica e não foi realizado; checklists históricos/customizados não são modificados pela autorização de implementação.

## Resolução visual aplicada — 2026-10-24

No recorte de componentes existente, Style vazio resolve o default do legado na apresentação; os defaults genéricos v1 acima são históricos. Variantes de botão mapeiam `botao_confirmar`, `botao_cancelar`, `botao_sim` e `botao_nao`; rótulo usa `label_texto_azul`, título `label_titulo_1`, campo `textField/type1`. Relatório/teclado/menu usam partes derivadas com suas próprias skins. Seleção mantém marcadores quando nenhuma skin de toggle foi escolhida. Overrides explícitos são preservados.

Novo teclado usa 200 × 350 e 14 teclas; carregar/gerar cena não altera bounds de instâncias anteriores. LIMPAR produz Backspace; TAB/setas são comandos de runtime, sem novos campos persistidos. Este incremento mantém o schema do código fornecido, sem promover o alvo documental de formatos posteriores a implementação concluída.

## Menu/XML e skins por item — 1.0.0.2

A fonte ativa é `material-legado/menu/menu.xml`, preservada em ISO-8859-1. Factory mantém ordem, espaços e metadados dos elementos ativos; comentários não geram itens. Variant vincula a chave exata do skin.xml, não um caminho inferido. O resolvedor usa imagens default/touched/disabled originais, com disabled prevalecendo sobre pressão. Conteúdo e limites em 033/FR-011–FR-012; a existência de uma skin não autoriza inserir item no seed.

## Menu e referências de movimentos — 1.0.0.3

Title vazio permanece vazio; não derivar nome como título ausente na fábrica. Itens continuam identificados por UUID: ramo abre submenu e folha simula seleção ou segue destino explícito. Código, nome, opção e atalho são editáveis como metadados literais; nome de movimento distingue referências que compartilham código. Sua exposição não comprova validação do cadastro movimento.xml ausente. Rodapé e faixas são primitivas derivadas, nunca novos itens persistidos.
