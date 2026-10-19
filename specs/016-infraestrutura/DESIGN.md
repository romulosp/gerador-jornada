# DESIGN — Indicadores — InfraInfo

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [004-tela](../004-tela/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/InfraInfo/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir representar indicadores ilustrativos de infraestrutura na tela desenhada. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

78 × 36, três linhas de 12; referência x=722, y=6.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #424344 | Style.TextColor |
| fonte | 11 | Style.FontSize |
| espaçamento | 12 | Bounds.Y |
| fundo | transparente | ausência de caixa |

Conteúdo e opções configuráveis: Textos dos três indicadores, geometria, fonte e cor.

## Estados

Ilustrativo; personalizado; selecionado; removido.

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

- Sem rede, indicador removido e texto extenso. Não representar conectividade como verificada nem inserir integrações externas.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/application/presets.go](../../apps/gerador-jornada/internal/application/presets.go)
- [internal/adapter/desktop/presets.go](../../apps/gerador-jornada/internal/adapter/desktop/presets.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir as linhas GATEWAY, DNS e HTTP no cabeçalho ou em outra posição escolhida.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.
