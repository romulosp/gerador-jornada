package catalog

import (
	"fmt"
	"strconv"
	"strings"
)

// BarcodeDefaultLabel contém o rótulo padrão da entrada de linha digitável.
const BarcodeDefaultLabel = "Número do Código de Barras:"

// BarcodeInputProperties define a configuração persistível da entrada de código de barras.
type BarcodeInputProperties struct {
	// Label contém o rótulo Unicode exibido acima dos oito campos; pode ser vazio.
	Label string
}

// NewBarcodeInputProperties cria a configuração padrão de uma entrada vazia.
func NewBarcodeInputProperties() *BarcodeInputProperties {
	return &BarcodeInputProperties{Label: BarcodeDefaultLabel}
}

// Validate aceita qualquer rótulo Unicode, inclusive vazio, sem validar valores de runtime.
func (p BarcodeInputProperties) Validate() error {
	return nil
}

// Barcode reúne as duas representações validadas e os dados derivados de um boleto.
type Barcode struct {
	// Read contém os 44 dígitos do código de barras de cobrança.
	Read string
	// Digitized contém os 47 dígitos da linha digitável correspondente.
	Digitized string
	// Bank contém os três dígitos do banco.
	Bank string
	// ValueCents contém o valor em centavos exatos, sem arredondamento.
	ValueCents int64
	// Factor contém o fator de vencimento numérico, sem convertê-lo em data.
	Factor int
	// HasFactor informa se o fator é diferente de zero.
	HasFactor bool
	// IsCaixa informa se o banco é a Caixa Econômica Federal (104).
	IsCaixa bool
	// Sigcb informa a classificação experimental sob a premissa SISPL local.
	Sigcb bool
	// TitleNumber contém o número do título somente quando Sigcb é verdadeiro.
	TitleNumber string
}

// NormalizeBarcode remove somente espaços ASCII e pontos de apresentação.
// Enter final, letras, sinais, quebras internas e demais símbolos são recusados.
func NormalizeBarcode(raw string) (string, error) {
	// cleaned contém a representação canônica enquanto caracteres de apresentação são removidos.
	cleaned := strings.Map(func(value rune) rune {
		if value == ' ' || value == '.' {
			return -1
		}
		return value
	}, raw)
	if len(cleaned) != 44 && len(cleaned) != 47 {
		return "", fmt.Errorf("código de barras deve ter 44 ou 47 dígitos")
	}
	// index percorre cada posição para garantir conteúdo ASCII numérico.
	for index := 0; index < len(cleaned); index++ {
		if cleaned[index] < '0' || cleaned[index] > '9' {
			return "", fmt.Errorf("código de barras aceita somente dígitos")
		}
	}
	return cleaned, nil
}

// ParseBarcode valida as duas representações e devolve os dados derivados.
// Entradas de 44 dígitos validam o DV geral; entradas de 47 validam os três DVs locais e o geral.
func ParseBarcode(raw string) (Barcode, error) {
	// normalized contém somente dígitos e um comprimento aceito.
	normalized, err := NormalizeBarcode(raw)
	if err != nil {
		return Barcode{}, err
	}
	if len(normalized) == 44 {
		if !validGeneralDV(normalized) {
			return Barcode{}, fmt.Errorf("DV geral inválido")
		}
		return deriveBarcode(normalized, readToDigitized(normalized)), nil
	}
	// digitized preserva a linha informada e suas posições de DV.
	digitized := normalized
	// localBlocks contém os três pares base/DV da linha digitável.
	localBlocks := []struct {
		// base contém os dígitos sem o DV local.
		base string
		// dv contém o dígito verificador informado.
		dv byte
	}{
		{digitized[0:9], digitized[9]},
		{digitized[10:20], digitized[20]},
		{digitized[21:31], digitized[31]},
	}
	// block enumera os três blocos para localizar um DV local inválido.
	for index, block := range localBlocks {
		if byte('0'+mod10(block.base)) != block.dv {
			return Barcode{}, fmt.Errorf("DV local do bloco %d inválido", index+1)
		}
	}
	// read contém a representação de 44 posições convertida da linha digitável.
	read := digitizedToRead(digitized)
	if !validGeneralDV(read) {
		return Barcode{}, fmt.Errorf("DV geral inválido")
	}
	return deriveBarcode(read, digitized), nil
}

// deriveBarcode monta o resultado aceito sem consultar serviços ou relógio.
func deriveBarcode(read, digitized string) Barcode {
	// cents contém o valor em centavos convertido de dez dígitos exatos.
	cents, _ := strconv.ParseInt(read[9:19], 10, 64)
	// factor contém o fator de vencimento como número, sem interpretar o ciclo.
	factor, _ := strconv.Atoi(read[5:9])
	// result reúne as representações e os derivados somente após a validação.
	result := Barcode{
		Read: read, Digitized: digitized, Bank: read[0:3], ValueCents: cents,
		Factor: factor, HasFactor: factor != 0, IsCaixa: read[0:3] == "104",
	}
	if result.IsCaixa && (read[29] == '1' || read[29] == '2') && read[33] >= '1' && read[33] <= '5' {
		// firstBase e firstDigit reproduzem o recorte lido 20..26.
		firstBase, firstDigit := read[19:25], int(read[25]-'0')
		// secondBase e secondDigit reproduzem o recorte lido 20..44.
		secondBase, secondDigit := read[19:43], int(read[43]-'0')
		if mod11Sispl(firstBase) == firstDigit && mod11Sispl(secondBase) == secondDigit {
			result.Sigcb = true
			result.TitleNumber = read[26:44]
		}
	}
	return result
}

// readToDigitized converte os 44 dígitos em 47 sem alterar os valores.
func readToDigitized(read string) string {
	// c1Base reúne as posições 1..4 e 20..24 do código de barras.
	c1Base := read[0:4] + read[19:24]
	// c2Base e c3Base reúnem os dez dígitos de cada bloco seguinte.
	c2Base := read[24:34]
	c3Base := read[34:44]
	return c1Base + strconv.Itoa(mod10(c1Base)) + c2Base + strconv.Itoa(mod10(c2Base)) + c3Base + strconv.Itoa(mod10(c3Base)) + read[4:5] + read[5:19]
}

// digitizedToRead converte os 47 dígitos em 44 removendo os DVs locais.
func digitizedToRead(digitized string) string {
	// c1 contém o primeiro bloco completo e suas posições de base.
	c1 := digitized[0:10]
	return c1[0:4] + digitized[32:33] + digitized[33:47] + c1[4:9] + digitized[10:20] + digitized[21:31]
}

// validGeneralDV compara o DV geral do boleto com a posição 5 do código lido.
func validGeneralDV(read string) bool {
	return mod11Boleto(read) == int(read[4]-'0')
}

// mod10 calcula o dígito verificador módulo 10 de uma base sem o DV.
func mod10(base string) int {
	// sum acumula os algarismos de cada produto, não o produto inteiro.
	sum := 0
	// weight alterna 2 e 1 da direita para a esquerda.
	weight := 2
	// index percorre a base da direita para a esquerda.
	for index := len(base) - 1; index >= 0; index-- {
		// product contém o valor ponderado do algarismo atual.
		product := int(base[index]-'0') * weight
		sum += product%10 + product/10
		if weight == 2 {
			weight = 1
		} else {
			weight = 2
		}
	}
	return (10 - sum%10) % 10
}

// mod11Boleto calcula o DV geral removendo a posição 5 e aplicando pesos 2..9.
func mod11Boleto(read string) int {
	// digits contém as 43 posições que participam do DV geral.
	digits := read[0:4] + read[5:44]
	// candidate contém o resultado antes da normalização dos casos especiais.
	candidate := 11 - weightedSum(digits)%11
	if candidate == 0 || candidate == 10 || candidate == 11 {
		return 1
	}
	return candidate
}

// mod11Sispl calcula o DV experimental SISPL sob a premissa local documentada.
func mod11Sispl(base string) int {
	// candidate contém o resultado antes da normalização SISPL.
	candidate := 11 - weightedSum(base)%11
	if candidate == 10 || candidate == 11 {
		return 0
	}
	return candidate
}

// weightedSum soma os algarismos com pesos cíclicos de 2 a 9 da direita para a esquerda.
func weightedSum(digits string) int {
	// sum acumula os produtos ponderados.
	sum := 0
	// weight percorre o ciclo 2..9.
	weight := 2
	// index percorre a base da direita para a esquerda.
	for index := len(digits) - 1; index >= 0; index-- {
		sum += int(digits[index]-'0') * weight
		weight++
		if weight > 9 {
			weight = 2
		}
	}
	return sum
}

// BarcodeFieldLimits contém os limites dos oito campos visuais na ordem exibida.
// Os seis primeiros formam três pares de validação; o sétimo é o DV geral e o oitavo é fator/valor.
var BarcodeFieldLimits = [8]int{5, 5, 5, 6, 5, 6, 1, 14}

// BarcodeFieldsFromDigitized divide uma linha digitável de 47 posições nos oito campos visuais.
func BarcodeFieldsFromDigitized(digitized string) ([]string, error) {
	// normalized valida comprimento e conteúdo antes de recortar.
	normalized, err := NormalizeBarcode(digitized)
	if err != nil || len(normalized) != 47 {
		return nil, fmt.Errorf("linha digitável de 47 dígitos é necessária")
	}
	return []string{
		normalized[0:5], normalized[5:10], normalized[10:15], normalized[15:21],
		normalized[21:26], normalized[26:32], normalized[32:33], normalized[33:47],
	}, nil
}

// BarcodeFieldsFromRead converte o código lido de 44 posições nos oito campos visuais.
func BarcodeFieldsFromRead(read string) ([]string, error) {
	// barcode valida o DV geral e devolve a linha digitável correspondente.
	barcode, err := ParseBarcode(read)
	if err != nil {
		return nil, err
	}
	return BarcodeFieldsFromDigitized(barcode.Digitized)
}

// ValidateTypedBlock confere o DV local do par 1, 2 ou 3 assim que o segundo campo completa.
// O bloqueio não apaga os blocos anteriores já válidos.
func ValidateTypedBlock(fields []string, block int) error {
	if len(fields) != 8 || block < 1 || block > 3 {
		return fmt.Errorf("par de campos inválido")
	}
	// first e second identificam o par visual do bloco lógico.
	first, second := 2*(block-1), 2*(block-1)+1
	if len(fields[first]) > BarcodeFieldLimits[first] || len(fields[second]) > BarcodeFieldLimits[second] {
		return fmt.Errorf("comprimento do bloco %d excede o limite", block)
	}
	// combined contém os dois campos sem separador.
	combined := fields[first] + fields[second]
	// expected repete o tamanho lógico do bloco: 10, 11 e 11 posições.
	expected := []int{10, 11, 11}[block-1]
	if len(combined) != expected {
		return fmt.Errorf("bloco %d incompleto", block)
	}
	// base exclui o DV local e dv contém o último dígito informado.
	base, dv := combined[:len(combined)-1], int(combined[len(combined)-1]-'0')
	if mod10(base) != dv {
		return fmt.Errorf("Bloco %d inválido", block)
	}
	return nil
}

// CompleteLastField preenche o campo de fator/valor com zeros à direita até 14 posições.
// A conclusão explícita é a única operação autorizada a completar o campo.
func CompleteLastField(fields []string) []string {
	// result separa a entrada antes de completar o último campo.
	result := append([]string{}, fields...)
	if len(result) == 8 && len(result[7]) < 14 {
		result[7] += strings.Repeat("0", 14-len(result[7]))
	}
	return result
}

// ParseTypedFields valida a entrada digitada nos oito campos e devolve os dados derivados.
func ParseTypedFields(fields []string) (Barcode, error) {
	if len(fields) != 8 {
		return Barcode{}, fmt.Errorf("oito campos são necessários")
	}
	// value acumula a linha digitável reconstruída.
	value := ""
	// index percorre os oito campos visuais.
	for index, field := range fields {
		if len(field) > BarcodeFieldLimits[index] {
			return Barcode{}, fmt.Errorf("campo %d excede o limite", index+1)
		}
		value += field
	}
	if len(value) != 47 {
		return Barcode{}, fmt.Errorf("entrada incompleta")
	}
	return ParseBarcode(value)
}

// BarcodeLabelText devolve o texto visível do tipo de rótulo escolhido.
func BarcodeLabelText(kind string) (string, bool) {
	switch kind {
	case "cpf":
		return "CPF:", true
	case "cnpj":
		return "CNPJ:", true
	case "name":
		return "Nome:", true
	case "tradeName":
		return "Nome fantasia:", true
	case "corporateName":
		return "Razão social:", true
	}
	return "", false
}
