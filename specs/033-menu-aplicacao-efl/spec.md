# Feature Specification: Menu aplicação — árvore editável

**Feature Branch**: `033-menu-aplicacao-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-23
**Status**: Implementação parcial na revisão 1.0.0.3; requisitos remanescentes e homologação/aceite pendentes
**Input**: Adicionar Menu aplicação em Itens legados, junto aos componentes 031/032; editar árvore de serviços e jogos e simular sua navegação local. Usuário autorizou implementar os três componentes e fará o teste manual depois.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md#menu-aplicacao-2026-10-23); instâncias em [004-tela](../004-tela/spec.md#menu-aplicacao-2026-10-23); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md#menu-aplicacao-2026-10-23).
**Referência**: fonte integral recuperada e preservada em `material-legado/menu/menu.xml`, skins em `material-legado/elemetos_componente/xml/skin.xml` e classes Java relevantes. Os anexos de imagem estão realmente vazios. [DESIGN.md](DESIGN.md) distingue geometria derivada do código de fidelidade visual ainda não verificada.

<a id="estado-verificado-1-0-0-3"></a>
## Estado vigente — 1.0.0.3

Conferência adicional com `material-legado/menu/PnlMenu.java` e `UIMenuControle.java`: preservar ausência de tituloPrincipal; aplicar faixa horizontal do jogo em y150/27px; centralizar os dois botões de Mega-Sena em x220/400, y310, 164×52; compartilhar pintura e clique. Rodapé local com Cancelar/Voltar no submenu e Sair na raiz, com as skins do XML e sem widgets opacos externos. Grade padrão usa a fórmula Java; layouts financeiro/digital e exceções visuais seguem a referência preservada. Casos especiais por nome servem somente à apresentação, nunca à resolução de movimentos ou rotas.

Propriedades passam a editar código, nome, opção e atalho já existentes no formato. A seleção simulada identifica a folha e seus metadados, sem fundir normal/especial que compartilham código. Voltar/Cancelar limpam a seleção; componente desabilitado ou sessão encerrada não ativa itens. Novas instâncias não inventam título; projetos existentes não são migrados.

**Limite aberto:** movimento.xml/movimentos.xml e a implementação de controle/InfoMovimento não estão acessíveis. O cadastro não foi validado e não se afirma configuração completa. Java disponível demonstra lookup por nomeMovimento quando presente; nomes/códigos permanecem referências inertes no app. Títulos de folhas ainda não possuem campo próprio no modelo. Validações e revisão executável na [seção 25](../001-gerador-jornada/quickstart.md#25-submenus-java-e-movimentos--1003); Office/DPI/aceite continuam pendentes.

<a id="estado-verificado-1-0-0-1"></a>
<a id="estado-verificado-1-0-0-2"></a>
## Estado verificado — 2026-10-24 / 1.0.0.2

A fonte completa foi recuperada de `D:/desenvolvimento/ia/estudo/sistema-prototipo-android/componentes/src/java/menu/menu.xml` e preservada sem transcodificação em [material-legado/menu/menu.xml](../../material-legado/menu/menu.xml). O importador de desenvolvimento gera `menu_legacy_generated.go`; a fábrica instancia novos UUIDs e conserva a ordem, os nomes, espaços e metadados ativos. Isso substitui os relatórios sintéticos e a transcrição parcial descritos na 1.0.0.1.

A raiz permanece **12 serviços + 10 jogos**, mas a lista correta é: MEGA-SENA, LOTOFÁCIL, QUINA, LOTOMANIA, LOTECA, DUPLA-SENA, TIMEMANIA, DIA DE SORTE, SUPER SETE, MILIONÁRIA. LotoGol está comentado; Federal não é jogo root. Quatro jogos possuem ramos normal/especial. BOLÃO tem dez filhos, com BOLÕES DIGITAIS (`803`) por último: VENDER COTAS (`791`), BOLÕES MARKETPLACE (`790`) e BOLÕES CONCLUÍDOS (`792`). RELATÓRIOS tem 27 títulos descendentes e 148 folhas reais.

Cada item de jogo recebe sua chave de skin, com imagens default/touched/disabled do XML. A mesma cena serve canvas, prévia e PPT; pressão é transitória. A marca já contida no PNG não ganha texto sobreposto; legenda personalizada continua editável em caixa separada. Escolhas explícitas prevalecem, e skins indisponíveis geram diagnóstico. Arquivos existentes conservam sua árvore: a correção do seed vale para nova instância.

**Limites:** `tituloPrincipal` de ramos alimenta o título do submenu; o modelo atual não possui título separado para folhas. Legendas ilustrativas do layout digital são complemento da referência visual existente, não atributos do XML. Integridade global completa, recorte/foco/DPI, Office e aceite humano continuam pendentes. Evidências, hashes e roteiro na [seção 24 do quickstart](../001-gerador-jornada/quickstart.md#24-menu-xml-e-skins-dos-jogos--1002).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Inserir e editar o menu (Priority: P1)

Como autor, quero inserir um menu único e editar seus itens sem reconstruir telas para cada submenu.

**Why this priority**: A árvore configurável é o conteúdo principal da capacidade.
**Independent Test**: Em uma tela vazia, inserir Menu aplicação, editar dois ramos, salvar/reabrir e desfazer uma alteração, sem executar operações de negócio.

**Acceptance Scenarios**:

1. **Given** Itens legados, **When** inserir Menu aplicação, **Then** nasce uma única instância com árvore inicial de serviços à esquerda e jogos à direita baseada no Java/XML fornecido; Código de barras e Código de barras processado ficam no mesmo grupo.
2. **Given** um item em qualquer nível, **When** editar legenda, ordem, categoria, estilo, habilitação ou destino explícito, adicionar/remover item ou submenu e salvar, **Then** a alteração inteira é confirmada de uma vez, sem mudar os demais itens; Descartar e Undo/Redo preservam o último estado confirmado.
3. **Given** nomes repetidos, Unicode e menus vazios, **When** salvar/reabrir ou duplicar o componente, **Then** conteúdo e ordem permanecem; a cópia é independente e não reutiliza identidades internas da original.
4. **Given** uma árvore com Bolão aninhado, **When** selecionar seu layout especial, **Then** sua composição mantém três botões e legendas; outras árvores continuam editáveis, não viram telas adicionais.

### User Story 2 - Explorar serviços e submenus localmente (Priority: P1)

Como usuário da prévia, quero abrir submenus, voltar e selecionar um serviço sem acionar o sistema real.

**Why this priority**: Diferencia navegação local de uma rota de jornada ou operação financeira.
**Independent Test**: Usar raiz → submenu → submenu aninhado, Voltar, Cancelar, uma folha sem destino e Sair; repetir com duas instâncias.

**Acceptance Scenarios**:

1. **Given** a raiz visível, **When** ativar item habilitado com submenu, **Then** a mesma interface mostra apenas o submenu selecionado e mantém o caminho local; item desabilitado não produz efeito.
2. **Given** um submenu aninhado, **When** escolher Voltar, **Then** um nível é removido do caminho; na raiz, Voltar não muda nada. Cancelar retorna diretamente à raiz, inclusive quando o menu inicial configurado era outro.
3. **Given** uma folha sem tela de destino, **When** selecioná-la, **Then** a prévia informa a seleção simulada sem executar movimento, consultar conta ou inferir rota pelo texto/metadado.
4. **Given** uma folha com destino explicitamente escolhido pelo autor, **When** ativá-la na prévia completa, **Then** somente a tela existente configurada pode ser aberta; na Prévia Tela Atual a navegação entre telas permanece desativada.
5. **Given** a prévia aberta, **When** escolher Sair, **Then** encerra apenas a sessão de prévia, não o editor nem uma sessão de autenticação. Reabrir/revisitar restaura o menu inicial configurado, sem persistir o caminho percorrido.

### User Story 3 - Fixar o estado apresentado e exportado (Priority: P2)

Como autor, quero escolher qual menu aparece inicialmente no canvas, na prévia e na apresentação estática.

**Why this priority**: Permite exportar um submenu sem gravar a exploração de uma sessão.
**Independent Test**: Selecionar um submenu inicial, exportar sua representação e compará-la ao canvas antes e depois de explorar outro ramo na prévia.

**Acceptance Scenarios**:

1. **Given** um submenu existente escolhido como inicial, **When** salvar/reabrir ou gerar quadro estático, **Then** esse submenu, sua ordem e suas legendas são exibidos; o caminho da prévia livre não interfere.
2. **Given** remoção do submenu inicial sem escolha de substituto válido, **When** confirmar, **Then** o erro é localizado e nenhuma alteração parcial é publicada.
3. **Given** quadro válido, **When** exportar PPT, **Then** textos/formas continuam editáveis e representam o estado escolhido, sem prometer árvore interativa ou operação bancária no PowerPoint.

### Edge Cases

- Menus vazios, todas as folhas desabilitadas, rótulos repetidos/longos/acentuados, reordenação e remoção de ramo com descendentes.
- Profundidade máxima de 16 níveis contando a raiz como 1; no máximo 1000 itens por árvore, contando todos os descendentes. Limites 17/1001 são recusados antes de publicar ou recursar sem limite.
- Identidades repetidas, ciclos em memória, destino inexistente e item contendo simultaneamente submenu e destino são inválidos. Autor pode manter rascunho incompleto conforme regras gerais, sem exportá-lo como válido.
- Instância ou ancestral desabilitado impede ativar descendentes; nenhuma busca de rota por nome, atalho ou código de movimento.
- Janela 1366 × 768 e DPI 100/150/200%; rolagem/recorte não tornam itens sobrepostos nem alteram sua identidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Oferecer Menu aplicação, Código de barras e Código de barras processado em Itens legados somente após suporte técnico completo dos consumidores correspondentes.
- **FR-002**: Manter um único componente com dois estados de visualização local: raiz ou submenu selecionado; suportar aninhamento sem criar uma tela de jornada por submenu.
- **FR-003**: Permitir editar a árvore, suas legendas, ordem, adição/remoção, categorias, estilos e habilitação por item, bem como escolher explicitamente o destino de uma folha; aplicar confirmação atômica e histórico existente.
- **FR-004**: Usar raiz serviços/esquerda e jogos/direita, grade padrão, grade financeira com botões de 140 × 80 unidades lógicas e composição Bolão digital de três botões com legendas, conforme referência textual, sem declarar validação de imagem ausente.
- **FR-005**: Fazer Voltar remover um nível local, Cancelar restaurar raiz e Sair encerrar somente a prévia; revisita reinicia no menu inicial configurado e não muda o documento.
- **FR-006**: Executar navegação entre telas apenas para folha habilitada com destino explicitamente autorado e válido; sem destino, informar item selecionado em simulação. Prévia Tela Atual não segue destinos.
- **FR-007**: Preservar árvore e escolha inicial no projeto; manter caminho, seleção simulada, foco e estado de saída somente durante a sessão, isolados por instância/visita.
- **FR-008**: Permitir menus vazios em autoria; exigir identidades únicas, até 16 níveis e 1000 itens por árvore, e recusar submenu misturado com destino ou referência inicial/destino inválida sem perda parcial.
- **FR-009**: Renderizar canvas, prévia, quadro e PPT pelo mesmo conteúdo confirmado; a saída estática usa o menu inicial escolhido e mantém textos/formas editáveis, sem ler widgets ou capturar tela.
- **FR-011**: Derivar novas instâncias da árvore ativa do XML preservado, respeitando comentários, ordem, espaços e metadados; conferir cada nó com teste independente contra a fonte. Não reconstruir árvore a partir da lista de assets nem migrar documentos existentes silenciosamente.
- **FR-012**: Vincular jogos normais/especiais à chave original e aos estados default/touched/disabled definidos em skin.xml. Preservar a marca do PNG sem texto duplicado, legenda autorada editável, precedência de variantes e bloqueio de skin indisponível nos três consumidores.
- **FR-010**: Não executar Java/EFL, timers de negócio, periféricos, contas, autenticação, rede, QR ou despacho de movimentos. Metadados legados são descrição, nunca comandos executáveis; digitação de roteiro não é suportada neste composto.

### Key Entities *(include if feature involves data)*

- **Menu aplicação**: componente único com raiz e seleção inicial persistidas.
- **Submenu**: identidade, nome, título, organização visual e itens ordenados; pode aparecer aninhado em um item.
- **Item de menu**: identidade, legenda, categoria, estilo/habilitação, submenu opcional ou destino explícito de folha; metadados legados sem execução.
- **Caminho local**: percurso de uma sessão, distinto da árvore e das rotas de Jornada. Modelo técnico canônico em [003/data-model.md](../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os cenários de US1 conservam integralmente texto, ordem e habilitação após confirmação, Undo/Redo e salvar/reabrir; nenhuma identidade da duplicação colide.
- **SC-002**: Todos os percursos de US2 produzem o nível/tela esperado, com zero alterações no projeto e zero operações externas; duas instâncias permanecem independentes.
- **SC-003**: Árvores nos limites 16/1000 são preservadas; violações 17/1001, identidades repetidas e referências inválidas são localizadas sem publicação parcial.
- **SC-004**: Canvas, quadro e PPT mostram o mesmo menu inicial, conteúdo e ordem; todos os textos exportados permanecem editáveis. Comparação visual com imagem só será afirmada após recebê-la.

## Assumptions

### Histórico: correção após teste parcial — 2026-10-23

A transcrição abaixo foi superada pela fonte recuperada e pelo estado 1.0.0.2 acima; nomes abreviados e lista antiga não são mais a fábrica vigente.

- Usuário aprovou quase, mas faltava **PAGA PREMIADOS**; pediu RELATÓRIOS abaixo de APOIO, MEIOS abaixo de RELATÓRIOS e OFERTAR abaixo de PAGA. Não é aceite integral.
- XML recebido na conversa confirma **12 serviços + 10 jogos root**. Ordem: boleto/convênios/benefícios/financeiros/atendimento/resultados/apoio/bolão/relatórios/paga premiados/meios de pagamento/ofertar. **BOLÕES DIGITAIS (`803`) é último filho de BOLÃO**, não serviço root.
- Duas colunas alternadas: esquerda **APOIO→RELATÓRIOS→MEIOS**; direita **BOLÃO→PAGA PREMIADOS→OFERTAR**. Preservar x=20/200, topo150, pitch73, botões164×68; clip serviços largura400 até y600, última linha y515–583. Rodapé Sair reservado somente para jogos à direita.
- PAGA PREMIADOS: PROGNÓSTICOS `1500`, LOTERIA FEDERAL `1510`, INSTANTÂNEA `1513`. MEIOS: PIX → RECEBIMENTO/CONSULTAR RECEBIMENTO opção `99`, CARTÃO CRÉDITO/DÉBITO `99`, ENTRADA EM DINHEIRO `2`. Preservar literais/espaços do XML.
- RELATÓRIOS integral do XML: **148 folhas/27 títulos**, zero inicial em `080206`, dias `1..8`, auditoria `9..16`, pagamento `[0,2,3,4,5,6,7,1]`. Metadados permanecem texto inerte.
- Factory corrigida somente para **novas instâncias**: inserir novo **Menu da aplicação** para receber árvore corrigida. Não reescrever menus editados/salvos ao abrir/salvar/renderizar. Clipping corrigido afeta existentes, mas preserva integralmente a árvore autorada.
- Código/factory/fixtures e clipping/regressões corrigidos; principal ainda validando, sem novo build ou PASS Go afirmado. Hash anterior não é atual. Rastreabilidade [003/T022](../003-catalogo/tasks.md#t022-correcao-menu-xml) e [001/T145–T146](../001-gerador-jornada/tasks.md#t145-correcao-menu-xml). Não implementar outros ramos APOIO/BENEFÍCIOS/RESULTADOS agora.

### Premissas preservadas

- Java/XML do usuário é fonte do seed de serviços e jogos e da geometria, não runtime a importar ou transcrever integralmente. Não duplicar o XML enorme nos documentos; registrar diferenças reais quando a fonte for conferida pelo executor.
- Sem imagem recebida, não há homologação visual. Se a transcrição da árvore não estiver acessível ao executor, solicitar a fonte antes de alegar seed fiel; não substituir por menu genérico e chamá-lo de legado.
- A autorização atual permite implementar 031/032/033 e preparar revisão executável; não encerra revisões/aceites históricos abertos de 001–006, não altera checklists de revisores e não comprova teste manual.
- Extensão de clique demonstrativo com identificação do item só será habilitada após decisão explícita da coordenação e contrato sincronizado; até lá, clique é visual e a exportação é estática. Nenhum novo tipo de etapa é presumido.
- Esta pasta contém contrato de componente, não outro domínio ou backlog. Plano/modelo/tarefas permanecem canônicos em 001–006; ponteiro ativo 031 é preservado.