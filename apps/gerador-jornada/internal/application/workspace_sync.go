package application

import "context"

// GitWorkspaceReader expõe consultas e recuperação local necessárias aos casos de uso.
type GitWorkspaceReader interface {
	Binding(string) (GitWorkspaceBinding, error)
	Review(context.Context, string, GitHubCredential) (GitPublicationReview, error)
	Recover(context.Context, string) error
}

// GitWorkspaceWriter delimita as duas operações confirmadas que alteram metadados Git.
type GitWorkspaceWriter interface {
	Bind(context.Context, string, GitHubRepository, GitHubCredential) error
	Publish(context.Context, GitPublicationReview, string, GitHubAccount, GitHubCredential) error
}

// WorkspaceSync coordena autorização e persistência sem importar Git, filesystem ou widgets.
type WorkspaceSync struct {
	// connection revalida a identidade e a autorização em cada ação remota.
	connection *GitHubConnection
	// reader e writer isolam o mecanismo Git escolhido pela composição.
	reader GitWorkspaceReader
	writer GitWorkspaceWriter
}

// NewWorkspaceSync aceita portas pequenas e retorna o coordenador concreto.
func NewWorkspaceSync(connection *GitHubConnection, reader GitWorkspaceReader, writer GitWorkspaceWriter) *WorkspaceSync {
	return &WorkspaceSync{connection: connection, reader: reader, writer: writer}
}

// Bind associa a raiz ao destino revalidado sem criar commit ou publicar.
func (s *WorkspaceSync) Bind(ctx context.Context, root string, repo GitHubRepository) error {
	return s.connection.WithRepository(ctx, repo, func(token GitHubCredential, _ GitHubAccount, fresh GitHubRepository) error {
		return s.writer.Bind(ctx, root, fresh, token)
	})
}

// ReviewPublication consulta alterações somente depois de conferir autorização do destino persistido.
func (s *WorkspaceSync) ReviewPublication(ctx context.Context, root string) (GitPublicationReview, error) {
	// binding não contém credenciais e fixa o repositório do workspace.
	binding, err := s.reader.Binding(root)
	if err != nil {
		return GitPublicationReview{}, err
	}
	// review é a única projeção que poderá ser entregue à UI.
	var review GitPublicationReview
	err = s.connection.WithRepository(ctx, binding.Repository, func(token GitHubCredential, account GitHubAccount, _ GitHubRepository) error {
		var readErr error
		review, readErr = s.reader.Review(ctx, root, token)
		review.AccountID = account.ID
		return readErr
	})
	return review, err
}

// Publish usa somente a revisão explicitamente confirmada e a conta atual revalidada.
func (s *WorkspaceSync) Publish(ctx context.Context, review GitPublicationReview, message string) error {
	return s.connection.WithRepository(ctx, review.Repository, func(token GitHubCredential, account GitHubAccount, _ GitHubRepository) error {
		if review.AccountID != account.ID {
			return SyncProblem("staleReview", "A conta mudou após a revisão; revise novamente antes de publicar.")
		}
		return s.writer.Publish(ctx, review, message, account, token)
	})
}
