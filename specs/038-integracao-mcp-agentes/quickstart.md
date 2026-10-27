# Validação e aceite — SPEC 038 v0.4

**Estado de entrega:** em_andamento — análise após o segundo replanejamento concluída sem novos achados na conversa. Usuário invocou speckit-implement novamente após o gate de checklist; implementação autorizada, mantendo os sete marcadores pendentes intactos. Início pela prova G0, sem homologação ou aprovação funcional.
**Revisão executável da 038:** 1.0.0.18, autoria transitória parcial T015–T025 no EXE Windows padrão; sem salvamento MCP/revisão congelada/biblioteca. 1.0.0.17 preservada como fundação e 1.0.0.16 como G0 histórica. **Perfil:** stdio/MCP 2025-11-25 e HostApproval-v1 desenhados; host real ainda não qualificado. **Execução observada:** G0 e testes técnicos da fundação; nenhum aceite funcional/humano.
**Constituição:** emenda 1.4.0 aprovada/aplicada em 2026-10-25; implementação da spec v0.4 autorizada nesta sequência por invocação reiterada de speckit-implement. Não autoriza commit, merge, push ou aceite funcional automático.

## Início da implementação — baseline T001

- Feature lógica 038-integracao-mcp-agentes; branch Git main preservada, HEAD 16bebff76e97ecba8b02f137c37fd402a9fbf5b2, com alterações anteriores de 037 e documentação preservadas. SHA-256 agregado inicial dos caminhos/bytes versionados e não ignorados sob apps/gerador-jornada: cc19e55a26b719487a37f1f7e6aa69c51fe9e7927c86a03fc69c9cfdd6c2be26.
- Autorização: nova invocação de speckit-implement após apresentação do gate 26/33, sete pendências. Os marcadores de checklists/requirements.md são somente leitura; nenhuma pendência de software será marcada como concluída por autorização.
- Toolchain de desenvolvimento: scripts/env.ps1 seleciona Go 1.27.2 Windows/amd64 e MinGW64/CGO=1. Go global é Windows/386 e não será usado para compilar o aplicativo.
- Ordem: T001–T004 preparam SDK/fixtures; T005 deve provar stdio no EXE Windows padrão sem Fyne/splash/preferências/efeito persistente antes da extração de workspace. Depois I1 autoria existente, I2 biblioteca comum, I3 evidências e análise com captura vinculada, respeitando o DAG. Abertura manual permanece independente; ambos os caminhos conservam estilo legado e os cinco domínios.
- Gates ainda abertos após G0: integração distribuída completa de autoria, limites/deadlines das futuras ferramentas, host humano real, proteção de substituição contra escritor externo, paridade visual/DPI/PowerPoint e aceite do usuário. Startup mínimo e EOF de G0 foram comprovados abaixo. Não reutilizar hashes/resultados de executável anterior como evidência da 038.

G0 já possui flag MCP, harness e build próprios conforme a evidência abaixo. Os roteiros posteriores são expectativas futuras até execução específica; não utilizar G0 ou executáveis anteriores como prova de autoria, biblioteca, consentimento ou paridade visual completos.

## Evidência técnica G0 — T002–T005

- Artefato: apps/gerador-jornada/dist/windows-amd64-038-g0/JornadaStudio.exe; versão 1.0.0.16, PE Windows/amd64 com subsistema windowsgui pelo scripts/build.ps1 padrão. SHA-256: 0ebb1ece6e7e2cb1918068b0fc3dd27736c723b84d2e5553e4c074e99efa7ab0. Pacote local: dist/JornadaStudio-1.0.0.16-windows-amd64.zip, com licenças e inventário; nenhum envio/publicação GitHub.
- Cliente identificado: harness jornada-g0-harness versão 1 usando SDK Go MCP v1.8.0/CommandTransport, processo externo e pipes reais. Perfil servidor restrito a 2025-11-25. Somente aplicacao_capacidades oferecida, startupProbe=true e persistentWrites=false.
- Antes da correção, TestDistributedStartupProbe falhou no EXE 1.0.0.15: flag provided but not defined: -mcp. Após implementação e build: PASS, negociação/listTools/callTool/EOF, ausência de janelas de topo Win32 do PID e de arquivos no workspace/configuração isolados. Encerramento com sucesso sem processo residual. Isso é prova técnica, não host humano homologado.
- Comandos locais: carregar scripts/env.ps1; go test ./internal/application -run '^TestMCPFixture$' -count=1; scripts/build.ps1 -OutputDirectory dist/windows-amd64-038-g0; definir JORNADA_MCP_PROBE_EXE com caminho absoluto do EXE; go test ./internal/adapter/mcp -run '^TestDistributedStartupProbe$' -count=1 -v. Teste sem a variável pula o cenário e não comprova G0.
- Verificações: fixture determinística PASS; 12 fixtures dos quatro schemas validadas documentalmente; coletor de licença preservou o LICENSE original misto do SDK; go vet dos pacotes mcp/application/cmd PASS. Casos wire/semânticos da fixture ainda dependem da implementação posterior. Nenhum resultado observado foi copiado da matriz Q038.
- Auditoria de startup: a ramificação MCP retorna antes de desktop.InitApplicationLog, app.NewWithID e LoadStartupResources. Não foram encontrados inicializadores desktop que chamem esses pontos; a consulta de janelas/configuração conferiu o processo real. Abertura manual, visual/DPI/Office, cliente de produção e consentimento humano continuam com seus testes próprios pendentes.

## Pré-condições para execução futura

- Emenda e desenvolvimento autorizados; análise prévia concluída. Executar cada roteiro somente após sua implementação e identificar a revisão real, sem equiparar autorização a aceite.
- Identificar revisão/hash do EXE distribuído, Windows x64, cliente e versão, modo de transporte e mecanismo de consentimento aprovado.
- Preparar workspace e cópias de jornadas de teste, sem dados pessoais; autorizar materiais do legado separadamente.
- Disponibilizar fixture com duas telas, botão/campo, cabeçalho composto, vínculos internos/externos e recursos; biblioteca vazia e versões válida/incompatível/com recurso ausente.
- Implementar limites/políticas de [mcp.md](contracts/mcp.md) e [consentimento](contracts/consentimento.md); fornecer mecanismo de teste de host sem confirmação e evidência forjada, sem alterar arquivos reais do usuário.
- Para visual: manter mesmos dados e parâmetros nos dois canais, skins/fontes instaladas, DPI 100/150/200% e PowerPoint/ambiente aplicável identificado. Não exigir Java/EFL suspenso.

## Comandos de referência — somente após implementação autorizada

Diretório de trabalho de desenvolvimento: `apps/gerador-jornada`. Primeiro executar os checks existentes pertinentes, depois o harness e só então o EXE padrão distribuído. Builds locais são autorizados na etapa atual de implementação; os comandos futuros abaixo dependem de suas tarefas, sem presumir execução.

```powershell
# Checks existentes; execução futura, não realizada nesta fase.
powershell -ExecutionPolicy Bypass -File scripts/check.ps1

# Pacotes e teste de integração serão criados pelas tarefas de implementação.
go test ./internal/adapter/mcp ./internal/adapter/consent ./internal/adapter/componentstore ./internal/adapter/workspaceguard
go test ./internal/adapter/mcp -run TestDistributedStdio -args -jornada-exe '<caminho-do-EXE-da-revisao>'
```

`TestDistributedStdio` e a flag `-jornada-exe` são contratos de harness propostos, ainda inexistentes. O teste deve iniciar o EXE via `exec.Command`/cliente MCP oficial com pipes, não chamar funções no mesmo processo e não usar binário de console como substituto. Verificar initialize/descoberta conforme o perfil 2025-11-25, tools/resources, cancelamento, limite de frame, EOF, stderr e encerramento.

Configuração de cliente de referência **futura**:

```json
{
  "command": "<caminho-absoluto-do-JornadaStudio.exe-da-revisao>",
  "args": ["--mcp", "--workspace", "<workspace-de-teste-autorizado>"],
  "transport": "stdio"
}
```

Esses campos descrevem a configuração lógica; cada cliente tem formato próprio. Sem approval-profile qualificado, esperar `persistentWrites=false`, R0/R1 disponíveis e R2 recusada. Depois de qualificar host real, operador configura chave pública e escopos fora das ferramentas do modelo e acrescenta `--approval-profile <configuração-do-operador>`; não gerar autorização na própria conversa do agente.

O harness usa chave efêmera e decisões de teste, não representa humano. Teste manual deve identificar host/versão reais, mostrar a revisão original ao usuário, provar negação/aceitação e proteção da chave/emissor. Não presumir que Codex, ChatGPT ou cliente genérico implementa HostApproval-v1. Se a qualificação falhar, registrar bloqueio de R2 e não declarar o MVP completo de salvar aprovado.

Modo manual é iniciado abrindo o mesmo EXE sem `--mcp`; confirmar janela/editor/biblioteca sem rede ou LLM. Trocar o modo não muda o estilo legado, as regras ou o formato da jornada. Para substituição, executar cenário de writer externo antes de habilitar capacidade; ausência de prova oferece novo destino exclusivo, nunca overwrite inseguro.

## Cenários planejados — nenhum executado

| ID | Requisito/cenário | Ação futura e resultado esperado | Resultado observado |
| --- | --- | --- | --- |
| Q038-01 | FR-011/FR-034, US2.5 | Revisar R2 e aprovar por canal confiável; gravar só conteúdo/destino autorizados sem abrir GUI | Não executado |
| Q038-02 | FR-011, US2.5 | Negar/cancelar ou pedir correção; nenhum efeito durável e revisão antiga não autoriza novo conteúdo | Não executado |
| Q038-03 | FR-011, US2.6 | Enviar `approved: true`/texto de aprovação pelo modelo; recusar escrita e preservar bytes anteriores | Não executado |
| Q038-04 | FR-011, US2.6 | Fornecer token ou arquivo gravável pelo agente sem origem independente; não reconhecer consentimento | Não executado |
| Q038-05 | FR-011/FR-012, US2.1 | Mudar conteúdo, base ou destino após revisão/aprovação; invalidar e exigir nova revisão | Não executado |
| Q038-06 | FR-011, US2.7 | Tentar autorização expirada, revogada, de outra sessão/escopo; recusar sem gravação | Não executado |
| Q038-07 | FR-011/FR-016, US2.4/US2.7 | Repetir resultado perdido com mesma solicitação; um efeito; reutilizar autorização em nova solicitação não produz outro efeito | Não executado |
| Q038-08 | FR-034, US2.8 | Usar host sem confirmação confiável/canal alternativo; R2 indisponível, R0/R1 autorizadas utilizáveis, sem fallback de GUI | Não executado |
| Q038-09 | FR-012/FR-013, US2.2/US2.3 | Salvar GUI→MCP e MCP→GUI e interferir com operação Git; detectar conflitos e preservar ambos os trabalhos | Não executado |
| Q038-10 | FR-017/FR-035, SC-008/SC-010 | Cancelar, desconectar/EOF ou falhar durante publicação; preservar última gravação e informar resultado/recuperação sem processo residual | Não executado |
| Q038-11 | FR-027–FR-029, US5.1/US5.6 | Consultar legado autorizado e apresentar hipótese; Go separa evidência/validade de equivalência e autorização, sem executar Java | Não executado |
| Q038-12 | FR-018/FR-030, US5.2–US5.5 | Propor cabeçalho, duas legendas e ComboBox; publicar apenas declarativo suportado/autorizado e identificar desenvolvimento pendente | Não executado |
| Q038-13 | FR-020/FR-025, US3.6/US4.5 | Exercitar estados vazios, incompatíveis, ausentes, conflito e falhas de importação/publicação; ações corretas e dados preservados | Não executado |
| Q038-14 | FR-020/FR-030, US3.7 | Criar seleção→rascunho→validar→revisar→cancelar/confirmar e nova versão/retirada; nenhuma publicação acidental ou alteração antiga | Não executado |
| Q038-15 | FR-002/FR-007/FR-022, US3.1–US3.5 | Criar manual e instanciar MCP, depois inverso; novas identidades, vínculos corretos e equivalência de conteúdo e recursos | Não executado |
| Q038-16 | FR-024/FR-032, SC-005/SC-009 | Abrir sem biblioteca de origem e comparar default legado/overrides em criação/canvas/prévia/PPT, fontes/alpha/partes/DPI e limitações reais | Não executado |
| Q038-17 | FR-001, US1.6/US1.7, SC-012 | Abrir software manual sem LLM/MCP/rede; usar biblioteca por teclado e confirmar que início/fim do MCP não elimina modo manual | Não executado |
| Q038-18 | FR-002, SC-013 | Inspecionar dependências/modelos/casos de uso por Projeto/Catálogo/Tela/Jornada/Demonstração e exercitar os dois adaptadores; sem regras ou aparência paralelas | Não executado |
| Q038-19 | FR-004/FR-037 | Aplicar schemas/DTOs com campos extras, duplicatas, frame de entrada/saída, revisão máxima sem duplicação Content/StructuredContent, pacotes aninhados excessivos, propriedade de família errada e assinatura fora de tamanho; recusar sem frame parcial ou efeito indevido | Não executado |
| Q038-20 | FR-012/FR-035 | Interferir com arquivo externo entre conferência e publicação no EXE Windows; provar proteção ou desabilitar substituição com novo destino, sem afirmar CAS por hash+rename | Não executado |
| Q038-21 | FR-009–FR-011/FR-016/FR-037; U1 | Com oito sessões e retenção próxima ao teto, admitir a 1.024ª edição; negar/expirar repetidamente e preparar revisão nova sem renovar replay. Finalizar/descarte ainda disponíveis; com 63 commits, reservar por alvo distinto e recusar novo alvo antes de mutar se exceder 64; recuperar intent na posição original, sem 65ª. Verificar staging/memória reservados e resultado sem gravação quando falta consentimento ou há falha externa | Não executado |
| Q038-22 | FR-020/FR-021/FR-024/FR-030; I1 | Executar testes/codec T068/T070 e 002/T018 antes de US3 publicar T056 ou entregar T066; pacote final validado/hash antes da revisão, sem exigir UI de transferência US4 nem regenerar bytes após consentimento | Não executado |
| Q038-23 | FR-004/FR-021/FR-029/FR-037; I2 | Usar ASCII, acentos, emojis e combinações em GUI/MCP/pacote/proposta: 65.536 pontos de código aceitos, 65.537 recusados onde esse é o limite; nomes e campos específicos menores respeitados. 20.000 emojis=80.000 bytes UTF-8 não excedem o limite do campo; frame/pacote em bytes validado separadamente. Rejeitar UTF-8/surrogates inválidos sem substituir/truncar/normalizar; leitura histórica não trunca valores | Não executado |
| Q038-24 | FR-027–FR-031; I3 | Conferir defaults unknown/observed/proposed, tipos string/integer/boolean/enum/reference/collection, motivo/evidência exigidos, enum/ref inválidos e objetos/coleções heterogêneas recusados. Registros complexos permanecem unknown com limitação; valor obrigatório não resolvido bloqueia publicação. Evidência rotulada pelo agente não prova homologação nem altera default instalado | Não executado |
| Q038-25 | FR-009/FR-010/FR-016/FR-028/FR-029; U2 | Criar sessão blank sem arquivo, propor seleção/legado com sessionId/requestId/expectedRevision, conferir ownership/revisão própria/CAS/replay sem mudar Project.Revision/undo. Alterar documento invalida validação/review afetadas; limite 32 anexos/memória preserva estado. Salvar documento não salva proposta; dirty agregado impede fechamento silencioso e publicação limpa apenas o anexo correspondente | Não executado |

## Cenários adicionais do segundo replanejamento

| ID | Requisito/cenário | Ação futura e resultado esperado | Resultado observado |
| --- | --- | --- | --- |
| Q038-26 | FR-018/FR-028–FR-030; U3; US5.9 | Criar análise legacy com materialization=null: validade estrutural possível, publicação bloqueada. Construir manualmente e por I1, capturar selection, vincular mesma sessão/revisão/hash e revisar pacote/proveniência. Negar vínculos para newKind, outra sessão, revisão/hash errados; alterar captura/análise/evidência invalida Review. Publicar somente após autorização; uma reserva por captura, análise continua transitória e descarte ao fechar é informado | Não executado |
| Q038-27 | FR-011/FR-020/FR-025/FR-036/FR-037; U4; US3.9 | Para publish/withdraw/import/export, executar componente_revisar e confirmar com a ferramenta R2 correspondente; negar operação trocada/campos extras/reviewId inexistente/assinatura falsa ou expirada. Admissão de importação usa componente_preparar_importacao e quota de edição, revisão não cria anexo novo; exportação regular .jcomp é exclusiva no workspace autorizado. Replay não renova validade ou duplica efeito; alteração de índice/base/pacote exige nova revisão | Não executado |

## Registro a preencher somente após execução

Para cada cenário registrar revisão/hash, ambiente/cliente/versão, mecanismo/política de consentimento, executor, data, ação observada, resultado real e evidência pertinente. Metadados de auditoria não devem expor credenciais, segredos de autorização ou conteúdo pessoal. Resultado esperado da tabela não deve ser copiado como observado.

Checks automatizados, contrato comportamental e revisão de Markdown não comprovam teste do EXE ou consentimento independente. Office/DPI e aceite visual são verificações próprias. Usuário deve testar a mesma revisão entregue e aprová-la explicitamente; nenhuma aprovação de software foi recebida.

## Fundação de autoria — T006/T009

Registro e anexos implementados em internal/application/authoring.go e authoring_draft.go. go test -race ./internal/application -run '^TestAuthoring|^TestMCPFixture$' -count=1: PASS; go vet ./internal/application: PASS. Casos: sessões 8/9 e 30 min, quatro chamadas/busy, isolamento/erro sem mutação, 1.024 edições e revisão/descarte ainda disponíveis, anexos 32/33 com CAS/ownership/revisão independente, quota de retenção incluindo histórico/recursos, staging antes da admissão e 64 posições com intent/reconciliação sem 65ª. Livro de posições é coordenação em memória; não demonstra gravação/journal/consentimento de US2. O núcleo ainda não é exposto pelo MCP G0, que conserva descoberta mínima.

Cliente pretendido informado pelo usuário: Codex, modelo descrito literalmente como “deepseek v4 v1 flash”. Produto/versão/configuração efetivos e suporte a HostApproval-v1 serão conferidos no gate T042; não registrar esse relato como host qualificado.

## Fundação de IO e protocolo — evidências posteriores ao build G0

Fontes posteriores a 1.0.0.16, ainda não representados pelo EXE G0. Executor: agente de implementação no Windows/amd64, Go 1.27.2/CGO=1, em 2026-10-25. Nenhuma homologação humana implícita.

- 038/T007/T011: go test -race ./internal/adapter/workspaceguard -count=1 PASS; go vet PASS. Exclusão real em subprocesso, namespace estável com JORNADA_CONFIG_DIR diferente, names/ADS/dispositivos, leitura/base, criação exclusiva/cancelamento, junction e raiz renomeada: PASS. Symlink de arquivo não pôde ser criado por privilégio do ambiente; junction real foi usada. Substituição existente e todos os writers ainda dependem de US2.
- 034/T026–T028 e 038/T012: regressão falhou inicialmente por OpenWorkspace ausente no pacote extraído; depois go test -race ./internal/adapter/workspace PASS e go test ./internal/adapter/desktop -run 'TestWorkspace|TestLauncher' -count=1 PASS. Importação conserva bytes, homônimos não são substituídos, fundo incorporado/exemplo preservados, descendentes ignorados, corrupção não omitida, base muda ao trocar pasta e roots de finalidade não concedem escrita. Fonte/coordenador compartilhados testam replay sem reler origem e rascunho sem persistir.
- 038/T008/T010 parcial: go test -race ./internal/adapter/mcp PASS; testes Unicode/framing e decoder fechado, campo obrigatório/case exata, respostas grandes/Content curto, cursores autenticados por revisão/escopo e orçamento de página: PASS. authoring_validation_test comprova 65.536/65.537 pontos de código, emoji e caracteres combinados; integração de conteúdo à GUI/MCP ainda pertence a T022.
- Auditoria parcial T013: teste de allowlist PASS; nomes livres, caminhos, mensagens internas e segredos não são serializados. O novo servidor/handlers ainda precisa compor a auditoria em stderr. Nenhum resultado esperado de Q038 foi promovido a observado por esses testes.

## Fundação composta — T010/T013/T014, EXE 1.0.0.17

- Artefato apps/gerador-jornada/dist/windows-amd64-038-foundation/JornadaStudio.exe, PE Windows/amd64 com subsistema windowsgui padrão; SHA-256 86e622956a0a6e8548e083feb2181b40f3d1be5e9a8ae70d18c646c13a35a0cb. Build local scripts/build.ps1; sem envio externo.
- TestDistributedFoundation falhou primeiro contra 1.0.0.16 (startupProbe=true). PASS contra 1.0.0.17 via JORNADA_MCP_EXE: negociação/listTools/capacidades, startupProbe=false, persistentWrites=false, ausência de janelas Win32/preferências/arquivos e EOF com exit normal. TestDistributedMCPRejectsMixedFlags PASS para --open/--save-example junto ao MCP e --legacy-root sem MCP, sem efeitos persistentes. A prova G0 separada continuou PASS no seu próprio EXE.
- Cliente técnico SDK Go v1.8.0, Windows/amd64; executor agente em 2026-10-25. TestServerFoundationUsesStrictBoundary PASS -race: campos extras recusados sem eco, capacidades limitadas e auditoria em stderr. Testes de orçamento comprovam quatro chamadas/busy, cancelamento livre e prazos 30/120 segundos. DTOs de sete famílias e dez Kinds suportados para seleção, uint64 máximo/overflow, null em mutação e nomes de campo exatos testados. As ferramentas de edição ainda são futuras e não são anunciadas.
- go test ./... e go vet ./... PASS após extração; testes focados -race e vet de aplicação/MCP/workspace/guard PASS após composição. Próxima parcela US1 adicionará ferramentas e suas regressões, sem reaproveitar o hash desta fundação como prova de fontes posteriores. Host humano/consentimento, substituição Windows, biblioteca, visual/DPI/Office e aceite continuam pendentes.

## Autoria transitória — fontes posteriores à fundação 1.0.0.17

Revisão de fonte 1.0.0.18 em desenvolvimento; sem presumir build/homologação até evidência própria. Vinte ferramentas registradas para descoberta, workspace, rascunho, tela, componente, preset, navegação e validação; revisão congelada/controles de sessão/consentimento/biblioteca ainda dependem das tarefas abertas.

- Testes inicialmente RED por catalogo_listar/componente_inserir/jornada_criar/jornada_abrir inexistentes. Após integração, discovery/edit_contract/authoring_flow e comparação desktop/MCP PASS. O teste completo de patch encontrou schema de ponteiro para união aceitando só null e autorreferência indevida na verificação de ID de linha nova; corrigidos, com regressões retidas. Não foi marcado aceite humano.
- Dez Kinds inseridos, sete patches específicos e campos comuns de compostos, movimento com replay, remoção recusando referências externas, bindings próprios explícitos, duas telas/navegação/duplicação/IDs/FileID: PASS -race. Validação por alvo é consulta e não exporta. Cursores e respostas respeitam o orçamento; erro não transcreve dados internos.
- TestManualMCPPresetKeepsLegacyDecoration usa editor Fyne de teste e servidor stdio real: faixa com mesmos bounds/estilo/texto/recurso e origem intacta, PASS. TestManualTextUsesSameCodePointLimit usa o botão Salvar alterações do formulário: 65.536 emoji aceitos e excesso recusado sem alteração, PASS. Não constitui captura nativa, DPI, Office ou homologação visual humana.
- TestMCPTextLimitAndMetadataQuery e TestHistoricalContentRemainsReadableAndUnsaved: Unicode integral, limite por pontos de código, leitura de pacote histórico excedente sem truncar e nenhuma escrita na origem, PASS. ValidateComponentEdit é compartilhado com formulário e preserva campo histórico inalterado.
- TestDuplicateRegeneratesAllInstanceIDs: reproduzido RED ao reexecutar a chamada anterior RegenerateCompoundIDs; PASS ao incluir opções/linhas e seleção em RegenerateInstanceIDs. Corrige cópia de tela/seleção nos dois canais. Mapa de resposta é limitado antes de publicar a edição, sem erro de tamanho posterior a um efeito aceito.
- scripts/check.ps1 PASS após a primeira parcela de autoria: gofmt, skins/menu --check (205 disponíveis, 8 indisponíveis, 224 recursos), go vet ./..., três testes de empacotamento, dois de configuração OAuth e go test ./.... Testes focados posteriores de IDs/movimento/stream também PASS; novo check/build ainda devem identificar exatamente a revisão entregue.
- Pré-requisito de T026: projectzip.Write passou -race com cancelamento/short write/round-trip e golden dos bytes anteriores. Baseline determinístico: 439 bytes, SHA-256 071d6970f7dd6b0ca04ffc5ee6db6d14af35c620be5b01dc21c2ba596cc6d55b. Encode e Write compartilham uma codificação; Write não fecha destino nem publica ZIP parcial. Não comprova staging/revisão/consentimento ainda ausentes.

## Executável de autoria transitória — 1.0.0.18

Artefato local: apps/gerador-jornada/dist/windows-amd64-038-authoring/JornadaStudio.exe. SHA-256: 03d0f192ff81fe87f8b3b94cbd25a118d742e9bb86b57427799e5cb7c1ad753c. Build scripts/build.ps1 padrão Windows/amd64, subsistema windowsgui. Distribuição local dist/JornadaStudio-1.0.0.18-windows-amd64.zip com inventário/licenças; nenhum commit/merge/envio externo.

Executor: agente de implementação, Windows/amd64, Go 1.27.2/CGO=1, 2026-10-25; cliente técnico distributed-authoring-harness 1 com SDK MCP Go v1.8.0. TestDistributedAuthoringDraft falhou no EXE 1.0.0.17 por jornada_criar ausente; PASS no 1.0.0.18 para criar rascunho Unicode, consultar tela/revisão, inserir botão, criar segunda tela, ligar navegação e validar draft. persistentWrites=false/host não qualificado, ausência de janelas Win32 do PID, auditoria sem nomes dos projetos, EOF normal e workspace/configuração isolados vazios: PASS. Não houve salvamento ou aprovação humana nesse teste.

Comandos: carregar scripts/env.ps1, definir JORNADA_MCP_EXE com caminho absoluto do EXE 1.0.0.18 e executar go test ./internal/adapter/mcp -run '^TestDistributed' -count=1 -v. JORNADA_MCP_PROBE_EXE pode apontar separadamente a G0 histórica. Startup/flags incompatíveis continuam PASS na revisão nova. Sem as variáveis, os testes pulam e não constituem evidência.

Roteiro manual disponível: abrir o mesmo EXE normalmente, criar jornada com fundo legado, inserir botão/campo e faixa de título, salvar pelo editor, reabrir e conferir defaults/overrides. Para teste MCP de produção, configurar processo stdio com workspace de cópias; conferir aplicacao_capacidades e usar somente operações anunciadas. Criar/editar deve modificar apenas rascunho até que as tarefas de R2 sejam entregues. Não aguardar salvamento automático no EOF. O usuário ainda não executou/aprovou essa revisão; comparação visual, DPI, PowerPoint, host de produção e biblioteca têm gates próprios.

## Revisão interna e controles de sessão — 1.0.0.19, 2026-10-27

Executor: agente de implementação; Windows/amd64, Go 1.27.2/CGO=1, SDK MCP Go v1.8.0. Esta seção supera o estado dos fontes descrito nas seções históricas de 1.0.0.17/18, sem alterar suas evidências.

Artefato local: apps/gerador-jornada/dist/windows-amd64-038-sessions/JornadaStudio.exe. SHA-256 c3be011bc6208d78d5c629fd6c0c1d52d841f404972c3b5a0b5da3b44d263053. Build scripts/build.ps1 padrão Windows/amd64/windowsgui. Pacote local dist/JornadaStudio-1.0.0.19-windows-amd64.zip, com README/licenças/inventário de hashes. Sem commit, merge ou publicação externa.

- T026: Authoring→Freezer real prepara ZIP e revisão conforme authorization.schema.json, com ordem canônica, IDs distintos, hash do pacote final, desafio/validade fixos e proposal/materialization null para jornada. Replay conserva bytes/artefato; repreparo aposenta o anterior; descarte libera staging sem publicar destino. TestFreezerFollowsSessionReviewLifecycle PASS. Raiz privada dentro do workspace e base alterada por terceiro são recusadas; falha depois da codificação limpa o temporário. Guard Windows comprova segundo escritor recusado, delete-on-close e queda de subprocesso sem defer. Casos -race PASS. Não há chave privada, aprovação humana ou verificador de consentimento nessa preparação.
- Orçamento: 64 MiB globais reservados para uma finalização pesada e 1 MiB/sessão para metadados/revisão dentro dos 256 MiB; limites de staging preservados. TestReviewRetirementKeepsBudgetUntilPhysicalRelease comprova reserva retida e consulta disponível enquanto outro artefato é liberado. TestIdleExpirationPreservesUncertainIntent conserva posição/prioridade de recovery. Esses testes não comprovam journal persistente.
- Unicode do manifesto: regressão RED para surrogate isolado aceito pelo decoder anterior; após jsonwire, UTF-8 inválido, surrogate isolado e duplicatas são recusados, inclusive em objeto aninhado. Texto Go inválido não é normalizado silenciosamente ao serializar. Write mantém o golden pré-extração de 439 bytes e Verify usa o mesmo handle sem reconstruir imagens em memória.
- T028: testes de sessão pelo SDK PASS -race para dirty agregado de anexo, fechamento recusado sem descarte, resposta estável de descarte/fechamento, revisão monotônica, quotas/reservas liberadas, quota de edição cheia e sessão expirada sem recriação. Varredura ociosa e encerramento da goroutine testados. SessionStatus não copia o documento para consultar controles.
- TestDistributedSessionControls RED contra 1.0.0.18 por sessao_obter inexistente, depois PASS no EXE 1.0.0.19. Admissão informa prazo/política; sessão nova está dirty e reservada; fechamento exige descarte; replay terminal é exato e consulta posterior retorna sessionExpired; EOF termina normalmente, sem janela Win32 do PID ou arquivos no workspace/configuração. persistentWrites=false e host não qualificado. Auditoria não transcreve nome do projeto.
- TestDistributedAuthoringDraft/Foundation/MCPRejectsMixedFlags PASS no EXE novo; StartupProbe PASS separadamente no G0 1.0.0.16. scripts/check.ps1 PASS para gofmt, recursos legados --check, go vet ./..., empacotamento/OAuth e go test ./.... Suites de aplicação, MCP, consent, projectzip e workspaceguard PASS -race; vet dos pacotes alterados PASS. A ausência das variáveis de EXE pula os testes distribuídos e não vale como prova.

Reproduzir: carregar scripts/env.ps1, definir JORNADA_MCP_EXE com o caminho absoluto do artefato acima e executar go test ./internal/adapter/mcp -run '^TestDistributed' -count=1 -v; fornecer JORNADA_MCP_PROBE_EXE com o G0 somente para sua prova histórica.

Roteiro manual deste incremento: abrir normalmente o EXE 1.0.0.19 e criar/reabrir uma cópia de jornada no editor. No cliente MCP, iniciar o mesmo EXE com --mcp --workspace apontando uma pasta de cópias; conferir 23 ferramentas e persistentWrites=false. Criar rascunho, consultar sessao_obter, editar e tentar encerrar com discardUnsaved=false: deve recusar sem perda de conteúdo. Consultar os contadores correntes e descartar por CAS; conferir base restaurada e reservas liberadas. Encerrar a sessão limpa; somente o replay terminal idêntico deve retornar o fechamento conhecido, sem reabrir a sessão. Outra sessão nova pode ser encerrada com discardUnsaved=true. Confirmar que nenhum fluxo MCP gravou .jornada. A preparação de revisão permanece interna até T035/T038 e não deve ser procurada na descoberta desta revisão.

Próximos gates: T027 paridade das cenas/canvas/prévia/PPT, T029/T030 integração ampliada, T031 execução humana; depois consentimento/journal/coordenação de writers e qualificação real do host em US2. Biblioteca manual/MCP e propostas do legado permanecem nos incrementos seguintes. Nenhum Q038 inteiro, Office/DPI ou aceite funcional é promovido a aprovado por estes testes.
