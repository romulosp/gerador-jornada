# DESIGN — Seleção — EFLRadio

## Aplicação na revisão 1.0.0.1

Seleção 180 × 60 com opções identificadas e orientação vertical/horizontal. Sem skin explícita usa círculo e texto #4864A0/tamanho 16; uma skin toggle escolhida altera a decoração da opção. Escolhas da prévia são locais.

**Revisão ainda necessária:** Grupos tipados mais amplos, toggleButton independente e matriz completa de estados de habilitação/foco não estão homologados. Não confundir o visual de marcador padrão com uma implementação de toda a família Java de toggles. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLToggleButton.java](../../material-legado/elemetos_componente/component/EFLToggleButton.java), [EFLButtonGroup.java](../../material-legado/elemetos_componente/component/EFLButtonGroup.java) e [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [Radio React](../../material-legado/src/components/ListComponents/EFLRadio/index.js) é referência secundária histórica, não prova de construtor Java.

## Responsabilidade visual

O componente deve permitir oferecer uma escolha única entre opções identificadas. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Legado: 110 × 44 para duas opções; distribuir itens dentro da caixa configurada sem sobreposição.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #4864A0 | Style.TextColor |
| marcador | círculo com centro preenchido quando escolhido | renderer de seleção |
| fonte | 16 na referência legada | Style.FontSize |
| espaçamento | distribuído pela orientação | layout derivado dos bounds |

Conteúdo e opções configuráveis: Opções com IDs, textos, orientação, opção inicial, geometria e estilo tipográfico.

## Estados

Sem escolha; opção selecionada; foco; opção removida; seleção no editor.

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

- Textos repetidos com IDs distintos, nenhuma opção, opção inicial removida e largura insuficiente. O template Java legado inconsistente não comprova a API EFL.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback exploratório explícito não homologa. Conferir glifos, métricas e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Adicionar, editar e reordenar opções sem usar suas legendas como identidade.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

Marcador circular/tokens/tamanho acima descrevem Radio inicial, não todos os toggles Java. Novas skins toggleButton/grupos têm dimensões internas, faces e estados próprios; manter seleção simples de arquivos antigos e separar selected de touched/pressed com evidência. Negrito/itálico declarados no XML não provam aplicação pelo builder. Mesma cena/cor/fonte em canvas/prévia/PPT, sem choices nativas opacas alterando a aparência. [Matriz C05/D-11](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19); 003/T011 e 001/T111–T122. EFLToggleButtonUI ausente impede homologação dependente; nenhuma aprovação por similaridade.
