# Incorpora o ícone do Windows sem dependências adicionais no computador do usuário.
$ErrorActionPreference = 'Stop'
# iconAppRoot delimita os caminhos do aplicativo, inclusive com espaços ou acentos.
$iconAppRoot = Split-Path -Parent $PSScriptRoot
# sourcePath aponta para a cópia versionada da imagem original, sem alterações visuais.
$sourcePath = Join-Path $iconAppRoot 'internal\adapter\builtinassets\assets\app-icon.png'
# outputDirectory reúne apenas arquivos intermediários reproduzíveis e ignorados pelo Git.
$outputDirectory = Join-Path $iconAppRoot 'build\windows-icon'
New-Item -ItemType Directory -Path $outputDirectory -Force | Out-Null
Add-Type -AssemblyName System.Drawing
# source carrega a imagem que será reduzida proporcionalmente para as escalas do Windows.
$source = [System.Drawing.Image]::FromFile($sourcePath)
# images armazena cada PNG do contêiner ICO, do tamanho pequeno ao de alta resolução.
$images = [System.Collections.Generic.List[byte[]]]::new()
# sizes contém os tamanhos utilizados por Explorer, janela e barra de tarefas.
$sizes = @(16, 24, 32, 48, 64, 128, 256)
try {
    if ($source.Width -ne $source.Height) { throw 'A imagem do ícone deve ser quadrada.' }
    # size determina a resolução de cada representação sem recortar a imagem.
    foreach ($size in $sizes) {
        # bitmap recebe o quadro com interpolação de alta qualidade.
        $bitmap = [System.Drawing.Bitmap]::new($size, $size)
        # graphics executa somente a redução para a representação do sistema operacional.
        $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
        # stream recebe os bytes PNG preservando o canal alfa da imagem.
        $stream = [System.IO.MemoryStream]::new()
        try {
            $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
            $graphics.DrawImage($source, 0, 0, $size, $size)
            $bitmap.Save($stream, [System.Drawing.Imaging.ImageFormat]::Png)
            $images.Add($stream.ToArray())
        } finally {
            $stream.Dispose()
            $graphics.Dispose()
            $bitmap.Dispose()
        }
    }
} finally { $source.Dispose() }
# iconPath recebe o diretório ICO e seus quadros PNG, suportados no Windows 10/11.
$iconPath = Join-Path $outputDirectory 'jornada-studio.ico'
# writer grava números little-endian conforme o formato ICO.
$writer = [System.IO.BinaryWriter]::new([System.IO.File]::Create($iconPath))
try {
    $writer.Write([uint16]0)
    $writer.Write([uint16]1)
    $writer.Write([uint16]$sizes.Count)
    # offset aponta para o primeiro quadro, após o cabeçalho e as entradas de 16 bytes.
    $offset = 6 + 16 * $sizes.Count
    # index relaciona resolução, comprimento e posição de cada quadro no arquivo.
    for ($index = 0; $index -lt $sizes.Count; $index++) {
        # dimension representa 256 por zero, como exigido pelo diretório ICO.
        $dimension = if ($sizes[$index] -eq 256) { 0 } else { $sizes[$index] }
        $writer.Write([byte]$dimension)
        $writer.Write([byte]$dimension)
        $writer.Write([byte]0)
        $writer.Write([byte]0)
        $writer.Write([uint16]1)
        $writer.Write([uint16]32)
        $writer.Write([uint32]$images[$index].Length)
        $writer.Write([uint32]$offset)
        $offset += $images[$index].Length
    }
    # frame mantém a ordem das entradas do diretório.
    foreach ($frame in $images) { $writer.Write([byte[]]$frame) }
} finally { $writer.Dispose() }
# versionText lê a fonte única da versão exibida na interface e no Windows.
$versionText = [System.IO.File]::ReadAllText((Join-Path $iconAppRoot 'internal\adapter\buildinfo\VERSION')).Trim()
if ($versionText -notmatch '^\d+\.\d+\.\d+\.\d+$') { throw 'A versão deve conter quatro partes numéricas.' }
# part valida o limite de 16 bits de cada parte do recurso VERSIONINFO.
foreach ($part in $versionText.Split('.')) {
    if ([uint64]$part -gt 65535) { throw 'Cada parte da versão deve estar entre 0 e 65535.' }
}
# versionNumbers adapta a versão textual ao formato numérico exigido pelo compilador de recursos.
$versionNumbers = $versionText.Replace('.', ',')
# resourcePath descreve ícone e versão do executável, sem alterar o formato dos projetos.
$resourcePath = Join-Path $outputDirectory 'icon.rc'
# resourceText mantém as versões de produto e arquivo sincronizadas com a interface.
$resourceText = @"
1 ICON "jornada-studio.ico"
1 VERSIONINFO
FILEVERSION $versionNumbers
PRODUCTVERSION $versionNumbers
FILEFLAGSMASK 0x3fL
FILEFLAGS 0
FILEOS 0x40004L
FILETYPE 0x1L
FILESUBTYPE 0
BEGIN
    BLOCK "StringFileInfo"
    BEGIN
        BLOCK "041604b0"
        BEGIN
            VALUE "FileDescription", "Gerador de Jornada\0"
            VALUE "FileVersion", "$versionText\0"
            VALUE "InternalName", "JornadaStudio\0"
            VALUE "OriginalFilename", "JornadaStudio.exe\0"
            VALUE "ProductName", "Jornada Studio\0"
            VALUE "ProductVersion", "$versionText\0"
        END
    END
    BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x0416, 1200
    END
END
"@
[System.IO.File]::WriteAllText($resourcePath, $resourceText, [System.Text.Encoding]::ASCII)
# compiler utiliza a mesma distribuição MinGW já necessária ao Fyne.
$compiler = Join-Path $iconAppRoot '.tools\mingw64\bin\windres.exe'
# targetPath é carregado automaticamente pelo compilador Go apenas no alvo Windows amd64.
$targetPath = Join-Path $iconAppRoot 'cmd\gerador-jornada\icon_windows_amd64.syso'
Push-Location $outputDirectory
try {
    & $compiler --input icon.rc --output $targetPath --output-format coff --target pe-x86-64
    if ($LASTEXITCODE -ne 0) { throw 'Falha ao incorporar o ícone no executável Windows.' }
} finally { Pop-Location }
