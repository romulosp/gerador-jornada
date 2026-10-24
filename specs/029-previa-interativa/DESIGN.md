# DESIGN — Prévia interativa — Preview

## Aplicação na revisão 1.0.0.1

Prévia completa parte da inicial e permite rotas; Prévia Tela Atual isola uma tela. Cena, decorações/fontes e regiões de teclado são compartilhadas; campos/opções/menu/código de barras mantêm runtime separado.

**Revisão ainda necessária:** Máscaras incrementais, todas as novas famílias, scroll/modal de todos os tipos e bloqueios integrais de perfis ainda não foram entregues. Os testes executados não homologam os critérios geométricos/visuais no Office ou em todos os DPIs. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Preview/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir percorrer a jornada e experimentar entradas antes de exportar. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Conteúdo lógico igual ao editor; janela acomoda a tela ou oferece rolagem; nenhuma alça de edição aparece.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| fundo | derivado da tela | Screen.Background |
| controles | derivados do componente | Style e propriedades |
| aviso de modo | #64748B, 13 | cabeçalho da prévia |

Conteúdo e opções configuráveis: Snapshot do projeto; tela corrente, foco, valores e escolhas runtime não são gravados no documento.

## Estados

Aberta; navegando; entrada alterada; tela revisitada; erro de referência; fechada.

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

- Ciclo de navegação, fechar durante callback, projeto inválido e ausência de imagem local. Prévia não executa Java nem regras do sistema real.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte/asset/delegate/referência ausente bloqueia fidelidade; fallback exploratório explícito não homologa. Conferir face real/glifos/métricas e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/application/preview.go](../../apps/gerador-jornada/internal/application/preview.go)
- [internal/adapter/desktop/preview.go](../../apps/gerador-jornada/internal/adapter/desktop/preview.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Abrir a tela inicial em uma janela independente com o mesmo fundo e controles.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Dois escopos de visualização

Manter Prévia como execução completa e acrescentar Prévia Tela Atual no cabeçalho. A janela individual usa título `Prévia Tela Atual · <nome da tela>` e indicação `Tela atual · campos e escolhas interativos; navegação entre telas desativada`. Preservar a renderização existente de fundo, tipografia, alinhamento e controles; nenhuma sobreposição de seleção ou alça do editor aparece. Dimensões seguem a tela com rolagem disponível. Não executar roteiro nem trocar de tela por botão nesse modo.

A janela usa snapshot do conteúdo confirmado, independente de futuras edições. Para visualizar uma propriedade ainda em rascunho, o autor primeiro usa Salvar alterações. Testar tela ativa diferente da inicial, campos, teclado, escolhas, fundo e fechamento sem alteração no projeto.


## Tema da superfície de execução

A base inicial usa escopo de tema claro fixo para isolar widgets interativos. No alvo pós-v1 de 2026-10-19, a superfície pinta controles pela mesma cena/skin/fontes do canvas/PPT; tema claro isolado não substitui fidelidade. O cabeçalho explicativo acompanha o tema do editor. Conferir a prévia aberta ao alternar temas, sem modificar o perfil visual dos componentes.

## Cena e interação comuns — expansão 2026-10-19

FR-012–FR-014: compartilhar fonte/face/RGBA/imagem/hash/baseline/spacement/alinhamento/fit/clip/estados por componente. Hit regions recebem clique/foco/tecla sem reinstalar pintura opaca nativa Fyne. Selected, pressed, focused e disabled são distintos conforme referência; XML bold/italic não prova estilo efetivo do Font criado pelo builder.

Estado local cobre entrada/máscaras/cursor/TAB/modificadores, grupos/seleções, tabelas/Reports, scroll e modal onde contratados; restaura iniciais por visita, sem gravar uso/alterar bounds. Rascunho não confirmado do inspector não entra na prévia. Recurso/delegate/capacidade faltante gera bloqueio por entrada, não visual aproximado chamado fiel. [Matriz](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19), 001/T111–T116 e T120–T122; Java continua pausado.

## Entrada e ficha processada — planejamento 2026-10-23

FR-015 usa [031/DESIGN](../031-codigo-barras-efl/DESIGN.md) para campos/estados/modal e [032/DESIGN](../032-codigo-barras-processado-efl/DESIGN.md) para a ficha de leitura. Modo de leitura deve ser perceptível; Enter submete uma captura, sem dividir a entrada antes de validá-la. Modal recebe foco, bloqueia entrada subjacente e tem a única ação “OK, entendi”; ao fechar retorna ao campo 1 vazio. Erro manual por bloco é local e identifica o par, sem limpar outros blocos válidos. Cor não substitui mensagem.

Concluir/ler/reiniciar são comandos locais da interação, não botões de Jornada gerados automaticamente. O inspector não é exportado nem reproduzido como combo dentro da ficha. Comparar foco/TAB/teclado numérico, captura parcial/cancelada e duas instâncias; estado não muda documento ou ficha parametrizada. Cena/áreas de interação compartilham transformação, sem widget opaco cobrindo a skin. Trabalho em 001/T135/T137; anexos, implementação e revisão visual continuam pendentes.

## Vista local de Menu aplicação — 2026-10-23

FR-016 consome [033/DESIGN](../033-menu-aplicacao-efl/DESIGN.md): raiz ou um submenu visível por instância, título/itens/hit regions da mesma cena. Aninhamento mantém caminho, não janelas/painéis sobrepostos; enabled efetivo bloqueia interação. Voltar pop, Cancelar raiz, Sair somente prévia; comandos identificados semanticamente, nunca inferidos de legenda/Movement*. Inicial vem do snapshot; revisita o restaura sem herdar caminho. Tela Atual permite árvore local mas não segue destinos de tela. Folha sem destino mostra seleção simulada, sem operação de negócio. Layout root/standard/financial140x80/bolaoDigital três botões+legendas deriva da fonte Java/XML; anexos vazios não validam visual. Sessão001/T140 e cena001/T142 são os donos; comparação DPI/PPT e aceite em001/T144 continuam pendentes.
