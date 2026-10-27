# Specification Quality Checklist: SPEC 038

**Purpose**: Verificar redação, coerência e estado real dos apontamentos e do planejamento; registrar somente aprovações efetivamente recebidas, sem equipará-las a implementação ou aceite funcional.
**Created**: 2026-10-25. **Revisão:** v0.4 em 2026-10-25.
**Feature**: [SPEC 038](../spec.md). Base: parecer técnico de correção fornecido pelo usuário e orientações sobre modos manual/MCP, legado e domínio.

## Content Quality

- [x] Seções obrigatórias preenchidas e IDs FR-001–FR-039, US1–US5 e SC-001–SC-010 preservados; SC-011–SC-013 acrescentados sem renumerar anteriores.
- [x] Modo manual e MCP com LLM externa explicitados como independentes e complementares; um não substitui o outro.
- [x] Estilo do legado, aparência, recursos/defaults e paridade vinculam documentação e implementação futura.
- [x] Modelo dividido nos cinco domínios existentes; adaptadores não possuem modelo de negócio paralelo.

## Requirement Completeness

- [x] FR-011 e FR-034 definem consentimento verificável sem GUI, recusa/correção, revisão vinculada e falha segura.
- [x] FR-027–FR-030 separam interpretação do agente, validação determinística do Go e decisões do usuário.
- [x] Fluxos declarativo e funcional distinguidos; proposta válida não significa equivalência comprovada, autorização ou Kind instalado.
- [x] Consentimento classifica R0–R3 e recusa autoaprovação, arquivo gravável pelo agente e token sem origem independente.
- [x] Cenários incluem negação, falta de canal, expiração/revogação, escopo, reutilização e alteração concorrente, sem alegar execução.
- [x] DESIGN cobre biblioteca vazia, detalhes/estados, importação, incompatibilidade, recursos ausentes, inserção, revisão, publicação e versão.
- [x] DESIGN 023/028 referenciam o fluxo único da 038 sem criar outra biblioteca ou duplicar regras de domínio.
- [x] Emenda constitucional contém texto, justificativa e impacto; aprovação explícita registrada e versão vigente atualizada para 1.4.0 sem alterar princípios I–VII.

## Coerência documental e limites

- [x] Spec, research, modelo conceitual, contratos comportamentais, DESIGN e roteiro de aceitação refletem as mesmas responsabilidades e limites.
- [x] Critérios de sucesso e casos de aceitação contemplam falsa aprovação, manual/MCP, visual legado e divisão por domínio.
- [x] Pendências técnicas têm destino explícito no planejamento; ausência de mecanismo/esquema/limite não autoriza comportamento irrestrito.
- [x] Proposta, decisão pendente, definição técnica, implementação, teste e aprovação são estados distintos.

## Matriz de conformidade dos apontamentos

| Item | Evidência documental | Estado e condição restante |
| --- | --- | --- |
| P0 — texto da exceção | [research](../research.md#proposta-de-exceção-constitucional) e constituição | Emenda 1.4.0 aprovada/aplicada em 2026-10-25 |
| P0 — limites/offline/impacto | Texto e relatório de impacto no research; FR-001/FR-038 | Corrigido documentalmente; sem prova técnica de software |
| P0 — versionamento | 1.3.0 base histórica; 1.4.0 vigente, expansão compatível MINOR | Decisão explícita recebida e relatório de impacto atualizado; não autoriza código |
| P1 — responsabilidades | FR-027–FR-030, US5 e [contrato](../contracts/propostas-legado.md) | Corrigido documentalmente; implementação e testes futuros |
| P1 — consentimento | FR-011/FR-034, US2, SC-011 e [contrato](../contracts/consentimento.md) | HostApproval-v1/Ed25519 desenhado; implementação/host real não validados |
| P1 — alternativas/ameaças | research e contrato tratam host MCP-only, pin de chave e alternativas | Modelo de confiança/limites definidos no plano; prova real pendente, R2 falha fechado |
| P2 — experiência manual | [DESIGN](../DESIGN.md), 023/028, US3/US4 e SC-012 | Desenho proposto; implementação e homologação não executadas |
| Requisitos do usuário | FR-001/FR-002/FR-032, US1 e [modelo](../data-model.md) | Ambos os modos, legado e domínios vinculam spec e implementação futura |

## Gates e trabalho ainda pendentes

- [x] Usuário aprovou explicitamente emenda e planejamento: resposta “Aprovar emenda e planejamento”, em 2026-10-25.
- [x] Emenda 1.4.0 aplicada e relatório de impacto/metadados/referências desta feature sincronizados.
- [x] SPEC 038 v0.4 autorizada para planejamento pela mesma resposta explícita; nenhuma autorização de implementação.
- [x] Plano técnico produzido com mecanismo, perfil de host, dados, fronteiras, alternativas e gates de viabilidade definidos; não demonstra execução.
- [x] Contratos técnicos/esquemas, pacote/parâmetros, erros e limites concretos documentados; verificação de handlers/codec ainda futura.
- [x] Usuário invocou explicitamente speckit-tasks, adotando o plano disponível para decomposição documental em 2026-10-25; isso não é aprovação técnica de software.
- [ ] Host real foi integrado/qualificado e fronteira de consentimento demonstrada; nenhum cliente real foi homologado.
- [ ] Substituição de arquivo sob interferência externa foi comprovada no EXE Windows; sem essa prova a capacidade deve ficar bloqueada.
- [x] Tasks geradas por essa autorização, com dependências, UI manual, testes de falsa aprovação, domínios e gates; todas as 115 tarefas novas continuam abertas.
- [x] Análise spec/plan/tasks inicial executada: cinco achados U1/I1/I2/I3/U2, sem crítico; não representou liberação para implementar.
- [x] Replanejamento documental dos cinco achados realizado pela invocação de speckit-plan; contratos/modelo/aceitação e tarefas existentes sincronizados, sem alterar IDs ou checkboxes de implementação.
- [x] Reanálise das correções iniciais executada; identificou U3/U4 HIGH e I4 LOW adicionais, sem crítico.
- [ ] Análise subsequente ao segundo replanejamento solicitada para a sequência; relatório somente na conversa, sem antecipar resultado neste checklist.
- [ ] Implementação dos dois canais e paridade visual/domínio realizada; não autorizada nesta etapa.
- [ ] Testes automatizados/cliente real/EXE, provas de consentimento, regressão, DPI e Office aplicável executados; não executados.
- [ ] Revisão executável, comandos reais e evidências registrados em quickstart; roteiro atual contém somente cenários futuros.
- [ ] Usuário homologou a revisão funcional; nenhuma aprovação recebida.

## Rastreio do replanejamento documental

| Achado da análise | Correção proposta e evidência | Estado restante |
| --- | --- | --- |
| U1 — quota bloqueia finalização | research D6/D10, mcp/consentimento/modelo; T006/T009/T026/T028/T033/T037/T041; US2.9 e Q038-21 | Edição/controle separados, reserva por alvo e recovery na posição original; prova funcional futura |
| I1 — publicação precede codec | research D7, biblioteca/plano e DAG; T055/T056 dependem de T070; T068/T070/002-T018 antecipados; US3.8 e Q038-22 | Dependência documental corrigida; codec/publicação não implementados |
| I2 — bytes versus caracteres | research D8, mcp/biblioteca/schemas/modelo/DESIGN; T008/T010 e 003/T026; Q038-23 | Pontos de código comuns, frames em bytes; guardas Unicode/entradas reais ainda futuras |
| I3 — defaults sem representação | research D9, proposal.schema.json/propostas/modelo/DESIGN; T079/T082/T083; US5.7 e Q038-24 | Estados e tipos explícitos, limites de coleção e lacunas declarados; validação semântica/homologação futura |
| U2 — R1 de proposta sem sessão/revisão | research D10, propostas/biblioteca/mcp/modelo/authorization.schema.json; T055/T079/T082/T084; US5.8 e Q038-25 | Mesmo coordenador, revisão própria, ownership e dirty agregado; aplicação futura |

As correções são de desenho, não encerramento automático da análise nem aprovação da implementação. Nenhuma emenda adicional é proposta; a 1.4.0 já aprovada permanece vigente. Gates EXE Windows, proteção contra escritor externo, host real e paridade visual continuam abertos.

## Segundo replanejamento — U3/U4/I4

A análise seguinte identificou U3/U4 HIGH e I4 LOW, sem CRITICAL. O usuário autorizou speckit-plan e speckit-analyze em sequência; esta seção registra o desenho corrigido antes da análise subsequente, sem antecipar seu resultado.

| Achado | Ajuste documental | Evidência e comprovação futura |
| --- | --- | --- |
| U3 — análise sem conteúdo publicável vinculado | materialization null ou captura selection/revisão/templateSha256, validação do servidor e vínculo no Review; análise não constrói template nem reserva segundo commit | research D11, proposta/consentimento/biblioteca/modelo/DESIGN; T079/T082–T087; US5.9 e Q038-26 |
| U4 — preparação da revisão não alcançável por MCP | componente_revisar/ReviewRequest e quatro commits fechados; admissão componente_preparar_importacao separada como R1 edição | research D12, library-operations.schema.json e contratos; T003/T056/T059/T071–T073/T076; US3.9 e Q038-27 |
| I4 — exemplo de paralelismo antigo | exemplo US4 exige T070 e não condiciona T068 a T057; dependências locais/IDs preservados | tasks, seção Parallel Examples; codec antecipado continua pré-requisito de US3 |

Quatro schemas de planejamento; fixture/schema não prova handler, assinatura confiável, GUI ou executável. A próxima análise será somente leitura e terá relatório na conversa; sua conclusão não será marcada antecipadamente neste checklist. Nenhuma tarefa de implementação foi executada e nenhuma autorização de desenvolvimento foi recebida.

## Limites da evidência

Os itens de qualidade marcados representam verificações documentais; as três aprovações/aplicações registradas em gates correspondem à resposta real do usuário e à emenda efetuada. O [quickstart](../quickstart.md) registra todos os cenários como não executados. A constituição ativa é 1.4.0: conflito de escopo resolvido, mas conformidade do software depende de implementação/testes futuros. Não declarar autorização de implementação ou aceite funcional.

Pesquisa e plano foram decompostos em [tasks.md](../tasks.md) mediante a invocação explícita dessa fase pelo usuário. Escolhas de desenho e checklist sintático não comprovam mecanismo real, Windows, fronteira de SO ou cobertura semântica da análise. Nenhuma tarefa de código/build/teste ou homologação foi executada. Nenhuma marca de qualidade substitui aprovação humana, teste ou evidência executável.
