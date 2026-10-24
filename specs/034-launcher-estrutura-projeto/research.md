# Research: Tela inicial (launcher) e Estrutura de Projeto

Registro das decisões de projeto para a feature 034, com racional e alternativas rejeitadas. Nenhuma decisão aqui muda princípios da constituição; todas são extensões autorizadas dos domínios existentes.

## Decisões

### 1. Workspace = pasta do sistema (sem novo arquivo de workspace)

- **Decision**: O "projeto principal" é a pasta do sistema aberta/em uso; a árvore da Estrutura projeto lista os arquivos `.jornada` presentes nela, com o nome da pasta como raiz.
- **Rationale**: Decisão do usuário (pergunta de clarificação). Evita introduzir um novo formato de arquivo e um novo domínio; reutiliza o filesystem e o codec `.jornada` existente.
- **Alternatives considered**:
  - *Workspace persistente (arquivo/pasta com manifesto)*: estrutura mais forte, mas cria um artefato novo e muda mais o modelo — rejeitada pelo usuário.
  - *Arquivo `.jornada` estendido embutindo sub-jornadas*: menos arquivos, porém mexe no formato ZIP atual e no round-trip — rejeitada.

### 2. Lista de projetos recentes = JSON local por usuário

- **Decision**: Manter a lista de recentes em um arquivo JSON local (diretório de configuração por usuário), com caminho do workspace/`.jornada` principal, nome e timestamp da última abertura.
- **Rationale**: É preferência de UI, não dado de negócio; a constituição exige mantê-la fora das entidades e contratos persistidos da jornada. Sem rede/telemetria.
- **Alternatives considered**: registrar no registro do sistema operacional; banco local; reabrir a pasta por busca. Rejeitadas por acoplamento de plataforma ou complexidade desnecessária.

### 3. Destino entre jornadas como novo campo tipado em `MenuItem`

- **Decision**: Adicionar `MenuItem.TargetJourney` (nome de arquivo `.jornada` relativo ao workspace), mutuamente exclusivo com `TargetScreenID` e `Submenu`. Bump de `catalogVersion` (3 → 4), com leitura dos pares históricos e backup antes de save incompatível, no mesmo padrão dos incrementos de códigos de barras/menu (FR-011/FR-012 do 002-projeto).
- **Rationale**: O modelo atual (`TargetScreenID`) referencia uma tela do mesmo documento; jornada é outro arquivo. Um campo explícito mantém a regra "sem inferência por nome/atalho/metadado" (033/005) e a validação tipada.
- **Alternatives considered**:
  - *Referenciar por UUID do projeto-alvo*: mais robusto a renomeação, porém exige abrir cada arquivo para ler o UUID; nome relativo é mais simples e coerente com "workspace = pasta". Escolhido nome relativo; ausência do arquivo = destino ausente, sem substituição.
  - *Reutilizar `TargetScreenID` com um namespace*: ambiguidade e risco de interpretação errada — rejeitada.

### 4. Navegação entre jornadas na prévia

- **Decision**: Ao ativar uma folha com `TargetJourney`, a prévia carrega a jornada destino e mostra a tela inicial dela; um comando "Voltar" explícito retorna à jornada de origem, com pilha simples de jornadas visitadas. A navegação continua proibida na "Prévia Tela Atual".
- **Rationale**: Decisão do usuário. Reutiliza o isolamento de `application.Preview` (snapshot) e adiciona uma pilha transitória de navegação, sem persistir caminho percorrido.
- **Alternatives considered**:
  - *Sessão de prévia separada*: perde o vínculo de origem — rejeitada.
  - *Somente na execução final*: contradiz o pedido de ver o destino ao clicar na prévia — rejeitada.

### 5. Tela inicial substitui o carregamento direto do exemplo

- **Decision**: Sem documento aberto (ou `-open` ausente), o programa mostra a tela inicial em vez de abrir `application.Example()` automaticamente. "Novo com exemplo" continua usando `application.Example()` + fundo padrão; "Novo em branco" usa `application.NewDocument(nome)` + fundo padrão.
- **Rationale**: Alinha com o pedido do usuário e evita criar/descartar documento indesejado; o `Example()` permanece a fonte do "cenário de carregamento atual".
- **Alternatives considered**: manter o exemplo como padrão e mostrar o launcher só ao clicar em "Novo" — rejeitada porque o launcher é a porta de entrada pedida.

## Riscos e pendências abertas

- A migração de catálogo (campo novo) exige atualizar o codec `projectzip` e os contratos em 002/003/005; o plano de tarefas deve incluir esses ponteiros (não duplicar o backlog).
- Nome de item de menu nem sempre é um nome de arquivo válido; normalizar para nome de arquivo sem quebrar o nome de exibição (que vem de `Project.Name`).

## Decisões confirmadas na implementação — 1.0.0.4

- Publicação nova usa temporário com Sync/Close e movimento exclusivo no Windows (MoveFileEx sem REPLACE_EXISTING), preservando homônimos e volumes sem suporte a hard links. Importação lê uma vez, limita a 250 MiB e valida exatamente os bytes copiados.
- Criação pelo menu primeiro valida um snapshot com o vínculo, publica o destino e salva a origem atomicamente. Falha na origem remove apenas o destino recém-publicado; a sessão e o histórico só mudam depois do sucesso. A confirmação informa que alterações pendentes da origem serão salvas. A origem deve estar salva antes de usar o atalho.
- Undo/redo conserva o vínculo como uma alteração de Session; desfazer não remove o arquivo criado, evitando apagar conteúdo autorado depois. Pilha de prévia limitada a 32 visitas, preservando tela, menu e entradas da origem.
- Fyne Tree usa Root explícita `workspace`; conferir só ChildUIDs não prova renderização. A regressão percorre labels produzidas pelo renderer; captura nativa confirma raiz e descendentes.
- Recentes usam até 20 entradas no config por usuário; JSON corrompido é diagnosticado e preservado, sem bloquear criar/abrir projetos. O tema do launcher respeita a preferência existente do editor.

## Falha nativa de transição identificada no uso — 1.0.0.5

O smoke original abria o launcher vazio e o editor com `-open`, sem passar pela troca entre duas janelas durante `App.Run`. O backend de teste também não expõe HWND. A reprodução no driver real mostrou `Editor.New → initTheme → Select.SetSelected → applyTheme → SetTheme → setDarkMode → GetWin32Window(nil)`. A gravação já havia terminado, explicando o projeto nos recentes após reabrir o EXE.

A solução prepara o tema global antes de registrar a nova janela, ignora reaplicação do mesmo tema e instala OnChanged somente depois da seleção inicial. Não é tratamento genérico de panic nem alteração dos arquivos. A regressão unitária observa a quantidade de janelas em cada SetTheme e confere troca explícita de tema após Show. O smoke nativo exercita cliques de criação e reabertura com App.Run ativo. O título recebe a cor do tema diretamente no construtor.
