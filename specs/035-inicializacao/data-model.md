# Data Model: Inicialização

## startup.Session

Estado transitório em `internal/startup/session.go`, sem dependências de adaptadores.

| Campo | Significado |
| --- | --- |
| Loaded | Preparação obrigatória concluída com sucesso |
| Playing | Evento real de início recebido |
| VideoEnded | Evento real de fim recebido após início |
| Err | Primeira falha; impede liberação |
| Closed | Inicialização cancelada/fechada |
| Transitioning | Autorização já consumida; evita repetição |

Eventos: `ResourcesCompleted(error)`, `PlaybackStarted()`, `PlaybackCompleted()`, `Fail(error)`, `Close()`. Consulta com efeito: `BeginTransition()` autoriza uma única vez se Loaded, Playing e VideoEnded, sem erro/cancelamento. Eventos são serializados pelo coordenador; o domínio não faz sincronização nem cria goroutines. Não há entidade de arquivo, widget, áudio, duração ou caminho neste modelo.

## startup.Preference

Modelo puro em `internal/startup/preference.go`: `Presented` indica apresentação concluída, `ShowOnStartup` representa a escolha explícita, inicialmente falsa. `ShouldShow()` é `!Presented || ShowOnStartup`. Não contém caminhos, versão de build nem identificador de projeto. É dado por usuário, persistido pelo adaptador 036 fora do `.jornada`; a sessão de reprodução continua transitória.
