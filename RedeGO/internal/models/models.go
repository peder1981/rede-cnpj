package models

// Node representa um nó no grafo (empresa ou pessoa)
type Node struct {
	ID          string                 `json:"id"`
	Label       string                 `json:"label"`
	Type        string                 `json:"tipo"`
	Icon        string                 `json:"icone,omitempty"`
	Color       string                 `json:"cor,omitempty"`
	Fixed       bool                   `json:"fixed,omitempty"`
	Note        string                 `json:"nota,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Flags       []string               `json:"flags,omitempty"`
	Camada      int                    `json:"camada,omitempty"`
	X           float64                `json:"x,omitempty"`
	Y           float64                `json:"y,omitempty"`
}

// Edge representa uma ligação entre nós
type Edge struct {
	From        string                 `json:"de"`
	To          string                 `json:"para"`
	Label       string                 `json:"label,omitempty"`
	Type        string                 `json:"tipo,omitempty"`
	Value       float64                `json:"valor,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Qualificacao string                `json:"qualificacao,omitempty"`
}

// Graph representa o grafo completo
type Graph struct {
	Nodes []Node `json:"no"`
	Edges []Edge `json:"ligacao"`
}

// CNPJData representa dados de uma empresa
type CNPJData struct {
	CNPJ                string  `json:"cnpj"`
	RazaoSocial         string  `json:"razao_social"`
	NomeFantasia        string  `json:"nome_fantasia,omitempty"`
	SituacaoCadastral   string  `json:"situacao_cadastral"`
	DataSituacao        string  `json:"data_situacao,omitempty"`
	MotivoSituacao      string  `json:"motivo_situacao,omitempty"`
	NaturezaJuridica    string  `json:"natureza_juridica,omitempty"`
	CNAEPrincipal       string  `json:"cnae_principal,omitempty"`
	CNAESecundario      string  `json:"cnae_secundario,omitempty"`
	CapitalSocial       float64 `json:"capital_social,omitempty"`
	Porte               string  `json:"porte,omitempty"`
	DataAbertura        string  `json:"data_abertura,omitempty"`
	Logradouro          string  `json:"logradouro,omitempty"`
	Numero              string  `json:"numero,omitempty"`
	Complemento         string  `json:"complemento,omitempty"`
	Bairro              string  `json:"bairro,omitempty"`
	Municipio           string  `json:"municipio,omitempty"`
	UF                  string  `json:"uf,omitempty"`
	CEP                 string  `json:"cep,omitempty"`
	Email               string  `json:"email,omitempty"`
	Telefone1           string  `json:"telefone1,omitempty"`
	Telefone2           string  `json:"telefone2,omitempty"`
	Matriz              bool    `json:"matriz,omitempty"`
}

// SocioData representa dados de um sócio
type SocioData struct {
	CNPJ                string `json:"cnpj"`
	IdentificadorSocio  string `json:"identificador_socio"`
	NomeSocio           string `json:"nome_socio"`
	CPFCNPJSocio        string `json:"cpf_cnpj_socio,omitempty"`
	QualificacaoSocio   string `json:"qualificacao_socio,omitempty"`
	DataEntrada         string `json:"data_entrada,omitempty"`
	Pais                string `json:"pais,omitempty"`
	RepresentanteLegal  string `json:"representante_legal,omitempty"`
	NomeRepresentante   string `json:"nome_representante,omitempty"`
	QualificacaoRepr    string `json:"qualificacao_repr,omitempty"`
	FaixaEtaria         string `json:"faixa_etaria,omitempty"`
}

// SearchResult representa resultado de busca
type SearchResult struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"tipo"`
	Score float64 `json:"score,omitempty"`
}

// CamadaRequest representa requisição para buscar camadas
type CamadaRequest struct {
	IDs               []string `json:"ids"`
	Camada            int      `json:"camada"`
	CriterioCaminhos  string   `json:"criterioCaminhos,omitempty"`
}

// LinkRequest representa requisição para buscar links
type LinkRequest struct {
	IDs          []string `json:"ids"`
	Camada       int      `json:"camada"`
	NumeroItens  int      `json:"numeroItens"`
	ValorMinimo  int      `json:"valorMinimo"`
	ValorMaximo  int      `json:"valorMaximo"`
}

// MapaRequest representa requisição para gerar mapa
type MapaRequest struct {
	Nodes []Node `json:"no"`
}

// ExportRequest representa requisição para exportar dados
type ExportRequest struct {
	Nodes []Node `json:"no"`
	Edges []Edge `json:"ligacao"`
}

// APIResponse representa resposta padrão da API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// FileUploadResponse representa resposta de upload de arquivo
type FileUploadResponse struct {
	NomeArquivoServidor string `json:"nomeArquivoServidor"`
	Mensagem            string `json:"mensagem,omitempty"`
}

// DadosPublicosResponse representa resposta sobre dados públicos disponíveis
type DadosPublicosResponse struct {
	AnoMesSendoUsado   string `json:"ano_mes_sendo_usado"`
	AnoMesDisponivel   string `json:"ano_mes_disponivel"`
	URL                string `json:"url"`
	Mensagem           string `json:"mensagem,omitempty"`
}
