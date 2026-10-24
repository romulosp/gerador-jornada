# Tasks: teclado numérico Java/EFL

**Reconciliação documental — 2026-10-24 / 1.0.0.1:** a spec identifica código existente e lacunas atuais. Evidências datadas anteriores são históricas; uma referência a arquivo ausente não comprova capacidade atual. Tarefas reabertas conservam a descrição anterior e o motivo. Testes de um recorte não concluem tarefas de expansão ou homologação. Checklists de qualidade e aceite humano não foram marcados por esta revisão documental.


**Date**: 2026-10-19 | **Input**: [spec.md](spec.md), [plan.md](plan.md), [DESIGN.md](DESIGN.md) e contratos compartilhados.
**Estado**: novas tarefas abertas, sem implementação/teste/homologação atribuídos por esta documentação. IDs locais a 013.

- [ ] T001 Reconciliar contrato das 14 teclas/layout/estados/skins/fontes e foco/cursor/TAB, separando legenda LIMPAR, VK_BACK_SPACE e efeito efetivo; registrar fontes/bloqueios C15 e fixtures nova/antiga. FR-003/FR-008–FR-009; usa 001/T108/003-T009, sem duplicar importador.
- [ ] T002 Verificar o slice numérico dos donos 003/T013, 004/T009–T012, 002/T009–T012 e 001/T111–T116, incluindo default novo 200 × 350 e preservação de antigos 200 × 280/IDs/alvos/overrides. FR-001–FR-005/FR-008–FR-009; depende T001 e entregas por lote; não repetir implementação nesses arquivos.
- [ ] T003 Criar/executar teste integrado específico da fixture: arranjo/14 teclas/hit regions, entrada/cursor/setas/TAB/limpar conforme contrato, alvos ausentes/incompatíveis, reset sem edição, v1/v2→v3 sem resize e texto/formas PPT editáveis. Integrar galeria 001/T117/regressão T118 sem duplicar testes gerais de modelo/codec; registrar revisão/resultado. Depende T002.
- [ ] T004 Homologar manualmente o numérico com referência independente no Windows/DPI e PPT das tarefas 001/T120/T121; registrar fonte/skin/efeitos de foco/comandos/recorte, revisão/hash/ambiente/executor e observado. Depende T003, recursos/referências T119 e executável T125 da coordenadora. Falta de referência/delegate mantém bloqueio, não gera aprovação.
- [ ] T005 Obter aceite explícito do usuário da revisão numérica testada, ligado a 001/T122, com cobertura/bloqueios; não aprovar AlphaKeyboard/calculadora ou todo acervo pelo reuso. Depende T004; Java permanece pausado.

T001→T002→T003→T004→T005. Modelo/consumidores têm donos únicos referenciados; tarefas locais detalham contrato/fixture/verificação/manual/aceite. Nada foi marcado concluído nesta revisão.

### Resultado parcial de implementação — 2026-10-24

T001–T003 têm agora evidência técnica para o recorte numeric existente: default 200 × 350, 14 regiões Java, preservação dos bounds antigos, skins normais/auxiliares, inserção no cursor, setas, TAB entre campos habilitados e Backspace sob a legenda LIMPAR. Donos: 001/T148–T150; `presentation/legacy_style_test.go` e `desktop/legacy_preview_test.go`. As caixas locais continuam abertas porque incluem contratos/migrações/perfis mais amplos, alvos compostos e validação completa ainda não demonstrados por esses casos. T004/T005 seguem manuais, sem aceite atribuído.
