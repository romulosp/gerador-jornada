# Feature Specification: Conexão do workspace com GitHub

**Feature Branch**: `037-integracao-github` (identificador documental; nenhuma branch Git criada)

**Created**: 2026-10-24

**Status**: Testes manuais de 1.0.0.12/13 reprovados em 2026-10-25; correção 1.0.0.14 remove a digitação extra e melhora o diagnóstico. Checks aprovados; reteste e aceite humano pendentes.

**Input**: Conectar o GitHub pelo navegador padrão, retornar ao Jornada Studio, listar os repositórios do usuário e vincular um deles ao workspace local. Inicializar Git na pasta selecionada e publicar os arquivos `.jornada`, incluindo jornadas interligadas. Decisões confirmadas pelo usuário: publicação por botão separado de Salvar; integração de repositórios com conteúdo existente; somente conexão de projeto local nesta versão. Escopo da etapa inicial: especificação, revisão de qualidade e alinhamento constitucional; a continuação autorizada implementa a integração e registra suas evidências separadamente.

## Histórico técnico do incremento — 2026-10-24 / código 1.0.0.12

A implementação inclui OAuth loopback/PKCE, listagem REST completa, visibilidade e escrita do repositório, coordenação `application.WorkspaceSync`, cofre Windows versionado e operações Git incorporadas ao executável. A publicação inicial com HEAD/remoto ausentes usa árvore vazia e commit sem pais. Retentativa, integração disjunta e recuperação preservam os OIDs conhecidos.

O painel compacto mantém estado truncado com detalhes no menu ⋮ e uma ação principal conforme conta/vínculo; a lista completa é carregada em background antes da busca global. Isso preserva a cobertura de todos os repositórios autorizados de FR-003, sem prometer paginação incremental visível. Caminhos reais e decisões estão em [plan.md](plan.md) e [DESIGN.md](DESIGN.md).

Testes locais exercitam Git bare real, queda em subprocesso, pacotes ZIP interligados, callback de salvamento Fyne, OAuth/REST simulados e WinCred com token sintético. Cadastro OAuth externo foi fornecido, mas seus atributos presentes não comprovam autorização no navegador, política de branch ou push no GitHub real. Esta seção não marca requisito homologado nem aprovação humana; build, checks globais e evidências aplicáveis devem ser consultados em [quickstart.md](quickstart.md).

## User Scenarios & Testing *(mandatory)*

### Correções exigidas no teste manual de 2026-10-25

O usuário escolheu `C:\Users\RomuloPenha\JornadaStudio\Projetos\Meu projeto`, mas a aplicação criou outra pasta `Meu projeto` dentro dela. Também surgiram `.agents`, `.github`, `.specify`, `.vscode`, `apps`, `specs` e outros arquivos do software. O destino esperado, `romulosp/repositorio_jornadas`, permanece vazio na captura fornecida. A inspeção local encontrou o vínculo e origin apontando para `romulosp/gerador-jornada`; o commit registrado alterou duas jornadas nesse outro destino. Sucesso de transporte para esse vínculo não satisfaz o envio solicitado. Esse resultado reprova o fluxo e exige confirmação inequívoca do repositório antes do vínculo e indicação do destino na conclusão.

O produto salvo é o pacote `.jornada`, que contém o projeto e seus recursos necessários para reconstruir a visão editável ao reabrir. A expressão "foto" descreve esse estado salvo; não introduz captura raster nem dependência do código-fonte do software. Cada jornada mantém seu arquivo; vínculos entre jornadas conservam os nomes. A integração não deve copiar nem exportar o código-fonte do Jornada Studio para o workspace.

No reteste de 1.0.0.13, o usuário rejeitou a segunda digitação do nome após selecionar o repositório. A tentativa aceita no log falhou durante o vínculo, antes do push, com erro genérico insuficiente para determinar sua causa exata. A 1.0.0.14 usa diretamente a seleção e abre revisão após conectar; colisão entre jornadas mantém o bloqueio e ganha código/fase de diagnóstico, sem substituir arquivos.

### User Story 1 - Conectar a conta e escolher o repositório (Priority: P1)

Como autor com um projeto local aberto, quero conectar minha conta GitHub pelo navegador e escolher um repositório acessível para trabalhar com esse workspace.

**Why this priority**: Identifica a conta e o destino antes de qualquer envio do projeto.

**Independent Test**: Abrir um workspace local, conectar uma conta de teste e conferir retorno ao aplicativo, identidade exibida e seleção entre repositórios públicos, privados e de organização autorizados.

**Acceptance Scenarios**:

1. **Given** um workspace local aberto e nenhuma conta conectada, **When** o autor aciona “Conectar GitHub”, **Then** o navegador padrão abre a autorização do GitHub sem solicitar senha na aplicação.
2. **Given** a autorização em andamento e o aplicativo aberto, **When** o GitHub confirma a autorização, **Then** o aplicativo retoma o fluxo na mesma sessão, mostra a conta conectada e apresenta o seletor de repositórios, sem exigir copiar códigos ou credenciais manualmente.
3. **Given** autorização negada, cancelada ou não concluída, **When** o autor retorna ao aplicativo, **Then** vê o estado e pode tentar novamente, mantendo o projeto local sem vínculo novo ou publicação.
4. **Given** uma conta com repositórios distribuídos em várias páginas de resultados, **When** o autor percorre a lista, **Then** todos os repositórios acessíveis pela autorização podem ser encontrados sem omissões por paginação, com proprietário, nome e indicação de escrita.
5. **Given** um repositório sem escrita ou uma lista vazia, **When** o seletor é apresentado, **Then** o destino sem escrita não pode ser confirmado e a ausência de destinos é explicada com opção de atualizar a listagem.

---

### User Story 2 - Vincular a pasta local preservando os conteúdos (Priority: P1)

Como autor, quero associar meu workspace a um repositório vazio ou existente, mantendo minhas jornadas e o conteúdo que já está no GitHub.

**Why this priority**: Define a unidade de versionamento e impede que a conexão substitua trabalho existente.

**Independent Test**: Vincular cópias de um workspace com duas jornadas a repositórios de teste vazio e com conteúdo; comparar arquivos e históricos antes e depois, incluindo um caso de colisão.

**Acceptance Scenarios**:

1. **Given** um workspace sem Git, por exemplo `C:\Users\RomuloPenha\JornadaStudio\Projetos\Meu projeto`, **When** o autor confirma um repositório vazio, **Then** o Git é inicializado nessa pasta, o destino fica vinculado e a publicação continua dependendo de ação explícita.
2. **Given** um repositório remoto com arquivos que não colidem com os locais, **When** o vínculo é concluído, **Then** seu histórico é preservado usando a branch padrão e somente os `.jornada` da raiz são materializados no workspace; demais arquivos remotos permanecem no repositório, sem cópia local ou exclusão remota.
3. **Given** uma colisão de arquivos com conteúdos diferentes ou históricos que não possam ser integrados com segurança, **When** a integração é tentada, **Then** ela é interrompida com indicação do impedimento e da ação necessária, sem substituir arquivos locais ou remotos nem concluir um vínculo inválido.
4. **Given** uma pasta que já contém configuração Git compatível com o repositório selecionado e a branch padrão, **When** o autor confirma o vínculo, **Then** a configuração e o histórico existentes são preservados, sem reinicialização destrutiva.
5. **Given** um destino Git diferente, outra branch ativa, operações Git incompletas ou uma pasta contida em outro repositório, **When** o vínculo não pode usar exatamente a raiz selecionada e a branch padrão com segurança, **Then** a operação é bloqueada com diagnóstico sem reconfigurar o repositório existente.
6. **Given** uma falha de rede, permissão ou escrita local durante a integração, **When** a operação termina, **Then** os arquivos anteriores continuam recuperáveis e o aplicativo informa que a conexão do workspace não foi concluída, permitindo nova tentativa.
7. **Given** a pasta escolhida já chamada `Meu projeto`, **When** o autor cria o projeto com esse nome, **Then** o arquivo `Meu projeto.jornada` é salvo diretamente nessa pasta, sem criar `Meu projeto\Meu projeto`.
8. **Given** vários repositórios na listagem, **When** o autor filtra e seleciona um destino, **Then** o seletor mostra seu proprietário/nome completo e conecta exatamente essa identidade, sem pedir para redigitar o nome nem aproveitar seleção anterior; após sucesso, abre a revisão da publicação, sem envio automático.

---

### User Story 3 - Publicar o conjunto de jornadas (Priority: P1)

Como autor, quero salvar localmente durante a edição e escolher quando publicar as jornadas do workspace no GitHub, com uma mensagem de commit revisável.

**Why this priority**: Entrega o envio versionado solicitado sem transformar cada salvamento em uma operação de rede.

**Independent Test**: Editar e salvar duas jornadas interligadas, publicar o conjunto e verificar os arquivos no GitHub e a navegação local. Repetir sem alterações e após uma falha de envio.

**Acceptance Scenarios**:

1. **Given** um workspace conectado, **When** o autor usa Salvar, **Then** o arquivo local é gravado sem iniciar commit ou push, inclusive sem internet.
2. **Given** alterações em arquivos `.jornada` na raiz do workspace, **When** o autor aciona “Publicar no GitHub”, **Then** vê as inclusões, alterações e exclusões elegíveis, revisa uma mensagem de commit editável e confirma a publicação.
3. **Given** conteúdo ainda não salvo na edição, **When** o autor inicia a publicação, **Then** o aplicativo oferece salvar antes de continuar; cancelar ou falhar no salvamento interrompe a publicação.
4. **Given** duas jornadas interligadas salvas e alteradas no mesmo workspace, **When** a publicação é confirmada, **Then** suas alterações são incluídas no mesmo commit e enviadas juntas, preservando nomes de arquivos, destinos e conteúdo dos pacotes `.jornada`.
5. **Given** backups, exportações, temporários, preferências, credenciais ou mudanças de outros arquivos presentes no workspace, **When** a publicação é confirmada, **Then** esses conteúdos não são adicionados ao commit pela ferramenta, mesmo se houver alterações previamente preparadas por outro cliente Git; a publicação é bloqueada se essa separação não puder ser garantida.
6. **Given** nenhuma alteração elegível nem commit pendente, **When** o autor solicita publicar, **Then** é informado de que não há novidades, sem criar commit vazio.
7. **Given** um commit criado e push malsucedido, **When** o autor tenta publicar novamente sem novas alterações, **Then** o mesmo commit é reenviado sem duplicação, e só é marcado como publicado quando sua presença no destino é confirmada.
8. **Given** alterações remotas conflitantes, envio rejeitado, falta de escrita ou branch protegida, **When** a publicação é tentada, **Then** o aplicativo explica o impedimento e preserva arquivos e commits locais, sem force push ou escolha automática de versão.
9. **Given** novidades remotas integráveis sem conflito, **When** o autor publica, **Then** a publicação preserva essas novidades e a história existente; se a integração puder descartar alterações locais, ela é interrompida.

---

### User Story 4 - Retomar o trabalho e desconectar (Priority: P2)

Como autor, quero reconhecer a conta e o repositório em uso, acompanhar publicações pendentes e desconectar sem perder meu trabalho.

**Why this priority**: Permite uso contínuo e recuperação de falhas sem impor nova configuração em cada sessão.

**Independent Test**: Reiniciar a aplicação após conectar, após falha no push e após desconectar; conferir identidade, vínculo, estado de publicação e integridade das jornadas.

**Acceptance Scenarios**:

1. **Given** uma conta e um workspace vinculados, **When** o aplicativo é reiniciado e o workspace reaberto, **Then** a identidade e o destino são lembrados; sem internet, o estado remoto é indicado como não verificado e a edição local permanece disponível.
2. **Given** mudanças locais ou commits ainda não enviados, **When** o workspace é reaberto, **Then** o aplicativo distingue alterações locais, publicação pendente e última publicação confirmada.
3. **Given** autorização revogada ou expirada sem recuperação automática disponível, **When** uma operação remota é solicitada, **Then** o aplicativo pede reconexão e preserva arquivos e commits pendentes.
4. **Given** uma conta conectada, **When** o autor desconecta, **Then** as credenciais mantidas pela aplicação são removidas e as operações remotas exigem nova conexão, preservando jornadas, histórico Git e informações locais do vínculo.
5. **Given** reconexão com outra conta, **When** o autor tenta publicar um workspace já vinculado, **Then** as permissões são reavaliadas para o mesmo destino; a troca de conta não redireciona a publicação silenciosamente.

### Edge Cases

- Retorno de autenticação atrasado, repetido, inválido ou recebido após cancelar: não conectar outra sessão nem repetir o vínculo; orientar nova tentativa quando a sessão não existir mais.
- Aplicativo encerrado durante a autorização: retomar exige nova tentativa nesta versão; não prometer inicialização automática de um processo fechado.
- Nomes com espaços, acentos e diferenças de caixa no Windows: preservar os nomes e bloquear colisões que o sistema de arquivos não consiga representar.
- Falha após o servidor receber o push, antes da confirmação local: verificar o destino na nova tentativa antes de reenviar ou declarar falha definitiva.
- Repositório removido, renomeado, inacessível ou branch padrão alterada: diagnosticar e revalidar o destino antes de publicar, sem enviar para outro destino implicitamente.
- Operações de vínculo ou publicação concorrentes no mesmo workspace: impedir sobreposição e alterações inconsistentes; operações externas detectadas exigem reavaliação antes de prosseguir.
- Jornada inválida recebida do remoto: preservar o arquivo e apresentar o diagnóstico previsto no workspace; não substituir a sessão de edição por conteúdo ilegível.
- Exclusão de jornada ainda referenciada: informar os destinos locais afetados na revisão de publicação; não renomear ou reparar vínculos por inferência.
- Mensagem de commit vazia: solicitar uma mensagem não vazia antes da confirmação.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O sistema DEVE disponibilizar “Conectar GitHub” no contexto do projeto local e realizar a autorização no navegador padrão, sem coletar a senha GitHub na aplicação. Cobertura: US1.1–US1.3.
- **FR-002**: O sistema DEVE retomar a sessão solicitante após autorização válida e tratar cancelamento, negação, expiração e retornos inválidos sem alterar o projeto ou aceitar outra sessão. Cobertura: US1.2–US1.3 e Edge Cases.
- **FR-003**: O seletor DEVE permitir percorrer todos os repositórios acessíveis à autorização, inclusive privados e de organizações, mostrar proprietário, nome e escrita e bloquear confirmação sem escrita. Cobertura: US1.4–US1.5.
- **FR-004**: O sistema DEVE manter uma conta ativa por usuário da aplicação e um repositório vinculado por workspace, identificando conta e destino nas ações remotas. Cobertura: US1.2, US2.1 e US4.1–US4.5.
- **FR-005**: O vínculo DEVE usar exatamente a pasta do workspace como raiz Git, inicializando-a quando necessário e preservando configurações compatíveis. Configuração incompatível DEVE bloquear a operação com diagnóstico. Cobertura: US2.1, US2.4–US2.5.
- **FR-006**: O vínculo DEVE usar a branch padrão do repositório, inclusive no primeiro envio a um repositório vazio, sem oferecer gerenciamento de branches nesta versão. Cobertura: US2.1–US2.5 e Edge Cases.
- **FR-007**: Antes da primeira publicação, o sistema DEVE obter e integrar o histórico remoto sem descartar arquivos locais ou remotos. Somente arquivos `.jornada` da raiz DEVEM ser materializados no workspace; código-fonte, README, pastas e demais conteúdo remoto DEVEM permanecer preservados no histórico/remoto, sem cópia para a pasta local. Colisões entre jornadas e históricos incompatíveis DEVEM interromper a operação. Cobertura: US2.2–US2.3.
- **FR-008**: Salvar DEVE permanecer local e disponível offline; apenas a confirmação explícita de “Publicar no GitHub” autoriza a sequência de add, commit e push. Cobertura: US2.1 e US3.1–US3.2.
- **FR-009**: A publicação DEVE oferecer salvar a edição pendente, mostrar os arquivos elegíveis alterados e exigir mensagem de commit não vazia e editável. Cancelamento ou erro de salvamento DEVE impedir o envio. Cobertura: US3.2–US3.3 e Edge Cases.
- **FR-010**: O commit da ferramenta DEVE abranger somente inclusões, alterações e exclusões de arquivos `.jornada` diretamente na raiz do workspace. Arquivos alheios, mesmo previamente preparados por outro cliente, DEVEM ser excluídos desse commit ou impedir a publicação segura, sem descarte de mudanças existentes. Cobertura: US3.2, US3.4–US3.5.
- **FR-011**: A publicação DEVE preservar os bytes salvos dos pacotes, nomes dos arquivos e destinos entre jornadas, reunindo as alterações elegíveis do workspace no mesmo commit. Referências afetadas por exclusão DEVEM ser informadas sem reparo automático. Cobertura: US3.4 e Edge Cases.
- **FR-012**: Sem alterações elegíveis nem commits pendentes, o sistema DEVE informar que não há novidades sem criar commit vazio. Cobertura: US3.6.
- **FR-013**: Falha de push DEVE conservar o commit local como pendente; nova tentativa sem novas alterações DEVE reutilizá-lo, verificando o destino quando o resultado anterior for incerto. Cobertura: US3.7, US4.2 e Edge Cases.
- **FR-014**: O sistema DEVE preservar novidades remotas integráveis e bloquear conflitos, rejeições, restrições de branch ou perda de permissão com diagnóstico acionável, sem force push, descarte de alterações ou resolução automática de conflitos `.jornada`. Cobertura: US3.8–US3.9.
- **FR-015**: O sistema DEVE mostrar conta, destino e estados distinguíveis de desconectado, conectando, selecionando repositório, integrando, conectado, publicando e falha, além de alterações locais, envio pendente e última publicação confirmada; sem consulta remota recente, não DEVE afirmar sincronismo atual. Cobertura: US1–US4.
- **FR-016**: Conta, vínculo e pendências DEVEM ser retomáveis após reiniciar; revogação, expiração e troca de conta DEVEM revalidar autorização e permissões sem alterar o destino silenciosamente. Cobertura: US4.1–US4.3, US4.5.
- **FR-017**: Credenciais DEVEM ser protegidas no computador e não podem integrar `.jornada`, arquivos versionados, URLs persistidas do remoto ou diagnósticos. Desconectar DEVE removê-las da aplicação sem excluir jornadas, histórico Git ou vínculo local. Cobertura: US3.5 e US4.4; verificação por inspeção dos artefatos e diagnósticos da homologação.
- **FR-018**: Falhas ou cancelamentos de operações remotas NÃO DEVEM impedir edição e salvamento local nem destruir o estado anterior recuperável. Operações no mesmo workspace NÃO DEVEM se sobrepor. Cobertura: US1.3, US2.6, US3.1, US3.8, US4.1 e Edge Cases.

- **FR-019**: O sistema DEVE registrar etapas de autenticação, callback, rede, cofre, seleção e Git com correlação, duração, status HTTP, request ID validado, classificação de causa e informação de DNS/TLS disponível. NÃO DEVE registrar Client ID, Client Secret, tokens, código OAuth, state, verificador/desafio PKCE, cabeçalhos de autorização, URLs com query, corpos livres, conteúdo de jornadas ou mensagens de commit. Erros na aplicação DEVEM apresentar o identificador para consulta no Log de aplicação. Falhas ocorridas apenas na página externa, antes de qualquer callback, devem ser distinguíveis pela ausência de retorno e pelo timeout, sem afirmar acesso ao conteúdo do navegador.

- **FR-020**: Criar um projeto DEVE usar exatamente a pasta escolhida, salvando nela o `.jornada` sem acrescentar automaticamente uma subpasta com o nome do projeto. A pasta contém os pacotes de jornada; os metadados técnicos `.git` são locais, não são conteúdo publicável e não fazem parte do pacote.
- **FR-021**: O `.jornada` DEVE continuar autossuficiente para reconstruir o projeto persistido, incluindo telas, componentes, propriedades e recursos incorporados, sem exigir pastas do software ao reabrir. Jornadas interligadas usam os respectivos pacotes no mesmo workspace.
- **FR-022**: O seletor DEVE mostrar o proprietário/nome completo escolhido e usar essa seleção diretamente, sem exigir redigitar o nome em outra tela. A ação de conectar DEVE levar à revisão da publicação após sucesso, sem publicar automaticamente. A conclusão de publicação DEVE identificar o repositório, branch e commit confirmado no servidor; um commit apenas local ou uma referência remota em cache NÃO DEVE ser apresentado como envio concluído. Vínculo anterior incompatível DEVE ser bloqueado e identificado sem redirecionamento silencioso.

### Key Entities *(include if feature involves data)*

**Incremento de diagnóstico solicitado em 2026-10-24:** logs detalhados integram esta feature, com os limites abaixo.



- **Conta conectada**: Identidade GitHub reconhecida, estado de autorização e permissões efetivamente concedidas; uma ativa por usuário da aplicação. Sua credencial não é conteúdo de projeto.
- **Repositório disponível**: Identidade do destino, proprietário, nome, visibilidade, possibilidade de escrita e branch padrão. A listagem respeita o acesso autorizado.
- **Vínculo do workspace**: Associação local entre pasta existente e um repositório, com a branch padrão utilizada. Não faz parte do documento `.jornada`.
- **Conjunto de alterações**: Inclusões, modificações e exclusões elegíveis das jornadas na raiz do workspace apresentadas antes da confirmação.
- **Publicação**: Mensagem, commit local, destino e resultado do envio, distinguindo pendência, confirmação e resultado incerto. Não confunde salvamento local com envio remoto.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: No cenário autorizado com o aplicativo aberto, o autor conecta a conta, retorna à aplicação e seleciona o repositório sem copiar credenciais ou executar comandos externos; todos os repositórios da conta de teste acessíveis pela autorização são encontráveis, inclusive quando há mais de uma página.
- **SC-002**: No cenário com duas jornadas interligadas, 100% dos arquivos elegíveis selecionados para publicação aparecem no destino com o conteúdo salvo, e a navegação entre elas continua funcionando após reabrir o workspace.
- **SC-003**: Nos cenários de integração sem colisão e de bloqueio por colisão, nenhum arquivo anterior local ou remoto é perdido ou substituído sem decisão do autor.
- **SC-004**: Em todos os cenários offline, de autorização negada/revogada e de envio rejeitado, o autor consegue continuar editando e salvando localmente, com impedimento remoto identificado.
- **SC-005**: Uma falha de envio seguida de nova tentativa sem novas alterações resulta em uma única publicação do commit, inclusive quando o primeiro resultado é incerto.
- **SC-006**: Reiniciar preserva identidade, vínculo e pendências; desconectar remove a credencial mantida pela aplicação sem perder jornadas ou histórico local.
- **SC-007**: Nenhuma credencial, preferência, exportação, cópia de recuperação ou alteração de arquivo alheio às jornadas elegíveis integra os commits produzidos pela ferramenta nos cenários de aceitação.

- **SC-008**: Os testes de diagnóstico DEVEM reconstruir uma tentativa com callback e erro OAuth pelo mesmo identificador, preservando status/código seguro do provedor e comprovando ausência de todos os valores confidenciais sintéticos. Validação real continua separada da simulação.

- **SC-009**: Criar projeto em pasta existente e vincular remoto com código-fonte resulta em zero subpastas de projeto adicionais e zero arquivos remotos alheios materializados. Publicar modifica somente `.jornada` na raiz; abrir os pacotes restaura os dados persistidos e os vínculos.
- **SC-010**: Nos testes com dois repositórios, incluindo um vazio, selecionar e confirmar o segundo publica exclusivamente nele; a conclusão aponta o commit acessível nesse destino. Trocar o filtro invalida a seleção e remoto vazio não confirma um OID apenas em cache.

## Assumptions



- A origem desta spec é o plano aprovado pelo usuário em 2026-10-24. Os requisitos abaixo continuam normativos; presença de código ou revisão documental não representa homologação ou aceite.
- O alvo continua Windows 10/11 x64, com uso local independente de rede. GitHub.com é o serviço desta versão; outros provedores e GitHub Enterprise ficam fora do escopo.
- O fluxo começa com um projeto local aberto. Abrir/clonar um projeto apenas remoto, criar repositórios no GitHub, gerenciar branches, resolver conflitos dentro do editor e colaborar simultaneamente ficam fora desta versão.
- O aplicativo permanece aberto durante a autorização. Retomar a janela existente atende ao retorno solicitado; iniciar um aplicativo encerrado por um protocolo do sistema não faz parte desta versão.
- O workspace e os destinos por nome de arquivo seguem o [contrato de navegação da feature 034](../034-launcher-estrutura-projeto/contracts/workspace-navigation.md). A listagem de jornadas é restrita à raiz; esta feature não introduz subpastas de jornadas.
- A persistência de projeto continua sob o [domínio Projeto](../002-projeto/spec.md) e o [contrato canônico do formato](../001-gerador-jornada/contracts/project-format.md). A conexão GitHub não altera o formato ou a versão de `.jornada` para armazenar Git, caminhos locais ou credenciais.
- Requisitos de integração pertencem à feature 037, com referências aos domínios existentes, sem duplicar seus backlogs. Alterações posteriores de contratos e comportamento desses domínios deverão sincronizar seus artefatos conforme a constituição.
- A autorização efetiva depende do cadastro/configuração da aplicação no GitHub, do acesso concedido e das políticas de organizações e branches. Repositório listado com escrita pode ter publicação bloqueada por política; isso deve gerar diagnóstico, não contorno da proteção.
- O fluxo desktop de autenticação, retorno pelo navegador, armazenamento seguro, autoria, mecanismo Git e dependências estão definidos em [research.md](research.md) e [plan.md](plan.md), produzidos na continuação documental de 2026-10-24. O desenho preserva a distribuição autônoma e não presume que um segredo embarcado permaneça confidencial; a implementação local está presente, mas execução real no GitHub e aceite ainda precisam de evidência.
- Contratos técnicos, [DESIGN.md](DESIGN.md), roteiro manual em [quickstart.md](quickstart.md) e [tarefas rastreáveis](tasks.md) acompanham o planejamento. A geração de tarefas não concluiu código ou homologação por si só; a continuação possui testes técnicos locais. Build/checks globais, execução real no GitHub e aprovação devem manter registros próprios, sem inferir um resultado do outro.
