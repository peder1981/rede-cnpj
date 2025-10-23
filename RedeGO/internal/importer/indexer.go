package importer

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Indexer cria os índices de busca (rede_search.db)
type Indexer struct {
	dbDir string
}

// NewIndexer cria um novo indexer
func NewIndexer(dbDir string) *Indexer {
	return &Indexer{dbDir: dbDir}
}

// CreateIndexes cria os índices de busca full-text
func (i *Indexer) CreateIndexes() error {
	fmt.Println("🔍 Criando índices de busca...")
	
	cnpjDB := filepath.Join(i.dbDir, "cnpj.db")
	redeDB := filepath.Join(i.dbDir, "rede.db")
	searchDB := filepath.Join(i.dbDir, "rede_search.db")

	// Abre conexão
	db, err := sql.Open("sqlite3", searchDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Anexa bancos
	_, err = db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as cnpj", cnpjDB))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as rede", redeDB))
	if err != nil {
		return err
	}

	// SQL para criar índices FTS5
	sqlSearch := `
-- Remove tabela antiga
DROP TABLE IF EXISTS id_search;

-- Cria tabela FTS5 (Full-Text Search)
CREATE VIRTUAL TABLE id_search USING fts5(id_descricao);

-- Insere dados para busca
INSERT INTO id_search
SELECT id_descricao
FROM (
    -- PJ com razão social
    SELECT 'PJ_' || te.cnpj || '-' || t.razao_social as id_descricao
    FROM cnpj.estabelecimento te 
    LEFT JOIN cnpj.empresas t ON t.cnpj_basico=te.cnpj_basico
    WHERE te.matriz_filial='1'
    
    UNION ALL
    
    -- PJ com nome fantasia
    SELECT 'PJ_' || te.cnpj || '-' || te.nome_fantasia as id_descricao 
    FROM cnpj.estabelecimento te 
    WHERE te.nome_fantasia IS NOT NULL AND te.nome_fantasia <> ''
    
    UNION ALL
    
    -- PF e PE da tabela de ligação
    SELECT id1 as id_descricao
    FROM rede.ligacao
    WHERE substr(id1,1,3)<>'PJ_'
    
    UNION ALL
    
    SELECT id2 as id_descricao
    FROM rede.ligacao
    WHERE substr(id2,1,3)<>'PJ_'
) as tunion
GROUP BY id_descricao;
`

	// Executa SQL
	fmt.Println("  Executando SQL de indexação...")
	start := time.Now()
	
	if _, err := db.Exec(sqlSearch); err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	// Estatísticas
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM id_search").Scan(&count); err != nil {
		return err
	}

	fmt.Printf("  ✅ %d entradas indexadas em %v\n", count, time.Since(start))
	return nil
}
