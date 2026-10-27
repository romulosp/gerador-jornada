package application

// SyncError preserva classificação de falha sem incluir respostas livres, tokens ou códigos OAuth.
type SyncError struct {
	// Code é estável para apresentação e testes de contratos internos.
	Code string
	// Message é a orientação sanitizada em português.
	Message string
	// cause conserva diagnóstico interno sem concatená-lo à mensagem ou serializá-lo.
	cause error
}

// Error expõe somente a mensagem sanitizada e nunca concatena resposta de transporte.
func (e *SyncError) Error() string { return e.Message }

// Unwrap permite testar a causa no backend sem expor respostas brutas na interface.
func (e *SyncError) Unwrap() error { return e.cause }

// WrapSyncProblem conserva diagnóstico interno e mantém Error sanitizado.
func WrapSyncProblem(code, message string, cause error) error {
	return &SyncError{Code: code, Message: message, cause: cause}
}

// SyncProblem cria uma falha tipada sem acoplar o consumidor ao adaptador que a originou.
func SyncProblem(code, message string) error { return &SyncError{Code: code, Message: message} }
