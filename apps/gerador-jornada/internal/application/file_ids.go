package application

import (
	"fmt"
	"strconv"
	"strings"

	"gerador-jornada.local/editor/internal/screen"
)

// uniqueFileID encontra uma sugestão disponível sem ultrapassar o limite do nome-base.
func uniqueFileID(base string, used map[string]bool) string {
	// index acrescenta um sufixo numérico somente durante criação, duplicação ou migração.
	for index := 0; ; index++ {
		// suffix diferencia cópias sem alterar o prefixo obrigatório de telas.
		suffix := ""
		if index > 0 {
			suffix = strconv.Itoa(index + 1)
		}
		// candidate reserva espaço para o sufixo dentro dos 40 caracteres.
		candidate := base[:min(len(base), screen.FileIDLimit-len(suffix))] + suffix
		if !used[strings.ToLower(candidate)] {
			return candidate
		}
	}
}

// UniqueScreenFileID sugere um ID livre no projeto para uma tela nova ou duplicada.
func UniqueScreenFileID(document Document, base string) string {
	// used registra nomes de arquivos existentes sem distinguir caixa.
	used := map[string]bool{}
	for _, page := range document.Screens {
		used[strings.ToLower(page.FileID)] = true
	}
	return uniqueFileID(screen.SuggestFileID(base, true), used)
}

// UniqueComponentFileID sugere um ID livre na tela para uma instância nova ou duplicada.
func UniqueComponentFileID(page screen.Screen, base string) string {
	// used limita as colisões ao contexto da tela; outras telas têm diretórios próprios.
	used := map[string]bool{}
	for _, component := range page.Components {
		used[strings.ToLower(component.FileID)] = true
	}
	return uniqueFileID(screen.SuggestFileID(base, false), used)
}

// EnsureFileIDs preenche somente IDs ausentes de projetos antigos, preservando UUIDs e valores explícitos.
// Colisões das sugestões recebem sufixos determinísticos; IDs explícitos nunca são renomeados.
func EnsureFileIDs(document *Document) {
	// pages reserva todos os nomes explícitos antes de criar sugestões para os ausentes.
	pages := map[string]bool{}
	for _, page := range document.Screens {
		if page.FileID != "" {
			pages[strings.ToLower(page.FileID)] = true
		}
	}
	for index := range document.Screens {
		// page é a composição em migração, sem alterar nome amigável ou referências.
		page := &document.Screens[index]
		if page.FileID == "" {
			page.FileID = uniqueFileID(screen.SuggestFileID(page.Name, true), pages)
			pages[strings.ToLower(page.FileID)] = true
		}
		// components reserva os IDs existentes no contexto local antes das sugestões.
		components := map[string]bool{}
		for _, component := range page.Components {
			if component.FileID != "" {
				components[strings.ToLower(component.FileID)] = true
			}
		}
		for componentIndex := range page.Components {
			// component recebe apenas o novo atributo ausente, sem trocar sua identidade interna.
			component := &page.Components[componentIndex]
			if component.FileID == "" {
				component.FileID = uniqueFileID(screen.SuggestFileID(component.Name, false), components)
				components[strings.ToLower(component.FileID)] = true
			}
		}
	}
}

// ValidateFileIDs valida nomes de saída e colisões sem exigir geometria ou navegação completas.
func (d Document) ValidateFileIDs() error {
	// pages protege os arquivos de tela em sistemas de arquivos Windows.
	pages := map[string]bool{}
	for _, page := range d.Screens {
		if err := screen.ValidateFileID(page.FileID, true); err != nil {
			return fmt.Errorf("tela %s: %w", page.Name, err)
		}
		if pages[strings.ToLower(page.FileID)] {
			return fmt.Errorf("ID de tela repetido no projeto: %s", page.FileID)
		}
		pages[strings.ToLower(page.FileID)] = true
		// components protege nomes dos arquivos de cada tela, sem impedir reuso em outra tela.
		components := map[string]bool{}
		for _, component := range page.Components {
			if err := screen.ValidateFileID(component.FileID, false); err != nil {
				return fmt.Errorf("%s / %s: %w", page.Name, component.Name, err)
			}
			if components[strings.ToLower(component.FileID)] {
				return fmt.Errorf("ID de componente repetido em %s: %s", page.Name, component.FileID)
			}
			components[strings.ToLower(component.FileID)] = true
		}
	}
	return nil
}
