# DESIGN — Faixa de título da página

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [004-tela](../004-tela/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Design/page.css). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir adicionar uma faixa decorativa com título de página editável. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Referência: largura da página e altura 28, origem x=0, y=100; texto 18 com margem direita 48. Adaptar a páginas de outras dimensões.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| recurso | faixatituloprincipal.png | Style.BackgroundAssetID |
| texto | #FFFFFF | Style.TextColor |
| fonte | 18, destaque | Style.FontSize / família title |
| alinhamento alvo | à direita | Style.TextAlign; esquerda / centro / direita |

Conteúdo e opções configuráveis: Texto do título, faixa de fundo, posição, dimensões, fonte e cor.

## Estados

Visível; removida; texto vazio; selecionada.

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

- Sem imagem, título muito longo e página estreita. Remoção deve deixar de exportar tanto a faixa quanto o texto removido.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/application/presets.go](../../apps/gerador-jornada/internal/application/presets.go)
- [internal/adapter/builtinassets/background.go](../../apps/gerador-jornada/internal/adapter/builtinassets/background.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir a faixa original com um título editável por cima.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Oferecer alinhamento horizontal do texto da faixa: esquerda, centro e direita; preservar X, Y, largura e altura da faixa ao mudar somente essa propriedade.
- **FR-009**: Conservar alinhamento ao salvar/reabrir e reproduzi-lo no editor, prévia, demonstração e PPT, mantendo texto editável. A faixa pronta inicia com texto à direita.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.
