# Contract: serviço de exportação

**Version**: 1 | **Requirements**: FR-026, FR-028, NFR-004, NFR-005.

## Estado de implementação — 2026-10-18

Este contrato define o alvo de evolução; não descreve uma API já integralmente implementada.
Hoje desktop.startExport coordena um job/snapshot e callback `func(float64)`, chamando
pptx.Export ou javaefl.Export com contexto. PPT publica por atomicfile; Java tem staging
e gate de perfil. Ainda faltam o caso de uso ExportProject na aplicação, request/result
estruturados, Artifact, eventos por fase e todas as localizações de Diagnostic.
Pendências mantidas em 001-gerador-jornada/T009, T030, T038, T040, T046 e T068; não considerar o contrato
satisfeito apenas por existir exportação funcional. Requisitos abaixo permanecem vigentes.

## Fronteira interna

Não há API HTTP pública nesta versão. O caso de uso ExportProject recebe:

| Entrada | Conteúdo |
| --- | --- |
| context | Cancelamento da operação, independente da janela. |
| snapshot | Projeto, revisão, catálogo/perfil resolvidos e bytes dos recursos, imutáveis. |
| target | `pptx` ou `javaEfl`; cada job tem um alvo. |
| options | PPT: includeExploration e scenarioIds; Java: packageName e profileId resolvidos. |
| destination | Caminho escolhido no diálogo, fora do modelo de domínio. |
| progress | Callback/canal de eventos; não manipula widgets. |

Exportador produz artefatos relativos em staging e diagnósticos; o caso de uso publica a saída. Interface no consumidor com uma operação Export, usando structs concretas para request/result. Não adicionar uma interface para cada tipo de dados.

Somente um job ativo por sessão na v1. Uma segunda solicitação deve informar que há trabalho em andamento e permitir cancelar o atual; nunca sobrescrever seu contexto. Edição pode continuar; progresso identifica a revisão do snapshot.

## Resultado e diagnóstico

Result contém status `succeeded`, `failed` ou `cancelled`, snapshotRevision, target, artifacts e diagnostics. artifacts só contém saídas publicadas no sucesso. Cada Artifact contém relativePath, mediaType, sizeBytes e sha256.

Diagnostic contém code estável, severity `error`/`warning`, message em português, target e location. Location tem campos opcionais screenId, componentId, scenarioId, stepId e property. Mensagem ajuda a corrigir, sem transcrever valores sensíveis completos.

Erro de domínio pode gerar vários diagnósticos numa validação. Erro técnico retorna também error contextualizado para logging local. Cancelamento não é apresentado como corrupção ou defeito do projeto.

## Pipeline

1. Validar estrutura, catálogo, recursos e conteúdo requerido pelo alvo.
2. Resolver plano visual/eventos ou perfil Java; verificar capacidades antes de renderizar.
3. Gerar em staging no mesmo volume do destino.
4. Conferir estrutura e integridade dos artefatos.
5. Publicar saída completa. Para arquivo PPT, substituição atômica; para Java, novo diretório completo.
6. Limpar staging e concluir. Se falhar antes de publicar, manter destino anterior.

PPT pode substituir arquivo existente após confirmação no diálogo. Java só exporta para diretório novo/inexistente; se já existir, pedir outro nome, sem mesclar ou substituir fontes editados.

Publicação é a fronteira de commit: cancelamento confirmado antes dela impede publicar. Se já ocorreu commit, reportar sucesso, não um cancelamento falso. Progresso não deve anunciar 100% antes da publicação.

## Progresso e responsividade

Event contém phase, completed, total e message curta. Fases: validating, generating, verifying, publishing. Pode ser indeterminado na validação; percentuais são derivados somente quando total é conhecido.

Atualizar por componente/slide/arquivo, agregando notificações para evitar saturar a UI. Verificar context entre unidades de trabalho e laços longos, com resposta ao cancelamento em até 1 s no conjunto nominal.

Não passar `context.Background()` em cascata ignorando o contexto do job. Não compartilhar ponteiros mutáveis da sessão com workers.

## Determinismo

Ordenar telas/componentes conforme projeto, imports Java lexicalmente, arquivos por caminho e IDs Open XML por travessia determinística. Normalizar timestamps ZIP; não inserir data corrente nos fontes. A mesma revisão/opções deve produzir conteúdo equivalente; temporários e destino absoluto não fazem parte da saída.

## Testes

Injetar falhas de escrita, validação, publicação e cancelamento. Provar ausência de artefatos publicados em falha, preservação do destino anterior, snapshot estável durante edição e mensagens de localização. Validar cada alvo também pelos contratos específicos.
