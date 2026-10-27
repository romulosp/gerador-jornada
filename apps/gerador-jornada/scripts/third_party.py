"""Reúne licenças de módulos fixados no pacote sem copiar configurações ou credenciais."""
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import stat


def collect(root, output, environment):
    """Copia documentos legais dos módulos resolvidos e registra somente nome/versão/arquivos."""
    # go usa a mesma toolchain do inventário e não instala ferramentas no usuário final.
    go = environment.get("JORNADA_GO", "go")
    # result contém diretórios do cache apenas durante a montagem do pacote.
    result = subprocess.run([go, "list", "-m", "-json", "all"], cwd=root, env=environment,
                            check=True, text=True, capture_output=True)
    # decoder percorre os objetos JSON concatenados de go list.
    decoder = json.JSONDecoder()
    remaining = result.stdout.lstrip()
    inventory = []
    destination = output / "licenses"
    destination.mkdir(parents=True, exist_ok=True)
    while remaining:
        # module/offset identificam um módulo resolvido sem publicar seu diretório local.
        module, offset = decoder.raw_decode(remaining)
        remaining = remaining[offset:].lstrip()
        if module.get("Main") or not module.get("Dir"):
            continue
        # key mantém os caminhos da distribuição curtos e estáveis no Windows.
        key = hashlib.sha256((module["Path"] + "@" + module.get("Version", "")).encode()).hexdigest()[:20]
        folder = Path(module["Dir"])
        names = []
        for source in sorted(folder.iterdir()):
            if source.is_file() and source.name.upper().startswith(("LICENSE", "LICENCE", "COPYING", "NOTICE")):
                (destination / key).mkdir(exist_ok=True)
                # target não conserva read-only do cache de módulos, permitindo recompilar no Windows.
                target = destination / key / source.name
                if target.exists():
                    target.chmod(stat.S_IREAD | stat.S_IWRITE)
                shutil.copy2(source, target)
                target.chmod(stat.S_IREAD | stat.S_IWRITE)
                names.append(key + "/" + source.name)
        inventory.append({"module": module["Path"], "version": module.get("Version", ""), "files": names})
    # octicons é recurso vetorial estático, com licença distribuída junto ao ativo.
    icon_license = root / "internal/adapter/desktop/assets/github-icon-LICENSE.txt"
    if icon_license.exists():
        shutil.copy2(icon_license, destination / "github-octicons-LICENSE.txt")
        inventory.append({"module": "github.com/primer/octicons", "version": "static-2026-10-24", "files": ["github-octicons-LICENSE.txt"]})
    (destination / "index.json").write_text(json.dumps(inventory, indent=2) + "\n", encoding="utf-8")
