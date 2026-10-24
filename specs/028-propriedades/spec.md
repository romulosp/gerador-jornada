# Feature Specification: Formulários de propriedades — Properties

**Feature Branch**: `028-propriedades` (identificador documental; nenhuma branch criada)
**Created**: 2026-10-18
**Status**: Implementação parcial na revisão 1.0.0.1; requisitos remanescentes e homologação/aceite pendentes
**Input**: Pedido do usuário: especificar cada componente do legado e permitir personalização em DESIGN.md.
**Domínio responsável**: [001-gerador-jornada](../001-gerador-jornada/spec.md); integração em [001-gerador-jornada](../001-gerador-jornada/spec.md).
**Referência observada**: [legado](../../material-legado/src/components/Properties/index.js). Contrato visual específico em [DESIGN.md](DESIGN.md).

<a id="estado-verificado-1-0-0-1"></a>
## Estado verificado — 2026-10-24 / 1.0.0.1

**Implementação existente:** Inspector atual edita identidade, geometria, conteúdo tipado, skin compatível, família física, tamanho/cor de texto, fundo/borda, alinhamento, habilitação, pai e sobreposição. Salvar publica uma transação; Descartar não publica.

**Pendências e limites:** Não há painel Recursos/referência nem inspector_appearance.go/visualprofile atuais. Overrides com presença tipada/RGBA/restauração por campo e diagnóstico de todas as revisões/recursos continuam planejados.

**Fontes verificadas:** [inspector.go](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go), [style_validation.go](../../apps/gerador-jornada/internal/presentation/style_validation.go), [zz_smoke_test.go](../../apps/gerador-jornada/internal/adapter/desktop/zz_smoke_test.go). Os checks e as capturas desta revisão estão na [seção 22 do quickstart](../001-gerador-jornada/quickstart.md#22-reimplementação-visual-do-legado--2026-10-24). Critérios de aceitação abaixo são requisitos a verificar, não resultados automaticamente aprovados. As seções de expansão preservam o escopo planejado; não descrevem capacidades disponíveis apenas por constarem nesta spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Usar o componente em seu fluxo principal (Priority: P1)

Como autor da jornada, quero personalizar conteúdo, geometria e aparência com validação antes de aplicar.

**Why this priority**: Essa é a responsabilidade principal deste componente; sua aparência deve apoiar uma operação verificável.
**Independent Test**: Preparar uma jornada mínima, exercitar apenas os três comportamentos abaixo e observar seus efeitos sem depender da criação manual de outros componentes.

**Acceptance Scenarios**:

1. **Given** uma jornada mínima válida, **When** o autor usa o componente, **Then** o sistema deve editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
2. **Given** o componente no estado de referência, **When** o autor configura seu comportamento, **Then** o sistema deve personalizar cor, fonte, preenchimento, borda e imagens quando suportados pelo componente.
3. **Given** conteúdo configurado, **When** o fluxo principal prossegue, **Then** o sistema deve aplicar todos os valores válidos em uma ação reversível e conservar o estado anterior em caso de erro.

---

### User Story 2 - Personalizar a aparência com previsibilidade (Priority: P2)

Como autor ou responsável pelo visual, quero ajustar a aparência deste componente sem alterar sua identidade nem suas regras.

**Why this priority**: Permite reproduzir o visual legado e evoluí-lo de forma controlada.
**Independent Test**: Usar o estado de referência do DESIGN.md, alterar uma propriedade visual por vez e comparar resultado e comportamento antes/depois.

**Acceptance Scenarios**:

1. **Given** o estado de referência, **When** uma propriedade visual permitida muda, **Then** a aparência correspondente muda e os três comportamentos da US1 permanecem válidos.
2. **Given** uma personalização fora dos limites, **When** o autor tenta confirmá-la, **Then** o autor recebe diagnóstico no contexto e o último estado válido permanece íntegro.

---

### User Story 3 - Preservar conteúdo e contexto (Priority: P2)

Como autor, quero trocar de seleção ou modo e voltar ao painel sem perda ou alteração implícita dos dados confirmados.

**Why this priority**: O mesmo trabalho precisa continuar coerente ao mudar de contexto.
**Independent Test**: Usar uma fixture com texto acentuado e geometria conhecida, executar o percurso abaixo e comparar os dados antes/depois.

**Acceptance Scenarios**:

1. **Given** uma configuração confirmada, **When** o autor executa o percurso de preservação, **Then** o sistema deve manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
2. **Given** uma operação cancelada ou recusada, **When** o autor retorna ao trabalho, **Then** nenhum dado confirmado foi perdido nem parcialmente substituído.

### Edge Cases

- Inteiro inválido, cor fora do formato, nome vazio, referências inválidas e texto com caracteres especiais. O Properties legado era um placeholder; o comportamento funcional estava em ToolsRight.
- Janela 1366 × 768 e DPI 100%, 150% e 200%; texto longo e navegação por teclado.
- Uma referência ausente deve identificar a origem do problema; não ser substituída silenciosamente por outra entidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Editar nome lógico, dimensões e propriedades específicas sem misturar dados de outra família.
- **FR-002**: Personalizar cor, fonte, preenchimento, borda e imagens quando suportados pelo componente.
- **FR-003**: Aplicar todos os valores válidos em uma ação reversível e conservar o estado anterior em caso de erro.
- **FR-004**: Permitir personalização visual conforme DESIGN.md, distinguindo conteúdo persistido de decoração transitória do editor.
- **FR-005**: Manter seleção, foco, zoom e rolagem fora do documento; salvar apenas alterações confirmadas nos dados da jornada.
- **FR-006**: Recusar alterações inválidas com diagnóstico contextual, sem publicar estado parcial; operações de edição confirmadas devem respeitar o histórico do projeto.
- **FR-007**: Manter os controles e mensagens em português; preservar texto Unicode e não depender de caminhos externos para recursos incorporados.

### Key Entities *(include if feature involves data)*

- **Configuração do componente**: Valores editáveis do modelo e rascunho de formulário separado; tipos de entrada apropriados a número, texto, opção e cor.
- **Estado de interação**: Não alterado; alterado; inválido; aplicado; seleção trocada.
- **Referências compartilhadas**: entidades do domínio responsável; não duplicar telas, catálogo, navegação ou roteiro nesta feature.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Todos os três cenários de uso principal passam com a fixture de referência e com uma personalização visual válida.
- **SC-002**: Nenhum texto, identificador ou geometria confirmada se perde no percurso de preservação; para representações de tela, diferença geométrica máxima de uma unidade lógica.
- **SC-003**: Todos os casos inválidos da fixture preservam o último estado válido e apresentam uma causa localizável.

## Assumptions

- Requisitos e IDs desta feature são locais; citar `028-propriedades/FR-xxx` ao referenciá-los em outra feature.
- Componentes de interface não são novos domínios de negócio. O domínio responsável continua sendo 001-gerador-jornada.
- O DESIGN.md é um contrato de personalização e revisão, não um arquivo de configuração lido automaticamente pelo executável.
- A entrega desta especificação não significa implementação integral ou homologação. Lacunas devem ser tratadas no planejamento da feature coordenadora.
- A saída Java continua condicionada ao perfil EFL real e ao ambiente Java 6; documentos e dublês não substituem a homologação. A interface do editor não é exportada como parte da tela da jornada.

## Atualização solicitada em 2026-10-18 — Confirmação de alterações

- **FR-008**: Manter botão Salvar alterações visível fora da rolagem do formulário, inclusive quando o último campo estiver fora da área visível.
- **FR-009**: Aplicar o rascunho válido em uma ação reversível e atualizar a tela imediatamente após salvar; Descartá-lo pelo botão Descartar restaura os valores confirmados.
- **FR-010**: Exibir Alinhamento do texto, com esquerda, centro e direita, em controles textuais. A alteração modifica texto dentro da caixa, sem mover o elemento.
- **FR-011**: Com seleção múltipla, mostrar quantidade e comandos coletivos; não exibir formulário de um único componente como se editasse todos.

**Clarificação de FR-008 — redimensionamento, 2026-10-21 (sem novo FR):** rodapé fixo não significa divisor travado pela largura do formulário. Manter cabeçalho e Salvar alterações/Descartar fora da rolagem, permitindo que somente o corpo role **horizontal e verticalmente**. Seletores longos, caminhos de fontes e versões de perfil não devem impor seu mínimo horizontal ao painel; conservar valores completos e acesso horizontal aos campos, inclusive com overflow. A largura mínima continua limitada naturalmente pelos controles fixos, sem largura fixa, botões estreitados ou colapso ilimitado. Rolar/redimensionar não reconstrói nem publica o rascunho e não altera Save/Discard, documento ou histórico.

**Aceitação da clarificação (pendente):** em janela 1280 × 720, testar propriedades de tela e das 16 famílias, além de seleção múltipla sem ações individuais fictícias. Medir deslocamento/largura efetivos do painel ao mover o divisor, não somente `Offset`; alcançar campos longos nos dois eixos, conservar rascunho e manter o rodapé imóvel. Salvar/Descartar devem responder ao hit-test dos botões renderizados após rolar, confirmando uma transação ou descartando sem publicar revisão, respectivamente. Quatro regressões preparadas em inspector_split_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_split_test.go`), não executadas por política Windows, sem contorno. **Evidência final do executor principal (2026-10-21):** gofmt em dois arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` concluiu a compilação de 15 suítes, **sem executá-las**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade PASS: CRC, todos os SHA256SUMS, sidecar SHA-256, PE amd64, ambas as cópias do EXE idênticas, README definitivo, manifesto com aprovações pendentes e ausência de backups terminados em `~`. Identidades finais e evidências no [quickstart, seção 11](../001-gerador-jornada/quickstart.md#11-redimensionamento-do-painel-de-propriedades--2026-10-21). [001/T129](../001-gerador-jornada/tasks.md#phase-23-redimensionamento-do-painel-de-propriedades--2026-10-21) permanece aberta como verificação mista: sem execução dos testes Go, validação de UI/arraste nativos, DPI, revisão manual ou aceite; `manualApproval` e `redistributionApproval` continuam `pending`.

### Aceitação da atualização

1. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-008, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
2. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-009, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
3. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-010, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.
4. **Given** uma fixture de referência e o estado inicial registrado, **When** o autor executa a operação descrita em FR-011, **Then** o resultado atende integralmente ao requisito e o estado anterior pode ser recuperado quando houver edição confirmada.

### Cenário concreto do rodapé

Selecionar um componente com formulário maior que o painel, rolar até o início e alterar o texto. Salvar alterações continua visível na base do painel; clicar nele publica uma revisão e atualiza o canvas. Alterar novamente e clicar Descartar restaura o último texto confirmado, sem nova revisão. Uma largura inválida não publica nenhum campo do rascunho.

## Campo ID comum

- **FR-012**: Exibir campo ID nas propriedades de tela e de todas as famílias, independente de Nome. Explicar nome-base sem extensão e limite de 40 caracteres; exemplificar PnlIdentificacao na tela. Remover acentos/espaços/símbolos imediatamente na edição; não truncar silenciosamente valor digitado acima do limite. Validar prefixo Pnl somente para telas, início em letra e nomes Windows reservados.
- **FR-013**: Salvar alterações deve confirmar FileID, validar colisão em seu contexto e produzir uma ação reversível, sem mudar UUID ou texto exibido do componente. Se inválido, manter o último estado do projeto e informar o motivo; Descartar retorna ao ID confirmado.


## Proteção de rascunhos no encerramento

- **FR-014**: Ao solicitar encerramento pelo sistema ou Sair (025-barra-ferramentas/FR-016), detectar valores do formulário atual diferentes de sua montagem, mesmo sem revisão de projeto criada. Não substitui callbacks de validação nem publica o rascunho. **Não** preserva o formulário e a sessão; **Sim** autoriza descartá-lo e encerrar. Reconstrução por Salvar/Descartar renova a referência de comparação. A preferência de tema não altera essa referência.

## Aparência e famílias pós-v1 — 2026-10-19

- **FR-015**: Expor família/variante/skin/versão, fonte/face real e overrides tipados suportados, com origem e valor efetivo versus herdado; ausência não equivale a zero/falso. Restaurar ao padrão remove override explicitamente, sem trocar bounds/IDs/FileIDs/conteúdo. Formulários específicos abrangem novos tipos, duas legendas, variantes Report e grupos/ParentID quando aplicável, sem JSON livre. **Sem skin explícita, o seletor mostra “Padrão do legado” e a cena resolve o default da família; “sem skin” não pode significar aparência genérica.** A oferta de skins é filtrada pelas definições operacionais com recursos presentes; escolha antiga não resolvida mostra diagnóstico e não é apagada.
- **FR-016**: Mostrar recursos/delegates/capacidades indisponíveis e diagnóstico de fidelidade bloqueada por campo/skin/estado; não apresentar variante incompleta como homologada nem fonte por nome como licença verificada. Fonte XML bold/italic declarada e face/efeito Java comprovados são informações distintas. Salvar/Descartar e rascunho continuam transacionais; UI/chrome não recebe skin da jornada.

**Clarificação de FR-016 — bugfix de Recursos / referência, 2026-10-21 (sem novo requisito):** a consulta deve usar o componente selecionado **no rascunho recebido**, mesmo ainda não armazenado, e a revisão exata de `SkinVersion`, não todos os avisos do documento. Abranger todas as `group.Definitions`, partes compostas e seus recursos em todos os estados, inclusive inativos. Diagnóstico de biblioteca com tipo/nome pertence somente à definição exata; compartilhar caminho de fonte/imagem não torna outra variante pertinente. Aviso sem definição pode ser associado pelo recurso efetivamente dependente. Conservar falhas globais de integridade do perfil/recursos (`PROFILE_INVALID`), faltas reais, licenças e portabilidade, sem fallback de perfil ou revisão.

Preferir o diagnóstico local detalhado, omitir duplicatas genéricas e não mostrar `TEXT_METRICS_REQUIRED` quando a face já foi medida com sucesso pelo desktop via SFNT. Deduplicar por código/tipo/nome/estado/recurso, nunca pela mensagem; estados distintos continuam distintos. Exibir texto completo, não editável, com quebra de linhas e rolagem vertical na área de referência de 160 unidades lógicas do [DESIGN](DESIGN.md#recursos--referência--correção-de-2026-10-21). Dica exata: **Avisos somente deste componente. Role a lista para consultar os detalhes.** Seletores de skin/fonte/versão atualizam os avisos do rascunho sem publicar documento, perfil, recursos, ParentID ou histórico. O isolamento limpa ParentID apenas na cópia local de diagnóstico, não valida ancestralidade nem altera o vínculo real.

**Aceitação da correção (pendente):** uma seleção saudável não recebe avisos de controles defeituosos na mesma/outra tela; selecionar o controle defeituoso conserva sua falha. Revisões corrente/arquivada com o mesmo caminho não misturam avisos; revisão desconhecida não recebe substituição. Conferir partes compostas/estados futuros, corrupção de recurso não usado ainda bloqueante, ausência de aviso de métrica já medida, deduplicação com mensagens diferentes e estados distintos, lista longa integralmente rolável e atualização pelos três seletores sem salvar. Integração e dois arquivos de regressões em [001/T128](../001-gerador-jornada/tasks.md#phase-22-correção-de-recursos--referência--2026-10-21). Evidência final do executor principal para esta correção: gofmt em quatro arquivos PASS; `go vet ./...` PASS; `go test -c -o ./build/compile-tests/ ./...` PASS, com 15 suítes **somente compiladas, não executadas**; quatro testes Python de empacotamento PASS; build Windows pela tarefa process do VS Code na raiz e empacotamento Python PASS. Integridade do pacote PASS (CRC, SHA256SUMS, sidecar SHA-256, PE amd64, README exato, inventário exato, ausência de backups terminados em `~` e ambas as cópias do EXE idênticas); identidades finais somente na [seção 10 do quickstart](../001-gerador-jornada/quickstart.md#10-correção-de-recursos--referência--2026-10-21). Testes Go não executados por política Windows, sem contorno; UI nativa, DPI, revisão manual e aceite não realizados. T128 permanece aberta como verificação mista; resultados técnicos não equivalem a homologação.

**Clarificação adicional de FR-016 — disponibilidade local versus portabilidade, 2026-10-21 (sem novo FR):** conforme análise visual registrada pelo executor principal, a captura fornecida pelo usuário foi efetivamente visualizada e contém `PROFILE_RESOURCE_NOT_EMBEDDED` para os caminhos lógicos /azul/botao_aux1_pressed.png e /azul/botao_pressed.png. Os dois PNGs já estão incorporados ao aplicativo, com bytes/hash/dimensões conferidos por Python; não são imagens perdidas. O sintoma observado é a apresentação de uma limitação de portabilidade do arquivo do projeto como se fosse indisponibilidade no runtime local. A correção anterior de escopo/rolagem não tratava essa distinção; sua evidência e seus hashes permanecem históricos.

Alterar **somente a apresentação** em `componentDiagnostics`: para cada `PROFILE_RESOURCE_NOT_EMBEDDED` com caminho, consultar `sceneResourceData(document, group.SkinVersion, path)` na revisão exata, com verificação de SHA-256, tamanho e MIME, nunca por homônimo ou fallback. Com bytes exatos resolvidos, exibir **uma linha por caminho**, independentemente de estados/origens repetidos: **Disponível no aplicativo · PATH (não incluído no arquivo do projeto)**. Havendo ao menos uma dessas linhas, acrescentar **uma única vez**: **Aviso de portabilidade, não de carregamento: os recursos acima dependem deste acervo exato do aplicativo. Não foram copiados para o projeto; disponibilidade não comprova licença de redistribuição nem fidelidade integral.** Se a resolução falhar, conservar código e localização, com a mensagem **Recurso indisponível nesta instalação e não incluído no projeto; nenhuma substituição foi aplicada. Portabilidade e licença de redistribuição não comprovadas.** Não ocultar `PROFILE_INVALID`, recursos realmente ausentes/bloqueados, falhas de carregamento ou avisos legais.

`visualprofile.Diagnostics`, `visualprofile.ComponentDiagnostics`, `Capture`, diagnósticos persistidos, hashes/licenças/status do perfil e validação estrita do PPT permanecem inalterados. Não copiar recursos para o projeto, conceder licença, mudar modelo/schema ou enfraquecer gates de exportação; disponibilidade local não é portabilidade autônoma nem homologação. Manter o escopo de seleção/rascunho/revisão e o layout rolável existentes.

**Aceitação adicional (pendente):** verificar uma linha por caminho e um aviso comum, builtin/captura/salvar-reabrir sem copiar bytes, estados inativos, revisão corrente/arquivada exata, rejeição de hash/tamanho/MIME divergentes, faltas/bloqueios/corrupção reais e avisos legais preservados; seletores não publicam rascunho nem alteram perfil/documento. Nove testes preparados em resource_availability_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/resource_availability_test.go`), incluindo o subteste de recurso já incorporado e validado com aviso persistido obsoleto, sem anunciar dependência externa inexistente. **Evidência final do executor principal (2026-10-21), após a última edição:** gofmt dos dois arquivos Go PASS; `go vet ./...` PASS (`RESOURCE_FINAL_VET_PASS`); `go test -c -o ./build/compile-tests/ ./...` PASS (`RESOURCE_FINAL_COMPILE_PASS`), com **15 suítes somente compiladas e ZERO execuções**; quatro testes Python de empacotamento PASS; tarefa VS Code de build Windows na raiz terminou com exit 0 e empacotamento Python sem `--check` concluiu o novo pacote. Revisão Explore somente leitura não apontou defeito concreto; não é teste executado.

CRC, SHA256SUMS/bytes locais, inventário exato, ausência de backups `~`, sidecar SHA-256, manifesto/build-info, produto/PE amd64, novo texto literal do aviso no EXE e README congelado conferidos. **Histórico preservado:** inicialmente somente o EXE canônico e o ZIP foram atualizados; a primeira cópia de compatibilidade foi bloqueada por estar em uso, conservando a revisão anterior. A checagem ad hoc inicial falhou por exigir incorretamente cópias idênticas naquele momento, não por falha de build/ZIP; a leitura independente posterior confirmou a distribuição canônica.

**Entrega de compatibilidade RESOLVIDA (2026-10-21):** após o usuário informar **“fechei a aplicação”**, o executor principal realizou `Copy-Item` do [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) para a [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) e `Get-FileHash` independentemente nos dois caminhos, com PASS: `COMPATIBILITY_COPY_UPDATED_AND_VERIFIED_PASS SHA256=48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d`. **Ambos agora têm 34.608.128 bytes e o mesmo SHA-256**; a divergência anterior é somente histórica. O ZIP **não foi refeito**: conserva 20.408.876 bytes, SHA-256 `b40cc96c901b51d00ca042831c419c222f44f90709b7d0328baf70695d8bdf01` e os checks de integridade já PASS. README congelado preservado; não há necessidade de reempacotar por esta atualização documental. Identidades atuais e roteiro na [seção 12 do quickstart](../001-gerador-jornada/quickstart.md#12-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21); evidências e hashes anteriores permanecem históricos.

[001/T130](../001-gerador-jornada/tasks.md#phase-24-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21) permanece aberta como **verificação mista**: 15 suítes Go somente compiladas, não executadas por política Windows, sem contorno; quatro testes Python e build/empacotamento PASS. Nenhum aplicativo ou teste Go foi iniciado nesta finalização, nenhum processo foi encerrado pelo agente e nenhuma política foi alterada. UI nativa, DPI 100/150/200%, Office e aceite pendentes. A captura analisada registra o sintoma anterior, não valida a UI nova; `manualApproval` e `redistributionApproval` continuam `pending`, sem alterar T120–T122 ou os demais gates.

**Aceitação**: override falso/zero permitido distingue herança; restaurar remove somente o campo escolhido, Undo restaura e round-trip mantém valores. Seleção/tema não publicam rascunho. Modelos em 003/T010–T015 e 004/T009–T010; integração única 001/T113, testes T118 e manuais/aceite T120–T122. Sem implementação duplicada nesta feature nem homologação nesta revisão.

## Alinhamento textual universal — 2026-10-20

- **FR-017**: Estender FR-010 a todas as 16 famílias capazes de apresentar texto, inclusive com conteúdo vazio: exibir sempre **Alinhamento do texto** após as propriedades de conteúdo e antes da aparência avançada, com **Padrão da skin**, **Esquerda**, **Centro** e **Direita**. Exibir dica de que altera somente o texto, sem mover ou redimensionar o componente; manter **Posição do elemento** separada dos controles textuais. Aplicar a escolha a cada região textual da família, incluindo legendas secundárias, células, áreas de relatório, teclas e popup de combo, inclusive **Sem seleção**. Carregar o valor confirmado ao selecionar/reabrir, confirmar por Salvar alterações em uma ação reversível e preservar Descartar/Undo/Redo, sem alterar bounds, pais/filhos, hit areas ou fundos.

**Aceitação adicional (pendente):** em cada família, selecionar conteúdo preenchido e vazio e conferir visibilidade/ordem/dica do campo; alternar as quatro opções, salvar, descartar outro rascunho, desfazer/refazer e salvar/reabrir. Conferir texto em canvas/prévia/PPT e abrir o popup do combo com e sem seleção. Geometria, hierarquia e áreas de interação permanecem iguais. Padrão da skin remove apenas o override textual; não altera a migração histórica de alinhamento vazio antigo para esquerda.

Integração e regressões em [001/T126](../001-gerador-jornada/tasks.md#phase-20-alinhamento-textual-universal--2026-10-20); invariantes em [004/FR-017](../004-tela/spec.md#alinhamento-textual-universal--2026-10-20). Após produção e quatro novos arquivos de testes, o executor principal informou em 20/10: gofmt e `go vet ./...` PASS, 15 suítes Go compiladas sem execução, quatro testes Python de empacotamento PASS e tarefa VS Code de compilação Windows PASS. Registro técnico e identidade do EXE na [seção 8 do quickstart](../001-gerador-jornada/quickstart.md#8-alinhamento-textual-universal--2026-10-20). Testes Go não executados pelo bloqueio prévio de política Windows, sem contorno nem PASS de execução; tarefas mistas, manuais e aceite permanecem abertos.

## Sobreposição de componentes — 2026-10-20

- **FR-018**: Na seleção individual das 16 famílias, oferecer **Sobreposição** com **Manter posição atual**, **À frente do objeto**, **Atrás do objeto**, **Trazer para a frente** e **Enviar para o fundo**. Iniciar em Manter posição atual, sem reordenar ao salvar outras propriedades; mostrar camada atual contando **1 = fundo**. Habilitar **Objeto de referência** somente nos modos relativos, sem seleção implícita: listar outros componentes da mesma tela, excluindo o próprio, com **ordinal · nome** para desambiguar nomes repetidos. A escolha é rascunho até **Salvar alterações** e integra a mesma transação atômica das demais propriedades; referência ausente/inválida ou qualquer outro campo inválido não publica estado parcial. Preservar Descartar, Undo/Redo e proteção de rascunho. Colocar imediatamente à frente/atrás da referência ou no extremo solicitado, sem vínculo permanente; os botões **Frente ↑ / Fundo ↓** continuam de um passo. Persistir somente a ordem de Components, sem novos campos/schema; não modificar bounds, estilo, conteúdo, ParentID ou camadas dos filhos por consequência da operação.

**Aceitação adicional (pendente):** salvar outra propriedade no modo padrão mantém a ordem; selecionar cada modo relativo com referência válida e conferir adjacência e estabilidade dos demais objetos. Sem referência válida, salvar recusa também as outras alterações do rascunho. Conferir extremos, tela com um único componente, referência excluída/desatualizada, nomes repetidos, Descartar, Undo/Redo e salvar/reabrir. Testar pai com filhos: somente o componente escolhido muda na lista plana, sem mover a subárvore. Modos/referência não são relações persistidas e voltam ao padrão neutro ao remontar o formulário.

Implementação em inspector_layer.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_layer.go`), integrada a [componentInspector](../../apps/gerador-jornada/internal/adapter/desktop/inspector.go); invariantes em [004/FR-018](../004-tela/spec.md#sobreposição-de-componentes--2026-10-20), integração [001/T127](../001-gerador-jornada/tasks.md#phase-21-sobreposição-de-componentes--2026-10-20). Produção e regressões implementadas. O executor principal informou gofmt dos sete arquivos e vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e build/pacote Windows com integridade PASS. [Evidências, roteiro e identidade — seção 9](../001-gerador-jornada/quickstart.md#9-sobreposição-de-componentes--2026-10-20). Testes Go não executados por política Windows, sem contorno; T127/004-T014, UI nativa, DPI, Office e aprovação manual permanecem abertos. Resultados anteriores são históricos, não PASS deste incremento.

## Log da consulta visível — clarificação de FR-016, 2026-10-21

> **Removido em 2026-10-23 ([001/T134](../001-gerador-jornada/tasks.md#t134-remocao-recursos-referencia)):** o campo **Recursos / referência** (FR-016 e clarificações de 2026-10-19/21) deixa de existir no inspector, assim como o registro automático e o botão de diagnóstico da prévia; Skin/Face da fonte/Versão continuam aplicando o rascunho sem publicar. **Registrar diagnóstico agora** captura o componente selecionado ou, sem seleção, os da tela atual (estado salvo). O texto a seguir é histórico, sem novo FR.

Sem novo FR local: consumir [001/FR-057](../001-gerador-jornada/spec.md#log-de-aplicação--2026-10-21). Registrar automaticamente a mesma consulta de **Recursos / referência** mostrada na seleção/rascunho e ao trocar Skin/Face da fonte/Versão do perfil, sem salvar. **Log de aplicação** no cabeçalho mostra o caminho efetivo e permite **Registrar diagnóstico agora**, **Copiar caminho** e **Abrir pasta**. Captura explícita força o rascunho visível, mesmo dentro da janela de dedup, ou os componentes da tela atual sem seleção; o diagnóstico da prévia também registra a consulta.

Conservar diagnóstico técnico original código/tipo/nome/estado/recurso, incluindo condições não transcritas pela apresentação amigável, com redação/hash de strings personalizadas. Mensagem builtin exata só quando permitida. `displayed` significa presença do código no relatório, não resolução; `report_sha256` correlaciona o texto efetivamente mostrado, **nunca o registra em bruto**. Recursos têm revisão exata, inclusive arquivada, partes/estados futuros, SHA/tamanho/MIME/status/licença esperados e hash/tamanho resolvidos ou código de erro. Não alterar disponibilidade, perfil, licenças, políticas/fallback, documento ou histórico. Privacidade, limites de armazenamento e erros consultáveis do logger pertencem a 001/FR-057, não a um segundo requisito/implementação em Propriedades.

**Estado final técnico — pronto para coleta (2026-10-21):** produção concluída e 11 arquivos Go formatados (sete existentes e quatro novos). Executor principal observou `APPLICATION_LOG_VET_PASS` e `APPLICATION_LOG_COMPILE_ONLY_PASS`: **15 suítes somente compiladas, ZERO execuções Go**, incluindo 12 testes novos (9 de núcleo + 3 de UI). Quatro testes Python de empacotamento PASS; tarefa Windows na raiz exit 0; pacote sem `--check` e integridade somente leitura PASS. Ambos os EXEs atualizados/idênticos; identidades reais e checks na [seção 13](../001-gerador-jornada/quickstart.md#13-log-de-aplicação--2026-10-21), sem repetir hashes aqui. README da aplicação congelado já empacotado; produto/schema/catálogo **1.0.0.0/3/2** inalterados.

**Coleta e aceite pendentes:** fechar a instância antiga e abrir o novo EXE identificado para reproduzir/enviar o log. A falha continua relatada, **causa não mapeada e solução não comprovada**. Nenhum aplicativo foi executado nesta entrega; presença de marcadores no EXE não prova startup/uso do logger. Startup/runtime, execução dos testes Go, UI nativa/DPI 100/150/200%, Office e aceite pendentes, política Windows sem contorno. Integração única [001/T131 aberta como verificação mista](../001-gerador-jornada/tasks.md#phase-25-log-de-aplicação--2026-10-21); T130/T120–T122 e aprovações `pending` preservados. Histórico acima intacto, sem herdar seus PASS/hashes; não criar logs artificiais nem afirmar que o usuário já abriu o novo EXE.

## Parametrização dos códigos de barras — planejamento 2026-10-23

- **FR-019**: Consumir propriedades tipadas de [031](../031-codigo-barras-efl/spec.md) e [032](../032-codigo-barras-processado-efl/spec.md). Na entrada, editar rótulo e propriedades comuns sem configurar DVs/quantidade de campos/valores runtime. No processado, editar blocos/títulos/linhas ordenados, combo por linha CPF/CNPJ/Nome/Nome fantasia/Razão social e valor literal. Trocar rótulo conserva valor; confirmar tudo em uma transação via Tela, Descartar restaura confirmado. Preservar rodapé fixo, rascunho e foco; seleção de bloco/popup não é dado do projeto nem controle da jornada.

**Aceitação**: duas linhas/blocos com tipos distintos são editáveis independentemente; alteração cancelada não muda modelo/canvas/prévia. Tipo inválido/remoção da última linha ou bloco informa contexto e não publica alteração parcial. Domínio de dados é [003/FR-015](../003-catalogo/spec.md#codigos-barras-2026-10-23), comandos [004/FR-019](../004-tela/spec.md#codigos-barras-2026-10-23); implementação do formulário somente em [001/T136](../001-gerador-jornada/tasks.md#tarefas-codigos-barras-2026-10-23), ainda planejada. Não reintroduzir o campo Recursos / referência ou gerador Java retirados por decisões anteriores.

## Autoria do Menu aplicação — 2026-10-23

- **FR-020**: Editar [033](../033-menu-aplicacao-efl/spec.md) por árvore de menus/submenus/itens com legenda, ordem, adicionar/remover, categoria, estilo/variante, habilitação por item, submenu ou destino de folha explicitamente escolhido e **Menu inicial** por identidade/caminho. Títulos repetidos não confundem seleção. Menus vazios são válidos; não configurar caminho runtime ou tratar metadados de movimento como comandos. Preservar formulário em rascunho, rodapé fixo, Salvar alterações/Descartar e uma transação reversível via Tela; limites/UUIDs/inicial removido/submenu+destino inválidos recusam confirmação sem perda parcial.

**Aceitação**: criar submenu aninhado, reordenar, desabilitar item, escolher destino e inicial, salvar/reabrir e desfazer/refazer; rascunho descartado não altera canvas. Campo Movement* é metadado inerte quando exibido, nunca logoff/rede/operação. Dados em 003/FR-016, comandos em 004/FR-020, destinos em 005/FR-009 e formulário único em [001/T141](../001-gerador-jornada/tasks.md#tarefas-menu-aplicacao-2026-10-23); não duplicar implementações nem checklists de revisão. Autorizado para implementação, não testado manualmente nesta documentação.

## Metadados dos itens de menu — 1.0.0.3

O editor de árvore apresenta Código do movimento, Nome do movimento, Opção do movimento e Atalho legado para a folha/ramo selecionado. Trocar seleção preserva o rascunho; confirmar publica atomicamente e Descartar não altera o original. Conservar texto literal, zeros e Unicode; distinguir normal/especial mesmo com código compartilhado. Esses campos não iniciam movimentos nem validam o cadastro XML ausente.
