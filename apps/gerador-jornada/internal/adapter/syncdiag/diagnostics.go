// Package syncdiag emite diagnósticos estruturados por allowlist, sem URLs livres, corpos ou credenciais.
package syncdiag

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"gerador-jornada.local/editor/internal/application"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Logger encaminha somente campos técnicos explicitamente permitidos ao log local existente.
type Logger struct {
	// sink recebe eventos sanitizados, sem telemetria ou serviço externo.
	sink func(string, map[string]any)
	// mutex/sequence correlacionam eventos emitidos por callback e transporte concorrentes.
	mutex    sync.Mutex
	sequence uint64
}

// New recebe a fronteira do logger do desktop sem importar Fyne ou filesystem.
func New(sink func(string, map[string]any)) *Logger { return &Logger{sink: sink} }

// allowedFields define o contrato de metadados; campos novos exigem revisão explícita.
var allowedFields = map[string]bool{
	"action": true, "success": true, "duration_ms": true, "elapsed_ms": true, "oauth_configured": true,
	"callback_host": true, "callback_port": true, "callback_path": true, "method": true,
	"method_valid": true, "host_valid": true, "path_valid": true, "state_present": true, "state_valid": true,
	"code_present": true, "provider_error": true, "replay": true, "reason": true, "query_parameter_count": true,
	"http_route": true, "http_status": true, "request_id": true, "response_bytes": true,
	"response_format": true,
	"rate_limit":      true, "rate_remaining": true, "rate_reset": true, "retry_after_seconds": true,
	"sso_required": true, "sso_partial": true, "page": true, "items": true, "total_items": true,
	"repo_scope": true, "refresh_available": true, "expires_in_seconds": true, "credential_key_present": true,
	"workspace_key": true, "repository_key": true, "repository_id": true, "branch_present": true, "pending": true,
	"head_key": true, "remote_key": true, "pending_key": true, "commit_key": true, "head_present": true, "remote_present": true,
	"change_count": true, "remote_change_count": true, "commit_created": true, "confirmed": true,
	"checkpoint": true, "error_type": true, "error_code": true, "error_class": true, "cause_types": true,
	"os_error_code": true, "timeout": true, "canceled": true, "dns_addresses": true, "dns_not_found": true,
	"network": true, "connection_reused": true, "connection_was_idle": true, "connection_idle_ms": true,
	"tls_version": true, "tls_cipher": true, "tls_verified_chains": true, "tls_peer_certificates": true,
	"stack_frames": true, "panic_type": true,
}

// requestIDPattern aceita o formato técnico GitHub com grupos hexadecimais, não tokens ou Client ID.
var requestIDPattern = regexp.MustCompile(`^[A-Fa-f0-9]{1,16}(:[A-Fa-f0-9]{1,16}){2,5}$`)

// Event filtra nomes e tipos de campos antes de encaminhar o registro ao logger da aplicação.
func (l *Logger) Event(ctx context.Context, event string, fields map[string]any) {
	if l == nil || l.sink == nil {
		return
	}
	// safe não reaproveita mapas do chamador, evitando mutação concorrente durante serialização.
	safe := map[string]any{"operation_id": application.SyncOperationID(ctx)}
	for key, value := range fields {
		if !allowedFields[key] {
			continue
		}
		switch typed := value.(type) {
		case string:
			if key == "request_id" && !requestIDPattern.MatchString(typed) {
				continue
			}
			if len(typed) > 4096 {
				typed = typed[:4096]
			}
			safe[key] = typed
		case bool, int, int64, uint64, uint16, uint32:
			safe[key] = typed
		case []string:
			if key == "cause_types" || key == "stack_frames" {
				safe[key] = append([]string{}, typed...)
			}
		}
	}
	l.mutex.Lock()
	l.sequence++
	safe["sequence"] = l.sequence
	l.mutex.Unlock()
	l.sink("github."+event, safe)
}

// ErrorFields classifica causas sem serializar Error(), PathError.Path ou url.Error.URL.
func ErrorFields(err error) map[string]any {
	// fields transporta somente tipos, códigos e classes técnicas conhecidas.
	fields := map[string]any{"success": err == nil}
	if err == nil {
		return fields
	}
	fields["error_type"] = fmt.Sprintf("%T", err)
	var failure *application.SyncError
	if errors.As(err, &failure) {
		fields["error_code"] = failure.Code
	}
	// types limita a cadeia de causas para proteger contra ciclos de implementação externa.
	types := []string{}
	cause := err
	for depth := 0; cause != nil && depth < 8; depth++ {
		types = append(types, fmt.Sprintf("%T", cause))
		cause = errors.Unwrap(cause)
	}
	fields["cause_types"] = types
	if errors.Is(err, context.Canceled) {
		fields["canceled"] = true
		fields["error_class"] = "context_canceled"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		fields["timeout"] = true
		fields["error_class"] = "deadline_exceeded"
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		fields["error_class"] = "unexpected_eof"
	}
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		fields["error_class"] = "git_object_not_found"
	}
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		fields["error_class"] = "git_reference_not_found"
	}
	if errors.Is(err, transport.ErrAuthenticationRequired) {
		fields["error_class"] = "git_authentication_required"
	}
	if errors.Is(err, transport.ErrAuthorizationFailed) {
		fields["error_class"] = "git_authorization_failed"
	}
	if errors.Is(err, transport.ErrRepositoryNotFound) {
		fields["error_class"] = "git_repository_not_found"
	}
	// errno não inclui o caminho ou argumento que originou a falha do sistema.
	var errno syscall.Errno
	if errors.As(err, &errno) {
		fields["os_error_code"] = uint64(errno)
	}
	var network net.Error
	if errors.As(err, &network) {
		fields["timeout"] = network.Timeout()
	}
	var dns *net.DNSError
	if errors.As(err, &dns) {
		fields["error_class"] = "dns"
		fields["dns_not_found"] = dns.IsNotFound
	}
	var authority x509.UnknownAuthorityError
	if errors.As(err, &authority) {
		fields["error_class"] = "tls_unknown_authority"
	}
	var certificate x509.CertificateInvalidError
	if errors.As(err, &certificate) {
		fields["error_class"] = "tls_invalid_certificate"
	}
	var hostname x509.HostnameError
	if errors.As(err, &hostname) {
		fields["error_class"] = "tls_hostname_mismatch"
	}
	return fields
}

// Complete registra duração e resultado sem imprimir uma mensagem de erro potencialmente confidencial.
func (l *Logger) Complete(ctx context.Context, event string, started time.Time, err error) {
	// fields preserva classificação e duração em um único evento de término.
	fields := ErrorFields(err)
	fields["duration_ms"] = time.Since(started).Milliseconds()
	l.Event(ctx, event, fields)
}

// ProviderError aceita somente códigos OAuth publicados; descrições e valores desconhecidos são omitidos.
func ProviderError(value string) string {
	switch value {
	case "access_denied", "incorrect_client_credentials", "redirect_uri_mismatch", "bad_verification_code", "bad_refresh_token", "unsupported_grant_type", "application_suspended", "missing_code", "unverified_user_email", "expired_token":
		return value
	}
	if value == "" {
		return "none"
	}
	return "unrecognized"
}

// HTTPFields extrai apenas status, request ID validado e contadores de limite/SSO.
func HTTPFields(response *http.Response) map[string]any {
	// fields não copia headers livres nem o corpo da resposta.
	fields := map[string]any{"http_status": response.StatusCode, "response_bytes": response.ContentLength}
	// format classifica o corpo sem copiar headers livres, HTML, descrições ou tokens.
	format := "other"
	contentType := strings.ToLower(response.Header.Get("Content-Type"))
	if strings.Contains(contentType, "json") {
		format = "json"
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		format = "form"
	} else if strings.HasPrefix(contentType, "text/html") {
		format = "html"
	}
	fields["response_format"] = format
	if id := response.Header.Get("X-GitHub-Request-Id"); requestIDPattern.MatchString(id) {
		fields["request_id"] = id
	}
	for header, key := range map[string]string{"X-RateLimit-Limit": "rate_limit", "X-RateLimit-Remaining": "rate_remaining", "X-RateLimit-Reset": "rate_reset", "Retry-After": "retry_after_seconds"} {
		if number, err := strconv.ParseInt(response.Header.Get(header), 10, 64); err == nil && number >= 0 {
			fields[key] = number
		}
	}
	sso := strings.ToLower(response.Header.Get("X-GitHub-SSO"))
	fields["sso_required"] = strings.HasPrefix(sso, "required")
	fields["sso_partial"] = strings.HasPrefix(sso, "partial")
	return fields
}

// HTTPTrace observa DNS/TCP/TLS e latência sem registrar endereço livre ou conteúdo da requisição.
func (l *Logger) HTTPTrace(ctx context.Context, started time.Time) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) { l.Event(ctx, "http.dns.start", nil) },
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fields := ErrorFields(info.Err)
			fields["dns_addresses"] = len(info.Addrs)
			l.Event(ctx, "http.dns.complete", fields)
		},
		ConnectStart: func(network, address string) {
			kind := "other"
			if network == "tcp" || network == "tcp4" || network == "tcp6" {
				kind = network
			}
			l.Event(ctx, "http.connect.start", map[string]any{"network": kind})
		},
		ConnectDone:       func(network, address string, err error) { l.Complete(ctx, "http.connect.complete", started, err) },
		TLSHandshakeStart: func() { l.Event(ctx, "http.tls.start", nil) },
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			fields := ErrorFields(err)
			fields["tls_version"] = state.Version
			fields["tls_cipher"] = state.CipherSuite
			fields["tls_verified_chains"] = len(state.VerifiedChains)
			fields["tls_peer_certificates"] = len(state.PeerCertificates)
			l.Event(ctx, "http.tls.complete", fields)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			l.Event(ctx, "http.connection", map[string]any{"connection_reused": info.Reused, "connection_was_idle": info.WasIdle, "connection_idle_ms": info.IdleTime.Milliseconds()})
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) { l.Complete(ctx, "http.request.sent", started, info.Err) },
		GotFirstResponseByte: func() {
			l.Event(ctx, "http.response.first_byte", map[string]any{"elapsed_ms": time.Since(started).Milliseconds()})
		},
	}
}

// PanicFields devolve somente tipo e frames de código, sem o valor do panic ou argumentos.
func PanicFields(value any) map[string]any {
	// callers contém endereços só durante a conversão; nenhum endereço é persistido.
	callers := make([]uintptr, 32)
	count := runtime.Callers(3, callers)
	frames := runtime.CallersFrames(callers[:count])
	names := []string{}
	for {
		frame, more := frames.Next()
		names = append(names, fmt.Sprintf("%s:%d", frame.Function, frame.Line))
		if !more {
			break
		}
	}
	return map[string]any{"success": false, "panic_type": reflect.TypeOf(value).String(), "stack_frames": names}
}
