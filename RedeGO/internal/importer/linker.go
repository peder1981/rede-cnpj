package importer

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Linker cria as tabelas de ligação (rede.db)
type Linker struct {
	dbDir string
}

// NewLinker cria um novo linker
func NewLinker(dbDir string) *Linker {
	return &Linker{dbDir: dbDir}
}

// CreateLinks cria as tabelas de ligação
func (l *Linker) CreateLinks() error {
	fmt.Println("🔗 Criando tabelas de ligação...")
	
	cnpjDB := filepath.Join(l.dbDir, "cnpj.db")
	redeDB := filepath.Join(l.dbDir, "rede.db")

	// Abre conexão
	db, err := sql.Open("sqlite3", redeDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Anexa banco CNPJ
	_, err = db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' as cnpj", cnpjDB))
	if err != nil {
		return err
	}

	// SQL para criar tabelas de ligação
	sqlLigacao := `
-- Remove tabelas antigas
DROP TABLE IF EXISTS ligacao;
DROP TABLE IF EXISTS ligacao1;

-- PJ->PJ vínculo sócio pessoa jurídica
CREATE TABLE ligacao1 AS
SELECT 'PJ_'||t.cnpj_cpf_socio as origem, 
       'PJ_'||t.cnpj as destino, 
       sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_socio
WHERE length(t.cnpj_cpf_socio)=14;

-- PF->PJ vínculo de sócio pessoa física
INSERT INTO ligacao1
SELECT 'PF_'||t.cnpj_cpf_socio||'-'||t.nome_socio as origem, 
       'PJ_'||t.cnpj as destino, 
       sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_socio
WHERE length(t.cnpj_cpf_socio)=11 AND t.nome_socio<>'';

-- PE->PJ empresa sócia no exterior
INSERT INTO ligacao1
SELECT 'PE_'||t.nome_socio as origem, 
       'PJ_'||t.cnpj as destino, 
       sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_socio
WHERE length(t.cnpj_cpf_socio)<>14 
  AND length(t.cnpj_cpf_socio)<>11 
  AND t.cnpj_cpf_socio='';

-- PF->PE representante legal de empresa sócia no exterior
INSERT INTO ligacao1
SELECT 'PF_'||t.representante_legal||'-'||t.nome_representante as origem, 
       'PE_'||t.nome_socio as destino, 
       'rep-sócio-'||sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_representante_legal
WHERE length(t.cnpj_cpf_socio)<>14 
  AND length(t.cnpj_cpf_socio)<>11 
  AND t.cnpj_cpf_socio='' 
  AND t.representante_legal<>'***000000**';

-- PF->PJ representante legal PJ->PJ
INSERT INTO ligacao1
SELECT 'PF_'||t.representante_legal||'-'||t.nome_representante as origem, 
       'PJ_'||t.cnpj_cpf_socio as destino, 
       'rep-sócio-'||sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_representante_legal
WHERE length(t.cnpj_cpf_socio)=14 
  AND t.representante_legal<>'***000000**';

-- PF->PF representante legal de sócio PF
INSERT INTO ligacao1
SELECT 'PF_'||t.representante_legal||'-'||t.nome_representante as origem, 
       'PF_'||t.cnpj_cpf_socio||'-'||t.nome_socio as destino, 
       'rep-sócio-'||sq.descricao as tipo, 
       'socios' as base
FROM cnpj.socios t
LEFT JOIN cnpj.qualificacao_socio sq ON sq.codigo=t.qualificacao_representante_legal
WHERE length(t.cnpj_cpf_socio)=11 
  AND t.representante_legal<>'***000000**';

-- Tabela temporária de filiais
CREATE TABLE tfilial AS 
SELECT cnpj, cnpj_basico
FROM cnpj.estabelecimento t
WHERE t.matriz_filial = '2';

CREATE INDEX idx_filiais ON tfilial(cnpj_basico);

-- PJ filial -> PJ matriz
INSERT INTO ligacao1
SELECT 'PJ_'||tf.cnpj as origem, 
       'PJ_'||t.cnpj as destino, 
       'filial' as tipo, 
       'estabelecimento' as base
FROM tfilial tf
LEFT JOIN cnpj.estabelecimento t ON t.cnpj_basico=tf.cnpj_basico 
WHERE t.matriz_filial = '1';

DROP TABLE IF EXISTS tfilial;

-- Cria tabela final de ligação (remove duplicatas)
CREATE TABLE ligacao AS
SELECT origem as id1, 
       destino as id2, 
       tipo as descricao, 
       base as comentario 
FROM ligacao1 
GROUP BY origem, destino, tipo, base;

DROP TABLE IF EXISTS ligacao1;

-- Cria índices
CREATE INDEX idx_ligacao_origem ON ligacao(id1);
CREATE INDEX idx_ligacao_destino ON ligacao(id2);
`

	// Executa SQL
	fmt.Println("  Executando SQL de ligação...")
	start := time.Now()
	
	if _, err := db.Exec(sqlLigacao); err != nil {
		return fmt.Errorf("erro ao executar SQL: %w", err)
	}

	// Estatísticas
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM ligacao").Scan(&count); err != nil {
		return err
	}

	fmt.Printf("  ✅ %d ligações criadas em %v\n", count, time.Since(start))
	return nil
}
