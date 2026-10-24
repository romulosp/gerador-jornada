# DESIGN — Catálogo / ListComponents

## Aplicação na revisão 1.0.0.1

Catálogo 3 oferece dez Kinds: sete famílias base e barcodeInput/barcodeProcessed/applicationMenu. A biblioteca visual contém 205 chaves tipo/nome, 224 recursos referenciados, seis fontes físicas e oito definições retiradas por recursos ausentes. Skin não cria um novo Kind.

**Revisão ainda necessária:** O metamodelo declarativo de capacidades, as nove famílias adicionais da expansão e as variantes completas de Report/twoLabel/alphaKeyboard não foram implementados. Quantidade de skins não comprova quantidade de componentes nem homologação. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Spec**: [spec.md](spec.md) | **Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Referência**: [ListComponents/index.js](../../material-legado/src/components/ListComponents/index.js).

## Responsabilidade

ListComponents registra sete famílias: EFLButton, EFLLabel, EFLLabelTitulo1, EFLTextField, EFLRadio, EFLReportEvent e EFLVirtualKeyboardNumeric. JComponent é uma base técnica do legado, não um oitavo controle arrastável. DateInfo, LotericaInfo e InfraInfo compõem o cabeçalho; as demais pastas de interface não são famílias serializadas.

## Organização visual

A aba Paleta (componentes) deve listar todas as famílias com nomes em português e permitir identificar o nome EFL correspondente. Os conjuntos informativos e a faixa de título aparecem em uma seção separada. Cada item aceita clique e arraste para a imagem da tela. A lista tem rolagem própria, independente das páginas.

| Parte | Token de referência | Personalização |
| --- | --- | --- |
| Item disponível | `#E8EEF6`, texto `#0F172A` | Tema da paleta, sem alterar cores das telas |
| Seção | 11, negrito, `#64748B` | Tema do painel |
| Área de ação | Altura mínima 36 | Layout da paleta |
| Nome de família | Nome amigável e identificação EFL acessível | Definição do catálogo/apresentação |
| Conteúdo do controle criado | Defaults tipados de cada família | DESIGN individual e propriedades da instância |

## Designs individuais

- [Botão](../007-botao-efl/DESIGN.md).
- [Rótulo](../008-rotulo-efl/DESIGN.md).
- [Título](../009-titulo-efl/DESIGN.md).
- [Campo de texto](../010-campo-texto-efl/DESIGN.md).
- [Seleção](../011-selecao-efl/DESIGN.md).
- [Relatório/lista](../012-relatorio-efl/DESIGN.md).
- [Teclado numérico](../013-teclado-numerico-efl/DESIGN.md).

## Limites

Conjuntos prontos expandem-se em componentes das famílias existentes; não introduzem novos discriminadores no arquivo .jornada. A definição de uma família não contém widgets nem comandos de exportação. Alterar tokens aqui não modifica automaticamente o EXE: implementar os adaptadores responsáveis e revisar o DESIGN da família afetada.

O catálogo representa intenção funcional, enquanto as imagens originais dão referência visual. Templates Java antigos e suas assinaturas não comprovam compatibilidade EFL.

## Identificação das instâncias

Todas as famílias compartilham o campo ID descrito no [design de propriedades](../028-propriedades/DESIGN.md#id-do-arquivo). Não mostrar o ID como legenda no desenho, PPT ou prévia; é metadado de nomeação de arquivos, distinto do conteúdo visual. Conjuntos prontos atribuem IDs aos controles efetivamente criados, sem inventar uma entidade de agrupamento persistida.

## Expansão autorizada — 2026-10-19

As sete famílias acima são histórico da base, não limite do catálogo 2. Galeria organiza todo o acervo por semântica, skins/variantes/previews do mesmo resolver, origem e disponibilidade. Incluir ToggleButton, TextArea, PasswordField, Combo/UF, Table, AlphaKeyboard, ProgressBar e Panel/Dialog; duas legendas é variante button, subclasses Report são variantes, Frame/UF e 014–018 são presets. Scroll é capacidade do contêiner/composto, não item por seta. JComponent continua base técnica, não controle arrastável. ParentID opcional é composição autorizada de Tela, não conversão da seleção transitória em grupo.

Definições com recursos ausentes ficam fora da oferta de novas escolhas; uma referência antiga indisponível continua preservada e mostra diagnóstico. O comportamento de manter toda entrada incompleta visível é histórico anterior a FR-058. Skin da jornada não altera chrome. Selecionar variante não redimensiona instância existente; button atual 110 × 60, teclado novo 200 × 350, antigo 200 × 280 preservado. Fonte exige face real, separando flags XML e efeito Java. Modelos 003/T009–T016 e paleta/inspector 001/T113, sem duplicação; [matriz integral](../001-gerador-jornada/research.md#matriz-de-cobertura-javaefl--expansão-2026-10-19). Homologações abertas.
