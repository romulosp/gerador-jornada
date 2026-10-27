<!--
Sync Impact Report
Version change: 1.3.0 -> 1.4.0 (expansão compatível: MCP local opcional).
Modified principles: nenhum; preservados I a VII, inclusive distribuição e uso offline.
Added sections: nenhuma.
Updated sections: Restrições e evolução; exceção estreita para a feature 038,
aprovada explicitamente pelo usuário em 2026-10-25 junto ao avanço do planejamento,
sem autorização de implementação. MCP opt-in sem portas de rede, consentimento
verificável e modos manual/MCP complementares com o mesmo domínio e estilo legado.
Removed sections: nenhuma.
Deferred placeholders: nenhum.
Sync impact: templates e comandos leem a constituição em runtime e não foram alterados.
Spec/research/checklist da 038 devem refletir a aprovação; contratos, modelo e DESIGN
preservam suas obrigações. Plano técnico autorizado em operação posterior; tarefas,
implementação e homologação não são aprovadas por esta emenda.
Histórico preservado: 1.2.0 -> 1.3.0 aprovou em 2026-10-24 a integração GitHub
opcional da 037, com conexão voluntária e publicação explícita; limites mantidos.
Não há alteração de formato .jornada, dependências ou código nesta etapa.
Nota de revisão: este relatório é temporário e deve ser removido antes do commit da emenda.
-->
# Gerador de Jornada Constitution

## Core Principles

### I. Domínio e responsabilidade explícitos

O sistema deve ser um monólito modular organizado por Projeto, Catálogo, Tela, Jornada e Demonstração. Regras de domínio não devem importar Fyne, filesystem, PresentationML ou templates Java. Adaptadores traduzem essas tecnologias; casos de uso coordenam o trabalho. Não criar um objeto central com toda a edição, persistência e exportação.

Cada tipo, função e método deve ter uma responsabilidade identificável. Usar composição de structs, dependências explícitas por construtores e interfaces pequenas definidas pelo consumidor, somente quando justificadas por fronteira ou teste. Mudanças nas dependências entre módulos exigem atualização do plano e do registro de decisão em research.md quando alterarem uma decisão.

A divisão por domínio deve existir também nas especificações GitHub Spec Kit.
Projeto, Catálogo, Tela, Jornada e Demonstração devem manter spec.md, plan.md,
data-model.md e tasks.md próprios, com referências aos contratos compartilhados.
A feature coordenadora reúne integração, adaptadores e sprints, sem duplicar tarefas
de implementação dos domínios. IDs de requisito e tarefa são locais à feature;
referências entre features devem identificar o caminho e o ID. Essa organização
documental não cria microserviços nem módulos Go independentes.

### II. Modelo compartilhado e contratos verificáveis

Editor, prévia e exportadores devem derivar do mesmo projeto validado, com identificadores estáveis e propriedades tipadas. Nenhum exportador pode ler widgets da UI ou tratar captura de tela como representação funcional. Estado transitório do editor não é dado de negócio.

Formato de projeto, catálogo, ações, roteiro e exportadores devem seguir os contratos em Markdown da especificação. Versões incompatíveis são recusadas sem modificar o arquivo original. Falhas devem informar a localização do problema e preservar os dados anteriores.

### III. Distribuição offline e ambientes separados

O editor deve executar em Windows 10/11 x64 com janela própria, sem instalação de Go, Node, JVM, Office ou WebView. Recursos e templates necessários devem acompanhar o executável. Build pode exigir Go, compilador C e ferramentas de empacotamento; isso não constitui pré-requisito do usuário final.

Fontes Java são uma saída Java 1.6/Swing com EFL. Compilação e execução dessa saída ocorrem fora do editor. Não prometer autonomia do código Java nem compatibilidade EFL sem evidência com a biblioteca real. PowerPoint é necessário para homologar e reproduzir o PPT conforme a matriz definida, não para gerar o arquivo.

### IV. Documentação de cada elemento do código

Documentar em português pacotes, tipos, interfaces, funções, métodos, campos, constantes e variáveis de pacote, públicos e privados, com comentários no formato GoDoc quando aplicável. Comentários iniciam pelo identificador; explicam propósito, unidades, restrições, efeitos e erros relevantes. Parâmetros e retornos são explicados no comentário da função. Identificadores podem permanecer em inglês; o texto explicativo de GoDoc, comentários locais, docstrings e Javadoc gerado deve estar em português.

Toda variável local declarada deve ter seu papel documentado próximo à declaração; o comentário não precisa repetir o nome nem explicar sintaxe trivial. Parâmetros de closures e variáveis de iteração seguem a mesma exigência de finalidade. Comentários locais não aparecem em `go doc`. Código de teste e exemplos também precisam de clareza de responsabilidade e variáveis.

Fontes Java geradas devem conter Javadoc em classes, interfaces, construtores e métodos, além de comentários de campos, constantes e variáveis locais. Conteúdo fornecido pelo usuário deve ser escapado inclusive quando inserido em comentários. Nenhuma documentação pode afirmar comportamento inexistente.

### V. Compatibilidade demonstrada e testes proporcionais

Testar invariantes, persistência, navegação, temporização e contratos de saída. Usar testes de referência para artefatos determinísticos, testes de integração para fronteiras e homologação visual nas ferramentas reais. Testes com dublês EFL não substituem integração Java/EFL; validação de XML não substitui execução do PPT no PowerPoint.

Toda entrega funcional deve ter requisitos, tarefas e evidências rastreáveis. Testes de regressão de um defeito devem falhar antes da correção quando tecnicamente viável. Não criar testes que apenas repitam detalhes internos sem verificar uma regra ou contrato.

Antes da aprovação funcional, disponibilizar uma revisão executável do incremento
e um roteiro de teste manual com passos, resultados esperados e ambiente conhecido.
Registrar em quickstart.md a revisão/build, executor, resultado observado e evidência
dos cenários aplicáveis. Testes automatizados aprovados não substituem essa validação.
A aprovação final exige manifestação explícita do usuário sobre a revisão testada;
silêncio, aprovação do plano e autorização para implementar não constituem aceite.

### VI. Simplicidade sustentada por trade-offs

Um módulo Go, composição manual e persistência em arquivo são as escolhas iniciais. Não adicionar microserviços, banco, framework de DI, motor de plugins executáveis ou servidor web sem novo requisito e decisão registrada em research.md. Usar biblioteca padrão quando suficiente. Dependências devem ter versão fixada, licença identificada e custo de distribuição conhecido.

Decisões arquiteturais devem declarar contexto, alternativas rejeitadas, benefícios, custos e condição para revisão. O legado informa o domínio; seus defeitos e sua organização não são requisitos de compatibilidade.

Manter os artefatos do padrão GitHub Spec Kit. Pesquisa, análise do legado e decisões
ficam em research.md; divisão em sprints, dependências e rastreabilidade ficam em
tasks.md. Não criar índices, ADRs ou relatórios de planejamento avulsos quando esse
conteúdo pertence aos artefatos existentes. Recursos necessários às skills seguem
o formato próprio de skills, sem se tornarem artefatos adicionais da feature.

### VII. Especificações sincronizadas e rastreabilidade por domínio

Cada incremento deve atualizar os artefatos correspondentes à sua responsabilidade
antes da entrega para teste: spec.md para comportamento e aceitação; plan.md para
arquitetura e caminhos reais; data-model.md e contracts/ para dados e formatos;
DESIGN.md para alterações visuais; tasks.md para trabalho e evidências; quickstart.md
para revisão executável e roteiro manual. Atualizar somente os artefatos afetados,
sem copiar regras integrais para múltiplos documentos.

Projeto, Catálogo, Tela, Jornada e Demonstração mantêm seus próprios requisitos e
backlogs. A coordenadora referencia a integração. Features de componentes e
adaptadores, como barra, propriedades e prévia, detalham suas responsabilidades
sem criar domínios artificiais. Preferências de tema, seleção e estado da janela
permanecem fora das entidades e dos contratos persistidos da jornada.

Toda mudança deve permitir seguir o requisito qualificado pelo caminho da feature
até a tarefa, o código e a evidência pertinente. Ao mudar uma assinatura, caminho,
versão de formato ou decisão técnica, corrigir as referências afetadas. Distinguir
versão do software, formato do projeto, catálogo e constituição; valores históricos
podem permanecer como histórico identificado, nunca como instrução vigente conflitante.

Uma tarefa deve refletir seu escopo integral. Código e checks suficientes permitem
concluir uma tarefa técnica; homologação e aprovação exigem evidências próprias.
Se uma tarefa mistura implementação e teste manual ainda não executado, separar
explicitamente esses trabalhos ou mantê-la aberta com seu resultado parcial.
Não marcar tarefas em lote por existir código semelhante e não manter trabalho
comprovadamente concluído descrito como exclusivamente futuro sem registrar a situação.

Antes de declarar conformidade entre specs e implementação, revisar os requisitos,
contratos e tarefas afetados e registrar divergências nos artefatos existentes.
Distinguir comportamento não implementado, documentação desatualizada, teste ainda
não executado e capacidade adiada pelo usuário. Não modificar a spec apenas para
apagar uma lacuna funcional: mudanças de escopo precisam ter justificativa e origem
identificadas. Essa regra evita que uma entrega pareça aprovada ou completa por
atualização de texto sem implementação ou validação correspondente.

## Restrições e evolução

- Escopo funcional planejado da primeira versão: sete famílias EFL, navegação explícita, roteiro finito, PPTX editável e Java fonte. Planejamento não comprova implementação ou homologação.
- Entrega corrente: compilar somente Windows x64. Linux/macOS ficam como base de configuração e documentação para uso futuro, sem novos builds obrigatórios por incremento, conforme orientação do usuário em 2026-10-18.
- Java/EFL permanece visível e temporariamente desabilitado por orientação do usuário. Preservar o código e as pendências de integração real; a existência do gerador não comprova compatibilidade EFL nem conclui sua homologação.
- Fora da primeira versão: importação `.page`, integrações externas além das exceções GitHub e MCP local explicitamente delimitadas nesta seção, banco de negócio, colaboração simultânea, gravação de uso e compilação Java pelo editor.
- A integração GitHub da feature `037-integracao-github` é opcional e depende de conexão voluntária da conta e vínculo explícito do workspace. Salvar permanece local; somente a confirmação de publicação autoriza add, commit e push das jornadas elegíveis. Indisponibilidade da rede ou da autorização não deve impedir edição e salvamento local. Credenciais e vínculo Git ficam fora do formato `.jornada`; credenciais não podem ser versionadas. Esta exceção GitHub não autoriza telemetria nem outras integrações externas; a exceção MCP local possui limites próprios no parágrafo seguinte.
- A feature `038-integracao-mcp-agentes` poderá implementar, após aprovação explícita de sua especificação e autorização de desenvolvimento, uma interface MCP local, opcional e iniciada explicitamente para agentes/LLMs externos por clientes compatíveis. A abertura normal do editor não inicia MCP. A primeira entrega usa comunicação local sem portas de rede, tendo `stdio` como transporte candidato sujeito à validação. Editor manual e MCP são modos independentes e complementares: a GUI continua integralmente utilizável offline e sem IA, e o MCP não exige abrir a GUI. Ambos utilizam o mesmo núcleo de Projeto, Catálogo, Tela, Jornada e Demonstração e o mesmo estilo legado definido pelo projeto.
- O aplicativo e o servidor MCP local não dependem de internet e não enviam espontaneamente jornadas, fontes legadas, imagens, fontes tipográficas, logs ou outros dados a serviços remotos. Leituras, alterações e exportações restringem-se aos materiais, destinos e operações autorizados. Persistência exige validação, revisão, proteção contra conflitos e autorização verificável conforme o risco; o agente não pode produzir sua própria prova independente de consentimento. Interpretar linguagem natural e o legado é responsabilidade do agente externo, sem IA embutida obrigatória no executável.
- Um cliente com IA hospedada na nuvem pode transmitir ao seu provedor informações recebidas conforme suas funcionalidades, permissões e políticas. A ausência de transmissão automática pelo Jornada Studio não garante execução local de qualquer cliente. O autor controla os materiais disponibilizados e escolhe o cliente; IA inteiramente local pode ser utilizada por cliente compatível, sem ser requisito inicial.
- A exceção MCP não autoriza servidor HTTP, exposição em rede local/internet, serviço residente, telemetria, SaaS obrigatório, acesso irrestrito a arquivos, execução de conteúdo de projeto, alteração/compilação automática de Go, novos executáveis ou instaladores. Esses acréscimos exigem escopo e decisão próprios, análise de impacto e aprovação explícita. Decisões e evidências permanecem rastreáveis na 038 e nos domínios afetados; salvar localmente não autoriza publicação GitHub.
- Conteúdo de projeto é dado, nunca código executável, template arbitrário ou caminho externo a resolver automaticamente.
- Projeto inválido pode ser preservado como rascunho quando a estrutura for legível; exportação exige validação do alvo. Dados estruturalmente corrompidos não substituem a sessão ativa.
- Logs e diagnósticos locais devem evitar transcrever os textos digitados no projeto. Telemetria e atualização automática pela rede não fazem parte da primeira versão.

## Processo e critérios de conclusão

O planejamento de referência segue as oito sprints da especificação; cada entrega deve registrar seu estado efetivo. Antes de desenvolver uma capacidade, revisar seus contratos e testes de aceitação. Antes de concluir a tarefa, verificar documentação, limites entre módulos e validação pertinente. Cada sprint produz uma demonstração e um registro de evidências com versão do editor e ambiente.

Uma tarefa de implementação só pode ser marcada como concluída após existir implementação e evidência. Uma entrega exclusivamente documental não conclui tarefas de código; uma entrega com implementação deve indicar a revisão e as verificações realizadas. Dependências externas ausentes devem aparecer como condições de homologação, com impacto e responsável funcional identificado, sem inventar APIs.

Separar conclusão técnica de teste manual e aprovação. Uma versão pronta para teste
manual deve ter build ou harness utilizável, comando/caminho real e checks pertinentes
executados. Teste não executado permanece pendente; não copiar resultado esperado
como resultado observado nem alterar checklists de revisão para contornar um gate.

Falhas manuais devem retornar à correção vinculada ao requisito e ao cenário, com
regressões e reteste do fluxo afetado. Mudanças relevantes após o teste exigem nova
identificação de revisão e reteste dos cenários afetados antes do aceite atualizado.
Aprovar um domínio ou sprint não aprova automaticamente o produto inteiro.

Registrar o estado de entrega em quickstart.md: em andamento, bloqueado por condição
concreta, pronto para teste manual, teste manual reprovado, pronto para aprovação ou
aprovado. Não manter um backlog paralelo. Aprovação funcional não autoriza por si só
publicação, deploy, push, merge ou comunicação externa.

## Governance

A revisão de conformidade deve incluir sincronização por domínio conforme o princípio
VII, além de arquitetura, documentação em português e evidências. Constatar uma
inconsistência exige corrigi-la no fluxo autorizado ou registrá-la como pendência
localizável; não declarar conformidade integral enquanto houver divergências conhecidas.


Esta constituição prevalece sobre recomendações genéricas de skills. As decisões de arquitetura e DI já foram dadas pelo usuário ao aprovar o plano; não é necessário reabri-las por instrução genérica de uma skill. Adaptações para desktop prevalecem sobre recomendações de serviços, como processos sem estado ou configuração exclusivamente por variáveis de ambiente.

Alterações exigem justificativa, atualização dos contratos afetados e análise de impacto nas tarefas. Mudança incompatível de princípio incrementa versão maior; novo princípio compatível, versão menor; esclarecimento, correção. A validação é feita por revisão e checks automatizados pertinentes, sem adicionar uma etapa de aprovação por alteração rotineira.

**Version**: 1.4.0 | **Ratified**: 2026-10-17 | **Last Amended**: 2026-10-25
