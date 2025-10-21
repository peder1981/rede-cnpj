package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
)

var (
	dbReceita *sql.DB
	dbRede    *sql.DB
	dbSearch  *sql.DB
	dbLocal   *sql.DB
	dbMutex   sync.Mutex
	once      sync.Once
)

// InitDatabases inicializa as conexões com os bancos de dados
func InitDatabases(cfg *config.Config) error {
	var err error

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

// GetDBReceita retorna a conexão com o banco da Receita
func GetDBReceita() *sql.DB {
	return dbReceita
}

// GetDBRede retorna a conexão com o banco de Rede
func GetDBRede() *sql.DB {
	return dbRede
}

// GetDBSearch retorna a conexão com o banco de Search
func GetDBSearch() *sql.DB {
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
	rows, err := db.Query("SELECT codigo, descricao FROM qualificacao_socio")
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
	rows, err = db.Query("SELECT codigo, descricao FROM motivo")
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
	rows, err = db.Query("SELECT codigo, descricao FROM cnae")
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
	rows, err = db.Query("SELECT codigo, descricao FROM natureza_juridica")
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
