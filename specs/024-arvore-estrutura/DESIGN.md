# DESIGN — Estrutura — StructureTab

## Aplicação na revisão 1.0.0.1

Estrutura lista os componentes da tela ativa na ordem de camadas e permite selecionar objetos sobrepostos. Nomes/seleção são sincronizados com o canvas; ParentID existe no modelo.

**Revisão ainda necessária:** A lista atual é plana: não há árvore expansível de contêineres com arraste entre nós. Requisitos de novas famílias/contêineres continuam no backlog; não deduzir UI hierárquica da existência de ParentID. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [004-tela](../004-tela/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/StructureTab/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir selecionar e organizar componentes pelo nome, inclusive quando estiverem cobertos. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Lista vertical com rolagem; primeiro plano aparece antes do fundo; nomes longos precisam de identificação acessível.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| seleção | #2563EB / #FFFFFF | tema de lista |
| normal | #E8EEF6 / #0F172A | tema de lista |
| texto | 14 | tema do editor |

Conteúdo e opções configuráveis: Nomes e ordem lidos do projeto; seleção é transitória e não cria revisão de documento.

## Estados

Vazia; item ativo; item comum; componente coberto; tela trocada.

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

- Itens sobrepostos, componente removido por undo, nomes longos e mudança de tela com seleção anterior.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/tools.go](../../apps/gerador-jornada/internal/adapter/desktop/tools.go)
- [internal/adapter/desktop/operations.go](../../apps/gerador-jornada/internal/adapter/desktop/operations.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Listar somente os componentes da tela ativa por sua ordem de camadas.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.
