# DESIGN — Biblioteca manual e fábrica de componentes — SPEC 038 v0.4

**Estado:** desenho funcional integrante do plano proposto, sem UI implementada ou aceite visual. Spec v0.4/emenda 1.4.0 aprovadas para planejamento em 2026-10-25; isso não aprova o software.
**Requisitos:** FR-001–FR-004, FR-018–FR-032; **aceitação:** US3/US4/US5, SC-003–SC-009/SC-012/SC-013.
**Base:** [DESIGN do editor](../001-gerador-jornada/DESIGN.md), [paleta 023](../023-paleta-componentes/DESIGN.md), [propriedades 028](../028-propriedades/DESIGN.md), [paridade visual](../../.agents/skills/jornada-visual-parity/SKILL.md).

## Responsabilidade e estilo

Este documento é dono dos novos fluxos da biblioteca: acesso, listagem, detalhes, criação, importação, revisão, publicação, versão e inserção. 023 é dona da integração com a paleta e 028 dos formulários/edição de instâncias. Não duplicar os fluxos ou criar biblioteca separada para GUI e MCP.

O autor abre o software e usa a biblioteca manual sem cliente MCP/LLM, inclusive offline. O agente externo usa as mesmas definições, versões e casos de uso. Estado de editor, seleção e filtros não pertencem à jornada.

Manter a linguagem visual/tokens e padrões de painéis já definidos no projeto. Os componentes representados usam aparência do legado: mesmas skins, fontes físicas, cores/alpha, métricas, estados e partes internas na criação manual/MCP, canvas, prévia e PPT aplicável. Não propor tema paralelo, ícones/imagens inventados ou miniatura genérica como prova de fidelidade. Geometria detalhada dos novos controles será refinada no plano, respeitando as referências da paleta e DPI existentes.

## Entrada, listagem e detalhes

- Acesso em seção identificada como **Biblioteca** da paleta existente; preservar famílias/presets atuais. Oferecer ações explícitas **Criar da seleção** e **Importar definição**, utilizáveis sem agente.
- Listar nome, descrição resumida, categoria (família nativa/composição/variante), versão, origem, estado, compatibilidade e disponibilidade. Proposta de novo Kind aparece como proposta pendente de engenharia, nunca opção inserível.
- Filtro por nome/categoria/estado não altera dados. Lista tem rolagem própria e não invade o canvas; detalhes preservam contexto da seleção. Janela menor deve permitir rolar/reorganizar os detalhes sem ocultar ações críticas.
- Detalhes mostram identidade, descrição, versão, propriedades/parâmetros, Kinds, recursos, dependências, alvos suportados, diagnósticos e operações permitidas. Origem informativa não é selo de autorização/homologação.
- Compatibilidade e estado de publicação são dimensões distintas: uma versão publicada pode estar incompatível no ambiente atual. Não usar apenas cor para explicar estado ou esconder o motivo de indisponibilidade.

## Estados observáveis

| Estado | Conteúdo e ações |
| --- | --- |
| Biblioteca vazia | “Nenhum componente reutilizável”; oferecer criar da seleção/importar; explicar seleção necessária sem conteúdo fictício |
| Publicados disponíveis | Nome/versão/categoria/compatibilidade visíveis; detalhes e Inserir habilitados conforme tela e parâmetros |
| Rascunho | Editar, validar e cancelar; não aparece como publicado nem altera versão anterior |
| Validado | Mostrar resultado e alvos; habilitar Revisar publicação, sem publicar automaticamente |
| Aguardando aprovação | Resumo imutável e ações Confirmar, Cancelar ou Corrigir; confirmar somente com permissões e revisão atuais |
| Publicado | Versão/conteúdo fixos; inserir/exportar/criar nova versão; não editar conteúdo in place |
| Retirado da oferta | Mostrar estado/motivo e detalhes auditáveis; impedir novas inserções, preservar instâncias existentes |
| Incompatível | Mostrar Kind/catálogo/versão requerido; Inserir indisponível, detalhes e recuperação permitidos |
| Recurso ausente | Localizar recurso e dependência afetada; não usar substituto; permitir corrigir rascunho/importação |
| Falha ao importar | Motivo contextual; cancelar/corrigir e tentar novamente, mantendo biblioteca anterior |
| Falha ao publicar | Mostrar motivo e estado real; manter rascunho, não anunciar sucesso ou repetir efeito sem reconciliação |
| Conflito de versão | Identidade/versão e divergência de conteúdo; recusar substituição e permitir cancelar ou preparar outra versão |

Carregamento, operação em andamento e falta de permissão devem ter feedback explícito, bloquear somente ações conflitantes e conservar foco/contexto. Mensagem de erro não substitui os dados anteriores por lista vazia.

## Criação e edição de rascunho

1. Autor seleciona componentes da tela e usa **Criar da seleção**; mostrar seleção, geometria relativa e vínculos internos/externos.
2. Autor informa nome/descrição/versão e parâmetros tipados suportados. Não permitir acrescentar propriedade arbitrária ou comportamento não instalado.
3. Informar referências externas que exigem parâmetro; se não puderem ser resolvidas, recusar captura completa sem salvar definição parcial.
4. Validar esquema, tipos, referências, recursos, portabilidade e limites. Mostrar diagnóstico junto ao campo/parte afetado, preservando entradas.
5. **Revisar publicação** abre o resumo descrito abaixo. Cancelar/corrigir conserva o rascunho conforme política; fechar não publica.

A primeira biblioteca trabalha com cópias independentes. Alterar parâmetro da definição em rascunho não edita automaticamente a seleção de origem ou instâncias antigas.

Para análise do legado, mostrar separadamente **Análise** e **Captura vinculada** (identidade, revisão e integridade), com ação **Vincular captura** para uma seleção já capturada da mesma sessão. Sem vínculo, permitir analisar/corrigir e indicar “Captura necessária para publicar”; novo Kind continua pendente de engenharia. Vincular usa atualização CAS da proposta, sem criar conteúdo a partir da descrição. Autor pode construir/capturar no editor sem LLM; agente faz a mesma construção pelos comandos compartilhados. Alterar análise/captura invalida revisão anterior. Publicar persiste captura/proveniência e conserva análise transitória: avisar que ela precisa de descarte explícito ao encerrar e não foi salva integralmente no pacote.

## Importação e exportação

Importar inicia por seleção de origem autorizada. Mostrar identidade, versão, categoria, recursos, integridade, dependências e diferenças/conflitos antes de qualquer escrita na biblioteca. A prévia é leitura/rascunho; confirmar importação persistente usa consentimento e revalidação. Falta de recursos, conteúdo inválido ou colisão bloqueia confirmação e mostra o motivo; não instalar fontes/skins/código automaticamente.

Definição publicada com conteúdo idêntico pode ser reconhecida como já existente; não criar duplicata silenciosa. Exportação usa destino autorizado e revisão/consentimento aplicável; dados portáteis não contêm caminhos pessoais ou arquivos não selecionados. Contrato físico do pacote será definido no plano.

## Revisão e publicação local

Antes de confirmar, mostrar nome, categoria, identidade/versão, mudanças, destino da biblioteca, dependências, recursos, integridade, resultado de validação, alvos/limitações e impactos em oferta. Esclarecer que publicação é local e não altera jornadas anteriores ou envia ao GitHub.

A ação Revisar publicação usa a preparação comum de aplicação correspondente a componente_revisar. Admissão de importação, revisão e confirmação são etapas distintas; prévia não publica nem contorna quotas. Exportação v1 informa criação exclusiva de .jcomp no workspace autorizado e conflito de homônimo. No MCP, as mesmas etapas usam DTOs fechados e HostApproval-v1; a GUI confirma por evento humano, sem depender de assinatura enviada pelo agente.

**Confirmar publicação**, **Cancelar** e **Corrigir** são ações distintas. Confirmar não pode ser acionado por uma edição incidental do formulário. Se revisão, conteúdo, dependências ou permissões mudarem, invalidar confirmação e pedir nova revisão. Falha preserva rascunho e versão anterior; resultado incerto exige reconciliação explícita. O fluxo manual confirma por evento humano da GUI; o MCP segue [consentimento sem GUI](contracts/consentimento.md), usando a mesma revisão/validação de aplicação e sua própria fronteira de confiança.

## Inserção, versões e propriedades

- Selecionar versão publicada compatível e tela; configurar parâmetros válidos e inserir por clique ou drop conforme 023. Sem tela ou espaço insuficiente, apresentar motivo, preservar seleção e aplicar as regras de geometria existentes.
- Uma inserção é operação reversível do conjunto inteiro; gera componentes com novas identidades, recursos necessários e vínculos internos corretos. Não deixar metade do conjunto após erro.
- Após inserir, os componentes são independentes e editáveis no inspector 028; não mostrar controle de atualização automática da definição inexistente. Novo conteúdo da biblioteca não muda instâncias anteriores.
- **Criar nova versão** parte de cópia em rascunho; conteúdo publicado não é editado. **Retirar da oferta** exige revisão/consentimento, mostra alcance e preserva dados/instâncias; não equivale a excluir fisicamente arquivos.
- Prévia da definição utiliza a mesma apresentação do legado. Mostrar qual estado/alvo está representado; miniatura não prova interação, homologação no PPT ou suporte a Java.

## Teclado, foco e mensagens

Tab/Shift+Tab percorrem filtros, lista, detalhes e ações em ordem lógica; Enter ativa apenas a ação focada e permitida. Escape cancela a etapa atual sem publicar; rascunho não salvo segue política visível. Depois de erro, foco vai ao diagnóstico/campo relevante; depois de confirmar/cancelar, retorna ao contexto de origem.

Foco perceptível, textos legíveis e rótulos de categoria/estado devem acompanhar ícones/cores. Ação indisponível tem motivo acessível nos detalhes. Não atribuir atalhos globais novos sem verificar conflitos do editor. Textos longos, Unicode e 100%/150%/200% DPI fazem parte da revisão futura.

Formulários aplicam os mesmos limites de pontos de código Unicode de MCP/biblioteca/propostas, com mensagem no campo; não cortar texto, contar bytes/UTF-16 como caracteres ou normalizar silenciosamente. Propriedade sugerida distingue valor desconhecido, observado com referência de evidência e proposto com motivo; não exibir sugestão da LLM como default instalado/homologado.

Rascunhos de proposta pertencem à sessão, separados de alterações da jornada. Mostrar os dois tipos de pendência ao encerrar; salvar jornada não salva proposta, publicar proposta não salva jornada. Oferecer voltar/corrigir ou descarte explícito; sem persistência oculta de novo Kind. Quota de edição esgotada desabilita novas edições com motivo, mantendo consultar/revisar/finalizar autorizado/descarte; negação/expiração preserva trabalho para nova revisão. O modo manual continua utilizável sem MCP/LLM, com a mesma coordenação e o mesmo legado visual.

## Matriz de verificação e responsabilidades futuras

| Cenário | Resultado a comprovar |
| --- | --- |
| Biblioteca vazia → criar/importar | Ações compreensíveis e dados anteriores preservados |
| Importação inválida/conflito | Nenhuma substituição, mensagem e correção contextual |
| Rascunho → validar → revisar → cancelar/confirmar | Sem publicação por edição incidental; apenas versão autorizada persiste |
| Incompatibilidade/recurso ausente | Sem inserção/fallback; motivo e recuperação visíveis |
| Inserção manual × MCP | Mesmo conteúdo, recursos, vínculos e aparência, descontadas novas identidades |
| Default legado e overrides | Mesmas fontes/alpha/partes em criação, canvas, prévia e PPT aplicável |
| Nova versão/retirada | Instâncias antigas preservadas; oferta atualizada sem apagar dados |
| Navegação por teclado/DPI | Foco, leitura, ações permitidas e mensagens utilizáveis sem depender de cor |

O [plano](plan.md) referencia este DESIGN e 023/028; tarefas futuras separarão UI, casos de uso por domínio, testes e homologação. O [quickstart](quickstart.md) contém comandos de referência e cenários futuros; revisão executável, harness efetivo, capturas e Office serão registrados somente quando existirem. Este desenho não é prova de implementação ou aprovação visual.
