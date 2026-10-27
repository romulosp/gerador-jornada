# Pesquisa e decisões propostas — SPEC 038 v0.4

**Data:** 2026-10-25. **Estado atual:** implementação parcial autorizada após o planejamento; emenda 1.4.0 vigente. G0 Windows/stdio e fundação com testes registrados no quickstart. Decisões anteriores preservadas; nenhuma delas é aceite funcional.
**Base:** [spec v0.4](spec.md), parecer de correção e autorização explícita “Aprovar emenda e planejamento” recebida em 2026-10-25. A constituição 1.3.0 é a base histórica; vigente 1.4.0. Modos manual/MCP complementares, estilo legado e modelos por domínio vinculam documentação e implementação futura.

## Proposta de exceção constitucional

### Contexto e decisão submetida à aprovação

Na base 1.3.0, a constituição excluía integrações externas além da exceção GitHub. MCP local continua sendo integração com cliente externo, sem autorização implícita pela natureza local. O usuário aprovou explicitamente a exceção abaixo e o planejamento em 2026-10-25; o texto foi aplicado à constituição 1.4.0 mantendo os princípios I–VII. Os blocos seguintes preservam a proposta exata aprovada, não uma emenda ainda pendente.

### Texto exato proposto para a seção Restrições e evolução

Substituir somente a linha que começa com “Fora da primeira versão” por:

> - Fora da primeira versão: importação `.page`, integrações externas além das exceções GitHub e MCP local explicitamente delimitadas nesta seção, banco de negócio, colaboração simultânea, gravação de uso e compilação Java pelo editor.

No parágrafo da exceção GitHub, substituir apenas sua última frase por:

> Esta exceção GitHub não autoriza telemetria nem outras integrações externas; a exceção MCP local possui limites próprios no parágrafo seguinte.

Após a exceção GitHub, acrescentar:

> - A feature `038-integracao-mcp-agentes` poderá implementar, após aprovação explícita de sua especificação e autorização de desenvolvimento, uma interface MCP local, opcional e iniciada explicitamente para agentes/LLMs externos por clientes compatíveis. A abertura normal do editor não inicia MCP. A primeira entrega usa comunicação local sem portas de rede, tendo `stdio` como transporte candidato sujeito à validação. Editor manual e MCP são modos independentes e complementares: a GUI continua integralmente utilizável offline e sem IA, e o MCP não exige abrir a GUI. Ambos utilizam o mesmo núcleo de Projeto, Catálogo, Tela, Jornada e Demonstração e o mesmo estilo legado definido pelo projeto.
> - O aplicativo e o servidor MCP local não dependem de internet e não enviam espontaneamente jornadas, fontes legadas, imagens, fontes tipográficas, logs ou outros dados a serviços remotos. Leituras, alterações e exportações restringem-se aos materiais, destinos e operações autorizados. Persistência exige validação, revisão, proteção contra conflitos e autorização verificável conforme o risco; o agente não pode produzir sua própria prova independente de consentimento. Interpretar linguagem natural e o legado é responsabilidade do agente externo, sem IA embutida obrigatória no executável.
> - Um cliente com IA hospedada na nuvem pode transmitir ao seu provedor informações recebidas conforme suas funcionalidades, permissões e políticas. A ausência de transmissão automática pelo Jornada Studio não garante execução local de qualquer cliente. O autor controla os materiais disponibilizados e escolhe o cliente; IA inteiramente local pode ser utilizada por cliente compatível, sem ser requisito inicial.
> - A exceção MCP não autoriza servidor HTTP, exposição em rede local/internet, serviço residente, telemetria, SaaS obrigatório, acesso irrestrito a arquivos, execução de conteúdo de projeto, alteração/compilação automática de Go, novos executáveis ou instaladores. Esses acréscimos exigem escopo e decisão próprios, análise de impacto e aprovação explícita. Decisões e evidências permanecem rastreáveis na 038 e nos domínios afetados; salvar localmente não autoriza publicação GitHub.

### Relatório de impacto da emenda aprovada e aplicada

| Item | Proposta e impacto |
| --- | --- |
| Versão base | 1.3.0, referência histórica da revisão anterior |
| Versão vigente | 1.4.0, aprovada e aplicada em 2026-10-25 como expansão compatível |
| Classificação | MINOR: amplia materialmente o escopo permitido, preservando os princípios e as garantias anteriores. Não é correção editorial/PATCH |
| Condição de revisão | Se a aprovação exigir remover/redefinir de forma incompatível um princípio, reavaliar como MAJOR antes de aplicar; não impor esse resultado silenciosamente |
| Princípios I–VII | Nenhum removido ou renomeado; domínio, modelo compartilhado, distribuição offline, documentação, evidências, simplicidade e rastreabilidade preservados |
| Seções afetadas | Restrições e evolução; relatório de impacto temporário e metadados da versão somente após aprovação |
| Datas de governança | Ratified preservada em 2026-10-17; Last Amended atualizado para a aplicação real em 2026-10-25 |
| Formato `.jornada` | Sem alteração automática de schema/catalogVersion; qualquer evolução terá contrato próprio |
| Interface manual | Biblioteca será utilizável sem LLM; mesmo estilo legado e regras de domínio da autoria MCP |
| GUI/MCP | Adaptadores do mesmo núcleo; nova fronteira de arquivos/consentimento, sem serviço de domínio separado |
| Dependências | SDK, licença, transitivas, empacotamento offline e transporte ainda precisam de pesquisa técnica; nada adicionado agora |
| GitHub | Exceção 037 e separação entre salvar e publicar preservadas; concorrência com escritores deve ser projetada |
| Templates/workflow | Leem constituição em runtime; não foram modificados. Referências vigentes da 038 sincronizadas com 1.4.0; 1.3.0 identificada como histórico |
| Artefatos dependentes | Spec/checklist/research da 038; contratos e modelo devem ser reconferidos; plano/tarefas futuros devem incorporar gates e testes; desenhos 023/028 referenciam o desenho único 038 |
| Riscos | Perda por concorrência, falsa aprovação, cliente remoto recebendo dados, deriva visual e dependências ausentes; tratados como requisitos, ainda não mitigados por software nesta etapa |
| Aprovação | Resposta explícita do usuário “Aprovar emenda e planejamento”, em 2026-10-25. Não autoriza implementação ou aceite funcional |

O relatório temporário no topo da constituição foi atualizado com os efeitos reais e preserva a referência à emenda GitHub anterior. Não registra implementação ou homologação inexistentes. Remoção antes de eventual commit segue a governança; nenhum commit está autorizado nesta etapa.

### Justificativa e alternativas

**Proposta:** exceção local estreita. Permite usar os casos de uso reais por outro adaptador e preserva edição manual independente, regras de domínio e aparência do legado. O custo adicional concentra-se em transporte, consentimento, sessões e concorrência, não em duplicar o editor.

**Alternativas não recomendadas:** tratar MCP como já autorizado por ser local ignora a restrição vigente; liberar genericamente integrações externas amplia desnecessariamente o escopo; adotar HTTP/serviço residente muda a distribuição e a fronteira de exposição; embutir LLM cria dependência dispensável ao editor. A decisão pode ser revista diante de requisito remoto explícito e nova análise.

**Fechamento P0:** aprovação explícita recebida, emenda aplicada e referências desta feature sincronizadas. O conflito de escopo MCP local está resolvido na norma vigente 1.4.0; isso não comprova os gates de software ou autoriza implementar.

## Responsabilidades na fábrica

Agente interpreta e formula hipóteses; Go fornece evidências autorizadas e valida deterministicamente; usuário decide equivalência, autorizações e homologação. Isso evita exigir um intérprete semântico universal de Java ou LLM dentro do executável. O [contrato de propostas](contracts/propostas-legado.md) detalha o comportamento sem escolher parser/SDK nem confundir estrutura válida com equivalência funcional comprovada.

Uma composição ou variante já representável pode ser declarativa. Duas legendas continuam variante Button, mas dependem de suporte específico; ComboBox permanece uma família planejada em 003. Nenhum material Java é executado para validar uma proposta. Novos Kinds exigem engenharia autorizada e entrega compatível pelos dois modos.

## Consentimento headless: alternativas para o planejamento

O [contrato](contracts/consentimento.md) define comportamento obrigatório. O plano seleciona a alternativa A com envelope HostApproval-v1 e pin de chave pública. B/C permanecem alternativas futuras; nenhuma integração de host foi tecnicamente homologada nesta etapa.

| Alternativa | Garantia a demonstrar | Limite/custo |
| --- | --- | --- |
| A — aprovação pelo host/cliente | O host apresenta a revisão ao humano e fornece decisão verificável que as ferramentas do modelo não podem forjar ou ampliar | Depende de cliente/versão concretos; UI de confirmação ou texto no chat não bastam sem verificar a integração |
| B — canal local separado | Usuário decide por canal fora das ferramentas do agente, e servidor vincula a decisão à revisão específica | Mesmo usuário/mesmos privilégios de SO podem permitir adulteração; arquivo gravável pelo agente não é isolamento |
| C — política pré-configurada | Autor configura permissões delimitadas fora do controle do agente; servidor impede ampliação e revogações são respeitadas | Leitura e rascunhos podem usar escopo prévio; automatizar R2 exige política específica, revisão e decisão futura. A base desta revisão exige consentimento por revisão para R2 |

Não depender da GUI Fyne nem imprimir perguntas em stdout/ler respostas humanas no stdin do protocolo. Sem alternativa confiável aprovada, manter R0/R1 autorizadas e recusar R2. R3 não ganha ferramenta genérica de engenharia/exclusão.

### Modelo de confiança a completar

O MCP `stdio` é uma comunicação entre processos locais; ele não isola automaticamente processos com privilégios equivalentes. Os requisitos protegem contra autoaprovação via argumentos/ferramentas e contra revisões obsoletas dentro do fluxo controlado. Não se promete resistência a um agente com acesso irrestrito às chaves, ao host, ao executável ou ao mesmo armazenamento de autorização.

O plano deverá identificar ator humano, host confiado, privilégios de SO e ferramentas externas do agente, fronteira de proteção, armazenamento/verificador e limitações assumidas. Se a separação necessária não puder ser demonstrada no ambiente alvo, o cliente não será habilitado para operações persistentes. Não apresentar token ou criptografia como garantia sem comprovar origem e proteção do emissor.

## Estilo legado, domínio e experiência manual

As orientações do usuário vinculam documentação e implementação. [data-model.md](data-model.md) separa responsáveis existentes, sem novo domínio MCP/IA ou biblioteca paralela. [DESIGN.md](DESIGN.md) é dono dos novos fluxos da biblioteca; 023/028 apenas integram acesso e edição. Default legado, fonte real, alpha e partes internas devem atravessar criação manual/MCP, canvas, prévia e PPT com evidências. Design não substitui teste no executável e no Office aplicável.

## Decisões técnicas do planejamento — propostas para revisão

### D1 — SDK e transporte

**Decisão:** fixar `github.com/modelcontextprotocol/go-sdk v1.8.0` na implementação autorizada e usar stdio no mesmo EXE, com perfil explícito MCP 2025-11-25 e frame 1 MiB. Release e `go.mod` oficiais conferidos em 2026-10-25: SDK declara Go mínimo 1.25.0; o projeto declara 1.27.2. Nenhum módulo foi baixado/adicionado ao projeto ou compilado. [Release](https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.8.0), [go.mod](https://raw.githubusercontent.com/modelcontextprotocol/go-sdk/v1.8.0/go.mod).

**Justificativa:** SDK oficial já oferece servidor, tools/resources e stdio; evita manter implementação própria de protocolo. Versões recentes do SDK suportam também revisão 2026-07-28 com ciclo diferente; restringir o perfil inicial evita assumir que initialize é obrigatório em qualquer versão. `StdioTransport.MaxLineLength` deve ser configurado explicitamente, não herdar limite incidental. [Transporte do tag](https://raw.githubusercontent.com/modelcontextprotocol/go-sdk/v1.8.0/mcp/transport.go), [opções do servidor](https://raw.githubusercontent.com/modelcontextprotocol/go-sdk/v1.8.0/mcp/server.go).

**Custo/licença:** o LICENSE do tag descreve transição Apache-2.0 com contribuições históricas MIT e documentação não normativa CC-BY-4.0; preservar o arquivo real e avisos aplicáveis no inventário, sem declarar licença única simplificada. O módulo declara dependências como jwt/v5, jsonschema-go, segmentio/encoding, uritemplate e pacotes x/*; conjunto ligado e impacto no EXE devem ser medidos na implementação. Usar SDK não cria runtime instalado nem rede obrigatória no usuário final. [LICENSE](https://raw.githubusercontent.com/modelcontextprotocol/go-sdk/v1.8.0/LICENSE).

**Alternativas:** JSON-RPC manual aumenta manutenção; HTTP/new EXE violam o escopo inicial; protocolo mais novo só deve ser anunciado após teste específico. Revisar pin por incompatibilidade/vulnerabilidade comprovada, mantendo versão/licença verificadas.

**Correção da revisão técnica:** usar IOTransport stdio com writer de frame limitado, além do MaxLineLength de entrada. O helper AddTool pode duplicar objeto em StructuredContent e TextContent se Content for nil; revisão de 512 KiB em base64 excederia o frame de 1 MiB com essa duplicação. Definir Content curto explícito, paginação e guarda do frame completo antes de emitir qualquer byte. Teste futuro Q038-19 deve cobrir esse limite. Essa conclusão foi obtida por leitura do SDK, não execução. [Servidor do tag](https://raw.githubusercontent.com/modelcontextprotocol/go-sdk/v1.8.0/mcp/server.go).

### D2 — HostApproval-v1

**Decisão:** host confiável assina decisão humana Ed25519 sobre desafio/revisão exatos; servidor verifica chave pública fixa no startup, escopo, validade e uso único. Escolha de algoritmo usa biblioteca padrão, sem nova dependência. [Go Ed25519](https://pkg.go.dev/crypto/ed25519).

**Justificativa:** separa solicitação do modelo da autoridade de consentimento. Chave privada e emissor ficam fora das ferramentas do agente. O perfil de referência exige host gerenciado com modelo limitado às ferramentas MCP autorizadas, sem shell/fs/processos que atinjam o emissor/configuração. O operador registra a chave e qualifica host/versão/ambiente; `clientInfo` ou flag de autocertificação não habilitam R2. Não se atribui suporte desse perfil a Codex, ChatGPT ou outro produto sem evidência.

**Provas e limite:** harness Go com signer efêmero valida mecânica/negativas, não aprovação humana. Entrega completa exige host real qualificado com humano, no mesmo EXE distribuído. `elicitation` pode apresentar interação, mas não atesta por si só independência da origem. [MCP elicitation](https://modelcontextprotocol.io/specification/2025-11-25/client/elicitation). Cliente sem essa integração oferece R0/R1 e recusa R2. Mesmo usuário com acesso irrestrito às chaves/host é fora da fronteira; nesse ambiente não habilitar gravação. B (canal separado) e C (política prévia R2) não estão selecionadas.

### D3 — Windows e concorrência

**Decisão:** gate entre processos compartilhado por GUI/MCP/Git, identidade/hash da base e criação exclusiva. Substituição requer provider que comprove proteção contra interferência Windows; na ausência dessa prova a capacidade de substituir fica desabilitada e o usuário escolhe novo destino. Nenhum hash+rename é anunciado como CAS completo.

**Evidência local:** `Session` não tem mutex; `gitCanWrite` é guarda da UI; `projectzip.Save` termina em `MoveFileEx(REPLACE_EXISTING)`; gate Git é privado. É necessário unificar escritores e testar races/recuperação. Bloquear escrita/deleção por handle interfere no próprio rename; não assumir solução sem prova. [CreateFileW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew).

**Transporte:** `-H=windowsgui` não proíbe handles herdados, mas GUI pode iniciar sem handles válidos. A implementação Go e a documentação Windows dão base estrutural, não homologação do JornadaStudio. Auditar startup e testar EXE padrão com pipes; binário de console não substitui evidência. [Go linker](https://pkg.go.dev/cmd/link), [Windows pipes](https://learn.microsoft.com/en-us/windows/win32/procthread/creating-a-child-process-with-redirected-input-and-output), [Go 1.27.2 subprocessos Windows](https://raw.githubusercontent.com/golang/go/go1.27.2/src/syscall/exec_windows.go).

### D4 — Biblioteca e recursos

**Decisão:** `.jcomp` ZIP v1 contém manifest de definição e template `.jornada` válido de uma tela; store imutável por SHA-256 e índice de oferta separado. Preservar codec existente, copiar imagens e remapear instâncias, sem vínculo vivo. Referências externas à seleção são recusadas inicialmente. Tipos/propriedades continuam nos domínios; backend de arquivos é adaptador.

**Justificativa:** evita serializador paralelo de dez Kinds e evolução incidental de `.jornada`. Parâmetros v1 restringem-se a propriedades simples permitidas; não instalar skins/fontes/código. Alternativas de template executável, biblioteca viva no documento e pacote por widget foram rejeitadas por compatibilidade/complexidade. Reabrir decisão quando houver requisito explícito de vínculo, fontes externas ou composição aninhada.

### D5 — Limites, sessão e apresentação

**Decisão:** 8 sessões, 30 min de inatividade, 5 min por revisão/consentimento, 256 componentes/composição e deadlines 30/120 s; quotas completas em mcp.md. Sem eviction de IDs que permita repetir efeito. Orçamento de retenção inclui histórico/recursos, sem prometer heap total idêntico à quota.

**Justificativa:** impedir alocação/trabalho sem limite com valores iniciais proporcionais ao desktop; ajustar por evidência posterior. Casos de uso comuns aplicam recursos/defaults do legado (inclusive decoração de preset hoje desktop). A interface manual usa DESIGN 038/023/028 e o modelo por domínio; MCP não desenha outra aparência.

## Replanejamento após análise cruzada — U1/I1/I2/I3/U2

O usuário invocou speckit-plan após a análise. As decisões abaixo corrigem o desenho e sincronizam dependências das tarefas existentes, sem código, mudança de IDs ou aceite automático. D5 mantém seus limites gerais; a política específica de controle/finalização passa a ser a de D6/D10.

### D6 — U1: quota de edição separada da finalização

**Decisão:** 1.024 IDs de edição por sessão, sem eviction; preparar/repreparar revisão, descartar e encerrar usam controle separado com controlGeneration uint64. Uma revisão corrente na v1, até 16 registros recentes limitados e high-water para recusar gerações aposentadas sem novo efeito. Replay corrente devolve resultado sem renovar validade. Reservar uma posição de commit por alvo persistível sujo antes de admitir alteração; commits consumidos, intents incertos e reservas somam no máximo 64. Recuperação usa a posição original. Ao consumir 64, selar novas edições; anexos não publicáveis continuam explicitamente transitórios.

**Justificativa:** classificar preparação como edição esgotável impede salvar a 1.024ª mudança. Uma reserva global também falha quando há documento e várias definições pendentes. Negação/expiração não deve esgotar quota ou exigir nona sessão. Geração delimita retenção sem transformar replay antigo em nova autorização; contrato não promete conservar resposta completa de controle aposentado indefinidamente.

**Alternativas consideradas:** aumentar a quota apenas adia o bloqueio; reservar N revisões falha após N negativas/expirações; cache ilimitado rompe memória; eviction sem high-water permite repetir efeito; gravar automaticamente viola consentimento. A simplificação de uma revisão corrente substitui a proposta inicial de 16 simultâneas, sem reduzir operações disponíveis sequencialmente. Rever limites por medição futura, preservando progresso e idempotência.

### D7 — I1: codec compartilhado antes da publicação

**Decisão:** T068/T070 e 002/T018 são pré-requisitos de T055/T056. T054 guarda bytes opacos e não depende do codec; testes/codec vêm depois dele e antes de validar/congelar/publicar US3. Manter IDs e associação US4 para rastreabilidade; explicitar execução antecipada no DAG.

**Justificativa:** US3 publica pacote `.jcomp` com hash dos bytes exatos e precisa do codec antes da revisão, ainda que import/export sejam entregues depois. Não criar segundo serializador provisório. **Alternativas:** adiar publicação até toda US4 aumenta acoplamento; antecipar somente testes/codec resolve a dependência sem ampliar o incremento. Conferir DAG ao regenerar tarefas e manter publicação bloqueada sem codec.

### D8 — I2: uma unidade de comprimento textual

**Decisão:** conteúdo até 65.536 pontos de código Unicode válidos, nomes até 256 e limites menores específicos preservados, em editor/MCP/biblioteca/proposta; rejeitar Unicode inválido, sem normalizar/truncar. Frames/pacotes continuam em bytes. maxLength do schema é coerente com a regra de caracteres da [especificação JSON Schema 2020-12](https://json-schema.org/draft/2020-12/json-schema-validation#name-maxlength); os casos de astral/combinação discriminam contagem de bytes, UTF-16 e grafemas.

**Justificativa:** 20.000 emojis ocupam 80.000 bytes UTF-8 e continuam dentro do campo; o orçamento do frame é conferido separadamente. Novas entradas/alterações usam a mesma validação de aplicação/domínio nos dois canais; leitura histórica não trunca valores. **Alternativas:** bytes exigiriam restrição semântica extra nos pacotes e fariam idiomas distintos consumir limites diferentes; UTF-16 é dependente de cliente. Não migrar formatos históricos incidentalmente.

### D9 — I3: defaults de proposta representam conhecimento e tipo

**Decisão:** PropertyDescriptor fechado exige default.status unknown/observed/proposed, valor tipado quando conhecido, motivo para desconhecido/sugerido e evidências para observado. Tipos string/integer/boolean/enum/reference e collection escalar homogênea até 256 itens; referências são IDs lógicos. Enum/evidências/capacidade instalada exigem validação semântica além do schema. Registros complexos ficam unknown com limitação explícita; não se afirma cobertura integral do legado.

**Justificativa:** promessa textual de defaults precisa de representação verificável; ausência de prova não é string vazia. Unknown registra lacuna, mas não libera publicação se faltar default obrigatório do contrato instalado. Valores alegados não alteram catálogo/template automaticamente nem comprovam homologação. **Alternativas:** qualquer JSON reabre objetos executáveis/sem tipo; remover defaults do contrato esconderia o requisito; inferir valor obrigatório inventa evidência. Ampliar coleções exige contrato próprio.

### D10 — U2: proposta anexada à sessão existente e reserva de recursos

**Decisão:** wrapper legacy exige sessionId/requestId/expectedRevision/proposal. Reusar AuthoringSession para até 32 anexos; expectedRevision é do documento, proposal.revision é independente (criação=1, atualização=current+1), com baseRevision e ownership. Alterar anexo não usa Session.Change/Project.Revision/undo. documentDirty/proposalDirty compõem dirty e impedem encerramento silencioso. Sessão blank em memória permite analisar sem arquivo. Reservar memória/metadados de finalização dentro dos 256 MiB e staging até 250 MiB/alvo, total 1 GiB/processo, antes de admitir trabalho persistível; IO em blocos limitados, sem clone integral adicional para revisar.

**Justificativa:** evita revisão paralela de projeto para texto de proposta, cache ilimitado e perda de rascunho no fechamento. O teto de 32 anexos limita metadados; staging fora do heap permite finalizar pacotes dentro do limite existente, com reserva global para não multiplicar 250 MiB por todo alvo potencial. Reserva é orçamento lógico/checagem de recursos, não garantia contra disco cheio ou falha de SO; falha externa não autoriza salvar sem consentimento. **Alternativas:** registro próprio no adaptador duplica ciclo de vida; inserir proposta no manifesto altera formato sem necessidade; armazenar bytes inteiros duplicados na revisão pode exceder a memória. Ajustar valores por evidência, sem retirar reserva de rascunho já admitido.

## Fechamento do segundo replanejamento — U3/U4/I4

### D11 — U3: análise vinculada a captura concreta

**Decisão:** proposal.materialization obrigatório, null ou referência de captura selection/revisão/templateSha256 da mesma AuthoringSession. Autor constrói manualmente ou agente pelos comandos I1, captura via I2 e atualiza a análise por CAS. NewKind exige null; composição/variante sem captura é analisável mas não publicável. Review vincula ambas as revisões, hash do template, analysisSha256 dos bytes retidos e pacote final, com proveniência/evidências sanitizadas. A captura é o único alvo persistível; análise continua transitória e exige descarte informado ao fechar.

**Justificativa:** descrições de propriedades/estados não definem template executável. A ponte explícita mantém interpretação no agente e construção/validação nos domínios existentes; evita duplicar codec ou inferir geometry/Kind de texto. Mudanças invalidam revisão sem recaptura automática, e um efeito não consome duas reservas. Fonte real/default legado percorre a mesma criação e apresentação dos dois canais.

**Alternativas consideradas:** geração automática a partir de statements exigiria interpretação fora do contrato; incluir outro template livre no schema duplica Tela/Projeto; publicar só classificação perderia conteúdo concreto. Novos tipos e recursos externos continuam ampliação separada. Revisitar quando houver requisito de materialização diferente, com contratos e evidências próprios.

### D12 — U4/I4: chamada de revisão alcançável e ordem de execução

**Decisão:** componente_revisar é controle R1, ReviewRequest fechado com quatro operations; retorna ReviewResponse existente, com geração nos reviewBytes. Publicar/importar recebem anexo/revisão; retirar/exportar recebem definição/versão/hash; só exportar aceita filename .jcomp na raiz workspace autorizada, criação exclusiva v1. Os quatro commits R2 aceitam apenas CommitRequest e conferem operation. Importação admite bytes/anexo em componente_preparar_importacao, R1 edição antes da revisão, consumindo quota/reserva. Schema library-operations define os pedidos sem mudar protocolo assinado/autoridade. Exemplo de paralelismo corrige I4, separando T068/T070 antecipados da portabilidade posterior a T057.

**Justificativa:** exigir reviewId sem chamada de preparação deixa o fluxo incompleto; jornada_revisar_alteracoes continua somente create/replace. Separar admissão de importação preserva D6; consulta R0 não vira mutação e controle não cria conteúdo isento. Exportação reutiliza finalidade/raiz existente, evitando deslocar a lacuna para rootIds/flags não definidos. Hashes e transição do índice entram na revisão/registro congelados, com gate e recuperação existentes.

**Alternativas consideradas:** reutilizar chamada de jornada sem união discriminada mistura alvos; autoaprovar publicação viola consentimento; receber pacote novo na revisão contorna quota; novas raízes ou overwrite de exportação exigem escopo adicional. IDs das tarefas são preservados e os handlers futuros selecionam schemas específicos. Essas decisões não qualificam host, Windows ou software; dependem da implementação/testes já planejados.

### D13 — handles Windows e namespace comum de coordenação

**Contexto/decisão:** resolução textual da raiz seguida de concatenação/reabertura permite redirecionamento entre verificação e efeito. O adaptador workspaceguard conserva volume/file ID e handle da pasta; leitura e criação exclusiva usam NtCreateFile relativo a RootDirectory com OBJ_DONT_REPARSE/FILE_OPEN_REPARSE_POINT. Rename de criação também é relativo ao mesmo handle e não substitui destino. Nova operação revalida a identidade fixada na seleção. O lease de leitura não adquire autoridade de escrita.

Gate operacional comum usa a pasta conhecida do Windows LocalAppData/JornadaStudio/workspace-locks, independente de JORNADA_CONFIG_DIR. Preferências e biblioteca mantêm sua configuração anterior; somente o namespace da exclusão não pode variar entre perfis. Arquivo aberto sem compartilhamento mantém a trava até fechamento/queda do processo. A raiz em si aceita compartilhamento normal para não bloquear Explorer/monitoramento.

**Justificativa e alternativas:** hash mais MoveFileEx não é CAS; comparação de string, mutex apenas em memória ou lock por perfil permitiriam redirecionamento/escritores não coordenados. Não exigir administrador, Git nem exclusividade artificial sobre o diretório. A API nativa é localizada no adaptador Windows e segue os bindings/testes de x/sys v0.47.0 já fixado. Outros SOs recusam o provider, sem fallback de segurança inventado. Custo: código específico Windows e testes reais de NT; revisar se o contrato ampliar sistemas de destino.

**Evidência/limites:** testes Windows -race comprovam leitura/base, criação exclusiva, homônimo, cancelamento, gate entre processos, alias, junction real e troca do caminho após abertura sem redirecionar escrita/inventário. O ambiente não permitiu criar symlink de arquivo; a alternativa de teste usou junction de diretório e não comprova aquela ramificação específica. Substituição existente continua indisponível até T043. Conectar todos os writers GUI/Git é trabalho de T040 e 037/T080–T081; o novo guard isolado não comprova essa coordenação integral.

Fontes primárias: [NtCreateFile](https://github.com/MicrosoftDocs/sdk-api/blob/docs/sdk-api-src/content/winternl/nf-winternl-ntcreatefile.md), [OBJECT_ATTRIBUTES](https://learn.microsoft.com/en-us/windows/win32/api/ntdef/ns-ntdef-_object_attributes), [FILE_RENAME_INFORMATION](https://github.com/MicrosoftDocs/windows-driver-docs-ddi/blob/staging/wdk-ddi-src/content/ntifs/ns-ntifs-_file_rename_information.md) e teste NtCreateFile/NtSetInformationFile no módulo local golang.org/x/sys/windows/syscall_windows_test.go.

### D14 — staging em handle privado, orçamento de finalização e JSON comum

**Decisão implementada:** projectzip.Write codifica em fluxo no handle privado e projectzip.Verify relê o mesmo ReaderAt, validando manifesto, recursos, CRC/hash/dimensões sem reconstruir imagens/ZIP em memória. Encode usa a mesma codificação e conserva o golden pré-extração. O provider consent mantém somente handle/ID/hash/tamanho; AuthoringReviewFreezer, definido na aplicação, fornece duas operações e recebe aliases somente leitura sob a trava da sessão. O resumo é derivado dos dados pelo aplicativo; revisão canônica vincula IDs separados de preparação/commit, revisão, geração, base, destino, hash final, desafio e validade fixa. Proposta/materialização são null para jornada; extensões I2/I3 continuam futuras.

O temporário Windows usa FILE_DELETE_ON_CLOSE e DELETE access; kernel o remove também na queda do processo. Raiz privada deve estar fora do workspace físico, inclusive na resolução de ancestrais ainda não existentes. GetFinalPathNameByHandle/VOLUME_NAME_NT fornece apenas nomes normalizados para comparação; não se passa esse nome NT a filepath ou IO Win32 ordinário. Outro escritor comum não obtém acesso ao arquivo congelado enquanto seu handle permanece aberto. Isso não pretende isolar um agente com privilégios equivalentes sobre todo o processo/SO, nem comprova substituição atômica de destino existente.

Reserva global de 64 MiB para um encode/verify pesado por vez, mais 1 MiB por sessão para revisão/controle, incluídos nos 256 MiB. Uma segunda preparação retorna busy sem avançar controle; consultas e edições continuam sujeitas às quatro vagas normais. Reservas de staging mantêm 250 MiB/alvo e 1 GiB/processo. Descarte/fechamento destacam artefatos; a limpeza executa fora da trava global e só devolve orçamento depois de Release confirmar. Falha conserva cobrança para nova tentativa. Intenção incerta não é descartada por inatividade/EOF; exige recovery. Varredura a cada 30 s não espera sessão ocupada.

**Justificativa/alternativas:** quatro decodificações pesadas simultâneas multiplicariam buffers de manifesto/JSON; reservar um slot global conserva o caminho de finalização depois da quota de edição. A reserva é contabilidade operacional, não promessa de heap Go exato ou disco disponível. ZIP inteiro em bytes.Buffer e clone de imagens na revisão foram evitados. Preparar revisão não é publicar nem consentir; referência de host/chave de harness não qualifica R2. A exposição da revisão por ferramenta/configuração continua em T035/T038.

Validação estrita antes do encoding/json fica no adaptador jsonwire, compartilhada por MCP e manifesto .jornada: duplicatas, UTF-8 inválido e surrogate isolado são recusados, sem conversão silenciosa para U+FFFD. ValidateStrings protege também o caminho Go→JSON; não percorre bytes binários de recursos. O domínio não importa esse adaptador. Reutilizar um decoder permissivo no ZIP violaria o mesmo requisito Unicode já exigido nos dois canais; a regressão de surrogate falhou antes da correção e passou após a extração.

**Evidências/limites:** testes Windows -race do guard comprovam negação de segundo escritor e limpeza pelo kernel em subprocesso encerrado sem defer. Provider passa schema normativo/ordem de campos/ZIP real/hash/base alterada/destino privado inválido e rollback após falha. Integração real Authoring→Freezer comprova replay byte a byte, repreparo, aposentadoria e descarte sem jornada no workspace. Testes do núcleo comprovam quotas cheias, expiração fixa, reserva retida durante limpeza e consulta não bloqueada por IO. Host humano, assinaturas, journal, publicação, writers GUI/Git e Office/DPI permanecem gates abertos.

Fontes primárias: [NtCreateFile e FILE_DELETE_ON_CLOSE](https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntcreatefile), [GetFinalPathNameByHandleW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfinalpathnamebyhandlew). Implementação isolada em workspaceguard/consent, sem dependência nova ou alteração de formato .jornada.

## Evidências técnicas ainda pendentes — não confundidas com decisões

- Pin/licença adicionados e G0 stdio comprovada no EXE Windows/amd64 padrão sem janela/preferências. A integração distribuída completa e host humano continuam pendentes, conforme quickstart.
- Alternativa A/HostApproval-v1 foi selecionada no plano; host real/versão ainda precisa ser integrado e qualificado. Expiração/revogação/consumo são contratados, mas não demonstrados.
- Diferenças/revisão canônica de jornada e replay estão testados no núcleo/provider; exposição MCP, consentimento, coordenação GUI/MCP/Git e recuperação persistente permanecem pendentes.
- Formato físico `.jcomp` e schemas de definição/proposta/consentimento foram definidos; round-trip, deduplicação/remapeamento, integridade e compatibilidade ainda precisam de implementação/teste.
- Limites, paginação, prazos e erros de wire estão em mcp.md; schemas técnicos existem e ainda devem ser conferidos contra handlers/codec na implementação.
- Geometria detalhada e integração dos controles visuais ao espaço existente; fixtures, DPI/Office, comandos/harness e tarefas por domínio.

Histórico: tarefas foram inicialmente geradas sem execução e passaram por análise antes da autorização posterior para implementar. Atualmente há dependências, código, build G0 e testes técnicos, com escopo e revisões no quickstart; não há homologação de host humano, paridade visual completa ou aceite funcional.
