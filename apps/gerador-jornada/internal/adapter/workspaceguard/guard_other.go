//go:build !windows

package workspaceguard

import (
	"context"
	"os"
)

// NewDefault recusa namespace de coordenação ainda não homologado em outra plataforma.
func NewDefault() (*Guard, error) { return nil, ErrUnsupported }

// Open não oferece garantias de filesystem Windows em plataforma ainda não homologada.
func (g *Guard) Open(root string) (*Lease, error) { return nil, ErrUnsupported }

// Acquire falha fechado sem provider comprovado para a plataforma futura.
func (g *Guard) Acquire(root string) (*Lease, error) { return nil, ErrUnsupported }

// Snapshot conserva a fronteira de API sem prometer leitura protegida ainda não entregue.
func (lease *Lease) Snapshot(name string) ([]byte, Baseline, error) {
	return nil, Baseline{}, ErrUnsupported
}

// CreateNew recusa persistência sem garantias equivalentes comprovadas.
func (lease *Lease) CreateNew(ctx context.Context, name string, data []byte) error {
	return ErrUnsupported
}

// CreateTransient recusa staging sem garantia nativa de encerramento ainda comprovada.
func (lease *Lease) CreateTransient(name string) (*os.File, error) { return nil, ErrUnsupported }

// BaselineContext não promete identidade/hash protegidos em plataforma ainda não homologada.
func (lease *Lease) BaselineContext(ctx context.Context, name string) (Baseline, error) {
	return Baseline{}, ErrUnsupported
}

// PhysicalPath não simula identidade de caminho NT em outra plataforma.
func (lease *Lease) PhysicalPath() (string, error) { return "", ErrUnsupported }
