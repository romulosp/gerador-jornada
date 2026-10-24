# DESIGN: Menu aplicação

## Aplicação na revisão 1.0.0.1

Menu 800 × 600 implementa árvore/inicial, raiz 12 serviços + 10 jogos, vistas de submenu, estilos de botão/faixa, comandos locais e destinos autorados. Contagem de RELATÓRIOS é 27 títulos/148 folhas na fixture atual.

**Revisão ainda necessária:** O XML integral de menus não está no workspace. menu_reports.go completa a contagem com títulos e tokens sintéticos: isso não cumpre a fidelidade literal exigida pela referência. Integridade global dos destinos, foco/recorte completo e homologação visual permanecem pendentes. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md). **Fonte**: Java/XML fornecidos pelo usuário na conversa; imagens anexadas realmente vazias. Não há screenshot a reconstruir por suposição.

## Identidade e composição

Uma instância selecionável/movível/redimensionável no canvas. Menus e itens internos não entram como componentes independentes na lista de Tela. O editor mantém sua decoração/tema; cores e fontes da jornada vêm do modelo/cena, não de widgets nativos opacos.

Dois estados visuais, não duas janelas: **raiz** e **submenu selecionado**. Aninhamento troca o conteúdo do segundo estado, sem empilhar painéis visíveis. Título vem do submenu selecionado; rótulo de item não substitui sua identidade.

**Correção solicitada em 2026-10-23:** a raiz não exibe faixa nem título. Somente submenus exibem a faixa principal, com o texto branco do token `label_faixa_titulo_principal`, não o azul de `label_titulo_1`. A distinção usa a identidade da raiz, inclusive quando um submenu é escolhido como menu inicial. Posições dos botões permanecem preservadas. Usuário também relatou um botão incorreto, mas os anexos chegaram vazios: botão específico e causa continuam sem identificação, aguardando captura ou nome; não atribuir correção ou aceite por suposição.

## Layouts contratados

| Layout | Composição | Evidência e limite |
| --- | --- | --- |
| `root` | Serviços à esquerda e jogos à direita, ordem relativa dos itens preservada dentro das categorias | Seed e organização derivados do Java/XML do usuário; não inventar catálogo completo a partir de nomes parciais. |
| `standard` | Grade de itens do submenu selecionado, título e comandos locais | Espaçamentos/colunas/bounds devem ser extraídos da referência Java, não de imagem inexistente. |
| `financial` | Grade financeira com cada botão 140 × 80 unidades lógicas na escala-base | Medida fornecida na referência textual, ainda sem inspeção visual; não cria domínio financeiro. |
| `bolao` | Bolão digital (`bolaoDigital` na origem): três botões e legendas da referência | Textos/variantes saem do seed; layout não limita globalmente a árvore a três itens. Edições fora da composição de três posições precisam diagnóstico visual, nunca poda silenciosa dos dados. |

Não fixar medidas externas, cores, ícones, número de colunas ou margens sem evidência. O executor deve obter os bounds/offsets no Java fornecido e registrar a proveniência no artefato existente de pesquisa antes de alegar paridade. Quando o conteúdo exceder bounds, aplicar recorte/rolagem contratados e diagnóstico apropriado; nunca descartar itens persistidos ou ampliar o componente sem decisão do autor. Geometria desconhecida é limitação de fidelidade, não prova de suporte ausente de todo o menu.

## Comandos e regiões de interação

### Correção da raiz — XML/teste parcial, 2026-10-23

Raiz com **12 serviços + 10 jogos**; ordem alternada dos serviços: boleto/convênios/benefícios/financeiros/atendimento/resultados/apoio/bolão/relatórios/paga premiados/meios/ofertar. Últimas linhas:

| y | x=20 | x=200 |
| --- | --- | --- |
| 369 | APOIO | BOLÃO |
| 442 | RELATÓRIOS | PAGA PREMIADOS |
| 515 | MEIOS DE PAGAMENTO | OFERTAR |

Preservar geometria-base: x20/200, topo150, pitch73, botões **164×68**. Clip serviços **largura400 até y600**; última linha termina em **583**. Rodapé **Sair somente para jogos à direita**, sem cortar serviços. Pintura e hit regions devem concordar com o recorte, sem deslocar/reduzir botões ou podar dados.

**BOLÕES DIGITAIS (`803`)** é último filho de BOLÃO, não serviço root. PAGA PREMIADOS: PROGNÓSTICOS `1500`, LOTERIA FEDERAL `1510`, INSTANTÂNEA `1513`. MEIOS: PIX → RECEBIMENTO/CONSULTAR RECEBIMENTO opção `99`, CARTÃO CRÉDITO/DÉBITO `99`, ENTRADA EM DINHEIRO `2`. RELATÓRIOS completo do XML: **148 folhas/27 títulos**, `080206`, dias `1..8`, auditoria `9..16`, pagamento `[0,2,3,4,5,6,7,1]`; preservar zeros/espaços/literais, sem executar movimentos.

Seed corrigido somente para **novas instâncias**; inserir novo **Menu da aplicação** para recebê-lo. Não reescrever árvores editadas/salvas ao carregar/salvar/renderizar. Clipping afeta também existentes sem destruir autoria. Não ampliar outros ramos APOIO/BENEFÍCIOS/RESULTADOS agora.

Usuário aprovou quase, com ausência de Paga e posições corrigidas acima; não é homologação integral. Código/regressões presentes, principal ainda validando; novo build e manual posteriores pendentes, sem PASS Go ou hash atual estabelecidos aqui.

- Botões internos têm hit regions identificadas por componente e item, não por legenda ou posição da lista.
- Habilitado/pressionado/desabilitado respeitam estado e estilo explícitos. Instância desabilitada bloqueia todas as regiões; ramo desabilitado não pode ser atravessado.
- **Voltar**: remove o último submenu do caminho; raiz é limite sem efeito.
- **Cancelar**: restaura a raiz, mesmo que a escolha inicial seja submenu.
- **Sair**: fecha somente a sessão de prévia; não encerra editor, projeto ou autenticação.
- Folha com destino autorado usa validação de Jornada; sem destino informa seleção simulada. Comandos de controle não são destinos inferidos por texto. Na Prévia Tela Atual, submenus operam e rotas entre telas permanecem desativadas.

## Autoria e estado inicial

Formulário em árvore, com seleção de submenu/item e edição de legenda, categoria, variante visual, habilitação, ordem, adição/remoção e destino explícito. Mostrar submenu e destino como alternativas mutuamente exclusivas. Menu e Submenu mantêm modelos separados, sem JSON livre nem centenas de botões independentes no canvas.

Escolha **Menu inicial** lista identidades existentes com caminho/nome para desambiguar títulos repetidos; raiz é o default. Essa escolha é configuração e define o estado fixo exportado. Seleção do nó no inspector, expansão de árvore e caminho explorado na prévia são transitórios. Salvar alterações/Descartar ficam no rodapé fixo existente; confirmar faz uma transação, erro mantém o rascunho e identifica o nó.

Menus vazios são válidos para autoria. Se o submenu inicial for removido, exigir escolha explícita de outro menu antes da confirmação; não trocar silenciosamente para raiz. Não permitir estruturas excedendo os limites de [spec.md](spec.md), mesmo em rascunhos de UI que serão confirmados.

## Cena e exportação

Canvas e PPT estático usam o menu inicial confirmado; prévia usa o caminho da sessão a partir desse mesmo valor. O caminho é inicializado pelos ancestrais do submenu escolhido para que Voltar tenha significado. Textos, formas, título, ordem e enabled vêm da mesma cena. Não ler estado de widgets, fotografar a prévia ou exportar o inspector.

PPT conserva textos/formas editáveis; não se promete interação da árvore, Cancelar/Sair, despacho ou autenticação dentro do arquivo. `typeText` não tem capacidade neste composto. Eventual clique de roteiro por item exige contrato aprovado de metadados/quadros e não é habilitado por esta redação.

## Revisão pendente

Revisar os quatro layouts, menus vazios/aninhados, rótulos longos/Unicode, múltiplas instâncias, recorte, hit-test e DPI 100/150/200%; comparar canvas/prévia/PPT com fonte Java/XML e, quando recebida, imagem real. Nenhum desses resultados foi observado nesta entrega documental. Roteiro/evidências ficam no [quickstart coordenador](../001-gerador-jornada/quickstart.md#18-menu-aplicacao-2026-10-23).