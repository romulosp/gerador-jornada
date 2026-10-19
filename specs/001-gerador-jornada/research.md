# Research: Gerador de Jornada

**Date**: 2026-10-17
**Input**: [spec.md](spec.md), plano aprovado e inspeção estática de `material-legado`.
**Status**: decisões documentadas; provas de execução e homologações ainda não realizadas.

## Findings: análise do legado

A análise é de código-fonte, sem instalar ou executar o projeto React. Não há JAR EFL, fonte Java independente, exemplo `.page` ou apresentação PPTX no material recebido. O manifesto declara React 16.14, EJS e html2canvas, embora o campo de engine Node já tenha sido atualizado: a idade da arquitetura não pode ser deduzida apenas desse campo.

| Capacidade / evidência | Conclusão para a nova versão |
| --- | --- |
| [App.js](../../material-legado/src/App.js), estado e `_importPage` | Telas/componentes são reconstruídos como instâncias; edição e persistência estão acopladas ao estado central. Substituir por modelo de dados independente. |
| [PageService](../../material-legado/src/service/PageService/index.js), `pagesSave` | Persistência em localStorage; erro de leitura pode iniciar página aleatória. Preservar dados e apresentar recuperação explícita. |
| [Header](../../material-legado/src/components/Header/index.js), `onSave`/`onImport` | Exporta/importa array JSON com extensão `.page`. Nova versão usará pacote versionado; importação antiga ficou fora do escopo. |
| Header, `onScreenshot` | Exporta PNG via html2canvas. Não é exportação PPT nem atende a editabilidade por componente. |
| [Design](../../material-legado/src/components/Design/index.js) | Há drag/drop, arraste e movimentação por setas; dimensões e alinhamento assumem área 800 × 600. Preservar intenção e remover acoplamento a dimensões fixas. |
| [ToolsRight](../../material-legado/src/components/ToolsRight/index.js) e [JComponent](../../material-legado/src/components/ListComponents/JComponent/index.js) | Schema visual inclui tipos, tamanho e posição. Nova versão terá definições tipadas, sem funções serializadas. |
| [Catálogo](../../material-legado/src/components/ListComponents/index.js) | Sete componentes registrados. São as sete famílias iniciais do novo catálogo. |
| [EFLButton](../../material-legado/src/components/ListComponents/EFLButton/index.js), `events.click` | A prévia implementa navegação quando `eventType=NAVEGAR`. `EXIBIR` aparece como opção, mas não tem comportamento correspondente nesse handler. Não importar opção sem semântica. |
| EFLButton, `code` | Listener gerado contém TODO: a navegação da prévia não é transportada para Java. Nova geração precisa implementar o roteamento. |
| [Code](../../material-legado/src/components/Code/index.js) e [base](../../material-legado/src/components/Code/base.js) | EJS agrega imports, declarações e métodos dentro da UI. `UIControle`, `uiManager` e skins dependem de contexto não fornecido. Separar geração e perfil de integração. |
| [EFLRadio](../../material-legado/src/components/ListComponents/EFLRadio/index.js), `code` | Declara grupo de botões, mas gera getter/construtor de EFLButton: inconsistência de template, não contrato válido. |
| [EFLLabel](../../material-legado/src/components/ListComponents/EFLLabel/index.js) | Prévia oferece fonte/cor, mas geração usa skin fixa sem mapear toda personalização. Exigir matriz de capacidades por alvo. |
| [EFLTextField](../../material-legado/src/components/ListComponents/EFLTextField/index.js) | Lista formatos de documentos EFL; não há definição real desses documentos. Não presumir validação CPF/data apenas pelo nome. |
| [EFLReportEvent](../../material-legado/src/components/ListComponents/EFLReportEvent/index.js) | Lista de opções e listener de seleção existem no template, com ação TODO. Não declarar integração de negócio já implementada. |

Não foram identificados modelo de roteiro, animações, empacotamento Windows ou separação de exportadores. Essas são capacidades novas. Não foram executados testes de compilação do legado; os defeitos acima são observações estáticas.

## Decisions

### D-01 — Go/Fyne para a janela desktop

**Decision**: Go com Fyne v2 estável, versões exatas fixadas na Sprint 1. O editor não usará navegador nem WebView.

**Rationale**: Uma UI desenhada em Go atende a janela própria e permite distribuir sem instalar ambiente de linguagem. Ferramentas Go/C e drivers são considerações do desenvolvimento e da matriz de máquinas. A documentação Fyne distingue esses pré-requisitos do uso final e descreve o empacotamento EXE. [Preparação](https://docs.fyne.io/started/quick/), [empacotamento](https://docs.fyne.io/started/packaging/)

**Alternatives considered**: navegador local foi recusado pelo usuário; Wails/WebView exige aceitar um componente de execução externo; Qt amplia distribuição e ferramentas.

**Trade-off**: Canvas, drag/drop e representação EFL exigem widgets/renderer próprios. Validar DPI e driver gráfico na Sprint 1 e no Windows limpo antes de declarar distribuição homologada.

**Review when**: a prova técnica mostrar incompatibilidade incontornável com as máquinas Windows de destino. Mudança exige revisar plan.md e a restrição de runtimes, sem fallback silencioso.

### D-02 — Monólito modular e composição manual

**Decision**: um módulo Go, cinco módulos de domínio e adaptadores por fronteira; injeção por construtores.

**Rationale**: Aplicação local e individual não precisa de infraestrutura distribuída. Interfaces de persistência e exportação permitem testar falhas reais; structs concretas bastam nas regras internas.

**Alternatives considered**: microserviços aumentam distribuição e falhas; camadas globais `models/services/controllers` espalham uma mudança de domínio; framework DI adiciona indireção sem necessidade atual.

**Trade-off**: A disciplina de importações precisa ser verificada. Se uma responsabilidade crescer, separar arquivos/tipos antes de criar novo serviço ou módulo Go.

### D-03 — Arquivo portátil e modelo tipado

**Decision**: `.jornada` como ZIP de manifesto JSON v1 e recursos incorporados; DTOs de serialização ficam no adaptador.

**Rationale**: Um projeto pode ser copiado offline e incluir suas imagens. As sete famílias conhecidas permitem propriedades tipadas, eliminando métodos serializados do modelo legado.

**Alternatives considered**: JSON solto perde portabilidade de imagens; SQLite não resolve uma necessidade de consulta/multiusuário desta versão; manter `.page` reproduz um contrato não versionado.

**Trade-off**: São necessários limites de ZIP, validação, autosave e escrita atômica. Sem migração automática de formatos futuros ou de `.page`.

### D-04 — Plano visual compartilhado

**Decision**: expandir componentes em primitivas independentes de UI para renderizar desktop e PPT. Java recebe a mesma geometria e semântica, com mapeamento EFL próprio.

**Rationale**: Captura de tela não preserva controle editável. Um plano comum evita dois algoritmos incompatíveis de layout.

**Alternatives considered**: extrair geometria de widgets Fyne acopla exportação à janela; rasterizar tudo inviabiliza editabilidade.

**Trade-off**: Métricas de fonte variam entre plataformas; homologar geometria com tolerâncias e registrar diferenças de skin, sem prometer rasterização idêntica.

### D-05 — PPTX diretamente em Go

**Decision**: implementar subconjunto explícito de OPC/PresentationML usando ZIP/XML padrão, sem Office COM, Node, Python ou SDK .NET no editor.

**Rationale**: O formato contém links internos e uma árvore temporal de animações, permitindo cumprir as duas saídas de PPT com controle dos IDs e alvos. [Animações Microsoft](https://learn.microsoft.com/en-us/office/open-xml/presentation/working-with-animation), [ação de navegação interna](https://learn.microsoft.com/en-us/openspecs/office_standards/ms-oe376/7ff3db24-b7b9-4ffe-aa78-3ec47cab2489)

**Alternatives considered**: screenshots eliminam editabilidade; automação Office viola execução sem Office; bibliotecas que só criam slides estáticos não resolvem a demonstração.

**Trade-off**: Implementar um subconjunto exige mais testes que um exportador estático. A Sprint 1 deve demonstrar links, trajetória de cursor, clique e digitação editável, abrindo sem reparo no PowerPoint.

**Review when**: uma biblioteca Go mantida cobrir esse subconjunto com evidências e reduzir esforço sem mudar os contratos.

### D-06 — Exploração separada da demonstração

**Decision**: um PPTX com menu e dois percursos. Exploração tem slides manuais; demonstração tem ocorrências de telas e tempos finitos. Roteiros são listas explícitas, sem gravação ou scripts.

**Rationale**: Clique livre e avanço automático têm comportamentos diferentes. Separar os percursos impede que um temporizador interrompa a exploração e permite revisitar telas com estados próprios.

**Trade-off**: O arquivo terá slides adicionais. Campos de PPT não aceitam digitação real durante apresentação; a demonstração usa texto editável e efeitos.

### D-07 — Fontes Java 1.6 com perfil EFL

**Decision**: templates internos por componente, classes por tela e roteamento separado. Perfil EFL versionado, homologado no sistema real antes de publicar a capacidade Java.

**Rationale**: EFL foi explicitamente escolhido pelo usuário; substituí-lo por Swing puro não atende o alvo. A documentação Swing é referência para Java SE 6 e threading, não prova das assinaturas proprietárias. [Swing Java SE 6](https://docs.oracle.com/javase/6/docs/api/javax/swing/package-summary.html)

**Alternatives considered**: fontes Swing puras eliminam a biblioteca solicitada; compilar Java no editor reintroduz JDK na distribuição; copiar templates antigos preserva erros.

**Trade-off**: É necessária participação da equipe EFL. A geração em si não usa JDK, mas a homologação usa compilador, API SE 6 e bibliotecas reais.

### D-08 — Documentação completa, com GoDoc e Javadoc

**Decision**: seguir a constituição, ampliando o mínimo GoDoc para elementos privados e comentários de variáveis locais, conforme pedido do usuário.

**Rationale**: Comentários de declaração dão contexto a APIs; comentários locais explicam finalidade, mas não são publicados em `go doc`. Não tratar esses mecanismos como equivalentes. [Go Doc Comments](https://go.dev/doc/comment)

**Trade-off**: Há custo contínuo de manutenção. Comentários precisam acompanhar mudanças; linter cobre apenas parte da exigência, complementada por revisão.

### D-09 — Alvo Java 6 com ferramenta executada no JDK 17

**Context**: o usuário disponibilizou `C:\Desenvolvimento\jdk-17.0.11` e reafirmou o alvo Java 1.6. O `javac --help` local aceita source/target/release de 7 a 17; não compila diretamente para Java 6.

**Decision**: executar Eclipse ECJ 3.33.0 no JDK 17 com source/target 1.6; verificar classfile major 50 e as referências de API com Animal Sniffer 1.24, ASM 9.7 e assinatura java16:1.1. Ferramenta externa ao editor, com versões, origem e hashes fixados em `apps/gerador-jornada/scripts/java6-tools.json`. Sem dependência Maven ou mudança de JAVA_HOME global.

**Rationale**: resolve a compilação estática no ambiente fornecido sem elevar o alvo. ECJ permite execução independente e configuração de source/target; Animal Sniffer compara referências de classes, métodos e campos com assinaturas da plataforma. O verificador combina símbolos dos fontes e JARs reais com os de SE 6 e verifica também o conteúdo das dependências, sem ignorar pacotes. [ECJ](https://help.eclipse.org/latest/topic/org.eclipse.jdt.doc.user/tasks/task-using_batch_compiler.htm), [API Animal Sniffer](https://www.mojohaus.org/animal-sniffer/animal-sniffer/index.html), [assinaturas Java SE](https://www.mojohaus.org/signatures/index.html).

**Alternatives considered**: `javac --release 6` não é aceito pelo JDK 17 local; somente source/target no ECJ permite chamadas a APIs modernas; exigir outra instalação Java para toda compilação adiciona uma dependência que o ambiente fornecido já permite evitar.

**Trade-off**: ferramentas adicionais são baixadas uma vez para desenvolvimento, com modo offline após preparação. A verificação estática cobre sintaxe, bytecode e símbolos; não demonstra disponibilidade operacional da JVM 6, comportamento visual nem bootstrap EFL. Um perfil Swing exclusivo de teste comprova a estrutura gerada, sem liberar o perfil EFL do produto. Revisar a decisão se o ambiente EFL fornecer uma toolchain oficial ou exigir APIs específicas de fornecedor.

## External dependencies and validation gates

| Gate | Evidência necessária | Responsável funcional | Efeito se ausente |
| --- | --- | --- | --- |
| G-01 Desktop | EXE com canvas editável e DPI testado em Windows-alvo | Equipe do editor | Impede concluir prova técnica desktop e release |
| G-02 PPT | Arquivo de prova sem reparo, formas editáveis e efeitos executados nas duas versões-alvo | Equipe do editor e homologação | Impede concluir exportação correspondente |
| G-03 EFL | JARs/versionamento, dependências transitivas, exemplo funcional de bootstrap, gestores UI e skins; APIs inspecionadas | Mantenedores do sistema EFL | Impede homologar perfil e concluir US5 |
| G-04 Java 6 | Compilador/toolchain com suporte real a source/target 1.6 e API SE 6; execução no ambiente de destino | Mantenedores EFL e build | Impede declarar compatibilidade Java |
| G-05 Release | Máquina Windows limpa e offline; versões e dependências efetivas registradas | Homologação | Impede release validado |

G-03 não é resolvido adivinhando classes. Levantar construtores, constantes, listeners, EDT, vinculação do teclado e ciclo de vida das sete famílias. Evidências observadas no legado ficam marcadas como não verificadas até esse levantamento.

## Imported skills: origem e integridade

Cópias integrais em 2026-10-17, da origem `D:/desenvolvimento/ia/aula02/.agents/skills/`. Cada diretório contém apenas SKILL.md; não foram alterados os textos.

| Skill | Cópia local | SHA-256 verificado |
| --- | --- | --- |
| golang-project-layout | [SKILL.md](../../.agents/skills/golang-project-layout/SKILL.md) | `E90E8A0B56F5C775AF39F0009627EEDAB918CE89CD6BB3D79B4F29334A727630` |
| golang-structs-interfaces | [SKILL.md](../../.agents/skills/golang-structs-interfaces/SKILL.md) | `9249AB52A43ADA4DA4DB808993FAB95E677720CFE2DC587057CB8C5B4FAFBA25` |
| golang-documentation | [SKILL.md](../../.agents/skills/golang-documentation/SKILL.md) | `4BE01B6A5EFF702B00DCEF7A73478BDBB048D4057F5AF8DF73402A886905E0BD` |

Referências a `references/`, `assets/` e outras skills `samber/cc-skills-golang@...` não vieram acompanhadas dos recursos. Não são dependências locais existentes. Esta especificação contém as decisões necessárias sem depender delas.

A arquitetura e a composição manual já foram aprovadas: a orientação genérica da skill para perguntar essas preferências não reabre decisões. Recomendações web/12-Factor ou de publicação pública não são aplicadas automaticamente ao desktop. A constituição prevalece sobre orientações genéricas.

## Research conclusion

As escolhas de produto e arquitetura estão definidas e há implementação em andamento. Há evidência automatizada e builds registrados em quickstart.md; G-01 a G-05 só são concluídos pela evidência específica exigida em cada linha, incluindo execução real e homologação quando aplicáveis. Os registros de decisão permanecem neste arquivo para manter a estrutura GitHub Spec Kit solicitada.

Atualização da implementação: a parte estática de G-04 tem ferramenta executável e testes, conforme D-09 e quickstart.md. G-03 e a execução Java 6/EFL de G-04 continuam pendentes. Não foram encontrados JARs EFL no material legado nem na busca por nomes EFL sob `C:\Desenvolvimento` e `D:\desenvolvimento`; isso não identifica bibliotecas com outros nomes nem substitui o inventário dos mantenedores.

## Inventário de componentes e designs individuais

Em 2026-10-18, o usuário solicitou uma spec e um DESIGN.md por componente. Foram revisadas as dezesseis pastas diretamente sob components e as sete famílias registradas em ListComponents. Pastas de ferramentas da interface possuem responsabilidade própria, sem virar famílias arrastáveis ou novos domínios de negócio. O cabeçalho composto e a faixa de título também têm contratos individuais.

As 24 features abaixo são especificações individuais em estado Draft; não equivalem a 24 capacidades já homologadas. Reutilizam os domínios e contratos existentes. Sua criação não muda a feature coordenadora selecionada nem cria branches Git. Novos plan.md/tasks.md individuais serão produzidos quando essas features passarem por planejamento; o pedido atual entrega spec e design.

| Componente | Responsável | Evidência no legado | Especificação | Personalização |
| --- | --- | --- | --- | --- |
| Botão — EFLButton | 003-catalogo | `ListComponents/EFLButton/index.js` | [spec](../007-botao-efl/spec.md) | [design](../007-botao-efl/DESIGN.md) |
| Rótulo — EFLLabel | 003-catalogo | `ListComponents/EFLLabel/index.js` | [spec](../008-rotulo-efl/spec.md) | [design](../008-rotulo-efl/DESIGN.md) |
| Título — EFLLabelTitulo1 | 003-catalogo | `ListComponents/EFLLabelTitulo1/index.js` | [spec](../009-titulo-efl/spec.md) | [design](../009-titulo-efl/DESIGN.md) |
| Campo de texto — EFLTextField | 003-catalogo | `ListComponents/EFLTextField/index.js` | [spec](../010-campo-texto-efl/spec.md) | [design](../010-campo-texto-efl/DESIGN.md) |
| Seleção — EFLRadio | 003-catalogo | `ListComponents/EFLRadio/index.js` | [spec](../011-selecao-efl/spec.md) | [design](../011-selecao-efl/DESIGN.md) |
| Relatório/lista — EFLReportEvent | 003-catalogo | `ListComponents/EFLReportEvent/index.js` | [spec](../012-relatorio-efl/spec.md) | [design](../012-relatorio-efl/DESIGN.md) |
| Teclado numérico — EFLVirtualKeyboardNumeric | 003-catalogo | `ListComponents/EFLVirtualKeyboardNumeric/index.js` | [spec](../013-teclado-numerico-efl/spec.md) | [design](../013-teclado-numerico-efl/DESIGN.md) |
| Data e hora — DateInfo | 004-tela | `DateInfo/index.js` | [spec](../014-data-hora/spec.md) | [design](../014-data-hora/DESIGN.md) |
| Dados da lotérica — LotericaInfo | 004-tela | `LotericaInfo/index.js` | [spec](../015-dados-loterica/spec.md) | [design](../015-dados-loterica/DESIGN.md) |
| Indicadores — InfraInfo | 004-tela | `InfraInfo/index.js` | [spec](../016-infraestrutura/spec.md) | [design](../016-infraestrutura/DESIGN.md) |
| Cabeçalho informativo da tela | 004-tela | `Design/index.js` | [spec](../017-cabecalho-tela/spec.md) | [design](../017-cabecalho-tela/DESIGN.md) |
| Faixa de título da página | 004-tela | `Design/page.css` | [spec](../018-faixa-titulo/spec.md) | [design](../018-faixa-titulo/DESIGN.md) |
| Área de desenho — Design | 004-tela | `Design/index.js` | [spec](../019-canvas-edicao/spec.md) | [design](../019-canvas-edicao/DESIGN.md) |
| Elemento editável — PageElement | 004-tela | `PageElement/index.js` | [spec](../020-elemento-edicao/spec.md) | [design](../020-elemento-edicao/DESIGN.md) |
| Elemento em execução — PageElementPreview | 001-gerador-jornada | `PageElementPreview/index.js` | [spec](../021-elemento-previa/spec.md) | [design](../021-elemento-previa/DESIGN.md) |
| Lista de telas — PageTab | 004-tela | `PageTab/index.js` | [spec](../022-lista-telas/spec.md) | [design](../022-lista-telas/DESIGN.md) |
| Paleta — ComponentTab | 003-catalogo | `ComponentTab/index.js` | [spec](../023-paleta-componentes/spec.md) | [design](../023-paleta-componentes/DESIGN.md) |
| Estrutura — StructureTab | 004-tela | `StructureTab/index.js` | [spec](../024-arvore-estrutura/spec.md) | [design](../024-arvore-estrutura/DESIGN.md) |
| Barra de ferramentas — Header | 001-gerador-jornada | `Header/index.js` | [spec](../025-barra-ferramentas/spec.md) | [design](../025-barra-ferramentas/DESIGN.md) |
| Painel esquerdo — ToolsLeft | 001-gerador-jornada | `ToolsLeft/index.js` | [spec](../026-painel-esquerdo/spec.md) | [design](../026-painel-esquerdo/DESIGN.md) |
| Painel direito — ToolsRight | 001-gerador-jornada | `ToolsRight/index.js` | [spec](../027-painel-direito/spec.md) | [design](../027-painel-direito/DESIGN.md) |
| Formulários de propriedades — Properties | 001-gerador-jornada | `Properties/index.js` | [spec](../028-propriedades/spec.md) | [design](../028-propriedades/DESIGN.md) |
| Prévia interativa — Preview | 001-gerador-jornada | `Preview/index.js` | [spec](../029-previa-interativa/spec.md) | [design](../029-previa-interativa/DESIGN.md) |
| Código e diagnóstico — Code | 001-gerador-jornada | `Code/index.js` | [spec](../030-visualizacao-codigo/spec.md) | [design](../030-visualizacao-codigo/DESIGN.md) |
| Catálogo / ListComponents | 003-catalogo | `ListComponents/index.js` | [spec](../003-catalogo/spec.md) | [design](../003-catalogo/DESIGN.md) |

Properties/index.js era somente um contêiner vazio; os formulários funcionais estavam em ToolsRight. JComponent é uma base abstrata, mapeada ao modelo compartilhado do Catálogo, sem um controle fictício na paleta. O relógio de DateInfo era atualizado no navegador; na jornada nova os valores são ilustrativos e editáveis para manter demonstrações determinísticas. Code usava templates com inconsistências; a visualização/publicação Java permanece condicionada à homologação EFL real.

Foi tentada execução de uma cópia isolada do legado em .tools/legacy-preview, sem alterar material-legado. A inicialização encontrou dependência AJV inconsistente; os navegadores desta sessão não estão disponíveis pelo conector. A análise usada para as specs foi do código, CSS e imagens locais, sem afirmar teste visual do legado em execução.

## Revisão de interação e transições — 2026-10-18

O usuário solicitou suspensão temporária de Java/EFL, transições de apresentação, confirmação de propriedades acessível, seleção múltipla/atalhos, cursor contextual, redução do custo de seleção e alinhamento do texto da faixa. Esses pontos foram registrados nas specs individuais e no domínio responsável, com rastreabilidade na spec coordenadora.

Decisão de edição: confirmar por botão fixo Salvar alterações, em vez de publicar cada caractere ou número intermediário inválido. Descartar reconstrói o rascunho a partir dos dados confirmados. O custo é uma confirmação explícita; o benefício é uma edição transacional compreensível, sem centenas de entradas de Undo.

Decisão de renderização: cache de formas e imagens por sessão/tela/revisão/zoom, reaproveitando posições e contornos de seleção. Se a seleção não muda, não reconstruir o inspector. Não introduzir índices espaciais ou paralelismo no domínio sem medição que os justifique.

No mesmo microbenchmark `BenchmarkSelect500`, Windows amd64, Intel Core i5-13600KF, Go 1.27.2 e driver Fyne de teste, dez iterações: antes 14.251.660 ns/op, 1.776.739 B/op e 27.175 alocações/op; após a correção 704.540 ns/op, 483.956 B/op e 9.004 alocações/op. O tempo por seleção medido caiu de aproximadamente 14,25 ms para 0,70 ms. Trata-se de benchmark de teste com rótulos, não latência garantida da janela nativa ou de projetos com muitas imagens.

Decisão de transições: fade de 300 ms por padrão, push/wipe opcionais e opção sem efeito. São recursos disponíveis no PresentationML das versões-alvo; duração usa p14:dur com namespace de compatibilidade. Não adotar Morph nem rasterizar controles. A execução visual/temporal real permanece dependente de homologação no PowerPoint. [Transições](https://learn.microsoft.com/en-us/dotnet/api/documentformat.openxml.presentation.transition?view=openxml-3.0.1), [Push](https://learn.microsoft.com/en-us/dotnet/api/documentformat.openxml.presentation.pushtransition?view=openxml-3.0.1).

Na API Fyne instalada, Ctrl+A chega como ShortcutSelectAll e Shift+setas passa por evento de tecla com modificador do driver, sem CustomShortcut de Shift. O adaptador trata esses caminhos e preserva foco textual. Casos de integração em janela Windows continuam no roteiro manual.

## ID de arquivo separado de UUID e nome amigável — 2026-10-18

Confirmação do usuário: Pnl obrigatório somente em telas, sem prefixo obrigatório em componentes. Adotado FileID ASCII alfanumérico de até 40 caracteres, com início em letra; limite inclui prefixo e exclui extensão. A decisão preserva UUIDs internos, nomes amigáveis Unicode e navegação ao editar o arquivo de destino. A UI remove acentos/especiais; validação de domínio/aplicação continua protegendo gravação e importação independentes da UI.

Manifesto passou a schemaVersion 2 para impedir perda silenciosa do atributo em leitores antigos. Leitura v1 migra nomes ausentes com sugestões/sufixos determinísticos; leitura não grava nem renomeia IDs explícitos. Duplicação gera nomes disponíveis, enquanto exportação preserva o nome confirmado e recusa colisões/palavras reservadas Java. ID de componente é local à tela; diretórios/pacotes separados permitem repetição em outras telas. Variáveis Java usam prefixo técnico component para evitar colisões com host/navigator/factory, sem alterar o nome-base do arquivo principal.

Verificado: suíte Go e vet; nove testes Python da toolchain Java 6; compilação ECJ + assinatura SE 6 dos fontes estruturais com arquivos individuais, incluindo ID livre host para detectar ocultação de parâmetros. Não altera a suspensão do botão Java/EFL nem comprova compatibilidade com EFL real.

## Builds nativos e documentação de produto — 2026-10-18

A documentação oficial do [Fyne para preparação](https://docs.fyne.io/started/quick/), [cross-compilation](https://docs.fyne.io/started/cross-compiling/) e [pacotes desktop](https://docs.fyne.io/started/packaging/) confirma o custo das ferramentas/bibliotecas nativas. Adotado build por sistema com CGO, evitando alegar que mudar GOOS basta para distribuir uma GUI. Docker foi iniciado pelo usuário e disponibiliza compilação Linux local; não fornece o SDK Apple.

A [referência oficial de runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners) distingue macos-15-intel de macos-15 ARM64. A matriz explicita arquitetura e verifica o formato do binário; o macOS não é tratado como um único pacote universal. Assinatura ad-hoc serve ao teste interno e não representa assinatura Developer ID/notarização. Ações do workflow fixadas por commits consultados nos repositórios oficiais, sem execução remota nesta entrega.

Tradeoff: builds nativos/Docker exigem mais infraestrutura do que uma compilação Go sem CGO, mas respeitam as dependências gráficas e produzem artefatos verificáveis. O usuário final não precisa instalar a toolchain; Linux ainda requer bibliotecas gráficas compatíveis. O manual separa essas duas classes de pré-requisitos e os estados configurado/gerado/homologado.

## Maximização para tela cheia sem estado de negócio

A API pública Fyne 2.8.2 fornece Window.SetFullScreen/FullScreen e driver.NativeWindow, mas não expõe evento de maximização; conferido no módulo local. No Windows, adotada captura restrita de WM_SYSCOMMAND/SC_MAXIMIZE com [SetWindowLongPtrW](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowlongptrw), delegação de mensagens com [CallWindowProcW](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callwindowprocw) e máscara dos comandos descrita em [WM_SYSCOMMAND](https://learn.microsoft.com/en-us/windows/win32/menurc/wm-syscommand).

O custo é manter um adaptador específico e exigir teste nativo. Em troca, a intenção é capturada antes de a maximização comum substituir a geometria, sem polling ou dependências nos domínios. F12/Escape no HWND cobrem campos focados, pois o driver Fyne não encaminha suas teclas simples ao canvas quando um widget está focado; CustomShortcut sem modificadores também não é emitido nesse driver. Todos os demais eventos continuam no Fyne. Botão interno permite retornar quando a decoração do sistema está oculta. Linux/macOS não recebem uma implementação Win32 artificial e permanecem sujeitos a homologação própria.


## Temas de interface e encerramento (2026-10-18)

Escolha: paletas Fyne e Preferences já disponíveis, sem dependência adicional ou tema no projeto. Claro permanece padrão, Escuro atende ambientes com pouca luz e Azul oferece alternativa clara. Custo: manter tokens/contraste e conferir visualmente os estados nos três temas. Alterar apenas a interface evita contaminar as saídas da jornada. A preferência é local ao usuário, não um dado de domínio compartilhável.

Sair chama a mesma solicitação do fechamento nativo, com confirmação única para documento alterado ou rascunho do formulário; a janela principal é master para também encerrar com prévias abertas. Alternativa de chamar Quit diretamente foi descartada porque ignoraria a confirmação. A captura dos valores do formulário é local ao inspector, sem observar teclas ou substituir validações.
