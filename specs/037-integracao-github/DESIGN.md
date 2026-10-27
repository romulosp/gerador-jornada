# Design: GitHub no workspace

Interface administrativa implementada para a revisão de código 1.0.0.12. Reutiliza tema, tipografia e controles Fyne da aplicação; não altera componentes EFL, suas skins ou a cena das jornadas. Testes do painel/callback não representam homologação visual em Windows/DPI reais; o resultado manual pertence ao quickstart.

## Entrada compacta

Na aba Estrutura projeto, a seção GitHub mantém uma linha de estado truncada por reticências e menu ⋮, seguida de uma ação principal. O texto completo fica em **⋮ → Detalhes da conexão**, sem ocupar a altura da árvore.

- Sem conta: **Conectar GitHub**, com ícone Octicons vetorial incorporado e licença MIT em `desktop/assets/github-icon-LICENSE.txt`.
- Conta conectada sem vínculo: a ação principal passa a **Selecionar repositório**.
- Conta e workspace vinculados: o botão principal é **Publicar no GitHub**; não são empilhados todos os botões administrativos no painel.
- Sem pasta definida: explicar que é necessário abrir/salvar uma jornada e desabilitar a ação correspondente. Sem cadastro OAuth no build: informar configuração indisponível, mantendo os fluxos locais.

O menu ⋮ oferece Detalhes da conexão, Conectar ou trocar conta, Selecionar repositório, Publicar, Atualizar estado, Ver repositório no navegador e Desconectar. A elegibilidade acompanha o estado da conta/vínculo/operação. Logout informa que arquivos, histórico e vínculo local permanecem.

**Justificativa do ajuste**: o desenho inicial de vários botões/linhas consumia o espaço da árvore em painel estreito. Status compacto, detalhes explícitos e uma ação principal preservam acesso às mesmas operações sem esconder as jornadas. O ícone acompanha o contraste do tema sem modificar a aparência do documento.

## Autorização e operação

Conectar abre navegador padrão e diálogo cancelável de espera. Sucesso retorna à mesma janela e abre seletor. Operações de rede são executadas em background; um diálogo modal impede alterações concorrentes na sessão enquanto o trabalho está ativo. Cancelamento/falha encerra a tentativa; transação local interrompida é recuperada ou recebe diagnóstico antes de liberar novas escritas.

## Seletor de repositórios

Diálogo inicialmente 720×520 unidades Fyne, com busca, lista rolável e informação de destino/branch. Cada linha mostra proprietário/nome, visibilidade e escrita. Arquivado, desabilitado ou sem escrita não permite confirmar; o impedimento é mostrado em texto. Atualizar lista refaz a consulta e Cancelar não cria vínculo.

**Paginação efetiva**: o cliente REST percorre todas as páginas em background e somente depois apresenta a lista completa para busca global. Não há botão Carregar mais nem lista parcial incremental visível. Essa adaptação preserva FR-003: todos os destinos efetivamente autorizados continuam encontráveis, e uma falha de página não é apresentada como lista completa. Durante a carga, o diálogo de operação permanece cancelável; erro permite nova tentativa.

O seletor mostra o destino completo e usa diretamente o item selecionado ao clicar Conectar repositório; não abre outro formulário nem pede para redigitar owner/name. Alterar o filtro limpa seleção e detalhes antigos. Após conectar, abre a revisão da publicação; o usuário confirma o envio em Publicar no GitHub. Somente `.jornada` da raiz será integrado no disco, sem copiar código-fonte ou pastas remotas. Colisões entre jornadas são bloqueadas com diagnóstico que identifica o arquivo e preserva os dados. A jornada ativa poderá ser recarregada somente com pacote válido e sessão limpa. Novo projeto usa o campo Pasta do projeto como destino final, sem acrescentar outra pasta homônima. Mudança da 1.0.0.14 substitui a digitação extra rejeitada no teste de 1.0.0.13.

## Revisão e publicação

O comando principal Publicar exige aplicar/descartar rascunhos de propriedades e salvar a edição pendente. A revisão apresenta pasta, destino/branch, autoria pública, inclusões/alterações/exclusões de todas as jornadas elegíveis e mensagem editável não vazia; não oferece subconjunto parcial. Exclusões listam referências locais afetadas, sem reparo inferido.

A revisão também informa que arquivos locais podem ser atualizados pela integração remota e que a permissão OAuth pode abranger outros repositórios, embora a operação use só o destino exibido. **Publicar e atualizar local** confirma; Cancelar não publica.

Com commit pendente, a revisão informa que ele será reenviado/verificado primeiro e que novas edições aguardam outra publicação. Erros de envio preservam a pendência; nova ação Publicar refaz a verificação. A data de última publicação e a contagem de alterações locais ficam nos detalhes. Atualizar estado consulta o remoto sem confirmar commit/push; não transforma informação antiga em promessa permanente de sincronismo.

Depois de publicar, a confirmação identifica repositório, branch, OID e Ver commit no GitHub. Atualização sem novo commit recebe mensagem própria. Vínculo incompatível identifica os dois destinos e orienta usar outra pasta, preservando o histórico anterior. Revisão motivada pelo teste manual reprovado de 2026-10-25; não altera skins, cena, prévia ou exportadores.

## Acessibilidade e validação pendente

Estados têm rótulos textuais além de cor e mensagens em português; busca, lista, mensagem e ações devem permanecer alcançáveis por teclado. Teste de layout confere que o status não esconde a árvore, mas foco nativo, DPI 100/125/150%, temas claro/escuro e nomes longos ainda precisam da matriz visual manual.

Capturar entrada, seletor, revisão, pendência e conflito na build entregue e comparar editor/prévia/exportações. Nenhuma conclusão de homologação visual ou aprovação humana é derivada apenas de testes Fyne em memória.
