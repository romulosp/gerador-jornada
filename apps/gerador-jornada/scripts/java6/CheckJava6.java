import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.util.Collections;
import org.codehaus.mojo.animal_sniffer.SignatureBuilder;
import org.codehaus.mojo.animal_sniffer.SignatureChecker;
import org.codehaus.mojo.animal_sniffer.logging.PrintWriterLogger;

/** Verifica APIs SE 6 dos fontes compilados e dos JARs; ferramenta executada apenas no JDK 17. */
public final class CheckJava6 {
    /** Impede a criação de instâncias de uma ferramenta sem estado. */
    private CheckJava6() { }

    /** Recebe assinatura SE 6, diretório de classes e JARs externos; falhas encerram com código 1. */
    public static void main(String[] args) throws Exception {
        // logger preserva os símbolos e arquivos responsáveis por incompatibilidades.
        var logger = new PrintWriterLogger(System.err);
        // platform contém exclusivamente as classes e assinaturas públicas de Java SE 6.
        var platform = SignatureChecker.loadClasses(new FileInputStream(args[0]));
        // bytes recebe as assinaturas das classes próprias e das dependências reais.
        var bytes = new ByteArrayOutputStream();
        // builder permite resolver herança e métodos das dependências sem ignorar suas referências SE.
        var builder = new SignatureBuilder(bytes, logger);
        // index percorre o diretório compilado e cada JAR transitivo informado.
        for (int index = 1; index < args.length; index++) {
            builder.process(new File(args[index]));
        }
        builder.close();
        // application contém somente os símbolos provenientes dos artefatos em validação.
        var application = SignatureChecker.loadClasses(new ByteArrayInputStream(bytes.toByteArray()));
        // name detecta substituições indevidas de classes da plataforma.
        for (var name : application.keySet()) {
            if (platform.containsKey(name)) {
                throw new IllegalArgumentException("Classe da plataforma redefinida: " + name);
            }
        }
        platform.putAll(application);
        // checker verifica inclusive dependências; nenhum pacote é dispensado da verificação.
        var checker = new SignatureChecker(platform, Collections.emptySet(), logger);
        checker.setAnnotationTypes(Collections.emptyList());
        checker.setSourcePath(Collections.emptyList());
        // index inclui os JARs no exame de chamadas a APIs posteriores ao Java 6.
        for (int index = 1; index < args.length; index++) {
            checker.process(new File(args[index]));
        }
        if (checker.isSignatureBroken()) {
            System.exit(1);
        }
    }
}
