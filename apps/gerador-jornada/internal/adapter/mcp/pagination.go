package mcp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"gerador-jornada.local/editor/internal/application"
)

// pageRequest fixa o limite de consulta sem permitir negociação do orçamento do frame.
type pageRequest struct {
	// Cursor é opaco e está vinculado à revisão/filtro da consulta anterior.
	Cursor string `json:"cursor,omitempty"`
	// Limit assume 100 quando omitido e rejeita valores fora de 1–100.
	Limit int `json:"limit,omitempty"`
}

// pageResult conserva valores integrais, contagem e continuação verificável.
type pageResult[T any] struct {
	// Items não contém valor truncado; item indivisível grande exige consulta mais específica.
	Items []T `json:"items"`
	// NextCursor vazio indica fim; Revision permite identificar a base da consulta.
	NextCursor string `json:"nextCursor,omitempty"`
	Revision   string `json:"revision"`
	// Total conta itens antes do recorte da página.
	Total int `json:"total"`
}

// pagination guarda somente chave efêmera de processo, sem registro crescente de cursores.
type pagination struct {
	// key autentica posição/escopo/revisão, sem ser autoridade para consentimento.
	key [32]byte
}

// newPagination exige entropia do SO e não persiste segredo nem consulta rede.
func newPagination() *pagination {
	// pages nasce na composição e pode servir consultas concorrentes, pois key é imutável.
	pages := &pagination{}
	if _, err := rand.Read(pages.key[:]); err != nil {
		panic("entropia de paginação indisponível")
	}
	return pages
}

// cursor codifica a posição e a base em formato limitado, sem expor caminhos/filtros.
func (pages *pagination) cursor(scope, revision string, position int) string {
	// binding evita ambiguidades de concatenação e oculta a referência lógica do escopo.
	binding, _ := json.Marshal([]string{scope, revision})
	digest := sha256.Sum256(binding)
	data := make([]byte, 40, 72)
	binary.BigEndian.PutUint64(data[:8], uint64(position))
	copy(data[8:], digest[:])
	mac := hmac.New(sha256.New, pages.key[:])
	mac.Write(data)
	data = append(data, mac.Sum(nil)...)
	return base64.RawURLEncoding.EncodeToString(data)
}

// position revalida tamanho, autenticidade, revisão e escopo antes de acessar itens.
func (pages *pagination) position(cursor, scope, revision string, count int) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	if len(cursor) != 96 {
		return 0, &application.AuthoringError{Code: "invalidArgument"}
	}
	// data é sempre limitado antes de decodificar; expected inclui o vínculo atual completo.
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(data) != 72 {
		return 0, &application.AuthoringError{Code: "invalidArgument"}
	}
	position := binary.BigEndian.Uint64(data[:8])
	if position > uint64(count) {
		return 0, &application.AuthoringError{Code: "invalidArgument"}
	}
	expected := pages.cursor(scope, revision, int(position))
	if !hmac.Equal([]byte(cursor), []byte(expected)) {
		return 0, &application.AuthoringError{Code: "invalidArgument"}
	}
	return int(position), nil
}

// paginate respeita quantidade e bytes antes de responder, preservando cada valor integral.
func paginate[T any](pages *pagination, scope, revision string, items []T, request pageRequest) (pageResult[T], error) {
	// limit usa o padrão apenas quando omitido, mantendo limite superior não negociável.
	limit := request.Limit
	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return pageResult[T]{}, &application.AuthoringError{Code: "invalidArgument"}
	}
	start, err := pages.position(request.Cursor, scope, revision, len(items))
	if err != nil {
		return pageResult[T]{}, err
	}
	// result/size reservam metadados de página e envelope além do orçamento de item.
	result := pageResult[T]{Items: []T{}, Revision: revision, Total: len(items)}
	size := 8 << 10
	for index := start; index < len(items) && len(result.Items) < limit; index++ {
		// encoded mede escapes Unicode/JSON antes de alocar o frame de saída.
		encoded, err := json.Marshal(items[index])
		if err != nil {
			return pageResult[T]{}, &application.AuthoringError{Code: "invalidArgument"}
		}
		if size+len(encoded)+1 > maxStructuredBytes {
			if len(result.Items) == 0 {
				return pageResult[T]{}, &application.AuthoringError{Code: "resourceLimit"}
			}
			break
		}
		size += len(encoded) + 1
		result.Items = append(result.Items, items[index])
	}
	if next := start + len(result.Items); next < len(items) {
		result.NextCursor = pages.cursor(scope, revision, next)
	}
	return result, nil
}
