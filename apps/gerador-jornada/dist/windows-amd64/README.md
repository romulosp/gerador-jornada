<!-- Cópia do README da aplicação; links de código/specs referem-se ao repositório. -->

# Jornada Studio — Gerador de Jornada

<img src="JornadaStudio.png" alt="Ícone do Jornada Studio" width="112">

**Software 1.0.0.11 · Go 1.27.2 · Fyne 2.8.2 · Desktop · Português**

O Jornada Studio permite desenhar telas, conectar etapas de atendimento e apresentar uma jornada antes de integrar o sistema definitivo. O autor posiciona componentes, configura suas propriedades, testa a navegação e gera uma apresentação PowerPoint com telas editáveis e demonstração guiada.

A solução moderniza a experiência do gerador legado sem repetir sua implementação em Node/React. O editor atual é uma aplicação Go com janela própria, recursos incorporados e funcionamento offline. Seu público inclui analistas, designers de interação, responsáveis pelo atendimento e desenvolvedores que precisam discutir um fluxo usando o mesmo projeto visual.

## Um exemplo de uso

Uma jornada de atendimento pode conter **Identificação → Confirmação → Conclusão**. Na primeira tela, um campo recebe um número pelo teclado virtual; um botão leva à confirmação. O roteiro demonstrativo pode mover o cursor, clicar, simular a digitação, esperar e navegar. O PPT resultante permite explorar as telas ou assistir à demonstração.

O projeto `.jornada` guarda a composição editável. O PPT é uma saída de apresentação; ele não substitui o projeto de edição nem executa as regras do sistema Java.

## Começar a usar

1. No Windows x64, execute `dist/windows-amd64/JornadaStudio.exe`. Na primeira execução, aguarde o vídeo com áudio terminar; depois, a apresentação depende da opção no launcher.
2. No launcher, escolha um recente, crie um projeto em branco/com exemplo ou abra uma `.jornada`. Confirme nome e pasta.
3. Na aba **Telas**, escolha uma tela. Na **Paleta**, insira e posicione os componentes sobre o fundo.
4. Edite as propriedades e clique em **Salvar alterações**. **Descartar** retorna aos valores confirmados.
5. Configure destinos na aba **Navegação** e etapas na aba **Demonstração**.
6. Use **Prévia Tela Atual** para revisar somente a tela selecionada ou **Prévia** para percorrer o projeto.
7. Salve o `.jornada` e exporte o **PowerPoint**.

### Splash, carregamento e autoria — 1.0.0.11

A apresentação abre centralizada em janela normal de 960×540 (unidades de interface), com bordas e barra de título, sem fullscreen ou maximização automática, reproduzindo o vídeo v2 com áudio. Não há barra, legenda de carga, autoria ou opção de pular. O vídeo se ajusta proporcionalmente à área interna da janela, sem contar a moldura nem cortar bordas. Faixas pretas preservam a proporção quando necessário. O usuário ainda relatou falha na revisão fullscreen; esta revisão remove a troca de modo e permite conferir a reprodução em janela normal. Fonte continua intacta e incorporada ao EXE. A autoria **Desenvolvido por Rômulo Penha - Canais Parceiro - 03** fica no rodapé do launcher.

Quando exibida, carregamento e reprodução acontecem em paralelo; o fade para o launcher só começa depois de ambos concluídos. O término vem do player, sem usar os aproximadamente dez segundos do arquivo como temporizador de liberação. -open respeita a mesma preferência de apresentação; salvar/exportar por flags sem interface continua sem vídeo. Fechar cancela a inicialização e interrompe áudio. Falha de mídia/recursos apresenta diagnóstico e impede avanço. Configuração de volume/mudo do Windows continua respeitada.

Na primeira execução com esta configuração de usuário, o vídeo aparece uma vez. Na seleção de projetos, **Mostrar splash screen ao iniciar** fica **desmarcado por padrão**. As próximas aberturas vão direto para a tela inicial depois de preparar os recursos. Marcar repete o vídeo nas próximas inicializações; desmarcar volta a suprimi-lo. A alteração é salva imediatamente e sobrevive a reinicializações, atualizações do EXE e mudanças de projeto. Uma apresentação cancelada ou com falha não consome a primeira exibição.

Foi necessário criar a política pura `startup.Preference` e o adaptador `desktop/startup_preferences.go`, além do coordenador `desktop.StartDesktop`. A preferência fica em `%APPDATA%/JornadaStudio/startup-preferences.json` (`presented` e `show_on_startup`), com publicação atômica; `JORNADA_CONFIG_DIR` isola testes. Não fica no `.jornada` nem na pasta do EXE. A construção do launcher não grava nem marca a opção. Mesmo sem vídeo, os recursos obrigatórios são carregados e suas falhas impedem abrir a edição.

O modelo do ciclo fica em internal/startup e na [spec de domínio 035](../../specs/035-inicializacao/spec.md); janela/carregador/player ficam nos adaptadores e na [spec 036](../../specs/036-splash-screen/spec.md). A janela nativa atende WM_PAINT/WM_SIZE e só inicia a reprodução após mídia pronta. A geometria usa o cliente físico, sem duplicar DPI, e a razão retornada pelo player é convertida em resolução. [Verificação nativa e roteiro humano](../../specs/036-splash-screen/quickstart.md). Formato de projeto permanece 3/4, sem alteração no modelo da jornada.

### Correção da abertura — 1.0.0.5

Corrigido o fechamento ao criar ou reabrir pelo launcher: o tema era reaplicado quando a janela nativa do editor ainda não existia. A configuração agora antecede NewWindow e não dispara callbacks durante o preenchimento inicial do seletor. Testados criar/reabrir com driver Windows nos três temas e com exemplo; o log ganhou fases técnicas de abertura sem conteúdo dos projetos. [Evidências e reprodução](../../specs/034-launcher-estrutura-projeto/quickstart.md#correção-de-encerramento--1005).

### Launcher, workspace e destinos entre jornadas (034)

O launcher mantém recentes locais e cria projetos em branco (somente fundo incorporado) ou com exemplo editável. Cada projeto fica numa pasta; ela é o workspace, sem novo formato. A aba inicial **Estrutura projeto** lista apenas as `.jornada` da raiz. Use **⋮ → Criar jornada** ou **Importar jornada existente**; arquivos de mesmo nome recebem sufixo e a origem importada fica intacta. **Atualizar** relê mudanças externas. A paleta e as demais abas continuam disponíveis, inclusive pelo menu **…** quando o painel estiver estreito.

Em uma jornada salva, a Prévia sugere criar o destino de uma folha de menu sem rota. Confirmar salva as alterações e o vínculo da origem, cria o arquivo e atualiza a árvore. Cancelar não grava nada. Na próxima ativação, a prévia abre `Journey.StartScreenID` do destino; **Voltar à jornada anterior** restaura tela, submenu e dados transitórios da origem. Prévia Tela Atual não segue destinos. Arquivo removido recebe diagnóstico e nova sugestão; corrupção não é tratada como ausência. A pilha comporta até 32 visitas.

O campo editável **Jornada destino** aceita apenas um nome de arquivo `.jornada` da mesma pasta. `Submenu`, `TargetScreenID` e `TargetJourney` são exclusivos. Catálogo 4 torna esse campo obrigatório no wire, inclusive vazio. Abrir catálogos antigos não regrava os arquivos; salvar solicita confirmação e conserva backup integral. Desfazer/refazer restaura o vínculo no documento; o arquivo criado permanece na pasta. A nova navegação está restrita à prévia, sem gerar fluxo entre arquivos no PPT ou Java.

A implementação exigiu separar launcher/sessão, criar recentes atômicos, resolver arquivos por nome relativo, validar importações antes de publicar e coordenar criação/vínculo com rollback caso a origem não possa ser salva. [Detalhes e evidências C1/C2/C3](../../specs/034-launcher-estrutura-projeto/quickstart.md). A regra de gerar apenas o EXE permanece vigente.

### Projeto de teste dos compostos 031–033

O executável aceita dois utilitários locais que não abrem a janela:

```powershell
JornadaStudio.exe -save-compounds .\dist\teste-compostos.jornada
JornadaStudio.exe -export-compounds .\dist\teste-compostos.pptx
```

O primeiro salva uma revisão com uma instância de código de barras, uma ficha processada e o menu da aplicação em duas telas; o segundo exporta a apresentação correspondente. Os valores digitáveis começam vazios, o menu inicia na raiz e nenhum roteiro é criado. Use uma **cópia** do arquivo para a sessão manual e registre a revisão/hash do EXE, o ambiente, o DPI e as capturas exigidas pelo [roteiro de teste humano](../../specs/001-gerador-jornada/quickstart.md#21-teste-humano-consolidado--features-001033-2026-10-24).

### Requisitos para executar

| Sistema | Arquitetura e requisitos de execução | Situação |
| --- | --- | --- |
| Windows | Windows 10/11 x64, sessão gráfica, driver gráfico funcional e permissão de gravação nos destinos escolhidos | Build local disponível; homologação manual continua separada |
| Linux | x86-64, ambiente gráfico X11/Wayland compatível, bibliotecas OpenGL/Mesa, Wayland/X11 do sistema e glibc compatível com a base de compilação | Build local gerado em contêiner Debian 12; interface precisa de teste na distribuição de destino |
| macOS Intel | macOS 15 ou posterior, Intel de 64 bits | Build nativo configurado; geração/homologação dependem de um Mac/runner |
| macOS Apple Silicon | macOS 15 ou posterior, ARM64 | Build nativo configurado; geração/homologação dependem de um Mac/runner |

**Para usar o editor, não é necessário instalar Go, Node, JVM ou PowerPoint.** No Windows, ícone, fontes e imagens internas acompanham o EXE. No Linux, o programa utiliza bibliotecas gráficas do sistema; não é uma distribuição Linux totalmente estática. No macOS, a distribuição é um bundle `.app` para teste interno, sem notarização ou assinatura Developer ID.

O primeiro conjunto de referência usa 800 × 600 unidades lógicas; as dimensões são configuráveis. A interface foi planejada para telas a partir de 1366 × 768, com zoom e painéis ajustáveis. DPI, acessibilidade e comportamento nativo precisam de teste por plataforma.

Para assistir e homologar as animações exportadas, use um visualizador PPTX compatível. A matriz de homologação prevista é PowerPoint desktop 2016/Microsoft 365 no Windows. Gerar o arquivo não exige Office instalado; compatibilidade visual com outros visualizadores precisa ser verificada.

### Onde ficam os executáveis

O comando de compilação de rotina gera somente este executável dentro de `apps/gerador-jornada/dist`:

```text
dist/
└── windows-amd64/JornadaStudio.exe
```

**Regra vigente, definida pelo usuário:** gerar por enquanto apenas `JornadaStudio.exe` para Windows x64. ZIP, TAR.GZ, instaladores, bundles `.app`, builds Linux/macOS, outras arquiteturas ou qualquer outro formato de compilação/distribuição exigem solicitação explícita do usuário para o formato desejado. Corrigir, testar ou recompilar não autoriza gerar essas saídas adicionais. A regra também se aplica a agentes e automações de entrega.

Pacotes, hashes, inventários e fixtures de revisões anteriores podem continuar em `dist` como histórico; não devem ser regenerados automaticamente. Recursos intermediários e arquivos temporários dos testes não são formatos adicionais de entrega. Salvar `.jornada` e exportar PPT continuam funcionalidades do aplicativo.

## Componentes visuais

As sete famílias têm propriedades tipadas. Os tamanhos abaixo são os defaults da paleta, em unidades lógicas, e podem ser alterados no editor.

| Componente | Propósito | Propriedades específicas disponíveis | Tamanho inicial |
| --- | --- | --- | --- |
| **Botão** | Representar uma ação e conectar telas | Texto, variante `confirm/cancel/yes/no` e destino de navegação | 110 × 60 |
| **Rótulo** | Exibir instruções, valores ilustrativos e legendas | Texto Unicode | 150 × 20 |
| **Título** | Destacar uma informação na composição | Texto Unicode; aceita decoração e alinhamento textual | 300 × 28 |
| **Campo de texto** | Simular entrada de dados na prévia | Valor inicial e formato `text/number/cellular/cep/cpf/date/ddd/nis` | 150 × 25 |
| **Seleção** | Apresentar opções de escolha única | Uma opção por linha e orientação vertical/horizontal | 180 × 60 |
| **Relatório/lista** | Representar resultados e informações de consulta | Cabeçalho e linhas de texto ordenadas | 300 × 345 |
| **Teclado numérico** | Alimentar um campo da mesma tela | Campo de destino; 0–9, setas, TAB e LIMPAR (Backspace) | 200 × 350 |

**Compostos do legado (031–033):**

| Componente | Propósito | Propriedades específicas | Tamanho inicial |
| --- | --- | --- | --- |
| **Código de barras** | Simular a linha digitável de um boleto | Rótulo; oito campos derivados e validação local por DV | 610 × 60 |
| **Código de barras processado** | Exibir uma ficha ilustrativa de leitura | Blocos com título e linhas CPF/CNPJ/Nome/Nome fantasia/Razão social | 530 × 120 |
| **Menu da aplicação** | Editar e explorar a árvore de serviços e jogos | Submenu, item, legenda, categoria, variante, habilitação, destino e menu inicial | 800 × 600 |

Os três aparecem em **ELEMENTOS DO LEGADO**, não entre as sete famílias primitivas. A árvore inicial do menu e os dados derivados do boleto são ilustrativos: não executam movimento bancário, consulta, autenticação, QR ou rede. O PPT preserva textos e formas editáveis, mas não calcula DV, não abre modal e não opera a árvore. Imagens de referência e as classes reais de módulo SISPL/SIGCB não estão no workspace; fidelidade visual e equivalência bancária continuam pendentes de evidência.

### Skins do acervo legado

O inspector oferece **Padrão do legado** ou uma definição compatível com a família. Sem skin explícita, o resolvedor aplica os defaults do acervo; não restaura a aparência genérica anterior. Família incorporada, tamanho/cor do texto, fundo, borda e alinhamento são personalizações confirmadas pelo formulário.

A biblioteca corrente tem **205 chaves tipo/nome, 224 recursos referenciados e seis fontes físicas**. Oito definições com imagens ausentes ficam fora da oferta. Esses números não são quantidade de tipos inseríveis: o catálogo atual tem sete famílias e três compostos. Imagens e fontes acompanham o EXE; runtime lê JSON estático e recursos incorporados, sem interpretar XML Java.

[O que foi necessário para aplicar o estilo](../../README.md#o-que-foi-necessário-para-recuperar-o-estilo-do-legado) explica a importação, correção de alpha, famílias SFNT, defaults, cena compartilhada, captura de texto transparente e exportação. Para corrigir os jogos, também foi necessário recuperar o `menu.xml` integral, preservar sua codificação ISO-8859-1 e os itens ativos, substituir os relatórios sintéticos, vincular cada item à chave do `skin.xml` e aplicar os PNGs normal/pressionado/desabilitado sem duplicar o nome já contido na imagem. A geração desses dados ocorre no desenvolvimento; o EXE usa os recursos incorporados. [A seção 24 do quickstart](../../specs/001-gerador-jornada/quickstart.md#24-menu-xml-e-skins-dos-jogos--1002) registra os checks, capturas e roteiro manual da revisão 1.0.0.2.


### Conjuntos prontos do legado

| Conjunto | Conteúdo criado | Dimensão de referência |
| --- | --- | --- |
| **Data e hora** | Dois rótulos de data/horário de Brasília | 126 × 32 |
| **Dados da lotérica** | Proprietário/usuário, terminal e identificação da unidade | 222 × 36 |
| **Infraestrutura** | Rótulos GATEWAY, DNS e HTTP | 78 × 36 |
| **Cabeçalho completo** | Oito rótulos combinando os três conjuntos anteriores | 426 × 40 |
| **Faixa de título** | Título editável sobre a faixa decorativa incorporada | 800 × 28 |

Esses conjuntos geram controles comuns editáveis, incluídos como uma ação de histórico. Data/hora e indicadores são exemplos fixos; não há relógio vivo, autenticação ou monitoramento de rede. A faixa permite alinhar **o texto** à esquerda, ao centro ou à direita, sem mover o elemento. Conjuntos não constituem um novo agrupamento persistido.

A árvore atual vem de `material-legado/menu/menu.xml`, recuperado e preservado no repositório. RELATÓRIOS conserva 27 títulos descendentes e 148 folhas reais. Milionária está ativa e LotoGol está comentado na raiz. Jogos usam suas skins exclusivas nos estados normal, pressionado e desabilitado. Projetos salvos não são reescritos: inserir novo Menu aplicação para obter a árvore corrigida. [Detalhes e evidências da 1.0.0.2](../../specs/001-gerador-jornada/quickstart.md#24-menu-xml-e-skins-dos-jogos--1002).

### Submenus e metadados dos movimentos — 1.0.0.3

O submenu da Mega-Sena usa a faixa original e dois botões centralizados de 164 × 52, com opções normal/especial. Cancelar/Voltar aparecem no rodapé do próprio menu; Sair aparece na raiz e encerra somente a prévia. Faixa principal só aparece quando há título configurado. O desenho e as regiões clicáveis usam as mesmas medidas extraídas de `material-legado/menu/PnlMenu.java`.

Em Propriedades, selecionar o item permite editar **Código do movimento**, **Nome do movimento**, **Opção do movimento** e **Atalho legado**, preservando espaços, zeros e acentos. Normal/especial não são identificados apenas pelo código compartilhado. A prévia mostra os metadados da folha selecionada; destino de tela continua explícito. Voltar/Cancelar limpam a seleção anterior sem modificar o documento.

O cadastro `movimento.xml`/`movimentos.xml` continua indisponível nos materiais acessíveis. Suas referências são preservadas, mas o cadastro e suas classes/telas não foram validados ou reconstruídos. Não confundir metadados editáveis com execução do sistema legado. [Verificação e limite da 1.0.0.3](../../specs/001-gerador-jornada/quickstart.md#25-submenus-java-e-movimentos--1003). Inserir novo Menu aplicação para receber os defaults corrigidos; arquivos existentes conservam os títulos e variantes autorados.

### Propriedades comuns e identificação

Cada instância possui nome amigável, **ID de arquivo**, posição X/Y, largura/altura, cor do texto, tamanho lógico da fonte e alinhamento textual. A imagem decorativa é configurada pela ação específica. Alguns estilos previstos no modelo ainda não têm edição completa ou equivalência em todos os alvos.

- **Tela:** ID com prefixo `Pnl`, por exemplo `PnlIdentificacao` para `01 · Identificação`.
- **Componente:** ID livre, por exemplo `BtnContinuar`; `Pnl` não é obrigatório.
- IDs têm até **40 caracteres**, incluindo prefixo e excluindo extensão; somente letras ASCII e números, começando por letra. Não aceitam acentos, espaços, símbolos, caminhos ou nomes reservados do Windows.
- IDs de tela são únicos no projeto; IDs de componente, na tela, sem distinguir caixa. Renomear a legenda não renomeia um ID confirmado.
- O ID de arquivo é independente do **UUID interno**, usado em navegação e vínculos. A interface remove acentos/símbolos durante a edição; excesso de tamanho e colisões impedem salvar.
- O gerador de fontes utiliza o nome-base confirmado, como `PnlIdentificacao.java` e `BtnContinuar.java`. A ação pública Java/EFL permanece desabilitada.

### Hierarquia, habilitação e camadas

- **Pai:** o campo **Pai** do inspector associa o componente a um contêiner da mesma tela. Pai ausente, autorreferência e ciclo são recusados; mover o pai desloca os descendentes uma vez, sem duplicar o deslocamento quando o filho também está selecionado.
- **Habilitado:** o campo **Habilitado** controla a interação. Componente desabilitado continua visível em contraste reduzido, mas não recebe clique, foco, digitação ou navegação na prévia; a instância e o documento não são apagados.
- **Sobreposição:** o inspector oferece **Manter posição atual**, **À frente do objeto**, **Atrás do objeto**, **Trazer para a frente** e **Enviar para o fundo**, com **Objeto de referência** opcional. A posição é persistida apenas como a ordem canônica da tela, sem campo extra; **Frente ↑ / Fundo ↓** continuam movendo um passo.
- Remover um pai mantém os filhos na tela como raízes, sem referência inválida. Duplicar tela remapeia pai, filhos, teclado e destinos internos de menu.

### Log de aplicação

O cabeçalho oferece **Log de aplicação**. O popup mostra o caminho efetivo e permite **Registrar diagnóstico agora**, **Copiar caminho** e **Abrir pasta**. O arquivo JSONL registra startup, erros de interface e contagens técnicas (telas, componentes, famílias e habilitados); **não** grava textos digitados, nomes de projeto/componente, caminhos de arquivo ou conteúdo de recursos. No Windows, a pasta fica em `%LOCALAPPDATA%\GeradorJornada\logs`; há rotação de 2 MiB por segmento e retenção de 20 arquivos.

## Áreas da aplicação e edição

| Área | Responsabilidade |
| --- | --- |
| **Telas** | Adicionar, selecionar, duplicar, excluir e ordenar páginas da jornada |
| **Paleta** | Disponibilizar famílias e conjuntos prontos com rolagem independente |
| **Estrutura** | Selecionar elementos por nome, inclusive quando estão sobrepostos |
| **Canvas** | Posicionar sobre o fundo, selecionar, mover, redimensionar e controlar camadas |
| **Propriedades** | Editar rascunho individual, salvar/descartar em rodapé fixo ou operar seleção coletiva |
| **Navegação** | Definir tela inicial e destinos declarativos de botões |
| **Demonstração** | Configurar roteiro explícito e reproduzi-lo separadamente |
| **Prévia** | Executar o projeto completo com dados temporários |
| **Prévia Tela Atual** | Abrir apenas a tela ativa, sem executar navegação/roteiros de outras telas |

O botão **Fundo** permite aplicar, trocar ou remover uma imagem PNG/JPEG. Os recursos são incorporados ao projeto, sem dependência do caminho original. O fundo Loterias CAIXA e a faixa do legado acompanham os novos projetos desta revisão.

O inventário documental contém **24 responsabilidades individuais**, incluindo as famílias/conjuntos e as ferramentas do editor. Renderização de elemento em edição/prévia, lista de telas, paleta, árvore, barra de ferramentas e painéis esquerdo/direito são responsabilidades de interface, e não novos tipos inseríveis na tela. A visualização completa de código é tratada na [spec própria](../../specs/030-visualizacao-codigo/spec.md); sua descrição não significa disponibilidade integral na interface atual, especialmente enquanto Java/EFL está suspenso.

Componentes suportam seleção múltipla por Ctrl/Shift+clique e caixa de seleção. O movimento coletivo preserva as distâncias e respeita os limites do conjunto. A alça de redimensionamento é individual. O cursor indica movimento ou redimensionamento conforme o ponto selecionado. Alinhar o elemento e alinhar seu texto são operações diferentes.

Botão direito abre **Alterar propriedades** e **Excluir elemento/seleção**. A tecla Delete, o menu e **Excluir** ao lado de Ajustar compartilham confirmação **Sim/Não**; sem seleção, Excluir fica desabilitado. A operação pode ser desfeita.

### Temas e saída do aplicativo

No topo, **Tema** permite escolher **Claro** (padrão), **Escuro** ou **Azul**. A mudança é imediata e a escolha fica salva para o usuário nas próximas aberturas. Ela personaliza os painéis e controles do editor; o fundo e as cores das telas desenhadas, prévias e exportações permanecem os definidos no projeto. Rascunhos, seleção, zoom e histórico são preservados.

Em tela cheia, o botão **Sair**, junto de **Restaurar**, encerra o programa sem precisar restaurar a janela. Ele fica oculto no modo normal, onde o fechamento do sistema continua disponível. Se houver alterações no projeto ou no formulário de propriedades, a confirmação **Sim/Não** protege a saída; **Não** mantém a sessão aberta, **Sim** descarta as alterações e encerra, inclusive com prévias abertas. Salve antes de sair para conservar suas mudanças.

### Atalhos atuais

| Atalho | Ação |
| --- | --- |
| Ctrl+S | Salvar projeto |
| Ctrl+Z / Ctrl+Y | Desfazer / refazer |
| Ctrl+A | Selecionar todos os componentes da tela |
| Ctrl+D | Duplicar seleção |
| Ctrl/Shift+clique | Adicionar/remover componente da seleção |
| Setas / Shift+setas | Mover 1 / 10 unidades lógicas |
| Delete | Solicitar exclusão da seleção |
| F12 | Alternar tela cheia quando encaminhado ao canvas; captura nativa na janela Windows |
| Escape | Sair da tela cheia; no modo normal, limpar seleção/cancelar manipulação |

O launcher usa uma janela de boas-vindas. O editor de projetos abre em **tela cheia por padrão**, com a ação **Restaurar** disponível. Essa escolha inicial não depende do último estado da sessão. No Windows, **Maximizar** entra em tela cheia do editor, sem decoração do sistema/barra de tarefas. O botão **Tela cheia/Restaurar**, junto à marca, permite alternar sem depender da barra oculta; F12 alterna e Escape retorna, inclusive com campos focados pela captura nativa. Documento, seleção, zoom e rascunhos são preservados. Não é modo quiosque: fechar, minimizar e alternar aplicações continuam disponíveis. Linux/macOS mantêm a maximização do sistema e têm o botão de tela cheia; o encaminhamento das teclas pelo canvas depende do foco/driver e exige homologação.

Atalhos de desenho respeitam o foco dos campos: Delete e setas continuam editando texto. Os atalhos implementados usam Ctrl; equivalências com Command e ergonomia nativa do macOS ainda precisam de homologação.

## Saídas e armazenamento

### Projeto `.jornada`

É um ZIP portátil contendo manifesto JSON e recursos locais PNG/JPEG. O formato atual grava **schemaVersion 3/catalogVersion 4**, com `fileId` separado do UUID. O reader reconhece os pares 1/1, 2/1, 3/2, 3/3 e 3/4, dentro dos campos suportados. Na versão 1, IDs de arquivo são sugeridos deterministicamente; o próximo salvamento grava 3/4. Não regrava o original durante a leitura. Versões desconhecidas, campos desconhecidos, estruturas inválidas e recursos inconsistentes são recusados. Reconhecer o par 3/2 não significa aceitar appearanceProfile de implementações anteriores: esse campo não existe na árvore atual. A sessão original não deve ser substituída por uma leitura recusada.

A gravação usa arquivo temporário e publicação atômica. Após alterações, existe recuperação automática com debounce de 2 segundos. Ela não substitui Salvar nem grava seleção, zoom, foco, undo/redo ou valores temporários da prévia.

| Sistema | Diretório de recuperação por usuário |
| --- | --- |
| Windows | `%LOCALAPPDATA%/GeradorJornada/recovery` |
| Linux | `$XDG_CONFIG_HOME/GeradorJornada/recovery`, ou `~/.config/GeradorJornada/recovery` |
| macOS | `~/Library/Application Support/GeradorJornada/recovery` |

### PowerPoint `.pptx`

Gerado diretamente em Go com PresentationML/Open XML, sem automação do Office. Inclui exploração por links e demonstração finita com cursor, clique, digitação simulada, espera e navegação. Uma revisita pode gerar outra ocorrência do slide. As transições disponíveis são **Fade suave**, **Deslizar**, **Revelar** e **Sem transição**; Fade de 300 ms é o padrão.

Textos e controles são formas editáveis; imagens ficam para fundos/decorações. O PPT não executa Java, macros ou regras de negócio. Exportação usa snapshot, permite cancelamento e admite um trabalho por vez. Estrutura ZIP/XML verificada automaticamente não substitui abrir e testar no PowerPoint.

### Java 1.6 / EFL

O botão é mantido visível, porém **temporariamente desabilitado**. O gerador prepara fontes por tela e componente, fábricas e contrato de navegação separado. Os arquivos principais usam os IDs confirmados e são organizados para integração Swing/EFL.

O alvo permanece Java SE 6, mesmo com JDK 17 no desenvolvimento. Testes estruturais usam ECJ 3.33.0, Animal Sniffer 1.24 e assinaturas SE 6 para conferir sintaxe, APIs e bytecode major 50. Isso não homologa EFL: faltam bibliotecas reais e bootstrap do hospedeiro. A aplicação Java gerada dependerá do ambiente Java/EFL; essa dependência não pertence ao editor Go.

## Tecnologias e abordagem

| Elemento | Versão/abordagem |
| --- | --- |
| Linguagem do editor | Go **1.27.2**, conforme `go.mod` |
| Interface | Fyne **2.8.2**, janela desktop com renderer próprio de composição |
| Compilação gráfica | CGO + compilador C nativo; MinGW **16.2.0** na preparação Windows fixada |
| Unicode | `golang.org/x/text` **0.22.0** |
| Formato de projeto | ZIP + JSON, writer **3/4**; reader reconhece **1/1, 2/1, 3/2, 3/3 e 3/4** nos campos suportados |
| Catálogo | catalogVersion **4**, sete famílias base e três compostos tipados |
| PPT | `archive/zip` e XML Go; PresentationML/Open XML |
| Integração legada | Fontes Java 1.6/Swing/EFL, ação temporariamente suspensa |
| Build/distribuição | PowerShell no Windows; Python 3.10+; Docker opcional para Linux; runners nativos no GitHub Actions |

O código segue **monólito modular**, um módulo Go e composição manual de dependências. Domínios não importam Fyne, filesystem, XML ou templates Java. Há funções/tipos coesos e adaptadores específicos, sem um serviço central que faça toda a edição e exportação. A documentação GoDoc/Javadoc e os comentários do projeto são em português.

| Domínio | Responsabilidade | Pacote |
| --- | --- | --- |
| Projeto | Identidade, revisão e metadados/recursos | `internal/project` |
| Catálogo | Famílias, defaults e propriedades tipadas | `internal/catalog` |
| Tela | Componentes, geometria, aparência e regras de FileID | `internal/screen` |
| Jornada | Tela inicial e destinos de navegação | `internal/journey` |
| Demonstração | Etapas e interpretação temporal | `internal/scenario` |
| Aplicação | Sessão, transações, histórico, snapshots e coordenação | `internal/application` |

`internal/adapter` contém desktop, persistência, PPT, Java, recursos incorporados e metadados de build. `internal/presentation` traduz a composição em primitivas compartilhadas de apresentação. A separação permite mudar o visual ou um exportador sem transferir essas responsabilidades para o domínio.

## Compilar em `dist` — 64 bits

A entrega corrente e as recompilações de manutenção geram somente **`JornadaStudio.exe` para Windows x64**. Outros formatos ou plataformas só podem ser gerados após solicitação explícita do usuário. As instruções opcionais abaixo são referência e não autorizam executá-las na manutenção de rotina.

### Windows x64

No diretório desta aplicação, com Python 3.10+ e PowerShell:

```powershell
python -X utf8 scripts/bootstrap.py
./scripts/check.ps1
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
```

O bootstrap baixa Go/MinGW para `.tools` usando versões e SHA-256 de [`toolchains-windows.json`](scripts/toolchains-windows.json). A preparação inicial exige rede. `go.mod`/`go.sum` fixam as dependências. A compilação incorpora ícone multirresolução, FileVersion/ProductVersion e suporte C ligado estaticamente no Windows.

[`scripts/build.ps1`](scripts/build.ps1) gera somente o EXE no diretório indicado, sem ZIP, inventário ou pacote adicional. A tarefa padrão do VS Code usa esse mesmo caminho. `-Console` fica reservado a diagnóstico solicitado. **Não usar `scripts/build.py` na rotina**, pois ele também chama o empacotador; usá-lo somente quando o usuário solicitar explicitamente a distribuição correspondente.

### Referência opcional: Linux x64 usando Docker — somente sob solicitação explícita

Com Docker Desktop iniciado e mecanismo de contêineres Linux:

```powershell
python -X utf8 scripts/build.py --target linux --arch amd64 --docker --check
```

O [`Dockerfile.linux`](scripts/Dockerfile.linux) fixa Go 1.27.2 em Debian 12 por digest, instala dependências de build e utiliza volumes de cache para módulos/compilação. Os resultados ficam no `dist` do host. A primeira preparação exige rede e espaço para imagem/cache. O contêiner compila e testa; ele não abre a interface gráfica no Windows.

O pacote desse contêiner deve ser testado em Debian 12 ou sistema com bibliotecas compatíveis; não presumir compatibilidade com distribuições mais antigas. O inventário registra a libc do ambiente de build.

### Referência opcional: Linux x64 nativo — somente sob solicitação explícita

Em Debian/Ubuntu, prepare as bibliotecas de desenvolvimento conforme a [documentação Fyne](https://docs.fyne.io/started/quick/), além de Go 1.27.2 e Python 3.10+:

```bash
sudo apt-get update
sudo apt-get install gcc pkg-config libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev
python3 scripts/build.py --target linux --arch amd64 --check
```

As dependências de desenvolvimento pertencem ao build; o usuário final recebe o executável e precisa das bibliotecas gráficas/runtime compatíveis da sua distribuição.

### Referência opcional: macOS de 64 bits — somente sob solicitação explícita

Em um Mac da arquitetura correspondente, instale Go 1.27.2, Python 3.10+ e as ferramentas nativas Apple (`xcode-select --install`). Então execute **um** dos comandos:

```bash
# Mac Intel
python3 scripts/build.py --target darwin --arch amd64 --check

# Mac Apple Silicon
python3 scripts/build.py --target darwin --arch arm64 --check
```

O build cria `.app` com ícone `.icns`, Info.plist e assinatura ad-hoc local. A versão completa do produto é 1.0.0.11; os campos convencionais do bundle usam três partes. Não gera notarização ou assinatura de distribuição Apple.

O script recusa simplesmente trocar GOOS/GOARCH sem toolchain nativo. Fyne utiliza CGO e exige compilador/bibliotecas do alvo; consulte a [compilação entre plataformas](https://docs.fyne.io/started/cross-compiling/). Docker deste projeto atende Linux; o macOS utiliza Mac ou runner próprio.

### Referência opcional: GitHub Actions — somente sob solicitação explícita

O workflow [`desktop-distribution.yml`](../../.github/workflows/desktop-distribution.yml) possui matriz Windows x64, Linux x64, macOS Intel e macOS ARM64, usando runners nativos e ações fixadas por commit. Não executar esse workflow na rotina: ele gera múltiplas distribuições e depende de solicitação explícita do usuário para esse escopo.

Cada job verifica o código, compila e disponibiliza seu pacote com hash como artifact por 14 dias. Não publica release nem marca aprovação humana. O workflow foi preparado neste repositório; sua execução remota e os builds macOS ainda precisam de evidência. Os nomes dos runners seguem a [referência oficial do GitHub](https://docs.github.com/en/actions/reference/runners/github-hosted-runners).

### Verificações e ferramentas de apoio

```powershell
# Windows: conferir formatação, vet, importadores e testes Go/Python, sem gerar distribuição
./scripts/check.ps1

# Acrescentar validação estrutural Java 6 com o JDK informado
./scripts/check.ps1 -Jdk17 'C:\Desenvolvimento\jdk-17.0.11'

# Consultar os inventários realmente gerados
python scripts/build.py --list
```

Nos sistemas Unix, `build.py --check` executa vet e testes Go antes da distribuição. O workflow também exercita o reconhecimento dos cabeçalhos PE/ELF/Mach-O, recusando pacotes de arquitetura incorreta. O script não executa o editor para validar o cabeçalho.

O comando principal também oferece `--save-example <caminho.jornada>` e `--export-example <caminho.pptx>` para conjuntos de referência, e `--smoke-screenshot <caminho.png>` para prova visual controlada. Essas saídas de apoio não representam aprovação funcional.

## Especificações, qualidade e contribuição

O planejamento usa **GitHub Spec Kit**, com constituição, requisitos identificados, plano, contratos, tarefas e roteiro de teste manual. Cada responsabilidade visual tem uma spec e um `DESIGN.md`; editar o Markdown registra o contrato, mas não muda o visual automaticamente no executável.

- [Especificação integrada](../../specs/001-gerador-jornada/spec.md).
- [Inventário das specs e designs individuais](../../specs/001-gerador-jornada/research.md#inventário-de-componentes-e-designs-individuais).
- [Contrato do projeto e migração](../../specs/001-gerador-jornada/contracts/project-format.md).
- [Roteiro manual e evidências por revisão](../../specs/001-gerador-jornada/quickstart.md).
- [Constituição do projeto](../../.specify/memory/constitution.md).

Para contribuir, atualize primeiro a spec do domínio afetado, mantenha a regra fora dos widgets/exportadores e documente tipos/métodos em português. Execute os checks pertinentes e entregue o pacote para teste manual antes de considerar o incremento aprovado. Não versione `.tools`, `dist`, `build`, caches, relatórios locais ou credenciais; versione código, recursos, scripts, manifestos de ferramentas e specs.

### Estado e limites da versão

O editor, persistência, prévias e exportador PPT têm implementação e verificações automáticas. As specs também descrevem evoluções que ainda não estão implementadas; uma caixa documental ou um build concluído não comprova todos os comportamentos.

Permanecem separados: homologação do PowerPoint, EFL real, DPI/driver gráficos, teste nativo em máquina limpa e aprovação humana por plataforma. Não há importação `.page`, colaboração simultânea, banco de negócio, integração externa, gravação automática do roteiro ou compilação Java dentro do editor nesta versão.

## Responsáveis e licença

O repositório concentra o conhecimento do projeto e as decisões discutidas com seus mantenedores. A aplicação ainda não possui um arquivo LICENSE próprio definido neste repositório. Recursos visuais do legado, marcas e dependências mantêm suas condições de origem; este README não define uma licença para eles.

## Revisão visual 1.0.0.2 — 2026-10-24

As sete famílias atuais e os compostos de código de barras/menu usam as decorações e fontes disponíveis do legado. Defaults são aplicados mesmo sem escolher skin; cores/fonte/alinhamento podem ser personalizados no inspector. A prévia mantém essa cena durante a interação, com captura de texto transparente, cursor e valores locais. Teclado novo usa 14 teclas em 200 × 350; instâncias antigas conservam seus bounds. LIMPAR emite Backspace e TAB percorre os campos habilitados.

A biblioteca é gerada em desenvolvimento por `python -X utf8 scripts/import_legacy_skin.py`; `--check` confere o resultado sem escrever. Runtime carrega somente os dados JSON e recursos incorporados. Oito definições com imagens ausentes não são oferecidas; escolhas antigas indisponíveis exibem diagnóstico e impedem exportação sem apagar o projeto.

O PPT conserva decorações originais e texto/formas editáveis, com famílias SFNT reais. As fontes não são incorporadas ao PPT: validar sua disponibilidade e licença no computador do Office. A revisão não conclui a expansão integral dos tipos descritos nas specs. Build, capturas, resultados e roteiro manual estão na [seção 22 do quickstart](../../specs/001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24).
