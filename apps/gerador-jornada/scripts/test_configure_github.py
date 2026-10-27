"""Confere geração isolada, par completo e limpeza sem credenciais reais."""
from pathlib import Path
import tempfile
import unittest
from configure_github import prepare, clean


class GitHubConfigurationTests(unittest.TestCase):
    """Valida configuração ausente, completa, inválida e arquivo alheio."""

    def test_profile_and_clean(self):
        """Gera somente com par completo e remove apenas seu próprio artefato."""
        with tempfile.TemporaryDirectory() as directory:
            # destination isola o artefato de qualquer configuração real do build.
            destination = Path(directory) / "config_generated.go"
            self.assertFalse(prepare(destination, {}))
            with self.assertRaises(ValueError):
                prepare(destination, {"JORNADA_GITHUB_CLIENT_ID": "fixture"})
            self.assertTrue(prepare(destination, {"JORNADA_GITHUB_CLIENT_ID": "fixture", "JORNADA_GITHUB_CLIENT_SECRET": "synthetic"}))
            self.assertIn('ClientID = "fixture"', destination.read_text())
            clean(destination)
            self.assertFalse(destination.exists())
            destination.write_text("arquivo manual", encoding="utf-8")
            with self.assertRaises(ValueError):
                clean(destination)

    def test_external_properties(self):
        """Consome somente os atributos permitidos sem copiar o arquivo externo."""
        with tempfile.TemporaryDirectory() as directory:
            # properties e generated são fixtures, não credenciais reais.
            properties = Path(directory) / "external.properties"
            generated = Path(directory) / "config_generated.go"
            properties.write_text("# teste\nCLI_ID_GITHUB_OAUTH=fixture\nSECRET_GITHUB_OAUTH=synthetic\nOTHER_KEY=ignored\n", encoding="utf-8")
            self.assertTrue(prepare(generated, {"JORNADA_GITHUB_PROPERTIES": str(properties)}))
            self.assertNotIn("OTHER_KEY", generated.read_text(encoding="utf-8"))
            clean(generated)


if __name__ == "__main__":
    unittest.main()
