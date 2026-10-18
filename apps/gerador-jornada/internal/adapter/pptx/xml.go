// Package pptx gera apresentações editáveis diretamente em PresentationML.
package pptx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"gerador-jornada.local/editor/internal/presentation"
	"gerador-jornada.local/editor/internal/screen"
)

// emu converte coordenadas lógicas para a unidade exata do Office.
func emu(value int) int64 { return int64(value) * 9525 }

// escape preserva texto de usuário sem permitir injeção de XML.
func escape(value string) string {
	// output recebe os caracteres XML escapados.
	output := &bytes.Buffer{}
	xml.EscapeText(output, []byte(value))
	return output.String()
}

// color remove somente o marcador da notação RGB do projeto.
func color(value string) string { return strings.TrimPrefix(value, "#") }

// solid produz preenchimento RGB explícito.
func solid(value string) string {
	if value == "" {
		return "<a:noFill/>"
	}
	return `<a:solidFill><a:srgbClr val="` + color(value) + `"/></a:solidFill>`
}

// transform posiciona uma forma sem depender da janela do editor.
func transform(bounds screen.Bounds) string {
	return fmt.Sprintf(`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`, emu(bounds.X), emu(bounds.Y), emu(bounds.Width), emu(bounds.Height))
}

// shape mantém cada controle ou texto como uma forma editável.
func shape(id int, name string, primitive presentation.Primitive, link, geometry string) string {
	if geometry == "" {
		geometry = "rect"
	}
	// hyperlink cria navegação interna sem relação externa.
	hyperlink := ""
	if link != "" {
		hyperlink = `<a:hlinkClick r:id="` + link + `" action="ppaction://hlinksldjump"/>`
	}
	// border define uma borda fina sem mudar o limite geométrico.
	border := `<a:ln w="9525">` + solid(primitive.Stroke) + `</a:ln>`
	// size converte pixels lógicos em centésimos de ponto.
	size := primitive.FontSize * 75
	if size == 0 {
		size = 900
	}
	// weight aplica destaque tipográfico quando solicitado.
	weight := ""
	if primitive.Bold {
		weight = ` b="1"`
	}
	// foreground resolve o padrão legível de texto.
	foreground := primitive.Color
	if foreground == "" {
		foreground = "#334155"
	}
	return fmt.Sprintf(`<p:sp><p:nvSpPr><p:cNvPr id="%d" name="%s">%s</p:cNvPr><p:cNvSpPr txBox="1"/><p:nvPr/></p:nvSpPr><p:spPr>%s<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>%s%s</p:spPr><p:txBody><a:bodyPr wrap="none" lIns="0" tIns="0" rIns="0" bIns="0" anchor="t"/><a:lstStyle/><a:p><a:pPr/><a:r><a:rPr lang="pt-BR" sz="%d"%s>%s<a:latin typeface="Arial"/></a:rPr><a:t xml:space="preserve">%s</a:t></a:r><a:endParaRPr lang="pt-BR"/></a:p></p:txBody></p:sp>`, id, escape(name), hyperlink, transform(primitive.Bounds), geometry, solid(primitive.Fill), border, size, weight, solid(foreground), escape(primitive.Text))
}

// groupStart mantém a geometria das primitivas no mesmo sistema lógico.
func groupStart(id int, group presentation.Group, link string) string {
	// action preserva o vínculo no grupo inteiro.
	action := ""
	if link != "" {
		action = `<a:hlinkClick r:id="` + link + `" action="ppaction://hlinksldjump"/>`
	}
	return fmt.Sprintf(`<p:grpSp><p:nvGrpSpPr><p:cNvPr id="%d" name="%s">%s</p:cNvPr><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/><a:chOff x="%d" y="%d"/><a:chExt cx="%d" cy="%d"/></a:xfrm></p:grpSpPr>`, id, escape(group.ComponentID), action, emu(group.Bounds.X), emu(group.Bounds.Y), emu(group.Bounds.Width), emu(group.Bounds.Height), emu(group.Bounds.X), emu(group.Bounds.Y), emu(group.Bounds.Width), emu(group.Bounds.Height))
}

// treeStart define o grupo raiz obrigatório de cada slide.
const treeStart = `<p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/><a:chOff x="0" y="0"/><a:chExt cx="0" cy="0"/></a:xfrm></p:grpSpPr>`

// slideNamespaces declara apenas os namespaces usados pelo pacote.
const slideNamespaces = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`

// relation cria um relacionamento OPC interno.
func relation(id, kind, target string) string {
	return `<Relationship Id="` + escape(id) + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/` + kind + `" Target="` + escape(target) + `"/>`
}

// relations envolve os vínculos da parte atual.
func relations(content string) string {
	return `<?xml version="1.0" encoding="UTF-8"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + content + `</Relationships>`
}
