# DESIGN — Rótulo — EFLLabel

## Aplicação na revisão 1.0.0.1

Rótulo 150 × 20; default label_texto_azul, Futura Medium física, tamanho 12 e #5064AA. Fonte/cor/alinhamento explícitos prevalecem; fundo do rótulo padrão é transparente.

**Revisão ainda necessária:** Variantes compostas e cobertura integral dos rótulos do acervo não são homologadas por esta implementação. Longas palavras, recorte e DPI ainda exigem revisão visual. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLLabel.java](../../material-legado/elemetos_componente/component/EFLLabel.java)/[skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLLabel/index.js) é referência secundária histórica.

## Responsabilidade visual

O componente deve permitir posicionar uma informação textual editável sem transformá-la em imagem. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Referência legada: 150 × 12; permitir aumentar altura e largura sem truncar o conteúdo persistido.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #5064AA (label_texto_azul) | Style.TextColor |
| fonte | 12 | Style.FontSize |
| fundo | transparente | ausência de preenchimento |
| realce de edição | #2563EB | tema do editor; não exportar |

Conteúdo e opções configuráveis: Texto, nome lógico, posição, tamanho, fonte, tamanho de fonte e cor.

## Estados

Normal; vazio; selecionado para edição; texto maior que a caixa.

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

- Acentos, aspas, emoji, quebras, caixa pequena e texto vazio.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback explícito de rascunho não homologa. Conferir glifos, métricas, baseline e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Inserir texto Unicode e alterar seu conteúdo preservando os identificadores.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

Tokens/150 × 12 acima são referência React histórica, não default universal Java (catálogo inicial label 150 × 20). Cada skin label tem face/tamanho/RGBA/métricas/alinhamento próprios a inventariar, sem resize de instâncias existentes. ValorSurpresaSuperSete é composição label alinhado à direita + toggle opcional comprovada por construtor, sem resultado de negócio automático. SkinUtil.TEXTBOX_BORDER=10 não prova padding simétrico. Negrito/itálico XML são declarados, não comprovados pelo builder; aplicar D-11 da [matriz C03–C04](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19). Todas skins e partes textuais devem ser verificadas em canvas/prévia/PPT; 003/T011, 004/T011 e 001/T111–T122 mantêm tarefas abertas.
