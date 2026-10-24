# Executa verificações reais sem confundir ausência de homologação com aprovação.
param([string]$Jdk17)

. (Join-Path $PSScriptRoot 'env.ps1')
# previousTestJdk preserva a configuração do processo chamador após a validação.
$previousTestJdk = $env:JORNADA_TEST_JDK
if ($Jdk17) { $env:JORNADA_TEST_JDK = $Jdk17 }
Push-Location $appRoot
try {
    $goFiles = Get-ChildItem internal,cmd -Recurse -Filter '*.go' -File | Select-Object -ExpandProperty FullName
    $unformatted = & (Join-Path $goBin 'gofmt.exe') -l $goFiles
    if ($unformatted) { throw ('Arquivos sem gofmt: ' + ($unformatted -join ', ')) }
    python -X utf8 scripts/import_legacy_skin.py --check
    if ($LASTEXITCODE -ne 0) { throw 'Biblioteca visual desatualizada.' }
    python -X utf8 scripts/import_legacy_menu.py --check
    if ($LASTEXITCODE -ne 0) { throw 'Árvore de menu desatualizada.' }
    go vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'go vet reprovado.' }
    python -X utf8 -m unittest discover -s scripts -p test_build.py -v
    if ($LASTEXITCODE -ne 0) { throw 'Validação do empacotamento de 64 bits reprovada.' }
    if ($env:JORNADA_TEST_JDK) {
        python -X utf8 -m unittest discover -s scripts -p test_java6.py -v
        if ($LASTEXITCODE -ne 0) { throw 'Validação da toolchain Java 6 reprovada.' }
    }
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw 'Testes Go reprovados.' }
    Write-Output 'Verificações automáticas aprovadas. Homologação manual, PowerPoint e EFL continuam independentes.'
} finally {
    $env:JORNADA_TEST_JDK = $previousTestJdk
    Pop-Location
}
