// Package workspaceguard coordena escritores e acesso relativo à raiz física autorizada.
package workspaceguard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// ErrBusy informa que outro processo já opera a mesma raiz física.
var ErrBusy = errors.New("workspace ocupado")

// ErrUnsupported impede anunciar substituição sem provider comprovado.
var ErrUnsupported = errors.New("operação de arquivo sem provider comprovado")

// MaxPackageBytes preserva o limite de leitura/gravação dos pacotes existentes.
const MaxPackageBytes = 250 << 20

// Baseline identifica exatamente o arquivo/bytes usados ao abrir o documento.
type Baseline struct {
	// Exists distingue ausência esperada de arquivo existente.
	Exists bool
	// Identity inclui volume e identidade permanente do arquivo.
	Identity string
	// SHA256 identifica o conteúdo lido pelo mesmo handle.
	SHA256 string
	// SizeBytes conta bytes comprimidos, antes da decodificação do pacote.
	SizeBytes int64
}

// Guard seleciona o namespace privado comum aos escritores da aplicação.
type Guard struct {
	// PrivateRoot é configuração operacional, nunca diretório escolhido pela ferramenta.
	PrivateRoot string
}

// Lease conserva handle da raiz; chamadas não voltam a interpretar caminhos arbitrários.
type Lease struct {
	// Root/Identity identificam caminho canônico e volume/arquivo da pasta autorizada.
	Root, Identity string
	// mu impede Close simultâneo a leitura/publicação com os mesmos handles.
	mu sync.Mutex
	// directory conserva a raiz física; lock mantém exclusão entre processos.
	directory, lock *os.File
	// closed torna o encerramento idempotente sem reutilizar handle liberado.
	closed bool
	// writable só é habilitado depois de adquirir a exclusão operacional.
	writable bool
}

// New resolve configuração sem criar pasta; arquivos de trava só nascem ao adquirir escrita.
func New(privateRoot string) (*Guard, error) {
	if privateRoot == "" {
		return nil, os.ErrInvalid
	}
	// absolute evita que troca de diretório de execução altere o namespace de coordenação.
	absolute, err := filepath.Abs(privateRoot)
	if err != nil {
		return nil, err
	}
	return &Guard{PrivateRoot: absolute}, nil
}

// ValidChildName aceita somente um nome Windows regular, sem ADS ou separador.
// O consumidor aplica ainda sua extensão/limite de domínio específico.
func ValidChildName(name string) bool {
	if name == "" || name == "." || name == ".." || !utf8.ValidString(name) || filepath.Base(name) != name || strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return false
	}
	if len(utf16.Encode([]rune(name))) > 255 {
		return false
	}
	// character bloqueia controles e sintaxe interpretada pelo filesystem.
	for _, character := range name {
		if unicode.IsControl(character) || strings.ContainsRune(`/\:<>"|?*`, character) {
			return false
		}
	}
	// base detecta nomes de dispositivo mesmo quando seguidos de extensão.
	base := strings.ToLower(strings.SplitN(name, ".", 2)[0])
	if base == "con" || base == "prn" || base == "aux" || base == "nul" || base == "conin$" || base == "conout$" || base == "clock$" {
		return false
	}
	if strings.HasPrefix(base, "com") || strings.HasPrefix(base, "lpt") {
		// suffix inclui as variantes superscritas também reservadas no Windows.
		suffix := strings.TrimPrefix(strings.TrimPrefix(base, "com"), "lpt")
		if len([]rune(suffix)) == 1 && strings.ContainsRune("123456789¹²³", []rune(suffix)[0]) {
			return false
		}
	}
	return true
}

// Close libera exclusão e raiz; arquivo de trava remanescente não conserva bloqueio.
func (lease *Lease) Close() error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if lease.closed {
		return nil
	}
	lease.closed = true
	// result conserva falha real de fechamento sem impedir liberação do outro handle.
	var result error
	if lease.lock != nil {
		result = lease.lock.Close()
	}
	if lease.directory != nil {
		if err := lease.directory.Close(); result == nil {
			result = err
		}
	}
	return result
}

// Entries lista o diretório do handle autorizado, mesmo se seu antigo caminho mudar.
func (lease *Lease) Entries() ([]os.DirEntry, error) {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if lease.closed {
		return nil, os.ErrClosed
	}
	return lease.directory.ReadDir(-1)
}
