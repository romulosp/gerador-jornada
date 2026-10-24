package desktop

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
)

// RecentProject é uma preferência local, independente do documento e do formato .jornada.
type RecentProject struct {
	// Path identifica a jornada principal ou pasta escolhida pelo autor.
	Path string
	// Name conserva o nome de exibição, sem substituir identidade por nome de arquivo.
	Name string
	// OpenedAt ordena a última abertura confirmada em horário absoluto.
	OpenedAt time.Time
}

// Available consulta a existência atual sem apagar entradas movidas ou indisponíveis.
func (r RecentProject) Available() bool { _, err := os.Stat(r.Path); return err == nil }

// recentFile resolve a preferência por usuário; a variável opcional isola testes e instalações portáteis.
func recentFile() (string, error) {
	// directory evita escrever configurações no workspace ou no diretório do executável.
	directory := os.Getenv("JORNADA_CONFIG_DIR")
	if directory == "" {
		// userDirectory é fornecido pelo sistema operacional, sem rede ou telemetria.
		userDirectory, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		directory = filepath.Join(userDirectory, "JornadaStudio")
	}
	return filepath.Join(directory, "recent-projects.json"), nil
}

// loadRecents lê a preferência sem transformar ausência de arquivo em falha do launcher.
func loadRecents() ([]RecentProject, error) {
	// path e data identificam apenas a preferência local do usuário.
	path, err := recentFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []RecentProject{}, nil
	}
	if err != nil {
		return nil, err
	}
	// entries mantém inclusive referências indisponíveis para diagnóstico individual.
	entries := []RecentProject{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// rememberProject move uma abertura confirmada para o topo e limita o histórico a vinte entradas.
func rememberProject(path, name string) error {
	// absolute normaliza a identidade local sem alterar o arquivo de projeto.
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// entries conserva os demais recentes; arquivo corrompido é diagnosticado, nunca sobrescrito.
	entries, err := loadRecents()
	if err != nil {
		return err
	}
	// next evita duplicatas de caixa no Windows e mantém ordem de última abertura.
	next := []RecentProject{{Path: absolute, Name: name, OpenedAt: time.Now().UTC()}}
	for _, entry := range entries {
		if !strings.EqualFold(filepath.Clean(entry.Path), absolute) && len(next) < 20 {
			next = append(next, entry)
		}
	}
	// data é publicada atomicamente somente após a serialização completa.
	data, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	config, err := recentFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(config), 0755); err != nil {
		return err
	}
	return atomicfile.Write(context.Background(), config, data)
}
