package crossdata

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// CrossDataEngine motor de cruzamento de dados
type CrossDataEngine struct {
	cnpjDB string
	redeDB string
}

// NewCrossDataEngine cria novo motor
func NewCrossDataEngine(cnpjDB, redeDB string) *CrossDataEngine {
	return &CrossDataEngine{
		cnpjDB: cnpjDB,
		redeDB: redeDB,
	}
}

// Empresa dados completos de empresa SEM CENSURA
type Empresa struct {
	CNPJ                    string  `json:"cnpj"`
	CNPJBasico              string  `json:"cnpj_basico"`
	RazaoSocial             string  `json:"razao_social"`
	NomeFantasia            string  `json:"nome_fantasia"`
	MatrizFilial            string  `json:"matriz_filial"`
	SituacaoCadastral       string  `json:"situacao_cadastral"`
	DataSituacaoCadastral   string  `json:"data_situacao_cadastral"`
	MotivoSituacaoCadastral string  `json:"motivo_situacao_cadastral"`
	DataInicioAtividades    string  `json:"data_inicio_atividades"`
	CNAEFiscal              string  `json:"cnae_fiscal"`
	CNAESecundaria          string  `json:"cnae_secundaria"`
	NaturezaJuridica        string  `json:"natureza_juridica"`
	CapitalSocial           float64 `json:"capital_social"`
	PorteEmpresa            string  `json:"porte_empresa"`
	// ENDEREÇO COMPLETO - SEM CENSURA
	TipoLogradouro  string `json:"tipo_logradouro"`
	Logradouro      string `json:"logradouro"`
	Numero          string `json:"numero"`
	Complemento     string `json:"complemento"`
	Bairro          string `json:"bairro"`
	CEP             string `json:"cep"`
	UF              string `json:"uf"`
	Municipio       string `json:"municipio"`
	// CONTATOS COMPLETOS - SEM CENSURA
	DDD1               string `json:"ddd1"`
	Telefone1          string `json:"telefone1"`
	DDD2               string `json:"ddd2"`
	Telefone2          string `json:"telefone2"`
	DDDFax             string `json:"ddd_fax"`
	Fax                string `json:"fax"`
	CorreioEletronico  string `json:"correio_eletronico"`
	OpcaoMEI           string `json:"opcao_mei"`
}

// Socio dados completos de sócio SEM CENSURA
type Socio struct {
	CNPJ                         string `json:"cnpj"`
	CNPJBasico                   string `json:"cnpj_basico"`
	IdentificadorSocio           string `json:"identificador_socio"` // 1=PF, 2=PJ, 3=Estrangeiro
	// DADOS PESSOAIS COMPLETOS - SEM CENSURA
	NomeSocio                    string `json:"nome_socio"`
	CNPJCPFSocio                 string `json:"cnpj_cpf_socio"` // CPF/CNPJ SEM MÁSCARA
	QualificacaoSocio            string `json:"qualificacao_socio"`
	DataEntradaSociedade         string `json:"data_entrada_sociedade"`
	Pais                         string `json:"pais"`
	// REPRESENTANTE LEGAL - SEM CENSURA
	RepresentanteLegal           string `json:"representante_legal"` // CPF do representante
	NomeRepresentante            string `json:"nome_representante"`
	QualificacaoRepresentante    string `json:"qualificacao_representante"`
	FaixaEtaria                  string `json:"faixa_etaria"`
}

// 1. CPF → Empresas
func (c *CrossDataEngine) EmpresasPorCPF(cpf string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT DISTINCT 
			s.cnpj_cpf_socio as cpf,
			s.nome_socio,
			est.cnpj,
			e.razao_social,
			est.nome_fantasia,
			s.qualificacao_socio,
			s.data_entrada_sociedade,
			est.situacao_cadastral,
			e.capital_social,
			est.correio_eletronico,
			est.telefone1,
			est.ddd1,
			est.logradouro,
			est.numero,
			est.bairro,
			est.cep,
			est.uf
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.cnpj_cpf_socio = ?
		ORDER BY s.data_entrada_sociedade DESC
	`

	rows, err := db.Query(query, cpf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 2. CNPJ → Sócios (TODOS OS DADOS SEM CENSURA)
func (c *CrossDataEngine) SociosPorCNPJ(cnpj string) ([]Socio, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.cnpj,
			s.cnpj_basico,
			s.identificador_de_socio,
			s.nome_socio,
			s.cnpj_cpf_socio,
			s.qualificacao_socio,
			s.data_entrada_sociedade,
			s.pais,
			s.representante_legal,
			s.nome_representante,
			s.qualificacao_representante_legal,
			s.faixa_etaria
		FROM socios s
		WHERE s.cnpj = ?
		ORDER BY s.qualificacao_socio, s.nome_socio
	`

	rows, err := db.Query(query, cnpj)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var socios []Socio
	for rows.Next() {
		var s Socio
		err := rows.Scan(
			&s.CNPJ, &s.CNPJBasico, &s.IdentificadorSocio,
			&s.NomeSocio, &s.CNPJCPFSocio, &s.QualificacaoSocio,
			&s.DataEntradaSociedade, &s.Pais,
			&s.RepresentanteLegal, &s.NomeRepresentante,
			&s.QualificacaoRepresentante, &s.FaixaEtaria,
		)
		if err != nil {
			continue
		}
		socios = append(socios, s)
	}

	return socios, nil
}

// 3. Sócios em Comum entre duas empresas
func (c *CrossDataEngine) SociosEmComum(cnpj1, cnpj2 string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT DISTINCT
			s1.cnpj as cnpj1,
			s2.cnpj as cnpj2,
			s1.nome_socio,
			s1.cnpj_cpf_socio,
			s1.qualificacao_socio as qualif_empresa1,
			s2.qualificacao_socio as qualif_empresa2,
			s1.data_entrada_sociedade as data_entrada1,
			s2.data_entrada_sociedade as data_entrada2
		FROM socios s1
		JOIN socios s2 ON s1.cnpj_cpf_socio = s2.cnpj_cpf_socio
		WHERE s1.cnpj = ? AND s2.cnpj = ? AND s1.cnpj != s2.cnpj
	`

	rows, err := db.Query(query, cnpj1, cnpj2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 4. Rede de Empresas de uma Pessoa (2º grau)
func (c *CrossDataEngine) RedeEmpresasPessoa(cpf string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		WITH empresas_pessoa AS (
			SELECT DISTINCT cnpj, qualificacao_socio
			FROM socios
			WHERE cnpj_cpf_socio = ?
		)
		SELECT 
			ep.cnpj,
			e.razao_social,
			est.nome_fantasia,
			ep.qualificacao_socio,
			s2.nome_socio as outros_socios,
			s2.cnpj_cpf_socio as cpf_outros_socios,
			s2.qualificacao_socio as qualif_outros_socios,
			est.situacao_cadastral,
			est.correio_eletronico,
			est.telefone1
		FROM empresas_pessoa ep
		JOIN estabelecimento est ON ep.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		LEFT JOIN socios s2 ON ep.cnpj = s2.cnpj AND s2.cnpj_cpf_socio != ?
		ORDER BY ep.cnpj, s2.nome_socio
	`

	rows, err := db.Query(query, cpf, cpf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 5. Empresas no Mesmo Endereço
func (c *CrossDataEngine) EmpresasMesmoEndereco(cep, logradouro, numero string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			est.cnpj,
			e.razao_social,
			est.nome_fantasia,
			est.situacao_cadastral,
			est.correio_eletronico,
			est.telefone1,
			est.ddd1,
			est.logradouro,
			est.numero,
			est.complemento,
			est.bairro,
			est.cep,
			est.uf
		FROM estabelecimento est
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE est.cep = ? 
		  AND est.logradouro = ?
		  AND est.numero = ?
		ORDER BY e.razao_social
	`

	rows, err := db.Query(query, cep, logradouro, numero)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 6. Empresas com Mesmo Email ou Telefone
func (c *CrossDataEngine) EmpresasMesmoContato(email, telefone string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			est.cnpj,
			e.razao_social,
			est.nome_fantasia,
			est.correio_eletronico,
			est.telefone1,
			est.ddd1,
			est.telefone2,
			est.ddd2,
			est.situacao_cadastral,
			est.logradouro,
			est.numero,
			est.bairro,
			est.cep,
			est.uf
		FROM estabelecimento est
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE est.correio_eletronico = ? OR est.telefone1 = ?
		ORDER BY e.razao_social
	`

	rows, err := db.Query(query, email, telefone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 7. Representantes Legais (Menores com Representantes)
func (c *CrossDataEngine) RepresentantesLegais() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.cnpj,
			e.razao_social,
			est.nome_fantasia,
			s.nome_socio as socio_menor,
			s.cnpj_cpf_socio as cpf_menor,
			s.faixa_etaria,
			s.representante_legal as cpf_representante,
			s.nome_representante,
			s.qualificacao_representante_legal,
			est.situacao_cadastral
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.representante_legal IS NOT NULL AND s.representante_legal != ''
		ORDER BY s.nome_socio
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 8. Empresas Estrangeiras
func (c *CrossDataEngine) EmpresasEstrangeiras() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			est.cnpj,
			e.razao_social,
			est.nome_fantasia,
			est.nome_cidade_exterior,
			p.descricao as pais,
			est.correio_eletronico,
			est.telefone1,
			est.situacao_cadastral
		FROM estabelecimento est
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		LEFT JOIN pais p ON est.pais = p.codigo
		WHERE est.uf = 'EX'
		ORDER BY p.descricao, e.razao_social
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 9. Sócios Estrangeiros
func (c *CrossDataEngine) SociosEstrangeiros() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.cnpj,
			e.razao_social,
			est.nome_fantasia,
			s.nome_socio,
			s.cnpj_cpf_socio,
			p.descricao as pais,
			s.qualificacao_socio,
			s.data_entrada_sociedade,
			est.situacao_cadastral
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		LEFT JOIN pais p ON s.pais = p.codigo
		WHERE s.identificador_de_socio = '3'
		ORDER BY p.descricao, s.nome_socio
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 10. Timeline de Atividades de uma Pessoa
func (c *CrossDataEngine) TimelinePessoa(cpf string) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.cnpj_cpf_socio,
			s.nome_socio,
			s.cnpj,
			e.razao_social,
			est.nome_fantasia,
			s.qualificacao_socio,
			s.data_entrada_sociedade,
			est.data_inicio_atividades,
			est.situacao_cadastral,
			est.data_situacao_cadastral,
			est.correio_eletronico,
			est.telefone1
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.cnpj_cpf_socio = ?
		ORDER BY s.data_entrada_sociedade, est.data_inicio_atividades
	`

	rows, err := db.Query(query, cpf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 11. Empresas Baixadas com Sócios Ativos
func (c *CrossDataEngine) SociosEmpresasBaixadas() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.cnpj_cpf_socio,
			s.nome_socio,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '02' THEN s.cnpj END) as empresas_ativas,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '08' THEN s.cnpj END) as empresas_baixadas,
			COUNT(DISTINCT s.cnpj) as total_empresas
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		GROUP BY s.cnpj_cpf_socio, s.nome_socio
		HAVING empresas_baixadas > 0 AND empresas_ativas > 0
		ORDER BY empresas_baixadas DESC, empresas_ativas DESC
		LIMIT 1000
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMaps(rows)
}

// 12. Dados Completos de Empresa (SEM CENSURA)
func (c *CrossDataEngine) DadosCompletosEmpresa(cnpj string) (*Empresa, error) {
	db, err := sql.Open("sqlite3", c.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			est.cnpj,
			est.cnpj_basico,
			e.razao_social,
			est.nome_fantasia,
			est.matriz_filial,
			est.situacao_cadastral,
			est.data_situacao_cadastral,
			est.motivo_situacao_cadastral,
			est.data_inicio_atividades,
			est.cnae_fiscal,
			est.cnae_fiscal_secundaria,
			e.natureza_juridica,
			e.capital_social,
			e.porte_empresa,
			est.tipo_logradouro,
			est.logradouro,
			est.numero,
			est.complemento,
			est.bairro,
			est.cep,
			est.uf,
			est.municipio,
			est.ddd1,
			est.telefone1,
			est.ddd2,
			est.telefone2,
			est.ddd_fax,
			est.fax,
			est.correio_eletronico,
			COALESCE(sim.opcao_mei, '') as opcao_mei
		FROM estabelecimento est
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		LEFT JOIN simples sim ON est.cnpj_basico = sim.cnpj_basico
		WHERE est.cnpj = ?
	`

	var emp Empresa
	err = db.QueryRow(query, cnpj).Scan(
		&emp.CNPJ, &emp.CNPJBasico, &emp.RazaoSocial, &emp.NomeFantasia,
		&emp.MatrizFilial, &emp.SituacaoCadastral, &emp.DataSituacaoCadastral,
		&emp.MotivoSituacaoCadastral, &emp.DataInicioAtividades,
		&emp.CNAEFiscal, &emp.CNAESecundaria, &emp.NaturezaJuridica,
		&emp.CapitalSocial, &emp.PorteEmpresa,
		&emp.TipoLogradouro, &emp.Logradouro, &emp.Numero, &emp.Complemento,
		&emp.Bairro, &emp.CEP, &emp.UF, &emp.Municipio,
		&emp.DDD1, &emp.Telefone1, &emp.DDD2, &emp.Telefone2,
		&emp.DDDFax, &emp.Fax, &emp.CorreioEletronico, &emp.OpcaoMEI,
	)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// scanToMaps converte rows para slice de maps
func scanToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	return results, rows.Err()
}
