<!--
Sync Impact Report
Version change: 1.1.1 -> 1.2.0 (novo princípio de sincronização e rastreabilidade por domínio).
Modified principles: nenhum princípio removido ou enfraquecido; preservados I a VI.
Added sections: VII. Especificações sincronizadas e rastreabilidade por domínio.
Updated sections: Restrições e evolução; Processo e critérios de conclusão; Governance.
Removed sections: nenhuma.
Deferred placeholders: nenhum.
Sync impact: somente a constituição alterada; templates e comandos leem as regras em runtime.
Follow-up documental concluído em 2026-10-19, por autorização posterior do usuário:
spec.md/plan.md/data-model.md/tasks.md da coordenadora e dos cinco domínios reconciliados,
com contratos, evidências e gates manuais/externos pendentes explicitados.
Corrigidas as menções a entrega somente documental, formato atual, toolchain e caminhos
antigos; tarefas técnicas/pendentes reconciliadas sem renumeração. A emenda por si só
não conclui implementação nem homologação; evidências ficam em quickstart.md.
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
- Fora da primeira versão: importação `.page`, integração externa, banco de negócio, colaboração simultânea, gravação de uso e compilação Java pelo editor.
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

**Version**: 1.2.0 | **Ratified**: 2026-10-17 | **Last Amended**: 2026-10-18
