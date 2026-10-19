package javaefl

import (
	"fmt"

	"gerador-jornada.local/editor/internal/screen"
)

// componentSource gera um arquivo com o ID exato, isolando a construção da instância do painel.
func componentSource(packageName, rootPackage string, component screen.Component, factory Factory, construction string) string {
	return fmt.Sprintf("package %s;\n\n%s\n\n/** Constrói %s com o perfil de integração selecionado. */\npublic final class %s {\n    /** Impede instâncias de uma fábrica que oferece somente criação estática. */\n    private %s() { }\n    /** Cria o controle com geometria em pixels lógicos. */\n    public static %s create(%s.integration.HostBindings host) {\n        // component é uma instância independente da família EFL.\n        %s component = %s;\n        component.setBounds(%d, %d, %d, %d);\n        return component;\n    }\n}\n", packageName, imports(factory.Imports), Comment(component.Name), component.FileID, component.FileID, factory.Type, rootPackage, factory.Type, construction, component.Bounds.X, component.Bounds.Y, component.Bounds.Width, component.Bounds.Height)
}
