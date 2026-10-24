# Implementation Plan: Domínio Inicialização

Go no módulo existente, composição manual. Novo pacote `internal/startup` para a responsabilidade explicitamente solicitada pelo usuário; nenhuma biblioteca gráfica, COM, filesystem ou mídia pode ser importada por esse núcleo. O adaptador desktop traduz notificações em eventos de domínio na thread de UI. Fonte em `internal/startup/session.go`; regressões proporcionais em `session_test.go`. Implementação de apresentação e reprodução descrita em 036/plan, sem duplicar o núcleo.

Constituição: mantém monólito modular, modelo puro, teste de invariantes e documentação em português. A extensão foi solicitada explicitamente; não cria novo serviço ou formato persistido.

## Política opcional — 1.0.0.10

`startup/preference.go` decide primeira exibição ou repetição explícita; `preference_test.go` cobre primeira execução, conclusão, marcar e desmarcar. `Session` permanece responsável somente pelo caminho com reprodução. O coordenador `desktop.StartDesktop` decide qual caminho executar e continua exigindo os recursos quando não existe apresentação. Persistência local e checkbox pertencem aos adaptadores 036 e launcher 034.
