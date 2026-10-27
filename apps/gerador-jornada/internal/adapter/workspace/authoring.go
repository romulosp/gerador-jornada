package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"

	"gerador-jornada.local/editor/internal/application"
)

// ListJourneys oferece inventário seguro para autoria, mantendo erros técnicos fora da resposta.
func (w *Workspace) ListJourneys(ctx context.Context) (application.AuthoringInventory, error) {
	if err := ctx.Err(); err != nil {
		return application.AuthoringInventory{}, err
	}
	// lease conserva a raiz física durante inventário, hashes e decodificação.
	lease, err := w.lease(false)
	if err != nil {
		return application.AuthoringInventory{}, err
	}
	defer lease.Close()
	files, err := lease.Entries()
	if err != nil {
		return application.AuthoringInventory{}, err
	}
	// inventory/hash vinculam o cursor ao conteúdo real, não apenas à lista de nomes.
	inventory := application.AuthoringInventory{Entries: []application.AuthoringJourney{}}
	hash := sha256.New()
	// file é ordenado antes do digest para não depender da enumeração do filesystem.
	sort.Slice(files, func(i, j int) bool { return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name()) })
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return application.AuthoringInventory{}, err
		}
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".jornada") {
			continue
		}
		// entry preserva o arquivo ilegível como diagnóstico individual.
		entry := application.AuthoringJourney{Filename: file.Name(), Name: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))}
		document, baseline, readErr := loadSnapshot(lease, file.Name())
		if readErr != nil {
			entry.Diagnostic = "invalidPackage"
		} else {
			entry.Name = document.Project.Name
		}
		// encoded separa campos sem ambiguidade, inclusive nomes que contêm delimitadores.
		encoded, encodeErr := json.Marshal(struct {
			Entry          application.AuthoringJourney
			Identity, Hash string
		}{entry, baseline.Identity, baseline.SHA256})
		if encodeErr != nil {
			return application.AuthoringInventory{}, encodeErr
		}
		hash.Write(encoded)
		inventory.Entries = append(inventory.Entries, entry)
	}
	inventory.Revision = hex.EncodeToString(hash.Sum(nil))
	return inventory, nil
}

// OpenDocument adapta a base física ao contrato de aplicação sem expor handles ou caminhos.
func (w *Workspace) OpenDocument(ctx context.Context, name string) (application.Document, application.AuthoringBaseline, error) {
	if err := ctx.Err(); err != nil {
		return application.Document{}, application.AuthoringBaseline{}, err
	}
	// document/baseline são lidos juntos e não dependem do caminho retornado por Resolve.
	document, baseline, err := w.LoadSnapshot(name)
	return document, application.AuthoringBaseline{Exists: baseline.Exists, Identity: baseline.Identity, SHA256: baseline.SHA256, SizeBytes: baseline.SizeBytes, RootIdentity: w.identity, Filename: name}, err
}

// NewDocument reutiliza exatamente o construtor legado utilizado pelo launcher manual.
func (w *Workspace) NewDocument(name string, example bool) application.Document {
	return NewJourneyDocument(name, example)
}

// A asserção impede divergência silenciosa entre o adaptador e a fronteira consumidora.
var _ application.AuthoringDocumentSource = (*Workspace)(nil)
