# DESIGN — Gerador de Jornada

## Aplicação na revisão 1.0.0.1

Editor Go/Fyne com sete famílias base, três compostos, autoria de telas/jornada/roteiros, persistência 3/3, prévia e PPT. A revisão 1.0.0.1 aplica skins e fontes do legado ao recorte existente.

**Revisão ainda necessária:** FR-049–FR-054 e a expansão de 16 famílias não estão integralmente atendidos. Faltam tipos adicionais, perfil visual portátil por documento, inventário completo de direitos/recursos e homologação Windows/DPI/Office. Java/EFL continua desabilitado. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Implementação**: [apps/gerador-jornada](../../apps/gerador-jornada).

## Duas superfícies visuais

O editor tem seu próprio tema, painéis, abas, formulários e foco. A tela desenhada tem dimensões lógicas, fundo, componentes e suas propriedades. Personalizar o tema do editor não deve alterar a tela criada nem as saídas.

O fundo solicitado é o arquivo original [fundo.jpg](../../material-legado/src/assets/images/fundo.jpg), de 800 × 600. Uma cópia idêntica fica incorporada ao EXE e é registrada como recurso do projeto. SHA-256: `AA2B4C1DDC31D8BFDE78B3B80210CA4D18E09EA192B4AE2E8C3771541868A1A2`.

## Tokens comuns da interface

| Token | Referência | Responsabilidade |
| --- | --- | --- |
| Superfície principal | `#F8FAFC` | Área externa à tela desenhada |
| Texto principal | `#0F172A` | Ações e campos |
| Texto secundário | `#64748B` | Instruções e rótulos de seção |
| Ação/seleção/foco | `#2563EB` | Estado perceptível com indicação além de cor |
| Fundo de ação secundária | `#E8EEF6` | Botões e listas |
| Borda | `#CBD5E1` | Separadores e caixas |
| Corpo da interface | 14 unidades | Texto de controles do editor |
| Título do projeto | 24, negrito | Hierarquia do cabeçalho |
| Área mínima de ação | 36 unidades de altura | Botões e entradas |

Esses valores pertencem ao tema da interface. Overrides de texto, fonte, preenchimento, borda e imagem de um componente pertencem ao documento quando suportados pela família. Não acrescentar parâmetros decorativos de Fyne às entidades de domínio.

## Composição da janela

- Cabeçalho compacto com identidade e ações de arquivo, prévia e exportação; botões não esticam para ocupar toda sua altura.
- Área de edição em três painéis redimensionáveis. À esquerda, abas Telas, Paleta (componentes) e Estrutura com rolagem independente.
- Centro com histórico, zoom efetivo, Ajustar e Fundo. O zoom modifica apenas a visualização. A imagem da tela é a área de drop, com margem externa de 24 unidades.
- À direita, propriedades da tela ou do componente selecionado. Rótulos e valores precisam continuar alcançáveis com rolagem.
- A superfície da jornada mantém o fundo atrás dos componentes. Contorno de seleção e alças ficam na frente e não são exportados.

## Contrato de personalização por componente

Cada feature individual tem spec, DESIGN e checklist de qualidade. O inventário histórico fica em [research.md](research.md#inventário-de-componentes-e-designs-individuais); a expansão Java/EFL acrescenta a matriz integral referenciada abaixo. O catálogo, inicialmente de sete famílias, tem seu [DESIGN.md](../003-catalogo/DESIGN.md).

Alterar um DESIGN.md declara uma intenção visual para implementação e revisão. O EXE não lê Markdown como configuração. Antes de considerar uma alteração pronta, implementar o token correspondente, testar comportamento e comparar os estados de referência. Não apresentar um token planejado como controle de personalização já disponível.

Não criar um novo domínio para cada widget. Família pertence a Catálogo; composição e geometria pertencem a Tela; navegação a Jornada; temporização a Demonstração; persistência a Projeto. Painéis, código e renderização são adaptadores coordenados pela feature 001.

## Recursos e saídas

Imagens podem decorar fundo, faixa e skins. Textos, botões e demais controles devem continuar editáveis no PPT. A prévia usa o mesmo recurso do projeto, sem acessar o caminho original do legado. O Java usa somente estilos mapeados pelo perfil real; uma propriedade sem mapeamento produz diagnóstico.

Data/hora e indicadores informativos são valores ilustrativos editáveis nesta versão. Não são relógio vivo, autenticação, medição de rede nem consulta a equipamento. Essa escolha permite repetir roteiros e exportações de modo determinístico.

## Matriz de revisão

Validar janela 1366 × 768 e DPI 100%, 150% e 200%; zoom 25%, 100%, Ajustar e 400%; acentos, emoji, textos longos e estados vazios. Confirmar contraste, foco por teclado e mensagens contextuais. Cada DESIGN específico acrescenta suas dimensões, tokens, estados e casos de borda.

Registrar build, fixture e resultado em [quickstart.md](quickstart.md). Validação XML, testes de widgets e documentação revisada não equivalem a homologação do PowerPoint, EFL ou aprovação humana.

## Ícone do programa — 2026-10-18

Fonte escolhida pelo usuário: `sugestao_nao_final_icone.png`, 1254 × 1254. Cópia incorporada: `apps/gerador-jornada/internal/adapter/builtinassets/assets/app-icon.png`, SHA-256 `E070BA60470E35331007AFB91F3DB12F8DFB4114F513DA7441C39A269060F8EA`. Preservar a composição e as cores azul/amarelo do original, incluindo seu fundo escuro; não aplicar recortes ou remover o fundo.

O adaptador desktop fornece a imagem ao Fyne. O build gera um grupo ICO com quadros 16, 24, 32, 48, 64, 128 e 256 px, incorporado aos recursos PE por MinGW windres. PNG e scripts são versionados; ICO, RC e SYso gerados são descartáveis e ficam ignorados pelo Git. Nenhuma ferramenta de conversão é necessária no computador do usuário.

Validar Explorer, título da janela e barra de tarefas nos DPIs definidos acima. Este recurso identifica o editor e não é adicionado aos projetos ou às apresentações exportadas.

## IDs de arquivo no painel

Tela e instâncias visuais têm campo ID, separado de Nome, com regras e mensagens do [painel de propriedades](../028-propriedades/DESIGN.md#id-do-arquivo). Prefixo Pnl é obrigatório somente nas telas; exemplos PnlIdentificacao e BtnContinuar. O atributo é metadado e não é desenhado no PPT/prévia. Texto visível e nome amigável continuam aceitando acentos.


## Temas e saída em tela cheia

O isolamento de tema continua obrigatório no alvo pós-v1; skins/fontes da jornada não alteram a chrome.

A barra de ferramentas define as paletas e controles em [025/DESIGN.md](../025-barra-ferramentas/DESIGN.md#sair-e-seletor-de-tema). Claro é padrão; Escuro e Azul são preferências locais. As superfícies do editor mudam, enquanto fundo, texto e componentes da jornada mantêm os valores do projeto. Sair aparece à direita somente em tela cheia e usa a confirmação de encerramento com Sim/Não. Rascunhos são protegidos conforme [028-propriedades](../028-propriedades/spec.md).

## Paridade Java/EFL pós-v1 — 2026-10-19

Escopo autorizado é todo o acervo visual, preservando a história das sete famílias. Referência primária Java/XML/builder/pintura e [matriz de evidência](research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19); React é secundária identificada. Biblioteca skin estática gerada de dados puros + recursos no adaptador, cena/métricas comuns a canvas/miniatura/prévia/PPT, sem executar Java ou ler XML legado em runtime.

Aplicar aparência nova a antigos conservando bounds/IDs/FileIDs/conteúdo/vínculos/ordem/overrides e left histórico. Fonte/face/RGBA/spacement/estados exigem evidência; bold/italic XML não são estilo aplicado pelo builder. Não resolver paridade com tema global, widget nativo opaco, rasterização ou substituição silenciosa. Fontes em outras_fontes existem por nome, mas licenças/métricas/recursos/delegates continuam gates. Defaults geométricos novos não redimensionam antigos.

Piloto Report e teclado têm [plano 012](../012-relatorio-efl/plan.md) e [plano 013](../013-teclado-numerico-efl/plan.md). Demais famílias seguem donos de domínio/adaptadores, sem backlog avulso. Comparação Windows/DPI e PowerPoint e aceite explícito da revisão são etapas separadas. Este DESIGN não comprova implementação/homologação ou reativa Java.

## Aplicação efetiva dos estilos — 2026-10-24

A cena agora resolve os defaults do acervo mesmo sem escolha explícita de skin. Imagens de botão/relatório/faixa continuam decorações; legendas, campos, opções e linhas permanecem texto e formas. O tema Claro/Escuro/Azul personaliza o editor; a jornada usa as fontes e cores resolvidas de sua cena.

A implementação corrente está em `internal/presentation/{scene,legacy_parts,menus,barcodes}.go` e no renderer Fyne comum `internal/adapter/desktop/scene_render.go`. Prévia usa botões/opções/teclas com renderer próprio e entrada transparente sobre a decoração EFL. Editor e prévia carregam as faces físicas incorporadas. Skins ausentes mostram diagnóstico e bloqueiam PPT, preservando o documento.

Referência renderizada em memória: `apps/gerador-jornada/build/legacy-preview.png`. Inspeção visual do agente não equivale a teste Windows em DPI 100/150/200% ou aceite do usuário; veja seção 22 do quickstart.
