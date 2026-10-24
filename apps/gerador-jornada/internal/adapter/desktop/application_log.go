package desktop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"gerador-jornada.local/editor/internal/adapter/buildinfo"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// logSegmentLimit é o tamanho máximo de um segmento antes da rotação.
const logSegmentLimit = 2 << 20

// logRetention limita os arquivos próprios removíveis mantidos na pasta.
const logRetention = 20

// applicationLog grava eventos JSONL seguros sem conteúdo do projeto.
type applicationLog struct {
	// mu serializa escrita, rotação e consulta de caminho.
	mu sync.Mutex
	// dir contém a pasta efetiva do log.
	dir string
	// path contém o arquivo corrente.
	path string
	// file mantém o descritor aberto para escrita síncrona.
	file *os.File
	// size contém o tamanho corrente do segmento.
	size int64
	// limit permite reduzir o segmento em testes.
	limit int64
	// session identifica a execução sem expor dados do usuário.
	session string
	// lastRecord associa uma chave de snapshot ao instante da última captura.
	lastRecord map[string]time.Time
}

// appLog guarda o logger do processo atual.
var appLog struct {
	// mu protege a criação e a leitura do ponteiro.
	mu sync.Mutex
	// current contém a instância ativa, quando inicializada.
	current *applicationLog
}

// InitApplicationLog cria o logger antes da janela Fyne e registra o startup.
func InitApplicationLog() {
	appLog.mu.Lock()
	defer appLog.mu.Unlock()
	if appLog.current != nil {
		return
	}
	// dir resolve o cache do usuário com fallback para o temporário.
	dir, err := os.UserCacheDir()
	if err == nil {
		dir = filepath.Join(dir, "GeradorJornada", "logs")
	} else {
		dir = filepath.Join(os.TempDir(), "GeradorJornada", "logs")
	}
	// logger registra o startup; falha de gravação não impede a abertura do editor.
	logger, err := newApplicationLog(dir)
	if err != nil {
		return
	}
	appLog.current = logger
	// startup identifica binário e ambiente sem caminho do projeto.
	logger.write("startup", map[string]any{
		"pid": os.Getpid(), "session": logger.session, "version": buildinfo.Version(),
		"go": runtime.Version(), "os": runtime.GOOS, "arch": runtime.GOARCH,
		"executable": executableIdentity(),
	})
}

// ApplicationLogPath devolve o caminho efetivo do arquivo corrente.
func ApplicationLogPath() string {
	appLog.mu.Lock()
	defer appLog.mu.Unlock()
	if appLog.current == nil {
		return ""
	}
	return appLog.current.path
}

// LogApplication registra um evento com campos escalares seguros.
func LogApplication(event string, fields map[string]any) {
	appLog.mu.Lock()
	logger := appLog.current
	appLog.mu.Unlock()
	if logger == nil {
		return
	}
	logger.write(event, fields)
}

// LogApplicationSnapshot registra apenas contagens e identidades técnicas da tela atual.
// force ignora a janela de deduplicação; nenhum nome, texto ou caminho do projeto é gravado.
func LogApplicationSnapshot(document application.Document, screenID string, force bool) {
	// fields contém somente dados estruturais permitidos.
	fields := map[string]any{"screens": len(document.Screens), "revision": document.Project.Revision}
	// index localiza a tela corrente sem expor seu nome.
	for index := range document.Screens {
		if document.Screens[index].ID != screenID {
			continue
		}
		page := document.Screens[index]
		fields["components"] = len(page.Components)
		// kinds agrega a quantidade por família estável.
		kinds := map[string]int{}
		// enabled conta componentes interativos.
		enabled := 0
		// component percorre a composição.
		for _, component := range page.Components {
			kinds[string(component.Kind)]++
			if component.Enabled {
				enabled++
			}
		}
		fields["kinds"] = kinds
		fields["enabled"] = enabled
	}
	// key identifica a captura para deduplicação.
	key := fmt.Sprintf("%s/%d/%v", screenID, document.Project.Revision, fields["components"])
	LogApplicationDedup("component_snapshot", fields, key, force)
}

// LogApplicationDedup grava um evento no máximo uma vez a cada dois segundos por chave.
func LogApplicationDedup(event string, fields map[string]any, key string, force bool) {
	appLog.mu.Lock()
	logger := appLog.current
	appLog.mu.Unlock()
	if logger == nil {
		return
	}
	logger.writeDedup(event, fields, key, force, 2*time.Second)
}

// newApplicationLog abre a pasta, remove excedentes e cria o arquivo da sessão.
func newApplicationLog(dir string) (*applicationLog, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	// logger recebe a identidade da execução.
	logger := &applicationLog{dir: dir, limit: logSegmentLimit, session: project.NewID(), lastRecord: map[string]time.Time{}}
	// stamp ordena o nome do arquivo.
	stamp := time.Now().UTC().Format("20060102-150405")
	logger.path = filepath.Join(dir, fmt.Sprintf("jornada-%s-%d.jsonl", stamp, os.Getpid()))
	// file cria o arquivo exclusivo da sessão.
	file, err := os.OpenFile(logger.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	logger.file = file
	logger.retain()
	return logger, nil
}

// write serializa uma linha JSONL e sincroniza o arquivo.
func (l *applicationLog) write(event string, fields map[string]any) {
	l.writeDedup(event, fields, "", true, 0)
}

// writeDedup aplica a janela de deduplicação antes de escrever.
func (l *applicationLog) writeDedup(event string, fields map[string]any, key string, force bool, window time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return
	}
	if key != "" && !force && window > 0 {
		if last, ok := l.lastRecord[key]; ok && time.Since(last) < window {
			return
		}
		l.lastRecord[key] = time.Now()
	}
	// record contém somente o evento e os campos recebidos.
	record := map[string]any{"ts": time.Now().UTC().Format(time.RFC3339Nano), "event": event, "session": l.session}
	for name, value := range fields {
		record[name] = value
	}
	line, err := json.Marshal(record)
	if err != nil {
		return
	}
	line = append(line, '\n')
	if l.size+int64(len(line)) > l.limit {
		l.rotate()
	}
	if _, err = l.file.Write(line); err != nil {
		return
	}
	l.size += int64(len(line))
	_ = l.file.Sync()
}

// rotate fecha o segmento corrente e preserva até três backups.
func (l *applicationLog) rotate() {
	if l.file != nil {
		_ = l.file.Close()
	}
	// index percorre os backups do mais antigo ao mais novo.
	for index := 2; index >= 1; index-- {
		_ = os.Rename(fmt.Sprintf("%s.%d", l.path, index), fmt.Sprintf("%s.%d", l.path, index+1))
	}
	_ = os.Rename(l.path, l.path+".1")
	// file reabre o caminho canônico para a continuação da sessão.
	file, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		l.file = nil
		return
	}
	l.file = file
	l.size = 0
}

// retain remove arquivos próprios mais antigos, preservando os vinte mais recentes.
func (l *applicationLog) retain() {
	// entries reúne somente os arquivos criados por este logger.
	entries, err := filepath.Glob(filepath.Join(l.dir, "jornada-*.jsonl*"))
	if err != nil {
		return
	}
	if len(entries) <= logRetention {
		return
	}
	// sort ordena do mais antigo para o mais novo.
	sort.Slice(entries, func(i, j int) bool {
		// left e right comparam a data de modificação sem abrir os arquivos.
		left, leftErr := os.Stat(entries[i])
		right, rightErr := os.Stat(entries[j])
		if leftErr != nil || rightErr != nil {
			return entries[i] < entries[j]
		}
		return left.ModTime().Before(right.ModTime())
	})
	// remove o excedente, sem falhar quando um arquivo está bloqueado.
	for _, path := range entries[:len(entries)-logRetention] {
		_ = os.Remove(path)
	}
}

// CloseApplicationLog fecha o arquivo corrente ao encerrar o processo.
func CloseApplicationLog() {
	appLog.mu.Lock()
	defer appLog.mu.Unlock()
	if appLog.current != nil && appLog.current.file != nil {
		_ = appLog.current.file.Close()
		appLog.current.file = nil
	}
}

// executableIdentity devolve basename e SHA-256 do executável, sem caminho completo.
func executableIdentity() map[string]any {
	// path contém o caminho do processo corrente.
	path, err := os.Executable()
	if err != nil {
		return map[string]any{"basename": "", "sha256": ""}
	}
	// result preserva somente dados técnicos do binário.
	result := map[string]any{"basename": filepath.Base(path)}
	// file permite calcular o hash sem manter o executável aberto.
	file, err := os.Open(path)
	if err != nil {
		result["sha256"] = ""
		return result
	}
	defer file.Close()
	// hash calcula o resumo do binário em streaming.
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		result["sha256"] = ""
		return result
	}
	result["sha256"] = hex.EncodeToString(hash.Sum(nil))
	return result
}
