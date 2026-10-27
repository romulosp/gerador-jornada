//go:build windows

// Package credential isola tokens no Windows Credential Manager.
package credential

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"gerador-jornada.local/editor/internal/application"
	"golang.org/x/sys/windows"
)

// Store acessa somente credenciais genéricas criadas pela aplicação.
type Store struct{}

// protectedRecord permite recusar blobs de versão desconhecida sem interpretar credenciais incorretamente.
type protectedRecord struct {
	// Version identifica o contrato privado do cofre.
	Version int `json:"version"`
	// Token permanece exclusivamente no armazenamento protegido.
	Token application.GitHubCredential `json:"token"`
}

// record corresponde a CREDENTIALW no ABI Windows x64.
type record struct {
	// Flags e Kind identificam o tipo genérico.
	Flags uint32
	// Kind seleciona CRED_TYPE_GENERIC sem envolver credenciais de domínio.
	Kind uint32
	// Target e Comment são strings UTF-16 alocadas pelo chamador.
	Target *uint16
	// Comment é o comentário opcional do registro nativo.
	Comment *uint16
	// Modified é FILETIME mantido pelo sistema.
	Modified windows.Filetime
	// Size e Blob identificam os bytes secretos.
	Size uint32
	// Blob aponta para os bytes secretos controlados por Size.
	Blob *byte
	// Persist define persistência local, sem roaming.
	Persist uint32
	// Count e Attributes permanecem vazios.
	Count uint32
	// Attributes fica vazio porque a aplicação não publica atributos adicionais.
	Attributes uintptr
	// Alias e User são metadados opcionais não utilizados.
	Alias *uint16
	// User é metadado opcional não utilizado na credencial genérica.
	User *uint16
}

// calls carrega somente as funções WinCred documentadas.
var calls = windows.NewLazySystemDLL("advapi32.dll")

// Load devolve o token apenas ao coordenador backend e libera o blob nativo.
// key identifica o registro exato; ausência, falha WinCred ou blob inválido retornam erro.
func (Store) Load(key string) (application.GitHubCredential, error) {
	// target permanece vivo durante a chamada nativa.
	target, err := windows.UTF16PtrFromString(key)
	if err != nil {
		return application.GitHubCredential{}, err
	}
	// pointer é liberado exclusivamente por CredFree.
	var pointer *record
	// result e failure distinguem ausência de falha nativa.
	result, _, failure := calls.NewProc("CredReadW").Call(uintptr(unsafe.Pointer(target)), 1, 0, uintptr(unsafe.Pointer(&pointer)))
	if result == 0 {
		return application.GitHubCredential{}, application.WrapSyncProblem("credentialStoreUnavailable", "Credencial indisponível; reconecte o GitHub.", failure)
	}
	defer calls.NewProc("CredFree").Call(uintptr(unsafe.Pointer(pointer)))
	// token nunca é serializado em preferências comuns.
	var token application.GitHubCredential
	if pointer.Size > 2560 {
		return token, fmt.Errorf("Credencial protegida inválida.")
	}
	// stored valida a versão antes de disponibilizar qualquer token ao coordenador.
	var stored protectedRecord
	if err = json.Unmarshal(unsafe.Slice(pointer.Blob, int(pointer.Size)), &stored); err != nil || stored.Version != 1 || stored.Token.AccessToken == "" {
		return token, fmt.Errorf("Credencial protegida ilegível; reconecte o GitHub.")
	}
	return stored.Token, nil
}

// Save publica o blob completo no cofre local do usuário.
// key identifica o destino e token fornece os bytes; exceder o limite ou falhar WinCred retorna erro.
func (Store) Save(key string, token application.GitHubCredential) error {
	// data limita o blob ao máximo documentado por WinCred.
	data, err := json.Marshal(protectedRecord{Version: 1, Token: token})
	if err != nil {
		return err
	}
	if len(data) > 2560 {
		return fmt.Errorf("Credencial excede o limite do cofre Windows.")
	}
	// target e entry são válidos durante toda a chamada.
	target, err := windows.UTF16PtrFromString(key)
	if err != nil {
		return err
	}
	// entry apresenta o blob e o destino ao ABI Windows durante a chamada.
	entry := record{Kind: 1, Target: target, Size: uint32(len(data)), Blob: &data[0], Persist: 2}
	// result não expõe o conteúdo do blob em erro.
	result, _, failure := calls.NewProc("CredWriteW").Call(uintptr(unsafe.Pointer(&entry)), 0)
	if result == 0 {
		return application.WrapSyncProblem("credentialStoreUnavailable", "Não foi possível proteger a credencial no Windows.", failure)
	}
	return nil
}

// Delete é idempotente e nunca remove namespaces alheios.
// key seleciona somente o registro solicitado; inexistência é sucesso e outras falhas são preservadas.
func (Store) Delete(key string) error {
	// target identifica exatamente a credencial solicitada pelo coordenador.
	target, err := windows.UTF16PtrFromString(key)
	if err != nil {
		return err
	}
	// result e failure preservam erro de remoção, sem declarar logout fictício.
	result, _, failure := calls.NewProc("CredDeleteW").Call(uintptr(unsafe.Pointer(target)), 1, 0)
	if result == 0 && !errors.Is(failure, windows.ERROR_NOT_FOUND) {
		return application.WrapSyncProblem("credentialStoreUnavailable", "Não foi possível remover a credencial protegida.", failure)
	}
	return nil
}
