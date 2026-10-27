//go:build !windows

// Package credential isola tokens no armazenamento protegido da plataforma.
package credential

import (
	"fmt"
	"gerador-jornada.local/editor/internal/application"
)

// Store não oferece fallback inseguro em plataformas ainda não homologadas.
type Store struct{}

// Load exige uma plataforma suportada.
func (Store) Load(string) (application.GitHubCredential, error) {
	return application.GitHubCredential{}, fmt.Errorf("Cofre GitHub disponível somente no Windows.")
}

// Save não escreve tokens em texto puro.
func (Store) Save(string, application.GitHubCredential) error {
	return fmt.Errorf("Cofre GitHub disponível somente no Windows.")
}

// Delete não modifica dados em plataformas não suportadas.
func (Store) Delete(string) error { return nil }
