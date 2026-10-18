# Specification Quality Checklist: Gerador de Jornada

**Purpose**: validar qualidade e consistência dos artefatos GitHub Spec Kit, sem confundir revisão documental com homologação do produto.
**Created**: 2026-10-17
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] Especificação centrada no usuário e nas saídas esperadas.
- [x] Arquitetura, pacotes e algoritmos separados em plan.md, research.md e contratos.
- [x] Restrições explícitas do usuário (Windows, Go, Java 1.6/EFL, PPT editável) preservadas como restrições do produto.
- [x] Seções obrigatórias do template preenchidas, sem placeholders de scaffold.
- [x] Domínios possuem features Spec Kit próprias e responsabilidades identificadas.
- [x] Sem arquivos avulsos de README, ADR, análise ou sprints; conteúdo incorporado nos artefatos previstos.

## Requirement Completeness

- [x] Requisitos identificados e verificáveis, com cenários de aceitação.
- [x] Histórias priorizadas e testes independentes definidos.
- [x] Limites de execução do editor e da aplicação Java explicitados.
- [x] PPT exploratório separado de demonstração; diferença entre entrada real e simulada documentada.
- [x] Formatos, referências, rascunhos inválidos, falhas e cancelamento tratados.
- [x] Resultados mensuráveis definidos sem afirmar medição já realizada.
- [x] Sem marcadores de esclarecimento pendente; dependências externas registradas como gates de implementação.

## Feature Readiness

- [x] Planos, modelo, pesquisa e contratos cobrem as histórias.
- [x] Backlog no formato checkbox, IDs, histórias e caminhos.
- [x] Sprints/dependências/matriz de rastreabilidade em tasks.md.
- [x] Próximas implementações dos cinco domínios possuem backlog próprio, sem serviços/módulos Go separados.
- [x] Constituição concreta, versionada e com regras de documentação.
- [x] Três skills copiadas sem alteração, origem/hashes/recursos ausentes registrados em research.md.
- [x] Feature coordenadora registrada para resolução pelos scripts do Spec Kit.

## Implementation Evidence

Estes itens permanecem abertos por se referirem ao produto futuro:

- [ ] Prova e homologação de janela desktop/DPI/EXE offline.
- [ ] Homologação de PPT exploratório/animado no PowerPoint 2016 e Microsoft 365.
- [ ] Disponibilização e validação real de JARs, APIs e bootstrap EFL.
- [ ] Compilação/execução Java SE 6/EFL e classfile 50.0 comprovados.
- [ ] Suítes de domínio, contrato e integração aprovadas.
- [ ] Metas de desempenho e cenário completo aprovados.

## Notes

Esta lista revisa os documentos produzidos; caixas marcadas não representam funcionalidades implementadas. As dependências EFL/PowerPoint pertencem à implementação/homologação futura. Não há exceção de arquitetura à constituição.

A referência a tecnologias em spec.md registra escolhas explícitas do usuário e restrições de saída; implementação interna permanece nos documentos técnicos. A validação documental automatizada é de estrutura, links, IDs e resolução de artefatos, não de conteúdo semântico de um sistema ainda inexistente.

Verificação executada em 2026-10-17: seis features resolvidas por `check_prerequisites.py --require-spec --require-tasks --include-tasks`; 38 arquivos Markdown em specs; 116 tarefas desmarcadas (76 de coordenação/integração e oito por domínio); IDs sequenciais por feature, links locais existentes, fences e exemplo JSON válidos; hashes das três skills iguais às origens. Nenhum teste da aplicação foi executado.
