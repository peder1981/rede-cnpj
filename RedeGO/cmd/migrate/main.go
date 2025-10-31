package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	sqliteDB   = "bases/cnpj.db"
	postgresDB = "postgresql://rede_user:rede_cnpj_2025@localhost:5433/rede_cnpj?sslmode=disable"
	batchSize  = 10000
)

type MigrationStats struct {
	TableName    string
	TotalRows    int64
	MigratedRows int64
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
}

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

// sanitizeNumericString normaliza campos numéricos
// Converte strings vazias ou inválidas em NULL
func sanitizeNumericString(ns sql.NullString) sql.NullString {
	if !ns.Valid {
		return sql.NullString{Valid: false}
	}
	
	trimmed := strings.TrimSpace(ns.String)
	
	// Valores inválidos
	if trimmed == "" || trimmed == "0" || trimmed == "00" {
		return sql.NullString{Valid: false}
	}
	
	return sql.NullString{
		String: trimmed,
		Valid:  true,
	}
}

// sanitizeDateString normaliza e valida campos de data
// Converte valores inválidos (vazios, "0", "00000000", etc) em NULL
func sanitizeDateString(ns sql.NullString) sql.NullString {
	// Se não é válido, retorna NULL
	if !ns.Valid {
		return sql.NullString{Valid: false}
	}
	
	// Remove espaços em branco
	trimmed := strings.TrimSpace(ns.String)
	
	// Valores inválidos que devem ser convertidos para NULL
	invalidValues := []string{"", "0", "00", "000", "0000", "00000000"}
	for _, invalid := range invalidValues {
		if trimmed == invalid {
			return sql.NullString{Valid: false}
		}
	}
	
	// Valida formato de data (YYYYMMDD - 8 dígitos)
	if len(trimmed) == 8 {
		// Verifica se todos são dígitos
		for _, c := range trimmed {
			if c < '0' || c > '9' {
				// Formato inválido, retorna NULL
				return sql.NullString{Valid: false}
			}
		}
		
		// Valida componentes da data
		year := trimmed[0:4]
		month := trimmed[4:6]
		day := trimmed[6:8]
		
		// Verifica se ano é válido (maior que 1900 e menor que 2100)
		if year < "1900" || year > "2100" {
			return sql.NullString{Valid: false}
		}
		
		// Verifica se mês é válido (01-12)
		if month < "01" || month > "12" {
			return sql.NullString{Valid: false}
		}
		
		// Verifica se dia é válido (01-31)
		if day < "01" || day > "31" {
			return sql.NullString{Valid: false}
		}
		
		// Data válida, retorna formatada como YYYY-MM-DD para PostgreSQL
		return sql.NullString{
			String: year + "-" + month + "-" + day,
			Valid:  true,
		}
	}
	
	// Se não tem 8 dígitos, tenta outros formatos comuns
	// Formato YYYY-MM-DD já está correto
	if len(trimmed) == 10 && trimmed[4] == '-' && trimmed[7] == '-' {
		return sql.NullString{
			String: trimmed,
			Valid:  true,
		}
	}
	
	// Formato inválido, retorna NULL
	return sql.NullString{Valid: false}
}

func main() {
	printHeader()

	// Conectar SQLite
	log.Println("📂 Conectando ao SQLite...")
	srcDB, err := sql.Open("sqlite3", sqliteDB)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar SQLite: %v", err)
	}
	defer srcDB.Close()

	// Conectar PostgreSQL
	log.Println("🐘 Conectando ao PostgreSQL...")
	dstDB, err := sql.Open("postgres", postgresDB)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar PostgreSQL: %v", err)
	}
	defer dstDB.Close()

	// Testar conexões
	if err := srcDB.Ping(); err != nil {
		log.Fatalf("❌ SQLite não acessível: %v", err)
	}
	if err := dstDB.Ping(); err != nil {
		log.Fatalf("❌ PostgreSQL não acessível: %v", err)
	}

	log.Println("✅ Conexões estabelecidas!")
	log.Println("")

	// Migrar tabelas
	stats := []MigrationStats{}

	// 1. Empresas
	stat := migrateEmpresas(srcDB, dstDB)
	stats = append(stats, stat)

	// 2. Estabelecimentos
	stat = migrateEstabelecimentos(srcDB, dstDB)
	stats = append(stats, stat)

	// 3. Sócios
	stat = migrateSocios(srcDB, dstDB)
	stats = append(stats, stat)

	// 4. Simples
	stat = migrateSimples(srcDB, dstDB)
	stats = append(stats, stat)

	// 5. Tabelas de códigos
	migrateLookupTables(srcDB, dstDB)

	// 6. Rede (se existir)
	if tableExists(srcDB, "ligacao") {
		stat = migrateRede(srcDB, dstDB)
		stats = append(stats, stat)
	}

	// Resumo final
	printSummary(stats)
}

func migrateEmpresas(src, dst *sql.DB) MigrationStats {
	stat := MigrationStats{
		TableName: "empresas",
		StartTime: time.Now(),
	}

	log.Println("📊 Migrando tabela: empresas")

	// Criar normalizador específico para empresas
	normalizer := GetEmpresasNormalizer()

	// Contar registros
	err := src.QueryRow("SELECT COUNT(*) FROM empresas").Scan(&stat.TotalRows)
	if err != nil {
		log.Printf("⚠️  Erro ao contar empresas: %v", err)
		return stat
	}

	log.Printf("   Total de registros: %d", stat.TotalRows)

	// Query de origem
	rows, err := src.Query(`
		SELECT 
			cnpj_basico,
			razao_social,
			natureza_juridica,
			qualificacao_responsavel,
			capital_social,
			porte_empresa,
			ente_federativo_responsavel
		FROM empresas
	`)
	if err != nil {
		log.Printf("❌ Erro ao ler empresas: %v", err)
		return stat
	}
	defer rows.Close()

	// Preparar insert
	stmt, err := dst.Prepare(`
		INSERT INTO receita.empresas (
			cnpj_basico,
			razao_social,
			natureza_juridica,
			qualificacao_responsavel,
			capital_social,
			porte_empresa,
			ente_federativo_responsavel
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (cnpj_basico) DO NOTHING
	`)
	if err != nil {
		log.Printf("❌ Erro ao preparar insert: %v", err)
		return stat
	}
	defer stmt.Close()

	// Migrar em batches
	tx, _ := dst.Begin()
	count := int64(0)
	batchCount := 0

	for rows.Next() {
		var cnpj, razao, natureza, qualif, porte, ente sql.NullString
		var capital sql.NullFloat64

		err := rows.Scan(&cnpj, &razao, &natureza, &qualif, &capital, &porte, &ente)
		if err != nil {
			log.Printf("⚠️  Erro ao ler linha: %v", err)
			continue
		}

		// Normaliza campos usando metadados específicos da tabela
		cnpjNorm := normalizer.NormalizeNullString("cnpj_basico", cnpj)
		razaoNorm := normalizer.NormalizeNullString("razao_social", razao)
		naturezaNorm := normalizer.NormalizeNullString("natureza_juridica", natureza)
		qualifNorm := normalizer.NormalizeNullString("qualificacao_responsavel", qualif)
		porteNorm := normalizer.NormalizeNullString("porte_empresa", porte)
		enteNorm := normalizer.NormalizeNullString("ente_federativo_responsavel", ente)
		capitalNorm := normalizer.NormalizeFloat64("capital_social", capital)

		_, err = stmt.Exec(cnpjNorm, razaoNorm, naturezaNorm, qualifNorm, capitalNorm, porteNorm, enteNorm)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir: %v", err)
			continue
		}

		count++
		batchCount++

		if batchCount >= batchSize {
			tx.Commit()
			tx, _ = dst.Begin()
			batchCount = 0
			log.Printf("   Progresso: %d/%d (%.1f%%)", count, stat.TotalRows, float64(count)/float64(stat.TotalRows)*100)
		}
	}

	tx.Commit()
	stat.MigratedRows = count
	stat.EndTime = time.Now()
	stat.Duration = stat.EndTime.Sub(stat.StartTime)

	log.Printf("✅ Empresas migradas: %d em %v", count, stat.Duration)
	log.Println("")

	return stat
}

func migrateEstabelecimentos(src, dst *sql.DB) MigrationStats {
	stat := MigrationStats{
		TableName: "estabelecimentos",
		StartTime: time.Now(),
	}

	log.Println("📊 Migrando tabela: estabelecimento")

	// Criar normalizador específico para estabelecimento
	normalizer := GetEstabelecimentoNormalizer()

	// Contar registros
	err := src.QueryRow("SELECT COUNT(*) FROM estabelecimento").Scan(&stat.TotalRows)
	if err != nil {
		log.Printf("⚠️  Erro ao contar estabelecimentos: %v", err)
		return stat
	}

	log.Printf("   Total de registros: %d", stat.TotalRows)

	// Query de origem
	rows, err := src.Query(`
		SELECT 
			cnpj, cnpj_basico, cnpj_ordem, cnpj_dv,
			matriz_filial, nome_fantasia, situacao_cadastral,
			data_situacao_cadastral, motivo_situacao_cadastral,
			nome_cidade_exterior, pais, data_inicio_atividades,
			cnae_fiscal, cnae_fiscal_secundaria,
			tipo_logradouro, logradouro, numero, complemento,
			bairro, cep, uf, municipio,
			ddd1, telefone1, ddd2, telefone2,
			ddd_fax, fax, correio_eletronico,
			situacao_especial, data_situacao_especial
		FROM estabelecimento
	`)
	if err != nil {
		log.Printf("❌ Erro ao ler estabelecimentos: %v", err)
		return stat
	}
	defer rows.Close()

	// Preparar insert
	stmt, err := dst.Prepare(`
		INSERT INTO receita.estabelecimento (
			cnpj, cnpj_basico, cnpj_ordem, cnpj_dv,
			matriz_filial, nome_fantasia, situacao_cadastral,
			data_situacao_cadastral, motivo_situacao_cadastral,
			nome_cidade_exterior, pais, data_inicio_atividades,
			cnae_fiscal, cnae_fiscal_secundaria,
			tipo_logradouro, logradouro, numero, complemento,
			bairro, cep, uf, municipio,
			ddd1, telefone1, ddd2, telefone2,
			ddd_fax, fax, correio_eletronico,
			situacao_especial, data_situacao_especial
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31
		)
		ON CONFLICT (cnpj, uf) DO NOTHING
	`)
	if err != nil {
		log.Printf("❌ Erro ao preparar insert: %v", err)
		return stat
	}
	defer stmt.Close()

	// Migrar em batches
	tx, _ := dst.Begin()
	count := int64(0)
	batchCount := 0

	for rows.Next() {
		var cnpj, cnpjBasico, cnpjOrdem, cnpjDv, matrizFilial, nomeFantasia string
		var situacaoCadastral, motivoSituacao, nomeCidadeExterior, pais string
		var cnae, cnaeSecundaria, tipoLogradouro, logradouro, numero string
		var complemento, bairro, cep, uf, municipio string
		var ddd1, tel1, ddd2, tel2, dddFax, fax, email string
		var situacaoEspecial string
		var dataSituacao, dataInicio, dataEspecial sql.NullString

		err := rows.Scan(
			&cnpj, &cnpjBasico, &cnpjOrdem, &cnpjDv,
			&matrizFilial, &nomeFantasia, &situacaoCadastral,
			&dataSituacao, &motivoSituacao,
			&nomeCidadeExterior, &pais, &dataInicio,
			&cnae, &cnaeSecundaria,
			&tipoLogradouro, &logradouro, &numero, &complemento,
			&bairro, &cep, &uf, &municipio,
			&ddd1, &tel1, &ddd2, &tel2,
			&dddFax, &fax, &email,
			&situacaoEspecial, &dataEspecial,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao ler linha: %v", err)
			continue
		}

		// Normaliza campos usando metadados específicos da tabela
		cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)
		cnpjBasicoNorm := normalizer.NormalizeString("cnpj_basico", cnpjBasico)
		cnpjOrdemNorm := normalizer.NormalizeString("cnpj_ordem", cnpjOrdem)
		cnpjDvNorm := normalizer.NormalizeString("cnpj_dv", cnpjDv)
		matrizFilialNorm := normalizer.NormalizeString("matriz_filial", matrizFilial)
		nomeFantasiaNorm := normalizer.NormalizeString("nome_fantasia", nomeFantasia)
		situacaoCadastralNorm := normalizer.NormalizeString("situacao_cadastral", situacaoCadastral)
		motivoSituacaoNorm := normalizer.NormalizeString("motivo_situacao_cadastral", motivoSituacao)
		nomeCidadeExteriorNorm := normalizer.NormalizeString("nome_cidade_exterior", nomeCidadeExterior)
		paisNorm := normalizer.NormalizeString("pais", pais)
		cnaeNorm := normalizer.NormalizeString("cnae_fiscal", cnae)
		cnaeSecundariaNorm := normalizer.NormalizeString("cnae_fiscal_secundaria", cnaeSecundaria)
		tipoLogradouroNorm := normalizer.NormalizeString("tipo_logradouro", tipoLogradouro)
		logradouroNorm := normalizer.NormalizeString("logradouro", logradouro)
		numeroNorm := normalizer.NormalizeString("numero", numero)
		complementoNorm := normalizer.NormalizeString("complemento", complemento)
		bairroNorm := normalizer.NormalizeString("bairro", bairro)
		cepNorm := normalizer.NormalizeString("cep", cep)
		ufNorm := normalizer.NormalizeString("uf", uf)
		municipioNorm := normalizer.NormalizeString("municipio", municipio)
		ddd1Norm := normalizer.NormalizeString("ddd1", ddd1)
		tel1Norm := normalizer.NormalizeString("telefone1", tel1)
		ddd2Norm := normalizer.NormalizeString("ddd2", ddd2)
		tel2Norm := normalizer.NormalizeString("telefone2", tel2)
		dddFaxNorm := normalizer.NormalizeString("ddd_fax", dddFax)
		faxNorm := normalizer.NormalizeString("fax", fax)
		emailNorm := normalizer.NormalizeString("correio_eletronico", email)
		situacaoEspecialNorm := normalizer.NormalizeString("situacao_especial", situacaoEspecial)
		dataSituacaoNorm := normalizer.NormalizeNullString("data_situacao_cadastral", dataSituacao)
		dataInicioNorm := normalizer.NormalizeNullString("data_inicio_atividades", dataInicio)
		dataEspecialNorm := normalizer.NormalizeNullString("data_situacao_especial", dataEspecial)

		// Validar campos obrigatórios antes de inserir
		if !cnpjNorm.Valid || cnpjNorm.String == "" {
			log.Printf("⚠️  CNPJ vazio ou inválido: original='%s'", cnpj)
			continue
		}
		if !cnpjBasicoNorm.Valid || cnpjBasicoNorm.String == "" {
			log.Printf("⚠️  CNPJ Básico vazio: CNPJ=%s, original='%s'", cnpj, cnpjBasico)
			continue
		}
		if !cnpjOrdemNorm.Valid || cnpjOrdemNorm.String == "" {
			log.Printf("⚠️  CNPJ Ordem vazio: CNPJ=%s, original='%s'", cnpj, cnpjOrdem)
			continue
		}
		if !cnpjDvNorm.Valid || cnpjDvNorm.String == "" {
			log.Printf("⚠️  CNPJ DV vazio: CNPJ=%s, original='%s'", cnpj, cnpjDv)
			continue
		}
		if !ufNorm.Valid || ufNorm.String == "" {
			log.Printf("⚠️  UF vazio: CNPJ=%s, original='%s'", cnpj, uf)
			continue
		}

		_, err = stmt.Exec(
			cnpjNorm, cnpjBasicoNorm, cnpjOrdemNorm, cnpjDvNorm,
			matrizFilialNorm, nomeFantasiaNorm, situacaoCadastralNorm,
			dataSituacaoNorm, motivoSituacaoNorm,
			nomeCidadeExteriorNorm, paisNorm, dataInicioNorm,
			cnaeNorm, cnaeSecundariaNorm,
			tipoLogradouroNorm, logradouroNorm, numeroNorm, complementoNorm,
			bairroNorm, cepNorm, ufNorm, municipioNorm,
			ddd1Norm, tel1Norm, ddd2Norm, tel2Norm,
			dddFaxNorm, faxNorm, emailNorm,
			situacaoEspecialNorm, dataEspecialNorm,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir CNPJ %s: %v", cnpj, err)
			log.Printf("    Campos: basico='%s', ordem='%s', dv='%s', uf='%s'",
				cnpjBasicoNorm.String, cnpjOrdemNorm.String, cnpjDvNorm.String, ufNorm.String)
			log.Printf("    Datas: situacao=%v, inicio=%v, especial=%v", 
				dataSituacaoNorm, dataInicioNorm, dataEspecialNorm)
			continue
		}

		count++
		batchCount++

		if batchCount >= batchSize {
			tx.Commit()
			tx, _ = dst.Begin()
			batchCount = 0
			log.Printf("   Progresso: %d/%d (%.1f%%)", count, stat.TotalRows, float64(count)/float64(stat.TotalRows)*100)
		}
	}

	tx.Commit()
	stat.MigratedRows = count
	stat.EndTime = time.Now()
	stat.Duration = stat.EndTime.Sub(stat.StartTime)

	log.Printf("✅ Estabelecimentos migrados: %d em %v", count, stat.Duration)
	log.Println("")

	return stat
}

func migrateSocios(src, dst *sql.DB) MigrationStats {
	stat := MigrationStats{
		TableName: "socios",
		StartTime: time.Now(),
	}

	log.Println("📊 Migrando tabela: socios")

	// Criar normalizador específico para socios
	normalizer := GetSociosNormalizer()

	// Contar registros
	err := src.QueryRow("SELECT COUNT(*) FROM socios").Scan(&stat.TotalRows)
	if err != nil {
		log.Printf("⚠️  Erro ao contar sócios: %v", err)
		return stat
	}

	log.Printf("   Total de registros: %d", stat.TotalRows)

	// Query de origem
	rows, err := src.Query(`
		SELECT 
			cnpj, cnpj_basico, identificador_de_socio,
			nome_socio, cnpj_cpf_socio, qualificacao_socio,
			data_entrada_sociedade, pais,
			representante_legal, nome_representante,
			qualificacao_representante_legal, faixa_etaria
		FROM socios
	`)
	if err != nil {
		log.Printf("❌ Erro ao ler sócios: %v", err)
		return stat
	}
	defer rows.Close()

	// Preparar insert
	stmt, err := dst.Prepare(`
		INSERT INTO receita.socios (
			cnpj, cnpj_basico, identificador_de_socio,
			nome_socio, cnpj_cpf_socio, qualificacao_socio,
			data_entrada_sociedade, pais,
			representante_legal, nome_representante,
			qualificacao_representante_legal, faixa_etaria
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)
	if err != nil {
		log.Printf("❌ Erro ao preparar insert: %v", err)
		return stat
	}
	defer stmt.Close()

	// Migrar em batches
	tx, _ := dst.Begin()
	count := int64(0)
	batchCount := 0

	for rows.Next() {
		var cnpj, cnpjBasico, identificador, nome, cpfCnpj string
		var qualif, pais, repLegal, nomeRep, qualifRep, faixaEtaria string
		var dataEntrada sql.NullString

		err := rows.Scan(
			&cnpj, &cnpjBasico, &identificador,
			&nome, &cpfCnpj, &qualif,
			&dataEntrada, &pais,
			&repLegal, &nomeRep,
			&qualifRep, &faixaEtaria,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao ler linha: %v", err)
			continue
		}

		// Normaliza campos usando metadados específicos da tabela
		cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)
		cnpjBasicoNorm := normalizer.NormalizeString("cnpj_basico", cnpjBasico)
		identificadorNorm := normalizer.NormalizeString("identificador_de_socio", identificador)
		nomeNorm := normalizer.NormalizeString("nome_socio", nome)
		cpfCnpjNorm := normalizer.NormalizeString("cnpj_cpf_socio", cpfCnpj)
		qualifNorm := normalizer.NormalizeString("qualificacao_socio", qualif)
		paisNorm := normalizer.NormalizeString("pais", pais)
		repLegalNorm := normalizer.NormalizeString("representante_legal", repLegal)
		nomeRepNorm := normalizer.NormalizeString("nome_representante", nomeRep)
		qualifRepNorm := normalizer.NormalizeString("qualificacao_representante_legal", qualifRep)
		faixaEtariaNorm := normalizer.NormalizeString("faixa_etaria", faixaEtaria)
		dataEntradaNorm := normalizer.NormalizeNullString("data_entrada_sociedade", dataEntrada)

		_, err = stmt.Exec(
			cnpjNorm, cnpjBasicoNorm, identificadorNorm,
			nomeNorm, cpfCnpjNorm, qualifNorm,
			dataEntradaNorm, paisNorm,
			repLegalNorm, nomeRepNorm,
			qualifRepNorm, faixaEtariaNorm,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir: %v", err)
			continue
		}

		count++
		batchCount++

		if batchCount >= batchSize {
			tx.Commit()
			tx, _ = dst.Begin()
			batchCount = 0
			log.Printf("   Progresso: %d/%d (%.1f%%)", count, stat.TotalRows, float64(count)/float64(stat.TotalRows)*100)
		}
	}

	tx.Commit()
	stat.MigratedRows = count
	stat.EndTime = time.Now()
	stat.Duration = stat.EndTime.Sub(stat.StartTime)

	log.Printf("✅ Sócios migrados: %d em %v", count, stat.Duration)
	log.Println("")

	return stat
}

func migrateSimples(src, dst *sql.DB) MigrationStats {
	stat := MigrationStats{
		TableName: "simples",
		StartTime: time.Now(),
	}

	log.Println("📊 Migrando tabela: simples")

	// Criar normalizador específico para simples
	normalizer := GetSimplesNormalizer()

	// Contar registros
	err := src.QueryRow("SELECT COUNT(*) FROM simples").Scan(&stat.TotalRows)
	if err != nil {
		log.Printf("⚠️  Erro ao contar simples: %v", err)
		return stat
	}

	log.Printf("   Total de registros: %d", stat.TotalRows)

	// Query de origem
	rows, err := src.Query(`
		SELECT 
			cnpj_basico, opcao_simples, data_opcao_simples,
			data_exclusao_simples, opcao_mei, data_opcao_mei,
			data_exclusao_mei
		FROM simples
	`)
	if err != nil {
		log.Printf("❌ Erro ao ler simples: %v", err)
		return stat
	}
	defer rows.Close()

	// Preparar insert
	stmt, err := dst.Prepare(`
		INSERT INTO receita.simples (
			cnpj_basico, opcao_simples, data_opcao_simples,
			data_exclusao_simples, opcao_mei, data_opcao_mei,
			data_exclusao_mei
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (cnpj_basico) DO NOTHING
	`)
	if err != nil {
		log.Printf("❌ Erro ao preparar insert: %v", err)
		return stat
	}
	defer stmt.Close()

	// Migrar em batches
	tx, _ := dst.Begin()
	count := int64(0)
	batchCount := 0

	for rows.Next() {
		var cnpjBasico, opcaoSimples, opcaoMei string
		var dataOpcaoSimples, dataExclusaoSimples, dataOpcaoMei, dataExclusaoMei sql.NullString

		err := rows.Scan(
			&cnpjBasico, &opcaoSimples, &dataOpcaoSimples,
			&dataExclusaoSimples, &opcaoMei, &dataOpcaoMei,
			&dataExclusaoMei,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao ler linha: %v", err)
			continue
		}

		// Normaliza campos usando metadados específicos da tabela
		cnpjBasicoNorm := normalizer.NormalizeString("cnpj_basico", cnpjBasico)
		opcaoSimplesNorm := normalizer.NormalizeString("opcao_simples", opcaoSimples)
		dataOpcaoSimplesNorm := normalizer.NormalizeNullString("data_opcao_simples", dataOpcaoSimples)
		dataExclusaoSimplesNorm := normalizer.NormalizeNullString("data_exclusao_simples", dataExclusaoSimples)
		opcaoMeiNorm := normalizer.NormalizeString("opcao_mei", opcaoMei)
		dataOpcaoMeiNorm := normalizer.NormalizeNullString("data_opcao_mei", dataOpcaoMei)
		dataExclusaoMeiNorm := normalizer.NormalizeNullString("data_exclusao_mei", dataExclusaoMei)

		_, err = stmt.Exec(
			cnpjBasicoNorm, opcaoSimplesNorm, dataOpcaoSimplesNorm,
			dataExclusaoSimplesNorm, opcaoMeiNorm, dataOpcaoMeiNorm,
			dataExclusaoMeiNorm,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir: %v", err)
			continue
		}

		count++
		batchCount++

		if batchCount >= batchSize {
			tx.Commit()
			tx, _ = dst.Begin()
			batchCount = 0
			log.Printf("   Progresso: %d/%d (%.1f%%)", count, stat.TotalRows, float64(count)/float64(stat.TotalRows)*100)
		}
	}

	tx.Commit()
	stat.MigratedRows = count
	stat.EndTime = time.Now()
	stat.Duration = stat.EndTime.Sub(stat.StartTime)

	log.Printf("✅ Simples migrados: %d em %v", count, stat.Duration)
	log.Println("")

	return stat
}

func migrateLookupTables(src, dst *sql.DB) {
	tables := []string{"cnae", "motivo", "municipio", "natureza_juridica", "pais", "qualificacao_socio"}

	for _, table := range tables {
		log.Printf("📊 Migrando tabela: %s", table)

		rows, err := src.Query(fmt.Sprintf("SELECT codigo, descricao FROM %s", table))
		if err != nil {
			log.Printf("⚠️  Tabela %s não encontrada", table)
			continue
		}

		stmt, err := dst.Prepare(fmt.Sprintf(`
			INSERT INTO receita.%s (codigo, descricao) 
			VALUES ($1, $2) 
			ON CONFLICT (codigo) DO NOTHING
		`, table))
		if err != nil {
			log.Printf("⚠️  Erro ao preparar insert para %s: %v", table, err)
			rows.Close()
			continue
		}

		count := 0
		for rows.Next() {
			var codigo, descricao string
			err := rows.Scan(&codigo, &descricao)
			if err != nil {
				log.Printf("⚠️  Erro ao ler linha: %v", err)
				continue
			}
			
			// Sanitiza strings antes de inserir
			codigo = sanitizeString(codigo)
			descricao = sanitizeString(descricao)
			
			_, err = stmt.Exec(codigo, descricao)
			if err != nil {
				log.Printf("⚠️  Erro ao inserir em %s: %v", table, err)
				continue
			}
			count++
		}

		rows.Close()
		stmt.Close()
		log.Printf("✅ %s: %d registros", table, count)
	}
	log.Println("")
}

func migrateRede(src, dst *sql.DB) MigrationStats {
	stat := MigrationStats{
		TableName: "ligacao",
		StartTime: time.Now(),
	}

	log.Println("📊 Migrando tabela: ligacao (rede)")

	// Verificar se existe no SQLite
	var count int64
	err := src.QueryRow("SELECT COUNT(*) FROM ligacao").Scan(&count)
	if err != nil {
		log.Printf("⚠️  Tabela ligacao não encontrada no SQLite")
		return stat
	}

	stat.TotalRows = count
	log.Printf("   Total de registros: %d", count)

	// Migrar
	rows, err := src.Query("SELECT id1, id2, descricao, cnpj, peso FROM ligacao")
	if err != nil {
		log.Printf("❌ Erro ao ler ligações: %v", err)
		return stat
	}
	defer rows.Close()

	stmt, err := dst.Prepare(`
		INSERT INTO rede.ligacao (id1, id2, descricao, cnpj, peso)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		log.Printf("❌ Erro ao preparar insert: %v", err)
		return stat
	}
	defer stmt.Close()

	tx, _ := dst.Begin()
	migrated := int64(0)

	for rows.Next() {
		var id1, id2, descricao, cnpj string
		var peso int
		err := rows.Scan(&id1, &id2, &descricao, &cnpj, &peso)
		if err != nil {
			log.Printf("⚠️  Erro ao ler linha: %v", err)
			continue
		}
		
		// Sanitiza strings antes de inserir
		id1 = sanitizeString(id1)
		id2 = sanitizeString(id2)
		descricao = sanitizeString(descricao)
		cnpj = sanitizeString(cnpj)
		
		_, err = stmt.Exec(id1, id2, descricao, cnpj, peso)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir: %v", err)
			continue
		}
		migrated++

		if migrated%batchSize == 0 {
			tx.Commit()
			tx, _ = dst.Begin()
			log.Printf("   Progresso: %d/%d (%.1f%%)", migrated, count, float64(migrated)/float64(count)*100)
		}
	}

	tx.Commit()
	stat.MigratedRows = migrated
	stat.EndTime = time.Now()
	stat.Duration = stat.EndTime.Sub(stat.StartTime)

	log.Printf("✅ Ligações migradas: %d em %v", migrated, stat.Duration)
	log.Println("")

	return stat
}

func tableExists(db *sql.DB, tableName string) bool {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
		tableName,
	).Scan(&count)
	return err == nil && count > 0
}

func printHeader() {
	fmt.Println("")
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║     🔄 RedeCNPJ - Migração SQLite → PostgreSQL                ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println("")
}

func printSummary(stats []MigrationStats) {
	fmt.Println("")
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║     ✅ MIGRAÇÃO CONCLUÍDA COM SUCESSO!                        ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Println("📊 Resumo da Migração:")
	fmt.Println("")

	totalRows := int64(0)
	totalDuration := time.Duration(0)

	for _, stat := range stats {
		if stat.MigratedRows > 0 {
			fmt.Printf("   %-20s: %10d registros em %v\n",
				stat.TableName,
				stat.MigratedRows,
				stat.Duration.Round(time.Second),
			)
			totalRows += stat.MigratedRows
			totalDuration += stat.Duration
		}
	}

	fmt.Println("")
	fmt.Printf("   %-20s: %10d registros\n", "TOTAL", totalRows)
	fmt.Printf("   %-20s: %v\n", "Tempo Total", totalDuration.Round(time.Second))
	fmt.Println("")
}
