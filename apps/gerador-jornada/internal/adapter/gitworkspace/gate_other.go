//go:build !windows

package gitworkspace

import "fmt"

// lock recusa mutação em plataforma ainda não homologada, sem prometer trava entre processos.
// O parâmetro não é usado para I/O e o retorno nunca permite iniciar uma mutação.
func lock(string) (func(), error) {
	return nil, fmt.Errorf("Integração Git disponível somente no Windows.")
}
