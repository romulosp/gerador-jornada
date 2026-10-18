# Configura somente o processo atual com as ferramentas locais do projeto.
$ErrorActionPreference = 'Stop'
$appRoot = Split-Path -Parent $PSScriptRoot
$goBin = Join-Path $appRoot '.tools\go\bin'
$cBin = Join-Path $appRoot '.tools\mingw64\bin'
if (-not (Test-Path -LiteralPath (Join-Path $goBin 'go.exe'))) {
    throw 'Ferramentas locais ausentes. Execute python -X utf8 scripts/bootstrap.py.'
}
$env:PATH = $goBin + ';' + $cBin + ';' + $env:PATH
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '1'
$env:CC = Join-Path $cBin 'gcc.exe'
$env:CXX = Join-Path $cBin 'g++.exe'
