"""Gera dados visuais determinísticos do acervo EFL; não executa código Java."""

import hashlib
import argparse
import json
from pathlib import Path
import struct
import xml.etree.ElementTree as ET


def font_family(data):
    """Lê o nome real da família SFNT, sem deduzi-lo do nome do arquivo."""
    # table_count limita a busca à tabela de nomes declarada no arquivo.
    table_count = struct.unpack_from(">H", data, 4)[0]
    for index in range(table_count):
        # offset localiza o cabeçalho da tabela corrente.
        offset = 12 + index * 16
        tag, _, table_offset, _ = struct.unpack_from(">4sIII", data, offset)
        if tag != b"name":
            continue
        _, count, strings_offset = struct.unpack_from(">HHH", data, table_offset)
        for record in range(count):
            # Os registros Windows Unicode preservam os nomes das faces originais.
            platform, _, _, name_id, length, position = struct.unpack_from(
                ">HHHHHH", data, table_offset + 6 + record * 12
            )
            if platform == 3 and name_id == 1:
                start = table_offset + strings_offset + position
                return data[start:start + length].decode("utf-16-be")
    raise ValueError("fonte sem família Unicode")


def color(element):
    """Preserva ausência, preto opaco e alfa explicitamente zero como valores distintos."""
    if element is None:
        return dict(R=0, G=0, B=0, A=0)
    return dict(R=int(element.get("red", 0)), G=int(element.get("green", 0)),
                B=int(element.get("blue", 0)), A=int(element.get("alpha", 255)))


def font_style(data):
    """Extrai flags físicas SFNT; flags XML não implicam síntese no builder Java."""
    # index percorre o diretório SFNT até a tabela head.
    for index in range(struct.unpack_from(">H", data, 4)[0]):
        tag, _, offset, _ = struct.unpack_from(">4sIII", data, 12 + index * 16)
        if tag == b"head":
            # flags contém macStyle, que descreve a face física incorporada.
            flags = struct.unpack_from(">H", data, offset + 44)[0]
            return bool(flags & 1), bool(flags & 2)
    return False, False


def generate(check=False):
    """Produz o manifesto a partir do XML e dos bytes disponíveis, com origem e hashes."""
    # app_root e source mantêm a importação restrita aos recursos deste repositório.
    app_root = Path(__file__).resolve().parents[1]
    source = app_root.parents[1] / "material-legado/elemetos_componente/xml"
    destination = app_root / "internal/skin/assets"
    xml = (source / "skin.xml").read_bytes()
    definitions, resources, unavailable = {}, {}, []
    # element percorre também definições internas, como tableHeader.
    for element in ET.fromstring(xml).iter():
        name = element.get("name")
        if not name or element.tag == "state":
            continue
        states = {}
        missing = []
        for state in element.findall("state"):
            font = state.find("font")
            converted = dict(Name=state.get("name"), ImagePath=state.get("imagePath", ""),
                             Background=color(state.find("backgroundColor")), Font=None,
                             Spacement=float(state.get("spacement", 1)))
            if font is not None:
                converted["Font"] = dict(Type=font.get("type", ""), Size=int(font.get("size", 0)),
                                         Bold=font.get("bold") == "true", Italic=font.get("italic") == "true",
                                         Color=color(font.find("color")))
            for path in [converted["ImagePath"], font.get("type", "") if font is not None else ""]:
                if not path:
                    continue
                # O acervo incorporado já contém recursos complementares resolvidos anteriormente.
                target = destination / path.lstrip("/")
                original = source / path.lstrip("/")
                if original.is_file() and not check:
                    target.parent.mkdir(parents=True, exist_ok=True)
                    target.write_bytes(original.read_bytes())
                if not target.is_file():
                    missing.append(path)
                    continue
                data = target.read_bytes()
                # bold e italic refletem a face física, separada dos atributos XML.
                bold, italic = font_style(data) if path.endswith(".ttf") else (False, False)
                resources[path] = dict(SHA256=hashlib.sha256(data).hexdigest(), Size=len(data),
                                       Family=font_family(data) if path.endswith(".ttf") else "",
                                       Bold=bold, Italic=italic)
            states[converted["Name"]] = converted
        if missing:
            unavailable.append(dict(Type=element.tag, Name=name, Missing=sorted(set(missing))))
        else:
            definitions.setdefault(element.tag, {})[name] = dict(Type=element.tag, Name=name, States=states)
    manifest = dict(SourceSHA256=hashlib.sha256(xml).hexdigest(), Definitions=definitions,
                    Resources=resources, Unavailable=unavailable)
    # output é determinístico e pode ser conferido sem alterar o acervo incorporado.
    output = json.dumps(manifest, ensure_ascii=False, sort_keys=True, indent=2) + "\n"
    if check:
        if (destination / "library.json").read_text(encoding="utf-8") != output:
            raise SystemExit("Biblioteca visual desatualizada; execute scripts/import_legacy_skin.py.")
    else:
        (destination / "library.json").write_text(output, encoding="utf-8")
    print(f"Skins: {sum(map(len, definitions.values()))}; indisponíveis: {len(unavailable)}; recursos: {len(resources)}")


if __name__ == "__main__":
    # parser separa regeneração explícita da verificação somente leitura usada no build.
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", help="Confere o manifesto sem regravá-lo.")
    generate(check=parser.parse_args().check)
