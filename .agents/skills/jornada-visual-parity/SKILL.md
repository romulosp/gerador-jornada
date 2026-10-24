---
name: jornada-visual-parity
description: 'Use em qualquer alteração visual do Jornada Studio: componentes, skins, fontes, alpha, cena compartilhada, prévia interativa, PPT ou importação de material-legado. Garante que o estilo legado atravesse criação, canvas, prévia e PPT com evidência verificável, em vez de ficar restrito a uma camada ou a uma skin explícita.'
---

# Jornada Visual Parity

Esta skill existe porque o projeto já tinha requisitos de fidelidade visual que foram cumpridos apenas parcialmente: a skin era aplicada em uma camada, mas não no caminho completo. Use-a antes de implementar ou revisar qualquer aparência.

## Quando usar

- Inserir/alterar componente, família, variante, preset ou composto.
- Mexer em `internal/skin`, `material-legado`, `Style.Skin`, fontes, cores, imagens ou alpha.
- Alterar `presentation`, canvas, prévia interativa, inspector de aparência ou exportação PPT.
- Revisar uma entrega que afirma “skin aplicada”, “estilo do legado” ou “paridade visual”.

## Fontes obrigatórias, nesta ordem

1. `material-legado/elemetos_componente/xml/skin.xml` e os bytes em `material-legado/elemetos_componente/xml/{azul,icon,fonts,treinamento}`.
2. Classes Java relevantes: `SkinBuilderImpl`, `SkinUtil`, `EFLVirtualKeyboardNumeric`, `EFLReportEvent`, `TextFieldSkin`, `LabelSkin`, `ButtonSkin`.
3. [spec 001, seção “Reimplementação visual”](../../../specs/001-gerador-jornada/spec.md#reimplementação-visual-solicitada--2026-10-24), [029/FR-012](../../../specs/029-previa-interativa/spec.md#paridade-interativa-pós-v1--2026-10-19), [028/FR-015/FR-016](../../../specs/028-propriedades/spec.md#aparência-e-famílias-pós-v1--2026-10-19), [003/FR-009–FR-012](../../../specs/003-catalogo/spec.md).
4. Contrato de catálogo/PPT e importador `apps/gerador-jornada/scripts/import_legacy_skin.py`.
5. Evidência vigente: [quickstart §22](../../../specs/001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24).

Não substitua fonte ausente por “parecida”, nem XML ausente por transcrição inventada.

## Invariantes que não podem falhar

1. **Default é skin do legado.** Componente novo, com `Style.Skin == ""`, resolve a definição default da família. Ausência de skin **não** pode cair em cores genéricas (`#4864A0`, `#2563EB`). Overrides explícitos de cor/fonte continuam prevalecendo.
2. **Uma cena, três consumidores.** Canvas, prévia e PPT consomem as mesmas primitivas de `presentation`. A prévia não pode trocar a pintura por `widget.Button`, `widget.Entry`, `widget.Select` opacos.
3. **Controle nativo é só input transparente.** Foco, cursor, seleção, TAB e clipboard podem usar widget nativo, mas fundo/borda/texto visíveis vêm da cena. Use `legacyControl`/`legacyEntry` como referência.
4. **Fonte é bytes, não nome.** A face usada no Fyne precisa de `canvas.Text.FontSource` com os bytes SFNT incorporados. Família e flags físicas vêm do cabeçalho SFNT, não do atributo XML.
5. **Alpha é RGBA.** Distinga alpha ausente de `alpha="0"` explícito; preserve alfa parcial em fundos/bordas. `Hex()` sozinho não representa transparência.
6. **Importação é de desenvolvimento.** Gere `library.json` determinístico, com tipo/nome/estado, recurso, tamanho, SHA-256, família e flags. Runtime lê JSON estático e recursos embutidos; não interpreta XML Java.
7. **Definição só é ofertada se todos os recursos existirem.** Definição com imagem/fonte ausente sai inteira da oferta e fica auditável em `Unavailable`. Escolha antiga não resolve → diagnóstico e bloqueio de PPT; nunca fallback silencioso.
8. **Compostos são recompostos.** Relatório, menu, teclado, ficha e entrada têm partes (cabeçalho/corpo/rodapé, setas, teclas, campos) com skins/fontes próprias; não basta pintar a caixa externa.
9. **Layout usa a face real.** Meça com a fonte incorporada; transporte tamanho, cor, alinhamento, centralização vertical, quebra por palavra e `spacement` (menu 0,8).

## Red flags (já causaram retrabalho)

- `if component.Style.Skin == "" { ... cores genéricas ... }`.
- `widget.NewButton(component.Text)` ou `widget.NewEntry()` desenhando o componente na prévia.
- `Primitive.FontFamily` preenchido e nenhum `FontSource`.
- `Color.Hex()` retornando `""` para `A == 0` sem distinguir ausência.
- Parser de XML em runtime ou cópia de assets por basename sem manifesto/hash.
- Ofertar todas as definições, inclusive as com 1–3 imagens faltantes.
- Teste que só usa `Style.Skin` explícito e nunca insere um componente vazio.
- Build gerado sem captura da prévia, do canvas e do PPT na mesma fixture.

## Matriz mínima de verificação

| Cenário | Asserção |
| --- | --- |
| Inserir cada família com `Style.Skin == ""` | A cena contém decoração/fonte do default legado; não há cor genérica. |
| Canvas × prévia × PPT | As três saídas derivam das mesmas primitivas; a prévia mostra a mesma decoração. |
| Botão/campo na prévia | Cursor, TAB, Backspace e clipboard funcionam; fundo/borda/texto continuam EFL. |
| Fonte | `FontSource` não nulo; família SFNT correta; tamanho/cor/alinhamento conferidos. |
| Alpha | Fundo transparente, preto opaco e alfa parcial produzem o resultado esperado. |
| Importador | `--check` PASS; contagens e SHA-256 conferem; definições indisponíveis fora da oferta. |
| Recurso removido | Diagnóstico localizado; PPT recusa publicar; projeto antigo preservado. |
| Compostos | Relatório com cabeçalho/rodapé, teclado 14 teclas, menu com setas, ficha com partes. |

Execute `go test ./...`, `go vet ./...`, `gofmt -l internal cmd`, `python scripts/import_legacy_skin.py --check` e o build Windows. Teste automatizado não substitui DPI/Office/aceite humano.

## Definição de pronto

- Default e override cobertos por teste de cena, não só por teste de parser.
- Prévia e canvas usam o mesmo renderer; nenhum widget opaco cobre a skin.
- Font bytes e alpha verificados por asserção.
- Manifesto íntegro, sem definições quebradas na oferta.
- Fixtures `estilos-legado.jornada` e `compostos-legado.jornada`, capturas e PPT gerados.
- Limitações (fonte não embutida no PPT, 16 famílias, Office/DPI) explícitas; sem alegar homologação.
