"""Prepara configuração pública de cliente OAuth para um build, sem registrar seus valores."""
import argparse
import json
import os
from pathlib import Path
import re

# ROOT identifica o módulo, independentemente do diretório do processo.
ROOT = Path(__file__).resolve().parents[1]
# GENERATED é o único arquivo que esta ferramenta pode criar ou remover.
GENERATED = ROOT / "internal/adapter/github/config_generated.go"
# MARKER distingue um artefato da ferramenta de um arquivo escrito manualmente.
MARKER = "// Arquivo gerado por configure_github.py; não versionar.\n"


def properties_values(path):
    """Lê somente os dois atributos autorizados de um arquivo externo ao repositório."""
    # content aceita UTF-8/BOM e o encoding tradicional de properties sem imprimir dados.
    content = Path(path).read_bytes()
    try:
        text = content.decode("utf-8-sig")
    except UnicodeDecodeError:
        text = content.decode("latin-1")
    # selected não transporta outros atributos eventualmente presentes no arquivo.
    selected = {}
    for line in text.splitlines():
        stripped = line.strip()
        if not stripped or stripped.startswith(("#", "!", ";", "[")):
            continue
        # match suporta os separadores comuns sem avaliar código ou expansão de shell.
        match = re.match(r"^([^\s=:]+)\s*(?:=|:)\s*(.*?)\s*$", stripped)
        if match and match[1] in ("CLI_ID_GITHUB_OAUTH", "SECRET_GITHUB_OAUTH"):
            if match[1] in selected:
                raise ValueError("Atributo OAuth duplicado no arquivo externo.")
            selected[match[1]] = match[2]
    return selected


def prepare(destination=GENERATED, environment=None):
    """Gera init Go a partir do par completo; ausência do par mantém o stub offline."""
    # values permite testes sem ler credenciais do processo real.
    values = os.environ if environment is None else environment
    # client_id e client_secret são do aplicativo público, nunca tokens de usuário.
    client_id = values.get("JORNADA_GITHUB_CLIENT_ID", "")
    client_secret = values.get("JORNADA_GITHUB_CLIENT_SECRET", "")
    if not client_id and not client_secret and values.get("JORNADA_GITHUB_PROPERTIES"):
        # selected recebe apenas as chaves explicitamente autorizadas pelo proprietário.
        selected = properties_values(values["JORNADA_GITHUB_PROPERTIES"])
        client_id = selected.get("CLI_ID_GITHUB_OAUTH", "")
        client_secret = selected.get("SECRET_GITHUB_OAUTH", "")
        if not client_id or not client_secret:
            raise ValueError("Arquivo externo sem o par completo de atributos OAuth.")
    if bool(client_id) != bool(client_secret):
        raise ValueError("Forneça o par completo de configuração OAuth GitHub.")
    if destination.exists():
        raise ValueError("Configuração gerada já existe; limpe o build anterior antes de continuar.")
    if not client_id:
        return False
    if not all(re.fullmatch(r"[A-Za-z0-9_-]{1,256}", value) for value in (client_id, client_secret)):
        raise ValueError("Formato de configuração OAuth inválido.")
    destination.parent.mkdir(parents=True, exist_ok=True)
    # source usa codificação de string, sem avaliar texto como shell ou código arbitrário.
    source = MARKER + "package github\nfunc init() { ClientID = " + json.dumps(client_id) + "; ClientSecret = " + json.dumps(client_secret) + " }\n"
    with destination.open("x", encoding="utf-8", newline="\n") as generated:
        generated.write(source)
    return True


def clean(destination=GENERATED):
    """Remove somente o arquivo gerado com marcador conhecido."""
    if destination.exists():
        if not destination.read_text(encoding="utf-8").startswith(MARKER):
            raise ValueError("Arquivo sem marcador de build; remoção recusada.")
        destination.unlink()


if __name__ == "__main__":
    # parser oferece operações explícitas sem aceitar caminhos arbitrários na CLI.
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--clean", action="store_true")
    arguments = parser.parse_args()
    try:
        if arguments.clean:
            clean()
        else:
            prepare()
    except (ValueError, OSError) as error:
        parser.exit(1, f"Configuração OAuth não preparada: {error}\n")
