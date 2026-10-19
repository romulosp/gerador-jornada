# Contract: geração Java 1.6/Swing com EFL

**Version**: 1 | **Requirements**: FR-022 a FR-025, NFR-003, NFR-008; US5.
**Scope**: fontes e recursos para integração no sistema EFL; sem JAR, EXE Java ou compilador no editor.

## Entrada e pré-condição

Receber projeto validado, packageName e profileId. packageName deve conter segmentos ASCII Java válidos em minúsculas, sem palavras reservadas; profileId deve selecionar perfil EFL embutido e homologado.

O perfil não é template arbitrário de projeto. Ele pertence à versão do editor e documenta bibliotecas, assinaturas, constantes de skin, construção, contexto UI, listeners e ciclo de vida do hospedeiro.

**Estado atual**: não existe perfil homologado porque não foram fornecidos JARs/fontes da EFL. A implementação pode desenvolver testes e templates, mas a saída Java publicável deve retornar `JAVA_PROFILE_UNVERIFIED` enquanto esse gate não for cumprido. Não gerar silenciosamente Swing puro nem apresentar uma saída parcial como projeto compilável.

## Saída prevista

Diretório novo contendo:

```text
src/<packagePath>/
├── screen/         # <fileIdDaTela>.java por tela
├── component/      # fábricas por tela e arquivos individuais em <fileIdDaTelaEmMinusculas>/
├── navigation/     # controlador e contrato de navegação
└── integration/    # adaptação ao contexto real do hospedeiro
resources/          # recursos referenciados pelo projeto
sources.list        # caminhos relativos dos arquivos .java
README.md           # dependências, build, integração e limites
dependencies.md     # versões EFL/Java e bibliotecas requeridas
```

As classes de tela montam o painel e delegam construção à fábrica; a fábrica cria/configura componentes em métodos coesos. Navegação resolve IDs estáveis e destinos fora da criação visual. Listeners delegam ao controlador, sem TODO em uma capacidade solicitada.

Não gerar uma classe única para todas as telas. Não adicionar persistência ou regras de negócio. O roteiro demonstrativo pertence ao PPT/prévia e não gera automação de cursor ou digitação no aplicativo Java.

A saída integra o contexto do sistema EFL: não inventar `UIControle`, `uiManager` ou uma janela executável com construtores adivinhados. A camada integration e as instruções precisam corresponder ao bootstrap real validado no perfil.

## Matriz EFL a homologar

| Família | Evidência observada no template legado | Evidência real necessária |
| --- | --- | --- |
| button | EFLButton, ButtonSkin, SwingWorkerActionListener | Construtores, contexto UI, constantes e thread do callback |
| label/title | EFLLabel e LabelSkin | Construtores, skins, fonte, cor e alinhamento |
| textField | EFLTextField, TextFieldSkin, documentos EFL | Construtores, setDocument, leitura/escrita e formatos |
| radio | EFLButtonGroup declarado; template inconsistente | API correta de grupo/opções/seleção/orientação |
| report | EFLReportEvent, listener, painéis de cabeçalho/linha | Criação e povoamento sem apagar cabeçalho indevidamente |
| numericKeyboard | EFLVirtualKeyboardNumeric | Construção e vínculo efetivo com campo da mesma tela |
| screen | EFLPanel e PanelSkin | Layout absoluto, contexto, ciclo de vida e inclusão no hospedeiro |

Esses nomes não garantem assinaturas. A equipe mantenedora fornece JARs, versões, dependências transitivas e um exemplo funcional antes da conclusão de T005/T006 e da homologação de US5.

## Identificadores e texto

Preservar nomes visíveis como texto. Usar FileID confirmado exatamente, inclusive caixa, como classe/nome-base: `screen/PnlIdentificacao.java` e `component/pnlidentificacao/BtnContinuar.java`. A fábrica individual do componente cria sua instância; a fábrica da tela delega a ela. Arquivos de suporte derivados, como `<IDDaTela>Factory.java`, não alteram o nome-base principal confirmado.

O valor deve obedecer ao contrato de FileID: ASCII alfanumérico, início em letra, até 40 caracteres sem extensão, Pnl obrigatório somente nas telas; recusar nomes Windows reservados. Palavras reservadas Java e colisões do escopo gerado produzem diagnóstico; não normalizar/renomear silenciosamente IDs confirmados durante exportação. UUIDs internos continuam resolvendo navegação e vínculos; Name aparece somente como conteúdo/comentário. IDs de componente podem se repetir em telas diferentes porque seus arquivos usam contextos de pacote separados.

Escapar aspas, barras, controles e quebras em literais Java. Preservar Unicode no arquivo UTF-8 e no build com encoding explícito. Texto que se parece com escape Unicode, HTML ou `*/` precisa de tratamento próprio ao entrar em Javadoc, sem encerrar comentário ou alterar sintaxe.

Gerar imports únicos em ordem lexical. Usar somente sintaxe/API SE 6: sem diamond, lambdas, streams, try-with-resources, var, records ou switch moderno. Bibliotecas EFL também devem ser compatíveis com o ambiente de destino.

## Threading e compatibilidade

O ambiente informado pelo usuário contém JDK 17 em `C:\Desenvolvimento\jdk-17.0.11`. A consulta local a `bin/javac.exe --help` confirmou suporte a `--release`, `--source` e `--target` somente de 7 a 17. Portanto, esse compilador não atende ao alvo 1.6. O alvo dos fontes permanece Java SE 6; não elevar a versão da saída em função do JDK instalado.

Para validação estática, a ferramenta externa `apps/gerador-jornada/scripts/java6.py` executa ECJ 3.33.0 no JDK 17 informado, com `-source 1.6 -target 1.6 -encoding UTF-8 -proc:none`. Exige classfile major 50 na saída própria e recusa bytecode posterior ao major 50 nas dependências. Animal Sniffer 1.24 e a assinatura `org.codehaus.mojo.signature:java16:1.1` verificam classes, métodos e campos da plataforma, inclusive chamadas dentro dos JARs fornecidos. Não basta emitir bytecode antigo usando APIs recentes.

Versões, URLs e SHA-256 estão fixados em `scripts/java6-tools.json`; ferramentas são preparadas em `.tools/java6`, sem modificar a instalação Java. Depois da preparação, `--offline` impede downloads. Informar todos os JARs transitivos com opções `--classpath` repetidas. O relatório registra hashes dos fontes, dependências e ferramentas; falhas também substituem o relatório com status reprovado e encerram com código diferente de zero.

Essa ferramenta pertence ao ambiente de desenvolvimento/homologação, sem ser incorporada ao EXE. A validação estática foi exercitada com fontes do gerador usando um perfil estrutural de teste; isso não homologa EFL. Execução na JVM Java 6 e no hospedeiro EFL real continua necessária para aprovação de US5. A disponibilidade do JDK 17 não fornece os JARs EFL nem encerra essa pendência.

Toda criação e mutação de Swing deve ocorrer na Event Dispatch Thread. Um listener EFL denominado background não é evidência de execução na EDT: o adaptador deve despachar a mudança de tela para ela explicitamente quando necessário.

Bounds são aplicados em unidades inteiras; não usar screenshots nem deslocamentos fixos do cabeçalho legado. Estilo alterado só pode ser exportado se mapeado pelo perfil; caso contrário, erro `JAVA_PROPERTY_UNSUPPORTED`.

Homologar com compilação usando API Java SE 6 e target 1.6, verificar classfile major 50 e executar no ambiente Java/EFL real. Compilar com um JDK moderno sem target 6 ou apenas contra dublês não comprova o requisito.

## Documentação e build entregue

Javadoc de pacotes, classes, interfaces, construtores e métodos; comentários de campos, constantes e variáveis, incluindo unidades e contexto dos listeners. Não inserir data corrente nos fontes nem comentários genéricos que escondam integração pendente.

README gerado deve informar comandos compatíveis com a toolchain homologada, classpath externo, encoding UTF-8, sources.list e ponto de integração no hospedeiro. Não distribuir JARs proprietários no pacote de fontes; inventariar versões necessárias.

## Diagnósticos e testes

`JAVA_PACKAGE_INVALID`, `JAVA_IDENTIFIER_INVALID`, `JAVA_NAME_COLLISION`, `JAVA_PROFILE_UNVERIFIED`, `JAVA_PROPERTY_UNSUPPORTED` e `JAVA_RESOURCE_INVALID`.

Testar todas as famílias, acentos/aspas/barras/quebras, colisões e nomes reservados; comparar fontes de referência; compilar com bibliotecas reais; executar navegação A → B → A e digitação/seleção/teclado. Registrar perfil, versão dos JARs e toolchain na evidência de homologação.

Referência pública: [Swing Java SE 6](https://docs.oracle.com/javase/6/docs/api/javax/swing/package-summary.html). As APIs EFL são dependência externa ainda não comprovada.
