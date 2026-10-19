# Feature Specification: Gerador de Jornada

**Feature Branch**: `001-gerador-jornada` (identificador da feature; não há branch Git criada)
**Created**: 2026-10-17
**Status**: Draft — especificação documentada; produto não implementado
**Input**: Editor desktop offline para Windows, com componentes EFL arrastáveis, jornadas navegáveis, exportação PPTX editável e animada e geração de fontes Java 1.6/Swing com EFL.

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
- A entrega atual é documental. Os testes e capacidades descritos são trabalho futuro.
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
- **FR-044**: Preparar distribuição em dist para Windows amd64, Linux amd64 e macOS amd64/arm64. Usar compilação nativa CGO, Docker Linux opcional a partir do Windows e runners nativos para macOS. Cada pacote efetivamente gerado deve identificar software 1.0.0.0, plataforma, arquitetura, inventário e hashes; recusar cabeçalho incompatível e não criar artefato falso para alvo sem toolchain. Configuração/compilação não equivalem a homologação humana.
- **FR-045**: Armazenar recuperação em diretório nativo por usuário; preservar LOCALAPPDATA no Windows, respeitar XDG_CONFIG_HOME no Linux e Application Support no macOS. Não depender do diretório de execução nem inserir essas regras de plataforma nos domínios.

Windows 10/11 x64 permanece a plataforma inicial. Os pacotes Linux/macOS ampliam os alvos de build; suporte final, DPI, dependências de runtime e operação em máquina limpa precisam de evidência por plataforma. MacOS usa macOS 15 como alvo inicial do bundle. Gate Java/EFL e homologação PowerPoint permanecem independentes. Fontes: T097–T100; 002-projeto/FR-007.
