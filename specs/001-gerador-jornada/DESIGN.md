# DESIGN — Gerador de Jornada

**Status**: referência visual e regras de personalização; homologação manual pendente.
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

Cada feature individual tem `spec.md`, `DESIGN.md` e checklist de qualidade. O inventário completo fica em [research.md](research.md#inventário-de-componentes-e-designs-individuais). O catálogo das sete famílias também tem seu [DESIGN.md](../003-catalogo/DESIGN.md).

Alterar um DESIGN.md declara uma intenção visual para implementação e revisão. O EXE não lê Markdown como configuração. Antes de considerar uma alteração pronta, implementar o token correspondente, testar comportamento e comparar os estados de referência. Não apresentar um token planejado como controle de personalização já disponível.

Não criar um novo domínio para cada widget. Família pertence a Catálogo; composição e geometria pertencem a Tela; navegação a Jornada; temporização a Demonstração; persistência a Projeto. Painéis, código e renderização são adaptadores coordenados pela feature 001.

## Recursos e saídas

Imagens podem decorar fundo, faixa e skins. Textos, botões e demais controles devem continuar editáveis no PPT. A prévia usa o mesmo recurso do projeto, sem acessar o caminho original do legado. O Java usa somente estilos mapeados pelo perfil real; uma propriedade sem mapeamento produz diagnóstico.

Data/hora e indicadores informativos são valores ilustrativos editáveis nesta versão. Não são relógio vivo, autenticação, medição de rede nem consulta a equipamento. Essa escolha permite repetir roteiros e exportações de modo determinístico.

## Matriz de revisão

Validar janela 1366 × 768 e DPI 100%, 150% e 200%; zoom 25%, 100%, Ajustar e 400%; acentos, emoji, textos longos e estados vazios. Confirmar contraste, foco por teclado e mensagens contextuais. Cada DESIGN específico acrescenta suas dimensões, tokens, estados e casos de borda.

Registrar build, fixture e resultado em [quickstart.md](quickstart.md). Validação XML, testes de widgets e documentação revisada não equivalem a homologação do PowerPoint, EFL ou aprovação humana.
