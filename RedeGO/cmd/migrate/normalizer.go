package main

import (
	"database/sql"
	"regexp"
	"strings"
)

// FieldType representa o tipo de um campo no banco de dados
type FieldType string

const (
	FieldTypeDate        FieldType = "DATE"
	FieldTypeVarchar     FieldType = "VARCHAR"
	FieldTypeText        FieldType = "TEXT"
	FieldTypeNumeric     FieldType = "NUMERIC"
	FieldTypeInteger     FieldType = "INTEGER"
	FieldTypeCNPJ        FieldType = "CNPJ"
	FieldTypeCPF         FieldType = "CPF"
	FieldTypeCEP         FieldType = "CEP"
	FieldTypeEmail       FieldType = "EMAIL"
	FieldTypeUF          FieldType = "UF"
	FieldTypeCode        FieldType = "CODE"
	FieldTypePhone       FieldType = "PHONE"
)

// FieldMetadata contém metadados sobre um campo
type FieldMetadata struct {
	Name       string
	Type       FieldType
	MaxLength  int
	Required   bool
	Pattern    *regexp.Regexp
	ValidValues []string
}

// Normalizer é responsável por normalizar dados baseado em metadados
type Normalizer struct {
	metadata map[string]FieldMetadata
}

// NewNormalizer cria um novo normalizador
func NewNormalizer() *Normalizer {
	return &Normalizer{
		metadata: make(map[string]FieldMetadata),
	}
}

// RegisterField registra metadados de um campo
func (n *Normalizer) RegisterField(field FieldMetadata) {
	n.metadata[field.Name] = field
}

// NormalizeString normaliza uma string baseado nos metadados do campo
func (n *Normalizer) NormalizeString(fieldName string, value string) sql.NullString {
	meta, exists := n.metadata[fieldName]
	if !exists {
		// Se não tem metadados, apenas sanitiza UTF-8
		return sql.NullString{
			String: sanitizeString(value),
			Valid:  value != "",
		}
	}

	// Remove espaços em branco
	trimmed := strings.TrimSpace(value)

	// Se campo obrigatório e vazio, retorna inválido
	if meta.Required && trimmed == "" {
		return sql.NullString{Valid: false}
	}

	// Se campo opcional e vazio, retorna NULL
	if !meta.Required && trimmed == "" {
		return sql.NullString{Valid: false}
	}

	// Normaliza baseado no tipo
	switch meta.Type {
	case FieldTypeDate:
		return n.normalizeDate(trimmed)
	case FieldTypeCNPJ:
		return n.normalizeCNPJ(trimmed)
	case FieldTypeCPF:
		return n.normalizeCPF(trimmed)
	case FieldTypeCEP:
		return n.normalizeCEP(trimmed)
	case FieldTypeEmail:
		return n.normalizeEmail(trimmed)
	case FieldTypeUF:
		return n.normalizeUF(trimmed)
	case FieldTypeCode:
		return n.normalizeCode(trimmed, meta)
	case FieldTypePhone:
		return n.normalizePhone(trimmed)
	case FieldTypeVarchar, FieldTypeText:
		return n.normalizeText(trimmed, meta)
	default:
		return sql.NullString{
			String: sanitizeString(trimmed),
			Valid:  true,
		}
	}
}

// NormalizeNullString normaliza um sql.NullString
func (n *Normalizer) NormalizeNullString(fieldName string, ns sql.NullString) sql.NullString {
	if !ns.Valid {
		return sql.NullString{Valid: false}
	}
	return n.NormalizeString(fieldName, ns.String)
}

// normalizeDate normaliza campos de data
func (n *Normalizer) normalizeDate(value string) sql.NullString {
	// Valores inválidos
	invalidValues := []string{"", "0", "00", "000", "0000", "00000000"}
	for _, invalid := range invalidValues {
		if value == invalid {
			return sql.NullString{Valid: false}
		}
	}

	// Formato YYYYMMDD (8 dígitos)
	if len(value) == 8 {
		// Verifica se todos são dígitos
		for _, c := range value {
			if c < '0' || c > '9' {
				return sql.NullString{Valid: false}
			}
		}

		year := value[0:4]
		month := value[4:6]
		day := value[6:8]

		// Valida componentes
		if year < "1900" || year > "2100" {
			return sql.NullString{Valid: false}
		}
		if month < "01" || month > "12" {
			return sql.NullString{Valid: false}
		}
		if day < "01" || day > "31" {
			return sql.NullString{Valid: false}
		}

		// Retorna formatado como YYYY-MM-DD
		return sql.NullString{
			String: year + "-" + month + "-" + day,
			Valid:  true,
		}
	}

	// Formato YYYY-MM-DD (já correto)
	if len(value) == 10 && value[4] == '-' && value[7] == '-' {
		return sql.NullString{String: value, Valid: true}
	}

	// Formato inválido
	return sql.NullString{Valid: false}
}

// normalizeCNPJ normaliza CNPJ (remove máscara, valida tamanho)
func (n *Normalizer) normalizeCNPJ(value string) sql.NullString {
	// Remove caracteres não numéricos
	cnpj := regexp.MustCompile(`[^0-9]`).ReplaceAllString(value, "")

	// CNPJ deve ter 14 dígitos
	if len(cnpj) != 14 {
		return sql.NullString{Valid: false}
	}

	// CNPJ não pode ser todo zeros
	if cnpj == "00000000000000" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: cnpj, Valid: true}
}

// normalizeCPF normaliza CPF (remove máscara, valida tamanho)
func (n *Normalizer) normalizeCPF(value string) sql.NullString {
	// Remove caracteres não numéricos
	cpf := regexp.MustCompile(`[^0-9]`).ReplaceAllString(value, "")

	// CPF deve ter 11 dígitos
	if len(cpf) != 11 {
		return sql.NullString{Valid: false}
	}

	// CPF não pode ser todo zeros
	if cpf == "00000000000" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: cpf, Valid: true}
}

// normalizeCEP normaliza CEP (remove máscara, valida tamanho)
func (n *Normalizer) normalizeCEP(value string) sql.NullString {
	// Remove caracteres não numéricos
	cep := regexp.MustCompile(`[^0-9]`).ReplaceAllString(value, "")

	// CEP deve ter 8 dígitos
	if len(cep) != 8 {
		return sql.NullString{Valid: false}
	}

	// CEP não pode ser todo zeros
	if cep == "00000000" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: cep, Valid: true}
}

// normalizeEmail normaliza email
func (n *Normalizer) normalizeEmail(value string) sql.NullString {
	// Email vazio é válido (campo opcional)
	if value == "" {
		return sql.NullString{Valid: false}
	}

	// Converte para minúsculas
	email := strings.ToLower(value)

	// Validação básica de email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: sanitizeString(email), Valid: true}
}

// normalizeUF normaliza UF (sigla de estado)
func (n *Normalizer) normalizeUF(value string) sql.NullString {
	uf := strings.ToUpper(strings.TrimSpace(value))

	// Lista de UFs válidas
	validUFs := map[string]bool{
		"AC": true, "AL": true, "AP": true, "AM": true, "BA": true,
		"CE": true, "DF": true, "ES": true, "GO": true, "MA": true,
		"MT": true, "MS": true, "MG": true, "PA": true, "PB": true,
		"PR": true, "PE": true, "PI": true, "RJ": true, "RN": true,
		"RS": true, "RO": true, "RR": true, "SC": true, "SP": true,
		"SE": true, "TO": true, "EX": true, // EX = Exterior
	}

	if !validUFs[uf] {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: uf, Valid: true}
}

// normalizeCode normaliza códigos (CNAE, natureza jurídica, etc)
func (n *Normalizer) normalizeCode(value string, meta FieldMetadata) sql.NullString {
	// Remove espaços
	code := strings.TrimSpace(value)

	// Código vazio ou "0" é inválido
	if code == "" || code == "0" || code == "00" {
		return sql.NullString{Valid: false}
	}

	// Valida tamanho máximo
	if meta.MaxLength > 0 && len(code) > meta.MaxLength {
		code = code[:meta.MaxLength]
	}

	// Valida padrão se especificado
	if meta.Pattern != nil && !meta.Pattern.MatchString(code) {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: code, Valid: true}
}

// normalizePhone normaliza telefone (remove caracteres não numéricos)
func (n *Normalizer) normalizePhone(value string) sql.NullString {
	// Remove caracteres não numéricos
	phone := regexp.MustCompile(`[^0-9]`).ReplaceAllString(value, "")

	// Telefone vazio é válido (campo opcional)
	if phone == "" || phone == "0" {
		return sql.NullString{Valid: false}
	}

	// Telefone deve ter entre 8 e 11 dígitos
	if len(phone) < 8 || len(phone) > 11 {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: phone, Valid: true}
}

// normalizeText normaliza campos de texto
func (n *Normalizer) normalizeText(value string, meta FieldMetadata) sql.NullString {
	// Sanitiza UTF-8
	text := sanitizeString(value)

	// Valida tamanho máximo
	if meta.MaxLength > 0 && len(text) > meta.MaxLength {
		text = text[:meta.MaxLength]
	}

	// Valida padrão se especificado
	if meta.Pattern != nil && !meta.Pattern.MatchString(text) {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: text, Valid: true}
}

// NormalizeFloat64 normaliza valores numéricos
func (n *Normalizer) NormalizeFloat64(fieldName string, value sql.NullFloat64) sql.NullFloat64 {
	if !value.Valid {
		return sql.NullFloat64{Valid: false}
	}

	// Valores negativos em campos que não permitem
	meta, exists := n.metadata[fieldName]
	if exists && meta.Type == FieldTypeNumeric {
		// Capital social não pode ser negativo
		if value.Float64 < 0 {
			return sql.NullFloat64{Valid: false}
		}
	}

	return value
}
