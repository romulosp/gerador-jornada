# Contrato comportamental — consentimento e operação sem GUI

**Estado:** contrato técnico HostApproval-v1 para revisão do plano; mecanismo selecionado, sem implementação/homologação de confiança.
**Responsável:** integração 038; **requisitos:** FR-009–FR-017, FR-020, FR-030, FR-034–FR-039; **aceitação:** US2, SC-002/SC-008/SC-010/SC-011.
Este documento fixa comportamento e envelope de [authorization.schema.json](authorization.schema.json), com operações/erros/limites de [mcp.md](mcp.md). Não presume que o protocolo ou um cliente genérico ofereça aprovação humana verificável.

## Classes de operação

| Classe | Efeito e exemplos | Política da proposta inicial |
| --- | --- | --- |
| R0 — leitura | Descobrir catálogo, ler jornada/legado autorizado, consultar revisão | Permissão prévia de leitura com escopo explícito; sem leitura arbitrária |
| R1 — transitória | Criar/editar rascunho em memória, propor definição, simular sem persistir | Sessão e limites autorizados; sem gravação durável de documento/definição |
| R2 — persistente | Criar/salvar/substituir jornada, importar, exportar arquivo, publicar/criar versão ou retirar versão da oferta | Revisão e consentimento verificável por operação na base inicial; destino, efeito e conteúdo vinculados |
| R3 — alto impacto/engenharia | Exclusão física de arquivos, apagar definições/dados, mudar fontes Go, catálogo compilado, construir/distribuir EXE | Fora da permissão genérica MCP; fluxo de engenharia/escopo e autorização próprios |

Sobrescrever versão publicada é proibido por FR-025/FR-030; nem aprovação R3 contorna imutabilidade. Retirar da oferta é R2, reversível conforme política futura e sem apagar instâncias; exclusão física não faz parte da primeira biblioteca. Uma política prévia que permita R0/R1 não autoriza R2/R3. Se uma ação combinar efeitos, aplicar o mais restritivo; trocar o nome da ferramenta não reduz sua classe.

Logs operacionais mínimos não tornam edição de rascunho um salvamento implícito: somente metadados permitidos são registrados. Recuperação durável de rascunhos exige contrato próprio; não gravar conteúdo transitório sem informar sua política.

## Fluxo obrigatório

1. **Preparar:** operações R1 alteram apenas o rascunho isolado. Nenhum arquivo definitivo é modificado.
2. **Revisar:** servidor produz revisão imutável com identidade da solicitação/sessão, operação/classe, destino resolvido, base persistida ou ausência esperada, revisão do rascunho, identidade verificável do conteúdo, diferenças, dependências, diagnósticos, impactos e validade da revisão.
3. **Apresentar:** host/canal autorizado exibe ao humano a revisão de forma compreensível. Aprovar, negar e pedir correção são decisões distintas. A revisão vem do aplicativo; resumo livre produzido pelo modelo não a substitui.
4. **Autorizar:** canal controlado pelo usuário comunica decisão por mecanismo cuja origem e proteção possam ser verificadas conforme o modelo de confiança aprovado. Ferramentas disponíveis ao agente não podem produzir a prova de decisão por conta própria.
5. **Conferir:** antes da publicação, servidor confere permissão atual, origem, escopo, validade/revogação/uso da autorização, conteúdo/revisão, destino, base persistida, recursos e conflitos GUI/MCP/Git. Alteração relevante invalida a revisão e exige nova apresentação/decisão.
6. **Publicar:** efeito durável e registro de consumo devem ser coordenados com a política de transação/recuperação; nenhum reaproveitamento da mesma autorização pode produzir outro efeito.
7. **Responder:** retornar resultado estruturado com identidade da operação, estado preservado/gravado, nova revisão e destino permitido. Resposta perdida não autoriza repetição cega.

Solicitar correção invalida a proposta anterior; negar/cancelar impede persistência e conserva o rascunho conforme FR-010. Expirar ou revogar autorização impede a execução. Operações independentes precisam de suas próprias revisões; consentimento não se transfere a outro destino, conteúdo ou sessão.

## Evidência de autorização e limites de confiança

Aceitar somente a decisão verificada pelo mecanismo aprovado no planejamento. `approved: true`, mensagem “o usuário aprovou”, ferramenta de aprovação invocável livremente pelo modelo, token apenas reenviado ou arquivo gravável pelo agente não são prova independente. Um identificador opaco também não é prova por si só.

Um host pode ser parte confiável se sua implementação/versão demonstrar apresentação humana, origem da decisão e vinculação da revisão, com barreira que o agente não possa contornar pelas ferramentas concedidas. Canal separado exige proteção equivalente, não apenas outro nome de arquivo. Política prévia deve ser configurada pelo usuário fora do controle do agente, com escopo e revogação verificáveis; não libera R2 automaticamente nesta proposta.

Processos locais podem compartilhar privilégios de SO. Não se promete isolamento contra comprometimento do host/servidor ou agente com acesso equivalente aos segredos e arquivos do mecanismo. O plano documentará essas fronteiras e os acessos externos do agente. Se a independência necessária não existir, R2 fica indisponível nesse ambiente. Detalhes e alternativas: [research](../research.md#consentimento-headless-alternativas-para-o-planejamento).

## Operação sem GUI

O servidor não cria janela/diálogo Fyne nem depende de eventos da GUI para consentir. A interação humana usa cliente compatível confiável ou canal local independente autorizado. Não iniciar a GUI como fallback; não usar stdout/stdin do transporte como terminal humano. Sem canal adequado, recusar R2 com motivo e conservar R0/R1 dentro do escopo. A abertura normal do software mantém o editor manual disponível, inclusive sem cliente ou LLM.

## Falhas e resultados observáveis

| Condição | Resultado obrigatório |
| --- | --- |
| Mecanismo ausente/não comprovado | Operação persistente indisponível, sem gravação; indicar necessidade de canal aprovado |
| Decisão negada/cancelada | Nenhum efeito durável; preservar rascunho segundo política |
| Conteúdo/base/destino/revisão alterado | Revisão obsoleta; exigir revisão e consentimento novos |
| Expiração/revogação/escopo diferente | Recusar sem consumir outra autorização ou ampliar acesso |
| Evidência forjada/autoaprovação | Recusar e registrar motivo sem expor segredos ou conteúdo do projeto |
| Autorização já consumida | Não produzir outro efeito; solicitação idêntica devolve resultado comprovado, nova solicitação é recusada |
| Falha com resultado incerto | Informar necessidade de reconciliação; preservar dados e evidência de recuperação, sem afirmar sucesso/fracasso não comprovado |

Auditoria local registra apenas metadados necessários: identidade técnica da revisão/operação, classe, origem verificada/política, decisão, validade, estado de consumo e resultado. Não registrar segredo de autorização, credenciais ou texto integral das jornadas. Armazenamento, retenção, proteção e identidade humana exata são decisões técnicas pendentes.

## Cenários mínimos de verificação futura

US2.5–US2.8 e [Q038-01–Q038-10](../quickstart.md) devem provar aprovação real, recusa, correção, falsa aprovação, evidência gravável pelo agente, validade/escopo/repetição, concorrência e cliente sem canal. Esses cenários foram especificados; nenhum foi executado nesta revisão.

## Perfil técnico selecionado — HostApproval-v1

Operador configura no startup `profileVersion=1`, `hostId`, chaves `{keyId,publicKey,allowedOperations,allowedTargetRoots}` e limites de validade. PublicKey é Ed25519 de 32 bytes em base64 canônico. Nenhuma chave pública enviada pelo chamador é aceita; chave privada não é armazenada pelo servidor. O host não é qualificado por nome em `clientInfo`, campo `trusted=true` ou arquivo gravável pelo modelo.

O host elegível mantém a LLM restrita às ferramentas MCP concedidas, sem acesso ao emissor, chave privada ou configuração por shell/fs/processos. Ele apresenta a revisão recebida diretamente do servidor, coleta decisão humana e assina. O operador deve registrar produto/versão/ambiente e demonstrar essa fronteira. Não há produto real já homologado nesta fase. Cliente genérico sem integração utiliza R0/R1; R2 é recusada, não uma gravação automaticamente confirmada.

### Revisão congelada

Servidor produz um objeto Review fechado conforme schema, codificado em UTF-8 por DTO com ordem fixa dos campos `required` do schema; chaves duplicadas/números fora dos intervalos são recusados. `reviewBytes` é base64 dos bytes exatos congelados (até 512 KiB); `reviewHash` é SHA-256 desses bytes. O host decodifica/apresenta esses mesmos bytes e verifica hash/schema; não reconstrói o objeto a partir de um resumo da LLM. Servidor conserva a revisão original; a prova é comparada à revisão guardada, não a um payload substituído pelo cliente.

Review inclui operação, destino físico autorizado, base/ausência, conteúdo final, sessão/revisão, impactos, diagnósticos, host/key permitidos, expiração e desafio aleatório de 32 bytes (base64url sem padding). Sua requestId é a identidade do **commit**, gerada pelo servidor na preparação e devolvida ao cliente; é diferente da requestId R1 que pediu a preparação. O próximo `jornada_salvar` usa a requestId do Review, e o journal R2 mantém a idempotência desse commit.

Review inclui também `controlGeneration` corrente após a preparação, vinculada pelo reviewHash assinado; `proposal` contém proposalId/revision/baseRevision quando o alvo for um anexo, ou null nas demais operações, conforme schema fechado. Identidade/revisão/base entram nos bytes e resumo imutáveis. Uma revisão corrente por sessão na v1; preparar novamente usa a geração exposta em sessao_obter, aposenta a anterior e não consome os 1.024 IDs de edição. Replay corrente não renova validade; geração aposentada retorna controlRetired sem novo efeito. Uma nova revisão exige nova decisão humana, inclusive após negação/expiração.

O campo obrigatório Review.materialization é null salvo publicação originada de análise do legado: nesse caso fixa captureProposalId/captureRevision/templateSha256/analysisSha256 conforme propostas-legado.md. Review.proposal identifica a análise; pacote final e proveniência vêm da captura validada e evidências autorizadas. Captura e análise mudadas invalidam a revisão. A assinatura continua sobre reviewHash e os demais campos existentes; não há novo canal ou autoridade de consentimento.

`componente_revisar` prepara as quatro operações persistentes da biblioteca segundo ReviewRequest de library-operations.schema.json; não reutilizar os argumentos de jornada_revisar_alteracoes para publicar definição. As confirmações componente_publicar/retirar/importar/exportar aceitam CommitRequest e conferem operation do Review. Admissão de importação usa componente_preparar_importacao, edição R1 com quota própria, antes da revisão de controle. GUI chama a mesma preparação de aplicação e confirma por evento humano; MCP mantém HostApproval-v1.

Reserva de finalização segue mcp.md: uma posição por alvo persistível sujo, com `commits consumidos + intents incertos + reservas <= 64`, memória de controle e staging assegurados antes da edição. Negar/expirar não consome a posição reservada; recuperação/replay de publicação usa a posição original. Quota de edição cheia não impede rever/gravar o rascunho já admitido. Indisponibilidade de consentimento, disco ou sistema operacional continua impedindo persistência; reserva de quota não representa autorização ou promessa contra falha externa.

Destino (`targetIdentity`) é identificador estável da operação em raiz física permitida (hash da raiz + nome relativo ou identidade definitionId/version), com localização compreensível no resumo humano. Não transformar string do cliente em caminho. `baseExists=false` exige baseHash=null; revisão existente conserva identidade de volume/arquivo e tamanho no servidor além do hash. Não é permitido regenerar os bytes finais após consentimento.

### Envelope assinado

Host envia Authorization com os campos fechados do schema, assinatura Ed25519 de 64 bytes (base64 padrão), actorId atribuído pelo host e decisão accept/decline/cancel. Assinar os bytes UTF-8 abaixo, uma linha por campo, LF inclusive após a última linha, sem espaços adicionais:

```text
JornadaStudio/HostApproval/v1
<hostId>
<keyId>
<serverInstanceId>
<sessionId>
<requestId-do-commit>
<reviewId>
<reviewHash>
<challenge>
<actorId>
<decision>
<issuedAt>
<expiresAt>
```

Campos são ASCII conforme schema, sem LF/CR; timestamps são segundos Unix UTC em decimal sem sinal/zeros adicionais. issuedAt não pode estar mais de 30 s no futuro e deve estar dentro da validade da revisão; expiresAt deve ser maior que issuedAt e não exceder a expiração do Review nem 5 min depois de issuedAt. Não aceitar valores uint64/revisão acima do intervalo mesmo quando o regex permite 20 dígitos. Validar decodificação base64 canônica/tamanho real; `contentEncoding` do schema não executa essa verificação automaticamente.

Verificar assinatura, chave/host pinados, actorId conforme política do host, todos os vínculos com a revisão guardada, validade/revogação, estado de consumo e base ainda atual. decision diferente de accept nunca grava. Um token/envelope apenas reenviado pelo modelo só é utilizável se essas verificações demonstrarem origem independente; a posse ou o reenvio isolados não são consentimento.

Configuração/chave são fixas por instância. Na v1, revogar um host/chave requer o operador encerrar a instância e reiniciar com a chave removida; o novo serverInstanceId invalida envelopes anteriores. Cancelamento da sessão/instância antes da publicação também invalida revisão pendente. Não oferecer ao modelo ferramenta para alterar a lista confiada. Revogação não desfaz publicação já concluída; resultado incerto segue recuperação, sem sobrescrever dados de terceiro.

### Estados, replay e recuperação

Fluxo: draft → reviewed → awaitingConsent → authorized → committing → committed. Saídas: denied, expired, stale, canceled, recoveryRequired. Editar conteúdo invalida revisão/autorização; autorização não pode ser transferida entre sessões/processos. serverInstanceId muda no reinício; revisão antiga não reaproveita a assinatura.

Antes de gravar, persistir intenção com requestId/reviewHash/contentHash/base/destino e estado de consumo sob coordenação da raiz. Depois, registrar resultado. Falha entre publicação e confirmação exige comparar os bytes publicados com a intenção antes de devolver o resultado; nunca executar outro efeito por replay. Em reinício, journal detecta intenção pendente; novas mutações conflitantes aguardam recuperação. Não restaurar arquivo alterado por terceiro para aparentar rollback.

O harness Go com signer efêmero comprova o contrato e ataques de replay, não interação humana. Teste com host real qualificado e prova Windows/concorrência continuam gates de entrega. Não habilitar R2 em ambiente onde agente dispõe de shell/fs/processos irrestritos sob a mesma autoridade de SO do emissor.
