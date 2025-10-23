package search

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// SearchOptions opções de busca avançada
type SearchOptions struct {
	Query      string
	Limit      int
	UseGlob    bool  // Usar wildcards (* ?)
	MatchAll   bool  // Match todas as palavras
	RandomTest bool  // Busca aleatória para teste
}

// AdvancedSearch busca avançada com FTS5
type AdvancedSearch struct {
	dbPath string
}

// NewAdvancedSearch cria um novo buscador avançado
func NewAdvancedSearch(dbPath string) *AdvancedSearch {
	return &AdvancedSearch{dbPath: dbPath}
}

// Search executa busca avançada
func (s *AdvancedSearch) Search(opts SearchOptions) ([]string, error) {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Busca aleatória para teste
	if opts.RandomTest {
		return s.searchRandom(db)
	}

	// Busca com wildcards
	if opts.UseGlob && (strings.Contains(opts.Query, "*") || strings.Contains(opts.Query, "?")) {
		return s.searchWithGlob(db, opts)
	}

	// Busca normal (match flexível)
	return s.searchNormal(db, opts)
}

// searchRandom busca aleatória
func (s *AdvancedSearch) searchRandom(db *sql.DB) ([]string, error) {
	query := `
		SELECT id_descricao 
		FROM id_search 
		WHERE rowid > (abs(random()) % (SELECT max(rowid) FROM id_search)) 
		LIMIT 1
	`

	var result string
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		return nil, err
	}

	return []string{result}, nil
}

// searchWithGlob busca com wildcards
func (s *AdvancedSearch) searchWithGlob(db *sql.DB, opts SearchOptions) ([]string, error) {
	// Remove caracteres especiais exceto * e ?
	cleaned := cleanForMatch(opts.Query)
	
	// Prepara para MATCH (substitui * por espaço, ? por *)
	matchQuery := strings.ReplaceAll(cleaned, "*", " ")
	matchQuery = strings.ReplaceAll(matchQuery, "?", "*")
	matchQuery = strings.TrimLeft(matchQuery, "*")
	matchQuery = strings.TrimSpace(matchQuery)

	// Prepara para GLOB
	globQuery := "*-" + strings.ReplaceAll(cleaned, " ", "?")

	if matchQuery == "" {
		return nil, fmt.Errorf("query vazia após limpeza")
	}

	query := `
		SELECT DISTINCT id_descricao 
		FROM id_search
		WHERE id_descricao MATCH ? 
		  AND id_descricao GLOB ?
		LIMIT ?
	`

	rows, err := db.Query(query, matchQuery, globQuery, opts.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanResults(rows)
}

// searchNormal busca normal (match flexível)
func (s *AdvancedSearch) searchNormal(db *sql.DB, opts SearchOptions) ([]string, error) {
	cleaned := cleanForMatch(opts.Query)
	if cleaned == "" {
		return nil, fmt.Errorf("query vazia")
	}

	query := `
		SELECT DISTINCT id_descricao 
		FROM id_search
		WHERE id_descricao MATCH ?
		LIMIT ?
	`

	rows, err := db.Query(query, cleaned, opts.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanResults(rows)
}

// cleanForMatch limpa string para busca
func cleanForMatch(s string) string {
	s = strings.ToUpper(s)
	
	// Remove acentos (simplificado)
	replacements := map[string]string{
		"Á": "A", "À": "A", "Ã": "A", "Â": "A",
		"É": "E", "Ê": "E",
		"Í": "I",
		"Ó": "O", "Õ": "O", "Ô": "O",
		"Ú": "U",
		"Ç": "C",
	}
	
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	// Remove caracteres especiais (exceto * e ?)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' || r == '*' || r == '?' {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}

	return strings.TrimSpace(result.String())
}

// scanResults lê resultados da query
func scanResults(rows *sql.Rows) ([]string, error) {
	var results []string
	
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		results = append(results, id)
	}

	return results, rows.Err()
}
