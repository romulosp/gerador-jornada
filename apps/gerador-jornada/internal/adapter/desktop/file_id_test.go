package desktop

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/application"
)

// TestFileIDProperties cobre o atributo comum, normalização imediata e salvamento com UUID preservado.
func TestFileIDProperties(t *testing.T) {
	// app mantém a interação dos formulários isolada da janela nativa.
	app := test.NewApp()
	defer app.Quit()
	// editor fornece a tela de identificação e todos os tipos iniciais.
	editor := New(app, application.Example())
	defer editor.close()
	// pageID conserva a identidade interna, distinta do atributo editável.
	pageID := editor.PageID
	// form contém as propriedades da tela enquanto nenhum controle está selecionado.
	form := editor.Inspector.Objects[1].(*widget.Form)
	// fileID localiza o atributo pelo nome apresentado, sem depender da ordem dos campos.
	var fileID *widget.Entry
	for _, item := range form.Items {
		if item.Text == "ID" {
			fileID = item.Widget.(*widget.Entry)
		}
	}
	if fileID == nil || fileID.Text != "PnlIdentificacao" {
		t.Fatal("tela precisa apresentar a sugestão PnlIdentificacao")
	}
	fileID.SetText("Pnl Identificação/Final!")
	if fileID.Text != "PnlIdentificacaoFinal" {
		t.Fatal("campo permitiu acentos ou caracteres especiais")
	}
	test.Tap(editor.InspectorActions.Objects[0].(*widget.Button))
	if editor.PageID != pageID || editor.page().FileID != "PnlIdentificacaoFinal" {
		t.Fatal("salvar ID alterou o UUID ou perdeu o valor")
	}
	clickBoard(editor, 50, 50, 0)
	form = editor.Inspector.Objects[1].(*widget.Form)
	fileID = nil
	for _, item := range form.Items {
		if item.Text == "ID" {
			fileID = item.Widget.(*widget.Entry)
		}
	}
	if fileID == nil {
		t.Fatal("componente não apresenta atributo ID")
	}
	fileID.SetText("Titulo Ação!")
	if fileID.Text != "TituloAcao" || fileID.Validate() != nil {
		t.Fatal("componente não deve exigir Pnl")
	}
	fileID.SetText(strings.Repeat("A", 41))
	if fileID.Validate() == nil {
		t.Fatal("limite de 40 caracteres não foi validado")
	}
	fileID.SetText("TituloAcao")
	test.Tap(editor.InspectorActions.Objects[0].(*widget.Button))
	if editor.page().ComponentByID(editor.Selected).FileID != "TituloAcao" {
		t.Fatal("ID do componente não foi salvo")
	}
}
