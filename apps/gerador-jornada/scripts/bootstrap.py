"""Prepara ferramentas Go e MinGW locais com versões fixadas sem alterar o Windows."""
from pathlib import Path
import hashlib
import json
import urllib.request
import zipfile


def fetch_json(url):
    """Lê os metadados das ferramentas disponibilizados pelo fornecedor."""
    with urllib.request.urlopen(url) as response:
        return json.load(response)


def install(url, archive, destination, checksum=None):
    """Baixa e verifica o pacote antes de extrair as ferramentas locais de compilação."""
    if not archive.exists():
        urllib.request.urlretrieve(url, archive)
    if checksum and hashlib.sha256(archive.read_bytes()).hexdigest() != checksum:
        raise ValueError('Soma de verificação divergente: ' + str(archive))
    with zipfile.ZipFile(archive) as package:
        package.extractall(destination)


def main():
    """Resolve versões estáveis uma vez e preserva suas versões exatas localmente."""
    # root identifica a raiz da aplicação a partir deste script.
    root = Path(__file__).resolve().parents[1]
    # tools isola as ferramentas locais do ambiente do Windows.
    tools = root / '.tools'
    tools.mkdir(exist_ok=True)
    # manifest registra as versões e os endereços exatos das ferramentas.
    manifest = tools / 'versions.json'
    if manifest.exists():
        # versions preserva os metadados fixados para instalações reproduzíveis.
        versions = json.loads(manifest.read_text(encoding='utf-8'))
    else:
        # releases reúne as versões Go publicadas pelo fornecedor.
        releases = fetch_json('https://go.dev/dl/?mode=json')
        # release seleciona a primeira versão estável disponível.
        release = next(item for item in releases if item['stable'])
        # go_file seleciona o pacote Go para Windows x64.
        go_file = next(item for item in release['files'] if item['os'] == 'windows' and item['arch'] == 'amd64' and item['kind'] == 'archive')
        # mingw reúne os metadados oficiais da distribuição MinGW.
        mingw = fetch_json('https://api.github.com/repos/brechtsanders/winlibs_mingw/releases/latest')
        # asset seleciona a distribuição x64 sem o pacote LLVM.
        asset = next(item for item in mingw['assets'] if item['name'].startswith('winlibs-x86_64-posix-seh-gcc-') and item['name'].endswith('.zip') and '-llvm-' not in item['name'])
        # versions preserva os metadados fixados para instalações reproduzíveis.
        versions = {'go': release['version'], 'goUrl': 'https://go.dev/dl/' + go_file['filename'], 'goSha256': go_file['sha256'], 'mingw': mingw['tag_name'], 'mingwUrl': asset['browser_download_url'], 'mingwSha256': (asset.get('digest') or '').removeprefix('sha256:')}
        manifest.write_text(json.dumps(versions, indent=2), encoding='utf-8')
    print(json.dumps(versions), flush=True)
    if not (tools / 'go/bin/go.exe').exists():
        install(versions['goUrl'], tools / 'go.zip', tools, versions['goSha256'])
        print('Go preparado', flush=True)
    if not (tools / 'mingw64/bin/gcc.exe').exists():
        install(versions['mingwUrl'], tools / 'mingw.zip', tools, versions['mingwSha256'] or None)
        print('MinGW preparado', flush=True)


if __name__ == '__main__':
    main()
