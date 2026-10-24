# Contract: exportação PowerPoint

## Implementação conferida — 1.0.0.1

`presentation.Component/ComponentState` resolve os dez Kinds atuais; o adaptador PPT consome suas primitivas. Decorações originais são mídias; legendas/valores são formas/textos editáveis. A revisão 1.0.0.1 transporta família SFNT real, flags físicas de fonte, alpha parcial de primitivas, centralização e espaçamento entre linhas. Skin/fonte explícita não resolvida bloqueia Build antes de publicar.

Não há embedding de fontes no PPT. A fonte precisa estar disponível/licenciada no computador de destino; XML válido não comprova ausência de substituição no Office. Recorte e equivalência de métricas de todos os tipos/estados permanecem requisitos a homologar. Menus e códigos de barras são estáticos no PPT: não calculam DV nem executam a árvore. Evidência estrutural/arquivos gerados e roteiro manual na seção 22 do quickstart.

Código atual: [wire do projeto](../../../apps/gerador-jornada/internal/adapter/projectzip/wire.go), [cena](../../../apps/gerador-jornada/internal/presentation/scene.go) e [exportador PPT](../../../apps/gerador-jornada/internal/adapter/pptx/export.go). Evidência executada na [seção 22 do quickstart](../quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). O restante do contrato conserva regras normativas e evolução planejada; divergências apontadas aqui não foram apagadas dos requisitos.


**Version**: 1 | **Requirements**: FR-014 a FR-016, FR-020, FR-021, FR-027; US3 e US4.
**Compatibility target**: PowerPoint desktop 2016 e Microsoft 365 no Windows; extensão `.pptx`.

## Entrada e saída

Receber snapshot validado, includeExploration boolean e scenarioIds ordenados. Exigir ao menos um percurso selecionado. O padrão inclui exploração e todos os roteiros não vazios. Validar todos os roteiros selecionados; não omitir inválidos.

Gerar um arquivo PPTX sem macros, links externos obrigatórios, vídeos dos controles ou dependência do Office para criação. Usar ZIP/OPC, PresentationML e DrawingML. A conformidade estrutural não significa animação homologada.

Partes mínimas: tipos de conteúdo, relacionamentos raiz, apresentação, slide master, layout, tema, slides/relacionamentos e media para recursos decorativos. IDs e relacionamentos devem ser únicos e resolvíveis; todos os textos devem passar por escape XML.

## Organização dos slides

1. Menu inicial com entrada Explorar, quando selecionada, e uma entrada por roteiro selecionado.
2. Bloco de exploração: um slide por tela na ordem do projeto; o botão Explorar aponta à tela inicial, independentemente dessa ordem.
3. Blocos de demonstração: ocorrências de visitas na ordem do roteiro, seguidas de slide final com links de Menu e Repetir.

Exploração não tem avanço temporizado. Botões navigate têm ação de link interno para o slide da tela destino; botões none não têm link. Cada tela exploratória oferece acesso ao menu fora da área lógica do desenho.

Demonstração não mantém links interativos nos botões do desenho: eles são parte da simulação. O avanço é o roteiro. Seu último slide não avança sozinho para outro roteiro. Repetir aponta à primeira ocorrência.

O tamanho do slide é único para todo o arquivo. Usar a maior largura/altura lógica das telas mais uma faixa inferior de 40 unidades para navegação de apresentação. Telas menores são centralizadas sem distorção; menu/final usam a área total.

## Coordenadas e editabilidade

Converter 1 unidade lógica para 9525 EMU (96 DPI). Tamanho de fonte em pontos é fontSize × 0,75; representar em centésimos de ponto quando exigido pelo formato. Não converter via dimensão da janela nem screenshot.

Uma família vira um grupo com formas/textos internos. Botão e campo têm caixa e texto separados; radio contém opções; report contém título/linhas; teclado contém teclas e textos. IDs de formas derivam de travessia determinística, mantendo mapa para os alvos de animação.

Textos e contornos de controles não podem ser rasterizados. Imagens são admitidas para background/skin decorativa sem texto funcional. Na base v1, fontes eram referenciadas por nome, não redistribuídas automaticamente. A expansão de 2026-10-19 abaixo exige face real e incorporação permitida ou requisito de ambiente explícito; substituição não pode ser anunciada como fidelidade homologada.

Geometria: erro máximo de 1 unidade lógica em posição/tamanho após conversão. Conteúdo, ordem visual, cores e destinos devem corresponder; diferenças de antialiasing não são falhas. Verificar legibilidade e overflow separadamente.

## Compilação das animações

O roteiro interpretado fornece ocorrências, tempos e estados. Cada ocorrência gera seu slide; auto-navegação também cria novo slide. Incluir apenas efeitos necessários e existentes no subconjunto validado na Sprint 1.

| Etapa | Representação |
| --- | --- |
| moveCursor | Cursor vetorial independente e trajetória linear entre posições, usando animação de movimento. |
| click | Cursor no alvo e círculo/realce vetorial com entrada/saída temporizadas. |
| typeText | Caixas de texto editáveis com prefixos progressivos, alternadas por visibilidade; a caixa final contém o valor completo. |
| wait | Intervalo temporal sem mudança de estado. |
| navigate | Encerramento temporizado da ocorrência e avanço para a próxima visita do roteiro. |

Para typeText, segmentar em clusters Unicode e atribuir duração proporcional a cada prefixo. Estado inicial do campo permanece visível até a primeira troca; a última caixa fica visível até fim da ocorrência. Texto vazio troca para caixa vazia ao fim da etapa. Preservar alinhamento e geometria de todas as caixas, sem deixar múltiplos prefixos visíveis.

Editar o conteúdo depois no PowerPoint pode exigir ajustar os prefixos da demonstração; o editor de origem continua sendo a fonte para regeneração consistente. Não prometer atualização automática entre caixas após edição manual do PPT.

A duração de cada ocorrência é soma de suas etapas. Para uma última ocorrência sem navigate, ao fim das etapas manter estado final por 1000 ms e avançar ao slide final. Usar `timing`/alvos de formas e `transition` com tempos definidos, sem Morph ou recursos exclusivos do Microsoft 365.

Não remover animação para contornar problema de compatibilidade sem emitir erro. Uma saída estática não atende a US4.

## Transições visuais — atualização de 2026-10-18

O diálogo de exportação permite `fade`, `push` (deslizar), `wipe` (revelar) e `none`. O padrão é fade com duração 300 ms; a escolha pertence à operação de exportação, não ao modelo de navegação. As transições são nativas, sem Morph, e não substituem o timing de cursor/clique/digitação.

Exploração, menu e encerramento não recebem `advTm`, nem avançam por clique no fundo. Links continuam sendo o mecanismo de navegação exploratória. Nas ocorrências demonstrativas, preservar `advTm` calculado pelo roteiro. Conferir a percepção e duração da entrada no PowerPoint separadamente do tempo das etapas; a conformidade estrutural não equivale a homologação visual.

Usar `p:fade`, `p:push` ou `p:wipe` no `p:transition`, com duração `p14:dur` e declaração de compatibilidade do namespace. [Referência Microsoft](https://learn.microsoft.com/en-us/dotnet/api/documentformat.openxml.presentation.transition?view=openxml-3.0.1).

## Diagnósticos

Além dos erros gerais: `PPT_TARGET_INVALID`, `PPT_RELATIONSHIP_INVALID`, `PPT_ANIMATION_UNSUPPORTED`, `PPT_TEXT_OVERFLOW` (warning) e `PPT_FONT_SUBSTITUTED` (warning de homologação). Apontar ocorrência/etapa quando houver falha temporal.

## Aceitação específica

Abrir sem reparo nas duas versões-alvo; inspecionar formas/textos; percorrer todos os botões; comprovar independência dos percursos; reproduzir cursor, clique e digitação com acentos; revisar tempos conforme SC-004. Reabrir após salvar pelo PowerPoint e repetir a navegação.

Fonte primária: [animações PresentationML](https://learn.microsoft.com/en-us/office/open-xml/presentation/working-with-animation) e [links internos](https://learn.microsoft.com/en-us/openspecs/office_standards/ms-oe376/7ff3db24-b7b9-4ffe-aa78-3ec47cab2489). O código será Go; o SDK .NET citado na documentação não é dependência do editor.

## Paridade de todo o acervo — expansão 2026-10-19

**Estado implementado na continuação:** exportador consome a biblioteca imutável
do perfil de `application.Document`, inclusive revisão arquivada exata, e recursos
validados por `visualprofile.ResourceData`; cache identifica SHA/face exata, não só
família/caminho. Referência desconhecida ou bytes divergentes não recebem fallback
do perfil corrente. Texto/formas continuam editáveis. **Não incorpora fontes**:
`Options.EmbedFonts` e `Options.StrictFidelity` retornam
`PPT_FONT_EMBEDDING_UNSUPPORTED`, `FontEmbedding=false` e `CompleteFidelity=false`.
O snapshot do perfil não altera essas capacidades. As exigências abaixo continuam
sendo alvo de paridade/homologação, não evidência de embedding ou aprovação Office.

FR-053–FR-054: consumir a mesma cena/resolver tipográfico e snapshot do perfil de canvas/prévia. Usar nome interno/face efetiva da fonte, não nome do arquivo e não Arial fixa. Calibrar tamanho lógico→pontos, baseline, parágrafos/quebra/line-height/spacement, alpha e clipping com a referência. Bold/italic XML não são automaticamente estilos efetivos: aplicar a decisão registrada por skin/face em research/D-11. Preservar bounds existentes inclusive teclados antigos, sem resize na exportação.

Imagem de skin é decoração; legendas duplas, células, linhas, cabeçalhos e teclas permanecem textos/formas editáveis por componente/grupo. Recorte deve mostrar conteúdo visível do estado inicial/quadro, não rasterizar o controle para esconder overflow. Todo conteúdo continua no projeto; a limitação estática do PPT é informada, não perda silenciosa. GIF/animação somente após uso/semântica comprovados e snapshot estático explicitado; não habilitar formato pelo nome sem contrato.

Incorporar fonte no PPTX quando formato/licença permitem, validando partes/relacionamentos de fonte e abertura real. Se incorporação for proibida/inviável, bloquear o modo de fidelidade portátil e informar as faces que precisam estar instaladas no ambiente de homologação; usuário deve conhecer essa dependência. Substituição detectada (`PPT_FONT_SUBSTITUTED`) impede concluir fidelidade mesmo quando só produz warning exploratório. Ausência de fonte/recurso/estado/capacidade exigida gera diagnóstico localizado antes de publicar pacote parcial; `PPT_COMPONENT_UNSUPPORTED`, `PPT_RESOURCE_MISSING` e `PPT_FONT_UNAVAILABLE` identificam causa e componente/skin/estado.

Cada nova família/variante requer capacidade explícita e fixture; sucesso OPC/XML não prova visual/edição nem homologa execução. Manuais nas duas versões-alvo conferem geometrias, fontes/recursos, seleção/editabilidade, navegação, tempos/revisitas e revisão/hash. A tolerância geométrica de até uma unidade não dispensa fonte/face/asset corretos; diferenças de rasterização ficam medidas e aceitas explicitamente, sem promessa de pixel idêntico não medida. Responsabilidades: 001/T116 (implementação), T118 (automatizado), T121 (manual) e T122 (aprovação).
