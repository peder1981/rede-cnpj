package main

import (
	"database/sql"
	"testing"
)

func TestSanitizeDateString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected sql.NullString
	}{
		{
			name:     "Valor NULL",
			input:    sql.NullString{Valid: false},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "String vazia",
			input:    sql.NullString{String: "", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Valor zero",
			input:    sql.NullString{String: "0", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Valor 00000000",
			input:    sql.NullString{String: "00000000", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Data válida YYYYMMDD",
			input:    sql.NullString{String: "20231015", Valid: true},
			expected: sql.NullString{String: "2023-10-15", Valid: true},
		},
		{
			name:     "Data válida YYYY-MM-DD",
			input:    sql.NullString{String: "2023-10-15", Valid: true},
			expected: sql.NullString{String: "2023-10-15", Valid: true},
		},
		{
			name:     "Data com espaços",
			input:    sql.NullString{String: "  20231015  ", Valid: true},
			expected: sql.NullString{String: "2023-10-15", Valid: true},
		},
		{
			name:     "Ano inválido (muito antigo)",
			input:    sql.NullString{String: "18991231", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Ano inválido (muito futuro)",
			input:    sql.NullString{String: "21010101", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Mês inválido",
			input:    sql.NullString{String: "20231315", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Dia inválido",
			input:    sql.NullString{String: "20231032", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Formato inválido (letras)",
			input:    sql.NullString{String: "2023AB15", Valid: true},
			expected: sql.NullString{Valid: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeDateString(tt.input)
			
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}

func TestSanitizeNumericString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected sql.NullString
	}{
		{
			name:     "Valor NULL",
			input:    sql.NullString{Valid: false},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "String vazia",
			input:    sql.NullString{String: "", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Valor zero",
			input:    sql.NullString{String: "0", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Valor 00",
			input:    sql.NullString{String: "00", Valid: true},
			expected: sql.NullString{Valid: false},
		},
		{
			name:     "Valor válido",
			input:    sql.NullString{String: "123", Valid: true},
			expected: sql.NullString{String: "123", Valid: true},
		},
		{
			name:     "Valor com espaços",
			input:    sql.NullString{String: "  456  ", Valid: true},
			expected: sql.NullString{String: "456", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeNumericString(tt.input)
			
			if result.Valid != tt.expected.Valid {
				t.Errorf("Valid: esperado %v, obtido %v", tt.expected.Valid, result.Valid)
			}
			
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("String: esperado %q, obtido %q", tt.expected.String, result.String)
			}
		})
	}
}
