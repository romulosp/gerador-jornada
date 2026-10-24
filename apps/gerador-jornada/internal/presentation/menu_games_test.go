package presentation

import (
	"reflect"
	"testing"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// TestMenuGameImages verifica os caminhos literais do XML, inclusive nomes não derivados da skin.
func TestMenuGameImages(t *testing.T) {
	// cases mantém referências independentes para os doze jogos normais do legado.
	cases := []struct{ label, variant, image string }{
		{"MEGA-SENA", "botao_megasena", "botao_megasena"},
		{"LOTOFÁCIL", "botao_lotofacil", "botao_lotofacil"},
		{"QUINA", "botao_quina", "botao_quina"},
		{"LOTOMANIA", "botao_lotomania", "botao_lotomania"},
		{"LOTECA", "botao_loteca", "botao_loteca"},
		{"LOTOGOL", "botao_lotogol", "botao_lotogol"},
		{"DUPLA SENA", "botao_duplasena", "botao_dupla_sena"},
		{"DIA DE SORTE", "botao_dia_de_sorte", "botao_dia_de_sorte"},
		{"SUPER SETE", "botao_super_sete", "botao_super_sete"},
		{"TIMEMANIA", "botao_timemania", "botao_timemania"},
		{"+MILIONÁRIA", "botao_milionaria", "botao_milionaria"},
		{"LOTERIA FEDERAL", "botao_federal", "botao_federal"},
		{"MEGA-SENA ESPECIAL", "botao_megasena_especial", "botao_megasena_especial"},
		{"LOTOFÁCIL ESPECIAL", "botao_lotofacil_especial", "botao_lotofacil_especial"},
		{"QUINA ESPECIAL", "botao_quina_especial", "botao_quina_especial"},
		{"DUPLA-SENA ESPECIAL", "botao_duplasena_especial", "botao_dupla_sena_especial"},
	}
	// component fornece estilo e geometria reais do composto, sem skin explícita.
	component, _ := newMenuComponent(t)
	// entry percorre cada jogo e cada estado que deve selecionar uma imagem exclusiva.
	for _, entry := range cases {
		t.Run(entry.label, func(t *testing.T) {
			// variant vazia representa projetos antigos; a explícita representa novas fábricas.
			for _, variant := range []string{"", entry.variant} {
				// state associa interação a caminhos originais, sem fallback para botão de serviço.
				for _, state := range []struct {
					name, suffix string
					enabled      bool
				}{{"default", "", true}, {"touched", "_pressed", true}, {"touched", "_disabled", false}} {
					// item e original permitem comprovar que resolver aparência não reescreve o projeto.
					item := catalog.MenuItem{ID: "jogo", Label: entry.label, Category: "JOGOS", Enabled: state.enabled, Variant: variant}
					original := item
					// group é a mesma cena usada por canvas, prévia e PPT.
					group := ApplicationMenuItem(component, item, state.name)
					if len(group.Shapes) != 1 || group.Shapes[0].SkinAsset != "/"+entry.image+state.suffix+".png" || group.Shapes[0].Text != "" {
						t.Fatalf("skin %s / %s divergente: %+v", variant, state.name, group.Shapes)
					}
					if !reflect.DeepEqual(original, item) {
						t.Fatal("resolução alterou dados confirmados")
					}
				}
			}
		})
	}
}

// TestMenuGameOverrides preserva estilos autorados e legendas personalizadas separadas do logotipo.
func TestMenuGameOverrides(t *testing.T) {
	// component e item permitem variar a precedência sem alterar identidade ou destino.
	component, _ := newMenuComponent(t)
	component.Bounds = screen.Bounds{Width: 164, Height: 68}
	item := catalog.MenuItem{ID: "milionaria", Label: "+MILIONÁRIA", Category: "JOGOS", Enabled: true, Variant: "botao_milionaria"}
	component.Enabled = false
	if shape := ApplicationMenuItem(component, item, "touched").Shapes[0]; shape.SkinAsset != "/botao_milionaria_disabled.png" {
		t.Fatalf("composto desabilitado não propagou estado: %+v", shape)
	}
	component.Enabled = true
	item.Label = "Apostar agora"
	// group mantém a marca original e uma legenda autorada editável em caixas separadas.
	group := ApplicationMenuItem(component, item, "default")
	if len(group.Shapes) != 2 || group.Shapes[0].Text != "" || group.Shapes[1].Text != item.Label || group.Shapes[1].SkinAsset != "" || group.Shapes[0].Bounds.Y+group.Shapes[0].Bounds.Height > group.Shapes[1].Bounds.Y {
		t.Fatalf("legenda sobreposta ao logotipo: %+v", group.Shapes)
	}
	item.Variant = "botao_menu"
	if shape := ApplicationMenuItem(component, item, "default").Shapes[0]; shape.SkinAsset != "/azul/botao_menu_released.png" || shape.Text != item.Label {
		t.Fatalf("variante explícita foi substituída: %+v", shape)
	}
	item.Variant, item.Label = "botao_milionaria", "+MILIONÁRIA"
	component.Style.FillColor = "#123456"
	if shape := ApplicationMenuItem(component, item, "default").Shapes[0]; shape.SkinAsset != "" || shape.Text != item.Label || shape.Fill != "#123456" {
		t.Fatalf("override de fundo perdeu texto: %+v", shape)
	}
}
