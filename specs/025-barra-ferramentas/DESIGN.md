# DESIGN — Barra de ferramentas — Header

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Header/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir acessar arquivos, alinhamento, prévia e exportação sem abandonar o desenho. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Ações com altura mínima 36 e ícones acompanhados de texto nas ações principais; cabeçalho não deve ocupar a área útil com botões esticados.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| superfície | #F8FAFC | tema do cabeçalho |
| título | #0F172A, 24, negrito | tema do editor |
| ação | #E8EEF6 / #0F172A | tema dos botões |
| destaque | #2563EB | marca e foco |

Conteúdo e opções configuráveis: Comandos de edição e arquivo; habilitação e estado ocupado transitórios.

## Estados

Pronto; sem seleção; alterações não salvas; exportação ativa; erro recuperável.

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

- Janela estreita, ação sem projeto, salvamento cancelado e duas exportações concorrentes. Importação .page continua fora do escopo.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/editor.go](../../apps/gerador-jornada/internal/adapter/desktop/editor.go)
- [internal/adapter/desktop/operations.go](../../apps/gerador-jornada/internal/adapter/desktop/operations.go)
- [internal/adapter/desktop/files.go](../../apps/gerador-jornada/internal/adapter/desktop/files.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Disponibilizar novo, abrir, salvar, prévia e as duas saídas com nomes claros.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Manter Java/EFL visível e desabilitado temporariamente, sem abrir diálogo de exportação; preservar código e contratos de geração para retomada posterior.
- **FR-009**: Oferecer atalhos de seleção e movimento previstos em 020-elemento-edicao, respeitando foco de entrada de texto e histórico.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.

## Identidade do programa

FR-010 usa o ícone definido no [design integrado](../001-gerador-jornada/DESIGN.md#ícone-do-programa--2026-10-18), com o mesmo desenho no EXE e no Fyne. Não converter a marca em botão de ação. A identificação visual da janela deve permanecer disponível offline e sem arquivo externo ao EXE.

## Excluir componente

Ordem no grupo do canvas: zoom → Ajustar → Excluir → espaço flexível → Fundo. Excluir tem texto e ícone de lixeira e usa o estado desabilitado do tema sem seleção válida. O botão permanece visível, inclusive desabilitado. Confirmação compartilhada com menu e Delete: título Excluir componente, quantidade selecionada, respostas Sim/Não. Conferir largura da barra em 1366 × 768 e DPI 100/150/200%.

## Versão da distribuição

Exibir `JORNADA STUDIO · 1.0.0.0` no cabeçalho e `Gerador de Jornada 1.0.0.0` no título da janela. A identificação discreta mantém os estilos tipográficos existentes. A versão vem de `internal/adapter/buildinfo/VERSION`, também utilizada pelo build dos metadados Windows.

## Prévia da tela ativa

Acrescentar o botão com ícone de olho e legenda Prévia Tela Atual imediatamente após Prévia. O primeiro continua abrindo o projeto completo. A nova ação fica visível e desabilitada sem tela ativa; a disponibilidade não depende da seleção de componente. Conferir a largura do cabeçalho na resolução nominal e DPIs homologados.
