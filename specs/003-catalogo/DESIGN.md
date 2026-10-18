# DESIGN — Catálogo / ListComponents

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
