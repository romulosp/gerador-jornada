# Data Model: Domínio Catálogo

## Extensão planejada — definições reutilizáveis 038

Definição, versão imutável, requisitos, parâmetros tipados e disponibilidade são conceitos deste domínio conforme [modelo 038](../038-integracao-mcp-agentes/data-model.md); wire/schema/arquivos ficam nos adaptadores. Oferta published/withdrawn é separada do conteúdo e da compatibilidade calculada. Classificações existentes de composição/variante/Kind são preservadas, com mesmo default legado nos dois modos. [Contrato biblioteca v1](../038-integracao-mcp-agentes/contracts/biblioteca.md) e [schema](../038-integracao-mcp-agentes/contracts/library.schema.json) detalham fronteira externa. Não há structs/capacidades novas implementadas nesta etapa; tarefas futuras qualificam domínio e integração sem duplicação.

## Situação corrente — revisão 1.0.0.1

Para novos dados de parâmetros da biblioteca 038, limites textuais contam pontos de código Unicode válidos: texto/default até 65.536, nomes e campos específicos conforme schema, sem normalizar/truncar. Ambos os adaptadores usam a mesma validação; bytes UTF-8 de frame/pacote são orçamento independente. A regra não implica truncar conteúdo de documento histórico na leitura. Descritores de proposta/evidência e suas revisões são coordenação de aplicação, não novas propriedades persistidas do Catálogo.

catalog.Definition contém Kind/Name/Width/Height. Properties é uma união tipada de dez Kinds; catalogVersion=4 não significa implementação das 16 famílias planejadas. As skins estáticas vêm de internal/skin e não são entidades persistidas de Catálogo.

Este quadro e a [spec](spec.md#estado-verificado-1-0-0-1) descrevem a implementação atual. As seções datadas de expansão conservam decisões/requisitos e registros históricos: tipos, arquivos, testes e perfis ali citados só são evidência atual quando existem na árvore e estão cobertos pela seção 22 do [quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Não atribuir ao build atual resultados de uma revisão antiga.


**Date**: 2026-10-17 | **Spec**: [spec.md](spec.md)
**Owner**: `internal/catalog`; mesmo modelo do monólito, sem entidades paralelas.

## Entidades e campos

No código atual, catalog.Definition mantém Kind, Name, Width e Height; Properties é a união tipada de famílias. O modelo declarativo completo de ComponentDefinition com propertyDefinitions/capabilities continua planejado em T004/T006, não é um tipo Go existente. Kind é enum estável. Propriedades específicas ficam em structs por família; Option/Row têm IDs e texto, sem separação por vírgulas.

## Invariantes

catalogVersion atual é 4 (versões 1, 2 e 3 são históricas). Enums devem pertencer à lista do contrato; fontSize positivo; Style usa cores #RRGGBB e não tem FillAlpha; alpha importado existe na primitiva de apresentação, não como override RGBA normalizado da instância; radio tem ao menos uma opção e selectedOptionId existente ou nulo; formatos não implicam validação CPF/data. Teclado aponta a campo da mesma tela, validado pela aplicação.

## Limites e relacionamento com outros domínios

Catálogo é definição embutida e sem IO. Screen consome suas definições; apresentação consulta metadados para expandir formas. Perfil EFL no adaptador resolve suporte efetivo, sem importar o adaptador no domínio.

Não implementa widget, drag/drop, perfis EFL externos, plugins executáveis ou templates por projeto.

## Estado e validação

Comando inválido devolve erro sem publicação parcial. Comandos válidos são confirmados pela aplicação, que incrementa revisão e coordena histórico. Este domínio não modifica diretamente a sessão da UI.

Testes obrigatórios: sete defaults, enums desconhecidos, cores inválidas, opções com vírgulas, campo de família errada e catálogo incompatível. Campos e enums do wire format são definidos pelo [contrato compartilhado](../001-gerador-jornada/contracts/component-catalog.md); a tradução DTO/domínio fica no adaptador.

## Rastreabilidade

Requisitos FR deste domínio cobrem FR-003, FR-005, FR-015, FR-022; NFR-002 da feature coordenadora. [Backlog local](tasks.md) detalha os arquivos e testes. Implementação e integrações verificadas estão identificadas em tasks.md; cenários não exercitados permanecem pendentes, conforme o registro de evidências da coordenadora.


## Correspondência com o código atual

Catálogo atual com dez Kinds, união tipada, Lookup e validação de enums. As 16 famílias eram o alvo da expansão de catálogo 2; não estão integralmente implementadas nesta árvore. Definition contém Kind/Name/Width/Height; metadados completos de propriedades/capacidades ainda não são um modelo declarativo. Inspector usa formulários explícitos por família. Perfil EFL segue no adaptador e não está homologado.

## Modelo alvo de catálogo 2 — 2026-10-19

A base histórica tinha catalogVersion=1; a expansão de catálogo 2 foi planejada para button/label/title/textField/radio/report/numericKeyboard e toggleButton/textArea/passwordField/comboBox/table/alphaKeyboard/progressBar/panel/dialog. UF é preset comboBox, Frame preset de tela, duas legendas variante button, oito classes Report variantes de report. EFLValorSurpresaSuperSete é composição label/toggle, não cálculo de loteria.

Definition evolui para metadados tipados de default/propriedades/variantes/estados/capacidades por alvo; não inclui widget/factory/renderizador/importador. Propriedades específicas constam no [contrato catálogo 2](../001-gerador-jornada/contracts/component-catalog.md#catálogo-2--expansão-autorizada-2026-10-19). Option/Row/Column/Cell e grupos preservam IDs/ordem/referências; seleção inicial/enabled pertencem aos dados quando aplicáveis; foco/pressed/seleção corrente/scroll/modal são runtime da aplicação. Validar selected IDs, limites/valores, tipo de linha/célula, alvo compatível de teclado e grupo da mesma tela.

O perfil por documento descrito nos registros de 2026-10-19 é histórico/plano de expansão; não existe na revisão 1.0.0.1. A implementação atual usa Style e biblioteca incorporada global, sem captura/selamento de revisões privadas. O histórico de build e validação continua no quickstart; não o repetir como evidência vigente deste modelo.

SkinDefinition/SkinState/FontFace, RGBA normalizado, aliases/revisões privadas, descritores de capacidades e galeria das 16 famílias são alvos ainda incompletos. O slice atual de skin tem Definition/State/Font/Resource e dados gerados; não possui Profile/Archived/Capture. Fonte desconhecida deve continuar literal no rascunho, com diagnóstico e bloqueio da saída fiel, sem associação fictícia.

Todo slice/map/propriedade nova entra em Copy/Validate/defaults, criação e fixtures; Kind novo não basta. Button novo mantém 110 × 60; keyboard novo 200 × 350 sem resize de antigos. Defaults desconhecidos aguardam evidência T009. T010–T016 locais não são marcados concluídos por este modelo documental.

<a id="modelo-codigos-barras-2026-10-23"></a>
## Modelo planejado — códigos de barras (2026-10-23)

FR-013–FR-015, [contrato canônico](../001-gerador-jornada/contracts/component-catalog.md#contrato-codigos-barras-2026-10-23). Nomes abaixo são tipos/variantes **a implementar**, não símbolos Go existentes:

- `barcodeInput`: propriedades tipadas com `label` (default “Número do Código de Barras:”). Não persistir valor inicial nesta v1: os oito campos sempre começam vazios. Estado de habilitação/bounds/estilo continua comum à instância de Tela. Oito regiões são derivadas, sem IDs de componentes adicionais.
- `barcodeProcessed`: propriedades tipadas `blocks`, lista ordenada não vazia. Cada bloco contém `id` UUID, `title` string opcional/vazia e `rows`, lista ordenada não vazia. Cada linha contém `id` UUID, `labelKind` enum `cpf/cnpj/name/tradeName/corporateName` e `value` string Unicode, inclusive vazia. Rótulo visível é derivado do enum; não gravar um texto de rótulo concorrente. Defaults: um bloco com CNPJ/Nome e valores vazios; nova linha é `name` vazia. Repetir tipo/título é permitido, repetir UUID não.
- `BarcodeValue`/resultado local: cinco blocos lógicos, representações de 44/47 dígitos e erros localizados; valor monetário em centavos exatos, banco textual, fator e identificadores derivados. Esses conceitos puros são reutilizados pela aplicação, não propriedades a gravar nem entidade de pagamento. Sem atualização de data por relógio da máquina.

Cópia profunda conserva IDs para snapshot; duplicação em Tela gera novos IDs de instância/blocos/linhas. Validar união discriminada, enum, lista mínima e identidade; a aplicação verifica unicidade global do documento. Definição de tipo/default não incorpora widgets, arquivos, payload de leitora, mensagens modais ou funções de UI. Não usar structs anônimas/JSON livre para contornar tipagem.

Versão efetiva do writer: `catalogVersion=4`; a revisão 3 introduziu os compostos e a 4 acrescentou TargetJourney e não equivale à expansão integral dos tipos planejados. Os contratos de Projeto discriminam leitores e escritores; não enriquecer perfis antigos ou inventar definições Java/XML para estes compostos novos.

<a id="modelo-menu-aplicacao-2026-10-23"></a>
## Modelo canônico planejado — Menu aplicação

API Go acordada, **não declaração de implementação validada**:

- `Kind ApplicationMenu = "applicationMenu"`; união `Properties.Menu *ApplicationMenuProperties` presente somente neste Kind.
- `ApplicationMenuProperties{Root SubmenuProperties; InitialMenuID string}`: raiz por valor e UUID do submenu inicial; factory seleciona `Root.ID` por padrão. InitialMenuID deve resolver raiz ou descendente existente, nunca título/índice. Não usar vazio como fallback implícito.
- `SubmenuProperties{ID, Name, Title, Layout string; Items []MenuItem}`: UUID único, nome/título Unicode (inclusive vazios), ordem do slice significativa. Layout é exatamente `root`, `standard`, `financial` ou `bolao`; raiz usa `root`, descendentes usam os outros layouts. Lista vazia é válida.
- `MenuItem{ID, Label, Legend, Category, Variant string; Enabled bool; Submenu *SubmenuProperties; TargetScreenID, TargetJourney, MovementCode, MovementOption, MovementName, Shortcut string}`: UUID próprio; Label é o rótulo do botão e Legend é a legenda ilustrativa independente de MovementName, Unicode literal, inclusive vazia, sem preenchimento implícito a partir de metadados. Categoria/metadados também são literais, inclusive vazios; Category é texto editável, não domínio de negócio. Variant identifica estilo suportado pelo catálogo e não código executável; desconhecido tem diagnóstico de estilo, sem troca silenciosa. Enabled é explícito; factory de item e seed usam true. Submenu ausente/nil representa folha; TargetScreenID vazio representa ausência de destino. Submenu, TargetScreenID e TargetJourney são mutuamente exclusivos. TargetJourney contém apenas o nome relativo de uma `.jornada` na raiz do workspace; vazio não infere destino.

Layout `bolao` sempre representa três posições de botões com Legend, independentemente do Name editável; excedentes geram diagnóstico e recorte, sem poda da árvore. No seed, BOLÃO ordinário usa `standard`, preservando jogos e rótulos; BOLÕES DIGITAIS usa `bolao` com Legend, na ordem: “Visualizar cotas digitais para venda na lotérica.”, “Consulta bolões de concursos abertos para comercialização.” e “Consulta bolões de concursos encerrados e não prescritos.”. Navegação usa somente itens pertencentes ao nó autorado, sem importar irmãos por nome.

Menus/submenus e itens têm IDs distintos entre si e dos demais UUIDs do documento, sem FileID próprio. Limites por árvore: profundidade <=16 (raiz=1), total de MenuItem <=1000 somando descendentes; menus vazios não violam limite mínimo. Ciclos, UUID inválido/repetido e inicial não resolvido são erros estruturais; referências de destino são validadas em Jornada. Nenhum limite é relaxado por desabilitar um ramo. Cópia profunda conserva IDs/ordem e não compartilha ponteiros/slices; duplicação/regeneração/remapeio pertence a Tela.

Factory recebe callback `supplyID` que entrega UUIDs; não lê relógio, arquivo ou rede. Seed único deve transcrever serviços/jogos e submenus aninhados da fonte Java/XML fornecida, incluindo Bolão digital; não copiar sua árvore inteira neste documento. MovementCode/MovementOption/MovementName/Shortcut não disparam ação, navegação, QR ou lookup. Categoria/variante não inferem destino. Caminho atual, modo raiz/submenu, foco e mensagem de seleção não são propriedades persistidas.

## Vínculo visual e fonte integral — 1.0.0.2

Sem mudança da união ou versão: `MenuItem.Variant` guarda a chave XML para preservar a marca após renomear. Variante vazia em projeto anterior pode resolver marca conhecida somente na apresentação; não altera dados persistidos. A fábrica usa dados estáticos do XML preservado e mantém os metadados nos ramos e folhas. Títulos principais dos ramos alimentam Submenu.Title; não existe título separado em MenuItem. Legendas digitais são complementos visuais documentados, não substitutos dos nomes XML.

## Ajuste dos defaults — 1.0.0.3

Sem novos campos/versões. Submenu.Title conserva literalmente tituloPrincipal, inclusive vazio; Name não o preenche. Variant pode referenciar menu2/menu2_seta/menu_peq ou marca reutilizada no Bolão conforme Java. Os campos Movement* e Shortcut existentes passam a ter autoria no formulário, mantendo identidade, cópia e confirmação atômica. O cadastro XML de movimentos permanece indisponível; não é equivalente à árvore de MenuItem.

## Parcela instalada integrada à 038

038/T019 implementa catalog/capabilities.go: PropertyDefinition e enums puros das sete famílias I1, com DefaultProperties compartilhado por screen.NewComponent e descoberta. Dez Kinds existentes preservados; compostos não anunciam patch avançado MCP. InstanceIDs/RegenerateInstanceIDs incluem opções/linhas/partes para cópias independentes. Não implementa o metamodelo completo do inspector, famílias novas, ReusableDefinition ou homologação. T004 amplo permanece aberto. Evidência de testes em 038/quickstart.
