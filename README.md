# Jornada Studio — Gerador de Jornada

Editor desktop offline em **Go 1.27.2 / Fyne 2.8.2**, revisão **1.0.0.1**, para compor telas, conectar jornadas e preparar demonstrações. O projeto editável é salvo em `.jornada`; a saída PowerPoint mantém textos e formas editáveis. Java 1.6/EFL permanece visível e desabilitado.

O código fica em [apps/gerador-jornada](apps/gerador-jornada). O [guia da aplicação](apps/gerador-jornada/README.md) detalha operação, propriedades, atalhos e distribuição. As [33 specs](specs/001-gerador-jornada/spec.md) separam implementação existente, requisitos ainda pendentes e homologação manual.

## Executar a revisão atual

No Windows x64, abrir [JornadaStudio.exe](apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) ou extrair o [pacote 1.0.0.1](apps/gerador-jornada/dist/JornadaStudio-1.0.0.1-windows-amd64.zip). Não é necessário instalar Go, Node, Java ou Office para usar o editor. O diretório `dist` contém artefatos gerados, não o código-fonte canônico.

```powershell
# A partir da raiz do repositório:
./apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe --open ./apps/gerador-jornada/dist/estilos-legado.jornada
```

O [projeto das sete famílias](apps/gerador-jornada/dist/estilos-legado.jornada) e o [projeto dos compostos](apps/gerador-jornada/dist/compostos-legado.jornada) são fixtures de teste. Capturas verificadas: [prévia em memória](apps/gerador-jornada/build/legacy-preview.png) e [menu no executável Windows](apps/gerador-jornada/build/legacy-menu-native.png). Arquivos em `build` são evidências locais descartáveis.

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

## Estado real e limites

| Identidade | Valor na revisão atual |
| --- | --- |
| Software | 1.0.0.1 |
| Arquivo `.jornada` gravado | schemaVersion **3** / catalogVersion **3** |
| Pares reconhecidos pelo reader | 1/1, 2/1, 3/2 e 3/3, com os campos suportados; campos desconhecidos são recusados |
| Componentes inseríveis | Sete famílias base e três compostos: entrada de código de barras, ficha processada e menu |
| Biblioteca visual operacional | 205 chaves tipo/nome, 224 recursos referenciados, seis fontes físicas; oito definições indisponíveis fora da oferta |
| Entrega compilada nesta revisão | Windows amd64; Linux/macOS permanecem configuração futura |

O número do catálogo não comprova entrega da expansão para 16 famílias. Os tipos adicionais, variantes completas de relatórios, máscaras incrementais, aparência tipada/RGBA e perfil visual portátil por documento continuam no backlog. **Não há `Document.AppearanceProfile`, `skin.Profile` nem adaptador `visualprofile` na árvore atual**; suas descrições anteriores são histórico/plano de expansão.

O ramo RELATÓRIOS do menu tem uma fixture estrutural de 27 títulos e 148 folhas. [menu_reports.go](apps/gerador-jornada/internal/catalog/menu_reports.go) usa títulos e tokens sintéticos para completar essa contagem porque o XML integral não está no workspace. A contagem não comprova fidelidade literal. A ficha processada também não tem imagem independente para homologação. Nenhum componente consulta sistemas bancários ou executa pagamentos.

Fontes físicas acompanham o editor, mas **não são incorporadas ao PPT**. Disponibilidade/licença das fontes no computador do Office, fidelidade visual Windows/DPI 100/150/200% e aprovação humana continuam pendentes. O [quickstart, seção 22](specs/001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24) contém hashes, comandos executados, capturas e roteiro manual.

## Reproduzir a importação e o build

A partir de `apps/gerador-jornada`, com Python/PowerShell e as ferramentas de build:

```powershell
# Preparação local de Go/MinGW, quando ainda não disponíveis:
python -X utf8 scripts/bootstrap.py

# Regenerar os dados e incorporar os recursos disponíveis do acervo:
python -X utf8 scripts/import_legacy_skin.py

# Conferir o manifesto contra o XML atual e os bytes incorporados, sem regravar:
python -X utf8 scripts/import_legacy_skin.py --check

# Formatação, vet, testes Go/Python e pacote Windows:
python -X utf8 scripts/build.py --target windows --arch amd64 --check
```

O importador não substitui auditoria de proveniência, direitos ou cobertura integral do acervo. Arquivos complementares já incorporados são reutilizados; recursos ausentes não são inventados. Na revisão 1.0.0.1 passaram testes Go nos 13 pacotes com testes, `go vet`, gofmt, três testes Python de empacotamento, importador `--check`, build Windows e smoke nativo. Essas evidências são técnicas e não equivalem a homologação no PowerPoint ou aceite humano.

## Specs e manutenção

- [Spec coordenadora](specs/001-gerador-jornada/spec.md): integração e estado do produto; os demais domínios mantêm seus próprios contratos e backlogs.
- [Catálogo](specs/003-catalogo/spec.md), [Tela](specs/004-tela/spec.md) e [Propriedades](specs/028-propriedades/spec.md): famílias, instâncias e personalização.
- [Contrato `.jornada`](specs/001-gerador-jornada/contracts/project-format.md) e [contrato PPT](specs/001-gerador-jornada/contracts/pptx-export.md).
- [Research](specs/001-gerador-jornada/research.md): fontes, decisões, histórico e divergências conhecidas.
- [Tasks](specs/001-gerador-jornada/tasks.md) e [quickstart](specs/001-gerador-jornada/quickstart.md): trabalho remanescente e evidência por revisão.
- [Constituição](.specify/memory/constitution.md): monólito modular, domínios separados, dados compartilhados e documentação em português.

Ao alterar componentes, aplicar as skills Go em `.agents/skills`, **e obrigatoriamente `jornada-visual-parity` para qualquer mudança de aparência/skin/fonte/prévia/PPT**, atualizar os artefatos do domínio afetado e conferir código/cena/prévia/PPT. Requisitos planejados continuam requisitos, com lacunas explícitas; não devem ser removidos apenas para fazer a documentação parecer conforme ao código. Os registros antigos do quickstart são históricos e não aprovam a revisão atual.
