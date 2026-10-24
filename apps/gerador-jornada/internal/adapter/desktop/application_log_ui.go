package desktop

import (
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showApplicationLog apresenta o caminho efetivo e as ações de coleta do log.
// O popup captura somente contagens técnicas; nenhum texto do projeto é transcrito.
func (e *Editor) showApplicationLog() {
	// Init garante o arquivo mesmo quando o popup é aberto após uma falha de startup.
	InitApplicationLog()
	// path contém o destino efetivo mostrado ao usuário.
	path := ApplicationLogPath()
	// status informa o resultado da última ação sem fechar o diálogo.
	status := widget.NewLabel("")
	status.Wrapping = fyne.TextWrapWord
	// content reúne caminho, aviso de privacidade e mensagens.
	content := container.NewVBox()
	if path == "" {
		content.Add(widget.NewLabel("Log de aplicação indisponível nesta execução."))
	} else {
		content.Add(widget.NewLabel("Arquivo efetivo:"))
		content.Add(widget.NewLabel(path))
		content.Add(widget.NewLabel("Somente eventos técnicos e contagens são registrados; textos do projeto, nomes e caminhos de arquivo não são gravados."))
	}
	content.Add(status)
	// popup permite repetir a captura ou abrir a pasta, além de copiar o caminho.
	popup := dialog.NewCustom("Log de aplicação", "Fechar", content, e.Window)
	if path != "" {
		// capture força uma nova fotografia das contagens da tela atual.
		capture := widget.NewButton("Registrar diagnóstico agora", func() {
			LogApplicationSnapshot(e.Session.Document, e.PageID, true)
			status.SetText("Diagnóstico registrado no arquivo acima.")
		})
		// copyPath facilita colar o caminho em uma conversa ou ferramenta.
		copyPath := widget.NewButton("Copiar caminho", func() {
			e.App.Clipboard().SetContent(path)
			status.SetText("Caminho copiado.")
		})
		// openFolder abre somente a pasta do log no explorador do sistema.
		openFolder := widget.NewButton("Abrir pasta", func() {
			if err := exec.Command("explorer", filepath.Dir(path)).Start(); err != nil {
				status.SetText("Não foi possível abrir a pasta: " + err.Error())
			}
		})
		content.Add(container.NewGridWithColumns(3, capture, copyPath, openFolder))
	}
	// Abrir o popup já registra a seleção/tela atual, conforme o contrato.
	LogApplicationSnapshot(e.Session.Document, e.PageID, true)
	popup.Show()
}
