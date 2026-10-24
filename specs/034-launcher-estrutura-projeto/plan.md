# Implementation Plan: Tela inicial (launcher) e Estrutura de Projeto

**Branch**: `034-launcher-estrutura-projeto` | **Date**: 2026-10-24 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/034-launcher-estrutura-projeto/spec.md`

## Summary

Estender o editor desktop (módulo Go `gerador-jornada.local/editor`, UI Fyne) com: (1) uma tela inicial de boas-vindas com projetos recentes e criação de projeto em branco/com exemplo; (2) a aba "Estrutura projeto" no painel esquerdo, listando em árvore os arquivos `.jornada` da pasta do sistema em uso; e (3) a criação rápida de uma jornada a partir de uma folha do Menu aplicação sem destino, com navegação entre jornadas na prévia (substituir prévia + botão voltar com empilhamento simples). O workspace é a pasta do sistema (decisão do usuário), sem novo formato de arquivo de workspace.

## Technical Context

**Language/Version**: Go (módulo único `gerador-jornada.local/editor`), Fyne v2.

**Primary Dependencies**: Fyne v2 (janela/controles), biblioteca padrão Go (`os`, `path/filepath`, `encoding/json`). Sem dependências de rede.

**Storage**: Sistema de arquivos. Jornada = arquivo `.jornada` (ZIP, codec `projectzip`; writer implementado schema 3/catalog 4, leitor reconhece 1/1, 2/1, 3/2, 3/3 e 3/4). Workspace = pasta do sistema com arquivos `.jornada`. Lista de recentes = arquivo JSON local por usuário (preferência de UI, fora do contrato da jornada).

**Testing**: `go test ./...` (stdlib); os checks Python existentes (java6, packaging) não são afetados nesta entrega. Homologação manual registrada em quickstart.md.

**Target Platform**: Windows 10/11 x64, desktop offline, janela nativa Fyne.

**Project Type**: desktop-app (módulo Go monólito modular).

**Performance Goals**: Tela inicial e árvore do workspace responsivas (< 1 s) para workspace com até ~100 `.jornada`; navegação entre jornadas equivalente à troca de tela existente.

**Constraints**: Offline, sem telemetria/rede; UI em português; preservar round-trip do `.jornada`; não criar microserviços/banco/domínio novo; limites de menu existentes (16 níveis, 1000 itens).

**Scale/Scope**: Usuário único; workspace com dezenas a centenas de jornadas; menu aplicação com os limites já validados.

## Constitution Check

*GATE: passa antes da Fase 0. Re-checar após a Fase 1.*

| Princípio | Avaliação |
|---|---|
| I. Domínios explícitos | PASS — estende Projeto (workspace) e Jornada (navegação entre jornadas); tela inicial/aba são componentes de interface, sem domínio novo. |
| II. Modelo compartilhado + contratos | PASS — destino entre jornadas é explícito e tipado; recentes/workspace não entram no documento persistido. |
| III. Distribuição offline | PASS — tudo local; nenhuma rede. |
| IV. Documentação em português | PASS — artefatos e mensagens em PT. |
| V. Compatibilidade + testes | PASS — bump de catálogo para o novo campo, com leitura dos pares históricos e backup antes de save incompatível (mesmo padrão de FR-011/FR-012 do 002). |
| VI. Simplicidade | PASS — sem novos serviços; pasta como workspace e JSON local de recentes. |
| VII. Rastreabilidade por domínio | PASS — alterações de formato referenciam 002-projeto/003-catalogo/005-jornada sem duplicar backlog. |

Nenhuma violação a justificar.

## Project Structure

### Documentation (this feature)

```text
specs/034-launcher-estrutura-projeto/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── workspace-navigation.md
└── tasks.md   # fase 2 ($speckit-tasks)
```

### Source Code (repository root)

```text
apps/gerador-jornada/
├── cmd/gerador-jornada/main.go          # inicia a tela inicial quando não há documento
└── internal/
    ├── application/                      # Document, Session, Preview, Example
    ├── journey/                          # Journey, Action (none/navigate)
    ├── catalog/                          # MenuItem, ApplicationMenuProperties
    ├── project/                          # IDs, validação
    └── adapter/
        ├── desktop/                      # Editor, files.go, tools.go, preview_compounds_view.go
        │   ├── launcher.go               # (novo) tela inicial + recentes
        │   ├── workspace.go              # árvore e ações do workspace
        │   ├── workspace_store.go        # resolução, validação e publicação exclusiva
        │   ├── workspace_preview.go      # confirmação, vínculo e rollback
        │   ├── recents.go                # JSON local por usuário
        │   └── ...
        └── projectzip/                   # codec .jornada (ZIP)
```

**Structure Decision**: manter o monólito modular existente. O workspace (pasta) e a lista de recentes são responsabilidades do adaptador desktop (filesystem), sem novo pacote de domínio. O novo campo de destino entre jornadas vive em `catalog.MenuItem` e é validado/serializado pelos mesmos domínios, com bump de `catalogVersion`.

## Complexity Tracking

> Sem violações a justificar.

## Opção de apresentação — 1.0.0.10

`Launcher.ShowSplash` é um widget.Check abaixo de abrir projeto, construído pelo adaptador `startup_preferences.go`. Ler antes de associar OnChanged evita gravar no construtor. Usar update atômico do campo, reverter em falha e preservar primeiro uso. Modelo puro e startup coordenados por 035/036, fora das entidades de autoria.
