# Feature Specification: 038 — Integração MCP, Agentes de IA e Fábrica de Componentes Reutilizáveis

**Feature Branch**: `038-integracao-mcp-agentes` (identificador; sem branch Git criada)
**Created**: 2026-10-25
**Version**: 0.4 — revisão documental dos apontamentos P0/P1/P1/P2
**Status**: implementação autorizada por invocações reiteradas de speckit-implement e pedidos de continuidade, após análise dos replanejamentos. Autoria transitória e controles de sessão comprovados no EXE 1.0.0.19; revisão/staging de jornada implementados internamente. Exposição da revisão, paridade visual completa, persistência, biblioteca, host real e aceite funcional ainda pendentes. Ver tarefas e evidências no quickstart.
**Input**: Propostas iniciais, parecer técnico de correção da v0.3 e autorização explícita para revisão documental dos quatro apontamentos em 2026-10-25. O usuário reafirmou a independência dos modos manual e MCP com LLM externa, a aparência do legado e a divisão por domínios como requisitos da spec e da implementação futura.
**Constituição de referência**: [1.4.0](../../.specify/memory/constitution.md), emenda MCP aprovada explicitamente pelo usuário e aplicada em 2026-10-25.

O Jornada Studio deve permitir autoria e reutilização de componentes pelo editor manual e por agentes compatíveis com MCP. No modo manual, o autor abre o software e utiliza sua interface; no modo MCP, uma LLM externa interpreta solicitações e usa um cliente compatível. Um modo não substitui, desativa ou condiciona o outro. Os dois canais utilizam as mesmas regras de domínio e o mesmo estilo visual do legado, produzindo jornadas editáveis.

Essa equivalência é obrigação da especificação e da implementação: não basta declarar paridade documental, compartilhar nomes ou implementar aparência em apenas um caminho. Projeto, Catálogo, Tela, Jornada e Demonstração conservam seus modelos e regras; GUI e MCP são adaptadores desse núcleo. A responsabilidade dos dados está dividida no [modelo por domínio](data-model.md), e os comportamentos visuais seguem o [DESIGN da biblioteca](DESIGN.md).

Esta feature coordena a integração; não cria outro domínio nem duplica o backlog de expansão do Catálogo. A aprovação “Aprovar emenda e planejamento” resolveu o gate constitucional. Após a fase técnica, a invocação explícita de speckit-tasks autorizou a decomposição do [plano](plan.md) em [tarefas](tasks.md) e apontadores nos responsáveis. Nenhum documento comprova disponibilidade de software ou viabilidade experimental. O estado documental anterior foi superado pela autorização de implementação; somente tarefas com evidência registrada podem ser marcadas concluídas. A revisão 1.0.0.16 comprova G0, sem representar os fontes posteriores nem a entrega completa da feature.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Criar e continuar uma jornada por qualquer canal (Priority: P1)

Como autor, quero solicitar a um agente uma jornada com componentes existentes, conferir as alterações e continuar a edição manual no Jornada Studio.

**Why this priority**: Entrega o primeiro fluxo útil sem depender da fábrica ou da expansão de famílias.

**Independent Test**: Em workspace autorizado, descobrir capacidades, criar duas telas, inserir botão e campo, vincular a navegação, revisar, salvar e reabrir no editor.

**Acceptance Scenarios**:

1. **Given** uma versão com integração habilitada, **When** o cliente inicia o modo MCP local, **Then** nenhuma janela ou splash aparece e o cliente consegue consultar capacidades e componentes efetivamente disponíveis.
2. **Given** um rascunho criado pelo agente, **When** o autor revisa e autoriza sua gravação, **Then** a jornada pode ser aberta e editada manualmente, preservando conteúdo, geometria, recursos e referências.
3. **Given** uma jornada criada manualmente, **When** o agente a abre e realiza uma alteração autorizada, **Then** os dados não abrangidos pela alteração permanecem preservados após salvar e reabrir.
4. **Given** propriedade, família ou variante indisponível, **When** qualquer canal tenta aplicá-la, **Then** a operação é recusada com diagnóstico localizado e sem alteração parcial.
5. **Given** alterações ainda não salvas, **When** a sessão é encerrada ou descartada, **Then** o destino em disco permanece inalterado e o resultado informa o descarte conforme a política apresentada.
6. **Given** ausência de cliente MCP, LLM ou rede, **When** o autor abre o software, **Then** o modo manual e a biblioteca instalada continuam utilizáveis; iniciar ou encerrar o MCP não desabilita esse modo.
7. **Given** mesma família e parâmetros nos dois modos, **When** o autor ou agente insere o componente sem skin explícita, **Then** ambos resolvem o mesmo default do legado, incluindo fontes, cores, alpha, estados e partes internas; overrides válidos produzem o mesmo resultado.

### User Story 2 - Proteger o trabalho durante revisões e gravações (Priority: P1)

Como autor, quero que ações da GUI, de outro agente ou da integração Git não sobrescrevam alterações que não foram revisadas.

**Why this priority**: A gravação pelo agente não pode causar perda de trabalho ou publicação implícita.

**Independent Test**: Abrir o mesmo arquivo em sessões independentes, alterar e salvar por um canal, e tentar salvar o estado anterior pelo outro, repetindo com a ordem invertida.

**Acceptance Scenarios**:

1. **Given** uma revisão apresentada ao autor, **When** o arquivo ou rascunho muda antes da confirmação, **Then** a revisão fica obsoleta e exige nova comparação; a confirmação anterior não autoriza gravar conteúdo diferente.
2. **Given** uma GUI e uma sessão MCP com a mesma base, **When** uma salva e a outra tenta sobrescrever, **Then** a segunda detecta o conflito, preserva ambos os trabalhos e oferece recarga ou salvamento em novo destino autorizado.
3. **Given** operação Git concorrente ou recuperação pendente, **When** uma gravação incompatível é solicitada, **Then** o sistema coordena ou recusa a operação sem corromper jornadas nem alterar o histórico Git implicitamente.
4. **Given** uma chamada cujo resultado foi perdido, **When** o cliente repete a mesma solicitação de mutação identificada, **Then** o sistema não duplica telas, componentes, versões ou arquivos e informa o resultado conhecido ou a necessidade de reconciliação.
5. **Given** operação persistente revisada sem GUI do Jornada Studio, **When** o autor a aprova por canal confiável ou a recusa, **Then** o sistema verifica a origem e o vínculo da decisão antes de executar ou recusar, sem abrir janela Fyne.
6. **Given** `approved: true`, mensagem de aprovação ou arquivo/token produzido ou gravável pelo agente, **When** não existe evidência independente segundo o mecanismo de confiança aprovado, **Then** a operação persistente é recusada e o arquivo anterior permanece íntegro.
7. **Given** autorização expirada, revogada, já consumida ou para outro conteúdo/destino/sessão, **When** é apresentada, **Then** não libera um novo efeito; repetição idempotente apenas devolve resultado comprovado.
8. **Given** cliente sem confirmação confiável e ausência de canal alternativo autorizado, **When** é solicitada gravação/publicação, **Then** ela falha de forma segura; leitura e rascunho continuam possíveis dentro das permissões existentes.
9. **Given** rascunho admitido e quota de novas edições esgotada, **When** o autor nega/expira uma revisão e solicita outra, **Then** ainda pode revisar/finalizar com autorização ou descartar o trabalho; replay não renova validade nem repete efeito. Reservas de finalização impedem admitir um alvo pendente que dependa de uma posição inexistente de commit.

### User Story 3 - Criar e reutilizar uma composição (Priority: P2)

Como autor, quero salvar um conjunto de componentes existentes na biblioteca e inseri-lo em outras jornadas manualmente ou com auxílio de agente.

**Why this priority**: Permite reutilização sem transformar cada conjunto em uma família nova ou exigir nova compilação para cada composição.

**Independent Test**: Criar um cabeçalho no editor, publicá-lo na biblioteca local, inseri-lo por MCP em outra jornada e repetir a criação pelo MCP com inserção manual.

**Acceptance Scenarios**:

1. **Given** um conjunto válido com propriedades, recursos e referências internas, **When** o autor o transforma em definição reutilizável, **Then** pode revisar limites, dependências e conteúdo antes de publicar uma versão.
2. **Given** uma versão publicada, **When** o autor a procura no editor ou pelo MCP, **Then** ambos oferecem a mesma identidade, versão, parâmetros permitidos, disponibilidade e limitações.
3. **Given** uma definição com vínculos internos, **When** ela é instanciada duas vezes, **Then** cada cópia possui identidades independentes e referências remapeadas apenas para os destinos correspondentes.
4. **Given** a mesma versão e os mesmos parâmetros, **When** cada canal insere uma cópia, **Then** os resultados são equivalentes em conteúdo, aparência, recursos e comportamento, descontadas identidades novas.
5. **Given** uma instância inserida, **When** a definição original é atualizada ou retirada da oferta, **Then** a instância existente não muda silenciosamente e continua editável dentro das capacidades da versão do aplicativo.
6. **Given** biblioteca vazia, rascunho, versão publicada, incompatibilidade ou recurso ausente, **When** o autor abre a biblioteca manual, **Then** vê o estado e as ações possíveis descritos no DESIGN, com foco por teclado e sem depender apenas de cor.
7. **Given** definição validada para publicar, **When** o autor revisa nome, versão, dependências, recursos, diagnósticos e impactos e cancela, **Then** nada é publicado; confirmar executa apenas a revisão autorizada e preserva versões anteriores.
8. **Given** criação de biblioteca antes da entrega de importação/exportação, **When** uma versão é preparada para publicar, **Then** seu pacote definitivo já é produzido/validado antes da revisão; a etapa posterior de transferência não é pré-requisito oculto da publicação.
9. **Given** rascunho validado ou versão existente, **When** o agente solicita revisão de publicação/retirada/transferência pela operação de controle definida, **Then** recebe uma revisão identificada sem gravar; só a confirmação correspondente, com autorização verificável, executa o efeito. Preparar revisão não admite pacote novo fora da quota de edição.

### User Story 4 - Transportar e evoluir a biblioteca sem quebrar jornadas (Priority: P2)

Como autor, quero levar definições para outro workspace ou computador e preservar as jornadas já produzidas.

**Why this priority**: Uma biblioteca pessoal não pode virar dependência invisível dos documentos entregues.

**Independent Test**: Exportar/importar uma definição e abrir uma jornada instanciada em ambiente sem a biblioteca de origem, com a mesma versão compatível do aplicativo.

**Acceptance Scenarios**:

1. **Given** uma composição baseada em capacidades instaladas e imagens incorporadas, **When** a jornada é aberta sem a biblioteca de origem, **Then** seus componentes e imagens permanecem disponíveis sem acesso a caminhos da máquina original.
2. **Given** uma definição exportada com seus recursos autorizados, **When** é importada em outro ambiente compatível, **Then** sua integridade, dependências e identidade de versão são verificadas antes da oferta.
3. **Given** colisão de identidade e versão com conteúdo diferente, **When** ocorre importação, **Then** o sistema recusa substituição silenciosa e informa a colisão.
4. **Given** definição incompatível ou recurso ausente, **When** ela é consultada, **Then** fica indisponível para novas inserções, com justificativa; não há substituição por recurso semelhante.
5. **Given** importação inválida ou conflito de versão, **When** o autor usa o fluxo manual, **Then** recebe motivo e opção de cancelar/corrigir sem substituir a biblioteca existente; uma prévia de importação não publica nem instala conteúdo automaticamente.

### User Story 5 - Analisar legado e propor evolução de componentes (Priority: P3)

Como autor, quero fornecer Java, XML, imagens e configurações para que um agente apresente uma proposta verificável de componente equivalente, com diferenças e dependências explícitas.

**Why this priority**: A análise deve orientar reutilização ou desenvolvimento sem prometer implementação automática de capacidades inexistentes.

**Independent Test**: Analisar o ComboBox do acervo com suas dependências e um cabeçalho composto; classificar o primeiro como necessidade de família planejada e o segundo como composição de famílias existentes.

**Acceptance Scenarios**:

1. **Given** arquivos autorizados do legado, **When** o agente os analisa a partir das evidências fornecidas pelo aplicativo, **Then** apresenta propriedades, estados, comportamento observado, dependências, origem, recursos faltantes e limites da evidência, separando fatos e hipóteses; nenhum material é executado como código.
2. **Given** uma proposta, **When** ela é classificada, **Then** diferencia composição, variante e novo Kind e referencia a especificação responsável quando já existe requisito equivalente.
3. **Given** composição suportada, **When** sua validação é aprovada e o autor autoriza a publicação, **Then** ela entra na biblioteca e pode ser utilizada pelos dois canais.
4. **Given** proposta que exige novo comportamento, **When** o autor a aprova, **Then** ela continua pendente de desenvolvimento; não aparece como componente funcional instalado.
5. **Given** novo Kind integrado por entrega posterior, **When** sua versão é disponibilizada, **Then** ele passa pelos mesmos critérios de descoberta, edição manual, MCP, persistência, prévia e exportação aplicáveis às famílias existentes.
6. **Given** proposta com esquema válido e hipótese de equivalência não comprovada, **When** o Go valida os contratos, **Then** retorna validade estrutural separada de autorização e homologação; não converte a hipótese do agente em fato ou componente instalado.
7. **Given** propriedade de proposta sem valor inicial comprovado, **When** o agente a descreve, **Then** distingue default desconhecido, observado com evidência ou sugerido com motivo, com tipo verificável; não publica se faltar valor obrigatório do contrato instalado.
8. **Given** proposta pendente na sessão, **When** é editada, o documento muda ou o autor tenta encerrar, **Then** a proposta conserva revisão própria, alterações relevantes exigem revalidação e trabalho pendente impede encerramento silencioso; editar proposta não altera por si só a revisão/undo do documento.
9. **Given** proposta analítica classificada como composição/variante suportada, **When** ainda não referencia uma captura concreta, **Then** continua analisável mas não publicável. Autor manual ou agente constrói/captura o conteúdo e vincula identidade/revisão/hash; revisão humana identifica análise, captura e pacote final. Alteração de qualquer vínculo relevante impede usar a autorização anterior, sem geração automática a partir de texto.

### Edge Cases

- Rascunho legível com navegação incompleta versus conteúdo estruturalmente corrompido; somente o primeiro pode ser preservado como rascunho.
- Documento, rascunho ou definição alterados depois da revisão; confirmação repetida, expirada ou para outro destino.
- Chamadas simultâneas, sessões expiradas, cancelamento, desconexão, encerramento inesperado e repetição após resposta perdida.
- Caminhos externos, nomes reservados Windows, links/junctions, mudança do destino durante a operação e arquivos não regulares.
- Entradas enormes, pacotes com expansão excessiva, definições cíclicas, profundidade excessiva, nomes colidentes e dependências indisponíveis.
- Referências externas à seleção capturada, teclado vinculado a campo fora do conjunto, rotas ausentes e identidades internas de compostos.
- Inserção fora dos limites, tela menor que o conjunto, ordem de camadas, estado inicial e conteúdo Unicode.
- Publicação interrompida, biblioteca ausente, versão retirada, importação com hash divergente e falta de permissão de escrita.
- Skin default, overrides, fonte real, alpha zero/parcial, estados desabilitados e recursos de partes internas de compostos.
- Executável Windows distribuído iniciado por cliente MCP com pipes, sem console/rede ou com diretório de execução diferente.
- Texto do legado/documento contendo instruções: é dado não confiável, não autorização para comandos, acessos ou transmissão de conteúdo.

## Requirements *(mandatory)*

### Functional Requirements

#### Capacidades e autoria comum

- **FR-001**: Manter dois modos independentes e complementares: editor manual aberto pelo autor e integração MCP local opcional utilizada por LLM/agente externo. O modo manual não exige IA, cliente MCP, conexão de conta ou rede; o modo MCP não exige abrir a GUI. Um não substitui nem desativa o outro. MCP deve permanecer desativado na abertura padrão e iniciar somente por ação explícita.
- **FR-002**: Editor e MCP devem utilizar os mesmos contratos de criação, alteração, validação, persistência e instanciação, organizados pelos modelos de Projeto, Catálogo, Tela, Jornada e Demonstração. A implementação deve manter regras nesses domínios/casos de uso compartilhados; adaptadores GUI/MCP não podem criar regras ou catálogos paralelos. Decoração, recursos e defaults do legado devem ser aplicados em ambos os caminhos, com evidência de comportamento equivalente.
- **FR-003**: A descoberta deve informar versão do aplicativo, formatos suportados, Kinds, variantes, presets, definições reutilizáveis e capacidades por operação/alvo, distinguindo disponível, indisponível e planejado. Backlog não pode ser anunciado como capacidade instalada.
- **FR-004**: Cada item disponível deve apresentar identidade estável, descrição, propriedades tipadas, defaults, valores permitidos, dimensões e restrições. Campos desconhecidos ou propriedades de outra família devem ser recusados sem mutação parcial.
- **FR-005**: O agente deve poder listar jornadas da raiz autorizada, criar documento vazio ou de exemplo, abrir e consultar documento, criar/configurar/duplicar telas, inserir/configurar/mover/remover componentes e inserir presets suportados, respeitando referências e autorização da operação.
- **FR-006**: Permitir definir início e navegação explícita entre telas. Destinos de menu, quando habilitados, respeitam as alternativas do domínio: submenu, tela local ou jornada da mesma pasta. Metadados de negócio não podem inferir navegação.
- **FR-007**: Preservar UUIDs existentes, FileIDs válidos, limites geométricos, camadas e referências não afetadas. Cópias e instanciações devem gerar novas identidades e remapear referências internas, inclusive de compostos.
- **FR-008**: Autoria de demonstrações e simulação, quando habilitadas, devem expor apenas ações e alvos realmente aceitos pelo runtime. Prévia lógica não implica screenshot; suporte na prévia não implica suporte no roteiro ou no PPT.

#### Sessões, revisão e persistência

- **FR-009**: A edição MCP deve ocorrer em rascunho isolado, com identidade de sessão e revisão esperada por mutação. Serializar ou rejeitar alterações concorrentes na mesma sessão, sem compartilhar estado transitório com a GUI.
- **FR-010**: Permitir consultar, validar, revisar diferenças, salvar, descartar e encerrar rascunhos. Chamadas de edição não salvam automaticamente. Informar prazo de inatividade e política de descarte ao abrir a sessão; sessão expirada não pode ser reutilizada silenciosamente.
- **FR-011**: Operações persistentes devem exigir revisão de operação, destino, diferenças, base, conteúdo identificável, diagnósticos, dependências e impactos, seguida de consentimento verificável controlado pelo usuário, sem depender da GUI Fyne. O usuário deve poder aprovar, negar ou pedir correção. A aplicação deve verificar origem confiável, escopo, validade, revogação e uso da autorização, além da revisão e destino ainda atuais. Declaração do agente, `approved: true`, token reenviado pelo agente ou arquivo gravável por ele não são prova independente de consentimento. Sem mecanismo confiável aprovado e disponível, recusar a operação preservando o estado. Classes R0–R3, ciclo de revisão e limites de confiança estão no [contrato de consentimento](contracts/consentimento.md); mecanismo concreto será definido e validado no planejamento.
- **FR-012**: GUI e MCP devem detectar base persistida alterada antes de substituir arquivo existente, inclusive após gravação pelo outro canal ou mudança externa. Verificação e publicação devem compor operação protegida contra interferência; revisão em memória ou substituição atômica isoladas não bastam.
- **FR-013**: Coordenar gravações com operações Git e recuperação pendente da feature 037 sem exigir vínculo Git para salvar localmente. Salvar pelo MCP não autoriza add, commit, push ou publicação remota.
- **FR-014**: Permitir salvar rascunho estruturalmente legível com diagnósticos após revisão explícita dos problemas; corrupção estrutural não pode substituir documento válido. Execução e exportação exigem validação pertinente ao alvo, conforme a constituição e os contratos existentes.
- **FR-015**: Manter o formato `.jornada` vigente e os pares históricos reconhecidos enquanto não houver evolução aprovada. Dados transitórios do agente/sessão não entram no manifesto. Metadados novos de biblioteca no documento exigem contrato e compatibilidade explícitos antes do desenvolvimento.
- **FR-016**: Mutações repetidas com a mesma identidade de solicitação devem produzir um único efeito durante a sessão e devolver resultado consistente; reutilizar a identidade com parâmetros diferentes deve ser recusado. Após reinício, exigir reconciliação quando o resultado anterior não puder ser comprovado, sem repetir gravações automaticamente.
- **FR-017**: Falhas, cancelamento e encerramento devem preservar a última gravação íntegra e informar o resultado conhecido. Operações entre documentos ou sobre biblioteca devem prever publicação integral ou recuperação explícita; arquivo órfão não pode ser apresentado como sucesso.

#### Taxonomia e biblioteca reutilizável

- **FR-018**: Classificar propostas como composição de Kinds suportados, variante de família ou novo Kind, considerando propriedades, interação e representação, não apenas nome Java ou aparência. Skin não cria Kind; variante pode exigir desenvolvimento adicional.
- **FR-019**: Preservar classificações de [003/FR-007–FR-012](../003-catalogo/spec.md): ComboBox é família planejada, UF é preset, duas legendas são variante Button e apresentações Report são variantes da família correspondente. A 038 coordena; implementação de família permanece no domínio responsável.
- **FR-020**: Persistir definições fora da jornada individual e permitir criar, consultar, validar, revisar, publicar, exportar/importar e retirar versões da oferta pelos dois canais, com os mesmos limites e autorizações. Publicação significa registro local, não envio ao GitHub.
- **FR-021**: Cada definição deve identificar tipo, origem, versão imutável, conteúdo verificável, propriedades permitidas, geometria relativa, componentes, referências, dependências, recursos e limitações por alvo. Nome exibido não é identidade de versão.
- **FR-022**: A primeira biblioteca deve instanciar cópias independentes de componentes existentes. Após inserção, cada componente permanece editável normalmente; não há vínculo vivo, atualização automática ou necessidade de consultar a biblioteca original ao abrir a jornada.
- **FR-023**: Captura e instanciação devem validar o conjunto inteiro antes de publicar o resultado. Remapear referências internas; referências externas devem ser declaradas como parâmetros de destino ou recusadas com diagnóstico, nunca copiadas como vínculos acidentais.
- **FR-024**: Conservar recursos portáteis autorizados e dependências explícitas de capacidades instaladas. Jornadas instanciadas devem incorporar suas imagens e abrir sem a biblioteca de origem em versão compatível do aplicativo; caminhos da máquina de criação não podem ser dependência oculta.
- **FR-025**: Nova versão não pode alterar definições publicadas ou instâncias anteriores. Recusar importação de mesma identidade/versão com conteúdo diferente; retirada da oferta não apaga instâncias ou dados existentes. Versão incompatível ou recurso ausente deve impedir novas inserções sem fallback silencioso.
- **FR-026**: A primeira biblioteca deve usar Kinds, skins e fontes já instalados, além das imagens portáteis suportadas. Instalar fontes, skins externas ou novos comportamentos é ampliação posterior condicionada a contratos de recursos, licença, compatibilidade e validação; não é capacidade presumida da primeira biblioteca.

#### Legado, propostas e novos Kinds

- **FR-027**: O aplicativo/servidor MCP deve fornecer evidências e metadados do legado apenas de fontes autorizadas, separadas do workspace de jornadas quando necessário: referência, tipo, conteúdo permitido, origem, integridade, dependências identificáveis, erros e recursos ausentes. Deve informar capacidades e contratos instalados. Java, XML, imagens e textos são dados não confiáveis; não executar código, templates arbitrários ou instruções contidas nesses materiais. Fornecer evidência não exige compreender semanticamente toda a linguagem Java.
- **FR-028**: O agente/LLM externo é responsável por interpretar a solicitação e o legado e elaborar proposta estruturada com objetivo, evidências, propriedades, estados, comportamento, dependências, classificação, capacidades reutilizáveis/ausentes, diferenças, hipóteses não comprovadas e plano de validação. O aplicativo não depende de inferência por IA para operar; o usuário revisa a aceitabilidade da equivalência e decide sobre autorizações e homologação aplicáveis. O agente não é autoridade final de validação ou publicação.
- **FR-029**: O aplicativo Go deve validar deterministicamente propostas contra contratos implementados: esquema, identidades, versão, Kinds, propriedades, referências, recursos/integridade, compatibilidade, portabilidade, permissões e limites. Retornar dimensões separadas de validade contratual, dependência de desenvolvimento, autorização pendente e homologação necessária. Validade estrutural não comprova equivalência funcional/visual nem transforma hipótese do agente em fato. Contrato e responsabilidades em [propostas do legado](contracts/propostas-legado.md).
- **FR-030**: Publicar composição ou variante declarativa suportada somente após validação e autorizações aplicáveis, vinculadas à versão/conteúdo revisados conforme FR-011. Registrar evidências, verificar dependências, preservar versões imutáveis e impedir substituição silenciosa. Novo Kind ou variante funcional dependente de código segue desenvolvimento Go autorizado, testes, integração ao catálogo e distribuição pelo fluxo Spec-Driven; aprovar a proposta declarativa não autoriza esse trabalho. MCP não possui permissão genérica para alterar fontes, instalar código executável, compilar ou distribuir o aplicativo.
- **FR-031**: Novo Kind ou variante funcional só aparece como disponível quando a versão executada suportar contratos, edição manual, instanciação MCP, validação, cópia, persistência e consumidores aplicáveis. Capacidade de saída ausente deve ficar explícita e não pode gerar perda silenciosa.

#### Aparência, limites e operação local

- **FR-032**: O estilo definido no projeto é o legado: criação manual e MCP, canvas, prévia e PPT devem consumir a mesma representação visual, incluindo default sem skin explícita, overrides válidos, estados, partes de compostos, bytes/fontes reais, cores, alpha e métricas, conforme a [skill de paridade visual](../../.agents/skills/jornada-visual-parity/SKILL.md). A implementação deve demonstrar essa cadeia por fixtures e evidências; não basta declarar conformidade na spec ou aplicar estilo apenas na GUI. Recurso indisponível gera diagnóstico e bloqueia o alvo afetado; não usar aparência genérica ou substituto parecido. Biblioteca e inspector seguem o DESIGN vigente do projeto, sem tema paralelo para MCP ou reutilizáveis.
- **FR-033**: Distinguir aparência estática, edição, interação de prévia, demonstração e exportação. PPT não promete interação não implementada; Java/EFL permanece indisponível enquanto durar sua suspensão e ausência de homologação.
- **FR-034**: O modo MCP deve funcionar sem criar janelas, diálogos, splash ou dependência de eventos Fyne, nem carregar/alterar preferências visuais. Headless não elimina interação humana: revisão e consentimento ocorrem no cliente externo confiável ou canal local autorizado independente das ferramentas do agente, conforme FR-011. Não abrir GUI nem aceitar autoaprovação para contornar falta de mecanismo; negar persistência quando a autorização requerida não puder ser verificada. Não misturar perguntas ao usuário ou logs com mensagens do protocolo. Esse modo não substitui a abertura normal do editor manual.
- **FR-035**: Homologar com o executável Windows x64 efetivamente distribuído e um cliente MCP identificado, cobrindo negociação, descoberta, chamada, cancelamento, desconexão/EOF e encerramento sem processos residuais. Teste com outro artefato não substitui essa evidência.
- **FR-036**: Restringir leituras, gravações e exportações aos locais autorizados por finalidade, validando nomes, arquivos regulares, resolução de links e interferência durante a operação. URI ou identificador não pode funcionar como caminho arbitrário.
- **FR-037**: Documentar e aplicar limites para mensagens, pacotes, recursos descomprimidos, componentes por definição, profundidade, sessões e duração das operações. Exceder limite produz erro identificável e preserva o estado; ciclos são recusados. Valores concretos foram propostos em [contracts/mcp.md](contracts/mcp.md) e devem ser implementados/verificados antes da entrega.
- **FR-038**: A integração local não deve transmitir documentos ou legado a serviço externo por iniciativa própria. Uso de cliente/agente remoto depende da escolha e política do autor; execução offline do editor não significa que um provedor externo de IA funcione offline.
- **FR-039**: Diagnósticos devem distinguir erro de contrato, autorização, conflito de revisão, capacidade indisponível, dependência ausente, limite, cancelamento e falha de persistência, com localização e estado preservado. Consultas grandes devem oferecer resumos e paginação sem omitir diagnósticos relevantes.

### Vocabulário de classificação

| Categoria | Critério | Exemplo e limite |
| --- | --- | --- |
| Composição reutilizável | Combina instâncias de Kinds suportados com layout, parâmetros e referências declarativas | Cabeçalho; expande em componentes editáveis sem introduzir novo comportamento |
| Variante | Especializa uma família sem mudar sua identidade semântica; pode afetar aparência ou propriedades suportadas | Skin de botão ou botão com duas legendas; a segunda hipótese ainda exige suporte específico, não apenas imagem |
| Novo Kind | Introduz semântica, propriedades ou interação que as famílias instaladas não representam | ComboBox; exige entrega de código e compatibilidade dos consumidores antes da oferta |

Os compostos legados `barcodeInput`, `barcodeProcessed` e `applicationMenu` já são Kinds implementados. Sua existência não comprova carregamento de composições arbitrárias. A aprovação de uma proposta visual não dispensa classificar o comportamento necessário.

### Key Entities *(include if feature involves data)*

- **Capacidade instalada**: operação/família/variante suportada pela versão executada, com restrições por alvo.
- **Sessão de autoria**: rascunho isolado, base persistida, revisão atual e ciclo de vida; não integra a jornada salva.
- **Revisão para gravação**: comparação vinculada ao conteúdo, destino, base e autorização correspondente.
- **Definição reutilizável**: composição ou variante declarativa identificada, com parâmetros, dependências e recursos.
- **Versão da definição**: conteúdo imutável verificável e situação na oferta; independente da versão do aplicativo e do catálogo.
- **Instância expandida**: cópia editável de componentes existentes, sem dependência viva da biblioteca original.
- **Proposta de componente**: classificação, fontes do legado, comportamento, diferenças e requisitos de integração; não representa componente instalado.
- **Resultado de validação**: diagnósticos localizados e capacidades verificadas por alvo; aprovação documental não equivale a execução ou aceite.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: US1 cria duas telas com botão/campo e navegação, salva e reabre nos dois sentidos manual/agente, preservando 100% das propriedades e referências previstas na fixture de aceitação.
- **SC-002**: Nos conflitos US2, nenhuma tentativa com base obsoleta substitui a gravação atual; ambos os trabalhos permanecem recuperáveis.
- **SC-003**: Todos os itens oferecidos na matriz de aceitação podem ser inseridos e configurados pelos dois canais; nenhum item apenas planejado é apresentado como utilizável.
- **SC-004**: As duas inserções da mesma composição em US3 possuem zero colisões de identidade e 100% dos vínculos internos apontando para a cópia correta.
- **SC-005**: A jornada US4 abre e permanece editável sem a biblioteca de origem, com todas as imagens necessárias e sem leitura de caminhos da máquina original.
- **SC-006**: Atualizar, retirar ou importar versão colidente não modifica nenhuma instância anterior nos cenários US3/US4.
- **SC-007**: Cada proposta do conjunto US5 informa classificação, dependências e diferenças; nenhum caso dependente de desenvolvimento aparece como instalado antes da entrega correspondente.
- **SC-008**: Nos cenários de limites, cancelamento, entrada inválida e autorização negada, nenhum arquivo anterior é perdido nem ocorre publicação parcial apresentada como sucesso.
- **SC-009**: Na matriz visual de cada incremento, canvas, prévia e apresentação preservam conteúdo, geometria, imagens, estados iniciais, cores e tipografia previstos; limitações reais ficam identificadas e o aceite visual é registrado para a revisão testada.
- **SC-010**: A revisão distribuída passa pelo fluxo local do cliente identificado sem janela/splash, corrupção da comunicação ou processo residual após encerramento; edição e salvamento manual continuam utilizáveis sem rede.
- **SC-011**: Todos os cenários negativos de consentimento do [roteiro](quickstart.md) recusam efeitos persistentes sem autorização válida e independente; aprovação legítima executa apenas a revisão/destino autorizados, sem GUI ou efeito duplicado.
- **SC-012**: Todos os estados e fluxos manuais definidos no [DESIGN](DESIGN.md) têm comportamento observável, ação de recuperação e navegação por teclado previstos; sua implementação e homologação visual serão comprovadas separadamente por revisão executável.
- **SC-013**: Na revisão arquitetural da implementação, regras e dados de Projeto, Catálogo, Tela, Jornada e Demonstração permanecem nos responsáveis do [modelo por domínio](data-model.md); GUI e MCP usam os mesmos casos de uso, sem modelo de negócio paralelo ou tipos de adaptação nos domínios.

## Assumptions

- A autorização atual cobre aplicação da emenda 1.4.0 e planejamento da v0.4, concedidos explicitamente em 2026-10-25. Não autoriza implementação, builds, commits, publicação ou aceite funcional. As orientações sobre os dois modos, legado e domínios vinculam a implementação futura; não ampliam o escopo desta fase documental.
- Cópias independentes, ausência de atualização automática e restrição inicial a skins/fontes instaladas são propostas conservadoras para revisão desta spec, não decisões de implementação já homologadas.
- Reutilizam-se monólito modular, domínios, contratos tipados, persistência em arquivo e distribuição Windows x64 existentes. Não se introduzem banco, serviço web ou motor de plugins executáveis.
- O plano seleciona MCP local `stdio`, modo futuro `JornadaStudio.exe --mcp` e SDK oficial Go v1.8.0 pesquisado em fonte primária. Ainda não são interfaces instaladas ou viabilidade comprovada; o go.mod não foi alterado. Decisões técnicas completas e gates constam do plano e research.
- No código inspecionado em 2026-10-25, há dez Kinds, cinco presets fixos, biblioteca visual incorporada e escrita de manifesto 3/4. Isso não comprova biblioteca extensível, metamodelo completo de propriedades ou homologação de novos tipos.
- Não há promessa de IA embutida no editor; o agente é cliente externo opcional. O autor não precisa de IA para criar e reutilizar composições.

## Escopo incremental e capacidades candidatas

| Incremento | Resultado funcional | Fora desse incremento |
| --- | --- | --- |
| I1 — autoria MCP existente | Descobrir capacidades, listar workspace, abrir/criar jornada, editar telas/componentes existentes, navegar, validar, revisar, salvar e encerrar sessão; continuar manualmente | Biblioteca nova, novos Kinds, importação externa, demonstração avançada e exportação pelo agente |
| I2 — biblioteca comum | Criar/revisar/publicar/instanciar composições pelos dois canais, versionar, retirar da oferta e exportar/importar definições portáteis | Atualização vinculada de instâncias, composições aninhadas, instalação de skins/fontes externas |
| I3 — fábrica assistida | Aplicativo fornece evidências; LLM externa interpreta e propõe composição/variante/Kind; Go valida contratos e usuário decide autorizações; encaminhar desenvolvimento ao domínio responsável | Execução de código legado, LLM embutida obrigatória, implementação/compilação automática pelo servidor MCP |
| Ampliações posteriores | Menus avançados, criação coordenada entre jornadas, demonstrações, simulação e exportação existentes acessíveis ao agente | Servidor remoto, chat interno, renderização headless e plugins executáveis sem especificação própria |

I2 não precisa de IA para funcionar; I3 pode propor desenvolvimento em outro incremento, sem obrigar a implementar todas as famílias planejadas. Cada entrega deve ser verificável; terminar I1 não conclui a feature inteira.

O catálogo candidato original é preservado abaixo como vocabulário funcional para contratos futuros. Nomes não são APIs implementadas e a tabela não determina oferta em `tools/list`.

| Grupo | Ferramentas candidatas | Disponibilidade pretendida |
| --- | --- | --- |
| Descoberta | `catalogo_listar`, `catalogo_descrever`, `catalogo_listar_presets`, `catalogo_listar_skins`, `aplicacao_capacidades` | I1, com apenas itens instalados |
| Workspace/projeto | `workspace_listar_jornadas`, `jornada_criar`, `jornada_abrir`, `jornada_obter` | I1 |
| Importação de jornada | `jornada_importar` | Ampliação posterior, com origem autorizada |
| Telas/componentes | `tela_criar`, `tela_configurar`, `tela_duplicar`, `componente_inserir`, `componente_atualizar`, `componente_mover`, `preset_inserir`, `componente_remover` | I1 |
| Navegação | `jornada_definir_inicio`, `navegacao_vincular` | I1 |
| Menus/documentos relacionados | `menu_configurar_destino`, `jornada_criar_e_vincular` | Ampliação posterior |
| Demonstração/prévia | `demonstracao_criar`, `demonstracao_configurar_etapas`, `demonstracao_validar`, `previa_simular` | Ampliação posterior |
| Validação/gravação | `jornada_validar`, `jornada_revisar_alteracoes`, `jornada_salvar` | I1 |
| Ciclo de vida adicional | Consultar sessão, descartar rascunho e encerrar sessão; nomes finais nos contratos | I1, necessários ao FR-010 |
| Exportação | `jornada_exportar_pptx` | Ampliação posterior, preservando validação do alvo |
| Biblioteca | `componente_propor`, `componente_validar`, `componente_revisar`, `componente_publicar`, `componente_listar_customizados`, `componente_instanciar`, `componente_versionar`, `componente_retirar`, `componente_preparar_importacao`, `componente_importar`, `componente_exportar`; contratos em biblioteca.md | I2 para composições suportadas; I3 vincula análise a captura I2, sem execução de descrições |
| Análise | `legado_analisar_componente` | I3 |

Resources candidatos, somente leitura: `jornada://catalogo`, `jornada://presets`, `jornada://workspace`, `jornada://documento/{id}`, `jornada://diagnosticos/{id}` e `jornada://formatos`. Identificadores devem ser validados; consultas respeitam autorização e paginação; nenhuma URI resolve caminhos arbitrários. Prompts opcionais para construir jornada, criar fluxo entre jornadas, revisar jornada e criar demonstração ficam para ampliação; não substituem regras nem validação.

## Dependências, responsabilidades e condições para planejamento

| Responsável | Regras afetadas e referência |
| --- | --- |
| 038 — coordenação/adaptador | Contratos MCP, capacidades, sessões externas, autorização e integração entre canais; sem dependência de tipos MCP nos domínios |
| [001 — coordenadora](../001-gerador-jornada/spec.md) | Integração geral e evidências, por referência; sem duplicar requisitos de domínio |
| [002 — Projeto](../002-projeto/spec.md) | Compatibilidade, recursos e eventual evolução do pacote; versão atual não muda apenas pela presença do MCP |
| [003 — Catálogo](../003-catalogo/spec.md) | Taxonomia, propriedades, definições, variantes e famílias já planejadas, especialmente FR-007–FR-012 |
| [004 — Tela](../004-tela/spec.md) | Instâncias, geometria, FileIDs, cópia e remapeamento |
| [005 — Jornada](../005-jornada/spec.md) e [006 — Demonstração](../006-demonstracao/spec.md) | Navegação, ações e roteiro; capacidades não presumidas para novos Kinds |
| [023 — Paleta](../023-paleta-componentes/spec.md), [028 — Propriedades](../028-propriedades/spec.md), [029 — Prévia](../029-previa-interativa/spec.md) | Oferta manual, parâmetros editáveis, interação e paridade visual |
| [034 — Workspace](../034-launcher-estrutura-projeto/spec.md) e [037 — GitHub](../037-integracao-github/spec.md) | Resolução de arquivos, gravação, concorrência e recuperação; salvar continua local |

**P0 — resolvido no escopo normativo:** o usuário aprovou explicitamente a emenda e o planejamento em 2026-10-25; a constituição vigente passou a 1.4.0 com a exceção local delimitada. [Research](research.md#proposta-de-exceção-constitucional) registra texto, justificativa, impacto e decisão. Isso resolve o conflito constitucional de escopo; conformidade da implementação e gates Windows/host/visual permanecem futuros. Nenhuma implementação foi autorizada.

**Passagem documental:** v0.4 autorizada para planejamento, incluindo suas premissas; speckit-tasks gerou a decomposição e speckit-analyze encontrou U1/I1/I2/I3/U2. A nova invocação de speckit-plan autoriza corrigir o desenho e sincronizar os artefatos existentes, sem código. Correções iniciais registradas em research D6–D10 e Q038-21–Q038-25. A reanálise apontou U3/U4/I4, tratados neste segundo replanejamento em D11/D12 e Q038-26/Q038-27; análise somente leitura será executada na sequência solicitada, sem autorização de implementação ou aceite técnico.

**P1 — desenho técnico e validação futura:** SDK/versão/licença foram pesquisados e o plano seleciona HostApproval-v1, coordenação comum e limites explícitos. Ainda é necessário comprovar o EXE Windows, qualificar host real, demonstrar proteção de substituição contra interferência e medir custo/compatibilidade. Não há comprovação experimental nesta etapa; capacidades não comprovadas falham fechadas, sem reduzir o requisito.

**P1 — artefatos e responsabilidades:** plan.md, research, modelo, contratos/esquemas, DESIGN e quickstart proposto foram decompostos em tasks.md, com backlogs de domínio preservados. O plano referencia os desenhos 038/[023](../023-paleta-componentes/DESIGN.md)/[028](../028-propriedades/DESIGN.md). Spec, plano e tarefas estão disponíveis para análise formal; quickstart receberá revisão executável/evidências somente quando existirem.

**P2 — evidências de entrega:** implementação, testes automatizados, integração com cliente real, executável distribuído, validação visual/PowerPoint aplicável e aceite humano serão registrados para cada revisão em `quickstart.md`. Checklist de requisitos não comprova nenhum desses resultados.

## Resultado da revisão v0.3 → v0.4

| Apontamento | Correção documental | Estado real |
| --- | --- | --- |
| P0 — exceção constitucional | Texto, justificativa, limites e impacto em research | Emenda 1.4.0 aprovada/aplicada em 2026-10-25; implementação não comprovada |
| P1 — responsabilidades | FR-027–FR-030, US5 e contrato separam agente, Go e usuário | Corrigido documentalmente; validação de implementação não executada |
| P1 — consentimento sem GUI | FR-011/FR-034, US2, SC-011 e HostApproval-v1 impõem verificação e recusa segura | Mecanismo desenhado no plano; host real e demonstração técnica pendentes de implementação/homologação |
| P2 — experiência manual | DESIGN 038, vínculos 023/028, US3/US4, SC-012 e roteiro | Desenho proposto; UI, tarefas e homologação ainda pendentes |
| Orientações do usuário | FR-001/FR-002/FR-032, US1 e SC-013 vinculam modos independentes, estilo legado e divisão por domínio à implementação | Requisitos documentados; não representam software alterado ou aprovado |

Os IDs FR-001–FR-039, US1–US5 e SC-001–SC-010 foram preservados. SC-011–SC-013 são acréscimos. O [checklist](checklists/requirements.md) distingue aprovação da spec/emenda para planejamento, plano proposto, implementação e testes. A autorização documental foi registrada; não há aprovação de software.
