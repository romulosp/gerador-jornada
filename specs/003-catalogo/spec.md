# Feature Specification: Domínio Catálogo

**Feature Branch**: `003-catalogo` (identificador; sem branch Git criada)
**Created**: 2026-10-17
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Decomposição por domínio da [feature coordenadora](../001-gerador-jornada/spec.md).

**Instrução vigente — 2026-10-21:** a [clarificação FR-012](#builtin-filtrado-fr-012) delimita a disponibilidade operacional do builtin após a mudança de escopo. Inventário original e perfis/documentos existentes continuam preservados; verificação ainda aberta, sem alterar o checklist 13/2. Detalhes e evidências atuais na [seção 14 da coordenadora](../001-gerador-jornada/quickstart.md#14-builtin-filtrado).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Catálogo 3 oferece dez Kinds: sete famílias base e barcodeInput/barcodeProcessed/applicationMenu. A biblioteca visual contém 205 chaves tipo/nome, 224 recursos referenciados, seis fontes físicas e oito definições retiradas por recursos ausentes. Skin não cria um novo Kind.

**Pendências e limites:** O metamodelo declarativo de capacidades, as nove famílias adicionais da expansão e as variantes completas de Report/twoLabel/alphaKeyboard não foram implementados. Quantidade de skins não comprova quantidade de componentes nem homologação.

**Fontes verificadas:** [model.go](../../apps/gerador-jornada/internal/catalog/model.go), [model.go](../../apps/gerador-jornada/internal/skin/model.go), [compound_test.go](../../apps/gerador-jornada/internal/catalog/compound_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Escolher um componente conhecido (Priority: P1)

Definir as sete famílias, seus defaults, propriedades tipadas e capacidades sem executar código do projeto.

**Why this priority**: Uma mesma definição deve orientar paleta, inspector e validação, evitando interpretações diferentes entre editor e saídas.
**Independent Test**: Executar os casos do domínio com fixtures e sem UI/exportadores.

**Acceptance Scenarios**:

1. **Given** o catálogo versão 1, **When** a paleta solicita componentes disponíveis, **Then** recebe exatamente botão, rótulo, título, campo, seleção, relatório e teclado numérico, com identificadores estáveis.

---

### User Story 2 - Editar apenas propriedades válidas (Priority: P2)

Como autor da jornada, quero receber comportamento previsível e diagnóstico quando uma alteração viola o contrato deste domínio.

**Why this priority**: Evita que operações inválidas sejam propagadas para persistência ou saídas.
**Independent Test**: Construir o estado de referência inválido e verificar diagnóstico/estado preservado, sem executar outro domínio.

**Acceptance Scenarios**:

1. **Given** um campo de texto, **When** se tenta atribuir uma propriedade de outra família ou enum desconhecido, **Then** a operação é recusada sem alterar o componente; propriedades suportadas ficam preservadas.

### Edge Cases

Cobrir sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível. Alterações inválidas não publicam estado parcialmente aplicado. Rascunho legível pode ter diagnóstico, mas isso não equivale a ser executável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Disponibilizar exatamente os kinds button, label, title, textField, radio, report e numericKeyboard no catálogo inicial.
- **FR-002**: Definir defaults, tipos, enums e estilos conforme o contrato de catálogo; versão inicial igual a 1.
- **FR-003**: Validar propriedades por structs de família, sem map[string]any atravessando os casos de uso.
- **FR-004**: Expor definições para paleta/inspector sem depender de Fyne, JSON, XML ou templates Java.
- **FR-005**: Tratar capacidade Java como mapeamento de perfil do adaptador; nome EFL observado não comprova API nem permite ignorar propriedade alterada.

### Key Entities

ComponentDefinition, Kind, PropertyDefinition e structs ButtonProperties, LabelProperties, TitleProperties, TextFieldProperties, RadioProperties, ReportProperties e NumericKeyboardProperties. Modelo e invariantes em [data-model.md](data-model.md).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: As sete famílias têm defaults válidos e round-trip sem perda; 100% dos enums/propriedades fora do contrato são diagnosticados nos testes de referência.
- **SC-002**: Todos os estados inválidos do conjunto de testes produzem diagnóstico verificável sem alteração parcial do estado anterior.

## Assumptions

- Domínio do mesmo monólito Go; não é microserviço ou aplicação independente.
- Requisitos pai cobertos: FR-003, FR-005, FR-015, FR-022; NFR-002. IDs FR/US/SC são locais a esta feature.
- Contratos compartilhados: [component-catalog.md](../001-gerador-jornada/contracts/component-catalog.md).
- Não implementa widget, drag/drop, perfis EFL externos, plugins executáveis ou templates por projeto.
- Esta especificação define o alvo; evidências de implementação e testes estão no quickstart da feature coordenadora, sem homologação implícita.

## Especificações individuais dos componentes

O detalhamento abaixo mantém este domínio como responsável pelo modelo e pelas regras compartilhadas. Cada feature tem contrato visual próprio; sua existência documental não conclui implementação ou homologação.

- [Botão — EFLButton](../007-botao-efl/spec.md) — [DESIGN.md](../007-botao-efl/DESIGN.md).
- [Rótulo — EFLLabel](../008-rotulo-efl/spec.md) — [DESIGN.md](../008-rotulo-efl/DESIGN.md).
- [Título — EFLLabelTitulo1](../009-titulo-efl/spec.md) — [DESIGN.md](../009-titulo-efl/DESIGN.md).
- [Campo de texto — EFLTextField](../010-campo-texto-efl/spec.md) — [DESIGN.md](../010-campo-texto-efl/DESIGN.md).
- [Seleção — EFLRadio](../011-selecao-efl/spec.md) — [DESIGN.md](../011-selecao-efl/DESIGN.md).
- [Relatório/lista — EFLReportEvent](../012-relatorio-efl/spec.md) — [DESIGN.md](../012-relatorio-efl/DESIGN.md).
- [Teclado numérico — EFLVirtualKeyboardNumeric](../013-teclado-numerico-efl/spec.md) — [DESIGN.md](../013-teclado-numerico-efl/DESIGN.md).
- [Paleta — ComponentTab](../023-paleta-componentes/spec.md) — [DESIGN.md](../023-paleta-componentes/DESIGN.md).

## Atributo compartilhado das instâncias

- **FR-006**: Todas as sete famílias e componentes criados por conjuntos prontos expõem ID de arquivo como atributo comum da instância no domínio Tela. Não duplicar essa regra nas variantes Properties do catálogo. Componentes não exigem Pnl; seus IDs seguem as regras de 004-tela/FR-011–FR-012, enquanto os tipos e versões do catálogo continuam estáveis.

## Expansão pós-v1 autorizada — 2026-10-19

FR-001/FR-002 e US1 descrevem o catálogo inicial de sete famílias/versão 1, não um limite permanente. O usuário aprovou todo o acervo visual Java/EFL; catálogo alvo 2, com preservação dos sete discriminadores existentes. Não alterar checklist (13 marcados/2 abertos, prosseguimento autorizado).

- **FR-007**: Classificar todas as classes, helpers, máscaras e cada tipo/nome/estado de skin do acervo na matriz coordenadora, sem descartar referências incompletas; skins viram variantes/presets da família semântica, não Kind por entrada XML.
- **FR-008**: Acrescentar tipos ToggleButton, TextArea, PasswordField, ComboBox (UF como preset), Table, AlphaKeyboard, ProgressBar, Panel e Dialog, com propriedades/defaults/validação/cópia/capacidades tipadas; duas legendas são variante Button e Report/Report2/Report3/Event/Event2/Notificacao/SelecaoMultipla/Internal são variantes Report. Conservar Radio simples e presets compostos/Frame.
- **FR-009**: Definir skin/version, estados suportados e aparência/overrides com presença explícita, fonte/face real, RGBA e métricas/espaçamento comprovados; biblioteca gerada é dado puro estático, não importação XML/runtime nem dependência Fyne/PPT no domínio. Declaração XML bold/italic não equivale à pintura Java.
- **FR-010**: Expor disponibilidade/capacidades por família/variante/propriedade/estado/alvo, com referência de evidência e bloqueio; dependências/delegates ausentes impedem fidelidade total, não desaparecem do catálogo. Nenhuma liberação apenas por adicionar enum ou existir controle semelhante.
- **FR-011**: Preservar IDs de opções/linhas, formatos/defaults semânticos anteriores e overrides na migração; button default atual 110 × 60 permanece, teclado novo 200 × 350 substitui default inicial 200 × 280 somente em novas inserções. Máscaras e comandos locais seguem Java comprovado, sem negócio/periféricos e sem reativar Java.

**Aceitação**: cada nova família/variante tem default válido e contrato completo de criação/Validate/deep-copy/duplicação/remoção/Undo/codec/inspector/cena/preview/export; teste de referência cobre cada classe/skin/estado, não só as sete atuais. Password/Dialog ou imagens/delegates faltantes permanecem bloqueados para homologação. Galeria/paleta mostra essa diferença. Matriz [research.md](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19); contrato [component-catalog.md](../001-gerador-jornada/contracts/component-catalog.md#catálogo-2--expansão-autorizada-2026-10-19); T009–T016 locais e 001/T108–T125.

<a id="builtin-filtrado-fr-012"></a>
## Disponibilidade do builtin filtrado — clarificação de FR-007/FR-010 (2026-10-21)

- **FR-012**: Conforme a decisão posterior do usuário de excluir do builtin o que não for encontrado e reimportar em futura change, oferecer para novos projetos somente as definições ativas do builtin filtrado por [001/FR-058](../001-gerador-jornada/spec.md#builtin-filtrado-fr-058). A retirada afeta a **definição inteira**, não só um estado: qualquer referência `imagePath`, `slider` ou `font:type` ausente/inválida, caminho recusado ou tipo incompatível adia a definição e, havendo filho afetado, toda a composição pai/filhos. Consumir identidades tipadas e vínculos reais, sem excluir homônimos de outro tipo ou composições independentes. Não repetir o parser no Catálogo nem criar Kind por entrada adiada.

Esta clarificação prevalece sobre a oferta anterior de variantes incompletas como bloqueadas **somente para o builtin operacional**. As oito variantes listadas em [FR-058](../001-gerador-jornada/spec.md#builtin-filtrado-fr-058) ficam fora de novas inserções/trocas de variante nessa revisão, mas permanecem auditáveis em `deferred_definitions` (`type`, `name`, `definition`, `reasons`), com `source_tree`/`xml_sha256` e acervo original intactos. **206 definições ativas + 8 adiadas = 214 originais**, com **8 fontes / 270 recursos / 31 diagnósticos ativos / 16 referências ausentes adiadas**. Definições incluem partes/variantes: não são 206 famílias nem 206 opções em cada seletor; o catálogo continua com 16 famílias.

**Correção de estado — 2026-10-24 / 1.0.0.1:** os números 206+8/270/31/16 famílias acima pertencem ao perfil amplo de 16 famílias que **não foi entregue**. A biblioteca operacional efetiva desta revisão, para as sete famílias base e os três compostos, é **205 chaves tipo/nome, 224 recursos, 6 fontes físicas e 8 definições indisponíveis fora da oferta**, geradas por `scripts/import_legacy_skin.py` em `library.json`. Não usar os números antigos como inventário do executável atual.

Skin **`legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a`**; manifest/Profile/project/catalog/produto **1/1/3/2/1.0.0.0** inalterados. Não reduzir metadados por uso, apagar fontes/bytes disponíveis, migrar/podar/enriquecer perfis próprios ou arquivados nem alterar documentos existentes para aplicar esse filtro. A resolução exata continua sem fallback ou upgrade silencioso; capturar/salvar/reabrir deve conservar o perfil original completo e permitir inserções válidas nele.

**Aceitação local pendente:** novo projeto oferece variantes ativas por família e nenhuma das oito adiadas; falha de filho não deixa órfãos e não atinge identidade independente homônima. Projeto antigo/perfil próprio/arquivado conserva variantes, estados e referências no round-trip. Reimportação exige change futura com recursos exatos, origem/hash/autorização e regeneração/validação, sem substituir recursos por nome parecido. Licenças/fontes, `DELEGATE_MISSING`, `FIDELITY_UNVERIFIED` e gates Java/PPT permanecem; a redução não prova 100% de fidelidade ou aceite integral.

**Decisão posterior prevalente — dependências ausentes:** [001/FR-059](../001-gerador-jornada/spec.md#dependencias-ausentes-fr-059) retira as referências/call sites de delegates externos ausentes e a ação/integração Java/EFL, em vez de manter avisos ou botão desabilitado. Retira também a oferta/helper/ações da calculadora sem layout. Os valores antigos permanecem legíveis e não sofrem conversão silenciosa; os layouts oferecidos são Padrão/numeric/alpha. O builtin 206+8, recursos disponíveis e perfis/documentos anteriores não são podados por essa mudança. Licenças, escala, métricas, portabilidade e limitações reais de PPT permanecem. Esta decisão supera somente as exigências históricas de preservar delegates/Java como dependências do produto. Evidência vigente: [001/T133](../001-gerador-jornada/tasks.md#t133-dependencias-ausentes)/[seção 15](../001-gerador-jornada/quickstart.md#15-dependencias-ausentes-removidas), sem herdar PASS ou hashes da entrega filtrada abaixo.

Verificação local em [T017 aberta](tasks.md#t017-builtin-filtrado), integração única em [001/T132 aberta](../001-gerador-jornada/tasks.md#t132-builtin-filtrado). **19 Python do importador PASS**, incluindo descoberta/execução dos quatro novos métodos; **4 Python de empacotamento e importador `--check` PASS**. Regressões Go de biblioteca/seletores/Capture/perfis próprios e arquivados implementadas; gofmt/vet e **compilação de 15 suítes PASS, ZERO execuções Go**, respeitando a política Windows. Build/pacote e integridade somente leitura PASS; dois EXEs idênticos, identidades reais e limites na [seção 14](../001-gerador-jornada/quickstart.md#14-builtin-filtrado). Execução Go, roteiros de UI/round-trip, DPI/Office e aceite pendentes; registros 214/47 e hashes anteriores são históricos, sem alterar caixas antigas.

<a id="codigos-barras-2026-10-23"></a>
## Incremento planejado — códigos de barras (2026-10-23)

Detalhes de uso nos componentes [031 — entrada validada](../031-codigo-barras-efl/spec.md) e [032 — processado parametrizável](../032-codigo-barras-processado-efl/spec.md), sem criar domínio financeiro ou duplicar modelo/tarefas nessas pastas. Esta seção é planejamento, não ampliação já implementada das 16 famílias.

- **FR-013**: Definir duas famílias semânticas tipadas, entrada de boleto e apresentação processada, com defaults, propriedades, capacidades e cópia/validação coerentes. Não criar um Kind por campo/bloco/skin nem oferecer enum sem suporte dos consumidores. Entrada começa vazia; processado começa com um bloco e linhas CNPJ/Nome de valores vazios.
- **FR-014**: Manter no núcleo de Catálogo as regras puras locais de validação por três blocos Módulo 10 e DV geral Módulo 11, conversão 44 ↔ 47 e derivados banco/valor/fator/CAIXA/SIGCB. Não importar UI/IO ou consultar serviços; distinguir SIGCB sob premissa de classe não recebida de equivalência SISPL comprovada. Contrato matemático e fixtures na [seção compartilhada](../001-gerador-jornada/contracts/component-catalog.md#contrato-codigos-barras-2026-10-23); 48 dígitos não são boleto suportado.
- **FR-015**: Representar blocos/linhas processados como dados ordenados com IDs e tipos de rótulo CPF/CNPJ/Nome/Nome fantasia/Razão social; cada valor é texto ilustrativo independente. Trocar tipo não apaga/converte valor; exigir ao menos um bloco e uma linha por bloco, IDs válidos/únicos e enum conhecido. Estado de digitação, foco, erros e popup não pertencem às propriedades persistidas.

**Aceitação documental do alvo:** 031/US1–US3 cobrem FR-013–FR-014; 032/US1–US2 cobrem FR-013/FR-015. Defaults/cópia/enum/estrutura e matemática têm fixtures sem UI. Derivação não fornece nome/CPF/CNPJ para o processado nem comprova processamento bancário. Catálogo 3 é alvo planejado para reconhecer esses tipos; catálogo 2 e perfis existentes não são alterados nesta entrega. Implementação em [T018–T020](tasks.md#tarefas-codigos-barras-2026-10-23), integrações em 001, persistência em 002 e instâncias em 004.

<a id="menu-aplicacao-2026-10-23"></a>
## Incremento autorizado — Menu aplicação (2026-10-23)

- **FR-016**: Definir Menu aplicação como família composta com modelos separados de menu/submenu e itens ordenados; defaults baseados na árvore Java/XML do usuário, propriedades editáveis de legenda/categoria/estilo/habilitação e submenu opcional ou destino explícito de folha. Aplicar as regras puras de [033/FR-002–FR-004/FR-008](../033-menu-aplicacao-efl/spec.md): layouts raiz/padrão/financeiro/Bolão, menus vazios válidos, até 16 níveis contando raiz e 1000 itens por árvore, UUIDs globalmente únicos e ausência de submenu+destino. Metadados de movimento/atalho são dados inertes, não autenticação, dispatch ou rota inferida.

**Aceitação**: defaults recebem gerador de IDs externo e itens habilitados; snapshot conserva identidade com cópia profunda. Duplicações e referências globais são verificadas por Tela/aplicação/Jornada. Catálogo contém só regras puras, sem Fyne/IO/Java. [T021](tasks.md#tarefas-menu-aplicacao-2026-10-23) acrescenta este modelo ao alvo catálogo 3 **já planejado para 031/032**, sem catálogo 4 e sem novo domínio financeiro. Implementação autorizada não significa executada ou homologada nesta entrega.
