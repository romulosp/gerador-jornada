# Contrato biblioteca v1 — pacote declarativo e instâncias independentes

**Estado:** desenho técnico para revisão do plano; não implementado. **Requisitos:** FR-018–FR-026/FR-030–FR-033. **Schema:** [library.schema.json](library.schema.json).

## Formato e responsabilidades

Pacote `.jcomp` é ZIP com exatamente `definition.json` e `template.jornada` na raiz. Não extrair automaticamente arquivos para caminhos do pacote. Recusar diretórios, links, nomes duplicados/case-insensitive, caminhos externos, entradas extras ou expansão acima dos limites de [mcp](mcp.md).

`definition.json` é metadado de Catálogo, com schemaVersion=1, definitionId UUID, version SemVer x.y.z, name, description, category=`composition` ou `variant`, source, templateSha256, requirements e parameters. Origem é informação, não confiança. `template.jornada` usa o codec estrito existente, manifesto 3/4 e uma única tela válida; contém apenas a seleção capturada, sua ordem, referências internas e recursos necessários. Não introduzir outro formato para as propriedades de dez Kinds nem aceitar campos desconhecidos no template.

O template recebe identidades próprias, geometria relativa normalizada e início na única tela; referências à tela capturada são remapeadas para ela. Referências a componentes/telas/arquivos fora da seleção são recusadas com diagnóstico na v1, antes de publicar captura. Teclado→campo interno é preservado. Parâmetros de destino externo e aninhamento de definições são futuras ampliações, não placeholders executáveis.

Instanciar expande em componentes `screen.Component` independentes e editáveis na tela destino, com novos UUIDs/FileIDs/IDs internos, referências remapeadas e recursos copiados/deduplicados por hash. O componente instanciado não depende do arquivo `.jcomp` ou biblioteca original; o formato `.jornada` não recebe estado de autorização ou agente.

## Conteúdo e integridade

`requirements` identifica catálogo suportado, Kinds e recursos instalados exigidos; skin/fontes incluem tipo/nome/path lógico e SHA-256 dos bytes necessários. Conferir contra biblioteca instalada antes da oferta/inserção; nome parecido não basta. Recursos de imagem já incluídos no template conservam metadados de hash/tamanho e são validados pelo codec. Não instalar fontes, skins novas ou código durante importação.

`templateSha256` refere-se aos bytes exatos do template. `packageSha256` é calculado sobre o arquivo externo completo, fora de seu manifesto, usado no endereço imutável do store e revisão de importação/exportação. Cada publicação congela os bytes antes da revisão; não regenerar ZIP depois do consentimento. ZIP diferente com mesmo definitionId/version é colisão, mesmo quando alguém alegar equivalência; exigir nova versão ou reapresentar pacote idêntico.

`parameters` contém nome, tipo, valor default opcional, required e bindings `{componentId,property}`. Somente patches de propriedade/estilo existentes, tipados e validados no catálogo são permitidos; property é um identificador allowlisted de contrato, nunca expressão, script ou caminho JSON arbitrário. Na v1: text, initialValue, buttonVariant, fieldFormat e enabled. O tipo e a família alvo devem coincidir; um parâmetro pode vincular vários componentes compatíveis. Destinos/IDs e referências não são parâmetros livres na v1.

## Armazenamento e estados

Raiz é `components/` sob a configuração privada do operador, com:

- `objects/<packageSha256>.jcomp`: bytes imutáveis; criação exclusiva e integridade conferida.
- `index.json`: índice versionado de definitionId/version/packageSha256 e estado de oferta `published` ou `withdrawn`; não contém conteúdo mutável de versão.
- `operations/`: journal/intenção e recuperação; mesma disciplina de autorização, gate e consumo da persistência.

Rascunho/validação/revisão pertencem à sessão de aplicação; não aparecem como itens publicados. Fluxo: draft → validated → reviewed → published. Mudança de conteúdo invalida validação/revisão. Incompatibilidade e recurso ausente são disponibilidade calculada, separada do estado publicado; withdrawn impede novas inserções, sem apagar blob ou instância antiga. Nova versão volta a draft e não altera anterior.

## Operações manuais/MCP de I2

Todas usam os mesmos casos de uso; GUI apresenta a revisão ao autor e MCP segue HostApproval-v1 para R2. Mutação transitória recebe requestId e revisão esperada da sessão.

| Operação candidata | Classe / entrada | Resultado |
| --- | --- | --- |
| `componente_propor` | R1; origin=selection, requestId, expectedRevision, seleção `{sessionId,screenId,componentIds}`, metadados e parâmetros | proposalId/revision/baseRevision/templateSha256 da captura, diagnóstico de referências e dependências, sem publicar |
| `componente_validar` | R0; proposalId ou definição publicada | validade, suporte, autorização e homologação em dimensões separadas |
| `componente_listar_customizados` | R0; cursor/limit/filtro opcional | versões/estados/compatibilidade e motivos; itens incompatíveis não são inseríveis |
| `componente_instanciar` | R1; sessão/revisão, definitionId/version/packageSha256, screenId, x/y, parameters | mapa de IDs e recursos, operação integral reversível |
| `componente_publicar` | R2; CommitRequest de proposta congelada | versão/index e hash publicados, sem Git |
| `componente_versionar` | R1; identidade/versão-base/nova versão | nova proposta independente; não substitui blob antigo |
| `componente_revisar` | R1 controle; ReviewRequest de library-operations.schema.json | ReviewResponse do contrato de autorização; geração nova dentro dos reviewBytes, sem publicar |
| `componente_retirar` | R2; CommitRequest de revisão da identidade/versão | estado withdrawn; instâncias e conteúdo preservados |
| `componente_preparar_importacao` | R1 edição; PrepareImportRequest de library-operations.schema.json | admite anexo de importação, proposalId/revision/baseRevision/packageSha256 e diagnósticos; sem publicar |
| `componente_importar` | R2; CommitRequest de library-operations.schema.json | confirma exclusivamente importDefinition previamente revisada; conflito não sobrescreve |
| `componente_exportar` | R2; CommitRequest para escrever destino autorizado | cópia dos bytes revisados; sem acesso a arquivos não selecionados |
| `biblioteca_listar_importaveis` | R0; cursor/limit nas import-root selecionadas pelo operador | materialIds/tamanho/hash de pacotes `.jcomp`, sem aceitar diretório arbitrário |

Preparar/revisar publicações/importações/exportações reutiliza Review/Authorization do contrato de consentimento com operation específica; as ferramentas R2 não aceitam conteúdo novo fora da revisão. Preparação é controle R1 com sessionId/requestId/expectedRevision/controlGeneration; não consome quota de edição. Arquivos de entrada são IDs de material autorizado, não caminhos arbitrários enviados pela LLM. Casos de falha, teclado/estados e confirmação visual são definidos no [DESIGN](../DESIGN.md).

Propostas de seleção/versão/importação pertencem à AuthoringSession, com ownership, baseRevision, revisão própria, dirty agregado e limites de anexos definidos em [propostas](propostas-legado.md#sessão-revisões-e-descarte). Todas as edições R1 recebem campos comuns de sessão, incluindo componente_versionar; consulta de proposta identifica sessão/proposalId/revision. Editar um anexo não altera Project.Revision ou undo do documento. Publicar limpa apenas o alvo confirmado; salvar jornada não publica seus anexos.

Strings usam pontos de código Unicode válidos, conforme mcp.md e maxLength do schema, em GUI/MCP/pacote; text/initialValue/default textual até 65.536, nomes até 256, preservados limites específicos menores. Bytes UTF-8 contam somente para frame e pacotes. Não truncar, normalizar ou aceitar Unicode inválido; validar os mesmos valores na entrada manual e MCP.

O codec `.jcomp` de T068/T070, apoiado em 002/T018, é pré-requisito compartilhado da captura validada/publicação US3 (T055/T056) e da transferência US4. Deve produzir/validar bytes finais antes da revisão US3; não adiar sua existência até a entrega da UI de importação/exportação.

## Compatibilidade e validação

Recusar schemaVersion desconhecido, Kind não instalado, propriedades/recursos ausentes, hash divergente, versões colidentes, parâmetros ilegais, captura incompleta, ciclos, template com várias telas e pacotes excessivos. Preservar biblioteca/jornada existentes. Limites e erros seguem mcp.md; nova versão da biblioteca é independente da versão do aplicativo/catálogo/constituição.

Validar paridade do mesmo pacote/versão/parâmetros por GUI e MCP, round-trip sem biblioteca, referências internas, criação/retirada e aparência do legado nas saídas. Não afirmar fidelidade de uma variante sem suporte real no renderer/runtime nem usar uma imagem para simular novo Kind funcional.

## Preparação de revisão e confirmação — U4

Entradas fechadas em [library-operations.schema.json](library-operations.schema.json). Os handlers selecionam o $defs correspondente; o oneOf raiz permite conferir fixtures, não cria uma ferramenta genérica. `componente_validar` continua R0 e não cria Review nem avança geração.

`componente_revisar` exige sessionId/requestId/expectedRevision/controlGeneration/operation/subject. Servidor resolve destinos/bases/hash/bytes a partir do estado autorizado, não de alegações do cliente:

| operation | subject fechado | target | Review e efeito congelado | Confirmação R2 |
| --- | --- | --- | --- | --- |
| publishDefinition | proposalId/proposalRevision de captura selection ou análise legacy com materialização válida | ausente | proposal referencia o anexo; materialization vincula análise/captura quando aplicável; contentHash do pacote final e impactos do índice | componente_publicar |
| importDefinition | proposalId/proposalRevision de anexo de importação admitido | ausente | proposal identifica anexo, materialization=null; pacote original congelado e transição de índice | componente_importar |
| withdrawDefinition | definitionId/version/packageSha256 existente | ausente | proposal/materialization=null; versão/blob ficam intactos; contentHash dos bytes finais do índice, baseHash do índice anterior | componente_retirar |
| exportDefinition | definitionId/version/packageSha256 existente | `{filename}` | proposal/materialization=null; bytes exatos do pacote; destino é filename regular .jcomp sob workspace já autorizado, criação exclusiva | componente_exportar |

Não aceitar campos de outra variante. Filename tem até 120 pontos de código, sem separadores/caminho absoluto, e passa pelas regras Windows de nomes reservados, pontos/espaços finais, arquivo regular e raiz física. Exportação v1 não substitui homônimo; arquivo existente gera conflito sem escrita. Não introduzir nova raiz/flag implícita. Perfil do operador deve autorizar a raiz efetiva de cada operação, inclusive biblioteca privada; falta de escopo/host confiável recusa preparação persistente e commit, preservando o anexo.

Admissão de importação é etapa R1 edição separada: recebe materialId do inventário import-root e expectedPackageSha256; confere origem/hash, ZIP/template/recursos/colisão e quotas antes de criar anexo com os bytes originais congelados. Consome ID de edição e reserva por alvo; replay não cria segundo anexo. Conteúdo idêntico já existente pode retornar identidade/estado existentes sem mudar oferta ou criar duplicata. Se admitido, anexo importado conserva bytes originais; editá-los requer nova definição/versão, não uma importação silenciosamente reescrita. Somente depois `componente_revisar(importDefinition)` prepara o consentimento. Não carregar pacote novo dentro de um controle isento da quota de edição.

Resposta de componente_revisar usa `{ok:true,result:{reviewBytes,reviewHash}}`, com ReviewResponse já definido em authorization.schema.json. Review inclui controlGeneration após a preparação; cliente pode decodificá-la ou consultar sessao_obter. requestId de preparação é do cliente; requestId de commit vem do Review, é gerada pelo servidor e nunca é reutilizada como edição. Uma revisão corrente, replay sem renovar validade, high-water e reservas seguem mcp.md. Repreparar exige geração atual e nova decisão humana.

Os quatro commits recebem **somente** sessionId/requestId/reviewId/authorization, conforme CommitRequest. Comparar operação do Review com a ferramenta chamada, conferir captura/análise quando aplicáveis, pacote, evidências, recursos, base/índice/destino e consentimento sob o gate antes do efeito. Operação trocada ou conteúdo/destino novo é recusado; nenhum bool/texto de aprovação libera gravação.

Em publicação/importação, Review.baseHash registra a precondição do índice; servidor conserva também ausência esperada do blob/identidade-versão, ou integridade do blob idêntico, para transação/recuperação. contentHash é o pacote final; bytes exatos da transição do índice e seus hashes ficam congelados no registro de operação, com impactos descritos na revisão. Exportação usa baseExists=false/baseHash=null e criação exclusiva; retirada congela índice como conteúdo. Alteração concorrente do índice/base/recursos/proposta torna a revisão obsoleta; não regenerar bytes após consentimento. Journal/recovery mantêm a posição reservada e não apresentam publicação parcial como sucesso.

O resumo imutável inclui também indexContentHash da transição de índice quando houver e as precondições de colisão/oferta. O verificador compara esses hashes com o registro congelado antes de aplicar; não recalcular uma transição diferente mantendo a mesma revisão.

GUI usa os mesmos pedidos/estados de aplicação, com confirmação humana da interface. LLM usa as ferramentas e HostApproval-v1; nenhuma modalidade é requisito da outra. Materialização do legado segue [propostas](propostas-legado.md#materialização-explícita-de-proposta-analítica--u3).
