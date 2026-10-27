# Feature Specification: Gerador de Jornada

**Feature Branch**: `001-gerador-jornada` (identificador da feature; não há branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Editor desktop offline para Windows, com componentes EFL arrastáveis, jornadas navegáveis, exportação PPTX editável e animada e geração de fontes Java 1.6/Swing com EFL.

**Referência vigente — 1.0.0.1:** a seção de estado verificado e a [seção 22 do quickstart](quickstart.md#22-reimplementação-visual-do-legado--2026-10-24) descrevem a árvore atual. FR-058 mantém o objetivo de filtrar dependências ausentes. Os inventários de 214/206 definições e seus PASS são registros históricos; a biblioteca estática corrente tem 205 chaves tipadas, 224 recursos e oito definições indisponíveis.

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Editor Go/Fyne com sete famílias base, três compostos, autoria de telas/jornada/roteiros, persistência 3/3, prévia e PPT. A revisão 1.0.0.1 aplica skins e fontes do legado ao recorte existente.

**Pendências e limites:** FR-049–FR-054 e a expansão de 16 famílias não estão integralmente atendidos. Faltam tipos adicionais, perfil visual portátil por documento, inventário completo de direitos/recursos e homologação Windows/DPI/Office. Java/EFL continua desabilitado.

**Fontes verificadas:** [document.go](../../apps/gerador-jornada/internal/application/document.go), [scene.go](../../apps/gerador-jornada/internal/presentation/scene.go), [legacy_preview_test.go](../../apps/gerador-jornada/internal/adapter/desktop/legacy_preview_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Desenhar e preservar um projeto (Priority: P1)

Como autor de jornadas, quero compor telas visualmente e salvar o projeto para continuar a edição sem perder trabalho.

**Why this priority**: É a base de todas as saídas e já permite validar a composição com o usuário.
**Independent Test**: Criar uma tela, inserir as sete famílias de componentes, editar, desfazer uma ação, salvar, fechar e reabrir com conteúdo e geometria preservados.

**Acceptance Scenarios**:

1. **Given** um projeto novo, **When** o autor cria uma tela, **Then** recebe uma área inicial de 800 × 600, nome editável e catálogo de componentes.
2. **Given** uma tela, **When** o autor arrasta um componente da paleta, **Then** ele é inserido na posição indicada e suas propriedades ficam disponíveis.
3. **Given** um componente selecionado, **When** o autor move, redimensiona, alinha ou duplica, **Then** a operação é aplicada e pode ser desfeita/refeita como uma única ação.
4. **Given** um projeto salvo, **When** é reaberto, **Then** telas, componentes, recursos e referências são preservados; seleção e zoom não alteram os dados.
5. **Given** erro ao salvar ou arquivo corrompido ao abrir, **When** a operação falha, **Then** a sessão ativa e a última versão válida permanecem disponíveis.

---

### User Story 2 - Conectar e experimentar telas (Priority: P1)

Como autor, quero definir a tela inicial e o destino dos botões para experimentar o fluxo antes de exportá-lo.

**Why this priority**: Uma coleção de desenhos precisa de navegação explícita para representar uma jornada.
**Independent Test**: Carregar um projeto de três telas, conectar seus botões e percorrer ida, retorno e encerramento da prévia sem exportação.

**Acceptance Scenarios**:

1. **Given** três telas, **When** o autor define a inicial e associa um botão a outra tela, **Then** a prévia abre a inicial e o clique leva ao destino.
2. **Given** uma tela referenciada, **When** o autor tenta removê-la, **Then** vê as referências afetadas e pode cancelar ou remover deixando diagnósticos localizados.
3. **Given** um destino inexistente, **When** solicita prévia ou exportação, **Then** recebe diagnóstico com tela e componente de origem; o rascunho continua salvável.
4. **Given** campos com valores iniciais, **When** o autor interage na prévia e volta a uma tela, **Then** os valores iniciais são restaurados sem modificar o projeto.

---

### User Story 3 - Exportar apresentação navegável e editável (Priority: P1)

Como analista, quero compartilhar a jornada em PowerPoint, permitindo que outra pessoa explore os botões e edite os textos e controles.

**Why this priority**: É uma entrega central do produto, utilizável sem geração Java ou roteiro animado.
**Independent Test**: Exportar projeto de três telas sem roteiro; abrir no PowerPoint, navegar pelos botões e editar rótulo e campo como elementos separados.

**Acceptance Scenarios**:

1. **Given** jornada válida sem roteiro, **When** exportada, **Then** o PPTX oferece exploração a partir da tela inicial e não avança automaticamente.
2. **Given** componentes de texto, botão, campo, seleção, relatório e teclado, **When** o PPT é editado, **Then** seus textos e formas permanecem editáveis, sem uma imagem única da tela.
3. **Given** botões com destinos, **When** clicados em modo de apresentação, **Then** abrem o slide correto, inclusive em ciclos entre telas.
4. **Given** o editor em máquina sem Office, **When** exporta, **Then** produz o arquivo sem iniciar ou instalar PowerPoint.

---

### User Story 4 - Montar uma demonstração guiada (Priority: P2)

Como apresentador, quero configurar etapas de cursor, clique, digitação, espera e navegação para demonstrar um percurso reproduzível.

**Why this priority**: Acrescenta explicação temporal a uma jornada que já pode ser explorada.
**Independent Test**: Em projeto de referência, configurar roteiro que preenche um campo e visita novamente uma tela; comparar prévia e execução da demonstração no PPT.

**Acceptance Scenarios**:

1. **Given** uma jornada válida, **When** o autor adiciona e reordena etapas explícitas, **Then** pode reproduzir o roteiro e identificar uma etapa inválida.
2. **Given** um roteiro válido, **When** exportado com o PPT, **Then** o arquivo oferece exploração manual e demonstração temporizada como percursos separados.
3. **Given** duas visitas à mesma tela, **When** o roteiro é exportado, **Then** cada visita reproduz seu estado e suas etapas sem herdar animações da visita anterior.
4. **Given** uma etapa cujo alvo foi removido, **When** inicia a demonstração ou exporta seu percurso, **Then** recebe diagnóstico indicando a etapa; a exploração válida pode ser exportada separadamente.

---

### User Story 5 - Gerar fontes Java para o sistema EFL (Priority: P2)

Como desenvolvedor do sistema de destino, quero receber fontes Java 1.6/Swing com EFL, com telas e navegação separadas, para integrar a jornada ao sistema existente.

**Why this priority**: Reutiliza a definição visual no sistema corporativo e depende do contrato da biblioteca real.
**Independent Test**: Gerar projeto de referência e compilar/executar suas telas e navegação no ambiente Java 1.6/EFL de homologação, fora do editor.

**Acceptance Scenarios**:

1. **Given** jornada válida e perfil EFL homologado, **When** gera fontes, **Then** recebe telas, controlador de navegação, recursos e instruções de build e integração.
2. **Given** aspas, acentos, barras e quebras de linha, **When** gera textos e comentários, **Then** os fontes preservam o conteúdo sem produzir sintaxe inválida.
3. **Given** nomes que colidem após normalização, **When** solicita geração, **Then** recebe diagnóstico com os identificadores envolvidos e pode corrigir sem alteração silenciosa dos nomes.
4. **Given** editor sem Java instalado, **When** gera fontes, **Then** a geração funciona sem compilar nem executar Java.
5. **Given** assinatura EFL não homologada, **When** solicita exportação Java, **Then** recebe diagnóstico de capacidade indisponível; o sistema não inventa APIs nem anuncia compilação validada.

---

### User Story 6 - Usar a ferramenta offline no Windows (Priority: P1)

Como usuário final, quero abrir o editor em uma janela própria sem preparar ambiente de desenvolvimento.

**Why this priority**: É uma restrição de distribuição de todas as funcionalidades.
**Independent Test**: Copiar a distribuição para Windows limpo, desligar a rede, abrir, salvar, reabrir e exportar PPTX e fontes sem instalar runtimes.

**Acceptance Scenarios**:

1. **Given** Windows 10/11 x64 sem Go, Node, Java, Office ou WebView adicional, **When** abre o EXE, **Then** obtém uma janela própria e os recursos necessários estão disponíveis offline.
2. **Given** exportação em andamento, **When** o usuário cancela, **Then** o editor permanece utilizável, o destino anterior é preservado e arquivos parciais não são publicados.
3. **Given** caminho com espaços e acentos, **When** salva e exporta, **Then** as operações funcionam ou apresentam erro localizado sem perda do projeto.

### Edge Cases

- Projeto sem telas é um rascunho válido para salvar; prévia e exportação exigem uma tela inicial existente.
- Telas podem ter ciclos; roteiros são listas finitas sem saltos ou repetição implícita.
- Troca de tamanho de tela não reposiciona componentes silenciosamente; elementos fora dos limites são destacados e bloqueiam exportação.
- Remoção com referências mostra o impacto; desfazer restaura o elemento e as referências em uma operação.
- Arquivos incompatíveis, recursos ausentes, identificadores repetidos e valores inválidos não podem ser ignorados.
- Campo interativo do PPT é uma representação editável no modo de edição: não recebe entrada real em apresentação. Digitação é demonstrativa.
- Fontes e estilos sem equivalente EFL devem produzir diagnóstico específico do alvo Java, sem bloquear PPT compatível.
- Falta de permissão, destino existente, disco cheio e cancelamento não podem destruir a saída anterior.

## Requirements *(mandatory)*

### Functional Requirements

| ID | Requisito verificável |
| --- | --- |
| FR-001 | Criar, renomear, salvar e reabrir projetos locais com telas, jornada, roteiros e recursos. |
| FR-002 | Criar, ordenar, renomear, duplicar e remover telas; dimensões configuráveis, padrão 800 × 600. |
| FR-003 | Disponibilizar sete famílias: botão, rótulo, título, campo, seleção, relatório/lista e teclado numérico. |
| FR-004 | Inserir componentes por arrastar e soltar, selecionar, mover e redimensionar dentro da área da tela. |
| FR-005 | Editar nomes, textos, dimensões, posição e propriedades específicas previstas no catálogo. |
| FR-006 | Alinhar, duplicar, ordenar camadas e remover componentes por comandos e atalhos. |
| FR-007 | Oferecer zoom proporcional de 25% a 400%, sem alterar coordenadas persistidas. |
| FR-008 | Desfazer/refazer as últimas 100 operações confirmadas da sessão; arraste completo conta como uma operação. |
| FR-009 | Salvar com substituição atômica e oferecer recuperação da última cópia automática válida após interrupção. |
| FR-010 | Validar formato e versão na abertura, preservando sessão ativa quando o arquivo não puder ser interpretado. |
| FR-011 | Definir uma tela inicial e ações de botão sem ação ou navegar para tela existente. |
| FR-012 | Oferecer prévia navegável com entrada em campos e seleção, isolada da edição do projeto. |
| FR-013 | Mostrar referências afetadas antes de remover telas/componentes e localizar referências quebradas. |
| FR-014 | Exportar PPTX com menu de entrada, uma tela por slide de exploração e links internos corretos. |
| FR-015 | Preservar textos e controles como formas editáveis; permitir imagens apenas em recursos decorativos. |
| FR-016 | Preservar proporção, posição, conteúdo e ordem visual dos componentes na exportação. |
| FR-017 | Criar, nomear, duplicar, remover e ordenar roteiros de demonstração. |
| FR-018 | Configurar etapas explícitas de mover cursor, clicar, digitar, esperar e navegar, com duração em milissegundos. |
| FR-019 | Validar e reproduzir roteiros na prévia; pausar, continuar e reiniciar a partir do começo. |
| FR-020 | Exportar demonstração com animações nativas, digitação simulada e ocorrências distintas para revisitas. |
| FR-021 | Manter exploração manual e demonstração temporizada em percursos separados no mesmo PPTX. |
| FR-022 | Gerar fontes Java 1.6/Swing com EFL a partir do projeto validado, usando perfil de biblioteca homologado. |
| FR-023 | Separar fontes por tela, fábricas de componentes e controlador de navegação, sem listeners de navegação vazios. |
| FR-024 | Validar pacotes e identificadores Java, colisões e escape de textos, comentários e recursos. |
| FR-025 | Incluir recursos locais, inventário de dependências e instruções de compilação e integração, sem exigir JDK no editor. |
| FR-026 | Validar por alvo; informar severidade, código e localização dos problemas; nenhum artefato parcial deve ser publicado. |
| FR-027 | Informar os limites entre entrada real na prévia/Java e entrada apenas simulada no PPT. |
| FR-028 | Exibir progresso de exportação, permitir cancelamento e preservar a saída anterior em falha. |

### Non-Functional Requirements

| ID | Restrição / critério |
| --- | --- |
| NFR-001 | Janela própria, funcionamento offline em Windows 10/11 x64, sem runtimes adicionais instalados pelo usuário. |
| NFR-002 | Módulos de domínio independentes de UI e formatos de saída, composição manual e responsabilidades coesas. |
| NFR-003 | Documentação em português de cada método e variável: GoDoc onde aplicável, comentários locais, docstrings de ferramentas e Javadoc na saída Java. Identificadores podem permanecer em inglês. |
| NFR-004 | Mesma revisão de projeto e mesmas opções produzem conteúdo equivalente, com ordenação determinística de artefatos. |
| NFR-005 | Abertura, salvamento e exportação mantêm dados anteriores em caso de falha; recursos não podem escapar do pacote. |
| NFR-006 | Interface e diagnósticos em português, textos com Unicode e caminhos com espaços e acentos. |
| NFR-007 | No conjunto nominal de 10 telas × 50 componentes e roteiro de 50 etapas: abrir/salvar até 5 s, exportar cada alvo até 30 s; operações longas não bloqueiam interação/cancelamento por mais de 1 s. |
| NFR-008 | Homologação real de PPT no PowerPoint desktop 2016 e Microsoft 365 e de fontes no ambiente Java 1.6/EFL; dublês e validação estrutural não substituem esses testes. |

Os tempos são metas futuras em Windows x64, CPU de quatro núcleos, 8 GB de RAM e SSD, com aquecimento já concluído; medir três execuções e registrar a mediana. Não são resultados medidos nesta entrega.

### Key Entities

- **Projeto**: documento persistido que reúne telas, jornada, roteiros e recursos.
- **Tela**: área lógica com dimensões, componentes e ordem de apresentação.
- **Componente**: instância de uma família do catálogo, com geometria e propriedades.
- **Catálogo**: definição versionada de propriedades e capacidades de cada família.
- **Jornada**: tela inicial e ações que conectam botões a telas.
- **Roteiro / etapa**: sequência finita de ações demonstrativas e seus tempos.
- **Recurso**: imagem local embutida e identificada, sem dependência de URL externa.
- **Diagnóstico / artefato**: problema localizado ou saída produzida por uma exportação.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: O cenário de três telas pode ser criado, salvo e reaberto sem perda de conteúdo, geometria, links ou etapas.
- **SC-002**: Todos os botões configurados do cenário de referência levam ao destino esperado na prévia, apresentação exploratória e aplicação Java homologada.
- **SC-003**: Todos os textos e controles do cenário exportado continuam selecionáveis e editáveis individualmente no PowerPoint, sem reparo ao abrir.
- **SC-004**: A demonstração de referência reproduz todas as etapas na ordem, inclusive uma revisita, com diferença total de tempo de até 10% ou 500 ms, o maior valor, em relação ao roteiro.
- **SC-005**: O editor completa abertura, edição, salvamento e geração das duas saídas em máquina limpa e offline sem instalar ferramentas.
- **SC-006**: Todas as falhas previstas de abertura/salvamento/exportação mantêm o último projeto e artefatos válidos, com diagnóstico que permite localizar a causa.
- **SC-007**: As operações nominais atendem aos tempos de NFR-007 no ambiente de referência registrado.

## Assumptions

- Público inicial: autores de jornadas e desenvolvedores que integram telas ao sistema EFL; uso individual em uma sessão por projeto.
- A tecnologia antiga permanece apenas como alvo solicitado de geração: Java 1.6/Swing com EFL. O editor não reutilizará React nem Node.
- A biblioteca EFL e exemplos funcionais do sistema hospedeiro serão fornecidos pela equipe mantenedora antes de concluir a homologação Java. Sua ausência está registrada em research.md.
- Não há promessa de identidade pixel a pixel entre renderizadores; geometria, conteúdo, hierarquia e navegação seguem os contratos e tolerâncias de homologação.
- Fora da primeira versão: importação `.page`, compilação Java no editor, geração Go da jornada, banco de negócio, APIs externas, multiusuário, plugins executáveis e gravação automática.
- A entrega atual contém implementação em apps/gerador-jornada, testes automatizados e EXE Windows disponível para teste manual. Requisitos ainda não atendidos ou não homologados permanecem no backlog; testes verdes não constituem aprovação humana.
- A implementação é decomposta nas features de domínio [Projeto](../002-projeto/spec.md), [Catálogo](../003-catalogo/spec.md), [Tela](../004-tela/spec.md), [Jornada](../005-jornada/spec.md) e [Demonstração](../006-demonstracao/spec.md); esta feature mantém a visão integrada e as duas saídas.

## Revisão de usabilidade e apresentação — 2026-10-18

- **FR-029**: Manter a ação Java/EFL na tela, desabilitada temporariamente por decisão do usuário, preservando o gerador para retomada.
- **FR-030**: Exportar transições fade/deslizar/revelar ou nenhuma; fade de 300 ms é o padrão. Manter separados efeitos de entrada e animações do roteiro.
- **FR-031**: Manter Salvar alterações e Descartar acessíveis fora da rolagem do inspector; confirmar o formulário em uma edição reversível e renderizar o estado confirmado imediatamente.
- **FR-032**: Oferecer seleção múltipla, caixa de seleção, Ctrl+A, Ctrl/Shift+clique, setas, Shift+setas, Ctrl+D e Delete coletivos, respeitando foco textual e os limites geométricos do conjunto.
- **FR-033**: Oferecer alinhamento horizontal do texto dentro da faixa, independente da posição/tamanho do elemento; persistir e reproduzir o alinhamento nas representações.
- **FR-034**: Evitar reconstruções de primitivas/imagens e de formulários em seleção inalterada; medir o custo de seleção com 500 controles e registrar ambiente e limites da medição.

### Rastreabilidade da revisão

| Requisitos | Responsabilidade | Validação |
| --- | --- | --- |
| FR-029 | 025-barra-ferramentas/FR-008; 030-visualizacao-codigo/FR-008 | Java/EFL aparece desabilitado no EXE |
| FR-030 | 006-demonstracao/FR-007–FR-008; contrato PPT | Opcões e efeitos no arquivo; homologação PowerPoint separada |
| FR-031 | 027-painel-direito/FR-008–FR-009; 028-propriedades/FR-008–FR-011 | Rodapé, confirmação, descartar e histórico |
| FR-032 | 004-tela/FR-006–FR-007; 020-elemento-edicao/FR-008–FR-013 | Grupo, foco, bordas, atalhos e Undo |
| FR-033 | 004-tela/FR-008; 018-faixa-titulo/FR-008–FR-009 | Texto alinha sem mover bounds; persistência e PPT |
| FR-034 | 019-canvas-edicao/FR-009–FR-010 | Benchmark antes/depois e revisão manual |

## Identidade da aplicação — 2026-10-18

- **FR-035**: Usar a imagem fornecida pelo usuário como ícone da janela, da barra de tarefas e do executável Windows, incorporada ao programa e disponível offline; preservar a imagem original e gerar representações adequadas às escalas do sistema. Responsabilidade do adaptador desktop/distribuição, sem alterar os modelos de domínio. Detalhamento: 025-barra-ferramentas/FR-010.

**Aceitação**: com o EXE copiado para outra pasta, sem imagens externas, o Windows apresenta o ícone escolhido no Explorer e na janela/barra de tarefas. Verificar tamanhos pequeno/grande e DPI 100/150/200%; inspeção de recursos PE comprova incorporação, mas não substitui essa revisão visual.

## Ações contextuais e exclusão — 2026-10-18

- **FR-038**: Identificar esta distribuição como versão 1.0.0.0 no título/cabeçalho da aplicação e nas propriedades de versão do arquivo/produto do EXE Windows. Manter uma fonte única para esses valores; a versão do software não modifica a versão do manifesto `.jornada`, do catálogo ou de contratos de domínio.

- **FR-036**: Clique com botão direito sobre um componente deve selecionar o alvo e abrir as ações Alterar propriedades e Excluir elemento. Ao clicar num membro de uma seleção múltipla, preservar o grupo para exclusão; Alterar propriedades deve editar somente o alvo clicado. Clique direito vazio não oferece ações de componente.
- **FR-037**: Oferecer Excluir imediatamente à direita de Ajustar; habilitar somente com seleção válida da tela ativa. Menu, botão e Delete compartilham confirmação Sim/Não; Não não altera documento/histórico, Sim remove o conjunto em uma única ação reversível. Delete em campo de texto não remove componentes. Não acumular confirmações simultâneas.

Responsabilidades: domínio Tela mantém remoção, referências e histórico por meio da aplicação; adaptador desktop mantém seleção, menu, foco, confirmação e disponibilidade das ações. Rastreabilidade: 004-tela/FR-009, 020-elemento-edicao/FR-014–FR-016, 025-barra-ferramentas/FR-011–FR-012; T088–T090.

## Prévia da tela atual — 2026-10-18

- **FR-039**: Oferecer botão separado Prévia Tela Atual, mantendo Prévia para o projeto completo. Abrir somente a tela ativa em janela independente, com fundo, controles e valores temporários interativos, sem navegar para outras telas, mudar a tela inicial, persistir dados ou alterar histórico. Validar a tela solicitada sem exigir uma jornada completa válida ou outras telas válidas; sem tela ativa, desabilitar a ação. Exibir valores confirmados no modelo; rascunhos de propriedades precisam de Salvar alterações antes da prévia.

Rastreabilidade: 004-tela/FR-010 (composição), 029-previa-interativa/FR-008–FR-010 (sessão de prévia) e 025-barra-ferramentas/FR-014 (entrada na interface). Implementação: T092.

## ID editável para arquivos — 2026-10-18

- **FR-040**: Oferecer atributo ID em todas as telas e instâncias de componentes visuais, separado do UUID interno. Nas telas, prefixo Pnl obrigatório (`01 · Identificação` sugere `PnlIdentificacao`); nos componentes, nome livre sem prefixo obrigatório, conforme confirmação do usuário. O ID é nome-base sem extensão, com 1–40 caracteres ASCII alfanuméricos e início em letra. Não aceitar acentos, espaços, símbolos, caminhos ou nomes de dispositivos Windows. O limite inclui Pnl quando presente. Nome amigável pode continuar acentuado.
- **FR-041**: Preservar IDs confirmados ao renomear legendas e ao exportar; gerar arquivos principais `<ID>.java`. IDs de tela são únicos no projeto e IDs de componente na tela, sem distinguir caixa. Edição inválida/duplicada não altera documento ou histórico; duplicação/criação e migração podem sugerir sufixos numéricos dentro do limite. UUIDs continuam resolvendo vínculos e navegação.
- **FR-042**: Gravar `.jornada` com schemaVersion 2 e `fileId` obrigatório para Screen/Component. Ler v1 por migração determinística, preservando nomes, UUIDs, geometria e vínculos; versão desconhecida e v2 inválida devem ser recusadas com mensagem. A versão do software permanece 1.0.0.0 e catalogVersion permanece 1. A ação Java/EFL continua visível e desabilitada; alterações no gerador não homologam o perfil proprietário.

**Rastreabilidade**: 004-tela/FR-011–FR-012; 003-catalogo/FR-006; 002-projeto/FR-006; 028-propriedades/FR-012–FR-013; contratos de projeto/catálogo/Java; T093–T096.

## Documentação e distribuição de 64 bits — 2026-10-18

- **FR-043**: Manter README completo em apps/gerador-jornada, com propósito, público, componentes/famílias/conjuntos, propriedades, atalhos, tecnologias e versões, abordagem por domínio, pré-requisitos separados de uso/build, estado real das saídas e links para specs. README da raiz deve orientar a descoberta desse guia.
- **FR-044**: Entregar e recompilar a distribuição corrente somente para Windows amd64 em dist. Manter Linux amd64 e macOS amd64/arm64 como base de configuração e documentação para uso futuro, sem exigir novos builds ou homologação desses alvos nesta etapa. Usar compilação nativa CGO, Docker Linux opcional a partir do Windows e runners nativos para macOS. Cada pacote efetivamente gerado deve identificar software 1.0.0.0, plataforma, arquitetura, inventário e hashes; recusar cabeçalho incompatível e não criar artefato falso para alvo sem toolchain. Configuração/compilação não equivalem a homologação humana.
- **FR-045**: Armazenar recuperação em diretório nativo por usuário; preservar LOCALAPPDATA no Windows, respeitar XDG_CONFIG_HOME no Linux e Application Support no macOS. Não depender do diretório de execução nem inserir essas regras de plataforma nos domínios.

Windows 10/11 x64 é o alvo da entrega corrente. Por orientação do usuário em 2026-10-18, a base Linux/macOS fica como referência futura e sua execução sai do escopo atual; suporte final, DPI, dependências de runtime e operação em máquina limpa precisam de evidência por plataforma. MacOS usa macOS 15 como alvo inicial do bundle. Gate Java/EFL e homologação PowerPoint permanecem independentes. Fontes: T097–T100; 002-projeto/FR-007.

## Maximizar como tela cheia — 2026-10-18

- **FR-046**: Na janela principal Windows, converter Maximizar (botão/comando nativo, incluindo duplo clique na barra) em tela cheia do editor, ocultando a decoração da janela e a barra de tarefas enquanto estiver em primeiro plano. Usar a API de tela cheia do Fyne; não esconder globalmente a barra do Windows nem alterar outras aplicações. A captura pertence somente ao adaptador desktop.
- **FR-047**: Abrir o editor em tela cheia por padrão a cada inicialização, com Restaurar disponível, sem persistir essa escolha no projeto. Oferecer botão Tela cheia/Restaurar dentro da aplicação e alternância por F12; Escape em tela cheia retorna à janela normal. No Windows, capturar as teclas no HWND para funcionar também com propriedades focadas e evitar alternância repetida ao segurar F12. A mudança preserva documento, revisão, histórico, seleção, zoom e rascunho, sem confirmar/cancelar operações de negócio. Linux/macOS mantêm maximização nativa e usam o botão/API de tela cheia; o recebimento de teclas pelo canvas depende do foco/driver e deve ser homologado.

**Aceitação**: abrir em tela cheia, usar Restaurar para obter uma janela normal, maximizar e confirmar ausência das barras do sistema; restaurar por botão/F12/Escape, repetir e conferir posição/dimensões normais e dados inalterados. Testar campo focado, monitores/DPI distintos, Alt+Tab, minimizar/fechar e projeto não salvo. Não persistir o modo em `.jornada`. Fontes: 025-barra-ferramentas/FR-015; 019-canvas-edicao/FR-011; T101–T102.


**Sequência de abertura de FR-047**: primeiro criar/exibir a janela nativa, depois solicitar a entrada real em tela cheia pela mesma operação do F12. Não considerar maximização comum nem um indicador interno como evidência de tela cheia. Conferir ausência de decoração/barra de tarefas no teste manual Windows.


- **FR-048**: Integrar Sair somente em tela cheia e os temas Claro (padrão), Escuro e Azul com preferência local. Responsabilidade do adaptador desktop e de 025-barra-ferramentas/FR-016–FR-018; nenhuma mudança nos modelos Projeto, Catálogo, Tela, Jornada e Demonstração ou contratos de exportação.

## Expansão pós-v1 autorizada — 2026-10-19

O usuário aprovou iniciar a implementação do plano de paridade de **todo o acervo visual Java/EFL**, não somente das sete famílias iniciais. FR-003, FR-042 e as restrições descritas como v1 registram o escopo/histórico inicial; os requisitos abaixo os estendem intencionalmente, sem apagar esse histórico. A autorização para implementar não é homologação nem aceite do produto. O checklist de Catálogo permanece com 13 itens marcados e 2 abertos, com prosseguimento expressamente autorizado; não alterar suas caixas.

- **FR-049**: Inventariar e reproduzir todas as classes visuais, composições, nomes de skin e estados do acervo Java em editor, prévia e PPTX. Manter skins como variantes/presets de famílias semânticas, incluindo ToggleButton, TextArea, PasswordField, ComboBox/UF, Table, AlphaKeyboard, ProgressBar e Panel/Dialog, botão de duas legendas e todas as subclasses de Report. Classificar também helpers, máscaras, contêineres/scroll e classes adicionais descobertas; nenhuma referência pode sumir do inventário por não pertencer às sete famílias. Matriz e bloqueios em [research.md](research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19).
- **FR-050**: Gerar biblioteca estática de skins como dados puros e recursos incorporados por adaptador, a partir de importação de desenvolvimento reproduzível. Editor não interpreta XML Java nem executa Java em runtime; mantém fonte/face real, RGBA, imagens, hashes, dimensões, alinhamento e espaçamento comprovados. Distinguir declaração XML, transformação Java e evidência visual; bold/italic declarados não provam sua aplicação pelo builder.
- **FR-051**: Evoluir para schemaVersion 3/catalogVersion 2, com leitura explícita de v1/v2/v3 e escrita v3. Aplicar nova aparência também aos projetos antigos, preservando UUIDs, FileIDs, conteúdo, opções/linhas, vínculos, ordem, bounds e overrides explícitos; conservar alinhamento histórico vazio como esquerda. Novos defaults geométricos só valem para inserções novas. Abrir não sobrescreve original; primeiro salvamento migrado avisa incompatibilidade e preserva cópia anterior.
- **FR-052**: Compartilhar resolução de skin, métricas/fontes, estados e composição entre canvas, miniaturas, prévia e PPTX. Prévia opera foco/cursor, máscaras, seleção, rolagem e modal locais em estado isolado; não cobrir cena fiel com widget nativo opaco de outra aparência. O tema/chrome do editor não muda com o perfil da jornada.
- **FR-053**: Exportar recursos decorativos e texto/formas editáveis com typeface real, alpha, recorte e espaçamento da cena. Não rasterizar controle/tela nem substituir fonte por Arial silenciosamente. Incorporar fontes somente com licença/formato permitidos; indisponibilidade bloqueia o modo fiel ou exige ambiente explicitamente documentado, nunca equivalência fictícia.
- **FR-054**: Encerrar a expansão somente após cobertura classe/skin/estado, testes automatizados, homologações manuais Windows/DPI e PowerPoint e aprovação explícita da revisão testada, com rastreabilidade separada. Recursos, delegates ou referências ausentes bloqueiam fidelidade total das entradas dependentes; parciais podem ser entregues para teste, não como acervo homologado.

### Aceitação adicional

**Estado técnico da continuação (não aceite):** perfil por documento/cópia profunda,
captura integral dos metadados da biblioteca original em cópias privadas (definições,
fontes, aliases, recursos e todas as revisões arquivadas, inclusive não usadas e todos
os estados), sem enriquecer perfis privados com builtin nem fazer upgrade. Uso corrente
só limita diagnósticos e acréscimo de entrada missing. Desvio arquitetural justificado
pelo bug crítico save/reopen → AddComponent `defaultconfirm`, preserva inserir famílias/
trocar variantes após reopen. Não é incorporação binária integral: alvo SUBSET de bytes
usados e permitidos permanece; payload selado mais amplo altera o hash. Resolução exata de
até 16 revisões arquivadas diretas e integração desktop/PPT estão implementados.
Captura é de metadados e conserva bytes já incorporados, não copia novos binários
por declaração JSON allowed. Wire atual Style/Fonts/Aliases/Resources é transitório;
appearance normalizado/overrides com presença completa/RGBA/descritores continuam
alvo aberto. Galerias estáticas16 Java/XML não são screenshots Java/homologação DPI.
Após a correção integral de Capture, conforme main: gofmt/vet PASS, **15 suítes Go
compiladas, não executadas** (bloqueio prévio por política Windows preservado);
Python **28 PASS** (15 importador + 4 packaging + 9 Java 6, erro esperado de
`String.strip` incluído, não homologação EFL) e importer `--check` PASS, mesmo
inventário 214/8/270/47. Novo EXE Windows amd64 construído pela tarefa VS Code,
hash/tamanho no [quickstart](quickstart.md#7-incremento-do-acervo-javaefl--2026-10-19-evidência-vigente).
Distribuição Windows gerada e verificada; checksum do ZIP no sidecar externo
`.zip.sha256`. Reempacotar documentos preserva o EXE. Embedding PPT,
Windows/DPI/PowerPoint e aceite continuam pendentes.

1. Dada fixture v1/v2 com customizações, abrir e salvar em v3 mantém todos os dados e bounds, aplica a nova aparência nos campos não sobrescritos e não transforma teclado antigo 200 × 280 em 200 × 350.
2. Dada entrada da matriz, localizar classe/skin/estado, origem, fonte/face, imagens/hashes, geometria, comportamento local, alvos, tarefas e evidências; entrada incompleta continua bloqueada e visível.
3. Dada a mesma configuração validada, comparar editor/prévia/PPT e referência Java nos estados aplicáveis; screenshots produzidos pelo novo renderer não são referência independente do legado.

Sem importador `.page`, serviços/periféricos reais, negócio externo ou reativação Java. `JAVA_PROFILE_UNVERIFIED` permanece; paridade dos controles Go não homologa a saída Java. Responsabilidades locais: 002/FR-008–FR-010, 003/FR-007–FR-011, 004/FR-013–FR-016, 005/FR-006–FR-007 e 006/FR-009–FR-010. Integração/adaptadores: T108–T125 em [tasks.md](tasks.md).

## Clarificação dos diagnósticos de recursos — correção de 2026-10-21

**Removido em 2026-10-23 (pedido do usuário; [T134](tasks.md#t134-remocao-recursos-referencia)):** o campo **Recursos / referência** do inspector e o botão de diagnóstico da prévia deixam de existir; o texto abaixo é histórico. O diagnóstico técnico e as travas de exportação/empacotamento permanecem.

Correção do escopo de diagnóstico previsto em **FR-026, FR-050 e FR-052**, detalhada em [028/FR-016](../028-propriedades/spec.md#aparência-e-famílias-pós-v1--2026-10-19), sem novo ID funcional: **Recursos / referência** deve consultar somente o componente selecionado no rascunho e sua revisão exata, incluindo definições de partes compostas e recursos de todos os estados. Um caminho compartilhado não autoriza atribuir avisos de outra definição ou revisão à seleção. Integridade do perfil e de todos os recursos permanece global; corrupção não usada continua bloqueante, sem fallback ou reparação silenciosa.

O desktop prefere diagnóstico local detalhado, não repete versões genéricas da condição nem exige métricas já medidas com SFNT. Identidade de deduplicação: código/tipo/nome/estado/recurso, não mensagem. Texto informativo completo, não editável, com quebra e rolagem vertical de referência de 160 unidades lógicas; skin/fonte/versão atualizam a consulta do rascunho sem confirmar dados. A prévia completa já usa `componentDiagnostics` por componente: o ajuste elimina a repetição dos avisos do documento inteiro para cada componente, sem introduzir outro fluxo de prévia.

Produção e dois novos arquivos de regressões implementados. **Evidência final do executor principal para esta correção:** gofmt em quatro arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` PASS, com 15 suítes **somente compiladas, não executadas**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade do pacote PASS (CRC, SHA256SUMS, sidecar SHA-256, PE amd64, README exato, inventário exato, ausência de backups terminados em `~` e ambas as cópias do EXE idênticas). Testes Go não executados por política Windows, sem contorno; UI nativa, DPI, revisão manual e aceite não realizados. [T128](tasks.md#phase-22-correção-de-recursos--referência--2026-10-21) permanece aberta como verificação mista; identidades finais somente na [seção 10 do quickstart](quickstart.md#10-correção-de-recursos--referência--2026-10-21), sem equiparar compilação a execução ou build a homologação. Nenhuma alteração de schema, catálogo, domínio ou camadas integra esta correção; históricos e arquivos de camadas modificados pelo usuário são preservados.

## Sobreposição de componentes — 2026-10-20

- **FR-056**: Permitir colocar um componente à frente/atrás de outro da mesma tela ou nos extremos, em todas as 16 famílias, pelo formulário **Sobreposição** confirmado em **Salvar alterações** junto às demais propriedades, atomicamente e com Descartar/Undo/Redo preservados. Manter posição atual é o padrão neutro; Frente ↑/Fundo ↓ continuam comandos de um passo. Reordenar somente a entrada escolhida na lista canônica fundo→frente, sem nova relação persistida, campos/schema, alteração geométrica ou deslocamento das camadas dos filhos.

Contrato do formulário em [028/FR-018](../028-propriedades/spec.md#sobreposição-de-componentes--2026-10-20) e invariantes em [004/FR-018](../004-tela/spec.md#sobreposição-de-componentes--2026-10-20)/T014; coordenação única em [T127](tasks.md#phase-21-sobreposição-de-componentes--2026-10-20). Produção e quatro arquivos de regressões implementados; testes Go não executados por política Windows, sem contorno. Vet/compilação integral/build aguardam o executor principal, sem PASS antecipado; [seção 9](quickstart.md#9-sobreposição-de-componentes--2026-10-20) mantém identidade nova, manuais e aceite pendentes. Registros/hashes anteriores são históricos, não evidência deste incremento.

## Alinhamento textual universal — 2026-10-20

- **FR-055**: Estender o alinhamento textual de FR-033 a todos os elementos capazes de apresentar texto nas 16 famílias, inclusive vazios. Expor **Alinhamento do texto** após o conteúdo e antes da aparência avançada, com **Padrão da skin / Esquerda / Centro / Direita** e dica de que altera somente o texto, sem mover ou redimensionar o componente. Separar a seção **Posição do elemento**. Reutilizar `Style.TextAlign` e a resolução textual compartilhada em canvas/prévia/PPT, inclusive popup do combo e **Sem seleção**; alinhar nas regiões internas de texto, preservando bounds, pais/filhos, fundos e hit areas. Carregar a escolha confirmada, aplicar por Salvar alterações e preservar Descartar, Undo/Redo e salvar/reabrir, sem novo schema, catálogo ou domínio.

**Aceitação adicional (pendente):** exercitar as 16 famílias × esquerda/centro/direita, restaurar Padrão da skin, conferir vazio/preenchido e regiões de compostos (células, áreas de relatório, legendas de teclas), popup, histórico e round-trip. Nenhuma escolha textual muda geometria, hierarquia, conteúdo ou áreas clicáveis. Padrão da skin remove somente o override; migração histórica de vazio antigo para esquerda permanece.

Rastreabilidade: [028/FR-017](../028-propriedades/spec.md#alinhamento-textual-universal--2026-10-20), [004/FR-017](../004-tela/spec.md#alinhamento-textual-universal--2026-10-20) e [T126](tasks.md#phase-20-alinhamento-textual-universal--2026-10-20). Após produção e quatro novos arquivos de testes, o executor principal informou em 20/10: **gofmt/vet PASS, 15 suítes Go compiladas sem execução, quatro testes Python de empacotamento PASS e tarefa VS Code de compilação Windows PASS**. Execução Go não realizada pelo bloqueio prévio de política Windows, sem contorno nem PASS dos testes. T126/004-T013 permanecem abertas com verificação parcial; manuais Windows/DPI/PowerPoint e aceite pendentes. Gates Java/EFL, fontes/licenças e PPT/embedding não mudam. Os 28 Python e referências de 19/10 continuam históricos, não reexecutados nesta evidência. Registro técnico e identidade do EXE na [seção 8 do quickstart](quickstart.md#8-alinhamento-textual-universal--2026-10-20); novo ZIP aguarda geração e confirmação do executor principal, sem hash antecipado.

## Log de aplicação — 2026-10-21

- **FR-057**: Disponibilizar log local offline em JSONL por processo/sessão, iniciado antes do Fyne, identificando o executável realmente aberto por SHA-256/basename, versão do produto/Go, SO, arquitetura, PID e sessão. Registrar erros apresentados por diálogo/status com classe/tipo seguro, errno/código de perfil quando disponíveis e fingerprint, sem mensagem arbitrária bruta. Permitir captura manual do diagnóstico do componente selecionado ou, sem seleção, dos componentes da tela atual (estado salvo, sem rascunho nem publicação de alterações); não há registro automático por seleção/troca de skin/fonte/versão nem botão de diagnóstico na prévia (removidos em 2026-10-23, [T134](tasks.md#t134-remocao-recursos-referencia)). Oferecer **Log de aplicação** no cabeçalho com caminho efetivo, **Registrar diagnóstico agora**, **Copiar caminho** e **Abrir pasta**.

**Contrato técnico e privacidade:** manter códigos e localização original código/tipo/nome/estado/recurso, inclusive condições técnicas que a apresentação amigável não transcreve. Somente valores conhecidos/permitidos e mensagens builtin exatas podem aparecer literalmente; strings, mensagens, caminhos e metadados personalizados desconhecidos ficam em hash ou omitidos. `displayed` indica presença do código no relatório efetivamente mostrado, não disponibilidade; `report_sha256` correlaciona esse relatório sem armazená-lo. Resolução registra revisão exata, inclusive arquivada, SHA-256/tamanho/MIME esperados, status/licença e hash/tamanho observados ou código de erro, abrangendo partes/estados futuros, sem fallback ou mudança de política de recursos.

Usar cache local por usuário (**%LOCALAPPDATA%/GeradorJornada/logs** no Windows, `os.UserCacheDir`), com fallback para **os.TempDir()/GeradorJornada/logs** e destino visível. Escrita síncrona com mutex/Sync, segmentos de até **2 MiB + três backups**, retenção de **20 arquivos próprios quando removíveis** e deduplicação de snapshots por **2 segundos/até 256 identidades**, contornável pela captura explícita. Falhas do logger ficam em `Err` e visíveis na UI. Não registrar conteúdos de entradas/senhas, nomes de projeto/componente ou `FilePath`; standard log/Fyne registra só hash e contagem de bytes, nunca texto bruto. Sem rede, novo domínio ou mudança de schemaVersion **3**, catalogVersion **2** ou produto **1.0.0.0**.

**Estado final técnico — pronto para coleta (2026-10-21):** implementação de produção concluída e 11 arquivos Go formatados (sete existentes e quatro novos, incluindo dois de testes). O executor principal observou `APPLICATION_LOG_VET_PASS` e `APPLICATION_LOG_COMPILE_ONLY_PASS`: **15 suítes somente compiladas, ZERO execuções Go**, incluindo os 12 testes novos (9 de núcleo + 3 de UI). Quatro testes Python de empacotamento PASS; tarefa Windows na raiz exit 0; empacotamento sem `--check` e integridade somente leitura PASS. Ambos os EXEs foram atualizados e são idênticos; identidades reais dos dois EXEs e do ZIP e limites das verificações na [seção 13](quickstart.md#13-log-de-aplicação--2026-10-21). README da aplicação congelado já incluído no pacote; produto/schema/catálogo **1.0.0.0/3/2** inalterados.

**Aceitação e diagnóstico real pendentes:** fechar a instância antiga, abrir o novo EXE identificado e correlacionar startup, relatório visível/rascunho, diagnósticos técnicos e resolução exata; repetir captura dentro da janela de deduplicação e conferir privacidade, fallback, rotação/retenção, falha de gravação visível e documento/histórico intactos. A falha continua relatada pelo usuário; **causa não mapeada e solução não comprovada** até reproduzir e enviar o log real. Nenhum aplicativo ou teste Go foi executado nesta entrega; presença de marcadores no EXE não prova startup ou uso do logger. Integração única em [T131 aberta como verificação mista](tasks.md#phase-25-log-de-aplicação--2026-10-21), clarificação em [028/FR-016](../028-propriedades/spec.md#log-da-consulta-visível--clarificação-de-fr-016-2026-10-21). Política Windows preservada sem contorno; execução dos testes Go, startup/runtime do logger, UI nativa/DPI 100/150/200%, Office e aceite pendentes. T130 e gates T120–T122 permanecem abertos; `manualApproval` e `redistributionApproval` continuam `pending`. Históricos/PASS/hashes anteriores não validam LOG; não criar logs artificiais nem presumir que o usuário já abriu o novo EXE.

<a id="builtin-filtrado-fr-058"></a>
## Redução do builtin por dependências indisponíveis — 2026-10-21

**Decisão posterior do usuário:** “pra ficar facil oq não achar exclui caso eu ter necessidade eu faço uma importação em uma change”. A exclusão é **do builtin operacional gerado**, não do acervo legado, do inventário auditável nem de documentos existentes. Esta decisão altera o escopo operacional de FR-049/FR-050/FR-054; as oito variantes abaixo ficam adiadas para uma change futura, sem declarar o acervo integralmente homologado.

- **FR-058**: Na importação de desenvolvimento, retirar a **definição inteira** do builtin se qualquer referência de recurso `imagePath`, `slider` ou `font:type`, em qualquer estado ou fora dele, estiver ausente/inválida, tiver caminho lógico recusado ou tipo incompatível (`RESOURCE_MISSING`, `RESOURCE_PATH_INVALID`, `RESOURCE_TYPE_MISMATCH`). Não conservar só o estado default nem trocar recursos por homônimos. Se uma parte filha falhar, adiar **toda a composição, pai e todos os filhos**, sem publicar órfãos; usar a identidade tipada `(Type, Name)` e o vínculo de composição extraído do XML, não nomes simples ou prefixos textuais. Preservar `source_tree` integral e `xml_sha256` do XML original, origens, conflitos, seleções/pins, bytes disponíveis e todas as fontes. Registrar cada definição adiada no campo aditivo `deferred_definitions`, com `type`, `name`, `definition` e `reasons`, sem elevar a versão do schema.

### Inventário ativo e adiamentos

Geração existente: **206 definições ativas, 8 fontes, 270 recursos e 31 diagnósticos ativos**. As **8 definições adiadas** conservam os motivos das **16 referências ausentes**; **214 definições originais = 206 ativas + 8 adiadas**, sem apagar ou editar o legado. Menos diagnósticos ativos não significa recursos recuperados ou fidelidade comprovada.

| Tipo XML | Variantes adiadas, identificadas por tipo e nome |
| --- | --- |
| `button` | `botao_jogos_futuros`, `botao_loteca_especial`, `botao_lotogol_especial`, `botao_timemania_especial` |
| `label` | `label_faixa_titulo_secundario_federal` |
| `panel` | `painel_titulo_loteca_especial`, `painel_titulo_lotogol_especial`, `painel_titulo_timemania_especial` |

Nova revisão de skin: **`legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a`**. Permanecem **manifest schema 1 / Profile schema 1 / project schemaVersion 3 / catalogVersion 2 / produto 1.0.0.0**; a revisão de conteúdo não é versão do produto nem hash de um novo EXE.

### Preservação e critérios de aceitação ainda abertos

1. Novos projetos usam o builtin filtrado; a biblioteca ativa continua **integral dentro desse escopo**, inclusive variantes não usadas e todos os estados. Não é subset por uso nem limpeza de bytes/fontes sem consumidores.
2. Perfis próprios/privados, revisões arquivadas e documentos existentes conservam definições, referências, conteúdo, IDs/FileIDs, bounds, overrides e ordem. A redução não provoca nova migração, upgrade, troca de perfil/revisão, remoção retroativa ou fallback builtin por nome; seleção/resolução exata e round-trip permanecem obrigatórios.
3. As oito entradas não são oferecidas para nova inserção/troca de variante no builtin filtrado. Um homônimo de outro tipo e composições independentes permanecem disponíveis se válidos. A lista do seletor é por família: **206 é o total de definições ativas, não 206 famílias ou opções em cada seletor**. Definições adiadas continuam consultáveis no manifesto, não como opções operacionais prontas.
4. Preservar fontes existentes, avisos/licenças e gates `DELEGATE_MISSING`, `FIDELITY_UNVERIFIED`, Java/EFL e PPT/font embedding. Falta de delegate, restrição legal e ausência de homologação não são corrigidas pelo filtro; não anunciar “100% fiel” ou redistribuição autorizada.
5. Reimportar somente em change futura, com os recursos **exatos**, origem/hash e autorização verificáveis; regenerar, conferir determinismo/integridade e validar novamente composição, estados, seletor, inserção e preservação de projetos antes de liberar cada entrada. Não inventar recurso substituto, alias ou pin para forçar disponibilidade.

<a id="dependencias-ausentes-fr-059"></a>

### FR-059 — Retirada de referências e ações sem implementação

**Decisão posterior e prevalente do usuário:** “o que não tiver no projeto apaga a referencia completa e a ação se precisar vou solicitar pra anexar”. Substitui, somente nesse alcance, as exigências anteriores de manter Java/EFL visível/desabilitado e de conservar referências de pintura a delegates externos ausentes.

- Excluir o produtor/chamadas de delegates Java da apresentação local e avisos que exigiam classes inexistentes para senha/diálogo, preservando os controles Go, seus recursos, geometria e interação implementada. Não resolver por filtro genérico de diagnósticos.
- Excluir botão/callback/adaptador Java/EFL e suporte de validação Java 6; não deixar ação inoperante. Integração e visualização Java dependem de nova change com a biblioteca real.
- Oferecer somente Padrão/numeric/alpha no teclado. Excluir o helper de calculadora e impedir desenho, hit areas e comandos dessa variante. Ler valores antigos sem converter ou apagar o documento ao salvar outra propriedade; substituição exige escolha explícita.
- Manter IDs, FileIDs, ordem, dimensões, vínculos, estilos, perfis próprios/arquivados e bytes originais. Recursos realmente presentes não são removidos por licença desconhecida. Licenças, escala, métricas, portabilidade, limitações reais de layout/PPT e aprovação manual continuam independentes.

Coordenação em [T133](tasks.md#t133-dependencias-ausentes), evidências e limites na [seção 15](quickstart.md#15-dependencias-ausentes-removidas). Os registros abaixo e anteriores da entrega filtrada são históricos, não reintroduzem a ação Java nem o vínculo de delegates.

**Verificação mista — entrega técnica concluída, homologação pendente:** **19 testes Python do importador PASS**, incluindo descoberta e execução explícitas dos quatro métodos novos após correção do aninhamento; **4 testes Python de empacotamento e importador `--check` PASS**. Regressões Go de biblioteca, seletores, Capture e perfis próprios/arquivados implementadas; gofmt/vet e compilação das **15 suítes Go PASS, ZERO execuções Go**, respeitando a política Windows. Tarefa Windows, empacotamento sem `--check` e integridade somente leitura PASS: 549 referências XML ativas, 270 recursos/264 blobs, origens/pin claro/XML original, CRC/checksums, manifesto, README congelado e dois EXEs idênticos. Identidades reais na [seção 14](quickstart.md#14-builtin-filtrado). Coordenação em [T132 aberta](tasks.md#t132-builtin-filtrado) e clarificação em [003/FR-012](../003-catalogo/spec.md#builtin-filtrado-fr-012); execução Go, roteiros de UI/round-trip, DPI/Office e aceite continuam pendentes, sem fechar caixas anteriores ou aprovar fidelidade integral.

<a id="codigos-barras-2026-10-23"></a>
### FR-060 — Entrada de código de barras e apresentação processada

**Origem**: pedido de 2026-10-23 e reforço explícito para seguir GitHub Spec Kit, decompondo especificações e implementação pelos domínios existentes. **Status**: planejamento documental, sem código, novo build ou aprovação funcional.

Integrar os componentes [031 — entrada validada](../031-codigo-barras-efl/spec.md) e [032 — processado parametrizável](../032-codigo-barras-processado-efl/spec.md), respeitando os seus contratos visuais e o modelo tipado. Entrada usa oito campos vazios, valida pares/DV geral e captura integral local, bloqueando sua conclusão quando inválida. Processado permite blocos/linhas independentes, valores configuráveis e escolha no inspector entre CPF/CNPJ/Nome/Nome fantasia/Razão social. Não derivar dados pessoais de boleto, consultar serviços, processar pagamento, criar rota implícita ou reintroduzir Java/EFL.

**Fronteiras obrigatórias**: [002/FR-011](../002-projeto/spec.md#codigos-barras-2026-10-23) preserva o documento/compatibilidade; [003/FR-013–FR-015](../003-catalogo/spec.md#codigos-barras-2026-10-23) define tipos/regras puros; [004/FR-019](../004-tela/spec.md#codigos-barras-2026-10-23) compõe instâncias; [005/FR-008](../005-jornada/spec.md#codigos-barras-2026-10-23) mantém rotas explícitas; [006/FR-011](../006-demonstracao/spec.md#codigos-barras-2026-10-23) mantém quadros/capacidades. A coordenadora implementa somente casos de uso/runtime e adaptadores de UI/cena/PPT, não cópias dos modelos ou algoritmos desses domínios.

**Aceitação**: cenários 031/US1–US4 e 032/US1–US3 rastreados até as tarefas de domínio e [T135–T139](tasks.md#tarefas-codigos-barras-2026-10-23), sem compartilhar runtime entre instâncias ou persistir capturas. PPT contém textos/formas editáveis, não interação bancária. Anexos não recebidos e classes de módulo não fornecidas são limitações expressas para fidelidade/equivalência, não resultados de teste. Plano, contratos e roteiro manual acompanham esta mudança; execução e aceite continuam pendentes.

<a id="menu-aplicacao-2026-10-23"></a>
### FR-061 — Menu aplicação e implementação autorizada dos três compostos

Integrar [033 — Menu aplicação](../033-menu-aplicacao-efl/spec.md) em **Itens legados**, ao lado de 031/032. Uma interface por instância alterna raiz/submenu selecionado, com árvore autorável, inicial fixo escolhido, layouts root/standard/financial/bolao e caminho local isolado. Voltar remove nível, Cancelar restaura raiz, Sair encerra somente prévia. Folha sem destino informa seleção simulada; somente destino explicitamente autorado e validado por Jornada pode mudar tela na prévia completa. Prévia Tela Atual mantém interação local sem transição entre telas. Não executar metadados de movimento, rede, contas, autenticação, periféricos, QR, timers de negócio ou Java/EFL.

**Donos únicos**: [003/FR-016](../003-catalogo/spec.md#menu-aplicacao-2026-10-23) modelos/regras puros; [004/FR-020](../004-tela/spec.md#menu-aplicacao-2026-10-23) composição/IDs; [002/FR-012](../002-projeto/spec.md#menu-aplicacao-2026-10-23) codec/compatibilidade; [005/FR-009](../005-jornada/spec.md#menu-aplicacao-2026-10-23) destinos; [006/FR-012](../006-demonstracao/spec.md#menu-aplicacao-2026-10-23) quadros/capacidades. Coordenadora é runtime/adaptadores, não novo domínio financeiro. Catálogo 3 planejado passa a cobrir **031/032/033 juntos**, mantendo schema3; base validada anterior 3/2 é histórica até validação principal do incremento.

**Autorização e aceite**: usuário autorizou implementar todos os incrementos 031/032/033 e fará teste manual posterior. Essa autorização permite prosseguir sem alterar checklists históricos/customizados ou marcar revisões abertas como aprovadas; não é homologação, aceite funcional, licença ou liberação de distribuição. Esta entrega é documental, sem Go/build/testes. Novas tarefas [T140–T144](tasks.md#tarefas-menu-aplicacao-2026-10-23) permanecem abertas. Imagens realmente vazias e módulos de cálculo não fornecidos mantêm seus gates de fidelidade/equivalência, não impedem implementação local explicitamente autorizada.

## Reimplementação visual solicitada — 2026-10-24

O pedido do usuário nesta revisão é reaplicar aos componentes existentes o estilo do acervo `material-legado`, utilizando as skills Go locais. Este incremento executa o recorte de FR-050/FR-052/FR-053 nas sete famílias atuais, nos conjuntos informativos e nos compostos 031/032/033; não declara entregue a expansão integral de tipos do backlog pós-v1.

- Sem override, resolver o estilo do legado: botão por variante confirm/cancel/yes/no, rótulo `label_texto_azul`, título `label_titulo_1`, entrada `textField/type1`, seleção com marcadores e relatório com suas partes originais. Skin explícita e personalizações de cor/fonte/alinhamento prevalecem.
- A prévia usa a mesma cena do editor/PPT. Captura nativa de entrada permanece transparente sobre a decoração EFL; foco, cursor, seleção, clipboard e valores digitados são locais à visita.
- Teclado novo usa 200 × 350 e 14 teclas. Bounds antigos são mantidos; setas editam a posição do cursor, TAB avança pelos campos habilitados e LIMPAR emite Backspace. Não interpreta a legenda como apagar todo o campo.
- Importação é uma operação de desenvolvimento, com dados JSON e bytes incorporados. Runtime não interpreta XML Java. Preservar preto sem alpha como opaco e alpha explicitamente zero como transparente; famílias e flags físicas de fonte vêm do SFNT, separadas das flags XML.
- Recursos ausentes não recebem substitutos silenciosos: oito definições são retiradas da oferta; escolhas antigas não resolvidas continuam preservadas e mostram diagnóstico. PPT recusa aparência não resolvida antes de publicar.

Aprovação visual Windows/DPI e Office permanece separada dos testes executados. Este recorte não aprova o acervo completo nem altera a suspensão de Java/EFL. Evidências e limites na seção 22 do quickstart; tarefas T148–T150 e T147.

### Gate de fidelidade visual — critérios objetivos (2026-10-24)

A seção anterior descreve a intenção. Este gate define o que **reprova** a entrega, para que “skin aplicada” não seja aceita por inspeção parcial. A skill [jornada-visual-parity](../../.agents/skills/jornada-visual-parity/SKILL.md) é obrigatória em qualquer alteração visual.

| Critério | Evidência mínima | Reprova se |
| --- | --- | --- |
| Default sem skin explícita | teste de cena inserindo cada família com `Style.Skin == ""` | a cena cai em cor genérica (`#4864A0`/`#2563EB`) ou não contém decoração do legado |
| Cena compartilhada | teste/inspeção dos mesmos `presentation.Group` em canvas, prévia e PPT | a prévia usa `widget.Button`/`Entry`/`Select` opaco no lugar da pintura |
| Fonte física | asserção de `canvas.Text.FontSource` não nulo e família SFNT esperada | só o nome/tipo foi transportado; Fyne usa a fonte do tema |
| Alpha | testes de fundo transparente, preto sem alpha e alfa parcial | `alpha=0` explícito vira ausência ou transparência vira branco |
| Importação | `import_legacy_skin.py --check` PASS; 205 chaves/224 recursos/6 fontes; 8 indisponíveis fora da oferta | runtime interpreta XML, oferta definição com recurso faltante ou não valida SHA/tamanho |
| Interação | teste de cursor/TAB/Backspace/clipboard com decoração EFL visível | o controle nativo cobre ou substitui a skin |
| Compostos | fixture de relatório/teclado/menu/ficha com partes próprias | só a caixa externa recebe estilo; partes internas ficam genéricas |
| Diagnóstico | PPT recusa skin/fonte não resolvida e projeto antigo é preservado | fallback silencioso, remoção de escolha ou publicação parcial |

O gate é critério de conclusão técnica; DPI 100/150/200%, PowerPoint real e aceite humano continuam em T148–T150/quickstart §22.

## Correção do menu e skins de jogos — 1.0.0.2

A árvore passa a derivar da fonte integral recuperada, preservada em `material-legado/menu/menu.xml`, substituindo o preenchimento sintético anterior. O inventário de PNGs não determina itens do menu. [033/FR-011–FR-012](../033-menu-aplicacao-efl/spec.md#estado-verificado-1-0-0-2) define fidelidade de conteúdo e estados originais; evidência na seção 24 do quickstart. Documentos salvos não são migrados automaticamente.

## Correção dos submenus e limites de movimentos — 1.0.0.3

Conferência com PnlMenu/UIMenuControle preservados corrige faixas, geometria, rodapé e navegação normal/especial. Propriedades editam os metadados existentes; testes exercitam a janela reconstruída. O cadastro movimento.xml/movimentos.xml não foi recuperado, portanto sua configuração e referências externas não são consideradas completas. Entrega somente EXE Windows, conforme regra vigente do README; detalhes na spec 033 e quickstart seção 25.

## Launcher e estrutura de projeto — incremento 034 / 1.0.0.4

Implementados launcher sem documento implícito, recentes locais, criação em branco/com exemplo, árvore de workspace com criar/importar e destino explícito entre jornadas em `MenuItem.TargetJourney`. Writer **3/4**, leitores **1/1, 2/1, 3/2, 3/3 e 3/4**. A [spec 034](../034-launcher-estrutura-projeto/spec.md) é o backlog canônico deste incremento; o [quickstart 034](../034-launcher-estrutura-projeto/quickstart.md) registra validação técnica C1/C2/C3. O estado inicial 1.0.0.1 acima continua histórico, sem converter pendências dos outros domínios em entrega. Aceite humano permanece pendente.

## Inicialização e splash — integração 035/036 / 1.0.0.9

O novo modelo de domínio do ciclo de inicialização está na [035](../035-inicializacao/spec.md) e em internal/startup, separado de UI/mídia/IO. A apresentação e os adaptadores têm [spec/modelo/plano próprios na 036](../036-splash-screen/spec.md). Carregar recursos e reproduzir vídeo v2 com áudio em paralelo; liberar transição somente após ambos. Splash fullscreen somente vídeo, quadro integral sem ampliação acima de 100%, autoria exclusivamente no launcher. Falha/cancelamento impedem avanço. Modelo persistido dos domínios de autoria continua 3/4; aprovação humana ainda não concedida.

## Inicialização opcional e fullscreen — 1.0.0.10

Integrar 035/FR-008–FR-010 (política por usuário), 036/FR-011/FR-013/FR-014 (mídia fullscreen e armazenamento fora do projeto) e 034/FR-014 (checkbox no launcher). Primeira apresentação obrigatória uma vez; repetição desmarcada por padrão e optativa nas próximas aberturas. Vídeo ocupa a maior área proporcional da tela, sem limite de resolução original nem corte. Não alterar formato .jornada ou regra de entrega somente EXE; evidência atual em 036/quickstart.

## Vídeo em janela — 1.0.0.11

O pedido vigente substitui o fullscreen do incremento anterior: apresentação em janela normal centralizada, 960×540 unidades Fyne inicialmente, mantendo áudio, vídeo inteiro, espera de recursos e opção de primeiro uso/repetição. Detalhes em 036/FR-001/FR-011 e DESIGN. Não alterar o modelo de domínio 035 nem formato de projeto. Exibição no ambiente do usuário ainda pendente de teste.

## Arquivos no workspace — 2026-10-25 / 1.0.0.15

- **FR-062**: Exportar PowerPoint DEVE usar automaticamente a pasta do workspace atual e o nome do arquivo `.jornada`, com extensão `.pptx`, sem solicitar pasta ou nome de saída. As opções de conteúdo/transição permanecem disponíveis; apresentação existente exige apenas confirmação de substituição e só é alterada após geração completa. Sem workspace definido, orientar definir/salvar o projeto, sem exportar para pasta implícita. O primeiro salvamento de jornada com workspace já configurado DEVE usar essa pasta, preservando homônimos. Formatos confirmados: `.jornada` e `.pptx`; `.ppt` antigo permanece fora deste incremento.

**Aceitação**: abrir `Meu projeto (2).jornada`, exportar e encontrar `Meu projeto (2).pptx` ao lado, sem seletor de destino. Cancelar substituição preserva a apresentação anterior. Salvar documento novo com workspace definido cria nele o pacote, sem pedir a pasta novamente. A mudança é de destino de arquivos, sem alteração de formato, cena ou publicação Git.
