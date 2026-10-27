# Feature Specification: Domínio Projeto

## Integração planejada — 038, 2026-10-25

A [038 v0.4](../038-integracao-mcp-agentes/spec.md) foi autorizada para planejamento sob constituição 1.4.0, sem autorização de implementação. Projeto conserva identidade, formato e recursos das jornadas instanciadas da biblioteca; sessão MCP, consentimento e biblioteca pessoal ficam fora do manifesto. Regras específicas e aceitação: 038/FR-014–FR-015/FR-024/FR-026, US4 e SC-005. Esta referência não declara novo formato ou código entregue. Detalhamento em [modelo 038](../038-integracao-mcp-agentes/data-model.md) e [contrato da biblioteca](../038-integracao-mcp-agentes/contracts/biblioteca.md); tarefas deste domínio serão definidas na fase própria, sem duplicar integração.

**Feature Branch**: `002-projeto` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

## Incremento implementado — 2026-10-24 / 1.0.0.4

O writer agora grava **schemaVersion 3/catalogVersion 4** e lê 1/1, 2/1, 3/2, 3/3 e 3/4. `TargetJourney` é preservado no MenuItem; migração antiga acrescenta apenas string vazia. Primeiro save incompatível avisa e conserva backup integral antes de substituir o original. Workspace é a pasta do documento; recentes são preferência local fora do ZIP. Testes de migração/backup e importação sem alterar origem passaram.

Backlog e evidência específicos em [034/tasks](../034-launcher-estrutura-projeto/tasks.md) e [034/quickstart](../034-launcher-estrutura-projeto/quickstart.md), sem duplicar tarefas deste domínio. Homologação humana continua pendente. O estado 1.0.0.1 abaixo é histórico.

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** O writer grava schemaVersion=3/catalogVersion=3; o reader reconhece 1/1, 2/1, 3/2 e 3/3 com os campos efetivamente aceitos pelo codec. Arquivos ZIP guardam configuração e imagens PNG/JPEG; Session conserva histórico separado.

**Pendências e limites:** appearanceProfile, revisões privadas de skins e bytes de fonte por documento não existem no modelo atual. Reconhecer um par de versões não permite ler campos desconhecidos: um arquivo antigo com esses campos será recusado. A matriz completa de falhas/recuperação e portabilidade do perfil permanece pendente.

**Fontes verificadas:** [wire.go](../../apps/gerador-jornada/internal/adapter/projectzip/wire.go), [compound_test.go](../../apps/gerador-jornada/internal/adapter/projectzip/compound_test.go), [skin_wire_test.go](../../apps/gerador-jornada/internal/adapter/projectzip/skin_wire_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preservar a identidade do documento (Priority: P1)

Definir o documento de uma jornada, sua identidade, revisão e recursos, sem conhecer ZIP, janela ou exportadores.

**Why this priority**: O usuário precisa continuar e transportar o trabalho sem perder referências ou depender do navegador.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** um projeto com telas e roteiros, **When** sua identidade, nome ou revisão são consultados após edição e persistência, **Then** os IDs permanecem estáveis e a revisão salva representa exatamente o snapshot gravado.

---

### User Story 2 - Gerenciar recursos incorporados (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um recurso utilizado por duas telas, **When** uma tela é removida, **Then** o recurso ainda referenciado permanece disponível; dados de recurso ausente ou alterado geram diagnóstico localizado.

### Edge Cases

Cobrir IDs duplicados, nome vazio, revisão negativa, recurso compartilhado e isolamento de snapshot. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Representar projeto com UUID estável, nome não vazio, revisão não negativa e coleções ordenadas de telas/roteiros/recursos.
- **FR-002**: Manter UUIDs internos únicos globalmente no projeto; mudanças de nome não alteram referências.
- **FR-003**: Representar imagens PNG/JPEG incorporadas por ID, caminho interno, mediaType, hash SHA-256 e tamanho, sem URLs externas.
- **FR-004**: Permitir documento sem telas como rascunho; sua ausência de inicial deve ser resolvida antes de prévia/exportação.
- **FR-005**: Fornecer visão consistente para snapshot sem compartilhar mapas ou slices mutáveis com edição; autosave e IO pertencem à aplicação/adaptador.

### Key Entities

Project e Asset; metadados JavaOptions são dados do documento, validados pelo alvo Java. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: IDs e ordem preservados em 100% das operações do conjunto de referência; snapshot anterior não muda após nova edição.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-001, FR-009, FR-010; NFR-004, NFR-005. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [project-format.md](../001-gerador-jornada/contracts/project-format.md).
- Não implementa filesystem, ZIP, autosave, undo, UI, decodificação de imagens ou geração de artefatos.
- Evidências da base implementada estão no backlog/quickstart da coordenadora; esta atualização é somente documental e não executa novos testes.

## Persistência de IDs de arquivo

- **FR-006**: Gravar manifesto schemaVersion 2 com fileId de telas/componentes. Aceitar schemaVersion 1 por migração explícita: preencher somente IDs ausentes com sugestões sem acentos e sufixos determinísticos para colisões; preservar UUIDs, nomes e referências. Em v2, exigir valores válidos e recusar duplicidade no escopo. Não renomear IDs explícitos silenciosamente. Leitores antigos devem recusar a versão 2 como incompatível.

## Recuperação nos alvos de distribuição

- **FR-007**: O adaptador desktop deve resolver a pasta de recuperação por usuário com APIs nativas, mantendo o caminho anterior no Windows. Linux usa configuração XDG e macOS Application Support. Falha de resolução deve informar indisponibilidade, sem gravar caminho relativo no diretório corrente; Project/Document não carregam caminhos do sistema operacional.

## Expansão pós-v1 autorizada — 2026-10-19

FR-003/FR-006 registram imagens/formato da base inicial. A extensão abaixo é intencional, sem apagar seu histórico.

- **FR-008**: Representar perfil visual versionado e recursos tipados (imagens/fontes permitidas) com proveniência/hash/licença e snapshot das definições usadas; transportar somente dependências efetivamente usadas e licenciadas, sem caminhos externos, instalação de fontes obrigatória ou execução de conteúdo.
- **FR-009**: Leitor novo aceita schema 1/catalog 1, schema 2/catalog 1 e schema 3/catalog 2 por migração explícita; escritor grava schema 3/catalog 2. Preservar IDs/FileIDs, conteúdo, opções/linhas, vínculos, ordem, bounds e overrides, aplicando nova aparência onde não houver override. Leitura não sobrescreve; primeiro save migrado avisa incompatibilidade e conserva cópia íntegra do original antes de publicação atômica.
- **FR-010**: Validar perfil/fontes/imagens por conteúdo, MIME/tamanho/hash/referências e limites ZIP/safePath também nos novos recursos, sem ampliar allowlist silenciosamente. Snapshot/cópia profunda/reabertura não dependem do acervo original nem de recurso mutável global.

FR-008 conserva o alvo normativo SUBSET de bytes usados e permitidos. A restrição
original de snapshot às definições usadas não descreve Capture corrente: há desvio
arquitetural explícito para snapshot integral de metadados, sem incorporação binária
integral, conforme estado implementado abaixo e contrato compartilhado.

**Estado implementado da continuação, sem alterar FR-008–FR-010:** perfil em
`application.Document.AppearanceProfile *skin.Profile`, não em Project, cópia
profunda e captura/Library imutável com revisão exata/arquivada; perfil opcional
no reader schema3 em transição e não admitido em v1/v2. Capture conserva todos os
metadados da biblioteca original em cópias privadas (definições/fontes/aliases/recursos/
revisões arquivadas inclusive não usadas e todos os estados), sem enriquecer perfis
privados com builtin nem fazer upgrade. Uso só limita diagnósticos e acréscimo de
entrada missing. Desvio arquitetural justificado pelo bug crítico save/reopen →
AddComponent `defaultconfirm`, para inserir famílias/trocar variantes após reopen;
não é incorporação binária integral. Alvo SUBSET de bytes usados e permitidos
permanece, payload selado mais amplo altera o hash. Não acrescenta binários sem
autorização mesmo com JSON allowed;
portabilidade completa licenciada permanece alvo. ResourceData verifica
AssetID/hash/tamanho/MIME ou builtin de bytes exatos. Style/DTOs
Fonts/Aliases/Resources não cumprem ainda toda a normalização appearance/overrides/RGBA
e descritores-alvo. Após correção integral de Capture, conforme main: gofmt/vet PASS,
**15 suítes Go compiladas, não executadas**, bloqueio prévio por política Windows
preservado; **28 Python PASS** (15 importador + 4 packaging + 9 Java 6, erro esperado
de `String.strip` incluído, não homologação EFL) e importer `--check` PASS, mesmo
inventário 214/8/270/47. Distribuição Windows amd64 **gerada e verificada**;
checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos preserva o EXE.
Hash/tamanho e manual pendente no [quickstart](../001-gerador-jornada/quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
Não é aceite nem conclusão integral de T009–T012.

**Aceitação**: v1/v2 com teclado 200 × 280 e personalizações abre com aparência nova sem resize; arquivo original permanece byte a byte até salvar. Salvar v3 preserva cópia antiga e round-trip conserva todas as identidades/referências. TTF corrompido, dependência ausente, licença impeditiva ou falha de backup/publicação são localizados e não substituem dados anteriores. Regras de formato no [contrato compartilhado](../001-gerador-jornada/contracts/project-format.md#expansão-pós-v1--schema-3-e-catálogo-2-2026-10-19); tarefas T009–T012 locais. Java continua pausado.

<a id="codigos-barras-2026-10-23"></a>
## Incremento planejado — persistência dos códigos de barras (2026-10-23)

- **FR-011**: Preservar as configurações tipadas de [031](../031-codigo-barras-efl/spec.md) e [032](../032-codigo-barras-processado-efl/spec.md), incluindo rótulo da entrada, blocos/linhas/seletores/valores e todos os IDs/ordem/estilo. Não persistir código digitado, foco, validação, modal ou estado de captura. Planejar `schemaVersion=3/catalogVersion=3` para os tipos novos, com leitura dos pares históricos 1/1, 2/1 e 3/2, sem prometer compatibilidade do leitor antigo; a versão atual 3/2 não muda nesta etapa documental. Primeiro save que alterar a compatibilidade exige aviso e cópia original íntegra antes da publicação, sem fallback, poda ou enriquecimento de perfis existentes.

**Aceitação**: salvar/reabrir duas instâncias parametrizadas conserva tipos/IDs/ordem/Unicode e valores vazios; runtime de entrada não aparece no manifesto. Campos/enums desconhecidos ou versão incompatível recusam a leitura sem substituir a sessão. Abrir arquivo antigo não o regrava, migração não redimensiona componentes e falha de backup impede sobrescrita. Contrato em [project-format.md](../001-gerador-jornada/contracts/project-format.md#formato-codigos-barras-2026-10-23); implementação planejada em [T013–T014](tasks.md#tarefas-codigos-barras-2026-10-23), sem novo modelo financeiro em Project.

<a id="menu-aplicacao-2026-10-23"></a>
## Incremento autorizado — persistência de Menu aplicação

- **FR-012**: Preservar árvore completa, UUIDs/ordem, categorias/variantes/enabled por item, metadados inertes, destinos explicitamente autorados e escolha inicial de [033](../033-menu-aplicacao-efl/spec.md), incluindo menus vazios. Não gravar caminho percorrido, foco, seleção simulada, saída da prévia ou expansão do inspector. Estender o mesmo alvo schema3/catalog3 de FR-011 para os três novos tipos; leitores históricos 1/1, 2/1 e 3/2 continuam, com novo 3/3 e backup antes de save incompatível.

**Aceitação**: round-trip de raiz e submenus aninhados/financeiro/Bolão conserva todos os dados confirmados, inclusive false explícito e strings vazias/Unicode, sem enriquecer perfis ou criar entradas ZIP. Falha de versão/estrutura/backup não substitui sessão/original. [T015](tasks.md#tarefas-menu-aplicacao-2026-10-23) estende codec existente de T013–T014; não repete implementação do backup/versionamento. Catálogo 3 planejado não é alegação sobre o writer validado atual 3/2.

## Integração com Inicialização 035 e Splash 036 — 1.0.0.9

A abertura por -open passa pela preparação assíncrona e só entrega o documento ao editor depois do vídeo e dos recursos. Falha de leitura preserva o original e impede avanço com diagnóstico. Não existe novo campo persistido; reader/writer continuam 3/4. Modelo do ciclo em 035 e adaptação em 036, com backlog próprio.
