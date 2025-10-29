package crossdata

// EmpresaInfo contém informações de uma empresa
type EmpresaInfo struct {
	CNPJBasico              string
	RazaoSocial             string
	NomeFantasia            string
	SituacaoCadastral       string
	DataSituacaoCadastral   string
	DataInicioAtividades    string
	UF                      string
	Municipio               string
	QualificacaoSocio       string
	QualificacaoDescricao   string
	DataEntradaSociedade    string
}

// SocioInfo contém informações de um sócio
type SocioInfo struct {
	Nome                   string
	CPFCNPJ                string
	Identificador          string
	Qualificacao           string
	QualificacaoDescricao  string
	DataEntrada            string
	Pais                   string
	PaisDescricao          string
	RepresentanteLegal     string
	NomeRepresentante      string
	FaixaEtaria            string
}

// SocioComum representa um sócio compartilhado entre empresas
type SocioComum struct {
	Nome                   string
	CPFCNPJ                string
	Identificador          string
	Qualificacao           string
	QualificacaoDescricao  string
	DataEntradaEmp1        string
	DataEntradaEmp2        string
}

// RedeGrau2 representa empresas de 2º grau
type RedeGrau2 struct {
	CNPJBasico             string
	RazaoSocial            string
	NomeFantasia           string
	SituacaoCadastral      string
	UF                     string
	NomeSocio              string
	CPFSocio               string
	Qualificacao           string
	QualificacaoDescricao  string
}

// EmpresaEndereco representa empresa com endereço
type EmpresaEndereco struct {
	CNPJBasico             string
	RazaoSocial            string
	NomeFantasia           string
	SituacaoCadastral      string
	Logradouro             string
	Numero                 string
	Complemento            string
	Bairro                 string
	CEP                    string
	Municipio              string
	UF                     string
	DataInicioAtividades   string
}

// EmpresaContato representa empresa com contato
type EmpresaContato struct {
	CNPJBasico        string
	RazaoSocial       string
	NomeFantasia      string
	SituacaoCadastral string
	Email             string
	Telefone          string
	UF                string
	TipoMatch         string
}

// RepresentanteLegal representa menor e seu representante
type RepresentanteLegal struct {
	NomeSocio               string
	CPFSocio                string
	FaixaEtaria             string
	CPFRepresentante        string
	NomeRepresentante       string
	QualificacaoRepresentante string
	QualificacaoDescricao   string
	DataEntrada             string
}

// EmpresaEstrangeira representa empresa com sede no exterior
type EmpresaEstrangeira struct {
	CNPJBasico             string
	RazaoSocial            string
	NomeFantasia           string
	SituacaoCadastral      string
	CidadeExterior         string
	CodigoPais             string
	PaisDescricao          string
	DataInicioAtividades   string
	CNAE                   string
}

// SocioEstrangeiro representa sócio estrangeiro
type SocioEstrangeiro struct {
	Nome                   string
	Identificacao          string
	CodigoPais             string
	PaisDescricao          string
	Qualificacao           string
	QualificacaoDescricao  string
	DataEntrada            string
}

// Timeline representa histórico de eventos
type Timeline struct {
	CNPJBasico   string
	RazaoSocial  string
	Eventos      []EventoTimeline
}

// EventoTimeline representa um evento na timeline
type EventoTimeline struct {
	Data      string
	Tipo      string
	Descricao string
}

// EmpresaBaixada representa empresa encerrada
type EmpresaBaixada struct {
	CNPJBasico             string
	RazaoSocial            string
	NomeFantasia           string
	SituacaoCadastral      string
	DataSituacaoCadastral  string
	MotivoSituacao         string
	MotivoDescricao        string
	DataInicioAtividades   string
	UF                     string
	Municipio              string
}

// DadosCompletos representa todos os dados de uma empresa
type DadosCompletos struct {
	Empresa         EmpresaCompleta
	Estabelecimentos []EstabelecimentoCompleto
	Socios          []SocioCompleto
	Simples         *SimplesCompleto
}

// EmpresaCompleta com todos os campos
type EmpresaCompleta struct {
	CNPJBasico                  string
	RazaoSocial                 string
	NaturezaJuridica            string
	NaturezaJuridicaDescricao   string
	QualificacaoResponsavel     string
	QualificacaoDescricao       string
	CapitalSocial               float64
	PorteEmpresa                string
	EnteFederativoResponsavel   string
}

// EstabelecimentoCompleto com todos os campos
type EstabelecimentoCompleto struct {
	CNPJ                      string
	CNPJBasico                string
	CNPJOrdem                 string
	CNPJDV                    string
	MatrizFilial              string
	NomeFantasia              string
	SituacaoCadastral         string
	DataSituacaoCadastral     string
	MotivoSituacaoCadastral   string
	MotivoDescricao           string
	NomeCidadeExterior        string
	Pais                      string
	PaisDescricao             string
	DataInicioAtividades      string
	CNAEFiscal                string
	CNAEFiscalDescricao       string
	CNAEFiscalSecundaria      string
	TipoLogradouro            string
	Logradouro                string
	Numero                    string
	Complemento               string
	Bairro                    string
	CEP                       string
	UF                        string
	Municipio                 string
	MunicipioDescricao        string
	DDD1                      string
	Telefone1                 string
	DDD2                      string
	Telefone2                 string
	DDDFax                    string
	Fax                       string
	CorreioEletronico         string
	SituacaoEspecial          string
	DataSituacaoEspecial      string
}

// SocioCompleto com todos os campos
type SocioCompleto struct {
	CNPJ                            string
	CNPJBasico                      string
	IdentificadorDeSocio            string
	NomeSocio                       string
	CNPJCPFSocio                    string
	QualificacaoSocio               string
	QualificacaoDescricao           string
	DataEntradaSociedade            string
	Pais                            string
	PaisDescricao                   string
	RepresentanteLegal              string
	NomeRepresentante               string
	QualificacaoRepresentanteLegal  string
	QualificacaoRepresentanteDescricao string
	FaixaEtaria                     string
}

// SimplesCompleto com todos os campos
type SimplesCompleto struct {
	CNPJBasico            string
	OpcaoSimples          string
	DataOpcaoSimples      string
	DataExclusaoSimples   string
	OpcaoMEI              string
	DataOpcaoMEI          string
	DataExclusaoMEI       string
}
