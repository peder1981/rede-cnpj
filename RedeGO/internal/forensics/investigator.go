package forensics

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Investigator motor de investigação forense
type Investigator struct {
	cnpjDB string
	redeDB string
}

// NewInvestigator cria novo investigador
func NewInvestigator(cnpjDB, redeDB string) *Investigator {
	return &Investigator{
		cnpjDB: cnpjDB,
		redeDB: redeDB,
	}
}

// SuspectProfile perfil de suspeito
type SuspectProfile struct {
	CPF                  string                   `json:"cpf"`
	Nome                 string                   `json:"nome"`
	TotalEmpresas        int                      `json:"total_empresas"`
	EmpresasAtivas       int                      `json:"empresas_ativas"`
	EmpresasBaixadas     int                      `json:"empresas_baixadas"`
	EmpresasSuspensas    int                      `json:"empresas_suspensas"`
	CapitalSocialTotal   float64                  `json:"capital_social_total"`
	EnderecosDiferentes  int                      `json:"enderecos_diferentes"`
	TelefonesDiferentes  int                      `json:"telefones_diferentes"`
	EmailsDiferentes     int                      `json:"emails_diferentes"`
	PrimeiraEmpresa      string                   `json:"primeira_empresa"`
	UltimaEmpresa        string                   `json:"ultima_empresa"`
	PeriodoAtividade     string                   `json:"periodo_atividade"`
	RedeBancaria         int                      `json:"rede_bancaria"` // Empresas de outros sócios
	Score                int                      `json:"score_risco"`   // 0-100
	Flags                []string                 `json:"flags"`
	Empresas             []map[string]interface{} `json:"empresas"`
}

// CompanyCluster cluster de empresas suspeitas
type CompanyCluster struct {
	TipoCluster     string                   `json:"tipo_cluster"`
	Criterio        string                   `json:"criterio"`
	ValorComum      string                   `json:"valor_comum"`
	TotalEmpresas   int                      `json:"total_empresas"`
	TotalSocios     int                      `json:"total_socios"`
	Empresas        []map[string]interface{} `json:"empresas"`
	Socios          []map[string]interface{} `json:"socios"`
	Score           int                      `json:"score_risco"`
	Flags           []string                 `json:"flags"`
}

// 1. PERFIL COMPLETO DE SUSPEITO
func (inv *Investigator) InvestigatePerson(cpf string) (*SuspectProfile, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	profile := &SuspectProfile{
		CPF:   cpf,
		Flags: []string{},
	}

	// Query principal
	query := `
		SELECT 
			s.nome_socio,
			COUNT(DISTINCT s.cnpj) as total_empresas,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '02' THEN s.cnpj END) as ativas,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '08' THEN s.cnpj END) as baixadas,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral IN ('03','04') THEN s.cnpj END) as suspensas,
			SUM(DISTINCT e.capital_social) as capital_total,
			COUNT(DISTINCT est.cep || est.logradouro) as enderecos,
			COUNT(DISTINCT est.telefone1) as telefones,
			COUNT(DISTINCT est.correio_eletronico) as emails,
			MIN(s.data_entrada_sociedade) as primeira,
			MAX(s.data_entrada_sociedade) as ultima
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.cnpj_cpf_socio = ?
		GROUP BY s.nome_socio
	`

	var capitalTotal sql.NullFloat64
	var primeira, ultima sql.NullString
	
	err = db.QueryRow(query, cpf).Scan(
		&profile.Nome,
		&profile.TotalEmpresas,
		&profile.EmpresasAtivas,
		&profile.EmpresasBaixadas,
		&profile.EmpresasSuspensas,
		&capitalTotal,
		&profile.EnderecosDiferentes,
		&profile.TelefonesDiferentes,
		&profile.EmailsDiferentes,
		&primeira,
		&ultima,
	)

	if err != nil {
		return nil, err
	}

	if capitalTotal.Valid {
		profile.CapitalSocialTotal = capitalTotal.Float64
	}

	if primeira.Valid && ultima.Valid {
		profile.PrimeiraEmpresa = primeira.String
		profile.UltimaEmpresa = ultima.String
		
		// Calcula período
		if len(primeira.String) == 8 && len(ultima.String) == 8 {
			p, _ := time.Parse("20060102", primeira.String)
			u, _ := time.Parse("20060102", ultima.String)
			anos := u.Year() - p.Year()
			profile.PeriodoAtividade = fmt.Sprintf("%d anos", anos)
		}
	}

	// Busca empresas detalhadas
	empresasQuery := `
		SELECT 
			s.cnpj,
			e.razao_social,
			est.nome_fantasia,
			s.qualificacao_socio,
			s.data_entrada_sociedade,
			est.situacao_cadastral,
			e.capital_social,
			est.correio_eletronico,
			est.telefone1,
			est.cep,
			est.logradouro,
			est.numero,
			est.uf
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.cnpj_cpf_socio = ?
		ORDER BY s.data_entrada_sociedade DESC
	`

	rows, err := db.Query(empresasQuery, cpf)
	if err != nil {
		return profile, nil
	}
	defer rows.Close()

	for rows.Next() {
		var emp map[string]interface{} = make(map[string]interface{})
		var cnpj, razao, fantasia, qualif, data, situacao, email, tel, cep, logr, num, uf string
		var capital float64

		rows.Scan(&cnpj, &razao, &fantasia, &qualif, &data, &situacao, &capital,
			&email, &tel, &cep, &logr, &num, &uf)

		emp["cnpj"] = cnpj
		emp["razao_social"] = razao
		emp["nome_fantasia"] = fantasia
		emp["qualificacao"] = qualif
		emp["data_entrada"] = data
		emp["situacao"] = situacao
		emp["capital_social"] = capital
		emp["email"] = email
		emp["telefone"] = tel
		emp["endereco"] = fmt.Sprintf("%s, %s - %s - %s", logr, num, cep, uf)

		profile.Empresas = append(profile.Empresas, emp)
	}

	// Calcula rede bancária (empresas de outros sócios)
	redeQuery := `
		SELECT COUNT(DISTINCT s2.cnpj)
		FROM socios s1
		JOIN socios s2 ON s1.cnpj = s2.cnpj
		WHERE s1.cnpj_cpf_socio = ? AND s2.cnpj_cpf_socio != ?
	`
	db.QueryRow(redeQuery, cpf, cpf).Scan(&profile.RedeBancaria)

	// Calcula score e flags
	profile.calculateRiskScore()

	return profile, nil
}

// calculateRiskScore calcula score de risco
func (p *SuspectProfile) calculateRiskScore() {
	score := 0
	
	// Muitas empresas baixadas
	if p.EmpresasBaixadas > 5 {
		score += 20
		p.Flags = append(p.Flags, fmt.Sprintf("ALTO: %d empresas baixadas", p.EmpresasBaixadas))
	} else if p.EmpresasBaixadas > 2 {
		score += 10
		p.Flags = append(p.Flags, fmt.Sprintf("MÉDIO: %d empresas baixadas", p.EmpresasBaixadas))
	}

	// Muitas empresas ativas
	if p.EmpresasAtivas > 10 {
		score += 15
		p.Flags = append(p.Flags, fmt.Sprintf("ALTO: %d empresas ativas simultaneamente", p.EmpresasAtivas))
	}

	// Empresas suspensas
	if p.EmpresasSuspensas > 0 {
		score += 15
		p.Flags = append(p.Flags, fmt.Sprintf("CRÍTICO: %d empresas suspensas", p.EmpresasSuspensas))
	}

	// Muitos endereços diferentes
	if p.EnderecosDiferentes > 10 {
		score += 10
		p.Flags = append(p.Flags, fmt.Sprintf("MÉDIO: %d endereços diferentes", p.EnderecosDiferentes))
	}

	// Muitos telefones diferentes
	if p.TelefonesDiferentes > 5 {
		score += 5
		p.Flags = append(p.Flags, fmt.Sprintf("BAIXO: %d telefones diferentes", p.TelefonesDiferentes))
	}

	// Rede bancária grande
	if p.RedeBancaria > 50 {
		score += 20
		p.Flags = append(p.Flags, fmt.Sprintf("ALTO: Rede de %d empresas conectadas", p.RedeBancaria))
	}

	// Capital social muito alto
	if p.CapitalSocialTotal > 10000000 { // 10 milhões
		score += 10
		p.Flags = append(p.Flags, fmt.Sprintf("INFO: Capital social total R$ %.2f milhões", p.CapitalSocialTotal/1000000))
	}

	if score > 100 {
		score = 100
	}

	p.Score = score
}

// 2. DETECTAR EMPRESAS DE FACHADA (MESMO ENDEREÇO)
func (inv *Investigator) DetectShellCompanies(minEmpresas int) ([]CompanyCluster, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			est.cep,
			est.logradouro,
			est.numero,
			est.uf,
			COUNT(DISTINCT est.cnpj) as total_empresas,
			COUNT(DISTINCT s.cnpj_cpf_socio) as total_socios
		FROM estabelecimento est
		LEFT JOIN socios s ON est.cnpj = s.cnpj
		WHERE est.situacao_cadastral = '02'
		  AND est.cep IS NOT NULL
		  AND est.logradouro IS NOT NULL
		GROUP BY est.cep, est.logradouro, est.numero, est.uf
		HAVING total_empresas >= ?
		ORDER BY total_empresas DESC
		LIMIT 100
	`

	rows, err := db.Query(query, minEmpresas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []CompanyCluster

	for rows.Next() {
		var cep, logr, num, uf string
		var totalEmp, totalSoc int

		rows.Scan(&cep, &logr, &num, &uf, &totalEmp, &totalSoc)

		cluster := CompanyCluster{
			TipoCluster:   "MESMO_ENDERECO",
			Criterio:      "Empresas no mesmo endereço físico",
			ValorComum:    fmt.Sprintf("%s, %s - %s - %s", logr, num, cep, uf),
			TotalEmpresas: totalEmp,
			TotalSocios:   totalSoc,
			Flags:         []string{},
		}

		// Score baseado na quantidade
		if totalEmp > 50 {
			cluster.Score = 90
			cluster.Flags = append(cluster.Flags, "CRÍTICO: Mais de 50 empresas no mesmo endereço")
		} else if totalEmp > 20 {
			cluster.Score = 70
			cluster.Flags = append(cluster.Flags, "ALTO: Mais de 20 empresas no mesmo endereço")
		} else if totalEmp > 10 {
			cluster.Score = 50
			cluster.Flags = append(cluster.Flags, "MÉDIO: Mais de 10 empresas no mesmo endereço")
		}

		// Busca empresas do cluster
		empQuery := `
			SELECT cnpj, e.razao_social, est.nome_fantasia, est.correio_eletronico, est.telefone1
			FROM estabelecimento est
			JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
			WHERE est.cep = ? AND est.logradouro = ? AND est.numero = ?
			  AND est.situacao_cadastral = '02'
			LIMIT 50
		`

		empRows, _ := db.Query(empQuery, cep, logr, num)
		for empRows.Next() {
			var cnpj, razao, fantasia, email, tel string
			empRows.Scan(&cnpj, &razao, &fantasia, &email, &tel)
			cluster.Empresas = append(cluster.Empresas, map[string]interface{}{
				"cnpj":          cnpj,
				"razao_social":  razao,
				"nome_fantasia": fantasia,
				"email":         email,
				"telefone":      tel,
			})
		}
		empRows.Close()

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// 3. DETECTAR LARANJAS (MESMO TELEFONE/EMAIL)
func (inv *Investigator) DetectFrontmen(criterio string, valor string) (*CompanyCluster, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string

	if criterio == "telefone" {
		query = `
			SELECT 
				est.cnpj,
				e.razao_social,
				est.nome_fantasia,
				est.correio_eletronico,
				est.telefone1,
				est.situacao_cadastral
			FROM estabelecimento est
			JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
			WHERE est.telefone1 = ?
		`
	} else {
		query = `
			SELECT 
				est.cnpj,
				e.razao_social,
				est.nome_fantasia,
				est.correio_eletronico,
				est.telefone1,
				est.situacao_cadastral
			FROM estabelecimento est
			JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
			WHERE est.correio_eletronico = ?
		`
	}

	rows, err := db.Query(query, valor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cluster := &CompanyCluster{
		TipoCluster: strings.ToUpper(criterio),
		Criterio:    fmt.Sprintf("Empresas com mesmo %s", criterio),
		ValorComum:  valor,
		Flags:       []string{},
	}

	for rows.Next() {
		var cnpj, razao, fantasia, email, tel, situacao string
		rows.Scan(&cnpj, &razao, &fantasia, &email, &tel, &situacao)

		cluster.Empresas = append(cluster.Empresas, map[string]interface{}{
			"cnpj":          cnpj,
			"razao_social":  razao,
			"nome_fantasia": fantasia,
			"email":         email,
			"telefone":      tel,
			"situacao":      situacao,
		})
	}

	cluster.TotalEmpresas = len(cluster.Empresas)

	// Score
	if cluster.TotalEmpresas > 10 {
		cluster.Score = 90
		cluster.Flags = append(cluster.Flags, fmt.Sprintf("CRÍTICO: %d empresas com mesmo %s", cluster.TotalEmpresas, criterio))
	} else if cluster.TotalEmpresas > 5 {
		cluster.Score = 70
		cluster.Flags = append(cluster.Flags, fmt.Sprintf("ALTO: %d empresas com mesmo %s", cluster.TotalEmpresas, criterio))
	} else if cluster.TotalEmpresas > 2 {
		cluster.Score = 50
		cluster.Flags = append(cluster.Flags, fmt.Sprintf("MÉDIO: %d empresas com mesmo %s", cluster.TotalEmpresas, criterio))
	}

	return cluster, nil
}

// 4. ANÁLISE TEMPORAL (EMPRESAS ABERTAS EM MASSA)
func (inv *Investigator) DetectMassRegistration(cpf string, diasJanela int) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			s.data_entrada_sociedade,
			COUNT(*) as total,
			GROUP_CONCAT(est.cnpj, ', ') as cnpjs,
			GROUP_CONCAT(e.razao_social, ' | ') as empresas
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
		WHERE s.cnpj_cpf_socio = ?
		GROUP BY s.data_entrada_sociedade
		HAVING total >= 2
		ORDER BY total DESC
	`

	rows, err := db.Query(query, cpf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var data, cnpjs, empresas string
		var total int

		rows.Scan(&data, &total, &cnpjs, &empresas)

		results = append(results, map[string]interface{}{
			"data":      data,
			"total":     total,
			"cnpjs":     cnpjs,
			"empresas":  empresas,
			"flag":      fmt.Sprintf("SUSPEITO: %d empresas abertas na mesma data", total),
		})
	}

	return results, nil
}

// 5. CADEIA DE CONTROLE (EMPRESAS DE EMPRESAS)
func (inv *Investigator) TraceOwnershipChain(cnpj string, maxNivel int) ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Busca sócios PJ
	query := `
		WITH RECURSIVE cadeia(cnpj, cnpj_socio, nome_socio, nivel) AS (
			SELECT s.cnpj, s.cnpj_cpf_socio, s.nome_socio, 1
			FROM socios s
			WHERE s.cnpj = ? AND s.identificador_de_socio = '2'
			
			UNION ALL
			
			SELECT s.cnpj, s.cnpj_cpf_socio, s.nome_socio, c.nivel + 1
			FROM socios s
			JOIN cadeia c ON s.cnpj = c.cnpj_socio
			WHERE s.identificador_de_socio = '2' AND c.nivel < ?
		)
		SELECT DISTINCT cnpj, cnpj_socio, nome_socio, nivel
		FROM cadeia
		ORDER BY nivel, cnpj
	`

	rows, err := db.Query(query, cnpj, maxNivel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var cnpjEmp, cnpjSocio, nome string
		var nivel int

		rows.Scan(&cnpjEmp, &cnpjSocio, &nome, &nivel)

		results = append(results, map[string]interface{}{
			"cnpj_empresa":     cnpjEmp,
			"cnpj_socio":       cnpjSocio,
			"nome_socio":       nome,
			"nivel":            nivel,
			"tipo":             "PESSOA_JURIDICA",
		})
	}

	return results, nil
}

// 6. PADRÃO DE ATIVIDADE SUSPEITA
func (inv *Investigator) DetectSuspiciousPatterns() ([]map[string]interface{}, error) {
	db, err := sql.Open("sqlite3", inv.cnpjDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Pessoas com muitas empresas baixadas rapidamente
	query := `
		SELECT 
			s.cnpj_cpf_socio,
			s.nome_socio,
			COUNT(DISTINCT s.cnpj) as total_empresas,
			COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '08' THEN s.cnpj END) as baixadas,
			MIN(est.data_situacao_cadastral) as primeira_baixa,
			MAX(est.data_situacao_cadastral) as ultima_baixa
		FROM socios s
		JOIN estabelecimento est ON s.cnpj = est.cnpj
		WHERE est.situacao_cadastral = '08'
		GROUP BY s.cnpj_cpf_socio, s.nome_socio
		HAVING baixadas >= 5
		ORDER BY baixadas DESC
		LIMIT 100
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var cpf, nome, primeira, ultima string
		var total, baixadas int

		rows.Scan(&cpf, &nome, &total, &baixadas, &primeira, &ultima)

		score := 0
		if baixadas > 10 {
			score = 90
		} else if baixadas > 7 {
			score = 70
		} else {
			score = 50
		}

		results = append(results, map[string]interface{}{
			"cpf":            cpf,
			"nome":           nome,
			"total_empresas": total,
			"baixadas":       baixadas,
			"primeira_baixa": primeira,
			"ultima_baixa":   ultima,
			"score":          score,
			"flag":           fmt.Sprintf("ALTO RISCO: %d empresas baixadas", baixadas),
		})
	}

	return results, nil
}
