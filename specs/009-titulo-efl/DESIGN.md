# DESIGN — Título — EFLLabelTitulo1

## Aplicação na revisão 1.0.0.1

Título 300 × 28; default label_titulo_1, família física Futura Medium, tamanho 14 e #3E578D. O destaque genérico de tamanho 22 foi substituído pela skin; bounds antigos continuam explícitos.

**Revisão ainda necessária:** Referência Java/Office e todos os estados/variantes continuam sem homologação visual. A faixa de página é um preset separado desta família de título simples. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLLabel.java](../../material-legado/elemetos_componente/component/EFLLabel.java)/[skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLLabelTitulo1/index.js) é referência secundária histórica.

## Responsabilidade visual

O componente deve permitir destacar o assunto de uma área da tela com texto independente do título da página. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Referência React histórica: 150 × 40. O catálogo atual mantém bounds iniciais 300 × 28 e resolve label_titulo_1, fonte física Futura Medium de tamanho 14; o destaque genérico anterior de 22 foi substituído. Bounds e tamanho de fonte personalizados são preservados.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #3E578D (label_titulo_1) | Style.TextColor |
| tamanho | 14 na referência legada | Style.FontSize |
| peso histórico React/XML | negrito declarado; efeito depende de face e transformação Java | aparência comprovada por skin, não Bold implícito |
| fundo | transparente | decoração opcional |

Conteúdo e opções configuráveis: Texto, posição, dimensões, fonte, tamanho e cor; imagem decorativa opcional.

## Estados

Normal; selecionado; vazio; texto extenso.

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

- Título longo, múltiplos títulos, remoção da imagem decorativa e fontes indisponíveis.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback exploratório explicitado não homologa. Conferir glifos, métricas e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir um título independente, com identidade e texto próprios.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

14 React/22 do renderer inicial não definem um único tamanho Java de título. Variantes label_titulo e demais skins aplicáveis têm fonte/face/tamanho/alpha próprios, por inventário. Override explícito prevalece; campo não sobrescrito adota nova skin, sem mudar bounds. SkinBuilderImpl calcula flags bold/italic sem aplicar style ao Font; fonte física pode já ser face pesada/itálica. Distinguir declaração, transformação e referência visual na [matriz C03/D-11](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19). Fontes em outras_fontes não são ausentes por nome, mas licença/métricas e EFLLabelUI seguem pendentes. Tarefas 003/T011 e 001/T109–T122 conforme responsabilidade.
