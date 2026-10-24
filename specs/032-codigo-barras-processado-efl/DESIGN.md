# DESIGN — Código de barras processado — blocos parametrizáveis

## Aplicação na revisão 1.0.0.1

Ficha composta 530 × 120 com blocos/linhas identificados, cinco rótulos permitidos e valores literais; inspector administra rascunho, tipos e ordem. Salvar/reabrir, cópia e PPT usam as propriedades confirmadas.

**Revisão ainda necessária:** A imagem de referência continua ausente; posições/cores do renderer são proposta documentada, não reprodução homologada do anexo. Diagnóstico integral de overflow e matriz completa de identidade/edição permanecem pendentes. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [003-catalogo](../003-catalogo/spec.md#codigos-barras-2026-10-23)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).

## Responsabilidade visual

Apresentar blocos de informação com rótulos e valores configurados pelo autor. “Processado” é a representação de resultado em uma jornada ilustrativa; não exibir indicador de aprovação financeira gerado automaticamente por um DV correto.

O segundo anexo não chegou. Não há evidência para reproduzir seus tamanhos, títulos de seção, bordas, cores ou quantidade de blocos. O layout abaixo é uma **proposta inicial explícita**, ajustável quando a referência for recebida, não uma cópia homologada.

## Geometria e composição

Proposta de nova inserção: **530 × 120 unidades lógicas**, um bloco sem título e duas linhas. Margens internas de 12, espaçamento de 12 entre blocos, altura-base de linha 24 e intervalo rótulo/valor de 8. A região de rótulo usa 40% da largura interna; a região de valor recebe o restante após o intervalo. Blocos empilham verticalmente na ordem configurada; linhas mantêm sua ordem interna. Um título não vazio reserva uma linha-base de 24 antes das linhas; título vazio não reserva essa linha.

Valores longos podem quebrar linha dentro de sua região; a altura cresce pelo número de linhas medido pela cena comum. Se o conteúdo ultrapassar os bounds confirmados, produzir diagnóstico de overflow para ajuste pelo autor, sem aumentar bounds ou cortar dados persistidos silenciosamente. A proposta de default não prova métricas da imagem; o gate de fidelidade visual continua aberto.

Rótulos alinhados à esquerda e valores à esquerda por padrão; override de alinhamento atua na região de cada texto, nunca desloca o bloco inteiro. Não desenhar borda/título financeiro, ícone de aprovação ou fundo de Report por suposta semelhança com anexo não visto.

O composto é uma única instância de Tela; blocos/linhas são dados internos ordenados. No PPT, cada título, rótulo e valor permanece texto editável dentro do agrupamento correspondente, nunca captura rasterizada.

## Tokens e pontos de personalização

| Parte | Proposta inicial / edição |
| --- | --- |
| Fonte e tamanho | defaults comuns de texto, 12 unidades; face real resolvida pela biblioteca/perfil, com diagnóstico se ausente |
| Rótulo e valor | proposta atual: rótulo #334155, valor #0F172A e título #1E3A8A; fonte física Futura Medium, sem atribuir esses tokens ao anexo ausente |
| Fundo | transparente, sem caixa bancária obrigatória |
| Título de bloco | texto livre, opcional; vazio não gera texto substituto |
| Conteúdo | lista ordenada de blocos e linhas configuráveis |
| Tipo do rótulo | enum de cinco escolhas; não string duplicada do texto exibido |
| Valor | texto literal configurado; sem máscara fiscal automática |
| Bounds/estilo | propriedades comuns da instância, sem coordenadas absolutas por linha gravadas no domínio |

Não serializar cores/medidas propostas em dezenas de campos se as propriedades comuns/defaults bastarem. O modelo canônico está no [Catálogo](../003-catalogo/data-model.md#modelo-codigos-barras-2026-10-23); implementação de layout pertence à apresentação.

## Estados

| Situação | Resultado |
| --- | --- |
| Nova instância | Um bloco, CNPJ e Nome, valores vazios; sem nome/documento inventado |
| Configurada | Títulos/linhas/valores na ordem confirmada |
| Rascunho no inspector | Alterações ainda não confirmadas não modificam cena, prévia nem documento |
| Combo aberto | Popup do editor, com opção corrente e navegação por teclado; não parte da tela da jornada |
| Valor vazio | Manter rótulo e área de valor em branco, sem placeholder numérico |
| Erro de estrutura | Localizar bloco/linha e conservar configuração confirmada |
| Overflow | Diagnóstico textual ao autor; preservar conteúdo integral |
| Seleção no canvas | Contorno/alças transitórios no composto; não exportar |
| Prévia | Somente leitura dos valores confirmados; nenhum combo de parametrização ou processamento externo |

## Parametrização no painel de propriedades

Lista de blocos com seleção local, título e comandos de adicionar/remover/reordenar. Dentro do bloco, lista de linhas com seletor de rótulo e editor de valor. Não identificar blocos pelo título nem linhas pelo texto: usar IDs internos estáveis, permitindo títulos repetidos e linhas com o mesmo tipo.

| Opção do combo | Texto na jornada |
| --- | --- |
| CPF | CPF: |
| CNPJ | CNPJ: |
| Nome | Nome: |
| Nome fantasia | Nome fantasia: |
| Razão social | Razão social: |

Troca de opção conserva valor literal. Uma confirmação aplica o rascunho completo de forma atômica; cancelar não altera o conteúdo anterior. O autor pode corrigir o valor depois de trocar o significado do rótulo; a UI não procura documento/nome nem converte CPF em CNPJ.

## Personalização segura

1. Configurar somente o bloco/linha selecionado e confirmar explicitamente.
2. Usar uma transação reversível e cópia profunda, sem compartilhar listas entre original/cópia.
3. Manter popup, seleção do inspector e rolagem fora do projeto.
4. Usar a mesma cena e medição de texto no canvas, prévia e PPT.
5. Rever esta proposta com o anexo quando recebido; registrar a nova referência e retestar o layout, sem chamar a proposta atual de paridade comprovada.

## Limites e acessibilidade

- Seletores identificam bloco e linha; foco é visível e sua ordem acompanha a lista exibida.
- Mensagens de erro e overflow não dependem somente de cor.
- Testar textos longos e acentuados; preservar vírgulas, pontuação e zeros como dados, não separadores.
- CPF/CNPJ são dados ilustrativos parametrizados, não fiscalmente validados. Não registrar valores em logs.
- Nenhum controle de autoria é desenhado dentro do slide ou passado como interação da prévia.

## Referências de implementação

[Plano coordenador](../001-gerador-jornada/plan.md#plano-codigos-barras-2026-10-23): Catálogo define enum/blocos/linhas; Tela compõe e duplica; Projeto serializa; Demonstração preserva quadro; apresentação e desktop desenham/editam. Reutilizar primitivas de texto, sem colocar validações de domínio no formulário ou forçar este tipo para Report.

## Roteiro de revisão visual

- Usar dois blocos com escolhas CPF/CNPJ e Nome/Nome fantasia/Razão social; trocar uma linha e comparar as outras.
- Confirmar/cancelar, reordenar, duplicar, desfazer/refazer e salvar/reabrir.
- Comparar editor/prévia/PPT, textos editáveis e ausência de combos de autoria na tela.
- Conferir valores vazios, títulos ausentes, textos longos e DPI 100%, 150% e 200%.
- Registrar resultado observado e revisão no [quickstart](../001-gerador-jornada/quickstart.md#17-codigos-barras-2026-10-23). Teste e aprovação visual permanecem pendentes, em especial pela ausência da imagem.

## ID de arquivo compartilhado

Usar o [campo ID](../028-propriedades/DESIGN.md#id-do-arquivo) da instância composta. UUIDs internos de blocos/linhas não são nomes visuais nem arquivos independentes.