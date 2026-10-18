# Data Model: Domínio Catálogo

**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/catalog`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

ComponentDefinition mantém kind, displayName, defaultBounds, propertyDefinitions e capabilities. Kind é enum estável. Propriedades específicas ficam em structs por família; Option/Row têm IDs e texto, sem separação por vírgulas.

## Invariantes

catalogVersion é 1. Enums devem pertencer à lista do contrato; fontSize positivo; cor #RRGGBB; radio tem ao menos uma opção e selectedOptionId existente ou nulo; formatos não implicam validação CPF/data. Teclado aponta a campo da mesma tela, validado pela aplicação.

## Limites e relacionamento com outros domínios

Catálogo é definição embutida e sem IO. Screen consome suas definições; apresentação consulta metadados para expandir formas. Perfil EFL no adaptador resolve suporte efetivo, sem importar o adaptador no domínio.

Não implementa widget, drag/drop, perfis EFL externos, plugins executáveis ou templates por projeto.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-003, FR-005, FR-015, FR-022; NFR-002 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Nenhuma persistência ou integração foi executada nesta entrega.
