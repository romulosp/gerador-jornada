# DESIGN — Botão — EFLButton

## Aplicação na revisão 1.0.0.1

Botão 110 × 60; variantes confirm/cancel/yes/no resolvem botao_confirmar/botao_cancelar/botao_sim/botao_nao quando Skin está vazia. Imagem e legenda permanecem separadas. Prévia usa legacyControl, com estado de pressão e foco.

**Revisão ainda necessária:** Duas legendas tipadas, cobertura de todos os estados/skins e comparação independente de métricas Java/Office permanecem pendentes. Hover não é um estado demonstrado por este incremento. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLButton.java](../../material-legado/elemetos_componente/component/EFLButton.java) e [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLButton/index.js) é referência secundária histórica. Delegate ausente não é evidência de pintura homologada.

## Responsabilidade visual

O componente deve permitir inserir uma ação com legenda e destino de navegação identificável. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Catálogo atual e referência React: **110 × 60**. A menção anterior a 160 × 44 era divergência documental, corrigida em 2026-10-19, não um default a migrar. Dimensões são editáveis e não variam ao trocar legenda/skin; arquivos existentes preservam exatamente seus bounds.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto padrão | #00FF00 em confirmar/sim; #FF8060 em cancelar/não | cor da skin; Style.TextColor é override explícito |
| decoração padrão | imagem da skin de cada variante; botao_aux2_released.png em confirmar/cancelar | resolvedor da skin; Style.FillColor é override explícito |
| fonte histórica React | 12, negrito declarado; efeito Java depende de face/skin | Aparência/face efetiva; não presumir Bold do widget |
| variantes | confirmar / cancelar / sim / não | ButtonProperties.Variant |

Conteúdo e opções configuráveis: Legenda, variante, destino, posição, tamanho, fonte, cor do texto, preenchimento, borda e imagem decorativa.

## Estados

Normal; ponteiro sobre o controle; pressionado; foco; selecionado no editor; destino ausente.

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

- Legenda vazia ou longa, destino removido, auto-navegação e sobreposição de controles.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente mantém fidelidade bloqueada; fallback explicitado só permite rascunho exploratório, não homologação. Verificar acentos, glifos, métricas e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir por clique ou arraste e editar a legenda sem mudar a identidade do botão.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

FR-008 inclui todas variantes XML button e [EFLTwoLabelButton](../../material-legado/elemetos_componente/component/EFLTwoLabelButton.java), com duas áreas textuais/estilos/métricas próprios e editáveis, não uma imagem de legenda. Tabela acima registra a base histórica; novos valores resolvem override presente → skin/estado → default comprovado. Inventariar default/touched/disabled sem inventar herança/foco/pressed. Fontes reais estão em xml/fonts e outras_fontes, mas presença não prova licença/métricas; SkinBuilderImpl ignora flags bold/italic calculadas ao obter Font. Matriz C01–C02/D-11 em [research.md](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19); regras 003/T011, cena/desktop/PPT 001/T111–T116 e homologação separada T120–T122.
