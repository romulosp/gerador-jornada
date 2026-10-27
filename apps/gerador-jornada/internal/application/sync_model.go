package application

import (
	"context"
	"time"
)

// GitHubAccount identifica a conta por ID estável, sem credenciais.
type GitHubAccount struct {
	// ID não muda quando o login é renomeado.
	ID int64 `json:"id"`
	// Login é o nome público atual.
	Login string `json:"login"`
	// Name é a autoria pública apresentada no commit.
	Name string `json:"name"`
}

// GitHubRepository representa um destino autorizado pelo provedor.
type GitHubRepository struct {
	// ID identifica o repositório independentemente de renomeações.
	ID int64 `json:"id"`
	// Owner e Name formam o nome qualificado.
	Owner string `json:"owner"`
	// Name identifica o repositório dentro do proprietário.
	Name string `json:"name"`
	// Branch é a branch padrão informada pelo servidor.
	Branch string `json:"branch"`
	// Visibility distingue público, privado e interno.
	Visibility string `json:"visibility"`
	// Writable indica permissão básica; políticas ainda podem impedir push.
	Writable bool `json:"writable"`
	// Archived e Disabled impedem selecionar o destino.
	Archived bool `json:"archived"`
	// Disabled informa que o provedor desabilitou o destino.
	Disabled bool `json:"disabled"`
}

// FullName apresenta o destino sem credenciais ou URL arbitrária.
func (r GitHubRepository) FullName() string { return r.Owner + "/" + r.Name }

// GitHubCredential é transportada somente entre adaptadores, nunca para widgets ou preferências.
type GitHubCredential struct {
	// AccessToken autoriza chamadas HTTPS.
	AccessToken string `json:"accessToken"`
	// RefreshToken e ExpiresAt são opcionais conforme o cadastro OAuth.
	RefreshToken string `json:"refreshToken,omitempty"`
	// ExpiresAt é o vencimento opcional; zero indica prazo não informado.
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
	// Scopes registra o alcance efetivamente concedido.
	Scopes string `json:"scopes"`
}

// GitHubAccountPreference guarda apenas a referência ao cofre.
type GitHubAccountPreference struct {
	// Version permite recusar estado incompatível.
	Version int `json:"version"`
	// Account é a identidade ativa.
	Account GitHubAccount `json:"account"`
	// CredentialKey identifica a credencial protegida do usuário.
	CredentialKey string `json:"credentialKey"`
}

// GitWorkspaceBinding associa uma raiz canônica ao destino sem modificar .jornada.
type GitWorkspaceBinding struct {
	// Version identifica a versão deste estado local.
	Version int `json:"version"`
	// Root é a pasta validada no sistema de arquivos.
	Root string `json:"canonicalRoot"`
	// Repository identifica o destino confirmado.
	Repository GitHubRepository `json:"repository"`
	// Pending e Published são OIDs, vazios quando não há commit.
	Pending string `json:"pendingOID,omitempty"`
	// Published é o último OID confirmado no destino, ou vazio antes do primeiro envio.
	Published string `json:"lastPublishedOID,omitempty"`
	// PublishedAt é fato histórico, não garantia de sincronismo atual.
	PublishedAt time.Time `json:"lastPublishedAt,omitempty"`
}

// GitChange descreve um arquivo que será incluído no commit.
type GitChange struct {
	// Path é um nome relativo de arquivo regular na raiz.
	Path string
	// Kind informa inclusão, alteração ou exclusão.
	Kind string
	// Digest identifica os bytes revisados; vazio indica exclusão.
	Digest string
}

// GitPublicationReview é um snapshot imutável a revalidar na confirmação.
type GitPublicationReview struct {
	// AccountID fixa a conta apresentada na revisão e impede publicar após troca de identidade.
	AccountID int64
	// Root e Repository fixam o destino revisado.
	Root string
	// Repository fixa o destino autorizado apresentado na revisão.
	Repository GitHubRepository
	// Head é o OID local; vazio representa referência ainda sem commit.
	Head string
	// Remote é o OID remoto; vazio representa referência ainda sem commit.
	Remote string
	// Changes contém o conjunto completo das jornadas elegíveis.
	Changes []GitChange
	// RemoteChanges informa quais arquivos serão atualizados localmente após confirmação.
	RemoteChanges []string
	// Pending distingue reenvio de novas edições ainda não commitadas.
	Pending bool
}

// CredentialVault delimita acesso ao armazenamento protegido do sistema.
type CredentialVault interface {
	// Load lê a chave fornecida no cofre ou informa ausência/falha de acesso.
	Load(string) (GitHubCredential, error)
	// Save protege a credencial completa sob a chave fornecida.
	Save(string, GitHubCredential) error
	// Delete remove somente a chave fornecida e informa falha de remoção.
	Delete(string) error
}

// SyncPreferences publica somente dados não secretos e versionados.
type SyncPreferences interface {
	// Read decodifica a preferência nomeada no destino, recusando estado inválido.
	Read(string, any) error
	// Write publica o valor não secreto no nome fornecido e respeita cancelamento.
	Write(context.Context, string, any) error
}
