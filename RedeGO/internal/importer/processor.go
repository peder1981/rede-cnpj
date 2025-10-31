package importer

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
)

// Processor processa os arquivos ZIP e importa para o banco de dados
type Processor struct {
	zipDir string
	csvDir string
	dbDir  string
	dbMgr  *DatabaseManager
	cfg    *config.Config
}

// NewProcessor cria um novo processor
func NewProcessor(zipDir, csvDir, dbDir string) *Processor {
	return &Processor{
		zipDir: zipDir,
		csvDir: csvDir,
		dbDir:  dbDir,
	}
}

// NewProcessorWithConfig cria um novo processor com configuração
func NewProcessorWithConfig(zipDir, csvDir, dbDir string, cfg *config.Config) *Processor {
	return &Processor{
		zipDir: zipDir,
		csvDir: csvDir,
		dbDir:  dbDir,
		cfg:    cfg,
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
	
	// Cria gerenciador de banco de dados
	var err error
	if p.cfg != nil {
		p.dbMgr, err = NewDatabaseManager(p.cfg, dbPath)
	} else {
		// Fallback para SQLite se não houver config
		p.dbMgr, err = NewDatabaseManager(&config.Config{}, dbPath)
	}
	if err != nil {
		return fmt.Errorf("erro ao inicializar banco: %w", err)
	}
	defer p.dbMgr.Close()
	
	// Cria schemas se PostgreSQL
	if err := p.dbMgr.CreateSchemas(); err != nil {
		return fmt.Errorf("erro ao criar schemas: %w", err)
	}

	// Cria tabelas
	if err := p.createTables(); err != nil {
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
		
		if err := p.processZipFile(zipFile); err != nil {
			return fmt.Errorf("erro ao processar %s: %w", zipFile, err)
		}
	}

	// Cria índices
	fmt.Println("\n🔍 Criando índices...")
	if err := p.createIndexes(); err != nil {
		return err
	}

	// Estatísticas
	if err := p.printStats(); err != nil {
		return err
	}

	fmt.Println("\n✅ Processamento concluído!")
	return nil
}

// createTables cria as tabelas no banco
func (p *Processor) createTables() error {
	db := p.dbMgr.GetDB()
	
	var tables map[string]string
	if p.dbMgr.IsPostgreSQL() {
		tables = GetTableSchemasPostgreSQL()
	} else {
		tables = GetTableSchemasSQLite()
	}
	
	for name, sqlStmt := range tables {
		fmt.Printf("  Criando tabela %s...\n", name)
		if _, err := db.Exec(sqlStmt); err != nil {
			return fmt.Errorf("erro ao criar tabela %s: %w", name, err)
		}
	}
	
	return nil
}

// processZipFile processa um arquivo ZIP
func (p *Processor) processZipFile(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err := p.processCSVFromZip(f); err != nil {
			return err
		}
	}

	return nil
}

// processCSVFromZip processa um CSV dentro do ZIP
func (p *Processor) processCSVFromZip(f *zip.File) error {
	db := p.dbMgr.GetDB()
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

	// Prepara statement com normalizador
	stmt, normalizer, err := p.prepareInsertStatement(tx, tableName)
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

		if err := p.insertRecordWithNormalization(stmt, normalizer, tableName, record); err != nil {
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

// createIndexes cria os índices
func (p *Processor) createIndexes() error {
	db := p.dbMgr.GetDB()
	
	var indexes []string
	if p.dbMgr.IsPostgreSQL() {
		indexes = GetIndexesPostgreSQL()
	} else {
		indexes = GetIndexesSQLite()
	}
	
	for _, idx := range indexes {
		fmt.Printf("  %s\n", idx)
		if _, err := db.Exec(idx); err != nil {
			return err
		}
	}
	
	// Atualiza CNPJ dos sócios (apenas se não for PostgreSQL, pois lá já é feito durante import)
	if !p.dbMgr.IsPostgreSQL() {
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
	
	return nil
}

// printStats imprime estatísticas
func (p *Processor) printStats() error {
	db := p.dbMgr.GetDB()
	fmt.Println("\n📊 Estatísticas:")
	
	tables := []string{"empresas", "estabelecimento", "socios", "simples"}
	for _, table := range tables {
		var count int
		tableWithPrefix := p.dbMgr.TablePrefix(table)
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableWithPrefix)).Scan(&count)
		if err != nil {
			return err
		}
		fmt.Printf("  %s: %d registros\n", table, count)
	}

	// Data de referência
	var dataRef string
	tableEstab := p.dbMgr.TablePrefix("estabelecimento")
	err := db.QueryRow(fmt.Sprintf("SELECT data_situacao_cadastral FROM %s LIMIT 1", tableEstab)).Scan(&dataRef)
	if err == nil && len(dataRef) >= 8 {
		fmt.Printf("\n📅 Data de referência: %s/%s/%s\n", dataRef[6:8], dataRef[4:6], dataRef[0:4])
	}

	return nil
}
