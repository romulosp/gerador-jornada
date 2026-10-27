package application

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestAuthoringLimits preserva sessões, revisões e um efeito por requestId.
func TestAuthoringLimits(t *testing.T) {
	// now permite expirar sessões sem esperar tempo real.
	now := time.Unix(1800000000, 0)
	// authoring usa o contrato padrão e relógio controlado somente pelo teste.
	authoring := NewAuthoring(DefaultAuthoringLimits(), func() time.Time { return now })
	// created conserva a primeira identidade enquanto preenche o limite de sessões.
	var created AuthoringSessionInfo
	// index cria oito documentos já persistidos, sem consumir staging de edição.
	for index := 1; index <= 8; index++ {
		// info/err observam a admissão sem alterar qualquer arquivo.
		info, err := authoring.Admit(fmt.Sprintf("00000000-0000-4000-8000-%012x", index), fmt.Sprint(index), NewDocument("Teste"), AuthoringBaseline{Exists: true})
		if err != nil {
			t.Fatal(err)
		}
		if index == 1 {
			created = info
		}
	}
	if _, err := authoring.Admit("00000000-0000-4000-8000-000000000009", "ninth", NewDocument("Outra"), AuthoringBaseline{Exists: true}); authoringCode(err) != "resourceLimit" {
		t.Fatalf("nona sessão: %v", err)
	}
	// first aplica uma edição real pelo histórico compartilhado.
	first, err := authoring.Edit(created.ID, "10000000-0000-4000-8000-000000000001", "same", created.Revision, func(session *Session) (map[string]string, error) {
		return nil, session.Change(func(document *Document) error { document.Project.Name = "Alterado"; return nil })
	})
	if err != nil {
		t.Fatal(err)
	}
	// replay usa a revisão antiga: precisa devolver resultado antes da comparação CAS.
	replay, err := authoring.Edit(created.ID, "10000000-0000-4000-8000-000000000001", "same", created.Revision, nil)
	if err != nil || first.Revision != replay.Revision {
		t.Fatalf("replay: %+v (%v)", replay, err)
	}
	if _, err := authoring.Edit(created.ID, "10000000-0000-4000-8000-000000000001", "different", created.Revision, nil); authoringCode(err) != "invalidArgument" {
		t.Fatalf("mesmo ID e outro payload: %v", err)
	}
	// snapshot é uma cópia, nunca a memória mutável mantida pelo registro.
	snapshot, err := authoring.Inspect(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Document.Project.Name = "Fora da sessão"
	// preserved deve manter o conteúdo e a revisão depois de falha parcial.
	preserved, err := authoring.Edit(created.ID, "10000000-0000-4000-8000-000000000002", "failed", first.Revision, func(session *Session) (map[string]string, error) {
		session.Document.Project.Name = "Não publicar"
		return nil, errors.New("falha sintética")
	})
	_ = preserved
	if err == nil {
		t.Fatal("falha do comando foi ignorada")
	}
	snapshot, err = authoring.Inspect(created.ID)
	if err != nil || snapshot.Document.Project.Name != "Alterado" || snapshot.Revision != first.Revision {
		t.Fatalf("estado parcialmente alterado: %+v (%v)", snapshot, err)
	}
	now = now.Add(31 * time.Minute)
	if _, err := authoring.Inspect(created.ID); authoringCode(err) != "sessionExpired" {
		t.Fatalf("sessão expirada reutilizada: %v", err)
	}
}

// TestAuthoringFinalization mantém controle disponível depois da quota de edição.
func TestAuthoringFinalization(t *testing.T) {
	// limits conserva as 1.024 identidades do contrato para testar a fronteira real.
	limits := DefaultAuthoringLimits()
	// authoring usa uma raiz lógica sem adaptador de gravação.
	authoring := NewAuthoring(limits, time.Now)
	// info/err identificam o documento existente que receberá trabalho pendente.
	info, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "open", NewDocument("Teste"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	// revision acompanha todas as edições aceitas por CAS.
	revision := info.Revision
	// index produz IDs distintos sem liberar o cache da sessão.
	for index := 1; index <= limits.EditRequests; index++ {
		// result contém a revisão que deve ser usada pela próxima operação.
		result, err := authoring.Edit(info.ID, fmt.Sprintf("10000000-0000-4000-8000-%012x", index), fmt.Sprint(index), revision, func(session *Session) (map[string]string, error) {
			return nil, session.Change(func(document *Document) error { document.Project.Name += "x"; return nil })
		})
		if err != nil {
			t.Fatal(err)
		}
		revision = result.Revision
	}
	if _, err := authoring.Edit(info.ID, fmt.Sprintf("10000000-0000-4000-8000-%012x", limits.EditRequests+1), "over", revision, nil); authoringCode(err) != "resourceLimit" {
		t.Fatalf("quota de edição não aplicada: %v", err)
	}
	// before/ticket conferem preparação independente da quota de conteúdo.
	before, _ := authoring.Inspect(info.ID)
	ticket, err := authoring.Control(info.ID, "20000000-0000-4000-8000-000000000001", "review", revision, before.ControlGeneration, "review")
	if err != nil {
		t.Fatal(err)
	}
	// same deve repetir os IDs da preparação, sem gerar outra revisão.
	same, err := authoring.Control(info.ID, "20000000-0000-4000-8000-000000000001", "review", revision, before.ControlGeneration, "review")
	if err != nil || same != ticket {
		t.Fatalf("replay de controle: %+v (%v)", same, err)
	}
	// next aposenta a revisão anterior, sem apagar o rascunho.
	next, err := authoring.Control(info.ID, "20000000-0000-4000-8000-000000000002", "review-again", revision, ticket.Generation, "review")
	if err != nil || next.Generation <= ticket.Generation {
		t.Fatalf("nova revisão: %+v (%v)", next, err)
	}
	if _, err := authoring.Control(info.ID, "20000000-0000-4000-8000-000000000001", "review", revision, before.ControlGeneration, "review"); authoringCode(err) != "controlRetired" {
		t.Fatalf("controle aposentado produziu efeito: %v", err)
	}
	if _, err := authoring.Control(info.ID, "20000000-0000-4000-8000-000000000003", "discard", revision, next.Generation, "discard"); err != nil {
		t.Fatal(err)
	}
	// final não deve continuar dirty após descarte explícito.
	final, _ := authoring.Inspect(info.ID)
	if final.Dirty {
		t.Fatal("descarte não liberou o rascunho")
	}
}

// TestAuthoringBusy recusa uma segunda mutação sem esperar ou misturar snapshots.
func TestAuthoringBusy(t *testing.T) {
	// authoring/info identificam a sessão usada pelas duas goroutines.
	authoring := NewAuthoring(DefaultAuthoringLimits(), time.Now)
	info, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "open", NewDocument("Teste"), AuthoringBaseline{Exists: true})
	if err != nil {
		t.Fatal(err)
	}
	// entered/release controlam a seção de mutação sem depender de sleep.
	entered, release := make(chan struct{}), make(chan struct{})
	// pending assegura término da goroutine antes de finalizar o teste.
	var pending sync.WaitGroup
	pending.Add(1)
	go func() {
		defer pending.Done()
		_, err := authoring.Edit(info.ID, "10000000-0000-4000-8000-000000000001", "first", info.Revision, func(session *Session) (map[string]string, error) {
			close(entered)
			<-release
			return nil, session.Change(func(document *Document) error { document.Project.Name = "Primeiro"; return nil })
		})
		if err != nil {
			t.Error(err)
		}
	}()
	<-entered
	if _, err := authoring.Edit(info.ID, "10000000-0000-4000-8000-000000000002", "second", info.Revision, nil); authoringCode(err) != "busy" {
		t.Errorf("concorrência: %v", err)
	}
	close(release)
	pending.Wait()
}
