package main

import "regexp"

// GetEmpresasNormalizer retorna o normalizador para a tabela empresas
func GetEmpresasNormalizer() *Normalizer {
	n := NewNormalizer()

	// cnpj_basico VARCHAR(8) PRIMARY KEY
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_basico",
		Type:      FieldTypeCode,
		MaxLength: 8,
		Required:  true,
		Pattern:   regexp.MustCompile(`^\d{8}$`),
	})

	// razao_social TEXT NOT NULL
	n.RegisterField(FieldMetadata{
		Name:     "razao_social",
		Type:     FieldTypeText,
		Required: true,
	})

	// natureza_juridica VARCHAR(4)
	n.RegisterField(FieldMetadata{
		Name:      "natureza_juridica",
		Type:      FieldTypeCode,
		MaxLength: 4,
		Required:  false,
	})

	// qualificacao_responsavel VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "qualificacao_responsavel",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// capital_social NUMERIC(15,2)
	n.RegisterField(FieldMetadata{
		Name:     "capital_social",
		Type:     FieldTypeNumeric,
		Required: false,
	})

	// porte_empresa VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "porte_empresa",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// ente_federativo_responsavel TEXT
	n.RegisterField(FieldMetadata{
		Name:     "ente_federativo_responsavel",
		Type:     FieldTypeText,
		Required: false,
	})

	return n
}

// GetEstabelecimentoNormalizer retorna o normalizador para a tabela estabelecimento
func GetEstabelecimentoNormalizer() *Normalizer {
	n := NewNormalizer()

	// cnpj VARCHAR(14) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj",
		Type:      FieldTypeCNPJ,
		MaxLength: 14,
		Required:  true,
	})

	// cnpj_basico VARCHAR(8) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_basico",
		Type:      FieldTypeCode,
		MaxLength: 8,
		Required:  true,
	})

	// cnpj_ordem VARCHAR(4) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_ordem",
		Type:      FieldTypeCode,
		MaxLength: 4,
		Required:  true,
	})

	// cnpj_dv VARCHAR(2) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_dv",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  true,
	})

	// matriz_filial VARCHAR(1)
	n.RegisterField(FieldMetadata{
		Name:        "matriz_filial",
		Type:        FieldTypeCode,
		MaxLength:   1,
		Required:    false,
		ValidValues: []string{"1", "2"}, // 1=Matriz, 2=Filial
	})

	// nome_fantasia TEXT
	n.RegisterField(FieldMetadata{
		Name:     "nome_fantasia",
		Type:     FieldTypeText,
		Required: false,
	})

	// situacao_cadastral VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "situacao_cadastral",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// data_situacao_cadastral DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_situacao_cadastral",
		Type:     FieldTypeDate,
		Required: false,
	})

	// motivo_situacao_cadastral VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "motivo_situacao_cadastral",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// nome_cidade_exterior TEXT
	n.RegisterField(FieldMetadata{
		Name:     "nome_cidade_exterior",
		Type:     FieldTypeText,
		Required: false,
	})

	// pais VARCHAR(3)
	n.RegisterField(FieldMetadata{
		Name:      "pais",
		Type:      FieldTypeCode,
		MaxLength: 3,
		Required:  false,
	})

	// data_inicio_atividades DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_inicio_atividades",
		Type:     FieldTypeDate,
		Required: false,
	})

	// cnae_fiscal VARCHAR(7)
	n.RegisterField(FieldMetadata{
		Name:      "cnae_fiscal",
		Type:      FieldTypeCode,
		MaxLength: 7,
		Required:  false,
	})

	// cnae_fiscal_secundaria TEXT
	n.RegisterField(FieldMetadata{
		Name:     "cnae_fiscal_secundaria",
		Type:     FieldTypeText,
		Required: false,
	})

	// tipo_logradouro TEXT
	n.RegisterField(FieldMetadata{
		Name:     "tipo_logradouro",
		Type:     FieldTypeText,
		Required: false,
	})

	// logradouro TEXT
	n.RegisterField(FieldMetadata{
		Name:     "logradouro",
		Type:     FieldTypeText,
		Required: false,
	})

	// numero TEXT
	n.RegisterField(FieldMetadata{
		Name:     "numero",
		Type:     FieldTypeText,
		Required: false,
	})

	// complemento TEXT
	n.RegisterField(FieldMetadata{
		Name:     "complemento",
		Type:     FieldTypeText,
		Required: false,
	})

	// bairro TEXT
	n.RegisterField(FieldMetadata{
		Name:     "bairro",
		Type:     FieldTypeText,
		Required: false,
	})

	// cep VARCHAR(8)
	n.RegisterField(FieldMetadata{
		Name:      "cep",
		Type:      FieldTypeCEP,
		MaxLength: 8,
		Required:  false,
	})

	// uf VARCHAR(2) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "uf",
		Type:      FieldTypeUF,
		MaxLength: 2,
		Required:  true,
	})

	// municipio VARCHAR(4)
	n.RegisterField(FieldMetadata{
		Name:      "municipio",
		Type:      FieldTypeCode,
		MaxLength: 4,
		Required:  false,
	})

	// ddd1 VARCHAR(4)
	n.RegisterField(FieldMetadata{
		Name:      "ddd1",
		Type:      FieldTypePhone,
		MaxLength: 4,
		Required:  false,
	})

	// telefone1 VARCHAR(8)
	n.RegisterField(FieldMetadata{
		Name:      "telefone1",
		Type:      FieldTypePhone,
		MaxLength: 8,
		Required:  false,
	})

	// ddd2 VARCHAR(4)
	n.RegisterField(FieldMetadata{
		Name:      "ddd2",
		Type:      FieldTypePhone,
		MaxLength: 4,
		Required:  false,
	})

	// telefone2 VARCHAR(8)
	n.RegisterField(FieldMetadata{
		Name:      "telefone2",
		Type:      FieldTypePhone,
		MaxLength: 8,
		Required:  false,
	})

	// ddd_fax VARCHAR(4)
	n.RegisterField(FieldMetadata{
		Name:      "ddd_fax",
		Type:      FieldTypePhone,
		MaxLength: 4,
		Required:  false,
	})

	// fax VARCHAR(8)
	n.RegisterField(FieldMetadata{
		Name:      "fax",
		Type:      FieldTypePhone,
		MaxLength: 8,
		Required:  false,
	})

	// correio_eletronico TEXT
	n.RegisterField(FieldMetadata{
		Name:     "correio_eletronico",
		Type:     FieldTypeEmail,
		Required: false,
	})

	// situacao_especial TEXT
	n.RegisterField(FieldMetadata{
		Name:     "situacao_especial",
		Type:     FieldTypeText,
		Required: false,
	})

	// data_situacao_especial DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_situacao_especial",
		Type:     FieldTypeDate,
		Required: false,
	})

	return n
}

// GetSociosNormalizer retorna o normalizador para a tabela socios
func GetSociosNormalizer() *Normalizer {
	n := NewNormalizer()

	// cnpj VARCHAR(14) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj",
		Type:      FieldTypeCNPJ,
		MaxLength: 14,
		Required:  true,
	})

	// cnpj_basico VARCHAR(8) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_basico",
		Type:      FieldTypeCode,
		MaxLength: 8,
		Required:  true,
	})

	// identificador_de_socio VARCHAR(1) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:        "identificador_de_socio",
		Type:        FieldTypeCode,
		MaxLength:   1,
		Required:    true,
		ValidValues: []string{"1", "2", "3"}, // 1=PF, 2=PJ, 3=Estrangeiro
	})

	// nome_socio TEXT NOT NULL
	n.RegisterField(FieldMetadata{
		Name:     "nome_socio",
		Type:     FieldTypeText,
		Required: true,
	})

	// cnpj_cpf_socio VARCHAR(14) NOT NULL
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_cpf_socio",
		Type:      FieldTypeCode, // Pode ser CPF ou CNPJ
		MaxLength: 14,
		Required:  true,
	})

	// qualificacao_socio VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "qualificacao_socio",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// data_entrada_sociedade DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_entrada_sociedade",
		Type:     FieldTypeDate,
		Required: false,
	})

	// pais VARCHAR(3)
	n.RegisterField(FieldMetadata{
		Name:      "pais",
		Type:      FieldTypeCode,
		MaxLength: 3,
		Required:  false,
	})

	// representante_legal VARCHAR(11)
	n.RegisterField(FieldMetadata{
		Name:      "representante_legal",
		Type:      FieldTypeCPF,
		MaxLength: 11,
		Required:  false,
	})

	// nome_representante TEXT
	n.RegisterField(FieldMetadata{
		Name:     "nome_representante",
		Type:     FieldTypeText,
		Required: false,
	})

	// qualificacao_representante_legal VARCHAR(2)
	n.RegisterField(FieldMetadata{
		Name:      "qualificacao_representante_legal",
		Type:      FieldTypeCode,
		MaxLength: 2,
		Required:  false,
	})

	// faixa_etaria VARCHAR(1)
	n.RegisterField(FieldMetadata{
		Name:      "faixa_etaria",
		Type:      FieldTypeCode,
		MaxLength: 1,
		Required:  false,
	})

	return n
}

// GetSimplesNormalizer retorna o normalizador para a tabela simples
func GetSimplesNormalizer() *Normalizer {
	n := NewNormalizer()

	// cnpj_basico VARCHAR(8) PRIMARY KEY
	n.RegisterField(FieldMetadata{
		Name:      "cnpj_basico",
		Type:      FieldTypeCode,
		MaxLength: 8,
		Required:  true,
	})

	// opcao_simples VARCHAR(1)
	n.RegisterField(FieldMetadata{
		Name:        "opcao_simples",
		Type:        FieldTypeCode,
		MaxLength:   1,
		Required:    false,
		ValidValues: []string{"S", "N"},
	})

	// data_opcao_simples DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_opcao_simples",
		Type:     FieldTypeDate,
		Required: false,
	})

	// data_exclusao_simples DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_exclusao_simples",
		Type:     FieldTypeDate,
		Required: false,
	})

	// opcao_mei VARCHAR(1)
	n.RegisterField(FieldMetadata{
		Name:        "opcao_mei",
		Type:        FieldTypeCode,
		MaxLength:   1,
		Required:    false,
		ValidValues: []string{"S", "N"},
	})

	// data_opcao_mei DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_opcao_mei",
		Type:     FieldTypeDate,
		Required: false,
	})

	// data_exclusao_mei DATE
	n.RegisterField(FieldMetadata{
		Name:     "data_exclusao_mei",
		Type:     FieldTypeDate,
		Required: false,
	})

	return n
}
