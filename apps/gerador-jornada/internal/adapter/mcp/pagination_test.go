package mcp

import (
	"strings"
	"testing"
)

// TestPaginationBindsRevisionAndScope não reutiliza cursor entre bases ou operações diferentes.
func TestPaginationBindsRevisionAndScope(t *testing.T) {
	// pages conserva uma chave efêmera por processo, sem cache ilimitado de cursores.
	pages := newPagination()
	items := []string{"a", "b", "c"}
	first, err := paginate(pages, "session/screens", "12", items, pageRequest{Limit: 2})
	if err != nil || len(first.Items) != 2 || first.NextCursor == "" || first.Total != 3 {
		t.Fatalf("primeira página: %+v (%v)", first, err)
	}
	second, err := paginate(pages, "session/screens", "12", items, pageRequest{Cursor: first.NextCursor, Limit: 2})
	if err != nil || len(second.Items) != 1 || second.Items[0] != "c" || second.NextCursor != "" {
		t.Fatal("continuação inválida")
	}
	if _, err := paginate(pages, "session/screens", "13", items, pageRequest{Cursor: first.NextCursor}); err == nil {
		t.Fatal("cursor aceitou revisão diferente")
	}
	if _, err := paginate(pages, "other/screens", "12", items, pageRequest{Cursor: first.NextCursor}); err == nil {
		t.Fatal("cursor atravessou escopo")
	}
	if _, err := paginate(pages, "session/screens", "12", items, pageRequest{Cursor: first.NextCursor + "x"}); err == nil {
		t.Fatal("cursor adulterado aceito")
	}
}

// TestPaginationHonorsWireBudget não trunca valores nem encerra a conexão por página grande.
func TestPaginationHonorsWireBudget(t *testing.T) {
	// items força reduzir a quantidade por página abaixo do limite de 100 itens.
	items := []string{strings.Repeat("😀", 140000), strings.Repeat("x", 600000)}
	page, err := paginate(newPagination(), "large", "1", items, pageRequest{Limit: 100})
	if err != nil || len(page.Items) != 1 || page.Items[0] != items[0] || page.NextCursor == "" {
		t.Fatal("valor truncado ou orçamento ignorado")
	}
	if _, err := paginate(newPagination(), "large", "1", []string{strings.Repeat("x", MaxFrameBytes)}, pageRequest{}); err == nil {
		t.Fatal("item indivisível sem orçamento aceito")
	}
	if _, err := paginate(newPagination(), "a", "1", []string{}, pageRequest{Limit: 101}); err == nil {
		t.Fatal("limite superior ignorado")
	}
}
