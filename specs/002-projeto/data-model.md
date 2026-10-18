# Data Model: Domínio Projeto

**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/project`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

Project contém id, name, revision, catalogVersion e assets. Asset contém id, path, mediaType, sha256 e sizeBytes. A aplicação compõe Project com screens, journey, scenarios e javaOptions em Document; o DTO de persistência é achatado conforme o contrato compartilhado. Revisão é incrementada por comandos confirmados de aplicação; operações que falham não publicam nova revisão.

## Invariantes

UUID v4 canônico; nome não vazio; revision inteiro não negativo; unicidade global de IDs; arrays mantêm ordem. Recursos referenciados não são removidos implicitamente. Projeto vazio é salvável, mas não executável. Erro estrutural impede nova sessão; referência quebrada legível permanece reparável.

## Limites e relacionamento com outros domínios

A aplicação resolve referências entre módulos e constrói snapshot. O adaptador projectzip converte DTOs e bytes; Project não importa screen/journey/scenario para executar suas regras nem seus adaptadores. A composição do documento é montada na aplicação, evitando um agregado com todos os comportamentos.

Não implementa filesystem, ZIP, autosave, undo, UI, decodificação de imagens ou geração de artefatos.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: IDs duplicados, nome vazio, revisão negativa, recurso compartilhado e isolamento de snapshot. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/project-format.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-001, FR-009, FR-010; NFR-004, NFR-005 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Nenhuma persistência ou integração foi executada nesta entrega.
