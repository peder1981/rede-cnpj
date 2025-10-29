package main

import (
	"database/sql"
	"testing"
)

func TestNormalizerDate(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:     "data_teste",
		Type:     FieldTypeDate,
		Required: false,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "Data válida YYYYMMDD",
			input:    "20231015",
			expected: sql.NullString{String: "2023-10-15", Valid: true},
		},
		{
			name:     "Data inválida zero",
			input:    "0",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Data vazia",
			input:    "",
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("data_teste", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNormalizerCNPJ(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:      "cnpj",
		Type:      FieldTypeCNPJ,
		MaxLength: 14,
		Required:  true,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "CNPJ válido com máscara",
			input:    "12.345.678/0001-90",
			expected: sql.NullString{String: "12345678000190", Valid: true},
		},
		{
			name:     "CNPJ válido sem máscara",
			input:    "12345678000190",
			expected: sql.NullString{String: "12345678000190", Valid: true},
		},
		{
			name:     "CNPJ inválido (tamanho errado)",
			input:    "123456",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "CNPJ todo zeros",
			input:    "00000000000000",
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("cnpj", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNormalizerCEP(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:      "cep",
		Type:      FieldTypeCEP,
		MaxLength: 8,
		Required:  false,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "CEP válido com máscara",
			input:    "01310-100",
			expected: sql.NullString{String: "01310100", Valid: true},
		},
		{
			name:     "CEP válido sem máscara",
			input:    "01310100",
			expected: sql.NullString{String: "01310100", Valid: true},
		},
		{
			name:     "CEP inválido (tamanho errado)",
			input:    "123",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "CEP todo zeros",
			input:    "00000000",
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("cep", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNormalizerUF(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:      "uf",
		Type:      FieldTypeUF,
		MaxLength: 2,
		Required:  true,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "UF válida maiúscula",
			input:    "SP",
			expected: sql.NullString{String: "SP", Valid: true},
		},
		{
			name:     "UF válida minúscula",
			input:    "sp",
			expected: sql.NullString{String: "SP", Valid: true},
		},
		{
			name:     "UF inválida",
			input:    "XX",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "UF exterior",
			input:    "EX",
			expected: sql.NullString{String: "EX", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("uf", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNormalizerEmail(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:     "email",
		Type:     FieldTypeEmail,
		Required: false,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "Email válido",
			input:    "teste@exemplo.com",
			expected: sql.NullString{String: "teste@exemplo.com", Valid: true},
		},
		{
			name:     "Email válido maiúscula",
			input:    "TESTE@EXEMPLO.COM",
			expected: sql.NullString{String: "teste@exemplo.com", Valid: true},
		},
		{
			name:     "Email inválido (sem @)",
			input:    "teste.exemplo.com",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Email vazio",
			input:    "",
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("email", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNormalizerPhone(t *testing.T) {
	n := NewNormalizer()
	n.RegisterField(FieldMetadata{
		Name:      "telefone",
		Type:      FieldTypePhone,
		MaxLength: 11,
		Required:  false,
	})

	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "Telefone válido com máscara",
			input:    "(11) 98765-4321",
			expected: sql.NullString{String: "11987654321", Valid: true},
		},
		{
			name:     "Telefone válido sem máscara",
			input:    "11987654321",
			expected: sql.NullString{String: "11987654321", Valid: true},
		},
		{
			name:     "Telefone fixo",
			input:    "1133334444",
			expected: sql.NullString{String: "1133334444", Valid: true},
		},
		{
			name:     "Telefone inválido (muito curto)",
			input:    "123",
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Telefone zero",
			input:    "0",
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := n.NormalizeString("telefone", tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestEstabelecimentoNormalizer(t *testing.T) {
	n := GetEstabelecimentoNormalizer()

	// Testa normalização de CNPJ
	cnpj := n.NormalizeString("cnpj", "12.345.678/0001-90")
	if !cnpj.Valid || cnpj.String != "12345678000190" {
		t.Errorf("CNPJ: esperado '12345678000190', obtido %v", cnpj)
	}

	// Testa normalização de UF
	uf := n.NormalizeString("uf", "sp")
	if !uf.Valid || uf.String != "SP" {
		t.Errorf("UF: esperado 'SP', obtido %v", uf)
	}

	// Testa normalização de CEP
	cep := n.NormalizeString("cep", "01310-100")
	if !cep.Valid || cep.String != "01310100" {
		t.Errorf("CEP: esperado '01310100', obtido %v", cep)
	}

	// Testa normalização de data
	data := n.NormalizeNullString("data_situacao_cadastral", sql.NullString{String: "20231015", Valid: true})
	if !data.Valid || data.String != "2023-10-15" {
		t.Errorf("Data: esperado '2023-10-15', obtido %v", data)
	}

	// Testa normalização de data inválida
	dataInvalida := n.NormalizeNullString("data_situacao_cadastral", sql.NullString{String: "0", Valid: true})
	if dataInvalida.Valid {
		t.Errorf("Data inválida deveria ser NULL, obtido %v", dataInvalida)
	}
}
