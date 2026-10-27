# Modelo por domínio e coordenação — SPEC 038 v0.4

**Estado:** desenho técnico com implementação parcial. Authoring/AuthoringDraft/AuthoringWorkspace e mudanças tipadas de tela/componente/navegação já existem na aplicação, com snapshots, CAS, replay, quotas e fonte de documentos; Review assinado, biblioteca e proposta específica do legado ainda exigem suas tarefas. Nenhuma migração ou extensão do formato .jornada foi feita neste incremento.
**Requisitos:** FR-001–FR-004, FR-007, FR-015, FR-018–FR-033; **aceitação:** US1/US3/US5 e SC-013.

O projeto conserva um monólito modular com cinco domínios. A biblioteca reutilizável estende responsabilidades do Catálogo, e a autoria coordena instâncias e recursos dos demais domínios. MCP, editor Fyne, filesystem e exportadores são adaptações; IA externa não é entidade de negócio. Essa divisão deve existir na implementação e nas evidências, não apenas neste documento.

## Projeto

**Fonte de verdade:** [002/data-model](../002-projeto/data-model.md) e contratos de formato existentes.

Responsável por identidade/revisão do projeto, catálogo/formato compatível e metadados de recursos incorporados. Imagens instanciadas devem acompanhar a jornada com integridade verificável. Documento não recebe seleção da UI, sessão MCP, autorização, credencial, prompt ou estado de LLM.

Não se altera schema/catalogVersion pela simples adição de MCP. Copiar composição não cria dependência oculta de caminho externo ou da biblioteca de origem. Evolução de formato, fontes externas e metadados persistidos de proveniência exigem contratos próprios antes de implementação.

## Catálogo

**Fonte de verdade:** [003/data-model](../003-catalogo/data-model.md) e [003/spec](../003-catalogo/spec.md), especialmente FR-007–FR-012.

Responsável por Kind, variantes, defaults, propriedades tipadas, capacidades e disponibilidade. Definição reutilizável descreve conteúdo semântico, parâmetros e composição suportada; não armazena código executável ou tipos MCP/Fyne.

| Conceito estendido | Dados e invariantes |
| --- | --- |
| Definição reutilizável | Identidade estável, categoria, descrição, origem informativa e versões; nome não é identidade |
| Versão | Identidade da versão, conteúdo verificável, capacidades/requisitos e recursos; conteúdo publicado é imutável |
| Oferta | Estado de publicação/retirada e compatibilidade calculada com motivos; retirar da oferta não altera conteúdo publicado |
| Parâmetro | Propriedade de Kind suportado ou referência externa explicitamente requerida; tipo, default e restrições devem ser validados |
| Composição | Componentes existentes, geometria relativa e vínculos internos; sem ciclos nem comportamento novo implícito |

Variante preserva família semântica; novo Kind precisa de suporte real no catálogo e consumidores compilados. Skins/fontes do legado já instaladas são referências tipadas; recursos visuais ausentes bloqueiam a oferta afetada. Declaração de capacidade não substitui evidência.

O armazenamento da biblioteca é adaptador; a validação de definição não depende de acesso direto ao filesystem. Caminho, codec `.jcomp` v1 e schema estão desenhados nos contratos; sua implementação e validação experimental continuam pendentes.

## Tela

**Fonte de verdade:** [004/data-model](../004-tela/data-model.md).

Responsável por tela, dimensões, fundo, componente, UUID, FileID, bounds, estilo confirmado e ordem de camadas. Instanciar uma composição produz componentes independentes editáveis; gera novas identidades e remapeia partes internas, opções, teclado e referências correspondentes.

A identidade de uma definição/versionamento não é UUID da instância. Não reutilizar identidades da biblioteca dentro da jornada. Referência externa à composição exige parâmetro resolvido ou diagnóstico antes de publicar a operação. Nenhum vínculo vivo ou atualização silenciosa da biblioteca é introduzido na primeira entrega.

## Jornada

**Fonte de verdade:** [005/data-model](../005-jornada/data-model.md).

Responsável por início, ações e destinos explícitos. Remapeamento de cópia e instanciação deve preservar a intenção dos vínculos, sem inferir navegação de nome, texto, código de negócio ou atalho. Destinos para outra jornada continuam restritos às regras de arquivo da mesma pasta; disponibilidade de arquivo é verificação do caso de uso/adaptador, não importação de filesystem pelo domínio.

## Demonstração

**Fonte de verdade:** [006/data-model](../006-demonstracao/data-model.md).

Responsável por roteiro finito, passos, alvos, ordem e tempo. Novo Kind não ganha ações de demonstração automaticamente. Estado de prévia é transitório e não altera valores iniciais salvos, roteiro ou histórico do editor. Capacidade de prévia, demonstração e PPT é declarada separadamente.

## Aparência do legado compartilhada

Catálogo define família/defaults e capacidades; Tela conserva escolhas/overrides confirmados; o caso de uso resolve os recursos necessários; apresentação deriva a cena comum para canvas, prévia e PPT. Adaptadores pintam ou exportam a mesma representação e preservam as limitações de cada alvo.

Nos dois modos de autoria, ausência de skin explícita resolve o mesmo default do legado. Fontes reais, cores/alpha, métricas, imagens, estados e decoração das partes de compostos percorrem a mesma cadeia. Não aplicar decoração somente ao inserir pela GUI, nem criar aparência genérica para itens MCP. Foco, seleção, zoom e tema operacional não entram no documento ou nas imagens exportadas.

## Coordenação de aplicação e adaptadores

| Modelo de coordenação | Papel e limite |
| --- | --- |
| Documento | Compõe os cinco domínios e recursos; valida referências entre eles sem se tornar objeto central de UI, IO e exportação |
| Sessão de autoria | Cópia isolada, revisão, base, histórico e ciclo de vida; estado externo ao manifesto |
| Proposta do legado | Evidências, hipóteses e classificação; a validação contratual é determinística e a interpretação é do agente externo |
| Revisão de operação | Conteúdo/base/destino/efeito imutáveis para consentimento; não é autorização por si só |
| Decisão de autorização | Origem confiável, escopo, validade e consumo, verificados segundo mecanismo futuro; não é entidade persistida na jornada |

Gui e MCP traduzem entradas para os mesmos casos de uso. Não mover regras de negócio para DTO de ferramenta, widget ou biblioteca de IA. Filesystem, armazenamento da biblioteca, consentimento, transporte e codecs entram por fronteiras explícitas; interfaces pequenas apenas quando justificadas, conforme a constituição. O plano definirá caminhos e assinaturas concretos, evitando abstrações ou subdomínios artificiais.

## Evolução e verificação pendentes

Parâmetros/esquemas, pacote v1, versões, recursos, locking e autorização foram desenhados no [plano](plan.md) e contratos: [MCP](contracts/mcp.md), [biblioteca](contracts/biblioteca.md), [propostas](contracts/propostas-legado.md) e [consentimento](contracts/consentimento.md). São decisões técnicas propostas, ainda sem implementação/validação experimental. Structs concretas devem seguir essas fronteiras e não importar os DTOs de wire nos domínios.

O plano/tarefas devem manter requisitos e implementação nos respectivos domínios e rastrear integração em 038; não duplicar tarefas de novos Kinds do Catálogo. [Q038-18](quickstart.md) exige inspeção das dependências e comparação dos dois caminhos na implementação futura. Nenhuma validação arquitetural do código novo foi executada, pois não houve código novo nesta etapa.

## Estruturas e invariantes de implementação propostas

Os nomes abaixo são conceitos/estruturas planejadas, não tipos já adicionados. O schema externo é responsabilidade do adaptador; modelos internos usam structs tipadas e identidades já adotadas.

| Dono / estrutura | Campos/relações essenciais | Invariantes |
| --- | --- | --- |
| Catálogo / ReusableDefinition | definitionId, nome/descrição/categoria e origem informativa | Identidade distinta do nome e de UUID de componente |
| Catálogo / DefinitionVersion | definitionId, version, templateDigest, requirements, parâmetros/bindings | Conteúdo imutável depois de publicar; Kind existente e propriedade permitida; fonte/skin instalada identificada por integridade |
| Catálogo / Availability | published/withdrawn, compatibility, reasons | Estado de oferta separado de conteúdo; publicado pode estar indisponível; retirada não apaga instâncias |
| Tela / TemplateSelection | uma tela e componentes selecionados, origem relativa e vínculos internos | Limite 256; sem aninhamento de definições; referências externas recusadas inicialmente |
| Tela / InstantiationResult | novos componentes e mapa UUID/FileID/IDs internos | Nenhuma colisão; ordem/geometria preservadas; cada vínculo aponta para a cópia correspondente |
| Projeto / ResourceTransfer | assets metadados e bytes já validados no template | Imagens incorporadas, integridade conferida, deduplicação por hash sem colisão de ID/caminho |
| Jornada / Bindings | ações da seleção, tela de destino e componentes remapeados | Apenas referências explícitas válidas; sem resolver arquivo externo oculto |
| Demonstração / Scenario | roteiro existente e alvos suportados | Sem ações novas implícitas por instalar definição ou Kind |

O pacote `.jcomp` não exige campo novo em Component/Project para ligar instância à biblioteca. Sua versão usa SemVer x.y.z; schema do pacote=1; template usa projeto 3/4. Essas versões são independentes de constituição 1.4.0, aplicativo e protocolo MCP.

## Estados externos ao domínio persistido

| Coordenação de aplicação | Campos/estados | Política |
| --- | --- | --- |
| AuthoringSession | sessionId, Document, revisão/base, documentDirty, proposalDirty, dirty agregado, proposals, controlGeneration, reservas, lastActivity, histórico/cache | Mutex por sessão, snapshot independente, expiração 30 min; edição não grava automaticamente; até 32 anexos, sem registro paralelo de sessão |
| FileBaseline | exists, identidade volume/arquivo, size, SHA-256 | Capturada ao abrir; revalidada sob gate de raiz; ausência é precondição de criação exclusiva |
| RequestRecord | requestId, digest dos argumentos, revisão/resultado; controle inclui geração recebida | Até 1.024 IDs de edição sem eviction; controle separado, uma revisão corrente e até 16 registros recentes com high-water; replay aposentado é terminal sem efeito |
| OperationReview | Review fechado, bytes/hash, finalBytes em staging, controlGeneration, proposal opcional, expiração/desafio | Congelado; efeito/destino/base não mudam após revisão; revisão invalida em edição relevante; replay não renova prazo |
| AuthorizationDecision | host/key/actor, assinatura, issued/expires, decision | Não é dado de projeto; verificador consulta chave pinada, origem e escopo, nunca chave do chamador |
| CommitIntent | requestId, reviewHash, contentHash, base, target, posição reservada, fase e resultado | Journal antes do efeito e confirmação após; revisão já consumida não produz novo efeito; recovery reutiliza posição original |
| LegacyProposal / anexo de seleção | proposalId, revision, baseRevision, schemaVersion, propriedades/defaults tipados, evidências, observações/hipóteses, dependências e classificação | Sessão proprietária única, criação=1/atualização=current+1; sem Session.Change/Project.Revision/undo por editar anexo; proposta não instala tipo; dimensões de validade/suporte/autorização/homologação separadas |
| Vínculo de materialização | materialization null ou captureProposalId/captureRevision/templateSha256 na proposta analítica | Captura origin=selection da mesma sessão; hash/revisão conferidos no servidor, sem ciclo ou template paralelo; newKind exige null |
| Pedido de revisão da biblioteca | sessionId/requestId/expectedRevision/controlGeneration/operation/subject/target condicionado | Controle R1 separado de edição; quatro variantes fechadas conforme library-operations.schema.json; resposta ReviewResponse e confirmação R2 por operação correspondente |

A análise do legado é anexo transitório, mesmo vinculada. O único alvo persistível de sua publicação é a captura I2; não reservar dois commits/stagings para análise e captura do mesmo efeito. Ambos consomem memória/anexos e IDs ao editar. Publicar a captura não salva a análise integral: ela conserva proposalDirty até descarte explícito. O pacote contém proveniência sanitizada, e o journal metadados de IDs/revisões/hashes, sem persistir texto integral da análise como log.

Review.proposal referencia a análise em publicação vinculada; Review.materialization fixa captura/revisão/templateSha256/analysisSha256 dos bytes retidos, enquanto contentHash fixa pacote final. Publicação direta de seleção referencia a captura e usa materialization=null. Alterar qualquer uma das revisões/evidências invalida review; mudanças no documento exigem revalidação, sem recaptura silenciosa. Ver [materialização](contracts/propostas-legado.md#materialização-explícita-de-proposta-analítica--u3).

Anexo de importação é admitido por componente_preparar_importacao como R1 edição com materialId/hash autorizado e bytes originais congelados. A revisão de controle subsequente não cria anexo oculto nem contorna quota. Withdrawal congela a transição do índice; exportação v1 cria filename regular .jcomp na raiz já autorizada, sem substituir homônimo. Casos de uso comuns atendem GUI e MCP; não há nova entidade nos cinco domínios nem alteração de manifesto.

Propriedades de proposta possuem default unknown (reason, sem value), observed (value/evidenceRefs) ou proposed (value/reason), conforme schema. A estrutura preserva evidência e sugestão separadas; defaults exigidos não resolvidos impedem publicação. Coleção v1 contém até 256 valores escalares homogêneos tipados; registros complexos ficam unknown com limitação explícita, sem serialização em texto. Nenhum default da proposta sobrescreve o catálogo ou template automaticamente.

Invariantes de orçamento: edição e controle separados; `commits consumidos + intents incertos + reservas de alvos persistíveis sujos <= 64`, reserva por alvo antes de admitir conteúdo pendente. Reserva global de finalização de 64 MiB para um encode/verify por vez, mais 1 MiB/sessão para revisão/metadados, dentro dos 256 MiB; staging de até 250 MiB/alvo, total 1 GiB/processo, antes da edição/importação. Conteúdo admitido mantém caminho de revisão/commit/descarte mesmo após esgotar a quota de edição ou ocupar oito sessões. Negação/expiração não consome commit; resultado incerto conserva reserva; 64 commits selam novas edições. Limpeza física fora da trava global conserva a cobrança até fechar o artefato. Esses limites não prometem heap exato ou gravação diante de falha externa ou sem autorização.

Textos usam pontos de código Unicode válidos em ambos os adaptadores, pacote e proposta; conteúdo até 65.536, nomes até 256 e limites específicos menores preservados. Frames/pacotes são contados em bytes UTF-8; sem normalização/truncamento silencioso ou contagem por UTF-16. Regras comuns validam novas entradas/alterações; leitura de documento histórico não trunca conteúdo.

Fluxos de sessão: active → expired/closed; descartar invalida revisões e restaura base sem escrita. Fluxo de operação: draft → reviewed → awaitingConsent → authorized → committing → committed; denied/expired/stale/canceled/recoveryRequired são saídas explícitas. Depois de reinício muda serverInstanceId; journal recupera resultado, mas assinatura anterior não autoriza novo efeito.

Anexo aceita expectedRevision do documento e sua própria revision sob o mesmo mutex. Alterar documento invalida validação/review afetadas e exige revalidação da base do anexo; alterar anexo incrementa somente sua revisão e a geração de controle. Salvar documento não salva anexos; publicar limpa apenas o alvo persistível confirmado, conservando a análise transitória vinculada. dirty agregado protege ambos no fechamento; proposta não publicável de novo Kind continua transitória e exige descarte explícito, sem falsa promessa de persistência no manifesto.

O adaptador workspaceguard coordena raiz e escritores da aplicação; `Session` não resolve filesystem ou Git. O adaptador consent verifica envelopes; Catálogo/Tela não conhecem chave pública ou assinatura. A oferta manual e MCP lê a mesma Availability e instancia os mesmos tipos; GUI não inventa defaults ou libera dependência ausente para pintar algo parecido.

## Recorte implementado de autoria

AuthoringSessionInfo conserva também o nome original no cache de admissão, sem reler origem no replay; seu custo permanece contado no orçamento. InspectMetadata copia coleções/metadados sem imagens e Diagnostics valida presença real dos recursos sob a proteção de sessão. EditContext recusa candidato cancelado antes de publicar. ConfigureScreen/ValidateComponentEdit/BindNavigation/SetJourneyStart são compartilhados com os formulários. Texto histórico inalterado não é truncado; novas alterações usam o limite comum.

InstanceIDs/RegenerateInstanceIDs incluem opções e linhas, além dos compostos; DuplicateScreenWithMapping retorna o remapeamento e preserva referências internas. Resposta grande é recusada pelo coordenador antes da mutação se ultrapassar seu orçamento, e nunca apenas depois de publicar. Caches não carregam recursos binários em resultados de comando. AuthoringFrozenReview conserva só referência de artefato, hashes, revisão e bytes JSON pequenos; a aplicação recebe o freezer por interface e não importa filesystem/ZIP. O adaptador conserva o único handle privado de ZIP; não retém documento de domínio. SessionStatus dispensa cópia do documento; fechamento conserva apenas resultado terminal e tombstone, sem artefatos. Extensões de proposta/materialização e commit persistente continuam nas tarefas correspondentes.
