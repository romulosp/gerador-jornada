# DESIGN — Formulários de propriedades — Properties

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Properties/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir personalizar conteúdo, geometria e aparência com validação antes de aplicar. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Seções Conteúdo, Geometria, Aparência e Ações; rótulos em português, unidades explícitas e controles com área mínima 36.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #0F172A | tema de formulário |
| entrada | #FFFFFF | tema de campo |
| borda | #CBD5E1 | tema de campo |
| foco | #2563EB | tema de foco |

Conteúdo e opções configuráveis: Valores editáveis do modelo e rascunho de formulário separado; tipos de entrada apropriados a número, texto, opção e cor.

## Estados

Não alterado; alterado; inválido; aplicado; seleção trocada.

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

- Inteiro inválido, cor fora do formato, nome vazio, referências inválidas e texto com caracteres especiais. O Properties legado era um placeholder; o comportamento funcional estava em ToolsRight.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/inspector.go](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Manter botão Salvar alterações visível fora da rolagem do formulário, inclusive quando o último campo estiver fora da área visível.
- **FR-009**: Aplicar o rascunho válido em uma ação reversível e atualizar a tela imediatamente após salvar; Descártá-lo pelo botão Descartar restaura os valores confirmados.
- **FR-010**: Exibir Alinhamento do texto, com esquerda, centro e direita, em controles textuais. A alteração modifica texto dentro da caixa, sem mover o elemento.
- **FR-011**: Com seleção múltipla, mostrar quantidade e comandos coletivos; não exibir formulário de um único componente como se editasse todos.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.

## ID do arquivo

Adicionar ID logo após Tela/Nome, usando o mesmo campo de entrada dos dados comuns. Tela: dica `Pnl + nome, até 40 caracteres`, exemplo `PnlIdentificacao → PnlIdentificacao.java`. Componente: `até 40 letras/números, sem acentos; Pnl não obrigatório`. O ID técnico UUID não aparece como campo editável. Texto amigável continua podendo ter acentos. Limpeza de acentos/pontuação acontece na digitação; erro de tamanho/prefixo/colisão é explícito e não salva valor inválido. Nome e ID são independentes após a sugestão inicial.
