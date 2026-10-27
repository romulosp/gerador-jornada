package desktop

import (
	"context"

	"gerador-jornada.local/editor/internal/adapter/workspace"
	"gerador-jornada.local/editor/internal/application"
)

// Workspace conserva a fachada manual sobre o armazenamento compartilhado com MCP.
type Workspace struct {
	// Workspace fornece os métodos públicos sem replicar regras de armazenamento.
	*workspace.Workspace
}

// WorkspaceEntry conserva a identidade e o diagnóstico usados pela árvore manual.
type WorkspaceEntry = workspace.WorkspaceEntry

// OpenWorkspace seleciona uma pasta existente sem criar configuração de negócio.
func OpenWorkspace(root string) (*Workspace, error) {
	// store valida a raiz física usando o mesmo adaptador consumido pelo MCP.
	store, err := workspace.OpenWorkspace(root)
	if err != nil {
		return nil, err
	}
	return &Workspace{Workspace: store}, nil
}

// JourneyFilename compartilha a normalização do nome físico com os demais canais.
func JourneyFilename(name string) (string, error) { return workspace.JourneyFilename(name) }

// NewJourneyDocument mantém os defaults e o fundo legado nos dois modos do launcher.
func NewJourneyDocument(name string, example bool) application.Document {
	return workspace.NewJourneyDocument(name, example)
}

// uniqueName preserva chamadas internas manuais; a publicação continua exclusiva.
func (w *Workspace) uniqueName(name string) (string, error) { return w.UniqueName(name) }

// publishNew preserva chamadas internas manuais sobre a publicação ancorada e exclusiva.
func (w *Workspace) publishNew(name string, data []byte) error {
	return w.PublishNew(context.Background(), name, data)
}
