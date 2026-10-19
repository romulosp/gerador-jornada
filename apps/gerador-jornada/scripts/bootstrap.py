"""Prepara ferramentas Go e MinGW locais com versões fixadas sem alterar o Windows."""
from pathlib import Path
import hashlib
import json
import urllib.request
import zipfile


def install(url, archive, destination, checksum=None):
    """Baixa e verifica o pacote antes de extrair as ferramentas locais de compilação."""
    if not archive.exists():
        urllib.request.urlretrieve(url, archive)
    if checksum and hashlib.sha256(archive.read_bytes()).hexdigest() != checksum:
        raise ValueError('Soma de verificação divergente: ' + str(archive))
    with zipfile.ZipFile(archive) as package:
        package.extractall(destination)


def main():
    """Instala as versões fixadas no manifesto versionado, conferindo os arquivos baixados."""
    # root identifica a raiz da aplicação a partir deste script.
    root = Path(__file__).resolve().parents[1]
    # tools isola as ferramentas locais do ambiente do Windows.
    tools = root / '.tools'
    tools.mkdir(exist_ok=True)
    # manifest versionado impede trocar silenciosamente a toolchain em uma instalação nova.
    manifest = root / 'scripts/toolchains-windows.json'
    # versions contém versões e hashes dos pacotes já utilizados pelo projeto.
    versions = json.loads(manifest.read_text(encoding='utf-8'))
    (tools / 'versions.json').write_text(json.dumps(versions, indent=2), encoding='utf-8')
    print(json.dumps(versions), flush=True)
    if not (tools / 'go/bin/go.exe').exists():
        install(versions['goUrl'], tools / 'go.zip', tools, versions['goSha256'])
        print('Go preparado', flush=True)
    if not (tools / 'mingw64/bin/gcc.exe').exists():
        install(versions['mingwUrl'], tools / 'mingw.zip', tools, versions['mingwSha256'] or None)
        print('MinGW preparado', flush=True)


if __name__ == '__main__':
    main()
