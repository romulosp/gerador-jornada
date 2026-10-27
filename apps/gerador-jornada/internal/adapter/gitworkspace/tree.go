// Package gitworkspace versiona exclusivamente jornadas e preserva conteúdo remoto existente.
package gitworkspace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// entry mantém os bytes exatos e o modo de um arquivo regular, sem aplicar filtros Git.
type entry struct {
	// Data e Mode compõem a identidade usada na integração conservadora.
	Data []byte
	// Mode preserva a classificação Git regular/executável durante a integração.
	Mode filemode.FileMode
}

// journeyPath delimita os únicos arquivos materializados e publicados pelo Studio.
// name usa caminhos relativos Git; subpastas e demais extensões permanecem só nos objetos Git.
func journeyPath(name string) bool {
	return !strings.ContainsAny(name, `/\`) && strings.EqualFold(filepath.Ext(name), ".jornada")
}

// digest identifica os bytes revisados independentemente do algoritmo de OID Git.
// data é lido sem transformação; o retorno é SHA-256 hexadecimal de comprimento fixo.
// sum conserva o SHA-256 calculado antes de convertê-lo para o retorno hexadecimal.
func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

// validPath recusa caminhos que o Windows não pode representar sem ambiguidade.
// name é relativo à raiz Git e usa barras; erros não criam nem resolvem arquivos.
func validPath(name string) error {
	if name == "" || strings.ContainsAny(name, `\:`) || strings.HasPrefix(name, "/") {
		return fmt.Errorf("Caminho remoto inválido: %s", name)
	}
	// part é validada separadamente para impedir prefixos reservados e traversal.
	for _, part := range strings.Split(name, "/") {
		if part == "" || part == "." || part == ".." || strings.EqualFold(part, ".git") || strings.TrimRight(part, " .") != part || strings.ContainsAny(part, `<>"|?*`) {
			return fmt.Errorf("Caminho incompatível com Windows: %s", name)
		}
		// char é o próximo caractere que pode invalidar a representação do caminho no Windows.
		for _, char := range part {
			if char < 32 {
				return fmt.Errorf("Caminho remoto com caractere de controle.")
			}
		}
		// base identifica dispositivos DOS mesmo quando possuem uma extensão.
		base := strings.ToUpper(strings.Split(part, ".")[0])
		if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && base[3] >= '1' && base[3] <= '9' {
			return fmt.Errorf("Nome Windows reservado: %s", name)
		}
	}
	return nil
}

// validateTree verifica colisões sem distinção de caixa e prefixos arquivo/diretório.
// files contém a árvore futura; qualquer nome incompatível impede a materialização inteira.
func validateTree(files map[string]entry) error {
	// names inclui diretórios para detectar variantes de caixa em qualquer nível.
	names := map[string]string{}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range files {
		// err impede aceitar um caminho cuja representação no Windows seja ambígua.
		if err := validPath(name); err != nil {
			return err
		}
		// parts conserva a grafia exata do caminho autorizado.
		parts := strings.Split(name, "/")
		// count seleciona cada prefixo do caminho para detectar colisões de arquivo e pasta.
		for count := 1; count <= len(parts); count++ {
			// path fixa o arquivo ou prefixo exato desta etapa dentro da raiz ou do namespace privado.
			path := strings.Join(parts[:count], "/")
			// lower permite comparar o caminho conforme a política Windows sem distinção de caixa.
			lower := strings.ToLower(path)
			// previous guarda o valor anterior para recusar colisão ou recuperar somente bytes conhecidos.
			if previous, ok := names[lower]; ok && previous != path {
				return fmt.Errorf("Colisão de nomes Windows: %s / %s", previous, path)
			}
			names[lower] = path
			if count < len(parts) {
				// ok permite conferir o resultado desta etapa antes de prosseguir.
				if _, ok := files[path]; ok {
					return fmt.Errorf("Conflito de arquivo e pasta: %s", path)
				}
			}
		}
	}
	return nil
}

// loadTree lê objetos sem checkout de conteúdo não validado; OID vazio é árvore inicial vazia.
// repo fornece os objetos do commit oid; o retorno preserva bytes e modos regulares.
// Symlinks, submódulos, LFS e conteúdo acima dos limites resultam em erro antes do checkout.
func loadTree(repo *git.Repository, oid string) (map[string]entry, error) {
	// files não compartilha bytes mutáveis com o armazenamento Git.
	files := map[string]entry{}
	if oid == "" {
		return files, nil
	}
	// commit fixa a árvore do OID solicitado antes da leitura dos blobs.
	commit, err := repo.CommitObject(plumbing.NewHash(oid))
	if err != nil {
		return nil, err
	}
	// total limita a soma dos blobs para não materializar uma árvore remota ilimitada.
	total := int64(0)
	// walk percorre cada hash de árvore sob o prefixo relativo sem criar arquivos locais.
	var walk func(plumbing.Hash, string) error
	walk = func(hash plumbing.Hash, prefix string) error {
		// tree reúne as entradas Git deste nível, sem aplicar checkout ou filtros.
		tree, err := repo.TreeObject(hash)
		if err != nil {
			return err
		}
		// item identifica modo e conteúdo da próxima entrada Git considerada.
		for _, item := range tree.Entries {
			// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
			name := prefix + item.Name
			if err = validPath(name); err != nil {
				return err
			}
			if item.Mode == filemode.Dir {
				if err = walk(item.Hash, name+"/"); err != nil {
					return err
				}
				continue
			}
			if item.Mode != filemode.Regular && item.Mode != filemode.Executable {
				return fmt.Errorf("Arquivo Git não suportado (symlink/submódulo): %s", name)
			}
			// blob representa conteúdo Git imutável, separado dos arquivos do workspace.
			blob, err := repo.BlobObject(item.Hash)
			if err != nil {
				return err
			}
			if blob.Size > 250<<20 || total+blob.Size > 1<<30 {
				return fmt.Errorf("Conteúdo remoto excede o limite de materialização.")
			}
			// reader limita o acesso aos bytes do blob e deve ser fechado após a leitura.
			reader, err := blob.Reader()
			if err != nil {
				return err
			}
			// data conserva os bytes completos desta etapa para validar integridade antes da publicação.
			data, err := io.ReadAll(io.LimitReader(reader, (250<<20)+1))
			reader.Close()
			if err != nil {
				return err
			}
			total += int64(len(data))
			if bytes.HasPrefix(data, []byte("version https://git-lfs.github.com/spec/v1")) {
				return fmt.Errorf("Arquivo Git LFS não suportado: %s", name)
			}
			files[name] = entry{Data: data, Mode: item.Mode}
		}
		return nil
	}
	if err = walk(commit.TreeHash, ""); err != nil {
		return nil, err
	}
	return files, validateTree(files)
}

// equal compara presença, modo e conteúdo sem tentar mesclar ZIPs.
// leftOK/rightOK distinguem ausência de arquivo vazio; left/right só são usados quando presentes.
func equal(left entry, leftOK bool, right entry, rightOK bool) bool {
	return leftOK == rightOK && (!leftOK || left.Mode == right.Mode && bytes.Equal(left.Data, right.Data))
}

// merge aplica a regra B/L/R por caminho, incluindo exclusões, sem inferir renomeações.
// base é a ancestral comum e local/remote são as árvores concorrentes; retorna a união
// somente se cada caminho puder ser preservado sem escolher conteúdo divergente.
func merge(base, local, remote map[string]entry) (map[string]entry, error) {
	// paths reúne também os removidos de um dos lados.
	paths := map[string]bool{}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range base {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range local {
		paths[name] = true
	}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range remote {
		paths[name] = true
	}
	// result acumula somente entradas aceitas pelas regras conservadoras de integração.
	result := map[string]entry{}
	// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
	for name := range paths {
		// before e bOK identificam conteúdo e presença da base comum deste caminho.
		before, bOK := base[name]
		// left e lOK identificam conteúdo e presença locais usados na comparação com a base.
		left, lOK := local[name]
		// right e rOK identificam conteúdo e presença remotos usados na comparação com a base.
		right, rOK := remote[name]
		switch {
		case equal(left, lOK, before, bOK):
			if rOK {
				result[name] = right
			}
		case equal(right, rOK, before, bOK) || equal(left, lOK, right, rOK):
			if lOK {
				result[name] = left
			}
		default:
			return nil, fmt.Errorf("Conflito em %s; resolva externamente e tente novamente. Nenhum arquivo foi substituído.", name)
		}
	}
	return result, validateTree(result)
}

// fileData recusa links e junctions em todos os componentes existentes do caminho.
// root já deve ser canônica e name relativo. Retorna bytes/presença ou erro de tipo,
// tamanho, caminho ou leitura; ausência não é confundida com arquivo vazio.
func fileData(root, name string) ([]byte, bool, error) {
	// err impede aceitar um caminho cuja representação no Windows seja ambígua.
	if err := validPath(name); err != nil {
		return nil, false, err
	}
	// current nunca pode escapar da raiz canônica já validada.
	current := root
	// position e part permitem conferir diretórios intermediários antes do arquivo final.
	for position, part := range strings.Split(name, "/") {
		current = filepath.Join(current, part)
		// info descreve o objeto real no disco antes de aceitar tipo e tamanho.
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, err
		}
		// resolved permite recusar links ou junctions antes de ler o arquivo de trabalho.
		resolved, err := filepath.EvalSymlinks(current)
		if err != nil {
			return nil, false, err
		}
		if info.Mode()&os.ModeSymlink != 0 || !strings.EqualFold(resolved, current) {
			return nil, false, fmt.Errorf("Link no workspace não suportado: %s", name)
		}
		if position < len(strings.Split(name, "/"))-1 {
			if !info.IsDir() {
				return nil, false, fmt.Errorf("Conflito de pasta: %s", name)
			}
			continue
		}
		if !info.Mode().IsRegular() {
			return nil, false, fmt.Errorf("Arquivo não regular no workspace: %s", name)
		}
		if info.Size() > 250<<20 {
			return nil, false, fmt.Errorf("Arquivo excede 250 MiB: %s", name)
		}
	}
	// data representa uma leitura integral cujo digest será revalidado antes da mutação.
	data, err := os.ReadFile(current)
	return data, err == nil, err
}

// writeTree codifica blobs e árvores ordenadas sem tocar no índice ou worktree.
// repo recebe somente objetos imutáveis de files; o retorno é o OID da árvore raiz.
func writeTree(repo *git.Repository, files map[string]entry) (plumbing.Hash, error) {
	// encode escreve objetos imutáveis antes de qualquer referência mutável.
	var encode func(string) (plumbing.Hash, error)
	encode = func(prefix string) (plumbing.Hash, error) {
		// prefix delimita o nível relativo que esta chamada recursiva codifica.
		// tree reúne as entradas Git deste nível, sem aplicar checkout ou filtros.
		tree := &object.Tree{}
		// directories reúne somente filhos que precisam de uma árvore Git própria.
		directories := map[string]bool{}
		// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
		for name, item := range files {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			// relative preserva o nome do filho dentro do nível Git em codificação.
			relative := strings.TrimPrefix(name, prefix)
			// offset delimita a primeira subpasta antes de decidir entre blob e árvore.
			if offset := strings.Index(relative, "/"); offset >= 0 {
				directories[relative[:offset]] = true
				continue
			}
			// blob representa conteúdo Git imutável, separado dos arquivos do workspace.
			blob := repo.Storer.NewEncodedObject()
			blob.SetType(plumbing.BlobObject)
			// writer recebe os bytes exatos do objeto imutável antes de publicá-lo no armazenamento Git.
			writer, err := blob.Writer()
			if err != nil {
				return plumbing.ZeroHash, err
			}
			if _, err = writer.Write(item.Data); err != nil {
				writer.Close()
				return plumbing.ZeroHash, err
			}
			if err = writer.Close(); err != nil {
				return plumbing.ZeroHash, err
			}
			// hash identifica o objeto codificado sem alterar qualquer referência mutável.
			hash, err := repo.Storer.SetEncodedObject(blob)
			if err != nil {
				return plumbing.ZeroHash, err
			}
			tree.Entries = append(tree.Entries, object.TreeEntry{Name: relative, Mode: item.Mode, Hash: hash})
		}
		// name identifica o próximo caminho a validar ou integrar no conjunto de arquivos.
		for name := range directories {
			// hash identifica o objeto codificado sem alterar qualquer referência mutável.
			hash, err := encode(prefix + name + "/")
			if err != nil {
				return plumbing.ZeroHash, err
			}
			tree.Entries = append(tree.Entries, object.TreeEntry{Name: name, Mode: filemode.Dir, Hash: hash})
		}
		sort.Slice(tree.Entries, func(i, j int) bool {
			// i/j selecionam as entradas; left/right comparam diretórios com a barra exigida pelo Git.
			left, right := tree.Entries[i].Name, tree.Entries[j].Name
			if tree.Entries[i].Mode == filemode.Dir {
				left += "/"
			}
			if tree.Entries[j].Mode == filemode.Dir {
				right += "/"
			}
			return left < right
		})
		// encoded recebe a árvore completa antes da publicação imutável no armazenamento Git.
		encoded := repo.Storer.NewEncodedObject()
		// err impede publicar objeto de árvore que não pôde ser codificado.
		if err := tree.Encode(encoded); err != nil {
			return plumbing.ZeroHash, err
		}
		return repo.Storer.SetEncodedObject(encoded)
	}
	return encode("")
}
