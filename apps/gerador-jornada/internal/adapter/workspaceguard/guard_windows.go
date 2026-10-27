//go:build windows

package workspaceguard

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"gerador-jornada.local/editor/internal/project"
	"golang.org/x/sys/windows"
)

// NewDefault usa a pasta conhecida do SO para coordenar o mesmo usuário entre configurações.
// JORNADA_CONFIG_DIR continua controlando preferências/biblioteca, não o namespace de locks.
func NewDefault() (*Guard, error) {
	// root vem da API de pastas conhecidas, sem depender de APPDATA fornecido pelo cliente.
	root, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	if err != nil {
		return nil, err
	}
	return New(filepath.Join(root, "JornadaStudio"))
}

// Open conserva a pasta física para consultas sem criar trava ou preferência.
func (g *Guard) Open(root string) (*Lease, error) { return openDirectory(root, false) }

// Acquire coordena a raiz por identidade física, inclusive caixa/aliases distintos.
func (g *Guard) Acquire(root string) (*Lease, error) {
	// lease é fechado em toda falha anterior à obtenção da trava privada.
	lease, err := openDirectory(root, true)
	if err != nil {
		return nil, err
	}
	// relative impede gravar a coordenação operacional dentro do workspace.
	relative, err := filepath.Rel(lease.Root, g.PrivateRoot)
	if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		lease.Close()
		return nil, os.ErrPermission
	}
	// directory/key não expõem caminhos pessoais no nome da trava.
	directory := filepath.Join(g.PrivateRoot, "workspace-locks")
	if err := os.MkdirAll(directory, 0700); err != nil {
		lease.Close()
		return nil, err
	}
	key := sha256.Sum256([]byte(lease.Identity))
	path := filepath.Join(directory, fmt.Sprintf("%x.lock", key))
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		lease.Close()
		return nil, err
	}
	// handle sem compartilhamento é liberado pelo SO inclusive no crash do processo.
	handle, err := windows.CreateFile(name, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		lease.Close()
		if err == windows.ERROR_SHARING_VIOLATION || err == windows.ERROR_LOCK_VIOLATION {
			return nil, ErrBusy
		}
		return nil, err
	}
	lease.lock = os.NewFile(uintptr(handle), path)
	lease.writable = true
	return lease, nil
}

// openDirectory captura um handle da autorização inicial sem ler descendentes.
func openDirectory(root string, write bool) (*Lease, error) {
	// absolute/canonical resolvem somente a raiz fornecida pelo operador.
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, err
	}
	name, err := windows.UTF16PtrFromString(canonical)
	if err != nil {
		return nil, err
	}
	// access exige escrita somente quando a aplicação solicita operação persistente.
	access := uint32(windows.FILE_GENERIC_READ)
	if write {
		access |= windows.FILE_GENERIC_WRITE
	}
	// handle fixa a raiz física mesmo que um caminho ancestral seja renomeado depois.
	handle, err := windows.CreateFile(name, access, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return nil, err
	}
	// info comprova pasta/identidade antes de entregar a capacidade de acesso relativo.
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		windows.CloseHandle(handle)
		return nil, err
	}
	if info.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY == 0 || info.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		windows.CloseHandle(handle)
		return nil, os.ErrInvalid
	}
	return &Lease{Root: canonical, Identity: fileIdentity(info), directory: os.NewFile(uintptr(handle), canonical)}, nil
}

// fileIdentity conserva volume/índice de arquivo sem depender da grafia do caminho.
func fileIdentity(info windows.ByHandleFileInformation) string {
	return fmt.Sprintf("%08x:%08x%08x", info.VolumeSerialNumber, info.FileIndexHigh, info.FileIndexLow)
}

// PhysicalPath retorna nome NT normalizado do diretório aberto, resolvendo aliases de volume.
// Serve somente a validações operacionais; a capacidade continua vinculada ao handle/Identity.
func (lease *Lease) PhysicalPath() (string, error) {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if lease.closed {
		return "", os.ErrClosed
	}
	// buffer limita o resultado ao máximo de caminho Unicode Win32; flag 2 é VOLUME_NAME_NT.
	buffer := make([]uint16, 32768)
	length, err := windows.GetFinalPathNameByHandle(windows.Handle(lease.directory.Fd()), &buffer[0], uint32(len(buffer)), 2)
	if err != nil || length == 0 || length >= uint32(len(buffer)) {
		return "", ErrUnsupported
	}
	return windows.UTF16ToString(buffer[:length]), nil
}

// openChild usa nome relativo ao handle de raiz e não segue reparse points.
// desired/disposition são internos; o cliente não escolhe flags de acesso nativas.
func (lease *Lease) openChild(name string, desired, disposition uint32) (*os.File, error) {
	return lease.openChildOptions(name, desired, disposition, 0)
}

// openChildOptions acrescenta opções internas de staging, mantendo nome/raiz/sharing fechados.
func (lease *Lease) openChildOptions(name string, desired, disposition, extraOptions uint32) (*os.File, error) {
	if lease.closed || !ValidChildName(name) {
		return nil, os.ErrInvalid
	}
	// object/attributes ancoram a abertura no diretório já autorizado.
	object, err := windows.NewNTUnicodeString(name)
	if err != nil {
		return nil, err
	}
	attributes := windows.OBJECT_ATTRIBUTES{RootDirectory: windows.Handle(lease.directory.Fd()), ObjectName: object, Attributes: windows.OBJ_CASE_INSENSITIVE | windows.OBJ_DONT_REPARSE}
	attributes.Length = uint32(unsafe.Sizeof(attributes))
	// status/handle recebem a operação NT sem resolver novamente um caminho completo.
	var status windows.IO_STATUS_BLOCK
	var handle windows.Handle
	err = windows.NtCreateFile(&handle, desired|windows.SYNCHRONIZE, &attributes, &status, nil, windows.FILE_ATTRIBUTE_NORMAL, windows.FILE_SHARE_READ, disposition, windows.FILE_NON_DIRECTORY_FILE|windows.FILE_OPEN_REPARSE_POINT|windows.FILE_SYNCHRONOUS_IO_NONALERT|extraOptions, 0, 0)
	if err != nil {
		return nil, ntError(err)
	}
	// info recusa link mesmo se um filesystem permitir abri-lo como objeto regular.
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		windows.CloseHandle(handle)
		return nil, err
	}
	if info.FileAttributes&(windows.FILE_ATTRIBUTE_DIRECTORY|windows.FILE_ATTRIBUTE_REPARSE_POINT) != 0 {
		windows.CloseHandle(handle)
		return nil, os.ErrPermission
	}
	return os.NewFile(uintptr(handle), name), nil
}

// CreateTransient cria arquivo privado exclusivo que o kernel apaga ao fechar o último handle.
// O consumidor limita tamanho/trabalho e conserva o handle; nunca publica este nome como jornada final.
func (lease *Lease) CreateTransient(name string) (*os.File, error) {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if !lease.writable {
		return nil, os.ErrPermission
	}
	return lease.openChildOptions(name, windows.FILE_GENERIC_READ|windows.FILE_GENERIC_WRITE|windows.DELETE, windows.FILE_CREATE, windows.FILE_DELETE_ON_CLOSE)
}

// BaselineContext calcula identidade/hash em blocos de 64 KiB, sem alocar o pacote inteiro.
func (lease *Lease) BaselineContext(ctx context.Context, name string) (Baseline, error) {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return Baseline{}, err
	}
	// file conserva a mesma abertura durante identidade, tamanho e leitura, sem seguir links.
	file, err := lease.openChild(name, windows.FILE_GENERIC_READ, windows.FILE_OPEN)
	if err != nil {
		return Baseline{}, err
	}
	defer file.Close()
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(windows.Handle(file.Fd()), &info); err != nil {
		return Baseline{}, err
	}
	size := int64(uint64(info.FileSizeHigh)<<32 | uint64(info.FileSizeLow))
	if size > MaxPackageBytes {
		return Baseline{}, fmt.Errorf("pacote excede 250 MiB")
	}
	// digest/buffer/total são limitados independentemente da quantidade de bytes do pacote.
	digest := sha256.New()
	buffer := make([]byte, 64<<10)
	total := int64(0)
	for {
		if err := ctx.Err(); err != nil {
			return Baseline{}, err
		}
		read, err := file.Read(buffer)
		if read > 0 {
			total += int64(read)
			if total > MaxPackageBytes {
				return Baseline{}, fmt.Errorf("pacote excede 250 MiB")
			}
			_, _ = digest.Write(buffer[:read])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return Baseline{}, err
		}
	}
	return Baseline{Exists: true, Identity: fileIdentity(info), SHA256: fmt.Sprintf("%x", digest.Sum(nil)), SizeBytes: total}, nil
}

// ntError traduz status NT para códigos Win32 reconhecidos por errors.Is/os.IsNotExist.
func ntError(err error) error {
	if status, ok := err.(windows.NTStatus); ok {
		return status.Errno()
	}
	return err
}

// Snapshot lê dados e identidade pelo mesmo handle com escrita/deleção bloqueadas.
func (lease *Lease) Snapshot(name string) ([]byte, Baseline, error) {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	// file mantém os bytes estáveis durante hash e leitura, sem seguir links.
	file, err := lease.openChild(name, windows.FILE_GENERIC_READ, windows.FILE_OPEN)
	if err != nil {
		return nil, Baseline{}, err
	}
	defer file.Close()
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(windows.Handle(file.Fd()), &info); err != nil {
		return nil, Baseline{}, err
	}
	// size limita alocação antes da leitura e o LimitReader protege crescimento inesperado.
	size := int64(uint64(info.FileSizeHigh)<<32 | uint64(info.FileSizeLow))
	if size > MaxPackageBytes {
		return nil, Baseline{}, fmt.Errorf("pacote excede 250 MiB")
	}
	data, err := io.ReadAll(io.LimitReader(file, MaxPackageBytes+1))
	if err != nil {
		return nil, Baseline{}, err
	}
	if len(data) > MaxPackageBytes {
		return nil, Baseline{}, fmt.Errorf("pacote excede 250 MiB")
	}
	digest := sha256.Sum256(data)
	return data, Baseline{Exists: true, Identity: fileIdentity(info), SHA256: fmt.Sprintf("%x", digest), SizeBytes: int64(len(data))}, nil
}

// renameInformation replica o layout documentado de FILE_RENAME_INFORMATION.
type renameInformation struct {
	// ReplaceIfExists permanece zero: criação nunca sobrescreve outro arquivo.
	ReplaceIfExists uint32
	// RootDirectory ancora o destino no mesmo handle autorizado.
	RootDirectory windows.Handle
	// FileNameLength conta bytes UTF-16 sem terminador.
	FileNameLength uint32
	// FileName inicia a cauda variável reservada no buffer nativo.
	FileName [1]uint16
}

// CreateNew publica arquivo completo por rename relativo exclusivo, sem CAS de substituição.
// Falha anterior ao rename remove somente o temporário criado pelo próprio handle.
func (lease *Lease) CreateNew(ctx context.Context, name string, data []byte) error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if !lease.writable || lease.closed {
		return os.ErrPermission
	}
	if !ValidChildName(name) || len(data) > MaxPackageBytes {
		return os.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// temporaryName não é nome de jornada e nunca é anunciado como publicação concluída.
	temporaryName := ".jornada-" + project.NewID() + ".tmp"
	file, err := lease.openChild(temporaryName, windows.FILE_GENERIC_READ|windows.FILE_GENERIC_WRITE|windows.DELETE, windows.FILE_CREATE)
	if err != nil {
		return err
	}
	// published define se o fechamento deve preservar a entrada final ou apagar temporário.
	published := false
	defer func() {
		if !published {
			// disposition remove exclusivamente o arquivo do handle desta operação.
			disposition := byte(1)
			var status windows.IO_STATUS_BLOCK
			_ = windows.NtSetInformationFile(windows.Handle(file.Fd()), &status, &disposition, 1, windows.FileDispositionInformation)
		}
		_ = file.Close()
	}()
	// offset limita cada escrita para observar cancelamento sem clonar todo o pacote.
	for offset := 0; offset < len(data); {
		if err := ctx.Err(); err != nil {
			return err
		}
		end := min(offset+(1<<20), len(data))
		written, err := file.Write(data[offset:end])
		if err != nil {
			return err
		}
		if written != end-offset {
			return io.ErrShortWrite
		}
		offset = end
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// filename/header/buffer montam a cauda variável com alinhamento nativo x64.
	filename, err := windows.UTF16FromString(name)
	if err != nil {
		return err
	}
	var header renameInformation
	filenameBytes := (len(filename) - 1) * 2
	buffer := make([]byte, max(int(unsafe.Sizeof(header)), int(unsafe.Offsetof(header.FileName))+filenameBytes))
	information := (*renameInformation)(unsafe.Pointer(&buffer[0]))
	information.RootDirectory = windows.Handle(lease.directory.Fd())
	information.FileNameLength = uint32(filenameBytes)
	copy(unsafe.Slice(&information.FileName[0], len(filename)-1), filename[:len(filename)-1])
	// status é separado do buffer porque o kernel escreve somente o resultado da operação.
	var status windows.IO_STATUS_BLOCK
	if err := windows.NtSetInformationFile(windows.Handle(file.Fd()), &status, &buffer[0], uint32(len(buffer)), windows.FileRenameInformation); err != nil {
		return ntError(err)
	}
	published = true
	return nil
}
