# Contrato MCP v1 — autoria e descoberta

**Estado:** contrato normativo com implementação parcial. Artefatos 1.0.0.17/18 são históricos. Fontes 1.0.0.19 registram descoberta/autoria/validação e os três controles de sessão (23 ferramentas); revisão/staging de jornada estão implementados internamente, com exposição condicionada a T035/T038. Consentimento, gravação R2 e biblioteca continuam pendentes. Estado executável e testes próprios no quickstart. **Perfil:** MCP 2025-11-25, stdio UTF-8/JSON-RPC; SDK oficial Go v1.8.0 fixado. O servidor restringe versões explicitamente; não anuncia revisões novas por herança do SDK.
**Referências:** [spec](../spec.md), [plano](../plan.md), [consentimento](consentimento.md), [biblioteca](biblioteca.md), [propostas](propostas-legado.md).

## Processo e escopo

Linha instalada na fundação: `JornadaStudio.exe --mcp --workspace <raiz>`; `--legacy-root <raiz>` e `--import-root <raiz>` repetíveis por finalidade; `--approval-profile <arquivo-do-operador>` é futuro e será recusado enquanto não houver suporte ao verificador. `--mcp` é exclusivo das flags GUI/exportação e não abre console/Fyne, splash ou preferências. Configuração de aprovação e chaves públicas são carregadas pelo operador no startup, nunca por ferramenta MCP. Sem perfil qualificado, R2 indisponível.

Somente stdin/stdout do protocolo; stderr para falha/log sem conteúdo pessoal. Sem portas HTTP, sampling, rede necessária ou execução arbitrária. EOF/cancelamento libera sessões/locks; publicação já iniciada exige resultado/recuperação, não abandono silencioso.

Usar IOTransport do SDK ligado a stdin/stdout com limite de leitura e writer que bufferiza uma linha completa até LF (incluído no orçamento de 1 MiB). Excesso não pode produzir bytes parciais no stdout; fechar conexão com erro limitado em stderr. Resultados são paginados/pré-limitados para evitar atingir essa guarda; resultado de commit com resposta perdida segue journal/idempotência. IDs JSON-RPC recebidos são inteiro decimal seguro sem fração/expoente, ou string de até 128 bytes UTF-8 após decodificação; entrada fora dessa política é recusada sem efeito.

Raiz física existente e autorização são fixadas no startup. Jornada é nome regular da raiz, validado pelas regras Windows existentes; IDs de sessão/revisão são UUID v4 e nunca caminhos. Ferramentas não alteram raízes, perfil, chave pública ou política. Resources utilizam a mesma autorização.

## Tipos e validação comum

- `sessionId`, `requestId`, `reviewId`, `screenId`, `componentId`: UUID v4 canônico minúsculo. IDs referenciados devem pertencer ao documento/sessão em questão.
- `expectedRevision`: string decimal sem sinal, intervalo uint64; o adaptador converte para o tipo Go. Representação em string evita perda em clientes JSON/JavaScript. Versão do manifesto permanece número conforme codec existente.
- Coordenadas/dimensões: inteiros lógicos; `x/y >= 0`, `width/height > 0`, com clamp/validação iguais ao editor. Campos não enviados permanecem; `null` não representa apagar genericamente.
- Textos de conteúdo: até 65.536 pontos de código Unicode válidos por campo, em GUI, MCP, biblioteca e proposta; nome até 256 pontos de código, preservando limites menores específicos dos contratos/campos. Não contar bytes, unidades UTF-16 ou grafemas como caracteres; não normalizar/truncar silenciosamente. Rejeitar Unicode inválido. Limites agregados de frame/pacote continuam em bytes e são independentes. Strings vazias somente onde o domínio permite.
- Kinds: `button`, `label`, `title`, `textField`, `radio`, `report`, `numericKeyboard`, `barcodeInput`, `barcodeProcessed`, `applicationMenu`. Skin é nome compatível ofertado, não Kind.
- Objetos são fechados: recusar campos desconhecidos, números não finitos, IDs duplicados, JSON com chave duplicada e valores fora do contrato. Schema não substitui validação de domínio nem autorização.

A guarda de frame valida UTF-8 e pares de escapes surrogate antes do decoder/SDK, evitando substituição silenciosa por U+FFFD de Unicode inválido. Depois da decodificação, contar pontos de código no validador comum de novas entradas/alterações. Leitura de documento histórico não trunca conteúdo; diagnóstico de limite não permite apagar ou normalizar dados antigos.

Cada mutação R1 recebe `sessionId`, `requestId`, `expectedRevision`, salvo criação/abertura antes de existir sessão. Falha preserva documento/revisão/histórico. Edição do documento incrementa sua revisão e retorna IDs criados/remapeados. Propostas anexadas e operações de controle têm as exceções de revisão/cota abaixo: não alteram Project.Revision apenas por existir novo rascunho de proposta. Repetição de edição com mesmo requestId/argumentos devolve resultado; mesmo ID com argumentos diferentes falha.

Criação/abertura usam cache por processo antes de existir sessionId. Preparar revisão é controle R1, sem consumir os 1.024 IDs de edição. Recebe também `controlGeneration` corrente; aceitar uma nova solicitação avança a geração, invalida revisão anterior e gera requestId própria do commit. Identidade de controle é `(sessionId, controlGeneration recebida, requestId)`; a de edição permanece `(sessionId, requestId)`. Replay da preparação ainda corrente devolve o mesmo resultado, sem renovar prazo; geração aposentada recebe `controlRetired` sem efeito, e não exige conservar indefinidamente a resposta completa. A chamada de salvar usa requestId do Review/journal R2, não a ID de preparação.

## Ferramentas de I1

Em cada linha, somente os campos indicados são aceitos, além dos campos comuns de mutação quando aplicáveis. `?` indica opcional; demais campos são obrigatórios. A oferta real só ocorre após implementação/teste; tabela é desenho, não inventário do EXE atual.

| Ferramenta / classe | Entrada específica | Saída específica |
| --- | --- | --- |
| `aplicacao_capacidades` / R0 | objeto vazio | applicationVersion, schema/catalog suportados, profileVersion, tools e capacidades por alvo/família, limits, persistentWrites, approvalProfile e motivos de bloqueio |
| `catalogo_listar` / R0 | cursor?, limit? | itens com Kind/nome/dimensões, operações suportadas e próximo cursor |
| `catalogo_descrever` / R0 | kind | propriedades tipadas, defaults/enums/restrições, skins disponíveis e capacidades por alvo |
| `catalogo_listar_presets` / R0 | cursor?, limit? | presets instalados, dimensões/restrições |
| `catalogo_listar_skins` / R0 | kind, cursor?, limit? | nomes/estados/recursos disponíveis e diagnósticos auditáveis, sem oferecer definição quebrada |
| `workspace_listar_jornadas` / R0 | cursor?, limit? | filename, nome autorizado, diagnóstico por entrada ilegível; sem percorrer subpastas |
| `jornada_criar` / R1 | requestId, name, mode=`blank` ou `example` | sessionId, revision, resumo e estado de rascunho, sem arquivo novo |
| `jornada_abrir` / R1 | requestId, filename | sessão com snapshot e base persistida, sem substituir GUI |
| `jornada_obter` / R0 | sessionId, screenId?, componentId?, cursor?, limit? | snapshot/resumo filtrado, revisão e diagnósticos/paginação; não retorna pacote em base64 |
| `tela_criar` / R1 | name, width?, height?, backgroundColor? | screenId; fundo/defaults e FileID do mesmo caso de uso manual |
| `tela_configurar` / R1 | screenId, changes | mudanças fechadas: name, fileId, width, height, backgroundColor; inválidas não reposicionam silenciosamente os controles |
| `tela_duplicar` / R1 | screenId | novo screenId e mapa de IDs conforme operação existente |
| `componente_inserir` / R1 | screenId, kind, x, y | componentId com defaults/capacidades reais |
| `componente_atualizar` / R1 | screenId, componentId, changes | mudanças fechadas e validação descritas abaixo |
| `componente_mover` / R1 | screenId, componentId, x, y | geometria resultante, mesma regra manual |
| `preset_inserir` / R1 | screenId, presetId, x, y | componentIds; inclui decoração/recursos pelo caso de uso comum |
| `componente_remover` / R1 | screenId, componentId, referencePolicy=`reject` ou `removeOwnedReferences` | referências afetadas e IDs removidos; política não autoriza apagar arquivos |
| `jornada_definir_inicio` / R1 | screenId | início confirmado no rascunho |
| `navegacao_vincular` / R1 | sourceScreenId, sourceComponentId, targetScreenId?, action=`navigate` ou `none` | vínculo; action none exige ausência de target e remove ação conforme domínio |
| `jornada_validar` / R0 | sessionId, target=`draft`/`preview`/`pptx`, cursor?, limit? | diagnósticos localizados por alvo; validar PPT não exporta |
| `jornada_revisar_alteracoes` / R1 controle | filename, intent=`create` ou `replace`, allowDraft, controlGeneration, além dos campos comuns | Review e geração nova, diferenças/challenge; sem persistência de jornada |
| `jornada_salvar` / R2 | sessionId, requestId, reviewId, authorization | resultado do commit/journal; não aceita novo conteúdo ou destino diferente do revisado |
| `sessao_obter` / R0 | sessionId | revisão, documentDirty, proposalDirty, dirty agregado, controlGeneration, quotas/reservas, expiração e capacidades |
| `sessao_descartar` / R1 controle | campos comuns, controlGeneration | invalida propostas/reviews e restaura base conhecida sem gravar; sessão nova sem base volta ao estado inicial |
| `sessao_encerrar` / R1 controle | campos comuns, controlGeneration, discardUnsaved | fecha/libera memória; trabalho pendente sem discardUnsaved=true é recusado com estado preservado |

`sessao_descartar/encerrar` podem descartar trabalho transitório conforme política informada; não são prova de consentimento para R2. Uma sessão expirada retorna erro identificável; nenhuma chamada a recria implicitamente.

Admissão retorna expiresAt inicial UTC/RFC3339, idleSeconds=1800 e unsavedPolicy="explicitDiscardRequired;idleOrEOFDiscardsWithoutSaving". Replay conserva essa resposta original; sessao_obter informa prazo corrente. A consulta devolve quotas {editRequests,proposals,commitPositions,consumedCommits,uncertainIntents,reservedTargets,stagingReservedBytes}, processUsage e capacidades reais. Controle retorna {sessionId,revision,controlGeneration,closed,discarded}, conservados no replay sem consultar estado posterior. Sessão fechada aceita somente o replay terminal conhecido de fechamento com mesmo pedido/geração/digest, sem reabrir conteúdo; outras chamadas retornam sessionExpired. Descarte restaura a base inicial e produz revisão monotônica, remove anexos/review e libera reservas, mantendo o cache de edição para impedir reaproveitamento de IDs. Sessões ociosas são varridas a cada 30 s; acesso após o prazo já recusa imediatamente, sem depender da varredura. Intenção incerta bloqueia descarte/fechamento e preserva reserva até reconciliação.

### Mudanças tipadas de componente

`changes` aceita somente `name`, `fileId`, `enabled`, `bounds`, `style`, `properties`. `bounds` é atualização integral de x/y/width/height, e `style` aceita campos reais de `screen.Style`: fontFamily/fontSize/textColor/fillColor/borderColor/backgroundAssetId/textAlign/skin. Recursos precisam pertencer ao documento; não se aceita caminho externo. Cor/fontes são validadas no domínio/apresentação aplicáveis; `fillAlpha` inexistente não é inventado.

`properties` contém discriminador kind igual ao componente e patch fechado por família:

| Família I1 | Campos permitidos |
| --- | --- |
| button | text; variant em confirm/cancel/yes/no |
| label/title | text |
| textField | initialValue; format em text/number/cellular/cep/cpf/date/ddd/nis |
| radio | options (lista integral `{id,text}`), orientation vertical/horizontal, selectedOptionId; preservar IDs válidos |
| report | title; rows (lista integral `{id,text}`) |
| numericKeyboard | targetComponentId da mesma tela, vazio para desvincular |

Em I1, três Kinds legados aceitam inserção/defaults e alterações comuns; `properties` avançadas são recusadas com capabilityUnavailable até contratos/testes próprios. Isso deve estar explícito em descoberta, sem anunciar edição avançada inexistente. A GUI mantém seus recursos existentes. O DTO MCP não atravessa a aplicação como map livre: adaptar para structs por família.

## Limites iniciais normativos de implementação

| Recurso | Limite |
| --- | ---: |
| Frame de entrada e resposta MCP | 1 MiB UTF-8, antes de decodificar/alocar conteúdo sem limite |
| Profundidade sintática JSON do frame | até 64 níveis; preserva a árvore de menu permitida de 16 níveis e evita recursão ilimitada |
| Sessões por processo / chamadas simultâneas | 8 / 4; uma mutação por sessão, concorrência incompatível retorna busy |
| Inatividade / revisão e consentimento | 30 min / máximo 5 min, sem renovação por replay |
| IDs de edição R1 lembrados por sessão | 1.024, sem eviction; esgotamento bloqueia só novas edições de conteúdo, incluindo propostas, mantendo consulta/revisão/finalização/descarte/encerramento |
| Controle de sessão | uma revisão corrente na v1; até 16 registros recentes de controle mais high-water de geração, sem cache ilimitado de revisões expiradas; respostas aposentadas são terminais |
| Journal/commits por sessão | 64 posições entre commits/intents e reservas; reservar uma posição por alvo persistível sujo antes da edição; recovery/replay usa a posição original |
| Orçamento de snapshots/textos/recursos | 256 MiB somados, incluindo histórico e reserva de metadados de finalização; validar antes de admitir edição, sem prometer heap Go exato |
| Finalização pesada | reserva global de 64 MiB e uma preparação encode/verify por vez; outra retorna busy sem consumir controle. Reserva adicional de 1 MiB/sessão para revisão/metadados, dentro dos 256 MiB |
| Staging reservado à finalização | até o limite do pacote por alvo (250 MiB), total 1 GiB/processo; reserva antes da edição/importação que crie alvo pendente, IO em blocos limitados sem clonar todo documento para revisar |
| Consulta | até 100 itens/página e orçamento do frame; cursor vinculado à revisão, inválido se base muda |
| Composição | 256 componentes; sem composição aninhada; limites de menu do domínio preservados |
| Anexos transitórios de proposta | até 32/sessão, incluindo seleção/versão/importação/legado; revisão própria int32 positiva e baseRevision do documento; incluídos na quota de memória/edição |
| Pacote `.jornada` / total descomprimido | 250 MiB; até 1.000 entradas; manifesto 10 MiB; imagem 25 MiB, conforme codec existente |
| Pacote biblioteca | dois arquivos na raiz, `.jcomp` até 250 MiB; definition.json 1 MiB; template e descompressão interna respeitam limites `.jornada` |
| Evidência textual legado | arquivo até 10 MiB; resposta por trechos de até 64 KiB; imagens somente metadados/recurso autorizado e limites existentes |
| Inventário de material autorizado | até 2.000 arquivos/raiz, profundidade 16, sem seguir links/junctions; cursor/limite de consulta, roots separados por finalidade |
| Deadlines de trabalho / pacote | 30 s / 120 s; espera humana não ocupa uma chamada/trava, revisão expira separadamente |

Histórico/undo também conta no orçamento de retenção. O driver não deve solicitar LLM dentro do deadline do servidor. Alterações desses valores exigem atualização de contrato e evidência, não negociação pelo agente. Dados grandes são paginados/resumidos com próximo cursor e contagem; não truncar valor silenciosamente nem transportar imagem/pacote completo para contornar o frame.

### Finalização sem bloquear rascunho admitido

A quota de edição e a de controle são separadas. Depois da 1.024ª edição, novas edições retornam resourceLimit, mas uma revisão nova da última versão continua preparável. Controle de revisão/descarte/encerramento nunca exige novo ID de edição nem nona sessão. A v1 permite uma revisão corrente por sessão; substituir, negar, cancelar ou expirar aposenta a anterior sem gastar posição de commit ou perder a reserva do alvo sujo.

`controlGeneration` é contador uint64 do servidor exposto em sessao_obter. Edição aceita ou novo controle avança a geração e invalida validação/review afetadas. Um replay do último controle com mesmos argumentos devolve o mesmo resultado ou seu estado terminal, sem renovar autorização. Conferir esse registro antes de recusar geração anterior: sua geração recebida já foi incrementada pelo primeiro sucesso. Fora desse replay corrente, gerações abaixo do high-water são recusadas sem recriar trabalho; requestId reutilizada na mesma geração com outro payload falha. Nova preparação deliberada lê a geração corrente e usa nova identidade de controle. Fechamento/expiração não recriam sessão ao receber replay; ID antiga retorna sessionExpired/closed. Não zerar cotas ou reaproveitar IDs para contornar limites.

Antes de admitir mudança que crie ou mantenha trabalho persistível pendente, assegurar uma reserva para cada alvo sujo distinto, reusada nas edições seguintes do mesmo alvo. Invariante: `commits consumidos + intents incertos + reservas de alvos sujos <= 64`. Documento e duas definições sujas são três alvos, não uma reserva global. Proposta de novo Kind que não pode ser publicada continua transitória e explicitamente não persistível; não prometer que sua estrutura será salva no `.jornada`.

Com 63 commits e um alvo pendente, a reserva permite o 64º commit; sem reserva livre, recusar a edição de um novo alvo **antes** de mudar conteúdo/undo. Negação/expiração não consome commit; falha sem efeito comprovado pode devolver posição à reserva do mesmo alvo. Intenção incerta mantém sua posição até reconciliação; esta não precisa de uma 65ª posição. Após consumir as 64 posições, sessão fica selada para edições; os alvos persistíveis já devem estar finalizados por esse invariante. Propostas transitórias não publicáveis continuam consultáveis e exigem descarte explícito ao fechar. Para continuar a autoria, fechar após finalizar/descartar e abrir sessão nova sobre base persistida, sem transferir autorização antiga.

Reserva de memória/staging/controle deve ser contabilizada antes da admissão, reaproveitando snapshots imutáveis e arquivos congelados; não duplicar documentos integrais ao preparar revisão. Estar com oito sessões ou quota de edição cheia não remove a reserva das existentes. Chamadas podem retornar busy durante operação com deadline, mas finalizar/recovery não depende de vaga de edição/commit nova. Não segurar lock enquanto aguarda humano. EOF/expiração não passam a salvar rascunho implicitamente; journal recupera publicação, não conteúdo transitório.

Operações persistentes sem edição prévia, como exportar ou retirar da oferta, reservam sua posição/recursos na preparação antes de congelar/revisar; falta de quota recusa essa operação sem consumir reservas dos alvos já sujos. Cache de criação/abertura anterior à sessão também entra no orçamento de retenção do processo, sem eviction que permita repetir efeito; quando não houver recurso, recusar nova criação/abertura preservando o caminho de finalização das sessões existentes.

## Resultados, erros e resources

Sucesso estruturado: `{ok:true, revision?, result}`. Falha de ferramenta: `isError=true`, conteúdo estruturado `{ok:false,error:{code,message,location?,retryable,statePreserved,recoveryId?}}`. Violações JSON-RPC/esquema usam erro de protocolo apropriado; erro de domínio não é uma exceção interna exposta ao modelo.

No SDK v1.8.0, fornecer `CallToolResult.Content` textual explícito e curto (até 1 KiB), conservando objeto somente em StructuredContent. Não deixar Content=nil em resposta grande: AddTool pode duplicar o JSON como TextContent. O limite vale para o frame JSON-RPC final, não apenas o objeto de negócio; a revisão máxima deve caber sem duplicação. Resources aplicam o mesmo orçamento e paginação, sem segunda cópia integral em texto.

Codes contratuais: `invalidArgument`, `unknownProperty`, `invalidReference`, `capabilityUnavailable`, `permissionDenied`, `authorizationUnavailable`, `invalidAuthorization`, `authorizationExpired`, `authorizationConsumed`, `staleReview`, `controlRetired`, `externalModification`, `sessionExpired`, `busy`, `resourceLimit`, `canceled`, `persistenceFailed`, `recoveryRequired`. Traduzir diagnósticos existentes sem perder localização/código de origem; mensagens não devem transcrever conteúdo pessoal.

Resources candidatos preservados: `jornada://catalogo`, `jornada://presets`, `jornada://workspace`, `jornada://documento/{id}`, `jornada://diagnosticos/{id}`, `jornada://formatos`. Documento/diagnóstico usa sessionId autorizado, não filename interpretado como URI de filesystem. Paginação e limites iguais aos tools; leitura nunca grava. Resources não são prova de consentimento.

Ferramentas I2 seguem [biblioteca](biblioteca.md); I3 segue [propostas](propostas-legado.md). Menus avançados, jornadas relacionadas, roteiro, simulação e PPT pelo agente permanecem posteriores conforme a spec. Anotações MCP de leitura/idempotência descrevem ferramentas, não concedem autorização nem comprovam ação humana.

I2 declara `componente_revisar` como R1 controle, com entrada ReviewRequest de library-operations.schema.json e saída ReviewResponse de authorization.schema.json; cobre publishDefinition/withdrawDefinition/importDefinition/exportDefinition. `componente_preparar_importacao` é R1 edição, PrepareImportRequest, não controle isento. Os commits R2 usam CommitRequest e não aceitam conteúdo/destino novo. Preparação de jornada permanece restrita a create/replace; consulta componente_validar não cria revisão. I3 publica apenas através de captura I2 vinculada explicitamente por proposal.materialization, nunca executando os statements da proposta.

Materiais de entrada são descobertos por ferramentas R0 paginadas `biblioteca_listar_importaveis` (I2, apenas arquivos `.jcomp` em import-root) e `legado_listar_materiais` (I3, arquivos de evidência em legacy-root). Retornar materialId, referência relativa sanitizada, tipo, tamanho/hash e diagnóstico de indisponibilidade; identidade/hash ficam registrados no servidor e revalidados ao consultar/importar. Nenhuma ferramenta cria nova raiz, segue instrução de caminho de arquivo ou oferece material fora da finalidade autorizada. Falta de raiz/inventário significa capacidade indisponível com motivo, não filesystem irrestrito.

### Detalhes da autoria implementada em T019–T025

jornada_obter retorna items/cursor/total/revision e contexto project. Sem screenId, items são resumos de tela; com screenId, componentes; com componentId, uma entidade da tela. view distingue esses recortes. Estilos/propriedades I1 são tipados, sem recursos binários. Compostos avançados são identificados como summaryOnly/capabilityUnavailable, não expostos como propriedades editáveis completas.

removeOwnedReferences elimina somente bindings cuja origem pertence ao componente removido. Teclado de outro componente, roteiro e vínculo ParentID de filho devem ser resolvidos antes: a ferramenta recusa com invalidReference e preserva estado. reject recusa também binding próprio. Nenhuma política remove arquivos ou entidades externas. Movimento retorna bounds efetivamente confirmado e conserva essa resposta no replay; cópia retorna screenId/idMap, incluindo identidades internas remapeadas.

Metadados de defaults de rádio descrevem optionTexts com IDs gerados por inserção; nunca ofertam UUID de exemplo como identidade reutilizável. Defaults avançados apontam o construtor instalado e não prometem edição avançada. As ferramentas de consulta/rascunho não habilitam R2; persistentWrites continua false.
