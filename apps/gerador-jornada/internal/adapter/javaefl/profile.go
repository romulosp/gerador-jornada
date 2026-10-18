// Package javaefl gera fontes somente com contratos EFL explicitamente homologados.
package javaefl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// Factory descreve código de construção cuja API foi comprovada no perfil.
type Factory struct {
	// Type é o tipo concreto do componente EFL.
	Type string
	// Imports contém as classes necessárias à construção.
	Imports []string
	// Create emite a construção e configuração após validação do componente.
	Create func(screen.Component) (string, error)
	// Click emite o listener de navegação com a assinatura comprovada.
	Click func(variable, destination string) string
}

// Profile pertence ao build do editor, nunca ao arquivo de projeto.
type Profile struct {
	// ID identifica a versão do ambiente EFL.
	ID string
	// Verified exige evidência externa antes da publicação de fontes.
	Verified bool
	// Evidence identifica a homologação real da biblioteca.
	Evidence string
	// PanelType informa a classe de painel.
	PanelType string
	// PanelImports informa imports de painel e skin.
	PanelImports []string
	// PanelSuper informa os argumentos comprovados do construtor base.
	PanelSuper string
	// HostSource fornece o contrato de integração verificado com o hospedeiro.
	HostSource string
	// Factories mapeia todas as famílias do catálogo.
	Factories map[catalog.Kind]Factory
}

// CurrentProfile conserva o bloqueio explícito enquanto não há JARs e bootstrap homologados.
func CurrentProfile() Profile {
	return Profile{ID: "efl-legado-pendente", Factories: map[catalog.Kind]Factory{}}
}

// Export publica um diretório novo somente após validar todo o projeto e perfil.
func Export(ctx context.Context, destination string, document application.Document, profile Profile) error {
	if !profile.Verified || profile.Evidence == "" {
		return fmt.Errorf("JAVA_PROFILE_UNVERIFIED: forneça os JARs EFL e um exemplo funcional do sistema hospedeiro para homologar as APIs")
	}
	// problems reúne os diagnósticos que impedem execução ou exportação.
	if problems := document.Validate(); len(problems) > 0 {
		return fmt.Errorf("%s", application.ProblemsText(problems))
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := ValidatePackageName(document.JavaOptions.PackageName); err != nil {
		return err
	}
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		return fmt.Errorf("escolha um diretório de saída novo")
	}
	// files contém todos os fontes antes de criar a área temporária.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	files, err := Sources(document, profile)
	if err != nil {
		return err
	}
	// staging é irmão do destino e nunca mescla fontes existentes.
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	staging, err := os.MkdirTemp(filepath.Dir(destination), ".java-jornada-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	// names mantém publicação e inventário determinísticos.
	names := []string{}
	// name identifica uma entrada relativa do pacote ou arquivo gerado.
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	// name identifica um arquivo já validado do pacote.
	for _, name := range names {
		if err = ctx.Err(); err != nil {
			return err
		}
		// path permanece dentro da área construída pelo próprio gerador.
		path := filepath.Join(staging, filepath.FromSlash(name))
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err = os.WriteFile(path, []byte(files[name]), 0644); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return os.Rename(staging, destination)
}

// imports ordena e elimina classes repetidas de um arquivo Java.
func imports(values []string) string {
	// unique contém apenas declarações conhecidas pelo perfil.
	unique := map[string]bool{}
	// value é o nome completo de uma classe.
	for _, value := range values {
		if value != "" {
			unique[value] = true
		}
	}
	// sorted mantém a saída estável entre execuções.
	sorted := []string{}
	// value contém o valor corrente em validação, serialização ou edição.
	for value := range unique {
		sorted = append(sorted, value)
	}
	sort.Strings(sorted)
	// output acumula instruções de importação Java.
	output := []string{}
	// value contém o valor corrente em validação, serialização ou edição.
	for _, value := range sorted {
		output = append(output, "import "+value+";")
	}
	return strings.Join(output, "\n")
}
