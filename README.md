# Jornada Studio — Gerador de Jornada

Editor desktop offline em **Go 1.27.2 / Fyne 2.8.2**, revisão **1.0.0.11**, para compor telas, conectar jornadas e preparar demonstrações. O projeto editável é salvo em `.jornada`; a saída PowerPoint mantém textos e formas editáveis. Java 1.6/EFL permanece visível e desabilitado.

O código fica em [apps/gerador-jornada](apps/gerador-jornada). O [guia da aplicação](apps/gerador-jornada/README.md) detalha operação, propriedades, atalhos e distribuição. As [36 specs](specs/001-gerador-jornada/spec.md) separam implementação existente, requisitos ainda pendentes e homologação manual.

## Executar a revisão atual

No Windows x64, abrir [JornadaStudio.exe](apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe). Não é necessário instalar Go, Node, Java ou Office para usar o editor. O diretório `dist` contém artefatos gerados, não o código-fonte canônico.

## Regra vigente de compilação e entrega

**Por solicitação do usuário, gerar por enquanto somente `JornadaStudio.exe` para Windows x64**, em `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`. Esta regra vale para manutenção, correções, agentes e automações de entrega.

ZIP, TAR.GZ, instaladores, bundles `.app`, builds Linux/macOS, outras arquiteturas ou qualquer outro formato de compilação/distribuição **somente devem ser gerados quando o usuário solicitar explicitamente o formato desejado**. Um pedido para corrigir, testar ou recompilar o app não autoriza essas saídas adicionais. Pacotes e hashes de entregas anteriores são históricos e não precisam ser regenerados.

Usar `scripts/check.ps1` para as verificações e `scripts/build.ps1 -OutputDirectory dist/windows-amd64` para gerar apenas o EXE. **Não usar `scripts/build.py` na compilação de rotina**, pois ele também empacota a distribuição. Recursos intermediários e arquivos temporários dos testes não constituem formatos adicionais de entrega. Esta regra não altera a funcionalidade de salvar projetos `.jornada` ou exportar PPT dentro do aplicativo.

```powershell
# A partir da raiz do repositório:
./apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe
```

Esta revisão entrega apenas o EXE em `dist`; as fixtures de revisões anteriores não estão presentes no checkout atual. Capturas do incremento 034: [launcher](apps/gerador-jornada/build/034-launcher-native.png) e [estrutura do projeto no Windows](apps/gerador-jornada/build/034-workspace-native.png). Capturas dos submenus anteriores permanecem como histórico. Arquivos em `build` são evidências locais e fixtures temporárias de teste, não formatos adicionais de distribuição.

## Inicialização com vídeo e áudio — 1.0.0.11

A splash usa o vídeo **video-splash-programa-v2.mp4** incorporado ao EXE, com áudio original. A splash abre centralizada em uma janela normal de 960×540 (unidades de interface), com barra de título e bordas, sem fullscreen nem maximização automática. Sua área de conteúdo mostra somente o vídeo: não há progresso, legenda de carga, autoria ou botão para pular. O vídeo é ajustado proporcionalmente à área interna da janela, sem contar bordas e barra de título nem cortar o quadro. Faixas pretas podem aparecer para preservar a proporção. Por solicitação do usuário após a falha relatada em fullscreen, foi removida a entrada/saída desse modo; o player passa a usar a janela normal durante toda a apresentação e transição.

Os recursos necessários carregam em paralelo: skins e integridade, imagens decodificadas, fontes, defaults dos componentes, fundo/ícone e preferências locais. Quando a splash é exibida, a tela seguinte só é liberada depois do **evento real de término do vídeo e do carregamento**. Uma carga concluída em poucos milissegundos não corta a apresentação. Se a carga terminar depois, o último quadro permanece até ficar pronta. O fade leva à seleção de projetos; `-open` segue para o editor respeitando a mesma preferência de apresentação. Utilitários de exportação/gravação sem interface não reproduzem vídeo.

Na primeira execução com esta configuração de usuário, o vídeo aparece uma vez. Na seleção de projetos, **Mostrar splash screen ao iniciar** fica **desmarcado por padrão**. As próximas aberturas vão direto para a tela inicial depois de preparar os recursos. Marcar repete o vídeo nas próximas inicializações; desmarcar volta a suprimi-lo. A alteração é salva imediatamente e sobrevive a reinicializações, atualizações do EXE e mudanças de projeto. Uma apresentação cancelada ou com falha não consome a primeira exibição.

Foi necessário criar a política pura `startup.Preference` e o adaptador `desktop/startup_preferences.go`, além do coordenador `desktop.StartDesktop`. A preferência fica em `%APPDATA%/JornadaStudio/startup-preferences.json` (`presented` e `show_on_startup`), com publicação atômica; `JORNADA_CONFIG_DIR` isola testes. Não fica no `.jornada` nem na pasta do EXE. A construção do launcher não grava nem marca a opção. Mesmo sem vídeo, os recursos obrigatórios são carregados e suas falhas impedem abrir a edição.

A autoria foi movida para o rodapé da seleção de projetos, substituindo a mensagem anterior:

**Desenvolvido por Rômulo Penha - Canais Parceiro - 03**

Foi necessário separar a implementação por responsabilidade: o modelo puro de eventos e liberação fica em [`internal/startup`](apps/gerador-jornada/internal/startup/session.go), documentado pela [spec 035](specs/035-inicializacao/spec.md). A [spec 036](specs/036-splash-screen/spec.md) e seus próprios modelo/plano tratam janela, mídia, cache e transição nos adaptadores `desktop`, `builtinassets` e `nativevideo`. Não há campo novo no projeto, serviço, navegador ou player externo.

A integração Windows usa uma janela de vídeo própria, com tratamento de pintura/tamanho e início após mídia pronta. A implementação anterior corrigiu a interpretação da razão 16:9 como tamanho em pixels; a ausência de imagem ainda foi relatada pelo usuário e motivou esta revisão em janela normal. Os eventos de reprodução sozinhos não provaram que havia imagem. Testes agora verificam capturas reais em aberturas consecutivas, além de eventos, áudio ativo e geometria. MP4 conserva seus bytes; a cópia de execução vai para o cache por hash. [Roteiro, capturas e evidências](specs/036-splash-screen/quickstart.md). A regra de entregar somente JornadaStudio.exe continua vigente.

## Correção ao criar e reabrir projetos — 1.0.0.5

Foi corrigido o encerramento na passagem do launcher para o editor, tanto após criar quanto ao escolher um recente. A gravação do `.jornada` terminava, mas o construtor do editor disparava `SetTheme` ao inicializar o seletor de tema. Com o driver Windows já rodando, o Fyne tentava atualizar a aparência nativa de uma janela ainda sem viewport (`setDarkMode → GetWin32Window`) e ocorria acesso nulo.

O tema agora é preparado antes de `NewWindow`; o seletor recebe seu valor inicial sem disparar a troca global, e reaplicações redundantes são evitadas. A cor do título já nasce no tema escolhido. A troca manual de tema continua funcionando. O log registra `editor_open_start` e `editor_open_ready` na transição, sem nomes ou caminhos de projetos.

A reprodução anterior capturou o panic; a versão corrigida percorreu oito fluxos com driver gráfico real: criar em branco e reabrir nos temas Claro/Escuro/Azul, além de criar com exemplo e reabrir. Há teste de regressão da ordem de inicialização e uma ferramenta local `scripts/launcher_smoke` com build tag própria. [Resultados e comandos](specs/034-launcher-estrutura-projeto/quickstart.md#correção-de-encerramento--1005). O arquivo criado pode ser reaberto com a versão corrigida.

## Launcher e estrutura de projeto — 1.0.0.4

Sem `-open`, o EXE agora apresenta a tela inicial: projetos recentes, **Novo projeto em branco**, **Novo projeto com exemplo** e **Abrir projeto existente**. Nome e pasta são confirmados antes da criação. O modo em branco contém apenas o fundo Loterias CAIXA; o exemplo permanece editável. Recentes ficam no JSON por usuário `%APPDATA%/JornadaStudio/recent-projects.json`; entradas removidas continuam indicadas como indisponíveis. `JORNADA_CONFIG_DIR` permite isolar essa preferência em testes.

Foi necessário separar a inicialização do editor da criação do exemplo, acrescentar o armazenamento atômico de recentes e usar a própria pasta do projeto como workspace. A aba **Estrutura projeto** abre no painel esquerdo e mostra uma árvore real com a pasta e suas `.jornada` da raiz. O menu **⋮** cria, importa e atualiza a lista. Importações são validadas antes da cópia, preservam os bytes da origem e recebem sufixo em caso de colisão. Não existe novo arquivo de workspace. Abas que não couberem no painel ficam no menu de transbordamento **…**.

Na **Prévia**, clicar numa folha sem destino do Menu aplicação oferece criar uma jornada com seu rótulo. Cancelar mantém arquivos e histórico intactos. Confirmar cria o arquivo, salva as alterações e o vínculo da origem e atualiza a árvore; a próxima ativação abre a tela inicial do destino. **Voltar à jornada anterior** restaura a origem e o submenu visitado. A jornada de origem precisa estar salva na pasta; a Prévia Tela Atual conserva seu escopo sem navegar. Nome, código de movimento e atalho nunca são usados para inferir destino.

Para persistir esse vínculo, foi adicionado `MenuItem.TargetJourney` e o catálogo passou de **3 para 4**, com schema **3** preservado. O leitor aceita os pares históricos e insere somente destino vazio na migração. O primeiro salvamento incompatível avisa e guarda o original em `.bak`; se já houver backup diferente, usa `.catalog-<versão>.bak`, sem substituir cópias anteriores. Falha de backup impede a sobrescrita. O vínculo participa de desfazer/refazer; desfazer o vínculo conserva o arquivo criado, para não apagar trabalho.

A [spec 034](specs/034-launcher-estrutura-projeto/spec.md) e seu [quickstart](specs/034-launcher-estrutura-projeto/quickstart.md) registram C1/C2/C3, regressões e limites. Capturas do EXE: [launcher](apps/gerador-jornada/build/034-launcher-native.png) e [árvore do projeto](apps/gerador-jornada/build/034-workspace-native.png). Testes técnicos não representam aceite humano ou homologação de DPI/Office. A entrega continua gerando somente `JornadaStudio.exe`.

## O que foi necessário para recuperar o estilo do legado

O código já criava os componentes, mas o estilo ainda não atravessava todo o caminho de renderização. Sem uma skin explícita, a cena usava cores genéricas; a prévia trocava o desenho por controles padrão Fyne; informar o nome de uma fonte também não carregava seus bytes no desktop. Foi preciso corrigir esses pontos no código, além de revisar os documentos.

| Mudança necessária | Aplicação concreta |
| --- | --- |
| Levantar a referência correta | Usar `material-legado/elemetos_componente/xml/skin.xml`, imagens/fontes disponíveis e classes Java, incluindo `EFLVirtualKeyboardNumeric`, `EFLReportEvent`, `SkinBuilderImpl` e `SkinUtil`. React é referência secundária; anexos ausentes não viraram evidência inventada. |
| Importar durante o desenvolvimento | Criar [import_legacy_skin.py](apps/gerador-jornada/scripts/import_legacy_skin.py), que gera `internal/skin/assets/library.json` e incorpora os recursos disponíveis. O aplicativo carrega JSON estático, sem interpretar XML Java ou executar Java em runtime. |
| Preservar cores e transparência | Distinguir alpha ausente de alpha explicitamente zero. Preto sem alpha continua preto opaco; fundos transparentes não se tornam brancos. O modelo Style continua RGB; alpha importado é aplicado pelas primitivas da cena. |
| Usar a fonte física | Ler família e flags físicas SFNT e aplicar os bytes com `canvas.Text.FontSource`. `futura_mh.ttf` declara a família real **Futura Medium**. Flags XML de peso/itálico são separadas da face física usada pelo builder Java. |
| Resolver defaults por componente | Aplicar skins do legado mesmo com Style.Skin vazio, preservar overrides explícitos e mapear confirmar/cancelar/sim/não para suas definições. A opção do inspector passou a ser **Padrão do legado**. |
| Compartilhar a cena | Centralizar desenho em [presentation](apps/gerador-jornada/internal/presentation/scene.go) e no [renderer Fyne comum](apps/gerador-jornada/internal/adapter/desktop/scene_render.go), consumidos pelo canvas e pelas prévias. O exportador PPT consome as mesmas primitivas. |
| Manter o estilo durante a interação | Criar [legacyControl](apps/gerador-jornada/internal/adapter/desktop/legacy_control.go) para botões/opções/teclas e [legacyEntry](apps/gerador-jornada/internal/adapter/desktop/legacy_entry.go) para texto. O Entry conserva cursor, seleção e clipboard com fundo/borda transparentes sobre a decoração EFL. |
| Recompor os componentes compostos | Separar as partes de cabeçalho/corpo/rodapé do relatório e aplicar fontes/decorações às fichas e menus. Imagens são recursos decorativos; conteúdo autorado permanece texto/formas editáveis. |
| Corrigir o teclado numérico | Substituir a grade anterior de 12 teclas pelo arranjo Java de 14 teclas: 7–9, 4–6, 1–3, < / 0 / >, TAB e LIMPAR. Novo default 200 × 350; instâncias antigas conservam seus bounds. LIMPAR emite Backspace, sem apagar todo o campo por causa da legenda. |
| Preservar tipografia/layout | Medir a face incorporada para quebra por palavras; transportar tamanho, cor, alinhamento, centralização vertical e `spacement` entre linhas. No menu, o fator original é 0,8. A equivalência final com Java/Office ainda requer teste visual. |
| Expor personalização coerente | Oferecer skin compatível, família incorporada, tamanho/cor do texto, fundo, borda e alinhamento no inspector. Salvar confirma a transação; Descartar mantém os valores confirmados. |
| Impedir substituições silenciosas | Filtrar as oito definições com recursos ausentes, mostrar diagnóstico de skin/fonte não resolvida e impedir publicação do PPT nesse caso. Integridade dos recursos incorporados é verificada por tamanho e SHA-256. |
| Verificar antes de entregar | Criar regressões que falharam antes da correção, exercitar cursor/TAB/Backspace/isolamento, validar PPT ZIP/XML e gerar build/captura Windows identificados. |

**Alterar `DESIGN.md` registra a intenção visual; não modifica o EXE.** Toda mudança de token precisa chegar ao importador, ao resolvedor ou ao renderer correspondente, ser validada e recompilada. O tema Claro/Escuro/Azul é do editor; a aparência da jornada vem do documento e de suas skins.

## Correção dos jogos e da árvore do menu — 1.0.0.2

`menu.xml` define itens, ordem e metadados; `skin.xml` define imagem, fonte e estados. A existência de um PNG não acrescenta um jogo ao menu. Foi necessário recuperar a fonte completa, respeitar ISO-8859-1, espaços e itens comentados, gerar [dados Go estáticos](apps/gerador-jornada/internal/catalog/menu_legacy_generated.go) e remover os relatórios sintéticos. [import_legacy_menu.py](apps/gerador-jornada/scripts/import_legacy_menu.py) faz a geração durante o desenvolvimento; `--check` confere a reprodução. O aplicativo continua offline, sem ler XML em runtime.

Cada jogo usa sua chave explícita: por exemplo, `botao_megasena`, `botao_milionaria` e `botao_duplasena`. Os caminhos de cada estado vêm do manifesto gerado do XML, inclusive a exceção `botao_duplasena` → `botao_dupla_sena.png`. A prévia troca para `touched` enquanto o mouse está pressionado, restaura `default` ao soltar/sair e usa `disabled` quando o item ou o composto está desabilitado. Canvas e PPT usam o estado estático confirmado.

Os PNGs já contêm ícone e nome da marca: esse nome não é desenhado novamente por cima. Uma legenda personalizada continua editável em caixa separada. `Variant` já existe no formato 3/3 e é preservada após renomear. Projetos anteriores com variante vazia reconhecem os nomes conhecidos somente para aparência, sem reescrever árvore, IDs, destinos ou metadados. **Para receber a árvore XML corrigida, inserir um novo Menu aplicação**; arquivos editados não são migrados automaticamente.

Evidências e teste manual da revisão estão na [seção 24 do quickstart](specs/001-gerador-jornada/quickstart.md#24-menu-xml-e-skins-dos-jogos--1002). A cópia recuperada do XML tem SHA-256 `f54ad100a67ec1ec947ee74d35123a908e0fc06004813c82f1d541317dee6bd8`.

## Revisão dos submenus e movimentos — 1.0.0.3

A revisão anterior conferia a árvore, mas ainda tratava o submenu dos jogos como grade genérica e inventava um título principal quando o XML não o declarava. Foi necessário preservar e conferir também [PnlMenu.java](material-legado/menu/PnlMenu.java) e [UIMenuControle.java](material-legado/menu/UIMenuControle.java). Mega-Sena agora apresenta a faixa verde original, as opções normal/especial centralizadas e Cancelar/Voltar no rodapé. Botões dos jogos usam 164 × 52; serviços da raiz continuam 164 × 68. Grades e legendas digitais seguem as posições extraídas desse Java, e a mesma geometria define pintura e clique.

`itensMenu` corresponde à coleção de filhos percorrida pelo legado (`getItemMenuXml()`), representada por `Items` no app. Ramo abre submenu; folha seleciona um movimento ou segue o destino explicitamente autorado. Mega-Sena normal/especial compartilham `702`, mas usam `nomeMovimento` distintos. Esses valores agora ficam visíveis na simulação e podem ser editados em Propriedades: código, nome, opção e atalho. Os metadados não executam Java ou sistemas externos.

**Configuração de movimentos ainda não verificada:** `movimento.xml`/`movimentos.xml` não foi recuperado dos arquivos, anexos ou registros locais acessíveis nesta sessão. `UIMenuControle` consulta esse cadastro por nome quando `nomeMovimento` existe. Preservar as referências do menu não equivale a conferir esse cadastro; não se afirma configuração integral nem se inventam classes, telas ou destinos ausentes. Novas instâncias preservam `tituloPrincipal` vazio em vez de copiar o nome. Títulos autorados de arquivos existentes permanecem preservados.

Teste novo percorre os controles da janela reconstruída, as quatro duplas normal/especial, a seleção distinta dos movimentos, Voltar/Cancelar/Sair e a confirmação dos metadados pelo formulário. Evidências da revisão e pendências na [seção 25 do quickstart](specs/001-gerador-jornada/quickstart.md#25-submenus-java-e-movimentos--1003). A entrega desta revisão gera **somente `JornadaStudio.exe`**, conforme a regra acima.

## Estado real e limites

| Identidade | Valor na revisão atual |
| --- | --- |
| Software | 1.0.0.11 |
| Arquivo `.jornada` gravado | schemaVersion **3** / catalogVersion **4** |
| Pares reconhecidos pelo reader | 1/1, 2/1, 3/2, 3/3 e 3/4, com os campos suportados; campos desconhecidos são recusados |
| Componentes inseríveis | Sete famílias base e três compostos: entrada de código de barras, ficha processada e menu |
| Biblioteca visual operacional | 205 chaves tipo/nome, 224 recursos referenciados, seis fontes físicas; oito definições indisponíveis fora da oferta |
| Entrega compilada nesta revisão | Windows amd64; Linux/macOS permanecem configuração futura |

O número do catálogo não comprova entrega da expansão para 16 famílias. Os tipos adicionais, variantes completas de relatórios, máscaras incrementais, aparência tipada/RGBA e perfil visual portátil por documento continuam no backlog. **Não há `Document.AppearanceProfile`, `skin.Profile` nem adaptador `visualprofile` na árvore atual**; suas descrições anteriores são histórico/plano de expansão.

A fábrica do menu agora deriva dos bytes originais de [menu.xml](material-legado/menu/menu.xml), recuperados de `D:/desenvolvimento/ia/estudo/sistema-prototipo-android/componentes/src/java/menu/menu.xml`. Os 27 títulos descendentes e 148 folhas de RELATÓRIOS são os reais da fonte, substituindo o preenchimento sintético anterior. Os dez jogos ativos incluem **MILIONÁRIA**; LotoGol está comentado no XML. Mega-Sena, Lotofácil, Quina e Dupla-Sena conservam os submenus normal/especial. Nenhum componente executa movimentos financeiros. A ficha processada ainda não tem imagem independente para homologação.

Fontes físicas acompanham o editor, mas **não são incorporadas ao PPT**. Disponibilidade/licença das fontes no computador do Office, fidelidade visual Windows/DPI 100/150/200% e aprovação humana continuam pendentes. O [quickstart, seção 22](specs/001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24) contém hashes, comandos executados, capturas e roteiro manual.

## Reproduzir a importação e o build

A partir de `apps/gerador-jornada`, com Python/PowerShell e as ferramentas de build:

```powershell
# Preparação local de Go/MinGW, quando ainda não disponíveis:
python -X utf8 scripts/bootstrap.py

# Regenerar os dados e incorporar os recursos disponíveis do acervo:
python -X utf8 scripts/import_legacy_skin.py

# Conferir a árvore e as skins, sem regravar:
python -X utf8 scripts/import_legacy_menu.py --check

# Conferir o manifesto contra o XML atual e os bytes incorporados:
python -X utf8 scripts/import_legacy_skin.py --check

# Conferir formatação, vet e testes Go/Python:
./scripts/check.ps1

# Gerar somente JornadaStudio.exe, sem empacotar:
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
```

O importador não substitui auditoria de proveniência, direitos ou cobertura integral do acervo. Arquivos complementares já incorporados são reutilizados; recursos ausentes não são inventados. Na revisão 1.0.0.2 passaram testes Go nos 13 pacotes com testes, `go vet`, gofmt, três testes Python de empacotamento, importador `--check`, build Windows e smoke nativo. Essas evidências são técnicas e não equivalem a homologação no PowerPoint ou aceite humano.

## Specs e manutenção

- [Spec coordenadora](specs/001-gerador-jornada/spec.md): integração e estado do produto; os demais domínios mantêm seus próprios contratos e backlogs.
- [Catálogo](specs/003-catalogo/spec.md), [Tela](specs/004-tela/spec.md) e [Propriedades](specs/028-propriedades/spec.md): famílias, instâncias e personalização.
- [Contrato `.jornada`](specs/001-gerador-jornada/contracts/project-format.md) e [contrato PPT](specs/001-gerador-jornada/contracts/pptx-export.md).
- [Research](specs/001-gerador-jornada/research.md): fontes, decisões, histórico e divergências conhecidas.
- [Tasks](specs/001-gerador-jornada/tasks.md) e [quickstart](specs/001-gerador-jornada/quickstart.md): trabalho remanescente e evidência por revisão.
- [Constituição](.specify/memory/constitution.md): monólito modular, domínios separados, dados compartilhados e documentação em português.

Ao alterar componentes, aplicar as skills Go em `.agents/skills`, **e obrigatoriamente `jornada-visual-parity` para qualquer mudança de aparência/skin/fonte/prévia/PPT**, atualizar os artefatos do domínio afetado e conferir código/cena/prévia/PPT. Requisitos planejados continuam requisitos, com lacunas explícitas; não devem ser removidos apenas para fazer a documentação parecer conforme ao código. Os registros antigos do quickstart são históricos e não aprovam a revisão atual.
