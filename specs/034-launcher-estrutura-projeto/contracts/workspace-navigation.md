# Contract: Workspace e navegação entre jornadas

Contrato da feature 034. Complementa o contrato canônico `project-format.md` da coordenadora (001); as regras de versionamento/backup permanecem canônicas lá e em 002/003/005.

## 1. Workspace (pasta do sistema)

- O workspace é uma pasta do sistema que contém um ou mais arquivos `.jornada`.
- O nó raiz da árvore representa a pasta (nome = nome da pasta).
- A listagem considera somente arquivos `.jornada` na raiz da pasta (não recursiva), identificados por nome de arquivo relativo.
- Criar/importar adiciona um arquivo `.jornada` à pasta; importar preserva o arquivo de origem (cópia para dentro do workspace).
- A pasta em uso é resolvida pelo adaptador desktop; `Project`/`Document` não armazenam caminhos do sistema operacional.

## 2. Destino entre jornadas (`MenuItem.TargetJourney`)

- `TargetJourney` é um nome de arquivo `.jornada` relativo ao workspace, aplicável apenas a folhas (`Submenu == nil`).
- Exclusividade: `Submenu`, `TargetScreenID` e `TargetJourney` são mutuamente exclusivos — no máximo um por item.
- Validação: nome normalizado válido para o sistema de arquivos; na execução, o arquivo deve existir e ser um `.jornada` legível. Referência ausente/inválida é diagnosticada, nunca substituída ou inferida.
- Versionamento: o novo campo exige `catalogVersion` novo (3 → 4); o writer grava a nova versão e o reader aceita os pares históricos 1/1, 2/1, 3/2 e 3/3 por migração explícita. Primeiro save incompatível avisa e conserva cópia original antes da publicação atômica.

## 3. Navegação na prévia

- Ativar uma folha com `TargetJourney` carrega a jornada destino e mostra a tela inicial dela.
- "Voltar" restaura a jornada de origem (pilha simples); o caminho percorrido é transitório e não é persistido.
- "Prévia Tela Atual" não segue destinos de jornada nem de tela.
- Nenhuma rota é inferida por Label/Shortcut/MovementCode; só destino explicitamente confirmado pelo autor habilita navegação.

## 4. Lista de projetos recentes

- Não é parte do contrato `.jornada`; é preferência local de UI (JSON por usuário), com caminho, nome e timestamp. Não é telemetria nem sincronização.
