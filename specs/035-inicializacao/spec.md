# Feature Specification: Domínio Inicialização

**Created**: 2026-10-24
**Status**: Implementação técnica 1.0.0.10; validação final registrada no quickstart; aceite humano pendente
**Input**: Separar o modelo de domínio e a implementação do ciclo de inicialização da apresentação visual da splash.

## Responsabilidade

O domínio governa quando a execução pode sair da inicialização. Não conhece Fyne, Windows, vídeo MP4, relógio, arquivos ou skins. Apresentação e mídia pertencem à [036](../036-splash-screen/spec.md); os domínios de autoria existentes continuam independentes.

## Requisitos

- **FR-001**: Quando a política exigir apresentação, liberar a tela seguinte somente depois de recursos obrigatórios carregados e reprodução realmente iniciada e concluída. Sem apresentação, exigir os recursos e não inventar eventos de vídeo.
- **FR-002**: Quando houver apresentação, aceitar qualquer ordem entre término do vídeo e término do carregamento; uma condição isolada nunca libera a execução.
- **FR-003**: Um tempo estimado ou percentual completo não substitui o evento de término real.
- **FR-004**: Falha de recursos ou reprodução bloqueia o avanço; erro não representa sucesso.
- **FR-005**: Cancelamento impede callbacks tardios de abrir a tela seguinte.
- **FR-006**: Autorizar a transição exatamente uma vez, mesmo recebendo notificações repetidas.
- **FR-007**: Não persistir sessão de inicialização em `.jornada`, recentes, roteiros ou histórico de edição.

## Cenários e aceite

1. Recursos terminam antes: aguardar fim real da reprodução.
2. Vídeo termina antes: aguardar os recursos restantes.
3. Fim de reprodução recebido antes de início: não considerar vídeo concluído.
4. Erro/cancelamento antes ou depois de uma condição: não liberar.
5. Ambas concluídas: primeira solicitação libera; repetidas não liberam novamente.

**Critério técnico**: testes determinísticos das duas ordens, falhas, cancelamento e duplicidade. **Aceite humano** permanece separado dos testes.

## Política por usuário — 1.0.0.10

- **FR-008**: Exibir na primeira execução; repetição desabilitada por padrão. Depois da primeira apresentação concluída, exibir somente por opção explicitamente habilitada. Desmarcar suprime as próximas apresentações.
- **FR-009**: Registrar primeiro uso somente após recursos e vídeo concluídos; cancelamento/falha permite tentar a primeira apresentação novamente.
- **FR-010**: Preferência independente de versão do software/projeto; persistência é responsabilidade do adaptador. O domínio não acessa arquivos.

Cenários: configuração ausente → mostrar e repetição falsa; conclusão → não mostrar por padrão; marcar → mostrar novamente; desmarcar → não mostrar. Sem vídeo, o carregamento permanece obrigatório e falhas impedem a tela seguinte.
