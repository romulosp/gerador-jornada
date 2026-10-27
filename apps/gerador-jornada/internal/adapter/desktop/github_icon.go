package desktop

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// githubMark incorpora o Octicon mark-github-16, distribuído com a licença MIT adjacente.
// Fonte: github.com/primer/octicons, icons/mark-github-16.svg, obtido em 2026-10-24.
//
//go:embed assets/github.svg
var githubMark []byte

// githubIcon aplica o contraste do tema ao ícone vetorial sem alterar componentes da jornada.
func githubIcon() fyne.Resource {
	return theme.NewThemedResource(fyne.NewStaticResource("github.svg", githubMark))
}
