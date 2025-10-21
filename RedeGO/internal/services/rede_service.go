package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/database"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	"github.com/peder1981/rede-cnpj/RedeGO/pkg/cpfcnpj"
)

// RedeService gerencia operações relacionadas à rede de relacionamentos
type RedeService struct {
	cfg *config.Config
	db  *sql.DB
}

// NewRedeService cria uma nova instância do serviço
func NewRedeService(cfg *config.Config) *RedeService {
	return &RedeService{
		cfg: cfg,
		db:  database.GetDBRede(),
	}
}

// CamadasRede busca camadas de relacionamentos a partir de uma lista de IDs
func (s *RedeService) CamadasRede(camada int, listaIDs []string, grupo string, criterioCaminhos string) (*models.Graph, error) {
	database.Lock()
	defer database.Unlock()

	startTime := time.Now()
	maxDuration := time.Duration(s.cfg.TempoMaximoConsulta) * time.Second

	graph := &models.Graph{
		Nodes: make([]models.Node, 0),
		Edges: make([]models.Edge, 0),
	}

	if len(listaIDs) == 0 && grupo == "" {
		return graph, nil
	}

	// Mapa para evitar duplicatas
	nodeMap := make(map[string]bool)
	edgeMap := make(map[string]bool)

	// Processa cada ID da lista
	for _, id := range listaIDs {
		if time.Since(startTime) > maxDuration {
			break
		}

		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		// Valida e normaliza CPF/CNPJ
		if validCNPJ := cpfcnpj.ValidarCNPJ(id); validCNPJ != "" {
			id = validCNPJ
		} else if validCPF := cpfcnpj.ValidarCPF(id); validCPF != "" {
			id = validCPF
		}

		// Adiciona nó inicial se não existir
		if !nodeMap[id] {
			node := s.createNodeFromID(id)
			graph.Nodes = append(graph.Nodes, node)
			nodeMap[id] = true
		}

		// Busca relacionamentos
		if err := s.buscarRelacionamentos(id, camada, graph, nodeMap, edgeMap, startTime, maxDuration); err != nil {
			return nil, err
		}
	}

	return graph, nil
}

// buscarRelacionamentos busca recursivamente os relacionamentos de um ID
func (s *RedeService) buscarRelacionamentos(id string, camada int, graph *models.Graph, nodeMap map[string]bool, edgeMap map[string]bool, startTime time.Time, maxDuration time.Duration) error {
	if camada <= 0 || time.Since(startTime) > maxDuration {
		return nil
	}

	// Busca sócios se for CNPJ
	if len(id) == 14 {
		if err := s.buscarSocios(id, graph, nodeMap, edgeMap); err != nil {
			return err
		}
	}

	// Busca empresas se for CPF ou nome de sócio
	if len(id) == 11 || len(id) > 14 {
		if err := s.buscarEmpresas(id, graph, nodeMap, edgeMap); err != nil {
			return err
		}
	}

	return nil
}

// buscarSocios busca os sócios de uma empresa
func (s *RedeService) buscarSocios(cnpj string, graph *models.Graph, nodeMap map[string]bool, edgeMap map[string]bool) error {
	db := database.GetDBRede()
	if db == nil {
		return fmt.Errorf("banco de dados de rede não disponível")
	}

	query := `
		SELECT identificador_socio, nome_socio, qualificacao_socio, cpf_cnpj_socio
		FROM socios
		WHERE cnpj = ?
		LIMIT ?
	`

	rows, err := db.Query(query, cnpj, s.cfg.LimiteRegistrosCamada)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var identificador, nome, qualificacao string
		var cpfCnpjSocio sql.NullString

		if err := rows.Scan(&identificador, &nome, &qualificacao, &cpfCnpjSocio); err != nil {
			continue
		}

		// Cria nó do sócio se não existir
		socioID := identificador
		if cpfCnpjSocio.Valid && cpfCnpjSocio.String != "" {
			socioID = cpfCnpjSocio.String
		}

		if !nodeMap[socioID] {
			node := models.Node{
				ID:    socioID,
				Label: nome,
				Type:  "PF",
				Icon:  "pessoa",
			}
			graph.Nodes = append(graph.Nodes, node)
			nodeMap[socioID] = true
		}

		// Cria aresta
		edgeKey := cnpj + "->" + socioID
		if !edgeMap[edgeKey] {
			edge := models.Edge{
				From:         cnpj,
				To:           socioID,
				Label:        qualificacao,
				Type:         "socio",
				Qualificacao: qualificacao,
			}
			graph.Edges = append(graph.Edges, edge)
			edgeMap[edgeKey] = true
		}
	}

	return nil
}

// buscarEmpresas busca as empresas de um sócio
func (s *RedeService) buscarEmpresas(socioID string, graph *models.Graph, nodeMap map[string]bool, edgeMap map[string]bool) error {
	db := database.GetDBRede()
	if db == nil {
		return fmt.Errorf("banco de dados de rede não disponível")
	}

	query := `
		SELECT cnpj, qualificacao_socio
		FROM socios
		WHERE identificador_socio = ? OR cpf_cnpj_socio = ?
		LIMIT ?
	`

	rows, err := db.Query(query, socioID, socioID, s.cfg.LimiteRegistrosCamada)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cnpj, qualificacao string

		if err := rows.Scan(&cnpj, &qualificacao); err != nil {
			continue
		}

		// Cria nó da empresa se não existir
		if !nodeMap[cnpj] {
			node := s.createNodeFromID(cnpj)
			graph.Nodes = append(graph.Nodes, node)
			nodeMap[cnpj] = true
		}

		// Cria aresta
		edgeKey := cnpj + "->" + socioID
		if !edgeMap[edgeKey] {
			edge := models.Edge{
				From:         cnpj,
				To:           socioID,
				Label:        qualificacao,
				Type:         "socio",
				Qualificacao: qualificacao,
			}
			graph.Edges = append(graph.Edges, edge)
			edgeMap[edgeKey] = true
		}
	}

	return nil
}

// createNodeFromID cria um nó a partir de um ID
func (s *RedeService) createNodeFromID(id string) models.Node {
	node := models.Node{
		ID:   id,
		Type: "PJ",
		Icon: "empresa",
	}

	// Busca dados da empresa se for CNPJ
	if len(id) == 14 {
		if dados := s.GetDadosCNPJ(id); dados != nil {
			node.Label = dados.RazaoSocial
			if dados.SituacaoCadastral == "02" {
				node.Color = "green"
			} else {
				node.Color = "red"
			}
		} else {
			node.Label = cpfcnpj.CNPJFormatado(id)
		}
	} else {
		node.Label = id
		node.Type = "PF"
		node.Icon = "pessoa"
	}

	return node
}

// GetDadosCNPJ busca dados detalhados de um CNPJ
func (s *RedeService) GetDadosCNPJ(cnpj string) *models.CNPJData {
	db := database.GetDBReceita()
	if db == nil {
		return nil
	}

	query := `
		SELECT 
			cnpj, razao_social, nome_fantasia, situacao_cadastral,
			data_situacao, motivo_situacao, natureza_juridica,
			cnae_principal, capital_social, porte, data_abertura,
			logradouro, numero, complemento, bairro, municipio, uf, cep,
			email, telefone1, telefone2
		FROM empresas
		WHERE cnpj = ?
	`

	var dados models.CNPJData
	var nomeFantasia, dataAbertura, email, telefone1, telefone2 sql.NullString
	var complemento, motivoSituacao, naturezaJuridica, cnae sql.NullString

	err := db.QueryRow(query, cnpj).Scan(
		&dados.CNPJ, &dados.RazaoSocial, &nomeFantasia, &dados.SituacaoCadastral,
		&dados.DataSituacao, &motivoSituacao, &naturezaJuridica,
		&cnae, &dados.CapitalSocial, &dados.Porte, &dataAbertura,
		&dados.Logradouro, &dados.Numero, &complemento, &dados.Bairro,
		&dados.Municipio, &dados.UF, &dados.CEP,
		&email, &telefone1, &telefone2,
	)

	if err != nil {
		return nil
	}

	// Preenche campos opcionais
	if nomeFantasia.Valid {
		dados.NomeFantasia = nomeFantasia.String
	}
	if dataAbertura.Valid {
		dados.DataAbertura = dataAbertura.String
	}
	if email.Valid {
		dados.Email = email.String
	}
	if telefone1.Valid {
		dados.Telefone1 = telefone1.String
	}
	if telefone2.Valid {
		dados.Telefone2 = telefone2.String
	}

	// Traduz códigos usando dicionários
	dic := database.GetDicionarios()
	if dic != nil {
		if situacao, ok := dic.SituacaoCadastral[dados.SituacaoCadastral]; ok {
			dados.SituacaoCadastral = situacao
		}
		if motivoSituacao.Valid {
			if motivo, ok := dic.MotivoSituacao[motivoSituacao.String]; ok {
				dados.MotivoSituacao = motivo
			}
		}
		if naturezaJuridica.Valid {
			if nat, ok := dic.NaturezaJuridica[naturezaJuridica.String]; ok {
				dados.NaturezaJuridica = nat
			}
		}
		if cnae.Valid {
			if descCnae, ok := dic.CNAE[cnae.String]; ok {
				dados.CNAEPrincipal = descCnae
			}
		}
		if porte, ok := dic.PorteEmpresa[dados.Porte]; ok {
			dados.Porte = porte
		}
	}

	return &dados
}

// BuscaPorNome busca empresas ou sócios por nome
func (s *RedeService) BuscaPorNome(nome string, limite int) ([]models.SearchResult, error) {
	db := database.GetDBSearch()
	if db == nil {
		db = database.GetDBReceita()
	}
	if db == nil {
		return nil, fmt.Errorf("banco de dados não disponível")
	}

	results := make([]models.SearchResult, 0)
	nome = strings.ToUpper(strings.TrimSpace(nome))

	// Busca em empresas
	query := `
		SELECT cnpj, razao_social
		FROM empresas
		WHERE razao_social LIKE ? OR nome_fantasia LIKE ?
		LIMIT ?
	`

	rows, err := db.Query(query, "%"+nome+"%", "%"+nome+"%", limite)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cnpj, razao string
			if err := rows.Scan(&cnpj, &razao); err == nil {
				results = append(results, models.SearchResult{
					ID:    cnpj,
					Label: razao,
					Type:  "PJ",
				})
			}
		}
	}

	// Busca em sócios se ainda tiver espaço
	if len(results) < limite {
		query = `
			SELECT DISTINCT identificador_socio, nome_socio
			FROM socios
			WHERE nome_socio LIKE ?
			LIMIT ?
		`

		rows, err = db.Query(query, "%"+nome+"%", limite-len(results))
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, nomeSocio string
				if err := rows.Scan(&id, &nomeSocio); err == nil {
					results = append(results, models.SearchResult{
						ID:    id,
						Label: nomeSocio,
						Type:  "PF",
					})
				}
			}
		}
	}

	return results, nil
}
