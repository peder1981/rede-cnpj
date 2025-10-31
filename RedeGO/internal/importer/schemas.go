package importer

// GetTableSchemas retorna os schemas das tabelas para SQLite
func GetTableSchemasSQLite() map[string]string {
	return map[string]string{
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
}

// GetTableSchemasPostgreSQL retorna os schemas das tabelas para PostgreSQL
func GetTableSchemasPostgreSQL() map[string]string {
	return map[string]string{
		"receita.cnae": `CREATE TABLE IF NOT EXISTS receita.cnae (
			codigo VARCHAR(7) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.motivo": `CREATE TABLE IF NOT EXISTS receita.motivo (
			codigo VARCHAR(2) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.municipio": `CREATE TABLE IF NOT EXISTS receita.municipio (
			codigo VARCHAR(4) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.natureza_juridica": `CREATE TABLE IF NOT EXISTS receita.natureza_juridica (
			codigo VARCHAR(4) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.pais": `CREATE TABLE IF NOT EXISTS receita.pais (
			codigo VARCHAR(3) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.qualificacao_socio": `CREATE TABLE IF NOT EXISTS receita.qualificacao_socio (
			codigo VARCHAR(2) PRIMARY KEY,
			descricao TEXT NOT NULL
		)`,
		
		"receita.empresas": `CREATE TABLE IF NOT EXISTS receita.empresas (
			cnpj_basico VARCHAR(8) PRIMARY KEY,
			razao_social TEXT NOT NULL,
			natureza_juridica VARCHAR(4),
			qualificacao_responsavel VARCHAR(2),
			capital_social NUMERIC(15,2),
			porte_empresa VARCHAR(2),
			ente_federativo_responsavel TEXT
		)`,
		
		"receita.estabelecimento": `CREATE TABLE IF NOT EXISTS receita.estabelecimento (
			cnpj VARCHAR(14) NOT NULL,
			cnpj_basico VARCHAR(8) NOT NULL,
			cnpj_ordem VARCHAR(4) NOT NULL,
			cnpj_dv VARCHAR(2) NOT NULL,
			matriz_filial VARCHAR(1),
			nome_fantasia TEXT,
			situacao_cadastral VARCHAR(2),
			data_situacao_cadastral DATE,
			motivo_situacao_cadastral VARCHAR(2),
			nome_cidade_exterior TEXT,
			pais VARCHAR(3),
			data_inicio_atividades DATE,
			cnae_fiscal VARCHAR(7),
			cnae_fiscal_secundaria TEXT,
			tipo_logradouro TEXT,
			logradouro TEXT,
			numero TEXT,
			complemento TEXT,
			bairro TEXT,
			cep VARCHAR(8),
			uf VARCHAR(2) NOT NULL,
			municipio VARCHAR(4),
			ddd1 VARCHAR(4),
			telefone1 VARCHAR(8),
			ddd2 VARCHAR(4),
			telefone2 VARCHAR(8),
			ddd_fax VARCHAR(4),
			fax VARCHAR(8),
			correio_eletronico TEXT,
			situacao_especial TEXT,
			data_situacao_especial DATE,
			PRIMARY KEY (cnpj, uf)
		)`,
		
		"receita.socios": `CREATE TABLE IF NOT EXISTS receita.socios (
			cnpj VARCHAR(14) NOT NULL,
			cnpj_basico VARCHAR(8) NOT NULL,
			identificador_de_socio VARCHAR(1) NOT NULL,
			nome_socio TEXT NOT NULL,
			cnpj_cpf_socio VARCHAR(14) NOT NULL,
			qualificacao_socio VARCHAR(2),
			data_entrada_sociedade DATE,
			pais VARCHAR(3),
			representante_legal VARCHAR(11),
			nome_representante TEXT,
			qualificacao_representante_legal VARCHAR(2),
			faixa_etaria VARCHAR(1)
		)`,
		
		"receita.simples": `CREATE TABLE IF NOT EXISTS receita.simples (
			cnpj_basico VARCHAR(8) PRIMARY KEY,
			opcao_simples VARCHAR(1),
			data_opcao_simples DATE,
			data_exclusao_simples DATE,
			opcao_mei VARCHAR(1),
			data_opcao_mei DATE,
			data_exclusao_mei DATE
		)`,
	}
}

// GetIndexesSQLite retorna os índices para SQLite
func GetIndexesSQLite() []string {
	return []string{
		"CREATE INDEX IF NOT EXISTS idx_empresas_cnpj_basico ON empresas(cnpj_basico)",
		"CREATE INDEX IF NOT EXISTS idx_empresas_razao_social ON empresas(razao_social)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_cnpj_basico ON estabelecimento(cnpj_basico)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_cnpj ON estabelecimento(cnpj)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_nomefantasia ON estabelecimento(nome_fantasia)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj_basico ON socios(cnpj_basico)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj ON socios(cnpj)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj_cpf_socio ON socios(cnpj_cpf_socio)",
		"CREATE INDEX IF NOT EXISTS idx_socios_nome_socio ON socios(nome_socio)",
		"CREATE INDEX IF NOT EXISTS idx_simples_cnpj_basico ON simples(cnpj_basico)",
	}
}

// GetIndexesPostgreSQL retorna os índices para PostgreSQL
func GetIndexesPostgreSQL() []string {
	return []string{
		"CREATE INDEX IF NOT EXISTS idx_empresas_razao_social ON receita.empresas(razao_social)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_cnpj_basico ON receita.estabelecimento(cnpj_basico)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_cnpj ON receita.estabelecimento(cnpj)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_nomefantasia ON receita.estabelecimento(nome_fantasia)",
		"CREATE INDEX IF NOT EXISTS idx_estabelecimento_uf ON receita.estabelecimento(uf)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj_basico ON receita.socios(cnpj_basico)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj ON receita.socios(cnpj)",
		"CREATE INDEX IF NOT EXISTS idx_socios_cnpj_cpf_socio ON receita.socios(cnpj_cpf_socio)",
		"CREATE INDEX IF NOT EXISTS idx_socios_nome_socio ON receita.socios(nome_socio)",
	}
}
