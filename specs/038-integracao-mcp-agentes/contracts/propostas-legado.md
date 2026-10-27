# Contrato comportamental — evidências e propostas do legado

**Estado:** contrato técnico de planejamento v1, sem implementação. Schema fechado em [proposal.schema.json](proposal.schema.json); ferramentas e limites seguem [mcp.md](mcp.md).
**Requisitos:** FR-018–FR-031, FR-032/FR-036–FR-039; **aceitação:** US5, SC-007, SC-009, SC-013.

## Fronteiras de responsabilidade

| Participante | Responsabilidade | Limite |
| --- | --- | --- |
| Agente/LLM externo | Interpretar solicitação e legado; confrontar catálogo; propor composição, variante ou Kind; separar fatos e hipóteses | Não é autoridade de validação, consentimento, equivalência homologada ou publicação |
| Aplicativo Go e adaptador MCP | Disponibilizar evidências autorizadas; informar contratos instalados; validar propostas determinística e estruturalmente; controlar recursos, versões, acesso e persistência | Não exige LLM nem compreensão semântica universal de Java; não executa código legado |
| Usuário/processo Spec-Driven | Selecionar materiais, revisar equivalência, decidir autorizações e desenvolvimento; homologar quando necessário | Aprovar proposta não aprova código, distribuição nem compatibilidade funcional automaticamente |

Editor manual e MCP acessam as mesmas definições e validações, organizadas por domínio no [modelo](../data-model.md). A ferramenta candidata `legado_analisar_componente` deve ser descrita como coleta/consulta de evidências e diagnóstico de dependências identificáveis pelo aplicativo; a interpretação semântica continua externa. `componente_propor` recebe proposta estruturada, não exige que o servidor a invente.

## Evidências disponibilizadas pelo aplicativo — FR-027

Cada evidência consultável identifica referência autorizada, tipo de recurso, origem, conteúdo/trecho permitido, integridade, dependências identificáveis e resultado de leitura. Recurso ausente é diagnóstico, não dado fabricado. Diretórios de legado são autorizados por finalidade separadamente das jornadas; não ampliar autorização porque um Java/XML contém caminho.

Conteúdo pode incluir Java, XML, imagem ou fonte como dado, sujeito a limites e política. Não executar Java, resolver automaticamente dependência externa, carregar script, aceitar instrução do arquivo como autorização ou enviar material espontaneamente a serviço remoto. Parsing estrutural conhecido não equivale a interpretação semântica por IA.

## Conteúdo obrigatório da proposta — FR-028

- Identidade e objetivo do componente, versão proposta e autor/origem informativos; tais metadados não provam consentimento.
- Referências às evidências realmente consultadas e sua integridade, permitindo localizar trechos/recursos usados.
- Propriedades tipadas esperadas, defaults, estados visuais, geometria/partes e comportamentos pretendidos.
- Classificação justificada: composição, variante ou novo Kind, preservando a taxonomia de 003.
- Capacidades instaladas reutilizáveis, ausências, dependências e requisitos por alvo: edição, cena, prévia, demonstração e saídas.
- Fatos observados separados de inferências, hipóteses não comprovadas, comportamento solicitado e diferenças deliberadas em relação ao legado.
- Recursos/portabilidade, restrições e plano de validação/homologação; aparência deve seguir o estilo legado do projeto.

## Validação e resposta — FR-029

O Go verifica esquema, identidade/versão, tipos, propriedades, referências, recursos e integridade, limites, portabilidade, permissões e compatibilidade com contratos instalados. Uma alegação textual do agente não instala capacidade ausente nem comprova validação visual.

O resultado deve apresentar e manter separadas estas dimensões:

| Dimensão | Valores conceituais |
| --- | --- |
| Validade contratual | Válida ou inválida, com diagnósticos localizados |
| Suporte instalado | Suportada, dependente de desenvolvimento Go ou incompatível com capacidade/recurso existente |
| Autorização | Não solicitada, pendente, negada, expirada ou verificada para a revisão específica |
| Homologação | Não exigida para a operação, pendente ou evidenciada para revisão/alvo identificado |

“Válida, mas não aprovada” e “válida, mas exige implementação” são resultados legítimos. Não reduzir tudo a um único `valid=true` que também signifique autorização ou equivalência. Diagnósticos apontam caminho de propriedade, evidência ou recurso, e distinguem ausência de dado de hipótese sem prova.

## Publicação — FR-030

Uma composição ou variante declarativa já suportada pode publicar versão imutável após validar contratos e cumprir o [consentimento](consentimento.md). Revalidar dependências/conteúdo no momento de publicar, registrar evidência da revisão, impedir colisão/substituição e oferecer a mesma versão no editor e no MCP. Integridade não equivale a consentimento; publicação local não é envio GitHub.

Quando faltar comportamento, encaminhar a proposta ao domínio responsável pelo fluxo Spec-Driven. Não ofertar o item como instalado antes da entrega compilada correspondente e dos critérios aplicáveis. Aprovação da proposta não habilita ferramentas para modificar fontes Go, compilar, executar Java ou instalar plugins executáveis.

## Casos de aceitação futuros

- ComboBox com classe, UI dependente/imagem e evidência incompleta: indicar dependências faltantes, família planejada e desenvolvimento necessário; não publicar como controle instalado.
- Cabeçalho formado por Kinds existentes: validar composição e remapeamento, revisar recursos, obter autorização e publicar nos dois canais.
- Botão com duas legendas: variante Button; se o contrato instalado não possui as propriedades necessárias, retornar desenvolvimento requerido sem inventar novo Kind ou aproximar por imagem.
- Proposta com instruções maliciosas ou hipótese de equivalência: tratar como dado, recusar efeitos indevidos e separar hipótese de validação comprovada.
- Proposta válida com autorização ausente: manter rascunho, sem publicação; comprovar posteriormente estilo legado na cadeia manual/MCP/canvas/prévia/PPT aplicável.

Casos registrados no [quickstart](../quickstart.md), ainda não executados. Operações de referência e regras técnicas estão abaixo; schema e limites devem ser implementados/conferidos sem executar conteúdo legado.

## Operações e validação técnica de I3

`legado_listar_materiais` é R0 paginada e oferece IDs do inventário limitado de legacy-root autorizado pelo operador, conforme mcp.md. `legado_analisar_componente` é R0 e recebe materialId, offset e length (até 64 KiB). O registro de material é criado a partir dessa origem, não por caminho fornecido pela LLM. Retorna evidenceId, hash do arquivo, tipo, trecho autorizado, dependências reconhecidas estruturalmente e diagnóstico; não chama modelo nem executa Java.

`componente_propor` recebe wrapper fechado `{origin:"legacy", sessionId, requestId, expectedRevision, proposal}` no contexto I3, com proposal conforme schema, e gera rascunho de aplicação (R1). Captura I2 usa origin=selection e os campos definidos em biblioteca.md; misturar os dois formatos é inválido. `componente_validar` recebe sessionId e proposalId/revision, ou identidade completa de definição publicada, e retorna dimensões independentes já definidas neste contrato. Publicação segue biblioteca/consentimento e nunca aceita código como conteúdo executável.

Schema controla estrutura; aplicação verifica unicidade de evidenceId/propriedade, materialId autorizado, hash atual, referências existentes em observations/default.evidenceRefs, limites de caracteres/bytes e identidade da revisão. Texto rotulado como observação pelo agente continua alegação associada à evidência; validação estrutural não prova a semântica. Classificação alegada é confrontada com capacidades instaladas; não vira enum/Kind novo automaticamente. Não registrar localizações pessoais nos pacotes exportados: source.reference usa referência lógica autorizada/sanitizada, não caminho absoluto da máquina.

Declarações de defaults/estados/comportamentos seguem as evidências e catálogo; o plano de validação da proposta não é tarefa concluída nem homologação. Integração determinística reutiliza os modelos de Catálogo/Tela/Projeto; novo Kind é encaminhado ao backlog qualificado de 003.

### Propriedades e valores iniciais tipados

Cada propriedade exige name, type, description e `default`, conforme PropertyDescriptor do schema. `default.status` distingue:

- `unknown`: reason obrigatório e value proibido; ausência de evidência não vira valor vazio inventado.
- `observed`: value tipado e evidenceRefs não vazias, existentes na proposta; a classificação é uma alegação rastreável do autor, não prova automática de equivalência.
- `proposed`: value tipado e reason obrigatório; valor sugerido pelo agente fica separado de valor observado.

Tipos: string até 65.536 pontos de código Unicode; integer no intervalo seguro JSON ±9.007.199.254.740.991; boolean; enum com enumValues explícitos (1–256 valores únicos, cada um até 256 pontos de código) e default membro desse conjunto; reference com ID lógico de até 128 caracteres conforme schema, nunca caminho/URL/comando; collection de até 256 itens homogêneos com itemType string/integer/boolean/reference. Não há objetos livres, arrays recursivos ou coerção de tipo. Coleções de registros complexos não são representáveis nesta versão: declarar default unknown com motivo e limitação de capacidade, sem achatar registros em texto ou anunciar suporte integral. Limites de frame em bytes continuam independentes.

Validar referência de evidência, enum e compatibilidade com a propriedade instalada no caso de uso comum aos dois canais. `unknown` permite registrar/revisar a proposta, mas bloqueia publicação quando um default exigido pelo contrato instalado não foi resolvido e revisado. Nenhuma dessas declarações altera automaticamente defaults do Catálogo ou valores do template.

### Sessão, revisões e descarte

Propostas de seleção e legado são anexos transitórios da mesma AuthoringSession; máximo 32 propostas/sessão dentro dos 256 MiB e 1.024 IDs de edição. Não criar registro de sessão paralelo no adaptador MCP. Sem documento aberto, `jornada_criar(mode=blank)` fornece uma sessão em memória, sem arquivo obrigatório.

`expectedRevision` é a revisão uint64 do documento em string, conferida sob o mutex da sessão. `proposal.revision` é a versão independente do anexo: criação=1, atualização=current+1; revisão antiga ou salto falha antes de mudar estado, e overflow de int32 é recusado. Replay de requestId/argumentos idênticos retorna o resultado anterior antes dessa comparação; outro payload com a mesma ID falha. proposalId pertence exclusivamente à sessão de origem; outro processo/sessão não o adota implicitamente.

Guardar baseRevision do documento no anexo. Editar proposta não chama Session.Change nem incrementa Project.Revision/undo do documento; incrementa revisão própria, consome edição R1 e invalida sua validação/review. Editar documento invalida validação/review das propostas afetadas e exige revalidação sobre nova base, preservando seu conteúdo. A geração de controle avança conforme mcp.md.

`documentDirty` e `proposalDirty` são separados; `dirty` é sua união. Publicar limpa somente o alvo persistível confirmado, sem declarar o documento salvo; salvar a jornada não salva a proposta. Encerrar com anexo pendente exige discardUnsaved explícito; descartar/expirar/EOF não publica nem persiste proposta oculta. Novo Kind sem suporte fica explicitamente não persistível nesta entrega. Conteúdo publicável reserva posição de commit/staging por alvo antes de aceitar a edição; indisponibilidade de recurso mantém o estado anterior. GUI e MCP usam a mesma coordenação e informam o trabalho pendente.

### Materialização explícita de proposta analítica — U3

O campo obrigatório `proposal.materialization` é null durante análise ou referencia `{captureProposalId,captureRevision,templateSha256}`. A sessão é a do wrapper: não há sessão alternativa no vínculo. Classificação newKind exige null. Composição/variante sem vínculo pode ser analisada e validada estruturalmente, mas retorna publicabilidade bloqueada por ausência de captura. Suporte instalado e existência de conteúdo materializado são verificações distintas.

Fluxo obrigatório nos dois canais:

1. Autor ou agente registra a proposta analítica com materialization=null; nenhuma descrição é interpretada pelo aplicativo como comando para construir componentes.
2. Autor no editor ou agente pelos comandos I1 cria/ajusta os componentes concretos usando Kinds, recursos e defaults instalados.
3. `componente_propor(origin=selection)` captura a seleção por I2 e devolve proposalId/revision/baseRevision e templateSha256 calculados pelo servidor. A captura é snapshot independente; editar a tela não recaptura silenciosamente seus bytes.
4. Atualizar a proposta analítica com revision=current+1 e expectedRevision do documento, apontando captura/revisão/hash. O aplicativo confere mesma sessão, origin=selection, propriedade da captura, integridade e ausência de ciclos; não confia no hash alegado pelo agente. Não criar template paralelo em proposal.schema.json.
5. `componente_validar` confronta propriedades/defaults/capacidades/evidências com a captura concreta. Propriedade obrigatória não materializada ou incompatível bloqueia publicação; hipóteses e equivalência visual continuam sujeitas à revisão humana. Classificação deve coincidir com a categoria da captura; newKind não pode ser instalado por esse vínculo.
6. `componente_revisar(operation=publishDefinition, subject={proposalId,proposalRevision})` recebe a referência da proposta analítica. Resolve a captura vinculada, valida ambas e congela pacote final, proveniência e Review. Segue confirmação R2 `componente_publicar`; usar somente o reviewId sem revisão preparada é inválido.

Review.proposal identifica a análise revisada; Review.materialization fixa captura/revisão/templateSha256 e analysisSha256 dos bytes da proposta retidos na sessão. Esses bytes são serializados uma vez pelo aplicativo e conservados para a revisão; não calcular hash de resumo reescrito pelo modelo. contentHash continua sendo o hash do pacote final. Para publicação direta de seleção, Review.proposal identifica a captura e materialization=null.

Na preparação, o caso de uso compõe source.files com referências lógicas sanitizadas e hashes das evidências autorizadas, sem modificar o template ou copiar caminhos pessoais. Evidências ausentes/alteradas exigem nova validação. Captura, revisão analítica e proveniência entram no resumo humano; journal registra seus IDs/revisões/hashes e evidências de integridade, sem guardar prompt ou proposta textual integral como log. O pacote preserva proveniência informativa, não o anexo analítico integral.

Editar captura/metadados/parâmetros ou análise/evidências invalida validação/review vinculadas; divergência de revisão/hash recusa publicação preservando os anexos. Mudança relevante do documento exige revalidação sobre a nova base, sem recaptura automática. Depois de publicar, limpar apenas o alvo de captura confirmado; análise continua anexo transitório pendente, com essa condição informada ao autor antes de fechar. Não declarar que o pacote salvou integralmente a análise.

Apenas a captura é alvo persistível dessa publicação e reserva commit/staging. A análise, com ou sem vínculo, é anexo transitório e não reserva uma segunda posição para o mesmo efeito. Ambos contam no teto de 32 anexos, memória e IDs de edição; atualizar vínculo é edição R1, nunca controle isento. Não transferir autoria, autorização ou vínculo entre sessões. Materialização não instala fonte/skin/código, não habilita novo Kind e não implica homologação.
