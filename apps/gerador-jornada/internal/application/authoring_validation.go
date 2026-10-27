package application

import (
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/project"
	"gerador-jornada.local/editor/internal/screen"
	"reflect"
	"strings"
	"unicode/utf8"
)

// ValidateAuthoringText limita novas entradas por pontos de código sem normalizar nem truncar.
// Não deve ser aplicado destrutivamente à leitura histórica; diagnóstico preserva os bytes originais.
func ValidateAuthoringText(value string, limit int) error {
	if limit < 0 || !utf8.ValidString(value) || utf8.RuneCountInString(value) > limit {
		return &AuthoringError{Code: "invalidArgument"}
	}
	return nil
}

// ValidateComponentEdit valida apenas novos conteúdos textuais sem truncar dados históricos.
// Deve ser chamado no candidato de Session.Change; não altera o componente nem executa IO.
func ValidateComponentEdit(document *Document, page *screen.Screen, before, after screen.Component) error {
	if before.Kind != after.Kind || before.ID != after.ID {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if after.Name != before.Name {
		if err := ValidateAuthoringName(after.Name); err != nil {
			return err
		}
	}
	if err := validateTextDelta(reflect.ValueOf(before), reflect.ValueOf(after)); err != nil {
		return err
	}
	if err := after.Properties.Validate(after.Kind); err != nil {
		return &AuthoringError{Code: "invalidArgument"}
	}
	// options só exige UUID novo/coleção válida quando a lista é efetivamente alterada.
	var options, previous []catalog.Option
	if after.Properties.Radio != nil {
		options = after.Properties.Radio.Options
		if before.Properties.Radio != nil {
			previous = before.Properties.Radio.Options
		}
	}
	if after.Properties.Report != nil {
		options = after.Properties.Report.Rows
		if before.Properties.Report != nil {
			previous = before.Properties.Report.Rows
		}
	}
	if !reflect.DeepEqual(options, previous) {
		// seen impede repetir uma identidade dentro da nova coleção.
		seen, owned, existing := map[string]bool{}, map[string]bool{}, documentEntityIDs(document, before.ID)
		for _, option := range previous {
			owned[option.ID] = true
		}
		for _, option := range options {
			if !project.ValidID(option.ID) || seen[option.ID] || (existing[option.ID] && !owned[option.ID]) {
				return &AuthoringError{Code: "invalidArgument"}
			}
			seen[option.ID] = true
		}
	}
	if after.Properties.Keyboard != nil && after.Properties.Keyboard.TargetComponentID != "" {
		// target não pode ser resolvido em outra tela ou por nome de arquivo.
		target := page.ComponentByID(after.Properties.Keyboard.TargetComponentID)
		if target == nil || (target.Kind != catalog.TextField && target.Kind != catalog.BarcodeInput) {
			return &AuthoringError{Code: "invalidReference"}
		}
	}
	if after.Style.BackgroundAssetID != before.Style.BackgroundAssetID && after.Style.BackgroundAssetID != "" {
		// found verifica propriedade do recurso sem permitir caminho externo no comando.
		found := false
		for _, asset := range document.Project.Assets {
			found = found || (asset.ID == after.Style.BackgroundAssetID && strings.HasPrefix(asset.MediaType, "image/"))
		}
		if !found {
			return &AuthoringError{Code: "invalidReference"}
		}
	}
	return nil
}

// validateTextDelta percorre dados tipados; coleções identificadas comparam conteúdo pelo ID.
// Strings antigas inalteradas permanecem legíveis, inclusive ao reordenar uma coleção identificada.
func validateTextDelta(before, after reflect.Value) error {
	return validateTextDeltaDepth(before, after, 0)
}

// validateTextDeltaDepth limita recursão de dados de compostos antes de validar suas regras específicas.
func validateTextDeltaDepth(before, after reflect.Value, depth int) error {
	if depth > 256 {
		return &AuthoringError{Code: "invalidArgument"}
	}
	if !after.IsValid() {
		return nil
	}
	if !before.IsValid() || before.Type() != after.Type() {
		before = reflect.Zero(after.Type())
	}
	switch after.Kind() {
	case reflect.Pointer:
		if after.IsNil() {
			return nil
		}
		if before.IsNil() {
			return validateTextDeltaDepth(reflect.Value{}, after.Elem(), depth+1)
		}
		return validateTextDeltaDepth(before.Elem(), after.Elem(), depth+1)
	case reflect.String:
		if before.String() != after.String() {
			return ValidateAuthoringText(after.String(), 65536)
		}
	case reflect.Struct:
		// index identifica campo de dados do modelo, sem consultar tags de transporte.
		for index := 0; index < after.NumField(); index++ {
			if err := validateTextDeltaDepth(before.Field(index), after.Field(index), depth+1); err != nil {
				return err
			}
		}
	case reflect.Slice:
		// identified conserva a correspondência em reordenações de opções/partes compostas.
		identified := map[string]reflect.Value{}
		for index := 0; index < before.Len(); index++ {
			value := before.Index(index)
			if value.Kind() == reflect.Struct {
				if id := value.FieldByName("ID"); id.IsValid() && id.Kind() == reflect.String {
					identified[id.String()] = value
				}
			}
		}
		for index := 0; index < after.Len(); index++ {
			// current/old distinguem inclusão de item e alteração de conteúdo já existente.
			current, old := after.Index(index), reflect.Value{}
			if current.Kind() == reflect.Struct {
				if id := current.FieldByName("ID"); id.IsValid() && id.Kind() == reflect.String {
					old = identified[id.String()]
				}
			} else if index < before.Len() {
				old = before.Index(index)
			}
			if err := validateTextDeltaDepth(old, current, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// documentEntityIDs reúne identidades de entidades sem confundir referências com novos proprietários.
func documentEntityIDs(document *Document, optionOwner string) map[string]bool {
	// ids contém somente escalares, sem cópia de recursos ou conteúdo de componentes.
	ids := map[string]bool{document.Project.ID: true}
	for _, asset := range document.Project.Assets {
		ids[asset.ID] = true
	}
	for _, page := range document.Screens {
		ids[page.ID] = true
		for _, component := range page.Components {
			ids[component.ID] = true
			for _, id := range component.Properties.CompoundIDs() {
				ids[id] = true
			}
			if component.ID != optionOwner && component.Properties.Radio != nil {
				for _, option := range component.Properties.Radio.Options {
					ids[option.ID] = true
				}
			}
			if component.ID != optionOwner && component.Properties.Report != nil {
				for _, row := range component.Properties.Report.Rows {
					ids[row.ID] = true
				}
			}
		}
	}
	for _, scenario := range document.Scenarios {
		ids[scenario.ID] = true
		for _, step := range scenario.Steps {
			ids[step.ID] = true
		}
	}
	return ids
}

// ValidateAuthoringName exige conteúdo não vazio e no máximo 256 pontos de código Unicode.
// Nome de arquivo e identificadores mantêm as regras próprias, independentes do nome de exibição.
func ValidateAuthoringName(value string) error {
	if strings.TrimSpace(value) == "" {
		return &AuthoringError{Code: "invalidArgument"}
	}
	return ValidateAuthoringText(value, 256)
}
