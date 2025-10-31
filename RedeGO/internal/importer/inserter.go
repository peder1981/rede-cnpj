package importer

import (
	"database/sql"
	"fmt"
	"strings"
)

// prepareInsertStatement prepara o statement de insert com normalização
func (p *Processor) prepareInsertStatement(tx *sql.Tx, tableName string) (*sql.Stmt, *Normalizer, error) {
	// Cria normalizador específico para a tabela
	var normalizer *Normalizer
	switch tableName {
	case "empresas":
		normalizer = GetEmpresasNormalizer()
	case "estabelecimento":
		normalizer = GetEstabelecimentoNormalizer()
	case "socios":
		normalizer = GetSociosNormalizer()
	case "simples":
		normalizer = GetSimplesNormalizer()
	default:
		normalizer = NewNormalizer() // Normalizador vazio para tabelas de lookup
	}
	
	// Prepara query baseado no tipo de banco
	var query string
	tableWithPrefix := p.dbMgr.TablePrefix(tableName)
	
	if p.dbMgr.IsPostgreSQL() {
		query = p.getInsertQueryPostgreSQL(tableWithPrefix, tableName)
	} else {
		query = p.getInsertQuerySQLite(tableName)
	}
	
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao preparar statement para %s: %w", tableName, err)
	}
	
	return stmt, normalizer, nil
}

// getInsertQuerySQLite retorna a query de insert para SQLite
func (p *Processor) getInsertQuerySQLite(tableName string) string {
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
	return queries[tableName]
}

// getInsertQueryPostgreSQL retorna a query de insert para PostgreSQL
func (p *Processor) getInsertQueryPostgreSQL(tableWithPrefix, tableName string) string {
	var query string
	
	switch tableName {
	case "cnae", "motivo", "municipio", "natureza_juridica", "pais", "qualificacao_socio":
		query = fmt.Sprintf(`
			INSERT INTO %s (codigo, descricao) 
			VALUES ($1, $2) 
			ON CONFLICT (codigo) DO NOTHING
		`, tableWithPrefix)
		
	case "empresas":
		query = fmt.Sprintf(`
			INSERT INTO %s (
				cnpj_basico, razao_social, natureza_juridica,
				qualificacao_responsavel, capital_social, porte_empresa,
				ente_federativo_responsavel
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (cnpj_basico) DO NOTHING
		`, tableWithPrefix)
		
	case "estabelecimento":
		query = fmt.Sprintf(`
			INSERT INTO %s (
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
		`, tableWithPrefix)
		
	case "socios":
		query = fmt.Sprintf(`
			INSERT INTO %s (
				cnpj, cnpj_basico, identificador_de_socio,
				nome_socio, cnpj_cpf_socio, qualificacao_socio,
				data_entrada_sociedade, pais,
				representante_legal, nome_representante,
				qualificacao_representante_legal, faixa_etaria
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, tableWithPrefix)
		
	case "simples":
		query = fmt.Sprintf(`
			INSERT INTO %s (
				cnpj_basico, opcao_simples, data_opcao_simples,
				data_exclusao_simples, opcao_mei, data_opcao_mei,
				data_exclusao_mei
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (cnpj_basico) DO NOTHING
		`, tableWithPrefix)
	}
	
	return query
}

// insertRecordWithNormalization insere um registro com normalização
func (p *Processor) insertRecordWithNormalization(stmt *sql.Stmt, normalizer *Normalizer, tableName string, record []string) error {
	// Tabelas de lookup (sem normalização complexa)
	if tableName == "cnae" || tableName == "motivo" || tableName == "municipio" || 
	   tableName == "natureza_juridica" || tableName == "pais" || tableName == "qualificacao_socio" {
		if len(record) < 2 {
			return nil
		}
		codigo := sanitizeString(record[0])
		descricao := sanitizeString(record[1])
		_, err := stmt.Exec(codigo, descricao)
		return err
	}
	
	// Tabelas principais com normalização
	switch tableName {
	case "empresas":
		return p.insertEmpresa(stmt, normalizer, record)
	case "estabelecimento":
		return p.insertEstabelecimento(stmt, normalizer, record)
	case "socios":
		return p.insertSocio(stmt, normalizer, record)
	case "simples":
		return p.insertSimples(stmt, normalizer, record)
	}
	
	return nil
}

// insertEmpresa insere uma empresa com normalização
func (p *Processor) insertEmpresa(stmt *sql.Stmt, normalizer *Normalizer, record []string) error {
	if len(record) < 7 {
		return nil
	}
	
	cnpjBasico := normalizer.NormalizeString("cnpj_basico", record[0])
	razaoSocial := normalizer.NormalizeString("razao_social", record[1])
	natureza := normalizer.NormalizeString("natureza_juridica", record[2])
	qualif := normalizer.NormalizeString("qualificacao_responsavel", record[3])
	
	// Capital social - converte string para float
	capitalStr := strings.ReplaceAll(record[4], ",", ".")
	var capital sql.NullFloat64
	if capitalStr != "" && capitalStr != "0" {
		fmt.Sscanf(capitalStr, "%f", &capital.Float64)
		capital.Valid = true
	}
	capital = normalizer.NormalizeFloat64("capital_social", capital)
	
	porte := normalizer.NormalizeString("porte_empresa", record[5])
	ente := normalizer.NormalizeString("ente_federativo_responsavel", record[6])
	
	_, err := stmt.Exec(cnpjBasico, razaoSocial, natureza, qualif, capital, porte, ente)
	return err
}

// insertEstabelecimento insere um estabelecimento com normalização
func (p *Processor) insertEstabelecimento(stmt *sql.Stmt, normalizer *Normalizer, record []string) error {
	if len(record) < 30 {
		return nil
	}
	
	// Monta CNPJ completo
	cnpj := record[0] + record[1] + record[2]
	
	// Normaliza todos os campos
	cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)
	cnpjBasico := normalizer.NormalizeString("cnpj_basico", record[0])
	cnpjOrdem := normalizer.NormalizeString("cnpj_ordem", record[1])
	cnpjDv := normalizer.NormalizeString("cnpj_dv", record[2])
	matrizFilial := normalizer.NormalizeString("matriz_filial", record[3])
	nomeFantasia := normalizer.NormalizeString("nome_fantasia", record[4])
	situacao := normalizer.NormalizeString("situacao_cadastral", record[5])
	dataSituacao := normalizer.NormalizeString("data_situacao_cadastral", record[6])
	motivo := normalizer.NormalizeString("motivo_situacao_cadastral", record[7])
	cidadeExterior := normalizer.NormalizeString("nome_cidade_exterior", record[8])
	pais := normalizer.NormalizeString("pais", record[9])
	dataInicio := normalizer.NormalizeString("data_inicio_atividades", record[10])
	cnae := normalizer.NormalizeString("cnae_fiscal", record[11])
	cnaeSecundaria := normalizer.NormalizeString("cnae_fiscal_secundaria", record[12])
	tipoLogradouro := normalizer.NormalizeString("tipo_logradouro", record[13])
	logradouro := normalizer.NormalizeString("logradouro", record[14])
	numero := normalizer.NormalizeString("numero", record[15])
	complemento := normalizer.NormalizeString("complemento", record[16])
	bairro := normalizer.NormalizeString("bairro", record[17])
	cep := normalizer.NormalizeString("cep", record[18])
	uf := normalizer.NormalizeString("uf", record[19])
	municipio := normalizer.NormalizeString("municipio", record[20])
	ddd1 := normalizer.NormalizeString("ddd1", record[21])
	tel1 := normalizer.NormalizeString("telefone1", record[22])
	ddd2 := normalizer.NormalizeString("ddd2", record[23])
	tel2 := normalizer.NormalizeString("telefone2", record[24])
	dddFax := normalizer.NormalizeString("ddd_fax", record[25])
	fax := normalizer.NormalizeString("fax", record[26])
	email := normalizer.NormalizeString("correio_eletronico", record[27])
	situacaoEspecial := normalizer.NormalizeString("situacao_especial", record[28])
	dataEspecial := normalizer.NormalizeString("data_situacao_especial", record[29])
	
	// Valida campos obrigatórios
	if !cnpjNorm.Valid || !cnpjBasico.Valid || !cnpjOrdem.Valid || !cnpjDv.Valid || !uf.Valid {
		return nil // Ignora registro inválido
	}
	
	_, err := stmt.Exec(
		cnpjNorm, cnpjBasico, cnpjOrdem, cnpjDv,
		matrizFilial, nomeFantasia, situacao,
		dataSituacao, motivo,
		cidadeExterior, pais, dataInicio,
		cnae, cnaeSecundaria,
		tipoLogradouro, logradouro, numero, complemento,
		bairro, cep, uf, municipio,
		ddd1, tel1, ddd2, tel2,
		dddFax, fax, email,
		situacaoEspecial, dataEspecial,
	)
	return err
}

// insertSocio insere um sócio com normalização
func (p *Processor) insertSocio(stmt *sql.Stmt, normalizer *Normalizer, record []string) error {
	if len(record) < 11 {
		return nil
	}
	
	// Para PostgreSQL, precisamos buscar o CNPJ da matriz
	// Para SQLite, deixamos vazio e atualizamos depois
	cnpj := sql.NullString{Valid: false}
	if p.dbMgr.IsPostgreSQL() {
		// Busca CNPJ da matriz
		db := p.dbMgr.GetDB()
		cnpjBasico := record[0]
		tableEstab := p.dbMgr.TablePrefix("estabelecimento")
		err := db.QueryRow(
			fmt.Sprintf("SELECT cnpj FROM %s WHERE cnpj_basico = $1 AND matriz_filial = '1' LIMIT 1", tableEstab),
			cnpjBasico,
		).Scan(&cnpj)
		if err != nil {
			cnpj = sql.NullString{Valid: false}
		}
	}
	
	cnpjBasico := normalizer.NormalizeString("cnpj_basico", record[0])
	identificador := normalizer.NormalizeString("identificador_de_socio", record[1])
	nome := normalizer.NormalizeString("nome_socio", record[2])
	cpfCnpj := normalizer.NormalizeString("cnpj_cpf_socio", record[3])
	qualif := normalizer.NormalizeString("qualificacao_socio", record[4])
	dataEntrada := normalizer.NormalizeString("data_entrada_sociedade", record[5])
	pais := normalizer.NormalizeString("pais", record[6])
	repLegal := normalizer.NormalizeString("representante_legal", record[7])
	nomeRep := normalizer.NormalizeString("nome_representante", record[8])
	qualifRep := normalizer.NormalizeString("qualificacao_representante_legal", record[9])
	faixaEtaria := normalizer.NormalizeString("faixa_etaria", record[10])
	
	_, err := stmt.Exec(
		cnpj, cnpjBasico, identificador,
		nome, cpfCnpj, qualif,
		dataEntrada, pais,
		repLegal, nomeRep,
		qualifRep, faixaEtaria,
	)
	return err
}

// insertSimples insere um registro do Simples com normalização
func (p *Processor) insertSimples(stmt *sql.Stmt, normalizer *Normalizer, record []string) error {
	if len(record) < 7 {
		return nil
	}
	
	cnpjBasico := normalizer.NormalizeString("cnpj_basico", record[0])
	opcaoSimples := normalizer.NormalizeString("opcao_simples", record[1])
	dataOpcaoSimples := normalizer.NormalizeString("data_opcao_simples", record[2])
	dataExclusaoSimples := normalizer.NormalizeString("data_exclusao_simples", record[3])
	opcaoMei := normalizer.NormalizeString("opcao_mei", record[4])
	dataOpcaoMei := normalizer.NormalizeString("data_opcao_mei", record[5])
	dataExclusaoMei := normalizer.NormalizeString("data_exclusao_mei", record[6])
	
	_, err := stmt.Exec(
		cnpjBasico, opcaoSimples, dataOpcaoSimples,
		dataExclusaoSimples, opcaoMei, dataOpcaoMei,
		dataExclusaoMei,
	)
	return err
}
