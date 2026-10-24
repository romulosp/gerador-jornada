# Data Model: Tela inicial (launcher) e Estrutura de Projeto

Modelo de dados da feature 034. Extensões pontuais sobre os modelos canônicos dos domínios Projeto/Jornada/Catálogo; os invariantes existentes permanecem e são referenciados, não duplicados.

## Entidades

### Workspace (pasta do sistema)

- **Caminho**: caminho absoluto da pasta aberta/em uso.
- **Rótulo raiz**: nome da pasta (exibido como "projeto principal" na árvore).
- **Jornadas**: lista dos arquivos `.jornada` presentes na pasta (não recursiva), cada um identificado por nome de arquivo relativo.
- **Regras**: não é uma entidade persistida de negócio; é resolvido pelo adaptador desktop. A lista é recalculada ao abrir/criar/importar/remover uma jornada.

### Jornada

- Reutiliza `application.Document` (Project com UUID/nome/revisão, Screens ordenadas, Journey, Scenarios, Resources) persistido como `.jornada` (ZIP).
- **Nome de arquivo**: nome relativo ao workspace (normalizado para o sistema de arquivos); o nome de exibição continua sendo `Project.Name`.
- **Relação**: descendente do Workspace na árvore; pode ser destino de uma folha de menu (via `TargetJourney`) ou ser a jornada de origem da navegação.

### Projeto recente

- **Campos**: caminho do workspace (ou da `.jornada` principal), nome de exibição, timestamp da última abertura.
- **Persistência**: JSON local por usuário (preferência de UI; fora do contrato `.jornada`).
- **Regras**: entrada indisponível (movida/excluída) é indicada sem quebrar as demais.

### Destino entre jornadas (novo campo)

- **Campo**: `MenuItem.TargetJourney string` — nome de arquivo `.jornada` relativo ao workspace.
- **Exclusividade**: mutuamente exclusivo com `TargetScreenID` e `Submenu` (folha: `Submenu == nil`). No máximo um dos destinos por item.
- **Validação**: nome normalizado válido; o arquivo deve existir e ser um `.jornada` legível quando a jornada é executada. Ausência = destino ausente, sem inferência nem substituição.
- **Versão**: exige `catalogVersion` novo (3 → 4), com leitura dos pares históricos e backup antes de save incompatível.

### Navegação na prévia (transitório)

- **Pilha de jornadas**: sequência de jornadas abertas na prévia, com a tela corrente de cada uma; "Voltar" remove o topo e restaura a origem. Não persistida; caminho percorrido nunca entra no documento ou no histórico do editor.
- **Entrada**: carregar a jornada destino e entrar na tela inicial dela (`Journey.StartScreenID`), mantendo o isolamento de `application.Preview`.

## Transições relevantes

### Criação rápida a partir de folha do menu

1. Folha sem destino é ativada na prévia → `MenuSimulated`.
2. Autor aceita a sugestão → nasce um `.jornada` em branco (fundo padrão) no workspace, com `Project.Name` = rótulo da folha.
3. A folha recebe `TargetJourney` = nome do arquivo criado.
4. A árvore da Estrutura projeto reflete a nova jornada.
5. Nova ativação da folha → navegação para a jornada criada.

### Cancelamento

- Cancelar a sugestão/criação não publica estado parcial: nenhum arquivo é criado e `TargetJourney` não muda; histórico (undo/redo) preserva o último estado confirmado.

## Invariantes preservados (referências)

- IDs locais de feature; referências entre features por caminho (constituição VII).
- Menu sem destino não infere rota por Label/Shortcut/MovementCode (033/005).
- Prévia Tela Atual não segue destinos (inclui os de jornada).
- `Project`/`Document` não carregam caminhos do sistema operacional (002/FR-007).

## Opção visual da splash

`Launcher.ShowSplash` é estado de UI, vinculado à Preference definida em 035. Não integra Workspace, RecentProject ou Document. Persistência e campos JSON são definidos em 036/data-model.md; versões do formato permanecem 3/4.
