# DESIGN — Data e hora — DateInfo

## Aplicação na revisão 1.0.0.1

Data e hora é preset 126 × 32 de dois rótulos editáveis, com Arial 11/#424344 e valores ilustrativos estáticos. Cada rótulo recebe identidade e bounds próprios.

**Revisão ainda necessária:** Não há relógio vivo, consulta à máquina ou sincronização de horário. Comparação visual nativa e variações de texto/DPI continuam pendentes. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [004-tela](../004-tela/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/DateInfo/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir compor a identificação temporal da jornada com valores ilustrativos editáveis. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

126 × 32, fonte 11 e duas linhas. Posição sugerida no cabeçalho 800 × 600: x=374, y=6.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #424344 | Style.TextColor dos rótulos |
| fonte | 11 | Style.FontSize |
| linha da data | BSB - 01/01/2026 | LabelProperties.Text |
| linha da hora | 12:00:00 | LabelProperties.Text |

Conteúdo e opções configuráveis: Texto de cidade/data, texto de horário, posição, fonte, cor e espaçamento entre linhas.

## Estados

Exemplo; valores personalizados; seleção de uma linha; conjunto removido.

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

- Data ilustrativa inválida, horário longo e múltiplas inserções. O relógio vivo do legado é substituído por conteúdo determinístico nesta versão.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/referência ausente bloqueia fidelidade; fallback explícito de rascunho não homologa. Conferir face real/glifos/métricas e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/application/presets.go](../../apps/gerador-jornada/internal/application/presets.go)
- [internal/adapter/desktop/presets.go](../../apps/gerador-jornada/internal/adapter/desktop/presets.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir duas linhas, uma com cidade/data e outra com horário.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Aparência pós-v1 — 2026-10-19

Geometria/tokens acima são referência React da composição, não classe de relógio Java homologada. Reusar label/skin/face/alpha/métricas comprovados; fonte declarada XML não prova bold/italic efetivos do builder. Preservar textos/bounds/overrides antigos, sem hora automática ou mudança de tema do editor. 004/T011 e 001/T111–T116; [matriz C23/D-11](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19). Comparação manual por revisão/DPI permanece pendente.
