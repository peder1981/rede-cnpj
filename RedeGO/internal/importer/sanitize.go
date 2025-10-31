package importer

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// sanitizeString corrige problemas de encoding em strings
func sanitizeString(s string) string {
	// Se já é UTF-8 válido, retorna como está
	if utf8.ValidString(s) {
		return s
	}

	// Tenta converter de ISO-8859-1 (Latin1) para UTF-8
	decoder := charmap.ISO8859_1.NewDecoder()
	result, _, err := transform.String(decoder, s)
	if err == nil && utf8.ValidString(result) {
		return result
	}

	// Se ainda não funcionou, remove caracteres inválidos
	return strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return -1 // Remove o caractere
		}
		return r
	}, s)
}
