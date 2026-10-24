package catalog

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"strings"
	"testing"

	"gerador-jornada.local/editor/internal/project"
)

// xmlMenuReference lê a fonte independente, sem reutilizar os dados gerados pela fábrica.
type xmlMenuReference struct {
	// Name e Title preservam atributos literais, inclusive espaços.
	Name  string `xml:"nome,attr"`
	Title string `xml:"tituloPrincipal,attr"`
	// Code, Option, Movement e Shortcut são metadados inertes que precisam sobreviver à importação.
	Code     string `xml:"codigoMovimento,attr"`
	Option   string `xml:"opcao,attr"`
	Movement string `xml:"nomeMovimento,attr"`
	Shortcut string `xml:"atalho,attr"`
	// Children contém somente elementos ativos; comentários são ignorados pelo decoder padrão.
	Children []xmlMenuReference `xml:"itemMenuXml"`
}

// TestMenuFactoryMatchesOriginalXML compara cada caminho e metadado com os bytes originais.
func TestMenuFactoryMatchesOriginalXML(t *testing.T) {
	// data é a fonte ISO-8859-1, incorporada ao repositório sem transcodificação.
	data, err := os.ReadFile("../../../../material-legado/menu/menu.xml")
	if err != nil {
		t.Fatal(err)
	}
	// decoder converte somente a codificação declarada, mantendo caracteres e espaços literais.
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(_ string, source io.Reader) (io.Reader, error) {
		// content e text convertem os bytes Latin-1 diretamente em pontos Unicode.
		content, err := io.ReadAll(source)
		if err != nil {
			return nil, err
		}
		text := make([]rune, len(content))
		// index e value conservam inclusive os espaços significativos nos nomes.
		for index, value := range content {
			text[index] = rune(value)
		}
		return strings.NewReader(string(text)), nil
	}
	// source contém a raiz externa menuXml, cujo único filho é Principal.
	var source xmlMenuReference
	if err := decoder.Decode(&source); err != nil {
		t.Fatal(err)
	}
	// menu é uma instância real, com UUIDs e layout derivados pelos mesmos defaults do editor.
	menu := NewApplicationMenuProperties(project.NewID)
	// compare visita cada nó na mesma ordem e confere os campos úteis à representação do menu.
	var compare func(xmlMenuReference, SubmenuProperties)
	compare = func(original xmlMenuReference, actual SubmenuProperties) {
		if original.Name != actual.Name || len(original.Children) != len(actual.Items) {
			t.Fatalf("nó %q divergente", original.Name)
		}
		// index e child identificam o caminho exato, sem localizar itens pelo nome editável.
		for index, child := range original.Children {
			// item precisa conservar todos os metadados, inclusive em ramos com descendentes.
			item := actual.Items[index]
			if item.Label != child.Name || item.MovementCode != child.Code || item.MovementOption != child.Option || item.MovementName != child.Movement || item.Shortcut != child.Shortcut || !item.Enabled {
				t.Fatalf("item %q / %d divergente: %+v", original.Name, index, item)
			}
			if len(child.Children) == 0 {
				if item.Submenu != nil {
					t.Fatalf("folha %q ganhou submenu", child.Name)
				}
				continue
			}
			if item.Submenu == nil {
				t.Fatalf("ramo %q perdido", child.Name)
			}
			if item.Submenu.Title != child.Title {
				t.Fatalf("título %q perdido", child.Title)
			}
			compare(child, *item.Submenu)
		}
	}
	compare(source.Children[0], menu.Root)
}
