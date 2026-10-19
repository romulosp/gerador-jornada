"""Verifica que a distribuição recusa binários incompatíveis antes de criar pacotes."""
from pathlib import Path
import struct
import tempfile
import unittest

from build import verify_binary


class DistributionBinaryTest(unittest.TestCase):
    """Protege o contrato de 64 bits nos três formatos nativos suportados."""

    def check_header(self, data, target, arch, valid):
        """Usa cabeçalhos de fixture para testar a identificação sem executar binários."""
        with tempfile.TemporaryDirectory(prefix="jornada-header-") as directory:
            # binary é uma fixture isolada, nunca um arquivo real de distribuição.
            binary = Path(directory) / "fixture"
            binary.write_bytes(data)
            if valid:
                verify_binary(binary, target, arch)
            else:
                with self.assertRaises((ValueError, struct.error)):
                    verify_binary(binary, target, arch)

    def test_windows_requires_amd64_pe(self):
        """Recusa PE de 32 bits e confirma a máquina AMD64 no mesmo contêiner PE."""
        # data simula um PE com endereço de assinatura independente do tamanho padrão.
        data = bytearray(512)
        data[:2] = b"MZ"
        struct.pack_into("<I", data, 60, 128)
        data[128:132] = b"PE\0\0"
        struct.pack_into("<H", data, 132, 0x14C)
        self.check_header(data, "windows", "amd64", False)
        struct.pack_into("<H", data, 132, 0x8664)
        self.check_header(data, "windows", "amd64", True)

    def test_linux_rejects_32bits_and_other_cpu(self):
        """Exige ELF64 little-endian x86-64, impedindo renomear outro alvo como Linux x64."""
        # data alterna classe e CPU mantendo a assinatura ELF válida.
        data = bytearray(64)
        data[:6] = b"\x7fELF\x02\x01"
        struct.pack_into("<H", data, 18, 62)
        self.check_header(data, "linux", "amd64", True)
        data[4] = 1
        self.check_header(data, "linux", "amd64", False)
        data[4] = 2
        struct.pack_into("<H", data, 18, 183)
        self.check_header(data, "linux", "amd64", False)

    def test_macos_distinguishes_intel_and_apple_silicon(self):
        """Diferencia os dois formatos Mach-O de 64 bits e recusa arquivos truncados."""
        # data fornece um Mach-O de arquitetura única, sem alegar validação de bundle Apple.
        data = bytearray(64)
        data[:4] = b"\xcf\xfa\xed\xfe"
        struct.pack_into("<I", data, 4, 0x01000007)
        self.check_header(data, "darwin", "amd64", True)
        self.check_header(data, "darwin", "arm64", False)
        struct.pack_into("<I", data, 4, 0x0100000C)
        self.check_header(data, "darwin", "arm64", True)
        self.check_header(b"bad", "darwin", "arm64", False)


if __name__ == "__main__":
    unittest.main()
