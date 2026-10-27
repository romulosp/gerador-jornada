package gitworkspace

import (
	"gerador-jornada.local/editor/internal/application"
	git "github.com/go-git/go-git/v5"
)

// LocalChanges verifica jornadas alteradas sem rede ou credencial, mantendo publicação como fato histórico.
func (e *Engine) LocalChanges(root string) ([]application.GitChange, error) {
	// binding conserva a raiz e o destino previamente confirmados.
	binding, err := e.Binding(root)
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(binding.Root)
	if err != nil {
		return nil, err
	}
	oid, err := head(repo)
	if err != nil {
		return nil, err
	}
	files, err := loadTree(repo, oid)
	if err != nil {
		return nil, err
	}
	_, changes, err := working(binding.Root, files)
	return changes, err
}
