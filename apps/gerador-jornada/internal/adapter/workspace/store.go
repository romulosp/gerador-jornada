// Package workspace fornece armazenamento compartilhado entre editor manual e MCP.
package workspace

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/adapter/workspaceguard"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// Workspace representa uma pasta existente, nunca um novo arquivo de configuração de negócio.
type Workspace struct {
	// Root é absoluto e tem os links da própria pasta resolvidos.
	Root string
	// identity conserva a pasta física escolhida ao abrir, recusando troca do caminho.
	identity string
	// guard compartilha o namespace operacional entre processos e perfis.
	guard *workspaceguard.Guard
}

// WorkspaceEntry identifica uma jornada e conserva diagnóstico de entradas ilegíveis.
type WorkspaceEntry struct {
	// Filename é relativo à raiz, sem subdiretórios.
	Filename string
	// Name conserva o nome de exibição quando o pacote é legível.
	Name string
	// Error impede abrir entrada inválida sem interromper os demais arquivos.
	Error string
}

// OpenWorkspace valida uma pasta sem criar arquivos nem percorrer seus descendentes.
func OpenWorkspace(root string) (*Workspace, error) {
	// guard resolve apenas a configuração operacional, sem criar arquivos em leitura.
	guard, err := workspaceguard.NewDefault()
	if err != nil {
		return nil, err
	}
	lease, err := guard.Open(root)
	if err != nil {
		return nil, err
	}
	defer lease.Close()
	return &Workspace{Root: lease.Root, identity: lease.Identity, guard: guard}, nil
}

// lease revalida a identidade da raiz, inclusive se o caminho escolhido foi substituído.
func (w *Workspace) lease(write bool) (*workspaceguard.Lease, error) {
	if w == nil || w.guard == nil {
		return nil, os.ErrInvalid
	}
	// lease mantém o diretório físico aberto até o fim da operação.
	var lease *workspaceguard.Lease
	var err error
	if write {
		lease, err = w.guard.Acquire(w.Root)
	} else {
		lease, err = w.guard.Open(w.Root)
	}
	if err != nil {
		return nil, err
	}
	if lease.Identity != w.identity {
		lease.Close()
		return nil, fmt.Errorf("A pasta do workspace foi substituída; selecione-a novamente.")
	}
	return lease, nil
}

// Resolve conserva a API do editor para um arquivo validado na raiz.
// Consumidores de bytes devem preferir LoadSnapshot, que não reabre o caminho retornado.
func (w *Workspace) Resolve(name string) (string, error) {
	if !catalog.ValidJourneyFilename(name) {
		return "", fmt.Errorf("Nome de jornada inválido: %s", name)
	}
	// lease e Snapshot rejeitam links e nomes interpretáveis pelo sistema de arquivos.
	lease, err := w.lease(false)
	if err != nil {
		return "", err
	}
	defer lease.Close()
	if _, _, err := lease.Snapshot(name); err != nil {
		return "", err
	}
	return filepath.Join(w.Root, name), nil
}

// Load lê os mesmos bytes que identificam a base, sem reabrir arquivo entre hash e decode.
func (w *Workspace) Load(name string) (application.Document, error) {
	// document é uma cópia independente do armazenamento e de outras sessões.
	document, _, err := w.LoadSnapshot(name)
	return document, err
}

// LoadSnapshot retorna documento e precondição do arquivo lido pelo mesmo handle.
func (w *Workspace) LoadSnapshot(name string) (application.Document, workspaceguard.Baseline, error) {
	if !catalog.ValidJourneyFilename(name) {
		return application.Document{}, workspaceguard.Baseline{}, os.ErrInvalid
	}
	// lease restringe a leitura à raiz física autorizada no startup ou pelo editor.
	lease, err := w.lease(false)
	if err != nil {
		return application.Document{}, workspaceguard.Baseline{}, err
	}
	defer lease.Close()
	return loadSnapshot(lease, name)
}

// loadSnapshot permite inventariar entradas usando um único handle de raiz.
func loadSnapshot(lease *workspaceguard.Lease, name string) (application.Document, workspaceguard.Baseline, error) {
	// data e baseline representam exatamente a mesma abertura do arquivo.
	data, baseline, err := lease.Snapshot(name)
	if errors.Is(err, os.ErrNotExist) {
		return application.Document{}, baseline, fmt.Errorf("%w: %s", application.ErrJourneyMissing, name)
	}
	if err != nil {
		return application.Document{}, baseline, err
	}
	document, err := projectzip.Decode(data)
	return document, baseline, err
}

// List considera apenas a raiz e não omite silenciosamente um pacote inválido.
func (w *Workspace) List() ([]WorkspaceEntry, error) {
	// files preserva a listagem independente de qualquer documento aberto.
	lease, err := w.lease(false)
	if err != nil {
		return nil, err
	}
	defer lease.Close()
	files, err := lease.Entries()
	if err != nil {
		return nil, err
	}
	entries := []WorkspaceEntry{}
	for _, file := range files {
		if file.IsDir() || !strings.EqualFold(filepath.Ext(file.Name()), ".jornada") {
			continue
		}
		// entry usa o nome do arquivo como identidade, mesmo com nomes de projeto repetidos.
		entry := WorkspaceEntry{Filename: file.Name(), Name: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))}
		document, _, err := loadSnapshot(lease, file.Name())
		if err != nil {
			entry.Error = err.Error()
		} else {
			entry.Name = document.Project.Name
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Filename) < strings.ToLower(entries[j].Filename)
	})
	return entries, nil
}

// JourneyFilename normaliza somente o nome físico; o nome do projeto continua Unicode literal.
func JourneyFilename(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("Informe o nome da jornada.")
	}
	// clean substitui caracteres proibidos sem interpretar o nome como caminho.
	clean := strings.Map(func(char rune) rune {
		if unicode.IsControl(char) || strings.ContainsRune(`/\:<>"|?*`, char) {
			return '_'
		}
		return char
	}, name)
	clean = strings.Trim(clean, " .")
	if len([]rune(clean)) > 100 {
		clean = string([]rune(clean)[:100])
	}
	if clean == "" {
		clean = "jornada"
	}
	if !catalog.ValidJourneyFilename(clean + ".jornada") {
		clean = "jornada_" + clean
	}
	return clean + ".jornada", nil
}

// UniqueName sugere um nome livre; somente publicação exclusiva reserva o destino.
func (w *Workspace) UniqueName(name string) (string, error) {
	// lease não grava para sugerir nomes usados pelos fluxos manuais existentes.
	lease, err := w.lease(false)
	if err != nil {
		return "", err
	}
	defer lease.Close()
	return uniqueName(lease, name)
}

// uniqueName compara todas as entradas da raiz física, sem percorrer descendentes.
func uniqueName(lease *workspaceguard.Lease, name string) (string, error) {
	if !catalog.ValidJourneyFilename(name) {
		return "", os.ErrInvalid
	}
	// entries reserva a comparação sem diferenças de caixa antes da publicação exclusiva.
	entries, err := lease.Entries()
	if err != nil {
		return "", err
	}
	used := map[string]bool{}
	for _, entry := range entries {
		used[strings.ToLower(entry.Name())] = true
	}
	// base conserva o nome físico normalizado; candidate recebe sufixo em caso de colisão.
	base := strings.TrimSuffix(name, filepath.Ext(name))
	for index := 1; index < 10000; index++ {
		candidate := name
		if index > 1 {
			candidate = fmt.Sprintf("%s (%d).jornada", base, index)
		}
		if !used[strings.ToLower(candidate)] {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("Não há nome disponível para a jornada.")
}

// PublishNew publica bytes completos sob gate, sem substituir homônimo concorrente.
func (w *Workspace) PublishNew(ctx context.Context, name string, data []byte) error {
	if !catalog.ValidJourneyFilename(name) {
		return os.ErrInvalid
	}
	// lease coordena escritores e ancora o temporário/rename no mesmo diretório físico.
	lease, err := w.lease(true)
	if err != nil {
		return err
	}
	defer lease.Close()
	return lease.CreateNew(ctx, name, data)
}

// NewJourneyDocument oferece os dois modos do launcher com o mesmo fundo incorporado.
func NewJourneyDocument(name string, example bool) application.Document {
	// document usa identidades novas; o nome escolhido prevalece também no exemplo.
	document := application.NewDocument(name)
	if example {
		document = application.Example()
		document.Project.Name = name
	}
	builtinassets.DefaultBackgrounds(&document)
	return document
}

// Create grava uma nova jornada com o nome escolhido, sem sobrescrever homônimos.
func (w *Workspace) Create(name string, example bool) (string, application.Document, error) {
	// filename é independente do nome exibido no documento.
	filename, err := JourneyFilename(name)
	if err != nil {
		return "", application.Document{}, err
	}
	// lease mantém comparação e publicação sob a mesma exclusão.
	lease, err := w.lease(true)
	if err != nil {
		return "", application.Document{}, err
	}
	defer lease.Close()
	filename, err = uniqueName(lease, filename)
	if err != nil {
		return "", application.Document{}, err
	}
	document := NewJourneyDocument(name, example)
	data, err := projectzip.Encode(document)
	if err != nil {
		return "", application.Document{}, err
	}
	if _, err := projectzip.Decode(data); err != nil {
		return "", application.Document{}, err
	}
	if err := lease.CreateNew(context.Background(), filename, data); err != nil {
		return "", application.Document{}, err
	}
	return filename, document, nil
}

// Import valida antes de publicar e conserva integralmente os bytes da origem.
func (w *Workspace) Import(source string) (string, error) {
	// file é lido uma única vez e com limite, validando exatamente os bytes que serão publicados.
	file, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer file.Close()
	const maxImportSize = 250 * 1024 * 1024
	data, err := io.ReadAll(io.LimitReader(file, maxImportSize+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxImportSize {
		return "", fmt.Errorf("A jornada excede o limite de importação de 250 MiB.")
	}
	if _, err := projectzip.Decode(data); err != nil {
		return "", fmt.Errorf("Importar %s: %w", filepath.Base(source), err)
	}
	filename := filepath.Base(source)
	if !catalog.ValidJourneyFilename(filename) {
		filename, err = JourneyFilename(strings.TrimSuffix(filename, filepath.Ext(filename)))
		if err != nil {
			return "", err
		}
	}
	// lease reserva a exclusão somente depois de validar a origem inteira.
	lease, err := w.lease(true)
	if err != nil {
		return "", err
	}
	defer lease.Close()
	filename, err = uniqueName(lease, filename)
	if err != nil {
		return "", err
	}
	return filename, lease.CreateNew(context.Background(), filename, data)
}
