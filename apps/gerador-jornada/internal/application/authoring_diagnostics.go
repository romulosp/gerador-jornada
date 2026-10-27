package application

import "context"

// Diagnostics valida o documento protegido sem duplicar recursos binários para uma consulta R0.
// Retorna revisão e valores de diagnóstico independentes; não altera conteúdo, histórico ou oferta.
func (authoring *Authoring) Diagnostics(ctx context.Context, id string) (uint64, []Diagnostic, error) {
	// state/release preservam a mesma revisão durante a validação entre domínios.
	state, release, err := authoring.acquire(id)
	if err != nil {
		return 0, nil, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return 0, nil, err
	}
	// diagnostics não retém referências a recursos, widgets ou textos completos dos componentes.
	diagnostics := state.model.Document.Validate()
	if err := ctx.Err(); err != nil {
		return 0, nil, err
	}
	return state.model.Document.Project.Revision, diagnostics, nil
}
