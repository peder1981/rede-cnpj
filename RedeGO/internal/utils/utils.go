package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAcentos remove acentos de uma string
func RemoveAcentos(s string) string {
	if s == "" {
		return ""
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SecureFilename sanitiza um nome de arquivo
func SecureFilename(filename string) string {
	// Remove caracteres perigosos
	filename = filepath.Base(filename)
	
	// Remove caracteres não alfanuméricos (exceto . - _)
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = reg.ReplaceAllString(filename, "_")
	
	// Remove múltiplos underscores consecutivos
	reg = regexp.MustCompile(`_+`)
	filename = reg.ReplaceAllString(filename, "_")
	
	// Remove underscores no início e fim
	filename = strings.Trim(filename, "_")
	
	return filename
}

// GenerateToken gera um token hexadecimal aleatório
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// NomeArquivoNovo gera um nome de arquivo único
func NomeArquivoNovo(caminho string) string {
	if _, err := os.Stat(caminho); os.IsNotExist(err) {
		return caminho
	}

	ext := filepath.Ext(caminho)
	base := strings.TrimSuffix(caminho, ext)

	for i := 1; i <= 100; i++ {
		novoCaminho := fmt.Sprintf("%s%04d%s", base, i, ext)
		if _, err := os.Stat(novoCaminho); os.IsNotExist(err) {
			return novoCaminho
		}
	}

	return caminho
}

// FileExists verifica se um arquivo existe
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDir garante que um diretório existe
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsValidExtension verifica se a extensão do arquivo é permitida
func IsValidExtension(filename string, allowedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedExts {
		if ext == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}

// TruncateString trunca uma string para um tamanho máximo
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// NormalizeSpaces normaliza espaços em branco em uma string
func NormalizeSpaces(s string) string {
	// Remove espaços extras
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// SplitByDelimiter divide uma string por delimitador, ignorando vazios
func SplitByDelimiter(s string, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	
	return result
}

// Contains verifica se um slice contém um elemento
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Unique remove duplicatas de um slice
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)
	
	for _, item := range slice {
		if _, exists := keys[item]; !exists {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// FormatFileSize formata tamanho de arquivo em formato legível
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// IsNumeric verifica se uma string contém apenas dígitos
func IsNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return len(s) > 0
}

// ExtractDigits extrai apenas dígitos de uma string
func ExtractDigits(s string) string {
	var result strings.Builder
	for _, c := range s {
		if unicode.IsDigit(c) {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// PadLeft adiciona padding à esquerda de uma string
func PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(pad, length-len(s)) + s
}

// PadRight adiciona padding à direita de uma string
func PadRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(pad, length-len(s))
}
