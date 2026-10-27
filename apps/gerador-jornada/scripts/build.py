"""Compila e empacota o editor por plataforma em dist, sem simular builds ausentes."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import platform
import plistlib
import shutil
import struct
import subprocess
import sys
import tarfile
import tempfile
import zipfile
from third_party import collect as collect_licenses

# ROOT permite chamar o script de qualquer diretório do workspace.
ROOT = Path(__file__).resolve().parents[1]
# TARGETS contém somente arquiteturas de 64 bits previstas para distribuição.
TARGETS = (("windows", "amd64"), ("linux", "amd64"), ("darwin", "amd64"), ("darwin", "arm64"))


def run(arguments, env=None, capture=False):
    """Executa argumentos sem shell, interrompendo a distribuição quando uma etapa falha."""
    return subprocess.run([str(item) for item in arguments], cwd=ROOT, env=env,
                          check=True, text=True, capture_output=capture)


def host_target():
    """Identifica o host real para impedir rotular binário de outra arquitetura como nativo."""
    # system e machine normalizam nomes diferentes usados pelos três sistemas operacionais.
    system = {"Windows": "windows", "Linux": "linux", "Darwin": "darwin"}.get(platform.system())
    machine = {"AMD64": "amd64", "x86_64": "amd64", "arm64": "arm64", "aarch64": "arm64"}.get(platform.machine())
    return system, machine


def verify_binary(path, target, arch):
    """Confere PE, ELF ou Mach-O e a arquitetura sem executar o programa gerado."""
    # data contém apenas os cabeçalhos necessários à identificação do executável.
    with path.open("rb") as source:
        data = source.read(4096)
    if len(data) < 64:
        raise ValueError("Cabeçalho de executável truncado.")
    if target == "windows":
        if data[:2] != b"MZ":
            raise ValueError("Executável Windows sem cabeçalho PE.")
        # offset localiza a assinatura PE indicada no cabeçalho DOS.
        offset = struct.unpack_from("<I", data, 60)[0]
        valid = data[offset:offset + 4] == b"PE\0\0" and struct.unpack_from("<H", data, offset + 4)[0] == 0x8664
    elif target == "linux":
        valid = data[:6] == b"\x7fELF\x02\x01" and struct.unpack_from("<H", data, 18)[0] == 62
    else:
        # cpu diferencia Intel x64 e Apple Silicon em executáveis Mach-O de 64 bits.
        cpu = 0x01000007 if arch == "amd64" else 0x0100000C
        valid = data[:4] == b"\xcf\xfa\xed\xfe" and struct.unpack_from("<I", data, 4)[0] == cpu
    if not valid:
        raise ValueError(f"Cabeçalho não corresponde a {target}/{arch} de 64 bits.")


def mac_bundle(output, version, env):
    """Monta .app com ícone, metadados nativos e assinatura ad-hoc sem credenciais Apple."""
    # contents define a estrutura padrão do pacote do macOS.
    contents = output / "JornadaStudio.app/Contents"
    (contents / "MacOS").mkdir(parents=True, exist_ok=True)
    (contents / "Resources").mkdir(parents=True, exist_ok=True)
    # binary é o único executável do bundle, sempre da arquitetura solicitada.
    binary = contents / "MacOS/JornadaStudio"
    run(["go", "build", "-trimpath", "-ldflags=-s -w", "-o", binary, "./cmd/gerador-jornada"], env)
    # metadata utiliza três partes nos campos Apple e conserva a versão completa em campo próprio.
    metadata = {"CFBundleName": "JornadaStudio", "CFBundleDisplayName": "Jornada Studio",
                "CFBundleIdentifier": "br.local.geradorjornada", "CFBundleExecutable": "JornadaStudio",
                "CFBundlePackageType": "APPL", "CFBundleShortVersionString": ".".join(version.split(".")[:3]),
                "CFBundleVersion": ".".join(version.split(".")[:3]), "JornadaStudioFullVersion": version,
                "CFBundleIconFile": "JornadaStudio.icns", "NSHighResolutionCapable": True,
                "LSMinimumSystemVersion": env["MACOSX_DEPLOYMENT_TARGET"]}
    with (contents / "Info.plist").open("wb") as destination:
        plistlib.dump(metadata, destination)
    # temporary isola as representações intermediárias do ícone sem alterar a imagem original.
    with tempfile.TemporaryDirectory(prefix="jornada-icon-") as temporary:
        # iconset é o diretório exigido pelo conversor nativo iconutil.
        iconset = Path(temporary) / "JornadaStudio.iconset"
        iconset.mkdir()
        for size in (16, 32, 128, 256, 512):
            for scale in (1, 2):
                # suffix diferencia o recurso Retina mantendo o mesmo desenho fornecido.
                suffix = "@2x" if scale == 2 else ""
                run(["sips", "-z", size * scale, size * scale,
                     ROOT / "internal/adapter/builtinassets/assets/app-icon.png", "--out",
                     iconset / f"icon_{size}x{size}{suffix}.png"], capture=True)
        run(["iconutil", "-c", "icns", iconset, "-o", contents / "Resources/JornadaStudio.icns"])
    run(["codesign", "--force", "--sign", "-", "--timestamp=none", output / "JornadaStudio.app"])
    return binary


def package_output(output, target, arch, version, binary, env):
    """Gera inventário, hashes e arquivo de distribuição somente após validar o binário."""
    verify_binary(binary, target, arch)
    collect_licenses(ROOT, output, env)
    # guide conserva o manual completo, identificando os links que pertencem ao repositório-fonte.
    guide = (ROOT / "README.md").read_text(encoding="utf-8").replace(
        'src="internal/adapter/builtinassets/assets/app-icon.png"', 'src="JornadaStudio.png"')
    (output / "README.md").write_text(
        "> Pacote de distribuição: os links de código e specs deste guia são relativos ao repositório-fonte.\n\n" + guide,
        encoding="utf-8")
    shutil.copy2(ROOT / "internal/adapter/builtinassets/assets/app-icon.png", output / "JornadaStudio.png")
    # info identifica ferramentas e assinatura sem alegar homologação manual.
    info = {"version": version, "os": target, "arch": arch, "cgo": True,
            "go": run([env.get("JORNADA_GO", "go"), "version"], env, capture=True).stdout.strip(),
            "buildHost": platform.platform(), "libc": platform.libc_ver(),
            "signing": "ad-hoc" if target == "darwin" else "unsigned",
            "manualApproval": "pending", "binary": binary.relative_to(output).as_posix()}
    # compiler identifica a ferramenta C usada pela compilação nativa CGO.
    compiler = ROOT / ".tools/mingw64/bin/gcc.exe" if target == "windows" else "cc"
    info["compiler"] = run([compiler, "--version"], env, capture=True).stdout.splitlines()[0]
    if target == "darwin":
        info["appleSDK"] = run(["xcrun", "--show-sdk-version"], env, capture=True).stdout.strip()
    elif target == "linux":
        # libraries registra dependências efetivas sem iniciar a janela ou alegar runtime estático.
        libraries = run(["ldd", binary], env, capture=True).stdout
        if "not found" in libraries:
            raise ValueError("O ambiente de build não resolve todas as bibliotecas Linux.")
        (output / "runtime-libraries.txt").write_text(libraries, encoding="utf-8")
    (output / "build-info.json").write_text(json.dumps(info, indent=2) + "\n", encoding="utf-8")
    # files limita o inventário aos arquivos da plataforma, excluindo o próprio relatório de hashes.
    files = sorted(path for path in output.rglob("*") if path.is_file() and path.name != "SHA256SUMS.txt")
    (output / "SHA256SUMS.txt").write_text("".join(
        f"{hashlib.sha256(path.read_bytes()).hexdigest()}  {path.relative_to(output).as_posix()}\n" for path in files), encoding="utf-8")
    # archive mantém o nome/versão/arquitetura visíveis antes da extração.
    archive = ROOT / "dist" / f"JornadaStudio-{version}-{target}-{arch}"
    if target == "linux":
        with tarfile.open(str(archive) + ".tar.gz", "w:gz") as package:
            package.add(output, arcname=output.name)
        archive = Path(str(archive) + ".tar.gz")
    else:
        archive = Path(str(archive) + ".zip")
        with zipfile.ZipFile(archive, "w", zipfile.ZIP_DEFLATED) as package:
            for path in sorted(output.rglob("*")):
                if path.is_file():
                    package.write(path, path.relative_to(output.parent))
    Path(str(archive) + ".sha256").write_text(
        f"{hashlib.sha256(archive.read_bytes()).hexdigest()}  {archive.name}\n", encoding="utf-8")
    print(f"Distribuição gerada: {archive}")


def docker_linux(check):
    """Compila Linux amd64 num contêiner nativo, preservando a pasta dist do workspace."""
    run(["docker", "build", "--quiet", "--platform", "linux/amd64", "-f", "scripts/Dockerfile.linux",
         "-t", "jornada-builder-linux:go1.27.2", "scripts"])
    run(["docker", "run", "--rm", "--platform", "linux/amd64", "--mount",
         f"type=bind,source={ROOT},target=/workspace", "--mount", "type=volume,source=jornada-go-mod,target=/go/pkg/mod",
         "--mount", "type=volume,source=jornada-go-build,target=/root/.cache/go-build",
         "jornada-builder-linux:go1.27.2"] + (["--check"] if check else []))


def main():
    """Valida o alvo e executa somente o build disponível no host ou no contêiner Linux."""
    # parser mantém os alvos explícitos, sem produzir arquivos vazios para plataformas ausentes.
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--target", choices=("windows", "linux", "darwin"))
    parser.add_argument("--arch", choices=("amd64", "arm64"))
    parser.add_argument("--docker", action="store_true", help="Compila Linux amd64 com Docker.")
    parser.add_argument("--check", action="store_true", help="Executa vet e testes Go antes de compilar.")
    parser.add_argument("--list", action="store_true", help="Lista plataformas e inventários disponíveis em dist.")
    args = parser.parse_args()
    if args.list:
        for target, arch in TARGETS:
            print(f"{target}/{arch}: " + ("inventário gerado" if (ROOT / f"dist/{target}-{arch}/build-info.json").exists() else "a compilar"))
        return
    # target e arch usam o host por padrão, sempre com arquitetura de 64 bits validada.
    target, arch = args.target or host_target()[0], args.arch or host_target()[1]
    if (target, arch) not in TARGETS:
        parser.error("Alvo não suportado; use --list para consultar as plataformas.")
    if args.docker:
        if (target, arch) != ("linux", "amd64"):
            parser.error("--docker suporta somente --target linux --arch amd64.")
        docker_linux(args.check)
        return
    if (target, arch) != host_target():
        parser.error("Fyne/CGO exige host/toolchain nativo. Use Docker para Linux ou o runner nativo do workflow.")
    # version é a fonte única já compartilhada pela interface e pelos metadados Windows.
    version = (ROOT / "internal/adapter/buildinfo/VERSION").read_text(encoding="utf-8").strip()
    # output isola binários, metadados e pacotes por sistema/arquitetura.
    output = ROOT / f"dist/{target}-{arch}"
    output.mkdir(parents=True, exist_ok=True)
    # env impede que variáveis de cross-compilation do terminal alterem o alvo real.
    env = dict(os.environ, GOOS=target, GOARCH=arch, CGO_ENABLED="1")
    if args.check:
        if target == "windows":
            run(["powershell", "-NoProfile", "-File", ROOT / "scripts/check.ps1"], env)
        else:
            run(["go", "vet", "./..."], env)
            run(["go", "test", "./...", "-count=1"], env)
    if target == "windows":
        run(["powershell", "-NoProfile", "-File", ROOT / "scripts/build.ps1", "-OutputDirectory", output], env)
        binary = output / "JornadaStudio.exe"
        # JORNADA_GO aponta para a mesma ferramenta usada por build.ps1, também fora do PATH.
        env["JORNADA_GO"] = str(ROOT / ".tools/go/bin/go.exe")
        try:
            shutil.copy2(binary, ROOT / "dist/JornadaStudio.exe")
        except OSError as error:
            if isinstance(error, PermissionError) or getattr(error, "winerror", None) == 32:
                print("O EXE no caminho antigo está em uso; a distribuição atualizada está em dist/windows-amd64.")
            else:
                raise
    elif target == "darwin":
        env.setdefault("MACOSX_DEPLOYMENT_TARGET", "15.0")
        binary = mac_bundle(output, version, env)
    else:
        binary = output / "JornadaStudio"
        run(["go", "build", "-trimpath", "-ldflags=-s -w", "-o", binary, "./cmd/gerador-jornada"], env)
        binary.chmod(0o755)
    package_output(output, target, arch, version, binary, env)


if __name__ == "__main__":
    try:
        main()
    except (subprocess.CalledProcessError, OSError, ValueError) as error:
        print(f"Falha na distribuição: {error}", file=sys.stderr)
        sys.exit(1)
