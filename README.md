# Jornada Studio — Gerador de Jornada

Editor desktop em **Go 1.27.2/Fyne 2.8.2**, versão **1.0.0.0**, para desenhar telas, conectar navegação e preparar demonstrações de atendimento. Salva projetos portáteis `.jornada` e gera apresentações PowerPoint editáveis com navegação e animações. O editor funciona offline; a ação Java 1.6/EFL está temporariamente desabilitada.

O código da aplicação está em [`apps/gerador-jornada`](apps/gerador-jornada). O [**README completo da aplicação**](apps/gerador-jornada/README.md) descreve propósito, componentes, propriedades, atalhos, arquitetura por domínio, tecnologias, pré-requisitos e comandos de distribuição.

| Material | Local |
| --- | --- |
| Guia completo | [README da aplicação](apps/gerador-jornada/README.md) |
| Builds de 64 bits | [Comandos Windows, Linux e macOS](apps/gerador-jornada/README.md#compilar-em-dist--64-bits) |
| Matriz nativa para geração de pacotes | [Workflow de distribuição](.github/workflows/desktop-distribution.yml) |
| Planejamento GitHub Spec Kit | [Feature coordenadora](specs/001-gerador-jornada/spec.md) |
| Contratos, testes manuais e evidências | [Quickstart](specs/001-gerador-jornada/quickstart.md) |
| Material analisado do sistema anterior | `material-legado/` |

Distribuições são geradas em `apps/gerador-jornada/dist`, por sistema/arquitetura, com versão e SHA-256. A existência de configuração de build não significa que todos os pacotes foram compilados ou homologados: o README explica como consultar os artefatos realmente gerados e o que ainda precisa de validação.
