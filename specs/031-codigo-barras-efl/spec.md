# Feature Specification: Código de barras — entrada validada

**Feature Branch**: `031-codigo-barras-efl` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-23
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Componente de oito campos baseado nas classes Java fornecidas na conversa, com validação por bloco, captura integral por leitora e mensagem de erro. Campos inicialmente em branco. Decomposição e implementação futura obrigatoriamente por domínio, conforme GitHub Spec Kit.
**Domínio responsável**: [003-catalogo](../003-catalogo/spec.md#codigos-barras-2026-10-23). Instâncias em [004-tela](../004-tela/spec.md#codigos-barras-2026-10-23); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md#codigos-barras-2026-10-23).
**Referência**: `CodigoBarrasNovaPlataformaCobranca` e `AbstractPnlNovaPlataformaCobranca`, transcritas pelo usuário; não são arquivos presentes no workspace. Nenhuma imagem chegou nos anexos. [DESIGN.md](DESIGN.md) registra a geometria do código, não fidelidade a uma imagem não recebida.

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Entrada composta 610 × 60 com oito campos derivados, validação matemática local, captura/colagem 44/47, comandos explícitos e mensagem modal. Oito valores e seus erros/foco ficam somente na sessão de prévia; PPT é estático/editável.

**Pendências e limites:** Não há homologação bancária/SISPL, referência visual independente nem matriz completa de teclado/TAB/habilitação. A classificação CAIXA/SIGCB é premissa local e não consulta ou autorização de pagamento.

**Fontes verificadas:** [barcode.go](../../apps/gerador-jornada/internal/catalog/barcode.go), [preview_compounds.go](../../apps/gerador-jornada/internal/application/preview_compounds.go), [barcodes_test.go](../../apps/gerador-jornada/internal/presentation/barcodes_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Digitar e corrigir o código por blocos (Priority: P1)

Como usuário da prévia, quero digitar uma linha digitável e corrigir o bloco inválido antes de continuar, sem aceitar um código com dígitos verificadores incorretos.

**Why this priority**: Evita propagar uma entrada inválida e reproduz o fluxo principal solicitado.
**Independent Test**: Em uma tela com apenas este componente habilitado, usar a fixture F01 do contrato compartilhado, alterar seus DVs e observar avanço, bloqueio, recuperação e conclusão, sem leitora física.

**Acceptance Scenarios**:

1. **Given** uma nova instância, **When** a prévia inicia, **Then** os oito campos estão vazios e somente o primeiro recebe digitação; os comprimentos são 5, 5, 5, 6, 5, 6, 1 e 14, nessa ordem.
2. **Given** o primeiro, terceiro ou quinto campo visual completo, **When** seu limite é atingido, **Then** o próximo campo visual é habilitado; ainda não há DV suficiente para validar o bloco correspondente.
3. **Given** um par de campos formando um dos três primeiros blocos, **When** o segundo campo do par completa, **Then** o DV Módulo 10 é validado antes do avanço; erro impede ir ao bloco seguinte, identifica o bloco e devolve foco ao primeiro campo do par, conservando os demais blocos já válidos.
4. **Given** os três primeiros blocos válidos, **When** o campo de um dígito e o último campo são informados, **Then** o primeiro é tratado como DV geral e o último como fator/valor; a conclusão exige também DV geral Módulo 11 correto, não somente os três DVs locais.
5. **Given** o último campo ativo e incompleto, **When** o usuário solicita concluir a entrada, **Then** ele é completado com zeros à direita até 14 posições e o código inteiro é validado; se falhar, a conclusão permanece bloqueada e a correção recomeça no primeiro campo, sem apagar automaticamente o rascunho digitado.
6. **Given** uma entrada válida, **When** a conclusão é liberada, **Then** os oito campos ficam somente para leitura e o usuário pode prosseguir com as ações explícitas já existentes na tela, sem navegação automática causada pelo componente.
7. **Given** uma entrada em andamento ou concluída, **When** o usuário solicita reiniciar a digitação, **Then** todo o estado desta entrada é limpo e o foco volta ao primeiro campo; outras instâncias e o documento não mudam.

---

### User Story 2 - Receber leitura integral e tratar código inválido (Priority: P1)

Como usuário da prévia, quero receber o código integral de uma leitora ou colagem e obter uma confirmação clara quando ele for inválido.

**Why this priority**: Uma captura integral não pode ser truncada no primeiro campo nem aceita parcialmente.
**Independent Test**: Submeter strings integrais de 44 e 47 dígitos, válidas e adulteradas, incluindo captura iniciada durante uma digitação parcial. Colagem permite exercitar o contrato sem periférico.

**Acceptance Scenarios**:

1. **Given** captura integral de 44 dígitos por leitura ou colagem, **When** ela é submetida, **Then** o DV geral é validado antes de converter e preencher os oito campos; uma captura válida libera a conclusão e fica somente para leitura. A colagem de F01 lido deve produzir exatamente os mesmos oito valores que a leitora.
2. **Given** uma linha digitável completa de 47 dígitos colada, **When** a colagem é submetida, **Then** são validados os três DVs locais e o geral antes de publicar o preenchimento; nenhum campo recebe um trecho de uma entrada recusada.
3. **Given** uma captura integral inválida, incompleta, excedente ou não numérica, **When** ela é submetida, **Then** a conclusão é bloqueada, os oito campos e resultados derivados desta instância são limpos e aparece exatamente “Código de barras inválido”, com botão “OK, entendi”.
4. **Given** o diálogo de erro aberto, **When** o usuário aciona “OK, entendi” ou Enter, **Then** o diálogo fecha e o foco volta ao primeiro campo vazio; nenhuma tecla acumulada na captura anterior é reaplicada.
5. **Given** captura por dispositivo configurado como teclado, **When** o usuário ativa o modo de leitura e encerra com Enter, **Then** toda a string é tratada como uma única entrada de 44 dígitos; não se infere a origem pela velocidade das teclas.

---

### User Story 3 - Consultar os dados derivados sem simular processamento bancário (Priority: P2)

Como autor da jornada, quero conferir banco, valor e fator presentes no código, sem confundir validação matemática com processamento real de um boleto.

**Why this priority**: Permite uma demonstração coerente sem inventar dados que não existem no código.
**Independent Test**: Consultar F01 e F02 do contrato compartilhado e entradas vazias/inválidas, sem UI bancária, rede ou consulta externa.

**Acceptance Scenarios**:

1. **Given** F01 aceita, **When** consultar os derivados, **Then** banco é `001`, valor é R$ 1,00, fator é `3737` e documento CAIXA é falso; zeros à esquerda permanecem nas representações textuais.
2. **Given** entrada vazia, parcial ou inválida, **When** consultar os derivados, **Then** eles estão indisponíveis, sem exceção por substring nem reaproveitamento de resultado de uma entrada anterior.
3. **Given** um código válido do banco `104`, **When** verificar documento CAIXA, **Then** a identificação não depende de consulta externa; a classificação SIGCB usa as regras locais e mantém indicação de premissa não homologada até comparação com `Modulo11Sispl` real.
4. **Given** fator não nulo, **When** consultá-lo, **Then** o fator numérico é exposto, sem converter silenciosamente em uma data de vencimento que depende do ciclo do fator.

---

### User Story 4 - Compor, preservar e apresentar a entrada (Priority: P2)

Como autor, quero inserir um componente composto, personalizar sua aparência e salvar a composição sem gravar o que foi digitado durante a prévia.

**Why this priority**: Mantém o mesmo modelo entre editor, prévia e exportação.
**Independent Test**: Inserir, personalizar, duplicar, desfazer/refazer, salvar/reabrir e exportar uma tela com duas instâncias; conferir IDs, configuração e isolamento do runtime.

**Acceptance Scenarios**:

1. **Given** uma composição válida, **When** inserir ou duplicar o componente, **Then** ele é uma instância de Tela, não oito componentes soltos; mantém ID de arquivo e propriedades comuns, com nova identidade na cópia.
2. **Given** uma personalização confirmada, **When** salvar/reabrir e exportar, **Then** rótulo, geometria e estilo permanecem; valores digitados, campos habilitados, mensagem e resultados derivados da prévia não entram no projeto.
3. **Given** um PPT exportado, **When** inspecionar os campos, **Then** textos e formas são editáveis, mas não há digitação, leitora, cálculo de DV ou modal executável no PowerPoint; essa limitação é informada antes da exportação.

### Edge Cases

- Blocos visuais não são os blocos de validação: os seis primeiros campos formam três pares de 10, 11 e 11 dígitos. Os campos de 1 e 14 posições completam os 47.
- Aceitar somente dígitos ASCII no conteúdo canônico. Na colagem integral, remover apenas espaços e pontos de apresentação; sinais, letras, quebras internas e símbolos não são descartados para transformar texto inválido em válido. Enter final do dispositivo é terminador, não conteúdo.
- Recusar 43, 45, 46, 48 ou mais posições como captura integral incompatível; não truncar. Código de arrecadação de 48 posições e payload PIX/QR textual não são boleto de cobrança de 44/47 dígitos.
- Campo visual incompleto não permite TAB/clique para campos posteriores. Correção permite apagar e redigitar o par recusado; campos posteriores ficam bloqueados até revalidação.
- Instância ou contêiner desabilitado não recebe digitação, colagem ou leitura. Teclado numérico vinculado usa o subcampo ativo, com os mesmos limites e validações, sem um segundo algoritmo na UI.
- Captura inválida limpa apenas o runtime da instância destinatária; cancelar uma captura antes de submetê-la conserva o rascunho anterior. Nenhuma entrada aparece em logs.
- Campo final vazio na ação explícita de concluir equivale a 14 zeros antes da validação geral, não a aprovação automática. Conclusão e reinício são operações locais, não novas ações de Jornada.
- Fator `0000` significa ausente; os demais são expostos como números, inclusive os que começam em zero. Data civil/ciclo de vencimento fica fora do escopo.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; foco legível, valores completos e mensagem não dependente somente de cor.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Oferecer a entrada como componente composto do Catálogo, com oito campos numéricos nos comprimentos e ordem da US1, inicialmente em branco, sem exemplo preenchido ou placeholder numérico.
- **FR-002**: Controlar o avanço sequencial e validar Módulo 10 ao completar cada um dos três pares; bloco inválido impede avanço, identifica a causa e retorna ao início do par para correção.
- **FR-003**: Exigir 47 dígitos e DV geral Módulo 11 correto para concluir a entrada digitada; campo de um dígito não tem DV próprio e campo de 14 posições não recebe Módulo 10.
- **FR-004**: Permitir a ação explícita de concluir o último campo incompleto preenchendo zeros à direita; executar a mesma validação geral, sem mascarar falha nem liberar conclusão antes dela.
- **FR-005**: Receber leitura de 44 dígitos e colagem integral de 44 ou 47; validar a representação apropriada e converter 44 ↔ 47 sem perder zeros ou aceitar 48 posições como boleto deste componente.
- **FR-006**: Recusar uma captura integral de forma atômica, limpar seu runtime e exibir o diálogo “Código de barras inválido” / “OK, entendi”; após confirmação devolver foco ao primeiro campo. Erro manual de bloco conserva os demais blocos válidos e informa o bloco localmente.
- **FR-007**: Oferecer reinício e modo de leitura explícitos, respeitar habilitação e aplicar o mesmo comportamento ao teclado físico e numérico vinculado; não depender de câmera, driver ou heurística de velocidade.
- **FR-008**: Expor somente para entrada aceita as representações lida/digitada, banco, valor monetário exato, fator de vencimento, identificação CAIXA e classificação SIGCB; separar classificação experimental de homologação e não tratar DV válido como boleto pagável.
- **FR-009**: Manter foco, conteúdo capturado, cache, erros, captura em andamento e conclusão somente na sessão da prévia, reiniciando em cada visita, fechamento ou reinício, sem alterar revisão, histórico ou configuração do documento.
- **FR-010**: Permitir composição e personalização conforme [DESIGN.md](DESIGN.md), com cópia profunda, ID comum, histórico e persistência a cargo dos domínios respectivos, sem gravar subcampos artificiais como componentes independentes.
- **FR-011**: Compartilhar a representação visual entre editor, prévia e PPT editável, explicando a ausência de validação/interação real no PPT; capacidade de roteiro não implementada deve produzir diagnóstico, nunca ser simulada silenciosamente.
- **FR-012**: Restringir o componente à validação matemática local de boleto de cobrança; não consultar boleto, CPF/CNPJ, beneficiário, banco ou vencimento real, não executar pagamento e não preencher automaticamente o [componente processado](../032-codigo-barras-processado-efl/spec.md).

### Key Entities

- **Definição de entrada**: tipo semântico, rótulo e aparência, pertencentes ao Catálogo; geometria/UUID/FileID pertencem à instância de Tela.
- **Boleto validável**: cinco blocos lógicos, representações canônicas de 44/47 dígitos e resultado de validação matemática local, sem IO ou UI.
- **Sessão de entrada**: oito valores visuais, campo ativo, erro, modo de captura e resultado aceito; estado transitório coordenado pela aplicação.
- **Dados derivados**: banco, valor, fator e identificadores locais, calculados sem constituir cadastro bancário ou entidade paralela de Projeto.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% dos pares com DV local incorreto da matriz de referência impedem avanço ao bloco seguinte; 100% dos códigos com DV geral incorreto impedem a conclusão.
- **SC-002**: F01 produz exatamente as mesmas 44/47 posições no percurso de ida e volta, banco `001`, valor R$ 1,00 e fator `3737`; nenhuma posição zero é perdida.
- **SC-003**: Toda captura integral recusada apresenta o texto e a única ação especificados, deixa os oito campos vazios e restaura o primeiro foco após a confirmação.
- **SC-004**: Salvar/reabrir e duplicar preservam 100% da configuração e as identidades devidas; nenhuma alteração runtime de duas instâncias modifica a outra ou o documento.
- **SC-005**: Editor, prévia e PPT preservam geometria com diferença máxima de uma unidade lógica; todos os campos/textos do PPT permanecem editáveis, sem prometer cálculo durante apresentação.

## Assumptions

- Os oito campos, rótulo e posições foram extraídos do código fornecido, não de imagem. Aprovação visual depende do anexo ainda ausente.
- As classes `Modulo10`, `Modulo11Bloqueto` e `Modulo11Sispl` não foram recebidas. O [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md#contrato-codigos-barras-2026-10-23) define premissas matemáticas verificáveis; equivalência SISPL e classificação SIGCB não estão homologadas. A comparação com a classe real é condição para anunciar essa equivalência, não para escrever a spec.
- “Pistola QR” designa aqui um dispositivo que entrega os dígitos de boleto como teclado. Não inclui leitura por câmera, SDK de hardware ou interpretação de PIX.
- Ações locais de concluir/reiniciar/ler podem ser expostas nos controles de interação da prévia; não são botões de navegação adicionados automaticamente à tela. Não introduzir condicionais no grafo de Jornada nesta entrega.
- Código matematicamente válido não comprova existência, titularidade, valor devido, vigência ou autorização de pagamento. Fator numérico não resolve o ciclo de vencimento.
- Requisitos são locais: usar `031-codigo-barras-efl/FR-xxx`. Modelo, tarefas e implementação permanecem nos domínios existentes, conforme a constituição; este detalhamento não cria um domínio de código de barras.
- Esta etapa entrega documentação. Os backlogs de domínio continuam abertos; não há novo EXE, teste funcional ou homologação. Java/EFL continua retirado conforme [001/FR-059](../001-gerador-jornada/spec.md#dependencias-ausentes-fr-059).

## Divisão por domínio e rastreabilidade

| Responsabilidade | Artefato e trabalho responsável |
| --- | --- |
| Tipos, defaults, validações e derivados puros | [003/FR-013–FR-015 e T018–T020](../003-catalogo/tasks.md#tarefas-codigos-barras-2026-10-23); FR-001–FR-008 desta spec |
| Instâncias, identidade, cópia e comandos de composição | [004/FR-019 e T015–T016](../004-tela/tasks.md#tarefas-codigos-barras-2026-10-23); FR-010 |
| Serialização, compatibilidade e round-trip | [002/FR-011 e T013–T014](../002-projeto/tasks.md#tarefas-codigos-barras-2026-10-23); FR-009–FR-010 |
| Rotas explícitas, sem navegação por DV | [005/FR-008 e T011](../005-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); FR-007/FR-012 |
| Quadros e capacidades de roteiro, sem execução no PPT | [006/FR-011 e T011](../006-demonstracao/tasks.md#tarefas-codigos-barras-2026-10-23); FR-011 |
| Interface, prévia, cena e exportadores | [001/FR-060 e T135–T139](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); integração, FR-009/FR-011 e evidência |

Os planos, modelos e tarefas canônicos desses domínios são a próxima entrada do fluxo `plan → tasks → analyze → implement`. Não manter backlog concorrente nesta pasta. [Checklist](checklists/requirements.md) avalia a qualidade da spec, não implementação nem aprovação humana.

## ID das instâncias

Aplicar [004/FR-011–FR-012](../004-tela/spec.md#identidade-de-arquivo--2026-10-18). UUID e FileID do composto não se confundem com índices internos de campos, valores digitados ou nomes visuais. Nenhuma alteração de rótulo renomeia identidades.