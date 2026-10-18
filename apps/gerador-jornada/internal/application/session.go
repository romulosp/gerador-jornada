package application

import "fmt"

// Session controla o estado mutável da edição e um histórico reversível limitado.
type Session struct {
	// Document compõe os domínios de negócio sem estado visual; na sessão, contém a revisão editada.
	Document Document
	// UndoStack contém cópias independentes das revisões anteriores.
	UndoStack []Document
	// RedoStack contém cópias independentes das revisões desfeitas.
	RedoStack []Document
	// SavedRevision identifica a última revisão salva explicitamente.
	SavedRevision uint64
	// Sequence identifica edições de forma crescente, inclusive após desfazer.
	Sequence uint64
}

// NewSession cria uma sessão de edição com dados independentes do documento recebido.
func NewSession(document Document) *Session {
	return &Session{Document: document.Copy(), Sequence: document.Project.Revision, SavedRevision: document.Project.Revision}
}

// Change aplica uma operação atômica; erros preservam o documento e o histórico.
func (s *Session) Change(operation func(*Document) error) error {
	// next permanece isolado até a operação terminar com sucesso.
	next := s.Document.Copy()
	// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
	if err := operation(&next); err != nil {
		return err
	}
	s.UndoStack = append(s.UndoStack, s.Document.Copy())
	if len(s.UndoStack) > 100 {
		s.UndoStack = s.UndoStack[1:]
	}
	s.RedoStack = nil
	s.Sequence++
	next.Project.Revision = s.Sequence
	s.Document = next
	return nil
}

// Undo restaura um estado anterior com uma nova revisão monotonicamente crescente.
func (s *Session) Undo() error {
	if len(s.UndoStack) == 0 {
		return fmt.Errorf("Nada para desfazer")
	}
	s.RedoStack = append(s.RedoStack, s.Document.Copy())
	s.Document = s.UndoStack[len(s.UndoStack)-1]
	s.UndoStack = s.UndoStack[:len(s.UndoStack)-1]
	s.Sequence++
	s.Document.Project.Revision = s.Sequence
	return nil
}

// Redo reaplica o próximo estado desfeito com uma nova revisão.
func (s *Session) Redo() error {
	if len(s.RedoStack) == 0 {
		return fmt.Errorf("Nada para refazer")
	}
	s.UndoStack = append(s.UndoStack, s.Document.Copy())
	s.Document = s.RedoStack[len(s.RedoStack)-1]
	s.RedoStack = s.RedoStack[:len(s.RedoStack)-1]
	s.Sequence++
	s.Document.Project.Revision = s.Sequence
	return nil
}

// Dirty informa se a revisão atual ainda precisa ser salva explicitamente.
func (s *Session) Dirty() bool { return s.Document.Project.Revision != s.SavedRevision }
