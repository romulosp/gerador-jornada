package github

import "gerador-jornada.local/editor/internal/application"

// fixtureToken mantém a credencial sintética fora de logs e dados reais.
func fixtureToken() application.GitHubCredential {
	return application.GitHubCredential{AccessToken: "fixture"}
}
