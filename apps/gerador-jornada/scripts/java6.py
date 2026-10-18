"""Valida fontes Java 6 com ECJ no JDK 17, bytecode major 50 e assinaturas SE 6."""

import argparse
import hashlib
import json
import os
from pathlib import Path
import struct
import subprocess
import sys
import tempfile
import urllib.request
import zipfile


# ROOT mantém ferramentas e configuração dentro da aplicação, sem instalações globais.
ROOT = Path(__file__).resolve().parents[1]
# LOCK fixa fornecedores, versões e hashes das ferramentas de desenvolvimento.
LOCK = ROOT / "scripts/java6-tools.json"


def prepare_tools(offline=False):
    """Verifica hashes e baixa somente artefatos ausentes, salvo no modo offline."""
    # tools isola os binários auxiliares que não fazem parte do EXE distribuído.
    tools = ROOT / ".tools/java6"
    tools.mkdir(parents=True, exist_ok=True)
    # entries preserva o inventário versionado junto aos fontes.
    entries = json.loads(LOCK.read_text(encoding="utf-8"))
    # name e metadata identificam cada download previamente fixado.
    for name, metadata in entries.items():
        # target nunca recebe conteúdo cuja soma não corresponda ao manifesto.
        target = tools / name
        if not target.exists():
            if offline:
                raise ValueError("Ferramenta ausente no modo offline: " + name)
            # response limita a conexão de preparação a trinta segundos por tentativa.
            with urllib.request.urlopen(metadata["url"], timeout=30) as response:
                # content permanece em memória até a conferência de integridade.
                content = response.read()
            if hashlib.sha256(content).hexdigest() != metadata["sha256"]:
                raise ValueError("Hash de download divergente: " + name)
            target.write_bytes(content)
        if hashlib.sha256(target.read_bytes()).hexdigest() != metadata["sha256"]:
            raise ValueError("Ferramenta local alterada: " + str(target))
    return tools, entries


def run(command, cwd=None):
    """Executa argumentos sem shell e preserva diagnósticos em caso de falha."""
    # result mantém stdout e stderr disponíveis também nos testes negativos.
    result = subprocess.run(command, cwd=cwd, capture_output=True, text=True,
                            encoding="utf-8", errors="replace", timeout=180)
    if result.returncode:
        raise ValueError(result.stdout + result.stderr)
    return result.stdout + result.stderr


def class_major(content, name, exact=False):
    """Recusa classes malformadas ou posteriores ao Java 6; saída própria exige major 50."""
    if len(content) < 8 or content[:4] != b"\xca\xfe\xba\xbe":
        raise ValueError("Classfile inválido: " + name)
    # major corresponde à versão binária declarada no cabeçalho da classe.
    major = struct.unpack(">H", content[6:8])[0]
    if major > 50 or (exact and major != 50):
        raise ValueError("JAVA_BYTECODE_INCOMPATIBLE: %s possui major %s; alvo 50" % (name, major))
    return major


def dependency_classes(jars):
    """Valida o bytecode de todos os JARs transitivos fornecidos, sem carregar seu código."""
    # count informa no relatório quantas classes externas foram examinadas.
    count = 0
    # jar representa uma dependência explícita, sem descoberta silenciosa no ambiente.
    for jar in jars:
        # archive permite ler cabeçalhos sem executar bibliotecas externas.
        with zipfile.ZipFile(jar) as archive:
            # entry percorre inclusive classes que não sejam utilizadas diretamente pelo projeto.
            for entry in archive.infolist():
                if entry.filename.endswith(".class"):
                    # stream lê somente os oito bytes necessários para conferir a versão.
                    with archive.open(entry) as stream:
                        class_major(stream.read(8), str(jar) + "!" + entry.filename)
                    count += 1
    return count


def source_files(project):
    """Lê sources.list com caminhos locais seguros ou descobre fontes sob src/."""
    # inventory utiliza o mesmo contrato de saída do gerador Java.
    inventory = project / "sources.list"
    # sources mantém a ordem lexical e elimina entradas repetidas.
    sources = set()
    # names aceita espaços nos caminhos sem interpretá-los como opções do compilador.
    names = inventory.read_text(encoding="utf-8-sig").splitlines() if inventory.exists() else [
        str(path.relative_to(project)) for path in (project / "src").rglob("*.java")]
    # name identifica uma entrada do inventário que precisa permanecer sob o diretório do projeto.
    for name in names:
        if not name.strip():
            continue
        # path resolve também links para impedir a leitura de fontes externos ao projeto.
        path = (project / name).resolve()
        if not path.is_relative_to(project) or path.suffix != ".java" or not path.is_file():
            raise ValueError("Fonte inválido em sources.list: " + name)
        sources.add(path)
    if not sources:
        raise ValueError("Nenhum fonte Java encontrado em sources.list ou src/.")
    return sorted(sources)


def check(project, jdk_home, jars, offline=False):
    """Compila com sintaxe 1.6 e verifica APIs de classes próprias e dependências."""
    # project normaliza aliases Windows e links antes de verificar os caminhos dos fontes.
    project = project.resolve()
    # java e javac pertencem ao JDK informado, sem alterar JAVA_HOME globalmente.
    java = jdk_home / "bin/java.exe"
    javac = jdk_home / "bin/javac.exe"
    if not java.is_file() or not javac.is_file():
        raise ValueError("Informe --jdk-home com um JDK 17 completo.")
    # version permite identificar a ferramenta realmente utilizada no relatório.
    version = run([str(java), "-version"])
    if 'version "17.' not in version:
        raise ValueError("Esta ferramenta foi homologada para execução no JDK 17: " + version)
    # tools e entries preservam os artefatos fixados e sua procedência.
    tools, entries = prepare_tools(offline)
    # sources são examinados antes de qualquer compilação.
    sources = source_files(project)
    # external_count verifica também bytecode de bibliotecas EFL e transitivas.
    external_count = dependency_classes(jars)
    # temporary evita misturar classes antigas com o resultado da validação atual.
    with tempfile.TemporaryDirectory(prefix="jornada-java6-") as temporary:
        # work contém somente os artefatos temporários desta execução.
        work = Path(temporary)
        # classes recebe a saída própria, que deve ter major exatamente 50.
        classes = work / "classes"
        classes.mkdir()
        # arguments preserva caminhos com espaços e encoding UTF-8 para o ECJ.
        arguments = work / "sources.args"
        arguments.write_text("\n".join('"' + str(path).replace("\\", "/") + '"' for path in sources), encoding="utf-8")
        # command fixa as opções sem reaproveitar flags arbitrárias do projeto gerado.
        command = [str(java), "-Dfile.encoding=UTF-8", "-jar", str(tools / "ecj-3.33.0.jar"),
                   "-source", "1.6", "-target", "1.6", "-encoding", "UTF-8", "-proc:none",
                   "-g:lines,source", "-d", str(classes)]
        if jars:
            command += ["-classpath", os.pathsep.join(str(jar) for jar in jars)]
        command += ["@" + str(arguments)]
        run(command)
        # own_classes permite conferir que a compilação produziu bytecode verificável.
        own_classes = sorted(classes.rglob("*.class"))
        if not own_classes:
            raise ValueError("Compilação não produziu classes.")
        # path identifica cada classe própria, incluindo classes anônimas Swing.
        for path in own_classes:
            class_major(path.read_bytes()[:8], str(path), exact=True)
        # checker_classpath fornece somente as dependências fixadas do verificador.
        checker_classpath = os.pathsep.join(str(tools / name) for name in entries
                                           if name.endswith(".jar") and not name.startswith("ecj-"))
        # helper é uma ferramenta JDK 17; nunca integra os fontes Java 6 gerados.
        helper = ROOT / "scripts/java6/CheckJava6.java"
        run([str(javac), "-encoding", "UTF-8", "-proc:none", "-classpath", checker_classpath,
             "-d", str(work), str(helper)])
        run([str(java), "-Dfile.encoding=UTF-8", "-classpath", str(work) + os.pathsep + checker_classpath,
             "CheckJava6", str(tools / "java16-1.1.signature"), str(classes)] + [str(jar) for jar in jars])
    return {"status": "compatibilidade_estatica_aprovada", "source": "1.6", "target": "1.6",
            "classfileMajor": 50, "sources": len(sources), "classes": len(own_classes),
            "dependencyClasses": external_count, "jdk": version.strip(),
            "tools": entries, "sourcesSha256": {str(path.relative_to(project)): hashlib.sha256(path.read_bytes()).hexdigest()
                                               for path in sources},
            "dependenciesSha256": {str(jar): hashlib.sha256(jar.read_bytes()).hexdigest() for jar in jars},
            "limits": "Não executa Java 6 nem homologa bootstrap, navegação ou comportamento EFL."}


def publish_report(path, report):
    """Substitui o relatório integralmente; uma reprovação nunca preserva sucesso obsoleto."""
    # temporary permanece no mesmo volume do destino para permitir substituição atômica.
    temporary = None
    try:
        # stream recebe o JSON completo antes de publicar o novo resultado.
        with tempfile.NamedTemporaryFile(mode="w", encoding="utf-8", dir=path.parent,
                                         prefix=".java6-report-", suffix=".tmp", delete=False) as stream:
            temporary = Path(stream.name)
            stream.write(json.dumps(report, indent=2, ensure_ascii=False) + "\n")
        os.replace(temporary, path)
    finally:
        if temporary is not None:
            temporary.unlink(missing_ok=True)


def main():
    """Oferece preparação e validação explícitas, sem compilar dentro do editor desktop."""
    # parser mantém as opções descritas em português para quem prepara a homologação.
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=["prepare", "check"])
    parser.add_argument("project", nargs="?", type=Path)
    parser.add_argument("--jdk-home", type=Path, default=Path(os.environ.get("JAVA_HOME", ".")))
    parser.add_argument("--classpath", action="append", type=Path, default=[], help="JAR externo; repetir para transitivas")
    parser.add_argument("--offline", action="store_true", help="Recusar downloads; usar cache já verificado")
    parser.add_argument("--report", type=Path, help="Gravar resultado JSON; falhas também substituem o relatório e retornam código 1")
    # args representa somente argumentos locais informados pelo desenvolvedor.
    args = parser.parse_args()
    try:
        if args.command == "prepare":
            prepare_tools(args.offline)
            print("Ferramentas Java 6 preparadas e hashes conferidos.")
            return 0
        if args.project is None:
            parser.error("check exige o diretório do projeto Java")
        # report só é produzido depois de todas as verificações desta execução.
        report = check(args.project.resolve(), args.jdk_home.resolve(),
                       [jar.resolve() for jar in args.classpath], args.offline)
        # output preserva português e pode ser arquivado como evidência de validação estática.
        output = json.dumps(report, indent=2, ensure_ascii=False)
        if args.report:
            publish_report(args.report, report)
        print(output)
        return 0
    except (ValueError, OSError, subprocess.TimeoutExpired, zipfile.BadZipFile) as error:
        print("Validação Java 6 reprovada: " + str(error), file=sys.stderr)
        if args.report:
            try:
                publish_report(args.report, {"status": "reprovado", "error": str(error)})
            except OSError as report_error:
                print("Não foi possível atualizar o relatório: " + str(report_error), file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
