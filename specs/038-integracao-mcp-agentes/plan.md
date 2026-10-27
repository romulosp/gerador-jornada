# Implementation Plan: 038 — MCP, agentes e componentes reutilizáveis

**Branch**: `038-integracao-mcp-agentes` (identificador Spec Kit; branch Git atual `main`, sem checkout criado) | **Date**: 2026-10-25 | **Spec**: [v0.4](spec.md)
**Input**: SPEC 038 v0.4, emenda constitucional 1.4.0 e autorização explícita do usuário para emenda e planejamento em 2026-10-25.
**Estado atual**: implementação em andamento, autorizada posteriormente pelos comandos speckit-implement e pedidos de continuidade. Segundo replanejamento U3/U4/I4 e análise foram concluídos antes dessa autorização. G0/fundação/autoria Windows executadas; T015–T025, revisão/staging de jornada e controles de sessão implementados nos fontes. Exposição MCP da revisão depende de T035/T038; R2, biblioteca e paridade visual completa ainda pendentes. [Tasks](tasks.md) e [quickstart](quickstart.md) distinguem código/testes de homologação. Não há autorização para commits, merge ou publicação, nem aceite funcional.

## Summary

Preservar dois modos complementares: o autor abre o editor manual e uma LLM externa utiliza MCP opcional, sem que um modo substitua o outro. Os dois chamam casos de uso compartilhados sobre Projeto, Catálogo, Tela, Jornada e Demonstração; criação, canvas, prévia e PPT aplicável mantêm o estilo legado real.

O mesmo EXE Windows receberá ramificação `--mcp` antes da inicialização desktop, com transporte stdio e sem rede/janela. O SDK oficial Go está fixado em v1.8.0 no adaptador MCP; G0 compila no EXE Windows padrão. A biblioteca declarativa será um adaptador de arquivos sobre definições do Catálogo, instanciadas como cópias independentes; novos Kinds permanecem desenvolvimento Go próprio.

Consentimento segue HostApproval-v1: host confiável apresenta a revisão humana e assina a decisão; servidor verifica chave pública fixada no startup e vínculo com revisão/conteúdo/destino. Cliente genérico sem essa integração opera apenas leitura/rascunho. Nenhum produto externo foi declarado compatível com essa assinatura; entrega completa de gravação depende de homologação com host real qualificado, além do harness técnico.

## Technical Context

**Language/Version**: Go 1.27.2 conforme `apps/gerador-jornada/go.mod`; um módulo. Python 3 somente para scripts de desenvolvimento/verificação existentes, sem requisito no computador final.

**Primary Dependencies**: Fyne v2.8.2 e dependências atuais preservadas; `github.com/modelcontextprotocol/go-sdk v1.8.0` e `github.com/google/jsonschema-go v0.4.3` fixados apenas no adaptador MCP. Usar `crypto/ed25519`, `crypto/sha256`, `encoding/json`, `archive/zip` e `context` da biblioteca padrão. Go mínimo do SDK: 1.25.0; compatibilidade de compilação comprovada em G0 com o toolchain local; homologação de cliente é separada.

**Storage**: `.jornada` ZIP vigente, manifesto 3/4 e combinações históricas; biblioteca privada em `<UserConfigDir>/JornadaStudio/components` ou raiz já selecionada por `JORNADA_CONFIG_DIR`. Preferências/autorizações/journals fora das jornadas e do workspace versionado. O store atual tem envelope/limite próprios: não colocar pacotes grandes em `userconfig.Store`.

**Testing**: testes Go de domínio/contratos/adaptadores, harness cliente MCP Go em testes, falhas injetadas e testes reais Windows/host/Office/DPI. G0 e testes automatizados da fundação executados conforme quickstart; validação humana, DPI/Office e host qualificado continuam pendentes. Verificações de schemas não substituem handlers ou homologação.

**Target Platform**: Windows 10/11 x64, mesmo `JornadaStudio.exe` distribuído. Não criar executável auxiliar, serviço residente ou HTTP.

**Project Type**: monólito desktop com adaptadores manual/MCP e runtime de prévia existente; IA permanece externa.

**Performance Goals**: limites contratuais de [mcp.md](contracts/mcp.md), operações comuns com deadline 30 s, pacotes 120 s, consultas paginadas; validar em fixture real antes de alegar desempenho. Não estabelecer latência de LLM como garantia do aplicativo.

**Constraints**: offline local, stdout exclusivamente MCP, sem imports Fyne/MCP/IO nos domínios; documentação GoDoc em português; efeitos persistentes revisados/autorizados; mesma aparência do legado nos dois modos.

**Scale/Scope**: 8 sessões/processo, 30 min de inatividade, 256 componentes/composição sem aninhamento inicial; até 1.000 itens/16 níveis de menu conforme domínio vigente. Limites de pacote atuais preservados e novos limites explicitados no contrato. I1 autoria existente; I2 biblioteca; I3 evidências/propostas. Operações avançadas da spec continuam futuras, não aparecem em descoberta antes de implementadas.

## Constitution Check

*Gate pré-pesquisa: emenda 1.4.0 aprovada/aplicada em 2026-10-25; v0.4 autorizada para planejamento. Gate pós-desenho: conformidade do desenho abaixo, sem equiparar a evidência de software.*

| Princípio/limite | Decisão verificável de desenho | Estado |
| --- | --- | --- |
| I — domínios explícitos | Catálogo possui definição; Tela instância; Projeto recursos; Jornada vínculos; Demonstração roteiro. Adaptadores traduzem contratos | Atendido no desenho; testar dependências na implementação |
| II — modelo comum | Mesmos casos de uso/representação e instâncias expandidas; protocolo/consentimento fora do manifesto | Atendido no desenho; round-trip futuro |
| III — offline/distribuição | stdio opt-in no mesmo EXE, sem LLM embutida/rede obrigatória; cliente externo pode utilizar nuvem sob sua política | Gate Windows/host real ainda não executado |
| IV — documentação | Novos pacotes/tipos/variáveis/comentários em português na implementação; responsabilidades/caminhos registrados aqui | Obrigação de implementação, não check concluído |
| V — evidências | Harness não substitui usuário; testes de falsa aprovação, race, GUI/MCP e Office/DPI em quickstart | Roteiro definido; execução pendente |
| VI — simplicidade | Um módulo, composição manual, arquivos, SDK oficial; sem DI framework/banco/plugin executável | Atendido no desenho; pin/licenças serão adicionados no código somente com autorização |
| VII — rastreabilidade | Backlogs por domínio gerados; 038 integra sem duplicar novos Kinds; DESIGN 038/023/028 vinculados | 99 tarefas de integração e 16 novas nos responsáveis; marcadores atualizados individualmente por evidência. Correções U1/I1/I2/I3/U2 e U3/U4/I4 registradas; análise prévia sem novos achados relevantes |
| Exceção 038 | Opt-in, sem portas e sem publicação Git implícita; consentimento falha fechado | Autorização constitucional vigente; limites técnicos devem ser comprovados |

Não há violação constitucional de desenho a justificar. Existem gates de viabilidade e homologação deliberados: o plano não os declara resolvidos por documentação. Se um gate falhar, registrar decisão de correção no research; não reduzir requisito para marcar a entrega como concluída.

## Project Structure

### Documentation (this feature)

```text
specs/038-integracao-mcp-agentes/
  spec.md / plan.md / research.md / data-model.md / DESIGN.md / quickstart.md
  checklists/requirements.md
  contracts/
    mcp.md / consentimento.md / propostas-legado.md / biblioteca.md
    authorization.schema.json / library.schema.json / proposal.schema.json
    library-operations.schema.json   # pedidos fechados de revisão/importação/commit
  tasks.md   # execução parcial rastreada individualmente, sem equiparar à homologação
```

### Source Code (apps/gerador-jornada) — caminhos de implementação, com estado em tasks.md

```text
cmd/gerador-jornada/main.go             # seleção de modo antes da composição desktop
internal/application/
  session.go / edit_screen.go / presets.go / document.go  # existentes, reutilizar
  authoring.go / authoring_draft.go / authoring_workspace.go  # registro/fonte existentes
  authoring_screen.go / authoring_component.go / authoring_navigation.go  # casos comuns GUI/MCP
  authoring_validation.go / authoring_diagnostics.go  # validação comum e consulta sem copiar imagens
  authoring_review.go                  # revisão, reserva global e aposentadoria implementadas
  component_library.go / legacy_proposal.go # propostos: casos de uso, sem IO/MCP
internal/catalog/
  capabilities.go                     # metadados/defaults instalados implementados
  reusable.go                         # definições de biblioteca ainda futuras
internal/screen/                       # instâncias e remapeamento, sem modelo paralelo
internal/project/ / journey/ / scenario/ # domínios existentes, contratos preservados
internal/presentation/ / skin/         # cena e biblioteca legado existentes
internal/adapter/
  mcp/                                # SDK/stdio/autoria e controles implementados; R2 futuro
  workspace/                          # fonte extraída e compartilhada com desktop
  workspaceguard/                     # handles/gate/criação implementados; coordenação integral futura
  componentstore/                     # proposto: pacotes imutáveis da biblioteca
  consent/                            # revisão/staging implementados; verificador e journal futuros
  jsonwire/                           # validação estrita de Unicode/JSON comum ao MCP e manifesto
  desktop/                            # paleta, inspector e mesmos casos de uso
  projectzip/ / atomicfile/ / pptx/ / gitworkspace/ # existentes, fronteiras preservadas
```

**Structure Decision**: manter caminhos reais existentes e adicionar somente fronteiras justificadas. Não mover o pacote desktop inteiro, criar `pkg/`, novo módulo, serviço MCP de domínio ou motor de plugins. Interfaces pequenas são definidas pelo consumidor de aplicação apenas para IO, consentimento e teste; tipos de domínio permanecem concretos. A árvore distingue arquivos implementados das extensões futuras.

## Decisões técnicas e fluxo

### 1. Inicialização e transporte

Selecionar flags antes de `desktop.InitApplicationLog`, `app.NewWithID`, splash, preferências e operações GitHub. `--mcp` incompatível com flags GUI/exportação deve falhar em stderr/código de saída sem iniciar outra modalidade. Imports desktop ainda podem executar inicializadores: auditar e testar o processo real; retornar antes de `app.NewWithID` não comprova startup sem efeitos.

Composição disponível: `JornadaStudio.exe --mcp --workspace <pasta-autorizada>`; `--legacy-root` e `--import-root` repetíveis apenas quando explicitamente fornecidos por finalidade. `--approval-profile <configuração-do-operador>` continua futuro em T035. Inventários especializados R0 oferecerão materialIds, sem ferramenta para ampliar raízes. Sem perfil confiável, anunciar `persistentWrites=false` e manter R0/R1.

Usar `mcp.IOTransport` sobre stdin/stdout com MaxLineLength=1 MiB e writer limitado por frame completo antes de encaminhar bytes ao stdout; o wire continua stdio. Essa composição é necessária porque MaxLineLength limita entrada, não saída. Writer guarda até o LF em buffer limitado e não emite frame parcial; excesso encerra a conexão com erro seguro/stderr. Handlers pré-limitam/paginam resultados e usam Content textual curto explícito, evitando que AddTool duplique objetos grandes em texto e StructuredContent. Testar revisão máxima e envelope JSON-RPC completo, incluindo id limitado a 128 bytes.

Perfil MCP 2025-11-25 explicitamente permitido; não anunciar revisão 2026-07-28 até teste próprio. Reconhecer initialize/initialized conforme negociação. Não usar logging remoto/sampling/URL elicitation como dependência; aprovação segue contrato da aplicação. Resources leem IDs autorizados, não caminhos arbitrários.

### 2. Autoria comum e catálogo

`application.Authoring` coordena o registro de sessões e chama operações específicas existentes; não concentra domínio, persistência e widgets em um objeto novo. Mutex por sessão protege `Session.Change`, undo/redo, leitura e revisão; toda resposta usa snapshot independente. Propriedades são traduzidas para structs tipadas por família na borda.

Propostas são anexos da mesma sessão, com revisão própria, baseRevision e dirty separado/agregado; editar anexo não chama Session.Change nem altera Project.Revision/undo. Wrapper I3 exige sessionId/requestId/expectedRevision. Todos os canais validam conteúdo textual em pontos de código Unicode, com os limites específicos de mcp.md/library.schema.json; bytes de frame permanecem independentes. Coleções/defaults de proposta seguem o descriptor fechado e suas limitações em propostas-legado.md.

Metadados de propriedades/capacidades devem ser validados contra defaults/validação reais. I1 oferece inserção dos dez Kinds instalados, atualização comum e propriedades das sete famílias base; edição avançada de menu/barcodes só é anunciada quando tiver handler e testes próprios. `catalogo_listar` une `Definitions()` e `LegacyDefinitions()`, sem confundir skins/partes com Kinds.

Extrair a decoração de `pageTitle` de `desktop.insertPreset` para operação compartilhada com os recursos de `builtinassets`; ambos os canais chamam a mesma inserção integral. Defaults de fundo e estilo também devem ser aplicados por casos de uso comuns. Não adicionar estilo genérico ao MCP.

### 3. Revisão, consentimento e persistência

Produzir bytes finais de gravação antes da revisão e congelá-los em staging privado da sessão; revisão contém hash desses bytes, base, destino e diagnósticos. Mensagens MCP não transportam o pacote inteiro. Declaração `allowDraft` descreve a intenção e não autoriza a gravação: diagnósticos devem constar da revisão humana.

Implementação de preparação: AuthoringReviewFreezer é fronteira de duas operações definida pela aplicação. O adaptador consent usa handle privado Windows delete-on-close, projectzip.Write/Verify em fluxo e JSON canônico. Reserva global de 64 MiB com uma finalização pesada por vez e 1 MiB/sessão dentro dos 256 MiB; staging conserva a cobrança até liberação física. SessionStatus consulta ciclo de vida sem copiar documento; varredura a cada 30 s remove sessões ociosas sem esperar mutex ocupado. Intenção incerta permanece para recovery. Decisão e limites em research D14; nada disso qualifica consentimento.

Escolher HostApproval-v1, verificação Ed25519 da biblioteca padrão e chave pública pinada pelo operador. O host obtém a revisão do servidor e assina decisão somente após ação humana independente. Chave privada, emissor e configuração não são ferramentas do agente. Desenho e bytes assinados em [consentimento](contracts/consentimento.md); rejeitar replay, assinatura inválida, expiração, revisão obsoleta e chave enviada pelo chamador.

Um harness Go de contrato usa chave efêmera e decisões de teste; não é cliente humano homologado nem EXE adicional distribuído. Host externo qualificado deverá demonstrar isolamento das ferramentas do modelo e aprovação humana. Não foi verificado suporte de produto específico a esse perfil. Modelo com shell/fs/processos irrestritos sob o mesmo usuário não satisfaz essa separação e não recebe R2 nesse perfil.

Separar os 1.024 IDs de edição dos controles de revisão/descarte/fechamento. Uma revisão corrente, controlGeneration e registros recentes limitados permitem re-revisar após negar/expirar sem cache ilimitado ou renovar consentimento por replay. Reservar posição de commit **por alvo persistível sujo**, metadados/memória e staging antes de admitir edição; commits/intents/reservas somam no máximo 64. Documento e definições não compartilham uma única reserva. Recuperação usa a posição original; quota cheia não torna o rascunho já admitido irrevisável. Regras completas em mcp.md e decisões D6/D10 do research; capacidade de IO/autorização continuam gates reais.

### 4. Concorrência e falha segura Windows

Extrair gate comum por raiz física para GUI/MCP/Git, reutilizando a estratégia de handle Windows sem compartilhamento já usada pelo Git, com namespace comum em configuração privada. Todo escritor da aplicação deve passar por ele, inclusive criação/importação e fase de aplicação Git; não manter `gitCanWrite` como única proteção. Ordem de trava: raiz → sessão/registro necessário; não executar callbacks Fyne ou esperar consentimento segurando a trava de raiz.

Ao abrir, guardar `{exists, volume/file identity, SHA-256, size}` da base; no commit, revalidar destino e base sob gate e conferir autorização. Criação é exclusiva, sem substituição automática de homônimo. Gravação, consumo de autorização e recuperação usam journal com intenção/hash esperado antes de publicar e resultado depois.

**Limite explícito:** hash seguido de `MoveFileEx(REPLACE_EXISTING)` não é CAS contra escritor externo arbitrário. A implementação deve demonstrar no Windows um método que preserve esse requisito, com handle/precondições e testes de race, ou anunciar substituição indisponível e oferecer novo destino exclusivo. Não permitir sobrescrita desprotegida como fallback. GUI também deve detectar base obsoleta e respeitar a mesma política; Git deve participar da coordenação. Gate de arquivo externo permanece obrigatório para habilitar substituição, não é assumido como resolvido neste plano.

### 5. Biblioteca declarativa e divisão por domínio

Escolher pacote `.jcomp` ZIP v1 com `definition.json` e `template.jornada`. O template é uma jornada válida 3/4 de uma única tela contendo somente os componentes capturados, recursos necessários e vínculos internos. Reutilizar o codec estrito do projeto em vez de inventar outra serialização de dez famílias. Normalizar geometria relativa sem alterar proporção/ordem; rejeitar captura com referências externas não resolvidas na primeira entrega. Isso cumpre FR-023 pela recusa explícita; portas externas ficam ampliação futura, não vínculo acidental.

Definição identifica versão, dependências e parâmetros por caminhos tipados permitidos. Referências de skin/fonte instalada levam identidade/integridade da biblioteca e recursos exigidos; imagens ficam incorporadas no template e são copiadas/remapeadas à jornada de destino. Nenhuma fonte ou código externo é instalado.

Conteúdo publicado imutável em `components/objects/<packageSha256>.jcomp`, com índice separado de oferta/retirada e escrita transacional. Definição/versão com conteúdo diferente colide; retirada altera índice, não blob ou instâncias. Exportar copia pacote integral revisado; importar valida ZIP externo e template interno com limites cumulativos. Contrato e schema em [biblioteca](contracts/biblioteca.md).

Codec T068/T070 e fronteira 002/T018 são pré-requisitos compartilhados da publicação US3 T055/T056, executados antes de congelar bytes/revisar/publicar; a numeração/etiqueta US4 não impõe adiamento até import/export. Store T054 recebe bytes opacos e precede testes/codec; depois vêm captura validada e publicação. UI/transferência US4 continuam incremento posterior. DAG de tarefas explicita essa antecipação sem ciclo.

### 6. Fábrica e aparência

I3 fornece evidências paginadas/referenciadas, sem executar Java. Agente interpreta, entrega proposta estruturada e Go valida as dimensões separadas do contrato/schema; usuário decide equivalência e desenvolvimento. Novo Kind continua no backlog de Catálogo correspondente.

Aplicar o [DESIGN](DESIGN.md) na paleta/propriedades existentes. Incluir biblioteca vazia, importação, conflito, recurso ausente, revisão e retirada. Os componentes instanciados usam a cena e defaults reais do legado tanto pelo manual quanto pelo MCP; verificar fontes, alpha e partes compostas nas saídas. Não alterar o renderer apenas para encobrir ausência de recurso.

### 7. Materialização e revisão alcançável — segundo replanejamento

U3: proposta analítica inclui materialization null ou captura selection/revisão/templateSha256. Autor manual ou agente por I1 constrói componentes reais; I2 captura e devolve identidades/hashes do servidor. Atualização CAS vincula análise à captura da mesma sessão, sem executar statements, duplicar template ou instalar Kind. Validar contra a captura concreta antes de preparar/publicar. Review identifica análise, captura, analysisSha256 e pacote final, com proveniência sanitizada. A captura é um alvo persistível; análise continua transitória após publicação, com descarte explícito e sem segunda reserva para o mesmo efeito. Contrato detalhado em propostas-legado.md; não altera manifesto nem cria domínio MCP.

U4: componente_revisar é R1 controle com ReviewRequest fechado, quatro operations e resposta ReviewResponse existente. Publicar/importar apontam anexo; retirar/exportar apontam definição/versão/hash. Exportação v1 cria filename regular .jcomp no workspace autorizado, sem substituir homônimo. Operador/perfil autorizam a raiz efetiva; nenhum caminho amplia escopo. componente_preparar_importacao admite bytes originais/anexo como R1 edição antes da revisão, com quotas/reservas. Quatro commits usam somente CommitRequest e verificam operation do Review; controle separado não serve para injetar conteúdo novo. GUI compartilha preparação e confirma por evento humano, MCP usa HostApproval-v1. Detalhes e schemas em biblioteca.md/library-operations.schema.json.

I4: T068/T070 continuam antecipadas antes de T055/T056; exemplo de paralelismo não as condiciona a T057. T056/T059 integram revisão/publicação, T071/T072 admissão/revisão/transferência e T079/T082–T087 ligação análise/captura. Não há nova tarefa ou renumeração; critérios de aceitação e Q038-26/Q038-27 verificam as duas lacunas sem fingir execução.

## Sequência de implementação futura e gates de entrega

| Ordem | Escopo verificável | Condição para avançar |
| --- | --- | --- |
| G0 | Prova stdio no EXE padrão, limites, startup sem GUI e perfil de protocolo | Artefato distribuído funciona no harness e cliente real; sem novo EXE/HTTP como atalho |
| G1 | Extração mínima de workspace, gate de escritores e autoria comum | Regressões GUI/Git, rascunho, cópia e aparência preservadas |
| G2 | Ferramentas R0/R1 e propriedades/capacidades | Schemas estritos, validação de domínio e limites comprovados |
| G3 | Revisão/HostApproval/journal e gravação | Harness negativo + humano em host qualificado; substituição só após gate Windows específico |
| G4 | Biblioteca manual/MCP e pacotes v1 | Portabilidade, versões, remapeamento e estados do DESIGN demonstrados |
| G5 | Evidências/propostas do legado | Fronteiras agente/Go/usuário e indisponibilidade de Kind ausente demonstradas |
| G6 | Entrega por incremento | Quickstart com revisão real, comandos, resultado observado e aceite humano; nenhum incremento herda aceite automaticamente |

Essa sequência foi convertida em [tarefas qualificadas por domínio](tasks.md) mediante a invocação explícita da fase pelo usuário. G0 é prova mínima antes da extração; provas de autoria completa se repetem no EXE integrado. Nenhuma tarefa ou teste foi executado e não se impõe implementar todos os incrementos de uma vez.

## Rastreabilidade para a fase tasks

| Escopo | Requisitos 038 | Dono do trabalho futuro |
| --- | --- | --- |
| Adaptador e protocolo | FR-001–FR-005, FR-009–FR-011, FR-016–FR-017, FR-034–FR-039 | 038 integração; sem tipos MCP no domínio |
| Pacote/recursos | FR-014–FR-015, FR-024, FR-026 | 002/Projeto + projectzip; integração 038 |
| Definições/capacidades | FR-003–FR-004, FR-018–FR-021, FR-025–FR-031 | 003/Catálogo; não duplicar expansão de Kinds |
| Instanciação/cópia | FR-002, FR-005, FR-007, FR-022–FR-024 | 004/Tela e casos de uso |
| Navegação/roteiro | FR-006, FR-008 | 005/006; ampliações posteriores mantêm contratos próprios |
| Workspace e escritores | FR-012–FR-013, FR-017, FR-036 | 034/037 e adaptador comum; salvar continua local |
| Manual e visual | FR-002, FR-020–FR-026, FR-032–FR-033 | 023/028/029 e presentation; DESIGN 038 coordena |

## Validation and Integration

Casos e comandos planejados em [quickstart](quickstart.md). Em implementação autorizada: testes focados de domínio, decoder estrito/schema, sessão e replay, publicação/journal, concorrência e escopo; checks Go existentes e importador legado; build somente Windows x64 e roteiro humano no mesmo artefato. Testar todos os defaults sem `Style.Skin`, overrides, fontes, alpha e compostos pela mesma fixture nos dois canais.

Não executar builds/testes de aplicativo nesta fase. O plano resolve escolhas de desenho e identifica gates experimentais, sem alegar que a pesquisa documental prova o resultado no EXE. Tamanho incremental do SDK, memória/latência reais e confiança do host devem constar como resultados observados somente após medição/homologação.

### Histórico da conferência do primeiro replanejamento — 2026-10-25

Os três schemas passaram no metasschema Draft 2020-12 com validador documental temporário fora do projeto. Foram conferidos 52 exemplos em memória: 30 válidos/inválidos de PropertyDescriptor, 18 de limites Unicode em biblioteca/proposta e quatro de vínculos controlGeneration/proposal do Review. Isso não executa nem valida handlers Go, consentimento humano ou persistência real; enum/evidências/ownership continuam exigindo os testes semânticos planejados.

Conferência das 115 tarefas novas: IDs/checkboxes preservados, referências de dependência existentes e DAG sem ciclos, incluindo T055/T056 → T070. Links Markdown relativos da 038 conferidos. Fingerprints dos arquivos sob apps/gerador-jornada e da constituição iguais ao início deste replanejamento; alterações anteriores do workspace preservadas. Nenhum build, teste de aplicativo, commit, merge ou checkout executado. Branch Git permanece main; 038 é o identificador lógico da feature. Não há hooks de plano em .specify/extensions.yml, pois o arquivo não existe.

U1/I1/I2/I3/U2 estão tratados documentalmente em research D6–D10, contratos/modelo/DESIGN, critérios de aceitação e Q038-21–Q038-25, com sincronização mínima de tarefas. Nova invocação de speckit-analyze deve confrontar novamente spec/plan/tasks antes de eventual autorização de implementar. A emenda 1.4.0 já aprovada permanece vigente; nenhuma aprovação constitucional adicional é solicitada neste desenho. Host real, startup stdio no EXE Windows, proteção contra escritor externo, memória/staging e paridade visual continuam provas futuras.

## Complexity Tracking

### Fechamento do segundo replanejamento — antes da análise somente leitura

Os quatro schemas passaram no metasschema Draft 2020-12. Nesta revisão foram conferidos 43 novos exemplos positivos/negativos de materialization, Review vinculado, quatro ReviewRequest, PrepareImportRequest e CommitRequest, incluindo resolução local do $id de Authorization, campos trocados e destinos externos. São verificações documentais em memória, não testes de assinatura, aplicação ou EXE. As 115 tarefas novas preservam IDs/checkboxes; DAG sem ciclo ou referência ausente e links relativos da 038 conferidos. Fingerprints de apps/gerador-jornada e da constituição iguais ao início desta revisão. Nenhum build, commit, merge, checkout ou alteração Go.

Gates constitucionais pré/pós-desenho mantêm o estado da tabela: desenho compatível com a exceção 1.4.0 já aprovada, evidências de software ainda pendentes. Nenhuma alteração adicional de constituição/autoridade/escopo de engenharia foi proposta. Próxima operação explicitamente solicitada pelo usuário é speckit-analyze; nessa fase todos os artefatos ficam somente leitura e o relatório será emitido na conversa, sem marcar antecipadamente sua conclusão neste arquivo. Branch Git preservada main; identificador lógico 038-integracao-mcp-agentes. Artefatos revisados: spec, plan, research D11/D12, data-model, DESIGN, contratos/esquemas, quickstart Q038-26/Q038-27, checklist e tarefas existentes.

Nenhuma exceção adicional solicitada: um módulo e o mesmo executável. Três fronteiras novas são justificadas por IO/protocolo/consentimento; gate comum elimina duplicação insegura. SDK oficial é preferível a implementar JSON-RPC/MCP manualmente; template `.jornada` reutilizado evita segundo codec de componentes. Não adotar servidor remoto, DI framework, banco ou engine de plugins.

## Estado da parcela de autoria T019–T025

MCP já usa DTOs por operação/família e os mesmos comandos de criação/configuração/navegação/presets do editor. Catálogo centraliza defaults e enums; modelos conservam propriedades tipadas. Cópia remapeia opções/linhas/partes além de componentes. InspectMetadata/Diagnostics dispensam duplicação de imagens nas consultas. O codec possui Write em fluxo com bytes golden preservados, pré-requisito de staging; ainda não há Review congelado completo nem consentimento verificado. O build 1.0.0.17 não representa esses fontes posteriores; evidência da próxima revisão será registrada separadamente.
