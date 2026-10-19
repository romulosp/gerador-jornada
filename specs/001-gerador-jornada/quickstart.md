# Quickstart: validação do Gerador de Jornada

**Status**: implementação em andamento em apps/gerador-jornada; incremento atual do editor disponível para teste manual. A suíte existente está aprovada, mas não cobre integralmente todos os requisitos: lacunas de implementação/cobertura estão nos tasks.md de cada domínio. O registro vigente de evidências está no final deste documento; seções anteriores são histórico de revisões/builds. Os cenários manuais abaixo não são considerados executados por existir uma suíte automatizada.
**Input**: [spec.md](spec.md), [plan.md](plan.md), [tasks.md](tasks.md).

## 1. Revisar a documentação com Spec Kit

A feature coordenadora está registrada em `.specify/feature.json`. As features dos cinco domínios têm seus próprios spec.md, plan.md, data-model.md e tasks.md; não são produtos separados.

Na raiz do projeto, em PowerShell, o check existente pode verificar os artefatos:

```powershell
$env:PYTHONUTF8 = '1'
python .specify/scripts/python/check_prerequisites.py --json --require-spec --require-tasks --include-tasks
```

Para inspecionar outro domínio sem trocar a feature registrada:

```powershell
$env:SPECIFY_FEATURE_DIRECTORY = 'specs/002-projeto'
$env:SPECIFY_FEATURE_NO_PERSIST = '1'
python .specify/scripts/python/check_prerequisites.py --json --require-spec --require-tasks --include-tasks
Remove-Item Env:SPECIFY_FEATURE_DIRECTORY
Remove-Item Env:SPECIFY_FEATURE_NO_PERSIST
```

O check confirma presença e resolução dos documentos; não homologa requisitos nem executa a aplicação. Não rodar setup de uma feature para sobrescrever documentos preenchidos.

## 2. Ambiente futuro de desenvolvimento

Ferramentas locais preparadas: Go 1.27.2 amd64, Fyne 2.8.2 e MinGW GCC 16.2.0 UCRT. A partir da raiz do workspace, executar em PowerShell:

```powershell
Set-Location apps/gerador-jornada
. ./scripts/env.ps1
go test ./...
go vet ./...
go run ./cmd/gerador-jornada
```

scripts/check.ps1 executa verificações; scripts/build.ps1 compila o EXE e scripts/build.py empacota com metadados/hashes. Inventário completo de licenças da release ainda está em T067. Não recompilar Linux/macOS na entrega corrente. Testes que precisam de PowerPoint ou EFL devem identificar pré-requisitos e registrar ausência como não executados, sem aprovação falsa.

Toolchain Java SE 6 e JARs EFL ficam no ambiente de homologação da saída, não no editor distribuído. A biblioteca real e o bootstrap do hospedeiro são gates G-03/G-04 de research.md.

### 2.1. Validar Java 6 usando o JDK 17 informado

A ferramenta está implementada em `apps/gerador-jornada/scripts/java6.py`. Executar os comandos abaixo a partir de `apps/gerador-jornada`:

```powershell
python -X utf8 scripts/java6.py prepare
./scripts/check.ps1 -Jdk17 'C:\Desenvolvimento\jdk-17.0.11'
```

O primeiro comando baixa ferramentas em `.tools/java6`, com hashes fixados em `scripts/java6-tools.json`. O segundo executa testes Go, vet, testes positivos/negativos da toolchain e compilação dos fontes do gerador com perfil exclusivo de teste. Esse perfil não homologa EFL nem libera exportação Java no produto.

Quando os fontes e todos os JARs reais estiverem disponíveis, usar o comando abaixo substituindo os caminhos de exemplo. Repetir `--classpath` para cada dependência transitiva:

```powershell
python -X utf8 scripts/java6.py check 'D:\saida\fontes-jornada' --jdk-home 'C:\Desenvolvimento\jdk-17.0.11' --classpath 'D:\bibliotecas\efl.jar' --offline --report 'D:\saida\compatibilidade-java6.json'
```

O verificador usa ECJ com source/target 1.6, exige major 50 e verifica APIs SE 6 nos fontes compilados e JARs. Uma execução reprovada retorna código 1 e grava relatório reprovado, substituindo uma aprovação anterior. Não usar o `javac` 17 com `--release 6`, opção indisponível nesse JDK. A ferramenta não executa o programa na JVM 6: o teste funcional Java 6/EFL de V-07 permanece obrigatório. Não há JDK ou compilador embutido no EXE distribuído.

## 3. Conjunto de referência

Criar projeto Atendimento com telas A: Identificação, B: Confirmação e C: Conclusão, todas inicialmente 800 × 600.

- A: title, label, textField numérico, numericKeyboard vinculado ao campo e botão Continuar para B.
- B: report com duas linhas, radio com duas opções e botão Confirmar para C.
- C: label/title e botão Reiniciar para A.
- Tela inicial da jornada: A. Nomes lógicos diferentes por tela, textos com acentos e pelo menos uma imagem decorativa incorporada.
- Roteiro: iniciar em A, mover cursor ao campo, clicar, digitar 123, esperar, clicar em Continuar, navegar; esperar em B, clicar em Confirmar, navegar; esperar em C, clicar em Reiniciar e navegar para A.

Cada etapa tem tempo explícito e tela corrente correta. Duração sugerida: movimento 400 ms, clique 200 ms, digitação 900 ms, espera 500 ms e navigate 100 ms. Conferir que a revisita a A restaura valores iniciais.

## 4. Cenários de validação

### V-01 — edição e histórico (US1)

Inserir sete famílias por drop; mover/redimensionar/alinha-las; mudar textos e opções; alterar camadas; duplicar tela e componente. Conferir novos IDs e remapeamento de referências internas.

Aplicar zoom 25/100/200/400%; conferir que bounds persistidos não mudam. Undo de arraste deve ser uma ação; novo comando após undo remove redo. Testar 100 entradas, foco de campo, Escape e atalhos.

**Expected**: geometria coerente, sem alteração oculta de dados, e operações desfeitas/refeitas integralmente.

### V-02 — arquivo, falha e recuperação (US1)

Salvar/reabrir o conjunto, comparar manifesto e recursos. Forçar interrupção após autosave, recuperar e salvar nova versão. Testar manifesto truncado, versão desconhecida, hash errado, path traversal e limites ZIP.

Simular disco cheio, arquivo bloqueado e falha antes de substituir. Alterar o projeto durante save e verificar que a revisão mais nova segue marcada como alterada.

**Expected**: último arquivo válido intacto; arquivo ilegível não substitui sessão; erro localizado; referências inválidas preserváveis como rascunho.

### V-03 — jornada e prévia (US2)

Percorrer A → B → C → A. Preencher e selecionar, usar teclado numérico e reiniciar. Remover B após revisar impacto e tentar executar; desfazer a remoção.

**Expected**: navegação correta e valores iniciais ao reentrar; prévia não altera projeto; referência quebrada localizável; undo restaura fluxo.

### V-04 — roteiro explícito (US4)

Criar, duplicar e reordenar roteiros/etapas. Reproduzir, pausar, continuar e reiniciar. Alterar um destino de botão e verificar etapas subsequentes. Testar alvo ausente, typeText vazio, Unicode e duração fora de 1..600000 ms.

**Expected**: clique apenas visual, navigate separado, sequência determinística, sem consumo de tempo durante pausa nem efeitos residuais após restart.

### V-05 — PPT exploratório (US3)

Exportar somente exploração em máquina sem Office. Abrir o arquivo nas máquinas de homologação PowerPoint 2016 e Microsoft 365, sem reparar. Usar modo apresentação para testar botões e menu.

No modo de edição, selecionar textos/formas internos de cada família, alterar um rótulo e salvar/reabrir. Comparar posição e tamanho com tolerância de 1 unidade lógica; verificar ordem, legibilidade e recursos.

**Expected**: todos os destinos corretos, sem avanço automático, controles editáveis e nenhum screenshot substituindo a tela.

### V-06 — PPT demonstrativo (US4)

Exportar exploração e roteiro. Entrar no roteiro pelo menu; verificar movimento de cursor, realce, texto progressivo, waits e revisita a A. Conferir que apenas um prefixo de texto está visível por instante.

Medir duração das etapas, excluindo tempo de seleção do menu e os 1000 ms finais de permanência previstos no contrato. Aplicar tolerância SC-004. Testar replay e depois exploração manual, sem avanço temporal herdado.

**Expected**: execução consistente nas duas versões do PowerPoint, sem macros, dados interativos falsos ou animações omitidas.

### V-07 — fontes Java/EFL (US5)

Somente com perfil homologado: gerar para diretório novo, inspecionar as sete famílias, Javadoc, listeners e separação por tela. Gerar novamente sem alterações e comparar conteúdo.

Usar toolchain e comandos registrados pelo perfil para compilar contra API SE 6 e EFL real; verificar classfile major 50. Executar no hospedeiro e percorrer A → B → C → A na EDT, com campos, seleção e teclado.

Testar colisões pós-normalização, palavra reservada, pacote inválido, caracteres especiais e estilo sem capacidade EFL.

**Expected**: fontes válidos e integração funcional. Sem JARs/perfil, registrar não executado e verificar diagnóstico JAVA_PROFILE_UNVERIFIED; não contar dublês como homologação.

### V-08 — distribuição e cancelamento (US6)

Copiar distribuição para Windows 10/11 x64 limpos, com usuário comum e rede desligada; sem Go/Node/Java/Office/WebView adicional instalado como requisito. Abrir EXE em janela própria; salvar/reabrir e gerar PPT/fontes.

Testar DPI 100/150/200%, paths com espaços/acentos, diretório sem permissão e exportação cancelada. Durante exportação, editar o projeto e verificar revisão do snapshot. Tentar segundo job.

**Expected**: dependências efetivas registradas, edição utilizável, um job por vez, cancelamento em até 1 s no conjunto nominal e nenhuma publicação parcial.

### V-09 — integração, documentação e desempenho

Executar fluxo completo do conjunto de referência e comparar conteúdo/navegação das duas saídas. Medir conjunto nominal de 10 × 50 componentes e 50 etapas conforme spec.md.

Revisar dependências entre pacotes e documentação de elementos públicos/privados, campos, variáveis de pacote e locais; verificar Javadoc gerado. Registrar quais itens tiveram check automático e quais revisão humana.

**Expected**: SC-001 a SC-007 comprovados, com falhas e limites explicitados. Não exigir identidade de rasterização entre Fyne/PPT/Swing.

## 5. Registro de evidência e conclusão

Para cada cenário registrar data, versão do editor, revisão do projeto, versão do catálogo/perfil, toolchain, SO, DPI, PowerPoint/JARs quando aplicável, resultado e localização do artefato de evidência. Textos de projeto não devem ser copiados integralmente para logs.

Resultados possíveis: aprovado, reprovado ou não executado com causa. Campo vazio, requisito documentado ou teste com dublê não é aprovação. Atualizar tarefas apenas quando existir a evidência correspondente.

### Evidência da correção Java 6 — 2026-10-17

**Estado do incremento**: verificação estática concluída; integração Java/EFL bloqueada por ausência de bibliotecas e exemplo do hospedeiro. Feature completa permanece em_andamento, com teste manual e aprovação pendentes.

**Executor**: agente, ambiente Windows do workspace, JDK `C:\Desenvolvimento\jdk-17.0.11`. Comando executado: `./scripts/check.ps1 -Jdk17 'C:\Desenvolvimento\jdk-17.0.11'`.

| Verificação | Resultado observado |
| --- | --- |
| gofmt e `go vet ./...` | Aprovados |
| `go test ./...`, com integração Java habilitada | Aprovado; fontes de três telas/sete famílias compilados para major 50 com perfil estrutural de teste |
| Nove testes reais ECJ/Animal Sniffer | Aprovados: Swing/Unicode, dependência compatível, rejeição de lambda, API moderna, método moderno, API moderna dentro de JAR, bytecode moderno, caminho externo e relatório obsoleto |
| Integridade das ferramentas | SHA-256 conferidos contra `scripts/java6-tools.json` |
| Perfil EFL real e execução Java 6/EFL | Não executados; JARs e bootstrap não fornecidos, `JAVA_PROFILE_UNVERIFIED` preservado |

Esta evidência cobre T077 e a parte estática do gate G-04. Não representa execução dos cenários manuais V-01 a V-09 nem aceite da aplicação.

Build desktop recompilado após a correção com `scripts/build.ps1`: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `FE310C287190B337DCD566633C1009E0C441D4C2988B41E9E9558ACF5C5B1139`. O build não inclui as ferramentas de validação Java e conserva o diagnóstico de perfil EFL não homologado.

### Incremento visual e specs individuais — 2026-10-18

**Build atual para revisão**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `690934461D9E1543F4CD63FB35B5AA0F795119B57E98ACAD7E7A4192BE58C4CB`. Substitui a revisão de 2026-10-17 para os cenários de interface afetados.

**Implementado neste incremento**: fundo original incorporado nos novos projetos/telas, botão Fundo para aplicar/trocar/remover imagem em um projeto existente, paleta com rolagem independente, abas Telas/Paleta/Estrutura, seleção por estrutura, conjuntos editáveis de data/hora, lotérica, infraestrutura, cabeçalho e faixa de título. O zoom inicial se ajusta à janela, e sua indicação mostra o valor efetivo.

**Evidência automatizada**: `go test ./...` e `go vet ./...` passaram. Testes verificam os bytes originais em um projeto salvo/reaberto e no PPT, arraste sobre o fundo em zoom 50%, desfazer a movimentação e o cabeçalho como ação única, além da imagem na prévia. São verificações de integração e de pacote; não homologam a apresentação no PowerPoint.

**Evidência visual pelo agente**: janela nativa inicializada pelo comando `--smoke-screenshot`, captura [editor-fundo-legado.png](../../apps/gerador-jornada/build/editor-fundo-legado.png) inspecionada com fundo original e sete famílias acessíveis. Não é execução dos cenários manuais pelo usuário.

**Exemplo portátil**: [jornada-com-fundo.jornada](../../apps/gerador-jornada/dist/jornada-com-fundo.jornada). A imagem acompanha o arquivo; abrir projetos anteriores conserva o conteúdo deles.

**Reteste manual solicitado para esta revisão**:

1. Abrir o EXE e conferir o fundo; em um projeto anterior sem imagem, clicar **Fundo → Usar fundo Loterias CAIXA**.
2. Na aba Paleta, arrastar cada família para a imagem; mover, redimensionar e desfazer. Rolar a mesma aba para os conjuntos auxiliares e testar sua inclusão.
3. Alternar Telas/Paleta/Estrutura, selecionar um componente coberto pela estrutura e conferir as propriedades.
4. Salvar/reabrir, abrir a prévia e gerar PPT; conferir a imagem e os componentes nas três representações. Homologação do PPT nas versões-alvo continua pendente.
5. Revisar os contratos individuais no [inventário de specs e designs](research.md#inventário-de-componentes-e-designs-individuais), distinguindo o alvo visual de lacunas ainda documentadas.

**Estado**: este incremento de editor/fundo está pronto_para_teste_manual. Feature completa permanece em_andamento; EFL real, demais diferenças do legado e aceite humano continuam pendentes. As 24 specs novas têm estrutura e links conferidos, sem representar aprovação funcional.

### Revisão de edição, PPT e ícone — 2026-10-18

**Build mais recente para revisão**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `767F5716F4C6DCF1C53CEA68B38FD6DD9FEA79813921690CC878B5B2CD2E1E35`. Esta revisão substitui os hashes anteriores para novos testes. Compilada com `scripts/build.ps1`, incluindo o ícone fornecido pelo usuário na janela e nos recursos PE do EXE.

**Verificações realizadas**: `go test ./... -count=1` e `go vet ./...` passaram após a incorporação do ícone. Regressões cobrem seleção/movimento coletivo, limites de grupo, duplicação com referências remapeadas, remoção e Undo, formulário com confirmação e descarte, alinhamento sem alteração de geometria, persistência e transições PPT. Inspeção estática do EXE confirmou um grupo de ícones com sete imagens; o PNG incorporado tem o mesmo SHA-256 do original.

**Desempenho observado**: `BenchmarkSelect500`, dez iterações com 500 rótulos e driver de teste Fyne, Intel Core i5-13600KF, Go 1.27.2/windows/amd64: seleção passou de 14.251.660 ns/op, 1.776.739 B/op e 27.175 alocações/op para 704.540 ns/op, 483.956 B/op e 9.004 alocações/op. A medição compara o mesmo cenário; não comprova latência da janela nativa ou metas completas de desempenho.

**Evidência anterior desta revisão de interface**: [revisao-edicao.png](../../apps/gerador-jornada/build/revisao-edicao.png) mostra o rodapé fixo e Java/EFL desabilitado, antes da inclusão do novo ícone. `jornada-fade-revisao-2026-10-18.pptx` (artefato histórico, não está mais presente no dist atual) foi registrado como gerado antes da troca do ícone e contém nove slides com fade e oito movimentos de cursor, conferidos no ZIP/XML. Não houve homologação visual no PowerPoint.

**Limite de execução**: a revisão automática bloqueou a tentativa de executar o EXE para nova captura e regeneração da amostra, com a mensagem `blocked by policy`, sem detalhar a causa. Não foi tentado outro mecanismo de execução. O EXE final com o novo ícone foi compilado e inspecionado estaticamente; sua inicialização e aparência precisam do teste manual abaixo.

**Roteiro manual desta revisão**:

1. Abrir o EXE acima; conferir o ícone no Explorer, na janela e na barra de tarefas, também com DPI 150/200%. Copiar apenas o EXE para outra pasta e verificar o funcionamento offline.
2. Selecionar componentes com Ctrl/Shift+clique, depois com caixa de seleção em área vazia. Arrastar o grupo, usar setas (1 unidade) e Shift+setas (10), Ctrl+A, Ctrl+D, Delete, Escape e desfazer/refazer; conferir limites e distâncias relativas. Com um campo de texto focado, as setas devem editar o texto.
3. Conferir cursor de movimento sobre componentes e de redimensionamento na alça da seleção individual. Selecionar repetidamente e verificar resposta com muitos componentes.
4. Alterar propriedades, conferir **Salvar alterações** e **Descartar** sempre acessíveis no rodapé; salvar deve renderizar o novo estado e permitir Undo. Descartar conserva o estado anterior.
5. Inserir uma faixa de título e alterar **Alinhamento do texto** para esquerda/centro/direita; a posição do componente deve permanecer igual. Salvar/reabrir, prévia e PPT devem preservar a escolha.
6. Conferir **Java / EFL** visível e desabilitado. A suspensão é temporária e não remove a implementação.
7. Exportar PPT com Fade suave, Deslizar, Revelar e Sem transição; executar no PowerPoint desktop 2016/Microsoft 365, conferir efeitos e texto editável, links e roteiro com cursor/clique/digitação. Exploração deve manter navegação por ações, sem avanço automático. Registrar versão do Office e resultado.

**Estado**: T080–T087 têm implementação/evidências automáticas registradas. Teste manual e aprovação humana permanecem pendentes, assim como integração EFL real e capacidades futuras das specs individuais.

### Menu contextual e exclusão — 2026-10-18

**Build atual**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `FACB56BF41E62D5C62B22796B67F9EF5E4FEE0299AB6E4E5EA4DF29D11F24519`. Substitui a revisão anterior para teste. Mantém o ícone incorporado e as correções anteriores.

**Evidências**: `go test ./... -count=1`, `go vet ./...` e `scripts/build.ps1` passaram. `component_menu_test.go` exercita menu real do widget, foco em propriedades, Delete com foco textual, confirmação/cancelamento por menu e botão, Delete no grupo, confirmação única, Undo, clique vazio e troca de tela. O teste usa o driver Fyne de testes; não é execução do EXE no Windows nativo.

**Reteste manual**:

1. Sem seleção, conferir **Excluir** desabilitado imediatamente à direita de **Ajustar**. Selecionar um componente habilita; clique vazio, Escape e troca de tela sem seleção desabilitam.
2. Clicar com botão direito em um componente, inclusive com zoom diferente de 100%, e escolher **Alterar propriedades**. O painel direito deve mostrar somente o alvo e focar o formulário; salvar a alteração continua explícito.
3. Selecionar o componente e apertar **Delete**. Conferir pergunta e respostas **Sim/Não**. Não conserva o elemento; Sim remove, limpa seleção e desabilita Excluir. Desfazer deve restaurá-lo.
4. Repetir o cenário pelo botão **Excluir** e pelo item **Excluir elemento** do menu. Nenhuma entrada pode apagar sem confirmação.
5. Selecionar dois componentes e clicar com botão direito num deles: o conjunto permanece selecionado e **Excluir seleção** confirma a quantidade antes de removê-los em uma única ação reversível.
6. Com um campo de propriedades focado, usar Delete para editar texto; não deve aparecer exclusão de componente. Clicar com botão direito em área vazia não abre menu de elemento.

**Limites**: o bloqueio anterior de execução automática do EXE (`blocked by policy`, sem causa detalhada) não foi contornado. Esta revisão foi compilada e testada no driver de testes; teste visual nativo e aprovação humana permanecem pendentes.

### Versão 1.0.0.0 — 2026-10-18

EXE atualizado: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `53BEA5BAE15D09E4BF2A6CDBB28D2C447C1E9BF5A4E750B7E100F7E322636447`. Fonte única em `internal/adapter/buildinfo/VERSION`; cabeçalho/título e recurso VERSIONINFO usam o mesmo valor. `go test ./... -count=1`, `go vet ./...` e build passaram. Leitura estática com `System.Diagnostics.FileVersionInfo` confirmou FileVersion e ProductVersion iguais a `1.0.0.0`.

Teste manual: abrir Propriedades → Detalhes do EXE e conferir as duas versões; iniciar o aplicativo e conferir cabeçalho/título. Esta revisão não foi executada para inspeção visual; permanece o limite de execução automática já registrado, sem antecipar aprovação humana.

### Prévia Tela Atual — 2026-10-18

**EXE atualizado (versão 1.0.0.0)**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `C7ACDAEB961611C72615A690DE0C85A5E99D10AC2BC2A5634C19660690990472`. Encerrar a janela antiga e abrir este EXE para carregar os novos botões.

**Verificações**: `go test ./... -count=1`, `go vet ./...` e `scripts/build.ps1` passaram. `internal/application/screen_preview_test.go` comprova escopo de uma tela diferente da inicial, bloqueio da navegação entre telas, isolamento do documento/valores, validação da tela local e independência de outras telas/jornada incompletas. `internal/adapter/desktop/screen_preview_test.go` aciona o botão real no driver de testes, verifica o título da janela, preservação da prévia completa e desabilitação sem tela atual. São evidências automatizadas, sem homologação visual do EXE.

**Roteiro manual**:

1. Selecionar uma tela diferente da inicial e clicar **Prévia Tela Atual**, ao lado de **Prévia**. Conferir nome da tela no título, fundo, componentes e dimensões.
2. Operar campos, escolhas e teclado; clicar em um botão de navegação e conferir que continua na mesma tela.
3. Fechar a janela e confirmar que a tela inicial, seleção, valores salvos e histórico não foram modificados. Reabrir a prévia deve usar os valores confirmados no modelo.
4. Abrir **Prévia** e confirmar que o projeto completo continua começando na tela inicial e navegando normalmente.
5. Conferir o cabeçalho na resolução/DPI usados, incluindo as duas opções de prévia. Para ver alterações ainda em rascunho no painel, primeiro clicar **Salvar alterações**.

**Estado**: implementação de T092 concluída; teste visual nativo e aprovação do usuário continuam pendentes. A restrição anterior de execução automática permanece registrada; não foi utilizado outro mecanismo para contorná-la.

### ID de arquivo nas telas e componentes — 2026-10-18

**EXE atual (1.0.0.0)**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), SHA-256 `7778A00CDC21078CEBF3DBB8FF0F61514355B49684E030343E82F196188827F2`. Fechar a janela antiga antes de abrir esta revisão.

**Implementado**: campo ID comum em todas as famílias e nas telas. Tela `01 · Identificação` sugere `PnlIdentificacao`; componente aceita, por exemplo, `BtnContinuar`, sem Pnl obrigatório. Limite de 40 inclui prefixo e exclui extensão; não permite acentos, espaços, símbolos, começo numérico ou dispositivos Windows. Limpeza de acentos/símbolos é imediata; tamanho/prefixo/colisão inválidos impedem salvar sem truncar explicitamente o valor digitado. Renomear Nome não altera ID já confirmado, e alterar ID nunca muda UUIDs ou vínculos. Duplicação e migração sugerem sufixos livres dentro do limite.

**Persistência**: novos salvamentos gravam schemaVersion 2, com fileId separado de id UUID. Projetos v1 abrem por migração determinística e conservam referências; nenhum arquivo original é regravado durante a leitura. V2 incompleta/duplicada/inválida e versões desconhecidas são recusadas. catalogVersion e versão do software permanecem 1 e 1.0.0.0.

**Evidências**: `scripts/check.ps1 -Jdk17 'C:\Desenvolvimento\jdk-17.0.11'` e `scripts/build.ps1` passaram. Check inclui gofmt, go vet, testes Go, nove testes reais de ECJ/Animal Sniffer e compilação dos fontes estruturais do gerador para Java 1.6/classfile 50, com arquivos individuais nomeados por FileID. Novos testes cobrem o exemplo pedido, 40/41 caracteres, prefixo, nomes Windows, colisões de caixa, transação/Undo, duplicação no limite, migração v1, round-trip v2 e formulário. Gate real EFL preservado; botão continua visível e desabilitado.

**Reteste manual**:

1. Selecionar a tela Identificação sem componente selecionado e localizar ID nas propriedades. Conferir PnlIdentificacao e manter Nome com acentos.
2. Editar para `PnlConfirmação Final!`: campo deve mostrar PnlConfirmacaoFinal. Salvar alterações, mudar de tela e voltar, salvar/reabrir o projeto e conferir preservação.
3. Selecionar cada família, conferir ID comum e atribuir nome livre, por exemplo BtnContinuar. Confirmar que textos/legendas e a geometria não mudam e a navegação continua funcionando.
4. Tentar ID de 41 caracteres, ID de tela sem Pnl, Pnl sozinho, CON e um ID já usado no mesmo contexto. Salvar deve informar o erro e conservar o último estado válido. ID válido de 40 deve ser aceito.
5. Duplicar componente/tela com ID de 40 e conferir sugestão independente dentro do limite; desfazer deve restaurar o estado. Renomear um Nome amigável não renomeia automaticamente um ID confirmado.
6. Abrir um `.jornada` anterior e conferir sugestões e vínculos; salvar/reabrir no editor atualizado. O manifesto resultante usa formato 2 e é incompatível com leitores antigos que só aceitam 1.

**Estado**: T093–T096 implementadas e verificadas automaticamente. Interface nativa, teste manual e aprovação do usuário continuam pendentes, com o limite de execução automática anteriormente registrado. Integração EFL real continua independente e não homologada.

### README e distribuições de 64 bits — 2026-10-18

**Documentação**: [README completo da aplicação](../../apps/gerador-jornada/README.md) e [portal na raiz](../../README.md). Guia cobre propósito/público, sete famílias, cinco conjuntos, responsabilidades do editor, propriedades/IDs, atalhos, formatos/saídas, versões, arquitetura por domínio, requisitos de uso/build, comandos por plataforma, limitações e contribuição. Links locais conferidos.

**Pacotes gerados**:

| Alvo | Arquivo | SHA-256 do executável | Evidência |
| --- | --- | --- | --- |
| Windows amd64 | [ZIP](../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip), [EXE atual](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | `051E243F55548E18C1A539D7307B197D233F22CA86049FF51D5536E22C07A88E` | Build nativo, checks Go/vet e PE x64 conferidos |
| Linux amd64 | [TAR.GZ](../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-linux-amd64.tar.gz) | `6785223DEF2419721890E2C09E0EB279E3E6EEDF0C67B5E1B3B935286B9EDB2B` | Go/vet e suite passaram no contêiner Debian 12, ELF64 x86-64 e dependências conferidos |
| macOS amd64/arm64 | Workflow/bundle configurados; pacotes ainda não gerados | Não disponível | Sem Mac/SDK Apple disponível nesta estação; fixtures de cabeçalho não comprovam build Apple |

Versão do software mantida em 1.0.0.0. Windows: Go 1.27.2, GCC/MinGW 16.2.0; Linux: Go 1.27.2, GCC 12.2.0, glibc 2.36. O usuário iniciou o Docker para o build Linux. O empacotamento não abriu a interface nem homologou os sistemas-alvo.

**Comandos realizados**: `python scripts/build.py --target windows --arch amd64 --check` executou os checks e compilou; o caminho legado estava em uso. O empacotador foi ajustado para preservar o processo/arquivo antigo e entregar o pacote novo. Reexecução de `python scripts/build.py --target windows --arch amd64` gerou ZIP e inventário. `python scripts/build.py --target linux --arch amd64 --docker --check` executou vet/testes e gerou Linux; reempacotamento final incluiu dependências de runtime e README atualizado. Testes Python do reconhecimento PE/ELF/Mach-O passaram; hashes individuais e de ZIP/TAR, cabeçalhos reais e bit executável do TAR foram conferidos. Workflow YAML validado com quatro alvos, SHA de ações e gatilho manual.

**Caminho correto desta revisão Windows**: usar `dist/windows-amd64/JornadaStudio.exe`. O `dist/JornadaStudio.exe` antigo estava aberto e foi conservado; não afirmar que sua revisão é igual à do pacote atual. Nenhum processo do usuário foi encerrado.

**Recuperação**: adaptador usa caminho por usuário em Windows/Linux/macOS, preservando LOCALAPPDATA no Windows e evitando diretório relativo no Linux/macOS. O contêiner não substitui teste de uma sessão gráfica Linux.

**Roteiro manual**:

1. Extrair o pacote Windows, abrir o EXE do pacote, conferir versão/ícone, criar/editar/salvar/reabrir, prévias e exportação PPT. Conferir recuperação na pasta do usuário.
2. Em desktop Linux x64 com bibliotecas compatíveis, extrair o TAR e abrir o executável, repetir o fluxo e registrar distro, sessão gráfica, driver e DPI. Consultar runtime-libraries.txt/build-info.json; não contar os testes headless como aprovação da janela.
3. Em Mac/runner Intel e ARM64, executar os comandos do README ou o workflow manual, gerar e conferir os bundles, assinar/notarizar conforme a estratégia futura de distribuição e registrar evidências. Nenhuma release remota foi publicada nem workflow remoto executado nesta entrega.
4. Registrar aceite humano e resultados por plataforma separadamente. EFL real e PowerPoint desktop continuam gates próprios.

**Estado**: T097–T099 concluídas no escopo descrito. T100 permanece pendente para builds/homologação macOS e aprovação em máquinas limpas; disponibilizar configuração e compilar Windows/Linux não conclui essa validação.


### Maximização em tela cheia e entrega Windows — 2026-10-18

**Escopo atualizado pelo usuário**: recompilar somente o EXE Windows. Configuração e instruções Linux/macOS permanecem como base futura; builds e homologação desses alvos não são exigidos nesta entrega. T100 foi adiada, sem declarar execução ou aprovação.

**EXE 1.0.0.0 atualizado**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), também disponível em dist/windows-amd64. Ambos possuem SHA-256 `B839E989BEA8B3AF1F8D6A95714614D51D3409A0E47E750A45325A28710F6FCF`. A cópia principal foi atualizada nesta revisão.

**Comportamento implementado**: Maximizar no Windows entra em tela cheia pela API Fyne, mantendo os controles do editor e ocultando a decoração da janela/barra de tarefas. Botão Tela cheia/Restaurar e F12 alternam; Escape retorna. Estado de janela fica no adaptador desktop, fora do domínio e do manifesto. Não confirma nem cancela alterações de propriedades ou diálogos.

**Evidências**: testes Go e vet passaram, incluindo alternância, tradução de mensagens Windows, repetição de teclas, foco, rascunho, seleção, zoom, documento/histórico preservados e ações após fechamento. Os testes de mensagens são unitários, não exercitam o procedimento nativo do HWND. Build Windows com checks passou e o EXE foi recompilado após a definição do escopo. Uma verificação Linux anterior já havia sido concluída; não é requisito repetir esse build.

**Roteiro manual Windows**:

1. Encerrar a janela antiga, abrir o EXE atualizado e usar Maximizar ou duplo clique na barra de título. Conferir ausência da barra de título e da barra de tarefas, com ferramentas do editor acessíveis.
2. Retornar por Restaurar, F12 e Escape; repetir com uma propriedade em edição e componentes selecionados, conferindo rascunho, seleção, zoom e histórico preservados.
3. Abrir um diálogo, entrar em tela cheia e pressionar Escape: apenas sair de tela cheia, sem confirmar nem cancelar a operação do diálogo.
4. Conferir minimizar, fechar e Alt+Tab, além dos monitores e escalas DPI utilizados. Não considerar testes headless como homologação dessas interações nativas.

**Estado**: T101–T102 implementadas e verificadas automaticamente; teste visual e aprovação humana continuam pendentes. A execução automática do EXE foi bloqueada anteriormente com “blocked by policy”; não foi usado outro mecanismo para contornar a restrição.


### Tela cheia por padrão — 2026-10-18

O editor inicia em tela cheia em cada execução, com Restaurar visível. F12, Escape e Restaurar continuam permitindo voltar à janela. Estado visual transitório, sem alteração no domínio ou formato .jornada.

**Validação**: `python scripts/build.py --target windows --arch amd64 --check` passou (formatação, vet, testes Go, testes do empacotador e build Windows). O teste do adaptador verifica modo inicial, indicador atômico, legenda e restauração, além da preservação de edição já coberta. Não foram executados novos builds Linux/macOS.

**Executável desta revisão**: [dist/windows-amd64/JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), SHA-256 `39FD6B46E7545B65E9C55943ABBA9E202F204DCAF741B9C277DE8B6747FCA4C5`. A cópia dist/JornadaStudio.exe estava em uso e foi preservada; usar o caminho windows-amd64 para testar a alteração atual. Nenhum processo do usuário foi encerrado.

**Teste manual pendente**: fechar a janela antiga, abrir o EXE desta revisão e conferir tela cheia desde a abertura, controles acessíveis e Restaurar visível. Sair por Escape/F12/Restaurar, reabrir e conferir o padrão novamente. Os testes automatizados não homologam visualmente as barras do Windows.


### Correção da sequência de abertura em tela cheia — 2026-10-18

O usuário informou que ainda via as barras na abertura. A entrada em tela cheia foi movida da construção do controlador para ShowWindow: mostrar a janela nativa, instalar captura, solicitar tela cheia pela mesma operação usada por F12. A montagem e o dimensionamento normal ocorrem antes dessa sequência. O teste existente agora exige estado normal antes de ShowWindow e tela cheia/Restaurar após a abertura, mantendo os testes de retorno e edição.

`python scripts/build.py --target windows --arch amd64 --check` passou: formatação, vet, testes e build Windows. Os dois caminhos dist/JornadaStudio.exe e dist/windows-amd64/JornadaStudio.exe foram atualizados, SHA-256 `16749AD24B6D11B8BE966667E8984B58293EB114E49218DC582A397829F4F6F3`. Não houve novos builds Linux/macOS.

Teste manual pendente: encerrar a instância antiga, abrir este EXE e conferir que a tela inteira é ocupada, sem decoração/barra de tarefas; testar F12, Escape e Restaurar. Testes headless não comprovam o efeito visual no Windows. A causa visual não foi reproduzida em uma sessão nativa automática; a ordem de inicialização foi corrigida e a homologação continua humana.


### Sair em tela cheia e temas — 2026-10-18

**Implementado**: Sair no cabeçalho somente em tela cheia; confirmar Sim/Não quando houver documento alterado ou rascunho ainda não aplicado. Botão e fechamento nativo compartilham o fluxo, sem confirmações duplicadas. A janela principal encerra o aplicativo mesmo com prévias abertas. Tema no topo oferece Claro (padrão/fallback), Escuro e Azul, com preferência local salva por usuário. Cores do projeto e controles runtime de prévia ficam isolados da paleta do editor.

**Rastreabilidade por responsabilidade**: 025-barra-ferramentas/FR-016–FR-018 (controles/temas); 028-propriedades/FR-014 (rascunho); 029-previa-interativa/FR-011 (isolamento visual); 001-gerador-jornada/FR-048 integra. DESIGN.md específicos, plano, pesquisa e modelo de dados foram atualizados. Não foram acrescentadas entidades ou dependências de interface aos domínios existentes.

**Evidências**: `python scripts/build.py --target windows --arch amd64 --check` passou, incluindo formatação, vet, testes Go e do empacotador. Testes desktop verificam Sair visível/oculto, documento limpo, confirmação/cancelamento de alterações e rascunhos, não empilhamento de diálogo, preferência inválida/recuperada, troca sem reconstrução do formulário nem mudança de documento/seleção/zoom/histórico, isolamento dos botões da prévia e contraste mínimo 4,5:1 dos tokens principais nos três temas. O teste de preferência utiliza o armazenamento do driver de testes entre novas janelas; persistência após reiniciar o processo permanece na verificação manual. Nenhuma inspeção de janela nativa ou homologação visual foi executada nesta etapa.

**EXE Windows 1.0.0.0**: [JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe), também em dist/windows-amd64/JornadaStudio.exe. SHA-256 `0D512AE99FDB9433D11D42AFA37BD825FE0B66F72B04E0AE039AC0A59A5F8CDC`. Ambas as cópias atualizadas; somente Windows foi compilado.

**Roteiro manual**:

1. Fechar a versão anterior e abrir o novo EXE. Em tela cheia, conferir Tema, Restaurar e Sair acessíveis no topo. Restaurar/F12/Escape ocultam Sair; voltar a tela cheia reapresenta o botão.
2. Trocar entre Claro, Escuro e Azul; conferir menus, entradas, foco, diálogos, texto do cabeçalho e estados desabilitados. Conferir resolução/DPI utilizados, especialmente 1366×768 e escalas 100/150/200%.
3. Com seleção e rascunho no inspector, trocar tema: campos e conteúdo do canvas devem permanecer. Abrir as duas prévias e alternar tema no editor: fundo/primitivas e controles da jornada conservam aparência anterior.
4. Sair com documento alterado e, separadamente, apenas com rascunho: Não conserva tudo; Sim descarta e encerra. Repetir pelo fechamento do sistema em janela normal. Conferir encerramento também com prévias abertas.
5. Escolher Escuro ou Azul, fechar/reabrir e conferir a preferência. Salvar/reabrir .jornada e exportar PPT: cores do projeto não devem mudar por causa do tema do editor.

**Estado**: T105–T107 implementadas com verificações automatizadas; aceite visual, persistência entre processos e aprovação do usuário permanecem pendentes.


## Reconciliação documental por domínio — 2026-10-18

**Concluída em**: 2026-10-19 (America/Sao_Paulo). Revisão iniciada em 2026-10-18.
**Estado da documentação**: reconciliada com a implementação e as evidências disponíveis,
conforme constituição 1.2.0. Isso não afirma implementação integral nem aprovação do produto.
**Estado do produto completo**: em andamento; incremento Windows atual pronto para teste
manual, com pendências funcionais/de cobertura nos backlogs e gates externos abaixo.

**Correções documentais realizadas**:

- spec.md não descreve mais a entrega atual como exclusivamente documental.
- plan.md registra schemaVersion 2 (reader 1 com migração), toolchain fixada em
  scripts/toolchains-windows.json e árvore/caminhos de código existentes.
- Modelos diferenciam tipos Go implementados de contratos/metadados ainda planejados.
- Planos e tarefas dos cinco domínios registram o código existente e testes atuais;
  deixaram de tratar todo o trabalho como futuro. IDs foram preservados, sem renumerar.
- Backlog coordenador distingue conclusão técnica, implementação/testes parciais,
  dependência externa e validação manual. Contratos de exportação/relógio indicam
  explicitamente a diferença entre arquitetura alvo e implementação atual.

### Evidências executadas

Executor: agente, nesta estação Windows. Comando executado em apps/gerador-jornada:

```powershell
./scripts/check.ps1 -Jdk17 'C:\Desenvolvimento\jdk-17.0.11'
```

Resultado aprovado: formatação Go, go vet, suíte Go existente, três testes Python do
empacotador e nove testes da toolchain Java 6. O teste Go do gerador também compilou
fontes estruturais com ECJ/Animal Sniffer e classfile major 50. Nenhum teste de JAR
EFL real, PowerPoint, máquina limpa ou revisão visual nativa foi executado nesta passagem.
Checks não auditam automaticamente todo GoDoc/comentário nem fecham a matriz de testes
que ainda está aberta em cada domínio. A mensagem de reprovação Java emitida por um
cenário negativo foi esperada e o próprio teste terminou aprovado.

| Responsabilidade | Código e evidência existentes | Pendência preservada |
| --- | --- | --- |
| [Projeto](../002-projeto/tasks.md) | Identidade, snapshot/rollback, round-trip e manifesto v2/migração | Unicidade global/UUID de roteiros e etapas (001-gerador-jornada/T008), índice/remoção segura de recursos, falhas de gravação/recuperação e matriz completa |
| [Catálogo](../003-catalogo/tasks.md) | Sete famílias, união tipada, enum e round-trip | Metadados completos de propriedades/capacidades e matriz de validações |
| [Tela](../004-tela/tasks.md) | Bounds, seleção coletiva, FileID, duplicação/remoção reversíveis | Todos os alinhamentos/combinações de referências/camadas e revisão integral de comandos |
| [Jornada](../005-jornada/tasks.md) | Binding none/navigate, validação cruzada e reset da prévia | Warning de inalcançável, diagnóstico clicável e matriz completa A-B-C-A/remoções |
| [Demonstração](../006-demonstracao/tasks.md) | Interpretador, grafemas, ocorrências e animações/transições estruturais | Reordenação da lista de roteiros, relógio injetável e testes completos de tempo/pausa |
| Exportação | Snapshot/contexto/job no desktop e publicação pelos adaptadores | Contrato completo ExportProject/Result/Artifact/diagnóstico, testes de falha/commit e responsividade |
| Distribuição | EXE/pacote Windows 1.0.0.0 e recursos incorporados | Matriz Windows limpa/offline/DPI e inventário completo de licenças; Linux/macOS adiados |

### Revisão executável preservada

Esta correção alterou documentação; não alterou fontes do aplicativo nem recompilou
qualquer alvo. [EXE atual](../../apps/gerador-jornada/dist/JornadaStudio.exe), também
em dist/windows-amd64/JornadaStudio.exe: versão 1.0.0.0, SHA-256
`0D512AE99FDB9433D11D42AFA37BD825FE0B66F72B04E0AE039AC0A59A5F8CDC`.
Os roteiros manuais V-01 a V-09 e de cada incremento acima permanecem disponíveis;
nenhum resultado esperado foi convertido em resultado observado nesta reconciliação.

### Gates e aceite vigentes

- G-01/G-05: validar Windows real, escalas e máquina limpa/offline; ainda pendente.
- G-02: abrir/reproduzir/editar PPT nas versões 2016 e Microsoft 365; ainda pendente.
- G-03/G-04: faltam JARs/contratos/bootstrap e execução EFL real. Apenas a parte estática
  Java 6 possui evidência. Botão Java/EFL continua visível e desabilitado por orientação
  do usuário, sem reativação por esta revisão.
- Gaps de comportamento/testes: usar tasks.md da responsabilidade indicada na tabela;
  não eliminar requisito para fazer a documentação parecer integralmente atendida.
- Teste manual e aprovação integral: pendentes. Comentários positivos do usuário sobre
  incrementos não são aceite de todos os cenários/requisitos desta revisão.


**Conferência final da reconciliação (2026-10-19)**: 40 documentos Markdown da
coordenadora e dos cinco domínios conferidos; links locais/âncoras e caminhos de
código dos tasks.md válidos. IDs e ordem das 107 tarefas da coordenadora e das oito
por domínio preservados, sem duplicidades. O resolvedor Spec Kit confirmou spec,
plan e tasks da feature ativa. As caixas são evidências técnicas limitadas ao seu
escopo, não percentual de conclusão do produto. As duas cópias do EXE mantêm o hash
acima. Não foram executados novos builds nem homologações nesta revisão documental.
