# DESIGN — Relatório/lista — EFLReportEvent

## Aplicação na revisão 1.0.0.1

Relatório 300 × 345 com título e linhas simples. Cena recompõe painel, cabeçalho e rodapé originais; título/linhas usam label_report_titulo/label_report e permanecem editáveis no PPT.

**Revisão ainda necessária:** As oito classes Java de Report, linhas estruturadas, seleção/rolagem e diagnóstico integral de overflow não estão implementados. Hoje a cena limita as linhas ao corpo visível; isso não conclui FR-008–FR-009. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLReportEvent.java](../../material-legado/elemetos_componente/component/EFLReportEvent.java) e demais classes Report/[skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLReportEvent/index.js) é referência secundária histórica.

## Responsabilidade visual

O componente deve permitir organizar um cabeçalho e linhas de informação visual. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Base do catálogo/React: 300 × 345; topo 26, título 18 e rodapé 32 na referência React. Java confirma TAMANHO_CABECALHO=26 e TAMANHO_RODAPE=32, mas também ALTURA_ENTRADA_SIMPLES=20 e ALTURA_ENTRADA_TITULO=26, além de alturas customizadas. **Título 18 não é medida universal Java comprovada**. Renderer atual com header 28/corpo simplificado diverge; reconciliar por parte/classe no piloto 012/T001, sem inventar que a alteração já ocorreu. Bounds externos existentes não são redimensionados na migração.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| topo e rodapé | label_tabela_cabecalho.png e icon_tabela_rodape_2.png | recursos originais incorporados; textos independentes |
| título | #3C3C3C sobre faixa clara | cabeçalho editável |
| corpo | #B5C7DF | Style.FillColor / renderer |
| linhas | #A6B9D7 | token de alternância |
| texto | #5064AA, 12 | Style.TextColor / FontSize |

Conteúdo e opções configuráveis: Título, linhas, ordem, dimensões, texto, cores e decoração de cabeçalho/corpo/rodapé.

## Estados

Lista vazia; preenchida; linha destacada; conteúdo maior que a área; seleção de edição.

| Situação | Resultado visual e regra |
| --- | --- |
| Referência | Mostrar os tokens e a geometria definidos acima, com o conteúdo da fixture. |
| Seleção no editor | Contorno e alças ficam em camada transitória; não aparecem no PPT ou nos fontes Java. |
| Foco por teclado | Exibir indicação perceptível sem deslocar a composição; texto digitado deve seguir o foco correto. |
| Operação inválida | Explicar a causa junto ao contexto; conservar a aparência do último estado confirmado. |
| Estado vazio ou ausente | Distinguir conteúdo vazio permitido de recurso/referência quebrados; não preencher dados fictícios silenciosamente. |

## Personalização segura

1. Alterar os valores de referência deste DESIGN.md e identificar o token e o requisito afetados.
2. Confirmar se a propriedade pertence ao documento ou somente ao tema da interface; não acrescentar estado de UI ao arquivo .jornada.
3. Atualizar a implementação responsável por essa propriedade. Alterar Markdown, isoladamente, não muda o EXE.
4. Repetir os cenários da spec e comparar antes/depois em 100%, 150% e 200% de DPI.
5. Para elementos de tela, salvar/reabrir e conferir a mesma personalização no editor, prévia e PPT. Estilo Java sem mapeamento EFL homologado deve produzir diagnóstico, nunca aproximação silenciosa.

## Limites e acessibilidade

- Sem linhas, linhas longas, ordem alterada e conteúdo excedente. Definir rolagem na prévia e representação do recorte visível no PPT; não descartar dados do projeto.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback exploratório explícito não homologa. Conferir glifos, métricas, quebra e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Editar título e linhas ordenadas, preservando IDs independentes dos textos.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Piloto pós-v1 e variantes — 2026-10-19

Tokens acima documentam a referência inicial, não substituem skins panel/label/scroll por cores aproximadas. Inventariar por parte/estado/hash fonte/face/RGBA/alinhamento/baseline/spacement/recorte. XML bold/italic é declaração: builder calcula style sem aplicá-lo ao Font. Título/cabeçalho/corpo/rodapé são partes internas derivadas, não dezenas de Component persistidos. Seleção de linha e scroll runtime não alteram rows/IDs iniciais; PPT usa recorte do quadro e mantém textos/formas visíveis editáveis, conteúdo completo permanece no projeto.

Report/2/3/Event/Event2/Notificacao/SelecaoMultipla/Internal precisam contratos/fixtures próprios; altura/linhas/seleção não são herdadas por analogia. [Matriz C11–C14](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19), [plano local](plan.md) e [tarefas locais](tasks.md). Fontes por nome disponíveis em outras_fontes não resolvem licença/métricas/asset/delegate; ausência mantém homologação bloqueada. Não rasterizar para ocultar divergência, nem apresentar resultado do renderer novo como referência antiga.
