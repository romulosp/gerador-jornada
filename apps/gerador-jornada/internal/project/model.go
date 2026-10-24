// Package project define a identidade do documento e os metadados dos recursos incorporados.
package project

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

// idPattern reconhece identificadores UUID versão quatro canônicos.
var idPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// Project contém a identidade e os metadados dos recursos sem importar outros domínios.
type Project struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Name contém o nome exibido ao usuário; nomes de componentes são únicos na tela.
	Name string
	// Revision aumenta após cada edição confirmada, inclusive ao desfazer.
	Revision uint64
	// CatalogVersion seleciona as definições incorporadas dos componentes.
	CatalogVersion int
	// Assets descreve as imagens armazenadas no pacote do documento.
	Assets []Asset
}

// Asset identifica uma imagem incorporada imutável.
type Asset struct {
	// ID identifica a entidade de forma estável, independentemente de nome, rótulo ou ordenação.
	ID string
	// Path contém um caminho normalizado relativo ao pacote.
	Path string
	// MediaType aceita image/png ou image/jpeg.
	MediaType string
	// SHA256 identifica o conteúdo do recurso pelo resumo criptográfico.
	SHA256 string
	// SizeBytes expressa a quantidade de bytes descomprimidos.
	SizeBytes int64
}

// NewID cria um UUID v4 criptograficamente aleatório; falha de entropia provoca pânico.
func NewID() string {
	// bytes contém a entropia e os bits de versão do UUID.
	var bytes [16]byte
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	bytes[6] = bytes[6]&0x0f | 0x40
	bytes[8] = bytes[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// ValidID informa se o identificador é um UUID v4 canônico.
func ValidID(id string) bool { return idPattern.MatchString(id) }

// New cria metadados vazios com identidade nova e catálogo na versão corrente.
func New(name string) Project {
	return Project{ID: NewID(), Name: name, CatalogVersion: 4, Assets: []Asset{}}
}

// ValidCatalogVersion informa se a versão de catálogo é reconhecida pelo leitor.
func ValidCatalogVersion(version int) bool { return version >= 1 && version <= 4 }

// RequiresCatalogMigration informa se um documento antigo precisa de aviso e cópia antes do save.
func RequiresCatalogMigration(version int) bool { return version < 4 }

// Validate verifica identidade, nome não vazio e versão do catálogo sem resolver outros domínios.
func (p Project) Validate() error {
	if !ValidID(p.ID) || strings.TrimSpace(p.Name) == "" || !ValidCatalogVersion(p.CatalogVersion) {
		return fmt.Errorf("PROJECT_SCHEMA_INVALID: identidade, nome ou catálogo inválido")
	}
	return nil
}
