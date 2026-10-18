# DESIGN — Paleta — ComponentTab

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/ComponentTab/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir encontrar todas as famílias e conjuntos do legado para inserir na tela. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Rolagem própria; item clicável de pelo menos 36 unidades de altura; painel utilizável em janela 1366 × 768.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| item | #E8EEF6 | tema da paleta |
| legenda | #0F172A, 14 | tema da paleta |
| ícone | + | recurso vetorial do editor |
| seção | #64748B, 11 | tema do painel |

Conteúdo e opções configuráveis: Definições do catálogo e conjuntos disponíveis; não persistir itens de paleta dentro do projeto.

## Estados

Disponível; arrastando; drop aceito; drop recusado; sem tela.

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

- Paleta maior que a janela, rolagem antes do arraste, drop em painel lateral e nenhum projeto aberto.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/palette.go](../../apps/gerador-jornada/internal/adapter/desktop/palette.go)
- [internal/adapter/desktop/presets.go](../../apps/gerador-jornada/internal/adapter/desktop/presets.go)
- [internal/adapter/desktop/tools.go](../../apps/gerador-jornada/internal/adapter/desktop/tools.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Exibir as sete famílias e os conjuntos auxiliares em uma aba própria.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.
