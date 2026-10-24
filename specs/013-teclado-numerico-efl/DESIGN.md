# DESIGN — Teclado numérico — EFLVirtualKeyboardNumeric

## Aplicação na revisão 1.0.0.1

Novo teclado 200 × 350 com 14 teclas Java: 7–9, 4–6, 1–3, < / 0 / >, TAB e LIMPAR. Bounds antigos são preservados e as partes escaladas. LIMPAR emite Backspace; setas respeitam cursor e TAB percorre campos habilitados na prévia.

**Revisão ainda necessária:** Teclado alfanumérico e variantes/overrides tipados de toda a expansão não estão disponíveis. O alvo barcodeInput existe, mas a matriz completa de bloqueio/avanço por par com TAB/setas ainda requer validação própria. Não tratar LIMPAR como clear-all. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: primária [EFLVirtualKeyboardNumeric.java](../../material-legado/elemetos_componente/component/EFLVirtualKeyboardNumeric.java)/[skin.xml](../../material-legado/elemetos_componente/xml/skin.xml); [React](../../material-legado/src/components/ListComponents/EFLVirtualKeyboardNumeric/index.js) é referência secundária histórica.

## Responsabilidade visual

O componente deve permitir compor uma entrada numérica com um teclado vinculado a um campo da tela. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Default React e do editor anterior: **200 × 280** (histórico). Java chama setSize(200,350); default atual: **200 × 350**, somente para novas inserções. **Migração não redimensiona teclados existentes**, nem altera IDs/alvo/overrides. Componente mantém o arranjo interno de referência abaixo escalado aos bounds explícitos, inclusive margens e área livre, sem reflow pela ordem das teclas.

| Parte Java | Bounds internos x,y,w,h na base 200 × 350 |
| --- | --- |
| 7, 8, 9 | x=5/64/123, y=5, 52 × 52 |
| 4, 5, 6 | x=5/64/123, y=60, 52 × 52 |
| 1, 2, 3 | x=5/64/123, y=115, 52 × 52 |
| <, 0, > | x=5/64/123, y=170, 52 × 52 |
| TAB | x=5, y=225, 80 × 60 |
| LIMPAR (texto de Messages) | x=98, y=225, 80 × 60; código emitido VK_BACK_SPACE |

As 14 teclas e espaço restante até altura 350 são referência de construtor, não uma grade 4 × 3. O layout ascendente anterior de 12 teclas foi substituído na revisão 1.0.0.1. O teste de prévia demonstra LIMPAR como Backspace no cursor/seleção do campo; equivalência com todos os alvos Java ainda exige evidência. A legenda não significa apagar tudo. Ação clear-all do editor, se mantida, é contrato local explícito distinto, não uma suposta equivalência Java.

A ordem é fundo da tela → decoração do componente → formas e textos → indicadores de edição. Textos permanecem editáveis. O fundo não recebe seleção ou arraste de componente.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| tecla normal | skin botao_released.png | recurso decorativo opcional |
| tecla auxiliar | skin botao_aux1_released.png | recurso decorativo opcional |
| texto numérico histórico | #FFFFFF, 24, negrito declarado; face/efeito Java a comprovar | aparência/face e overrides |
| texto auxiliar | #FFFFFF, 12 | token específico do renderer |

Conteúdo e opções configuráveis: Campo de destino, posição, tamanho, fonte, cor e skins das teclas normais e auxiliares.

## Estados

Sem alvo; alvo válido; alvo removido; tecla pressionada; foco; selecionado no editor.

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

- Alvo fora da tela, alvo removido, campo vazio ao apagar, Unicode, cursor no início/fim e TAB sem próximo campo. No PPT as teclas são formas; a digitação vem do roteiro.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate ausente bloqueia fidelidade; fallback exploratório explicitado não homologa. Conferir face/métricas e rótulos com referência independente.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/catalog/model.go](../../apps/gerador-jornada/internal/catalog/model.go)
- [internal/presentation/scene.go](../../apps/gerador-jornada/internal/presentation/scene.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Arrastar o teclado, redimensioná-lo e escolher um campo de destino da mesma tela.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## ID de arquivo compartilhado

Aplicar o [design do campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) na instância. O atributo não é texto desenhado no componente e não altera sua tipografia ou geometria. Todas as cores, textos e legendas continuam independentes do nome-base ASCII.

## Paridade pós-v1 — 2026-10-19

Usar skins botao_teclado/botao_teclado_auxiliar_1 e seus estados originais por evidência, sem tratar imagem aproximada como homologada. Fonte real e flags XML são registros distintos: builder calcula bold/italic sem aplicar style ao Font. Teclas são partes derivadas, não componentes persistidos individuais; hit regions usam a mesma cena/métricas no desktop, com foco/cursor/TAB isolados da edição. PPT mantém formas/textos editáveis e roteiro demonstra digitação, não teclado executável. [Matriz C15–C17](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19), [plano](plan.md) e [tarefas](tasks.md). Numérico não homologa alfanumérico/calculadora por reuso de renderer.
