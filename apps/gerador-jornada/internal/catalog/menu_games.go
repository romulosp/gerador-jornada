package catalog

import "strings"

// MenuGame relaciona a marca de um jogo à chave exata do skin.xml; imagens pertencem à biblioteca visual.
type MenuGame struct {
	// Label identifica o jogo no menu e no editor de propriedades.
	Label string
	// Variant conserva a chave XML mesmo quando o autor altera a legenda.
	Variant string
}

// MenuGames devolve as marcas normais suportadas pelo skin.xml, sem definir itens ou ordem do menu.xml.
// A seleção não deriva caminhos de arquivos: botao_duplasena, por exemplo, usa botao_dupla_sena.png.
func MenuGames() []MenuGame {
	return []MenuGame{
		{"MEGA-SENA", "botao_megasena"}, {"LOTOFÁCIL", "botao_lotofacil"},
		{"QUINA", "botao_quina"}, {"LOTOMANIA", "botao_lotomania"},
		{"LOTECA", "botao_loteca"}, {"LOTOGOL", "botao_lotogol"},
		{"DUPLA SENA", "botao_duplasena"}, {"DIA DE SORTE", "botao_dia_de_sorte"},
		{"SUPER SETE", "botao_super_sete"}, {"TIMEMANIA", "botao_timemania"},
		{"+MILIONÁRIA", "botao_milionaria"}, {"LOTERIA FEDERAL", "botao_federal"},
	}
}

// MenuGameSkin reconhece somente marcas conhecidas de projetos antigos sem variante persistida.
// Não infere navegação, destino ou movimento a partir da legenda.
func MenuGameSkin(label string) string {
	// O mapa Java reutiliza as marcas nas opções de Bolão, sem alterar seu nome persistido.
	if strings.HasPrefix(strings.ToUpper(label), "BOLÃO ") {
		return MenuGameSkin(label[len("BOLÃO "):])
	}
	if strings.HasSuffix(strings.ToUpper(label), " ESPECIAL") {
		if base := MenuGameSkin(label[:len(label)-len(" ESPECIAL")]); base != "" {
			return base + "_especial"
		}
	}
	// key tolera grafias históricas de acentos, espaços, hífen e sinal de mais.
	key := menuGameKey(label)
	// game compara a identidade visual legada, sem modificar o item original.
	for _, game := range MenuGames() {
		if key == menuGameKey(game.Label) || (game.Variant == "botao_federal" && key == "FEDERAL") {
			return game.Variant
		}
	}
	return ""
}

// MenuGameLabel identifica imagens de marca, inclusive variantes especiais existentes no XML.
// A disponibilidade da variante continua sendo validada pela biblioteca de skins.
func MenuGameLabel(variant string) (string, bool) {
	// base conserva a marca dos concursos especiais sem inventar recursos ausentes.
	base := strings.TrimSuffix(variant, "_especial")
	// game relaciona a chave persistida à marca já desenhada no PNG.
	for _, game := range MenuGames() {
		if game.Variant == base {
			return game.Label, true
		}
	}
	return "", false
}

// MenuGameCaptionMatches informa se a legenda apenas repete a marca incorporada na imagem.
func MenuGameCaptionMatches(label, variant string) bool {
	// brand identifica a marca da variante selecionada explicitamente ou por compatibilidade.
	brand, ok := MenuGameLabel(variant)
	return ok && (menuGameKey(label) == menuGameKey(brand) || MenuGameSkin(label) == variant || MenuGameSkin(label) == strings.TrimSuffix(variant, "_especial"))
}

// menuGameKey normaliza somente comparação visual de marcas conhecidas.
func menuGameKey(value string) string {
	return strings.NewReplacer("Á", "A", "Ã", "A", "Â", "A", "É", "E", "Ê", "E", "Í", "I", "Ó", "O", "Ô", "O", "Ú", "U", "+", "", "-", "", " ", "").Replace(strings.ToUpper(strings.Join(strings.Fields(value), " ")))
}
