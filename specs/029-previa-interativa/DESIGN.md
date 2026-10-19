# DESIGN — Prévia interativa — Preview

**Status**: contrato visual de referência; personalizações e homologação pendentes.
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
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
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

A superfície com os controles da jornada usa escopo de tema claro fixo, preservando a aparência anterior dos widgets interativos. Imagens e primitivas continuam usando os estilos do modelo. O cabeçalho explicativo da janela acompanha o tema do editor. Conferir a prévia aberta ao alternar os três temas, incluindo botões, campo, seleção e teclado numérico.
