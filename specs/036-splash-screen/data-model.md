# Data Model: Adaptadores da splash

Modelo puro referenciado de 035/data-model.md; não duplicar suas regras aqui.

- `desktop.Splash`: janela, área de vídeo, contexto/cancelamento, sessão `startup.Session`, player e callbacks next/finished.
- `desktop.splashPlayback`: fronteira local Snapshot/Resize/Close; estado traduzido de mídia inclui Started, Ended, Audio, Err.
- `nativevideo.Player`: ponteiro nativo privado. Callback COM e estado atômico ficam em C++; nunca guarda widgets, documento ou ponteiro Go em callback COM.
- Recurso v2: MP4 original incorporado, 1280×720, duração de contêiner 10,006 s, vídeo e áudio; SHA-256 `3aa7a762e16da881509be1e7802c84f24e5d5848e14f8db4986332846c7bf3b3`.
- Cache de reprodução: arquivo por hash sob cache do usuário; escrita atômica, conferência antes de reutilizar. Cache não é arquivo de projeto.
- `skinImages` e `sceneFonts`: caches operacionais reutilizados por canvas/prévia, sem alterar pixels/alpha ou `.jornada`.
- `StudioCredit`: constante de apresentação do launcher. A splash não contém autoria ou outros textos.

Nenhum campo novo em Project/Document/MenuItem; formato permanece schema 3/catalog 4.

## Preferência persistida — 1.0.0.10

`startup-preferences.json` em `%APPDATA%/JornadaStudio` ou `JORNADA_CONFIG_DIR`, com `{"presented":false,"show_on_startup":false}` como valor ausente. Tipo de regra em 035; IO em `desktop/startup_preferences.go`. Registro da primeira apresentação altera somente presented. Checkbox altera somente show_on_startup. Publicação atômica; falha não confirma seleção e não descarta arquivo existente. Não contém versão nem caminho do projeto; atualizações do EXE não reiniciam o primeiro uso.

## Apresentação em janela — 1.0.0.11

Somente configuração visual do adaptador: tamanho inicial 960×540 unidades Fyne; fullscreen=false em splash_open. Nenhum campo novo na política por usuário, sessão de domínio ou projeto. Geometria técnica do teste distingue retângulo externo e área cliente física da janela, considerando a moldura padrão do Windows.
