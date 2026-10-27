# Quickstart: validação do Gerador de Jornada

**Revisão executável vigente — remoção do campo Recursos / referência (2026-10-23):** pacote reconstruído e pronto para teste manual, sem aceite; ver [seção 16](#16-remocao-do-campo-recursos--referencia) para identidades e roteiro. Go somente compilado, não executado. Os registros abaixo são históricos.

**Instrução vigente — builtin filtrado (2026-10-21), EM VERIFICAÇÃO:** a [seção 14](#14-builtin-filtrado) registra a mudança de escopo, 206 definições ativas/8 adiadas e o roteiro atual. Somente **15 Python do importador PASS** e **importador `--check` PASS** foram informados para este incremento; descoberta dos quatro novos métodos e integração Go/gofmt/vet/compile-only/build Windows/empacotamento ainda pendentes. Não há novo EXE/ZIP identificado ou hash de distribuição validado para a skin filtrada. Não executar os roteiros novos contra um binário histórico esperando a nova revisão. Esta edição altera apenas os sete documentos autorizados, sem builds/testes.

**Histórico antes da redução:** todas as indicações anteriores de “vigente”, PASS, identidades e seções 1–13 abaixo pertencem às suas próprias revisões. Em particular, **214 definições/47 diagnósticos** descrevem o builtin **antes da redução**; não substituir essas contagens históricas nem usar seus hashes como evidência da skin filtrada. A seção 13 continua registro da entrega LOG, não comprovação de build deste incremento; verificar os bytes reais após o executor principal concluir a nova distribuição.

**Entrega vigente — LOG (2026-10-21): checks técnicos PASS, pronta para coleta real.** Produção concluída e 11 arquivos Go formatados; `APPLICATION_LOG_VET_PASS` e `APPLICATION_LOG_COMPILE_ONLY_PASS`: **15 suítes somente compiladas, ZERO execuções Go**, incluindo 12 testes novos (9 de núcleo + 3 de UI). Quatro testes Python de empacotamento, tarefa Windows na raiz (exit 0), pacote sem `--check` e integridade somente leitura PASS. **Ambos os EXEs estão atualizados e idênticos**; identidades reais dos dois EXEs e do ZIP na [seção 13](#13-log-de-aplicação--2026-10-21). README da aplicação congelado já empacotado; esta finalização altera somente seis documentos, não código/testes/artefatos. **Fechar a instância antiga e abrir o novo EXE identificado para reproduzir/enviar o log.** A falha continua relatada: causa não mapeada e solução não comprovada. Nenhum aplicativo foi iniciado nesta entrega; marcadores no binário não provam execução e não há evidência de uso do novo EXE. Startup/runtime do logger, execução Go, UI nativa/DPI 100/150/200%, Office e aceite pendentes, política Windows sem contorno. T131 permanece aberta como verificação mista, assim como T130 e gates T120–T122; aprovações `pending`.

**Histórico anterior ao LOG:** todos os avisos abaixo de “vigente”, “atual” ou “mais recente”, seus PASS/tamanhos/hashes e as seções 1–12 referem-se exclusivamente às revisões anteriores. Permanecem preservados, sem validar LOG; os caminhos de distribuição agora contêm a entrega identificada na seção 13, não os bytes históricos.

**Revisão executável vigente — disponibilidade local e portabilidade de recursos (2026-10-21):** gofmt dos dois Go, vet integral, compilação das **15 suítes Go sem execução**, quatro testes Python de empacotamento, build Windows e integridade do pacote PASS, conforme registros finais do executor principal. **Compatibilidade RESOLVIDA:** após o usuário informar **“fechei a aplicação”**, o executor principal copiou o [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) para a [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) e verificou os hashes independentemente, com PASS. **Ambos agora têm 34.608.128 bytes e SHA-256 `48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d`.** ZIP não refeito: 20.408.876 bytes e SHA-256 `b40cc96c901b51d00ca042831c419c222f44f90709b7d0328baf70695d8bdf01`, checks já PASS. A [seção 12](#12-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21) registra a primeira cópia bloqueada, a checagem inicial falha pela exigência incorreta de igualdade naquele momento e a resolução posterior, com a [tabela de identidades atuais](#identidade-atual-da-distribuição-canônica-e-cópia-de-compatibilidade). Nenhum aplicativo ou teste Go foi iniciado nesta finalização, nenhum processo foi encerrado pelo agente e nenhuma política foi alterada. UI nativa/DPI/Office e aceite pendentes, T130 aberta como verificação mista e aprovações `pending`. README congelado e pacote preservados, sem necessidade de reempacotar. Captura anterior e revisão Explore somente leitura não validam a UI nova. Os avisos abaixo e todos os PASS/tamanhos/hashes anteriores, inclusive os da seção 11, ficam históricos; ambos os caminhos do EXE e o destino do ZIP agora correspondem à seção 12.

**Revisão executável vigente — redimensionamento do painel de propriedades (2026-10-21):** correção analisada, compilada e empacotada, pronta para teste manual sem aceite. A [seção 11](#11-redimensionamento-do-painel-de-propriedades--2026-10-21) identifica os artefatos atuais e checks realizados. Suítes Go somente compiladas, não executadas por política Windows, sem contorno; revisão nativa/DPI pendente. Os registros anteriores, inclusive a seção 10, são históricos; os caminhos de distribuição foram atualizados com esta revisão.

**Revisão executável vigente — Recursos / referência (2026-10-21):** correção implementada, analisada e compilada; pacote Windows pronto para teste manual, sem aceite. A [seção 10](#10-correção-de-recursos--referência--2026-10-21) identifica os artefatos atuais e os checks realizados. Testes Go somente compilados, não executados por política Windows, sem contorno. Os registros anteriores são históricos; seus caminhos de distribuição foram atualizados com esta revisão.

**Revisão executável vigente — sobreposição (2026-10-20):** pronta para teste manual, sem aceite. A [seção 9](#9-sobreposição-de-componentes--2026-10-20) identifica os artefatos atuais e os checks realmente realizados. Os hashes das seções anteriores são históricos; os mesmos caminhos de distribuição foram atualizados.

**Incremento mais recente — sobreposição (2026-10-20):** produção e quatro arquivos de regressões implementados; vet/compilação integral/build aguardam o executor principal. Testes Go não executados por política Windows, sem contorno. Identidade do novo EXE/ZIP, homologação manual e aceite permanecem pendentes na [seção 9](#9-sobreposição-de-componentes--2026-10-20). Todos os PASS, tamanhos e hashes anteriores abaixo, inclusive os de alinhamento textual da seção 8, são registros históricos exclusivos das respectivas revisões; não identificam nem aprovam um build de sobreposição. Destinos podem ser sobrescritos: conferir bytes reais antes do teste.

**Status (2026-10-19, pós-correção integral de Capture)**: implementação em andamento em [apps/gerador-jornada](../../apps/gerador-jornada/README.md), com catálogo v2/16 famílias, schema v3 e perfil visual por documento/captura/resolução arquivada implementados. Código congelado/validado pelo executor principal: gofmt/vet PASS e **15 suítes Go somente compiladas**, sem execução; bloqueio prévio por política Windows preservado, **sem PASS dos testes Go**. Python **28 PASS** (15 importador + 4 packaging + 9 Java 6, incluindo erro esperado de `String.strip`, não homologação EFL); importador `--check` PASS, mesmo inventário 214/8/270/47. **Distribuição Windows amd64 gerada e verificada.** Identidade do EXE e caminho do pacote na seção 7; checksum do ZIP no sidecar externo `.zip.sha256`. Reempacotar documentos preserva o EXE. Artefatos anteriores e seção 6 são históricos. Cenários manuais abaixo são roteiros, não execuções ou aceite.
**Input**: [spec.md](spec.md), [plan.md](plan.md), [tasks.md](tasks.md).

**Incremento de 2026-10-20:** após produção e quatro novos arquivos de testes de alinhamento textual universal, gofmt/vet PASS, 15 suítes Go somente compiladas, quatro testes Python de empacotamento PASS e build/empacotamento Windows PASS. Testes Go não executados pelo bloqueio prévio Windows, sem contorno. Evidências e identidades finais do EXE/ZIP na [seção 8](#8-alinhamento-textual-universal--2026-10-20); manuais Windows/DPI/PowerPoint e aceite pendentes. Os 28 Python, referências e hashes de 19/10 abaixo permanecem históricos, não reexecutados em 20/10.

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

## 2. Ambiente de desenvolvimento

Ferramentas locais preparadas: Go 1.27.2 amd64, Fyne 2.8.2 e MinGW GCC 16.2.0 UCRT. A partir da raiz do workspace, executar em PowerShell:

```powershell
Set-Location apps/gerador-jornada
. ./scripts/env.ps1
go test ./...
go vet ./...
go run ./cmd/gerador-jornada
```

[scripts/check.ps1](../../apps/gerador-jornada/scripts/check.ps1) verifica formatação e vet, integra testes do importador e `import_legacy_skin.py --check`, além dos testes do empacotador, **antes da suíte Go**. Os scripts de build compilam/empacotam com metadados e hashes. **Nesta estação, a política de grupo bloqueia executar os testes Go:** registrar o bloqueio, sem contornar ou converter compilação em PASS de testes. Diagnósticos do editor (`get_errors`) não substituem ferramentas Go. A distribuição Windows pós-Capture foi gerada e verificada, com identidade do EXE e caminho do pacote na seção 7; checksum do ZIP no sidecar externo `.zip.sha256`. O reempacotamento dos documentos preserva o EXE. Inventário completo de licenças permanece pendente. Não recompilar Linux/macOS na entrega corrente. Testes que precisam de PowerPoint ou EFL devem identificar pré-requisitos e registrar ausência como não executados, sem aprovação falsa.

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

Esse conjunto básico não cobre as 16 famílias nem todas as classes/skins do acervo. Ampliar a validação com os controles novos, estados, máscaras e contêineres abaixo. Galerias estáticas das 16 famílias com proveniência Java/XML estão implementadas; screenshots/goldens de execução real do legado e homologação Windows/DPI continuam pendentes.

## 4. Cenários de validação

### V-01 — edição e histórico (US1)

Inserir as 16 famílias; mover/redimensionar/alinhá-las; mudar textos, opções, skins e overrides no inspector; alterar camadas; duplicar tela e componente. Conferir novos IDs e remapeamento de referências internas. Comparar os defaults com [Definitions](../../apps/gerador-jornada/internal/catalog/model.go), incluindo tabela 300 × 200 e teclado numérico 200 × 350; componentes antigos não devem ser redimensionados pela migração.

Aplicar zoom 25/100/200/400%; conferir que bounds persistidos não mudam. Undo de arraste deve ser uma ação; novo comando após undo remove redo. Testar 100 entradas, foco de campo, Escape e atalhos.

**Expected**: geometria coerente, sem alteração oculta de dados, e operações desfeitas/refeitas integralmente.

### V-02 — arquivo, falha e recuperação (US1)

Salvar/reabrir o conjunto, comparar manifesto e recursos. Forçar interrupção após autosave, recuperar e salvar nova versão. Testar manifesto truncado, versão desconhecida, hash errado, path traversal e limites ZIP.

Abrir fixtures schema 1/2/3, conferir migração em memória para schema 3/catalogVersion 2 com bounds, IDs, vínculos e campos Style preservados. Abrir não regrava o arquivo; verificar cópia do original antes do primeiro salvamento migrado. Exercitar assets PNG/JPEG/TTF/JSON seguros e rejeições. Conferir captura de **todos os metadados da biblioteca original** em cópias privadas: definições, fontes, aliases, recursos e todas as revisões arquivadas, inclusive não usadas e todos os estados; sem enriquecer perfil privado com builtin nem fazer upgrade. Uso corrente só limita diagnósticos e acréscimo de entrada missing. Schema 3 sem perfil é aceito em transição, schemas 1/2 com o campo são recusados. Conferir hash raiz/arquivos (até 16, nível único), falha para revisão desconhecida, corrupção e bytes embedded divergentes. Recurso referenced só resolve builtin exato em hash/tamanho/MIME. Verificar que captura **não acrescenta bytes sem autorização**, nem com licença `allowed` no JSON, e informa portabilidade incompleta. Estilo normalizado/overrides com presença completa/RGBA e descritores-alvo permanecem pendentes; este roteiro não presume sua implementação.

Regressões novas do bug crítico save/reopen → `AddComponent` com `defaultconfirm`:
salvar/reabrir, addButton, inserir outra família/trocar variant após reopen, transient
states nas definições e custom namespace, sem builtin/upgrade de perfil privado e
preservando revisões não usadas. Estado runtime continua fora do documento.
Incluídas na compilação pós-correção integral de Capture pelo main (15 suítes),
**sem execução ou PASS dos testes Go**; evidências e novo EXE na seção 7.

Simular disco cheio, arquivo bloqueado e falha antes de substituir. Alterar o projeto durante save e verificar que a revisão mais nova segue marcada como alterada.

**Expected**: último arquivo válido intacto; arquivo ilegível não substitui sessão; erro localizado; referências inválidas preserváveis como rascunho.

### V-03 — jornada e prévia (US2)

Percorrer A → B → C → A. Preencher e selecionar, usar teclado numérico e reiniciar. Remover B após revisar impacto e tentar executar; desfazer a remoção.

Exercitar foco/cursor por grafemas e os oito formatos `text/number/cellular/cep/cpf/date/ddd/nis`; área multilinha, senha, combo/UF, tabela, seleção de relatório, toggle/grupos, teclados numérico/alfanumérico, rolagem e modal. Conferir que a variante `calculator` não promete nem executa cálculos e que não há consultas/negócio externo.

**Expected**: navegação correta e valores iniciais ao reentrar; prévia não altera projeto; referência quebrada localizável; undo restaura fluxo.

### V-04 — roteiro explícito (US4)

Criar, duplicar e reordenar roteiros/etapas. Reproduzir, pausar, continuar e reiniciar. Alterar um destino de botão e verificar etapas subsequentes. Testar alvo ausente, typeText vazio, Unicode e duração fora de 1..600000 ms.

**Expected**: clique apenas visual, navigate separado, sequência determinística, sem consumo de tempo durante pausa nem efeitos residuais após restart.

### V-05 — PPT exploratório (US3)

Exportar somente exploração em máquina sem Office. Abrir o arquivo nas máquinas de homologação PowerPoint 2016 e Microsoft 365, sem reparar. Usar modo apresentação para testar botões e menu.

No modo de edição, selecionar textos/formas internos de cada família, alterar um rótulo e salvar/reabrir. Comparar posição e tamanho com tolerância de 1 unidade lógica; verificar ordem, fontes internas/métricas, mídia, alpha e clipping da cena compartilhada. Instalar as faces originais no destino somente com autorização de licença; sem elas não há promessa de reprodução fiel. Não há embedding de fontes: exigir `Options.EmbedFonts` ou `Options.StrictFidelity` deve falhar com `PPT_FONT_EMBEDDING_UNSUPPORTED`, não produzir uma falsa saída portátil.

**Expected**: todos os destinos corretos, sem avanço automático, controles editáveis e nenhum screenshot substituindo a tela.

### V-06 — PPT demonstrativo (US4)

Exportar exploração e roteiro. Entrar no roteiro pelo menu; verificar movimento de cursor, realce, texto progressivo, waits e revisita a A. Conferir que apenas um prefixo de texto está visível por instante.

Medir duração das etapas, excluindo tempo de seleção do menu e os 1000 ms finais de permanência previstos no contrato. Aplicar tolerância SC-004. Testar replay e depois exploração manual, sem avanço temporal herdado.

**Expected**: execução consistente nas duas versões do PowerPoint, sem macros, dados interativos falsos ou animações omitidas.

### V-07 — fontes Java/EFL (US5)

Somente com perfil homologado: gerar para diretório novo, inspecionar as capacidades declaradas, Javadoc, listeners e separação por tela. A evidência estrutural histórica cobre sete famílias, não as 16 atuais. Gerar novamente sem alterações e comparar conteúdo. A ação pública Java/EFL permanece visível e desabilitada; este roteiro não autoriza remover gates.

Usar toolchain e comandos registrados pelo perfil para compilar contra API SE 6 e EFL real; verificar classfile major 50. Executar no hospedeiro e percorrer A → B → C → A na EDT, com campos, seleção e teclado.

Testar colisões pós-normalização, palavra reservada, pacote inválido, caracteres especiais e estilo sem capacidade EFL.

**Expected**: fontes válidos e integração funcional. Sem JARs/perfil, registrar não executado e verificar diagnóstico JAVA_PROFILE_UNVERIFIED; não contar dublês como homologação.

### V-08 — distribuição e cancelamento (US6)

Após conclusão e identificação do novo build, copiar distribuição para Windows 10/11 x64 limpos, com usuário comum e rede desligada; sem Go/Node/Java/Office/WebView adicional instalado como requisito. Abrir EXE em janela própria; salvar/reabrir, gerar PPT e conferir Java/EFL visível e desabilitado.

Testar DPI 100/150/200%, paths com espaços/acentos, diretório sem permissão e exportação cancelada. Durante exportação, editar o projeto e verificar revisão do snapshot. Tentar segundo job.

**Expected**: dependências efetivas registradas, edição utilizável, um job por vez, cancelamento em até 1 s no conjunto nominal e nenhuma publicação parcial.

### V-09 — integração, documentação e desempenho

Executar fluxo completo do conjunto de referência e comparar canvas, prévia e PPT. Comparar saída Java somente no ambiente de perfil homologado de V-07, sem reativar a ação pública. Medir conjunto nominal de 10 × 50 componentes e 50 etapas conforme spec.md.

Revisar dependências entre pacotes e documentação de elementos públicos/privados, campos, variáveis de pacote e locais; verificar Javadoc gerado. Registrar quais itens tiveram check automático e quais revisão humana.

**Expected**: SC-001 a SC-007 comprovados, com falhas e limites explicitados. Não exigir identidade de rasterização entre Fyne/PPT/Swing.

## 5. Registro de evidência e conclusão

Para cada cenário registrar data, versão do editor, revisão do projeto, versão do catálogo/perfil, toolchain, SO, DPI, PowerPoint/JARs quando aplicável, resultado e localização do artefato de evidência. Textos de projeto não devem ser copiados integralmente para logs.

Resultados possíveis: aprovado, reprovado ou não executado com causa. Campo vazio, requisito documentado ou teste com dublê não é aprovação. Atualizar tarefas apenas quando existir a evidência correspondente.

## 6. Histórico de revisões e builds

As aprovações, capturas e hashes abaixo pertencem **somente às revisões descritas**. Expressões como “atual” ou “mais recente” são históricas; não identificam o novo build nem aprovam a suíte deste incremento. A evidência vigente está na seção 7.

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

**Evidência visual histórica pelo agente**: janela nativa inicializada pelo comando `--smoke-screenshot`, captura então inspecionada com fundo original e sete famílias acessíveis. A captura histórica não está mais disponível no workspace; este registro não é nova evidência visual nem execução dos cenários manuais pelo usuário.

**Exemplo portátil**: jornada-com-fundo.jornada (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/jornada-com-fundo.jornada`). A imagem acompanha o arquivo; abrir projetos anteriores conserva o conteúdo deles.

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

**Evidência anterior desta revisão de interface**: a captura então registrada mostrava rodapé fixo e Java/EFL desabilitado, antes da inclusão do novo ícone; ela não está mais disponível no workspace. O PPT histórico, também ausente do dist atual, foi registrado como gerado antes da troca do ícone e contendo nove slides com fade e oito movimentos de cursor, conferidos no ZIP/XML. Não houve homologação visual no PowerPoint; esses registros não validam a continuação atual.

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
| Windows amd64 | ZIP (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), [EXE atual](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | `051E243F55548E18C1A539D7307B197D233F22CA86049FF51D5536E22C07A88E` | Build nativo, checks Go/vet e PE x64 conferidos |
| Linux amd64 | TAR.GZ (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-linux-amd64.tar.gz`) | `6785223DEF2419721890E2C09E0EB279E3E6EEDF0C67B5E1B3B935286B9EDB2B` | Go/vet e suite passaram no contêiner Debian 12, ELF64 x86-64 e dependências conferidos |
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


<a id="reconciliação-documental-por-domínio--2026-10-18"></a>
### Reconciliação documental por domínio — 2026-10-18 (histórico)

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

#### Evidências executadas nessa reconciliação

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

#### Revisão executável preservada nessa reconciliação

Esta correção alterou documentação; não alterou fontes do aplicativo nem recompilou
qualquer alvo. [EXE atual](../../apps/gerador-jornada/dist/JornadaStudio.exe), também
em dist/windows-amd64/JornadaStudio.exe: versão 1.0.0.0, SHA-256
`0D512AE99FDB9433D11D42AFA37BD825FE0B66F72B04E0AE039AC0A59A5F8CDC`.
Os roteiros manuais V-01 a V-09 e de cada incremento acima permanecem disponíveis;
nenhum resultado esperado foi convertido em resultado observado nesta reconciliação.

#### Gates e aceite registrados nessa reconciliação

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

## 7. Incremento do acervo Java/EFL — 2026-10-19 (evidência vigente)

**Identidade final da distribuição pós-Capture:**

| Artefato | Tamanho | SHA-256 |
| --- | --- | --- |
| [EXE Windows amd64](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | 34.589.696 bytes | `2ecee6c700b06c245d704f497ea00d143b801ea17812e8870246fde477c001ed` |
| ZIP Windows amd64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 20.401.442 bytes | `b7ce06bbe1b2729f798bb553095ab67ccaea3e6898a58b2026e18ecfed94b1f2` |

Ambas as cópias do EXE são idênticas. Após o último reempacotamento, o executor
principal conferiu CRC de todas as entradas, cada SHA256SUMS, sidecar do ZIP,
cabeçalho PE amd64, README definitivo e manifesto/inventário com aprovações
`pending`. Conferiu também links/codificação de 42 documentos e as referências
estáticas de quatro estados XML e 14 teclas Java da galeria. Isso não executa
as suítes Go, a interface nativa ou o PowerPoint e não constitui aceite humano.

**Estado da continuação:** implementação parcial, sem paridade integral homologada.
Resultados do executor principal após a correção integral de Capture, com código
congelado: gofmt/vet PASS; **15 suítes Go compiladas,
não executadas** por política Windows; **28 testes Python e importador check PASS**.
**Distribuição Windows amd64 gerada e verificada**. Hashes antigos
permanecem históricos identificados abaixo; capturas, execução nativa,
PowerPoint e aprovação humana continuam pendentes.

### Implementação e inventário

- **catalogVersion 2 / 16 famílias tipadas:** preservadas `button`, `label`, `title`, `textField`, `radio`, `report`, `numericKeyboard`; novas `toggleButton`, `textArea`, `passwordField`, `comboBox`, `table`, `alphaKeyboard`, `progressBar`, `panel`, `dialog`. Defaults e tabela completa no [README](../../apps/gerador-jornada/README.md#componentes-visuais), conferidos com o catálogo real. Skins/overrides estão disponíveis no inspector; 16 famílias não comprovam paridade de todo o acervo.
- **Perfil:** `legacy-skin-v1-0a668345039c24654d61a3c6bf9b78366245552604a0caffaf4074546f6f4f12`, com **214 definições, 8 fontes originais e 270 recursos disponíveis**, após importar os 18 pacotes adicionais de `outros_materiais`.
- **Fontes:** 8/8 faces aceitas; não faltam as Futura/Vera fornecidas. [futura_mh.ttf](../../material-legado/elemetos_componente/outras_fontes/futura_mh.ttf) não está corrompida: o parser ignora somente a sentinela histórica U+FFFF (não-caractere) do `cmap` formato 4, com aviso e bytes preservados. Aceitação técnica não comprova licença.
- **Cena/runtime:** canvas e prévia usam SFNT/`Fyne.FontSource`, sem Arial silencioso. Há foco/cursor por grafemas, máscaras dos oito formatos, seleção de tabela/relatório/toggle, teclados, rolagem e modal, sem negócio externo. A variante calculadora **não calcula**. Delegates ausentes e layouts incompletos permanecem diagnosticados.
- **Projeto:** schema 3, reader 1/2/3, migração em memória preservando bounds/IDs/vínculos/Style, sem regravar ao abrir e com backup original antes do primeiro save migrado. `Document.AppearanceProfile *skin.Profile` e cópia profunda implementados; perfil opcional na leitura v3 em transição e proibido em v1/v2. Captura automática de metadados e resolução arquivada existem, **não transporte automático de binários**.
- **Perfil por documento:** schema próprio 1, `id/version/sha256/definitions/fonts/aliases/resources/diagnostics/archived`, até 16 arquivos em nível único. Hash JSON determinístico exclui somente o hash próprio e inclui arquivos/hashes. `Library` imutável e `ForVersion` exato; desconhecidos não usam fallback. `Capture` conserva todos os metadados da biblioteca original em cópias privadas: definições, fontes, aliases, recursos e todas as revisões arquivadas, inclusive não usadas e todos os estados. Não enriquece perfil privado com builtin nem faz upgrade; uso só limita diagnósticos e acréscimo de entrada missing. Conserva bytes já incorporados, nunca acrescenta novos sem autorização, mesmo com JSON `allowed`. Builtin permanece `referenced/unknown`. `ResourceData` valida AssetID/hash/tamanho/MIME para embedded e bytes builtin exatos para referenced.
- **Desvio arquitetural justificado:** snapshot integral de metadados corrige o bug crítico save/reopen → `AddComponent`/`defaultconfirm` do catálogo podado e preserva inserir famílias/trocar variantes após reopen. Não é incorporação integral binária; alvo normativo SUBSET de bytes usados e permitidos permanece. Payload selado mais amplo altera o hash do perfil, distinto da versão do inventário e do hash do EXE. Regressões save/reopen/addButton/variant/transient states/custom namespace incluídas na compilação das 15 suítes pelo main após a correção; sem execução ou PASS de testes Go.
- **Lacunas de modelo mantidas:** wire ainda usa `Style` (`SkinName/SkinVersion`, RGB, `FillAlpha` opcional) e DTOs `Fonts/Aliases/Resources`, não `component.appearance` normalizado, presença tipada completa/RGBA ou descritores-alvo `fontFaces/resourceIds/requestedFamily/unresolved`. Família desconhecida conserva string literal em Style e diagnóstico. Portabilidade licenciada completa continua pendente.
- **Desktop/PPT e galeria:** resolução por documento e cache SHA/face exata integrados; galerias estáticas das 16 famílias baseadas em Java/XML, não screenshots Java nem homologação Windows/DPI.
- **PPT:** cena compartilhada editável, nomes internos de fontes/métricas, mídia, alpha e clipping. **Nenhum embedding de fontes**; `Options.EmbedFonts`/`Options.StrictFidelity` falham com `PPT_FONT_EMBEDDING_UNSUPPORTED`; `FontEmbedding=false`/`CompleteFidelity=false`. Sem fontes originais instaladas no destino, não prometer reprodução fiel.
- **Java/EFL:** ação visível **desabilitada**, gate `JAVA_PROFILE_UNVERIFIED` mantido; expansão do catálogo não homologa a saída Java.

Os **47 diagnósticos** do inventário são:

| Código | Quantidade |
| --- | --- |
| `RESOURCE_MISSING` | 16 referências a **16 caminhos ausentes** |
| `JAVA_FONT_STYLE_NOT_APPLIED` | 14 |
| `FONT_LICENSE_UNKNOWN` | 8 |
| `FONT_EMBEDDING_RESTRICTED` | 7 |
| `FONT_ALIAS_AMBIGUOUS` | 1 |
| `FONT_CMAP_SENTINEL_NONCANONICAL` | 1 |

### Evidências desta revisão

**Marco pós-correção integral de Capture:** resultados fornecidos pelo main após
a mudança e suas regressões save/reopen/addButton/variant/transient states/custom
namespace. A compilação inclui essas regressões, mas não comprova seu resultado
em execução.

Ambiente registrado: Windows 11 10.0.22631 do workspace, Go 1.27.2, Fyne 2.8.2 e MinGW GCC 16.2.0. Executor: main. Resultados após a correção:

| Verificação | Resultado observado |
| --- | --- |
| `python -X utf8 -m unittest discover -s scripts -p 'test_*.py'` | **28 PASS**: 15 importador + 4 packaging + 9 Java 6; erro esperado de `String.strip` incluído, **não homologação EFL** |
| `python -X utf8 scripts/import_legacy_skin.py --check` | **PASS**, sem alterar outputs; mesmo inventário `legacy-skin-v1-0a668345039c24654d61a3c6bf9b78366245552604a0caffaf4074546f6f4f12`: 214 definições, 8 fontes, 270 recursos, 47 diagnósticos |
| `gofmt` | **PASS**, após a correção integral de Capture |
| `go vet ./...` | **PASS** |
| `go test -c -o ./build/compile-tests/ ./...` | **PASS de compilação de 15 suítes, SEM executá-las** |
| Execução dos testes Go | **Não realizada**; bloqueio prévio por política de grupo do Windows preservado, sem contorno ou PASS atual da suíte |
| Tarefa VS Code `process: Jornada Studio: compilar Windows` | **PASS**, novo EXE Windows amd64 construído; identificação abaixo |
| `python scripts/build.py --target windows --arch amd64` | **PASS**, distribuição Windows gerada |
| Integridade do pacote | **PASS**: CRC do ZIP, SHA256SUMS de cada arquivo, PE amd64, README com normalização de quebras de linha, manifesto/inventário/build-info e sidecar `.zip.sha256` |
| Cópias do EXE | **PASS**, bytes idênticos em `dist/windows-amd64` e na cópia de compatibilidade `dist` |
| Referências estáticas da galeria | **PASS**: 4 estados confrontados com XML e 14 geometrias de teclas com Java original; **não homologação visual** |
| Capturas, execução manual, PowerPoint e aceite | **Não realizados nesta revisão** |

[scripts/check.ps1](../../apps/gerador-jornada/scripts/check.ps1) agora inclui os testes do importador e `--check` antes da suíte Go. Isso não afirma execução integral aprovada do script. `get_errors` não substitui vet, compilação ou execução Go. Aprovações e hashes da seção 6 não são evidência deste código.

### Limites e próximos registros

Permanecem pendentes: 16 recursos/caminhos ausentes, delegates/layouts/referências
incompletos, 8 licenças desconhecidas e 7 fontes com embedding restrito; portabilidade
binária licenciada, aparência/overrides/RGBA e descritores normalizados do contrato;
execução Go conforme política e gates dependentes de T118/T123.
Galerias estáticas já existem, mas screenshots/goldens reais do
legado, embedding PPT, Windows/DPI 100/150/200%, máquina limpa/offline e PowerPoint
2016/Microsoft 365 continuam abertos, assim como **aceite expresso**. Renderer novo
não é referência de execução Java. Calculadora continua sem cálculos, Java desabilitado.

<a id="novo-exe-pós-capture--pacote-e-homologação-pendentes"></a>

### Distribuição Windows pós-Capture — gerada e verificada; homologação pendente

**EXE Windows amd64 1.0.0.0:** [JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe),
34.589.696 bytes, SHA-256
`2ECEE6C700B06C245D704F497EA00D143B801EA17812E8870246FDE477C001ED`.
Construção pela tarefa VS Code acima aprovada pelo main, não execução do aplicativo.
`manualApproval=pending`; nenhuma homologação Windows/DPI/PowerPoint ou EFL atribuída.
Não há novo build Linux/macOS.

**Pacote Windows:** apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`).
Checksum externo: apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256`).
Empacotamento e verificações de integridade **PASS**, incluindo manifesto e seu hash
em build-info, inventário e README. `manualApproval=pending` e
`legacySkin.redistributionApproval=pending` conferidos; licença não é aprovada por hash.
O reempacotamento para incorporar os documentos atualizados preserva os bytes do EXE
e atualiza SHA256SUMS e o sidecar externo. Hash e tamanho finais do ZIP serão
registrados externamente após essa etapa, sem incluí-los no README empacotado.
**Entrega de build para revisão realizada.** T125 permanece aberto pelo escopo
dependente de T118/T123 e demais gates, não por falta de pacote;
T118/T120–T124 não são concluídos por compilação/build.

**Roteiro manual pendente para este EXE:** fechar a versão anterior e abrir o novo
executável; salvar/reabrir `.jornada`, inserir Button/`defaultconfirm` e outras famílias,
trocar variantes e conferir estados/namespace privado/revisões arquivadas sem upgrade
ou cópia automática de bytes. Executar também V-01 a V-09, galeria em DPI 100/150/200%
e PowerPoint 2016/Microsoft 365, registrando ambiente, hash, observado e divergências.
Nenhum desses passos foi executado ou aprovado nesta atualização.

### Artefatos históricos da etapa anterior — NÃO representam esta continuação

Hashes abaixo identificam exclusivamente os artefatos anteriores, não o código com
perfil por documento desta continuação. EXE e distribuição atuais estão identificados
acima; o checksum do ZIP fica no sidecar externo. Caminhos de destino podem ser
sobrescritos: os hashes abaixo são registros históricos, não hashes dos bytes atuais.

**EXE anterior:** [JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), também copiado para [dist/JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe); ambos com SHA-256 `76b7d4c89004c882db780ac4876c794d5141d76c1f68d22b5e40d1d7619d6166`, 34.442.752 bytes. Produto permanece 1.0.0.0, sem assinatura, aprovação manual pendente.

**Pacote anterior:** JornadaStudio-1.0.0.0-windows-amd64.zip (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), SHA-256 `63c7140424af7fe7d99c953621116038cef05a8421e1069461ecb9ebd5e808e8`, 20.343.057 bytes. Na etapa anterior foram conferidos CRC de todas as entradas, hashes de SHA256SUMS, manifesto e README/inventário. Uma checagem ad hoc inicial de texto acentuado falhou por conversão no pipe do PowerShell 5.1; repetida com literal ASCII sobre os mesmos arquivos, aprovada sem modificar o pacote. Integridade desse ZIP não comprova execução gráfica nem valida esta continuação.

**Importação adicional:** 18 pacotes de [outros_materiais](../../material-legado/elemetos_componente/outros_materiais) têm raiz de classpath `src/images` comprovada nas configurações Eclipse/Ant. As 271 cópias de recursos em `bin` foram comparadas aos originais: **nenhuma diferença**, sem duplicar sua importação. Metadados de IDE/build e bancos de thumbnails não são conteúdo executável nem recurso do editor. Duplicatas byte a byte preservam todas as origens no manifesto. O usuário escolheu a versão **clara original** de `/azul/label_lot_azul.png`; a resolução está fixada em legacy_resource_choices.json (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/scripts/legacy_resource_choices.json`). Pin de origem/hash divergente produz `RESOURCE_CHOICE_STALE`, sem substituir por homônimo. As duas candidatas continuam auditáveis no inventário; escolha de imagem não é homologação nem permissão jurídica.

O pacote histórico levou o manifesto de recursos e seu SHA-256 `04b436ef501a2f8649ff231f06325edfc089d1ff515d974ff329f79b09b12797` em build-info, com 47 diagnósticos e `redistributionApproval=pending`. Não é evidência de novo ZIP. Os roteiros V-01 a V-09 continuam **pendentes de execução e aceite**; não há novo screenshot ou aprovação manual. Políticas de segurança e gates Java/EFL permanecem intactos.

## 8. Alinhamento textual universal — 2026-10-20

**Registro histórico:** identidades desta seção pertencem à revisão de alinhamento. Os caminhos de distribuição agora contêm a revisão de sobreposição identificada na seção 9, que preserva a funcionalidade textual.

**Escopo:** FR-055, 028/FR-017, 004/FR-017; integração em 001/T126 e invariantes em 004/T013. O campo **Alinhamento do texto** fica após conteúdo e antes da aparência avançada nas 16 famílias, inclusive vazias, com **Padrão da skin / Esquerda / Centro / Direita**. A dica distingue texto de movimento/resize; **Posição do elemento** continua separada. `Style.TextAlign`/resolvedor compartilhado são reutilizados, inclusive no popup do combo e **Sem seleção**, sem novo schema, catálogo ou domínio.

### Estado de verificação deste incremento

| Verificação | Resultado de 20/10 informado pelo executor principal, após produção e quatro novos arquivos de testes |
| --- | --- |
| Regressões desktop/presentation/projectzip/PPT | Quatro novos arquivos para 48 combinações (16 famílias × esquerda/centro/direita), incluídos na compilação; testes não executados. Fixtures de fonte usam `goregular` permitida, sem liberar fontes/licenças do legado |
| Formatação | `gofmt` PASS |
| Análise estática | `go vet ./...` PASS |
| Compilação das suítes | `go test -c -o ./build/compile-tests/ ./...` PASS: **15 suítes somente compiladas, sem execução** |
| Execução de testes Go | Não realizada pelo bloqueio prévio de política Windows, mantido sem contorno e sem PASS de testes |
| Testes Python de empacotamento | `python -X utf8 -m unittest discover -s scripts -p test_build.py`: **4 PASS**; não é reexecução dos 28 Python de 19/10 |
| Compilação Windows | Tarefa VS Code `process: Jornada Studio: compilar Windows` PASS; EXE Windows amd64, versão de produto **1.0.0.0**, **34.590.208 bytes** |
| Identidade do EXE | [JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), SHA-256 **A26B96049E4F1D3758600787959410557C97B11D576D280983B2550A37B8506A** |
| Novo ZIP de distribuição | Pacote Windows amd64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), **20.401.996 bytes**, SHA-256 `4686e89d05738b2b9a2cd944192cf374e3164328ca42d6cf828af0da26e5d146` |
| Integridade da distribuição | PASS: CRC, SHA256SUMS de todas as entradas, sidecar do ZIP, PE amd64, README definitivo e inventário. Ambas as cópias do EXE são idênticas; aprovações manual e de redistribuição permanecem `pending` |
| Interface nativa, DPI, PowerPoint e aceite | Pendentes; os passos abaixo são roteiro, não resultados observados |

Os 28 Python, verificações de referências/importador e hashes de 19/10 nas seções anteriores são **históricos**; não foram declarados reexecutados em 20/10. T126 e 004/T013 permanecem abertas com verificação parcial: compilação não comprova execução das regressões ou homologação.

O empacotador passou a excluir backups terminados em `~`, inclusive o
`JornadaStudio.exe~` que o linker pode conservar ao atualizar uma revisão em uso.
Não apaga esses arquivos locais nem encerra processos: apenas os exclui do
inventário e da distribuição. A regressão em
[test_build.py](../../apps/gerador-jornada/scripts/test_build.py) confere a
preservação local e ausência do backup no ZIP/checksums; quatro testes Python
passaram novamente após essa correção.

Gates inalterados: Java/EFL visível/desabilitado e `JAVA_PROFILE_UNVERIFIED`; fontes/licenças/recursos pendentes; PPT sem embedding, `FontEmbedding=false`/`CompleteFidelity=false` e `PPT_FONT_EMBEDDING_UNSUPPORTED` em EmbedFonts/StrictFidelity. A fonte permitida de teste não remove esses limites.

### Roteiro manual aberto — executar somente na revisão identificada

1. Usar o EXE identificado acima e registrar o ambiente da execução manual; se houver outro build, registrar sua identidade real antes de testar. Preparar as 16 famílias com texto curto/acentuado, longo e vazio; manter exemplos com pais/filhos, fundo, células de tabela, áreas de relatório e legendas de teclas. Registrar bounds, hierarquia e regiões clicáveis antes da edição.
2. Selecionar cada família e conferir campo sempre visível, ordem após conteúdo/antes de aparência avançada, quatro opções e dica. Conferir **Posição do elemento** separada. Selecionar novamente deve carregar a escolha confirmada, inclusive em elemento vazio.
3. Para cada uma das 48 combinações (16 × três overrides), escolher Esquerda/Centro/Direita e **Salvar alterações**. Comparar canvas e prévias: somente texto se alinha dentro de cada região; bounds, pais/filhos, fundos, células/áreas/teclas e hit areas não mudam. Abrir popup do combo com e sem seleção e conferir opções e **Sem seleção**.
4. Alterar um rascunho e **Descartar**; conferir último valor confirmado. Salvar outra escolha, **Desfazer/Refazer**, salvar/reabrir o projeto e conferir valor/conteúdo/geometria. Escolher **Padrão da skin** remove somente o override textual; validar herança sem alterar a migração histórica de vazio antigo para esquerda.
5. Repetir a revisão visual em Windows com DPI 100/150/200% e zoom pertinente. Exercitar clique/edição/seleção nos mesmos pontos para confirmar preservação de hit areas e isolamento dos valores temporários da prévia.
6. Exportar PPT e abrir no PowerPoint 2016/Microsoft 365. Conferir alinhamento por região, texto ainda editável, geometria/fundos/hierarquia visual e ausência de reparo. Registrar fontes instaladas/licenças e divergências; XML/compilação não substituem esse passo. Não exigir nem declarar embedding implementado.

### Registro técnico informado e roteiro manual ainda não executado

Resultados técnicos informados pelo executor principal em 20/10; preencher os campos manuais somente após observação, sem checksum fictício nem aprovação por expectativa:

| Campo | Evidência técnica / conteúdo manual ainda a registrar |
| --- | --- |
| Data, executor e revisão do código | 2026-10-20, executor principal, após produção e quatro novos arquivos de testes; identidade binária registrada acima, sem revisão de código adicional informada |
| Versão, caminho, tamanho e SHA-256 do EXE | Produto 1.0.0.0, Windows amd64, 34.590.208 bytes e SHA-256 acima; identidade técnica, não evidência de teste manual do executável |
| Novo pacote ZIP | Gerado e verificado pelo executor principal; tamanho e hash final na tabela acima, sem execução do aplicativo |
| Comandos e resultados de gofmt/vet/compilação/build/Python | PASS conforme tabela técnica acima; 15 suítes Go somente compiladas, quatro Python de empacotamento executados |
| Fixture/projeto, perfil/skin e fonte | Aguardando identificação da fixture efetivamente usada |
| Windows, resolução, DPI, zoom e versão do PowerPoint | Aguardando ambiente efetivamente testado |
| Família, opção, conteúdo vazio/preenchido e popup | Aguardando resultado de cada caso e cobertura das 48 combinações |
| Resultado manual observado e evidências | Não executado; registrar aprovação/reprovação/bloqueio e localização da evidência após testar. Execução Go também não realizada pelo bloqueio prévio Windows, sem contorno |
| Divergências, reteste e aceite explícito | Pendentes; T120–T122 permanecem abertos |

## 9. Sobreposição de componentes — 2026-10-20

**Escopo:** 001/FR-056/T127, 028/FR-018 e 004/FR-018/T014. Implementados [screen.PlaceComponent](../../apps/gerador-jornada/internal/screen/layer.go) e o formulário inspector_layer.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_layer.go`), integrado a `componentInspector`. A operação reordena somente a origem na lista plana `Components`, fundo → frente, preservando a ordem relativa dos demais. Referência vazia indica extremo; referência válida indica adjacência antes/depois. Não há novos campos/schema/catálogo ou relação permanente frente/atrás; ParentID, bounds, estilo, conteúdo e camadas dos filhos não são alterados pela operação. Filhos não acompanham automaticamente a camada do pai.

### Estado técnico — verificado pelo executor principal

| Verificação | Estado deste incremento |
| --- | --- |
| Produção e regressões | Implementadas; quatro arquivos: desktop/layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/layer_test.go`), screen/layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/screen/layer_test.go`), projectzip/layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/projectzip/layer_test.go`), presentation/layer_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/presentation/layer_test.go`) |
| Formatação | `gofmt` PASS nos sete arquivos Go deste incremento; os quatro testes de alinhamento existentes não foram editados |
| Análise estática | `go vet ./...` PASS |
| Compilação integral das suítes | `go test -c -o ./build/compile-tests/ ./...` PASS; **15 suítes somente compiladas, não executadas** |
| Execução Go | Não realizada por política Windows; sem contorno e sem PASS dos testes. Compilar não é executar |
| Testes Python de empacotamento | `python -X utf8 -m unittest discover -s scripts -p test_build.py`: **4 PASS**, reexecutados neste incremento |
| Build Windows e empacotamento | PASS: tarefa VS Code `process: Jornada Studio: compilar Windows`, executada na raiz do workspace, e `python -X utf8 scripts/build.py --target windows --arch amd64` no módulo |
| Integridade | PASS: CRC, cada SHA256SUMS, sidecar, PE amd64, README definitivo, inventário/build-info e cópia de compatibilidade idêntica. Nenhum backup terminado em `~` distribuído |
| Interface nativa, DPI, PowerPoint e aceite | Pendentes; roteiro abaixo não executado nem aprovado |

Os testes existentes de alinhamento textual e seus quatro arquivos não são substituídos por estes testes de camada. Vet/compilação/build e os quatro Python de empacotamento foram reexecutados conforme tabela; os outros resultados históricos, os 28 Python gerais e as identidades anteriores não são promovidos a evidência deste incremento. Gates Java/EFL, fontes/licenças/recursos e PPT sem embedding permanecem inalterados; `manualApproval` e `redistributionApproval` continuam pendentes.

### Identidade observada e verificações manuais ainda pendentes

| Campo | Evidência técnica / pendência manual |
| --- | --- |
| Data, executor, revisão do código e comandos | 2026-10-20, executor principal, incremento de sobreposição; comandos e escopo na tabela acima. Sem revisão Git disponível; identidade binária abaixo |
| EXE Windows amd64 | [Executável](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), produto **1.0.0.0**, **34.600.448 bytes**, SHA-256 `48d84e77af6a64665f5cbbfec91926bee740f789d128bbae87cc1ff33c739710`; PE amd64 confirmado |
| Cópia de compatibilidade do EXE | [Cópia](../../apps/gerador-jornada/dist/JornadaStudio.exe) com bytes e SHA-256 idênticos ao executável acima |
| ZIP e sidecar | Pacote (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), **20.406.025 bytes**, SHA-256 `67571d8cbe9a45ca1a91741591d1783fcf26c05c3cbc1f7f379143f91ed2d009`; sidecar (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256`) correspondente confirmado |
| Integridade e conteúdo do pacote | PASS conforme tabela técnica, sem executar o aplicativo; aprovações manual/redistribuição `pending`, inventário legado preservado |
| Fixture, perfil/skin, ambiente manual | A preencher: projeto/revisão, Windows, resolução, DPI, zoom e PowerPoint efetivamente usados |
| Resultado observado, divergências e reteste | Não executado; preencher por caso somente após teste, com evidência localizável |
| Aceite explícito da revisão identificada | Pendente; T127/004-T014 e T120–T122 permanecem abertos |

### Roteiro manual aberto — após identificar o novo build

1. Registrar a identidade real do EXE e fechar a instância antiga antes de abrir a revisão de sobreposição. Criar fixture com objetos parcialmente sobrepostos, referências em ambos os lados da origem e exemplos das 16 famílias; registrar ordem canônica, bounds, estilo, conteúdo, vínculos e ParentID antes de editar. Incluir pai com filhos e tela com um único componente. Para nomes repetidos, usar fixture de rascunho apropriada, sem dispensar a validação de nomes do domínio.
2. Selecionar individualmente cada família. Conferir **Sobreposição**, a dica **Camada N de total (1 = fundo)** e os cinco modos: **Manter posição atual**, **À frente do objeto**, **Atrás do objeto**, **Trazer para a frente**, **Enviar para o fundo**. Iniciar no modo neutro e **Objeto de referência** desabilitado, com **Selecione um objeto**. Salvar somente outra propriedade deve manter a ordem.
3. Escolher cada modo relativo: habilitar a referência, listar somente outros objetos da mesma tela, excluir o próprio e identificar por **ordinal · nome**, inclusive em rascunho com nomes repetidos. A escolha não deve mudar o canvas/documento confirmado antes de **Salvar alterações**. Salvar deve colocar imediatamente antes/depois da referência, preservando a ordem dos demais. Repetir com origem de ambos os lados e com adjacência já satisfeita.
4. Escolher os modos de extremo e salvar: Trazer leva ao fim/frente absoluta, Enviar ao início/fundo absoluto. Conferir que **Frente ↑ / Fundo ↓** existentes continuam movendo somente um passo. Selecionar pai e mudar camada: somente sua entrada muda, sem acompanhar filhos ou alterar ParentID/geometria. Na tela unitária, extremos são no-op; modo relativo sem referência válida não salva.
5. Editar camada e conteúdo/aparência no mesmo rascunho; **Descartar** mantém todos os valores confirmados. Salvar edição válida confirma tudo numa única ação; **Desfazer/Refazer** restaura o conjunto e a ordem. Tentar salvar modo relativo sem referência e com outro campo alterado; nenhum campo deve ser publicado. Exercitar referência removida/desatualizada em fixture controlada e diagnóstico sem estado parcial; não substituir silenciosamente a referência. Conferir proteção do rascunho ao solicitar saída e cancelar.
6. Salvar/reabrir o projeto: conservar apenas a nova ordem canônica e os dados confirmados, sem campo extra, versão nova ou relação permanente com a referência. Ao remontar o formulário, modo neutro e dica atualizada; edições posteriores são independentes da referência usada antes. Comparar bounds, estilo, conteúdo, IDs, vínculos, ParentID e posição relativa dos demais/filhos. Conferir canvas, seleção do objeto frontal no ponto sobreposto e prévias, respeitando a composição de contêineres existente.
7. Repetir revisão de controles, foco, dica e ordem visível em Windows DPI 100/150/200% e zoom pertinente. Exportar PPT e abrir no PowerPoint 2016/Microsoft 365: conferir ordem visual, objetos/textos editáveis e ausência de reparo, sem promover compilação/cena a homologação Office. Registrar limitações de composição/fontes e divergências reais; não declarar movimento de subárvore nem embedding implementados.

**Conclusão deste registro:** produção, regressões, análise estática, compilação e pacote Windows entregues para teste manual. Sem execução das suítes Go, homologação nativa/Office/DPI ou aprovação do produto. Resultados futuros devem identificar a revisão efetivamente verificada e não apagar os registros históricos.

## 10. Correção de Recursos / referência — 2026-10-21

**Escopo:** clarificação de 001/FR-026/FR-050/FR-052 e [028/FR-016](../028-propriedades/spec.md#aparência-e-famílias-pós-v1--2026-10-19); verificação integrada [T128](tasks.md#phase-22-correção-de-recursos--referência--2026-10-21), aberta. Produção e dois novos arquivos de regressões implementados, analisados e compilados; build e pacote Windows gerados e verificados conforme tabela abaixo. Os testes Go não foram executados por política Windows, sem contorno; validação nativa/DPI e aceite permanecem pendentes. As evidências e hashes das seções 6–9 são históricos, não validam esta correção.

### Causa conhecida e comportamento corrigido

A consulta desktop `componentDiagnostics` acrescentava diagnósticos de `visualprofile.Diagnostics(document)` do documento inteiro a cada componente. Assim, falhas alheias na mesma/outra tela apareciam na seleção saudável e a prévia completa, que já chama a consulta por componente, repetia avisos globais a cada iteração. O filtro por caminho compartilhado também podia atribuir avisos de outra definição à seleção; a coleta só das primitivas atuais não abrangia todos os recursos de partes/estados futuros. Avisos genéricos podiam repetir condições locais ou continuar pedindo métricas após medição SFNT; comparar mensagens não identificava corretamente duplicatas. A lista sem rolagem própria aumentava o formulário com todos os detalhes.

A nova API visualprofile.ComponentDiagnostics (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/visualprofile/component_diagnostics.go`) substitui somente `Screens` na consulta temporária pelo componente/rascunho recebido, com revisão exata. Limpa ParentID apenas na cópia local para não fabricar falha de ancestral excluído pelo isolamento; não altera vínculo, documento, perfil, recursos ou histórico. Perfil e recursos completos continuam na validação global, inclusive corrupção não usada; sem fallback ou correção silenciosa.

O inspector_appearance.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_appearance.go`) reúne todas as `group.Definitions`, partes compostas e recursos de todos os estados; diagnóstico com tipo/nome exige definição exata, não apenas fonte/imagem compartilhada. Prefere diagnóstico local detalhado, omite duplicatas genéricas e `TEXT_METRICS_REQUIRED` já resolvido por medição SFNT. Deduplica por código/tipo/nome/estado/recurso, nunca mensagem. `widget.Label` não editável conserva texto completo e quebra de palavras em `container.NewVScroll`, mínimo configurado de 160 unidades lógicas (não máximo rígido de layout). Dica: **Avisos somente deste componente. Role a lista para consultar os detalhes.** Skin/fonte/versão atualizam somente o rascunho da consulta. A prévia completa recebe o ajuste pelo caminho já existente, sem avisos do documento inteiro multiplicados por componente.

Sem alterações de schemas, catálogo, domínios ou camadas. README da raiz e arquivos de camadas modificados pelo usuário ficam preservados; nenhum hash histórico é substituído. Gates Java/EFL, licenças/recursos e PPT sem embedding continuam inalterados.

### Regressões implementadas — não executadas

| Arquivo | Cobertura preparada, não resultado de execução |
| --- | --- |
| visualprofile/component_diagnostics_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/visualprofile/component_diagnostics_test.go`) | Seleção isolada na mesma/outra tela, rascunho ainda não armazenado, revisão corrente/arquivada exata com caminho compartilhado, dependências/estados futuros, licenças/portabilidade, integridade global inclusive recurso/revisão não usados e ParentID somente local; argumentos intactos |
| desktop/resource_diagnostics_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/resource_diagnostics_test.go`) | Escopo local, métricas SFNT, deduplicação independente de mensagem e estados distintos, revisões exatas, definições/partes/estados futuros versus variante alheia, lista longa rolável sem corte e atualização dos seletores somente no rascunho |

### Roteiro manual pendente — após identificar o novo build

1. Após o executor principal concluir o build, registrar versão, caminho, tamanho e SHA-256 reais do EXE/pacote, data/executor, Windows, resolução, DPI e zoom. Não usar o hash de sobreposição como identidade desta correção. Fechar a instância anterior e abrir a revisão identificada; registrar fixture e perfil/skin/versão efetivamente usados.
2. Preparar um componente saudável e controles com fonte/revisão ausentes na mesma tela e em outra tela. Selecionar o saudável: **Recursos / referência** não deve listar condições alheias. Selecionar cada defeituoso: sua falta real deve continuar visível. Abrir **Prévia** completa com vários componentes e conferir que os avisos do documento inteiro não reaparecem para cada um; não interpretar isso como remoção das falhas próprias ou da validação global.
3. Com fixture de perfil corrente e arquivado compartilhando caminho de fonte/imagem, alternar **Versão do perfil**, **Skin** e **Face da fonte** sem salvar: o texto deve refletir somente rascunho/revisão/definições selecionados. Voltar à opção válida remove a falha correspondente; revisão ausente continua explícita, sem fallback. Conferir documento e histórico inalterados até confirmação e **Descartar** preservando valores confirmados.
4. Usar composto com partes e recurso ausente apenas em estado ainda inativo; conferir diagnóstico pertinente. Uma variante alheia que compartilha fonte não pode contaminar a seleção. Com face real carregada/medida com sucesso via SFNT, não deve permanecer `TEXT_METRICS_REQUIRED`; fonte ou imagem realmente indisponível continua diagnosticada, assim como licença/portabilidade. Mesma condição/localização com mensagens diferentes aparece uma vez; estados distintos permanecem separados.
5. Em fixture controlada, verificar perfil/recurso corrompido fora do uso selecionado: `PROFILE_INVALID` continua bloqueante, sem substituir ou reparar dados. Selecionar filho com pai fora do recorte local: a consulta de recursos não fabrica `PARENT_CLIP_UNRESOLVED` pelo isolamento; ParentID real e a validação de hierarquia fora dessa consulta permanecem preservados.
6. Usar lista longa, incluindo cem avisos e mensagens extensas. Conferir texto completo não editável, quebra de palavras, dica exata, área de referência de 160 unidades e acesso ao último detalhe por rolagem vertical, sem crescimento intrínseco proporcional ao número de avisos. Conferir rodapé fixo, foco e seletores em DPI 100/150/200%. Registrar observado/divergências por caso; o roteiro não constitui aprovação.

### Resultados técnicos observados e verificação manual pendente

| Campo | Estado / resultado efetivamente informado |
| --- | --- |
| Produção e dois novos arquivos de regressões | Implementados; leitura/cobertura preparada não equivale a execução aprovada |
| Formatação e análise estática | PASS: `gofmt` nos quatro arquivos desta correção e `go vet ./...` integral. Arquivos de camadas modificados pelo usuário não foram editados |
| Compilação das suítes | PASS: `go test -c -o ./build/compile-tests/ ./...`; **15 suítes somente compiladas, não executadas** |
| Execução dos testes Go | Não executada por política Windows; sem contorno ou PASS; compilação não substitui execução |
| Testes Python | `python -X utf8 -m unittest discover -s scripts -p test_build.py`: **4 PASS**, reexecutados nesta correção, não os 28 históricos |
| Build/empacotamento Windows | PASS: tarefa VS Code `process: Jornada Studio: compilar Windows` na raiz do workspace e `python -X utf8 scripts/build.py --target windows --arch amd64` no módulo |
| Integridade da distribuição | PASS: CRC, cada entrada SHA256SUMS e seus bytes locais, sidecar, PE amd64, README definitivo, inventário/build-info, cópias idênticas do EXE e ausência de backups `~`; sem executar o aplicativo. Aprovações manual e redistribuição continuam `pending` |
| EXE Windows amd64 | [Executável](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), produto **1.0.0.0**, **34.605.056 bytes**, SHA-256 `634f460781c9a5795205009f6235bc43d5b4a82bc51a707da07b66314a89176c` |
| Cópia de compatibilidade | [Executável alternativo](../../apps/gerador-jornada/dist/JornadaStudio.exe), bytes e hash idênticos ao executável acima |
| ZIP e sidecar | Distribuição (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), **20.406.619 bytes**, SHA-256 `c630b3a163233f2200d3c0e93024d97df8f3f3323ae7484cfee0083c6b6b958d`; sidecar (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256`) correspondente confirmado |
| Data e executor | 2026-10-21, executor principal, após a correção de Recursos / referência; comandos e escopos identificados acima |
| Fixture, perfil/skin/revisão e fontes | A preencher com os dados efetivamente testados |
| Windows, resolução, DPI e zoom | A preencher após teste nativo |
| Resultado por passo, evidência, divergências e reteste | Não executados; registrar aprovado/reprovado/bloqueado somente após observação |
| Aceite explícito | Pendente; T128 e homologações/aceite T120–T122 permanecem abertos |

## 11. Redimensionamento do painel de propriedades — 2026-10-21

**Escopo:** clarificação do redimensionamento/rolagem existentes em [027/FR-003/FR-008](../027-painel-direito/spec.md#functional-requirements) e [028/FR-008](../028-propriedades/spec.md#atualização-solicitada-em-2026-10-18--confirmação-de-alterações), sem novo FR. [T129](tasks.md#phase-23-redimensionamento-do-painel-de-propriedades--2026-10-21) permanece aberta como verificação mista. Correção implementada, analisada, compilada e empacotada conforme checks abaixo; os quatro testes novos não foram executados por política Windows, sem contorno. Os resultados e hashes históricos, inclusive de Recursos / referência, não validam este incremento.

### Causa e correção delimitada

O `Split` Fyne limita o arraste pelo `MinSize` dos filhos. No painel direito, `container.NewVScroll(e.Inspector)` desacoplava somente a altura; a largura mínima do formulário continuava propagada ao painel por seletores longos, caminhos de fontes e versões de perfil. Assim, mudar `Offset` podia não mover visualmente o divisor como esperado.

Em [Editor.design](../../apps/gerador-jornada/internal/adapter/desktop/editor.go), o corpo passa a `container.NewScroll(e.Inspector)` (**ambos os eixos**), desacoplando as duas dimensões mínimas sem truncar valores, perder dados ou impedir acesso horizontal. O painel `right` com `container.NewBorder` permanece com cabeçalho no topo e `InspectorActions` na base, fora da rolagem. **Salvar alterações / Descartar continuam fixos e disponíveis**; Save/Discard, rascunho, documento e histórico não mudam. O mínimo dos controles fixos ainda limita naturalmente o estreitamento: não há largura fixa, botões estreitados ou colapso ilimitado. A rolagem vertical interna de Recursos / referência não é substituída pela correção do corpo externo.

Sem alteração de código nesta atualização documental, schemas, catálogo ou camadas do usuário. README da raiz e registros históricos preservados. Gates Java/EFL, fontes/licenças/recursos e PPT sem embedding continuam inalterados.

### Quatro regressões implementadas — não executadas

Em inspector_split_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_split_test.go`), os testes montam a interface de produção, não uma réplica do layout. Usam janela Fyne de teste **1280 × 720**, com driver em memória; isso não é execução ou aprovação de janela nativa Windows.

| Teste | Cobertura preparada, não resultado observado |
| --- | --- |
| `TestInspectorSplitAllSelections` | Tela, 16 famílias e seleção múltipla; corpo ScrollBoth, mínimos naturais, geometria efetiva do painel em diferentes offsets, ações corretas por seleção e dados intactos |
| `TestInspectorSplitLongValuesAndResize` | Versão builtin extensa, caminhos/fontes e seletores longos; comparar largura/posição efetivas ao mudar Offset e tamanho disponível, conservando valores e rascunho |
| `TestInspectorSplitScrollKeepsDraftHistoryAndFooter` | Alcançar campos nos dois eixos, deslocamento real do conteúdo, ida/volta da rolagem e resize sem perder rascunho ou modificar documento/pilhas; rodapé imóvel |
| `TestInspectorSplitStickySaveDiscard` | Tela e 16 famílias; hit-test dos botões renderizados após rolar, Descartar sem revisão e Salvar em uma transação, preservando dados/overrides longos e rodapé |

### Checks e identidades verificados pelo executor principal

| Verificação / identidade | Estado deste incremento |
| --- | --- |
| Produção e quatro testes | Implementados; não equivale a execução aprovada |
| Data e executor dos checks finais | 2026-10-21, executor principal, após revisão do ajuste e das quatro regressões |
| Formatação e vet integral finais | PASS: `gofmt` em editor.go e inspector_split_test.go; `go vet ./...` integral |
| Compilação integral final das suítes Go | PASS: `go test -c -o ./build/compile-tests/ ./...`; **15 suítes somente compiladas, não executadas** |
| Execução dos testes Go | Não realizada por política Windows; sem contorno ou PASS de testes |
| Testes Python de empacotamento | PASS: **4 testes reexecutados**, `python -X utf8 -m unittest discover -s scripts -p test_build.py`; não é reexecução dos 28 históricos |
| Build Windows e empacotamento | PASS: tarefa VS Code `process: Jornada Studio: compilar Windows` na raiz e `python -X utf8 scripts/build.py --target windows --arch amd64` no módulo |
| Integridade da nova distribuição | PASS: CRC, cada SHA256SUMS e bytes locais, sidecar, PE amd64, README definitivo, inventário/build-info, cópias do EXE idênticas e ausência de backups `~`; aplicativo não executado |
| Novo EXE e cópia de compatibilidade | [EXE Windows](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) e [cópia](../../apps/gerador-jornada/dist/JornadaStudio.exe) idênticos; produto **1.0.0.0**, **34.605.056 bytes**, SHA-256 `3ca955c811402faa1646d69dfd2aa7e29b8fb005e3f3a9e76e4bf82d80df036d` |
| Novo ZIP e sidecar | ZIP Windows (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`), **20.406.978 bytes**, SHA-256 `d0dd117bd2d451aa812d1fbd1f29e4d8766e2f12fc2726d40bddf1eab80f36c2`; sidecar (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256`) correspondente confirmado |
| Windows nativo, DPI e aceite | Não realizados; `manualApproval` e `redistributionApproval` permanecem pendentes |

### Roteiro manual pendente — somente após identificar o novo build

1. Aguardar o build desta correção. Registrar data/executor, revisão do código, comandos/resultados finais e identidade real do EXE/ZIP: caminho, produto/arquitetura, tamanho e SHA-256. Conferir sidecar/integridade e fechar a instância antiga antes de abrir a nova. Não identificar esta correção por hashes das seções anteriores. Registrar Windows, resolução, DPI, zoom e fixture/perfil/skin/revisão/fontes usados.
2. Restaurar a janela para **1280 × 720** e preparar tela com exemplos das 16 famílias. Sem componente selecionado, conferir propriedades da tela; selecionar cada família e depois seleção múltipla. A seleção múltipla deve mostrar seus comandos próprios, sem Salvar/Descartar individuais fictícios. Registrar as larguras/posições visíveis antes/depois de arrastar o divisor em ambos os sentidos: alteração de `Offset` isoladamente não comprova movimento real.
3. Usar versões de perfil extensas, seletores longos e overrides de caminhos de fontes desconhecidos em fixture controlada, sem encurtar dados nem introduzir fallback. Estreitar/ampliar o painel: o formulário não deve travar o divisor na sua largura intrínseca. Ao atingir o mínimo natural do cabeçalho/rodapé, o limite deve ser respeitado; não exigir largura zero nem botões reduzidos.
4. Rolar o **corpo** horizontal e verticalmente até alcançar campos de aparência e sua borda direita; voltar ao início e repetir após redimensionar. Conferir valores completos, edição acessível e rolagem interna independente de Recursos / referência. Cabeçalho e rodapé permanecem fixos e fora do recorte/rolagem do corpo.
5. Alterar um campo válido e manter o rascunho aberto. Rolar nos dois eixos e mover o divisor: não deve haver publicação, perda/reconstrução dos campos ou alteração do histórico. Clicar **Descartar** no botão visível renderizado após rolar: restaurar dados confirmados sem revisão. Editar novamente e clicar **Salvar alterações** após rolar: confirmar uma única transação com os demais valores/overrides preservados. Conferir Undo/Redo e proteção de rascunho existentes; campo inválido não publica estado parcial.
6. Repetir com propriedades de tela e das 16 famílias, janela maior e DPI Windows **100/150/200%**. Registrar foco, alcance dos campos, movimento efetivo do divisor, geometria fixa do rodapé e cliques em ambos os botões. Salvar/reabrir a fixture confirma preservação dos dados, não persistência de offsets transitórios. Anotar divergências e retestes por caso; testes em memória não substituem esta revisão nativa.

### Registro de execução e aprovação ainda aberto

Checks técnicos e identidades estão registrados acima. Preencher somente após observação manual: fixture/perfil/skin/fonte; Windows/resolução/DPI/zoom; seleção e valores longos usados; larguras/posições antes/depois; alcance horizontal/vertical; rascunho/histórico e hit-test de Salvar/Descartar; resultado/evidência/divergência/reteste por passo e aceite explícito da revisão identificada. **Nenhuma execução Go ou aprovação nativa é atribuída por compilação ou por este roteiro.** T129 e T120–T122 continuam abertos.

## 12. Disponibilidade local e portabilidade de recursos — 2026-10-21

**Escopo:** clarificação adicional de [028/FR-016](../028-propriedades/spec.md#aparência-e-famílias-pós-v1--2026-10-19), sem novo FR, com apresentação definida no [DESIGN](../028-propriedades/DESIGN.md#disponibilidade-local-e-portabilidade--clarificação-de-fr-016-2026-10-21). [T130](tasks.md#phase-24-disponibilidade-local-e-portabilidade-de-recursos--2026-10-21) permanece aberta como verificação mista. Somente a apresentação desktop em `componentDiagnostics` mudou na produção; nove regressões foram acrescentadas, incluindo o subteste de recurso já incorporado e válido com aviso persistido obsoleto. Após a última edição, o executor principal concluiu gofmt, vet integral, compilação integral, quatro testes Python, build e empacotamento com integridade canônica conferida. São checks técnicos, **não execução dos testes Go/aplicativo nem homologação**. Nenhuma evidência das seções anteriores é promovida a PASS deste incremento.

### Sintoma observado e causa delimitada

Conforme análise visual registrada pelo executor principal, a captura fornecida pelo usuário foi **efetivamente visualizada** e contém `PROFILE_RESOURCE_NOT_EMBEDDED` para os caminhos lógicos /azul/botao_aux1_pressed.png e /azul/botao_pressed.png. Isso é evidência observada do sintoma anterior, não execução nativa nem verificação da interface da nova compilação.

Os dois PNGs **já estão incorporados ao aplicativo**. O executor principal conferiu por Python os bytes, SHA-256, assinatura PNG e dimensões dos blobs existentes. A causa não é perda dessas imagens: o aviso técnico de **não incorporação ao arquivo do projeto** estava sendo apresentado como se significasse **indisponibilidade no runtime local**. O isolamento/rolagem de Recursos / referência da seção 10 tratava outro problema e não resolvia essa distinção; o ajuste do divisor da seção 11 também não a tratava.

| Caminho lógico do recurso (não caminho de disco) | Bytes originais | Dimensões PNG | SHA-256 verificado pelo executor principal |
| --- | --- | --- | --- |
| /azul/botao_aux1_pressed.png | 3.530 | 80 × 60 | `bd2a562a7dda9c16bb257368b79f3f859f40c7ee388a84da722a14b4491d1dec` |
| /azul/botao_pressed.png | 3.396 | 52 × 52 | `25007821e119733e4feeaf516e4f15dcce9bd9a0fa3242c85a88e8625da0999b` |

Esses hashes identificam somente os recursos originais, **não EXE, ZIP ou perfil**. Nenhum recurso foi recuperado, recriado, recodificado ou copiado para o projeto por esta mudança.

### Apresentação implementada e invariantes preservadas

Em inspector_appearance.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_appearance.go`), `componentDiagnostics` continua a consultar seleção/rascunho e revisão exata. Para cada `PROFILE_RESOURCE_NOT_EMBEDDED` com caminho, chama `sceneResourceData(document, group.SkinVersion, path)`, o mesmo resolvedor de bytes da pintura, que exige identidade de SHA-256, tamanho e MIME. Caminho igual, versão atual ou nome parecido não substituem a revisão requerida.

- **Bytes exatos resolvidos:** substituir somente a linha de apresentação dessa condição por **Disponível no aplicativo · PATH (não incluído no arquivo do projeto)**, uma vez por caminho, ainda que citado por vários estados/origens. `PATH` é o caminho lógico efetivamente consultado.
- **Uma ou mais linhas disponíveis:** acrescentar uma única mensagem comum: **Aviso de portabilidade, não de carregamento: os recursos acima dependem deste acervo exato do aplicativo. Não foram copiados para o projeto; disponibilidade não comprova licença de redistribuição nem fidelidade integral.**
- **Não resolvido:** conservar `PROFILE_RESOURCE_NOT_EMBEDDED`, localização e caminho, com **Recurso indisponível nesta instalação e não incluído no projeto; nenhuma substituição foi aplicada. Portabilidade e licença de redistribuição não comprovadas.**
- **Já incorporado e validado no projeto:** não reapresentar um aviso persistido obsoleto de não incorporação nem descrevê-lo como dependência externa; manter o perfil e as demais condições, inclusive legais, intactos.
- **Demais condições:** manter falhas reais de imagem/fonte, `PROFILE_INVALID`, recursos ausentes/bloqueados e avisos legais. Não converter ausência ou corrupção em disponibilidade nem esconder condições bloqueantes.

Esta é uma alteração de apresentação, não um novo modelo de recursos: `visualprofile.Diagnostics`, `visualprofile.ComponentDiagnostics`, `visualprofile.Capture` e diagnósticos persistidos mantêm códigos/mensagens técnicas. Perfil, hashes, recursos, licenças/status, documento e histórico não são modificados pela consulta. Não há cópia de bytes para o projeto, licença nova, alteração de schema/catálogo ou fallback por homônimo. Validação estrita e limites do PPT, inclusive ausência de embedding de fontes, continuam inalterados; Java/EFL permanece desabilitado. Disponibilidade local não comprova portabilidade autônoma, redistribuição autorizada ou fidelidade integral.

Texto não editável, quebra de palavras, rolagem vertical de referência de 160 unidades, dica existente, corpo bidirecional e rodapé fixo são preservados. Alterar Skin/Face da fonte/Versão do perfil atualiza a consulta do rascunho, sem confirmar alterações. Esta finalização altera somente a [spec de propriedades](../028-propriedades/spec.md), seu [DESIGN](../028-propriedades/DESIGN.md), [tasks.md](tasks.md) e este quickstart. [README da aplicação](../../apps/gerador-jornada/README.md) já empacotado e congelado, [README da raiz](../../README.md), código/testes, configuração de tarefas e artefatos permanecem intocados; nenhum reempacotamento é realizado nesta edição documental.

### Nove regressões implementadas — não executadas

Em resource_availability_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/resource_availability_test.go`), a cobertura preparada é:

| Teste | Cobertura preparada, não resultado de execução |
| --- | --- |
| `TestDisponibilidadeBuiltinCapturaEPacote` | Default e variantes dos botões; builtin/captura/salvar-reabrir, estados futuros, APIs técnicas e avisos legais preservados, sem copiar bytes ou alterar pacote determinístico |
| `TestDisponibilidadePNGExato` | SHA-256, tamanho, decodificação PNG e dimensões dos dois recursos pressionados originais |
| `TestDisponibilidadeDeduplicaCaminhosNaoAvisosLegais` | Uma linha por caminho e um aviso comum, apesar de múltiplos estados/origens; ressalvas jurídicas não reescritas |
| `TestDisponibilidadeReferenciaDivergenteNaoUsaHomonomo` | Recusa de hash/tamanho/MIME divergentes e referência própria homônima; código técnico e indisponibilidade mantidos |
| `TestDisponibilidadeRevisaoExataNaoMisturaMetadados` | Alternância de revisão ativa/arquivada com mesmo caminho e identidades divergentes, inclusive caches |
| `TestDisponibilidadeEstadoFuturoDistingueReferencias` | Mistura de referência disponível e referência não resolvida em estado ainda inativo, sem ocultar esta última |
| `TestDisponibilidadeAusenciaBuiltinConservaAvisos` | Ausências reais do inventário continuam diagnosticadas com localização/mensagem originais, antes/depois de Capture |
| `TestDisponibilidadeFalhasReaisConservamCodigos` | Perfil inválido, recursos ausentes/bloqueados e bytes ausentes/corrompidos não viram sucesso; bytes já incorporados e validados não são descritos como externos por aviso obsoleto, sem apagar ressalvas legais |
| `TestDisponibilidadeFormularioAtualizaRascunho` | Skin/revisão atualizam a mesma legenda sem publicar dados; aplicação explícita continua separada e origem/pacote permanecem intactos |

### Checks técnicos finais — observados pelo executor principal

| Verificação / evidência | Estado efetivamente informado neste incremento |
| --- | --- |
| Captura do usuário | Efetivamente visualizada pelo executor principal; código e dois caminhos identificados no sintoma anterior. Não é UI da nova revisão verificada |
| Bytes dos dois PNGs já incorporados | Conferência Python pelo executor principal aprovada: SHA-256, tamanho, assinatura PNG e dimensões da tabela acima |
| Produção e nove regressões | Implementadas, incluindo o subteste de aviso persistido obsoleto em recurso embedded e validado; incluídas na compilação final, **não executadas** |
| Data e executor dos checks finais | 2026-10-21, executor principal, após a última edição dos nove testes e do subteste de aviso obsoleto; sem nova revisão Git informada |
| Formatação final | PASS: gofmt nos inspector_appearance.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/inspector_appearance.go`) e resource_availability_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/resource_availability_test.go`) |
| Vet integral final | PASS: `go vet ./...`, marcador observado `RESOURCE_FINAL_VET_PASS` |
| Compilação integral final das suítes Go | PASS: `go test -c -o ./build/compile-tests/ ./...`, marcador observado `RESOURCE_FINAL_COMPILE_PASS`; **15 suítes somente compiladas, ZERO execuções** |
| Execução das suítes Go e aplicativo novo | Testes Go não executados por política Windows, sem contorno nem PASS de execução; aplicativo novo não iniciado nesta finalização. Compilação e leitura do EXE não substituem execução |
| Testes Python de empacotamento | `python -X utf8 -m unittest discover -s scripts -p test_build.py`: **4 PASS**, executados neste incremento; não são os 28 históricos nem a conferência dos PNGs |
| Revisão Explore somente leitura | Nenhum defeito concreto apontado; **não é teste executado**, execução Go ou revisão de UI nativa |
| Build Windows pela tarefa VS Code | PASS: `process: Jornada Studio: compilar Windows`, na raiz do workspace, terminou com **exit 0** |
| Empacotamento Windows posterior à tarefa | PASS: `python -X utf8 scripts/build.py --target windows --arch amd64`, no módulo, **sem `--check`**, concluiu o novo pacote |
| Leitura independente da distribuição canônica | PASS: CRC de todas as entradas, SHA256SUMS de todos os bytes arquivados e comparação com o output local, conjunto exato de entradas, nenhum backup terminado em `~` e sidecar SHA-256 do ZIP |
| Manifesto, metadados e PE | PASS: manifesto legado exato e seu SHA-256 em build-info, produto **1.0.0.0**, alvo **Windows amd64** e cabeçalho **PE AMD64** |
| Texto do novo aviso e README congelado no pacote | PASS: literal **Aviso de portabilidade, não de carregamento:** dentro do novo EXE; README empacotado corresponde ao guia congelado, normalizando somente fins de linha e referência do ícone do pacote |
| Cópia de compatibilidade — estado atual | **RESOLVIDO / PASS** após o usuário informar “fechei a aplicação”: `Copy-Item` do canônico para a cópia e `Get-FileHash` independente em ambos; **34.608.128 bytes e SHA-256 idênticos**, marcador `COMPATIBILITY_COPY_UPDATED_AND_VERIFIED_PASS SHA256=48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d` |
| ZIP após a resolução da cópia | **Não refeito**; mantém 20.408.876 bytes, SHA-256 `b40cc96c901b51d00ca042831c419c222f44f90709b7d0328baf70695d8bdf01` e os checks de integridade já PASS; README congelado e pacote não alterados nesta finalização |
| Teste completo nativo, DPI e Office | Pendentes; análise da captura anterior não equivale a nova execução nativa |
| Aceite e redistribuição | `manualApproval=pending` e `legacySkin.redistributionApproval=pending` conferidos; T130 e T120–T122 continuam abertos, sem promoção de aprovação por checks técnicos |

**Histórico da primeira cópia e checagem ad hoc — preservado, não estado atual:** inicialmente somente o EXE canônico e o ZIP foram atualizados; a primeira cópia de compatibilidade foi bloqueada por estar em uso e conservou a revisão anterior, **34.605.056 bytes**, SHA-256 `3ca955c811402faa1646d69dfd2aa7e29b8fb005e3f3a9e76e4bf82d80df036d`. A checagem ad hoc inicial **falhou** por exigir incorretamente igualdade entre o canônico e essa cópia naquele momento; não foi falha de build ou do ZIP. A leitura independente posterior validou os bytes da distribuição canônica sem essa suposição. O empacotador informou então: **O EXE no caminho antigo está em uso; a distribuição atualizada está em dist/windows-amd64.** Naquela etapa não houve atualização nem PASS de compatibilidade; mensagem e hash anteriores ficam somente como registro histórico.

**Resolução posterior pelo executor principal (2026-10-21) — RESOLVIDO:** o usuário informou **“fechei a aplicação”**. O executor principal então realizou `Copy-Item` do [canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) para a [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) e `Get-FileHash` independentemente em ambos, com PASS: `COMPATIBILITY_COPY_UPDATED_AND_VERIFIED_PASS SHA256=48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d`. Ambos agora têm **34.608.128 bytes e o mesmo SHA-256**; a cópia não permanece antiga ou bloqueada. Nenhum aplicativo ou teste Go foi iniciado nessa operação ou nesta edição documental; nenhum processo foi encerrado pelo agente, nenhum backup foi removido e nenhuma política foi alterada. **O ZIP não foi refeito** e seus checks já PASS permanecem válidos; a atualização documental não altera o README congelado nem exige reempacotamento.

<a id="identidade-atual-da-distribuição-canônica-e-cópia-antiga-bloqueada"></a>

### Identidade atual da distribuição canônica e cópia de compatibilidade

| Artefato | Situação | Tamanho | SHA-256 |
| --- | --- | --- | --- |
| [apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | **Canônico atualizado**, produto 1.0.0.0, Windows amd64 | **34.608.128 bytes** | `48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d` |
| [apps/gerador-jornada/dist/JornadaStudio.exe](../../apps/gerador-jornada/dist/JornadaStudio.exe) | **Compatibilidade atualizada e verificada — PASS**, idêntica ao canônico após o usuário fechar a aplicação | **34.608.128 bytes** | `48cdf163c61be34f46b035bd9fb29d1f1a89dc90da24debfc475c728d87a6a0d` |
| apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | **ZIP preservado, não refeito** após a cópia; contém o EXE canônico novo, checks já PASS e sidecar (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip.sha256`) correspondente confirmado | **20.408.876 bytes** | `b40cc96c901b51d00ca042831c419c222f44f90709b7d0328baf70695d8bdf01` |

Os histogramas, PASS, tamanhos e hashes históricos das seções 6–11 permanecem intactos e exclusivos de suas revisões. **Ambos os caminhos de entrega do EXE agora contêm a revisão acima**; a identidade da seção 11 não corresponde mais a nenhum desses dois arquivos atuais. O ZIP mantém os bytes e a identidade já verificados antes da resolução da cópia. O README congelado foi conferido no pacote, não editado nesta finalização. Não é necessário repetir cópia ou reempacotamento; os dois caminhos são equivalentes para um teste manual futuro, ainda pendente.

**Aviso de link histórico preservado:** o link de reconciliação documental de 2026-10-18 em [Evidências e interpretação do backlog](tasks.md#evidências-e-interpretação-do-backlog--reconciliação-2026-10-18) continua apontando para a âncora antiga quebrada deste quickstart. A referência histórica não foi removida nem corrigida; isso não afeta os novos links de identidade e resolução.

### Roteiro manual pendente — após identificar o novo build

1. Para teste manual futuro, conferir a identidade real contra a tabela atual acima e o sidecar, registrando data/executor, comandos/resultados e eventual revisão de código disponível. O usuário já informou que fechou a aplicação, e a cópia foi resolvida; o [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) e a [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) **agora são idênticos**, sem nova cópia ou reempacotamento exigidos. Abrir um deles somente na execução manual posterior; nenhum executável é iniciado nesta finalização. Registrar Windows, resolução, DPI, zoom e fixture/perfil/skin/revisão/fontes; não usar hashes históricos como identidade atual.
2. Selecionar as variantes dos dois botões da captura e consultar **Recursos / referência**. Conferir a linha disponível de cada PNG pertinente, inclusive quando só usado em estado inativo, sem linha bruta de `PROFILE_RESOURCE_NOT_EMBEDDED` para o mesmo recurso localmente resolvido. Outras dependências disponíveis podem aparecer; cada caminho aparece uma vez e o aviso comum apenas uma vez. Isso não aprova visualmente o estado pressionado nem o acervo integral.
3. Capturar/salvar/reabrir a fixture e conferir apresentação coerente com a mesma revisão e metadados. Inspecionar os recursos do projeto: nenhuma cópia automática, mudança de licença/status ou promoção dos diagnósticos persistidos. APIs técnicas continuam informando a limitação de portabilidade. Salvar/reabrir não autoriza redistribuição.
4. Em fixtures controladas, alternar revisão corrente/arquivada com caminho igual e hash/tamanho/MIME diferentes. Conferir recusa de homônimo e mensagem de indisponibilidade com código/localização preservados, inclusive quando só um estado futuro não resolve. Misturar caminhos disponíveis e indisponíveis: somente os primeiros recebem a apresentação amigável.
5. Conferir faltas reais, recursos bloqueados, bytes ausentes/corrompidos e perfil inválido, inclusive integridade fora do uso imediato. `PROFILE_INVALID`, falhas de carregamento e avisos legais continuam visíveis e bloqueios permanecem; não aplicar substituições. Verificar limites do PPT estrito sem tratar disponibilidade como licença ou embedding implementado; homologação Office continua separada.
6. Alterar Skin/Face da fonte/Versão do perfil sem salvar e conferir consulta da revisão exata, preservando rascunho/documento/histórico. Conferir Descartar e confirmação explícita existentes, rolagem/texto completo/dica/rodapé e divisor em DPI **100/150/200%**. Registrar evidência, observado, divergências e retestes por caso; nenhum passo planejado constitui execução ou aceite.

**Conclusão técnica, homologação pendente:** causa observada e PNGs conferidos; apresentação/nove regressões com subteste de aviso obsoleto implementados; gofmt, vet integral, compilação das 15 suítes sem execução, quatro testes Python, build Windows e pacote com integridade canônica PASS, conforme executor principal. **Compatibilidade RESOLVIDA / PASS após o usuário fechar a aplicação: ambos os EXEs agora têm 34.608.128 bytes e o mesmo SHA-256 da tabela.** Primeira cópia bloqueada e checagem inicial falha preservadas como histórico, sem falha de build/ZIP. ZIP não refeito, README congelado e checks já PASS preservados. **ZERO execuções Go/aplicativo novo nesta finalização**, sem encerrar processos pelo agente, alterar ou contornar a política Windows; UI nativa, DPI 100/150/200%, Office e aceite ainda pendentes. Captura anterior e revisão Explore somente leitura não validam a UI nova. T130 permanece aberta como verificação mista; T120–T122, aprovações `pending` e gates Java/EFL, fontes/licenças/recursos e PPT sem embedding ficam inalterados. Nenhuma tarefa é encerrada por esta edição documental.

## 13. Log de aplicação — 2026-10-21

**Escopo:** [001/FR-057](spec.md#log-de-aplicação--2026-10-21), [T131 aberta](tasks.md#phase-25-log-de-aplicação--2026-10-21) e [clarificação de 028/FR-016](../028-propriedades/spec.md#log-da-consulta-visível--clarificação-de-fr-016-2026-10-21), sem requisito duplicado em Propriedades. A falha continua relatada pelo usuário; nenhuma causa nova foi comprovada e **não há declaração de solução**. A seção 12, inclusive resolução da cópia de compatibilidade, checks e hashes, é histórica e permanece intacta: igualdade de EXEs antigos não comprova uso do logger novo.

### Implementação concluída para coleta, não evidência de execução

**Escopo técnico concluído pelo executor principal:** 11 arquivos Go formatados: sete existentes (bootstrap, editor, inspector, aparência, prévia, arquivos e modo de janela) e quatro novos (logger, UI e dois arquivos de testes). Produto **1.0.0.0**, Windows **amd64**, schemaVersion **3** e catalogVersion **2** inalterados. A conclusão da produção não comprova reprodução da falha ou funcionamento real do logger.

- [application_log.go](../../apps/gerador-jornada/internal/adapter/desktop/application_log.go): JSONL por processo/sessão, startup anterior ao Fyne; timestamp UTC, PID/sessão, basename/SHA-256 do EXE, versão do produto/Go, SO e arquitetura. Erros apresentados por diálogo/status são registrados com classe/tipo seguro, errno/código de perfil quando reconhecidos e fingerprint, não mensagem arbitrária bruta.
- Local: **%LOCALAPPDATA%/GeradorJornada/logs** no Windows, via `os.UserCacheDir`, ou **os.TempDir()/GeradorJornada/logs** se cache indisponível/não gravável; usar o caminho efetivo mostrado, não um destino presumido. Nome por timestamp/PID; escrita síncrona com mutex/Sync, até **2 MiB por segmento + 3 backups**, retenção de **20 arquivos próprios removíveis**, sem forçar remoção de bloqueados. Dedup de snapshots: **2 segundos/256 identidades**, sem impedir captura explícita.
- Recursos / referência registra automaticamente a consulta da seleção/rascunho e mudanças de Skin/Face da fonte/Versão do perfil. O popup **Log de aplicação**, no cabeçalho, mostra caminho e **Registrar diagnóstico agora / Copiar caminho / Abrir pasta**. Abrir o popup já força captura; o botão repete o rascunho visível ou, sem seleção, os componentes da tela atual. O botão de diagnóstico da prévia também registra suas consultas, sem salvar ou alterar documento/histórico.
- Cada diagnóstico conserva código/tipo/nome/estado/recurso sob allowlist: valores builtin conhecidos e mensagens exatas permitidas, desconhecidos/customizados em hash ou omitidos. `displayed` indica somente presença do código no relatório visível; `report_sha256` correlaciona esse texto **sem relatório bruto**. Diagnósticos técnicos continuam registrados mesmo quando a UI apresenta disponibilidade amigável em vez do código.
- Resolução usa a revisão exata, inclusive arquivada e partes/estados futuros: SHA-256/tamanho/MIME esperados, status/licença; hash/tamanho observados e resolvido/código de erro. Não troca revisão por homônimo, cria fallback, copia recursos, muda políticas ou concede licença.
- Privacidade: sem conteúdo de entradas/senhas, nomes de projeto/componente ou `FilePath`; mensagens, caminhos e metadados personalizados desconhecidos em hash/omissão. Standard log/Fyne registra somente hash/bytes. Sem rede. `Err` distingue falha do logger de erro da aplicação registrado; falha de gravação fica visível/consultável na UI.

### Verificação mista — checks técnicos PASS, coleta real PENDENTE

| Item | Estado deste incremento LOG |
| --- | --- |
| Produção / GoDoc / formatação | Implementação concluída no desktop/bootstrap, comentários em português; **11 arquivos Go formatados**, sem execução do aplicativo |
| Nove testes de núcleo | [application_log_test.go](../../apps/gerador-jornada/internal/adapter/desktop/application_log_test.go): identidade, redação, mensagens builtin, resolução, dedup/force, concorrência/rotação/retenção, falhas/nil/fechamento, revisão arquivada/imutabilidade e cache/fallback; **SOMENTE COMPILADOS, NÃO EXECUTADOS** |
| Três testes de UI | application_log_ui_test.go (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/application_log_ui_test.go`): relatório/rascunho/revisão, integração/captura/caminho/erro e compatibilidade/falha visível; **SOMENTE COMPILADOS, NÃO EXECUTADOS**, sem homologação Windows nativa |
| Vet integral | **PASS**, executor principal observou `APPLICATION_LOG_VET_PASS` |
| Compilação integral das suítes | **PASS**, `APPLICATION_LOG_COMPILE_ONLY_PASS`: **15 suítes, ZERO execuções Go**; não é PASS dos testes |
| Execução Go | Não realizada; bloqueio prévio por política Windows preservado **sem contorno** |
| Python de empacotamento | **4 PASS** em [test_build.py](../../apps/gerador-jornada/scripts/test_build.py) |
| Build Windows | **PASS**, tarefa `process: Jornada Studio: compilar Windows` no workspace raiz, **exit 0** |
| Empacotamento | **PASS**, `python -X utf8 scripts/build.py --target windows --arch amd64`, **sem `--check`** |
| Integridade somente leitura | **PASS**: todos os CRC/SHA256SUMS e bytes locais; inventário exato sem backups `~`; manifesto idêntico ao incorporado e hash coerente no build-info; PE AMD64/ícone; README congelado exato após normalização CRLF; sidecar SHA-256 do ZIP |
| Identidade dos artefatos | **PASS**, ambos os EXEs atualizados/idênticos e ZIP com tamanho/SHA-256 reais na tabela abaixo |
| Presença do logger no binário | **PASS somente estático** para os cinco marcadores: `Log de aplicação`, `component_snapshot`, `resource_resolution`, `executable_sha256`, `report_sha256`; **não prova execução** |
| Startup/runtime do logger e erro real correlacionado | **PENDENTES**: nenhum aplicativo iniciado nesta entrega, nenhum log real coletado; não afirmar que o usuário já abriu o novo EXE |
| UI nativa, DPI 100/150/200%, Office, reteste e aceite | **PENDENTES**; T131/T120–T122 abertas, aprovações `pending` e gates Java/fontes/PPT inalterados |

Checks e identidades acima foram observados pelo **executor principal**, não reexecutados nesta finalização documental. `manualApproval` e `redistributionApproval` continuam **`pending`**; T130 também permanece aberta. Preservar as seções 1–12 e seus hashes históricos; não usar seus PASS como evidência deste incremento.

### Identidades reais da entrega LOG — 2026-10-21

| Artefato atualizado | Tamanho em bytes | SHA-256 |
| --- | ---: | --- |
| [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | 34693632 | `ebfc1c7a6051d8512718e47e724fd543b68face33819882f62aa12d8ac172f84` |
| [EXE de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) | 34693632 | `ebfc1c7a6051d8512718e47e724fd543b68face33819882f62aa12d8ac172f84` |
| ZIP Windows amd64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 20442016 | `6ac2fb80d9d6a54a02fb26db1729fde05e762b7ee7f5a5babf212d9b34359b90` |

As duas cópias do EXE são **iguais e atualizadas**, sem divergência de entrega. O ZIP e seu sidecar foram conferidos. O [README da aplicação](../../apps/gerador-jornada/README.md) está **congelado e já empacotado**, com texto exato conferido após normalização CRLF; não editar nem reempacotar por esta atualização dos seis documentos. Hash/bytes/PE/marcadores não identificam por si sós o processo que o usuário tem aberto: confirmar a identidade pelo startup real.

### Coletar o erro real — roteiro ainda não executado

1. Usar a entrega LOG identificada na tabela acima: checks/build/empacotamento técnicos já PASS. Manter README/código/testes/artefatos intocados; não usar hashes da seção 12 nem reempacotar para a coleta. O objetivo é diagnosticar a falha ainda relatada, não declarar solução a partir dos checks.
2. **Fechar a instância antiga e abrir o novo EXE identificado**, preferencialmente o canônico da tabela. Abrir **Log de aplicação**, conferir arquivo/caminho efetivo e correlacionar startup com SHA-256/basename, PID/sessão/versão/Go/SO/arquitetura. Comparar `executable_sha256` ao hash desta entrega. Se não existir arquivo/startup, registrar falha/bloqueio; não criar um log fictício nem presumir que o EXE antigo já o produza. Esta abertura é instrução ao usuário, não execução realizada ou comprovada nesta entrega.
3. Reproduzir a falha em **Recursos / referência**. Selecionar o componente e alternar Skin/Face da fonte/Versão do perfil, inclusive revisão arquivada pertinente, mantendo o rascunho sem salvar. Registrar observado e horário; em **Log de aplicação**, clicar **Registrar diagnóstico agora** mesmo dentro de 2 s. Conferir correspondência de `report_sha256`/`displayed`, códigos técnicos originais e resolução esperado/observado, sem inferir sucesso só da mensagem amigável.
4. Usar **Copiar caminho / Abrir pasta** e coletar o arquivo indicado, incluindo segmentos/backups pertinentes se houver rotação. Sem seleção, conferir captura dos componentes da tela atual; conferir também o diagnóstico da prévia e eventual erro de diálogo/status. Se a escrita falhar, registrar a falha exibida e `Err`, não uma captura bem-sucedida presumida. Não publicar rascunho para obter diagnóstico.
5. Antes de compartilhar voluntariamente o log, revisar a privacidade com fixture controlada: entradas/senha/nomes/`FilePath`, erro/mensagem e metadados personalizados não podem aparecer em bruto; confirmar hash/bytes do standard log. Registrar ambiente/DPI/zoom, perfil/revisão por identidade técnica, erro observado, evento/fingerprint, evidência, divergência e reteste, sem anexar conteúdo privado do projeto ao log.

**Checks técnicos PASS; evidência real PENDENTE:** ainda não há arquivo de startup/erro proveniente do novo EXE, execução dos 12 testes ou aceite. A falha continua relatada; **causa não mapeada e problema não resolvido de forma comprovada** até reproduzir/enviar o log. O logger é instrumento de diagnóstico, não correção demonstrada; não criar logs artificiais nem afirmar uso do EXE novo sem evidência. Âncora histórica de reconciliação já quebrada permanece avisada/preservada; validação UTF-8/links desta finalização fica a cargo do executor principal, sem apagar o histórico.

<a id="14-builtin-filtrado"></a>
## 14. Builtin filtrado por dependências indisponíveis — 2026-10-21

**Estado: ENTREGA TÉCNICA CONCLUÍDA, homologação pendente.** Decisão do usuário: “pra ficar facil oq não achar exclui caso eu ter necessidade eu faço uma importação em uma change”. Escopo em [001/FR-058](spec.md#builtin-filtrado-fr-058), clarificação em [003/FR-012](../003-catalogo/spec.md#builtin-filtrado-fr-012); [001/T132](tasks.md#t132-builtin-filtrado) e [003/T017](../003-catalogo/tasks.md#t017-builtin-filtrado) permanecem **abertas como verificação mista**. Os registros 214/47 e hashes das seções anteriores são históricos antes da redução; esta seção registra verificações e identidades reais do builtin filtrado.

### Geração existente e preservação

| Item | Estado da geração filtrada |
| --- | --- |
| Skin version | `legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a` |
| Catálogo operacional | **206 definições ativas / 8 fontes / 270 recursos / 31 diagnósticos ativos**, nas mesmas 16 famílias |
| Adiamentos auditáveis | **8 definições / 16 referências missing**; **206 + 8 = 214 definições originais**, sem apagar legado |
| Manifesto | Campo aditivo `deferred_definitions`: `type`, `name`, `definition`, `reasons`; `source_tree` e `xml_sha256` originais preservados, assim como origens, conflitos, seleções/pins, bytes disponíveis e fontes |
| Versões inalteradas | **manifest schema 1 / Profile schema 1 / project schemaVersion 3 / catalogVersion 2 / produto 1.0.0.0** |
| EXE/ZIP desta skin | **Gerados e verificados**, identidades reais abaixo; não equivalem a homologação |

O [importador](../../apps/gerador-jornada/scripts/import_legacy_skin.py) adia **a definição inteira**, incluindo estados futuros, quando `imagePath`, `slider` ou `font:type` aponta para recurso ausente/inválido, caminho recusado ou tipo incompatível. Se um filho falha, pai e todos os filhos são adiados, pela identidade tipada/vínculo real de composição; não por basename/prefixo nem homônimos de outro tipo. As oito variantes `button`/`label`/`panel` estão nomeadas em [FR-058](spec.md#builtin-filtrado-fr-058); elas não são opções operacionais do novo builtin, mas continuam no manifesto com as razões originais.

Novos projetos usam o builtin filtrado **completo**, não subset de metadados por uso. Perfis próprios/privados, revisões arquivadas e documentos existentes são preservados sem nova migração, poda, enriquecimento, upgrade ou fallback builtin. Fontes/bytes disponíveis não são removidos por estarem sem uso. A redução não recupera os 16 recursos ausentes, não sana delegates/licenças e não elimina `DELEGATE_MISSING`/`FIDELITY_UNVERIFIED` ou gates Java/PPT; nenhuma promessa de 100% de fidelidade ou autorização de redistribuição.

### Evidência técnica do executor principal

| Verificação | Resultado informado ou pendência deste incremento |
| --- | --- |
| Python do importador | **19 testes PASS**, executados com `unittest discover -s scripts -p test_import_legacy_skin.py -v`; inclui todos os quatro métodos novos |
| Descoberta dos quatro novos métodos | **PASS explícito** por `getTestCaseNames`, total 19 e quatro nomes exigidos no nível da classe. Aninhamento inicial corrigido; a primeira execução de 15 casos era parcial, não evidência desses quatro |
| Python do empacotamento | **4 testes PASS** de `test_build.py`; testes Java históricos não reexecutados neste incremento |
| Importador `--check` | **PASS**, geração filtrada idêntica, sem alterações |
| Integração Go | Regressões de biblioteca (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/skin/builtin_filter_test.go`) e seletores/Capture/perfis próprios e arquivados (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/desktop/builtin_filter_test.go`) implementadas; cenário anterior incompleto preserva avisos, sem usar variante excluída como builtin |
| gofmt / vet / compile-only | **PASS**: quatro testes Go formatados; `go vet ./...` e `go test -c -o ./build/compile-tests/ ./...`, **15 suítes somente compiladas, ZERO execuções Go** |
| Execução de testes Go | **NÃO REALIZADA**, política Windows preservada sem contorno; compile-only não é execução/PASS de regressões |
| Build Windows / empacotamento | **PASS**: tarefa VS Code `process: Jornada Studio: compilar Windows` na raiz e `python -X utf8 scripts/build.py --target windows --arch amd64`, sem `--check` |
| Integridade do inventário | **PASS somente leitura**: 549 referências XML ativas resolvem por caminho/tipo, todos os 270 recursos/264 blobs têm bytes/hash/tamanho exatos e todas as origens coincidem; `source_tree`/`xml_sha256` originais e pin da imagem clara preservados; 8 adiadas/16 referências nos motivos, nenhuma falta operacional publicada |
| Integridade da entrega | **PASS somente leitura**: PE amd64, nova skin e todos os blobs disponíveis no EXE, dois EXEs idênticos, CRC de todas as entradas, inventário exato de seis arquivos sem backup `~`, todos os SHA256SUMS/bytes locais, manifesto/build-info, ícone, README congelado normalizado CRLF e sidecar |
| Inserção/seletor, projeto anterior/próprio/arquivado, UI nativa/DPI, Office e aceite | **PENDENTES**, roteiros abaixo; `manualApproval`/`redistributionApproval` continuam `pending` |

### Identidade real dos artefatos — não executados

| Artefato | Bytes | SHA-256 |
| --- | ---: | --- |
| [EXE Windows canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) e [cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe), idênticos | 34.690.560 | `87fa3e1876d2b4efe0d0c3ba71c254d6614e80f944f8dfffe2961ab635eaa0cb` |
| ZIP Windows (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 20.442.829 | `36dc8d1496b0ffe32a8e8174687cbdad44385a8199c55706d1d0bf1d993760d4` |
| Manifesto (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/legacyskin/assets/manifest.json`) | 1.039.093 | `60713ea85e1b27c38fd8ef37d9c4c74a9c1bbe074bc0269a00e2da8b84a6b3fb` |

Sidecar externo `.zip.sha256` conferido. O [README da aplicação](../../apps/gerador-jornada/README.md) foi atualizado antes do pacote e está congelado, SHA-256 `58d4e5ea19f015a1a2aa9e042d51a4d98d6479ed80faf94d07a8b22757bc8284`; não editar sem reempacotar. Ambos os caminhos do EXE agora contêm esta revisão; hashes das entregas anteriores não identificam seus bytes atuais. Produto **1.0.0.0**, Windows amd64; `manualApproval` e `legacySkin.redistributionApproval` permanecem `pending`. Nenhum aplicativo/teste Go foi executado ou processo encerrado; sem novos builds Linux/macOS.

### Roteiro pendente — nova inserção e seletor 206/8

1. Usar o EXE/ZIP identificado acima e registrar ambiente, DPI/zoom e skin. Checks/build/pacote estão concluídos, mas abertura/uso da nova revisão não foram comprovados; não reutilizar o SHA-256 da entrega LOG como aprovação deste incremento. Versão de produto 1.0.0.0 sozinha não distingue revisões.
2. Criar projeto novo e tela; conferir o perfil builtin filtrado e o total de **206 definições ativas**, incluindo partes/variantes. Abrir os seletores por família, inserir controles válidos e trocar variantes ativas; **206 não é a quantidade de opções de cada seletor** nem número de famílias.
3. Conferir ausência das oito variantes de FR-058 para novas inserções/trocas nesse builtin e presença das oito projeções/motivos em `deferred_definitions`. Reconciliar **206 ativas + 8 adiadas = 214 originais**, **16 referências missing**, 8 fontes, 270 recursos e 31 diagnósticos ativos; o XML/acervo original não pode ter sido alterado. A contagem menor de avisos não é recuperação de recurso.
4. Nas fixtures controladas do executor principal, verificar falha em estado não atual, `slider` fora de estado, `font:type`, caminho inválido e tipo incompatível: nenhuma definição parcial; falha de filho retira toda composição, não deixa órfãos nem atinge homônimos de outro tipo/prefixos independentes. Não editar o legado para o teste manual.
5. Salvar/reabrir o projeto novo e **inserir outra família/trocar variante não usada antes de salvar**; conferir biblioteca ativa completa, estados e referências, sem subset por uso. Comparar canvas/prévia e diagnósticos; não suprimir avisos legais/fidelidade remanescentes.

### Roteiro pendente — preservar projeto antigo, perfil próprio e revisão arquivada

1. Separar cópia segura de projeto anterior à redução e registrar seu hash original/perfil/revisões. Incluir caso com variante agora adiada e casos de perfil próprio e revisão arquivada exata; usar fixtures existentes/controladas, não converter o original para o novo builtin.
2. Abrir sem regravar e conferir definições/estados/fontes/aliases/recursos/revisões completos, conteúdo, UUIDs/FileIDs, bounds, overrides, vínculos e ordem. Ausência ou corrupção deve continuar diagnosticada pela revisão original; não substituir por homônimo/builtin novo para esconder falha.
3. Salvar **uma cópia**, reabrir, inserir variante válida desse perfil e comparar a integridade/snapshot sem poda, enriquecimento, upgrade, migração extra ou fallback. Confirmar original intacto e metadados não usados/arquivados preservados. Se recurso exato não resolver, registrar bloqueio: preservar não equivale a prometer disponibilidade/fidelidade.

### Reimportação em futura change

1. Abrir change específica se houver necessidade das variantes adiadas; identificar definição/composição e cada referência/estado em `deferred_definitions` e no XML original.
2. Obter **recursos exatos**, com origem, SHA-256/tamanho/MIME/tipo e autorização/licença verificáveis. Respeitar os pins existentes; não aceitar homônimo, placeholder, alias ou substituição de fonte/imagem para satisfazer o caminho.
3. Regenerar pelo fluxo do importador, verificar determinismo/`--check`, novo manifesto/skin, bytes e contagens, e executar/registrar os testes pertinentes quando permitido. A composição só retorna quando todas as suas dependências resolvem, em todos os estados; recursos restaurados não homologam delegates/licenças automaticamente.
4. Revalidar disponibilidade, inserção/seletor, cena/prévia/PPT e round-trip de novos e antigos perfis/documentos; então identificar novo build/pacote. Preservar revisões anteriores e gates; homologação Windows/DPI/Office e aceite exigem evidência própria, sem fechar T120–T122 ou as caixas anteriores por documentação.

**Pendências de homologação:** execução Go quando permitida; roteiros de inserção/seletor, salvar/reabrir e preservação; UI nativa/DPI/Office/aceite. Os testes Go novos cobrem consulta/Capture por código compilado, não novo ciclo de disco executado. Descoberta/19 Python, gofmt/vet/compile-only, build/pacote/integridade estão concluídos com o alcance acima, sem fechar as tarefas mistas. Nenhuma reimportação de recurso ausente nem execução nativa foi realizada; README da raiz, arquivos protegidos de camadas/alinhamento/divisor e constituição preservados.

<a id="15-dependencias-ausentes-removidas"></a>

## 15. Dependências ausentes: referências e ações removidas

**Escopo prevalente:** [FR-059](spec.md#dependencias-ausentes-fr-059)/[T133 aberta](tasks.md#t133-dependencias-ausentes). O usuário pediu apagar referências e ações sem implementação, deixando eventual reintegração para uma change. Foram removidos o produtor/chamadas dos delegates Java ausentes, avisos de classes externas para controles locais, botão/callback/adaptador Java/EFL, ferramentas Java 6 e opção/helper/ações da calculadora sem layout. Não é ocultação de mensagens. Os controles Go implementados e seus recursos permanecem.

**Evidência de origem:** o JSONL fornecido pelo usuário identifica o EXE anterior `87fa3e1876d2b4efe0d0c3ba71c254d6614e80f944f8dfffe2961ab635eaa0cb`, sessão `bd6e7a80ff9d411e67b1feae32dafbe1`, snapshot `2026-10-22T01:01:06.9338116Z`. Os seis recursos consultados resolveram com hash/tamanho exatos; `DELEGATE_MISSING` vinha das referências fixas, não de PNGs ou fonte faltantes. Escala antiga, sentinela da fonte, licença/embedding e portabilidade são condições distintas e continuam verdadeiras. Não há execução nativa desta nova revisão como evidência.

### Validação técnica vigente

- **27 Python executados PASS:** 19 importador, quatro empacotamento, quatro contratos estruturais de remoção, com os nomes dos quatro novos casos presentes na descoberta. Testes estruturais não executam o comportamento Go.
- Importador `--check` PASS: mesma skin `legacy-skin-v1-8c7d5015ddc5d9b4b8903730576d8b03e2351f9ab86bb5eace97b2eb57120b0a`, 206 definições, oito adiadas, oito fontes, 270 recursos e 31 diagnósticos. Manifesto inalterado.
- gofmt e `go vet ./...` PASS; `go test -c -o ./build/compile-local-only/ ./...` PASS com **14 suítes, ZERO execuções**. A redução de 15 para 14 decorre da exclusão do pacote Java, não de testes executados ou desativados para contornar a política.
- Cinco novas regressões Go cobrem controles locais sem delegates, geometria/estados do teclado, cabeçalho sem Java, seletor sem calculadora com preservação legada e guarda de ações antigas. Regressões de cena/PPT/geometria foram adaptadas; execução comportamental permanece pendente.

### Entrega e integridade

Tarefa VS Code Windows e `python -B -X utf8 scripts/build.py --target windows --arch amd64` **PASS**, sem `--check`. A conferência somente leitura confirmou PE AMD64, igualdade das duas cópias do EXE, ausência dos nomes das seis classes de pintura externas, botão/import Java e helper de calculadora, e presença da guarda `KEYBOARD_VARIANT_UNAVAILABLE`. Todos os **264 blobs originais** estão íntegros e incorporados no EXE. Não foi iniciado aplicativo, teste Go ou outro processo para contornar a política; nenhum processo do usuário foi encerrado.

| Artefato vigente | Bytes | SHA-256 |
| --- | ---: | --- |
| [EXE Windows x64](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) (canônico; identificador `sispl-efl-uc-…`, botão Nova tela e PPT sobre `modelo_jornada.pptx` com o identificador na capa) | 35.974.656 | `65c9c9d18d4ef39de8617ad48e0855a5f0878d52118a4f6c8bccbfa30f9416c0` |
| [Cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) (idêntica à canônica neste build) | 35.974.656 | `65c9c9d18d4ef39de8617ad48e0855a5f0878d52118a4f6c8bccbfa30f9416c0` |
| ZIP Windows x64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 21.419.209 | `0a590b1e7e101006fb7438a6b8fe2cb2c572f582b6f9391ce256b1e6c1c21f35` |
| Manifesto original incorporado (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/internal/adapter/legacyskin/assets/manifest.json`) | 1.039.093 | `60713ea85e1b27c38fd8ef37d9c4c74a9c1bbe074bc0269a00e2da8b84a6b3fb` |
| [README da aplicação congelado](../../apps/gerador-jornada/README.md) | 50.724 | `a76c566c79b594f1455eed8cf20dadc64254680a2e20343089d2efb7370aa391` |

**Observações em Recursos / referência (relato com o EXE `12dee2c5…`, sessão `7bc85d45…`):** o log mostrou seis recursos resolvidos por hash/tamanho exatos e nenhum `DELEGATE_MISSING`; as linhas restantes eram `KEYBOARD_LOCAL_SCALE`, `FONT_CMAP_SENTINEL_NONCANONICAL`, `FONT_EMBEDDING_RESTRICTED`, `FONT_LICENSE_UNKNOWN` e `FIDELITY_UNVERIFIED`, condições informativas e não falhas de carregamento. O campo passou a apresentá-las em bloco “Observações (não são erros…)”, em texto simples e sem código; sem falhas reais exibe “Nenhum erro”. Códigos continuam no Log de aplicação, nas APIs de perfil e nas travas de exportação/empacotamento (`redistributionApproval: pending`); nada foi aprovado nem redimensionado. A cópia de compatibilidade estava bloqueada pelo EXE em execução e não foi atualizada; use o EXE canônico. Go não foi executado (vet e compilação das 14 suítes apenas).

**Integridade do pacote PASS:** CRC ZIP, inventário exato de seis arquivos sem backups `~`, todos os cinco hashes de `SHA256SUMS`, sidecar do ZIP, manifesto byte a byte, ícone e README exatos com o prefixo/ajuste de imagem e normalização CRLF previstos pelo empacotador. Metadados registram Windows/amd64, produto 1.0.0.0, Go 1.27.2 e GCC 16.2.0; `manualApproval` e `redistributionApproval` continuam `pending`. A skin clara original conserva o pin `475dc455692a317373efe150080fcded40e61a176c0d4fee09420e23126e1c54`. O manifesto e seus recursos não mudaram nesta retirada de código. Os hashes da seção 14 identificam exclusivamente o EXE anterior; a integridade não constitui homologação visual, funcional, legal ou Office.

### Roteiro manual pendente

1. Usar o EXE identificado nesta seção; conferir no cabeçalho a ausência de Java/EFL e a presença de Novo/Abrir/Salvar/Prévia/Prévia Tela Atual/PowerPoint.
2. Inserir e usar teclados numeric/alpha, campos, senha, diálogo e demais famílias locais. Consultar o log/Recursos e confirmar ausência de dependência de classes Java; avisos reais não devem desaparecer por filtro.
3. Abrir arquivo antigo com `calculator`: leitura e salvamento neutro mantêm o valor; não aparecem teclas, hit areas ou ações de outro layout. Escolher explicitamente numeric/alpha/Padrão para substituí-lo; salvar/reabrir e desfazer/refazer sem alterar outros dados.
4. Conferir dimensões antigas, IDs/FileIDs, ordem, vínculos, estilos e perfis/revisões, incluindo seleção da skin clara original. Não aceitar resize, upgrade, fallback ou poda de documento.
5. Executar testes Go quando permitido e revisar Windows/DPI/Office e aceite. Build e verificações estruturais não homologam fidelidade, licença, embedding ou portabilidade.


<a id="16-remocao-do-campo-recursos--referencia"></a>

## 16. Remoção do campo Recursos / referência — 2026-10-23

**Escopo:** a pedido do usuário (o campo “nunca funcionou” e não é necessário), foram removidos o item **Recursos / referência** do inspector, o botão **Diagnóstico de recursos / referência…** da prévia e o registro automático do inspector no log. Mantidos: seletores Skin/Face da fonte/Versão (rascunho sem publicar), Salvar/Descartar, **Log de aplicação** (captura manual do componente selecionado ou da tela atual, estado salvo), `componentDiagnostics*`, travas de exportação/empacotamento e aprovações `pending`. Sem alteração de schema/catálogo/domínio/perfil. Tarefa: [T134](tasks.md#t134-remocao-recursos-referencia). As seções 10, 12 e 13 descrevem o comportamento anterior e são históricas.

### Validação técnica realizada

- `go vet ./...` PASS; `gofmt -l internal cmd` sem arquivos listados.
- `go test -c -o ./build/compile-tests-rr/ ./...` PASS: suítes **somente compiladas, ZERO execuções Go** (política Windows, sem contorno).
- Python `unittest discover` em `scripts/test_*.py`: **OK** (27 testes).
- Tarefa VS Code de compilação Windows e `python -B -X utf8 scripts/build.py --target windows --arch amd64` (sem `--check`) PASS. O ZIP foi reconstruído depois da edição do README; o README embutido contém a nota de remoção (53.530 bytes após prefixo/normalização CRLF do empacotador, contra 53.028 no disco).
- Não foi iniciado aplicativo; a UI nativa não foi observada.

### Artefatos

| Artefato | Bytes | SHA-256 |
| --- | ---: | --- |
| [EXE Windows x64](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) (canônico) | 35.964.416 | `4e0526c4abc0e9559383093b11bdc2cd00619df562929e39c3174b4592f5ac7a` |
| [Cópia de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) (idêntica) | 35.964.416 | `4e0526c4abc0e9559383093b11bdc2cd00619df562929e39c3174b4592f5ac7a` |
| ZIP Windows x64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 22.443.500 | `38d86c337f305ac67e267616ab7af21d9b6612e87858fc6e9069780aead781bc` |

Os hashes das seções anteriores identificam apenas suas revisões. Fechar instâncias abertas antes de testar o novo EXE.

### Roteiro manual pendente

1. Abrir o EXE acima, selecionar componentes de várias famílias: **Propriedades não mostra Recursos / referência**; Skin, Face da fonte e Versão do perfil continuam presentes.
2. Trocar Skin/Face/Versão: o rascunho muda sem publicar; **Salvar alterações** aplica e **Descartar** reverte; desfazer/refazer intactos.
3. Abrir a prévia: sem botão de diagnóstico de recursos; a prévia continua funcionando.
4. Cabeçalho → **Log de aplicação**: o popup abre; **Registrar diagnóstico agora** grava (com e sem seleção), **Copiar caminho** e **Abrir pasta** funcionam.
5. Salvar/reabrir projeto existente: sem mudança de dados. Execução Go, DPI/Office e aceite seguem pendentes.

<a id="17-codigos-barras-2026-10-23"></a>
## 17. Códigos de barras — especificação por domínio (2026-10-23)

**Registro da etapa documental original de 031/032:** requisitos/planos/modelos/contratos e tarefas acrescentados para 031/032 e domínios 002–006 sob 001/FR-060. Naquela etapa, implementação não iniciada e sem build novo utilizável desses componentes, testes funcionais, manual, PowerPoint ou aceite. Os EXEs/hashes anteriores àquela etapa não são sua evidência. **Atualização de autorização**: usuário autorizou posteriormente implementar 031/032/033 e testar manualmente depois; estado/evidências do novo lote devem ser registrados na seção18, sem converter este registro histórico em resultado atual de implementação.

### Pré-requisitos para uma revisão executável

Concluir tarefas técnicas dos [domínios e integração](tasks.md#tarefas-codigos-barras-2026-10-23), verificar compatibilidade planejada 3/3 e identificar o novo build. Receber as imagens citadas para revisar fidelidade e as classes/vetores de módulo para comparar equivalência SISPL, sobretudo SIGCB. Enquanto ausentes, registrar essas limitações, não usar o adjetivo homologado. Não contornar política Windows que bloqueie a execução de testes Go.

### Roteiro manual futuro — ainda não executado

1. Em Windows 10/11 x64, registrar revisão/hash/data do build novo, executor, resolução e DPI. Inserir duas entradas e um processado. Entradas começam com oito campos vazios; processado com um bloco/CNPJ/Nome vazios.
2. Digitar F01 do [contrato](contracts/component-catalog.md#contrato-codigos-barras-2026-10-23) pelos oito campos. Conferir avanço em pares, campo de DV geral e aceitação final. Adulterar cada DV local e conferir bloqueio/foco no par; adulterar somente DV geral e conferir bloqueio final/retorno ao primeiro campo. Cobrir conclusão explícita do campo final incompleto com zeros à direita.
3. Submeter F01 de 44 posições pelo modo de leitura, depois por colagem integral de 44 e de 47 posições; conferir que os três percursos produzem exatamente os mesmos oito valores. Submeter comprimento inválido/caractere/DV errado durante entrada parcial: limpar somente a instância destinatária, abrir “Código de barras inválido” / “OK, entendi”, Enter devolve foco ao primeiro campo. Cancelar captura não submetida preserva rascunho.
4. Usar teclado numérico vinculado, TAB, apagar, corrigir e reiniciar; conferir que não atravessa bloco inválido, desabilitado não recebe entrada e uma instância não altera outra. Revisitar A→B→A deve restaurar entrada vazia sem mudar rotas/configuração. Prévia Tela Atual continua sem navegar.
5. Na ficha processada, adicionar segundo bloco e configurar linhas com CPF/CNPJ e Nome/Nome fantasia/Razão social. Alterar cada seletor/valor isoladamente; confirmar/cancelar/reordenar/remover/desfazer/refazer e comparar os demais itens. Valor vazio e Unicode/pontuação permanecem literais; remover último bloco/linha é recusado sem perda.
6. Duplicar componente e tela, alterar a cópia e verificar independência/remapeio de IDs/teclado. Salvar/reabrir, inclusive original histórico e perfil próprio/arquivado: preservar identidades/ordem/valores/bounds, sem enriquecer perfil, gravar runtime ou perder original na migração. Verificar aviso/cópia íntegra do original no primeiro save incompatível.
7. Exportar exploração com os compostos e inspecionar no PowerPoint real textos/formas editáveis e ordem/Unicode/geométrica. Não há teclado, DV, combo de autoria ou processamento executável no slide. Roteiro selecionado tentando typeText nos compostos deve receber diagnóstico localizado, sem omissão silenciosa.
8. Comparar estados/legibilidade/overflow em DPI 100%, 150% e 200% e contra os anexos quando recebidos; registrar diferenças. F01 deve fornecer banco001/R$1,00/fator3737; F02 é somente teste da premissa experimental, não homologação SISPL ou pagamento.

### Registro de execução e aprovação

Executor, revisão executável, ambiente, resultados observados e evidências: **pendentes**. T138 deve registrar verificações técnicas realmente executadas, distinguindo compile-only de testes; T139 registra manual, reteste e aceite explícito. Checklists de qualidade de 031/032 não encerram essas tarefas. Não copiar resultado esperado acima como resultado observado.

### Verificação documental observada nesta etapa

Agente, workspace Windows, 2026-10-23: pré-requisitos oficiais `check_prerequisites.py --require-spec --require-tasks --include-tasks` aprovados para 001–006, com seleção temporária por ambiente e sem gravar o ponteiro ativo (031). Verificados 42 documentos alterados/novos, 122 links locais/âncoras, IDs sem duplicação e diff sem erro de whitespace. As 14 tarefas de código novas permanecem abertas. Conferência aritmética de F01/F02 e limites documentados aprovada em script isolado; não é teste do produto ou comparação com classes bancárias reais.

Revisão independente documental apontou dois ajustes de clareza (aceite de colagem44 e dependências da oferta); ambos incorporados. Checklists embutidos de qualidade 031/032 revisados em 16/16 itens cada. Não executados testes Go, aplicativo, build novo, interação nativa, DPI, PowerPoint ou homologação. Não se atribui a este incremento evidência de executáveis anteriores. Uma remoção externa de arquivo temporário Office detectada no workspace foi preservada, sem operação do agente sobre ela.

<a id="18-menu-aplicacao-2026-10-23"></a>
## 18. Menu aplicação e lote autorizado 031/032/033 — 2026-10-23

**Estado: implementado e compilado para teste manual; homologação pendente.** Usuário autorizou implementar 031/032/033 com agentes menores por domínio e testar depois, mesmo com pendências históricas nos checklists. Modelos, regras, codec/backup, formulários, sessão local, cena/PPT, índice de navegação e quadros foram integrados pelo executor principal. Evidências abaixo são desta revisão, não herdadas das seções anteriores. Não houve execução do aplicativo, de testes Go ou de PowerPoint nesta etapa.

### Pré-requisitos e gates

Concluir trilhas técnicas dos [domínios e integração](tasks.md#tarefas-menu-aplicacao-2026-10-23), identificar revisão executável utilizável e registrar verificações permitidas. Catálogo3 é alvo comum aos três tipos, schema3 preservado; leitura1/1 2/1 3/2 e3/3, backup antes de save incompatível. Não contornar bloqueio de testes Go. Checklists built-in históricos de001–006 têm itens abertos de revisão/aceitação; autorização atual permite implementação, **não os aprova**. Checklists customizados/revisores não foram alterados. Aceite exige manifestação explícita sobre a revisão testada, conforme constituição V/VII.

O seed e as demais medidas devem ser conferidos com o Java/XML fornecido na conversa; não há imagem recebida para comparação visual. Caso a transcrição não esteja disponível ao executor, pedir sua reposição antes de alegar seed fiel. Receber módulos/vetores independentes de031 para declarar equivalência SISPL continua gate próprio de T139. Clique demonstrativo por item permanece decisão principal; baseline click visual/estado estático e ausência de typeText estão contratados.

### Roteiro manual futuro do menu — resultados ainda não observados

1. Registrar Windows10/11x64, resolução, DPI, revisão/hash/data do novo build e executor. Em Itens legados localizar as três opções; inserir duas instâncias de Menu aplicação e conferir independência/defaults habilitados, serviços à esquerda e jogos à direita, conforme fonte original.
2. No inspector, alterar legenda/categoria/variante/enabled, adicionar/remover/reordenar itens/submenus aninhados, escolher uma folha com destino explícito e outra sem destino. Salvar em uma transação, desfazer/refazer; outro rascunho descartado não deve alterar conteúdo. Testar Unicode/rótulos repetidos, menu vazio e remoção do submenu inicial exigindo substituto válido.
3. Exercitar root→standard→financial e ramo Bolão digital; conferir grade financeira140x80 e três botões/legendas Bolão da referência. Voltar retorna um nível e na raiz não muda nada; Cancelar vai à raiz mesmo quando inicial é outro submenu. Desabilitar ramo/instância impede ativação dos seus itens.
4. Selecionar folha sem destino: somente informação de seleção simulada, nenhum movimento/conta/QR/rede. Folha com destino abre apenas tela autorada na prévia completa. Na Prévia Tela Atual, submenus operam mas destinos entre telas não são seguidos; destinos externos não impedem exploração local. Sair fecha somente prévia e mantém editor/projeto/rascunhos intactos.
5. Escolher submenu inicial aninhado, confirmar e reabrir prévia: inicia nele, Voltar encontra ancestral. Explorar outro ramo e fechar/revisitar: restaura inicial confirmado, não último caminho. Comparar duas instâncias sem estado compartilhado.
6. Duplicar componente e tela; IDs internos novos, escolha inicial remapeada, cópia editável independentemente. Destinos externos conservados e autorreferentes remapeados somente ao duplicar tela. Salvar/reabrir árvores/false/strings vazias/Unicode/ordem/estilos/metadados e conferir ausência de caminho/foco/seleção simulada no manifesto.
7. Abrir cópia de projeto histórico, salvar incompatível com aviso/backup íntegro, sem redimensionar/enriquecer perfil/alterar original ao abrir. Conferir falha de backup sem overwrite; estrutura inválida/limites17níveis ou1001itens/UUIDs repetidos/submenu+destino não publicam estado parcial. Fixtures16/1000 devem ser preservadas na validação técnica, sem necessidade de criar1000itens manualmente.
8. Exportar menu inicial escolhido e barcodes conforme seção17; abrir no PowerPoint real e conferir texto/formas editáveis, ordem/legendas/estilo e estado fixo independente da exploração livre. Roteiro com typeText dos compostos é recusado; click visual não opera árvore silenciosamente. Não prometer interação ou dispatch no slide.
9. Conferir hit regions, overflow, títulos longos e legibilidade nos quatro layouts, DPI100/150/200%, canvas/prévia/PPT. Anexos vazios impedem afirmar comparação visual; registrar divergências e imagens quando recebidas. Resultado de referência Java não é execução/homologação EFL.

### Registro de execução e aprovação do lote

| Campo | Estado desta entrega para teste manual |
| --- | --- |
| Revisão executável / hash / ambiente | Windows amd64, Go 1.27.2; EXE final `b313093dadbcef6c1f2a002be8f2d6999ce355377f3d8bedb24b6967ac62aef3`; executor agente, 2026-10-23 |
| Checks técnicos e testes efetivamente executados | Formatação e vet aprovados; 14 suítes Go compiladas, ZERO execuções; 27 testes Python executados e aprovados; importador em modo check aprovado |
| Teste manual do usuário / resultado observado / evidências | Pendente; roteiro acima é expectativa, não observação |
| Correções / reteste / aceite explícito | Pendentes; T139 para barcode e T144 para menu/lote não são fechadas por autorização |
| Fidelidade visual / SISPL / Office / licenças | Gates não resolvidos por documentação |

Nova trilha técnica: 001/T140–T143, 002/T015, 003/T021, 004/T017, 005/T012, 006/T012; manual001/T144. Estados de implementação atualizados nos backlogs por domínio; verificações mistas e manual continuam abertos. Checklists de revisores não foram alterados. Esta seção é o local canônico das evidências.

### Revisão integrada disponível e verificações observadas

| Artefato final | Bytes | SHA-256 |
| --- | --- | --- |
| [EXE canônico Windows](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) | 36.359.680 | `b313093dadbcef6c1f2a002be8f2d6999ce355377f3d8bedb24b6967ac62aef3` |
| [EXE de compatibilidade](../../apps/gerador-jornada/dist/JornadaStudio.exe) | 36.359.680 | `b313093dadbcef6c1f2a002be8f2d6999ce355377f3d8bedb24b6967ac62aef3` |
| ZIP Windows amd64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 21.549.823 | `89b591cf898f6c60b975db6eecd9ef8ec4677f0ac6f968a48f0124bcb0ad2b2f` |

Empacotamento final sem executar checks Go bloqueados. Conferência independente aprovou CRC das seis entradas do ZIP, cinco hashes do inventário, sidecar do ZIP, cabeçalho PE amd64, igualdade dos dois EXEs e do EXE compactado, README definitivo com nota/caminho de ícone transformados pelo empacotador e manifesto visual byte a byte. Aprovação manual e autorização de redistribuição continuam `pending`; não se publicou release remota nem se iniciou a janela nativa.

- Abrir o [executável Windows atualizado](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), não uma instância antiga já aberta. A tarefa oficial Windows concluiu a compilação final com o hash acima.
- Em **Paleta → ELEMENTOS DO LEGADO**, inserir **Código de barras**, **Código de barras processado** ou **Menu da aplicação**. Confirmar/descartar pelo rodapé de Propriedades; menus/submenus são uma instância composta, não novas telas.
- Entrada: **Ler · 44 + Enter**, **Cancelar leitura**, **Colar 44/47**, **Concluir · Enter** e **Reiniciar** ficam na prévia. A leitura é explícita, não depende de velocidade ou driver; valores capturados não são persistidos. Validação é matemática local, não pagamento nem homologação bancária.
- Processado: editar blocos, títulos, linhas, rótulos e valores ilustrativos independentes. A entrada de boleto não preenche essa ficha automaticamente.
- Menu: editar **Submenu em edição**, **Item selecionado**, **Legenda explicativa** e **Menu inicial**. Inicial sempre referencia UUID existente; categoria/legendas conservam texto literal. Layout **bolao** é o digital de três botões/legendas e não depende do nome editável do nó. Usar **standard** para a grade do Bolão ordinário.
- Prévia: Voltar remove um nível, Cancelar retorna à raiz e Sair fecha somente a prévia. Destinos explícitos de folhas navegam somente no modo completo. No roteiro, click é visual, sem dispatch de item; typeText nos três compostos é recusado com ID da etapa.
- Compatibilidade: writer schema3/catalog3; readers 1/1, 2/1, 3/2 e 3/3. O primeiro save incompatível também avisa e preserva backup integral de 3/2; abrir não regrava o original. Após save bem-sucedido a origem é atualizada inclusive no histórico.
- Código final congelado: `go vet ./...` aprovado e `go test -c -o ./build/compile-tests-compounds-final/ ./...` compilou 14 suítes. **Não executadas** por política Windows; não foi usado contorno. Os casos de round-trip, prévia, UI e PPT presentes nas suítes ainda não são resultados observados.
- Python: descoberta `test_*.py` executou **27 testes** (19 importador, 4 distribuição, 4 remoção) com resultado OK. Importador `--check` confirmou 206 definições ativas, 8 adiadas, 8 fontes, 270 recursos e 31 diagnósticos, sem regenerar o acervo.
- Spec Kit: pré-requisitos oficiais aprovados para os seis domínios com seleção temporária e sem mudar o ponteiro ativo. Nenhum hook before/after registrado, pois configuração de extensões não existe.
- Pendentes: roteiro manual das seções17–18, recorte/overflow de árvores extensas, DPI 100/150/200%, Office real, imagens originais e equivalência SISPL/SIGCB. Geometrias/legibilidade ainda não homologadas; nenhum resultado esperado foi convertido em observado. Remoções externas de temporários Office ou apresentações pessoais foram preservadas.

Revisão documental final indexou 51 documentos no escopo, 226 tarefas em sete backlogs e 72 âncoras explícitas sem IDs/âncoras duplicados por arquivo. Não encontrou novos links quebrados nas referências conferidas; identificou a omissão de `legend` na tabela wire, corrigida nesta integração. Essa revisão é estática, não execução de fixtures nem homologação. Os quantitativos 42/122 da seção17 continuam históricos da etapa documental original.

<a id="19-menu-ordem-xml"></a>
## 19. Correção do menu após teste parcial — ordem, recorte e XML

Esta revisão substitui os artefatos da seção 18, cujos hashes ficam históricos. Implementação separada em **003/T022** (fábrica/árvore) e **001/T145** (cena/regressões); **T146 permanece aberta** para verificação comportamental/manual. Checklists e aceite anteriores não foram promovidos.

- Raiz: 12 serviços e 10 jogos. Na esquerda, **APOIO → RELATÓRIOS → MEIOS DE PAGAMENTO**; na direita dos serviços, **BOLÃO → PAGA PREMIADOS → OFERTAR**. Bolões Digitais passa ao último item de Bolão, como no XML.
- Posições e tamanhos conservados: serviços 164 × 68, colunas x20/x200 e passo vertical 73. O recorte esquerdo agora chega a y600, mostrando integralmente a última linha y515–583; Sair permanece à direita, sem sobreposição.
- Fábrica corrigida para Paga Premiados, Meios de Pagamento e Relatórios. Fixtures cobrem 148 folhas/27 títulos de Relatórios, códigos literais como `080206`, espaços e sequências distintas das opções diárias. Regressão de apresentação confere as seis posições finais e regiões clicáveis completas; fixtures foram compiladas, não executadas.
- **Inserir um novo Menu da aplicação para receber a árvore corrigida.** Instâncias autoradas/salvas não são reescritas automaticamente. O ajuste de recorte também beneficia instâncias existentes, mas não modifica sua ordem ou conteúdo persistido.

### Checks observados nesta correção

`gofmt` e `go vet ./...` concluídos; **14 suítes Go compiladas somente**, em build/compile-tests-menu-order, sem nenhuma execução Go ou início automático do aplicativo. **27 testes Python executados e aprovados**. Tarefa Windows oficial e empacotamento sem `--check` concluídos.

Conferidos independentemente: CRC das seis entradas do ZIP, cinco hashes do inventário, sidecar, PE amd64, igualdade dos três conteúdos EXE (canônico/cópia/ZIP), README com transformações do empacotador e manifesto visual byte a byte; nenhum backup `~` incluído. `manualApproval` e `redistributionApproval` continuam `pending`.

| Artefato desta correção | Bytes | SHA-256 |
| --- | --- | --- |
| [EXE canônico](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe) e [cópia](../../apps/gerador-jornada/dist/JornadaStudio.exe) | 36.369.920 | `8c2edeab9132c0fbb938145d789cd3bda43c3c726a78c1c02265b4d10c8f39fe` |
| ZIP Windows amd64 (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 21.554.628 | `bfaad837eaaa6eed6bd6eee9feac76c4b10e94dab71a568bdd879eca87997bd3` |

### Reteste manual pendente

Fechar a instância antiga e abrir o EXE identificado acima. Inserir novo menu; conferir as duas colunas e a sexta linha inteira no canvas e na prévia. Abrir Paga Premiados (três opções), Meios de Pagamento/PIX e todos os níveis de Relatórios; conferir Voltar/Cancelar/Sair e Bolão/Bolões Digitais. Repetir com escala/DPI 100/150/200%, exportação PPT editável e salvar/reabrir uma instância editada, sem perda de autoria. Registrar ambiente, perfil/fontes, resultado observado e divergências. A compilação não substitui esse reteste nem o aceite explícito.

## 20. Faixa somente nos submenus — 2026-10-23

Correção adicional de 001/T145/T146 e 033/DESIGN: remover faixa/título da raiz e usar o token branco `label_faixa_titulo_principal` nos submenus. Regra aplicada à cena compartilhada, sem alterar os dados ou as posições dos botões. Também vale para menus existentes; não exige reinserção para esta correção visual. Regressão verifica raiz sem faixa e submenu com título branco, inclusive navegação por identidade.

Observados neste incremento: gofmt e vet integral concluídos; 14 suítes Go somente compiladas, sem execução por política Windows. Build oficial e empacotamento concluídos. CRC, hashes do inventário, sidecar e igualdade canônico/EXE compactado conferidos. Não foram reexecutados os testes Python da seção 19 nem iniciado o aplicativo.

- [EXE canônico atualizado](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe): 36.369.920 bytes, SHA-256 `ca6f3785885c1b19902903f5c684e692300ee5ed1b3abfcfb6052c3e72c03cb3`.
- ZIP atualizado (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`): 21.554.693 bytes, SHA-256 `3d04b1276f538b7f8d8854fb249166086028a37d67f93ea13fbaedca5ac13bc6`.
- **Cópia de compatibilidade antiga bloqueada por uso:** o empacotador não conseguiu atualizar o EXE diretamente em dist. Seu hash continua `8c2edeab9132c0fbb938145d789cd3bda43c3c726a78c1c02265b4d10c8f39fe`. A primeira conferência que exigia igualdade com essa cópia falhou; a conferência independente posterior aprovou somente canônico e ZIP. Nenhum processo foi encerrado nem política contornada.

Teste visual/Office/DPI e aceite continuam pendentes em T146. Usuário relatou também um botão incorreto, mas os anexos desta mensagem vieram vazios: necessário reenviar captura ou identificar botão para investigar sem inventar a causa. Não declarar esse defeito corrigido. Para reteste da faixa, fechar instância antiga e abrir o canônico acima; conferir raiz sem faixa e todos os submenus com título branco.

<a id="21-teste-humano-001-033-2026-10-24"></a>
## 21. Teste humano consolidado — features 001–033 (2026-10-24)

Esta seção consolida em uma única matriz o teste humano das features 001–033. A matriz é **EXPECTATIVA**: descreve o que deve ser observado, não o que foi observado. O status inicial de **todos** os cenários é **NÃO EXECUTADO**, e nenhum resultado esperado desta seção pode ser lido como observado, aprovado, homologado ou aceito. Ela não substitui as evidências históricas das seções [17](#17-codigos-barras-2026-10-23) e [18](#18-menu-aplicacao-2026-10-23) nem as correções das seções [19](#19-menu-ordem-xml) e 20.

### 21.1 Como executar

**Revisão congelada.** Antes de abrir qualquer executável, registrar a saída dos dois comandos no workspace:

```text
git rev-parse HEAD      -> 13e09d531f9c02b5a549b3d79c000f458bcbf1dd
git status --porcelain   -> 262 entradas na consulta final (213 modificadas + 49 não rastreadas)
```

O HEAD desta revisão é o commit `13e09d531f9c02b5a549b3d79c000f458bcbf1dd` («Versão estavel», 2026-10-19 09:23:50 -0300). Como `git status --porcelain` mudou durante a sessão (296–299 entradas antes das implementações e 262 ao final, com 213 arquivos modificados e 49 não rastreados), a árvore de trabalho **não** é idêntica ao HEAD: o pacote não é reproduzível apenas a partir do commit. Congelar, portanto, a dupla **HEAD + árvore de trabalho**, anexando a listagem completa de `git status --porcelain` (ou `git diff` e um `git stash create`) ao registro de execução. A contagem deve ser relida no momento do teste; revisão diferente desta exige novo registro e não herda cenários.

**Artefato executável.** A revisão de teste desta seção foi compilada em 2026-10-24 e o EXE canônico está presente em `apps/gerador-jornada/dist/windows-amd64/`; a cópia de compatibilidade também foi atualizada. O pacote ZIP contém a mesma árvore.

| Artefato | Caminho | Bytes | SHA-256 | Situação nesta revisão |
| --- | --- | --- | --- | --- |
| EXE Windows amd64 | `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe` | 32.754.176 | `953f9e102ff1cbbf8d007cf5fbebcfcb86b2a458c2533ff59f90401a710605f1` | presente; compilado nesta revisão |
| Cópia de compatibilidade | `apps/gerador-jornada/dist/JornadaStudio.exe` | 32.754.176 | `953f9e102ff1cbbf8d007cf5fbebcfcb86b2a458c2533ff59f90401a710605f1` | presente e idêntica ao canônico |
| Pacote Windows amd64 | ZIP (referência histórica/proposta sem arquivo nesta árvore: `../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.0-windows-amd64.zip`) | 19.811.611 | `a74cbb5115a253351315cf29a2ba740ce960917ddabf18134e334e4aea20f3d8` | presente |
| Inventário do pacote | [SHA256SUMS.txt](../../apps/gerador-jornada/dist/windows-amd64/SHA256SUMS.txt) | 330 | `d51bc8d6729d5f0dcf45922626ad90e2b2a19dfb11195266f0035e7d274926e0` | presente; registra o EXE acima |
| Metadados do build | [build-info.json](../../apps/gerador-jornada/dist/windows-amd64/build-info.json) | 401 | `67a9d4390f26a1c90661538075e0f7b835853ec0c931cbc45fa473d80ab95a7c` | `go1.27.2 windows/amd64`, `signing: unsigned`, `manualApproval: pending` |

Os hashes `0d512ae9…` (EXE do ZIP anterior) e `bbddb406…` (ZIP anterior) permanecem **históricos** e não identificam esta revisão. Da mesma forma, `b313093d…`/`8c2edeab…`/`ca6f3785…` e os ZIPs `bfaad837…`/`3d04b127…` das seções 18–20 continuam históricos.

**Ambiente e executor (preencher no teste).**

| Campo | Valor exigido pela matriz |
| --- | --- |
| SO | Windows 10/11 x64; Linux não cobre os cenários de janela nativa |
| Resolução | 1366 × 768 (mínima da matriz) |
| DPI / escala | 100%, 150% e 200%, em três passagens independentes |
| PowerPoint | versão real obtida em Ajuda → Sobre; obrigatório nos cenários de PPT |
| Executor | nome de quem executa o teste humano |
| Data e hora | início e fim, com fuso |
| Perfil / fontes | perfil próprio ou arquivado usado e diagnóstico de fonte ausente |
| Build aberto | revisão congelada acima e SHA-256 conferido do EXE executado |

**Preparação.**

1. Registrar a revisão congelada conforme acima e anexar a listagem ao registro.
2. Abrir o EXE canônico em `apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe` ou a cópia idêntica em `apps/gerador-jornada/dist/JornadaStudio.exe`; não usar o ZIP anterior como substituto do build novo.
3. Conferir o SHA-256 do arquivo que será aberto com `Get-FileHash <caminho> -Algorithm SHA256`. Divergência de `953f9e10…` significa **outra revisão**: registrar e não preencher a matriz como se fosse esta.
4. Encerrar instâncias antigas do JornadaStudio antes de abrir a nova, para não atribuir à revisão atual uma janela já aberta.
5. Percorrer a matriz na ordem 001 → 033, uma passagem por DPI, registrando cada cenário imediatamente após executá-lo.
6. Não usar `typeText` nos compostos 031/032/033: o roteiro deve receber diagnóstico localizado, nunca simulação silenciosa.

**Regra de leitura da matriz.** Status inicial de tudo: **NÃO EXECUTADO**. Só existe resultado quando houver observação e evidência anexada; ausência de evidência mantém o cenário em NÃO EXECUTADO. Um roteiro descrito, um requisito documentado, um teste compilado ou uma revisão de código não são execução. Onde o componente exigir anexo de referência ausente, o resultado máximo admitido é NÃO VERIFICÁVEL, nunca OBSERVADO-CONFORME.

### 21.2 Matriz-resumo por feature (001–033)

Cada linha é um cenário de aceitação derivado de `specs/00X/spec.md`; para 031–033 o detalhamento integral está em 21.3, 21.4 e 21.5. A coluna Evidência indica o arquivo a produzir, não um arquivo existente.

| ID | Feature | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- |
| H001-01 | [001](spec.md) US1 | Criar tela; arrastar componentes; mover, redimensionar, alinhar, duplicar; desfazer/refazer; salvar e reabrir | Área inicial 800 × 600 com nome editável; inserção na posição indicada; cada operação desfeita/refeita como ação única; reabertura preserva telas, componentes, recursos e referências, sem gravar seleção/zoom | `ev-H001-01_<rev>_dpi100.png` + `.txt` |
| H001-02 | 001 US2 | Definir tela inicial; botão de A com destino B; abrir prévia e clicar; tentar remover tela referenciada; criar destino inexistente | Prévia abre a inicial e o clique leva a B; remoção mostra referências afetadas e pode ser cancelada; destino inexistente gera diagnóstico com tela e componente de origem e o rascunho continua salvável | `ev-H001-02_*.png` + `.txt` |
| H001-03 | 001 US3 + [contrato PPT](contracts/pptx-export.md) | Exportar em máquina sem Office; abrir no PowerPoint real; editar textos e formas; clicar botões com destino em modo apresentação; percorrer ciclo entre telas | PPTX é gerado sem iniciar ou instalar PowerPoint; exploração parte da tela inicial sem avanço automático; textos e formas editáveis, sem imagem única da tela; botões abrem o slide correto inclusive em ciclo | `ev-H001-03_*.pptx` + PNG |
| H002-01 | [002](../002-projeto/spec.md) US1 | Editar projeto com telas e roteiros; salvar; consultar identidade, nome e revisão; reabrir | IDs permanecem estáveis e a revisão salva representa exatamente o snapshot gravado, sem alterar identidade ao renomear | `ev-H002-01_*.txt` |
| H002-02 | 002 US2 | Usar um recurso em duas telas; remover uma delas; depois remover o recurso ainda referenciado | Recurso ainda referenciado permanece disponível; recurso ausente ou alterado gera diagnóstico localizado | `ev-H002-02_*.png` |
| H002-03 | 002 FR-011 / 031–033 | Abrir projeto 1/1, 2/1, 3/2 e 3/3; salvar em 3/3; forçar primeiro save incompatível de cópia histórica | Leitores aceitam as quatro versões; primeiro save incompatível avisa e preserva backup íntegro; abrir não regrava o original; após save bem-sucedido a origem é atualizada, inclusive no histórico | `ev-H002-03_*.txt` + cópia do backup |
| H003-01 | [003](../003-catalogo/spec.md) US1 | Abrir a paleta e listar os componentes do catálogo versão 1 | Paleta recebe exatamente botão, rótulo, título, campo, seleção, relatório e teclado numérico, com identificadores estáveis | `ev-H003-01_*.png` |
| H003-02 | 003 US2 | Em um campo de texto, tentar atribuir propriedade de outra família e enum desconhecido | Operação recusada sem alterar o componente; propriedades suportadas preservadas | `ev-H003-02_*.png` |
| H003-03 | 003 FR-016 / 023 | Verificar o grupo **ELEMENTOS DO LEGADO** com Código de barras, Código de barras processado e Menu da aplicação | Oferta restrita aos tipos com consumidores completos; componentes primitivos permanecem separados dos conjuntos prontos | `ev-H003-03_*.png` |
| H004-01 | [004](../004-tela/spec.md) US1 | Inserir componente válido em tela 800 × 600 e movê-lo com zoom e DPI diferentes | Bounds armazenados em unidades lógicas, independentes de zoom ou DPI da janela | `ev-H004-01_*.txt` (bounds) |
| H004-02 | 004 US2 | Duplicar tela com teclado vinculado e botão que aponta para a própria tela | Tela e componentes recebem novos IDs e vínculos internos são remapeados para a cópia | `ev-H004-02_*.txt` |
| H005-01 | [005](../005-jornada/spec.md) US1 | Criar telas A e B; apontar ação navigate do botão de A para B; remover a ação de outro botão | Resolvedor retorna B pelo ID; ausência de ação equivale a none, sem inferência por texto ou posição | `ev-H005-01_*.txt` |
| H005-02 | 005 US2 | Remover o destino de uma ação e validar a jornada | Diagnóstico identifica origem e destino ausente; grafo inválido não é corrigido silenciosamente | `ev-H005-02_*.png` |
| H006-01 | [006](../006-demonstracao/spec.md) US1 | Interpretar roteiro A → B; executar etapa de clique e etapa de navigate | Clique produz destaque visual e navigate segue o binding explícito, sem disparo oculto de navegação | `ev-H006-01_*.png` |
| H006-02 | 006 US2 | Reproduzir roteiro A → B → A; pausar; retomar | Três ocorrências produzidas e a revisita a A começa com valores iniciais; pausa não consome tempo | `ev-H006-02_*.mp4` ou PNG por etapa |
| H007-01 | [007](../007-botao-efl/spec.md) US1 | Inserir botão por clique e por arraste; editar a legenda | Inserção nas duas formas; legenda editável sem mudar a identidade do botão | `ev-H007-01_*.png` |
| H007-02 | 007 US1 | Escolher variante confirmar, cancelar, sim ou não e vincular o clique a tela existente; abrir prévia e demonstração | Prévia executa a navegação; na demonstração o clique é representado sem navegar implicitamente | `ev-H007-02_*.png` |
| H007-03 | 007 US2/US3 | Alterar propriedade visual permitida e tentar personalização fora dos limites; salvar/reabrir e exportar | Aparência muda para valor permitido; valor fora do limite gera diagnóstico e preserva o último estado válido; texto editável nas saídas | `ev-H007-03_*.png` + `.pptx` |
| H008-01 | [008](../008-rotulo-efl/spec.md) US1 | Inserir rótulo com texto Unicode; alterar o conteúdo; mover, redimensionar e ordenar | Texto preservado com identificadores estáveis; ordem e geometria acompanham os demais controles | `ev-H008-01_*.png` |
| H008-02 | 008 US3 | Exportar PPT e gerar fontes com o rótulo acentuado | Texto editável na apresentação e conteúdo correspondente nos fontes gerados, sem perda de Unicode | `ev-H008-02_*.pptx` + `.java` |
| H009-01 | [009](../009-titulo-efl/spec.md) US1 | Inserir título independente e alterar texto, fonte, cor e dimensões | Título com identidade e texto próprios; alteração não afeta títulos de outras telas | `ev-H009-01_*.png` |
| H009-02 | 009 US3 | Exportar e inspecionar o título nas saídas | Título permanece texto editável, mantendo destaque visual | `ev-H009-02_*.pptx` |
| H010-01 | [010](../010-campo-texto-efl/spec.md) US1 | Editar valor inicial e percorrer os tipos texto, número, celular, CEP, CPF, data, DDD e NIS | Cada tipo aplicado ao campo; valor inicial configurado sem máscara inventada | `ev-H010-01_*.png` |
| H010-02 | 010 US1 | Digitar na prévia e reabrir a visita | Valor digitado não é gravado no projeto e o valor inicial é restaurado a cada visita | `ev-H010-02_*.txt` |
| H010-03 | 010 US2 | Vincular teclado numérico da mesma tela e digitar por etapas explícitas | Digitação chega ao campo vinculado, sem segundo algoritmo na UI | `ev-H010-03_*.png` |
| H011-01 | [011](../011-selecao-efl/spec.md) US1 | Adicionar, editar e reordenar opções de seleção sem usar legendas como identidade | Ordem e conteúdo preservados; rename de legenda não troca a identidade da opção | `ev-H011-01_*.png` |
| H011-02 | 011 US1 | Alternar orientação horizontal/vertical e definir opção inicial | Orientação aplicada e opção inicial opcional refletida no estado confirmado | `ev-H011-02_*.png` |
| H011-03 | 011 US2 | Selecionar uma opção na prévia, sair e voltar à tela | Uma única opção selecionada por visita; estado inicial restaurado ao retornar | `ev-H011-03_*.png` |
| H012-01 | [012](../012-relatorio-efl/spec.md) US1 | Editar título e linhas ordenadas do relatório | Título e linhas com IDs independentes dos textos; edição de texto não perde ordem | `ev-H012-01_*.png` |
| H012-02 | 012 US1 | Atualizar conteúdo mantendo o cabeçalho configurado | Cabeçalho separado das linhas e preservado durante a atualização | `ev-H012-02_*.png` |
| H012-03 | 012 US2 | Selecionar uma linha na prévia | Representação estática sem conexão a banco ou execução de regra de negócio | `ev-H012-03_*.png` |
| H013-01 | [013](../013-teclado-numerico-efl/spec.md) US1 | Arrastar o teclado, redimensionar e escolher campo de destino da mesma tela | Teclado posicionado e vinculado ao campo escolhido, sem vínculo implícito | `ev-H013-01_*.png` |
| H013-02 | 013 US1 | Digitar números, apagar e limpar o campo na prévia; reabrir a visita | Campo alterado apenas em runtime; valores iniciais persistidos intactos após reabrir | `ev-H013-02_*.png` |
| H013-03 | 013 US2 | Conferir arranjo 7–9, 4–6, 1–3, controles de cursor, 0, TAB e LIMPAR | Arranjo legado reproduzido; diferenças da versão atual (doze teclas) explicitadas, não ocultadas | `ev-H013-03_*.png` |
| H014-01 | [014](../014-data-hora/spec.md) US1 | Inserir o componente em tela vazia | Duas linhas: cidade/data e horário, em posições distintas | `ev-H014-01_*.png` |
| H014-02 | 014 US1 | Editar os valores ilustrativos e aguardar alguns minutos | Texto permanece o configurado; nenhuma consulta ao relógio da máquina nem alteração do projeto | `ev-H014-02_*.png` |
| H014-03 | 014 US2 | Mover e remover as linhas; desfazer a inclusão do conjunto | Linhas tratadas como componentes comuns; inclusão revertida em uma única ação | `ev-H014-03_*.png` |
| H015-01 | [015](../015-dados-loterica/spec.md) US1 | Inserir o componente | Três linhas de informação com textos ilustrativos independentes | `ev-H015-01_*.png` |
| H015-02 | 015 US1 | Alterar os valores e observar rede e arquivos | Nenhuma busca de dados reais de usuário, terminal ou estabelecimento | `ev-H015-02_*.txt` |
| H015-03 | 015 US2 | Salvar/reabrir e exportar; comparar visitas | Conteúdo, geometria e ordem preservados no arquivo e na apresentação | `ev-H015-03_*.pptx` |
| H016-01 | [016](../016-infraestrutura/spec.md) US1 | Inserir o componente no cabeçalho e em outra posição | Linhas GATEWAY, DNS e HTTP inseridas na posição escolhida | `ev-H016-01_*.png` |
| H016-02 | 016 US1 | Personalizar texto e aparência com a rede desconectada | Personalização aplicada sem executar teste de conectividade | `ev-H016-02_*.png` |
| H016-03 | 016 US2 | Exportar e inspecionar os indicadores | Indicadores preservados como texto editável nas representações da jornada | `ev-H016-03_*.pptx` |
| H017-01 | [017](../017-cabecalho-tela/spec.md) US1 | Inserir o conjunto de cabeçalho em jornada mínima válida | Oito rótulos alinhados no topo da tela | `ev-H017-01_*.png` |
| H017-02 | 017 US1 | Editar cada linha do cabeçalho | Linhas editáveis individualmente, sem converter o cabeçalho em imagem | `ev-H017-02_*.png` |
| H017-03 | 017 US1 | Desfazer a inclusão do conjunto | Inclusão completa revertida em uma ação, preservando os demais componentes | `ev-H017-03_*.png` |
| H018-01 | [018](../018-faixa-titulo/spec.md) US1 | Inserir a faixa e editar o título sobreposto | Faixa original com título editável por cima | `ev-H018-01_*.png` |
| H018-02 | 018 US1 | Mostrar, ocultar e remover a faixa | Faixa tratada como componente explícito, sem gerar texto invisível no modelo | `ev-H018-02_*.txt` |
| H018-03 | 018 US2 | Alterar título e largura da faixa | Separação entre recurso decorativo e texto preservada | `ev-H018-03_*.png` |
| H019-01 | [019](../019-canvas-edicao/spec.md) US1 | Abrir tela vazia e soltar componente dentro da superfície | Fundo original 800 × 600 e drop aceito dentro da superfície | `ev-H019-01_*.png` |
| H019-02 | 019 US1 | Selecionar, mover e redimensionar componentes | Somente componentes são afetados; o fundo permanece atrás deles | `ev-H019-02_*.png` |
| H019-03 | 019 US1 | Ajustar o zoom para ver a tela inteira | Dimensões lógicas da tela inalteradas pelo zoom | `ev-H019-03_*.txt` (bounds antes/depois) |
| H020-01 | [020](../020-elemento-edicao/spec.md) US1 | Selecionar componente sobreposto e confirmar um arraste e um resize | Componente superior sob o ponteiro selecionado; arraste e resize confirmados como uma única operação de histórico cada | `ev-H020-01_*.png` |
| H020-02 | 020 US1 | Abrir o menu de contexto de um controle e confirmar Excluir seleção com dois controles | Ambos removidos em uma única transação, restaurável por desfazer | `ev-H020-02_*.png` |
| H021-01 | [021](../021-elemento-previa/spec.md) US1 | Comparar o mesmo componente no editor e na prévia | Geometria e conteúdo idênticos entre editor e prévia | `ev-H021-01_*.png` (lado a lado) |
| H021-02 | 021 US1 | Operar botão, campo e seleção na prévia | Interação funciona com estado isolado da edição | `ev-H021-02_*.png` |
| H021-03 | 021 US1 | Digitar, fechar a prévia e inspecionar o documento | Snapshot independente preservado; dados digitados não entram no documento de edição | `ev-H021-03_*.txt` |
| H022-01 | [022](../022-lista-telas/spec.md) US1 | Criar, selecionar e renomear telas | IDs estáveis após renomear; seleção refletida no canvas | `ev-H022-01_*.png` |
| H022-02 | 022 US1 | Duplicar, remover e reordenar telas | Ordem resultante da ação explícita, sem ordenação implícita por título ou timestamp | `ev-H022-02_*.png` |
| H022-03 | 022 US2 | Solicitar remoção de tela referenciada | Referências afetadas informadas antes de confirmar; cancelamento preserva tudo | `ev-H022-03_*.png` |
| H023-01 | [023](../023-paleta-componentes/spec.md) US1 | Abrir a aba Paleta | Sete famílias e conjuntos auxiliares exibidos em aba própria | `ev-H023-01_*.png` |
| H023-02 | 023 US1 | Inserir componente por clique e por drop em coordenada da imagem | Inserção nas duas formas, com o componente na coordenada indicada | `ev-H023-02_*.png` |
| H023-03 | 023 US1 | Comparar família primitiva e conjunto pronto | Famílias primitivas separadas dos conjuntos prontos, com nomes amigáveis relacionados aos nomes EFL | `ev-H023-03_*.png` |
| H024-01 | [024](../024-arvore-estrutura/spec.md) US1 | Abrir Estrutura com a tela ativa selecionada | Apenas componentes da tela ativa, na ordem de camadas | `ev-H024-01_*.png` |
| H024-02 | 024 US1 | Selecionar item na lista, no canvas e no painel de propriedades alternadamente | Seleção sincronizada nos três locais | `ev-H024-02_*.png` |
| H024-03 | 024 US2 | Renomear, duplicar, remover, mudar camada e desfazer | Lista atualizada após cada operação, sem itens órfãos | `ev-H024-03_*.png` |
| H025-01 | [025](../025-barra-ferramentas/spec.md) US1 | Percorrer novo, abrir, salvar, prévia e as duas saídas | Comandos disponíveis e nomeados de forma clara, sem ação ambígua | `ev-H025-01_*.png` |
| H025-02 | 025 US1 | Usar desfazer/refazer, zoom, ajustar tela e configurar fundo | Cada comando produz o efeito esperado e pode ser revertido quando aplicável | `ev-H025-02_*.png` |
| H025-03 | 025 US2 | Alinhar componentes à esquerda, direita, topo, base e centros | Alinhamento conforme dimensões reais da tela, sem offsets ocultos do legado | `ev-H025-03_*.png` |
| H026-01 | [026](../026-painel-esquerdo/spec.md) US1 | Alternar as abas Telas, Paleta e Estrutura | Abas com responsabilidades próprias e conteúdo correto em cada uma | `ev-H026-01_*.png` |
| H026-02 | 026 US1 | Rolar cada conteúdo e arrastar a divisória | Rolagem independente por conteúdo e divisória redimensionável | `ev-H026-02_*.png` |
| H026-03 | 026 US2 | Trocar de aba com documento aberto e alcançar o último item da paleta | Documento inalterado e último item alcançável | `ev-H026-03_*.png` |
| H027-01 | [027](../027-painel-direito/spec.md) US1 | Abrir o painel sem seleção | Configuração da tela exibida | `ev-H027-01_*.png` |
| H027-02 | 027 US1 | Selecionar controles de famílias diferentes | Propriedades comuns mais as da família correta, sem misturar famílias | `ev-H027-02_*.png` |
| H027-03 | 027 US2 | Editar valores e aplicar; provocar erro de validação | Formulário, validação e aplicação transacional separados, com rolagem própria | `ev-H027-03_*.png` |
| H028-01 | [028](../028-propriedades/spec.md) US1 | Editar nome lógico, dimensões e propriedades específicas | Dados de outra família não aparecem nem são sobrescritos | `ev-H028-01_*.png` |
| H028-02 | 028 US1 | Personalizar cor, fonte, preenchimento, borda e imagem quando suportados | Alterações aplicadas apenas nos componentes que suportam a propriedade | `ev-H028-02_*.png` |
| H028-03 | 028 US1 | Aplicar valores válidos e, depois, introduzir valor inválido | Aplicação reversível em uma ação; erro conserva o estado anterior íntegro | `ev-H028-03_*.png` |
| H029-01 | [029](../029-previa-interativa/spec.md) US1 | Abrir a prévia do projeto | Janela independente com o mesmo fundo e controles da tela inicial | `ev-H029-01_*.png` |
| H029-02 | 029 US1 | Navegar por botões e operar campo, seleção e teclado | Navegação e interação funcionam com estado isolado do editor | `ev-H029-02_*.png` |
| H029-03 | 029 US1 | Sair e voltar a uma tela já visitada | Valores iniciais restaurados a cada visita | `ev-H029-03_*.png` |
| H030-01 | [030](../030-visualizacao-codigo/spec.md) US1 | Abrir visualização de código com perfil homologado | Arquivos organizados por tela, fábrica, navegação e integração | `ev-H030-01_*.png` |
| H030-02 | 030 US1 | Abrir a visualização sem bibliotecas reais disponíveis | Diagnóstico claro de perfil não verificado, sem inventar APIs | `ev-H030-02_*.png` |
| H030-03 | 030 US2 | Editar conteúdo e tentar executá-lo como template | Alvo Java 1.6 e documentação em português preservados; conteúdo editado não é executado | `ev-H030-03_*.txt` |
| H031-01 a H031-E10 | [031](../031-codigo-barras-efl/spec.md) | Ver matriz completa em 21.3 | Ver 21.3 | `ev-H031-*` |
| H032-01 a H032-E07 | [032](../032-codigo-barras-processado-efl/spec.md) | Ver matriz completa em 21.4 | Ver 21.4 | `ev-H032-*` |
| H033-01 a H033-E08 | [033](../033-menu-aplicacao-efl/spec.md) | Ver matriz completa em 21.5 | Ver 21.5 | `ev-H033-*` |

### 21.3 Matriz detalhada — 031 Código de barras

Fonte: [spec.md](../031-codigo-barras-efl/spec.md) (US1–US4, casos de borda e FR-001–FR-012) e [DESIGN.md](../031-codigo-barras-efl/DESIGN.md) (geometria, estados e textos). Dados de referência do [contrato compartilhado](contracts/component-catalog.md#contrato-codigos-barras-2026-10-23):

| Fixture | Conteúdo/uso | Observação |
| --- | --- | --- |
| F01 linha digitável | `00190500954014481606906809350314337370000000100` (47) | Três DVs locais e geral válidos sob as premissas; banco `001`, valor 100 centavos, fator `3737`, não CAIXA/SIGCB |
| F01 lido | `00193373700000001000500940144816060680935031` (44) | Conversão exata ida/volta com a linha acima |
| F02 lido sintético | `10498100000000001001234560000100010000000003` (44) | Banco `104`, 100 centavos, fator `1000`; SIGCB verdadeiro **somente sob a premissa SISPL experimental**; título `000100010000000003` |
| F02 digitado sintético | `10491234566000010001000000000034810000000000100` (47) | Conversão correspondente de 47; não é vetor independente de homologação |
| F03 DV local | Trocar o dígito da posição 10 da linha F01 de `5` para `4` | Deve falhar no bloco 1 e não avançar; repetir a mutação nos blocos 2 e 3 |
| F04 DV geral | Trocar somente a posição 5 do lido F01 de `3` para `4` e converter estruturalmente para linha | DVs locais permanecem válidos; conclusão e leitura falham pelo DV geral |
| F05 inválidos estruturais | Vazio como captura submetida; 43, 45, 46 e 48 dígitos; sinal, letras, Unicode numérico e quebra interna | Recusa da captura com o modal especificado, sem preenchimento parcial |

Textos exigidos, literais: `Código de barras inválido`, botão `OK, entendi` e o template de erro local `Bloco N inválido` (N = 1, 2 ou 3, conforme o par recusado).

| ID | Requisito | Pré-condição | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- | --- |
| H031-01 | 031/US1.1, FR-001 | Nova instância em tela com apenas o componente habilitado | Abrir a prévia | Oito campos vazios, somente o campo 1 recebe digitação, comprimentos 5, 5, 5, 6, 5, 6, 1 e 14 nessa ordem, sem placeholder numérico | `ev-H031-01_*.png` |
| H031-02 | 031/US1.2, FR-002 | Prévia aberta | Digitar até completar o 1º, o 3º e o 5º campo visual | Ao atingir o limite, o próximo campo visual é habilitado; nenhum DV é validado ainda nesse ponto | `ev-H031-02_*.png` |
| H031-03 | 031/US1.3, FR-002 | Campos 1 e 2 disponíveis | Digitar F01 no par do bloco 1 | Ao completar o 2º campo, o DV Módulo 10 é validado antes do avanço | `ev-H031-03_*.png` |
| H031-04 | 031/US1.3, FR-002 | Três blocos com F01 válido | Aplicar F03 no bloco 1; repetir nos blocos 2 e 3 | Erro impede ir ao bloco seguinte, identifica o bloco, devolve foco ao primeiro campo do par e conserva os demais blocos válidos; mensagem no template `Bloco N inválido` | `ev-H031-04_*.png` |
| H031-05 | 031/US1.4, FR-003 | Três primeiros blocos válidos | Informar o campo de 1 dígito e o de 14 posições | O campo de 1 dígito é tratado como DV geral e o de 14 como fator/valor; a conclusão exige DV geral Módulo 11 correto, não apenas os três DVs locais | `ev-H031-05_*.png` |
| H031-06 | 031/US1.5, FR-004 | Último campo ativo e incompleto | Solicitar concluir a entrada | Campo completado com zeros à direita até 14 posições e código inteiro validado; falha mantém a conclusão bloqueada e a correção recomeça no campo 1, sem apagar o rascunho | `ev-H031-06_*.png` |
| H031-07 | 031/US1.6 | Entrada válida | Concluir | Oito campos somente para leitura e ações explícitas da tela disponíveis, sem navegação automática causada pelo componente | `ev-H031-07_*.png` |
| H031-08 | 031/US1.7, FR-007 | Entrada em andamento; depois entrada concluída | Reiniciar a digitação nos dois estados | Estado desta entrada limpo, foco no campo 1, outras instâncias e o documento inalterados | `ev-H031-08_*.png` |
| H031-09 | 031/US2.1, FR-005 | Modo de leitura disponível | Submeter F01 lido (44) por leitura e por colagem | DV geral validado antes de converter e preencher; oito valores idênticos nos dois percursos; conclusão liberada e campos somente para leitura | `ev-H031-09_*.png` |
| H031-10 | 031/US2.1, FR-005 | Prévia aberta | Percorrer F01 por digitação, colagem de 44 e colagem de 47 | Os três percursos produzem exatamente os mesmos oito valores | `ev-H031-10_*.txt` (oito campos nos três percursos) |
| H031-11 | 031/US2.2, FR-005 | Linha digitável de 47 disponível | Colar F01 de 47 | Três DVs locais e o geral validados antes de publicar o preenchimento; nenhum campo recebe trecho de entrada recusada | `ev-H031-11_*.png` |
| H031-12 | 031/US2.3, FR-006 | Prévia aberta | Submeter captura inválida, incompleta, excedente e não numérica (F05) | Conclusão bloqueada, oito campos e derivados desta instância limpos, modal com exatamente `Código de barras inválido` e botão `OK, entendi` | `ev-H031-12_*.png` |
| H031-13 | 031/US2.4 | Modal de erro aberto | Acionar `OK, entendi`; repetir com Enter | Diálogo fecha e o foco volta ao primeiro campo vazio; nenhuma tecla acumulada da captura é reaplicada | `ev-H031-13_*.png` |
| H031-14 | 031/US2.5, FR-007 | Dispositivo configurado como teclado | Ativar o modo de leitura e encerrar com Enter | Toda a string é tratada como uma única entrada de 44 dígitos, sem inferência pela velocidade das teclas | `ev-H031-14_*.png` |
| H031-15 | 031/US3.1, FR-008 | F01 aceita | Consultar os dados derivados | Banco `001`, valor R$ 1,00, fator `3737`, documento CAIXA falso; zeros à esquerda preservados nas representações textuais | `ev-H031-15_*.png` |
| H031-16 | 031/US3.2, FR-008 | Sem entrada aceita | Consultar derivados com entrada vazia, parcial e inválida | Derivados indisponíveis, sem exceção por substring nem reaproveitamento de resultado anterior | `ev-H031-16_*.png` |
| H031-17 | 031/US3.3, FR-008 | F02 válido (banco `104`) | Verificar identificação CAIXA e classificação SIGCB | Identificação sem consulta externa; classificação SIGCB por regras locais com indicação de premissa não homologada | `ev-H031-17_*.txt` |
| H031-18 | 031/US3.4 | Entrada válida com fator não nulo (`3737` e fator iniciado em zero) | Consultar o fator | Fator exposto como número, inclusive quando começa em zero; nenhuma conversão silenciosa em data de vencimento | `ev-H031-18_*.txt` |
| H031-19 | 031/US4.1, FR-010 | Tela com componentes legados | Inserir e duplicar o componente por **ELEMENTOS DO LEGADO** | Instância única de Tela, não oito componentes soltos; ID de arquivo e propriedades comuns preservados e nova identidade na cópia | `ev-H031-19_*.png` |
| H031-20 | 031/US4.2, FR-009/FR-010 | Personalização confirmada | Alterar rótulo, geometria e estilo; salvar, reabrir e exportar | Rótulo, geometria e estilo permanecem; valores digitados, campos habilitados, mensagem e derivados da prévia não entram no projeto | `ev-H031-20_*.txt` + `.pptx` |
| H031-21 | 031/US4.3, FR-011 | PPT exportado com o componente | Inspecionar os campos no PowerPoint real | Textos e formas editáveis, sem digitação, leitora, cálculo de DV ou modal executável; limitação informada antes da exportação | `ev-H031-21_*.pptx` |

Casos de borda e limites da 031:

| ID | Regra | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- |
| H031-E01 | Blocos visuais não são blocos de validação | Contar dígitos por par: 10, 11 e 11 nos seis primeiros campos, mais 1 e 14 | Três pares de validação identificados; campos de 1 e 14 completam 47 | `ev-H031-E01_*.txt` |
| H031-E02 | Somente dígitos ASCII no conteúdo canônico | Colar com espaços e pontos de apresentação; depois com sinais, letras, quebras internas e símbolos | Apenas espaços e pontos de apresentação são removidos; os demais caracteres tornam a captura inválida; Enter final é terminador, não conteúdo | `ev-H031-E02_*.png` |
| H031-E03 | Captura integral incompatível | Submeter 43, 45, 46 e 48 posições; também payload PIX/QR textual | Recusa sem truncar; 48 posições de arrecadação e PIX/QR textual não são boleto de 44/47 deste componente | `ev-H031-E03_*.png` |
| H031-E04 | Campo visual incompleto não atravessa | Com o par incompleto, tentar TAB e clique em campo posterior; apagar e redigitar o par | TAB/clique não avançam; campos posteriores bloqueados até revalidação; par recusado pode ser corrigido | `ev-H031-E04_*.png` |
| H031-E05 | Instância ou contêiner desabilitado | Com a instância desabilitada, tentar digitação, colagem e leitura; depois usar teclado numérico vinculado | Nenhuma entrada é aceita no estado desabilitado; teclado numérico escreve no subcampo ativo com os mesmos limites e validações | `ev-H031-E05_*.png` |
| H031-E06 | Isolamento de runtime | Submeter captura inválida em uma instância com outra instância vizinha; cancelar uma captura antes de submeter | Apenas o runtime da instância destinatária é limpo; cancelar preserva o rascunho anterior; nenhuma entrada aparece em log | `ev-H031-E06_*.txt` |
| H031-E07 | Campo final vazio na conclusão | Acionar concluir com o último campo vazio | Equivale a 14 zeros antes da validação geral, nunca aprovação automática | `ev-H031-E07_*.png` |
| H031-E08 | Fator ausente | Aceitar código com fator `0000` e outro com fator iniciado em zero | `0000` significa ausente; os demais são expostos como números; data civil e ciclo de vencimento fora do escopo | `ev-H031-E08_*.txt` |
| H031-E09 | DPI e legibilidade | Repetir H031-01, H031-04 e H031-12 em 100%, 150% e 200% em 1366 × 768 | Foco legível, valores completos nos campos de 1 e 14 dígitos e mensagem que não depende apenas de cor | `ev-H031-E09_dpi100/150/200.png` |
| H031-E10 | Editor não simula boleto | Selecionar, mover e redimensionar o composto no canvas; exportar | Campos vazios no editor, sem digitação de boleto no canvas; seleção e alças não entram no PPT | `ev-H031-E10_*.png` |

### 21.4 Matriz detalhada — 032 Código de barras processado

Fonte: [spec.md](../032-codigo-barras-processado-efl/spec.md) (US1–US3, casos de borda e FR-001–FR-010) e [DESIGN.md](../032-codigo-barras-processado-efl/DESIGN.md). Os cinco textos de rótulo são derivados do tipo escolhido, seguidos de dois-pontos:

| Opção do combo | Texto exibido na jornada |
| --- | --- |
| CPF | `CPF:` |
| CNPJ | `CNPJ:` |
| Nome | `Nome:` |
| Nome fantasia | `Nome fantasia:` |
| Razão social | `Razão social:` |

| ID | Requisito | Pré-condição | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- | --- |
| H032-01 | 032/US1.1, FR-005 | Nova instância | Consultar as propriedades | Um bloco com título vazio e duas linhas iniciais, CNPJ e Nome, ambas com valor vazio; nenhum documento ou nome fictício inserido | `ev-H032-01_*.png` |
| H032-02 | 032/US1.2, FR-002 | Dois blocos configurados | Alterar o título e o valor de uma linha do segundo bloco | Somente essa linha e esse bloco mudam; o primeiro bloco permanece integralmente igual | `ev-H032-02_*.png` |
| H032-03 | 032/US1.3, FR-002 | Blocos e linhas existentes | Adicionar, reordenar e remover um item e confirmar; desfazer e refazer | Ordem visual acompanha a configuração; operação reversível sem perder conteúdo nem identidades dos itens restantes | `ev-H032-03_*.png` |
| H032-04 | 032/US1.4 | Bloco com valores vazios | Visualizar a tela na prévia | Rótulos permanecem com valores em branco, sem dado inventado e sem falso estado de processamento concluído | `ev-H032-04_*.png` |
| H032-05 | 032/US2.1, FR-003/FR-004 | Linha com CNPJ e valor configurado | Selecionar CPF no combo de propriedades | Rótulo passa a `CPF:` e o valor permanece disponível para edição, sem apagar, consultar ou converter silenciosamente | `ev-H032-05_*.png` |
| H032-06 | 032/US2.2 | Linha Nome com valor contendo acentos, espaços e vírgulas | Escolher Nome fantasia e depois Razão social | Texto correspondente exibido e valor preservado literal, inclusive acentos, espaços e vírgulas | `ev-H032-06_*.png` |
| H032-07 | 032/US2.3 | Dois blocos com escolhas diferentes | Confirmar a configuração e percorrer as cinco opções nos dois blocos | Cada linha conserva sua própria escolha e valor, sem seletor global substituir os demais | `ev-H032-07_*.png` |
| H032-08 | 032/US2.4 | Rascunho com seleções alteradas | Cancelar a edição | Rótulos, valores, ordem e geometria confirmados não mudam; popup do combo não entra no documento nem na tela exportada | `ev-H032-08_*.txt` |
| H032-09 | 032/US3.1, FR-007 | Configuração confirmada com mais de um bloco | Salvar e reabrir; desfazer e refazer uma alteração | Títulos, tipos de rótulo, valores, ordem e estilo permanecem exatamente configurados | `ev-H032-09_*.png` |
| H032-10 | 032/US3.2 | Instância duplicada | Editar uma linha da cópia e comparar com o original | Original inalterado; cópia com novos IDs de componente, blocos e linhas, sem perder valores ou seletores | `ev-H032-10_*.txt` |
| H032-11 | 032/US3.3, FR-008 | Tela na prévia e no PPT | Observar rótulos e valores nos dois alvos | Mesmos dados configurados como conteúdo de leitura; no PPT, rótulos e valores são textos editáveis, não screenshot | `ev-H032-11_*.pptx` |
| H032-12 | 032/US3.4, FR-009 | Entrada de código de barras válida na mesma tela | Concluir e depois reiniciar a entrada | Nomes e documentos da ficha não são preenchidos, apagados ou alterados; nenhuma navegação é disparada | `ev-H032-12_*.png` |

Casos de borda e limites da 032:

| ID | Regra | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- |
| H032-E01 | Estrutura mínima | Remover todos os blocos; depois todas as linhas de um bloco; confirmar | Configuração estruturalmente vazia recusada, com ao menos um bloco e uma linha; último estado confirmado preservado | `ev-H032-E01_*.png` |
| H032-E02 | Nova linha e enum | Adicionar linha e conferir default; tentar enum desconhecido | Nova linha começa como Nome, valor vazio e novo ID; enum desconhecido gera erro localizado, não rótulo por aproximação | `ev-H032-E02_*.txt` |
| H032-E03 | Valor ilustrativo | Trocar CPF por CNPJ com valor pontuado e com zeros à esquerda | Nenhum DV calculado, comprimento não alterado, zeros não completados; pontuação e máscara permanecem literais | `ev-H032-E03_*.png` |
| H032-E04 | Independência de linhas | Repetir o mesmo tipo em linhas e combinações iguais entre blocos | Linhas repetidas e combinações iguais permitidas; seletor altera uma linha por vez, sem inferir pessoa física ou jurídica | `ev-H032-E04_*.png` |
| H032-E05 | Título vazio e texto longo | Deixar título vazio; inserir nome muito longo | Título vazio não gera título fictício nem reserva linha; nome longo quebra ou recebe diagnóstico de overflow, nunca truncamento persistido | `ev-H032-E05_*.png` |
| H032-E06 | Identidades e duplicação | Reordenar, remover e duplicar com itens repetidos | Itens restantes mantêm identidade; duplicação não compartilha listas mutáveis nem UUIDs internos | `ev-H032-E06_*.txt` |
| H032-E07 | DPI e Unicode | Repetir H032-01, H032-05 e H032-11 em 100%, 150% e 200% com texto acentuado | Seletores navegáveis por teclado, foco visível, textos completos e exportação sem perda de Unicode | `ev-H032-E07_dpi100/150/200.png` |

### 21.5 Matriz detalhada — 033 Menu aplicação

Fonte: [spec.md](../033-menu-aplicacao-efl/spec.md) (US1–US3, casos de borda e FR-001–FR-010) e [DESIGN.md](../033-menu-aplicacao-efl/DESIGN.md). Referências numéricas que a execução deve conferir:

| Referência | Conteúdo esperado |
| --- | --- |
| Raiz | 12 serviços + 10 jogos; serviços à esquerda, jogos à direita |
| Ordem dos serviços | boleto, convênios, benefícios, financeiros, atendimento, resultados, apoio, bolão, relatórios, paga premiados, meios de pagamento, ofertar |
| Últimas linhas | y 369: APOIO e BOLÃO; y 442: RELATÓRIOS e PAGA PREMIADOS; y 515: MEIOS DE PAGAMENTO e OFERTAR |
| Geometria-base | x 20/200, topo 150, pitch 73, botões 164 × 68; clip serviços largura 400 até y 600, última linha termina em 583 |
| PAGA PREMIADOS | PROGNÓSTICOS `1500`, LOTERIA FEDERAL `1510`, INSTANTÂNEA `1513` |
| MEIOS DE PAGAMENTO | PIX → RECEBIMENTO e CONSULTAR RECEBIMENTO opção `99`; CARTÃO CRÉDITO e DÉBITO `99`; ENTRADA EM DINHEIRO `2` |
| RELATÓRIOS | 148 folhas e 27 títulos; zero inicial em `080206`; dias `1..8`; auditoria `9..16`; pagamento `[0,2,3,4,5,6,7,1]` |
| BOLÕES DIGITAIS | código `803`, último filho de BOLÃO (não é serviço root); layout `bolao` |
| Grade financeira | layout `financial` com botões 140 × 80 unidades lógicas |

Legendas exatas do bolão, na ordem, para o layout `bolao` de três botões:

1. `Visualizar cotas digitais para venda na lotérica.`
2. `Consulta bolões de concursos abertos para comercialização.`
3. `Consulta bolões de concursos encerrados e não prescritos.`

| ID | Requisito | Pré-condição | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- | --- |
| H033-01 | 033/US1.1, FR-001/FR-004 | Itens legados disponíveis | Inserir **Menu da aplicação** em tela vazia | Uma única instância com árvore inicial de serviços à esquerda e jogos à direita, baseada no Java/XML fornecido; Código de barras e Código de barras processado no mesmo grupo | `ev-H033-01_*.png` |
| H033-02 | 033/US1.2, FR-003 | Menu inserido | Editar legenda, ordem, categoria, estilo, habilitação e destino explícito; adicionar e remover item e submenu; salvar | Alteração inteira confirmada de uma vez, sem mudar os demais itens; Descartar e desfazer/refazer preservam o último estado confirmado | `ev-H033-02_*.png` |
| H033-03 | 033/US1.3, FR-007 | Árvore com nomes repetidos, Unicode e menus vazios | Salvar, reabrir e duplicar o componente | Conteúdo e ordem permanecem; cópia independente sem reutilizar identidades internas da original | `ev-H033-03_*.txt` |
| H033-04 | 033/US1.4, FR-004 | Árvore com Bolão aninhado | Selecionar o layout `bolao` e conferir a composição | Três botões com as legendas da tabela acima, na ordem; outras árvores continuam editáveis e não viram telas adicionais | `ev-H033-04_*.png` |
| H033-05 | 033/US2.1, FR-002/FR-008 | Raiz visível na prévia | Ativar item habilitado com submenu; depois um item desabilitado | A mesma interface mostra apenas o submenu selecionado e mantém o caminho local; item desabilitado não produz efeito | `ev-H033-05_*.png` |
| H033-06 | 033/US2.2, FR-005 | Submenu aninhado aberto | Usar Voltar; usar Voltar na raiz; usar Cancelar a partir de outro submenu inicial | Voltar remove um nível e não muda nada na raiz; Cancelar retorna à raiz mesmo quando o menu inicial configurado era outro submenu | `ev-H033-06_*.png` |
| H033-07 | 033/US2.3, FR-006 | Folha sem tela de destino | Selecioná-la na prévia | Informação de seleção simulada, sem executar movimento, consultar conta ou inferir rota por texto ou metadado | `ev-H033-07_*.png` |
| H033-08 | 033/US2.4, FR-006 | Folha com destino autorado; Prévia Tela Atual disponível | Ativar a folha na prévia completa; repetir na Prévia Tela Atual | Somente a tela existente configurada é aberta no modo completo; na Prévia Tela Atual a navegação entre telas permanece desativada | `ev-H033-08_*.png` |
| H033-09 | 033/US2.5, FR-005/FR-007 | Prévia aberta com caminho percorrido | Escolher Sair; reabrir ou revisitar a prévia | Sair encerra apenas a sessão de prévia, sem afetar editor ou autenticação; reabrir restaura o menu inicial configurado sem persistir o caminho | `ev-H033-09_*.png` |
| H033-10 | 033/US3.1, FR-007/FR-009 | Submenu escolhido como inicial | Salvar, reabrir, gerar quadro estático e percorrer outro ramo na prévia livre | Submenu inicial, sua ordem e suas legendas exibidos no canvas e no quadro; o caminho da prévia livre não interfere | `ev-H033-10_*.png` |
| H033-11 | 033/US3.2, FR-008 | Submenu inicial escolhido | Removê-lo sem escolher substituto válido e confirmar | Erro localizado e nenhuma alteração parcial publicada; não trocar silenciosamente para raiz | `ev-H033-11_*.png` |
| H033-12 | 033/US3.3, FR-009 | Quadro válido | Exportar e abrir no PowerPoint real | Textos e formas editáveis representando o estado escolhido, sem prometer árvore interativa nem operação bancária no arquivo | `ev-H033-12_*.pptx` |

Casos de borda e limites da 033:

| ID | Regra | Passos | Resultado esperado | Evidência |
| --- | --- | --- | --- | --- |
| H033-E01 | Faixa e título | Abrir a raiz e um submenu com faixa | Raiz sem faixa nem título; submenus com faixa e título branco do token `label_faixa_titulo_principal`, inclusive quando o inicial é submenu | `ev-H033-E01_*.png` |
| H033-E02 | Menus vazios e rótulos extremos | Criar menu vazio; usar rótulos repetidos, longos e acentuados | Menu vazio válido na autoria; rótulos preservados literalmente; texto longo quebra ou recebe diagnóstico de overflow | `ev-H033-E02_*.png` |
| H033-E03 | Profundidade | Montar árvore com 16 níveis contando a raiz; depois 17 | 16 níveis preservados; 17 recusados antes de publicar, sem recursão sem limite | `ev-H033-E03_*.txt` |
| H033-E04 | Contagem de itens | Validar fixture com 1000 itens; depois 1001 | 1000 itens preservados; 1001 recusado com erro localizado. As fixtures 16/1000 são a única evidência; não criar 1000 itens manualmente | `ev-H033-E04_*.txt` |
| H033-E05 | Estruturas inválidas | Inserir IDs repetidos, ciclo em memória, destino inexistente e item com submenu e destino simultâneos | Cada caso recusado com localização do nó, sem publicação parcial e sem estado parcial no projeto | `ev-H033-E05_*.png` |
| H033-E06 | Habilitação efetiva | Desabilitar ramo e instância; tentar ativar seus itens | Instância desabilitada bloqueia todas as regiões; ramo desabilitado não é atravessado; nenhuma busca de rota por nome, atalho ou código de movimento | `ev-H033-E06_*.png` |
| H033-E07 | Recorte da grade de serviços | Rolar e recortar a grade de serviços até a última linha | Última linha y 515–583 visível integralmente, sem sobreposição, sem mover ou reduzir botões e sem podar dados | `ev-H033-E07_*.png` |
| H033-E08 | DPI, PPT e roteiro | Repetir H033-01, H033-07 e H033-12 em 100%, 150% e 200%; tentar `typeText` no composto | Hit regions e legibilidade coerentes nos quatro layouts; clique é visual sem dispatch de item; `typeText` recusado com identificação da etapa | `ev-H033-E08_dpi100/150/200.png` |

### 21.6 Limitações e pendências desta matriz

- **Imagens e anexos ausentes.** Os anexos citados para 031, 032 e o botão incorreto do menu chegaram vazios. Enquanto não houver imagem recebida, a fidelidade visual de 031, 032 e 033 permanece **PENDENTE**: os cenários de comparação visual só podem ser registrados como NÃO VERIFICÁVEL, nunca como OBSERVADO-CONFORME.
- **Classes SISPL e SIGCB reais ausentes.** `Modulo10`, `Modulo11Bloqueto` e `Modulo11Sispl` não estão no workspace. A equivalência com a regra bancária real e a classificação SIGCB permanecem **PENDENTES**; a marcação de premissa não homologada deve continuar visível no resultado (H031-17).
- **Seed do menu sem o XML completo no workspace.** A árvore de 12 serviços + 10 jogos, as posições, os literais de PAGA PREMIADOS, MEIOS DE PAGAMENTO e RELATÓRIOS derivam de transcrição da conversa, não de um arquivo XML local. Sem a fonte disponível ao executor, **não afirmar fidelidade** ao legado: registrar a ausência e pedir reposição do XML. O seed não reescreve instâncias já autoradas; é preciso inserir um novo Menu da aplicação para recebê-lo.
- **Botão incorreto sem captura.** O relato do botão errado não é reproduzível sem captura ou nome do botão. Não declarar esse defeito identificado, corrigido ou aceito; a pendência é obter a imagem para investigar.
- **1000/1001 apenas fixture.** Os limites de 16 níveis e 1000 itens são verificados por fixtures compiladas, não por interação real com uma árvore de 1000 itens. H033-E04 é o cenário máximo admitido e não comprova usabilidade em escala.
- **Executável desta revisão.** O EXE canônico e a cópia de compatibilidade existem em `apps/gerador-jornada/dist` e têm o mesmo SHA-256 `d43cd0eb…`; o pacote ZIP tem `cb02b71f…`. Abrir o EXE e percorrer a matriz continua sendo trabalho do executor humano.
- **Verificação técnica desta revisão, não manual.** Foram executados `gofmt`, `go vet ./...`, `go test ./...` (todas as suítes disponíveis PASS, incluindo as novas fixtures de 031–033), três testes Python de empacotamento PASS, build Windows amd64 e uma captura automática de inicialização (`dist/evidencias/startup.png`). Isso comprova compilação, testes de domínio/contrato e inicialização da janela; **não** comprova nenhum cenário manual desta matriz, DPI, PowerPoint ou aceite. Nenhum deles fecha T138, T139, T143, T144 ou T146.
- **Office e licenças.** PowerPoint real é necessário para os cenários de PPT e não foi usado. Fonte declarada não prova licença nem métricas; ausência de licença continua fora do escopo desta matriz.

### 21.7 Modelo de evidência e regra de aceite

Cada cenário executado gera um par de evidências na pasta de registro escolhida pelo executor:

| Campo | Regra |
| --- | --- |
| Nome do arquivo | `ev-<ID-do-cenário>_<revisão-curta>_dpi<100/150/200>_<executor>.<ext>`; usar o ID da matriz (`H031-04`, `H033-E04`) e não renomear depois |
| Hash do EXE | SHA-256 do arquivo efetivamente aberto, conferido com `Get-FileHash`; divergência de `953f9e10…` implica outro registro |
| DPI e resolução | `100%`, `150%` ou `200%` em 1366 × 768; uma evidência por passagem quando a matriz exigir três |
| Data e hora | início e fim da execução, com fuso |
| Executor | nome de quem observou; agente não substitui executor humano |
| Conteúdo | captura da tela ou do PPT, saída textual da observação e, quando houver erro, o texto exibido literalmente |

Status permitidos por cenário, sem sinônimos:

| Status | Quando usar |
| --- | --- |
| **OBSERVADO-CONFORME** | Execução com evidência anexada e resultado igual ao esperado |
| **OBSERVADO-DIVERGENTE** | Execução com evidência e resultado diferente do esperado; descrever a diferença e abrir correção por requisito |
| **NÃO EXECUTADO** | Sem execução ou sem evidência anexada; estado inicial de todos os cenários desta seção |
| **NÃO VERIFICÁVEL** | Execução impedida por dependência ausente, como imagem de referência, classe SISPL real, XML completo do menu ou PowerPoint |

**Aceite explícito.** A aprovação só existe com manifestação explícita do usuário sobre a revisão testada, conforme a constituição. Autorização para implementar, plano aprovado, checklist de qualidade, compilação, `go vet`, teste Python ou revisão documental **não** são aceite. T139 é o dono do manual e do aceite de código de barras (031/032, incluindo SISPL/SIGCB); T144 é o dono do manual e do aceite do menu e do lote 031/032/033; T146 cobre a verificação comportamental do menu corrigido. Nenhuma delas é fechada por esta seção.

**Separação entre compilação e teste.** Compilar fixtures, compilar suítes e empacotar o EXE não é executar comportamento. `go test -c` produz binários de teste, não resultados; um build concluído não demonstra nenhum cenário desta matriz. Só conta como resultado o que tiver sido executado na revisão registrada, na janela e no ambiente declarados, com evidência anexada.

### 21.8 Registro de execução

Tabela a preencher pelo executor humano. **Nenhum cenário manual foi executado nesta sessão**: todas as linhas permanecem ausentes e todos os cenários em **NÃO EXECUTADO**. A verificação técnica (compilação, testes Go/Python, build e captura de inicialização) está registrada em 21.6 e não substitui esta tabela.

| Cenário | Executor | Revisão | Data | Resultado observado | Evidência |
| --- | --- | --- | --- | --- | --- |
|  |  |  |  |  |  |
|  |  |  |  |  |  |
|  |  |  |  |  |  |

Nota de encerramento desta sessão: a matriz 001–033 foi especificada e a revisão Windows foi compilada, testada tecnicamente e aberta para capturas de inicialização e dos compostos. **Nenhum cenário manual foi executado**; a tabela 21.8 permanece vazia e todos os cenários continuam em NÃO EXECUTADO. O EXE desta revisão tem 32.754.176 bytes e SHA-256 `953f9e10…`; o ZIP tem 19.811.611 bytes e SHA-256 `a74cbb51…`. As seções 17–20 e seus hashes permanecem históricos e intactos.

### 21.9 Evidência técnica desta revisão — 2026-10-24 (não é teste humano)

Esta subseção registra somente o que foi **executado** pelo executor técnico nesta revisão. Ela não preenche a tabela 21.8, não substitui os gates manuais T139/T144/T146 e não autoriza aceite.

| Verificação | Comando / artefato | Resultado observado |
| --- | --- | --- |
| Formatação | `gofmt -l internal cmd` | vazio |
| Análise estática | `go vet ./...` | PASS (sem diagnósticos) |
| Testes Go | `go test ./...` | PASS em todos os pacotes com testes; inclui as fixtures novas de catálogo, apresentação, PPTX, persistência, runtime e inspector |
| Testes Python | `python -X utf8 -m unittest discover -s scripts -p test_build.py -v` | 3 testes PASS |
| Build Windows | `python -X utf8 scripts/build.py --target windows --arch amd64` | distribuição gerada |
| EXE canônico | `dist/windows-amd64/JornadaStudio.exe` | 32.754.176 bytes · `953f9e102ff1cbbf8d007cf5fbebcfcb86b2a458c2533ff59f90401a710605f1` |
| Cópia idêntica | `dist/JornadaStudio.exe` | mesmo tamanho e hash do canônico |
| Pacote | `dist/JornadaStudio-1.0.0.0-windows-amd64.zip` | 19.811.611 bytes · `a74cbb5115a253351315cf29a2ba740ce960917ddabf18134e334e4aea20f3d8` |
| Inicialização real | `JornadaStudio.exe -smoke-screenshot dist/evidencias/startup.png` | janela 1366 × 768 capturada; processo encerrou com código 0; PNG 147.901 bytes · `d621bf3e0c126f8edfe708424c0bd24bfce79a554527cb3e8263f665b10876b9` |
| Compostos 031/032 no canvas | `JornadaStudio.exe -open dist/teste-compostos.jornada -screen 0 -smoke-screenshot dist/evidencias/compostos-031-032.png` | PNG 149.270 bytes · `a9f0b5d56f8e3ec1ccc64aad9a3dc6f00b7cde8a2f44c3872bcce3cfa5c8fb65`; oito campos e ficha processada visíveis |
| Compostos 033 no canvas | `JornadaStudio.exe -open dist/teste-compostos.jornada -screen 1 -smoke-screenshot dist/evidencias/compostos-033.png` | PNG 152.019 bytes · `dc79a411a737cfa76a935a4d6bdee3bcfaa183f8af5a8eb88dbd14f21ad8b4db`; 12 serviços + 10 jogos, última linha inteira |
| Projeto de teste | `JornadaStudio.exe -save-compounds dist/teste-compostos.jornada` | 21.705 bytes · `17487e9105013adb0c45d4417e922af258dbf8bc7b61ca1ace7d328eaf079878`; manifesto schema 3/catalog 3; 12 serviços + 10 jogos; BOLÕES DIGITAIS como último filho de BOLÃO; RELATÓRIOS com 27 títulos/148 folhas; 031 e 032 na primeira tela |
| PPT de teste | `JornadaStudio.exe -export-compounds dist/teste-compostos.pptx` | 25.120 bytes · `bdaea3dd347cd74d9281fe071a87f34c7570e0df00e6cc390460977177130c98`; slides com textos/formas editáveis (`txBox="1"`), sem screenshot dos compostos |
| Log de aplicação | execução de `-save-compounds` e inspeção de `%LOCALAPPDATA%\GeradorJornada\logs` | JSONL de startup criado com `executable.sha256` igual ao EXE acima, PID/sessão/versão/Go/SO/arquitetura; sem texto do projeto |
| Hierarquia/habilitação | `go test ./internal/screen ./internal/application ./internal/adapter/projectzip` | PASS: validação de pai/ciclo, movimento de subárvore uma vez, camadas relativas/extremos, remapeamento ao duplicar, órfãos ao remover e default de projetos antigos |
| Skins do acervo | `material-legado/elemetos_componente/xml/skin.xml` embutido em `internal/skin` | 10 tipos interpretados (48 painéis, 91 rótulos, 55 botões, campo/teclado etc.), 224 recursos resolvidos; testes de parser/render/PPT PASS; fixture usa `label_titulo_1` e `botao_confirmar`; `style.skin` persistido; imagem incorporada ao PPTX |
| Runtime sem persistência | inspeção do `manifest.json` gerado | nenhuma chave de valores/foco/caminho; somente `label`, `blocks` e `root/initialMenuId` autorados |

Limites desta evidência: o seed do menu foi gerado a partir da transcrição documentada (não do XML completo ausente); as imagens de referência e as classes SISPL/SIGCB reais continuam ausentes; a captura de inicialização não exercita os cenários da matriz. O executor humano deve abrir o EXE acima e registrar cada cenário em 21.8.

## 22. Reimplementação visual do legado — 2026-10-24

**Estado: pronto_para_teste_manual.** Versão **1.0.0.1**, Windows amd64. Executor dos checks e capturas: Codex. Teste/aceite humano desta revisão: **pendente**. As evidências abaixo pertencem a esta revisão; não herdam aprovação de builds anteriores.

### Entrega concreta

- Executável: [JornadaStudio.exe](../../apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe), **33.093.120 bytes**; FileVersion/ProductVersion `1.0.0.1`.
- Pacote: JornadaStudio-1.0.0.1-windows-amd64.zip (`../../apps/gerador-jornada/dist/JornadaStudio-1.0.0.1-windows-amd64.zip`; artefato histórico ausente no checkout atual).
- Projeto das sete famílias: estilos-legado.jornada (`../../apps/gerador-jornada/dist/estilos-legado.jornada`; artefato histórico ausente no checkout atual).
- Projeto dos compostos: compostos-legado.jornada (`../../apps/gerador-jornada/dist/compostos-legado.jornada`; artefato histórico ausente no checkout atual).
- Saídas editáveis: estilos-legado.pptx (`../../apps/gerador-jornada/dist/estilos-legado.pptx`; artefato histórico ausente no checkout atual) e compostos-legado.pptx (`../../apps/gerador-jornada/dist/compostos-legado.pptx`; artefato histórico ausente no checkout atual).
- Captura Fyne em memória da prévia: [legacy-preview.png](../../apps/gerador-jornada/build/legacy-preview.png), 800 × 600, após exercício real dos controles de cursor/TAB/Backspace.
- Captura de startup/renderização do executável Windows: [legacy-menu-native.png](../../apps/gerador-jornada/build/legacy-menu-native.png), framebuffer 1366 × 768, zoom do editor 60%, tema Escuro existente preservado. Isso não é teste humano de todos os fluxos.

| Artefato | SHA-256 |
| --- | --- |
| EXE | `709a81dd06090398db3df64e076db5ccd505edbd3bfc75735991e2ccdc5ffd2a` |
| ZIP | `39256c9ae9ce82a4804e8ec217bff6920d1d931bd0a9aa0a9666dcf759ab067b` |

### Verificações executadas e observações

| Verificação | Resultado observado |
| --- | --- |
| Regressões antes da correção | Falharam botão sem skin, teclado com altura 280/12 teclas e relatório sem decorações originais. |
| `go test ./...` | PASS nos 13 pacotes com testes; novos casos exercitaram aparência real, alpha/preto, fonte SFNT, 14 teclas, bounds antigos, cursor, TAB, Backspace, isolamento runtime e recurso indisponível. |
| `go vet ./...` e gofmt | PASS; executados também pelo check do empacotador. |
| `python -X utf8 scripts/import_legacy_skin.py --check` | PASS; 205 chaves tipadas operacionais, oito definições indisponíveis, 224 recursos referenciados. |
| Testes Python do empacotamento | PASS, três testes de identificação PE/ELF/Mach-O. |
| `python -X utf8 scripts/build.py --target windows --arch amd64 --check` | PASS; build PE amd64, pacote 1.0.0.1 e metadados nativos corretos. |
| Execução do EXE | PASS: `--save-example`, `--save-compounds`, `--open ... --export-example ...` e startup/captura de menu via `--smoke-screenshot`; todos encerraram com código 0. |
| PPTs reais produzidos pelo EXE | PASS de ZIP/XML: 27 partes XML/rels no exemplo e 15 no composto. Decorações originais incorporadas e textos/formas editáveis presentes. Não abertos no PowerPoint. |
| Integridade da distribuição | PASS: CRC do ZIP, SHA256SUMS, EXE contido idêntico ao entregue, cópia de compatibilidade idêntica e ausência de backups terminados em `~`. |
| Inspeção das capturas | O agente observou decoração azul original, fontes físicas, fundo Loterias, relatório com cabeçalho/rodapé e teclado 7–9/4–6/1–3/setas/0/TAB/LIMPAR; no menu nativo aparecem decoração e setas dos submenus. |
| Aceite Windows/DPI/Office | PENDENTE; captura automatizada não substitui matriz DPI 100/150/200%, uso humano e Office real. |

### Como iniciar e repetir

A partir de `apps/gerador-jornada`:

```powershell
./dist/windows-amd64/JornadaStudio.exe --open ./dist/estilos-legado.jornada
# Ou abrir os compostos, com a tela do menu selecionada:
./dist/windows-amd64/JornadaStudio.exe --open ./dist/compostos-legado.jornada --screen 1
```

Checks e reprodução da captura automatizada da prévia:

```powershell
./scripts/check.ps1
. ./scripts/env.ps1
$env:JORNADA_VISUAL_CAPTURE = "$PWD/build/legacy-preview.png"
go test ./internal/adapter/desktop -run TestLegacyPreviewRuntime -count=1
Remove-Item Env:JORNADA_VISUAL_CAPTURE
```

### Roteiro humano desta revisão

1. Abrir `estilos-legado.jornada`. Comparar botão, rótulo, título, campo, seleção, relatório e teclado com o acervo. Alterar tema do editor: a jornada conserva suas cores/fontes.
2. Inserir novo botão sem selecionar skin: deve usar o estilo original. Em Propriedades, testar confirm/cancel/yes/no, uma skin disponível, família da fonte e cores. Salvar alterações; Descartar deve preservar a última configuração confirmada.
3. Abrir Prévia Tela Atual: clicar no campo e digitar, selecionar texto e colar. O campo deve manter decoração EFL; foco usa fundo azul e texto claro. Inserir no meio do valor com o teclado virtual; setas movem cursor, LIMPAR apaga um caractere anterior/seleção e TAB avança entre campos habilitados. Fechar/reabrir a prévia reinicia os valores; o documento não recebe dados runtime.
4. Inserir teclado novo: default 200 × 350, 14 teclas e margens Java. Salvar/reabrir um teclado com bounds antigos ou personalizados: dimensões permanecem iguais. Conferir movimento, resize e Undo/Redo.
5. Na segunda tela do exemplo, conferir relatório e seleção. Em Prévia, selecionar opções sem perder o estilo. Desabilitar componente no inspector e verificar que continua visível sem interagir.
6. Abrir `compostos-legado.jornada`; conferir fichas e menu. Menu raiz não tem faixa; entrar em submenu mostra sua faixa branca sobre decoração original. Conferir rótulos longos, itens desabilitados e navegação previamente configurada.
7. Salvar/reabrir as personalizações. Exportar PPT e abrir no Office com as fontes reais disponíveis/licenciadas. Comparar imagens de decoração, fonte, texto editável, geometria, quebra de linha e navegação. Repetir nos DPIs 100%, 150% e 200% do Windows.
8. Registrar versão/hash, ambiente/DPI, executor, cenário, resultado observado e evidência. Corrigir/retestar falhas antes de aprovar explicitamente esta revisão.

### Limites do incremento

Este build reaplica estilos às famílias/compostos existentes no código atual. Não entrega todos os tipos, variantes de relatório, importação/perfis portáteis e expansões de formato presentes no backlog das specs; estes não foram silenciosamente concluídos. O inventário estático atual não equivale ao inventário histórico do perfil completo.

As oito definições com recursos ausentes ficam fora da oferta. Escolhas antigas indisponíveis conservam seus dados, exibem diagnóstico e bloqueiam PPT em vez de usar um substituto. As fontes físicas acompanham o editor, mas não são incorporadas dentro do PPT; disponibilidade/licença e equivalência no Office de destino continuam pendentes de validação. Java/EFL permanece desabilitado. Nenhum aceite humano, publicação ou merge foi registrado.

## 23. Reconciliação documental das 33 specs — 2026-10-24

**Escopo:** documentação sobre a revisão executável 1.0.0.1 da seção 22. Este trabalho não altera comportamento, não cria um novo teste funcional e não registra aceite humano. Os resultados técnicos e hashes da seção 22 continuam pertencendo à execução ali identificada.

- As 33 `spec.md` agora têm estado verificado próprio, fontes presentes e limites explícitos; os 29 DESIGNs existentes indicam a aplicação visual corrente e o que ainda é requisito de expansão.
- Planos/modelos dos seis domínios e contratos foram reconciliados com writer 3/3, dez Kinds e Style/biblioteca global atuais. Perfil visual portátil, aparência tipada, nove tipos adicionais e outras capacidades futuras continuam requisitos abertos, não foram apagados para corresponder ao código.
- Defaults/tokens corrigidos: skins do legado, fontes físicas, transparência, título 14 e teclado 200 × 350/14 teclas. Reader reconhece pares de versões nos campos suportados; não aceita appearanceProfile por reconhecer 3/2.
- O README da raiz explica causas, etapas, arquivos e comandos usados para recuperar o estilo; o guia da aplicação foi sincronizado. O menu RELATÓRIOS é descrito como fixture estrutural/sintética, sem alegar transcrição literal do XML integral ausente.
- IDs anteriores T144/T145/T146 preservados; as três tarefas visuais acrescentadas na última revisão foram corrigidas para T148/T149/T150. T147 continua sendo o aceite visual daquele recorte. T151 registra esta reconciliação documental.
- Oito tarefas foram reabertas com motivo concreto na própria linha: 001/T108, T140, T145; 002/T013; 003/T022; 004/T016–T017; 005/T012. A existência de implementação parcial/testes próximos não conclui seus escopos integrais.
- 77 links para artefatos/fontes históricos ou propostos que não existem na árvore foram convertidos em referências textuais explícitas, preservando a informação e sem inventar novos caminhos de evidência. Uma âncora histórica ganhou alias para manter as referências antigas válidas.

**Verificação documental:** revisão de caminhos/âncoras locais nos 127 Markdown de `specs` e dos READMEs de origem; unicidade de IDs nos oito backlogs; presença do quadro de estado em todas as 33 specs; revisão de whitespace. Os checklists de qualidade e os registros de aprovação humana foram preservados. Resultado final da checagem: PASS, sem links locais/âncoras inválidos e sem IDs de tarefa duplicados.

**Pendências:** concluir as tarefas reabertas e o restante do backlog nos respectivos donos, depois executar os cenários manuais Windows/DPI/PowerPoint da revisão escolhida. Esta atualização documental não declara conformidade funcional integral.


<a id="24-menu-xml-e-skins-dos-jogos--1002"></a>
## 24. Menu XML e skins dos jogos — 1.0.0.2

**Executor:** agente nesta sessão em Windows amd64, 2026-10-24. Revisão técnica disponível para teste humano; sem aceite presumido.

### Fonte recuperada e comportamento

Cópia original recuperada de `D:/desenvolvimento/ia/estudo/sistema-prototipo-android/componentes/src/java/menu/menu.xml`, preservada em `material-legado/menu/menu.xml`: 30.514 bytes, SHA-256 `f54ad100a67ec1ec947ee74d35123a908e0fc06004813c82f1d541317dee6bd8`. Codificação ISO-8859-1 respeitada, sem transcodificar a fonte nem incluir elementos comentados.

A fábrica preserva 22 itens raiz: dez jogos primeiro e doze serviços depois, conforme ordem XML. A cena separa as categorias nas metades sem alterar o modelo. Jogos: MEGA-SENA, LOTOFÁCIL, QUINA, LOTOMANIA, LOTECA, DUPLA-SENA, TIMEMANIA, DIA DE SORTE, SUPER SETE, MILIONÁRIA. LotoGol comentado e Federal não são acrescentados à raiz. Mega-Sena, Lotofácil, Quina e Dupla-Sena mantêm seus dois filhos normal/especial; RELATÓRIOS usa os 27 títulos descendentes e 148 folhas reais. BOLÕES DIGITAIS803 é último filho de BOLÃO, com VENDER COTAS791, BOLÕES MARKETPLACE790 e BOLÕES CONCLUÍDOS792.

O PNG original inclui ícone e nome: a marca não ganha texto duplicado. Legenda personalizada fica editável em caixa separada. Variante explícita prevalece; item/composto desabilitado usa disabled, inclusive durante pressão. Prévia usa touched enquanto pressionado e restaura default ao soltar/sair. Projeto existente não é reescrito nem recebe novos itens; variante vazia reconhecida só resolve a aparência. Inserir novo Menu aplicação para obter a árvore corrigida.

### Verificações executadas

- `python -X utf8 scripts/import_legacy_menu.py --check`: PASS, dados estáticos conferidos com o XML. Não há leitura de XML no aplicativo.
- `menu_xml_test.go`: PASS, comparação independente de cada caminho, ordem, nome com espaços, códigos, opções, nomes de movimento e atalhos; títulos principais aplicados aos submenus. O modelo não armazena título separado de folha.
- `menu_games_test.go` de presentation: PASS, doze marcas normais disponíveis e quatro especiais ativos, paths literais dos três estados, variantes vazias/explícitas, disabled do composto, override e legenda personalizada.
- `desktop/TestGameMenuPreview`: PASS, recursos realmente renderizados, pressão/liberação/saída, bloqueio disabled e abertura do submenu normal/especial pelo controle de Mega-Sena. Estado confirmado preservado.
- `pptx/TestMenuGamePPTImages`: PASS, comparação byte a byte de todos os recursos na fixture, incluindo Milionária disabled e os acréscimos explícitos LotoGol/Federal somente nesse teste; texto autorado editável e ausência de marca sobreposta.
- Build final `python -X utf8 scripts/build.py --target windows --arch amd64 --check`: exit 0. gofmt, `go vet ./...`, `go test ./...` (13 pacotes com testes), importador de skins (205 chaves/224 recursos, oito indisponíveis) e três testes Python de empacotamento aprovados. Log local: `apps/gerador-jornada/build/games-build.log`.
- EXE final gerou as duas fixtures `.jornada` e os dois PPTs; capturou canvas com novo XML e abriu uma cópia da fixture anterior para verificar skins sem migrar a árvore.

### Artefatos identificados

| Arquivo | SHA-256 |
| --- | --- |
| `dist/windows-amd64/JornadaStudio.exe` | `20f0bbcaf25b81ef572b12ca7a1b9754361a78ea8e0a3c022636c30ba5927a27` |
| `dist/JornadaStudio-1.0.0.2-windows-amd64.zip` | `519d680adef9053f4945b05b7305d5fc023f539c709fccedd923bf295c321ab1` |

Caminhos relativos a `apps/gerador-jornada`. O EXE tem 33.133.568 bytes. Fixtures atualizadas: `dist/estilos-legado.jornada`, `dist/compostos-legado.jornada` e PPTs homônimos. Evidências locais:

- [Canvas nativo, menu atual](../../apps/gerador-jornada/build/game-menu-native.png).
- [Prévia realmente renderizada](../../apps/gerador-jornada/build/game-menu-preview.png), [pressionado](../../apps/gerador-jornada/build/game-menu-preview-pressed.png), [desabilitado](../../apps/gerador-jornada/build/game-menu-preview-disabled.png).
- [Projeto anterior preservado](../../apps/gerador-jornada/build/game-menu-compatibility.png), cuja árvore continua antiga e cujas skins já resolvem.
- [PPT renderizado no LibreOffice](../../apps/gerador-jornada/build/game-menu-ppt.png), através de PDF temporário de três páginas. Imagens dos jogos conferidas visualmente. **A fonte dos textos editáveis foi substituída pelo visualizador**, pois a face Futura não está incorporada ao PPT; esta captura não comprova paridade tipográfica nem homologação no PowerPoint.

### Teste humano pendente

1. Abrir o EXE identificado e `dist/compostos-legado.jornada`, tela 02. Conferir Milionária ativa, ausência de LotoGol e ordem dos jogos/serviços; comparar com o XML preservado.
2. Na Prévia, abrir Mega-Sena e os demais ramos normal/especial, pressionar/soltar os botões e conferir os ícones específicos. Desabilitar um jogo no editor, confirmar e verificar imagem disabled e ausência de ação.
3. Explorar BOLÃO/digital, RELATÓRIOS, PAGA PREMIADOS e MEIOS; conferir nomes, códigos, opções, espaços e voltar/cancelar. Não executar movimentos externos.
4. Renomear uma legenda sem trocar sua variante: manter marca e texto autorado separado. Salvar/reabrir e testar Undo/Redo; abrir cópia de projeto anterior e conferir que a árvore não foi substituída.
5. Abrir o PPT no PowerPoint com as fontes necessárias disponíveis e conferir imagens/textos, navegação e DPI 100/150/200%. Registrar ambiente, capturas e aceite explícito. T147 e demais gates manuais continuam abertos.

<a id="25-submenus-java-e-movimentos--1003"></a>
## 25. Submenus Java e movimentos — 1.0.0.3

**Executor:** agente nesta sessão, Windows amd64, 2026-10-24. Revisão executável para teste manual; não é aceite do usuário nem homologação Office/DPI.

### Fontes e diferença encontrada

O menu.xml tinha sido preservado, mas faltava conferir a pintura de `PnlMenu.java` e a semântica de `UIMenuControle.java`. Ambos foram copiados integralmente, sem transcodificação, para `material-legado/menu`; origem/hashes na seção correspondente de research.md. O Java percorre `getItemMenuXml()` (coleção de itens/filhos), desenha a faixa específica e centraliza as duas opções dos jogos. A fábrica anterior copiava Name para Title ausente, e o renderer usava uma grade genérica no topo. O teste anterior abria o runtime, mas não inspecionava a janela reconstruída.

Corrigidos: ausência literal de tituloPrincipal; faixa verde Mega-Sena (0,150,800,27); botões normal/especial em (220,310,164,52)/(400,310,164,52); jogos root 164×52/pitch77; grade dinâmica Java; composição financeira/digital; legendas com fonte física; rodapé compartilhado com Cancelar/Voltar no submenu e Sair na raiz. O comando financeiro omite Cancelar conforme Java. Mesma geometria de pintura/clique para canvas/prévia/PPT; comandos derivados não acrescentam itens persistidos.

MovementCode/MovementName/MovementOption/Shortcut passam a ser editáveis no formulário e aparecem na seleção simulada. Normal/especial são folhas distintas por UUID e nomeMovimento, mesmo com código 702 compartilhado. Voltar/Cancelar limpam a seleção; componente desabilitado e sessão encerrada bloqueiam ativação. Nenhuma rota ou operação é inferida desses metadados.

**Pendência explícita:** movimento.xml/movimentos.xml não foi recuperado nos materiais, anexos e registros locais acessíveis. UIMenuControle consulta InfoMovimento por nome quando nomeMovimento existe. Referências preservadas não validam esse cadastro; classes/telas/associações ausentes não foram inventadas. T154 continua aberta. Título separado de folha também não foi acrescentado ao modelo atual.

### Verificações executadas

- Regressão `presentation/TestMegaSubmenuLegacyScene` falhou antes da correção com “Mega-Sena sem tituloPrincipal ganhou faixa principal inventada”; passou depois, conferindo faixa, caixas, clique e comandos.
- `desktop/TestGameSubmenuRedraw`: PASS, clique nos controles reais de `window.Content()`, reconstrução da página, quatro duplas normal/especial, metadados distintos, Voltar/Cancelar/Sair e isolamento do documento. A captura mostra o conteúdo já reconstruído, não só Path.
- `desktop/TestMenuMovementFields`: PASS, selecionar Mega-Sena Especial no formulário, editar código com zeros/Unicode/opção/atalho, confirmar cópia independente e conservar a folha normal/origem.
- `pptx/TestMegaSubmenuPPT`: PASS, pacote em memória incorpora bytes originais da faixa, Mega-Sena, Mega da Virada e rodapé; texto dos comandos permanece editável e não aparece título inventado. Não foi gerado PPT de entrega.
- `scripts/check.ps1`: PASS, gofmt, importadores --check, go vet ./..., go test ./... nos 13 pacotes com testes e três testes Python do empacotador. Testes do empacotador usam dados temporários e não criam um pacote de distribuição.
- `scripts/build.ps1 -OutputDirectory dist/windows-amd64`: PASS. Único artefato de compilação entregue: `JornadaStudio.exe`. Nenhum ZIP, instalador ou build de outra plataforma foi gerado. Inventários/pacotes anteriores são históricos.
- Smoke do EXE: uma fixture temporária em build/ recebeu o submenu Mega-Sena como inicial e foi aberta no editor nativo, sem migrar arquivos do usuário. Captura revisada visualmente.

### Identidade e evidências

`apps/gerador-jornada/dist/windows-amd64/JornadaStudio.exe`, software **1.0.0.3**, **33.153.536 bytes**, SHA-256 `6e5d11ead077cc97078a3d9bb3f14f56e54ee7f2a8f9e92cd6587bf3e4959f3d`. A cópia de compatibilidade já existente em `dist/JornadaStudio.exe` também foi atualizada com os mesmos bytes/hash, para evitar abrir a revisão anterior por esse caminho; nenhum novo formato foi produzido.

- [Submenu reconstruído na prévia](../../apps/gerador-jornada/build/submenu-megasena.png).
- [Canvas do EXE Windows 1.0.0.3](../../apps/gerador-jornada/build/submenu-megasena-native.png).
- [Raiz com comandos e altura dos jogos corrigidas](../../apps/gerador-jornada/build/submenu-raiz.png).
- Log local: `apps/gerador-jornada/build/submenus-check-build.log`.

Reprodução dos checks/build, a partir da aplicação:

```powershell
./scripts/check.ps1
./scripts/build.ps1 -OutputDirectory dist/windows-amd64
```

Captura opcional do teste, exclusivamente em build/:

```powershell
. ./scripts/env.ps1
$env:JORNADA_SUBMENU_CAPTURE = "$PWD/build"
go test ./internal/adapter/desktop -run TestGameSubmenuRedraw -count=1
Remove-Item Env:JORNADA_SUBMENU_CAPTURE
```

### Teste humano pendente

1. Abrir o EXE identificado, inserir novo Menu aplicação e iniciar Prévia Tela Atual. Clicar Mega-Sena: conferir faixa verde, dois botões centralizados, Mega da Virada como imagem especial e Cancelar/Voltar no rodapé. Repetir Lotofácil, Quina e Dupla-Sena.
2. Selecionar normal/especial: conferir nomesMovimento distintos apesar do código compartilhado; nenhuma operação externa deve iniciar. Testar Voltar, Cancelar e Sair e verificar que somente a prévia fecha.
3. Em Propriedades, selecionar uma folha e editar código/nome/opção/atalho. Descartar/confirmar, Undo/Redo e salvar/reabrir devem preservar o estado esperado. Conferir financeiro/digital e múltiplas instâncias em DPI 100/150/200%.
4. Arquivos existentes mantêm títulos/variantes autorados. Para comparar defaults corrigidos, usar uma nova instância; se um menu antigo contém título automático herdado igual ao nome, limpar explicitamente seu Título no formulário para reproduzir a ausência no XML. Fixtures/pacotes de dist anteriores não foram regenerados.
5. Conferir o cadastro de movimentos somente após recuperar sua fonte real; T154 não pode ser concluída pela contagem de itens do menu. Office e aceite humano continuam em T147 e nos gates anteriores.

## 26. Launcher e estrutura de projeto — 1.0.0.4

Feature [034](../034-launcher-estrutura-projeto/spec.md) implementada tecnicamente. Base Git `2319beb05 + árvore de trabalho local (sem commit de entrega)`. Writer schema **3**/catalog **4**, leitores 1/1, 2/1, 3/2, 3/3 e 3/4. Não confundir com a versão do produto.

C1 (launcher/recentes/branco/exemplo), C2 (workspace/tree/criar/importar) e C3 (sugestão/cancelamento/vínculo/ida e volta) passaram em testes com arquivos reais temporários e controles Fyne. A captura detectou uma raiz de Tree não renderizada; a correção foi verificada por labels reais e nova captura. Listagem de 100 jornadas com fundo: **54,96 ms** no ambiente atual. `scripts/check.ps1` e `scripts/build.ps1 -OutputDirectory dist/windows-amd64` passaram após a correção final; startup nativo sem `-open` e com `-open` passou.

EXE **33334784 bytes**, SHA-256 `366f6a7276d781bb7193d15fde70dbd12f0a0f8ec7837744495e2655f47a9378`. Entrega apenas `JornadaStudio.exe`; cópia de compatibilidade idêntica. Capturas: [launcher nativo](../../apps/gerador-jornada/build/034-launcher-native.png), [árvore nativa](../../apps/gerador-jornada/build/034-workspace-native.png), [prévia C3 no backend de teste](../../apps/gerador-jornada/build/034-preview-test.png). Logs, testes e restrições detalhados em [034/quickstart — resultado técnico](../034-launcher-estrutura-projeto/quickstart.md#resultado-técnico--2026-10-24--1004).

O README explica as alterações necessárias, a preferência local de recentes, a pasta como workspace, a confirmação que salva a origem, a migração com backup e a preservação do arquivo ao desfazer vínculo. Homologação manual/DPI/Office/aceite continuam separados; 001/T156 permanece aberto. A ausência de `movimento.xml`/`movimentos.xml` continua em T154, sem atribuir sua recuperação a esta entrega.

## 27. Correção da transição launcher → editor — 1.0.0.5

Relato do usuário reproduzido: o `.jornada` era criado, mas o processo fechava ao criar/abrir pelos recentes. Causa demonstrada no driver Fyne/Windows: SetTheme durante New do editor tentava acessar HWND de janela ainda não exibida. A versão 1.0.0.4 testou startup e gravação, mas deixou essa passagem nativa descoberta. Corrigidos tema antes de NewWindow, inicialização do seletor sem callback e contraste inicial do título; não se alterou o formato 3/4.

Regressão unitária, oito fluxos completos no driver nativo e uma reabertura final passaram. Check completo, build somente EXE e startup final aprovados. EXE 1.0.0.5: **33336320 bytes**, SHA-256 `2576917419e883744ff411ee362506935f808f6f9a01eb83242b70b563845568`. Logs, capturas e comandos no [quickstart 034](../034-launcher-estrutura-projeto/quickstart.md#correção-de-encerramento--1005). A documentação anterior registra o estado histórico da 1.0.0.4; aceite manual continua separado.

## 28. Inicialização separada e splash nativa — 1.0.0.9

Modelos e implementação separados: [035 domínio Inicialização](../035-inicializacao/spec.md) / internal/startup; [036 apresentação e adaptadores](../036-splash-screen/spec.md) / desktop, builtinassets e nativevideo. O MP4 v2 original com áudio é incorporado; apresentação fullscreen somente vídeo; autoria movida para o launcher. Recursos e reprodução completam em qualquer ordem, mas ambos são obrigatórios. Wire preservado em 3/4.

Corrigidas as falhas observadas na revisão intermediária: escala física/DPI, razão 16/9 interpretada como dimensão e superfície de vídeo sem protocolo de pintura. A revisão final usa janela nativa própria com WM_PAINT/WM_SIZE e Play após mídia pronta. Quatro inicializações do EXE confirmaram imagem no início e no fim, trilha de áudio ativa e transição posterior. Criar/reabrir/carga lenta de 14 s/cancelamento/falha e -open passaram na integração. Check completo e build somente EXE aprovados. [Evidências, capturas e roteiro](../036-splash-screen/quickstart.md).

EXE **1.0.0.9**, **37851648 bytes**, SHA-256 `af4400c379ad87096a693f30ee97275fa565f6ba0fa1d7196ae4ef388f4aaa03`; cópia de compatibilidade idêntica. Aprovação humana continua em T159/036-T009. Capturas de vídeo nativo exigem a janela real; um Canvas.Capture GL ou evento PLAY não comprova imagem no player.

## 29. Primeira apresentação e fullscreen — 1.0.0.10

Integração 035/T005, 036/T011–T013 e 034/T024: vídeo em fullscreen na maior área proporcional (sem limite de pixels originais), preferência por usuário “Mostrar splash screen ao iniciar”, padrão desmarcado após primeira apresentação completa, marcar/desmarcar persistidos fora do projeto. Carga de recursos continua obrigatória sem vídeo. [Modelo, revisão, roteiro e evidências atuais](../036-splash-screen/quickstart.md#revisão-atual--10010).

Checks completos e build passaram. Quatro reproduções reais com captura inicial/final e duas aberturas sem vídeo passaram; -open sem splash abriu o editor preservando o projeto. Testes Go cobrem política pura, checkbox real, persistência e falhas. EXE Windows amd64: 37860864 bytes, SHA-256 `c71eb5d0a0c342e2576e0d748214ded693a7831bd51bfc819c7f979a3bedb7ad`, cópia de compatibilidade idêntica. Nenhum pacote ou outro formato gerado. Aprovação humana permanece em T159/036-T009/034-T025, sem encerrar pendências antigas.

## 30. Vídeo em janela normal — 1.0.0.11

Por solicitação posterior do usuário, remover fullscreen da splash e reproduzir em janela normal centralizada, 960×540 unidades Fyne inicialmente. Não alterar áudio, política de primeira exibição/repetição ou domínio 035. Detalhes e limites em [036/quickstart](../036-splash-screen/quickstart.md#revisão-vigente--10011-em-janela).

Checks completos/build aprovados; duas reproduções do ciclo nativo com capturas inicial/final em janela normal e dois caminhos sem vídeo aprovados. Duas tentativas extras foram canceladas por fechamento e não contam como repetições aprovadas. Imagem no computador do usuário permanece pendente. EXE Windows amd64 37860864 bytes, SHA-256 `74aa094afa2e98a6e8d5f213ac735307b4ada6c8591b0c4c78bbeb2bff506c85`, cópia de compatibilidade idêntica; nenhum pacote ou outro formato gerado. Manter T159/036-T009 abertos.

## 31. Jornada e PowerPoint na pasta do workspace — 1.0.0.15

Requisito 001/FR-062, tarefa T162. A exportação mantém a tela de opções, com ação Exportar, e gera automaticamente `<nome-do-arquivo-atual>.pptx` na raiz do workspace. Não há seletor de pasta/arquivo; se já existe apresentação, pede apenas confirmação para substituí-la. Sem workspace, informa que é necessário definir/salvar o projeto. Salvar uma jornada nova com workspace configurado usa diretamente essa pasta e preserva homônimos. Os formatos deste incremento são `.jornada` e `.pptx`; não há `.ppt` legado.

Estado: **pronto_para_teste_manual**. Testes focados aprovados: `TestWorkspacePowerPointOutput` gerou um PPTX ZIP real ao lado do arquivo atual, sem diálogo de destino, e preservou a saída ao cancelar substituição; `TestFirstSaveUsesConfiguredWorkspace` gravou a jornada no primeiro save sem novo seletor. `scripts/check.ps1` passou integralmente; build PE AMD64 e smoke nativo com PATH sem Git aprovados em 2026-10-25. Captura em `apps/gerador-jornada/build/001-workspace-1.0.0.15.png`. EXE: `apps/gerador-jornada/dist/windows-amd64-1.0.0.15/JornadaStudio.exe`, SHA-256 `a685221d45dbe561fdd73e49edef5ad49441ca6f01bb3c27307bcaea4ac468b3`. Configuração OAuth externa incorporada e gerado Go removido. Não houve mudança de renderer nem homologação do conteúdo no PowerPoint desktop; aprovação humana pendente.

Roteiro: abrir jornada salva, usar PowerPoint → Exportar e conferir `.pptx` na mesma pasta; repetir e recusar substituição para conferir preservação; aceitar para atualizar. Sem projeto salvo/workspace definido, esperar orientação para definir a pasta. Reteste humano pendente.
