# Contract: Vínculo e publicação Git

**Cobertura**: spec 037 / FR-005–014 e FR-018. Preservar [workspace e destinos da 034](../../034-launcher-estrutura-projeto/contracts/workspace-navigation.md), sem mudança no formato `.jornada`.

## Vínculo

O seletor entrega o destino por valor, sem redigitação em outro formulário. Após vínculo bem-sucedido, abrir revisão; commit/push só após confirmação em Publicar no GitHub. Colisão retorna código seguro `collision` e diagnóstico de fase `workspace_integration`; falha impede abrir a revisão automaticamente.

1. Canonicalizar raiz exata. Bloquear Git em ancestral, bare, linked worktree, sparse checkout, submódulos, operação Git incompleta e índice com staging preexistente. Git existente exige branch padrão e origin compatível; não alterar config. Origin SSH exige ajuste externo explícito; persistir somente HTTPS sem userinfo.
2. Revalidar ID/branch via REST. Remoto vazio usa default_branch; se omitida, `main`, exibida antes da confirmação. Fetch privado e inventário completo antes de materializar.
3. Validar nomes Windows, colisões de caixa, arquivo/diretório, traversal, `.git`, symlinks/junctions, submódulos e ponteiros LFS. Bloquear conteúdo não suportado sem executá-lo ou aplicar hooks/filtros.
4. Sem Git local: árvore remota completa é base, preservada nos objetos/commits. Materializar somente `.jornada` regulares da raiz; homônimos de jornadas só são compatíveis se tipo/bytes idênticos. Diferença bloqueia todo vínculo, sem sobrescrever. Demais arquivos e subpastas remotos não são copiados para o workspace e não são excluídos do remoto. Jornadas locais exclusivas ficam não publicadas. O índice completo usa skip-worktree nos caminhos não elegíveis, sem aplicar sparse checkout de terceiros.
5. Preparar metadados privadamente e aplicar com journal, branch/índice na base remota e origin sem token; gravar vínculo por último. Não criar commit ou push. Arquivos locais anteriores nunca são removidos para acomodar remoto.
6. Git local existente exige ancestralidade compatível. Sem ancestral comum ou commits locais adiantados desconhecidos, bloquear e orientar resolução externa; nunca enviar histórico não revisado implicitamente.

## Integrar novidades remotas

Fetch antes de revisar e novamente ao confirmar. Se o estado mudar, pedir nova revisão antes de commitar. Integração materializa somente jornadas da raiz com preflight e journal. Arquivos não elegíveis existentes localmente ficam intocados, mesmo quando o remoto os altera ou exclui. Remoto vazio retorna OID ausente, nunca ref antiga em cache como prova de entrega.

| Relação | Comportamento |
|---|---|
| HEAD = remoto | Usar base atual. |
| HEAD ancestral do remoto | Fast-forward com preflight de mudanças locais; mesmo caminho local divergente bloqueia. |
| Remoto ancestral de HEAD | Reenviar apenas pendentes comprovados, preservando OID. |
| Divergentes com uma base comum | Comparação conservadora por caminho e commit com dois pais se não houver conflito. |
| Sem base ou múltiplas bases | Bloquear integração e orientar resolução externa. |

Comparar base B, local L e remoto R por presença/tipo/modo/OID: L=B adota R; R=B adota L; L=R aceita resultado; demais casos conflitam. Ausência inclui exclusões; não inferir rename ou merge textual de ZIP. Validar conflitos arquivo/diretório e representação no Windows também sobre o resultado.

No caso divergente, construir árvore preservando remoto e só deltas locais elegíveis. Exigir que os commits locais sejam comprovados no journal e restritos a jornadas. Criar merge com primeiro pai R e segundo L; o primeiro pai remoto é escolha deliberada para auditar que o delta contra R contém só `.jornada` da raiz. OIDs locais são preservados como ancestrais, sem rebase. Esse algoritmo usa plumbing próprio, não pressupõe merge geral da biblioteca.

Antes de alterar arquivo da jornada ativa, salvar edição e confirmar recarga; cancelar impede integração. Jornada inválida recebida permanece em disco com diagnóstico, sem substituir sessão válida automaticamente.

## Publicação

**Primeiro envio**: HEAD e referência remota ausentes são representados por string vazia. Usar árvore vazia, criar commit raiz com zero pais e publicar a branch padrão explicitamente, sem chamar MergeBase com OID nulo. AccountID da revisão deve coincidir com a conta revalidada antes de qualquer mutação.

1. Resolver edição não salva por callback de conclusão/cancelamento/erro. Revisão imutável contém HEAD, conta, destino e hashes dos bytes. Listar todas as inclusões/alterações/exclusões `.jornada` regulares da raiz; sem seleção parcial. Exclusão só para rastreados; alertar referências afetadas. Mensagem padrão `Atualiza jornadas do workspace`, editável e não vazia; mostrar autoria.
2. Pending OID existente é reenviado primeiro; novas edições ficam para próxima publicação. Se for preciso integrar remoto, apresentar nova revisão para o merge antes de confirmar.
3. Bloquear índice com staging alheio. Construir árvore baseada em HEAD mais os blobs revisados; sem add recursivo, All, filtros ou conversão de bytes ZIP. Auditar diff antes de mover ref. Arquivos não elegíveis são herdados sem entrar no delta.
4. Sob trava, conferir root, HEAD, índice, conta e hashes; alteração externa devolve staleReview. Registrar intenção, pais e OID calculado antes de atualizar HEAD por comparação do valor anterior. Worktree/índice/ref devem ter recuperação documentada para cada ponto de falha.
5. Push com refspec explícito da branch padrão, sem tags/force; token só no transporte. Confirmar OID alcançável na branch remota antes de marcar publicado.

Histórico remoto pode conter outros arquivos; a restrição vale para deltas produzidos/enviados pela ferramenta, sem apagar conteúdo herdado. Commits locais alheios ou desconhecidos bloqueiam envio automático.

## Falhas e recuperação

- Mutex por raiz e trava entre processos da aplicação. Checar locks Git existentes. Salvar/criar/importar aguardam a breve aplicação transacional; depois de falha remota voltam a funcionar. Cliente Git externo não respeita trava privada: revalidar snapshots e refs, usar criação exclusiva e interromper se houver interferência.
- Journal salva pré-imagens/hashes e fases antes de cada mutação. Restaurar caminho somente se ainda contiver os bytes que a operação gravou; se houve mudança externa, preservar backups e diagnosticar recuperação assistida.
- Erro de push mantém pending OID. Fetch: OID ancestral do remoto = entregue; remoto ancestral do OID = reenviar; divergência = integrar ou conflito; sem rede = pendência não verificada. Não recriar commit em cada tentativa.
- Reinício após intenção: HEAD planejado comprova commit; HEAD anterior permite retomar/abortar passos locais; outro HEAD bloqueia recuperação automática. Não reconhecer commit apenas por nome/autor.
- Nunca reset hard, apagar `.git` existente ou sobrescrever arquivo anterior para recuperar. Logout cancela rede e remove credenciais, preservando vínculo, histórico e pendências; não afirma revogar o OAuth App no site.

## Operações do coordenador

**Implementação atual**: `application.WorkspaceSync` expõe Bind/ReviewPublication/Publish e coordena `GitHubConnection`; Engine fornece Binding/Bind/Review/Publish/Recover/LocalChanges. Revisões são valores revalidados, não IDs opacos expostos pela UI. Os nomes abaixo representam o desenho inicial das capacidades, não uma API HTTP externa.

`InspectWorkspace(ctx, root)`, `PlanBinding(ctx, root, repositoryID)`, `ApplyBinding(ctx, planID)`, `ReviewPublication(ctx, root)`, `Publish(ctx, reviewID, message)`, `RetryPublication(ctx, root)` e `Status(ctx, root)`. IDs de revisão/plano são opacos e vinculados ao snapshot. Status local não exige rede; consulta remota é explícita.

Erros adicionais: incompatibleRepository, unsupportedGitFeature, collision, conflict, dirtyIndex, staleReview, externalModification, recoveryRequired, pushRejected, outcomeUnknown. Retornar estado preservado e ação possível, sem executar orientação externa automaticamente.
