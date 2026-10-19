# DESIGN — Código e diagnóstico — Code

**Status**: contrato visual de referência; personalizações e homologação pendentes.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Code/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir inspecionar os fontes por tela e entender impedimentos antes da geração Java. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Área de código com rolagem horizontal/vertical, árvore de arquivos e fonte monoespaçada 13; diagnóstico deve ser legível sem abrir arquivo.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| fundo do código | #161A20 | tema do visualizador |
| texto | #E2E8F0 | tema do visualizador |
| fonte | monoespaçada, 13 | tema do visualizador |
| diagnóstico | #92400E sobre #FFFBEB | painel de diagnóstico |

Conteúdo e opções configuráveis: Arquivos gerados e diagnósticos derivados do snapshot; seleção de arquivo e rolagem transitórias.

## Estados

Perfil indisponível; validação recusada; fontes disponíveis; arquivo selecionado; exportação cancelada.

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

- Sem JARs EFL, nomes em colisão, pacote inválido, arquivo extenso e código Unicode. O legado gerava templates inconsistentes; essa falha não deve ser reproduzida.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fonte ausente usa fallback documentado; o conteúdo e a geometria não podem mudar silenciosamente. Verificar acentos, emoji e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/javaefl/sources.go](../../apps/gerador-jornada/internal/adapter/javaefl/sources.go)
- [internal/adapter/javaefl/profile.go](../../apps/gerador-jornada/internal/adapter/javaefl/profile.go)
- [internal/adapter/desktop/export.go](../../apps/gerador-jornada/internal/adapter/desktop/export.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Apresentar arquivos organizados por tela, fábrica, navegação e integração quando houver perfil homologado.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Enquanto a capacidade Java/EFL estiver suspensa por decisão do usuário, sua ação permanece visível e desabilitada. Não remover o gerador nem desbloquear perfil sem homologação.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.
