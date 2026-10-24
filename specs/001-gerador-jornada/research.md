# Research: Gerador de Jornada

## Referência corrente — revisão 1.0.0.1

O estado do código está na [spec coordenadora](spec.md#estado-verificado-1-0-0-1) e na seção 22 do [quickstart](quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). As análises e evidências antigas abaixo são históricas ou requisitos de expansão; seus números, APIs e arquivos não devem ser apresentados como implementação da árvore atual. A reconciliação documental das 33 features está na [seção 23](quickstart.md#23-reconciliação-documental-das-33-specs--2026-10-24).

**Constatações desta revisão:** writer 3/3, dez Kinds; biblioteca global de 205 chaves/224 recursos/seis fontes, oito definições indisponíveis; sem AppearanceProfile/Capture/visualprofile, máscara incremental ou entrega das 16 famílias. Há runtime de destinos de menu, mas não grafo global completo. RELATÓRIOS tem títulos/tokens sintéticos, apesar da contagem 27/148. Os requisitos para resolver essas lacunas foram mantidos e oito tarefas tiveram o estado corrigido para abertas.

**Causa e correção do estilo:** o resolvedor não aplicava default de skin, a prévia o substituía por controles opacos e o desktop não carregava os bytes da face. Importação estática, correção de alpha, família/flags SFNT, cena comum e controles transparentes foram necessários para produzir as capturas da revisão. O [README da raiz](../../README.md#o-que-foi-necessário-para-recuperar-o-estilo-do-legado) documenta o procedimento reproduzível e seus limites.


**Date**: 2026-10-17
**Input**: [spec.md](spec.md), plano aprovado e inspeção estática de `material-legado`.
**Status**: decisões documentadas; distribuição Windows pós-Capture gerada e verificada. Checksum do ZIP no sidecar externo `.zip.sha256`; reempacotar documentos preserva o EXE. Execução Go bloqueada por política Windows, homologações manuais/EFL e aceite pendentes.

## Findings: análise do legado

A análise é de código-fonte, sem instalar ou executar o projeto React. Não há JAR EFL, fonte Java independente, exemplo `.page` ou apresentação PPTX no material recebido. O manifesto declara React 16.14, EJS e html2canvas, embora o campo de engine Node já tenha sido atualizado: a idade da arquitetura não pode ser deduzida apenas desse campo.

**Contexto histórico**: o parágrafo acima descreve o recebimento inicial de 2026-10-17. Em 2026-10-19 há fontes Java de componentes/skin em [elemetos_componente](../../material-legado/elemetos_componente). A investigação pós-v1 abaixo usa esse acervo como referência primária; não implica presença de JARs/delegates completos ou execução homologada.

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

## Matriz de cobertura Java/EFL — expansão 2026-10-19

**Origem/estado**: plano aprovado pelo usuário para implementação de todo o acervo visual. Esta é a matriz coordenadora, não um relatório de homologação. Inspeção estática e arquivos presentes são evidências de referência, não de execução do novo editor. Todas as entradas abaixo permanecem **abertas para implementação/verificação e não homologadas**, inclusive as que já têm uma representação simplificada nas sete famílias. T108 deve expandir esta classificação por **cada** chave XML e estado; não se presume contagem completa por busca textual truncada.

Referência primária: [skin.xml](../../material-legado/elemetos_componente/xml/skin.xml) → reader JAXB → [SkinBuilderImpl](../../material-legado/elemetos_componente/builder/impl/SkinBuilderImpl.java) → SkinState → paint/delegates/construtores. React/CSS ficam como referência secundária identificada. O nome `type1` não é global: identidade composta inclui tipo da skin e estado. Nomes originais, inclusive grafias históricas, são preservados no inventário.

### Cobertura por classe/responsabilidade

| Chave | Classe/referência Java | Mapeamento semântico e comportamento local previsto | Skin/estado e dependência de evidência | Dono / tarefas |
| --- | --- | --- | --- | --- |
| C01 | EFLButton | button; click/none/navigate | Todas skins button; default/touched/disabled declarados; EFLButtonUI ausente | 003/T011; 007/FR-008; 001/T111–T116 |
| C02 | EFLTwoLabelButton | variante button de duas legendas, não Kind novo | botao_two_label e partes/textos próprios; medir linhas separadas | 003/T011; 007/FR-008 |
| C03 | EFLLabel e skins de título | label/title com variantes | Todas label, incluindo títulos/biometria/notificação; EFLLabelUI ausente | 003/T011; 008–009/FR-008 |
| C04 | EFLValorSurpresaSuperSete | composição de label alinhado à direita + toggle opcional; sem negócio de loteria | Construtores comprovam add(botao), não cálculo externo | 003/T011; 004/T011; 008/FR-008 |
| C05 | EFLToggleButton, EFLButtonGroup | toggleButton e grupos; radio preserva seleção simples; multisseleção só por contrato tipado | Todas toggleButton; delegate ausente; distinguir selected de touched | 003/T011; 011/FR-008 |
| C06 | EFLTextField, EFLTextFieldAdapter | textField, documentos/máscaras e eventos locais | textField/type1, estados declarados; EFLTextFieldUI ausente | 003/T012; 010/FR-008 |
| C07 | EFLTextArea | textArea, multiline e recorte/rolagem | Todas textArea; pintura e documentos a verificar | 003/T012 |
| C08 | PasswordFieldSkin/XML, EFLPasswordDocument | passwordField, valores apenas ilustrativos | Todas passwordField; sem classe visual PasswordField fornecida; mascaramento não homologado | 003/T012; 001/T119 |
| C09 | EFLComboBox, EFLComboBoxUf | comboBox + preset UF, seleção local com IDs | Composição das skins aplicáveis; EFLComboBoxUI ausente | 003/T012 |
| C10 | EFLTable e TableHeaderSkin | table com colunas/células/seleção tipadas | Todas table e tableHeader; estados/linhas/cabeçalho não resumidos a Report | 003/T012 |
| C11 | EFLReportEvent | piloto report, cabeçalho/corpo/rodapé/seleção/scroll | Panels/labels/scroll compostos; topo 26/rodapé 32 em Java; título 18 é referência React a reconciliar | 003/T014; 012/T001–T005 |
| C12 | EFLReport, EFLReport2, EFLReport3 | variantes report; tipos de linha próprios | Cada classe requer fixture e contrato; não aprovada por reuso do piloto | 003/T014 |
| C13 | EFLReportEvent2, EFLReportNotificacao | variantes report/evento/notificação | Estados e linhas próprios; eventos externos só como contrato | 003/T014 |
| C14 | EFLReportSelecaoMultipla, EFLInternalReport | variantes report com seleção múltipla/interna | Seleção inicial tipada versus runtime; sem lógica de negócio | 003/T014 |
| C15 | EFLVirtualKeyboardNumeric | numericKeyboard; 14 teclas, foco/cursor/TAB | 200 × 350 em Java; teclas 52 × 52 e auxiliares 80 × 60; LIMPAR emite VK_BACK_SPACE, não prova clear-all | 003/T013; 013/T001–T005 |
| C16 | EFLVirtualKeyboard | alphaKeyboard, foco e entrada local | 800 × 200 como referência Java; layout/modificadores a inventariar integralmente | 003/T013 |
| C17 | EFLVirtualKeyboardCalculator | variante/helper teclado calculadora | Operações locais comprovadas, sem serviços/periféricos; dependências a classificar | 003/T013 |
| C18 | EFLProgressBar | progressBar de valor ilustrativo | panel/label de progresso, escala/estados a medir | 003/T014 |
| C19 | EFLPanel, EFLPanel2 | panel e variantes; composição de autor | Todas panel; EFLPanelUI ausente; preservar transparent/alpha | 003/T015; 004/T009–T011 |
| C20 | EFLScrollPane, EFLScrollPane2, EFLScrollBarUI | capacidade scroll/clip de contêiner/composto | Skins scroll e botões direcionais; não criar Kind por seta | 003/T015; 001/T114–T116 |
| C21 | DialogSkin/XML | dialog como composição/overlay local | dialog_notificacao e dependências; sem classe visual/delegate suficiente | 003/T015; 001/T119 |
| C22 | EFLFrame | preset de tela 800 × 600 | Sem recriação da janela Windows; conteúdo do frame, não chrome do SO | 004/T011 |
| C23 | composições 014–018 | Data/hora, lotérica, infraestrutura, cabeçalho, faixa | Reusar fontes/skins comprovadas; conservar valores estáticos; referências React identificadas | 004/T011; 001/T112–T116 |
| C24 | ComponentSkin/XML, SkinState/XML, classes *Skin/*SkinXML e beans de estado | Infraestrutura de dados, não controles na paleta | Todo tipo reconhecido pelo reader/builder deve ter classificação e diagnóstico de incompletude | 001/T108–T110; 003/T009 |
| C25 | SkinConfig, SkinRender, SkinPainter, SkinUtil, FontLoader, ImageLoader, reader JAXB e builder | Contrato de carga/pintura/métricas; não portar execução Java | SkinPainter é interface; esticar imagem aos bounds onde comprovado não prova nine-slice | 001/T108–T111 |
| C26 | documentos de entrada em java/ | Regras locais tipadas, não famílias visuais novas | Inventariar implementações/limites; nome de documento não prova validação fiscal | 003/T012; 010/FR-008 |

C26 abrange EFLDocument, EFLPlainDocument, EFLAlfaDocument, EFLAlfaNumDocument, EFLNumericDocument, EFLMonetaryDocument, EFLMaskedDocument, EFLMaskCelularDocument, EFLCelularDocument, EFLCepDocument, EFLCpfDocument, EFLCpfCnpjDocumento, EFLDateDocument, EFLShortDateDocument, EFLDddDocument, EFLNisDocument, EFLNsbDocument, EFLNsbcDocument, EFLNsbiDocument, EFLNsbMascaradoDocument, EFLNumeroContratoDocument, EFLPasswordDocument, EFLRicDocument e EFLUserDocument. EFLReportEventListener e adapters/listeners são contratos de evento, não novos Kinds. Dependência referenciada sem fonte também entra no inventário como ausente; lista acima não autoriza omiti-la.

### Registro granular obrigatório e estados independentes

Cada linha granular em T108/003-T009 deve registrar: chave Cxx; classe e método/construtor; tipo/nome/estado XML original; caminho lógico e físico; hash da origem/imagem/TTF; licença; tamanho/formato; fonte declarada (arquivo, tamanho, bold/italic), face efetiva/nome interno e evidência; RGBA; bounds internos; spacement/baseline/alinhamento/quebra/fit/clip; herança de estado comprovada; comportamento local/evento; consumidores canvas/miniatura/prévia/PPT; tarefa qualificada; status de implementação; teste automatizado/revisão; homologação manual por ambiente; aprovação; bloqueio, impacto e responsável. Campos não comprovados ficam `pendente`, nunca preenchidos por inferência de similaridade.

Os nomes de **todas** as entradas XML (panel, label, button, toggleButton, textField, textArea, passwordField, dialog, table/headers e scroll) devem ser expandidos por parser, incluindo variantes de marketplace/biometria/notificações e seus estados. Essa saída é manifesto de recursos rastreado nesta matriz, não backlog avulso. Detectar chaves duplicadas, colisões e referências não resolvidas; não usar apenas uma amostra de skins.

| Gate novo | Situação conhecida / impacto | Responsável funcional / tarefa |
| --- | --- | --- |
| G-06 Recursos visuais | Após importar 18 pacotes de outros_materiais em 2026-10-19, /icon/dupla_sena.png está disponível. Restam 16 referências a 16 caminhos ausentes, listados no manifesto gerado; não usar diagnóstico anterior de 138 caminhos. Bloqueia somente entradas dependentes, mas impede declarar fidelidade de todo acervo | Mantenedor do acervo + equipe editor; T108/T119 |
| G-07 Fontes | Os seis nomes referenciados possuem arquivo homônimo; 8/8 faces do inventário aceitas tecnicamente, com sentinela não canônica de futura_mh preservada/diagnosticada. Importador check e testes Python PASS não homologam métricas visuais nem direitos: 8 licenças desconhecidas e 7 fontes com embedding restrito permanecem bloqueios | Equipe recursos + responsável por licenças; T109 |
| G-08 Evidência de pintura | EFLButtonUI, EFLLabelUI, EFLPanelUI, EFLTextFieldUI, EFLToggleButtonUI e EFLComboBoxUI ausentes; PasswordField/Dialog sem referência visual suficiente | Mantenedor EFL/homologador; T119 |
| G-09 Paridade por superfície | Automatizados não equivalem a comparações reais Windows DPI 100/150/200% e PowerPoint 2016/365; nenhuma nova homologação executada nesta revisão | Equipe editor/homologação/usuário; T120–T122 |

### D-10 — Biblioteca estática de skins, com adaptador de recursos

**Contexto/decisão**: converter o acervo durante desenvolvimento em dados puros versionados (`internal/skin`), com recursos em `internal/adapter/legacyskin`; runtime consome definições estáticas/bytes resolvidos por adaptador. Usar biblioteca padrão no conversor quando suficiente; parser TTF apropriado com versão/licença fixadas se necessário. Nada de executar Java, interpretar XML legado a cada abertura, plugins ou registry central que importe Fyne/PPT no domínio.

**Alternativas/custo**: copiar manualmente sete estilos não cobre o acervo nem permite regeneração auditável; portar Swing/JVM viola distribuição. Geração determinística exige verificar diff/hash/licença e manter adapters, mas permite reproduzir e incorporar recursos offline. Revisar se a fonte do acervo mudar; não substituir imagem por nome parecido.

### D-11 — Fonte declarada não é fonte efetivamente pintada

**Evidência**: SkinBuilderImpl.createState calcula `style` por bold/italic, mas chama FontLoader.getFont(type,size) sem aplicar esse valor ao Font. Registrar flags XML e essa transformação separadamente. Não ativar negrito Fyne para "corrigir" o legado sem decisão visual. Fontes reais podem já ser faces pesadas/itálicas; sintese e flags não são equivalentes.

**Disponibilidade**: [arial.ttf](../../material-legado/elemetos_componente/xml/fonts/arial.ttf) original e, em [outras_fontes](../../material-legado/elemetos_componente/outras_fontes), futurbci.ttf, futura_mh.ttf, futurah.ttf, futuram.ttf e veramobd.ttf cobrem os seis nomes referenciados; futurabi.ttf e futurak.ttf são adicionais, sem substituição automática. Presença não comprova metadados/licença/cmap/OS2/métricas nem embedding. Fyne FontSource por controle é a estratégia; TextStyle sintético precisa de prova, não de suposição.

**Trade-off/revisão**: fontes embutidas custam inventário e licença, mas evitam instalação obrigatória no editor; PPT precisa de incorporação permitida ou ambiente explicitado. Referência visual homologada resolve divergências do builder sem perpetuar defeitos por acidente. Revisar quando face/metadados ou evidência de execução estiverem disponíveis.

### D-12 — Evolução semântica e compatibilidade intencionais

Schema 3/catalog 2 acrescentam famílias, aparência/overrides e perfil com recursos. Reader aceita v1/v2/v3 explicitamente, writer v3; editor antigo recusa v3. Não trocar discriminadores existentes. Duas legendas/subclasses Report são variantes; UF/Frame/presets informativos são composições, não Kind por nome XML. ParentID opcional mantém lista plana/bounds absolutos com regras de contêiner. Alternativa de JSON livre ou migração silenciosa perderia validação/dados; custo aceito é codec/cópia/validação completos e backup antes do primeiro save migrado.

Aplicar aparência nova em arquivos antigos foi escolha explícita do usuário, preservando overrides/IDs/bounds/conteúdo. Não redimensionar teclado antigo; novo default 200 × 350 só para novas inserções. Corrigir DESIGN 007: 110 × 60 é default atual; 160 × 44 era divergência documental, não estado a migrar. Histórico v1 permanece identificado nos artefatos. Revisar contrato antes de qualquer próximo formato.

### D-13 — Cena única, estados locais e PPT editável

Resolver tipografia/layout/pintura em presentation e compartilhar com desktop interativo/PPT. Eventos e hit regions não substituem a cena por widgets Fyne opacos. Fontes/imagens em cache por hash; sem recarga por frame. PPT mantém texto/formas, imagem só decorativa e snapshot do recorte visível; incapacidade de representar um alvo é diagnóstico antes de publicar, não rasterização. Alternativa de três renderers independentes perpetua a divergência atual; custo é ampliar primitivas, métricas e widgets custom. Revisar diferenças somente com fixture/referência por ambiente.

### D-14 — Entregas parciais não encerram fidelidade total

Piloto ReportEvent prova o slice arquitetural; não aprova todas as subclasses. Separar importador/fontes/skin, domínios/persistência, apresentação, desktop/prévia e PPT, testes automatizados, aquisição de referência, manuais e aceite. Bloqueios mantêm entrada e tarefa abertas. G-03/G-04/`JAVA_PROFILE_UNVERIFIED` continuam independentes e Java pausado. Não exigir `.page`, integrações/periféricos ou reativação Java para este incremento. Revisar conclusão apenas quando todas as entradas estiverem comprovadas e o usuário aprovar a revisão executável testada.

### D-15 — Pacotes adicionais, procedência e conflitos de classpath

**Evidência/decisão (2026-10-19):** os 18 pacotes fornecidos em `outros_materiais` declaram `src/images` como raiz de recursos. O importador usa o caminho relativo a essa raiz para reproduzir `/fonts`, `/icon`, `/azul` e demais referências, sem heurística por basename. As 271 cópias em `bin` conferem byte a byte; importar só fontes/imagens originais, nunca scripts Ant, metadados IDE ou thumbnails. Todas as origens equivalentes são preservadas, com deduplicação de bytes por SHA-256.

Um conflito real foi encontrado em `/azul/label_lot_azul.png`: XML/comum claro versus Milionária escuro. O usuário escolheu **manter a versão clara original**. legacy_resource_choices.json (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/scripts/legacy_resource_choices.json`) fixa origem/hash/justificativa; o manifesto conserva ambas as candidatas. Novos conflitos ficam bloqueados; pin ausente ou alterado produz diagnóstico, nunca fallback. Essa decisão não autoriza redistribuição nem aprova fidelidade visual.

Inventário vigente: 214 definições, 8 fontes, 270 recursos, 47 diagnósticos; 16 referências a 16 caminhos continuam ausentes. Versão e evidências técnicas/build estão no [quickstart](quickstart.md#7-incremento-do-acervo-javaefl--2026-10-19-evidência-vigente). Não confundir a redução das faltas com entrega dos delegates ou encerramento de G-06–G-08.

### D-16 — Snapshot de metadados por documento e resolução de bytes exatos

**Implementado na continuação:** `application.Document.AppearanceProfile *skin.Profile`
isola o perfil do projeto em edição sem acrescentar responsabilidades de outros
domínios a `project.Project`. `Profile` schema 1 conserva definições, fontes,
aliases, recursos, diagnósticos e até 16 revisões arquivadas diretas. Hash SHA-256
do JSON determinístico exclui somente o hash próprio, incluindo os arquivos e seus
hashes; selar/validar não comprova licença. `Library` mantém cópias privadas imutáveis
e revisão/estado exatos; desconhecidos falham em vez de usar a revisão corrente.

`visualprofile.Capture` conserva todos os metadados da biblioteca original em cópias
privadas: definições, fontes, aliases, recursos e todas as revisões arquivadas,
inclusive versões não usadas e todos os estados. Não enriquece perfis privados com
builtin nem faz upgrade. O uso corrente limita apenas diagnósticos e acréscimo de
entrada missing. Mantém recursos `referenced/unknown` e bytes já
incorporados; **não acrescenta binários**, mesmo com `allowed` no JSON. `ResourceData`
confere AssetID/hash/tamanho/MIME para embedded e só permite builtin com bytes
exatamente coincidentes para referenced. Não resolve pelo nome da família, homônimo,
fonte instalada ou caminho externo. Desktop/PPT usam o perfil e cache SHA/face exata.

**Alternativas rejeitadas:** biblioteca global mutável ou fallback de revisão
quebrariam reabertura estável; copiar fontes baseado só no JSON confundiria intenção
com autorização. Custo aceito: recursos não incorporados continuam dependentes de
bytes exatos disponíveis e geram diagnóstico de portabilidade incompleta. Restam
16 recursos ausentes, 8 licenças desconhecidas e 7 fontes com embedding restrito.

**Desvio arquitetural justificado:** a poda anterior causava bug crítico
save/reopen → `AddComponent` com `defaultconfirm`. Snapshot integral de metadados
preserva editabilidade futura (inserir famílias/trocar variantes após reopen),
todos os estados e namespace customizado. Não é incorporação integral binária;
o alvo normativo SUBSET de bytes usados e permitidos permanece. Payload selado
mais amplo altera o hash. Regressões save/reopen/addButton/variant/transient states/
custom namespace incluídas na compilação pós-correção integral de Capture pelo main,
sem execução ou PASS de testes Go. Hash do perfil não se confunde com a versão
inalterada do inventário ou com o novo hash do EXE registrado no quickstart.

**Transição, não redução do alvo D-12:** reader schema 3 aceita perfil ausente;
v1/v2 não admitem o campo. DTOs atuais `Fonts/Aliases/Resources` e `Style`
(`SkinName/SkinVersion`, RGB, `FillAlpha` opcional) ainda não são
`fontFaces/resourceIds`, `component.appearance`, overrides com presença completa e
RGBA canônico. Família desconhecida fica literal em Style com diagnóstico, sem
descritor unresolved implementado. Manter esses alvos e tarefas abertos; revisar
compatibilidade antes de mudar o wire, não declarar convergência por renomear docs.

**Evidência e revisão:** galerias estáticas das 16 famílias baseadas em Java/XML
não são screenshots de execução Java ou homologação Windows/DPI. Conforme main após
correção integral de Capture: gofmt/vet PASS, **15 suítes Go compiladas, não executadas**,
bloqueio prévio por política Windows preservado. Python **28 PASS** (15 importador +
4 packaging + 9 Java 6, incluindo erro esperado de `String.strip`, não homologação EFL);
importer `--check` PASS, mesmo inventário 214 definições/8 fontes/270 recursos/47
diagnósticos. Novo EXE Windows amd64 construído pela tarefa VS Code, hash/tamanho no
[quickstart](quickstart.md#novo-exe-pós-capture--pacote-e-homologação-pendentes).
Distribuição Windows gerada e verificada; checksum do ZIP no sidecar externo
`.zip.sha256`. Reempacotar documentos preserva o EXE. Conferência estática de 4
estados da galeria com XML e 14 geometrias de teclas com Java original PASS,
não homologação visual. Embedding PPT, Windows/DPI/PowerPoint e aceite continuam
pendentes, Java/EFL desabilitado.

<a id="pesquisa-codigos-barras-2026-10-23"></a>
## Pesquisa e decisões — códigos de barras (2026-10-23)

**Origem**: classes transcritas pelo usuário e solicitação de componente processado parametrizável; reforço explícito de GitHub Spec Kit e divisão por modelo de domínio. Não chegaram imagens nem implementações dos três módulos. Esta pesquisa documenta decisões, não testes do produto ou homologação bancária.

### D-17 — Dois componentes, domínios existentes

**Decisão**: 031 detalha entrada e 032 detalha apresentação processada; Catálogo é dono dos tipos/regras, Tela das instâncias, Projeto do formato, Jornada das rotas explícitas, Demonstração dos quadros e a coordenadora do runtime/adaptadores. Backlogs e modelos únicos nos respectivos domínios, conforme constituição I/VII.
**Justificativa**: permite reutilizar regra pura sem dependência de Fyne/PPT/IO e evita acumular persistência/processamento no componente visual.
**Alternativas rejeitadas**: novo domínio financeiro ou serviço externo, uma spec monolítica com todo o código, e repetição das mesmas tarefas em 031/032 e nos domínios.
**Custo/revisão**: sincronizar referências qualificadas entre artefatos; expandir domínio somente se surgir negócio real aprovado, não por existir um novo componente.

### D-18 — Contrato matemático explícito, não transcrição dos defeitos

**Decisão**: aceitar somente 44/47 dígitos, regras/vetores no [contrato](contracts/component-catalog.md#contrato-codigos-barras-2026-10-23). Mod10 e Mod11 de boleto usam premissas usuais; Mod11 SISPL é experimental até comparação com classe real. Não copiar setter48/recorte47, validateCampo5 usando campo1, cache de conversão incoerente ou teste de fator apenas pelo primeiro algarismo.
**Justificativa**: nomes de classes não provam suas implementações. Centavos exatos e zeros textuais evitam perda de precisão/identidade; fator numérico não resolve o ciclo de vencimento (reinício em 2025).
**Alternativas rejeitadas**: inventar fonte da classe, dizer que SIGCB já foi homologado, calcular data pela data atual ou aceitar 48 descartando um dígito.
**Custo/gate**: solicitar código/vetores independentes antes de anunciar equivalência SISPL. Verificação isolada de F01/F02 na redação apenas conferiu aritmética/conversão sob essas premissas; F02 é sintética e não um oráculo independente.

### D-19 — Parametrização processada e captura sem inferências

**Decisão**: combos por linha no inspector escolhem CPF/CNPJ/Nome/Nome fantasia/Razão social; valores textuais configurados permanecem literais e independentes. Não obter nome/documento da entrada ou chamar serviço. Modo de leitura da entrada é explícito, dispositivo tipo teclado com terminador Enter e colagem integral, sem heurística temporal.
**Justificativa**: o pedido informa os rótulos/opções, não uma consulta bancária ou SDK de QR. A leitura de 44 dígitos não fornece os dados pessoais exibidos pela ficha.
**Alternativas rejeitadas**: combo global sobrescrevendo blocos, popup de autoria no PPT, CPF/CNPJ validado sem solicitação, input truncado no primeiro campo e detecção de leitora pela rapidez das teclas.
**Custo/revisão**: a localização dos combos é premissa explícita de UX a confrontar com anexos quando recebidos. Novo requisito de combo runtime, PIX/câmera ou integração com ficha exige contrato e tarefas próprios.

### D-20 — Compatibilidade e evidência antes da oferta

**Decisão**: catálogo 3 planejado para as duas novas famílias; envelope schema3 preservado, reader dos pares anteriores e primeiro save incompatível com aviso/cópia original. Nenhuma versão do executável ou perfil é alterada por este documento. 031 usa medidas do Java transcrito; 032 tem proposta visual identificada, não referência observada.
**Justificativa**: leitores estritos de catálogo 2 não reconhecem os novos tipos; aumentar catálogo é mais preciso que afirmar compatibilidade ou mudar schema sem necessidade. Ausência de imagem impede alegar fidelidade, não redigir cenários testáveis.
**Alternativas rejeitadas**: gravar os novos tipos como Report/catalog2, fallback silencioso, atualizar perfis arquivados, inventar campos bancários ou screenshot do anexo ausente.
**Custo/gates**: testes de compatibilidade, revisão dos anexos/skin/fontes/DPI e PPT real. Roteiro typeText dos compostos não contratado nesta entrega recebe diagnóstico; representação estática é suportável sem prometer interação. Nenhum novo build, teste Go ou aprovação manual decorre desta etapa documental.

<a id="pesquisa-menu-aplicacao-2026-10-23"></a>
## Pesquisa e decisões — Menu aplicação e autorização conjunta

### D-21 — Menu único, árvore tipada e fronteiras existentes

**Decisão**: 033 contém spec/DESIGN/checklist; plano/modelo/tarefas nos donos001–006. ApplicationMenu e SubmenuProperties são modelos separados no Catálogo, instância única em Tela, wire no adaptador Projeto, destinos em Jornada e quadros em Demonstração. API acordada em [003/data-model](../003-catalogo/data-model.md#modelo-menu-aplicacao-2026-10-23), sem domínio financeiro. Raiz/submenu são dois estados visuais locais, aninhamento usa caminho da sessão.
**Justificativa**: preserva árvore autorável sem centenas de widgets/componentes persistidos ou telas artificiais. Limites16/1000 e UUIDs globais tornam validação/cópia/persistência limitadas e verificáveis. Factory recebe callback supplyID e cria Enabled=true; reader conserva false explícito.
**Alternativas rejeitadas**: Menu como Report renomeado, uma tela por submenu, árvore em JSON livre, singleton mutável de UI, cópia de regra no exportador, execução Java/XML ou despacho financeiro pelos metadados.
**Custo/condição de revisão**: travessia de referências e remapeio profundo na duplicação; manter diagnóstico localizado por item. Duplicação de componente conserva destinos; de tela remapeia só autorreferências, como política explícita deste desenho, a confirmar no teste manual.

### D-22 — Fonte do seed, layout e estado inicial fixo

**Decisão**: um seed serviços/esquerda+jogos/direita proveniente do Java/XML fornecido, com submenus aninhados e Bolão digital. Layouts root/standard/financial/bolao; financeiro140x80 e Bolão três botões+legendas são medidas/composição textuais fornecidas, não imagem observada. Registrar extração de outras geometrias aqui quando fonte for conferida pelo executor. Árvore XML enorme não é repetida nos documentos. InitialMenuID persistido escolhe raiz/submenu para canvas/quadro/PPT; caminho runtime deriva ancestrais mas não é salvo.
**Justificativa**: distingue autoria reprodutível de exploração efêmera; fonte legada informa aparência, não execução. Imagens anexadas estão realmente vazias; esta etapa não consegue fazer comparação visual.
**Alternativas rejeitadas**: inventar medidas/ícones por screenshot inexistente, duplicar XML integral em várias specs, chamar seed genérico de fiel ou capturar estado da prévia para exportar.
**Custo/gate**: executor deve ter acesso à transcrição original para reproduzir seed/geometrias; se não tiver, pedir a fonte. Falta de imagem limita alegação de fidelidade, não proíbe implementar simulação autorizada. Layout Bolão tem referência inicial de três botões; alterações não podem podar itens silenciosamente e exigem revisão de overflow/capacidade.

### D-23 — Rotas, demonstração e segurança da simulação

**Decisão**: TargetScreenID explicitamente autorado é a única rota de folha; validação/impacto em Jornada sem Action duplicada. Voltar/Cancelar/Sair são comandos locais; Sair fecha apenas prévia. Sem destino, informar seleção simulada. Prévia Tela Atual mantém menus mas ignora transição entre telas. Movement*/Shortcut são dados inertes. Demonstração conserva click visual e estado inicial fixo; typeText nos três compostos não é contratado.
**Justificativa**: reproduzir organização visual não autoriza negócio, autenticação, timers, periféricos, rede ou QR. Clique visual não deve causar efeitos ocultos ou um formato persistido novo por acidente.
**Alternativas rejeitadas**: inferir tela pelo texto/código, logoff real, usar bindings paralelos, preencher ficha automaticamente ou navegar no PPT por programação do legado.
**Custo/decisão restante**: coordenação principal decide se estenderá click por item metadata; antes de habilitar exige sincronizar Step/wire/capacidades/quadros e validar compatibilidade. Baseline sem essa extensão é completo para estado fixo e prévia livre, não bloqueia a redação.

### D-24 — Catálogo3 conjunto e autorização não é homologação

**Decisão**: complementar D-20 para incluir 033 no mesmo alvo schema3/catalog3 de031/032, readers1/1 2/1 3/2 3/3 e backup antes de save incompatível. Não é catálogo4, versão de produto ou schema de perfil novo. Usuário autorizou implementar os três incrementos e fará teste manual depois; preservar ponteiro031 e todos os checklists históricos/customizados.
**Justificativa**: catálogo2 não reconhece tipos novos, mas a extensão ainda coordenada do lote não exige três migrações. Autorização permite prosseguir no escopo aprovado apesar de itens históricos abertos, sem fingir revisão/aceite realizado.
**Alternativas rejeitadas**: editar checklist de revisor para contornar gate, fechar tarefas técnicas em lote por documentos, declarar homologação por build ou reutilizar hashes/observações de revisão anterior.
**Custo/gates**: novos T140–T144 e tarefas de domínio permanecem abertos até validação principal; seção18 do quickstart recebe revisão e evidências futuras. Esta etapa não alterou Go, executou builds/testes/aplicativo ou resolveu gates de imagens, módulos reais SISPL, DPI, PowerPoint, licença e aceite humano.

## Reaplicação do estilo ao código atual — 2026-10-24

Inspeção encontrou duas causas do visual genérico: a cena só consultava skins explicitamente selecionadas, e a prévia substituía a cena por Button/Entry padrão Fyne. O renderer desktop também não aplicava os bytes de `Primitive.FontFamily`; o nome inferido `Futura MH` divergia da família real `Futura Medium`. O parser antigo perdia preto sem alpha e tratava branco alpha=0 como opaco.

Fontes primárias usadas: `elemetos_componente/xml/skin.xml`, `component/EFLVirtualKeyboardNumeric.java`, `component/EFLReportEvent.java`, `builder/impl/SkinBuilderImpl.java` e `SkinUtil.java`; React é referência secundária. O builder calcula flags XML, mas instancia a fonte física sem aplicar o style calculado. Por isso o importador conserva ambas as informações e o renderer usa o arquivo físico. O menu declara `spacement=0.8`; o texto é medido com essa face, com quebra por palavras e distância entre linhas correspondente. A validação final de métricas Java/Office permanece pendente.

O importador determinístico produziu 205 chaves tipadas operacionais e oito definições indisponíveis (os mesmos oito nomes de FR-058), usando 224 recursos. Os recursos complementares já presentes em `internal/skin/assets` foram preservados; não se inventaram imagens para os ausentes. Não equiparar esse resultado ao inventário histórico de perfil completo de 206 definições/270 recursos: a reconstrução atual não recompõe todos os contratos e metadados daquela expansão.

Regressões de botão sem skin, teclado antigo de 12 teclas/default 280 e relatório genérico falharam antes da correção. Após a implementação passaram as suítes Go e os casos de cursor/TAB/Backspace/runtime, imagem original no PPT e bloqueio de recurso indisponível. Captura Fyne em memória e build Windows têm evidências próprias na seção 22 do quickstart; não representam homologação humana.

## Fonte integral de menus recuperada — 2026-10-24 / 1.0.0.2

Localizado `D:/desenvolvimento/ia/estudo/sistema-prototipo-android/componentes/src/java/menu/menu.xml`, 30.514 bytes, SHA-256 `f54ad100a67ec1ec947ee74d35123a908e0fc06004813c82f1d541317dee6bd8`. Cópia exata preservada em `material-legado/menu/menu.xml`. A declaração ISO-8859-1 explica os acentos; ler como UTF-8 corromperia nomes. A árvore confirma 12 serviços/10 jogos e 27 títulos descendentes/148 folhas em RELATÓRIOS, mas contradiz a transcrição anterior: LotoGol comentado, Milionária ativa; quatro jogos com ramos normal/especial; nomes/códigos reais do digital e serviços completos.

Decisão: gerar literal Go de desenvolvimento com ordem/atributos ativos, sem parser runtime, sem transcrever contagens sintéticas e sem migrar projetos. A fábrica é pura e cria novos UUIDs. `tituloPrincipal` de folha permanece somente na referência estática porque o modelo não oferece título de folha; títulos de submenu são aplicados. Legendas digitais visuais existentes são complemento e não atributos XML. O XML integral prevalece sobre a antiga restrição de não preencher ramos cujo conteúdo não estava disponível.

Para aparência, `skin.xml` e `ButtonSkin.java`/`EFLButton.java` confirmam as chaves por jogo e estados touched/disabled. Os PNGs originais contêm marca e ícone, logo não cabe sobrepor o mesmo nome. Variant explícita é estável para novos itens; compatibilidade de nome conhecido só resolve aparência de antigos, sem alterar navegação. Alternativas rejeitadas: construir menu a partir de assets; trocar todas as folhas por botao_menu; interpretar XML em runtime; reescrever documentos existentes. Revisar quando outra fonte XML autorizada for recebida.

## Revisão de PnlMenu e UIMenuControle — 1.0.0.3

Fontes copiadas sem transcodificação para material-legado/menu: PnlMenu.java SHA-256 `4f43113a004f4b31dc5024a621e026ff18eb75682901079370b87320ac2b91dc`; UIMenuControle.java SHA-256 `730d90b807233d8264eea0d7a1ab937d17130cebfddd9930c2dea3aa5cf687b7`. Origem: mesmo diretório componentes/src/java/br/gov/caixa/sispl/infra/ui do projeto sistema-prototipo-android, com PnlMenu em pnl.

PnlMenu linhas 218–248 define faixas; 251–419 define grades e iteração getItemMenuXml; 464–527 define comandos; 550–570 preserva nomeMovimento. UIMenuControle.selecionaCasoUso usa getInfoMovimento(nomeMovimento) quando o nome existe e não decide só por tipo/opção. A implementação Go anterior conferia árvore e runtime, mas seu teste de clique não inspecionava o conteúdo reconstruído. A nova regressão usa os controles reais de window.Content após cada redraw.

A busca em materiais, projeto legado, anexos e registros locais não recuperou movimento.xml/movimentos.xml, ItemMenuXml ou InfoMovimento completos. Não transcrever um cadastro fictício a partir dos códigos do menu. ItensMenu é a coleção de filhos observável no Java; o app preserva-a em Items. O incremento corrige apresentação/navegação e autoria dos metadados; a conferência do cadastro de movimentos permanece trabalho aberto, dependente da fonte real. Nenhuma operação financeira é executada.
