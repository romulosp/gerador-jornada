# DESIGN — Formulários de propriedades — Properties

## Aplicação na revisão 1.0.0.1

Inspector atual edita identidade, geometria, conteúdo tipado, skin compatível, família física, tamanho/cor de texto, fundo/borda, alinhamento, habilitação, pai e sobreposição. Salvar publica uma transação; Descartar não publica.

**Revisão ainda necessária:** Não há painel Recursos/referência nem inspector_appearance.go/visualprofile atuais. Overrides com presença tipada/RGBA/restauração por campo e diagnóstico de todas as revisões/recursos continuam planejados. Estado funcional e evidências na [spec](spec.md#estado-verificado-1-0-0-1). Tokens e matrizes de expansão abaixo continuam normativos; não implicam aprovação visual.


**Status**: referência visual e implementação parcial reconciliadas para 1.0.0.1; homologação pendente.
**Spec**: [spec.md](spec.md) | **Domínio**: [001-gerador-jornada](../001-gerador-jornada/spec.md)
**Base comum**: [DESIGN do editor](../001-gerador-jornada/DESIGN.md).
**Fonte visual/funcional**: [legado](../../material-legado/src/components/Properties/index.js). Revisar o comportamento, não apenas a existência da pasta.

## Responsabilidade visual

O componente deve permitir personalizar conteúdo, geometria e aparência com validação antes de aplicar. Sua aparência não pode introduzir um comportamento diferente do descrito na spec.

## Geometria e composição

Seções Conteúdo, Geometria, Aparência e Ações; rótulos em português, unidades explícitas e controles com área mínima 36.

### FR-008: rodapé fixo sem travar o divisor — 2026-10-21

Somente o **corpo do inspector** deve rolar nos dois eixos. Em [Editor.design](../../apps/gerador-jornada/internal/adapter/desktop/editor.go), `container.NewScroll(e.Inspector)` substitui `container.NewVScroll(e.Inspector)`: o `Split` do Fyne limita o arraste pelos mínimos dos filhos, e a rolagem apenas vertical propagava ao painel a largura mínima do formulário, inflada por seletores/caminhos de fontes/versões longos. A rolagem bidirecional desacopla as duas dimensões mínimas e conserva os dados integrais e o acesso horizontal aos campos, sem encurtar valores.

O `container.NewBorder` do painel direito permanece com cabeçalho e rodapé fixos, irmãos do corpo rolável. **Salvar alterações / Descartar** não são removidos, estreitados nem movidos para dentro da rolagem. O mínimo natural dos controles fixos ainda limita o estreitamento; não impor largura fixa nem prometer colapso ilimitado. Rascunho, Save/Discard e histórico permanecem inalterados. A rolagem interna vertical de Recursos / referência continua independente; esta correção altera o corpo externo, não a consulta de FR-016.

Conferir largura/posição reais ao mover o divisor em janela 1280 × 720, com tela, 16 famílias e seleção múltipla; alcançar campos extensos nos dois eixos e confirmar cabeçalho/rodapé fixos. Salvar/Descartar devem receber cliques nos botões renderizados depois da rolagem, conservando a transação/descarte existentes. Quatro regressões preparadas (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_split_test.go`), **não executadas por política Windows**, sem contorno. **Evidência final do executor principal (2026-10-21):** gofmt em dois arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` concluiu a compilação de 15 suítes, **sem executá-las**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade PASS: CRC, todos os SHA256SUMS, sidecar SHA-256, PE amd64, ambas as cópias do EXE idênticas, README definitivo, manifesto com aprovações pendentes e ausência de backups terminados em `~`. Identidades finais e evidências no [quickstart, seção 11](../001-gerador-jornada/quickstart.md#11-redimensionamento-do-painel-de-propriedades--2026-10-21). [T129](../001-gerador-jornada/tasks.md#phase-23-redimensionamento-do-painel-de-propriedades--2026-10-21) permanece aberta como verificação mista: sem execução dos testes Go, validação de UI/arraste nativos, DPI, revisão manual ou aceite; `manualApproval` e `redistributionApproval` continuam `pending`.

A superfície pertence ao editor ou à prévia. Não deve ser rasterizada nem incluída nos slides como decoração da jornada. Estados de interação não alteram as coordenadas persistidas.

## Tokens e pontos de personalização

Valores abaixo são a referência desta feature; diferenças do código atual devem ser verificadas e implementadas antes de declarar conformidade. Propriedades ainda sem controle visual ficam como trabalho explícito, sem inventar que o formulário já as oferece.

| Token / parte | Valor de referência | Onde personalizar / responsabilidade |
| --- | --- | --- |
| texto | #0F172A | tema de formulário |
| entrada | #FFFFFF | tema de campo |
| borda | #CBD5E1 | tema de campo |
| foco | #2563EB | tema de foco |

Conteúdo e opções configuráveis: Valores editáveis do modelo e rascunho de formulário separado; tipos de entrada apropriados a número, texto, opção e cor.

## Estados

Não alterado; alterado; inválido; aplicado; seleção trocada.

| Situação | Resultado visual e regra |
| --- | --- |
| Referência | Mostrar os tokens e a geometria definidos acima, com o conteúdo da fixture. |
| Seleção no editor | Contorno e alças ficam em camada transitória; não aparecem no PPT ou nos fontes Java. |
| Foco por teclado | Exibir indicação perceptível sem deslocar a composição; texto digitado deve seguir o foco correto. |
| Operação inválida | Explicar a causa junto ao contexto; conservar a aparência do último estado confirmado. |
| Estado vazio ou ausente | Distinguir conteúdo vazio permitido de recurso/referência quebrados; não preencher dados fictícios silenciosamente. |

## Personalização segura

1. Alterar os valores de referência deste DESIGN.md e identificar o token e o requisito afetados.
2. Confirmar se a propriedade pertence ao documento ou somente ao tema da interface; não acrescentar estado de UI ao arquivo .jornada.
3. Atualizar a implementação responsável por essa propriedade. Alterar Markdown, isoladamente, não muda o EXE.
4. Repetir os cenários da spec e comparar antes/depois em 100%, 150% e 200% de DPI.
5. Para elementos de tela, salvar/reabrir e conferir a mesma personalização no editor, prévia e PPT. Estilo Java sem mapeamento EFL homologado deve produzir diagnóstico, nunca aproximação silenciosa.

## Limites e acessibilidade

- Inteiro inválido, cor fora do formato, nome vazio, referências inválidas e texto com caracteres especiais. O Properties legado era um placeholder; o comportamento funcional estava em ToolsRight.
- Manter fonte legível, foco visível e mensagens que não dependam exclusivamente de cor.
- Fontes da chrome seguem o tema do editor; fontes do conteúdo da jornada exigem face/recurso comprovados. Fallback exploratório explícito do conteúdo não homologa fidelidade. Conferir glifos/métricas e texto longo.
- Imagens são apenas decoração; nunca substituir todos os textos e controles por uma captura de tela.

## Referências de implementação

Os caminhos abaixo identificam o trabalho existente; não comprovam cobertura integral desta spec:

- [internal/adapter/desktop/inspector.go](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go)

## Roteiro de revisão visual

- Capturar o estado de referência e uma variação de cada token permitido.
- Exercitar: Editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
- Conferir os estados descritos, casos inválidos e retorno ao estado anterior.
- Registrar versão/hash do build, fixture, resolução, DPI, executor e resultado em [quickstart.md](../001-gerador-jornada/quickstart.md).
- Homologação humana permanece pendente até o usuário testar a revisão entregue.

## Revisão de interação — 2026-10-18

- **FR-008**: Manter botão Salvar alterações visível fora da rolagem do formulário, inclusive quando o último campo estiver fora da área visível.
- **FR-009**: Aplicar o rascunho válido em uma ação reversível e atualizar a tela imediatamente após salvar; Descártá-lo pelo botão Descartar restaura os valores confirmados.
- **FR-010**: Exibir Alinhamento do texto, com esquerda, centro e direita, em controles textuais. A alteração modifica texto dentro da caixa, sem mover o elemento.
- **FR-011**: Com seleção múltipla, mostrar quantidade e comandos coletivos; não exibir formulário de um único componente como se editasse todos.

Os novos requisitos são normativos; verificar todos os estados afetados antes de considerar esta revisão homologada.

## ID do arquivo

Adicionar ID logo após Tela/Nome, usando o mesmo campo de entrada dos dados comuns. Tela: dica `Pnl + nome, até 40 caracteres`, exemplo `PnlIdentificacao → PnlIdentificacao.java`. Componente: `até 40 letras/números, sem acentos; Pnl não obrigatório`. O ID técnico UUID não aparece como campo editável. Texto amigável continua podendo ter acentos. Limpeza de acentos/pontuação acontece na digitação; erro de tamanho/prefixo/colisão é explícito e não salva valor inválido. Nome e ID são independentes após a sugestão inicial.


## Encerramento com rascunho

O diálogo de saída protege também o formulário em edição. Tema não remonta o formulário nem aplica valores. Responder Não conserva os campos e o modo de janela; Sim encerra descartando valores ainda não aplicados. A confirmação pertence à integração desktop, não às propriedades persistidas dos componentes.

## Perfil visual e overrides — expansão 2026-10-19

FR-015–FR-016: Aparência mostra skin/versão/origem/status e valor efetivo, distinguindo **Herdado** de **Personalizado** por campo. Fonte lista faces internas/recurso real e disponibilidade/licença, não só string de família. Estilo XML declarado não prova Bold efetivo do builder. Restaurar ao padrão remove override sem zerar campo/alterar geometria. Distinguir falso/alpha zero explícitos de ausência.

Formulários tipados incluem segunda legenda de button, variantes/linhas Report, novos campos/seleções/tabela/teclados/progresso/contêiner/dialog, capacidades e ParentID válido. Campo indisponível explica bloqueio; não persiste foco/pressed/scroll/cursor runtime. Indicadores de origem não entram no desenho/PPT. Preservar rodapé Salvar/Descartar e proteção do rascunho. Implementação 001/T113 e regras nos domínios; este DESIGN não conclui técnica ou aprovação.

## Recursos / referência — correção de 2026-10-21

> **Removido em 2026-10-23:** o campo foi retirado do inspector; esta seção é histórica.

Clarificação visual de **FR-016**, sem FR novo: **Recursos / referência** é uma consulta informativa, não uma entrada editável. Em inspector_appearance.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_appearance.go`), usar `container.NewVScroll(label)` com `widget.Label`, `TextWrapWord` e texto completo, sem truncamento ou corte de quantidade de avisos. A área rolável tem altura de referência **160 unidades lógicas**, configurada por `SetMinSize(fyne.NewSize(0, 160))`: limita o crescimento intrínseco do formulário pela lista, não é um máximo rígido de layout. Listas maiores rolam dentro dessa área, sem empurrar todas as ações para baixo. Manter o rodapé Salvar/Descartar fixo existente.

**HintText exato:** Avisos somente deste componente. Role a lista para consultar os detalhes.

Mostrar somente dependências da seleção/rascunho e da revisão exata, incluindo partes compostas e estados futuros; não confundir um caminho compartilhado com a identidade de outra definição. Preservar avisos reais e falhas globais de integridade, sem fallback. Preferir diagnóstico local, suprimir aviso genérico de métrica já medida com SFNT e deduplicar por código/tipo/nome/estado/recurso, não por mensagem. Mudar Skin/Face da fonte/Versão do perfil atualiza o texto sem salvar nem alterar documento/histórico. A limpeza de ParentID para diagnóstico ocorre somente na cópia temporária; não é edição de hierarquia.

Revisão nativa pendente: conferir lista curta/longa (incluindo cem avisos), quebra de texto integral, acesso ao último detalhe por rolagem, dica/foco e seletores em DPI 100/150/200%. Evidência final do executor principal para esta correção: gofmt em quatro arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` PASS, com 15 suítes **somente compiladas, não executadas**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade do pacote PASS (CRC, SHA256SUMS, sidecar SHA-256, PE amd64, README exato, inventário exato, ausência de backups terminados em `~` e ambas as cópias do EXE idênticas). [T128](../001-gerador-jornada/tasks.md#phase-22-correção-de-recursos--referência--2026-10-21) permanece aberta como verificação mista: testes Go não executados por política Windows, sem contorno; UI nativa, DPI, revisão manual e aceite não realizados. Identidades finais somente no [quickstart, seção 10](../001-gerador-jornada/quickstart.md#10-correção-de-recursos--referência--2026-10-21); não herdar hashes/PASS históricos nem tratar compilação como homologação visual.

### Disponibilidade local e portabilidade — clarificação de FR-016, 2026-10-21

A captura do usuário efetivamente visualizada pelo executor principal mostra `PROFILE_RESOURCE_NOT_EMBEDDED` para os caminhos lógicos /azul/botao_aux1_pressed.png e /azul/botao_pressed.png. Os PNGs já pertencem ao acervo incorporado, com identidade binária conferida por Python; o problema é confundir **não incluído no arquivo do projeto** com **indisponível nesta instalação**. A correção histórica de escopo/rolagem, acima, não resolvia essa apresentação.

Em `componentDiagnostics`, reclassificar visualmente somente `PROFILE_RESOURCE_NOT_EMBEDDED` cujo caminho seja resolvido por `sceneResourceData(document, group.SkinVersion, path)`, com SHA-256/tamanho/MIME exatos. Não basta encontrar o mesmo nome em outra revisão. Conservar a linha de identificação do perfil e os demais diagnósticos; depois deles, apresentar:

| Condição | Apresentação exigida |
| --- | --- |
| Referência com bytes exatos disponíveis localmente | **Disponível no aplicativo · PATH (não incluído no arquivo do projeto)**, uma linha por caminho, mesmo com vários estados/origens |
| Uma ou mais referências disponíveis | Uma única mensagem comum: **Aviso de portabilidade, não de carregamento: os recursos acima dependem deste acervo exato do aplicativo. Não foram copiados para o projeto; disponibilidade não comprova licença de redistribuição nem fidelidade integral.** |
| Referência não resolvida | Manter `PROFILE_RESOURCE_NOT_EMBEDDED`, localização e caminho, com **Recurso indisponível nesta instalação e não incluído no projeto; nenhuma substituição foi aplicada. Portabilidade e licença de redistribuição não comprovadas.** |
| Recurso já incorporado ao projeto, validado, mas com aviso persistido obsoleto de não incorporação | Não reapresentar essa condição nem anunciar dependência externa inexistente; preservar o perfil e os demais avisos, inclusive legais |
| Falta real, bloqueio, corrupção ou ressalva legal | Conservar os códigos e detalhes pertinentes; não converter `PROFILE_INVALID`, falhas de carregamento ou licenças em sucesso |

`PATH` representa o caminho lógico consultado, não um caminho externo de disco. Esta apresentação não autoriza copiar recursos nem muda hashes, licenças, status, modelo/schema ou diagnósticos globais/persistidos. `visualprofile.Diagnostics`, `visualprofile.ComponentDiagnostics`, `Capture` e PPT estrito permanecem intactos; disponibilidade local não homologa fidelidade ou redistribuição. A consulta continua por seleção/rascunho/revisão, mantendo texto completo não editável, quebra de palavras, rolagem de referência de 160 unidades, dica exata e rodapé fixo. Não acrescentar entrada editável ou alterar o divisor.

Nove regressões implementadas em resource_availability_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/resource_availability_test.go`), incluindo o subteste de bytes já incorporados e válidos com aviso persistido obsoleto, **não executadas por política Windows**, sem contorno. **Checks finais do executor principal (2026-10-21), após a última edição:** gofmt dos dois arquivos Go PASS; `go vet ./...` PASS (`RESOURCE_FINAL_VET_PASS`); `go test -c -o ./build/compile-tests/ ./...` PASS (`RESOURCE_FINAL_COMPILE_PASS`), com **15 suítes somente compiladas, ZERO execuções**; quatro testes Python de empacotamento PASS; tarefa VS Code de build Windows na raiz exit 0 e empacotamento Python sem `--check` concluído. Revisão Explore somente leitura sem defeito concreto apontado, não teste executado.

Integridade canônica conferida: CRC, SHA256SUMS/bytes locais, conjunto exato de entradas sem backups `~`, sidecar SHA-256, manifesto/build-info, produto/PE amd64, texto literal do novo aviso dentro do EXE e README congelado. **Histórico preservado:** inicialmente EXE canônico e ZIP foram atualizados, mas a primeira cópia de compatibilidade foi bloqueada por estar em uso e conservou identidade anterior. A checagem ad hoc inicial falhou somente pela exigência incorreta de igualdade das cópias naquele momento; não houve falha de build/ZIP, e a leitura independente validou a distribuição canônica.

**Compatibilidade RESOLVIDA após o usuário informar “fechei a aplicação” (2026-10-21):** o executor principal efetuou `Copy-Item` do [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) para a [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe), com `Get-FileHash` independente nos dois caminhos PASS e marcador `COMPATIBILITY_COPY_UPDATED_AND_VERIFIED_PASS SHA256=48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d`. Ambos agora têm **34.608.128 bytes e o mesmo SHA-256**; não há divergência de entrega atual. O ZIP não foi refeito: **20.408.876 bytes**, SHA-256 `b40cc96c901b51d00ca042831c419c222f44f90709b7d0328baf70695d8bdf01`, com checks já PASS. Nenhum aplicativo ou teste Go foi iniciado nesta finalização, nenhum processo foi encerrado pelo agente e nenhuma política foi alterada. README congelado e pacote preservados, sem necessidade de reempacotar. Identidades atuais e roteiro na [seção 12 do quickstart](../001-gerador-jornada/quickstart.md#12-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21); hashes anteriores permanecem históricos.

Revisar em UI nativa as duas linhas de PNG, outras dependências disponíveis e o aviso comum único, mistura de disponíveis/indisponíveis e atualização dos seletores sem salvar, nos DPI 100/150/200%. A captura analisada é evidência do sintoma anterior, não verificação visual da nova revisão. [T130](../001-gerador-jornada/tasks.md#phase-24-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21) permanece aberta como **verificação mista**: 15 suítes Go somente compiladas, não executadas; quatro testes Python e build/empacotamento PASS; aplicativo novo, revisão nativa/DPI/Office e aceite não realizados. `manualApproval` e `redistributionApproval` continuam `pending`. T120–T122 e os demais gates permanecem inalterados; igualdade das cópias e resultados técnicos não são homologação.

## Alinhamento textual universal — 2026-10-20

FR-017 amplia FR-010: em seleção individual, manter o `FormItem` **Alinhamento do texto** visível nas 16 famílias, inclusive vazias, após os campos de conteúdo e antes da aparência avançada. Opções: **Padrão da skin / Esquerda / Centro / Direita**. Dica: altera somente o texto; não move nem redimensiona o componente. Mostrar **Posição do elemento** como seção distinta junto aos botões de geometria, sem confundir alinhamento textual com posicionamento.

O alinhamento atua dentro das regiões textuais existentes: legendas, células, áreas de relatório e teclas mantêm suas caixas, recorte e decoração. No combo, incluir o popup de opções e **Sem seleção**. Não deslocar fundos, contêineres, filhos ou hit areas. Manter carregamento do valor, rodapé fixo, Save/Discard e Undo/Redo; **Padrão da skin** retira somente o override. Sem novo schema, catálogo ou domínio.

Revisão manual aberta: conferir ordem/visibilidade/foco, conteúdos vazios e longos, as quatro opções, confirmação/descarte/histórico e reabertura em DPI 100/150/200%, incluindo popup. Evidências técnicas de 20/10 informadas pelo executor principal no [quickstart, seção 8](../001-gerador-jornada/quickstart.md#8-alinhamento-textual-universal--2026-10-20): gofmt/vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e compilação Windows pela tarefa VS Code PASS. Isso não substitui a revisão nativa/DPI/PowerPoint ou o aceite; execução Go permanece não realizada pelo bloqueio prévio Windows, sem contorno.

## Sobreposição de componentes — 2026-10-20

FR-018: inserir os campos de camada após a geometria comum e antes das propriedades específicas, em todas as 16 famílias. Não confundir sobreposição com X/Y, alinhamento textual ou ParentID.

| Controle | Conteúdo e estado |
| --- | --- |
| **Sobreposição** | Manter posição atual / À frente do objeto / Atrás do objeto / Trazer para a frente / Enviar para o fundo; iniciar sempre no modo neutro Manter posição atual |
| Dica | Camada N de total (1 = fundo). Altera somente a ordem dos objetos ao salvar. |
| **Objeto de referência** | Iniciar em Selecione um objeto, desabilitado; habilitar somente nos dois modos relativos. Listar outros objetos da mesma tela como ordinal · nome, sem UUID exposto e sem o próprio componente |

A escolha permanece em rascunho, detectável pela proteção existente; o canvas confirmado não muda antes de **Salvar alterações**. Rodapé fixo **Salvar alterações / Descartar** permanece o mesmo. Salvar confirma camada e demais campos em uma transação: falha conserva todos os valores confirmados. Descartar restaura o formulário neutro; Undo/Redo restaura a edição completa. Ao remontar, mostrar a nova camada atual e não uma relação permanente com a referência.

**Frente ↑ / Fundo ↓** existentes permanecem comandos imediatos de um passo; **Trazer para a frente / Enviar para o fundo** do formulário são extremos confirmados ao salvar. Somente a ordem plana muda; pais/filhos conservam ParentID, geometria e suas próprias posições na lista, sem arrastar subárvore nem criar restrição relacional persistida.

Revisar foco, dica, nomes repetidos, tela unitária, erro de referência, rascunho, descarte/histórico e sobreposição real em DPI 100/150/200%. Implementação em inspector_layer.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_layer.go`); regressões implementadas, **não executadas** pela política Windows, sem contorno. O executor principal informou gofmt dos sete arquivos e vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e build/pacote Windows com integridade PASS. Evidências e identidade na [seção 9 do quickstart](../001-gerador-jornada/quickstart.md#9-sobreposição-de-componentes--2026-10-20); T127/004-T014 permanecem abertas, sem execução Go, revisão nativa/DPI/Office ou aprovação manual. Hashes e PASS anteriores são históricos.

## Log da consulta visível — clarificação de FR-016, 2026-10-21

> **Removido em 2026-10-23:** campo **Recursos / referência**, registro automático do inspector e botão de diagnóstico da prévia não existem mais; o painel mantém divisor, corpo rolável e rodapé fixo. O Log de aplicação captura manualmente o componente selecionado ou a tela atual. O texto a seguir é histórico.

Consumir [001/FR-057](../001-gerador-jornada/spec.md#log-de-aplicação--2026-10-21) e a [clarificação da spec](spec.md#log-da-consulta-visível--clarificação-de-fr-016-2026-10-21), sem novo FR ou logger duplicado. Manter **Recursos / referência** informativo, completo/rolável, dica, divisor, corpo bidirecional e rodapé fixo existentes. A consulta automática acompanha o relatório realmente mostrado ao selecionar e trocar Skin/Face da fonte/Versão do perfil, usando o rascunho sem salvar.

No cabeçalho, **Log de aplicação** abre popup de referência **720 × 340 unidades lógicas**, com caminho efetivo legível/quebrável, aviso de privacidade, resultado da gravação e ações **Registrar diagnóstico agora / Copiar caminho / Abrir pasta**. Abertura e registro manual forçam captura do rascunho visível ou dos componentes da tela atual sem seleção; o botão de diagnóstico da prévia também registra suas consultas. Mostrar falha de inicialização/gravação consultável, distinguindo-a do erro da aplicação; não anunciar sucesso se `Err` indicar falha. O caminho exibido é o real, inclusive fallback temporário, e só vai à área de transferência por ação explícita.

Log técnico não é cópia do painel: códigos/localização originais e mensagens builtin permitidas são preservados; strings personalizadas desconhecidas viram hash/omissão. `displayed` apenas indica código presente no relatório; `report_sha256` correlaciona a legenda sem texto bruto. Registrar resolução esperado/observado da revisão exata, inclusive arquivada/estados futuros, sem fallback ou mudança de recursos/licenças. Conteúdo de entradas/senhas, nomes de projeto/componente e `FilePath` não entram no arquivo; standard log só hash/bytes, sem rede. Abrir/capturar não altera formulário, rascunho, documento ou histórico.

**Estado final técnico — pronto para coleta (2026-10-21):** produção concluída e 11 arquivos Go formatados (sete existentes e quatro novos). Executor principal: `APPLICATION_LOG_VET_PASS`, `APPLICATION_LOG_COMPILE_ONLY_PASS`, **15 suítes somente compiladas, ZERO execuções Go**; os 12 testes novos (9 de núcleo + 3 de UI) não foram executados. Quatro testes Python de empacotamento, tarefa Windows na raiz (exit 0), pacote sem `--check` e integridade somente leitura PASS. Ambos os EXEs atualizados/idênticos e ZIP identificado na [seção 13](../001-gerador-jornada/quickstart.md#13-log-de-aplicação--2026-10-21); README congelado já empacotado, produto/schema/catálogo **1.0.0.0/3/2** inalterados. Hashes acima permanecem históricos, sem validar LOG.

**Revisão real pendente:** fechar a instância antiga e abrir o novo EXE identificado para reproduzir/enviar o log. Falha continua relatada, causa não mapeada e solução não comprovada; nenhum aplicativo foi iniciado nesta entrega. Marcadores no binário não provam execução, startup/runtime ou uso do logger. [T131](../001-gerador-jornada/tasks.md#phase-25-log-de-aplicação--2026-10-21) permanece aberta como verificação mista, assim como T130 e gates T120–T122; execução Go bloqueada por política Windows sem contorno, UI nativa/DPI 100/150/200%, Office e aceite pendentes, aprovações `pending`. Não criar logs artificiais nem afirmar que o usuário já abriu o novo EXE; diagnóstico do erro real depende da coleta.

## Formulários dos códigos de barras — planejamento 2026-10-23

FR-019: conservar layout/rodapé/divisor existentes. Formulário de 031 mostra rótulo e campos comuns; não oferecer editor de algoritmo ou de valores capturados. 032 apresenta lista de blocos, título opcional e lista de linhas do bloco selecionado, cada uma com **Rótulo** (combo) e **Valor** (texto literal). Exibir ações de adicionar/remover/reordenar com contexto do item; não usar título/rótulo como chave e não permitir remover último bloco/linha. A área pode rolar, mantendo Salvar alterações/Descartar acessíveis fora da rolagem.

Combo altera o rascunho da linha selecionada, conserva o valor e informa que os dados são ilustrativos. Todos os campos devem integrar a proteção de rascunho existente, inclusive ao reordenar/remover; não publicar por mudar seleção no popup. Remontar após Salvar/Descartar deve restaurar o estado confirmado. Não desenhar esses combos no componente processado da jornada/PPT. Contrato detalhado em [032/DESIGN](../032-codigo-barras-processado-efl/DESIGN.md), trabalho em [001/T136](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23); nenhuma revisão visual executada nesta entrega.

## Árvore de Menu aplicação — 2026-10-23

FR-020 acrescenta árvore de menus/submenus/itens no corpo rolável existente, sem alterar divisor/rodapé fixo. Mostrar identidade contextual por caminho para títulos repetidos; editar legenda/categoria/variante/enabled, ordem/add/remove e submenu ou destino de folha em alternativas exclusivas. **Menu inicial** seleciona raiz/submenu por ID e define estado fixo do canvas/PPT, não caminho da prévia. Selecionar/expandir nós não publica revisão. Salvar confirma toda a árvore em uma ação e Descartar restaura confirmado, inclusive listas e inicial; erro indica nó e preserva rascunho. Menus vazios são válidos; limites16/1000 e remoção de inicial inválida não são contornados por widgets. Detalhes visuais em [033/DESIGN](../033-menu-aplicacao-efl/DESIGN.md), formulário único001/T141, sem implementação/regressão duplicada ou homologação por este texto.
