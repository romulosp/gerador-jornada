package pptx

import (
	"fmt"
	"strings"
)

// baseParts cria apresentação, tema, master, layout e seus vínculos.
func (d *deck) baseParts() map[string]string {
	// types declara os tipos MIME de todas as partes.
	types := `<?xml version="1.0" encoding="UTF-8"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Default Extension="png" ContentType="image/png"/><Default Extension="jpg" ContentType="image/jpeg"/>`
	types += `<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/><Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/><Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>`
	// slideIDs preserva a ordem de navegação do plano.
	slideIDs := ""
	// links registra master e slides no namespace da apresentação.
	links := relation("rIdMaster", "slideMaster", "slideMasters/slideMaster1.xml")
	// index determina os nomes estáveis das partes.
	for index := range d.Slides {
		types += fmt.Sprintf(`<Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, index+1)
		slideIDs += fmt.Sprintf(`<p:sldId id="%d" r:id="rIdSlide%d"/>`, 256+index, index+1)
		links += relation(fmt.Sprintf("rIdSlide%d", index+1), "slide", fmt.Sprintf("slides/slide%d.xml", index+1))
	}
	types += "</Types>"
	// master é vazio, pois cada controle possui formatação explícita.
	master := `<?xml version="1.0" encoding="UTF-8"?><p:sldMaster ` + slideNamespaces + `><p:cSld>` + treeStart + `</p:spTree></p:cSld><p:clrMap accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" bg1="lt1" bg2="lt2" folHlink="folHlink" hlink="hlink" tx1="dk1" tx2="dk2"/><p:sldLayoutIdLst><p:sldLayoutId id="2147483649" r:id="rIdLayout"/></p:sldLayoutIdLst><p:txStyles><p:titleStyle/><p:bodyStyle/><p:otherStyle/></p:txStyles></p:sldMaster>`
	// layout não impõe placeholders sobre o desenho.
	layout := `<?xml version="1.0" encoding="UTF-8"?><p:sldLayout ` + slideNamespaces + ` type="blank" preserve="1"><p:cSld name="Tela livre">` + treeStart + `</p:spTree></p:cSld><p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr></p:sldLayout>`
	return map[string]string{
		"[Content_Types].xml":                          types,
		"_rels/.rels":                                  relations(relation("rIdPresentation", "officeDocument", "ppt/presentation.xml")),
		"ppt/presentation.xml":                         fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><p:presentation %s><p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rIdMaster"/></p:sldMasterIdLst><p:sldIdLst>%s</p:sldIdLst><p:sldSz cx="%d" cy="%d" type="custom"/><p:notesSz cx="6858000" cy="9144000"/><p:defaultTextStyle/></p:presentation>`, slideNamespaces, slideIDs, emu(d.Width), emu(d.Height)),
		"ppt/_rels/presentation.xml.rels":              relations(links),
		"ppt/slideMasters/slideMaster1.xml":            master,
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": relations(relation("rIdLayout", "slideLayout", "../slideLayouts/slideLayout1.xml") + relation("rIdTheme", "theme", "../theme/theme1.xml")),
		"ppt/slideLayouts/slideLayout1.xml":            layout,
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels": relations(relation("rIdMaster", "slideMaster", "../slideMasters/slideMaster1.xml")),
		"ppt/theme/theme1.xml":                         themeXML(),
	}
}

// themeXML fornece os esquemas obrigatórios sem depender de um arquivo Office externo.
func themeXML() string {
	// scheme define cores e links consistentes com a identidade do editor.
	scheme := ""
	// item declara cada slot obrigatório do tema.
	for _, item := range []struct{ name, value string }{{"dk1", "0F172A"}, {"lt1", "FFFFFF"}, {"dk2", "334155"}, {"lt2", "F8FAFC"}, {"accent1", "2563EB"}, {"accent2", "0D9488"}, {"accent3", "F59E0B"}, {"accent4", "7C3AED"}, {"accent5", "E11D48"}, {"accent6", "64748B"}, {"hlink", "2563EB"}, {"folHlink", "7C3AED"}} {
		scheme += "<a:" + item.name + "><a:srgbClr val=\"" + item.value + "\"/></a:" + item.name + ">"
	}
	// fill é um preenchimento de tema repetido nos três níveis exigidos.
	fill := `<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>`
	// line descreve um contorno sem efeitos adicionais.
	line := `<a:ln w="9525">` + fill + `<a:prstDash val="solid"/></a:ln>`
	return `<?xml version="1.0" encoding="UTF-8"?><a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Jornada"><a:themeElements><a:clrScheme name="Jornada">` + scheme + `</a:clrScheme><a:fontScheme name="Jornada"><a:majorFont><a:latin typeface="Arial"/><a:ea typeface=""/><a:cs typeface=""/></a:majorFont><a:minorFont><a:latin typeface="Arial"/><a:ea typeface=""/><a:cs typeface=""/></a:minorFont></a:fontScheme><a:fmtScheme name="Jornada"><a:fillStyleLst>` + strings.Repeat(fill, 3) + `</a:fillStyleLst><a:lnStyleLst>` + strings.Repeat(line, 3) + `</a:lnStyleLst><a:effectStyleLst>` + strings.Repeat(`<a:effectStyle><a:effectLst/></a:effectStyle>`, 3) + `</a:effectStyleLst><a:bgFillStyleLst>` + strings.Repeat(fill, 3) + `</a:bgFillStyleLst></a:fmtScheme></a:themeElements><a:objectDefaults/><a:extraClrSchemeLst/></a:theme>`
}
