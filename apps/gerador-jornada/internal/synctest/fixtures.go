// Package synctest fornece dados descartáveis de integração sem conta GitHub ou credencial real.
package synctest

import (
	"context"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
	"path/filepath"
	"testing"
)

// Journeys grava dois pacotes íntegros com destino explícito para verificar bytes e navegação.
func Journeys(t testing.TB, root string) {
	t.Helper()
	// source e destination possuem identidades próprias e não dependem de caminhos externos.
	source := application.NewDocument("Principal")
	destination := application.NewDocument("Destino")
	// menu usa a mesma família de navegação entre jornadas do editor.
	menu, err := screen.NewComponent(catalog.ApplicationMenu, "menu", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	menu.Properties.Menu.Root.Items[0].Submenu.Items[1].TargetJourney = "Destino.jornada"
	menu.Properties.Menu.Root.Items[0].Submenu.Items[1].TargetScreenID = ""
	source.Screens[0].Components = append(source.Screens[0].Components, menu)
	if err = projectzip.Save(context.Background(), filepath.Join(root, "Principal.jornada"), source); err != nil {
		t.Fatal(err)
	}
	if err = projectzip.Save(context.Background(), filepath.Join(root, "Destino.jornada"), destination); err != nil {
		t.Fatal(err)
	}
}
