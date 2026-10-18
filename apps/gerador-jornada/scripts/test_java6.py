"""Testes reais da toolchain ECJ/assinaturas, incluindo incompatibilidades silenciosas."""

import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import zipfile

import java6


class Java6ToolchainTest(unittest.TestCase):
    """Distingue bytecode antigo, sintaxe antiga e APIs efetivamente disponíveis no SE 6."""

    @classmethod
    def setUpClass(cls):
        """Exige um JDK explícito e prepara as dependências fixadas antes dos cenários."""
        # jdk informa o ambiente de validação sem modificar a instalação do Windows.
        cls.jdk = Path(os.environ.get("JORNADA_TEST_JDK", "C:/Desenvolvimento/jdk-17.0.11"))
        if not (cls.jdk / "bin/java.exe").exists():
            raise RuntimeError("Defina JORNADA_TEST_JDK para executar os testes reais com JDK 17.")
        java6.prepare_tools()

    def setUp(self):
        """Cria um projeto independente para cada teste e elimina os arquivos ao final."""
        # temporary impede reaproveitamento de classes entre os cenários positivos e negativos.
        self.temporary = tempfile.TemporaryDirectory(prefix="jornada-teste-java6-")
        self.addCleanup(self.temporary.cleanup)
        # project inclui espaço no caminho para exercitar o arquivo de argumentos ECJ.
        self.project = Path(self.temporary.name) / "projeto com espaço"
        (self.project / "src").mkdir(parents=True)

    def write_source(self, content):
        """Grava a classe de teste em UTF-8 com o inventário usado pelo gerador."""
        (self.project / "src/Exemplo.java").write_text(content, encoding="utf-8")
        (self.project / "sources.list").write_text("src/Exemplo.java\n", encoding="utf-8")

    def test_swing_java6_and_unicode(self):
        """Aceita Swing, genéricos Java 5/6, acentos e par substituto de emoji."""
        self.write_source('import javax.swing.JLabel; import java.util.ArrayList; '
                          'public class Exemplo { public JLabel tela() { '
                          'ArrayList<String> nomes = new ArrayList<String>(); '
                          'nomes.add("Olá 😀"); return new JLabel(nomes.get(0)); } }')
        # report comprova que tanto a compilação quanto a verificação de API foram executadas.
        report = java6.check(self.project, self.jdk, [], offline=True)
        self.assertEqual(report["classfileMajor"], 50)
        self.assertEqual(report["classes"], 1)
        self.assertEqual(report["status"], "compatibilidade_estatica_aprovada")
        self.assertEqual(len(report["sourcesSha256"]["src\\Exemplo.java"]), 64)

    def test_rejects_lambda_syntax(self):
        """ECJ deve recusar uma construção de linguagem posterior ao Java 6."""
        self.write_source('public class Exemplo { Runnable tarefa = () -> { }; }')
        with self.assertRaisesRegex(ValueError, "Lambda|lambda"):
            java6.check(self.project, self.jdk, [], offline=True)

    def test_rejects_new_static_api_despite_old_bytecode(self):
        """Sintaxe 1.6 não basta: Objects.requireNonNull só existe a partir do Java 7."""
        self.write_source('public class Exemplo { Object valor = java.util.Objects.requireNonNull("x"); }')
        with self.assertRaisesRegex(ValueError, "Undefined reference.*java.util.Objects"):
            java6.check(self.project, self.jdk, [], offline=True)

    def test_rejects_new_method_on_old_class(self):
        """Também recusa métodos modernos em classes que já existiam no Java 6."""
        self.write_source('public class Exemplo { String valor = " x ".strip(); }')
        with self.assertRaisesRegex(ValueError, "Undefined reference.*strip"):
            java6.check(self.project, self.jdk, [], offline=True)

    def make_dependency(self, source, modern=False):
        """Produz um JAR de teste sem apresentá-lo como biblioteca EFL real."""
        # source_path mantém a dependência fora do inventário do projeto principal.
        source_path = Path(self.temporary.name) / "Dependency.java"
        source_path.write_text(source, encoding="utf-8")
        # output permite inspecionar o bytecode da dependência separadamente.
        output = Path(self.temporary.name) / "dependency-classes"
        output.mkdir()
        if modern:
            java6.run([str(self.jdk / "bin/javac.exe"), "-d", str(output), str(source_path)])
        else:
            # tools contém o compilador já preparado e verificado pelo manifesto.
            tools, _ = java6.prepare_tools(True)
            java6.run([str(self.jdk / "bin/java.exe"), "-jar", str(tools / "ecj-3.33.0.jar"),
                       "-1.6", "-proc:none", "-d", str(output), str(source_path)])
        # jar representa uma biblioteca externa recebida no classpath.
        jar = Path(self.temporary.name) / "dependency.jar"
        with zipfile.ZipFile(jar, "w") as archive:
            archive.write(output / "Dependency.class", "Dependency.class")
        return jar

    def test_rejects_dependency_using_modern_api(self):
        """Não ignora referências inválidas dentro de bibliotecas com bytecode major 50."""
        # jar contém propositalmente uma chamada incompatível escondida na dependência.
        jar = self.make_dependency('public class Dependency { public String valor() { return " x ".strip(); } }')
        self.write_source('public class Exemplo { String valor = new Dependency().valor(); }')
        with self.assertRaisesRegex(ValueError, "Undefined reference.*strip"):
            java6.check(self.project, self.jdk, [jar], offline=True)

    def test_accepts_compatible_dependency(self):
        """Resolve métodos externos e sua herança sem ignorar o conteúdo do JAR."""
        # jar usa somente APIs existentes na plataforma de destino.
        jar = self.make_dependency('public class Dependency extends javax.swing.JLabel { public String valor() { return getText(); } }')
        self.write_source('public class Exemplo { String valor = new Dependency().valor(); }')
        # report registra também o inventário de classes externas realmente verificadas.
        report = java6.check(self.project, self.jdk, [jar], offline=True)
        self.assertEqual(report["dependencyClasses"], 1)

    def test_rejects_dependency_with_modern_bytecode(self):
        """Recusa um JAR compilado pelo javac 17 antes de verificar suas APIs."""
        # jar declara major 61, incompatível com a JVM de destino.
        jar = self.make_dependency('public class Dependency { }', modern=True)
        self.write_source('public class Exemplo { Dependency valor; }')
        with self.assertRaisesRegex(ValueError, "JAVA_BYTECODE_INCOMPATIBLE"):
            java6.check(self.project, self.jdk, [jar], offline=True)

    def test_rejects_external_source(self):
        """sources.list não pode escapar do diretório de fontes entregue."""
        (self.project / "sources.list").write_text("../externo.java\n", encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "Fonte inválido"):
            java6.check(self.project, self.jdk, [], offline=True)

    def test_failure_replaces_previous_report(self):
        """Uma execução reprovada não pode deixar um relatório anterior com status aprovado."""
        # report contém propositalmente uma aprovação obsoleta.
        report = self.project / "compatibilidade.json"
        report.write_text('{"status":"compatibilidade_estatica_aprovada"}', encoding="utf-8")
        self.write_source('public class Exemplo { String valor = "x".strip(); }')
        # arguments reproduz a chamada CLI com relatório persistido.
        arguments = ["java6.py", "check", str(self.project), "--jdk-home", str(self.jdk),
                     "--offline", "--report", str(report)]
        with patch("sys.argv", arguments):
            self.assertEqual(java6.main(), 1)
        self.assertEqual(json.loads(report.read_text(encoding="utf-8"))["status"], "reprovado")


if __name__ == "__main__":
    unittest.main()
