package application

import (
	"context"
	"gerador-jornada.local/editor/internal/project"
)

// syncTraceKey impede colisão com dados de contexto de outros casos de uso.
type syncTraceKey struct{}

// WithSyncOperation atribui correlação independente de state, PKCE ou credenciais.
func WithSyncOperation(ctx context.Context) context.Context {
	if SyncOperationID(ctx) != "" {
		return ctx
	}
	return context.WithValue(ctx, syncTraceKey{}, project.NewID())
}

// SyncOperationID retorna somente o identificador técnico da ação, nunca um segredo de autenticação.
func SyncOperationID(ctx context.Context) string {
	value, _ := ctx.Value(syncTraceKey{}).(string)
	return value
}
