# DESIGN — Campo de texto — EFLTextField

## Aplicação na revisão 1.0.0.1

Campo 150 × 25, textField/type1: Arial 12, fundo branco/texto preto; foco azul #10508E/texto branco. legacyEntry usa captura Fyne transparente sobre a cena, preservando cursor, seleção e clipboard.

**Revisão ainda necessária:** Os formatos existentes validam intenção/comprimento/data; não implementam um motor incremental de máscaras nem consulta/validação fiscal. TextArea/PasswordField continuam tipos planejados, não variantes disponíveis deste campo. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLTextField.java](../../material-legado/elemetos_componente/component/EFLTextField.java), [documentos](../../material-legado/elemetos_componente/java) e [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLTextField/index.js) é referência secundária histórica.

## Responsabilidade visual

O componente deve permitir desenhar e experimentar entradas de texto com intenção de formato definida. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Referência legada: 150 × 25; manter área de entrada configurável e geometria comum às representações.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| fundo | #FFFFFF | Style.FillColor |
| texto | #000000 no estado normal; #FFFFFF em foco | Style.TextColor |
| borda | #000000 na referência | Style.BorderColor |
| fonte | 12 | Style.FontSize |

Conteúdo e opções configuráveis: Valor inicial, formato, posição, dimensões, fonte, texto, fundo e borda.

## Estados

Vazio; preenchido; foco; entrada recusada; selecionado no editor.

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

- Valor vazio, colagem com caracteres inválidos em campo numérico, teclado sem alvo e formatos sem validação fiscal. Formato não significa consulta de CPF ou CEP.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback exploratório explícito não homologa. Conferir glifos, métricas, cursor e texto longo com face real.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/application/preview.go](../../apps/gerador-jornada/internal/application/preview.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Editar valor inicial e escolher texto, número, celular, CEP, CPF, data, DDD ou NIS.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

Inventariar textField default/focused/disabled, face/RGBA/spacement/alinhamento e pintura efetiva; tokens acima são base histórica, não prova de skin homogênea. Fonte XML bold/italic não é estilo aplicado automaticamente pelo builder; TEXTBOX_BORDER=10 não prova padding de dez em cada lado. Máscaras e documentos recebem fixtures locais sem consultas externas. TextArea/PasswordField têm tipos/skins próprios em 003/T012 e não herdam visual por analogia; Password/Dialog têm referência insuficiente para homologação. Estado de edição runtime não entra no arquivo. [Matriz C06–C08/C26](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19), 001/T111–T116; mesma face/cena nas três superfícies, sem Entry opaca de skin Fyne cobrindo o desenho.
