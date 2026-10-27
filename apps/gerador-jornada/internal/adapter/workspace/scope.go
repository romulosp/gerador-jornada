package workspace

import (
	"errors"
	"path/filepath"

	"gerador-jornada.local/editor/internal/adapter/workspaceguard"
	"gerador-jornada.local/editor/internal/project"
)

// MaterialPurpose impede confundir evidência legada e pacote importável.
type MaterialPurpose string

const (
	// LegacyEvidence permite somente evidências para análise externa, sem execução.
	LegacyEvidence MaterialPurpose = "legacy"
	// LibraryImport permite somente materiais de entrada da biblioteca.
	LibraryImport MaterialPurpose = "import"
)

// MaterialRoot descreve uma autorização de startup sem divulgar o caminho absoluto.
type MaterialRoot struct {
	// ID é opaco, local ao processo, e não pode conceder uma raiz nova.
	ID string
	// Name é apenas referência de exibição; nunca é reinterpretado como caminho.
	Name string
	// Purpose restringe as operações permitidas sobre os materiais dessa raiz.
	Purpose MaterialPurpose
}

// materialRootState conserva identidade/caminho exclusivamente no adaptador.
type materialRootState struct {
	// info pode ser exposta; path/identity permanecem fixos desde o startup.
	info           MaterialRoot
	path, identity string
}

// MaterialScopes é imutável depois da composição; ferramentas não possuem método de concessão.
type MaterialScopes struct {
	// guard abre apenas handles de leitura; roots conserva a ordem de configuração.
	guard *workspaceguard.Guard
	roots []materialRootState
}

// NewMaterialScopes valida raízes concedidas pelo operador, sem inventariar nem gravar arquivos.
func NewMaterialScopes(legacyRoots, importRoots []string) (*MaterialScopes, error) {
	// guard resolve o namespace sem criar arquivos em uma operação de leitura.
	guard, err := workspaceguard.NewDefault()
	if err != nil {
		return nil, err
	}
	scopes := &MaterialScopes{guard: guard}
	// group conserva finalidades separadas mesmo quando o operador escolhe a mesma pasta.
	for _, group := range []struct {
		purpose MaterialPurpose
		paths   []string
	}{{LegacyEvidence, legacyRoots}, {LibraryImport, importRoots}} {
		for _, path := range group.paths {
			// lease captura a identidade física e é fechado sem reter bloqueio ou preferências.
			lease, err := guard.Open(path)
			if err != nil {
				return nil, err
			}
			scopes.roots = append(scopes.roots, materialRootState{info: MaterialRoot{ID: project.NewID(), Name: filepath.Base(lease.Root), Purpose: group.purpose}, path: lease.Root, identity: lease.Identity})
			if err := lease.Close(); err != nil {
				return nil, err
			}
		}
	}
	return scopes, nil
}

// Roots retorna cópia das autorizações da finalidade, sem aceitar caminhos do agente.
func (scopes *MaterialScopes) Roots(purpose MaterialPurpose) []MaterialRoot {
	// result não compartilha a configuração imutável do adaptador.
	result := []MaterialRoot{}
	for _, root := range scopes.roots {
		if root.info.Purpose == purpose {
			result = append(result, root.info)
		}
	}
	return result
}

// Open oferece a raiz já autorizada somente à finalidade correta e revalida sua identidade.
// O consumidor deve fechar o lease e limitar sua leitura; ele não concede permissão de escrita.
func (scopes *MaterialScopes) Open(purpose MaterialPurpose, rootID string) (*workspaceguard.Lease, error) {
	for _, root := range scopes.roots {
		if root.info.ID != rootID || root.info.Purpose != purpose {
			continue
		}
		// lease impede que a troca posterior de um caminho altere a autorização original.
		lease, err := scopes.guard.Open(root.path)
		if err != nil {
			return nil, err
		}
		if lease.Identity != root.identity {
			lease.Close()
			return nil, errors.New("raiz autorizada substituída")
		}
		return lease, nil
	}
	return nil, errors.New("raiz não autorizada para esta finalidade")
}
