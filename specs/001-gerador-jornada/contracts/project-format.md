# Contract: formato de projeto .jornada

**Version**: 2 | **Consumers**: editor e adaptador projectzip | **Requirements**: FR-001, FR-009, FR-010, NFR-005.

## Pacote e wire format

Arquivo ZIP com `manifest.json` UTF-8 na raiz e imagens em `assets/<UUID>.png` ou `.jpg`. Sem links externos, código, templates arbitrários, JARs ou scripts no projeto.

Manifesto usa camelCase, sem funções. Os campos de raiz são:

| Campo | Tipo | Regra |
| --- | --- | --- |
| schemaVersion | integer | Writer: 2. Reader: 1 com migração, ou 2 com validação de fileId. |
| catalogVersion | integer | Exatamente 1. |
| id, name | string | UUID e nome não vazio. |
| revision | integer | Inteiro não negativo; identifica a revisão salva. |
| screens | array | Screen na ordem de apresentação. Pode ser vazio em rascunho. |
| journey | object | startScreenId (UUID ou string vazia quando ausente; reader também aceita null) e actions (array). |
| scenarios | array | Scenario; vazio quando não há demonstração. |
| assets | array | Asset; vazio quando não há imagens. |
| javaOptions | object | packageName e profileId, strings; podem estar vazios no rascunho. |

Screen contém `id`, `name`, `fileId`, `size:{width,height}`, `background:{color,assetId}` e `components`. Component contém `id`, `name`, `fileId`, `kind`, `bounds:{x,y,width,height}`, `style` e `properties`. assetId vazio significa sem imagem; o reader também aceita null. Style permite os campos descritos no modelo, omitindo os não personalizados.

Todos os arrays são obrigatórios, inclusive vazios. Propriedades e campos de outras entidades seguem [catálogo](component-catalog.md) e [jornada/roteiro](navigation-scenario.md). Campos desconhecidos são recusados pelo leitor, com localização; não são ignorados silenciosamente.

### Exemplo mínimo válido como rascunho

```json
{
  "schemaVersion": 2,
  "catalogVersion": 1,
  "id": "8be7fa20-bf51-4d26-9bc9-cae9fd075d01",
  "name": "Atendimento",
  "revision": 0,
  "screens": [],
  "journey": { "startScreenId": "", "actions": [] },
  "scenarios": [],
  "assets": [],
  "javaOptions": { "packageName": "", "profileId": "" }
}
```

Esse arquivo pode ser salvo/reaberto; exportação retorna `JOURNEY_START_MISSING`.

## Ausência de referências opcionais no wire format

O writer atual representa ausência em campos Go string por `""`, incluindo
startScreenId, targetScreenId de ação none, selectedOptionId, assetId e alvos opcionais.
O reader também aceita `null` para esses campos, normalizando-o para a string vazia.
Nos modelos conceituais, “nulo” descreve ausência de referência, não exige JSON null.
Opcionais tipados por ponteiro, como Step.point e Step.text, conservam JSON null;
texto vazio de typeText continua sendo um valor explícito diferente de null.
Essa descrição registra o comportamento do writer/reader já existente, sem migração
adicional nem mudança de schemaVersion 2.

## Recursos e limites de leitura

Asset contém `id`, `path`, `mediaType`, `sha256` hexadecimal minúsculo e `sizeBytes`. Aceitar somente PNG/JPEG decodificáveis e dimensões até 8192 × 8192. Sem SVG executável ou URLs remotas nas versões 1 e 2.

Limites iniciais para leitura: manifesto descomprimido até 10 MiB, recurso até 25 MiB, pacote até 250 MiB descomprimidos no total e até 1000 entradas. Verificar os bytes efetivamente lidos, não apenas o tamanho informado no cabeçalho. Mensagem deve indicar o limite excedido.

Recusar caminhos absolutos, unidades Windows, `..`, barras invertidas, nomes duplicados ou equivalentes sem distinção de maiúsculas, symlinks e entradas fora de manifest/assets. Nenhuma extração pode escapar da área temporária. Hash e tamanho precisam corresponder aos bytes.

Versão desconhecida gera `PROJECT_VERSION_UNSUPPORTED`; corrupção/ZIP inválido, `PROJECT_CORRUPT`; tipo/campo desconhecido, `PROJECT_SCHEMA_INVALID`. Erros de referência em conteúdo legível podem abrir rascunho com diagnósticos de domínio.

## Gravação e publicação

1. Capturar snapshot da revisão a salvar e validar estrutura/recursos. Permitir diagnósticos semânticos para salvar rascunho.
2. Escrever ZIP temporário no mesmo volume do destino, fechar streams, sincronizar e reler sua integridade.
3. Publicar com operação de substituição atômica compatível com Windows. Se não houver garantia de substituir com segurança, falhar preservando o arquivo anterior.
4. Atualizar revisão salva somente após publicação. Edições posteriores à revisão capturada permanecem pendentes.

Solicitar confirmação apenas para substituir um destino escolhido que já exista; autosave não dispara esse diálogo. Falta de permissão, disco cheio e arquivo bloqueado são erros recuperáveis.

Ordenar entradas e manter metadados ZIP normalizados para comparação determinística. Não gravar seleção, cursor, zoom, janela, undo/redo ou valores runtime.

## Autosave e recuperação

Após comando confirmado, debounce de 2 s grava cópia no diretório de recuperação do usuário: Windows `%LOCALAPPDATA%/GeradorJornada/recovery/<projectId>.jornada`; Linux `$XDG_CONFIG_HOME/GeradorJornada/recovery` (ou `~/.config`); macOS `~/Library/Application Support/GeradorJornada/recovery`. A resolução pertence ao adaptador; erro não deve cair em diretório relativo. Máximo uma cópia válida por projeto, substituída com o mesmo procedimento seguro.

Autosave não sobrescreve arquivo explícito nem marca sessão limpa. Em interrupção, oferecer recuperar a revisão mais recente válida; recuperação abre como alterado e exige salvar explicitamente. Descartar recuperação requer ação do usuário. Autosave falho é indicado sem interromper edição.

Salvar/fechar com sucesso pode remover a recuperação apenas quando ela não contém revisão posterior à salva. Em diretório sem permissão, explicar que recuperação automática está indisponível.

## Testes de contrato

Round-trip de todas as famílias; preservação de IDs/ordem; campos desconhecidos; versão futura; JSON truncado; UUID repetido; path traversal; recurso alterado; limites ZIP; falha antes da substituição; arquivo bloqueado; recuperação e edição concorrente ao salvamento.

## Nome-base de arquivos e migração v1 → v2

Screen.fileId e Component.fileId são obrigatórios na versão 2. O campo `id` continua sendo UUID interno, e `name` continua sendo nome amigável Unicode. FileID não substitui UUID em ações, roteiros, recursos ou vínculos de teclado.

Regra comum: `^[A-Za-z][A-Za-z0-9]{0,39}$`, sem extensão. Telas exigem adicionalmente prefixo exato Pnl e ao menos um caractere após ele; prefixo conta no limite de 40. Componentes não exigem Pnl. Recusar dispositivos Windows CON, PRN, AUX, NUL, COM1–COM9 e LPT1–LPT9 sem distinguir caixa. Unicidade: telas por projeto e componentes por tela, sem distinguir caixa.

O leitor aceita v1 e preenche somente fileIds ausentes, sem alterar UUIDs ou conteúdo. Sugestão de tela remove a ordenação numérica inicial e acentos: `01 · Identificação` → `PnlIdentificacao`. Sugestões vazias usam PnlTela/Componente; colisões de sugestões recebem sufixos 2, 3, … com espaço reservado dentro de 40 caracteres. Valores explícitos não são renomeados. O próximo salvamento grava v2; o arquivo original não é escrito durante a leitura. V2 com valor ausente, inválido ou duplicado gera PROJECT_SCHEMA_INVALID; versões fora de 1/2 geram PROJECT_VERSION_UNSUPPORTED. Software 1.0.0.0 e catalogVersion 1 permanecem independentes da versão do manifesto.

Exercitar migração, estabilidade das sugestões, UUID/navegação preservados, round-trip v2, colisões de caixa, 40/41 caracteres, acentos, caminhos e nomes Windows reservados.
