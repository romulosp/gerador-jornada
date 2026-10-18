# Compila o EXE x64 com recursos do Fyne e bibliotecas de suporte ligadas estaticamente.
param([switch]$Console)
. (Join-Path $PSScriptRoot 'env.ps1')
Push-Location $appRoot
try {
    New-Item -ItemType Directory -Path 'dist' -Force | Out-Null
    $linkFlags = '-s -w -extldflags=-static'
    if (-not $Console) { $linkFlags = '-H=windowsgui ' + $linkFlags }
    go build -trimpath -ldflags $linkFlags -o dist/JornadaStudio.exe ./cmd/gerador-jornada
    if ($LASTEXITCODE -ne 0) { throw 'Falha na compilação Go/Fyne.' }
    Get-FileHash -Algorithm SHA256 -LiteralPath 'dist/JornadaStudio.exe'
} finally {
    Pop-Location
}
