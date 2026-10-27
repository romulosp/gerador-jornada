package syncdiag

import (
	"context"
	"encoding/json"
	"errors"
	"gerador-jornada.local/editor/internal/application"
	"strings"
	"testing"
)

// TestAllowlistAndErrorsDoNotSerializeMessages recusa campos livres e conteúdo textual de causas técnicas.
func TestAllowlistAndErrorsDoNotSerializeMessages(t *testing.T) {
	// output representa somente o registro que seria entregue ao logger JSONL.
	output := ""
	logger := New(func(event string, fields map[string]any) {
		data, _ := json.Marshal(fields)
		output += event + string(data)
	})
	ctx := application.WithSyncOperation(context.Background())
	logger.Event(ctx, "fixture", map[string]any{"client_id": "confidential-id", "client_secret": "confidential-secret", "Authorization": "confidential-token", "message": "confidential-body", "body": "confidential-body", "request_id": "confidential-token", "http_status": 401})
	logger.Event(ctx, "failure", ErrorFields(errors.New("confidential-secret and confidential-token")))
	if strings.Contains(output, "confidential") {
		t.Fatal("segredo ou campo livre aceito")
	}
	if !strings.Contains(output, "operation_id") || !strings.Contains(output, "401") || !strings.Contains(output, "error_type") {
		t.Fatal("metadados técnicos seguros ausentes")
	}
}
