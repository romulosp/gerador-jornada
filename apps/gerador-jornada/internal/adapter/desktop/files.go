package desktop

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"gerador-jornada.local/editor/internal/adapter/builtinassets"
	"gerador-jornada.local/editor/internal/adapter/projectzip"
	"gerador-jornada.local/editor/internal/application"
	"gerador-jornada.local/editor/internal/project"
)

// confirmDiscard protege alterações ainda não salvas ao substituir a sessão.
func (e *Editor) confirmDiscard(action func()) {
	if e.Session.Dirty() {
		dialog.ShowConfirm("Alterações não salvas", "Continuar sem salvar as alterações atuais?", func(yes bool) {
			if yes {
				action()
			}
		}, e.Window)
	} else {
		action()
	}
}

// newProject oferece um documento em branco ou o exemplo completo.
func (e *Editor) newProject() {
	e.confirmDiscard(func() {
		// name é o título inicial escolhido pelo autor.
		name := entry("Minha jornada")
		// example permite começar a partir de um conjunto de referência editável.
		example := widget.NewCheck("Iniciar com uma jornada de exemplo", nil)
		dialog.ShowCustomConfirm("Novo projeto", "Criar", "Cancelar", container.NewVBox(widget.NewForm(widget.NewFormItem("Nome", name)), example), func(yes bool) {
			if !yes || strings.TrimSpace(name.Text) == "" {
				return
			}
			// document é criado sem reutilizar IDs da sessão anterior.
			document := application.NewDocument(name.Text)
			if example.Checked {
				document = application.Example()
				document.Project.Name = name.Text
			}
			builtinassets.DefaultBackgrounds(&document)
			e.loadDocument(document, "")
		}, e.Window)
	})
}

// loadDocument troca a sessão somente após um carregamento bem-sucedido.
func (e *Editor) loadDocument(document application.Document, path string) {
	if e.Autosave != nil {
		e.Autosave.Stop()
	}
	e.Session = application.NewSession(document)
	e.FilePath = path
	e.Selected = ""
	e.PageID = ""
	if len(document.Screens) > 0 {
		e.PageID = document.Journey.StartScreenID
		if document.Screen(e.PageID) == nil {
			e.PageID = document.Screens[0].ID
		}
	}
	if path != "" {
		e.SetProjectPath(path)
	} else {
		e.refreshWorkspace()
	}
	e.refresh()
}

// openProject usa o seletor nativo apenas para leitura do arquivo escolhido.
func (e *Editor) openProject() {
	e.confirmDiscard(func() {
		// chooser mantém o arquivo original intacto durante validação.
		chooser := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if !e.check(err) || reader == nil {
				return
			}
			defer reader.Close()
			// data é limitada ao tamanho comprimido permitido.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			data, err := io.ReadAll(io.LimitReader(reader, (250<<20)+1))
			if !e.check(err) {
				return
			}
			// document só substitui a sessão se o pacote for estruturalmente legível.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			document, err := projectzip.Decode(data)
			if !e.check(err) {
				return
			}
			e.loadDocument(document, reader.URI().Path())
			// problems reúne os diagnósticos que impedem execução ou exportação.
			if problems := document.Validate(); len(problems) > 0 {
				dialog.ShowInformation("Projeto aberto como rascunho", application.ProblemsText(problems), e.Window)
			}
		}, e.Window)
		chooser.SetFilter(storage.NewExtensionFileFilter([]string{".jornada"}))
		chooser.Show()
	})
}

// chooseOutput solicita um caminho sem truncar um arquivo existente pelo diálogo.
func (e *Editor) chooseOutput(title, extension string, accept func(string)) {
	e.chooseOutputResult(title, extension, accept, func() {})
}

// chooseOutputResult notifica cancelamento sem abrir ou truncar o arquivo de destino.
func (e *Editor) chooseOutputResult(title, extension string, accept func(string), canceled func()) {
	// directory é um destino local selecionável.
	directory, _ := os.Getwd()
	if e.FilePath != "" {
		directory = filepath.Dir(e.FilePath)
	} else if e.WorkspacePath != "" {
		directory = e.WorkspacePath
	}
	// path é editável e não abre um stream de gravação.
	path := entry(filepath.Join(directory, "minha-jornada"+extension))
	// browse permite escolher uma pasta sem criar arquivos.
	browse := widget.NewButton("Escolher pasta…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if e.check(err) && uri != nil {
				path.SetText(filepath.Join(uri.Path(), "minha-jornada"+extension))
			}
		}, e.Window)
	})
	dialog.ShowCustomConfirm(title, "Continuar", "Cancelar", container.NewVBox(widget.NewForm(widget.NewFormItem("Arquivo", path)), browse), func(yes bool) {
		if !yes {
			canceled()
			return
		}
		// destination inclui a extensão esperada sem alterar nomes já compatíveis.
		destination := path.Text
		if !strings.HasSuffix(strings.ToLower(destination), extension) {
			destination += extension
		}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if _, err := os.Stat(destination); err == nil {
			dialog.ShowConfirm("Substituir arquivo?", "O destino já existe. Substituir somente após concluir a gravação?", func(replace bool) {
				if replace {
					accept(destination)
				} else {
					canceled()
				}
			}, e.Window)
		} else {
			accept(destination)
		}
	}, e.Window)
}

// saveProject publica snapshot explícito sem marcar edições posteriores como salvas.
func (e *Editor) saveProject(saveAs bool) {
	e.saveWithCompletion(saveAs, func(error) {})
}

// errSaveCanceled distingue cancelamento de falha de gravação sem apresentá-lo como erro técnico.
var errSaveCanceled = errors.New("Salvamento cancelado")

// saveWithCompletion só libera publicação após gravação íntegra e confirmação de migração.
func (e *Editor) saveWithCompletion(saveAs bool, completed func(error)) {
	if !e.gitCanWrite() {
		completed(errSaveCanceled)
		return
	}
	// save captura todos os dados antes de iniciar IO.
	save := func(path string) {
		// snapshot é independente do estado que pode continuar sendo editado.
		snapshot := e.Session.Document.Copy()
		// err conserva o resultado para o chamador distinguir falha de cancelamento.
		err := projectzip.Save(context.Background(), path, snapshot)
		if e.check(err) {
			e.Session.Document.Project.CatalogVersion = 4
			e.SetProjectPath(path)
			e.Session.SavedRevision = snapshot.Project.Revision
			e.Status.SetText("Projeto salvo · " + path)
		}
		completed(err)
	}
	// confirmed avisa antes da primeira escrita incompatível e conserva o arquivo legado no codec.
	confirmed := func(path string) {
		if _, err := os.Stat(path); err == nil && project.RequiresCatalogMigration(e.Session.Document.Project.CatalogVersion) {
			dialog.ShowConfirm("Atualizar formato do projeto", "O arquivo será salvo no catálogo 4. Uma cópia do original será preservada antes da gravação. Continuar?", func(yes bool) {
				if yes {
					save(path)
				} else {
					completed(errSaveCanceled)
				}
			}, e.Window)
		} else {
			save(path)
		}
	}
	if e.FilePath != "" && !saveAs {
		confirmed(e.FilePath)
	} else if e.WorkspacePath != "" && !saveAs {
		// workspace já foi escolhido; o primeiro salvamento não precisa pedir a pasta novamente.
		workspace, err := e.workspace()
		if !e.check(err) {
			completed(err)
			return
		}
		// filename é exclusivo, preservando qualquer jornada homônima já existente.
		filename, err := JourneyFilename(e.Session.Document.Project.Name)
		if err == nil {
			filename, err = workspace.uniqueName(filename)
		}
		if !e.check(err) {
			completed(err)
			return
		}
		confirmed(filepath.Join(workspace.Root, filename))
	} else {
		e.chooseOutputResult("Salvar projeto", ".jornada", confirmed, func() { completed(errSaveCanceled) })
	}
}

// scheduleRecovery captura a revisão na UI antes de iniciar o timer de autosave.
func (e *Editor) scheduleRecovery() {
	if e.Autosave != nil {
		e.Autosave.Stop()
	}
	// snapshot não compartilha estado com o callback do timer.
	snapshot := e.Session.Document.Copy()
	e.Autosave = time.AfterFunc(2*time.Second, func() {
		// directory mantém a recuperação separada do arquivo do usuário.
		directory, err := recoveryDirectory()
		if err != nil {
			fyne.Do(func() { e.Status.SetText("Recuperação automática indisponível: " + err.Error()) })
			return
		}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if err := os.MkdirAll(directory, 0700); err != nil {
			fyne.Do(func() { e.Status.SetText("Recuperação automática indisponível: " + err.Error()) })
			return
		}
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		if err := projectzip.Save(context.Background(), filepath.Join(directory, snapshot.Project.ID+".jornada"), snapshot); err != nil {
			fyne.Do(func() { e.Status.SetText("Falha na recuperação automática: " + err.Error()) })
		}
	})
}

// OfferRecovery oferece a cópia válida mais recente após inicializar a janela.
func (e *Editor) OfferRecovery() {
	// directory contém apenas cópias automáticas de projeto.
	directory, err := recoveryDirectory()
	if err != nil {
		e.Status.SetText("Recuperação automática indisponível: " + err.Error())
		return
	}
	// files é ordenado por modificação sem alterar os arquivos.
	files, _ := os.ReadDir(directory)
	sort.Slice(files, func(i, j int) bool {
		// left fornece os metadados do primeiro arquivo para ordenar as recuperações por data.
		left, _ := files[i].Info()
		// right fornece os metadados do segundo arquivo para ordenar as recuperações por data.
		right, _ := files[j].Info()
		return left.ModTime().After(right.ModTime())
	})
	// file é a próxima cópia candidata a recuperação.
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jornada") {
			continue
		}
		// document só é oferecido quando o pacote está íntegro.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		document, err := projectzip.Load(filepath.Join(directory, file.Name()))
		if err != nil {
			continue
		}
		dialog.ShowConfirm("Recuperar trabalho?", "Há uma cópia de recuperação de “"+document.Project.Name+"”. Abrir como projeto não salvo?", func(yes bool) {
			if yes {
				e.loadDocument(document, "")
				e.Session.SavedRevision = ^uint64(0)
				e.Status.SetText("Recuperação aberta · salve para escolher um destino")
			}
		}, e.Window)
		break
	}
}

// importImage incorpora bytes e metadados, nunca uma dependência de caminho externo.
func (e *Editor) importImage(componentID string, background bool) {
	// chooser aceita apenas formatos do contrato.
	chooser := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if !e.check(err) || reader == nil {
			return
		}
		defer reader.Close()
		// data é limitada antes de decodificar a imagem.
		// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
		data, err := io.ReadAll(io.LimitReader(reader, (25<<20)+1))
		if !e.check(err) {
			return
		}
		if len(data) > 25<<20 {
			e.check(fmt.Errorf("imagem excede 25 MiB"))
			return
		}
		// id e digest tornam o recurso independente do nome original.
		id := project.NewID()
		// digest calcula o SHA-256 para identificar o conteúdo da imagem incorporada.
		digest := sha256.Sum256(data)
		// extension e mediaType determinam o formato incorporado.
		extension := ".png"
		// mediaType registra o formato MIME da imagem importada.
		mediaType := "image/png"
		if strings.HasSuffix(strings.ToLower(reader.URI().Name()), ".jpg") || strings.HasSuffix(strings.ToLower(reader.URI().Name()), ".jpeg") {
			extension = ".jpg"
			mediaType = "image/jpeg"
		}
		// asset fica disponível apenas no snapshot provisório.
		asset := project.Asset{ID: id, Path: "assets/" + id + extension, MediaType: mediaType, SHA256: hex.EncodeToString(digest[:]), SizeBytes: int64(len(data))}
		if e.check(e.Session.Change(func(document *application.Document) error {
			document.Project.Assets = append(document.Project.Assets, asset)
			document.Resources[asset.Path] = data
			if background {
				document.Screen(e.PageID).Background.AssetID = id
			} else {
				document.Screen(e.PageID).ComponentByID(componentID).Style.BackgroundAssetID = id
			}
			// validation usa o mesmo leitor de recursos do formato portátil.
			// err registra uma falha da operação para interromper o fluxo sem publicar estado parcial.
			// packageBytes contém o pacote portátil completo antes da publicação.
			packageBytes, err := projectzip.Encode(*document)
			if err != nil {
				return err
			}
			_, err = projectzip.Decode(packageBytes)
			return err
		})) {
			e.changed()
		}
	}, e.Window)
	chooser.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	chooser.Show()
}
