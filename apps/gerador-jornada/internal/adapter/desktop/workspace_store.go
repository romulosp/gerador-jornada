package desktop

import (
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
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/catalog"
)

// Workspace representa uma pasta existente, nunca um novo arquivo de configuração de negócio.
type Workspace struct {
	// Root é absoluto e tem os links da própria pasta resolvidos.
	Root string
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
	// absolute e canonical impedem que referências relativas dependam do diretório de execução.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(canonical)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("Workspace não é uma pasta: %s", root)
	}
	return &Workspace{Root: canonical}, nil
}

// Resolve aceita somente um arquivo da raiz, inclusive contra links que escapariam do workspace.
func (w *Workspace) Resolve(name string) (string, error) {
	if !catalog.ValidJourneyFilename(name) {
		return "", fmt.Errorf("Nome de jornada inválido: %s", name)
	}
	// target só é consultado após a validação portátil do nome.
	target := filepath.Join(w.Root, name)
	info, err := os.Lstat(target)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("Jornada não é um arquivo regular: %s", name)
	}
	// resolved também recusa junctions ou links fora da raiz em plataformas que os suportam.
	resolved, err := filepath.EvalSymlinks(target)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(filepath.Dir(resolved), w.Root) {
		return "", fmt.Errorf("Jornada fora do workspace: %s", name)
	}
	return resolved, nil
}

// Load lê o destino explícito, distinguindo ausência de corrupção sem substituir a prévia.
func (w *Workspace) Load(name string) (application.Document, error) {
	// path é restrito à raiz antes de delegar ao codec limitado do projeto.
	path, err := w.Resolve(name)
	if errors.Is(err, os.ErrNotExist) {
		return application.Document{}, fmt.Errorf("%w: %s", application.ErrJourneyMissing, name)
	}
	if err != nil {
		return application.Document{}, err
	}
	return projectzip.Load(path)
}

// List considera apenas a raiz e não omite silenciosamente um pacote inválido.
func (w *Workspace) List() ([]WorkspaceEntry, error) {
	// files preserva a listagem independente de qualquer documento aberto.
	files, err := os.ReadDir(w.Root)
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
		document, err := w.Load(file.Name())
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

// uniqueName conserva arquivos existentes, inclusive nomes equivalentes no Windows.
func (w *Workspace) uniqueName(name string) (string, error) {
	// entries reserva a comparação sem diferenças de caixa antes da publicação exclusiva.
	entries, err := os.ReadDir(w.Root)
	if err != nil {
		return "", err
	}
	used := map[string]bool{}
	for _, entry := range entries {
		used[strings.ToLower(entry.Name())] = true
	}
	base := strings.TrimSuffix(name, filepath.Ext(name))
	for index := 1; index < 10000; index++ {
		// candidate recebe sufixo somente em caso de colisão; nenhuma jornada é sobrescrita.
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

// publishNew publica bytes completos de forma exclusiva; falhas não substituem arquivo existente.
func (w *Workspace) publishNew(name string, data []byte) error {
	if !catalog.ValidJourneyFilename(name) {
		return fmt.Errorf("Nome de jornada inválido: %s", name)
	}
	// temporary fica na mesma pasta para publicar sem cruzar volumes.
	temporary, err := os.CreateTemp(w.Root, ".jornada-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return publishWorkspaceFile(temporary.Name(), filepath.Join(w.Root, name))
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
	filename, err = w.uniqueName(filename)
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
	if err := w.publishNew(filename, data); err != nil {
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
	filename, err = w.uniqueName(filename)
	if err != nil {
		return "", err
	}
	return filename, w.publishNew(filename, data)
}
