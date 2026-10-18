# Quickstart: validação do Gerador de Jornada

**Status**: implementação em andamento em apps/gerador-jornada. Testes automatizados dos domínios e dos formatos estão aprovados; build desktop e homologação manual são registrados na seção de aceite. Os cenários manuais abaixo não são considerados executados por existir uma suíte automatizada.
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

Os scripts scripts/check.ps1 e scripts/build.ps1 da aplicação executam checks e build. O pacote de release terá inventário/licenças próprios. Testes que precisam de PowerPoint ou EFL devem identificar pré-requisitos e registrar ausência como não executados, sem aprovação falsa.

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
