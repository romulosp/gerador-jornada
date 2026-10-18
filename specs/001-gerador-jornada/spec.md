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
