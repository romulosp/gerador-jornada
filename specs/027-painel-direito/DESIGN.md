# DESIGN — Painel direito — ToolsRight

## Aplicação na revisão 1.0.0.1

Painel direito usa formulário em área com rolagem nos dois eixos e rodapé fixo Salvar alterações/Descartar. O divisor faz parte do layout Fyne; conteúdo longo não é um novo dado do projeto.

**Revisão ainda necessária:** A suíte histórica inspector_split_test.go não está na árvore atual. Há implementação do layout, mas falta repetir sua matriz específica de largura/rolagem/foco em janela nativa/DPI. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/ToolsRight/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir apresentar somente as propriedades pertinentes à seleção atual. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Largura ajustável; rótulos e campos legíveis sem cobrir o canvas; altura controlada por rolagem.

### Clarificação do divisor e da rolagem — 2026-10-21

O divisor deve mudar a largura **efetiva** do painel, não apenas seu `Offset`. O `Split` do Fyne limita o arraste pelo `MinSize` dos filhos: o antigo `container.NewVScroll(e.Inspector)` propagava a largura mínima do formulário, ampliada por seletores, caminhos de fontes e versões de perfil longos. Em [Editor.design](../../apps/gerador-jornada/internal/adapter/desktop/editor.go), o corpo passa a `container.NewScroll(e.Inspector)`, com ambos os eixos roláveis, desacoplando as duas dimensões mínimas do conteúdo. Valores e acesso horizontal permanecem íntegros.

A composição `right` com `container.NewBorder` permanece: cabeçalho no topo, `InspectorActions` na base e somente o corpo no centro rolável. Salvar alterações/Descartar não foram removidos nem incorporados à rolagem. O mínimo natural ainda é definido pelos controles fixos de cabeçalho/rodapé; não adotar largura fixa, reduzir botões ou exigir colapso ilimitado. Layout não publica nem apaga rascunho, documento ou histórico; Save/Discard permanecem iguais.

Revisão pendente: janela 1280 × 720 com tela, 16 famílias e seleção múltipla; observar deslocamento e largura reais ao mover o divisor, testar valores longos e alcançar os campos nos dois eixos. Cabeçalho/rodapé devem permanecer fixos e os botões renderizados devem receber os cliques. Quatro regressões em inspector_split_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_split_test.go`), não executadas por política Windows, sem contorno. **Evidência final do executor principal (2026-10-21):** gofmt em dois arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` concluiu a compilação de 15 suítes, **sem executá-las**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade PASS: CRC, todos os SHA256SUMS, sidecar SHA-256, PE amd64, ambas as cópias do EXE idênticas, README definitivo, manifesto com aprovações pendentes e ausência de backups terminados em `~`. Identidades finais e evidências no [quickstart, seção 11](../001-gerador-jornada/quickstart.md#11-redimensionamento-do-painel-de-propriedades--2026-10-21). [T129](../001-gerador-jornada/tasks.md#phase-23-redimensionamento-do-painel-de-propriedades--2026-10-21) permanece aberta como verificação mista: sem execução dos testes Go, validação de UI/arraste nativos, DPI, revisão manual ou aceite; `manualApproval` e `redistributionApproval` continuam `pending`.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| superfície | #F8FAFC | tema do painel |
| rótulo | #0F172A, 14 | tema dos formulários |
| campo | #FFFFFF / #CBD5E1 | tema de entrada |
| erro | mensagem textual junto ao contexto | diagnóstico do formulário |

Conteúdo e opções configuráveis: Rascunho de formulário e foco transitórios; valores confirmados pertencem à tela ou ao componente.

## Estados

Propriedades da tela; de componente; sem tela; rascunho inválido; confirmado.

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

- Trocar seleção com rascunho, componente removido, propriedade não suportada e campos fora da área visível.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/inspector.go](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go)
- [internal/adapter/desktop/editor.go](../../apps/gerador-jornada/internal/adapter/desktop/editor.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Mostrar configuração de tela quando nenhum controle estiver selecionado.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Reservar rodapé fixo para Salvar alterações e Descartar, ligado ao formulário corrente e fora da sua rolagem.
- **FR-009**: Distinguir propriedade de texto, geometria do elemento e comandos para seleção múltipla; trocar seleção não deve publicar rascunho incompleto.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.
