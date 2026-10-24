# Tasks: piloto Relatório Java/EFL

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


**Date**: 2026-10-19 | **Input**: [spec.md](spec.md), [plan.md](plan.md), [DESIGN.md](DESIGN.md) e contratos compartilhados.
**Estado**: todas as tarefas abertas inicialmente; nenhuma implementação, execução ou homologação atribuída por esta atualização documental. IDs locais a 012.

- [ ] T001 Reconciliar contrato de cada parte Event (topo/título/corpo/linhas/rodapé/seleção/scroll) e variantes, distinguindo Java/XML/React e flags de fonte declaradas/efetivas; detalhar fixture/medidas/bloqueios C11–C14 na matriz coordenadora. FR-008–FR-009. Usa inventário 001/T108/003-T009; não reimplementa parser/modelo.
- [ ] T002 Verificar integração do piloto editor→prévia→save/reopen v3→PPT editável com os donos 003/T014/T015, 004/T009–T012, 002/T009–T012 e 001/T111–T116; registrar capacidades/divergências por parte, sem marcar implementação alheia por semelhança. FR-008; depende T001 e entregas correspondentes por lote.
- [ ] T003 Criar/executar teste integrado específico da fixture Event com título/linhas Unicode, vazio, seleção/overflow/recorte, preservação de IDs/bounds/overrides e texto/formas PPT editáveis. Integrar à galeria 001/T117/regressão T118, sem duplicar mecanismos de teste de catálogo/ZIP. Registrar revisão/resultados; FR-008–FR-009; depende T002.
- [ ] T004 Homologar manualmente a fixture piloto e variações suportadas nos ambientes de 001/T120/T121, com referência independente e fonte/skin efetivas, revisão/hash, DPI e observado; divergência retorna para dono e reteste. FR-008–FR-009; depende T003, recursos/referências de 001/T119 e executável T125. Outras classes não são homologadas por este piloto.
- [ ] T005 Obter aceite explícito do usuário para a revisão testada do piloto, vinculado a 001/T122; listar cobertura/bloqueios e não anunciar todo Report/acervo aprovado. FR-008–FR-009; depende T004. Java continua pausado.

T001→T002→T003→T004→T005. Implementação pertence aos backlogs referenciados; esta lista detalha somente contrato/fixture/verificação/manual/aceite locais. Evidências futuras serão registradas pelo executor na coordenadora, não criadas ficticiamente aqui.