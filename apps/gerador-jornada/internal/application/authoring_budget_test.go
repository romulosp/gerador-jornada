package application

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestAuthoringConcurrentBudget limita quatro chamadas mesmo em sessões distintas.
func TestAuthoringConcurrentBudget(t *testing.T) {
	// authoring/infos conservam quatro documentos de teste sem IO.
	authoring := NewAuthoring(DefaultAuthoringLimits(), time.Now)
	infos := make([]AuthoringSessionInfo, 4)
	// index/err identificam cada admissão antes de iniciar a concorrência.
	for index := range infos {
		var err error
		infos[index], err = authoring.Admit(fmt.Sprintf("00000000-0000-4000-8000-%012x", index+1), fmt.Sprint(index), NewDocument("Teste"), AuthoringBaseline{Exists: true})
		if err != nil {
			t.Fatal(err)
		}
	}
	// entered/release garantem quatro chamadas em andamento sem depender de sleep.
	entered, release := make(chan struct{}, 4), make(chan struct{})
	// pending assegura que o teste não deixe goroutine de mutação para trás.
	var pending sync.WaitGroup
	defer func() { close(release); pending.Wait() }()
	// info é capturado por valor, uma sessão independente por goroutine.
	for _, info := range infos {
		pending.Add(1)
		go func(info AuthoringSessionInfo) {
			defer pending.Done()
			_, err := authoring.Edit(info.ID, "10000000-0000-4000-8000-000000000001", "edit", info.Revision, func(session *Session) (map[string]string, error) {
				entered <- struct{}{}
				<-release
				return nil, session.Change(func(document *Document) error { document.Project.Name = "Alterado"; return nil })
			})
			if err != nil {
				t.Error(err)
			}
		}(info)
	}
	// index confirma entrada de todas as chamadas antes da quinta consulta.
	for index := 0; index < 4; index++ {
		<-entered
	}
	if _, err := authoring.Inspect(infos[0].ID); authoringCode(err) != "busy" {
		t.Fatalf("quinta chamada admitida: %v", err)
	}
}

// TestAuthoringStagingAdmission reserva finalização antes de admitir documento novo.
func TestAuthoringStagingAdmission(t *testing.T) {
	// limits permite exatamente um pacote pendente no processo.
	limits := DefaultAuthoringLimits()
	limits.StagingBytes = limits.TargetStagingBytes
	authoring := NewAuthoring(limits, time.Now)
	// first é um documento novo e já precisa de reserva para salvamento futuro.
	first, err := authoring.Admit("00000000-0000-4000-8000-000000000001", "new", NewDocument("Novo"), AuthoringBaseline{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := authoring.Admit("00000000-0000-4000-8000-000000000002", "another", NewDocument("Outro"), AuthoringBaseline{}); authoringCode(err) != "resourceLimit" {
		t.Fatalf("admitiu alvo sem staging: %v", err)
	}
	// state comprova que recusa do segundo alvo não remove reserva ou conteúdo do primeiro.
	state, err := authoring.Inspect(first.ID)
	if err != nil || !state.DocumentDirty || state.CommitPositions != 1 {
		t.Fatalf("primeiro alvo alterado: %+v (%v)", state, err)
	}
	if _, err := authoring.Control(first.ID, "30000000-0000-4000-8000-000000000001", "discard", state.Revision, state.ControlGeneration, "discard"); err != nil {
		t.Fatal(err)
	}
	if _, err := authoring.Admit("00000000-0000-4000-8000-000000000003", "after-discard", NewDocument("Terceiro"), AuthoringBaseline{}); err != nil {
		t.Fatalf("descarte não liberou reserva: %v", err)
	}
}
