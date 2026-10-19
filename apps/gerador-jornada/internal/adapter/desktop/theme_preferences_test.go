package desktop

import (
	"image/color"
	"math"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// TestThemePreferencePreservesEditing protege preferências, contraste do cabeçalho e isolamento do projeto.
func TestThemePreferencePreservesEditing(t *testing.T) {
	// app mantém as preferências compartilhadas por duas aberturas do editor neste teste.
	app := test.NewApp()
	defer app.Quit()
	app.Preferences().SetString(themePreferenceKey, "valor desconhecido")
	// editor permite conferir o padrão e alterar aparência com um formulário em andamento.
	editor := New(app, application.Example())
	defer editor.close()
	if editor.ThemeSelect.Selected != "Claro" {
		t.Fatal("preferência ausente ou inválida deve usar Claro")
	}
	clickBoard(editor, 50, 50, 0)
	editor.setZoom(.5)
	// field conserva um rascunho que a troca de tema não pode publicar ou reconstruir.
	field := editor.Inspector.Objects[1].(*widget.Form).Items[0].Widget.(*widget.Entry)
	field.SetText("Rascunho preservado")
	// before e selected permitem comparar dados e estado após cada troca.
	before := editor.Session.Document.Copy()
	selected := editor.Selected
	for _, name := range []string{"Escuro", "Azul", "Claro", "Escuro"} {
		editor.ThemeSelect.SetSelected(name)
		if app.Preferences().String(themePreferenceKey) != name {
			t.Fatal("tema escolhido não foi salvo na preferência local")
		}
		if !reflect.DeepEqual(before, editor.Session.Document) || editor.Selected != selected || editor.Zoom != .5 || len(editor.Session.UndoStack) != 0 || field.Text != "Rascunho preservado" {
			t.Fatal("troca de tema alterou documento, geometria, seleção, zoom, histórico ou rascunho")
		}
		if editor.Inspector.Objects[1].(*widget.Form).Items[0].Widget != field {
			t.Fatal("troca de tema reconstruiu o formulário")
		}
		if editor.Title.Color != app.Settings().Theme().Color(theme.ColorNameForeground, theme.VariantLight) {
			t.Fatal("título manteve cor fixa de outro tema")
		}
	}
	// reopened representa uma nova janela lendo a preferência já salva pelo usuário.
	reopened := New(app, application.Example())
	defer reopened.close()
	if reopened.ThemeSelect.Selected != "Escuro" {
		t.Fatal("nova abertura não recuperou o tema salvo")
	}
	editor.previewCurrentScreen()
	// windows identifica a prévia real criada pelo adaptador, sem executar o EXE.
	windows := app.Driver().AllWindows()
	// preview permite verificar o tema dos controles runtime depois da mudança do editor.
	preview := windows[len(windows)-1]
	defer preview.Close()
	// surface mantém um escopo visual estável para os controles da jornada.
	surface := preview.Content().(*fyne.Container).Objects[0].(*container.Scroll).Content.(*container.ThemeOverride)
	for _, name := range []string{"Azul", "Escuro"} {
		editor.ThemeSelect.SetSelected(name)
		for _, object := range surface.Content.(*fyne.Container).Objects {
			if button, ok := object.(*widget.Button); ok {
				if button.Theme().Color(theme.ColorNamePrimary, theme.VariantLight) != (jornadaTheme{}).Color(theme.ColorNamePrimary, theme.VariantLight) {
					t.Fatal("tema do editor contaminou os controles da prévia")
				}
			}
		}
	}
}

// TestThemeTextContrast exige contraste de texto normal e de ações primárias nas três paletas.
func TestThemeTextContrast(t *testing.T) {
	for _, name := range []string{"Claro", "Escuro", "Azul"} {
		// palette fornece cores reais utilizadas pelos controles, incluindo estados herdados do Fyne.
		palette := jornadaTheme{name: name}
		for _, surface := range []fyne.ThemeColorName{theme.ColorNameBackground, theme.ColorNameInputBackground, theme.ColorNameButton, theme.ColorNameMenuBackground, theme.ColorNameSelection} {
			// background e foreground representam a combinação realmente exibida pelo tema.
			background := palette.Color(surface, theme.VariantLight)
			foreground := palette.Color(theme.ColorNameForeground, theme.VariantLight)
			if contrastRatio(foreground, background) < 4.5 {
				t.Fatalf("tema %s: texto sem contraste suficiente em %s", name, surface)
			}
		}
		if contrastRatio(palette.Color(theme.ColorNameForegroundOnPrimary, theme.VariantLight), palette.Color(theme.ColorNamePrimary, theme.VariantLight)) < 4.5 {
			t.Fatalf("tema %s: ação primária sem contraste suficiente", name)
		}
	}
}

// contrastRatio calcula luminância relativa sRGB para verificar os tokens de texto, sem homologar o desenho.
func contrastRatio(first, second color.Color) float64 {
	// luminance lineariza canais sRGB antes da ponderação de luminância.
	luminance := func(value color.Color) float64 {
		// red, green e blue são canais de 16 bits retornados pela interface color.Color.
		red, green, blue, _ := value.RGBA()
		// linear converte um canal sRGB para a escala linear de intensidade.
		linear := func(channel uint32) float64 {
			// normalized ajusta o canal para o intervalo unitário.
			normalized := float64(channel) / 65535
			if normalized <= .04045 {
				return normalized / 12.92
			}
			return math.Pow((normalized+.055)/1.055, 2.4)
		}
		return .2126*linear(red) + .7152*linear(green) + .0722*linear(blue)
	}
	// left e right representam as luminâncias a ordenar para o quociente.
	left, right := luminance(first), luminance(second)
	return (math.Max(left, right) + .05) / (math.Min(left, right) + .05)
}
