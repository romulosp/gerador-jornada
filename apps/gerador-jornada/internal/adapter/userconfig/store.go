// Package userconfig persiste preferências privadas sem modificar documentos ou repositórios.
package userconfig

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gerador-jornada.local/editor/internal/adapter/atomicfile"
)

// Store mantém dados não secretos na configuração do usuário.
type Store struct {
	// Root isola a instalação e os testes.
	Root string
}

// New resolve o mesmo diretório utilizado pela lista de recentes.
// Retorna o armazenamento sem criar diretórios, ou erro ao resolver seu caminho absoluto.
func New() (*Store, error) {
	// root privilegia a configuração portátil explícita.
	root := os.Getenv("JORNADA_CONFIG_DIR")
	if root == "" {
		// base é obtida pelo sistema operacional, sem inferir o workspace.
		base, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		root = filepath.Join(base, "JornadaStudio")
	}
	// absolute impede que o diretório mude com o processo.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Store{Root: absolute}, nil
}

// Key gera uma chave estável para caminhos canônicos Windows, sem expô-los no nome.
// value deve ser o caminho canônico confirmado pelo adaptador consumidor.
func Key(value string) string {
	// digest usa comparação sem distinção de caixa do workspace Windows.
	digest := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(value))))
	return hex.EncodeToString(digest[:])
}

// path rejeita nomes que escapariam da configuração privada.
// name é relativo à configuração e o retorno é o destino completo ou erro de validação.
func (s *Store) path(name string) (string, error) {
	if filepath.IsAbs(name) || strings.Contains(name, "..") {
		return "", fmt.Errorf("Nome de preferência inválido.")
	}
	return filepath.Join(s.Root, filepath.FromSlash(name)), nil
}

// Read recusa corrupção e versões desconhecidas sem modificar o original.
// name seleciona o arquivo e destination recebe o valor; ausência permanece distinguível.
func (s *Store) Read(name string, destination any) error {
	// path restringe o acesso ao namespace local.
	path, err := s.path(name)
	if err != nil {
		return err
	}
	// data não contém credenciais; limite evita leitura ilimitada de corrupção.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) > 16<<20 {
		return fmt.Errorf("Preferência local excede o limite.")
	}
	// header valida a versão antes de aplicar valores ao consumidor.
	var header struct {
		// Version delimita o formato persistido antes de decodificar o valor do consumidor.
		Version int `json:"version"`
	}
	if json.Unmarshal(data, &header) != nil || header.Version != 1 {
		return fmt.Errorf("Preferência local inválida ou de versão incompatível; arquivo preservado.")
	}
	return json.Unmarshal(data, destination)
}

// Write publica JSON atomicamente e não substitui preferências corrompidas.
// ctx permite cancelar; name é relativo e value deve conter o contrato versionado não secreto.
func (s *Store) Write(ctx context.Context, name string, value any) error {
	// path é validado antes de qualquer criação de diretório.
	path, err := s.path(name)
	if err != nil {
		return err
	}
	// previous verifica corrupção sem confundir ausência com falha.
	var previous map[string]any
	if err = s.Read(name, &previous); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	// data representa integralmente o próximo estado.
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return atomicfile.Write(ctx, path, data)
}
