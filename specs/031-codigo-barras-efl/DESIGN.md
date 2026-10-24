# DESIGN — Código de barras — entrada validada

## Aplicação na revisão 1.0.0.1

Entrada composta 610 × 60 com oito campos derivados, validação matemática local, captura/colagem 44/47, comandos explícitos e mensagem modal. Oito valores e seus erros/foco ficam somente na sessão de prévia; PPT é estático/editável.

**Revisão ainda necessária:** Não há homologação bancária/SISPL, referência visual independente nem matriz completa de teclado/TAB/habilitação. A classificação CAIXA/SIGCB é premissa local e não consulta ou autorização de pagamento. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md#codigos-barras-2026-10-23)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).

## Responsabilidade visual

Representar uma única entrada composta, com rótulo “Número do Código de Barras:” e oito campos. Campos vazios não mostram números fictícios. Nenhum desenho de barras ou QR é exigido: o componente é a linha digitável da tela legada, não um gerador de simbologia.

## Geometria e composição

Origem relativa ao composto. Referência de nova inserção: **610 × 60 unidades lógicas**, extraída de `getPnlCodigoBarras()` fornecido pelo usuário. Rótulo em `(30, 2)`, área `200 × 20`. Campos em `y=20`, altura 20, com intervalo de 5 unidades:

| Campo visual | x | Largura | Limite | Papel |
| --- | --- | --- | --- | --- |
| 1 | 30 | 50 | 5 | primeira parte do bloco 1 |
| 2 | 85 | 50 | 5 | final e DV do bloco 1 |
| 3 | 140 | 50 | 5 | primeira parte do bloco 2 |
| 4 | 195 | 60 | 6 | final e DV do bloco 2 |
| 5 | 260 | 50 | 5 | primeira parte do bloco 3 |
| 6 | 315 | 60 | 6 | final e DV do bloco 3 |
| 7 | 380 | 15 | 1 | DV geral |
| 8 | 400 | 110 | 14 | fator e valor |

Três pontos separadores com área `5 × 10`, em `(81,27)`, `(191,27)` e `(311,27)`. Não acrescentar dígito, hífen ou pontuação ao conteúdo canônico.

Não criar Report, teclado, título de página, CONFIRMAR ou CANCELAR como filhos obrigatórios. São componentes independentes já existentes. Os oito campos e pontos são partes internas derivadas do composto, sem UUID/FileID de componentes adicionais.

Mover o composto preserva posições relativas. Redimensionamento muda sua representação proporcionalmente a essa referência, sem alterar limites de dígitos, DVs ou controles externos; overflow e legibilidade são diagnosticados, nunca resolvidos cortando dígitos. Aplicar o mesmo cálculo à pintura e às áreas de interação. Contêiner pai continua sujeito à regra de Tela: resize do pai não redimensiona seus filhos implicitamente.

## Tokens e pontos de personalização

[TextFieldSkin.java](../../material-legado/elemetos_componente/component/textfield/TextFieldSkin.java) confirma `TEXT_LOGIN = type1`. Estados e cores abaixo são dados de [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml#L217-L249), não evidência de execução Java:

| Parte / estado | Referência |
| --- | --- |
| Rótulo | `LABEL_TEXTO_AZUL`, RGB `#5064AA`, tamanho 12, fundo transparente; recurso declarado futura_mh.ttf, com família SFNT real Futura Medium em [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml#L892-L901) |
| Campo normal | `TEXT_LOGIN/type1`, fundo `#FFFFFF`, texto `#000000`, Arial declarada tamanho 12 |
| Campo com foco | fundo `#10508E`, texto `#FFFFFF`, mesma fonte/tamanho |
| Campo desabilitado | fundo `#E6EEFF`, texto `#808080`, mesma fonte/tamanho |
| Pontos | `LABEL_FORMULARIO` informado no Java; métricas/face efetivas exigem revisão própria |
| Erro | mensagem textual identifica o bloco; não inferir uma skin de erro legada que não foi fornecida |

Permitir rótulo, bounds e overrides comuns de texto/fundo/borda/fonte suportados. Não expor número de campos, limites, pesos dos módulos ou números de exemplo como tokens visuais. Alinhamento atua dentro de cada campo/região, não move o conjunto. Fonte declarada não prova licença, resolução ou baseline; não substituir face silenciosamente.

## Estados

| Situação | Resultado e interação |
| --- | --- |
| Editor | Campos vazios; selecionar/mover/redimensionar o composto, sem ativar digitação de boleto no canvas |
| Início | Somente campo 1 habilitado; foco visível |
| Par em andamento | Somente subcampo ativo recebe entrada; próximo par permanece bloqueado |
| Erro de bloco | Mensagem local “Bloco N inválido”; primeiro campo do par ativo, conteúdo do par disponível para substituição/correção; preservar blocos anteriores válidos |
| Erro geral digitado | “Código de barras inválido” local; não concluir; foco no campo 1, rascunho preservado para correção sequencial |
| Captura integral | Modo de leitura explícito com indicação visível e Enter para submeter; não distribuir caracteres em oito campos antes da validação |
| Captura inválida | Limpar entrada runtime e abrir modal com “Código de barras inválido” / “OK, entendi”; nenhuma tecla modifica campos sob o modal |
| Aceito | Oito campos preenchidos e somente para leitura; conclusão local liberada, sem navegar |
| Desabilitado | Sem entrada ou captura; visual de estado sem depender do tema do editor |
| Reinício/revisita | Oito campos vazios, sem resultados/erro/modo de captura herdados |

Concluir entrada, ler e reiniciar são comandos de interação da prévia, acessíveis por teclado/menu local; não criam botões persistidos nem são exportados como decoração do editor. O modal deve receber foco, ter botão como ação padrão e devolver foco ao campo 1 ao fechar. Não duplicar mensagem a cada caractere após o erro.

## Personalização segura

1. Editar apenas configuração persistível e confirmar uma transação do histórico.
2. Não gravar foco, valores de teste ou skin de erro transitória no documento.
3. Usar a mesma cena no canvas, prévia e PPT, sem Entry Fyne opaca escondendo a pintura.
4. Validar fonte, métricas, clipping e campos de 1 e 14 dígitos em cada DPI.
5. Preservar bounds/IDs/estilos existentes em futura migração; o default 610 × 60 só vale para novas inserções.

## Limites e acessibilidade

- Mensagem de erro deve ser textual e associada ao bloco; cor não é o único sinal.
- TAB/clique não atravessam um bloco ainda não validado; comandos de cancelar/fechar a prévia continuam acessíveis.
- Teclado numérico vinculado escreve no campo ativo do composto. A distribuição e validação são compartilhadas com o teclado físico.
- Campo de 14 posições deve manter todos os dígitos legíveis ou apresentar diagnóstico de overflow; não truncar o código para caber.
- A mensagem não expõe o código completo em log. Conteúdo dos cenários é ilustrativo.

## Referências de implementação

Trabalho **futuro**, com donos e caminhos no [plano coordenador](../001-gerador-jornada/plan.md#plano-codigos-barras-2026-10-23): regras puras em Catálogo, instância em Tela, wire no adaptador de Projeto, estado de sessão na aplicação, cena em apresentação e interação/exportação nos adaptadores. Não criar uma classe central que acumule essas responsabilidades.

## Roteiro de revisão visual

- Comparar oito campos/pontos e estados descritos; imagem do usuário continua ausente.
- Exercitar F01 válida e DVs locais/geral adulterados; capturas integrais inválidas devem abrir o modal especificado.
- Conferir 1366 × 768 e DPI 100%, 150% e 200%, teclado físico e teclado numérico.
- Salvar/reabrir sem valores runtime, duplicar sem compartilhamento e conferir PPT com textos/formas editáveis.
- Registrar revisão, executor, resultado observado e evidência no [quickstart](../001-gerador-jornada/quickstart.md#17-codigos-barras-2026-10-23); nenhum roteiro foi executado nesta etapa documental.

## ID de arquivo compartilhado

Aplicar o [campo ID de arquivo](../028-propriedades/DESIGN.md#id-do-arquivo) à instância composta, nunca aos dígitos. Seleção/alças não fazem parte do PPT. Java/EFL não é reintroduzido por esta referência visual.