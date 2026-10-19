# Compila o EXE x64 com recursos do Fyne e bibliotecas de suporte ligadas estaticamente.
param([switch]$Console, [string]$OutputDirectory = 'dist')
. (Join-Path $PSScriptRoot 'env.ps1')
Push-Location $appRoot
try {
    New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
    & (Join-Path $PSScriptRoot 'windows-icon.ps1')
    $linkFlags = '-s -w -extldflags=-static'
    if (-not $Console) { $linkFlags = '-H=windowsgui ' + $linkFlags }
    $executablePath = Join-Path $OutputDirectory 'JornadaStudio.exe'
    go build -trimpath -ldflags $linkFlags -o $executablePath ./cmd/gerador-jornada
    if ($LASTEXITCODE -ne 0) { throw 'Falha na compilação Go/Fyne.' }
    Get-FileHash -Algorithm SHA256 -LiteralPath $executablePath
} finally {
    Pop-Location
}
