package importer

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Processor processa os arquivos ZIP e cria o banco cnpj.db
type Processor struct {
	zipDir string
	csvDir string
	dbDir  string
}

// NewProcessor cria um novo processor
func NewProcessor(zipDir, csvDir, dbDir string) *Processor {
	return &Processor{
		zipDir: zipDir,
		csvDir: csvDir,
		dbDir:  dbDir,
	}
}

// Process executa todo o processamento
func (p *Processor) Process() error {
	fmt.Println("📦 Iniciando processamento dos arquivos...")
	
	// Cria diretórios
	if err := os.MkdirAll(p.csvDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(p.dbDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(p.dbDir, "cnpj.db")
	
	// Remove banco antigo se existir
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("⚠️  Removendo banco antigo: %s\n", dbPath)
		if err := os.Remove(dbPath); err != nil {
			return err
		}
	}

	// Abre conexão com banco
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Configura SQLite para performance
	if _, err := db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA cache_size = -64000;
		PRAGMA temp_store = MEMORY;
	`); err != nil {
		return err
	}

	// Cria tabelas
	if err := p.createTables(db); err != nil {
		return err
	}

	// Descompacta e processa arquivos
	zipFiles, err := filepath.Glob(filepath.Join(p.zipDir, "*.zip"))
	if err != nil {
		return err
	}

	fmt.Printf("📋 Encontrados %d arquivos ZIP para processar\n\n", len(zipFiles))

	for i, zipFile := range zipFiles {
		fmt.Printf("[%d/%d] 📦 Processando %s...\n", i+1, len(zipFiles), filepath.Base(zipFile))
		
		if err := p.processZipFile(db, zipFile); err != nil {
			return fmt.Errorf("erro ao processar %s: %w", zipFile, err)
		}
	}

	// Cria índices
	fmt.Println("\n🔍 Criando índices...")
	if err := p.createIndexes(db); err != nil {
		return err
	}

	// Estatísticas
	if err := p.printStats(db); err != nil {
		return err
	}

	fmt.Println("\n✅ Processamento concluído!")
	return nil
}

// createTables cria as tabelas no banco
func (p *Processor) createTables(db *sql.DB) error {
	tables := map[string]string{
		"cnae": `CREATE TABLE cnae (codigo TEXT, descricao TEXT)`,
		"motivo": `CREATE TABLE motivo (codigo TEXT, descricao TEXT)`,
		"municipio": `CREATE TABLE municipio (codigo TEXT, descricao TEXT)`,
		"natureza_juridica": `CREATE TABLE natureza_juridica (codigo TEXT, descricao TEXT)`,
		"pais": `CREATE TABLE pais (codigo TEXT, descricao TEXT)`,
		"qualificacao_socio": `CREATE TABLE qualificacao_socio (codigo TEXT, descricao TEXT)`,
		
		"empresas": `CREATE TABLE empresas (
			cnpj_basico TEXT,
			razao_social TEXT,
			natureza_juridica TEXT,
			qualificacao_responsavel TEXT,
			capital_social REAL,
			porte_empresa TEXT,
			ente_federativo_responsavel TEXT
		)`,
		
		"estabelecimento": `CREATE TABLE estabelecimento (
			cnpj_basico TEXT,
			cnpj_ordem TEXT,
			cnpj_dv TEXT,
			matriz_filial TEXT,
			nome_fantasia TEXT,
			situacao_cadastral TEXT,
			data_situacao_cadastral TEXT,
			motivo_situacao_cadastral TEXT,
			nome_cidade_exterior TEXT,
			pais TEXT,
			data_inicio_atividades TEXT,
			cnae_fiscal TEXT,
			cnae_fiscal_secundaria TEXT,
			tipo_logradouro TEXT,
			logradouro TEXT,
			numero TEXT,
			complemento TEXT,
			bairro TEXT,
			cep TEXT,
			uf TEXT,
			municipio TEXT,
			ddd1 TEXT,
			telefone1 TEXT,
			ddd2 TEXT,
			telefone2 TEXT,
			ddd_fax TEXT,
			fax TEXT,
			correio_eletronico TEXT,
			situacao_especial TEXT,
			data_situacao_especial TEXT,
			cnpj TEXT
		)`,
		
		"socios": `CREATE TABLE socios (
			cnpj TEXT,
			cnpj_basico TEXT,
			identificador_de_socio TEXT,
			nome_socio TEXT,
			cnpj_cpf_socio TEXT,
			qualificacao_socio TEXT,
			data_entrada_sociedade TEXT,
			pais TEXT,
			representante_legal TEXT,
			nome_representante TEXT,
			qualificacao_representante_legal TEXT,
			faixa_etaria TEXT
		)`,
		
		"simples": `CREATE TABLE simples (
			cnpj_basico TEXT,
			opcao_simples TEXT,
			data_opcao_simples TEXT,
			data_exclusao_simples TEXT,
			opcao_mei TEXT,
			data_opcao_mei TEXT,
			data_exclusao_mei TEXT
		)`,
	}

	for name, sql := range tables {
		fmt.Printf("  Criando tabela %s...\n", name)
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("erro ao criar tabela %s: %w", name, err)
		}
	}

	return nil
}

// processZipFile processa um arquivo ZIP
func (p *Processor) processZipFile(db *sql.DB, zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := p.processCSVFromZip(db, f); err != nil {
			return err
		}
	}

	return nil
}

// processCSVFromZip processa um CSV dentro do ZIP
func (p *Processor) processCSVFromZip(db *sql.DB, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Determina tipo de arquivo
	tableName := p.getTableName(f.Name)
	if tableName == "" {
		return nil // Ignora arquivos desconhecidos
	}

	fmt.Printf("    Importando %s para tabela %s...\n", f.Name, tableName)

	// Lê CSV
	reader := csv.NewReader(rc)
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Inicia transação
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepara statement
	stmt, err := p.prepareInsert(tx, tableName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Importa linhas
	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Ignora linhas com erro
		}

		if err := p.insertRecord(stmt, tableName, record); err != nil {
			continue // Ignora registros com erro
		}

		count++
		if count%100000 == 0 {
			fmt.Printf("      %d registros...\n", count)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("      ✅ %d registros importados\n", count)
	return nil
}

// getTableName determina o nome da tabela baseado no nome do arquivo
func (p *Processor) getTableName(filename string) string {
	upper := strings.ToUpper(filename)
	
	if strings.Contains(upper, "CNAE") {
		return "cnae"
	} else if strings.Contains(upper, "MOTI") {
		return "motivo"
	} else if strings.Contains(upper, "MUNIC") {
		return "municipio"
	} else if strings.Contains(upper, "NATJU") {
		return "natureza_juridica"
	} else if strings.Contains(upper, "PAIS") {
		return "pais"
	} else if strings.Contains(upper, "QUALS") {
		return "qualificacao_socio"
	} else if strings.Contains(upper, "EMPRE") {
		return "empresas"
	} else if strings.Contains(upper, "ESTABELE") {
		return "estabelecimento"
	} else if strings.Contains(upper, "SOCIO") {
		return "socios"
	} else if strings.Contains(upper, "SIMPLES") {
		return "simples"
	}
	
	return ""
}

// prepareInsert prepara o statement de insert
func (p *Processor) prepareInsert(tx *sql.Tx, tableName string) (*sql.Stmt, error) {
	queries := map[string]string{
		"cnae":                 "INSERT INTO cnae VALUES (?, ?)",
		"motivo":               "INSERT INTO motivo VALUES (?, ?)",
		"municipio":            "INSERT INTO municipio VALUES (?, ?)",
		"natureza_juridica":    "INSERT INTO natureza_juridica VALUES (?, ?)",
		"pais":                 "INSERT INTO pais VALUES (?, ?)",
		"qualificacao_socio":   "INSERT INTO qualificacao_socio VALUES (?, ?)",
		"empresas":             "INSERT INTO empresas VALUES (?, ?, ?, ?, ?, ?, ?)",
		"estabelecimento":      "INSERT INTO estabelecimento VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"socios":               "INSERT INTO socios VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"simples":              "INSERT INTO simples VALUES (?, ?, ?, ?, ?, ?, ?)",
	}

	query, ok := queries[tableName]
	if !ok {
		return nil, fmt.Errorf("tabela desconhecida: %s", tableName)
	}

	return tx.Prepare(query)
}

// insertRecord insere um registro
func (p *Processor) insertRecord(stmt *sql.Stmt, tableName string, record []string) error {
	// Converte para interface{}
	values := make([]interface{}, len(record))
	for i, v := range record {
		values[i] = v
	}

	// Tratamento especial para estabelecimento (adiciona CNPJ completo)
	if tableName == "estabelecimento" && len(record) >= 3 {
		cnpj := record[0] + record[1] + record[2]
		values = append(values, cnpj)
	}

	// Tratamento especial para socios (adiciona CNPJ da matriz)
	if tableName == "socios" && len(record) >= 1 {
		// Precisamos buscar o CNPJ da matriz, mas por ora vamos deixar vazio
		// Será preenchido depois com UPDATE
		values = append([]interface{}{""}, values...)
	}

	// Tratamento especial para empresas (converte capital social)
	if tableName == "empresas" && len(record) >= 5 {
		capitalStr := strings.ReplaceAll(record[4], ",", ".")
		values[4] = capitalStr
	}

	_, err := stmt.Exec(values...)
	return err
}

// createIndexes cria os índices
func (p *Processor) createIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX idx_empresas_cnpj_basico ON empresas(cnpj_basico)",
		"CREATE INDEX idx_empresas_razao_social ON empresas(razao_social)",
		"CREATE INDEX idx_estabelecimento_cnpj_basico ON estabelecimento(cnpj_basico)",
		"CREATE INDEX idx_estabelecimento_cnpj ON estabelecimento(cnpj)",
		"CREATE INDEX idx_estabelecimento_nomefantasia ON estabelecimento(nome_fantasia)",
		"CREATE INDEX idx_socios_cnpj_basico ON socios(cnpj_basico)",
		"CREATE INDEX idx_socios_cnpj ON socios(cnpj)",
		"CREATE INDEX idx_socios_cnpj_cpf_socio ON socios(cnpj_cpf_socio)",
		"CREATE INDEX idx_socios_nome_socio ON socios(nome_socio)",
		"CREATE INDEX idx_simples_cnpj_basico ON simples(cnpj_basico)",
	}

	for _, idx := range indexes {
		fmt.Printf("  %s\n", idx)
		if _, err := db.Exec(idx); err != nil {
			return err
		}
	}

	// Atualiza CNPJ dos sócios
	fmt.Println("  Atualizando CNPJ dos sócios...")
	_, err := db.Exec(`
		UPDATE socios
		SET cnpj = (
			SELECT cnpj FROM estabelecimento 
			WHERE estabelecimento.cnpj_basico = socios.cnpj_basico 
			AND estabelecimento.matriz_filial = '1'
			LIMIT 1
		)
	`)

	return err
}

// printStats imprime estatísticas
func (p *Processor) printStats(db *sql.DB) error {
	fmt.Println("\n📊 Estatísticas:")
	
	tables := []string{"empresas", "estabelecimento", "socios", "simples"}
	for _, table := range tables {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			return err
		}
		fmt.Printf("  %s: %d registros\n", table, count)
	}

	// Data de referência
	var dataRef string
	err := db.QueryRow("SELECT data_situacao_cadastral FROM estabelecimento LIMIT 1").Scan(&dataRef)
	if err == nil && len(dataRef) >= 8 {
		fmt.Printf("\n📅 Data de referência: %s/%s/%s\n", dataRef[6:8], dataRef[4:6], dataRef[0:4])
	}

	return nil
}
