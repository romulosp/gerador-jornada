package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// ScreenChanges contém somente alterações reais de tela, independentes da origem manual/MCP.
type ScreenChanges struct {
	// Name/FileID distinguem nome de exibição e identificador de arquivo gerado.
	Name, FileID *string
	// Width/Height são dimensões lógicas de 1 a 8192; omissão preserva valor atual.
	Width, Height *int
	// BackgroundColor é RGB validado, sem referência de arquivo externo.
	BackgroundColor *string
}

// ConfigureScreen valida o candidato antes de aplicá-lo, sem deslocar controles silenciosamente.
// Deve ser chamado dentro de Session.Change para conservar histórico e atomicidade.
func ConfigureScreen(document *Document, pageID string, changes ScreenChanges) error {
	// page/candidate isolam os campos escalares até validar a alteração inteira.
	page := document.Screen(pageID)
	if page == nil {
		return &AuthoringError{Code: "invalidReference"}
	}
	candidate := *page
	if changes.Name != nil && *changes.Name != candidate.Name {
		if err := ValidateAuthoringName(*changes.Name); err != nil {
			return err
		}
		candidate.Name = *changes.Name
	}
	if changes.FileID != nil {
		candidate.FileID = *changes.FileID
	}
	if changes.Width != nil {
		candidate.Size.Width = *changes.Width
	}
	if changes.Height != nil {
		candidate.Size.Height = *changes.Height
	}
	if changes.BackgroundColor != nil {
		candidate.Background.Color = *changes.BackgroundColor
	}
	if candidate.Size.Width < 1 || candidate.Size.Width > 8192 || candidate.Size.Height < 1 || candidate.Size.Height > 8192 || !catalog.ValidColor(candidate.Background.Color) {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if err := screen.ValidateFileID(candidate.FileID, true); err != nil {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if candidate.Size != page.Size {
		// component precisa continuar dentro da nova tela; nenhum clamp é aplicado pela configuração.
		for _, component := range page.Components {
			if !component.Bounds.Fits(candidate.Size) {
				return &AuthoringError{Code: "invalidArgument"}
			}
		}
	}
	*page = candidate
	return nil
}

// CreateScreen cria composição com defaults comuns e uma única entrada de histórico.
// defaultBackground só incorpora dados do adaptador; origem de fundo já existente é preferida.
func (session *Session) CreateScreen(name, backgroundSourceID string, changes ScreenChanges, defaultBackground func(*Document, string) error) (string, error) {
	if err := ValidateAuthoringName(name); err != nil {
		return "", err
	}
	// created conserva a identidade estável atribuída antes de devolver o resultado.
	created := screen.New(name)
	err := session.Change(func(document *Document) error {
		created.FileID = UniqueScreenFileID(*document, created.FileID)
		document.Screens = append(document.Screens, created)
		if source := document.Screen(backgroundSourceID); source != nil && source.Background.AssetID != "" {
			document.Screen(created.ID).Background = source.Background
		} else if defaultBackground != nil {
			if err := defaultBackground(document, created.ID); err != nil {
				return err
			}
		}
		if err := ConfigureScreen(document, created.ID, changes); err != nil {
			return err
		}
		if document.Journey.StartScreenID == "" {
			document.Journey.StartScreenID = created.ID
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return created.ID, nil
}
