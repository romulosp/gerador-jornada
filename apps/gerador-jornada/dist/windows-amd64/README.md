> Pacote de distribuição: os links de código e specs deste guia são relativos ao repositório-fonte.

# Jornada Studio — Gerador de Jornada

<img src="JornadaStudio.png" alt="Ícone do Jornada Studio" width="112">

**Software 1.0.0.0 · Go 1.27.2 · Fyne 2.8.2 · Desktop · Português**

O Jornada Studio permite desenhar telas, conectar etapas de atendimento e apresentar uma jornada antes de integrar o sistema definitivo. O autor posiciona componentes, configura suas propriedades, testa a navegação e gera uma apresentação PowerPoint com telas editáveis e demonstração guiada.

A solução moderniza a experiência do gerador legado sem repetir sua implementação em Node/React. O editor atual é uma aplicação Go com janela própria, recursos incorporados e funcionamento offline. Seu público inclui analistas, designers de interação, responsáveis pelo atendimento e desenvolvedores que precisam discutir um fluxo usando o mesmo projeto visual.

## Um exemplo de uso

Uma jornada de atendimento pode conter **Identificação → Confirmação → Conclusão**. Na primeira tela, um campo recebe um número pelo teclado virtual; um botão leva à confirmação. O roteiro demonstrativo pode mover o cursor, clicar, simular a digitação, esperar e navegar. O PPT resultante permite explorar as telas ou assistir à demonstração.

O projeto `.jornada` guarda a composição editável. O PPT é uma saída de apresentação; ele não substitui o projeto de edição nem executa as regras do sistema Java.

## Começar a usar

1. Abra o pacote correspondente ao seu sistema e execute o programa.
2. Crie um projeto ou abra um arquivo `.jornada` existente.
3. Na aba **Telas**, escolha uma tela. Na **Paleta**, insira e posicione os componentes sobre o fundo.
4. Edite as propriedades e clique em **Salvar alterações**. **Descartar** retorna aos valores confirmados.
5. Configure destinos na aba **Navegação** e etapas na aba **Demonstração**.
6. Use **Prévia Tela Atual** para revisar somente a tela selecionada ou **Prévia** para percorrer o projeto.
7. Salve o `.jornada` e exporte o **PowerPoint**.

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

Os comandos de build gravam somente dentro de `apps/gerador-jornada/dist`:

```text
dist/
├── JornadaStudio.exe                         # atalho de compatibilidade do build Windows
├── windows-amd64/JornadaStudio.exe
├── linux-amd64/JornadaStudio
├── darwin-amd64/JornadaStudio.app/
├── darwin-arm64/JornadaStudio.app/
├── JornadaStudio-1.0.0.0-windows-amd64.zip
├── JornadaStudio-1.0.0.0-linux-amd64.tar.gz
├── JornadaStudio-1.0.0.0-darwin-amd64.zip
└── JornadaStudio-1.0.0.0-darwin-arm64.zip
```

Essa árvore mostra os **destinos previstos**, não afirma que todos os pacotes já existem. Consulte os disponíveis com `python scripts/build.py --list`. Cada pacote gerado contém `build-info.json` e `SHA256SUMS.txt`; o arquivo compactado tem um `.sha256` ao lado. O inventário registra versão, sistema, arquitetura, Go e compilador C, mantendo aprovação manual como pendente. Linux também inclui `runtime-libraries.txt`; macOS registra o SDK Apple quando compilado.

No Windows, extraia o ZIP e abra `windows-amd64/JornadaStudio.exe`. No Linux, extraia o TAR.GZ, entre em `linux-amd64` e execute `./JornadaStudio`; se o meio de cópia removeu a permissão, aplique `chmod +x JornadaStudio`. No macOS, extraia o ZIP no Mac e abra `JornadaStudio.app` da arquitetura correta.

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
| **Teclado numérico** | Alimentar um campo da mesma tela | Campo de destino; 0–9, Limpar e apagar último caractere | 200 × 280 |

**Comportamentos e limites importantes:**

- Botões disparam a navegação configurada na prévia completa. No modo Tela Atual, a navegação entre telas é desativada. Variantes e skins ainda não têm equivalência visual integral em todos os renderizadores.
- Rótulo e título são textos editáveis, inclusive no PPT; não são capturas de tela.
- O campo numérico restringe a entrada a dígitos na prévia. As demais opções de formato representam intenção de apresentação; não garantem máscara completa, consulta externa ou validação fiscal de CPF/NIS.
- A seleção mantém escolhas temporárias na prévia. As opções têm identidade interna e preservam sua ordem; vírgulas podem fazer parte do texto.
- O relatório contém dados ilustrativos locais; não consulta banco, API ou sistema de negócio.
- O teclado implementado tem **12 teclas**. A especificação de evolução contempla diferenças do legado ainda não entregues. Sem vínculo, o teclado pode ser desenhado, mas não altera um campo.

### Conjuntos prontos do legado

| Conjunto | Conteúdo criado | Dimensão de referência |
| --- | --- | --- |
| **Data e hora** | Dois rótulos de data/horário de Brasília | 126 × 32 |
| **Dados da lotérica** | Proprietário/usuário, terminal e identificação da unidade | 222 × 36 |
| **Infraestrutura** | Rótulos GATEWAY, DNS e HTTP | 78 × 36 |
| **Cabeçalho completo** | Oito rótulos combinando os três conjuntos anteriores | 426 × 40 |
| **Faixa de título** | Título editável sobre a faixa decorativa incorporada | 800 × 28 |

Esses conjuntos geram controles comuns editáveis, incluídos como uma ação de histórico. Data/hora e indicadores são exemplos fixos; não há relógio vivo, autenticação ou monitoramento de rede. A faixa permite alinhar **o texto** à esquerda, ao centro ou à direita, sem mover o elemento. Conjuntos não constituem um novo agrupamento persistido.

### Propriedades comuns e identificação

Cada instância possui nome amigável, **ID de arquivo**, posição X/Y, largura/altura, cor do texto, tamanho lógico da fonte e alinhamento textual. A imagem decorativa é configurada pela ação específica. Alguns estilos previstos no modelo ainda não têm edição completa ou equivalência em todos os alvos.

- **Tela:** ID com prefixo `Pnl`, por exemplo `PnlIdentificacao` para `01 · Identificação`.
- **Componente:** ID livre, por exemplo `BtnContinuar`; `Pnl` não é obrigatório.
- IDs têm até **40 caracteres**, incluindo prefixo e excluindo extensão; somente letras ASCII e números, começando por letra. Não aceitam acentos, espaços, símbolos, caminhos ou nomes reservados do Windows.
- IDs de tela são únicos no projeto; IDs de componente, na tela, sem distinguir caixa. Renomear a legenda não renomeia um ID confirmado.
- O ID de arquivo é independente do **UUID interno**, usado em navegação e vínculos. A interface remove acentos/símbolos durante a edição; excesso de tamanho e colisões impedem salvar.
- O gerador de fontes utiliza o nome-base confirmado, como `PnlIdentificacao.java` e `BtnContinuar.java`. A ação pública Java/EFL permanece desabilitada.

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

O aplicativo abre em **tela cheia por padrão**, com a ação **Restaurar** disponível. Essa escolha inicial não depende do último estado da sessão. No Windows, **Maximizar** entra em tela cheia do editor, sem decoração do sistema/barra de tarefas. O botão **Tela cheia/Restaurar**, junto à marca, permite alternar sem depender da barra oculta; F12 alterna e Escape retorna, inclusive com campos focados pela captura nativa. Documento, seleção, zoom e rascunhos são preservados. Não é modo quiosque: fechar, minimizar e alternar aplicações continuam disponíveis. Linux/macOS mantêm a maximização do sistema e têm o botão de tela cheia; o encaminhamento das teclas pelo canvas depende do foco/driver e exige homologação.

Atalhos de desenho respeitam o foco dos campos: Delete e setas continuam editando texto. Os atalhos implementados usam Ctrl; equivalências com Command e ergonomia nativa do macOS ainda precisam de homologação.

## Saídas e armazenamento

### Projeto `.jornada`

É um ZIP portátil contendo manifesto JSON e recursos locais PNG/JPEG. O formato atual tem **schemaVersion 2**, com `fileId` separado do UUID. O leitor migra a versão 1 por sugestões determinísticas; o próximo salvamento grava a versão 2. Não regrava o original durante a leitura. Versões desconhecidas, estruturas inválidas e recursos inconsistentes são recusados.

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
| Formato de projeto | ZIP + JSON, schemaVersion **2**, leitor compatível com migração v1 |
| Catálogo | catalogVersion **1**, sete famílias tipadas |
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

A entrega corrente e as recompilações de manutenção abrangem somente **Windows x64**, conforme orientação do usuário. As instruções e configurações Linux/macOS abaixo ficam como referência futura; não exigem novos builds ou homologação nesta etapa.

### Windows x64

No diretório desta aplicação, com Python 3.10+ e PowerShell:

```powershell
python -X utf8 scripts/bootstrap.py
python -X utf8 scripts/build.py --target windows --arch amd64 --check
```

O bootstrap baixa Go/MinGW para `.tools` usando versões e SHA-256 de [`toolchains-windows.json`](scripts/toolchains-windows.json). A preparação inicial exige rede. `go.mod`/`go.sum` fixam as dependências. A compilação incorpora ícone multirresolução, FileVersion/ProductVersion e suporte C ligado estaticamente no Windows.

O comando anterior [`scripts/build.ps1`](scripts/build.ps1) continua disponível e gera diretamente `dist/JornadaStudio.exe`; `-Console` habilita console para diagnóstico. O empacotador também tenta atualizar esse caminho de compatibilidade; se estiver aberto/em uso, conserva o antigo e entrega o novo EXE em `dist/windows-amd64`, sem encerrar processos do usuário.

### Linux x64 usando Docker a partir do Windows

Com Docker Desktop iniciado e mecanismo de contêineres Linux:

```powershell
python -X utf8 scripts/build.py --target linux --arch amd64 --docker --check
```

O [`Dockerfile.linux`](scripts/Dockerfile.linux) fixa Go 1.27.2 em Debian 12 por digest, instala dependências de build e utiliza volumes de cache para módulos/compilação. Os resultados ficam no `dist` do host. A primeira preparação exige rede e espaço para imagem/cache. O contêiner compila e testa; ele não abre a interface gráfica no Windows.

O pacote desse contêiner deve ser testado em Debian 12 ou sistema com bibliotecas compatíveis; não presumir compatibilidade com distribuições mais antigas. O inventário registra a libc do ambiente de build.

### Linux x64 nativo

Em Debian/Ubuntu, prepare as bibliotecas de desenvolvimento conforme a [documentação Fyne](https://docs.fyne.io/started/quick/), além de Go 1.27.2 e Python 3.10+:

```bash
sudo apt-get update
sudo apt-get install gcc pkg-config libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev
python3 scripts/build.py --target linux --arch amd64 --check
```

As dependências de desenvolvimento pertencem ao build; o usuário final recebe o executável e precisa das bibliotecas gráficas/runtime compatíveis da sua distribuição.

### macOS de 64 bits

Em um Mac da arquitetura correspondente, instale Go 1.27.2, Python 3.10+ e as ferramentas nativas Apple (`xcode-select --install`). Então execute **um** dos comandos:

```bash
# Mac Intel
python3 scripts/build.py --target darwin --arch amd64 --check

# Mac Apple Silicon
python3 scripts/build.py --target darwin --arch arm64 --check
```

O build cria `.app` com ícone `.icns`, Info.plist e assinatura ad-hoc local. A versão completa do produto continua 1.0.0.0; os campos convencionais do bundle usam três partes. Não gera notarização ou assinatura de distribuição Apple.

O script recusa simplesmente trocar GOOS/GOARCH sem toolchain nativo. Fyne utiliza CGO e exige compilador/bibliotecas do alvo; consulte a [compilação entre plataformas](https://docs.fyne.io/started/cross-compiling/). Docker deste projeto atende Linux; o macOS utiliza Mac ou runner próprio.

### Todas as plataformas pelo GitHub Actions

O workflow [`desktop-distribution.yml`](../../.github/workflows/desktop-distribution.yml) possui matriz Windows x64, Linux x64, macOS Intel e macOS ARM64, usando runners nativos e ações fixadas por commit. Depois de disponibilizar o repositório no GitHub, execute manualmente **Actions → Distribuição desktop 64 bits → Run workflow**.

Cada job verifica o código, compila e disponibiliza seu pacote com hash como artifact por 14 dias. Não publica release nem marca aprovação humana. O workflow foi preparado neste repositório; sua execução remota e os builds macOS ainda precisam de evidência. Os nomes dos runners seguem a [referência oficial do GitHub](https://docs.github.com/en/actions/reference/runners/github-hosted-runners).

### Verificações e ferramentas de apoio

```powershell
# Windows: formatação, vet, testes Go e empacotamento
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
