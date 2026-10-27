package application

import (
	"fmt"
	"testing"
	"time"
)

// TestAuthoringDrafts separa revisão/dirty dos anexos e do documento.
func TestAuthoringDrafts(t *testing.T) {
	// authoring/info formam uma sessão existente sem edição de documento.
	authoring := NewAuthoring(DefaultAuthoringLimits(), time.Now)
	info, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "open", NewDocument("Teste"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	// first registra resposta original para testar replay depois de atualizar o anexo.
	var first AuthoringDraft
	// index preenche os 32 anexos transitórios sem reservar commits de análise.
	for index := 1; index <= 32; index++ {
		// draft/record usam dados genéricos de coordenação; semântica é do caso de uso específico.
		draft := AuthoringDraft{ID: fmt.Sprintf("20000000-0000-4000-8000-%012x", index), Revision: 1, Origin: "legacy", Data: []byte(`{"description":"teste"}`)}
		record, err := authoring.StoreDraft(info.ID, fmt.Sprintf("10000000-0000-4000-8000-%012x", index), fmt.Sprint(index), info.Revision, draft)
		if err != nil {
			t.Fatal(err)
		}
		if index == 1 {
			first = record
		}
	}
	if _, err := authoring.StoreDraft(info.ID, "10000000-0000-4000-8000-000000000033", "overflow", info.Revision, AuthoringDraft{ID: "20000000-0000-4000-8000-000000000033", Revision: 1, Origin: "legacy", Data: []byte(`{}`)}); authoringCode(err) != "resourceLimit" {
		t.Fatalf("33º anexo: %v", err)
	}
	// snapshot mantém revisão de documento e quota de commit independentes dos anexos.
	snapshot, _ := authoring.Inspect(info.ID)
	if snapshot.Revision != info.Revision || snapshot.DocumentDirty || !snapshot.ProposalDirty || !snapshot.Dirty || snapshot.CommitPositions != 0 {
		t.Fatalf("estado misturado: %+v", snapshot)
	}
	// updated aumenta somente revisão da proposta; replay original não retorna revisão nova.
	updated := first
	updated.Revision = 2
	updated.Data = []byte(`{"description":"corrigido"}`)
	if _, err := authoring.StoreDraft(info.ID, "10000000-0000-4000-8000-000000000034", "update", info.Revision, updated); err != nil {
		t.Fatal(err)
	}
	// replay conserva identidade/hash/revisão que o cliente recebeu originalmente.
	replay, err := authoring.StoreDraft(info.ID, "10000000-0000-4000-8000-000000000001", "1", info.Revision, AuthoringDraft{ID: first.ID})
	if err != nil || replay.Revision != first.Revision || replay.SHA256 != first.SHA256 {
		t.Fatalf("replay do anexo: %+v (%v)", replay, err)
	}
	// other confere ownership sem registro paralelo ou busca por UUID global.
	other, err := authoring.Admit("00000000-0000-4000-8000-000000000002", "other", NewDocument("Outra"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := authoring.Draft(other.ID, first.ID); authoringCode(err) != "invalidReference" {
		t.Fatalf("anexo de outra sessão: %v", err)
	}
	// copied não pode alterar os bytes conservados para revisão/hash.
	copied, err := authoring.Draft(info.ID, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	copied.Data[0] = 'x'
	stored, _ := authoring.Draft(info.ID, first.ID)
	if stored.Data[0] != '{' {
		t.Fatal("bytes do anexo compartilhados")
	}
	snapshot, _ = authoring.Inspect(info.ID)
	if _, err := authoring.Control(info.ID, "30000000-0000-4000-8000-000000000001", "close", snapshot.Revision, snapshot.ControlGeneration, "close"); err == nil {
		t.Fatal("fechou análise pendente sem descarte")
	}
	authoring.CloseAll()
	if _, err := authoring.Inspect(info.ID); authoringCode(err) != "sessionExpired" {
		t.Fatalf("EOF recriou sessão: %v", err)
	}
}

// TestAuthoringCommitReservations prova a 64ª posição e recuperação sem 65ª.
func TestAuthoringCommitReservations(t *testing.T) {
	// authoring/info isolam o livro de posições sem efetuar persistência em teste.
	authoring := NewAuthoring(DefaultAuthoringLimits(), time.Now)
	info, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "open", NewDocument("Teste"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	// index consome 63 posições, liberando staging depois de cada efeito comprovado.
	for index := 1; index <= 63; index++ {
		// target/commit são identidades lógicas independentes de caminhos e de autorização.
		target := fmt.Sprintf("definition:%d", index)
		commit := fmt.Sprintf("40000000-0000-4000-8000-%012x", index)
		if err := authoring.ReserveTarget(info.ID, target); err != nil {
			t.Fatal(err)
		}
		if err := authoring.BeginIntent(info.ID, target, commit); err != nil {
			t.Fatal(err)
		}
		if err := authoring.ResolveIntent(info.ID, commit, true); err != nil {
			t.Fatal(err)
		}
	}
	if err := authoring.ReserveTarget(info.ID, "last"); err != nil {
		t.Fatal(err)
	}
	if err := authoring.ReserveTarget(info.ID, "would-be-65"); authoringCode(err) != "resourceLimit" {
		t.Fatalf("reserva indevida: %v", err)
	}
	// commit identifica a posição ainda incerta; recuperação não pede outra reserva.
	commit := "40000000-0000-4000-8000-000000000064"
	if err := authoring.BeginIntent(info.ID, "last", commit); err != nil {
		t.Fatal(err)
	}
	if err := authoring.BeginIntent(info.ID, "last", commit); err != nil {
		t.Fatal(err)
	}
	if err := authoring.ReserveTarget(info.ID, "last"); authoringCode(err) != "recoveryRequired" {
		t.Fatalf("duplicou alvo incerto: %v", err)
	}
	if _, err := authoring.Edit(info.ID, "50000000-0000-4000-8000-000000000001", "edit", info.Revision, nil); authoringCode(err) != "recoveryRequired" {
		t.Fatalf("editou com publicação incerta: %v", err)
	}
	if err := authoring.ResolveIntent(info.ID, commit, false); err != nil {
		t.Fatal(err)
	}
	if err := authoring.BeginIntent(info.ID, "last", commit); err != nil {
		t.Fatal(err)
	}
	if err := authoring.ResolveIntent(info.ID, commit, true); err != nil {
		t.Fatal(err)
	}
	if err := authoring.ResolveIntent(info.ID, commit, true); err != nil {
		t.Fatal(err)
	}
	snapshot, _ := authoring.Inspect(info.ID)
	if snapshot.CommitPositions != 64 {
		t.Fatalf("quota alterada: %d", snapshot.CommitPositions)
	}
	if err := authoring.ReserveTarget(info.ID, "extra"); authoringCode(err) != "resourceLimit" {
		t.Fatalf("65º commit admitido: %v", err)
	}
}

// TestAuthoringMemoryAdmission preserva trabalho anterior quando falta orçamento.
func TestAuthoringMemoryAdmission(t *testing.T) {
	// limits/authoring reduzem orçamento de conteúdo, mantendo reserva de controle.
	limits := DefaultAuthoringLimits()
	limits.RetainedBytes = 5 << 20
	// A fixture escala as reservas para testar esgotamento de conteúdo sem alocar buffers de produção.
	limits.FinalizationReserveBytes = 1 << 20
	limits.FinalizationWorkBytes = 128 << 10
	authoring := NewAuthoring(limits, time.Now)
	// document contém recurso real contado também no histórico, sem base64 duplicado.
	document := NewDocument("Memória")
	document.Resources["assets/teste.bin"] = make([]byte, 100<<10)
	info, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "open", document, AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	// revision conserva a última revisão aceita antes do limite.
	revision := info.Revision
	for index := 1; index < 100; index++ {
		// result/err distinguem a mutação aceita da recusa sem alteração parcial.
		result, err := authoring.Edit(info.ID, fmt.Sprintf("10000000-0000-4000-8000-%012x", index), fmt.Sprint(index), revision, func(session *Session) (map[string]string, error) {
			return nil, session.Change(func(d *Document) error { d.Project.Name += "x"; return nil })
		})
		if authoringCode(err) == "resourceLimit" {
			// snapshot/control comprovam que esgotar retenção não impede descarte/finalização.
			snapshot, _ := authoring.Inspect(info.ID)
			if snapshot.Revision != revision {
				t.Fatal("publicou mutação acima da quota")
			}
			if _, err := authoring.Control(info.ID, "30000000-0000-4000-8000-000000000001", "discard", revision, snapshot.ControlGeneration, "discard"); err != nil {
				t.Fatal(err)
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		revision = result.Revision
	}
	t.Fatal("histórico/recursos não foram contabilizados")
}
