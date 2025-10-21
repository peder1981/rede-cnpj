package cpfcnpj

import "testing"

func TestValidarCPF(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"CPF válido", "12345678909", "12345678909"},
		{"CPF com pontuação", "123.456.789-09", "12345678909"},
		{"CPF inválido", "12345678900", ""},
		{"CPF vazio", "", ""},
		{"CPF com zeros à esquerda", "00012345678909", "12345678909"},
		{"CPF curto", "123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarCPF(tt.input)
			if result != tt.expected {
				t.Errorf("ValidarCPF(%q) = %q, esperado %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidarCNPJ(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"CNPJ válido", "11222333000181", "11222333000181"},
		{"CNPJ com pontuação", "11.222.333/0001-81", "11222333000181"},
		{"CNPJ inválido", "11222333000180", ""},
		{"CNPJ vazio", "", ""},
		{"CNPJ com 8 dígitos (radical)", "11222333", "11222333000181"},
		{"CNPJ curto", "123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarCNPJ(tt.input)
			if result != tt.expected {
				t.Errorf("ValidarCNPJ(%q) = %q, esperado %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCNPJFormatado(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"CNPJ válido", "11222333000181", "11.222.333/0001-81"},
		{"CNPJ inválido", "123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CNPJFormatado(tt.input)
			if result != tt.expected {
				t.Errorf("CNPJFormatado(%q) = %q, esperado %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveCPFFinal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Nome com CPF", "JOAO DA SILVA12345678909", "JOAO DA SILVA"},
		{"Nome sem CPF", "JOAO DA SILVA", "JOAO DA SILVA"},
		{"Nome curto", "JOAO", "JOAO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveCPFFinal(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveCPFFinal(%q) = %q, esperado %q", tt.input, result, tt.expected)
			}
		})
	}
}
