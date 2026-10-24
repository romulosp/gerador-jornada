# Research: Inicialização

Separar regra de liberação da tecnologia evita que término de timer, abertura de janela ou carregamento de uma imagem seja interpretado como fim de reprodução. Usar eventos explícitos e uma autorização consumível é suficiente; máquina de estados com framework, relógio no domínio ou novas goroutines seriam desnecessários. O modelo é transitório e não altera os cinco domínios de autoria nem o wire 3/4. A mídia produz eventos reais no adaptador 036; uma interface de arquivos ou UI no domínio seria dependência invertida indevida.

## Primeira exibição e repetição optativa

Separar Preference da Session evita misturar uma política durável por usuário com eventos transitórios de mídia. Ausência equivale a Presented=false/ShowOnStartup=false; conclusão não habilita repetição. Preferência não é dependente da versão, do projeto ou da atualização do EXE. Cancelamento/falha não conta como apresentação completa. O adaptador confirma a conclusão antes de montar o launcher, preservando qualquer escolha já gravada.
