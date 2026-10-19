package javaefl

import (
	"fmt"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// Sources monta fontes pequenos por tela e contrato de navegação.
// Não publica nem afirma homologação; Export aplica o gate de evidência real.
func Sources(document application.Document, profile Profile) (map[string]string, error) {
	if err := document.ValidateFileIDs(); err != nil {
		return nil, err
	}
	// packageName é um namespace validado antes da geração.
	packageName := document.JavaOptions.PackageName
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := ValidatePackageName(packageName); err != nil {
		return nil, err
	}
	if profile.PanelType == "" || profile.HostSource == "" {
		return nil, fmt.Errorf("JAVA_PROFILE_UNVERIFIED: contrato de integração incompleto")
	}
	// files é o conjunto de saída ainda não publicado.
	files := map[string]string{}
	// base é o caminho do pacote definido pelo usuário.
	base := "src/" + strings.ReplaceAll(packageName, ".", "/") + "/"
	// names detecta colisões também em sistemas Windows sem distinção de caixa.
	names := map[string]bool{}
	// sourceList mantém o inventário de compilação.
	sourceList := []string{}
	// page fornece uma classe e fábrica próprias.
	for _, page := range document.Screens {
		// name é o nome normalizado sem renomear a entidade de origem.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		_, err := Identifier(page.FileID, true)
		if err != nil {
			return nil, err
		}
		// name preserva exatamente o ID editável como nome de classe e arquivo da tela.
		name := page.FileID
		if names[strings.ToLower(name)] {
			return nil, fmt.Errorf("JAVA_NAME_COLLISION: %s", page.ID)
		}
		names[strings.ToLower(name)] = true
		// declarations e initialization formam um painel coeso.
		declarations := ""
		// initialization acumula a inicialização do componente no fonte Java.
		initialization := ""
		// factoryMethods mantém um método de construção por componente.
		factoryMethods := ""
		// factoryImports contém somente APIs comprovadas pelo perfil.
		factoryImports := []string{packageName + ".integration.HostBindings"}
		// fields detecta colisões locais após normalização.
		fields := map[string]bool{}
		// component é uma instância do painel.
		for _, component := range page.Components {
			// variable identifica o campo Java da instância.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			_, err := Identifier(component.FileID, false)
			if err != nil {
				return nil, err
			}
			// variable evita colisões com parâmetros host/navigator e a fábrica local do painel.
			// O prefixo técnico não modifica o nome-base confirmado do arquivo individual.
			variable := "component" + component.FileID
			if fields[strings.ToLower(variable)] {
				return nil, fmt.Errorf("JAVA_NAME_COLLISION: %s", component.ID)
			}
			fields[strings.ToLower(variable)] = true
			// factory contém a API real de criação desta família.
			// ok informa se a consulta encontrou um valor válido.
			factory, ok := profile.Factories[component.Kind]
			if !ok || factory.Create == nil {
				return nil, fmt.Errorf("JAVA_PROPERTY_UNSUPPORTED: %s", component.Kind)
			}
			// construction é produzido exclusivamente por código de perfil homologado.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			construction, err := factory.Create(component)
			if err != nil {
				return nil, err
			}
			factoryImports = append(factoryImports, factory.Imports...)
			declarations += fmt.Sprintf("    /** %s. */\n    private final %s %s;\n", Comment(component.Name), factory.Type, variable)
			initialization += fmt.Sprintf("        %s = factory.create%s();\n        add(%s);\n", variable, component.FileID, variable)
			if component.Kind == catalog.Button {
				// target resolve o destino da operação na coleção ou no mapa corrente.
				if target := document.Journey.Destination(component.ID); target != "" {
					if factory.Click == nil {
						return nil, fmt.Errorf("JAVA_PROPERTY_UNSUPPORTED: listener não homologado")
					}
					initialization += factory.Click(variable, Literal(target))
				}
			}
			// componentPackage separa arquivos de telas diferentes que reutilizam o mesmo ID local.
			componentPackage := packageName + ".component." + strings.ToLower(name)
			// imported impede que a classe individual esconda um tipo simples exigido pelo perfil.
			for _, imported := range factory.Imports {
				if imported[strings.LastIndex(imported, ".")+1:] == component.FileID {
					return nil, fmt.Errorf("JAVA_NAME_COLLISION: ID %s conflita com import do perfil", component.FileID)
				}
			}
			files[base+"component/"+strings.ToLower(name)+"/"+component.FileID+".java"] = componentSource(componentPackage, packageName, component, factory, construction)
			factoryMethods += fmt.Sprintf("    /** Cria %s por sua fábrica individual. */\n    public %s create%s() {\n        return %s.%s.create(host);\n    }\n", Comment(component.Name), factory.Type, component.FileID, componentPackage, component.FileID)
		}
		// panelImports reúne o painel, fábricas e controles tipados.
		panelImports := append([]string{}, profile.PanelImports...)
		panelImports = append(panelImports, factoryImports...)
		panelImports = append(panelImports, packageName+".component."+name+"Factory", packageName+".integration.HostBindings", packageName+".navigation.Navigator")
		// panel contém somente montagem visual e delegação de navegação.
		panel := fmt.Sprintf("package %s.screen;\n\n%s\n\n/** Painel %s gerado a partir da jornada. */\npublic final class %s extends %s {\n%s\n    /** Monta a tela com o contexto real do hospedeiro e seu navegador. */\n    public %s(HostBindings host, final Navigator navigator) {\n        super(%s);\n        setLayout(null);\n        // factory encapsula a construção tipada desta tela.\n        %sFactory factory = new %sFactory(host);\n%s    }\n}\n", packageName, imports(panelImports), Comment(page.Name), name, profile.PanelType, declarations, name, profile.PanelSuper, name, name, initialization)
		// factorySource não concentra controles de outras telas.
		factorySource := fmt.Sprintf("package %s.component;\n\n%s\n\n/** Constrói os componentes da tela %s. */\npublic final class %sFactory {\n    /** host fornece dependências EFL do sistema de destino. */\n    private final HostBindings host;\n    /** Recebe o contexto homologado do sistema hospedeiro. */\n    public %sFactory(HostBindings host) { this.host = host; }\n%s}\n", packageName, imports(factoryImports), Comment(page.Name), name, name, factoryMethods)
		files[base+"screen/"+name+".java"] = panel
		files[base+"component/"+name+"Factory.java"] = factorySource
	}
	files[base+"integration/HostBindings.java"] = "package " + packageName + ".integration;\n" + profile.HostSource
	files[base+"navigation/Navigator.java"] = "package " + packageName + ".navigation;\n\n/** Recebe destinos por ID e mantém a navegação na EDT. */\npublic interface Navigator {\n    /** Solicita a tela identificada na thread de eventos Swing. */\n    void show(String screenId);\n}\n"
	// name percorre todos os fontes para o inventário determinístico.
	for name := range files {
		if strings.HasSuffix(name, ".java") {
			sourceList = append(sourceList, name)
		}
	}
	sort.Strings(sourceList)
	files["sources.list"] = strings.Join(sourceList, "\n") + "\n"
	files["README.md"] = "# Fontes Java 1.6 / EFL\n\nPerfil: " + profile.ID + "\n\n" +
		"O alvo permanece Java SE 6: sintaxe 1.6, APIs SE 6 e classfile major 50. O javac do JDK 17 não aceita esse alvo. " +
		"No repositório do editor, scripts/java6.py executa ECJ 3.33.0 no JDK 17 e Animal Sniffer com assinaturas java16:1.1.\n\n" +
		"```text\npython <repositorio-do-editor>/apps/gerador-jornada/scripts/java6.py check . --jdk-home <diretorio-do-jdk-17> --classpath <efl.jar> --report compatibilidade.json\n```\n\n" +
		"Repetir --classpath para todos os JARs transitivos. A primeira preparação baixa ferramentas com hashes fixados; depois usar --offline. " +
		"O verificador compila em UTF-8, confere bytecode e APIs inclusive das dependências. O relatório estático não comprova execução em Java 6 nem comportamento EFL. " +
		"Integrar HostBindings e Navigator ao bootstrap do sistema real; executar alterações Swing na EDT e homologar no ambiente Java 6/EFL de destino.\n"
	files["dependencies.md"] = "# Dependências\n\nPerfil: " + profile.ID + "\n\nEvidência: " + profile.Evidence + "\n\nOs JARs proprietários não são redistribuídos pelo gerador.\n"
	// asset inclui somente recursos declarados no projeto.
	for _, asset := range document.Project.Assets {
		files["resources/"+asset.Path] = string(document.Resources[asset.Path])
	}
	return files, nil
}
