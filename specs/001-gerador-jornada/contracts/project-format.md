# Contract: formato de projeto .jornada

**Version**: 1 | **Consumers**: editor e adaptador projectzip | **Requirements**: FR-001, FR-009, FR-010, NFR-005.

## Pacote e wire format

Arquivo ZIP com `manifest.json` UTF-8 na raiz e imagens em `assets/<UUID>.png` ou `.jpg`. Sem links externos, código, templates arbitrários, JARs ou scripts no projeto.

Manifesto usa camelCase, sem funções. Os campos de raiz são:

| Campo | Tipo | Regra |
| --- | --- | --- |
| schemaVersion | integer | Exatamente 1. |
| catalogVersion | integer | Exatamente 1. |
| id, name | string | UUID e nome não vazio. |
| revision | integer | Inteiro não negativo; identifica a revisão salva. |
| screens | array | Screen na ordem de apresentação. Pode ser vazio em rascunho. |
| journey | object | startScreenId (UUID ou null) e actions (array). |
| scenarios | array | Scenario; vazio quando não há demonstração. |
| assets | array | Asset; vazio quando não há imagens. |
| javaOptions | object | packageName e profileId, strings; podem estar vazios no rascunho. |

Screen contém `id`, `name`, `size:{width,height}`, `background:{color,assetId}` e `components`. Component contém `id`, `name`, `kind`, `bounds:{x,y,width,height}`, `style` e `properties`. assetId nulo significa sem imagem. Style permite os campos descritos no modelo, omitindo os não personalizados.

Todos os arrays são obrigatórios, inclusive vazios. Propriedades e campos de outras entidades seguem [catálogo](component-catalog.md) e [jornada/roteiro](navigation-scenario.md). Campos desconhecidos são recusados no leitor v1, com localização; não são ignorados silenciosamente.

### Exemplo mínimo válido como rascunho

```json
{
  "schemaVersion": 1,
  "catalogVersion": 1,
  "id": "8be7fa20-bf51-4d26-9bc9-cae9fd075d01",
  "name": "Atendimento",
  "revision": 0,
  "screens": [],
  "journey": { "startScreenId": null, "actions": [] },
  "scenarios": [],
  "assets": [],
  "javaOptions": { "packageName": "", "profileId": "" }
}
```

Esse arquivo pode ser salvo/reaberto; exportação retorna `JOURNEY_START_MISSING`.

## Recursos e limites de leitura

Asset contém `id`, `path`, `mediaType`, `sha256` hexadecimal minúsculo e `sizeBytes`. Aceitar somente PNG/JPEG decodificáveis e dimensões até 8192 × 8192. Sem SVG executável ou URLs remotas na v1.

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

Após comando confirmado, debounce de 2 s grava cópia em `%LOCALAPPDATA%/GeradorJornada/recovery/<projectId>.jornada`. Máximo uma cópia válida por projeto, substituída com o mesmo procedimento seguro.

Autosave não sobrescreve arquivo explícito nem marca sessão limpa. Em interrupção, oferecer recuperar a revisão mais recente válida; recuperação abre como alterado e exige salvar explicitamente. Descartar recuperação requer ação do usuário. Autosave falho é indicado sem interromper edição.

Salvar/fechar com sucesso pode remover a recuperação apenas quando ela não contém revisão posterior à salva. Em diretório sem permissão, explicar que recuperação automática está indisponível.

## Testes de contrato

Round-trip de todas as famílias; preservação de IDs/ordem; campos desconhecidos; versão futura; JSON truncado; UUID repetido; path traversal; recurso alterado; limites ZIP; falha antes da substituição; arquivo bloqueado; recuperação e edição concorrente ao salvamento.
