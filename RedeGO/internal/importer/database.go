package importer

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
)

// DatabaseType representa o tipo de banco de dados
type DatabaseType int

const (
	DatabaseTypeSQLite DatabaseType = iota
	DatabaseTypePostgreSQL
)

// DatabaseManager gerencia a conexão com o banco de dados
type DatabaseManager struct {
	db       *sql.DB
	dbType   DatabaseType
	connStr  string
}

// NewDatabaseManager cria um novo gerenciador de banco de dados
func NewDatabaseManager(cfg *config.Config, dbPath string) (*DatabaseManager, error) {
	dm := &DatabaseManager{}
	
	// Verifica se PostgreSQL está configurado
	if cfg.PostgresURL != "" {
		dm.dbType = DatabaseTypePostgreSQL
		dm.connStr = cfg.PostgresURL
		
		db, err := sql.Open("postgres", dm.connStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao conectar PostgreSQL: %w", err)
		}
		
		// Testa conexão
		if err := db.Ping(); err != nil {
			db.Close()
			return nil, fmt.Errorf("erro ao pingar PostgreSQL: %w", err)
		}
		
		// Configurações otimizadas para importação em massa
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		
		dm.db = db
		fmt.Println("✅ Usando PostgreSQL para importação direta")
		
	} else {
		// Fallback para SQLite
		dm.dbType = DatabaseTypeSQLite
		dm.connStr = dbPath
		
		// Remove banco antigo se existir
		if _, err := os.Stat(dbPath); err == nil {
			fmt.Printf("⚠️  Removendo banco SQLite antigo: %s\n", dbPath)
			if err := os.Remove(dbPath); err != nil {
				return nil, fmt.Errorf("erro ao remover banco antigo: %w", err)
			}
		}
		
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao abrir SQLite: %w", err)
		}
		
		// Configurações otimizadas para SQLite
		if _, err := db.Exec(`
			PRAGMA journal_mode = WAL;
			PRAGMA synchronous = NORMAL;
			PRAGMA cache_size = -64000;
			PRAGMA temp_store = MEMORY;
		`); err != nil {
			db.Close()
			return nil, fmt.Errorf("erro ao configurar SQLite: %w", err)
		}
		
		dm.db = db
		fmt.Println("⚠️  Usando SQLite (modo legado)")
	}
	
	return dm, nil
}

// GetDB retorna a conexão com o banco
func (dm *DatabaseManager) GetDB() *sql.DB {
	return dm.db
}

// IsPostgreSQL retorna true se estiver usando PostgreSQL
func (dm *DatabaseManager) IsPostgreSQL() bool {
	return dm.dbType == DatabaseTypePostgreSQL
}

// Close fecha a conexão com o banco
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

// TablePrefix retorna o prefixo do schema para PostgreSQL
func (dm *DatabaseManager) TablePrefix(table string) string {
	if !dm.IsPostgreSQL() {
		return table
	}
	
	// Mapear tabelas para schemas PostgreSQL
	switch table {
	case "empresas", "estabelecimento", "socios", "simples",
		"cnae", "motivo", "municipio", "natureza_juridica", "pais", "qualificacao_socio":
		return "receita." + table
	case "ligacao":
		return "rede." + table
	default:
		return table
	}
}

// AdaptPlaceholder adapta placeholders para o banco correto
// SQLite usa ?, PostgreSQL usa $1, $2, etc
func (dm *DatabaseManager) AdaptPlaceholder(query string) string {
	if !dm.IsPostgreSQL() {
		return query
	}
	
	// Converte ? para $1, $2, etc
	result := ""
	paramCount := 1
	for _, char := range query {
		if char == '?' {
			result += fmt.Sprintf("$%d", paramCount)
			paramCount++
		} else {
			result += string(char)
		}
	}
	
	return result
}

// CreateSchemas cria os schemas necessários no PostgreSQL
func (dm *DatabaseManager) CreateSchemas() error {
	if !dm.IsPostgreSQL() {
		return nil // SQLite não precisa de schemas
	}
	
	schemas := []string{
		"CREATE SCHEMA IF NOT EXISTS receita",
		"CREATE SCHEMA IF NOT EXISTS rede",
	}
	
	for _, schema := range schemas {
		if _, err := dm.db.Exec(schema); err != nil {
			return fmt.Errorf("erro ao criar schema: %w", err)
		}
	}
	
	fmt.Println("✅ Schemas PostgreSQL criados/verificados")
	return nil
}
