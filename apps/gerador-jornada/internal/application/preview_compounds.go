package application

import (
	"fmt"

	"gerador-jornada.local/editor/internal/catalog"
	"gerador-jornada.local/editor/internal/screen"
)

// BarcodeFieldAllowed informa se o campo pode receber foco sem atravessar bloco inválido.
func (p *Preview) BarcodeFieldAllowed(componentID string, index int) bool {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil || runtime.Locked || runtime.Reading || index < 0 || index > 7 {
		return false
	}
	if index == 0 {
		return true
	}
	// blockValid verifica o par concluído que libera o próximo bloco.
	blockValid := func(block int) bool {
		return catalog.ValidateTypedBlock(runtime.Values[:], block) == nil
	}
	switch index {
	case 1:
		return len(runtime.Values[0]) == catalog.BarcodeFieldLimits[0]
	case 2:
		return blockValid(1)
	case 3:
		return blockValid(1) && len(runtime.Values[2]) == catalog.BarcodeFieldLimits[2]
	case 4:
		return blockValid(2)
	case 5:
		return blockValid(2) && len(runtime.Values[4]) == catalog.BarcodeFieldLimits[4]
	case 6:
		return blockValid(3)
	case 7:
		return blockValid(3) && len(runtime.Values[6]) == catalog.BarcodeFieldLimits[6]
	}
	return false
}

// Barcode devolve o estado transitório da entrada indicada.
func (p *Preview) Barcode(componentID string) (*BarcodeRuntime, bool) {
	// runtime preserva a identidade por componente na visita atual.
	runtime, ok := p.Barcodes[componentID]
	return runtime, ok
}

// barcodeComponent resolve a entrada e seu estado sem aceitar outro tipo.
func (p *Preview) barcodeComponent(componentID string) (*BarcodeRuntime, screen.Component, error) {
	// page identifica a composição corrente.
	page := p.Document.Screen(p.ScreenID)
	if page == nil {
		return nil, screen.Component{}, fmt.Errorf("tela inexistente")
	}
	// component precisa ser uma entrada de código de barras da visita.
	component := page.ComponentByID(componentID)
	if component == nil || component.Kind != catalog.BarcodeInput {
		return nil, screen.Component{}, fmt.Errorf("entrada de código de barras inexistente")
	}
	// runtime é criado sob demanda para preservar a sessão de documentos anteriores.
	runtime := p.Barcodes[componentID]
	if runtime == nil {
		runtime = &BarcodeRuntime{}
		p.Barcodes[componentID] = runtime
	}
	return runtime, *component, nil
}

// TypeBarcodeDigit acrescenta um dígito ASCII ao campo ativo ou à captura integral.
// Retorna erro somente para componente inexistente; erros de conteúdo são locais.
func (p *Preview) TypeBarcodeDigit(componentID string, value rune) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	if runtime.Locked || value < '0' || value > '9' {
		return nil
	}
	if runtime.Reading {
		if len(runtime.Capture) < 47 {
			runtime.Capture += string(value)
		}
		return nil
	}
	if !p.BarcodeFieldAllowed(componentID, runtime.Active) {
		return nil
	}
	// limit contém o máximo de dígitos do campo ativo.
	limit := catalog.BarcodeFieldLimits[runtime.Active]
	if len(runtime.Values[runtime.Active]) >= limit {
		return nil
	}
	runtime.Values[runtime.Active] += string(value)
	if len(runtime.Values[runtime.Active]) < limit {
		runtime.Message = ""
		return nil
	}
	// active é o campo que acabou de completar.
	active := runtime.Active
	if active == 1 || active == 3 || active == 5 {
		// block identifica o par 1, 2 ou 3 validado pelo DV local.
		block := active/2 + 1
		if err := catalog.ValidateTypedBlock(runtime.Values[:], block); err != nil {
			runtime.Message = err.Error()
			runtime.Active = active - 1
			return nil
		}
		runtime.Message = ""
	}
	if runtime.Active < 7 {
		runtime.Active++
	}
	return nil
}

// BarcodeBackspace remove o último caractere do campo ativo ou da captura.
func (p *Preview) BarcodeBackspace(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	if runtime.Locked {
		return nil
	}
	if runtime.Reading {
		if len(runtime.Capture) > 0 {
			// runes preserva a remoção de um único caractere.
			runes := []rune(runtime.Capture)
			runtime.Capture = string(runes[:len(runes)-1])
		}
		return nil
	}
	if len(runtime.Values[runtime.Active]) > 0 {
		// runes preserva a remoção de um único caractere.
		runes := []rune(runtime.Values[runtime.Active])
		runtime.Values[runtime.Active] = string(runes[:len(runes)-1])
	} else if runtime.Active > 0 {
		runtime.Active--
	}
	runtime.Message = ""
	return nil
}

// BarcodeClear limpa a digitação ou a captura sem apagar outras instâncias.
func (p *Preview) BarcodeClear(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil || runtime.Locked {
		return err
	}
	if runtime.Reading {
		runtime.Capture = ""
		return nil
	}
	runtime.Values = [8]string{}
	runtime.Active = 0
	runtime.Message = ""
	return nil
}

// BarcodeFinish conclui a entrada digitada, completando o último campo com zeros.
func (p *Preview) BarcodeFinish(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil || runtime.Locked {
		return err
	}
	// fields é uma cópia completável sem destruir o rascunho em caso de erro.
	fields := make([]string, 8)
	copy(fields, runtime.Values[:])
	fields = catalog.CompleteLastField(fields)
	if barcode, err := catalog.ParseTypedFields(fields); err == nil {
		copy(runtime.Values[:], fields)
		runtime.Locked = true
		runtime.Result = barcode
		runtime.Message = ""
		runtime.Active = 0
		return nil
	}
	runtime.Message = "Código de barras inválido"
	runtime.Active = 0
	return nil
}

// BarcodeRestart limpa todo o estado transitório da instância.
func (p *Preview) BarcodeRestart(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	*runtime = BarcodeRuntime{}
	return nil
}

// BarcodePaste valida uma captura integral antes de publicar qualquer campo.
func (p *Preview) BarcodePaste(componentID, raw string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	return p.applyBarcodeCapture(runtime, raw)
}

// BarcodeBeginRead inicia o modo explícito de leitura, sem inferir velocidade.
func (p *Preview) BarcodeBeginRead(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil || runtime.Locked {
		return err
	}
	runtime.Reading = true
	runtime.Capture = ""
	runtime.Message = ""
	runtime.Active = 0
	return nil
}

// BarcodeCancelRead cancela a captura antes de submeter e preserva o rascunho.
func (p *Preview) BarcodeCancelRead(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	runtime.Reading = false
	runtime.Capture = ""
	return nil
}

// BarcodeSubmitRead submete a captura integral acumulada pelo modo de leitura.
func (p *Preview) BarcodeSubmitRead(componentID, raw string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	runtime.Reading = false
	return p.applyBarcodeCapture(runtime, raw)
}

// applyBarcodeCapture aceita 44/47 dígitos ou limpa somente a instância destinatária.
func (p *Preview) applyBarcodeCapture(runtime *BarcodeRuntime, raw string) error {
	// barcode valida a captura antes de qualquer preenchimento parcial.
	barcode, err := catalog.ParseBarcode(raw)
	if err != nil {
		*runtime = BarcodeRuntime{Modal: true, Message: "Código de barras inválido"}
		return nil
	}
	// fields contém os oito valores derivados da linha digitável validada.
	fields, err := catalog.BarcodeFieldsFromDigitized(barcode.Digitized)
	if err != nil {
		*runtime = BarcodeRuntime{Modal: true, Message: "Código de barras inválido"}
		return nil
	}
	// runtime publica os valores somente depois da validação integral.
	*runtime = BarcodeRuntime{Locked: true, Result: barcode, Message: ""}
	copy(runtime.Values[:], fields)
	return nil
}

// BarcodeDismissModal fecha o diálogo de captura inválida e devolve o foco ao primeiro campo.
func (p *Preview) BarcodeDismissModal(componentID string) error {
	// runtime contém o estado da instância consultada.
	runtime, _, err := p.barcodeComponent(componentID)
	if err != nil {
		return err
	}
	runtime.Modal = false
	runtime.Message = ""
	runtime.Active = 0
	return nil
}

// MenuState devolve o caminho local da instância de menu indicada.
func (p *Preview) MenuState(componentID string) (*MenuRuntime, bool) {
	// runtime preserva a identidade por componente na visita atual.
	runtime, ok := p.Menus[componentID]
	return runtime, ok
}

// menuComponent resolve a instância e o estado sem aceitar outro tipo.
func (p *Preview) menuComponent(componentID string) (*MenuRuntime, screen.Component, error) {
	// page identifica a composição corrente.
	page := p.Document.Screen(p.ScreenID)
	if page == nil {
		return nil, screen.Component{}, fmt.Errorf("tela inexistente")
	}
	// component precisa ser um menu da aplicação da visita.
	component := page.ComponentByID(componentID)
	if component == nil || component.Kind != catalog.ApplicationMenu {
		return nil, screen.Component{}, fmt.Errorf("menu da aplicação inexistente")
	}
	// runtime é criado sob demanda para preservar a sessão de documentos anteriores.
	runtime := p.Menus[componentID]
	if runtime == nil {
		runtime = &MenuRuntime{Path: initialMenuPath(component.Properties.Menu)}
		p.Menus[componentID] = runtime
	}
	return runtime, *component, nil
}

// initialMenuPath devolve o caminho inicial confirmado do menu.
func initialMenuPath(menu *catalog.ApplicationMenuProperties) []string {
	if menu == nil {
		return nil
	}
	// path contém a cadeia de nós até o menu inicial.
	path := menu.PathTo(menu.InitialMenuID)
	if len(path) == 0 {
		return []string{menu.Root.ID}
	}
	// ids preserva somente as identidades dos nós.
	ids := make([]string, len(path))
	for index, node := range path {
		ids[index] = node.ID
	}
	return ids
}

// MenuNode devolve o submenu visível no estado atual.
func (p *Preview) MenuNode(componentID string) *catalog.SubmenuProperties {
	// runtime e component identificam a árvore consultada.
	runtime, component, err := p.menuComponent(componentID)
	if err != nil || component.Properties.Menu == nil {
		return nil
	}
	if len(runtime.Path) == 0 {
		return &component.Properties.Menu.Root
	}
	return component.Properties.Menu.FindSubmenu(runtime.Path[len(runtime.Path)-1])
}

// MenuActivate executa um item habilitado do nó visível, sem inferir rotas por texto.
func (p *Preview) MenuActivate(componentID, itemID string) (MenuOutcome, error) {
	// runtime identifica a árvore consultada na visita.
	runtime, component, err := p.menuComponent(componentID)
	if err != nil {
		return MenuNoop, err
	}
	if !component.Enabled || runtime.Exited {
		return MenuBlocked, nil
	}
	node := p.MenuNode(componentID)
	if node == nil {
		return MenuNoop, nil
	}
	// item é localizado pela identidade, nunca pela legenda.
	for _, item := range node.Items {
		if item.ID != itemID {
			continue
		}
		if !item.Enabled {
			return MenuBlocked, nil
		}
		if item.Submenu != nil {
			runtime.Path = append(runtime.Path, item.Submenu.ID)
			runtime.SimulatedItemID = ""
			return MenuOpened, nil
		}
		if item.TargetJourney != "" && p.AllowScreenNavigation {
			if err := p.EnterJourney(item.TargetJourney); err != nil {
				return MenuBlocked, err
			}
			return MenuNavigated, nil
		}
		if item.TargetScreenID != "" && p.AllowScreenNavigation {
			if p.Document.Screen(item.TargetScreenID) == nil {
				return MenuBlocked, fmt.Errorf("destino de menu inexistente")
			}
			if err := p.Enter(item.TargetScreenID); err != nil {
				return MenuBlocked, err
			}
			return MenuNavigated, nil
		}
		runtime.SimulatedItemID = item.ID
		return MenuSimulated, nil
	}
	return MenuNoop, nil
}

// MenuSelectedItem devolve a folha simulada do nó atual com seus metadados literais.
// Nome de movimento, código e opção são independentes; nenhum deles resolve ou executa movimento XML.
func (p *Preview) MenuSelectedItem(componentID string) (catalog.MenuItem, bool) {
	// runtime e node identificam a seleção local, sem procurar itens por código ou legenda.
	runtime, ok := p.Menus[componentID]
	if !ok || runtime.SimulatedItemID == "" {
		return catalog.MenuItem{}, false
	}
	// node conserva o contexto mesmo quando outras folhas compartilham o código de movimento.
	node := p.MenuNode(componentID)
	if node != nil {
		// item compara somente identidade confirmada dentro do nó visível.
		for _, item := range node.Items {
			if item.ID == runtime.SimulatedItemID && item.Submenu == nil {
				return item, true
			}
		}
	}
	return catalog.MenuItem{}, false
}

// MenuBack remove um nível do caminho local; na raiz não muda nada.
func (p *Preview) MenuBack(componentID string) error {
	// runtime contém o caminho local da instância.
	runtime, _, err := p.menuComponent(componentID)
	if err != nil {
		return err
	}
	if len(runtime.Path) > 1 {
		runtime.Path = runtime.Path[:len(runtime.Path)-1]
	}
	runtime.SimulatedItemID = ""
	return nil
}

// MenuCancel retorna diretamente à raiz, mesmo quando o inicial é outro submenu.
func (p *Preview) MenuCancel(componentID string) error {
	// runtime contém o caminho local da instância.
	runtime, _, err := p.menuComponent(componentID)
	if err != nil {
		return err
	}
	if len(runtime.Path) > 1 {
		runtime.Path = runtime.Path[:1]
	}
	runtime.SimulatedItemID = ""
	return nil
}

// MenuExit sinaliza o fechamento somente da sessão de prévia.
func (p *Preview) MenuExit(componentID string) error {
	// runtime contém o caminho local da instância.
	runtime, _, err := p.menuComponent(componentID)
	if err != nil {
		return err
	}
	runtime.Exited = true
	return nil
}
