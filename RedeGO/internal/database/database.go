package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
)

var (
	db        *sql.DB // PostgreSQL connection
	dbReceita *sql.DB // Legacy SQLite (deprecated)
	dbRede    *sql.DB // Legacy SQLite (deprecated)
	dbSearch  *sql.DB // Legacy SQLite (deprecated)
	dbLocal   *sql.DB
	dbMutex   sync.Mutex
	once      sync.Once
	usePostgres bool
)

// InitDatabases inicializa as conexões com os bancos de dados
func InitDatabases(cfg *config.Config) error {
	var err error

	// Verificar se deve usar PostgreSQL
	if cfg.PostgresURL != "" {
		usePostgres = true
		db, err = sql.Open("postgres", cfg.PostgresURL)
		if err != nil {
			return fmt.Errorf("erro ao conectar PostgreSQL: %w", err)
		}
		
		// Configurações otimizadas para PostgreSQL
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
		
		// Testar conexão
		if err := db.Ping(); err != nil {
			return fmt.Errorf("erro ao pingar PostgreSQL: %w", err)
		}
		
		fmt.Println("✅ Conectado ao PostgreSQL")
		return nil
	}

	// Fallback para SQLite (legacy)
	fmt.Println("⚠️  Usando SQLite (modo legado)")
	usePostgres = false

	// Base da Receita (CNPJ)
	if cfg.BaseReceita != "" {
		dbReceita, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", cfg.BaseReceita))
		if err != nil {
			return fmt.Errorf("erro ao abrir base receita: %w", err)
		}
		dbReceita.SetMaxOpenConns(10)
		dbReceita.SetMaxIdleConns(5)
	}

	// Base de Rede
	if cfg.BaseRede != "" {
		dbRede, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", cfg.BaseRede))
		if err != nil {
			return fmt.Errorf("erro ao abrir base rede: %w", err)
		}
		dbRede.SetMaxOpenConns(10)
		dbRede.SetMaxIdleConns(5)
	}

	// Base de Search
	if cfg.BaseRedeSearch != "" {
		dbSearch, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", cfg.BaseRedeSearch))
		if err != nil {
			return fmt.Errorf("erro ao abrir base search: %w", err)
		}
		dbSearch.SetMaxOpenConns(10)
		dbSearch.SetMaxIdleConns(5)
	}

	// Base Local (se configurada)
	if cfg.BaseLocal != "" {
		dbLocal, err = sql.Open("sqlite3", cfg.BaseLocal)
		if err != nil {
			return fmt.Errorf("erro ao abrir base local: %w", err)
		}
		dbLocal.SetMaxOpenConns(5)
		dbLocal.SetMaxIdleConns(2)
	}

	return nil
}

// GetDB retorna a conexão principal (PostgreSQL ou SQLite)
func GetDB() *sql.DB {
	if usePostgres {
		return db
	}
	return dbReceita
}

// IsPostgres retorna true se estiver usando PostgreSQL
func IsPostgres() bool {
	return usePostgres
}

// GetDBReceita retorna a conexão com o banco da Receita (legacy)
func GetDBReceita() *sql.DB {
	if usePostgres {
		return db
	}
	return dbReceita
}

// GetDBRede retorna a conexão com o banco de Rede (legacy)
func GetDBRede() *sql.DB {
	if usePostgres {
		return db
	}
	return dbRede
}

// GetDBSearch retorna a conexão com o banco de Search (legacy)
func GetDBSearch() *sql.DB {
	if usePostgres {
		return db
	}
	return dbSearch
}

// GetDBLocal retorna a conexão com o banco Local
func GetDBLocal() *sql.DB {
	return dbLocal
}

// Lock adquire o mutex global do banco de dados
func Lock() {
	dbMutex.Lock()
}

// Unlock libera o mutex global do banco de dados
func Unlock() {
	dbMutex.Unlock()
}

// Close fecha todas as conexões de banco de dados
func Close() {
	if db != nil {
		db.Close()
	}
	if dbReceita != nil {
		dbReceita.Close()
	}
	if dbRede != nil {
		dbRede.Close()
	}
	if dbSearch != nil {
		dbSearch.Close()
	}
	if dbLocal != nil {
		dbLocal.Close()
	}
}

// AdaptQuery adapta uma query SQLite para PostgreSQL
func AdaptQuery(query string) string {
	if !usePostgres {
		return query
	}
	
	// Substituir placeholders ? por $1, $2, etc
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

// TablePrefix retorna o prefixo do schema para PostgreSQL
func TablePrefix(table string) string {
	if !usePostgres {
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

// DicionariosCodigosCNPJ armazena os dicionários de códigos
type DicionariosCodigosCNPJ struct {
	QualificacaoSocio  map[string]string
	MotivoSituacao     map[string]string
	CNAE               map[string]string
	NaturezaJuridica   map[string]string
	SituacaoCadastral  map[string]string
	PorteEmpresa       map[string]string
}

var dicionarios *DicionariosCodigosCNPJ

// LoadDicionarios carrega os dicionários de códigos do banco de dados
func LoadDicionarios() (*DicionariosCodigosCNPJ, error) {
	if dicionarios != nil {
		return dicionarios, nil
	}

	db := GetDBReceita()
	if db == nil {
		return nil, fmt.Errorf("banco de dados da receita não inicializado")
	}

	dic := &DicionariosCodigosCNPJ{
		QualificacaoSocio: make(map[string]string),
		MotivoSituacao:    make(map[string]string),
		CNAE:              make(map[string]string),
		NaturezaJuridica:  make(map[string]string),
		SituacaoCadastral: map[string]string{
			"01": "Nula",
			"02": "Ativa",
			"03": "Suspensa",
			"04": "Inapta",
			"08": "Baixada",
		},
		PorteEmpresa: map[string]string{
			"00": "Não informado",
			"01": "Micro empresa",
			"03": "Empresa de pequeno porte",
			"05": "Demais (Médio ou Grande porte)",
		},
	}

	// Carrega qualificação de sócio
	rows, err := db.Query(fmt.Sprintf("SELECT codigo, descricao FROM %s", TablePrefix("qualificacao_socio")))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var codigo, descricao string
			if err := rows.Scan(&codigo, &descricao); err == nil {
				dic.QualificacaoSocio[codigo] = descricao
			}
		}
	}

	// Carrega motivo de situação
	rows, err = db.Query(fmt.Sprintf("SELECT codigo, descricao FROM %s", TablePrefix("motivo")))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var codigo, descricao string
			if err := rows.Scan(&codigo, &descricao); err == nil {
				dic.MotivoSituacao[codigo] = descricao
			}
		}
	}

	// Carrega CNAE
	rows, err = db.Query(fmt.Sprintf("SELECT codigo, descricao FROM %s", TablePrefix("cnae")))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var codigo, descricao string
			if err := rows.Scan(&codigo, &descricao); err == nil {
				dic.CNAE[codigo] = descricao
			}
		}
	}

	// Carrega natureza jurídica
	rows, err = db.Query(fmt.Sprintf("SELECT codigo, descricao FROM %s", TablePrefix("natureza_juridica")))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var codigo, descricao string
			if err := rows.Scan(&codigo, &descricao); err == nil {
				dic.NaturezaJuridica[codigo] = descricao
			}
		}
	}

	dicionarios = dic
	return dic, nil
}

// GetDicionarios retorna os dicionários carregados
func GetDicionarios() *DicionariosCodigosCNPJ {
	if dicionarios == nil {
		LoadDicionarios()
	}
	return dicionarios
}
