package cpfcnpj

import (
	"regexp"
	"strconv"
	"strings"
)

var digitRegex = regexp.MustCompile(`\d`)

// ValidarCPF valida CPFs, retornando apenas a string de números válida
func ValidarCPF(cpf string) string {
	// Extrai apenas dígitos
	digits := digitRegex.FindAllString(cpf, -1)
	cpf = strings.Join(digits, "")

	if cpf == "" {
		return ""
	}

	// Remove zeros à esquerda extras
	if len(cpf) > 11 {
		zeros := strings.Repeat("0", len(cpf)-11)
		if cpf[:len(cpf)-11] == zeros {
			cpf = cpf[len(cpf)-11:]
		} else {
			return ""
		}
	}

	if len(cpf) < 3 {
		return ""
	}

	// Preenche com zeros à esquerda se necessário
	if len(cpf) < 11 {
		cpf = strings.Repeat("0", 11-len(cpf)) + cpf
	}

	// Converte para slice de inteiros
	numbers := make([]int, 11)
	for i, c := range cpf {
		numbers[i], _ = strconv.Atoi(string(c))
	}

	// Validação do primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		sum += numbers[i] * (10 - i)
	}
	expectedDigit := (sum * 10 % 11) % 10
	if numbers[9] != expectedDigit {
		return ""
	}

	// Validação do segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		sum += numbers[i] * (11 - i)
	}
	expectedDigit = (sum * 10 % 11) % 10
	if numbers[10] != expectedDigit {
		return ""
	}

	return cpf
}

// ValidarCNPJ valida CNPJs, retornando apenas a string de números válida
func ValidarCNPJ(cnpj string) string {
	// Extrai apenas dígitos
	digits := digitRegex.FindAllString(cnpj, -1)
	cnpj = strings.Join(digits, "")

	if cnpj == "" {
		return ""
	}

	cnpjOriginal := cnpj

	// Remove zeros à esquerda extras
	if len(cnpj) > 14 {
		zeros := strings.Repeat("0", len(cnpj)-14)
		if cnpj[:len(cnpj)-14] == zeros {
			cnpj = cnpj[len(cnpj)-14:]
		} else {
			return ""
		}
	}

	if len(cnpj) < 3 {
		return ""
	}

	// Se tem 8 dígitos (radical), adiciona 000100
	if len(cnpjOriginal) == 8 {
		cnpj += "000100"
	}

	// Preenche com zeros à esquerda
	if len(cnpj) < 14 {
		cnpj = strings.Repeat("0", 14-len(cnpj)) + cnpj
	}

	// Converte para slice de inteiros
	inteiros := make([]int, 14)
	for i, c := range cnpj {
		inteiros[i], _ = strconv.Atoi(string(c))
	}

	// Gera os 2 dígitos verificadores
	novo := make([]int, 12)
	copy(novo, inteiros[:12])

	prod := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	for len(novo) < 14 {
		r := 0
		for i, v := range novo {
			r += v * prod[i]
		}
		r = r % 11

		var f int
		if r > 1 {
			f = 11 - r
		} else {
			f = 0
		}

		novo = append(novo, f)
		prod = append([]int{6}, prod...)
	}

	// Verifica se o número gerado coincide com o original
	valid := true
	for i := range inteiros {
		if novo[i] != inteiros[i] {
			valid = false
			break
		}
	}

	if valid {
		return cnpj
	}

	// Se tinha 8 dígitos da matriz, retorna com dígitos verificadores corretos
	if len(cnpjOriginal) == 8 {
		result := cnpj[:12] + strconv.Itoa(novo[12]) + strconv.Itoa(novo[13])
		return result
	}

	return ""
}

// CNPJFormatado formata um CNPJ no padrão XX.XXX.XXX/XXXX-XX
func CNPJFormatado(cnpj string) string {
	if len(cnpj) != 14 {
		return cnpj
	}
	return cnpj[:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:]
}

// RemoveCPFFinal remove CPF do final do nome se existir
func RemoveCPFFinal(nome string) string {
	if len(nome) >= 11 {
		sufixo := nome[len(nome)-11:]
		if digitRegex.MatchString(sufixo) && len(digitRegex.FindAllString(sufixo, -1)) == 11 {
			return strings.TrimSpace(nome[:len(nome)-11])
		}
	}
	return nome
}
